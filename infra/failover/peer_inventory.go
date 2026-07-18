package failover

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
)

// IngestPeerInventory connects read-only to each peer Docker API and upserts
// HTTP failover candidates so never-local services still get Traefik routes.
// Does not start/stop containers on peers (ReplicaEnsure remains separate).
func (a *Agent) IngestPeerInventory(ctx context.Context) error {
	peers := filterPeers(a.cfg.PeerHosts, a.cfg.LocalNode)
	if len(peers) == 0 {
		return nil
	}
	if a.log == nil {
		a.log = log.Default()
	}

	var firstErr error
	for _, peer := range peers {
		if err := a.ingestPeerHost(ctx, peer); err != nil {
			a.log.Printf("peer inventory %s: %v", peer, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (a *Agent) ingestPeerHost(ctx context.Context, peer string) error {
	rdc, err := NewRemoteDockerClient(peer, a.cfg.RemoteDockerPort, a.cfg.RemoteDockerTLS)
	if err != nil {
		return err
	}
	cli, err := rdc.CreateClient(ctx)
	if err != nil {
		return err
	}
	defer cli.Close()

	list, err := cli.ContainerList(ctx, types.ContainerListOptions{All: false})
	if err != nil {
		return fmt.Errorf("list containers: %w", err)
	}

	seen := make(map[string]bool)
	for _, c := range list {
		if !IsHTTPFailoverCandidate(c.Labels) {
			continue
		}
		inspect, err := cli.ContainerInspect(ctx, c.ID)
		if err != nil {
			a.log.Printf("peer %s inspect %s: %v", peer, shortID(c.ID), err)
			continue
		}
		name, err := a.upsertFromPeerInspect(peer, inspect)
		if err != nil {
			a.log.Printf("peer %s upsert: %v", peer, err)
			continue
		}
		if name != "" {
			seen[name] = true
		}
	}

	// Peer services previously running on this peer but missing now → crashed on peer.
	for _, e := range a.registry.List() {
		if seen[e.Name] {
			continue
		}
		a.markPeerMissing(e.Name, peer)
	}
	return nil
}

func (a *Agent) markPeerMissing(name, peer string) {
	e, ok := a.registry.Get(name)
	if !ok || e.Nodes == nil {
		return
	}
	np, ok := e.Nodes[peer]
	if !ok || np == nil {
		return
	}
	if np.Status != StatusRunning && np.Status != StatusUnhealthy {
		return
	}
	np.Status = StatusCrashed
	np.LastSeen = time.Now().UTC()
	e.Status = aggregateStatus(e, a.cfg.LocalNode)
	e.UpdatedAt = time.Now().UTC()
	a.registry.Upsert(e)
}

// upsertFromPeerInspect merges a peer container into the registry without
// overwriting local ContainerID / local node placement.
func (a *Agent) upsertFromPeerInspect(peer string, inspect types.ContainerJSON) (string, error) {
	labels := inspect.Config.Labels
	if !IsHTTPFailoverCandidate(labels) {
		return "", nil
	}

	name := ComposeServiceName(labels, inspect.Name)
	port := ResolvePort(labels, firstExposedPort(inspect))
	status := StatusRunning
	if inspect.State != nil {
		if inspect.State.Health != nil && inspect.State.Health.Status == "unhealthy" {
			status = StatusUnhealthy
		}
		if !inspect.State.Running {
			if inspect.State.ExitCode == 0 && !inspect.State.OOMKilled {
				status = StatusIntentionallyStopped
			} else {
				status = StatusCrashed
			}
		}
	}

	existing, ok := a.registry.Get(name)
	var entry *ServiceEntry
	if ok {
		// Shallow copy fields we may mutate
		entry = &ServiceEntry{
			Name:            existing.Name,
			Port:            existing.Port,
			HealthPath:      existing.HealthPath,
			HealthInterval:  existing.HealthInterval,
			HealthTimeout:   existing.HealthTimeout,
			ReplicaEligible: existing.ReplicaEligible,
			ComposeService:  existing.ComposeService,
			Image:           existing.Image,
			ContainerID:     existing.ContainerID,
			Labels:          existing.Labels,
			Middlewares:     append([]string(nil), existing.Middlewares...),
			Nodes:           cloneNodes(existing.Nodes),
			Status:          existing.Status,
		}
		if port > 0 {
			entry.Port = port
		}
		if mws := MiddlewaresFromLabels(labels); len(mws) > 0 && len(entry.Middlewares) == 0 {
			entry.Middlewares = mws
		}
		if entry.HealthPath == "" {
			entry.HealthPath = HealthPathFromLabels(labels)
		}
	} else {
		entry = &ServiceEntry{
			Name:            name,
			Port:            port,
			HealthPath:      HealthPathFromLabels(labels),
			HealthInterval:  labelOr(labels, "kuma.healthcheck.interval", "15s"),
			HealthTimeout:   labelOr(labels, "kuma.healthcheck.timeout", "5s"),
			ReplicaEligible: IsReplicaEligible(name, labels),
			ComposeService:  name,
			Image:           inspect.Config.Image,
			Labels:          labels,
			Middlewares:     MiddlewaresFromLabels(labels),
			Nodes:           map[string]*NodePlacement{},
			Status:          status,
		}
	}

	if entry.Nodes == nil {
		entry.Nodes = map[string]*NodePlacement{}
	}
	entry.Nodes[peer] = &NodePlacement{
		Status:   status,
		LastSeen: time.Now().UTC(),
		Priority: 2,
	}

	// Aggregate Status: prefer local running; else any peer running.
	entry.Status = aggregateStatus(entry, a.cfg.LocalNode)

	a.registry.Upsert(entry)
	return name, nil
}

func cloneNodes(in map[string]*NodePlacement) map[string]*NodePlacement {
	out := make(map[string]*NodePlacement, len(in))
	for k, v := range in {
		if v == nil {
			continue
		}
		cp := *v
		out[k] = &cp
	}
	return out
}

func aggregateStatus(e *ServiceEntry, localNode string) ServiceStatus {
	if e.Nodes != nil {
		if np, ok := e.Nodes[localNode]; ok && np != nil {
			if np.Status == StatusRunning || np.Status == StatusUnhealthy {
				return np.Status
			}
		}
		for _, np := range e.Nodes {
			if np != nil && (np.Status == StatusRunning || np.Status == StatusUnhealthy) {
				return np.Status
			}
		}
		for _, np := range e.Nodes {
			if np != nil {
				return np.Status
			}
		}
	}
	return e.Status
}
