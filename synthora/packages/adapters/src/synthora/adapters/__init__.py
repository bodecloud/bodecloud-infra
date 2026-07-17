"""Synthora adapters: LLM providers, search engines, strategies, embeddings."""

from synthora.adapters.embeddings import (
    EmbeddingRegistry,
    HashEmbeddings,
    OllamaEmbeddings,
    OpenAIEmbeddings,
    embedding_registry,
)
from synthora.adapters.llm import LLMProviderRegistry, llm_registry
from synthora.adapters.search_engines import (
    SearchEngineRegistry,
    search_engine_registry,
)
from synthora.adapters.strategies import SearchStrategyRegistry, strategy_registry
from synthora.adapters.summarize import summarize_page

__all__ = [
    "EmbeddingRegistry",
    "HashEmbeddings",
    "LLMProviderRegistry",
    "OllamaEmbeddings",
    "OpenAIEmbeddings",
    "SearchEngineRegistry",
    "SearchStrategyRegistry",
    "embedding_registry",
    "llm_registry",
    "search_engine_registry",
    "strategy_registry",
    "summarize_page",
]
