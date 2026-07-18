package paas

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

const helmRawFilesExtension = "helm.files"
const helmAppExtensionsFile = ".bolabaden.app.extensions.json"
const helmServiceExtensionsFile = ".bolabaden.service.extensions.json"
const helmAppNameFile = ".bolabaden.app.name"
const helmAppModelsFile = ".bolabaden.app.models.json"
const helmAppIncludesFile = ".bolabaden.app.includes.json"
const helmAppNetworksFile = ".bolabaden.app.networks.json"
const helmAppVolumesFile = ".bolabaden.app.volumes.json"
const helmAppConfigsFile = ".bolabaden.app.configs.json"
const helmAppSecretsFile = ".bolabaden.app.secrets.json"
const helmAppRoutesFile = ".bolabaden.app.routes.json"
const helmAppPoliciesFile = ".bolabaden.app.policies.json"
const helmCanonicalRawResourcesFile = ".bolabaden.canonical.raw-resources.json"
const helmKubernetesRawResourcesFile = ".bolabaden.kubernetes.raw-resources.json"

// ParseHelmChart parses a Helm chart directory
func ParseHelmChart(chartPath string) (*Application, error) {
	// Check if it's a chart directory
	chartYaml := filepath.Join(chartPath, "Chart.yaml")
	if _, err := os.Stat(chartYaml); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a Helm chart directory: %s", chartPath)
	}

	// Read Chart.yaml
	chartData, err := os.ReadFile(chartYaml)
	if err != nil {
		return nil, fmt.Errorf("failed to read Chart.yaml: %w", err)
	}

	var chartMeta map[string]interface{}
	if err := yaml.Unmarshal(chartData, &chartMeta); err != nil {
		return nil, fmt.Errorf("failed to parse Chart.yaml: %w", err)
	}

	app := &Application{
		Platform: PlatformHelm,
		Services: make(map[string]*Service),
	}

	chrt, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load Helm chart: %w", err)
	}

	rendered, err := renderLoadedHelmChart(chrt)
	if err != nil {
		return nil, err
	}
	k8sApp, err := ParseKubernetesYAML(rendered)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rendered Helm manifests: %w", err)
	}
	mergeRenderedHelmApp(app, k8sApp)
	if extensions, err := helmApplicationExtensionsFromChart(chartPath); err == nil && len(extensions) > 0 {
		if app.Extensions == nil {
			app.Extensions = map[string]interface{}{}
		}
		for key, value := range extensions {
			if key == helmRawFilesExtension || key == "chart" {
				continue
			}
			if key == helmAppRoutesFile || key == helmAppPoliciesFile || key == helmAppNameFile || key == helmAppIncludesFile || key == helmAppNetworksFile || key == helmAppVolumesFile || key == helmAppConfigsFile || key == helmAppSecretsFile || key == helmCanonicalRawResourcesFile || key == helmKubernetesRawResourcesFile {
				continue
			}
			app.Extensions[key] = value
		}
		for key, value := range copyStringInterfaceMap(app.Extensions) {
			if canonical := composeApplicationCanonicalKey(key); canonical != "" {
				app.Extensions[canonical] = value
			}
		}
	}
	if serviceExtensions, err := helmApplicationServiceExtensionsFromChart(chartPath); err == nil && len(serviceExtensions) > 0 {
		for serviceName, extensions := range serviceExtensions {
			service := app.Services[serviceName]
			if service == nil {
				continue
			}
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			for key, value := range extensions {
				service.Extensions[key] = value
			}
			for key, value := range copyStringInterfaceMap(service.Extensions) {
				if canonical := composeApplicationCanonicalKey(key); canonical != "" {
					service.Extensions[canonical] = value
				}
			}
		}
	}
	if name, err := helmApplicationNameFromChart(chartPath); err == nil && name != "" {
		app.Name = name
	}
	if version := toString(chartMeta["appVersion"]); version != "" {
		app.Version = version
	}
	if models, err := helmApplicationModelsFromChart(chartPath); err == nil && len(models) > 0 {
		app.Models = models
	}
	if includes, err := helmApplicationIncludesFromChart(chartPath); err == nil && len(includes) > 0 {
		app.IncludeEntries = mergeIncludeEntries(app.IncludeEntries, includes)
		app.Includes = mergeUniqueStrings(app.Includes, composeIncludePaths(includes))
	}
	if networks, err := helmApplicationNetworksFromChart(chartPath); err == nil && len(networks) > 0 {
		if app.Networks == nil {
			app.Networks = map[string]*Network{}
		}
		for name, network := range networks {
			app.Networks[name] = cloneNetwork(network)
		}
	}
	if volumes, err := helmApplicationVolumesFromChart(chartPath); err == nil && len(volumes) > 0 {
		if app.Volumes == nil {
			app.Volumes = map[string]*Volume{}
		}
		for name, volume := range volumes {
			app.Volumes[name] = cloneVolume(volume)
		}
	}
	if configs, err := helmApplicationConfigsFromChart(chartPath); err == nil && len(configs) > 0 {
		if app.Configs == nil {
			app.Configs = map[string]*Config{}
		}
		for name, config := range configs {
			app.Configs[name] = cloneConfig(config)
		}
	}
	if secrets, err := helmApplicationSecretsFromChart(chartPath); err == nil && len(secrets) > 0 {
		if app.Secrets == nil {
			app.Secrets = map[string]*Secret{}
		}
		for name, secret := range secrets {
			app.Secrets[name] = cloneSecret(secret)
		}
	}
	if routes, err := helmApplicationRoutesFromChart(chartPath); err == nil && len(routes) > 0 {
		if app.Extensions == nil {
			app.Extensions = map[string]interface{}{}
		}
		app.Extensions[helmAppRoutesFile] = routes
	}
	if policies, err := helmApplicationPoliciesFromChart(chartPath); err == nil && len(policies) > 0 {
		if app.Extensions == nil {
			app.Extensions = map[string]interface{}{}
		}
		app.Extensions[helmAppPoliciesFile] = policies
	}
	if resources, err := helmApplicationCanonicalRawResourcesFromChart(chartPath); err == nil && len(resources) > 0 {
		if app.Extensions == nil {
			app.Extensions = map[string]interface{}{}
		}
		app.Extensions[composeCanonicalRawResourcesExtension] = resources
	}
	if resources, err := helmApplicationKubernetesRawResourcesFromChart(chartPath); err == nil && len(resources) > 0 {
		if app.Extensions == nil {
			app.Extensions = map[string]interface{}{}
		}
		app.Extensions[composeKubernetesRawResourcesExtension] = resources
	}
	rehydrateComposeApplicationExtensions(app)
	syncPortableApplicationState(app)

	// Store chart metadata
	if app.Extensions == nil {
		app.Extensions = make(map[string]interface{})
	}
	app.Extensions["chart"] = chartMeta
	app.Extensions[helmRawFilesExtension] = helmRawFilesFromChart(chrt)
	app.AttachCanonical()
	app.Canonical.AddResource(ResourceKindRaw, PlatformHelm, "chart", "HelmChart", chartMeta)
	app.Canonical.AddResource(ResourceKindRaw, PlatformHelm, "files", "HelmRawFiles", app.Extensions[helmRawFilesExtension])

	return app, nil
}

func renderHelmChart(chartPath string) (string, error) {
	chrt, err := loader.Load(chartPath)
	if err != nil {
		return "", fmt.Errorf("failed to load Helm chart: %w", err)
	}
	return renderLoadedHelmChart(chrt)
}

func renderLoadedHelmChart(chrt *chart.Chart) (string, error) {
	releaseName := chrt.Name()
	values, err := chartutil.ToRenderValues(chrt, chrt.Values, chartutil.ReleaseOptions{
		Name:      releaseName,
		Namespace: "default",
		Revision:  1,
		IsInstall: true,
	}, chartutil.DefaultCapabilities)
	if err != nil {
		return "", fmt.Errorf("failed to prepare Helm render values: %w", err)
	}
	renderedFiles, err := engine.Render(chrt, values)
	if err != nil {
		return "", fmt.Errorf("failed to render Helm chart: %w", err)
	}
	keys := make([]string, 0, len(renderedFiles))
	for key := range renderedFiles {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var rendered strings.Builder
	for _, key := range keys {
		if shouldSkipRenderedHelmFile(key) {
			continue
		}
		content := strings.TrimSpace(renderedFiles[key])
		if content == "" {
			continue
		}
		if rendered.Len() > 0 {
			rendered.WriteString("\n---\n")
		}
		rendered.WriteString(content)
	}
	return rendered.String(), nil
}

func helmRawFilesFromChart(chrt *chart.Chart) map[string]string {
	files := map[string]string{}
	if chrt == nil {
		return files
	}
	helmCollectRawFiles(files, chrt, "")
	if chrt.Lock != nil {
		if data, err := yaml.Marshal(chrt.Lock); err == nil {
			files["Chart.lock"] = base64.StdEncoding.EncodeToString(data)
		}
	}
	return files
}

func helmCollectRawFiles(files map[string]string, chrt *chart.Chart, prefix string) {
	if chrt == nil {
		return
	}
	for _, file := range chrt.Raw {
		if file == nil || file.Name == "" {
			continue
		}
		name := filepath.ToSlash(filepath.Clean(filepath.Join(prefix, file.Name)))
		if !isSafeHelmChartFileName(name) {
			continue
		}
		files[name] = base64.StdEncoding.EncodeToString(file.Data)
	}
	for _, dep := range chrt.Dependencies() {
		helmCollectRawFiles(files, dep, filepath.Join(prefix, "charts", dep.Name()))
	}
}

func shouldSkipRenderedHelmFile(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, "_") || strings.EqualFold(base, "NOTES.txt")
}

func mergeRenderedHelmApp(app, rendered *Application) {
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	if rendered != nil {
		if rendered.Namespace != "" {
			app.Namespace = rendered.Namespace
		}
		if len(rendered.KubernetesOpaqueManifests) > 0 {
			app.KubernetesOpaqueManifests = cloneKubernetesOpaqueManifestSpecs(rendered.KubernetesOpaqueManifests)
		}
		for key, value := range rendered.Extensions {
			app.Extensions[key] = value
		}
	}
	for name, service := range rendered.Services {
		service.Platform = PlatformHelm
		app.Services[name] = service
	}
	if app.Networks == nil {
		app.Networks = make(map[string]*Network)
	}
	for name, network := range rendered.Networks {
		app.Networks[name] = network
	}
	if app.Volumes == nil {
		app.Volumes = make(map[string]*Volume)
	}
	for name, volume := range rendered.Volumes {
		app.Volumes[name] = volume
	}
	if app.Configs == nil {
		app.Configs = make(map[string]*Config)
	}
	for name, config := range rendered.Configs {
		app.Configs[name] = config
	}
	if app.Secrets == nil {
		app.Secrets = make(map[string]*Secret)
	}
	for name, secret := range rendered.Secrets {
		app.Secrets[name] = secret
	}
	if len(rendered.Models) > 0 {
		if app.Models == nil {
			app.Models = map[string]*ComposeModel{}
		}
		for name, model := range rendered.Models {
			if _, exists := app.Models[name]; !exists {
				app.Models[name] = cloneComposeModel(model)
			}
		}
	}
}

