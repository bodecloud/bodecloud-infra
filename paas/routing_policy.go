package paas

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var traefikHostPattern = regexp.MustCompile("Host\\(`([^`]+)`\\)")

const kubernetesPortableRouteAnnotation = "bolabaden.dev/portable-route"
const kubernetesPortablePolicyAnnotation = "bolabaden.dev/portable-policy"

func routesFromServiceLabels(platform Platform, serviceName string, labels map[string]string) []*RouteSpec {
	if len(labels) == 0 || !labelTruthy(labels["traefik.enable"]) {
		return nil
	}

	routeNames := map[string]struct{}{}
	const routerPrefix = "traefik.http.routers."
	for key := range labels {
		if !strings.HasPrefix(key, routerPrefix) {
			continue
		}
		rest := strings.TrimPrefix(key, routerPrefix)
		parts := strings.SplitN(rest, ".", 2)
		if len(parts) == 2 {
			routeNames[parts[0]] = struct{}{}
		}
	}

	var routes []*RouteSpec
	for name := range routeNames {
		route := &RouteSpec{
			Name:     name,
			Service:  serviceName,
			Protocol: "http",
			Source:   platform,
			Raw:      labels,
			Metadata: map[string]string{"provider": "traefik"},
		}
		base := routerPrefix + name + "."
		route.Hosts = hostsFromTraefikRule(labels[base+"rule"])
		route.Paths = pathsFromTraefikRule(labels[base+"rule"])
		route.Entrypoints = splitCSV(labels[base+"entrypoints"])
		route.Middlewares = splitCSV(labels[base+"middlewares"])
		route.TLS = labelTruthy(labels[base+"tls"])
		if serviceLabel := labels[base+"service"]; serviceLabel != "" {
			route.Metadata["traefik.service"] = serviceLabel
		}
		routes = append(routes, route)
	}
	sort.Slice(routes, func(i, j int) bool { return routes[i].Name < routes[j].Name })
	return routes
}

