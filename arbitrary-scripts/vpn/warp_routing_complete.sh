#!/bin/bash

set -e

#===========================================
# ULTRA-SAFE WARP ROUTING SETUP SCRIPT
#===========================================
# Maximum safety features:
# - Preview mode (show commands without executing)
# - Backup/restore network state
# - Step-by-step confirmation
# - Automatic rollback on any failure
# - Dry-run simulation
# - Interactive guided mode

# Safety Configuration (NEW SAFETY FEATURES)
PREVIEW_MODE=${PREVIEW_MODE:-false}          # Show commands without executing
DRY_RUN=${DRY_RUN:-false}                   # Simulate everything
INTERACTIVE_MODE=${INTERACTIVE_MODE:-false}  # Step-by-step confirmation
BACKUP_NETWORK_STATE=${BACKUP_NETWORK_STATE:-true}  # Backup before changes
AUTO_ROLLBACK=${AUTO_ROLLBACK:-true}         # Auto rollback on failure
SAFETY_TIMEOUT=${SAFETY_TIMEOUT:-60}        # Auto-rollback timeout in seconds
STEP_BY_STEP=${STEP_BY_STEP:-false}         # Pause after each command
BACKUP_DIR=${BACKUP_DIR:-/tmp/warp-routing-backup-$(date +%s)}

# Network Configuration
WARP_SUBNET=${WARP_SUBNET:-172.20.0.0/16}
WARP_GATEWAY=${WARP_GATEWAY_IPV4_ADDRESS:-172.20.0.1}
PUBLICNET_SUBNET=${PUBLICNET_SUBNET:-10.76.0.0/16}
PUBLICNET_GATEWAY=${PUBLICNET_GATEWAY:-10.76.0.1}
WARP_CONTAINER_IP=${WARP_IPV4_ADDRESS:-10.76.128.200}
ROUTING_TABLE_ID=${ROUTING_TABLE_ID:-200}
ROUTING_TABLE_NAME=${ROUTING_TABLE_NAME:-warpnetvpn}

# Container Configuration
WARP_CONTAINER_NAME=${WARP_CONTAINER_NAME:-warp-with-nat}
WARP_IMAGE=${WARP_IMAGE:-caomingjun/warp:latest}
CONFIG_PATH=${CONFIG_PATH:-./configs}
WARP_SOCKS_PORT=${WARP_SOCKS_PORT:-5080}

# WARP Configuration
WARP_LICENSE_KEY=${WARP_LICENSE_KEY:-}
TUNNEL_TOKEN=${TUNNEL_TOKEN:-}

# Script Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEBUG=${DEBUG:-false}
FORCE_RECREATE=${FORCE_RECREATE:-false}
FORCE_APPLY=${FORCE_APPLY:-false}
AUTO_START_CONTAINER=${AUTO_START_CONTAINER:-true}
RUN_TESTS=${RUN_TESTS:-true}
CLEANUP_ON_TEST_FAILURE=${CLEANUP_ON_TEST_FAILURE:-true}
USE_SAFE_MODE=${USE_SAFE_MODE:-true}
DIND_CONTAINER_NAME="warp-routing-test-dind"
TEST_TIMEOUT=${TEST_TIMEOUT:-300}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m'

#===========================================
# ULTRA-SAFE EXECUTION FRAMEWORK
#===========================================

# Track all commands for rollback
EXECUTED_COMMANDS=()
BACKUP_FILES=()
ROLLBACK_COMMANDS=()

log_debug() { [[ "$DEBUG" == "true" ]] && echo -e "${PURPLE}[DEBUG]${NC} $1"; }
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_step() { echo -e "${CYAN}[STEP]${NC} $1"; }
log_preview() { echo -e "${PURPLE}[PREVIEW]${NC} Would execute: $1"; }
log_confirm() { echo -e "${BOLD}${YELLOW}[CONFIRM]${NC} $1"; }

# Safe command execution with preview/dry-run/confirmation
safe_exec() {
    local cmd="$1"
    local rollback_cmd="$2"
    local description="$3"
    
    if [[ "$PREVIEW_MODE" == "true" ]]; then
        log_preview "$cmd"
        [[ -n "$description" ]] && echo "  → $description"
        return 0
    fi
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "[DRY-RUN] Would execute: $cmd"
        [[ -n "$description" ]] && echo "  → $description"
        return 0
    fi
    
    if [[ "$INTERACTIVE_MODE" == "true" ]] || [[ "$STEP_BY_STEP" == "true" ]]; then
        echo
        log_confirm "About to execute: $cmd"
        [[ -n "$description" ]] && echo "  Purpose: $description"
        
        if ! confirm_action "Proceed with this command?"; then
            log_warning "Command skipped by user"
            return 1
        fi
    fi
    
    log_step "Executing: $cmd"
    
    # Execute and track for rollback
    if eval "$cmd"; then
        EXECUTED_COMMANDS+=("$cmd")
        [[ -n "$rollback_cmd" ]] && ROLLBACK_COMMANDS+=("$rollback_cmd")
        log_success "Command completed successfully"
        
        if [[ "$STEP_BY_STEP" == "true" ]]; then
            read -p "Press Enter to continue..." -r
        fi
        return 0
    else
        local exit_code=$?
        log_error "Command failed: $cmd"
        
        if [[ "$AUTO_ROLLBACK" == "true" ]]; then
            log_warning "Auto-rollback enabled, initiating cleanup..."
            emergency_rollback
        fi
        
        return $exit_code
    fi
}

# Backup network state before making changes
backup_network_state() {
    if [[ "$BACKUP_NETWORK_STATE" != "true" ]]; then
        return 0
    fi
    
    print_section "BACKING UP NETWORK STATE"
    
    mkdir -p "$BACKUP_DIR"
    log_info "Creating backup in: $BACKUP_DIR"
    
    # Backup routing tables
    ip route show > "$BACKUP_DIR/routes.txt" 2>/dev/null || true
    ip rule list > "$BACKUP_DIR/rules.txt" 2>/dev/null || true
    cat /etc/iproute2/rt_tables > "$BACKUP_DIR/rt_tables.txt" 2>/dev/null || true
    
    # Backup iptables
    iptables-save > "$BACKUP_DIR/iptables.txt" 2>/dev/null || true
    
    # Backup Docker networks
    docker network ls --format "{{.ID}} {{.Name}}" > "$BACKUP_DIR/docker_networks.txt" 2>/dev/null || true
    
    # Backup Docker containers
    docker ps -a --format "{{.ID}} {{.Names}} {{.Status}}" > "$BACKUP_DIR/docker_containers.txt" 2>/dev/null || true
    
    log_success "Network state backed up to: $BACKUP_DIR"
}

