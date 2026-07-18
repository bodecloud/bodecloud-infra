package paas

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/resource"
)

const composeRawYAMLExtension = "compose.raw.yaml"
const composeAppRoutesExtension = "x-bolabaden-canonical-routes"
const composeAppPoliciesExtension = "x-bolabaden-canonical-policies"
const composeCanonicalRawResourcesExtension = "x-bolabaden-canonical-raw-resources"
const composeKubernetesRawResourcesExtension = "x-bolabaden-kubernetes-raw-resources"
const composeFailoverExtension = "x-bolabaden-failover"
const meshExtensionKey = "x-bolabaden-mesh"
const legacyMeshExtensionKey = "bolabaden.mesh"
const composeCompatExtension = "x-bolabaden-compose-compat"
const composeBlkioConfigExtensionsExtension = "x-bolabaden-blkio-extensions"
const composeKubernetesNamespaceExtension = "x-kubernetes-namespace"
const composeSwarmJobExtension = "x-swarm-job"
const composeMemoryLimitExtension = "x-bolabaden-compose-memory-limit"
const composeMemoryReservationExtension = "x-bolabaden-compose-memory-reservation"
const composeMemorySwapExtension = "x-bolabaden-compose-memory-swap"
const composeShmSizeExtension = "x-bolabaden-compose-shm-size"

// ParseDockerCompose parses a Docker Compose YAML file or content
func ParseDockerCompose(content string) (*Application, error) {
	if app, err := parseDockerComposeSpec(content); err == nil {
		return app, nil
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse docker-compose YAML: %w", err)
	}

	app := &Application{
		Platform: PlatformDockerCompose,
		Services: make(map[string]*Service),
		Networks: make(map[string]*Network),
		Volumes:  make(map[string]*Volume),
		Configs:  make(map[string]*Config),
		Secrets:  make(map[string]*Secret),
		Models:   make(map[string]*ComposeModel),
	}

	// Parse version
	if version := toString(data["version"]); version != "" {
		app.Version = version
	}
	if name := toString(data["name"]); name != "" {
		app.Name = name
	}
	if strings.EqualFold(toString(data["x-platform"]), string(PlatformDockerSwarm)) {
		app.Platform = PlatformDockerSwarm
	}

	// Parse includes
	if includes, includeEntries := parseComposeIncludeEntries(data["include"]); len(includeEntries) > 0 {
		app.Includes = includes
		app.IncludeEntries = includeEntries
	}

	// Parse services
	if servicesData, ok := asMap(data["services"]); ok {
		for name, serviceData := range servicesData {
			serviceMap, ok := asMap(serviceData)
			if !ok {
				return nil, fmt.Errorf("service %s must be a mapping", name)
			}
			service, err := parseService(name, serviceMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse service %s: %w", name, err)
			}
			service.Platform = app.Platform
			app.Services[name] = service
		}
	}

	// Parse networks
	if networksData, ok := asMap(data["networks"]); ok {
		for name, networkData := range networksData {
			networkMap, _ := asMap(networkData)
			network, err := parseNetwork(name, networkMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse network %s: %w", name, err)
			}
			app.Networks[name] = network
		}
	}

	// Parse volumes
	if volumesData, ok := asMap(data["volumes"]); ok {
		for name, volumeData := range volumesData {
			volumeMap, _ := asMap(volumeData)
			volume, err := parseVolume(name, volumeMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse volume %s: %w", name, err)
			}
			app.Volumes[name] = volume
		}
	}

	// Parse configs
	if configsData, ok := asMap(data["configs"]); ok {
		for name, configData := range configsData {
			configMap, _ := asMap(configData)
			config, err := parseConfig(name, configMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse config %s: %w", name, err)
			}
			app.Configs[name] = config
		}
	}

	// Parse secrets
	if secretsData, ok := asMap(data["secrets"]); ok {
		for name, secretData := range secretsData {
			secretMap, _ := asMap(secretData)
			secret, err := parseSecret(name, secretMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse secret %s: %w", name, err)
			}
			app.Secrets[name] = secret
		}
	}

	// Parse top-level models
	if modelsData, ok := asMap(data["models"]); ok {
		for name, modelData := range modelsData {
			modelMap, _ := asMap(modelData)
			if model := parseComposeTopLevelModel(name, modelMap); model != nil {
				app.Models[name] = model
			}
		}
	}

	// Store extensions
	app.Extensions = make(map[string]interface{})
	for key, value := range data {
		if !isStandardKey(key) {
			app.Extensions[key] = value
			if alias := composeApplicationCanonicalKey(key); alias != "" {
				app.Extensions[alias] = value
			}
		}
	}
	rehydrateComposeApplicationExtensions(app)
	syncPortableApplicationState(app)
	app.Extensions[composeRawYAMLExtension] = content
	app.AttachCanonical()
	app.Canonical.AddResource(ResourceKindRaw, app.Platform, "compose-yaml", "ComposeYAML", content)

	return app, nil
}

// ParseDockerComposeFile parses a Docker Compose file from disk using its real
// filename so compose-go can resolve include/extends references correctly.
func ParseDockerComposeFile(filename string) (*Application, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read docker-compose file %s: %w", filename, err)
	}
	if app, err := parseDockerComposeFile(filename, string(content)); err == nil {
		return app, nil
	}
	return ParseDockerCompose(string(content))
}

func parseComposeIncludeEntries(value interface{}) ([]string, []interface{}) {
	items, ok := value.([]interface{})
	if !ok || len(items) == 0 {
		return nil, nil
	}
	entries := make([]interface{}, 0, len(items))
	var includes []string
	for _, item := range items {
		switch typed := item.(type) {
		case string:
			entries = append(entries, typed)
			appendUniqueString(&includes, typed)
		default:
			copied := deepCopyValue(item)
			entries = append(entries, copied)
			includes = mergeUniqueStrings(includes, composeIncludePaths([]interface{}{copied}))
		}
	}
	return includes, entries
}

func composeIncludePaths(entries []interface{}) []string {
	var paths []string
	for _, entry := range entries {
		switch typed := entry.(type) {
		case string:
			appendUniqueString(&paths, typed)
		default:
			mapped, ok := asMap(typed)
			if !ok {
				continue
			}
			values, err := toStringSlice(mapped["path"])
			if err != nil {
				continue
			}
			paths = mergeUniqueStrings(paths, values)
		}
	}
	return paths
}

func parseService(name string, data map[string]interface{}) (*Service, error) {
	service := &Service{
		Name:        name,
		Platform:    PlatformDockerCompose,
		Environment: make(map[string]string),
		Labels:      make(map[string]string),
		Extensions:  make(map[string]interface{}),
	}

	for key, value := range data {
		if strings.HasPrefix(key, "x-kubernetes-") {
			service.Extensions[key] = value
		}
		if alias := kubernetesComposeExtensionAlias(key); alias != "" && alias != key {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			if _, exists := service.Extensions[alias]; !exists {
				service.Extensions[alias] = value
			}
		}
		switch key {
		case "image":
			service.Image = toString(value)
		case "container_name":
			service.ContainerName = toString(value)
		case "hostname":
			service.Hostname = toString(value)
		case "ports":
			ports, err := parsePorts(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ports: %w", err)
			}
			service.Ports = ports
		case "expose":
			expose, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse expose: %w", err)
			}
			service.Expose = expose
		case "networks":
			networks, attachments, err := parseServiceNetworks(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse networks: %w", err)
			}
			service.Networks = networks
			service.NetworkAttachments = attachments
		case "dns":
			dns, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse dns: %w", err)
			}
			service.DNS = dns
		case "dns_search":
			search, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse dns_search: %w", err)
			}
			service.DNSSearch = search
		case "dns_opt":
			options, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse dns_opt: %w", err)
			}
			service.DNSOptions = options
		case "extra_hosts":
			hosts, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse extra_hosts: %w", err)
			}
			service.ExtraHosts = hosts
		case "environment":
			env, err := parseEnvironment(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse environment: %w", err)
			}
			service.Environment = env
		case "env_file":
			envFileRefs, err := parseEnvFileRefs(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse env_file: %w", err)
			}
			service.EnvFileRefs = envFileRefs
			service.EnvFile = envFilePaths(envFileRefs)
		case "x-env-sources":
			envSources, err := parseEnvSources(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse x-env-sources: %w", err)
			}
			service.EnvSources = envSources
		case "x-env-from":
			envFrom, err := parseEnvFromSources(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse x-env-from: %w", err)
			}
			service.EnvFrom = envFrom
		case "x-kubernetes-imagePullSecrets":
			secrets, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse x-kubernetes-imagePullSecrets: %w", err)
			}
			service.ImagePullSecrets = secrets
		case "x-kubernetes-image-pull-secrets":
			secrets, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse x-kubernetes-image-pull-secrets: %w", err)
			}
			service.ImagePullSecrets = secrets
		case "x-kubernetes-imagePullPolicy":
			service.ImagePullPolicy = toString(value)
		case "x-kubernetes-image-pull-policy":
			service.ImagePullPolicy = toString(value)
		case "x-kubernetes-dnsPolicy":
			service.DNSPolicy = toString(value)
		case "x-kubernetes-dns-policy":
			service.DNSPolicy = toString(value)
		case "x-kubernetes-schedulerName":
			service.SchedulerName = toString(value)
		case "x-kubernetes-scheduler-name":
			service.SchedulerName = toString(value)
		case "x-kubernetes-hostname":
			if service.Hostname == "" {
				service.Hostname = toString(value)
			}
		case "x-kubernetes-terminationMessagePath":
			service.TerminationMessagePath = toString(value)
		case "x-kubernetes-termination-message-path":
			service.TerminationMessagePath = toString(value)
		case "x-kubernetes-terminationMessagePolicy":
			service.TerminationMessagePolicy = toString(value)
		case "x-kubernetes-termination-message-policy":
			service.TerminationMessagePolicy = toString(value)
		case "x-kubernetes-hostAliases":
			if aliases := kubernetesHostAliasesFromExtension(value); len(aliases) > 0 {
				service.HostAliases = aliases
				for _, alias := range aliases {
					for _, hostname := range alias.Hostnames {
						appendUniqueString(&service.ExtraHosts, hostname+"="+alias.IP)
					}
				}
			}
		case "x-kubernetes-host-aliases":
			if aliases := kubernetesHostAliasesFromExtension(value); len(aliases) > 0 {
				service.HostAliases = aliases
				for _, alias := range aliases {
					for _, hostname := range alias.Hostnames {
						appendUniqueString(&service.ExtraHosts, hostname+"="+alias.IP)
					}
				}
			}
		case "x-kubernetes-fsGroup":
			if fsGroup := int64(toInt(value)); fsGroup > 0 {
				service.FSGroup = &fsGroup
			}
		case "x-kubernetes-fs-group":
			if fsGroup := int64(toInt(value)); fsGroup > 0 {
				service.FSGroup = &fsGroup
			}
		case "x-kubernetes-seLinuxOptions":
			if options, ok := asMap(value); ok {
				service.SELinuxOptions = parseKubernetesSELinuxOptions(options)
			}
		case "x-kubernetes-se-linux-options":
			if options, ok := asMap(value); ok {
				service.SELinuxOptions = parseKubernetesSELinuxOptions(options)
			}
		case "x-kubernetes-windowsOptions":
			if options, ok := asMap(value); ok {
				service.WindowsOptions = parseKubernetesWindowsSecurityContextOptions(options)
			}
		case "x-kubernetes-windows-options":
			if options, ok := asMap(value); ok {
				service.WindowsOptions = parseKubernetesWindowsSecurityContextOptions(options)
			}
		case "x-kubernetes-fsGroupChangePolicy":
			service.FSGroupChangePolicy = toString(value)
		case "x-kubernetes-fs-group-change-policy":
			service.FSGroupChangePolicy = toString(value)
		case "x-kubernetes-runAsNonRoot":
			if str := toString(value); str != "" {
				runAsNonRoot := strings.EqualFold(str, "true")
				service.RunAsNonRoot = &runAsNonRoot
			}
		case "x-kubernetes-run-as-non-root":
			if str := toString(value); str != "" {
				runAsNonRoot := strings.EqualFold(str, "true")
				service.RunAsNonRoot = &runAsNonRoot
			}
		case "x-kubernetes-supplementalGroups":
			if groups, err := toInt64Slice(value); err == nil && len(groups) > 0 {
				service.SupplementalGroups = groups
			}
		case "x-kubernetes-supplementalGroupsPolicy":
			service.SupplementalGroupsPolicy = toString(value)
		case "x-kubernetes-supplemental-groups-policy":
			service.SupplementalGroupsPolicy = toString(value)
		case "x-kubernetes-activeDeadlineSeconds":
			if seconds := int64(toInt(value)); seconds > 0 {
				service.ActiveDeadlineSeconds = &seconds
			}
		case "x-kubernetes-active-deadline-seconds":
			if seconds := int64(toInt(value)); seconds > 0 {
				service.ActiveDeadlineSeconds = &seconds
			}
		case "x-kubernetes-tolerations":
			if tolerations := parseKubernetesTolerationsExtension(value); len(tolerations) > 0 {
				service.Tolerations = tolerations
			}
		case "x-kubernetes-restartPolicy":
			if restartPolicy := toString(value); restartPolicy != "" {
				service.PodRestartPolicy = restartPolicy
			}
		case "x-kubernetes-restart-policy":
			if restartPolicy := toString(value); restartPolicy != "" {
				service.PodRestartPolicy = restartPolicy
			}
		case "x-kubernetes-startup-probe":
			if probe, ok := asMap(value); ok {
				service.StartupProbe = parseKubernetesProbe(probe)
			}
		case "x-kubernetes-seccomp-profile":
			if profile, ok := asMap(value); ok {
				service.SeccompProfile = parseKubernetesSeccompProfile(profile)
			}
		case "x-kubernetes-seccompProfile":
			if profile, ok := asMap(value); ok {
				service.SeccompProfile = parseKubernetesSeccompProfile(profile)
			}
		case "x-kubernetes-affinity":
			if affinity, ok := asMap(value); ok {
				service.Affinity = cloneMap(affinity)
			}
		case "x-kubernetes-readiness-gates":
			if gates, ok := kubernetesReadinessGatesFromExtension(value); ok {
				service.ReadinessGates = gates
			}
		case "x-kubernetes-readinessGates":
			if gates, ok := kubernetesReadinessGatesFromExtension(value); ok {
				service.ReadinessGates = gates
			}
		case "x-kubernetes-allowPrivilegeEscalation":
			service.AllowPrivilegeEscalation = boolPtrFromInterface(value)
		case "x-kubernetes-allow-privilege-escalation":
			service.AllowPrivilegeEscalation = boolPtrFromInterface(value)
		case "x-kubernetes-procMount":
			service.ProcMount = toString(value)
		case "x-kubernetes-proc-mount":
			service.ProcMount = toString(value)
		case "x-kubernetes-init-containers":
			if containers, ok := kubernetesMapSliceFromExtension(value); ok {
				service.InitContainers = containers
			}
		case "x-kubernetes-initContainers":
			if containers, ok := kubernetesMapSliceFromExtension(value); ok {
				service.InitContainers = containers
			}
		case "x-kubernetes-scheduling-gates":
			if gates, ok := kubernetesReadinessGatesFromExtension(value); ok {
				service.SchedulingGates = gates
			}
		case "x-kubernetes-schedulingGates":
			if gates, ok := kubernetesReadinessGatesFromExtension(value); ok {
				service.SchedulingGates = gates
			}
		case "x-kubernetes-hostUsers", "x-kubernetes-host-users":
			service.HostUsers = boolPtrFromInterface(value)
		case "x-kubernetes-group":
			if group := toString(value); group != "" {
				service.Group = group
				service.Extensions["kubernetes.group"] = group
			}
		case "x-kubernetes-pidMode":
			if mode := toString(value); mode != "" {
				service.PIDMode = mode
				if strings.EqualFold(mode, "host") {
					service.HostPID = boolPtr(true)
				}
				service.Extensions["kubernetes.pidMode"] = mode
			}
		case "x-kubernetes-ipcMode":
			if mode := toString(value); mode != "" {
				service.IPCMode = mode
				if strings.EqualFold(mode, "host") {
					service.HostIPC = boolPtr(true)
				}
				service.Extensions["kubernetes.ipcMode"] = mode
			}
		case "x-kubernetes-hostNetwork", "x-kubernetes-host-network":
			service.HostNetwork = toBool(value)
			service.HostNetworkSet = true
		case "x-kubernetes-hostPID", "x-kubernetes-host-pid":
			service.HostPID = boolPtrFromInterface(value)
			if service.HostPID != nil && *service.HostPID {
				service.PIDMode = "host"
			}
		case "x-kubernetes-hostIPC", "x-kubernetes-host-ipc":
			service.HostIPC = boolPtrFromInterface(value)
			if service.HostIPC != nil && *service.HostIPC {
				service.IPCMode = "host"
			}
		case "x-kubernetes-priorityClassName":
			service.PriorityClassName = toString(value)
		case "x-kubernetes-priority-class-name":
			service.PriorityClassName = toString(value)
		case "x-kubernetes-runtimeClassName":
			service.RuntimeClassName = toString(value)
		case "x-kubernetes-runtime-class-name":
			service.RuntimeClassName = toString(value)
		case "x-kubernetes-nodeName":
			service.NodeName = toString(value)
		case "x-kubernetes-node-name":
			service.NodeName = toString(value)
		case "x-kubernetes-subdomain":
			service.Subdomain = toString(value)
		case "x-kubernetes-sub-domain":
			service.Subdomain = toString(value)
		case "x-kubernetes-os":
			service.OSName = toString(value)
		case "x-kubernetes-setHostnameAsFQDN", "x-kubernetes-set-hostname-as-fqdn":
			service.SetHostnameAsFQDN = boolPtrFromInterface(value)
		case "x-kubernetes-shareProcessNamespace", "x-kubernetes-share-process-namespace":
			service.ShareProcessNamespace = boolPtrFromInterface(value)
		case "x-kubernetes-enableServiceLinks", "x-kubernetes-enable-service-links":
			service.EnableServiceLinks = boolPtrFromInterface(value)
		case "x-kubernetes-serviceAccountName", "x-kubernetes-service-account-name":
			service.ServiceAccountName = toString(value)
		case "x-kubernetes-automountServiceAccountToken", "x-kubernetes-automount-service-account-token":
			service.AutomountServiceAccountToken = boolPtrFromInterface(value)
		case "x-kubernetes-topology-spread-constraints":
			if constraints, ok := kubernetesTopologySpreadConstraintsFromExtension(value); ok {
				service.TopologySpreadConstraints = constraints
			}
		case "x-kubernetes-topologySpreadConstraints":
			if constraints, ok := kubernetesTopologySpreadConstraintsFromExtension(value); ok {
				service.TopologySpreadConstraints = constraints
			}
		case "x-nomad-spread":
			if spreads := nomadSpreadSpecsFromAny(value); len(spreads) > 0 {
				service.Spreads = spreads
				service.Extensions["x-nomad-spread"] = deepCopyValue(value)
			}
		case "x-nomad-connect":
			if connect := nomadConnectSpecFromAny(value); connect != nil {
				service.Connect = connect
				service.Extensions["x-nomad-connect"] = deepCopyValue(value)
			}
		case "x-nomad-restart":
			if restart, ok := asMap(value); ok {
				applyNomadRestartExtensions(service, restart)
				service.Extensions["x-nomad-restart"] = deepCopyValue(value)
			}
		case "command":
			command, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse command: %w", err)
			}
			service.Command = command
		case "entrypoint":
			entrypoint, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse entrypoint: %w", err)
			}
			service.Entrypoint = entrypoint
		case "working_dir":
			service.WorkingDir = toString(value)
		case "attach":
			flag := toBool(value)
			ensureComposeCompat(service).Attach = &flag
		case "annotations":
			ensureComposeCompat(service).Annotations = toStringMapLoose(value)
		case "blkio_config":
			if mapped, ok := asMap(value); ok {
				ensureComposeCompat(service).BlkioConfig = copyStringInterfaceMap(mapped)
			}
		case "credential_spec":
			if mapped, ok := asMap(value); ok {
				ensureComposeCompat(service).CredentialSpec = copyStringInterfaceMap(mapped)
			}
		case "provider":
			if mapped, ok := asMap(value); ok {
				ensureComposeCompat(service).Provider = copyStringInterfaceMap(mapped)
			}
		case "extends":
			ensureComposeCompat(service).Extends = composeExtendsRawMap(value)
		case "platform":
			ensureComposeCompat(service).Platform = toString(value)
		case composeCompatExtension:
			parseComposeCompatExtension(service, value)
		case "pull_refresh_after":
			ensureComposeCompat(service).PullRefreshAfter = toString(value)
		case "logging":
			driver, options, extensions, err := parseComposeLogging(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse logging: %w", err)
			}
			if driver != "" {
				service.LogDriver = driver
			}
			if len(options) > 0 {
				service.LogOpt = options
			}
			if len(extensions) > 0 {
				service.LogExtensions = extensions
			}
		case "log_driver":
			service.LogDriver = toString(value)
		case "log_opt":
			options, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse log_opt: %w", err)
			}
			service.LogOpt = options
		case "build":
			build, err := parseBuildConfig(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse build: %w", err)
			}
			service.Build = build
		case "develop":
			develop, err := parseDevelopConfig(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse develop: %w", err)
			}
			service.Develop = develop
		case "pre_start":
			hooks, err := parseServiceHookList(value, "pre_start")
			if err != nil {
				return nil, err
			}
			ensureLifecycle(service).PreStart = hooks
		case "post_start":
			hooks, err := parseServiceHookList(value, "post_start")
			if err != nil {
				return nil, err
			}
			ensureLifecycle(service).PostStart = hooks
		case "pre_stop":
			hooks, err := parseServiceHookList(value, "pre_stop")
			if err != nil {
				return nil, err
			}
			ensureLifecycle(service).PreStop = hooks
		case "devices":
			devices, mappings, err := parseDevicesWithMappings(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse devices: %w", err)
			}
			service.Devices = devices
			service.DeviceMappings = mappings
		case "profiles":
			profiles, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse profiles: %w", err)
			}
			service.Profiles = profiles
		case "volumes":
			volumes, err := parseVolumes(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse volumes: %w", err)
			}
			service.Volumes = volumes
		case "configs":
			configs, err := parseFileRefs(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse configs: %w", err)
			}
			service.Configs = configs
		case "secrets":
			secrets, err := parseFileRefs(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse secrets: %w", err)
			}
			service.Secrets = secrets
		case "depends_on":
			dependencies, err := parseDependencySpecs(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse depends_on: %w", err)
			}
			service.Dependencies = dependencies
			service.DependsOn = dependencyNames(dependencies)
		case "links":
			links, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse links: %w", err)
			}
			service.Links = links
		case "x-bolabaden-links":
			links, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse x-bolabaden-links: %w", err)
			}
			service.Links = links
		case "restart":
			service.Restart = toString(value)
		case "privileged":
			service.Privileged = toBool(value)
			service.PrivilegedSet = true
		case "user":
			service.User = toString(value)
		case "group":
			service.Group = toString(value)
		case "pid":
			service.PIDMode = toString(value)
			if strings.EqualFold(service.PIDMode, "host") {
				service.HostPID = boolPtr(true)
			}
		case "ipc":
			service.IPCMode = toString(value)
			if strings.EqualFold(service.IPCMode, "host") {
				service.HostIPC = boolPtr(true)
			}
		case "pids_limit":
			service.PidsLimit = int64(toInt(value))
			service.pidsLimitSet = true
		case "shm_size":
			if raw := toString(value); raw != "" {
				if _, ok := service.Extensions[composeShmSizeExtension]; !ok {
					service.Extensions[composeShmSizeExtension] = raw
				}
				if quantity, err := resource.ParseQuantity(raw); err == nil {
					service.ShmSize = quantity.Value()
				} else {
					service.ShmSize = int64(toInt(value))
				}
				service.shmSizeSet = true
			}
		case "mac_address":
			ensureComposeCompat(service).MacAddress = toString(value)
		case "domainname":
			ensureComposeCompat(service).DomainName = toString(value)
		case "cgroup_parent":
			ensureComposeCompat(service).CgroupParent = toString(value)
		case "cgroup":
			ensureComposeCompat(service).Cgroup = toString(value)
		case "cpu_count":
			compat := ensureComposeCompat(service)
			compat.CPUCount = int64(toInt(value))
			compat.CPUCountSet = true
		case "cpu_percent":
			if f, err := strconv.ParseFloat(toString(value), 32); err == nil {
				compat := ensureComposeCompat(service)
				compat.CPUPercent = float32(f)
				compat.CPUPercentSet = true
			}
		case "cpu_period":
			compat := ensureComposeCompat(service)
			compat.CPUPeriod = int64(toInt(value))
			compat.CPUPeriodSet = true
		case "cpu_rt_period":
			compat := ensureComposeCompat(service)
			compat.CPURTPeriod = int64(toInt(value))
			compat.CPURTPeriodSet = true
		case "cpu_rt_runtime":
			compat := ensureComposeCompat(service)
			compat.CPURTRuntime = int64(toInt(value))
			compat.CPURTRuntimeSet = true
		case "cpuset":
			ensureComposeCompat(service).CPUSet = toString(value)
		case "device_cgroup_rules":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse device_cgroup_rules: %w", err)
			}
			ensureComposeCompat(service).DeviceCgroupRules = values
		case "gpus":
			if items, ok := value.([]interface{}); ok {
				compat := ensureComposeCompat(service)
				compat.Gpus = compat.Gpus[:0]
				for _, item := range items {
					if mapped, ok := asMap(item); ok {
						compat.Gpus = append(compat.Gpus, copyStringInterfaceMap(mapped))
					}
				}
			}
		case "network_mode":
			ensureComposeCompat(service).NetworkMode = toString(value)
		case "oom_kill_disable":
			compat := ensureComposeCompat(service)
			compat.OomKillDisable = toBool(value)
			compat.OomKillDisableSet = true
		case "oom_score_adj":
			compat := ensureComposeCompat(service)
			compat.OomScoreAdj = int64(toInt(value))
			compat.OomScoreAdjSet = true
		case "scale":
			scale := toInt(value)
			ensureComposeCompat(service).Scale = &scale
		case "models":
			if mapped, ok := asMap(value); ok {
				ensureComposeCompat(service).Models = parseComposeModels(mapped)
			}
		case "storage_opt":
			options, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse storage_opt: %w", err)
			}
			ensureComposeCompat(service).StorageOpt = options
		case "use_api_socket":
			compat := ensureComposeCompat(service)
			compat.UseAPISocket = toBool(value)
			compat.UseAPISocketSet = true
		case "uts":
			ensureComposeCompat(service).Uts = toString(value)
		case "volume_driver":
			ensureComposeCompat(service).VolumeDriver = toString(value)
		case "label_file":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse label_file: %w", err)
			}
			ensureComposeCompat(service).LabelFiles = values
		case "external_links":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse external_links: %w", err)
			}
			ensureComposeCompat(service).ExternalLinks = values
		case "tmpfs":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse tmpfs: %w", err)
			}
			ensureComposeCompat(service).Tmpfs = values
		case "isolation":
			ensureComposeCompat(service).Isolation = toString(value)
		case "volumes_from":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse volumes_from: %w", err)
			}
			ensureComposeCompat(service).VolumesFrom = values
		case "group_add":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse group_add: %w", err)
			}
			service.GroupAdd = values
		case "sysctls":
			values, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse sysctls: %w", err)
			}
			service.Sysctls = values
		case "runtime":
			service.Runtime = toString(value)
		case "cap_add":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cap_add: %w", err)
			}
			service.CapAdd = values
		case "cap_drop":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cap_drop: %w", err)
			}
			service.CapDrop = values
		case "security_opt":
			values, err := toStringSlice(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse security_opt: %w", err)
			}
			service.SecurityOpt = values
		case "read_only":
			service.ReadOnlyRootFS = toBool(value)
			service.ReadOnlyRootFSSet = true
		case "init":
			init := toBool(value)
			service.Init = &init
		case "tty":
			service.Tty = toBool(value)
			service.TtySet = true
		case "stdin_open":
			service.StdinOpen = toBool(value)
			service.StdinOpenSet = true
		case "stop_signal":
			service.StopSignal = toString(value)
		case "stop_grace_period":
			service.StopGracePeriod = toString(value)
		case "cpu_shares":
			if intVal, ok := value.(int); ok {
				service.CPUShares = intVal
			}
		case "cpu_quota":
			if intVal, ok := value.(int); ok {
				service.CPUQuota = intVal
			}
		case "cpus":
			service.CPUs = toString(value)
		case "mem_limit":
			raw := toString(value)
			if raw != "" {
				if _, ok := service.Extensions[composeMemoryLimitExtension]; !ok {
					service.Extensions[composeMemoryLimitExtension] = raw
				}
				if quantity, err := resource.ParseQuantity(raw); err == nil {
					normalized := fmt.Sprintf("%d", quantity.Value())
					service.MemoryLimit = normalized
					service.MemLimit = normalized
					break
				}
			}
			service.MemoryLimit = raw
			service.MemLimit = service.MemoryLimit
		case "mem_reservation":
			raw := toString(value)
			if raw != "" {
				if _, ok := service.Extensions[composeMemoryReservationExtension]; !ok {
					service.Extensions[composeMemoryReservationExtension] = raw
				}
				if quantity, err := resource.ParseQuantity(raw); err == nil {
					service.MemReservation = fmt.Sprintf("%d", quantity.Value())
					break
				}
			}
			service.MemReservation = raw
		case "memswap_limit":
			raw := toString(value)
			if raw != "" {
				if _, ok := service.Extensions[composeMemorySwapExtension]; !ok {
					service.Extensions[composeMemorySwapExtension] = raw
				}
				if quantity, err := resource.ParseQuantity(raw); err == nil {
					service.MemorySwap = fmt.Sprintf("%d", quantity.Value())
					break
				}
			}
			service.MemorySwap = raw
		case "mem_swappiness":
			ensureComposeCompat(service).MemSwappiness = toString(value)
		case "ulimits":
			ulimits, err := parseUlimits(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ulimits: %w", err)
			}
			service.Ulimits = ulimits
		case "userns_mode":
			service.UserNSMode = toString(value)
		case "pull_policy":
			service.PullPolicy = toString(value)
		case "healthcheck":
			healthcheck, err := parseHealthCheck(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse healthcheck: %w", err)
			}
			service.HealthCheck = healthcheck
		case "deploy":
			deploy, err := parseDeploySpec(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse deploy: %w", err)
			}
			service.Deploy = deploy
			if deploy.Replicas > 0 {
				service.Replicas = deploy.Replicas
			}
		case composeFailoverExtension:
			failover, err := parseFailoverSpec(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", composeFailoverExtension, err)
			}
			service.Failover = failover
		case "labels":
			labels, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
			service.Labels = labels
		default:
			service.Extensions[key] = value
		}
	}

	return service, nil
}

