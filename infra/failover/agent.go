package failover

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// AgentConfig holds runtime configuration for the failover agent.
type AgentConfig struct {
	Enabled                     bool
	Domain                      string
	LocalNode                   string
	MainHost                    string
	PeerHosts                   []string
	RegistryPath                string
	TraefikOutputPath           string
	RemoteDockerPort            int
	RemoteDockerTLS             bool
	MaxLocalRestarts            int
	ReconcileInterval           time.Duration
	StopReplicasOnIntentional   bool
	ReplicaEnsure               bool // FAILOVER_REPLICA_ENSURE — off until peer Docker proven
	ReplicaPullNever            bool // FAILOVER_REPLICA_PULL=never — skip ImagePull on peers
	ReplicaEnsureStrict         bool // fail /healthz when allowlist ensure errors
	ComposeEnsureServices       []string // FAILOVER_COMPOSE_ENSURE_SERVICES — minimal recreate path
	HealthListenAddr            string
}

// Agent is the Compose-first failover/replica maintenance loop.
type Agent struct {
	cfg              AgentConfig
	docker           *client.Client
	registry         *Registry
	restarts         map[string]int
	mu               sync.Mutex
	log              *log.Logger
	allowlistEnsureErr string
}

// NewAgent constructs an Agent.
func NewAgent(cfg AgentConfig, docker *client.Client, logger *log.Logger) (*Agent, error) {
	if logger == nil {
		logger = log.Default()
	}
	if cfg.RemoteDockerPort == 0 {
		cfg.RemoteDockerPort = 2375
	}
	if cfg.MaxLocalRestarts == 0 {
		cfg.MaxLocalRestarts = 3
	}
	if cfg.ReconcileInterval == 0 {
		cfg.ReconcileInterval = 30 * time.Second
	}
	if cfg.RegistryPath == "" {
		return nil, fmt.Errorf("registry path required")
	}
	if cfg.TraefikOutputPath == "" {
		return nil, fmt.Errorf("traefik output path required")
	}

	reg, err := NewRegistry(cfg.RegistryPath)
	if err != nil {
		return nil, err
	}

	return &Agent{
		cfg:      cfg,
		docker:   docker,
		registry: reg,
		restarts: make(map[string]int),
		log:      logger,
	}, nil
}

// Registry returns the placement registry.
func (a *Agent) Registry() *Registry {
	return a.registry
}

// IsMainHost reports whether this node is the primary placement authority.
func (a *Agent) IsMainHost() bool {
	return a.cfg.LocalNode != "" && a.cfg.LocalNode == a.cfg.MainHost
}

// Run blocks until ctx is cancelled.
func (a *Agent) Run(ctx context.Context) error {
	if !a.cfg.Enabled {
		a.log.Printf("failover-agent disabled (FAILOVER_ENABLED=false)")
		<-ctx.Done()
		return ctx.Err()
	}

	if err := a.Reconcile(ctx); err != nil {
		a.log.Printf("initial reconcile: %v", err)
	}

	eventCh, errCh := a.docker.Events(ctx, types.EventsOptions{
		Filters: filters.NewArgs(
			filters.Arg("type", "container"),
			filters.Arg("event", "start"),
			filters.Arg("event", "die"),
			filters.Arg("event", "oom"),
			filters.Arg("event", "health_status"),
			filters.Arg("event", "destroy"),
			filters.Arg("event", "stop"),
		),
	})

	ticker := time.NewTicker(a.cfg.ReconcileInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			if err != nil && ctx.Err() == nil {
				a.log.Printf("docker events error: %v — reconnecting in 5s", err)
				time.Sleep(5 * time.Second)
				eventCh, errCh = a.docker.Events(ctx, types.EventsOptions{
					Filters: filters.NewArgs(
						filters.Arg("type", "container"),
						filters.Arg("event", "start"),
						filters.Arg("event", "die"),
						filters.Arg("event", "oom"),
						filters.Arg("event", "health_status"),
						filters.Arg("event", "destroy"),
						filters.Arg("event", "stop"),
					),
				})
			}
		case ev := <-eventCh:
			a.handleEvent(ctx, ev)
		case <-ticker.C:
			if err := a.Reconcile(ctx); err != nil {
				a.log.Printf("reconcile: %v", err)
			}
		}
	}
}