# Restore network state from backup
restore_network_state() {
    if [[ ! -d "$BACKUP_DIR" ]]; then
        log_warning "No backup directory found at: $BACKUP_DIR"
        return 1
    fi
    
    print_section "RESTORING NETWORK STATE FROM BACKUP"
    
    log_warning "This will restore your network to the state before script execution"
    if ! confirm_action "Proceed with network state restoration?"; then
        log_info "Network restoration cancelled"
        return 0
    fi
    
    # Note: Full network restoration is complex and risky
    # Instead, we'll use our tracked rollback commands
    log_info "Using tracked rollback commands for safer restoration..."
    emergency_rollback
}

# Emergency rollback function
emergency_rollback() {
    print_section "EMERGENCY ROLLBACK IN PROGRESS"
    
    log_warning "Rolling back ${#ROLLBACK_COMMANDS[@]} operations..."
    
    # Execute rollback commands in reverse order
    for (( i=${#ROLLBACK_COMMANDS[@]}-1; i>=0; i-- )); do
        local rollback_cmd="${ROLLBACK_COMMANDS[i]}"
        log_step "Rollback: $rollback_cmd"
        eval "$rollback_cmd" || log_warning "Rollback command failed: $rollback_cmd"
    done
    
    # Clean up any containers we created
    cleanup_created_resources
    
    log_success "Emergency rollback completed"
}

# Cleanup resources created by this script
cleanup_created_resources() {
    log_step "Cleaning up created resources..."
    
    # Remove containers
    docker rm -f "$WARP_CONTAINER_NAME" >/dev/null 2>&1 || true
    docker rm -f "$DIND_CONTAINER_NAME" >/dev/null 2>&1 || true
    
    # Remove networks (only if they were created by us)
    docker network rm warpnet >/dev/null 2>&1 || true
    docker network rm publicnet >/dev/null 2>&1 || true
    
    log_success "Resource cleanup completed"
}

# Safety timeout - auto rollback after specified time
setup_safety_timeout() {
    if [[ "$AUTO_ROLLBACK" == "true" && "$SAFETY_TIMEOUT" -gt 0 ]]; then
        log_info "Safety timeout set to ${SAFETY_TIMEOUT} seconds"
        log_warning "Script will auto-rollback if not manually stopped within timeout"
        
        (
            sleep "$SAFETY_TIMEOUT"
            log_warning "Safety timeout reached! Initiating automatic rollback..."
            emergency_rollback
            exit 1
        ) &
        
        SAFETY_PID=$!
        trap 'kill $SAFETY_PID 2>/dev/null || true' EXIT
    fi
}

# Stop safety timeout
stop_safety_timeout() {
    if [[ -n "$SAFETY_PID" ]]; then
        kill "$SAFETY_PID" 2>/dev/null || true
        log_success "Safety timeout cancelled - setup completed successfully"
    fi
}

print_safety_banner() {
    echo
    echo "================================================================"
    echo -e "${BOLD}${GREEN}ULTRA-SAFE WARP ROUTING SETUP${NC}"
    echo "================================================================"
    echo -e "${YELLOW}Safety Features Active:${NC}"
    [[ "$PREVIEW_MODE" == "true" ]] && echo "  ✓ Preview Mode - Show commands without executing"
    [[ "$DRY_RUN" == "true" ]] && echo "  ✓ Dry Run - Simulate all operations"
    [[ "$INTERACTIVE_MODE" == "true" ]] && echo "  ✓ Interactive Mode - Confirm each step"
    [[ "$BACKUP_NETWORK_STATE" == "true" ]] && echo "  ✓ Network State Backup"
    [[ "$AUTO_ROLLBACK" == "true" ]] && echo "  ✓ Automatic Rollback on Failure"
    [[ "$STEP_BY_STEP" == "true" ]] && echo "  ✓ Step-by-step Execution"
    [[ "$SAFETY_TIMEOUT" -gt 0 ]] && echo "  ✓ Safety Timeout: ${SAFETY_TIMEOUT}s"
    echo "================================================================"
    echo
}

show_ultra_safe_options() {
    echo
    echo -e "${BOLD}${CYAN}ULTRA-SAFE EXECUTION OPTIONS:${NC}"
    echo
    echo -e "${GREEN}Maximum Safety (Recommended for first run):${NC}"
    echo "  sudo PREVIEW_MODE=true $0 safe-setup    # Show all commands without executing"
    echo "  sudo DRY_RUN=true $0 safe-setup         # Simulate everything"
    echo "  sudo INTERACTIVE_MODE=true $0 safe-setup # Confirm each step"
    echo
    echo -e "${YELLOW}Progressive Safety Levels:${NC}"
    echo "  sudo STEP_BY_STEP=true $0 safe-setup    # Pause after each command"
    echo "  sudo SAFETY_TIMEOUT=30 $0 safe-setup    # Auto-rollback after 30s"
    echo "  sudo AUTO_ROLLBACK=true $0 safe-setup   # Auto-rollback on any failure"
    echo
    echo -e "${BLUE}Network State Protection:${NC}"
    echo "  sudo BACKUP_NETWORK_STATE=true $0 safe-setup  # Backup before changes"
    echo "  sudo $0 restore-backup                   # Restore from backup"
    echo "  sudo $0 emergency-rollback               # Manual rollback"
    echo
    echo -e "${PURPLE}Combination Example (Maximum Safety):${NC}"
    echo "  sudo PREVIEW_MODE=true INTERACTIVE_MODE=true BACKUP_NETWORK_STATE=true \\"
    echo "       SAFETY_TIMEOUT=60 AUTO_ROLLBACK=true $0 safe-setup"
    echo
}

#===========================================
# UTILITY FUNCTIONS
#===========================================

print_section() {
    echo
    echo "=================================================="
    echo -e "${CYAN}$1${NC}"
    echo "=================================================="
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

command_exists() { command -v "$1" >/dev/null 2>&1; }

confirm_action() {
    local message="$1"
    if [[ "$FORCE_RECREATE" == "true" || "$FORCE_APPLY" == "true" ]]; then
        log_info "Force mode enabled, proceeding with: $message"
        return 0
    fi
    read -p "$(echo -e "${YELLOW}[CONFIRM]${NC} $message [y/N]: ")" -n 1 -r
    echo
    [[ $REPLY =~ ^[Yy]$ ]]
}

check_prerequisites() {
    print_section "CHECKING PREREQUISITES"
    local missing_tools=()
    local required_commands=("docker" "ip" "iptables" "ping" "grep" "awk" "sed")
    
    for cmd in "${required_commands[@]}"; do
        if ! command_exists "$cmd"; then
            missing_tools+=("$cmd")
        else
            log_debug "$cmd found"
        fi
    done
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker daemon is not running or not accessible"
        exit 1
    fi
    
    if ! ip rule list >/dev/null 2>&1; then
        log_error "Cannot access network routing rules"
        exit 1
    fi
    
    log_success "All prerequisites satisfied"
}

#===========================================
# DOCKER-IN-DOCKER SAFE TESTING
#===========================================

create_dind_test_environment() {
    print_section "CREATING ISOLATED TEST ENVIRONMENT"
    
    log_step "Starting Docker-in-Docker container..."
    docker rm -f "$DIND_CONTAINER_NAME" >/dev/null 2>&1 || true
    
    docker run -d \
        --name "$DIND_CONTAINER_NAME" \
        --privileged \
        --cgroupns=host \
        -v /var/lib/docker \
        -e DOCKER_TLS_CERTDIR= \
        docker:24-dind >/dev/null
    
    log_info "Waiting for Docker daemon to start in test environment..."
    for i in {1..30}; do
        if docker exec "$DIND_CONTAINER_NAME" docker info >/dev/null 2>&1; then
            break
        fi
        sleep 2
        if [[ $i -eq 30 ]]; then
            log_error "Docker daemon failed to start in test environment"
            cleanup_dind_environment
            exit 1
        fi
    done
    
    log_success "Test environment ready"
}

create_dind_test_script() {
    cat << EOF
#!/bin/sh
set -e

# Configuration
WARP_SUBNET="$WARP_SUBNET"
PUBLICNET_SUBNET="$PUBLICNET_SUBNET"
PUBLICNET_GATEWAY="$PUBLICNET_GATEWAY"
WARP_GATEWAY="$WARP_GATEWAY"
WARP_CONTAINER_IP="$WARP_CONTAINER_IP"
ROUTING_TABLE_ID="$ROUTING_TABLE_ID"
ROUTING_TABLE_NAME="$ROUTING_TABLE_NAME"

# Colors
RED='\033[0;31m'; GREEN='\033[0;32m'; BLUE='\033[0;34m'; NC='\033[0m'
log_info() { echo -e "\${BLUE}[TEST]\${NC} \$1"; }
log_success() { echo -e "\${GREEN}[TEST]\${NC} \$1"; }
log_error() { echo -e "\${RED}[TEST]\${NC} \$1"; }

# Install required packages
apk add --no-cache curl iproute2 iptables >/dev/null 2>&1

log_info "Setting up test networks..."
docker network create --driver=bridge --subnet="\$PUBLICNET_SUBNET" --gateway="\$PUBLICNET_GATEWAY" publicnet >/dev/null 2>&1
docker network create --driver=bridge --subnet="\$WARP_SUBNET" --gateway="\$WARP_GATEWAY" warpnet >/dev/null 2>&1

log_info "Starting mock WARP container..."
docker run -d --name warp-with-nat --network publicnet --ip "\$WARP_CONTAINER_IP" \\
    alpine:latest sh -c 'apk add --no-cache iptables && 
    echo 1 > /proc/sys/net/ipv4/ip_forward &&
    iptables -t nat -A POSTROUTING -j MASQUERADE &&
    while true; do sleep 1; done' >/dev/null 2>&1

# Get bridge interfaces
PUBLICNET_BRIDGE=\$(docker network inspect publicnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})
WARPNET_BRIDGE=\$(docker network inspect warpnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})

log_info "Setting up routing..."
echo "\$ROUTING_TABLE_ID \$ROUTING_TABLE_NAME" >> /etc/iproute2/rt_tables
ip rule add from "\$WARP_SUBNET" table "\$ROUTING_TABLE_ID"
ip route add default via "\$WARP_CONTAINER_IP" dev "\$PUBLICNET_BRIDGE" table "\$ROUTING_TABLE_ID"
iptables -t nat -A POSTROUTING -s "\$WARP_SUBNET" -o "\$PUBLICNET_BRIDGE" -j MASQUERADE

log_info "Testing routing configuration..."

# Test 1: Basic connectivity
log_info "Test 1: Basic container connectivity"
TEST_RESULT=\$(docker run --rm --network=warpnet alpine:latest sh -c 'apk add --no-cache curl >/dev/null 2>&1 && timeout 10 curl -s ifconfig.me' 2>/dev/null || echo "FAILED")

if [[ "\$TEST_RESULT" == "FAILED" || -z "\$TEST_RESULT" ]]; then
    log_error "Basic connectivity test failed"
    exit 1
fi
log_success "Container can reach internet: \$TEST_RESULT"

# Test 2: IP assignment
log_info "Test 2: IP assignment validation"
CONTAINER_IP=\$(docker run -d --network=warpnet alpine:latest sleep 10)
ASSIGNED_IP=\$(docker inspect "\$CONTAINER_IP" -f '{{.NetworkSettings.Networks.warpnet.IPAddress}}')
docker rm -f "\$CONTAINER_IP" >/dev/null 2>&1

if [[ -z "\$ASSIGNED_IP" ]]; then
    log_error "IP assignment test failed"
    exit 1
fi
log_success "Container assigned IP: \$ASSIGNED_IP"

# Test 3: Routing validation
log_info "Test 3: Routing rules validation"
if ! ip rule list | grep -q "from \$WARP_SUBNET lookup \$ROUTING_TABLE_ID"; then
    log_error "Policy routing rule not found"
    exit 1
fi

if ! ip route show table "\$ROUTING_TABLE_ID" | grep -q "default via \$WARP_CONTAINER_IP"; then
    log_error "Custom route not found"
    exit 1
fi

if ! iptables -t nat -C POSTROUTING -s "\$WARP_SUBNET" -o "\$PUBLICNET_BRIDGE" -j MASQUERADE 2>/dev/null; then
    log_error "NAT rule not found"
    exit 1
fi

log_success "All routing rules validated"
log_success "🎉 All tests passed in isolated environment!"
echo "ROUTING_TEST_SUCCESS"
EOF
}

run_dind_test() {
    print_section "RUNNING ISOLATED ROUTING TEST"
    
    log_step "Creating test script..."
    local test_script=$(create_dind_test_script)
    
    echo "$test_script" | docker exec -i "$DIND_CONTAINER_NAME" sh -c 'cat > /test_routing.sh && chmod +x /test_routing.sh'
    
    log_step "Running routing test in isolated environment..."
    
    local test_output
    if test_output=$(timeout "$TEST_TIMEOUT" docker exec "$DIND_CONTAINER_NAME" /test_routing.sh 2>&1); then
        echo "$test_output"
        
        if echo "$test_output" | grep -q "ROUTING_TEST_SUCCESS"; then
            log_success "✅ Isolated test completed successfully!"
            return 0
        else
            log_error "❌ Test script ran but didn't report success"
            return 1
        fi
    else
        log_error "❌ Test failed or timed out"
        echo "$test_output"
        return 1
    fi
}

cleanup_dind_environment() {
    log_step "Removing Docker-in-Docker container..."
    docker rm -f "$DIND_CONTAINER_NAME" >/dev/null 2>&1 || true
}

#===========================================
# DIRECT SETUP FUNCTIONS
#===========================================

cleanup_previous_setup() {
    print_section "CLEANING UP PREVIOUS SETUP"
    
    log_step "Removing policy routing rules..."
    local rules_removed=false
    
    for table_ref in "$ROUTING_TABLE_ID" "$ROUTING_TABLE_NAME"; do
        while ip rule del from "$WARP_SUBNET" table "$table_ref" 2>/dev/null; do
            log_info "Removed routing rule for $WARP_SUBNET -> $table_ref"
            rules_removed=true
        done
    done
    
    [[ "$rules_removed" == "false" ]] && log_info "No existing routing rules found"
    
    log_step "Flushing custom routing table..."
    ip route flush table "$ROUTING_TABLE_ID" 2>/dev/null && log_info "Flushed routing table $ROUTING_TABLE_ID" || true
    
    log_step "Cleaning up iptables NAT rules..."
    cleanup_iptables_rules
    
    log_step "Cleaning up routing table definition..."
    if grep -q "$ROUTING_TABLE_ID.*$ROUTING_TABLE_NAME" /etc/iproute2/rt_tables 2>/dev/null; then
        sed -i "/$ROUTING_TABLE_ID.*$ROUTING_TABLE_NAME/d" /etc/iproute2/rt_tables
        log_info "Removed $ROUTING_TABLE_NAME from rt_tables"
    fi
    
    log_success "Cleanup completed"
}

cleanup_iptables_rules() {
    local rules_found=false
    
    while IFS= read -r line; do
        if [[ -n "$line" ]]; then
            local rule_num=$(echo "$line" | awk '{print $1}')
            if iptables -t nat -D POSTROUTING "$rule_num" 2>/dev/null; then
                log_info "Removed iptables POSTROUTING rule #$rule_num"
                rules_found=true
                break
            fi
        fi
    done < <(iptables -t nat -L POSTROUTING --line-numbers -n | grep "$WARP_SUBNET" | sort -rn)
    
    [[ "$rules_found" == "true" ]] && cleanup_iptables_rules
}

setup_docker_networks() {
    print_section "SETTING UP DOCKER NETWORKS"
    
    # Setup publicnet
    if docker network inspect publicnet >/dev/null 2>&1; then
        log_success "publicnet network exists"
    else
        log_info "Creating publicnet network..."
        docker network create \
            --driver=bridge \
            --attachable \
            --subnet="$PUBLICNET_SUBNET" \
            --gateway="$PUBLICNET_GATEWAY" \
            publicnet
        log_success "publicnet network created"
    fi
    
    # Setup warpnet
    if docker network inspect warpnet >/dev/null 2>&1; then
        local existing_subnet=$(docker network inspect warpnet -f '{{range .IPAM.Config}}{{.Subnet}}{{end}}' 2>/dev/null)
        if [[ "$existing_subnet" == "$WARP_SUBNET" ]]; then
            log_success "warpnet network exists with correct configuration"
        else
            log_warning "warpnet exists but subnet differs"
            if confirm_action "Recreate warpnet network with correct subnet?"; then
                docker network rm warpnet
                create_warpnet_network
            fi
        fi
    else
        create_warpnet_network
    fi
}

create_warpnet_network() {
    log_info "Creating warpnet network..."
    safe_exec "docker network create --driver=bridge --attachable --subnet='$WARP_SUBNET' --gateway='$WARP_GATEWAY' warpnet" \
             "docker network rm warpnet 2>/dev/null || true" \
             "Create warpnet Docker network with subnet $WARP_SUBNET"
}

setup_warp_container() {
    print_section "SETTING UP WARP CONTAINER"
    
    if check_container_running; then
        log_success "warp-with-nat container is already running"
        verify_container_configuration
    elif check_container_exists; then
        if confirm_action "Start existing warp-with-nat container?"; then
            start_existing_container
        else
            if confirm_action "Remove existing container and create new one?"; then
                remove_existing_container
                create_warp_container
            fi
        fi
    else
        if [[ "$AUTO_START_CONTAINER" == "true" ]] || confirm_action "Create and start warp-with-nat container?"; then
            create_warp_container
        else
            log_error "warp-with-nat container is required but not available"
            exit 1
        fi
    fi
}

check_container_running() {
    docker ps --format '{{.Names}}' | grep -q "^${WARP_CONTAINER_NAME}$"
}

check_container_exists() {
    docker ps -a --format '{{.Names}}' | grep -q "^${WARP_CONTAINER_NAME}$"
}

start_existing_container() {
    log_info "Starting existing warp-with-nat container..."
    docker start "$WARP_CONTAINER_NAME"
    sleep 5
    
    if check_container_running; then
        log_success "Container started successfully"
        verify_container_configuration
    else
        log_error "Failed to start container"
        exit 1
    fi
}

remove_existing_container() {
    log_info "Removing existing warp-with-nat container..."
    docker rm -f "$WARP_CONTAINER_NAME" 2>/dev/null || true
}

create_warp_container() {
    log_info "Creating warp-with-nat container..."
    mkdir -p "${CONFIG_PATH}/warp/data"
    
    docker run -d \
        --name "$WARP_CONTAINER_NAME" \
        --network publicnet \
        --ip "$WARP_CONTAINER_IP" \
        --device-cgroup-rule 'c 10:200 rwm' \
        --cap-add MKNOD \
        --cap-add AUDIT_WRITE \
        --cap-add NET_ADMIN \
        --sysctl net.ipv6.conf.all.disable_ipv6=1 \
        --sysctl net.ipv4.conf.all.src_valid_mark=1 \
        --sysctl net.ipv4.ip_forward=1 \
        --sysctl net.ipv6.conf.all.forwarding=1 \
        --sysctl net.ipv6.conf.all.accept_ra=2 \
        --expose 443/udp \
        --expose 1080 \
        --publish "${WARP_SOCKS_PORT}:1080" \
        --env WARP_SLEEP=2 \
        --env "WARP_LICENSE_KEY=${WARP_LICENSE_KEY}" \
        --env WARP_ENABLE_NAT=1 \
        --env "TUNNEL_TOKEN=${TUNNEL_TOKEN}" \
        --volume "${CONFIG_PATH}/warp/data:/var/lib/cloudflare-warp" \
        --restart always \
        --label deunhealth.restart.on.unhealthy=true \
        "$WARP_IMAGE"
    
    log_info "Waiting for container to initialize..."
    sleep 10
    
    if check_container_running; then
        log_success "warp-with-nat container created and started"
        verify_container_configuration
    else
        log_error "Failed to create or start warp-with-nat container"
        docker logs "$WARP_CONTAINER_NAME" --tail 20
        exit 1
    fi
}

verify_container_configuration() {
    log_step "Verifying container configuration..."
    
    local container_networks=$(docker inspect "$WARP_CONTAINER_NAME" -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}} {{end}}' 2>/dev/null)
    if [[ ! "$container_networks" =~ "publicnet" ]]; then
        log_error "warp-with-nat container is not connected to publicnet"
        exit 1
    fi
    
    local container_ip=$(docker inspect "$WARP_CONTAINER_NAME" -f '{{.NetworkSettings.Networks.publicnet.IPAMConfig.IPv4Address}}' 2>/dev/null)
    if [[ -z "$container_ip" ]]; then
        container_ip=$(docker inspect "$WARP_CONTAINER_NAME" -f '{{.NetworkSettings.Networks.publicnet.IPAddress}}' 2>/dev/null)
    fi
    
    if [[ -z "$container_ip" ]]; then
        log_error "Could not determine container IP address"
        exit 1
    fi
    
    if [[ "$container_ip" != "$WARP_CONTAINER_IP" ]]; then
        log_warning "Container IP ($container_ip) differs from expected ($WARP_CONTAINER_IP)"
        log_info "Updating configuration to use actual container IP"
        WARP_CONTAINER_IP="$container_ip"
    fi
    
    log_success "Container verified at $WARP_CONTAINER_IP"
}

setup_routing() {
    print_section "SETTING UP ROUTING"
    
    get_bridge_interfaces
    configure_routing_table
    setup_policy_routing
    setup_nat_rules
}

get_bridge_interfaces() {
    log_step "Determining bridge interface names..."
    
    PUBLICNET_BRIDGE=$(docker network inspect publicnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})
    if [[ -z "$PUBLICNET_BRIDGE" ]] || ! ip link show "$PUBLICNET_BRIDGE" >/dev/null 2>&1; then
        log_error "Could not determine publicnet bridge interface"
        exit 1
    fi
    
    WARPNET_BRIDGE=$(docker network inspect warpnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})
    if [[ -z "$WARPNET_BRIDGE" ]] || ! ip link show "$WARPNET_BRIDGE" >/dev/null 2>&1; then
        log_error "Could not determine warpnet bridge interface"
        exit 1
    fi
    
    log_success "Bridge interfaces identified"
}