// SerializeHelmChart generates a Helm chart from an Application
func SerializeHelmChart(app *Application, chartPath string) error {
	emitApp := cloneApplication(app)
	syncPortableApplicationState(emitApp)
	if err := restoreHelmRawChart(emitApp, chartPath); err == nil {
		if err := writeHelmApplicationExtensions(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmServiceExtensions(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationName(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationModels(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationIncludes(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationNetworks(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationVolumes(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationConfigs(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationSecrets(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationRoutes(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmApplicationPolicies(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmCanonicalRawResources(emitApp, chartPath); err != nil {
			return err
		}
		if err := writeHelmKubernetesRawResources(emitApp, chartPath); err != nil {
			return err
		}
		return nil
	} else if !isNoHelmRawFilesError(err) {
		return err
	}

	// Create chart directory
	if err := os.MkdirAll(chartPath, 0755); err != nil {
		return fmt.Errorf("failed to create chart directory: %w", err)
	}

	// Create Chart.yaml
	chartYaml := map[string]interface{}{
		"apiVersion":  "v2",
		"name":        filepath.Base(chartPath),
		"description": fmt.Sprintf("Generated Helm chart for %s", emitApp.Platform),
		"type":        "application",
		"version":     "0.1.0",
		"appVersion":  "1.0.0",
	}

	if chartMeta, ok := emitApp.Extensions["chart"].(map[string]interface{}); ok {
		// Merge with existing metadata
		for k, v := range chartMeta {
			chartYaml[k] = v
		}
	}
	if emitApp.Version != "" {
		chartYaml["appVersion"] = emitApp.Version
	}

	chartData, err := yaml.Marshal(chartYaml)
	if err != nil {
		return fmt.Errorf("failed to marshal Chart.yaml: %w", err)
	}

	if err := os.WriteFile(filepath.Join(chartPath, "Chart.yaml"), chartData, 0644); err != nil {
		return fmt.Errorf("failed to write Chart.yaml: %w", err)
	}
	if err := writeHelmValuesFile(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationExtensions(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmServiceExtensions(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationName(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationModels(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationIncludes(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationNetworks(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationVolumes(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationConfigs(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationSecrets(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationRoutes(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmApplicationPolicies(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmCanonicalRawResources(emitApp, chartPath); err != nil {
		return err
	}
	if err := writeHelmKubernetesRawResources(emitApp, chartPath); err != nil {
		return err
	}

	// Create templates directory
	templatesDir := filepath.Join(chartPath, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Generate Kubernetes YAML files
	k8sContent, err := SerializeKubernetesYAML(app)
	if err != nil {
		return fmt.Errorf("failed to serialize Kubernetes YAML: %w", err)
	}

	// Split into individual resources and create template files
	documents := strings.Split(k8sContent, "\n---\n")
	for i, doc := range documents {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		resourceType, resourceName := helmResourceTypeAndName(doc)

		if resourceType == "" {
			resourceType = "resource"
		}
		if resourceName == "" {
			resourceName = fmt.Sprintf("%d", i)
		}
		if resourceType == "Deployment" {
			doc = helmTemplateDeploymentDoc(doc, resourceName, emitApp.Services[resourceName], emitApp)
		}

		filename := fmt.Sprintf("%s-%s.yaml", strings.ToLower(resourceType), resourceName)
		filepath := filepath.Join(templatesDir, filename)

		if err := os.WriteFile(filepath, []byte(doc), 0644); err != nil {
			return fmt.Errorf("failed to write template %s: %w", filename, err)
		}
	}

	return nil
}

func helmResourceTypeAndName(doc string) (string, string) {
	var resource map[string]interface{}
	if err := yaml.Unmarshal([]byte(doc), &resource); err != nil {
		return "", ""
	}
	resourceType := toString(resource["kind"])
	resourceName := ""
	if metadata, ok := asMap(resource["metadata"]); ok {
		resourceName = toString(metadata["name"])
	}
	return resourceType, resourceName
}

func writeHelmValuesFile(app *Application, chartPath string) error {
	values := helmValuesForApplication(app)
	if len(values) == 0 {
		return nil
	}
	data, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to marshal Helm values.yaml: %w", err)
	}
	return os.WriteFile(filepath.Join(chartPath, "values.yaml"), data, 0644)
}

func helmValuesForApplication(app *Application) map[string]interface{} {
	if app == nil || len(app.Services) == 0 {
		return nil
	}
	services := map[string]interface{}{}
	names := make([]string, 0, len(app.Services))
	for name, service := range app.Services {
		if service == nil {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		service := app.Services[name]
		replicas := service.Replicas
		if replicas <= 0 {
			replicas = 1
		}
		imagePullSecrets := append([]string{}, service.ImagePullSecrets...)
		var automount *bool
		if service.AutomountServiceAccountToken != nil {
			value := *service.AutomountServiceAccountToken
			automount = &value
		}
		var enableServiceLinks *bool
		if service.EnableServiceLinks != nil {
			value := *service.EnableServiceLinks
			enableServiceLinks = &value
		}
		var hostPID *bool
		if service.HostPID != nil {
			value := *service.HostPID
			hostPID = &value
		}
		var hostIPC *bool
		if service.HostIPC != nil {
			value := *service.HostIPC
			hostIPC = &value
		}
		var hostUsers *bool
		if service.HostUsers != nil {
			value := *service.HostUsers
			hostUsers = &value
		}
		var shareProcessNamespace *bool
		if service.ShareProcessNamespace != nil {
			value := *service.ShareProcessNamespace
			shareProcessNamespace = &value
		}
		var setHostnameAsFQDN *bool
		if service.SetHostnameAsFQDN != nil {
			value := *service.SetHostnameAsFQDN
			setHostnameAsFQDN = &value
		}
		var activeDeadlineSeconds *int64
		if service.ActiveDeadlineSeconds != nil {
			value := *service.ActiveDeadlineSeconds
			activeDeadlineSeconds = &value
		}
		terminationGracePeriodSeconds := durationSeconds(service.StopGracePeriod)
		entrypoint := append([]string{}, service.Entrypoint...)
		command := append([]string{}, service.Command...)
		var runAsUser *int64
		if user := strings.TrimSpace(service.User); user != "" {
			if parsed := parseInt(user); parsed > 0 {
				value := int64(parsed)
				runAsUser = &value
			}
		}
		var runAsGroup *int64
		if group := strings.TrimSpace(service.Group); group != "" {
			if parsed := parseInt(group); parsed > 0 {
				value := int64(parsed)
				runAsGroup = &value
			}
		}
		var fsGroup *int64
		if service.FSGroup != nil {
			value := *service.FSGroup
			fsGroup = &value
		}
		var selinuxOptions map[string]interface{}
		if serialized := serializeKubernetesSELinuxOptions(service.SELinuxOptions); len(serialized) > 0 {
			selinuxOptions = serialized
		}
		var windowsOptions map[string]interface{}
		if serialized := serializeKubernetesWindowsSecurityContextOptions(service.WindowsOptions); len(serialized) > 0 {
			windowsOptions = serialized
		}
		var fsGroupChangePolicy string
		if service.FSGroupChangePolicy != "" {
			fsGroupChangePolicy = service.FSGroupChangePolicy
		}
		var runAsNonRoot *bool
		if service.RunAsNonRoot != nil {
			value := *service.RunAsNonRoot
			runAsNonRoot = &value
		}
		var supplementalGroups []int64
		if len(service.SupplementalGroups) > 0 {
			supplementalGroups = append([]int64{}, service.SupplementalGroups...)
		}
		var supplementalGroupsPolicy string
		if service.SupplementalGroupsPolicy != "" {
			supplementalGroupsPolicy = service.SupplementalGroupsPolicy
		}
		var seccompProfileType string
		var seccompProfileLocalhost string
		if service.SeccompProfile != nil {
			seccompProfileType = service.SeccompProfile.Type
			seccompProfileLocalhost = service.SeccompProfile.LocalhostProfile
		}
		var initValue *bool
		if service.Init != nil {
			value := *service.Init
			initValue = &value
		}
		var sysctls []map[string]interface{}
		if len(service.Sysctls) > 0 {
			sysctls = kubernetesSysctlsFromMap(service.Sysctls)
		}
		affinity := map[string]interface{}{}
		if len(service.Affinity) > 0 {
			affinity = cloneMap(service.Affinity)
		}
		var readinessGates []map[string]interface{}
		if len(service.ReadinessGates) > 0 {
			readinessGates = cloneMapSlice(service.ReadinessGates)
		}
		var schedulingGates []map[string]interface{}
		if len(service.SchedulingGates) > 0 {
			schedulingGates = cloneMapSlice(service.SchedulingGates)
		}
		var topologySpreadConstraints []map[string]interface{}
		if len(service.TopologySpreadConstraints) > 0 {
			topologySpreadConstraints = cloneMapSlice(service.TopologySpreadConstraints)
		}
		var readinessProbe map[string]interface{}
		if probe := serializeKubernetesProbe(service.HealthCheck); len(probe) > 0 {
			readinessProbe = probe
		}
		var livenessProbe map[string]interface{}
		if probe := serializeKubernetesProbe(service.HealthCheck); len(probe) > 0 {
			livenessProbe = cloneMap(probe)
		}
		var startupProbe map[string]interface{}
		if probe := serializeKubernetesProbeExtension(service.StartupProbe); len(probe) > 0 {
			startupProbe = probe
		}
		var portableHealthcheck string
		if service.HealthCheck != nil {
			if raw, err := json.Marshal(service.HealthCheck); err == nil {
				portableHealthcheck = string(raw)
			}
		}
		var portableStartupProbe string
		if service.StartupProbe != nil {
			if raw, err := json.Marshal(service.StartupProbe); err == nil {
				portableStartupProbe = string(raw)
			}
		}
		var portableUlimits string
		if !isEmptyUlimits(service.Ulimits) {
			if raw, err := json.Marshal(service.Ulimits); err == nil {
				portableUlimits = string(raw)
			}
		}
		var portableLogging string
		if service.LogDriver != "" || len(service.LogOpt) > 0 || len(service.LogExtensions) > 0 {
			logging := map[string]interface{}{}
			if service.LogDriver != "" {
				logging["driver"] = service.LogDriver
			}
			if len(service.LogOpt) > 0 {
				logging["options"] = copyStringMap(service.LogOpt)
			}
			if len(service.LogExtensions) > 0 {
				logging["extensions"] = copyStringInterfaceMap(service.LogExtensions)
			}
			if raw, err := json.Marshal(logging); err == nil {
				portableLogging = string(raw)
			}
		}
		var portableComposeCompat string
		if !isEmptyComposeCompat(service.ComposeCompat) {
			if raw, err := json.Marshal(service.ComposeCompat); err == nil {
				portableComposeCompat = string(raw)
			}
		}
		var portableNomadConnect string
		if connect := nomadConnectSpecToMap(service.Connect); len(connect) > 0 {
			if raw, err := json.Marshal(connect); err == nil {
				portableNomadConnect = string(raw)
			}
		}
		var portableNomadSpread string
		if spreads := nomadSpreadSpecsToMap(service.Spreads); len(spreads) > 0 {
			if raw, err := json.Marshal(spreads); err == nil {
				portableNomadSpread = string(raw)
			}
		}
		var portableNomadRestart string
		if restart := nomadRestartBlockForService(service); len(restart) > 0 {
			if raw, err := json.Marshal(restart); err == nil {
				portableNomadRestart = string(raw)
			}
		}
		var portableNomadUpdate string
		if update := nomadSchedulerExtensionMap(service, nomadUpdateExtensionKey, "x-nomad-update"); len(update) > 0 {
			if raw, err := json.Marshal(update); err == nil {
				portableNomadUpdate = string(raw)
			}
		}
		var portableNomadMigrate string
		if migrate := nomadSchedulerExtensionMap(service, nomadMigrateExtensionKey, "x-nomad-migrate"); len(migrate) > 0 {
			if raw, err := json.Marshal(migrate); err == nil {
				portableNomadMigrate = string(raw)
			}
		}
		var portableNomadReschedule string
		if reschedule := nomadSchedulerExtensionMap(service, nomadRescheduleExtensionKey, "x-nomad-reschedule"); len(reschedule) > 0 {
			if raw, err := json.Marshal(reschedule); err == nil {
				portableNomadReschedule = string(raw)
			}
		}
		var portableMesh string
		if mesh := meshSpecToMap(app.Mesh); len(mesh) > 0 {
			if raw, err := json.Marshal(mesh); err == nil {
				portableMesh = string(raw)
			}
		}
		var dnsConfig map[string]interface{}
		if len(service.DNS) > 0 || len(service.DNSSearch) > 0 || len(service.DNSOptions) > 0 {
			dnsConfig = map[string]interface{}{}
			if len(service.DNS) > 0 {
				dnsConfig["nameservers"] = append([]string{}, service.DNS...)
			}
			if len(service.DNSSearch) > 0 {
				dnsConfig["searches"] = append([]string{}, service.DNSSearch...)
			}
			if len(service.DNSOptions) > 0 {
				var options []map[string]interface{}
				for _, option := range service.DNSOptions {
					name, value := splitDNSOption(option)
					if name == "" {
						continue
					}
					item := map[string]interface{}{"name": name}
					if value != "" {
						item["value"] = value
					}
					options = append(options, item)
				}
				if len(options) > 0 {
					dnsConfig["options"] = options
				}
			}
		}
		hostAliases := kubernetesHostAliasesFromService(service)
		podRestartPolicy := service.PodRestartPolicy
		var profiles []string
		if len(service.Profiles) > 0 {
			profiles = append([]string{}, service.Profiles...)
		}
		var dependenciesJSON string
		if dependencies := serviceDependencies(service); len(dependencies) > 0 {
			if raw, err := json.Marshal(dependencies); err == nil {
				dependenciesJSON = string(raw)
			}
		}
		var portableDevelop string
		if !isEmptyDevelopConfig(service.Develop) {
			if raw, err := json.Marshal(service.Develop); err == nil {
				portableDevelop = string(raw)
			}
		}
		var lifecycle map[string]interface{}
		if serialized := serializeKubernetesLifecycle(service.Lifecycle); len(serialized) > 0 {
			lifecycle = serialized
		}
		var portableLifecycle string
		if !isEmptyLifecycleHooks(service.Lifecycle) {
			if raw, err := json.Marshal(service.Lifecycle); err == nil {
				portableLifecycle = string(raw)
			}
		}
		var deploySpec string
		var deployMode string
		var deployEndpointMode string
		var deployLabels string
		var deployResources string
		var deployJob string
		var deployPlacementConstraints []string
		var deployPlacementPreferences []string
		var deployMaxReplicasPerNode int
		var deployUpdateParallelism *int
		var deployUpdateDelay string
		var deployUpdateMonitor string
		var deployUpdateFailureRate string
		var deployUpdateOrder string
		var deployUpdateOnFailure string
		var deployUpdateHealthCheck string
		var deployUpdateMinHealthyTime string
		var deployUpdateHealthyDeadline string
		var deployUpdateProgressDeadline string
		var deployUpdateAutoRevert bool
		var deployUpdateAutoPromote bool
		var deployUpdateCanary int
		var deployUpdateStagger string
		var deployRollbackParallelism *int
		var deployRollbackDelay string
		var deployRollbackMonitor string
		var deployRollbackFailureRate string
		var deployRollbackOrder string
		var deployRollbackOnFailure string
		var deployRollbackHealthCheck string
		var deployRollbackMinHealthyTime string
		var deployRollbackHealthyDeadline string
		var deployRollbackProgressDeadline string
		var deployRollbackAutoRevert bool
		var deployRollbackAutoPromote bool
		var deployRollbackCanary int
		var deployRollbackStagger string
		var deployRestartCondition string
		var deployRestartDelay string
		var deployRestartAttempts int
		var deployRestartWindow string
		var failoverSpec string
		if service.Failover != nil {
			if raw, err := json.Marshal(service.Failover); err == nil {
				failoverSpec = string(raw)
			}
		}
		if deploy := cloneDeploySpec(service.Deploy); deploy != nil && !isEmptyDeploySpec(deploy) {
			deployMode = deploy.Mode
			deployEndpointMode = deploy.EndpointMode
			if len(deploy.Labels) > 0 {
				if raw, err := json.Marshal(deploy.Labels); err == nil {
					deployLabels = string(raw)
				}
			}
			if deploy.Resources != nil && !isEmptyResourceSpec(deploy.Resources) {
				if raw, err := json.Marshal(deploy.Resources); err == nil {
					deployResources = string(raw)
				}
			}
			if deploy.Job != nil && !isEmptySwarmJobSpec(deploy.Job) {
				if raw, err := json.Marshal(deploy.Job); err == nil {
					deployJob = string(raw)
				}
			}
			if deploy.Placement != nil {
				deployPlacementConstraints = append([]string{}, deploy.Placement.Constraints...)
				deployPlacementPreferences = append([]string{}, deploy.Placement.Preferences...)
				deployMaxReplicasPerNode = deploy.Placement.MaxReplicasPerNode
			}
			if deploy.UpdateConfig != nil {
				value := deploy.UpdateConfig.Parallelism
				deployUpdateParallelism = &value
				deployUpdateDelay = deploy.UpdateConfig.Delay
				deployUpdateMonitor = deploy.UpdateConfig.Monitor
				deployUpdateFailureRate = deploy.UpdateConfig.MaxFailureRatio
				deployUpdateOrder = deploy.UpdateConfig.Order
				deployUpdateOnFailure = deploy.UpdateConfig.OnFailure
				deployUpdateHealthCheck = deploy.UpdateConfig.HealthCheck
				deployUpdateMinHealthyTime = deploy.UpdateConfig.MinHealthyTime
				deployUpdateHealthyDeadline = deploy.UpdateConfig.HealthyDeadline
				deployUpdateProgressDeadline = deploy.UpdateConfig.ProgressDeadline
				deployUpdateAutoRevert = deploy.UpdateConfig.AutoRevert
				deployUpdateAutoPromote = deploy.UpdateConfig.AutoPromote
				deployUpdateCanary = deploy.UpdateConfig.Canary
				deployUpdateStagger = deploy.UpdateConfig.Stagger
			}
			if deploy.RollbackConfig != nil {
				value := deploy.RollbackConfig.Parallelism
				deployRollbackParallelism = &value
				deployRollbackDelay = deploy.RollbackConfig.Delay
				deployRollbackMonitor = deploy.RollbackConfig.Monitor
				deployRollbackFailureRate = deploy.RollbackConfig.MaxFailureRatio
				deployRollbackOrder = deploy.RollbackConfig.Order
				deployRollbackOnFailure = deploy.RollbackConfig.OnFailure
				deployRollbackHealthCheck = deploy.RollbackConfig.HealthCheck
				deployRollbackMinHealthyTime = deploy.RollbackConfig.MinHealthyTime
				deployRollbackHealthyDeadline = deploy.RollbackConfig.HealthyDeadline
				deployRollbackProgressDeadline = deploy.RollbackConfig.ProgressDeadline
				deployRollbackAutoRevert = deploy.RollbackConfig.AutoRevert
				deployRollbackAutoPromote = deploy.RollbackConfig.AutoPromote
				deployRollbackCanary = deploy.RollbackConfig.Canary
				deployRollbackStagger = deploy.RollbackConfig.Stagger
			}
			if deploy.RestartPolicy != nil {
				deployRestartCondition = deploy.RestartPolicy.Condition
				deployRestartDelay = deploy.RestartPolicy.Delay
				deployRestartAttempts = deploy.RestartPolicy.MaxAttempts
				deployRestartWindow = deploy.RestartPolicy.Window
			}
			if raw, err := json.Marshal(deploy); err == nil {
				deploySpec = string(raw)
			}
		}
		services[name] = map[string]interface{}{
			"image":                             service.Image,
			"replicas":                          replicas,
			"image_pull_policy":                 service.ImagePullPolicy,
			"service_account":                   service.ServiceAccountName,
			"automount_service_account_token":   automount,
			"image_pull_secrets":                imagePullSecrets,
			"runtime_class_name":                service.RuntimeClassName,
			"enable_service_links":              enableServiceLinks,
			"priority_class_name":               service.PriorityClassName,
			"scheduler_name":                    service.SchedulerName,
			"node_name":                         service.NodeName,
			"hostname":                          service.Hostname,
			"node_selector":                     copyStringMap(service.NodeSelector),
			"subdomain":                         service.Subdomain,
			"os_name":                           service.OSName,
			"dns_policy":                        service.DNSPolicy,
			"host_network":                      service.HostNetwork,
			"host_pid":                          hostPID,
			"host_ipc":                          hostIPC,
			"host_users":                        hostUsers,
			"share_process_namespace":           shareProcessNamespace,
			"set_hostname_as_fqdn":              setHostnameAsFQDN,
			"active_deadline_seconds":           activeDeadlineSeconds,
			"termination_grace_period_seconds":  terminationGracePeriodSeconds,
			"entrypoint":                        entrypoint,
			"command":                           command,
			"working_dir":                       service.WorkingDir,
			"tty":                               service.Tty,
			"stdin_open":                        service.StdinOpen,
			"termination_message_path":          service.TerminationMessagePath,
			"termination_message_policy":        service.TerminationMessagePolicy,
			"privileged":                        service.Privileged,
			"read_only_root_filesystem":         service.ReadOnlyRootFS,
			"allow_privilege_escalation":        service.AllowPrivilegeEscalation,
			"proc_mount":                        service.ProcMount,
			"run_as_user":                       runAsUser,
			"run_as_group":                      runAsGroup,
			"fs_group":                          fsGroup,
			"se_linux_options":                  selinuxOptions,
			"windows_options":                   windowsOptions,
			"fs_group_change_policy":            fsGroupChangePolicy,
			"run_as_non_root":                   runAsNonRoot,
			"supplemental_groups":               supplementalGroups,
			"supplemental_groups_policy":        supplementalGroupsPolicy,
			"cap_add":                           append([]string{}, service.CapAdd...),
			"cap_drop":                          append([]string{}, service.CapDrop...),
			"security_opt":                      append([]string{}, service.SecurityOpt...),
			"init":                              initValue,
			"stop_signal":                       service.StopSignal,
			"pid_mode":                          service.PIDMode,
			"ipc_mode":                          service.IPCMode,
			"pids_limit":                        service.PidsLimit,
			"shm_size":                          service.ShmSize,
			"runtime":                           service.Runtime,
			"userns_mode":                       service.UserNSMode,
			"group_add":                         append([]string{}, service.GroupAdd...),
			"sysctls":                           sysctls,
			"seccomp_profile_type":              seccompProfileType,
			"seccomp_profile_localhost_profile": seccompProfileLocalhost,
			"affinity":                          affinity,
			"readiness_gates":                   readinessGates,
			"init_containers":                   cloneMapSlice(service.InitContainers),
			"resource_claims":                   cloneMapSlice(service.ResourceClaims),
			"ephemeral_containers":              cloneMapSlice(service.EphemeralContainers),
			"scheduling_gates":                  schedulingGates,
			"tolerations":                       serializeKubernetesTolerationsNative(service.Tolerations),
			"topology_spread_constraints":       topologySpreadConstraints,
			"profiles":                          profiles,
			"dependencies":                      dependenciesJSON,
			"portable_develop":                  portableDevelop,
			"portable_lifecycle":                portableLifecycle,
			"portable_healthcheck":              portableHealthcheck,
			"portable_startup_probe":            portableStartupProbe,
			"portable_ulimits":                  portableUlimits,
			"portable_logging":                  portableLogging,
			"portable_compose_compat":           portableComposeCompat,
			"portable_nomad_connect":            portableNomadConnect,
			"portable_nomad_spread":             portableNomadSpread,
			"portable_nomad_restart":            portableNomadRestart,
			"portable_nomad_update":             portableNomadUpdate,
			"portable_nomad_migrate":            portableNomadMigrate,
			"portable_nomad_reschedule":         portableNomadReschedule,
			"portable_mesh":                     portableMesh,
			"dns_config":                        dnsConfig,
			"host_aliases":                      hostAliases,
			"links":                             append([]string{}, service.Links...),
			"pod_restart_policy":                podRestartPolicy,
			"lifecycle":                         lifecycle,
			"readiness_probe":                   readinessProbe,
			"liveness_probe":                    livenessProbe,
			"startup_probe":                     startupProbe,
			"deploy_spec":                       deploySpec,
			"deploy_mode":                       deployMode,
			"deploy_endpoint_mode":              deployEndpointMode,
			"deploy_labels":                     deployLabels,
			"deploy_resources":                  deployResources,
			"deploy_job":                        deployJob,
			"deploy_placement_constraints":      deployPlacementConstraints,
			"deploy_placement_preferences":      deployPlacementPreferences,
			"deploy_max_replicas_per_node":      deployMaxReplicasPerNode,
			"deploy_update_parallelism":         deployUpdateParallelism,
			"deploy_update_delay":               deployUpdateDelay,
			"deploy_update_monitor":             deployUpdateMonitor,
			"deploy_update_failure_rate":        deployUpdateFailureRate,
			"deploy_update_order":               deployUpdateOrder,
			"deploy_update_on_failure":          deployUpdateOnFailure,
			"deploy_update_health_check":        deployUpdateHealthCheck,
			"deploy_update_min_healthy_time":    deployUpdateMinHealthyTime,
			"deploy_update_healthy_deadline":    deployUpdateHealthyDeadline,
			"deploy_update_progress_deadline":   deployUpdateProgressDeadline,
			"deploy_update_auto_revert":         deployUpdateAutoRevert,
			"deploy_update_auto_promote":        deployUpdateAutoPromote,
			"deploy_update_canary":              deployUpdateCanary,
			"deploy_update_stagger":             deployUpdateStagger,
			"deploy_rollback_parallelism":       deployRollbackParallelism,
			"deploy_rollback_delay":             deployRollbackDelay,
			"deploy_rollback_monitor":           deployRollbackMonitor,
			"deploy_rollback_failure_rate":      deployRollbackFailureRate,
			"deploy_rollback_order":             deployRollbackOrder,
			"deploy_rollback_on_failure":        deployRollbackOnFailure,
			"deploy_rollback_health_check":      deployRollbackHealthCheck,
			"deploy_rollback_min_healthy_time":  deployRollbackMinHealthyTime,
			"deploy_rollback_healthy_deadline":  deployRollbackHealthyDeadline,
			"deploy_rollback_progress_deadline": deployRollbackProgressDeadline,
			"deploy_rollback_auto_revert":       deployRollbackAutoRevert,
			"deploy_rollback_auto_promote":      deployRollbackAutoPromote,
			"deploy_rollback_canary":            deployRollbackCanary,
			"deploy_rollback_stagger":           deployRollbackStagger,
			"deploy_restart_condition":          deployRestartCondition,
			"deploy_restart_delay":              deployRestartDelay,
			"deploy_restart_attempts":           deployRestartAttempts,
			"deploy_restart_window":             deployRestartWindow,
			"failover_spec":                     failoverSpec,
		}
	}
	if len(services) == 0 {
		return nil
	}
	return map[string]interface{}{"services": services}
}

func helmTemplateDeploymentDoc(doc, serviceName string, service *Service, app *Application) string {
	if strings.TrimSpace(doc) == "" || service == nil {
		return doc
	}
	templated := doc
	if service.Image != "" {
		templated = strings.ReplaceAll(templated, "image: "+service.Image, fmt.Sprintf("image: {{ index .Values.services %q \"image\" | quote }}", serviceName))
	}
	replicas := service.Replicas
	if replicas <= 0 {
		replicas = 1
	}
	templated = strings.ReplaceAll(templated, fmt.Sprintf("replicas: %d", replicas), fmt.Sprintf("replicas: {{ index .Values.services %q \"replicas\" }}", serviceName))
	if service.ImagePullPolicy != "" {
		templated = strings.ReplaceAll(templated, "imagePullPolicy: "+service.ImagePullPolicy, fmt.Sprintf("imagePullPolicy: {{ index .Values.services %q \"image_pull_policy\" | quote }}", serviceName))
	}
	if service.ServiceAccountName != "" {
		templated = strings.ReplaceAll(templated, "serviceAccountName: "+service.ServiceAccountName, fmt.Sprintf("serviceAccountName: {{ index .Values.services %q \"service_account\" | quote }}", serviceName))
	}
	if service.RuntimeClassName != "" {
		templated = strings.ReplaceAll(templated, "runtimeClassName: "+service.RuntimeClassName, fmt.Sprintf("runtimeClassName: {{ index .Values.services %q \"runtime_class_name\" | quote }}", serviceName))
	}
	if service.PriorityClassName != "" {
		templated = strings.ReplaceAll(templated, "priorityClassName: "+service.PriorityClassName, fmt.Sprintf("priorityClassName: {{ index .Values.services %q \"priority_class_name\" | quote }}", serviceName))
	}
	if service.SchedulerName != "" {
		templated = strings.ReplaceAll(templated, "schedulerName: "+service.SchedulerName, fmt.Sprintf("schedulerName: {{ index .Values.services %q \"scheduler_name\" | quote }}", serviceName))
	}
	if service.NodeName != "" {
		templated = strings.ReplaceAll(templated, "nodeName: "+service.NodeName, fmt.Sprintf("nodeName: {{ index .Values.services %q \"node_name\" | quote }}", serviceName))
	}
	if service.Hostname != "" {
		templated = strings.ReplaceAll(templated, "hostname: "+service.Hostname, fmt.Sprintf("hostname: {{ index .Values.services %q \"hostname\" | quote }}", serviceName))
	}
	if service.Subdomain != "" {
		templated = strings.ReplaceAll(templated, "subdomain: "+service.Subdomain, fmt.Sprintf("subdomain: {{ index .Values.services %q \"subdomain\" | quote }}", serviceName))
	}
	if service.OSName != "" {
		templated = replaceHelmIndentedBlock(templated, "os:", func(indent string) string {
			childIndent := indent + "  "
			var b strings.Builder
			b.WriteString(fmt.Sprintf("%sos:\n", indent))
			b.WriteString(fmt.Sprintf("%sname: {{ index .Values.services %q \"os_name\" | quote }}\n", childIndent, serviceName))
			return b.String()
		})
	}
	if service.DNSPolicy != "" {
		dnsTemplate := fmt.Sprintf("dnsPolicy: {{ index .Values.services %q \"dns_policy\" | quote }}", serviceName)
		templated = strings.ReplaceAll(templated, "dnsPolicy: "+service.DNSPolicy, dnsTemplate)
		if len(service.NodeSelector) > 0 {
			nodeSelectorBlock := fmt.Sprintf("\n            {{- with index $.Values.services %q \"node_selector\" }}\n            nodeSelector:\n            {{ toYaml . | nindent 14 }}\n            {{- end }}", serviceName)
			templated = strings.Replace(templated, dnsTemplate, dnsTemplate+nodeSelectorBlock, 1)
		}
	} else if len(service.NodeSelector) > 0 {
		templated = replaceHelmIndentedBlock(templated, "nodeName:", func(indent string) string {
			childIndent := indent + "  "
			var b strings.Builder
			b.WriteString(fmt.Sprintf("%snodeName: {{ index .Values.services %q \"node_name\" | quote }}\n", indent, serviceName))
			b.WriteString(fmt.Sprintf("%s{{- with index $.Values.services %q \"node_selector\" }}\n", indent, serviceName))
			b.WriteString(fmt.Sprintf("%snodeSelector:\n", indent))
			b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
			return b.String()
		})
	}
	if service.PodRestartPolicy != "" {
		templated = strings.ReplaceAll(templated, "restartPolicy: "+service.PodRestartPolicy, fmt.Sprintf("restartPolicy: {{ index .Values.services %q \"pod_restart_policy\" | quote }}", serviceName))
	}
	if service.HostNetworkSet || service.HostNetwork {
		templated = strings.ReplaceAll(templated, "hostNetwork: true", fmt.Sprintf("hostNetwork: {{ index .Values.services %q \"host_network\" }}", serviceName))
		templated = strings.ReplaceAll(templated, "hostNetwork: false", fmt.Sprintf("hostNetwork: {{ index .Values.services %q \"host_network\" }}", serviceName))
	}
	if service.HostPID != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("hostPID: %t", *service.HostPID), fmt.Sprintf("hostPID: {{ index .Values.services %q \"host_pid\" }}", serviceName))
	}
	if service.HostIPC != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("hostIPC: %t", *service.HostIPC), fmt.Sprintf("hostIPC: {{ index .Values.services %q \"host_ipc\" }}", serviceName))
	}
	if service.HostUsers != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("hostUsers: %t", *service.HostUsers), fmt.Sprintf("hostUsers: {{ index .Values.services %q \"host_users\" }}", serviceName))
	}
	if service.ShareProcessNamespace != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("shareProcessNamespace: %t", *service.ShareProcessNamespace), fmt.Sprintf("shareProcessNamespace: {{ index .Values.services %q \"share_process_namespace\" }}", serviceName))
	}
	if service.SetHostnameAsFQDN != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("setHostnameAsFQDN: %t", *service.SetHostnameAsFQDN), fmt.Sprintf("setHostnameAsFQDN: {{ index .Values.services %q \"set_hostname_as_fqdn\" }}", serviceName))
	}
	if service.ActiveDeadlineSeconds != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("activeDeadlineSeconds: %d", *service.ActiveDeadlineSeconds), fmt.Sprintf("activeDeadlineSeconds: {{ index .Values.services %q \"active_deadline_seconds\" }}", serviceName))
	}
	if seconds := durationSeconds(service.StopGracePeriod); seconds > 0 {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("terminationGracePeriodSeconds: %d", seconds), fmt.Sprintf("terminationGracePeriodSeconds: {{ index .Values.services %q \"termination_grace_period_seconds\" }}", serviceName))
	}
	if len(service.Entrypoint) > 0 {
		templated = helmTemplateStringSliceField(templated, serviceName, "", "command", "entrypoint", service.Entrypoint)
	}
	if len(service.Command) > 0 {
		templated = helmTemplateStringSliceField(templated, serviceName, "- ", "args", "command", service.Command)
	}
	if service.WorkingDir != "" {
		templated = strings.ReplaceAll(templated, "workingDir: "+service.WorkingDir, fmt.Sprintf("workingDir: {{ index .Values.services %q \"working_dir\" | quote }}", serviceName))
	}
	if service.TtySet || service.Tty {
		templated = strings.ReplaceAll(templated, "tty: true", fmt.Sprintf("tty: {{ index .Values.services %q \"tty\" }}", serviceName))
		templated = strings.ReplaceAll(templated, "tty: false", fmt.Sprintf("tty: {{ index .Values.services %q \"tty\" }}", serviceName))
	}
	if service.StdinOpenSet || service.StdinOpen {
		templated = strings.ReplaceAll(templated, "stdin: true", fmt.Sprintf("stdin: {{ index .Values.services %q \"stdin_open\" }}", serviceName))
		templated = strings.ReplaceAll(templated, "stdin: false", fmt.Sprintf("stdin: {{ index .Values.services %q \"stdin_open\" }}", serviceName))
	}
	if service.TerminationMessagePath != "" {
		templated = strings.ReplaceAll(templated, "terminationMessagePath: "+service.TerminationMessagePath, fmt.Sprintf("terminationMessagePath: {{ index .Values.services %q \"termination_message_path\" | quote }}", serviceName))
	}
	if service.TerminationMessagePolicy != "" {
		templated = strings.ReplaceAll(templated, "terminationMessagePolicy: "+service.TerminationMessagePolicy, fmt.Sprintf("terminationMessagePolicy: {{ index .Values.services %q \"termination_message_policy\" | quote }}", serviceName))
	}
	if service.PrivilegedSet || service.Privileged {
		templated = strings.ReplaceAll(templated, "privileged: true", fmt.Sprintf("privileged: {{ index .Values.services %q \"privileged\" }}", serviceName))
		templated = strings.ReplaceAll(templated, "privileged: false", fmt.Sprintf("privileged: {{ index .Values.services %q \"privileged\" }}", serviceName))
	}
	if service.ReadOnlyRootFSSet || service.ReadOnlyRootFS {
		templated = strings.ReplaceAll(templated, "readOnlyRootFilesystem: true", fmt.Sprintf("readOnlyRootFilesystem: {{ index .Values.services %q \"read_only_root_filesystem\" }}", serviceName))
		templated = strings.ReplaceAll(templated, "readOnlyRootFilesystem: false", fmt.Sprintf("readOnlyRootFilesystem: {{ index .Values.services %q \"read_only_root_filesystem\" }}", serviceName))
	}
	if service.AllowPrivilegeEscalation != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("allowPrivilegeEscalation: %t", *service.AllowPrivilegeEscalation), fmt.Sprintf("allowPrivilegeEscalation: {{ index .Values.services %q \"allow_privilege_escalation\" }}", serviceName))
	}
	if service.ProcMount != "" {
		templated = strings.ReplaceAll(templated, "procMount: "+service.ProcMount, fmt.Sprintf("procMount: {{ index .Values.services %q \"proc_mount\" | quote }}", serviceName))
	}
	if user := strings.TrimSpace(service.User); user != "" {
		templated = strings.ReplaceAll(templated, "runAsUser: "+user, fmt.Sprintf("runAsUser: {{ index .Values.services %q \"run_as_user\" }}", serviceName))
	}
	if group := strings.TrimSpace(service.Group); group != "" {
		templated = strings.ReplaceAll(templated, "runAsGroup: "+group, fmt.Sprintf("runAsGroup: {{ index .Values.services %q \"run_as_group\" }}", serviceName))
	}
	if service.FSGroup != nil && *service.FSGroup > 0 {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("fsGroup: %d", *service.FSGroup), fmt.Sprintf("fsGroup: {{ index .Values.services %q \"fs_group\" }}", serviceName))
	}
	if service.SELinuxOptions != nil {
		templated = replaceHelmIndentedBlock(templated, "seLinuxOptions:", func(indent string) string {
			childIndent := indent + "  "
			var b strings.Builder
			b.WriteString(fmt.Sprintf("%sseLinuxOptions:\n", indent))
			b.WriteString(fmt.Sprintf("%s{{- with index .Values.services %q \"se_linux_options\" }}\n", indent, serviceName))
			b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
			return b.String()
		})
	}
	if service.WindowsOptions != nil {
		templated = replaceHelmIndentedBlock(templated, "windowsOptions:", func(indent string) string {
			childIndent := indent + "  "
			var b strings.Builder
			b.WriteString(fmt.Sprintf("%swindowsOptions:\n", indent))
			b.WriteString(fmt.Sprintf("%s{{- with index .Values.services %q \"windows_options\" }}\n", indent, serviceName))
			b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
			return b.String()
		})
	}
	if service.FSGroupChangePolicy != "" {
		templated = strings.ReplaceAll(templated, "fsGroupChangePolicy: "+service.FSGroupChangePolicy, fmt.Sprintf("fsGroupChangePolicy: {{ index .Values.services %q \"fs_group_change_policy\" | quote }}", serviceName))
	}
	if service.RunAsNonRoot != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("runAsNonRoot: %t", *service.RunAsNonRoot), fmt.Sprintf("runAsNonRoot: {{ index .Values.services %q \"run_as_non_root\" }}", serviceName))
	}
	if len(service.SupplementalGroups) > 0 {
		templated = helmTemplateIntSliceField(templated, serviceName, "supplementalGroups", "supplemental_groups", service.SupplementalGroups)
	}
	if service.SupplementalGroupsPolicy != "" {
		templated = strings.ReplaceAll(templated, "supplementalGroupsPolicy: "+service.SupplementalGroupsPolicy, fmt.Sprintf("supplementalGroupsPolicy: {{ index .Values.services %q \"supplemental_groups_policy\" | quote }}", serviceName))
	}
	if len(service.CapAdd) > 0 || len(service.CapDrop) > 0 {
		templated = helmTemplateCapabilitiesBlock(templated, serviceName, service.CapAdd, service.CapDrop)
	}
	if len(service.SecurityOpt) > 0 {
		templated = helmTemplateAnnotationListField(templated, serviceName, "bolabaden.dev/runtime-security-opt", "security_opt")
	}
	if service.Init != nil {
		templated = helmTemplateAnnotationBoolField(templated, serviceName, "bolabaden.dev/runtime-init", "init", *service.Init)
	}
	if service.StopSignal != "" {
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/runtime-stop-signal", "stop_signal", service.StopSignal)
	}
	if service.Deploy != nil && !isEmptyDeploySpec(service.Deploy) {
		deploySpecJSON := ""
		if raw, err := json.Marshal(service.Deploy); err == nil {
			deploySpecJSON = string(raw)
		}
		deployLabelsJSON := ""
		if len(service.Deploy.Labels) > 0 {
			if raw, err := json.Marshal(service.Deploy.Labels); err == nil {
				deployLabelsJSON = string(raw)
			}
		}
		deployResourcesJSON := ""
		if service.Deploy.Resources != nil && !isEmptyResourceSpec(service.Deploy.Resources) {
			if raw, err := json.Marshal(service.Deploy.Resources); err == nil {
				deployResourcesJSON = string(raw)
			}
		}
		deployJobJSON := ""
		if service.Deploy.Job != nil && !isEmptySwarmJobSpec(service.Deploy.Job) {
			if raw, err := json.Marshal(service.Deploy.Job); err == nil {
				deployJobJSON = string(raw)
			}
		}
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-spec", "deploy_spec", deploySpecJSON)
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-mode", "deploy_mode", service.Deploy.Mode)
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-endpoint-mode", "deploy_endpoint_mode", service.Deploy.EndpointMode)
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-labels", "deploy_labels", deployLabelsJSON)
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-resources", "deploy_resources", deployResourcesJSON)
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-job", "deploy_job", deployJobJSON)
		templated = helmTemplateAnnotationListField(templated, serviceName, "bolabaden.dev/deploy-placement-constraints", "deploy_placement_constraints")
		templated = helmTemplateAnnotationListField(templated, serviceName, "bolabaden.dev/deploy-placement-preferences", "deploy_placement_preferences")
		if service.Deploy.Placement != nil && service.Deploy.Placement.MaxReplicasPerNode > 0 {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-max-replicas-per-node", "deploy_max_replicas_per_node", fmt.Sprintf("%d", service.Deploy.Placement.MaxReplicasPerNode))
		}
		if service.Deploy.UpdateConfig != nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-parallelism", "deploy_update_parallelism", fmt.Sprintf("%d", service.Deploy.UpdateConfig.Parallelism))
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-delay", "deploy_update_delay", service.Deploy.UpdateConfig.Delay)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-monitor", "deploy_update_monitor", service.Deploy.UpdateConfig.Monitor)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-max-failure-ratio", "deploy_update_failure_rate", service.Deploy.UpdateConfig.MaxFailureRatio)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-order", "deploy_update_order", service.Deploy.UpdateConfig.Order)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-failure-action", "deploy_update_on_failure", service.Deploy.UpdateConfig.OnFailure)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-health-check", "deploy_update_health_check", service.Deploy.UpdateConfig.HealthCheck)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-min-healthy-time", "deploy_update_min_healthy_time", service.Deploy.UpdateConfig.MinHealthyTime)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-healthy-deadline", "deploy_update_healthy_deadline", service.Deploy.UpdateConfig.HealthyDeadline)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-progress-deadline", "deploy_update_progress_deadline", service.Deploy.UpdateConfig.ProgressDeadline)
			if service.Deploy.UpdateConfig.AutoRevertSet || service.Deploy.UpdateConfig.AutoRevert {
				templated = helmTemplateAnnotationBoolField(templated, serviceName, "bolabaden.dev/deploy-update-auto-revert", "deploy_update_auto_revert", service.Deploy.UpdateConfig.AutoRevert)
			}
			if service.Deploy.UpdateConfig.AutoPromoteSet || service.Deploy.UpdateConfig.AutoPromote {
				templated = helmTemplateAnnotationBoolField(templated, serviceName, "bolabaden.dev/deploy-update-auto-promote", "deploy_update_auto_promote", service.Deploy.UpdateConfig.AutoPromote)
			}
			if service.Deploy.UpdateConfig.CanarySet || service.Deploy.UpdateConfig.Canary > 0 {
				templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-canary", "deploy_update_canary", fmt.Sprintf("%d", service.Deploy.UpdateConfig.Canary))
			}
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-update-stagger", "deploy_update_stagger", service.Deploy.UpdateConfig.Stagger)
		}
		if service.Deploy.RollbackConfig != nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-parallelism", "deploy_rollback_parallelism", fmt.Sprintf("%d", service.Deploy.RollbackConfig.Parallelism))
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-delay", "deploy_rollback_delay", service.Deploy.RollbackConfig.Delay)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-monitor", "deploy_rollback_monitor", service.Deploy.RollbackConfig.Monitor)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-max-failure-ratio", "deploy_rollback_failure_rate", service.Deploy.RollbackConfig.MaxFailureRatio)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-order", "deploy_rollback_order", service.Deploy.RollbackConfig.Order)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-failure-action", "deploy_rollback_on_failure", service.Deploy.RollbackConfig.OnFailure)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-health-check", "deploy_rollback_health_check", service.Deploy.RollbackConfig.HealthCheck)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-min-healthy-time", "deploy_rollback_min_healthy_time", service.Deploy.RollbackConfig.MinHealthyTime)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-healthy-deadline", "deploy_rollback_healthy_deadline", service.Deploy.RollbackConfig.HealthyDeadline)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-progress-deadline", "deploy_rollback_progress_deadline", service.Deploy.RollbackConfig.ProgressDeadline)
			if service.Deploy.RollbackConfig.AutoRevertSet || service.Deploy.RollbackConfig.AutoRevert {
				templated = helmTemplateAnnotationBoolField(templated, serviceName, "bolabaden.dev/deploy-rollback-auto-revert", "deploy_rollback_auto_revert", service.Deploy.RollbackConfig.AutoRevert)
			}
			if service.Deploy.RollbackConfig.AutoPromoteSet || service.Deploy.RollbackConfig.AutoPromote {
				templated = helmTemplateAnnotationBoolField(templated, serviceName, "bolabaden.dev/deploy-rollback-auto-promote", "deploy_rollback_auto_promote", service.Deploy.RollbackConfig.AutoPromote)
			}
			if service.Deploy.RollbackConfig.CanarySet || service.Deploy.RollbackConfig.Canary > 0 {
				templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-canary", "deploy_rollback_canary", fmt.Sprintf("%d", service.Deploy.RollbackConfig.Canary))
			}
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-rollback-stagger", "deploy_rollback_stagger", service.Deploy.RollbackConfig.Stagger)
		}
		if service.Deploy.RestartPolicy != nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-restart-condition", "deploy_restart_condition", service.Deploy.RestartPolicy.Condition)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-restart-delay", "deploy_restart_delay", service.Deploy.RestartPolicy.Delay)
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-restart-max-attempts", "deploy_restart_attempts", fmt.Sprintf("%d", service.Deploy.RestartPolicy.MaxAttempts))
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/deploy-restart-window", "deploy_restart_window", service.Deploy.RestartPolicy.Window)
		}
	}
	if service.Failover != nil {
		if raw, err := json.Marshal(service.Failover); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-failover", "failover_spec", string(raw))
		}
	}
	if service.Runtime != "" {
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-runtime", "runtime", service.Runtime)
	}
	if service.PIDMode != "" {
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-pid-mode", "pid_mode", service.PIDMode)
	}
	if service.IPCMode != "" {
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-ipc-mode", "ipc_mode", service.IPCMode)
	}
	if service.pidsLimitSet || service.PidsLimit > 0 {
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-pids-limit", "pids_limit", fmt.Sprintf("%d", service.PidsLimit))
	}
	if service.shmSizeSet || service.ShmSize > 0 {
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-shm-size", "shm_size", fmt.Sprintf("%d", service.ShmSize))
	}
	if service.UserNSMode != "" {
		templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-userns-mode", "userns_mode", service.UserNSMode)
	}
	if len(service.GroupAdd) > 0 {
		templated = helmTemplateAnnotationListField(templated, serviceName, "bolabaden.dev/portable-group-add", "group_add")
	}
	if service.SeccompProfile != nil && (service.SeccompProfile.Type != "" || service.SeccompProfile.LocalhostProfile != "") {
		templated = helmTemplateSeccompProfileBlock(templated, serviceName, service.SeccompProfile)
	}
	if service.HealthCheck != nil && !isEmptyHealthCheck(service.HealthCheck) {
		if raw, err := json.Marshal(service.HealthCheck); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-healthcheck", "portable_healthcheck", string(raw))
		}
	}
	if service.StartupProbe != nil && !isEmptyHealthCheck(service.StartupProbe) {
		if raw, err := json.Marshal(service.StartupProbe); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-startup-probe", "portable_startup_probe", string(raw))
		}
	}
	if service.Ulimits != nil && !isEmptyUlimits(service.Ulimits) {
		if raw, err := json.Marshal(service.Ulimits); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-ulimits", "portable_ulimits", string(raw))
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
		if len(service.LogExtensions) > 0 {
			logging["extensions"] = copyStringInterfaceMap(service.LogExtensions)
		}
		if raw, err := json.Marshal(logging); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-logging", "portable_logging", string(raw))
		}
	}
	if !isEmptyComposeCompat(service.ComposeCompat) {
		if raw, err := json.Marshal(service.ComposeCompat); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-compose-compat", "portable_compose_compat", string(raw))
		}
	}
	if connect := nomadConnectSpecToMap(service.Connect); len(connect) > 0 {
		if raw, err := json.Marshal(connect); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-nomad-connect", "portable_nomad_connect", string(raw))
		}
	}
	if spreads := nomadSpreadSpecsToMap(service.Spreads); len(spreads) > 0 {
		if raw, err := json.Marshal(spreads); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-nomad-spread", "portable_nomad_spread", string(raw))
		}
	}
	if restart := nomadRestartBlockForService(service); len(restart) > 0 {
		if raw, err := json.Marshal(restart); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-nomad-restart", "portable_nomad_restart", string(raw))
		}
	}
	if update := nomadSchedulerExtensionMap(service, nomadUpdateExtensionKey, "x-nomad-update"); len(update) > 0 {
		if raw, err := json.Marshal(update); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-nomad-update", "portable_nomad_update", string(raw))
		}
	}
	if migrate := nomadSchedulerExtensionMap(service, nomadMigrateExtensionKey, "x-nomad-migrate"); len(migrate) > 0 {
		if raw, err := json.Marshal(migrate); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-nomad-migrate", "portable_nomad_migrate", string(raw))
		}
	}
	if reschedule := nomadSchedulerExtensionMap(service, nomadRescheduleExtensionKey, "x-nomad-reschedule"); len(reschedule) > 0 {
		if raw, err := json.Marshal(reschedule); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-nomad-reschedule", "portable_nomad_reschedule", string(raw))
		}
	}
	if mesh := meshSpecToMap(app.Mesh); len(mesh) > 0 {
		if raw, err := json.Marshal(mesh); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-mesh", "portable_mesh", string(raw))
		}
	}
	if dnsConfig := map[string]interface{}{}; len(service.DNS) > 0 || len(service.DNSSearch) > 0 || len(service.DNSOptions) > 0 {
		if len(service.DNS) > 0 {
			dnsConfig["nameservers"] = append([]string{}, service.DNS...)
		}
		if len(service.DNSSearch) > 0 {
			dnsConfig["searches"] = append([]string{}, service.DNSSearch...)
		}
		if len(service.DNSOptions) > 0 {
			var options []map[string]interface{}
			for _, option := range service.DNSOptions {
				name, value := splitDNSOption(option)
				if name == "" {
					continue
				}
				item := map[string]interface{}{"name": name}
				if value != "" {
					item["value"] = value
				}
				options = append(options, item)
			}
			if len(options) > 0 {
				dnsConfig["options"] = options
			}
		}
		templated = helmTemplateYAMLMapBlock(templated, serviceName, "dnsConfig:", "dns_config", dnsConfig)
	}
	if hostAliases := kubernetesHostAliasesFromService(service); len(hostAliases) > 0 {
		templated = helmTemplateYAMLListBlock(templated, serviceName, "hostAliases:", "host_aliases", hostAliases)
	}
	if len(service.Profiles) > 0 {
		templated = helmTemplateAnnotationListField(templated, serviceName, "bolabaden.dev/portable-profiles", "profiles")
	}
	if len(service.Links) > 0 {
		templated = helmTemplateAnnotationListField(templated, serviceName, "bolabaden.dev/portable-links", "links")
	}
	if dependencies := serviceDependencies(service); len(dependencies) > 0 {
		if raw, err := json.Marshal(dependencies); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/dependencies", "dependencies", string(raw))
		}
	}
	if !isEmptyDevelopConfig(service.Develop) {
		if raw, err := json.Marshal(service.Develop); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-develop", "portable_develop", string(raw))
		}
	}
	if service.Lifecycle != nil && !isEmptyLifecycleHooks(service.Lifecycle) {
		if raw, err := json.Marshal(service.Lifecycle); err == nil {
			templated = helmTemplateAnnotationStringField(templated, serviceName, "bolabaden.dev/portable-lifecycle", "portable_lifecycle", string(raw))
		}
		if lifecycle := serializeKubernetesLifecycle(service.Lifecycle); len(lifecycle) > 0 {
			templated = helmTemplateYAMLMapBlock(templated, serviceName, "lifecycle:", "lifecycle", lifecycle)
		}
	}
	if probe := serializeKubernetesProbe(service.HealthCheck); len(probe) > 0 {
		templated = helmTemplateYAMLMapBlock(templated, serviceName, "readinessProbe:", "readiness_probe", probe)
		templated = helmTemplateYAMLMapBlock(templated, serviceName, "livenessProbe:", "liveness_probe", cloneMap(probe))
	}
	if probe := serializeKubernetesProbeExtension(service.StartupProbe); len(probe) > 0 {
		templated = helmTemplateYAMLMapBlock(templated, serviceName, "startupProbe:", "startup_probe", probe)
	}
	if len(service.Sysctls) > 0 {
		templated = helmTemplateYAMLListBlock(templated, serviceName, "sysctls:", "sysctls", kubernetesSysctlsFromMap(service.Sysctls))
	}
	if len(service.Affinity) > 0 {
		templated = helmTemplateYAMLMapBlock(templated, serviceName, "affinity:", "affinity", service.Affinity)
	}
	if len(service.ReadinessGates) > 0 {
		templated = helmTemplateYAMLListBlock(templated, serviceName, "readinessGates:", "readiness_gates", service.ReadinessGates)
	}
	if len(service.InitContainers) > 0 {
		templated = helmTemplateYAMLListBlock(templated, serviceName, "initContainers:", "init_containers", service.InitContainers)
	}
	if len(service.SchedulingGates) > 0 {
		templated = helmTemplateYAMLListBlock(templated, serviceName, "schedulingGates:", "scheduling_gates", service.SchedulingGates)
	}
	if len(service.Tolerations) > 0 {
		templated = helmTemplateYAMLListBlockWithFollowup(templated, serviceName, "tolerations:", "tolerations", serializeKubernetesTolerationsNative(service.Tolerations), func(indent string) string {
			if len(service.ResourceClaims) == 0 {
				return ""
			}
			childIndent := indent + "  "
			var b strings.Builder
			b.WriteString(fmt.Sprintf("%s{{- with index $.Values.services %q \"resource_claims\" }}\n", indent, serviceName))
			b.WriteString(fmt.Sprintf("%sresourceClaims:\n", indent))
			b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
			if len(service.EphemeralContainers) > 0 {
				b.WriteString(fmt.Sprintf("%s{{- with index $.Values.services %q \"ephemeral_containers\" }}\n", indent, serviceName))
				b.WriteString(fmt.Sprintf("%sephemeralContainers:\n", indent))
				b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
				b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
			}
			return b.String()
		})
	}
	if len(service.ResourceClaims) > 0 {
		if len(service.Tolerations) == 0 {
			templated = replaceHelmIndentedBlock(templated, "dnsPolicy:", func(indent string) string {
				childIndent := indent + "  "
				var b strings.Builder
				b.WriteString(fmt.Sprintf("%sdnsPolicy: {{ index .Values.services %q \"dns_policy\" | quote }}\n", indent, serviceName))
				b.WriteString(fmt.Sprintf("%s{{- with index $.Values.services %q \"resource_claims\" }}\n", indent, serviceName))
				b.WriteString(fmt.Sprintf("%sresourceClaims:\n", indent))
				b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
				b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
				if len(service.EphemeralContainers) > 0 {
					b.WriteString(fmt.Sprintf("%s{{- with index $.Values.services %q \"ephemeral_containers\" }}\n", indent, serviceName))
					b.WriteString(fmt.Sprintf("%sephemeralContainers:\n", indent))
					b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
					b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
				}
				return b.String()
			})
		}
	}
	if len(service.ResourceClaims) == 0 && len(service.EphemeralContainers) > 0 {
		templated = insertHelmIndentedBlockAfterAny(templated, []string{"dnsPolicy:", "nodeName:", "schedulerName:"}, func(indent string) string {
			childIndent := indent + "  "
			var b strings.Builder
			b.WriteString(fmt.Sprintf("%s{{- with index $.Values.services %q \"ephemeral_containers\" }}\n", indent, serviceName))
			b.WriteString(fmt.Sprintf("%sephemeralContainers:\n", indent))
			b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
			return b.String()
		})
	}
	if len(service.TopologySpreadConstraints) > 0 {
		templated = helmTemplateYAMLListBlock(templated, serviceName, "topologySpreadConstraints:", "topology_spread_constraints", service.TopologySpreadConstraints)
	}
	if len(service.ImagePullSecrets) > 0 {
		templated = helmTemplateImagePullSecrets(templated, serviceName, service.ImagePullSecrets)
	}
	if service.EnableServiceLinks != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("enableServiceLinks: %t", *service.EnableServiceLinks), fmt.Sprintf("enableServiceLinks: {{ index .Values.services %q \"enable_service_links\" }}", serviceName))
	}
	if service.AutomountServiceAccountToken != nil {
		templated = strings.ReplaceAll(templated, fmt.Sprintf("automountServiceAccountToken: %t", *service.AutomountServiceAccountToken), fmt.Sprintf("automountServiceAccountToken: {{ index .Values.services %q \"automount_service_account_token\" }}", serviceName))
	}
	return templated
}

