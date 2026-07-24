#!/bin/bash
set -euo pipefail

# ============================================================================
# VPS Bootstrap Script - Idempotent & Configurable
# ============================================================================
# This script can be run multiple times safely to ensure a known working state.
# It doesn't skip steps - instead, operations are naturally idempotent.
# All configuration can be provided via environment variables or config file.
# ============================================================================

# Enable debug mode if requested
[ "${DEBUG:-false}" = "true" ] && set -x

# ============================================================================
# CONFIGURATION - Override with environment variables or config file
# ============================================================================

# Load explicit + bootstrap config (explicit wins over bootstrap defaults)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REQUIRED_EXPLICIT_FILE="${REQUIRED_EXPLICIT_FILE:-/etc/required-explicit.env}"
if [ -f "$REQUIRED_EXPLICIT_FILE" ]; then
  # shellcheck source=/dev/null
  source "$REQUIRED_EXPLICIT_FILE"
elif [ -f "${SCRIPT_DIR}/required-explicit.env" ]; then
  # shellcheck source=/dev/null
  source "${SCRIPT_DIR}/required-explicit.env"
fi
CONFIG_FILE="${BOOTSTRAP_CONFIG_FILE:-/etc/bootstrap-config.env}"
[ -f "$CONFIG_FILE" ] && source "$CONFIG_FILE"

# Core Configuration
DOMAIN="${DOMAIN:-bolabaden.org}"
HOSTNAME_ARG="${1:-${BOOTSTRAP_HOSTNAME:-${TS_HOSTNAME:-$(hostname -s 2>/dev/null || hostname | cut -d'.' -f1)}}}"
TS_HOSTNAME="${TS_HOSTNAME:-$HOSTNAME_ARG}"
HOSTNAME_SHORT=$(echo "$HOSTNAME_ARG" | cut -d'.' -f1)
FQDN="${HOSTNAME_SHORT}.${DOMAIN}"

# User Configuration
PRIMARY_USER="${PRIMARY_USER:-ubuntu}"
ADMIN_USERS="${ADMIN_USERS:-root ubuntu}"
PASSWORD_HASH="${PASSWORD_HASH:-\$6\$pWurw/L0tau67C7g\$kiM8cWIAg97/je2BQLKAm/FRuTz1Xu.g0UC59HuqK0d2jkLqw1FcDcB8YH.Iv0PEh3DhyMPosfmEWCi/AnmrX.}"
GITHUB_USERS="${GITHUB_USERS:-th3w1zard1}" # Comma-separated list

# DNS Configuration
DNS_SERVERS="${DNS_SERVERS:-1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4}"
IFS=',' read -ra DNS_ARRAY <<<"$DNS_SERVERS"

# Tailscale Configuration
ENABLE_TAILSCALE="${ENABLE_TAILSCALE:-true}"
TAILSCALE_AUTH_KEY="${TAILSCALE_AUTH_KEY:-}"
TAILSCALE_LOGIN_SERVER="${TAILSCALE_LOGIN_SERVER:-https://headscale.${DOMAIN}}"
TAILSCALE_ADVERTISE_EXIT="${TAILSCALE_ADVERTISE_EXIT:-true}"

# Docker Configuration
ENABLE_DOCKER="${ENABLE_DOCKER:-true}"
DOCKER_VERSION="${DOCKER_VERSION:-27.0}"

# Nomad/Consul Configuration
ENABLE_NOMAD="${ENABLE_NOMAD:-true}"
ENABLE_CONSUL="${ENABLE_CONSUL:-true}"
NOMAD_DATACENTER="${NOMAD_DATACENTER:-dc1}"
NOMAD_BOOTSTRAP_EXPECT="${NOMAD_BOOTSTRAP_EXPECT:-1}"
NOMAD_NODE_CLASS="${NOMAD_NODE_CLASS:-balanced}"
NOMAD_SERVERS="${NOMAD_SERVERS:-}" # Comma-separated IPs, auto-detected if empty

# Swap Configuration
SWAP_SIZE="${SWAP_SIZE:-4G}"
SWAP_FILE="${SWAP_FILE:-/swapfile}"

# SSH Configuration
SSH_PERMIT_ROOT="${SSH_PERMIT_ROOT:-yes}"
SSH_PASSWORD_AUTH="${SSH_PASSWORD_AUTH:-yes}"

