package paas

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

// PaaS represents the Platform as a Service system
type PaaS struct {
	// Configuration
	config *PaaSConfig
}

// PaaSConfig holds PaaS configuration
type PaaSConfig struct {
	// Default platform for conversions
	DefaultPlatform Platform

	// Working directory for temporary files
	WorkDir string

	// Whether to preserve extensions during conversion
	PreserveExtensions bool
}

// New creates a new PaaS instance
func New(config *PaaSConfig) *PaaS {
	if config == nil {
		config = &PaaSConfig{
			DefaultPlatform:    PlatformDockerCompose,
			WorkDir:            "/tmp/paas",
			PreserveExtensions: true,
		}
	}

	if config.WorkDir == "" {
		config.WorkDir = "/tmp/paas"
	}

	// Create working directory
	os.MkdirAll(config.WorkDir, 0755)

	return &PaaS{
		config: config,
	}
}

// NormalizePlatformAlias maps common platform aliases to their canonical
// platform identifiers.
func NormalizePlatformAlias(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "docker", "compose", "docker compose", "docker-compose", "docker compose spec", "docker compose-spec", "compose-spec":
		return string(PlatformDockerCompose)
	case "swarm", "docker swarm", "docker swarm stack", "stack", "stack-file", "stackfile", "docker-stack", "docker-swarm", "mirantis", "mirantis swarm", "mirantis-swarm", "mirantis stack", "mirantis-stack":
		return string(PlatformDockerSwarm)
	case "nomad", "nomad hcl", "nomad-hcl", "nomad hcl job", "hcl", "jobspec", "nomad-jobspec":
		return string(PlatformNomad)
	case "nomad-json", "nomad json", "json-nomad", "nomadjson":
		return string(PlatformNomad)
	case "kubernetes", "kubernetes yaml", "kubernetes manifests", "kubernetes-manifest", "kubernetes-manifests", "k8s", "k8", "kube":
		return string(PlatformKubernetes)
	case "kubernetes-json", "kubernetes json", "k8s-json", "kube-json":
		return string(PlatformKubernetes)
	case "helm", "helm chart", "helm charts", "helm-chart":
		return string(PlatformHelm)
	default:
		return strings.ToLower(strings.TrimSpace(format))
	}
}

// ParsePlatform parses a user-facing platform string into a canonical
// platform identifier.
func ParsePlatform(platform string) (Platform, error) {
	switch NormalizePlatformAlias(platform) {
	case string(PlatformDockerCompose):
		return PlatformDockerCompose, nil
	case string(PlatformDockerSwarm):
		return PlatformDockerSwarm, nil
	case string(PlatformNomad):
		return PlatformNomad, nil
	case string(PlatformKubernetes):
		return PlatformKubernetes, nil
	case string(PlatformHelm):
		return PlatformHelm, nil
	default:
		return PlatformDockerCompose, fmt.Errorf("unknown platform: %s", platform)
	}
}

// DetectPlatformFromFilename detects a platform from a filename or directory
// path. Helm charts are detected from directories containing Chart.yaml.
func DetectPlatformFromFilename(filename string, defaultPlatform Platform) Platform {
	if info, err := os.Stat(filename); err == nil && info.IsDir() {
		chartYaml := filepath.Join(filename, "Chart.yaml")
		if _, err := os.Stat(chartYaml); err == nil {
			return PlatformHelm
		}
	}

	ext := strings.ToLower(filepath.Ext(filename))
	lowerName := strings.ToLower(filename)

	switch ext {
	case ".yml", ".yaml":
		if strings.Contains(lowerName, "k8s") || strings.Contains(lowerName, "kubernetes") {
			return PlatformKubernetes
		}
		if strings.Contains(lowerName, "swarm") || strings.Contains(lowerName, "stack") ||
			strings.Contains(lowerName, "mirantis-swarm") || strings.Contains(lowerName, "mirantis-stack") {
			return PlatformDockerSwarm
		}
		return PlatformDockerCompose
	case ".json":
		if strings.Contains(lowerName, "nomad") || strings.Contains(lowerName, "job") {
			return PlatformNomad
		}
		if strings.Contains(lowerName, "k8s") || strings.Contains(lowerName, "kubernetes") {
			return PlatformKubernetes
		}
		return defaultPlatform
	case ".hcl", ".nomad":
		return PlatformNomad
	default:
		if defaultPlatform == "" {
			return PlatformDockerCompose
		}
		return defaultPlatform
	}
}

