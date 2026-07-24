#!/usr/bin/env bash
# Shared helpers for failover-ci (4-node mesh).
# Backends: multipass | qemu (needs /dev/kvm) | dind (nested fallback)
set -euo pipefail

FAILOVER_CI_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${FAILOVER_CI_ROOT}/../.." && pwd)"
STATE_DIR="${FAILOVER_CI_STATE_DIR:-${FAILOVER_CI_ROOT}/state}"
ENV_FILE="${FAILOVER_CI_ENV:-${FAILOVER_CI_ROOT}/env/test.env}"

NODES=(ci-node1 ci-node2 ci-node3 ci-node4)
HS_NODES=(ci-node1 ci-node2)
COREDNS_NODES=(ci-node1 ci-node2)
DIND_NETWORK="${DIND_NETWORK:-failover-ci-net}"
DIND_IMAGE="${DIND_IMAGE:-docker:27-dind}"

log()  { printf '[failover-ci] %s\n' "$*"; }
die()  { printf '[failover-ci] ERROR: %s\n' "$*" >&2; exit 1; }
need() { command -v "$1" >/dev/null 2>&1 || die "required command missing: $1"; }
warn() { printf '[failover-ci] WARN: %s\n' "$*" >&2; }

load_env() {
  if [[ -f "${ENV_FILE}" ]]; then
    # shellcheck disable=SC1090
    set -a; source "${ENV_FILE}"; set +a
  elif [[ -f "${FAILOVER_CI_ROOT}/env/test.env.example" ]]; then
    log "using env/test.env.example (copy to env/test.env to customize)"
    # shellcheck disable=SC1090
    set -a; source "${FAILOVER_CI_ROOT}/env/test.env.example"; set +a
  fi
  DOMAIN="${DOMAIN:-ci.bolabaden.test}"
  MAIN_HOST="${MAIN_HOST:-ci-node1}"
  VM_CPUS="${VM_CPUS:-2}"
  VM_MEM="${VM_MEM:-4096}"
  VM_DISK="${VM_DISK:-40G}"
  MULTIPASS_IMAGE="${MULTIPASS_IMAGE:-22.04}"
  HEADSCALE_MAGICDNS="${HEADSCALE_MAGICDNS:-100.100.100.100}"
  HEADSCALE_BASE_DOMAIN="${HEADSCALE_BASE_DOMAIN:-myscale.${DOMAIN}}"
  GOOGLE_DNS_1="${GOOGLE_DNS_1:-8.8.8.8}"
  GOOGLE_DNS_2="${GOOGLE_DNS_2:-8.8.4.4}"
  if [[ -z "${COMPOSE_UP_FLAGS+x}" || -z "${COMPOSE_UP_FLAGS}" ]]; then
    COMPOSE_UP_FLAGS="-d --remove-orphans --pull=always"
  fi
  mkdir -p "${STATE_DIR}"
}

# Detect: on nested hosts (no /dev/kvm) prefer DinD; else multipass → qemu → dind → none
detect_backend() {
  if [[ -n "${FAILOVER_CI_BACKEND:-}" ]]; then
    echo "${FAILOVER_CI_BACKEND}"
    return
  fi
  # Nested QEMU guests usually lack /dev/kvm; Multipass/libvirt are poor fits.
  if [[ ! -e /dev/kvm ]] \
    && command -v docker >/dev/null 2>&1 \
    && docker info >/dev/null 2>&1; then
    warn "no /dev/kvm (nested guest?) — using privileged Docker-in-Docker nodes"
    echo dind
    return
  fi
  if command -v multipass >/dev/null 2>&1 && multipass version >/dev/null 2>&1; then
    echo multipass
    return
  fi
  if [[ -e /dev/kvm ]] \
    && command -v virt-install >/dev/null 2>&1 \
    && command -v virsh >/dev/null 2>&1 \
    && command -v qemu-img >/dev/null 2>&1; then
    echo qemu
    return
  fi
  if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
    warn "QEMU/libvirt/multipass unavailable — falling back to privileged Docker-in-Docker nodes"
    echo dind
    return
  fi
  echo none
}

backend() {
  cat "${STATE_DIR}/backend" 2>/dev/null || detect_backend
}

