package paas

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	composeloader "github.com/compose-spec/compose-go/v2/loader"
	composetypes "github.com/compose-spec/compose-go/v2/types"
	"gopkg.in/yaml.v3"
)

func parseDockerComposeSpec(content string) (*Application, error) {
	return parseDockerComposeSpecWithFile(content, "")
}

func parseDockerComposeFile(filename, content string) (*Application, error) {
	return parseDockerComposeSpecWithFile(content, filename)
}

func parseDockerComposeSpecWithFile(content, filename string) (*Application, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configFile := composetypes.ConfigFile{Filename: "compose.yml", Content: []byte(content)}
	if filename != "" {
		configFile.Filename = filename
		if dir := filepath.Dir(filename); dir != "" {
			workingDir = dir
		}
	}
	project, err := composeloader.LoadWithContext(context.Background(), composetypes.ConfigDetails{
		WorkingDir: workingDir,
		ConfigFiles: []composetypes.ConfigFile{
			configFile,
		},
		Environment: composetypes.Mapping{},
	}, func(options *composeloader.Options) {
		options.SetProjectName("paas", true)
		options.Profiles = []string{"*"}
		options.SkipConsistencyCheck = true
		options.ResolvePaths = filename != ""
	})
	if err != nil {
		return nil, err
	}

	app := applicationFromComposeProject(project)
	if app.Extensions == nil {
		app.Extensions = map[string]interface{}{}
	}
	app.Extensions[composeRawYAMLExtension] = content
	var raw map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &raw); err == nil {
		if version := toString(raw["version"]); version != "" {
			app.Version = version
		}
		if name := toString(raw["name"]); name != "" {
			app.Name = name
		}
		if platform := toString(raw["x-platform"]); platform == string(PlatformDockerSwarm) {
			app.Platform = PlatformDockerSwarm
			for _, service := range app.Services {
				service.Platform = PlatformDockerSwarm
			}
		}
		if includes, includeEntries := parseComposeIncludeEntries(raw["include"]); len(includeEntries) > 0 {
			app.Includes = includes
			app.IncludeEntries = includeEntries
		}
		if routes := raw[composeAppRoutesExtension]; routes != nil {
			app.Extensions[composeAppRoutesExtension] = deepCopyValue(routes)
		}
		if policies := raw[composeAppPoliciesExtension]; policies != nil {
			app.Extensions[composeAppPoliciesExtension] = deepCopyValue(policies)
		}
		if resources := raw[composeCanonicalRawResourcesExtension]; resources != nil {
			app.Extensions[composeCanonicalRawResourcesExtension] = deepCopyValue(resources)
		}
		if resources := raw[composeKubernetesRawResourcesExtension]; resources != nil {
			app.Extensions[composeKubernetesRawResourcesExtension] = deepCopyValue(resources)
		}
		mergeRawComposeServices(app, raw)
		mergeRawComposeNetworks(app, raw)
		mergeRawComposeVolumes(app, raw)
		mergeRawComposeConfigs(app, raw)
		mergeRawComposeSecrets(app, raw)
	}
	syncPortableApplicationState(app)
	app.AttachCanonical()
	if project != nil {
		app.Canonical.AddResource(ResourceKindRaw, app.Platform, "compose-project", "compose-go.Project", project)
	}
	app.Canonical.AddResource(ResourceKindRaw, app.Platform, "compose-yaml", "ComposeYAML", content)
	return app, nil
}

func mergeRawComposeServices(app *Application, raw map[string]interface{}) {
	if app == nil || len(raw) == 0 {
		return
	}
	rawServices, ok := asMap(raw["services"])
	if !ok || len(rawServices) == 0 {
		return
	}
	for name, value := range rawServices {
		serviceMap, ok := asMap(value)
		if !ok {
			continue
		}
		service := app.Services[name]
		if service == nil {
			continue
		}
		mergeRawComposeServiceCompat(service, serviceMap)
	}
}

func mergeRawComposeServiceCompat(service *Service, raw map[string]interface{}) {
	if service == nil || len(raw) == 0 {
		return
	}
	if volumesRaw, ok := raw["volumes"]; ok {
		mergeRawComposeServiceVolumes(service, volumesRaw)
	}
	if configsRaw, ok := raw["configs"]; ok {
		mergeRawComposeServiceFileRefs(service.Configs, configsRaw)
	}
	if secretsRaw, ok := raw["secrets"]; ok {
		mergeRawComposeServiceFileRefs(service.Secrets, secretsRaw)
	}
	if envFileRaw, ok := raw["env_file"]; ok {
		mergeRawComposeServiceEnvFiles(service.EnvFileRefs, envFileRaw)
	}
	if buildRaw, ok := raw["build"]; ok {
		mergeRawComposeBuild(service, buildRaw)
	}
	if healthcheckRaw, ok := raw["healthcheck"]; ok {
		if parsed, err := parseHealthCheck(healthcheckRaw); err == nil {
			service.HealthCheck = mergeHealthCheckSpec(service.HealthCheck, parsed)
		}
	}
	if deployRaw, ok := raw["deploy"]; ok {
		mergeRawComposeDeploy(service, deployRaw)
	}
	if ulimitsRaw, ok := raw["ulimits"]; ok {
		mergeRawComposeUlimits(service, ulimitsRaw)
	}
	if pidsLimit, ok := raw["pids_limit"]; ok {
		service.PidsLimit = int64(toInt(pidsLimit))
		service.pidsLimitSet = true
	}
	if shmSize, ok := raw["shm_size"]; ok {
		service.ShmSize = int64(toInt(shmSize))
		service.shmSizeSet = true
	}
	if blkioRaw, ok := raw["blkio_config"]; ok {
		mergeRawComposeBlkioConfig(service, blkioRaw)
	}
	if blkioExtensionsRaw, ok := raw[composeBlkioConfigExtensionsExtension]; ok {
		mergeRawComposeBlkioConfigExtensions(service, blkioExtensionsRaw)
	}
	if value, ok := raw["privileged"]; ok {
		service.Privileged = toBool(value)
		service.PrivilegedSet = true
	}
	if value, ok := raw["read_only"]; ok {
		service.ReadOnlyRootFS = toBool(value)
		service.ReadOnlyRootFSSet = true
	}
	if value, ok := raw["tty"]; ok {
		service.Tty = toBool(value)
		service.TtySet = true
	}
	if value, ok := raw["stdin_open"]; ok {
		service.StdinOpen = toBool(value)
		service.StdinOpenSet = true
	}
	if value, ok := raw["cpu_count"]; ok {
		compat := ensureComposeCompat(service)
		compat.CPUCount = int64(toInt(value))
		compat.CPUCountSet = true
	}
	if value, ok := raw["cpu_percent"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprint(value), 32); err == nil {
			compat := ensureComposeCompat(service)
			compat.CPUPercent = float32(f)
			compat.CPUPercentSet = true
		}
	}
	if value, ok := raw["cpu_period"]; ok {
		compat := ensureComposeCompat(service)
		compat.CPUPeriod = int64(toInt(value))
		compat.CPUPeriodSet = true
	}
	if value, ok := raw["cpu_rt_period"]; ok {
		compat := ensureComposeCompat(service)
		compat.CPURTPeriod = int64(toInt(value))
		compat.CPURTPeriodSet = true
	}
	if value, ok := raw["cpu_rt_runtime"]; ok {
		compat := ensureComposeCompat(service)
		compat.CPURTRuntime = int64(toInt(value))
		compat.CPURTRuntimeSet = true
	}
	if value, ok := raw["oom_score_adj"]; ok {
		compat := ensureComposeCompat(service)
		compat.OomScoreAdj = int64(toInt(value))
		compat.OomScoreAdjSet = true
	}
	if value, ok := raw["oom_kill_disable"]; ok {
		compat := ensureComposeCompat(service)
		compat.OomKillDisable = toBool(value)
		compat.OomKillDisableSet = true
	}
	if value, ok := raw["use_api_socket"]; ok {
		compat := ensureComposeCompat(service)
		compat.UseAPISocket = toBool(value)
		compat.UseAPISocketSet = true
	}
	compat := service.ComposeCompat
	for _, key := range []string{"extends", "platform", "pull_refresh_after", "mem_swappiness"} {
		value, ok := raw[key]
		if !ok {
			continue
		}
		if compat == nil {
			compat = ensureComposeCompat(service)
		}
		switch key {
		case "extends":
			compat.Extends = composeExtendsRawMap(value)
		case "platform":
			compat.Platform = toString(value)
		case "pull_refresh_after":
			compat.PullRefreshAfter = toString(value)
		case "mem_swappiness":
			compat.MemSwappiness = toString(value)
		}
	}
	if isEmptyComposeCompat(compat) {
		service.ComposeCompat = nil
	}
	hydrateServiceFailoverFromExtensions(service)
}

func mergeRawComposeServiceVolumes(service *Service, raw interface{}) {
	if service == nil {
		return
	}
	rawVolumes, ok := raw.([]interface{})
	if !ok || len(rawVolumes) == 0 || len(service.Volumes) == 0 {
		return
	}
	for i, item := range rawVolumes {
		if i >= len(service.Volumes) {
			break
		}
		rawMount, ok := asMap(item)
		if !ok {
			continue
		}
		mergeRawComposeVolumeMount(&service.Volumes[i], rawMount)
	}
}

func mergeRawComposeVolumeMount(mount *VolumeMount, raw map[string]interface{}) {
	if mount == nil || len(raw) == 0 {
		return
	}
	if propagation, ok := raw["x-kubernetes-mountPropagation"]; ok {
		mount.MountPropagation = toString(propagation)
	} else if propagation, ok := raw["x-kubernetes-mount-propagation"]; ok {
		mount.MountPropagation = toString(propagation)
	}
	for key, value := range raw {
		if mount.Extensions == nil {
			mount.Extensions = map[string]interface{}{}
		}
		mount.Extensions[key] = deepCopyValue(value)
	}
	if bindRaw, ok := asMap(raw["bind"]); ok {
		if mount.BindExtensions == nil {
			mount.BindExtensions = map[string]interface{}{}
		}
		for key, value := range bindRaw {
			mount.BindExtensions[key] = deepCopyValue(value)
		}
	}
	if volumeRaw, ok := asMap(raw["volume"]); ok {
		if mount.VolumeExtensions == nil {
			mount.VolumeExtensions = map[string]interface{}{}
		}
		for key, value := range volumeRaw {
			mount.VolumeExtensions[key] = deepCopyValue(value)
		}
	}
	if tmpfsRaw, ok := asMap(raw["tmpfs"]); ok {
		if mount.TmpfsExtensions == nil {
			mount.TmpfsExtensions = map[string]interface{}{}
		}
		for key, value := range tmpfsRaw {
			mount.TmpfsExtensions[key] = deepCopyValue(value)
		}
	}
	if imageRaw, ok := asMap(raw["image"]); ok {
		if mount.ImageExtensions == nil {
			mount.ImageExtensions = map[string]interface{}{}
		}
		for key, value := range imageRaw {
			mount.ImageExtensions[key] = deepCopyValue(value)
		}
	}
	if value, ok := raw["x-kubernetes-subPath"]; ok {
		mount.SubPath = toString(value)
	}
	if value, ok := raw["x-kubernetes-subPathExpr"]; ok {
		mount.SubPathExpr = strings.ReplaceAll(toString(value), "$$", "$")
	}
	if value, ok := raw["x-kubernetes-mountPropagation"]; ok {
		mount.MountPropagation = toString(value)
	}
	if value, ok := raw["x-kubernetes-mount-propagation"]; ok {
		mount.MountPropagation = toString(value)
	}
	if value, ok := raw["x-kubernetes-recursiveReadOnly"]; ok {
		mount.RecursiveReadOnly = toString(value)
	}
	if value, ok := raw["x-kubernetes-recursive-read-only"]; ok {
		mount.RecursiveReadOnly = toString(value)
	}
}

func mergeRawComposeServiceFileRefs(refs []FileRef, raw interface{}) {
	rawRefs, ok := raw.([]interface{})
	if !ok || len(rawRefs) == 0 || len(refs) == 0 {
		return
	}
	for i, item := range rawRefs {
		if i >= len(refs) {
			break
		}
		rawRef, ok := asMap(item)
		if !ok {
			continue
		}
		for key, value := range rawRef {
			if refs[i].Extensions == nil {
				refs[i].Extensions = map[string]interface{}{}
			}
			refs[i].Extensions[key] = deepCopyValue(value)
		}
	}
}

func mergeRawComposeServiceEnvFiles(refs []EnvFileRef, raw interface{}) {
	rawRefs, ok := raw.([]interface{})
	if !ok || len(rawRefs) == 0 || len(refs) == 0 {
		return
	}
	for i, item := range rawRefs {
		if i >= len(refs) {
			break
		}
		rawRef, ok := asMap(item)
		if !ok {
			continue
		}
		for key, value := range rawRef {
			if refs[i].Extensions == nil {
				refs[i].Extensions = map[string]interface{}{}
			}
			refs[i].Extensions[key] = deepCopyValue(value)
		}
	}
}

func mergeRawComposeUlimits(service *Service, raw interface{}) {
	if service == nil || raw == nil {
		return
	}
	parsed, err := parseUlimits(raw)
	if err != nil || parsed == nil {
		return
	}
	service.Ulimits = mergeUlimits(service.Ulimits, parsed)
}

func mergeUlimits(base, overlay *Ulimits) *Ulimits {
	if base == nil {
		return cloneUlimits(overlay)
	}
	if overlay == nil {
		return cloneUlimits(base)
	}
	result := cloneUlimits(base)
	if len(overlay.Limits) > 0 {
		if result.Limits == nil {
			result.Limits = map[string]UlimitSpec{}
		}
		for name, limit := range overlay.Limits {
			existing := result.Limits[name]
			result.Limits[name] = mergeUlimitSpec(existing, limit)
		}
	}
	if overlay.Nofile != nil {
		if result.Nofile == nil {
			nofile := *overlay.Nofile
			nofile.Extensions = copyStringInterfaceMap(overlay.Nofile.Extensions)
			result.Nofile = &nofile
		} else {
			if overlay.Nofile.Soft > 0 {
				result.Nofile.Soft = overlay.Nofile.Soft
			}
			if overlay.Nofile.Hard > 0 {
				result.Nofile.Hard = overlay.Nofile.Hard
			}
			if len(overlay.Nofile.Extensions) > 0 {
				if result.Nofile.Extensions == nil {
					result.Nofile.Extensions = map[string]interface{}{}
				}
				for key, value := range overlay.Nofile.Extensions {
					result.Nofile.Extensions[key] = value
				}
			}
		}
	}
	return result
}