// LoadFile loads an application from a file
func (p *PaaS) LoadFile(filename string) (*Application, error) {
	// Check if it's a Helm chart directory
	if info, err := os.Stat(filename); err == nil && info.IsDir() {
		chartYaml := filepath.Join(filename, "Chart.yaml")
		if _, err := os.Stat(chartYaml); err == nil {
			return ParseHelmChart(filename)
		}
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	platform := DetectPlatformFromFilename(filename, p.config.DefaultPlatform)
	if platform == PlatformDockerCompose || platform == PlatformDockerSwarm {
		if app, err := ParseDockerComposeFile(filename); err == nil {
			if platform == PlatformDockerSwarm || app.Platform != PlatformDockerSwarm {
				retagApplicationForPlatform(app, platform)
			}
			return app, nil
		}
		if app, err := ParseKubernetesYAML(string(content)); err == nil {
			return app, nil
		}
		if app, err := ParseNomadJSON(string(content)); err == nil {
			return app, nil
		}
		return p.LoadContent(string(content), platform)
	}
	if platform == PlatformKubernetes {
		if app, err := ParseKubernetesYAML(string(content)); err == nil {
			return app, nil
		}
		if app, err := ParseDockerCompose(string(content)); err == nil {
			return app, nil
		}
		return p.LoadContent(string(content), platform)
	}
	if platform == PlatformNomad {
		if app, err := ParseNomadHCL(string(content)); err == nil {
			return app, nil
		}
		if app, err := ParseNomadJSON(string(content)); err == nil {
			return app, nil
		}
		return p.LoadContent(string(content), platform)
	}
	if platform == PlatformHelm {
		return p.LoadContent(string(content), platform)
	}
	return p.LoadContent(string(content), platform)
}

// LoadContent loads an application from content string
func (p *PaaS) LoadContent(content string, platform Platform) (*Application, error) {
	switch platform {
	case PlatformDockerCompose, PlatformDockerSwarm:
		app, err := ParseDockerCompose(content)
		if err != nil {
			if k8sApp, k8sErr := ParseKubernetesYAML(content); k8sErr == nil {
				retagApplicationForPlatform(k8sApp, platform)
				return k8sApp, nil
			}
			if nomadApp, nomadErr := ParseNomadJSON(content); nomadErr == nil {
				retagApplicationForPlatform(nomadApp, platform)
				return nomadApp, nil
			}
			return nil, err
		}
		if platform == PlatformDockerSwarm || app.Platform != PlatformDockerSwarm {
			retagApplicationForPlatform(app, platform)
		}
		return app, nil
	case PlatformNomad:
		if app, err := ParseNomadHCL(content); err == nil {
			return app, nil
		}
		return ParseNomadJSON(content)
	case PlatformKubernetes:
		if app, err := ParseKubernetesYAML(content); err == nil {
			return app, nil
		}
		if app, err := ParseDockerCompose(content); err == nil {
			return app, nil
		}
		return ParseNomadJSON(content)
	case PlatformHelm:
		// For Helm, we need a directory, not content
		return nil, fmt.Errorf("Helm charts must be loaded from directories, not content strings")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

// LoadContentFormat loads an application using a user-facing format alias.
// It supports aliases such as "nomad-json" that do not map 1:1 to a single
// Platform enum value.
func (p *PaaS) LoadContentFormat(content, format string) (*Application, error) {
	rawFormat := strings.ToLower(strings.TrimSpace(format))
	switch rawFormat {
	case "nomad-json", "json-nomad", "nomadjson":
		return ParseNomadJSON(content)
	case "kubernetes-json", "k8s-json", "kube-json":
		return ParseKubernetesYAML(kubernetesJSONToYAMLContent(content))
	}
	platform, err := ParsePlatform(format)
	if err != nil && strings.TrimSpace(format) != "" {
		platform = Platform(NormalizePlatformAlias(format))
	}
	return p.LoadContent(content, platform)
}

// LoadFileFormat loads an application from a file using a user-facing format
// alias instead of filename heuristics.
func (p *PaaS) LoadFileFormat(filename, format string) (*Application, error) {
	normalized := NormalizePlatformAlias(format)
	switch normalized {
	case string(PlatformDockerCompose), string(PlatformDockerSwarm):
		app, err := ParseDockerComposeFile(filename)
		if err != nil {
			return nil, err
		}
		retagApplicationForPlatform(app, Platform(normalized))
		return app, nil
	case string(PlatformHelm):
		return ParseHelmChart(filename)
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return p.LoadContentFormat(string(content), format)
}

func retagApplicationForPlatform(app *Application, platform Platform) {
	if app == nil {
		return
	}
	from := app.Platform
	if from == platform {
		return
	}
	app.Platform = platform
	for _, service := range app.Services {
		if service != nil {
			service.Platform = platform
		}
	}
	retagCanonicalApplication(app.Canonical, from, platform)
}

// SaveFile saves an application to a file
func (p *PaaS) SaveFile(app *Application, filename string) error {
	if app != nil && app.Platform == PlatformHelm {
		return SerializeHelmChart(app, filename)
	}
	if app != nil && app.Platform == PlatformNomad && strings.EqualFold(filepath.Ext(filename), ".json") {
		return writeContentToFile(filename, func() (string, error) { return SerializeNomadJSON(app) })
	}
	if app != nil && app.Platform == PlatformKubernetes && strings.EqualFold(filepath.Ext(filename), ".json") {
		return writeContentToFile(filename, func() (string, error) { return SerializeKubernetesJSON(app) })
	}
	platform := DetectPlatformFromFilename(filename, p.config.DefaultPlatform)

	// Special handling for Helm charts
	if platform == PlatformHelm {
		return SerializeHelmChart(app, filename)
	}

	content, err := p.SaveContent(app, platform)
	if err == nil && platform == PlatformNomad && strings.EqualFold(filepath.Ext(filename), ".json") {
		content, err = SerializeNomadJSON(app)
	}
	if err == nil && platform == PlatformKubernetes && strings.EqualFold(filepath.Ext(filename), ".json") {
		content, err = SerializeKubernetesJSON(app)
	}
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func writeContentToFile(filename string, render func() (string, error)) error {
	content, err := render()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

// SaveContent saves an application to a content string
func (p *PaaS) SaveContent(app *Application, platform Platform) (string, error) {
	emitApp := cloneApplication(app)
	switch platform {
	case PlatformDockerCompose, PlatformDockerSwarm:
		return SerializeDockerCompose(emitApp)
	case PlatformNomad:
		return SerializeNomadHCL(emitApp)
	case PlatformKubernetes:
		return SerializeKubernetesYAML(emitApp)
	case PlatformHelm:
		return "", fmt.Errorf("Helm charts must be saved to directories, not content strings")
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
}

// SaveContentFormat saves an application using a user-facing format alias.
// It supports aliases such as "nomad-json" that do not map 1:1 to a single
// Platform enum value.
func (p *PaaS) SaveContentFormat(app *Application, format string) (string, error) {
	rawFormat := strings.ToLower(strings.TrimSpace(format))
	switch rawFormat {
	case "nomad-json", "json-nomad", "nomadjson":
		return SerializeNomadJSON(app)
	case "kubernetes-json", "k8s-json", "kube-json":
		return SerializeKubernetesJSON(app)
	}
	platform, err := ParsePlatform(format)
	if err != nil && strings.TrimSpace(format) != "" {
		platform = Platform(NormalizePlatformAlias(format))
	}
	return p.SaveContent(app, platform)
}

// SaveFileFormat saves an application to a file using a user-facing format
// alias instead of filename heuristics.
func (p *PaaS) SaveFileFormat(app *Application, format, filename string) error {
	if NormalizePlatformAlias(format) == string(PlatformHelm) {
		return SerializeHelmChart(app, filename)
	}
	content, err := p.SaveContentFormat(app, format)
	if err != nil {
		return err
	}
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("filename is required")
	}
	if dir := filepath.Dir(filename); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

// RestoreSource writes the preserved source for a specific format back to disk.
// It is a convenience wrapper over the format-specific restore helpers.
func RestoreSource(app *Application, format, filename string) error {
	rawFormat := strings.ToLower(strings.TrimSpace(format))
	switch rawFormat {
	case "nomad-json", "json-nomad", "nomadjson":
		return RestoreNomadJSONSource(app, filename)
	case "kubernetes-json", "k8s-json", "kube-json":
		return RestoreKubernetesJSONSource(app, filename)
	}
	format = NormalizePlatformAlias(format)
	switch format {
	case string(PlatformDockerCompose), string(PlatformDockerSwarm):
		return RestoreDockerComposeSource(app, filename)
	case string(PlatformNomad):
		return RestoreNomadSource(app, filename)
	case string(PlatformKubernetes):
		return RestoreKubernetesSource(app, filename)
	case string(PlatformHelm):
		return SerializeHelmChart(app, filename)
	default:
		return fmt.Errorf("unsupported restore format: %s", format)
	}
}

// RestoreSourceContent returns the preserved source text for the requested
// format when the format is text-based. Helm charts are directory restores and
// therefore do not have a single content payload.
func RestoreSourceContent(app *Application, format string) (string, error) {
	rawFormat := strings.ToLower(strings.TrimSpace(format))
	switch rawFormat {
	case "nomad-json", "json-nomad", "nomadjson":
		return nomadJSONSourceContent(app)
	case "kubernetes-json", "k8s-json", "kube-json":
		return kubernetesJSONSourceContent(app)
	}
	format = NormalizePlatformAlias(format)
	switch format {
	case string(PlatformDockerCompose), string(PlatformDockerSwarm):
		return dockerComposeSourceContent(app)
	case string(PlatformNomad):
		return nomadSourceContent(app)
	case string(PlatformKubernetes):
		return kubernetesSourceContent(app)
	case string(PlatformHelm):
		return "", fmt.Errorf("helm restore requires a directory path, not content")
	default:
		return "", fmt.Errorf("unsupported restore format: %s", format)
	}
}

// Convert converts an application from one platform to another
func (p *PaaS) Convert(app *Application, from, to Platform) (*Application, error) {
	syncPortableApplicationState(app)
	sourceCanonical := canonicalForApplication(app)

	// If already in target format, return as-is
	if app.Platform == to {
		cloned := cloneApplication(app)
		cloned.AttachCanonical()
		MergeCanonicalResources(cloned.Canonical, sourceCanonical)
		mergeCanonicalModelsIntoApplication(cloned, cloned.Canonical)
		return cloned, nil
	}

	// For now, create a new application with the target platform
	converted := &Application{
		Name:                      app.Name,
		Version:                   app.Version,
		Platform:                  to,
		Namespace:                 applicationNamespace(app),
		Services:                  make(map[string]*Service),
		Networks:                  make(map[string]*Network),
		Volumes:                   make(map[string]*Volume),
		Configs:                   make(map[string]*Config),
		Secrets:                   make(map[string]*Secret),
		KubernetesOpaqueManifests: cloneKubernetesOpaqueManifestSpecs(app.KubernetesOpaqueManifests),
		Models:                    cloneComposeModels(app.Models),
		Includes:                  append([]string{}, app.Includes...),
		IncludeEntries:            cloneInterfaceSlice(app.IncludeEntries),
		Extensions:                convertedExtensions(app.Extensions, from, to),
		SourceFiles:               append([]string{}, app.SourceFiles...),
	}

	// Deep copy services with platform conversion
	for name, service := range app.Services {
		convertedService := cloneService(service)
		convertedService.Platform = to
		convertedService.Extensions = convertedServiceExtensions(convertedService.Extensions, from, to)

		// Convert platform-specific attributes
		if err := p.convertServiceAttributes(convertedService, from, to); err != nil {
			return nil, fmt.Errorf("failed to convert service %s: %w", name, err)
		}

		converted.Services[name] = convertedService
	}

	// Copy other resources (networks, volumes, etc.)
	for name, network := range app.Networks {
		converted.Networks[name] = cloneNetwork(network)
	}

	for name, volume := range app.Volumes {
		converted.Volumes[name] = cloneVolume(volume)
	}

	for name, config := range app.Configs {
		converted.Configs[name] = cloneConfig(config)
	}

	for name, secret := range app.Secrets {
		converted.Secrets[name] = cloneSecret(secret)
	}

	converted.AttachCanonical()
	MergeCanonicalResources(converted.Canonical, sourceCanonical)
	mergeCanonicalModelsIntoApplication(converted, converted.Canonical)
	syncPortableApplicationState(converted)
	if to == PlatformDockerCompose || to == PlatformDockerSwarm {
		stripKubernetesServiceExtensionKeys(converted)
	}
	return converted, nil
}

func applicationModelsForEmit(app *Application) map[string]*ComposeModel {
	if app == nil {
		return nil
	}
	models := cloneComposeModels(app.Models)
	models = mergeComposeModels(models, canonicalModelsForApplication(app))
	return models
}

func applicationNetworksForEmit(app *Application) map[string]*Network {
	if app == nil {
		return nil
	}
	networks := copyNetworks(app.Networks)
	canonical := canonicalForApplication(app)
	if canonical != nil && len(canonical.Networks) > 0 {
		if networks == nil {
			networks = map[string]*Network{}
		}
		for name, network := range canonical.Networks {
			if network == nil {
				continue
			}
			if _, exists := networks[name]; !exists {
				networks[name] = cloneNetwork(network)
			}
		}
	}
	if len(networks) == 0 {
		return nil
	}
	return networks
}

func applicationVolumesForEmit(app *Application) map[string]*Volume {
	if app == nil {
		return nil
	}
	volumes := copyVolumes(app.Volumes)
	canonical := canonicalForApplication(app)
	if canonical != nil && len(canonical.Volumes) > 0 {
		if volumes == nil {
			volumes = map[string]*Volume{}
		}
		for name, volume := range canonical.Volumes {
			if volume == nil {
				continue
			}
			if _, exists := volumes[name]; !exists {
				volumes[name] = cloneVolume(volume)
			}
		}
	}
	if len(volumes) == 0 {
		return nil
	}
	return volumes
}

func applicationConfigsForEmit(app *Application) map[string]*Config {
	if app == nil {
		return nil
	}
	configs := copyConfigs(app.Configs)
	canonical := canonicalForApplication(app)
	if canonical != nil && len(canonical.Configs) > 0 {
		if configs == nil {
			configs = map[string]*Config{}
		}
		for name, config := range canonical.Configs {
			if config == nil {
				continue
			}
			if _, exists := configs[name]; !exists {
				configs[name] = cloneConfig(config)
			}
		}
	}
	if len(configs) == 0 {
		return nil
	}
	return configs
}

func applicationSecretsForEmit(app *Application) map[string]*Secret {
	if app == nil {
		return nil
	}
	secrets := copySecrets(app.Secrets)
	canonical := canonicalForApplication(app)
	if canonical != nil && len(canonical.Secrets) > 0 {
		if secrets == nil {
			secrets = map[string]*Secret{}
		}
		for name, secret := range canonical.Secrets {
			if secret == nil {
				continue
			}
			if _, exists := secrets[name]; !exists {
				secrets[name] = cloneSecret(secret)
			}
		}
	}
	if len(secrets) == 0 {
		return nil
	}
	return secrets
}

func mergeCanonicalModelsIntoApplication(app *Application, canonical *CanonicalApplication) {
	if app == nil {
		return
	}
	models := canonicalModels(canonical)
	if len(models) == 0 {
		return
	}
	if app.Models == nil {
		app.Models = map[string]*ComposeModel{}
	}
	app.Models = mergeComposeModels(app.Models, models)
}

func canonicalModelsForApplication(app *Application) map[string]*ComposeModel {
	if app == nil {
		return nil
	}
	return canonicalModels(canonicalForApplication(app))
}

func canonicalRoutesForApplication(app *Application) map[string]*RouteSpec {
	canonical := canonicalForApplication(app)
	if canonical == nil || len(canonical.Routes) == 0 {
		return nil
	}
	routes := map[string]*RouteSpec{}
	for _, route := range canonical.Routes {
		if route == nil || route.Name == "" {
			continue
		}
		if _, exists := routes[route.Name]; !exists {
			routes[route.Name] = clonePortableRouteSpec(route)
		}
	}
	if len(routes) == 0 {
		return nil
	}
	return routes
}

func canonicalPoliciesForApplication(app *Application) map[string]*PolicySpec {
	canonical := canonicalForApplication(app)
	if canonical == nil || len(canonical.Policies) == 0 {
		return nil
	}
	policies := map[string]*PolicySpec{}
	for name, policy := range canonical.Policies {
		if policy == nil {
			continue
		}
		if _, exists := policies[name]; !exists {
			policies[name] = clonePortablePolicySpec(policy)
		}
	}
	if len(policies) == 0 {
		return nil
	}
	return policies
}

func canonicalKubernetesRawResourcesForApplication(app *Application) []interface{} {
	canonical := canonicalForApplication(app)
	if canonical == nil || len(canonical.Resources) == 0 {
		return nil
	}
	type rawResource struct {
		key      string
		ordinal  int
		resource interface{}
	}
	var resources []rawResource
	for _, resource := range canonical.Resources {
		if resource == nil || resource.Platform != PlatformKubernetes {
			if resource == nil || !kubernetesRawResourceLike(resource.Raw) {
				continue
			}
		}
		if resource == nil {
			continue
		}
		if resource.Kind != ResourceKindRaw && resource.Kind != ResourceKindUnknown {
			continue
		}
		raw, ok := asMap(resource.Raw)
		if !ok || toString(raw["kind"]) == "" {
			continue
		}
		resources = append(resources, rawResource{
			key:      kubernetesDocumentKeyFromMap(raw),
			ordinal:  resource.Ordinal,
			resource: deepCopyValue(raw),
		})
	}
	sort.Slice(resources, func(i, j int) bool {
		if resources[i].ordinal != resources[j].ordinal {
			return resources[i].ordinal < resources[j].ordinal
		}
		return resources[i].key < resources[j].key
	})
	result := make([]interface{}, 0, len(resources))
	seen := map[string]struct{}{}
	for _, resource := range resources {
		if resource.key != "" {
			if _, exists := seen[resource.key]; exists {
				continue
			}
			seen[resource.key] = struct{}{}
		}
		result = append(result, resource.resource)
	}
	return result
}

func kubernetesRawResourceLike(raw interface{}) bool {
	mapped, ok := asMap(raw)
	if !ok {
		return false
	}
	if toString(mapped["kind"]) == "" {
		return false
	}
	if toString(mapped["apiVersion"]) == "" {
		return false
	}
	return true
}

func canonicalRawResourcesForBridge(app *Application, target Platform) []map[string]interface{} {
	canonical := canonicalForApplication(app)
	if canonical == nil || len(canonical.Resources) == 0 {
		return nil
	}
	ids := make([]string, 0, len(canonical.Resources))
	for id := range canonical.Resources {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	resources := make([]map[string]interface{}, 0, len(ids))
	seen := map[string]struct{}{}
	for _, id := range ids {
		resource := canonical.Resources[id]
		if resource == nil {
			continue
		}
		if target != "" && resource.Platform == target {
			continue
		}
		if resource.NativeKind == "compose-go.Project" {
			continue
		}
		cloned := cloneCanonicalResource(resource)
		if cloned == nil || cloned.Raw == nil {
			continue
		}
		rawJSON, err := json.Marshal(cloned.Raw)
		if err != nil {
			continue
		}
		bridgeKey := cloned.ID
		if mapped, ok := asMap(cloned.Raw); ok {
			if key := kubernetesDocumentKeyFromMap(mapped); key != "" {
				bridgeKey = string(cloned.Platform) + ":" + key
			}
		}
		if _, exists := seen[bridgeKey]; exists {
			continue
		}
		seen[bridgeKey] = struct{}{}
		item := map[string]interface{}{
			"id":              cloned.ID,
			"name":            cloned.Name,
			"kind":            string(cloned.Kind),
			"platform":        string(cloned.Platform),
			"native_kind":     cloned.NativeKind,
			"raw_json_base64": base64.StdEncoding.EncodeToString(rawJSON),
		}
		if cloned.Ordinal != 0 {
			item["ordinal"] = cloned.Ordinal
		}
		if cloned.APIVersion != "" {
			item["api_version"] = cloned.APIVersion
		}
		if len(cloned.Metadata) > 0 {
			item["metadata"] = copyStringMap(cloned.Metadata)
		}
		if len(cloned.Annotations) > 0 {
			item["annotations"] = copyStringMap(cloned.Annotations)
		}
		if len(cloned.Extensions) > 0 {
			item["extensions"] = copyStringInterfaceMap(cloned.Extensions)
		}
		resources = append(resources, item)
	}
	if len(resources) == 0 {
		return nil
	}
	return resources
}

func canonicalModels(canonical *CanonicalApplication) map[string]*ComposeModel {
	if canonical == nil || len(canonical.Models) == 0 {
		return nil
	}
	return cloneComposeModels(canonical.Models)
}

func mergeComposeModels(dst map[string]*ComposeModel, incoming map[string]*ComposeModel) map[string]*ComposeModel {
	if len(incoming) == 0 {
		return dst
	}
	if dst == nil {
		dst = map[string]*ComposeModel{}
	}
	for name, model := range incoming {
		if model == nil {
			continue
		}
		if _, exists := dst[name]; !exists {
			dst[name] = cloneComposeModel(model)
		}
	}
	return dst
}

func convertedExtensions(extensions map[string]interface{}, from, to Platform) map[string]interface{} {
	if len(extensions) == 0 {
		return nil
	}
	copied := copyStringInterfaceMap(extensions)
	stripPortableCanonicalExtensionKeys(copied)
	delete(copied, legacyMeshExtensionKey)
	if from == to {
		return copied
	}
	if to != PlatformDockerCompose && to != PlatformDockerSwarm {
		delete(copied, composeRawYAMLExtension)
	}
	if to != PlatformNomad {
		delete(copied, nomadRawHCLExtension)
	}
	if to != PlatformHelm {
		delete(copied, helmRawFilesExtension)
		delete(copied, "chart")
	}
	if to != PlatformDockerCompose && to != PlatformDockerSwarm {
		delete(copied, "x-platform")
	}
	if to == PlatformDockerCompose || to == PlatformDockerSwarm {
		for _, key := range []string{
			legacyMeshExtensionKey,
			"kubernetes.raw",
			"kubernetes.resources",
			"kubernetes.workloads",
			"kubernetes.routes",
			"kubernetes.policies",
			"kubernetes.namespaces",
			"kubernetes.configMaps",
			"kubernetes.secretResources",
			"kubernetes.hpas",
			"kubernetes.pdbs",
			"kubernetes.serviceAccounts",
			"kubernetes.services",
			"kubernetes.serviceResources",
			"kubernetes.ingresses",
			"kubernetes.networkPolicies",
			"kubernetes.persistentVolumes",
			"kubernetes.persistentVolumeClaims",
			"kubernetes.rbac",
			"kubernetes.rbacResources",
			"kubernetes.resourceQuotas",
			"kubernetes.limitRanges",
			"kubernetes.priorityClasses",
			"kubernetes.runtimeClasses",
			"kubernetes.storageClasses",
			"kubernetes.ingressClasses",
			"kubernetes.mutatingWebhookConfigurations",
			"kubernetes.validatingWebhookConfigurations",
			"kubernetes.customResourceDefinitions",
			"kubernetes.customResources",
		} {
			delete(copied, key)
		}
	}
	if to != PlatformKubernetes && to != PlatformHelm && to != PlatformNomad {
		delete(copied, "ingress")
	}
	return copied
}

func applicationNamespace(app *Application) string {
	if app == nil {
		return ""
	}
	if app.Namespace != "" {
		return app.Namespace
	}
	if value, ok := applicationExtensionValueForKey(app, "kubernetes.namespace"); ok {
		return toString(value)
	}
	return ""
}

func applicationExtensionValueForKey(app *Application, key string) (interface{}, bool) {
	if app == nil || app.Extensions == nil {
		return nil, false
	}
	if value, ok := app.Extensions[key]; ok {
		return value, true
	}
	if key == meshExtensionKey || key == legacyMeshExtensionKey {
		if value, ok := app.Extensions[meshExtensionKey]; ok {
			return value, true
		}
		if value, ok := app.Extensions[legacyMeshExtensionKey]; ok {
			return value, true
		}
	}
	if key == "kubernetes.services" {
		if value, ok := app.Extensions["kubernetes.serviceResources"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["x-kubernetes-services"]; ok {
			return value, true
		}
	}
	if key == "kubernetes.serviceResources" {
		if value, ok := app.Extensions["kubernetes.services"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["x-kubernetes-services"]; ok {
			return value, true
		}
	}
	if key == "kubernetes.hpas" || key == "kubernetes.horizontalPodAutoscalers" {
		if value, ok := app.Extensions["kubernetes.hpas"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["kubernetes.horizontalPodAutoscalers"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["x-kubernetes-hpas"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["x-kubernetes-horizontal-pod-autoscalers"]; ok {
			return value, true
		}
	}
	if key == "kubernetes.pdbs" || key == "kubernetes.podDisruptionBudgets" {
		if value, ok := app.Extensions["kubernetes.pdbs"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["kubernetes.podDisruptionBudgets"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["x-kubernetes-pdbs"]; ok {
			return value, true
		}
		if value, ok := app.Extensions["x-kubernetes-pod-disruption-budgets"]; ok {
			return value, true
		}
	}
	if alias := composeApplicationExtensionKey(key); alias != key {
		if value, ok := app.Extensions[alias]; ok {
			return value, true
		}
	}
	return nil, false
}

func composeApplicationCanonicalKey(key string) string {
	if !strings.HasPrefix(key, "x-kubernetes-") {
		if strings.HasPrefix(key, "x-bolabaden-") {
			return "bolabaden." + strings.ReplaceAll(strings.TrimPrefix(key, "x-bolabaden-"), "-", ".")
		}
		if strings.HasPrefix(key, "x-nomad-") {
			return "nomad." + strings.ReplaceAll(strings.TrimPrefix(key, "x-nomad-"), "-", ".")
		}
		return ""
	}
	switch key {
	case "x-kubernetes-namespace":
		return "kubernetes.namespace"
	case "x-kubernetes-namespaces":
		return kubernetesNamespacesExtensionKey
	case "x-kubernetes-workloads":
		return kubernetesWorkloadsExtensionKey
	case "x-kubernetes-config-maps":
		return kubernetesConfigMapsExtensionKey
	case "x-kubernetes-secret-resources":
		return kubernetesSecretsExtensionKey
	case "x-kubernetes-source-documents":
		return kubernetesSourceDocumentsExtensionKey
	case "x-kubernetes-service-targets":
		return "kubernetes.serviceTargets"
	case "x-kubernetes-service-port-targets":
		return "kubernetes.servicePortTargets"
	case "x-kubernetes-node-selector":
		return "kubernetes.nodeSelector"
	case "x-kubernetes-nodeSelector":
		return "kubernetes.nodeSelector"
	case "x-kubernetes-resource-claims":
		return "kubernetes.resourceClaims"
	case "x-kubernetes-resourceClaims":
		return "kubernetes.resourceClaims"
	case "x-kubernetes-readiness-gates":
		return "kubernetes.readinessGates"
	case "x-kubernetes-readinessGates":
		return "kubernetes.readinessGates"
	case "x-kubernetes-init-containers":
		return "kubernetes.initContainers"
	case "x-kubernetes-initContainers":
		return "kubernetes.initContainers"
	case "x-kubernetes-scheduling-gates":
		return "kubernetes.schedulingGates"
	case "x-kubernetes-schedulingGates":
		return "kubernetes.schedulingGates"
	case "x-kubernetes-topology-spread-constraints":
		return "kubernetes.topologySpreadConstraints"
	case "x-kubernetes-topologySpreadConstraints":
		return "kubernetes.topologySpreadConstraints"
	case "x-kubernetes-hostAliases":
		return "kubernetes.hostAliases"
	case "x-kubernetes-host-aliases":
		return "kubernetes.hostAliases"
	case "x-kubernetes-job-pod-failure-policy":
		return "kubernetes.job.podFailurePolicy"
	case "x-kubernetes-job-success-policy":
		return "kubernetes.job.successPolicy"
	case "x-kubernetes-image-pull-secrets":
		return "x-kubernetes-imagePullSecrets"
	case "x-kubernetes-image-pull-policy":
		return "x-kubernetes-imagePullPolicy"
	case "x-kubernetes-dns-policy":
		return "x-kubernetes-dnsPolicy"
	case "x-kubernetes-scheduler-name":
		return "x-kubernetes-schedulerName"
	case "x-kubernetes-termination-message-path":
		return "kubernetes.terminationMessagePath"
	case "x-kubernetes-termination-message-policy":
		return "kubernetes.terminationMessagePolicy"
	case "x-kubernetes-allow-privilege-escalation":
		return "kubernetes.allowPrivilegeEscalation"
	case "x-kubernetes-proc-mount":
		return "kubernetes.procMount"
	case "x-kubernetes-seccompProfile", "x-kubernetes-seccomp-profile":
		return "kubernetes.seccompProfile"
	case "x-kubernetes-se-linux-options":
		return "kubernetes.seLinuxOptions"
	case "x-kubernetes-windows-options":
		return "kubernetes.windowsOptions"
	case "x-kubernetes-fs-group-change-policy":
		return "kubernetes.fsGroupChangePolicy"
	case "x-kubernetes-supplemental-groups-policy":
		return "kubernetes.supplementalGroupsPolicy"
	case "x-kubernetes-binary-data":
		return "kubernetes.binaryData"
	case "x-kubernetes-binaryData":
		return "kubernetes.binaryData"
	case "x-kubernetes-string-data":
		return "kubernetes.stringData"
	case "x-kubernetes-stringData":
		return "kubernetes.stringData"
	case "x-kubernetes-mount-propagation":
		return "kubernetes.mountPropagation"
	case "x-kubernetes-mountPropagation":
		return "kubernetes.mountPropagation"
	case "x-kubernetes-recursive-read-only":
		return "kubernetes.recursiveReadOnly"
	case "x-kubernetes-recursiveReadOnly":
		return "kubernetes.recursiveReadOnly"
	case "x-kubernetes-services":
		return "kubernetes.services"
	case "x-kubernetes-horizontal-pod-autoscalers":
		return "kubernetes.horizontalPodAutoscalers"
	case "x-kubernetes-pod-disruption-budgets":
		return "kubernetes.podDisruptionBudgets"
	case "x-kubernetes-routes":
		return "kubernetes.routes"
	case "x-kubernetes-policies":
		return "kubernetes.policies"
	case "x-kubernetes-resources":
		return "kubernetes.resources"
	case "x-kubernetes-hpas":
		return "kubernetes.hpas"
	case "x-kubernetes-pdbs":
		return "kubernetes.pdbs"
	case "x-kubernetes-network-policies":
		return "kubernetes.networkPolicies"
	case "x-kubernetes-persistent-volumes":
		return "kubernetes.persistentVolumes"
	case "x-kubernetes-persistent-volume-claims":
		return "kubernetes.persistentVolumeClaims"
	case "x-kubernetes-service-accounts":
		return "kubernetes.serviceAccounts"
	case "x-kubernetes-resource-quotas":
		return "kubernetes.resourceQuotas"
	case "x-kubernetes-limit-ranges":
		return "kubernetes.limitRanges"
	case "x-kubernetes-priority-classes":
		return "kubernetes.priorityClasses"
	case "x-kubernetes-runtime-classes":
		return "kubernetes.runtimeClasses"
	case "x-kubernetes-storage-classes":
		return "kubernetes.storageClasses"
	case "x-kubernetes-ingress-classes":
		return "kubernetes.ingressClasses"
	case "x-kubernetes-mutating-webhook-configurations":
		return "kubernetes.mutatingWebhookConfigurations"
	case "x-kubernetes-validating-webhook-configurations":
		return "kubernetes.validatingWebhookConfigurations"
	case "x-kubernetes-custom-resource-definitions":
		return "kubernetes.customResourceDefinitions"
	case "x-kubernetes-custom-resources":
		return "kubernetes.customResources"
	default:
		return "kubernetes." + strings.ReplaceAll(strings.TrimPrefix(key, "x-kubernetes-"), "-", ".")
	}
}

func preserveKubernetesApplicationExtension(key string, to Platform) bool {
	return to == PlatformDockerCompose || to == PlatformDockerSwarm
}

func stripPortableCanonicalExtensionKeys(extensions map[string]interface{}) {
	if len(extensions) == 0 {
		return
	}
	for _, key := range []string{
		"name",
		"version",
		legacyMeshExtensionKey,
		composeAppRoutesExtension,
		composeAppPoliciesExtension,
		composeCanonicalRawResourcesExtension,
		composeKubernetesRawResourcesExtension,
		"kubernetes.serviceResources",
		nomadAppRoutesMetaKey,
		nomadAppPoliciesMetaKey,
		nomadAppNetworksMetaKey,
		nomadAppVolumesMetaKey,
		nomadAppConfigsMetaKey,
		nomadAppSecretsMetaKey,
		nomadCanonicalRawResourcesMetaKey,
		nomadKubernetesRawResourcesMetaKey,
		helmAppRoutesFile,
		helmAppPoliciesFile,
		helmAppNetworksFile,
		helmAppVolumesFile,
		helmAppConfigsFile,
		helmAppSecretsFile,
		helmCanonicalRawResourcesFile,
		helmKubernetesRawResourcesFile,
	} {
		delete(extensions, key)
	}
}

func convertedServiceExtensions(extensions map[string]interface{}, from, to Platform) map[string]interface{} {
	if len(extensions) == 0 {
		return nil
	}
	copied := copyStringInterfaceMap(extensions)
	if to == PlatformDockerCompose || to == PlatformDockerSwarm {
		mapped := make(map[string]interface{}, len(copied))
		for key, value := range copied {
			mapped[composeServiceExtensionKey(key)] = value
		}
		copied = mapped
	}
	if len(copied) == 0 {
		return nil
	}
	return copied
}

func composeServiceExtensionKey(key string) string {
	switch key {
	case "kubernetes.nodeSelector":
		return "x-kubernetes-node-selector"
	case "kubernetes.resourceClaims":
		return "x-kubernetes-resource-claims"
	case "kubernetes.seccompProfile":
		return "x-kubernetes-seccomp-profile"
	case "kubernetes.readinessGates":
		return "x-kubernetes-readiness-gates"
	case "kubernetes.initContainers":
		return "x-kubernetes-init-containers"
	case "kubernetes.schedulingGates":
		return "x-kubernetes-scheduling-gates"
	case "kubernetes.hostAliases":
		return "x-kubernetes-hostAliases"
	case "x-kubernetes-host-aliases":
		return "x-kubernetes-hostAliases"
	case "kubernetes.workload.labels":
		return "x-kubernetes-workload-labels"
	case "kubernetes.workload.annotations":
		return "x-kubernetes-workload-annotations"
	case "kubernetes.template.labels":
		return "x-kubernetes-template-labels"
	case "kubernetes.template.annotations":
		return "x-kubernetes-template-annotations"
	case "kubernetes.topologySpreadConstraints":
		return "x-kubernetes-topology-spread-constraints"
	case "x-kubernetes-nodeSelector":
		return "x-kubernetes-node-selector"
	case "x-kubernetes-resourceClaims":
		return "x-kubernetes-resource-claims"
	case "x-kubernetes-seccompProfile":
		return "x-kubernetes-seccomp-profile"
	case "x-kubernetes-readinessGates":
		return "x-kubernetes-readiness-gates"
	case "x-kubernetes-readiness-gates":
		return "x-kubernetes-readiness-gates"
	case "x-kubernetes-initContainers":
		return "x-kubernetes-init-containers"
	case "x-kubernetes-init-containers":
		return "x-kubernetes-init-containers"
	case "x-kubernetes-schedulingGates":
		return "x-kubernetes-scheduling-gates"
	case "x-kubernetes-scheduling-gates":
		return "x-kubernetes-scheduling-gates"
	case "x-kubernetes-topologySpreadConstraints":
		return "x-kubernetes-topology-spread-constraints"
	case "x-kubernetes-topology-spread-constraints":
		return "x-kubernetes-topology-spread-constraints"
	case "x-kubernetes-hostAliases":
		return "x-kubernetes-hostAliases"
	case "x-kubernetes-job-podFailurePolicy":
		return "x-kubernetes-job-podFailurePolicy"
	case "x-kubernetes-job-successPolicy":
		return "x-kubernetes-job-successPolicy"
	case "kubernetes.job.podFailurePolicy":
		return "x-kubernetes-job-podFailurePolicy"
	case "kubernetes.job.successPolicy":
		return "x-kubernetes-job-successPolicy"
	case "x-kubernetes-priority-class-name":
		return "x-kubernetes-priorityClassName"
	case "x-kubernetes-runtime-class-name":
		return "x-kubernetes-runtimeClassName"
	case "x-kubernetes-node-name":
		return "x-kubernetes-nodeName"
	case "x-kubernetes-sub-domain":
		return "x-kubernetes-subdomain"
	case "x-kubernetes-image-pull-secrets":
		return "x-kubernetes-imagePullSecrets"
	case "x-kubernetes-image-pull-policy":
		return "x-kubernetes-imagePullPolicy"
	case "x-kubernetes-dns-policy":
		return "x-kubernetes-dnsPolicy"
	case "x-kubernetes-scheduler-name":
		return "x-kubernetes-schedulerName"
	case "x-kubernetes-termination-message-path":
		return "x-kubernetes-terminationMessagePath"
	case "x-kubernetes-termination-message-policy":
		return "x-kubernetes-terminationMessagePolicy"
	case "x-kubernetes-allow-privilege-escalation":
		return "x-kubernetes-allowPrivilegeEscalation"
	case "x-kubernetes-proc-mount":
		return "x-kubernetes-procMount"
	case "x-kubernetes-fs-group":
		return "x-kubernetes-fsGroup"
	case "x-kubernetes-fs-group-change-policy":
		return "x-kubernetes-fsGroupChangePolicy"
	case "x-kubernetes-supplemental-groups-policy":
		return "x-kubernetes-supplementalGroupsPolicy"
	case "x-kubernetes-run-as-non-root":
		return "x-kubernetes-runAsNonRoot"
	case "x-kubernetes-active-deadline-seconds":
		return "x-kubernetes-activeDeadlineSeconds"
	case "x-kubernetes-restart-policy":
		return "x-kubernetes-restartPolicy"
	case "x-kubernetes-service-account-name":
		return "x-kubernetes-serviceAccountName"
	case "x-kubernetes-serviceAccountName":
		return "x-kubernetes-serviceAccountName"
	case "x-kubernetes-automount-service-account-token":
		return "x-kubernetes-automountServiceAccountToken"
	case "x-kubernetes-set-hostname-as-fqdn":
		return "x-kubernetes-setHostnameAsFQDN"
	case "x-kubernetes-share-process-namespace":
		return "x-kubernetes-shareProcessNamespace"
	case "x-kubernetes-enable-service-links":
		return "x-kubernetes-enableServiceLinks"
	case "x-kubernetes-host-users":
		return "x-kubernetes-hostUsers"
	case "x-kubernetes-host-network":
		return "x-kubernetes-hostNetwork"
	case "x-kubernetes-host-pid":
		return "x-kubernetes-hostPID"
	case "x-kubernetes-host-ipc":
		return "x-kubernetes-hostIPC"
	default:
		if strings.HasPrefix(key, "x-") {
			return key
		}
		return "x-" + strings.ReplaceAll(key, ".", "-")
	}
}

func preserveKubernetesServiceExtension(key string, to Platform) bool {
	return strings.HasPrefix(key, "kubernetes.") &&
		(to == PlatformDockerCompose || to == PlatformDockerSwarm || to == PlatformNomad || to == PlatformHelm || to == PlatformKubernetes)
}

func cloneService(service *Service) *Service {
	if service == nil {
		return nil
	}
	cloned := *service
	cloned.Ports = clonePortMappings(service.Ports)
	cloned.Networks = append([]string{}, service.Networks...)
	cloned.NetworkAttachments = cloneNetworkAttachments(service.NetworkAttachments)
	cloned.Expose = append([]string{}, service.Expose...)
	cloned.DNS = append([]string{}, service.DNS...)
	cloned.DNSSearch = append([]string{}, service.DNSSearch...)
	cloned.DNSOptions = append([]string{}, service.DNSOptions...)
	cloned.ExtraHosts = append([]string{}, service.ExtraHosts...)
	cloned.HostAliases = cloneHostAliases(service.HostAliases)
	cloned.HostNetworkSet = service.HostNetworkSet
	cloned.Environment = copyStringMap(service.Environment)
	cloned.EnvFile = append([]string{}, service.EnvFile...)
	cloned.EnvFileRefs = cloneEnvFileRefs(service.EnvFileRefs)
	if len(service.EnvSources) > 0 {
		cloned.EnvSources = make([]EnvSource, len(service.EnvSources))
		for i, source := range service.EnvSources {
			cloned.EnvSources[i] = source
			cloned.EnvSources[i].Extensions = cloneMap(source.Extensions)
		}
	}
	if len(service.EnvFrom) > 0 {
		cloned.EnvFrom = make([]EnvFromSource, len(service.EnvFrom))
		for i, source := range service.EnvFrom {
			cloned.EnvFrom[i] = source
			cloned.EnvFrom[i].Extensions = cloneMap(source.Extensions)
		}
	}
	cloned.Command = append([]string{}, service.Command...)
	cloned.Entrypoint = append([]string{}, service.Entrypoint...)
	cloned.Build = cloneBuildConfig(service.Build)
	cloned.Develop = cloneDevelopConfig(service.Develop)
	cloned.Lifecycle = cloneLifecycleHooks(service.Lifecycle)
	cloned.Devices = append([]string{}, service.Devices...)
	cloned.DeviceMappings = cloneDeviceMappings(service.DeviceMappings)
	cloned.Profiles = append([]string{}, service.Profiles...)
	cloned.GroupAdd = append([]string{}, service.GroupAdd...)
	cloned.Sysctls = copyStringMap(service.Sysctls)
	cloned.Volumes = cloneVolumeMounts(service.Volumes)
	cloned.Configs = cloneFileRefs(service.Configs)
	cloned.Secrets = cloneFileRefs(service.Secrets)
	cloned.DependsOn = append([]string{}, service.DependsOn...)
	cloned.Dependencies = cloneDependencySpecs(service.Dependencies)
	cloned.Links = append([]string{}, service.Links...)
	cloned.CapAdd = append([]string{}, service.CapAdd...)
	cloned.CapDrop = append([]string{}, service.CapDrop...)
	cloned.SecurityOpt = append([]string{}, service.SecurityOpt...)
	cloned.Tolerations = cloneTolerations(service.Tolerations)
	cloned.Init = cloneBoolPtr(service.Init)
	cloned.Ulimits = cloneUlimits(service.Ulimits)
	cloned.Deploy = cloneDeploySpec(service.Deploy)
	cloned.Failover = cloneFailoverSpec(service.Failover)
	cloned.Connect = cloneNomadConnectSpec(service.Connect)
	cloned.Spreads = cloneNomadSpreadSpecs(service.Spreads)
	cloned.HealthCheck = cloneHealthCheck(service.HealthCheck)
	cloned.StartupProbe = cloneHealthCheck(service.StartupProbe)
	cloned.SeccompProfile = cloneSeccompProfile(service.SeccompProfile)
	cloned.SELinuxOptions = cloneSELinuxOptions(service.SELinuxOptions)
	cloned.WindowsOptions = cloneWindowsSecurityContextOptions(service.WindowsOptions)
	cloned.AllowPrivilegeEscalation = cloneBoolPtr(service.AllowPrivilegeEscalation)
	cloned.ProcMount = service.ProcMount
	cloned.Affinity = copyStringInterfaceMap(service.Affinity)
	cloned.NodeSelector = copyStringMap(service.NodeSelector)
	cloned.ReadinessGates = cloneMapSlice(service.ReadinessGates)
	cloned.InitContainers = cloneMapSlice(service.InitContainers)
	cloned.ResourceClaims = cloneMapSlice(service.ResourceClaims)
	cloned.EphemeralContainers = cloneMapSlice(service.EphemeralContainers)
	cloned.SchedulingGates = cloneMapSlice(service.SchedulingGates)
	cloned.HostUsers = cloneBoolPtr(service.HostUsers)
	cloned.SupplementalGroups = append([]int64{}, service.SupplementalGroups...)
	cloned.TopologySpreadConstraints = cloneMapSlice(service.TopologySpreadConstraints)
	cloned.Runtime = service.Runtime
	cloned.LogDriver = service.LogDriver
	cloned.LogOpt = copyStringMap(service.LogOpt)
	cloned.LogExtensions = copyStringInterfaceMap(service.LogExtensions)
	cloned.ComposeCompat = cloneComposeCompat(service.ComposeCompat)
	cloned.PidsLimit = service.PidsLimit
	cloned.ShmSize = service.ShmSize
	cloned.PIDMode = service.PIDMode
	cloned.IPCMode = service.IPCMode
	cloned.Labels = copyStringMap(service.Labels)
	cloned.Extensions = canonicalizeServiceExtensionMap(service.Extensions)
	return &cloned
}

func cloneDeviceMappings(devices []DeviceMappingSpec) []DeviceMappingSpec {
	if len(devices) == 0 {
		return nil
	}
	cloned := make([]DeviceMappingSpec, len(devices))
	for i, device := range devices {
		cloned[i] = device
		cloned[i].Extensions = copyStringInterfaceMap(device.Extensions)
	}
	return cloned
}

func cloneNetworkAttachments(attachments map[string]*NetworkAttachment) map[string]*NetworkAttachment {
	if len(attachments) == 0 {
		return nil
	}
	cloned := make(map[string]*NetworkAttachment, len(attachments))
	for name, attachment := range attachments {
		if attachment == nil {
			continue
		}
		copied := *attachment
		copied.Aliases = append([]string{}, attachment.Aliases...)
		copied.LinkLocalIPs = append([]string{}, attachment.LinkLocalIPs...)
		copied.DriverOpts = copyStringMap(attachment.DriverOpts)
		copied.Extensions = copyStringInterfaceMap(attachment.Extensions)
		cloned[name] = &copied
	}
	return cloned
}

func cloneFailoverSpec(failover *FailoverSpec) *FailoverSpec {
	if failover == nil {
		return nil
	}
	cloned := *failover
	cloned.Extensions = copyStringInterfaceMap(failover.Extensions)
	if len(failover.Nodes) > 0 {
		cloned.Nodes = make(map[string]*FailoverNode, len(failover.Nodes))
		for name, node := range failover.Nodes {
			if node == nil {
				continue
			}
			copied := *node
			copied.Extensions = copyStringInterfaceMap(node.Extensions)
			cloned.Nodes[name] = &copied
		}
	}
	return &cloned
}

func cloneNomadConnectSpec(connect *NomadConnectSpec) *NomadConnectSpec {
	if connect == nil {
		return nil
	}
	cloned := *connect
	cloned.Extensions = copyStringInterfaceMap(connect.Extensions)
	cloned.Gateway = cloneMap(connect.Gateway)
	if connect.SidecarService != nil {
		cloned.SidecarService = cloneNomadConnectSidecarService(connect.SidecarService)
	}
	return &cloned
}

func cloneNomadConnectSidecarService(sidecar *NomadConnectSidecarService) *NomadConnectSidecarService {
	if sidecar == nil {
		return nil
	}
	cloned := *sidecar
	cloned.Extensions = copyStringInterfaceMap(sidecar.Extensions)
	cloned.Check = cloneMap(sidecar.Check)
	if sidecar.Proxy != nil {
		cloned.Proxy = cloneNomadConnectProxy(sidecar.Proxy)
	}
	return &cloned
}

func cloneNomadConnectProxy(proxy *NomadConnectProxy) *NomadConnectProxy {
	if proxy == nil {
		return nil
	}
	cloned := *proxy
	cloned.Extensions = copyStringInterfaceMap(proxy.Extensions)
	cloned.Config = cloneMap(proxy.Config)
	if len(proxy.Upstreams) > 0 {
		cloned.Upstreams = make([]NomadConnectUpstream, len(proxy.Upstreams))
		for i, upstream := range proxy.Upstreams {
			cloned.Upstreams[i] = upstream
			cloned.Upstreams[i].Extensions = copyStringInterfaceMap(upstream.Extensions)
		}
	}
	return &cloned
}

func cloneNomadSpreadSpecs(spreads []NomadSpreadSpec) []NomadSpreadSpec {
	if len(spreads) == 0 {
		return nil
	}
	cloned := make([]NomadSpreadSpec, len(spreads))
	for i, spread := range spreads {
		cloned[i] = spread
		cloned[i].Extensions = copyStringInterfaceMap(spread.Extensions)
		if len(spread.Targets) > 0 {
			cloned[i].Targets = make([]NomadSpreadTarget, len(spread.Targets))
			for j, target := range spread.Targets {
				cloned[i].Targets[j] = target
				cloned[i].Targets[j].Extensions = copyStringInterfaceMap(target.Extensions)
			}
		}
	}
	return cloned
}

func cloneHostAliases(aliases []HostAlias) []HostAlias {
	if len(aliases) == 0 {
		return nil
	}
	cloned := make([]HostAlias, 0, len(aliases))
	for _, alias := range aliases {
		cloned = append(cloned, HostAlias{
			IP:         alias.IP,
			Hostnames:  append([]string{}, alias.Hostnames...),
			Extensions: cloneMap(alias.Extensions),
		})
	}
	return cloned
}

func cloneTolerations(tolerations []Toleration) []Toleration {
	if len(tolerations) == 0 {
		return nil
	}
	cloned := make([]Toleration, len(tolerations))
	for i, toleration := range tolerations {
		cloned[i] = toleration
		cloned[i].Extensions = copyStringInterfaceMap(toleration.Extensions)
		if toleration.TolerationSeconds != nil {
			value := *toleration.TolerationSeconds
			cloned[i].TolerationSeconds = &value
		}
	}
	return cloned
}

func mergeTolerations(existing, preserved []Toleration) []Toleration {
	if len(existing) == 0 {
		return cloneTolerations(preserved)
	}
	if len(preserved) == 0 {
		return existing
	}
	result := cloneTolerations(existing)
	for i, toleration := range preserved {
		if i >= len(result) {
			result = append(result, toleration)
			continue
		}
		if len(toleration.Extensions) > 0 {
			if result[i].Extensions == nil {
				result[i].Extensions = map[string]interface{}{}
			}
			for key, value := range toleration.Extensions {
				if _, exists := result[i].Extensions[key]; !exists {
					result[i].Extensions[key] = value
				}
			}
		}
	}
	return result
}

func cloneApplication(app *Application) *Application {
	if app == nil {
		return nil
	}
	cloned := &Application{
		Name:                                    app.Name,
		Version:                                 app.Version,
		Services:                                make(map[string]*Service, len(app.Services)),
		Networks:                                make(map[string]*Network, len(app.Networks)),
		Volumes:                                 make(map[string]*Volume, len(app.Volumes)),
		Configs:                                 make(map[string]*Config, len(app.Configs)),
		Secrets:                                 make(map[string]*Secret, len(app.Secrets)),
		KubernetesServices:                      cloneKubernetesServiceSpecs(app.KubernetesServices),
		KubernetesServiceAccounts:               cloneKubernetesServiceAccountSpecs(app.KubernetesServiceAccounts),
		KubernetesHPAs:                          cloneKubernetesHPASpecs(app.KubernetesHPAs),
		KubernetesPDBs:                          cloneKubernetesPDBSpecs(app.KubernetesPDBs),
		KubernetesResourceQuotas:                cloneKubernetesResourceQuotaSpecs(app.KubernetesResourceQuotas),
		KubernetesLimitRanges:                   cloneKubernetesLimitRangeSpecs(app.KubernetesLimitRanges),
		KubernetesStorageClasses:                cloneKubernetesStorageClassSpecs(app.KubernetesStorageClasses),
		KubernetesIngressClasses:                cloneKubernetesIngressClassSpecs(app.KubernetesIngressClasses),
		KubernetesMutatingWebhookConfigurations: cloneKubernetesWebhookConfigurationSpecs(app.KubernetesMutatingWebhookConfigurations),
		KubernetesValidatingWebhookConfigurations: cloneKubernetesWebhookConfigurationSpecs(app.KubernetesValidatingWebhookConfigurations),
		KubernetesCustomResourceDefinitions:       cloneKubernetesCustomResourceDefinitionSpecs(app.KubernetesCustomResourceDefinitions),
		KubernetesPriorityClasses:                 cloneKubernetesPriorityClassSpecs(app.KubernetesPriorityClasses),
		KubernetesRuntimeClasses:                  cloneKubernetesRuntimeClassSpecs(app.KubernetesRuntimeClasses),
		KubernetesOpaqueManifests:                 cloneKubernetesOpaqueManifestSpecs(app.KubernetesOpaqueManifests),
		Includes:                                  append([]string{}, app.Includes...),
		IncludeEntries:                            cloneInterfaceSlice(app.IncludeEntries),
		Extensions:                                copyStringInterfaceMap(app.Extensions),
		Platform:                                  app.Platform,
		SourceFiles:                               append([]string{}, app.SourceFiles...),
		Mesh:                                      cloneMeshSpec(app.Mesh),
		Namespace:                                 app.Namespace,
	}
	for name, service := range app.Services {
		cloned.Services[name] = cloneService(service)
	}
	for name, network := range app.Networks {
		cloned.Networks[name] = cloneNetwork(network)
	}
	for name, volume := range app.Volumes {
		cloned.Volumes[name] = cloneVolume(volume)
	}
	for name, config := range app.Configs {
		cloned.Configs[name] = cloneConfig(config)
	}
	for name, secret := range app.Secrets {
		cloned.Secrets[name] = cloneSecret(secret)
	}
	cloned.Models = cloneComposeModels(app.Models)
	cloned.Canonical = cloneCanonicalApplication(app.Canonical)
	return cloned
}

func cloneCanonicalApplication(canonical *CanonicalApplication) *CanonicalApplication {
	if canonical == nil {
		return nil
	}
	cloned := &CanonicalApplication{
		Name:       canonical.Name,
		Source:     canonical.Source,
		Services:   copyServices(canonical.Services),
		Networks:   copyNetworks(canonical.Networks),
		Volumes:    copyVolumes(canonical.Volumes),
		Configs:    copyConfigs(canonical.Configs),
		Secrets:    copySecrets(canonical.Secrets),
		Models:     cloneComposeModels(canonical.Models),
		Mesh:       cloneMeshSpec(canonical.Mesh),
		Routes:     cloneRouteSpecMap(canonical.Routes),
		Policies:   clonePolicySpecMap(canonical.Policies),
		Resources:  map[string]*CanonicalResource{},
		Extensions: copyStringInterfaceMap(canonical.Extensions),
	}
	if len(canonical.Resources) > 0 {
		cloned.Resources = make(map[string]*CanonicalResource, len(canonical.Resources))
		for id, resource := range canonical.Resources {
			cloned.Resources[id] = cloneCanonicalResource(resource)
		}
	}
	if len(cloned.Services) == 0 {
		cloned.Services = nil
	}
	if len(cloned.Networks) == 0 {
		cloned.Networks = nil
	}
	if len(cloned.Volumes) == 0 {
		cloned.Volumes = nil
	}
	if len(cloned.Configs) == 0 {
		cloned.Configs = nil
	}
	if len(cloned.Secrets) == 0 {
		cloned.Secrets = nil
	}
	if len(cloned.Models) == 0 {
		cloned.Models = nil
	}
	if len(cloned.Routes) == 0 {
		cloned.Routes = nil
	}
	if len(cloned.Policies) == 0 {
		cloned.Policies = nil
	}
	if len(cloned.Resources) == 0 {
		cloned.Resources = nil
	}
	if len(cloned.Extensions) == 0 {
		cloned.Extensions = nil
	}
	return cloned
}

func cloneKubernetesServiceSpecs(services map[string]*KubernetesServiceSpec) map[string]*KubernetesServiceSpec {
	if len(services) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesServiceSpec, len(services))
	for name, service := range services {
		if service == nil {
			continue
		}
		copied := *service
		copied.Selector = copyStringMap(service.Selector)
		copied.Ports = cloneKubernetesServicePorts(service.Ports)
		copied.LoadBalancerSourceRanges = append([]string{}, service.LoadBalancerSourceRanges...)
		copied.ExternalIPs = append([]string{}, service.ExternalIPs...)
		copied.IPFamilies = append([]string{}, service.IPFamilies...)
		copied.Annotations = copyStringMap(service.Annotations)
		copied.Labels = copyStringMap(service.Labels)
		copied.Extensions = copyStringInterfaceMap(service.Extensions)
		if raw, ok := deepCopyValue(service.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesServicePorts(ports []KubernetesServicePort) []KubernetesServicePort {
	if len(ports) == 0 {
		return nil
	}
	cloned := make([]KubernetesServicePort, len(ports))
	for i, port := range ports {
		cloned[i] = port
		cloned[i].Extensions = copyStringInterfaceMap(port.Extensions)
	}
	return cloned
}

func cloneKubernetesServiceAccountSpecs(accounts map[string]*KubernetesServiceAccountSpec) map[string]*KubernetesServiceAccountSpec {
	if len(accounts) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesServiceAccountSpec, len(accounts))
	for name, account := range accounts {
		if account == nil {
			continue
		}
		copied := *account
		copied.Labels = copyStringMap(account.Labels)
		copied.Annotations = copyStringMap(account.Annotations)
		copied.Secrets = append([]string{}, account.Secrets...)
		copied.ImagePullSecrets = append([]string{}, account.ImagePullSecrets...)
		copied.Extensions = copyStringInterfaceMap(account.Extensions)
		if raw, ok := deepCopyValue(account.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesHPASpecs(hpas map[string]*KubernetesHorizontalPodAutoscalerSpec) map[string]*KubernetesHorizontalPodAutoscalerSpec {
	if len(hpas) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesHorizontalPodAutoscalerSpec, len(hpas))
	for name, hpa := range hpas {
		if hpa == nil {
			continue
		}
		copied := *hpa
		copied.ScaleTarget = copyStringMap(hpa.ScaleTarget)
		copied.Extensions = copyStringInterfaceMap(hpa.Extensions)
		if raw, ok := deepCopyValue(hpa.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		if len(hpa.Metrics) > 0 {
			copied.Metrics = cloneMapSlice(hpa.Metrics)
		}
		if len(hpa.Behavior) > 0 {
			copied.Behavior = copyStringInterfaceMap(hpa.Behavior)
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesPDBSpecs(pdbs map[string]*KubernetesPodDisruptionBudgetSpec) map[string]*KubernetesPodDisruptionBudgetSpec {
	if len(pdbs) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesPodDisruptionBudgetSpec, len(pdbs))
	for name, pdb := range pdbs {
		if pdb == nil {
			continue
		}
		copied := *pdb
		copied.Selector = copyStringMap(pdb.Selector)
		copied.Annotations = copyStringMap(pdb.Annotations)
		copied.Labels = copyStringMap(pdb.Labels)
		copied.Extensions = copyStringInterfaceMap(pdb.Extensions)
		if raw, ok := deepCopyValue(pdb.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesResourceQuotaSpecs(quotas map[string]*KubernetesResourceQuotaSpec) map[string]*KubernetesResourceQuotaSpec {
	if len(quotas) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesResourceQuotaSpec, len(quotas))
	for name, quota := range quotas {
		if quota == nil {
			continue
		}
		copied := *quota
		copied.Scopes = append([]string{}, quota.Scopes...)
		copied.ScopeSelector = copyStringInterfaceMap(quota.ScopeSelector)
		copied.Hard = copyStringInterfaceMap(quota.Hard)
		copied.Annotations = copyStringMap(quota.Annotations)
		copied.Labels = copyStringMap(quota.Labels)
		copied.Extensions = copyStringInterfaceMap(quota.Extensions)
		if raw, ok := deepCopyValue(quota.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesLimitRangeSpecs(ranges map[string]*KubernetesLimitRangeSpec) map[string]*KubernetesLimitRangeSpec {
	if len(ranges) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesLimitRangeSpec, len(ranges))
	for name, limitRange := range ranges {
		if limitRange == nil {
			continue
		}
		copied := *limitRange
		copied.Limits = cloneMapSlice(limitRange.Limits)
		copied.Annotations = copyStringMap(limitRange.Annotations)
		copied.Labels = copyStringMap(limitRange.Labels)
		copied.Extensions = copyStringInterfaceMap(limitRange.Extensions)
		if raw, ok := deepCopyValue(limitRange.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesStorageClassSpecs(classes map[string]*KubernetesStorageClassSpec) map[string]*KubernetesStorageClassSpec {
	if len(classes) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesStorageClassSpec, len(classes))
	for name, class := range classes {
		if class == nil {
			continue
		}
		copied := *class
		copied.Parameters = copyStringMap(class.Parameters)
		copied.MountOptions = append([]string{}, class.MountOptions...)
		copied.AllowedTopologies = cloneMapSlice(class.AllowedTopologies)
		copied.Annotations = copyStringMap(class.Annotations)
		copied.Labels = copyStringMap(class.Labels)
		copied.Extensions = copyStringInterfaceMap(class.Extensions)
		if raw, ok := deepCopyValue(class.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesIngressClassSpecs(classes map[string]*KubernetesIngressClassSpec) map[string]*KubernetesIngressClassSpec {
	if len(classes) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesIngressClassSpec, len(classes))
	for name, class := range classes {
		if class == nil {
			continue
		}
		copied := *class
		copied.Parameters = copyStringInterfaceMap(class.Parameters)
		copied.Annotations = copyStringMap(class.Annotations)
		copied.Labels = copyStringMap(class.Labels)
		copied.Extensions = copyStringInterfaceMap(class.Extensions)
		if raw, ok := deepCopyValue(class.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesWebhookConfigurationSpecs(configs map[string]*KubernetesWebhookConfigurationSpec) map[string]*KubernetesWebhookConfigurationSpec {
	if len(configs) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesWebhookConfigurationSpec, len(configs))
	for name, config := range configs {
		if config == nil {
			continue
		}
		copied := *config
		copied.Webhooks = cloneMapSlice(config.Webhooks)
		copied.Annotations = copyStringMap(config.Annotations)
		copied.Labels = copyStringMap(config.Labels)
		copied.Extensions = copyStringInterfaceMap(config.Extensions)
		if raw, ok := deepCopyValue(config.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesCustomResourceDefinitionSpecs(crds map[string]*KubernetesCustomResourceDefinitionSpec) map[string]*KubernetesCustomResourceDefinitionSpec {
	if len(crds) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesCustomResourceDefinitionSpec, len(crds))
	for name, crd := range crds {
		if crd == nil {
			continue
		}
		copied := *crd
		copied.Names = copyStringInterfaceMap(crd.Names)
		copied.Conversion = copyStringInterfaceMap(crd.Conversion)
		copied.Validation = copyStringInterfaceMap(crd.Validation)
		copied.AdditionalPrinterColumns = cloneMapSlice(crd.AdditionalPrinterColumns)
		copied.Extensions = copyStringInterfaceMap(crd.Extensions)
		if raw, ok := deepCopyValue(crd.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		if len(crd.Versions) > 0 {
			copied.Versions = cloneMapSlice(crd.Versions)
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesPriorityClassSpecs(classes map[string]*KubernetesPriorityClassSpec) map[string]*KubernetesPriorityClassSpec {
	if len(classes) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesPriorityClassSpec, len(classes))
	for name, class := range classes {
		if class == nil {
			continue
		}
		copied := *class
		copied.Extensions = copyStringInterfaceMap(class.Extensions)
		if raw, ok := deepCopyValue(class.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesRuntimeClassSpecs(classes map[string]*KubernetesRuntimeClassSpec) map[string]*KubernetesRuntimeClassSpec {
	if len(classes) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesRuntimeClassSpec, len(classes))
	for name, class := range classes {
		if class == nil {
			continue
		}
		copied := *class
		copied.Overhead = copyStringInterfaceMap(class.Overhead)
		copied.Scheduling = copyStringInterfaceMap(class.Scheduling)
		copied.Extensions = copyStringInterfaceMap(class.Extensions)
		if raw, ok := deepCopyValue(class.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[name] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneKubernetesOpaqueManifestSpecs(manifests map[string]*KubernetesOpaqueManifestSpec) map[string]*KubernetesOpaqueManifestSpec {
	if len(manifests) == 0 {
		return nil
	}
	cloned := make(map[string]*KubernetesOpaqueManifestSpec, len(manifests))
	for key, manifest := range manifests {
		if manifest == nil {
			continue
		}
		copied := *manifest
		copied.Metadata = copyStringInterfaceMap(manifest.Metadata)
		copied.Spec = copyStringInterfaceMap(manifest.Spec)
		copied.Extensions = copyStringInterfaceMap(manifest.Extensions)
		if raw, ok := deepCopyValue(manifest.Raw).(map[string]interface{}); ok {
			copied.Raw = raw
		}
		cloned[key] = &copied
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneMeshSpec(mesh *MeshSpec) *MeshSpec {
	if mesh == nil {
		return nil
	}
	cloned := *mesh
	cloned.ControlPlaneURLs = append([]string{}, mesh.ControlPlaneURLs...)
	cloned.Extensions = copyStringInterfaceMap(mesh.Extensions)
	if len(mesh.Nodes) > 0 {
		cloned.Nodes = make(map[string]*MeshNode, len(mesh.Nodes))
		for name, node := range mesh.Nodes {
			cloned.Nodes[name] = cloneMeshNode(node)
		}
	}
	normalizeMeshSpec(&cloned)
	return &cloned
}

func (mesh *MeshSpec) MarshalJSON() ([]byte, error) {
	if mesh == nil {
		return []byte("null"), nil
	}
	return json.Marshal(serializeMeshSpec(mesh))
}

func (mesh *MeshSpec) UnmarshalJSON(data []byte) error {
	if mesh == nil {
		return fmt.Errorf("nil MeshSpec")
	}
	parsed, err := parseMeshSpecMapBytes(data)
	if err != nil {
		return err
	}
	*mesh = *parsed
	return nil
}

func cloneMeshNode(node *MeshNode) *MeshNode {
	if node == nil {
		return nil
	}
	cloned := *node
	cloned.Extensions = copyStringInterfaceMap(node.Extensions)
	return &cloned
}

func meshSpecFromAny(value interface{}) *MeshSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case *MeshSpec:
		return cloneMeshSpec(typed)
	case MeshSpec:
		return cloneMeshSpec(&typed)
	case map[string]interface{}:
		mesh, err := parseMeshSpecMap(typed)
		if err != nil || mesh == nil {
			return nil
		}
		return mesh
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var mesh MeshSpec
	if err := json.Unmarshal(raw, &mesh); err != nil || isEmptyMeshSpec(&mesh) {
		return nil
	}
	normalizeMeshSpec(&mesh)
	return &mesh
}

func meshSpecToMap(mesh *MeshSpec) map[string]interface{} {
	if mesh == nil || isEmptyMeshSpec(mesh) {
		return nil
	}
	return serializeMeshSpec(mesh)
}

func isEmptyMeshSpec(mesh *MeshSpec) bool {
	return mesh == nil ||
		(!(mesh.EnabledSet || mesh.Enabled) &&
			mesh.Provider == "" &&
			mesh.ControlPlaneURL == "" &&
			len(mesh.ControlPlaneURLs) == 0 &&
			mesh.DNSZone == "" &&
			!(mesh.MagicDNSSet || mesh.MagicDNS) &&
			len(mesh.Nodes) == 0 &&
			len(mesh.Extensions) == 0)
}

func serializeMeshSpec(mesh *MeshSpec) map[string]interface{} {
	if mesh == nil || isEmptyMeshSpec(mesh) {
		return nil
	}
	cloned := cloneMeshSpec(mesh)
	if cloned == nil {
		return nil
	}
	result := map[string]interface{}{}
	if cloned.EnabledSet || cloned.Enabled {
		result["enabled"] = cloned.Enabled
	}
	if cloned.Provider != "" {
		result["provider"] = cloned.Provider
	}
	if cloned.ControlPlaneURL != "" {
		result["control_plane_url"] = cloned.ControlPlaneURL
	}
	if len(cloned.ControlPlaneURLs) > 0 {
		result["control_plane_urls"] = append([]string{}, cloned.ControlPlaneURLs...)
	}
	if cloned.DNSZone != "" {
		result["dns_zone"] = cloned.DNSZone
	}
	if cloned.MagicDNSSet || cloned.MagicDNS {
		result["magic_dns"] = cloned.MagicDNS
	}
	if len(cloned.Nodes) > 0 {
		nodes := map[string]interface{}{}
		for name, node := range cloned.Nodes {
			if node == nil {
				continue
			}
			nodeData := map[string]interface{}{}
			if node.Name != "" {
				nodeData["name"] = node.Name
			}
			if node.Hostname != "" {
				nodeData["hostname"] = node.Hostname
			}
			if node.Address != "" {
				nodeData["address"] = node.Address
			}
			if node.Status != "" {
				nodeData["status"] = node.Status
			}
			if node.Role != "" {
				nodeData["role"] = node.Role
			}
			if node.Priority != 0 {
				nodeData["priority"] = node.Priority
			}
			if node.URL != "" {
				nodeData["url"] = node.URL
			}
			for key, value := range node.Extensions {
				nodeData[key] = deepCopyValue(value)
			}
			nodes[name] = nodeData
		}
		result["nodes"] = nodes
	}
	for key, value := range cloned.Extensions {
		result[key] = deepCopyValue(value)
	}
	return result
}

func parseMeshSpecMapBytes(data []byte) (*MeshSpec, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return parseMeshSpecMap(raw)
}

func parseMeshSpecMap(data map[string]interface{}) (*MeshSpec, error) {
	if len(data) == 0 {
		return nil, nil
	}
	mesh := &MeshSpec{Extensions: map[string]interface{}{}}
	for key, value := range data {
		switch key {
		case "enabled":
			mesh.Enabled = toBool(value)
			mesh.EnabledSet = true
		case "provider":
			mesh.Provider = toString(value)
		case "control_plane_url":
			mesh.ControlPlaneURL = toString(value)
		case "control_plane_urls":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, err
			}
			mesh.ControlPlaneURLs = values
		case "dns_zone":
			mesh.DNSZone = toString(value)
		case "magic_dns":
			mesh.MagicDNS = toBool(value)
			mesh.MagicDNSSet = true
		case "nodes":
			nodesRaw, ok := asMap(value)
			if !ok {
				continue
			}
			mesh.Nodes = map[string]*MeshNode{}
			for name, rawNode := range nodesRaw {
				nodeMap, ok := asMap(rawNode)
				if !ok {
					continue
				}
				mesh.Nodes[name] = &MeshNode{
					Name:       toString(nodeMap["name"]),
					Hostname:   toString(nodeMap["hostname"]),
					Address:    toString(nodeMap["address"]),
					Status:     toString(nodeMap["status"]),
					Role:       toString(nodeMap["role"]),
					Priority:   toInt(nodeMap["priority"]),
					URL:        toString(nodeMap["url"]),
					Extensions: map[string]interface{}{},
				}
				for nodeKey, nodeValue := range nodeMap {
					switch nodeKey {
					case "name", "hostname", "address", "status", "role", "priority", "url":
						continue
					default:
						mesh.Nodes[name].Extensions[nodeKey] = deepCopyValue(nodeValue)
					}
				}
				if len(mesh.Nodes[name].Extensions) == 0 {
					mesh.Nodes[name].Extensions = nil
				}
			}
		default:
			mesh.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(mesh.Extensions) == 0 {
		mesh.Extensions = nil
	}
	if isEmptyMeshSpec(mesh) {
		return nil, nil
	}
	normalizeMeshSpec(mesh)
	return mesh, nil
}

func isEmptyNomadConnectSpec(connect *NomadConnectSpec) bool {
	return connect == nil ||
		(!(connect.NativeSet || connect.Native) &&
			isEmptyNomadConnectSidecarService(connect.SidecarService) &&
			len(connect.Gateway) == 0 &&
			len(connect.Extensions) == 0)
}

func isEmptyNomadConnectSidecarService(sidecar *NomadConnectSidecarService) bool {
	return sidecar == nil ||
		(len(sidecar.Tags) == 0 &&
			len(sidecar.Check) == 0 &&
			isEmptyNomadConnectProxy(sidecar.Proxy) &&
			len(sidecar.Extensions) == 0)
}

func isEmptyNomadConnectProxy(proxy *NomadConnectProxy) bool {
	return proxy == nil ||
		(len(proxy.Upstreams) == 0 &&
			len(proxy.Config) == 0 &&
			len(proxy.Extensions) == 0)
}

func isEmptyNomadSpreadSpec(spread *NomadSpreadSpec) bool {
	return spread == nil ||
		(spread.Attribute == "" &&
			spread.Portable == "" &&
			spread.Weight == 0 &&
			len(spread.Targets) == 0 &&
			len(spread.Extensions) == 0)
}

func applicationMeshForCanonical(app *Application) *MeshSpec {
	if app == nil {
		return nil
	}
	if app.Mesh != nil {
		return cloneMeshSpec(app.Mesh)
	}
	value, ok := applicationExtensionValueForKey(app, meshExtensionKey)
	if !ok {
		return nil
	}
	return meshSpecFromAny(value)
}

func rehydratePortableMesh(app *Application) {
	if app == nil {
		return
	}
	if app.Mesh == nil {
		if value, ok := applicationExtensionValueForKey(app, meshExtensionKey); ok {
			app.Mesh = meshSpecFromAny(value)
		} else if value, ok := applicationExtensionValueForKey(app, legacyMeshExtensionKey); ok {
			app.Mesh = meshSpecFromAny(value)
		}
	}
	if app.Mesh == nil {
		return
	}
	normalizeMeshSpec(app.Mesh)
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	meshMap := meshSpecToMap(app.Mesh)
	if meshMap == nil {
		return
	}
	if _, exists := app.Extensions[meshExtensionKey]; !exists {
		app.Extensions[meshExtensionKey] = deepCopyValue(meshMap)
	}
	delete(app.Extensions, legacyMeshExtensionKey)
	if composeKey := composeApplicationExtensionKey(meshExtensionKey); composeKey != "" {
		if _, exists := app.Extensions[composeKey]; !exists {
			app.Extensions[composeKey] = deepCopyValue(meshMap)
		}
	}
}

func rehydratePortableNomadConnect(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if service.Connect == nil {
			if value, ok := firstExtensionValue(service.Extensions, nomadConnectExtensionKey, "x-nomad-connect"); ok {
				service.Connect = nomadConnectSpecFromAny(value)
			}
		}
		if service.Connect == nil {
			continue
		}
		normalizeNomadConnectSpec(service.Connect)
		if isEmptyNomadConnectSpec(service.Connect) {
			continue
		}
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		connectMap := nomadConnectSpecToMap(service.Connect)
		if connectMap == nil {
			continue
		}
		if app.Platform != PlatformDockerCompose && app.Platform != PlatformDockerSwarm {
			if _, exists := service.Extensions[nomadConnectExtensionKey]; !exists {
				service.Extensions[nomadConnectExtensionKey] = deepCopyValue(connectMap)
			}
		}
		if _, exists := service.Extensions["x-nomad-connect"]; !exists {
			service.Extensions["x-nomad-connect"] = deepCopyValue(connectMap)
		}
	}
}

func rehydratePortableNomadSpread(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if len(service.Spreads) == 0 {
			if value, ok := firstExtensionValue(service.Extensions, nomadSpreadExtensionKey, "x-nomad-spread"); ok {
				service.Spreads = nomadSpreadSpecsFromAny(value)
			}
		}
		if len(service.Spreads) == 0 {
			continue
		}
		service.Spreads = cloneNomadSpreadSpecs(service.Spreads)
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		spreadMap := nomadSpreadSpecsToMap(service.Spreads)
		if spreadMap == nil {
			continue
		}
		if _, exists := service.Extensions[nomadSpreadExtensionKey]; !exists {
			service.Extensions[nomadSpreadExtensionKey] = deepCopyValue(spreadMap)
		}
		if _, exists := service.Extensions["x-nomad-spread"]; !exists {
			service.Extensions["x-nomad-spread"] = deepCopyValue(spreadMap)
		}
	}
}

func rehydratePortableKubernetesServices(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesServices) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, "kubernetes.services", "kubernetes.serviceResources", "x-kubernetes-services"); ok {
			app.KubernetesServices = kubernetesServiceSpecsFromAny(value)
		}
	}
	if len(app.KubernetesServices) == 0 {
		return
	}
	app.KubernetesServices = cloneKubernetesServiceSpecs(app.KubernetesServices)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesServicesExtensionKey, composeKubernetesServicesExtensionKey, app.KubernetesServices)
}

func rehydratePortableKubernetesServiceAccounts(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesServiceAccounts) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesServiceAccountsExtensionKey, composeKubernetesServiceAccountsExtensionKey); ok {
			app.KubernetesServiceAccounts = kubernetesServiceAccountSpecsFromAny(value)
		}
	}
	if len(app.KubernetesServiceAccounts) == 0 {
		return
	}
	app.KubernetesServiceAccounts = cloneKubernetesServiceAccountSpecs(app.KubernetesServiceAccounts)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesServiceAccountsExtensionKey, composeKubernetesServiceAccountsExtensionKey, app.KubernetesServiceAccounts)
}

func rehydratePortableKubernetesServiceExtras(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if len(service.ResourceClaims) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.resourceClaims"); ok {
				if claims, ok := kubernetesMapSliceFromExtension(value); ok && len(claims) > 0 {
					service.ResourceClaims = claims
				}
			} else if value, ok := firstExtensionValue(service.Extensions, "kubernetes.resourceClaims", "x-kubernetes-resource-claims"); ok {
				if claims, ok := kubernetesMapSliceFromExtension(value); ok && len(claims) > 0 {
					service.ResourceClaims = claims
				}
			}
		}
		if len(service.ResourceClaims) > 0 {
			claims := cloneMapSlice(service.ResourceClaims)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.resourceClaims"]; !exists {
				service.Extensions["kubernetes.resourceClaims"] = deepCopyValue(claims)
			}
			if _, exists := service.Extensions["x-kubernetes-resource-claims"]; !exists {
				service.Extensions["x-kubernetes-resource-claims"] = deepCopyValue(claims)
			}
			delete(service.Extensions, "x-kubernetes-resourceClaims")
		}
		if len(service.EphemeralContainers) == 0 {
			if value, ok := firstExtensionValue(service.Extensions, "kubernetes.ephemeralContainers", "x-kubernetes-ephemeral-containers"); ok {
				if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
					service.EphemeralContainers = containers
				}
			}
		}
		if len(service.EphemeralContainers) > 0 {
			containers := cloneMapSlice(service.EphemeralContainers)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.ephemeralContainers"]; !exists {
				service.Extensions["kubernetes.ephemeralContainers"] = deepCopyValue(containers)
			}
			if _, exists := service.Extensions["x-kubernetes-ephemeral-containers"]; !exists {
				service.Extensions["x-kubernetes-ephemeral-containers"] = deepCopyValue(containers)
			}
		}
		if len(service.NodeSelector) == 0 && service.Deploy != nil && service.Deploy.Placement != nil {
			if selector := nodeSelectorFromPlacement(service.Deploy.Placement); len(selector) > 0 {
				service.NodeSelector = selector
			}
		}
		if len(service.NodeSelector) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.nodeSelector"); ok {
				if selector, err := toStringMap(value); err == nil && len(selector) > 0 {
					service.NodeSelector = copyStringMap(selector)
				}
			}
		}
		if len(service.NodeSelector) > 0 {
			selector := copyStringMap(service.NodeSelector)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.nodeSelector"]; !exists {
				service.Extensions["kubernetes.nodeSelector"] = deepCopyValue(selector)
			}
			if _, exists := service.Extensions["x-kubernetes-node-selector"]; !exists {
				service.Extensions["x-kubernetes-node-selector"] = deepCopyValue(selector)
			}
			delete(service.Extensions, "x-kubernetes-nodeSelector")
		}
		if len(service.HostAliases) == 0 {
			if value, ok := firstExtensionValue(service.Extensions, "kubernetes.hostAliases", "x-kubernetes-hostAliases", "x-kubernetes-host-aliases"); ok {
				if aliases := kubernetesHostAliasesFromExtension(value); len(aliases) > 0 {
					service.HostAliases = aliases
				}
			}
		}
		if len(service.HostAliases) > 0 {
			aliases := kubernetesHostAliasesFromService(service)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.hostAliases"]; !exists {
				service.Extensions["kubernetes.hostAliases"] = deepCopyValue(aliases)
			}
			if _, exists := service.Extensions["x-kubernetes-hostAliases"]; !exists {
				service.Extensions["x-kubernetes-hostAliases"] = deepCopyValue(aliases)
			}
			if _, exists := service.Extensions["x-kubernetes-host-aliases"]; !exists {
				service.Extensions["x-kubernetes-host-aliases"] = deepCopyValue(aliases)
			}
			for _, alias := range service.HostAliases {
				for _, hostname := range alias.Hostnames {
					appendUniqueString(&service.ExtraHosts, hostname+"="+alias.IP)
				}
			}
		}
	}
}

func rehydratePortableKubernetesPodPolicy(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if service.AllowPrivilegeEscalation == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.allowPrivilegeEscalation"); ok {
				service.AllowPrivilegeEscalation = boolPtr(value)
			}
		}
		if service.StartupProbe == nil {
			if value, ok := extensionValueForKey(service, "kubernetes.startupProbe"); ok {
				if probe, ok := asMap(value); ok && len(probe) > 0 {
					service.StartupProbe = parseKubernetesProbe(probe)
				}
			}
		}
		if service.ProcMount == "" {
			service.ProcMount = extensionStringValue(service, "kubernetes.procMount")
		}
		if service.RunAsNonRoot == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.runAsNonRoot"); ok {
				service.RunAsNonRoot = boolPtr(value)
			}
		}
		if service.ActiveDeadlineSeconds == nil {
			if value, ok := extensionValueForKey(service, "kubernetes.activeDeadlineSeconds"); ok {
				if seconds := int64(toInt(value)); seconds > 0 {
					service.ActiveDeadlineSeconds = &seconds
				}
			}
		}
		if len(service.SupplementalGroups) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.supplementalGroups"); ok {
				if groups, err := toInt64Slice(value); err == nil && len(groups) > 0 {
					service.SupplementalGroups = append([]int64{}, groups...)
				}
			}
		}
		if service.PodRestartPolicy == "" {
			service.PodRestartPolicy = extensionStringValue(service, "kubernetes.restartPolicy")
		}
		if len(service.Tolerations) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.tolerations"); ok {
				if tolerations, ok := kubernetesTolerationsFromExtension(value); ok && len(tolerations) > 0 {
					service.Tolerations = tolerations
				}
			}
		}
		if service.Lifecycle == nil {
			if value, ok := extensionValueForKey(service, "kubernetes.lifecycle"); ok {
				if lifecycle, ok := asMap(value); ok && len(lifecycle) > 0 {
					parsed := &LifecycleHooks{
						Extensions: map[string]interface{}{
							"kubernetes.lifecycle": cloneMap(lifecycle),
						},
					}
					if postStart, ok := asMap(lifecycle["postStart"]); ok {
						if hook := parseKubernetesLifecycleHook(postStart); hook != nil {
							parsed.PostStart = append(parsed.PostStart, *hook)
						}
					}
					if preStop, ok := asMap(lifecycle["preStop"]); ok {
						if hook := parseKubernetesLifecycleHook(preStop); hook != nil {
							parsed.PreStop = append(parsed.PreStop, *hook)
						}
					}
					if !isEmptyLifecycleHooks(parsed) {
						service.Lifecycle = parsed
					}
				}
			}
		}
	}
}

func rehydratePortableKubernetesHPAs(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesHPAs) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesHPAsExtensionKey, composeKubernetesHPAsExtensionKey); ok {
			app.KubernetesHPAs = kubernetesHPASpecsFromAny(value)
		}
	}
	if len(app.KubernetesHPAs) == 0 {
		return
	}
	app.KubernetesHPAs = cloneKubernetesHPASpecs(app.KubernetesHPAs)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesHPAsExtensionKey, composeKubernetesHPAsExtensionKey, app.KubernetesHPAs)
}

func rehydratePortableKubernetesPDBs(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesPDBs) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesPDBsExtensionKey, composeKubernetesPDBsExtensionKey); ok {
			app.KubernetesPDBs = kubernetesPDBSpecsFromAny(value)
		}
	}
	if len(app.KubernetesPDBs) == 0 {
		return
	}
	app.KubernetesPDBs = cloneKubernetesPDBSpecs(app.KubernetesPDBs)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesPDBsExtensionKey, composeKubernetesPDBsExtensionKey, app.KubernetesPDBs)
}

func rehydratePortableKubernetesResourceQuotas(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesResourceQuotas) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesResourceQuotasExtensionKey, composeKubernetesResourceQuotasExtensionKey); ok {
			app.KubernetesResourceQuotas = kubernetesResourceQuotaSpecsFromAny(value)
		}
	}
	if len(app.KubernetesResourceQuotas) == 0 {
		return
	}
	app.KubernetesResourceQuotas = cloneKubernetesResourceQuotaSpecs(app.KubernetesResourceQuotas)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesResourceQuotasExtensionKey, composeKubernetesResourceQuotasExtensionKey, app.KubernetesResourceQuotas)
}

func rehydratePortableKubernetesLimitRanges(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesLimitRanges) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesLimitRangesExtensionKey, composeKubernetesLimitRangesExtensionKey); ok {
			app.KubernetesLimitRanges = kubernetesLimitRangeSpecsFromAny(value)
		}
	}
	if len(app.KubernetesLimitRanges) == 0 {
		return
	}
	app.KubernetesLimitRanges = cloneKubernetesLimitRangeSpecs(app.KubernetesLimitRanges)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesLimitRangesExtensionKey, composeKubernetesLimitRangesExtensionKey, app.KubernetesLimitRanges)
}

