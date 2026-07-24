# Failover CI — 4-node Traefik peer-forward + dual-DNS (+ Module 5 dry-run)

Hands-off local driver. **Fully functional without GitHub Actions or git push.**

## Honesty contract (read this)

This suite proves **Tier-A Traefik HA** (edge + allowlisted apps), **not** whole-stack HA.

| Claim | Reality |
|---|---|
| Peer “full stack ×4” | **Banned.** Peers run an **HA-critical curated set** (Traefik, whoami, ci-probe, failover-agent, Headscale UI/client as shaped, bolabaden-nextjs, Autokuma) via image sync + `--pull=never`. |
| Dual-primary / HA Headscale | **Banned.** `headscale-server` is a **single-node admitted SPOF**; `prove-headscale-spof.sh` fail-closes when MagicDNS dies with it. |
| Green whoami alone = stack HA | **Banned.** Expanded `prove-failover` requires bolabaden + Autokuma peer URLs and chaos gates. |

## Topology

| VM | Headscale | Probes / Tier-A |
|---|---|---|
| ci-node1 | **server + UI** (+ CoreDNS) | bolabaden + Autokuma; **no** whoami/ci-probe |
| ci-node2 | **UI only** (+ CoreDNS) | **whoami** + bolabaden + Autokuma |
| ci-node3 | no | **whoami** + bolabaden + Autokuma |
| ci-node4 | no | **ci-probe** (+ bolabaden + Autokuma) |

DNS: **CI CoreDNS zone mirrors production Cloudflare semantics** (see `prove-production-dns.sh`):
- **Global multi-A** (`cloudflare-multi-ddns`) → wildcard `*` A records to all Traefik node IPs
- **Node-direct** (`favonia cloudflare-ddns`) → `*.ci-nodeN.$DOMAIN` → that node only
- **MagicDNS** (Headscale) → `prove-dns.sh` hard gate when Tailscale is up
- **Docker embedded DNS** (`127.0.0.11`) → compose service names on each node (`prove-production-dns.sh`)

Live Cloudflare API is **not** a hard gate; Module 5 sync is proven via `--dry-run` + CoreDNS parity (`prove-module5-ddns.sh`). Optional live zone: `CF_LIVE_MULTI_DDNS=1`.

**Chaos:** deterministic kills in `prove-failover.sh` + **seeded random rounds** in `prove-chaos-random.sh` (`CHAOS_ROUNDS`, `CHAOS_SEED`).

## Backends

Auto-detect order:

1. **DinD** — when `/dev/kvm` is missing (nested guest) and Docker works: WARN + 4 privileged `docker:dind` nodes on bridge `failover-ci-net`, Docker API `tcp://0.0.0.0:2375` (private net only)
2. **Multipass** — preferred on bare metal / KVM hosts when available
3. **QEMU/KVM + libvirt** — needs `/dev/kvm`, `virt-install`, `virsh`, `qemu-img`
4. **DinD** again — last resort if Multipass/QEMU tools are missing

Force a backend: `FAILOVER_CI_BACKEND=dind|multipass|qemu`

```bash
sudo snap install multipass
# or: sudo apt install qemu-kvm libvirt-daemon-system virtinst cloud-image-utils
# DinD needs only a working Docker daemon (privileged containers)
```

On nested VMs without KVM, `./provision-vms.sh` should WARN and launch DinD automatically.

## Quick start

```bash
cd arbitrary-scripts/failover-ci
cp env/test.env.example env/test.env
chmod +x *.sh
./validate-local.sh   # no VMs — unit/Module5/builds/CoreDNS smoke
./run-all.sh          # 4-node mesh + expanded proves (DinD / Multipass / QEMU)
# nested guest without KVM auto-selects DinD; or force:
#   FAILOVER_CI_BACKEND=dind ./provision-vms.sh
```

Teardown: `./teardown.sh`

## Compose modes

**Default:** root [`docker-compose.yml`](../../docker-compose.yml) (+ ci-probes / extra_hosts overlays). Success criterion for peers is the **HA-critical curated set**, not every media service healthy.

**Fast debug:** `FAILOVER_CI_MINIMAL=1` uses [`compose/docker-compose.ci-stack.yml`](compose/docker-compose.ci-stack.yml) plus [`compose/docker-compose.ci-tier-a.yml`](compose/docker-compose.ci-tier-a.yml) (Tier-A stubs) and ci-probe overlay.

```bash
FAILOVER_CI_MINIMAL=1 ./compose-up-all.sh   # lean stack
FAILOVER_CI_MINIMAL=0 ./compose-up-all.sh   # root compose; peers get curated critical set
```