# Timezone Configuration
TZ="${TZ:-}" # Auto-detect via GeoIP if empty

# ============================================================================
# HELPER FUNCTIONS
# ============================================================================

log() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*"
}

log_info() {
  log "[*] $*"
}

log_success() {
  log "  [OK] $*"
}

log_warn() {
  log "  [WARN] $*"
}

log_error() {
  log "  [ERROR] $*"
}

# ============================================================================
# MAIN SCRIPT
# ============================================================================

log_info "========================================"
log_info "VPS Bootstrap Script Starting"
log_info "========================================"
log_info "Hostname: $HOSTNAME_SHORT"
log_info "Domain: $DOMAIN"
log_info "FQDN: $FQDN"
log_info "========================================"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  log_error "This script must be run as root"
  exit 1
fi

# ============================================================================
# HOSTNAME CONFIGURATION
# ============================================================================

log_info "Configuring hostname..."
hostnamectl set-hostname "$HOSTNAME_SHORT"
echo "$HOSTNAME_SHORT" >/etc/hostname

# Update /etc/hosts - remove existing entries and add new one
sed -i "/127.0.1.1.*${HOSTNAME_SHORT}/d" /etc/hosts
if ! grep -q "127.0.1.1 ${FQDN} ${HOSTNAME_SHORT}" /etc/hosts; then
  echo "127.0.1.1 ${FQDN} ${HOSTNAME_SHORT}" >>/etc/hosts
fi
log_success "Hostname configured"

# ============================================================================
# SYSTEM UPDATE
# ============================================================================

log_info "Updating system packages..."
apt-get update -qq
apt-get autoremove -y -qq
DEBIAN_FRONTEND=noninteractive apt-get upgrade -y -qq
log_success "System packages updated"

# ============================================================================
# ESSENTIAL PACKAGES
# ============================================================================

log_info "Installing essential packages..."
DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
  curl wget git htop nano vim unzip jq bc yq \
  iptables-persistent python3 python3-pip \
  nodejs npm python3-venv pipx plocate whois sshpass ansible \
  dnsutils bind9-host ethtool ca-certificates gnupg lsb-release 2>&1 | grep -v "^Preconfiguring" || true

apt-get autoremove -y -qq
log_success "Essential packages installed"

# ============================================================================
# DOCKER INSTALLATION
# ============================================================================

if [ "$ENABLE_DOCKER" = "true" ]; then
  log_info "Installing Docker..."

  # Remove conflicting packages
  DEBIAN_FRONTEND=noninteractive apt-get remove -y -qq docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc 2>/dev/null || true

  # Clean up existing Docker sources
  rm -f /etc/apt/sources.list.d/docker.list /etc/apt/sources.list.d/docker.sources

  # Add Docker's official GPG key
  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
  chmod a+r /etc/apt/keyrings/docker.asc

  # Add repository
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" |
    tee /etc/apt/sources.list.d/docker.list >/dev/null

  apt-get update -qq
  DEBIAN_FRONTEND=noninteractive apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

  # Clean up duplicates
  [ -f /etc/apt/sources.list.d/docker.sources ] && rm -f /etc/apt/sources.list.d/docker.list

  # Enable and start Docker
  systemctl enable docker.service containerd.service
  systemctl start docker.service containerd.service

  log_success "Docker installed: $(docker --version)"
fi

# ============================================================================
# SSH CONFIGURATION
# ============================================================================

log_info "Configuring SSH..."

# Backup SSH config if it's the first time or weekly
BACKUP_COUNT=$(find /etc/ssh -name "sshd_config.backup.*" 2>/dev/null | wc -l)
OLD_BACKUP_COUNT=$(find /etc/ssh -name "sshd_config.backup.*" -mtime +7 2>/dev/null | wc -l)
if [ "$BACKUP_COUNT" -eq 0 ] || [ "$OLD_BACKUP_COUNT" -gt 0 ]; then
  cp /etc/ssh/sshd_config /etc/ssh/sshd_config.backup.$(date +%Y%m%d_%H%M%S)
fi

