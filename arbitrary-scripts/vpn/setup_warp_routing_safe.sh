#!/bin/bash

set -e

#===========================================
# SAFE WARP ROUTING SETUP WITH DIND TESTING
#===========================================

# Configuration
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
DEBUG=${DEBUG:-false}
FORCE_APPLY=${FORCE_APPLY:-false}
DIND_CONTAINER_NAME="warp-routing-test-dind"
TEST_TIMEOUT=${TEST_TIMEOUT:-300}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m'

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

check_prerequisites() {
    print_section "CHECKING PREREQUISITES"
    
    local missing_tools=()
    local required_commands=("docker" "ip" "iptables")
    
    for cmd in "${required_commands[@]}"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing_tools+=("$cmd")
        fi
    done
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker daemon is not running"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

create_dind_test_environment() {
    print_section "CREATING ISOLATED TEST ENVIRONMENT"
    
    log_step "Starting Docker-in-Docker container..."
    
    # Clean up any existing test container
    docker rm -f "$DIND_CONTAINER_NAME" >/dev/null 2>&1 || true
    
    # Create DinD container with necessary privileges
    docker run -d \
        --name "$DIND_CONTAINER_NAME" \
        --privileged \
        --cgroupns=host \
        -v /var/lib/docker \
        -e DOCKER_TLS_CERTDIR= \
        docker:24-dind >/dev/null
    
    # Wait for Docker daemon to start inside DinD
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

create_test_script() {
    local script_content=$(cat << 'EOF'
#!/bin/sh
set -e

# Configuration passed from main script
WARP_SUBNET="WARP_SUBNET_PLACEHOLDER"
PUBLICNET_SUBNET="PUBLICNET_SUBNET_PLACEHOLDER"
WARP_CONTAINER_IP="WARP_CONTAINER_IP_PLACEHOLDER"
ROUTING_TABLE_ID="ROUTING_TABLE_ID_PLACEHOLDER"
ROUTING_TABLE_NAME="ROUTING_TABLE_NAME_PLACEHOLDER"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[TEST]${NC} $1"; }
log_success() { echo -e "${GREEN}[TEST]${NC} $1"; }
log_error() { echo -e "${RED}[TEST]${NC} $1"; }

# Install required packages
apk add --no-cache curl iproute2 iptables >/dev/null 2>&1

log_info "Setting up test networks..."

# Create test networks
docker network create --driver=bridge --subnet="$PUBLICNET_SUBNET" --gateway="PUBLICNET_GATEWAY_PLACEHOLDER" publicnet >/dev/null 2>&1
docker network create --driver=bridge --subnet="$WARP_SUBNET" --gateway="WARP_GATEWAY_PLACEHOLDER" warpnet >/dev/null 2>&1

# Start mock WARP container (simple proxy for testing)
log_info "Starting mock WARP container..."
docker run -d --name warp-with-nat --network publicnet --ip "$WARP_CONTAINER_IP" \
    alpine:latest sh -c 'apk add --no-cache iptables && 
    echo 1 > /proc/sys/net/ipv4/ip_forward &&
    iptables -t nat -A POSTROUTING -j MASQUERADE &&
    while true; do sleep 1; done' >/dev/null 2>&1

# Get bridge interfaces
PUBLICNET_BRIDGE=$(docker network inspect publicnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})
WARPNET_BRIDGE=$(docker network inspect warpnet -f '{{.Id}}' | cut -c1-12 | xargs -I{} echo br-{})

log_info "Setting up routing..."

# Add routing table
echo "$ROUTING_TABLE_ID $ROUTING_TABLE_NAME" >> /etc/iproute2/rt_tables

# Set up policy routing
ip rule add from "$WARP_SUBNET" table "$ROUTING_TABLE_ID"
ip route add default via "$WARP_CONTAINER_IP" dev "$PUBLICNET_BRIDGE" table "$ROUTING_TABLE_ID"

# Set up NAT
iptables -t nat -A POSTROUTING -s "$WARP_SUBNET" -o "$PUBLICNET_BRIDGE" -j MASQUERADE

log_info "Testing routing configuration..."

# Test 1: Basic connectivity
log_info "Test 1: Basic container connectivity"
TEST_RESULT=$(docker run --rm --network=warpnet alpine:latest sh -c 'apk add --no-cache curl >/dev/null 2>&1 && timeout 10 curl -s ifconfig.me' 2>/dev/null || echo "FAILED")

if [[ "$TEST_RESULT" == "FAILED" || -z "$TEST_RESULT" ]]; then
    log_error "Basic connectivity test failed"
    exit 1
fi
log_success "Container can reach internet: $TEST_RESULT"

# Test 2: IP assignment
log_info "Test 2: IP assignment validation"
CONTAINER_IP=$(docker run -d --network=warpnet alpine:latest sleep 10)
ASSIGNED_IP=$(docker inspect "$CONTAINER_IP" -f '{{.NetworkSettings.Networks.warpnet.IPAddress}}')
docker rm -f "$CONTAINER_IP" >/dev/null 2>&1

if [[ -z "$ASSIGNED_IP" ]]; then
    log_error "IP assignment test failed"
    exit 1
fi
log_success "Container assigned IP: $ASSIGNED_IP"

# Test 3: Routing validation
log_info "Test 3: Routing rules validation"
if ! ip rule list | grep -q "from $WARP_SUBNET lookup $ROUTING_TABLE_ID"; then
    log_error "Policy routing rule not found"
    exit 1
fi

if ! ip route show table "$ROUTING_TABLE_ID" | grep -q "default via $WARP_CONTAINER_IP"; then
    log_error "Custom route not found"
    exit 1
fi

if ! iptables -t nat -C POSTROUTING -s "$WARP_SUBNET" -o "$PUBLICNET_BRIDGE" -j MASQUERADE 2>/dev/null; then
    log_error "NAT rule not found"
    exit 1
fi

log_success "All routing rules validated"

log_success "🎉 All tests passed in isolated environment!"
echo "ROUTING_TEST_SUCCESS"
EOF
)
    
    # Replace placeholders with actual values
    script_content="${script_content//WARP_SUBNET_PLACEHOLDER/$WARP_SUBNET}"
    script_content="${script_content//PUBLICNET_SUBNET_PLACEHOLDER/$PUBLICNET_SUBNET}"
    script_content="${script_content//PUBLICNET_GATEWAY_PLACEHOLDER/$PUBLICNET_GATEWAY}"
    script_content="${script_content//WARP_GATEWAY_PLACEHOLDER/$WARP_GATEWAY}"
    script_content="${script_content//WARP_CONTAINER_IP_PLACEHOLDER/$WARP_CONTAINER_IP}"
    script_content="${script_content//ROUTING_TABLE_ID_PLACEHOLDER/$ROUTING_TABLE_ID}"
    script_content="${script_content//ROUTING_TABLE_NAME_PLACEHOLDER/$ROUTING_TABLE_NAME}"
    
    echo "$script_content"
}

run_isolated_test() {
    print_section "RUNNING ISOLATED ROUTING TEST"
    
    log_step "Creating test script..."
    local test_script=$(create_test_script)
    
    # Copy test script to DinD container
    echo "$test_script" | docker exec -i "$DIND_CONTAINER_NAME" sh -c 'cat > /test_routing.sh && chmod +x /test_routing.sh'
    
    log_step "Running routing test in isolated environment..."
    
    # Run the test with timeout
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

apply_to_host() {
    print_section "APPLYING VALIDATED CONFIGURATION TO HOST"
    
    if [[ "$FORCE_APPLY" != "true" ]]; then
        echo -n "Tests passed in isolated environment. Apply to host? [y/N]: "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            log_info "Skipping host application"
            return 0
        fi
    fi
    
    log_step "Applying configuration to host..."
    
    # Use the original setup script but with tests disabled
    export RUN_TESTS=false
    export AUTO_START_CONTAINER=true
    
    # Check if the original script exists
    if [[ -f "./setup_warp_routing.sh" ]]; then
        log_info "Running validated setup on host..."
        if ./setup_warp_routing.sh; then
            log_success "✅ Configuration successfully applied to host!"
            
            # Run a quick verification test on the host
            log_step "Running verification test on host..."
            HOST_IP=$(timeout 10 curl -s ifconfig.me 2>/dev/null || echo "UNKNOWN")
            CONTAINER_IP=$(docker run --rm --network=warpnet alpine:latest sh -c 'apk add --no-cache curl >/dev/null 2>&1 && timeout 15 curl -s ifconfig.me' 2>/dev/null || echo "FAILED")
            
            if [[ "$CONTAINER_IP" != "FAILED" && "$CONTAINER_IP" != "$HOST_IP" && -n "$CONTAINER_IP" ]]; then
                log_success "🎉 Host verification successful!"
                log_info "Host IP: $HOST_IP"
                log_info "Container IP (via WARP): $CONTAINER_IP"
            else
                log_warning "⚠ Host verification had issues, but setup may still be working"
            fi
        else
            log_error "❌ Failed to apply configuration to host"
            return 1
        fi
    else
        log_error "setup_warp_routing.sh not found in current directory"
        return 1
    fi
}

cleanup_dind_environment() {
    print_section "CLEANING UP TEST ENVIRONMENT"
    
    log_step "Removing Docker-in-Docker container..."
    docker rm -f "$DIND_CONTAINER_NAME" >/dev/null 2>&1 || true
    log_success "Test environment cleaned up"
}

show_help() {
    echo "Safe WARP Routing Setup with Isolated Testing"
    echo
    echo "This script tests the WARP routing configuration in an isolated"
    echo "Docker-in-Docker environment before applying to the host."
    echo
    echo "Usage: $0 [COMMAND]"
    echo
    echo "Commands:"
    echo "  test-only   - Only run isolated test, don't apply to host"
    echo "  setup       - Test in isolation, then apply to host (default)"
    echo "  help        - Show this help"
    echo
    echo "Environment Variables:"
    echo "  FORCE_APPLY=true         - Apply to host without confirmation"
    echo "  DEBUG=true               - Enable debug output"
    echo "  TEST_TIMEOUT=300         - Test timeout in seconds"
    echo
    echo "Examples:"
    echo "  sudo $0                  # Test then apply"
    echo "  sudo $0 test-only        # Test only"
    echo "  sudo FORCE_APPLY=true $0 # Test and auto-apply"
}

main() {
    print_section "SAFE WARP ROUTING SETUP"
    log_info "This script will test routing in isolation before applying to host"
    echo
    
    check_root
    check_prerequisites
    create_dind_test_environment
    
    if run_isolated_test; then
        cleanup_dind_environment
        
        if [[ "${1:-setup}" != "test-only" ]]; then
            apply_to_host
        else
            log_success "✅ Isolated test completed successfully!"
            log_info "Run without 'test-only' to apply to host"
        fi
    else
        log_error "❌ Isolated test failed - NOT applying to host"
        cleanup_dind_environment
        exit 1
    fi
}

cleanup_only() {
    cleanup_dind_environment
}

# Handle command line arguments
case "${1:-setup}" in
    setup|"")
        main
        ;;
    test-only)
        main test-only
        ;;
    cleanup)
        check_root
        cleanup_only
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac 