# Brainstorm: Module 4 — Service Failover & Auto-Redeploy (Next-Gen)

> **Status**: Code landed — **runtime whoami proof pending**; do **not** claim Track 2 closed\
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

**4-VM CI driver:** [`arbitrary-scripts/failover-ci/`](../../arbitrary-scripts/failover-ci/) — CoreDNS (not Cloudflare), dual DNS (MagicDNS → CoreDNS → Google), Headscale on two nodes, heterogeneous whoami/ci-probe placement. Run `./run-all.sh` then `prove-dns.sh` + `prove-failover.sh` before claiming Track 2 closed.

**Sole file owner:** do not run `scripts/osvc_ingress_sync.py` against the same `failover-fallbacks.yaml` on the same node.

**Peer replica ensure** (`FAILOVER_REPLICA_ENSURE`) is enabled on **CI main (`ci-node1`) only**; prod `.env.example` stays `false` until Tailscale Docker API is standard. Image-based ensure works without local ContainerID.

## Feasibility verdict (2026-07-18)

**Conditional** against STRATEGY Track 2 / `.github/copilot-instructions.md`:

| Met | Gap |
|---|---|
| Registry-backed routes survive crash (direction correct) | Whoami kill proof not run |
| Local + `service.node.domain` peer URLs | Always-on replicas gated off / fragile ExportContainerConfig |
| Compose-first, no Swarm | Dual-writer conflict documented; agent now sole owner |
| Intentional stop ≠ crash (exit 0) | Middleware copied from labels when present; auth continuity still peer-parity dependent |

## Verification (whoami) — required before claiming success

On main host (`FAILOVER_MAIN_HOST=micklethefickle`) with peers configured:

1. Confirm agent healthy and file present: inspect `failover-fallbacks.yaml` for `whoami` + peer URLs.
2. Crash: `docker kill whoami` → registry status `crashed`, YAML still lists `https://whoami.<peer>.$DOMAIN` → `curl -k https://whoami.$DOMAIN` should succeed via peer.
3. Intentional stop: `docker compose stop whoami` → status `intentionally_stopped`; no peer redeploy.

After step 2 is proven in production, capture the learning with `/ce-compound`.

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
