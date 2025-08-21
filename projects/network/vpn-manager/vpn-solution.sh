#!/bin/bash

# vpn-solution.sh - Host-side VPN failover solution
# This script runs on the HOST and manages:
# 1. A DinD container that reads an external docker-compose.yml with VPN services
# 2. Host-side routing setup (ip route/ip rule) to route traffic through the VPN

set -e

# Default configuration
DEFAULT_DOCKER_NETWORK_NAME="vpn-network"
DEFAULT_DOCKER_NETWORK_CIDR="10.45.0.0/16"
DEFAULT_DOCKER_NETWORK_GATEWAY="10.45.0.1"
DEFAULT_VPN_MANAGER_IMAGE="vpn-fallback:latest"
DEFAULT_VPN_COMPOSE_FILE=""

# Configuration variables (can be overridden by arguments)
DOCKER_NETWORK_NAME="${DEFAULT_DOCKER_NETWORK_NAME}"
DOCKER_NETWORK_CIDR="${DEFAULT_DOCKER_NETWORK_CIDR}"
DOCKER_NETWORK_GATEWAY="${DEFAULT_DOCKER_NETWORK_GATEWAY}"
VPN_MANAGER_IMAGE="${DEFAULT_VPN_MANAGER_IMAGE}"
VPN_COMPOSE_FILE="${DEFAULT_VPN_COMPOSE_FILE}"
VPN_COMPOSE_CONTENT=""
ACTION=""

# Runtime variables
VPN_MANAGER_CONTAINER="vpn-fallback-manager"
HOST_ROUTING_TABLE="vpn"
HOST_ROUTING_TABLE_ID="100"

usage() {
    cat << EOF
Usage: $0 <action> [options]

Actions:
    start       Start the VPN solution
    stop        Stop the VPN solution
    restart     Restart the VPN solution
    status      Show status of the VPN solution
    logs        Show logs from the VPN manager

Options:
    --network-name <name>       Docker network name (default: ${DEFAULT_DOCKER_NETWORK_NAME})
    --network-cidr <cidr>       Docker network CIDR (default: ${DEFAULT_DOCKER_NETWORK_CIDR})
    --network-gateway <ip>      Docker network gateway (default: ${DEFAULT_DOCKER_NETWORK_GATEWAY})
    --compose-file <path>       Path to docker-compose.yml with VPN services (required)
    --compose-content <yaml>    Docker compose content as string (alternative to --compose-file)
    --image <image>             VPN manager Docker image (default: ${DEFAULT_VPN_MANAGER_IMAGE})

Examples:
    $0 start --compose-file ./my-vpns.yml
    $0 start --compose-file /path/to/vpns.yml --network-name my-vpn
    $0 stop
    $0 status
    $0 logs

The docker-compose.yml should contain VPN services in the order you want them to be tried.
Each service in the compose file will be treated as a VPN option in the failover list.
The services will be tried in the order they appear in the file (top to bottom).

Example compose file:
    services:
      gluetun-de:
        image: ghcr.io/qdm12/gluetun
        cap_add: [NET_ADMIN]
        environment:
          VPN_SERVICE_PROVIDER: custom
          # ... other gluetun config
      warp:
        image: caomingjun/warp:latest
        cap_add: [NET_ADMIN, MKNOD, AUDIT_WRITE]
        # ... other warp config
EOF
}

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] HOST: $1"
}

