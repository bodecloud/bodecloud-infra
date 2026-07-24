#!/usr/bin/env bash
# Auto-generate implicit .env, secrets, Constellation config, and placement paths.
# Requires explicit config (see required-explicit.env.example).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib-zero-spof.sh
source "${SCRIPT_DIR}/lib-zero-spof.sh"

zero_spof_load_explicit
zero_spof_validate_explicit
zero_spof_need_cmd openssl
zero_spof_need_cmd python3

REPO_ROOT="$(zero_spof_repo_root)"
ROOT_DIR="${ZERO_SPOF_REPO_ROOT:-$REPO_ROOT}"
SECRETS_DIR="${SECRETS_DIR:-${ROOT_DIR}/secrets}"
SECRETS_PATH="${SECRETS_PATH:-$SECRETS_DIR}"
CONFIG_PATH="${CONFIG_PATH:-${ROOT_DIR}/volumes}"
CERTS_DIR="${CERTS_DIR:-${ROOT_DIR}/certs}"
DATA_DIR="${DATA_DIR:-${ROOT_DIR}/data}"
STACK_NAME="${STACK_NAME:-${TS_HOSTNAME}}"
TZ="${TZ:-$(timedatectl show -p Timezone --value 2>/dev/null || echo Etc/UTC)}"
EXTERNAL_IP="$(zero_spof_discover_external_ip)"
ACME_RESOLVER_EMAIL="${ACME_RESOLVER_EMAIL:-admin@${DOMAIN}}"

# Failover peers / main host
if [[ -z "${FAILOVER_PEER_HOSTS:-}" ]]; then
  FAILOVER_PEER_HOSTS="$(zero_spof_discover_tailscale_peers)"
