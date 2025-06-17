#!/bin/bash

# Create the directory for kubelet configuration overrides if it doesn't exist
sudo mkdir -p /etc/systemd/system/k3s.service.d/

# Create the override.conf file with the allowed-unsafe-sysctls flag
cat <<EOF | sudo tee /etc/systemd/system/k3s.service.d/override.conf
[Service]
ExecStart=
ExecStart=/usr/local/bin/k3s server --kubelet-arg=allowed-unsafe-sysctls=net.ipv4.ip_forward,net.ipv4.conf.all.forwarding,net.ipv6.conf.all.forwarding,net.ipv6.conf.all.disable_ipv6,net.ipv4.conf.all.src_valid_mark
EOF

# Reload systemd configuration
sudo systemctl daemon-reload

# Restart k3s service
sudo systemctl restart k3s

# Wait for k3s to restart
echo "Waiting for k3s to restart..."
sleep 30

# Clean up any existing failed pods
echo "Cleaning up existing failed pods..."
kubectl delete pods -n vpn-gateway --field-selector=status.phase=Failed

# Apply the fixed warp-gateway deployment
echo "Applying fixed warp-gateway deployment..."
kubectl delete deployment warp-gateway -n vpn-gateway --ignore-not-found
kubectl apply -f ./k8s/my-media-stack/raw_v2/fixed-warp-gateway-deployment.yaml

# Check the status of the pods
echo "Checking pod status..."
sleep 10
kubectl get pods -n vpn-gateway

echo "Done! The warp-gateway should now be able to use the required sysctls."
