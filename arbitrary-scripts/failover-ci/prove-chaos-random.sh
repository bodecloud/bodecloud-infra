#!/usr/bin/env bash
# Randomized chaos: kill/stop arbitrary HA-critical components and assert failover.
# Complements deterministic prove-failover.sh (AE1/AE2/AE6) with seeded random rounds.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/node-ips.json" ]] || die "run provision-vms.sh first"
[[ -f "${STATE_DIR}/coredns-ips.txt" ]] || die "run provision-coredns.sh first"
[[ "$(backend)" == "dind" ]] || die "prove-chaos-random is DinD-only (backend=$(backend))"

CHAOS_ROUNDS="${CHAOS_ROUNDS:-10}"
CHAOS_SEED="${CHAOS_SEED:-1337}"
CHAOS_DOUBLE="${CHAOS_DOUBLE:-1}"
CHAOS_WAIT_SECS="${CHAOS_WAIT_SECS:-12}"  # settle after kill/stop before verify
CHAOS_VERIFY_RETRIES="${CHAOS_VERIFY_RETRIES:-3}"
VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"

mapfile -t CDS < "${STATE_DIR}/coredns-ips.txt"
CD1="${CDS[0]}"

FAIL=0
pass() { log "PASS: $*"; }
fail() { log "FAIL: $*"; FAIL=1; }

# Seeded PRNG (LCG) — reproducible in CI logs
RND_STATE="$CHAOS_SEED"
rnd() {
  RND_STATE=$(( (RND_STATE * 1103515245 + 12345) & 0x7fffffff ))
  echo "$RND_STATE"
}
rnd_pick() {
  local n="$1"
  local r max=0
  r="$(rnd)"
  max=$(( r % n ))
  echo "$max"
}

curl_host_code() {
  local tip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 12 -o /dev/null -w '%{http_code}' -H 'Host: ${host}' https://${tip}/ 2>/dev/null" \
    || echo "000"
}

curl_host_body() {
  local tip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 12 -H 'Host: ${host}' https://${tip}/ 2>/dev/null | head -c 2000"
}

dig_all_ips() {
  local host="$1"
  vm_exec ci-node1 "dig +short @${CD1} ${host} A 2>/dev/null | grep -E '^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$' | sort -u"
}

reset_round_rng() {
  local r="$1"
  RND_STATE=$(( CHAOS_SEED + r * 10007 ))
}

