#!/usr/bin/env bash
# Prove dual DNS: Tailscale/Headscale MagicDNS priority + Google path.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/coredns-ips.txt" ]] || die "run provision-coredns.sh first"
mapfile -t CDS < "${STATE_DIR}/coredns-ips.txt"
CD1="${CDS[0]}"
FAIL=0

pass() { log "PASS: $*"; }
fail() { log "FAIL: $*"; FAIL=1; }

run_dig() {
  vm_exec ci-node1 "dig +short $* 2>/dev/null | head -5"
}

log "=== prove-dns: Dual DNS comprehension ==="

ANS="$(run_dig "@${CD1}" "whoami.${DOMAIN}" A || true)"
if echo "$ANS" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+'; then
  pass "CoreDNS answers whoami.${DOMAIN} → $ANS"
else
  fail "CoreDNS did not answer whoami.${DOMAIN} (got: $ANS)"
fi

ANS3="$(run_dig "@${CD1}" "whoami.ci-node3.${DOMAIN}" A || true)"
N3="$(node_ip_from_state ci-node3)"
if echo "$ANS3" | grep -q "$N3"; then
  pass "whoami.ci-node3.${DOMAIN} → ${N3}"
else
  fail "whoami.ci-node3.${DOMAIN} expected ${N3}, got: $ANS3"
fi

# Tier-A DNS plurality (node-direct wildcards cover bolabaden/autokuma)
for svc in bolabaden-nextjs autokuma; do
  GANS_TIER="$(run_dig "@${CD1}" "${svc}.${DOMAIN}" A || true)"
  if echo "$GANS_TIER" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+'; then
    pass "CoreDNS answers ${svc}.${DOMAIN} → $GANS_TIER"
  else
    fail "CoreDNS did not answer ${svc}.${DOMAIN} (got: $GANS_TIER)"
  fi
  for node in ci-node1 ci-node2 ci-node3; do
    nip="$(node_ip_from_state "$node")"
    # Only require node-direct when the service is running on that node
    if vm_exec "$node" "docker inspect -f '{{.State.Running}}' ${svc} 2>/dev/null | grep -qx true"; then
      NANS="$(run_dig "@${CD1}" "${svc}.${node}.${DOMAIN}" A || true)"
      if echo "$NANS" | grep -q "$nip"; then
        pass "${svc}.${node}.${DOMAIN} → ${nip}"
      else
        fail "${svc}.${node}.${DOMAIN} expected ${nip}, got: $NANS"
      fi
    fi
  done
done

GANS="$(run_dig "@${GOOGLE_DNS_1}" google.com A || true)"
if echo "$GANS" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+'; then
  pass "Google ${GOOGLE_DNS_1} resolves google.com"
else
  fail "Google DNS failed for google.com: $GANS"
fi

MESH_VIA_GOOGLE="$(run_dig "@${GOOGLE_DNS_1}" "whoami.${DOMAIN}" A || true)"
PRIV="$(node_ip_from_state ci-node1)|$(node_ip_from_state ci-node3)"
if echo "$MESH_VIA_GOOGLE" | grep -qE "$PRIV"; then
  fail "Google returned private Traefik IP for mesh name: $MESH_VIA_GOOGLE"
else
  pass "Google does not return private Traefik IPs for whoami.${DOMAIN}"
fi

TS_IP="$(vm_exec ci-node1 "tailscale ip -4 2>/dev/null | head -1" || true)"
if [[ -n "$TS_IP" ]]; then
  MDNS_HOST="ci-node3.${HEADSCALE_BASE_DOMAIN}"
  MDNS_ANS="$(run_dig "@100.100.100.100" "${MDNS_HOST}" A || true)"
  if [[ -z "$MDNS_ANS" ]]; then
    MDNS_ANS="$(run_dig "@100.100.100.100" "ci-node3.${DOMAIN}" A || true)"
  fi
  if [[ -z "$MDNS_ANS" ]]; then
    MDNS_ANS="$(run_dig "@100.100.100.100" "ci-node3" A || true)"
  fi
  if [[ -n "$MDNS_ANS" ]]; then
    pass "MagicDNS @100.100.100.100 answered → $MDNS_ANS"
  else
    fail "Tailscale is up but MagicDNS did not answer ci-node3 (priority gate — no soft-skip)"
  fi
else
  # DinD mesh requires Tailscale; soft-skip only when backend is not DinD
  if [[ "$(backend)" == "dind" ]]; then
    fail "DinD requires Tailscale — MagicDNS hard gate cannot soft-skip"
  else
    log "WARN: tailscale not up — skipping MagicDNS hard gate (non-DinD)"
  fi
fi

DEF="$(vm_exec ci-node1 "getent hosts google.com 2>/dev/null | head -1 || dig +short google.com A | head -1" || true)"
if echo "$DEF" | grep -Eq '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+|([0-9a-fA-F]{0,4}:){2,}[0-9a-fA-F]{0,4}'; then
  pass "default resolv resolves google.com"
else
  fail "default resolv failed for google.com: $DEF"
fi

[[ "$FAIL" -eq 0 ]] || die "prove-dns FAILED"
log "prove-dns ALL PASSED"
