#!/bin/bash
# validate-ha.sh - Validate High Availability and identify Single Points of Failure

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
}

success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] SUCCESS:${NC} $1"
}

# Check cluster health
check_cluster_health() {
    log "🏥 Checking cluster health..."
    
    # Check nodes
    NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
    READY_NODES=$(kubectl get nodes --no-headers | grep -c " Ready ")
    
    if [ "$NODE_COUNT" -lt 3 ]; then
        warn "Only $NODE_COUNT nodes in cluster. Recommended: 3+ for HA"
    else
        success "$NODE_COUNT nodes in cluster (good for HA)"
    fi
    
    if [ "$READY_NODES" -ne "$NODE_COUNT" ]; then
        error "$((NODE_COUNT - READY_NODES)) nodes are not ready!"
    else
        success "All $READY_NODES nodes are ready"
    fi
    
    # Check node resources
    log "Node resource utilization:"
    kubectl top nodes 2>/dev/null || warn "Metrics server not available - cannot check resource usage"
}

# Check for single points of failure
check_single_points_of_failure() {
    log "🔍 Checking for Single Points of Failure..."
    
    SPOF_COUNT=0
    
    # Get all deployments across all namespaces
    while IFS= read -r line; do
        NAMESPACE=$(echo "$line" | awk '{print $1}')
        DEPLOYMENT=$(echo "$line" | awk '{print $2}')
        REPLICAS=$(echo "$line" | awk '{print $3}' | cut -d'/' -f2)
        
        if [ "$REPLICAS" -eq 1 ]; then
            warn "SPOF: $NAMESPACE/$DEPLOYMENT has only 1 replica"
            ((SPOF_COUNT++))
        fi
    done < <(kubectl get deployments --all-namespaces --no-headers 2>/dev/null)
    
    # Check StatefulSets
    while IFS= read -r line; do
        NAMESPACE=$(echo "$line" | awk '{print $1}')
        STATEFULSET=$(echo "$line" | awk '{print $2}')
        REPLICAS=$(echo "$line" | awk '{print $3}' | cut -d'/' -f2)
        
        if [ "$REPLICAS" -eq 1 ]; then
            warn "SPOF: $NAMESPACE/$STATEFULSET (StatefulSet) has only 1 replica"
            ((SPOF_COUNT++))
        fi
    done < <(kubectl get statefulsets --all-namespaces --no-headers 2>/dev/null)
    
    if [ "$SPOF_COUNT" -eq 0 ]; then
        success "No single points of failure detected in deployments/statefulsets"
    else
        error "Found $SPOF_COUNT potential single points of failure"
    fi
}

# Check pod distribution across nodes
check_pod_distribution() {
    log "📊 Checking pod distribution across nodes..."
    
    # Get pod distribution
    log "Pods per node:"
    kubectl get pods --all-namespaces --field-selector=status.phase=Running -o wide --no-headers | \
        awk '{print $8}' | sort | uniq -c | sort -nr
    
    # Check for anti-affinity rules
    log "Checking anti-affinity configurations..."
    ANTI_AFFINITY_COUNT=$(kubectl get deployments --all-namespaces -o yaml | grep -c "podAntiAffinity" || true)
    if [ "$ANTI_AFFINITY_COUNT" -eq 0 ]; then
        warn "No pod anti-affinity rules found - pods may cluster on same nodes"
    else
        success "$ANTI_AFFINITY_COUNT deployments have anti-affinity rules"
    fi
}

# Check PodDisruptionBudgets
check_pod_disruption_budgets() {
    log "🛡️ Checking PodDisruptionBudgets..."
    
    PDB_COUNT=$(kubectl get pdb --all-namespaces --no-headers 2>/dev/null | wc -l)
    if [ "$PDB_COUNT" -eq 0 ]; then
        warn "No PodDisruptionBudgets found - cluster updates may cause downtime"
    else
        success "$PDB_COUNT PodDisruptionBudgets configured"
        kubectl get pdb --all-namespaces
    fi
}

