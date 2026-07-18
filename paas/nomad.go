package paas

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/tmccombs/hcl2json/convert"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/api/resource"
)

const nomadRawHCLExtension = "nomad.raw.hcl"
const nomadExtensionsMetaKey = "bolabaden_extensions"
const nomadAppExtensionsMetaKey = "bolabaden_app_extensions"
const nomadAppNameMetaKey = "bolabaden_app_name"
const nomadAppVersionMetaKey = "bolabaden_app_version"
const nomadAppModelsMetaKey = "bolabaden_app_models"
const nomadAppIncludesMetaKey = "bolabaden_app_includes"
const nomadAppNetworksMetaKey = "bolabaden_app_networks"
const nomadAppVolumesMetaKey = "bolabaden_app_volumes"
const nomadAppConfigsMetaKey = "bolabaden_app_configs"
const nomadAppSecretsMetaKey = "bolabaden_app_secrets"
const nomadAppRoutesMetaKey = "bolabaden_app_routes"
const nomadAppPoliciesMetaKey = "bolabaden_app_policies"
const nomadCanonicalRawResourcesMetaKey = "bolabaden_canonical_raw_resources"
const nomadKubernetesRawResourcesMetaKey = "bolabaden_kubernetes_raw_resources"
const nomadComposeCompatMetaKey = "bolabaden_compose_compat"
const nomadLinksMetaKey = "bolabaden_links"
const nomadNetworkAttachmentsMetaKey = "bolabaden_network_attachments"
const nomadFailoverMetaKey = "bolabaden_failover"
const nomadVolumesMetaKey = "bolabaden_volumes"
const nomadTmpfsMetaKey = "bolabaden_tmpfs"
const nomadSpreadExtensionKey = "nomad.spread"
const nomadConnectExtensionKey = "nomad.connect"
const nomadRestartExtensionKey = "nomad.restart"
const nomadUpdateExtensionKey = "nomad.update"
const nomadMigrateExtensionKey = "nomad.migrate"
const nomadRescheduleExtensionKey = "nomad.reschedule"

// ParseNomadHCL parses Nomad job HCL into the portable application model.
func ParseNomadHCL(content string) (*Application, error) {
	return parseNomadContent(content, false)
}

// ParseNomadJSON parses Nomad job JSON into the portable application model.
func ParseNomadJSON(content string) (*Application, error) {
	return parseNomadJSONContent(content)
}

func parseNomadContent(content string, jsonMode bool) (*Application, error) {
	parser := hclparse.NewParser()
	var (
		file  *hcl.File
		diags hcl.Diagnostics
	)
	if jsonMode {
		file, diags = parser.ParseJSON([]byte(content), "nomad.json")
	} else {
		file, diags = parser.ParseHCL([]byte(content), "nomad.hcl")
	}
	if diags.HasErrors() {
		kind := "HCL"
		if jsonMode {
			kind = "JSON"
		}
		return nil, fmt.Errorf("failed to parse Nomad %s: %s", kind, diags.Error())
	}

	app := &Application{
		Platform: PlatformNomad,
		Services: make(map[string]*Service),
		Networks: make(map[string]*Network),
		Volumes:  make(map[string]*Volume),
		Configs:  make(map[string]*Config),
		Secrets:  make(map[string]*Secret),
		Models:   make(map[string]*ComposeModel),
		Extensions: map[string]interface{}{
			nomadRawHCLExtension: content,
		},
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		rehydrateComposeApplicationExtensions(app)
		syncPortableApplicationState(app)
		app.AttachCanonical()
		app.Canonical.AddResource(ResourceKindRaw, PlatformNomad, "hcl", "NomadHCL", content)
		return app, nil
	}
	for _, job := range body.Blocks {
		if job.Type != "job" {
			continue
		}
		applyNomadApplicationPortableMeta(app, job)
		parseNomadJob(app, job)
	}

	rehydrateComposeApplicationExtensions(app)
	syncPortableApplicationState(app)
	app.AttachCanonical()
	app.Canonical.AddResource(ResourceKindRaw, PlatformNomad, "hcl", "NomadHCL", content)
	return app, nil
}

func parseNomadJSONContent(content string) (*Application, error) {
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(content), &root); err != nil {
		return nil, fmt.Errorf("failed to parse Nomad JSON: %w", err)
	}

	app := &Application{
		Platform: PlatformNomad,
		Services: make(map[string]*Service),
		Networks: make(map[string]*Network),
		Volumes:  make(map[string]*Volume),
		Configs:  make(map[string]*Config),
		Secrets:  make(map[string]*Secret),
		Models:   make(map[string]*ComposeModel),
		Extensions: map[string]interface{}{
			nomadRawHCLExtension: content,
		},
	}

	if jobs, ok := asMap(root["job"]); ok {
		for jobName, rawJobs := range jobs {
			for _, rawJob := range nomadJSONEntrySlice(rawJobs) {
				job, ok := asMap(rawJob)
				if !ok {
					continue
				}
				parseNomadJSONJob(app, jobName, job)
			}
		}
	}

	rehydrateComposeApplicationExtensions(app)
	syncPortableApplicationState(app)
	app.AttachCanonical()
	app.Canonical.AddResource(ResourceKindRaw, PlatformNomad, "json", "NomadJSON", content)
	return app, nil
}

func parseNomadJSONJob(app *Application, jobName string, job map[string]interface{}) {
	if app == nil || len(job) == 0 {
		return
	}
	if app.Name == "" && strings.TrimSpace(jobName) != "" {
		app.Name = jobName
	}
	jobType := toString(job["type"])
	jobRestart := nomadJSONSchedulerBlockExtensions(job["restart"])
	jobUpdate := nomadJSONSchedulerBlockExtensions(job["update"])
	jobMigrate := nomadJSONSchedulerBlockExtensions(job["migrate"])
	jobReschedule := nomadJSONSchedulerBlockExtensions(job["reschedule"])
	if groups, ok := asMap(job["group"]); ok {
		for groupName, rawGroups := range groups {
			for _, rawGroup := range nomadJSONEntrySlice(rawGroups) {
				group, ok := asMap(rawGroup)
				if !ok {
					continue
				}
				parseNomadJSONGroup(app, jobType, groupName, group, jobRestart, jobUpdate, jobMigrate, jobReschedule)
			}
		}
	}
}

func parseNomadJSONGroup(app *Application, jobType, groupName string, group map[string]interface{}, jobRestart, jobUpdate, jobMigrate, jobReschedule map[string]interface{}) {
	if app == nil || len(group) == 0 {
		return
	}
	count := toInt(group["count"])
	groupRestart := mergeNomadSchedulerBlocks(jobRestart, nomadJSONSchedulerBlockExtensions(group["restart"]))
	groupUpdate := mergeNomadSchedulerBlocks(jobUpdate, nomadJSONSchedulerBlockExtensions(group["update"]))
	groupMigrate := mergeNomadSchedulerBlocks(jobMigrate, nomadJSONSchedulerBlockExtensions(group["migrate"]))
	groupReschedule := mergeNomadSchedulerBlocks(jobReschedule, nomadJSONSchedulerBlockExtensions(group["reschedule"]))
	if tasks, ok := asMap(group["task"]); ok {
		for taskName, rawTasks := range tasks {
			for _, rawTask := range nomadJSONEntrySlice(rawTasks) {
				task, ok := asMap(rawTask)
				if !ok {
					continue
				}
				service := parseNomadJSONTask(taskName, count, task, groupRestart)
				if service == nil {
					continue
				}
				if strings.EqualFold(jobType, "system") {
					deploy := ensureServiceDeploy(service)
					deploy.Mode = "global"
				} else if strings.EqualFold(jobType, "batch") {
					deploy := ensureServiceDeploy(service)
					if deploy.Mode == "" {
						deploy.Mode = "replicated-job"
					}
				}
				if service.Name == "" {
					service.Name = taskName
				}
				if service.Deploy == nil && count > 0 {
					service.Deploy = &DeploySpec{Replicas: count}
				}
				if service.Replicas == 0 && count > 0 {
					service.Replicas = count
				}
				applyNomadSchedulerBlockExtensions(service, groupUpdate, groupMigrate, groupReschedule)
				app.Services[service.Name] = service
			}
		}
	}
}

func parseNomadJSONTask(taskName string, count int, task map[string]interface{}, inheritedRestart map[string]interface{}) *Service {
	if len(task) == 0 {
		return nil
	}
	service := &Service{
		Name:        taskName,
		Platform:    PlatformNomad,
		Environment: map[string]string{},
		Extensions:  map[string]interface{}{},
	}
	if driver := toString(task["driver"]); driver != "" {
		service.Extensions["nomad.driver"] = driver
	}
	if config, ok := nomadJSONFirstObject(task["config"]); ok {
		if image := toString(config["image"]); image != "" {
			service.Image = image
		}
		if ports := interfaceSliceToStringSliceLoose(config["ports"]); len(ports) > 0 {
			service.Extensions["nomad.ports"] = ports
		}
	}
	if env, ok := nomadJSONFirstObject(task["env"]); ok {
		for key, value := range env {
			service.Environment[key] = toString(value)
		}
	}
	if count > 0 {
		service.Replicas = count
		service.Deploy = &DeploySpec{Replicas: count}
	}
	if restart := mergeNomadSchedulerBlocks(inheritedRestart, nomadJSONSchedulerBlockExtensions(task["restart"])); len(restart) > 0 {
		applyNomadRestartExtensions(service, restart)
	}
	return service
}

func nomadJSONEntrySlice(value interface{}) []interface{} {
	switch typed := value.(type) {
	case nil:
		return nil
	case []interface{}:
		return typed
	default:
		return []interface{}{typed}
	}
}

func nomadJSONFirstObject(value interface{}) (map[string]interface{}, bool) {
	for _, item := range nomadJSONEntrySlice(value) {
		if mapped, ok := asMap(item); ok {
			return mapped, true
		}
	}
	return nil, false
}

func parseNomadJob(app *Application, job *hclsyntax.Block) {
	jobType := stringAttribute(job.Body, "type")
	jobRestart := parseNomadRestartBlockExtensions(job.Body)
	jobUpdate := parseNomadSchedulerBlockExtensions(job.Body, "update")
	jobMigrate := parseNomadSchedulerBlockExtensions(job.Body, "migrate")
	jobReschedule := parseNomadSchedulerBlockExtensions(job.Body, "reschedule")
	for _, group := range job.Body.Blocks {
		if group.Type != "group" {
			continue
		}
		groupCount := intAttribute(group.Body, "count")
		groupRestart := mergeNomadSchedulerBlocks(jobRestart, parseNomadRestartBlockExtensions(group.Body))
		ports := parseNomadGroupPorts(group)
		dependencies := parseNomadGroupDependencies(group)
		groupServices := nomadGroupServiceBlocks(group)
		placement := parseNomadGroupPlacement(group)
		spreads := parseNomadGroupSpreadExtensions(group)
		groupUpdate := mergeNomadSchedulerBlocks(jobUpdate, parseNomadSchedulerBlockExtensions(group.Body, "update"))
		groupMigrate := mergeNomadSchedulerBlocks(jobMigrate, parseNomadSchedulerBlockExtensions(group.Body, "migrate"))
		groupReschedule := mergeNomadSchedulerBlocks(jobReschedule, parseNomadSchedulerBlockExtensions(group.Body, "reschedule"))
		portable := parseNomadGroupPortableMeta(group)
		for _, task := range group.Body.Blocks {
			if task.Type != "task" || len(task.Labels) == 0 {
				continue
			}
			if isNomadDependencyTask(task) {
				continue
			}
			service := parseNomadTask(task, groupCount, ports, groupRestart)
			if strings.EqualFold(jobType, "system") {
				deploy := ensureServiceDeploy(service)
				deploy.Mode = "global"
			} else if strings.EqualFold(jobType, "batch") {
				deploy := ensureServiceDeploy(service)
				if deploy.Mode == "" {
					deploy.Mode = "replicated-job"
				}
			}
			applyNomadPortableMeta(service, portable)
			for _, groupService := range groupServices {
				parseNomadServiceBlock(service, groupService.Body)
			}
			for _, dependency := range dependencies {
				service.Dependencies = appendUniqueDependency(service.Dependencies, dependency)
				service.DependsOn = appendUniqueName(service.DependsOn, dependency.Name)
			}
			applyNomadPlacement(service, placement)
			applyNomadSchedulerBlockExtensions(service, groupUpdate, groupMigrate, groupReschedule)
			applyNomadSpreadExtensions(service, spreads)
			app.Services[service.Name] = service
		}
	}
}

func nomadGroupServiceBlocks(group *hclsyntax.Block) []*hclsyntax.Block {
	if group == nil || group.Body == nil {
		return nil
	}
	var services []*hclsyntax.Block
	for _, block := range group.Body.Blocks {
		if block.Type == "service" {
			services = append(services, block)
		}
	}
	return services
}

