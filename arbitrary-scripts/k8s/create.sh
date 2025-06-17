#!/bin/bash

# Install k3s
cd ..
curl -sfL https://get.k3s.io | sh -
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
echo 'export KUBECONFIG=/etc/rancher/k3s/k3s.yaml' >> ~/.bashrc

# Create namespace
kubectl create namespace bolabaden

kubectl create secret generic plex-config -n bolabaden \
 --from-literal=PLEX_TOKEN=18rS6gM-Wec22uq5Gw-U \
 --from-literal=PLEX_CLAIM=claim-RmAd7Q1DfvybLKDY8zmQ \
 --from-literal=PLEX_UID=1001 \
 --from-literal=PLEX_GID=989

kubectl create secret generic cloudflare-api-token -n bolabaden \
 --from-literal=CF_API_TOKEN=fp6tjZWb66p1rCyZBM9k0hxrE0n0Y7XctSTwP52I \
 --from-literal=CF_EMAIL=boden.crouch@gmail.com \
 --from-literal=CF_ZONE_ID=164c8d72507295b51851d9b05f0e37a1

kubectl create secret generic gatus-smtp-config -n bolabaden \
 --from-literal=SMTP_USERNAME=halomastar@gmail.com \
 --from-literal=SMTP_PASSWORD=re_Ns9xDLUJ_JxnRW8m9hpUgyjecDgNVmSzS \
 --from-literal=SMTP_HOST=smtp.resend.com

# Phase 1: Deploy ConfigMaps
echo "🚀 Phase 1: Deploying ConfigMaps..."
find . -name "*configmap.yml" -exec kubectl apply -f {} \;

# Phase 2: Deploy Persistent Volume Claims
echo "🗄️ Phase 2: Deploying Persistent Volume Claims..."
find . -name "*persistentvolumeclaim.yml" -exec kubectl apply -f {} \;

echo "✅ Checking PVC status..."
kubectl get pvc -n bolabaden

# Phase 3: Deploy Core Services
echo "🔧 Phase 3: Deploying Core Services..."
kubectl apply -f brunner56-filebrowser-deployment.yml -f brunner56-filebrowser-service.yml
kubectl apply -f brunner56-homer-deployment.yml -f brunner56-homer-service.yml
kubectl apply -f brunner56-gatus-deployment.yml -f brunner56-gatus-service.yml

# Phase 4: Deploy Media Services
echo "🎬 Phase 4: Deploying Media Services..."
kubectl apply -f brunner56-zurg-deployment.yml -f brunner56-zurg-service.yml
kubectl apply -f brunner56-plex-deployment.yml -f brunner56-plex-service.yml
kubectl apply -f brunner56-riven-deployment.yml -f brunner56-riven-service.yml
kubectl apply -f brunner56-rivenfrontend-deployment.yml -f brunner56-rivenfrontend-service.yml

# Phase 5: Deploy Remaining Services
echo "🛠️ Phase 5: Deploying Remaining Services..."
kubectl apply -f brunner56-rclonefm-deployment.yml -f brunner56-rclonefm-service.yml
kubectl apply -f brunner56-rcloneui-deployment.yml -f brunner56-rcloneui-service.yml
kubectl apply -f brunner56-wizarr-deployment.yml -f brunner56-wizarr-service.yml
kubectl apply -f brunner56-kubernetesdashboard-deployment.yml -f brunner56-kubernetesdashboard-service.yml
kubectl apply -f brunner56-traefikforwardauth-deployment.yml -f brunner56-traefikforwardauth-service.yml
kubectl apply -f brunner56-starter-deployment.yml

# Phase 6: Deploy Standalone Services
echo "🔄 Phase 6: Deploying Standalone Services..."
find . -name "*-service.yml" ! -name "brunner56-*" -exec kubectl apply -f {} \;
find . -name "*-pod.yml" -exec kubectl apply -f {} \;

# Deploy priority classes and service accounts
kubectl apply -f priority-classes.yml
kubectl apply -f service-accounts.yml

# Restart deployments to pick up new configurations
echo "🔄 Restarting deployments..."
kubectl rollout restart deployment -n bolabaden
sleep 10

# Final status check
echo "🎯 FINAL DEPLOYMENT STATUS"
echo "========================="
echo "✅ k3s cluster: $(kubectl get nodes --no-headers | wc -l) node(s)"
echo "✅ Namespace: bolabaden"
echo "✅ ConfigMaps: $(kubectl get configmaps -n bolabaden --no-headers | wc -l)"
echo "✅ Services: $(kubectl get services -n bolabaden --no-headers | wc -l)"
echo "✅ Deployments: $(kubectl get deployments -n bolabaden --no-headers | wc -l)"
echo "✅ PVCs: $(kubectl get pvc -n bolabaden --no-headers | wc -l)"
echo "✅ Running Pods:"
kubectl get pods -n bolabaden | grep Running

# Start port forwarding to Homer dashboard
kubectl port-forward -n bolabaden service/brunner56-homer 8080:8080 --address=0.0.0.0

# Print final message
echo -e "\n${GREEN}🎉 Deployment completed successfully!${NC}"
echo -e "You can access the Homer dashboard at http://localhost:8080"
echo -e "You can access the Riven frontend at http://localhost:3000"
echo -e "You can access the Riven backend at http://localhost:8080"
echo -e "You can access the Riven frontend at http://localhost:3000"