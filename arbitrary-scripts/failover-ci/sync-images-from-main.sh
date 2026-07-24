#!/usr/bin/env bash
# Stream Docker images from MAIN_HOST (ci-node1) to peer DinD nodes.
# Avoids Docker Hub 429s when bringing full compose up on peers.
#
# Env:
#   FAILOVER_CI_IMAGE_MAX_MB  Skip images larger than this (default 450). 0 = no limit.
#   FAILOVER_CI_SYNC_TARGETS  Space-separated peer names (default: all except MAIN).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

MAIN="${MAIN_HOST:-ci-node1}"
MAX_MB="${FAILOVER_CI_IMAGE_MAX_MB:-450}"
VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"

TARGETS=()
if [[ -n "${FAILOVER_CI_SYNC_TARGETS:-}" ]]; then
  # shellcheck disable=SC2206
  TARGETS=(${FAILOVER_CI_SYNC_TARGETS})
else
  for n in "${NODES[@]}"; do
    [[ "$n" == "$MAIN" ]] && continue
    TARGETS+=("$n")
  done
fi

[[ "$(backend)" == "dind" ]] || die "sync-images-from-main is DinD-only (backend=$(backend))"

log "listing images on ${MAIN} (max ${MAX_MB} MiB; 0=unlimited)"
LIST_FILE="${STATE_DIR}/sync-images.list"
vm_exec "$MAIN" "bash -s" <<EOS >"${LIST_FILE}"
set -euo pipefail
MAX_MB=${MAX_MB}
docker images --format '{{.Repository}}:{{.Tag}} {{.Size}}' | while read -r ref size; do
  [[ "\$ref" == *"<none>"* ]] && continue
  [[ "\$ref" == *":<none>" ]] && continue
  if [[ "\$MAX_MB" != "0" ]]; then
    num=\$(printf '%s' "\$size" | sed -E 's/([0-9.]+).*/\\1/')
    unit=\$(printf '%s' "\$size" | sed -E 's/[0-9.]+//')
    mb=0
    case "\$unit" in
      GB|GiB) mb=\$(awk -v n="\$num" 'BEGIN{printf "%d", n*1024}') ;;
      MB|MiB) mb=\$(awk -v n="\$num" 'BEGIN{printf "%d", n}') ;;
      KB|KiB|B) mb=1 ;;
      *) mb=99999 ;;
    esac
    if [[ "\$mb" -gt "\$MAX_MB" ]]; then
      echo "[failover-ci] skip >\${MAX_MB}MiB \$ref (\$size)" >&2
      continue
    fi
  fi
  printf '%s\n' "\$ref"
done
EOS

# Always include HA-critical refs even if over size cap (honesty: curated set, not full×4)
while IFS= read -r c; do
  [[ -z "$c" ]] && continue
  # Prefer :latest and untagged variants present on MAIN
  for cand in "$c" "${c%:latest}" "${c##*/}"; do
    if vm_exec "$MAIN" "docker image inspect $(printf '%q' "$cand") >/dev/null 2>&1"; then
      grep -qxF "$cand" "${LIST_FILE}" 2>/dev/null || echo "$cand" >>"${LIST_FILE}"
      break
    fi
  done
  # Also try without docker.io/ prefix
  bare="${c#docker.io/}"
  if [[ "$bare" != "$c" ]] && vm_exec "$MAIN" "docker image inspect $(printf '%q' "$bare") >/dev/null 2>&1"; then
    grep -qxF "$bare" "${LIST_FILE}" 2>/dev/null || echo "$bare" >>"${LIST_FILE}"
  fi
done < <(ha_critical_images)
log "HA-critical peer set (not full stack ×4): $(ha_critical_services)"
sort -u -o "${LIST_FILE}" "${LIST_FILE}"
mapfile -t IMAGES < "${LIST_FILE}"
[[ "${#IMAGES[@]}" -gt 0 ]] || die "no images selected for sync"

log "syncing ${#IMAGES[@]} images → ${TARGETS[*]}"
head -40 "${LIST_FILE}"
[[ "${#IMAGES[@]}" -gt 40 ]] && log "... ($(( ${#IMAGES[@]} - 40 )) more)"

# Push list into MAIN for xargs docker save
docker exec -i "$MAIN" sh -lc 'cat > /tmp/failover-ci-sync-images.list' < "${LIST_FILE}"

sync_one() {
  local peer="$1"
  log "loading images into ${peer}"
  # Stream save→load via host; avoids DinD /tmp docker-cp quirks and Hub pulls
  if ! docker exec "$MAIN" sh -lc \
    'xargs -a /tmp/failover-ci-sync-images.list -r docker save' \
    | docker exec -i "$peer" docker load; then
    warn "bulk save/load failed on ${peer}; retrying one-by-one"
    local ref
    for ref in "${IMAGES[@]}"; do
      docker exec "$MAIN" docker save "$ref" \
        | docker exec -i "$peer" docker load \
        || warn "failed to load ${ref} on ${peer}"
    done
  fi
  local count
  count="$(vm_exec "$peer" 'docker images -q | wc -l' | tr -d '[:space:]')"
  log "${peer}: ${count} images present"
}

for peer in "${TARGETS[@]}"; do
  sync_one "$peer"
  free -h | head -2 || true
  df -h / | tail -1 || true
done

log "image sync complete"
log "next: FAILOVER_CI_MINIMAL=0 ./compose-up-critical-full.sh ${TARGETS[*]}"
# Or full up with pull=never on peers
log "      (peers should use --pull=never after sync)"