save_backend() {
  echo "$1" > "${STATE_DIR}/backend"
}

vm_exec() {
  local name="$1"; shift
  local b
  b="$(backend)"
  case "$b" in
    multipass)
      multipass exec "$name" -- bash -lc "$*"
      ;;
    qemu)
      local ssh_cfg="${STATE_DIR}/ssh/${name}"
      [[ -f "$ssh_cfg" ]] || die "missing SSH config for $name ($ssh_cfg)"
      # shellcheck disable=SC1090
      source "$ssh_cfg"
      ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -i "${SSH_KEY}" "${SSH_USER}@${SSH_HOST}" "$*"
      ;;
    dind)
      # -i so callers can pipe heredocs (compose-up, resolvers, mesh)
      docker exec -i "$name" sh -lc "$*"
      ;;
    *)
      die "no VM backend configured (install multipass/qemu+kvm or docker for dind)"
      ;;
  esac
}

vm_ip() {
  local name="$1"
  local b
  b="$(backend)"
  case "$b" in
    multipass)
      multipass info "$name" --format csv 2>/dev/null | awk -F, 'NR==2{print $3}'
      ;;
    qemu)
      local ssh_cfg="${STATE_DIR}/ssh/${name}"
      # shellcheck disable=SC1090
      source "$ssh_cfg"
      echo "${SSH_HOST}"
      ;;
    dind)
      docker inspect -f "{{(index .NetworkSettings.Networks \"${DIND_NETWORK}\").IPAddress}}" "$name"
      ;;
    *)
      die "no VM backend"
      ;;
  esac
}

vm_transfer() {
  local name="$1" src="$2" dest="$3"
  local b
  b="$(backend)"
  case "$b" in
    multipass)
      multipass transfer "$src" "${name}:${dest}"
      ;;
    qemu)
      local ssh_cfg="${STATE_DIR}/ssh/${name}"
      # shellcheck disable=SC1090
      source "$ssh_cfg"
      scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -i "${SSH_KEY}" "$src" "${SSH_USER}@${SSH_HOST}:${dest}"
      ;;
    dind)
      # docker cp into privileged DinD is unreliable; stream via stdin/tar
      local dest_dir
      dest_dir="$(dirname "$dest")"
      docker exec "$name" sh -lc "mkdir -p \"$dest_dir\""
      if [[ -d "$src" ]]; then
        tar -C "$src" -cf - . | docker exec -i "$name" sh -lc "mkdir -p \"$dest\" && tar -C \"$dest\" -xf -"
      else
        docker exec -i "$name" sh -lc "cat > \"$dest\"" < "$src"
      fi
      ;;
    *)
      die "no VM backend"
      ;;
  esac
}

# Ensure image exists on host Docker (GHA DinD seed path; avoids Hub pull inside nested docker).
ensure_host_image() {
  local ref="$1" cand
  for cand in "$ref" "docker.io/${ref}" "${ref#docker.io/}"; do
    if docker image inspect "$cand" >/dev/null 2>&1; then
      return 0
    fi
  done
  log "pulling $ref on host for DinD seed"
  docker pull "$ref"
}

# Load an image from the host Docker into a DinD node's inner docker (Hub 429 workaround).
dind_seed_image_from_host() {
  local node="$1" ref="$2"
  [[ "$(backend)" == "dind" ]] || return 0
  ensure_host_image "$ref" || return 1
  local cand ref_found=""
  for cand in "$ref" "docker.io/${ref}" "${ref#docker.io/}"; do
    if docker image inspect "$cand" >/dev/null 2>&1; then
      ref_found="$cand"
      break
    fi
  done
  [[ -n "$ref_found" ]] || { warn "host missing image $ref — DinD inner docker may pull from Hub"; return 1; }
  if docker exec "$node" docker image inspect "$ref" >/dev/null 2>&1; then
    log "DinD $node already has $ref"
    return 0
  fi
  log "seeding $ref from host → DinD $node (via $ref_found)"
  docker save "$ref_found" | docker exec -i "$node" docker load
  docker exec "$node" sh -c "docker image inspect '$ref' >/dev/null 2>&1 || docker tag '$ref_found' '$ref' 2>/dev/null || true"
}