func helmTemplateImagePullSecrets(doc, serviceName string, secrets []string) string {
	if strings.TrimSpace(doc) == "" || len(secrets) == 0 {
		return doc
	}
	re := regexp.MustCompile(`(?m)^([ \t]*)imagePullSecrets:\n((?:[ \t]*- name: .*\n)+)`)
	matches := re.FindAllStringSubmatchIndex(doc, -1)
	if len(matches) == 0 {
		return doc
	}
	var b strings.Builder
	last := 0
	for _, loc := range matches {
		if len(loc) < 4 || loc[0] < 0 || loc[1] < 0 || loc[2] < 0 || loc[3] < 0 {
			continue
		}
		b.WriteString(doc[last:loc[0]])
		indent := doc[loc[2]:loc[3]]
		itemIndent := indent + "  "
		b.WriteString(fmt.Sprintf("%simagePullSecrets:\n%s{{- range $secret := index .Values.services %q \"image_pull_secrets\" }}\n%s- name: {{ $secret | quote }}\n%s{{- end }}\n", indent, indent, serviceName, itemIndent, indent))
		last = loc[1]
	}
	b.WriteString(doc[last:])
	return b.String()
}

func helmTemplateStringSliceField(doc, serviceName, fieldPrefix, fieldName, valuesKey string, values []string) string {
	if strings.TrimSpace(doc) == "" || len(values) == 0 {
		return doc
	}
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^([ \t]*)(%s)%s:\n((?:[ \t]*- .*\n)+)`, regexp.QuoteMeta(fieldPrefix), regexp.QuoteMeta(fieldName)))
	matches := re.FindAllStringSubmatchIndex(doc, -1)
	if len(matches) == 0 {
		return doc
	}
	var b strings.Builder
	last := 0
	for _, loc := range matches {
		if len(loc) < 4 || loc[0] < 0 || loc[1] < 0 || loc[2] < 0 || loc[3] < 0 {
			continue
		}
		b.WriteString(doc[last:loc[0]])
		indent := doc[loc[2]:loc[3]]
		itemIndent := indent + "  "
		b.WriteString(fmt.Sprintf("%s%s%s:\n%s{{- range $item := index .Values.services %q %q }}\n%s- {{ $item | quote }}\n%s{{- end }}\n", indent, fieldPrefix, fieldName, indent, serviceName, valuesKey, itemIndent, indent))
		last = loc[1]
	}
	b.WriteString(doc[last:])
	return b.String()
}

func helmTemplateIntSliceField(doc, serviceName, fieldName, valuesKey string, values []int64) string {
	if strings.TrimSpace(doc) == "" || len(values) == 0 {
		return doc
	}
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^([ \t]*)%s:\n((?:[ \t]*- .*\n)+)`, regexp.QuoteMeta(fieldName)))
	matches := re.FindAllStringSubmatchIndex(doc, -1)
	if len(matches) == 0 {
		return doc
	}
	var b strings.Builder
	last := 0
	for _, loc := range matches {
		if len(loc) < 4 || loc[0] < 0 || loc[1] < 0 || loc[2] < 0 || loc[3] < 0 {
			continue
		}
		b.WriteString(doc[last:loc[0]])
		indent := doc[loc[2]:loc[3]]
		itemIndent := indent + "  "
		b.WriteString(fmt.Sprintf("%s%s:\n%s{{- range $item := index .Values.services %q %q }}\n%s- {{ $item }}\n%s{{- end }}\n", indent, fieldName, indent, serviceName, valuesKey, itemIndent, indent))
		last = loc[1]
	}
	b.WriteString(doc[last:])
	return b.String()
}