# Check storage high availability
check_storage_ha() {
    log "💾 Checking storage high availability..."
    
    # Check storage classes
    log "Available storage classes:"
    kubectl get storageclass
    
    # Check for distributed storage
    LONGHORN_PODS=$(kubectl get pods -n longhorn-system --no-headers 2>/dev/null | wc -l || echo "0")
    ROOK_PODS=$(kubectl get pods -n rook-ceph --no-headers 2>/dev/null | wc -l || echo "0")
    
    if [ "$LONGHORN_PODS" -gt 0 ]; then
        success "Longhorn distributed storage detected ($LONGHORN_PODS pods)"
    elif [ "$ROOK_PODS" -gt 0 ]; then
        success "Rook-Ceph distributed storage detected ($ROOK_PODS pods)"
    else
        warn "No distributed storage system detected - consider Longhorn or Rook-Ceph"
    fi
    
    # Check PVCs
    PVC_COUNT=$(kubectl get pvc --all-namespaces --no-headers 2>/dev/null | wc -l)
    log "Total PersistentVolumeClaims: $PVC_COUNT"
}

# Check network high availability
check_network_ha() {
    log "🌐 Checking network high availability..."
    
    # Check for load balancer
    METALLB_PODS=$(kubectl get pods -n metallb-system --no-headers 2>/dev/null | wc -l || echo "0")
    if [ "$METALLB_PODS" -gt 0 ]; then
        success "MetalLB load balancer detected ($METALLB_PODS pods)"
    else
        warn "No MetalLB load balancer detected - external access may be limited"
    fi
    
    # Check ingress controllers
    INGRESS_CONTROLLERS=$(kubectl get pods --all-namespaces -l app.kubernetes.io/name=traefik --no-headers 2>/dev/null | wc -l || echo "0")
    if [ "$INGRESS_CONTROLLERS" -gt 1 ]; then
        success "Multiple ingress controller replicas detected ($INGRESS_CONTROLLERS)"
    elif [ "$INGRESS_CONTROLLERS" -eq 1 ]; then
        warn "Only 1 ingress controller replica - consider scaling up"
    else
        warn "No Traefik ingress controllers detected"
    fi
}

