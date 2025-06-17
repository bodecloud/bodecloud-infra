#!/bin/bash

echo "Starting cleanup of error pods and deployment issues..."

# First, find all deployments with more than 1 replica
echo "Scaling down all deployments to 1 replica..."
kubectl get deployments -o json | jq -r '.items[] | select(.spec.replicas > 1) | .metadata.name' | while read deploy; do
  echo "Scaling $deploy to 1 replica"
  kubectl scale deployment $deploy --replicas=1
done

# Delete pods in Error state
echo "Deleting pods in Error state..."
kubectl get pods -o json | jq -r '.items[] | select(.status.phase == "Error" or .status.phase == "Failed") | .metadata.name' | while read pod; do
  echo "Deleting Error pod: $pod"
  kubectl delete pod $pod --force --grace-period=0
done

# Delete all pods in CrashLoopBackOff state
echo "Deleting pods in CrashLoopBackOff state..."
kubectl get pods -o json | jq -r '.items[] | select(.status.containerStatuses[]?.state.waiting.reason == "CrashLoopBackOff") | .metadata.name' | while read pod; do
  echo "Deleting CrashLoopBackOff pod: $pod"
  kubectl delete pod $pod --force --grace-period=0
done

# Delete all completed jobs
echo "Deleting completed jobs..."
kubectl delete jobs --all

# Delete any failed pods
echo "Delete any remaining failed pods..."
kubectl get pods --field-selector=status.phase=Failed --all-namespaces -o name | xargs -r kubectl delete --force --grace-period=0

echo "Cleanup complete. Check cluster status with: kubectl get pods"