func (a *Agent) handleEvent(ctx context.Context, ev events.Message) {
	id := ev.Actor.ID
	if id == "" {
		id = ev.ID
	}
	action := string(ev.Action)
	if action == "" {
		action = ev.Status
	}

	switch {
	case action == "start" || strings.HasPrefix(action, "start"):
		if err := a.ingestContainer(ctx, id); err != nil {
			a.log.Printf("ingest on start %s: %v", shortID(id), err)
		}
		_ = a.writeTraefik()
		if a.IsMainHost() && a.cfg.ReplicaEnsure {
			a.ensurePeersForID(ctx, id)
		}
	case action == "die" || action == "oom" || strings.HasPrefix(action, "die"):
		a.onDieOrOOM(ctx, id, action)
	case strings.HasPrefix(action, "health_status"):
		a.onHealth(ctx, id, action)
	case action == "stop":
		a.onIntentionalStop(ctx, id)
	case action == "destroy":
		// destroy after intentional stop: keep registry unless unmanaged
		a.log.Printf("container destroyed %s — registry retained", shortID(id))
		_ = a.writeTraefik()
	}
}

func (a *Agent) onDieOrOOM(ctx context.Context, id, action string) {
	// docker stop / compose stop emit die with ExitCode 0 — must not peer-redeploy.
	if action != "oom" {
		inspect, err := a.docker.ContainerInspect(ctx, id)
		if err == nil && inspect.State != nil {
			kind := ClassifyExit(
				inspect.State.Running,
				inspect.State.ExitCode,
				inspect.State.OOMKilled,
				true,
			)
			if kind == ExitIntentional {
				a.onIntentionalStop(ctx, id)
				return
			}
		}
	}
	a.onCrash(ctx, id, action)
}

func (a *Agent) onCrash(ctx context.Context, id, action string) {
	name, entry := a.lookupByContainerID(id)
	if entry == nil {
		// Try ingest from inspect of stopped container
		if err := a.ingestContainer(ctx, id); err != nil {
			a.log.Printf("crash for unknown container %s (%s)", shortID(id), action)
			return
		}
		name, entry = a.lookupByContainerID(id)
	}
	if entry == nil {
		return
	}

	a.registry.SetStatus(name, StatusCrashed, a.cfg.LocalNode)
	a.log.Printf("service %s crashed (%s) — keeping routes", name, action)
	_ = a.registry.Save()
	_ = a.writeTraefik()

	a.mu.Lock()
	a.restarts[name]++
	count := a.restarts[name]
	a.mu.Unlock()

	if count <= a.cfg.MaxLocalRestarts {
		a.log.Printf("attempting local restart %d/%d for %s", count, a.cfg.MaxLocalRestarts, name)
		if err := a.docker.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
			a.log.Printf("local restart failed for %s: %v", name, err)
		}
	}

	if a.IsMainHost() && a.cfg.ReplicaEnsure && entry.ReplicaEligible {
		if err := a.EnsureReplicas(ctx, entry); err != nil {
			a.log.Printf("peer ensure after crash %s: %v", name, err)
		}
	}
}

func (a *Agent) onHealth(ctx context.Context, id, action string) {
	unhealthy := strings.Contains(action, "unhealthy")
	name, entry := a.lookupByContainerID(id)
	if entry == nil {
		_ = a.ingestContainer(ctx, id)
		name, entry = a.lookupByContainerID(id)
	}
	if entry == nil {
		return
	}
	if unhealthy {
		a.registry.SetStatus(name, StatusUnhealthy, a.cfg.LocalNode)
		a.log.Printf("service %s unhealthy — keeping routes", name)
		_ = a.registry.Save()
		_ = a.writeTraefik()
		if a.IsMainHost() && a.cfg.ReplicaEnsure && entry.ReplicaEligible {
			_ = a.EnsureReplicas(ctx, entry)
		}
	} else if strings.Contains(action, "healthy") {
		a.registry.SetStatus(name, StatusRunning, a.cfg.LocalNode)
		a.mu.Lock()
		a.restarts[name] = 0
		a.mu.Unlock()
		_ = a.registry.Save()
		_ = a.writeTraefik()
	}
}