func helmTemplateCapabilitiesBlock(doc, serviceName string, add, drop []string) string {
	if strings.TrimSpace(doc) == "" || (len(add) == 0 && len(drop) == 0) {
		return doc
	}
	return replaceHelmIndentedBlock(doc, "capabilities:", func(indent string) string {
		childIndent := indent + "  "
		listIndent := childIndent + "  "
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%scapabilities:\n", indent))
		b.WriteString(fmt.Sprintf("%s{{- if or (index .Values.services %q \"cap_add\") (index .Values.services %q \"cap_drop\") }}\n", childIndent, serviceName, serviceName))
		if len(add) > 0 {
			b.WriteString(fmt.Sprintf("%s{{- if index .Values.services %q \"cap_add\" }}\n", childIndent, serviceName))
			b.WriteString(fmt.Sprintf("%sadd:\n", childIndent))
			b.WriteString(fmt.Sprintf("%s{{- range $item := index .Values.services %q \"cap_add\" }}\n", childIndent, serviceName))
			b.WriteString(fmt.Sprintf("%s- {{ $item | quote }}\n", listIndent))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", childIndent))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", childIndent))
		}
		if len(drop) > 0 {
			b.WriteString(fmt.Sprintf("%s{{- if index .Values.services %q \"cap_drop\" }}\n", childIndent, serviceName))
			b.WriteString(fmt.Sprintf("%sdrop:\n", childIndent))
			b.WriteString(fmt.Sprintf("%s{{- range $item := index .Values.services %q \"cap_drop\" }}\n", childIndent, serviceName))
			b.WriteString(fmt.Sprintf("%s- {{ $item | quote }}\n", listIndent))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", childIndent))
			b.WriteString(fmt.Sprintf("%s{{- end }}\n", childIndent))
		}
		b.WriteString(fmt.Sprintf("%s{{- end }}\n", childIndent))
		return b.String()
	})
}