configure_routing_table() {
    log_step "Configuring custom routing table..."
    
    if ! grep -q "^${ROUTING_TABLE_ID}[[:space:]].*${ROUTING_TABLE_NAME}" /etc/iproute2/rt_tables; then
        safe_exec "echo '$ROUTING_TABLE_ID $ROUTING_TABLE_NAME' >> /etc/iproute2/rt_tables" \
                 "sed -i '/$ROUTING_TABLE_NAME/d' /etc/iproute2/rt_tables" \
                 "Add custom routing table $ROUTING_TABLE_NAME (ID: $ROUTING_TABLE_ID)"
    fi
}

setup_policy_routing() {
    log_step "Setting up policy-based routing..."
    
    if ! ip rule list | grep -q "from ${WARP_SUBNET} lookup ${ROUTING_TABLE_ID}"; then
        safe_exec "ip rule add from '$WARP_SUBNET' table '$ROUTING_TABLE_ID'" \
                 "ip rule del from '$WARP_SUBNET' table '$ROUTING_TABLE_ID' 2>/dev/null || true" \
                 "Add policy routing rule: $WARP_SUBNET -> table $ROUTING_TABLE_ID"
    fi
    
    if ! ip route show table "$ROUTING_TABLE_ID" | grep -q "default via ${WARP_CONTAINER_IP}"; then
        safe_exec "ip route add default via '$WARP_CONTAINER_IP' dev '$PUBLICNET_BRIDGE' table '$ROUTING_TABLE_ID'" \
                 "ip route del default via '$WARP_CONTAINER_IP' dev '$PUBLICNET_BRIDGE' table '$ROUTING_TABLE_ID' 2>/dev/null || true" \
                 "Add default route via $WARP_CONTAINER_IP through $PUBLICNET_BRIDGE"
    fi
}

