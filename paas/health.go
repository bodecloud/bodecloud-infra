package paas

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var healthURLPattern = regexp.MustCompile(`https?://[^\s'"\\]+`)

func normalizeHealthCheck(health *HealthCheck) *HealthCheck {
	if health == nil {
		return nil
	}
	if health.Type == "" && len(health.Test) > 0 {
		inferHealthCheckProbe(health)
	}
	return health
}

func isEmptyHealthCheck(health *HealthCheck) bool {
	return health == nil ||
		(len(health.Test) == 0 &&
			health.Type == "" &&
			health.Path == "" &&
			health.Port == "" &&
			health.Interval == "" &&
			health.Timeout == "" &&
			health.Retries == 0 &&
			health.StartPeriod == "" &&
			health.StartInterval == "" &&
			!(health.DisableSet || health.Disable) &&
			len(health.Extensions) == 0)
}

func mergeHealthCheckSpec(native, portable *HealthCheck) *HealthCheck {
	if native == nil {
		return cloneHealthCheck(portable)
	}
	if portable == nil {
		return cloneHealthCheck(native)
	}
	merged := cloneHealthCheck(native)
	if len(merged.Test) == 0 {
		merged.Test = append([]string{}, portable.Test...)
	}
	if merged.Type == "" {
		merged.Type = portable.Type
	}
	if merged.Path == "" {
		merged.Path = portable.Path
	}
	if merged.Port == "" {
		merged.Port = portable.Port
	}
	if merged.Interval == "" {
		merged.Interval = portable.Interval
	}
	if merged.Timeout == "" {
		merged.Timeout = portable.Timeout
	}
	if merged.Retries == 0 {
		merged.Retries = portable.Retries
	}
	if merged.StartPeriod == "" {
		merged.StartPeriod = portable.StartPeriod
	}
	if merged.StartInterval == "" {
		merged.StartInterval = portable.StartInterval
	}
	if portable.DisableSet {
		merged.Disable = portable.Disable
		merged.DisableSet = true
	} else if portable.Disable {
		merged.Disable = true
		merged.DisableSet = true
	}
	if len(portable.Extensions) > 0 {
		merged.Extensions = copyStringInterfaceMap(portable.Extensions)
	}
	return normalizeHealthCheck(merged)
}