parse_arguments() {
    if [ $# -eq 0 ]; then
        usage
        exit 1
    fi

    ACTION="$1"
    shift

    while [[ $# -gt 0 ]]; do
        case $1 in
            --network-name)
                DOCKER_NETWORK_NAME="$2"
                shift 2
                ;;
            --network-cidr)
                DOCKER_NETWORK_CIDR="$2"
                shift 2
                ;;
            --network-gateway)
                DOCKER_NETWORK_GATEWAY="$2"
                shift 2
                ;;
            --compose-file)
                VPN_COMPOSE_FILE="$2"
                shift 2
                ;;
            --compose-content)
                VPN_COMPOSE_CONTENT="$2"
                shift 2
                ;;
            --image)
                VPN_MANAGER_IMAGE="$2"
                shift 2
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done

    case "$ACTION" in
        start|stop|restart|status|logs)
            ;;
        *)
            echo "Unknown action: $ACTION"
            usage
            exit 1
            ;;
    esac

    # Validate required arguments for start action
    if [ "$ACTION" = "start" ]; then
        if [ -z "$VPN_COMPOSE_FILE" ] && [ -z "$VPN_COMPOSE_CONTENT" ]; then
            echo "ERROR: Either --compose-file or --compose-content is required for start action"
            usage
            exit 1
        fi
        
        if [ -n "$VPN_COMPOSE_FILE" ] && [ ! -f "$VPN_COMPOSE_FILE" ]; then
            echo "ERROR: Compose file not found: $VPN_COMPOSE_FILE"
            exit 1
        fi
    fi
}

check_privileges() {
    if [ "$EUID" -ne 0 ]; then
        echo "ERROR: This script must be run as root (for ip route/ip rule commands)"
        exit 1
    fi
}

create_docker_network() {
    log "Creating Docker network: $DOCKER_NETWORK_NAME"
    
    # Remove existing network if it exists
    if docker network inspect "$DOCKER_NETWORK_NAME" >/dev/null 2>&1; then
        log "Network $DOCKER_NETWORK_NAME already exists, removing..."
        docker network rm "$DOCKER_NETWORK_NAME" || true
    fi

    # Create the network that external containers will connect to
    if ! docker network create \
        -d bridge \
        --attachable \
        --subnet "$DOCKER_NETWORK_CIDR" \
        --gateway "$DOCKER_NETWORK_GATEWAY" \
        -o com.docker.network.bridge.name="br_${DOCKER_NETWORK_NAME}" \
        "$DOCKER_NETWORK_NAME"; then
        log "ERROR: Could not create $DOCKER_NETWORK_NAME"
        exit 1
    fi

    log "Docker network created successfully"
}

setup_host_routing() {
    log "Setting up host routing for $DOCKER_NETWORK_NAME"

    local docker_bridge="br_${DOCKER_NETWORK_NAME}"
    local docker_net="$DOCKER_NETWORK_CIDR"

    # Dynamically determine the local network from the default interface
    local local_interface
    local_interface=$(ip route | grep default | awk '{print $5}' | head -1)
    if [ -z "$local_interface" ]; then
        log "ERROR: Could not determine local interface"
        exit 1
    fi

    local local_net
    local_net=$(ip route | grep "$local_interface" | grep -v default | awk '{print $1}' | head -1)

    # Dynamically determine the default gateway
    local local_gateway
    local_gateway=$(ip route | grep default | awk '{print $3}' | head -1)

    log "docker_net=$docker_net"
    log "local_net=$local_net"
    log "local_gateway=$local_gateway"
    log "docker_bridge=$docker_bridge"

    # Check if routing table 'vpn' exists, create if missing with proper table ID
    if ! grep -q "^[0-9]\+[[:space:]]\+${HOST_ROUTING_TABLE}$" /etc/iproute2/rt_tables; then
        # Find an available table ID
        local table_id=$HOST_ROUTING_TABLE_ID
        while grep -q "^${table_id}[[:space:]]" /etc/iproute2/rt_tables; do
            table_id=$((table_id + 1))
            if [ $table_id -gt 255 ]; then
                log "ERROR: No available table IDs found"
                exit 1
            fi
        done
        echo "${table_id}     ${HOST_ROUTING_TABLE}" >> /etc/iproute2/rt_tables
        log "Created routing table '${HOST_ROUTING_TABLE}' with ID ${table_id}"
        HOST_ROUTING_TABLE_ID=$table_id
    fi

    # Remove any previous routes in the routing table
    ip rule | sed -n "s/.*\(from[ \t]*[0-9\.\/]*\).*${HOST_ROUTING_TABLE}/\1/p" | while read RULE; do
        [ -n "$RULE" ] && ip rule del $RULE 2>/dev/null || true
    done
    ip route flush table $HOST_ROUTING_TABLE 2>/dev/null || true

    # Traffic coming FROM the docker network should go through the VPN table
    ip rule add from ${docker_net} lookup $HOST_ROUTING_TABLE

    # Get the VPN manager container IP (it will act as our gateway)
    local vpn_manager_ip
    vpn_manager_ip=$(get_vpn_manager_ip)
    
    if [ -z "$vpn_manager_ip" ]; then
        log "ERROR: Could not get VPN manager container IP"
        exit 1
    fi

    log "VPN manager container IP: $vpn_manager_ip"

    # Set up routes in the VPN table to go through the VPN manager container
    ip route add default via $vpn_manager_ip table $HOST_ROUTING_TABLE
    ip route add $local_net dev $local_interface table $HOST_ROUTING_TABLE
    ip route add $docker_net dev $docker_bridge table $HOST_ROUTING_TABLE

    log "Host routing configured successfully"
}

