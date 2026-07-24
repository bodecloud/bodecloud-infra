#!/usr/bin/env bash
# Production DNS parity inside DinD CI (no live Cloudflare required).
#
# Production uses:
#   - cloudflare-multi-ddns → global FQDNs get multi-A (all Traefik node IPs)
#   - favonia cloudflare-ddns → node-direct ($TS_HOSTNAME.$DOMAIN, *.$TS_HOSTNAME.$DOMAIN)
#   - Docker embedded DNS (127.0.0.11) for compose service names on each node
#
# CI mirrors CF semantics in CoreDNS (provision-coredns.sh zone) and proves client
# ingress via DNS-resolved Traefik IPs — not only direct tip-IP curls with Host headers.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/node-ips.json" ]] || die "run provision-vms.sh first"
[[ -f "${STATE_DIR}/coredns-ips.txt" ]] || die "run provision-coredns.sh first"
[[ "$(backend)" == "dind" ]] || die "prove-production-dns is DinD-only (backend=$(backend))"

REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
SYNC="${REPO_ROOT}/scripts/cloudflare_multi_record_ddns.py"
[[ -f "$SYNC" ]] || die "missing $SYNC"

mapfile -t CDS < "${STATE_DIR}/coredns-ips.txt"
CD1="${CDS[0]}"
FAIL=0
pass() { log "PASS: $*"; }
fail() { log "FAIL: $*"; FAIL=1; }

dig_a() {
  vm_exec ci-node1 "dig +short @$1 $2 A 2>/dev/null | grep -E '^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$' | sort -u"
}

curl_via_ip_host() {
  local ip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 10 -o /tmp/prod-dns.body -w '%{http_code}' -H 'Host: ${host}' https://${ip}/ 2>/dev/null" \
    || echo "000"
}

curl_body_via_ip_host() {
  local ip="$1" host="$2"
  vm_exec ci-node1 "curl -sk --connect-timeout 10 -H 'Host: ${host}' https://${ip}/ 2>/dev/null | head -c 1500"
}

expected_node_ips() {
  python3 - "${STATE_DIR}/node-ips.json" <<'PY'
import json, sys
ips = json.load(open(sys.argv[1]))
for node in ("ci-node1", "ci-node2", "ci-node3", "ci-node4"):
    print(ips[node])
PY
}

log "=== prove-production-dns: CF multi-A + favonia node-direct + Docker DNS parity ==="

# --- 1. Module 5 dry-run plan must match live node-ips (production multi-ddns input) ---
log "--- Module 5 plan vs node-ips.json ---"
NODE_IPS_JSON="$(cat "${STATE_DIR}/node-ips.json")"
for global_name in "whoami.${DOMAIN}" "bolabaden-nextjs.${DOMAIN}" "autokuma.${DOMAIN}"; do
  OUT="$(python3 "$SYNC" --dry-run --domain "${DOMAIN}" --names "${global_name}" --node-ips "$NODE_IPS_JSON")"
  if python3 -c '
import json, sys
name, raw = sys.argv[1], sys.argv[2]
r = json.loads(raw)
ips = (r.get("planned") or {}).get(name) or []
assert r.get("dry_run") is True and len(ips) >= 3, ips
' "$global_name" "$OUT"; then
    pass "Module5 dry-run plans multi-A for ${global_name} (>=3 IPs)"
  else
    fail "Module5 dry-run plan invalid for ${global_name}"
  fi
done

# --- 2. CoreDNS zone matches Cloudflare-multi-ddns semantics (global multi-A) ---
log "--- CoreDNS global multi-A (cloudflare-multi-ddns equivalent) ---"
mapfile -t WANT_IPS < <(expected_node_ips | sort -u)
for global_name in "whoami.${DOMAIN}" "bolabaden-nextjs.${DOMAIN}" "autokuma.${DOMAIN}" "${DOMAIN}"; do
  mapfile -t GOT_IPS < <(dig_a "${CD1}" "${global_name}")
  missing=0
  for ip in "${WANT_IPS[@]}"; do
    if ! printf '%s\n' "${GOT_IPS[@]}" | grep -qx "$ip"; then
      missing=1
      break
    fi
  done
  if [[ "${#GOT_IPS[@]}" -ge 3 && "$missing" -eq 0 ]]; then
    pass "CoreDNS multi-A ${global_name} → ${GOT_IPS[*]}"
  else
    fail "CoreDNS multi-A ${global_name} expected ${WANT_IPS[*]}, got ${GOT_IPS[*]:-empty}"
  fi
done

