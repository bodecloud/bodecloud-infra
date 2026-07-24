#!/usr/bin/env bash
# Deploy CoreDNS on ci-node1 and ci-node2 with identical zone for $DOMAIN.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/node-ips.json" ]] || die "run provision-vms.sh first"

N1="$(node_ip_from_state ci-node1)"
N2="$(node_ip_from_state ci-node2)"
N3="$(node_ip_from_state ci-node3)"
N4="$(node_ip_from_state ci-node4)"

ZONE_DIR="${STATE_DIR}/coredns"
mkdir -p "${ZONE_DIR}"

# Zone: global names → multi A (any node Traefik); node-direct → that node.
# TTL low for CI churn.
cat > "${ZONE_DIR}/${DOMAIN}.db" <<EOF
\$ORIGIN ${DOMAIN}.
\$TTL 30
@       IN SOA ns1.${DOMAIN}. admin.${DOMAIN}. (
        $(date +%Y%m%d%H) ; serial
        60 ; refresh
        30 ; retry
        3600 ; expire
        30 ) ; minimum
@       IN NS ns1.${DOMAIN}.
ns1     IN A ${N1}
ns2     IN A ${N2}

; Wildcard global → all Traefik entry IPs (any-node + peer-forward)
@       IN A ${N1}
@       IN A ${N2}
@       IN A ${N3}
@       IN A ${N4}
*       IN A ${N1}
*       IN A ${N2}
*       IN A ${N3}
*       IN A ${N4}

; Node-direct wildcards
*.ci-node1 IN A ${N1}
*.ci-node2 IN A ${N2}
*.ci-node3 IN A ${N3}
*.ci-node4 IN A ${N4}

ci-node1 IN A ${N1}
ci-node2 IN A ${N2}
ci-node3 IN A ${N3}
ci-node4 IN A ${N4}
EOF

COREFILE="${ZONE_DIR}/Corefile"
sed \
  -e "s/__DOMAIN__/${DOMAIN}/g" \
  -e "s/__GOOGLE_DNS_1__/${GOOGLE_DNS_1}/g" \
  -e "s/__GOOGLE_DNS_2__/${GOOGLE_DNS_2}/g" \
  "${SCRIPT_DIR}/coredns/Corefile.tmpl" > "${COREFILE}"

for name in "${COREDNS_NODES[@]}"; do
  log "deploying CoreDNS on $name"
  if [[ "$(backend)" == "dind" ]]; then
    dind_seed_image_from_host "$name" "coredns/coredns:1.11.3" \
      || dind_seed_image_from_host "$name" "docker.io/coredns/coredns:1.11.3" \
      || die "CoreDNS image missing on host — run: docker pull coredns/coredns:1.11.3"
  fi
  vm_exec "$name" "sudo mkdir -p /opt/coredns/zones && sudo chown -R \$(whoami):\$(whoami) /opt/coredns"
  vm_transfer "$name" "${COREFILE}" "/opt/coredns/Corefile"
  vm_transfer "$name" "${ZONE_DIR}/${DOMAIN}.db" "/opt/coredns/zones/${DOMAIN}.db"
  vm_exec "$name" "docker rm -f coredns-ci >/dev/null 2>&1 || true
    docker run -d --name coredns-ci --restart=always \
      -p 53:53/udp -p 53:53/tcp -p 8181:8080 \
      -v /opt/coredns/Corefile:/Corefile:ro \
      -v /opt/coredns/zones:/etc/coredns/zones:ro \
      coredns/coredns:1.11.3 -conf /Corefile"
  vm_exec "$name" "docker inspect -f '{{.State.Running}}' coredns-ci | grep -qx true" \
    || die "CoreDNS failed to start on $name"
done

# Persist CoreDNS IPs for resolver config
printf '%s\n' "$N1" "$N2" > "${STATE_DIR}/coredns-ips.txt"
log "CoreDNS listening on ${N1} and ${N2} for ${DOMAIN}"