func mergeUlimitSpec(base, overlay UlimitSpec) UlimitSpec {
	result := base
	if overlay.Single > 0 {
		result.Single = overlay.Single
	}
	if overlay.Soft > 0 {
		result.Soft = overlay.Soft
	}
	if overlay.Hard > 0 {
		result.Hard = overlay.Hard
	}
	if len(overlay.Extensions) > 0 {
		if result.Extensions == nil {
			result.Extensions = map[string]interface{}{}
		}
		for key, value := range overlay.Extensions {
			result.Extensions[key] = deepCopyValue(value)
		}
	}
	return result
}

func mergeRawComposeBlkioConfig(service *Service, raw interface{}) {
	if service == nil || raw == nil {
		return
	}
	mapped, ok := asMap(raw)
	if !ok || len(mapped) == 0 {
		return
	}
	compat := ensureComposeCompat(service)
	compat.BlkioConfig = mergeComposeBlkioConfigMap(compat.BlkioConfig, mapped)
}

func mergeRawComposeBlkioConfigExtensions(service *Service, raw interface{}) {
	mergeComposeBlkioConfigExtensions(ensureComposeCompat(service), raw)
}

func mergeComposeBlkioConfigExtensions(compat *ComposeCompat, raw interface{}) {
	if compat == nil || raw == nil {
		return
	}
	mapped, ok := asMap(raw)
	if !ok || len(mapped) == 0 {
		return
	}
	if compat.BlkioConfig == nil {
		compat.BlkioConfig = map[string]interface{}{}
	}
	for key, value := range mapped {
		switch key {
		case "weight_device", "device_read_bps", "device_read_iops", "device_write_bps", "device_write_iops":
			extraList, ok := value.([]interface{})
			if !ok || len(extraList) == 0 {
				continue
			}
			merged := blkioExtensionList(compat.BlkioConfig[key], extraList)
			if len(merged) > 0 {
				compat.BlkioConfig[key] = merged
			}
		default:
			compat.BlkioConfig[key] = deepCopyValue(value)
		}
	}
}

func blkioExtensionList(existing interface{}, extras []interface{}) []interface{} {
	merged, ok := existing.([]interface{})
	if !ok || len(merged) == 0 {
		merged = make([]interface{}, 0, len(extras))
	}
	for index, extra := range extras {
		extraMap, ok := asMap(extra)
		if !ok || len(extraMap) == 0 {
			continue
		}
		if index >= len(merged) {
			merged = append(merged, map[string]interface{}{})
		}
		official, ok := asMap(merged[index])
		if !ok || official == nil {
			official = map[string]interface{}{}
		}
		for key, value := range extraMap {
			official[key] = deepCopyValue(value)
		}
		merged[index] = official
	}
	return merged
}

func mergeComposeBlkioConfigMap(base, overlay map[string]interface{}) map[string]interface{} {
	if len(overlay) == 0 {
		return copyStringInterfaceMap(base)
	}
	result := copyStringInterfaceMap(base)
	if result == nil {
		result = map[string]interface{}{}
	}
	for key, value := range overlay {
		switch key {
		case "weight":
			result[key] = deepCopyValue(value)
		case "weight_device":
			result[key] = mergeComposeBlkioDeviceList(result[key], value)
		case "device_read_bps", "device_read_iops", "device_write_bps", "device_write_iops":
			result[key] = mergeComposeBlkioDeviceList(result[key], value)
		default:
			result[key] = deepCopyValue(value)
		}
	}
	return result
}

func mergeComposeBlkioDeviceList(existing interface{}, overlay interface{}) []interface{} {
	items, ok := overlay.([]interface{})
	if !ok || len(items) == 0 {
		if current, ok := existing.([]interface{}); ok {
			return current
		}
		return nil
	}
	merged, ok := existing.([]interface{})
	if !ok || len(merged) == 0 {
		merged = make([]interface{}, 0, len(items))
	}
	for index, item := range items {
		officialMap, ok := asMap(item)
		if !ok || len(officialMap) == 0 {
			continue
		}
		if index >= len(merged) {
			merged = append(merged, map[string]interface{}{})
		}
		existingMap, ok := asMap(merged[index])
		if !ok || existingMap == nil {
			existingMap = map[string]interface{}{}
		}
		for key, value := range officialMap {
			existingMap[key] = deepCopyValue(value)
		}
		merged[index] = existingMap
	}
	return merged
}

func mergeInterfaceMaps(base, overlay map[string]interface{}) map[string]interface{} {
	if len(base) == 0 {
		return copyStringInterfaceMap(overlay)
	}
	if len(overlay) == 0 {
		return copyStringInterfaceMap(base)
	}
	result := copyStringInterfaceMap(base)
	for key, value := range overlay {
		result[key] = deepCopyValue(value)
	}
	return result
}

func mergeRawComposeDeploy(service *Service, raw interface{}) {
	if service == nil || raw == nil {
		return
	}
	parsed, err := parseDeploySpec(raw)
	if err != nil || parsed == nil {
		return
	}
	if service.Deploy == nil {
		service.Deploy = parsed
		return
	}
	if parsed.Job != nil && service.Deploy.Job == nil {
		service.Deploy.Job = cloneSwarmJobSpec(parsed.Job)
	}
	if len(parsed.Extensions) > 0 {
		if service.Deploy.Extensions == nil {
			service.Deploy.Extensions = map[string]interface{}{}
		}
		for key, value := range parsed.Extensions {
			service.Deploy.Extensions[key] = value
		}
	}
	mergeRawComposePlacementSpec(&service.Deploy.Placement, parsed.Placement)
	mergeRawComposeResourceSpec(&service.Deploy.Resources, parsed.Resources)
	if parsed.Resources != nil && len(parsed.Resources.Extensions) > 0 {
		if claims, ok := parsed.Resources.Extensions["x-kubernetes-claims"]; ok {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-claims"] = deepCopyValue(claims)
		}
		if claims, ok := parsed.Resources.Extensions["kubernetes.claims"]; ok {
			if service.Extensions == nil {
				service.Extensions = map[string]interface{}{}
			}
			service.Extensions["x-kubernetes-claims"] = deepCopyValue(claims)
		}
	}
	mergeRawComposeUpdatePolicy(&service.Deploy.UpdateConfig, parsed.UpdateConfig)
	mergeRawComposeMigratePolicy(&service.Deploy.MigrateConfig, parsed.MigrateConfig)
	mergeRawComposeReschedulePolicy(&service.Deploy.RescheduleConfig, parsed.RescheduleConfig)
	mergeRawComposeUpdatePolicy(&service.Deploy.RollbackConfig, parsed.RollbackConfig)
	mergeRawComposeRestartPolicy(&service.Deploy.RestartPolicy, parsed.RestartPolicy)
}

func mergeRawComposePlacementSpec(target **PlacementSpec, raw *PlacementSpec) {
	if raw == nil {
		return
	}
	if *target == nil {
		*target = clonePlacementSpec(raw)
		return
	}
	if len(raw.Extensions) > 0 {
		(*target).Extensions = copyStringInterfaceMap(raw.Extensions)
	}
	if hasPlacementPreferenceExtensions(raw.PreferenceExtensions) {
		if len((*target).PreferenceExtensions) < len((*target).Preferences) {
			extensions := make([]map[string]interface{}, len((*target).Preferences))
			copy(extensions, (*target).PreferenceExtensions)
			(*target).PreferenceExtensions = extensions
		}
		for i, extension := range raw.PreferenceExtensions {
			if i >= len((*target).PreferenceExtensions) {
				break
			}
			if len(extension) > 0 {
				(*target).PreferenceExtensions[i] = copyStringInterfaceMap(extension)
			}
		}
	}
}

func mergeRawComposeResourceSpec(target **ResourceSpec, raw *ResourceSpec) {
	if raw == nil {
		return
	}
	if *target == nil {
		*target = cloneResourceSpec(raw)
		return
	}
	if len(raw.Extensions) > 0 {
		(*target).Extensions = copyStringInterfaceMap(raw.Extensions)
	}
	if len(raw.LimitExtensions) > 0 {
		(*target).LimitExtensions = copyStringInterfaceMap(raw.LimitExtensions)
	}
	if len(raw.ReservationExtensions) > 0 {
		(*target).ReservationExtensions = copyStringInterfaceMap(raw.ReservationExtensions)
	}
	for i := range raw.Devices {
		if i >= len((*target).Devices) {
			(*target).Devices = append((*target).Devices, cloneResourceDevices(raw.Devices[i:i+1])...)
			continue
		}
		if len(raw.Devices[i].Extensions) > 0 {
			(*target).Devices[i].Extensions = copyStringInterfaceMap(raw.Devices[i].Extensions)
		}
	}
	for i := range raw.GenericResources {
		if i >= len((*target).GenericResources) {
			(*target).GenericResources = append((*target).GenericResources, cloneGenericResources(raw.GenericResources[i:i+1])...)
			continue
		}
		if len(raw.GenericResources[i].Extensions) > 0 {
			(*target).GenericResources[i].Extensions = copyStringInterfaceMap(raw.GenericResources[i].Extensions)
		}
		if len(raw.GenericResources[i].DiscreteExtensions) > 0 {
			(*target).GenericResources[i].DiscreteExtensions = copyStringInterfaceMap(raw.GenericResources[i].DiscreteExtensions)
		}
	}
}

func mergeRawComposeUpdatePolicy(target **UpdatePolicy, raw *UpdatePolicy) {
	if raw == nil {
		return
	}
	if *target == nil {
		*target = cloneUpdatePolicy(raw)
		return
	}
	merged := cloneUpdatePolicy(*target)
	if raw.ParallelismSet || raw.Parallelism > 0 {
		merged.Parallelism = raw.Parallelism
		merged.ParallelismSet = raw.ParallelismSet || raw.Parallelism > 0
	}
	if raw.Delay != "" {
		merged.Delay = raw.Delay
	}
	if raw.Monitor != "" {
		merged.Monitor = raw.Monitor
	}
	if raw.MaxFailureRatio != "" {
		merged.MaxFailureRatio = raw.MaxFailureRatio
	}
	if raw.Order != "" {
		merged.Order = raw.Order
	}
	if raw.OnFailure != "" {
		merged.OnFailure = raw.OnFailure
	}
	if raw.HealthCheck != "" {
		merged.HealthCheck = raw.HealthCheck
	}
	if raw.MinHealthyTime != "" {
		merged.MinHealthyTime = raw.MinHealthyTime
	}
	if raw.HealthyDeadline != "" {
		merged.HealthyDeadline = raw.HealthyDeadline
	}
	if raw.ProgressDeadline != "" {
		merged.ProgressDeadline = raw.ProgressDeadline
	}
	if raw.AutoRevertSet {
		merged.AutoRevert = raw.AutoRevert
		merged.AutoRevertSet = true
	}
	if raw.AutoPromoteSet {
		merged.AutoPromote = raw.AutoPromote
		merged.AutoPromoteSet = true
	}
	if raw.CanarySet {
		merged.Canary = raw.Canary
		merged.CanarySet = true
	}
	if raw.Stagger != "" {
		merged.Stagger = raw.Stagger
	}
	if len(raw.Extensions) > 0 {
		if merged.Extensions == nil {
			merged.Extensions = map[string]interface{}{}
		}
		for key, value := range raw.Extensions {
			merged.Extensions[key] = value
		}
	}
	*target = merged
}

func mergeRawComposeMigratePolicy(target **MigratePolicy, raw *MigratePolicy) {
	if raw == nil {
		return
	}
	if *target == nil {
		*target = cloneMigratePolicy(raw)
		return
	}
	merged := cloneMigratePolicy(*target)
	if raw.MaxParallel > 0 {
		merged.MaxParallel = raw.MaxParallel
	}
	if raw.HealthCheck != "" {
		merged.HealthCheck = raw.HealthCheck
	}
	if raw.MinHealthyTime != "" {
		merged.MinHealthyTime = raw.MinHealthyTime
	}
	if raw.HealthyDeadline != "" {
		merged.HealthyDeadline = raw.HealthyDeadline
	}
	if len(raw.Extensions) > 0 {
		if merged.Extensions == nil {
			merged.Extensions = map[string]interface{}{}
		}
		for key, value := range raw.Extensions {
			merged.Extensions[key] = value
		}
	}
	*target = merged
}

func mergeRawComposeReschedulePolicy(target **ReschedulePolicy, raw *ReschedulePolicy) {
	if raw == nil {
		return
	}
	if *target == nil {
		*target = cloneReschedulePolicy(raw)
		return
	}
	merged := cloneReschedulePolicy(*target)
	if raw.Attempts > 0 {
		merged.Attempts = raw.Attempts
	}
	if raw.Interval != "" {
		merged.Interval = raw.Interval
	}
	if raw.Delay != "" {
		merged.Delay = raw.Delay
	}
	if raw.DelayFunction != "" {
		merged.DelayFunction = raw.DelayFunction
	}
	if raw.MaxDelay != "" {
		merged.MaxDelay = raw.MaxDelay
	}
	if raw.Unlimited {
		merged.Unlimited = true
	}
	if len(raw.Extensions) > 0 {
		if merged.Extensions == nil {
			merged.Extensions = map[string]interface{}{}
		}
		for key, value := range raw.Extensions {
			merged.Extensions[key] = value
		}
	}
	*target = merged
}

func mergeRawComposeRestartPolicy(target **RestartPolicy, raw *RestartPolicy) {
	if raw == nil {
		return
	}
	if *target == nil {
		*target = cloneRestartPolicy(raw)
		return
	}
	if len(raw.Extensions) == 0 {
		return
	}
	if (*target).Extensions == nil {
		(*target).Extensions = map[string]interface{}{}
	}
	for key, value := range raw.Extensions {
		(*target).Extensions[key] = value
	}
}

func mergeRawComposeBuild(service *Service, raw interface{}) {
	if service == nil || raw == nil {
		return
	}
	if service.Build == nil {
		if parsed, err := parseBuildConfig(raw); err == nil {
			service.Build = parsed
		}
	}
	if service.Build == nil {
		return
	}
	if service.Build.Extensions == nil {
		service.Build.Extensions = map[string]interface{}{}
	}
	switch typed := raw.(type) {
	case string:
		if typed != "" {
			service.Build.Extensions["compose.build"] = typed
		}
	default:
		if mapped, ok := asMap(raw); ok && len(mapped) > 0 {
			service.Build.Extensions["compose.build"] = copyStringInterfaceMap(mapped)
		}
	}
	if len(service.Build.Extensions) == 0 {
		service.Build.Extensions = nil
	}
}

func composeExtendsRawMap(value interface{}) map[string]interface{} {
	switch typed := value.(type) {
	case string:
		if typed == "" {
			return nil
		}
		return map[string]interface{}{"service": typed}
	case map[string]interface{}:
		return copyStringInterfaceMap(typed)
	default:
		mapped, ok := asMap(value)
		if !ok {
			return nil
		}
		return copyStringInterfaceMap(mapped)
	}
}

