# Constellation Agent Roadmap

This roadmap is now written to distinguish between:

- **implemented code direction**
- **plausible but unverified capability**
- **future work**

That distinction matters because older versions mixed "exists in the tree" with
"operationally proven."

## What Constellation is trying to become

Constellation is the repo's attempt to build a Docker-native coordination layer
that can reduce:

- hidden placement knowledge
- brittle node-local routing assumptions
- manual failover reconstruction
- some of the operator rent imposed by multi-node Compose

## Implemented direction visible in the tree

These items appear to have real implementation effort behind them and are fair
to discuss as coded project direction:

- gossip-oriented cluster state exchange
- Raft-backed coordination for leadership-style operations
- dynamic Traefik configuration generation
- Cloudflare DNS control logic
- service health monitoring
- Tailscale-aware peer discovery assumptions
- stateful service handling attempts
- install, verify, and deployment helper surfaces

## Direction that still needs stronger verification

These may be partly coded or locally testable, but should not be narrated as
operationally settled without stronger evidence:

- universal automatic failover
- zero-SPOF behavior
- production-grade stateful HA
- reliable wrong-node service resolution
- full repo-wide replacement of Compose-managed truth

## Near-term roadmap

### 1. Clarify runtime authority

Constellation needs a crisper answer to:

- what truth it owns
- what the root Compose runtime still owns
- where coexistence ends and takeover begins

### 2. Prove representative failover paths

Focus on evidence, not slogans:

- a small set of stateless HTTP services
- wrong-node entry behavior
- route survival after backend loss
- health-aware convergence after recovery

### 3. Separate stateless and stateful proof

Do not let ingress success stand in for database resilience.

Needed:

- explicit stateful-service test plans
- topology-specific recovery expectations
- honest limits where stateful HA is still incomplete

### 4. Reduce documentation overclaim

The docs should keep marking the difference between:

- intended architecture
- coded implementation
- verified behavior

### 5. Decide promotion thresholds

Constellation should justify itself against other options by explicit pain
removed:

- placement truth
- convergence truth
- ingress truth
- failover truth

If it cannot do that cleanly, the repo should keep other coordination options
open.

## Medium-term roadmap

- stronger operational verification tooling
- clearer cluster bootstrap story
- lower secret and env convergence ambiguity
- better explanation of direct-node vs global routing semantics
- more explicit boundaries around what services are safe candidates for
  Constellation-managed failover

## Longer-term decision point

The long-term question is not "can more features be added?"

It is:

> does Constellation earn the right to become a trusted control layer for the
> repo's real multi-node dream, or should some responsibilities be promoted into
> another platform?

That remains open.
