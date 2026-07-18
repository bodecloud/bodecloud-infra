#!/usr/bin/env bash
# Provision 4 nodes: Multipass → QEMU/KVM (/dev/kvm) → DinD nested fallback.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

BACKEND="$(detect_backend)"
[[ "$BACKEND" != "none" ]] || die "need Multipass, QEMU+/dev/kvm, or Docker (DinD fallback). See README.md"
save_backend "$BACKEND"
log "backend=$BACKEND"

CLOUD_INIT="${SCRIPT_DIR}/cloud-init/node.yaml"
SSH_DIR="${STATE_DIR}/ssh"
IMG_DIR="${STATE_DIR}/images"
mkdir -p "${SSH_DIR}" "${IMG_DIR}"

provision_multipass() {
  local name
  for name in "${NODES[@]}"; do
    if multipass info "$name" >/dev/null 2>&1; then
      log "VM $name already exists"
      continue
    fi
    log "launching $name (cpus=${VM_CPUS} mem=${VM_MEM}M disk=${VM_DISK})"
    multipass launch "${MULTIPASS_IMAGE}" \
      --name "$name" \
      --cpus "${VM_CPUS}" \
      --memory "${VM_MEM}M" \
      --disk "${VM_DISK}" \
      --cloud-init "${CLOUD_INIT}"
  done
}

ensure_ssh_key() {
  if [[ ! -f "${SSH_DIR}/id_ed25519" ]]; then
    ssh-keygen -t ed25519 -N "" -f "${SSH_DIR}/id_ed25519" -C "failover-ci"
  fi
}

download_cloud_image() {
  local img="${IMG_DIR}/jammy-server-cloudimg-amd64.img"
  if [[ -f "$img" ]]; then
    echo "$img"
    return
  fi
  local url="${QEMU_CLOUD_IMAGE_URL:-https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img}"
  log "downloading Ubuntu cloud image → $img"
  curl -fsSL -o "${img}.partial" "$url"
  mv "${img}.partial" "$img"
  echo "$img"
}

wait_dom_ip() {
  local name="$1"
  local ip=""
  local i
  for i in $(seq 1 60); do
    ip="$(virsh domifaddr "$name" --source lease 2>/dev/null | awk '/ipv4/ {print $4}' | head -1 | cut -d/ -f1 || true)"
    if [[ -n "$ip" ]]; then
      echo "$ip"
      return 0
    fi
    sleep 3
  done
  return 1
}

provision_qemu_node() {
  local name="$1"
  local base_img="$2"
  local disk="${IMG_DIR}/${name}.qcow2"
  local seed="${IMG_DIR}/${name}-seed.iso"
  local pub
  pub="$(cat "${SSH_DIR}/id_ed25519.pub")"

  if virsh dominfo "$name" >/dev/null 2>&1; then
    log "libvirt domain $name already exists"
    local ip
    ip="$(wait_dom_ip "$name" || true)"
    [[ -n "$ip" ]] || die "no IP for existing domain $name"
    cat > "${SSH_DIR}/${name}" <<EOF
SSH_USER=ubuntu
SSH_HOST=${ip}
SSH_KEY=${SSH_DIR}/id_ed25519
EOF
    return
  fi

  log "creating qcow2 for $name"
  qemu-img create -f qcow2 -F qcow2 -b "${base_img}" "${disk}" "${VM_DISK}"

  local ud="${IMG_DIR}/${name}-user-data"
  local md="${IMG_DIR}/${name}-meta-data"
  cat > "$ud" <<EOF
#cloud-config
hostname: ${name}
manage_etc_hosts: true
users:
  - name: ubuntu
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: sudo, docker
    shell: /bin/bash
    ssh_authorized_keys:
      - ${pub}
package_update: true
packages:
  - docker.io
  - docker-compose-v2
  - curl
  - wget
  - dnsutils
  - jq
  - python3
  - ca-certificates
runcmd:
  - systemctl enable --now docker
  - mkdir -p /opt/my-media-stack /opt/failover-ci /root/.docker
  - echo '{}' > /root/.docker/config.json
EOF
  cat > "$md" <<EOF
instance-id: ${name}
local-hostname: ${name}
EOF

  if command -v cloud-localds >/dev/null 2>&1; then
    cloud-localds "$seed" "$ud" "$md"
  else
    need genisoimage
    genisoimage -output "$seed" -volid cidata -joliet -rock "$ud" "$md"
  fi

  local mem_mib="${VM_MEM}"
  log "virt-install $name"
  virt-install \
    --name "$name" \
    --memory "$mem_mib" \
    --vcpus "${VM_CPUS}" \
    --disk "path=${disk},format=qcow2,bus=virtio" \
    --disk "path=${seed},device=cdrom" \
    --os-variant ubuntu22.04 \
    --import \
    --network network=default,model=virtio \
    --graphics none \
    --noautoconsole \
    --wait 0

  local ip
  ip="$(wait_dom_ip "$name")" || die "timeout waiting for DHCP on $name"
  log "$name IP=$ip"
  for i in $(seq 1 40); do
    if ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=5 \
      -i "${SSH_DIR}/id_ed25519" "ubuntu@${ip}" "true" 2>/dev/null; then
      break
    fi
    sleep 3
  done
  cat > "${SSH_DIR}/${name}" <<EOF
SSH_USER=ubuntu
SSH_HOST=${ip}
SSH_KEY=${SSH_DIR}/id_ed25519
EOF
}

