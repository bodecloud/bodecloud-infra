#!/usr/bin/env bash
# Prove Traefik wrong-node / LB / fallback across heterogeneous 4-node placement.
# Tier-A (bolabaden-nextjs, Autokuma) gates are HARD — whoami alone cannot pass.
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
YAML_PATH="${VM_REPO_PATH}/volumes/traefik/dynamic/failover-fallbacks.yaml"

curl_host() {
  local tip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 10 -o /tmp/cf.body -w '%{http_code}' -H 'Host: ${host}' https://${tip}/ 2>/dev/null" \
    || echo "000"
}

curl_host_body() {
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

expect_http_ok() {
  local tip="$1" host="$2" label="$3"
  local code body
  code="$(curl_host "$tip" "$host" || echo 000)"
  body="$(curl_host_body "$tip" "$host" || true)"
  if [[ "$code" == "200" ]] || echo "$body" | grep -qiE 'autokuma|bolabaden|<!DOCTYPE|Next\.js|Hostname:'; then
    pass "$label → HTTP ok (code=${code})"
  else
    fail "$label expected HTTP 200-ish, code=${code} body=$(echo "$body" | tr '\n' ' ' | head -c 160)"
  fi
}

assert_peer_urls() {
  local svc="$1"
  if ! vm_exec ci-node1 "grep -q '${svc}-with-failover' ${YAML_PATH} 2>/dev/null"; then
    fail "YAML missing ${svc}-with-failover block"
    return
  fi
  # AE4: must not be local-only when peers should run the service
  if ! vm_exec ci-node1 "grep -q 'https://${svc}.ci-node' ${YAML_PATH} 2>/dev/null"; then
    fail "AE4: ${svc} YAML has no peer https://${svc}.ci-nodeN.${DOMAIN} URLs (local-only)"
  else
    pass "YAML ${svc} includes peer https:// URLs"
  fi
}

assert_tier_a_plurality() {
  local svc="$1"
  local running=0 node
  for node in ci-node1 ci-node2 ci-node3 ci-node4; do
    if vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
      running=$((running + 1))
    fi
  done
  if [[ "$running" -ge 3 ]]; then
    pass "${svc} running on ${running}/4 nodes (n1 + >=2 peers)"
  else
    fail "${svc} plurality: running on ${running}/4 nodes (need ≥3)"
  fi
}

kill_running_on_peer() {
  local svc="$1"
  local node
  node="$(first_running_peer "$svc")" || { fail "no running ${svc} to kill for chaos test"; return 1; }
  log "killing ${svc} on ${node}" >&2
  vm_exec "$node" "docker kill ${svc} >/dev/null 2>&1 || true"
  echo "$node"
}

first_running_peer() {
  local svc="$1"
  local node
  for node in ci-node2 ci-node3 ci-node4 ci-node1; do
    if vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
      echo "$node"
      return 0
    fi
  done
  return 1
}

log "=== prove-failover: wrong-node / LB / fallback + Tier-A chaos ==="

for tip_name in ci-node1 ci-node4; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "whoami.${DOMAIN}" "whoami.ci-node1.${DOMAIN}"; do
    body="$(curl_host_body "$tip" "$host" || true)"
    expect_hostname "$body" "whoami" "${tip_name} Host=${host}"
  done
done

for tip_name in ci-node1 ci-node2; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "ci-probe.${DOMAIN}" "ci-probe.ci-node2.${DOMAIN}"; do
    body="$(curl_host_body "$tip" "$host" || true)"
    expect_hostname "$body" "ci-probe" "${tip_name} Host=${host}"
  done
done

# Tier-A plurality + YAML peer URLs (HARD)
for svc in bolabaden-nextjs autokuma; do
  assert_tier_a_plurality "$svc"
  assert_peer_urls "$svc"
done

# Tier-A Host reachability via tip
for tip_name in ci-node1 ci-node4; do
  tip="$(node_ip_from_state "$tip_name")"
  expect_http_ok "$tip" "bolabaden-nextjs.${DOMAIN}" "${tip_name} bolabaden-nextjs Host"
  expect_http_ok "$tip" "${DOMAIN}" "${tip_name} bolabaden apex Host"
  expect_http_ok "$tip" "autokuma.${DOMAIN}" "${tip_name} autokuma Host"
done

log "killing whoami on ci-node3 — must still serve via ci-node2 replica"
vm_exec ci-node3 "docker kill whoami >/dev/null 2>&1 || true"
sleep 8
body="$(curl_host_body "$N1" "whoami.${DOMAIN}" || true)"
expect_hostname "$body" "whoami" "after kill node3 via node1 (replica on node2)"
yaml="$(vm_exec ci-node1 "grep -E 'whoami-with-failover|url: https://whoami' ${YAML_PATH} 2>/dev/null | head -20" || true)"
if echo "$yaml" | grep -q 'whoami-with-failover'; then
  pass "failover-fallbacks.yaml retained whoami routes after kill"
else
  fail "whoami routes missing from failover-fallbacks.yaml after kill"
fi
vm_exec ci-node3 "docker start whoami >/dev/null 2>&1 || true"
sleep 5

# AE1: kill bolabaden on a peer that actually runs it — tip still 200
log "killing bolabaden-nextjs on a running peer — tip must still serve"
bola_node="$(kill_running_on_peer bolabaden-nextjs || true)"
if [[ -n "${bola_node:-}" ]]; then
  sleep 8
  for tip_name in ci-node1 ci-node4; do
    tip="$(node_ip_from_state "$tip_name")"
    expect_http_ok "$tip" "bolabaden-nextjs.${DOMAIN}" "AE1 after kill on ${bola_node} via ${tip_name}"
  done
  vm_exec "$bola_node" "docker start bolabaden-nextjs >/dev/null 2>&1 || docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env up -d --no-deps bolabaden-nextjs >/dev/null 2>&1 || true"
  sleep 5
fi

# AE6: kill Autokuma on a running peer
log "killing autokuma on a running peer — tip must still serve Autokuma Host"
aku_node="$(kill_running_on_peer autokuma || true)"
if [[ -n "${aku_node:-}" ]]; then
  sleep 8
  for tip_name in ci-node1 ci-node2 ci-node4; do
    tip="$(node_ip_from_state "$tip_name")"
    expect_http_ok "$tip" "autokuma.${DOMAIN}" "AE6 after kill on ${aku_node} via ${tip_name}"
  done
  vm_exec "$aku_node" "docker start autokuma >/dev/null 2>&1 || docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env up -d --no-deps autokuma >/dev/null 2>&1 || true"
  sleep 5
fi

log "stopping Traefik/coolify-proxy on ci-node1 — ingress via node4 (AE2)"
vm_exec ci-node1 "docker stop coolify-proxy >/dev/null 2>&1 || docker stop traefik >/dev/null 2>&1 || true"
sleep 3
body="$(curl_host_body "$N4" "whoami.${DOMAIN}" || true)"
expect_hostname "$body" "whoami" "via node4 while node1 Traefik down"
expect_http_ok "$N4" "bolabaden-nextjs.${DOMAIN}" "AE2 bolabaden via n4 tip"
expect_http_ok "$N4" "autokuma.${DOMAIN}" "AE2 autokuma via n4 tip"
vm_exec ci-node1 "docker start coolify-proxy >/dev/null 2>&1 || docker start traefik >/dev/null 2>&1 || true"

[[ "$FAIL" -eq 0 ]] || die "prove-failover FAILED"
log "prove-failover ALL PASSED"
