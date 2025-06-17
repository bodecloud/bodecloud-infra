#!/bin/bash

# Script to export Kubernetes resources using kubernetes-manifest-exporter approach
# Based on https://github.com/jonchen727/kubernetes-manifest-exporter

# Set the output directory
OUTPUT_DIR="$(pwd)"
NAMESPACED_DIR="${OUTPUT_DIR}/manifest-exporter/namespaces"
CLUSTER_DIR="${OUTPUT_DIR}/manifest-exporter/cluster"

# Create directories
mkdir -p "${NAMESPACED_DIR}"
mkdir -p "${CLUSTER_DIR}"

# Check if kubectl-neat is installed (for cleaner YAML output)
if ! command -v kubectl-neat &> /dev/null; then
  echo "kubectl-neat not found. Installing kubectl krew first..."
  
  # Install kubectl krew if not already installed
  if ! command -v kubectl-krew &> /dev/null; then
    echo "Installing kubectl krew..."
    (
      set -x; cd "$(mktemp -d)" &&
      OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
      ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
      KREW="krew-${OS}_${ARCH}" &&
      curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
      tar zxvf "${KREW}.tar.gz" &&
      ./"${KREW}" install krew
    )
    
    # Add krew to PATH for this session
    export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
  fi
  
  # Install kubectl-neat
  echo "Installing kubectl-neat..."
  kubectl krew install neat
  
  # Add to PATH for this session
  export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
fi

echo "Exporting namespaced resources..."
# Export namespaced resources
for namespace in $(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}'); do
  echo "Processing namespace: $namespace"
  mkdir -p "${NAMESPACED_DIR}/${namespace}"
  
  # Get all resource types available in this namespace
  for resource in $(kubectl api-resources --namespaced=true --verbs=list -o name | grep -v "events\|events.events.k8s.io"); do
    echo "  Exporting $resource"
    
    # Get all resources of this type in the namespace
    if kubectl get $resource -n $namespace 2>/dev/null | grep -v "No resources found" > /dev/null; then
      mkdir -p "${NAMESPACED_DIR}/${namespace}/${resource}"
      
      # Get all resource names
      for resourceitem in $(kubectl get $resource -n $namespace -o name 2>/dev/null); do
        resourcename=$(echo $resourceitem | cut -d/ -f2)
        echo "    Exporting $resourcename"
        
        # Export the resource definition and clean it with kubectl-neat
        kubectl get $resource $resourcename -n $namespace -o yaml | kubectl neat > "${NAMESPACED_DIR}/${namespace}/${resource}/${resourcename}.yaml"
      done
    fi
  done
done

echo "Exporting cluster-wide resources..."
# Export cluster-wide resources
for resource in $(kubectl api-resources --namespaced=false --verbs=list -o name | grep -v "componentstatuses\|nodes"); do
  echo "  Exporting $resource"
  mkdir -p "${CLUSTER_DIR}/${resource}"
  
  # Get all resources of this type
  if kubectl get $resource 2>/dev/null | grep -v "No resources found" > /dev/null; then
    # Get all resource names
    for resourceitem in $(kubectl get $resource -o name 2>/dev/null); do
      resourcename=$(echo $resourceitem | cut -d/ -f2)
      echo "    Exporting $resourcename"
      
      # Export the resource definition and clean it with kubectl-neat
      kubectl get $resource $resourcename -o yaml | kubectl neat > "${CLUSTER_DIR}/${resource}/${resourcename}.yaml"
    done
  fi
done

echo "Export complete!"
echo "Exported resources to:"
echo "  - Namespaced resources: ${NAMESPACED_DIR}"
echo "  - Cluster-wide resources: ${CLUSTER_DIR}" 