func parseFailoverSpec(value interface{}) (*FailoverSpec, error) {
	mapped, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("failover must be a map")
	}
	return failoverSpecFromMap(mapped)
}

func serializeFailoverSpec(failover *FailoverSpec) map[string]interface{} {
	if failover == nil {
		return nil
	}
	data := map[string]interface{}{}
	if failover.EnabledSet || failover.Enabled {
		data["enabled"] = true
		if failover.EnabledSet {
			data["enabled"] = failover.Enabled
		}
	}
	if failover.Type != "" {
		data["type"] = failover.Type
	}
	if failover.Port != "" {
		data["port"] = failover.Port
	}
	if failover.HealthcheckPath != "" {
		data["healthcheck_path"] = failover.HealthcheckPath
	}
	if failover.HealthcheckInterval != "" {
		data["healthcheck_interval"] = failover.HealthcheckInterval
	}
	if failover.MaxRetries > 0 {
		data["max_retries"] = failover.MaxRetries
	}
	if failover.RedeployOnPeer {
		data["redeploy_on_peer"] = true
	}
	if failover.Singleton {
		data["singleton"] = true
	}
	if failover.SingletonElection != "" {
		data["singleton_election"] = failover.SingletonElection
	}
	if failover.Strategy != "" {
		data["strategy"] = failover.Strategy
	}
	if failover.PreferLocal {
		data["prefer_local"] = true
	}
	if len(failover.Nodes) > 0 {
		nodes := map[string]interface{}{}
		for name, node := range failover.Nodes {
			nodes[name] = serializeFailoverNode(node)
		}
		data["nodes"] = nodes
	}
	for key, value := range failover.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	return data
}

func serializeFailoverNode(node *FailoverNode) map[string]interface{} {
	data := map[string]interface{}{}
	if node == nil {
		return data
	}
	if node.Status != "" {
		data["status"] = node.Status
	}
	if node.LastSeen != "" {
		data["last_seen"] = node.LastSeen
	}
	if node.Priority > 0 {
		data["priority"] = node.Priority
	}
	if node.URL != "" {
		data["url"] = node.URL
	}
	if node.Weight > 0 {
		data["weight"] = node.Weight
	}
	for key, value := range node.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	return data
}

func serializeServiceNetworks(service *Service) interface{} {
	if service == nil || len(service.Networks) == 0 {
		return nil
	}
	if len(service.NetworkAttachments) == 0 {
		return append([]string{}, service.Networks...)
	}
	networks := map[string]interface{}{}
	for _, name := range service.Networks {
		attachment := service.NetworkAttachments[name]
		if networkAttachmentEmpty(attachment) {
			networks[name] = map[string]interface{}{}
			continue
		}
		networks[name] = serializeNetworkAttachment(attachment)
	}
	for name, attachment := range service.NetworkAttachments {
		if _, exists := networks[name]; exists {
			continue
		}
		networks[name] = serializeNetworkAttachment(attachment)
	}
	return networks
}

func serializeNetworkAttachment(attachment *NetworkAttachment) map[string]interface{} {
	data := map[string]interface{}{}
	if attachment == nil {
		return data
	}
	if len(attachment.Aliases) > 0 {
		data["aliases"] = append([]string{}, attachment.Aliases...)
	}
	if attachment.InterfaceName != "" {
		data["interface_name"] = attachment.InterfaceName
	}
	if attachment.IPv4Address != "" {
		data["ipv4_address"] = attachment.IPv4Address
	}
	if attachment.IPv6Address != "" {
		data["ipv6_address"] = attachment.IPv6Address
	}
	if len(attachment.LinkLocalIPs) > 0 {
		data["link_local_ips"] = append([]string{}, attachment.LinkLocalIPs...)
	}
	if attachment.MacAddress != "" {
		data["mac_address"] = attachment.MacAddress
	}
	if len(attachment.DriverOpts) > 0 {
		data["driver_opts"] = copyStringMap(attachment.DriverOpts)
	}
	if attachment.GWPriority != 0 {
		data["gw_priority"] = attachment.GWPriority
	}
	if attachment.Priority != 0 {
		data["priority"] = attachment.Priority
	}
	for key, value := range attachment.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	return data
}

func parseServiceNetworks(value interface{}) ([]string, map[string]*NetworkAttachment, error) {
	switch typed := value.(type) {
	case []interface{}:
		networks, err := toStringSlice(value)
		return networks, nil, err
	case []string:
		return append([]string{}, typed...), nil, nil
	case map[string]interface{}:
		names := make([]string, 0, len(typed))
		attachments := map[string]*NetworkAttachment{}
		for name, raw := range typed {
			names = append(names, name)
			attachment := &NetworkAttachment{Name: name, Extensions: map[string]interface{}{}}
			if raw != nil {
				mapped, ok := asMap(raw)
				if !ok {
					return nil, nil, fmt.Errorf("network %s must be a map", name)
				}
				for key, rawValue := range mapped {
					switch key {
					case "aliases":
						values, err := toStringSlice(rawValue)
						if err != nil {
							return nil, nil, fmt.Errorf("network %s aliases: %w", name, err)
						}
						attachment.Aliases = values
					case "interface_name":
						attachment.InterfaceName = toString(rawValue)
					case "ipv4_address":
						attachment.IPv4Address = toString(rawValue)
					case "ipv6_address":
						attachment.IPv6Address = toString(rawValue)
					case "link_local_ips":
						values, err := toStringSlice(rawValue)
						if err != nil {
							return nil, nil, fmt.Errorf("network %s link_local_ips: %w", name, err)
						}
						attachment.LinkLocalIPs = values
					case "mac_address":
						attachment.MacAddress = toString(rawValue)
					case "driver_opts":
						options, err := toStringMap(rawValue)
						if err != nil {
							return nil, nil, fmt.Errorf("network %s driver_opts: %w", name, err)
						}
						attachment.DriverOpts = options
					case "gw_priority":
						attachment.GWPriority = toInt(rawValue)
					case "priority":
						attachment.Priority = toInt(rawValue)
					default:
						attachment.Extensions[key] = rawValue
					}
				}
			}
			if len(attachment.Extensions) == 0 {
				attachment.Extensions = nil
			}
			if !networkAttachmentEmpty(attachment) {
				attachments[name] = attachment
			}
		}
		sort.Strings(names)
		if len(attachments) == 0 {
			attachments = nil
		}
		return names, attachments, nil
	default:
		networks, err := toStringSlice(value)
		return networks, nil, err
	}
}

func networkAttachmentEmpty(attachment *NetworkAttachment) bool {
	return attachment == nil ||
		(len(attachment.Aliases) == 0 &&
			attachment.InterfaceName == "" &&
			attachment.IPv4Address == "" &&
			attachment.IPv6Address == "" &&
			len(attachment.LinkLocalIPs) == 0 &&
			attachment.MacAddress == "" &&
			len(attachment.DriverOpts) == 0 &&
			attachment.GWPriority == 0 &&
			attachment.Priority == 0 &&
			len(attachment.Extensions) == 0)
}

func parseDeploySpec(value interface{}) (*DeploySpec, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("deploy must be a map")
	}

	deploy := &DeploySpec{Labels: map[string]string{}, Extensions: map[string]interface{}{}}
	for key, val := range data {
		switch key {
		case "mode":
			deploy.Mode = toString(val)
		case "endpoint_mode":
			deploy.EndpointMode = toString(val)
		case "replicas":
			deploy.Replicas = toInt(val)
		case composeSwarmJobExtension, "x-swarm-job-spec":
			job, err := parseSwarmJobSpec(val)
			if err != nil {
				return nil, err
			}
			deploy.Job = job
		case "labels":
			labels, err := toStringMap(val)
			if err != nil {
				return nil, err
			}
			deploy.Labels = labels
		case "placement":
			placement, err := parsePlacementSpec(val)
			if err != nil {
				return nil, err
			}
			deploy.Placement = placement
		case "resources":
			resources, err := parseResourceSpec(val)
			if err != nil {
				return nil, err
			}
			deploy.Resources = resources
		case "update_config", "x-nomad-update":
			update, err := parseUpdatePolicy(val)
			if err != nil {
				return nil, err
			}
			deploy.UpdateConfig = update
			if key != "update_config" {
				deploy.Extensions[key] = val
			}
		case "migrate", "x-nomad-migrate":
			migrate, err := parseMigratePolicy(val)
			if err != nil {
				return nil, err
			}
			deploy.MigrateConfig = migrate
			if key != "migrate" {
				deploy.Extensions[key] = val
			}
		case "reschedule", "x-nomad-reschedule":
			reschedule, err := parseReschedulePolicy(val)
			if err != nil {
				return nil, err
			}
			deploy.RescheduleConfig = reschedule
			if key != "reschedule" {
				deploy.Extensions[key] = val
			}
		case "rollback_config":
			rollback, err := parseUpdatePolicy(val)
			if err != nil {
				return nil, err
			}
			deploy.RollbackConfig = rollback
		case "restart_policy":
			restart, err := parseRestartPolicy(val)
			if err != nil {
				return nil, err
			}
			deploy.RestartPolicy = restart
		default:
			deploy.Extensions[key] = val
		}
	}
	if len(deploy.Labels) == 0 {
		deploy.Labels = nil
	}
	if len(deploy.Extensions) == 0 {
		deploy.Extensions = nil
	}
	return deploy, nil
}

func parseSwarmJobSpec(value interface{}) (*SwarmJobSpec, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("%s must be a map", composeSwarmJobExtension)
	}
	return parseSwarmJobSpecMap(data)
}

func parseSwarmJobSpecJSON(data []byte) (*SwarmJobSpec, error) {
	var mapped map[string]interface{}
	if err := json.Unmarshal(data, &mapped); err != nil {
		return nil, err
	}
	return parseSwarmJobSpecMap(mapped)
}

func parseSwarmJobSpecMap(data map[string]interface{}) (*SwarmJobSpec, error) {
	job := &SwarmJobSpec{Extensions: map[string]interface{}{}}
	for key, val := range data {
		switch key {
		case "max_concurrent", "maxConcurrent":
			job.MaxConcurrent = toInt(val)
			job.maxConcurrentSet = true
		case "total_completions", "totalCompletions":
			job.TotalCompletions = toInt(val)
			job.totalCompletionsSet = true
		case "completion_mode", "completionMode":
			job.CompletionMode = toString(val)
			job.completionModeSet = true
		case "suspend":
			flag := toBool(val)
			job.Suspend = &flag
		case "backoff_limit", "backoffLimit":
			job.BackoffLimit = toInt(val)
			job.backoffLimitSet = true
		case "backoff_limit_per_index", "backoffLimitPerIndex":
			job.BackoffLimitPerIndex = toInt(val)
			job.backoffLimitPerIndexSet = true
		case "ttl_seconds_after_finished", "ttlSecondsAfterFinished":
			job.TTLSecondsAfterFinished = toInt(val)
			job.ttlSecondsAfterFinishedSet = true
		case "extensions":
			if nested, ok := asMap(val); ok {
				for nestedKey, nestedValue := range nested {
					job.Extensions[nestedKey] = deepCopyValue(nestedValue)
				}
			}
		default:
			job.Extensions[key] = deepCopyValue(val)
		}
	}
	if len(job.Extensions) == 0 {
		job.Extensions = nil
	}
	if isEmptySwarmJobSpec(job) {
		return nil, nil
	}
	return job, nil
}

func parseDevelopConfig(value interface{}) (*DevelopConfig, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("develop must be a map")
	}
	develop := &DevelopConfig{
		Extensions: map[string]interface{}{},
	}
	for key, val := range data {
		switch key {
		case "watch":
			items, ok := val.([]interface{})
			if !ok {
				return nil, fmt.Errorf("develop.watch must be a list")
			}
			for _, item := range items {
				watch, err := parseDevelopWatch(item)
				if err != nil {
					return nil, err
				}
				develop.Watch = append(develop.Watch, watch)
			}
		default:
			develop.Extensions[key] = val
		}
	}
	if len(develop.Extensions) == 0 {
		develop.Extensions = nil
	}
	if len(develop.Watch) == 0 && len(develop.Extensions) == 0 {
		return nil, nil
	}
	return develop, nil
}