# Check VPN gateway system HA
check_vpn_gateway_ha() {
    log "🔒 Checking VPN Gateway system HA..."
    
    # Check if VPN gateway namespace exists
    if ! kubectl get namespace vpn-gateway &>/dev/null; then
        warn "VPN gateway namespace not found - system not deployed"
        return
    fi
    
    # Check controller
    CONTROLLER_REPLICAS=$(kubectl get deployment vpn-gateway-controller -n vpn-gateway -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0")
    if [ "$CONTROLLER_REPLICAS" -eq 1 ]; then
        success "VPN gateway controller running (single replica by design)"
    else
        warn "VPN gateway controller not running or misconfigured"
    fi
    
    # Check DaemonSet
    EXPECTED_DAEMONSET=$(kubectl get nodes --no-headers | wc -l)
    ACTUAL_DAEMONSET=$(kubectl get pods -n vpn-gateway -l app=pod-gateway --field-selector=status.phase=Running --no-headers 2>/dev/null | wc -l || echo "0")
    
    if [ "$ACTUAL_DAEMONSET" -eq "$EXPECTED_DAEMONSET" ]; then
        success "Pod-gateway DaemonSet running on all $ACTUAL_DAEMONSET nodes"
    else
        warn "Pod-gateway DaemonSet only running on $ACTUAL_DAEMONSET/$EXPECTED_DAEMONSET nodes"
    fi
    
    # Check available VPN gateways
    AVAILABLE_GATEWAYS=$(kubectl get deployments -n vpn-gateway -l component=vpn-gateway --no-headers 2>/dev/null | wc -l || echo "0")
    ACTIVE_GATEWAYS=$(kubectl get pods -n vpn-gateway -l component=vpn-gateway --field-selector=status.phase=Running --no-headers 2>/dev/null | wc -l || echo "0")
    
    log "VPN Gateways: $ACTIVE_GATEWAYS active out of $AVAILABLE_GATEWAYS available"
    
    if [ "$AVAILABLE_GATEWAYS" -gt 5 ]; then
        success "Multiple VPN gateway options available for failover"
    else
        warn "Limited VPN gateway options - consider adding more providers"
    fi
}

# Check monitoring and alerting
check_monitoring() {
    log "📊 Checking monitoring and alerting..."
    
    # Check for Prometheus
    PROMETHEUS_PODS=$(kubectl get pods --all-namespaces -l app=prometheus --no-headers 2>/dev/null | wc -l || echo "0")
    if [ "$PROMETHEUS_PODS" -gt 0 ]; then
        success "Prometheus monitoring detected ($PROMETHEUS_PODS pods)"
    else
        warn "No Prometheus monitoring detected"
    fi
    
    # Check for Grafana
    GRAFANA_PODS=$(kubectl get pods --all-namespaces -l app=grafana --no-headers 2>/dev/null | wc -l || echo "0")
    if [ "$GRAFANA_PODS" -gt 0 ]; then
        success "Grafana dashboards detected ($GRAFANA_PODS pods)"
    else
        warn "No Grafana dashboards detected"
    fi
    
    # Check for AlertManager
    ALERTMANAGER_PODS=$(kubectl get pods --all-namespaces -l app=alertmanager --no-headers 2>/dev/null | wc -l || echo "0")
    if [ "$ALERTMANAGER_PODS" -gt 0 ]; then
        success "AlertManager detected ($ALERTMANAGER_PODS pods)"
    else
        warn "No AlertManager detected - no alerting configured"
    fi
}

# Generate HA recommendations
generate_recommendations() {
    log "💡 Generating High Availability recommendations..."
    
    echo
    echo "=== HIGH AVAILABILITY RECOMMENDATIONS ==="
    echo
    
    # Node recommendations
    NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
    if [ "$NODE_COUNT" -lt 3 ]; then
        echo "🔴 CRITICAL: Add more nodes (current: $NODE_COUNT, recommended: 3+)"
    fi
    
    # Deployment recommendations
    SINGLE_REPLICA_DEPLOYMENTS=$(kubectl get deployments --all-namespaces --no-headers | awk '{print $4}' | grep -c "1/1" || true)
    if [ "$SINGLE_REPLICA_DEPLOYMENTS" -gt 0 ]; then
        echo "🟡 MEDIUM: Scale up $SINGLE_REPLICA_DEPLOYMENTS single-replica deployments"
    fi
    
    # Storage recommendations
    LONGHORN_PODS=$(kubectl get pods -n longhorn-system --no-headers 2>/dev/null | wc -l || echo "0")
    ROOK_PODS=$(kubectl get pods -n rook-ceph --no-headers 2>/dev/null | wc -l || echo "0")
    if [ "$LONGHORN_PODS" -eq 0 ] && [ "$ROOK_PODS" -eq 0 ]; then
        echo "🟡 MEDIUM: Deploy distributed storage (Longhorn or Rook-Ceph)"
    fi
    
    # Monitoring recommendations
    PROMETHEUS_PODS=$(kubectl get pods --all-namespaces -l app=prometheus --no-headers 2>/dev/null | wc -l || echo "0")
    if [ "$PROMETHEUS_PODS" -eq 0 ]; then
        echo "🟡 MEDIUM: Deploy monitoring stack (Prometheus + Grafana)"
    fi
    
    # PDB recommendations
    PDB_COUNT=$(kubectl get pdb --all-namespaces --no-headers 2>/dev/null | wc -l)
    if [ "$PDB_COUNT" -eq 0 ]; then
        echo "🟡 MEDIUM: Configure PodDisruptionBudgets for critical services"
    fi
    
    echo
    echo "=== NEXT STEPS ==="
    echo "1. Address CRITICAL issues first"
    echo "2. Scale up single-replica deployments"
    echo "3. Configure pod anti-affinity rules"
    echo "4. Set up distributed storage"
    echo "5. Deploy monitoring and alerting"
    echo "6. Test failover scenarios"
    echo
}

# Main function
main() {
    log "🚀 Starting High Availability Validation"
    echo
    
    check_cluster_health
    echo
    check_single_points_of_failure
    echo
    check_pod_distribution
    echo
    check_pod_disruption_budgets
    echo
    check_storage_ha
    echo
    check_network_ha
    echo
    check_vpn_gateway_ha
    echo
    check_monitoring
    echo
    generate_recommendations
    
    log "✅ High Availability validation complete"
}

# Run main function
main "$@" 