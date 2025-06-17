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

# Override port selection: support single port or range spec in first argument
if [ "$#" -eq 1 ]; then
    if [[ "$1" =~ ^([0-9]+)-([0-9]+)$ ]]; then
        START_PORT="${BASH_REMATCH[1]}"
        END_PORT="${BASH_REMATCH[2]}"
    elif [[ "$1" =~ ^[0-9]+$ ]]; then
        START_PORT="$1"
        END_PORT="$1"
    fi
fi

# --- Ensure Tools ---
# Check for required host tools
command -v docker >/dev/null 2>&1 || { echo >&2 "Error: docker command not found. Aborting."; exit 1; }
command -v nmap >/dev/null 2>&1 || { echo >&2 "Error: nmap command not found. Aborting."; exit 1; }
command -v nc >/dev/null 2>&1 || { echo >&2 "Error: nc command not found. Aborting."; exit 1; }
command -v seq >/dev/null 2>&1 || { echo >&2 "Error: seq command not found. Aborting."; exit 1; }
command -v awk >/dev/null 2>&1 || { echo >&2 "Error: awk command not found. Aborting."; exit 1; }
command -v grep >/dev/null 2>&1 || { echo >&2 "Error: grep command not found. Aborting."; exit 1; }
command -v tee >/dev/null 2>&1 || { echo >&2 "Error: tee command not found. Aborting."; exit 1; }

# Ensure nmap is installed in the container
echo -e "${BLUE}Ensuring nmap is installed in container $CONTAINER...${NC}"
docker exec "$CONTAINER" sh -c '
    if ! command -v nmap >/dev/null 2>&1; then
        if command -v apk >/dev/null 2>&1; then
            apk add --no-cache nmap
        elif command -v apt-get >/dev/null 2>&1; then
            apt-get update -qq && apt-get install -y --no-install-recommends nmap
        elif command -v yum >/dev/null 2>&1; then
            yum install -y -q nmap
        else
            echo >&2 "${RED}No package manager found to install nmap${NC}"
            exit 1
        fi
    fi
'
if [ $? -ne 0 ]; then exit 1; fi
echo -e "${GREEN}Tools check complete in container.${NC}"

# --- Get HOST External IP ---
echo -e "${CYAN}Fetching host external IP...${NC}"
HOST_IP=$(curl -s --max-time 10 ifconfig.me)
# Validate HOST_IP
if ! [[ "$HOST_IP" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo >&2 "${YELLOW}Warning: Could not retrieve a valid host IP, defaulting to 127.0.0.1${NC}"
    HOST_IP="127.0.0.1"
fi
echo -e "${GREEN}Host external IP: $HOST_IP${NC}" | tee "$RESULTS_FILE"

# --- Initialize Results File ---
echo -e "${BLUE}Scanning Host IP $HOST_IP from port $START_PORT to $END_PORT (step: $STEP)${NC}" | tee -a "$RESULTS_FILE"
echo "Scan started: $(date)" | tee -a "$RESULTS_FILE"
echo "----------------------------------------" | tee -a "$RESULTS_FILE"

# --- Main Scan Loop ---
for PORT in $(seq "$START_PORT" "$STEP" "$END_PORT"); do
    echo -e "${YELLOW}-------------------- Testing Port $PORT --------------------${NC}" | tee -a "$RESULTS_FILE"

    # Define the host listener command
    listener_cmd="while true; do echo 'Listening on TCP port $PORT' | nc -l -p $PORT -v; done"
    # Check if host already has a service listening
    if nc -z -w1 127.0.0.1 $PORT >/dev/null 2>&1; then
        echo -e "${CYAN}[$(date +%T)] Found existing host service on port $PORT, skipping listener spawn${NC}" | tee -a "$RESULTS_FILE"
        HOST_PID=""
    else
        echo -e "${CYAN}[$(date +%T)] Starting netcat listener on host for port $PORT...${NC}" | tee -a "$RESULTS_FILE"
        sh -c "$listener_cmd" > /dev/null 2>&1 &
        HOST_PID=$!
        sleep "$WAIT_TIME"
    fi

    # Verify host listener
    if [ -n "$HOST_PID" ] && ! ps -p $HOST_PID >/dev/null; then
        echo >&2 "${RED}[$(date +%T)] Host listener failed to start for port $PORT. Continuing${NC}"
        HOST_PID=""
    fi

    # --- Verify host listener (BEFORE) ---
    if nc -z -w1 127.0.0.1 $PORT >/dev/null 2>&1; then
        echo -e "${GREEN}[$(date +%T)] Host listener active for port $PORT.${NC}" | tee -a "$RESULTS_FILE"
    else
        echo >&2 "${RED}[$(date +%T)] Host listener not active for port $PORT. Continuing${NC}"
    fi

    # --- Scan host from inside container --- 
    echo -e "${CYAN}[$(date +%T)] Scanning host $HOST_IP:$PORT with nmap from container...${NC}" | tee -a "$RESULTS_FILE"
    NMAP_OUTPUT=$(docker exec "$CONTAINER" nmap -Pn --max-retries 0 --host-timeout "${SCAN_TIMEOUT}s" -p "$PORT" "$HOST_IP")

    # Log the raw Nmap output immediately
    echo -e "${MAGENTA}---------- NMAP RAW OUTPUT START -----------${NC}" | tee -a "$RESULTS_FILE"
    echo "$NMAP_OUTPUT" | tee -a "$RESULTS_FILE"
    echo -e "${MAGENTA}----------- NMAP RAW OUTPUT END ------------${NC}" | tee -a "$RESULTS_FILE"

    # --- Verify host listener (AFTER) ---
    if nc -z -w1 127.0.0.1 $PORT >/dev/null 2>&1; then
        echo -e "${GREEN}[$(date +%T)] Host listener still active for port $PORT.${NC}" | tee -a "$RESULTS_FILE"
    else
        echo >&2 "${YELLOW}[$(date +%T)] Host listener no longer active for port $PORT.${NC}"
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

    # --- Cleanup Listener (Remains at the end) ---
    if [ -n "$HOST_PID" ]; then
        echo "[$(date +%T)] Scheduling listener cleanup for port $PORT (Host PID $HOST_PID)..." | tee -a "$RESULTS_FILE"
        {
            # Kill host-side listener process only
            kill "$HOST_PID" > /dev/null 2>&1
            wait "$HOST_PID" > /dev/null 2>&1 || true
            echo -e "${GREEN}[$(date +%T)] Host listener cleanup completed for port $PORT.${NC}" | tee -a "$RESULTS_FILE"
        } &
    fi
    # Continue immediately without blocking
    sleep 0.1
done

echo -e "${BLUE}----------------------------------------${NC}" | tee -a "$RESULTS_FILE"
echo -e "${BLUE}Scan complete: $(date)${NC}" | tee -a "$RESULTS_FILE"
echo -e "${BLUE}Full results saved to $RESULTS_FILE${NC}"