func serializeHealthCheck(health *HealthCheck) map[string]interface{} {
	health = normalizeHealthCheck(health)
	if health == nil {
		return nil
	}
	result := map[string]interface{}{}
	if len(health.Extensions) > 0 {
		for key, value := range health.Extensions {
			result[key] = deepCopyValue(value)
		}
	}
	if len(health.Test) > 0 {
		result["test"] = append([]string{}, health.Test...)
	}
	if health.Type != "" {
		result["type"] = health.Type
	}
	if health.Path != "" {
		result["path"] = health.Path
	}
	if health.Port != "" {
		result["port"] = health.Port
	}
	if health.Interval != "" {
		result["interval"] = health.Interval
	}
	if health.Timeout != "" {
		result["timeout"] = health.Timeout
	}
	if health.Retries > 0 {
		result["retries"] = health.Retries
	}
	if health.StartPeriod != "" {
		result["start_period"] = health.StartPeriod
	}
	if health.StartInterval != "" {
		result["start_interval"] = health.StartInterval
	}
	if health.DisableSet || health.Disable {
		result["disable"] = health.Disable
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func parseHealthCheckJSON(data []byte) (*HealthCheck, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return parseHealthCheckMap(raw)
}

func parseHealthCheckMap(data map[string]interface{}) (*HealthCheck, error) {
	if len(data) == 0 {
		return nil, nil
	}
	health := &HealthCheck{Extensions: map[string]interface{}{}}
	for key, val := range data {
		switch key {
		case "test":
			if test, err := toStringSlice(val); err == nil {
				health.Test = test
			}
		case "type":
			health.Type = toString(val)
		case "path":
			health.Path = toString(val)
		case "port":
			health.Port = toString(val)
		case "interval":
			health.Interval = toString(val)
		case "timeout":
			health.Timeout = toString(val)
		case "retries":
			health.Retries = toInt(val)
		case "start_period":
			health.StartPeriod = toString(val)
		case "start_interval":
			health.StartInterval = toString(val)
		case "disable":
			health.Disable = toBool(val)
			health.DisableSet = true
		default:
			health.Extensions[key] = deepCopyValue(val)
		}
	}
	if len(health.Extensions) == 0 {
		health.Extensions = nil
	}
	if isEmptyHealthCheck(health) {
		return nil, nil
	}
	return normalizeHealthCheck(health), nil
}

func mergeLifecycleHooks(native, portable *LifecycleHooks) *LifecycleHooks {
	if native == nil {
		return cloneLifecycleHooks(portable)
	}
	if portable == nil {
		return cloneLifecycleHooks(native)
	}
	merged := cloneLifecycleHooks(native)
	merged.PreStart = mergeLifecycleHookSlice(merged.PreStart, portable.PreStart)
	merged.PostStart = mergeLifecycleHookSlice(merged.PostStart, portable.PostStart)
	merged.PreStop = mergeLifecycleHookSlice(merged.PreStop, portable.PreStop)
	if len(merged.PreStart) == 0 && len(merged.PostStart) == 0 && len(merged.PreStop) == 0 {
		merged.PreStart = cloneServiceHooks(portable.PreStart)
		merged.PostStart = cloneServiceHooks(portable.PostStart)
		merged.PreStop = cloneServiceHooks(portable.PreStop)
	}
	if len(portable.Extensions) > 0 {
		if merged.Extensions == nil {
			merged.Extensions = map[string]interface{}{}
		}
		for key, value := range portable.Extensions {
			if _, exists := merged.Extensions[key]; !exists {
				merged.Extensions[key] = deepCopyValue(value)
			}
		}
	}
	return merged
}

func mergeLifecycleHookSlice(native, portable []ServiceHook) []ServiceHook {
	if len(native) == 0 {
		return cloneServiceHooks(portable)
	}
	if len(portable) == 0 {
		return cloneServiceHooks(native)
	}
	merged := cloneServiceHooks(native)
	for i := range portable {
		if i >= len(merged) {
			merged = append(merged, *cloneServiceHook(&portable[i]))
			continue
		}
		merged[i] = *mergeServiceHook(&merged[i], &portable[i])
	}
	return merged
}

func mergeServiceHook(native, portable *ServiceHook) *ServiceHook {
	if native == nil {
		return cloneServiceHook(portable)
	}
	if portable == nil {
		return cloneServiceHook(native)
	}
	merged := cloneServiceHook(native)
	if len(merged.Command) == 0 {
		merged.Command = append([]string{}, portable.Command...)
	}
	if merged.Image == "" {
		merged.Image = portable.Image
	}
	if merged.User == "" {
		merged.User = portable.User
	}
	if !merged.Privileged {
		merged.Privileged = portable.Privileged
	}
	if merged.WorkingDir == "" {
		merged.WorkingDir = portable.WorkingDir
	}
	if len(portable.Environment) > 0 {
		if merged.Environment == nil {
			merged.Environment = map[string]*string{}
		}
		for key, value := range portable.Environment {
			if _, exists := merged.Environment[key]; !exists {
				merged.Environment[key] = value
			}
		}
	}
	if !merged.PerReplica {
		merged.PerReplica = portable.PerReplica
	}
	if len(portable.Extensions) > 0 {
		if merged.Extensions == nil {
			merged.Extensions = map[string]interface{}{}
		}
		for key, value := range portable.Extensions {
			if _, exists := merged.Extensions[key]; !exists {
				merged.Extensions[key] = deepCopyValue(value)
			}
		}
	}
	return merged
}

func inferHealthCheckProbe(health *HealthCheck) {
	raw := strings.Join(health.Test, " ")
	matched := healthURLPattern.FindString(raw)
	if matched == "" {
		if health.Type == "" {
			health.Type = "exec"
		}
		return
	}
	parsed, err := url.Parse(matched)
	if err != nil {
		health.Type = "exec"
		return
	}
	if health.Type == "" {
		health.Type = "http"
	}
	if health.Path == "" {
		health.Path = parsed.EscapedPath()
		if health.Path == "" {
			health.Path = "/"
		}
		if parsed.RawQuery != "" {
			health.Path += "?" + parsed.RawQuery
		}
	}
	if health.Port == "" {
		health.Port = parsed.Port()
		if health.Port == "" {
			switch parsed.Scheme {
			case "https":
				health.Port = "443"
			default:
				health.Port = "80"
			}
		}
	}
}

func healthCheckCommand(health *HealthCheck) []string {
	if health == nil || len(health.Test) == 0 {
		return nil
	}
	test := append([]string(nil), health.Test...)
	if len(test) == 1 {
		if strings.EqualFold(test[0], "NONE") {
			return nil
		}
		return []string{"sh", "-c", test[0]}
	}
	switch strings.ToUpper(test[0]) {
	case "NONE":
		return nil
	case "CMD":
		return test[1:]
	case "CMD-SHELL":
		return []string{"sh", "-c", strings.Join(test[1:], " ")}
	default:
		return test
	}
}

func durationSeconds(value string) int {
	if value == "" {
		return 0
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return parseInt(value)
	}
	seconds := int(duration.Seconds())
	if seconds == 0 && duration > 0 {
		return 1
	}
	return seconds
}

func secondsDuration(value interface{}) string {
	seconds := toInt(value)
	if seconds <= 0 {
		return ""
	}
	return strconv.Itoa(seconds) + "s"
}
