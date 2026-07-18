#!/usr/bin/env bash
# Configure resolver order: MagicDNS > CoreDNS mesh > Google 8.8.8.8 / 8.8.4.4
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/coredns-ips.txt" ]] || die "run provision-coredns.sh first"
mapfile -t CDS < "${STATE_DIR}/coredns-ips.txt"
CD1="${CDS[0]:-}"
CD2="${CDS[1]:-}"
[[ -n "$CD1" ]] || die "empty CoreDNS IP list"

backend_is_dind=0
[[ "$(backend)" == "dind" ]] && backend_is_dind=1

# Prefer MagicDNS when Tailscale is up on this backend; else CoreDNS-first on DinD cold start
for name in "${NODES[@]}"; do
  log "configuring resolvers on $name"
  ts_up=0
  if vm_exec "$name" "tailscale ip -4 >/dev/null 2>&1"; then
    ts_up=1
  fi
  vm_exec "$name" "bash -s" <<EOS
set -euo pipefail
MAGIC=100.100.100.100
TS_UP=${ts_up}
BACKEND_IS_DIND=${backend_is_dind}

if [[ "\$TS_UP" == "1" ]]; then
  if systemctl is-active systemd-resolved >/dev/null 2>&1; then
    sudo mkdir -p /etc/systemd/resolved.conf.d
    sudo tee /etc/systemd/resolved.conf.d/failover-ci.conf >/dev/null <<EOF
[Resolve]
DNS=\${MAGIC} ${CD1} ${CD2} ${GOOGLE_DNS_1} ${GOOGLE_DNS_2}
FallbackDNS=${GOOGLE_DNS_1} ${GOOGLE_DNS_2}
Domains=~${DOMAIN}
DNSSEC=no
EOF
    sudo systemctl restart systemd-resolved || true
  fi
  sudo tee /etc/resolv.conf >/dev/null <<EOF
# failover-ci: MagicDNS > CoreDNS > Google
nameserver 100.100.100.100
nameserver ${CD1}
nameserver ${CD2}
nameserver ${GOOGLE_DNS_1}
nameserver ${GOOGLE_DNS_2}
search ${DOMAIN}
options timeout:2 attempts:2
EOF
elif [[ "\$BACKEND_IS_DIND" == "1" ]]; then
  sudo tee /etc/resolv.conf >/dev/null <<EOF
# failover-ci DinD (no Tailscale yet): CoreDNS > Google
nameserver ${CD1}
nameserver ${CD2}
nameserver ${GOOGLE_DNS_1}
nameserver ${GOOGLE_DNS_2}
search ${DOMAIN}
options timeout:2 attempts:2
EOF
else
  if systemctl is-active systemd-resolved >/dev/null 2>&1; then
    sudo mkdir -p /etc/systemd/resolved.conf.d
    sudo tee /etc/systemd/resolved.conf.d/failover-ci.conf >/dev/null <<EOF
[Resolve]
DNS=\${MAGIC} ${CD1} ${CD2} ${GOOGLE_DNS_1} ${GOOGLE_DNS_2}
FallbackDNS=${GOOGLE_DNS_1} ${GOOGLE_DNS_2}
Domains=~${DOMAIN}
DNSSEC=no
EOF
    sudo systemctl restart systemd-resolved || true
  fi
  sudo tee /etc/resolv.conf >/dev/null <<EOF
# failover-ci dual DNS: MagicDNS > CoreDNS > Google
nameserver 100.100.100.100
nameserver ${CD1}
nameserver ${CD2}
nameserver ${GOOGLE_DNS_1}
nameserver ${GOOGLE_DNS_2}
search ${DOMAIN}
options timeout:2 attempts:2
EOF
fi

mkdir -p /opt/failover-ci
cat > /opt/failover-ci/resolvers.env <<EOF
DOMAIN=${DOMAIN}
COREDNS_1=${CD1}
COREDNS_2=${CD2}
GOOGLE_DNS_1=${GOOGLE_DNS_1}
GOOGLE_DNS_2=${GOOGLE_DNS_2}
MAGICDNS=100.100.100.100
EOF
EOS
done

log "configure-resolvers complete"
