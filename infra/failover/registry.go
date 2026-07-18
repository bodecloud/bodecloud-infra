package failover

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// ServiceStatus is the placement status for a registered service.
type ServiceStatus string

const (
	StatusRunning             ServiceStatus = "running"
	StatusUnhealthy           ServiceStatus = "unhealthy"
	StatusCrashed             ServiceStatus = "crashed"
	StatusIntentionallyStopped ServiceStatus = "intentionally_stopped"
	StatusUnknown             ServiceStatus = "unknown"
)

// NodePlacement tracks one node's view of a service instance.
type NodePlacement struct {
	Status   ServiceStatus `yaml:"status"`
	LastSeen time.Time     `yaml:"last_seen"`
	Priority int           `yaml:"priority,omitempty"`
}

// ServiceEntry is one service in the placement registry.
// Entries are never deleted on crash/unhealthy — only status changes.
type ServiceEntry struct {
	Name              string                    `yaml:"name"`
	Port              int                       `yaml:"port"`
	HealthPath        string                    `yaml:"health_path"`
	HealthInterval    string                    `yaml:"health_interval,omitempty"`
	HealthTimeout     string                    `yaml:"health_timeout,omitempty"`
	ReplicaEligible   bool                      `yaml:"replica_eligible"`
	ComposeService    string                    `yaml:"compose_service,omitempty"`
	Image             string                    `yaml:"image,omitempty"`
	ContainerID       string                    `yaml:"container_id,omitempty"`
	Labels            map[string]string         `yaml:"labels,omitempty"`
	Middlewares       []string                  `yaml:"middlewares,omitempty"`
	Nodes             map[string]*NodePlacement `yaml:"nodes"`
	Status            ServiceStatus             `yaml:"status"`
	UpdatedAt         time.Time                 `yaml:"updated_at"`
}

// RegistryFile is the on-disk services.yaml shape.
type RegistryFile struct {
	Services map[string]*ServiceEntry `yaml:"services"`
}

// Registry is a thread-safe placement registry persisted to YAML.
type Registry struct {
	mu       sync.RWMutex
	path     string
	services map[string]*ServiceEntry
}

// NewRegistry loads or creates a placement registry at path.
func NewRegistry(path string) (*Registry, error) {
	r := &Registry{
		path:     path,
		services: make(map[string]*ServiceEntry),
	}
	if err := r.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return r, nil
}

// Load reads services.yaml from disk.
func (r *Registry) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.path)
	if err != nil {
		return err
	}
	var file RegistryFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("parse registry: %w", err)
	}
	if file.Services == nil {
		file.Services = make(map[string]*ServiceEntry)
	}
	r.services = file.Services
	return nil
}

// Save atomically writes the registry to disk.
func (r *Registry) Save() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	file := RegistryFile{Services: r.services}
	data, err := yaml.Marshal(&file)
	if err != nil {
		return err
	}
	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path)
}

// Upsert inserts or updates a service entry. Crash/unhealthy never deletes.
func (r *Registry) Upsert(entry *ServiceEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entry.Nodes == nil {
		entry.Nodes = make(map[string]*NodePlacement)
	}
	if entry.Labels == nil {
		entry.Labels = make(map[string]string)
	}
	entry.UpdatedAt = time.Now().UTC()
	r.services[entry.Name] = entry
}

// Get returns a copy-safe pointer to a service (caller must not mutate without Upsert).
func (r *Registry) Get(name string) (*ServiceEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.services[name]
	return e, ok
}

// SetStatus updates status for a service without removing the entry.
func (r *Registry) SetStatus(name string, status ServiceStatus, node string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	e, ok := r.services[name]
	if !ok {
		return false
	}
	e.Status = status
	e.UpdatedAt = time.Now().UTC()
	if node != "" {
		if e.Nodes == nil {
			e.Nodes = make(map[string]*NodePlacement)
		}
		np, exists := e.Nodes[node]
		if !exists {
			np = &NodePlacement{}
			e.Nodes[node] = np
		}
		np.Status = status
		np.LastSeen = time.Now().UTC()
	}
	return true
}

// Delete removes a service only on explicit deregister (not crash).
func (r *Registry) Delete(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services, name)
}

// List returns all registered services.
func (r *Registry) List() []*ServiceEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*ServiceEntry, 0, len(r.services))
	for _, e := range r.services {
		out = append(out, e)
	}
	return out
}

// Path returns the registry file path.
func (r *Registry) Path() string {
	return r.path
}
