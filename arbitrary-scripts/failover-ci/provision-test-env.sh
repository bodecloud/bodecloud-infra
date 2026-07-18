#!/usr/bin/env bash
# Sync repo into VMs, write per-node .env + placeholder secrets, disable Cloudflare DDNS.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"
load_env

[[ -f "${STATE_DIR}/node-ips.json" ]] || die "run provision-vms.sh first"
need tar
need python3

VM_REPO_PATH="${VM_REPO_PATH:-/opt/my-media-stack}"

# Bundle repo (exclude heavy/unneeded trees)
BUNDLE="${STATE_DIR}/repo-bundle.tgz"
log "creating repo bundle → $BUNDLE"
rm -f "${BUNDLE}"
tar -C "${REPO_ROOT}" \
  --exclude='.git' \
  --exclude='node_modules' \
  --exclude='**/node_modules' \
  --exclude='volumes' \
  --exclude='data' \
  --exclude='certs' \
  --exclude='src/zurg' \
  --exclude='src/hedgedoc' \
  --exclude='src/firecrawl' \
  --exclude='src/AIOStreams' \
  --exclude='projects' \
  --exclude='knowledgebase' \
  --exclude='arbitrary-scripts/failover-ci/state' \
  -czf "${BUNDLE}" \
  docker-compose.yml compose infra scripts placement package.json .env.example \
  arbitrary-scripts/failover-ci

# Auto-discover every ${SECRETS_PATH}/foo.txt referenced by root + included compose files
mapfile -t SECRET_NAMES < <(python3 - "$REPO_ROOT" <<'PY'
import re, sys
from pathlib import Path
root = Path(sys.argv[1])
files = [root / "docker-compose.yml"]
text = (root / "docker-compose.yml").read_text()
for m in re.finditer(r"^- compose/(\S+\.yml)\s*$", text, re.M):
    files.append(root / "compose" / m.group(1))
pat = re.compile(r"\$\{SECRETS_PATH[^}]*\}/([A-Za-z0-9_.-]+\.txt)")
names = set()
for p in files:
    if p.is_file():
        names.update(pat.findall(p.read_text()))
# Always include a few extras used by optional overlays / CI
names.update({"authentik-secret-key.txt", "gmail-app-password.txt"})
print("\n".join(sorted(names)))
PY
)
log "placeholder secrets: ${#SECRET_NAMES[@]} files"