func rehydratePortableKubernetesStorageClasses(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesStorageClasses) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesStorageClassesExtensionKey, composeKubernetesStorageClassesExtensionKey); ok {
			app.KubernetesStorageClasses = kubernetesStorageClassSpecsFromAny(value)
		}
	}
	if len(app.KubernetesStorageClasses) == 0 {
		return
	}
	app.KubernetesStorageClasses = cloneKubernetesStorageClassSpecs(app.KubernetesStorageClasses)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesStorageClassesExtensionKey, composeKubernetesStorageClassesExtensionKey, app.KubernetesStorageClasses)
}

func rehydratePortableKubernetesIngressClasses(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesIngressClasses) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesIngressClassesExtensionKey, composeKubernetesIngressClassesExtensionKey); ok {
			app.KubernetesIngressClasses = kubernetesIngressClassSpecsFromAny(value)
		}
	}
	if len(app.KubernetesIngressClasses) == 0 {
		return
	}
	app.KubernetesIngressClasses = cloneKubernetesIngressClassSpecs(app.KubernetesIngressClasses)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesIngressClassesExtensionKey, composeKubernetesIngressClassesExtensionKey, app.KubernetesIngressClasses)
}

func rehydratePortableKubernetesWebhookConfigurations(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesMutatingWebhookConfigurations) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesMutatingWebhooksExtensionKey, composeKubernetesMutatingWebhooksExtensionKey); ok {
			app.KubernetesMutatingWebhookConfigurations = kubernetesWebhookConfigurationSpecsFromAny(value)
		}
	}
	if len(app.KubernetesMutatingWebhookConfigurations) > 0 {
		app.KubernetesMutatingWebhookConfigurations = cloneKubernetesWebhookConfigurationSpecs(app.KubernetesMutatingWebhookConfigurations)
		propagateTypedKubernetesResourcesToExtensions(app, kubernetesMutatingWebhooksExtensionKey, composeKubernetesMutatingWebhooksExtensionKey, app.KubernetesMutatingWebhookConfigurations)
	}
	if len(app.KubernetesValidatingWebhookConfigurations) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesValidatingWebhooksExtensionKey, composeKubernetesValidatingWebhooksExtensionKey); ok {
			app.KubernetesValidatingWebhookConfigurations = kubernetesWebhookConfigurationSpecsFromAny(value)
		}
	}
	if len(app.KubernetesValidatingWebhookConfigurations) > 0 {
		app.KubernetesValidatingWebhookConfigurations = cloneKubernetesWebhookConfigurationSpecs(app.KubernetesValidatingWebhookConfigurations)
		propagateTypedKubernetesResourcesToExtensions(app, kubernetesValidatingWebhooksExtensionKey, composeKubernetesValidatingWebhooksExtensionKey, app.KubernetesValidatingWebhookConfigurations)
	}
}

