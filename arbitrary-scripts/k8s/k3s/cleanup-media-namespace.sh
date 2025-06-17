#!/bin/bash

echo "Cleaning up the media-services namespace..."

# Delete all failed/error pods
echo "Deleting error pods in media-services namespace..."
kubectl delete pod -n media-services --field-selector=status.phase=Failed
kubectl delete pod -n media-services --field-selector=status.phase=Error
kubectl delete pod -n media-services --field-selector=status.phase=Completed

# Handle the ContainerStatusUnknown and ImagePullBackOff pods with force delete
echo "Force deleting pods with ContainerStatusUnknown and ImagePullBackOff status..."
kubectl get pods -n media-services | grep "ContainerStatusUnknown\|ImagePullBackOff" | awk '{print $1}' | xargs -I {} kubectl delete pod -n media-services {} --force --grace-period=0

# Scale down the problematic deployments
echo "Scaling down problematic deployments in media-services namespace..."
kubectl scale deployment -n media-services plex --replicas=0
kubectl scale deployment -n media-services sonarr --replicas=0
kubectl scale deployment -n media-services tautulli --replicas=0

# Final check
echo "Current status of media-services namespace:"
kubectl get pods -n media-services 