func parseDevelopWatch(value interface{}) (DevelopWatch, error) {
	data, ok := asMap(value)
	if !ok {
		return DevelopWatch{}, fmt.Errorf("develop.watch item must be a map")
	}
	watch := DevelopWatch{
		Extensions: map[string]interface{}{},
	}
	for key, val := range data {
		switch key {
		case "path":
			watch.Path = toString(val)
		case "action":
			watch.Action = toString(val)
		case "target":
			watch.Target = toString(val)
		case "exec":
			hook, err := parseServiceHook(val, "develop.watch.exec")
			if err != nil {
				return DevelopWatch{}, err
			}
			watch.Exec = hook
		case "include":
			values, err := toStringSlice(val)
			if err != nil {
				return DevelopWatch{}, fmt.Errorf("develop.watch.include: %w", err)
			}
			watch.Include = values
		case "ignore":
			values, err := toStringSlice(val)
			if err != nil {
				return DevelopWatch{}, fmt.Errorf("develop.watch.ignore: %w", err)
			}
			watch.Ignore = values
		case "initial_sync":
			watch.InitialSync = toBool(val)
		default:
			watch.Extensions[key] = val
		}
	}
	if len(watch.Extensions) == 0 {
		watch.Extensions = nil
	}
	return watch, nil
}

func parseServiceHookList(value interface{}, field string) ([]ServiceHook, error) {
	items, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be a list", field)
	}
	hooks := make([]ServiceHook, 0, len(items))
	for _, item := range items {
		hook, err := parseServiceHook(item, field)
		if err != nil {
			return nil, err
		}
		if hook != nil {
			hooks = append(hooks, *hook)
		}
	}
	return hooks, nil
}

func parseServiceHook(value interface{}, field string) (*ServiceHook, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("%s item must be a map", field)
	}
	hook := &ServiceHook{
		Extensions: map[string]interface{}{},
	}
	for key, val := range data {
		switch key {
		case "command":
			values, err := toStringSlice(val)
			if err != nil {
				if single := toString(val); single != "" {
					values = []string{single}
				} else {
					return nil, fmt.Errorf("%s.command: %w", field, err)
				}
			}
			hook.Command = values
		case "image":
			hook.Image = toString(val)
		case "user":
			hook.User = toString(val)
		case "privileged":
			hook.Privileged = toBool(val)
		case "working_dir":
			hook.WorkingDir = toString(val)
		case "environment":
			env, err := toStringPtrMap(val)
			if err != nil {
				return nil, fmt.Errorf("%s.environment: %w", field, err)
			}
			hook.Environment = env
		case "per_replica":
			hook.PerReplica = toBool(val)
		default:
			hook.Extensions[key] = val
		}
	}
	if len(hook.Extensions) == 0 {
		hook.Extensions = nil
	}
	return hook, nil
}

func ensureLifecycle(service *Service) *LifecycleHooks {
	if service.Lifecycle == nil {
		service.Lifecycle = &LifecycleHooks{}
	}
	return service.Lifecycle
}

func ensureComposeCompat(service *Service) *ComposeCompat {
	if service.ComposeCompat == nil {
		service.ComposeCompat = &ComposeCompat{}
	}
	return service.ComposeCompat
}

func parseComposeCompatExtension(service *Service, value interface{}) {
	data, ok := asMap(value)
	if !ok || len(data) == 0 || service == nil {
		return
	}
	compat := composeCompatFromExtensionMap(data)
	if compat == nil {
		return
	}
	mergeComposeCompat(ensureComposeCompat(service), compat)
	if service.Failover == nil {
		for _, key := range []string{composeFailoverExtension, "failover"} {
			raw, ok := compat.Extensions[key]
			if !ok {
				continue
			}
			if failover, err := parseFailoverSpec(raw); err == nil {
				service.Failover = failover
				break
			}
		}
	}
}

func serializeComposeCompatExtension(compat *ComposeCompat) map[string]interface{} {
	if compat == nil {
		return nil
	}
	data := map[string]interface{}{}
	if compat.PullRefreshAfter != "" {
		data["pull_refresh_after"] = compat.PullRefreshAfter
	}
	if compat.MemSwappiness != "" {
		if value, err := strconv.Atoi(compat.MemSwappiness); err == nil {
			data["mem_swappiness"] = value
		} else {
			data["mem_swappiness"] = compat.MemSwappiness
		}
	}
	if compat.VolumeDriver != "" {
		data["volume_driver"] = compat.VolumeDriver
	}
	for key, value := range compat.Extensions {
		data[key] = deepCopyValue(value)
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

func parsePlacementSpec(value interface{}) (*PlacementSpec, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("placement must be a map")
	}
	placement := &PlacementSpec{Extensions: map[string]interface{}{}}
	for key, val := range data {
		switch key {
		case "constraints":
			values, err := toStringSlice(val)
			if err != nil {
				return nil, err
			}
			placement.Constraints = values
		case "preferences":
			values, extensions, err := parsePlacementPreferences(val)
			if err != nil {
				return nil, err
			}
			placement.Preferences = values
			placement.PreferenceExtensions = extensions
		case "max_replicas_per_node":
			placement.MaxReplicasPerNode = toInt(val)
		default:
			placement.Extensions[key] = val
		}
	}
	if len(placement.Extensions) == 0 {
		placement.Extensions = nil
	}
	return placement, nil
}

func parsePlacementPreferences(value interface{}) ([]string, []map[string]interface{}, error) {
	var preferences []string
	var extensions []map[string]interface{}
	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			preference, extension, err := parsePlacementPreference(item)
			if err != nil {
				return nil, nil, err
			}
			if preference != "" || len(extension) > 0 {
				preferences = append(preferences, preference)
				extensions = append(extensions, extension)
			}
		}
	case []string:
		for _, item := range v {
			if preference := normalizePlacementPreference(item); preference != "" {
				preferences = append(preferences, preference)
				extensions = append(extensions, nil)
			}
		}
	case string:
		if preference := normalizePlacementPreference(v); preference != "" {
			preferences = append(preferences, preference)
			extensions = append(extensions, nil)
		}
	default:
		return nil, nil, fmt.Errorf("cannot convert placement preferences to slice: %T", value)
	}
	if !hasPlacementPreferenceExtensions(extensions) {
		extensions = nil
	}
	return preferences, extensions, nil
}

func parsePlacementPreference(value interface{}) (string, map[string]interface{}, error) {
	if data, ok := asMap(value); ok {
		extensions := map[string]interface{}{}
		for key, val := range data {
			switch key {
			case "spread":
			default:
				extensions[key] = val
			}
		}
		if len(extensions) == 0 {
			extensions = nil
		}
		if raw := toString(data["x-bolabaden-preference"]); raw != "" {
			return normalizePlacementPreference(raw), extensions, nil
		}
		if spread := toString(data["spread"]); spread != "" {
			return "spread=" + spread, extensions, nil
		}
		return "", nil, nil
	}
	return normalizePlacementPreference(toString(value)), nil, nil
}

func hasPlacementPreferenceExtensions(extensions []map[string]interface{}) bool {
	for _, extension := range extensions {
		if len(extension) > 0 {
			return true
		}
	}
	return false
}

func normalizePlacementPreference(preference string) string {
	preference = strings.TrimSpace(preference)
	if preference == "" {
		return ""
	}
	if strings.HasPrefix(preference, "spread:") {
		return "spread=" + strings.TrimSpace(strings.TrimPrefix(preference, "spread:"))
	}
	return preference
}

func parseResourceSpec(value interface{}) (*ResourceSpec, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("resources must be a map")
	}
	return parseResourceSpecMap(data)
}

func parseResourceSpecJSON(data []byte) (*ResourceSpec, error) {
	var mapped map[string]interface{}
	if err := json.Unmarshal(data, &mapped); err != nil {
		return nil, err
	}
	return parseResourceSpecMap(mapped)
}

func parseResourceSpecMap(data map[string]interface{}) (*ResourceSpec, error) {
	resources := &ResourceSpec{
		Extensions:            map[string]interface{}{},
		LimitExtensions:       map[string]interface{}{},
		ReservationExtensions: map[string]interface{}{},
	}
	for key, val := range data {
		switch key {
		case "limits":
			limits, ok := asMap(val)
			if !ok {
				return nil, fmt.Errorf("limits must be a map")
			}
			resources.CPULimit = toString(limits["cpus"])
			resources.MemoryLimit = toString(limits["memory"])
			resources.EphemeralStorageLimit = toString(limits["ephemeral-storage"])
			if resources.EphemeralStorageLimit == "" {
				resources.EphemeralStorageLimit = toString(limits["x-kubernetes-ephemeral-storage-limit"])
				if resources.EphemeralStorageLimit != "" {
					if resources.LimitExtensions == nil {
						resources.LimitExtensions = map[string]interface{}{}
					}
					resources.LimitExtensions["x-kubernetes-ephemeral-storage-limit"] = resources.EphemeralStorageLimit
				}
			}
			resources.PidsLimit = int64(toInt(limits["pids"]))
			if _, ok := limits["pids"]; ok {
				resources.pidsLimitSet = true
			}
			for nestedKey, nestedValue := range limits {
				switch nestedKey {
				case "cpus", "memory", "ephemeral-storage", "pids":
				default:
					resources.LimitExtensions[nestedKey] = nestedValue
				}
			}
		case "reservations":
			reservations, ok := asMap(val)
			if !ok {
				return nil, fmt.Errorf("reservations must be a map")
			}
			resources.CPUReservation = toString(reservations["cpus"])
			resources.MemoryReservation = toString(reservations["memory"])
			resources.EphemeralStorageReservation = toString(reservations["ephemeral-storage"])
			if resources.EphemeralStorageReservation == "" {
				resources.EphemeralStorageReservation = toString(reservations["x-kubernetes-ephemeral-storage-reservation"])
				if resources.EphemeralStorageReservation != "" {
					if resources.ReservationExtensions == nil {
						resources.ReservationExtensions = map[string]interface{}{}
					}
					resources.ReservationExtensions["x-kubernetes-ephemeral-storage-reservation"] = resources.EphemeralStorageReservation
				}
			}
			resources.PidsReservation = int64(toInt(reservations["pids"]))
			if _, ok := reservations["pids"]; ok {
				resources.pidsReservationSet = true
			}
			for nestedKey, nestedValue := range reservations {
				switch nestedKey {
				case "cpus", "memory", "ephemeral-storage", "pids", "devices", "generic_resources":
				default:
					resources.ReservationExtensions[nestedKey] = nestedValue
				}
			}
			if devices, ok := reservations["devices"]; ok {
				parsed, err := parseResourceDevices(devices)
				if err != nil {
					return nil, err
				}
				resources.Devices = parsed
			}
			if generic, ok := reservations["generic_resources"]; ok {
				parsed, err := parseGenericResources(generic)
				if err != nil {
					return nil, err
				}
				resources.GenericResources = parsed
			}
		case "cpus", "memory", "pids", "devices", "generic_resources":
		default:
			resources.Extensions[key] = val
		}
	}
	if len(resources.Extensions) == 0 {
		resources.Extensions = nil
	}
	if len(resources.LimitExtensions) == 0 {
		resources.LimitExtensions = nil
	}
	if len(resources.ReservationExtensions) == 0 {
		resources.ReservationExtensions = nil
	}
	return resources, nil
}

func composeCPUValue(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if cpu, err := strconv.ParseFloat(value, 64); err == nil {
		return cpu
	}
	quantity, err := resource.ParseQuantity(value)
	if err != nil {
		return value
	}
	return quantity.AsApproximateFloat64()
}

func composeMemoryValue(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	quantity, err := resource.ParseQuantity(value)
	if err != nil {
		return value
	}
	return fmt.Sprintf("%d", quantity.Value())
}

func parseResourceDevices(value interface{}) ([]ResourceDevice, error) {
	values, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("resource devices must be a list")
	}
	devices := make([]ResourceDevice, 0, len(values))
	for _, item := range values {
		data, ok := asMap(item)
		if !ok {
			return nil, fmt.Errorf("resource device must be a map")
		}
		device := ResourceDevice{
			Driver:     toString(data["driver"]),
			Count:      toString(data["count"]),
			Extensions: map[string]interface{}{},
		}
		for key, val := range data {
			switch key {
			case "driver", "count", "capabilities", "device_ids", "options":
			default:
				device.Extensions[key] = val
			}
		}
		if capabilities, ok := data["capabilities"]; ok {
			values, err := toStringSlice(capabilities)
			if err != nil {
				return nil, err
			}
			device.Capabilities = values
		}
		if ids, ok := data["device_ids"]; ok {
			values, err := toStringSlice(ids)
			if err != nil {
				return nil, err
			}
			device.DeviceIDs = values
		}
		if options, ok := data["options"]; ok {
			values, err := toStringMap(options)
			if err != nil {
				return nil, err
			}
			device.Options = values
		}
		if len(device.Extensions) == 0 {
			device.Extensions = nil
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func parseGenericResources(value interface{}) ([]GenericResource, error) {
	values, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("generic_resources must be a list")
	}
	resources := make([]GenericResource, 0, len(values))
	for _, item := range values {
		data, ok := asMap(item)
		if !ok {
			return nil, fmt.Errorf("generic resource must be a map")
		}
		resource := GenericResource{
			Extensions:         map[string]interface{}{},
			DiscreteExtensions: map[string]interface{}{},
		}
		for key, val := range data {
			switch key {
			case "kind", "value", "discrete_resource_spec":
			default:
				resource.Extensions[key] = val
			}
		}
		if discrete, ok := asMap(data["discrete_resource_spec"]); ok {
			resource.Kind = toString(discrete["kind"])
			resource.Value = toString(discrete["value"])
			for key, val := range discrete {
				switch key {
				case "kind", "value":
				default:
					resource.DiscreteExtensions[key] = val
				}
			}
		}
		if len(resource.Extensions) == 0 {
			resource.Extensions = nil
		}
		if len(resource.DiscreteExtensions) == 0 {
			resource.DiscreteExtensions = nil
		}
		if resource.Kind == "" && resource.Value == "" && resource.Extensions == nil && resource.DiscreteExtensions == nil {
			continue
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func parseUpdatePolicy(value interface{}) (*UpdatePolicy, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("update_config must be a map")
	}
	_, parallelismSet := data["parallelism"]
	update := &UpdatePolicy{
		Parallelism:      toInt(data["parallelism"]),
		ParallelismSet:   parallelismSet,
		Delay:            toString(data["delay"]),
		Monitor:          toString(data["monitor"]),
		MaxFailureRatio:  toString(data["max_failure_ratio"]),
		Order:            toString(data["order"]),
		OnFailure:        toString(data["failure_action"]),
		HealthCheck:      toString(data["health_check"]),
		MinHealthyTime:   toString(data["min_healthy_time"]),
		HealthyDeadline:  toString(data["healthy_deadline"]),
		ProgressDeadline: toString(data["progress_deadline"]),
		AutoRevert:       toBool(data["auto_revert"]),
		AutoPromote:      toBool(data["auto_promote"]),
		Canary:           toInt(data["canary"]),
		Stagger:          toString(data["stagger"]),
		Extensions:       map[string]interface{}{},
	}
	if _, ok := data["auto_revert"]; ok {
		update.AutoRevertSet = true
	}
	if _, ok := data["auto_promote"]; ok {
		update.AutoPromoteSet = true
	}
	if _, ok := data["canary"]; ok {
		update.CanarySet = true
	}
	for key, val := range data {
		switch key {
		case "parallelism", "delay", "monitor", "max_failure_ratio", "order", "failure_action", "health_check", "min_healthy_time", "healthy_deadline", "progress_deadline", "auto_revert", "auto_promote", "canary", "stagger":
		default:
			update.Extensions[key] = val
		}
	}
	if len(update.Extensions) == 0 {
		update.Extensions = nil
	}
	return update, nil
}

func parseMigratePolicy(value interface{}) (*MigratePolicy, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("migrate must be a map")
	}
	migrate := &MigratePolicy{
		MaxParallel:     toInt(data["max_parallel"]),
		HealthCheck:     toString(data["health_check"]),
		MinHealthyTime:  toString(data["min_healthy_time"]),
		HealthyDeadline: toString(data["healthy_deadline"]),
		Extensions:      map[string]interface{}{},
	}
	for key, val := range data {
		switch key {
		case "max_parallel", "health_check", "min_healthy_time", "healthy_deadline":
		default:
			migrate.Extensions[key] = val
		}
	}
	if len(migrate.Extensions) == 0 {
		migrate.Extensions = nil
	}
	return migrate, nil
}

func parseReschedulePolicy(value interface{}) (*ReschedulePolicy, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("reschedule must be a map")
	}
	reschedule := &ReschedulePolicy{
		Attempts:      toInt(data["attempts"]),
		Interval:      toString(data["interval"]),
		Delay:         toString(data["delay"]),
		DelayFunction: toString(data["delay_function"]),
		MaxDelay:      toString(data["max_delay"]),
		Unlimited:     toBool(data["unlimited"]),
		Extensions:    map[string]interface{}{},
	}
	for key, val := range data {
		switch key {
		case "attempts", "interval", "delay", "delay_function", "max_delay", "unlimited":
		default:
			reschedule.Extensions[key] = val
		}
	}
	if len(reschedule.Extensions) == 0 {
		reschedule.Extensions = nil
	}
	return reschedule, nil
}

func parseRestartPolicy(value interface{}) (*RestartPolicy, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("restart_policy must be a map")
	}
	restart := &RestartPolicy{
		Condition:   toString(data["condition"]),
		Delay:       toString(data["delay"]),
		MaxAttempts: toInt(data["max_attempts"]),
		Window:      toString(data["window"]),
		Extensions:  map[string]interface{}{},
	}
	for key, val := range data {
		switch key {
		case "condition", "delay", "max_attempts", "window":
		default:
			restart.Extensions[key] = val
		}
	}
	if len(restart.Extensions) == 0 {
		restart.Extensions = nil
	}
	return restart, nil
}

func parsePorts(value interface{}) ([]PortMapping, error) {
	var ports []PortMapping

	switch v := value.(type) {
	case []interface{}:
		for _, portValue := range v {
			var (
				port PortMapping
				err  error
			)
			if portMap, ok := asMap(portValue); ok {
				port = parsePortMappingMap(portMap)
			} else {
				port, err = parsePortMapping(toString(portValue))
				if err != nil {
					return nil, err
				}
			}
			ports = append(ports, port)
		}
	case []string:
		for _, portStr := range v {
			port, err := parsePortMapping(portStr)
			if err != nil {
				return nil, err
			}
			ports = append(ports, port)
		}
	}

	return ports, nil
}

