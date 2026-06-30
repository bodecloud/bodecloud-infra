***

name: bolabaden.org Infrastructure
last_updated: 2026-06-30

***

# bolabaden.org Infrastructure Strategy

## What problem this strategy is actually solving

The problem is not "how do we run containers."

The real problem is:

> how do we make multiple ordinary Docker nodes behave like a resilient, operator-readable system without immediately paying the full tax of Swarm, Kubernetes, or another heavyweight orchestrator?

That problem breaks down into several linked pressures:

- requests may land on any node
- the relevant service may not be on the node that received the request
- routing, auth, middleware, and health semantics should survive failover
- secrets, env, and placement truth should stop depending on operator memory
- stateful services should be treated honestly instead of being lumped into the same story as stateless web routing

The manual synchronization bottleneck is the heart of the issue. The repo keeps paying tax in:

- routing drift
- service placement ambiguity
- env and secret sprawl
- partial failover behavior
- too much hidden operator knowledge

## Guiding policy

The repo's strategy is to stay Compose-first for as long as that remains the least painful honest choice, while deliberately adding only the extra coordination layers that solve a real failure mode.

That means:

- keep Compose readable and operator-facing
- reject fake HA language
- improve ingress, placement, and convergence before pretending the system has solved state
- keep multiple stronger-control-plane options open until one clearly earns its tax

This is not anti-automation.

It is anti-premature-platform-ideology.

The goal is not to avoid a control plane forever. The goal is to avoid adopting a heavier one before the repo can justify exactly what pain it removes.

## What this strategy believes

### 1. Compose is still the baseline human contract

The repo is built around the root [`docker-compose.yml`](docker-compose.yml) and included fragments under [`compose/`](compose). That remains the clearest live surface for the operator.

### 2. Compose alone is not enough for honest multi-node behavior

The stack already shows the need for stronger layers around:

- placement truth
- route generation
- health-aware failover
- DNS and ingress convergence
- bootstrap and node join consistency

### 3. The next useful layer is not automatically Kubernetes

The strategy keeps several futures open because the repo is still working out where extra complexity should live:

- Compose plus helpers and registry-driven routing
- failover-agent or Constellation-style convergence logic
- OpenSVC or similar infra-grade HA
- Nomad or a comparable lighter scheduler
- k3s or Kubernetes only when the capability threshold is clearly crossed

### 4. State is the honesty boundary

The strategy treats ingress failover and stateful resilience as different categories.

It is acceptable to improve:

- HTTP request continuity
- local-first then peer-fallback routing
- DDNS and ingress survivability

without pretending that Redis, MongoDB, Postgres, RabbitMQ, or other state-bearing systems are thereby solved.

## Coherent action

The coherent action for this repo is to reduce operator burden in the order that most directly supports the actual dream.

### Track 1: Placement and convergence truth

Build and verify the smallest honest current-state layer that answers:

- what runs where
- what can fail over where
- what should restart locally
- what should be picked up elsewhere

This includes work around:

- a real placement-truth surface such as the repeatedly referenced `services.yaml`
- safer route generation
- sync and convergence logic
- explicit separation between live placement truth and planned placement models

Why this matters:

Without placement truth, "multi-node" remains partly folklore.

### Track 2: Ingress continuity without fake guarantees

Harden the any-node request model so the stack gets better at:

- local-first serving
- peer forwarding
- middleware continuity
- auth continuity
- health-aware route survival

This includes work around:

- Traefik routing semantics
- replacement of known-brittle failover generation paths like `docker-gen-failover`
- Cloudflare node-entry behavior
- private peer reachability

Why this matters:

This is the nearest honest resilience layer for the repo.

### Track 3: Bootstrap and environment convergence

Reduce the amount of hidden memory required to bring up a node or keep one aligned.

This includes work around:

- bootstrap consistency
- secret and env hydration
- node identity and mesh join
- reducing tacit operator steps

Why this matters:

If joining or recovering a node still depends on private memory, the system keeps carrying human SPOFs even when containers look distributed.

### Track 4: Honest stateful resilience

Treat state-bearing systems as first-class architecture problems rather than as extensions of the HTTP routing story.

This includes work around:

- replica and election semantics
- storage topology
- failover correctness
- client reconnect behavior
- L4 exposure boundaries

Why this matters:

This is where the repo either becomes genuinely trustworthy or quietly starts lying.

### Track 5: Control-plane justification and escalation

Keep stronger orchestration options alive, but judge them by explicit pain removed rather than industry fashion.

This includes work around:

- documenting promotion thresholds
- clarifying whether missing semantics belong in metadata, active agent logic, or a stronger scheduler
- preserving optionality without hiding indecision

Why this matters:

The repo is not choosing products. It is choosing where complexity should live.

## Success metrics

These metrics are only useful if they map to the actual anti-SPOF goal.

- **Request continuity under node loss**
  Percent of representative HTTP requests that still succeed when an entry node or serving node fails.

- **Failover correctness, not just detection**
  Time from service or node failure to restored healthy routing, with middleware/auth behavior still preserved.

- **Placement truth confidence**
  How often the operator can answer "what runs where right now?" from authoritative artifacts instead of inference or memory.

- **Configuration and secret convergence time**
  Time required for updated configuration or secret material to become consistent across the nodes that need it.

- **Bootstrap memory removed**
  Number of hidden manual decisions still required to turn a fresh node into a trustworthy participant.

- **Stateful recovery realism**
  Verified service-specific recovery behavior for Redis, MongoDB, Postgres, RabbitMQ, and similar systems, measured separately from ingress success.

## Milestones

### Near-term

- make the failover and route-generation story less brittle than the current `docker-gen-failover` path
- tighten placement-truth and convergence documentation until the next implementation surface is unambiguous
- keep proving the difference between DNS-level survivability, HTTP-level failover, and stateful correctness

### Medium-term

- establish a genuinely trustworthy placement and convergence layer
- reduce node bootstrap and recovery memory burden materially
- prove representative failover behavior on real service classes rather than only on architecture diagrams

### Longer-term

- decide which workloads remain Compose-first
- decide which workloads require stronger scheduler or HA semantics
- justify any move toward OpenSVC, Nomad, k3s, Kubernetes, or deeper Constellation ownership with concrete pain removed

## What this strategy is explicitly not doing

- claiming the repo already has universal multi-node failover
- treating DNS failover as equivalent to service resilience
- treating stateless routing improvements as proof of stateful HA
- committing the whole repo to Kubernetes or Swarm by ideology
- pretending the future control surface is already fully settled

## One-line strategy

Keep the stack Compose-first and operator-readable while deliberately adding the smallest honest coordination layers needed to remove SPOFs, survive request-path failures, and reduce hidden operator burden.

## Short message

`bolabaden-infra` is trying to build a personal-cloud middle ground: stronger than raw multi-node Compose sprawl, lighter than full orchestrator religion, and honest about where routing success ends and real HA begins.
