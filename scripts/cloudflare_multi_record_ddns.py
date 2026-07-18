#!/usr/bin/env python3
"""Cloudflare multi-record DDNS (Module 5).

Owns *global* FQDNs only: ensures one A record per healthy node Traefik IP.
Node-direct records (*.$TS_HOSTNAME.$DOMAIN) remain owned by favonia cloudflare-ddns.

DNS plurality is not failover — Traefik peer-forward remains the failover path.

Usage:
  cloudflare_multi_record_ddns.py --dry-run --node-ips '{"n1":"1.2.3.4","n2":"1.2.3.5"}' \\
    --names whoami.example.org --domain example.org
  cloudflare_multi_record_ddns.py --loop   # production loop
  cloudflare_multi_record_ddns.py --health
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from typing import Any, Dict, List, Optional

COMMENT_PREFIX = "failover-multi-ddns"
CF_API = "https://api.cloudflare.com/client/v4"


def env(name: str, default: str = "") -> str:
    return os.environ.get(name, default).strip()


def cf_request(
    method: str,
    path: str,
    token: str,
    body: Optional[dict] = None,
) -> dict:
    url = f"{CF_API}{path}"
    data = None if body is None else json.dumps(body).encode()
    req = urllib.request.Request(
        url,
        data=data,
        method=method,
        headers={
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json",
        },
    )
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        payload = e.read().decode()
        raise RuntimeError(f"Cloudflare API {method} {path}: {e.code} {payload}") from e


def plan_records(
    names: List[str],
    node_ips: Dict[str, str],
) -> Dict[str, List[str]]:
    """Return {fqdn: [ip, ...]} sorted uniquely."""
    ips = sorted({ip for ip in node_ips.values() if ip})
    out: Dict[str, List[str]] = {}
    for name in names:
        name = name.strip().rstrip(".")
        if not name:
            continue
        out[name] = list(ips)
    return out


def list_a_records(token: str, zone_id: str, name: str) -> List[dict]:
    q = urllib.parse.urlencode({"type": "A", "name": name, "per_page": 100})
    data = cf_request("GET", f"/zones/{zone_id}/dns_records?{q}", token)
    return list(data.get("result") or [])


def sync_name(
    token: str,
    zone_id: str,
    fqdn: str,
    desired_ips: List[str],
    dry_run: bool,
) -> Dict[str, Any]:
    existing = list_a_records(token, zone_id, fqdn) if token and zone_id else []
    owned = [
        r
        for r in existing
        if str(r.get("comment") or "").startswith(COMMENT_PREFIX)
        or not r.get("comment")
    ]
    have = {r["content"]: r for r in owned}
    desired = set(desired_ips)
    actions: List[str] = []

    for ip in sorted(desired):
        if ip in have:
            continue
        actions.append(f"CREATE {fqdn} A {ip}")
        if not dry_run and token and zone_id:
            cf_request(
                "POST",
                f"/zones/{zone_id}/dns_records",
                token,
                {
                    "type": "A",
                    "name": fqdn,
                    "content": ip,
                    "ttl": 60,
                    "proxied": False,
                    "comment": f"{COMMENT_PREFIX} multi-A",
                },
            )

    for ip, rec in list(have.items()):
        if ip in desired:
            continue
        # Only delete records we tagged
        if not str(rec.get("comment") or "").startswith(COMMENT_PREFIX):
            continue
        actions.append(f"DELETE {fqdn} A {ip} id={rec.get('id')}")
        if not dry_run and token and zone_id:
            cf_request("DELETE", f"/zones/{zone_id}/dns_records/{rec['id']}", token)

    return {"fqdn": fqdn, "desired": desired_ips, "actions": actions}


def load_node_ips(path_or_json: str) -> Dict[str, str]:
    raw = path_or_json.strip()
    if raw.startswith("{"):
        return json.loads(raw)
    with open(raw, encoding="utf-8") as f:
        return json.load(f)


def parse_names(arg: str, domain: str) -> List[str]:
    if arg.strip():
        return [n.strip() for n in arg.split(",") if n.strip()]
    # Default: whoami.$DOMAIN when MULTI_DDNS_NAMES unset
    if domain:
        return [f"whoami.{domain}"]
    return []


def run_once(args: argparse.Namespace) -> dict:
    domain = args.domain or env("DOMAIN")
    names = parse_names(args.names or env("MULTI_DDNS_NAMES"), domain)
    node_ips = load_node_ips(args.node_ips) if args.node_ips else load_node_ips(
        env("NODE_IPS_JSON") or "{}"
    )
    if not node_ips and env("NODE_IPS_FILE"):
        node_ips = load_node_ips(env("NODE_IPS_FILE"))

    planned = plan_records(names, node_ips)
    token = args.token or env("CF_API_TOKEN") or env("CLOUDFLARE_DNS_API_TOKEN")
    zone_id = args.zone_id or env("CF_ZONE_ID") or env("CLOUDFLARE_ZONE_ID")
    dry = args.dry_run or not token or not zone_id

    results = []
    for fqdn, ips in planned.items():
        results.append(sync_name(token, zone_id, fqdn, ips, dry_run=dry))

    report = {
        "dry_run": dry,
        "domain": domain,
        "node_ips": node_ips,
        "planned": planned,
        "results": results,
    }
    print(json.dumps(report, indent=2))
    return report


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--dry-run", action="store_true")
    p.add_argument("--loop", action="store_true")
    p.add_argument("--interval", type=int, default=int(env("MULTI_DDNS_INTERVAL", "60") or "60"))
    p.add_argument("--health", action="store_true")
    p.add_argument("--domain", default="")
    p.add_argument("--names", default="", help="comma-separated global FQDNs")
    p.add_argument("--node-ips", default="", help="JSON object or path to JSON file")
    p.add_argument("--token", default="")
    p.add_argument("--zone-id", default="")
    args = p.parse_args()

    if args.health:
        print("ok")
        return 0

    if args.loop:
        while True:
            try:
                run_once(args)
            except Exception as exc:  # noqa: BLE001
                print(json.dumps({"error": str(exc)}), file=sys.stderr)
            time.sleep(max(10, args.interval))
    else:
        report = run_once(args)
        # dry-run success requires planned multi-A
        for fqdn, ips in (report.get("planned") or {}).items():
            if len(ips) < 1:
                print(f"ERROR: no IPs planned for {fqdn}", file=sys.stderr)
                return 1
        return 0


if __name__ == "__main__":
    raise SystemExit(main())