func rehydratePortableKubernetesCustomResourceDefinitions(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesCustomResourceDefinitions) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesCRDsExtensionKey, composeKubernetesCRDsExtensionKey); ok {
			app.KubernetesCustomResourceDefinitions = kubernetesCustomResourceDefinitionSpecsFromAny(value)
		}
	}
	if len(app.KubernetesCustomResourceDefinitions) == 0 {
		return
	}
	app.KubernetesCustomResourceDefinitions = cloneKubernetesCustomResourceDefinitionSpecs(app.KubernetesCustomResourceDefinitions)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesCRDsExtensionKey, composeKubernetesCRDsExtensionKey, app.KubernetesCustomResourceDefinitions)
}

func rehydratePortableKubernetesPriorityClasses(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesPriorityClasses) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesPriorityClassesExtensionKey, composeKubernetesPriorityClassesExtensionKey); ok {
			app.KubernetesPriorityClasses = kubernetesPriorityClassSpecsFromAny(value)
		}
	}
	if len(app.KubernetesPriorityClasses) == 0 {
		return
	}
	app.KubernetesPriorityClasses = cloneKubernetesPriorityClassSpecs(app.KubernetesPriorityClasses)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesPriorityClassesExtensionKey, composeKubernetesPriorityClassesExtensionKey, app.KubernetesPriorityClasses)
}

func rehydratePortableKubernetesRuntimeClasses(app *Application) {
	if app == nil {
		return
	}
	if len(app.KubernetesRuntimeClasses) == 0 {
		if value, ok := firstExtensionValue(app.Extensions, kubernetesRuntimeClassesExtensionKey, composeKubernetesRuntimeClassesExtensionKey); ok {
			app.KubernetesRuntimeClasses = kubernetesRuntimeClassSpecsFromAny(value)
		}
	}
	if len(app.KubernetesRuntimeClasses) == 0 {
		return
	}
	app.KubernetesRuntimeClasses = cloneKubernetesRuntimeClassSpecs(app.KubernetesRuntimeClasses)
	propagateTypedKubernetesResourcesToExtensions(app, kubernetesRuntimeClassesExtensionKey, composeKubernetesRuntimeClassesExtensionKey, app.KubernetesRuntimeClasses)
}

func rehydratePortableKubernetesOpaqueManifests(app *Application) {
	if app == nil {
		return
	}
	combined := map[string]*KubernetesOpaqueManifestSpec{}
	for key, spec := range app.KubernetesOpaqueManifests {
		if spec == nil {
			continue
		}
		combined[key] = cloneKubernetesOpaqueManifestSpecs(map[string]*KubernetesOpaqueManifestSpec{key: spec})[key]
	}
	if canonical := canonicalForApplication(app); canonical != nil && len(canonical.Resources) > 0 {
		for key, spec := range kubernetesOpaqueManifestSpecsFromCanonical(canonical) {
			if spec == nil {
				continue
			}
			if _, exists := combined[key]; !exists {
				combined[key] = cloneKubernetesOpaqueManifestSpecs(map[string]*KubernetesOpaqueManifestSpec{key: spec})[key]
			}
		}
	}
	resourceSpecs := kubernetesOpaqueManifestSpecsFromResources(
		app.Extensions[kubernetesResourcesExtensionKey],
		app.Extensions[composeKubernetesResourcesExtensionKey],
		app.Extensions["kubernetes.raw"],
		app.Extensions[kubernetesWorkloadsExtensionKey],
		app.Extensions[composeKubernetesWorkloadsExtensionKey],
		app.Extensions[kubernetesServicesExtensionKey],
		app.Extensions["kubernetes.serviceResources"],
		app.Extensions[composeKubernetesServicesExtensionKey],
		app.Extensions[kubernetesHPAsExtensionKey],
		app.Extensions[composeKubernetesHPAsExtensionKey],
		app.Extensions[kubernetesPDBsExtensionKey],
		app.Extensions[composeKubernetesPDBsExtensionKey],
		app.Extensions["kubernetes.namespaces"],
		app.Extensions[composeKubernetesNamespacesExtensionKey],
		app.Extensions[kubernetesConfigMapsExtensionKey],
		app.Extensions[composeKubernetesConfigMapsExtensionKey],
		app.Extensions[kubernetesSecretsExtensionKey],
		app.Extensions[composeKubernetesSecretsExtensionKey],
		app.Extensions[kubernetesServiceAccountsExtensionKey],
		app.Extensions[composeKubernetesServiceAccountsExtensionKey],
		app.Extensions[kubernetesIngressesExtensionKey],
		app.Extensions[composeKubernetesIngressesExtensionKey],
		app.Extensions[kubernetesNetworkPoliciesExtensionKey],
		app.Extensions[composeKubernetesNetworkPoliciesExtensionKey],
		app.Extensions[kubernetesPersistentVolumesExtensionKey],
		app.Extensions[composeKubernetesPersistentVolumesExtensionKey],
		app.Extensions[kubernetesPVCsExtensionKey],
		app.Extensions[composeKubernetesPVCsExtensionKey],
		app.Extensions[kubernetesRBACResourcesExtensionKey],
		app.Extensions[composeKubernetesRBACResourcesExtensionKey],
		app.Extensions[kubernetesResourceQuotasExtensionKey],
		app.Extensions[composeKubernetesResourceQuotasExtensionKey],
		app.Extensions[kubernetesLimitRangesExtensionKey],
		app.Extensions[composeKubernetesLimitRangesExtensionKey],
		app.Extensions[kubernetesPriorityClassesExtensionKey],
		app.Extensions[composeKubernetesPriorityClassesExtensionKey],
		app.Extensions[kubernetesRuntimeClassesExtensionKey],
		app.Extensions[composeKubernetesRuntimeClassesExtensionKey],
		app.Extensions[kubernetesStorageClassesExtensionKey],
		app.Extensions[composeKubernetesStorageClassesExtensionKey],
		app.Extensions[kubernetesIngressClassesExtensionKey],
		app.Extensions[composeKubernetesIngressClassesExtensionKey],
		app.Extensions[kubernetesMutatingWebhooksExtensionKey],
		app.Extensions[composeKubernetesMutatingWebhooksExtensionKey],
		app.Extensions[kubernetesValidatingWebhooksExtensionKey],
		app.Extensions[composeKubernetesValidatingWebhooksExtensionKey],
		app.Extensions[kubernetesCRDsExtensionKey],
		app.Extensions[composeKubernetesCRDsExtensionKey],
		app.Extensions[kubernetesCustomResourcesExtensionKey],
		app.Extensions[composeKubernetesCustomResourcesExtensionKey],
	)
	for key, spec := range resourceSpecs {
		if spec == nil {
			continue
		}
		if _, exists := combined[key]; !exists {
			combined[key] = spec
		}
	}
	if len(combined) == 0 {
		return
	}
	app.KubernetesOpaqueManifests = cloneKubernetesOpaqueManifestSpecs(combined)
}

func syncPortableApplicationState(app *Application) {
	if app == nil {
		return
	}
	rehydratePortableMesh(app)
	rehydratePortableFailover(app)
	rehydratePortableNomadConnect(app)
	rehydratePortableNomadSpread(app)
	rehydratePortableNomadScheduler(app)
	rehydratePortableKubernetesServices(app)
	reconcileKubernetesServices(app)
	rehydratePortableKubernetesServiceAccounts(app)
	reconcileKubernetesServiceAccounts(app)
	rehydratePortableKubernetesServiceExtras(app)
	rehydratePortableKubernetesPodPolicy(app)
	reconcileKubernetesRBACResources(app)
	rehydratePortableKubernetesHPAs(app)
	rehydratePortableKubernetesPDBs(app)
	reconcileKubernetesHorizontalPodAutoscalers(app)
	reconcileKubernetesPodDisruptionBudgets(app)
	rehydratePortableKubernetesResourceQuotas(app)
	rehydratePortableKubernetesLimitRanges(app)
	rehydratePortableKubernetesStorageClasses(app)
	rehydratePortableKubernetesIngressClasses(app)
	rehydratePortableKubernetesWebhookConfigurations(app)
	rehydratePortableKubernetesCustomResourceDefinitions(app)
	rehydratePortableKubernetesPriorityClasses(app)
	rehydratePortableKubernetesRuntimeClasses(app)
	rehydratePortableKubernetesServiceIdentity(app)
	rehydratePortableKubernetesSecurityContext(app)
	rehydratePortableKubernetesOpaqueManifests(app)
	rehydratePortableKubernetesWindowsSecurityContext(app)
	rehydratePortableKubernetesHostNetworkDNSPolicy(app)
}

func rehydratePortableFailover(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		hydrateServiceFailoverFromExtensions(service)
	}
}

func rehydratePortableKubernetesWindowsSecurityContext(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil || service.WindowsOptions == nil || service.WindowsOptions.HostProcess == nil || !*service.WindowsOptions.HostProcess {
			continue
		}
		if !service.HostNetwork {
			service.HostNetwork = true
		}
		service.HostNetworkSet = true
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		extensionKey := "kubernetes.hostNetwork"
		if app.Platform == PlatformDockerCompose || app.Platform == PlatformDockerSwarm {
			extensionKey = composeApplicationExtensionKey(extensionKey)
		}
		service.Extensions[extensionKey] = "true"
	}
	rehydratePortableKubernetesWindowsHostProcessResources(app)
}