func policiesFromServiceLabels(platform Platform, serviceName string, labels map[string]string) []*PolicySpec {
	if len(labels) == 0 {
		return nil
	}
	policies := map[string]*PolicySpec{}
	for key, value := range labels {
		if strings.Contains(key, ".middlewares") && value != "" {
			name := serviceName + "-middlewares"
			policies[name] = &PolicySpec{
				Name:     name,
				Type:     "middleware-chain",
				Target:   serviceName,
				Provider: "traefik",
				Rules:    splitCSV(value),
				Source:   platform,
				Raw:      labels,
			}
		}
		if strings.Contains(key, "ratelimit") || strings.Contains(key, "rate-limit") {
			name := serviceName + "-rate-limit"
			policies[name] = &PolicySpec{
				Name:     name,
				Type:     "rate-limit",
				Target:   serviceName,
				Provider: "traefik",
				Rules:    []string{key + "=" + value},
				Source:   platform,
				Raw:      labels,
			}
		}
		if strings.Contains(value, "nginx-auth") || strings.Contains(value, "authentik") || strings.Contains(value, "tinyauth") {
			name := serviceName + "-auth"
			policies[name] = &PolicySpec{
				Name:     name,
				Type:     "auth",
				Target:   serviceName,
				Provider: "traefik",
				Rules:    splitCSV(value),
				Source:   platform,
				Raw:      labels,
			}
		}
	}

	var result []*PolicySpec
	for _, policy := range policies {
		result = append(result, policy)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func routeFromKubernetesExtension(app *Application, raw interface{}) *RouteSpec {
	data, ok := asMap(raw)
	if !ok {
		return nil
	}
	kind := toString(data["kind"])
	metadata, _ := asMap(data["metadata"])
	spec, _ := asMap(data["spec"])
	name := toString(metadata["name"])
	if name == "" {
		return nil
	}
	route := &RouteSpec{
		Name:        name,
		Protocol:    "http",
		Source:      app.Platform,
		Raw:         raw,
		Metadata:    toStringMapLoose(metadata["labels"]),
		Annotations: toStringMapLoose(metadata["annotations"]),
	}
	if route.Metadata == nil {
		route.Metadata = map[string]string{}
	}
	if kind != "" {
		route.Metadata["kubernetes.kind"] = kind
	}
	annotations := route.Annotations
	if portableRoute := portableRouteFromKubernetesAnnotation(annotations[kubernetesPortableRouteAnnotation]); portableRoute != nil {
		mergePortableRouteIntoKubernetesRoute(route, portableRoute)
	}
	if tlsList, ok := spec["tls"].([]interface{}); ok && len(tlsList) > 0 {
		route.TLS = true
	}
	if kind == "HTTPRoute" {
		if hostnames, ok := spec["hostnames"].([]interface{}); ok {
			for _, host := range hostnames {
				if value := toString(host); value != "" {
					route.Hosts = appendUnique(route.Hosts, value)
				}
			}
		}
		if rules, ok := spec["rules"].([]interface{}); ok {
			for _, ruleValue := range rules {
				rule, ok := asMap(ruleValue)
				if !ok {
					continue
				}
				matches, _ := rule["matches"].([]interface{})
				for _, matchValue := range matches {
					match, ok := asMap(matchValue)
					if !ok {
						continue
					}
					if pathSpec, ok := asMap(match["path"]); ok {
						if path := toString(pathSpec["value"]); path != "" {
							route.Paths = appendUnique(route.Paths, path)
						}
					}
				}
				backendRefs, _ := rule["backendRefs"].([]interface{})
				for _, backendValue := range backendRefs {
					backend, ok := asMap(backendValue)
					if !ok {
						continue
					}
					if backendService := toString(backend["name"]); backendService != "" && route.Service == "" {
						route.Service = backendService
						route.Metadata["kubernetes.service"] = backendService
					}
					if backendPort := toString(backend["port"]); backendPort != "" && route.Port == "" {
						route.Port = backendPort
						route.Metadata["kubernetes.servicePort"] = backendPort
					}
				}
			}
		}
	} else if rules, ok := spec["rules"].([]interface{}); ok {
		for _, ruleValue := range rules {
			rule, ok := asMap(ruleValue)
			if !ok {
				continue
			}
			if host := toString(rule["host"]); host != "" {
				route.Hosts = appendUnique(route.Hosts, host)
			}
			httpSpec, _ := asMap(rule["http"])
			paths, _ := httpSpec["paths"].([]interface{})
			for _, pathValue := range paths {
				pathSpec, ok := asMap(pathValue)
				if !ok {
					continue
				}
				if path := toString(pathSpec["path"]); path != "" {
					route.Paths = appendUnique(route.Paths, path)
				}
				backend, _ := asMap(pathSpec["backend"])
				service, _ := asMap(backend["service"])
				backendService := toString(service["name"])
				if route.Service == "" {
					route.Service = backendService
					if backendService != "" {
						route.Metadata["kubernetes.service"] = backendService
					}
				}
				port, _ := asMap(service["port"])
				if route.Port == "" {
					backendPort := toString(port["number"])
					if backendPort == "" {
						backendPort = toString(port["name"])
					}
					if backendPort != "" {
						route.Port = backendPort
						route.Metadata["kubernetes.servicePort"] = backendPort
					}
				}
			}
		}
	}
	resolveKubernetesRouteBackend(app, route)
	return route
}

func portableRouteFromKubernetesAnnotation(raw string) *RouteSpec {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var mapped map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &mapped); err != nil {
		return nil
	}
	route := routeSpecFromExtensionMap(mapped)
	if route == nil {
		return nil
	}
	route.Raw = nil
	return route
}

func mergePortableRouteIntoKubernetesRoute(route, portable *RouteSpec) {
	if route == nil || portable == nil {
		return
	}
	if route.Protocol == "" && portable.Protocol != "" {
		route.Protocol = portable.Protocol
	}
	if len(route.Entrypoints) == 0 {
		route.Entrypoints = append([]string{}, portable.Entrypoints...)
	}
	if len(route.Middlewares) == 0 {
		route.Middlewares = append([]string{}, portable.Middlewares...)
	}
	if len(route.Hosts) == 0 {
		route.Hosts = append([]string{}, portable.Hosts...)
	}
	if len(route.Paths) == 0 {
		route.Paths = append([]string{}, portable.Paths...)
	}
	if route.Service == "" {
		route.Service = portable.Service
	}
	if route.Port == "" {
		route.Port = portable.Port
	}
	if portable.TLS {
		route.TLS = true
	}
	for key, value := range portable.Metadata {
		if route.Metadata == nil {
			route.Metadata = map[string]string{}
		}
		if _, exists := route.Metadata[key]; !exists {
			route.Metadata[key] = value
		}
	}
	for key, value := range portable.Annotations {
		if route.Annotations == nil {
			route.Annotations = map[string]string{}
		}
		if _, exists := route.Annotations[key]; !exists {
			route.Annotations[key] = value
		}
	}
	if len(portable.Extensions) > 0 {
		if route.Extensions == nil {
			route.Extensions = map[string]interface{}{}
		}
		for key, value := range portable.Extensions {
			if _, exists := route.Extensions[key]; !exists {
				route.Extensions[key] = deepCopyValue(value)
			}
		}
	}
	if route.Source == "" {
		route.Source = portable.Source
	}
}

func resolveKubernetesRouteBackend(app *Application, route *RouteSpec) {
	if app == nil || route == nil {
		return
	}
	backendService := route.Metadata["kubernetes.service"]
	if backendService == "" {
		backendService = route.Service
	}
	backendPort := route.Metadata["kubernetes.servicePort"]
	if backendPort == "" {
		backendPort = route.Port
	}
	if targetService := lookupKubernetesServiceTarget(app, backendService); targetService != "" {
		route.Service = targetService
	}
	if targetPort := lookupKubernetesServicePortTarget(app, backendService, backendPort); targetPort != "" {
		route.Port = targetPort
	}
}

func lookupKubernetesServiceTarget(app *Application, serviceName string) string {
	if app == nil || serviceName == "" {
		return ""
	}
	value, ok := applicationExtensionValueForKey(app, "kubernetes.serviceTargets")
	if !ok {
		return ""
	}
	targets := toStringMapLoose(value)
	return targets[serviceName]
}

func lookupKubernetesServicePortTarget(app *Application, serviceName, servicePort string) string {
	if app == nil || serviceName == "" || servicePort == "" {
		return ""
	}
	value, ok := applicationExtensionValueForKey(app, "kubernetes.servicePortTargets")
	if !ok {
		return ""
	}
	targets := toStringMapLoose(value)
	return targets[serviceName+":"+servicePort]
}

func policyFromKubernetesExtension(platform Platform, raw interface{}) *PolicySpec {
	data, ok := asMap(raw)
	if !ok {
		return nil
	}
	metadata, _ := asMap(data["metadata"])
	spec, _ := asMap(data["spec"])
	name := toString(metadata["name"])
	if name == "" {
		return nil
	}
	policy := &PolicySpec{
		Name:        name,
		Type:        strings.ToLower(toString(data["kind"])),
		Provider:    "kubernetes",
		Source:      platform,
		Raw:         raw,
		Metadata:    toStringMapLoose(metadata["labels"]),
		Annotations: toStringMapLoose(metadata["annotations"]),
	}
	if portablePolicy := portablePolicyFromKubernetesAnnotation(policy.Annotations[kubernetesPortablePolicyAnnotation]); portablePolicy != nil {
		mergePortablePolicyIntoKubernetesPolicy(policy, portablePolicy)
	}
	if selector, ok := asMap(spec["podSelector"]); ok {
		if matchLabels, ok := asMap(selector["matchLabels"]); ok {
			var rules []string
			for key, value := range matchLabels {
				rules = append(rules, key+"="+toString(value))
			}
			sort.Strings(rules)
			policy.Rules = rules
		}
	}
	return policy
}

func portablePolicyFromKubernetesAnnotation(raw string) *PolicySpec {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var mapped map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &mapped); err != nil {
		return nil
	}
	policy := policySpecFromExtensionMap(mapped)
	if policy == nil {
		return nil
	}
	policy.Raw = nil
	return policy
}

