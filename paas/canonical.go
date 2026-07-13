package paas

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ResourceKind is the portable type of a platform object.
type ResourceKind string

const (
	ResourceKindService ResourceKind = "service"
	ResourceKindNetwork ResourceKind = "network"
	ResourceKindVolume  ResourceKind = "volume"
	ResourceKindConfig  ResourceKind = "config"
	ResourceKindSecret  ResourceKind = "secret"
	ResourceKindModel   ResourceKind = "model"
	ResourceKindRoute   ResourceKind = "route"
	ResourceKindPolicy  ResourceKind = "policy"
	ResourceKindRaw     ResourceKind = "raw"
	ResourceKindUnknown ResourceKind = "unknown"
)

// CanonicalApplication is the universal orchestration abstraction layer.
//
// It deliberately stores both normalized resources and raw platform objects:
// the normalized maps make cross-platform planning possible, while Resources
// keeps source-specific Kubernetes, Compose, Swarm, Nomad, and Helm details
// available for loss-aware round trips and later adapter-specific emitters.
type CanonicalApplication struct {
	Name       string                        `json:"name,omitempty" yaml:"name,omitempty"`
	Source     Platform                      `json:"source" yaml:"source"`
	Services   map[string]*Service           `json:"services,omitempty" yaml:"services,omitempty"`
	Networks   map[string]*Network           `json:"networks,omitempty" yaml:"networks,omitempty"`
	Volumes    map[string]*Volume            `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Configs    map[string]*Config            `json:"configs,omitempty" yaml:"configs,omitempty"`
	Secrets    map[string]*Secret            `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Models     map[string]*ComposeModel      `json:"models,omitempty" yaml:"models,omitempty"`
	Mesh       *MeshSpec                     `json:"mesh,omitempty" yaml:"mesh,omitempty"`
	Routes     map[string]*RouteSpec         `json:"routes,omitempty" yaml:"routes,omitempty"`
	Policies   map[string]*PolicySpec        `json:"policies,omitempty" yaml:"policies,omitempty"`
	Resources  map[string]*CanonicalResource `json:"resources,omitempty" yaml:"resources,omitempty"`
	Extensions map[string]interface{}        `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// RouteSpec captures HTTP/TCP routing intent across Traefik, Kubernetes
// Ingress, Gateway-like resources, Nomad service tags, and future emitters.
type RouteSpec struct {
	Name        string                 `json:"name" yaml:"name"`
	Service     string                 `json:"service,omitempty" yaml:"service,omitempty"`
	Protocol    string                 `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Hosts       []string               `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Paths       []string               `json:"paths,omitempty" yaml:"paths,omitempty"`
	Port        string                 `json:"port,omitempty" yaml:"port,omitempty"`
	Entrypoints []string               `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Middlewares []string               `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	TLS         bool                   `json:"tls,omitempty" yaml:"tls,omitempty"`
	Source      Platform               `json:"source" yaml:"source"`
	Raw         interface{}            `json:"raw,omitempty" yaml:"raw,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// PolicySpec captures portable access/rate/placement policy signals.
type PolicySpec struct {
	Name        string                 `json:"name" yaml:"name"`
	Type        string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Target      string                 `json:"target,omitempty" yaml:"target,omitempty"`
	Provider    string                 `json:"provider,omitempty" yaml:"provider,omitempty"`
	Rules       []string               `json:"rules,omitempty" yaml:"rules,omitempty"`
	Source      Platform               `json:"source" yaml:"source"`
	Raw         interface{}            `json:"raw,omitempty" yaml:"raw,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Extension   map[string]string      `json:"extension,omitempty" yaml:"extension,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// CanonicalResource is a single source platform object with normalized identity.
type CanonicalResource struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Kind        ResourceKind           `json:"kind" yaml:"kind"`
	Platform    Platform               `json:"platform" yaml:"platform"`
	Ordinal     int                    `json:"ordinal,omitempty" yaml:"ordinal,omitempty"`
	APIVersion  string                 `json:"api_version,omitempty" yaml:"api_version,omitempty"`
	NativeKind  string                 `json:"native_kind,omitempty" yaml:"native_kind,omitempty"`
	Namespace   string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Raw         interface{}            `json:"raw,omitempty" yaml:"raw,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Canonicalize builds the universal abstraction for an Application.
func Canonicalize(app *Application) *CanonicalApplication {
	canonical := &CanonicalApplication{
		Source:     app.Platform,
		Services:   copyServices(app.Services),
		Networks:   copyNetworks(app.Networks),
		Volumes:    copyVolumes(app.Volumes),
		Configs:    copyConfigs(app.Configs),
		Secrets:    copySecrets(app.Secrets),
		Models:     cloneComposeModels(app.Models),
		Mesh:       applicationMeshForCanonical(app),
		Routes:     map[string]*RouteSpec{},
		Policies:   map[string]*PolicySpec{},
		Resources:  map[string]*CanonicalResource{},
		Extensions: map[string]interface{}{},
	}
	if app.Extensions != nil {
		canonical.Extensions = copyStringInterfaceMap(app.Extensions)
	}

	for name, service := range canonical.Services {
		canonical.AddResource(ResourceKindService, app.Platform, name, "Service", service)
	}
	for name, network := range canonical.Networks {
		canonical.AddResource(ResourceKindNetwork, app.Platform, name, "Network", network)
	}
	for name, volume := range canonical.Volumes {
		canonical.AddResource(ResourceKindVolume, app.Platform, name, "Volume", volume)
	}
	for name, config := range canonical.Configs {
		canonical.AddResource(ResourceKindConfig, app.Platform, name, "Config", config)
	}
	for name, secret := range canonical.Secrets {
		canonical.AddResource(ResourceKindSecret, app.Platform, name, "Secret", secret)
	}
	for name, model := range canonical.Models {
		canonical.AddResource(ResourceKindModel, app.Platform, name, "ComposeModel", model)
	}
	attachRouteAndPolicyResources(canonical, app)
	attachRawExtensionResources(canonical, app)
	attachOpaqueKubernetesManifests(canonical, app)
	return canonical
}

// AttachCanonical refreshes the canonical view on an Application.
func (app *Application) AttachCanonical() {
	app.Canonical = Canonicalize(app)
}

func canonicalForApplication(app *Application) *CanonicalApplication {
	if app == nil {
		return nil
	}
	if app.Canonical != nil {
		return app.Canonical
	}
	return Canonicalize(app)
}

// MergeCanonicalResources copies raw canonical resources and auxiliary
// canonical data from src into dst without overwriting resources that were
// produced by the current target conversion.
func MergeCanonicalResources(dst, src *CanonicalApplication) {
	if dst == nil || src == nil {
		return
	}
	if dst.Extensions == nil {
		dst.Extensions = map[string]interface{}{}
	}
	for key, value := range src.Extensions {
		if _, exists := dst.Extensions[key]; !exists {
			dst.Extensions[key] = deepCopyValue(value)
		}
	}
	if dst.Routes == nil {
		dst.Routes = map[string]*RouteSpec{}
	}
	for name, route := range src.Routes {
		if route == nil {
			continue
		}
		if _, exists := dst.Routes[name]; !exists {
			dst.Routes[name] = clonePortableRouteSpec(route)
		}
	}
	if dst.Policies == nil {
		dst.Policies = map[string]*PolicySpec{}
	}
	for name, policy := range src.Policies {
		if policy == nil {
			continue
		}
		if _, exists := dst.Policies[name]; !exists {
			dst.Policies[name] = clonePortablePolicySpec(policy)
		}
	}
	if dst.Models == nil {
		dst.Models = map[string]*ComposeModel{}
	}
	for name, model := range src.Models {
		if model == nil {
			continue
		}
		if _, exists := dst.Models[name]; !exists {
			dst.Models[name] = cloneComposeModel(model)
		}
	}
	if dst.Mesh == nil && src.Mesh != nil {
		dst.Mesh = cloneMeshSpec(src.Mesh)
	}
	if dst.Resources == nil {
		dst.Resources = map[string]*CanonicalResource{}
	}
	for id, resource := range src.Resources {
		if resource == nil {
			continue
		}
		cloned := cloneCanonicalResourceForConversion(resource)
		if cloned == nil {
			continue
		}
		if _, exists := dst.Resources[id]; !exists {
			dst.Resources[id] = cloned
			continue
		}
		if shadowID := canonicalShadowResourceID(id, dst.Resources); shadowID != "" {
			dst.Resources[shadowID] = cloned
		}
	}
}

func canonicalShadowResourceID(id string, resources map[string]*CanonicalResource) string {
	if id == "" {
		return ""
	}
	shadowID := id + ":source"
	if resources == nil {
		return shadowID
	}
	if _, exists := resources[shadowID]; !exists {
		return shadowID
	}
	for index := 2; index < 1000; index++ {
		candidate := fmt.Sprintf("%s:source-%d", id, index)
		if _, exists := resources[candidate]; !exists {
			return candidate
		}
	}
	return ""
}

func cloneCanonicalResource(resource *CanonicalResource) *CanonicalResource {
	if resource == nil {
		return nil
	}
	cloned := *resource
	cloned.Raw = deepCopyValue(resource.Raw)
	cloned.Metadata = copyStringMap(resource.Metadata)
	cloned.Annotations = copyStringMap(resource.Annotations)
	cloned.Extensions = canonicalizeOpaqueExtensionMap(resource.Extensions)
	return &cloned
}

func cloneCanonicalResourceForConversion(resource *CanonicalResource) *CanonicalResource {
	cloned := cloneCanonicalResource(resource)
	if cloned == nil {
		return nil
	}
	switch cloned.Kind {
	case ResourceKindRoute:
		if route, ok := cloned.Raw.(*RouteSpec); ok {
			cloned.Raw = clonePortableRouteSpec(route)
		} else if route, ok := cloned.Raw.(RouteSpec); ok {
			cloned.Raw = clonePortableRouteSpec(&route)
		}
	case ResourceKindPolicy:
		if policy, ok := cloned.Raw.(*PolicySpec); ok {
			cloned.Raw = clonePortablePolicySpec(policy)
		} else if policy, ok := cloned.Raw.(PolicySpec); ok {
			cloned.Raw = clonePortablePolicySpec(&policy)
		}
	}
	return cloned
}

// canonicalRawResourceValue returns the first raw canonical resource matching
// the requested platform and native kind.
func canonicalRawResourceValue(app *Application, platform Platform, nativeKinds ...string) interface{} {
	canonical := canonicalForApplication(app)
	if canonical == nil {
		return nil
	}
	for _, resource := range canonical.Resources {
		if resource == nil || resource.Platform != platform {
			continue
		}
		if resource.Kind != ResourceKindRaw && resource.Kind != ResourceKindUnknown {
			continue
		}
		if len(nativeKinds) > 0 {
			matched := false
			for _, nativeKind := range nativeKinds {
				if nativeKind == "" {
					continue
				}
				if strings.EqualFold(resource.NativeKind, nativeKind) || strings.EqualFold(resource.Name, nativeKind) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		return deepCopyValue(resource.Raw)
	}
	return nil
}

// AddResource appends or replaces a normalized resource entry.
func (c *CanonicalApplication) AddResource(kind ResourceKind, platform Platform, name, nativeKind string, raw interface{}) {
	if c.Resources == nil {
		c.Resources = map[string]*CanonicalResource{}
	}
	id := canonicalResourceID(platform, kind, name)
	c.Resources[id] = &CanonicalResource{
		ID:         id,
		Name:       name,
		Kind:       kind,
		Platform:   platform,
		NativeKind: nativeKind,
		Raw:        deepCopyValue(raw),
	}
}

// AddRawResource indexes a source-native object without flattening it into the
// portable model. This keeps API identity available for loss-aware round trips
// and adapter-specific emitters.
func (c *CanonicalApplication) AddRawResource(platform Platform, fallbackName string, ordinal int, raw interface{}) {
	if c.Resources == nil {
		c.Resources = map[string]*CanonicalResource{}
	}
	name := fallbackName
	resource := &CanonicalResource{
		Name:     name,
		Kind:     ResourceKindRaw,
		Platform: platform,
		Ordinal:  ordinal,
		Raw:      deepCopyValue(raw),
	}
	if mapped, ok := asMap(resource.Raw); ok {
		resource.NativeKind = toString(mapped["kind"])
		resource.APIVersion = toString(mapped["apiVersion"])
		if metadata, ok := asMap(mapped["metadata"]); ok {
			resource.Namespace = toString(metadata["namespace"])
			if metaName := toString(metadata["name"]); metaName != "" {
				name = canonicalRawResourceName(resource.NativeKind, metaName)
				resource.Name = metaName
			}
			resource.Metadata = toStringMapLoose(metadata["labels"])
			resource.Annotations = toStringMapLoose(metadata["annotations"])
		}
	}
	resource.ID = canonicalResourceID(platform, ResourceKindRaw, name)
	c.Resources[resource.ID] = resource
}

func canonicalRawResourceName(nativeKind, name string) string {
	normalizedKind := strings.ToLower(nativeKind)
	if normalizedKind == "" {
		normalizedKind = "resource"
	}
	return normalizedKind + "-" + name
}

// AddRoute records a route and indexes it as a canonical resource.
func (c *CanonicalApplication) AddRoute(route *RouteSpec) {
	if route == nil || route.Name == "" {
		return
	}
	if c.Routes == nil {
		c.Routes = map[string]*RouteSpec{}
	}
	cloned := cloneRouteSpec(route)
	c.Routes[route.Name] = cloned
	for _, host := range cloned.Hosts {
		if host == "" || host == route.Name {
			continue
		}
		if _, exists := c.Routes[host]; !exists {
			c.Routes[host] = cloneRouteSpec(cloned)
		}
	}
	c.AddResource(ResourceKindRoute, cloned.Source, cloned.Name, "Route", cloned)
}

// AddPolicy records a policy and indexes it as a canonical resource.
func (c *CanonicalApplication) AddPolicy(policy *PolicySpec) {
	if policy == nil || policy.Name == "" {
		return
	}
	if c.Policies == nil {
		c.Policies = map[string]*PolicySpec{}
	}
	cloned := clonePolicySpec(policy)
	c.Policies[policy.Name] = cloned
	c.AddResource(ResourceKindPolicy, cloned.Source, cloned.Name, "Policy", cloned)
}

// RawResources returns raw/non-normalized resources in stable order.
func (c *CanonicalApplication) RawResources() []*CanonicalResource {
	var resources []*CanonicalResource
	for _, resource := range c.Resources {
		if resource.Kind == ResourceKindRaw || resource.Kind == ResourceKindUnknown {
			resources = append(resources, resource)
		}
	}
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].ID < resources[j].ID
	})
	return resources
}

func canonicalResourceID(platform Platform, kind ResourceKind, name string) string {
	return fmt.Sprintf("%s:%s:%s", platform, kind, name)
}

func retagCanonicalApplication(canonical *CanonicalApplication, from, to Platform) {
	if canonical == nil || from == to {
		return
	}
	canonical.Source = to
	if len(canonical.Resources) == 0 {
		return
	}
	retagged := make(map[string]*CanonicalResource, len(canonical.Resources))
	for id, resource := range canonical.Resources {
		if resource == nil {
			retagged[id] = nil
			continue
		}
		if resource.Platform == from {
			resource.Platform = to
		}
		if strings.HasPrefix(id, string(from)+":") {
			id = string(to) + id[len(string(from)):]
		}
		resource.ID = id
		retagged[id] = resource
	}
	canonical.Resources = retagged
}

func attachRawExtensionResources(canonical *CanonicalApplication, app *Application) {
	if rawObjects, ok := applicationExtensionValueForKey(app, "kubernetes.raw"); ok {
		if typed, ok := rawObjects.([]interface{}); ok {
			for index, raw := range typed {
				canonical.AddRawResource(app.Platform, fmt.Sprintf("kubernetes-raw-%d", index), index, raw)
			}
		}
	}

	if rawServices, ok := firstExtensionValue(app.Extensions, "kubernetes.serviceResources", "kubernetes.services", "x-kubernetes-services"); ok {
		if typed, ok := rawServices.([]interface{}); ok {
			for index, raw := range typed {
				canonical.AddRawResource(app.Platform, fmt.Sprintf("kubernetes-service-%d", index), index, raw)
			}
		}
	}
	if rawRoutes, ok := applicationExtensionValueForKey(app, "kubernetes.routes"); ok {
		if typed, ok := rawRoutes.([]interface{}); ok {
			for index, raw := range typed {
				canonical.AddRawResource(app.Platform, fmt.Sprintf("kubernetes-route-%d", index), index, raw)
			}
		}
	}
	if rawPolicies, ok := applicationExtensionValueForKey(app, "kubernetes.policies"); ok {
		if typed, ok := rawPolicies.([]interface{}); ok {
			for index, raw := range typed {
				canonical.AddRawResource(app.Platform, fmt.Sprintf("kubernetes-policy-%d", index), index, raw)
			}
		}
	}
	if rawHPAs, ok := applicationExtensionValueForKey(app, "kubernetes.hpas"); ok {
		if typed, ok := rawHPAs.([]interface{}); ok {
			for index, raw := range typed {
				canonical.AddRawResource(app.Platform, fmt.Sprintf("kubernetes-hpa-%d", index), index, raw)
			}
		}
	}
	if rawPDBs, ok := applicationExtensionValueForKey(app, "kubernetes.pdbs"); ok {
		if typed, ok := rawPDBs.([]interface{}); ok {
			for index, raw := range typed {
				canonical.AddRawResource(app.Platform, fmt.Sprintf("kubernetes-pdb-%d", index), index, raw)
			}
		}
	}
	for _, item := range []struct {
		key            string
		fallbackPrefix string
	}{
		{key: kubernetesNamespacesExtensionKey, fallbackPrefix: "kubernetes-namespace"},
		{key: kubernetesServicesExtensionKey, fallbackPrefix: "kubernetes-service"},
		{key: kubernetesServiceAccountsExtensionKey, fallbackPrefix: "kubernetes-service-account"},
		{key: kubernetesRBACResourcesExtensionKey, fallbackPrefix: "kubernetes-rbac-resource"},
		{key: kubernetesResourceQuotasExtensionKey, fallbackPrefix: "kubernetes-resource-quota"},
		{key: kubernetesLimitRangesExtensionKey, fallbackPrefix: "kubernetes-limit-range"},
		{key: kubernetesPriorityClassesExtensionKey, fallbackPrefix: "kubernetes-priority-class"},
		{key: kubernetesRuntimeClassesExtensionKey, fallbackPrefix: "kubernetes-runtime-class"},
		{key: kubernetesStorageClassesExtensionKey, fallbackPrefix: "kubernetes-storage-class"},
		{key: kubernetesIngressClassesExtensionKey, fallbackPrefix: "kubernetes-ingress-class"},
		{key: kubernetesMutatingWebhooksExtensionKey, fallbackPrefix: "kubernetes-mutating-webhook"},
		{key: kubernetesValidatingWebhooksExtensionKey, fallbackPrefix: "kubernetes-validating-webhook"},
		{key: kubernetesCRDsExtensionKey, fallbackPrefix: "kubernetes-crd"},
		{key: kubernetesCustomResourcesExtensionKey, fallbackPrefix: "kubernetes-custom-resource"},
		{key: kubernetesNetworkPoliciesExtensionKey, fallbackPrefix: "kubernetes-network-policy"},
		{key: kubernetesPersistentVolumesExtensionKey, fallbackPrefix: "kubernetes-persistent-volume"},
		{key: kubernetesPVCsExtensionKey, fallbackPrefix: "kubernetes-persistent-volume-claim"},
		{key: kubernetesWorkloadsExtensionKey, fallbackPrefix: "kubernetes-workload"},
	} {
		if rawItems, ok := applicationExtensionValueForKey(app, item.key); ok {
			if typed, ok := rawItems.([]interface{}); ok {
				for index, raw := range typed {
					canonical.AddRawResource(app.Platform, fmt.Sprintf("%s-%d", item.fallbackPrefix, index), index, raw)
				}
			}
		}
	}
	for _, resource := range canonicalRawResourcesFromBridge(app.Extensions[composeCanonicalRawResourcesExtension]) {
		if resource == nil {
			continue
		}
		canonical.Resources[resource.ID] = resource
	}
	if rawResources, ok := app.Extensions[composeKubernetesRawResourcesExtension].([]interface{}); ok {
		for index, raw := range rawResources {
			canonical.AddRawResource(PlatformKubernetes, fmt.Sprintf("kubernetes-raw-bridge-%d", index), index, raw)
		}
	}

	rawResources, ok := applicationExtensionValueForKey(app, "kubernetes.resources")
	if ok {
		if typed, ok := rawResources.([]interface{}); ok {
			for index, raw := range typed {
				name := fmt.Sprintf("kubernetes-unknown-%d", index)
				resource := &CanonicalResource{
					ID:       canonicalResourceID(app.Platform, ResourceKindUnknown, name),
					Name:     name,
					Kind:     ResourceKindUnknown,
					Platform: app.Platform,
					Ordinal:  index,
					Raw:      deepCopyValue(raw),
				}
				if mapped, ok := asMap(resource.Raw); ok {
					resource.NativeKind = toString(mapped["kind"])
					resource.APIVersion = toString(mapped["apiVersion"])
					if metadata, ok := asMap(mapped["metadata"]); ok {
						resource.Namespace = toString(metadata["namespace"])
						if metaName := toString(metadata["name"]); metaName != "" {
							resource.Name = metaName
							resource.ID = canonicalResourceID(app.Platform, ResourceKindUnknown, metaName)
						}
						resource.Metadata = toStringMapLoose(metadata["labels"])
						resource.Annotations = toStringMapLoose(metadata["annotations"])
					}
				}
				canonical.Resources[resource.ID] = resource
			}
		}
	}
}

func attachOpaqueKubernetesManifests(canonical *CanonicalApplication, app *Application) {
	if canonical == nil || app == nil || len(app.KubernetesOpaqueManifests) == 0 {
		return
	}
	for key, manifest := range app.KubernetesOpaqueManifests {
		if manifest == nil || len(manifest.Raw) == 0 {
			continue
		}
		resource := &CanonicalResource{
			ID:         canonicalResourceID(PlatformKubernetes, ResourceKindRaw, key),
			Name:       key,
			Kind:       ResourceKindRaw,
			Platform:   PlatformKubernetes,
			APIVersion: manifest.APIVersion,
			NativeKind: manifest.Kind,
			Raw:        deepCopyValue(manifest.Raw),
			Extensions: canonicalizeOpaqueExtensionMap(manifest.Extensions),
		}
		if metadata, ok := asMap(manifest.Raw["metadata"]); ok {
			resource.Namespace = toString(metadata["namespace"])
			resource.Metadata = toStringMapLoose(metadata["labels"])
			resource.Annotations = toStringMapLoose(metadata["annotations"])
		}
		resource.ID = canonicalResourceID(PlatformKubernetes, ResourceKindRaw, key)
		canonical.Resources[resource.ID] = resource
	}
}

func attachRouteAndPolicyResources(canonical *CanonicalApplication, app *Application) {
	for name, service := range app.Services {
		for _, route := range routesFromServiceLabels(app.Platform, name, service.Labels) {
			canonical.AddRoute(route)
		}
		for _, policy := range policiesFromServiceLabels(app.Platform, name, service.Labels) {
			canonical.AddPolicy(policy)
		}
	}

	if routes, ok := applicationExtensionValueForKey(app, "kubernetes.routes"); ok {
		if typed, ok := routes.([]interface{}); ok {
			for _, raw := range typed {
				if route := routeFromKubernetesExtension(app, raw); route != nil {
					canonical.AddRoute(route)
				}
			}
		}
	}
	if policies, ok := applicationExtensionValueForKey(app, "kubernetes.policies"); ok {
		if typed, ok := policies.([]interface{}); ok {
			for _, raw := range typed {
				if policy := policyFromKubernetesExtension(app.Platform, raw); policy != nil {
					canonical.AddPolicy(policy)
				}
			}
		}
	}
	for _, route := range routeSpecsFromExtension(app.Extensions[nomadAppRoutesMetaKey]) {
		canonical.AddRoute(route)
		canonical.AddKubernetesRawResourceFromRoute(route)
	}
	for _, route := range routeSpecsFromExtension(app.Extensions[helmAppRoutesFile]) {
		canonical.AddRoute(route)
		canonical.AddKubernetesRawResourceFromRoute(route)
	}
	for _, route := range routeSpecsFromExtension(app.Extensions[composeAppRoutesExtension]) {
		canonical.AddRoute(route)
		canonical.AddKubernetesRawResourceFromRoute(route)
	}
	for _, policy := range policySpecsFromExtension(app.Extensions[nomadAppPoliciesMetaKey]) {
		canonical.AddPolicy(policy)
		canonical.AddKubernetesRawResourceFromPolicy(policy)
	}
	for _, policy := range policySpecsFromExtension(app.Extensions[helmAppPoliciesFile]) {
		canonical.AddPolicy(policy)
		canonical.AddKubernetesRawResourceFromPolicy(policy)
	}
	for _, policy := range policySpecsFromExtension(app.Extensions[composeAppPoliciesExtension]) {
		canonical.AddPolicy(policy)
		canonical.AddKubernetesRawResourceFromPolicy(policy)
	}
}

// AddKubernetesRawResourceFromRoute indexes a Kubernetes-native route payload
// carried inside a portable route sidecar. This lets Compose/Swarm exports
// reload and still restore the original Ingress object.
func (c *CanonicalApplication) AddKubernetesRawResourceFromRoute(route *RouteSpec) {
	if route == nil {
		return
	}
	raw, ok := asMap(route.Raw)
	if !ok || toString(raw["kind"]) == "" {
		return
	}
	c.AddRawResource(PlatformKubernetes, "kubernetes-route-"+route.Name, 0, raw)
}

// AddKubernetesRawResourceFromPolicy indexes a Kubernetes-native policy payload
// carried inside a portable policy sidecar.
func (c *CanonicalApplication) AddKubernetesRawResourceFromPolicy(policy *PolicySpec) {
	if policy == nil {
		return
	}
	raw, ok := asMap(policy.Raw)
	if !ok || toString(raw["kind"]) == "" {
		return
	}
	c.AddRawResource(PlatformKubernetes, "kubernetes-policy-"+policy.Name, 0, raw)
}

func canonicalRawResourcesFromBridge(value interface{}) []*CanonicalResource {
	items, ok := value.([]interface{})
	if !ok || len(items) == 0 {
		return nil
	}
	resources := make([]*CanonicalResource, 0, len(items))
	for _, item := range items {
		mapped, ok := asMap(item)
		if !ok {
			continue
		}
		if encodedRawJSON := toString(mapped["raw_json_base64"]); encodedRawJSON != "" {
			rawJSON, err := base64.StdEncoding.DecodeString(encodedRawJSON)
			if err != nil {
				continue
			}
			var rawValue interface{}
			if err := json.Unmarshal(rawJSON, &rawValue); err != nil {
				continue
			}
			mapped = copyStringInterfaceMap(mapped)
			mapped["raw"] = rawValue
			delete(mapped, "raw_json_base64")
		} else if rawJSON := toString(mapped["raw_json"]); rawJSON != "" {
			var rawValue interface{}
			if err := json.Unmarshal([]byte(rawJSON), &rawValue); err != nil {
				continue
			}
			mapped = copyStringInterfaceMap(mapped)
			mapped["raw"] = rawValue
			delete(mapped, "raw_json")
		}
		var resource CanonicalResource
		raw, err := json.Marshal(mapped)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(raw, &resource); err != nil {
			continue
		}
		if resource.ID == "" || resource.Kind == "" || resource.Platform == "" {
			continue
		}
		resource.Raw = deepCopyValue(resource.Raw)
		resource.Metadata = copyStringMap(resource.Metadata)
		resource.Extensions = canonicalizeOpaqueExtensionMap(resource.Extensions)
		resources = append(resources, &resource)
	}
	return resources
}

func routeSpecsFromExtension(value interface{}) []*RouteSpec {
	switch typed := value.(type) {
	case map[string]*RouteSpec:
		routes := make([]*RouteSpec, 0, len(typed))
		for _, route := range typed {
			routes = append(routes, cloneRouteSpec(route))
		}
		return routes
	case map[string]interface{}:
		return routeSpecsFromExtensionMap(typed)
	case []*RouteSpec:
		routes := make([]*RouteSpec, 0, len(typed))
		for _, route := range typed {
			routes = append(routes, cloneRouteSpec(route))
		}
		return routes
	case []interface{}:
		return routeSpecsFromExtensionSlice(typed)
	default:
		return nil
	}
}

func policySpecsFromExtension(value interface{}) []*PolicySpec {
	switch typed := value.(type) {
	case map[string]*PolicySpec:
		policies := make([]*PolicySpec, 0, len(typed))
		for _, policy := range typed {
			policies = append(policies, clonePolicySpec(policy))
		}
		return policies
	case map[string]interface{}:
		return policySpecsFromExtensionMap(typed)
	case []*PolicySpec:
		policies := make([]*PolicySpec, 0, len(typed))
		for _, policy := range typed {
			policies = append(policies, clonePolicySpec(policy))
		}
		return policies
	case []interface{}:
		return policySpecsFromExtensionSlice(typed)
	default:
		return nil
	}
}

func routeSpecsFromExtensionMap(typed map[string]interface{}) []*RouteSpec {
	if len(typed) == 0 {
		return nil
	}
	routes := make([]*RouteSpec, 0, len(typed))
	for _, raw := range typed {
		if route, ok := raw.(*RouteSpec); ok {
			routes = append(routes, cloneRouteSpec(route))
			continue
		}
		if mapped, ok := asMap(raw); ok {
			if data := routeSpecFromExtensionMap(mapped); data != nil {
				routes = append(routes, data)
			}
		}
	}
	return routes
}

func routeSpecsFromExtensionSlice(typed []interface{}) []*RouteSpec {
	if len(typed) == 0 {
		return nil
	}
	routes := make([]*RouteSpec, 0, len(typed))
	for _, raw := range typed {
		if route, ok := raw.(*RouteSpec); ok {
			routes = append(routes, cloneRouteSpec(route))
			continue
		}
		if mapped, ok := asMap(raw); ok {
			if data := routeSpecFromExtensionMap(mapped); data != nil {
				routes = append(routes, data)
			}
		}
	}
	return routes
}

func policySpecsFromExtensionMap(typed map[string]interface{}) []*PolicySpec {
	if len(typed) == 0 {
		return nil
	}
	policies := make([]*PolicySpec, 0, len(typed))
	for _, raw := range typed {
		if policy, ok := raw.(*PolicySpec); ok {
			policies = append(policies, clonePolicySpec(policy))
			continue
		}
		if mapped, ok := asMap(raw); ok {
			if data := policySpecFromExtensionMap(mapped); data != nil {
				policies = append(policies, data)
			}
		}
	}
	return policies
}

func policySpecsFromExtensionSlice(typed []interface{}) []*PolicySpec {
	if len(typed) == 0 {
		return nil
	}
	policies := make([]*PolicySpec, 0, len(typed))
	for _, raw := range typed {
		if policy, ok := raw.(*PolicySpec); ok {
			policies = append(policies, clonePolicySpec(policy))
			continue
		}
		if mapped, ok := asMap(raw); ok {
			if data := policySpecFromExtensionMap(mapped); data != nil {
				policies = append(policies, data)
			}
		}
	}
	return policies
}

func composeModelsFromExtensionMap(typed map[string]map[string]interface{}) map[string]*ComposeModel {
	if len(typed) == 0 {
		return nil
	}
	models := map[string]*ComposeModel{}
	for name, raw := range typed {
		if data := composeModelFromExtensionMap(raw); data != nil {
			models[name] = data
		}
	}
	if len(models) == 0 {
		return nil
	}
	return models
}

func composeModelFromExtensionMap(mapped map[string]interface{}) *ComposeModel {
	if len(mapped) == 0 {
		return nil
	}
	typed, err := jsonMarshalUnmarshal[ComposeModel](mapped)
	if err != nil || typed == nil {
		return nil
	}
	if typed.Extensions == nil {
		typed.Extensions = map[string]interface{}{}
	}
	for key, value := range mapped {
		switch key {
		case "name", "model", "context_size", "runtime_flags", "extensions":
			continue
		}
		if _, exists := typed.Extensions[key]; !exists {
			typed.Extensions[key] = deepCopyValue(value)
		}
	}
	return typed
}

func jsonMarshalUnmarshal[T any](value interface{}) (*T, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var typed T
	if err := json.Unmarshal(raw, &typed); err != nil {
		return nil, err
	}
	return &typed, nil
}

func copyServices(input map[string]*Service) map[string]*Service {
	output := make(map[string]*Service, len(input))
	for name, service := range input {
		output[name] = cloneService(service)
	}
	return output
}

func copyNetworks(input map[string]*Network) map[string]*Network {
	output := make(map[string]*Network, len(input))
	for name, network := range input {
		output[name] = cloneNetwork(network)
	}
	return output
}

func copyVolumes(input map[string]*Volume) map[string]*Volume {
	output := make(map[string]*Volume, len(input))
	for name, volume := range input {
		output[name] = cloneVolume(volume)
	}
	return output
}

func copyConfigs(input map[string]*Config) map[string]*Config {
	output := make(map[string]*Config, len(input))
	for name, config := range input {
		output[name] = cloneConfig(config)
	}
	return output
}

func copySecrets(input map[string]*Secret) map[string]*Secret {
	output := make(map[string]*Secret, len(input))
	for name, secret := range input {
		output[name] = cloneSecret(secret)
	}
	return output
}
