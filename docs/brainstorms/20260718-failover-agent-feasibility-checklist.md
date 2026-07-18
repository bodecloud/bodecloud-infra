# Failover-agent feasibility criteria checklist

Generated: 2026-07-18 (inline feasibility review after subagent ResourceExhausted)

## Must-pass (near-term Track 2)

Sources: `STRATEGY.md`, `.github/copilot-instructions.md`, `knowledgebase/architecture/ha-failover-routing.md`, `docs/brainstorms/20260604-failover-agent-exploration.md`

1. Replace brittle `docker-gen-failover` route generation — STRATEGY near-term milestone
2. Route persistence under the failure that triggered fallback (crash/die must not delete routes)
3. Local-first then peer-forward using `service.domain` + `service.node.domain`
4. Placement / current-state truth readable as registry (`services.yaml` idea)
5. Keep Compose-first; no Swarm/K8s requirement
6. Do not treat HTTP routing success as stateful HA proof
7. Prove representative failover (whoami) before claiming success

## Must-not-claim

- Universal multi-node / wrong-node success
- Middleware/auth continuity fully proven
- TCP/stateful HA solved
- DNS plurality == service resilience

## Dependencies docs name

- Peer Traefik + DNS wildcards `*.node.domain`
- Shared/identical edge middleware stack on peers (assumption)
- Peer Docker API reachable for replica ensure (plan prerequisite)
- OpenSVC sync historically also writes `failover-fallbacks.yaml`

## Ambiguities / tensions

| Surface | Approach |
|---|---|
| Compose `failover-agent` | Env peers + registry + file provider |
| `scripts/osvc_ingress_sync.py` | OpenSVC node list; same output file |
| Constellation gossip | Future; plan defers |
| `infra/traefik.go` GenerateFailoverConfig | Writes `.json` sibling, not yaml |

## Scorecard vs current implementation

| Criterion | Score | Notes |
|---|---|---|
| Replace docker-gen | PASS | Removed from coolify-proxy; agent added |
| Route persistence on crash | PARTIAL | Registry keeps entries; **die+exit0 misclassified as crash** |
| Local + peer URL shape | PASS | Matches osvc / DNS contract |
| Placement registry | PARTIAL | Writes `placement/services.yaml`; not yet operator-proven |
| Compose-first / no Swarm | PASS | |
| Stateful opt-out | PASS | Deny-list + labels |
| Always-on peer replicas | FAIL/UNKNOWN | Needs peer `:2375`; uses ExportContainerConfig not compose-up; secrets/networks fragile |
| Auth/middleware continuity | FAIL | File routers omit middlewares (same gap as osvc sync) |
| Dual writer safety | FAIL | Agent + osvc_ingress_sync share `failover-fallbacks.yaml` |
| Whoami runtime proof | FAIL | Unit tests only; brainstorm overclaims "Implemented" |
| Honesty vs STRATEGY | PARTIAL | Scope outs correct; status language too strong |

## Verdict

**Conditional** for Track 2 near-term: direction is correct (registry-backed Traefik writer replacing docker-gen). Not yet criteria-complete due to intentional-stop bug, dual-writer conflict, unproven replicas, missing runtime proof, and middleware gap.
