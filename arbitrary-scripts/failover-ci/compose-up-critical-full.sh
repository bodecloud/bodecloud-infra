#!/usr/bin/env bash
# Bring up prove-critical services via full root compose (pull=never).
# Used when Docker Hub rate-limits block full pulls on peer DinD nodes.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"
export FAILOVER_CI_MINIMAL=0
CF_ARGS="$(compose_f_args)"
TARGETS=("$@")
if [[ "${#TARGETS[@]}" -eq 0 ]]; then
  TARGETS=("${NODES[@]}")
fi

EXTRA_HOSTS_FILE="${STATE_DIR}/docker-compose.ci-extra-hosts.yml"
write_ci_extra_hosts_compose "${EXTRA_HOSTS_FILE}"
[[ "$(backend)" == "dind" ]] && cleanup_stale_dind_bridge_filters || true

for name in "${TARGETS[@]}"; do
  log "critical full-compose up on $name (pull=never)"
  replica_ensure=false
  [[ "$name" == "${MAIN_HOST:-ci-node1}" ]] && replica_ensure=true
  vm_transfer "$name" "${EXTRA_HOSTS_FILE}" "${VM_REPO_PATH}/compose/docker-compose.ci-extra-hosts.yml"
  vm_transfer "$name" "${FAILOVER_CI_ROOT}/compose/docker-compose.ci-dind-fixes.yml" \
    "${VM_REPO_PATH}/compose/docker-compose.ci-dind-fixes.yml"
  vm_exec "$name" "bash -s" <<EOS
set -euo pipefail
cd ${VM_REPO_PATH}
mkdir -p compose volumes/traefik/dynamic volumes/placement volumes/headscale/lib volumes/headscale/run ~/.docker /tmp/dev-secrets
echo '{}' > ~/.docker/config.json
cp -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-probes.yml compose/ 2>/dev/null || true
cp -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-dind-fixes.yml compose/ 2>/dev/null || true
grep -q '^FAILOVER_REPLICA_ENSURE=' .env 2>/dev/null \
  && sed -i 's/^FAILOVER_REPLICA_ENSURE=.*/FAILOVER_REPLICA_ENSURE=${replica_ensure}/' .env \
  || echo "FAILOVER_REPLICA_ENSURE=${replica_ensure}" >> .env
grep -q '^COMPOSE_PROJECT_NAME=' .env 2>/dev/null \
  && sed -i "s/^COMPOSE_PROJECT_NAME=.*/COMPOSE_PROJECT_NAME=${name}/" .env \
  || echo "COMPOSE_PROJECT_NAME=${name}" >> .env
sysctl -w net.ipv4.ip_forward=1 >/dev/null 2>&1 || true
sysctl -w fs.inotify.max_user_instances=1024 >/dev/null 2>&1 || true
sysctl -w fs.inotify.max_user_watches=524288 >/dev/null 2>&1 || true
iptables -t nat -C POSTROUTING -s 172.16.0.0/12 -o eth0 -j MASQUERADE 2>/dev/null \
  || iptables -t nat -A POSTROUTING -s 172.16.0.0/12 -o eth0 -j MASQUERADE 2>/dev/null || true
export COMPOSE_PROJECT_NAME=${name}
docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env \
  -f compose/docker-compose.ci-stack.yml -f compose/docker-compose.ci-probes.yml \
  down --remove-orphans 2>/dev/null || true
docker network inspect warp-nat-net >/dev/null 2>&1 || docker network create --driver bridge --attachable \
  --opt com.docker.network.bridge.name=br_warp-nat-net \
  --opt com.docker.network.bridge.enable_ip_masquerade=false \
  --subnet 10.0.2.0/24 --gateway 10.0.2.1 warp-nat-net
# stub secrets
python3 - "${VM_REPO_PATH}" <<'PY'
import re, sys
from pathlib import Path
root = Path(sys.argv[1])
sec = Path("/tmp/dev-secrets"); sec.mkdir(parents=True, exist_ok=True)
files = [root / "docker-compose.yml"]
text = files[0].read_text()
for m in re.finditer(r"^- compose/(\S+\.yml)\s*$", text, re.M):
    files.append(root / "compose" / m.group(1))
pat = re.compile(r"\$\{SECRETS_PATH[^}]*\}/([A-Za-z0-9_.-]+\.txt)")
for p in files:
    if p.is_file():
        for name in pat.findall(p.read_text()):
            f = sec / name
            if not f.exists():
                f.write_text("ci-placeholder\n")
PY
CF='${CF_ARGS}'
docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env \$CF \
  config --format json > /tmp/cc.json 2>/dev/null || true
python3 - "${VM_REPO_PATH}" /tmp/cc.json <<'PY'
import json, sys
from pathlib import Path
root = Path(sys.argv[1]).resolve()
cfg_path = Path(sys.argv[2])
if not cfg_path.is_file() or cfg_path.stat().st_size == 0:
    raise SystemExit(0)
allow = (str(root), "/tmp", "/opt", "/etc/", "/var/log", "/usr/share/zoneinfo")
skip = {"/var/run/docker.sock", "/run/docker.sock", "/dev/net/tun", "/sys", "/proc"}
cfg = json.loads(cfg_path.read_text())
def ens(src, force=False):
    if not src or src in skip or not src.startswith(allow):
        return
    p = Path(src)
    if p.exists():
        return
    looks = force or bool(p.suffix) or src.startswith("/etc/") or "/var/log/" in src
    try:
        p.parent.mkdir(parents=True, exist_ok=True)
        if looks:
            p.touch()
        else:
            p.mkdir(parents=True, exist_ok=True)
    except OSError:
        pass
for svc in (cfg.get("services") or {}).values():
    for m in svc.get("volumes") or []:
        if isinstance(m, dict) and m.get("type") == "bind":
            ens(m.get("source") or "")
for section in ("configs", "secrets"):
    for item in (cfg.get(section) or {}).values():
        if isinstance(item, dict) and item.get("file"):
            ens(item["file"], True)
Path("/etc/fuse.conf").touch()
PY
# Best-effort full stack without registry pulls; then force critical
set +e
docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env \$CF \
  up -d --remove-orphans --pull=never
set -e
# Clear name conflicts from prior CI-minimal / foreign project labels
docker rm -f traefik whoami headscale-server headscale failover-agent ci-probe \
  coolify-proxy 2>/dev/null || true
docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env \$CF \
  up -d --no-deps --remove-orphans --pull=never \
  traefik whoami headscale-server headscale failover-agent ci-probe
for crit in traefik whoami headscale-server failover-agent; do
  docker inspect -f '{{.State.Running}}' "\$crit" 2>/dev/null | grep -qx true \
    || { echo "ERROR: \$crit not running on ${name}" >&2; docker ps -a --format '{{.Names}} {{.Status}}' | head -30; exit 1; }
done
echo "[failover-ci] critical OK on ${name} containers=\$(docker ps -q | wc -l)"
EOS
done
log "compose-up-critical-full complete"
