#!/bin/bash
# Apply Kubernetes manifests in the correct order to avoid dependency issues
# This script ensures cert-manager is ready before applying dependent resources

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1"
    exit 1
}

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    error "kubectl is not installed or not in PATH"
fi

if ! kubectl cluster-info &> /dev/null; then
    error "Cannot connect to Kubernetes cluster"
fi

log "🚀 Starting ordered deployment of my-media-stack"

# Step 1: Apply namespaces first
log "📁 Step 1: Creating namespaces..."
kubectl apply -f 00-namespaces.yml
log "✅ Namespaces created"

# Step 2: Apply cert-manager installation
log "🔧 Step 2: Installing cert-manager..."
kubectl apply -f 00-cert-manager-install.yml

# Wait for cert-manager to be ready
log "⏳ Waiting for cert-manager to be ready..."
kubectl wait --for=condition=available deployment/cert-manager -n cert-manager --timeout=300s
log "✅ cert-manager is ready"

# Step 3: Apply cert-manager setup (ClusterIssuer and Certificate)
log "🔐 Step 3: Setting up cert-manager resources..."
kubectl apply -f 02-cert-manager-setup.yml

# Wait for the certificate to be ready
log "⏳ Waiting for certificate to be ready..."
kubectl wait --for=condition=ready certificate/proxy-injector-cert -n vpn-gateway --timeout=300s
log "✅ cert-manager setup completed"

# Step 4: Apply proxy injector
log "🔌 Step 4: Deploying cluster proxy injector..."
kubectl apply -f 03-cluster-proxy-injector.yml

# Wait for proxy injector to be ready
log "⏳ Waiting for proxy injector to be ready..."
kubectl wait --for=condition=available deployment/proxy-injector -n vpn-gateway --timeout=300s
log "✅ Proxy injector is ready"

# Step 5: Apply all other manifests
log "📦 Step 5: Applying remaining manifests..."

# Get all YAML files except the ones we've already applied
remaining_files=$(find . -name "*.yml" -o -name "*.yaml" | grep -v -E "(00-namespaces|00-cert-manager|02-cert-manager-setup|03-cluster-proxy-injector)" | sort)

if [ -n "$remaining_files" ]; then
    for file in $remaining_files; do
        log "📄 Applying $file..."
        kubectl apply -f "$file"
    done
    log "✅ All remaining manifests applied"
else
    log "ℹ️  No additional manifests to apply"
fi

log "🎉 Deployment completed successfully!"
echo
log "📋 Summary:"
log "  • Namespaces: Created"
log "  • cert-manager: Installed and configured"
log "  • Proxy injector: Deployed and ready"
log "  • All other resources: Applied"
echo
log "🔍 Useful commands:"
log "  • Check all pods: kubectl get pods -A"
log "  • Check cert-manager: kubectl get pods -n cert-manager"
log "  • Check proxy injector: kubectl get pods -n vpn-gateway"
log "  • Check certificates: kubectl get certificates -A"
echo
log "✅ Your media stack is ready!" 