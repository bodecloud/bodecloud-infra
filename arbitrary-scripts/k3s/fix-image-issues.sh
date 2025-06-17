#!/bin/bash

echo "Starting cleanup of deployments with bad images..."

# Bad images from our analysis
BAD_IMAGES=(
  "alpine/k8s:1.29.1"
  "filebrowser/filebrowser:latest"
  "twinproduction/gatus:latest"
  "curlimages/curl:latest"
  "prom/node-exporter:latest"
  "ghcr.io/elfhosted/notifiarr:latest"
  "sctx/overseerr:latest"
  "ghcr.io/elfhosted/riven:latest"
  "spoked/riven:latest"
  "ghcr.io/elfhosted/zurg:v0.9.3-hotfix.6"
  "debridmediamanager/zurg-testing:latest"
  "filebrowser/filebrowser:v2.31.2"
)

# First, delete any cronjobs that might be creating problem pods
echo "Deleting all cronjobs that might create problem pods..."
kubectl delete cronjob --all --ignore-not-found=true

# Delete all jobs related to cleanup-transcode specifically
echo "Deleting cleanup-transcode jobs..."
kubectl delete jobs -l job-name=cleanup-transcode --ignore-not-found=true

# Delete deployments with matching images
for img in "${BAD_IMAGES[@]}"; do
  echo "Looking for deployments using image: $img"
  
  # Find deployments using this image
  kubectl get deployments -o json | jq -r --arg img "$img" '.items[] | select(.spec.template.spec.containers[].image | contains($img)) | .metadata.name' | while read deployment; do
    if [ ! -z "$deployment" ]; then
      echo "Deleting deployment: $deployment"
      kubectl delete deployment "$deployment" --ignore-not-found=true
    fi
  done
  
  # Find StatefulSets using this image
  kubectl get statefulsets -o json | jq -r --arg img "$img" '.items[] | select(.spec.template.spec.containers[].image | contains($img)) | .metadata.name' | while read statefulset; do
    if [ ! -z "$statefulset" ]; then
      echo "Deleting StatefulSet: $statefulset"
      kubectl delete statefulset "$statefulset" --ignore-not-found=true
    fi
  done
  
  # Also delete any standalone pods with this image
  echo "Looking for standalone pods using image: $img"
  kubectl get pods -o json | jq -r --arg img "$img" '.items[] | select(.spec.containers[].image | contains($img)) | .metadata.name' | while read pod; do
    if [ ! -z "$pod" ]; then
      echo "Deleting pod: $pod"
      kubectl delete pod "$pod" --force --grace-period=0 --ignore-not-found=true
    fi
  done
done

# Delete pods in ImagePullBackOff state
echo "Deleting pods stuck in ImagePullBackOff state..."
kubectl get pods -A -o json | jq -r '.items[] | select(.status.containerStatuses?[]?.state.waiting?.reason == "ImagePullBackOff") | "\(.metadata.namespace) \(.metadata.name)"' | while read namespace pod; do
  if [ ! -z "$pod" ]; then
    echo "Deleting ImagePullBackOff pod: $pod in namespace: $namespace"
    kubectl delete pod "$pod" -n "$namespace" --force --grace-period=0 --ignore-not-found=true
  fi
done

# Delete failed cleanup-transcode pods specifically
echo "Deleting all cleanup-transcode pods..."
kubectl get pods -A -o json | jq -r '.items[] | select(.metadata.name | startswith("cleanup-transcode")) | "\(.metadata.namespace) \(.metadata.name)"' | while read namespace pod; do
  if [ ! -z "$pod" ]; then
    echo "Deleting cleanup-transcode pod: $pod in namespace: $namespace"
    kubectl delete pod "$pod" -n "$namespace" --force --grace-period=0 --ignore-not-found=true
  fi
done

echo "Cleanup complete. Check for any remaining issues with: kubectl get pods" 