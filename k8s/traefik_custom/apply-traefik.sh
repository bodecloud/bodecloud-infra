#!/bin/bash

# Apply Traefik manifests in the correct order
echo "Applying Traefik manifests..."

# Create namespace if it doesn't exist
kubectl create namespace my-media-stack --dry-run=client -o yaml | kubectl apply -f -

# Apply RBAC first
if [ -f "traefik-rbac.yaml" ]; then
    echo "Applying RBAC..."
    kubectl apply -f traefik-rbac.yaml
fi

# Apply secrets
echo "Applying secrets..."
kubectl apply -f traefik-secret.yaml

# Apply ConfigMap
echo "Applying ConfigMap..."
kubectl apply -f traefik-configmap.yaml

# Apply PVC
echo "Applying PVC..."
kubectl apply -f traefik-pvc.yaml

# Apply Service
echo "Applying Service..."
kubectl apply -f traefik-service.yaml

# Apply Deployment
echo "Applying Deployment..."
kubectl apply -f traefik-deployment.yaml

# Apply IngressRoute if it exists
if [ -f "traefik-ingressroute.yaml" ]; then
    echo "Applying IngressRoute..."
    kubectl apply -f traefik-ingressroute.yaml
fi

echo "Traefik deployment complete!"
echo "Check status with: kubectl get pods -n my-media-stack -l app=traefik" 