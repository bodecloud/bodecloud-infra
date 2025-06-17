#!/bin/bash

# ANSI color codes for colored output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Ensure the script is run as root
if [ "$EUID" -ne 0 ]; then
    echo >&2 "${RED}Error: This script must be run as root. Aborting.${NC}"
    exit 1
fi

# --- Configuration ---
CONTAINER="gluetun-premiumize" # CHANGE THIS to your gluetun container name if different
RESULTS_FILE="vpn_port_scan_results.txt"
START_PORT=${1:-1}
END_PORT=${2:-65535}
STEP=${3:-1}
SCAN_TIMEOUT=${4:-3} # Timeout for nmap scan (seconds)
LISTEN_TIMEOUT=${5:-5} # Max time the listener will run if not killed (seconds)
WAIT_TIME=${6:-0.5} # Time to wait for listener to start before scanning (seconds)

# --- Ensure Tools ---
# Check for required host tools
command -v docker >/dev/null 2>&1 || { echo >&2 "Error: docker command not found. Aborting."; exit 1; }
command -v nmap >/dev/null 2>&1 || { echo >&2 "Error: nmap command not found. Aborting."; exit 1; }
command -v nc >/dev/null 2>&1 || { echo >&2 "Error: nc command not found. Aborting."; exit 1; }
command -v seq >/dev/null 2>&1 || { echo >&2 "Error: seq command not found. Aborting."; exit 1; }
command -v awk >/dev/null 2>&1 || { echo >&2 "Error: awk command not found. Aborting."; exit 1; }
command -v grep >/dev/null 2>&1 || { echo >&2 "Error: grep command not found. Aborting."; exit 1; }
command -v tee >/dev/null 2>&1 || { echo >&2 "Error: tee command not found. Aborting."; exit 1; }

# Ensure netcat, timeout (coreutils), and curl are installed in the container
echo -e "${BLUE}Ensuring netcat, coreutils, and curl are installed in container $CONTAINER...${NC}"
# Combine checks to reduce docker exec calls
docker exec "$CONTAINER" sh -c '
    INSTALL_CMD=""
    command -v nc >/dev/null 2>&1 || INSTALL_CMD="$INSTALL_CMD netcat-openbsd"
    command -v timeout >/dev/null 2>&1 || INSTALL_CMD="$INSTALL_CMD coreutils"
    command -v curl >/dev/null 2>&1 || INSTALL_CMD="$INSTALL_CMD curl"

    if [ -n "$INSTALL_CMD" ]; then
        echo "Attempting to install: $INSTALL_CMD"
        # Prioritize apk for Alpine Linux containers
        if command -v apk >/dev/null 2>&1; then
            apk add --no-cache $INSTALL_CMD
        elif command -v apt-get >/dev/null 2>&1; then
            apt-get update -qq && apt-get install -y --no-install-recommends $INSTALL_CMD
        elif command -v yum >/dev/null 2>&1; then
            yum install -y -q $INSTALL_CMD
        else
            echo >&2 "No package manager found to install: $INSTALL_CMD"
            exit 1
        fi

        # Verify installation
        for cmd in $INSTALL_CMD; do
            command -v "$cmd" >/dev/null 2>&1 || { 
                echo >&2 "Failed to install $cmd"; 
                exit 1; 
            }
        done

        echo "Installation successful."
    else
        echo "Required tools already present."
    fi
'
if [ $? -ne 0 ]; then exit 1; fi
echo -e "${GREEN}Tools check complete in container.${NC}"

# --- Get VPN IP ---
echo -e "${CYAN}Fetching VPN external IP from container $CONTAINER...${NC}"
VPN_IP=$(docker exec "$CONTAINER" curl -s --max-time 10 ifconfig.me) # Increased timeout slightly

