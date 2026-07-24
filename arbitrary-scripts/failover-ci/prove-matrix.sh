#!/usr/bin/env bash
# Comprehensive DinD ingress matrix: every Traefik tip × Tier-A/canary Host combinations.
# Run after shape-placement + wait-tier-a-ready (no chaos — see prove-failover.sh).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/node-ips.json" ]] || die "run provision-vms.sh first"
[[ "$(backend)" == "dind" ]] || die "prove-matrix is DinD-only (backend=$(backend))"

FAIL=0
pass() { log "PASS: $*"; }
fail() { log "FAIL: $*"; FAIL=1; }

VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"
YAML_PATH="${VM_REPO_PATH}/volumes/traefik/dynamic/failover-fallbacks.yaml"
REGISTRY="${VM_REPO_PATH}/volumes/placement/services.yaml"

curl_code() {
  local tip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 10 -o /dev/null -w '%{http_code}' -H 'Host: ${host}' https://${tip}/ 2>/dev/null" \
    || echo "000"
}

curl_body() {
  local tip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 10 -H 'Host: ${host}' https://${tip}/ 2>/dev/null | head -c 1500"
}

first_running_node() {
  local svc="$1"
  local node
  for node in "${NODES[@]}"; do
    if vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
      echo "$node"
      return 0
    fi
  done
  return 1
}

log "=== prove-matrix: placement + registry + YAML + tip×Host grid (DinD) ==="

# --- 1. Placement matrix (shape contract) ---
log "--- placement matrix ---"
declare -A EXPECT_RUNNING EXPECT_STOPPED
# whoami: n2,n3 only
EXPECT_RUNNING[whoami@ci-node2]=1
EXPECT_RUNNING[whoami@ci-node3]=1
EXPECT_STOPPED[whoami@ci-node1]=1
EXPECT_STOPPED[whoami@ci-node4]=1
# ci-probe: n4 only
EXPECT_RUNNING[ci-probe@ci-node4]=1
EXPECT_STOPPED[ci-probe@ci-node1]=1
EXPECT_STOPPED[ci-probe@ci-node2]=1
EXPECT_STOPPED[ci-probe@ci-node3]=1
# headscale-server: n1 only
EXPECT_RUNNING[headscale-server@ci-node1]=1
EXPECT_STOPPED[headscale-server@ci-node2]=1
EXPECT_STOPPED[headscale-server@ci-node3]=1
EXPECT_STOPPED[headscale-server@ci-node4]=1
# headscale UI client: n1,n2
EXPECT_RUNNING[headscale@ci-node1]=1
EXPECT_RUNNING[headscale@ci-node2]=1
EXPECT_STOPPED[headscale@ci-node3]=1
EXPECT_STOPPED[headscale@ci-node4]=1
# Tier-A: all nodes after ensure
for node in "${NODES[@]}"; do
  EXPECT_RUNNING[bolabaden-nextjs@"$node"]=1
  EXPECT_RUNNING[autokuma@"$node"]=1
  EXPECT_RUNNING[traefik@"$node"]=1
  EXPECT_RUNNING[failover-agent@"$node"]=1
done

for key in "${!EXPECT_RUNNING[@]}"; do
  svc="${key%%@*}"; node="${key#*@}"
  if vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
    pass "placement ${svc} running on ${node}"
  else
    fail "placement ${svc} should run on ${node}"
  fi
done
for key in "${!EXPECT_STOPPED[@]}"; do
  svc="${key%%@*}"; node="${key#*@}"
  if vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
    fail "placement ${svc} should NOT run on ${node} (shaped off)"
  else
    pass "placement ${svc} absent on ${node} (expected)"
  fi
done

# --- 2. Registry intentionally_stopped (R4) ---
log "--- registry R4 ---"
reg="$(vm_exec ci-node1 "cat ${REGISTRY} 2>/dev/null" || true)"
marker="$(vm_exec ci-node1 "cat ${VM_REPO_PATH}/volumes/placement/shape-intentional-stops.txt 2>/dev/null" || true)"
# Agent inventory may merge over shape registry on n1 (full stack); shape marker is authoritative.
if echo "$reg" | grep -q 'intentionally_stopped' || echo "$marker" | grep -qE '@'; then
  pass "registry/marker contains intentionally_stopped or shape marker pairs"
else
  fail "registry missing intentionally_stopped after shape (and no shape marker)"
fi
for pair in whoami@ci-node1 whoami@ci-node4 ci-probe@ci-node1; do
  svc="${pair%%@*}"; node="${pair#*@}"
  if echo "$reg" | grep -A5 "${node}:" | grep -q 'intentionally_stopped'; then
    pass "R4 ${svc}@${node} marked intentionally_stopped"
  elif echo "$marker" | grep -qxF "${svc}@${node}"; then
    pass "R4 ${svc}@${node} in shape marker"
  else
    fail "R4 ${svc}@${node} not marked intentionally_stopped"
  fi
done

# --- 3. Agent health + YAML peer URLs ---
log "--- agent + YAML audit ---"
hz="$(vm_exec ci-node1 "wget -qO- --timeout=5 http://127.0.0.1:8082/healthz 2>/dev/null" || true)"
if [[ "$hz" == "ok" ]]; then
  pass "failover-agent /healthz ok (strict ensure satisfied)"
