#!/usr/bin/env bash
# Join nodes to Headscale (works on Multipass/QEMU/DinD).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

log "Headscale mesh: primary control plane expected on ci-node1"

if [[ "$(backend)" == "dind" ]]; then
  cleanup_stale_dind_bridge_filters
fi

# Ensure Tailscale present (DinD bootstrap installs static binaries)
for name in "${NODES[@]}"; do
  log "ensuring tailscale on $name"
  if [[ "$(backend)" == "dind" ]]; then
    if ! vm_exec "$name" "command -v tailscale >/dev/null 2>&1"; then
      arch="$(vm_exec "$name" "uname -m")"
      case "$arch" in aarch64|arm64) ts_arch=arm64 ;; *) ts_arch=amd64 ;; esac
      ver="${TAILSCALE_VERSION:-1.80.3}"
      vm_exec "$name" "bash -s" <<EOS
set -euo pipefail
cd /tmp
wget -q -O ts.tgz 'https://pkgs.tailscale.com/stable/tailscale_${ver}_${ts_arch}.tgz' \
  || curl -fsSL -o ts.tgz 'https://pkgs.tailscale.com/stable/tailscale_${ver}_${ts_arch}.tgz'
tar -xzf ts.tgz
cp -f tailscale_${ver}_${ts_arch}/tailscale tailscale_${ver}_${ts_arch}/tailscaled /usr/local/bin/
chmod +x /usr/local/bin/tailscale /usr/local/bin/tailscaled
mkdir -p /dev/net /var/lib/tailscale /var/run/tailscale
[[ -c /dev/net/tun ]] || mknod /dev/net/tun c 10 200 || true
EOS
    fi
  else
    vm_exec "$name" 'bash -s' <<'EOS'
set -euo pipefail
if command -v tailscale >/dev/null 2>&1; then
  exit 0
fi
curl -fsSL https://tailscale.com/install.sh | sh
EOS
  fi
done

log "ensuring headscale user + preauth key on ci-node1"
# Wait for headscale-server readiness
for i in $(seq 1 30); do
  if vm_exec ci-node1 "docker exec headscale-server headscale users list >/dev/null 2>&1"; then
    break
  fi
  sleep 2
done

