"""Repository implementations over SQLAlchemy async sessions."""

from __future__ import annotations

from typing import Optional

from sqlalchemy import select

from synthora.core.events import ProgressEvent, RunEventType
from synthora.core.models import (
    Artifact,
    ArtifactKind,
    Citation,
    KnowledgeEdge,
    KnowledgeNode,
    ResearchRun,
    RunConfig,
    RunStatus,
    Session,
    User,
    Workspace,
)
from synthora.persistence.database import Database
from synthora.persistence.tables import (
    ArtifactRow,
    CitationRow,
    KnowledgeEdgeRow,
    KnowledgeNodeRow,
    ResearchRunRow,
    RunEventRow,
    SessionRow,
    UserRow,
    WorkspaceRow,
)


class UserRepository:
    def __init__(self, db: Database) -> None:
        self.db = db

    async def create(self, user: User) -> User:
        async with self.db.session() as s:
            s.add(
                UserRow(
                    id=user.id,
                    username=user.username,
                    password_hash=user.password_hash,
                    created_at=user.created_at,
                )
            )
        return user

    async def get_by_username(self, username: str) -> Optional[User]:
        async with self.db.session() as s:
            row = (
                await s.execute(select(UserRow).where(UserRow.username == username))
            ).scalar_one_or_none()
        if row is None:
            return None
        return User(
            id=row.id,
            username=row.username,
            password_hash=row.password_hash,
            created_at=row.created_at,
        )


class WorkspaceRepository:
    def __init__(self, db: Database) -> None:
        self.db = db

    async def ensure_default(self) -> Workspace:
        async with self.db.session() as s:
            row = (
                await s.execute(
                    select(WorkspaceRow).where(WorkspaceRow.name == "default")
                )
            ).scalar_one_or_none()
            if row is None:
                ws = Workspace(name="default")
                s.add(WorkspaceRow(id=ws.id, name=ws.name, created_at=ws.created_at))
                return ws
            return Workspace(
                id=row.id, name=row.name, owner_id=row.owner_id, created_at=row.created_at
            )


class SessionRepository:
    def __init__(self, db: Database) -> None:
        self.db = db

    async def create(self, session: Session) -> Session:
        async with self.db.session() as s:
            s.add(
                SessionRow(
                    id=session.id,
                    workspace_id=session.workspace_id,
                    title=session.title,
                    tags=session.tags,
                    created_at=session.created_at,
                )
            )
        return session

    async def list_sessions(self, workspace_id: str, limit: int = 100) -> list[Session]:
        async with self.db.session() as s:
            rows = (
                (
                    await s.execute(
                        select(SessionRow)
                        .where(SessionRow.workspace_id == workspace_id)
                        .order_by(SessionRow.created_at.desc())
                        .limit(limit)
                    )
                )
                .scalars()
                .all()
            )
        return [
            Session(
                id=r.id,
                workspace_id=r.workspace_id,
                title=r.title,
                tags=list(r.tags or []),
                created_at=r.created_at,
            )
            for r in rows
        ]


def _run_from_row(row: ResearchRunRow) -> ResearchRun:
    return ResearchRun(
        id=row.id,
        session_id=row.session_id,
        workspace_id=row.workspace_id,
        question=row.question,
        brief=row.brief,
        pipeline_id=row.pipeline_id,
        status=RunStatus(row.status),
        config=RunConfig.model_validate(row.config or {}),
        error=row.error,
        created_at=row.created_at,
        started_at=row.started_at,
        finished_at=row.finished_at,
    )


class RunRepositorySQL:
    """Implements the RunRepository port over SQL."""

    def __init__(self, db: Database) -> None:
        self.db = db

    async def create(self, run: ResearchRun) -> ResearchRun:
        async with self.db.session() as s:
            s.add(
                ResearchRunRow(
                    id=run.id,
                    session_id=run.session_id,
                    workspace_id=run.workspace_id,
                    question=run.question,
                    brief=run.brief,
                    pipeline_id=run.pipeline_id,
                    status=run.status.value,
                    config=run.config.model_dump(mode="json"),
                    error=run.error,
                    created_at=run.created_at,
                    started_at=run.started_at,
                    finished_at=run.finished_at,
                )
            )
        return run

    async def get(self, run_id: str) -> Optional[ResearchRun]:
        async with self.db.session() as s:
            row = await s.get(ResearchRunRow, run_id)
        return _run_from_row(row) if row else None

    async def update(self, run: ResearchRun) -> ResearchRun:
        async with self.db.session() as s:
            row = await s.get(ResearchRunRow, run.id)
            if row is None:
                raise KeyError(f"run {run.id} not found")
            row.status = run.status.value
            row.brief = run.brief
            row.error = run.error
            row.started_at = run.started_at
            row.finished_at = run.finished_at
            row.config = run.config.model_dump(mode="json")
        return run

    async def list_runs(
        self, *, workspace_id: Optional[str] = None, limit: int = 50
    ) -> list[ResearchRun]:
        async with self.db.session() as s:
            q = select(ResearchRunRow).order_by(ResearchRunRow.created_at.desc()).limit(limit)
            if workspace_id:
                q = q.where(ResearchRunRow.workspace_id == workspace_id)
            rows = (await s.execute(q)).scalars().all()
        return [_run_from_row(r) for r in rows]


