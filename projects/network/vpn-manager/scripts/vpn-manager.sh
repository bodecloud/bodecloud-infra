#!/bin/bash

# VPN Manager - Docker-in-Docker Edition with Docker Compose
# This script runs inside a privileged DinD container and manages VPN containers
# defined in an external docker-compose.yml file

set -e

# --- Configuration from Environment ---
HEALTH_CHECK_INTERVAL="${HEALTH_CHECK_INTERVAL:-60}"
MAX_FAILURES="${MAX_FAILURES:-3}"
DOCKER_NETWORK_NAME="${DOCKER_NETWORK_NAME:-vpn-network}"
DOCKER_NETWORK_CIDR="${DOCKER_NETWORK_CIDR:-10.45.0.0/16}"

# Get VPN service list from entrypoint (set by entrypoint.sh)
VPN_SERVICE_LIST="${VPN_SERVICE_LIST:-}"

# Convert comma-separated VPN list to array
if [ -z "$VPN_SERVICE_LIST" ]; then
    log "ERROR: VPN_SERVICE_LIST not provided by entrypoint"
    exit 1
fi

IFS=',' read -ra VPN_ARRAY <<< "$VPN_SERVICE_LIST"

# --- Globals ---
CURRENT_VPN_INDEX=0
FAILURE_COUNT=0
ACTIVE_SERVICE=""
COMPOSE_PROJECT_NAME="vpn-fallback"

# --- Functions ---

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] VPN-Manager: $1"
}

wait_for_docker() {
    log "Waiting for Docker daemon to be ready..."
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if docker info >/dev/null 2>&1; then
            log "Docker daemon is ready"
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 2
    done
    
    log "ERROR: Docker daemon failed to start"
    exit 1
}

setup_container_routing() {
    log "Setting up container routing and NAT"
    
    # Enable IP forwarding
    echo 1 > /proc/sys/net/ipv4/ip_forward
    
    # Set up iptables for NAT and forwarding
    # Clear existing rules
    iptables -t nat -F 2>/dev/null || true
    iptables -t filter -F FORWARD 2>/dev/null || true
    
    # Allow forwarding between networks
    iptables -A FORWARD -i eth0 -o docker+ -j ACCEPT
    iptables -A FORWARD -i docker+ -o eth0 -j ACCEPT
    iptables -A FORWARD -i docker+ -o docker+ -j ACCEPT
    
    # Set up NAT for outgoing traffic
    iptables -t nat -A POSTROUTING -o docker+ -j MASQUERADE
    
    log "Container routing configured"
}

start_vpn_service() {
    local service_name="$1"
    
    log "Starting VPN service: $service_name"
    
    # Stop the service if it's running
    docker-compose -f /app/vpn-compose.yml -p "$COMPOSE_PROJECT_NAME" stop "$service_name" 2>/dev/null || true
    docker-compose -f /app/vpn-compose.yml -p "$COMPOSE_PROJECT_NAME" rm -f "$service_name" 2>/dev/null || true
    
    # Start the specific service
    if ! docker-compose -f /app/vpn-compose.yml -p "$COMPOSE_PROJECT_NAME" up -d "$service_name"; then
        log "ERROR: Failed to start service $service_name"
        return 1
    fi
    
    ACTIVE_SERVICE="$service_name"
    
    # Wait for service to establish VPN connection
    log "Waiting for VPN connection to establish..."
    sleep 15
    
    # Update routing to use the new active VPN
    update_routing_to_active_vpn
}

get_service_container_name() {
    local service_name="$1"
    echo "${COMPOSE_PROJECT_NAME}_${service_name}_1"
}

get_service_container_ip() {
    local service_name="$1"
    local container_name
    container_name=$(get_service_container_name "$service_name")
    
    # Try to get IP from any network the container is connected to
    docker inspect "$container_name" 2>/dev/null | jq -r '.[0].NetworkSettings.Networks | to_entries | .[0].value.IPAddress' 2>/dev/null || echo ""
}

update_routing_to_active_vpn() {
    if [ -z "$ACTIVE_SERVICE" ]; then
        log "ERROR: No active VPN service"
        return 1
    fi
    
    # Get the active VPN service's container IP
    local vpn_container_ip
    vpn_container_ip=$(get_service_container_ip "$ACTIVE_SERVICE")
    
    if [ -z "$vpn_container_ip" ] || [ "$vpn_container_ip" = "null" ]; then
        log "ERROR: Could not get IP for active VPN service $ACTIVE_SERVICE"
        return 1
    fi
    
    log "Updating routing to use VPN service $ACTIVE_SERVICE ($vpn_container_ip)"
    
    # Update default route to go through the active VPN container
    # Remove old default route
    ip route del default 2>/dev/null || true
    
    # Add new default route through active VPN container
    ip route add default via "$vpn_container_ip"
    
    # Update iptables to forward traffic through the active VPN
    iptables -t nat -F POSTROUTING
    iptables -t nat -A POSTROUTING -o docker+ -j MASQUERADE
    
    log "Routing updated to use $ACTIVE_SERVICE"
}

