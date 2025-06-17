#!/bin/bash

set -e  # Exit on any error

# Initialize step counter
STEP=1

echo "=== Server Setup Script ==="
echo "This script will:"
echo "1. Update system packages"
echo "2. Install GitHub CLI and plocate and pipx"
echo "3. Check external IP"
echo "4. Partition and format secondary disk (/dev/sdb)"
echo "5. Move Docker storage to new disk"
echo "6. Restart Docker service"
echo ""

read -p "Do you want to continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

echo "=== Step $STEP: Updating system packages ==="
sudo apt update && sudo apt upgrade -y
((STEP++))

echo "=== Step $STEP: Installing packages ==="
sudo apt install gh plocate python3-pip python3-venv pipx -y
((STEP++))

echo "=== Step $STEP: Adding External IP to loopback interface ==="
EXTERNAL_IP=$(curl -s ifconfig.me)
echo "External IP: $EXTERNAL_IP"

# Note: The original command seemed malformed, this is the corrected version
sudo ip addr add ${EXTERNAL_IP}/32 dev lo
((STEP++))

echo "=== Step $STEP: Verifying IP configuration ==="
curl -s ifconfig.me
echo ""
((STEP++))

echo "=== Step $STEP: Listing block devices ==="
lsblk
((STEP++))

echo ""
echo "=== Step $STEP: Partitioning secondary disk /dev/sdb ==="
# Check if /dev/sdb exists
if [ ! -b /dev/sdb ]; then
    echo "ERROR: /dev/sdb does not exist. Please check your disk configuration."
    exit 1
fi

# Warning about data loss
echo "WARNING: This will destroy all data on /dev/sdb (there shouldn't be any data if you're following the instructions)!"
read -p "Are you sure you want to continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

# Automated fdisk partitioning
sudo fdisk /dev/sdb << EOF
n
p
1


w
EOF
((STEP++))

echo "=== Step $STEP: Formatting partition with ext4 ==="
sudo mkfs.ext4 /dev/sdb1
((STEP++))

echo "=== Step $STEP: Creating mount point and mounting ==="
sudo mkdir -p /mnt/block200
sudo mount /dev/sdb1 /mnt/block200
((STEP++))

echo "=== Step $STEP: Moving Docker data to new disk ==="
# Stop Docker first to avoid issues
sudo systemctl stop docker

# Move Docker data (handle potential missing files gracefully)
if [ -d /var/lib/docker ]; then
    echo "Moving Docker data..."
    sudo mv /var/lib/docker /mnt/block200/docker 2>/dev/null || {
        echo "Some files couldn't be moved (this is often normal), continuing..."
        # If mv fails, try with rsync and then remove
        sudo rsync -av /var/lib/docker/ /mnt/block200/docker/
        sudo rm -rf /var/lib/docker
    }
else
    echo "Docker directory doesn't exist, creating new one..."
    sudo mkdir -p /mnt/block200/docker
fi
((STEP++))

echo "=== Step $STEP: Creating symlink ==="
sudo ln -s /mnt/block200/docker /var/lib/docker
((STEP++))

echo "=== Step $STEP: Adding to fstab for persistent mounting ==="
# Get UUID of the partition
UUID=$(sudo blkid -s UUID -o value /dev/sdb1)
echo "UUID=$UUID /mnt/block200 ext4 defaults 0 2" | sudo tee -a /etc/fstab
((STEP++))

echo "=== Step $STEP: Restarting Docker service ==="
sudo systemctl daemon-reload
sudo systemctl restart docker
sudo systemctl enable docker
((STEP++))

echo "=== Step $STEP: Verifying Docker status ==="
sudo systemctl status docker --no-pager --lines=100

echo "=== Setup Complete! ==="
echo "Docker storage has been moved to /mnt/block200/docker (symlinked to /var/lib/docker)"
echo "The disk will be automatically mounted on boot via fstab entry"
echo ""
echo "Summary:"
echo "- External IP: $EXTERNAL_IP"
echo "- New Docker storage: /mnt/block200/docker (symlinked to /var/lib/docker)"
echo "- Disk: /dev/sdb1 mounted at /mnt/block200" 