provision_qemu() {
  [[ -e /dev/kvm ]] || die "QEMU backend requires /dev/kvm (use DinD fallback on nested hosts)"
  need virt-install
  need virsh
  need qemu-img
  ensure_ssh_key
  local base
  base="$(download_cloud_image)"
  local name
  for name in "${NODES[@]}"; do
    provision_qemu_node "$name" "$base"
  done
}

wait_dind_ready() {
  local name="$1"
  local i
  for i in $(seq 1 90); do
    if docker exec "$name" docker info >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done
  return 1
}

bootstrap_dind_node() {
  local name="$1"
  # Alpine DinD: tools + sudo shim (scripts assume ubuntu/sudo)
  docker exec "$name" sh -lc '
    set -eu
    if command -v apk >/dev/null 2>&1; then
      apk add --no-cache bash curl wget bind-tools jq python3 tar ca-certificates \
        git openssh-client docker-cli-compose iptables 2>/dev/null \
        || apk add --no-cache bash curl wget bind-tools jq python3 tar ca-certificates \
             git openssh-client iptables 2>/dev/null || true
    fi
    mkdir -p /opt/my-media-stack /opt/failover-ci /root/.docker
    echo "{}" > /root/.docker/config.json
    # Nested containers (Traefik) must SNAT to reach peer DinD IPs on eth0
    sysctl -w net.ipv4.ip_forward=1 >/dev/null 2>&1 || true
    sysctl -w fs.inotify.max_user_instances=1024 >/dev/null 2>&1 || true
    sysctl -w fs.inotify.max_user_watches=524288 >/dev/null 2>&1 || true
    iptables -t nat -C POSTROUTING -s 172.16.0.0/12 -o eth0 -j MASQUERADE 2>/dev/null \
      || iptables -t nat -A POSTROUTING -s 172.16.0.0/12 -o eth0 -j MASQUERADE
    iptables -C FORWARD -j ACCEPT 2>/dev/null || iptables -I FORWARD -j ACCEPT
    mkdir -p /dev/net
    [[ -c /dev/net/tun ]] || mknod /dev/net/tun c 10 200 || true
  ' || true
  # Install Tailscale outside nested quoting mess
  if ! docker exec "$name" command -v tailscale >/dev/null 2>&1; then
    local arch ts_arch ver
    arch="$(docker exec "$name" uname -m)"
    case "$arch" in
      aarch64|arm64) ts_arch=arm64 ;;
      *) ts_arch=amd64 ;;
    esac
    ver="${TAILSCALE_VERSION:-1.80.3}"
    log "installing Tailscale ${ver} (${ts_arch}) on $name"
    docker exec "$name" sh -lc "
      set -eu
      cd /tmp
      wget -q -O ts.tgz 'https://pkgs.tailscale.com/stable/tailscale_${ver}_${ts_arch}.tgz' \
        || curl -fsSL -o ts.tgz 'https://pkgs.tailscale.com/stable/tailscale_${ver}_${ts_arch}.tgz'
      tar -xzf ts.tgz
      cp -f tailscale_${ver}_${ts_arch}/tailscale tailscale_${ver}_${ts_arch}/tailscaled /usr/local/bin/
      chmod +x /usr/local/bin/tailscale /usr/local/bin/tailscaled
    " || warn "Tailscale install failed on $name (mesh may soft-fail)"
  fi
  if ! docker exec "$name" command -v sudo >/dev/null 2>&1; then
    docker exec -i "$name" sh -c 'cat > /usr/local/bin/sudo && chmod +x /usr/local/bin/sudo' <<'SUDOEOF'
