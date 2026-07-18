#!/usr/bin/env bash
# Module 5 prove: dry-run multi-A plan; optional live CF test zone.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
SYNC="${REPO_ROOT}/scripts/cloudflare_multi_record_ddns.py"
[[ -f "$SYNC" ]] || die "missing $SYNC"

FAIL=0
pass() { log "PASS: $*"; }
fail() { log "FAIL: $*"; FAIL=1; }

FIXTURE='{"ci-node1":"10.0.0.1","ci-node2":"10.0.0.2","ci-node3":"10.0.0.3","ci-node4":"10.0.0.4"}'
OUT="$(python3 "$SYNC" --dry-run --domain "${DOMAIN}" --names "whoami.${DOMAIN}" --node-ips "$FIXTURE")"
echo "$OUT" | python3 -c '
import json,sys
r=json.load(sys.stdin)
planned=r.get("planned") or {}
key=list(planned.keys())[0] if planned else ""
ips=planned.get(key) or []
assert r.get("dry_run") is True
assert len(ips)==4, ips
print("planned", key, ips)
' && pass "dry-run emits 4 A records for whoami.${DOMAIN}" || fail "dry-run plan invalid"

if [[ -n "${CF_API_TOKEN:-}" && -n "${CF_ZONE_ID:-}" && "${CF_LIVE_MULTI_DDNS:-}" == "1" ]]; then
  log "live CF upsert under ci-multi.${DOMAIN}"
  python3 "$SYNC" --domain "${DOMAIN}" --names "ci-multi.${DOMAIN}" \
    --node-ips "$FIXTURE" --token "$CF_API_TOKEN" --zone-id "$CF_ZONE_ID" \
    && pass "live multi-A upsert" || fail "live upsert failed"
else
  log "skip live CF (set CF_LIVE_MULTI_DDNS=1 with token+zone to enable)"
fi

[[ "$FAIL" -eq 0 ]] || die "prove-module5-ddns FAILED"
log "prove-module5-ddns ALL PASSED"