## DinD specifics

- **`failover-agent` runs as `user: "0:0"`** — the DinD docker.sock is `root:docker` mode `660`; the image default `nobody` cannot write `failover-fallbacks.yaml`.
- **Tailscale + Headscale are required** on DinD (no soft-skip). Static Tailscale binaries are installed in bootstrap; `provision-mesh.sh` joins all 4 nodes and hard-fails if any node lacks a Tailscale IP. MagicDNS `@100.100.100.100` is a prove-dns hard gate when Tailscale is up; **Headscale down must fail** `prove-headscale-spof.sh`.
- **Traefik plugins:** DinD DNS often cannot reach `plugins.traefik.io`, which disables *all* Traefik plugins and leaves routers referencing `crowdsec@file` disabled (HTTPS 404). [`compose/docker-compose.ci-dind-fixes.yml`](compose/docker-compose.ci-dind-fixes.yml) stubs `crowdsec` / error middlewares as no-op `headers` and pins `DOCKER_API_VERSION=1.47`.
- **Stale host iptables:** Docker `raw` PREROUTING pins container IPs to a bridge. After DinD network recreate, orphan `! -i br-<dead> -j DROP` rules for `ci-node1`'s IP blackhole peer→n1 traffic (ARP still works). `cleanup_stale_dind_bridge_filters` in `lib.sh` removes those for our node IPs during provision/mesh.
- Peer Docker API stays on private-net `:2375` (Tailscale-bound bind is Multipass/QEMU only).
- Full stack creates external `warp-nat-net` inside each DinD before `compose up`.
- **Image sync:** after `ci-node1` is up, [`sync-images-from-main.sh`](sync-images-from-main.sh) streams a **curated** image list to peers (`FAILOVER_CI_IMAGE_MAX_MB=450` by default — do **not** clone ~21GB×4). HA-critical refs (including bolabaden-nextjs + Autokuma) bypass the size cap. Peers then `compose up --pull=never` for the curated service list. Set `FAILOVER_CI_SYNC_IMAGES=0` to skip; `FAILOVER_CI_IMAGE_MAX_MB=0` for unlimited (needs large disk).

## GitHub Actions

Workflow: [`.github/workflows/failover-mesh.yml`](../../.github/workflows/failover-mesh.yml)

| Trigger | Job |
|---|---|
| PR / push (path-filtered) | `validate-local.sh` on `ubuntu-latest` only |
| Weekly schedule | DinD `./run-all.sh` (expanded proves) |
| `workflow_dispatch` (`run_full_mesh` default **true**) | `./run-all.sh` on `ubuntu-latest` DinD or self-hosted |

`ci_minimal` input selects full (`0`, default) vs lean (`1`) compose. Job summaries must not claim Headscale HA or full-stack×4.

## Ownership split (DNS)

| Production | CI equivalent | Prove script |
|---|---|---|
| `cloudflare-multi-ddns` global multi-A | CoreDNS `*` → all node Traefik IPs | `prove-production-dns.sh`, `prove-module5-ddns.sh` |
| `favonia cloudflare-ddns` node-direct | CoreDNS `*.ci-nodeN` → single node IP | `prove-production-dns.sh` |
| Traefik peer-forward via public DNS | Same, plus DinD `extra_hosts` shim for bridge nets | `prove-matrix.sh`, `prove-production-dns.sh` |
| Docker embedded DNS for backends | `getent`/`nslookup 127.0.0.11` from Traefik | `prove-production-dns.sh` |
| MagicDNS (Headscale/Tailscale) | Same when mesh up | `prove-dns.sh` |

Do **not** run `scripts/osvc_ingress_sync.py` on CI nodes (`OSVC_INGRESS_SYNC_DISABLE=1`).

## Replica ensure

`FAILOVER_REPLICA_ENSURE=true` on **ci-node1 only** after peer Docker API is reachable
(Tailscale-bound on Multipass/QEMU; private-net `:2375` on DinD).

Allowlisted Tier-A services (`FAILOVER_COMPOSE_ENSURE_SERVICES`, default `bolabaden-nextjs,autokuma`) use minimal registry config + `FAILOVER_REPLICA_PULL=never` (no Hub pull on peers). Ensure failures for that allowlist fail `/healthz` when `FAILOVER_REPLICA_ENSURE_STRICT=true` (CI default on main). Shape marks `intentionally_stopped` so ensure never undoes whoami/ci-probe/HS placement.
