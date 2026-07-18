#!/usr/bin/env bash
# Tear down 4 nodes and local state.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

BACKEND="$(backend)"
log "teardown backend=$BACKEND"

case "$BACKEND" in
  multipass)
    for name in "${NODES[@]}"; do
      if multipass info "$name" >/dev/null 2>&1; then
        log "deleting $name"
        multipass delete "$name" || true
      fi
    done
    multipass purge || true
    ;;
  qemu)
    for name in "${NODES[@]}"; do
      if virsh dominfo "$name" >/dev/null 2>&1; then
        log "destroy/undefine $name"
        virsh destroy "$name" 2>/dev/null || true
        virsh undefine "$name" --remove-all-storage --nvram 2>/dev/null \
          || virsh undefine "$name" --remove-all-storage 2>/dev/null \
          || true
      fi
      rm -f "${STATE_DIR}/images/${name}.qcow2" "${STATE_DIR}/images/${name}-seed.iso" \
        "${STATE_DIR}/images/${name}-user-data" "${STATE_DIR}/images/${name}-meta-data"
    done
    rm -rf "${STATE_DIR}/ssh"
    ;;
  dind)
    for name in "${NODES[@]}"; do
      if docker inspect "$name" >/dev/null 2>&1; then
        log "removing DinD $name"
        docker rm -f "$name" >/dev/null || true
      fi
      docker volume rm "failover-ci-${name}-docker" >/dev/null 2>&1 || true
    done
    if docker network inspect "${DIND_NETWORK}" >/dev/null 2>&1; then
      log "removing network ${DIND_NETWORK}"
      docker network rm "${DIND_NETWORK}" >/dev/null || true
    fi
    rm -f "${STATE_DIR}/dind-hosts"
    ;;
  *)
    log "no backend recorded — nothing to delete"
    ;;
esac

rm -f "${STATE_DIR}/repo-bundle.tgz" "${STATE_DIR}/node-ips.json" \
  "${STATE_DIR}/coredns-ips.txt" "${STATE_DIR}/headscale-preauth.key" \
  "${STATE_DIR}/tailscale-ips.txt" "${STATE_DIR}"/env-*.env "${STATE_DIR}/backend"
log "teardown complete"
