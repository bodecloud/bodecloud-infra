"""MCP client bridge: load tools from servers listed in RunConfig.extra["mcp"].

Preferred path uses ``langchain-mcp-adapters`` when installed. Otherwise a
minimal HTTP JSON-RPC stub talks to MCP-ish ``tools/list`` / ``tools/call``
endpoints (including Synthora's own REST MCP surface).
"""

from __future__ import annotations

import logging
from dataclasses import dataclass, field
from typing import Any, Optional

import httpx

logger = logging.getLogger("synthora.adapters.mcp")


@dataclass
class MCPTool:
    """Callable tool handle returned to the researcher ReAct loop."""

    name: str
    description: str = ""
    input_schema: dict[str, Any] = field(default_factory=dict)
    _server_url: str = ""
    _transport: str = "http"
    _callable: Any = None

    async def ainvoke(self, arguments: Optional[dict[str, Any]] = None) -> str:
        args = arguments or {}
        if self._callable is not None:
            result = self._callable(args)
            if hasattr(result, "__await__"):
                result = await result
            return str(result)
        return await _http_tools_call(self._server_url, self.name, args)


async def load_mcp_tools(mcp_config: Optional[dict[str, Any]]) -> list[MCPTool]:
    """Load tools from ``{"servers": [{"url": "...", "transport": "http"}]}``."""
    if not mcp_config:
        return []
    servers = mcp_config.get("servers") or []
    if not isinstance(servers, list) or not servers:
        return []

    # Prefer langchain-mcp-adapters when available.
    try:
        tools = await _load_via_langchain(servers)
        if tools:
            return tools
    except Exception as exc:  # noqa: BLE001 — fall back to stub
        logger.debug("langchain-mcp-adapters unavailable: %s", exc)

    tools: list[MCPTool] = []
    for server in servers:
        if not isinstance(server, dict):
            continue
        url = str(server.get("url") or "").rstrip("/")
        transport = str(server.get("transport") or "http")
        if not url:
            continue
        try:
            listed = await _http_tools_list(url)
        except Exception as exc:  # noqa: BLE001
            logger.warning("MCP list failed for %s: %s", url, exc)
            continue
        for item in listed:
            name = str(item.get("name") or "")
            if not name:
                continue
            tools.append(
                MCPTool(
                    name=name,
                    description=str(item.get("description") or ""),
                    input_schema=dict(item.get("inputSchema") or item.get("input_schema") or {}),
                    _server_url=url,
                    _transport=transport,
                )
            )
    return tools


async def _load_via_langchain(servers: list[dict]) -> list[MCPTool]:
    from langchain_mcp_adapters.client import MultiServerMCPClient  # type: ignore

    connections: dict[str, dict[str, Any]] = {}
    for i, server in enumerate(servers):
        url = str(server.get("url") or "").rstrip("/")
        if not url:
            continue
        transport = str(server.get("transport") or "sse")
        key = str(server.get("name") or f"server_{i}")
        if transport in ("http", "streamable_http", "streamable-http"):
            connections[key] = {"url": url, "transport": "streamable_http"}
        else:
            connections[key] = {"url": url, "transport": transport or "sse"}
    if not connections:
        return []
    client = MultiServerMCPClient(connections)
    lc_tools = await client.get_tools()
    out: list[MCPTool] = []
    for tool in lc_tools:
        name = getattr(tool, "name", "") or ""
        desc = getattr(tool, "description", "") or ""
        schema = getattr(tool, "args_schema", None)
        input_schema: dict[str, Any] = {}
        if schema is not None and hasattr(schema, "model_json_schema"):
            input_schema = schema.model_json_schema()
        out.append(
            MCPTool(
                name=name,
                description=desc,
                input_schema=input_schema,
                _callable=getattr(tool, "ainvoke", None) or getattr(tool, "invoke", None),
            )
        )
    return out


async def _http_tools_list(base_url: str) -> list[dict[str, Any]]:
    """Call Synthora-style REST or JSON-RPC tools/list."""
    async with httpx.AsyncClient(timeout=20.0) as client:
        # Synthora REST surface
        resp = await client.post(f"{base_url}/api/v1/mcp/tools/list", json={})
        if resp.status_code < 400:
            data = resp.json()
            tools = data.get("tools") if isinstance(data, dict) else None
            if isinstance(tools, list):
                return tools
        # Minimal JSON-RPC
        resp = await client.post(
            base_url,
            json={
                "jsonrpc": "2.0",
                "id": 1,
                "method": "tools/list",
                "params": {},
            },
        )
        if resp.status_code >= 400:
            resp.raise_for_status()
        data = resp.json()
        result = data.get("result") if isinstance(data, dict) else None
        if isinstance(result, dict) and isinstance(result.get("tools"), list):
            return result["tools"]
        if isinstance(result, list):
            return result
        return []


async def _http_tools_call(
    base_url: str, name: str, arguments: dict[str, Any]
) -> str:
    async with httpx.AsyncClient(timeout=60.0) as client:
        resp = await client.post(
            f"{base_url}/api/v1/mcp/tools/call",
            json={"name": name, "arguments": arguments},
        )
        if resp.status_code < 400:
            data = resp.json()
            if isinstance(data, dict):
                if "content" in data:
                    return str(data["content"])
                if "result" in data:
                    return str(data["result"])
            return str(data)
        resp = await client.post(
            base_url,
            json={
                "jsonrpc": "2.0",
                "id": 1,
                "method": "tools/call",
                "params": {"name": name, "arguments": arguments},
            },
        )
        resp.raise_for_status()
        data = resp.json()
        if isinstance(data, dict):
            if "result" in data:
                return str(data["result"])
            if "error" in data:
                return f"error: {data['error']}"
        return str(data)
