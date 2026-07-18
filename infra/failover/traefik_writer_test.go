package failover

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFailoverYAMLOmitsDeadPeerURLs(t *testing.T) {
	dir := t.TempDir()
	regPath := filepath.Join(dir, "services.yaml")
	outPath := filepath.Join(dir, "failover-fallbacks.yaml")

	reg, err := NewRegistry(regPath)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{
		Name:       "whoami",
		Port:       80,
		HealthPath: "/",
		Status:     StatusCrashed,
		Nodes: map[string]*NodePlacement{
			"micklethefickle": {Status: StatusCrashed},
			"peer1":           {Status: StatusRunning},
		},
	})

	cfg := DefaultTraefikWriterConfig(
		"bolabaden.org",
		"micklethefickle",
		outPath,
		[]string{"peer1", "peer2"},
	)
	if err := WriteFailoverYAML(reg, cfg); err != nil {
		t.Fatalf("WriteFailoverYAML: %v", err)
	}

	body := string(mustRead(t, outPath))
	for _, want := range []string{
		"Host(`whoami.bolabaden.org`)",
		"Host(`whoami.micklethefickle.bolabaden.org`)",
		"Host(`whoami.peer1.bolabaden.org`)",
		"url: https://whoami.peer1.bolabaden.org",
		"whoami-with-failover@file",
		"serversTransport: failover-peers",
		"priority: 1000",
		"sole owner",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q in:\n%s", want, body)
		}
	}
	if strings.Contains(body, "url: https://whoami.peer2.bolabaden.org") {
		t.Fatalf("dead/unknown peer2 must not get backend URL:\n%s", body)
	}
	if strings.Contains(body, "url: http://whoami:80") {
		t.Fatalf("crashed local should omit http://whoami:80:\n%s", body)
	}
	if strings.Contains(body, "service: whoami-direct@file") {
		t.Fatalf("node-direct must not use local-only direct service:\n%s", body)
	}
}

func TestWriteFailoverYAMLIncludesLocalWhenRunning(t *testing.T) {
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
			"n1": {Status: StatusRunning},
			"n2": {Status: StatusRunning},
		},
	})
	out := filepath.Join(dir, "out.yaml")
	if err := WriteFailoverYAML(reg, DefaultTraefikWriterConfig("bolabaden.org", "n1", out, []string{"n2"})); err != nil {
		t.Fatalf("WriteFailoverYAML: %v", err)
	}
	body := string(mustRead(t, out))
	if !strings.Contains(body, "url: http://whoami:80") {
		t.Fatalf("running local should include http backend:\n%s", body)
	}
	if !strings.Contains(body, "url: https://whoami.n2.bolabaden.org") {
		t.Fatalf("missing peer URL:\n%s", body)
	}
	if !strings.Contains(body, "priority: 1000") {
		t.Fatalf("missing priority 1000:\n%s", body)
	}
}

func TestWriteFailoverYAMLPeerOnlyOmitsLocal(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(filepath.Join(dir, "services.yaml"))
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{
		Name:   "uptime-kuma",
		Port:   3001,
		Status: StatusRunning,
		Nodes: map[string]*NodePlacement{
			"n4": {Status: StatusRunning},
		},
	})
	out := filepath.Join(dir, "out.yaml")
	cfg := DefaultTraefikWriterConfig("bolabaden.org", "n1", out, []string{"n2", "n3", "n4"})
	if err := WriteFailoverYAML(reg, cfg); err != nil {
		t.Fatalf("WriteFailoverYAML: %v", err)
	}
	body := string(mustRead(t, out))
	if strings.Contains(body, "url: http://uptime-kuma:3001") {
		t.Fatalf("peer-only service must omit local URL:\n%s", body)
	}
	if !strings.Contains(body, "url: https://uptime-kuma.n4.bolabaden.org") {
		t.Fatalf("must include active peer n4 URL:\n%s", body)
	}
	if strings.Contains(body, "url: https://uptime-kuma.n2.bolabaden.org") {
		t.Fatalf("must not include inactive peer n2:\n%s", body)
	}
}

func TestWriteFailoverYAMLEmitsMiddlewares(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(filepath.Join(dir, "services.yaml"))
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{
		Name:        "app",
		Port:        8080,
		Status:      StatusRunning,
		Middlewares: []string{"tinyauth@file", "crowdsec@file"},
		Nodes: map[string]*NodePlacement{
			"n1": {Status: StatusRunning},
		},
	})
	out := filepath.Join(dir, "out.yaml")
	if err := WriteFailoverYAML(reg, DefaultTraefikWriterConfig("bolabaden.org", "n1", out, []string{"p2"})); err != nil {
		t.Fatalf("WriteFailoverYAML: %v", err)
	}
	body := string(mustRead(t, out))
	if !strings.Contains(body, "- tinyauth@file") || !strings.Contains(body, "- crowdsec@file") {
		t.Fatalf("middlewares missing:\n%s", body)
	}
}

func TestWriteFailoverYAMLIntentionalStopOmitsDeadPeers(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(filepath.Join(dir, "services.yaml"))
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{
		Name:   "whoami",
		Port:   80,
		Status: StatusIntentionallyStopped,
		Nodes: map[string]*NodePlacement{
			"micklethefickle": {Status: StatusIntentionallyStopped},
			"peer1":           {Status: StatusRunning},
		},
	})

	outPath := filepath.Join(dir, "failover-fallbacks.yaml")
	cfg := DefaultTraefikWriterConfig("bolabaden.org", "micklethefickle", outPath, []string{"peer1"})
	if err := WriteFailoverYAML(reg, cfg); err != nil {
		t.Fatalf("WriteFailoverYAML: %v", err)
	}
	body := string(mustRead(t, outPath))
	if !strings.Contains(body, "https://whoami.peer1.bolabaden.org") {
		t.Fatalf("active peer1 should remain:\n%s", body)
	}
	if strings.Contains(body, "url: http://whoami:80") {
		t.Fatalf("intentional stop should omit local URL:\n%s", body)
	}
}

func TestWriteFailoverYAMLSkipsZeroPort(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(filepath.Join(dir, "services.yaml"))
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{Name: "broken", Port: 0, Status: StatusRunning})
	out := filepath.Join(dir, "out.yaml")
	if err := WriteFailoverYAML(reg, DefaultTraefikWriterConfig("bolabaden.org", "n1", out, nil)); err != nil {
		t.Fatalf("WriteFailoverYAML: %v", err)
	}
	if strings.Contains(string(mustRead(t, out)), "broken") {
		t.Fatalf("zero-port service should be skipped")
	}
}

func TestLocalBackendActive(t *testing.T) {
	e := &ServiceEntry{
		Status: StatusRunning,
		Nodes: map[string]*NodePlacement{
			"n1": {Status: StatusCrashed},
			"n3": {Status: StatusRunning},
		},
	}
	if LocalBackendActive(e, "n1") {
		t.Fatal("crashed local should be inactive")
	}
	if LocalBackendActive(e, "n2") {
		t.Fatal("missing local node should be inactive")
	}
	if !LocalBackendActive(e, "n3") {
		t.Fatal("running node should be active")
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	return data
}