teardown_host_routing() {
    log "Tearing down host routing"

    local docker_net="$DOCKER_NETWORK_CIDR"

    # Remove ip rules
    ip rule | sed -n "s/.*\(from[ \t]*[0-9\.\/]*\).*${HOST_ROUTING_TABLE}/\1/p" | while read RULE; do
        [ -n "$RULE" ] && ip rule del $RULE 2>/dev/null || true
    done

    # Flush the routing table
    ip route flush table $HOST_ROUTING_TABLE 2>/dev/null || true

    log "Host routing cleaned up"
}

get_vpn_manager_ip() {
    if docker ps --format "table {{.Names}}" | grep -q "^${VPN_MANAGER_CONTAINER}$"; then
        docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$VPN_MANAGER_CONTAINER" 2>/dev/null
    fi
}

start_vpn_manager() {
    log "Starting VPN manager container"

    # Stop existing container if running
    if docker ps -a --format "table {{.Names}}" | grep -q "^${VPN_MANAGER_CONTAINER}$"; then
        log "Stopping existing VPN manager container"
        docker stop "$VPN_MANAGER_CONTAINER" >/dev/null 2>&1 || true
        docker rm "$VPN_MANAGER_CONTAINER" >/dev/null 2>&1 || true
    fi

    # Prepare docker run command
    local docker_run_cmd="docker run -d \
        --name $VPN_MANAGER_CONTAINER \
        --privileged \
        --network $DOCKER_NETWORK_NAME \
        --env DOCKER_NETWORK_NAME=$DOCKER_NETWORK_NAME \
        --env DOCKER_NETWORK_CIDR=$DOCKER_NETWORK_CIDR \
        --restart unless-stopped"

    # Add compose file or content
    if [ -n "$VPN_COMPOSE_FILE" ]; then
        # Mount the compose file
        local abs_compose_file
        abs_compose_file=$(realpath "$VPN_COMPOSE_FILE")
        docker_run_cmd="$docker_run_cmd --volume $abs_compose_file:/compose.yml:ro --env VPN_COMPOSE_FILE=/compose.yml"
    elif [ -n "$VPN_COMPOSE_CONTENT" ]; then
        # Pass compose content as environment variable
        docker_run_cmd="$docker_run_cmd --env VPN_COMPOSE=$VPN_COMPOSE_CONTENT"
    fi

    # Add image name
    docker_run_cmd="$docker_run_cmd $VPN_MANAGER_IMAGE"

    # Execute the docker run command
    eval $docker_run_cmd

    # Wait for container to be ready
    log "Waiting for VPN manager to be ready..."
    local max_attempts=60
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if docker logs "$VPN_MANAGER_CONTAINER" 2>&1 | grep -q "VPN Manager ready"; then
            log "VPN manager is ready"
            break
        fi
        attempt=$((attempt + 1))
        sleep 2
    done

    if [ $attempt -eq $max_attempts ]; then
        log "WARNING: VPN manager may not be fully ready yet"
        log "Check logs with: $0 logs"
    fi
}