#!/bin/sh
exec "$@"
SUDOEOF
  fi
}

provision_dind() {
  need docker
  docker info >/dev/null 2>&1 || die "docker daemon not available"

  if ! docker network inspect "${DIND_NETWORK}" >/dev/null 2>&1; then
    log "creating docker network ${DIND_NETWORK}"
    docker network create --driver bridge --attachable "${DIND_NETWORK}"
  fi

  local name
  for name in "${NODES[@]}"; do
    if docker inspect "$name" >/dev/null 2>&1; then
      log "DinD node $name already exists"
      if ! docker exec "$name" docker info >/dev/null 2>&1; then
        docker start "$name" >/dev/null || true
        wait_dind_ready "$name" || die "DinD $name not ready"
      fi
      bootstrap_dind_node "$name"
      continue
    fi
    log "launching DinD node $name (${DIND_IMAGE})"
    # Official entrypoint; args after image go to dockerd. TCP on private net only
    # (no -p publish). DOCKER_TLS_CERTDIR= disables TLS so peers use plain :2375.
    docker run -d \
      --name "$name" \
      --hostname "$name" \
      --privileged \
      --network "${DIND_NETWORK}" \
      --restart unless-stopped \
      --ulimit nofile=1048576:1048576 \
      -e DOCKER_TLS_CERTDIR= \
      -v "failover-ci-${name}-docker:/var/lib/docker" \
      "${DIND_IMAGE}" \
      dockerd \
      --host=unix:///var/run/docker.sock \
      --host=tcp://0.0.0.0:2375 \
      --tls=false \
      --default-ulimit nofile=1048576:1048576
    wait_dind_ready "$name" || die "timeout waiting for dockerd in $name"
    bootstrap_dind_node "$name"
    # inotify for Traefik file provider inside nested Docker
    docker exec "$name" sh -lc '
      sysctl -w fs.inotify.max_user_instances=1024 >/dev/null 2>&1 || true
      sysctl -w fs.inotify.max_user_watches=524288 >/dev/null 2>&1 || true
    ' || true
  done

  # Second pass: write /etc/hosts so peer names resolve inside each DinD
  write_node_ips
  cleanup_stale_dind_bridge_filters
  local n
  local hosts_file="${STATE_DIR}/dind-hosts"
  {
    echo "# failover-ci-peers"
    for n in "${NODES[@]}"; do
      echo "$(node_ip_from_state "$n") ${n}"
    done
  } > "${hosts_file}"
  for n in "${NODES[@]}"; do
    # sed -i on /etc/hosts fails in DinD (resource busy); append once only
    docker exec -i "$n" sh -lc '
      if grep -q failover-ci-peers /etc/hosts 2>/dev/null; then
        exit 0
      fi
      cat >> /etc/hosts
    ' < "${hosts_file}"
  done
}

case "$BACKEND" in
  multipass) provision_multipass; write_node_ips ;;
  qemu)      provision_qemu; write_node_ips ;;
  dind)      provision_dind ;;  # write_node_ips inside
  *)         die "unsupported backend $BACKEND" ;;
esac

log "nodes ready (backend=$BACKEND):"
for n in "${NODES[@]}"; do
  log "  $n → $(node_ip_from_state "$n")"
done
