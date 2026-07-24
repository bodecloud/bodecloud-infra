package failover

import (
	"testing"
)

func TestConfigFromRegistryEntry(t *testing.T) {
	cfg := ConfigFromRegistryEntry(&ServiceEntry{
		Name:           "whoami",
		ComposeService: "whoami",
		Image:          "docker.io/traefik/whoami:v1.10.3",
		Port:           80,
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers.whoami.rule": "Host(`whoami.example`)",
		},
	})
	if cfg.Name != "whoami" {
		t.Fatalf("name=%q", cfg.Name)
	}
	if cfg.Image != "docker.io/traefik/whoami:v1.10.3" {
		t.Fatalf("image=%q", cfg.Image)
	}
	if cfg.Labels["traefik.enable"] != "true" {
		t.Fatal("missing traefik.enable")
	}
	if _, ok := cfg.ExposedPorts["80/tcp"]; !ok {
		t.Fatal("missing exposed port")
	}
}

func TestPeerIntentionallyStopped(t *testing.T) {
	e := &ServiceEntry{
		Nodes: map[string]*NodePlacement{
			"n2": {Status: StatusIntentionallyStopped},
			"n3": {Status: StatusCrashed},
		},
	}
	if !peerIntentionallyStopped(e, "n2") {
		t.Fatal("n2 should be intentional")
	}
	if peerIntentionallyStopped(e, "n3") {
		t.Fatal("n3 crash is not intentional")
	}
	if peerIntentionallyStopped(e, "n4") {
		t.Fatal("missing peer is not intentional")
	}
}

func TestInComposeEnsureAllowlist(t *testing.T) {
	a := &Agent{cfg: AgentConfig{ComposeEnsureServices: []string{"bolabaden-nextjs", "autokuma"}}}
	if !a.inComposeEnsureAllowlist("bolabaden-nextjs", "bolabaden-nextjs") {
		t.Fatal("bolabaden should be allowlisted")
	}
	if !a.inComposeEnsureAllowlist("autokuma", "") {
		t.Fatal("autokuma should be allowlisted")
	}
	if a.inComposeEnsureAllowlist("whoami", "whoami") {
		t.Fatal("whoami must not be on Tier-A compose ensure allowlist")
	}
	if a.inComposeEnsureAllowlist("headscale-server", "headscale-server") {
		t.Fatal("headscale must not be allowlisted for ensure")
	}
}

func TestHealthOKStrictAllowlist(t *testing.T) {
	a := &Agent{cfg: AgentConfig{ReplicaEnsureStrict: true}}
	if !a.HealthOK() {
		t.Fatal("empty error should be healthy")
	}
	a.allowlistEnsureErr = "bolabaden-nextjs: connect peer docker: boom"
	if a.HealthOK() {
		t.Fatal("strict mode must fail health when allowlist ensure erred")
	}
	a.cfg.ReplicaEnsureStrict = false
	if !a.HealthOK() {
		t.Fatal("non-strict should ignore allowlist ensure err")
	}
}

func TestConfigFromRegistryEntryTierAPorts(t *testing.T) {
	bola := ConfigFromRegistryEntry(&ServiceEntry{
		Name: "bolabaden-nextjs", ComposeService: "bolabaden-nextjs",
		Image: "docker.io/bolabaden/bolabaden-nextjs:latest", Port: 3000,
		Labels: map[string]string{"traefik.enable": "true"},
	})
	if _, ok := bola.ExposedPorts["3000/tcp"]; !ok {
		t.Fatalf("bolabaden port missing: %#v", bola.ExposedPorts)
	}
	aku := ConfigFromRegistryEntry(&ServiceEntry{
		Name: "autokuma", ComposeService: "autokuma",
		Image: "ghcr.io/bigboot/autokuma:latest", Port: 8080,
		Labels: map[string]string{"traefik.enable": "true", "failover.replica": "true"},
	})
	if _, ok := aku.ExposedPorts["8080/tcp"]; !ok {
		t.Fatalf("autokuma port missing: %#v", aku.ExposedPorts)
	}
}