func applyNomadApplicationPortableMeta(app *Application, job *hclsyntax.Block) {
	if app == nil || job == nil || job.Body == nil {
		return
	}
	meta := parseNomadJobMeta(job)
	if len(meta) == 0 {
		return
	}
	if raw := nomadPortableMetaValue(meta, nomadAppExtensionsMetaKey, "bolabaden.app_extensions"); raw != "" {
		var extensions map[string]interface{}
		if json.Unmarshal([]byte(raw), &extensions) == nil && len(extensions) > 0 {
			if app.Extensions == nil {
				app.Extensions = map[string]interface{}{}
			}
			for key, value := range extensions {
				if key == nomadRawHCLExtension || key == nomadAppExtensionsMetaKey {
					continue
				}
				app.Extensions[key] = value
			}
		}
	}
	if name := nomadPortableMetaValue(meta, nomadAppNameMetaKey, "bolabaden.app_name"); name != "" && app.Name == "" {
		app.Name = name
	}
	if version := nomadPortableMetaValue(meta, nomadAppVersionMetaKey, "bolabaden.app_version"); version != "" && app.Version == "" {
		app.Version = version
	}
	if raw := nomadPortableMetaValue(meta, nomadAppModelsMetaKey, "bolabaden.app_models"); raw != "" {
		var decoded map[string]map[string]interface{}
		if json.Unmarshal([]byte(raw), &decoded) == nil && len(decoded) > 0 {
			models := composeModelsFromExtensionMap(decoded)
			if len(models) > 0 {
				if app.Models == nil {
					app.Models = map[string]*ComposeModel{}
				}
				for name, model := range models {
					if _, exists := app.Models[name]; !exists {
						app.Models[name] = cloneComposeModel(model)
					}
				}
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadAppIncludesMetaKey, "bolabaden.app_includes"); raw != "" {
		var includes []interface{}
		if json.Unmarshal([]byte(raw), &includes) == nil && len(includes) > 0 {
			app.IncludeEntries = mergeIncludeEntries(app.IncludeEntries, includes)
			app.Includes = mergeUniqueStrings(app.Includes, composeIncludePaths(includes))
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadAppNetworksMetaKey, "bolabaden.app_networks"); raw != "" {
		var networks map[string]*Network
		if json.Unmarshal([]byte(raw), &networks) == nil && len(networks) > 0 {
			if app.Networks == nil {
				app.Networks = map[string]*Network{}
			}
			for name, network := range networks {
				if network == nil {
					continue
				}
				if _, exists := app.Networks[name]; !exists {
					app.Networks[name] = cloneNetwork(network)
				}
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadAppVolumesMetaKey, "bolabaden.app_volumes"); raw != "" {
		var volumes map[string]*Volume
		if json.Unmarshal([]byte(raw), &volumes) == nil && len(volumes) > 0 {
			if app.Volumes == nil {
				app.Volumes = map[string]*Volume{}
			}
			for name, volume := range volumes {
				if volume == nil {
					continue
				}
				if _, exists := app.Volumes[name]; !exists {
					app.Volumes[name] = cloneVolume(volume)
				}
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadAppConfigsMetaKey, "bolabaden.app_configs"); raw != "" {
		var configs map[string]*Config
		if json.Unmarshal([]byte(raw), &configs) == nil && len(configs) > 0 {
			if app.Configs == nil {
				app.Configs = map[string]*Config{}
			}
			for name, config := range configs {
				if config == nil {
					continue
				}
				if _, exists := app.Configs[name]; !exists {
					app.Configs[name] = cloneConfig(config)
				}
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadAppSecretsMetaKey, "bolabaden.app_secrets"); raw != "" {
		var secrets map[string]*Secret
		if json.Unmarshal([]byte(raw), &secrets) == nil && len(secrets) > 0 {
			if app.Secrets == nil {
				app.Secrets = map[string]*Secret{}
			}
			for name, secret := range secrets {
				if secret == nil {
					continue
				}
				if _, exists := app.Secrets[name]; !exists {
					app.Secrets[name] = cloneSecret(secret)
				}
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadAppRoutesMetaKey, "bolabaden.app_routes"); raw != "" {
		var decoded map[string]map[string]interface{}
		if json.Unmarshal([]byte(raw), &decoded) == nil && len(decoded) > 0 {
			routes := map[string]*RouteSpec{}
			for name, mapped := range decoded {
				if route := routeSpecFromExtensionMap(mapped); route != nil {
					routes[name] = route
				}
			}
			if len(routes) > 0 {
				if app.Extensions == nil {
					app.Extensions = map[string]interface{}{}
				}
				app.Extensions[nomadAppRoutesMetaKey] = routes
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadAppPoliciesMetaKey, "bolabaden.app_policies"); raw != "" {
		var decoded map[string]map[string]interface{}
		if json.Unmarshal([]byte(raw), &decoded) == nil && len(decoded) > 0 {
			policies := map[string]*PolicySpec{}
			for name, mapped := range decoded {
				if policy := policySpecFromExtensionMap(mapped); policy != nil {
					policies[name] = policy
				}
			}
			if len(policies) > 0 {
				if app.Extensions == nil {
					app.Extensions = map[string]interface{}{}
				}
				app.Extensions[nomadAppPoliciesMetaKey] = policies
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadCanonicalRawResourcesMetaKey, "bolabaden.canonical_raw_resources"); raw != "" {
		var resources []interface{}
		if json.Unmarshal([]byte(raw), &resources) == nil && len(resources) > 0 {
			if app.Extensions == nil {
				app.Extensions = map[string]interface{}{}
			}
			app.Extensions[composeCanonicalRawResourcesExtension] = resources
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadKubernetesRawResourcesMetaKey, "bolabaden.kubernetes_raw_resources"); raw != "" {
		var resources []interface{}
		if json.Unmarshal([]byte(raw), &resources) == nil && len(resources) > 0 {
			if app.Extensions == nil {
				app.Extensions = map[string]interface{}{}
			}
			app.Extensions[composeKubernetesRawResourcesExtension] = resources
		}
	}
}

func parseNomadJobMeta(job *hclsyntax.Block) map[string]string {
	meta := map[string]string{}
	if job == nil || job.Body == nil {
		return meta
	}
	for _, block := range job.Body.Blocks {
		if block.Type != "meta" {
			continue
		}
		for key, attr := range block.Body.Attributes {
			if val := expressionString(attr.Expr); val != "" {
				meta[key] = val
			}
		}
	}
	return meta
}

func parseNomadTask(task *hclsyntax.Block, groupCount int, ports map[string]PortMapping, inheritedRestart map[string]interface{}) *Service {
	name := task.Labels[0]
	service := &Service{
		Name:        name,
		Platform:    PlatformNomad,
		Environment: map[string]string{},
		Labels:      map[string]string{},
		Extensions:  map[string]interface{}{},
	}
	if groupCount > 0 {
		service.Replicas = groupCount
		service.Deploy = &DeploySpec{Replicas: groupCount}
	}
	if driver := stringAttribute(task.Body, "driver"); driver != "" {
		service.Extensions["nomad.driver"] = driver
	}
	if user := stringAttribute(task.Body, "user"); user != "" {
		service.User = user
	}
	if value := stringAttribute(task.Body, "kill_signal"); value != "" {
		service.StopSignal = value
	}
	if value := stringAttribute(task.Body, "kill_timeout"); value != "" {
		service.StopGracePeriod = value
	}
	restartApplied := false

	for _, block := range task.Body.Blocks {
		switch block.Type {
		case "config":
			parseNomadConfig(service, block.Body, ports)
		case "env":
			for key, attr := range block.Body.Attributes {
				if val := expressionString(attr.Expr); val != "" {
					service.Environment[key] = val
				}
			}
		case "template":
			parseNomadTemplateBlock(service, block.Body)
		case "resources":
			service.Deploy = ensureServiceDeploy(service)
			service.Deploy.Resources = parseNomadResources(block.Body)
			if service.Deploy.Resources != nil {
				if service.Deploy.Resources.CPUReservation != "" {
					service.CPUs = service.Deploy.Resources.CPUReservation
				}
				if service.Deploy.Resources.MemoryReservation != "" {
					service.MemLimit = service.Deploy.Resources.MemoryReservation
					service.MemoryLimit = service.Deploy.Resources.MemoryReservation
					service.MemReservation = service.Deploy.Resources.MemoryReservation
				}
			}
		case "service":
			parseNomadServiceBlock(service, block.Body)
		case "meta":
			for key, attr := range block.Body.Attributes {
				if val := expressionString(attr.Expr); val != "" {
					if service.Extensions == nil {
						service.Extensions = map[string]interface{}{}
					}
					service.Extensions[key] = val
				}
			}
		case "restart":
			if restart := mergeNomadSchedulerBlocks(inheritedRestart, nomadAttributesToMap(block.Body)); len(restart) > 0 {
				applyNomadRestartExtensions(service, restart)
				restartApplied = true
			}
		}
	}
	if !restartApplied && len(inheritedRestart) > 0 {
		applyNomadRestartExtensions(service, inheritedRestart)
	}
	if rawJSON := toString(service.Extensions[nomadTmpfsMetaKey]); rawJSON != "" {
		var raw map[string]string
		if json.Unmarshal([]byte(rawJSON), &raw) == nil && len(raw) > 0 {
			for _, volume := range service.Volumes {
				if volume.Type != "tmpfs" || volume.Target == "" {
					continue
				}
				if size := raw[volume.Target]; size != "" {
					volume.Options = ensureStringMap(volume.Options)
					volume.Options["size"] = size
					if volume.TmpfsExtensions == nil {
						volume.TmpfsExtensions = map[string]interface{}{}
					}
					if _, exists := volume.TmpfsExtensions["size"]; !exists {
						volume.TmpfsExtensions["size"] = size
					}
				}
			}
		}
	}
	return service
}

func parseNomadGroupDependencies(group *hclsyntax.Block) []DependencySpec {
	var dependencies []DependencySpec
	for _, task := range group.Body.Blocks {
		if task.Type != "task" || len(task.Labels) == 0 || !isNomadDependencyTask(task) {
			continue
		}
		dependency := nomadDependencyFromTask(task)
		dependencies = appendUniqueDependency(dependencies, dependency)
	}
	return dependencies
}

func parseNomadGroupPlacement(group *hclsyntax.Block) *PlacementSpec {
	placement := &PlacementSpec{}
	for _, block := range group.Body.Blocks {
		switch block.Type {
		case "constraint":
			if constraint := nomadConstraintToPortable(block.Body); constraint != "" {
				appendUniqueString(&placement.Constraints, constraint)
			}
		case "affinity":
			if preference := nomadAffinityToPortable(block.Body); preference != "" {
				appendUniqueString(&placement.Preferences, preference)
			}
		case "spread":
			if preference := nomadSpreadToPortable(block.Body); preference != "" {
				appendUniqueString(&placement.Preferences, preference)
			}
		}
	}
	if len(placement.Constraints) == 0 && len(placement.Preferences) == 0 {
		return nil
	}
	return placement
}

func parseNomadGroupSpreadExtensions(group *hclsyntax.Block) []map[string]interface{} {
	if group == nil {
		return nil
	}
	var spreads []map[string]interface{}
	for _, block := range group.Body.Blocks {
		if block.Type != "spread" {
			continue
		}
		spread := map[string]interface{}{}
		if attribute := stringAttribute(block.Body, "attribute"); attribute != "" {
			spread["attribute"] = attribute
			spread["portable"] = nomadPortableAttribute(attribute)
		}
		if weight := intAttribute(block.Body, "weight"); weight > 0 {
			spread["weight"] = weight
		}
		for _, targetBlock := range block.Body.Blocks {
			if targetBlock.Type != "target" {
				continue
			}
			target := map[string]interface{}{}
			if len(targetBlock.Labels) > 0 {
				target["value"] = targetBlock.Labels[0]
			}
			if percent := intAttribute(targetBlock.Body, "percent"); percent > 0 {
				target["percent"] = percent
			}
			if len(target) > 0 {
				spread["targets"] = appendExtensionSlice(spread["targets"], target)
			}
		}
		if len(spread) > 0 {
			spreads = append(spreads, spread)
		}
	}
	return spreads
}

func parseNomadGroupPortableMeta(group *hclsyntax.Block) map[string]string {
	meta := map[string]string{}
	for _, block := range group.Body.Blocks {
		if block.Type != "meta" {
			continue
		}
		for key, attr := range block.Body.Attributes {
			if val := expressionString(attr.Expr); val != "" {
				meta[key] = val
			}
		}
	}
	if len(meta) == 0 {
		return nil
	}
	return meta
}

func parseNomadSchedulerBlockExtensions(body *hclsyntax.Body, blockType string) map[string]interface{} {
	if body == nil || blockType == "" {
		return nil
	}
	extensions := map[string]interface{}{}
	for _, block := range body.Blocks {
		if block.Type != blockType {
			continue
		}
		for key, value := range nomadAttributesToMap(block.Body) {
			extensions[key] = value
		}
	}
	if len(extensions) == 0 {
		return nil
	}
	return extensions
}

func parseNomadRestartBlockExtensions(body *hclsyntax.Body) map[string]interface{} {
	if body == nil {
		return nil
	}
	restart := map[string]interface{}{}
	for _, block := range body.Blocks {
		if block.Type != "restart" {
			continue
		}
		for key, value := range nomadAttributesToMap(block.Body) {
			restart[key] = value
		}
	}
	if len(restart) == 0 {
		return nil
	}
	return restart
}

func nomadJSONSchedulerBlockExtensions(value interface{}) map[string]interface{} {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]interface{}:
		if len(typed) == 0 {
			return nil
		}
		return cloneMap(typed)
	case []interface{}:
		var merged map[string]interface{}
		for _, item := range typed {
			if mapped, ok := asMap(item); ok && len(mapped) > 0 {
				merged = mergeNomadSchedulerBlocks(merged, mapped)
			}
		}
		if len(merged) == 0 {
			return nil
		}
		return merged
	default:
		if mapped, ok := asMap(typed); ok && len(mapped) > 0 {
			return cloneMap(mapped)
		}
		return nil
	}
}

func mergeNomadSchedulerBlocks(base, overlay map[string]interface{}) map[string]interface{} {
	if len(base) == 0 && len(overlay) == 0 {
		return nil
	}
	merged := map[string]interface{}{}
	for key, value := range base {
		merged[key] = deepCopyValue(value)
	}
	for key, value := range overlay {
		merged[key] = deepCopyValue(value)
	}
	return merged
}

func applyNomadSchedulerBlockExtensions(service *Service, update, migrate, reschedule map[string]interface{}) {
	if service == nil {
		return
	}
	if len(update) > 0 {
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions[nomadUpdateExtensionKey] = cloneMap(update)
		applyNomadUpdatePortableFields(service, update)
	}
	if len(migrate) > 0 {
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions[nomadMigrateExtensionKey] = cloneMap(migrate)
		deploy := ensureServiceDeploy(service)
		deploy.MigrateConfig = parseNomadMigratePolicyMap(migrate)
	}
	if len(reschedule) > 0 {
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions[nomadRescheduleExtensionKey] = cloneMap(reschedule)
		deploy := ensureServiceDeploy(service)
		deploy.RescheduleConfig = parseNomadReschedulePolicyMap(reschedule)
	}
}

func applyNomadRestartExtensions(service *Service, restart map[string]interface{}) {
	if service == nil || len(restart) == 0 {
		return
	}
	if service.Extensions == nil {
		service.Extensions = map[string]interface{}{}
	}
	service.Extensions[nomadRestartExtensionKey] = cloneMap(restart)
	deploy := ensureServiceDeploy(service)
	if attempts, ok := restart["attempts"]; ok {
		deploy.RestartPolicy = ensureRestartPolicy(deploy.RestartPolicy)
		deploy.RestartPolicy.MaxAttempts = toInt(attempts)
	}
	if delay, ok := restart["delay"]; ok {
		deploy.RestartPolicy = ensureRestartPolicy(deploy.RestartPolicy)
		deploy.RestartPolicy.Delay = toString(delay)
	}
	if interval, ok := restart["interval"]; ok {
		deploy.RestartPolicy = ensureRestartPolicy(deploy.RestartPolicy)
		deploy.RestartPolicy.Window = toString(interval)
	}
	if deploy.RestartPolicy != nil && isEmptyRestartPolicy(deploy.RestartPolicy) {
		deploy.RestartPolicy = nil
	}
}

func applyNomadUpdatePortableFields(service *Service, update map[string]interface{}) {
	if service == nil || len(update) == 0 {
		return
	}
	deploy := ensureServiceDeploy(service)
	if deploy.UpdateConfig == nil {
		deploy.UpdateConfig = &UpdatePolicy{}
	}
	if deploy.UpdateConfig.Extensions == nil {
		deploy.UpdateConfig.Extensions = map[string]interface{}{}
	}
	existingExactStrategy, hasExactStrategy := deploy.UpdateConfig.Extensions["kubernetes-deployment-strategy"]
	if !hasExactStrategy {
		existingExactStrategy, hasExactStrategy = deploy.UpdateConfig.Extensions["x-kubernetes-deployment-strategy"]
	}
	if parallelism, ok := update["max_parallel"]; ok {
		deploy.UpdateConfig.Parallelism = toInt(parallelism)
		deploy.UpdateConfig.ParallelismSet = true
	} else if value, ok := update["parallelism"]; ok {
		deploy.UpdateConfig.Parallelism = toInt(value)
		deploy.UpdateConfig.ParallelismSet = true
	}
	for key, value := range update {
		switch key {
		case "max_parallel", "parallelism":
			continue
		case "delay":
			deploy.UpdateConfig.Delay = toString(value)
		case "monitor":
			deploy.UpdateConfig.Monitor = toString(value)
		case "max_failure_ratio":
			deploy.UpdateConfig.MaxFailureRatio = toString(value)
		case "order":
			deploy.UpdateConfig.Order = toString(value)
		case "failure_action":
			deploy.UpdateConfig.OnFailure = toString(value)
		case "health_check":
			deploy.UpdateConfig.HealthCheck = toString(value)
		case "min_healthy_time":
			deploy.UpdateConfig.MinHealthyTime = toString(value)
		case "healthy_deadline":
			deploy.UpdateConfig.HealthyDeadline = toString(value)
		case "progress_deadline":
			deploy.UpdateConfig.ProgressDeadline = toString(value)
		case "auto_revert":
			deploy.UpdateConfig.AutoRevert = toBool(value)
			deploy.UpdateConfig.AutoRevertSet = true
		case "auto_promote":
			deploy.UpdateConfig.AutoPromote = toBool(value)
			deploy.UpdateConfig.AutoPromoteSet = true
		case "canary":
			deploy.UpdateConfig.Canary = toInt(value)
			deploy.UpdateConfig.CanarySet = true
		case "stagger":
			deploy.UpdateConfig.Stagger = toString(value)
		case "kubernetes-deployment-strategy", "x-kubernetes-deployment-strategy":
			if hasExactStrategy {
				deploy.UpdateConfig.Extensions["kubernetes-deployment-strategy"] = deepCopyValue(existingExactStrategy)
				continue
			}
			deploy.UpdateConfig.Extensions[key] = deepCopyValue(value)
		default:
			deploy.UpdateConfig.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(deploy.UpdateConfig.Extensions) == 0 {
		deploy.UpdateConfig.Extensions = nil
	}
}

func ensureRestartPolicy(policy *RestartPolicy) *RestartPolicy {
	if policy != nil {
		return policy
	}
	return &RestartPolicy{}
}

func applyNomadPortableMeta(service *Service, meta map[string]string) {
	if service == nil || len(meta) == 0 {
		return
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_build", "bolabaden.build"); raw != "" {
		var build BuildConfig
		if json.Unmarshal([]byte(raw), &build) == nil && buildConfigHasData(&build) {
			service.Build = &build
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_devices", "bolabaden.devices"); raw != "" {
		if devices := parseNomadJSONStringSlice(raw); len(devices) > 0 {
			service.Devices = devices
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_device_mappings", "bolabaden.device_mappings"); raw != "" {
		var mappings []DeviceMappingSpec
		if json.Unmarshal([]byte(raw), &mappings) == nil && len(mappings) > 0 {
			service.DeviceMappings = cloneDeviceMappings(mappings)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_expose", "bolabaden.expose"); raw != "" {
		if expose := parseNomadJSONStringSlice(raw); len(expose) > 0 {
			service.Expose = expose
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_healthcheck", "bolabaden.healthcheck"); raw != "" {
		var health HealthCheck
		if json.Unmarshal([]byte(raw), &health) == nil && !isEmptyHealthCheck(&health) {
			service.HealthCheck = mergeHealthCheckSpec(service.HealthCheck, &health)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_develop", "bolabaden.develop"); raw != "" {
		var develop DevelopConfig
		if json.Unmarshal([]byte(raw), &develop) == nil && !isEmptyDevelopConfig(&develop) {
			service.Develop = &develop
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_logging", "bolabaden.logging"); raw != "" {
		var logging struct {
			Driver     string                 `json:"driver"`
			Options    map[string]string      `json:"options"`
			Extensions map[string]interface{} `json:"extensions"`
		}
		if json.Unmarshal([]byte(raw), &logging) == nil {
			service.LogDriver = logging.Driver
			service.LogOpt = copyStringMap(logging.Options)
			service.LogExtensions = copyStringInterfaceMap(logging.Extensions)
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadComposeCompatMetaKey, "bolabaden.compose_compat"); raw != "" {
		var mapped map[string]interface{}
		if json.Unmarshal([]byte(raw), &mapped) == nil {
			if compat := composeCompatFromExtensionMap(mapped); compat != nil {
				service.ComposeCompat = compat
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadLinksMetaKey, "bolabaden.links"); raw != "" {
		if links := parseNomadJSONStringSlice(raw); len(links) > 0 {
			service.Links = links
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadVolumesMetaKey, "bolabaden.volumes"); raw != "" {
		var volumes []VolumeMount
		if json.Unmarshal([]byte(raw), &volumes) == nil && len(volumes) > 0 {
			service.Volumes = mergePortableVolumeMounts(service.Volumes, volumes)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_lifecycle", "bolabaden.lifecycle"); raw != "" {
		var lifecycle LifecycleHooks
		if json.Unmarshal([]byte(raw), &lifecycle) == nil && !isEmptyLifecycleHooks(&lifecycle) {
			service.Lifecycle = &lifecycle
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_ports", "bolabaden.ports"); raw != "" {
		var ports []PortMapping
		if json.Unmarshal([]byte(raw), &ports) == nil && len(ports) > 0 {
			service.Ports = mergePortablePorts(service.Ports, ports)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_dependencies", "bolabaden.dependencies"); raw != "" {
		var dependencies []DependencySpec
		if json.Unmarshal([]byte(raw), &dependencies) == nil && len(dependencies) > 0 {
			for _, dependency := range dependencies {
				service.Dependencies = appendUniqueDependency(service.Dependencies, dependency)
				service.DependsOn = appendUniqueName(service.DependsOn, dependency.Name)
			}
		}
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_deploy_endpoint_mode", "bolabaden.deploy_endpoint_mode"); value != "" {
		deploy := ensureServiceDeploy(service)
		deploy.EndpointMode = value
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_deploy_spec", "bolabaden.deploy_spec"); raw != "" {
		var deploy DeploySpec
		if json.Unmarshal([]byte(raw), &deploy) == nil && !isEmptyDeploySpec(&deploy) {
			service.Deploy = cloneDeploySpec(&deploy)
		}
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_deploy_mode", "bolabaden.deploy_mode"); value != "" {
		deploy := ensureServiceDeploy(service)
		deploy.Mode = value
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_deploy_labels", "bolabaden.deploy_labels"); raw != "" {
		var labels map[string]string
		if json.Unmarshal([]byte(raw), &labels) == nil && len(labels) > 0 {
			deploy := ensureServiceDeploy(service)
			deploy.Labels = copyStringMap(labels)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_deploy_resources", "bolabaden.deploy_resources"); raw != "" {
		var resources ResourceSpec
		if json.Unmarshal([]byte(raw), &resources) == nil {
			deploy := ensureServiceDeploy(service)
			deploy.Resources = mergeResourceSpec(deploy.Resources, &resources)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_deploy_job", "bolabaden.deploy_job"); raw != "" {
		var job SwarmJobSpec
		if json.Unmarshal([]byte(raw), &job) == nil && !isEmptySwarmJobSpec(&job) {
			deploy := ensureServiceDeploy(service)
			deploy.Job = &job
		}
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_deploy_max_replicas_per_node", "bolabaden.deploy_max_replicas_per_node"); value != "" {
		deploy := ensureServiceDeploy(service)
		if deploy.Placement == nil {
			deploy.Placement = &PlacementSpec{}
		}
		deploy.Placement.MaxReplicasPerNode = parseInt(value)
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_deploy_update_config", "bolabaden.deploy_update_config"); raw != "" {
		var update UpdatePolicy
		if json.Unmarshal([]byte(raw), &update) == nil && !isEmptyUpdatePolicy(&update) {
			deploy := ensureServiceDeploy(service)
			deploy.UpdateConfig = &update
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_deploy_rollback_config", "bolabaden.deploy_rollback_config"); raw != "" {
		var rollback UpdatePolicy
		if json.Unmarshal([]byte(raw), &rollback) == nil && !isEmptyUpdatePolicy(&rollback) {
			deploy := ensureServiceDeploy(service)
			deploy.RollbackConfig = &rollback
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_deploy_restart_policy", "bolabaden.deploy_restart_policy"); raw != "" {
		var restart RestartPolicy
		if json.Unmarshal([]byte(raw), &restart) == nil && !isEmptyRestartPolicy(&restart) {
			deploy := ensureServiceDeploy(service)
			deploy.RestartPolicy = &restart
		}
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_cpu_shares", "bolabaden.cpu_shares"); value != "" {
		service.CPUShares = parseInt(value)
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_cpu_quota", "bolabaden.cpu_quota"); value != "" {
		service.CPUQuota = parseInt(value)
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_mem_limit", "bolabaden.mem_limit"); value != "" {
		service.MemLimit = value
		service.MemoryLimit = value
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_memory_swap", "bolabaden.memory_swap"); value != "" {
		service.MemorySwap = value
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_mem_reservation", "bolabaden.mem_reservation"); value != "" {
		service.MemReservation = value
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_cpus", "bolabaden.cpus"); value != "" {
		service.CPUs = value
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_ulimits", "bolabaden.ulimits"); raw != "" {
		var limits Ulimits
		if json.Unmarshal([]byte(raw), &limits) == nil && limits.Nofile != nil {
			service.Ulimits = &limits
		}
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_userns_mode", "bolabaden.userns_mode"); value != "" {
		service.UserNSMode = value
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_pull_policy", "bolabaden.pull_policy"); value != "" {
		service.PullPolicy = value
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_init", "bolabaden.init"); value != "" {
		init := strings.EqualFold(value, "true")
		service.Init = &init
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_tty", "bolabaden.tty"); value != "" {
		service.Tty = strings.EqualFold(value, "true")
		service.TtySet = true
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_stdin_open", "bolabaden.stdin_open"); value != "" {
		service.StdinOpen = strings.EqualFold(value, "true")
		service.StdinOpenSet = true
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_stop_signal", "bolabaden.stop_signal"); value != "" {
		service.StopSignal = value
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_user", "bolabaden.user"); value != "" {
		service.User = value
	}
	if value := nomadPortableMetaValue(meta, "bolabaden_group", "bolabaden.group"); value != "" {
		service.Group = value
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_profiles", "bolabaden.profiles"); raw != "" {
		if profiles := parseNomadJSONStringSlice(raw); len(profiles) > 0 {
			service.Profiles = profiles
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_env_files", "bolabaden.env_files"); raw != "" {
		var refs []EnvFileRef
		if json.Unmarshal([]byte(raw), &refs) == nil && len(refs) > 0 {
			service.EnvFileRefs = cloneEnvFileRefs(refs)
			service.EnvFile = envFilePaths(refs)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_configs", "bolabaden.configs"); raw != "" {
		var refs []FileRef
		if json.Unmarshal([]byte(raw), &refs) == nil && len(refs) > 0 {
			service.Configs = mergePortableFileRefs(service.Configs, refs, false)
		}
	}
	if raw := nomadPortableMetaValue(meta, "bolabaden_secrets", "bolabaden.secrets"); raw != "" {
		var refs []FileRef
		if json.Unmarshal([]byte(raw), &refs) == nil && len(refs) > 0 {
			service.Secrets = mergePortableFileRefs(service.Secrets, refs, false)
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadNetworkAttachmentsMetaKey, "bolabaden.network_attachments"); raw != "" {
		var attachments map[string]*NetworkAttachment
		if json.Unmarshal([]byte(raw), &attachments) == nil && len(attachments) > 0 {
			service.NetworkAttachments = cloneNetworkAttachments(attachments)
			for name := range service.NetworkAttachments {
				appendUniqueString(&service.Networks, name)
			}
			sort.Strings(service.Networks)
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadFailoverMetaKey, "bolabaden.failover"); raw != "" {
		var mapped map[string]interface{}
		if json.Unmarshal([]byte(raw), &mapped) == nil {
			if failover, err := failoverSpecFromMap(mapped); err == nil {
				service.Failover = cloneFailoverSpec(failover)
			}
		}
	}
	if raw := nomadPortableMetaValue(meta, nomadExtensionsMetaKey, "bolabaden.extensions"); raw != "" {
		var extensions map[string]interface{}
		if json.Unmarshal([]byte(raw), &extensions) == nil && len(extensions) > 0 {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			for key, value := range extensions {
				if key == nomadExtensionsMetaKey || key == nomadRawHCLExtension {
					continue
				}
				service.Extensions[key] = value
			}
			applyNomadOpaqueKubernetesExtensions(service, extensions)
		}
	}
}

func nomadPortableMetaValue(meta map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := meta[key]; value != "" {
			return value
		}
	}
	return ""
}

func portablePortsNeedNomadMeta(ports []PortMapping) bool {
	for _, port := range ports {
		if port.HostIP != "" || port.TargetName != "" || port.AppProtocol != "" || port.Mode != "" || len(port.Extensions) > 0 {
			return true
		}
		if port.HostPort != "" && port.ContainerPort != "" && port.HostPort != port.ContainerPort {
			return true
		}
	}
	return false
}

func mergePortablePorts(existing, portable []PortMapping) []PortMapping {
	if len(portable) == 0 {
		return existing
	}
	if len(existing) == 0 {
		return clonePortMappings(portable)
	}
	result := clonePortMappings(existing)
	for _, port := range portable {
		index := matchingPortIndex(result, port)
		if index < 0 {
			result = append(result, port)
			continue
		}
		result[index] = mergePortMapping(result[index], port)
	}
	return result
}

func matchingPortIndex(ports []PortMapping, candidate PortMapping) int {
	for index, existing := range ports {
		if candidate.Name != "" && existing.Name == candidate.Name {
			return index
		}
		if candidate.ContainerPort != "" && existing.ContainerPort == candidate.ContainerPort {
			return index
		}
		if candidate.HostPort != "" && existing.HostPort == candidate.HostPort {
			return index
		}
	}
	return -1
}

func mergePortMapping(existing, portable PortMapping) PortMapping {
	if portable.Name != "" {
		existing.Name = portable.Name
	}
	if portable.TargetName != "" {
		existing.TargetName = portable.TargetName
	}
	if portable.HostIP != "" {
		existing.HostIP = portable.HostIP
	}
	if portable.HostPort != "" {
		existing.HostPort = portable.HostPort
	}
	if portable.ContainerPort != "" {
		existing.ContainerPort = portable.ContainerPort
	}
	if portable.Protocol != "" {
		existing.Protocol = portable.Protocol
	}
	if portable.AppProtocol != "" {
		existing.AppProtocol = portable.AppProtocol
	}
	if portable.Mode != "" {
		existing.Mode = portable.Mode
	}
	if portable.NodePort != "" {
		existing.NodePort = portable.NodePort
	}
	if len(portable.Extensions) > 0 {
		existing.Extensions = mergeInterfaceMaps(existing.Extensions, portable.Extensions)
	}
	return existing
}

func parseNomadJSONStringSlice(raw string) []string {
	var result []string
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(raw), &result); err == nil {
		return result
	}
	return splitNomadMetaList(raw)
}

func splitNomadMetaList(raw string) []string {
	var result []string
	for _, part := range strings.FieldsFunc(raw, func(r rune) bool { return r == '\n' || r == ',' }) {
		if trimmed := strings.TrimSpace(strings.Trim(part, `"`)); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func applyNomadPlacement(service *Service, placement *PlacementSpec) {
	if placement == nil {
		return
	}
	deploy := ensureServiceDeploy(service)
	if deploy.Placement == nil {
		deploy.Placement = &PlacementSpec{}
	}
	for _, constraint := range placement.Constraints {
		appendUniqueString(&deploy.Placement.Constraints, constraint)
	}
	for _, preference := range placement.Preferences {
		appendUniqueString(&deploy.Placement.Preferences, preference)
	}
}

func applyNomadSpreadExtensions(service *Service, spreads []map[string]interface{}) {
	if service == nil || len(spreads) == 0 {
		return
	}
	if len(service.Spreads) == 0 {
		service.Spreads = nomadSpreadSpecsFromAny(spreads)
	}
	if service.Extensions == nil {
		service.Extensions = map[string]interface{}{}
	}
	service.Extensions[nomadSpreadExtensionKey] = cloneMapSlice(spreads)
}

func nomadConstraintToPortable(body *hclsyntax.Body) string {
	attribute := stringAttribute(body, "attribute")
	operator := stringAttribute(body, "operator")
	value := stringAttribute(body, "value")
	if operator == "" {
		operator = "=="
	}
	return portablePlacementExpression(attribute, operator, value)
}

func nomadAffinityToPortable(body *hclsyntax.Body) string {
	attribute := stringAttribute(body, "attribute")
	value := stringAttribute(body, "value")
	operator := stringAttribute(body, "operator")
	if operator == "" {
		operator = "=="
	}
	expression := portablePlacementExpression(attribute, operator, value)
	if expression == "" {
		return ""
	}
	return "prefer:" + expression
}

func nomadSpreadToPortable(body *hclsyntax.Body) string {
	attribute := nomadPortableAttribute(stringAttribute(body, "attribute"))
	if attribute == "" {
		return ""
	}
	return "spread=" + attribute
}

func portablePlacementExpression(attribute, operator, value string) string {
	attribute = nomadPortableAttribute(attribute)
	operator = nomadPortableOperator(operator)
	value = strings.TrimSpace(value)
	if attribute == "" {
		switch operator {
		case "distinct_hosts":
			if value == "" {
				value = "true"
			}
			return "nomad.distinct_hosts == " + value
		case "distinct_property":
			if value == "" {
				return ""
			}
			return "nomad.distinct_property == " + value
		default:
			return ""
		}
	}
	if value == "" && (operator == "is_set" || operator == "is_not_set") {
		return fmt.Sprintf("%s %s", attribute, operator)
	}
	if value == "" {
		return ""
	}
	return fmt.Sprintf("%s %s %s", attribute, operator, value)
}

func nomadPortableAttribute(attribute string) string {
	attribute = strings.TrimSpace(attribute)
	attribute = strings.TrimPrefix(attribute, "${")
	attribute = strings.TrimSuffix(attribute, "}")
	switch {
	case strings.HasPrefix(attribute, "meta."):
		return "node.labels." + strings.TrimPrefix(attribute, "meta.")
	case attribute == "node.unique.name":
		return "node.hostname"
	case attribute == "node.datacenter":
		return "node.datacenter"
	case strings.HasPrefix(attribute, "attr."):
		return "nomad." + attribute
	default:
		return attribute
	}
}

func nomadPortableOperator(operator string) string {
	switch strings.TrimSpace(operator) {
	case "", "=":
		return "=="
	default:
		return strings.TrimSpace(operator)
	}
}

func isNomadDependencyTask(task *hclsyntax.Block) bool {
	if len(task.Labels) == 0 || !strings.HasPrefix(task.Labels[0], "wait-for-") {
		return false
	}
	for _, block := range task.Body.Blocks {
		if block.Type != "lifecycle" {
			continue
		}
		return stringAttribute(block.Body, "hook") == "prestart" && !boolAttribute(block.Body, "sidecar")
	}
	return false
}

func nomadDependencyFromTask(task *hclsyntax.Block) DependencySpec {
	dependency := DependencySpec{
		Name:      strings.TrimPrefix(task.Labels[0], "wait-for-"),
		Condition: "service_started",
	}
	env := nomadEnvMap(task)
	if name := env["BOLABADEN_DEPENDENCY_NAME"]; name != "" {
		dependency.Name = name
	}
	if condition := env["BOLABADEN_DEPENDENCY_CONDITION"]; condition != "" {
		dependency.Condition = condition
	}
	if restart := env["BOLABADEN_DEPENDENCY_RESTART"]; restart != "" {
		dependency.Restart = strings.EqualFold(restart, "true")
	}
	if required := env["BOLABADEN_DEPENDENCY_REQUIRED"]; required != "" {
		dependency.Required = boolPtr(strings.EqualFold(required, "true"))
	}
	return dependency
}

func nomadEnvMap(task *hclsyntax.Block) map[string]string {
	env := map[string]string{}
	for _, block := range task.Body.Blocks {
		if block.Type != "env" {
			continue
		}
		for key, attr := range block.Body.Attributes {
			if val := expressionString(attr.Expr); val != "" {
				env[key] = val
			}
		}
	}
	return env
}

func parseNomadConfig(service *Service, body *hclsyntax.Body, ports map[string]PortMapping) {
	if image := stringAttribute(body, "image"); image != "" {
		service.Image = image
	}
	if command := stringAttribute(body, "command"); command != "" {
		service.Entrypoint = []string{command}
	}
	if args := stringListAttribute(body, "args"); len(args) > 0 {
		service.Command = append(service.Command, args...)
	}
	if workDir := stringAttribute(body, "work_dir"); workDir != "" {
		service.WorkingDir = workDir
	}
	if hostname := stringAttribute(body, "hostname"); hostname != "" {
		service.Hostname = hostname
	}
	if _, ok := body.Attributes["privileged"]; ok {
		service.Privileged = boolAttribute(body, "privileged")
		service.PrivilegedSet = true
	}
	if _, ok := body.Attributes["readonly_rootfs"]; ok {
		service.ReadOnlyRootFS = boolAttribute(body, "readonly_rootfs")
		service.ReadOnlyRootFSSet = true
	}
	if runtime := stringAttribute(body, "runtime"); runtime != "" {
		service.Runtime = runtime
	}
	if value := stringAttribute(body, "pid_mode"); value != "" {
		service.PIDMode = value
		if strings.EqualFold(value, "host") {
			service.HostPID = boolPtr(true)
		}
	}
	if value := stringAttribute(body, "ipc_mode"); value != "" {
		service.IPCMode = value
		if strings.EqualFold(value, "host") {
			service.HostIPC = boolPtr(true)
		}
	}
	if _, ok := body.Attributes["pids_limit"]; ok {
		service.PidsLimit = int64(intAttribute(body, "pids_limit"))
		service.pidsLimitSet = true
		if service.Deploy == nil {
			service.Deploy = &DeploySpec{}
		}
		if service.Deploy.Resources == nil {
			service.Deploy.Resources = &ResourceSpec{}
		}
		service.Deploy.Resources.PidsLimit = service.PidsLimit
		service.Deploy.Resources.pidsLimitSet = true
	}
	if _, ok := body.Attributes["shm_size"]; ok {
		service.ShmSize = int64(intAttribute(body, "shm_size"))
		service.shmSizeSet = true
	}
	service.CapAdd = append(service.CapAdd, stringListAttribute(body, "cap_add")...)
	service.GroupAdd = append(service.GroupAdd, stringListAttribute(body, "group_add")...)
	if sysctls := stringMapAttribute(body, "sysctl"); len(sysctls) > 0 {
		service.Sysctls = sysctls
	}
	service.CapDrop = append(service.CapDrop, stringListAttribute(body, "cap_drop")...)
	service.SecurityOpt = append(service.SecurityOpt, stringListAttribute(body, "security_opt")...)
	service.DNS = append(service.DNS, stringListAttribute(body, "dns_servers")...)
	service.DNSSearch = append(service.DNSSearch, stringListAttribute(body, "dns_search_domains")...)
	service.DNSOptions = append(service.DNSOptions, stringListAttribute(body, "dns_options")...)
	service.ExtraHosts = append(service.ExtraHosts, stringListAttribute(body, "extra_hosts")...)
	for _, block := range body.Blocks {
		if block.Type == "mount" {
			service.Volumes = append(service.Volumes, parseNomadMount(block.Body))
		}
	}
	if portLabels := stringListAttribute(body, "ports"); len(portLabels) > 0 {
		for _, label := range portLabels {
			if port, ok := ports[label]; ok {
				if port.Name == "" {
					port.Name = label
				}
				service.Ports = append(service.Ports, port)
			}
		}
	}
}

func parseNomadMount(body *hclsyntax.Body) VolumeMount {
	mount := VolumeMount{
		Type:     stringAttribute(body, "type"),
		Source:   stringAttribute(body, "source"),
		Target:   stringAttribute(body, "target"),
		ReadOnly: boolAttribute(body, "readonly"),
	}
	if mount.Type == "" {
		mount.Type = "bind"
	}
	for _, block := range body.Blocks {
		switch block.Type {
		case "bind_options":
			mount.Propagation = stringAttribute(block.Body, "propagation")
		case "volume_options":
			mount.NoCopy = boolAttribute(block.Body, "no_copy")
			for _, nested := range block.Body.Blocks {
				if nested.Type == "driver_config" {
					if driver := stringAttribute(nested.Body, "name"); driver != "" {
						mount.Options = ensureStringMap(mount.Options)
						mount.Options["driver"] = driver
					}
				}
			}
		case "tmpfs_options":
			if size := intAttribute(block.Body, "size"); size > 0 {
				mount.Options = ensureStringMap(mount.Options)
				mount.Options["size"] = fmt.Sprintf("%d", size)
			}
		}
	}
	return mount
}

func parseNomadGroupPorts(group *hclsyntax.Block) map[string]PortMapping {
	ports := map[string]PortMapping{}
	for _, block := range group.Body.Blocks {
		if block.Type != "network" {
			continue
		}
		for _, portBlock := range block.Body.Blocks {
			if portBlock.Type != "port" || len(portBlock.Labels) == 0 {
				continue
			}
			label := portBlock.Labels[0]
			target := intAttribute(portBlock.Body, "to")
			static := intAttribute(portBlock.Body, "static")
			if target == 0 {
				target = static
			}
			if target > 0 {
				host := target
				if static > 0 {
					host = static
				}
				ports[label] = PortMapping{
					Name:          label,
					HostPort:      fmt.Sprintf("%d", host),
					ContainerPort: fmt.Sprintf("%d", target),
					Protocol:      "tcp",
				}
			}
		}
	}
	return ports
}

func parseNomadResources(body *hclsyntax.Body) *ResourceSpec {
	resources := &ResourceSpec{}
	if cpu := intAttribute(body, "cpu"); cpu > 0 {
		resources.CPUReservation = fmt.Sprintf("%dm", cpu)
	}
	if memory := intAttribute(body, "memory"); memory > 0 {
		resources.MemoryReservation = fmt.Sprintf("%dMi", memory)
	}
	if resources.CPUReservation == "" && resources.MemoryReservation == "" {
		return nil
	}
	return resources
}

func parseNomadServiceBlock(service *Service, body *hclsyntax.Body) {
	if serviceName := stringAttribute(body, "name"); serviceName != "" {
		service.Labels["nomad.service.name"] = serviceName
	}
	if portLabel := stringAttribute(body, "port"); portLabel != "" {
		service.Labels["nomad.service.port"] = portLabel
	}
	for _, block := range body.Blocks {
		switch block.Type {
		case "check":
			if health := parseNomadCheckBlock(block.Body, stringAttribute(body, "port")); health != nil {
				service.HealthCheck = health
			}
		case "connect":
			if connect := parseNomadConnectBlock(block.Body); len(connect) > 0 {
				if service.Extensions == nil {
					service.Extensions = map[string]interface{}{}
				}
				service.Extensions[nomadConnectExtensionKey] = connect
			}
		case "restart":
			if restart := nomadAttributesToMap(block.Body, "attempts", "delay", "interval", "mode", "render_templates"); len(restart) > 0 {
				applyNomadRestartExtensions(service, restart)
			}
		}
	}
}

func parseNomadConnectBlock(body *hclsyntax.Body) map[string]interface{} {
	if body == nil {
		return nil
	}
	connect := nomadAttributesToMap(body)
	if connect == nil {
		connect = map[string]interface{}{}
	}
	for _, block := range body.Blocks {
		switch block.Type {
		case "sidecar_service":
			if sidecar := parseNomadSidecarServiceBlock(block.Body); len(sidecar) > 0 {
				connect["sidecar_service"] = sidecar
			}
		case "gateway":
			if gateway := nomadGenericBlockFromBody(block.Body, "native"); len(gateway) > 0 {
				if proxy, ok := asMap(gateway["proxy"]); ok && len(proxy) > 0 {
					gateway["proxy"] = cloneMap(proxy)
				} else if proxyList, ok := gateway["proxy"].([]interface{}); ok && len(proxyList) == 1 {
					if proxy, ok := asMap(proxyList[0]); ok && len(proxy) > 0 {
						gateway["proxy"] = cloneMap(proxy)
					}
				}
				connect["gateway"] = gateway
			}
		case "extensions":
			if extensions := nomadGenericBlockFromBody(block.Body); len(extensions) > 0 {
				connect["extensions"] = extensions
			}
		}
	}
	return connect
}

func parseNomadSidecarServiceBlock(body *hclsyntax.Body) map[string]interface{} {
	sidecar := nomadAttributesToMap(body, "native")
	if sidecar == nil {
		sidecar = map[string]interface{}{}
	}
	for _, block := range body.Blocks {
		switch block.Type {
		case "proxy":
			if proxy := parseNomadProxyBlock(block.Body); len(proxy) > 0 {
				sidecar["proxy"] = proxy
			}
		case "check":
			if check := parseNomadConnectCheckBlock(block.Body, ""); len(check) > 0 {
				sidecar["check"] = check
			}
		case "extensions":
			if extensions := nomadGenericBlockFromBody(block.Body); len(extensions) > 0 {
				sidecar["extensions"] = extensions
			}
		}
	}
	return sidecar
}

func parseNomadProxyBlock(body *hclsyntax.Body) map[string]interface{} {
	proxy := nomadAttributesToMap(body, "native")
	if proxy == nil {
		proxy = map[string]interface{}{}
	}
	for _, block := range body.Blocks {
		switch block.Type {
		case "upstreams":
			if upstream := parseNomadConnectUpstreamBlock(block.Body); len(upstream) > 0 {
				proxy["upstreams"] = appendExtensionSlice(proxy["upstreams"], upstream)
			}
		case "config":
			if config := parseNomadConnectConfigBlock(block.Body); len(config) > 0 {
				proxy["config"] = config
			}
		case "extensions":
			if extensions := nomadGenericBlockFromBody(block.Body); len(extensions) > 0 {
				proxy["extensions"] = extensions
			}
		}
	}
	return proxy
}

func parseNomadConnectCheckBlock(body *hclsyntax.Body, servicePort string) map[string]interface{} {
	if body == nil {
		return nil
	}
	check := nomadAttributesToMap(body)
	if check == nil {
		check = map[string]interface{}{}
	}
	if servicePort != "" {
		if _, ok := check["port"]; !ok {
			check["port"] = servicePort
		}
	}
	for _, block := range body.Blocks {
		if block.Type != "extensions" {
			continue
		}
		if extensions := nomadGenericBlockFromBody(block.Body); len(extensions) > 0 {
			check["extensions"] = extensions
		}
	}
	if len(check) == 0 {
		return nil
	}
	return check
}

func parseNomadConnectConfigBlock(body *hclsyntax.Body) map[string]interface{} {
	if body == nil {
		return nil
	}
	config := nomadAttributesToMap(body)
	if config == nil {
		config = map[string]interface{}{}
	}
	for _, block := range body.Blocks {
		if block.Type != "extensions" {
			continue
		}
		if extensions := nomadGenericBlockFromBody(block.Body); len(extensions) > 0 {
			config["extensions"] = extensions
		}
	}
	if len(config) == 0 {
		return nil
	}
	return config
}

func parseNomadConnectUpstreamBlock(body *hclsyntax.Body) map[string]interface{} {
	if body == nil {
		return nil
	}
	upstream := nomadAttributesToMap(body, "native")
	if upstream == nil {
		upstream = map[string]interface{}{}
	}
	for _, block := range body.Blocks {
		if block.Type == "extensions" {
			if extensions := nomadGenericBlockFromBody(block.Body); len(extensions) > 0 {
				upstream["extensions"] = extensions
			}
		}
	}
	if len(upstream) == 0 {
		return nil
	}
	return upstream
}

func nomadGenericBlockFromBody(body *hclsyntax.Body, skip ...string) map[string]interface{} {
	if body == nil {
		return nil
	}
	result := nomadAttributesToMap(body, skip...)
	if result == nil {
		result = map[string]interface{}{}
	}
	for _, block := range body.Blocks {
		child := nomadGenericBlockFromBody(block.Body)
		if len(block.Labels) > 0 {
			if child == nil {
				child = map[string]interface{}{}
			}
			child["labels"] = append([]string{}, block.Labels...)
		}
		if len(child) == 0 {
			child = map[string]interface{}{}
		}
		if block.Type == "extensions" {
			merged := map[string]interface{}{}
			if existing, ok := asMap(result["extensions"]); ok {
				merged = cloneMap(existing)
			}
			for key, value := range child {
				merged[key] = deepCopyValue(value)
			}
			result["extensions"] = merged
			continue
		}
		result[block.Type] = appendExtensionSlice(result[block.Type], child)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func nomadAttributesToMap(body *hclsyntax.Body, skip ...string) map[string]interface{} {
	if body == nil || len(body.Attributes) == 0 {
		return nil
	}
	skipSet := map[string]bool{}
	for _, key := range skip {
		skipSet[key] = true
	}
	result := map[string]interface{}{}
	for key, attr := range body.Attributes {
		if skipSet[key] {
			continue
		}
		if value, ok := nomadAttributeValue(attr.Expr); ok {
			result[key] = value
		}
	}
	return result
}

func nomadAttributeValue(expr hclsyntax.Expression) (interface{}, bool) {
	value, diags := expr.Value(nil)
	if diags.HasErrors() {
		if text := templateTraversalString(expr); text != "" {
			return text, true
		}
		return nil, false
	}
	switch {
	case value.Type() == cty.String:
		return value.AsString(), true
	case value.Type() == cty.Bool:
		return value.True(), true
	case value.Type() == cty.Number:
		var i int
		if err := goctyInt(value, &i); err == nil {
			return i, true
		}
		return value.GoString(), true
	case value.CanIterateElements():
		var result []interface{}
		it := value.ElementIterator()
		for it.Next() {
			_, val := it.Element()
			if val.Type() == cty.String {
				result = append(result, val.AsString())
			} else if val.Type() == cty.Bool {
				result = append(result, val.True())
			} else if val.Type() == cty.Number {
				var i int
				if err := goctyInt(val, &i); err == nil {
					result = append(result, i)
				}
			}
		}
		if len(result) > 0 {
			return result, true
		}
	}
	return nil, false
}

func parseNomadCheckBlock(body *hclsyntax.Body, servicePort string) *HealthCheck {
	health := &HealthCheck{
		Type:     stringAttribute(body, "type"),
		Path:     stringAttribute(body, "path"),
		Port:     stringAttribute(body, "port"),
		Interval: stringAttribute(body, "interval"),
		Timeout:  stringAttribute(body, "timeout"),
	}
	if health.Port == "" {
		health.Port = servicePort
	}
	if health.Type == "script" {
		health.Type = "exec"
		command := stringAttribute(body, "command")
		args := stringListAttribute(body, "args")
		if command != "" {
			health.Test = append([]string{"CMD", command}, args...)
		}
	}
	if health.Type == "" && health.Path != "" {
		health.Type = "http"
	}
	if health.Type == "" && health.Port != "" {
		health.Type = "tcp"
	}
	if health.Type == "" && health.Path == "" && health.Port == "" && health.Interval == "" && health.Timeout == "" && len(health.Test) == 0 {
		return nil
	}
	return normalizeHealthCheck(health)
}

func parseNomadTemplateBlock(service *Service, body *hclsyntax.Body) {
	data := stringAttribute(body, "data")
	if data == "" {
		return
	}
	fields, ok := parseNomadBolabadenMarker(data)
	if !ok {
		return
	}
	switch fields["kind"] {
	case "file_ref":
		ref := FileRef{
			Source: fields["source"],
			Key:    fields["key"],
			Target: stringAttribute(body, "destination"),
			Mode:   fields["mode"],
		}
		if ref.Source == "" {
			return
		}
		if fields["read_only"] != "" {
			ref.ReadOnly = strings.EqualFold(fields["read_only"], "true")
		}
		if fields["optional"] != "" {
			ref.Optional = boolPtr(strings.EqualFold(fields["optional"], "true"))
		}
		switch fields["type"] {
		case "config":
			service.Configs = append(service.Configs, ref)
		case "secret":
			service.Secrets = append(service.Secrets, ref)
		}
	case "env_source":
		source := EnvSource{
			Name:       fields["name"],
			SourceType: fields["type"],
			Source:     fields["source"],
			Key:        fields["key"],
			Optional:   strings.EqualFold(fields["optional"], "true"),
		}
		if extensions := nomadExtensionsFromMarker(fields["extensions_b64"]); len(extensions) > 0 {
			source.Extensions = extensions
		}
		if source.Name != "" && source.Source != "" {
			service.EnvSources = append(service.EnvSources, source)
		}
	case "env_from":
		source := EnvFromSource{
			SourceType: fields["type"],
			Source:     fields["source"],
			Prefix:     fields["prefix"],
			Optional:   strings.EqualFold(fields["optional"], "true"),
		}
		if extensions := nomadExtensionsFromMarker(fields["extensions_b64"]); len(extensions) > 0 {
			source.Extensions = extensions
		}
		if source.Source != "" {
			service.EnvFrom = append(service.EnvFrom, source)
		}
	}
}

func parseNomadBolabadenMarker(data string) (map[string]string, bool) {
	const prefix = "bolabaden:"
	start := strings.Index(data, prefix)
	if start < 0 {
		return nil, false
	}
	line := data[start+len(prefix):]
	if idx := strings.IndexAny(line, "\r\n"); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(strings.TrimSuffix(line, "*/}}"))
	if line == "" {
		return nil, false
	}
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, false
	}
	fields := map[string]string{"kind": parts[0]}
	for _, part := range parts[1:] {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		fields[key] = strings.Trim(value, `"`)
	}
	return fields, true
}

// SerializeNomadHCL converts an Application to Nomad HCL.
func SerializeNomadHCL(app *Application) (string, error) {
	emitApp := cloneApplication(app)
	syncPortableApplicationState(emitApp)
	if emitApp != nil && emitApp.Platform == PlatformNomad {
		if raw := toString(emitApp.Extensions[nomadRawHCLExtension]); strings.TrimSpace(raw) != "" && !json.Valid([]byte(raw)) {
			return raw, nil
		}
	}

	var hcl strings.Builder

	hcl.WriteString("# Auto-generated Nomad HCL from PaaS converter\n\n")
	hcl.WriteString("job \"app\" {\n")
	hcl.WriteString("  datacenters = [\"dc1\"]\n")
	hcl.WriteString(fmt.Sprintf("  type = %q\n\n", nomadJobType(emitApp)))

	writeNomadApplicationMeta(&hcl, emitApp)

	for serviceName, service := range emitApp.Services {
		count := service.Replicas
		if service.Deploy != nil && service.Deploy.Replicas > 0 {
			count = service.Deploy.Replicas
		}
		if count == 0 {
			count = 1
		}

		hcl.WriteString(fmt.Sprintf("  group \"%s\" {\n", serviceName))
		if !nomadServiceIsGlobal(service) {
			hcl.WriteString(fmt.Sprintf("    count = %d\n\n", count))
		} else {
			hcl.WriteString("\n")
		}

		writeNomadPortableMeta(&hcl, service)
		writeNomadNetwork(&hcl, service)
		writeNomadPlacement(&hcl, service)
		writeNomadSchedulerBlocks(&hcl, service)
		writeNomadDependencyTasks(&hcl, service)

		hcl.WriteString(fmt.Sprintf("    task \"%s\" {\n", serviceName))
		driver := "docker"
		if value := strings.TrimSpace(toString(service.Extensions["nomad.driver"])); value != "" {
			driver = value
		}
		hcl.WriteString(fmt.Sprintf("      driver = %q\n\n", driver))
		if service.User != "" {
			hcl.WriteString(fmt.Sprintf("      user = %q\n\n", service.User))
		}
		if service.StopSignal != "" {
			hcl.WriteString(fmt.Sprintf("      kill_signal = %q\n", service.StopSignal))
		}
		if service.StopGracePeriod != "" {
			hcl.WriteString(fmt.Sprintf("      kill_timeout = %q\n", service.StopGracePeriod))
		}
		writeNomadRestart(&hcl, service)
		hcl.WriteString("      config {\n")
		if service.Image != "" {
			hcl.WriteString(fmt.Sprintf("        image = %q\n", service.Image))
		}
		writeNomadRuntimeConfig(&hcl, service)
		writeNomadVolumeMounts(&hcl, service)
		if len(service.Ports) > 0 {
			var labels []string
			for i := range service.Ports {
				labels = append(labels, fmt.Sprintf("%q", nomadPortLabelForPort(service.Ports[i], i)))
			}
			hcl.WriteString(fmt.Sprintf("        ports = [%s]\n", strings.Join(labels, ", ")))
		}
		hcl.WriteString("      }\n\n")

		if len(service.Environment) > 0 {
			hcl.WriteString("      env {\n")
			for _, key := range sortedMapKeys(service.Environment) {
				value := service.Environment[key]
				hcl.WriteString(fmt.Sprintf("        %s = %q\n", key, value))
			}
			hcl.WriteString("      }\n\n")
		}
		writeNomadTaskMeta(&hcl, service)
		writeNomadTemplates(&hcl, service)

		if hasNomadNativeResources(service.Deploy) {
			hcl.WriteString("      resources {\n")
			if cpu := nomadCPU(service.Deploy.Resources.CPUReservation); cpu != "" {
				hcl.WriteString(fmt.Sprintf("        cpu = %s\n", cpu))
			}
			if memory := nomadMemory(service.Deploy.Resources.MemoryReservation); memory != "" {
				hcl.WriteString(fmt.Sprintf("        memory = %s\n", memory))
			}
			hcl.WriteString("      }\n")
		}
		writeNomadService(&hcl, serviceName, service)

		hcl.WriteString("    }\n")
		hcl.WriteString("  }\n\n")
	}

	hcl.WriteString("}\n")

	return hcl.String(), nil
}

func writeNomadRestart(hcl *strings.Builder, service *Service) {
	if hcl == nil || service == nil {
		return
	}
	restart := nomadRestartBlockForService(service)
	if len(restart) == 0 {
		return
	}
	hcl.WriteString("      restart {\n")
	emitNomadSchedulerBlockAttributes(hcl, restart, []string{
		"attempts",
		"delay",
		"interval",
		"mode",
		"render_templates",
	}, "        ")
	hcl.WriteString("      }\n\n")
}

func nomadRestartBlockForService(service *Service) map[string]interface{} {
	if service == nil {
		return nil
	}
	if value, ok := firstExtensionValue(service.Extensions, nomadRestartExtensionKey, "x-nomad-restart"); ok {
		if restart, ok := asMap(value); ok && len(restart) > 0 {
			return cloneMap(restart)
		}
	}
	if service.Deploy == nil || service.Deploy.RestartPolicy == nil {
		return nil
	}
	restart := map[string]interface{}{}
	if service.Deploy.RestartPolicy.MaxAttempts > 0 {
		restart["attempts"] = service.Deploy.RestartPolicy.MaxAttempts
	}
	if service.Deploy.RestartPolicy.Delay != "" {
		restart["delay"] = service.Deploy.RestartPolicy.Delay
	}
	if service.Deploy.RestartPolicy.Window != "" {
		restart["interval"] = service.Deploy.RestartPolicy.Window
	}
	if len(restart) == 0 {
		return nil
	}
	return restart
}

// SerializeNomadJSON converts an Application to Nomad JSON.
func SerializeNomadJSON(app *Application) (string, error) {
	syncPortableApplicationState(app)
	if app != nil && app.Platform == PlatformNomad {
		if raw := toString(app.Extensions[nomadRawHCLExtension]); strings.TrimSpace(raw) != "" && json.Valid([]byte(raw)) {
			return raw, nil
		}
	}
	hclSource, err := SerializeNomadHCL(app)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(hclSource) == "" {
		return "", nil
	}
	jsonBytes, err := convert.Bytes([]byte(hclSource), "nomad.hcl", convert.Options{Simplify: true})
	if err != nil {
		return "", fmt.Errorf("failed to convert Nomad HCL to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

func nomadJobType(app *Application) string {
	if app == nil || len(app.Services) == 0 {
		return "service"
	}
	for _, service := range app.Services {
		if nomadServiceIsBatch(service) {
			continue
		}
		if !nomadServiceIsGlobal(service) {
			return "service"
		}
	}
	for _, service := range app.Services {
		if nomadServiceIsBatch(service) {
			return "batch"
		}
	}
	return "system"
}

func nomadServiceIsGlobal(service *Service) bool {
	return service != nil && service.Deploy != nil && strings.EqualFold(service.Deploy.Mode, "global")
}

func nomadServiceIsBatch(service *Service) bool {
	return service != nil && service.Deploy != nil && isSwarmJobMode(service.Deploy.Mode)
}

// RestoreNomadSource writes the preserved Nomad HCL source back to disk.
// It prefers the source extension when the original Nomad app is still intact,
// then falls back to the canonical raw resource preserved across conversions.
func RestoreNomadSource(app *Application, filename string) error {
	raw, err := nomadSourceContent(app)
	if err != nil {
		return err
	}
	if dir := filepath.Dir(filename); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create Nomad output directory: %w", err)
		}
	}
	if err := os.WriteFile(filename, []byte(raw), 0644); err != nil {
		return fmt.Errorf("failed to restore Nomad source: %w", err)
	}
	return nil
}

// RestoreNomadJSONSource writes a Nomad JSON rendering back to disk.
func RestoreNomadJSONSource(app *Application, filename string) error {
	raw, err := nomadJSONSourceContent(app)
	if err != nil {
		return err
	}
	if dir := filepath.Dir(filename); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create Nomad JSON output directory: %w", err)
		}
	}
	if err := os.WriteFile(filename, []byte(raw), 0644); err != nil {
		return fmt.Errorf("failed to restore Nomad JSON source: %w", err)
	}
	return nil
}

func nomadSourceContent(app *Application) (string, error) {
	raw := ""
	if app != nil && app.Extensions != nil {
		raw = toString(app.Extensions[nomadRawHCLExtension])
	}
	if strings.TrimSpace(raw) != "" && !json.Valid([]byte(raw)) {
		return raw, nil
	}
	if strings.TrimSpace(raw) == "" {
		if preserved := canonicalRawResourceValue(app, PlatformNomad, "NomadHCL"); preserved != nil {
			raw = toString(preserved)
		}
	}
	if strings.TrimSpace(raw) != "" && !json.Valid([]byte(raw)) {
		return raw, nil
	}
	serialized, err := SerializeNomadHCL(app)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(serialized) == "" {
		return "", fmt.Errorf("application does not contain raw Nomad HCL")
	}
	return serialized, nil
}

func nomadJSONSourceContent(app *Application) (string, error) {
	raw, err := SerializeNomadJSON(app)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("application does not contain renderable Nomad JSON")
	}
	return raw, nil
}

func writeNomadNetwork(hcl *strings.Builder, service *Service) {
	if len(service.Ports) == 0 {
		return
	}
	hcl.WriteString("    network {\n")
	for i, port := range service.Ports {
		label := nomadPortLabelForPort(port, i)
		target := port.ContainerPort
		if target == "" {
			target = port.HostPort
		}
		hcl.WriteString(fmt.Sprintf("      port \"%s\" {\n", label))
		hcl.WriteString(fmt.Sprintf("        to = %s\n", target))
		if port.HostPort != "" && port.HostPort != target {
			hcl.WriteString(fmt.Sprintf("        static = %s\n", port.HostPort))
		}
		hcl.WriteString("      }\n")
	}
	hcl.WriteString("    }\n\n")
}

func hasNomadNativeResources(deploy *DeploySpec) bool {
	return deploy != nil &&
		deploy.Resources != nil &&
		(deploy.Resources.CPUReservation != "" || deploy.Resources.MemoryReservation != "")
}

func writeNomadPortableMeta(hcl *strings.Builder, service *Service) {
	if service == nil {
		return
	}
	meta := map[string]string{}
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
	if buildConfigHasData(service.Build) {
		if raw, err := json.Marshal(service.Build); err == nil {
			meta["bolabaden_build"] = string(raw)
		}
	}
	if len(service.Devices) > 0 {
		if raw, err := json.Marshal(service.Devices); err == nil {
			meta["bolabaden_devices"] = string(raw)
		}
	}
	if len(service.DeviceMappings) > 0 {
		if raw, err := json.Marshal(service.DeviceMappings); err == nil {
			meta["bolabaden_device_mappings"] = string(raw)
		}
	}
	if len(service.Expose) > 0 {
		if raw, err := json.Marshal(service.Expose); err == nil {
			meta["bolabaden_expose"] = string(raw)
		}
	}
	if !isEmptyHealthCheck(service.HealthCheck) {
		if raw, err := json.Marshal(service.HealthCheck); err == nil {
			meta["bolabaden_healthcheck"] = string(raw)
		}
	}
	if !isEmptyDevelopConfig(service.Develop) {
		if raw, err := json.Marshal(service.Develop); err == nil {
			meta["bolabaden_develop"] = string(raw)
		}
	}
	if !isEmptyLifecycleHooks(service.Lifecycle) {
		if raw, err := json.Marshal(service.Lifecycle); err == nil {
			meta["bolabaden_lifecycle"] = string(raw)
		}
	}
	if service.Init != nil {
		meta["bolabaden_init"] = fmt.Sprintf("%t", *service.Init)
	}
	if service.TtySet {
		meta["bolabaden_tty"] = fmt.Sprintf("%t", service.Tty)
	} else if service.Tty {
		meta["bolabaden_tty"] = "true"
	}
	if service.StdinOpenSet {
		meta["bolabaden_stdin_open"] = fmt.Sprintf("%t", service.StdinOpen)
	} else if service.StdinOpen {
		meta["bolabaden_stdin_open"] = "true"
	}
	if service.StopSignal != "" {
		meta["bolabaden_stop_signal"] = service.StopSignal
	}
	if service.StopGracePeriod != "" {
		meta["bolabaden_stop_grace_period"] = service.StopGracePeriod
	}
	if service.User != "" {
		meta["bolabaden_user"] = service.User
	}
	if service.Group != "" {
		meta["bolabaden_group"] = service.Group
	}
	if portablePortsNeedNomadMeta(service.Ports) {
		if raw, err := json.Marshal(service.Ports); err == nil {
			meta["bolabaden_ports"] = string(raw)
		}
	}
	if dependencies := serviceDependencies(service); len(dependencies) > 0 {
		if raw, err := json.Marshal(dependencies); err == nil {
			meta["bolabaden_dependencies"] = string(raw)
		}
	}
	if deploy != nil && deploy.EndpointMode != "" {
		meta["bolabaden_deploy_endpoint_mode"] = deploy.EndpointMode
	}
	if deploy != nil && !isEmptyDeploySpec(deploy) {
		if raw, err := json.Marshal(deploy); err == nil {
			meta["bolabaden_deploy_spec"] = string(raw)
		}
	}
	if deploy != nil && deploy.Mode != "" {
		meta["bolabaden_deploy_mode"] = deploy.Mode
	}
	if deploy != nil && len(deploy.Labels) > 0 {
		if raw, err := json.Marshal(deploy.Labels); err == nil {
			meta["bolabaden_deploy_labels"] = string(raw)
		}
	}
	if deploy != nil && !isEmptyResourceSpec(deploy.Resources) {
		if raw, err := json.Marshal(deploy.Resources); err == nil {
			meta["bolabaden_deploy_resources"] = string(raw)
		}
	}
	if deploy != nil && !isEmptySwarmJobSpec(deploy.Job) {
		if raw, err := json.Marshal(deploy.Job); err == nil {
			meta["bolabaden_deploy_job"] = string(raw)
		}
	}
	if deploy != nil && deploy.Placement != nil && deploy.Placement.MaxReplicasPerNode > 0 {
		meta["bolabaden_deploy_max_replicas_per_node"] = fmt.Sprintf("%d", deploy.Placement.MaxReplicasPerNode)
	}
	if deploy != nil && !isEmptyUpdatePolicy(deploy.UpdateConfig) {
		if raw, err := json.Marshal(deploy.UpdateConfig); err == nil {
			meta["bolabaden_deploy_update_config"] = string(raw)
		}
	}
	if deploy != nil && !isEmptyUpdatePolicy(deploy.RollbackConfig) {
		if raw, err := json.Marshal(deploy.RollbackConfig); err == nil {
			meta["bolabaden_deploy_rollback_config"] = string(raw)
		}
	}
	if deploy != nil && !isEmptyRestartPolicy(deploy.RestartPolicy) {
		if raw, err := json.Marshal(deploy.RestartPolicy); err == nil {
			meta["bolabaden_deploy_restart_policy"] = string(raw)
		}
	}
	if service.CPUShares > 0 {
		meta["bolabaden_cpu_shares"] = fmt.Sprintf("%d", service.CPUShares)
	}
	if service.CPUQuota > 0 {
		meta["bolabaden_cpu_quota"] = fmt.Sprintf("%d", service.CPUQuota)
	}
	if service.MemLimit != "" {
		meta["bolabaden_mem_limit"] = service.MemLimit
	} else if service.MemoryLimit != "" {
		meta["bolabaden_mem_limit"] = service.MemoryLimit
	}
	if service.MemorySwap != "" {
		meta["bolabaden_memory_swap"] = service.MemorySwap
	}
	if service.MemReservation != "" {
		meta["bolabaden_mem_reservation"] = service.MemReservation
	}
	if service.CPUs != "" {
		meta["bolabaden_cpus"] = service.CPUs
	}
	if !isEmptyUlimits(service.Ulimits) {
		if raw, err := json.Marshal(service.Ulimits); err == nil {
			meta["bolabaden_ulimits"] = string(raw)
		}
	}
	if service.UserNSMode != "" {
		meta["bolabaden_userns_mode"] = service.UserNSMode
	}
	if service.LogDriver != "" || len(service.LogOpt) > 0 || len(service.LogExtensions) > 0 {
		if raw, err := json.Marshal(map[string]interface{}{
			"driver":     service.LogDriver,
			"options":    service.LogOpt,
			"extensions": service.LogExtensions,
		}); err == nil {
			meta["bolabaden_logging"] = string(raw)
		}
	}
	if !isEmptyComposeCompat(service.ComposeCompat) {
		if raw, err := json.Marshal(service.ComposeCompat); err == nil {
			meta[nomadComposeCompatMetaKey] = string(raw)
		}
	}
	if len(service.Links) > 0 {
		if raw, err := json.Marshal(service.Links); err == nil {
			meta[nomadLinksMetaKey] = string(raw)
		}
	}
	if service.PullPolicy != "" {
		meta["bolabaden_pull_policy"] = service.PullPolicy
	}
	if len(service.Profiles) > 0 {
		if raw, err := json.Marshal(service.Profiles); err == nil {
			meta["bolabaden_profiles"] = string(raw)
		}
	}
	if len(service.EnvFileRefs) > 0 {
		if raw, err := json.Marshal(service.EnvFileRefs); err == nil {
			meta["bolabaden_env_files"] = string(raw)
		}
	} else if len(service.EnvFile) > 0 {
		refs := make([]EnvFileRef, 0, len(service.EnvFile))
		for _, path := range service.EnvFile {
			refs = append(refs, EnvFileRef{Path: path})
		}
		if raw, err := json.Marshal(refs); err == nil {
			meta["bolabaden_env_files"] = string(raw)
		}
	}
	if len(service.Configs) > 0 {
		if raw, err := json.Marshal(service.Configs); err == nil {
			meta["bolabaden_configs"] = string(raw)
		}
	}
	if len(service.Secrets) > 0 {
		if raw, err := json.Marshal(service.Secrets); err == nil {
			meta["bolabaden_secrets"] = string(raw)
		}
	}
	if len(service.Volumes) > 0 {
		if raw, err := json.Marshal(service.Volumes); err == nil {
			meta[nomadVolumesMetaKey] = string(raw)
		}
	}
	if len(service.NetworkAttachments) > 0 {
		if raw, err := json.Marshal(service.NetworkAttachments); err == nil {
			meta[nomadNetworkAttachmentsMetaKey] = string(raw)
		}
	}
	if service.Failover != nil {
		if raw, err := json.Marshal(service.Failover); err == nil {
			meta[nomadFailoverMetaKey] = string(raw)
		}
	}
	extensions := jsonSerializableExtensionMap(service.Extensions)
	if kubernetesExtensions := nomadKubernetesExtensions(service); len(kubernetesExtensions) > 0 {
		if extensions == nil {
			extensions = map[string]interface{}{}
		}
		for key, value := range kubernetesExtensions {
			if _, exists := extensions[key]; !exists {
				extensions[key] = value
			}
		}
	}
	if len(extensions) > 0 {
		delete(extensions, nomadRawHCLExtension)
		if raw, err := json.Marshal(extensions); err == nil {
			meta[nomadExtensionsMetaKey] = string(raw)
		}
	}
	if len(meta) == 0 {
		return
	}
	hcl.WriteString("    meta {\n")
	for _, key := range sortedMapKeys(meta) {
		hcl.WriteString(fmt.Sprintf("      %s = %q\n", key, meta[key]))
	}
	hcl.WriteString("    }\n\n")
}

func nomadKubernetesExtensions(service *Service) map[string]interface{} {
	if service == nil {
		return nil
	}
	extensions := map[string]interface{}{}
	if len(service.ImagePullSecrets) > 0 {
		extensions["kubernetes.imagePullSecrets"] = append([]string{}, service.ImagePullSecrets...)
	}
	if service.ImagePullPolicy != "" {
		extensions["kubernetes.imagePullPolicy"] = service.ImagePullPolicy
	}
	if aliases := kubernetesHostAliasesFromService(service); len(aliases) > 0 {
		extensions["kubernetes.hostAliases"] = aliases
	}
	if service.DNSPolicy != "" {
		extensions["kubernetes.dnsPolicy"] = service.DNSPolicy
	}
	if service.SchedulerName != "" {
		extensions["kubernetes.schedulerName"] = service.SchedulerName
	}
	if service.TerminationMessagePath != "" {
		extensions["kubernetes.terminationMessagePath"] = service.TerminationMessagePath
	}
	if service.TerminationMessagePolicy != "" {
		extensions["kubernetes.terminationMessagePolicy"] = service.TerminationMessagePolicy
	}
	if service.HostNetworkSet {
		extensions["kubernetes.hostNetwork"] = service.HostNetwork
	} else if service.HostNetwork {
		extensions["kubernetes.hostNetwork"] = true
	}
	if service.HostPID != nil {
		extensions["kubernetes.hostPID"] = *service.HostPID
	}
	if service.HostIPC != nil {
		extensions["kubernetes.hostIPC"] = *service.HostIPC
	}
	if service.PriorityClassName != "" {
		extensions["kubernetes.priorityClassName"] = service.PriorityClassName
	}
	if service.RuntimeClassName != "" {
		extensions["kubernetes.runtimeClassName"] = service.RuntimeClassName
	}
	if service.NodeName != "" {
		extensions["kubernetes.nodeName"] = service.NodeName
	}
	if len(service.NodeSelector) > 0 {
		extensions["kubernetes.nodeSelector"] = copyStringMap(service.NodeSelector)
	}
	if service.Subdomain != "" {
		extensions["kubernetes.subdomain"] = service.Subdomain
	}
	if service.OSName != "" {
		extensions["kubernetes.os"] = service.OSName
	}
	if service.SetHostnameAsFQDN != nil {
		extensions["kubernetes.setHostnameAsFQDN"] = *service.SetHostnameAsFQDN
	}
	if service.HostUsers != nil {
		extensions["kubernetes.hostUsers"] = *service.HostUsers
	}
	if service.Group != "" {
		extensions["kubernetes.group"] = service.Group
	}
	if service.ShareProcessNamespace != nil {
		extensions["kubernetes.shareProcessNamespace"] = *service.ShareProcessNamespace
	}
	if service.EnableServiceLinks != nil {
		extensions["kubernetes.enableServiceLinks"] = *service.EnableServiceLinks
	}
	if service.ServiceAccountName != "" {
		extensions["kubernetes.serviceAccountName"] = service.ServiceAccountName
	}
	if service.AutomountServiceAccountToken != nil {
		extensions["kubernetes.automountServiceAccountToken"] = *service.AutomountServiceAccountToken
	}
	if service.FSGroup != nil {
		extensions["kubernetes.fsGroup"] = *service.FSGroup
	}
	if selinux := serializeKubernetesSELinuxOptions(service.SELinuxOptions); len(selinux) > 0 {
		extensions["kubernetes.seLinuxOptions"] = selinux
	}
	if windows := serializeKubernetesWindowsSecurityContextOptions(service.WindowsOptions); len(windows) > 0 {
		extensions["kubernetes.windowsOptions"] = windows
	}
	if service.FSGroupChangePolicy != "" {
		extensions["kubernetes.fsGroupChangePolicy"] = service.FSGroupChangePolicy
	}
	if service.RunAsNonRoot != nil {
		extensions["kubernetes.runAsNonRoot"] = *service.RunAsNonRoot
	}
	if len(service.SupplementalGroups) > 0 {
		extensions["kubernetes.supplementalGroups"] = append([]int64{}, service.SupplementalGroups...)
	}
	if service.SupplementalGroupsPolicy != "" {
		extensions["kubernetes.supplementalGroupsPolicy"] = service.SupplementalGroupsPolicy
	}
	if len(service.GroupAdd) > 0 {
		extensions["kubernetes.groupAdd"] = append([]string{}, service.GroupAdd...)
	}
	if service.Runtime != "" {
		extensions["kubernetes.runtime"] = service.Runtime
	}
	if service.PIDMode != "" {
		extensions["kubernetes.pidMode"] = service.PIDMode
	}
	if service.IPCMode != "" {
		extensions["kubernetes.ipcMode"] = service.IPCMode
	}
	if service.ActiveDeadlineSeconds != nil {
		extensions["kubernetes.activeDeadlineSeconds"] = *service.ActiveDeadlineSeconds
	}
	if service.PodRestartPolicy != "" {
		extensions["kubernetes.restartPolicy"] = service.PodRestartPolicy
	}
	if len(service.Tolerations) > 0 {
		extensions["kubernetes.tolerations"] = serializeKubernetesTolerations(service.Tolerations)
	}
	if service.StartupProbe != nil {
		if probe := serializeKubernetesProbeExtension(service.StartupProbe); len(probe) > 0 {
			extensions["x-kubernetes-startup-probe"] = probe
		}
	}
	if profile := serializeKubernetesSeccompProfileExtension(service.SeccompProfile); len(profile) > 0 {
		extensions["x-kubernetes-seccomp-profile"] = profile
	}
	if affinity := copyStringInterfaceMap(service.Affinity); len(affinity) > 0 {
		extensions["x-kubernetes-affinity"] = affinity
	}
	if len(service.ReadinessGates) > 0 {
		extensions["x-kubernetes-readiness-gates"] = cloneMapSlice(service.ReadinessGates)
	}
	if service.AllowPrivilegeEscalation != nil {
		extensions["x-kubernetes-allowPrivilegeEscalation"] = *service.AllowPrivilegeEscalation
	}
	if service.ProcMount != "" {
		extensions["x-kubernetes-procMount"] = service.ProcMount
	}
	if len(service.InitContainers) > 0 {
		extensions["x-kubernetes-init-containers"] = cloneMapSlice(service.InitContainers)
	}
	if len(service.ResourceClaims) > 0 {
		extensions["x-kubernetes-resource-claims"] = cloneMapSlice(service.ResourceClaims)
	}
	if len(service.EphemeralContainers) > 0 {
		extensions["x-kubernetes-ephemeral-containers"] = cloneMapSlice(service.EphemeralContainers)
	}
	if len(service.SchedulingGates) > 0 {
		extensions["x-kubernetes-scheduling-gates"] = cloneMapSlice(service.SchedulingGates)
	}
	if len(service.TopologySpreadConstraints) > 0 {
		extensions["x-kubernetes-topology-spread-constraints"] = cloneMapSlice(service.TopologySpreadConstraints)
	}
	return extensions
}

func writeNomadApplicationMeta(hcl *strings.Builder, app *Application) {
	if app == nil {
		return
	}
	meta := map[string]string{}
	if app.Name != "" {
		meta[nomadAppNameMetaKey] = app.Name
	}
	if app.Version != "" {
		meta[nomadAppVersionMetaKey] = app.Version
	}
	if len(app.Extensions) > 0 {
		extensions := jsonSerializableExtensionMap(app.Extensions)
		delete(extensions, nomadRawHCLExtension)
		delete(extensions, "x-platform")
		stripPortableCanonicalExtensionKeys(extensions)
		delete(extensions, composeKubernetesServicesExtensionKey)
		if len(extensions) > 0 {
			if raw, err := json.Marshal(extensions); err == nil {
				meta[nomadAppExtensionsMetaKey] = string(raw)
			}
		}
	}
	if models := applicationModelsForEmit(app); len(models) > 0 {
		if raw, err := json.Marshal(models); err == nil {
			meta[nomadAppModelsMetaKey] = string(raw)
		}
	}
	if len(app.IncludeEntries) > 0 {
		if raw, err := json.Marshal(app.IncludeEntries); err == nil {
			meta[nomadAppIncludesMetaKey] = string(raw)
		}
	}
	if networks := applicationNetworksForEmit(app); len(networks) > 0 {
		if raw, err := json.Marshal(networks); err == nil {
			meta[nomadAppNetworksMetaKey] = string(raw)
		}
	}
	if volumes := applicationVolumesForEmit(app); len(volumes) > 0 {
		if raw, err := json.Marshal(volumes); err == nil {
			meta[nomadAppVolumesMetaKey] = string(raw)
		}
	}
	if configs := applicationConfigsForEmit(app); len(configs) > 0 {
		if raw, err := json.Marshal(configs); err == nil {
			meta[nomadAppConfigsMetaKey] = string(raw)
		}
	}
	if secrets := applicationSecretsForEmit(app); len(secrets) > 0 {
		if raw, err := json.Marshal(secrets); err == nil {
			meta[nomadAppSecretsMetaKey] = string(raw)
		}
	}
	if routes := canonicalRoutesForApplication(app); len(routes) > 0 {
		if raw, err := json.Marshal(routes); err == nil {
			meta[nomadAppRoutesMetaKey] = string(raw)
		}
	}
	if policies := canonicalPoliciesForApplication(app); len(policies) > 0 {
		if raw, err := json.Marshal(policies); err == nil {
			meta[nomadAppPoliciesMetaKey] = string(raw)
		}
	}
	if resources := canonicalRawResourcesForBridge(app, PlatformNomad); len(resources) > 0 {
		if raw, err := json.Marshal(resources); err == nil {
			meta[nomadCanonicalRawResourcesMetaKey] = string(raw)
		}
	}
	if resources := canonicalKubernetesRawResourcesForApplication(app); len(resources) > 0 {
		if raw, err := json.Marshal(resources); err == nil {
			meta[nomadKubernetesRawResourcesMetaKey] = string(raw)
		}
	}
	if len(meta) == 0 {
		return
	}
	hcl.WriteString("  meta {\n")
	keys := make([]string, 0, len(meta))
	for key := range meta {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		hcl.WriteString(fmt.Sprintf("    %s = %q\n", key, meta[key]))
	}
	hcl.WriteString("  }\n\n")
}

func writeNomadPlacement(hcl *strings.Builder, service *Service) {
	if service == nil || service.Deploy == nil || service.Deploy.Placement == nil {
		return
	}
	for _, constraint := range service.Deploy.Placement.Constraints {
		attribute, operator, value := nomadConstraintFromPortable(constraint)
		if operator == "" {
			continue
		}
		hcl.WriteString("    constraint {\n")
		if attribute != "" {
			hcl.WriteString(fmt.Sprintf("      attribute = %q\n", attribute))
		}
		hcl.WriteString(fmt.Sprintf("      operator = %q\n", operator))
		if value != "" {
			hcl.WriteString(fmt.Sprintf("      value = %q\n", value))
		}
		hcl.WriteString("    }\n\n")
	}
	for _, preference := range service.Deploy.Placement.Preferences {
		if attribute := nomadSpreadAttributeFromPortable(preference); attribute != "" {
			hcl.WriteString("    spread {\n")
			hcl.WriteString(fmt.Sprintf("      attribute = %q\n", attribute))
			hcl.WriteString("    }\n\n")
			continue
		}
		attribute, operator, value := nomadAffinityFromPortable(preference)
		if attribute == "" || value == "" {
			continue
		}
		hcl.WriteString("    affinity {\n")
		hcl.WriteString(fmt.Sprintf("      attribute = %q\n", attribute))
		if operator != "" && operator != "=" {
			hcl.WriteString(fmt.Sprintf("      operator = %q\n", operator))
		}
		hcl.WriteString(fmt.Sprintf("      value = %q\n", value))
		hcl.WriteString("      weight = 50\n")
		hcl.WriteString("    }\n\n")
	}
	for _, spread := range service.Spreads {
		writeNomadSpreadBlock(hcl, spread)
	}
}

func writeNomadSchedulerBlocks(hcl *strings.Builder, service *Service) {
	if service == nil {
		return
	}
	writeNomadSchedulerBlock(hcl, service, nomadUpdateExtensionKey, "update", []string{
		"max_parallel",
		"health_check",
		"min_healthy_time",
		"healthy_deadline",
		"progress_deadline",
		"auto_revert",
		"auto_promote",
		"canary",
		"stagger",
	})
	writeNomadSchedulerBlock(hcl, service, nomadMigrateExtensionKey, "migrate", []string{
		"max_parallel",
		"health_check",
		"min_healthy_time",
		"healthy_deadline",
	})
	writeNomadSchedulerBlock(hcl, service, nomadRescheduleExtensionKey, "reschedule", []string{
		"attempts",
		"interval",
		"delay",
		"delay_function",
		"max_delay",
		"unlimited",
	})
}

func writeNomadSpreadBlock(hcl *strings.Builder, spread NomadSpreadSpec) {
	if hcl == nil || isEmptyNomadSpreadSpec(&spread) {
		return
	}
	hcl.WriteString("    spread {\n")
	if spread.Attribute != "" {
		hcl.WriteString(fmt.Sprintf("      attribute = %q\n", spread.Attribute))
	}
	if spread.Weight > 0 {
		hcl.WriteString(fmt.Sprintf("      weight = %d\n", spread.Weight))
	}
	if len(spread.Extensions) > 0 {
		keys := make([]string, 0, len(spread.Extensions))
		for key := range spread.Extensions {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			writeNomadGenericValue(hcl, "      ", key, spread.Extensions[key])
		}
	}
	for _, target := range spread.Targets {
		writeNomadSpreadTargetBlock(hcl, target)
	}
	hcl.WriteString("    }\n\n")
}

func writeNomadSpreadTargetBlock(hcl *strings.Builder, target NomadSpreadTarget) {
	if hcl == nil || (target.Value == "" && target.Percent == 0 && len(target.Extensions) == 0) {
		return
	}
	if target.Value != "" {
		hcl.WriteString(fmt.Sprintf("      target %q {\n", target.Value))
	} else {
		hcl.WriteString("      target {\n")
	}
	if target.Percent > 0 {
		hcl.WriteString(fmt.Sprintf("        percent = %d\n", target.Percent))
	}
	if len(target.Extensions) > 0 {
		keys := make([]string, 0, len(target.Extensions))
		for key := range target.Extensions {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			writeNomadGenericValue(hcl, "        ", key, target.Extensions[key])
		}
	}
	hcl.WriteString("      }\n")
}

func writeNomadSchedulerBlock(hcl *strings.Builder, service *Service, extensionKey, blockName string, orderedKeys []string) {
	if hcl == nil || service == nil || blockName == "" {
		return
	}
	block := nomadSchedulerBlockForService(service, extensionKey)
	if len(block) == 0 {
		return
	}
	hcl.WriteString(fmt.Sprintf("    %s {\n", blockName))
	emitNomadSchedulerBlockAttributes(hcl, block, orderedKeys, "      ")
	hcl.WriteString("    }\n\n")
}

func nomadSchedulerBlockForService(service *Service, extensionKey string) map[string]interface{} {
	if service == nil {
		return nil
	}
	switch extensionKey {
	case nomadMigrateExtensionKey:
		if service.Deploy != nil && !isEmptyMigratePolicy(service.Deploy.MigrateConfig) {
			return serializeMigratePolicy(service.Deploy.MigrateConfig)
		}
	case nomadRescheduleExtensionKey:
		if service.Deploy != nil && !isEmptyReschedulePolicy(service.Deploy.RescheduleConfig) {
			return serializeReschedulePolicy(service.Deploy.RescheduleConfig)
		}
	}
	if value, ok := firstExtensionValue(service.Extensions, extensionKey, composeServiceExtensionKey(extensionKey)); ok {
		if block, ok := asMap(value); ok && len(block) > 0 {
			return cloneMap(block)
		}
	}
	if service.Deploy == nil {
		return nil
	}
	switch extensionKey {
	case nomadUpdateExtensionKey:
		if !isEmptyUpdatePolicy(service.Deploy.UpdateConfig) {
			return nomadUpdatePolicyForHCL(service.Deploy.UpdateConfig)
		}
	}
	return nil
}

func nomadUpdatePolicyForHCL(policy *UpdatePolicy) map[string]interface{} {
	update := serializeNomadUpdatePolicy(policy)
	if len(update) == 0 {
		return nil
	}
	delete(update, "kubernetes-deployment-strategy")
	delete(update, "x-kubernetes-deployment-strategy")
	delete(update, "kubernetes.workload.updateStrategy")
	delete(update, "x-kubernetes-workload-updateStrategy")
	if nested, ok := asMap(update["extensions"]); ok && len(nested) > 0 {
		delete(update, "extensions")
		for key, value := range nested {
			update[key] = deepCopyValue(value)
		}
	}
	return update
}

func parseNomadMigratePolicyMap(mapped map[string]interface{}) *MigratePolicy {
	if len(mapped) == 0 {
		return nil
	}
	policy := &MigratePolicy{Extensions: map[string]interface{}{}}
	for key, value := range mapped {
		switch key {
		case "max_parallel":
			policy.MaxParallel = toInt(value)
		case "health_check":
			policy.HealthCheck = toString(value)
		case "min_healthy_time":
			policy.MinHealthyTime = toString(value)
		case "healthy_deadline":
			policy.HealthyDeadline = toString(value)
		default:
			policy.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(policy.Extensions) == 0 {
		policy.Extensions = nil
	}
	if isEmptyMigratePolicy(policy) {
		return nil
	}
	return policy
}

func parseNomadReschedulePolicyMap(mapped map[string]interface{}) *ReschedulePolicy {
	if len(mapped) == 0 {
		return nil
	}
	policy := &ReschedulePolicy{Extensions: map[string]interface{}{}}
	for key, value := range mapped {
		switch key {
		case "attempts":
			policy.Attempts = toInt(value)
		case "interval":
			policy.Interval = toString(value)
		case "delay":
			policy.Delay = toString(value)
		case "delay_function":
			policy.DelayFunction = toString(value)
		case "max_delay":
			policy.MaxDelay = toString(value)
		case "unlimited":
			policy.Unlimited = toBool(value)
		default:
			policy.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(policy.Extensions) == 0 {
		policy.Extensions = nil
	}
	if isEmptyReschedulePolicy(policy) {
		return nil
	}
	return policy
}

func emitNomadSchedulerBlockAttributes(hcl *strings.Builder, attrs map[string]interface{}, orderedKeys []string, indent string) {
	if hcl == nil || len(attrs) == 0 {
		return
	}
	emitted := map[string]bool{}
	writeAttr := func(key string) {
		if emitted[key] {
			return
		}
		value, ok := attrs[key]
		if !ok {
			return
		}
		emitted[key] = true
		switch typed := value.(type) {
		case string:
			if typed == "" {
				return
			}
			hcl.WriteString(fmt.Sprintf("%s%s = %q\n", indent, key, typed))
		case bool:
			hcl.WriteString(fmt.Sprintf("%s%s = %t\n", indent, key, typed))
		case int:
			hcl.WriteString(fmt.Sprintf("%s%s = %d\n", indent, key, typed))
		case int32:
			hcl.WriteString(fmt.Sprintf("%s%s = %d\n", indent, key, typed))
		case int64:
			hcl.WriteString(fmt.Sprintf("%s%s = %d\n", indent, key, typed))
		case float64:
			hcl.WriteString(fmt.Sprintf("%s%s = %v\n", indent, key, typed))
		default:
			if raw, err := json.Marshal(typed); err == nil {
				hcl.WriteString(fmt.Sprintf("%s%s = %s\n", indent, key, string(raw)))
			}
		}
	}
	for _, key := range orderedKeys {
		writeAttr(key)
	}
	if len(emitted) == len(attrs) {
		return
	}
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		if !emitted[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		writeAttr(key)
	}
}

func nomadSpreadAttributeFromPortable(preference string) string {
	preference = strings.TrimSpace(strings.TrimPrefix(preference, "prefer:"))
	if !strings.HasPrefix(preference, "spread=") {
		return ""
	}
	return nomadAttributeFromPortable(strings.TrimPrefix(preference, "spread="))
}

func applyNomadOpaqueKubernetesExtensions(service *Service, extensions map[string]interface{}) {
	if service == nil || len(extensions) == 0 {
		return
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-startup-probe", "kubernetes.startupProbe"); ok {
		if probe, ok := asMap(value); ok {
			service.StartupProbe = parseKubernetesProbe(probe)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-readiness-gates", "kubernetes.readinessGates"); ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); ok && len(gates) > 0 {
			service.ReadinessGates = gates
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-allowPrivilegeEscalation", "kubernetes.allowPrivilegeEscalation"); ok {
		service.AllowPrivilegeEscalation = boolPtrFromInterface(value)
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-procMount", "kubernetes.procMount"); ok {
		if text := toString(value); text != "" {
			service.ProcMount = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-init-containers", "kubernetes.initContainers"); ok {
		if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
			service.InitContainers = containers
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-resource-claims", "kubernetes.resourceClaims"); ok {
		if claims, ok := kubernetesMapSliceFromExtension(value); ok && len(claims) > 0 {
			service.ResourceClaims = claims
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-ephemeral-containers", "kubernetes.ephemeralContainers"); ok {
		if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
			service.EphemeralContainers = containers
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-scheduling-gates", "kubernetes.schedulingGates"); ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); ok && len(gates) > 0 {
			service.SchedulingGates = gates
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-seccomp-profile", "kubernetes.seccompProfile"); ok {
		if profile, ok := asMap(value); ok {
			service.SeccompProfile = parseKubernetesSeccompProfile(profile)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-affinity", "kubernetes.affinity"); ok {
		if affinity, ok := asMap(value); ok {
			service.Affinity = copyStringInterfaceMap(affinity)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-topology-spread-constraints", "kubernetes.topologySpreadConstraints"); ok {
		if constraints, ok := kubernetesTopologySpreadConstraintsFromExtension(value); ok && len(constraints) > 0 {
			service.TopologySpreadConstraints = constraints
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-node-selector", "x-kubernetes-nodeSelector", "kubernetes.nodeSelector"); ok {
		if selector, err := toStringMap(value); err == nil && len(selector) > 0 && len(service.NodeSelector) == 0 {
			service.NodeSelector = copyStringMap(selector)
		}
	}
	if secrets, ok := firstExtensionValue(extensions, "x-kubernetes-imagePullSecrets", "kubernetes.imagePullSecrets"); ok {
		if names, err := toStringSlice(secrets); err == nil && len(names) > 0 {
			service.ImagePullSecrets = names
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-imagePullPolicy", "kubernetes.imagePullPolicy"); ok {
		if policy := toString(value); policy != "" {
			service.ImagePullPolicy = policy
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-hostAliases", "kubernetes.hostAliases"); ok {
		if aliases := kubernetesHostAliasesFromExtension(value); len(aliases) > 0 {
			service.HostAliases = aliases
			for _, alias := range aliases {
				for _, hostname := range alias.Hostnames {
					appendUniqueString(&service.ExtraHosts, hostname+"="+alias.IP)
				}
			}
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-dnsPolicy", "kubernetes.dnsPolicy"); ok {
		if policy := toString(value); policy != "" {
			service.DNSPolicy = policy
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-schedulerName", "kubernetes.schedulerName"); ok {
		if scheduler := toString(value); scheduler != "" {
			service.SchedulerName = scheduler
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-terminationMessagePath", "kubernetes.terminationMessagePath"); ok {
		if path := toString(value); path != "" {
			service.TerminationMessagePath = path
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-terminationMessagePolicy", "kubernetes.terminationMessagePolicy"); ok {
		if policy := toString(value); policy != "" {
			service.TerminationMessagePolicy = policy
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-hostNetwork", "kubernetes.hostNetwork"); ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostNetwork = flag
			service.HostNetworkSet = true
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-hostPID", "kubernetes.hostPID"); ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostPID = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-hostIPC", "kubernetes.hostIPC"); ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostIPC = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-pidMode", "kubernetes.pidMode"); ok {
		if text := toString(value); text != "" {
			service.PIDMode = text
			if strings.EqualFold(text, "host") {
				service.HostPID = boolPtr(true)
			}
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-ipcMode", "kubernetes.ipcMode"); ok {
		if text := toString(value); text != "" {
			service.IPCMode = text
			if strings.EqualFold(text, "host") {
				service.HostIPC = boolPtr(true)
			}
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-priorityClassName", "kubernetes.priorityClassName"); ok {
		if text := toString(value); text != "" {
			service.PriorityClassName = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-runtimeClassName", "kubernetes.runtimeClassName"); ok {
		if text := toString(value); text != "" {
			service.RuntimeClassName = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-nodeName", "kubernetes.nodeName"); ok {
		if text := toString(value); text != "" {
			service.NodeName = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-subdomain", "kubernetes.subdomain"); ok {
		if text := toString(value); text != "" {
			service.Subdomain = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-os", "kubernetes.os"); ok {
		if text := toString(value); text != "" {
			service.OSName = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-setHostnameAsFQDN", "kubernetes.setHostnameAsFQDN"); ok {
		if text := toString(value); text != "" {
			flag := strings.EqualFold(text, "true")
			service.SetHostnameAsFQDN = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-hostUsers", "kubernetes.hostUsers"); ok {
		if text := toString(value); text != "" {
			flag := strings.EqualFold(text, "true")
			service.HostUsers = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-group", "kubernetes.group"); ok {
		if text := toString(value); text != "" {
			service.Group = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-shareProcessNamespace", "kubernetes.shareProcessNamespace"); ok {
		if text := toString(value); text != "" {
			flag := strings.EqualFold(text, "true")
			service.ShareProcessNamespace = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-enableServiceLinks", "kubernetes.enableServiceLinks"); ok {
		if text := toString(value); text != "" {
			flag := strings.EqualFold(text, "true")
			service.EnableServiceLinks = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-serviceAccountName", "kubernetes.serviceAccountName"); ok {
		if text := toString(value); text != "" {
			service.ServiceAccountName = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-automountServiceAccountToken", "kubernetes.automountServiceAccountToken"); ok {
		if text := toString(value); text != "" {
			flag := strings.EqualFold(text, "true")
			service.AutomountServiceAccountToken = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-fsGroup", "kubernetes.fsGroup"); ok {
		fsGroup := int64(toInt(value))
		if fsGroup > 0 {
			service.FSGroup = &fsGroup
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-seLinuxOptions", "kubernetes.seLinuxOptions"); ok {
		if options, ok := asMap(value); ok {
			service.SELinuxOptions = parseKubernetesSELinuxOptions(options)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-windowsOptions", "kubernetes.windowsOptions"); ok {
		if options, ok := asMap(value); ok {
			service.WindowsOptions = parseKubernetesWindowsSecurityContextOptions(options)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-fsGroupChangePolicy", "kubernetes.fsGroupChangePolicy"); ok {
		if text := toString(value); text != "" {
			service.FSGroupChangePolicy = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-runAsNonRoot", "kubernetes.runAsNonRoot"); ok {
		if text := toString(value); text != "" {
			flag := strings.EqualFold(text, "true")
			service.RunAsNonRoot = boolPtr(flag)
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-supplementalGroups", "kubernetes.supplementalGroups"); ok {
		if groups, err := toInt64Slice(value); err == nil && len(groups) > 0 {
			service.SupplementalGroups = groups
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-supplementalGroupsPolicy", "kubernetes.supplementalGroupsPolicy"); ok {
		if text := toString(value); text != "" {
			service.SupplementalGroupsPolicy = text
		}
	}
	if value, ok := firstExtensionValue(extensions, "kubernetes.groupAdd", "x-kubernetes-groupAdd"); ok {
		if groupAdd, err := toStringSlice(value); err == nil && len(groupAdd) > 0 {
			service.GroupAdd = groupAdd
		}
	}
	if value, ok := firstExtensionValue(extensions, "kubernetes.runtime", "x-kubernetes-runtime"); ok {
		if runtime := toString(value); runtime != "" {
			service.Runtime = runtime
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-activeDeadlineSeconds", "kubernetes.activeDeadlineSeconds"); ok {
		seconds := int64(toInt(value))
		if seconds > 0 {
			service.ActiveDeadlineSeconds = &seconds
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-tolerations", "kubernetes.tolerations"); ok {
		if tolerations := parseKubernetesTolerationsExtension(value); len(tolerations) > 0 {
			service.Tolerations = tolerations
		}
	}
	if value, ok := firstExtensionValue(extensions, "x-kubernetes-restartPolicy", "kubernetes.restartPolicy"); ok {
		if restartPolicy := toString(value); restartPolicy != "" {
			service.PodRestartPolicy = restartPolicy
		}
	}
}

func firstExtensionValue(extensions map[string]interface{}, keys ...string) (interface{}, bool) {
	for _, key := range keys {
		if value, ok := extensions[key]; ok {
			return value, true
		}
	}
	return nil, false
}

func jsonSerializableExtensionMap(extensions map[string]interface{}) map[string]interface{} {
	if len(extensions) == 0 {
		return nil
	}
	result := map[string]interface{}{}
	for key, value := range extensions {
		encoded, err := json.Marshal(value)
		if err != nil {
			continue
		}
		var decoded interface{}
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			continue
		}
		result[key] = decoded
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func nomadConstraintFromPortable(constraint string) (string, string, string) {
	constraint = strings.TrimSpace(constraint)
	if constraint == "" {
		return "", "", ""
	}
	if strings.HasPrefix(constraint, "nomad.distinct_hosts") {
		return "", "distinct_hosts", portableConstraintValue(constraint, "true")
	}
	if strings.HasPrefix(constraint, "nomad.distinct_property") {
		return "", "distinct_property", portableConstraintValue(constraint, "")
	}
	left, operator, right := splitPortableConstraint(constraint)
	if operator == "" {
		return "${meta.bolabaden.constraint}", "=", constraint
	}
	attribute := nomadAttributeFromPortable(left)
	if attribute == "" {
		return "${meta.bolabaden.constraint}", "=", constraint
	}
	return attribute, nomadOperatorFromPortable(operator), right
}

func nomadAffinityFromPortable(preference string) (string, string, string) {
	preference = strings.TrimSpace(strings.TrimPrefix(preference, "prefer:"))
	if strings.HasPrefix(preference, "spread=") {
		attribute := nomadAttributeFromPortable(strings.TrimPrefix(preference, "spread="))
		if attribute == "" {
			return "", "", ""
		}
		return attribute, "=", "true"
	}
	left, operator, right := splitPortableConstraint(preference)
	if operator == "" {
		return "", "", ""
	}
	return nomadAttributeFromPortable(left), nomadOperatorFromPortable(operator), right
}

func splitPortableConstraint(constraint string) (string, string, string) {
	for _, operator := range []string{"==", "!=", ">=", "<=", ">", "<", "=", "is_not_set", "is_set"} {
		if strings.Contains(constraint, " "+operator+" ") {
			parts := strings.SplitN(constraint, " "+operator+" ", 2)
			return strings.TrimSpace(parts[0]), operator, strings.TrimSpace(parts[1])
		}
		if strings.HasSuffix(constraint, " "+operator) {
			return strings.TrimSpace(strings.TrimSuffix(constraint, " "+operator)), operator, ""
		}
	}
	return "", "", ""
}

func nomadAttributeFromPortable(attribute string) string {
	attribute = strings.TrimSpace(attribute)
	switch {
	case strings.HasPrefix(attribute, "node.labels."):
		return "${meta." + strings.TrimPrefix(attribute, "node.labels.") + "}"
	case attribute == "node.hostname":
		return "${node.unique.name}"
	case attribute == "node.datacenter":
		return "${node.datacenter}"
	case strings.HasPrefix(attribute, "nomad.attr."):
		return "${" + strings.TrimPrefix(attribute, "nomad.") + "}"
	case strings.HasPrefix(attribute, "${"):
		return attribute
	default:
		return ""
	}
}

func nomadOperatorFromPortable(operator string) string {
	switch strings.TrimSpace(operator) {
	case "==":
		return "="
	default:
		return strings.TrimSpace(operator)
	}
}

func portableConstraintValue(constraint, fallback string) string {
	_, _, value := splitPortableConstraint(constraint)
	if value == "" {
		return fallback
	}
	return value
}

func writeNomadDependencyTasks(hcl *strings.Builder, service *Service) {
	for _, dependency := range serviceDependencies(service) {
		if dependency.Name == "" {
			continue
		}
		taskName := "wait-for-" + sanitizeKubernetesName(dependency.Name)
		hcl.WriteString(fmt.Sprintf("    task \"%s\" {\n", taskName))
		hcl.WriteString("      lifecycle {\n")
		hcl.WriteString("        hook = \"prestart\"\n")
		hcl.WriteString("        sidecar = false\n")
		hcl.WriteString("      }\n\n")
		hcl.WriteString("      driver = \"exec\"\n\n")
		hcl.WriteString("      config {\n")
		hcl.WriteString("        command = \"sh\"\n")
		hcl.WriteString(fmt.Sprintf("        args = [\"-c\", %q]\n", nomadDependencyWaitCommand(dependency)))
		hcl.WriteString("      }\n\n")
		hcl.WriteString("      env {\n")
		hcl.WriteString(fmt.Sprintf("        BOLABADEN_DEPENDENCY_NAME = %q\n", dependency.Name))
		if dependency.Condition != "" {
			hcl.WriteString(fmt.Sprintf("        BOLABADEN_DEPENDENCY_CONDITION = %q\n", dependency.Condition))
		}
		if dependency.Restart {
			hcl.WriteString("        BOLABADEN_DEPENDENCY_RESTART = \"true\"\n")
		}
		if dependency.Required != nil {
			hcl.WriteString(fmt.Sprintf("        BOLABADEN_DEPENDENCY_REQUIRED = %q\n", fmt.Sprintf("%t", *dependency.Required)))
		}
		hcl.WriteString("      }\n")
		hcl.WriteString("    }\n\n")
	}
}

func nomadDependencyWaitCommand(dependency DependencySpec) string {
	name := dependency.Name
	if name == "" {
		name = "dependency"
	}
	return fmt.Sprintf("until getent hosts %s.service.consul >/dev/null 2>&1 || getent hosts %s >/dev/null 2>&1; do echo waiting for %s; sleep 2; done", name, name, name)
}

func writeNomadRuntimeConfig(hcl *strings.Builder, service *Service) {
	if len(service.Entrypoint) > 0 {
		hcl.WriteString(fmt.Sprintf("        command = %q\n", service.Entrypoint[0]))
	}
	args := append([]string{}, service.Command...)
	if len(service.Entrypoint) > 1 {
		args = append(append([]string{}, service.Entrypoint[1:]...), args...)
	}
	if len(args) > 0 {
		hcl.WriteString(fmt.Sprintf("        args = [%s]\n", quotedStringList(args)))
	}
	if service.WorkingDir != "" {
		hcl.WriteString(fmt.Sprintf("        work_dir = %q\n", service.WorkingDir))
	}
	if service.Hostname != "" {
		hcl.WriteString(fmt.Sprintf("        hostname = %q\n", service.Hostname))
	}
	if service.PrivilegedSet {
		hcl.WriteString(fmt.Sprintf("        privileged = %t\n", service.Privileged))
	} else if service.Privileged {
		hcl.WriteString("        privileged = true\n")
	}
	if service.ReadOnlyRootFSSet {
		hcl.WriteString(fmt.Sprintf("        readonly_rootfs = %t\n", service.ReadOnlyRootFS))
	} else if service.ReadOnlyRootFS {
		hcl.WriteString("        readonly_rootfs = true\n")
	}
	if service.Runtime != "" {
		hcl.WriteString(fmt.Sprintf("        runtime = %q\n", service.Runtime))
	}
	if service.PIDMode != "" {
		hcl.WriteString(fmt.Sprintf("        pid_mode = %q\n", service.PIDMode))
	}
	if service.IPCMode != "" {
		hcl.WriteString(fmt.Sprintf("        ipc_mode = %q\n", service.IPCMode))
	}
	if service.pidsLimitSet || service.PidsLimit > 0 {
		hcl.WriteString(fmt.Sprintf("        pids_limit = %d\n", service.PidsLimit))
	}
	if service.shmSizeSet || service.ShmSize > 0 {
		hcl.WriteString(fmt.Sprintf("        shm_size = %d\n", service.ShmSize))
	}
	if len(service.CapAdd) > 0 {
		hcl.WriteString(fmt.Sprintf("        cap_add = [%s]\n", quotedStringList(service.CapAdd)))
	}
	if len(service.GroupAdd) > 0 {
		hcl.WriteString(fmt.Sprintf("        group_add = [%s]\n", quotedStringList(service.GroupAdd)))
	}
	if len(service.Sysctls) > 0 {
		hcl.WriteString("        sysctl = {\n")
		for _, key := range sortedMapKeys(service.Sysctls) {
			hcl.WriteString(fmt.Sprintf("          %q = %q\n", key, service.Sysctls[key]))
		}
		hcl.WriteString("        }\n")
	}
	if len(service.CapDrop) > 0 {
		hcl.WriteString(fmt.Sprintf("        cap_drop = [%s]\n", quotedStringList(service.CapDrop)))
	}
	if len(service.SecurityOpt) > 0 {
		hcl.WriteString(fmt.Sprintf("        security_opt = [%s]\n", quotedStringList(service.SecurityOpt)))
	}
	if len(service.DNS) > 0 {
		hcl.WriteString(fmt.Sprintf("        dns_servers = [%s]\n", quotedStringList(service.DNS)))
	}
	if len(service.DNSSearch) > 0 {
		hcl.WriteString(fmt.Sprintf("        dns_search_domains = [%s]\n", quotedStringList(service.DNSSearch)))
	}
	if len(service.DNSOptions) > 0 {
		hcl.WriteString(fmt.Sprintf("        dns_options = [%s]\n", quotedStringList(service.DNSOptions)))
	}
	if len(service.ExtraHosts) > 0 {
		hcl.WriteString(fmt.Sprintf("        extra_hosts = [%s]\n", quotedStringList(service.ExtraHosts)))
	}
}

func writeNomadVolumeMounts(hcl *strings.Builder, service *Service) {
	for _, volume := range service.Volumes {
		if volume.Target == "" {
			continue
		}
		mountType := nomadMountType(volume)
		hcl.WriteString("        mount {\n")
		hcl.WriteString(fmt.Sprintf("          type = %q\n", mountType))
		if volume.Source != "" {
			hcl.WriteString(fmt.Sprintf("          source = %q\n", volume.Source))
		}
		hcl.WriteString(fmt.Sprintf("          target = %q\n", volume.Target))
		if volume.ReadOnly {
			hcl.WriteString("          readonly = true\n")
		}
		switch mountType {
		case "bind":
			if volume.Propagation != "" {
				hcl.WriteString("          bind_options {\n")
				hcl.WriteString(fmt.Sprintf("            propagation = %q\n", volume.Propagation))
				hcl.WriteString("          }\n")
			}
		case "volume":
			if volume.NoCopy || volume.Options["driver"] != "" {
				hcl.WriteString("          volume_options {\n")
				if volume.NoCopy {
					hcl.WriteString("            no_copy = true\n")
				}
				if driver := volume.Options["driver"]; driver != "" {
					hcl.WriteString("            driver_config {\n")
					hcl.WriteString(fmt.Sprintf("              name = %q\n", driver))
					hcl.WriteString("            }\n")
				}
				hcl.WriteString("          }\n")
			}
		case "tmpfs":
			if size := nomadTmpfsSize(volume.Options["size"]); size != "" {
				hcl.WriteString("          tmpfs_options {\n")
				hcl.WriteString(fmt.Sprintf("            size = %s\n", size))
				hcl.WriteString("          }\n")
			}
		}
		hcl.WriteString("        }\n")
	}
}

func writeNomadTaskMeta(hcl *strings.Builder, service *Service) {
	if service == nil || len(service.Volumes) == 0 {
		return
	}
	tmpfs := map[string]string{}
	for _, volume := range service.Volumes {
		if volume.Type != "tmpfs" || volume.Target == "" {
			continue
		}
		if raw := rawTmpfsSize(volume); raw != "" {
			tmpfs[volume.Target] = raw
		}
	}
	if len(tmpfs) == 0 {
		return
	}
	raw, err := json.Marshal(tmpfs)
	if err != nil {
		return
	}
	hcl.WriteString("      meta {\n")
	hcl.WriteString(fmt.Sprintf("        %s = %q\n", nomadTmpfsMetaKey, string(raw)))
	hcl.WriteString("      }\n\n")
}

func nomadTmpfsSize(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if size := parseInt(value); size > 0 {
		return fmt.Sprintf("%d", size)
	}
	quantity, err := resource.ParseQuantity(value)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", quantity.Value())
}

func rawTmpfsSize(volume VolumeMount) string {
	if volume.TmpfsExtensions != nil {
		if raw := toString(volume.TmpfsExtensions["size"]); raw != "" {
			return raw
		}
	}
	if raw := volume.Options["size"]; raw != "" && parseInt(raw) <= 0 {
		return raw
	}
	return ""
}

func nomadMountType(volume VolumeMount) string {
	if volume.Type != "" {
		return volume.Type
	}
	if volume.Source == "" {
		return "tmpfs"
	}
	if strings.HasPrefix(volume.Source, "/") || strings.HasPrefix(volume.Source, ".") || strings.HasPrefix(volume.Source, "~") {
		return "bind"
	}
	return "volume"
}

func writeNomadTemplates(hcl *strings.Builder, service *Service) {
	for _, ref := range service.Configs {
		if ref.Source == "" {
			continue
		}
		destination := ref.Target
		if destination == "" {
			destination = "local/configs/" + sanitizeKubernetesName(ref.Source)
		}
		writeNomadTemplate(hcl, nomadFileRefTemplateData("config", ref), destination, false)
	}
	for _, ref := range service.Secrets {
		if ref.Source == "" {
			continue
		}
		destination := ref.Target
		if destination == "" {
			destination = "${NOMAD_SECRETS_DIR}/" + sanitizeKubernetesName(ref.Source)
		}
		writeNomadTemplate(hcl, nomadFileRefTemplateData("secret", ref), destination, false)
	}
	for _, source := range service.EnvSources {
		if source.Name == "" || source.Source == "" {
			continue
		}
		destination := "secrets/env-" + sanitizeKubernetesName(source.Name) + ".env"
		writeNomadTemplate(hcl, nomadEnvSourceTemplateData(source), destination, true)
	}
	for _, source := range service.EnvFrom {
		if source.Source == "" {
			continue
		}
		destination := "secrets/env-from-" + sanitizeKubernetesName(source.Source) + ".env"
		writeNomadTemplate(hcl, nomadEnvFromTemplateData(source), destination, true)
	}
}

func quotedStringList(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("%q", value))
	}
	return strings.Join(quoted, ", ")
}

func writeNomadTemplate(hcl *strings.Builder, data, destination string, env bool) {
	hcl.WriteString("      template {\n")
	hcl.WriteString("        data = <<EOH\n")
	hcl.WriteString(data)
	if !strings.HasSuffix(data, "\n") {
		hcl.WriteString("\n")
	}
	hcl.WriteString("EOH\n")
	hcl.WriteString(fmt.Sprintf("        destination = %q\n", destination))
	if env {
		hcl.WriteString("        env = true\n")
	}
	hcl.WriteString("        change_mode = \"restart\"\n")
	hcl.WriteString("      }\n\n")
}

func nomadFileRefTemplateData(sourceType string, ref FileRef) string {
	optional := ""
	if ref.Optional != nil {
		optional = fmt.Sprintf(" optional=%t", *ref.Optional)
	}
	return fmt.Sprintf("{{/* bolabaden:file_ref type=%s source=%s key=%s mode=%s read_only=%t%s */}}\n%s\n",
		sourceType,
		nomadMarkerValue(ref.Source),
		nomadMarkerValue(ref.Key),
		nomadMarkerValue(ref.Mode),
		ref.ReadOnly,
		optional,
		nomadTemplateLookup(sourceType, ref.Source, ref.Key),
	)
}

func nomadEnvSourceTemplateData(source EnvSource) string {
	key := source.Key
	if key == "" {
		key = source.Name
	}
	return fmt.Sprintf("{{/* bolabaden:env_source name=%s type=%s source=%s key=%s optional=%t extensions_b64=%s */}}\n%s=%s\n",
		nomadMarkerValue(source.Name),
		nomadMarkerValue(source.SourceType),
		nomadMarkerValue(source.Source),
		nomadMarkerValue(source.Key),
		source.Optional,
		nomadMarkerValue(nomadExtensionsBlob(source.Extensions)),
		source.Name,
		nomadTemplateLookup(source.SourceType, source.Source, key),
	)
}

func nomadEnvFromTemplateData(source EnvFromSource) string {
	return fmt.Sprintf("{{/* bolabaden:env_from type=%s source=%s prefix=%s optional=%t extensions_b64=%s */}}\n{{ range $key, $value := %s }}%s{{ $key }}={{ $value }}\n{{ end }}\n",
		nomadMarkerValue(source.SourceType),
		nomadMarkerValue(source.Source),
		nomadMarkerValue(source.Prefix),
		source.Optional,
		nomadMarkerValue(nomadExtensionsBlob(source.Extensions)),
		nomadTemplateMapLookup(source.SourceType, source.Source),
		source.Prefix,
	)
}

func nomadExtensionsBlob(extensions map[string]interface{}) string {
	if len(extensions) == 0 {
		return ""
	}
	data, err := json.Marshal(extensions)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

func nomadExtensionsFromMarker(value string) map[string]interface{} {
	if value == "" {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(data) == 0 {
		return nil
	}
	extensions := map[string]interface{}{}
	if err := json.Unmarshal(data, &extensions); err != nil || len(extensions) == 0 {
		return nil
	}
	return extensions
}

func nomadTemplateLookup(sourceType, source, key string) string {
	switch sourceType {
	case "secret":
		if key == "" {
			return fmt.Sprintf("{{ with secret %q }}{{ .Data.data.value }}{{ end }}", source)
		}
		return fmt.Sprintf("{{ with secret %q }}{{ index .Data.data %q }}{{ end }}", source, key)
	default:
		if key == "" {
			return fmt.Sprintf("{{ key %q }}", source)
		}
		return fmt.Sprintf("{{ key %q }}", strings.TrimSuffix(source, "/")+"/"+key)
	}
}

func nomadTemplateMapLookup(sourceType, source string) string {
	if sourceType == "secret" {
		return fmt.Sprintf("(secret %q).Data.data", source)
	}
	return fmt.Sprintf("tree %q", strings.TrimSuffix(source, "/")+"/")
}

func nomadMarkerValue(value string) string {
	return strings.NewReplacer(" ", "_", "\t", "_", "\n", "_", "\r", "_").Replace(value)
}

func writeNomadService(hcl *strings.Builder, serviceName string, service *Service) {
	hasConnect := service.Connect != nil
	if !hasConnect {
		if connectValue, ok := firstExtensionValue(service.Extensions, nomadConnectExtensionKey, "x-nomad-connect"); ok {
			if connect, ok := asMap(connectValue); ok && len(connect) > 0 {
				hasConnect = true
			}
		}
	}
	if service.HealthCheck == nil && len(service.Labels) == 0 && !hasConnect {
		return
	}
	name := service.Labels["nomad.service.name"]
	if name == "" {
		name = serviceName
	}
	port := service.Labels["nomad.service.port"]
	if port == "" {
		port = firstNomadPortLabel(service)
	}
	hcl.WriteString("\n      service {\n")
	hcl.WriteString(fmt.Sprintf("        name = %q\n", name))
	if port != "" {
		hcl.WriteString(fmt.Sprintf("        port = %q\n", port))
	}
	if service.HealthCheck != nil {
		writeNomadCheck(hcl, service, port)
	}
	writeNomadConnect(hcl, service)
	hcl.WriteString("      }\n")
}

func writeNomadConnect(hcl *strings.Builder, service *Service) {
	connectValue, ok := firstExtensionValue(service.Extensions, nomadConnectExtensionKey, "x-nomad-connect")
	if !ok {
		connectValue = nomadConnectSpecToMap(service.Connect)
	}
	connect, ok := asMap(connectValue)
	if !ok || len(connect) == 0 {
		return
	}
	hcl.WriteString("        connect {\n")
	if native, ok := connect["native"].(bool); ok {
		hcl.WriteString(fmt.Sprintf("          native = %t\n", native))
	}
	if sidecar, ok := asMap(connect["sidecar_service"]); ok && len(sidecar) > 0 {
		hcl.WriteString("          sidecar_service {\n")
		if tags, err := toStringSlice(sidecar["tags"]); err == nil && len(tags) > 0 {
			hcl.WriteString(fmt.Sprintf("            tags = [%s]\n", quotedStringList(tags)))
		}
		if check, ok := asMap(sidecar["check"]); ok && len(check) > 0 {
			writeNomadGenericBlock(hcl, "            ", "check", check)
		}
		if proxy, ok := asMap(sidecar["proxy"]); ok && len(proxy) > 0 {
			hcl.WriteString("            proxy {\n")
			if upstreams, ok := proxy["upstreams"].([]interface{}); ok {
				for _, item := range upstreams {
					if upstream, ok := asMap(item); ok {
						writeNomadConnectUpstream(hcl, upstream)
					}
				}
			}
			if config, ok := asMap(proxy["config"]); ok && len(config) > 0 {
				writeNomadGenericBlock(hcl, "              ", "config", config)
			}
			if extensions, ok := asMap(proxy["extensions"]); ok && len(extensions) > 0 {
				writeNomadGenericBlock(hcl, "              ", "extensions", extensions)
			}
			hcl.WriteString("            }\n")
		}
		if extensions, ok := asMap(sidecar["extensions"]); ok && len(extensions) > 0 {
			writeNomadGenericBlock(hcl, "            ", "extensions", extensions)
		}
		hcl.WriteString("          }\n")
	}
	if gateway, ok := asMap(connect["gateway"]); ok && len(gateway) > 0 {
		writeNomadGenericBlock(hcl, "          ", "gateway", gateway)
	}
	if extensions, ok := asMap(connect["extensions"]); ok && len(extensions) > 0 {
		writeNomadGenericBlock(hcl, "          ", "extensions", extensions)
	}
	hcl.WriteString("        }\n")
}

func writeNomadGenericBlock(hcl *strings.Builder, indent, name string, data map[string]interface{}) {
	if hcl == nil || len(data) == 0 {
		return
	}
	hcl.WriteString(fmt.Sprintf("%s%s {\n", indent, name))
	nextIndent := indent + "  "
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		writeNomadGenericValue(hcl, nextIndent, key, data[key])
	}
	hcl.WriteString(fmt.Sprintf("%s}\n", indent))
}

func writeNomadGenericValue(hcl *strings.Builder, indent, key string, value interface{}) {
	switch typed := value.(type) {
	case string:
		hcl.WriteString(fmt.Sprintf("%s%s = %q\n", indent, key, typed))
	case bool:
		hcl.WriteString(fmt.Sprintf("%s%s = %t\n", indent, key, typed))
	case int:
		hcl.WriteString(fmt.Sprintf("%s%s = %d\n", indent, key, typed))
	case int32:
		hcl.WriteString(fmt.Sprintf("%s%s = %d\n", indent, key, typed))
	case int64:
		hcl.WriteString(fmt.Sprintf("%s%s = %d\n", indent, key, typed))
	case float64:
		hcl.WriteString(fmt.Sprintf("%s%s = %v\n", indent, key, typed))
	case []string:
		hcl.WriteString(fmt.Sprintf("%s%s = [%s]\n", indent, key, quotedStringList(typed)))
	case []interface{}:
		if len(typed) == 0 {
			return
		}
		if blocks, ok := nomadGenericBlockSlice(typed); ok {
			for _, block := range blocks {
				writeNomadGenericBlock(hcl, indent, key, block)
			}
			return
		}
		if raw, err := json.Marshal(typed); err == nil {
			hcl.WriteString(fmt.Sprintf("%s%s = %s\n", indent, key, string(raw)))
		}
	case map[string]interface{}:
		writeNomadGenericBlock(hcl, indent, key, typed)
	default:
		if raw, err := json.Marshal(typed); err == nil {
			hcl.WriteString(fmt.Sprintf("%s%s = %s\n", indent, key, string(raw)))
		}
	}
}

func nomadGenericBlockSlice(values []interface{}) ([]map[string]interface{}, bool) {
	if len(values) == 0 {
		return nil, false
	}
	result := make([]map[string]interface{}, 0, len(values))
	for _, value := range values {
		mapped, ok := asMap(value)
		if !ok {
			return nil, false
		}
		result = append(result, mapped)
	}
	return result, true
}

func writeNomadConnectUpstream(hcl *strings.Builder, upstream map[string]interface{}) {
	destination := toString(upstream["destination_name"])
	localPort := parseInt(toString(upstream["local_bind_port"]))
	if destination == "" || localPort == 0 {
		return
	}
	hcl.WriteString("              upstreams {\n")
	hcl.WriteString(fmt.Sprintf("                destination_name = %q\n", destination))
	hcl.WriteString(fmt.Sprintf("                local_bind_port = %d\n", localPort))
	if address := toString(upstream["local_bind_address"]); address != "" {
		hcl.WriteString(fmt.Sprintf("                local_bind_address = %q\n", address))
	}
	if datacenter := toString(upstream["datacenter"]); datacenter != "" {
		hcl.WriteString(fmt.Sprintf("                datacenter = %q\n", datacenter))
	}
	if extensions, ok := asMap(upstream["extensions"]); ok && len(extensions) > 0 {
		writeNomadGenericBlock(hcl, "                ", "extensions", extensions)
	}
	hcl.WriteString("              }\n")
}

func writeNomadCheck(hcl *strings.Builder, service *Service, servicePort string) {
	health := normalizeHealthCheck(service.HealthCheck)
	if health == nil {
		return
	}
	if health.Disable {
		return
	}
	checkType := health.Type
	if checkType == "" {
		checkType = "tcp"
	}
	if checkType == "exec" {
		checkType = "script"
	}
	hcl.WriteString("        check {\n")
	hcl.WriteString("          name = \"health\"\n")
	hcl.WriteString(fmt.Sprintf("          type = %q\n", checkType))
	if checkType == "http" && health.Path != "" {
		hcl.WriteString(fmt.Sprintf("          path = %q\n", health.Path))
	}
	if port := nomadHealthPort(service, health, servicePort); port != "" {
		hcl.WriteString(fmt.Sprintf("          port = %q\n", port))
	}
	if health.Interval != "" {
		hcl.WriteString(fmt.Sprintf("          interval = %q\n", health.Interval))
	}
	if health.Timeout != "" {
		hcl.WriteString(fmt.Sprintf("          timeout = %q\n", health.Timeout))
	}
	if checkType == "script" {
		command := healthCheckCommand(health)
		if len(command) > 0 {
			hcl.WriteString(fmt.Sprintf("          command = %q\n", command[0]))
			if len(command) > 1 {
				var args []string
				for _, arg := range command[1:] {
					args = append(args, fmt.Sprintf("%q", arg))
				}
				hcl.WriteString(fmt.Sprintf("          args = [%s]\n", strings.Join(args, ", ")))
			}
		}
	}
	hcl.WriteString("        }\n")
}

func nomadHealthPort(service *Service, health *HealthCheck, servicePort string) string {
	if health == nil {
		return servicePort
	}
	if health.Port == "" {
		return servicePort
	}
	for i, port := range service.Ports {
		if health.Port == port.Name || health.Port == port.ContainerPort || health.Port == port.HostPort {
			return nomadPortLabelForPort(port, i)
		}
	}
	return health.Port
}

func firstNomadPortLabel(service *Service) string {
	if service == nil || len(service.Ports) == 0 {
		return ""
	}
	return nomadPortLabelForPort(service.Ports[0], 0)
}

func stringAttribute(body *hclsyntax.Body, name string) string {
	attr, ok := body.Attributes[name]
	if !ok {
		return ""
	}
	return expressionString(attr.Expr)
}

func intAttribute(body *hclsyntax.Body, name string) int {
	attr, ok := body.Attributes[name]
	if !ok {
		return 0
	}
	value, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return 0
	}
	if value.Type() == cty.Number {
		var i int
		if err := goctyInt(value, &i); err == nil {
			return i
		}
	}
	return parseInt(expressionString(attr.Expr))
}

func boolAttribute(body *hclsyntax.Body, name string) bool {
	attr, ok := body.Attributes[name]
	if !ok {
		return false
	}
	value, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return false
	}
	if value.Type() == cty.Bool {
		return value.True()
	}
	return strings.EqualFold(expressionString(attr.Expr), "true")
}

func stringListAttribute(body *hclsyntax.Body, name string) []string {
	attr, ok := body.Attributes[name]
	if !ok {
		return nil
	}
	value, diags := attr.Expr.Value(nil)
	if diags.HasErrors() || !value.CanIterateElements() {
		return nil
	}
	var result []string
	it := value.ElementIterator()
	for it.Next() {
		_, val := it.Element()
		if val.Type() == cty.String {
			result = append(result, val.AsString())
		}
	}
	return result
}

func stringMapAttribute(body *hclsyntax.Body, name string) map[string]string {
	attr, ok := body.Attributes[name]
	if !ok {
		return nil
	}
	value, diags := attr.Expr.Value(nil)
	if diags.HasErrors() || !value.CanIterateElements() {
		return nil
	}
	result := map[string]string{}
	it := value.ElementIterator()
	for it.Next() {
		keyValue, element := it.Element()
		key := keyValue.AsString()
		if key == "" {
			continue
		}
		if element.Type() == cty.String {
			result[key] = element.AsString()
			continue
		}
		result[key] = toString(element.GoString())
	}
	return result
}

func sortedMapKeys(input map[string]string) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func expressionString(expr hclsyntax.Expression) string {
	value, diags := expr.Value(nil)
	if !diags.HasErrors() && value.Type() == cty.String {
		return value.AsString()
	}
	if !diags.HasErrors() {
		return value.GoString()
	}
	return templateTraversalString(expr)
}

func templateTraversalString(expr hclsyntax.Expression) string {
	switch typed := expr.(type) {
	case *hclsyntax.TemplateWrapExpr:
		if traversal, ok := typed.Wrapped.(*hclsyntax.ScopeTraversalExpr); ok {
			return "${" + hclTraversalString(traversal.Traversal) + "}"
		}
	case *hclsyntax.ScopeTraversalExpr:
		return hclTraversalString(typed.Traversal)
	}
	return ""
}

func hclTraversalString(traversal hcl.Traversal) string {
	parts := make([]string, 0, len(traversal))
	for _, traverser := range traversal {
		switch step := traverser.(type) {
		case hcl.TraverseRoot:
			parts = append(parts, step.Name)
		case hcl.TraverseAttr:
			parts = append(parts, step.Name)
		default:
			return ""
		}
	}
	return strings.Join(parts, ".")
}

func goctyInt(value cty.Value, out *int) error {
	float, _ := value.AsBigFloat().Float64()
	*out = int(float)
	return nil
}

func nomadPortLabel(index int) string {
	return fmt.Sprintf("p%d", index)
}

func nomadPortLabelForPort(port PortMapping, index int) string {
	if port.Name != "" {
		return port.Name
	}
	return nomadPortLabel(index)
}

func nomadCPU(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	quantity, err := resource.ParseQuantity(value)
	if err != nil {
		return strings.TrimSuffix(value, "m")
	}
	return fmt.Sprintf("%d", quantity.MilliValue())
}

func nomadMemory(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	quantity, err := resource.ParseQuantity(value)
	if err != nil {
		value = strings.TrimSuffix(value, "Mi")
		value = strings.TrimSuffix(value, "M")
		return strings.TrimSpace(value)
	}
	const mib = int64(1024 * 1024)
	return fmt.Sprintf("%d", quantity.Value()/mib)
}
