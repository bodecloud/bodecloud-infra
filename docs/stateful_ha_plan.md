### Stateful HA (Zero-SPOF) Plan

This page is a planning document.
It describes the minimum kinds of topology truth the repo would need before it
could honestly claim low-SPOF behavior for stateful systems.

It is not proof that the current runtime already provides this.

## What this page is and is not allowed to prove

This legacy planning page is allowed to:

- explain why stateful HA is a different category from HTTP fallback
- describe the minimum replication, quorum, and storage truths the repo would
  need before using stronger HA language
- preserve the user's refusal to call reachability alone stateful resilience

This page is not allowed to:

- imply that current Compose fragments already deliver zero-SPOF stateful
  behavior
- treat L4 routing experiments as proof of datastore correctness
- let topology plausibility masquerade as replicated-write safety
- narrate Redis, Mongo, or Postgres plans as if the underlying replication
  contracts are already proven in this repo

## What still does not count as stateful HA progress here

The following still do not count as honest stateful progress:

- moving a container between nodes
- routing TCP to a reachable backend
- placing a proxy in front of a single datastore
- having a replication topology that has not been drilled under failure
- making node-local volumes feel interchangeable through documentation tone

Stateful HA only becomes real when the write path, failover path, and data
survival path are all owned strongly enough to survive a bad day.

## Strongest honest current answer

The strongest honest current answer is that this page correctly identifies the
minimum classes of truth required before the repo can stop lying about
stateful resilience. It remains a plan. It does not prove that the current
priority implementation already owns replicated datastore truth, quorum
survival, or shared-volume correctness.

You can’t get “zero SPOF including stateful services” by *only* moving containers between nodes.

Stateful HA requires **replication + quorum** (or replicated block storage) so that losing one node does not lose the data and does not stop writes.

This document is the pragmatic plan for this repo.

---

### 0) Ingress reality check (TCP vs hostname)

- **HTTP(S)** can be routed by hostname (Host header/SNI), and the repo already
  experiments with dynamic Traefik failover config generation.
- **Plain TCP** (like `redis://…`) cannot be routed by hostname unless you terminate TLS and use SNI.
  - So for TCP we load-balance by **port** (ex: 6379) and rely on DNS/LB to land you on any node.
  - See `scripts/osvc_l4_sync.py` + `compose/docker-compose.l4-ingress.yml`.

That does **not** mean the repo already has full TCP failover correctness.
It means the repo already knows that plain TCP and stateful HA cannot be faked
with the same story used for HTTP.

---

### 1) Redis (recommendation: Redis Sentinel + HAProxy master routing)

**Goal**: `redis://redis.<node>.bolabaden.org:6379` and
`redis://redis.bolabaden.org:6379` connect to the *current master* once this
topology exists and is actually verified.

Minimum topology (3 nodes):
- 1 Redis master
- 2 Redis replicas
- 3 Sentinels (one per node)

Failover:
- Sentinels elect a new master automatically.
- HAProxy performs a protocol-aware check and forwards clients to the master.

Notes:
- This is operationally simpler than Redis Cluster for “single logical Redis” usage.
- For strict zero-SPOF you want 3 nodes so Sentinel quorum survives 1 node loss.

---

### 2) MongoDB (recommendation: Replica Set)

Minimum topology (3 nodes):
- 1 primary
- 2 secondaries

Client behavior:
- Drivers handle failover if they can see multiple members.

That conditional matters.
Mongo resilience is not created by putting Traefik or DNS in front of a single
Mongo container.

Two ways to make it “one hostname”:
- **Best**: Mongo SRV records (`mongodb+srv://…`) with records generated from cluster membership.
- **Acceptable**: Provide a seed list in the connection string (`mongodb://host1,host2,host3/?replicaSet=rs0`).

We can generate SRV + A records from OpenSVC membership, but it requires Cloudflare DNS API automation.

---

### 3) Postgres (if you add it later)

Options:
- **Patroni + etcd** (strong HA, more moving parts)
- **Repmgr** (simpler)
- **Citus/Spilo** variants (opinionated)

Recommendation for “zero bullshit”:
- Use a managed Postgres if possible, or keep Postgres to a single node until the rest is stable.

---

### 4) Files/Volumes (the hardest part)

Anything that uses bind-mounts like `${CONFIG_PATH}/...` becomes node-local state.

To make it HA you need:
- **Distributed filesystem** (CephFS, GlusterFS) or
- **Replicated block devices** (DRBD) + a filesystem + single active writer at a time.

For this stack, the minimal workable path is:
- Keep most app config volumes node-local initially.
- Make only the truly critical datastores replicated (Redis/Mongo).
- Promote shared volumes later once ingress + scheduling are stable.

This is one of the most important anti-fantasy sections in the repo.

If volumes stay node-local, then some forms of node interchangeability are
still impossible no matter how polished the ingress story sounds.

---

### 5) What we will do next in this repo

- Introduce opt-in labels for L4 services (Redis, Mongo, etc.) so `scripts/osvc_l4_sync.py` can generate per-port HAProxy frontends/backends dynamically.
- Add a Cloudflare DNS automation script to keep:
  - per-node wildcard records (`*.node.domain`)
  - optional global wildcard (`*.domain`) via LB/VIP
  - optional Mongo SRV records for replica set discovery
