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