func parsePortMappingMap(data map[string]interface{}) PortMapping {
	port := PortMapping{
		Name:          toString(data["name"]),
		HostIP:        toString(data["host_ip"]),
		HostPort:      toString(data["published"]),
		ContainerPort: toString(data["target"]),
		Protocol:      toString(data["protocol"]),
		AppProtocol:   toString(data["app_protocol"]),
		Mode:          toString(data["mode"]),
		NodePort:      toString(data["x-kubernetes-node-port"]),
		Extensions:    map[string]interface{}{},
	}
	for key, value := range data {
		if key != "x-kubernetes-node-port" {
			port.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(port.Extensions) == 0 {
		port.Extensions = nil
	}
	if port.HostPort == "" {
		port.HostPort = toString(data["published_port"])
	}
	if port.Protocol == "" {
		port.Protocol = "tcp"
	}
	return port
}

func parsePortMapping(portStr string) (PortMapping, error) {
	// Handle formats like:
	// - "8080:80"
	// - "127.0.0.1:8080:80"
	// - "8080:80/tcp"
	// - "8080:80/udp"
	// - "${VAR:-8080}:${VAR:-80}"
	// - "${VAR:-8080}:${VAR:-80}/tcp"

	// First expand environment variables in the entire string
	expandedPortStr := expandEnvVars(portStr)

	parts := strings.Split(expandedPortStr, "/")
	protocol := "tcp"
	if len(parts) > 1 {
		protocol = parts[1]
	}

	hostParts := strings.Split(parts[0], ":")
	if len(hostParts) == 2 {
		return PortMapping{
			HostPort:      hostParts[0],
			ContainerPort: hostParts[1],
			Protocol:      protocol,
		}, nil
	} else if len(hostParts) == 3 {
		return PortMapping{
			HostIP:        hostParts[0],
			HostPort:      hostParts[1],
			ContainerPort: hostParts[2],
			Protocol:      protocol,
		}, nil
	}

	return PortMapping{}, fmt.Errorf("invalid port mapping format: %s", portStr)
}

// expandEnvVars handles basic environment variable substitution in strings
// Supports patterns like ${VAR:-default} and ${VAR}
func parseDependsOn(value interface{}) ([]string, error) {
	dependencies, err := parseDependencySpecs(value)
	if err != nil {
		return nil, err
	}
	return dependencyNames(dependencies), nil
}

func parseDependencySpecs(value interface{}) ([]DependencySpec, error) {
	switch v := value.(type) {
	case []interface{}:
		var dependencies []DependencySpec
		for _, item := range v {
			name := toString(item)
			if name != "" {
				dependencies = append(dependencies, DependencySpec{Name: name})
			}
		}
		return dependencies, nil
	case []string:
		var dependencies []DependencySpec
		for _, name := range v {
			if name != "" {
				dependencies = append(dependencies, DependencySpec{Name: name})
			}
		}
		return dependencies, nil
	default:
		data, ok := asMap(value)
		if !ok {
			return nil, fmt.Errorf("cannot convert to dependency specs: %T", value)
		}
		var dependencies []DependencySpec
		for serviceName, raw := range data {
			if serviceName == "" {
				continue
			}
			dependency := DependencySpec{Name: serviceName}
			if depMap, ok := asMap(raw); ok {
				dependency.Extensions = map[string]interface{}{}
				dependency.Condition = toString(depMap["condition"])
				dependency.Restart = toBool(depMap["restart"])
				if required, exists := depMap["required"]; exists {
					requiredValue := toBool(required)
					dependency.Required = &requiredValue
				}
				for key, value := range depMap {
					dependency.Extensions[key] = value
				}
				if len(dependency.Extensions) == 0 {
					dependency.Extensions = nil
				}
			}
			dependencies = append(dependencies, dependency)
		}
		return dependencies, nil
	}
}

func dependencyNames(dependencies []DependencySpec) []string {
	names := make([]string, 0, len(dependencies))
	seen := map[string]struct{}{}
	for _, dependency := range dependencies {
		if dependency.Name == "" {
			continue
		}
		if _, ok := seen[dependency.Name]; ok {
			continue
		}
		seen[dependency.Name] = struct{}{}
		names = append(names, dependency.Name)
	}
	return names
}

func dependencySpecsFromNames(names []string) []DependencySpec {
	dependencies := make([]DependencySpec, 0, len(names))
	for _, name := range names {
		if name != "" {
			dependencies = append(dependencies, DependencySpec{Name: name})
		}
	}
	return dependencies
}

func dependencyHasMetadata(dependency DependencySpec) bool {
	return dependency.Condition != "" || dependency.Restart || dependency.Required != nil || len(dependency.Extensions) > 0
}

func serviceDependencies(service *Service) []DependencySpec {
	if len(service.Dependencies) > 0 {
		return service.Dependencies
	}
	return dependencySpecsFromNames(service.DependsOn)
}

func serializeDependencySpecs(app *Application, dependencies []DependencySpec) interface{} {
	if len(dependencies) == 0 {
		return nil
	}
	longSyntax := false
	filtered := make([]DependencySpec, 0, len(dependencies))
	for _, dependency := range dependencies {
		if dependency.Name == "" {
			continue
		}
		if app != nil {
			if _, ok := app.Services[dependency.Name]; !ok {
				continue
			}
		}
		filtered = append(filtered, dependency)
		if dependencyHasMetadata(dependency) {
			longSyntax = true
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	if !longSyntax {
		return dependencyNames(filtered)
	}
	result := map[string]interface{}{}
	for _, dependency := range filtered {
		item := map[string]interface{}{}
		if dependency.Condition != "" {
			item["condition"] = dependency.Condition
		} else {
			// compose-go rejects long-form depends_on entries without a condition.
			// Default to the standard Compose behavior so metadata-only extensions stay compatible.
			item["condition"] = "service_started"
		}
		if dependency.Restart {
			item["restart"] = true
		}
		if dependency.Required != nil {
			item["required"] = *dependency.Required
		}
		for key, value := range dependency.Extensions {
			item[composeApplicationExtensionKey(key)] = value
		}
		result[dependency.Name] = item
	}
	return result
}

func composeModeReferenceExists(app *Application, mode string) bool {
	if app == nil {
		return false
	}
	if strings.HasPrefix(mode, "service:") {
		name := strings.TrimPrefix(mode, "service:")
		if name == "" {
			return false
		}
		_, ok := app.Services[name]
		return ok
	}
	return true
}

func boolPtr(value bool) *bool {
	return &value
}

func boolPtrFromInterface(value interface{}) *bool {
	if value == nil {
		return nil
	}
	if b, ok := value.(bool); ok {
		return &b
	}
	str := strings.ToLower(strings.TrimSpace(toString(value)))
	switch str {
	case "true", "1", "yes":
		return boolPtr(true)
	case "false", "0", "no":
		return boolPtr(false)
	default:
		return nil
	}
}

func stringInSlice(value string, values []string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func appendUniqueName(values []string, value string) []string {
	if value == "" || stringInSlice(value, values) {
		return values
	}
	return append(values, value)
}

func appendUniqueDependency(dependencies []DependencySpec, dependency DependencySpec) []DependencySpec {
	if dependency.Name == "" {
		return dependencies
	}
	for i, existing := range dependencies {
		if existing.Name == dependency.Name {
			if existing.Condition == "" {
				dependencies[i].Condition = dependency.Condition
			}
			if !existing.Restart {
				dependencies[i].Restart = dependency.Restart
			}
			if existing.Required == nil {
				dependencies[i].Required = dependency.Required
			}
			if len(existing.Extensions) == 0 && len(dependency.Extensions) > 0 {
				dependencies[i].Extensions = copyStringInterfaceMap(dependency.Extensions)
			}
			return dependencies
		}
	}
	return append(dependencies, dependency)
}

func expandEnvVars(s string) string {
	// Simple regex-based replacement for common patterns
	// ${VAR:-default} -> default
	// ${VAR} -> "" (empty if not set)

	result := s

	// Handle ${VAR:-default} pattern
	for strings.Contains(result, "${") && strings.Contains(result, ":-") && strings.Contains(result, "}") {
		start := strings.Index(result, "${")
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start

		varPart := result[start+2 : end]
		colonDash := strings.Index(varPart, ":-")
		if colonDash != -1 {
			varName := varPart[:colonDash]
			defaultValue := varPart[colonDash+2:]

			// Try to get environment variable
			if envValue := os.Getenv(varName); envValue != "" {
				result = strings.Replace(result, result[start:end+1], envValue, 1)
			} else {
				result = strings.Replace(result, result[start:end+1], defaultValue, 1)
			}
		} else {
			break
		}
	}
	return result
}

func parseVolumes(value interface{}) ([]VolumeMount, error) {
	var volumes []VolumeMount

	switch v := value.(type) {
	case []interface{}:
		for _, volValue := range v {
			var (
				volume VolumeMount
				err    error
			)
			if volMap, ok := asMap(volValue); ok {
				volume = parseVolumeMountMap(volMap)
			} else {
				volume, err = parseVolumeMount(toString(volValue))
				if err != nil {
					return nil, err
				}
			}
			volumes = append(volumes, volume)
		}
	case []string:
		for _, volStr := range v {
			volume, err := parseVolumeMount(volStr)
			if err != nil {
				return nil, err
			}
			volumes = append(volumes, volume)
		}
	}

	return volumes, nil
}

func parseFileRefs(value interface{}) ([]FileRef, error) {
	var refs []FileRef
	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			ref, err := parseFileRef(item)
			if err != nil {
				return nil, err
			}
			refs = append(refs, ref)
		}
	case []string:
		for _, item := range v {
			refs = append(refs, FileRef{Source: item})
		}
	}
	return refs, nil
}

func parseFileRef(value interface{}) (FileRef, error) {
	if mapped, ok := asMap(value); ok {
		ref := FileRef{Extensions: map[string]interface{}{}}
		source := toString(mapped["source"])
		if source == "" {
			source = toString(mapped["name"])
		}
		ref.Source = source
		ref.Key = toString(mapped["x-kubernetes-key"])
		ref.Target = toString(mapped["target"])
		ref.UID = toString(mapped["uid"])
		ref.GID = toString(mapped["gid"])
		ref.Mode = toString(mapped["mode"])
		ref.ReadOnly = toBool(mapped["read_only"])
		ref.Optional = boolPtrFromInterface(mapped["x-kubernetes-optional"])
		for key, val := range mapped {
			ref.Extensions[key] = val
		}
		if len(ref.Extensions) == 0 {
			ref.Extensions = nil
		}
		return ref, nil
	}
	source := toString(value)
	if source == "" {
		return FileRef{}, fmt.Errorf("file reference missing source")
	}
	return FileRef{Source: source}, nil
}

func parseVolumeMountMap(data map[string]interface{}) VolumeMount {
	mount := VolumeMount{
		Source:           toString(data["source"]),
		Target:           toString(data["target"]),
		Type:             toString(data["type"]),
		ReadOnly:         toBool(data["read_only"]),
		Consistency:      toString(data["consistency"]),
		Extensions:       map[string]interface{}{},
		BindExtensions:   map[string]interface{}{},
		VolumeExtensions: map[string]interface{}{},
		TmpfsExtensions:  map[string]interface{}{},
		ImageExtensions:  map[string]interface{}{},
	}
	if bind, ok := asMap(data["bind"]); ok {
		mount.Mode = toString(bind["selinux"])
		mount.Propagation = toString(bind["propagation"])
		if createHostPath, ok := bind["create_host_path"]; ok {
			v := toBool(createHostPath)
			mount.CreateHostPath = &v
		}
		if recursive := toString(bind["recursive"]); recursive != "" {
			mount.Options = ensureStringMap(mount.Options)
			mount.Options["recursive"] = recursive
		}
		for key, val := range bind {
			mount.BindExtensions[key] = val
		}
	}
	if volume, ok := asMap(data["volume"]); ok {
		mount.NoCopy = toBool(volume["nocopy"])
		if labels, err := toStringMap(volume["labels"]); err == nil && len(labels) > 0 {
			mount.VolumeLabels = labels
		}
		if subpath := toString(volume["subpath"]); subpath != "" {
			mount.Options = ensureStringMap(mount.Options)
			mount.Options["subpath"] = subpath
		}
		for key, val := range volume {
			mount.VolumeExtensions[key] = val
		}
	}
	if subPath := toString(data["x-kubernetes-subPath"]); subPath != "" {
		mount.SubPath = subPath
	}
	if subPathExpr := toString(data["x-kubernetes-subPathExpr"]); subPathExpr != "" {
		mount.SubPathExpr = strings.ReplaceAll(subPathExpr, "$$", "$")
	}
	if propagation, ok := data["x-kubernetes-mountPropagation"]; ok {
		mount.MountPropagation = toString(propagation)
	} else if propagation, ok := data["x-kubernetes-mount-propagation"]; ok {
		mount.MountPropagation = toString(propagation)
	}
	if recursiveReadOnly, ok := data["x-kubernetes-recursiveReadOnly"]; ok {
		mount.RecursiveReadOnly = toString(recursiveReadOnly)
	} else if recursiveReadOnly, ok := data["x-kubernetes-recursive-read-only"]; ok {
		mount.RecursiveReadOnly = toString(recursiveReadOnly)
	}
	if tmpfs, ok := asMap(data["tmpfs"]); ok {
		if size := toString(tmpfs["size"]); size != "" {
			mount.Options = ensureStringMap(mount.Options)
			mount.Options["size"] = size
		}
		if mode := toString(tmpfs["mode"]); mode != "" {
			mount.TmpfsMode = mode
		}
		for key, val := range tmpfs {
			mount.TmpfsExtensions[key] = val
		}
	}
	if image, ok := asMap(data["image"]); ok {
		if subpath := toString(image["subpath"]); subpath != "" {
			mount.ImageSubpath = subpath
		}
		for key, val := range image {
			mount.ImageExtensions[key] = val
		}
	}
	for key, val := range data {
		mount.Extensions[key] = val
	}
	return mount
}

func parseVolumeMount(volStr string) (VolumeMount, error) {
	// Handle formats like:
	// - "/host/path:/container/path"
	// - "/host/path:/container/path:ro"
	// - "volume_name:/container/path"
	// - "volume_name:/container/path:Z"

	parts := strings.Split(volStr, ":")
	if len(parts) < 2 {
		return VolumeMount{}, fmt.Errorf("invalid volume format: %s", volStr)
	}

	volume := VolumeMount{
		Source: parts[0],
		Target: parts[1],
	}

	if len(parts) > 2 {
		options := strings.Split(parts[2], ",")
		for _, option := range options {
			switch option {
			case "ro", "readonly":
				volume.ReadOnly = true
			case "rw":
				volume.ReadOnly = false
			case "Z", "z":
				volume.Mode = option
			case "rprivate", "private", "rshared", "shared", "rslave", "slave":
				volume.Propagation = option
			}
		}
	}

	// Determine type
	if strings.Contains(volume.Source, "/") || strings.HasPrefix(volume.Source, ".") || strings.HasPrefix(volume.Source, "~") {
		volume.Type = "bind"
	} else {
		volume.Type = "volume"
	}

	return volume, nil
}

func parseEnvironment(value interface{}) (map[string]string, error) {
	env := make(map[string]string)

	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			env[key] = toString(val)
		}
	case []interface{}:
		for _, envStr := range v {
			parts := strings.SplitN(toString(envStr), "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
	case []string:
		for _, envStr := range v {
			parts := strings.SplitN(envStr, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
	}

	return env, nil
}

func parseEnvFileRefs(value interface{}) ([]EnvFileRef, error) {
	var refs []EnvFileRef
	switch v := value.(type) {
	case string:
		refs = append(refs, EnvFileRef{Path: v})
	case []interface{}:
		for _, item := range v {
			ref, err := parseEnvFileRef(item)
			if err != nil {
				return nil, err
			}
			refs = append(refs, ref)
		}
	case []string:
		for _, item := range v {
			refs = append(refs, EnvFileRef{Path: item})
		}
	}
	return refs, nil
}

func parseEnvFileRef(value interface{}) (EnvFileRef, error) {
	if data, ok := asMap(value); ok {
		ref := EnvFileRef{
			Path:       toString(data["path"]),
			Format:     toString(data["format"]),
			Extensions: map[string]interface{}{},
		}
		if required, ok := data["required"]; ok {
			requiredValue := toBool(required)
			ref.Required = &requiredValue
		}
		for key, val := range data {
			ref.Extensions[key] = val
		}
		if len(ref.Extensions) == 0 {
			ref.Extensions = nil
		}
		if ref.Path == "" {
			return EnvFileRef{}, fmt.Errorf("env_file entry requires path")
		}
		return ref, nil
	}
	path := toString(value)
	if path == "" {
		return EnvFileRef{}, fmt.Errorf("env_file entry requires path")
	}
	return EnvFileRef{Path: path}, nil
}

func envFilePaths(refs []EnvFileRef) []string {
	paths := make([]string, 0, len(refs))
	for _, ref := range refs {
		if ref.Path != "" {
			paths = append(paths, ref.Path)
		}
	}
	return paths
}

func parseEnvSources(value interface{}) ([]EnvSource, error) {
	var sources []EnvSource
	items, ok := value.([]interface{})
	if !ok {
		return sources, nil
	}
	for _, item := range items {
		data, ok := asMap(item)
		if !ok {
			return nil, fmt.Errorf("env source must be a map")
		}
		source := EnvSource{
			Name:       toString(data["name"]),
			SourceType: toString(data["source_type"]),
			Source:     toString(data["source"]),
			Key:        toString(data["key"]),
			Optional:   toBool(data["optional"]),
			Extensions: map[string]interface{}{},
		}
		if source.Name == "" || source.SourceType == "" || source.Source == "" {
			return nil, fmt.Errorf("env source requires name, source_type, and source")
		}
		for key, value := range data {
			switch key {
			case "name", "source_type", "source", "key", "optional":
				continue
			default:
				source.Extensions[key] = deepCopyValue(value)
			}
		}
		if len(source.Extensions) == 0 {
			source.Extensions = nil
		}
		sources = append(sources, source)
	}
	return sources, nil
}

func parseEnvFromSources(value interface{}) ([]EnvFromSource, error) {
	var sources []EnvFromSource
	items, ok := value.([]interface{})
	if !ok {
		return sources, nil
	}
	for _, item := range items {
		data, ok := asMap(item)
		if !ok {
			return nil, fmt.Errorf("envFrom source must be a map")
		}
		source := EnvFromSource{
			SourceType: toString(data["source_type"]),
			Source:     toString(data["source"]),
			Prefix:     toString(data["prefix"]),
			Optional:   toBool(data["optional"]),
			Extensions: map[string]interface{}{},
		}
		if source.SourceType == "" || source.Source == "" {
			return nil, fmt.Errorf("envFrom source requires source_type and source")
		}
		for key, value := range data {
			switch key {
			case "source_type", "source", "prefix", "optional":
				continue
			default:
				source.Extensions[key] = deepCopyValue(value)
			}
		}
		if len(source.Extensions) == 0 {
			source.Extensions = nil
		}
		sources = append(sources, source)
	}
	return sources, nil
}

func parseHealthCheck(value interface{}) (*HealthCheck, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("healthcheck must be a map")
	}

	hc := &HealthCheck{Extensions: map[string]interface{}{}}

	for key, val := range data {
		switch key {
		case "test":
			test, err := toStringSlice(val)
			if err != nil {
				return nil, fmt.Errorf("failed to parse test: %w", err)
			}
			hc.Test = test
		case "interval":
			hc.Interval = toString(val)
		case "timeout":
			hc.Timeout = toString(val)
		case "retries":
			hc.Retries = toInt(val)
		case "start_period":
			hc.StartPeriod = toString(val)
		case "start_interval":
			hc.StartInterval = toString(val)
		case "disable":
			hc.Disable = toBool(val)
			hc.DisableSet = true
		default:
			hc.Extensions[key] = val
		}
	}
	if len(hc.Extensions) == 0 {
		hc.Extensions = nil
	}

	return normalizeHealthCheck(hc), nil
}

func parseNetwork(name string, data map[string]interface{}) (*Network, error) {
	network := &Network{Name: name, Extensions: map[string]interface{}{}}

	for key, value := range data {
		switch key {
		case "name":
			network.PlatformName = toString(value)
		case "driver":
			network.Driver = toString(value)
		case "driver_opts":
			opts, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver_opts: %w", err)
			}
			network.DriverOpts = opts
		case "attachable":
			network.Attachable = toBool(value)
			network.AttachableSet = true
		case "external":
			network.External, network.PlatformName, network.ExternalExtensions = parseComposeExternal(value, network.PlatformName)
			network.ExternalSet = true
		case "internal":
			network.Internal = toBool(value)
			network.InternalSet = true
		case "enable_ipv4":
			enable := toBool(value)
			network.EnableIPv4 = &enable
		case "enable_ipv6":
			enable := toBool(value)
			network.EnableIPv6 = &enable
		case "ipam":
			ipam, err := parseIPAM(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ipam: %w", err)
			}
			network.IPAM = ipam
		case "labels":
			labels, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
			network.Labels = labels
		default:
			network.Extensions[key] = value
		}
	}

	return network, nil
}

func parseVolume(name string, data map[string]interface{}) (*Volume, error) {
	volume := &Volume{Name: name, Extensions: map[string]interface{}{}}

	for key, value := range data {
		switch key {
		case "name":
			volume.PlatformName = toString(value)
		case "driver":
			volume.Driver = toString(value)
		case "driver_opts":
			opts, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver_opts: %w", err)
			}
			volume.DriverOpts = opts
		case "x-bolabaden-volume-driver":
			volume.Driver = toString(value)
		case "x-bolabaden-volume-driver-opts":
			opts, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse x-bolabaden-volume-driver-opts: %w", err)
			}
			volume.DriverOpts = opts
		case "external":
			volume.External, volume.PlatformName, volume.ExternalExtensions = parseComposeExternal(value, volume.PlatformName)
			volume.ExternalSet = true
		case "labels":
			labels, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
			volume.Labels = labels
		default:
			volume.Extensions[key] = value
			if key == "x-kubernetes-kind" && volume.Extensions["kubernetes.kind"] == nil {
				volume.Extensions["kubernetes.kind"] = value
			}
		}
	}

	return volume, nil
}

func parseConfig(name string, data map[string]interface{}) (*Config, error) {
	config := &Config{Name: name, Extensions: map[string]interface{}{}}

	for key, value := range data {
		switch key {
		case "name":
			config.PlatformName = toString(value)
		case "content":
			config.Content = toString(value)
		case "environment":
			config.Environment = toString(value)
		case "file":
			config.File = toString(value)
		case "template", "template_driver":
			config.Template = toString(value)
		case "mode":
			config.Mode = toString(value)
		case "external":
			config.External, config.PlatformName, config.ExternalExtensions = parseComposeExternal(value, config.PlatformName)
			config.ExternalSet = true
		case "x-bolabaden-config-content":
			config.Content = toString(value)
		case "x-bolabaden-config-environment":
			config.Environment = toString(value)
		case "x-bolabaden-config-file":
			config.File = toString(value)
		case "x-bolabaden-config-template":
			config.Template = toString(value)
		case "x-bolabaden-config-mode":
			config.Mode = toString(value)
		case "labels":
			labels, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
			config.Labels = labels
		default:
			config.Extensions[key] = value
			switch key {
			case "x-kubernetes-labels":
				if config.Extensions["kubernetes.labels"] == nil {
					config.Extensions["kubernetes.labels"] = value
				}
			case "x-kubernetes-annotations":
				if config.Extensions["kubernetes.annotations"] == nil {
					config.Extensions["kubernetes.annotations"] = value
				}
			case "x-kubernetes-immutable":
				if config.Extensions["kubernetes.immutable"] == nil {
					config.Extensions["kubernetes.immutable"] = value
				}
			case "x-kubernetes-data":
				if config.Extensions["kubernetes.data"] == nil {
					config.Extensions["kubernetes.data"] = value
				}
			case "x-kubernetes-binaryData", "x-kubernetes-binary-data":
				if config.Extensions["kubernetes.binaryData"] == nil {
					config.Extensions["kubernetes.binaryData"] = value
				}
			}
		}
	}

	return config, nil
}

func parseSecret(name string, data map[string]interface{}) (*Secret, error) {
	secret := &Secret{Name: name, Extensions: map[string]interface{}{}}

	for key, value := range data {
		switch key {
		case "name":
			secret.PlatformName = toString(value)
		case "file":
			secret.File = toString(value)
		case "environment":
			secret.Environment = toString(value)
		case "template", "template_driver":
			secret.Template = toString(value)
		case "external":
			secret.External, secret.PlatformName, secret.ExternalExtensions = parseComposeExternal(value, secret.PlatformName)
			secret.ExternalSet = true
		case "x-bolabaden-secret-file":
			secret.File = toString(value)
		case "x-bolabaden-secret-environment":
			secret.Environment = toString(value)
		case "x-bolabaden-secret-template":
			secret.Template = toString(value)
		case "x-bolabaden-secret-driver":
			secret.Driver = toString(value)
		case "x-bolabaden-secret-driver-opts":
			opts, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse x-bolabaden-secret-driver-opts: %w", err)
			}
			secret.DriverOpts = opts
		case "driver":
			secret.Driver = toString(value)
		case "driver_opts":
			opts, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver_opts: %w", err)
			}
			secret.DriverOpts = opts
		case "labels":
			labels, err := toStringMap(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
			secret.Labels = labels
		default:
			secret.Extensions[key] = value
			switch key {
			case "x-kubernetes-labels":
				if secret.Extensions["kubernetes.labels"] == nil {
					secret.Extensions["kubernetes.labels"] = value
				}
			case "x-kubernetes-annotations":
				if secret.Extensions["kubernetes.annotations"] == nil {
					secret.Extensions["kubernetes.annotations"] = value
				}
			case "x-kubernetes-immutable":
				if secret.Extensions["kubernetes.immutable"] == nil {
					secret.Extensions["kubernetes.immutable"] = value
				}
			case "x-kubernetes-data":
				if secret.Extensions["kubernetes.data"] == nil {
					secret.Extensions["kubernetes.data"] = value
				}
			case "x-kubernetes-stringData", "x-kubernetes-string-data":
				if secret.Extensions["kubernetes.stringData"] == nil {
					secret.Extensions["kubernetes.stringData"] = value
				}
			case "x-kubernetes-type":
				if secret.Extensions["kubernetes.type"] == nil {
					secret.Extensions["kubernetes.type"] = value
				}
			}
		}
	}

	return secret, nil
}

