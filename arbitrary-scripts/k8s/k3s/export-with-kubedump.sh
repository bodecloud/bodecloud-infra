#!/bin/bash

# Script to export Kubernetes resources using kube-dump
# Based on https://github.com/WoozyMasta/kube-dump

# Check if kube-dump is installed
if ! command -v kube-dump &> /dev/null; then
  echo "kube-dump not found. Installing..."
  
  # Install kube-dump
  curl -sL https://raw.githubusercontent.com/WoozyMasta/kube-dump/master/kube-dump -o /tmp/kube-dump
  chmod +x /tmp/kube-dump
  KUBE_DUMP_CMD="/tmp/kube-dump"
else
  KUBE_DUMP_CMD="kube-dump"
fi

# Set the output directory
OUTPUT_DIR="$(pwd)"

# Run kube-dump to export all resources
echo "Exporting all Kubernetes resources using kube-dump..."
$KUBE_DUMP_CMD dump -d "$OUTPUT_DIR/kube-dump" --clean --all-resources

echo "Export complete!"
echo "Exported resources to: $OUTPUT_DIR/kube-dump" 