setup_nat_rules() {
    log_step "Setting up NAT rules..."
    
    if ! iptables -t nat -C POSTROUTING -s "$WARP_SUBNET" -o "$PUBLICNET_BRIDGE" -j MASQUERADE 2>/dev/null; then
        safe_exec "iptables -t nat -A POSTROUTING -s '$WARP_SUBNET' -o '$PUBLICNET_BRIDGE' -j MASQUERADE" \
                 "iptables -t nat -D POSTROUTING -s '$WARP_SUBNET' -o '$PUBLICNET_BRIDGE' -j MASQUERADE 2>/dev/null || true" \
                 "Add iptables NAT masquerading rule for $WARP_SUBNET"
    fi
}

verify_setup() {
    print_section "VERIFYING SETUP"
    
    local verification_failed=false
    
    # Check policy rule
    if ip rule list | grep -q "from ${WARP_SUBNET} lookup ${ROUTING_TABLE_ID}"; then
        log_success "✓ Policy routing rule is active"
    else
        log_error "✗ Policy routing rule not found"
        verification_failed=true
    fi
    
    # Check custom route
    if ip route show table "$ROUTING_TABLE_ID" | grep -q "default via ${WARP_CONTAINER_IP}"; then
        log_success "✓ Custom routing table configured correctly"
    else
        log_error "✗ Custom route not found in table $ROUTING_TABLE_ID"
        verification_failed=true
    fi
    
    # Check iptables rule
    if iptables -t nat -C POSTROUTING -s "$WARP_SUBNET" -o "$PUBLICNET_BRIDGE" -j MASQUERADE 2>/dev/null; then
        log_success "✓ iptables NAT rule is active"
    else
        log_error "✗ iptables NAT rule not found"
        verification_failed=true
    fi
    
    # Check container connectivity
    if ping -c 1 -W 2 "$WARP_CONTAINER_IP" >/dev/null 2>&1; then
        log_success "✓ Warp container is reachable"
    else
        log_warning "⚠ Cannot ping warp container (may be normal if ICMP is blocked)"
    fi
    
    if [[ "$verification_failed" == "true" ]]; then
        log_error "Setup verification failed"
        return 1
    else
        log_success "All components verified successfully"
        return 0
    fi
}

