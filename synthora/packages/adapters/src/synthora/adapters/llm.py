"""LLM provider registry (R-LDR-4, R-ODR-5).

Model identifiers use ``provider:model`` strings, e.g. ``openai:gpt-4o``,
``ollama:llama3.1``, ``fake:scripted``. New providers register a factory
by name; every provider returns an object satisfying the ChatModel port.
"""

from __future__ import annotations

import os
from typing import Callable, Optional

import httpx
from synthora.core.ports import ChatModel

ProviderFactory = Callable[[str], ChatModel]


class OpenAICompatibleModel:
    """Chat model over any OpenAI-compatible endpoint (OpenAI, OpenRouter,
    vLLM, LM Studio, llama.cpp server)."""

    def __init__(
        self,
        model: str,
        *,
        api_key: Optional[str] = None,
        base_url: Optional[str] = None,
        timeout: float = 120.0,
    ) -> None:
        self.model = model
        self.api_key = api_key or os.environ.get("OPENAI_API_KEY", "")
        self.base_url = (
            base_url
            or os.environ.get("OPENAI_BASE_URL")
            or "https://api.openai.com/v1"
        ).rstrip("/")
        self.timeout = timeout

    async def complete(
        self,
        messages: list[dict[str, str]],
        *,
        temperature: float = 0.3,
        max_tokens: Optional[int] = None,
    ) -> str:
        payload: dict = {
            "model": self.model,
            "messages": messages,
            "temperature": temperature,
        }
        if max_tokens:
            payload["max_tokens"] = max_tokens
        headers = {"Content-Type": "application/json"}
        if self.api_key:
            headers["Authorization"] = f"Bearer {self.api_key}"
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.post(
                f"{self.base_url}/chat/completions", json=payload, headers=headers
            )
            resp.raise_for_status()
            data = resp.json()
        content = data["choices"][0]["message"]["content"] or ""
        return strip_think_tags(content)


class OllamaModel:
    """Chat model over a local Ollama server."""

    def __init__(
        self, model: str, *, base_url: Optional[str] = None, timeout: float = 300.0
    ) -> None:
        self.model = model
        self.base_url = (
            base_url or os.environ.get("OLLAMA_BASE_URL") or "http://localhost:11434"
        ).rstrip("/")
        self.timeout = timeout

    async def complete(
        self,
        messages: list[dict[str, str]],
        *,
        temperature: float = 0.3,
        max_tokens: Optional[int] = None,
    ) -> str:
        options: dict = {"temperature": temperature}
        if max_tokens:
            options["num_predict"] = max_tokens
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.post(
                f"{self.base_url}/api/chat",
                json={
                    "model": self.model,
                    "messages": messages,
                    "stream": False,
                    "options": options,
                },
            )
            resp.raise_for_status()
            data = resp.json()
        return strip_think_tags(data.get("message", {}).get("content", ""))


def strip_think_tags(text: str) -> str:
    """Remove <think>...</think> blocks emitted by reasoning models
    (mirrors Local Deep Research's think-tag wrapper)."""
    import re

    return re.sub(r"<think>.*?</think>", "", text, flags=re.DOTALL).strip()


class LLMProviderRegistry:
    def __init__(self) -> None:
        self._factories: dict[str, ProviderFactory] = {}

    def register(self, provider: str, factory: ProviderFactory) -> None:
        self._factories[provider] = factory

    def providers(self) -> list[str]:
        return sorted(self._factories)

    def resolve(self, model_id: str) -> ChatModel:
        """Resolve a ``provider:model`` identifier to a chat model."""
        provider, _, model = model_id.partition(":")
        if not model:
            provider, model = "openai", provider
        if provider not in self._factories:
            raise KeyError(
                f"unknown LLM provider '{provider}' (known: {self.providers()})"
            )
        return self._factories[provider](model)


llm_registry = LLMProviderRegistry()
llm_registry.register("openai", lambda m: OpenAICompatibleModel(m))
llm_registry.register(
    "openai-compatible",
    lambda m: OpenAICompatibleModel(m, base_url=os.environ.get("OPENAI_BASE_URL")),
)
llm_registry.register("ollama", lambda m: OllamaModel(m))
