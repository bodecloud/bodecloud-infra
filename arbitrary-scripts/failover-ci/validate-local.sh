#!/usr/bin/env bash
# Local gates that do not require Multipass/QEMU VMs.
# Run from repo root or from this directory.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "======== go test ./failover/ ========"
(cd "${ROOT}/infra" && go test ./failover/ -count=1)

echo "======== bash -n failover-ci scripts ========"
for s in "${CI_DIR}"/*.sh; do
  bash -n "$s"
done
# Explicitly cover new prove scripts (also matched by *.sh above)
for s in prove-headscale-spof.sh prove-failover.sh prove-dns.sh prove-matrix.sh prove-production-dns.sh prove-chaos-random.sh wait-tier-a-ready.sh sync-images-from-main.sh; do
  [[ -f "${CI_DIR}/${s}" ]] || { echo "FAIL: missing ${s}"; exit 1; }
  bash -n "${CI_DIR}/${s}"
done

echo "======== honesty contract (README) ========"
# Require explicit bans (wording that documents what we do NOT claim)
grep -Eqi 'Banned' "${CI_DIR}/README.md" \
  || { echo "FAIL: README missing Banned honesty rows"; exit 1; }
grep -Eq 'admitted SPOF|HA-critical curated' "${CI_DIR}/README.md" \
  || { echo "FAIL: README missing honesty contract language"; exit 1; }
# Affirmative overclaims (outside the honesty table) are banned
if grep -Eni 'we (run|prove|provide).*(dual-primary headscale|full stack ×4|full-stack ×4)' "${CI_DIR}/README.md"; then
  echo "FAIL: README contains affirmative overclaim"
  exit 1
fi
echo "README honesty contract OK"

echo "======== py_compile Module 5 syncer ========"
python3 -m py_compile "${ROOT}/scripts/cloudflare_multi_record_ddns.py"

echo "======== prove-module5-ddns (dry-run) ========"
bash "${CI_DIR}/prove-module5-ddns.sh"

echo "======== docker build failover-agent ========"
docker build -t local/failover-agent:test -f "${ROOT}/infra/Dockerfile.failover-agent" "${ROOT}/infra"

echo "======== docker build ci-probe ========"
docker build -t local/failover-ci-probe:test -f "${CI_DIR}/compose/Dockerfile.ci-probe" "${CI_DIR}/compose"

echo "======== CoreDNS dual-DNS smoke ========"
DOMAIN="${DOMAIN:-ci.bolabaden.test}"
TMP="${CI_DIR}/state/coredns-validate"
mkdir -p "$TMP"
cat > "${TMP}/${DOMAIN}.db" <<EOF
\$ORIGIN ${DOMAIN}.
\$TTL 30
@ IN SOA ns1.${DOMAIN}. admin.${DOMAIN}. ( 1 60 30 3600 30 )
@ IN NS ns1.${DOMAIN}.
ns1 IN A 127.0.0.1
* IN A 10.0.0.1
* IN A 10.0.0.2
*.ci-node3 IN A 10.0.0.3
ci-node3 IN A 10.0.0.3
EOF
sed -e "s/__DOMAIN__/${DOMAIN}/g" -e "s/__GOOGLE_DNS_1__/8.8.8.8/g" -e "s/__GOOGLE_DNS_2__/8.8.4.4/g" \
  "${CI_DIR}/coredns/Corefile.tmpl" > "${TMP}/Corefile"
sed -i 's|file /etc/coredns/zones/|file /zones/|g' "${TMP}/Corefile"
docker rm -f coredns-validate >/dev/null 2>&1 || true
docker run -d --name coredns-validate -p 1053:53/udp -p 1053:53/tcp \
  -v "${TMP}/Corefile:/Corefile:ro" -v "${TMP}:/zones:ro" \
  coredns/coredns:1.11.3 -conf /Corefile >/dev/null
sleep 1
ANS="$(dig +short @127.0.0.1 -p 1053 "whoami.${DOMAIN}" A | head -1)"
N3="$(dig +short @127.0.0.1 -p 1053 "whoami.ci-node3.${DOMAIN}" A | head -1)"
G="$(dig +short @8.8.8.8 google.com A | head -1)"
MG="$(dig +short @8.8.8.8 "whoami.${DOMAIN}" A || true)"
docker rm -f coredns-validate >/dev/null
[[ -n "$ANS" ]] || { echo "FAIL: CoreDNS whoami"; exit 1; }
[[ "$N3" == "10.0.0.3" ]] || { echo "FAIL: node-direct $N3"; exit 1; }
[[ -n "$G" ]] || { echo "FAIL: Google"; exit 1; }
if echo "$MG" | grep -qE '10\.0\.0\.'; then
  echo "FAIL: Google returned private mesh IP"
  exit 1
fi
echo "CoreDNS dual-DNS smoke OK (whoami=$ANS node3=$N3 google=$G)"

echo ""
echo "[failover-ci] validate-local ALL PASSED"
echo "Note: Tier-A chaos / Headscale SPOF / DinD mesh gates require ./run-all.sh (not this script)."
echo "Next (on a virt-capable host): ./run-all.sh"