func parseIPAM(value interface{}) (*IPAMConfig, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("ipam must be a map")
	}

	ipam := &IPAMConfig{Extensions: map[string]interface{}{}}

	for key, val := range data {
		switch key {
		case "driver":
			ipam.Driver = toString(val)
		case "config":
			config, err := parseIPAMSubnets(val)
			if err != nil {
				return nil, fmt.Errorf("failed to parse config: %w", err)
			}
			ipam.Config = config
		case "options":
			options, err := toStringMap(val)
			if err != nil {
				return nil, fmt.Errorf("failed to parse options: %w", err)
			}
			ipam.Options = options
		default:
			ipam.Extensions[key] = val
		}
	}

	return ipam, nil
}

func parseIPAMSubnets(value interface{}) ([]IPAMSubnet, error) {
	var subnets []IPAMSubnet

	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			if data, ok := asMap(item); ok {
				subnet := IPAMSubnet{
					AuxAddresses: map[string]string{},
					Extensions:   map[string]interface{}{},
				}
				for key, val := range data {
					switch key {
					case "subnet":
						subnet.Subnet = toString(val)
					case "gateway":
						subnet.Gateway = toString(val)
					case "ip_range":
						subnet.IPRange = toString(val)
					case "aux_addresses":
						aux, err := toStringMap(val)
						if err != nil {
							return nil, fmt.Errorf("failed to parse aux_addresses: %w", err)
						}
						subnet.AuxAddresses = aux
					default:
						subnet.Extensions[key] = val
					}
				}
				subnets = append(subnets, subnet)
			}
		}
	}

	return subnets, nil
}

// Helper functions
func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

func toBool(value interface{}) bool {
	if b, ok := value.(bool); ok {
		return b
	}
	if str := toString(value); str != "" {
		if str == "true" || str == "1" || str == "yes" {
			return true
		}
	}
	return false
}

func parseComposeExternal(value interface{}, currentName string) (bool, string, map[string]interface{}) {
	if data, ok := asMap(value); ok {
		extensions := map[string]interface{}{}
		for key, val := range data {
			extensions[key] = deepCopyValue(val)
		}
		if name := toString(data["name"]); name != "" {
			return true, name, extensions
		}
		return true, currentName, extensions
	}
	if name := toString(value); name != "" && name != "true" && name != "false" && name != "1" && name != "0" && name != "yes" && name != "no" {
		return true, name, nil
	}
	return toBool(value), currentName, nil
}

func toStringSlice(value interface{}) ([]string, error) {
	var result []string

	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			result = append(result, toString(item))
		}
	case []string:
		result = v
	case string:
		result = []string{v}
	default:
		return nil, fmt.Errorf("cannot convert to string slice: %T", value)
	}

	return result, nil
}

func toInt64Slice(value interface{}) ([]int64, error) {
	var result []int64

	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			result = append(result, int64(toInt(item)))
		}
	case []int64:
		result = append([]int64{}, v...)
	case []int:
		for _, item := range v {
			result = append(result, int64(item))
		}
	case []float64:
		for _, item := range v {
			result = append(result, int64(item))
		}
	default:
		return nil, fmt.Errorf("cannot convert to int64 slice: %T", value)
	}

	return result, nil
}

func serializeKubernetesTolerations(tolerations []Toleration) []map[string]interface{} {
	if len(tolerations) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(tolerations))
	for _, toleration := range tolerations {
		item := map[string]interface{}{}
		if toleration.Key != "" {
			item["key"] = toleration.Key
		}
		if toleration.Operator != "" {
			item["operator"] = toleration.Operator
		}
		if toleration.Value != "" {
			item["value"] = toleration.Value
		}
		if toleration.Effect != "" {
			item["effect"] = toleration.Effect
		}
		if toleration.TolerationSeconds != nil {
			item["tolerationSeconds"] = *toleration.TolerationSeconds
		}
		for key, value := range toleration.Extensions {
			item[key] = deepCopyValue(value)
		}
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func serializeKubernetesTolerationsNative(tolerations []Toleration) []map[string]interface{} {
	if len(tolerations) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(tolerations))
	for _, toleration := range tolerations {
		item := map[string]interface{}{}
		if toleration.Key != "" {
			item["key"] = toleration.Key
		}
		if toleration.Operator != "" {
			item["operator"] = toleration.Operator
		}
		if toleration.Value != "" {
			item["value"] = toleration.Value
		}
		if toleration.Effect != "" {
			item["effect"] = toleration.Effect
		}
		if toleration.TolerationSeconds != nil {
			item["tolerationSeconds"] = *toleration.TolerationSeconds
		}
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func parseKubernetesTolerationsExtension(value interface{}) []Toleration {
	items, ok := value.([]interface{})
	if !ok {
		return nil
	}
	var tolerations []Toleration
	for _, item := range items {
		tolerationMap, ok := asMap(item)
		if !ok {
			continue
		}
		toleration := Toleration{
			Key:        toString(tolerationMap["key"]),
			Operator:   toString(tolerationMap["operator"]),
			Value:      toString(tolerationMap["value"]),
			Effect:     toString(tolerationMap["effect"]),
			Extensions: map[string]interface{}{},
		}
		switch v := tolerationMap["tolerationSeconds"].(type) {
		case int64:
			toleration.TolerationSeconds = &v
		case int:
			value := int64(v)
			toleration.TolerationSeconds = &value
		case float64:
			value := int64(v)
			toleration.TolerationSeconds = &value
		}
		for key, value := range tolerationMap {
			switch key {
			case "key", "operator", "value", "effect", "tolerationSeconds":
			default:
				toleration.Extensions[key] = deepCopyValue(value)
			}
		}
		if len(toleration.Extensions) == 0 {
			toleration.Extensions = nil
		}
		if toleration.Key != "" || toleration.Operator != "" || toleration.Value != "" || toleration.Effect != "" || toleration.TolerationSeconds != nil {
			tolerations = append(tolerations, toleration)
		}
	}
	if len(tolerations) == 0 {
		return nil
	}
	return tolerations
}

func toStringMap(value interface{}) (map[string]string, error) {
	result := make(map[string]string)

	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			result[key] = toString(val)
		}
	case map[string]string:
		result = v
	case []interface{}:
		for _, item := range v {
			parts := strings.SplitN(toString(item), "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	case []string:
		for _, item := range v {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	default:
		return nil, fmt.Errorf("cannot convert to string map: %T", value)
	}

	return result, nil
}

func toStringPtrMap(value interface{}) (map[string]*string, error) {
	result := make(map[string]*string)

	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if val == nil {
				result[key] = nil
				continue
			}
			str := toString(val)
			result[key] = &str
		}
	case map[string]string:
		for key, val := range v {
			str := val
			result[key] = &str
		}
	case map[interface{}]interface{}:
		for key, val := range v {
			strKey := toString(key)
			if val == nil {
				result[strKey] = nil
				continue
			}
			str := toString(val)
			result[strKey] = &str
		}
	case []interface{}:
		for _, item := range v {
			parts := strings.SplitN(toString(item), "=", 2)
			if len(parts) == 1 {
				result[parts[0]] = nil
				continue
			}
			str := parts[1]
			result[parts[0]] = &str
		}
	case []string:
		for _, item := range v {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 1 {
				result[parts[0]] = nil
				continue
			}
			str := parts[1]
			result[parts[0]] = &str
		}
	default:
		return nil, fmt.Errorf("cannot convert to string pointer map: %T", value)
	}

	return result, nil
}

func asMap(value interface{}) (map[string]interface{}, bool) {
	switch v := value.(type) {
	case map[string]interface{}:
		return v, true
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			result[toString(key)] = val
		}
		return result, true
	case nil:
		return nil, false
	default:
		return nil, false
	}
}

func toInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		return parseInt(v)
	default:
		return 0
	}
}

func isStandardKey(key string) bool {
	standardKeys := []string{
		"name", "version", "services", "networks", "volumes", "configs", "secrets", "models", "include",
	}
	for _, k := range standardKeys {
		if k == key {
			return true
		}
	}
	return false
}

// SerializeDockerCompose converts an Application to Docker Compose YAML
func SerializeDockerCompose(app *Application) (string, error) {
	emitApp := cloneApplication(app)
	syncPortableApplicationState(emitApp)
	data := make(map[string]interface{})

	// Version
	if emitApp.Version != "" {
		data["version"] = emitApp.Version
	}
	if emitApp.Name != "" {
		data["name"] = emitApp.Name
	}

	// Includes
	if len(emitApp.IncludeEntries) > 0 {
		data["include"] = cloneInterfaceSlice(emitApp.IncludeEntries)
	} else if len(emitApp.Includes) > 0 {
		data["include"] = emitApp.Includes
	}

	// Services
	if len(emitApp.Services) > 0 {
		servicesData := make(map[string]interface{})
		for name, service := range emitApp.Services {
			serviceData, err := serializeService(emitApp, service)
			if err != nil {
				return "", fmt.Errorf("failed to serialize service %s: %w", name, err)
			}
			servicesData[name] = serviceData
		}
		applyCanonicalRoutesToComposeServices(emitApp, servicesData)
		data["services"] = servicesData
	}

	// Models
	if models := applicationModelsForEmit(emitApp); len(models) > 0 {
		modelsData := make(map[string]interface{})
		for name, model := range models {
			modelData, err := serializeComposeTopLevelModel(model)
			if err != nil {
				return "", fmt.Errorf("failed to serialize model %s: %w", name, err)
			}
			modelsData[name] = modelData
		}
		data["models"] = modelsData
	}

	// Networks
	if len(emitApp.Networks) > 0 {
		networksData := make(map[string]interface{})
		for name, network := range emitApp.Networks {
			networkData, err := serializeNetwork(network)
			if err != nil {
				return "", fmt.Errorf("failed to serialize network %s: %w", name, err)
			}
			networksData[name] = networkData
		}
		data["networks"] = networksData
	}

	// Volumes
	if len(emitApp.Volumes) > 0 {
		volumesData := make(map[string]interface{})
		for name, volume := range emitApp.Volumes {
			volumeData, err := serializeVolume(volume)
			if err != nil {
				return "", fmt.Errorf("failed to serialize volume %s: %w", name, err)
			}
			volumesData[name] = volumeData
		}
		data["volumes"] = volumesData
	}

	// Configs
	if len(emitApp.Configs) > 0 {
		configsData := make(map[string]interface{})
		for name, config := range emitApp.Configs {
			configData, err := serializeConfig(config)
			if err != nil {
				return "", fmt.Errorf("failed to serialize config %s: %w", name, err)
			}
			configsData[name] = configData
		}
		data["configs"] = configsData
	}

	// Secrets
	if len(emitApp.Secrets) > 0 {
		secretsData := make(map[string]interface{})
		for name, secret := range emitApp.Secrets {
			secretData, err := serializeSecret(secret)
			if err != nil {
				return "", fmt.Errorf("failed to serialize secret %s: %w", name, err)
			}
			secretsData[name] = secretData
		}
		data["secrets"] = secretsData
	}

	// Extensions
	for key, value := range convertedExtensions(emitApp.Extensions, emitApp.Platform, PlatformDockerCompose) {
		if key == composeRawYAMLExtension {
			continue
		}
		if key == "name" {
			continue
		}
		composeKey := composeApplicationExtensionKey(key)
		if composeKey == composeKubernetesNamespaceExtension && data[composeKey] != nil {
			continue
		}
		data[composeKey] = value
	}
	if routes := canonicalRoutesForApplication(emitApp); len(routes) > 0 {
		data[composeAppRoutesExtension] = sortedRouteSpecList(routes)
	}
	if policies := canonicalPoliciesForApplication(emitApp); len(policies) > 0 {
		data[composeAppPoliciesExtension] = sortedPolicySpecList(policies)
	}
	targetPlatform := PlatformDockerCompose
	if emitApp.Platform == PlatformDockerSwarm {
		targetPlatform = PlatformDockerSwarm
	}
	if resources := canonicalRawResourcesForBridge(emitApp, targetPlatform); len(resources) > 0 {
		data[composeCanonicalRawResourcesExtension] = resources
	}
	if resources := canonicalKubernetesRawResourcesForApplication(emitApp); len(resources) > 0 {
		data[composeKubernetesRawResourcesExtension] = resources
	}
	if emitApp.Platform == PlatformDockerSwarm {
		data["x-platform"] = string(PlatformDockerSwarm)
	} else {
		delete(data, "x-platform")
	}

	data = escapeComposeExtensionValues(data, false).(map[string]interface{})

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	return string(yamlData), nil
}

func escapeComposeExtensionValues(value interface{}, inExtension bool) interface{} {
	if value == nil {
		return nil
	}
	switch typed := value.(type) {
	case string:
		if inExtension {
			return strings.ReplaceAll(typed, "$", "$$")
		}
		return typed
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return value
		}
		cloned := make(map[string]interface{}, rv.Len())
		for _, key := range rv.MapKeys() {
			keyStr := key.String()
			cloned[keyStr] = escapeComposeExtensionValues(rv.MapIndex(key).Interface(), inExtension || strings.HasPrefix(keyStr, "x-"))
		}
		return cloned
	case reflect.Slice, reflect.Array:
		cloned := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			cloned[i] = escapeComposeExtensionValues(rv.Index(i).Interface(), inExtension)
		}
		return cloned
	default:
		return value
	}
}

func rehydrateComposeApplicationExtensions(app *Application) {
	if app == nil || app.Extensions == nil {
		return
	}
	if strings.EqualFold(toString(app.Extensions["x-platform"]), string(PlatformDockerSwarm)) {
		app.Platform = PlatformDockerSwarm
		for _, service := range app.Services {
			if service != nil {
				service.Platform = PlatformDockerSwarm
			}
		}
	}
	if app.Platform != PlatformDockerSwarm && composeAppLooksLikeSwarm(app) {
		app.Platform = PlatformDockerSwarm
		for _, service := range app.Services {
			if service != nil {
				service.Platform = PlatformDockerSwarm
			}
		}
	}
	namespace := ""
	if value, ok := applicationExtensionValueForKey(app, "kubernetes.namespace"); ok {
		namespace = toString(value)
	}
	if namespace == "" {
		namespace = toString(app.Extensions[composeKubernetesNamespaceExtension])
	}
	if namespace != "" {
		app.Namespace = namespace
		app.Extensions["kubernetes.namespace"] = namespace
		app.Extensions[composeKubernetesNamespaceExtension] = namespace
	}
	rehydratePortableMesh(app)
	for key, value := range copyStringInterfaceMap(app.Extensions) {
		if canonical := composeApplicationCanonicalKey(key); canonical != "" {
			app.Extensions[canonical] = value
		}
	}
}

func composeAppLooksLikeSwarm(app *Application) bool {
	if app == nil {
		return false
	}
	if strings.EqualFold(toString(app.Extensions["x-platform"]), string(PlatformDockerSwarm)) {
		return true
	}
	for _, service := range app.Services {
		if service == nil {
			continue
		}
		if service.Deploy != nil && (strings.EqualFold(service.Deploy.Mode, "global") || isSwarmJobMode(service.Deploy.Mode)) {
			return true
		}
		if service.Deploy != nil && !isEmptySwarmJobSpec(service.Deploy.Job) {
			return true
		}
		if _, ok := firstExtensionValue(service.Extensions, composeSwarmJobExtension, "x-swarm-job-spec"); ok {
			return true
		}
	}
	return false
}

func composeApplicationExtensionKey(key string) string {
	switch key {
	case "x-kubernetes-nodeSelector":
		key = "x-kubernetes-node-selector"
	case "x-kubernetes-resourceClaims":
		key = "x-kubernetes-resource-claims"
	}
	if key == "kubernetes.namespace" {
		return composeKubernetesNamespaceExtension
	}
	if key == "kubernetes.serviceResources" {
		return "x-kubernetes-services"
	}
	if key == "kubernetes.horizontalPodAutoscalers" || key == "kubernetes.hpas" {
		return "x-kubernetes-horizontal-pod-autoscalers"
	}
	if key == "kubernetes.podDisruptionBudgets" || key == "kubernetes.pdbs" {
		return "x-kubernetes-pod-disruption-budgets"
	}
	if key == "kubernetes.serviceTargets" {
		return "x-kubernetes-service-targets"
	}
	if key == "kubernetes.servicePortTargets" {
		return "x-kubernetes-service-port-targets"
	}
	if key == "kubernetes.nodeSelector" {
		return "x-kubernetes-node-selector"
	}
	if key == "kubernetes.resourceClaims" {
		return "x-kubernetes-resource-claims"
	}
	if key == "kubernetes.serviceAccounts" {
		return "x-kubernetes-service-accounts"
	}
	if key == "kubernetes.resourceQuotas" {
		return "x-kubernetes-resource-quotas"
	}
	if key == "kubernetes.limitRanges" {
		return "x-kubernetes-limit-ranges"
	}
	if key == "kubernetes.priorityClasses" {
		return "x-kubernetes-priority-classes"
	}
	if key == "kubernetes.runtimeClasses" {
		return "x-kubernetes-runtime-classes"
	}
	if key == "kubernetes.storageClasses" {
		return "x-kubernetes-storage-classes"
	}
	if key == "kubernetes.ingressClasses" {
		return "x-kubernetes-ingress-classes"
	}
	if key == "kubernetes.mutatingWebhookConfigurations" {
		return "x-kubernetes-mutating-webhook-configurations"
	}
	if key == "kubernetes.validatingWebhookConfigurations" {
		return "x-kubernetes-validating-webhook-configurations"
	}
	if key == "kubernetes.customResourceDefinitions" {
		return "x-kubernetes-custom-resource-definitions"
	}
	if key == "kubernetes.customResources" {
		return "x-kubernetes-custom-resources"
	}
	if key == kubernetesNamespacesExtensionKey {
		return composeKubernetesNamespacesExtensionKey
	}
	if key == kubernetesWorkloadsExtensionKey {
		return composeKubernetesWorkloadsExtensionKey
	}
	if key == kubernetesConfigMapsExtensionKey {
		return composeKubernetesConfigMapsExtensionKey
	}
	if key == kubernetesSecretsExtensionKey {
		return composeKubernetesSecretsExtensionKey
	}
	if key == kubernetesSourceDocumentsExtensionKey {
		return "x-kubernetes-source-documents"
	}
	if strings.HasPrefix(key, "x-") {
		return key
	}
	return "x-" + strings.ReplaceAll(key, ".", "-")
}

