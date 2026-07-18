# Failover CI — 4-VM wrong-node Traefik + dual-DNS (+ Module 5 dry-run)

Hands-off local driver. **Fully functional without GitHub Actions or git push.**

## Topology

| VM | Headscale | Probes |
|---|---|---|
| ci-node1 | **server + UI** (+ CoreDNS) | none |
| ci-node2 | **UI only** (+ CoreDNS) | **whoami** |
| ci-node3 | no | **whoami** |
| ci-node4 | no | **ci-probe** |

DNS: CoreDNS for mesh (not Cloudflare). Resolver order: MagicDNS → CoreDNS → Google `8.8.8.8`/`8.8.4.4`.  
Module 5 multi-record Cloudflare DDNS is proven with `--dry-run` (`prove-module5-ddns.sh`); optional live zone with `CF_LIVE_MULTI_DDNS=1`.

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
./run-all.sh          # full 4-node mesh (DinD / Multipass / QEMU)
# nested guest without KVM auto-selects DinD; or force:
#   FAILOVER_CI_BACKEND=dind ./provision-vms.sh
```

Teardown: `./teardown.sh`

## Compose modes

**Default:** full root [`docker-compose.yml`](../../docker-compose.yml) (+ ci-probes / extra_hosts overlays) on all 4 nodes, brought up **sequentially** (DinD RAM/disk).

**Fast debug:** `FAILOVER_CI_MINIMAL=1` uses [`compose/docker-compose.ci-stack.yml`](compose/docker-compose.ci-stack.yml) (Traefik, whoami, headscale, failover-agent, ci-probe only).

```bash
FAILOVER_CI_MINIMAL=1 ./compose-up-all.sh   # lean stack
FAILOVER_CI_MINIMAL=0 ./compose-up-all.sh   # full media stack (default)
```

## DinD specifics

- **`failover-agent` runs as `user: "0:0"`** — the DinD docker.sock is `root:docker` mode `660`; the image default `nobody` cannot write `failover-fallbacks.yaml`.
- **Tailscale + Headscale are required** on DinD (no soft-skip). Static Tailscale binaries are installed in bootstrap; `provision-mesh.sh` joins all 4 nodes and hard-fails if any node lacks a Tailscale IP. MagicDNS `@100.100.100.100` is a prove-dns hard gate when Tailscale is up.
- **Traefik plugins:** DinD DNS often cannot reach `plugins.traefik.io`, which disables *all* Traefik plugins and leaves routers referencing `crowdsec@file` disabled (HTTPS 404). [`compose/docker-compose.ci-dind-fixes.yml`](compose/docker-compose.ci-dind-fixes.yml) stubs `crowdsec` / error middlewares as no-op `headers` and pins `DOCKER_API_VERSION=1.47`.
- **Stale host iptables:** Docker `raw` PREROUTING pins container IPs to a bridge. After DinD network recreate, orphan `! -i br-<dead> -j DROP` rules for `ci-node1`'s IP blackhole peer→n1 traffic (ARP still works). `cleanup_stale_dind_bridge_filters` in `lib.sh` removes those for our node IPs during provision/mesh.
- Peer Docker API stays on private-net `:2375` (Tailscale-bound bind is Multipass/QEMU only).
- Full stack creates external `warp-nat-net` inside each DinD before `compose up`.

## GitHub Actions

Workflow: [`.github/workflows/failover-mesh.yml`](../../.github/workflows/failover-mesh.yml)

| Trigger | Job |
|---|---|
| PR / push (path-filtered) | `validate-local.sh` on `ubuntu-latest` |
| `workflow_dispatch` + `run_full_mesh` + `mesh_runner=dind-ubuntu-latest` | `./run-all.sh` with `FAILOVER_CI_BACKEND=dind` on `ubuntu-latest` (privileged Docker) |
| `workflow_dispatch` + `mesh_runner=self-hosted` | same on `[self-hosted, failover-mesh]` (Multipass/QEMU hosts) |

`ci_minimal` input selects full (`0`, default) vs lean (`1`) compose.

## Ownership split (DNS)

- **CoreDNS (CI)** / Traefik peer-forward: failover path
- **favonia `cloudflare-ddns`**: node-direct `*.$TS_HOSTNAME.$DOMAIN` only
- **`cloudflare-multi-ddns`** (profile `multi-ddns`): global multi-A for configured FQDNs

Do **not** run `scripts/osvc_ingress_sync.py` on CI nodes (`OSVC_INGRESS_SYNC_DISABLE=1`).

## Replica ensure

`FAILOVER_REPLICA_ENSURE=true` on **ci-node1 only** after peer Docker API is reachable
(Tailscale-bound on Multipass/QEMU; private-net `:2375` on DinD).
