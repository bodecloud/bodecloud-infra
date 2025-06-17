#!/bin/bash

# =============================================================================
# Media Stack K3s Cluster Deployment Script
# =============================================================================
# This script deploys a complete K3s-based media stack cluster
# Assumes K3s and Tailscale are already installed via cloud-init
# =============================================================================

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}==============================================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}==============================================================================${NC}"
}

# Function to check prerequisites
check_prerequisites() {
    print_header "CHECKING PREREQUISITES"
    
    # Check if .env file exists
    if [[ ! -f "${ROOT_DIR}/.env" ]]; then
        print_error ".env file not found at ${ROOT_DIR}/.env"
        exit 1
    fi
    print_status "Found .env file"
    
    # Check if ansible is installed
    if ! command -v ansible-playbook &> /dev/null; then
        print_error "ansible-playbook is not installed"
        print_status "Installing ansible..."
        sudo apt update
        sudo apt install -y ansible
    fi
    print_status "Ansible is available"
    
    # Check if ansible kubernetes collection is installed
    if ! ansible-galaxy collection list | grep -q kubernetes.core; then
        print_status "Installing kubernetes.core collection..."
        ansible-galaxy collection install kubernetes.core
    fi
    print_status "Kubernetes collection is available"
    
    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        print_status "Installing kubectl..."
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
        sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
        rm kubectl
    fi
    print_status "kubectl is available"
    
    # Source environment variables
    set -a
    source "${ROOT_DIR}/.env"
    set +a
    
    print_status "Environment variables loaded"
}

# Function to validate configuration
validate_configuration() {
    print_header "VALIDATING CONFIGURATION"
    
    # Check required environment variables
    required_vars=(
        "DOMAIN"
        "TS_BASE_DOMAIN"
        "TS_AUTHKEY"
        "CF_EMAIL"
        "CF_DNS_API_TOKEN"
        "EMAIL"
    )
    
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var:-}" ]]; then
            print_error "Required environment variable $var is not set"
            exit 1
        fi
        print_status "✓ $var is set"
    done
    
    print_status "Configuration validation passed"
}

# Function to test connectivity to nodes
test_connectivity() {
    print_header "TESTING NODE CONNECTIVITY"
    
    # Test connectivity using ansible ping
    if ansible all -i "${SCRIPT_DIR}/inventory/hosts.yml" -m ping --private-key ~/.ssh/id_rsa; then
        print_status "All nodes are reachable"
    else
        print_error "Some nodes are not reachable. Please check your Tailscale configuration and SSH keys."
        exit 1
    fi
}

# Function to deploy the cluster
deploy_cluster() {
    print_header "DEPLOYING KUBERNETES CLUSTER"
    
    # Change to the script directory
    cd "${SCRIPT_DIR}"
    
    # Run the ansible playbook
    if ansible-playbook \
        -i inventory/hosts.yml \
        site.yml \
        --private-key ~/.ssh/id_rsa \
        -v; then
        print_status "Cluster deployment completed successfully"
    else
        print_error "Cluster deployment failed"
        exit 1
    fi
}

# Function to verify deployment
verify_deployment() {
    print_header "VERIFYING DEPLOYMENT"
    
    # Set kubeconfig
    export KUBECONFIG="${SCRIPT_DIR}/kubeconfig"
    
    # Check if kubeconfig exists
    if [[ ! -f "${KUBECONFIG}" ]]; then
        print_error "Kubeconfig not found at ${KUBECONFIG}"
        exit 1
    fi
    
    # Test kubectl connectivity
    if kubectl cluster-info &> /dev/null; then
        print_status "✓ Kubectl connectivity working"
    else
        print_error "✗ Kubectl connectivity failed"
        exit 1
    fi
    
    # Check node status
    print_status "Cluster nodes:"
    kubectl get nodes -o wide
    
    # Check pod status
    print_status "Cluster pods:"
    kubectl get pods --all-namespaces
    
    print_status "Verification completed"
}

# Function to display access information
display_access_info() {
    print_header "CLUSTER ACCESS INFORMATION"
    
    source "${ROOT_DIR}/.env"
    
    echo -e "${GREEN}Cluster deployed successfully!${NC}"
    echo ""
    echo "Domain: ${DOMAIN}"
    echo "Tailscale Network: ${TS_BASE_DOMAIN}"
    echo "Kubeconfig: ${SCRIPT_DIR}/kubeconfig"
    echo ""
    echo "To access your cluster:"
    echo "  export KUBECONFIG=${SCRIPT_DIR}/kubeconfig"
    echo "  kubectl get nodes"
    echo ""
    echo "Service URLs (after services are deployed):"
    echo "  https://traefik.${DOMAIN}      - Traefik Dashboard"
    echo "  https://plex.${DOMAIN}         - Plex Media Server"
    echo "  https://prowlarr.${DOMAIN}     - Prowlarr Indexer"
    echo "  https://comet.${DOMAIN}        - Comet Stremio Addon"
    echo ""
    echo "Documentation: ${SCRIPT_DIR}/CLUSTER_ACCESS.md"
}

# Main execution
main() {
    print_header "MEDIA STACK K3S CLUSTER DEPLOYMENT"
    
    check_prerequisites
    validate_configuration
    test_connectivity
    deploy_cluster
    verify_deployment
    display_access_info
    
    print_header "DEPLOYMENT COMPLETE"
}

# Help function
show_help() {
    cat << EOF
Media Stack K3s Cluster Deployment Script

Usage: $0 [OPTION]

Options:
    -h, --help      Show this help message
    --check-only    Only run prerequisite and connectivity checks
    --verify-only   Only run deployment verification

This script deploys a complete K3s-based media stack cluster.
Prerequisites:
- K3s and Tailscale must be installed on all nodes via cloud-init
- .env file must exist in the project root
- SSH key-based authentication to all nodes

EOF
}

# Parse command line arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    --check-only)
        check_prerequisites
        validate_configuration
        test_connectivity
        print_status "Checks completed successfully"
        exit 0
        ;;
    --verify-only)
        verify_deployment
        exit 0
        ;;
    "")
        main
        ;;
    *)
        print_error "Unknown option: $1"
        show_help
        exit 1
        ;;
esac 