func rehydratePortableKubernetesSecurityContext(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if service.SELinuxOptions == nil {
			if value, ok := firstExtensionValue(service.Extensions, "kubernetes.seLinuxOptions", "x-kubernetes-seLinuxOptions"); ok {
				if options, ok := asMap(value); ok {
					service.SELinuxOptions = parseKubernetesSELinuxOptions(options)
				}
			}
		}
		if service.SELinuxOptions != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.seLinuxOptions"]; !exists {
				service.Extensions["kubernetes.seLinuxOptions"] = serializeKubernetesSELinuxOptions(service.SELinuxOptions)
			}
			if _, exists := service.Extensions["x-kubernetes-seLinuxOptions"]; !exists {
				service.Extensions["x-kubernetes-seLinuxOptions"] = serializeKubernetesSELinuxOptions(service.SELinuxOptions)
			}
		}
		if service.WindowsOptions == nil {
			if value, ok := firstExtensionValue(service.Extensions, "kubernetes.windowsOptions", "x-kubernetes-windowsOptions"); ok {
				if options, ok := asMap(value); ok {
					service.WindowsOptions = parseKubernetesWindowsSecurityContextOptions(options)
				}
			}
		}
		if service.WindowsOptions != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.windowsOptions"]; !exists {
				service.Extensions["kubernetes.windowsOptions"] = serializeKubernetesWindowsSecurityContextOptions(service.WindowsOptions)
			}
			if _, exists := service.Extensions["x-kubernetes-windowsOptions"]; !exists {
				service.Extensions["x-kubernetes-windowsOptions"] = serializeKubernetesWindowsSecurityContextOptions(service.WindowsOptions)
			}
		}
	}
}

func rehydratePortableKubernetesServiceIdentity(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if len(service.ImagePullSecrets) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.imagePullSecrets"); ok {
				if names, err := toStringSlice(value); err == nil && len(names) > 0 {
					service.ImagePullSecrets = append([]string{}, names...)
				}
			}
		}
		if len(service.ImagePullSecrets) > 0 {
			secrets := append([]string{}, service.ImagePullSecrets...)
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.imagePullSecrets"]; !exists {
				service.Extensions["kubernetes.imagePullSecrets"] = deepCopyValue(secrets)
			}
			if _, exists := service.Extensions["x-kubernetes-imagePullSecrets"]; !exists {
				service.Extensions["x-kubernetes-imagePullSecrets"] = deepCopyValue(secrets)
			}
		}
		if service.ImagePullPolicy == "" {
			service.ImagePullPolicy = extensionStringValue(service, "kubernetes.imagePullPolicy")
		}
		if service.ImagePullPolicy != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.imagePullPolicy"]; !exists {
				service.Extensions["kubernetes.imagePullPolicy"] = service.ImagePullPolicy
			}
			if _, exists := service.Extensions["x-kubernetes-imagePullPolicy"]; !exists {
				service.Extensions["x-kubernetes-imagePullPolicy"] = service.ImagePullPolicy
			}
		}
		if service.TerminationMessagePath == "" {
			service.TerminationMessagePath = extensionStringValue(service, "kubernetes.terminationMessagePath")
		}
		if service.TerminationMessagePath != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.terminationMessagePath"]; !exists {
				service.Extensions["kubernetes.terminationMessagePath"] = service.TerminationMessagePath
			}
			if _, exists := service.Extensions["x-kubernetes-terminationMessagePath"]; !exists {
				service.Extensions["x-kubernetes-terminationMessagePath"] = service.TerminationMessagePath
			}
		}
		if service.TerminationMessagePolicy == "" {
			service.TerminationMessagePolicy = extensionStringValue(service, "kubernetes.terminationMessagePolicy")
		}
		if service.TerminationMessagePolicy != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.terminationMessagePolicy"]; !exists {
				service.Extensions["kubernetes.terminationMessagePolicy"] = service.TerminationMessagePolicy
			}
			if _, exists := service.Extensions["x-kubernetes-terminationMessagePolicy"]; !exists {
				service.Extensions["x-kubernetes-terminationMessagePolicy"] = service.TerminationMessagePolicy
			}
		}
		if service.DNSPolicy == "" {
			service.DNSPolicy = extensionStringValue(service, "kubernetes.dnsPolicy")
		}
		if service.DNSPolicy != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.dnsPolicy"]; !exists {
				service.Extensions["kubernetes.dnsPolicy"] = service.DNSPolicy
			}
			if _, exists := service.Extensions["x-kubernetes-dnsPolicy"]; !exists {
				service.Extensions["x-kubernetes-dnsPolicy"] = service.DNSPolicy
			}
		}
		if service.SchedulerName == "" {
			service.SchedulerName = extensionStringValue(service, "kubernetes.schedulerName")
		}
		if service.SchedulerName != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.schedulerName"]; !exists {
				service.Extensions["kubernetes.schedulerName"] = service.SchedulerName
			}
			if _, exists := service.Extensions["x-kubernetes-schedulerName"]; !exists {
				service.Extensions["x-kubernetes-schedulerName"] = service.SchedulerName
			}
		}
		if service.Runtime == "" {
			service.Runtime = extensionStringValue(service, "kubernetes.runtime")
		}
		if service.Runtime != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.runtime"]; !exists {
				service.Extensions["kubernetes.runtime"] = service.Runtime
			}
			if _, exists := service.Extensions["x-kubernetes-runtime"]; !exists {
				service.Extensions["x-kubernetes-runtime"] = service.Runtime
			}
		}
		if !service.HostNetwork && !service.HostNetworkSet {
			if value, ok := extensionBoolValue(service, "kubernetes.hostNetwork"); ok {
				service.HostNetwork = value
				service.HostNetworkSet = true
			}
		}
		if service.HostNetworkSet || service.HostNetwork {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.hostNetwork"]; !exists {
				service.Extensions["kubernetes.hostNetwork"] = fmt.Sprintf("%t", service.HostNetwork)
			}
			if _, exists := service.Extensions["x-kubernetes-hostNetwork"]; !exists {
				service.Extensions["x-kubernetes-hostNetwork"] = fmt.Sprintf("%t", service.HostNetwork)
			}
		}
		if service.HostPID == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.hostPID"); ok {
				service.HostPID = boolPtr(value)
			}
		}
		if service.HostPID != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.hostPID"]; !exists {
				service.Extensions["kubernetes.hostPID"] = fmt.Sprintf("%t", *service.HostPID)
			}
			if _, exists := service.Extensions["x-kubernetes-hostPID"]; !exists {
				service.Extensions["x-kubernetes-hostPID"] = fmt.Sprintf("%t", *service.HostPID)
			}
		}
		if service.HostIPC == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.hostIPC"); ok {
				service.HostIPC = boolPtr(value)
			}
		}
		if service.HostIPC != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.hostIPC"]; !exists {
				service.Extensions["kubernetes.hostIPC"] = fmt.Sprintf("%t", *service.HostIPC)
			}
			if _, exists := service.Extensions["x-kubernetes-hostIPC"]; !exists {
				service.Extensions["x-kubernetes-hostIPC"] = fmt.Sprintf("%t", *service.HostIPC)
			}
		}
		if service.PIDMode == "" {
			if value := extensionStringValue(service, "kubernetes.pidMode"); value != "" {
				service.PIDMode = value
				if strings.EqualFold(value, "host") && service.HostPID == nil {
					service.HostPID = boolPtr(true)
				}
			}
		}
		if service.PIDMode != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.pidMode"]; !exists {
				service.Extensions["kubernetes.pidMode"] = service.PIDMode
			}
			if _, exists := service.Extensions["x-kubernetes-pidMode"]; !exists {
				service.Extensions["x-kubernetes-pidMode"] = service.PIDMode
			}
		}
		if service.IPCMode == "" {
			if value := extensionStringValue(service, "kubernetes.ipcMode"); value != "" {
				service.IPCMode = value
				if strings.EqualFold(value, "host") && service.HostIPC == nil {
					service.HostIPC = boolPtr(true)
				}
			}
		}
		if service.IPCMode != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.ipcMode"]; !exists {
				service.Extensions["kubernetes.ipcMode"] = service.IPCMode
			}
			if _, exists := service.Extensions["x-kubernetes-ipcMode"]; !exists {
				service.Extensions["x-kubernetes-ipcMode"] = service.IPCMode
			}
		}
		if service.PriorityClassName == "" {
			service.PriorityClassName = extensionStringValue(service, "kubernetes.priorityClassName")
		}
		if service.PriorityClassName != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.priorityClassName"]; !exists {
				service.Extensions["kubernetes.priorityClassName"] = service.PriorityClassName
			}
			if _, exists := service.Extensions["x-kubernetes-priorityClassName"]; !exists {
				service.Extensions["x-kubernetes-priorityClassName"] = service.PriorityClassName
			}
		}
		if service.RuntimeClassName == "" {
			service.RuntimeClassName = extensionStringValue(service, "kubernetes.runtimeClassName")
		}
		if service.RuntimeClassName != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.runtimeClassName"]; !exists {
				service.Extensions["kubernetes.runtimeClassName"] = service.RuntimeClassName
			}
			if _, exists := service.Extensions["x-kubernetes-runtimeClassName"]; !exists {
				service.Extensions["x-kubernetes-runtimeClassName"] = service.RuntimeClassName
			}
		}
		if service.NodeName == "" {
			service.NodeName = extensionStringValue(service, "kubernetes.nodeName")
		}
		if service.NodeName != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.nodeName"]; !exists {
				service.Extensions["kubernetes.nodeName"] = service.NodeName
			}
			if _, exists := service.Extensions["x-kubernetes-nodeName"]; !exists {
				service.Extensions["x-kubernetes-nodeName"] = service.NodeName
			}
		}
		if service.Subdomain == "" {
			service.Subdomain = extensionStringValue(service, "kubernetes.subdomain")
		}
		if service.Subdomain != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.subdomain"]; !exists {
				service.Extensions["kubernetes.subdomain"] = service.Subdomain
			}
			if _, exists := service.Extensions["x-kubernetes-subdomain"]; !exists {
				service.Extensions["x-kubernetes-subdomain"] = service.Subdomain
			}
		}
		if service.OSName == "" {
			service.OSName = extensionStringValue(service, "kubernetes.os")
		}
		if service.OSName != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.os"]; !exists {
				service.Extensions["kubernetes.os"] = service.OSName
			}
			if _, exists := service.Extensions["x-kubernetes-os"]; !exists {
				service.Extensions["x-kubernetes-os"] = service.OSName
			}
		}
		if service.ServiceAccountName == "" {
			service.ServiceAccountName = extensionStringValue(service, "kubernetes.serviceAccountName")
		}
		if service.ServiceAccountName != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.serviceAccountName"]; !exists {
				service.Extensions["kubernetes.serviceAccountName"] = service.ServiceAccountName
			}
			if _, exists := service.Extensions["x-kubernetes-serviceAccountName"]; !exists {
				service.Extensions["x-kubernetes-serviceAccountName"] = service.ServiceAccountName
			}
		}
		if service.FSGroupChangePolicy == "" {
			service.FSGroupChangePolicy = extensionStringValue(service, "kubernetes.fsGroupChangePolicy")
		}
		if service.FSGroupChangePolicy != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.fsGroupChangePolicy"]; !exists {
				service.Extensions["kubernetes.fsGroupChangePolicy"] = service.FSGroupChangePolicy
			}
			if _, exists := service.Extensions["x-kubernetes-fsGroupChangePolicy"]; !exists {
				service.Extensions["x-kubernetes-fsGroupChangePolicy"] = service.FSGroupChangePolicy
			}
		}
		if service.SupplementalGroupsPolicy == "" {
			service.SupplementalGroupsPolicy = extensionStringValue(service, "kubernetes.supplementalGroupsPolicy")
		}
		if service.SupplementalGroupsPolicy != "" {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.supplementalGroupsPolicy"]; !exists {
				service.Extensions["kubernetes.supplementalGroupsPolicy"] = service.SupplementalGroupsPolicy
			}
			if _, exists := service.Extensions["x-kubernetes-supplementalGroupsPolicy"]; !exists {
				service.Extensions["x-kubernetes-supplementalGroupsPolicy"] = service.SupplementalGroupsPolicy
			}
		}
		if service.HostUsers == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.hostUsers"); ok {
				service.HostUsers = boolPtr(value)
			}
		}
		if service.HostUsers != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.hostUsers"]; !exists {
				service.Extensions["kubernetes.hostUsers"] = fmt.Sprintf("%t", *service.HostUsers)
			}
			if _, exists := service.Extensions["x-kubernetes-hostUsers"]; !exists {
				service.Extensions["x-kubernetes-hostUsers"] = fmt.Sprintf("%t", *service.HostUsers)
			}
		}
		if service.SetHostnameAsFQDN == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.setHostnameAsFQDN"); ok {
				service.SetHostnameAsFQDN = boolPtr(value)
			}
		}
		if service.SetHostnameAsFQDN != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.setHostnameAsFQDN"]; !exists {
				service.Extensions["kubernetes.setHostnameAsFQDN"] = fmt.Sprintf("%t", *service.SetHostnameAsFQDN)
			}
			if _, exists := service.Extensions["x-kubernetes-setHostnameAsFQDN"]; !exists {
				service.Extensions["x-kubernetes-setHostnameAsFQDN"] = fmt.Sprintf("%t", *service.SetHostnameAsFQDN)
			}
		}
		if service.ShareProcessNamespace == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.shareProcessNamespace"); ok {
				service.ShareProcessNamespace = boolPtr(value)
			}
		}
		if service.ShareProcessNamespace != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.shareProcessNamespace"]; !exists {
				service.Extensions["kubernetes.shareProcessNamespace"] = fmt.Sprintf("%t", *service.ShareProcessNamespace)
			}
			if _, exists := service.Extensions["x-kubernetes-shareProcessNamespace"]; !exists {
				service.Extensions["x-kubernetes-shareProcessNamespace"] = fmt.Sprintf("%t", *service.ShareProcessNamespace)
			}
		}
		if service.EnableServiceLinks == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.enableServiceLinks"); ok {
				service.EnableServiceLinks = boolPtr(value)
			}
		}
		if service.EnableServiceLinks != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.enableServiceLinks"]; !exists {
				service.Extensions["kubernetes.enableServiceLinks"] = fmt.Sprintf("%t", *service.EnableServiceLinks)
			}
			if _, exists := service.Extensions["x-kubernetes-enableServiceLinks"]; !exists {
				service.Extensions["x-kubernetes-enableServiceLinks"] = fmt.Sprintf("%t", *service.EnableServiceLinks)
			}
		}
		if service.AutomountServiceAccountToken == nil {
			if value, ok := extensionBoolValue(service, "kubernetes.automountServiceAccountToken"); ok {
				service.AutomountServiceAccountToken = boolPtr(value)
			}
		}
		if service.AutomountServiceAccountToken != nil {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions["kubernetes.automountServiceAccountToken"]; !exists {
				service.Extensions["kubernetes.automountServiceAccountToken"] = fmt.Sprintf("%t", *service.AutomountServiceAccountToken)
			}
			if _, exists := service.Extensions["x-kubernetes-automountServiceAccountToken"]; !exists {
				service.Extensions["x-kubernetes-automountServiceAccountToken"] = fmt.Sprintf("%t", *service.AutomountServiceAccountToken)
			}
		}
		if service.SeccompProfile == nil {
			if value, ok := extensionValueForKey(service, "kubernetes.seccompProfile"); ok {
				if profile, ok := asMap(value); ok && len(profile) > 0 {
					service.SeccompProfile = parseKubernetesSeccompProfile(profile)
				}
			}
		}
		if service.Affinity == nil {
			if value, ok := extensionValueForKey(service, "kubernetes.affinity"); ok {
				if affinity, ok := asMap(value); ok && len(affinity) > 0 {
					service.Affinity = copyStringInterfaceMap(affinity)
				}
			}
		}
		if len(service.ReadinessGates) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.readinessGates"); ok {
				if gates, ok := kubernetesReadinessGatesFromExtension(value); ok && len(gates) > 0 {
					service.ReadinessGates = gates
				}
			}
		}
		if len(service.InitContainers) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.initContainers"); ok {
				if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
					service.InitContainers = containers
				}
			}
		}
		if len(service.SchedulingGates) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.schedulingGates"); ok {
				if gates, ok := kubernetesMapSliceFromExtension(value); ok && len(gates) > 0 {
					service.SchedulingGates = gates
				}
			}
		}
		if len(service.TopologySpreadConstraints) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.topologySpreadConstraints"); ok {
				if constraints, ok := kubernetesTopologySpreadConstraintsFromExtension(value); ok && len(constraints) > 0 {
					service.TopologySpreadConstraints = constraints
				}
			}
		}
		if len(service.GroupAdd) == 0 {
			if value, ok := extensionValueForKey(service, "kubernetes.groupAdd"); ok {
				if items, err := toStringSlice(value); err == nil && len(items) > 0 {
					service.GroupAdd = items
				}
			}
		}
	}
}

func rehydratePortableKubernetesHostNetworkDNSPolicy(app *Application) {
	if app == nil {
		return
	}
	const dnsPolicy = "ClusterFirstWithHostNet"
	serviceNames := map[string]struct{}{}
	for name, service := range app.Services {
		if service == nil || !service.HostNetwork || strings.TrimSpace(service.DNSPolicy) != "" {
			continue
		}
		service.DNSPolicy = dnsPolicy
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["kubernetes.dnsPolicy"] = dnsPolicy
		serviceNames[name] = struct{}{}
	}
	if len(serviceNames) == 0 {
		return
	}
	patchResource := func(resource map[string]interface{}, serviceName string) {
		if len(resource) == 0 || serviceName == "" {
			return
		}
		metadata, _ := asMap(resource["metadata"])
		if toString(metadata["name"]) != serviceName {
			return
		}
		podSpec := kubernetesWorkloadPodSpec(resource)
		if podSpec == nil {
			return
		}
		if value := strings.TrimSpace(toString(podSpec["dnsPolicy"])); value == "" {
			podSpec["dnsPolicy"] = dnsPolicy
		}
	}
	for _, resource := range kubernetesRawExtensionResources(app.Extensions["kubernetes.raw"]) {
		for serviceName := range serviceNames {
			patchResource(resource, serviceName)
		}
	}
	for _, resource := range kubernetesRawExtensionResources(app.Extensions[kubernetesWorkloadsExtensionKey]) {
		for serviceName := range serviceNames {
			patchResource(resource, serviceName)
		}
	}
	for _, resource := range kubernetesRawExtensionResources(app.Extensions[composeKubernetesWorkloadsExtensionKey]) {
		for serviceName := range serviceNames {
			patchResource(resource, serviceName)
		}
	}
	for _, manifest := range app.KubernetesOpaqueManifests {
		if manifest == nil {
			continue
		}
		for serviceName := range serviceNames {
			patchResource(manifest.Raw, serviceName)
		}
	}
}

func kubernetesTolerationsFromExtension(value interface{}) ([]Toleration, bool) {
	items, ok := extensionSliceValues(value)
	if !ok || len(items) == 0 {
		return nil, false
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return nil, false
	}
	var tolerations []Toleration
	if err := json.Unmarshal(raw, &tolerations); err != nil {
		return nil, false
	}
	if len(tolerations) == 0 {
		return nil, false
	}
	return cloneTolerations(tolerations), true
}

func rehydratePortableKubernetesWindowsHostProcessResources(app *Application) {
	if app == nil || len(app.Services) == 0 {
		return
	}
	serviceNames := map[string]struct{}{}
	for name, service := range app.Services {
		if service == nil || service.WindowsOptions == nil || service.WindowsOptions.HostProcess == nil || !*service.WindowsOptions.HostProcess {
			continue
		}
		serviceNames[name] = struct{}{}
	}
	if len(serviceNames) == 0 {
		return
	}
	patchResource := func(resource map[string]interface{}, serviceName string) {
		if len(resource) == 0 || serviceName == "" {
			return
		}
		metadata, _ := asMap(resource["metadata"])
		if toString(metadata["name"]) != serviceName {
			return
		}
		podSpec := kubernetesWorkloadPodSpec(resource)
		if podSpec == nil {
			return
		}
		podSpec["hostNetwork"] = true
	}
	for _, resource := range kubernetesRawExtensionResources(app.Extensions["kubernetes.raw"]) {
		for serviceName := range serviceNames {
			patchResource(resource, serviceName)
		}
	}
	for _, resource := range kubernetesRawExtensionResources(app.Extensions[kubernetesWorkloadsExtensionKey]) {
		for serviceName := range serviceNames {
			patchResource(resource, serviceName)
		}
	}
	for _, resource := range kubernetesRawExtensionResources(app.Extensions[composeKubernetesWorkloadsExtensionKey]) {
		for serviceName := range serviceNames {
			patchResource(resource, serviceName)
		}
	}
	for _, manifest := range app.KubernetesOpaqueManifests {
		if manifest == nil {
			continue
		}
		for serviceName := range serviceNames {
			patchResource(manifest.Raw, serviceName)
		}
	}
}

func kubernetesWorkloadPodSpec(resource map[string]interface{}) map[string]interface{} {
	if len(resource) == 0 {
		return nil
	}
	spec, _ := asMap(resource["spec"])
	if spec == nil {
		return nil
	}
	switch toString(resource["kind"]) {
	case "Pod":
		return spec
	case "Deployment", "StatefulSet", "DaemonSet":
		template, _ := asMap(spec["template"])
		if template == nil {
			return nil
		}
		podSpec, _ := asMap(template["spec"])
		return podSpec
	case "Job":
		template, _ := asMap(spec["template"])
		if template == nil {
			return nil
		}
		podSpec, _ := asMap(template["spec"])
		return podSpec
	case "CronJob":
		jobTemplate, _ := asMap(spec["jobTemplate"])
		if jobTemplate == nil {
			return nil
		}
		jobSpec, _ := asMap(jobTemplate["spec"])
		if jobSpec == nil {
			return nil
		}
		template, _ := asMap(jobSpec["template"])
		if template == nil {
			return nil
		}
		podSpec, _ := asMap(template["spec"])
		return podSpec
	default:
		return nil
	}
}

func rehydratePortableNomadScheduler(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		update, _ := firstExtensionValue(service.Extensions, nomadUpdateExtensionKey, "x-nomad-update")
		migrate, _ := firstExtensionValue(service.Extensions, nomadMigrateExtensionKey, "x-nomad-migrate")
		reschedule, _ := firstExtensionValue(service.Extensions, nomadRescheduleExtensionKey, "x-nomad-reschedule")
		if update == nil && migrate == nil && reschedule == nil {
			continue
		}
		applyNomadSchedulerBlockExtensions(service, schedulerBlockMap(update), schedulerBlockMap(migrate), schedulerBlockMap(reschedule))
	}
}

func nomadConnectSpecFromAny(value interface{}) *NomadConnectSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case *NomadConnectSpec:
		return cloneNomadConnectSpec(typed)
	case NomadConnectSpec:
		return cloneNomadConnectSpec(&typed)
	case map[string]interface{}:
		connect := &NomadConnectSpec{Extensions: map[string]interface{}{}}
		for key, raw := range typed {
			switch key {
			case "native":
				connect.Native = toBool(raw)
				connect.NativeSet = true
			case "sidecar_service":
				connect.SidecarService = nomadConnectSidecarServiceFromAny(raw)
			case "gateway":
				connect.Gateway = nomadConnectGatewayFromAny(raw)
			case "extensions":
				if ext, ok := asMap(raw); ok {
					for extKey, extValue := range ext {
						connect.Extensions[extKey] = deepCopyValue(extValue)
					}
				}
			default:
				connect.Extensions[key] = deepCopyValue(raw)
			}
		}
		if len(connect.Extensions) == 0 {
			connect.Extensions = nil
		}
		if isEmptyNomadConnectSpec(connect) {
			return nil
		}
		normalizeNomadConnectSpec(connect)
		return connect
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	connect, err := parseNomadConnectSpecMapBytes(raw)
	if err != nil || connect == nil {
		return nil
	}
	return connect
}

func nomadConnectSpecToMap(connect *NomadConnectSpec) map[string]interface{} {
	if connect == nil || isEmptyNomadConnectSpec(connect) {
		return nil
	}
	return serializeNomadConnectSpec(connect)
}

func parseNomadConnectSpecMapBytes(data []byte) (*NomadConnectSpec, error) {
	var mapped map[string]interface{}
	if err := json.Unmarshal(data, &mapped); err != nil {
		return nil, err
	}
	return parseNomadConnectSpecMap(mapped)
}

func parseNomadConnectSpecMap(mapped map[string]interface{}) (*NomadConnectSpec, error) {
	if len(mapped) == 0 {
		return nil, nil
	}
	connect := &NomadConnectSpec{Extensions: map[string]interface{}{}}
	for key, raw := range mapped {
		switch key {
		case "native":
			connect.Native = toBool(raw)
			connect.NativeSet = true
		case "sidecar_service":
			connect.SidecarService = nomadConnectSidecarServiceFromAny(raw)
		case "check":
			if check, ok := asMap(raw); ok && len(check) > 0 {
				connect.SidecarService = nomadConnectSidecarServiceFromAny(map[string]interface{}{"check": check})
			}
		case "gateway":
			connect.Gateway = nomadConnectGatewayFromAny(raw)
		case "extensions":
			if ext, ok := asMap(raw); ok {
				for extKey, extValue := range ext {
					connect.Extensions[extKey] = deepCopyValue(extValue)
				}
			}
		default:
			connect.Extensions[key] = deepCopyValue(raw)
		}
	}
	if len(connect.Extensions) == 0 {
		connect.Extensions = nil
	}
	if isEmptyNomadConnectSpec(connect) {
		return nil, nil
	}
	normalizeNomadConnectSpec(connect)
	return connect, nil
}