# Configure SSH settings
if grep -qE '^\s*PasswordAuthentication\s+' /etc/ssh/sshd_config; then
  sed -i "s/^\s*PasswordAuthentication\s.*/PasswordAuthentication $SSH_PASSWORD_AUTH/" /etc/ssh/sshd_config
else
  echo "PasswordAuthentication $SSH_PASSWORD_AUTH" >>/etc/ssh/sshd_config
fi

if grep -qE '^\s*PubkeyAuthentication\s+' /etc/ssh/sshd_config; then
  sed -i "s/^\s*PubkeyAuthentication\s.*/PubkeyAuthentication yes/" /etc/ssh/sshd_config
else
  echo "PubkeyAuthentication yes" >>/etc/ssh/sshd_config
fi

if grep -qE '^\s*PermitRootLogin\s+' /etc/ssh/sshd_config; then
  sed -i "s/^\s*PermitRootLogin\s.*/PermitRootLogin $SSH_PERMIT_ROOT/" /etc/ssh/sshd_config
else
  echo "PermitRootLogin $SSH_PERMIT_ROOT" >>/etc/ssh/sshd_config
fi

# Test and reload SSH
if sshd -t 2>/dev/null; then
  systemctl reload sshd || systemctl restart sshd
  log_success "SSH configured"
else
  log_error "SSH configuration test failed"
fi

# ============================================================================
# USER CONFIGURATION
# ============================================================================

log_info "Configuring users..."

for USER in $ADMIN_USERS; do
  log_info "  Setting up user: $USER"

  # Create user if doesn't exist
  if ! id "$USER" &>/dev/null; then
    useradd -m -s /bin/bash -G sudo "$USER" 2>/dev/null || true
  else
    usermod -aG sudo "$USER" 2>/dev/null || true
  fi

  # Add to docker group if Docker is enabled
  if [ "$ENABLE_DOCKER" = "true" ]; then
    usermod -aG docker "$USER" 2>/dev/null || true
  fi

  HOME_DIR=$(eval echo "~$USER")

  # Setup Docker directory
  if [ "$ENABLE_DOCKER" = "true" ]; then
    mkdir -p "$HOME_DIR/.docker"
    chown "$USER":"$USER" "$HOME_DIR/.docker" -R 2>/dev/null || true
    chmod g+rwx "$HOME_DIR/.docker" -R 2>/dev/null || true
  fi

  # Setup SSH keys
  mkdir -p "$HOME_DIR/.ssh"
  chmod 700 "$HOME_DIR/.ssh"
  chown "$USER":"$USER" "$HOME_DIR/.ssh"
  touch "$HOME_DIR/.ssh/authorized_keys"
  chmod 600 "$HOME_DIR/.ssh/authorized_keys"

  # Import SSH keys from GitHub
  IFS=',' read -ra GH_USERS <<<"$GITHUB_USERS"
  for gh_user in "${GH_USERS[@]}"; do
    gh_user=$(echo "$gh_user" | xargs) # Trim whitespace
    if [ -n "$gh_user" ]; then
      GITHUB_KEYS=$(curl -fsSL --max-time 10 "https://github.com/${gh_user}.keys" 2>/dev/null || true)
      if [ -n "$GITHUB_KEYS" ]; then
        while IFS= read -r key; do
          [ -n "$key" ] && ! grep -Fxq "$key" "$HOME_DIR/.ssh/authorized_keys" 2>/dev/null && echo "$key" >>"$HOME_DIR/.ssh/authorized_keys"
        done <<<"$GITHUB_KEYS"
      fi
    fi
  done

  chown "$USER":"$USER" "$HOME_DIR/.ssh/authorized_keys" 2>/dev/null || true

  # Set password
  if [ -n "$PASSWORD_HASH" ]; then
    echo "$USER:${PASSWORD_HASH}" | chpasswd -e 2>/dev/null || true
  fi

  # Setup pipx and uv
  sudo -u "$USER" -H bash -c 'pipx ensurepath 2>/dev/null || true' 2>/dev/null || true
  sudo -u "$USER" -H bash -c 'pipx list | grep -q uv || pipx install uv 2>/dev/null' 2>/dev/null || true
done

log_success "Users configured"

# ============================================================================
# NETWORK OPTIMIZATION
# ============================================================================

log_info "Optimizing network settings..."