fi
if [[ -z "${FAILOVER_MAIN_HOST:-}" ]]; then
  if [[ -n "${FAILOVER_PEER_HOSTS:-}" ]]; then
    FAILOVER_MAIN_HOST="$(printf '%s\n%s' "$TS_HOSTNAME" "${FAILOVER_PEER_HOSTS//,/}" | tr ',' '\n' | sort | head -1)"
  else
    FAILOVER_MAIN_HOST="${TS_HOSTNAME}"
  fi
fi
FAILOVER_REPLICA_ENSURE="${FAILOVER_REPLICA_ENSURE:-$([ "$TS_HOSTNAME" = "$FAILOVER_MAIN_HOST" ] && echo true || echo false)}"

zero_spof_log_info "Generating implicit environment for ${TS_HOSTNAME}.${DOMAIN}"
zero_spof_log_info "ROOT_DIR=${ROOT_DIR} SECRETS_PATH=${SECRETS_PATH}"

# --- Directories -------------------------------------------------------------
zero_spof_ensure_dir \
  "$SECRETS_DIR" "$CONFIG_PATH" "$CERTS_DIR" "$DATA_DIR" \
  "${CONFIG_PATH}/traefik/dynamic" \
  "${CONFIG_PATH}/traefik/certs" \
  "${CONFIG_PATH}/traefik/logs" \
  "${CONFIG_PATH}/traefik/plugins-local" \
  "${CONFIG_PATH}/placement" \
  "${CONFIG_PATH}/traefik/crowdsec/var/log" \
  "${CONFIG_PATH}/traefik/crowdsec/data" \
  "${CONFIG_PATH}/traefik/crowdsec/etc/crowdsec" \
  "${CONFIG_PATH}/traefik/crowdsec/plugins" \
  "${CONFIG_PATH}/headscale/config" \
  "${CONFIG_PATH}/headscale/lib" \
  "${CONFIG_PATH}/headscale/run"

touch "${CONFIG_PATH}/traefik/crowdsec/var/log/auth.log" \
      "${CONFIG_PATH}/traefik/crowdsec/var/log/syslog"
touch "${CERTS_DIR}/acme.json"
chmod 600 "${CERTS_DIR}/acme.json" 2>/dev/null || true

[[ -f "${CONFIG_PATH}/placement/services.yaml" ]] || \
  cp -n "${REPO_ROOT}/placement/services.yaml.example" "${CONFIG_PATH}/placement/services.yaml" 2>/dev/null || \
  echo 'services: {}' >"${CONFIG_PATH}/placement/services.yaml"

mkdir -p "${CONFIG_PATH}/placement"
cat >"${CONFIG_PATH}/placement/node-ips.json" <<EOF
{"${TS_HOSTNAME}":"${EXTERNAL_IP}"}
EOF

# Merge CoreDNS IPs from failover-ci state or operator file (implicit split-DNS).
COREDNS_IPS_FILE="${COREDNS_IPS_FILE:-${CONFIG_PATH}/placement/coredns-ips.txt}"
if [[ -z "${COREDNS_IPS:-}" && -f "${COREDNS_IPS_FILE}" ]]; then
  COREDNS_IPS="$(grep -v '^#' "${COREDNS_IPS_FILE}" | paste -sd, - || true)"
fi
HEADSCALE_MAGIC_BASE_DOMAIN="${HEADSCALE_MAGIC_BASE_DOMAIN:-myscale.${DOMAIN}}"
HEADSCALE_SERVER_URL="${HEADSCALE_SERVER_URL:-https://headscale.${DOMAIN}}"
HEADSCALE_POLICY_PATH="${HEADSCALE_POLICY_PATH:-/etc/headscale/acl.hujson}"

# --- Headscale config (implicit DNS / ACL / extra records) -------------------
export DOMAIN TS_HOSTNAME CONFIG_PATH EXTERNAL_IP COREDNS_IPS COREDNS_IPS_FILE \
  HEADSCALE_MAGIC_BASE_DOMAIN HEADSCALE_SERVER_URL HEADSCALE_POLICY_PATH \
  HEADSCALE_HTTP_PORT="${HEADSCALE_HTTP_PORT:-8081}" \
  HEADSCALE_METRICS_PORT="${HEADSCALE_METRICS_PORT:-8080}" \
  HEADSCALE_STUN_PORT="${HEADSCALE_STUN_PORT:-3478}" \
  HEADSCALE_PREFIX_ALLOCATION="${HEADSCALE_PREFIX_ALLOCATION:-sequential}" \
  HEADSCALE_MAGIC_DNS="${HEADSCALE_MAGIC_DNS:-true}" \
  HEADSCALE_OVERRIDE_LOCAL_DNS="${HEADSCALE_OVERRIDE_LOCAL_DNS:-true}"
python3 "${SCRIPT_DIR}/render-headscale-config.py" "${CONFIG_PATH}/headscale/config/config.yaml"
zero_spof_log_ok "Wrote Headscale config ${CONFIG_PATH}/headscale/config/config.yaml"

# --- Secrets: discover compose references + ensure files -------------------
mapfile -t SECRET_NAMES < <(python3 - "$REPO_ROOT" <<'PY'
import re, sys
from pathlib import Path
root = Path(sys.argv[1])
files = [root / "docker-compose.yml"]
text = files[0].read_text()
for m in re.finditer(r"^- compose/(\S+\.yml)\s*$", text, re.M):
    files.append(root / "compose" / m.group(1))
pat = re.compile(r"\$\{SECRETS_PATH[^}]*\}/([A-Za-z0-9_.-]+\.txt)")
names = set()
for p in files:
    if p.is_file():
        names.update(pat.findall(p.read_text()))
names.update({
    "signing_secret.txt", "cf-api-token.txt", "cf-api-key.txt",
    "nginx-auth-api-key.txt", "tinyauth-secret.txt",
    "tinyauth-google-client-secret.txt", "tinyauth-github-client-secret.txt",
    "crowdsec-lapi-key.txt", "sudo-password.txt",
})
print("\n".join(sorted(names)))
PY
)

SUDO_PASSWORD="${SUDO_PASSWORD:-$(zero_spof_rand 24)}"
for name in "${SECRET_NAMES[@]}"; do
  [[ -n "$name" ]] || continue
  dest="${SECRETS_DIR}/${name}"
  case "$name" in
    cf-api-token.txt)
      zero_spof_write_secret_if_missing "$dest" "${CLOUDFLARE_API_TOKEN:-}"
      ;;
    cf-api-key.txt)
      zero_spof_write_secret_if_missing "$dest" "${CLOUDFLARE_API_KEY:-placeholder-not-used}"
      ;;
    sudo-password.txt)
      zero_spof_write_secret_if_missing "$dest" "$SUDO_PASSWORD"
      ;;
    tinyauth-google-client-secret.txt|tinyauth-github-client-secret.txt)
      zero_spof_write_secret_if_missing "$dest" "optional-oauth-not-configured"
      ;;
    *)
      zero_spof_write_secret_if_missing "$dest" ""
      ;;
  esac
