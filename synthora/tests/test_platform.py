"""U6: platform integration — API, queue, executor, auth, WebSocket, SDK paths."""

from __future__ import annotations

import json

import fakeredis.aioredis
import pytest
import synthora.api.main as api_main
from fastapi.testclient import TestClient
from synthora.adapters import llm_registry, search_engine_registry
from synthora.api.settings import settings
from synthora.worker.executor import RunExecutor

from tests.conftest import FakeSearchEngine

SEARCH = json.dumps({"action": "search", "query": "q"})


class RoutingFakeModel:
    """Stateless fake LLM routed by system prompt (works for every role)."""

    async def complete(self, messages, *, temperature=0.3, max_tokens=None) -> str:
        system = messages[0]["content"]
        if "Rewrite the research request" in system:
            return "integration test brief"
        if "researcher with a search tool" in system:
            return SEARCH
        if "research supervisor" in system:
            return json.dumps({"action": "research_complete"})
        if "Compress these research findings" in system:
            return "compressed findings [1]"
        if "final research report" in system:
            return "# Integration Report\n\nA finding [1].\n\n## Sources"
        if "Decompose the research topic" in system:
            return "sub query"
        return "ok"


@pytest.fixture
def platform(tmp_path, monkeypatch):
    """API TestClient wired to a temp SQLite DB and fakeredis, with fake
    LLM/search providers registered."""
    settings.database_url = f"sqlite+aiosqlite:///{tmp_path}/synthora-test.db"
    settings.auth_mode = "none"

    fake_redis = fakeredis.aioredis.FakeRedis()
    monkeypatch.setattr(
        api_main.aioredis, "from_url", lambda url: fake_redis
    )

    llm_registry.register("fake", lambda m: RoutingFakeModel())
    search_engine_registry.register("fake", FakeSearchEngine)

    with TestClient(api_main.app) as client:
        yield client, api_main.app


def fake_run_config() -> dict:
    return {
        "planner_model": "fake:m",
        "researcher_model": "fake:m",
        "compressor_model": "fake:m",
        "writer_model": "fake:m",
        "critic_model": "fake:m",
        "search_engines": ["fake"],
        "search_strategy": "source_based",
        "max_react_tool_calls": 1,
    }


def make_executor(app) -> RunExecutor:
    return RunExecutor(app.state.db, app.state.queue)


# ------------------------------------------------------------------ health/catalog


def test_health_and_ready(platform):
    client, _ = platform
    assert client.get("/health").json() == {"status": "ok"}
    assert client.get("/ready").json() == {"status": "ready"}


def test_catalog_endpoints(platform):
    client, _ = platform
    pipelines = client.get("/api/v1/pipelines").json()["pipelines"]
    assert {p["id"] for p in pipelines} == {
        "fast_research",
        "deep_research",
        "academic_research",
        "autonomous_research",
    }
    providers = client.get("/api/v1/providers").json()
    assert "openai" in providers["llm_providers"]
    assert "searxng" in providers["search_engines"]
    assert "source_based" in providers["search_strategies"]


# ------------------------------------------------------------------ lifecycle


def test_start_progress_complete_lifecycle(platform):
    client, app = platform
    resp = client.post(
        "/api/v1/research",
        json={
            "question": "What is integration testing?",
            "pipeline_id": "fast_research",
            "config": fake_run_config(),
        },
    )
    assert resp.status_code == 202
    run_id = resp.json()["run_id"]
    assert client.get(f"/api/v1/research/{run_id}").json()["status"] == "queued"

    # worker turn: dequeue and execute
    executor = make_executor(app)

    async def drive():
        job = await app.state.queue.dequeue(timeout=1)
        assert job["run_id"] == run_id
        return await executor.execute(run_id)

    run = client.portal.call(drive)
    assert run.status.value == "completed"

    detail = client.get(f"/api/v1/research/{run_id}").json()
    assert detail["status"] == "completed"
    assert detail["brief"] == "integration test brief"

    report = client.get(f"/api/v1/research/{run_id}/report").json()
    assert report["report_markdown"].startswith("# Integration Report")
    assert report["citations"]

    events = client.get(f"/api/v1/research/{run_id}/events").json()["events"]
    types = [e["type"] for e in events]
    assert "status" in types and "done" in types

    runs = client.get("/api/v1/research").json()["runs"]
    assert runs[0]["id"] == run_id


