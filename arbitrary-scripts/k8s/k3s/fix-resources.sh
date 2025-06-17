#!/bin/bash

echo "Starting resource limit adjustment..."

# Get a list of all deployments and statefulsets
echo "Scanning for deployments with resource issues..."

# Patch deployments with high CPU limits
kubectl get deployments -A -o json | jq -r '.items[] | select(.spec.template.spec.containers[0].resources.limits.cpu != null) | select(.spec.template.spec.containers[0].resources.limits.cpu | startswith("1") or startswith("2") or startswith("3") or startswith("4") or startswith("5") or startswith("6") or startswith("7") or startswith("8") or startswith("9")) | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    echo "Patching high CPU limits for deployment $ns/$name"
    # Set reasonable CPU limits - 200m instead of multiple cores
    kubectl patch deployment -n $ns $name --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/limits/cpu", "value":"200m"}]'
    # Set reasonable CPU requests - 50m instead of potentially high values
    kubectl patch deployment -n $ns $name --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/requests/cpu", "value":"50m"}]'
done

# Patch statefulsets with high CPU limits
kubectl get statefulsets -A -o json | jq -r '.items[] | select(.spec.template.spec.containers[0].resources.limits.cpu != null) | select(.spec.template.spec.containers[0].resources.limits.cpu | startswith("1") or startswith("2") or startswith("3") or startswith("4") or startswith("5") or startswith("6") or startswith("7") or startswith("8") or startswith("9")) | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    echo "Patching high CPU limits for statefulset $ns/$name"
    # Set reasonable CPU limits - 200m instead of multiple cores
    kubectl patch statefulset -n $ns $name --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/limits/cpu", "value":"200m"}]'
    # Set reasonable CPU requests - 50m instead of potentially high values
    kubectl patch statefulset -n $ns $name --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/requests/cpu", "value":"50m"}]'
done

# Patch deployments with high memory limits (>1Gi)
kubectl get deployments -A -o json | jq -r '.items[] | select(.spec.template.spec.containers[0].resources.limits.memory != null) | select(.spec.template.spec.containers[0].resources.limits.memory | endswith("Gi") or endswith("G")) | .metadata.namespace + " " + .metadata.name' | while read ns name; do
    echo "Patching high memory limits for deployment $ns/$name"
    # Set reasonable memory limits - 512Mi instead of multiple Gi
    kubectl patch deployment -n $ns $name --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/limits/memory", "value":"512Mi"}]'
    # Set reasonable memory requests - 128Mi instead of potentially high values
    kubectl patch deployment -n $ns $name --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/resources/requests/memory", "value":"128Mi"}]'
done

# Fix replica counts for deployments with more than 1 replica
kubectl get deployments -A -o json | jq -r '.items[] | select(.spec.replicas > 1) | .metadata.namespace + " " + .metadata.name + " " + (.spec.replicas | tostring)' | while read ns name replicas; do
    echo "Reducing replicas for deployment $ns/$name from $replicas to 1"
    kubectl scale deployment -n $ns $name --replicas=1
done

echo "Resource adjustments completed." 