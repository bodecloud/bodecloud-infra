#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Setting up Docker Hub registry credentials to avoid pull rate limits...${NC}"

# Prompt for Docker Hub credentials
echo -e "${GREEN}Please enter your Docker Hub credentials:${NC}"
read -p "Docker Hub Username: " DOCKER_USERNAME
read -s -p "Docker Hub Password: " DOCKER_PASSWORD
echo

# Check if credentials were provided
if [ -z "$DOCKER_USERNAME" ] || [ -z "$DOCKER_PASSWORD" ]; then
    echo -e "${RED}Docker Hub username or password not provided. Exiting.${NC}"
    exit 1
fi

# Create registry secret in default namespace
echo -e "${GREEN}Creating Docker Hub registry secret in default namespace...${NC}"
kubectl create secret docker-registry regcred \
    --docker-server=https://index.docker.io/v1/ \
    --docker-username=$DOCKER_USERNAME \
    --docker-password=$DOCKER_PASSWORD \
    --docker-email=$DOCKER_USERNAME@example.com \
    --namespace=default

# Create registry secret in my-media-stack namespace
echo -e "${GREEN}Creating Docker Hub registry secret in my-media-stack namespace...${NC}"
kubectl create secret docker-registry regcred \
    --docker-server=https://index.docker.io/v1/ \
    --docker-username=$DOCKER_USERNAME \
    --docker-password=$DOCKER_PASSWORD \
    --docker-email=$DOCKER_USERNAME@example.com \
    --namespace=my-media-stack

# Patch default service account in default namespace
echo -e "${GREEN}Patching default service account in default namespace...${NC}"
kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "regcred"}]}' -n default

# Patch default service account in my-media-stack namespace
echo -e "${GREEN}Patching default service account in my-media-stack namespace...${NC}"
kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "regcred"}]}' -n my-media-stack

# Create a ConfigMap for containerd registry configuration
echo -e "${GREEN}Creating containerd registry configuration...${NC}"
cat <<EOF >registries.yaml
mirrors:
  "docker.io":
    endpoint:
      - "https://registry-1.docker.io"
configs:
  "docker.io":
    auth:
      username: $DOCKER_USERNAME
      password: $DOCKER_PASSWORD
EOF

# Apply the ConfigMap
sudo mkdir -p /etc/rancher/k3s
sudo mv registries.yaml /etc/rancher/k3s/registries.yaml

echo -e "${GREEN}Docker Hub registry credentials have been set up.${NC}"
echo -e "${YELLOW}Note: You need to restart k3s for the containerd registry configuration to take effect.${NC}"
echo -e "${YELLOW}If you haven't already restarted k3s after running the fix-storage-path.sh script, you should do so now.${NC}"
echo -e "${GREEN}sudo systemctl restart k3s${NC}"

echo -e "${GREEN}Done!${NC}"
