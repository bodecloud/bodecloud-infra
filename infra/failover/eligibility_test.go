package failover

import "testing"

func TestIsReplicaEligible(t *testing.T) {
	tests := []struct {
		name   string
		svc    string
		labels map[string]string
		want   bool
	}{
		{"whoami default", "whoami", map[string]string{"traefik.enable": "true"}, true},
		{"label false wins", "whoami", map[string]string{"failover.replica": "false"}, false},
		{"label true forces redis", "redis", map[string]string{"failover.replica": "true"}, true},
		{"deny redis", "redis", nil, false},
		{"deny mongodb", "mongodb", nil, false},
		{"deny postgres substring", "litellm-postgres", nil, false},
		{"deny headscale", "headscale-server", nil, false},
		{"managed false", "whoami", map[string]string{"failover.managed": "false"}, false},
		{"rabbitmq deny", "rabbitmq", nil, false},
		{"bolabaden default eligible", "bolabaden-nextjs", map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers.bolabaden-nextjs.rule": "Host(`example`)",
		}, true},
		{"bolabaden label true", "bolabaden-nextjs", map[string]string{"failover.replica": "true"}, true},
		{"autokuma label true", "autokuma", map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers.autokuma.rule": "Host(`autokuma.example`)",
			"failover.replica": "true",
		}, true},
		{"headscale never via label true blocked by deny? label wins", "headscale-server", map[string]string{"failover.replica": "true"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsReplicaEligible(tt.svc, tt.labels)
			if got != tt.want {
				t.Fatalf("IsReplicaEligible(%q)=%v want %v", tt.svc, got, tt.want)
			}
		})
	}
}

func TestIsTraefikEnabled(t *testing.T) {
	if IsTraefikEnabled(nil) {
		t.Fatal("nil labels should be false")
	}
	if !IsTraefikEnabled(map[string]string{"traefik.enable": "true"}) {
		t.Fatal("expected true")
	}
	if IsTraefikEnabled(map[string]string{"traefik.enable": "false"}) {
		t.Fatal("expected false")
	}
}

func TestIsHTTPFailoverCandidate(t *testing.T) {
	tcpOnly := map[string]string{
		"traefik.enable":                   "true",
		"traefik.tcp.routers.redis.rule":   "HostSNI(`*`)",
	}
	if IsHTTPFailoverCandidate(tcpOnly) {
		t.Fatal("TCP-only must not be HTTP failover candidate")
	}
	httpOK := map[string]string{
		"traefik.enable": "true",
		"traefik.http.routers.whoami.rule": "Host(`whoami.example`)",
	}
	if !IsHTTPFailoverCandidate(httpOK) {
		t.Fatal("HTTP labels should qualify")
	}
}

func TestClassifyExit(t *testing.T) {
	if ClassifyExit(false, 0, false, true) != ExitIntentional {
		t.Fatal("exit 0 should be intentional")
	}
	if ClassifyExit(false, 137, false, true) != ExitCrash {
		t.Fatal("nonzero should be crash")
	}
	if ClassifyExit(false, 0, true, true) != ExitCrash {
		t.Fatal("OOM should be crash")
	}
	if ClassifyExit(true, 0, false, true) != ExitUnknown {
		t.Fatal("still running → unknown")
	}
}

func TestMiddlewaresFromLabels(t *testing.T) {
	got := MiddlewaresFromLabels(map[string]string{
		"traefik.http.routers.app.middlewares": "tinyauth@file,crowdsec@file",
	})
	if len(got) != 2 || got[0] != "tinyauth@file" {
		t.Fatalf("got %#v", got)
	}
}

func TestResolvePortAndHealthPath(t *testing.T) {
	labels := map[string]string{
		"traefik.http.services.whoami.loadbalancer.server.port": "8080",
		"failover.health_path": "/healthz",
	}
	if got := ResolvePort(labels, 80); got != 8080 {
		t.Fatalf("ResolvePort=%d want 8080", got)
	}
	if got := ResolvePort(nil, 99); got != 99 {
		t.Fatalf("fallback=%d want 99", got)
	}
	if got := HealthPathFromLabels(labels); got != "/healthz" {
		t.Fatalf("HealthPath=%s want /healthz", got)
	}
}

func TestShouldEnsurePeers(t *testing.T) {
	// Intentional stop must not trigger peer ensure; crash/unhealthy may.
	cases := []struct {
		status ServiceStatus
		elig   bool
		want   bool
	}{
		{StatusRunning, true, true},
		{StatusCrashed, true, true},
		{StatusUnhealthy, true, true},
		{StatusIntentionallyStopped, true, false},
		{StatusCrashed, false, false},
		{StatusRunning, false, false},
	}
	for _, c := range cases {
		got := ShouldEnsurePeers(c.status, c.elig)
		if got != c.want {
			t.Errorf("ShouldEnsurePeers(%s, elig=%v)=%v want %v", c.status, c.elig, got, c.want)
		}
	}
}
