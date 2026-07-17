"""Search engine adapters and registry (R-LDR-4).

Initial engines: SearXNG (meta-search), Tavily, arXiv, Semantic Scholar,
and ``none``. Engines are constructed lazily by name so tests can register
fakes without network access.
"""

from __future__ import annotations

import os
from typing import Callable, Optional

import httpx
from defusedxml import ElementTree

from synthora.core.models import SearchResult
from synthora.core.ports import SearchEngine

EngineFactory = Callable[[], SearchEngine]


class SearxngEngine:
    name = "searxng"

    def __init__(self, base_url: Optional[str] = None, timeout: float = 20.0) -> None:
        self.base_url = (
            base_url or os.environ.get("SEARXNG_URL") or "http://localhost:8080"
        ).rstrip("/")
        self.timeout = timeout

    async def search(self, query: str, *, max_results: int = 5) -> list[SearchResult]:
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.get(
                f"{self.base_url}/search",
                params={"q": query, "format": "json"},
            )
            resp.raise_for_status()
            data = resp.json()
        results = []
        for item in data.get("results", [])[:max_results]:
            results.append(
                SearchResult(
                    url=item.get("url", ""),
                    title=item.get("title", ""),
                    snippet=item.get("content", ""),
                    content=item.get("content", ""),
                    engine=self.name,
                    score=float(item.get("score", 0.0) or 0.0),
                )
            )
        return results


class TavilyEngine:
    name = "tavily"

    def __init__(self, api_key: Optional[str] = None, timeout: float = 30.0) -> None:
        self.api_key = api_key or os.environ.get("TAVILY_API_KEY", "")
        self.timeout = timeout

    async def search(self, query: str, *, max_results: int = 5) -> list[SearchResult]:
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.post(
                "https://api.tavily.com/search",
                json={
                    "api_key": self.api_key,
                    "query": query,
                    "max_results": max_results,
                    "include_raw_content": True,
                },
            )
            resp.raise_for_status()
            data = resp.json()
        return [
            SearchResult(
                url=item.get("url", ""),
                title=item.get("title", ""),
                snippet=item.get("content", ""),
                content=item.get("raw_content") or item.get("content", ""),
                engine=self.name,
                score=float(item.get("score", 0.0) or 0.0),
            )
            for item in data.get("results", [])[:max_results]
        ]


class ArxivEngine:
    """Academic search over the arXiv Atom API (R-PIPE-4)."""

    name = "arxiv"
    _ns = {"atom": "http://www.w3.org/2005/Atom"}

    def __init__(self, timeout: float = 30.0) -> None:
        self.timeout = timeout

    async def search(self, query: str, *, max_results: int = 5) -> list[SearchResult]:
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.get(
                "https://export.arxiv.org/api/query",
                params={
                    "search_query": f"all:{query}",
                    "max_results": max_results,
                    "sortBy": "relevance",
                },
            )
            resp.raise_for_status()
            text = resp.text
        root = ElementTree.fromstring(text)
        results = []
        for entry in root.findall("atom:entry", self._ns):
            title = (entry.findtext("atom:title", "", self._ns) or "").strip()
            summary = (entry.findtext("atom:summary", "", self._ns) or "").strip()
            link = entry.findtext("atom:id", "", self._ns) or ""
            results.append(
                SearchResult(
                    url=link,
                    title=title,
                    snippet=summary[:500],
                    content=summary,
                    engine=self.name,
                    metadata={
                        "published": entry.findtext("atom:published", "", self._ns),
                        "authors": [
                            a.findtext("atom:name", "", self._ns)
                            for a in entry.findall("atom:author", self._ns)
                        ],
                    },
                )
            )
        return results


class SemanticScholarEngine:
    name = "semantic_scholar"

    def __init__(self, api_key: Optional[str] = None, timeout: float = 30.0) -> None:
        self.api_key = api_key or os.environ.get("SEMANTIC_SCHOLAR_API_KEY", "")
        self.timeout = timeout

    async def search(self, query: str, *, max_results: int = 5) -> list[SearchResult]:
        headers = {"x-api-key": self.api_key} if self.api_key else {}
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.get(
                "https://api.semanticscholar.org/graph/v1/paper/search",
                params={
                    "query": query,
                    "limit": max_results,
                    "fields": "title,abstract,url,year,citationCount,authors",
                },
                headers=headers,
            )
            resp.raise_for_status()
            data = resp.json()
        return [
            SearchResult(
                url=item.get("url") or "",
                title=item.get("title") or "",
                snippet=(item.get("abstract") or "")[:500],
                content=item.get("abstract") or "",
                engine=self.name,
                score=float(item.get("citationCount", 0) or 0),
                metadata={
                    "year": item.get("year"),
                    "authors": [
                        a.get("name") for a in (item.get("authors") or [])
                    ],
                },
            )
            for item in data.get("data", [])[:max_results]
        ]


class NullEngine:
    """No-op engine for offline runs."""

    name = "none"

    async def search(self, query: str, *, max_results: int = 5) -> list[SearchResult]:
        return []


class SearchEngineRegistry:
    def __init__(self) -> None:
        self._factories: dict[str, EngineFactory] = {}

    def register(self, name: str, factory: EngineFactory) -> None:
        self._factories[name] = factory

    def engines(self) -> list[str]:
        return sorted(self._factories)

    def resolve(self, name: str) -> SearchEngine:
        if name not in self._factories:
            raise KeyError(
                f"unknown search engine '{name}' (known: {self.engines()})"
            )
        return self._factories[name]()

    def resolve_many(self, names: list[str]) -> list[SearchEngine]:
        return [self.resolve(n) for n in names]


search_engine_registry = SearchEngineRegistry()
search_engine_registry.register("searxng", SearxngEngine)
search_engine_registry.register("tavily", TavilyEngine)
search_engine_registry.register("arxiv", ArxivEngine)
search_engine_registry.register("semantic_scholar", SemanticScholarEngine)
search_engine_registry.register("none", NullEngine)
