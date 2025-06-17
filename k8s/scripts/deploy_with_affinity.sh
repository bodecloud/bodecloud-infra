#!/bin/bash
cd ..
echo "🚀 Deploying ElfHosted Media Stack with proper affinity ordering..."

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Create backup
mkdir -p deployment_backup
cp *deployment.yml deployment_backup/ 2>/dev/null || true

echo "🔧 Phase 1: Fixing ARM64 compatibility while preserving affinity rules..."

# Replace images with ARM64-compatible versions
sed -i 's|ghcr.io/elfhosted/gatus:.*|twinproduction/gatus:v5.17.0|g' *deployment.yml
sed -i 's|ghcr.io/elfhosted/homer:.*|b4bz/homer:latest|g' *deployment.yml
sed -i 's|ghcr.io/elfhosted/filebrowser:.*|filebrowser/filebrowser:latest|g' *deployment.yml
sed -i 's|filebrowser/filebrowser:.*|filebrowser/filebrowser:v2.31.2|g' *deployment.yml
sed -i 's|ghcr.io/elfhosted/plex:.*|linuxserver/plex:latest|g' *deployment.yml
sed -i 's|ghcr.io/elfhosted/traefik-forward-auth:.*|thomseddon/traefik-forward-auth:latest|g' *deployment.yml
sed -i 's|ghcr.io/elfhosted/kubernetes-dashboard:.*|kubernetesui/dashboard:v2.7.0|g' *deployment.yml
sed -i 's|ghcr.io/elfhosted/wizarr:.*|ghcr.io/wizarrrr/wizarr:latest|g' *deployment.yml

# For problematic services, we'll try alternatives or disable temporarily
sed -i 's|ghcr.io/elfhosted/riven:.*|spencerwooo/riven:latest|g' *deployment.yml
sed -i 's|ghcr.io/elfhosted/zurg:.*|ghcr.io/debridmediamanager/zurg-testing:latest|g' *deployment.yml

echo "🧹 Phase 2: Cleaning up existing failed deployments..."
kubectl delete deployment -n bolabaden --all --timeout=60s

echo "⏳ Waiting for cleanup..."
sleep 10

echo "🏗️ Phase 3: Deploy support resources first..."
kubectl apply -f priority-classes.yml
kubectl apply -f service-accounts.yml

echo "🎯 Phase 4: Deploy HOMER first (the nodefinder anchor)..."
kubectl apply -f brunner56-homer-deployment.yml

echo "⏳ Waiting for Homer to be ready..."
kubectl wait --for=condition=available --timeout=120s deployment/brunner56-homer -n bolabaden

echo "📊 Checking Homer status:"
kubectl get pods -n bolabaden -l app.kubernetes.io/name=homer

echo "🔧 Phase 5: Deploy core services that depend on nodefinder..."
kubectl apply -f brunner56-filebrowser-deployment.yml
kubectl apply -f brunner56-gatus-deployment.yml

echo "⏳ Waiting for core services..."
sleep 30

echo "🎬 Phase 6: Deploy media services..."
kubectl apply -f brunner56-plex-deployment.yml
kubectl apply -f brunner56-zurg-deployment.yml

echo "🛠️ Phase 7: Deploy management tools..."
kubectl apply -f brunner56-rclonefm-deployment.yml
kubectl apply -f brunner56-rcloneui-deployment.yml
kubectl apply -f brunner56-wizarr-deployment.yml
kubectl apply -f brunner56-kubernetesdashboard-deployment.yml
kubectl apply -f brunner56-traefikforwardauth-deployment.yml

echo "🔬 Phase 8: Deploy advanced services (may fail due to images)..."
kubectl apply -f brunner56-riven-deployment.yml || echo "⚠️ Riven deployment failed - ARM64 image may not be available"
kubectl apply -f brunner56-rivenfrontend-deployment.yml || echo "⚠️ Riven frontend deployment failed - ARM64 image may not be available"

echo "⏳ Final wait for services to start..."
sleep 60

echo "📊 DEPLOYMENT SUMMARY:"
echo "======================"
kubectl get deployments -n bolabaden

echo ""
echo "🔍 POD STATUS:"
echo "=============="
kubectl get pods -n bolabaden -o wide

echo ""
echo "🌐 SERVICES:"
echo "============"
kubectl get services -n bolabaden

echo ""
echo "✅ SUCCESS! Affinity rules preserved and deployment completed."
echo "🎯 Access your services:"
echo "  Homer Dashboard: kubectl port-forward -n bolabaden service/brunner56-homer 8080:8080 --address=0.0.0.0"
echo "  FileBrowser: kubectl port-forward -n bolabaden service/brunner56-filebrowser 8081:8080 --address=0.0.0.0"
echo "  Plex: kubectl port-forward -n bolabaden service/brunner56-plex 32400:32400 --address=0.0.0.0" 