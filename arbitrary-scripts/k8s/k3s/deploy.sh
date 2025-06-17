#!/bin/bash

# =============================================================================
# MEDIA STACK K3S DEPLOYMENT SCRIPT
# =============================================================================
# This script deploys the complete K3s-based media stack using Ansible.
# It loads configuration from cluster-config.yml and validates prerequisites.
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
CONFIG_FILE="${SCRIPT_DIR}/cluster-config.yml"

# Logging function
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Function to check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if running on Ubuntu/Debian
    if [[ ! -f /etc/os-release ]]; then
        error "Cannot determine OS. This script is designed for Ubuntu/Debian systems."
        exit 1
    fi
    
    source /etc/os-release
    if [[ "$ID" != "ubuntu" && "$ID" != "debian" ]]; then
        warning "This script is optimized for Ubuntu/Debian. Your OS: $ID"
    fi
    
    # Check if cluster-config.yml exists
    if [[ ! -f "$CONFIG_FILE" ]]; then
        error "Configuration file not found: $CONFIG_FILE"
        error "Please create cluster-config.yml with your environment-specific values."
        exit 1
    fi
    
    # Check if Ansible is installed
    if ! command -v ansible-playbook &> /dev/null; then
        log "Installing Ansible..."
        sudo apt update
        sudo apt install -y software-properties-common
        sudo add-apt-repository --yes --update ppa:ansible/ansible
        sudo apt install -y ansible
    fi
    
    # Check if required Ansible collections are installed
    log "Installing required Ansible collections..."
    ansible-galaxy collection install kubernetes.core community.general --force
    
    # Install required Python packages
    log "Installing required Python packages..."
    pip3 install --user kubernetes pyyaml jinja2
    
    success "Prerequisites check completed"
}

# Function to validate configuration
validate_config() {
    log "Validating cluster configuration..."
    
    # Check if yq is installed for YAML parsing
    if ! command -v yq &> /dev/null; then
        log "Installing yq for YAML parsing..."
        sudo wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64
        sudo chmod +x /usr/local/bin/yq
    fi
    
    # Validate required fields in config
    local required_fields=(
        ".cluster.domain"
        ".cluster.tailscale.authkey"
        ".cluster.tailscale.nodes.master.hostname"
        ".cluster.tailscale.nodes.master.tailscale_ip"
        ".certificates.cloudflare.api_key"
        ".vpn.providers.premiumize.api_key"
    )
    
    for field in "${required_fields[@]}"; do
        if ! yq eval "$field" "$CONFIG_FILE" &> /dev/null; then
            error "Required configuration field missing: $field"
            exit 1
        fi
    done
    
    # Extract domain for validation
    local domain
    domain=$(yq eval '.cluster.domain' "$CONFIG_FILE")
    
    if [[ "$domain" == "null" || -z "$domain" ]]; then
        error "Domain not configured in cluster-config.yml"
        exit 1
    fi
    
    success "Configuration validation completed"
}

# Function to check node connectivity
check_node_connectivity() {
    log "Checking node connectivity..."
    
    # Extract node IPs from config
    local master_ip
    master_ip=$(yq eval '.cluster.tailscale.nodes.master.tailscale_ip' "$CONFIG_FILE")
    
    local agent_ips
    agent_ips=$(yq eval '.cluster.tailscale.nodes.agents[].tailscale_ip' "$CONFIG_FILE")
    
    # Test connectivity to master
    if ping -c 1 -W 5 "$master_ip" &> /dev/null; then
        success "Master node ($master_ip) is reachable"
    else
        error "Master node ($master_ip) is not reachable"
        exit 1
    fi
    
    # Test connectivity to agents
    while IFS= read -r agent_ip; do
        if ping -c 1 -W 5 "$agent_ip" &> /dev/null; then
            success "Agent node ($agent_ip) is reachable"
        else
            warning "Agent node ($agent_ip) is not reachable"
        fi
    done <<< "$agent_ips"
}

# Function to generate Kubernetes configuration
generate_k8s_config() {
    log "Generating Kubernetes configuration..."
    
    if ansible-playbook -i inventory/hosts.yml generate-k8s-config.yml; then
        success "Kubernetes configuration generated successfully"
    else
        error "Failed to generate Kubernetes configuration"
        exit 1
    fi
}

