#!/bin/bash

set -e  # Exit on any error

#===========================================
# CONFIGURATION SECTION
#===========================================

# Network Configuration
WARP_SUBNET=${WARP_SUBNET:-172.20.0.0/16}
WARP_GATEWAY=${WARP_GATEWAY_IPV4_ADDRESS:-172.20.0.1}
PUBLICNET_SUBNET=${PUBLICNET_SUBNET:-10.76.0.0/16}
PUBLICNET_GATEWAY=${PUBLICNET_GATEWAY:-10.76.0.1}
WARP_CONTAINER_IP=${WARP_IPV4_ADDRESS:-10.76.128.200}

# Routing Configuration
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
AUTO_START_CONTAINER=${AUTO_START_CONTAINER:-true}
RUN_TESTS=${RUN_TESTS:-true}
CLEANUP_ON_TEST_FAILURE=${CLEANUP_ON_TEST_FAILURE:-true}

#===========================================
# UTILITY FUNCTIONS
#===========================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Logging functions
log_debug() {
    [[ "$DEBUG" == "true" ]] && echo -e "${PURPLE}[DEBUG]${NC} $1"
}

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

log_step() {
    echo -e "${CYAN}[STEP]${NC} $1"
}

# Function to print section headers
print_section() {
    echo
    echo "=================================================="
    echo -e "${CYAN}$1${NC}"
    echo "=================================================="
}

# Function to check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to wait for user confirmation
confirm_action() {
    local message="$1"
    local default="${2:-n}"
    
    if [[ "$FORCE_RECREATE" == "true" ]]; then
        log_info "Force mode enabled, proceeding with: $message"
        return 0
    fi
    
    read -p "$(echo -e "${YELLOW}[CONFIRM]${NC} $message [y/N]: ")" -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        return 0
    else
        return 1
    fi
}

#===========================================
# PREREQUISITE CHECKS
#===========================================

check_prerequisites() {
    print_section "CHECKING PREREQUISITES"
    
    local missing_tools=()
    
    # Check required commands
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
        log_info "Please install the missing tools and try again"
        exit 1
    fi
    
    # Check Docker daemon
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker daemon is not running or not accessible"
        exit 1
    fi
    
    # Check if we can modify network rules
    if ! ip rule list >/dev/null 2>&1; then
        log_error "Cannot access network routing rules"
        exit 1
    fi
    
    log_success "All prerequisites satisfied"
}

#===========================================
# CLEANUP FUNCTIONS
#===========================================

cleanup_previous_setup() {
    print_section "CLEANING UP PREVIOUS SETUP"
    
    # Remove policy routing rules
    log_step "Removing policy routing rules..."
    local rules_removed=false
    
    # Try both numeric and named table references
    for table_ref in "$ROUTING_TABLE_ID" "$ROUTING_TABLE_NAME"; do
        while ip rule del from "$WARP_SUBNET" table "$table_ref" 2>/dev/null; do
            log_info "Removed routing rule for $WARP_SUBNET -> $table_ref"
            rules_removed=true
        done
    done
    
    if [[ "$rules_removed" == "false" ]]; then
        log_info "No existing routing rules found"
    fi
    
    # Flush the custom routing table
    log_step "Flushing custom routing table..."
    if ip route flush table "$ROUTING_TABLE_ID" 2>/dev/null; then
        log_info "Flushed routing table $ROUTING_TABLE_ID"
    fi
    
    # Clean up iptables rules
    log_step "Cleaning up iptables NAT rules..."
    cleanup_iptables_rules
    
    # Remove routing table definition
    log_step "Cleaning up routing table definition..."
    if grep -q "$ROUTING_TABLE_ID.*$ROUTING_TABLE_NAME" /etc/iproute2/rt_tables 2>/dev/null; then
        sed -i "/$ROUTING_TABLE_ID.*$ROUTING_TABLE_NAME/d" /etc/iproute2/rt_tables
        log_info "Removed $ROUTING_TABLE_NAME from rt_tables"
    fi
    
    log_success "Cleanup completed"
}

cleanup_iptables_rules() {
    local rules_found=false
    
    # Find and remove POSTROUTING rules for our subnet
    while IFS= read -r line; do
        if [[ -n "$line" ]]; then
            local rule_num=$(echo "$line" | awk '{print $1}')
            if iptables -t nat -D POSTROUTING "$rule_num" 2>/dev/null; then
                log_info "Removed iptables POSTROUTING rule #$rule_num"
                rules_found=true
                # Break after removing one rule as line numbers shift
                break
            fi
        fi
    done < <(iptables -t nat -L POSTROUTING --line-numbers -n | grep "$WARP_SUBNET" | sort -rn)
    
    # Repeat until no more rules are found
    if [[ "$rules_found" == "true" ]]; then
        cleanup_iptables_rules  # Recursive call to handle shifted line numbers
    fi
}

