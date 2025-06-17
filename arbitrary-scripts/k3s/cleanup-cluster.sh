#!/bin/bash

echo "Starting cluster cleanup..."

# 1. Delete pods in failed states
echo "Deleting pods in Error state..."
kubectl get pods --all-namespaces -o json | jq -r '.items[] | select(.status.phase == "Error") | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    kubectl delete pod -n $ns $name
done

echo "Deleting pods in ImagePullBackOff state..."
kubectl get pods --all-namespaces -o json | jq -r '.items[] | select(.status.containerStatuses[]?.state.waiting.reason == "ImagePullBackOff") | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    kubectl delete pod -n $ns $name
done

echo "Deleting pods in ContainerStatusUnknown state..."
kubectl get pods --all-namespaces -o json | jq -r '.items[] | select(.status.phase == "Unknown") | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    kubectl delete pod -n $ns $name
done

echo "Deleting evicted pods..."
kubectl get pods --all-namespaces -o json | jq -r '.items[] | select(.status.reason == "Evicted") | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    kubectl delete pod -n $ns $name
done

# 2. Delete duplicate pods for key services (radarr, sonarr, etc.)
echo "Deleting duplicate media service pods..."
for app in radarr sonarr lidarr prowlarr jellyfin plex; do
    echo "Cleaning up $app duplicates..."
    # Keep only the newest pod for each app
    kubectl get pods --all-namespaces -o json | jq -r ".items[] | select(.metadata.name | contains(\"$app\")) | .metadata.namespace + \" \" + .metadata.name + \" \" + .metadata.creationTimestamp" | sort -k3 | head -n -1 | while read ns name timestamp; do
        kubectl delete pod -n $ns $name
    done
done

# 3. Delete PVCs using non-existent storage classes
echo "Deleting PVCs with non-existent storage classes..."
kubectl get pvc --all-namespaces -o json | jq -r '.items[] | select(.spec.storageClassName == "ceph-block-ssd") | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    kubectl delete pvc -n $ns $name
done

echo "Cleanup completed." 