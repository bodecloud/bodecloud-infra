"""SQLAlchemy table definitions (R-LDR-1)."""

from __future__ import annotations

from datetime import datetime, timezone

from sqlalchemy import (
    JSON,
    Boolean,
    DateTime,
    Float,
    ForeignKey,
    Integer,
    String,
    Text,
)
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


def utcnow() -> datetime:
    return datetime.now(timezone.utc)


class Base(DeclarativeBase):
    pass


class UserRow(Base):
    __tablename__ = "users"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    username: Mapped[str] = mapped_column(String(255), unique=True, index=True)
    password_hash: Mapped[str | None] = mapped_column(String(255), nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)


class WorkspaceRow(Base):
    __tablename__ = "workspaces"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    name: Mapped[str] = mapped_column(String(255), default="default")
    owner_id: Mapped[str | None] = mapped_column(
        String(32), ForeignKey("users.id"), nullable=True
    )
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)


class SessionRow(Base):
    __tablename__ = "sessions"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    workspace_id: Mapped[str] = mapped_column(String(32), ForeignKey("workspaces.id"))
    title: Mapped[str] = mapped_column(String(512), default="Untitled research")
    tags: Mapped[list] = mapped_column(JSON, default=list)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)


class ResearchRunRow(Base):
    __tablename__ = "research_runs"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    session_id: Mapped[str | None] = mapped_column(
        String(32), ForeignKey("sessions.id"), nullable=True
    )
    workspace_id: Mapped[str] = mapped_column(String(64), default="default", index=True)
    question: Mapped[str] = mapped_column(Text)
    brief: Mapped[str | None] = mapped_column(Text, nullable=True)
    pipeline_id: Mapped[str] = mapped_column(String(64), default="deep_research")
    status: Mapped[str] = mapped_column(String(32), default="queued", index=True)
    config: Mapped[dict] = mapped_column(JSON, default=dict)
    error: Mapped[str | None] = mapped_column(Text, nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)
    started_at: Mapped[datetime | None] = mapped_column(
        DateTime(timezone=True), nullable=True
    )
    finished_at: Mapped[datetime | None] = mapped_column(
        DateTime(timezone=True), nullable=True
    )


class RunEventRow(Base):
    __tablename__ = "run_events"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    run_id: Mapped[str] = mapped_column(
        String(32), ForeignKey("research_runs.id"), index=True
    )
    type: Mapped[str] = mapped_column(String(64))
    message: Mapped[str] = mapped_column(Text, default="")
    node: Mapped[str | None] = mapped_column(String(128), nullable=True)
    payload: Mapped[dict] = mapped_column(JSON, default=dict)
    timestamp: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)


class ArtifactRow(Base):
    __tablename__ = "artifacts"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    run_id: Mapped[str] = mapped_column(
        String(32), ForeignKey("research_runs.id"), index=True
    )
    kind: Mapped[str] = mapped_column(String(64))
    content: Mapped[str] = mapped_column(Text, default="")
    meta: Mapped[dict] = mapped_column(JSON, default=dict)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)


class CitationRow(Base):
    __tablename__ = "citations"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    run_id: Mapped[str | None] = mapped_column(
        String(32), ForeignKey("research_runs.id"), index=True, nullable=True
    )
    url: Mapped[str] = mapped_column(Text)
    title: Mapped[str] = mapped_column(Text, default="")
    snippet: Mapped[str] = mapped_column(Text, default="")
    confidence: Mapped[float] = mapped_column(Float, default=1.0)
    index: Mapped[int | None] = mapped_column(Integer, nullable=True)
    verified: Mapped[bool] = mapped_column(Boolean, default=False)


class KnowledgeNodeRow(Base):
    __tablename__ = "knowledge_nodes"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    run_id: Mapped[str | None] = mapped_column(
        String(32), ForeignKey("research_runs.id"), index=True, nullable=True
    )
    name: Mapped[str] = mapped_column(Text)
    summary: Mapped[str] = mapped_column(Text, default="")
    parent_id: Mapped[str | None] = mapped_column(String(32), nullable=True)
    infos: Mapped[list] = mapped_column(JSON, default=list)


class KnowledgeEdgeRow(Base):
    __tablename__ = "knowledge_edges"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    run_id: Mapped[str | None] = mapped_column(
        String(32), ForeignKey("research_runs.id"), index=True, nullable=True
    )
    source_id: Mapped[str] = mapped_column(String(32))
    target_id: Mapped[str] = mapped_column(String(32))
    relation: Mapped[str] = mapped_column(String(128), default="related_to")


class DiscourseTurnRow(Base):
    __tablename__ = "discourse_turns"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    run_id: Mapped[str | None] = mapped_column(
        String(32), ForeignKey("research_runs.id"), index=True, nullable=True
    )
    speaker: Mapped[str] = mapped_column(String(255))
    role: Mapped[str] = mapped_column(String(32), default="expert")
    utterance: Mapped[str] = mapped_column(Text)
    intent: Mapped[str] = mapped_column(String(32), default="answer")
    citations: Mapped[list] = mapped_column(JSON, default=list)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)


class DocumentRow(Base):
    __tablename__ = "documents"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    workspace_id: Mapped[str] = mapped_column(String(64), default="default", index=True)
    title: Mapped[str] = mapped_column(Text, default="")
    url: Mapped[str | None] = mapped_column(Text, nullable=True)
    path: Mapped[str | None] = mapped_column(Text, nullable=True)
    meta: Mapped[dict] = mapped_column(JSON, default=dict)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=utcnow)


class ProviderSettingRow(Base):
    __tablename__ = "provider_settings"

    id: Mapped[str] = mapped_column(String(32), primary_key=True)
    workspace_id: Mapped[str] = mapped_column(String(64), default="default", index=True)
    key: Mapped[str] = mapped_column(String(255))
    value: Mapped[dict] = mapped_column(JSON, default=dict)
