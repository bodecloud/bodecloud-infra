"""U2: adapter registries, strategies with fakes."""

import pytest
from synthora.adapters import llm_registry, search_engine_registry, strategy_registry
from synthora.adapters.llm import OpenAICompatibleModel, strip_think_tags
from synthora.adapters.strategies import (
    FocusedIterationStrategy,
    SourceBasedStrategy,
    dedupe_results,
)
from synthora.core.models import SearchResult

from tests.conftest import FakeChatModel, FakeSearchEngine


def test_llm_registry_resolution():
    model = llm_registry.resolve("openai:gpt-4o-mini")
    assert isinstance(model, OpenAICompatibleModel)
    assert model.model == "gpt-4o-mini"
    # bare model string defaults to openai
    assert llm_registry.resolve("gpt-4o").model == "gpt-4o"
    with pytest.raises(KeyError):
        llm_registry.resolve("nonexistent:model")


def test_strip_think_tags():
    assert strip_think_tags("<think>internal</think>answer") == "answer"
    assert strip_think_tags("plain") == "plain"


def test_search_engine_registry():
    assert {"searxng", "tavily", "arxiv", "semantic_scholar", "none"} <= set(
        search_engine_registry.engines()
    )
    with pytest.raises(KeyError):
        search_engine_registry.resolve("missing")


async def test_null_engine():
    engine = search_engine_registry.resolve("none")
    assert await engine.search("anything") == []


def test_dedupe_prefers_higher_score():
    results = [
        SearchResult(url="https://x.com/a", score=0.2, title="low"),
        SearchResult(url="https://x.com/a/", score=0.9, title="high"),
        SearchResult(url="https://x.com/b", score=0.5, title="other"),
    ]
    unique = dedupe_results(results)
    assert len(unique) == 2
    assert unique[0].title == "high"


async def test_source_based_strategy_decomposes_and_merges():
    llm = FakeChatModel(responses=["query one\nquery two\nquery three"])
    engine = FakeSearchEngine()
    strategy = SourceBasedStrategy()
    results = await strategy.run("quantum error correction", engines=[engine], llm=llm)
    assert engine.queries == ["query one", "query two", "query three"]
    assert results and all(r.url.startswith("https://example.com") for r in results)


async def test_focused_iteration_refines_query():
    llm = FakeChatModel(responses=["refined query", "another refinement"])
    engine = FakeSearchEngine()
    strategy = FocusedIterationStrategy(iterations=2)
    await strategy.run("topic", engines=[engine], llm=llm)
    assert engine.queries[0] == "topic"
    assert engine.queries[1] == "refined query"


async def test_strategy_survives_engine_failure():
    class BrokenEngine:
        name = "broken"

        async def search(self, query, *, max_results=5):
            raise RuntimeError("boom")

    llm = FakeChatModel(responses=["q1\nq2"])
    good = FakeSearchEngine()
    strategy = SourceBasedStrategy()
    results = await strategy.run("t", engines=[BrokenEngine(), good], llm=llm)
    assert results  # good engine results still returned


def test_strategy_registry():
    assert {"source_based", "focused_iteration"} <= set(strategy_registry.strategies())
    assert strategy_registry.resolve("source_based").name == "source_based"