func mergePortablePolicyIntoKubernetesPolicy(policy, portable *PolicySpec) {
	if policy == nil || portable == nil {
		return
	}
	if policy.Type == "" && portable.Type != "" {
		policy.Type = portable.Type
	}
	if policy.Target == "" {
		policy.Target = portable.Target
	}
	if policy.Provider == "" {
		policy.Provider = portable.Provider
	}
	if len(policy.Rules) == 0 {
		policy.Rules = append([]string{}, portable.Rules...)
	}
	for key, value := range portable.Metadata {
		if policy.Metadata == nil {
			policy.Metadata = map[string]string{}
		}
		if _, exists := policy.Metadata[key]; !exists {
			policy.Metadata[key] = value
		}
	}
	for key, value := range portable.Annotations {
		if policy.Annotations == nil {
			policy.Annotations = map[string]string{}
		}
		if _, exists := policy.Annotations[key]; !exists {
			policy.Annotations[key] = value
		}
	}
	if len(portable.Extensions) > 0 {
		if policy.Extensions == nil {
			policy.Extensions = map[string]interface{}{}
		}
		for key, value := range portable.Extensions {
			if _, exists := policy.Extensions[key]; !exists {
				policy.Extensions[key] = deepCopyValue(value)
			}
		}
	}
	if len(portable.Extension) > 0 {
		if policy.Extension == nil {
			policy.Extension = map[string]string{}
		}
		for key, value := range portable.Extension {
			if _, exists := policy.Extension[key]; !exists {
				policy.Extension[key] = value
			}
		}
	}
	if policy.Source == "" {
		policy.Source = portable.Source
	}
}

