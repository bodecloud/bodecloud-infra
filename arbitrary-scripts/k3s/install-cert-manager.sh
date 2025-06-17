#!/bin/bash
# Comprehensive cert-manager installation script
# This script installs cert-manager and then creates the required issuers and certificates

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
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        error "kubectl is not installed or not in PATH"
    fi
    
    if ! kubectl cluster-info &> /dev/null; then
        error "Cannot connect to Kubernetes cluster"
    fi
    
    log "✅ kubectl is working and connected to cluster"
}

# Install cert-manager using the official method
install_cert_manager() {
    log "🔧 Installing cert-manager..."
    
    # Check if cert-manager is already installed
    if kubectl get namespace cert-manager &> /dev/null; then
        warn "cert-manager namespace already exists, checking if it's working..."
        if kubectl get deployment cert-manager -n cert-manager &> /dev/null; then
            log "✅ cert-manager is already installed"
            return 0
        fi
    fi
    
    # Install cert-manager using the official YAML
    log "📦 Applying cert-manager CRDs and deployment..."
    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.3/cert-manager.yaml
    
    # Wait for cert-manager to be ready
    log "⏳ Waiting for cert-manager to be ready..."
    kubectl wait --for=condition=available deployment/cert-manager -n cert-manager --timeout=300s
    kubectl wait --for=condition=available deployment/cert-manager-cainjector -n cert-manager --timeout=300s
    kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=300s
    
    log "✅ cert-manager installed and ready"
}

# Create the vpn-gateway namespace if it doesn't exist
create_vpn_namespace() {
    log "🔧 Creating vpn-gateway namespace..."
    kubectl create namespace vpn-gateway --dry-run=client -o yaml | kubectl apply -f -
    log "✅ vpn-gateway namespace ready"
}

# Create self-signed issuer
create_self_signed_issuer() {
    log "🔧 Creating self-signed ClusterIssuer..."
    
    cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
EOF
    
    log "✅ Self-signed ClusterIssuer created"
}

# Create certificate for proxy injector
create_proxy_injector_cert() {
    log "🔧 Creating proxy-injector certificate..."
    
    cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: proxy-injector-cert
  namespace: vpn-gateway
spec:
  secretName: proxy-injector-certs
  issuerRef:
    name: selfsigned-issuer
    kind: ClusterIssuer
  dnsNames:
  - proxy-injector.vpn-gateway.svc
  - proxy-injector.vpn-gateway.svc.cluster.local
EOF
    
    # Wait for certificate to be ready
    log "⏳ Waiting for certificate to be issued..."
    kubectl wait --for=condition=ready certificate/proxy-injector-cert -n vpn-gateway --timeout=120s
    
    log "✅ Proxy-injector certificate created and ready"
}

# Verify installation
verify_installation() {
    log "🔍 Verifying cert-manager installation..."
    
    # Check cert-manager pods
    if ! kubectl get pods -n cert-manager | grep -q "Running"; then
        error "cert-manager pods are not running"
    fi
    
    # Check ClusterIssuer
    if ! kubectl get clusterissuer selfsigned-issuer &> /dev/null; then
        error "selfsigned-issuer ClusterIssuer not found"
    fi
    
    # Check certificate
    if ! kubectl get certificate proxy-injector-cert -n vpn-gateway &> /dev/null; then
        error "proxy-injector-cert Certificate not found"
    fi
    
    # Check secret
    if ! kubectl get secret proxy-injector-certs -n vpn-gateway &> /dev/null; then
        error "proxy-injector-certs Secret not found"
    fi
    
    log "✅ All cert-manager components are working correctly"
}

# Main function
main() {
    log "🚀 Starting comprehensive cert-manager installation"
    
    check_kubectl
    install_cert_manager
    create_vpn_namespace
    create_self_signed_issuer
    create_proxy_injector_cert
    verify_installation
    
    log "🎉 cert-manager installation complete!"
    echo
    log "📋 Summary:"
    log "  • cert-manager: Installed and running"
    log "  • selfsigned-issuer: ClusterIssuer created"
    log "  • proxy-injector-cert: Certificate issued"
    log "  • proxy-injector-certs: Secret created"
    echo
    log "🔍 Useful commands:"
    log "  • Check cert-manager status: kubectl get pods -n cert-manager"
    log "  • Check certificates: kubectl get certificates -A"
    log "  • Check issuers: kubectl get clusterissuers"
    echo
    log "✅ You can now apply your other Kubernetes manifests that depend on cert-manager"
}

# Run main function
main "$@" 