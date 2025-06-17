#!/bin/bash

# Simple test script for WARP routing functionality
# This script runs basic tests to verify routing is working

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

echo "=========================================="
echo "WARP Routing Test Script"
echo "=========================================="
echo

# Check if warpnet exists
if ! docker network inspect warpnet >/dev/null 2>&1; then
    log_error "warpnet Docker network not found. Please run setup_warp_routing.sh first."
    exit 1
fi

log_info "Testing routing functionality..."
echo

# Test 1: Get host external IP
log_info "Getting host external IP..."
HOST_IP=$(timeout 10 curl -s ifconfig.me 2>/dev/null || echo "FAILED")
if [[ "$HOST_IP" == "FAILED" ]]; then
    log_error "Could not get host external IP"
    exit 1
fi
log_info "Host external IP: $HOST_IP"
echo

# Test 2: Get container external IP through warpnet
log_info "Testing container routing through warpnet..."
CONTAINER_IP=$(docker run --rm --network=warpnet alpine:latest sh -c '
    apk add --no-cache curl >/dev/null 2>&1 || exit 1
    curl -s --max-time 15 ifconfig.me 2>/dev/null || echo "FAILED"
' 2>/dev/null)

if [[ "$CONTAINER_IP" == "FAILED" || -z "$CONTAINER_IP" ]]; then
    log_error "Could not get container external IP through warpnet"
    log_error "This indicates routing is not working correctly"
    exit 1
fi
log_info "Container external IP: $CONTAINER_IP"
echo

# Test 3: Compare IPs
log_info "Comparing external IPs..."
if [[ "$HOST_IP" == "$CONTAINER_IP" ]]; then
    log_error "Host and container have the same external IP!"
    log_error "This indicates that routing through WARP is NOT working"
    echo
    log_info "Host IP:      $HOST_IP"
    log_info "Container IP: $CONTAINER_IP"
    exit 1
else
    log_success "Host and container have different external IPs!"
    log_success "Routing through WARP is working correctly"
    echo
    log_info "Host IP:      $HOST_IP"
    log_info "Container IP: $CONTAINER_IP (via WARP)"
fi
echo

# Test 4: Test multiple containers get different internal IPs
log_info "Testing internal IP assignment..."
CONTAINER_1_INTERNAL=$(docker run -d --network=warpnet alpine:latest sleep 10)
CONTAINER_2_INTERNAL=$(docker run -d --network=warpnet alpine:latest sleep 10)

IP_1=$(docker inspect "$CONTAINER_1_INTERNAL" -f '{{.NetworkSettings.Networks.warpnet.IPAddress}}')
IP_2=$(docker inspect "$CONTAINER_2_INTERNAL" -f '{{.NetworkSettings.Networks.warpnet.IPAddress}}')

docker rm -f "$CONTAINER_1_INTERNAL" "$CONTAINER_2_INTERNAL" >/dev/null 2>&1

if [[ "$IP_1" != "$IP_2" && -n "$IP_1" && -n "$IP_2" ]]; then
    log_success "Containers receive unique internal IPs"
    log_info "Container 1: $IP_1"
    log_info "Container 2: $IP_2"
else
    log_warning "Issue with internal IP assignment"
    log_info "Container 1: $IP_1"
    log_info "Container 2: $IP_2"
fi
echo

# Test 5: Test DNS resolution
log_info "Testing DNS resolution through warpnet..."
DNS_TEST=$(docker run --rm --network=warpnet alpine:latest sh -c '
    nslookup google.com >/dev/null 2>&1 && echo "SUCCESS" || echo "FAILED"
' 2>/dev/null)

if [[ "$DNS_TEST" == "SUCCESS" ]]; then
    log_success "DNS resolution working correctly"
else
    log_warning "DNS resolution may have issues"
fi
echo

# Summary
echo "=========================================="
echo "TEST SUMMARY"
echo "=========================================="
log_success "✓ WARP routing is working correctly!"
log_info "Containers on warpnet are routing through WARP"
log_info "External traffic shows different IP than host"
echo

log_info "You can now use warpnet for containers that need WARP routing:"
log_info "docker run --rm --network=warpnet alpine sh -c 'apk add curl && curl ifconfig.me'"
echo 