func mergeRawComposeNetworks(app *Application, raw map[string]interface{}) {
	if app == nil || len(raw) == 0 {
		return
	}
	rawNetworks, ok := asMap(raw["networks"])
	if !ok || len(rawNetworks) == 0 {
		return
	}
	if app.Networks == nil {
		app.Networks = map[string]*Network{}
	}
	for name, value := range rawNetworks {
		networkMap, ok := asMap(value)
		if !ok {
			continue
		}
		network, err := parseNetwork(name, networkMap)
		if err != nil || network == nil {
			continue
		}
		app.Networks[name] = mergeComposeNetworkDefinition(app.Networks[name], network)
	}
}

func mergeComposeNetworkDefinition(existing, raw *Network) *Network {
	if existing == nil {
		return cloneNetwork(raw)
	}
	if raw == nil {
		return cloneNetwork(existing)
	}
	merged := cloneNetwork(existing)
	if raw.Name != "" {
		merged.Name = raw.Name
	}
	if raw.PlatformName != "" {
		merged.PlatformName = raw.PlatformName
	}
	if raw.Driver != "" {
		merged.Driver = raw.Driver
	}
	if len(raw.DriverOpts) > 0 {
		merged.DriverOpts = copyStringMap(raw.DriverOpts)
	}
	if raw.AttachableSet {
		merged.Attachable = raw.Attachable
		merged.AttachableSet = true
	} else if raw.Attachable {
		merged.Attachable = true
	}
	if raw.ExternalSet {
		merged.External = raw.External
		merged.ExternalSet = true
	} else if raw.External {
		merged.External = true
	}
	if len(raw.ExternalExtensions) > 0 {
		merged.ExternalExtensions = copyStringInterfaceMap(raw.ExternalExtensions)
	}
	if raw.InternalSet {
		merged.Internal = raw.Internal
		merged.InternalSet = true
	} else if raw.Internal {
		merged.Internal = true
	}
	if raw.EnableIPv4 != nil {
		merged.EnableIPv4 = cloneBoolPtr(raw.EnableIPv4)
	}
	if raw.EnableIPv6 != nil {
		merged.EnableIPv6 = cloneBoolPtr(raw.EnableIPv6)
	}
	if raw.IPAM != nil {
		merged.IPAM = cloneIPAMConfig(raw.IPAM)
	}
	if len(raw.Labels) > 0 {
		merged.Labels = copyStringMap(raw.Labels)
	}
	if len(raw.Extensions) > 0 {
		merged.Extensions = copyStringInterfaceMap(raw.Extensions)
	}
	return merged
}

func mergeRawComposeVolumes(app *Application, raw map[string]interface{}) {
	if app == nil || len(raw) == 0 {
		return
	}
	rawVolumes, ok := asMap(raw["volumes"])
	if !ok || len(rawVolumes) == 0 {
		return
	}
	if app.Volumes == nil {
		app.Volumes = map[string]*Volume{}
	}
	for name, value := range rawVolumes {
		volumeMap, ok := asMap(value)
		if !ok {
			continue
		}
		volume, err := parseVolume(name, volumeMap)
		if err != nil || volume == nil {
			continue
		}
		app.Volumes[name] = cloneVolume(volume)
	}
}

func mergeRawComposeConfigs(app *Application, raw map[string]interface{}) {
	if app == nil || len(raw) == 0 {
		return
	}
	rawConfigs, ok := asMap(raw["configs"])
	if !ok || len(rawConfigs) == 0 {
		return
	}
	if app.Configs == nil {
		app.Configs = map[string]*Config{}
	}
	for name, value := range rawConfigs {
		configMap, ok := asMap(value)
		if !ok {
			continue
		}
		config, err := parseConfig(name, configMap)
		if err != nil || config == nil {
			continue
		}
		app.Configs[name] = cloneConfig(config)
	}
}

func mergeRawComposeSecrets(app *Application, raw map[string]interface{}) {
	if app == nil || len(raw) == 0 {
		return
	}
	rawSecrets, ok := asMap(raw["secrets"])
	if !ok || len(rawSecrets) == 0 {
		return
	}
	if app.Secrets == nil {
		app.Secrets = map[string]*Secret{}
	}
	for name, value := range rawSecrets {
		secretMap, ok := asMap(value)
		if !ok {
			continue
		}
		secret, err := parseSecret(name, secretMap)
		if err != nil || secret == nil {
			continue
		}
		app.Secrets[name] = cloneSecret(secret)
	}
}

func applicationFromComposeProject(project *composetypes.Project) *Application {
	app := &Application{
		Platform:   PlatformDockerCompose,
		Services:   make(map[string]*Service),
		Networks:   make(map[string]*Network),
		Volumes:    make(map[string]*Volume),
		Configs:    make(map[string]*Config),
		Secrets:    make(map[string]*Secret),
		Models:     make(map[string]*ComposeModel),
		Extensions: make(map[string]interface{}),
	}
	if project == nil {
		return app
	}
	if project.Name != "" {
		app.Name = project.Name
		app.Extensions["name"] = project.Name
	}
	for key, value := range project.Extensions {
		app.Extensions[key] = value
	}
	rehydrateComposeApplicationExtensions(app)
	for name, svc := range project.Services {
		app.Services[name] = serviceFromCompose(name, svc)
	}
	for name, network := range project.Networks {
		app.Networks[name] = networkFromCompose(name, network)
	}
	for name, volume := range project.Volumes {
		app.Volumes[name] = volumeFromCompose(name, volume)
	}
	for name, config := range project.Configs {
		app.Configs[name] = configFromCompose(name, config)
	}
	for name, secret := range project.Secrets {
		app.Secrets[name] = secretFromCompose(name, secret)
	}
	for name, model := range project.Models {
		if model := composeModelFromCompose(name, model); model != nil {
			app.Models[name] = model
		}
	}
	if app.Platform == PlatformDockerSwarm || composeAppLooksLikeSwarm(app) {
		app.Platform = PlatformDockerSwarm
		for _, service := range app.Services {
			if service != nil {
				service.Platform = PlatformDockerSwarm
			}
		}
	}
	return app
}