#===========================================
# DOCKER NETWORK FUNCTIONS
#===========================================

setup_docker_networks() {
    print_section "SETTING UP DOCKER NETWORKS"
    
    setup_publicnet_network
    setup_warpnet_network
}

setup_publicnet_network() {
    log_step "Verifying publicnet network..."
    
    if docker network inspect publicnet >/dev/null 2>&1; then
        log_success "publicnet network exists"
        
        # Verify configuration
        local existing_subnet=$(docker network inspect publicnet -f '{{range .IPAM.Config}}{{.Subnet}}{{end}}' 2>/dev/null)
        if [[ -n "$existing_subnet" && "$existing_subnet" != "$PUBLICNET_SUBNET" ]]; then
            log_warning "publicnet subnet ($existing_subnet) differs from expected ($PUBLICNET_SUBNET)"
        fi
    else
        log_info "Creating publicnet network..."
        docker network create \
            --driver=bridge \
            --attachable \
            --subnet="$PUBLICNET_SUBNET" \
            --gateway="$PUBLICNET_GATEWAY" \
            --opt com.docker.network.bridge.name=publicnet \
            publicnet
        log_success "publicnet network created"
    fi
}

setup_warpnet_network() {
    log_step "Setting up warpnet network..."
    
    if docker network inspect warpnet >/dev/null 2>&1; then
        local existing_subnet=$(docker network inspect warpnet -f '{{range .IPAM.Config}}{{.Subnet}}{{end}}' 2>/dev/null)
        
        if [[ "$existing_subnet" == "$WARP_SUBNET" ]]; then
            log_success "warpnet network exists with correct configuration"
        else
            log_warning "warpnet exists but subnet ($existing_subnet) differs from expected ($WARP_SUBNET)"
            
            if confirm_action "Recreate warpnet network with correct subnet?"; then
                docker network rm warpnet
                create_warpnet_network
            else
                log_warning "Continuing with existing warpnet configuration"
            fi
        fi
    else
        create_warpnet_network
    fi
}

create_warpnet_network() {
    log_info "Creating warpnet network..."
    docker network create \
        --driver=bridge \
        --attachable \
        --subnet="$WARP_SUBNET" \
        --gateway="$WARP_GATEWAY" \
        warpnet
    log_success "warpnet network created ($WARP_SUBNET)"
}

#===========================================
# CONTAINER MANAGEMENT FUNCTIONS
#===========================================

setup_warp_container() {
    print_section "SETTING UP WARP CONTAINER"
    
    if check_container_running; then
        log_success "warp-with-nat container is already running"
        verify_container_configuration
    elif check_container_exists; then
        log_info "warp-with-nat container exists but is not running"
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
    
    # Wait for container to be ready
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
    
    # Ensure config directory exists
    mkdir -p "${CONFIG_PATH}/warp/data"
    
    # Create the container with all necessary configuration
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
    
    # Wait for container to start and initialize
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
    
    # Check if container is on publicnet
    local container_networks=$(docker inspect "$WARP_CONTAINER_NAME" -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}} {{end}}' 2>/dev/null)
    if [[ ! "$container_networks" =~ "publicnet" ]]; then
        log_error "warp-with-nat container is not connected to publicnet"
        exit 1
    fi
    
    # Get and verify container IP
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

#===========================================
# ROUTING SETUP FUNCTIONS
#===========================================

setup_routing() {
    print_section "SETTING UP ROUTING"
    
    get_bridge_interfaces
    configure_routing_table
    setup_policy_routing
    setup_nat_rules
}

get_bridge_interfaces() {
    log_step "Determining bridge interface names..."
    
    # Get publicnet bridge interface
    PUBLICNET_BRIDGE=$(docker network inspect publicnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})
    if [[ -z "$PUBLICNET_BRIDGE" ]] || ! ip link show "$PUBLICNET_BRIDGE" >/dev/null 2>&1; then
        log_error "Could not determine publicnet bridge interface"
        exit 1
    fi
    log_debug "publicnet bridge: $PUBLICNET_BRIDGE"
    
    # Get warpnet bridge interface
    WARPNET_BRIDGE=$(docker network inspect warpnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})
    if [[ -z "$WARPNET_BRIDGE" ]] || ! ip link show "$WARPNET_BRIDGE" >/dev/null 2>&1; then
        log_error "Could not determine warpnet bridge interface"
        exit 1
    fi
    log_debug "warpnet bridge: $WARPNET_BRIDGE"
    
    log_success "Bridge interfaces identified"
}

