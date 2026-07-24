#!/usr/bin/env bash
# Inside each node: compose up CI-minimal (FAILOVER_CI_MINIMAL=1) or full root stack.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"
CF_ARGS="$(compose_f_args)"
USE_MINIMAL=0
if use_ci_minimal; then
  USE_MINIMAL=1
  log "using CI-minimal compose stack (backend=$(backend))"
else
  log "using full docker-compose.yml stack (sequential)"
fi

CD1=""; CD2=""
if [[ -f "${STATE_DIR}/coredns-ips.txt" ]]; then
  mapfile -t CDS < "${STATE_DIR}/coredns-ips.txt"
  CD1="${CDS[0]:-}"
  CD2="${CDS[1]:-}"
fi

EXTRA_HOSTS_FILE="${STATE_DIR}/docker-compose.ci-extra-hosts.yml"
write_ci_extra_hosts_compose "${EXTRA_HOSTS_FILE}"
log "wrote peer extra_hosts overlay → ${EXTRA_HOSTS_FILE}"

if [[ "$(backend)" == "dind" ]]; then
  cleanup_stale_dind_bridge_filters
  bash "${SCRIPT_DIR}/seed-dind-images-from-host.sh"
fi

for name in "${NODES[@]}"; do
  log "compose up on $name"
  replica_ensure=false
  if [[ "$name" == "${MAIN_HOST:-ci-node1}" ]]; then
    replica_ensure=true
  fi
  CRIT_SVCS="$(ha_critical_services_for_node "$name")"
  MUST_RUN="$(ha_critical_must_run_on_node "$name")"
  vm_transfer "$name" "${EXTRA_HOSTS_FILE}" "${VM_REPO_PATH}/compose/docker-compose.ci-extra-hosts.yml"
  vm_exec "$name" "bash -s" <<EOS
set -euo pipefail
cd ${VM_REPO_PATH}
mkdir -p compose volumes/traefik/dynamic volumes/placement volumes/headscale/lib volumes/headscale/run
mkdir -p ~/.docker && echo '{}' > ~/.docker/config.json
if [[ -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-probes.yml ]]; then
  cp -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-probes.yml compose/docker-compose.ci-probes.yml
fi
if [[ -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-stack.yml ]]; then
  cp -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-stack.yml compose/docker-compose.ci-stack.yml
fi
if [[ -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-tier-a.yml ]]; then
  cp -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-tier-a.yml compose/docker-compose.ci-tier-a.yml
fi
if [[ -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-dind-fixes.yml ]]; then
  cp -f arbitrary-scripts/failover-ci/compose/docker-compose.ci-dind-fixes.yml compose/docker-compose.ci-dind-fixes.yml
fi
[[ -f compose/docker-compose.ci-extra-hosts.yml ]] || printf '%s\n' 'services: {}' > compose/docker-compose.ci-extra-hosts.yml
if grep -q '^FAILOVER_REPLICA_ENSURE=' .env 2>/dev/null; then
  sed -i 's/^FAILOVER_REPLICA_ENSURE=.*/FAILOVER_REPLICA_ENSURE=${replica_ensure}/' .env
else
  echo "FAILOVER_REPLICA_ENSURE=${replica_ensure}" >> .env
fi
grep -q '^FAILOVER_REPLICA_PULL=' .env 2>/dev/null \
  && sed -i 's/^FAILOVER_REPLICA_PULL=.*/FAILOVER_REPLICA_PULL=never/' .env \
  || echo "FAILOVER_REPLICA_PULL=never" >> .env
grep -q '^FAILOVER_REPLICA_ENSURE_STRICT=' .env 2>/dev/null \
  && sed -i "s/^FAILOVER_REPLICA_ENSURE_STRICT=.*/FAILOVER_REPLICA_ENSURE_STRICT=${replica_ensure}/" .env \
  || echo "FAILOVER_REPLICA_ENSURE_STRICT=${replica_ensure}" >> .env
grep -q '^FAILOVER_COMPOSE_ENSURE_SERVICES=' .env 2>/dev/null \
  && sed -i 's/^FAILOVER_COMPOSE_ENSURE_SERVICES=.*/FAILOVER_COMPOSE_ENSURE_SERVICES=bolabaden-nextjs,autokuma/' .env \
  || echo "FAILOVER_COMPOSE_ENSURE_SERVICES=bolabaden-nextjs,autokuma" >> .env
grep -q '^COREDNS_1=' .env 2>/dev/null && sed -i 's/^COREDNS_1=.*/COREDNS_1=${CD1}/' .env || echo "COREDNS_1=${CD1}" >> .env
grep -q '^COREDNS_2=' .env 2>/dev/null && sed -i 's/^COREDNS_2=.*/COREDNS_2=${CD2}/' .env || echo "COREDNS_2=${CD2}" >> .env
grep -q '^COMPOSE_PROJECT_NAME=' .env 2>/dev/null && sed -i "s/^COMPOSE_PROJECT_NAME=.*/COMPOSE_PROJECT_NAME=${name}/" .env || echo "COMPOSE_PROJECT_NAME=${name}" >> .env
grep -q '^STACK_NAME=' .env 2>/dev/null || echo "STACK_NAME=${name}" >> .env
grep -q '^ROOT_PATH=' .env 2>/dev/null || echo "ROOT_PATH=${VM_REPO_PATH}" >> .env

