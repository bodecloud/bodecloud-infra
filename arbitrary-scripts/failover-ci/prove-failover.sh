#!/usr/bin/env bash
# Prove Traefik wrong-node / LB / fallback across heterogeneous 4-node placement.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/node-ips.json" ]] || die "run provision-vms.sh first"
FAIL=0
pass() { log "PASS: $*"; }
fail() { log "FAIL: $*"; FAIL=1; }

N1="$(node_ip_from_state ci-node1)"
N4="$(node_ip_from_state ci-node4)"
VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"

curl_host() {
  local tip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 10 -H 'Host: ${host}' https://${tip}/ 2>/dev/null | head -c 2000"
}

expect_hostname() {
  local body="$1" want="$2" label="$3"
  if echo "$body" | grep -qiE "Hostname:[[:space:]]*${want}|hostname[=:].*${want}"; then
    pass "$label → hostname ${want}"
  else
    fail "$label expected hostname ${want}, body=$(echo "$body" | tr '\n' ' ' | head -c 200)"
  fi
}

log "=== prove-failover: wrong-node / LB / fallback ==="

for tip_name in ci-node1 ci-node4; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "whoami.${DOMAIN}" "whoami.ci-node1.${DOMAIN}"; do
    body="$(curl_host "$tip" "$host" || true)"
    expect_hostname "$body" "whoami" "${tip_name} Host=${host}"
  done
done

for tip_name in ci-node1 ci-node2; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "ci-probe.${DOMAIN}" "ci-probe.ci-node2.${DOMAIN}"; do
    body="$(curl_host "$tip" "$host" || true)"
    expect_hostname "$body" "ci-probe" "${tip_name} Host=${host}"
  done
done

log "killing whoami on ci-node3 — must still serve via ci-node2 replica"
vm_exec ci-node3 "docker kill whoami || true"
sleep 8
body="$(curl_host "$N1" "whoami.${DOMAIN}" || true)"
expect_hostname "$body" "whoami" "after kill node3 via node1 (replica on node2)"
yaml="$(vm_exec ci-node1 "grep -E 'whoami-with-failover|url: https://whoami' ${VM_REPO_PATH}/volumes/traefik/dynamic/failover-fallbacks.yaml 2>/dev/null | head -20" || true)"
if echo "$yaml" | grep -q 'whoami-with-failover'; then
  pass "failover-fallbacks.yaml retained whoami routes after kill"
else
  fail "whoami routes missing from failover-fallbacks.yaml after kill"
fi
vm_exec ci-node3 "docker start whoami || true"
sleep 5

log "stopping Traefik/coolify-proxy on ci-node1 — ingress via node4"
vm_exec ci-node1 "docker stop coolify-proxy 2>/dev/null || docker stop traefik 2>/dev/null || true"
sleep 3
body="$(curl_host "$N4" "whoami.${DOMAIN}" || true)"
expect_hostname "$body" "whoami" "via node4 while node1 Traefik down"
vm_exec ci-node1 "docker start coolify-proxy 2>/dev/null || docker start traefik 2>/dev/null || true"

[[ "$FAIL" -eq 0 ]] || die "prove-failover FAILED"
log "prove-failover ALL PASSED"
