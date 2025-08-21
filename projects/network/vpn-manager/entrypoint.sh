#!/bin/bash

# entrypoint.sh - Handle VPN compose file loading and validation

set -e

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] ENTRYPOINT: $1"
}

# Function to load and validate compose file
load_compose_file() {
    local compose_file="$1"
    
    if [ ! -f "$compose_file" ]; then
        log "ERROR: VPN compose file not found: $compose_file"
        exit 1
    fi
    
    log "Loading VPN compose file: $compose_file"
    
    # Read the compose file contents
    VPN_COMPOSE=$(cat "$compose_file")
    export VPN_COMPOSE
    
    # Validate that it's a valid YAML file
    if ! echo "$VPN_COMPOSE" | yq eval '.' >/dev/null 2>&1; then
        log "ERROR: Invalid YAML in compose file: $compose_file"
        exit 1
    fi
    
    # Extract service names (these will be our VPN options)
    local services
    services=$(echo "$VPN_COMPOSE" | yq eval '.services | keys | .[]' - 2>/dev/null || echo "")
    
    if [ -z "$services" ]; then
        log "ERROR: No services found in compose file"
        exit 1
    fi
    
    log "Found VPN services: $(echo "$services" | tr '\n' ' ')"
    
    # Create the compose file in the working directory
    echo "$VPN_COMPOSE" > /app/vpn-compose.yml
    
    # Export the service list for the VPN manager
    export VPN_SERVICE_LIST="$(echo "$services" | tr '\n' ',' | sed 's/,$//')"
    
    log "VPN compose file loaded successfully"
    log "Service list: $VPN_SERVICE_LIST"
}

# Main entrypoint logic
log "VPN Fallback Container Starting"

# Check if VPN_COMPOSE_FILE is provided
if [ -n "$VPN_COMPOSE_FILE" ]; then
    load_compose_file "$VPN_COMPOSE_FILE"
elif [ -n "$VPN_COMPOSE" ]; then
    log "Using VPN_COMPOSE environment variable"
    
    # Validate the compose content
    if ! echo "$VPN_COMPOSE" | yq eval '.' >/dev/null 2>&1; then
        log "ERROR: Invalid YAML in VPN_COMPOSE environment variable"
        exit 1
    fi
    
    # Extract service names
    local services
    services=$(echo "$VPN_COMPOSE" | yq eval '.services | keys | .[]' - 2>/dev/null || echo "")
    
    if [ -z "$services" ]; then
        log "ERROR: No services found in VPN_COMPOSE"
        exit 1
    fi
    
    log "Found VPN services: $(echo "$services" | tr '\n' ' ')"
    
    # Create the compose file
    echo "$VPN_COMPOSE" > /app/vpn-compose.yml
    
    # Export the service list
    export VPN_SERVICE_LIST="$(echo "$services" | tr '\n' ',' | sed 's/,$//')"
    
    log "VPN compose content loaded successfully"
else
    log "ERROR: Either VPN_COMPOSE_FILE or VPN_COMPOSE must be provided"
    log "Usage examples:"
    log "  docker run -v /path/to/compose.yml:/compose.yml -e VPN_COMPOSE_FILE=/compose.yml vpn-fallback"
    log "  docker run -e VPN_COMPOSE='...' vpn-fallback"
    exit 1
fi

# Ensure the compose file exists and is readable
if [ ! -f "/app/vpn-compose.yml" ]; then
    log "ERROR: Failed to create vpn-compose.yml"
    exit 1
fi

log "Entrypoint setup complete. Starting supervisor..."

# Execute the original command (supervisord)
exec "$@" 