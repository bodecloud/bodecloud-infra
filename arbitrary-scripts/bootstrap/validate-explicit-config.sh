#!/usr/bin/env bash
# Validate that required explicit configuration is present before bootstrap/stack start.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib-zero-spof.sh
source "${SCRIPT_DIR}/lib-zero-spof.sh"

zero_spof_load_explicit
zero_spof_validate_explicit

zero_spof_log_ok "Explicit configuration valid"
zero_spof_log_info "DOMAIN=${DOMAIN} TS_HOSTNAME=${TS_HOSTNAME}"
[[ -n "${CLOUDFLARE_API_TOKEN:-}" ]] && zero_spof_log_info "CLOUDFLARE_API_TOKEN=(set)"
[[ -n "${TAILSCALE_AUTH_KEY:-}" ]] && zero_spof_log_info "TAILSCALE_AUTH_KEY=(set)"