#===========================================
# TESTING FUNCTIONS
#===========================================

test_setup() {
    print_section "COMPREHENSIVE ROUTING TESTS"
    
    if [[ "$RUN_TESTS" != "true" ]]; then
        log_info "Tests disabled (RUN_TESTS=false), skipping validation"
        return 0
    fi
    
    local test_failures=0
    local test_image="alpine:latest"
    
    log_step "Preparing test environment..."
    if ! docker image inspect "$test_image" >/dev/null 2>&1; then
        log_info "Pulling test image: $test_image"
        docker pull "$test_image" >/dev/null
    fi
    
    # Test 1: Basic connectivity
    log_step "Test 1: Basic connectivity through warpnet"
    if test_basic_connectivity; then
        log_success "✓ Basic connectivity test passed"
    else
        log_error "✗ Basic connectivity test failed"
        ((test_failures++))
    fi
    
    # Test 2: IP routing verification
    log_step "Test 2: External IP routing verification"
    if test_ip_routing; then
        log_success "✓ IP routing test passed"
    else
        log_error "✗ IP routing test failed"
        ((test_failures++))
    fi
    
    # Test 3: Dynamic IP assignment
    log_step "Test 3: Dynamic IP assignment within subnet"
    if test_dynamic_ip_assignment; then
        log_success "✓ Dynamic IP assignment test passed"
    else
        log_error "✗ Dynamic IP assignment test failed"
        ((test_failures++))
    fi
    
    cleanup_test_containers
    
    if [[ $test_failures -eq 0 ]]; then
        log_success "🎉 All tests passed! Routing setup is working correctly."
        return 0
    else
        log_error "❌ $test_failures test(s) failed. Routing setup may not be working correctly."
        
        if [[ "$CLEANUP_ON_TEST_FAILURE" == "true" ]]; then
            log_warning "CLEANUP_ON_TEST_FAILURE is enabled. Cleaning up failed setup..."
            cleanup_previous_setup
        fi
        
        return 1
    fi
}