peers_csv_for() {
  local self="$1"
  local out=()
  local n
  for n in "${NODES[@]}"; do
    [[ "$n" == "$self" ]] && continue
    out+=("$n")
  done
  local IFS=,
  echo "${out[*]}"
}

all_node_ips_json() {
  local n ip
  echo -n '{'
  local first=1
  for n in "${NODES[@]}"; do
    ip="$(vm_ip "$n")"
    [[ -n "$ip" ]] || die "no IP for $n"
    if [[ $first -eq 1 ]]; then first=0; else echo -n ','; fi
    printf '"%s":"%s"' "$n" "$ip"
  done
  echo '}'
}

write_node_ips() {
  all_node_ips_json > "${STATE_DIR}/node-ips.json"
  log "wrote ${STATE_DIR}/node-ips.json"
}

# Host bridge-nf + Docker raw PREROUTING pin each container IP to its bridge.
# After recreate, stale "! -i br-<dead>" DROP rules for the same IP blackhole
# peer traffic (e.g. n2→n1) while ARP still works. Only touch IPs we own.
cleanup_stale_dind_bridge_filters() {
  command -v iptables >/dev/null 2>&1 || return 0
  [[ -f "${STATE_DIR}/node-ips.json" ]] || return 0
  local br_id br line iface ip deleted=0
  br_id="$(docker network inspect -f '{{.Id}}' "${DIND_NETWORK}" 2>/dev/null | cut -c1-12 || true)"
  [[ -n "$br_id" ]] || return 0
  br="br-${br_id}"

  mapfile -t _ci_ips < <(python3 -c 'import json,sys; print("\n".join(json.load(open(sys.argv[1])).values()))' "${STATE_DIR}/node-ips.json")
  [[ "${#_ci_ips[@]}" -gt 0 ]] || return 0

  for ip in "${_ci_ips[@]}"; do
    while IFS= read -r line; do
      [[ "$line" == "-A PREROUTING"*"-d ${ip}/32"* ]] || continue
      [[ "$line" == *"! -i br-"* ]] || continue
      [[ "$line" == *"-j DROP"* ]] || continue
      iface="${line#*! -i }"
      iface="${iface%% *}"
      # Keep the live failover-ci-net pin; drop pins for missing/other bridges
      if [[ "$iface" == "$br" ]] && ip link show "$iface" >/dev/null 2>&1; then
        continue
      fi
      if [[ "$iface" == br-* ]] && ip link show "$iface" >/dev/null 2>&1 && [[ "$iface" != "$br" ]]; then
        # Competing live bridge for same IP — remove the non-CI pin
        :
      elif ip link show "$iface" >/dev/null 2>&1; then
        continue
      fi
      # shellcheck disable=SC2086
      if iptables -t raw -D ${line#-A }; then
        deleted=$((deleted + 1))
        warn "removed stale raw DROP ${ip} via ${iface}"
      fi
    done < <(iptables-save -t raw 2>/dev/null | grep "^-A PREROUTING" || true)
  done

  # Orphan MASQUERADE only for our subnet when the CI bridge iface is gone
  if ! ip link show "$br" >/dev/null 2>&1; then
    while IFS= read -r line; do
      [[ "$line" == "-A POSTROUTING"*"-s 10.0.3.0/24"* ]] || continue
      [[ "$line" == *"! -o ${br}"* ]] || continue
      # shellcheck disable=SC2086
      iptables -t nat -D ${line#-A } 2>/dev/null && deleted=$((deleted + 1)) || true
    done < <(iptables-save -t nat 2>/dev/null | grep '^-A POSTROUTING' || true)
  fi

  if [[ "$deleted" -gt 0 ]]; then
    log "cleaned ${deleted} stale DinD bridge iptables rule(s)"
  fi
}

node_ip_from_state() {
  local name="$1"
  python3 - "$name" "${STATE_DIR}/node-ips.json" <<'PY'
import json,sys
name,path=sys.argv[1],sys.argv[2]
print(json.load(open(path))[name])
PY
}

# CI-minimal only when explicitly requested (fast debug). Full root compose is default.
use_ci_minimal() {
  if [[ "${FAILOVER_CI_MINIMAL:-}" == "1" || "${FAILOVER_CI_MINIMAL:-}" == "true" ]]; then
    return 0
  fi
  return 1
}

# Echo compose -f args (relative to VM_REPO_PATH)
compose_f_args() {
  if use_ci_minimal; then
    echo "-f compose/docker-compose.ci-stack.yml -f compose/docker-compose.ci-probes.yml -f compose/docker-compose.ci-extra-hosts.yml"
  else
    local args="-f docker-compose.yml -f compose/docker-compose.ci-probes.yml -f compose/docker-compose.ci-extra-hosts.yml"
    if [[ "$(backend)" == "dind" ]]; then
      args+=" -f compose/docker-compose.ci-dind-fixes.yml"
    fi
    echo "$args"
  fi
}

# HA-critical curated peer set (NOT full media stack — honesty contract).
# Used by sync-images, compose-up-critical, and prove gates.
ha_critical_services() {
  echo "traefik whoami headscale-server headscale failover-agent ci-probe bolabaden-nextjs autokuma"
}

# Per-node HA-critical compose services — must match shape-placement.sh roles.
ha_critical_services_for_node() {
  local node="$1"
  case "$node" in
    ci-node1) echo "traefik headscale-server headscale failover-agent bolabaden-nextjs autokuma" ;;
    ci-node2) echo "traefik headscale failover-agent whoami bolabaden-nextjs autokuma" ;;
    ci-node3) echo "traefik failover-agent whoami bolabaden-nextjs autokuma" ;;
    ci-node4) echo "traefik failover-agent ci-probe bolabaden-nextjs autokuma" ;;
    *) ha_critical_services ;;
  esac
}