stop_vpn_manager() {
    log "Stopping VPN manager container"
    
    if docker ps --format "table {{.Names}}" | grep -q "^${VPN_MANAGER_CONTAINER}$"; then
        docker stop "$VPN_MANAGER_CONTAINER" >/dev/null 2>&1 || true
        docker rm "$VPN_MANAGER_CONTAINER" >/dev/null 2>&1 || true
        log "VPN manager container stopped"
    else
        log "VPN manager container is not running"
    fi
}

remove_docker_network() {
    log "Removing Docker network: $DOCKER_NETWORK_NAME"
    
    if docker network inspect "$DOCKER_NETWORK_NAME" >/dev/null 2>&1; then
        docker network rm "$DOCKER_NETWORK_NAME" || true
        log "Docker network removed"
    else
        log "Docker network does not exist"
    fi
}

start_solution() {
    log "Starting VPN solution"
    log "Network: $DOCKER_NETWORK_NAME ($DOCKER_NETWORK_CIDR)"
    if [ -n "$VPN_COMPOSE_FILE" ]; then
        log "Compose file: $VPN_COMPOSE_FILE"
    else
        log "Using compose content from environment variable"
    fi

    create_docker_network
    start_vpn_manager
    setup_host_routing

    log "VPN solution started successfully"
    log "Other containers can now connect to network '$DOCKER_NETWORK_NAME' to use VPN"
}

stop_solution() {
    log "Stopping VPN solution"

    teardown_host_routing
    stop_vpn_manager
    remove_docker_network

    log "VPN solution stopped"
}

show_status() {
    echo "=== VPN Solution Status ==="
    echo "Network Name: $DOCKER_NETWORK_NAME"
    echo "Network CIDR: $DOCKER_NETWORK_CIDR"
    if [ -n "$VPN_COMPOSE_FILE" ]; then
        echo "Compose File: $VPN_COMPOSE_FILE"
    else
        echo "Compose Content: [from environment variable]"
    fi
    echo ""

    # Check if network exists
    if docker network inspect "$DOCKER_NETWORK_NAME" >/dev/null 2>&1; then
        echo "✓ Docker network '$DOCKER_NETWORK_NAME' exists"
    else
        echo "✗ Docker network '$DOCKER_NETWORK_NAME' does not exist"
    fi

    # Check if VPN manager is running
    if docker ps --format "table {{.Names}}" | grep -q "^${VPN_MANAGER_CONTAINER}$"; then
        local vpn_manager_ip
        vpn_manager_ip=$(get_vpn_manager_ip)
        echo "✓ VPN manager container is running (IP: $vpn_manager_ip)"
        
        # Try to get current VPN from logs
        local current_vpn
        current_vpn=$(docker logs --tail 10 "$VPN_MANAGER_CONTAINER" 2>&1 | grep "Current VPN:" | tail -1 | sed 's/.*Current VPN: \([^,]*\).*/\1/' || echo "unknown")
        echo "  - Current active VPN: $current_vpn"
    else
        echo "✗ VPN manager container is not running"
    fi

    # Check routing table
    if grep -q "^[0-9]\+[[:space:]]\+${HOST_ROUTING_TABLE}$" /etc/iproute2/rt_tables; then
        echo "✓ Host routing table '$HOST_ROUTING_TABLE' exists"
        
        # Show active rules
        local rules_count
        rules_count=$(ip rule | grep -c "$HOST_ROUTING_TABLE" || echo "0")
        echo "  - Active routing rules: $rules_count"
    else
        echo "✗ Host routing table '$HOST_ROUTING_TABLE' does not exist"
    fi
}

show_logs() {
    if docker ps --format "table {{.Names}}" | grep -q "^${VPN_MANAGER_CONTAINER}$"; then
        echo "=== VPN Manager Logs ==="
        docker logs -f "$VPN_MANAGER_CONTAINER"
    else
        echo "VPN manager container is not running"
        exit 1
    fi
}

# Main execution
parse_arguments "$@"
check_privileges

case "$ACTION" in
    start)
        start_solution
        ;;
    stop)
        stop_solution
        ;;
    restart)
        stop_solution
        sleep 2
        start_solution
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
esac 