func normalizeComposeExtensionObject(input map[string]interface{}, allowedKeys []string) map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	result := map[string]interface{}{}
	allowed := map[string]struct{}{}
	for _, key := range allowedKeys {
		allowed[key] = struct{}{}
	}
	for key, value := range input {
		if _, ok := allowed[key]; ok {
			result[key] = deepCopyValue(value)
			continue
		}
		result[composeApplicationExtensionKey(key)] = deepCopyValue(value)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func serializeComposeBlkioConfig(input map[string]interface{}) (map[string]interface{}, map[string]interface{}) {
	if len(input) == 0 {
		return nil, nil
	}
	result := map[string]interface{}{}
	extensions := map[string]interface{}{}
	for key, value := range input {
		switch key {
		case "weight":
			result[key] = deepCopyValue(value)
		case "weight_device":
			official, extra := serializeComposeBlkioDeviceList(value, []string{"path", "weight"})
			if len(official) > 0 {
				result[key] = official
			}
			if len(extra) > 0 {
				extensions[key] = extra
			}
		case "device_read_bps", "device_read_iops", "device_write_bps", "device_write_iops":
			official, extra := serializeComposeBlkioDeviceList(value, []string{"path", "rate"})
			if len(official) > 0 {
				result[key] = official
			}
			if len(extra) > 0 {
				extensions[key] = extra
			}
		default:
			extensions[key] = deepCopyValue(value)
		}
	}
	if len(result) == 0 {
		result = nil
	}
	if len(extensions) == 0 {
		extensions = nil
	}
	return result, extensions
}

func serializeComposeBlkioDeviceList(value interface{}, allowedKeys []string) ([]interface{}, []interface{}) {
	items, ok := value.([]interface{})
	if !ok || len(items) == 0 {
		return nil, nil
	}
	official := make([]interface{}, 0, len(items))
	extensions := make([]interface{}, 0, len(items))
	allowed := map[string]struct{}{}
	for _, key := range allowedKeys {
		allowed[key] = struct{}{}
	}
	for _, item := range items {
		mapped, ok := asMap(item)
		if !ok || len(mapped) == 0 {
			official = append(official, deepCopyValue(item))
			extensions = append(extensions, map[string]interface{}{})
			continue
		}
		officialItem := map[string]interface{}{}
		extraItem := map[string]interface{}{}
		for key, value := range mapped {
			if _, ok := allowed[key]; ok {
				officialItem[key] = deepCopyValue(value)
				continue
			}
			extraItem[key] = deepCopyValue(value)
		}
		if len(officialItem) > 0 {
			official = append(official, officialItem)
		}
		if len(extraItem) > 0 {
			extensions = append(extensions, extraItem)
		} else {
			extensions = append(extensions, map[string]interface{}{})
		}
	}
	if len(official) == 0 {
		official = nil
	}
	if len(extensions) == 0 {
		extensions = nil
	}
	return official, extensions
}

func parseComposeTopLevelModel(name string, value map[string]interface{}) *ComposeModel {
	if len(value) == 0 {
		return nil
	}
	model := &ComposeModel{
		Name:         name,
		Model:        toString(value["model"]),
		ContextSize:  toInt(value["context_size"]),
		RuntimeFlags: nil,
		Extensions:   map[string]interface{}{},
	}
	if flags, err := toStringSlice(value["runtime_flags"]); err == nil && len(flags) > 0 {
		model.RuntimeFlags = append([]string{}, flags...)
	}
	for key, raw := range value {
		switch key {
		case "model", "context_size", "runtime_flags", "name":
			continue
		default:
			model.Extensions[key] = raw
		}
	}
	if len(model.Extensions) == 0 {
		model.Extensions = nil
	}
	if model.Name == "" && model.Model == "" && model.ContextSize == 0 && len(model.RuntimeFlags) == 0 && len(model.Extensions) == 0 {
		return nil
	}
	return model
}

func serializeComposeTopLevelModel(model *ComposeModel) (map[string]interface{}, error) {
	if model == nil {
		return nil, fmt.Errorf("nil model")
	}
	data := map[string]interface{}{}
	if model.Model != "" {
		data["model"] = model.Model
	}
	if model.ContextSize > 0 {
		data["context_size"] = model.ContextSize
	}
	if len(model.RuntimeFlags) > 0 {
		data["runtime_flags"] = append([]string{}, model.RuntimeFlags...)
	}
	for key, value := range model.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty compose model")
	}
	return data, nil
}

func RestoreDockerComposeSource(app *Application, filename string) error {
	raw, err := dockerComposeSourceContent(app)
	if err != nil {
		return err
	}
	if dir := filepath.Dir(filename); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create Compose output directory: %w", err)
		}
	}
	if err := os.WriteFile(filename, []byte(raw), 0644); err != nil {
		return fmt.Errorf("failed to restore Compose source: %w", err)
	}
	return nil
}

func dockerComposeSourceContent(app *Application) (string, error) {
	raw := ""
	if app != nil && app.Extensions != nil {
		raw = toString(app.Extensions[composeRawYAMLExtension])
	}
	if strings.TrimSpace(raw) == "" {
		if preserved := canonicalRawResourceValue(app, PlatformDockerCompose, "ComposeYAML"); preserved != nil {
			raw = toString(preserved)
		}
	}
	if strings.TrimSpace(raw) == "" {
		if preserved := canonicalRawResourceValue(app, PlatformDockerSwarm, "ComposeYAML"); preserved != nil {
			raw = toString(preserved)
		}
	}
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("application does not contain raw Compose YAML")
	}
	return raw, nil
}

func serializeService(app *Application, service *Service) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// Basic fields
	if service.Image != "" {
		data["image"] = service.Image
	}
	if service.ContainerName != "" {
		data["container_name"] = service.ContainerName
	}
	if service.Hostname != "" {
		data["hostname"] = service.Hostname
	}

	// Ports
	if len(service.Ports) > 0 {
		var ports []interface{}
		for _, port := range service.Ports {
			if portMappingNeedsLongSyntax(port) {
				ports = append(ports, serializePortMappingLong(port))
				continue
			}
			portStr := port.ContainerPort
			if port.HostPort != "" {
				if port.HostIP != "" {
					portStr = fmt.Sprintf("%s:%s:%s", port.HostIP, port.HostPort, port.ContainerPort)
				} else {
					portStr = fmt.Sprintf("%s:%s", port.HostPort, port.ContainerPort)
				}
			}
			if port.Protocol != "" && port.Protocol != "tcp" {
				portStr = fmt.Sprintf("%s/%s", portStr, port.Protocol)
			}
			ports = append(ports, portStr)
		}
		data["ports"] = ports
	}

	// Expose
	if len(service.Expose) > 0 {
		data["expose"] = service.Expose
	}

	// Networks
	if len(service.Networks) > 0 {
		data["networks"] = serializeServiceNetworks(service)
	}
	if len(service.DNS) > 0 {
		data["dns"] = service.DNS
	}
	if len(service.DNSSearch) > 0 {
		data["dns_search"] = service.DNSSearch
	}
	if len(service.DNSOptions) > 0 {
		data["dns_opt"] = service.DNSOptions
	}
	if len(service.ExtraHosts) > 0 {
		data["extra_hosts"] = service.ExtraHosts
	}
	if aliases := kubernetesHostAliasesFromService(service); len(aliases) > 0 {
		data["x-kubernetes-hostAliases"] = aliases
	}

	// Environment
	if len(service.Environment) > 0 {
		data["environment"] = service.Environment
	}

	// Env file
	if len(service.EnvFileRefs) > 0 {
		data["env_file"] = serializeEnvFileRefs(service.EnvFileRefs)
	} else if len(service.EnvFile) > 0 {
		data["env_file"] = service.EnvFile
	}
	if len(service.EnvSources) > 0 {
		data["x-env-sources"] = serializeEnvSources(service.EnvSources)
	}
	if len(service.EnvFrom) > 0 {
		data["x-env-from"] = serializeEnvFromSources(service.EnvFrom)
	}
	if len(service.ImagePullSecrets) > 0 {
		data["x-kubernetes-imagePullSecrets"] = append([]string{}, service.ImagePullSecrets...)
	}
	if service.ImagePullPolicy != "" {
		data["x-kubernetes-imagePullPolicy"] = service.ImagePullPolicy
	}
	if service.DNSPolicy != "" {
		data["x-kubernetes-dnsPolicy"] = service.DNSPolicy
	}
	if service.SchedulerName != "" {
		data["x-kubernetes-schedulerName"] = service.SchedulerName
	}
	if service.Hostname != "" {
		data["x-kubernetes-hostname"] = service.Hostname
	}
	if service.TerminationMessagePath != "" {
		data["x-kubernetes-terminationMessagePath"] = service.TerminationMessagePath
	}
	if service.TerminationMessagePolicy != "" {
		data["x-kubernetes-terminationMessagePolicy"] = service.TerminationMessagePolicy
	}

	// Command
	if len(service.Command) > 0 {
		data["command"] = service.Command
	}

	// Entrypoint
	if len(service.Entrypoint) > 0 {
		data["entrypoint"] = service.Entrypoint
	}

	// Working dir
	if service.WorkingDir != "" {
		data["working_dir"] = service.WorkingDir
	}
	if build := buildConfigToCompose(service.Build); build != nil {
		data["build"] = build
		if buildConfigHasData(service.Build) {
			if raw, err := json.Marshal(service.Build); err == nil {
				data["x-bolabaden-portable-build"] = string(raw)
			}
		}
	}
	if develop := developConfigToCompose(service.Develop); develop != nil {
		data["develop"] = develop
	}
	if service.Lifecycle != nil {
		if hooks := serviceHooksToCompose(service.Lifecycle.PreStart); len(hooks) > 0 {
			data["pre_start"] = hooks
		}
		if hooks := serviceHooksToCompose(service.Lifecycle.PostStart); len(hooks) > 0 {
			data["post_start"] = hooks
		}
		if hooks := serviceHooksToCompose(service.Lifecycle.PreStop); len(hooks) > 0 {
			data["pre_stop"] = hooks
		}
	}
	if len(service.Devices) > 0 || len(service.DeviceMappings) > 0 {
		data["devices"] = deviceMappingsToCompose(service.Devices, service.DeviceMappings)
	}
	if service.PrivilegedSet {
		data["privileged"] = service.Privileged
	} else if service.Privileged {
		data["privileged"] = true
	}
	if service.User != "" {
		data["user"] = service.User
	}
	if service.Group != "" {
		data["x-kubernetes-group"] = service.Group
	} else if value := extensionStringValue(service, "kubernetes.group"); value != "" {
		data["x-kubernetes-group"] = value
	}
	if len(service.GroupAdd) > 0 {
		data["group_add"] = append([]string{}, service.GroupAdd...)
	}
	if len(service.Sysctls) > 0 {
		data["sysctls"] = copyStringMap(service.Sysctls)
	}
	if service.Runtime != "" {
		data["runtime"] = service.Runtime
	}
	if compat := service.ComposeCompat; compat != nil {
		if compat.Attach != nil {
			data["attach"] = *compat.Attach
		}
		if len(compat.Annotations) > 0 {
			data["annotations"] = copyStringMap(compat.Annotations)
		}
		if len(compat.BlkioConfig) > 0 {
			blkioConfig, blkioExtensions := serializeComposeBlkioConfig(compat.BlkioConfig)
			if len(blkioConfig) > 0 {
				data["blkio_config"] = blkioConfig
			}
			if len(blkioExtensions) > 0 {
				data[composeBlkioConfigExtensionsExtension] = blkioExtensions
			}
		}
		if len(compat.CredentialSpec) > 0 {
			data["credential_spec"] = normalizeComposeExtensionObject(compat.CredentialSpec, []string{
				"config",
				"file",
				"registry",
			})
		}
		if len(compat.Provider) > 0 {
			data["provider"] = normalizeComposeExtensionObject(compat.Provider, []string{
				"type",
				"options",
			})
		}
		if len(compat.Extends) > 0 {
			data["extends"] = normalizeComposeExtensionObject(compat.Extends, []string{
				"service",
				"file",
			})
		}
		if compat.Platform != "" {
			data["platform"] = compat.Platform
		}
		if compatFields := serializeComposeCompatExtension(compat); len(compatFields) > 0 {
			data[composeCompatExtension] = compatFields
		}
		if compat.MacAddress != "" {
			data["mac_address"] = compat.MacAddress
		}
		if compat.DomainName != "" {
			data["domainname"] = compat.DomainName
		}
		if compat.CgroupParent != "" {
			data["cgroup_parent"] = compat.CgroupParent
		}
		if compat.Cgroup != "" {
			data["cgroup"] = compat.Cgroup
		}
		if compat.CPUCountSet || compat.CPUCount != 0 {
			data["cpu_count"] = compat.CPUCount
		}
		if compat.CPUPercentSet || compat.CPUPercent != 0 {
			data["cpu_percent"] = compat.CPUPercent
		}
		if compat.CPUPeriodSet || compat.CPUPeriod != 0 {
			data["cpu_period"] = compat.CPUPeriod
		}
		if compat.CPURTPeriodSet || compat.CPURTPeriod != 0 {
			data["cpu_rt_period"] = compat.CPURTPeriod
		}
		if compat.CPURTRuntimeSet || compat.CPURTRuntime != 0 {
			data["cpu_rt_runtime"] = compat.CPURTRuntime
		}
		if compat.CPUSet != "" {
			data["cpuset"] = compat.CPUSet
		}
		if len(compat.DeviceCgroupRules) > 0 {
			data["device_cgroup_rules"] = append([]string{}, compat.DeviceCgroupRules...)
		}
		if len(compat.Gpus) > 0 {
			data["gpus"] = cloneMapSlice(compat.Gpus)
		}
		if compat.NetworkMode != "" {
			data["network_mode"] = compat.NetworkMode
		}
		if compat.OomKillDisableSet {
			data["oom_kill_disable"] = compat.OomKillDisable
		} else if compat.OomKillDisable {
			data["oom_kill_disable"] = true
		}
		if compat.OomScoreAdjSet || compat.OomScoreAdj != 0 {
			data["oom_score_adj"] = compat.OomScoreAdj
		}
		if compat.Scale != nil {
			data["scale"] = *compat.Scale
		}
		if len(compat.Models) > 0 {
			models := map[string]interface{}{}
			for name, model := range compat.Models {
				normalized := map[string]interface{}{}
				for key, value := range model {
					switch key {
					case "endpoint_var", "model_var":
						normalized[key] = value
					default:
						normalized[composeApplicationExtensionKey(key)] = value
					}
				}
				models[name] = normalized
			}
			data["models"] = models
		}
		if len(compat.ExternalLinks) > 0 {
			data["external_links"] = append([]string{}, compat.ExternalLinks...)
		}
		if len(compat.LabelFiles) > 0 {
			data["label_file"] = append([]string{}, compat.LabelFiles...)
		}
		if len(compat.StorageOpt) > 0 {
			data["storage_opt"] = copyStringMap(compat.StorageOpt)
		}
		if compat.UseAPISocketSet {
			data["use_api_socket"] = compat.UseAPISocket
		} else if compat.UseAPISocket {
			data["use_api_socket"] = true
		}
		if compat.Isolation != "" {
			data["isolation"] = compat.Isolation
		}
		if len(compat.Tmpfs) > 0 {
			data["tmpfs"] = append([]string{}, compat.Tmpfs...)
		}
		if compat.Uts != "" {
			data["uts"] = compat.Uts
		}
		if len(compat.VolumesFrom) > 0 {
			data["volumes_from"] = append([]string{}, compat.VolumesFrom...)
		}
	}
	if service.LogDriver != "" || len(service.LogOpt) > 0 || len(service.LogExtensions) > 0 {
		logging := map[string]interface{}{}
		if service.LogDriver != "" {
			logging["driver"] = service.LogDriver
		}
		if len(service.LogOpt) > 0 {
			logging["options"] = copyStringMap(service.LogOpt)
		}
		for key, value := range service.LogExtensions {
			logging[composeApplicationExtensionKey(key)] = value
		}
		data["logging"] = logging
	}
	if service.pidsLimitSet || service.PidsLimit > 0 {
		data["pids_limit"] = service.PidsLimit
	}
	if service.shmSizeSet || service.ShmSize > 0 {
		data["shm_size"] = service.ShmSize
	}
	if len(service.CapAdd) > 0 {
		data["cap_add"] = service.CapAdd
	}
	if len(service.CapDrop) > 0 {
		data["cap_drop"] = service.CapDrop
	}
	if len(service.SecurityOpt) > 0 {
		data["security_opt"] = service.SecurityOpt
	}
	if service.ReadOnlyRootFSSet {
		data["read_only"] = service.ReadOnlyRootFS
	} else if service.ReadOnlyRootFS {
		data["read_only"] = true
	}
	if service.Init != nil {
		data["init"] = *service.Init
	}
	if service.TtySet {
		data["tty"] = service.Tty
	} else if service.Tty {
		data["tty"] = true
	}
	if service.StdinOpenSet {
		data["stdin_open"] = service.StdinOpen
	} else if service.StdinOpen {
		data["stdin_open"] = true
	}
	if service.StopSignal != "" {
		data["stop_signal"] = service.StopSignal
	}
	if service.StopGracePeriod != "" {
		data["stop_grace_period"] = service.StopGracePeriod
	}

	// Volumes
	if len(service.Volumes) > 0 {
		var volumes []interface{}
		for _, volume := range service.Volumes {
			if volumeMountNeedsLongSyntax(volume) {
				volumes = append(volumes, serializeVolumeMount(volume))
				continue
			}
			volStr := fmt.Sprintf("%s:%s", volume.Source, volume.Target)
			if volume.ReadOnly {
				volStr += ":ro"
			}
			if volume.Mode != "" {
				volStr += ":" + volume.Mode
			}
			volumes = append(volumes, volStr)
		}
		data["volumes"] = volumes
	}
	if len(service.Configs) > 0 {
		data["configs"] = serializeFileRefs(service.Configs)
	}
	if len(service.Secrets) > 0 {
		data["secrets"] = serializeFileRefs(service.Secrets)
	}

	// Dependencies
	if dependencies := serviceDependencies(service); len(dependencies) > 0 {
		if serialized := serializeDependencySpecs(app, dependencies); serialized != nil {
			data["depends_on"] = serialized
		}
	}

	// Links
	if len(service.Links) > 0 {
		data["links"] = service.Links
	}

	// Restart
	if service.Restart != "" {
		data["restart"] = service.Restart
	}

	// Resources
	if service.CPUShares > 0 {
		data["cpu_shares"] = service.CPUShares
	}
	if service.CPUQuota > 0 {
		data["cpu_quota"] = service.CPUQuota
	}
	if service.CPUs != "" {
		data["cpus"] = composeCPUValue(service.CPUs)
	}
	if service.MemLimit != "" {
		data["mem_limit"] = composeMemoryValue(service.MemLimit)
		if raw := toString(service.Extensions[composeMemoryLimitExtension]); raw != "" {
			data[composeMemoryLimitExtension] = raw
		}
	} else if service.MemoryLimit != "" {
		data["mem_limit"] = composeMemoryValue(service.MemoryLimit)
		if raw := toString(service.Extensions[composeMemoryLimitExtension]); raw != "" {
			data[composeMemoryLimitExtension] = raw
		}
	}
	if service.MemReservation != "" {
		data["mem_reservation"] = composeMemoryValue(service.MemReservation)
		if raw := toString(service.Extensions[composeMemoryReservationExtension]); raw != "" {
			data[composeMemoryReservationExtension] = raw
		}
	}
	if service.MemorySwap != "" {
		data["memswap_limit"] = composeMemoryValue(service.MemorySwap)
		if raw := toString(service.Extensions[composeMemorySwapExtension]); raw != "" {
			data[composeMemorySwapExtension] = raw
		}
	}
	if ulimits := ulimitsToCompose(service.Ulimits); ulimits != nil {
		data["ulimits"] = ulimits
	}
	if service.ShmSize > 0 {
		if raw := toString(service.Extensions[composeShmSizeExtension]); raw != "" {
			data[composeShmSizeExtension] = raw
		}
	}
	if service.UserNSMode != "" {
		data["userns_mode"] = service.UserNSMode
	}
	if service.PullPolicy != "" {
		data["pull_policy"] = service.PullPolicy
	}
	if len(service.Profiles) > 0 {
		data["profiles"] = service.Profiles
	}
	if connect := nomadConnectSpecToMap(service.Connect); len(connect) > 0 {
		data["x-nomad-connect"] = connect
	}
	if spreads := nomadSpreadSpecsToMap(service.Spreads); len(spreads) > 0 {
		data["x-nomad-spread"] = spreads
	}

	// Healthcheck
	if service.HealthCheck != nil {
		health := normalizeHealthCheck(service.HealthCheck)
		hc := make(map[string]interface{})
		for key, value := range health.Extensions {
			hc[composeApplicationExtensionKey(key)] = value
		}
		if health.DisableSet {
			hc["disable"] = health.Disable
		} else if health.Disable {
			hc["disable"] = true
		}
		if len(health.Test) > 0 {
			hc["test"] = health.Test
		}
		if health.Interval != "" {
			hc["interval"] = health.Interval
		}
		if health.Timeout != "" {
			hc["timeout"] = health.Timeout
		}
		if health.Retries > 0 {
			hc["retries"] = health.Retries
		}
		if health.StartPeriod != "" {
			hc["start_period"] = health.StartPeriod
		}
		if health.StartInterval != "" {
			hc["start_interval"] = health.StartInterval
		}
		data["healthcheck"] = hc
	}
	if service.StartupProbe != nil {
		if probe := serializeKubernetesProbeExtension(service.StartupProbe); probe != nil {
			data["x-kubernetes-startup-probe"] = probe
		}
	}
	if service.SeccompProfile != nil {
		if profile := serializeKubernetesSeccompProfileExtension(service.SeccompProfile); profile != nil {
			data["x-kubernetes-seccomp-profile"] = profile
		}
	}
	if len(service.Affinity) > 0 {
		data["x-kubernetes-affinity"] = copyStringInterfaceMap(service.Affinity)
	}
	if len(service.ReadinessGates) > 0 {
		data["x-kubernetes-readiness-gates"] = cloneMapSlice(service.ReadinessGates)
	}
	if service.AllowPrivilegeEscalation != nil {
		data["x-kubernetes-allowPrivilegeEscalation"] = *service.AllowPrivilegeEscalation
	}
	if service.ProcMount != "" {
		data["x-kubernetes-procMount"] = service.ProcMount
	}
	if len(service.InitContainers) > 0 {
		data["x-kubernetes-init-containers"] = cloneMapSlice(service.InitContainers)
	}
	if len(service.ResourceClaims) > 0 {
		data["x-kubernetes-resource-claims"] = cloneMapSlice(service.ResourceClaims)
	}
	if len(service.EphemeralContainers) > 0 {
		data["x-kubernetes-ephemeral-containers"] = cloneMapSlice(service.EphemeralContainers)
	}
	if len(service.SchedulingGates) > 0 {
		data["x-kubernetes-scheduling-gates"] = cloneMapSlice(service.SchedulingGates)
	}
	if service.HostUsers != nil {
		data["x-kubernetes-hostUsers"] = *service.HostUsers
	}
	if service.HostNetworkSet {
		data["x-kubernetes-hostNetwork"] = service.HostNetwork
	} else if service.HostNetwork {
		data["x-kubernetes-hostNetwork"] = true
	} else if value, ok := extensionBoolValue(service, "kubernetes.hostNetwork"); ok {
		data["x-kubernetes-hostNetwork"] = value
	}
	if service.PIDMode != "" {
		if composeModeReferenceExists(app, service.PIDMode) {
			data["pid"] = service.PIDMode
			if strings.EqualFold(service.PIDMode, "host") {
				data["x-kubernetes-hostPID"] = true
			}
		} else {
			data["x-kubernetes-pidMode"] = service.PIDMode
		}
	} else if value := extensionStringValue(service, "kubernetes.pidMode"); value != "" {
		data["pid"] = value
	} else if service.HostPID != nil {
		data["x-kubernetes-hostPID"] = *service.HostPID
	} else if value, ok := extensionBoolValue(service, "kubernetes.hostPID"); ok {
		data["x-kubernetes-hostPID"] = value
	}
	if service.IPCMode != "" {
		if composeModeReferenceExists(app, service.IPCMode) {
			data["ipc"] = service.IPCMode
			if strings.EqualFold(service.IPCMode, "host") {
				data["x-kubernetes-hostIPC"] = true
			}
		} else {
			data["x-kubernetes-ipcMode"] = service.IPCMode
		}
	} else if value := extensionStringValue(service, "kubernetes.ipcMode"); value != "" {
		data["ipc"] = value
	} else if service.HostIPC != nil {
		data["x-kubernetes-hostIPC"] = *service.HostIPC
	} else if value, ok := extensionBoolValue(service, "kubernetes.hostIPC"); ok {
		data["x-kubernetes-hostIPC"] = value
	}
	if service.PriorityClassName != "" {
		data["x-kubernetes-priorityClassName"] = service.PriorityClassName
	} else if value := extensionStringValue(service, "kubernetes.priorityClassName"); value != "" {
		data["x-kubernetes-priorityClassName"] = value
	}
	if service.RuntimeClassName != "" {
		data["x-kubernetes-runtimeClassName"] = service.RuntimeClassName
	} else if value := extensionStringValue(service, "kubernetes.runtimeClassName"); value != "" {
		data["x-kubernetes-runtimeClassName"] = value
	}
	if service.NodeName != "" {
		data["x-kubernetes-nodeName"] = service.NodeName
	} else if value := extensionStringValue(service, "kubernetes.nodeName"); value != "" {
		data["x-kubernetes-nodeName"] = value
	}
	if service.Subdomain != "" {
		data["x-kubernetes-subdomain"] = service.Subdomain
	} else if value := extensionStringValue(service, "kubernetes.subdomain"); value != "" {
		data["x-kubernetes-subdomain"] = value
	}
	if service.OSName != "" {
		data["x-kubernetes-os"] = service.OSName
	} else if value := extensionStringValue(service, "kubernetes.os"); value != "" {
		data["x-kubernetes-os"] = value
	}
	if service.SetHostnameAsFQDN != nil {
		data["x-kubernetes-setHostnameAsFQDN"] = *service.SetHostnameAsFQDN
	} else if value, ok := extensionBoolValue(service, "kubernetes.setHostnameAsFQDN"); ok {
		data["x-kubernetes-setHostnameAsFQDN"] = value
	}
	if service.ShareProcessNamespace != nil {
		data["x-kubernetes-shareProcessNamespace"] = *service.ShareProcessNamespace
	} else if value, ok := extensionBoolValue(service, "kubernetes.shareProcessNamespace"); ok {
		data["x-kubernetes-shareProcessNamespace"] = value
	}
	if service.EnableServiceLinks != nil {
		data["x-kubernetes-enableServiceLinks"] = *service.EnableServiceLinks
	} else if value, ok := extensionBoolValue(service, "kubernetes.enableServiceLinks"); ok {
		data["x-kubernetes-enableServiceLinks"] = value
	}
	if service.ServiceAccountName != "" {
		data["x-kubernetes-serviceAccountName"] = service.ServiceAccountName
	} else if value := extensionStringValue(service, "kubernetes.serviceAccountName"); value != "" {
		data["x-kubernetes-serviceAccountName"] = value
	}
	if service.AutomountServiceAccountToken != nil {
		data["x-kubernetes-automountServiceAccountToken"] = *service.AutomountServiceAccountToken
	} else if value, ok := extensionBoolValue(service, "kubernetes.automountServiceAccountToken"); ok {
		data["x-kubernetes-automountServiceAccountToken"] = value
	}
	if service.FSGroup != nil && *service.FSGroup > 0 {
		data["x-kubernetes-fsGroup"] = *service.FSGroup
	}
	if selinux := serializeKubernetesSELinuxOptions(service.SELinuxOptions); len(selinux) > 0 {
		data["x-kubernetes-seLinuxOptions"] = selinux
	}
	if windows := serializeKubernetesWindowsSecurityContextOptions(service.WindowsOptions); len(windows) > 0 {
		data["x-kubernetes-windowsOptions"] = windows
	}
	if service.FSGroupChangePolicy != "" {
		data["x-kubernetes-fsGroupChangePolicy"] = service.FSGroupChangePolicy
	}
	if service.RunAsNonRoot != nil {
		data["x-kubernetes-runAsNonRoot"] = *service.RunAsNonRoot
	}
	if len(service.SupplementalGroups) > 0 {
		data["x-kubernetes-supplementalGroups"] = append([]int64{}, service.SupplementalGroups...)
	}
	if service.SupplementalGroupsPolicy != "" {
		data["x-kubernetes-supplementalGroupsPolicy"] = service.SupplementalGroupsPolicy
	}
	if service.ActiveDeadlineSeconds != nil && *service.ActiveDeadlineSeconds > 0 {
		data["x-kubernetes-activeDeadlineSeconds"] = *service.ActiveDeadlineSeconds
	}
	if len(service.Tolerations) > 0 {
		data["x-kubernetes-tolerations"] = serializeKubernetesTolerations(service.Tolerations)
	}
	if service.PodRestartPolicy != "" {
		data["x-kubernetes-restartPolicy"] = service.PodRestartPolicy
	}
	if len(service.TopologySpreadConstraints) > 0 {
		data["x-kubernetes-topology-spread-constraints"] = cloneMapSlice(service.TopologySpreadConstraints)
	}
	if restart := nomadRestartBlockForService(service); len(restart) > 0 {
		data["x-nomad-restart"] = restart
	}

	deploy := cloneDeploySpec(service.Deploy)
	if claims := serviceKubernetesClaimsFromExtensions(service); len(claims) > 0 {
		if deploy == nil {
			deploy = &DeploySpec{}
		}
		if deploy.Resources == nil {
			deploy.Resources = &ResourceSpec{}
		}
		if deploy.Resources.Extensions == nil {
			deploy.Resources.Extensions = map[string]interface{}{}
		}
		if _, ok := deploy.Resources.Extensions["kubernetes.claims"]; !ok {
			deploy.Resources.Extensions["kubernetes.claims"] = cloneInterfaceSlice(claims)
		}
	}
	if deploy != nil {
		data["deploy"] = serializeDeploySpec(deploy)
	}
	if failover := serializeFailoverSpec(service.Failover); len(failover) > 0 {
		data[composeFailoverExtension] = failover
	}

	// Labels
	if len(service.Labels) > 0 {
		data["labels"] = service.Labels
	}

	// Extensions
	for key, value := range convertedServiceExtensions(service.Extensions, service.Platform, PlatformDockerCompose) {
		if key == "compose" {
			continue
		}
		if key == "kubernetes.serviceAccountName" || key == "x-kubernetes-serviceAccountName" {
			continue
		}
		if key == composeFailoverExtension && service.Failover != nil {
			continue
		}
		data[composeServiceExtensionKey(key)] = value
	}
	// Never leak canonical Kubernetes keys into Compose output; Compose users
	// must only see extension-form keys here.
	delete(data, "kubernetes.serviceAccountName")
	delete(data, "kubernetes.automountServiceAccountToken")

	return data, nil
}

