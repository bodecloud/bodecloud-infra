# Brainstorm: Module 4 — Service Failover & Auto-Redeploy (Next-Gen)

> **Status**: Code landed — **DinD CI proves Tier-A**; do **not** claim full Track 2 / middleware HA closed\
> **Date**: 2026-06-04 (updated 2026-07-18)\
> **Feasibility**: Conditional (see `/tmp/compound-engineering/failover-feasibility/` or repo notes below)\
> **Topic**: Replacing the broken `docker-gen-failover` with a registry-backed Go `failover-agent`.

## Live pointer

| Piece | Path |
|---|---|
| Agent daemon | [`infra/cmd/failover-agent/main.go`](../../infra/cmd/failover-agent/main.go) |
| Registry + Traefik writer | [`infra/failover/`](../../infra/failover/) |
| Compose service | [`compose/docker-compose.failover-agent.yml`](../../compose/docker-compose.failover-agent.yml) |
| Example registry | [`placement/services.yaml.example`](../../placement/services.yaml.example) |
| Env contract | `.env.example` (`FAILOVER_*`) |
| 4-VM dual-DNS CI | [`arbitrary-scripts/failover-ci/`](../../arbitrary-scripts/failover-ci/) |
| Module 5 CF DDNS (deferred) | [`docs/brainstorms/20260718-module5-cloudflare-ddns-followup.md`](20260718-module5-cloudflare-ddns-followup.md) |

`docker-gen-failover` has been removed from Compose. The agent writes `${CONFIG_PATH}/traefik/dynamic/failover-fallbacks.yaml` from the placement registry and never drops routes on crash/unhealthy.

**4-VM CI driver:** [`arbitrary-scripts/failover-ci/`](../../arbitrary-scripts/failover-ci/) — CoreDNS (not Cloudflare), dual DNS (MagicDNS → CoreDNS → Google), Headscale admitted SPOF, heterogeneous whoami/ci-probe placement, production DNS parity, seeded random chaos.

Run `./run-all.sh` (or GHA `Failover Mesh CI` schedule/dispatch) before claiming Tier-A ingress HA. Prove scripts: `prove-matrix`, `prove-dns`, `prove-production-dns`, `prove-failover`, `prove-chaos-random`, `prove-headscale-spof`, `prove-module5-ddns`.

**Sole file owner:** do not run `scripts/osvc_ingress_sync.py` against the same `failover-fallbacks.yaml` on the same node.

**Peer replica ensure** (`FAILOVER_REPLICA_ENSURE`) is enabled on **CI main (`ci-node1`) only**; prod `.env.example` stays `false` until Tailscale Docker API is standard. Image-based ensure works without local ContainerID.

## Feasibility verdict (2026-07-18, updated 2026-07-23)

**Conditional → CI-ready for Tier-A scope** against STRATEGY Track 2 / `.github/copilot-instructions.md`:

| Met | Gap |
|---|---|
| Registry-backed routes survive crash | Middleware/auth continuity not proven (deferred) |
| Local + `service.node.domain` peer URLs | Prod dual-writer (agent + osvc) still open |
| Compose-first, no Swarm | Self-hosted / GHA DinD mesh not yet green on CI |
| Tier-A DinD runtime proof (matrix, failover, chaos, HS SPOF) | Stateful HA, TCP failover still out of scope |
| Compose ensure allowlist for bolabaden + Autokuma on peers | ExportContainerConfig fallback for non-allowlist services |

## Verification — DinD CI (required before claiming Tier-A ingress HA)

Run [`arbitrary-scripts/failover-ci/run-all.sh`](../../arbitrary-scripts/failover-ci/run-all.sh) or GHA **Failover Mesh CI** (schedule / `workflow_dispatch`). All proves must pass, including Headscale SPOF (`prove-headscale-spof.sh`).

Production whoami kill drill (below) remains recommended before claiming prod Track 2 closure.

## 1. Diagnosis (The "Why")

### The Failure of `docker-gen`

* **The Bug**: When a container stops/dies, `docker-gen` excludes it from the template context, even with `-include-stopped`. This results in the deletion of the Traefik route, which is the exact opposite of "failover."
* **The Gap**: There is no mechanism to "hand off" a service to a peer node if the local instance stays down.
* **The Constraint**: We MUST maintain the "No Orchestrator" philosophy. No K8s control plane, no Swarm.

## 2. Approach (The "How") — v1 shipped shape

1. **Persistence**: Traefik routes come from `placement/services.yaml`; crash only changes status.
2. **Ordered servers**: local `http://svc:port` first + peer `https://svc.peer.$DOMAIN` (osvc-aligned; healthCheck drops dead).
3. **Optional replicas**: `FAILOVER_REPLICA_ENSURE=true` on main for HA-eligible services (default off).
4. **Stateful opt-out**: `failover.replica=false` or deny-list (`redis`, `mongo*`, `postgres*`, `rabbitmq`, `headscale*`).

Constellation gossip remains a future upgrade path; v1 does not require it.

## 3. Product Pressure Test (still open)

* **Q1 (Conflict)**: Split-brain after main recovers — acceptable for stateless HA; singletons use `failover.replica=false`.
* **Q2 (Security)**: Peer secret availability still depends on Module 1 secret sync.
* **Q3 (Storage)**: Stateful volume HA is out of scope for this agent.

***

*Originally Phase 1.1 ce-brainstorm; Compose `failover-agent` code landed under Conditional feasibility.*
