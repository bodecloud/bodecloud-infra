#!/bin/bash

# =============================================================================
# MEDIA STACK SETUP SCRIPT
# =============================================================================
# This script helps you set up the media stack configuration for the first time.
# It will guide you through creating your cluster-config.yml file.
# =============================================================================

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/cluster-config.yml"
EXAMPLE_FILE="${SCRIPT_DIR}/cluster-config.yml.example"

# Logging functions
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

prompt() {
    echo -e "${CYAN}[INPUT]${NC} $1"
}

# Function to display welcome message
show_welcome() {
    cat << 'EOF'

 ███╗   ███╗███████╗██████╗ ██╗ █████╗     ███████╗████████╗ █████╗  ██████╗██╗  ██╗
 ████╗ ████║██╔════╝██╔══██╗██║██╔══██╗    ██╔════╝╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝
 ██╔████╔██║█████╗  ██║  ██║██║███████║    ███████╗   ██║   ███████║██║     █████╔╝ 
 ██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║    ╚════██║   ██║   ██╔══██║██║     ██╔═██╗ 
 ██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║    ███████║   ██║   ██║  ██║╚██████╗██║  ██╗
 ╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝    ╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝
                                                                                      
                           K3s Deployment Setup

EOF

    echo -e "${GREEN}Welcome to the Media Stack K3s Deployment Setup!${NC}"
    echo
    echo "This script will help you create your cluster-config.yml file with your"
    echo "specific environment values. You'll need to have the following ready:"
    echo
    echo "📋 Required Information:"
    echo "  • Your domain name (e.g., yourdomain.org)"
    echo "  • Tailscale network details and node IPs"
    echo "  • Cloudflare API credentials (for SSL certificates)"
    echo "  • VPN provider credentials (Premiumize, Real-Debrid)"
    echo "  • API keys for various services"
    echo
    echo "🔧 Prerequisites:"
    echo "  • K3s cluster already set up and running"
    echo "  • All nodes connected to Tailscale"
    echo "  • SSH access to all nodes"
    echo
}

# Function to check if configuration exists
check_existing_config() {
    if [[ -f "$CONFIG_FILE" ]]; then
        warning "Configuration file already exists: $CONFIG_FILE"
        echo
        prompt "Do you want to:"
        echo "  1) Backup existing and create new configuration"
        echo "  2) Edit existing configuration"
        echo "  3) Exit"
        echo
        read -p "Enter your choice (1-3): " choice
        
        case $choice in
            1)
                local backup_file="${CONFIG_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
                mv "$CONFIG_FILE" "$backup_file"
                success "Existing configuration backed up to: $backup_file"
                return 0
                ;;
            2)
                log "Opening existing configuration for editing..."
                ${EDITOR:-nano} "$CONFIG_FILE"
                exit 0
                ;;
            3)
                log "Exiting setup."
                exit 0
                ;;
            *)
                error "Invalid choice. Exiting."
                exit 1
                ;;
        esac
    fi
}

