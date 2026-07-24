#!/usr/bin/env bash
# Shared helpers for zero-SPOF bootstrap (sourced, not executed directly).
set -euo pipefail

ZERO_SPOF_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

zero_spof_script_dir() {
  echo "$ZERO_SPOF_LIB_DIR"
}

zero_spof_repo_root() {
  local script_dir
  script_dir="$(zero_spof_script_dir)"
  cd "${script_dir}/../.." && pwd
}

# Default HA edge stack — matches services already enabled in coolify-proxy + core compose.
ZERO_SPOF_DEFAULT_SERVICES=(
  dockerproxy-ro
  dockerproxy-rw
  crowdsec
  traefik
  nginx-traefik-extensions
  tinyauth
  logrotate-traefik
  autokuma
  watchtower
  cloudflare-ddns
  failover-agent
)

zero_spof_log() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*"
}

zero_spof_log_info()  { zero_spof_log "[*] $*"; }
zero_spof_log_ok()    { zero_spof_log "  [OK] $*"; }
zero_spof_log_warn()  { zero_spof_log "  [WARN] $*"; }
zero_spof_log_error() { zero_spof_log "  [ERROR] $*" >&2; }

zero_spof_need_cmd() {
  local cmd=$1
  command -v "$cmd" >/dev/null 2>&1 || {
    zero_spof_log_error "required command not found: $cmd"
    exit 1
  }
}

zero_spof_rand() {
  local len=${1:-32}
  openssl rand -base64 48 | tr -d '/+=' | cut -c1-"$len"
}

zero_spof_load_explicit() {
  local script_dir repo_root
  script_dir="$(zero_spof_script_dir)"
  repo_root="$(zero_spof_repo_root)"

  if [[ -n "${REQUIRED_EXPLICIT_FILE:-}" && -f "${REQUIRED_EXPLICIT_FILE}" ]]; then
    # shellcheck source=/dev/null
    source "${REQUIRED_EXPLICIT_FILE}"
  elif [[ -f /etc/required-explicit.env ]]; then
    # shellcheck source=/dev/null
    source /etc/required-explicit.env
  elif [[ -f "${script_dir}/required-explicit.env" ]]; then
    # shellcheck source=/dev/null
    source "${script_dir}/required-explicit.env"
  fi

  if [[ -f "${BOOTSTRAP_CONFIG_FILE:-/etc/bootstrap-config.env}" ]]; then
    # shellcheck source=/dev/null
    source "${BOOTSTRAP_CONFIG_FILE:-/etc/bootstrap-config.env}"
  fi

  ZERO_SPOF_REPO_ROOT="${ZERO_SPOF_REPO_ROOT:-${ROOT_DIR:-${repo_root}}}"
  DOMAIN="${DOMAIN:-}"
  TS_HOSTNAME="${TS_HOSTNAME:-${BOOTSTRAP_HOSTNAME:-$(hostname -s 2>/dev/null || hostname | cut -d. -f1)}}"
  ENABLE_TAILSCALE="${ENABLE_TAILSCALE:-true}"
  ENABLE_CONSTELLATION="${ENABLE_CONSTELLATION:-true}"
  ENABLE_ZERO_SPOF_STACK="${ENABLE_ZERO_SPOF_STACK:-true}"
  ACME_RESOLVER_EMAIL="${ACME_RESOLVER_EMAIL:-admin@${DOMAIN}}"
  TAILSCALE_LOGIN_SERVER="${TAILSCALE_LOGIN_SERVER:-https://headscale.${DOMAIN}}"
}

zero_spof_validate_explicit() {
  local missing=()
  [[ -n "${DOMAIN:-}" ]] || missing+=("DOMAIN")
  [[ -n "${TS_HOSTNAME:-}" ]] || missing+=("TS_HOSTNAME (or BOOTSTRAP_HOSTNAME)")

  if [[ -z "${CLOUDFLARE_API_TOKEN:-}" && -z "${CLOUDFLARE_API_KEY:-}" ]]; then
    missing+=("CLOUDFLARE_API_TOKEN (or CLOUDFLARE_API_KEY + CLOUDFLARE_EMAIL)")
  fi

  if [[ "${ENABLE_TAILSCALE:-true}" == "true" && -z "${TAILSCALE_AUTH_KEY:-}" ]]; then
    missing+=("TAILSCALE_AUTH_KEY")
  fi

  if ((${#missing[@]} > 0)); then
    zero_spof_log_error "Missing required explicit configuration:"
    printf '  - %s\n' "${missing[@]}" >&2
    zero_spof_log_error "Copy required-explicit.env.example → /etc/required-explicit.env and fill in values."
    exit 1
  fi
}

zero_spof_discover_external_ip() {
  if [[ -n "${EXTERNAL_IP:-}" ]]; then
    echo "$EXTERNAL_IP"
    return
  fi
  curl -4fsS --max-time 5 https://ifconfig.me/ip 2>/dev/null \
    || curl -4fsS --max-time 5 https://api.ipify.org 2>/dev/null \
    || hostname -I 2>/dev/null | awk '{print $1}' \
    || echo "127.0.0.1"
}

zero_spof_discover_tailscale_peers() {
  local self="${TS_HOSTNAME}"
  if ! command -v tailscale >/dev/null 2>&1; then
    return 0
  fi
  if ! tailscale status >/dev/null 2>&1; then
    return 0
  fi
  tailscale status --json 2>/dev/null | python3 - "$self" <<'PY' || true
import json, sys
self = sys.argv[1]
try:
    data = json.load(sys.stdin)
except Exception:
    raise SystemExit(0)
peers = []
for p in (data.get("Peer") or {}).values():
    name = (p.get("HostName") or p.get("DNSName") or "").split(".")[0]
    if name and name != self:
        peers.append(name)
print(",".join(sorted(set(peers))))
PY
}

zero_spof_compose_arch_profile() {
  case "$(uname -m)" in
    arm*|aarch64) echo arm ;;
    x86_64|amd64|x86) echo x86_64 ;;
    ppc64le|ppc64) echo ppc64le ;;
    riscv64) echo riscv64 ;;
    *) zero_spof_log_error "Unsupported architecture: $(uname -m)"; exit 1 ;;
  esac
}

zero_spof_ensure_dir() {
  mkdir -p "$@"
}

zero_spof_write_secret_if_missing() {
  local file=$1
  local value=${2:-}
  zero_spof_ensure_dir "$(dirname "$file")"
  if [[ ! -f "$file" ]]; then
    if [[ -n "$value" ]]; then
      printf '%s\n' "$value" >"$file"
    else
      zero_spof_rand 32 >"$file"
    fi
    chmod 600 "$file" 2>/dev/null || true
  fi
}