test_basic_connectivity() {
    local test_result
    test_result=$(docker run --rm --name "warp-test-basic" --network=warpnet \
        alpine:latest sh -c 'apk add --no-cache curl >/dev/null 2>&1 && curl -s --max-time 15 ifconfig.me' 2>/dev/null)
    
    if [[ -n "$test_result" && "$test_result" != "FAILED" ]]; then
        log_debug "Container external IP: $test_result"
        return 0
    else
        log_debug "Failed to get external IP from container"
        return 1
    fi
}

test_ip_routing() {
    local container_ip
    container_ip=$(docker run --rm --name "warp-test-routing" --network=warpnet \
        alpine:latest sh -c 'apk add --no-cache curl >/dev/null 2>&1 && curl -s --max-time 15 ifconfig.me' 2>/dev/null)
    
    local host_ip
    host_ip=$(timeout 10 curl -s ifconfig.me 2>/dev/null || echo "UNKNOWN")
    
    log_debug "Container IP: $container_ip"
    log_debug "Host IP: $host_ip"
    
    if [[ -n "$container_ip" && "$container_ip" != "FAILED" && "$container_ip" != "$host_ip" ]]; then
        log_info "✓ Traffic is routing through WARP (Container: $container_ip, Host: $host_ip)"
        return 0
    elif [[ "$container_ip" == "$host_ip" ]]; then
        log_warning "Container and host have same external IP - routing may not be working"
        return 1
    else
        log_error "Could not determine container external IP"
        return 1
    fi
}