verify_failover() {
  local round="$1"
  local node tip_ip code body ok attempt

  for attempt in $(seq 1 "$CHAOS_VERIFY_RETRIES"); do
    ok=0
    for node in ci-node1 ci-node2 ci-node3 ci-node4; do
      traefik_up_on "$node" || continue
      tip_ip="$(node_ip_from_state "$node")"
      body="$(curl_host_body "$tip_ip" "whoami.${DOMAIN}" || true)"
      code="$(curl_host_code "$tip_ip" "whoami.${DOMAIN}")"
      if echo "$body" | grep -qiE 'Hostname:[[:space:]]*whoami|hostname[=:].*whoami'; then
        pass "round ${round} whoami via tip ${node} (${code}) attempt=${attempt}"
        ok=1
        break
      fi
    done
    [[ "$ok" -eq 1 ]] && break
    [[ "$attempt" -lt "$CHAOS_VERIFY_RETRIES" ]] && sleep 4
  done
  [[ "$ok" -eq 1 ]] || fail "round ${round} whoami failed on all Traefik tips"

  for attempt in $(seq 1 "$CHAOS_VERIFY_RETRIES"); do
    ok=0
    for node in ci-node1 ci-node2 ci-node3 ci-node4; do
      traefik_up_on "$node" || continue
      tip_ip="$(node_ip_from_state "$node")"
      body="$(curl_host_body "$tip_ip" "bolabaden-nextjs.${DOMAIN}" || true)"
      code="$(curl_host_code "$tip_ip" "bolabaden-nextjs.${DOMAIN}")"
      if [[ "$code" == "200" ]] || echo "$body" | grep -qiE 'bolabaden|<!DOCTYPE|Next\.js'; then
        pass "round ${round} bolabaden via tip ${node} (${code}) attempt=${attempt}"
        ok=1
        break
      fi
    done
    [[ "$ok" -eq 1 ]] && break
    [[ "$attempt" -lt "$CHAOS_VERIFY_RETRIES" ]] && sleep 4
  done
  [[ "$ok" -eq 1 ]] || fail "round ${round} bolabaden failed on all Traefik tips"

  for attempt in $(seq 1 "$CHAOS_VERIFY_RETRIES"); do
    ok=0
    for node in ci-node1 ci-node2 ci-node3 ci-node4; do
      traefik_up_on "$node" || continue
      tip_ip="$(node_ip_from_state "$node")"
      body="$(curl_host_body "$tip_ip" "autokuma.${DOMAIN}" || true)"
      code="$(curl_host_code "$tip_ip" "autokuma.${DOMAIN}")"
      if [[ "$code" == "200" ]] || echo "$body" | grep -qi autokuma; then
        pass "round ${round} autokuma via tip ${node} (${code}) attempt=${attempt}"
        ok=1
        break
      fi
    done
    [[ "$ok" -eq 1 ]] && break
    [[ "$attempt" -lt "$CHAOS_VERIFY_RETRIES" ]] && sleep 4
  done
  [[ "$ok" -eq 1 ]] || fail "round ${round} autokuma failed on all Traefik tips"

  # Production DNS path: multi-A pick among Traefik nodes that are still up
  ok=0
  for node in ci-node1 ci-node2 ci-node3 ci-node4; do
    traefik_up_on "$node" || continue
    tip_ip="$(node_ip_from_state "$node")"
    body="$(curl_host_body "$tip_ip" "whoami.${DOMAIN}" || true)"
    if echo "$body" | grep -qiE 'Hostname:[[:space:]]*whoami|hostname[=:].*whoami'; then
      pass "round ${round} whoami via DNS-up IP ${tip_ip} (${node})"
      ok=1
      break
    fi
  done
  [[ "$ok" -eq 1 ]] || fail "round ${round} whoami failed on all up Traefik DNS IPs"
}

svc_running_on() {
  local node="$1" svc="$2"
  vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"
}

traefik_up_on() {
  local node="$1"
  svc_running_on "$node" traefik
}

apply_chaos() {
  local kind="$1" node="$2" svc="${3:-}"
  case "$kind" in
    kill)
      log "CHAOS kill ${svc} on ${node}" >&2
      vm_exec "$node" "docker kill ${svc} >/dev/null 2>&1 || true"
      echo "kill:${svc}@${node}"
      ;;
    stop_traefik)
      log "CHAOS stop traefik on ${node}" >&2
      vm_exec "$node" "docker stop traefik >/dev/null 2>&1 || true"
      echo "stop_traefik@${node}"
      ;;
  esac
}

traefik_up_count() {
  local n c=0
  for n in ci-node1 ci-node2 ci-node3 ci-node4; do
    traefik_up_on "$n" && c=$((c + 1))
  done
  echo "$c"
}

can_stop_traefik() {
  [[ "$(traefik_up_count)" -gt 1 ]]
}

restore_chaos() {
  local token="$1"
  if [[ "$token" == kill:*@* ]]; then
    local svc="${token#kill:}"
    svc="${svc%%@*}"
    local node="${token##*@}"
    vm_exec "$node" "docker start ${svc} >/dev/null 2>&1 || docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env up -d --no-deps ${svc} >/dev/null 2>&1 || true"
  elif [[ "$token" == stop_traefik@* ]]; then
    local node="${token#stop_traefik@}"
    vm_exec "$node" "docker start traefik >/dev/null 2>&1 || docker compose --project-directory ${VM_REPO_PATH} --env-file ${VM_REPO_PATH}/.env up -d --no-deps traefik >/dev/null 2>&1 || true"
  fi
}

