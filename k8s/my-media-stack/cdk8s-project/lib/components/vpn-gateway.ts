import { Construct } from 'constructs';
import { ApiObject } from 'cdk8s';
import { KubeConfigMap } from '../../imports/k8s';

/**
 * Creates VPN Gateway components including pod-gateway configuration
 */
export function createVpnGateway(scope: Construct): void {
  // Pod Gateway ConfigMap
  new KubeConfigMap(scope, 'pod-gateway-config', {
    metadata: {
      name: 'pod-gateway-config',
      namespace: 'vpn-gateway',
    },
    data: {
      'settings.sh': `#!/bin/bash
# Pod-Gateway Settings for VPN Integration

# Gateway hostname (will be set dynamically by controller)
GATEWAY_NAME="\${ACTIVE_VPN_GATEWAY:-gluetun-premiumize-us}"

# K8S DNS configuration
K8S_DNS_IPS="10.96.0.10"
NOT_ROUTED_TO_GATEWAY_CIDRS="10.0.0.0/8 172.16.0.0/12 192.168.0.0/16"

# VXLAN Configuration
VXLAN_ID="42"
VXLAN_PORT="4789"
VXLAN_IP_NETWORK="172.16.1"
VXLAN_GATEWAY_FIRST_DYNAMIC_IP=20

# VPN Interface Configuration
VPN_INTERFACE="tun0"
VPN_BLOCK_OTHER_TRAFFIC=true
VPN_TRAFFIC_PORT="1194"
VPN_LOCAL_CIDRS="10.0.0.0/8 192.168.0.0/16 172.16.0.0/12"

# DNS Configuration
DNS_LOCAL_CIDRS="cluster.local svc.cluster.local"
RESOLV_CONF_COPY="/etc/resolv_copy.conf"

# Health and Connection Settings
CONNECTION_RETRY_COUNT=3
GATEWAY_ENABLE_DNSSEC=true
IPTABLES_NFT=no

# VPN Interface MTU
VPN_INTERFACE_MTU="1420"`,

      'gateway_init.sh': `#!/bin/bash
set -euo pipefail

# Load settings
source /config/settings.sh

echo "=== Initializing VPN Gateway Pod-Gateway ==="

# Enable IP forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward
echo 0 > /proc/sys/net/ipv6/conf/all/disable_ipv6

# Create VXLAN interface
VXLAN_GATEWAY_IP="\${VXLAN_IP_NETWORK}.1"

if ip link show vxlan0 >/dev/null 2>&1; then
    ip link del vxlan0
fi

ip link add vxlan0 type vxlan id $VXLAN_ID dev eth0 dstport $VXLAN_PORT
ip addr add \${VXLAN_GATEWAY_IP}/24 dev vxlan0
ip link set up dev vxlan0
ip link set mtu $VPN_INTERFACE_MTU dev vxlan0

# Configure iptables for NAT and forwarding
iptables -t nat -A POSTROUTING -j MASQUERADE

echo "✅ VPN Gateway Pod-Gateway initialization complete"`,

      'client_init.sh': `#!/bin/bash
set -euo pipefail

source /config/settings.sh

echo "=== Initializing Pod-Gateway Client ==="

# Set up routes for local traffic
K8S_GW_IP=$(ip route | awk '/default/ { print $3 }')
for local_cidr in $NOT_ROUTED_TO_GATEWAY_CIDRS; do
    ip route add $local_cidr via $K8S_GW_IP || true
done

echo "✅ Pod-Gateway Client initialization complete"`,
    },
  });

  // Pod Gateway HelmChartConfig for K3s
  new ApiObject(scope, 'pod-gateway-helmchartconfig', {
    apiVersion: 'helm.cattle.io/v1',
    kind: 'HelmChartConfig',
    metadata: {
      name: 'pod-gateway',
      namespace: 'kube-system',
    },
    spec: {
      valuesContent: `image:
  repository: angelnu/pod-gateway
  tag: v1.10.0

settings:
  VPN_INTERFACE: tun0
  VPN_BLOCK_OTHER_TRAFFIC: true
  VPN_TRAFFIC_PORT: "1194"
  VPN_LOCAL_CIDRS: "10.0.0.0/8 192.168.0.0/16 172.16.0.0/12"
  NOT_ROUTED_TO_GATEWAY_CIDRS: "10.0.0.0/8 172.16.0.0/12 192.168.0.0/16"
  
  VXLAN_ID: 42
  VXLAN_IP_NETWORK: "172.16.1"
  VXLAN_GATEWAY_FIRST_DYNAMIC_IP: 20

routed_namespaces:
  - my-media-stack

webhook:
  image:
    repository: angelnu/gateway-admision-controller
    tag: v3.10.0

publicPorts:
  - hostname: gluetun-premiumize-us
    IP: 10
    ports:
      - type: tcp
        port: 8080
      - type: udp  
        port: 6881

DNS: 172.16.1.1`,
    },
  });
} 