health_check() {
    local service_name="$1"
    local container_name
    container_name=$(get_service_container_name "$service_name")
    
    # Check if container is running
    if ! docker ps --format "table {{.Names}}" | grep -q "^${container_name}$"; then
        log "Container $container_name is not running"
        return 1
    fi
    
    # Get container IP
    local container_ip
    container_ip=$(get_service_container_ip "$service_name")
    
    if [ -z "$container_ip" ] || [ "$container_ip" = "null" ]; then
        log "Could not get IP for container $container_name"
        return 1
    fi
    
    # Test internet connectivity through the container
    # We'll try multiple methods since different VPN types have different access patterns
    local test_result=""
    
    # Method 1: Try direct connection through container's network namespace
    test_result=$(docker exec "$container_name" wget -qO- --timeout=10 --tries=1 "https://api.ipify.org" 2>/dev/null || echo "")
    
    # Method 2: If that fails, try with curl
    if [ -z "$test_result" ]; then
        test_result=$(docker exec "$container_name" curl -s --connect-timeout 10 --max-time 30 "https://api.ipify.org" 2>/dev/null || echo "")
    fi
    
    # Method 3: If container doesn't have wget/curl, try with a test container using the same network
    if [ -z "$test_result" ]; then
        test_result=$(docker run --rm --network container:"$container_name" \
            curlimages/curl:latest \
            --connect-timeout 10 \
            --max-time 30 \
            --silent \
            --fail \
            "https://api.ipify.org" 2>/dev/null || echo "")
    fi
    
    # Method 4: For SOCKS proxy containers (like WARP), try SOCKS5
    if [ -z "$test_result" ]; then
        # Check if port 1080 is open (common SOCKS port)
        if docker exec "$container_name" netstat -ln 2>/dev/null | grep -q ":1080"; then
            test_result=$(docker run --rm --network container:"$container_name" \
                curlimages/curl:latest \
                --connect-timeout 10 \
                --max-time 30 \
                --proxy "socks5://127.0.0.1:1080" \
                --silent \
                --fail \
                "https://api.ipify.org" 2>/dev/null || echo "")
        fi
    fi
    
    if [ -z "$test_result" ]; then
        log "Health check failed for $service_name - no response from external IP service"
        return 1
    fi
    
    # Validate that we got an IP address
    if echo "$test_result" | grep -qE '^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$'; then
        log "Health check passed for $service_name (External IP: $test_result)"
        return 0
    else
        log "Health check failed for $service_name - invalid response: $test_result"
        return 1
    fi
}

switch_to_next_vpn() {
    # Move to next VPN in the list
    CURRENT_VPN_INDEX=$(( (CURRENT_VPN_INDEX + 1) % ${#VPN_ARRAY[@]} ))
    local next_vpn="${VPN_ARRAY[$CURRENT_VPN_INDEX]}"
    
    log "Switching to VPN service: $next_vpn"
    
    # Stop current VPN service
    if [ -n "$ACTIVE_SERVICE" ]; then
        log "Stopping current service: $ACTIVE_SERVICE"
        docker-compose -f /app/vpn-compose.yml -p "$COMPOSE_PROJECT_NAME" stop "$ACTIVE_SERVICE" 2>/dev/null || true
    fi
    
    # Start new VPN service
    if start_vpn_service "$next_vpn"; then
        log "Successfully switched to $next_vpn"
        FAILURE_COUNT=0
        return 0
    else
        log "Failed to switch to $next_vpn"
        FAILURE_COUNT=$((FAILURE_COUNT + 1))
        return 1
    fi
}

cleanup() {
    log "Shutting down VPN manager..."
    
    # Stop all VPN services
    if [ -f "/app/vpn-compose.yml" ]; then
        docker-compose -f /app/vpn-compose.yml -p "$COMPOSE_PROJECT_NAME" down 2>/dev/null || true
    fi
    
    log "Cleanup complete"
    exit 0
}

trap cleanup SIGINT SIGTERM

# --- Main Execution ---

log "Starting VPN Manager (Docker Compose Edition)"
log "VPN Services: $VPN_SERVICE_LIST"
log "Network: $DOCKER_NETWORK_NAME ($DOCKER_NETWORK_CIDR)"

# Validate that we have the compose file
if [ ! -f "/app/vpn-compose.yml" ]; then
    log "ERROR: VPN compose file not found at /app/vpn-compose.yml"
    exit 1
fi

wait_for_docker
setup_container_routing

# Start with first VPN
CURRENT_VPN_INDEX=0
switch_to_next_vpn

log "VPN Manager ready - acting as gateway for $DOCKER_NETWORK_NAME"

# Main monitoring loop
while true; do
    local current_vpn="${VPN_ARRAY[$CURRENT_VPN_INDEX]}"
    log "Current VPN: $current_vpn, Failures: $FAILURE_COUNT/$MAX_FAILURES"
    
    if [ $FAILURE_COUNT -ge $MAX_FAILURES ]; then
        log "ERROR: Maximum failures reached. Exiting."
        exit 1
    fi
    
    if health_check "$current_vpn"; then
        log "Health check passed for $current_vpn"
        FAILURE_COUNT=0
    else
        log "Health check failed for $current_vpn"
        FAILURE_COUNT=$((FAILURE_COUNT + 1))
        
        if ! switch_to_next_vpn; then
            log "Failed to switch VPN"
        fi
    fi
    
    sleep $HEALTH_CHECK_INTERVAL
done 