"""Search strategies (R-LDR-4): how queries are generated and results merged.

Ported concepts from Local Deep Research's strategy layer:

- ``source_based``: decompose the topic into sub-queries, fan out across
  engines, deduplicate by URL, rank by score.
- ``focused_iteration``: iterate query refinement — search, ask the LLM
  what is still missing, search again with a refined query.
"""

from __future__ import annotations

import asyncio
from typing import Callable

from synthora.core.models import SearchResult
from synthora.core.ports import ChatModel, SearchEngine

StrategyFactory = Callable[[], "object"]


def dedupe_results(results: list[SearchResult]) -> list[SearchResult]:
    seen: set[str] = set()
    unique: list[SearchResult] = []
    for r in sorted(results, key=lambda r: r.score, reverse=True):
        key = r.url.rstrip("/")
        if key and key not in seen:
            seen.add(key)
            unique.append(r)
    return unique


async def _fan_out(
    queries: list[str], engines: list[SearchEngine], per_query: int
) -> list[SearchResult]:
    tasks = [
        engine.search(q, max_results=per_query) for q in queries for engine in engines
    ]
    batches = await asyncio.gather(*tasks, return_exceptions=True)
    results: list[SearchResult] = []
    for batch in batches:
        if isinstance(batch, BaseException):
            continue  # one engine failing must not sink the strategy
        results.extend(batch)
    return results


class SourceBasedStrategy:
    """Decompose into sub-queries, fan out, dedupe, rank."""

    name = "source_based"

    async def run(
        self,
        topic: str,
        *,
        engines: list[SearchEngine],
        llm: ChatModel,
        max_results: int = 8,
    ) -> list[SearchResult]:
        raw = await llm.complete(
            [
                {
                    "role": "system",
                    "content": (
                        "Decompose the research topic into 3 focused web search "
                        "queries. Return one query per line, no numbering."
                    ),
                },
                {"role": "user", "content": topic},
            ]
        )
        queries = [q.strip("-• ").strip() for q in raw.splitlines() if q.strip()][:3]
        if not queries:
            queries = [topic]
        results = await _fan_out(queries, engines, per_query=max(2, max_results // 2))
        return dedupe_results(results)[:max_results]


class FocusedIterationStrategy:
    """Search, reflect on gaps, refine the query, search again."""

    name = "focused_iteration"

    def __init__(self, iterations: int = 2) -> None:
        self.iterations = iterations

    async def run(
        self,
        topic: str,
        *,
        engines: list[SearchEngine],
        llm: ChatModel,
        max_results: int = 8,
    ) -> list[SearchResult]:
        collected: list[SearchResult] = []
        query = topic
        for _ in range(max(1, self.iterations)):
            batch = await _fan_out([query], engines, per_query=max_results)
            collected.extend(batch)
            summary = "\n".join(
                f"- {r.title}: {r.snippet[:150]}" for r in dedupe_results(collected)[:8]
            )
            query = (
                await llm.complete(
                    [
                        {
                            "role": "system",
                            "content": (
                                "Given a research topic and findings so far, produce "
                                "ONE refined search query targeting the biggest "
                                "remaining information gap. Reply with the query only."
                            ),
                        },
                        {
                            "role": "user",
                            "content": f"Topic: {topic}\nFindings:\n{summary}",
                        },
                    ]
                )
            ).strip()
            if not query:
                break
        return dedupe_results(collected)[:max_results]


class SearchStrategyRegistry:
    def __init__(self) -> None:
        self._factories: dict[str, StrategyFactory] = {}

    def register(self, name: str, factory: StrategyFactory) -> None:
        self._factories[name] = factory

    def strategies(self) -> list[str]:
        return sorted(self._factories)

    def resolve(self, name: str):
        if name not in self._factories:
            raise KeyError(
                f"unknown search strategy '{name}' (known: {self.strategies()})"
            )
        return self._factories[name]()


strategy_registry = SearchStrategyRegistry()
strategy_registry.register("source_based", SourceBasedStrategy)
strategy_registry.register("focused_iteration", FocusedIterationStrategy)