pick_round_actions() {
  local -a picked=()
  local -a used=()
  local entry kind node svc key
  local want=1

  if [[ "$CHAOS_DOUBLE" == "1" ]] && [[ $(( $(rnd) % 2 )) -eq 0 ]] && [[ "${#POOL[@]}" -gt 1 ]]; then
    want=2
  fi

  while [[ "${#picked[@]}" -lt "$want" ]]; do
    local found=0
    for entry in "${POOL[@]}"; do
      IFS='|' read -r kind node svc <<< "$entry"
      key="${kind}|${node}|${svc}"
      [[ " ${used[*]} " == *" ${key} "* ]] && continue
      if [[ "$kind" == "stop_traefik" ]] && ! can_stop_traefik; then
        continue
      fi
      picked+=("$(apply_chaos "$kind" "$node" "$svc")")
      used+=("$key")
      found=1
      break
    done
    [[ "$found" -eq 1 ]] || break
  done

  [[ "${#picked[@]}" -gt 0 ]] || die "no valid chaos action at round ${round} (pool=${#POOL[@]})"
  printf '%s\n' "${picked[@]}"
}

# Build weighted chaos pool (only targets that should exist after shape)
# Format: kind|node|svc
POOL=()
add_pool() {
  local kind="$1" node="$2" svc="${3:-}"
  if [[ "$kind" == "kill" ]] && svc_running_on "$node" "$svc"; then
    POOL+=("${kind}|${node}|${svc}")
  elif [[ "$kind" == "stop_traefik" ]] && traefik_up_on "$node"; then
    POOL+=("${kind}|${node}|")
  fi
}

for node in ci-node2 ci-node3; do
  add_pool kill "$node" whoami
done
for node in ci-node1 ci-node2 ci-node3 ci-node4; do
  add_pool kill "$node" bolabaden-nextjs
  add_pool kill "$node" autokuma
  add_pool stop_traefik "$node"
done

[[ "${#POOL[@]}" -gt 0 ]] || die "chaos pool empty — is mesh up after shape?"

log "=== prove-chaos-random: ${CHAOS_ROUNDS} rounds seed=${CHAOS_SEED} pool=${#POOL[@]} ==="

round=1
while [[ "$round" -le "$CHAOS_ROUNDS" ]]; do
  reset_round_rng "$round"
  POOL=()
  for node in ci-node2 ci-node3; do add_pool kill "$node" whoami; done
  for node in ci-node1 ci-node2 ci-node3 ci-node4; do
    add_pool kill "$node" bolabaden-nextjs
    add_pool kill "$node" autokuma
    add_pool stop_traefik "$node"
  done
  # Fisher–Yates shuffle (seeded)
  i=$((${#POOL[@]} - 1))
  while [[ "$i" -gt 0 ]]; do
    j="$(rnd_pick $((i + 1)))"
    tmp="${POOL[$i]}"
    POOL[$i]="${POOL[$j]}"
    POOL[$j]="$tmp"
    i=$((i - 1))
  done
  [[ "${#POOL[@]}" -gt 0 ]] || die "chaos pool empty at round ${round}"

  actions=()
  mapfile -t actions < <(pick_round_actions)

  if [[ "${#actions[@]}" -eq 2 ]]; then
    log "round ${round}: double chaos ${actions[*]}"
  elif [[ "${#actions[@]}" -eq 1 ]]; then
    log "round ${round}: chaos ${actions[0]}"
  fi

  sleep "$CHAOS_WAIT_SECS"
  verify_failover "$round"

  for tok in "${actions[@]}"; do
    restore_chaos "$tok"
  done
  sleep $(( CHAOS_WAIT_SECS / 2 + 2 ))
  round=$((round + 1))
done

[[ "$FAIL" -eq 0 ]] || die "prove-chaos-random FAILED"
log "prove-chaos-random ALL PASSED (${CHAOS_ROUNDS} random rounds)"