def test_unknown_pipeline_rejected(platform):
    client, _ = platform
    resp = client.post(
        "/api/v1/research", json={"question": "q?", "pipeline_id": "bogus"}
    )
    assert resp.status_code == 422


def test_cancel_before_execution(platform):
    client, app = platform
    run_id = client.post(
        "/api/v1/research",
        json={
            "question": "cancel me",
            "pipeline_id": "fast_research",
            "config": fake_run_config(),
        },
    ).json()["run_id"]
    assert (
        client.post(f"/api/v1/research/{run_id}/cancel").json()["cancel_requested"]
        is True
    )
    executor = make_executor(app)
    run = client.portal.call(executor.execute, run_id)
    assert run.status.value == "cancelled"
    # double-cancel of a finished run conflicts
    assert client.post(f"/api/v1/research/{run_id}/cancel").status_code == 409


def test_steer_lands_in_queue(platform):
    client, app = platform
    run_id = client.post(
        "/api/v1/research",
        json={
            "question": "steer me",
            "pipeline_id": "fast_research",
            "config": fake_run_config(),
        },
    ).json()["run_id"]
    client.post(f"/api/v1/research/{run_id}/steer", json={"message": "focus on cost"})

    async def drain():
        return await app.state.queue.drain_steering(run_id)

    assert client.portal.call(drain) == ["focus on cost"]


def test_websocket_replays_events(platform):
    client, app = platform
    run_id = client.post(
        "/api/v1/research",
        json={
            "question": "ws test",
            "pipeline_id": "fast_research",
            "config": fake_run_config(),
        },
    ).json()["run_id"]
    executor = make_executor(app)
    client.portal.call(executor.execute, run_id)

    with client.websocket_connect(f"/api/v1/research/{run_id}/events/ws") as ws:
        received = []
        while True:
            try:
                received.append(ws.receive_json())
            except Exception:
                break
        types = [e["type"] for e in received]
        assert types[0] == "status"
        assert types[-1] == "done"


def test_failed_run_reports_error(platform):
    client, app = platform

    class ExplodingModel:
        async def complete(self, messages, *, temperature=0.3, max_tokens=None):
            raise RuntimeError("provider exploded")

    llm_registry.register("exploding", lambda m: ExplodingModel())
    config = {**fake_run_config(), "planner_model": "exploding:m"}
    run_id = client.post(
        "/api/v1/research",
        json={"question": "will fail", "pipeline_id": "fast_research", "config": config},
    ).json()["run_id"]
    executor = make_executor(app)
    run = client.portal.call(executor.execute, run_id)
    assert run.status.value == "failed"
    assert "provider exploded" in run.error
    detail = client.get(f"/api/v1/research/{run_id}").json()
    assert detail["status"] == "failed"


# ------------------------------------------------------------------ auth


def test_session_auth_flow(platform):
    client, _ = platform
    settings.auth_mode = "session"
    try:
        # unauthenticated request rejected
        assert client.get("/api/v1/research").status_code == 401

        resp = client.post(
            "/api/v1/auth/register",
            json={"username": "alice", "password": "supersecret"},
        )
        assert resp.status_code == 201
        token = resp.json()["token"]

        # duplicate username rejected
        assert (
            client.post(
                "/api/v1/auth/register",
                json={"username": "alice", "password": "supersecret"},
            ).status_code
            == 409
        )
        # wrong password rejected
        assert (
            client.post(
                "/api/v1/auth/login",
                json={"username": "alice", "password": "wrongpassword"},
            ).status_code
            == 401
        )
        # login works
        login = client.post(
            "/api/v1/auth/login",
            json={"username": "alice", "password": "supersecret"},
        )
        assert login.status_code == 200

        headers = {"Authorization": f"Bearer {token}"}
        assert client.get("/api/v1/research", headers=headers).status_code == 200
    finally:
        settings.auth_mode = "none"


def test_auth_endpoints_disabled_in_none_mode(platform):
    client, _ = platform
    resp = client.post(
        "/api/v1/auth/register", json={"username": "bob", "password": "password123"}
    )
    assert resp.status_code == 400