configure_routing_table() {
    log_step "Configuring custom routing table..."
    
    # Add routing table to rt_tables if not present
    if ! grep -q "^${ROUTING_TABLE_ID}[[:space:]].*${ROUTING_TABLE_NAME}" /etc/iproute2/rt_tables; then
        echo "$ROUTING_TABLE_ID $ROUTING_TABLE_NAME" >> /etc/iproute2/rt_tables
        log_info "Added routing table $ROUTING_TABLE_NAME (ID: $ROUTING_TABLE_ID)"
    else
        log_debug "Routing table already exists in rt_tables"
    fi
}

setup_policy_routing() {
    log_step "Setting up policy-based routing..."
    
    # Add policy routing rule
    if ! ip rule list | grep -q "from ${WARP_SUBNET} lookup ${ROUTING_TABLE_ID}"; then
        ip rule add from "$WARP_SUBNET" table "$ROUTING_TABLE_ID"
        log_info "Added policy routing rule: $WARP_SUBNET -> table $ROUTING_TABLE_ID"
    else
        log_debug "Policy routing rule already exists"
    fi
    
    # Add route in custom table
    if ! ip route show table "$ROUTING_TABLE_ID" | grep -q "default via ${WARP_CONTAINER_IP}"; then
        ip route add default via "$WARP_CONTAINER_IP" dev "$PUBLICNET_BRIDGE" table "$ROUTING_TABLE_ID"
        log_info "Added route: default via $WARP_CONTAINER_IP dev $PUBLICNET_BRIDGE"
    else
        log_debug "Custom route already exists"
    fi
}

setup_nat_rules() {
    log_step "Setting up NAT rules..."
    
    # Check if NAT rule already exists
    if ! iptables -t nat -C POSTROUTING -s "$WARP_SUBNET" -o "$PUBLICNET_BRIDGE" -j MASQUERADE 2>/dev/null; then
        iptables -t nat -A POSTROUTING -s "$WARP_SUBNET" -o "$PUBLICNET_BRIDGE" -j MASQUERADE
        log_info "Added iptables NAT rule for $WARP_SUBNET"
    else
        log_debug "NAT rule already exists"
    fi
}

#===========================================
# VERIFICATION FUNCTIONS
#===========================================

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
# STATUS AND TESTING FUNCTIONS
#===========================================

show_status() {
    print_section "CURRENT CONFIGURATION STATUS"
    
    echo "Network Configuration:"
    echo "  warpnet subnet: $WARP_SUBNET"
    echo "  publicnet subnet: $PUBLICNET_SUBNET"
    echo "  warp container IP: $WARP_CONTAINER_IP"
    echo
    
    echo "Bridge Interfaces:"
    echo "  publicnet: $PUBLICNET_BRIDGE"
    echo "  warpnet: $WARPNET_BRIDGE"
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
        docker inspect "$WARP_CONTAINER_NAME" --format '  IP: {{.NetworkSettings.Networks.publicnet.IPAddress}}'
    else
        echo "  warp-with-nat: NOT RUNNING"
    fi
}

test_setup() {
    print_section "COMPREHENSIVE ROUTING TESTS"
    
    if [[ "$RUN_TESTS" != "true" ]]; then
        log_info "Tests disabled (RUN_TESTS=false), skipping validation"
        return 0
    fi
    
    local test_failures=0
    local test_image="alpine:latest"
    
    # Ensure test image is available
    log_step "Preparing test environment..."
    if ! docker image inspect "$test_image" >/dev/null 2>&1; then
        log_info "Pulling test image: $test_image"
        docker pull "$test_image" >/dev/null
    fi
    
    # Test 1: Basic connectivity through warpnet
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
    
    # Test 4: Container-to-container communication
    log_step "Test 4: Container-to-container communication on warpnet"
    if test_container_communication; then
        log_success "✓ Container communication test passed"
    else
        log_error "✗ Container communication test failed"
        ((test_failures++))
    fi
    
    # Test 5: Routing table validation
    log_step "Test 5: Routing table configuration validation"
    if test_routing_configuration; then
        log_success "✓ Routing configuration test passed"
    else
        log_error "✗ Routing configuration test failed"
        ((test_failures++))
    fi
    
    # Test 6: WARP container health
    log_step "Test 6: WARP container health check"
    if test_warp_container_health; then
        log_success "✓ WARP container health test passed"
    else
        log_error "✗ WARP container health test failed"
        ((test_failures++))
    fi
    
    # Cleanup test containers
    cleanup_test_containers
    
    # Summary
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
        "$test_image" sh -c 'apk add --no-cache curl >/dev/null 2>&1 && curl -s --max-time 15 ifconfig.me' 2>/dev/null)
    
    if [[ -n "$test_result" && "$test_result" != "FAILED" ]]; then
        log_debug "Container external IP: $test_result"
        return 0
    else
        log_debug "Failed to get external IP from container"
        return 1
    fi
}

