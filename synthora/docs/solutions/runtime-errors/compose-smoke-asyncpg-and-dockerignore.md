---
title: Compose smoke failed without asyncpg and bloated build context
date: 2026-07-17
problem_type: runtime-error
component: synthora
tags: [docker, compose, asyncpg, dockerignore, smoke]
---

# Compose smoke: missing asyncpg + huge build context

## Symptoms

`scripts/smoke.sh` failed with `ModuleNotFoundError: No module named 'asyncpg'`
in the API/worker containers. An earlier attempt also transferred a multi‑MB
build context (`.venv` / `node_modules` included).

## Root cause

1. `asyncpg` was only an *optional* extra on `synthora-persistence`, so
   `uv sync --frozen --no-dev` in Docker did not install it while
   `SYNTHORA_DATABASE_URL` used `postgresql+asyncpg://`.
2. No `.dockerignore`, so local `.venv` and `node_modules` were sent to the
   builder.

## Fix

- Make `asyncpg` and `aiosqlite` hard dependencies of `synthora-persistence`.
- Add `synthora/.dockerignore` excluding `.venv`, `node_modules`, caches.
- Prefer a slim context (~310KB) and re-run `scripts/smoke.sh`.

## Verification

`bash scripts/smoke.sh` → health/ready/pipelines + web UI probe passed.
