#!/bin/bash

echo "=== Comprehensive Cleanup for Image Pull Issues ==="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if kubectl is working
check_kubectl() {
  if ! kubectl cluster-info >/dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to Kubernetes cluster${NC}"
    echo "Please ensure kubectl is configured and the cluster is accessible"
    exit 1
  fi
  echo -e "${GREEN}✓ Kubernetes cluster is accessible${NC}"
}

# Function to cleanup specific problematic image
cleanup_problematic_images() {
  echo -e "${BLUE}Cleaning up pods with problematic images...${NC}"

  # List of problematic images
  PROBLEMATIC_IMAGES=(
    "ghcr.io/elfhosted/tooling:focal-20240530"
    "alpine/k8s:1.29.1"
    "prom/node-exporter:latest"
    "twinproduction/gatus:latest"
    "curlimages/curl:latest"
    "ghcr.io/elfhosted/notifiarr:latest"
    "sctx/overseerr:latest"
    "ghcr.io/elfhosted/riven:latest"
    "spoked/riven:latest"
    "ghcr.io/elfhosted/zurg:v0.9.3-hotfix.6"
    "debridmediamanager/zurg-testing:latest"
  )

  for image in "${PROBLEMATIC_IMAGES[@]}"; do
    echo -e "${YELLOW}Processing image: $image${NC}"

    # Find and delete pods using this image across all namespaces
    kubectl get pods --all-namespaces -o json |
      jq -r --arg img "$image" '.items[] | select(.spec.containers[]?.image | contains($img)) | "\(.metadata.namespace) \(.metadata.name)"' |
      while read namespace pod; do
        if [ ! -z "$pod" ]; then
          echo -e "${YELLOW}  Deleting pod: $pod in namespace: $namespace${NC}"
          kubectl delete pod "$pod" -n "$namespace" --force --grace-period=0 --ignore-not-found=true
        fi
      done

    # Find and delete deployments using this image
    kubectl get deployments --all-namespaces -o json |
      jq -r --arg img "$image" '.items[] | select(.spec.template.spec.containers[]?.image | contains($img)) | "\(.metadata.namespace) \(.metadata.name)"' |
      while read namespace deployment; do
        if [ ! -z "$deployment" ]; then
          echo -e "${YELLOW}  Deleting deployment: $deployment in namespace: $namespace${NC}"
          kubectl delete deployment "$deployment" -n "$namespace" --ignore-not-found=true
        fi
      done

    # Find and delete StatefulSets using this image
    kubectl get statefulsets --all-namespaces -o json |
      jq -r --arg img "$image" '.items[] | select(.spec.template.spec.containers[]?.image | contains($img)) | "\(.metadata.namespace) \(.metadata.name)"' |
      while read namespace statefulset; do
        if [ ! -z "$statefulset" ]; then
          echo -e "${YELLOW}  Deleting StatefulSet: $statefulset in namespace: $namespace${NC}"
          kubectl delete statefulset "$statefulset" -n "$namespace" --ignore-not-found=true
        fi
      done
  done
}

# Function to cleanup all cleanup-transcode related resources
cleanup_transcode_resources() {
  echo -e "${BLUE}Cleaning up cleanup-transcode resources...${NC}"

  # Delete all cleanup-transcode cronjobs
  echo -e "${YELLOW}Deleting cleanup-transcode cronjobs...${NC}"
  kubectl get cronjobs --all-namespaces -o name | grep cleanup-transcode | xargs -r kubectl delete --ignore-not-found=true

  # Delete all cleanup-transcode jobs
  echo -e "${YELLOW}Deleting cleanup-transcode jobs...${NC}"
  kubectl get jobs --all-namespaces -o name | grep cleanup-transcode | xargs -r kubectl delete --ignore-not-found=true

  # Delete all cleanup-transcode pods
  echo -e "${YELLOW}Deleting cleanup-transcode pods...${NC}"
  kubectl get pods --all-namespaces -o name | grep cleanup-transcode | xargs -r kubectl delete --force --grace-period=0 --ignore-not-found=true
}

# Function to cleanup ImagePullBackOff pods
cleanup_imagepullbackoff_pods() {
  echo -e "${BLUE}Cleaning up ImagePullBackOff pods...${NC}"

  kubectl get pods --all-namespaces -o json |
    jq -r '.items[] | select(.status.containerStatuses[]?.state.waiting?.reason == "ImagePullBackOff" or .status.initContainerStatuses[]?.state.waiting?.reason == "ImagePullBackOff") | "\(.metadata.namespace) \(.metadata.name)"' |
    while read namespace pod; do
      if [ ! -z "$pod" ]; then
        echo -e "${YELLOW}  Deleting ImagePullBackOff pod: $pod in namespace: $namespace${NC}"
        kubectl delete pod "$pod" -n "$namespace" --force --grace-period=0 --ignore-not-found=true
      fi
    done
}

# Function to cleanup ErrImagePull pods
cleanup_errimagepull_pods() {
  echo -e "${BLUE}Cleaning up ErrImagePull pods...${NC}"

  kubectl get pods --all-namespaces -o json |
    jq -r '.items[] | select(.status.containerStatuses[]?.state.waiting?.reason == "ErrImagePull" or .status.initContainerStatuses[]?.state.waiting?.reason == "ErrImagePull") | "\(.metadata.namespace) \(.metadata.name)"' |
    while read namespace pod; do
      if [ ! -z "$pod" ]; then
        echo -e "${YELLOW}  Deleting ErrImagePull pod: $pod in namespace: $namespace${NC}"
        kubectl delete pod "$pod" -n "$namespace" --force --grace-period=0 --ignore-not-found=true
      fi
    done
}

# Function to show status after cleanup
show_status() {
  echo -e "${BLUE}=== Post-cleanup Status ===${NC}"

  echo -e "${YELLOW}Remaining pods with issues:${NC}"
  kubectl get pods --all-namespaces --field-selector=status.phase=Failed 2>/dev/null | head -10 || echo "No failed pods found"

  echo -e "\n${YELLOW}Remaining ImagePullBackOff pods:${NC}"
  kubectl get pods --all-namespaces -o wide | grep ImagePullBackOff | head -5 || echo "No ImagePullBackOff pods found"

  echo -e "\n${YELLOW}Remaining cleanup-transcode resources:${NC}"
  kubectl get pods,jobs,cronjobs --all-namespaces | grep cleanup-transcode | head -5 || echo "No cleanup-transcode resources found"
}

# Main execution
main() {
  echo -e "${GREEN}Starting comprehensive cleanup...${NC}"

  check_kubectl

  cleanup_problematic_images

  cleanup_transcode_resources

  cleanup_imagepullbackoff_pods

  cleanup_errimagepull_pods

  # Wait a moment for resources to be fully deleted
  echo -e "${BLUE}Waiting for cleanup to complete...${NC}"
  sleep 5

  show_status

  echo -e "${GREEN}=== Cleanup Complete ===${NC}"
  echo -e "${BLUE}You may want to check your cluster with:${NC}"
  echo "  kubectl get pods --all-namespaces"
  echo "  kubectl get jobs --all-namespaces"
  echo "  kubectl get cronjobs --all-namespaces"
}

# Run main function
main "$@"