# Hard-fail if any of these are not running after critical compose up on a node.
ha_critical_must_run_on_node() {
  local node="$1"
  case "$node" in
    ci-node1) echo "traefik headscale-server failover-agent bolabaden-nextjs autokuma" ;;
    ci-node2) echo "traefik headscale whoami failover-agent bolabaden-nextjs autokuma" ;;
    ci-node3) echo "traefik whoami failover-agent bolabaden-nextjs autokuma" ;;
    ci-node4) echo "traefik ci-probe failover-agent bolabaden-nextjs autokuma" ;;
    *) echo "traefik failover-agent bolabaden-nextjs autokuma" ;;
  esac
}

# Images that must sync even when FAILOVER_CI_IMAGE_MAX_MB would skip them.
ha_critical_images() {
  cat <<'EOF'
traefik:latest
traefik/whoami:latest
bolabaden/failover-agent:latest
local/failover-ci-probe:latest
headscale/headscale:latest
ghcr.io/gurucomputing/headscale-ui:latest
coredns/coredns:1.11.3
docker.io/bolabaden/bolabaden-nextjs:latest
ghcr.io/bigboot/autokuma:latest
EOF
}

# Traefik (bridge) cannot resolve peer FQDNs via CoreDNS from nested nets;
# pin svc.<node>.$DOMAIN → node IP so peer HTTPS URLs work.
write_ci_extra_hosts_compose() {
  local out="${1:-${STATE_DIR}/docker-compose.ci-extra-hosts.yml}"
  local domain="${DOMAIN}"
  [[ -f "${STATE_DIR}/node-ips.json" ]] || die "node-ips.json required for extra_hosts"
  python3 - "$out" "$domain" "${STATE_DIR}/node-ips.json" <<'PY'
import json, sys
out, domain, path = sys.argv[1], sys.argv[2], sys.argv[3]
ips = json.load(open(path))
svcs = ("whoami", "ci-probe", "headscale", "headscale-server", "traefik", "failover-agent", "bolabaden-nextjs", "autokuma")
lines = [
    "# Generated — peer FQDN → DinD node IP for Traefik peer-forward",
    "services:",
    "  traefik:",
    "    extra_hosts:",
]
for node, ip in sorted(ips.items()):
    lines.append(f'      - "{node}:{ip}"')
    lines.append(f'      - "{node}.{domain}:{ip}"')
    for svc in svcs:
        lines.append(f'      - "{svc}.{node}.{domain}:{ip}"')
open(out, "w").write("\n".join(lines) + "\n")
print(out)
PY
}