vm_exec ci-node1 'bash -s' <<'EOS'
set -euo pipefail
create_user() {
  docker exec headscale-server headscale users list -o json 2>/dev/null \
    | python3 -c 'import sys,json; u=json.load(sys.stdin); sys.exit(0 if any((x.get("name") or "")=="failover-ci" for x in (u if isinstance(u,list) else [])) else 1)' \
    || docker exec headscale-server headscale users create failover-ci
  sleep 1
}
create_user
HS_UID=""
for _ in 1 2 3 4 5; do
  HS_UID=$(docker exec headscale-server headscale users list -o json 2>/dev/null \
    | python3 -c 'import sys,json; u=json.load(sys.stdin);
print(next((str(x.get("id")) for x in (u if isinstance(u,list) else []) if (x.get("name") or "")=="failover-ci"), ""))' || true)
  [[ -n "${HS_UID}" ]] && break
  sleep 1
done
[[ -n "${HS_UID}" ]] || { echo "failed to resolve headscale user id" >&2; exit 1; }

extract_key() {
  local json="$1"
  local k
  k=$(printf '%s' "$json" | python3 -c 'import sys,json; d=json.load(sys.stdin); print(d.get("key") or "")' 2>/dev/null || true)
  [[ -n "$k" ]] && { printf '%s' "$k"; return 0; }
  k=$(printf '%s' "$json" | jq -r '.key // empty' 2>/dev/null || true)
  [[ -n "$k" ]] && { printf '%s' "$k"; return 0; }
  printf '%s' "$json" | sed -n 's/.*"key"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1
}

create_preauth() {
  local who="$1"
  docker exec headscale-server headscale preauthkeys create --user "$who" --reusable --expiration 24h -o json 2>/dev/null
}

KEY=""
KEY_JSON=""
# Headscale 0.25.x expects username; 0.29+ expects numeric id — try both.
for who in failover-ci "$HS_UID"; do
  if KEY_JSON=$(create_preauth "$who"); then
    KEY=$(extract_key "$KEY_JSON" || true)
    [[ -n "$KEY" ]] && break
  fi
  sleep 1
done
[[ -n "$KEY" ]] || { echo "failed to create preauth key (uid=${HS_UID})" >&2; exit 1; }
mkdir -p /opt/failover-ci
printf '%s' "$KEY" > /opt/failover-ci/headscale-preauth.key
chmod 600 /opt/failover-ci/headscale-preauth.key
echo "$KEY"
EOS

KEY="$(vm_exec ci-node1 'cat /opt/failover-ci/headscale-preauth.key')"
[[ -n "$KEY" ]] || die "empty preauth key"
# Guard against accidental table-parse garbage (spaces / multi-field)
[[ "$KEY" != *" "* ]] || die "preauth key looks malformed: $KEY"
echo "$KEY" > "${STATE_DIR}/headscale-preauth.key"
chmod 600 "${STATE_DIR}/headscale-preauth.key"

LOGIN_SERVER="${HEADSCALE_LOGIN_SERVER:-http://$(node_ip_from_state ci-node1):8081}"
# Quick reachability gate before hanging on tailscale up
if [[ "$(backend)" == "dind" ]]; then
  for name in ci-node2 ci-node3 ci-node4; do
    if ! vm_exec "$name" "wget -qO- --timeout=5 '${LOGIN_SERVER}/health' >/dev/null 2>&1"; then
      cleanup_stale_dind_bridge_filters
      vm_exec "$name" "wget -qO- --timeout=5 '${LOGIN_SERVER}/health' >/dev/null 2>&1" \
        || die "$name cannot reach Headscale at ${LOGIN_SERVER}/health (check host iptables raw DROP)"
    fi
  done
fi

start_tailscaled_dind() {
  local name="$1"
  docker exec "$name" sh -lc 'pkill -9 tailscale tailscaled 2>/dev/null || true'
  sleep 1
  docker exec "$name" sh -lc 'mkdir -p /var/lib/tailscale /var/run/tailscale /dev/net
    [[ -c /dev/net/tun ]] || mknod /dev/net/tun c 10 200 || true'
  if docker exec "$name" sh -lc 'test -c /dev/net/tun'; then
    docker exec -d "$name" sh -lc \
      'exec /usr/local/bin/tailscaled --state=/var/lib/tailscale/tailscaled.state --socket=/var/run/tailscale/tailscaled.sock >>/var/log/tailscaled.log 2>&1'
  else
    docker exec -d "$name" sh -lc \
      'exec /usr/local/bin/tailscaled --tun=userspace-networking --state=/var/lib/tailscale/tailscaled.state --socket=/var/run/tailscale/tailscaled.sock >>/var/log/tailscaled.log 2>&1'
  fi
  local i
  for i in $(seq 1 20); do
    if docker exec "$name" sh -lc 'test -S /var/run/tailscale/tailscaled.sock'; then
      return 0
    fi
    sleep 0.5
  done
  die "tailscaled socket not ready on $name"
}

for name in "${NODES[@]}"; do
  log "tailscale up on $name → $LOGIN_SERVER"
  if [[ "$(backend)" == "dind" ]]; then
    start_tailscaled_dind "$name"
    # Hard timeout — hung auth was the main DinD failure mode
    if ! docker exec "$name" sh -lc \
      "timeout 75 /usr/local/bin/tailscale up --login-server='${LOGIN_SERVER}' --authkey='${KEY}' --hostname='${name}' --accept-dns=true --reset"; then
      warn "tailscale up failed/timed out on $name — retry once"
      start_tailscaled_dind "$name"
      docker exec "$name" sh -lc \
        "timeout 75 /usr/local/bin/tailscale up --login-server='${LOGIN_SERVER}' --authkey='${KEY}' --hostname='${name}' --accept-dns=true --reset" \
        || die "tailscale up failed on $name"
    fi
  else
    vm_exec "$name" "sudo timeout 75 tailscale up --login-server='${LOGIN_SERVER}' --authkey='${KEY}' --hostname='${name}' --accept-dns=true --reset || true"
  fi
done

: > "${STATE_DIR}/tailscale-ips.txt"
mesh_ok=1
for name in "${NODES[@]}"; do
  tip="$(vm_exec "$name" "tailscale ip -4 2>/dev/null | head -1" || true)"
  log "  $name tailscale=${tip:-unknown}"
  echo "${name}=${tip}" >> "${STATE_DIR}/tailscale-ips.txt"
  if [[ -z "$tip" || "$tip" == unknown ]]; then
    mesh_ok=0
  fi
done

[[ "$mesh_ok" -eq 1 ]] || die "one or more nodes missing Tailscale IPs — see ${STATE_DIR}/tailscale-ips.txt"
log "provision-mesh complete"
