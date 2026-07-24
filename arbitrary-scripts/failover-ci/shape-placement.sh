#!/usr/bin/env bash
# Heterogeneous placement after compose up:
#   ci-node1: headscale-server + UI; stop whoami + ci-probe
#   ci-node2: UI only (stop headscale-server); keep whoami; stop ci-probe
#   ci-node3: whoami; stop HS* + ci-probe
#   ci-node4: ci-probe; stop whoami + HS*
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"
DOCKER_PORT="${FAILOVER_REMOTE_DOCKER_PORT:-2375}"
CF_ARGS="$(compose_f_args)"
DC_BASE="docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS}"

stop_svcs() {
  local name="$1"; shift
  # restart:always / deunhealth will bounce plain docker stop — disable restart first
  vm_exec "$name" "cd ${VM_REPO_PATH} && export COMPOSE_PROJECT_NAME=${name}
    ${DC_BASE} stop $* 2>/dev/null || true
    for s in $*; do
      docker update --restart=no \$s 2>/dev/null || true
      docker stop \$s 2>/dev/null || true
      docker rm -f \$s 2>/dev/null || true
    done"
}

start_svcs() {
  local name="$1"; shift
  vm_exec "$name" "cd ${VM_REPO_PATH} && export COMPOSE_PROJECT_NAME=${name}
    for s in $*; do docker update --restart=always \$s 2>/dev/null || true; done
    ${DC_BASE} start $* 2>/dev/null || true
    for s in $*; do docker start \$s 2>/dev/null || true; done"
}

log "shaping ci-node1 (HS server+UI; no whoami/ci-probe)"
stop_svcs ci-node1 whoami ci-probe
start_svcs ci-node1 headscale-server headscale

log "shaping ci-node2 (UI only + whoami; no headscale-server / ci-probe)"
stop_svcs ci-node2 headscale-server ci-probe
start_svcs ci-node2 headscale whoami

log "shaping ci-node3 (whoami; no HS / ci-probe)"
stop_svcs ci-node3 headscale-server headscale ci-probe
start_svcs ci-node3 whoami

log "shaping ci-node4 (ci-probe only)"
stop_svcs ci-node4 headscale-server headscale whoami
start_svcs ci-node4 ci-probe

# Expose Docker API for peer inventory / replica ensure
for name in "${NODES[@]}"; do
  b="$(backend)"
  if [[ "$b" == "dind" ]]; then
    log "DinD $name: dockerd already listening on tcp://0.0.0.0:${DOCKER_PORT} (private net)"
    continue
  fi
  log "Docker TCP :${DOCKER_PORT} on Tailscale IP for $name"
  vm_exec "$name" "bash -s" <<EOS
set -euo pipefail
TIP=\$(tailscale ip -4 2>/dev/null | head -1 || true)
if [[ -z "\$TIP" ]]; then
  echo "WARN: no Tailscale IP yet on $name — binding Docker API skipped until mesh is up" >&2
  exit 0
fi
sudo mkdir -p /etc/systemd/system/docker.service.d
sudo tee /etc/systemd/system/docker.service.d/failover-ci.conf >/dev/null <<EOF
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// -H unix:///var/run/docker.sock -H tcp://\${TIP}:${DOCKER_PORT}
EOF
if [[ -f /etc/docker/daemon.json ]]; then
  python3 - <<'PY'
import json
path="/etc/docker/daemon.json"
try:
    data=json.load(open(path))
except Exception:
    raise SystemExit(0)
if "hosts" in data:
    data.pop("hosts", None)
    open(path,"w").write(json.dumps(data,indent=2)+"\\n")
    print("removed hosts from daemon.json")
PY
fi
sudo systemctl daemon-reload
sudo systemctl restart docker
EOS
done

# Persist intentionally_stopped so replica ensure never undoes shape (R4).
# Agent events/inventory usually classify compose stop, but ensure can race —
# write the registry on ci-node1 before ensure loops run.
log "marking shaped-off placements intentionally_stopped in registry"
vm_exec ci-node1 "python3 - ${VM_REPO_PATH}/volumes/placement/services.yaml" <<'PY'
import sys, time
from pathlib import Path

path = Path(sys.argv[1])
path.parent.mkdir(parents=True, exist_ok=True)
now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
# service -> nodes that must stay intentionally_stopped after shape
shaped = {
    "whoami": ["ci-node1", "ci-node4"],
    "ci-probe": ["ci-node1", "ci-node2", "ci-node3"],
    "headscale-server": ["ci-node2", "ci-node3", "ci-node4"],
    "headscale": ["ci-node3", "ci-node4"],
}

# Write a valid registry fragment (never append to "services: {}" — that breaks YAML).
lines = ["services:"]
for svc, nodes in shaped.items():
    lines.extend([
        f"  {svc}:",
        f"    name: {svc}",
        f"    status: running",
        f"    updated_at: {now}",
        "    nodes:",
    ])
    for node in nodes:
        lines.extend([
            f"      {node}:",
            "        status: intentionally_stopped",
            f"        last_seen: {now}",
        ])

text = "\n".join(lines) + "\n"
path.write_text(text)
marker = path.parent / "shape-intentional-stops.txt"
marker.write_text("\n".join(f"{s}@{n}" for s, ns in shaped.items() for n in ns) + "\n")
print(f"wrote valid registry {path} ({len(shaped)} services) + {marker}")
PY

log "restarting failover-agent on ci-node1 (reload registry + reconcile)"
vm_exec ci-node1 "docker restart failover-agent >/dev/null"
sleep 8

# Wait for failover-agent YAML with Tier-A + canary routes
log "waiting for failover-agent YAML on ci-node1"
for i in $(seq 1 60); do
  if vm_exec ci-node1 "grep -q 'whoami-with-failover' ${VM_REPO_PATH}/volumes/traefik/dynamic/failover-fallbacks.yaml 2>/dev/null \
    && grep -q 'bolabaden-nextjs-with-failover' ${VM_REPO_PATH}/volumes/traefik/dynamic/failover-fallbacks.yaml 2>/dev/null"; then
    log "failover-fallbacks.yaml present (whoami + bolabaden Tier-A)"
    break
  fi
  sleep 3
done

# R4 observation: whoami must remain absent on n1/n4 after shape (+ brief ensure window)
sleep 5
for node in ci-node1 ci-node4; do
  if vm_exec "$node" "docker inspect -f '{{.State.Running}}' whoami 2>/dev/null | grep -qx true"; then
    die "R4: whoami running on ${node} after shape — ensure undid intentional stop"
  fi
done
if vm_exec ci-node1 "grep -q intentionally_stopped ${VM_REPO_PATH}/volumes/placement/services.yaml 2>/dev/null \
  || test -f ${VM_REPO_PATH}/volumes/placement/shape-intentional-stops.txt"; then
  log "registry/shape marker shows intentionally_stopped"
else
  die "R4: placement registry missing intentionally_stopped after shape"
fi

log "shape-placement complete"
log "  n1: HS server+UI + Tier-A; n2: UI+whoami + Tier-A; n3: whoami + Tier-A; n4: ci-probe + Tier-A"