func routeSpecFromExtensionMap(mapped map[string]interface{}) *RouteSpec {
	if len(mapped) == 0 {
		return nil
	}
	typed, err := jsonMarshalUnmarshal[RouteSpec](mapped)
	if err != nil || typed == nil {
		return nil
	}
	if typed.Extensions == nil {
		typed.Extensions = map[string]interface{}{}
	}
	for key, value := range mapped {
		switch key {
		case "name", "service", "protocol", "hosts", "paths", "port", "entrypoints", "middlewares", "tls", "source", "raw", "metadata", "annotations", "extensions":
			continue
		}
		if _, exists := typed.Extensions[key]; !exists {
			typed.Extensions[key] = deepCopyValue(value)
		}
	}
	return typed
}

func policySpecFromExtensionMap(mapped map[string]interface{}) *PolicySpec {
	if len(mapped) == 0 {
		return nil
	}
	typed, err := jsonMarshalUnmarshal[PolicySpec](mapped)
	if err != nil || typed == nil {
		return nil
	}
	if typed.Extensions == nil {
		typed.Extensions = map[string]interface{}{}
	}
	for key, value := range mapped {
		switch key {
		case "name", "type", "target", "provider", "rules", "source", "raw", "metadata", "annotations", "extension", "extensions":
			continue
		}
		if _, exists := typed.Extensions[key]; !exists {
			typed.Extensions[key] = deepCopyValue(value)
		}
	}
	return typed
}

func labelTruthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "on":
		return true
	default:
		return false
	}
}

func hostsFromTraefikRule(rule string) []string {
	var hosts []string
	for _, match := range traefikHostPattern.FindAllStringSubmatch(rule, -1) {
		if len(match) > 1 {
			hosts = appendUnique(hosts, match[1])
		}
	}
	return hosts
}