func helmTemplateSeccompProfileBlock(doc, serviceName string, profile *SeccompProfile) string {
	if strings.TrimSpace(doc) == "" || profile == nil {
		return doc
	}
	return replaceHelmIndentedBlock(doc, "seccompProfile:", func(indent string) string {
		childIndent := indent + "  "
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%sseccompProfile:\n", indent))
		if profile.Type != "" {
			b.WriteString(fmt.Sprintf("%stype: {{ index .Values.services %q \"seccomp_profile_type\" | quote }}\n", childIndent, serviceName))
		}
		if profile.LocalhostProfile != "" {
			b.WriteString(fmt.Sprintf("%slocalhostProfile: {{ index .Values.services %q \"seccomp_profile_localhost_profile\" | quote }}\n", childIndent, serviceName))
		}
		return b.String()
	})
}

func helmTemplateYAMLMapBlock(doc, serviceName, fieldName, valuesKey string, value map[string]interface{}) string {
	if strings.TrimSpace(doc) == "" || len(value) == 0 {
		return doc
	}
	return replaceHelmIndentedBlock(doc, fieldName, func(indent string) string {
		childIndent := indent + "  "
		return fmt.Sprintf("%s{{- with index .Values.services %q %q }}\n%s%s\n%s{{ toYaml . | nindent %d }}\n%s{{- end }}\n", indent, serviceName, valuesKey, indent, fieldName, indent, len(childIndent), indent)
	})
}

