#!/bin/bash

set -euo pipefail

echo "🚀 VPN Gateway Kubernetes Deployment Script"
echo "============================================"

# Configuration
NAMESPACE="default"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIGS_DIR="$SCRIPT_DIR/../../configs/vpn-gateway"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}ℹ️  $1${NC}"
}

log_warn() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed or not in PATH"
        exit 1
    fi

    # Check cluster connectivity
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi

    # Check configs directory
    if [ ! -d "$CONFIGS_DIR" ]; then
        log_error "VPN configs directory not found: $CONFIGS_DIR"
        exit 1
    fi

    log_info "Prerequisites check passed ✅"
}

# Step 1: Resolve VPN hostnames and create ConfigMaps
resolve_vpn_ips() {
    log_info "Step 1: Resolving VPN hostnames and creating ConfigMaps..."

    if [ -f "$SCRIPT_DIR/resolve-and-apply-vpn-configs.sh" ]; then
        log_info "Running VPN hostname resolution script..."
        "$SCRIPT_DIR/resolve-and-apply-vpn-configs.sh"
    else
        log_error "VPN resolution script not found: $SCRIPT_DIR/resolve-and-apply-vpn-configs.sh"
        exit 1
    fi
}

# Step 2: Create CA certificate ConfigMap
create_ca_configmap() {
    log_info "Step 2: Creating CA certificate ConfigMap..."

    local ca_file="$CONFIGS_DIR/premiumize-ca.pem"
    if [ ! -f "$ca_file" ]; then
        log_error "CA certificate file not found: $ca_file"
        exit 1
    fi

    # Check if it's still a placeholder
    if grep -q "Replace with your Premiumize CA certificate" "$ca_file"; then
        log_warn "CA certificate appears to be a placeholder!"
        log_warn "Please replace the content of $ca_file with your actual Premiumize CA certificate"
        log_warn "You can download it from your Premiumize account"
        echo
        read -p "Do you want to continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Deployment cancelled. Please update the CA certificate first."
            exit 1
        fi
    fi

    kubectl create configmap vpn-ca-cert \
        --from-file=premiumize-ca.pem="$ca_file" \
        --dry-run=client -o yaml | kubectl apply -f -

    log_info "CA certificate ConfigMap created ✅"
}

# Step 3: Create auth credentials Secret
create_auth_secret() {
    log_info "Step 3: Creating auth credentials Secret..."

    local auth_file="$CONFIGS_DIR/auth.conf"
    if [ ! -f "$auth_file" ]; then
        log_error "Auth file not found: $auth_file"
        exit 1
    fi

    # Check if it's still a placeholder
    if grep -q "your_premiumize_username" "$auth_file"; then
        log_warn "Auth file appears to contain placeholder credentials!"
        log_warn "Please replace the content of $auth_file with your actual Premiumize username and password"
        echo
        read -p "Do you want to continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Deployment cancelled. Please update the auth credentials first."
            exit 1
        fi
    fi

    kubectl create secret generic vpn-auth \
        --from-file=auth.conf="$auth_file" \
        --dry-run=client -o yaml | kubectl apply -f -

    log_info "Auth credentials Secret created ✅"
}

# Step 4: Deploy Gluetun VPN gateways
deploy_gluetun() {
    log_info "Step 4: Deploying Gluetun VPN gateways..."

    # Check if deployment files exist
    if [ -f "$SCRIPT_DIR/gluetun-test.yaml" ]; then
        log_info "Deploying test Gluetun instance (Finland)..."
        kubectl apply -f "$SCRIPT_DIR/gluetun-test.yaml"
    fi

    if [ -f "$SCRIPT_DIR/gluetun-vpn-gateway.yaml" ]; then
        log_info "Deploying full Gluetun VPN gateway instances..."
        kubectl apply -f "$SCRIPT_DIR/gluetun-vpn-gateway.yaml"
    fi

    log_info "Gluetun deployments created ✅"
}

# Step 5: Verify deployment
verify_deployment() {
    log_info "Step 5: Verifying deployment..."

    echo
    log_info "ConfigMaps:"
    kubectl get configmaps | grep -E "(resolved-vpn-configs|vpn-ca-cert|vpn-templates)" || true

    echo
    log_info "Secrets:"
    kubectl get secrets | grep "vpn-auth" || true

    echo
    log_info "Deployments:"
    kubectl get deployments | grep "gluetun" || true

    echo
    log_info "Pods:"
    kubectl get pods | grep "gluetun" || true

    echo
    log_info "Services:"
    kubectl get services | grep "gluetun" || true
}

# Step 6: Show connection status
show_connection_status() {
    log_info "Step 6: Checking VPN connection status..."

    # Get all gluetun pods
    local pods=$(kubectl get pods -l app=gluetun-test -o jsonpath='{.items[*].metadata.name}' 2>/dev/null || echo "")

    if [ -z "$pods" ]; then
        log_warn "No Gluetun pods found"
        return
    fi

    for pod in $pods; do
        echo
        log_info "Checking pod: $pod"

        # Check if pod is running
        local status=$(kubectl get pod "$pod" -o jsonpath='{.status.phase}')
        if [ "$status" != "Running" ]; then
            log_warn "Pod $pod is not running (status: $status)"
            continue
        fi

        # Show recent logs
        log_info "Recent logs for $pod:"
        kubectl logs "$pod" --tail=5 | sed 's/^/  /'
    done
}

# Cleanup function
cleanup() {
    log_warn "Cleaning up existing resources..."

    # Delete existing deployments
    kubectl delete deployment -l app=gluetun-test --ignore-not-found=true
    kubectl delete deployment -l app=gluetun-vpn-gateway --ignore-not-found=true

    # Delete services
    kubectl delete service -l app=gluetun-test --ignore-not-found=true
    kubectl delete service -l app=gluetun-vpn-gateway --ignore-not-found=true

    # Delete ConfigMaps and Secrets
    kubectl delete configmap resolved-vpn-configs --ignore-not-found=true
    kubectl delete configmap vpn-ca-cert --ignore-not-found=true
    kubectl delete configmap vpn-templates --ignore-not-found=true
    kubectl delete secret vpn-auth --ignore-not-found=true

    # Delete jobs
    kubectl delete job vpn-ip-resolver --ignore-not-found=true

    log_info "Cleanup completed ✅"
}

# Main execution
main() {
    echo
    log_info "Starting VPN Gateway deployment..."
    echo

    # Parse command line arguments
    case "${1:-deploy}" in
        "deploy")
            check_prerequisites
            resolve_vpn_ips
            create_ca_configmap
            create_auth_secret
            deploy_gluetun
            verify_deployment
            show_connection_status

            echo
            log_info "🎉 VPN Gateway deployment completed!"
            log_info "Use 'kubectl logs -l app=gluetun-test' to monitor VPN connections"
            log_info "Use 'kubectl port-forward svc/gluetun-test-fi 8080:8080' to access HTTP proxy"
            ;;
        "cleanup")
            cleanup
            ;;
        "status")
            verify_deployment
            show_connection_status
            ;;
        "resolve")
            resolve_vpn_ips
            ;;
        *)
            echo "Usage: $0 [deploy|cleanup|status|resolve]"
            echo
            echo "Commands:"
            echo "  deploy   - Full deployment (default)"
            echo "  cleanup  - Remove all VPN gateway resources"
            echo "  status   - Show current deployment status"
            echo "  resolve  - Only resolve VPN IPs and update ConfigMaps"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"