func pathsFromTraefikRule(rule string) []string {
	var paths []string
	for _, token := range []string{"PathPrefix(`", "Path(`"} {
		remaining := rule
		for {
			index := strings.Index(remaining, token)
			if index == -1 {
				break
			}
			remaining = remaining[index+len(token):]
			end := strings.Index(remaining, "`)")
			if end == -1 {
				break
			}
			paths = appendUnique(paths, remaining[:end])
			remaining = remaining[end+2:]
		}
	}
	return paths
}

func splitCSV(value string) []string {
	if value == "" {
		return nil
	}
	var result []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func applyCanonicalRoutesToComposeServices(app *Application, servicesData map[string]interface{}) {
	canonical := canonicalForApplication(app)
	if canonical == nil {
		return
	}
	for _, route := range canonical.Routes {
		if route.Service == "" {
			continue
		}
		serviceData, ok := servicesData[route.Service].(map[string]interface{})
		if !ok {
			continue
		}
		labels := map[string]string{}
		if existing, ok := serviceData["labels"]; ok {
			if parsed, err := toStringMap(existing); err == nil {
				labels = parsed
			}
		}
		for key, value := range traefikLabelsFromRoute(route) {
			labels[key] = value
		}
		for _, policy := range canonical.Policies {
			if policy.Target == route.Service && (policy.Type == "auth" || policy.Type == "middleware-chain" || policy.Type == "rate-limit") {
				middlewareKey := fmt.Sprintf("traefik.http.routers.%s.middlewares", route.Name)
				labels[middlewareKey] = strings.Join(appendUniqueList(splitCSV(labels[middlewareKey]), policy.Rules), ",")
			}
		}
		if len(labels) > 0 {
			serviceData["labels"] = labels
		}
	}
}

func traefikLabelsFromRoute(route *RouteSpec) map[string]string {
	labels := map[string]string{"traefik.enable": "true"}
	base := fmt.Sprintf("traefik.http.routers.%s.", route.Name)
	if rule := traefikRuleFromRoute(route); rule != "" {
		labels[base+"rule"] = rule
	}
	if len(route.Entrypoints) > 0 {
		labels[base+"entrypoints"] = strings.Join(route.Entrypoints, ",")
	}
	if len(route.Middlewares) > 0 {
		labels[base+"middlewares"] = strings.Join(route.Middlewares, ",")
	}
	if route.TLS {
		labels[base+"tls"] = "true"
	}
	if route.Port != "" {
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", route.Name)] = route.Port
		labels[base+"service"] = route.Name
	}
	return labels
}

func traefikRuleFromRoute(route *RouteSpec) string {
	var parts []string
	for _, host := range route.Hosts {
		parts = append(parts, fmt.Sprintf("Host(`%s`)", host))
	}
	for _, path := range route.Paths {
		parts = append(parts, fmt.Sprintf("PathPrefix(`%s`)", path))
	}
	return strings.Join(parts, " && ")
}

func appendUniqueList(values []string, additions []string) []string {
	for _, addition := range additions {
		values = appendUnique(values, addition)
	}
	return values
}

func serializeKubernetesRoute(route *RouteSpec, app *Application) (string, error) {
	if route == nil || route.Name == "" {
		return "", nil
	}
	kind := route.Metadata["kubernetes.kind"]
	if kind != "" && kind != "Ingress" {
		return serializeKubernetesGatewayRoute(route, app, kind)
	}
	if route.Service == "" {
		return "", nil
	}
	serviceName := kubernetesRouteServiceName(route)
	port := kubernetesRouteServicePort(route)
	if port == "" {
		port = firstServicePort(app, serviceName)
	}
	if port == "" {
		port = "80"
	}
	paths := route.Paths
	if len(paths) == 0 {
		paths = []string{"/"}
	}
	hosts := route.Hosts
	if len(hosts) == 0 {
		hosts = []string{""}
	}

	var rules []map[string]interface{}
	for _, host := range hosts {
		var pathSpecs []map[string]interface{}
		for _, path := range paths {
			pathSpecs = append(pathSpecs, map[string]interface{}{
				"path":     path,
				"pathType": "Prefix",
				"backend": map[string]interface{}{
					"service": map[string]interface{}{
						"name": serviceName,
						"port": kubernetesServicePort(port),
					},
				},
			})
		}
		rule := map[string]interface{}{
			"http": map[string]interface{}{
				"paths": pathSpecs,
			},
		}
		if host != "" {
			rule["host"] = host
		}
		rules = append(rules, rule)
	}

	ingress := map[string]interface{}{
		"apiVersion": "networking.k8s.io/v1",
		"kind":       "Ingress",
		"metadata": map[string]interface{}{
			"name":      route.Name,
			"namespace": app.Namespace,
		},
		"spec": map[string]interface{}{
			"rules": rules,
		},
	}
	if len(route.Annotations) > 0 {
		ingress["metadata"].(map[string]interface{})["annotations"] = copyStringMap(route.Annotations)
	}
	if portable := clonePortableRouteSpec(route); portable != nil {
		raw, err := json.Marshal(portable)
		if err != nil {
			return "", err
		}
		metadata := ingress["metadata"].(map[string]interface{})
		annotations, _ := metadata["annotations"].(map[string]string)
		if annotations == nil {
			annotations = map[string]string{}
		}
		annotations[kubernetesPortableRouteAnnotation] = string(raw)
		metadata["annotations"] = annotations
	}
	if route.TLS && len(route.Hosts) > 0 {
		ingress["spec"].(map[string]interface{})["tls"] = []map[string]interface{}{
			{"hosts": route.Hosts},
		}
	}
	data, err := marshalYAMLString(ingress)
	return data, err
}

func serializeKubernetesGatewayRoute(route *RouteSpec, app *Application, kind string) (string, error) {
	if route == nil || route.Name == "" {
		return "", nil
	}
	if kind == "" {
		kind = "HTTPRoute"
	}
	routeResource := map[string]interface{}{
		"apiVersion": "gateway.networking.k8s.io/v1",
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":      route.Name,
			"namespace": app.Namespace,
		},
		"spec": map[string]interface{}{},
	}
	metadata := routeResource["metadata"].(map[string]interface{})
	if len(route.Annotations) > 0 {
		metadata["annotations"] = copyStringMap(route.Annotations)
	}
	if portable := clonePortableRouteSpec(route); portable != nil {
		raw, err := json.Marshal(portable)
		if err != nil {
			return "", err
		}
		annotations, _ := metadata["annotations"].(map[string]string)
		if annotations == nil {
			annotations = map[string]string{}
		}
		annotations[kubernetesPortableRouteAnnotation] = string(raw)
		metadata["annotations"] = annotations
	}
	spec := routeResource["spec"].(map[string]interface{})
	if len(route.Hosts) > 0 {
		hostnames := make([]string, 0, len(route.Hosts))
		for _, host := range route.Hosts {
			if strings.TrimSpace(host) != "" {
				hostnames = append(hostnames, host)
			}
		}
		if len(hostnames) > 0 {
			spec["hostnames"] = hostnames
		}
	}
	if len(route.Paths) > 0 || kubernetesRouteServiceName(route) != "" || kubernetesRouteServicePort(route) != "" {
		matches := []map[string]interface{}{}
		for _, path := range route.Paths {
			if strings.TrimSpace(path) == "" {
				continue
			}
			matches = append(matches, map[string]interface{}{
				"path": map[string]interface{}{
					"type":  "PathPrefix",
					"value": path,
				},
			})
		}
		if len(matches) == 0 {
			matches = append(matches, map[string]interface{}{})
		}
		backend := map[string]interface{}{}
		if serviceName := kubernetesRouteServiceName(route); serviceName != "" {
			backend["name"] = serviceName
		}
		if servicePort := kubernetesRouteServicePort(route); servicePort != "" {
			backend["port"] = kubernetesIntOrString(servicePort)
		}
		rule := map[string]interface{}{
			"matches":     matches,
			"backendRefs": []map[string]interface{}{backend},
		}
		spec["rules"] = []map[string]interface{}{rule}
	}
	data, err := marshalYAMLString(routeResource)
	return data, err
}

