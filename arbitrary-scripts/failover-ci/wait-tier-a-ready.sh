#!/usr/bin/env bash
# Block until HA-critical Tier-A services are plural on the mesh and agent YAML has peer URLs.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"
YAML_PATH="${VM_REPO_PATH}/volumes/traefik/dynamic/failover-fallbacks.yaml"
MIN_NODES="${FAILOVER_CI_TIER_A_MIN_NODES:-3}"
TIMEOUT="${FAILOVER_CI_TIER_A_WAIT_SEC:-300}"

log "waiting up to ${TIMEOUT}s for Tier-A plurality (≥${MIN_NODES} nodes) + peer URLs in YAML"

deadline=$((SECONDS + TIMEOUT))
while (( SECONDS < deadline )); do
  ok=1
  for svc in bolabaden-nextjs autokuma; do
    running=0
    for node in "${NODES[@]}"; do
      if vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
        running=$((running + 1))
      fi
    done
    if [[ "$running" -lt "$MIN_NODES" ]]; then
      ok=0
      log "  ${svc}: ${running}/${MIN_NODES} nodes running — waiting"
      break
    fi
    yaml="$(vm_exec ci-node1 "cat ${YAML_PATH} 2>/dev/null" || true)"
    if ! echo "$yaml" | grep -q "${svc}-with-failover"; then
      ok=0
      log "  ${svc}: YAML block missing — waiting"
      break
    fi
    if ! echo "$yaml" | grep -q "https://${svc}.ci-node"; then
      ok=0
      log "  ${svc}: no peer https URLs yet — waiting"
      break
    fi
  done
  if [[ "$ok" -eq 1 ]]; then
    hz="$(vm_exec ci-node1 "wget -qO- --timeout=3 http://127.0.0.1:8082/healthz 2>/dev/null" || true)"
    if [[ "$hz" != "ok" ]]; then
      log "  failover-agent /healthz not ok yet (${hz:-empty}) — waiting"
      ok=0
    fi
  fi
  if [[ "$ok" -eq 1 ]]; then
    log "Tier-A ready (plurality + peer YAML + agent healthz)"
    exit 0
  fi
  sleep 10
done

die "Tier-A not ready within ${TIMEOUT}s — check compose-up, image sync, and failover-agent ensure logs on ci-node1"