for name in "${NODES[@]}"; do
  ip="$(node_ip_from_state "$name")"
  peers="$(peers_csv_for "$name")"
  log "provisioning env on $name ($ip) peers=$peers"

  vm_exec "$name" "sudo mkdir -p ${VM_REPO_PATH} /tmp/dev-secrets && sudo chown -R \$(whoami):\$(whoami) ${VM_REPO_PATH} || true"
  vm_transfer "$name" "${BUNDLE}" "/tmp/repo-bundle.tgz"
  vm_exec "$name" "tar -xzf /tmp/repo-bundle.tgz -C ${VM_REPO_PATH}"

  # Placeholder secrets (full list from compose discovery)
  secrets_list="${STATE_DIR}/secret-names.txt"
  printf '%s\n' "${SECRET_NAMES[@]}" > "${secrets_list}"
  vm_transfer "$name" "${secrets_list}" "/tmp/failover-ci-secret-names.txt"
  vm_exec "$name" "bash -lc '
    mkdir -p /tmp/dev-secrets ~/.docker
    while IFS= read -r f; do
      [[ -n \"\$f\" ]] || continue
      echo ci-placeholder > \"/tmp/dev-secrets/\$f\"
    done < /tmp/failover-ci-secret-names.txt
    echo {} > ~/.docker/config.json
  '"

  # Per-node .env
  local_env="${STATE_DIR}/env-${name}.env"
  cat > "${local_env}" <<EOF
DOMAIN=${DOMAIN}
TS_HOSTNAME=${name}
STACK_NAME=${name}
SECRETS_PATH=/tmp/dev-secrets
SECRETS_DIR=/tmp/dev-secrets
ROOT_DIR=${VM_REPO_PATH}
ROOT_PATH=${VM_REPO_PATH}
CONFIG_PATH=${VM_REPO_PATH}/volumes
CERTS_DIR=${VM_REPO_PATH}/certs
CERTS_PATH=${VM_REPO_PATH}/certs
DATA_DIR=${VM_REPO_PATH}/data
DATA_PATH=${VM_REPO_PATH}/data
FAILOVER_ENABLED=true
FAILOVER_MAIN_HOST=${MAIN_HOST}
FAILOVER_PEER_HOSTS=${peers}
FAILOVER_REMOTE_DOCKER_PORT=${FAILOVER_REMOTE_DOCKER_PORT:-2375}
FAILOVER_REPLICA_ENSURE=$([ "$name" = "${MAIN_HOST}" ] && echo true || echo false)
OSVC_INGRESS_SYNC_DISABLE=1
CLOUDFLARE_DDNS_ENABLED=false
DISABLE_CLOUDFLARE_DDNS=true
CLOUDFLARE_MULTI_DDNS_ENABLE=false
CLOUDFLARE_EMAIL=ci@${DOMAIN}
CLOUDFLARE_ZONE_ID=ci
CF_ZONE_ID=ci
CF_API_KEY=ci
CLOUDFLARE_DNS_API_TOKEN=ci
CLOUDFLARE_ZONE_API_TOKEN=ci
ACME_RESOLVER_EMAIL=ci@${DOMAIN}
EXTERNAL_IP=${ip}
EOF
  # Append common placeholders from test-stack.yml style
  cat >> "${local_env}" <<'EOF'
SEARXNG_SECRET=ci
REDIS_PASSWORD=ci
SUDO_PASSWORD=ci
NGINX_AUTH_API_KEY=ci
PROWLARR_API_KEY=ci
REALDEBRID_TOKEN=ci
REALDEBRID_API_KEY=ci
TINYAUTH_GOOGLE_CLIENT_ID=ci
TINYAUTH_GITHUB_CLIENT_ID=ci
TINYAUTH_SECRET=ci
TINYAUTH_USERS=ci
TINYAUTH_OAUTH_WHITELIST=ci
TINYAUTH_GOOGLE_CLIENT_SECRET=ci
TINYAUTH_GITHUB_CLIENT_SECRET=ci
CROWDSEC_LAPI_KEY=ci
CROWDSEC_MACHINE_ID=ci
CROWDSEC_BOUNCER_ENABLED=false
AIOSTREAMS_SECRET_KEY=ci
AIOSTREAMS_ADDON_PASSWORD=ci
STREMTHRU_PROXY_AUTH=ci
STREMTHRU_STORE_AUTH=ci
STREMTHRU_AUTH_ADMIN=ci
LITELLM_MASTER_KEY=ci
FIRECRAWL_API_KEY=ci
FIRECRAWL_BULL_AUTH_KEY=ci
MCPO_API_KEY=ci
MCP_AUTH_TOKEN=ci
OPEN_WEBUI_SECRET_KEY=ci
GF_SECURITY_ADMIN_PASSWORD=ci
GRAFANA_PASSWORD=ci
MEDIAFUSION_SECRET_KEY=ci
MEDIAFLOW_PROXY_API_PASSWORD=ci
AUTOKUMA__KUMA__USERNAME=ci
AUTOKUMA__KUMA__PASSWORD=ci
JACKETT_API_KEY=ci
EOF

  vm_transfer "$name" "${local_env}" "${VM_REPO_PATH}/.env"
  # CI compose overlays always present under compose/
  vm_exec "$name" "mkdir -p ${VM_REPO_PATH}/compose ${VM_REPO_PATH}/arbitrary-scripts/failover-ci/compose
    cp ${VM_REPO_PATH}/arbitrary-scripts/failover-ci/compose/docker-compose.ci-probes.yml ${VM_REPO_PATH}/compose/ 2>/dev/null || true
    cp ${VM_REPO_PATH}/arbitrary-scripts/failover-ci/compose/docker-compose.ci-stack.yml ${VM_REPO_PATH}/compose/ 2>/dev/null || true
    cp ${VM_REPO_PATH}/arbitrary-scripts/failover-ci/compose/Dockerfile.ci-probe ${VM_REPO_PATH}/arbitrary-scripts/failover-ci/compose/ 2>/dev/null || true"
done

log "provision-test-env complete"
