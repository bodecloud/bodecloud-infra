#!/usr/bin/env bash
# Self-host smoke test: validate compose, bring the stack up, probe health.
set -euo pipefail
cd "$(dirname "$0")/.."

echo "==> validating compose file"
docker compose config --quiet

echo "==> building and starting stack"
docker compose up -d --build

cleanup() { docker compose down; }
trap cleanup EXIT

echo "==> waiting for API health"
for i in $(seq 1 60); do
  if curl -fsS "http://localhost:${SYNTHORA_API_PORT:-8000}/health" >/dev/null 2>&1; then
    break
  fi
  sleep 2
done
curl -fsS "http://localhost:${SYNTHORA_API_PORT:-8000}/health"
curl -fsS "http://localhost:${SYNTHORA_API_PORT:-8000}/ready"
curl -fsS "http://localhost:${SYNTHORA_API_PORT:-8000}/api/v1/pipelines"

echo "==> waiting for web UI"
curl -fsS "http://localhost:${SYNTHORA_WEB_PORT:-3000}/" >/dev/null

echo "smoke test passed"