func normalizeNomadConnectSpec(connect *NomadConnectSpec) {
	if connect == nil {
		return
	}
	if connect.SidecarService != nil {
		if len(connect.SidecarService.Tags) > 0 {
			connect.SidecarService.Tags = append([]string{}, connect.SidecarService.Tags...)
		}
		if len(connect.SidecarService.Check) > 0 {
			connect.SidecarService.Check = cloneMap(connect.SidecarService.Check)
		}
		if connect.SidecarService.Proxy != nil {
			if len(connect.SidecarService.Proxy.Upstreams) > 0 {
				connect.SidecarService.Proxy.Upstreams = append([]NomadConnectUpstream{}, connect.SidecarService.Proxy.Upstreams...)
			}
			if len(connect.SidecarService.Proxy.Config) > 0 {
				connect.SidecarService.Proxy.Config = cloneMap(connect.SidecarService.Proxy.Config)
			}
		}
	}
	if len(connect.Gateway) > 0 {
		connect.Gateway = cloneMap(connect.Gateway)
	}
}

func serializeNomadConnectSpec(connect *NomadConnectSpec) map[string]interface{} {
	if connect == nil || isEmptyNomadConnectSpec(connect) {
		return nil
	}
	cloned := cloneNomadConnectSpec(connect)
	if cloned == nil {
		return nil
	}
	result := map[string]interface{}{}
	if cloned.NativeSet || cloned.Native {
		result["native"] = cloned.Native
	}
	if cloned.SidecarService != nil {
		sidecar := map[string]interface{}{}
		if len(cloned.SidecarService.Tags) > 0 {
			sidecar["tags"] = append([]string{}, cloned.SidecarService.Tags...)
		}
		if cloned.SidecarService.Proxy != nil {
			proxy := map[string]interface{}{}
			if len(cloned.SidecarService.Proxy.Upstreams) > 0 {
				upstreams := make([]interface{}, 0, len(cloned.SidecarService.Proxy.Upstreams))
				for _, upstream := range cloned.SidecarService.Proxy.Upstreams {
					upstreamMap := map[string]interface{}{}
					if upstream.DestinationName != "" {
						upstreamMap["destination_name"] = upstream.DestinationName
					}
					if upstream.LocalBindPort != 0 {
						upstreamMap["local_bind_port"] = upstream.LocalBindPort
					}
					if upstream.LocalBindAddress != "" {
						upstreamMap["local_bind_address"] = upstream.LocalBindAddress
					}
					for key, value := range upstream.Extensions {
						if upstreamMap["extensions"] == nil {
							upstreamMap["extensions"] = map[string]interface{}{}
						}
						upstreamMap["extensions"].(map[string]interface{})[key] = deepCopyValue(value)
					}
					upstreams = append(upstreams, upstreamMap)
				}
				proxy["upstreams"] = upstreams
			}
			if len(cloned.SidecarService.Proxy.Config) > 0 {
				proxy["config"] = cloneMap(cloned.SidecarService.Proxy.Config)
			}
			for key, value := range cloned.SidecarService.Proxy.Extensions {
				if proxy["extensions"] == nil {
					proxy["extensions"] = map[string]interface{}{}
				}
				proxy["extensions"].(map[string]interface{})[key] = deepCopyValue(value)
			}
			if len(proxy) > 0 {
				sidecar["proxy"] = proxy
			}
		}
		if len(cloned.SidecarService.Check) > 0 {
			sidecar["check"] = cloneMap(cloned.SidecarService.Check)
		}
		for key, value := range cloned.SidecarService.Extensions {
			if sidecar["extensions"] == nil {
				sidecar["extensions"] = map[string]interface{}{}
			}
			sidecar["extensions"].(map[string]interface{})[key] = deepCopyValue(value)
		}
		if len(sidecar) > 0 {
			result["sidecar_service"] = sidecar
		}
	}
	if len(cloned.Gateway) > 0 {
		result["gateway"] = cloneMap(cloned.Gateway)
	}
	for key, value := range cloned.Extensions {
		if result["extensions"] == nil {
			result["extensions"] = map[string]interface{}{}
		}
		result["extensions"].(map[string]interface{})[key] = deepCopyValue(value)
	}
	return result
}

func nomadConnectSidecarServiceFromAny(value interface{}) *NomadConnectSidecarService {
	mapped, ok := asMap(value)
	if !ok || len(mapped) == 0 {
		return nil
	}
	sidecar := &NomadConnectSidecarService{Extensions: map[string]interface{}{}}
	if tags, err := toStringSlice(mapped["tags"]); err == nil {
		sidecar.Tags = tags
	}
	if proxyValue, ok := asMap(mapped["proxy"]); ok && len(proxyValue) > 0 {
		proxy := &NomadConnectProxy{Extensions: map[string]interface{}{}}
		if upstreams, ok := proxyValue["upstreams"].([]interface{}); ok {
			for _, item := range upstreams {
				if upstreamMap, ok := asMap(item); ok {
					upstream := NomadConnectUpstream{
						DestinationName:  toString(upstreamMap["destination_name"]),
						LocalBindPort:    toInt(upstreamMap["local_bind_port"]),
						LocalBindAddress: toString(upstreamMap["local_bind_address"]),
						Extensions:       map[string]interface{}{},
					}
					for key, raw := range upstreamMap {
						switch key {
						case "destination_name", "local_bind_port", "local_bind_address":
							continue
						case "extensions":
							if ext, ok := asMap(raw); ok {
								for extKey, extValue := range ext {
									upstream.Extensions[extKey] = deepCopyValue(extValue)
								}
							}
						default:
							upstream.Extensions[key] = deepCopyValue(raw)
						}
					}
					if len(upstream.Extensions) == 0 {
						upstream.Extensions = nil
					}
					proxy.Upstreams = append(proxy.Upstreams, upstream)
				}
			}
		}
		if configValue, ok := asMap(proxyValue["config"]); ok && len(configValue) > 0 {
			proxy.Config = cloneMap(configValue)
		}
		for key, raw := range proxyValue {
			switch key {
			case "upstreams", "config":
				continue
			case "extensions":
				if ext, ok := asMap(raw); ok {
					for extKey, extValue := range ext {
						proxy.Extensions[extKey] = deepCopyValue(extValue)
					}
				}
			default:
				proxy.Extensions[key] = deepCopyValue(raw)
			}
		}
		if len(proxy.Extensions) == 0 {
			proxy.Extensions = nil
		}
		sidecar.Proxy = proxy
	}
	if checkValue, ok := asMap(mapped["check"]); ok && len(checkValue) > 0 {
		sidecar.Check = cloneMap(checkValue)
	}
	for key, raw := range mapped {
		switch key {
		case "tags", "proxy", "check":
			continue
		case "extensions":
			if ext, ok := asMap(raw); ok {
				for extKey, extValue := range ext {
					sidecar.Extensions[extKey] = deepCopyValue(extValue)
				}
			}
		default:
			sidecar.Extensions[key] = deepCopyValue(raw)
		}
	}
	if len(sidecar.Extensions) == 0 {
		sidecar.Extensions = nil
	}
	if sidecar.Proxy == nil && len(sidecar.Tags) == 0 && len(sidecar.Check) == 0 && len(sidecar.Extensions) == 0 {
		return nil
	}
	return sidecar
}

func nomadConnectGatewayFromAny(value interface{}) map[string]interface{} {
	mapped, ok := asMap(value)
	if !ok || len(mapped) == 0 {
		return nil
	}
	cloned := cloneMap(mapped)
	if proxy, ok := asMap(cloned["proxy"]); ok && len(proxy) > 0 {
		cloned["proxy"] = cloneMap(proxy)
	} else if proxyList, ok := cloned["proxy"].([]interface{}); ok && len(proxyList) == 1 {
		if proxy, ok := asMap(proxyList[0]); ok && len(proxy) > 0 {
			cloned["proxy"] = cloneMap(proxy)
		}
	}
	return cloned
}

func nomadSpreadSpecsFromAny(value interface{}) []NomadSpreadSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case []NomadSpreadSpec:
		return cloneNomadSpreadSpecs(typed)
	case []map[string]interface{}:
		result := make([]NomadSpreadSpec, 0, len(typed))
		for _, item := range typed {
			if spec := nomadSpreadSpecFromMap(item); spec != nil {
				result = append(result, *spec)
			}
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var spreads []NomadSpreadSpec
	if err := json.Unmarshal(raw, &spreads); err != nil || len(spreads) == 0 {
		return nil
	}
	return cloneNomadSpreadSpecs(spreads)
}

func nomadSpreadSpecsToMap(spreads []NomadSpreadSpec) []map[string]interface{} {
	if len(spreads) == 0 {
		return nil
	}
	cloned := cloneNomadSpreadSpecs(spreads)
	if len(cloned) == 0 {
		return nil
	}
	raw, err := json.Marshal(cloned)
	if err != nil {
		return nil
	}
	var mapped []map[string]interface{}
	if err := json.Unmarshal(raw, &mapped); err != nil || len(mapped) == 0 {
		return nil
	}
	return mapped
}

func nomadSchedulerExtensionMap(service *Service, primaryKey, fallbackKey string) map[string]interface{} {
	if service == nil {
		return nil
	}
	if value, ok := firstExtensionValue(service.Extensions, primaryKey, fallbackKey); ok {
		if mapped, ok := asMap(value); ok && len(mapped) > 0 {
			return cloneMap(mapped)
		}
	}
	if service.Deploy != nil {
		switch primaryKey {
		case nomadUpdateExtensionKey:
			if !isEmptyUpdatePolicy(service.Deploy.UpdateConfig) {
				return serializeNomadUpdatePolicy(service.Deploy.UpdateConfig)
			}
		case nomadMigrateExtensionKey:
			if !isEmptyMigratePolicy(service.Deploy.MigrateConfig) {
				return serializeMigratePolicy(service.Deploy.MigrateConfig)
			}
		case nomadRescheduleExtensionKey:
			if !isEmptyReschedulePolicy(service.Deploy.RescheduleConfig) {
				return serializeReschedulePolicy(service.Deploy.RescheduleConfig)
			}
		}
	}
	return nil
}

func schedulerBlockMap(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}
	if mapped, ok := asMap(value); ok && len(mapped) > 0 {
		return cloneMap(mapped)
	}
	return nil
}

func nomadSpreadSpecFromMap(mapped map[string]interface{}) *NomadSpreadSpec {
	if len(mapped) == 0 {
		return nil
	}
	spec := &NomadSpreadSpec{
		Attribute:  toString(mapped["attribute"]),
		Portable:   toString(mapped["portable"]),
		Weight:     toInt(mapped["weight"]),
		Extensions: map[string]interface{}{},
	}
	if targets, ok := mapped["targets"].([]interface{}); ok {
		for _, item := range targets {
			targetMap, ok := asMap(item)
			if !ok {
				continue
			}
			target := NomadSpreadTarget{
				Value:      toString(targetMap["value"]),
				Percent:    toInt(targetMap["percent"]),
				Extensions: map[string]interface{}{},
			}
			for key, value := range targetMap {
				switch key {
				case "value", "percent", "extensions":
					if key == "extensions" {
						if ext, ok := asMap(value); ok {
							for extKey, extValue := range ext {
								target.Extensions[extKey] = deepCopyValue(extValue)
							}
						}
					}
				default:
					target.Extensions[key] = deepCopyValue(value)
				}
			}
			if len(target.Extensions) == 0 {
				target.Extensions = nil
			}
			spec.Targets = append(spec.Targets, target)
		}
	}
	for key, value := range mapped {
		switch key {
		case "attribute", "portable", "targets", "extensions":
			if key == "extensions" {
				if ext, ok := asMap(value); ok {
					for extKey, extValue := range ext {
						spec.Extensions[extKey] = deepCopyValue(extValue)
					}
				}
			}
		default:
			spec.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(spec.Extensions) == 0 {
		spec.Extensions = nil
	}
	if isEmptyNomadSpreadSpec(spec) {
		return nil
	}
	return spec
}

func kubernetesServiceSpecsFromAny(value interface{}) map[string]*KubernetesServiceSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesServiceSpec:
		return cloneKubernetesServiceSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesServiceSpec, len(typed))
		for _, item := range typed {
			spec := kubernetesServiceSpecFromMap(item)
			if spec == nil || spec.Name == "" {
				continue
			}
			result[spec.Name] = spec
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesServiceSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			spec := kubernetesServiceSpecFromMap(mapped)
			if spec == nil || spec.Name == "" {
				continue
			}
			result[spec.Name] = spec
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesServiceSpecsFromAny(typedList)
}

func kubernetesServiceSpecFromMap(mapped map[string]interface{}) *KubernetesServiceSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	service := &KubernetesServiceSpec{
		Name:        toString(metadata["name"]),
		Namespace:   toString(metadata["namespace"]),
		Type:        toString(specMap["type"]),
		Selector:    toStringMapLoose(specMap["selector"]),
		Annotations: toStringMapLoose(metadata["annotations"]),
		Labels:      toStringMapLoose(metadata["labels"]),
		Extensions:  map[string]interface{}{},
		Raw:         deepCopyValue(mapped).(map[string]interface{}),
	}
	service.ExternalName = toString(specMap["externalName"])
	service.SessionAffinity = toString(specMap["sessionAffinity"])
	service.LoadBalancerIP = toString(specMap["loadBalancerIP"])
	service.LoadBalancerClass = toString(specMap["loadBalancerClass"])
	service.IPFamilyPolicy = toString(specMap["ipFamilyPolicy"])
	service.ExternalTrafficPolicy = toString(specMap["externalTrafficPolicy"])
	service.InternalTrafficPolicy = toString(specMap["internalTrafficPolicy"])
	service.TrafficDistribution = toString(specMap["trafficDistribution"])
	service.HealthCheckNodePort = toInt(specMap["healthCheckNodePort"])
	service.LoadBalancerSourceRanges = appendStringSliceFromAny(specMap["loadBalancerSourceRanges"])
	service.ExternalIPs = appendStringSliceFromAny(specMap["externalIPs"])
	service.IPFamilies = appendStringSliceFromAny(specMap["ipFamilies"])
	if b, ok := boolPtrFromAny(specMap["publishNotReadyAddresses"]); ok {
		service.PublishNotReadyAddresses = b
	}
	if b, ok := boolPtrFromAny(specMap["allocateLoadBalancerNodePorts"]); ok {
		service.AllocateLoadBalancerNodePorts = b
	}
	if ports, ok := specMap["ports"].([]interface{}); ok {
		for _, portValue := range ports {
			if portMap, ok := asMap(portValue); ok {
				service.Ports = append(service.Ports, kubernetesServicePortFromMap(portMap))
			}
		}
	}
	if len(service.Extensions) == 0 {
		service.Extensions = nil
	}
	return service
}

func kubernetesServiceAccountSpecsFromAny(value interface{}) map[string]*KubernetesServiceAccountSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesServiceAccountSpec:
		return cloneKubernetesServiceAccountSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesServiceAccountSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesServiceAccountSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesServiceAccountSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesServiceAccountSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesServiceAccountSpecsFromAny(typedList)
}

func kubernetesServiceAccountSpecFromMap(mapped map[string]interface{}) *KubernetesServiceAccountSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	account := &KubernetesServiceAccountSpec{
		Name:        toString(metadata["name"]),
		Namespace:   toString(metadata["namespace"]),
		Labels:      toStringMapLoose(metadata["labels"]),
		Annotations: toStringMapLoose(metadata["annotations"]),
		Extensions:  map[string]interface{}{},
	}
	account.Secrets = appendNamedSliceFromAny(mapped["secrets"])
	account.ImagePullSecrets = appendNamedSliceFromAny(mapped["imagePullSecrets"])
	if b, ok := boolPtrFromAny(mapped["automountServiceAccountToken"]); ok {
		account.AutomountServiceAccountToken = b
	}
	account.Raw = nil
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		account.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "secrets", "imagePullSecrets", "automountServiceAccountToken":
		default:
			account.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(account.Extensions) == 0 {
		account.Extensions = nil
	}
	return account
}

func kubernetesHPASpecsFromAny(value interface{}) map[string]*KubernetesHorizontalPodAutoscalerSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesHorizontalPodAutoscalerSpec:
		return cloneKubernetesHPASpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesHorizontalPodAutoscalerSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesHPASpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesHorizontalPodAutoscalerSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesHPASpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesHPASpecsFromAny(typedList)
}

func kubernetesHPASpecFromMap(mapped map[string]interface{}) *KubernetesHorizontalPodAutoscalerSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	hpa := &KubernetesHorizontalPodAutoscalerSpec{
		Name:        toString(metadata["name"]),
		Namespace:   toString(metadata["namespace"]),
		MaxReplicas: toInt(specMap["maxReplicas"]),
		ScaleTarget: map[string]string{},
		Extensions:  map[string]interface{}{},
	}
	if min := specMap["minReplicas"]; min != nil {
		value := toInt(min)
		hpa.MinReplicas = &value
	}
	if scaleTarget, ok := asMap(specMap["scaleTargetRef"]); ok {
		hpa.ScaleTarget["apiVersion"] = toString(scaleTarget["apiVersion"])
		hpa.ScaleTarget["kind"] = toString(scaleTarget["kind"])
		hpa.ScaleTarget["name"] = toString(scaleTarget["name"])
		if ns := toString(scaleTarget["namespace"]); ns != "" {
			hpa.ScaleTarget["namespace"] = ns
		}
	}
	if metrics, ok := specMap["metrics"].([]interface{}); ok {
		for _, item := range metrics {
			if mappedMetric, ok := asMap(item); ok {
				hpa.Metrics = append(hpa.Metrics, mappedMetric)
			}
		}
	}
	if behavior, ok := asMap(specMap["behavior"]); ok {
		hpa.Behavior = behavior
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		hpa.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "spec":
		default:
			hpa.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(hpa.Extensions) == 0 {
		hpa.Extensions = nil
	}
	if hpa.Name == "" {
		hpa.Name = toString(metadata["name"])
	}
	return hpa
}

func kubernetesPDBSpecsFromAny(value interface{}) map[string]*KubernetesPodDisruptionBudgetSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesPodDisruptionBudgetSpec:
		return cloneKubernetesPDBSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesPodDisruptionBudgetSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesPDBSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesPodDisruptionBudgetSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesPDBSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesPDBSpecsFromAny(typedList)
}

func kubernetesPDBSpecFromMap(mapped map[string]interface{}) *KubernetesPodDisruptionBudgetSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	pdb := &KubernetesPodDisruptionBudgetSpec{
		Name:                       toString(metadata["name"]),
		Namespace:                  toString(metadata["namespace"]),
		MinAvailable:               toString(specMap["minAvailable"]),
		MaxUnavailable:             toString(specMap["maxUnavailable"]),
		UnhealthyPodEvictionPolicy: toString(specMap["unhealthyPodEvictionPolicy"]),
		Selector:                   toStringMapLoose(specMap["selector"]),
		Annotations:                toStringMapLoose(metadata["annotations"]),
		Labels:                     toStringMapLoose(metadata["labels"]),
		Extensions:                 map[string]interface{}{},
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		pdb.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "spec":
		default:
			pdb.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(pdb.Extensions) == 0 {
		pdb.Extensions = nil
	}
	return pdb
}

func kubernetesResourceQuotaSpecsFromAny(value interface{}) map[string]*KubernetesResourceQuotaSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesResourceQuotaSpec:
		return cloneKubernetesResourceQuotaSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesResourceQuotaSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesResourceQuotaSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesResourceQuotaSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesResourceQuotaSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesResourceQuotaSpecsFromAny(typedList)
}

func kubernetesResourceQuotaSpecFromMap(mapped map[string]interface{}) *KubernetesResourceQuotaSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	quota := &KubernetesResourceQuotaSpec{
		Name:          toString(metadata["name"]),
		Namespace:     toString(metadata["namespace"]),
		Scopes:        appendStringSliceFromAny(specMap["scopes"]),
		ScopeSelector: copyStringInterfaceMap(asMapOrEmpty(specMap["scopeSelector"])),
		Hard:          copyStringInterfaceMap(asMapOrEmpty(specMap["hard"])),
		Annotations:   toStringMapLoose(metadata["annotations"]),
		Labels:        toStringMapLoose(metadata["labels"]),
		Extensions:    map[string]interface{}{},
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		quota.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "spec":
		default:
			quota.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(quota.Extensions) == 0 {
		quota.Extensions = nil
	}
	if quota.Name == "" {
		quota.Name = toString(metadata["name"])
	}
	return quota
}

func kubernetesLimitRangeSpecsFromAny(value interface{}) map[string]*KubernetesLimitRangeSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesLimitRangeSpec:
		return cloneKubernetesLimitRangeSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesLimitRangeSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesLimitRangeSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesLimitRangeSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesLimitRangeSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesLimitRangeSpecsFromAny(typedList)
}

func kubernetesLimitRangeSpecFromMap(mapped map[string]interface{}) *KubernetesLimitRangeSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	limitRange := &KubernetesLimitRangeSpec{
		Name:        toString(metadata["name"]),
		Namespace:   toString(metadata["namespace"]),
		Annotations: toStringMapLoose(metadata["annotations"]),
		Labels:      toStringMapLoose(metadata["labels"]),
		Extensions:  map[string]interface{}{},
	}
	if limits, ok := specMap["limits"].([]interface{}); ok {
		for _, item := range limits {
			if limit, ok := asMap(item); ok {
				limitRange.Limits = append(limitRange.Limits, limit)
			}
		}
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		limitRange.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "spec":
		default:
			limitRange.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(limitRange.Extensions) == 0 {
		limitRange.Extensions = nil
	}
	if limitRange.Name == "" {
		limitRange.Name = toString(metadata["name"])
	}
	return limitRange
}

func kubernetesStorageClassSpecsFromAny(value interface{}) map[string]*KubernetesStorageClassSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesStorageClassSpec:
		return cloneKubernetesStorageClassSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesStorageClassSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesStorageClassSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesStorageClassSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesStorageClassSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesStorageClassSpecsFromAny(typedList)
}

func kubernetesStorageClassSpecFromMap(mapped map[string]interface{}) *KubernetesStorageClassSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	sc := &KubernetesStorageClassSpec{
		Name:                 toString(metadata["name"]),
		Namespace:            toString(metadata["namespace"]),
		Provisioner:          toString(mapped["provisioner"]),
		Parameters:           toStringMapLoose(mapped["parameters"]),
		ReclaimPolicy:        toString(mapped["reclaimPolicy"]),
		AllowVolumeExpansion: boolPtrFromInterface(mapped["allowVolumeExpansion"]),
		VolumeBindingMode:    toString(mapped["volumeBindingMode"]),
		MountOptions:         appendNamedSliceFromAny(mapped["mountOptions"]),
		Annotations:          toStringMapLoose(metadata["annotations"]),
		Labels:               toStringMapLoose(metadata["labels"]),
		Extensions:           map[string]interface{}{},
	}
	if allowedTopologies, ok := mapped["allowedTopologies"].([]interface{}); ok {
		for _, item := range allowedTopologies {
			if topo, ok := asMap(item); ok {
				sc.AllowedTopologies = append(sc.AllowedTopologies, topo)
			}
		}
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		sc.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "provisioner", "parameters", "reclaimPolicy", "allowVolumeExpansion", "volumeBindingMode", "mountOptions", "allowedTopologies":
		default:
			sc.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(sc.Extensions) == 0 {
		sc.Extensions = nil
	}
	if sc.Name == "" {
		sc.Name = toString(metadata["name"])
	}
	return sc
}

func kubernetesIngressClassSpecsFromAny(value interface{}) map[string]*KubernetesIngressClassSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesIngressClassSpec:
		return cloneKubernetesIngressClassSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesIngressClassSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesIngressClassSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesIngressClassSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesIngressClassSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesIngressClassSpecsFromAny(typedList)
}

func kubernetesIngressClassSpecFromMap(mapped map[string]interface{}) *KubernetesIngressClassSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	ic := &KubernetesIngressClassSpec{
		Name:        toString(metadata["name"]),
		Namespace:   toString(metadata["namespace"]),
		Controller:  toString(specMap["controller"]),
		Parameters:  copyStringInterfaceMap(asMapOrEmpty(specMap["parameters"])),
		Annotations: toStringMapLoose(metadata["annotations"]),
		Labels:      toStringMapLoose(metadata["labels"]),
		Extensions:  map[string]interface{}{},
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		ic.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "spec":
		default:
			ic.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(ic.Extensions) == 0 {
		ic.Extensions = nil
	}
	if ic.Name == "" {
		ic.Name = toString(metadata["name"])
	}
	return ic
}

func kubernetesWebhookConfigurationSpecsFromAny(value interface{}) map[string]*KubernetesWebhookConfigurationSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesWebhookConfigurationSpec:
		return cloneKubernetesWebhookConfigurationSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesWebhookConfigurationSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesWebhookConfigurationSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesWebhookConfigurationSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesWebhookConfigurationSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesWebhookConfigurationSpecsFromAny(typedList)
}

func kubernetesWebhookConfigurationSpecFromMap(mapped map[string]interface{}) *KubernetesWebhookConfigurationSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["webhooks"])
	webhooks := []map[string]interface{}{}
	if items, ok := mapped["webhooks"].([]interface{}); ok {
		for _, item := range items {
			if webhook, ok := asMap(item); ok {
				webhooks = append(webhooks, webhook)
			}
		}
	}
	if len(webhooks) == 0 && len(specMap) > 0 {
		webhooks = append(webhooks, specMap)
	}
	config := &KubernetesWebhookConfigurationSpec{
		Name:        toString(metadata["name"]),
		Namespace:   toString(metadata["namespace"]),
		Webhooks:    webhooks,
		Annotations: toStringMapLoose(metadata["annotations"]),
		Labels:      toStringMapLoose(metadata["labels"]),
		Extensions:  map[string]interface{}{},
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		config.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "webhooks":
		default:
			config.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(config.Extensions) == 0 {
		config.Extensions = nil
	}
	if config.Name == "" {
		config.Name = toString(metadata["name"])
	}
	return config
}

func kubernetesCustomResourceDefinitionSpecsFromAny(value interface{}) map[string]*KubernetesCustomResourceDefinitionSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesCustomResourceDefinitionSpec:
		return cloneKubernetesCustomResourceDefinitionSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesCustomResourceDefinitionSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesCustomResourceDefinitionSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesCustomResourceDefinitionSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesCustomResourceDefinitionSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesCustomResourceDefinitionSpecsFromAny(typedList)
}