func kubernetesRouteServiceName(route *RouteSpec) string {
	if route == nil {
		return ""
	}
	if serviceName := route.Metadata["kubernetes.service"]; serviceName != "" {
		return serviceName
	}
	return route.Service
}

func kubernetesRouteServicePort(route *RouteSpec) string {
	if route == nil {
		return ""
	}
	if servicePort := route.Metadata["kubernetes.servicePort"]; servicePort != "" {
		return servicePort
	}
	return route.Port
}

func serializeKubernetesPolicy(policy *PolicySpec, namespace string) (string, error) {
	if policy == nil || policy.Name == "" {
		return "", nil
	}
	if policy.Type != "networkpolicy" {
		return serializeKubernetesPortablePolicy(policy, namespace)
	}
	matchLabels := map[string]string{}
	for _, rule := range policy.Rules {
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) == 2 {
			matchLabels[parts[0]] = parts[1]
		}
	}
	spec := map[string]interface{}{
		"podSelector": map[string]interface{}{},
	}
	if len(matchLabels) > 0 {
		spec["podSelector"] = map[string]interface{}{"matchLabels": matchLabels}
	}
	networkPolicy := map[string]interface{}{
		"apiVersion": "networking.k8s.io/v1",
		"kind":       "NetworkPolicy",
		"metadata": map[string]interface{}{
			"name":      policy.Name,
			"namespace": namespace,
		},
		"spec": spec,
	}
	if len(policy.Annotations) > 0 {
		networkPolicy["metadata"].(map[string]interface{})["annotations"] = copyStringMap(policy.Annotations)
	}
	if portable := clonePortablePolicySpec(policy); portable != nil {
		raw, err := json.Marshal(portable)
		if err != nil {
			return "", err
		}
		metadata := networkPolicy["metadata"].(map[string]interface{})
		annotations, _ := metadata["annotations"].(map[string]string)
		if annotations == nil {
			annotations = map[string]string{}
		}
		annotations[kubernetesPortablePolicyAnnotation] = string(raw)
		metadata["annotations"] = annotations
	}
	data, err := marshalYAMLString(networkPolicy)
	return data, err
}