func serviceFromCompose(name string, svc composetypes.ServiceConfig) *Service {
	service := &Service{
		Name:           name,
		Image:          svc.Image,
		ContainerName:  svc.ContainerName,
		Hostname:       svc.Hostname,
		HostPID:        boolPtr(svc.Pid == "host"),
		HostIPC:        boolPtr(svc.Ipc == "host"),
		PIDMode:        svc.Pid,
		IPCMode:        svc.Ipc,
		DNS:            append([]string{}, []string(svc.DNS)...),
		DNSSearch:      append([]string{}, []string(svc.DNSSearch)...),
		DNSOptions:     append([]string{}, svc.DNSOpts...),
		ExtraHosts:     svc.ExtraHosts.AsList("="),
		WorkingDir:     svc.WorkingDir,
		Restart:        svc.Restart,
		Privileged:     svc.Privileged,
		User:           svc.User,
		Runtime:        svc.Runtime,
		LogDriver:      svc.LogDriver,
		LogOpt:         copyStringMap(svc.LogOpt),
		PidsLimit:      svc.PidsLimit,
		ShmSize:        int64(svc.ShmSize),
		GroupAdd:       append([]string{}, svc.GroupAdd...),
		Sysctls:        copyStringMap(map[string]string(svc.Sysctls)),
		CapAdd:         append([]string{}, svc.CapAdd...),
		CapDrop:        append([]string{}, svc.CapDrop...),
		SecurityOpt:    append([]string{}, svc.SecurityOpt...),
		ReadOnlyRootFS: svc.ReadOnly,
		Init:           composeBoolPtr(svc.Init),
		Tty:            svc.Tty,
		StdinOpen:      svc.StdinOpen,
		StopSignal:     svc.StopSignal,
		Platform:       PlatformDockerCompose,
		Environment:    map[string]string{},
		Labels:         map[string]string{},
		Extensions:     map[string]interface{}{},
	}
	service.ComposeCompat = composeCompatFromServiceConfig(svc)
	if svc.StopGracePeriod != nil {
		service.StopGracePeriod = svc.StopGracePeriod.String()
	}
	service.Command = []string(svc.Command)
	service.Entrypoint = []string(svc.Entrypoint)
	service.Build = buildFromCompose(svc.Build)
	service.Develop = developFromCompose(svc.Develop)
	service.Lifecycle = lifecycleFromCompose(svc.PreStart, svc.PostStart, svc.PreStop)
	service.Devices, service.DeviceMappings = devicesFromCompose(svc.Devices)
	service.Profiles = append([]string{}, svc.Profiles...)
	service.Expose = append([]string{}, stringOrNumberListToStrings(svc.Expose)...)
	service.CPUShares = int(svc.CPUShares)
	service.CPUQuota = int(svc.CPUQuota)
	if svc.MemLimit != 0 {
		memLimit := fmt.Sprintf("%d", svc.MemLimit)
		service.MemoryLimit = memLimit
		service.MemLimit = memLimit
	}
	if svc.MemSwapLimit != 0 {
		service.MemorySwap = fmt.Sprintf("%d", svc.MemSwapLimit)
	}
	if svc.MemReservation != 0 {
		service.MemReservation = fmt.Sprintf("%d", svc.MemReservation)
	}
	if svc.CPUS != 0 {
		service.CPUs = fmt.Sprintf("%g", svc.CPUS)
	}
	if len(svc.Ulimits) > 0 {
		service.Ulimits = ulimitsFromCompose(svc.Ulimits)
	}
	service.UserNSMode = svc.UserNSMode
	service.PullPolicy = svc.PullPolicy
	if svc.Logging != nil {
		if service.LogDriver == "" {
			service.LogDriver = svc.Logging.Driver
		}
		if len(service.LogOpt) == 0 {
			service.LogOpt = copyStringMap(map[string]string(svc.Logging.Options))
		}
		for key, value := range svc.Logging.Extensions {
			if service.LogExtensions == nil {
				service.LogExtensions = map[string]interface{}{}
			}
			service.LogExtensions[key] = deepCopyValue(value)
		}
	}
	if svc.Provider != nil {
		if compat := ensureComposeCompat(service); compat != nil {
			compat.Provider = map[string]interface{}{}
			if svc.Provider.Type != "" {
				compat.Provider["type"] = svc.Provider.Type
			}
			if len(svc.Provider.Options) > 0 {
				compat.Provider["options"] = cloneMultiOptions(svc.Provider.Options)
			}
			for key, value := range svc.Provider.Extensions {
				compat.Provider[key] = deepCopyValue(value)
			}
		}
	}
	if value, ok := svc.Extensions[composeBlkioConfigExtensionsExtension]; ok {
		mergeComposeBlkioConfigExtensions(ensureComposeCompat(service), value)
	}
	if value, ok := svc.Extensions["x-nomad-spread"]; ok {
		if len(service.Spreads) == 0 {
			service.Spreads = nomadSpreadSpecsFromAny(value)
		}
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["x-nomad-spread"] = deepCopyValue(value)
	}
	if value, ok := svc.Extensions["x-nomad-connect"]; ok {
		if service.Connect == nil {
			service.Connect = nomadConnectSpecFromAny(value)
		}
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["x-nomad-connect"] = deepCopyValue(value)
	}
	if value, ok := svc.Extensions["x-nomad-restart"]; ok {
		if restart, ok := asMap(value); ok {
			applyNomadRestartExtensions(service, restart)
		}
		if service.Extensions == nil {
			service.Extensions = map[string]interface{}{}
		}
		service.Extensions["x-nomad-restart"] = deepCopyValue(value)
	}
	if svc.Extends != nil {
		ensureComposeCompat(service).Extends = extendsConfigToMap(svc.Extends)
	}
	if svc.Platform != "" {
		ensureComposeCompat(service).Platform = svc.Platform
	}
	for key, value := range svc.Environment {
		if value != nil {
			service.Environment[key] = *value
		}
	}
	for _, envFile := range svc.EnvFiles {
		ref := EnvFileRef{
			Path:   envFile.Path,
			Format: envFile.Format,
		}
		if !bool(envFile.Required) {
			required := false
			ref.Required = &required
		}
		service.EnvFileRefs = append(service.EnvFileRefs, ref)
		service.EnvFile = append(service.EnvFile, envFile.Path)
	}
	for key, value := range svc.Extensions {
		service.Extensions[key] = value
		if alias := kubernetesComposeExtensionAlias(key); alias != "" {
			if _, exists := service.Extensions[alias]; !exists {
				service.Extensions[alias] = value
			}
		}
	}
	if value, ok := svc.Extensions[composeFailoverExtension]; ok {
		if failover, err := parseFailoverSpec(value); err == nil {
			service.Failover = failover
		}
	}
	if value, ok := svc.Extensions["x-env-sources"]; ok {
		if envSources, err := parseEnvSources(value); err == nil {
			service.EnvSources = envSources
		}
	}
	if value, ok := svc.Extensions["x-env-from"]; ok {
		if envFrom, err := parseEnvFromSources(value); err == nil {
			service.EnvFrom = envFrom
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-startup-probe"]; ok {
		if probe, ok := asMap(value); ok {
			service.StartupProbe = parseKubernetesProbe(probe)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-imagePullSecrets"]; ok {
		if secrets, err := toStringSlice(value); err == nil && len(secrets) > 0 {
			service.ImagePullSecrets = secrets
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-image-pull-secrets"]; ok {
		if secrets, err := toStringSlice(value); err == nil && len(secrets) > 0 {
			service.ImagePullSecrets = secrets
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-imagePullPolicy"]; ok {
		if policy := toString(value); policy != "" {
			service.ImagePullPolicy = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-image-pull-policy"]; ok {
		if policy := toString(value); policy != "" {
			service.ImagePullPolicy = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-hostAliases"]; ok {
		if aliases := kubernetesHostAliasesFromExtension(value); len(aliases) > 0 {
			service.HostAliases = aliases
			for _, alias := range aliases {
				for _, hostname := range alias.Hostnames {
					appendUniqueString(&service.ExtraHosts, hostname+"="+alias.IP)
				}
			}
		}
	}
	if value, ok := svc.Extensions["x-bolabaden-links"]; ok {
		if links, err := toStringSlice(value); err == nil && len(links) > 0 {
			service.Links = links
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-dnsPolicy"]; ok {
		if policy := toString(value); policy != "" {
			service.DNSPolicy = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-dns-policy"]; ok {
		if policy := toString(value); policy != "" {
			service.DNSPolicy = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-schedulerName"]; ok {
		if scheduler := toString(value); scheduler != "" {
			service.SchedulerName = scheduler
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-scheduler-name"]; ok {
		if scheduler := toString(value); scheduler != "" {
			service.SchedulerName = scheduler
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-terminationMessagePath"]; ok {
		if path := toString(value); path != "" {
			service.TerminationMessagePath = path
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-termination-message-path"]; ok {
		if path := toString(value); path != "" {
			service.TerminationMessagePath = path
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-terminationMessagePolicy"]; ok {
		if policy := toString(value); policy != "" {
			service.TerminationMessagePolicy = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-termination-message-policy"]; ok {
		if policy := toString(value); policy != "" {
			service.TerminationMessagePolicy = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-readiness-gates"]; ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); ok && len(gates) > 0 {
			service.ReadinessGates = gates
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-readinessGates"]; ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); ok && len(gates) > 0 {
			service.ReadinessGates = gates
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-allowPrivilegeEscalation"]; ok {
		if flag := boolPtrFromInterface(value); flag != nil {
			service.AllowPrivilegeEscalation = flag
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-allow-privilege-escalation"]; ok {
		if flag := boolPtrFromInterface(value); flag != nil {
			service.AllowPrivilegeEscalation = flag
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-procMount"]; ok {
		if text := toString(value); text != "" {
			service.ProcMount = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-proc-mount"]; ok {
		if text := toString(value); text != "" {
			service.ProcMount = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-init-containers"]; ok {
		if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
			service.InitContainers = containers
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-initContainers"]; ok {
		if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
			service.InitContainers = containers
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-resource-claims"]; ok {
		if claims, ok := kubernetesMapSliceFromExtension(value); ok && len(claims) > 0 {
			service.ResourceClaims = claims
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-ephemeral-containers"]; ok {
		if containers, ok := kubernetesMapSliceFromExtension(value); ok && len(containers) > 0 {
			service.EphemeralContainers = containers
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-scheduling-gates"]; ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); ok && len(gates) > 0 {
			service.SchedulingGates = gates
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-schedulingGates"]; ok {
		if gates, ok := kubernetesReadinessGatesFromExtension(value); ok && len(gates) > 0 {
			service.SchedulingGates = gates
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-hostUsers"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostUsers = boolPtr(flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-host-users"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostUsers = boolPtr(flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-group"]; ok {
		if text := toString(value); text != "" {
			service.Group = text
			service.Extensions["kubernetes.group"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-hostNetwork"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostNetwork = flag
			service.HostNetworkSet = true
			service.Extensions["kubernetes.hostNetwork"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-host-network"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostNetwork = flag
			service.HostNetworkSet = true
			service.Extensions["kubernetes.hostNetwork"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-hostPID"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostPID = boolPtr(flag)
			service.Extensions["kubernetes.hostPID"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-host-pid"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostPID = boolPtr(flag)
			service.Extensions["kubernetes.hostPID"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-hostIPC"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostIPC = boolPtr(flag)
			service.Extensions["kubernetes.hostIPC"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-host-ipc"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.HostIPC = boolPtr(flag)
			service.Extensions["kubernetes.hostIPC"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-pidMode"]; ok {
		if mode := toString(value); mode != "" {
			service.PIDMode = mode
			if strings.EqualFold(mode, "host") {
				service.HostPID = boolPtr(true)
			}
			service.Extensions["kubernetes.pidMode"] = mode
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-ipcMode"]; ok {
		if mode := toString(value); mode != "" {
			service.IPCMode = mode
			if strings.EqualFold(mode, "host") {
				service.HostIPC = boolPtr(true)
			}
			service.Extensions["kubernetes.ipcMode"] = mode
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-priorityClassName"]; ok {
		if text := toString(value); text != "" {
			service.PriorityClassName = text
			service.Extensions["kubernetes.priorityClassName"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-priority-class-name"]; ok {
		if text := toString(value); text != "" {
			service.PriorityClassName = text
			service.Extensions["kubernetes.priorityClassName"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-runtimeClassName"]; ok {
		if text := toString(value); text != "" {
			service.RuntimeClassName = text
			service.Extensions["kubernetes.runtimeClassName"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-runtime-class-name"]; ok {
		if text := toString(value); text != "" {
			service.RuntimeClassName = text
			service.Extensions["kubernetes.runtimeClassName"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-nodeName"]; ok {
		if text := toString(value); text != "" {
			service.NodeName = text
			service.Extensions["kubernetes.nodeName"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-node-name"]; ok {
		if text := toString(value); text != "" {
			service.NodeName = text
			service.Extensions["kubernetes.nodeName"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-subdomain"]; ok {
		if text := toString(value); text != "" {
			service.Subdomain = text
			service.Extensions["kubernetes.subdomain"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-sub-domain"]; ok {
		if text := toString(value); text != "" {
			service.Subdomain = text
			service.Extensions["kubernetes.subdomain"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-os"]; ok {
		if text := toString(value); text != "" {
			service.OSName = text
			service.Extensions["kubernetes.os"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-setHostnameAsFQDN"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.SetHostnameAsFQDN = boolPtr(flag)
			service.Extensions["kubernetes.setHostnameAsFQDN"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-set-hostname-as-fqdn"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.SetHostnameAsFQDN = boolPtr(flag)
			service.Extensions["kubernetes.setHostnameAsFQDN"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-shareProcessNamespace"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.ShareProcessNamespace = boolPtr(flag)
			service.Extensions["kubernetes.shareProcessNamespace"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-share-process-namespace"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.ShareProcessNamespace = boolPtr(flag)
			service.Extensions["kubernetes.shareProcessNamespace"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-enableServiceLinks"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.EnableServiceLinks = boolPtr(flag)
			service.Extensions["kubernetes.enableServiceLinks"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-enable-service-links"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.EnableServiceLinks = boolPtr(flag)
			service.Extensions["kubernetes.enableServiceLinks"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-serviceAccountName"]; ok {
		if text := toString(value); text != "" {
			service.ServiceAccountName = text
			service.Extensions["kubernetes.serviceAccountName"] = text
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-automountServiceAccountToken"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.AutomountServiceAccountToken = boolPtr(flag)
			service.Extensions["kubernetes.automountServiceAccountToken"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-automount-service-account-token"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.AutomountServiceAccountToken = boolPtr(flag)
			service.Extensions["kubernetes.automountServiceAccountToken"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-fsGroup"]; ok {
		if fsGroup := int64(toInt(value)); fsGroup > 0 {
			service.FSGroup = &fsGroup
			service.Extensions["kubernetes.fsGroup"] = fmt.Sprintf("%d", fsGroup)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-fs-group"]; ok {
		if fsGroup := int64(toInt(value)); fsGroup > 0 {
			service.FSGroup = &fsGroup
			service.Extensions["kubernetes.fsGroup"] = fmt.Sprintf("%d", fsGroup)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-seLinuxOptions"]; ok {
		if options, ok := asMap(value); ok {
			service.SELinuxOptions = parseKubernetesSELinuxOptions(options)
			if service.SELinuxOptions != nil {
				service.Extensions["kubernetes.seLinuxOptions"] = cloneMap(options)
			}
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-se-linux-options"]; ok {
		if options, ok := asMap(value); ok {
			service.SELinuxOptions = parseKubernetesSELinuxOptions(options)
			if service.SELinuxOptions != nil {
				service.Extensions["kubernetes.seLinuxOptions"] = cloneMap(options)
			}
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-windowsOptions"]; ok {
		if options, ok := asMap(value); ok {
			service.WindowsOptions = parseKubernetesWindowsSecurityContextOptions(options)
			if service.WindowsOptions != nil {
				service.Extensions["kubernetes.windowsOptions"] = cloneMap(options)
			}
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-windows-options"]; ok {
		if options, ok := asMap(value); ok {
			service.WindowsOptions = parseKubernetesWindowsSecurityContextOptions(options)
			if service.WindowsOptions != nil {
				service.Extensions["kubernetes.windowsOptions"] = cloneMap(options)
			}
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-fsGroupChangePolicy"]; ok {
		if policy := toString(value); policy != "" {
			service.FSGroupChangePolicy = policy
			service.Extensions["kubernetes.fsGroupChangePolicy"] = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-fs-group-change-policy"]; ok {
		if policy := toString(value); policy != "" {
			service.FSGroupChangePolicy = policy
			service.Extensions["kubernetes.fsGroupChangePolicy"] = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-runAsNonRoot"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.RunAsNonRoot = boolPtr(flag)
			service.Extensions["kubernetes.runAsNonRoot"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-run-as-non-root"]; ok {
		if flag := strings.EqualFold(toString(value), "true"); toString(value) != "" {
			service.RunAsNonRoot = boolPtr(flag)
			service.Extensions["kubernetes.runAsNonRoot"] = fmt.Sprintf("%t", flag)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-supplementalGroups"]; ok {
		if groups, err := toInt64Slice(value); err == nil && len(groups) > 0 {
			service.SupplementalGroups = groups
			service.Extensions["kubernetes.supplementalGroups"] = append([]int64{}, groups...)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-supplementalGroupsPolicy"]; ok {
		if policy := toString(value); policy != "" {
			service.SupplementalGroupsPolicy = policy
			service.Extensions["kubernetes.supplementalGroupsPolicy"] = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-supplemental-groups-policy"]; ok {
		if policy := toString(value); policy != "" {
			service.SupplementalGroupsPolicy = policy
			service.Extensions["kubernetes.supplementalGroupsPolicy"] = policy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-activeDeadlineSeconds"]; ok {
		if seconds := int64(toInt(value)); seconds > 0 {
			service.ActiveDeadlineSeconds = &seconds
			service.Extensions["kubernetes.activeDeadlineSeconds"] = fmt.Sprintf("%d", seconds)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-active-deadline-seconds"]; ok {
		if seconds := int64(toInt(value)); seconds > 0 {
			service.ActiveDeadlineSeconds = &seconds
			service.Extensions["kubernetes.activeDeadlineSeconds"] = fmt.Sprintf("%d", seconds)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-tolerations"]; ok {
		if tolerations := parseKubernetesTolerationsExtension(value); len(tolerations) > 0 {
			service.Tolerations = tolerations
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-restartPolicy"]; ok {
		if restartPolicy := toString(value); restartPolicy != "" {
			service.PodRestartPolicy = restartPolicy
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-seccomp-profile"]; ok {
		if profile, ok := asMap(value); ok {
			service.SeccompProfile = parseKubernetesSeccompProfile(profile)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-seccompProfile"]; ok {
		if profile, ok := asMap(value); ok {
			service.SeccompProfile = parseKubernetesSeccompProfile(profile)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-affinity"]; ok {
		if affinity, ok := asMap(value); ok {
			service.Affinity = copyStringInterfaceMap(affinity)
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-topology-spread-constraints"]; ok {
		if constraints, ok := kubernetesTopologySpreadConstraintsFromExtension(value); ok && len(constraints) > 0 {
			service.TopologySpreadConstraints = constraints
		}
	}
	if value, ok := svc.Extensions["x-kubernetes-topologySpreadConstraints"]; ok {
		if constraints, ok := kubernetesTopologySpreadConstraintsFromExtension(value); ok && len(constraints) > 0 {
			service.TopologySpreadConstraints = constraints
		}
	}
	for _, port := range svc.Ports {
		service.Ports = append(service.Ports, PortMapping{
			Name:          port.Name,
			HostIP:        port.HostIP,
			HostPort:      port.Published,
			ContainerPort: fmt.Sprintf("%d", port.Target),
			NodePort:      toString(port.Extensions["x-kubernetes-node-port"]),
			Protocol:      port.Protocol,
			AppProtocol:   port.AppProtocol,
			Mode:          port.Mode,
			Extensions:    composePortExtensions(port.Extensions),
		})
	}
	service.Networks, service.NetworkAttachments = networkAttachmentsFromCompose(svc.Networks)
	if len(service.NetworkAttachments) == 0 {
		service.NetworkAttachments = nil
	}
	for _, volume := range svc.Volumes {
		mount := VolumeMount{
			Source:           volume.Source,
			Target:           volume.Target,
			Type:             volume.Type,
			ReadOnly:         volume.ReadOnly,
			Consistency:      volume.Consistency,
			Extensions:       map[string]interface{}{},
			BindExtensions:   map[string]interface{}{},
			VolumeExtensions: map[string]interface{}{},
			TmpfsExtensions:  map[string]interface{}{},
			ImageExtensions:  map[string]interface{}{},
		}
		if volume.Bind != nil {
			mount.Mode = volume.Bind.SELinux
			mount.Propagation = volume.Bind.Propagation
			createHostPath := bool(volume.Bind.CreateHostPath)
			if !createHostPath {
				mount.CreateHostPath = &createHostPath
			}
			if volume.Bind.Recursive != "" {
				mount.Options = ensureStringMap(mount.Options)
				mount.Options["recursive"] = volume.Bind.Recursive
			}
			for key, value := range volume.Bind.Extensions {
				mount.BindExtensions[key] = value
			}
		}
		if volume.Volume != nil {
			mount.NoCopy = volume.Volume.NoCopy
			mount.SubPath = volume.Volume.Subpath
			if len(volume.Volume.Labels) > 0 {
				mount.VolumeLabels = map[string]string{}
				for key, value := range volume.Volume.Labels {
					mount.VolumeLabels[key] = value
				}
			}
			if volume.Volume.Subpath != "" {
				mount.Options = ensureStringMap(mount.Options)
				mount.Options["subpath"] = volume.Volume.Subpath
			}
			for key, value := range volume.Volume.Extensions {
				mount.VolumeExtensions[key] = value
			}
		}
		if volume.Tmpfs != nil && volume.Tmpfs.Size != 0 {
			mount.Options = ensureStringMap(mount.Options)
			mount.Options["size"] = fmt.Sprintf("%d", volume.Tmpfs.Size)
			if volume.Tmpfs.Mode != 0 {
				mount.TmpfsMode = fmt.Sprintf("%d", volume.Tmpfs.Mode)
			}
			for key, value := range volume.Tmpfs.Extensions {
				mount.TmpfsExtensions[key] = value
			}
		}
		if volume.Image != nil {
			mount.ImageSubpath = volume.Image.SubPath
			for key, value := range volume.Image.Extensions {
				mount.ImageExtensions[key] = value
			}
		}
		for key, value := range volume.Extensions {
			mount.Extensions[key] = value
		}
		service.Volumes = append(service.Volumes, mount)
	}
	for _, config := range svc.Configs {
		ref := fileRefFromCompose(composetypes.FileReferenceConfig(config))
		if ref.Source != "" {
			service.Configs = append(service.Configs, ref)
		}
	}
	for _, secret := range svc.Secrets {
		ref := fileRefFromCompose(composetypes.FileReferenceConfig(secret))
		if ref.Source != "" {
			service.Secrets = append(service.Secrets, ref)
		}
	}
	for name, dep := range svc.DependsOn {
		dependency := dependencyFromCompose(name, dep)
		service.Dependencies = appendUniqueDependency(service.Dependencies, dependency)
		service.DependsOn = appendUniqueName(service.DependsOn, name)
	}
	service.Links = append(service.Links, svc.Links...)
	for key, value := range svc.Labels {
		service.Labels[key] = value
	}
	service.HealthCheck = healthCheckFromCompose(svc.HealthCheck)
	service.Deploy = deployFromCompose(svc.Deploy)
	if svc.Scale != nil {
		service.Replicas = *svc.Scale
	} else if service.Deploy != nil && service.Deploy.Replicas > 0 {
		service.Replicas = service.Deploy.Replicas
	}
	hydrateServiceFailoverFromExtensions(service)
	return service
}

func hydrateServiceFailoverFromExtensions(service *Service) {
	if service == nil || service.Failover != nil {
		return
	}
	for _, candidate := range failoverExtensionCandidates(service) {
		if failover, err := parseFailoverSpec(candidate); err == nil && failover != nil {
			service.Failover = failover
			return
		}
		if raw, err := json.Marshal(candidate); err == nil {
			if failover, err := parseFailoverSpecMapBytes(raw); err == nil && failover != nil {
				service.Failover = failover
				return
			}
		}
	}
}

func failoverExtensionCandidates(service *Service) []interface{} {
	if service == nil {
		return nil
	}
	candidates := make([]interface{}, 0, 6)
	appendCandidate := func(value interface{}) {
		if value != nil {
			candidates = append(candidates, value)
		}
	}
	if value, ok := service.Extensions[composeFailoverExtension]; ok {
		appendCandidate(value)
	}
	if value, ok := service.Extensions["failover"]; ok {
		appendCandidate(value)
	}
	if compat := service.ComposeCompat; compat != nil {
		if value, ok := compat.Extensions[composeFailoverExtension]; ok {
			appendCandidate(value)
		}
		if value, ok := compat.Extensions["failover"]; ok {
			appendCandidate(value)
		}
		if value, ok := compat.Extensions["extensions"]; ok {
			if mapped, ok := asMap(value); ok {
				if nested, ok := mapped[composeFailoverExtension]; ok {
					appendCandidate(nested)
				}
				if nested, ok := mapped["failover"]; ok {
					appendCandidate(nested)
				}
			}
		}
	}
	if value, ok := service.Extensions[composeCompatExtension]; ok {
		if mapped, ok := asMap(value); ok {
			if nested, ok := mapped[composeFailoverExtension]; ok {
				appendCandidate(nested)
			}
			if nested, ok := mapped["failover"]; ok {
				appendCandidate(nested)
			}
			if extensions, ok := asMap(mapped["extensions"]); ok {
				if nested, ok := extensions[composeFailoverExtension]; ok {
					appendCandidate(nested)
				}
				if nested, ok := extensions["failover"]; ok {
					appendCandidate(nested)
				}
			}
		}
	}
	return candidates
}

func networkAttachmentsFromCompose(networks map[string]*composetypes.ServiceNetworkConfig) ([]string, map[string]*NetworkAttachment) {
	if len(networks) == 0 {
		return nil, nil
	}
	names := make([]string, 0, len(networks))
	attachments := map[string]*NetworkAttachment{}
	for name, network := range networks {
		names = append(names, name)
		if network == nil {
			continue
		}
		attachment := &NetworkAttachment{
			Name:          name,
			Aliases:       append([]string{}, network.Aliases...),
			InterfaceName: network.InterfaceName,
			IPv4Address:   network.Ipv4Address,
			IPv6Address:   network.Ipv6Address,
			LinkLocalIPs:  append([]string{}, network.LinkLocalIPs...),
			MacAddress:    network.MacAddress,
			DriverOpts:    copyStringMap(map[string]string(network.DriverOpts)),
			GWPriority:    network.GatewayPriority,
			Priority:      network.Priority,
			Extensions:    copyStringInterfaceMap(map[string]interface{}(network.Extensions)),
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
		return names, nil
	}
	return names, attachments
}

func composeCompatFromServiceConfig(svc composetypes.ServiceConfig) *ComposeCompat {
	compat := &ComposeCompat{
		Attach:            composeBoolPtr(svc.Attach),
		Annotations:       copyStringMap(map[string]string(svc.Annotations)),
		Extends:           extendsConfigToMap(svc.Extends),
		Platform:          svc.Platform,
		MemSwappiness:     composeUnitBytesString(svc.MemSwappiness),
		MacAddress:        svc.MacAddress,
		DomainName:        svc.DomainName,
		CgroupParent:      svc.CgroupParent,
		Cgroup:            svc.Cgroup,
		CPUCount:          svc.CPUCount,
		CPUPercent:        svc.CPUPercent,
		CPUPeriod:         svc.CPUPeriod,
		CPURTPeriod:       svc.CPURTPeriod,
		CPURTRuntime:      svc.CPURTRuntime,
		CPUSet:            svc.CPUSet,
		DeviceCgroupRules: append([]string{}, svc.DeviceCgroupRules...),
		NetworkMode:       firstNonEmpty(svc.NetworkMode, svc.Net),
		OomKillDisable:    svc.OomKillDisable,
		OomScoreAdj:       svc.OomScoreAdj,
		Scale:             svc.Scale,
		ExternalLinks:     append([]string{}, svc.ExternalLinks...),
		LabelFiles:        append([]string{}, svc.LabelFiles...),
		StorageOpt:        copyStringMap(map[string]string(svc.StorageOpt)),
		UseAPISocket:      svc.UseAPISocket,
		Isolation:         svc.Isolation,
		Tmpfs:             append([]string{}, svc.Tmpfs...),
		Uts:               svc.Uts,
		VolumeDriver:      svc.VolumeDriver,
		VolumesFrom:       append([]string{}, svc.VolumesFrom...),
		Models:            composeModelsToCompat(svc.Models),
		Extensions:        map[string]interface{}{},
	}
	if len(svc.Gpus) > 0 {
		compat.Gpus = make([]map[string]interface{}, 0, len(svc.Gpus))
		for _, gpu := range svc.Gpus {
			compat.Gpus = append(compat.Gpus, gpuRequestToMap(gpu))
		}
	}
	if svc.BlkioConfig != nil {
		compat.BlkioConfig = blkioConfigToMap(svc.BlkioConfig)
	}
	if svc.CredentialSpec != nil {
		compat.CredentialSpec = credentialSpecToMap(svc.CredentialSpec)
	}
	for key, value := range svc.Extensions {
		compat.Extensions[key] = deepCopyValue(value)
	}
	if len(compat.Extensions) == 0 {
		compat.Extensions = nil
	}
	if isEmptyComposeCompat(compat) {
		return nil
	}
	return compat
}

func composeCompatFromExtensionMap(mapped map[string]interface{}) *ComposeCompat {
	if len(mapped) == 0 {
		return nil
	}
	compat := &ComposeCompat{
		Annotations:    map[string]string{},
		BlkioConfig:    map[string]interface{}{},
		CredentialSpec: map[string]interface{}{},
		Provider:       map[string]interface{}{},
		Extends:        map[string]interface{}{},
		StorageOpt:     map[string]string{},
		Models:         map[string]map[string]interface{}{},
		Extensions:     map[string]interface{}{},
	}
	for key, value := range mapped {
		switch key {
		case "attach":
			flag := toBool(value)
			compat.Attach = &flag
		case "annotations":
			compat.Annotations = toStringMapLoose(value)
		case "blkio_config":
			if mappedValue, ok := asMap(value); ok {
				compat.BlkioConfig = copyStringInterfaceMap(mappedValue)
			}
		case "credential_spec":
			if mappedValue, ok := asMap(value); ok {
				compat.CredentialSpec = copyStringInterfaceMap(mappedValue)
			}
		case "provider":
			if mappedValue, ok := asMap(value); ok {
				compat.Provider = copyStringInterfaceMap(mappedValue)
			}
		case "extends":
			if mappedValue, ok := asMap(value); ok {
				compat.Extends = copyStringInterfaceMap(mappedValue)
			}
		case "platform":
			compat.Platform = toString(value)
		case "pull_refresh_after":
			compat.PullRefreshAfter = toString(value)
		case "mem_swappiness":
			compat.MemSwappiness = toString(value)
		case "mac_address":
			compat.MacAddress = toString(value)
		case "domain_name", "domainname":
			compat.DomainName = toString(value)
		case "cgroup_parent":
			compat.CgroupParent = toString(value)
		case "cgroup":
			compat.Cgroup = toString(value)
		case "cpu_count":
			compat.CPUCount = int64(toInt(value))
			compat.CPUCountSet = true
		case "cpu_percent":
			if f, err := strconv.ParseFloat(fmt.Sprint(value), 32); err == nil {
				compat.CPUPercent = float32(f)
				compat.CPUPercentSet = true
			}
		case "cpu_period":
			compat.CPUPeriod = int64(toInt(value))
			compat.CPUPeriodSet = true
		case "cpu_rt_period":
			compat.CPURTPeriod = int64(toInt(value))
			compat.CPURTPeriodSet = true
		case "cpu_rt_runtime":
			compat.CPURTRuntime = int64(toInt(value))
			compat.CPURTRuntimeSet = true
		case "cpuset":
			compat.CPUSet = toString(value)
		case "device_cgroup_rules":
			if rules, err := toStringSlice(value); err == nil {
				compat.DeviceCgroupRules = rules
			}
		case "gpus":
			if list, ok := value.([]interface{}); ok {
				compat.Gpus = make([]map[string]interface{}, 0, len(list))
				for _, item := range list {
					if mapped, ok := asMap(item); ok {
						compat.Gpus = append(compat.Gpus, copyStringInterfaceMap(mapped))
					}
				}
			}
		case "network_mode":
			compat.NetworkMode = toString(value)
		case "oom_kill_disable":
			compat.OomKillDisable = toBool(value)
			compat.OomKillDisableSet = true
		case "oom_score_adj":
			compat.OomScoreAdj = int64(toInt(value))
			compat.OomScoreAdjSet = true
		case "scale":
			scale := toInt(value)
			compat.Scale = &scale
		case "models":
			if models, ok := value.(map[string]interface{}); ok {
				compat.Models = map[string]map[string]interface{}{}
				for name, rawModel := range models {
					if mapped, ok := asMap(rawModel); ok {
						if model := composeModelFromExtensionMap(mapped); model != nil {
							compat.Models[name] = map[string]interface{}{
								"name":          model.Name,
								"model":         model.Model,
								"context_size":  model.ContextSize,
								"runtime_flags": append([]string{}, model.RuntimeFlags...),
							}
							if len(model.Extensions) > 0 {
								for nestedKey, nestedValue := range model.Extensions {
									compat.Models[name][nestedKey] = deepCopyValue(nestedValue)
								}
							}
						}
					}
				}
			}
		case "external_links":
			if values, err := toStringSlice(value); err == nil {
				compat.ExternalLinks = values
			}
		case "label_files":
			if values, err := toStringSlice(value); err == nil {
				compat.LabelFiles = values
			}
		case "storage_opt":
			compat.StorageOpt = toStringMapLoose(value)
		case "use_api_socket":
			compat.UseAPISocket = toBool(value)
			compat.UseAPISocketSet = true
		case "isolation":
			compat.Isolation = toString(value)
		case "tmpfs":
			if values, err := toStringSlice(value); err == nil {
				compat.Tmpfs = values
			}
		case "uts":
			compat.Uts = toString(value)
		case "volume_driver":
			compat.VolumeDriver = toString(value)
		case "volumes_from":
			if values, err := toStringSlice(value); err == nil {
				compat.VolumesFrom = values
			}
		case "extensions":
			if nested, ok := asMap(value); ok {
				for nestedKey, nestedValue := range nested {
					compat.Extensions[nestedKey] = deepCopyValue(nestedValue)
				}
			}
		default:
			compat.Extensions[key] = deepCopyValue(value)
		}
	}
	if len(compat.Annotations) == 0 {
		compat.Annotations = nil
	}
	if len(compat.BlkioConfig) == 0 {
		compat.BlkioConfig = nil
	}
	if len(compat.CredentialSpec) == 0 {
		compat.CredentialSpec = nil
	}
	if len(compat.Provider) == 0 {
		compat.Provider = nil
	}
	if len(compat.Extends) == 0 {
		compat.Extends = nil
	}
	if len(compat.StorageOpt) == 0 {
		compat.StorageOpt = nil
	}
	if len(compat.Models) == 0 {
		compat.Models = nil
	}
	if len(compat.Extensions) == 0 {
		compat.Extensions = nil
	}
	if isEmptyComposeCompat(compat) {
		return nil
	}
	return compat
}

func mergeComposeCompat(dst, src *ComposeCompat) {
	if dst == nil || src == nil {
		return
	}
	if dst.Attach == nil {
		dst.Attach = cloneBoolPtr(src.Attach)
	}
	if len(dst.Annotations) == 0 && len(src.Annotations) > 0 {
		dst.Annotations = copyStringMap(src.Annotations)
	}
	if len(dst.BlkioConfig) == 0 && len(src.BlkioConfig) > 0 {
		dst.BlkioConfig = copyStringInterfaceMap(src.BlkioConfig)
	}
	if len(dst.CredentialSpec) == 0 && len(src.CredentialSpec) > 0 {
		dst.CredentialSpec = copyStringInterfaceMap(src.CredentialSpec)
	}
	if len(dst.Provider) == 0 && len(src.Provider) > 0 {
		dst.Provider = copyStringInterfaceMap(src.Provider)
	}
	if len(dst.Extends) == 0 && len(src.Extends) > 0 {
		dst.Extends = copyStringInterfaceMap(src.Extends)
	}
	if dst.Platform == "" {
		dst.Platform = src.Platform
	}
	if dst.PullRefreshAfter == "" {
		dst.PullRefreshAfter = src.PullRefreshAfter
	}
	if dst.MemSwappiness == "" {
		dst.MemSwappiness = src.MemSwappiness
	}
	if dst.MacAddress == "" {
		dst.MacAddress = src.MacAddress
	}
	if dst.DomainName == "" {
		dst.DomainName = src.DomainName
	}
	if dst.CgroupParent == "" {
		dst.CgroupParent = src.CgroupParent
	}
	if dst.Cgroup == "" {
		dst.Cgroup = src.Cgroup
	}
	if !dst.CPUCountSet {
		if dst.CPUCount == 0 {
			dst.CPUCount = src.CPUCount
			dst.CPUCountSet = src.CPUCountSet
		} else {
			dst.CPUCountSet = true
		}
	}
	if !dst.CPUPercentSet {
		if dst.CPUPercent == 0 {
			dst.CPUPercent = src.CPUPercent
			dst.CPUPercentSet = src.CPUPercentSet
		} else {
			dst.CPUPercentSet = true
		}
	}
	if !dst.CPUPeriodSet {
		if dst.CPUPeriod == 0 {
			dst.CPUPeriod = src.CPUPeriod
			dst.CPUPeriodSet = src.CPUPeriodSet
		} else {
			dst.CPUPeriodSet = true
		}
	}
	if !dst.CPURTPeriodSet {
		if dst.CPURTPeriod == 0 {
			dst.CPURTPeriod = src.CPURTPeriod
			dst.CPURTPeriodSet = src.CPURTPeriodSet
		} else {
			dst.CPURTPeriodSet = true
		}
	}
	if !dst.CPURTRuntimeSet {
		if dst.CPURTRuntime == 0 {
			dst.CPURTRuntime = src.CPURTRuntime
			dst.CPURTRuntimeSet = src.CPURTRuntimeSet
		} else {
			dst.CPURTRuntimeSet = true
		}
	}
	if dst.CPUSet == "" {
		dst.CPUSet = src.CPUSet
	}
	if len(dst.DeviceCgroupRules) == 0 && len(src.DeviceCgroupRules) > 0 {
		dst.DeviceCgroupRules = append([]string{}, src.DeviceCgroupRules...)
	}
	if len(dst.Gpus) == 0 && len(src.Gpus) > 0 {
		dst.Gpus = cloneMapSlice(src.Gpus)
	}
	if dst.NetworkMode == "" {
		dst.NetworkMode = src.NetworkMode
	}
	if !dst.OomKillDisableSet {
		if !dst.OomKillDisable {
			dst.OomKillDisable = src.OomKillDisable
			dst.OomKillDisableSet = src.OomKillDisableSet
		} else {
			dst.OomKillDisableSet = true
		}
	}
	if !dst.OomScoreAdjSet {
		if dst.OomScoreAdj == 0 {
			dst.OomScoreAdj = src.OomScoreAdj
			dst.OomScoreAdjSet = src.OomScoreAdjSet
		} else {
			dst.OomScoreAdjSet = true
		}
	}
	if dst.Scale == nil {
		dst.Scale = cloneIntPtr(src.Scale)
	}
	if len(dst.Models) == 0 && len(src.Models) > 0 {
		dst.Models = cloneComposeModelMap(src.Models)
	}
	if len(dst.ExternalLinks) == 0 && len(src.ExternalLinks) > 0 {
		dst.ExternalLinks = append([]string{}, src.ExternalLinks...)
	}
	if len(dst.LabelFiles) == 0 && len(src.LabelFiles) > 0 {
		dst.LabelFiles = append([]string{}, src.LabelFiles...)
	}
	if len(dst.StorageOpt) == 0 && len(src.StorageOpt) > 0 {
		dst.StorageOpt = copyStringMap(src.StorageOpt)
	}
	if !dst.UseAPISocketSet {
		if !dst.UseAPISocket {
			dst.UseAPISocket = src.UseAPISocket
			dst.UseAPISocketSet = src.UseAPISocketSet
		} else {
			dst.UseAPISocketSet = true
		}
	}
	if dst.Isolation == "" {
		dst.Isolation = src.Isolation
	}
	if len(dst.Tmpfs) == 0 && len(src.Tmpfs) > 0 {
		dst.Tmpfs = append([]string{}, src.Tmpfs...)
	}
	if dst.Uts == "" {
		dst.Uts = src.Uts
	}
	if dst.VolumeDriver == "" {
		dst.VolumeDriver = src.VolumeDriver
	}
	if len(dst.VolumesFrom) == 0 && len(src.VolumesFrom) > 0 {
		dst.VolumesFrom = append([]string{}, src.VolumesFrom...)
	}
	if len(dst.Extensions) == 0 && len(src.Extensions) > 0 {
		dst.Extensions = copyStringInterfaceMap(src.Extensions)
	} else if len(src.Extensions) > 0 {
		if dst.Extensions == nil {
			dst.Extensions = map[string]interface{}{}
		}
		for key, value := range src.Extensions {
			if _, ok := dst.Extensions[key]; !ok {
				dst.Extensions[key] = deepCopyValue(value)
			}
		}
	}
	if len(dst.Annotations) == 0 && len(src.Annotations) > 0 {
		dst.Annotations = copyStringMap(src.Annotations)
	}
}

func composeUnitBytesString(value composetypes.UnitBytes) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

func extendsConfigToMap(extends *composetypes.ExtendsConfig) map[string]interface{} {
	if extends == nil || (extends.File == "" && extends.Service == "") {
		return nil
	}
	result := map[string]interface{}{}
	if extends.Service != "" {
		result["service"] = extends.Service
	}
	if extends.File != "" {
		result["file"] = extends.File
	}
	return result
}

func kubernetesComposeExtensionAlias(key string) string {
	if !strings.HasPrefix(key, "x-kubernetes-") {
		if !strings.HasPrefix(key, "x-nomad-") {
			return ""
		}
		return "nomad." + strings.ReplaceAll(strings.TrimPrefix(key, "x-nomad-"), "-", ".")
	}
	switch key {
	case "x-kubernetes-horizontal-pod-autoscalers":
		return "kubernetes.horizontalPodAutoscalers"
	case "x-kubernetes-pod-disruption-budgets":
		return "kubernetes.podDisruptionBudgets"
	case "x-kubernetes-topology-spread-constraints":
		return "kubernetes.topologySpreadConstraints"
	case "x-kubernetes-host-aliases":
		return "kubernetes.hostAliases"
	case "x-kubernetes-job-pod-failure-policy":
		return "kubernetes.job.podFailurePolicy"
	case "x-kubernetes-job-podFailurePolicy":
		return "kubernetes.job.podFailurePolicy"
	case "x-kubernetes-job-success-policy":
		return "kubernetes.job.successPolicy"
	case "x-kubernetes-job-successPolicy":
		return "kubernetes.job.successPolicy"
	case "x-kubernetes-priority-class-name":
		return "kubernetes.priorityClassName"
	case "x-kubernetes-runtime-class-name":
		return "kubernetes.runtimeClassName"
	case "x-kubernetes-node-name":
		return "kubernetes.nodeName"
	case "x-kubernetes-sub-domain":
		return "kubernetes.subdomain"
	case "x-kubernetes-service-account-name":
		return "kubernetes.serviceAccountName"
	case "x-kubernetes-fs-group":
		return "kubernetes.fsGroup"
	case "x-kubernetes-run-as-non-root":
		return "kubernetes.runAsNonRoot"
	case "x-kubernetes-active-deadline-seconds":
		return "kubernetes.activeDeadlineSeconds"
	case "x-kubernetes-restart-policy":
		return "kubernetes.restartPolicy"
	case "x-kubernetes-set-hostname-as-fqdn":
		return "kubernetes.setHostnameAsFQDN"
	case "x-kubernetes-enable-service-links":
		return "kubernetes.enableServiceLinks"
	case "x-kubernetes-host-users":
		return "kubernetes.hostUsers"
	case "x-kubernetes-termination-message-path":
		return "kubernetes.terminationMessagePath"
	case "x-kubernetes-termination-message-policy":
		return "kubernetes.terminationMessagePolicy"
	case "x-kubernetes-allow-privilege-escalation":
		return "kubernetes.allowPrivilegeEscalation"
	case "x-kubernetes-proc-mount":
		return "kubernetes.procMount"
	case "x-kubernetes-fs-group-change-policy":
		return "kubernetes.fsGroupChangePolicy"
	case "x-kubernetes-supplemental-groups-policy":
		return "kubernetes.supplementalGroupsPolicy"
	}
	return "kubernetes." + strings.ReplaceAll(strings.TrimPrefix(key, "x-kubernetes-"), "-", ".")
}

func composePortExtensions(extensions composetypes.Extensions) map[string]interface{} {
	if len(extensions) == 0 {
		return nil
	}
	result := map[string]interface{}{}
	for key, value := range extensions {
		if key == "x-kubernetes-node-port" {
			continue
		}
		result[key] = deepCopyValue(value)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func composeModelsToCompat(models map[string]*composetypes.ServiceModelConfig) map[string]map[string]interface{} {
	if len(models) == 0 {
		return nil
	}
	result := make(map[string]map[string]interface{}, len(models))
	for name, model := range models {
		if model == nil {
			continue
		}
		entry := map[string]interface{}{}
		if model.EndpointVariable != "" {
			entry["endpoint_var"] = model.EndpointVariable
		}
		if model.ModelVariable != "" {
			entry["model_var"] = model.ModelVariable
		}
		for key, value := range model.Extensions {
			entry[key] = deepCopyValue(value)
		}
		if len(entry) > 0 {
			result[name] = entry
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func composeModelFromCompose(name string, model composetypes.ModelConfig) *ComposeModel {
	item := &ComposeModel{
		Name:         name,
		Model:        model.Model,
		ContextSize:  model.ContextSize,
		RuntimeFlags: append([]string{}, model.RuntimeFlags...),
		Extensions:   map[string]interface{}{},
	}
	for key, value := range model.Extensions {
		item.Extensions[key] = deepCopyValue(value)
	}
	if len(item.Extensions) == 0 {
		item.Extensions = nil
	}
	if item.Name == "" && item.Model == "" && item.ContextSize == 0 && len(item.RuntimeFlags) == 0 && len(item.Extensions) == 0 {
		return nil
	}
	return item
}

func gpuRequestToMap(req composetypes.DeviceRequest) map[string]interface{} {
	result := map[string]interface{}{}
	if len(req.Capabilities) > 0 {
		result["capabilities"] = append([]string{}, req.Capabilities...)
	}
	if req.Driver != "" {
		result["driver"] = req.Driver
	}
	if req.Count != 0 {
		result["count"] = int64(req.Count)
	}
	if len(req.IDs) > 0 {
		result["device_ids"] = append([]string{}, req.IDs...)
	}
	if len(req.Options) > 0 {
		result["options"] = copyStringMap(map[string]string(req.Options))
	}
	return result
}

func blkioConfigToMap(cfg *composetypes.BlkioConfig) map[string]interface{} {
	if cfg == nil {
		return nil
	}
	result := map[string]interface{}{}
	if cfg.Weight != 0 {
		result["weight"] = cfg.Weight
	}
	if len(cfg.WeightDevice) > 0 {
		result["weight_device"] = blkioWeightDevicesToMaps(cfg.WeightDevice)
	}
	if len(cfg.DeviceReadBps) > 0 {
		result["device_read_bps"] = blkioThrottleDevicesToMaps(cfg.DeviceReadBps)
	}
	if len(cfg.DeviceReadIOps) > 0 {
		result["device_read_iops"] = blkioThrottleDevicesToMaps(cfg.DeviceReadIOps)
	}
	if len(cfg.DeviceWriteBps) > 0 {
		result["device_write_bps"] = blkioThrottleDevicesToMaps(cfg.DeviceWriteBps)
	}
	if len(cfg.DeviceWriteIOps) > 0 {
		result["device_write_iops"] = blkioThrottleDevicesToMaps(cfg.DeviceWriteIOps)
	}
	for key, value := range cfg.Extensions {
		result[composeApplicationExtensionKey(key)] = deepCopyValue(value)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func blkioWeightDevicesToMaps(devices []composetypes.WeightDevice) []interface{} {
	if len(devices) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(devices))
	for _, device := range devices {
		entry := map[string]interface{}{}
		if device.Path != "" {
			entry["path"] = device.Path
		}
		if device.Weight != 0 {
			entry["weight"] = device.Weight
		}
		for key, value := range device.Extensions {
			entry[composeApplicationExtensionKey(key)] = deepCopyValue(value)
		}
		result = append(result, entry)
	}
	return result
}

func blkioThrottleDevicesToMaps(devices []composetypes.ThrottleDevice) []interface{} {
	if len(devices) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(devices))
	for _, device := range devices {
		entry := map[string]interface{}{}
		if device.Path != "" {
			entry["path"] = device.Path
		}
		if device.Rate != 0 {
			entry["rate"] = int64(device.Rate)
		}
		for key, value := range device.Extensions {
			entry[composeApplicationExtensionKey(key)] = deepCopyValue(value)
		}
		result = append(result, entry)
	}
	return result
}

func credentialSpecToMap(cfg *composetypes.CredentialSpecConfig) map[string]interface{} {
	if cfg == nil {
		return nil
	}
	result := map[string]interface{}{}
	if cfg.Config != "" {
		result["config"] = cfg.Config
	}
	if cfg.File != "" {
		result["file"] = cfg.File
	}
	if cfg.Registry != "" {
		result["registry"] = cfg.Registry
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func cloneMultiOptions(input composetypes.MultiOptions) map[string][]string {
	if len(input) == 0 {
		return nil
	}
	result := map[string][]string{}
	for key, values := range input {
		if len(values) == 0 {
			continue
		}
		result[key] = append([]string{}, values...)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func developFromCompose(develop *composetypes.DevelopConfig) *DevelopConfig {
	if develop == nil {
		return nil
	}
	result := &DevelopConfig{
		Extensions: map[string]interface{}{},
	}
	for _, trigger := range develop.Watch {
		watch := DevelopWatch{
			Path:        trigger.Path,
			Action:      string(trigger.Action),
			Target:      trigger.Target,
			Exec:        developHookFromCompose(trigger.Exec),
			Include:     append([]string{}, trigger.Include...),
			Ignore:      append([]string{}, trigger.Ignore...),
			InitialSync: trigger.InitialSync,
			Extensions:  map[string]interface{}{},
		}
		for key, value := range trigger.Extensions {
			watch.Extensions[key] = deepCopyValue(value)
		}
		if len(watch.Extensions) == 0 {
			watch.Extensions = nil
		}
		result.Watch = append(result.Watch, watch)
	}
	for key, value := range develop.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	if len(result.Watch) == 0 && len(result.Extensions) == 0 {
		return nil
	}
	return result
}

func lifecycleFromCompose(preStart, postStart, preStop []composetypes.ServiceHook) *LifecycleHooks {
	lifecycle := &LifecycleHooks{
		PreStart:   serviceHooksFromCompose(preStart),
		PostStart:  serviceHooksFromCompose(postStart),
		PreStop:    serviceHooksFromCompose(preStop),
		Extensions: map[string]interface{}{},
	}
	if isEmptyLifecycleHooks(lifecycle) {
		return nil
	}
	if len(lifecycle.Extensions) == 0 {
		lifecycle.Extensions = nil
	}
	return lifecycle
}

func serviceHooksFromCompose(hooks []composetypes.ServiceHook) []ServiceHook {
	if len(hooks) == 0 {
		return nil
	}
	result := make([]ServiceHook, 0, len(hooks))
	for _, hook := range hooks {
		if converted := serviceHookFromCompose(hook); converted != nil {
			result = append(result, *converted)
		}
	}
	return result
}

func developHookFromCompose(hook composetypes.ServiceHook) *ServiceHook {
	return serviceHookFromCompose(hook)
}

func serviceHookFromCompose(hook composetypes.ServiceHook) *ServiceHook {
	if len(hook.Command) == 0 &&
		hook.Image == "" &&
		hook.User == "" &&
		!hook.Privileged &&
		hook.WorkingDir == "" &&
		len(hook.Environment) == 0 &&
		!hook.PerReplica &&
		len(hook.Extensions) == 0 {
		return nil
	}
	result := &ServiceHook{
		Command:     append([]string{}, hook.Command...),
		Image:       hook.Image,
		User:        hook.User,
		Privileged:  hook.Privileged,
		WorkingDir:  hook.WorkingDir,
		Environment: copyStringPtrMap(map[string]*string(hook.Environment)),
		PerReplica:  hook.PerReplica,
		Extensions:  map[string]interface{}{},
	}
	for key, value := range hook.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	return result
}

func deployFromCompose(deploy *composetypes.DeployConfig) *DeploySpec {
	if deploy == nil {
		return nil
	}
	result := &DeploySpec{
		Mode:         deploy.Mode,
		EndpointMode: deploy.EndpointMode,
		Labels:       map[string]string{},
		Extensions:   map[string]interface{}{},
	}
	if deploy.Replicas != nil {
		result.Replicas = *deploy.Replicas
	}
	if value, ok := deploy.Extensions[composeSwarmJobExtension]; ok {
		if job, err := parseSwarmJobSpec(value); err == nil {
			result.Job = job
		}
	}
	if value, ok := deploy.Extensions["x-nomad-update"]; ok {
		if update, err := parseUpdatePolicy(value); err == nil {
			result.UpdateConfig = update
		}
	}
	if value, ok := deploy.Extensions["x-nomad-migrate"]; ok {
		if migrate, err := parseMigratePolicy(value); err == nil {
			result.MigrateConfig = migrate
		}
	}
	if value, ok := deploy.Extensions["x-nomad-reschedule"]; ok {
		if reschedule, err := parseReschedulePolicy(value); err == nil {
			result.RescheduleConfig = reschedule
		}
	}
	for key, value := range deploy.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	for key, value := range deploy.Labels {
		result.Labels[key] = value
	}
	if len(result.Labels) == 0 {
		result.Labels = nil
	}
	result.Placement = placementFromCompose(deploy.Placement)
	result.Resources = resourcesFromCompose(deploy.Resources)
	if deploy.UpdateConfig != nil {
		result.UpdateConfig = updatePolicyFromCompose(deploy.UpdateConfig)
	}
	if deploy.RollbackConfig != nil {
		result.RollbackConfig = updatePolicyFromCompose(deploy.RollbackConfig)
	}
	if deploy.RestartPolicy != nil {
		result.RestartPolicy = restartPolicyFromCompose(deploy.RestartPolicy)
	}
	return result
}

func dependencyFromCompose(name string, dep composetypes.ServiceDependency) DependencySpec {
	result := DependencySpec{Name: name, Extensions: map[string]interface{}{}}
	if dep.Condition != "" && dep.Condition != composetypes.ServiceConditionStarted {
		result.Condition = dep.Condition
	}
	if dep.Restart {
		result.Restart = true
	}
	if !dep.Required {
		result.Required = boolPtr(false)
	}
	for key, value := range dep.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	return result
}

func placementFromCompose(placement composetypes.Placement) *PlacementSpec {
	if len(placement.Constraints) == 0 && len(placement.Preferences) == 0 && placement.MaxReplicas == 0 && len(placement.Extensions) == 0 {
		return nil
	}
	result := &PlacementSpec{
		Constraints:        append([]string{}, placement.Constraints...),
		MaxReplicasPerNode: int(placement.MaxReplicas),
		Extensions:         map[string]interface{}{},
	}
	for key, value := range placement.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	for _, pref := range placement.Preferences {
		extensions := map[string]interface{}{}
		for key, value := range pref.Extensions {
			extensions[key] = deepCopyValue(value)
		}
		if pref.Spread != "" || len(extensions) > 0 {
			if pref.Spread != "" {
				result.Preferences = append(result.Preferences, "spread="+pref.Spread)
			} else {
				result.Preferences = append(result.Preferences, "")
			}
			if len(extensions) == 0 {
				result.PreferenceExtensions = append(result.PreferenceExtensions, nil)
			} else {
				result.PreferenceExtensions = append(result.PreferenceExtensions, extensions)
			}
		}
	}
	if !hasPlacementPreferenceExtensions(result.PreferenceExtensions) {
		result.PreferenceExtensions = nil
	}
	return result
}

func resourcesFromCompose(resources composetypes.Resources) *ResourceSpec {
	result := &ResourceSpec{
		Extensions:            map[string]interface{}{},
		LimitExtensions:       map[string]interface{}{},
		ReservationExtensions: map[string]interface{}{},
	}
	for key, value := range resources.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if resources.Limits != nil {
		for key, value := range resources.Limits.Extensions {
			result.LimitExtensions[key] = deepCopyValue(value)
		}
		if resources.Limits.NanoCPUs != 0 {
			result.CPULimit = fmt.Sprintf("%g", resources.Limits.NanoCPUs.Value())
		}
		if resources.Limits.MemoryBytes != 0 {
			result.MemoryLimit = fmt.Sprintf("%d", resources.Limits.MemoryBytes)
		}
		if resources.Limits.Pids != 0 {
			result.PidsLimit = resources.Limits.Pids
		}
	}
	if resources.Reservations != nil {
		for key, value := range resources.Reservations.Extensions {
			result.ReservationExtensions[key] = deepCopyValue(value)
		}
		if resources.Reservations.NanoCPUs != 0 {
			result.CPUReservation = fmt.Sprintf("%g", resources.Reservations.NanoCPUs.Value())
		}
		if resources.Reservations.MemoryBytes != 0 {
			result.MemoryReservation = fmt.Sprintf("%d", resources.Reservations.MemoryBytes)
		}
		if resources.Reservations.Pids != 0 {
			result.PidsReservation = resources.Reservations.Pids
		}
		result.Devices = resourceDevicesFromCompose(resources.Reservations.Devices)
		result.GenericResources = genericResourcesFromCompose(resources.Reservations.GenericResources)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	if len(result.LimitExtensions) == 0 {
		result.LimitExtensions = nil
	}
	if len(result.ReservationExtensions) == 0 {
		result.ReservationExtensions = nil
	}
	if isEmptyResourceSpec(result) {
		return nil
	}
	return result
}

func resourceDevicesFromCompose(devices []composetypes.DeviceRequest) []ResourceDevice {
	if len(devices) == 0 {
		return nil
	}
	result := make([]ResourceDevice, 0, len(devices))
	for _, device := range devices {
		item := ResourceDevice{
			Capabilities: append([]string{}, device.Capabilities...),
			Driver:       device.Driver,
			DeviceIDs:    append([]string{}, device.IDs...),
		}
		if int64(device.Count) == -1 {
			item.Count = "all"
		} else if device.Count != 0 {
			item.Count = strconv.FormatInt(int64(device.Count), 10)
		}
		if len(device.Options) > 0 {
			item.Options = map[string]string{}
			for key, value := range device.Options {
				item.Options[key] = value
			}
		}
		result = append(result, item)
	}
	return result
}

func genericResourcesFromCompose(resources []composetypes.GenericResource) []GenericResource {
	if len(resources) == 0 {
		return nil
	}
	result := make([]GenericResource, 0, len(resources))
	for _, resource := range resources {
		if resource.DiscreteResourceSpec == nil {
			continue
		}
		item := GenericResource{
			Kind:               resource.DiscreteResourceSpec.Kind,
			Value:              strconv.FormatInt(resource.DiscreteResourceSpec.Value, 10),
			Extensions:         map[string]interface{}{},
			DiscreteExtensions: map[string]interface{}{},
		}
		for key, value := range resource.Extensions {
			item.Extensions[key] = deepCopyValue(value)
		}
		for key, value := range resource.DiscreteResourceSpec.Extensions {
			item.DiscreteExtensions[key] = deepCopyValue(value)
		}
		if len(item.Extensions) == 0 {
			item.Extensions = nil
		}
		if len(item.DiscreteExtensions) == 0 {
			item.DiscreteExtensions = nil
		}
		result = append(result, item)
	}
	return result
}

func buildFromCompose(build *composetypes.BuildConfig) *BuildConfig {
	if build == nil {
		return nil
	}
	result := &BuildConfig{Extensions: copyStringInterfaceMap(build.Extensions)}
	if build.Context != "" {
		result.Context = build.Context
	}
	if build.Dockerfile != "" {
		result.Dockerfile = build.Dockerfile
	}
	if build.DockerfileInline != "" {
		result.Extensions["dockerfile_inline"] = build.DockerfileInline
	}
	if result.Context == "" && result.Dockerfile == "" && len(result.Extensions) == 0 {
		return nil
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	return result
}

func devicesFromCompose(devices []composetypes.DeviceMapping) ([]string, []DeviceMappingSpec) {
	if len(devices) == 0 {
		return nil, nil
	}
	result := make([]string, 0, len(devices))
	mappings := make([]DeviceMappingSpec, 0, len(devices))
	for _, device := range devices {
		spec := DeviceMappingSpec{
			Source:      device.Source,
			Target:      device.Target,
			Permissions: device.Permissions,
			Extensions:  map[string]interface{}{},
		}
		for key, value := range device.Extensions {
			spec.Extensions[key] = value
		}
		if len(spec.Extensions) == 0 {
			spec.Extensions = nil
		}
		result = append(result, deviceMappingToString(spec.Source, spec.Target, spec.Permissions))
		mappings = append(mappings, spec)
	}
	return result, mappings
}

func stringOrNumberListToStrings(values composetypes.StringOrNumberList) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, fmt.Sprintf("%v", value))
	}
	return result
}

func ulimitsFromCompose(ulimits map[string]*composetypes.UlimitsConfig) *Ulimits {
	if len(ulimits) == 0 {
		return nil
	}
	result := &Ulimits{Limits: map[string]UlimitSpec{}}
	for name, limit := range ulimits {
		if name == "" || limit == nil {
			continue
		}
		spec := UlimitSpec{
			Single:     limit.Single,
			Soft:       limit.Soft,
			Hard:       limit.Hard,
			Extensions: map[string]interface{}{},
		}
		for key, value := range limit.Extensions {
			spec.Extensions[key] = value
		}
		if len(spec.Extensions) == 0 {
			spec.Extensions = nil
		}
		if spec.Single == 0 && spec.Soft == 0 && spec.Hard == 0 && len(spec.Extensions) == 0 {
			continue
		}
		result.Limits[name] = spec
		if name == "nofile" {
			result.Nofile = nofileLimitFromSpec(spec)
		}
	}
	if len(result.Limits) == 0 && result.Nofile == nil {
		return nil
	}
	return result
}

func buildConfigToCompose(build *BuildConfig) interface{} {
	if build == nil {
		return nil
	}
	if raw, ok := build.Extensions["compose.build"]; ok {
		switch typed := raw.(type) {
		case string:
			if typed != "" {
				return typed
			}
		case map[string]interface{}:
			result := copyStringInterfaceMap(typed)
			if build.Context != "" {
				result["context"] = build.Context
			}
			if build.Dockerfile != "" {
				result["dockerfile"] = build.Dockerfile
			}
			return result
		default:
			if mapped, ok := asMap(raw); ok && len(mapped) > 0 {
				result := copyStringInterfaceMap(mapped)
				if build.Context != "" {
					result["context"] = build.Context
				}
				if build.Dockerfile != "" {
					result["dockerfile"] = build.Dockerfile
				}
				return result
			}
		}
	}
	if build.Context == "" && build.Dockerfile == "" && len(build.Extensions) == 0 {
		return nil
	}
	if build.Dockerfile == "" {
		return build.Context
	}
	result := map[string]interface{}{}
	if build.Context != "" {
		result["context"] = build.Context
	}
	if build.Dockerfile != "" {
		result["dockerfile"] = build.Dockerfile
	}
	for key, value := range build.Extensions {
		if key == "compose.build" {
			continue
		}
		result[key] = value
	}
	return result
}

func buildConfigHasData(build *BuildConfig) bool {
	return build != nil &&
		(build.Context != "" ||
			build.Dockerfile != "" ||
			len(build.Extensions) > 0)
}

func devicesToCompose(devices []string) []interface{} {
	if len(devices) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(devices))
	for _, device := range devices {
		result = append(result, device)
	}
	return result
}

func deviceMappingsToCompose(devices []string, mappings []DeviceMappingSpec) []interface{} {
	if len(devices) == 0 && len(mappings) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(devices)+len(mappings))
	seen := map[string]struct{}{}
	for _, mapping := range mappings {
		device := deviceMappingToString(mapping.Source, mapping.Target, mapping.Permissions)
		if device == "" {
			continue
		}
		seen[device] = struct{}{}
		if len(mapping.Extensions) == 0 {
			result = append(result, device)
			continue
		}
		item := map[string]interface{}{"source": mapping.Source}
		if mapping.Target != "" {
			item["target"] = mapping.Target
		}
		if mapping.Permissions != "" && mapping.Permissions != "rwm" {
			item["permissions"] = mapping.Permissions
		}
		for key, value := range mapping.Extensions {
			item[composeApplicationExtensionKey(key)] = value
		}
		result = append(result, item)
	}
	for _, device := range devices {
		if _, ok := seen[device]; ok {
			continue
		}
		result = append(result, device)
	}
	return result
}

func ulimitsToCompose(ulimits *Ulimits) map[string]interface{} {
	if isEmptyUlimits(ulimits) {
		return nil
	}
	result := map[string]interface{}{}
	if len(ulimits.Limits) > 0 {
		names := make([]string, 0, len(ulimits.Limits))
		for name := range ulimits.Limits {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			if entry := ulimitSpecToCompose(ulimits.Limits[name]); entry != nil {
				result[name] = entry
			}
		}
	}
	if ulimits.Nofile != nil {
		if _, exists := result["nofile"]; !exists {
			if entry := ulimitSpecToCompose(ulimitSpecFromNofile(ulimits.Nofile)); entry != nil {
				result["nofile"] = entry
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func ulimitSpecToCompose(limit UlimitSpec) interface{} {
	if len(limit.Extensions) == 0 && limit.Single > 0 && limit.Soft == 0 && limit.Hard == 0 {
		return limit.Single
	}
	entry := map[string]interface{}{}
	if limit.Soft > 0 {
		entry["soft"] = limit.Soft
	}
	if limit.Hard > 0 {
		entry["hard"] = limit.Hard
	}
	if limit.Soft == 0 && limit.Hard == 0 && limit.Single > 0 {
		entry["soft"] = limit.Single
		entry["hard"] = limit.Single
	}
	for key, value := range limit.Extensions {
		entry[key] = value
	}
	if len(entry) == 0 {
		return nil
	}
	return entry
}

func ulimitSpecFromNofile(nofile *NofileLimit) UlimitSpec {
	if nofile == nil {
		return UlimitSpec{}
	}
	return UlimitSpec{Soft: nofile.Soft, Hard: nofile.Hard, Extensions: copyStringInterfaceMap(nofile.Extensions)}
}

func nofileLimitFromSpec(spec UlimitSpec) *NofileLimit {
	limit := &NofileLimit{Soft: spec.Soft, Hard: spec.Hard, Extensions: copyStringInterfaceMap(spec.Extensions)}
	if spec.Single > 0 {
		if limit.Soft == 0 {
			limit.Soft = spec.Single
		}
		if limit.Hard == 0 {
			limit.Hard = spec.Single
		}
	}
	if limit.Soft == 0 && limit.Hard == 0 && len(limit.Extensions) == 0 {
		return nil
	}
	return limit
}

func isEmptyUlimits(ulimits *Ulimits) bool {
	return ulimits == nil || (ulimits.Nofile == nil && len(ulimits.Limits) == 0)
}

func deviceMappingToString(source, target, permissions string) string {
	if source == "" {
		return ""
	}
	if target == "" {
		target = source
	}
	if permissions == "" {
		permissions = "rwm"
	}
	if target == source && permissions == "rwm" {
		return source
	}
	if permissions == "rwm" {
		return fmt.Sprintf("%s:%s", source, target)
	}
	return fmt.Sprintf("%s:%s:%s", source, target, permissions)
}

func parseBuildConfig(value interface{}) (*BuildConfig, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return nil, nil
		}
		return &BuildConfig{Context: v, Extensions: map[string]interface{}{"compose.build": v}}, nil
	}
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("build must be a map or string")
	}
	build := &BuildConfig{}
	if context := toString(data["context"]); context != "" {
		build.Context = context
	}
	if dockerfile := toString(data["dockerfile"]); dockerfile != "" {
		build.Dockerfile = dockerfile
	}
	build.Extensions = map[string]interface{}{"compose.build": copyStringInterfaceMap(data)}
	for key, value := range data {
		if key == "context" || key == "dockerfile" {
			continue
		}
		build.Extensions[key] = value
	}
	if build.Context == "" && build.Dockerfile == "" && len(build.Extensions) == 0 {
		return nil, nil
	}
	return build, nil
}

func parseDevices(value interface{}) ([]string, error) {
	devices, _, err := parseDevicesWithMappings(value)
	return devices, err
}

func parseDevicesWithMappings(value interface{}) ([]string, []DeviceMappingSpec, error) {
	switch v := value.(type) {
	case []interface{}:
		result := make([]string, 0, len(v))
		mappings := make([]DeviceMappingSpec, 0, len(v))
		for _, item := range v {
			device, mapping, err := parseDeviceMappingEntry(item)
			if err != nil {
				return nil, nil, err
			}
			if device != "" {
				result = append(result, device)
				mappings = append(mappings, mapping)
			}
		}
		return result, mappings, nil
	case []string:
		result := make([]string, 0, len(v))
		mappings := make([]DeviceMappingSpec, 0, len(v))
		for _, item := range v {
			device, mapping, err := parseDeviceMappingEntry(item)
			if err != nil {
				return nil, nil, err
			}
			if device != "" {
				result = append(result, device)
				mappings = append(mappings, mapping)
			}
		}
		return result, mappings, nil
	case string:
		device, mapping, err := parseDeviceMappingEntry(v)
		if err != nil {
			return nil, nil, err
		}
		if device == "" {
			return nil, nil, nil
		}
		return []string{device}, []DeviceMappingSpec{mapping}, nil
	default:
		return nil, nil, fmt.Errorf("devices must be a list or string, got %T", value)
	}
}

func parseDeviceEntry(value interface{}) (string, error) {
	device, _, err := parseDeviceMappingEntry(value)
	return device, err
}

func parseDeviceMappingEntry(value interface{}) (string, DeviceMappingSpec, error) {
	switch v := value.(type) {
	case string:
		parts := strings.Split(v, ":")
		mapping := DeviceMappingSpec{Source: parts[0]}
		switch len(parts) {
		case 1:
			return deviceMappingToString(parts[0], "", ""), mapping, nil
		case 2:
			mapping.Target = parts[1]
			return deviceMappingToString(parts[0], parts[1], ""), mapping, nil
		default:
			mapping.Target = parts[1]
			mapping.Permissions = strings.Join(parts[2:], ":")
			return deviceMappingToString(mapping.Source, mapping.Target, mapping.Permissions), mapping, nil
		}
	case map[string]interface{}:
		return deviceMappingFromMap(v)
	case map[interface{}]interface{}:
		mapped, _ := asMap(v)
		return deviceMappingFromMap(mapped)
	default:
		return "", DeviceMappingSpec{}, fmt.Errorf("cannot convert device entry: %T", value)
	}
}

func deviceMappingFromMap(data map[string]interface{}) (string, DeviceMappingSpec, error) {
	mapping := DeviceMappingSpec{
		Source:      toString(data["source"]),
		Target:      toString(data["target"]),
		Permissions: toString(data["permissions"]),
		Extensions:  map[string]interface{}{},
	}
	for key, value := range data {
		mapping.Extensions[key] = value
	}
	if len(mapping.Extensions) == 0 {
		mapping.Extensions = nil
	}
	return deviceMappingToString(mapping.Source, mapping.Target, mapping.Permissions), mapping, nil
}

func parseUlimits(value interface{}) (*Ulimits, error) {
	data, ok := asMap(value)
	if !ok {
		return nil, fmt.Errorf("ulimits must be a map")
	}
	result := &Ulimits{Limits: map[string]UlimitSpec{}}
	for name, raw := range data {
		if name == "" {
			continue
		}
		spec, err := parseUlimitSpec(raw)
		if err != nil {
			return nil, err
		}
		if spec.Single == 0 && spec.Soft == 0 && spec.Hard == 0 && len(spec.Extensions) == 0 {
			continue
		}
		result.Limits[name] = spec
		if name == "nofile" {
			result.Nofile = nofileLimitFromSpec(spec)
		}
	}
	if len(result.Limits) == 0 && result.Nofile == nil {
		return nil, nil
	}
	return result, nil
}

func parseNofileLimit(value interface{}) (*NofileLimit, error) {
	spec, err := parseUlimitSpec(value)
	if err != nil {
		return nil, err
	}
	return nofileLimitFromSpec(spec), nil
}

func parseUlimitSpec(value interface{}) (UlimitSpec, error) {
	limit := UlimitSpec{}
	switch v := value.(type) {
	case int, int64, float64, string:
		single := toInt(v)
		if single <= 0 {
			return limit, nil
		}
		limit.Single = single
		return limit, nil
	case map[string]interface{}:
		if single := toInt(v["single"]); single > 0 {
			limit.Single = single
		}
		if soft := toInt(v["soft"]); soft > 0 {
			limit.Soft = soft
		}
		if hard := toInt(v["hard"]); hard > 0 {
			limit.Hard = hard
		}
		limit.Extensions = ulimitExtensionsFromMap(v)
	case map[interface{}]interface{}:
		mapped, _ := asMap(v)
		return parseUlimitSpec(mapped)
	default:
		return limit, fmt.Errorf("invalid ulimit entry: %T", value)
	}
	return limit, nil
}

func ulimitExtensionsFromMap(data map[string]interface{}) map[string]interface{} {
	extensions := map[string]interface{}{}
	for key, value := range data {
		extensions[key] = deepCopyValue(value)
	}
	if len(extensions) == 0 {
		return nil
	}
	return extensions
}

func updatePolicyFromCompose(update *composetypes.UpdateConfig) *UpdatePolicy {
	result := &UpdatePolicy{
		Delay:      update.Delay.String(),
		Order:      update.Order,
		OnFailure:  update.FailureAction,
		Extensions: map[string]interface{}{},
	}
	if update.Parallelism != nil {
		result.Parallelism = int(*update.Parallelism)
		result.ParallelismSet = true
	}
	if update.Monitor != 0 {
		result.Monitor = update.Monitor.String()
	}
	if update.MaxFailureRatio != 0 {
		result.MaxFailureRatio = fmt.Sprintf("%g", update.MaxFailureRatio)
	}
	for key, value := range update.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	return result
}

func restartPolicyFromCompose(restart *composetypes.RestartPolicy) *RestartPolicy {
	result := &RestartPolicy{Condition: restart.Condition, Extensions: map[string]interface{}{}}
	if restart.Delay != nil {
		result.Delay = restart.Delay.String()
	}
	if restart.MaxAttempts != nil {
		result.MaxAttempts = int(*restart.MaxAttempts)
	}
	if restart.Window != nil {
		result.Window = restart.Window.String()
	}
	for key, value := range restart.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	return result
}

func healthCheckFromCompose(health *composetypes.HealthCheckConfig) *HealthCheck {
	if health == nil {
		return nil
	}
	result := &HealthCheck{Test: []string(health.Test), Disable: health.Disable, Extensions: map[string]interface{}{}}
	if health.Disable {
		result.DisableSet = true
	}
	if health.Interval != nil {
		result.Interval = health.Interval.String()
	}
	if health.Timeout != nil {
		result.Timeout = health.Timeout.String()
	}
	if health.Retries != nil {
		result.Retries = int(*health.Retries)
	}
	if health.StartPeriod != nil {
		result.StartPeriod = health.StartPeriod.String()
	}
	if health.StartInterval != nil {
		result.StartInterval = health.StartInterval.String()
	}
	for key, value := range health.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	return normalizeHealthCheck(result)
}

func networkFromCompose(name string, network composetypes.NetworkConfig) *Network {
	result := &Network{
		Name:         name,
		PlatformName: network.Name,
		Driver:       network.Driver,
		DriverOpts:   map[string]string{},
		Attachable:   network.Attachable,
		External:     bool(network.External),
		ExternalSet:  bool(network.External),
		Internal:     network.Internal,
		EnableIPv4:   network.EnableIPv4,
		EnableIPv6:   network.EnableIPv6,
		Labels:       map[string]string{},
		Extensions:   map[string]interface{}{},
	}
	for key, value := range network.DriverOpts {
		result.DriverOpts[key] = value
	}
	for key, value := range network.Labels {
		result.Labels[key] = value
	}
	for key, value := range network.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(network.Ipam.Config) > 0 || network.Ipam.Driver != "" {
		result.IPAM = &IPAMConfig{
			Driver:     network.Ipam.Driver,
			Options:    map[string]string{},
			Extensions: map[string]interface{}{},
		}
		for _, pool := range network.Ipam.Config {
			subnet := IPAMSubnet{
				Subnet:       pool.Subnet,
				Gateway:      pool.Gateway,
				IPRange:      pool.IPRange,
				AuxAddresses: map[string]string{},
				Extensions:   map[string]interface{}{},
			}
			for key, value := range pool.AuxiliaryAddresses {
				subnet.AuxAddresses[key] = value
			}
			for key, value := range pool.Extensions {
				subnet.Extensions[key] = deepCopyValue(value)
			}
			result.IPAM.Config = append(result.IPAM.Config, subnet)
		}
		for key, value := range network.Ipam.Options {
			result.IPAM.Options[key] = value
		}
		for key, value := range network.Ipam.Extensions {
			result.IPAM.Extensions[key] = deepCopyValue(value)
		}
	}
	return result
}

func volumeFromCompose(name string, volume composetypes.VolumeConfig) *Volume {
	result := &Volume{
		Name:         name,
		PlatformName: volume.Name,
		Driver:       volume.Driver,
		DriverOpts:   map[string]string{},
		External:     bool(volume.External),
		ExternalSet:  bool(volume.External),
		Labels:       map[string]string{},
		Extensions:   map[string]interface{}{},
	}
	for key, value := range volume.DriverOpts {
		result.DriverOpts[key] = value
	}
	for key, value := range volume.Labels {
		result.Labels[key] = value
	}
	for key, value := range volume.Extensions {
		result.Extensions[key] = deepCopyValue(value)
		if key == "x-kubernetes-kind" && result.Extensions["kubernetes.kind"] == nil {
			result.Extensions["kubernetes.kind"] = deepCopyValue(value)
		}
	}
	return result
}

func configFromCompose(name string, config composetypes.ConfigObjConfig) *Config {
	file := composetypes.FileObjectConfig(config)
	result := &Config{
		Name:         name,
		PlatformName: file.Name,
		Content:      file.Content,
		Environment:  file.Environment,
		File:         file.File,
		Template:     file.TemplateDriver,
		External:     bool(file.External),
		ExternalSet:  bool(file.External),
		Labels:       map[string]string{},
		Extensions:   map[string]interface{}{},
	}
	for key, value := range file.Labels {
		result.Labels[key] = value
	}
	for key, value := range file.Extensions {
		result.Extensions[key] = deepCopyValue(value)
		switch key {
		case "x-kubernetes-immutable":
			if result.Extensions["kubernetes.immutable"] == nil {
				result.Extensions["kubernetes.immutable"] = deepCopyValue(value)
			}
		case "x-kubernetes-binaryData", "x-kubernetes-binary-data":
			if result.Extensions["kubernetes.binaryData"] == nil {
				result.Extensions["kubernetes.binaryData"] = deepCopyValue(value)
			}
		}
	}
	return result
}

func secretFromCompose(name string, secret composetypes.SecretConfig) *Secret {
	file := composetypes.FileObjectConfig(secret)
	result := &Secret{
		Name:         name,
		PlatformName: file.Name,
		File:         file.File,
		Environment:  file.Environment,
		Template:     file.TemplateDriver,
		External:     bool(file.External),
		ExternalSet:  bool(file.External),
		Driver:       file.Driver,
		DriverOpts:   map[string]string{},
		Labels:       map[string]string{},
		Extensions:   map[string]interface{}{},
	}
	for key, value := range file.DriverOpts {
		result.DriverOpts[key] = value
	}
	for key, value := range file.Labels {
		result.Labels[key] = value
	}
	for key, value := range file.Extensions {
		result.Extensions[key] = deepCopyValue(value)
		switch key {
		case "x-kubernetes-immutable":
			if result.Extensions["kubernetes.immutable"] == nil {
				result.Extensions["kubernetes.immutable"] = deepCopyValue(value)
			}
		case "x-kubernetes-stringData", "x-kubernetes-string-data":
			if result.Extensions["kubernetes.stringData"] == nil {
				result.Extensions["kubernetes.stringData"] = deepCopyValue(value)
			}
		case "x-kubernetes-type":
			if result.Extensions["kubernetes.type"] == nil {
				result.Extensions["kubernetes.type"] = deepCopyValue(value)
			}
		}
	}
	return result
}

func fileRefFromCompose(ref composetypes.FileReferenceConfig) FileRef {
	result := FileRef{
		Source:     ref.Source,
		Key:        toString(ref.Extensions["x-kubernetes-key"]),
		Target:     ref.Target,
		UID:        ref.UID,
		GID:        ref.GID,
		Extensions: map[string]interface{}{},
	}
	if ref.Mode != nil {
		result.Mode = fmt.Sprintf("0%o", *ref.Mode)
	}
	result.Optional = boolPtrFromInterface(ref.Extensions["x-kubernetes-optional"])
	for key, value := range ref.Extensions {
		result.Extensions[key] = deepCopyValue(value)
	}
	if len(result.Extensions) == 0 {
		result.Extensions = nil
	}
	return result
}

func composeBoolPtr(value *bool) *bool {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}

func ensureStringMap(value map[string]string) map[string]string {
	if value == nil {
		return map[string]string{}
	}
	return value
}