# --- 3. Client ingress via DNS-resolved Traefik IP (production client path) ---
log "--- ingress via DNS-resolved Traefik IP (not hardcoded tip) ---"
for global_host in "whoami.${DOMAIN}" "bolabaden-nextjs.${DOMAIN}" "autokuma.${DOMAIN}"; do
  mapfile -t RESOLVED < <(dig_a "${CD1}" "${global_host}")
  [[ "${#RESOLVED[@]}" -gt 0 ]] || { fail "no A records for ${global_host}"; continue; }
  # Pick pseudo-random IP from multi-A set (seeded for reproducibility)
  seed="${PROD_DNS_SEED:-424242}"
  idx=$(( (seed + ${#global_host}) % ${#RESOLVED[@]} ))
  tip_ip="${RESOLVED[$idx]}"
  code="$(curl_via_ip_host "$tip_ip" "$global_host")"
  body="$(curl_body_via_ip_host "$tip_ip" "$global_host" || true)"
  case "$global_host" in
    whoami.*)
      if echo "$body" | grep -qiE 'Hostname:[[:space:]]*whoami|hostname[=:].*whoami'; then
        pass "DNS-resolved ${tip_ip} Host=${global_host} → whoami (${code})"
      else
        fail "DNS-resolved ${tip_ip} Host=${global_host} expected whoami code=${code}"
      fi
      ;;
    autokuma.*)
      if [[ "$code" == "200" ]] || echo "$body" | grep -qi autokuma; then
        pass "DNS-resolved ${tip_ip} Host=${global_host} → autokuma (${code})"
      else
        fail "DNS-resolved ${tip_ip} Host=${global_host} autokuma failed code=${code}"
      fi
      ;;
    *)
      if [[ "$code" == "200" ]] || echo "$body" | grep -qiE 'bolabaden|<!DOCTYPE|Next\.js'; then
        pass "DNS-resolved ${tip_ip} Host=${global_host} → bolabaden (${code})"
      else
        fail "DNS-resolved ${tip_ip} Host=${global_host} bolabaden failed code=${code}"
      fi
      ;;
  esac
done

# --- 4. favonia node-direct semantics ($TS_HOSTNAME.$DOMAIN / *.$TS_HOSTNAME.$DOMAIN) ---
log "--- favonia node-direct parity (single-node A) ---"
for node in ci-node1 ci-node2 ci-node3 ci-node4; do
  nip="$(node_ip_from_state "$node")"
  for host in "${node}.${DOMAIN}" "whoami.${node}.${DOMAIN}" "bolabaden-nextjs.${node}.${DOMAIN}"; do
    ans="$(dig_a "${CD1}" "${host}")"
    if echo "$ans" | grep -qx "$nip" && [[ "$(echo "$ans" | wc -l)" -eq 1 ]]; then
      pass "node-direct ${host} → ${nip} (favonia pattern)"
    else
      fail "node-direct ${host} expected only ${nip}, got: ${ans:-empty}"
    fi
  done
done

# --- 5. Production-like resolver (CoreDNS first, no MagicDNS) still resolves mesh names ---
log "--- production-like resolver profile (CoreDNS > Google, no MagicDNS) ---"
vm_exec ci-node1 "bash -s" <<EOS
set -euo pipefail
sudo tee /etc/resolv.conf >/dev/null <<EOF
# failover-ci production-dns prove: CoreDNS > Google (no MagicDNS)
nameserver ${CD1}
nameserver ${GOOGLE_DNS_1}
nameserver ${GOOGLE_DNS_2}
search ${DOMAIN}
options timeout:2 attempts:2
EOF
dig +short google.com A | grep -qE '^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$'
EOS
DEF_WHOAMI="$(vm_exec ci-node1 "dig +short whoami.${DOMAIN} A 2>/dev/null | grep -E '^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$' | head -1" || true)"
if echo "$DEF_WHOAMI" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$'; then
  pass "production-like resolv resolves whoami.${DOMAIN} → ${DEF_WHOAMI}"
  code="$(curl_via_ip_host "$DEF_WHOAMI" "whoami.${DOMAIN}")"
  body="$(curl_body_via_ip_host "$DEF_WHOAMI" "whoami.${DOMAIN}" || true)"
  if echo "$body" | grep -qi whoami; then
    pass "production-like resolv path serves whoami (${code})"
  else
    fail "production-like resolv path failed whoami HTTP code=${code}"
  fi
else
  fail "production-like resolv did not resolve whoami.${DOMAIN} (got: ${DEF_WHOAMI:-empty})"
fi
# Restore mesh resolver order for downstream proves
bash "${SCRIPT_DIR}/configure-resolvers.sh" >/dev/null

# --- 6. Docker embedded DNS (127.0.0.11) on each node Traefik network ---
log "--- Docker embedded DNS (compose service names) ---"
for node in ci-node1 ci-node2 ci-node3 ci-node4; do
  for svc in bolabaden-nextjs autokuma whoami; do
    # whoami only on n2/n3 after shape — skip absent
    if ! vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
      continue
    fi
    if vm_exec "$node" "docker exec traefik getent hosts ${svc} 2>/dev/null | grep -q ."; then
      pass "Docker DNS on ${node}: traefik resolves ${svc}"
    elif vm_exec "$node" "docker exec traefik sh -c 'nslookup ${svc} 127.0.0.11 2>/dev/null | grep -q Address'"; then
      pass "Docker DNS on ${node}: traefik nslookup ${svc} via 127.0.0.11"
    else
      fail "Docker DNS on ${node}: traefik cannot resolve ${svc}"
    fi
  done
done

log "--- honesty: CI vs production DNS ---"
log "  CoreDNS zone mirrors cloudflare-multi-ddns (global multi-A) + favonia (node-direct)."
log "  Traefik extra_hosts on DinD is a bridge-net shim; production uses public CF records instead."
log "  MagicDNS remains a Headscale path (prove-dns); production TS clients use it when enabled."

[[ "$FAIL" -eq 0 ]] || die "prove-production-dns FAILED"
log "prove-production-dns ALL PASSED"