sysctl -w net.ipv4.ip_forward=1 >/dev/null 2>&1 || true
sysctl -w fs.inotify.max_user_instances=1024 >/dev/null 2>&1 || true
sysctl -w fs.inotify.max_user_watches=524288 >/dev/null 2>&1 || true
iptables -t nat -C POSTROUTING -s 172.16.0.0/12 -o eth0 -j MASQUERADE 2>/dev/null \
  || iptables -t nat -A POSTROUTING -s 172.16.0.0/12 -o eth0 -j MASQUERADE 2>/dev/null || true

export COMPOSE_PROJECT_NAME=${name}
export COMPOSE_PARALLEL_LIMIT=2
USE_MINIMAL=${USE_MINIMAL}

if [[ "\$USE_MINIMAL" == "1" ]]; then
  docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} \
    build failover-agent ci-probe
  for crit in ${CRIT_SVCS}; do
    docker rm -f "\$crit" 2>/dev/null || true
  done
  docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} \
    up -d --remove-orphans --force-recreate --pull=missing ${CRIT_SVCS}
  for crit in ${MUST_RUN}; do
    docker inspect -f '{{.State.Running}}' "\$crit" 2>/dev/null | grep -qx true \
      || { echo "[failover-ci] ERROR: HA-critical \$crit not running on ${name}" >&2; exit 1; }
  done
else
  # Tear down prior CI-minimal project containers so fixed bridge names can bind
  docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env \
    -f compose/docker-compose.ci-stack.yml -f compose/docker-compose.ci-probes.yml \
    down --remove-orphans 2>/dev/null || true

  # External network required by root compose
  if ! docker network inspect warp-nat-net >/dev/null 2>&1; then
    docker network create --driver bridge --attachable \
      --opt com.docker.network.bridge.name=br_warp-nat-net \
      --opt com.docker.network.bridge.enable_ip_masquerade=false \
      --subnet 10.0.2.0/24 \
      --gateway 10.0.2.1 \
      warp-nat-net
  fi

  # Rebuild agent when forced or when local test image tag missing (branch code drift)
  FORCE_BUILD=${FAILOVER_CI_FORCE_AGENT_BUILD:-0}
  if [[ "\$FORCE_BUILD" == "1" ]] || ! docker image inspect local/failover-agent:ci-test >/dev/null 2>&1; then
    if docker image inspect bolabaden/failover-agent:latest >/dev/null 2>&1 \
      && docker image inspect local/failover-ci-probe:latest >/dev/null 2>&1; then
      echo "[failover-ci] using host-seeded failover-agent + ci-probe (skip inner build)"
      docker tag bolabaden/failover-agent:latest local/failover-agent:ci-test 2>/dev/null || true
    else
    echo "[failover-ci] building failover-agent (FORCE=\$FORCE_BUILD)"
    for _try in 1 2 3 4 5; do
      if docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} \
        build failover-agent ci-probe; then
        docker tag bolabaden/failover-agent:latest local/failover-agent:ci-test 2>/dev/null || true
        break
      fi
      echo "[failover-ci] build retry \$_try after registry throttle..." >&2
      sleep \$((_try * 20))
    done
    fi
  elif docker image inspect bolabaden/failover-agent:latest >/dev/null 2>&1 \
    && docker image inspect local/failover-ci-probe:latest >/dev/null 2>&1; then
    echo "[failover-ci] using preloaded failover-agent + ci-probe images (set FAILOVER_CI_FORCE_AGENT_BUILD=1 to rebuild)"
  else
    for _try in 1 2 3 4 5; do
      if docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} \
        build failover-agent ci-probe; then
        break
      fi
      echo "[failover-ci] build retry \$_try after registry throttle..." >&2
      sleep \$((_try * 20))
    done
  fi
  # Ensure every SECRETS_PATH file referenced by compose exists (idempotent stubs)
  python3 - "${VM_REPO_PATH}" <<'PY' || true
import re, sys
from pathlib import Path
root = Path(sys.argv[1])
sec = Path("/tmp/dev-secrets")
sec.mkdir(parents=True, exist_ok=True)
files = [root / "docker-compose.yml"]
text = files[0].read_text()
for m in re.finditer(r"^- compose/(\S+\.yml)\s*$", text, re.M):
    files.append(root / "compose" / m.group(1))