# Function to gather basic cluster information
gather_cluster_info() {
    log "Gathering cluster information..."
    echo
    
    # Domain configuration
    prompt "Enter your domain name (e.g., yourdomain.org):"
    read -p "> " DOMAIN
    
    prompt "Enter your subdomain (e.g., media):"
    read -p "> " SUBDOMAIN
    
    prompt "Enter your timezone (e.g., America/Chicago):"
    read -p "> " TIMEZONE
    
    # Tailscale configuration
    echo
    log "Tailscale Configuration"
    prompt "Enter your Tailscale base domain (e.g., your-tailnet.ts.net):"
    read -p "> " TS_BASE_DOMAIN
    
    prompt "Enter your master node hostname:"
    read -p "> " MASTER_HOSTNAME
    
    prompt "Enter your master node Tailscale IP:"
    read -p "> " MASTER_IP
    
    prompt "Enter your master node external IP:"
    read -p "> " MASTER_EXTERNAL_IP
    
    prompt "Enter your Tailscale auth key:"
    read -s -p "> " TS_AUTHKEY
    echo
    
    # Agent nodes
    echo
    log "Agent Nodes Configuration"
    AGENT_NODES=()
    
    while true; do
        prompt "Enter agent node hostname (or 'done' to finish):"
        read -p "> " agent_hostname
        
        if [[ "$agent_hostname" == "done" ]]; then
            break
        fi
        
        prompt "Enter Tailscale IP for $agent_hostname:"
        read -p "> " agent_ip
        
        AGENT_NODES+=("$agent_hostname:$agent_ip")
        success "Added agent node: $agent_hostname ($agent_ip)"
    done
    
    if [[ ${#AGENT_NODES[@]} -eq 0 ]]; then
        warning "No agent nodes configured. You can add them later."
    fi
}

# Function to gather certificate information
gather_cert_info() {
    echo
    log "Certificate Configuration"
    
    prompt "Enter your email address for Let's Encrypt:"
    read -p "> " CERT_EMAIL
    
    prompt "Enter your Cloudflare email:"
    read -p "> " CF_EMAIL
    
    prompt "Enter your Cloudflare API key:"
    read -s -p "> " CF_API_KEY
    echo
    
    prompt "Enter your Cloudflare DNS API token:"
    read -s -p "> " CF_DNS_TOKEN
    echo
    
    prompt "Enter your Cloudflare Zone ID:"
    read -p "> " CF_ZONE_ID
}

# Function to gather VPN information
gather_vpn_info() {
    echo
    log "VPN Provider Configuration"
    
    # Premiumize
    prompt "Do you have Premiumize? (y/n):"
    read -p "> " has_premiumize
    
    if [[ "$has_premiumize" =~ ^[Yy]$ ]]; then
        prompt "Enter your Premiumize customer ID:"
        read -p "> " PREMIUMIZE_ID
        
        prompt "Enter your Premiumize API key:"
        read -s -p "> " PREMIUMIZE_KEY
        echo
        PREMIUMIZE_ENABLED="true"
    else
        PREMIUMIZE_ENABLED="false"
        PREMIUMIZE_ID=""
        PREMIUMIZE_KEY=""
    fi
    
    # Real-Debrid
    prompt "Do you have Real-Debrid? (y/n):"
    read -p "> " has_realdebrid
    
    if [[ "$has_realdebrid" =~ ^[Yy]$ ]]; then
        prompt "Enter your Real-Debrid API key:"
        read -s -p "> " REALDEBRID_KEY
        echo
        REALDEBRID_ENABLED="true"
    else
        REALDEBRID_ENABLED="false"
        REALDEBRID_KEY=""
    fi
}

# Function to generate configuration file
generate_config() {
    log "Generating configuration file..."
    
    if [[ ! -f "$EXAMPLE_FILE" ]]; then
        error "Example configuration file not found: $EXAMPLE_FILE"
        exit 1
    fi
    
    # Copy example file
    cp "$EXAMPLE_FILE" "$CONFIG_FILE"
    
    # Replace placeholders with actual values
    sed -i "s/your-domain\.org/$DOMAIN/g" "$CONFIG_FILE"
    sed -i "s/your-subdomain/$SUBDOMAIN/g" "$CONFIG_FILE"
    sed -i "s/America\/Chicago/$TIMEZONE/g" "$CONFIG_FILE"
    sed -i "s/your-tailnet\.ts\.net/$TS_BASE_DOMAIN/g" "$CONFIG_FILE"
    sed -i "s/your-master-hostname/$MASTER_HOSTNAME/g" "$CONFIG_FILE"
    sed -i "s/your-master/$MASTER_HOSTNAME/g" "$CONFIG_FILE"
    sed -i "s/100\.x\.x\.x/$MASTER_IP/g" "$CONFIG_FILE"
    sed -i "s/x\.x\.x\.x/$MASTER_EXTERNAL_IP/g" "$CONFIG_FILE"
    sed -i "s/tskey-auth-XXXXXXXXXX/$TS_AUTHKEY/g" "$CONFIG_FILE"
    sed -i "s/your-email@example\.com/$CERT_EMAIL/g" "$CONFIG_FILE"
    sed -i "s/your-cloudflare-api-key/$CF_API_KEY/g" "$CONFIG_FILE"
    sed -i "s/your-dns-api-token/$CF_DNS_TOKEN/g" "$CONFIG_FILE"
    sed -i "s/your-zone-id/$CF_ZONE_ID/g" "$CONFIG_FILE"
    sed -i "s/your-zone-api-token/$CF_DNS_TOKEN/g" "$CONFIG_FILE"
    
    if [[ "$PREMIUMIZE_ENABLED" == "true" ]]; then
        sed -i "s/your-customer-id/$PREMIUMIZE_ID/g" "$CONFIG_FILE"
        sed -i "s/your-premiumize-api-key/$PREMIUMIZE_KEY/g" "$CONFIG_FILE"
    fi
    
    if [[ "$REALDEBRID_ENABLED" == "true" ]]; then
        sed -i "s/your-realdebrid-api-key/$REALDEBRID_KEY/g" "$CONFIG_FILE"
    fi
    
    # Update agent nodes
    if [[ ${#AGENT_NODES[@]} -gt 0 ]]; then
        # This is a simplified replacement - for complex YAML, manual editing might be needed
        log "Agent nodes configured. You may need to manually edit the agents section."
    fi
    
    # Generate secure tokens
    local k3s_token=$(openssl rand -hex 32)
    local tinyauth_secret=$(openssl rand -hex 32)
    local oauth_secret=$(openssl rand -hex 32)
    local auth_secret=$(openssl rand -hex 32)
    
    sed -i "s/CHANGE-THIS-TO-A-SECURE-TOKEN/$k3s_token/g" "$CONFIG_FILE"
    sed -i "s/your-tinyauth-secret/$tinyauth_secret/g" "$CONFIG_FILE"
    sed -i "s/your-oauth-secret/$oauth_secret/g" "$CONFIG_FILE"
    sed -i "s/your-auth-secret/$auth_secret/g" "$CONFIG_FILE"
    
    success "Configuration file generated: $CONFIG_FILE"
}

# Function to show next steps
show_next_steps() {
    echo
    success "Setup completed successfully!"
    echo
    log "Next steps:"
    echo
    echo "1. Review and edit your configuration:"
    echo "   ${EDITOR:-nano} $CONFIG_FILE"
    echo
    echo "2. Add any missing API keys and credentials"
    echo
    echo "3. Configure your agent nodes in the 'agents' section"
    echo
    echo "4. Deploy your media stack:"
    echo "   ./deploy.sh deploy"
    echo
    echo "5. Verify the deployment:"
    echo "   ./deploy.sh verify"
    echo
    warning "Important: Review the configuration file before deploying!"
    warning "Make sure all API keys and credentials are correct."
    echo
    log "For help and documentation, see: README.md"
}

# Function to show help
show_help() {
    cat << EOF
Media Stack Setup Script

USAGE:
    $0 [COMMAND]

COMMANDS:
    setup       Interactive setup (default)
    edit        Edit existing configuration
    validate    Validate configuration file
    help        Show this help message

EXAMPLES:
    $0              # Run interactive setup
    $0 setup        # Run interactive setup
    $0 edit         # Edit existing config
    $0 validate     # Validate config file

EOF
}

# Function to validate configuration
validate_config() {
    if [[ ! -f "$CONFIG_FILE" ]]; then
        error "Configuration file not found: $CONFIG_FILE"
        error "Run '$0 setup' to create one."
        exit 1
    fi
    
    log "Validating configuration file..."
    
    # Check if yq is available
    if ! command -v yq &> /dev/null; then
        warning "yq not found. Installing for validation..."
        sudo wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64
        sudo chmod +x /usr/local/bin/yq
    fi
    
    # Basic YAML syntax check
    if ! yq eval '.' "$CONFIG_FILE" > /dev/null 2>&1; then
        error "Invalid YAML syntax in configuration file"
        exit 1
    fi
    
    # Check required fields
    local required_fields=(
        ".cluster.domain"
        ".cluster.tailscale.authkey"
        ".cluster.tailscale.nodes.master.hostname"
        ".certificates.cloudflare.api_key"
    )
    
    local missing_fields=()
    for field in "${required_fields[@]}"; do
        local value=$(yq eval "$field" "$CONFIG_FILE" 2>/dev/null)
        if [[ "$value" == "null" || "$value" == "your-"* || -z "$value" ]]; then
            missing_fields+=("$field")
        fi
    done
    
    if [[ ${#missing_fields[@]} -gt 0 ]]; then
        error "Missing or incomplete configuration fields:"
        for field in "${missing_fields[@]}"; do
            echo "  - $field"
        done
        exit 1
    fi
    
    success "Configuration validation passed!"
}

# Main function
main() {
    local command="${1:-setup}"
    
    case "$command" in
        "setup")
            show_welcome
            check_existing_config
            gather_cluster_info
            gather_cert_info
            gather_vpn_info
            generate_config
            show_next_steps
            ;;
        "edit")
            if [[ ! -f "$CONFIG_FILE" ]]; then
                error "Configuration file not found. Run '$0 setup' first."
                exit 1
            fi
            ${EDITOR:-nano} "$CONFIG_FILE"
            ;;
        "validate")
            validate_config
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

# Run main function
main "$@" 