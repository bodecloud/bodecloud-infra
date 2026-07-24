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

NODE_IPS_FIXTURE='{"ci-node1":"10.0.3.2","ci-node2":"10.0.3.3","ci-node3":"10.0.3.4","ci-node4":"10.0.3.5"}'
if [[ -f "${STATE_DIR}/node-ips.json" ]]; then
  NODE_IPS_JSON="$(cat "${STATE_DIR}/node-ips.json")"
else
  NODE_IPS_JSON="$NODE_IPS_FIXTURE"
fi
[[ -f "${STATE_DIR}/coredns-ips.txt" ]] && CD1="$(head -1 "${STATE_DIR}/coredns-ips.txt")" || CD1=""

OUT="$(python3 "$SYNC" --dry-run --domain "${DOMAIN}" --names "whoami.${DOMAIN}" --node-ips "$NODE_IPS_JSON")"
echo "$OUT" | python3 -c '
import json,sys
r=json.load(sys.stdin)
planned=r.get("planned") or {}
key=list(planned.keys())[0] if planned else ""
ips=planned.get(key) or []
assert r.get("dry_run") is True
assert len(ips)>=3, ips
print("planned", key, ips)
' && pass "dry-run emits multi-A for whoami.${DOMAIN} (>=3 node IPs)" || fail "dry-run plan invalid"

# When mesh is up, CoreDNS multi-A must match Module 5 plan (CF parity).
if [[ -n "$CD1" && -f "${STATE_DIR}/node-ips.json" ]] \
  && command -v docker >/dev/null 2>&1 \
  && docker inspect ci-node1 >/dev/null 2>&1; then
  mapfile -t PLANNED_IPS < <(echo "$OUT" | python3 -c '
import json, sys
r = json.load(sys.stdin)
planned = r.get("planned") or {}
ips = next(iter(planned.values()), [])
print("\n".join(sorted(set(ips))))
')
  RESOLVED="$(vm_exec ci-node1 "dig +short @${CD1} whoami.${DOMAIN} A 2>/dev/null" | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$' | sort -u || true)"
  missing=0
  for ip in "${PLANNED_IPS[@]}"; do
    echo "$RESOLVED" | grep -qx "$ip" || missing=1
  done
  if [[ "$missing" -eq 0 && -n "$RESOLVED" ]]; then
    pass "CoreDNS multi-A matches Module5 plan"
  else
    fail "CoreDNS whoami.${DOMAIN} != Module5 plan (resolved=${RESOLVED})"
  fi
else
  log "skip CoreDNS↔Module5 compare (mesh not running — OK for validate-local)"
fi

if [[ -n "${CF_API_TOKEN:-}" && -n "${CF_ZONE_ID:-}" && "${CF_LIVE_MULTI_DDNS:-}" == "1" ]]; then
  log "live CF upsert under ci-multi.${DOMAIN}"
  python3 "$SYNC" --domain "${DOMAIN}" --names "ci-multi.${DOMAIN}" \
    --node-ips "$NODE_IPS_JSON" --token "$CF_API_TOKEN" --zone-id "$CF_ZONE_ID" \
    && pass "live multi-A upsert" || fail "live upsert failed"
else
  log "skip live CF (set CF_LIVE_MULTI_DDNS=1 with token+zone to enable)"
fi

[[ "$FAIL" -eq 0 ]] || die "prove-module5-ddns FAILED"
log "prove-module5-ddns ALL PASSED"
