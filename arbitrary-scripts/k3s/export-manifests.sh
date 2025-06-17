#!/bin/bash

# Script to export all Kubernetes resources to YAML manifests
# Based on information from various sources

# Set the output directory
OUTPUT_DIR="$(pwd)"
NAMESPACED_DIR="${OUTPUT_DIR}/namespaced"
CLUSTER_DIR="${OUTPUT_DIR}/cluster"
KUBECONFIG_DIR="${OUTPUT_DIR}/kubeconfig"
HELM_DIR="${OUTPUT_DIR}/helm-releases"
KUSTOMIZE_DIR="${OUTPUT_DIR}/kustomize"

# Create directories
mkdir -p "${NAMESPACED_DIR}"
mkdir -p "${CLUSTER_DIR}"
mkdir -p "${KUBECONFIG_DIR}"
mkdir -p "${HELM_DIR}"
mkdir -p "${KUSTOMIZE_DIR}"

echo "Exporting kubeconfig..."
# Export kubeconfig (sanitized)
kubectl config view --minify --flatten > "${KUBECONFIG_DIR}/config"
echo "Kubeconfig exported to ${KUBECONFIG_DIR}/config"

echo "Exporting Helm releases..."
# Export Helm releases if helm is installed
if command -v helm &> /dev/null; then
  for release in $(helm list --all-namespaces --short); do
    namespace=$(helm list --all-namespaces | grep $release | awk '{print $2}')
    echo "Exporting Helm release: $release in namespace: $namespace"
    mkdir -p "${HELM_DIR}/${namespace}"
    helm get manifest $release -n $namespace > "${HELM_DIR}/${namespace}/${release}.yaml"
    helm get values $release -n $namespace -o yaml > "${HELM_DIR}/${namespace}/${release}-values.yaml"
  done
else
  echo "Helm not installed, skipping Helm release export"
fi

echo "Exporting namespaced resources..."
# Export namespaced resources
for namespace in $(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}'); do
  echo "Processing namespace: $namespace"
  mkdir -p "${NAMESPACED_DIR}/${namespace}"
  
  # Get all resource types available in this namespace
  for resource in $(kubectl api-resources --namespaced=true --verbs=list -o name); do
    # Skip some resources that are typically not needed or problematic
    if [[ "$resource" == "events" || "$resource" == "events.events.k8s.io" ]]; then
      continue
    fi
    
    echo "  Exporting $resource"
    # Get all resources of this type in the namespace
    if kubectl get $resource -n $namespace -o name 2>/dev/null | grep -v "No resources found" > /dev/null; then
      for resourceitem in $(kubectl get $resource -n $namespace -o name 2>/dev/null); do
        resourcename=$(echo $resourceitem | cut -d/ -f2)
        mkdir -p "${NAMESPACED_DIR}/${namespace}/${resource}"
        # Export the resource definition, removing cluster-specific fields
        kubectl get $resource $resourcename -n $namespace -o yaml | \
          grep -v "^\s*creationTimestamp:" | \
          grep -v "^\s*resourceVersion:" | \
          grep -v "^\s*selfLink:" | \
          grep -v "^\s*uid:" | \
          grep -v "^\s*status:" | \
          grep -v "^\s*generation:" | \
          grep -v "^\s*managedFields:" > "${NAMESPACED_DIR}/${namespace}/${resource}/${resourcename}.yaml"
      done
    fi
  done
done

echo "Exporting cluster-wide resources..."
# Export cluster-wide resources
for resource in $(kubectl api-resources --namespaced=false --verbs=list -o name); do
  # Skip some resources that are typically not needed or problematic
  if [[ "$resource" == "componentstatuses" || "$resource" == "nodes" ]]; then
    continue
  fi
  
  echo "  Exporting $resource"
  mkdir -p "${CLUSTER_DIR}/${resource}"
  
  # Get all resources of this type
  if kubectl get $resource -o name 2>/dev/null | grep -v "No resources found" > /dev/null; then
    for resourceitem in $(kubectl get $resource -o name 2>/dev/null); do
      resourcename=$(echo $resourceitem | cut -d/ -f2)
      # Export the resource definition, removing cluster-specific fields
      kubectl get $resource $resourcename -o yaml | \
        grep -v "^\s*creationTimestamp:" | \
        grep -v "^\s*resourceVersion:" | \
        grep -v "^\s*selfLink:" | \
        grep -v "^\s*uid:" | \
        grep -v "^\s*status:" | \
        grep -v "^\s*generation:" | \
        grep -v "^\s*managedFields:" > "${CLUSTER_DIR}/${resource}/${resourcename}.yaml"
    done
  fi
done

echo "Creating Kustomize structure..."
# Create basic kustomize structure
for namespace in $(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}'); do
  mkdir -p "${KUSTOMIZE_DIR}/${namespace}"
  echo "apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:" > "${KUSTOMIZE_DIR}/${namespace}/kustomization.yaml"
  
  # Add references to all exported resources for this namespace
  find "${NAMESPACED_DIR}/${namespace}" -type f -name "*.yaml" | sort | while read file; do
    relative_path=$(realpath --relative-to="${KUSTOMIZE_DIR}/${namespace}" "$file")
    echo "- $relative_path" >> "${KUSTOMIZE_DIR}/${namespace}/kustomization.yaml"
  done
done

# Create root kustomization
echo "apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:" > "${KUSTOMIZE_DIR}/kustomization.yaml"

# Add references to all namespace kustomizations
for namespace in $(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}'); do
  echo "- ${namespace}" >> "${KUSTOMIZE_DIR}/kustomization.yaml"
done

echo "Export complete!"
echo "Exported resources to:"
echo "  - Namespaced resources: ${NAMESPACED_DIR}"
echo "  - Cluster-wide resources: ${CLUSTER_DIR}"
echo "  - Kubeconfig: ${KUBECONFIG_DIR}"
echo "  - Helm releases: ${HELM_DIR}"
echo "  - Kustomize structure: ${KUSTOMIZE_DIR}" 