# *** CRITICAL CHECK: Validate VPN_IP ***
if ! [[ "$VPN_IP" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo >&2 "${RED}Error: Could not retrieve a valid VPN external IP address. Value received: '$VPN_IP'. Aborting.${NC}"
    exit 1
fi
echo -e "${GREEN}VPN external IP: $VPN_IP${NC}" | tee "$RESULTS_FILE"

# --- Get Container IP ---
echo -e "${CYAN}Fetching internal IP for container $CONTAINER...${NC}"
CONTAINER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$CONTAINER" | head -n 1)

if [ -z "$CONTAINER_IP" ]; then
    echo >&2 "${RED}Error: Could not retrieve an internal IP address for container $CONTAINER. Aborting.${NC}"
    exit 1
fi
echo -e "${GREEN}Container internal IP: $CONTAINER_IP${NC}" | tee -a "$RESULTS_FILE"

# --- Initialize Results File ---
echo -e "${BLUE}Scanning VPN IP $VPN_IP from port $START_PORT to $END_PORT (step: $STEP)${NC}" | tee -a "$RESULTS_FILE"
echo "Scan started: $(date)" | tee -a "$RESULTS_FILE"
echo "----------------------------------------" | tee -a "$RESULTS_FILE"

# --- Main Scan Loop ---
for PORT in $(seq "$START_PORT" "$STEP" "$END_PORT"); do
    echo -e "${YELLOW}-------------------- Testing Port $PORT --------------------${NC}" | tee -a "$RESULTS_FILE"

    # Define the listener command to run inside the container using a loop
    listener_cmd="while true; do echo 'Listening on TCP port $PORT' | nc -l -p $PORT -v; done"

    echo -e "${CYAN}[$(date +%T)] Starting persistent listener on port $PORT inside $CONTAINER...${NC}" | tee -a "$RESULTS_FILE"

    # Start the listener IN THE BACKGROUND on the host, running the command inside the container
    # No timeout here, relying on while true loop and pkill for cleanup
    docker exec "$CONTAINER" sh -c "$listener_cmd" > /dev/null 2>&1 &
    HOST_PID=$! # Get the PID of the 'docker exec' command on the host

    # Give the listener a moment to start up inside the container
    sleep "$WAIT_TIME"

    # Check if the docker exec process is still running (basic check)
    if ps -p $HOST_PID > /dev/null; then
        # Listener process seems okay, proceed with checks:

        # --- Verify Listener Internally (BEFORE) ---
        echo -e "${CYAN}[$(date +%T)] Verifying listener internally (BEFORE) on ${CONTAINER_IP}:${PORT}...${NC}" | tee -a "$RESULTS_FILE"
        timeout 1 nc -z -w 1 "$CONTAINER_IP" "$PORT" >/dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}[$(date +%T)] Internal listener check (BEFORE): successful for port $PORT.${NC}" | tee -a "$RESULTS_FILE"
        else
            # Skip aborting for well-known ports (1-1024)
            if [ "$PORT" -gt 1024 ]; then
                echo >&2 "${RED}[$(date +%T)] CRITICAL: Internal listener check (BEFORE) FAILED for port $PORT. Aborting.${NC}"
                # Cleanup and exit
                if [ -n "$HOST_PID" ]; then
                    kill $HOST_PID >/dev/null 2>&1; wait $HOST_PID >/dev/null 2>&1
                    docker exec "$CONTAINER" pkill -f "nc -l -p $PORT" >/dev/null 2>&1
                fi
                exit 1
            else
                echo >&2 "${YELLOW}[$(date +%T)] WARNING: Internal listener check (BEFORE) FAILED for well-known port $PORT. Continuing scan.${NC}"
            fi
        fi

        # --- External Scan --- 
        echo -e "${CYAN}[$(date +%T)] Scanning external $VPN_IP:$PORT with nmap...${NC}" | tee -a "$RESULTS_FILE"
        NMAP_OUTPUT=$(nmap -Pn --max-retries 0 --host-timeout "${SCAN_TIMEOUT}s" -p "$PORT" "$VPN_IP")

        # Log the raw Nmap output immediately
        echo -e "${MAGENTA}---------- NMAP RAW OUTPUT START -----------${NC}" | tee -a "$RESULTS_FILE"
        echo "$NMAP_OUTPUT" | tee -a "$RESULTS_FILE"
        echo -e "${MAGENTA}----------- NMAP RAW OUTPUT END ------------${NC}" | tee -a "$RESULTS_FILE"

        # --- Verify Listener Internally (AFTER) ---
        echo -e "${CYAN}[$(date +%T)] Verifying listener internally (AFTER) on ${CONTAINER_IP}:${PORT}...${NC}" | tee -a "$RESULTS_FILE"
        timeout 1 nc -z -w 1 "$CONTAINER_IP" "$PORT" >/dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}[$(date +%T)] Internal listener check (AFTER): successful for port $PORT.${NC}" | tee -a "$RESULTS_FILE"
        else
            # Skip aborting for well-known ports (1-1024)
            if [ "$PORT" -gt 1024 ]; then
                echo >&2 "${RED}[$(date +%T)] CRITICAL: Internal listener check (AFTER) FAILED for port $PORT. Aborting.${NC}"
                # Cleanup and exit
                if [ -n "$HOST_PID" ]; then
                    kill $HOST_PID >/dev/null 2>&1; wait $HOST_PID >/dev/null 2>&1
                    docker exec "$CONTAINER" pkill -f "nc -l -p $PORT" >/dev/null 2>&1
                fi
                exit 1
            else
                echo >&2 "${YELLOW}[$(date +%T)] WARNING: Internal listener check (AFTER) FAILED for well-known port $PORT. Continuing scan.${NC}"
            fi
        fi

        # --- Parse Nmap Result & Log Final State --- 
        SCAN_RESULT=$(echo "$NMAP_OUTPUT" | grep "^$PORT/" | awk '{print $2}')
        # Handle cases where nmap doesn't output the state clearly for the port
        if [ -z "$SCAN_RESULT" ]; then
            if echo "$NMAP_OUTPUT" | grep -q "Host seems down"; then
                SCAN_RESULT="host_down"
            elif echo "$NMAP_OUTPUT" | grep -q "Note: Host seems down"; then # Handle slight variation
                SCAN_RESULT="host_down (note)"
            else
                SCAN_RESULT="unknown_nmap_output"
            fi
        fi
        # Log the final parsed Nmap state
        echo -e "${GREEN}[$(date +%T)] Nmap determined state for Port $PORT: $SCAN_RESULT${NC}" | tee -a "$RESULTS_FILE"

    else # Listener process failed to start or exited prematurely
        echo >&2 "${RED}[$(date +%T)] CRITICAL: Listener process failed to start for port $PORT. Aborting.${NC}"
        exit 1
    fi # End ps -p $HOST_PID check

    # --- Cleanup Listener (Remains at the end) ---
    if [ -n "$HOST_PID" ]; then # Check if HOST_PID was actually set
        echo "${YELLOW}[$(date +%T)] Attempting to stop listener for port $PORT (Host PID $HOST_PID)...${NC}" | tee -a "$RESULTS_FILE"
        # Then, forcefully kill any lingering nc process listening on that port INSIDE the container
        # Use pkill to find the specific nc listener process AND the sh process running the loop
        # Kill the shell running the 'while true' loop first
        docker exec "$CONTAINER" pkill -f "sh -c ${listener_cmd}" > /dev/null 2>&1
        # Then kill any remaining nc process just in case
        docker exec "$CONTAINER" pkill -f "nc -l -p $PORT -v" > /dev/null 2>&1
        if [ $? -eq 0 ]; then
             echo -e "${GREEN}[$(date +%T)] Successfully sent kill signal to internal nc process/loop for port $PORT.${NC}" | tee -a "$RESULTS_FILE"
        else
             # This is expected if the process already exited or was killed by the host kill
             echo "${YELLOW}[$(date +%T)] Internal nc process for port $PORT likely already exited or wasn't found (normal).${NC}" | tee -a "$RESULTS_FILE"
        fi
    fi
    # Brief pause to ensure port is released before next iteration
    sleep 0.1
done

echo -e "${BLUE}----------------------------------------${NC}" | tee -a "$RESULTS_FILE"
echo -e "${BLUE}Scan complete: $(date)${NC}" | tee -a "$RESULTS_FILE"
echo -e "${BLUE}Full results saved to $RESULTS_FILE${NC}"