done

# --- Write .env (implicit + explicit merge) ----------------------------------
ENV_FILE="${ROOT_DIR}/.env"
{
  cat <<EOF
# Generated by arbitrary-scripts/bootstrap/generate-implicit-env.sh
# Explicit inputs: DOMAIN, TS_HOSTNAME, CLOUDFLARE_*, TAILSCALE_AUTH_KEY
# Do not hand-edit generated secrets — re-run generate-implicit-env.sh instead.

DOMAIN=${DOMAIN}
TS_HOSTNAME=${TS_HOSTNAME}
STACK_NAME=${STACK_NAME}
ACME_RESOLVER_EMAIL=${ACME_RESOLVER_EMAIL}
EXTERNAL_IP=${EXTERNAL_IP}
TZ=${TZ}
PUID=${PUID:-1000}
PGID=${PGID:-1000}
UMASK=${UMASK:-002}

ROOT_DIR=${ROOT_DIR}
ROOT_PATH=${ROOT_DIR}
REPO_ROOT=${REPO_ROOT}
CONFIG_PATH=${CONFIG_PATH}
SECRETS_DIR=${SECRETS_DIR}
SECRETS_PATH=${SECRETS_PATH}
CERTS_DIR=${CERTS_DIR}
CERTS_PATH=${CERTS_DIR}
DATA_DIR=${DATA_DIR}
DATA_PATH=${DATA_DIR}
CREDENTIALS_DIRECTORY=${SECRETS_DIR}
SRC_DIR=${ROOT_DIR}/projects
SRC_PATH=${ROOT_DIR}/projects

SUDO_PASSWORD=${SUDO_PASSWORD}
WATCHTOWER_REPO_PASS=${WATCHTOWER_REPO_PASS:-${SUDO_PASSWORD}}

FAILOVER_ENABLED=true
FAILOVER_MAIN_HOST=${FAILOVER_MAIN_HOST}
FAILOVER_PEER_HOSTS=${FAILOVER_PEER_HOSTS}
FAILOVER_REMOTE_DOCKER_PORT=${FAILOVER_REMOTE_DOCKER_PORT:-2375}
FAILOVER_REMOTE_DOCKER_TLS=${FAILOVER_REMOTE_DOCKER_TLS:-false}
FAILOVER_REPLICA_ENSURE=${FAILOVER_REPLICA_ENSURE}
FAILOVER_MAX_LOCAL_RESTARTS=${FAILOVER_MAX_LOCAL_RESTARTS:-3}
FAILOVER_RECONCILE_SECONDS=${FAILOVER_RECONCILE_SECONDS:-30}
FAILOVER_COMPOSE_ENSURE_SERVICES=${FAILOVER_COMPOSE_ENSURE_SERVICES:-bolabaden-nextjs,autokuma}

CROWDSEC_BOUNCER_ENABLED=${CROWDSEC_BOUNCER_ENABLED:-true}
CROWDSEC_LAPI_KEY=${CROWDSEC_LAPI_KEY:-$(cat "${SECRETS_DIR}/crowdsec-lapi-key.txt" 2>/dev/null || zero_spof_rand 32)}

TAILSCALE_LOGIN_SERVER=${TAILSCALE_LOGIN_SERVER}
HEADSCALE_SERVER_URL=${HEADSCALE_SERVER_URL}
HEADSCALE_MAGIC_BASE_DOMAIN=${HEADSCALE_MAGIC_BASE_DOMAIN}
HEADSCALE_HTTP_PORT=${HEADSCALE_HTTP_PORT:-8081}
HEADSCALE_METRICS_PORT=${HEADSCALE_METRICS_PORT:-8080}
HEADSCALE_STUN_PORT=${HEADSCALE_STUN_PORT:-3478}
HEADSCALE_POLICY_PATH=${HEADSCALE_POLICY_PATH}
HEADSCALE_MAGIC_DNS=${HEADSCALE_MAGIC_DNS:-true}
HEADSCALE_OVERRIDE_LOCAL_DNS=${HEADSCALE_OVERRIDE_LOCAL_DNS:-true}
COREDNS_IPS=${COREDNS_IPS:-}
EOF
} >"${ENV_FILE}.generated"

