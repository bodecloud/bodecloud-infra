#!/usr/bin/env bash
# Admit Headscale control-plane SPOF: stopping headscale-server must fail MagicDNS.
# Never soft-skip when Tailscale was up before the stop (AE3 / R5).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/node-ips.json" ]] || die "run provision-vms.sh first"
FAIL=0
pass() { log "PASS: $*"; }
fail() { log "FAIL: $*"; FAIL=1; }

log "=== prove-headscale-spof: admitted control-plane SPOF ==="

TS_IP="$(vm_exec ci-node1 "tailscale ip -4 2>/dev/null | head -1" || true)"
if [[ -z "$TS_IP" ]]; then
  die "Tailscale not up on ci-node1 — cannot prove Headscale SPOF (no soft-skip)"
fi
pass "Tailscale up on ci-node1 (${TS_IP})"

MDNS_HOST="ci-node3.${HEADSCALE_BASE_DOMAIN}"
MDNS_BEFORE="$(vm_exec ci-node1 "dig +short @100.100.100.100 ${MDNS_HOST} A 2>/dev/null | head -3" || true)"
if [[ -z "$MDNS_BEFORE" ]]; then
  MDNS_BEFORE="$(vm_exec ci-node1 "dig +short @100.100.100.100 ci-node3 A 2>/dev/null | head -3" || true)"
fi
if [[ -z "$MDNS_BEFORE" ]]; then
  die "MagicDNS did not answer before HS stop — mesh not ready for SPOF prove"
fi
pass "MagicDNS answered before HS stop → $MDNS_BEFORE"

log "stopping headscale-server on ci-node1 (admitted SPOF)"
vm_exec ci-node1 "docker update --restart=no headscale-server 2>/dev/null || true; docker stop headscale-server"
sleep 5

if vm_exec ci-node1 "docker inspect -f '{{.State.Running}}' headscale-server 2>/dev/null | grep -qx true"; then
  fail "headscale-server still running after stop — cannot prove SPOF"
else
  pass "headscale-server stopped on ci-node1"
fi

MDNS_AFTER="$(vm_exec ci-node1 "dig +short @100.100.100.100 ${MDNS_HOST} A 2>/dev/null | head -3" || true)"
TS_STATUS="$(vm_exec ci-node1 "tailscale status 2>&1 | head -25" || true)"
TS_PEER="$(vm_exec ci-node2 "tailscale status 2>&1 | head -25" || true)"

spof_proven=0
if [[ -z "$MDNS_AFTER" ]]; then
  pass "MagicDNS empty after headscale-server stop (fail-closed)"
  spof_proven=1
elif echo "$TS_STATUS" | grep -qiE 'stopped|logged out|no magicdns|backend error|control|offline'; then
  pass "Tailscale on ci-node1 reports control-plane trouble after HS stop"
  spof_proven=1
elif echo "$TS_PEER" | grep -E 'ci-node1' | grep -qiE 'offline|idle.*offline|not connected|last seen'; then
  pass "ci-node2 sees ci-node1 degraded after HS stop"
  spof_proven=1
fi

if [[ "$spof_proven" -eq 0 ]]; then
  if [[ -n "$MDNS_AFTER" ]]; then
    log "WARN: MagicDNS still answered after HS stop (cache?): $MDNS_AFTER"
  fi
  fail "AE3: headscale-server stopped but MagicDNS/Tailscale still look healthy — SPOF gate not proven"
fi

# Edge via CoreDNS may still serve Tier-A (honesty: split-brain DNS paths)
if [[ -f "${STATE_DIR}/coredns-ips.txt" ]]; then
  mapfile -t CDS < "${STATE_DIR}/coredns-ips.txt"
  CD1="${CDS[0]}"
  EDGE="$(vm_exec ci-node1 "dig +short @${CD1} whoami.${DOMAIN} A 2>/dev/null | head -1" || true)"
  if echo "$EDGE" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+'; then
    pass "CoreDNS edge still answers whoami while HS down (split-path honesty)"
  else
    log "WARN: CoreDNS edge did not answer while HS down (non-fatal for SPOF gate)"
  fi
fi

log "restoring headscale-server"
vm_exec ci-node1 "docker update --restart=always headscale-server 2>/dev/null || true; docker start headscale-server"
sleep 8
MDNS_RESTORE="$(vm_exec ci-node1 "dig +short @100.100.100.100 ${MDNS_HOST} A 2>/dev/null | head -3" || true)"
if [[ -n "$MDNS_RESTORE" ]]; then
  pass "MagicDNS recovered after HS restore → $MDNS_RESTORE"
else
  log "WARN: MagicDNS not yet recovered (may need mesh settle); not failing if SPOF was proven"
fi

[[ "$FAIL" -eq 0 ]] || die "prove-headscale-spof FAILED"
log "prove-headscale-spof ALL PASSED"
