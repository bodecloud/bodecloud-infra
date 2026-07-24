#!/usr/bin/env bash
# Seed HA-critical images from host Docker into all DinD inner daemons (Hub 429 workaround).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ "$(backend)" == "dind" ]] || exit 0

REPO="${REPO_ROOT:-$(cd "${SCRIPT_DIR}/../.." && pwd)}"

build_host_agent_if_needed() {
  if [[ "${FAILOVER_CI_FORCE_AGENT_BUILD:-0}" == "1" ]] \
    || ! docker image inspect bolabaden/failover-agent:latest >/dev/null 2>&1; then
    log "building failover-agent on host (seed source)"
    docker build -t bolabaden/failover-agent:latest \
      -f "${REPO}/infra/Dockerfile.failover-agent" "${REPO}/infra"
    docker tag bolabaden/failover-agent:latest local/failover-agent:ci-test
  fi
  if ! docker image inspect local/failover-ci-probe:latest >/dev/null 2>&1; then
    log "building ci-probe on host (seed source)"
    docker build -t local/failover-ci-probe:latest \
      -f "${REPO}/arbitrary-scripts/failover-ci/compose/Dockerfile.ci-probe" \
      "${REPO}/arbitrary-scripts/failover-ci/compose"
  fi
}

log "seed-dind-images-from-host: HA-critical set → all DinD nodes"
build_host_agent_if_needed

mapfile -t SEED_REFS < <(ha_critical_images)
# Also tag variants compose may reference
SEED_REFS+=(
  "bolabaden/failover-agent:latest"
  "local/failover-ci-probe:latest"
)

for name in "${NODES[@]}"; do
  for ref in "${SEED_REFS[@]}"; do
    [[ -z "$ref" ]] && continue
    dind_seed_image_from_host "$name" "$ref" || true
  done
done
log "seed-dind-images-from-host complete"
