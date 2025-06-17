#!/bin/bash
# deploy-ha-stack.sh - Deploy entire stack with High Availability

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    log "🔍 Checking prerequisites..."
    
    if ! command -v kubectl &> /dev/null; then
        error "kubectl is not installed or not in PATH"
    fi
    
    if ! kubectl cluster-info &> /dev/null; then
        error "Cannot connect to Kubernetes cluster"
    fi
    
    # Check if we have multiple nodes for HA
    NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
    if [ "$NODE_COUNT" -lt 3 ]; then
        warn "Only $NODE_COUNT nodes detected. For true HA, you need at least 3 nodes."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log "✅ $NODE_COUNT nodes detected - good for HA"
    fi
}

# Deploy VPN Gateway System with HA
deploy_vpn_system() {
    log "🔒 Deploying VPN Gateway System with High Availability..."
    
    # Deploy the pod-gateway system
    kubectl apply -f k8s/my-media-stack/vpn-gateway-pod-gateway.yaml
    
    # Deploy all VPN gateways
    kubectl apply -f k8s/my-media-stack/vpn-gateways-all.yaml
    
    # Wait for controller to be ready
    log "⏳ Waiting for VPN controller to be ready..."
    kubectl wait --for=condition=available deployment/vpn-gateway-controller -n vpn-gateway --timeout=300s
    
    # Wait for DaemonSet to be ready
    log "⏳ Waiting for pod-gateway DaemonSet to be ready..."
    kubectl rollout status daemonset/pod-gateway -n vpn-gateway --timeout=300s
    
    log "✅ VPN Gateway System deployed successfully"
}

# Validate HA deployment
validate_ha_deployment() {
    log "🔍 Validating High Availability deployment..."
    
    # Check VPN gateway system
    log "Checking VPN gateway system..."
    kubectl get pods -n vpn-gateway
    
    # Check if controller is running
    CONTROLLER_PODS=$(kubectl get pods -n vpn-gateway -l app=vpn-gateway-controller --field-selector=status.phase=Running --no-headers | wc -l)
    if [ "$CONTROLLER_PODS" -eq 0 ]; then
        error "VPN gateway controller is not running"
    fi
    
    # Check if DaemonSet is running on all nodes
    EXPECTED_DAEMONSET_PODS=$(kubectl get nodes --no-headers | wc -l)
    ACTUAL_DAEMONSET_PODS=$(kubectl get pods -n vpn-gateway -l app=pod-gateway --field-selector=status.phase=Running --no-headers | wc -l)
    
    if [ "$ACTUAL_DAEMONSET_PODS" -ne "$EXPECTED_DAEMONSET_PODS" ]; then
        warn "Pod-gateway DaemonSet not running on all nodes ($ACTUAL_DAEMONSET_PODS/$EXPECTED_DAEMONSET_PODS)"
    else
        log "✅ Pod-gateway DaemonSet running on all $ACTUAL_DAEMONSET_PODS nodes"
    fi
    
    # Check which VPN gateway is active
    log "Checking active VPN gateway..."
    sleep 10  # Give controller time to start a gateway
    ACTIVE_GATEWAYS=$(kubectl get pods -n vpn-gateway -l component=vpn-gateway --field-selector=status.phase=Running --no-headers | wc -l)
    log "Active VPN gateways: $ACTIVE_GATEWAYS"
}

# Test VPN connectivity
test_vpn_connectivity() {
    log "🧪 Testing VPN connectivity..."
    
    # Deploy test pod
    cat << 'EOF' | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: vpn-test
  namespace: default
spec:
  containers:
  - name: test
    image: curlimages/curl:latest
    command: ["sleep", "3600"]
  restartPolicy: Never
EOF

    # Wait for test pod to be ready
    kubectl wait --for=condition=ready pod/vpn-test --timeout=60s
    
    # Test external connectivity
    log "Testing external connectivity through VPN..."
    if kubectl exec vpn-test -- timeout 10 curl -s http://httpbin.org/ip > /tmp/vpn-test-result.json; then
        EXTERNAL_IP=$(cat /tmp/vpn-test-result.json | grep -o '"origin": "[^"]*"' | cut -d'"' -f4)
        log "✅ VPN connectivity working - External IP: $EXTERNAL_IP"
    else
        error "❌ VPN connectivity test failed"
    fi
    
    # Cleanup test pod
    kubectl delete pod vpn-test --ignore-not-found=true
    rm -f /tmp/vpn-test-result.json
}

# Main deployment function
main() {
    log "🚀 Starting High Availability VPN Gateway Deployment"
    
    check_prerequisites
    deploy_vpn_system
    validate_ha_deployment
    test_vpn_connectivity
    
    log "🎉 High Availability VPN Gateway deployment complete!"
    echo
    log "📋 Summary:"
    log "  • VPN Gateway Controller: Managing failover automatically"
    log "  • Pod-Gateway DaemonSet: Running on all nodes for traffic interception"
    log "  • VPN Gateways: Multiple providers with priority-based failover"
    echo
    log "🔍 Useful commands:"
    log "  • Check VPN system status: kubectl get pods -n vpn-gateway"
    log "  • View controller logs: kubectl logs -f -n vpn-gateway deployment/vpn-gateway-controller"
    log "  • View active gateway: kubectl logs -n vpn-gateway deployment/vpn-gateway-controller | grep 'Active gateway'"
    log "  • Test connectivity: kubectl run test --rm -it --image=curlimages/curl -- curl http://httpbin.org/ip"
    echo
    log "🎯 Next steps:"
    log "  • Deploy your media stack services with HA patterns"
    log "  • Set up distributed storage (Longhorn/Rook-Ceph)"
    log "  • Configure monitoring and alerting"
    log "  • Test failover scenarios"
}

# Run main function
main "$@" 