func (a *Agent) onIntentionalStop(ctx context.Context, id string) {
	name, entry := a.lookupByContainerID(id)
	if entry == nil {
		_ = a.ingestContainer(ctx, id)
		name, entry = a.lookupByContainerID(id)
	}
	if entry == nil {
		return
	}

	// If stop event races with a true crash, reclassify.
	inspect, err := a.docker.ContainerInspect(ctx, id)
	if err == nil && inspect.State != nil {
		kind := ClassifyExit(
			inspect.State.Running,
			inspect.State.ExitCode,
			inspect.State.OOMKilled,
			true,
		)
		if kind == ExitCrash {
			a.onCrash(ctx, id, "die")
			return
		}
	}

	a.registry.SetStatus(name, StatusIntentionallyStopped, a.cfg.LocalNode)
	a.log.Printf("service %s intentionally stopped — not peer-redeploying", name)
	_ = a.registry.Save()
	_ = a.writeTraefik()

	if a.cfg.StopReplicasOnIntentional && a.IsMainHost() {
		a.log.Printf("FAILOVER_STOP_REPLICAS_ON_INTENTIONAL_STOP=true — peer stop not implemented in v1")
	}
}

// Reconcile scans running containers and refreshes registry + Traefik + peers.
func (a *Agent) Reconcile(ctx context.Context) error {
	list, err := a.docker.ContainerList(ctx, types.ContainerListOptions{All: false})
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	for _, c := range list {
		if err := a.ingestListed(ctx, c); err != nil {
			a.log.Printf("ingest %s: %v", shortID(c.ID), err)
			continue
		}
		name := ComposeServiceName(c.Labels, firstName(c.Names))
		seen[name] = true
	}

	// Mark previously running local placements not seen as crashed (keep routes).
	// Do not crash peer-only services that were never local.
	for _, e := range a.registry.List() {
		if seen[e.Name] {
			continue
		}
		localActive := false
		if e.Nodes != nil {
			if np, ok := e.Nodes[a.cfg.LocalNode]; ok && np != nil {
				localActive = np.Status == StatusRunning || np.Status == StatusUnhealthy
			}
		} else if e.Status == StatusRunning || e.Status == StatusUnhealthy {
			localActive = true
		}
		if localActive {
			a.registry.SetStatus(e.Name, StatusCrashed, a.cfg.LocalNode)
			a.log.Printf("reconcile: %s missing locally — marked crashed, routes kept", e.Name)
		}
	}

	// RO peer Docker inventory so never-local HTTP services get Traefik routes.
	if err := a.IngestPeerInventory(ctx); err != nil {
		a.log.Printf("peer inventory: %v", err)
	}

	if err := a.registry.Save(); err != nil {
		return err
	}
	if err := a.writeTraefik(); err != nil {
		return err
	}

	if a.IsMainHost() && a.cfg.ReplicaEnsure {
		var allowErrs []string
		for _, e := range a.registry.List() {
			if !ShouldEnsurePeers(e.Status, e.ReplicaEligible) {
				continue
			}
			if err := a.EnsureReplicas(ctx, e); err != nil {
				a.log.Printf("ensure replicas %s: %v", e.Name, err)
				if a.inComposeEnsureAllowlist(e.Name, e.ComposeService) {
					allowErrs = append(allowErrs, fmt.Sprintf("%s: %v", e.Name, err))
				}
			}
		}
		a.mu.Lock()
		if len(allowErrs) > 0 {
			a.allowlistEnsureErr = strings.Join(allowErrs, "; ")
		} else {
			a.allowlistEnsureErr = ""
		}
		a.mu.Unlock()
	}
	return nil
}

// AllowlistEnsureError returns the last allowlist ensure failure (empty if healthy).
func (a *Agent) AllowlistEnsureError() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.allowlistEnsureErr
}

// HealthOK is false when ReplicaEnsureStrict and an allowlist ensure failed.
func (a *Agent) HealthOK() bool {
	if !a.cfg.ReplicaEnsureStrict {
		return true
	}
	return a.AllowlistEnsureError() == ""
}

func (a *Agent) inComposeEnsureAllowlist(name, composeService string) bool {
	if len(a.cfg.ComposeEnsureServices) == 0 {
		return false
	}
	candidates := []string{sanitizeName(name), sanitizeName(composeService)}
	for _, want := range a.cfg.ComposeEnsureServices {
		want = sanitizeName(want)
		if want == "" {
			continue
		}
		for _, c := range candidates {
			if c == want {
				return true
			}
		}
	}
	return false
}

func (a *Agent) ingestContainer(ctx context.Context, id string) error {
	inspect, err := a.docker.ContainerInspect(ctx, id)
	if err != nil {
		return err
	}
	return a.upsertFromInspect(inspect)
}

func (a *Agent) ingestListed(ctx context.Context, c types.Container) error {
	if !IsHTTPFailoverCandidate(c.Labels) {
		return nil
	}
	inspect, err := a.docker.ContainerInspect(ctx, c.ID)
	if err != nil {
		return err
	}
	return a.upsertFromInspect(inspect)
}

