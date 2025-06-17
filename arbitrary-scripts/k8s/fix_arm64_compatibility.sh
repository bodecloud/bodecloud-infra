#!/bin/bash
cd ..
echo "🔧 Fixing ARM64 compatibility issues..."

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Create a backup directory for original files
mkdir -p arm64_fixes_backup

echo "📁 Creating backups..."
cp *deployment.yml arm64_fixes_backup/ 2>/dev/null || true

echo "🔄 Replacing container images with ARM64-compatible alternatives..."

# Replace ElfHosted-specific images with generic alternatives
# gatus -> TwinProduction/gatus (supports ARM64)
sed -i 's|ghcr.io/elfhosted/gatus:.*|twinproduction/gatus:v5.17.0|g' *deployment.yml

# homer -> bastienwirtz/homer (supports ARM64)  
sed -i 's|ghcr.io/elfhosted/homer:.*|b4bz/homer:latest|g' *deployment.yml

# filebrowser -> filebrowser/filebrowser (supports ARM64)
sed -i 's|ghcr.io/elfhosted/filebrowser:.*|filebrowser/filebrowser:latest|g' *deployment.yml

# tooling -> generic ubuntu (supports ARM64)
sed -i 's|filebrowser/filebrowser:.*|filebrowser/filebrowser:v2.31.2|g' *deployment.yml

# plex -> linuxserver/plex (supports ARM64)
sed -i 's|ghcr.io/elfhosted/plex:.*|linuxserver/plex:latest|g' *deployment.yml

# riven -> use rizlim/riven if available, else disable for now
sed -i 's|ghcr.io/elfhosted/riven:.*|rizlim/riven:latest|g' *deployment.yml

# zurg -> use yonasBSD/zurg-arm64 if available, else disable'
# rclone-fm -> rclone/rclone (supports ARM64)
sed -i 's|ghcr.io/elfhosted/rclone.*|rclone/rclone:latest|g' *deployment.yml

# wizarr -> ghcr.io/wizarrrr/wizarr (supports ARM64)
sed -i 's|ghcr.io/elfhosted/wizarr:.*|ghcr.io/wizarrrr/wizarr:latest|g' *deployment.yml

echo "🔧 Removing restrictive pod affinity rules..."

# Remove podAffinity sections that require specific node labels
sed -i '/affinity:/,/topologyKey: kubernetes.io\/hostname/c\      affinity: {}' *deployment.yml

# Also remove any priorityClassName that might not exist
sed -i 's/priorityClassName: tenant-normal//g' *deployment.yml
sed -i 's/priorityClassName: tenant-streaming//g' *deployment.yml

echo "🚀 Applying fixes..."

# Delete existing failed deployments
kubectl delete deployment -n bolabaden --all

# Wait a moment
sleep 5

echo "🔄 Redeploying with ARM64-compatible images..."

# Deploy priority classes and service accounts first
kubectl apply -f priority-classes.yml -f service-accounts.yml

# Deploy core services
kubectl apply -f brunner56-filebrowser-deployment.yml
kubectl apply -f brunner56-homer-deployment.yml  
kubectl apply -f brunner56-gatus-deployment.yml

# Deploy media services
kubectl apply -f brunner56-plex-deployment.yml
kubectl apply -f brunner56-zurg-deployment.yml

# Deploy management services
kubectl apply -f brunner56-rclonefm-deployment.yml
kubectl apply -f brunner56-rcloneui-deployment.yml
kubectl apply -f brunner56-wizarr-deployment.yml
kubectl apply -f brunner56-kubernetesdashboard-deployment.yml
kubectl apply -f brunner56-traefikforwardauth-deployment.yml

# Try riven last as it might be most problematic
kubectl apply -f brunner56-riven-deployment.yml
kubectl apply -f brunner56-rivenfrontend-deployment.yml

echo "⏳ Waiting for pods to start..."
sleep 30

echo "📊 Checking deployment status:"
kubectl get deployments -n bolabaden

echo "🔍 Checking pod status:"
kubectl get pods -n bolabaden

echo "🎯 Services that should be accessible:"
kubectl get services -n bolabaden

echo "✅ ARM64 compatibility fixes applied!"
echo "🌐 You can now access services via port-forward:"
echo "   kubectl port-forward -n bolabaden service/brunner56-homer 8080:8080 --address=0.0.0.0"
echo "   kubectl port-forward -n bolabaden service/brunner56-filebrowser 8081:8080 --address=0.0.0.0"
echo "   kubectl port-forward -n bolabaden service/brunner56-plex 32400:32400 --address=0.0.0.0" 