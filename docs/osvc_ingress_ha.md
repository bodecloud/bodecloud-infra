### OpenSVC HA Ingress (Dynamic, Zero Hardcoded Nodes/Services)

This page describes an **approach** for OpenSVC-assisted ingress failover.
It should not be read as proof that the current tracked root runtime already
delivers universal any-node success.

This repo already runs **Traefik** for HTTP(S). The missing piece for multi-node
HA is **dynamic failover/load-balancing across nodes** without hardcoding
node/service names.

This document describes the approach implemented in:
- `scripts/osvc_ingress_sync.py` (generates Traefik file-provider config)
- `scripts/osvc_ingress_sync.sh` (wrapper that loads `.env`)

## What this page is and is not allowed to prove

This legacy page is allowed to:

- explain the exact OpenSVC ingress bet this branch is making
- clarify how generated Traefik fallback config is supposed to reduce
  hand-written node/service folklore
- distinguish runtime-derived routing data from static per-node route files

This page is not allowed to:

- claim universal wrong-node success
- imply OpenSVC already governs the live root runtime
- treat generated fallback YAML as proof that peer forwarding survives real
  backend loss
- blur the difference between HTTP route generation and stateful or TCP
  correctness

## What still does not count as OpenSVC ingress proof here

The following still do not count as real closure:

- a fallback file being generated successfully
- peer hostnames appearing in that file
- one successful happy-path handoff
- local container labels being discoverable
- OpenSVC membership being queryable

Generated config can be strategically interesting and still leave the deepest
truths unowned.

## Strongest honest current answer

The strongest honest current answer is that this page describes a serious
attempt to move fallback routing out of handwritten folklore and toward
runtime-derived config. That is real progress. It is still not proof that the
priority Compose-first runtime already preserves wrong-node requests with full
middleware, auth, and backend correctness.

---

### What this enables

What follows is best read as intended behavior of this OpenSVC-based direction,
not as a blanket statement about today’s live Compose-first runtime.

- **Node-scoped hostnames (always hit that node first)**:
  - `https://<service>.<node>.bolabaden.org`
  - DNS resolves to `<node>`; the intended behavior is that Traefik on that
    node can route to a local container **or** fall back to another node that
    has the service.

- **Global hostnames (load-balance/failover across any node running the service)**:
  - `https://<service>.bolabaden.org`
  - DNS/LB must land you on *any* healthy node ingress; the intended behavior
    is that the receiving node can route locally or fall back to other nodes.

---

### DNS requirements (Cloudflare)

To satisfy `service.node.domain` without hardcoding service names:

- **Per-node records**
  - `A <node>.bolabaden.org -> <node public IP>`
  - `A *. <node>.bolabaden.org -> <node public IP>`

Your `cloudflare-ddns` container already supports managing:
- `$TS_HOSTNAME.$DOMAIN`
- `*.$TS_HOSTNAME.$DOMAIN`

To satisfy `service.domain` load balancing, you need **one** of:

- **Option A (Best: zero-SPOF ingress)**: Cloudflare Load Balancer
  - Create a LB for `*.bolabaden.org` (or for selected subdomains) with origins = your nodes.
  - Health checks should hit a known stable endpoint (ex: `https://whoami.<node>.bolabaden.org/`).

- **Option B (Self-hosted VIP)**: keepalived VRRP floating IP + node ingress
  - Requires network that supports a shared VIP.
  - Traefik binds the VIP and fails over with keepalived.

- **Option C (Good enough)**: DNS round-robin A records
  - `A *.bolabaden.org -> <node1 IP>, <node2 IP>, ...`
  - Not strictly “zero spof” because clients can cache a dead node until TTL.

---

### How the Traefik failover config is generated

Run on each node:

```bash
./scripts/osvc_ingress_sync.sh
```

It:
- Reads nodes from OpenSVC: `om node ls --format json`
- Reads Traefik-enabled **HTTP** containers from Docker (`traefik.enable=true` + `traefik.http.*` labels)
- Writes: `${CONFIG_PATH}/traefik/dynamic/failover-fallbacks.yaml`

Traefik is already configured with:
- `--providers.file.directory=/traefik/dynamic/`
- `--providers.file.watch=true`

So updating the file updates routing dynamically.

That is still weaker than proving correct failure behavior.

Dynamic generation alone does not prove:

- that the fallback route survives when the local backend disappears
- that auth and middleware continuity survive the peer hop
- that this path is the current authoritative runtime for the whole repo

---

### TCP (Redis/Mongo/etc.)

Plain TCP can’t be routed by hostname without TLS/SNI.

To get:
- `redis://redis.<node>.bolabaden.org:6379`
- `redis://redis.bolabaden.org:6379`

…you need an L4 load balancer per port (example: HAProxy in `network_mode: host`) and **true datastore HA** (Redis Sentinel/Cluster, Mongo replica set, etc.). This is planned next.

That last point is the real honesty wall:

HTTP failover experimentation is not proof of stateful HA.
