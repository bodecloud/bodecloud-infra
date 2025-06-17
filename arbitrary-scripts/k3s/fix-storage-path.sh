#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Fixing k3s storage paths to use /mnt/blockvolume instead of root filesystem...${NC}"

# Create necessary directories on blockvolume
echo -e "${GREEN}Creating necessary directories on /mnt/blockvolume...${NC}"
sudo mkdir -p /mnt/blockvolume/k3s-data
sudo mkdir -p /mnt/blockvolume/k3s-containerd

# Check if directories were created successfully
if [ ! -d "/mnt/blockvolume/k3s-data" ] || [ ! -d "/mnt/blockvolume/k3s-containerd" ]; then
  echo -e "${RED}Failed to create directories on /mnt/blockvolume. Please check if the volume is mounted correctly.${NC}"
  exit 1
fi

# Set proper permissions
sudo chmod 755 /mnt/blockvolume/k3s-data
sudo chmod 755 /mnt/blockvolume/k3s-containerd

# Backup existing override.conf
echo -e "${GREEN}Backing up existing k3s override configuration...${NC}"
sudo cp /etc/systemd/system/k3s.service.d/override.conf /etc/systemd/system/k3s.service.d/override.conf.bak

# Update the override.conf file to use the new paths
echo -e "${GREEN}Updating k3s configuration to use /mnt/blockvolume...${NC}"
sudo tee /etc/systemd/system/k3s.service.d/override.conf >/dev/null <<EOF
[Service]
ExecStart=
ExecStart=/usr/local/bin/k3s server \\
  --bind-address=0.0.0.0 \\
  --advertise-address=10.0.0.81 \\
  --tls-san=10.0.0.81 \\
  --tls-san=0.0.0.0 \\
  --kubelet-arg=allowed-unsafe-sysctls=net.* \\
  --data-dir=/mnt/blockvolume/k3s-data \\
  --default-local-storage-path=/mnt/blockvolume/k3s-data/storage \\
  --containerd-config=/mnt/blockvolume/k3s-containerd/config.toml \\
  --write-kubeconfig-mode 644
EOF

# Create containerd config directory and file
echo -e "${GREEN}Creating containerd configuration...${NC}"
sudo mkdir -p /mnt/blockvolume/k3s-containerd
sudo tee /mnt/blockvolume/k3s-containerd/config.toml >/dev/null <<EOF
version = 2
root = "/mnt/blockvolume/k3s-containerd/containerd"
state = "/mnt/blockvolume/k3s-containerd/run"

[grpc]
  address = "/run/k3s/containerd/containerd.sock"

[plugins."io.containerd.grpc.v1.cri"]
  stream_server_address = "127.0.0.1"
  stream_server_port = "10010"
  enable_selinux = false
  enable_unprivileged_ports = false
  enable_unprivileged_icmp = false
  sandbox_image = "rancher/mirrored-pause:3.6"

[plugins."io.containerd.grpc.v1.cri".containerd]
  snapshotter = "overlayfs"
  disable_snapshot_annotations = true

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  runtime_type = "io.containerd.runc.v2"

[plugins."io.containerd.grpc.v1.cri".registry]
  config_path = "/etc/rancher/k3s/registries.yaml"
EOF

# Reload systemd to apply changes
echo -e "${GREEN}Reloading systemd daemon...${NC}"
sudo systemctl daemon-reload

echo -e "${YELLOW}Configuration updated. You need to restart k3s for changes to take effect.${NC}"
echo -e "${YELLOW}This will temporarily disrupt your Kubernetes cluster.${NC}"
echo -e "${YELLOW}Run the following command to restart k3s:${NC}"
echo -e "${GREEN}sudo systemctl restart k3s${NC}"

# Ask if user wants to restart k3s now
read -p "Do you want to restart k3s now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo -e "${GREEN}Restarting k3s...${NC}"
  sudo systemctl restart k3s
  echo -e "${GREEN}Waiting for k3s to restart (this may take a few minutes)...${NC}"
  sleep 30
  echo -e "${GREEN}Checking k3s status:${NC}"
  sudo systemctl status k3s --no-pager
else
  echo -e "${YELLOW}Please restart k3s manually when ready using:${NC}"
  echo -e "${GREEN}sudo systemctl restart k3s${NC}"
fi

echo -e "${GREEN}Done!${NC}"
