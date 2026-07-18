package failover

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRegistryPersistsAcrossCrashStatus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "services.yaml")

	reg, err := NewRegistry(path)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}

	reg.Upsert(&ServiceEntry{
		Name:            "whoami",
		Port:            80,
		HealthPath:      "/",
		ReplicaEligible: true,
		Status:          StatusRunning,
		Nodes: map[string]*NodePlacement{
			"micklethefickle": {
				Status:   StatusRunning,
				LastSeen: time.Now().UTC(),
				Priority: 100,
			},
		},
	})
	if err := reg.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if !reg.SetStatus("whoami", StatusCrashed, "micklethefickle") {
		t.Fatal("SetStatus returned false")
	}
	if err := reg.Save(); err != nil {
		t.Fatalf("Save after crash: %v", err)
	}

	reg2, err := NewRegistry(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := reg2.Get("whoami")
	if !ok {
		t.Fatal("whoami missing after crash — routes would be dropped")
	}
	if e.Status != StatusCrashed {
		t.Fatalf("status=%s want crashed", e.Status)
	}
	if e.Port != 80 {
		t.Fatalf("port=%d want 80", e.Port)
	}
	if np := e.Nodes["micklethefickle"]; np == nil || np.Status != StatusCrashed {
		t.Fatalf("node placement missing or wrong: %+v", np)
	}

	reg2.Delete("whoami")
	if err := reg2.Save(); err != nil {
		t.Fatalf("Save after delete: %v", err)
	}
	reg3, err := NewRegistry(path)
	if err != nil {
		t.Fatalf("reload after delete: %v", err)
	}
	if _, ok := reg3.Get("whoami"); ok {
		t.Fatal("whoami should be gone after explicit Delete")
	}
}

func TestRegistryIntentionalStopDoesNotDelete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "services.yaml")
	reg, err := NewRegistry(path)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	reg.Upsert(&ServiceEntry{
		Name:   "whoami",
		Port:   80,
		Status: StatusRunning,
		Nodes: map[string]*NodePlacement{
			"micklethefickle": {Status: StatusRunning},
		},
	})
	reg.SetStatus("whoami", StatusIntentionallyStopped, "micklethefickle")
	if err := reg.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	body := string(data)
	if !strings.Contains(body, "intentionally_stopped") {
		t.Fatalf("registry file missing intentionally_stopped:\n%s", body)
	}
	if !strings.Contains(body, "whoami") {
		t.Fatal("whoami removed on intentional stop")
	}
}
