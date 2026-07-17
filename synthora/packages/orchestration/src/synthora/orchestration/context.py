"""Runtime research context: resolved providers, limits, and event emission.

Passed to graphs through LangGraph's ``config["configurable"]["synthora_ctx"]``
so nodes stay pure functions of (state, context) — R-ODR-6.
"""

from __future__ import annotations

import json
import re
from dataclasses import dataclass, field
from typing import Any, Optional

from synthora.core.events import ProgressEvent, RunEventType
from synthora.core.models import RunConfig
from synthora.core.ports import ChatModel, EventSink, SearchEngine, SearchStrategy


@dataclass
class ResearchContext:
    run_id: str
    config: RunConfig
    # role-split models (R-ODR-5)
    planner: ChatModel
    researcher: ChatModel
    compressor: ChatModel
    writer: ChatModel
    critic: ChatModel
    engines: list[SearchEngine] = field(default_factory=list)
    strategy: Optional[SearchStrategy] = None
    event_sink: Optional[EventSink] = None
    # user steering messages injected mid-run (R-STORM-6)
    steering: list[str] = field(default_factory=list)

    async def emit(
        self,
        type_: RunEventType,
        message: str = "",
        *,
        node: Optional[str] = None,
        payload: Optional[dict[str, Any]] = None,
    ) -> None:
        if self.event_sink is None:
            return
        await self.event_sink(
            ProgressEvent(
                run_id=self.run_id,
                type=type_,
                message=message,
                node=node,
                payload=payload or {},
            )
        )


def get_ctx(config: dict) -> ResearchContext:
    ctx = config.get("configurable", {}).get("synthora_ctx")
    if ctx is None:
        raise RuntimeError(
            "ResearchContext missing: pass config={'configurable': {'synthora_ctx': ctx}}"
        )
    return ctx


def parse_json_response(text: str) -> Optional[Any]:
    """Extract a JSON object/array from an LLM response.

    Handles bare JSON, fenced ```json blocks, and JSON embedded in prose.
    Provider-agnostic structured output: works with models lacking native
    tool calling (e.g. small local models).
    """
    text = text.strip()
    fence = re.search(r"```(?:json)?\s*(.*?)```", text, flags=re.DOTALL)
    if fence:
        text = fence.group(1).strip()
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        pass
    # fall back to the first {...} or [...] span
    for opener, closer in (("{", "}"), ("[", "]")):
        start = text.find(opener)
        if start == -1:
            continue
        depth = 0
        for i in range(start, len(text)):
            if text[i] == opener:
                depth += 1
            elif text[i] == closer:
                depth -= 1
                if depth == 0:
                    try:
                        return json.loads(text[start : i + 1])
                    except json.JSONDecodeError:
                        break
    return None