func serviceKubernetesClaimsFromExtensions(service *Service) []interface{} {
	if service == nil {
		return nil
	}
	if claims, ok := service.Extensions["x-kubernetes-claims"]; ok {
		if items, ok := claims.([]interface{}); ok && len(items) > 0 {
			return cloneInterfaceSlice(items)
		}
	}
	if claims, ok := service.Extensions["kubernetes.claims"]; ok {
		if items, ok := claims.([]interface{}); ok && len(items) > 0 {
			return cloneInterfaceSlice(items)
		}
	}
	return nil
}

func developConfigToCompose(develop *DevelopConfig) map[string]interface{} {
	if develop == nil {
		return nil
	}
	data := map[string]interface{}{}
	if len(develop.Watch) > 0 {
		watchItems := make([]interface{}, 0, len(develop.Watch))
		for _, watch := range develop.Watch {
			watchItems = append(watchItems, developWatchToCompose(watch))
		}
		data["watch"] = watchItems
	}
	for key, value := range develop.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

func developWatchToCompose(watch DevelopWatch) map[string]interface{} {
	data := map[string]interface{}{}
	if watch.Path != "" {
		data["path"] = watch.Path
	}
	if watch.Action != "" {
		data["action"] = watch.Action
	}
	if watch.Target != "" {
		data["target"] = watch.Target
	}
	if hook := developHookToCompose(watch.Exec); hook != nil {
		data["exec"] = hook
	}
	if len(watch.Include) > 0 {
		data["include"] = watch.Include
	}
	if len(watch.Ignore) > 0 {
		data["ignore"] = watch.Ignore
	}
	if watch.InitialSync {
		data["initial_sync"] = true
	}
	for key, value := range watch.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	return data
}

func serviceHooksToCompose(hooks []ServiceHook) []interface{} {
	if len(hooks) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(hooks))
	for _, hook := range hooks {
		if data := developHookToCompose(&hook); data != nil {
			result = append(result, data)
		}
	}
	return result
}

