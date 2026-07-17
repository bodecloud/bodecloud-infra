"""Runtime research context: resolved providers, limits, and event emission.

Passed to graphs through LangGraph's ``config["configurable"]["synthora_ctx"]``
so nodes stay pure functions of (state, context) — R-ODR-6.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any, Optional

from synthora.core.events import ProgressEvent, RunEventType
from synthora.core.models import RunConfig
from synthora.core.parsing import parse_json_response
from synthora.core.ports import ChatModel, EventSink, SearchEngine, SearchStrategy

__all__ = ["ResearchContext", "get_ctx", "parse_json_response"]


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
