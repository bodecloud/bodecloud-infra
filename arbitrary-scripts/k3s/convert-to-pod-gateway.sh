#!/bin/bash

# Ultra-simple pod-gateway converter - just add label "vpn=warp" to any deployment
set -e

NAMESPACE="my-media-stack"

echo "🚀 Auto-converting deployments with label 'vpn=warp' to use pod-gateway..."

# Find all deployments with vpn=warp label
DEPLOYMENTS=$(k3s kubectl get deployments -n "$NAMESPACE" -l vpn=warp -o name 2>/dev/null | sed 's|deployment.apps/||' || true)

if [ -z "$DEPLOYMENTS" ]; then
    echo "❌ No deployments found with label 'vpn=warp'"
    echo ""
    echo "💡 To use pod-gateway VPN routing, label your deployments:"
    echo "   kubectl label deployment <name> vpn=warp -n $NAMESPACE"
    echo ""
    echo "📋 Available deployments:"
    k3s kubectl get deployments -n "$NAMESPACE" -o name | sed 's|deployment.apps/|   - |'
    exit 0
fi

echo "📋 Found deployments with vpn=warp label:"
for deployment in $DEPLOYMENTS; do
    echo "   - $deployment"
done
echo ""

for deployment in $DEPLOYMENTS; do
    echo "🔧 Converting $deployment..."
    
    # Check if already converted
    if k3s kubectl get deployment "$deployment" -n "$NAMESPACE" -o yaml | grep -q "pod-gateway-client"; then
        echo "   ✅ Already converted"
        continue
    fi

    # Simple patch to add pod-gateway containers
    k3s kubectl patch deployment "$deployment" -n "$NAMESPACE" --type='merge' -p='{
        "spec": {
            "template": {
                "spec": {
                    "initContainers": [
                        {
                            "name": "pod-gateway-client",
                            "image": "ghcr.io/angelnu/pod-gateway:v1.10.0",
                            "command": ["/bin/client_init.sh"],
                            "securityContext": {"privileged": true},
                            "env": [
                                {"name": "gateway", "value": "warp-pod-gateway.my-media-stack.svc.cluster.local"},
                                {"name": "VXLAN_ID", "value": "42"},
                                {"name": "VXLAN_IP_NETWORK", "value": "172.16.0"}
                            ],
                            "volumeMounts": [{"name": "pod-gateway-config", "mountPath": "/config"}]
                        }
                    ],
                    "containers": [
                        {
                            "name": "pod-gateway-sidecar",
                            "image": "ghcr.io/angelnu/pod-gateway:v1.10.0",
                            "command": ["/bin/client_sidecar.sh"],
                            "securityContext": {"privileged": true},
                            "env": [
                                {"name": "gateway", "value": "warp-pod-gateway.my-media-stack.svc.cluster.local"},
                                {"name": "VXLAN_ID", "value": "42"},
                                {"name": "VXLAN_IP_NETWORK", "value": "172.16.0"}
                            ],
                            "volumeMounts": [{"name": "pod-gateway-config", "mountPath": "/config"}],
                            "resources": {"requests": {"memory": "16Mi", "cpu": "10m"}, "limits": {"memory": "32Mi", "cpu": "25m"}}
                        }
                    ],
                    "volumes": [
                        {
                            "name": "pod-gateway-config",
                            "configMap": {"name": "pod-gateway-config", "defaultMode": 493}
                        }
                    ]
                }
            }
        }
    }'
    
    echo "   ✅ Converted $deployment"
done

echo ""
echo "🎯 Done! All labeled deployments now route through WARP VPN"
echo ""
echo "💡 Usage:"
echo "   # Label any deployment for VPN routing:"
echo "   kubectl label deployment jackett vpn=warp -n $NAMESPACE"
echo ""
echo "   # Remove VPN routing:"
echo "   kubectl label deployment jackett vpn- -n $NAMESPACE"
echo "   kubectl rollout restart deployment/jackett -n $NAMESPACE" 