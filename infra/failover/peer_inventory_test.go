package failover

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func TestAggregateStatusPrefersLocal(t *testing.T) {
	e := &ServiceEntry{
		Status: StatusCrashed,
		Nodes: map[string]*NodePlacement{
			"n1": {Status: StatusRunning},
			"n3": {Status: StatusCrashed},
		},
	}
	if got := aggregateStatus(e, "n1"); got != StatusRunning {
		t.Fatalf("got %s want running", got)
	}
	e.Nodes["n1"].Status = StatusCrashed
	e.Nodes["n3"].Status = StatusRunning
	if got := aggregateStatus(e, "n1"); got != StatusRunning {
		t.Fatalf("got %s want peer running", got)
	}
}

func TestUpsertFromPeerInspectNeverLocal(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(filepath.Join(dir, "services.yaml"))
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	a := &Agent{
		cfg: AgentConfig{
			LocalNode: "n1",
			PeerHosts: []string{"n3", "n4"},
			Domain:    "bolabaden.org",
		},
		registry: reg,
		restarts: map[string]int{},
	}

	labels := map[string]string{
		"traefik.enable": "true",
		"traefik.http.routers.whoami.rule": "Host(`whoami.bolabaden.org`)",
		"traefik.http.services.whoami.loadbalancer.server.port": "80",
		"com.docker.compose.service": "whoami",
	}
	inspect := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:   "abc123peer",
			Name: "/whoami",
			State: &types.ContainerState{
				Running: true,
				Status:  "running",
			},
		},
		Config: &container.Config{
			Image:  "traefik/whoami:latest",
			Labels: labels,
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
		},
		NetworkSettings: &types.NetworkSettings{},
	}

	name, err := a.upsertFromPeerInspect("n3", inspect)
	if err != nil {
		t.Fatalf("upsertFromPeerInspect: %v", err)
	}
	if name != "whoami" {
		t.Fatalf("name=%q", name)
	}
	e, ok := reg.Get("whoami")
	if !ok {
		t.Fatal("missing registry entry")
	}
	if _, hasLocal := e.Nodes["n1"]; hasLocal {
		t.Fatal("must not invent local node placement")
	}
	if e.Nodes["n3"] == nil || e.Nodes["n3"].Status != StatusRunning {
		t.Fatalf("peer n3 placement: %+v", e.Nodes["n3"])
	}
	if e.ContainerID != "" {
		t.Fatal("peer ingest must not set local ContainerID")
	}
	if LocalBackendActive(e, "n1") {
		t.Fatal("never-local must omit local backend")
	}
}

func TestMarkPeerMissing(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(filepath.Join(dir, "services.yaml"))
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{
		Name:   "whoami",
		Port:   80,
		Status: StatusRunning,
		Nodes: map[string]*NodePlacement{
			"n3": {Status: StatusRunning, LastSeen: time.Now()},
		},
	})
	a := &Agent{cfg: AgentConfig{LocalNode: "n1"}, registry: reg, restarts: map[string]int{}}
	a.markPeerMissing("whoami", "n3")
	e, _ := reg.Get("whoami")
	if e.Nodes["n3"].Status != StatusCrashed {
		t.Fatalf("want crashed, got %s", e.Nodes["n3"].Status)
	}
}

func TestUpsertFromPeerPreservesLocalContainerID(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(filepath.Join(dir, "services.yaml"))
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{
		Name:        "whoami",
		Port:        80,
		Status:      StatusRunning,
		ContainerID: "local-id",
		Nodes: map[string]*NodePlacement{
			"n1": {Status: StatusRunning},
		},
	})
	a := &Agent{
		cfg:      AgentConfig{LocalNode: "n1", PeerHosts: []string{"n3"}},
		registry: reg,
		restarts: map[string]int{},
	}
	labels := map[string]string{
		"traefik.enable":                   "true",
		"traefik.http.routers.whoami.rule": "Host(`x`)",
		"com.docker.compose.service":       "whoami",
	}
	inspect := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    "peer-id",
			Name:  "/whoami",
			State: &types.ContainerState{Running: true},
		},
		Config: &container.Config{
			Image:  "traefik/whoami:latest",
			Labels: labels,
		},
		NetworkSettings: &types.NetworkSettings{},
	}
	if _, err := a.upsertFromPeerInspect("n3", inspect); err != nil {
		t.Fatal(err)
	}
	e, _ := reg.Get("whoami")
	if e.ContainerID != "local-id" {
		t.Fatalf("ContainerID overwritten: %s", e.ContainerID)
	}
	if e.Nodes["n1"].Status != StatusRunning || e.Nodes["n3"].Status != StatusRunning {
		t.Fatalf("nodes: %+v", e.Nodes)
	}
}