test_ip_routing() {
    # Get external IP from container on warpnet
    local container_ip
    container_ip=$(docker run --rm --name "warp-test-routing" --network=warpnet \
        "$test_image" sh -c 'apk add --no-cache curl >/dev/null 2>&1 && curl -s --max-time 15 ifconfig.me' 2>/dev/null)
    
    # Get direct external IP from host
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
    local container_names=()
    
    # Start multiple containers and check their IPs
    for i in $(seq 1 $container_count); do
        local container_name="warp-test-ip-$i"
        container_names+=("$container_name")
        
        # Start container in background
        docker run -d --name "$container_name" --network=warpnet \
            "$test_image" sleep 30 >/dev/null
        
        # Get assigned IP
        local assigned_ip
        assigned_ip=$(docker inspect "$container_name" -f '{{.NetworkSettings.Networks.warpnet.IPAddress}}' 2>/dev/null)
        
        if [[ -n "$assigned_ip" ]]; then
            assigned_ips+=("$assigned_ip")
            log_debug "Container $container_name assigned IP: $assigned_ip"
            
            # Verify IP is in correct subnet
            if ! ip route get "$assigned_ip" | grep -q "dev $WARPNET_BRIDGE"; then
                log_error "IP $assigned_ip is not in warpnet subnet"
                return 1
            fi
        else
            log_error "Could not get IP for container $container_name"
            return 1
        fi
    done
    
    # Check for unique IPs
    local unique_ips=($(printf '%s\n' "${assigned_ips[@]}" | sort -u))
    if [[ ${#unique_ips[@]} -eq ${#assigned_ips[@]} ]]; then
        log_debug "All containers received unique IPs: ${assigned_ips[*]}"
        return 0
    else
        log_error "IP assignment conflict detected"
        return 1
    fi
}

test_container_communication() {
    local server_container="warp-test-server"
    local client_container="warp-test-client"
    
    # Start a simple HTTP server container
    docker run -d --name "$server_container" --network=warpnet \
        "$test_image" sh -c 'echo "test-response" > /tmp/test.txt && while true; do echo -e "HTTP/1.1 200 OK\r\n\r\ntest-response" | nc -l -p 8080; done' >/dev/null
    
    sleep 2
    
    # Get server IP
    local server_ip
    server_ip=$(docker inspect "$server_container" -f '{{.NetworkSettings.Networks.warpnet.IPAddress}}' 2>/dev/null)
    
    if [[ -z "$server_ip" ]]; then
        log_error "Could not get server container IP"
        return 1
    fi
    
    # Test connection from client container
    local response
    response=$(docker run --rm --name "$client_container" --network=warpnet \
        "$test_image" sh -c "apk add --no-cache curl >/dev/null 2>&1 && timeout 5 curl -s http://$server_ip:8080/" 2>/dev/null)
    
    if [[ "$response" == "test-response" ]]; then
        log_debug "Container-to-container communication successful"
        return 0
    else
        log_debug "Container-to-container communication failed. Response: $response"
        return 1
    fi
}

test_routing_configuration() {
    # Test policy routing rule exists
    if ! ip rule list | grep -q "from ${WARP_SUBNET} lookup ${ROUTING_TABLE_ID}"; then
        log_error "Policy routing rule not found"
        return 1
    fi
    
    # Test custom route exists
    if ! ip route show table "$ROUTING_TABLE_ID" | grep -q "default via ${WARP_CONTAINER_IP}"; then
        log_error "Custom route not found"
        return 1
    fi
    
    # Test NAT rule exists
    if ! iptables -t nat -C POSTROUTING -s "$WARP_SUBNET" -o "$PUBLICNET_BRIDGE" -j MASQUERADE 2>/dev/null; then
        log_error "NAT rule not found"
        return 1
    fi
    
    # Test bridge interfaces exist
    if ! ip link show "$WARPNET_BRIDGE" >/dev/null 2>&1; then
        log_error "Warpnet bridge interface not found"
        return 1
    fi
    
    if ! ip link show "$PUBLICNET_BRIDGE" >/dev/null 2>&1; then
        log_error "Publicnet bridge interface not found"
        return 1
    fi
    
    log_debug "All routing components verified"
    return 0
}

test_warp_container_health() {
    # Check if container is running
    if ! check_container_running; then
        log_error "WARP container is not running"
        return 1
    fi
    
    # Check container health/logs for errors
    local container_logs
    container_logs=$(docker logs "$WARP_CONTAINER_NAME" --tail 10 2>/dev/null)
    
    # Look for common error patterns (adjust based on actual WARP container behavior)
    if echo "$container_logs" | grep -qi "error\|failed\|fatal"; then
        log_warning "Potential issues found in WARP container logs:"
        log_warning "$container_logs"
    fi
    
    # Test connectivity to WARP container
    if ping -c 1 -W 2 "$WARP_CONTAINER_IP" >/dev/null 2>&1; then
        log_debug "WARP container is reachable"
    else
        log_debug "WARP container ping failed (may be normal)"
    fi
    
    # Check if SOCKS proxy is responding (if enabled)
    if [[ -n "$WARP_SOCKS_PORT" ]]; then
        if timeout 2 bash -c "</dev/tcp/localhost/$WARP_SOCKS_PORT" 2>/dev/null; then
            log_debug "SOCKS proxy port is accessible"
        else
            log_debug "SOCKS proxy port not accessible (may be normal)"
        fi
    fi
    
    return 0
}

cleanup_test_containers() {
    log_debug "Cleaning up test containers..."
    
    # Remove any containers with our test naming pattern
    local test_containers
    test_containers=$(docker ps -aq --filter "name=warp-test-" 2>/dev/null || true)
    
    if [[ -n "$test_containers" ]]; then
        docker rm -f $test_containers >/dev/null 2>&1 || true
        log_debug "Removed test containers"
    fi
}

#===========================================
# MAIN EXECUTION FUNCTIONS
#===========================================

main() {
    print_section "DOCKER WARP NETWORK ROUTING SETUP"
    
    log_info "Starting setup with configuration:"
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
        
        # Run comprehensive tests
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
    
    # Optionally stop and remove the container
    if check_container_running; then
        if confirm_action "Stop and remove warp-with-nat container?"; then
            docker stop "$WARP_CONTAINER_NAME" >/dev/null 2>&1 || true
            docker rm "$WARP_CONTAINER_NAME" >/dev/null 2>&1 || true
            log_info "Container stopped and removed"
        fi
    fi
    
    log_success "Cleanup completed - normal Docker routing restored"
}

show_help() {
    echo "Docker WARP Network Routing Setup Script"
    echo
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo
    echo "Commands:"
    echo "  setup     - Set up routing (default)"
    echo "  cleanup   - Remove routing setup"
    echo "  status    - Show current status"
    echo "  test      - Test the routing setup"
    echo "  help      - Show this help message"
    echo
    echo "Environment Variables:"
    echo "  WARP_SUBNET              - warpnet subnet (default: 172.20.0.0/16)"
    echo "  PUBLICNET_SUBNET         - publicnet subnet (default: 10.76.0.0/16)"
    echo "  WARP_CONTAINER_IP        - warp container IP (default: 10.76.128.200)"
    echo "  ROUTING_TABLE_ID         - routing table ID (default: 200)"
    echo "  CONFIG_PATH              - config directory (default: ./configs)"
    echo "  FORCE_RECREATE           - skip confirmations (default: false)"
    echo "  AUTO_START_CONTAINER     - auto-start missing container (default: true)"
    echo "  RUN_TESTS                - run comprehensive tests (default: true)"
    echo "  CLEANUP_ON_TEST_FAILURE  - cleanup if tests fail (default: true)"
    echo "  DEBUG                    - enable debug output (default: false)"
    echo
    echo "Examples:"
    echo "  sudo $0                           # Setup with defaults and run tests"
    echo "  sudo FORCE_RECREATE=true $0       # Setup without prompts"
    echo "  sudo RUN_TESTS=false $0           # Setup without running tests"
    echo "  sudo DEBUG=true $0 status         # Show status with debug info"
    echo "  sudo $0 test                      # Run comprehensive tests only"
    echo "  sudo $0 cleanup                   # Remove setup"
}

#===========================================
# COMMAND LINE INTERFACE
#===========================================

case "${1:-setup}" in
    setup|"")
        main
        ;;
    cleanup|clean|remove)
        check_root
        cleanup_only
        ;;
    status|show)
        show_status
        ;;
    test|verify)
        test_setup
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