func serializeKubernetesPortablePolicy(policy *PolicySpec, namespace string) (string, error) {
	portable := clonePortablePolicySpec(policy)
	if portable == nil {
		return "", nil
	}
	raw, err := json.Marshal(portable)
	if err != nil {
		return "", err
	}
	configMap := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      "bolabaden-policy-" + policy.Name,
			"namespace": namespace,
			"labels": map[string]string{
				"bolabaden.dev/portable-policy": policy.Name,
			},
		},
		"data": map[string]string{
			"policy.json": string(raw),
		},
	}
	if len(policy.Annotations) > 0 {
		configMap["metadata"].(map[string]interface{})["annotations"] = copyStringMap(policy.Annotations)
	}
	return marshalYAMLString(configMap)
}

func marshalYAMLString(value interface{}) (string, error) {
	data, err := yaml.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func firstServicePort(app *Application, serviceName string) string {
	if app == nil {
		return ""
	}
	service := app.Services[serviceName]
	if service == nil || len(service.Ports) == 0 {
		return ""
	}
	if service.Ports[0].HostPort != "" {
		return service.Ports[0].HostPort
	}
	return service.Ports[0].ContainerPort
}

func kubernetesServicePort(port string) map[string]interface{} {
	if number := parseInt(port); number > 0 {
		return map[string]interface{}{"number": number}
	}
	return map[string]interface{}{"name": port}
}
