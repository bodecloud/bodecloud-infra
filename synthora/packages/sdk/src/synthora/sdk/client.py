"""Python client for the Synthora API (R-LDR-8)."""

from __future__ import annotations

import time
from typing import Any, Optional

import httpx


class SynthoraClient:
    """Synchronous client mirroring the REST API.

    Example:
        client = SynthoraClient("http://localhost:8000")
        run_id = client.start_research("What is X?", pipeline_id="fast_research")
        report = client.wait_for_report(run_id)
    """

    def __init__(
        self,
        base_url: str = "http://localhost:8000",
        *,
        token: Optional[str] = None,
        timeout: float = 30.0,
    ) -> None:
        self.base_url = base_url.rstrip("/")
        self.token = token
        self._client = httpx.Client(base_url=self.base_url, timeout=timeout)

    # -- auth ----------------------------------------------------------------

    def register(self, username: str, password: str) -> str:
        data = self._post("/api/v1/auth/register", {"username": username, "password": password})
        self.token = data["token"]
        return self.token

    def login(self, username: str, password: str) -> str:
        data = self._post("/api/v1/auth/login", {"username": username, "password": password})
        self.token = data["token"]
        return self.token

    # -- research ------------------------------------------------------------

    def start_research(
        self,
        question: str,
        *,
        pipeline_id: str = "deep_research",
        config: Optional[dict[str, Any]] = None,
    ) -> str:
        data = self._post(
            "/api/v1/research",
            {"question": question, "pipeline_id": pipeline_id, "config": config},
        )
        return data["run_id"]

    def get_run(self, run_id: str) -> dict:
        return self._get(f"/api/v1/research/{run_id}")

    def list_runs(self) -> list[dict]:
        return self._get("/api/v1/research")["runs"]

    def cancel(self, run_id: str) -> dict:
        return self._post(f"/api/v1/research/{run_id}/cancel", {})

    def steer(self, run_id: str, message: str) -> dict:
        return self._post(f"/api/v1/research/{run_id}/steer", {"message": message})

    def get_report(self, run_id: str) -> dict:
        return self._get(f"/api/v1/research/{run_id}/report")

    def get_events(self, run_id: str) -> list[dict]:
        return self._get(f"/api/v1/research/{run_id}/events")["events"]

    def get_knowledge_map(self, run_id: str) -> dict:
        return self._get(f"/api/v1/research/{run_id}/knowledge-map")

    def list_pipelines(self) -> list[dict]:
        return self._get("/api/v1/pipelines")["pipelines"]

    def list_providers(self) -> dict:
        return self._get("/api/v1/providers")

    def wait_for_report(
        self, run_id: str, *, poll_seconds: float = 2.0, timeout: float = 1800.0
    ) -> dict:
        deadline = time.monotonic() + timeout
        while time.monotonic() < deadline:
            run = self.get_run(run_id)
            if run["status"] == "completed":
                return self.get_report(run_id)
            if run["status"] in ("failed", "cancelled"):
                raise RuntimeError(
                    f"run {run_id} {run['status']}: {run.get('error')}"
                )
            time.sleep(poll_seconds)
        raise TimeoutError(f"run {run_id} did not finish within {timeout}s")

    # -- plumbing ----------------------------------------------------------

    def _headers(self) -> dict:
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    def _get(self, path: str) -> dict:
        resp = self._client.get(path, headers=self._headers())
        resp.raise_for_status()
        return resp.json()

    def _post(self, path: str, body: dict) -> dict:
        resp = self._client.post(path, json=body, headers=self._headers())
        resp.raise_for_status()
        return resp.json()

    def close(self) -> None:
        self._client.close()