func developHookToCompose(hook *ServiceHook) map[string]interface{} {
	if hook == nil {
		return nil
	}
	data := map[string]interface{}{}
	if len(hook.Command) > 0 {
		data["command"] = hook.Command
	}
	if hook.Image != "" {
		data["image"] = hook.Image
	}
	if hook.User != "" {
		data["user"] = hook.User
	}
	if hook.Privileged {
		data["privileged"] = true
	}
	if hook.WorkingDir != "" {
		data["working_dir"] = hook.WorkingDir
	}
	if len(hook.Environment) > 0 {
		data["environment"] = copyStringPtrMap(hook.Environment)
	}
	if hook.PerReplica {
		data["per_replica"] = true
	}
	for key, value := range hook.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

func parseComposeLogging(value interface{}) (string, map[string]string, map[string]interface{}, error) {
	data, ok := asMap(value)
	if !ok {
		return "", nil, nil, fmt.Errorf("logging must be a map")
	}
	driver := toString(data["driver"])
	var options map[string]string
	if rawOptions, ok := data["options"]; ok {
		parsed, err := toStringMap(rawOptions)
		if err != nil {
			return "", nil, nil, err
		}
		options = parsed
	}
	extensions := map[string]interface{}{}
	for key, value := range data {
		switch key {
		case "driver", "options":
		default:
			extensions[key] = value
		}
	}
	if len(extensions) == 0 {
		extensions = nil
	}
	return driver, options, extensions, nil
}

func parseComposeModels(value map[string]interface{}) map[string]map[string]interface{} {
	if len(value) == 0 {
		return nil
	}
	result := map[string]map[string]interface{}{}
	for name, raw := range value {
		mapped, ok := asMap(raw)
		if !ok {
			continue
		}
		model := map[string]interface{}{}
		if endpoint := toString(mapped["endpoint_var"]); endpoint != "" {
			model["endpoint_var"] = endpoint
		}
		if variable := toString(mapped["model_var"]); variable != "" {
			model["model_var"] = variable
		}
		for key, item := range mapped {
			switch key {
			case "endpoint_var", "model_var":
			default:
				model[composeApplicationExtensionKey(key)] = deepCopyValue(item)
			}
		}
		if len(model) > 0 {
			result[name] = model
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func portMappingNeedsLongSyntax(port PortMapping) bool {
	return port.Name != "" ||
		port.AppProtocol != "" ||
		port.NodePort != "" ||
		len(port.Extensions) > 0 ||
		(port.Mode != "" && port.Mode != "ingress")
}

func serializePortMappingLong(port PortMapping) map[string]interface{} {
	portData := map[string]interface{}{}
	if port.Name != "" {
		portData["name"] = port.Name
	}
	if parsed := parseInt(port.ContainerPort); parsed > 0 {
		portData["target"] = parsed
	} else if port.ContainerPort != "" {
		portData["target"] = port.ContainerPort
	}
	if port.HostIP != "" {
		portData["host_ip"] = port.HostIP
	}
	if port.HostPort != "" {
		portData["published"] = port.HostPort
	}
	if port.Protocol != "" && port.Protocol != "tcp" {
		portData["protocol"] = port.Protocol
	}
	if port.AppProtocol != "" {
		portData["app_protocol"] = port.AppProtocol
	}
	if port.NodePort != "" {
		portData["x-kubernetes-node-port"] = port.NodePort
	}
	if port.Mode != "" {
		portData["mode"] = port.Mode
	}
	for key, value := range port.Extensions {
		portData[composeApplicationExtensionKey(key)] = deepCopyValue(value)
	}
	return portData
}

func serializeFileRefs(refs []FileRef) []interface{} {
	result := make([]interface{}, 0, len(refs))
	for _, ref := range refs {
		if ref.Key == "" && ref.Target == "" && ref.UID == "" && ref.GID == "" && ref.Mode == "" && !ref.ReadOnly && ref.Optional == nil && len(ref.Extensions) == 0 {
			result = append(result, ref.Source)
			continue
		}
		item := map[string]interface{}{"source": ref.Source}
		if ref.Key != "" {
			item["x-kubernetes-key"] = ref.Key
		}
		if ref.Target != "" {
			item["target"] = ref.Target
		}
		if ref.UID != "" {
			item["uid"] = ref.UID
		}
		if ref.GID != "" {
			item["gid"] = ref.GID
		}
		if ref.Mode != "" {
			item["mode"] = ref.Mode
		}
		if ref.ReadOnly {
			item["read_only"] = true
		}
		if ref.Optional != nil {
			item["x-kubernetes-optional"] = *ref.Optional
		}
		for key, value := range ref.Extensions {
			if key == "x-kubernetes-key" || key == "x-kubernetes-optional" {
				continue
			}
			item[composeApplicationExtensionKey(key)] = value
		}
		result = append(result, item)
	}
	return result
}

func volumeMountNeedsLongSyntax(volume VolumeMount) bool {
	return volume.Type == "tmpfs" ||
		volume.Consistency != "" ||
		volume.Propagation != "" ||
		volume.MountPropagation != "" ||
		volume.RecursiveReadOnly != "" ||
		volume.SubPath != "" ||
		volume.SubPathExpr != "" ||
		volume.NoCopy ||
		volume.CreateHostPath != nil ||
		len(volume.BindExtensions) > 0 ||
		len(volume.VolumeLabels) > 0 ||
		len(volume.VolumeExtensions) > 0 ||
		len(volume.TmpfsExtensions) > 0 ||
		volume.TmpfsMode != "" ||
		volume.ImageSubpath != "" ||
		len(volume.ImageExtensions) > 0 ||
		len(volume.Options) > 0
}

func serializeVolumeMount(volume VolumeMount) map[string]interface{} {
	item := map[string]interface{}{}
	if volume.Type != "" {
		item["type"] = volume.Type
	}
	if volume.Source != "" {
		item["source"] = volume.Source
	}
	if volume.Target != "" {
		item["target"] = volume.Target
	}
	if volume.ReadOnly {
		item["read_only"] = true
	}
	if volume.Consistency != "" {
		item["consistency"] = volume.Consistency
	}
	if volume.MountPropagation != "" {
		item[composeApplicationExtensionKey("kubernetes.mountPropagation")] = volume.MountPropagation
	}
	if volume.RecursiveReadOnly != "" {
		item[composeApplicationExtensionKey("kubernetes.recursiveReadOnly")] = volume.RecursiveReadOnly
	}
	for key, value := range volume.Extensions {
		item[composeApplicationExtensionKey(key)] = value
	}
	switch volume.Type {
	case "bind":
		bind := map[string]interface{}{}
		if volume.Mode != "" {
			bind["selinux"] = volume.Mode
		}
		if volume.Propagation != "" {
			bind["propagation"] = volume.Propagation
		}
		if volume.CreateHostPath != nil {
			bind["create_host_path"] = *volume.CreateHostPath
		}
		if recursive := volume.Options["recursive"]; recursive != "" {
			bind["recursive"] = recursive
		}
		for key, value := range volume.BindExtensions {
			bind[composeApplicationExtensionKey(key)] = value
		}
		if len(bind) > 0 {
			item["bind"] = bind
		}
	case "volume":
		volumeOptions := map[string]interface{}{}
		if volume.NoCopy {
			volumeOptions["nocopy"] = true
		}
		if volume.SubPath != "" {
			volumeOptions["subpath"] = volume.SubPath
		} else if subpath := volume.Options["subpath"]; subpath != "" {
			volumeOptions["subpath"] = subpath
		}
		if len(volume.VolumeLabels) > 0 {
			volumeOptions["labels"] = volume.VolumeLabels
		}
		for key, value := range volume.VolumeExtensions {
			volumeOptions[composeApplicationExtensionKey(key)] = value
		}
		if len(volumeOptions) > 0 {
			item["volume"] = volumeOptions
		}
		if volume.SubPathExpr != "" {
			item[composeApplicationExtensionKey("kubernetes.subPathExpr")] = volume.SubPathExpr
		}
	case "tmpfs":
		tmpfs := map[string]interface{}{}
		if size := volume.Options["size"]; size != "" {
			if parsed := parseInt(size); parsed > 0 {
				tmpfs["size"] = parsed
			} else {
				tmpfs["size"] = size
			}
		}
		if len(tmpfs) > 0 {
			item["tmpfs"] = tmpfs
		}
		if volume.TmpfsMode != "" {
			itemTmpfs, ok := item["tmpfs"].(map[string]interface{})
			if !ok {
				itemTmpfs = map[string]interface{}{}
				item["tmpfs"] = itemTmpfs
			}
			if parsed := parseInt(volume.TmpfsMode); parsed > 0 {
				itemTmpfs["mode"] = parsed
			} else {
				itemTmpfs["mode"] = volume.TmpfsMode
			}
		}
		for key, value := range volume.TmpfsExtensions {
			itemTmpfs, ok := item["tmpfs"].(map[string]interface{})
			if !ok {
				itemTmpfs = map[string]interface{}{}
				item["tmpfs"] = itemTmpfs
			}
			itemTmpfs[composeApplicationExtensionKey(key)] = value
		}
	case "image":
		image := map[string]interface{}{}
		if volume.ImageSubpath != "" {
			image["subpath"] = volume.ImageSubpath
		}
		for key, value := range volume.ImageExtensions {
			image[composeApplicationExtensionKey(key)] = value
		}
		if len(image) > 0 {
			item["image"] = image
		}
	}
	return item
}

func serializeEnvFileRefs(refs []EnvFileRef) []interface{} {
	result := make([]interface{}, 0, len(refs))
	for _, ref := range refs {
		if ref.Format == "" && ref.Required == nil && len(ref.Extensions) == 0 {
			result = append(result, ref.Path)
			continue
		}
		item := map[string]interface{}{"path": ref.Path}
		if ref.Required != nil {
			item["required"] = *ref.Required
		}
		if ref.Format != "" {
			item["format"] = ref.Format
		}
		for key, value := range ref.Extensions {
			item[composeApplicationExtensionKey(key)] = value
		}
		result = append(result, item)
	}
	return result
}

func serializeEnvSources(sources []EnvSource) []interface{} {
	result := make([]interface{}, 0, len(sources))
	for _, source := range sources {
		item := map[string]interface{}{
			"name":        source.Name,
			"source_type": source.SourceType,
			"source":      source.Source,
		}
		if source.Key != "" {
			item["key"] = source.Key
		}
		if source.Optional {
			item["optional"] = true
		}
		for key, value := range source.Extensions {
			item[composeApplicationExtensionKey(key)] = deepCopyValue(value)
		}
		result = append(result, item)
	}
	return result
}

func serializeEnvFromSources(sources []EnvFromSource) []interface{} {
	result := make([]interface{}, 0, len(sources))
	for _, source := range sources {
		item := map[string]interface{}{
			"source_type": source.SourceType,
			"source":      source.Source,
		}
		if source.Prefix != "" {
			item["prefix"] = source.Prefix
		}
		if source.Optional {
			item["optional"] = true
		}
		for key, value := range source.Extensions {
			item[composeApplicationExtensionKey(key)] = deepCopyValue(value)
		}
		result = append(result, item)
	}
	return result
}

func serializeDeploySpec(deploy *DeploySpec) map[string]interface{} {
	data := make(map[string]interface{})
	for key, value := range deploy.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}
	if deploy.Mode != "" {
		data["mode"] = deploy.Mode
	}
	if deploy.EndpointMode != "" {
		data["endpoint_mode"] = deploy.EndpointMode
	}
	if deploy.Replicas > 0 {
		data["replicas"] = deploy.Replicas
	}
	if !isEmptySwarmJobSpec(deploy.Job) {
		data[composeSwarmJobExtension] = serializeSwarmJobSpec(deploy.Job)
	}
	if !isEmptyMigratePolicy(deploy.MigrateConfig) {
		data["x-nomad-migrate"] = serializeMigratePolicy(deploy.MigrateConfig)
	}
	if !isEmptyReschedulePolicy(deploy.RescheduleConfig) {
		data["x-nomad-reschedule"] = serializeReschedulePolicy(deploy.RescheduleConfig)
	}
	if len(deploy.Labels) > 0 {
		data["labels"] = deploy.Labels
	}
	if deploy.Placement != nil {
		placement := make(map[string]interface{})
		for key, value := range deploy.Placement.Extensions {
			placement[composeApplicationExtensionKey(key)] = value
		}
		if len(deploy.Placement.Constraints) > 0 {
			placement["constraints"] = deploy.Placement.Constraints
		}
		if len(deploy.Placement.Preferences) > 0 {
			placement["preferences"] = serializePlacementPreferences(deploy.Placement.Preferences, deploy.Placement.PreferenceExtensions)
		}
		if deploy.Placement.MaxReplicasPerNode > 0 {
			placement["max_replicas_per_node"] = deploy.Placement.MaxReplicasPerNode
		}
		if len(placement) > 0 {
			data["placement"] = placement
		}
	}
	if deploy.Resources != nil {
		resources := make(map[string]interface{})
		for key, value := range deploy.Resources.Extensions {
			resources[composeApplicationExtensionKey(key)] = value
		}
		limits := make(map[string]interface{})
		for key, value := range deploy.Resources.LimitExtensions {
			limits[composeApplicationExtensionKey(key)] = value
		}
		if deploy.Resources.CPULimit != "" {
			limits["cpus"] = composeCPUValue(deploy.Resources.CPULimit)
		}
		if deploy.Resources.MemoryLimit != "" {
			limits["memory"] = composeMemoryValue(deploy.Resources.MemoryLimit)
		}
		if deploy.Resources.EphemeralStorageLimit != "" {
			if raw := toString(deploy.Resources.LimitExtensions["x-kubernetes-ephemeral-storage-limit"]); raw != "" {
				limits["x-kubernetes-ephemeral-storage-limit"] = raw
			} else {
				limits["x-kubernetes-ephemeral-storage-limit"] = deploy.Resources.EphemeralStorageLimit
			}
		}
		if deploy.Resources.pidsLimitSet || deploy.Resources.PidsLimit > 0 {
			limits["pids"] = deploy.Resources.PidsLimit
		}
		if len(limits) > 0 {
			resources["limits"] = limits
		}
		reservations := make(map[string]interface{})
		for key, value := range deploy.Resources.ReservationExtensions {
			reservations[composeApplicationExtensionKey(key)] = value
		}
		if deploy.Resources.CPUReservation != "" {
			reservations["cpus"] = composeCPUValue(deploy.Resources.CPUReservation)
		}
		if deploy.Resources.MemoryReservation != "" {
			reservations["memory"] = composeMemoryValue(deploy.Resources.MemoryReservation)
		}
		if deploy.Resources.EphemeralStorageReservation != "" {
			if raw := toString(deploy.Resources.ReservationExtensions["x-kubernetes-ephemeral-storage-reservation"]); raw != "" {
				reservations["x-kubernetes-ephemeral-storage-reservation"] = raw
			} else {
				reservations["x-kubernetes-ephemeral-storage-reservation"] = deploy.Resources.EphemeralStorageReservation
			}
		}
		if deploy.Resources.pidsReservationSet || deploy.Resources.PidsReservation > 0 {
			reservations["pids"] = deploy.Resources.PidsReservation
		}
		if len(deploy.Resources.Devices) > 0 {
			reservations["devices"] = serializeResourceDevices(deploy.Resources.Devices)
		}
		if len(deploy.Resources.GenericResources) > 0 {
			reservations["generic_resources"] = serializeGenericResources(deploy.Resources.GenericResources)
		}
		if len(reservations) > 0 {
			resources["reservations"] = reservations
		}
		if len(resources) > 0 {
			data["resources"] = resources
		}
	}
	if deploy.UpdateConfig != nil {
		if update := serializeUpdatePolicy(deploy.UpdateConfig); len(update) > 0 {
			data["update_config"] = update
		}
		if nomadUpdate := serializeNomadUpdatePolicy(deploy.UpdateConfig); len(nomadUpdate) > 0 {
			data["x-nomad-update"] = nomadUpdate
		}
	}
	if deploy.RollbackConfig != nil {
		if rollback := serializeUpdatePolicy(deploy.RollbackConfig); len(rollback) > 0 {
			data["rollback_config"] = rollback
		}
	}
	if deploy.RestartPolicy != nil {
		restart := make(map[string]interface{})
		for key, value := range deploy.RestartPolicy.Extensions {
			restart[composeApplicationExtensionKey(key)] = value
		}
		if deploy.RestartPolicy.Condition != "" {
			restart["condition"] = deploy.RestartPolicy.Condition
		}
		if deploy.RestartPolicy.Delay != "" {
			restart["delay"] = deploy.RestartPolicy.Delay
		}
		if deploy.RestartPolicy.MaxAttempts > 0 {
			restart["max_attempts"] = deploy.RestartPolicy.MaxAttempts
		}
		if deploy.RestartPolicy.Window != "" {
			restart["window"] = deploy.RestartPolicy.Window
		}
		if len(restart) > 0 {
			data["restart_policy"] = restart
		}
	}
	return data
}

func serializeSwarmJobSpec(job *SwarmJobSpec) map[string]interface{} {
	data := map[string]interface{}{}
	if job == nil {
		return data
	}
	if job.maxConcurrentSet || job.MaxConcurrent > 0 {
		data["max_concurrent"] = job.MaxConcurrent
		data["parallelism"] = job.MaxConcurrent
	}
	if job.totalCompletionsSet || job.TotalCompletions > 0 {
		data["total_completions"] = job.TotalCompletions
		data["completions"] = job.TotalCompletions
	}
	if job.completionModeSet || job.CompletionMode != "" {
		data["completion_mode"] = job.CompletionMode
	}
	if job.Suspend != nil {
		data["suspend"] = *job.Suspend
	}
	if job.backoffLimitSet || job.BackoffLimit > 0 {
		data["backoff_limit"] = job.BackoffLimit
	}
	if job.backoffLimitPerIndexSet || job.BackoffLimitPerIndex > 0 {
		data["backoff_limit_per_index"] = job.BackoffLimitPerIndex
	}
	if job.ttlSecondsAfterFinishedSet || job.TTLSecondsAfterFinished > 0 {
		data["ttl_seconds_after_finished"] = job.TTLSecondsAfterFinished
	}
	for key, value := range job.Extensions {
		data[key] = deepCopyValue(value)
	}
	return data
}

func serializePlacementPreferences(preferences []string, extensions []map[string]interface{}) []map[string]interface{} {
	if len(preferences) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(preferences))
	for i, preference := range preferences {
		item := map[string]interface{}{}
		if i < len(extensions) {
			for key, value := range extensions[i] {
				item[composeApplicationExtensionKey(key)] = value
			}
		}
		if spread := composeSpreadFromPlacementPreference(preference); spread != "" {
			item["spread"] = spread
			if normalizePlacementPreference(preference) != "spread="+spread {
				item["x-bolabaden-preference"] = preference
			}
		} else if preference != "" {
			item["spread"] = preference
			item["x-bolabaden-preference"] = preference
		}
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	return result
}

func composeSpreadFromPlacementPreference(preference string) string {
	preference = strings.TrimSpace(strings.TrimPrefix(preference, "prefer:"))
	if strings.HasPrefix(preference, "spread=") {
		return strings.TrimSpace(strings.TrimPrefix(preference, "spread="))
	}
	left, operator, right := splitPortableConstraint(preference)
	if left != "" && (operator == "=" || operator == "==") && strings.EqualFold(right, "true") {
		return left
	}
	return ""
}

func serializeResourceDevices(devices []ResourceDevice) []interface{} {
	result := make([]interface{}, 0, len(devices))
	for _, device := range devices {
		item := map[string]interface{}{}
		for key, value := range device.Extensions {
			item[composeApplicationExtensionKey(key)] = value
		}
		if len(device.Capabilities) > 0 {
			item["capabilities"] = device.Capabilities
		}
		if device.Driver != "" {
			item["driver"] = device.Driver
		}
		if device.Count != "" {
			item["count"] = device.Count
		}
		if len(device.DeviceIDs) > 0 {
			item["device_ids"] = device.DeviceIDs
		}
		if len(device.Options) > 0 {
			item["options"] = device.Options
		}
		result = append(result, item)
	}
	return result
}

func serializeGenericResources(resources []GenericResource) []interface{} {
	result := make([]interface{}, 0, len(resources))
	for _, resource := range resources {
		discrete := map[string]interface{}{}
		for key, value := range resource.DiscreteExtensions {
			discrete[composeApplicationExtensionKey(key)] = value
		}
		if resource.Kind != "" {
			discrete["kind"] = resource.Kind
		}
		if resource.Value != "" {
			discrete["value"] = resource.Value
		}
		if len(discrete) == 0 && len(resource.Extensions) == 0 {
			continue
		}
		item := map[string]interface{}{}
		for key, value := range resource.Extensions {
			item[composeApplicationExtensionKey(key)] = value
		}
		if len(discrete) > 0 {
			item["discrete_resource_spec"] = discrete
		}
		result = append(result, item)
	}
	return result
}

func serializeUpdatePolicy(policy *UpdatePolicy) map[string]interface{} {
	update := make(map[string]interface{})
	if policy == nil {
		return update
	}
	for key, value := range policy.Extensions {
		update[composeApplicationExtensionKey(key)] = value
	}
	if policy.ParallelismSet || policy.Parallelism > 0 {
		update["parallelism"] = policy.Parallelism
	}
	if policy.Delay != "" {
		update["delay"] = policy.Delay
	}
	if policy.Monitor != "" {
		update["monitor"] = policy.Monitor
	}
	if policy.MaxFailureRatio != "" {
		update["max_failure_ratio"] = policy.MaxFailureRatio
	}
	if policy.Order != "" {
		update["order"] = policy.Order
	}
	if policy.OnFailure != "" {
		update["failure_action"] = policy.OnFailure
	}
	return update
}

func serializeNomadUpdatePolicy(policy *UpdatePolicy) map[string]interface{} {
	update := make(map[string]interface{})
	if policy == nil {
		return update
	}
	for key, value := range policy.Extensions {
		update[key] = value
	}
	if policy.ParallelismSet || policy.Parallelism > 0 {
		update["max_parallel"] = policy.Parallelism
	}
	if policy.Delay != "" {
		update["delay"] = policy.Delay
	}
	if policy.Monitor != "" {
		update["monitor"] = policy.Monitor
	}
	if policy.MaxFailureRatio != "" {
		update["max_failure_ratio"] = policy.MaxFailureRatio
	}
	if policy.Order != "" {
		update["order"] = policy.Order
	}
	if policy.OnFailure != "" {
		update["failure_action"] = policy.OnFailure
	}
	if policy.HealthCheck != "" {
		update["health_check"] = policy.HealthCheck
	}
	if policy.MinHealthyTime != "" {
		update["min_healthy_time"] = policy.MinHealthyTime
	}
	if policy.HealthyDeadline != "" {
		update["healthy_deadline"] = policy.HealthyDeadline
	}
	if policy.ProgressDeadline != "" {
		update["progress_deadline"] = policy.ProgressDeadline
	}
	if policy.AutoRevertSet || policy.AutoRevert {
		update["auto_revert"] = policy.AutoRevert
	}
	if policy.AutoPromoteSet || policy.AutoPromote {
		update["auto_promote"] = policy.AutoPromote
	}
	if policy.CanarySet || policy.Canary > 0 {
		update["canary"] = policy.Canary
	}
	if policy.Stagger != "" {
		update["stagger"] = policy.Stagger
	}
	return update
}

func serializeMigratePolicy(policy *MigratePolicy) map[string]interface{} {
	migrate := make(map[string]interface{})
	if policy == nil {
		return migrate
	}
	for key, value := range policy.Extensions {
		migrate[composeApplicationExtensionKey(key)] = value
	}
	if policy.MaxParallel > 0 {
		migrate["max_parallel"] = policy.MaxParallel
	}
	if policy.HealthCheck != "" {
		migrate["health_check"] = policy.HealthCheck
	}
	if policy.MinHealthyTime != "" {
		migrate["min_healthy_time"] = policy.MinHealthyTime
	}
	if policy.HealthyDeadline != "" {
		migrate["healthy_deadline"] = policy.HealthyDeadline
	}
	return migrate
}

func serializeReschedulePolicy(policy *ReschedulePolicy) map[string]interface{} {
	reschedule := make(map[string]interface{})
	if policy == nil {
		return reschedule
	}
	for key, value := range policy.Extensions {
		reschedule[composeApplicationExtensionKey(key)] = value
	}
	if policy.Attempts > 0 {
		reschedule["attempts"] = policy.Attempts
	}
	if policy.Interval != "" {
		reschedule["interval"] = policy.Interval
	}
	if policy.Delay != "" {
		reschedule["delay"] = policy.Delay
	}
	if policy.DelayFunction != "" {
		reschedule["delay_function"] = policy.DelayFunction
	}
	if policy.MaxDelay != "" {
		reschedule["max_delay"] = policy.MaxDelay
	}
	if policy.Unlimited {
		reschedule["unlimited"] = true
	}
	return reschedule
}

func serializeNetwork(network *Network) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	if network.PlatformName != "" {
		data["name"] = network.PlatformName
	}
	if network.Driver != "" {
		data["driver"] = network.Driver
	}
	if len(network.DriverOpts) > 0 {
		data["driver_opts"] = network.DriverOpts
	}
	if network.AttachableSet {
		data["attachable"] = network.Attachable
	} else if network.Attachable {
		data["attachable"] = true
	}
	if network.External || network.ExternalSet {
		data["external"] = serializeComposeExternal(network.External, network.ExternalSet, network.PlatformName, network.ExternalExtensions)
	}
	if network.InternalSet {
		data["internal"] = network.Internal
	} else if network.Internal {
		data["internal"] = true
	}
	if network.EnableIPv4 != nil {
		data["enable_ipv4"] = *network.EnableIPv4
	}
	if network.EnableIPv6 != nil {
		data["enable_ipv6"] = *network.EnableIPv6
	}
	if network.IPAM != nil {
		ipamData := make(map[string]interface{})
		if network.IPAM.Driver != "" {
			ipamData["driver"] = network.IPAM.Driver
		}
		if len(network.IPAM.Config) > 0 {
			var config []map[string]interface{}
			for _, subnet := range network.IPAM.Config {
				subnetData := make(map[string]interface{})
				if subnet.Subnet != "" {
					subnetData["subnet"] = subnet.Subnet
				}
				if subnet.Gateway != "" {
					subnetData["gateway"] = subnet.Gateway
				}
				if subnet.IPRange != "" {
					subnetData["ip_range"] = subnet.IPRange
				}
				if len(subnet.AuxAddresses) > 0 {
					subnetData["aux_addresses"] = subnet.AuxAddresses
				}
				for key, value := range subnet.Extensions {
					subnetData[composeApplicationExtensionKey(key)] = value
				}
				config = append(config, subnetData)
			}
			ipamData["config"] = config
		}
		if len(network.IPAM.Options) > 0 {
			ipamData["options"] = network.IPAM.Options
		}
		for key, value := range network.IPAM.Extensions {
			ipamData[composeApplicationExtensionKey(key)] = value
		}
		data["ipam"] = ipamData
	}
	if len(network.Labels) > 0 {
		data["labels"] = network.Labels
	}
	for key, value := range network.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}

	return data, nil
}

func serializeComposeExternal(external bool, externalSet bool, platformName string, extensions map[string]interface{}) interface{} {
	if !external && !externalSet {
		return false
	}
	if len(extensions) == 0 {
		if externalSet {
			return external
		}
		return true
	}
	data := map[string]interface{}{}
	if platformName != "" {
		data["name"] = platformName
	}
	for key, value := range extensions {
		data[composeApplicationExtensionKey(key)] = deepCopyValue(value)
	}
	if len(data) == 0 {
		if externalSet {
			return external
		}
		return true
	}
	return data
}

func serializeVolume(volume *Volume) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	if volume.PlatformName != "" {
		data["name"] = volume.PlatformName
	}
	if volume.Driver != "" && !(volume.External || volume.ExternalSet) {
		data["driver"] = volume.Driver
	}
	if len(volume.DriverOpts) > 0 && !(volume.External || volume.ExternalSet) {
		data["driver_opts"] = volume.DriverOpts
	}
	if volume.External || volume.ExternalSet {
		data["external"] = serializeComposeExternal(volume.External, volume.ExternalSet, volume.PlatformName, volume.ExternalExtensions)
		if volume.Driver != "" {
			data["x-bolabaden-volume-driver"] = volume.Driver
		}
		if len(volume.DriverOpts) > 0 {
			data["x-bolabaden-volume-driver-opts"] = volume.DriverOpts
		}
	}
	if len(volume.Labels) > 0 {
		data["labels"] = volume.Labels
	}
	if kind := volumeKubernetesKind(volume); kind != "" {
		data["x-kubernetes-kind"] = kind
	}
	for key, value := range volume.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}

	return data, nil
}

func serializeConfig(config *Config) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	if config.PlatformName != "" {
		data["name"] = config.PlatformName
	}
	if config.Content != "" && !(config.External || config.ExternalSet) {
		data["content"] = config.Content
	}
	if config.Environment != "" && !(config.External || config.ExternalSet) {
		data["environment"] = config.Environment
	}
	if config.File != "" && !(config.External || config.ExternalSet) {
		data["file"] = config.File
	}
	if config.Template != "" && !(config.External || config.ExternalSet) {
		data["template_driver"] = config.Template
	}
	if config.Mode != "" && !(config.External || config.ExternalSet) {
		data["mode"] = config.Mode
	}
	if config.External || config.ExternalSet {
		data["external"] = serializeComposeExternal(config.External, config.ExternalSet, config.PlatformName, config.ExternalExtensions)
		if config.Content != "" {
			data["x-bolabaden-config-content"] = config.Content
		}
		if config.Environment != "" {
			data["x-bolabaden-config-environment"] = config.Environment
		}
		if config.File != "" {
			data["x-bolabaden-config-file"] = config.File
		}
		if config.Template != "" {
			data["x-bolabaden-config-template"] = config.Template
		}
		if config.Mode != "" {
			data["x-bolabaden-config-mode"] = config.Mode
		}
	}
	if len(config.Labels) > 0 {
		data["labels"] = config.Labels
	}
	if kubernetesLabels := configKubernetesLabels(config); len(kubernetesLabels) > 0 {
		data["x-kubernetes-labels"] = kubernetesLabels
	}
	if kubernetesAnnotations := configKubernetesAnnotations(config); len(kubernetesAnnotations) > 0 {
		data["x-kubernetes-annotations"] = kubernetesAnnotations
	}
	if kubernetesImmutable := configKubernetesImmutable(config); kubernetesImmutable != nil {
		data["x-kubernetes-immutable"] = *kubernetesImmutable
	}
	if kubernetesData := configKubernetesData(config); len(kubernetesData) > 0 {
		data["x-kubernetes-data"] = kubernetesData
	}
	if kubernetesBinaryData := configKubernetesBinaryData(config); len(kubernetesBinaryData) > 0 {
		data["x-kubernetes-binaryData"] = kubernetesBinaryData
	}
	for key, value := range config.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}

	return data, nil
}

func serializeSecret(secret *Secret) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	if secret.PlatformName != "" {
		data["name"] = secret.PlatformName
	}
	if secret.File != "" && !(secret.External || secret.ExternalSet) {
		data["file"] = secret.File
	}
	if secret.Environment != "" && !(secret.External || secret.ExternalSet) {
		data["environment"] = secret.Environment
	}
	if secret.Template != "" && !(secret.External || secret.ExternalSet) {
		data["template_driver"] = secret.Template
	}
	if secret.External || secret.ExternalSet {
		data["external"] = serializeComposeExternal(secret.External, secret.ExternalSet, secret.PlatformName, secret.ExternalExtensions)
		if secret.File != "" {
			data["x-bolabaden-secret-file"] = secret.File
		}
		if secret.Environment != "" {
			data["x-bolabaden-secret-environment"] = secret.Environment
		}
		if secret.Template != "" {
			data["x-bolabaden-secret-template"] = secret.Template
		}
		if secret.Driver != "" {
			data["x-bolabaden-secret-driver"] = secret.Driver
		}
		if len(secret.DriverOpts) > 0 {
			data["x-bolabaden-secret-driver-opts"] = secret.DriverOpts
		}
	}
	if secret.Driver != "" && !(secret.External || secret.ExternalSet) {
		data["driver"] = secret.Driver
	}
	if len(secret.DriverOpts) > 0 && !(secret.External || secret.ExternalSet) {
		data["driver_opts"] = secret.DriverOpts
	}
	if len(secret.Labels) > 0 {
		data["labels"] = secret.Labels
	}
	if kubernetesLabels := secretKubernetesLabels(secret); len(kubernetesLabels) > 0 {
		data["x-kubernetes-labels"] = kubernetesLabels
	}
	if kubernetesAnnotations := secretKubernetesAnnotations(secret); len(kubernetesAnnotations) > 0 {
		data["x-kubernetes-annotations"] = kubernetesAnnotations
	}
	if kubernetesImmutable := secretKubernetesImmutable(secret); kubernetesImmutable != nil {
		data["x-kubernetes-immutable"] = *kubernetesImmutable
	}
	if kubernetesData := secretKubernetesData(secret); len(kubernetesData) > 0 {
		data["x-kubernetes-data"] = kubernetesData
	}
	if kubernetesStringData := secretKubernetesStringData(secret); len(kubernetesStringData) > 0 {
		data["x-kubernetes-stringData"] = kubernetesStringData
	}
	if kubernetesType := secretKubernetesType(secret); kubernetesType != "" && kubernetesType != "Opaque" {
		data["x-kubernetes-type"] = kubernetesType
	}
	for key, value := range secret.Extensions {
		data[composeApplicationExtensionKey(key)] = value
	}

	return data, nil
}
