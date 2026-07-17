"""Run executor: builds the research context, invokes the pipeline graph,
persists results, and streams progress (used by the worker process and by
integration tests)."""

from __future__ import annotations

import logging
from datetime import datetime, timezone

from synthora.adapters import llm_registry, search_engine_registry, strategy_registry
from synthora.core.events import ProgressEvent, RunEventType
from synthora.core.models import (
    Artifact,
    ArtifactKind,
    ResearchRun,
    RunStatus,
)
from synthora.orchestration import pipelines  # noqa: F401  (registers pipelines)
from synthora.orchestration.context import ResearchContext
from synthora.orchestration.registry import pipeline_registry
from synthora.persistence import (
    ArtifactRepository,
    CitationRepository,
    EventRepository,
    KnowledgeRepository,
    RunRepositorySQL,
)
from synthora.persistence.database import Database
from synthora.worker.queue import RedisJobQueue

logger = logging.getLogger("synthora.worker")


class RunCancelled(Exception):
    pass


def utcnow() -> datetime:
    return datetime.now(timezone.utc)


class RunExecutor:
    def __init__(
        self,
        db: Database,
        queue: RedisJobQueue,
        *,
        model_resolver=None,
        engine_resolver=None,
        strategy_resolver=None,
    ) -> None:
        self.db = db
        self.queue = queue
        self._resolve_model = model_resolver or llm_registry.resolve
        self._resolve_engines = engine_resolver or search_engine_registry.resolve_many
        self._resolve_strategy = strategy_resolver or strategy_registry.resolve
        self.runs = RunRepositorySQL(db)
        self.events = EventRepository(db)
        self.artifacts = ArtifactRepository(db)
        self.citations = CitationRepository(db)
        self.knowledge = KnowledgeRepository(db)

    def build_context(self, run: ResearchRun) -> ResearchContext:
        cfg = run.config

        async def sink(event: ProgressEvent) -> None:
            # cancellation checkpoint on every event boundary
            if await self.queue.is_cancelled(run.id):
                raise RunCancelled(run.id)
            # absorb steering messages pushed mid-run (R-STORM-6)
            for msg in await self.queue.drain_steering(run.id):
                ctx.steering.append(msg)
            await self.events.append(event)
            await self.queue.publish_event(run.id, event.to_wire())

        ctx = ResearchContext(
            run_id=run.id,
            config=cfg,
            planner=self._resolve_model(cfg.planner_model),
            researcher=self._resolve_model(cfg.researcher_model),
            compressor=self._resolve_model(cfg.compressor_model),
            writer=self._resolve_model(cfg.writer_model),
            critic=self._resolve_model(cfg.critic_model),
            engines=self._resolve_engines(cfg.search_engines),
            strategy=self._resolve_strategy(cfg.search_strategy),
            event_sink=sink,
        )
        return ctx

    async def _emit_status(self, run: ResearchRun, message: str = "") -> None:
        event = ProgressEvent(
            run_id=run.id,
            type=RunEventType.STATUS,
            message=message or run.status.value,
            payload={"status": run.status.value},
        )
        await self.events.append(event)
        await self.queue.publish_event(run.id, event.to_wire())

    async def execute(self, run_id: str, *, ctx: ResearchContext | None = None) -> ResearchRun:
        """Execute one research run to completion (or failure/cancellation).

        ``ctx`` can be injected for tests; production builds it from config.
        """
        run = await self.runs.get(run_id)
        if run is None:
            raise KeyError(f"run {run_id} not found")
        if await self.queue.is_cancelled(run.id):
            run.status = RunStatus.CANCELLED
            run.finished_at = utcnow()
            await self.runs.update(run)
            await self._emit_status(run)
            return run

        run.status = RunStatus.RUNNING
        run.started_at = utcnow()
        await self.runs.update(run)
        await self._emit_status(run, "Research started")

        ctx = ctx or self.build_context(run)
        graph = pipeline_registry.build(run.pipeline_id)
        try:
            result = await graph.ainvoke(
                {"question": run.question},
                config={
                    "configurable": {"synthora_ctx": ctx, "thread_id": run.id},
                    "recursion_limit": 150,
                },
            )
            await self._persist_result(run, result)
            run.brief = result.get("brief")
            run.status = RunStatus.COMPLETED
        except RunCancelled:
            run.status = RunStatus.CANCELLED
        except Exception as exc:  # surface real failures to the user
            logger.exception("run %s failed", run.id)
            run.error = f"{type(exc).__name__}: {exc}"
            run.status = RunStatus.FAILED
        run.finished_at = utcnow()
        await self.runs.update(run)
        final_type = (
            RunEventType.DONE
            if run.status == RunStatus.COMPLETED
            else RunEventType.ERROR
        )
        event = ProgressEvent(
            run_id=run.id,
            type=final_type,
            message=run.error or run.status.value,
            payload={"status": run.status.value},
        )
        await self.events.append(event)
        await self.queue.publish_event(run.id, event.to_wire())
        return run

    async def _persist_result(self, run: ResearchRun, result: dict) -> None:
        report = result.get("report", "")
        if report:
            await self.artifacts.save(
                Artifact(
                    run_id=run.id,
                    kind=ArtifactKind.REPORT_MARKDOWN,
                    content=report,
                )
            )
        outline = result.get("outline")
        if outline is not None:
            from synthora.intelligence.outline import outline_to_markdown

            await self.artifacts.save(
                Artifact(
                    run_id=run.id,
                    kind=ArtifactKind.OUTLINE,
                    content=outline_to_markdown(outline),
                )
            )
        citations = result.get("citations") or []
        if citations:
            unique = list({c.id: c for c in citations}.values())
            for c in unique:
                c.run_id = run.id
            await self.citations.save_many(unique)
        nodes = result.get("knowledge_nodes") or []
        edges = result.get("knowledge_edges") or []
        if nodes:
            await self.knowledge.save_map(run.id, nodes, edges)
        notes = result.get("notes") or []
        if notes:
            await self.artifacts.save(
                Artifact(
                    run_id=run.id,
                    kind=ArtifactKind.RAW_NOTES,
                    content="\n\n".join(notes),
                )
            )
