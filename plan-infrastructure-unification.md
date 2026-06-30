# Infrastructure Unification Plan

This file used to be a long essay about the pain of scaling multi-node
infrastructure.
Some of that diagnosis was real.
As a working document, it was too indirect and too easy to overread.

This rewrite keeps the core point and removes the theatrical sprawl.

## The actual problem

The repo does not mainly suffer from a lack of containers.

It suffers from too many partial truths spread across different surfaces:

- Compose definitions
- include selection
- env and secret material
- DNS and ingress behavior
- route-generation helpers
- node placement assumptions
- stateful-service behavior
- operator memory

That fragmentation is what "infrastructure unification" is supposed to fix.

## What unification means here

In `bolabaden-infra`, unification does **not** mean:

- invent one magical control surface
- flatten everything into a giant YAML or giant Go program
- assume Kubernetes, Nomad, OpenSVC, or Constellation automatically solves the
  whole stack

It means:

> reduce the amount of hidden operator reconstruction required to understand
> what runs where, how requests are routed, what fails over, and what only
> appears healthy.

That is the real job.

## The dream behind the word

The user is clearly asking for a middle layer that feels more coherent than raw
multi-node Compose, but less doctrinaire than immediately adopting a heavyweight
orchestrator.

The dream is roughly:

- any node can receive public traffic
- if the service exists locally, serve it locally
- if the service lives elsewhere, forward to the peer that hosts it
- keep routing, middleware, and auth semantics coherent across that handoff
- stop depending on private operator memory for placement truth
- stop pretending stateful HA is solved just because a web request still loads

That is the unification target.

## The five truths that actually need unifying

### 1. Convergence truth

Do the right nodes have the right:

- compose fragments
- secret material
- env values
- generated route data
- support files

If not, failover remains a performance.

### 2. Placement truth

Can the system answer:

- what runs where right now
- what host should answer this request
- what is eligible to take over

If not, the operator is still the hidden scheduler.

### 3. Ingress truth

Can any node receive a request and decide honestly whether to:

- serve locally
- forward to a peer
- fail because the service is actually unavailable

If not, the "any-node" story is still partly imaginary.

### 4. Failover truth

When a backend or node disappears:

- does the route survive
- is there a real next candidate
- does the system converge after handoff

If not, the repo still has decorative HA.

### 5. State truth

When writes and durable topology matter:

- who is primary
- who can safely take over
- what replication exists
- what reconnect behavior is expected

If this is fuzzy, the system may look resilient at L7 while still being
architecturally fragile.

## Why the repo keeps exploring multiple futures

Different coordination layers help with different missing truths.

### Compose plus helpers

Good for:

- readable authoring
- low migration shock
- small targeted improvements

Weak when:

- placement truth remains implicit
- route survival depends on brittle generators
- stateful behavior needs stronger coordination

### OpenSVC-style HA

Good for:

- infra-grade failover
- service takeover semantics
- stronger continuity at the node/service boundary

Still not a free answer for:

- developer workflow clarity
- full placement ergonomics
- stateful correctness across every backend

### Nomad

Good for:

- scheduler-backed placement truth
- relocation semantics
- clearer runtime state

Cost:

- stronger control-plane worldview
- another platform to operate

### k3s / Kubernetes

Good for:

- rich control-plane semantics
- workload scheduling
- broad ecosystem support

Cost:

- a much heavier operational and storage model
- easy over-adoption when only part of the stack truly needs it

### Constellation / custom agent layers

Good for:

- preserving repo-specific behavior
- experimenting with a Compose-adjacent control plane
- owning exactly the semantics the repo cares about

Risk:

- reinventing a control plane without enough proof
- drifting into "claimed complete" territory before runtime validation exists

## Current honest reading

The current repo already proves:

- Compose is still the primary human contract
- the user wants any-node entry with local-first then peer-forward behavior
- the repo no longer trusts raw Compose alone to express all runtime truth

The current repo does **not** yet prove:

- universally trusted placement truth
- broadly validated peer-aware failover
- stateful anti-SPOF behavior across all critical backends
- one settled control-plane choice

## Practical plan

The most useful unification work remains:

1. Make placement truth explicit.
2. Make route generation depend on that truth instead of folklore.
3. Reduce bootstrap and secret-convergence drift.
4. Separate ingress resilience from stateful resilience in both code and docs.
5. Promote only the semantics that have earned a stronger platform.

## Where the fuller repo-grounded version lives

The rendered knowledgebase version of this topic is here:

- [`knowledgebase/research/plan-infrastructure-unification.md`](knowledgebase/research/plan-infrastructure-unification.md)

That page is the better maintained explanatory surface.
This root file now exists to keep a concise, honest, operator-usable statement
near the repo root.
