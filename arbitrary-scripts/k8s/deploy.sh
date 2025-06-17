#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "This script should not be run as root"
        exit 1
    fi
}

# Install k3s
install_k3s() {
    log_info "Installing k3s..."
    
    if command -v k3s &> /dev/null; then
        log_warning "k3s is already installed"
        return 0
    fi
    
    # Install k3s with specific configuration
    curl -sfL https://get.k3s.io | sh -s - \
        --write-kubeconfig-mode 644 \
        --disable traefik \
        --disable servicelb \
        --disable local-storage \
        --cluster-init
    
    # Wait for k3s to be ready
    log_info "Waiting for k3s to be ready..."
    sudo systemctl enable k3s
    sudo systemctl start k3s
    
    # Wait for node to be ready
    timeout=300
    while ! kubectl get nodes | grep -q "Ready"; do
        if [ $timeout -le 0 ]; then
            log_error "Timeout waiting for k3s to be ready"
            exit 1
        fi
        sleep 5
        timeout=$((timeout - 5))
    done
    
    log_success "k3s installed and ready"
}

# Setup kubectl access
setup_kubectl() {
    log_info "Setting up kubectl access..."
    
    # Copy kubeconfig to user directory
    mkdir -p ~/.kube
    sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
    sudo chown $(id -u):$(id -g) ~/.kube/config
    
    # Verify kubectl works
    if kubectl get nodes &> /dev/null; then
        log_success "kubectl configured successfully"
    else
        log_error "Failed to configure kubectl"
        exit 1
    fi
}

# Install required tools
install_tools() {
    log_info "Installing required tools..."
    
    # Update package list
    sudo apt-get update
    
    # Install required packages
    sudo apt-get install -y \
        curl \
        wget \
        unzip \
        jq \
        openssl \
        uuid-runtime
    
    log_success "Required tools installed"
}

# Deploy Kubernetes manifests
deploy_manifests() {
    log_info "Deploying Kubernetes manifests..."
    
    # Deploy manifests in order
    local manifests=(
        "k8s/namespace.yaml"
        "k8s/secrets.yaml"
        "k8s/missing-secrets.yaml"
        "k8s/storage.yaml"
        "k8s/infrastructure.yaml"
        "k8s/vpn-sidecar.yaml"
        "k8s/warp-stremio-addons.yaml"
        "k8s/media-services.yaml"
        "k8s/monitoring.yaml"
        "k8s/traefik.yaml"
    )
    
    for manifest in "${manifests[@]}"; do
        if [[ -f "$manifest" ]]; then
            log_info "Applying $manifest..."
            kubectl apply -f "$manifest"
            sleep 5
        else
            log_warning "Manifest $manifest not found, skipping..."
        fi
    done
    
    log_success "All manifests applied"
}

# Wait for deployments to be ready
wait_for_deployments() {
    log_info "Waiting for deployments to be ready..."
    
    # Wait for infrastructure components
    kubectl wait --for=condition=available --timeout=600s deployment/mongodb -n infrastructure
    kubectl wait --for=condition=available --timeout=600s deployment/redis -n infrastructure
    kubectl wait --for=condition=available --timeout=600s deployment/traefik -n infrastructure
    
    # Wait for VPN components
    kubectl wait --for=condition=available --timeout=600s deployment/gluetun-airvpn -n vpn
    
    # Wait for media services
    kubectl wait --for=condition=available --timeout=600s deployment/stremio -n media-stack
    kubectl wait --for=condition=available --timeout=600s deployment/jellyfin -n media-stack
    
    # Wait for monitoring
    kubectl wait --for=condition=available --timeout=600s deployment/homepage -n monitoring
    
    log_success "All deployments are ready"
}

# Display access information
display_access_info() {
    log_info "Deployment completed successfully!"
    echo
    log_info "Access Information:"
    echo "===================="
    
    # Get LoadBalancer IP
    local lb_ip
    lb_ip=$(kubectl get svc traefik -n infrastructure -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "pending")
    
    if [[ "$lb_ip" == "pending" || -z "$lb_ip" ]]; then
        lb_ip=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="ExternalIP")].address}' 2>/dev/null || echo "localhost")
    fi
    
    echo "Load Balancer IP: $lb_ip"
    echo
    echo "Services (update your DNS to point to $lb_ip):"
    echo "- Traefik Dashboard: https://traefik.bolabaden.org"
    echo "- Homepage Dashboard: https://dashboard.bolabaden.org"
    echo "- Stremio: https://stremio.bolabaden.org"
    echo "- Jellyfin: https://jellyfin.bolabaden.org"
    echo "- Prowlarr: https://prowlarr.bolabaden.org"
    echo "- Radarr: https://radarr.bolabaden.org"
    echo "- Sonarr: https://sonarr.bolabaden.org"
    echo "- qBittorrent: https://qbittorrent.bolabaden.org"
    echo "- Logs (Dozzle): https://logs.bolabaden.org"
    echo "- Speedtest: https://speedtest.bolabaden.org"
    echo "- Code Server: https://code.bolabaden.org"
    echo "- AIOStreams: https://aiostreams.bolabaden.org"
    echo "- Torrentio: https://torrentio.bolabaden.org"
    echo "- MediaFusion: https://mediafusion.bolabaden.org"
    echo "- Comet: https://comet.bolabaden.org"
    echo
    echo "Stremio Addon URLs:"
    echo "- AIOStreams: https://aiostreams.bolabaden.org/manifest.json"
    echo "- MediaFusion: https://mediafusion.bolabaden.org/manifest.json"
    echo "- Comet: https://comet.bolabaden.org/manifest.json"
    echo
    log_warning "Remember to:"
    echo "1. Update DNS records to point your domains to $lb_ip"
    echo "2. Update the secrets in k8s/secrets.yaml with your actual API keys"
    echo "3. Configure your VPN credentials in the VPN secrets"
    echo "4. Set up your domain names in the ingress configurations"
    echo
    log_info "To check the status of your deployments:"
    echo "kubectl get pods --all-namespaces"
    echo
    log_info "To view logs for a specific service:"
    echo "kubectl logs -f deployment/<service-name> -n <namespace>"
}

# Cleanup function
cleanup() {
    log_info "Cleaning up..."
    # Add any cleanup tasks here if needed
}

# Main execution
main() {
    log_info "Starting Media Stack Kubernetes Deployment"
    echo "=========================================="
    
    # Set trap for cleanup
    trap cleanup EXIT
    
    # Check prerequisites
    check_root
    
    # Install components
    install_tools
    install_k3s
    setup_kubectl
    
    # Deploy the stack
    deploy_manifests
    wait_for_deployments
    
    # Show access information
    display_access_info
    
    log_success "Deployment completed successfully!"
}

# Run main function
main "$@" 