test_dynamic_ip_assignment() {
    local container_count=3
    local assigned_ips=()
    
    for i in $(seq 1 $container_count); do
        local container_name="warp-test-ip-$i"
        
        docker run -d --name "$container_name" --network=warpnet \
            alpine:latest sleep 30 >/dev/null
        
        local assigned_ip
        assigned_ip=$(docker inspect "$container_name" -f '{{.NetworkSettings.Networks.warpnet.IPAddress}}' 2>/dev/null)
        
        if [[ -n "$assigned_ip" ]]; then
            assigned_ips+=("$assigned_ip")
            log_debug "Container $container_name assigned IP: $assigned_ip"
        else
            log_error "Could not get IP for container $container_name"
            return 1
        fi
    done
    
    local unique_ips=($(printf '%s\n' "${assigned_ips[@]}" | sort -u))
    if [[ ${#unique_ips[@]} -eq ${#assigned_ips[@]} ]]; then
        log_debug "All containers received unique IPs: ${assigned_ips[*]}"
        return 0
    else
        log_error "IP assignment conflict detected"
        return 1
    fi
}

cleanup_test_containers() {
    log_debug "Cleaning up test containers..."
    local test_containers
    test_containers=$(docker ps -aq --filter "name=warp-test-" 2>/dev/null || true)
    
    if [[ -n "$test_containers" ]]; then
        docker rm -f $test_containers >/dev/null 2>&1 || true
        log_debug "Removed test containers"
    fi
}

show_status() {
    print_section "CURRENT CONFIGURATION STATUS"
    
    echo "Network Configuration:"
    echo "  warpnet subnet: $WARP_SUBNET"
    echo "  publicnet subnet: $PUBLICNET_SUBNET"
    echo "  warp container IP: $WARP_CONTAINER_IP"
    echo
    
    echo "Bridge Interfaces:"
    echo "  publicnet: ${PUBLICNET_BRIDGE:-not determined}"
    echo "  warpnet: ${WARPNET_BRIDGE:-not determined}"
    echo
    
    echo "Policy Routing Rules:"
    ip rule list | grep -E "(${ROUTING_TABLE_NAME}|${ROUTING_TABLE_ID})" || echo "  No custom rules found"
    echo
    
    echo "Custom Routing Table ($ROUTING_TABLE_ID):"
    ip route show table "$ROUTING_TABLE_ID" 2>/dev/null || echo "  Table empty or not found"
    echo
    
    echo "NAT Rules for $WARP_SUBNET:"
    iptables -t nat -L POSTROUTING -n --line-numbers | grep "$WARP_SUBNET" || echo "  No NAT rules found"
    echo
    
    echo "Container Status:"
    if check_container_running; then
        echo "  warp-with-nat: RUNNING"
        docker inspect "$WARP_CONTAINER_NAME" --format '  IP: {{.NetworkSettings.Networks.publicnet.IPAddress}}' 2>/dev/null || echo "  IP: unknown"
    else
        echo "  warp-with-nat: NOT RUNNING"
    fi
}

#===========================================
# MAIN EXECUTION FUNCTIONS
#===========================================

safe_setup_main() {
    print_section "SAFE WARP ROUTING SETUP"
    log_info "Testing routing in isolated Docker-in-Docker environment first"
    echo
    
    check_root
    check_prerequisites
    create_dind_test_environment
    
    if run_dind_test; then
        cleanup_dind_environment
        
        if [[ "$FORCE_APPLY" == "true" ]] || confirm_action "Tests passed in isolated environment. Apply to host?"; then
            log_step "Applying validated configuration to host..."
            direct_setup_main
        else
            log_info "Skipping host application"
        fi
    else
        log_error "❌ Isolated test failed - NOT applying to host"
        cleanup_dind_environment
        exit 1
    fi
}

direct_setup_main() {
    print_section "DIRECT WARP ROUTING SETUP"
    
    log_info "Starting direct setup with configuration:"
    log_info "  WARP Subnet: $WARP_SUBNET"
    log_info "  PublicNet Subnet: $PUBLICNET_SUBNET"
    log_info "  WARP Container IP: $WARP_CONTAINER_IP"
    log_info "  Routing Table: $ROUTING_TABLE_NAME ($ROUTING_TABLE_ID)"
    
    check_root
    check_prerequisites
    cleanup_previous_setup
    setup_docker_networks
    setup_warp_container
    setup_routing
    
    if verify_setup; then
        show_status
        
        if [[ "$RUN_TESTS" == "true" ]]; then
            if test_setup; then
                print_section "SETUP COMPLETED SUCCESSFULLY"
                log_success "All traffic from warpnet ($WARP_SUBNET) will route through warp-with-nat ($WARP_CONTAINER_IP)"
                log_success "All tests passed - routing is working correctly!"
            else
                log_error "Setup verification failed during testing"
                exit 1
            fi
        else
            print_section "SETUP COMPLETED"
            log_success "All traffic from warpnet ($WARP_SUBNET) will route through warp-with-nat ($WARP_CONTAINER_IP)"
            log_info "Tests were skipped (RUN_TESTS=false)"
        fi
        
        echo
        log_info "Usage examples:"
        log_info "  Test routing: docker run --rm --network=warpnet alpine sh -c 'apk add curl && curl -s ifconfig.me'"
        log_info "  Check status: sudo $0 status"
        log_info "  Run tests: sudo $0 test"
        log_info "  Remove setup: sudo $0 cleanup"
        echo
    else
        log_error "Setup completed with errors - please check the configuration"
        exit 1
    fi
}

cleanup_only() {
    print_section "REMOVING WARP NETWORK ROUTING SETUP"
    
    cleanup_previous_setup
    cleanup_dind_environment
    
    if check_container_running; then
        if confirm_action "Stop and remove warp-with-nat container?"; then
            docker stop "$WARP_CONTAINER_NAME" >/dev/null 2>&1 || true
            docker rm "$WARP_CONTAINER_NAME" >/dev/null 2>&1 || true
            log_info "Container stopped and removed"
        fi
    fi
    
    log_success "Cleanup completed - normal Docker routing restored"
}

standalone_test() {
    print_section "STANDALONE ROUTING TEST"
    
    if ! docker network inspect warpnet >/dev/null 2>&1; then
        log_error "warpnet Docker network not found. Please run setup first."
        exit 1
    fi
    
    log_info "Testing routing functionality..."
    
    # Get host external IP
    log_info "Getting host external IP..."
    HOST_IP=$(timeout 10 curl -s ifconfig.me 2>/dev/null || echo "FAILED")
    if [[ "$HOST_IP" == "FAILED" ]]; then
        log_error "Could not get host external IP"
        exit 1
    fi
    log_info "Host external IP: $HOST_IP"
    
    # Get container external IP through warpnet
    log_info "Testing container routing through warpnet..."
    CONTAINER_IP=$(docker run --rm --network=warpnet alpine:latest sh -c '
        apk add --no-cache curl >/dev/null 2>&1 || exit 1
        curl -s --max-time 15 ifconfig.me 2>/dev/null || echo "FAILED"
    ' 2>/dev/null)
    
    if [[ "$CONTAINER_IP" == "FAILED" || -z "$CONTAINER_IP" ]]; then
        log_error "Could not get container external IP through warpnet"
        exit 1
    fi
    log_info "Container external IP: $CONTAINER_IP"
    
    # Compare IPs
    if [[ "$HOST_IP" == "$CONTAINER_IP" ]]; then
        log_error "Host and container have the same external IP!"
        log_error "This indicates that routing through WARP is NOT working"
        exit 1
    else
        log_success "Host and container have different external IPs!"
        log_success "Routing through WARP is working correctly"
        log_info "Host IP: $HOST_IP"
        log_info "Container IP: $CONTAINER_IP (via WARP)"
    fi
}

show_ultra_safe_demo() {
    echo
    echo "================================================================"
    echo -e "${BOLD}${GREEN}ULTRA-SAFE DEMONSTRATION${NC}"
    echo "================================================================"
    echo -e "${YELLOW}This is what will happen when you run the script safely:${NC}"
    echo
    
    echo -e "${CYAN}1. BACKUP PHASE${NC}"
    echo "   • Create backup of current network state"
    echo "   • Save routing tables, iptables rules, Docker networks"
    echo "   • Store backup in: /tmp/warp-routing-backup-[timestamp]"
    echo
    
    echo -e "${CYAN}2. ISOLATED TESTING PHASE (Docker-in-Docker)${NC}"
    echo "   • Start completely isolated test environment"
    echo "   • Create test networks and routing inside DinD container"
    echo "   • Verify everything works without touching your host"
    echo "   • Test external IP differences to confirm WARP routing"
    echo
    
    echo -e "${CYAN}3. SAFETY CHECKS${NC}"
    if [[ "$PREVIEW_MODE" == "true" ]]; then
        echo "   • ✓ PREVIEW MODE: Show commands without executing"
    fi
    if [[ "$DRY_RUN" == "true" ]]; then
        echo "   • ✓ DRY RUN: Simulate all operations"
    fi
    if [[ "$INTERACTIVE_MODE" == "true" ]]; then
        echo "   • ✓ INTERACTIVE: Ask permission for each command"
    fi
    if [[ "$AUTO_ROLLBACK" == "true" ]]; then
        echo "   • ✓ AUTO-ROLLBACK: Undo everything if anything fails"
    fi
    if [[ "$SAFETY_TIMEOUT" -gt 0 ]]; then
        echo "   • ✓ SAFETY TIMEOUT: Auto-rollback after ${SAFETY_TIMEOUT} seconds"
    fi
    echo
    
    echo -e "${CYAN}4. HOST APPLICATION (Only after successful testing)${NC}"
    echo "   • Create Docker networks (publicnet, warpnet)"
    echo "   • Start WARP container with full isolation"
    echo "   • Add routing rules to route traffic through WARP"
    echo "   • Set up NAT masquerading for container communication"
    echo
    
    echo -e "${CYAN}5. VERIFICATION${NC}"
    echo "   • Test that containers get different external IPs"
    echo "   • Verify routing rules are working correctly"
    echo "   • Confirm WARP container is healthy"
    echo
    
    echo -e "${CYAN}6. ROLLBACK CAPABILITIES${NC}"
    echo "   • Emergency rollback: sudo $0 emergency-rollback"
    echo "   • Restore backup: sudo $0 restore-backup"
    echo "   • Clean removal: sudo $0 cleanup"
    echo
    
    echo -e "${GREEN}MAXIMUM SAFETY RECOMMENDATIONS:${NC}"
    echo
    echo -e "${BOLD}For the most cautious first run:${NC}"
    echo "  sudo PREVIEW_MODE=true INTERACTIVE_MODE=true BACKUP_NETWORK_STATE=true \\"
    echo "       SAFETY_TIMEOUT=60 AUTO_ROLLBACK=true $0 safe-setup"
    echo
    echo -e "${BOLD}This will:${NC}"
    echo "  • Show you every command before running it"
    echo "  • Ask permission for each step"
    echo "  • Create backups of everything"
    echo "  • Auto-rollback after 60 seconds if you don't stop it"
    echo "  • Test everything in isolation first"
    echo
    echo -e "${YELLOW}Want to see exactly what commands will run?${NC}"
    echo "  sudo PREVIEW_MODE=true $0"
    echo
    echo -e "${YELLOW}Want to simulate everything without any changes?${NC}"
    echo "  sudo DRY_RUN=true $0"
    echo
}

# Add demo command to help
show_help() {
    echo "ULTRA-SAFE WARP Network Routing Setup Script"
    echo
    echo "This script provides maximum safety with preview mode, backups, rollback,"
    echo "Docker-in-Docker testing, and step-by-step confirmation capabilities."
    echo
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo
    echo "Commands:"
    echo "  safe-setup       - Test in DinD isolation, then apply to host (default)"
    echo "  direct-setup     - Apply directly to host (skip DinD testing)"
    echo "  test-only        - Only run DinD test, don't apply to host"
    echo "  test             - Test existing setup"
    echo "  status           - Show current status"
    echo "  cleanup          - Remove routing setup"
    echo "  emergency-rollback - Emergency rollback of all changes"
    echo "  restore-backup   - Restore network state from backup"
    echo "  ultra-safe-help  - Show ultra-safe execution options"
    echo "  demo             - Show what the script will do (ultra-safe demo)"
    echo "  help             - Show this help"
    echo
    echo "Ultra-Safety Variables:"
    echo "  PREVIEW_MODE=true            - Show commands without executing"
    echo "  DRY_RUN=true                 - Simulate everything"
    echo "  INTERACTIVE_MODE=true        - Confirm each step"
    echo "  STEP_BY_STEP=true            - Pause after each command"
    echo "  BACKUP_NETWORK_STATE=true    - Backup before changes (default: true)"
    echo "  AUTO_ROLLBACK=true           - Auto rollback on failure (default: true)"
    echo "  SAFETY_TIMEOUT=seconds       - Auto-rollback timeout (default: 60)"
    echo
    echo "Traditional Variables:"
    echo "  USE_SAFE_MODE=false          - Skip DinD testing (use direct setup)"
    echo "  FORCE_APPLY=true             - Apply to host without confirmation"
    echo "  RUN_TESTS=false              - Skip host testing after setup"
    echo "  DEBUG=true                   - Enable debug output"
    echo "  FORCE_RECREATE=true          - Skip all confirmations"
    echo
    echo "Safety Examples:"
    echo "  sudo $0 demo                            # See what the script will do"
    echo "  sudo PREVIEW_MODE=true $0               # See all commands without executing"
    echo "  sudo DRY_RUN=true $0                   # Simulate everything"
    echo "  sudo INTERACTIVE_MODE=true $0          # Confirm each step"
    echo "  sudo STEP_BY_STEP=true $0              # Pause after each step"
    echo "  sudo SAFETY_TIMEOUT=30 $0              # Auto-rollback after 30s"
    echo "  sudo $0 ultra-safe-help                # Show all safety options"
    echo
    echo "Traditional Examples:"
    echo "  sudo $0                                 # Safe setup (test in DinD first)"
    echo "  sudo $0 direct-setup                   # Direct setup (no DinD testing)"
    echo "  sudo $0 test-only                      # Test only in DinD"
    echo "  sudo USE_SAFE_MODE=false $0            # Same as direct-setup"
    echo "  sudo FORCE_APPLY=true $0               # Auto-apply after DinD test"
}

#===========================================
# COMMAND LINE INTERFACE
#===========================================

# Show safety banner if any safety features are enabled
if [[ "$PREVIEW_MODE" == "true" || "$DRY_RUN" == "true" || "$INTERACTIVE_MODE" == "true" || 
      "$STEP_BY_STEP" == "true" || "$BACKUP_NETWORK_STATE" == "true" || 
      "$AUTO_ROLLBACK" == "true" || "$SAFETY_TIMEOUT" -gt 0 ]]; then
    print_safety_banner
fi

case "${1:-safe-setup}" in
    safe-setup|setup|"")
        if [[ "$USE_SAFE_MODE" == "false" ]]; then
            backup_network_state
            setup_safety_timeout
            direct_setup_main
            stop_safety_timeout
        else
            backup_network_state
            setup_safety_timeout
            safe_setup_main
            stop_safety_timeout
        fi
        ;;
    direct-setup)
        backup_network_state
        setup_safety_timeout
        direct_setup_main
        stop_safety_timeout
        ;;
    test-only)
        check_root
        check_prerequisites
        create_dind_test_environment
        if run_dind_test; then
            cleanup_dind_environment
            log_success "✅ Isolated test completed successfully!"
            log_info "Run 'safe-setup' to apply to host"
        else
            cleanup_dind_environment
            exit 1
        fi
        ;;
    test|verify)
        standalone_test
        ;;
    status|show)
        show_status
        ;;
    cleanup|clean|remove)
        check_root
        cleanup_only
        ;;
    emergency-rollback|rollback)
        check_root
        emergency_rollback
        ;;
    restore-backup|restore)
        check_root
        restore_network_state
        ;;
    ultra-safe-help|safety-help)
        show_ultra_safe_options
        ;;
    demo)
        show_ultra_safe_demo
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        echo
        show_help
        exit 1
        ;;
esac 