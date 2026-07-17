"""Synthora API gateway (R-LDR-3): REST + WebSocket + optional auth."""

from __future__ import annotations

import asyncio
import json
from contextlib import asynccontextmanager
from typing import Optional

import redis.asyncio as aioredis
import synthora.orchestration.pipelines  # noqa: F401  (registers pipelines)
from fastapi import Depends, FastAPI, HTTPException, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
from synthora.adapters import llm_registry, search_engine_registry, strategy_registry
from synthora.api.auth import (
    current_identity,
    hash_password,
    issue_token,
    verify_password,
)
from synthora.api.settings import settings
from synthora.core.models import (
    ArtifactKind,
    ResearchRun,
    RunConfig,
    RunStatus,
    User,
)
from synthora.orchestration.registry import pipeline_registry
from synthora.persistence import (
    ArtifactRepository,
    CitationRepository,
    EventRepository,
    KnowledgeRepository,
    RunRepositorySQL,
    UserRepository,
    WorkspaceRepository,
)
from synthora.persistence.database import Database
from synthora.worker.queue import RedisJobQueue, events_channel


@asynccontextmanager
async def lifespan(app: FastAPI):
    db = Database(settings.database_url)
    await db.create_all()
    redis = aioredis.from_url(settings.redis_url)
    app.state.db = db
    app.state.redis = redis
    app.state.queue = RedisJobQueue(redis)
    await WorkspaceRepository(db).ensure_default()
    yield
    await redis.aclose()
    await db.dispose()