func kubernetesCustomResourceDefinitionSpecFromMap(mapped map[string]interface{}) *KubernetesCustomResourceDefinitionSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	crd := &KubernetesCustomResourceDefinitionSpec{
		Name:                     toString(metadata["name"]),
		Namespace:                toString(metadata["namespace"]),
		Group:                    toString(specMap["group"]),
		Scope:                    toString(specMap["scope"]),
		Names:                    copyStringInterfaceMap(asMapOrEmpty(specMap["names"])),
		Conversion:               copyStringInterfaceMap(asMapOrEmpty(specMap["conversion"])),
		Validation:               copyStringInterfaceMap(asMapOrEmpty(specMap["validation"])),
		AdditionalPrinterColumns: []map[string]interface{}{},
		Extensions:               map[string]interface{}{},
	}
	if versions, ok := specMap["versions"].([]interface{}); ok {
		for _, item := range versions {
			if version, ok := asMap(item); ok {
				crd.Versions = append(crd.Versions, version)
			}
		}
	}
	if columns, ok := specMap["additionalPrinterColumns"].([]interface{}); ok {
		for _, item := range columns {
			if column, ok := asMap(item); ok {
				crd.AdditionalPrinterColumns = append(crd.AdditionalPrinterColumns, column)
			}
		}
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		crd.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "spec":
		default:
			crd.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(crd.Extensions) == 0 {
		crd.Extensions = nil
	}
	if crd.Name == "" {
		crd.Name = toString(metadata["name"])
	}
	return crd
}

func asMapOrEmpty(value interface{}) map[string]interface{} {
	mapped, _ := asMap(value)
	if len(mapped) == 0 {
		return nil
	}
	return mapped
}

func kubernetesPriorityClassSpecsFromAny(value interface{}) map[string]*KubernetesPriorityClassSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesPriorityClassSpec:
		return cloneKubernetesPriorityClassSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesPriorityClassSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesPriorityClassSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesPriorityClassSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesPriorityClassSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesPriorityClassSpecsFromAny(typedList)
}

func kubernetesPriorityClassSpecFromMap(mapped map[string]interface{}) *KubernetesPriorityClassSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	class := &KubernetesPriorityClassSpec{
		Name:             toString(metadata["name"]),
		Namespace:        toString(metadata["namespace"]),
		Value:            int32(toInt(mapped["value"])),
		GlobalDefault:    toBool(mapped["globalDefault"]),
		Description:      toString(mapped["description"]),
		PreemptionPolicy: toString(mapped["preemptionPolicy"]),
		Extensions:       map[string]interface{}{},
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		class.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "value", "globalDefault", "description", "preemptionPolicy":
		default:
			class.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(class.Extensions) == 0 {
		class.Extensions = nil
	}
	if class.Name == "" {
		class.Name = toString(metadata["name"])
	}
	return class
}

func kubernetesRuntimeClassSpecsFromAny(value interface{}) map[string]*KubernetesRuntimeClassSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesRuntimeClassSpec:
		return cloneKubernetesRuntimeClassSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesRuntimeClassSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesRuntimeClassSpecFromMap(item); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesRuntimeClassSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesRuntimeClassSpecFromMap(mapped); spec != nil && spec.Name != "" {
				result[spec.Name] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesRuntimeClassSpecsFromAny(typedList)
}

func kubernetesRuntimeClassSpecFromMap(mapped map[string]interface{}) *KubernetesRuntimeClassSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	overhead, _ := asMap(mapped["overhead"])
	scheduling, _ := asMap(mapped["scheduling"])
	class := &KubernetesRuntimeClassSpec{
		Name:       toString(metadata["name"]),
		Namespace:  toString(metadata["namespace"]),
		Handler:    toString(mapped["handler"]),
		Overhead:   copyStringInterfaceMap(overhead),
		Scheduling: copyStringInterfaceMap(scheduling),
		Extensions: map[string]interface{}{},
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		class.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "handler", "overhead", "scheduling":
		default:
			class.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(class.Extensions) == 0 {
		class.Extensions = nil
	}
	if class.Name == "" {
		class.Name = toString(metadata["name"])
	}
	return class
}

