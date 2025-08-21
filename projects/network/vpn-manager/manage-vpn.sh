#!/bin/bash

# vpn-manager/manage-vpn.sh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
# List of VPN services from your docker-compose.networks.yml to use for failover.
# Add or remove services from this list as needed.
VPN_SERVICES=(
    "gluetun-premiumize-de"
    "gluetun-premiumize-nl"
    "gluetun-premiumize-us"
    "warp"
)

# Network names
DOCKER_COMPOSE_PROJECT_NAME="my-media-stack"
PUBLIC_NETWORK="publicnet"

# Healthcheck settings
HEALTH_CHECK_INTERVAL=60 # seconds
MAX_FAILURES=3

# --- Globals ---
CURRENT_VPN_INDEX=0
FAILURE_COUNT=0

# --- Functions ---

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] - $1"
}

# This function will be defined in a separate health-check.sh script
# For now, we stub it.
# perform_health_check() {
#     ...
# }

# This function will be defined in a separate setup-routing.sh script
# For now, we stub it.
# setup_routing() {
#    ...
# }

# This function will be defined in a separate teardown-routing.sh script
# For now, we stub it.
# teardown_routing() {
#    ...
# }

cleanup() {
    log "Cleaning up..."
    local active_service=${VPN_SERVICES[$CURRENT_VPN_INDEX]}
    docker-compose -p "$DOCKER_COMPOSE_PROJECT_NAME" stop "$active_service" || true
    ./vpn-manager/teardown-routing.sh
    log "Cleanup complete."
    exit 0
}

trap cleanup SIGINT SIGTERM

start_next_vpn() {
    # Stop the current VPN if any is running
    if [ -n "${VPN_SERVICES[$CURRENT_VPN_INDEX]}" ]; then
        local old_service=${VPN_SERVICES[$CURRENT_VPN_INDEX]}
        log "Stopping old VPN service: $old_service"
        docker-compose -p "$DOCKER_COMPOSE_PROJECT_NAME" stop "$old_service" >/dev/null 2>&1 || true
        # Clean up old routing
        ./vpn-manager/teardown-routing.sh
    fi

    # Move to the next VPN in the list
    CURRENT_VPN_INDEX=$(( (CURRENT_VPN_INDEX + 1) % ${#VPN_SERVICES[@]} ))
    local new_service=${VPN_SERVICES[$CURRENT_VPN_INDEX]}
    log "Starting new VPN service: $new_service"

    # Start the new container in detached mode
    docker-compose -p "$DOCKER_COMPOSE_PROJECT_NAME" up -d "$new_service"

    # Give the container a moment to initialize
    log "Waiting for '$new_service' to settle..."
    sleep 20

    # Get the container's IP on the public network
    local container_ip
    container_ip=$(docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" "${DOCKER_COMPOSE_PROJECT_NAME}_${new_service}_1")

    if [ -z "$container_ip" ]; then
        log "ERROR: Could not get IP for '$new_service'. Is it running and attached to a network?"
        FAILURE_COUNT=$((FAILURE_COUNT + 1))
        return 1
    fi

    log "Container '$new_service' started with IP: $container_ip"

    # Set up the host routing to use the new container as a gateway
    if ./vpn-manager/setup-routing.sh "$container_ip"; then
        log "Successfully set up routing for $new_service."
        FAILURE_COUNT=0 # Reset failure count on successful switch
    else
        log "ERROR: Failed to set up routing for $new_service."
        FAILURE_COUNT=$((FAILURE_COUNT + 1))
        return 1
    fi
}

# --- Main Loop ---

# Initial startup
start_next_vpn

while true; do
    log "Current VPN: ${VPN_SERVICES[$CURRENT_VPN_INDEX]}. Failure count: $FAILURE_COUNT/$MAX_FAILURES"

    if [ $FAILURE_COUNT -ge $MAX_FAILURES ]; then
        log "ERROR: Maximum failure count reached. Exiting."
        exit 1
    fi

    local current_service=${VPN_SERVICES[$CURRENT_VPN_INDEX]}
    local container_ip
    container_ip=$(docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" "${DOCKER_COMPOSE_PROJECT_NAME}_${current_service}_1")

    if ! ./vpn-manager/health-check.sh "$current_service" "$container_ip"; then
        log "Health check FAILED for $current_service."
        FAILURE_COUNT=$((FAILURE_COUNT + 1))
        start_next_vpn
    else
        log "Health check PASSED for $current_service."
        FAILURE_COUNT=0
    fi

    log "Sleeping for $HEALTH_CHECK_INTERVAL seconds..."
    sleep $HEALTH_CHECK_INTERVAL
done 