pat = re.compile(r"\$\{SECRETS_PATH[^}]*\}/([A-Za-z0-9_.-]+\.txt)")
for p in files:
    if not p.is_file():
        continue
    for name in pat.findall(p.read_text()):
        f = sec / name
        if not f.exists():
            f.write_text("ci-placeholder\n")
PY

  # Materialize bind-mount sources under repo/tmp (skip host system paths)
  docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} \
    config --format json > /tmp/failover-ci-compose-config.json 2>/dev/null || true
  python3 - "${VM_REPO_PATH}" /tmp/failover-ci-compose-config.json <<'PY' || true
import json, sys
from pathlib import Path
root = Path(sys.argv[1]).resolve()
cfg_path = Path(sys.argv[2])
if not cfg_path.is_file() or cfg_path.stat().st_size == 0:
    raise SystemExit(0)
allow_prefixes = (str(root), "/tmp", "/opt", "/etc/", "/var/log", "/usr/share/zoneinfo")
skip_exact = {"/var/run/docker.sock", "/run/docker.sock", "/dev/net/tun", "/sys", "/proc"}
cfg = json.loads(cfg_path.read_text())

def ensure_path(src: str, force_file: bool = False) -> None:
    if not src or src in skip_exact or not src.startswith(allow_prefixes):
        return
    p = Path(src)
    if p.exists():
        return
    looks_file = force_file or bool(p.suffix) or src.endswith(
        (".log", ".conf", ".yml", ".yaml", ".json", ".txt", ".key", ".pem", ".crt", ".cfg")
    ) or src.startswith("/etc/") or "/var/log/" in src
    try:
        p.parent.mkdir(parents=True, exist_ok=True)
        if looks_file:
            p.touch()
        else:
            p.mkdir(parents=True, exist_ok=True)
    except OSError:
        pass

for svc in (cfg.get("services") or {}).values():
    for m in svc.get("volumes") or []:
        if isinstance(m, dict) and m.get("type") == "bind":
            ensure_path(m.get("source") or "")
# Compose configs/secrets with file: sources (crowdsec auth.log/syslog, etc.)
for section in ("configs", "secrets"):
    for item in (cfg.get(section) or {}).values():
        if isinstance(item, dict) and item.get("file"):
            ensure_path(item["file"], force_file=True)
PY

  # Full stack — best-effort; non-critical health failures are OK for prove gates.
  # Peers prefer --pull=never after image sync from MAIN (avoids Hub 429).
  PULL_POLICY=missing
  if [[ "${name}" != "${MAIN_HOST:-ci-node1}" && "$(backend)" == "dind" ]]; then
    PULL_POLICY=never
  fi
  set +e
  docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} \
    up -d --remove-orphans --pull=\$PULL_POLICY
  up_rc=$?
  set -e
  if [[ "\$up_rc" -ne 0 ]]; then
    echo "[failover-ci] WARN: full compose up exited \$up_rc on ${name} — forcing critical services" >&2
    docker network inspect warp-nat-net >/dev/null 2>&1 || \
      docker network create --driver bridge --attachable \
        --opt com.docker.network.bridge.name=br_warp-nat-net \
        --opt com.docker.network.bridge.enable_ip_masquerade=false \
        --subnet 10.0.2.0/24 --gateway 10.0.2.1 warp-nat-net
  fi
  # HA-critical curated set aligned with shape-placement (not full media stack)
  # Full-stack up may leave conflicting containers on peers — always recreate cleanly.
  for crit in ${CRIT_SVCS}; do
    docker rm -f "\$crit" 2>/dev/null || true
  done
  docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} \
    up -d --no-deps --remove-orphans --force-recreate --pull=\$PULL_POLICY \
    ${CRIT_SVCS}
  for crit in ${MUST_RUN}; do
    docker inspect -f '{{.State.Running}}' "\$crit" 2>/dev/null | grep -qx true \
      || { echo "[failover-ci] ERROR: HA-critical \$crit not running on ${name}" >&2; exit 1; }
  done
  echo "[failover-ci] peer/main HA-critical path done on ${name} (curated set, not full stack ×4)"
fi
docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env ${CF_ARGS} ps || true
EOS

  # After MAIN is up on DinD, seed peer image caches (Hub 429 workaround)
  if [[ "$USE_MINIMAL" != "1" && "$(backend)" == "dind" && "$name" == "${MAIN_HOST:-ci-node1}" ]]; then
    if [[ "${FAILOVER_CI_SYNC_IMAGES:-1}" == "1" ]]; then
      log "syncing images from ${name} → peers before peer compose-up"
      bash "${SCRIPT_DIR}/sync-images-from-main.sh" || warn "image sync failed — peers may hit registry 429"
    fi
  fi

  # Host-side headroom between sequential DinD nodes
  if [[ "$USE_MINIMAL" != "1" ]]; then
    log "post-${name}: host memory/disk"
    free -h | head -2 || true
    df -h / | tail -1 || true
    sleep 5
  fi
done

log "compose-up-all complete — run shape-placement.sh next"
