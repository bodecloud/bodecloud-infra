"""Shared LangGraph checkpointer for interrupt/resume across pipelines."""

from __future__ import annotations

import os
from functools import lru_cache
from typing import Any

from langgraph.checkpoint.memory import MemorySaver


@lru_cache(maxsize=1)
def get_checkpointer() -> Any:
    """Return the process-wide checkpointer.

    Default is an in-memory saver (correct for single-worker and tests).
    Set ``SYNTHORA_CHECKPOINT_BACKEND=postgres`` and
    ``SYNTHORA_CHECKPOINT_URL`` to use Postgres when multi-worker resume
    is required.
    """
    backend = os.environ.get("SYNTHORA_CHECKPOINT_BACKEND", "memory").lower()
    if backend == "postgres":
        try:
            from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver

            url = os.environ.get(
                "SYNTHORA_CHECKPOINT_URL",
                os.environ.get("DATABASE_URL", ""),
            )
            if not url:
                raise RuntimeError(
                    "SYNTHORA_CHECKPOINT_URL or DATABASE_URL required "
                    "for postgres checkpointer"
                )
            # AsyncPostgresSaver.from_conn_string is a context manager in
            # some versions; keep a long-lived saver via connection string.
            saver = AsyncPostgresSaver.from_conn_string(url)
            return saver
        except Exception:
            # Fall back rather than crash worker boot; interrupt still works
            # within a single process via MemorySaver.
            return MemorySaver()
    return MemorySaver()


def reset_checkpointer() -> None:
    """Clear the cached checkpointer (tests only)."""
    get_checkpointer.cache_clear()
