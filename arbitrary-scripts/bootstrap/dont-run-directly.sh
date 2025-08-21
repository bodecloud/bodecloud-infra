#!/bin/bash
set -euo pipefail
#set -x

# Set hostname and FQDN
hostnamectl set-hostname "$1"
echo "$1.bolabaden.org" > /etc/hostname
echo "127.0.1.1 $1.bolabaden.org $1" >> /etc/hosts

apt-get update
apt-get autoremove -y
apt-get upgrade -y
curl -sSL https://get.docker.com/ | CHANNEL=stable sh
apt-get install -y curl wget git htop nano vim unzip jq bc yq \
    iptables-persistent python3 python3-pip \
    nodejs npm python3-venv pipx plocate

apt-get autoremove -y

echo iptables-persistent iptables-persistent/autosave_v4 boolean true | sudo debconf-set-selections
echo iptables-persistent iptables-persistent/autosave_v6 boolean true | sudo debconf-set-selections

sudo usermod -aG docker $USER
if [ "$(cat /etc/passwd | grep -q "ubuntu")" ]; then
    sudo usermod -aG docker ubuntu
fi
if [ "$USER" != "root" ]; then
    sudo usermod -aG docker root
fi
sudo chown "$USER":"$USER" /home/"$USER"/.docker -R || true
sudo chmod g+rwx "$HOME/.docker" -R || true

# Enable and start services
systemctl enable docker.service
systemctl enable containerd.service
systemctl start docker.service
systemctl start containerd.service

# Python tools
pipx ensurepath
# Try register-python-argcomplete (v2)
if ! eval "$(register-python-argcomplete pipx 2>&1)"; then
  # If failed, try register-python-argcomplete3 (v3)
  if ! eval "$(register-python-argcomplete3 pipx 2>&1)"; then
    # If both fail, show errors from both commands
    echo "Both register-python-argcomplete and register-python-argcomplete3 failed:"
    echo "Output from register-python-argcomplete:"
    register-python-argcomplete pipx
    echo "Output from register-python-argcomplete3:"
    register-python-argcomplete3 pipx
  fi
fi

pipx install uv

# Node tools
# Install node globally with n
curl -L https://raw.githubusercontent.com/tj/n/master/bin/n -o /usr/local/bin/n
chmod +x /usr/local/bin/n

# Install LTS version of node globally
n lts

# Security updates
cat >/etc/apt/apt.conf.d/20auto-upgrades <<EOF
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
EOF

# Set timezone
timedatectl set-timezone America/Chicago

# SSH config
grep -qE '^\s*PasswordAuthentication\s+' /etc/ssh/sshd_config && \
  sed -i 's/^\s*PasswordAuthentication\s.*/PasswordAuthentication yes/' /etc/ssh/sshd_config || \
  echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config

grep -qE '^\s*PermitRootLogin\s+' /etc/ssh/sshd_config && \
  sed -i 's/^\s*PermitRootLogin\s.*/PermitRootLogin yes/' /etc/ssh/sshd_config || \
  echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config

rm -rvf /usr/share/keyrings/hashicorp-archive-keyring.gpg
apt-get update && apt-get install wget gpg coreutils
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
| tee /etc/apt/sources.list.d/hashicorp.list

apt-get update && apt-get install nomad


export ARCH_CNI=$( [ $(uname -m) = aarch64 ] && echo arm64 || echo amd64)
export CNI_PLUGIN_VERSION=v1.6.2
curl -L -o cni-plugins.tgz "https://github.com/containernetworking/plugins/releases/download/${CNI_PLUGIN_VERSION}/cni-plugins-linux-${ARCH_CNI}-${CNI_PLUGIN_VERSION}".tgz && \
  mkdir -p /opt/cni/bin && \
  tar -C /opt/cni/bin -xzf cni-plugins.tgz

apt-get install -y consul-cni

modprobe br_netfilter
ls /proc/sys/net/bridge -l
echo 1 > /proc/sys/net/bridge/bridge-nf-call-arptables
echo 1 > /proc/sys/net/bridge/bridge-nf-call-ip6tables
echo 1 > /proc/sys/net/bridge/bridge-nf-call-iptables

mkdir -p /etc/sysctl.d
echo "net.bridge.bridge-nf-call-arptables = 1" >> /etc/sysctl.d/bridge.conf
echo "net.bridge.bridge-nf-call-ip6tables = 1" >> /etc/sysctl.d/bridge.conf
echo "net.bridge.bridge-nf-call-iptables = 1" >> /etc/sysctl.d/bridge.conf

sysctl --system

cat /sys/fs/cgroup/cgroup.controllers

rm -rvf /usr/local/bin/nomad

nomad -v

sudo tee /etc/nomad.d/nomad.hcl > /dev/null <<'EOF'
datacenter = "dc1"  # "Nomad clusters can scale horizontally for increased capacity."
data_dir = "/nomad/data/"
bind_addr = "0.0.0.0"  # "Bind to all interfaces for multi-node communication." 
log_level = "DEBUG"  # "Set the log level to DEBUG for more detailed logging." 
log_json = true  # "Enable JSON logging for better machine readability." 
log_file = "/nomad/log/nomad.log"  # "Set the log file path." 
log_rotate_bytes = 10485760  # "Rotate the log file when it reaches 10MB." 
log_rotate_max_files = 5  # "Keep up to 5 rotated log files." 

server {
  enabled = true  # "Servers talk to each other and use a leader/follower load balancing method for HA." 
  bootstrap_expect = 5  # "Create a multi-server (HA) setup. Prepare, for example, three nodes... but for five nodes, set bootstrap_expect to 5 for quorum." Adapted from  for odd-number quorum to avoid SPOF.
  server_join {
    retry_join = [
      "170.9.225.137:4647",
      "149.130.219.117:4647",
      "149.130.222.229:4647",
      "150.136.84.225:4647",
      "172.245.88.16:4647"
    ]
  }
}

client {
  enabled = true  # "A single client process can handle running many allocations on a single node." 
  node_class = "balanced"  # Custom class for load spreading; "Nomad uses a bin packing algorithm, which means it tries to utilize all of a node's resources before placing tasks on a different node." 
}

advertise {
  http = "{{ GetPrivateIP }}:4646"
  rpc  = "{{ GetPrivateIP }}:4647"
  serf = "{{ GetPrivateIP }}:4648"
}

consul {
  address = "consul:8500"  # "Nomad integrates with Consul to provide service discovery and monitoring."  Assuming Consul agent on localhost.
  auto_advertise = true  # "Nomad can register services with Consul." 
  server_service_name = "nomad"  # "Consul allows services to easily register themselves in a central catalog." 
  client_service_name = "nomad-client"  # "Nomad integrates with Consul to provide service discovery." 
}

telemetry {
  collection_interval = "1s"  # For monitoring CPU/RAM balance; "Optimize the raft_multiplier." 
  publish_allocation_metrics = true  # "Nomad does not seem to balance the allocations across the clients... but with telemetry, it can monitor usage." 
}
EOF

sudo systemctl restart nomad && sudo systemctl enable nomad

nomad server members

nomad node status

nomad node status -self

sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstabsystemctl restart php7.4-fpm




# Final message
echo "VPS Server '$1' is ready!" | tee /var/log/cloud-init-complete.log
