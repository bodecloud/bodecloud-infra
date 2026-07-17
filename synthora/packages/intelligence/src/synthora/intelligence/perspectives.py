"""Perspective discovery (R-STORM-1).

STORM mines perspectives from related articles; we mine them directly from
the research brief plus optional retrieved context, producing N expert
personas with distinct focus areas.
"""

from __future__ import annotations

from synthora.core.models import Perspective
from synthora.core.parsing import parse_json_response
from synthora.core.ports import ChatModel


class PerspectiveEngine:
    def __init__(self, llm: ChatModel) -> None:
        self.llm = llm

    async def discover(
        self, brief: str, *, count: int = 3, context: str = ""
    ) -> list[Perspective]:
        raw = await self.llm.complete(
            [
                {
                    "role": "system",
                    "content": (
                        f"Identify {count} distinct expert perspectives for "
                        "researching the topic. Each expert must examine the topic "
                        "from a genuinely different angle (e.g. practitioner, "
                        "skeptic, historian, economist, engineer).\n"
                        'Reply JSON: [{"name": "...", "description": "...", '
                        '"focus": "..."}]'
                    ),
                },
                {
                    "role": "user",
                    "content": brief + (f"\n\nBackground:\n{context}" if context else ""),
                },
            ]
        )
        parsed = parse_json_response(raw)
        perspectives: list[Perspective] = []
        if isinstance(parsed, list):
            for item in parsed[:count]:
                if isinstance(item, dict) and item.get("name"):
                    perspectives.append(
                        Perspective(
                            name=str(item["name"]),
                            description=str(item.get("description", "")),
                            focus=str(item.get("focus", "")),
                        )
                    )
        if not perspectives:  # deterministic fallback keeps pipelines alive
            perspectives = [
                Perspective(
                    name=f"Expert {i + 1}",
                    description="General domain expert",
                    focus=brief[:100],
                )
                for i in range(count)
            ]
        return perspectives

    async def generate_questions(
        self, perspective: Perspective, brief: str, *, count: int = 3
    ) -> list[str]:
        """Perspective-guided question asking (the core STORM move)."""
        raw = await self.llm.complete(
            [
                {
                    "role": "system",
                    "content": (
                        f"You are {perspective.name}: {perspective.description}. "
                        f"Your focus: {perspective.focus}.\n"
                        f"Ask {count} incisive research questions about the topic "
                        "that only someone with your perspective would think to ask. "
                        "Return one question per line."
                    ),
                },
                {"role": "user", "content": brief},
            ]
        )
        questions = [q.strip("-• ").strip() for q in raw.splitlines() if q.strip()]
        return questions[:count] or [brief]
