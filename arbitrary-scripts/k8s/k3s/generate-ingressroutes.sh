#!/bin/bash

# Auto-generate IngressRoute resources for all services in my-media-stack namespace
# This script provides a semi-dynamic approach to expose services

set -e

NAMESPACE="my-media-stack"
DOMAIN="beatapostapita.duckdns.org"
OUTPUT_FILE="auto-generated-ingressroutes.yaml"

echo "🔍 Discovering services in namespace: $NAMESPACE"

# Clear the output file
> "$OUTPUT_FILE"

# Get all services in the namespace and generate IngressRoutes
k3s kubectl get services -n "$NAMESPACE" --no-headers | while read -r name type cluster_ip external_ip ports age; do
    # Skip headless services and services without ports
    if [[ "$cluster_ip" == "None" ]] || [[ "$ports" == "<none>" ]]; then
        echo "   Skipping headless/portless service: $name"
        continue
    fi
    
    # Extract the first port (assuming it's the main port)
    port=$(echo "$ports" | cut -d'/' -f1 | cut -d':' -f1)
    
    # Generate IngressRoute for this service
    echo "   Creating IngressRoute for: $name (port $port)"
    
    cat >> "$OUTPUT_FILE" << EOF
---
# Auto-generated IngressRoute for $name
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: ${name}-auto-ingressroute
  namespace: $NAMESPACE
  labels:
    auto-generated: "true"
spec:
  entryPoints:
    - websecure
  routes:
    - match: "Host(\`${name}.${DOMAIN}\`)"
      kind: Rule
      services:
        - name: $name
          port: $port
  tls:
    certResolver: beatapostapita_duckdns_letsencrypt

EOF
done

echo "✅ Generated IngressRoutes in: $OUTPUT_FILE"
echo ""
echo "🚀 To apply the auto-generated routes, run:"
echo "   k3s kubectl apply -f $OUTPUT_FILE"
echo ""
echo "🌍 Your services will be accessible at:"
k3s kubectl get services -n "$NAMESPACE" --no-headers | while read -r name type cluster_ip external_ip ports age; do
    if [[ "$cluster_ip" != "None" ]] && [[ "$ports" != "<none>" ]]; then
        echo "   - https://${name}.${DOMAIN}"
    fi
done 