func stringMapToInterfaceMap(input map[string]string) map[string]interface{} {
	if len(input) == 0 {
		return map[string]interface{}{}
	}
	output := make(map[string]interface{}, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func helmTemplateYAMLListBlock(doc, serviceName, fieldName, valuesKey string, value []map[string]interface{}) string {
	if strings.TrimSpace(doc) == "" || len(value) == 0 {
		return doc
	}
	return replaceHelmIndentedBlock(doc, fieldName, func(indent string) string {
		childIndent := indent + "  "
		return fmt.Sprintf("%s{{- with index .Values.services %q %q }}\n%s%s\n%s{{ toYaml . | nindent %d }}\n%s{{- end }}\n", indent, serviceName, valuesKey, indent, fieldName, indent, len(childIndent), indent)
	})
}

func helmTemplateYAMLListBlockWithFollowup(doc, serviceName, fieldName, valuesKey string, value []map[string]interface{}, followup func(indent string) string) string {
	if strings.TrimSpace(doc) == "" || len(value) == 0 {
		return doc
	}
	return replaceHelmIndentedBlock(doc, fieldName, func(indent string) string {
		childIndent := indent + "  "
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%s{{- with index .Values.services %q %q }}\n", indent, serviceName, valuesKey))
		b.WriteString(fmt.Sprintf("%s%s\n", indent, fieldName))
		b.WriteString(fmt.Sprintf("%s{{ toYaml . | nindent %d }}\n", indent, len(childIndent)))
		b.WriteString(fmt.Sprintf("%s{{- end }}\n", indent))
		if followup != nil {
			b.WriteString(followup(indent))
		}
		return b.String()
	})
}

func helmTemplateAnnotationListField(doc, serviceName, annotationKey, valuesKey string) string {
	if strings.TrimSpace(doc) == "" {
		return doc
	}
	return replaceHelmAnnotationBlock(doc, annotationKey, func(indent string) string {
		return fmt.Sprintf("%s%s: {{ index .Values.services %q %q | join \"\\n\" | quote }}\n", indent, annotationKey, serviceName, valuesKey)
	})
}

func helmTemplateAnnotationBoolField(doc, serviceName, annotationKey, valuesKey string, value bool) string {
	if strings.TrimSpace(doc) == "" {
		return doc
	}
	return replaceHelmAnnotationBlock(doc, annotationKey, func(indent string) string {
		return fmt.Sprintf("%s%s: {{ index .Values.services %q %q | quote }}\n", indent, annotationKey, serviceName, valuesKey)
	})
}

func helmTemplateAnnotationStringField(doc, serviceName, annotationKey, valuesKey, value string) string {
	if strings.TrimSpace(doc) == "" || strings.TrimSpace(value) == "" {
		return doc
	}
	return replaceHelmAnnotationBlock(doc, annotationKey, func(indent string) string {
		return fmt.Sprintf("%s%s: {{ index .Values.services %q %q | quote }}\n", indent, annotationKey, serviceName, valuesKey)
	})
}

func replaceHelmAnnotationBlock(doc, annotationKey string, replacement func(indent string) string) string {
	if strings.TrimSpace(doc) == "" {
		return doc
	}
	lines := strings.SplitAfter(doc, "\n")
	var b strings.Builder
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimRight(line, "\n")
		content := strings.TrimSpace(trimmed)
		leading := len(trimmed) - len(strings.TrimLeft(trimmed, " \t"))
		if strings.HasPrefix(content, annotationKey+":") {
			indent := trimmed[:leading]
			b.WriteString(replacement(indent))
			for i+1 < len(lines) {
				next := strings.TrimRight(lines[i+1], "\n")
				if strings.TrimSpace(next) == "" {
					b.WriteString(lines[i+1])
					i++
					continue
				}
				nextLeading := len(next) - len(strings.TrimLeft(next, " \t"))
				if nextLeading <= leading {
					break
				}
				i++
			}
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

func replaceHelmIndentedBlock(doc, fieldName string, replacement func(indent string) string) string {
	if strings.TrimSpace(doc) == "" {
		return doc
	}
	lines := strings.SplitAfter(doc, "\n")
	var b strings.Builder
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimRight(line, "\n")
		content := strings.TrimSpace(trimmed)
		leading := len(trimmed) - len(strings.TrimLeft(trimmed, " \t"))
		if content == fieldName && strings.HasSuffix(trimmed, ":") {
			indent := trimmed[:leading]
			b.WriteString(replacement(indent))
			for i+1 < len(lines) {
				next := strings.TrimRight(lines[i+1], "\n")
				if strings.TrimSpace(next) == "" {
					b.WriteString(lines[i+1])
					i++
					continue
				}
				nextLeading := len(next) - len(strings.TrimLeft(next, " \t"))
				if nextLeading <= leading {
					break
				}
				i++
			}
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

func insertHelmIndentedBlockAfterAny(doc string, fieldNames []string, block func(indent string) string) string {
	if strings.TrimSpace(doc) == "" || len(fieldNames) == 0 {
		return doc
	}
	lines := strings.SplitAfter(doc, "\n")
	var b strings.Builder
	inserted := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimRight(line, "\n")
		content := strings.TrimSpace(trimmed)
		leading := len(trimmed) - len(strings.TrimLeft(trimmed, " \t"))
		if inserted {
			b.WriteString(line)
			continue
		}
		matched := false
		for _, fieldName := range fieldNames {
			if strings.HasPrefix(content, fieldName) {
				indent := trimmed[:leading]
				b.WriteString(line)
				b.WriteString(block(indent))
				matched = true
				inserted = true
				break
			}
		}
		if matched {
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

func helmApplicationExtensionsFromChart(chartPath string) (map[string]interface{}, error) {
	path := filepath.Join(chartPath, helmAppExtensionsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app extensions file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	extensions := map[string]interface{}{}
	if err := json.Unmarshal(data, &extensions); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app extensions file: %w", err)
	}
	normalized := map[string]interface{}{}
	for key, value := range extensions {
		normalized[composeApplicationExtensionKey(key)] = value
	}
	return normalized, nil
}

func helmApplicationServiceExtensionsFromChart(chartPath string) (map[string]map[string]interface{}, error) {
	path := filepath.Join(chartPath, helmServiceExtensionsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm service extensions file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	raw := map[string]map[string]interface{}{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse Helm service extensions file: %w", err)
	}
	normalized := map[string]map[string]interface{}{}
	for serviceName, extensions := range raw {
		if len(extensions) == 0 {
			continue
		}
		mapped := map[string]interface{}{}
		for key, value := range extensions {
			mapped[composeApplicationExtensionKey(key)] = value
		}
		if len(mapped) > 0 {
			normalized[serviceName] = mapped
		}
	}
	return normalized, nil
}

func helmApplicationNameFromChart(chartPath string) (string, error) {
	path := filepath.Join(chartPath, helmAppNameFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read Helm app name file: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func helmApplicationModelsFromChart(chartPath string) (map[string]*ComposeModel, error) {
	path := filepath.Join(chartPath, helmAppModelsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app models file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var raw map[string]map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app models file: %w", err)
	}
	return composeModelsFromExtensionMap(raw), nil
}

func helmApplicationIncludesFromChart(chartPath string) ([]interface{}, error) {
	path := filepath.Join(chartPath, helmAppIncludesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app includes file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var includes []interface{}
	if err := json.Unmarshal(data, &includes); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app includes file: %w", err)
	}
	return cloneInterfaceSlice(includes), nil
}

func helmApplicationNetworksFromChart(chartPath string) (map[string]*Network, error) {
	path := filepath.Join(chartPath, helmAppNetworksFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app networks file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	networks := map[string]*Network{}
	if err := json.Unmarshal(data, &networks); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app networks file: %w", err)
	}
	return copyNetworks(networks), nil
}

func helmApplicationVolumesFromChart(chartPath string) (map[string]*Volume, error) {
	path := filepath.Join(chartPath, helmAppVolumesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app volumes file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	volumes := map[string]*Volume{}
	if err := json.Unmarshal(data, &volumes); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app volumes file: %w", err)
	}
	return copyVolumes(volumes), nil
}

func helmApplicationConfigsFromChart(chartPath string) (map[string]*Config, error) {
	path := filepath.Join(chartPath, helmAppConfigsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app configs file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	configs := map[string]*Config{}
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app configs file: %w", err)
	}
	return copyConfigs(configs), nil
}

func helmApplicationSecretsFromChart(chartPath string) (map[string]*Secret, error) {
	path := filepath.Join(chartPath, helmAppSecretsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app secrets file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	secrets := map[string]*Secret{}
	if err := json.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app secrets file: %w", err)
	}
	return copySecrets(secrets), nil
}

func helmApplicationRoutesFromChart(chartPath string) (map[string]*RouteSpec, error) {
	path := filepath.Join(chartPath, helmAppRoutesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app routes file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var raw map[string]map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app routes file: %w", err)
	}
	routes := map[string]*RouteSpec{}
	for name, mapped := range raw {
		if route := routeSpecFromExtensionMap(mapped); route != nil {
			routes[name] = route
		}
	}
	return cloneRouteSpecMap(routes), nil
}

func helmApplicationPoliciesFromChart(chartPath string) (map[string]*PolicySpec, error) {
	path := filepath.Join(chartPath, helmAppPoliciesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm app policies file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var raw map[string]map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse Helm app policies file: %w", err)
	}
	policies := map[string]*PolicySpec{}
	for name, mapped := range raw {
		if policy := policySpecFromExtensionMap(mapped); policy != nil {
			policies[name] = policy
		}
	}
	return clonePolicySpecMap(policies), nil
}

func helmApplicationCanonicalRawResourcesFromChart(chartPath string) ([]interface{}, error) {
	path := filepath.Join(chartPath, helmCanonicalRawResourcesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm canonical raw resources file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var resources []interface{}
	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, fmt.Errorf("failed to parse Helm canonical raw resources file: %w", err)
	}
	return cloneInterfaceSlice(resources), nil
}

func helmApplicationKubernetesRawResourcesFromChart(chartPath string) ([]interface{}, error) {
	path := filepath.Join(chartPath, helmKubernetesRawResourcesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read Helm Kubernetes raw resources file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var resources []interface{}
	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, fmt.Errorf("failed to parse Helm Kubernetes raw resources file: %w", err)
	}
	return cloneInterfaceSlice(resources), nil
}

func writeHelmApplicationExtensions(app *Application, chartPath string) error {
	if app == nil || len(app.Extensions) == 0 {
		return nil
	}
	extensions := copyStringInterfaceMap(app.Extensions)
	delete(extensions, "chart")
	delete(extensions, helmRawFilesExtension)
	delete(extensions, helmAppExtensionsFile)
	delete(extensions, "x-platform")
	stripPortableCanonicalExtensionKeys(extensions)
	delete(extensions, composeKubernetesServicesExtensionKey)
	normalized := map[string]interface{}{}
	for key, value := range extensions {
		normalized[composeApplicationExtensionKey(key)] = value
	}
	if len(normalized) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app extensions: %w", err)
	}
	path := filepath.Join(chartPath, helmAppExtensionsFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmServiceExtensions(app *Application, chartPath string) error {
	if app == nil || len(app.Services) == 0 {
		return nil
	}
	extensions := map[string]map[string]interface{}{}
	serviceNames := make([]string, 0, len(app.Services))
	for name, service := range app.Services {
		if service == nil || len(service.Extensions) == 0 {
			continue
		}
		serviceNames = append(serviceNames, name)
	}
	sort.Strings(serviceNames)
	for _, name := range serviceNames {
		service := app.Services[name]
		if service == nil || len(service.Extensions) == 0 {
			continue
		}
		mapped := copyStringInterfaceMap(service.Extensions)
		delete(mapped, "x-platform")
		stripPortableCanonicalExtensionKeys(mapped)
		normalized := map[string]interface{}{}
		for key, value := range mapped {
			normalized[composeApplicationExtensionKey(key)] = value
		}
		if len(normalized) > 0 {
			extensions[name] = normalized
		}
	}
	if len(extensions) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(extensions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm service extensions: %w", err)
	}
	path := filepath.Join(chartPath, helmServiceExtensionsFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationName(app *Application, chartPath string) error {
	if app == nil || app.Name == "" {
		return nil
	}
	path := filepath.Join(chartPath, helmAppNameFile)
	return os.WriteFile(path, []byte(app.Name+"\n"), 0644)
}

func writeHelmApplicationModels(app *Application, chartPath string) error {
	models := applicationModelsForEmit(app)
	if len(models) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app models: %w", err)
	}
	path := filepath.Join(chartPath, helmAppModelsFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationIncludes(app *Application, chartPath string) error {
	if app == nil || len(app.IncludeEntries) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(app.IncludeEntries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app includes: %w", err)
	}
	path := filepath.Join(chartPath, helmAppIncludesFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationNetworks(app *Application, chartPath string) error {
	networks := applicationNetworksForEmit(app)
	if len(networks) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(networks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app networks: %w", err)
	}
	path := filepath.Join(chartPath, helmAppNetworksFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationVolumes(app *Application, chartPath string) error {
	volumes := applicationVolumesForEmit(app)
	if len(volumes) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(volumes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app volumes: %w", err)
	}
	path := filepath.Join(chartPath, helmAppVolumesFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationConfigs(app *Application, chartPath string) error {
	configs := applicationConfigsForEmit(app)
	if len(configs) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app configs: %w", err)
	}
	path := filepath.Join(chartPath, helmAppConfigsFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationSecrets(app *Application, chartPath string) error {
	secrets := applicationSecretsForEmit(app)
	if len(secrets) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app secrets: %w", err)
	}
	path := filepath.Join(chartPath, helmAppSecretsFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationRoutes(app *Application, chartPath string) error {
	routes := canonicalRoutesForApplication(app)
	if len(routes) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app routes: %w", err)
	}
	path := filepath.Join(chartPath, helmAppRoutesFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmApplicationPolicies(app *Application, chartPath string) error {
	policies := canonicalPoliciesForApplication(app)
	if len(policies) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(policies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm app policies: %w", err)
	}
	path := filepath.Join(chartPath, helmAppPoliciesFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmCanonicalRawResources(app *Application, chartPath string) error {
	resources := canonicalRawResourcesForBridge(app, PlatformHelm)
	if len(resources) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm canonical raw resources: %w", err)
	}
	path := filepath.Join(chartPath, helmCanonicalRawResourcesFile)
	return os.WriteFile(path, data, 0644)
}

func writeHelmKubernetesRawResources(app *Application, chartPath string) error {
	resources := kubernetesOpaqueManifestSpecsForApplication(app)
	if resources == nil {
		resources = map[string]*KubernetesOpaqueManifestSpec{}
	}
	if app != nil {
		if len(app.Configs) > 0 {
			names := make([]string, 0, len(app.Configs))
			for name, config := range app.Configs {
				if config != nil {
					names = append(names, name)
				}
			}
			sort.Strings(names)
			for _, name := range names {
				if spec := helmKubernetesOpaqueManifestSpecForConfig(name, app.Configs[name], app.Namespace); spec != nil {
					key := kubernetesDocumentKeyFromMap(spec.Raw)
					if key == "" {
						key = spec.Kind + "/" + spec.Namespace + "/" + spec.Name
					}
					if _, exists := resources[key]; !exists {
						resources[key] = spec
					}
				}
			}
		}
		if len(app.Secrets) > 0 {
			names := make([]string, 0, len(app.Secrets))
			for name, secret := range app.Secrets {
				if secret != nil {
					names = append(names, name)
				}
			}
			sort.Strings(names)
			for _, name := range names {
				if spec := helmKubernetesOpaqueManifestSpecForSecret(name, app.Secrets[name], app.Namespace); spec != nil {
					key := kubernetesDocumentKeyFromMap(spec.Raw)
					if key == "" {
						key = spec.Kind + "/" + spec.Namespace + "/" + spec.Name
					}
					if _, exists := resources[key]; !exists {
						resources[key] = spec
					}
				}
			}
		}
	}
	if len(resources) == 0 {
		resourcesList := canonicalKubernetesRawResourcesForApplication(app)
		if len(resourcesList) == 0 {
			return nil
		}
		data, err := json.MarshalIndent(resourcesList, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal Helm Kubernetes raw resources: %w", err)
		}
		path := filepath.Join(chartPath, helmKubernetesRawResourcesFile)
		return os.WriteFile(path, data, 0644)
	}
	list := make([]interface{}, 0, len(resources))
	keys := make([]string, 0, len(resources))
	for key := range resources {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if resources[key] == nil {
			continue
		}
		list = append(list, resources[key])
	}
	if len(list) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Helm Kubernetes raw resources: %w", err)
	}
	path := filepath.Join(chartPath, helmKubernetesRawResourcesFile)
	return os.WriteFile(path, data, 0644)
}

func helmKubernetesOpaqueManifestSpecForConfig(name string, config *Config, namespace string) *KubernetesOpaqueManifestSpec {
	if config == nil {
		return nil
	}
	doc, err := serializeKubernetesConfigMap(name, config, namespace)
	if err != nil || strings.TrimSpace(doc) == "" {
		return nil
	}
	var mapped map[string]interface{}
	if err := yaml.Unmarshal([]byte(doc), &mapped); err != nil {
		return nil
	}
	return kubernetesOpaqueManifestSpecFromMap(mapped)
}

func helmKubernetesOpaqueManifestSpecForSecret(name string, secret *Secret, namespace string) *KubernetesOpaqueManifestSpec {
	if secret == nil {
		return nil
	}
	doc, err := serializeKubernetesSecret(name, secret, namespace)
	if err != nil || strings.TrimSpace(doc) == "" {
		return nil
	}
	var mapped map[string]interface{}
	if err := yaml.Unmarshal([]byte(doc), &mapped); err != nil {
		return nil
	}
	return kubernetesOpaqueManifestSpecFromMap(mapped)
}

type noHelmRawFilesError struct{}

func (noHelmRawFilesError) Error() string {
	return "application does not contain raw Helm chart files"
}

func isNoHelmRawFilesError(err error) bool {
	_, ok := err.(noHelmRawFilesError)
	return ok
}

func restoreHelmRawChart(app *Application, chartPath string) error {
	files, ok := helmRawFilesExtensionValue(app)
	if !ok || len(files) == 0 {
		if raw := helmCanonicalRawFiles(app); raw != nil {
			files = raw
			ok = len(files) > 0
		}
	}
	if !ok || len(files) == 0 {
		return noHelmRawFilesError{}
	}
	if err := os.MkdirAll(chartPath, 0755); err != nil {
		return fmt.Errorf("failed to create chart directory: %w", err)
	}
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if !isSafeHelmChartFileName(name) {
			return fmt.Errorf("unsafe Helm chart file path %q", name)
		}
		data, err := base64.StdEncoding.DecodeString(files[name])
		if err != nil {
			return fmt.Errorf("failed to decode Helm chart file %s: %w", name, err)
		}
		target := filepath.Join(chartPath, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("failed to create Helm chart file directory %s: %w", name, err)
		}
		if err := os.WriteFile(target, data, 0644); err != nil {
			return fmt.Errorf("failed to write Helm chart file %s: %w", name, err)
		}
	}
	return nil
}

func helmCanonicalRawFiles(app *Application) map[string]string {
	canonical := canonicalForApplication(app)
	if canonical == nil {
		return nil
	}
	resource := canonical.Resources["helm:raw:files"]
	if resource == nil || resource.Raw == nil {
		return nil
	}
	switch typed := resource.Raw.(type) {
	case map[string]string:
		return copyStringMap(typed)
	case map[string]interface{}:
		files := map[string]string{}
		for key, value := range typed {
			if encoded := toString(value); encoded != "" {
				files[key] = encoded
			}
		}
		if len(files) > 0 {
			return files
		}
	}
	return nil
}

func helmRawFilesExtensionValue(app *Application) (map[string]string, bool) {
	if app == nil || app.Extensions == nil {
		return nil, false
	}
	value, ok := app.Extensions[helmRawFilesExtension]
	if !ok {
		return nil, false
	}
	switch files := value.(type) {
	case map[string]string:
		return files, true
	case map[string]interface{}:
		result := map[string]string{}
		for key, value := range files {
			if encoded := toString(value); encoded != "" {
				result[key] = encoded
			}
		}
		return result, true
	default:
		return nil, false
	}
}

func isSafeHelmChartFileName(name string) bool {
	if name == "" || filepath.IsAbs(name) {
		return false
	}
	clean := filepath.ToSlash(filepath.Clean(name))
	return clean == name && clean != "." && !strings.HasPrefix(clean, "../") && clean != ".."
}