func (a *Agent) upsertFromInspect(inspect types.ContainerJSON) error {
	labels := inspect.Config.Labels
	if !IsHTTPFailoverCandidate(labels) {
		return nil
	}

	name := ComposeServiceName(labels, inspect.Name)
	port := ResolvePort(labels, firstExposedPort(inspect))
	entry := &ServiceEntry{
		Name:            name,
		Port:            port,
		HealthPath:      HealthPathFromLabels(labels),
		HealthInterval:  labelOr(labels, "kuma.healthcheck.interval", "15s"),
		HealthTimeout:   labelOr(labels, "kuma.healthcheck.timeout", "5s"),
		ReplicaEligible: IsReplicaEligible(name, labels),
		ComposeService:  ComposeServiceName(labels, inspect.Name),
		Image:           inspect.Config.Image,
		ContainerID:     inspect.ID,
		Labels:          labels,
		Middlewares:     MiddlewaresFromLabels(labels),
		Nodes: map[string]*NodePlacement{
			a.cfg.LocalNode: {
				Status:   StatusRunning,
				LastSeen: time.Now().UTC(),
				Priority: 1,
			},
		},
		Status: StatusRunning,
	}
	if inspect.State != nil {
		if inspect.State.Health != nil {
			switch inspect.State.Health.Status {
			case "unhealthy":
				entry.Status = StatusUnhealthy
				entry.Nodes[a.cfg.LocalNode].Status = StatusUnhealthy
			case "healthy", "starting":
				entry.Status = StatusRunning
			}
		}
		if !inspect.State.Running {
			if inspect.State.ExitCode == 0 && !inspect.State.OOMKilled {
				entry.Status = StatusIntentionallyStopped
			} else {
				entry.Status = StatusCrashed
			}
			entry.Nodes[a.cfg.LocalNode].Status = entry.Status
		}
	}

	// Preserve peer node entries from existing registry
	if existing, ok := a.registry.Get(name); ok && existing.Nodes != nil {
		for node, np := range existing.Nodes {
			if node == a.cfg.LocalNode {
				continue
			}
			entry.Nodes[node] = np
		}
	}

	a.registry.Upsert(entry)
	return nil
}

func (a *Agent) writeTraefik() error {
	cfg := DefaultTraefikWriterConfig(
		a.cfg.Domain,
		a.cfg.LocalNode,
		a.cfg.TraefikOutputPath,
		a.cfg.PeerHosts,
	)
	if err := WriteFailoverYAML(a.registry, cfg); err != nil {
		return err
	}
	a.log.Printf("wrote traefik failover config → %s (%d services)", a.cfg.TraefikOutputPath, len(a.registry.List()))
	return nil
}

func (a *Agent) lookupByContainerID(id string) (string, *ServiceEntry) {
	for _, e := range a.registry.List() {
		if e.ContainerID == id || strings.HasPrefix(e.ContainerID, id) || strings.HasPrefix(id, e.ContainerID) {
			return e.Name, e
		}
	}
	return "", nil
}

func (a *Agent) ensurePeersForID(ctx context.Context, id string) {
	_, entry := a.lookupByContainerID(id)
	if entry == nil {
		_ = a.ingestContainer(ctx, id)
		_, entry = a.lookupByContainerID(id)
	}
	if entry == nil || !entry.ReplicaEligible {
		return
	}
	if err := a.EnsureReplicas(ctx, entry); err != nil {
		a.log.Printf("ensure peers for %s: %v", entry.Name, err)
	}
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func firstName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func firstExposedPort(inspect types.ContainerJSON) int {
	if inspect.NetworkSettings != nil {
		for p := range inspect.NetworkSettings.Ports {
			if n := parsePort(string(p.Port())); n > 0 {
				return n
			}
		}
	}
	if inspect.Config != nil {
		for p := range inspect.Config.ExposedPorts {
			if n := parsePort(string(p.Port())); n > 0 {
				return n
			}
		}
	}
	return 80
}

func labelOr(labels map[string]string, key, def string) string {
	if labels != nil {
		if v := labels[key]; v != "" {
			return v
		}
	}
	return def
}

// RestartContainer is a thin helper for tests.
func (a *Agent) RestartContainer(ctx context.Context, id string) error {
	return a.docker.ContainerRestart(ctx, id, container.StopOptions{})
}
