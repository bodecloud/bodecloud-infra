# High Availability Kubernetes Infrastructure

## Reality check

This file is an HA target sketch from the `garden.io/ha-k8s/` exploration.
It should not be read as current repo truth, finished cluster proof, or
verified zero-SPOF reality.

The root priority implementation for `bolabaden-infra` is still the
Compose-first runtime built around
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

This page is useful because it records what the Kubernetes-shaped answer was
trying to buy:

- control-plane survivability
- storage replication
- ingress redundancy
- multi-node scheduling

It is not useful as proof that those things were actually achieved for the
repo as a whole.

For the assimilated reading of the Garden and k3s branch, see:

- [`../../knowledgebase/research/garden-k3s-exploration-evidence.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/garden-k3s-exploration-evidence.md)

## What this file can and cannot establish

This file can establish:

- which HA capabilities the branch intended Kubernetes to own
- which components were considered necessary for a serious HA cluster
- which topology assumptions the branch was willing to make

This file cannot establish:

- that the cluster was deployed exactly as described
- that the described storage stack matched the final branch implementation
- that ingress continuity was proven under real node loss
- that stateful services inherited honest HA just because the cluster design
  named replicated storage
- that this path superseded the repo's Compose-first direction

That distinction matters because a page like this can sound extremely complete
while still being mostly architecture intent.

## What this branch was trying to design

The branch aimed to describe a Kubernetes answer to the repo's anti-SPOF
pressure:

- multiple control-plane nodes instead of one sacred controller
- distributed workers instead of one sacred runtime box
- replicated storage instead of one sacred disk path
- replicated ingress instead of one sacred edge process

In other words, it was trying to buy the kind of system-owned redundancy the
repo cannot get from plain multi-node Compose alone.

That is the real value of preserving this file.
It shows what the Kubernetes-shaped answer was *supposed* to earn, even though
the repo has not promoted that answer into present authoritative runtime truth.

## Proposed cluster topology

### Control plane nodes

- `micklethefickle.bolabaden.org`
- `cloudserver1.bolabaden.org`
- `cloudserver2.bolabaden.org`

### Worker nodes

- `cloudserver3.bolabaden.org`
- `blackboar.bolabaden.org`

These names should be read as branch planning inputs, not as current verified
cluster membership.

## Intended HA components

### 1. Control plane HA

The branch intended:

- `kube-apiserver` load-balanced across three control-plane nodes
- `etcd` quorum across three nodes
- `kube-scheduler` with leader election
- `kube-controller-manager` with leader election

This is a classic "remove the sacred controller" move.
What it does **not** prove is that the repo validated these roles under real
failure on the current branch state.

### 2. Storage HA

The branch described:

- replicated CSI-backed storage
- HA storage classes
- automatic volume failover
- replication-based data protection

This is the section most likely to over-promise if read lazily.
Replicated storage design does not by itself prove:

- correct stateful application failover
- correct client reconnection semantics
- no split-brain risk for higher-level services
- that all critical repo backends were successfully mapped onto this model

### 3. Networking HA

The branch intended:

- Calico for multi-node networking
- optional service-mesh layers
- MetalLB or cloud load balancing for external access

That tells the reader which kinds of networking answers were being explored.
It does not prove the repo solved the more application-level questions the
user keeps caring about, such as preserving request meaning when traffic lands
on the wrong node.

### 4. DNS HA

The branch intended:

- multiple CoreDNS replicas
- anti-affinity
- DNS failover to healthy replicas

This can strengthen cluster-internal resolution.
It still should not be confused with end-to-end proof that public any-node
entry and service preservation are solved.

### 5. Ingress HA

The branch intended:

- multiple ingress-controller replicas
- traffic handoff across healthy nodes

This is closer to the repo's real anti-SPOF pressure than many other sections.
Even so, naming ingress HA is not the same thing as proving:

- wrong-node HTTP success
- auth continuity after handoff
- middleware continuity after handoff
- route persistence when a backend disappears during the exact failure being
  tested

### 6. Application HA

The branch intended:

- anti-affinity
- disruption budgets
- health checks
- optional auto-scaling

These are useful cluster-level controls.
They are not the same as repo-wide proof that all service classes now have
honest failover semantics.

## Intended failure scenarios

The branch treated these as target behaviors:

### Node failure

- reschedule pods onto healthy nodes
- retain storage access through replication
- keep services running through surviving replicas

### Control-plane failure

- maintain quorum with two of three control-plane nodes
- preserve API access through remaining nodes

### Storage failure

- fail over to replicated storage
- avoid data loss
- recover automatically when the failed component returns

### Network partition

- let quorum semantics prevent split-brain at the control-plane layer
- continue in the majority partition
- reconcile after healing

These scenarios show what the branch considered necessary for a serious HA
story.
They do not prove the branch achieved or tested them all.

## Historical deployment outline

The original page implied a straightforward deployment flow:

1. initialize cluster
2. deploy storage
3. deploy services

That flow is still useful as a branch artifact, but it should now be read as a
historical operating sketch rather than a verified recipe for current repo
truth.

## Why preserve this file at all?

Because deleting it would hide a real part of the repo's search space.

The repo keeps asking whether stronger anti-SPOF guarantees require a larger
control plane.
This file shows one of the strongest "yes, via Kubernetes" answers that the
branch explored.

That matters even when the answer is not promoted.

What matters just as much is refusing to let that exploration impersonate
present proof.
