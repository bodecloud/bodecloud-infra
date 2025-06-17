#!/bin/bash

cd "$(dirname "$0")"

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
  echo "Error: kubectl is not installed. Please install kubectl first."
  exit 1
fi

# Check if k3s is installed
if ! command -v k3s &> /dev/null; then
  echo "Error: k3s is not installed. Please install k3s first."
  echo "Visit https://k3s.io/ for installation instructions."
  exit 1
fi

# Check if k3s service is running and start it if not
if ! systemctl is-active --quiet k3s; then
  echo "K3s service is not running. Attempting to start it..."
  if [ "$EUID" -ne 0 ]; then
    echo "This script needs to be run with sudo to start the k3s service."
    echo "Please run: sudo $0"
    exit 1
  fi
  
  systemctl start k3s
  
  # Wait for K3s to be ready
  echo "Waiting for K3s to start (up to 60 seconds)..."
  TIMEOUT=60
  while [ $TIMEOUT -gt 0 ]; do
    if kubectl cluster-info &>/dev/null; then
      echo "K3s started successfully."
      break
    fi
    TIMEOUT=$((TIMEOUT-1))
    sleep 1
  done
  
  if [ $TIMEOUT -eq 0 ]; then
    echo "Timed out waiting for K3s to start. Please check your K3s installation."
    exit 1
  fi
else
  echo "K3s service is already running."
fi

# Check if kubectl can connect to the cluster
if ! kubectl cluster-info &> /dev/null; then
  echo "Error: kubectl cannot connect to the Kubernetes cluster."
  echo "Please ensure your cluster is running and kubectl is properly configured."
  echo "For K3s, you might need to set:"
  echo "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml"
  exit 1
fi

echo "Connected to Kubernetes cluster. Deploying Traefik..."

# Apply RBAC resources
kubectl apply -f 00-role.yml
kubectl apply -f 00-account.yml
kubectl apply -f 01-role-binding.yml

# Apply Traefik configuration
kubectl apply -f traefik-env-configmap.yml

# Apply Error Pages
kubectl apply -f 05-traefik-error-pages.yml
kubectl apply -f 06-traefik-middlewares.yml

# Wait for error pages to be ready
echo "Waiting for error pages to be ready..."
kubectl wait --for=condition=available deployment/traefik-error-pages --timeout=60s

# Apply Traefik resources
kubectl apply -f 02-traefik.yml
kubectl apply -f 02-traefik-services.yml

# Apply whoami resources
kubectl apply -f 03-whoami.yml
kubectl apply -f 03-whoami-services.yml
kubectl apply -f 04-whoami-ingress.yml

echo "Traefik and whoami application deployed successfully!"
echo "Access the Traefik dashboard at: http://localhost:8080"
echo "Access the whoami application at: http://whoami.${DOMAIN:-bolabaden.org}"

# Optional: Check for DNS resolution
echo "Testing DNS resolution to whoami.bolabaden.org..."
if dig +short whoami.bolabaden.org > /dev/null; then
  echo "DNS resolution successful"
else
  echo "Warning: Could not resolve whoami.bolabaden.org"
  echo "Make sure your DNS settings are correctly configured."
fi 