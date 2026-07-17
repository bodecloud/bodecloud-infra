"""Synthora adapters: LLM providers, search engines, search strategies (R-LDR-4)."""

from synthora.adapters.llm import LLMProviderRegistry, llm_registry
from synthora.adapters.search_engines import (
    SearchEngineRegistry,
    search_engine_registry,
)
from synthora.adapters.strategies import SearchStrategyRegistry, strategy_registry

__all__ = [
    "LLMProviderRegistry",
    "SearchEngineRegistry",
    "SearchStrategyRegistry",
    "llm_registry",
    "search_engine_registry",
    "strategy_registry",
]
