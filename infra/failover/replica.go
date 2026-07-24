package failover

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// EnsureReplicas starts the service on each peer that does not already run it.
// Disabled unless FAILOVER_REPLICA_ENSURE=true. Skips peers marked
// intentionally_stopped. Can create from registry Image when local ContainerID
// is empty (main host without a local instance).
func (a *Agent) EnsureReplicas(ctx context.Context, entry *ServiceEntry) error {
	if !a.cfg.ReplicaEnsure {
		return nil
	}
	if entry == nil || !entry.ReplicaEligible {
		return nil
	}
	if len(a.cfg.PeerHosts) == 0 {
		return nil
	}

	var errs []string
	for _, peer := range a.cfg.PeerHosts {
		peer = strings.TrimSpace(peer)
		if peer == "" || peer == a.cfg.LocalNode {
			continue
		}
		if peerIntentionallyStopped(entry, peer) {
			a.log.Printf("replica ensure skip %s on %s (intentionally_stopped)", entry.Name, peer)
			continue
		}
		if err := a.ensureReplicaOnPeer(ctx, entry, peer); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", peer, err))
			a.log.Printf("replica ensure %s on %s failed: %v", entry.Name, peer, err)
			continue
		}
		if existing, ok := a.registry.Get(entry.Name); ok {
			if existing.Nodes == nil {
				existing.Nodes = map[string]*NodePlacement{}
			}
			existing.Nodes[peer] = &NodePlacement{
				Status:   StatusRunning,
				LastSeen: time.Now().UTC(),
				Priority: 2,
			}
			a.registry.Upsert(existing)
		}
	}
	_ = a.registry.Save()
	if len(errs) > 0 {
		return fmt.Errorf("peer replica errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func peerIntentionallyStopped(entry *ServiceEntry, peer string) bool {
	if entry == nil || entry.Nodes == nil {
		return false
	}
	np, ok := entry.Nodes[peer]
	return ok && np != nil && np.Status == StatusIntentionallyStopped
}

func (a *Agent) ensureReplicaOnPeer(ctx context.Context, entry *ServiceEntry, peer string) error {
	rdc, err := NewRemoteDockerClient(peer, a.cfg.RemoteDockerPort, a.cfg.RemoteDockerTLS)
	if err != nil {
		return err
	}
	remote, err := rdc.CreateClient(ctx)
	if err != nil {
		return fmt.Errorf("connect peer docker: %w", err)
	}
	defer remote.Close()

	running, err := peerHasContainer(ctx, remote, entry.Name, entry.ComposeService)
	if err != nil {
		return err
	}
	if running {
		a.log.Printf("peer %s already has %s", peer, entry.Name)
		return nil
	}

	// Allowlisted Tier-A services: prefer registry/image recreate with generic
	// compose networks (backend/publicnet) instead of ExportContainerConfig, which
	// carries main-project network names (ci-node1_*) that do not exist on peers.
	var cfg *ContainerConfig
	if a.inComposeEnsureAllowlist(entry.Name, entry.ComposeService) {
		cfg = ConfigFromRegistryEntry(entry)
		a.log.Printf("compose-ensure allowlist: minimal config for %s on %s (pull=never path)", entry.Name, peer)
	} else {
		cfg, err = a.buildReplicaConfig(ctx, entry)
		if err != nil {
			return err
		}
	}

	if cfg.Image != "" && !a.cfg.ReplicaPullNever && !a.inComposeEnsureAllowlist(entry.Name, entry.ComposeService) {
		reader, err := remote.ImagePull(ctx, cfg.Image, types.ImagePullOptions{})
		if err != nil {
			a.log.Printf("image pull on %s for %s: %v (continuing)", peer, cfg.Image, err)
		} else {
			_ = reader.Close()
		}
	} else if cfg.Image != "" {
		a.log.Printf("skipping ImagePull for %s on %s (pull=never / allowlist)", cfg.Image, peer)
	}

	id, err := CreateContainerOnRemote(ctx, remote, cfg)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") || strings.Contains(err.Error(), "Conflict") {
			a.log.Printf("container name conflict on %s for %s — attempting start of existing", peer, cfg.Name)
			return startExistingByName(ctx, remote, cfg.Name)
		}
		return err
	}

	if err := remote.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("start on peer: %w", err)
	}
	a.log.Printf("started replica of %s on %s (%s)", entry.Name, peer, shortID(id))
	return nil
}

// buildReplicaConfig prefers exporting the local container; otherwise builds a
// minimal config from registry Image/labels/port.
func (a *Agent) buildReplicaConfig(ctx context.Context, entry *ServiceEntry) (*ContainerConfig, error) {
	if entry.ContainerID != "" && a.docker != nil {
		cfg, err := ExportContainerConfig(ctx, a.docker, entry.ContainerID)
		if err == nil {
			cfg.Name = strings.TrimPrefix(cfg.Name, "/")
			if entry.ComposeService != "" {
				cfg.Name = entry.ComposeService
			}
			return cfg, nil
		}
		a.log.Printf("export local config for %s failed: %v — falling back to image", entry.Name, err)
	}
	return ConfigFromRegistryEntry(entry), nil
}

// ConfigFromRegistryEntry builds a minimal recreate config from placement data.
func ConfigFromRegistryEntry(entry *ServiceEntry) *ContainerConfig {
	name := entry.ComposeService
	if name == "" {
		name = sanitizeName(entry.Name)
	}
	port := entry.Port
	if port <= 0 {
		port = 80
	}
	portStr := fmt.Sprintf("%d/tcp", port)
	labels := map[string]string{}
	for k, v := range entry.Labels {
		labels[k] = v
	}
	if labels["traefik.enable"] == "" {
		labels["traefik.enable"] = "true"
	}
	if labels["com.docker.compose.service"] == "" {
		labels["com.docker.compose.service"] = name
	}
	image := entry.Image
	if image == "" {
		image = "docker.io/traefik/whoami:latest"
	}
	return &ContainerConfig{
		Name:  name,
		Image: image,
		ExposedPorts: nat.PortSet{
			nat.Port(portStr): struct{}{},
		},
		Labels: labels,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Networks: []string{"backend", "publicnet"},
	}
}

func peerHasContainer(ctx context.Context, cli *client.Client, name, composeService string) (bool, error) {
	list, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false, err
	}
	want := []string{sanitizeName(name), sanitizeName(composeService)}
	for _, c := range list {
		candidates := append([]string{}, c.Names...)
		if svc := c.Labels["com.docker.compose.service"]; svc != "" {
			candidates = append(candidates, svc)
		}
		for _, n := range candidates {
			n = sanitizeName(n)
			for _, w := range want {
				if w != "" && n == w {
					return c.State == "running", nil
				}
			}
		}
	}
	return false, nil
}

func startExistingByName(ctx context.Context, cli *client.Client, name string) error {
	list, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	name = sanitizeName(name)
	for _, c := range list {
		for _, n := range c.Names {
			if sanitizeName(n) == name {
				return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
			}
		}
	}
	return fmt.Errorf("container %s not found on peer", name)
}
