#!/bin/bash

# Kubernetes Resource Monitoring Script
# Monitors CPU, memory, and pod resource usage across the cluster

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Kubernetes Resource Monitoring ===${NC}"
echo "Timestamp: $(date)"
echo

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}Error: kubectl is not installed or not in PATH${NC}"
    exit 1
fi

# Check cluster connectivity
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}Error: Cannot connect to Kubernetes cluster${NC}"
    exit 1
fi

echo -e "${GREEN}=== Node Resource Usage ===${NC}"
kubectl top nodes 2>/dev/null || echo -e "${YELLOW}Warning: Metrics server not available for node metrics${NC}"
echo

echo -e "${GREEN}=== Pod Resource Usage (Top 20) ===${NC}"
kubectl top pods --all-namespaces --sort-by=memory 2>/dev/null | head -21 || echo -e "${YELLOW}Warning: Metrics server not available for pod metrics${NC}"
echo

echo -e "${GREEN}=== Resource Quotas ===${NC}"
kubectl get resourcequota --all-namespaces -o wide 2>/dev/null || echo "No resource quotas found"
echo

echo -e "${GREEN}=== Pod Status Summary ===${NC}"
echo "Namespace breakdown:"
kubectl get pods --all-namespaces --no-headers | awk '{print $1}' | sort | uniq -c | sort -nr
echo

echo "Pod status summary:"
kubectl get pods --all-namespaces --no-headers | awk '{print $4}' | sort | uniq -c | sort -nr
echo

echo -e "${GREEN}=== Memory and CPU Pressure ===${NC}"
kubectl describe nodes | grep -E "(Name:|Conditions:)" -A 10 | grep -E "(MemoryPressure|DiskPressure|PIDPressure)" || echo "No resource pressure detected"
echo

echo -e "${GREEN}=== Failed/Pending Pods ===${NC}"
kubectl get pods --all-namespaces --field-selector=status.phase!=Running,status.phase!=Succeeded 2>/dev/null || echo "No failed or pending pods"
echo

echo -e "${GREEN}=== Recent Events (Last 10) ===${NC}"
kubectl get events --all-namespaces --sort-by='.lastTimestamp' | tail -10
echo

echo -e "${GREEN}=== Storage Usage ===${NC}"
kubectl get pv,pvc --all-namespaces 2>/dev/null || echo "No persistent volumes found"
echo

# Function to check resource usage percentage
check_resource_usage() {
    local namespace=$1
    echo -e "${BLUE}=== Resource Usage for namespace: $namespace ===${NC}"
    
    # Get resource quota if it exists
    if kubectl get resourcequota -n "$namespace" &>/dev/null; then
        kubectl describe resourcequota -n "$namespace"
    else
        echo "No resource quota set for namespace $namespace"
    fi
    echo
}

# Check specific namespaces
for ns in my-media-stack vpn-gateway; do
    if kubectl get namespace "$ns" &>/dev/null; then
        check_resource_usage "$ns"
    fi
done

echo -e "${GREEN}=== Recommendations ===${NC}"

# Check for high resource usage
echo "Checking for potential issues..."

# Check if any pods are using >80% of their memory limit
echo "Pods potentially hitting memory limits:"
kubectl top pods --all-namespaces --no-headers 2>/dev/null | while read -r namespace name cpu memory; do
    if [[ "$memory" =~ ([0-9]+)Mi ]]; then
        mem_usage=${BASH_REMATCH[1]}
        if [ "$mem_usage" -gt 800 ]; then
            echo -e "${YELLOW}  - $namespace/$name: ${memory} memory usage${NC}"
        fi
    fi
done || echo "Cannot check memory usage - metrics server unavailable"

# Check for pods with restart counts > 5
echo
echo "Pods with high restart counts:"
kubectl get pods --all-namespaces --no-headers | awk '$5 > 5 {print "  - " $1 "/" $2 ": " $5 " restarts"}' || echo "No pods with high restart counts"

echo
echo -e "${GREEN}=== Monitoring Complete ===${NC}"
echo "For continuous monitoring, run: watch -n 30 $0" 