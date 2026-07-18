# Module 5 — Cloudflare multi-record DDNS

> **Status**: In scope — syncer + dry-run CI prove  
> **Date**: 2026-07-18  
> **Related**: STRATEGY Track 2 Module 5; CI mesh DNS remains CoreDNS

## Ownership

| Writer | Records |
|---|---|
| `scripts/cloudflare_multi_record_ddns.py` / `cloudflare-multi-ddns` | Global FQDNs — **N A records** (one per healthy node IP), comment prefix `failover-multi-ddns` |
| favonia `cloudflare-ddns` | Node-direct `$TS_HOSTNAME.$DOMAIN` and `*.$TS_HOSTNAME.$DOMAIN` only |

DNS plurality ≠ failover. Traefik peer-forward remains the failover path.

## CI

```bash
arbitrary-scripts/failover-ci/prove-module5-ddns.sh
```

Always runs `--dry-run` against a 4-node fixture (expects 4 A IPs).  
Live upsert: `CF_LIVE_MULTI_DDNS=1` + `CF_API_TOKEN` + `CF_ZONE_ID`.

## Prod enable

```bash
# write volumes/placement/node-ips.json then:
docker compose --profile multi-ddns up -d cloudflare-multi-ddns
```
