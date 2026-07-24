#!/usr/bin/env bash
# Full local orchestration (no GHA, no git push).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

SKIP_COMPOSE="${SKIP_COMPOSE:-0}"
SKIP_PROVE="${SKIP_PROVE:-0}"

run() {
  echo ""
  echo "======== $* ========"
  bash "$@"
}

run ./provision-vms.sh
run ./provision-test-env.sh
run ./provision-coredns.sh
run ./configure-resolvers.sh

if [[ "$SKIP_COMPOSE" != "1" ]]; then
  run ./compose-up-all.sh
  run ./shape-placement.sh
  run ./provision-mesh.sh
  run ./configure-resolvers.sh
  run ./wait-tier-a-ready.sh
fi

if [[ "$SKIP_PROVE" != "1" ]]; then
  run ./prove-matrix.sh
  run ./prove-dns.sh
  run ./prove-production-dns.sh
  run ./prove-failover.sh
  run ./prove-chaos-random.sh
  run ./prove-headscale-spof.sh
  run ./prove-module5-ddns.sh
fi

echo ""
echo "[failover-ci] run-all COMPLETE"
