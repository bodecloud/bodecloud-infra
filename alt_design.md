# Alternative Design Note: Containerized HA Cluster Exploration

This file used to read like a finished HA deployment guide.
That was misleading.

It is **not** the current priority implementation.
It is **not** a validated runtime.
It is **not** proof that `bolabaden-infra` already has true anycast, VIP-based
multi-node failover, or host-independent HA.

The priority implementation for this repo is still the root
[`docker-compose.yml`](docker-compose.yml) plus its active include set.

## What this file actually is

This is a record of one stronger-clustering direction the repo explored while
searching for a middle ground between:

- raw multi-node Docker Compose sprawl
- and full orchestrator capture

The explored ingredients here were:

- Headscale for mesh identity and node connectivity
- CoreDNS for custom routing behavior
- Traefik for ingress and backend selection
- a shared VIP idea for "any node can answer" behavior
- health-aware failover inside a more opinionated cluster shape

That makes this file useful as **design pressure evidence**, not as live ops
guidance.

## Why it was attractive

This design was trying to solve a real pain:

- requests may hit the wrong node
- DNS-only failover is too weak
- operator-maintained upstream lists do not scale well
- the repo wants local-first serving with peer-aware fallback

In other words, the design was trying to buy stronger entry-node independence
than the current Compose-first runtime can honestly guarantee.

## Why it is not the adopted answer

Even if the individual technologies are real, this document overreached in
several ways:

- it spoke as if the cluster already existed
- it implied validated request-driven failover
- it implied the VIP and routing behavior were deployment-ready
- it hid the operational tax of introducing a stronger control surface
- it blurred "interesting architecture direction" with "current repo truth"

That is exactly the kind of documentation drift this repo is trying to stop.

## What this exploration still contributes

This alternative design still matters because it clarifies the user's real
demand:

> any node should be able to receive traffic, serve locally when possible, and
> forward intelligently when not, without lying about resilience.

That demand remains active even if this specific cluster design is not the
chosen path.

## How to read it now

Read this file as:

- a stronger-cluster exploration
- evidence that the repo has seriously considered more opinionated HA layers
- a contrast against the current Compose-first approach

Do **not** read it as:

- current deployment instructions
- a proven HA recipe
- an adopted architecture decision

## Where to look instead

For the current repo-grounded interpretation, use:

- [`README.md`](README.md)
- [`STRATEGY.md`](STRATEGY.md)
- [`knowledgebase/architecture/compose-first-architecture.md`](knowledgebase/architecture/compose-first-architecture.md)
- [`knowledgebase/architecture/ha-failover-routing.md`](knowledgebase/architecture/ha-failover-routing.md)
- [`knowledgebase/research/orchestration-research-2026.md`](knowledgebase/research/orchestration-research-2026.md)