# Preserve operator overrides from existing .env when keys are non-empty,
# except explicit keys which always reflect the current run.
EXPLICIT_ENV_KEYS="DOMAIN TS_HOSTNAME STACK_NAME ACME_RESOLVER_EMAIL TAILSCALE_LOGIN_SERVER HEADSCALE_SERVER_URL HEADSCALE_MAGIC_BASE_DOMAIN FAILOVER_MAIN_HOST FAILOVER_PEER_HOSTS FAILOVER_REPLICA_ENSURE"
if [[ -f "${ENV_FILE}" ]]; then
  python3 - "${ENV_FILE}" "${ENV_FILE}.generated" "$EXPLICIT_ENV_KEYS" <<'PY'
import sys
from pathlib import Path

existing = Path(sys.argv[1])
generated = Path(sys.argv[2])
explicit = set(sys.argv[3].split())
keep = {}
for line in existing.read_text().splitlines():
    if not line or line.lstrip().startswith("#") or "=" not in line:
        continue
    k, _, v = line.partition("=")
    k, v = k.strip(), v.strip()
    if k in explicit:
        continue
    if v and not v.startswith("$"):
        keep[k] = v
out = []
for line in generated.read_text().splitlines():
    if "=" in line and not line.lstrip().startswith("#"):
        k, _, _ = line.partition("=")
        k = k.strip()
        if k in keep:
            out.append(f"{k}={keep[k]}")
            continue
    out.append(line)
generated.write_text("\n".join(out) + "\n")
PY
fi
mv "${ENV_FILE}.generated" "${ENV_FILE}"
zero_spof_log_ok "Wrote ${ENV_FILE}"

# Iteratively seed any compose-required vars still missing
if command -v docker >/dev/null 2>&1; then
  seeded=0
  for _ in $(seq 1 50); do
    set +e
    err="$(docker compose --project-directory "$ROOT_DIR" --env-file "$ENV_FILE" -f "${ROOT_DIR}/docker-compose.yml" config 2>&1 >/dev/null)"
    set -e
    missing="$(echo "$err" | grep -oE 'required variable [A-Z_0-9]+' | head -1 | awk '{print $3}')" || true
    [[ -n "$missing" ]] || break
    val="$(zero_spof_rand 16)"
    if grep -q "^${missing}=" "$ENV_FILE"; then
      sed -i "s|^${missing}=.*|${missing}=${val}|" "$ENV_FILE"
    else
      echo "${missing}=${val}" >>"$ENV_FILE"
    fi
    seeded=$((seeded + 1))
  done
  [[ "$seeded" -gt 0 ]] && zero_spof_log_ok "Seeded ${seeded} compose-required variables"
fi

# --- Constellation config ----------------------------------------------------
CONSTELLATION_DIR="${CONSTELLATION_DIR:-/opt/constellation}"
CONSTELLATION_CONFIG="${CONSTELLATION_CONFIG:-${CONSTELLATION_DIR}/config.yaml}"
if [[ "${ENABLE_CONSTELLATION:-true}" == "true" ]]; then
  zero_spof_ensure_dir "$CONSTELLATION_DIR/data/raft/logs" \
    "$CONSTELLATION_DIR/data/raft/stable" \
    "$CONSTELLATION_DIR/data/raft/snapshots" \
    "$CONSTELLATION_DIR/volumes" \
    "$CONSTELLATION_DIR/secrets"
  cat >"${CONSTELLATION_CONFIG}" <<EOF
# Generated by generate-implicit-env.sh — merge with infra/config/examples/multi-node.yaml
domain: ${DOMAIN}
stack_name: ${STACK_NAME}
node_name: ${TS_HOSTNAME}

config_path: ${CONFIG_PATH}
secrets_path: ${SECRETS_DIR}
root_path: ${ROOT_DIR}
data_dir: ${DATA_DIR}

traefik:
  web_port: 80
  websecure_port: 443
  cert_resolver: letsencrypt
  http_provider_port: 8081

dns:
  provider: cloudflare
  domain: ${DOMAIN}

cluster:
  bind_port: 7946
  raft_port: 8300
  api_port: 8080
  priority: 100
EOF
  zero_spof_log_ok "Wrote Constellation config ${CONSTELLATION_CONFIG}"
fi

zero_spof_log_ok "Implicit environment generation complete"
