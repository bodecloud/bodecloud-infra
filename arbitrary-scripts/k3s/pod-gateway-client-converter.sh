#!/bin/bash

# Script to convert services to use pod-gateway for transparent WARP VPN routing
# This replaces the simple route manipulation with proper VXLAN tunneling

set -e

NAMESPACE="my-media-stack"
GATEWAY_NAME="warp-pod-gateway"

# Services that should route through VPN (same as Docker Compose network_mode: service:warp)
VPN_SERVICES=(
    "aiostreams"
    "mediafusion"
    "comet"
    "torrentio"
    "stremthru"
    "jackett"
    "prowlarr"
    "flaresolverr"
)

echo "🚀 Converting services to use pod-gateway for transparent WARP VPN routing..."
echo "   This provides VXLAN tunneling, DHCP, DNS, and automatic reconnection"
echo ""

for service in "${VPN_SERVICES[@]}"; do
    echo "🔍 Processing $service..."
    
    # Check if deployment exists
    if ! k3s kubectl get deployment "$service" -n "$NAMESPACE" &>/dev/null; then
        echo "   ⚠️  Deployment $service not found, skipping..."
        continue
    fi

    # Check if already has pod-gateway containers
    if k3s kubectl get deployment "$service" -n "$NAMESPACE" -o yaml | grep -q "pod-gateway.*client"; then
        echo "   ✅ $service already configured for pod-gateway VPN routing"
        continue
    fi

    # Remove old simple VPN routing if it exists
    if k3s kubectl get deployment "$service" -n "$NAMESPACE" -o yaml | grep -q "vpn-route-init"; then
        echo "   🧹 Removing old simple VPN routing from $service..."
        k3s kubectl patch deployment "$service" -n "$NAMESPACE" --type='json' -p='[{"op": "remove", "path": "/spec/template/spec/initContainers"}]' || true
        k3s kubectl label deployment "$service" -n "$NAMESPACE" vpn.routing- || true
    fi

    echo "   🔧 Adding pod-gateway VPN routing to $service..."

    # Create pod-gateway client patch
    POD_GATEWAY_PATCH='{
        "spec": {
            "template": {
                "metadata": {
                    "labels": {
                        "vpn.routing": "pod-gateway"
                    }
                },
                "spec": {
                    "initContainers": [
                        {
                            "name": "pod-gateway-client-init",
                            "image": "ghcr.io/angelnu/pod-gateway:v1.10.0",
                            "command": ["/bin/client_init.sh"],
                            "securityContext": {
                                "privileged": true,
                                "capabilities": {
                                    "add": ["NET_ADMIN", "SYS_MODULE"]
                                }
                            },
                            "env": [
                                {
                                    "name": "gateway",
                                    "value": "'$GATEWAY_NAME'"
                                },
                                {
                                    "name": "K8S_DNS_ips",
                                    "value": "10.43.0.10"
                                },
                                {
                                    "name": "VXLAN_ID",
                                    "value": "42"
                                },
                                {
                                    "name": "VXLAN_PORT",
                                    "value": "4789"
                                },
                                {
                                    "name": "VXLAN_IP_NETWORK",
                                    "value": "172.16.0"
                                },
                                {
                                    "name": "VPN_INTERFACE_MTU",
                                    "value": "1420"
                                },
                                {
                                    "name": "CONNECTION_RETRY_COUNT",
                                    "value": "3"
                                }
                            ],
                            "volumeMounts": [
                                {
                                    "name": "pod-gateway-client-config",
                                    "mountPath": "/config"
                                }
                            ]
                        }
                    ],
                    "containers": [
                        {
                            "name": "pod-gateway-client-sidecar",
                            "image": "ghcr.io/angelnu/pod-gateway:v1.10.0",
                            "command": ["/bin/client_sidecar.sh"],
                            "securityContext": {
                                "privileged": true,
                                "capabilities": {
                                    "add": ["NET_ADMIN"]
                                }
                            },
                            "env": [
                                {
                                    "name": "gateway",
                                    "value": "'$GATEWAY_NAME'"
                                },
                                {
                                    "name": "K8S_DNS_ips",
                                    "value": "10.43.0.10"
                                },
                                {
                                    "name": "VXLAN_ID",
                                    "value": "42"
                                },
                                {
                                    "name": "VXLAN_PORT",
                                    "value": "4789"
                                },
                                {
                                    "name": "VXLAN_IP_NETWORK",
                                    "value": "172.16.0"
                                },
                                {
                                    "name": "CONNECTION_RETRY_COUNT",
                                    "value": "3"
                                }
                            ],
                            "volumeMounts": [
                                {
                                    "name": "pod-gateway-client-config",
                                    "mountPath": "/config"
                                }
                            ],
                            "resources": {
                                "requests": {
                                    "memory": "32Mi",
                                    "cpu": "25m"
                                },
                                "limits": {
                                    "memory": "64Mi",
                                    "cpu": "50m"
                                }
                            }
                        }
                    ],
                    "volumes": [
                        {
                            "name": "pod-gateway-client-config",
                            "configMap": {
                                "name": "pod-gateway-client-config",
                                "defaultMode": 493
                            }
                        }
                    ]
                }
            }
        }
    }'

    # Apply the patch
    if k3s kubectl patch deployment "$service" -n "$NAMESPACE" --type='strategic' -p="$POD_GATEWAY_PATCH"; then
        echo "   ✅ $service configured for pod-gateway VPN routing"
        
        # Add pod-gateway routing label to deployment
        k3s kubectl label deployment "$service" -n "$NAMESPACE" vpn.routing=pod-gateway --overwrite
    else
        echo "   ❌ Failed to configure $service"
    fi
done

echo ""
echo "🎯 Pod-Gateway VPN Configuration Complete!"
echo ""
echo "📋 Summary:"
echo "   - All specified services now route through WARP VPN via pod-gateway"
echo "   - Services use VXLAN tunneling for robust connectivity"
echo "   - Automatic DHCP IP assignment and DNS resolution"
echo "   - Automatic reconnection if gateway restarts"
echo "   - Services are labeled with 'vpn.routing=pod-gateway'"
echo ""
echo "🚀 Next steps:"
echo "   1. Deploy pod-gateway client config: kubectl apply -f pod-gateway-client-config.yaml"
echo "   2. Deploy WARP pod-gateway: kubectl apply -f warp-pod-gateway-deployment.yaml"
echo "   3. Wait for gateway to be ready: kubectl wait --for=condition=ready pod -l app=warp-pod-gateway -n $NAMESPACE --timeout=300s"
echo "   4. Test VPN routing: kubectl exec -it deployment/jackett -n $NAMESPACE -- wget -qO- http://ifconfig.me"
echo ""
echo "🔍 Verify pod-gateway routing:"
echo "   kubectl exec -it deployment/jackett -n $NAMESPACE -- ip addr show vxlan0"
echo "   kubectl exec -it deployment/jackett -n $NAMESPACE -- ip route show"
echo "   kubectl exec -it deployment/jackett -n $NAMESPACE -- wget -qO- http://ifconfig.me"
echo ""
echo "🐛 Debug commands:"
echo "   kubectl logs -l app=warp-pod-gateway -n $NAMESPACE -c gateway-sidecar"
echo "   kubectl logs deployment/jackett -n $NAMESPACE -c pod-gateway-client-sidecar" 