# Detect primary interface
IFACE="$(ip route show default 0.0.0.0/0 | awk '{print $5}' | head -n1)"
if [ -n "$IFACE" ]; then
  log_info "  Primary interface: $IFACE"

  # Try to optimize NIC settings
  ethtool -K "$IFACE" gro on 2>/dev/null || true
  ethtool -K "$IFACE" lro off 2>/dev/null || true
  ethtool -K "$IFACE" rx-udp-gro-forwarding on 2>/dev/null || true
  ethtool -K "$IFACE" tx-udp-segmentation on 2>/dev/null || true
fi

# Enable IP forwarding
cat >/etc/sysctl.d/99-forwarding.conf <<EOF
net.ipv4.ip_forward = 1
net.ipv6.conf.all.forwarding = 1
EOF
sysctl -p /etc/sysctl.d/99-forwarding.conf || true

log_success "Network settings optimized"

# ============================================================================
# DNS CONFIGURATION
# ============================================================================

log_info "Configuring DNS..."

mkdir -p /etc/systemd/resolved.conf.d
cat >/etc/systemd/resolved.conf.d/custom.conf <<EOF
[Resolve]
DNS=${DNS_ARRAY[0]} ${DNS_ARRAY[1]:-}
FallbackDNS=${DNS_ARRAY[2]:-} ${DNS_ARRAY[3]:-}
DNSOverTLS=opportunistic
DNSSEC=allow-downgrade
EOF

if systemctl is-active --quiet systemd-resolved; then
  systemctl restart systemd-resolved
else
  # Fallback to /etc/resolv.conf
  : >/etc/resolv.conf
  for dns in "${DNS_ARRAY[@]}"; do
    echo "nameserver $dns" >>/etc/resolv.conf
  done
fi

log_success "DNS configured"

# ============================================================================
# TAILSCALE
# ============================================================================

if [ "$ENABLE_TAILSCALE" = "true" ]; then
  log_info "Installing and configuring Tailscale..."

  if [ -z "$TAILSCALE_AUTH_KEY" ]; then
    log_warn "TAILSCALE_AUTH_KEY not set — set it in /etc/required-explicit.env or export before running"
  fi

  # Install if not present
  if ! command -v tailscale >/dev/null 2>&1; then
    curl -fsSL https://tailscale.com/install.sh | sh
  fi

  # Configure
  tailscale down 2>/dev/null || true

  TS_CMD="tailscale up --login-server=${TAILSCALE_LOGIN_SERVER} --hostname=${HOSTNAME_SHORT} --operator=${PRIMARY_USER} --accept-dns=true"
  [ -n "$TAILSCALE_AUTH_KEY" ] && TS_CMD+=" --auth-key=${TAILSCALE_AUTH_KEY}"
  [ "$TAILSCALE_ADVERTISE_EXIT" = "true" ] && TS_CMD+=" --advertise-exit-node"
  TS_CMD+=" --reset"

  eval "$TS_CMD" || log_warn "Tailscale setup may need manual intervention"

  sleep 3
  tailscale status || true

  log_success "Tailscale configured"
fi

# ============================================================================
# NODE.JS TOOLS
# ============================================================================

log_info "Installing Node.js tools..."

if [ ! -f /usr/local/bin/n ]; then
  curl -fsSL https://raw.githubusercontent.com/tj/n/master/bin/n -o /usr/local/bin/n
  chmod +x /usr/local/bin/n
fi

/usr/local/bin/n lts 2>/dev/null || log_warn "Could not install Node.js LTS"

log_success "Node.js tools installed"

# ============================================================================
# SECURITY UPDATES
# ============================================================================

log_info "Configuring automatic security updates..."

cat >/etc/apt/apt.conf.d/20auto-upgrades <<EOF
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
EOF

log_success "Automatic updates configured"

# ============================================================================
# TIMEZONE
# ============================================================================

log_info "Configuring timezone..."

if [ -n "$TZ" ]; then
  timedatectl set-timezone "$TZ"
  log_info "  Timezone set to: $TZ"