func kubernetesOpaqueManifestSpecsFromAny(value interface{}) map[string]*KubernetesOpaqueManifestSpec {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]*KubernetesOpaqueManifestSpec:
		return cloneKubernetesOpaqueManifestSpecs(typed)
	case []map[string]interface{}:
		result := make(map[string]*KubernetesOpaqueManifestSpec, len(typed))
		for _, item := range typed {
			if spec := kubernetesOpaqueManifestSpecFromMap(item); spec != nil && spec.Name != "" {
				key := kubernetesDocumentKeyFromMap(spec.Raw)
				if key == "" {
					key = strings.Join([]string{spec.Kind, spec.Namespace, spec.Name}, "/")
				}
				result[key] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []interface{}:
		result := map[string]*KubernetesOpaqueManifestSpec{}
		for _, item := range typed {
			mapped, ok := asMap(item)
			if !ok {
				continue
			}
			if spec := kubernetesOpaqueManifestSpecFromMap(mapped); spec != nil && spec.Name != "" {
				key := kubernetesDocumentKeyFromMap(spec.Raw)
				if key == "" {
					key = strings.Join([]string{spec.Kind, spec.Namespace, spec.Name}, "/")
				}
				result[key] = spec
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var typedList []map[string]interface{}
	if err := json.Unmarshal(raw, &typedList); err != nil || len(typedList) == 0 {
		return nil
	}
	return kubernetesOpaqueManifestSpecsFromAny(typedList)
}

func kubernetesOpaqueManifestSpecsFromResources(values ...interface{}) map[string]*KubernetesOpaqueManifestSpec {
	resources := kubernetesExtensionResourceSlice(values...)
	if len(resources) == 0 {
		return nil
	}
	result := make(map[string]*KubernetesOpaqueManifestSpec, len(resources))
	for _, resource := range resources {
		spec := kubernetesOpaqueManifestSpecFromMap(resource)
		if spec == nil || spec.Name == "" {
			continue
		}
		key := kubernetesDocumentKeyFromMap(spec.Raw)
		if key == "" {
			key = strings.Join([]string{spec.Kind, spec.Namespace, spec.Name}, "/")
		}
		if key == "" {
			continue
		}
		result[key] = spec
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func kubernetesOpaqueManifestSpecsFromCanonical(canonical *CanonicalApplication) map[string]*KubernetesOpaqueManifestSpec {
	if canonical == nil || len(canonical.Resources) == 0 {
		return nil
	}
	result := map[string]*KubernetesOpaqueManifestSpec{}
	keys := make([]string, 0, len(canonical.Resources))
	for key := range canonical.Resources {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		resource := canonical.Resources[key]
		if resource == nil || resource.Platform != PlatformKubernetes {
			continue
		}
		if resource.Kind != ResourceKindRaw && resource.Kind != ResourceKindUnknown {
			continue
		}
		mapped, ok := asMap(resource.Raw)
		if !ok || len(mapped) == 0 {
			continue
		}
		spec := kubernetesOpaqueManifestSpecFromMap(mapped)
		if spec == nil {
			continue
		}
		if spec.Name == "" {
			spec.Name = resource.Name
		}
		if spec.Kind == "" {
			spec.Kind = resource.NativeKind
		}
		if spec.APIVersion == "" {
			spec.APIVersion = resource.APIVersion
		}
		if len(spec.Extensions) == 0 && len(resource.Extensions) > 0 {
			spec.Extensions = canonicalizeOpaqueExtensionMap(resource.Extensions)
		}
		manifestKey := kubernetesDocumentKeyFromMap(mapped)
		if manifestKey == "" {
			manifestKey = resource.ID
		}
		result[manifestKey] = spec
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func kubernetesOpaqueManifestSpecsForApplication(app *Application) map[string]*KubernetesOpaqueManifestSpec {
	if app == nil {
		return nil
	}
	result := map[string]*KubernetesOpaqueManifestSpec{}
	if len(app.KubernetesOpaqueManifests) > 0 {
		for key, manifest := range app.KubernetesOpaqueManifests {
			if manifest == nil {
				continue
			}
			cloned := cloneKubernetesOpaqueManifestSpecs(map[string]*KubernetesOpaqueManifestSpec{key: manifest})
			if cloned != nil {
				result[key] = cloned[key]
			}
		}
	}
	if canonical := canonicalForApplication(app); canonical != nil {
		for key, manifest := range kubernetesOpaqueManifestSpecsFromCanonical(canonical) {
			if manifest == nil {
				continue
			}
			if _, exists := result[key]; !exists {
				cloned := cloneKubernetesOpaqueManifestSpecs(map[string]*KubernetesOpaqueManifestSpec{key: manifest})
				if cloned != nil {
					result[key] = cloned[key]
				}
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func kubernetesOpaqueManifestSpecFromMap(mapped map[string]interface{}) *KubernetesOpaqueManifestSpec {
	if len(mapped) == 0 {
		return nil
	}
	metadata, _ := asMap(mapped["metadata"])
	specMap, _ := asMap(mapped["spec"])
	manifest := &KubernetesOpaqueManifestSpec{
		APIVersion: toString(mapped["apiVersion"]),
		Kind:       toString(mapped["kind"]),
		Name:       toString(metadata["name"]),
		Namespace:  toString(metadata["namespace"]),
		Metadata:   copyStringInterfaceMap(metadata),
		Spec:       copyStringInterfaceMap(specMap),
		Extensions: map[string]interface{}{},
	}
	if raw, ok := deepCopyValue(mapped).(map[string]interface{}); ok {
		manifest.Raw = raw
	}
	for key, value := range mapped {
		switch key {
		case "apiVersion", "kind", "metadata", "spec":
		default:
			manifest.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(manifest.Extensions) == 0 {
		manifest.Extensions = nil
	}
	return manifest
}

func propagateTypedKubernetesResourcesToExtensions(app *Application, key, composeKey string, specs interface{}) {
	if app == nil {
		return
	}
	items := typedKubernetesResourceSlice(specs)
	if len(items) == 0 {
		return
	}
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	if _, exists := app.Extensions[key]; !exists {
		app.Extensions[key] = deepCopyValue(items)
	}
	if _, exists := app.Extensions[composeKey]; !exists {
		app.Extensions[composeKey] = deepCopyValue(items)
	}
}

func typedKubernetesResourceSlice(specs interface{}) []interface{} {
	rv := reflect.ValueOf(specs)
	if !rv.IsValid() || rv.Kind() != reflect.Map || rv.Len() == 0 {
		return nil
	}
	names := rv.MapKeys()
	if len(names) == 0 {
		return nil
	}
	stringNames := make([]string, 0, len(names))
	for _, key := range names {
		if key.Kind() == reflect.String {
			stringNames = append(stringNames, key.String())
		}
	}
	sort.Strings(stringNames)
	result := make([]interface{}, 0, len(stringNames))
	for _, name := range stringNames {
		value := rv.MapIndex(reflect.ValueOf(name))
		if !value.IsValid() || (value.Kind() == reflect.Ptr && value.IsNil()) {
			continue
		}
		raw := rawKubernetesManifestFromTypedSpec(value.Interface())
		if len(raw) == 0 {
			continue
		}
		result = append(result, deepCopyValue(raw))
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func rawKubernetesManifestFromTypedSpec(spec interface{}) map[string]interface{} {
	if spec == nil {
		return nil
	}
	value := reflect.ValueOf(spec)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}
	field := value.FieldByName("Raw")
	if !field.IsValid() {
		return nil
	}
	raw, ok := field.Interface().(map[string]interface{})
	if !ok || len(raw) == 0 {
		return nil
	}
	return raw
}

func appendNamedSliceFromAny(value interface{}) []string {
	switch typed := value.(type) {
	case []string:
		return append([]string{}, typed...)
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			if mapped, ok := asMap(item); ok {
				if name := toString(mapped["name"]); name != "" {
					result = append(result, name)
				}
			} else if name := toString(item); name != "" {
				result = append(result, name)
			}
		}
		return result
	default:
		if name := toString(value); name != "" {
			return []string{name}
		}
		return nil
	}
}

func kubernetesServicePortFromMap(mapped map[string]interface{}) KubernetesServicePort {
	port := KubernetesServicePort{
		Name:        toString(mapped["name"]),
		Port:        toInt(mapped["port"]),
		TargetPort:  toString(mapped["targetPort"]),
		Protocol:    toString(mapped["protocol"]),
		NodePort:    toInt(mapped["nodePort"]),
		AppProtocol: toString(mapped["appProtocol"]),
		Extensions:  map[string]interface{}{},
	}
	for key, value := range mapped {
		switch key {
		case "name", "port", "targetPort", "protocol", "nodePort", "appProtocol":
		default:
			port.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(port.Extensions) == 0 {
		port.Extensions = nil
	}
	return port
}

func appendStringSliceFromAny(value interface{}) []string {
	switch typed := value.(type) {
	case []string:
		return append([]string{}, typed...)
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			if text := toString(item); text != "" {
				result = append(result, text)
			}
		}
		return result
	default:
		if text := toString(value); text != "" {
			return []string{text}
		}
		return nil
	}
}

func boolPtrFromAny(value interface{}) (*bool, bool) {
	switch typed := value.(type) {
	case bool:
		return &typed, true
	case string:
		if typed == "true" || typed == "false" {
			result := typed == "true"
			return &result, true
		}
	}
	return nil, false
}

func normalizeMeshSpec(mesh *MeshSpec) {
	if mesh == nil {
		return
	}
	if mesh.ControlPlaneURL == "" && len(mesh.ControlPlaneURLs) > 0 {
		mesh.ControlPlaneURL = mesh.ControlPlaneURLs[0]
	}
	if len(mesh.ControlPlaneURLs) == 0 && mesh.ControlPlaneURL != "" {
		mesh.ControlPlaneURLs = []string{mesh.ControlPlaneURL}
	}
	if len(mesh.ControlPlaneURLs) > 0 {
		mesh.ControlPlaneURLs = append([]string{}, mesh.ControlPlaneURLs...)
	}
}

func cloneEnvFileRefs(refs []EnvFileRef) []EnvFileRef {
	if len(refs) == 0 {
		return nil
	}
	cloned := make([]EnvFileRef, len(refs))
	for i, ref := range refs {
		cloned[i] = ref
		if ref.Required != nil {
			required := *ref.Required
			cloned[i].Required = &required
		}
		cloned[i].Extensions = copyStringInterfaceMap(ref.Extensions)
	}
	return cloned
}

func clonePortMappings(ports []PortMapping) []PortMapping {
	if len(ports) == 0 {
		return nil
	}
	cloned := make([]PortMapping, len(ports))
	for i, port := range ports {
		cloned[i] = port
		cloned[i].Extensions = copyStringInterfaceMap(port.Extensions)
	}
	return cloned
}

func cloneComposeModels(models map[string]*ComposeModel) map[string]*ComposeModel {
	if len(models) == 0 {
		return nil
	}
	cloned := make(map[string]*ComposeModel, len(models))
	for name, model := range models {
		if clonedModel := cloneComposeModel(model); clonedModel != nil {
			cloned[name] = clonedModel
		}
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func cloneComposeModel(model *ComposeModel) *ComposeModel {
	if model == nil {
		return nil
	}
	cloned := *model
	cloned.RuntimeFlags = append([]string{}, model.RuntimeFlags...)
	cloned.Extensions = canonicalizeOpaqueExtensionMap(model.Extensions)
	return &cloned
}

func cloneNetwork(network *Network) *Network {
	if network == nil {
		return nil
	}
	cloned := *network
	cloned.DriverOpts = copyStringMap(network.DriverOpts)
	cloned.ExternalExtensions = canonicalizeOpaqueExtensionMap(network.ExternalExtensions)
	cloned.EnableIPv4 = cloneBoolPtr(network.EnableIPv4)
	cloned.EnableIPv6 = cloneBoolPtr(network.EnableIPv6)
	cloned.IPAM = cloneIPAMConfig(network.IPAM)
	cloned.Labels = copyStringMap(network.Labels)
	cloned.Extensions = canonicalizeOpaqueExtensionMap(network.Extensions)
	return &cloned
}

func cloneVolume(volume *Volume) *Volume {
	if volume == nil {
		return nil
	}
	cloned := *volume
	cloned.DriverOpts = copyStringMap(volume.DriverOpts)
	cloned.ExternalExtensions = canonicalizeOpaqueExtensionMap(volume.ExternalExtensions)
	cloned.Labels = copyStringMap(volume.Labels)
	cloned.Extensions = canonicalizeOpaqueExtensionMap(volume.Extensions)
	return &cloned
}

func cloneConfig(config *Config) *Config {
	if config == nil {
		return nil
	}
	cloned := *config
	cloned.ExternalExtensions = canonicalizeOpaqueExtensionMap(config.ExternalExtensions)
	cloned.Labels = copyStringMap(config.Labels)
	cloned.Extensions = canonicalizeOpaqueExtensionMap(config.Extensions)
	return &cloned
}

func cloneSecret(secret *Secret) *Secret {
	if secret == nil {
		return nil
	}
	cloned := *secret
	cloned.DriverOpts = copyStringMap(secret.DriverOpts)
	cloned.ExternalExtensions = canonicalizeOpaqueExtensionMap(secret.ExternalExtensions)
	cloned.Labels = copyStringMap(secret.Labels)
	cloned.Extensions = canonicalizeOpaqueExtensionMap(secret.Extensions)
	return &cloned
}

func cloneDeploySpec(deploy *DeploySpec) *DeploySpec {
	if deploy == nil {
		return nil
	}
	cloned := *deploy
	cloned.Job = cloneSwarmJobSpec(deploy.Job)
	cloned.Placement = clonePlacementSpec(deploy.Placement)
	cloned.Resources = cloneResourceSpec(deploy.Resources)
	cloned.UpdateConfig = cloneUpdatePolicy(deploy.UpdateConfig)
	cloned.MigrateConfig = cloneMigratePolicy(deploy.MigrateConfig)
	cloned.RescheduleConfig = cloneReschedulePolicy(deploy.RescheduleConfig)
	cloned.RollbackConfig = cloneUpdatePolicy(deploy.RollbackConfig)
	cloned.RestartPolicy = cloneRestartPolicy(deploy.RestartPolicy)
	cloned.Labels = copyStringMap(deploy.Labels)
	cloned.Extensions = copyStringInterfaceMap(deploy.Extensions)
	return &cloned
}

func cloneSwarmJobSpec(job *SwarmJobSpec) *SwarmJobSpec {
	if job == nil {
		return nil
	}
	cloned := *job
	cloned.Extensions = copyStringInterfaceMap(job.Extensions)
	return &cloned
}

func isSwarmJobMode(mode string) bool {
	return strings.EqualFold(mode, "replicated-job") || strings.EqualFold(mode, "global-job")
}

func isEmptySwarmJobSpec(job *SwarmJobSpec) bool {
	return job == nil || (!job.maxConcurrentSet && job.MaxConcurrent == 0 && !job.totalCompletionsSet && job.TotalCompletions == 0 && !job.completionModeSet && job.CompletionMode == "" && job.Suspend == nil && !job.backoffLimitSet && job.BackoffLimit == 0 && !job.backoffLimitPerIndexSet && job.BackoffLimitPerIndex == 0 && !job.ttlSecondsAfterFinishedSet && job.TTLSecondsAfterFinished == 0 && len(job.Extensions) == 0)
}

func cloneComposeCompat(compat *ComposeCompat) *ComposeCompat {
	if compat == nil {
		return nil
	}
	cloned := *compat
	cloned.Attach = cloneBoolPtr(compat.Attach)
	cloned.Annotations = copyStringMap(compat.Annotations)
	cloned.BlkioConfig = copyStringInterfaceMap(compat.BlkioConfig)
	cloned.CredentialSpec = copyStringInterfaceMap(compat.CredentialSpec)
	cloned.Provider = copyStringInterfaceMap(compat.Provider)
	cloned.Extends = copyStringInterfaceMap(compat.Extends)
	cloned.Scale = cloneIntPtr(compat.Scale)
	cloned.StorageOpt = copyStringMap(compat.StorageOpt)
	cloned.VolumesFrom = append([]string{}, compat.VolumesFrom...)
	cloned.DeviceCgroupRules = append([]string{}, compat.DeviceCgroupRules...)
	cloned.Gpus = cloneMapSlice(compat.Gpus)
	cloned.Models = cloneComposeModelMap(compat.Models)
	cloned.ExternalLinks = append([]string{}, compat.ExternalLinks...)
	cloned.LabelFiles = append([]string{}, compat.LabelFiles...)
	cloned.Tmpfs = append([]string{}, compat.Tmpfs...)
	cloned.Extensions = copyStringInterfaceMap(compat.Extensions)
	return &cloned
}

func cloneComposeModelMap(input map[string]map[string]interface{}) map[string]map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	copied := make(map[string]map[string]interface{}, len(input))
	for key, value := range input {
		copied[key] = canonicalizeOpaqueExtensionMap(value)
	}
	return copied
}

func clonePlacementSpec(placement *PlacementSpec) *PlacementSpec {
	if placement == nil {
		return nil
	}
	cloned := *placement
	cloned.Constraints = append([]string{}, placement.Constraints...)
	cloned.Preferences = append([]string{}, placement.Preferences...)
	cloned.Extensions = copyStringInterfaceMap(placement.Extensions)
	if len(placement.PreferenceExtensions) > 0 {
		cloned.PreferenceExtensions = make([]map[string]interface{}, len(placement.PreferenceExtensions))
		for i, extensions := range placement.PreferenceExtensions {
			cloned.PreferenceExtensions[i] = copyStringInterfaceMap(extensions)
		}
	}
	return &cloned
}

func cloneResourceSpec(resources *ResourceSpec) *ResourceSpec {
	if resources == nil {
		return nil
	}
	cloned := *resources
	cloned.Extensions = copyStringInterfaceMap(resources.Extensions)
	cloned.LimitExtensions = copyStringInterfaceMap(resources.LimitExtensions)
	cloned.ReservationExtensions = copyStringInterfaceMap(resources.ReservationExtensions)
	cloned.Devices = cloneResourceDevices(resources.Devices)
	cloned.GenericResources = cloneGenericResources(resources.GenericResources)
	return &cloned
}

func cloneResourceDevices(devices []ResourceDevice) []ResourceDevice {
	if len(devices) == 0 {
		return nil
	}
	cloned := make([]ResourceDevice, len(devices))
	for i, device := range devices {
		cloned[i] = device
		cloned[i].Capabilities = append([]string{}, device.Capabilities...)
		cloned[i].DeviceIDs = append([]string{}, device.DeviceIDs...)
		cloned[i].Options = copyStringMap(device.Options)
		cloned[i].Extensions = copyStringInterfaceMap(device.Extensions)
	}
	return cloned
}

func cloneGenericResources(resources []GenericResource) []GenericResource {
	if len(resources) == 0 {
		return nil
	}
	cloned := make([]GenericResource, len(resources))
	for i, resource := range resources {
		cloned[i] = resource
		cloned[i].Extensions = copyStringInterfaceMap(resource.Extensions)
		cloned[i].DiscreteExtensions = copyStringInterfaceMap(resource.DiscreteExtensions)
	}
	return cloned
}

func cloneUpdatePolicy(policy *UpdatePolicy) *UpdatePolicy {
	if policy == nil {
		return nil
	}
	cloned := *policy
	cloned.Extensions = copyStringInterfaceMap(policy.Extensions)
	return &cloned
}

func cloneMigratePolicy(policy *MigratePolicy) *MigratePolicy {
	if policy == nil {
		return nil
	}
	cloned := *policy
	cloned.Extensions = copyStringInterfaceMap(policy.Extensions)
	return &cloned
}

func cloneReschedulePolicy(policy *ReschedulePolicy) *ReschedulePolicy {
	if policy == nil {
		return nil
	}
	cloned := *policy
	cloned.Extensions = copyStringInterfaceMap(policy.Extensions)
	return &cloned
}

func cloneRestartPolicy(policy *RestartPolicy) *RestartPolicy {
	if policy == nil {
		return nil
	}
	cloned := *policy
	cloned.Extensions = copyStringInterfaceMap(policy.Extensions)
	return &cloned
}

func cloneHealthCheck(health *HealthCheck) *HealthCheck {
	if health == nil {
		return nil
	}
	cloned := *health
	cloned.Test = append([]string{}, health.Test...)
	cloned.Extensions = copyStringInterfaceMap(health.Extensions)
	return &cloned
}

func cloneBuildConfig(build *BuildConfig) *BuildConfig {
	if build == nil {
		return nil
	}
	cloned := *build
	cloned.Extensions = copyStringInterfaceMap(build.Extensions)
	return &cloned
}

func cloneSeccompProfile(profile *SeccompProfile) *SeccompProfile {
	if profile == nil {
		return nil
	}
	cloned := *profile
	cloned.Extensions = copyStringInterfaceMap(profile.Extensions)
	return &cloned
}

func cloneSELinuxOptions(options *SELinuxOptions) *SELinuxOptions {
	if options == nil {
		return nil
	}
	cloned := *options
	cloned.Extensions = copyStringInterfaceMap(options.Extensions)
	return &cloned
}

func cloneWindowsSecurityContextOptions(options *WindowsSecurityContextOptions) *WindowsSecurityContextOptions {
	if options == nil {
		return nil
	}
	cloned := *options
	cloned.Extensions = copyStringInterfaceMap(options.Extensions)
	if options.GMSACredentialSpecName != nil {
		value := *options.GMSACredentialSpecName
		cloned.GMSACredentialSpecName = &value
	}
	if options.GMSACredentialSpec != nil {
		value := *options.GMSACredentialSpec
		cloned.GMSACredentialSpec = &value
	}
	if options.RunAsUserName != nil {
		value := *options.RunAsUserName
		cloned.RunAsUserName = &value
	}
	if options.HostProcess != nil {
		value := *options.HostProcess
		cloned.HostProcess = &value
	}
	return &cloned
}

func cloneMapSlice(values []map[string]interface{}) []map[string]interface{} {
	if len(values) == 0 {
		return nil
	}
	cloned := make([]map[string]interface{}, len(values))
	for i, value := range values {
		cloned[i] = cloneMap(value)
	}
	return cloned
}

func cloneDevelopConfig(develop *DevelopConfig) *DevelopConfig {
	if develop == nil {
		return nil
	}
	cloned := *develop
	if len(develop.Watch) > 0 {
		cloned.Watch = make([]DevelopWatch, len(develop.Watch))
		for i, watch := range develop.Watch {
			cloned.Watch[i] = watch
			cloned.Watch[i].Exec = cloneServiceHook(watch.Exec)
			cloned.Watch[i].Include = append([]string{}, watch.Include...)
			cloned.Watch[i].Ignore = append([]string{}, watch.Ignore...)
			cloned.Watch[i].Extensions = copyStringInterfaceMap(watch.Extensions)
		}
	}
	cloned.Extensions = copyStringInterfaceMap(develop.Extensions)
	return &cloned
}

func cloneLifecycleHooks(lifecycle *LifecycleHooks) *LifecycleHooks {
	if lifecycle == nil {
		return nil
	}
	cloned := *lifecycle
	cloned.PreStart = cloneServiceHooks(lifecycle.PreStart)
	cloned.PostStart = cloneServiceHooks(lifecycle.PostStart)
	cloned.PreStop = cloneServiceHooks(lifecycle.PreStop)
	cloned.Extensions = copyStringInterfaceMap(lifecycle.Extensions)
	return &cloned
}

func cloneServiceHooks(hooks []ServiceHook) []ServiceHook {
	if len(hooks) == 0 {
		return nil
	}
	cloned := make([]ServiceHook, len(hooks))
	for i, hook := range hooks {
		cloned[i] = *cloneServiceHook(&hook)
	}
	return cloned
}

func cloneServiceHook(hook *ServiceHook) *ServiceHook {
	if hook == nil {
		return nil
	}
	cloned := *hook
	cloned.Command = append([]string{}, hook.Command...)
	cloned.Environment = copyStringPtrMap(hook.Environment)
	cloned.Extensions = copyStringInterfaceMap(hook.Extensions)
	return &cloned
}

func cloneUlimits(limits *Ulimits) *Ulimits {
	if limits == nil {
		return nil
	}
	cloned := *limits
	if limits.Nofile != nil {
		nofile := *limits.Nofile
		nofile.Extensions = copyStringInterfaceMap(limits.Nofile.Extensions)
		cloned.Nofile = &nofile
	}
	cloned.Limits = cloneUlimitSpecs(limits.Limits)
	return &cloned
}

func cloneUlimitSpecs(limits map[string]UlimitSpec) map[string]UlimitSpec {
	if len(limits) == 0 {
		return nil
	}
	cloned := make(map[string]UlimitSpec, len(limits))
	for name, limit := range limits {
		copied := limit
		copied.Extensions = copyStringInterfaceMap(limit.Extensions)
		cloned[name] = copied
	}
	return cloned
}

func cloneIPAMConfig(ipam *IPAMConfig) *IPAMConfig {
	if ipam == nil {
		return nil
	}
	cloned := *ipam
	if len(ipam.Config) > 0 {
		cloned.Config = make([]IPAMSubnet, len(ipam.Config))
		for i, subnet := range ipam.Config {
			cloned.Config[i] = subnet
			cloned.Config[i].AuxAddresses = copyStringMap(subnet.AuxAddresses)
			cloned.Config[i].Extensions = copyStringInterfaceMap(subnet.Extensions)
		}
	}
	cloned.Options = copyStringMap(ipam.Options)
	cloned.Extensions = copyStringInterfaceMap(ipam.Extensions)
	return &cloned
}

func cloneVolumeMounts(volumes []VolumeMount) []VolumeMount {
	if len(volumes) == 0 {
		return nil
	}
	cloned := make([]VolumeMount, len(volumes))
	for i, volume := range volumes {
		cloned[i] = volume
		cloned[i].CreateHostPath = cloneBoolPtr(volume.CreateHostPath)
		cloned[i].VolumeLabels = copyStringMap(volume.VolumeLabels)
		cloned[i].Options = copyStringMap(volume.Options)
		cloned[i].Extensions = copyStringInterfaceMap(volume.Extensions)
		cloned[i].BindExtensions = copyStringInterfaceMap(volume.BindExtensions)
		cloned[i].VolumeExtensions = copyStringInterfaceMap(volume.VolumeExtensions)
		cloned[i].TmpfsExtensions = copyStringInterfaceMap(volume.TmpfsExtensions)
		cloned[i].ImageExtensions = copyStringInterfaceMap(volume.ImageExtensions)
	}
	return cloned
}

func mergePortableVolumeMounts(existing, preserved []VolumeMount) []VolumeMount {
	if len(preserved) == 0 {
		return existing
	}
	if len(existing) == 0 {
		return cloneVolumeMounts(preserved)
	}
	merged := cloneVolumeMounts(existing)
	for _, preservedMount := range preserved {
		next := make([]VolumeMount, 0, len(merged)+1)
		replaced := false
		clonedPreserved := cloneVolumeMounts([]VolumeMount{preservedMount})[0]
		for _, existingMount := range merged {
			if !volumeMountSameMount(existingMount, preservedMount) {
				next = append(next, existingMount)
				continue
			}
			if !replaced {
				next = append(next, clonedPreserved)
				replaced = true
			}
		}
		if !replaced {
			next = append(next, clonedPreserved)
		}
		merged = next
	}
	return merged
}

func volumeMountSameMount(a, b VolumeMount) bool {
	if a.Target != "" && b.Target != "" && a.Target == b.Target {
		return true
	}
	return a.Source != "" && b.Source != "" && a.Source == b.Source
}

func cloneFileRefs(refs []FileRef) []FileRef {
	if len(refs) == 0 {
		return nil
	}
	cloned := make([]FileRef, len(refs))
	for i, ref := range refs {
		cloned[i] = ref
		cloned[i].Optional = cloneBoolPtr(ref.Optional)
		cloned[i].Extensions = copyStringInterfaceMap(ref.Extensions)
	}
	return cloned
}

func cloneDependencySpecs(specs []DependencySpec) []DependencySpec {
	if len(specs) == 0 {
		return nil
	}
	cloned := make([]DependencySpec, len(specs))
	for i, spec := range specs {
		cloned[i] = spec
		cloned[i].Required = cloneBoolPtr(spec.Required)
		cloned[i].Extensions = copyStringInterfaceMap(spec.Extensions)
	}
	return cloned
}

func cloneBoolPtr(value *bool) *bool {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func copyStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	copied := make(map[string]string, len(input))
	for key, value := range input {
		copied[key] = value
	}
	return copied
}

func copyStringPtrMap(input map[string]*string) map[string]*string {
	if len(input) == 0 {
		return nil
	}
	copied := make(map[string]*string, len(input))
	for key, value := range input {
		if value == nil {
			copied[key] = nil
			continue
		}
		cloned := *value
		copied[key] = &cloned
	}
	return copied
}

func copyStringInterfaceMap(input map[string]interface{}) map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	copied := make(map[string]interface{}, len(input))
	for key, value := range input {
		copied[key] = deepCopyValue(value)
	}
	return copied
}

func deepCopyValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case *Service:
		return cloneService(typed)
	case *Network:
		return cloneNetwork(typed)
	case *Volume:
		return cloneVolume(typed)
	case *Config:
		return cloneConfig(typed)
	case *Secret:
		return cloneSecret(typed)
	case *DevelopConfig:
		return cloneDevelopConfig(typed)
	case *LifecycleHooks:
		return cloneLifecycleHooks(typed)
	case *ServiceHook:
		return cloneServiceHook(typed)
	case *ComposeModel:
		return cloneComposeModel(typed)
	case *RouteSpec:
		return cloneRouteSpec(typed)
	case *PolicySpec:
		return clonePolicySpec(typed)
	case ComposeModel:
		return cloneComposeModel(&typed)
	case RouteSpec:
		return cloneRouteSpec(&typed)
	case PolicySpec:
		return clonePolicySpec(&typed)
	case DevelopConfig:
		return cloneDevelopConfig(&typed)
	case LifecycleHooks:
		return cloneLifecycleHooks(&typed)
	case ServiceHook:
		return cloneServiceHook(&typed)
	case map[string]interface{}:
		return copyStringInterfaceMap(typed)
	case []interface{}:
		copied := make([]interface{}, len(typed))
		for i, item := range typed {
			copied[i] = deepCopyValue(item)
		}
		return copied
	case map[string]string:
		return copyStringMap(typed)
	case map[string]*string:
		return copyStringPtrMap(typed)
	case []string:
		return append([]string{}, typed...)
	case []map[string]interface{}:
		copied := make([]map[string]interface{}, len(typed))
		for i, item := range typed {
			copied[i] = copyStringInterfaceMap(item)
		}
		return copied
	default:
		return typed
	}
}

func cloneInterfaceSlice(values []interface{}) []interface{} {
	if len(values) == 0 {
		return nil
	}
	cloned := make([]interface{}, len(values))
	for i, value := range values {
		cloned[i] = deepCopyValue(value)
	}
	return cloned
}

func mergeIncludeEntries(existing, incoming []interface{}) []interface{} {
	if len(incoming) == 0 {
		return cloneInterfaceSlice(existing)
	}
	merged := cloneInterfaceSlice(existing)
	seen := map[string]struct{}{}
	for _, entry := range merged {
		seen[stableJSONKey(entry)] = struct{}{}
	}
	for _, entry := range incoming {
		key := stableJSONKey(entry)
		if _, ok := seen[key]; ok {
			continue
		}
		merged = append(merged, deepCopyValue(entry))
		seen[key] = struct{}{}
	}
	return merged
}

func stableJSONKey(value interface{}) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%#v", value)
	}
	return string(raw)
}

func cloneRouteSpec(route *RouteSpec) *RouteSpec {
	if route == nil {
		return nil
	}
	cloned := *route
	cloned.Hosts = append([]string{}, route.Hosts...)
	cloned.Paths = append([]string{}, route.Paths...)
	cloned.Entrypoints = append([]string{}, route.Entrypoints...)
	cloned.Middlewares = append([]string{}, route.Middlewares...)
	cloned.Metadata = copyStringMap(route.Metadata)
	cloned.Annotations = copyStringMap(route.Annotations)
	cloned.Extensions = copyStringInterfaceMap(route.Extensions)
	cloned.Raw = deepCopyValue(route.Raw)
	return &cloned
}

func clonePortableRouteSpec(route *RouteSpec) *RouteSpec {
	cloned := cloneRouteSpec(route)
	if cloned != nil {
		cloned.Raw = nil
	}
	return cloned
}

func cloneRouteSpecMap(routes map[string]*RouteSpec) map[string]*RouteSpec {
	if len(routes) == 0 {
		return nil
	}
	cloned := make(map[string]*RouteSpec, len(routes))
	for name, route := range routes {
		if clonedRoute := cloneRouteSpec(route); clonedRoute != nil {
			cloned[name] = clonedRoute
		}
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func sortedRouteSpecList(routes map[string]*RouteSpec) []*RouteSpec {
	if len(routes) == 0 {
		return nil
	}
	names := make([]string, 0, len(routes))
	for name, route := range routes {
		if route != nil {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	result := make([]*RouteSpec, 0, len(names))
	for _, name := range names {
		result = append(result, cloneRouteSpec(routes[name]))
	}
	return result
}

func clonePolicySpec(policy *PolicySpec) *PolicySpec {
	if policy == nil {
		return nil
	}
	cloned := *policy
	cloned.Rules = append([]string{}, policy.Rules...)
	cloned.Metadata = copyStringMap(policy.Metadata)
	cloned.Annotations = copyStringMap(policy.Annotations)
	cloned.Extension = canonicalizeStringOpaqueExtensionMap(policy.Extension)
	cloned.Extensions = copyStringInterfaceMap(policy.Extensions)
	cloned.Raw = deepCopyValue(policy.Raw)
	return &cloned
}

func clonePortablePolicySpec(policy *PolicySpec) *PolicySpec {
	cloned := clonePolicySpec(policy)
	if cloned != nil {
		cloned.Raw = nil
	}
	return cloned
}

func clonePolicySpecMap(policies map[string]*PolicySpec) map[string]*PolicySpec {
	if len(policies) == 0 {
		return nil
	}
	cloned := make(map[string]*PolicySpec, len(policies))
	for name, policy := range policies {
		if clonedPolicy := clonePolicySpec(policy); clonedPolicy != nil {
			cloned[name] = clonedPolicy
		}
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func sortedPolicySpecList(policies map[string]*PolicySpec) []*PolicySpec {
	if len(policies) == 0 {
		return nil
	}
	names := make([]string, 0, len(policies))
	for name, policy := range policies {
		if policy != nil {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	result := make([]*PolicySpec, 0, len(names))
	for _, name := range names {
		result = append(result, clonePolicySpec(policies[name]))
	}
	return result
}

// convertServiceAttributes converts platform-specific service attributes
func (p *PaaS) convertServiceAttributes(service *Service, from, to Platform) error {
	// Handle platform-specific conversions
	switch from {
	case PlatformDockerCompose, PlatformDockerSwarm:
		switch to {
		case PlatformNomad:
			return p.convertDockerToNomad(service)
		case PlatformKubernetes:
			return p.convertDockerToKubernetes(service)
		case PlatformDockerCompose, PlatformDockerSwarm:
			return nil
		}
	case PlatformNomad:
		switch to {
		case PlatformDockerCompose, PlatformDockerSwarm:
			return p.convertNomadToDocker(service)
		case PlatformKubernetes:
			return p.convertNomadToKubernetes(service)
		}
	case PlatformKubernetes:
		switch to {
		case PlatformDockerCompose, PlatformDockerSwarm:
			return p.convertKubernetesToDocker(service)
		case PlatformNomad:
			return p.convertKubernetesToNomad(service)
		}
	}

	return nil
}

// convertDockerToNomad converts Docker Compose service attributes to Nomad
func (p *PaaS) convertDockerToNomad(service *Service) error {
	ensureDeployFromReplicas(service)

	// Convert restart policies
	switch service.Restart {
	case "always":
		service.Restart = "" // Nomad default is restart
	case "no":
		service.Restart = "no"
	case "unless-stopped":
		service.Restart = "" // Nomad handles this differently
	}

	return nil
}

// convertDockerToKubernetes converts Docker Compose to Kubernetes
func (p *PaaS) convertDockerToKubernetes(service *Service) error {
	ensureDeployFromReplicas(service)

	// Convert restart policies
	switch service.Restart {
	case "always":
		service.Restart = "Always"
	case "no":
		service.Restart = "Never"
	case "unless-stopped":
		service.Restart = "OnFailure"
	}

	return nil
}

// convertNomadToDocker converts Nomad service attributes to Docker Compose
func (p *PaaS) convertNomadToDocker(service *Service) error {
	// Nomad restart policies are different from Docker
	// Nomad handles restarts at the job level, not service level
	service.Restart = "unless-stopped" // Default Docker behavior
	ensureDeployFromReplicas(service)

	return nil
}

// convertNomadToKubernetes converts Nomad to Kubernetes
func (p *PaaS) convertNomadToKubernetes(service *Service) error {
	// Similar to Docker conversion
	return p.convertNomadToDocker(service)
}

// convertKubernetesToDocker converts Kubernetes service attributes to Docker Compose
func (p *PaaS) convertKubernetesToDocker(service *Service) error {
	// Convert restart policies
	switch service.Restart {
	case "Always":
		service.Restart = "always"
	case "Never":
		service.Restart = "no"
	case "OnFailure":
		service.Restart = "unless-stopped"
	}
	ensureDeployFromReplicas(service)

	return nil
}

func stripKubernetesServiceExtensionKeys(app *Application) {
	if app == nil {
		return
	}
	for _, service := range app.Services {
		if service == nil || len(service.Extensions) == 0 {
			continue
		}
		for key := range service.Extensions {
			if strings.HasPrefix(key, "kubernetes.") {
				delete(service.Extensions, key)
			}
		}
	}
}

func ensureDeployFromReplicas(service *Service) {
	if service.Replicas <= 0 {
		return
	}
	if service.Deploy == nil {
		service.Deploy = &DeploySpec{}
	}
	if service.Deploy.Replicas == 0 {
		service.Deploy.Replicas = service.Replicas
	}
}

// convertKubernetesToNomad converts Kubernetes to Nomad
func (p *PaaS) convertKubernetesToNomad(service *Service) error {
	// Similar to Docker conversion
	return p.convertKubernetesToDocker(service)
}

// detectPlatform detects the platform from filename extension
// Validate validates an application
func (p *PaaS) Validate(app *Application) error {
	return app.Validate()
}

// RoundTrip performs a round-trip conversion test
func (p *PaaS) RoundTrip(app *Application, platforms ...Platform) error {
	original := app.Platform

	for _, targetPlatform := range platforms {
		converted, err := p.Convert(app, original, targetPlatform)
		if err != nil {
			return fmt.Errorf("failed to convert %s -> %s: %w", original, targetPlatform, err)
		}

		// Convert back
		backConverted, err := p.Convert(converted, targetPlatform, original)
		if err != nil {
			return fmt.Errorf("failed to convert back %s -> %s: %w", targetPlatform, original, err)
		}

		// Validate the round-trip
		if err := p.Validate(backConverted); err != nil {
			return fmt.Errorf("round-trip validation failed for %s: %w", targetPlatform, err)
		}

		// Basic structural comparison
		if len(backConverted.Services) != len(app.Services) {
			return fmt.Errorf("round-trip failed: service count mismatch for %s", targetPlatform)
		}
	}

	return nil
}

// ListServices returns a list of service names
func (p *PaaS) ListServices(app *Application) []string {
	var names []string
	for name := range app.Services {
		names = append(names, name)
	}
	return names
}

// GetService returns a service by name
func (p *PaaS) GetService(app *Application, name string) (*Service, bool) {
	service, exists := app.Services[name]
	return service, exists
}

// AddService adds a service to the application
func (p *PaaS) AddService(app *Application, service *Service) {
	app.Services[service.Name] = service
}

// RemoveService removes a service from the application
func (p *PaaS) RemoveService(app *Application, name string) {
	delete(app.Services, name)
}

// MergeApplications merges multiple applications
func (p *PaaS) MergeApplications(apps ...*Application) (*Application, error) {
	if len(apps) == 0 {
		return nil, fmt.Errorf("no applications to merge")
	}

	merged := &Application{
		Version:  apps[0].Version,
		Platform: apps[0].Platform,
		Services: make(map[string]*Service),
		Networks: make(map[string]*Network),
		Volumes:  make(map[string]*Volume),
		Configs:  make(map[string]*Config),
		Secrets:  make(map[string]*Secret),
		Models:   make(map[string]*ComposeModel),
	}

	for _, app := range apps {
		if merged.Name == "" && app.Name != "" {
			merged.Name = app.Name
		}
		if merged.Version == "" && app.Version != "" {
			merged.Version = app.Version
		}
		merged.Includes = mergeUniqueStrings(merged.Includes, app.Includes)
		merged.IncludeEntries = mergeIncludeEntries(merged.IncludeEntries, app.IncludeEntries)
		merged.SourceFiles = mergeUniqueStrings(merged.SourceFiles, app.SourceFiles)
		merged.Extensions = mergeApplicationExtensionMaps(merged.Extensions, app.Extensions)
		if merged.Mesh == nil && app.Mesh != nil {
			merged.Mesh = cloneMeshSpec(app.Mesh)
		}

		// Merge services
		for name, service := range app.Services {
			if _, exists := merged.Services[name]; exists {
				return nil, fmt.Errorf("service %s already exists", name)
			}
			merged.Services[name] = cloneService(service)
		}

		// Merge networks
		for name, network := range app.Networks {
			if _, exists := merged.Networks[name]; !exists {
				merged.Networks[name] = cloneNetwork(network)
			}
		}

		// Merge volumes
		for name, volume := range app.Volumes {
			if _, exists := merged.Volumes[name]; !exists {
				merged.Volumes[name] = cloneVolume(volume)
			}
		}

		// Merge configs
		for name, config := range app.Configs {
			if _, exists := merged.Configs[name]; !exists {
				merged.Configs[name] = cloneConfig(config)
			}
		}

		// Merge secrets
		for name, secret := range app.Secrets {
			if _, exists := merged.Secrets[name]; !exists {
				merged.Secrets[name] = cloneSecret(secret)
			}
		}

		// Merge compose-spec models
		for name, model := range app.Models {
			if _, exists := merged.Models[name]; !exists {
				merged.Models[name] = cloneComposeModel(model)
			}
		}
	}

	merged.AttachCanonical()
	return merged, nil
}

func mergeUniqueStrings(existing, incoming []string) []string {
	if len(incoming) == 0 {
		if len(existing) == 0 {
			return nil
		}
		return existing
	}
	if len(existing) == 0 {
		return append([]string{}, incoming...)
	}
	seen := make(map[string]struct{}, len(existing)+len(incoming))
	merged := append([]string{}, existing...)
	for _, value := range existing {
		seen[value] = struct{}{}
	}
	for _, value := range incoming {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		merged = append(merged, value)
	}
	return merged
}

func mergeApplicationExtensionMaps(existing, incoming map[string]interface{}) map[string]interface{} {
	if len(incoming) == 0 {
		return existing
	}
	incoming = canonicalizeApplicationExtensionMap(incoming)
	if len(existing) == 0 {
		return copyStringInterfaceMap(incoming)
	}
	merged := canonicalizeApplicationExtensionMap(existing)
	for key, value := range incoming {
		current, ok := merged[key]
		if !ok {
			merged[key] = deepCopyValue(value)
			continue
		}
		merged[key] = mergeApplicationExtensionValue(current, value)
	}
	return merged
}

func canonicalizeApplicationExtensionMap(extensions map[string]interface{}) map[string]interface{} {
	if len(extensions) == 0 {
		return nil
	}
	result := map[string]interface{}{}
	for key, value := range extensions {
		if canonical := composeApplicationCanonicalKey(key); canonical != "" {
			result[canonical] = deepCopyValue(value)
			continue
		}
		result[key] = deepCopyValue(value)
	}
	return result
}

func canonicalizeServiceExtensionMap(extensions map[string]interface{}) map[string]interface{} {
	return canonicalizeOpaqueExtensionMap(extensions)
}

func canonicalizeOpaqueExtensionMap(extensions map[string]interface{}) map[string]interface{} {
	if len(extensions) == 0 {
		return nil
	}
	result := copyStringInterfaceMap(extensions)
	for key, value := range extensions {
		if canonical := composeApplicationCanonicalKey(key); canonical != "" {
			result[canonical] = deepCopyValue(value)
		}
	}
	return result
}

func canonicalizeStringOpaqueExtensionMap(extensions map[string]string) map[string]string {
	if len(extensions) == 0 {
		return nil
	}
	result := copyStringMap(extensions)
	for key, value := range extensions {
		if canonical := composeApplicationCanonicalKey(key); canonical != "" {
			result[canonical] = value
		}
	}
	return result
}

func mergeApplicationExtensionValue(existing, incoming interface{}) interface{} {
	switch current := existing.(type) {
	case []interface{}:
		added, ok := extensionSliceValues(incoming)
		if !ok {
			return current
		}
		merged := append([]interface{}{}, current...)
		for _, item := range added {
			if !containsDeepEqual(merged, item) {
				merged = append(merged, deepCopyValue(item))
			}
		}
		return merged
	case []string:
		switch added := incoming.(type) {
		case []string:
			return mergeUniqueStrings(current, added)
		case []interface{}:
			merged := make([]interface{}, 0, len(current)+len(added))
			for _, item := range current {
				merged = append(merged, item)
			}
			for _, item := range added {
				if !containsDeepEqual(merged, item) {
					merged = append(merged, deepCopyValue(item))
				}
			}
			return merged
		default:
			return current
		}
	case map[string]interface{}:
		added, ok := extensionMapValues(incoming)
		if !ok {
			return current
		}
		return mergeApplicationExtensionMaps(current, added)
	case map[string]string:
		switch added := incoming.(type) {
		case map[string]string:
			merged := copyStringMap(current)
			for key, value := range added {
				if _, exists := merged[key]; !exists {
					merged[key] = value
				}
			}
			return merged
		case map[string]interface{}:
			merged := make(map[string]interface{}, len(current)+len(added))
			for key, value := range current {
				merged[key] = value
			}
			return mergeApplicationExtensionMaps(merged, added)
		default:
			return current
		}
	default:
		return current
	}
}

func extensionSliceValues(value interface{}) ([]interface{}, bool) {
	switch typed := value.(type) {
	case []interface{}:
		return typed, true
	case []map[string]interface{}:
		converted := make([]interface{}, len(typed))
		for i, item := range typed {
			converted[i] = item
		}
		return converted, true
	case []string:
		converted := make([]interface{}, len(typed))
		for i, item := range typed {
			converted[i] = item
		}
		return converted, true
	default:
		return nil, false
	}
}

func extensionMapValues(value interface{}) (map[string]interface{}, bool) {
	switch typed := value.(type) {
	case map[string]interface{}:
		return typed, true
	case map[string]string:
		converted := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			converted[key] = item
		}
		return converted, true
	default:
		return nil, false
	}
}

func containsDeepEqual(values []interface{}, candidate interface{}) bool {
	for _, value := range values {
		if reflect.DeepEqual(value, candidate) {
			return true
		}
	}
	return false
}
