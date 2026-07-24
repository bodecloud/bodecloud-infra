#!/usr/bin/env bash
# Start the default zero-SPOF edge stack with minimal explicit configuration.
#
# Usage:
#   sudo ./start-zero-spof-stack.sh              # full path: validate → generate → compose up
#   ./start-zero-spof-stack.sh --generate-only   # only write .env + secrets
#   ./start-zero-spof-stack.sh --no-constellation
#   ./start-zero-spof-stack.sh traefik autokuma  # subset override
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib-zero-spof.sh
source "${SCRIPT_DIR}/lib-zero-spof.sh"

GENERATE_ONLY=false
SKIP_CONSTELLATION=false
EXTRA_SERVICES=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --generate-only) GENERATE_ONLY=true; shift ;;
    --no-constellation) SKIP_CONSTELLATION=true; shift ;;
    -h|--help)
      cat <<EOF
Usage: $0 [options] [service ...]

Options:
  --generate-only       Run generate-implicit-env.sh only (no docker compose up)
  --no-constellation    Skip Constellation agent install/start
  -h, --help            Show this help

Default services (when none specified):
  ${ZERO_SPOF_DEFAULT_SERVICES[*]}

Explicit config: copy required-explicit.env.example → /etc/required-explicit.env
EOF
      exit 0
      ;;
    --*) zero_spof_log_error "Unknown option: $1"; exit 1 ;;
    *) EXTRA_SERVICES+=("$1"); shift ;;
  esac
done

zero_spof_load_explicit
bash "${SCRIPT_DIR}/validate-explicit-config.sh"
bash "${SCRIPT_DIR}/generate-implicit-env.sh"

if [[ "$GENERATE_ONLY" == "true" ]]; then
  zero_spof_log_ok "Generate-only mode; skipping compose and Constellation"
  exit 0
fi

ROOT_DIR="${ZERO_SPOF_REPO_ROOT:-$(zero_spof_repo_root)}"
# shellcheck source=/dev/null
source "${ROOT_DIR}/.env"

zero_spof_need_cmd docker

# Optional: host bootstrap (Docker, Tailscale) when run as root
if [[ "${EUID:-$(id -u)}" -eq 0 && "${SKIP_HOST_BOOTSTRAP:-false}" != "true" ]]; then
  if [[ -x "${SCRIPT_DIR}/dont-run-directly.sh" ]]; then
    zero_spof_log_info "Running host bootstrap (dont-run-directly.sh)..."
    ENABLE_ZERO_SPOF_STACK=false \
      BOOTSTRAP_HOSTNAME="${TS_HOSTNAME}" \
      bash "${SCRIPT_DIR}/dont-run-directly.sh" "${TS_HOSTNAME}" || \
      zero_spof_log_warn "Host bootstrap returned non-zero (continuing if Docker is available)"
  fi
fi

# Constellation agent (parallel control plane; does not replace Compose failover-agent)
if [[ "${ENABLE_CONSTELLATION:-true}" == "true" && "$SKIP_CONSTELLATION" != "true" ]]; then
  CONSTELLATION_DIR="${CONSTELLATION_DIR:-/opt/constellation}"
  CONSTELLATION_CONFIG="${CONSTELLATION_CONFIG:-${CONSTELLATION_DIR}/config.yaml}"
  REPO_ROOT="$(zero_spof_repo_root)"
  if [[ "${EUID:-$(id -u)}" -eq 0 ]]; then
    if [[ ! -x /usr/local/bin/constellation-agent && -f "${REPO_ROOT}/infra/scripts/install.sh" ]]; then
      zero_spof_log_info "Installing Constellation agent..."
      bash "${REPO_ROOT}/infra/scripts/install.sh" || zero_spof_log_warn "Constellation install failed"
    fi
    if command -v systemctl >/dev/null 2>&1; then
      mkdir -p /etc/systemd/system/constellation-agent.service.d
      cat >/etc/systemd/system/constellation-agent.service.d/override.conf <<EOF
[Service]
Environment=DOMAIN=${DOMAIN}
Environment=TS_HOSTNAME=${TS_HOSTNAME}
Environment=CONFIG_PATH=${CONFIG_PATH:-${ROOT_DIR}/volumes}
Environment=SECRETS_PATH=${SECRETS_PATH:-${ROOT_DIR}/secrets}
Environment=CONSTELLATION_CONFIG=${CONSTELLATION_CONFIG}
EOF
      systemctl daemon-reload
      systemctl enable constellation-agent.service 2>/dev/null || true
      systemctl restart constellation-agent.service 2>/dev/null || \
        zero_spof_log_warn "Constellation agent not started (build/install may be required)"
    fi
  else
    zero_spof_log_warn "Constellation install skipped (not root); run with sudo or install manually"
  fi
fi

# Docker networks referenced as external in root compose
for net in warp-nat-net; do
  if ! docker network inspect "$net" >/dev/null 2>&1; then
    zero_spof_log_info "Creating docker network ${net}..."
    docker network create --driver bridge \
      --opt "com.docker.network.bridge.name=br_${net//-/_}" \
      --opt com.docker.network.bridge.enable_ip_masquerade=false \
      --subnet "${WARP_NAT_NET_SUBNET:-10.0.2.0/24}" \
      --gateway "${WARP_NAT_NET_GATEWAY:-10.0.2.1}" \
      --attachable "$net" 2>/dev/null || true
  fi
done

ARCH_PROFILE="$(zero_spof_compose_arch_profile)"
SERVICES=("${ZERO_SPOF_DEFAULT_SERVICES[@]}")
if ((${#EXTRA_SERVICES[@]} > 0)); then
  SERVICES=("${EXTRA_SERVICES[@]}")
fi

COMPOSE_FILE="${ROOT_DIR}/docker-compose.yml"
zero_spof_log_info "Starting zero-SPOF services: ${SERVICES[*]}"

set +e
docker compose \
  --project-directory "$ROOT_DIR" \
  --env-file "${ROOT_DIR}/.env" \
  -f "$COMPOSE_FILE" \
  --profile main \
  --profile "$ARCH_PROFILE" \
  up -d --remove-orphans "${SERVICES[@]}"
rc=$?
set -e

if [[ "$rc" -ne 0 ]]; then
  zero_spof_log_warn "Some services may have failed to start (exit ${rc})"
  zero_spof_log_info "Check: docker compose --project-directory ${ROOT_DIR} ps"
  exit "$rc"
fi

zero_spof_log_ok "Zero-SPOF edge stack started on ${TS_HOSTNAME}.${DOMAIN}"
zero_spof_log_info "Traefik: https://traefik.${DOMAIN}  Failover: http://127.0.0.1:8082/healthz"
