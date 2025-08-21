#!/bin/bash

set -euo pipefail

# Configuration
OPENVPN_CONFIGS_DIR="../../configs/vpn-gateway"
NAMESPACE="default"
CONFIGMAP_NAME="resolved-vpn-configs"

echo "🔍 Starting VPN hostname resolution and ConfigMap generation..."

# Check if configs directory exists
if [ ! -d "$OPENVPN_CONFIGS_DIR" ]; then
    echo "❌ Error: OpenVPN configs directory not found: $OPENVPN_CONFIGS_DIR"
    exit 1
fi

# Verify kubectl is available
if ! command -v kubectl &>/dev/null; then
    echo "❌ Error: kubectl is not available in PATH"
    exit 1
fi

# Create temporary directory for resolved configs
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo "📁 Using temporary directory: $TEMP_DIR"

# Initialize ConfigMap YAML
cat >"$TEMP_DIR/configmap.yaml" <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: $CONFIGMAP_NAME
  namespace: $NAMESPACE
data:
EOF

# Process each .ovpn file
echo "🌐 Resolving VPN hostnames and generating configs..."

for OVPN_FILE in "$OPENVPN_CONFIGS_DIR"/*.ovpn; do
    if [ ! -f "$OVPN_FILE" ]; then
        echo "⚠️  No .ovpn files found in $OPENVPN_CONFIGS_DIR"
        continue
    fi

    BASENAME=$(basename "$OVPN_FILE")
    echo "📄 Processing: $BASENAME"

    # Extract hostname from the remote line
    HOSTNAME=$(grep "^remote " "$OVPN_FILE" | awk '{print $2}')

    if [ -z "$HOSTNAME" ]; then
        echo "⚠️  Warning: No hostname found in $BASENAME, skipping"
        continue
    fi

    echo "🔍 Resolving hostname: $HOSTNAME"

    # Resolve hostname to IP with retries
    IP_ADDRESS=""
    for attempt in 1 2 3; do
        IP_ADDRESS=$(getent hosts "$HOSTNAME" | awk '{print $1}' | head -n1)
        if [ -n "$IP_ADDRESS" ] && [[ "$IP_ADDRESS" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "✅ Resolved $HOSTNAME to $IP_ADDRESS"
            break
        fi
        echo "⏳ Attempt $attempt failed, retrying..."
        sleep 2
    done

    if [ -z "$IP_ADDRESS" ] || [[ ! "$IP_ADDRESS" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "❌ Error: Failed to resolve $HOSTNAME to a valid IP address"
        exit 1
    fi

    # Create the resolved config
    RESOLVED_CONFIG=$(sed "s/remote $HOSTNAME 1194/remote $IP_ADDRESS 1194/" "$OVPN_FILE")

    # Add to ConfigMap YAML with proper indentation
    echo "  $BASENAME: |" >>"$TEMP_DIR/configmap.yaml"
    echo "$RESOLVED_CONFIG" | sed 's/^/    /' >>"$TEMP_DIR/configmap.yaml"
done

echo "📝 Generated ConfigMap:"
cat "$TEMP_DIR/configmap.yaml"

echo ""
echo "🚀 Applying ConfigMap to Kubernetes..."

# Apply the ConfigMap
if kubectl apply -f "$TEMP_DIR/configmap.yaml"; then
    echo "✅ Successfully applied ConfigMap: $CONFIGMAP_NAME"

    # Verify the ConfigMap was created
    echo "🔍 Verifying ConfigMap..."
    kubectl get configmap "$CONFIGMAP_NAME" -n "$NAMESPACE" -o yaml | head -20

    echo ""
    echo "✅ VPN hostname resolution and ConfigMap creation completed successfully!"
    echo "📋 ConfigMap '$CONFIGMAP_NAME' contains resolved VPN configs ready for use by Gluetun"
else
    echo "❌ Error: Failed to apply ConfigMap"
    exit 1
fi
