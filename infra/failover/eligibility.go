package failover

import (
	"os"
	"strings"
)

// DefaultReplicaDenyPatterns are service name substrings that skip always-on replicas.
var DefaultReplicaDenyPatterns = []string{
	"redis",
	"mongodb",
	"mongo",
	"postgres",
	"postgresql",
	"rabbitmq",
	"headscale",
	"nuq-postgres",
	"litellm-postgres",
}

// IsReplicaEligible returns whether a service should be replicated to peers.
// Label failover.replica=false always wins. failover.replica=true forces include.
// Otherwise deny-list patterns and failover.managed=false exclude.
func IsReplicaEligible(name string, labels map[string]string) bool {
	if labels != nil {
		if v, ok := labels["failover.replica"]; ok {
			switch strings.ToLower(strings.TrimSpace(v)) {
			case "false", "0", "no":
				return false
			case "true", "1", "yes":
				return true
			}
		}
		if v, ok := labels["failover.managed"]; ok {
			switch strings.ToLower(strings.TrimSpace(v)) {
			case "false", "0", "no":
				return false
			}
		}
	}

	base := strings.ToLower(sanitizeName(name))
	for _, pat := range DefaultReplicaDenyPatterns {
		if strings.Contains(base, pat) {
			return false
		}
	}
	return true
}

// IsTraefikEnabled reports whether the container opts into Traefik routing.
func IsTraefikEnabled(labels map[string]string) bool {
	if labels == nil {
		return false
	}
	v, ok := labels["traefik.enable"]
	if !ok {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes":
		return true
	default:
		return false
	}
}

// HasHTTPTraefikLabels is true when the container defines HTTP Traefik routes.
// TCP-only services (redis/mongo with traefik.enable=true) must be excluded from
// HTTP failover YAML — matches scripts/osvc_ingress_sync.py.
func HasHTTPTraefikLabels(labels map[string]string) bool {
	if labels == nil {
		return false
	}
	for k := range labels {
		if strings.HasPrefix(k, "traefik.http.") {
			return true
		}
	}
	return false
}

// IsHTTPFailoverCandidate requires enable + at least one HTTP Traefik label.
func IsHTTPFailoverCandidate(labels map[string]string) bool {
	return IsTraefikEnabled(labels) && HasHTTPTraefikLabels(labels)
}

// MiddlewaresFromLabels collects traefik.http.routers.*.middlewares values.
func MiddlewaresFromLabels(labels map[string]string) []string {
	if labels == nil {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for k, v := range labels {
		if !strings.HasPrefix(k, "traefik.http.routers.") || !strings.HasSuffix(k, ".middlewares") {
			continue
		}
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" || seen[part] {
				continue
			}
			seen[part] = true
			out = append(out, part)
		}
	}
	return out
}

// ExitKind classifies a container exit for crash vs intentional stop.
type ExitKind string

const (
	ExitCrash       ExitKind = "crash"
	ExitIntentional ExitKind = "intentional"
	ExitUnknown     ExitKind = "unknown"
)

// ClassifyExit maps Docker exit metadata to crash vs intentional stop.
// ExitCode 0 without OOM is intentional (compose stop / docker stop).
// Non-zero exit or OOM is a crash. Unknown when state is unavailable.
func ClassifyExit(running bool, exitCode int, oomKilled bool, hasState bool) ExitKind {
	if !hasState {
		return ExitUnknown
	}
	if running {
		return ExitUnknown
	}
	if oomKilled || exitCode != 0 {
		return ExitCrash
	}
	return ExitIntentional
}

// ComposeServiceName extracts a compose service name from labels when present.
func ComposeServiceName(labels map[string]string, containerName string) string {
	if labels != nil {
		if v := labels["com.docker.compose.service"]; v != "" {
			return v
		}
	}
	return sanitizeName(containerName)
}

// ResolvePort picks the published/exposed HTTP port from common Traefik labels.
func ResolvePort(labels map[string]string, fallback int) int {
	if labels != nil {
		for k, v := range labels {
			if strings.HasSuffix(k, ".loadbalancer.server.port") ||
				strings.HasSuffix(k, ".loadBalancer.server.port") {
				if n := parsePort(v); n > 0 {
					return n
				}
			}
		}
	}
	if fallback > 0 {
		return fallback
	}
	return 80
}

func parsePort(s string) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// HealthPathFromLabels reads kuma/traefik health path hints.
func HealthPathFromLabels(labels map[string]string) string {
	if labels == nil {
		return "/"
	}
	for _, key := range []string{
		"kuma.healthcheck.path",
		"failover.health_path",
		"traefik.http.services.healthcheck.path",
	} {
		if v := labels[key]; v != "" {
			return v
		}
	}
	return "/"
}

// ShouldEnsurePeers reports whether peer replica ensure should run for this status.
// Intentional stops never peer-redeploy; crash/unhealthy/running do when eligible.
func ShouldEnsurePeers(status ServiceStatus, replicaEligible bool) bool {
	if !replicaEligible {
		return false
	}
	return status != StatusIntentionallyStopped
}

// EnsureDir creates a directory if missing.
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}