else
  # Auto-detect via GeoIP
  TZ_GEO=$(curl -fsSL --max-time 5 https://ipapi.co/timezone 2>/dev/null || curl -fsSL --max-time 5 https://ipinfo.io/timezone 2>/dev/null || true)
  if [ -n "$TZ_GEO" ]; then
    timedatectl set-timezone "$TZ_GEO"
    log_info "  Timezone auto-detected: $TZ_GEO"
  else
    log_warn "Could not determine timezone"
  fi
fi

log_success "Timezone configured"

# ============================================================================
# NOMAD & CONSUL
# ============================================================================

if [ "$ENABLE_NOMAD" = "true" ] || [ "$ENABLE_CONSUL" = "true" ]; then
  log_info "Installing HashiCorp tools..."

  # Add HashiCorp repository
  rm -f /usr/share/keyrings/hashicorp-archive-keyring.gpg 2>/dev/null || true
  apt-get update -qq
  apt-get install -y -qq wget gpg coreutils

  log_info "  Adding HashiCorp GPG key..."
  wget -4 -qO- https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg 2>/dev/null

  echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" |
    tee /etc/apt/sources.list.d/hashicorp.list >/dev/null

  apt-get update -qq

  # Install requested packages
  PACKAGES=()
  [ "$ENABLE_NOMAD" = "true" ] && PACKAGES+=("nomad")
  [ "$ENABLE_CONSUL" = "true" ] && PACKAGES+=("consul" "consul-cni")

  apt-get install -y -qq "${PACKAGES[@]}"

  # Install CNI plugins if Nomad is enabled
  if [ "$ENABLE_NOMAD" = "true" ] && [ ! -d /opt/cni/bin ]; then
    log_info "  Installing CNI plugins..."
    ARCH_CNI=$([ "$(uname -m)" = "aarch64" ] && echo "arm64" || echo "amd64")
    CNI_VERSION="v1.6.2"

    curl -fsSL -o /tmp/cni-plugins.tgz "https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-${ARCH_CNI}-${CNI_VERSION}.tgz"
    mkdir -p /opt/cni/bin
    tar -C /opt/cni/bin -xzf /tmp/cni-plugins.tgz
    rm -f /tmp/cni-plugins.tgz
  fi

  # Configure bridge networking
  if [ "$ENABLE_NOMAD" = "true" ]; then
    modprobe br_netfilter 2>/dev/null || true

    cat >/etc/sysctl.d/bridge.conf <<EOF
net.bridge.bridge-nf-call-arptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
    sysctl --system >/dev/null 2>&1 || true
  fi

  log_success "HashiCorp tools installed"
fi

# ============================================================================
# NOMAD CONFIGURATION
# ============================================================================

if [ "$ENABLE_NOMAD" = "true" ]; then
  log_info "Configuring Nomad..."

  mkdir -p /nomad/{data,log} /etc/nomad.d

  # Auto-discover Nomad servers if not provided
  if [ -z "$NOMAD_SERVERS" ] && command -v tailscale >/dev/null 2>&1; then
    DISCOVERED_SERVERS=$(tailscale status --json 2>/dev/null | jq -r '.Peer[].TailscaleIPs[0]' 2>/dev/null || true)
    if [ -n "$DISCOVERED_SERVERS" ]; then
      NOMAD_SERVERS=$(echo "$DISCOVERED_SERVERS" | tr '\n' ',' | sed 's/,$//')
    fi
  fi

  # Build retry_join array
  RETRY_JOIN_CONFIG=""
  if [ -n "$NOMAD_SERVERS" ]; then
    IFS=',' read -ra SERVER_ARRAY <<<"$NOMAD_SERVERS"
    RETRY_JOIN_CONFIG="    retry_join = ["
    for server in "${SERVER_ARRAY[@]}"; do
      server=$(echo "$server" | xargs)
      [ -n "$server" ] && RETRY_JOIN_CONFIG="${RETRY_JOIN_CONFIG}\n      \"${server}:4647\","
    done
    RETRY_JOIN_CONFIG="${RETRY_JOIN_CONFIG%,}\n    ]"
  fi

  # Get private IP
  PRIVATE_IP=$(ip route get 1.1.1.1 | awk '{print $7; exit}' 2>/dev/null || echo "{{ GetPrivateIP }}")

  cat >/etc/nomad.d/nomad.hcl <<EOF
datacenter = "${NOMAD_DATACENTER}"
data_dir = "/nomad/data/"
bind_addr = "0.0.0.0"
log_level = "INFO"
log_json = true
log_file = "/nomad/log/nomad.log"
log_rotate_bytes = 10485760
log_rotate_max_files = 5

server {
  enabled = true
  bootstrap_expect = ${NOMAD_BOOTSTRAP_EXPECT}
$([ -n "$RETRY_JOIN_CONFIG" ] && echo -e "  server_join {\n$RETRY_JOIN_CONFIG\n  }")
}

client {
  enabled = true
  node_class = "${NOMAD_NODE_CLASS}"
}

advertise {
  http = "${PRIVATE_IP}:4646"
  rpc  = "${PRIVATE_IP}:4647"
  serf = "${PRIVATE_IP}:4648"
}

$(
    [ "$ENABLE_CONSUL" = "true" ] && cat <<CONSUL
consul {
  address = "127.0.0.1:8500"
  auto_advertise = true
  server_service_name = "nomad"
  client_service_name = "nomad-client"
}
CONSUL
  )

telemetry {
  collection_interval = "1s"
  publish_allocation_metrics = true
}
EOF

  # Enable and start Nomad
  systemctl daemon-reload
  systemctl enable nomad
  systemctl restart nomad

  sleep 5
  nomad server members 2>/dev/null || true

  log_success "Nomad configured"
fi

# ============================================================================
# SWAP CONFIGURATION
# ============================================================================

log_info "Configuring swap..."

if [ ! -f "$SWAP_FILE" ]; then
  fallocate -l "$SWAP_SIZE" "$SWAP_FILE" 2>/dev/null || dd if=/dev/zero of="$SWAP_FILE" bs=1M count=$(echo "$SWAP_SIZE" | sed 's/G$//' | awk '{print $1*1024}') 2>/dev/null
  chmod 600 "$SWAP_FILE"
  mkswap "$SWAP_FILE"
  swapon "$SWAP_FILE"

  # Add to fstab if not present
  if ! grep -q "$SWAP_FILE" /etc/fstab; then
    echo "${SWAP_FILE} none swap sw 0 0" >>/etc/fstab
  fi

  log_success "Swap file created: $SWAP_SIZE"
else
  swapon "$SWAP_FILE" 2>/dev/null || log_info "  Swap already active"
fi

# ============================================================================
# CLEANUP
# ============================================================================

log_info "Cleaning up..."

# Remove old SSH config backups (older than 30 days)
find /etc/ssh -name "sshd_config.backup.*" -mtime +30 -delete 2>/dev/null || true

# Clean package manager cache
apt-get autoremove -y -qq
apt-get autoclean -y -qq

log_success "Cleanup complete"

# ============================================================================
# FINAL STATUS
# ============================================================================

log_info ""
log_info "========================================"
log_info "  VPS Bootstrap Complete!"
log_info "========================================"
log_info "  Hostname: ${HOSTNAME_SHORT}"
log_info "  FQDN: ${FQDN}"
[ "$ENABLE_DOCKER" = "true" ] && command -v docker >/dev/null 2>&1 && log_info "  Docker: $(docker --version 2>/dev/null)"
[ "$ENABLE_TAILSCALE" = "true" ] && command -v tailscale >/dev/null 2>&1 && log_info "  Tailscale: Connected"
[ "$ENABLE_NOMAD" = "true" ] && command -v nomad >/dev/null 2>&1 && log_info "  Nomad: $(nomad -v 2>/dev/null)"
log_info "========================================"
log_info ""

echo "VPS Server '${HOSTNAME_SHORT}' is ready!" | tee /var/log/cloud-init-complete.log

# Optional: start zero-SPOF compose stack after host bootstrap
if [ "${ENABLE_ZERO_SPOF_STACK:-false}" = "true" ] && [ -x "${SCRIPT_DIR}/start-zero-spof-stack.sh" ]; then
  log_info "Starting zero-SPOF Docker stack..."
  SKIP_HOST_BOOTSTRAP=true \
    REQUIRED_EXPLICIT_FILE="${REQUIRED_EXPLICIT_FILE:-/etc/required-explicit.env}" \
    bash "${SCRIPT_DIR}/start-zero-spof-stack.sh" || log_warn "Zero-SPOF stack start failed (see docker compose ps)"
fi