# Function to deploy the stack
deploy_stack() {
    local deployment_type="${1:-full}"
    
    log "Starting deployment (type: $deployment_type)..."
    
    case "$deployment_type" in
        "infrastructure")
            log "Deploying infrastructure components only..."
            ansible-playbook -i inventory/hosts.yml site.yml --tags "infrastructure"
            ;;
        "services")
            log "Deploying services only..."
            ansible-playbook -i inventory/hosts.yml site.yml --tags "applications"
            ;;
        "full")
            log "Deploying complete stack..."
            ansible-playbook -i inventory/hosts.yml site.yml
            ;;
        *)
            error "Invalid deployment type: $deployment_type"
            error "Valid types: infrastructure, services, full"
            exit 1
            ;;
    esac
}

# Function to verify deployment
verify_deployment() {
    log "Verifying deployment..."
    
    # Extract master IP for kubectl commands
    local master_ip
    master_ip=$(yq eval '.cluster.tailscale.nodes.master.tailscale_ip' "$CONFIG_FILE")
    
    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        log "Installing kubectl..."
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
        sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
        rm kubectl
    fi
    
    # Set up kubectl config to connect to the cluster
    export KUBECONFIG="$HOME/.kube/config"
    
    # Run verification playbook
    ansible-playbook -i inventory/hosts.yml site.yml --tags "verification"
    
    success "Deployment verification completed"
}

# Function to display help
show_help() {
    cat << EOF
Media Stack K3s Deployment Script

USAGE:
    $0 [COMMAND] [OPTIONS]

COMMANDS:
    deploy [TYPE]     Deploy the media stack
                      TYPE: infrastructure, services, full (default: full)
    
    config           Generate Kubernetes configuration only
    
    verify           Verify existing deployment
    
    check            Check prerequisites and connectivity
    
    help             Show this help message

EXAMPLES:
    $0 deploy                    # Deploy complete stack
    $0 deploy infrastructure     # Deploy infrastructure only
    $0 deploy services          # Deploy services only
    $0 config                   # Generate K8s config only
    $0 verify                   # Verify deployment
    $0 check                    # Check prerequisites

CONFIGURATION:
    Edit cluster-config.yml to customize your deployment.
    This file contains all environment-specific values.

EOF
}

# Function to display deployment summary
show_summary() {
    local domain
    domain=$(yq eval '.cluster.domain' "$CONFIG_FILE" 2>/dev/null || echo "unknown")
    
    cat << EOF

=============================================================================
DEPLOYMENT COMPLETED SUCCESSFULLY
=============================================================================

Your media stack is now deployed and accessible at:

🌐 Web Interfaces:
   • Traefik Dashboard: https://traefik.$domain
   • Plex Media Server: https://plex.$domain
   • Prowlarr (Indexer): https://prowlarr.$domain
   • Comet (Stremio): https://comet.$domain

🔧 Management:
   • Homepage Dashboard: https://homepage.$domain
   • Dozzle (Logs): https://dozzle.$domain

📊 Monitoring:
   • Speedtest Tracker: https://speedtest.$domain

🔐 Authentication:
   • TinyAuth: https://auth.$domain

Useful Commands:
   • Check cluster status: kubectl get nodes
   • View all pods: kubectl get pods -A
   • Check VPN gateways: kubectl get pods -n vpn-gateway
   • View service logs: kubectl logs -f deployment/<service-name>

Configuration:
   • All settings are stored in cluster-config.yml
   • Kubernetes configs are in namespace: media-stack-config
   • To modify services, edit cluster-config.yml and re-run deployment

=============================================================================

EOF
}

# Main execution
main() {
    local command="${1:-deploy}"
    
    case "$command" in
        "deploy")
            check_prerequisites
            validate_config
            check_node_connectivity
            generate_k8s_config
            deploy_stack "${2:-full}"
            verify_deployment
            show_summary
            ;;
        "config")
            check_prerequisites
            validate_config
            generate_k8s_config
            ;;
        "verify")
            check_prerequisites
            validate_config
            verify_deployment
            ;;
        "check")
            check_prerequisites
            validate_config
            check_node_connectivity
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@" 