app = FastAPI(title="Synthora", version="0.1.0", lifespan=lifespan)
app.add_middleware(
    CORSMiddleware,
    allow_origins=[o.strip() for o in settings.cors_origins.split(",")],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


def get_db() -> Database:
    return app.state.db


def get_queue() -> RedisJobQueue:
    return app.state.queue


# ---------------------------------------------------------------- schemas


class RegisterRequest(BaseModel):
    username: str = Field(min_length=3, max_length=64)
    password: str = Field(min_length=8)


class LoginRequest(BaseModel):
    username: str
    password: str


class StartResearchRequest(BaseModel):
    question: str = Field(min_length=3)
    pipeline_id: str = "deep_research"
    config: Optional[dict] = None


class SteerRequest(BaseModel):
    message: str = Field(min_length=1)


# ---------------------------------------------------------------- health


@app.get("/health")
async def health() -> dict:
    return {"status": "ok"}


@app.get("/ready")
async def ready() -> dict:
    try:
        await app.state.redis.ping()
        async with app.state.db.session() as _:
            pass
    except Exception as exc:
        raise HTTPException(status_code=503, detail=str(exc)) from exc
    return {"status": "ready"}


# ---------------------------------------------------------------- auth


@app.post("/api/v1/auth/register", status_code=201)
async def register(body: RegisterRequest) -> dict:
    if settings.auth_mode != "session":
        raise HTTPException(status_code=400, detail="auth disabled (AUTH_MODE=none)")
    if not settings.allow_registrations:
        raise HTTPException(status_code=403, detail="registrations disabled")
    users = UserRepository(get_db())
    if await users.get_by_username(body.username):
        raise HTTPException(status_code=409, detail="username taken")
    user = User(username=body.username, password_hash=hash_password(body.password))
    await users.create(user)
    return {"token": issue_token(user), "user_id": user.id}


@app.post("/api/v1/auth/login")
async def login(body: LoginRequest) -> dict:
    if settings.auth_mode != "session":
        raise HTTPException(status_code=400, detail="auth disabled (AUTH_MODE=none)")
    users = UserRepository(get_db())
    user = await users.get_by_username(body.username)
    if user is None or not verify_password(body.password, user.password_hash or ""):
        raise HTTPException(status_code=401, detail="invalid credentials")
    return {"token": issue_token(user), "user_id": user.id}


# ---------------------------------------------------------------- research


@app.post("/api/v1/research", status_code=202)
async def start_research(
    body: StartResearchRequest, identity: dict = Depends(current_identity)
) -> dict:
    try:
        pipeline_registry.get(body.pipeline_id)
    except KeyError as exc:
        raise HTTPException(status_code=422, detail=str(exc)) from exc
    config = RunConfig.model_validate(
        {**(body.config or {}), "pipeline_id": body.pipeline_id}
    )
    run = ResearchRun(
        question=body.question,
        pipeline_id=body.pipeline_id,
        workspace_id=identity["workspace_id"],
        config=config,
    )
    await RunRepositorySQL(get_db()).create(run)
    await get_queue().enqueue(run.id, {"pipeline_id": run.pipeline_id})
    return {"run_id": run.id, "status": run.status.value}


@app.get("/api/v1/research")
async def list_research(identity: dict = Depends(current_identity)) -> dict:
    runs = await RunRepositorySQL(get_db()).list_runs(
        workspace_id=identity["workspace_id"]
    )
    return {
        "runs": [
            {
                "id": r.id,
                "question": r.question,
                "pipeline_id": r.pipeline_id,
                "status": r.status.value,
                "created_at": r.created_at.isoformat(),
                "finished_at": r.finished_at.isoformat() if r.finished_at else None,
            }
            for r in runs
        ]
    }


async def _get_run_checked(run_id: str, identity: dict) -> ResearchRun:
    run = await RunRepositorySQL(get_db()).get(run_id)
    if run is None:
        raise HTTPException(status_code=404, detail="run not found")
    if settings.auth_mode == "session" and run.workspace_id != identity["workspace_id"]:
        raise HTTPException(status_code=404, detail="run not found")
    return run


@app.get("/api/v1/research/{run_id}")
async def get_research(
    run_id: str, identity: dict = Depends(current_identity)
) -> dict:
    run = await _get_run_checked(run_id, identity)
    return {
        "id": run.id,
        "question": run.question,
        "brief": run.brief,
        "pipeline_id": run.pipeline_id,
        "status": run.status.value,
        "error": run.error,
        "config": run.config.model_dump(mode="json"),
        "created_at": run.created_at.isoformat(),
        "started_at": run.started_at.isoformat() if run.started_at else None,
        "finished_at": run.finished_at.isoformat() if run.finished_at else None,
    }


@app.post("/api/v1/research/{run_id}/cancel")
async def cancel_research(
    run_id: str, identity: dict = Depends(current_identity)
) -> dict:
    run = await _get_run_checked(run_id, identity)
    if run.status in (RunStatus.COMPLETED, RunStatus.FAILED, RunStatus.CANCELLED):
        raise HTTPException(status_code=409, detail=f"run already {run.status.value}")
    await get_queue().request_cancel(run_id)
    return {"run_id": run_id, "cancel_requested": True}


@app.post("/api/v1/research/{run_id}/steer")
async def steer_research(
    run_id: str, body: SteerRequest, identity: dict = Depends(current_identity)
) -> dict:
    await _get_run_checked(run_id, identity)
    await get_queue().push_steering(run_id, body.message)
    return {"run_id": run_id, "steered": True}


@app.get("/api/v1/research/{run_id}/report")
async def get_report(
    run_id: str, identity: dict = Depends(current_identity)
) -> dict:
    run = await _get_run_checked(run_id, identity)
    artifacts = await ArtifactRepository(get_db()).list_for_run(run_id)
    report = next(
        (a for a in artifacts if a.kind == ArtifactKind.REPORT_MARKDOWN), None
    )
    if report is None:
        raise HTTPException(status_code=404, detail="report not ready")
    citations = await CitationRepository(get_db()).list_for_run(run_id)
    return {
        "run_id": run_id,
        "status": run.status.value,
        "report_markdown": report.content,
        "citations": [c.model_dump(mode="json") for c in citations],
        "artifacts": [
            {"id": a.id, "kind": a.kind.value} for a in artifacts
        ],
    }


@app.get("/api/v1/research/{run_id}/knowledge-map")
async def get_knowledge_map(
    run_id: str, identity: dict = Depends(current_identity)
) -> dict:
    await _get_run_checked(run_id, identity)
    nodes, edges = await KnowledgeRepository(get_db()).load_map(run_id)
    return {
        "nodes": [n.model_dump(mode="json") for n in nodes],
        "edges": [e.model_dump(mode="json") for e in edges],
    }


@app.get("/api/v1/research/{run_id}/events")
async def get_events(
    run_id: str, identity: dict = Depends(current_identity)
) -> dict:
    await _get_run_checked(run_id, identity)
    events = await EventRepository(get_db()).list_events(run_id)
    return {"events": [e.to_wire() for e in events]}


# ---------------------------------------------------------------- catalog


@app.get("/api/v1/pipelines")
async def list_pipelines() -> dict:
    return {
        "pipelines": [
            {
                "id": s.id,
                "name": s.name,
                "description": s.description,
                "tags": s.tags,
            }
            for s in pipeline_registry.list_specs()
        ]
    }


@app.get("/api/v1/providers")
async def list_providers() -> dict:
    return {
        "llm_providers": llm_registry.providers(),
        "search_engines": search_engine_registry.engines(),
        "search_strategies": strategy_registry.strategies(),
    }


# ---------------------------------------------------------------- websocket


@app.websocket("/api/v1/research/{run_id}/events/ws")
async def events_ws(websocket: WebSocket, run_id: str) -> None:
    """Replay persisted events, then stream live ones from Redis pub/sub."""
    await websocket.accept()
    events = EventRepository(app.state.db)
    try:
        replayed = await events.list_events(run_id)
        for event in replayed:
            await websocket.send_json(event.to_wire())
        if any(e.type.value in ("done", "error") for e in replayed):
            await websocket.close()
            return
        pubsub = app.state.redis.pubsub()
        await pubsub.subscribe(events_channel(run_id))
        try:
            while True:
                message = await pubsub.get_message(
                    ignore_subscribe_messages=True, timeout=1.0
                )
                if message is None:
                    await asyncio.sleep(0)
                    continue
                data = message["data"]
                if isinstance(data, bytes):
                    data = data.decode()
                payload = json.loads(data)
                await websocket.send_json(payload)
                if payload.get("type") in ("done", "error"):
                    break
        finally:
            await pubsub.unsubscribe(events_channel(run_id))
            await pubsub.aclose()
    except WebSocketDisconnect:
        pass