else
  fail "failover-agent /healthz expected ok, got: ${hz:-empty}"
fi
yaml_has() {
  vm_exec ci-node1 "grep -q '$1' ${YAML_PATH} 2>/dev/null"
}
for svc in whoami ci-probe bolabaden-nextjs autokuma; do
  if yaml_has "${svc}-with-failover"; then
    pass "YAML has ${svc}-with-failover"
  else
    fail "YAML missing ${svc}-with-failover"
  fi
  if [[ "$svc" == "whoami" || "$svc" == "bolabaden-nextjs" || "$svc" == "autokuma" ]]; then
    if vm_exec ci-node1 "grep -q 'https://${svc}.ci-node' ${YAML_PATH} 2>/dev/null"; then
      pass "YAML ${svc} has peer https URLs"
    else
      fail "YAML ${svc} missing peer https:// URLs"
    fi
  fi
done
if vm_exec ci-node1 "grep 'url: http://bolabaden-nextjs:3000' ${YAML_PATH} 2>/dev/null | grep -qv 'ci-node'"; then
  if vm_exec ci-node1 "grep -q 'https://bolabaden-nextjs.ci-node' ${YAML_PATH} 2>/dev/null"; then
    pass "bolabaden YAML not local-only (peer URLs present)"
  else
    fail "AE4 bolabaden YAML appears local-only"
  fi
fi

# --- 4. Tip × Host ingress matrix ---
log "--- tip × Host ingress matrix ---"
TIPS=(ci-node1 ci-node2 ci-node3 ci-node4)

# whoami: should answer from n2/n3 backends via any tip
WHOAMI_HOSTS=("whoami.${DOMAIN}" "whoami.ci-node2.${DOMAIN}" "whoami.ci-node3.${DOMAIN}")
for tip_name in "${TIPS[@]}"; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "${WHOAMI_HOSTS[@]}"; do
    body="$(curl_body "$tip" "$host" || true)"
    if echo "$body" | grep -qiE 'Hostname:[[:space:]]*whoami|hostname[=:].*whoami'; then
      pass "matrix ${tip_name} Host=${host} → whoami"
    else
      fail "matrix ${tip_name} Host=${host} expected whoami, got: $(echo "$body" | tr '\n' ' ' | head -c 120)"
    fi
  done
done

# ci-probe: n4 placement
PROBE_HOSTS=("ci-probe.${DOMAIN}" "ci-probe.ci-node4.${DOMAIN}")
for tip_name in "${TIPS[@]}"; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "${PROBE_HOSTS[@]}"; do
    body="$(curl_body "$tip" "$host" || true)"
    if echo "$body" | grep -qiE 'Hostname:[[:space:]]*ci-probe|hostname[=:].*ci-probe'; then
      pass "matrix ${tip_name} Host=${host} → ci-probe"
    else
      fail "matrix ${tip_name} Host=${host} expected ci-probe"
    fi
  done
done

# bolabaden Tier-A
BOLA_HOSTS=("bolabaden-nextjs.${DOMAIN}" "${DOMAIN}" "bolabaden-nextjs.ci-node1.${DOMAIN}" "bolabaden-nextjs.ci-node2.${DOMAIN}")
for tip_name in "${TIPS[@]}"; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "${BOLA_HOSTS[@]}"; do
    code="$(curl_code "$tip" "$host")"
    body="$(curl_body "$tip" "$host" || true)"
    if [[ "$code" == "200" ]] || echo "$body" | grep -qiE '<!DOCTYPE|Next\.js|bolabaden'; then
      pass "matrix ${tip_name} Host=${host} → bolabaden HTTP ok (${code})"
    else
      fail "matrix ${tip_name} Host=${host} bolabaden failed code=${code}"
    fi
  done
done

# autokuma Tier-A
AKU_HOSTS=("autokuma.${DOMAIN}" "autokuma.ci-node1.${DOMAIN}" "autokuma.ci-node3.${DOMAIN}")
for tip_name in "${TIPS[@]}"; do
  tip="$(node_ip_from_state "$tip_name")"
  for host in "${AKU_HOSTS[@]}"; do
    code="$(curl_code "$tip" "$host")"
    body="$(curl_body "$tip" "$host" || true)"
    if [[ "$code" == "200" ]] || echo "$body" | grep -qi 'autokuma'; then
      pass "matrix ${tip_name} Host=${host} → autokuma HTTP ok (${code})"
    else
      fail "matrix ${tip_name} Host=${host} autokuma failed code=${code}"
    fi
  done
done

# --- 5. Negative: shaped-off local Host should still peer-forward (not 404 on global) ---
log "--- negative / shaped-off sanity ---"
tip="$(node_ip_from_state ci-node1)"
code="$(curl_code "$tip" "whoami.ci-node1.${DOMAIN}")"
body="$(curl_body "$tip" "whoami.ci-node1.${DOMAIN}" || true)"
if echo "$body" | grep -qi whoami; then
  pass "whoami.ci-node1 Host peer-forwards despite n1 shaped off"
else
  fail "whoami.ci-node1 Host should peer-forward (code=${code})"
fi

[[ "$FAIL" -eq 0 ]] || die "prove-matrix FAILED"
log "prove-matrix ALL PASSED ($((${#TIPS[@]})) tips × multi-Host grid)"
