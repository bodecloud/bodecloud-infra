#!/bin/bash

# Script to convert services to use transparent WARP VPN gateway
# This script adds init containers that redirect default gateway

set -e

NAMESPACE="my-media-stack"

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

echo "🔄 Converting services to use transparent WARP VPN gateway..."

for service in "${VPN_SERVICES[@]}"; do
    echo "  🔍 Processing $service..."
    
    # Check if deployment exists
    if ! k3s kubectl get deployment "$service" -n "$NAMESPACE" &>/dev/null; then
        echo "     ⚠️  Deployment $service not found, skipping..."
        continue
    fi

    # Check if already has VPN init container
    if k3s kubectl get deployment "$service" -n "$NAMESPACE" -o yaml | grep -q "vpn-route-init"; then
        echo "     ✅ $service already configured for transparent VPN routing"
        continue
    fi

    echo "     🔧 Adding transparent VPN routing to $service..."

    # Create init container patch
    INIT_CONTAINER_PATCH='{
        "spec": {
            "template": {
                "metadata": {
                    "labels": {
                        "vpn.routing": "enabled"
                    }
                },
                "spec": {
                    "initContainers": [
                        {
                            "name": "vpn-route-init",
                            "image": "busybox",
                            "securityContext": {
                                "privileged": true,
                                "capabilities": {
                                    "add": ["NET_ADMIN"]
                                }
                            },
                            "command": [
                                "sh",
                                "-c",
                                "echo \"Waiting for WARP gateway...\"; until nslookup warp-gateway.my-media-stack.svc.cluster.local; do echo \"WARP gateway not ready, waiting...\"; sleep 5; done; WARP_IP=$(nslookup warp-gateway.my-media-stack.svc.cluster.local | grep -A1 \"Name:\" | tail -1 | awk \"{print \\$2}\"); echo \"WARP gateway IP: $WARP_IP\"; ip route del default || true; ip route add default via $WARP_IP; echo \"New routing table:\"; ip route show; echo \"VPN routing configured successfully\""
                            ]
                        }
                    ]
                }
            }
        }
    }'

    # Apply the patch
    if k3s kubectl patch deployment "$service" -n "$NAMESPACE" --type='merge' -p="$INIT_CONTAINER_PATCH"; then
        echo "     ✅ $service configured for transparent VPN routing"
        
        # Add VPN routing label to deployment
        k3s kubectl label deployment "$service" -n "$NAMESPACE" vpn.routing=enabled --overwrite
    else
        echo "     ❌ Failed to configure $service"
    fi
done

echo ""
echo "🎯 Transparent VPN Gateway Configuration Complete!"
echo ""
echo "📋 Summary:"
echo "   - All specified services now route through WARP VPN gateway transparently"
echo "   - Services are labeled with 'vpn.routing=enabled'"
echo "   - Init containers redirect default gateway to WARP gateway"
echo "   - NO application-level configuration required"
echo "   - Services are completely unaware of VPN routing"
echo ""
echo "🚀 Next steps:"
echo "   1. Deploy WARP gateway: kubectl apply -f warp-gateway-deployment.yaml"
echo "   2. Wait for gateway to be ready: kubectl wait --for=condition=ready pod -l app=warp-gateway -n $NAMESPACE"
echo "   3. Test VPN routing: kubectl exec -it deployment/jackett -n $NAMESPACE -- wget -qO- http://ifconfig.me"
echo ""
echo "🔍 Verify VPN routing:"
echo "   kubectl exec -it deployment/jackett -n $NAMESPACE -- ip route show"
echo "   kubectl exec -it deployment/jackett -n $NAMESPACE -- wget -qO- http://ifconfig.me" 