class EventRepository:
    def __init__(self, db: Database) -> None:
        self.db = db

    async def append(self, event: ProgressEvent) -> None:
        async with self.db.session() as s:
            s.add(
                RunEventRow(
                    run_id=event.run_id,
                    type=event.type.value,
                    message=event.message,
                    node=event.node,
                    payload=event.payload,
                    timestamp=event.timestamp,
                )
            )

    async def list_events(self, run_id: str, limit: int = 500) -> list[ProgressEvent]:
        async with self.db.session() as s:
            rows = (
                (
                    await s.execute(
                        select(RunEventRow)
                        .where(RunEventRow.run_id == run_id)
                        .order_by(RunEventRow.id.asc())
                        .limit(limit)
                    )
                )
                .scalars()
                .all()
            )
        return [
            ProgressEvent(
                run_id=r.run_id,
                type=RunEventType(r.type),
                message=r.message,
                node=r.node,
                payload=r.payload or {},
                timestamp=r.timestamp,
            )
            for r in rows
        ]


class ArtifactRepository:
    def __init__(self, db: Database) -> None:
        self.db = db

    async def save(self, artifact: Artifact) -> Artifact:
        async with self.db.session() as s:
            s.add(
                ArtifactRow(
                    id=artifact.id,
                    run_id=artifact.run_id,
                    kind=artifact.kind.value,
                    content=artifact.content,
                    meta=artifact.metadata,
                    created_at=artifact.created_at,
                )
            )
        return artifact

    async def list_for_run(self, run_id: str) -> list[Artifact]:
        async with self.db.session() as s:
            rows = (
                (
                    await s.execute(
                        select(ArtifactRow).where(ArtifactRow.run_id == run_id)
                    )
                )
                .scalars()
                .all()
            )
        return [
            Artifact(
                id=r.id,
                run_id=r.run_id,
                kind=ArtifactKind(r.kind),
                content=r.content,
                metadata=r.meta or {},
                created_at=r.created_at,
            )
            for r in rows
        ]


class CitationRepository:
    def __init__(self, db: Database) -> None:
        self.db = db

    async def save_many(self, citations: list[Citation]) -> None:
        async with self.db.session() as s:
            for c in citations:
                s.add(
                    CitationRow(
                        id=c.id,
                        run_id=c.run_id,
                        url=c.url,
                        title=c.title,
                        snippet=c.snippet,
                        confidence=c.confidence,
                        index=c.index,
                        verified=c.verified,
                    )
                )

    async def list_for_run(self, run_id: str) -> list[Citation]:
        async with self.db.session() as s:
            rows = (
                (
                    await s.execute(
                        select(CitationRow).where(CitationRow.run_id == run_id)
                    )
                )
                .scalars()
                .all()
            )
        return [
            Citation(
                id=r.id,
                run_id=r.run_id,
                url=r.url,
                title=r.title,
                snippet=r.snippet,
                confidence=r.confidence,
                index=r.index,
                verified=r.verified,
            )
            for r in rows
        ]


class KnowledgeRepository:
    def __init__(self, db: Database) -> None:
        self.db = db

    async def save_map(
        self,
        run_id: str,
        nodes: list[KnowledgeNode],
        edges: list[KnowledgeEdge],
    ) -> None:
        async with self.db.session() as s:
            for n in nodes:
                s.add(
                    KnowledgeNodeRow(
                        id=n.id,
                        run_id=run_id,
                        name=n.name,
                        summary=n.summary,
                        parent_id=n.parent_id,
                        infos=[c.model_dump(mode="json") for c in n.infos],
                    )
                )
            for e in edges:
                s.add(
                    KnowledgeEdgeRow(
                        id=e.id,
                        run_id=run_id,
                        source_id=e.source_id,
                        target_id=e.target_id,
                        relation=e.relation,
                    )
                )

    async def load_map(
        self, run_id: str
    ) -> tuple[list[KnowledgeNode], list[KnowledgeEdge]]:
        async with self.db.session() as s:
            node_rows = (
                (
                    await s.execute(
                        select(KnowledgeNodeRow).where(
                            KnowledgeNodeRow.run_id == run_id
                        )
                    )
                )
                .scalars()
                .all()
            )
            edge_rows = (
                (
                    await s.execute(
                        select(KnowledgeEdgeRow).where(
                            KnowledgeEdgeRow.run_id == run_id
                        )
                    )
                )
                .scalars()
                .all()
            )
        nodes = [
            KnowledgeNode(
                id=r.id,
                run_id=r.run_id,
                name=r.name,
                summary=r.summary,
                parent_id=r.parent_id,
                infos=[Citation.model_validate(c) for c in (r.infos or [])],
            )
            for r in node_rows
        ]
        edges = [
            KnowledgeEdge(
                id=r.id,
                run_id=r.run_id,
                source_id=r.source_id,
                target_id=r.target_id,
                relation=r.relation,
            )
            for r in edge_rows
        ]
        return nodes, edges
