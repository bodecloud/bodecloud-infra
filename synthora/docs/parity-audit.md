# Synthora parity audit checklist

Living checklist against Open Deep Research, STORM/Co-STORM, and Local Deep
Research. Status values: `done`, `in_progress`, `todo`. Updated as gap-closure
ships on `feat/synthora-full-parity`.

## Open Deep Research

| Capability | Status | Notes |
|---|---|---|
| Nested graphs | done | `graphs.py` / pipelines |
| Clarify interrupt | in_progress | needs checkpointer + resume API |
| Research brief | done | |
| Supervisor + parallel researchers | done | |
| Researcher ReAct + compress | done | |
| Final report | done | |
| Role-split models | done | |
| Token-limit retry/truncation | todo | |
| Page summarization before context | todo | |
| MCP tools into researchers | todo | |
| Native OpenAI/Anthropic web search | todo | |
| Studio / langgraph.json | done | |

## STORM / Co-STORM

| Capability | Status | Notes |
|---|---|---|
| Perspective discovery | done | |
| Perspective-guided questions | in_progress | wire into discourse |
| Discourse + moderator unknown-unknowns | done | persist turns |
| Knowledge map insert/reorganize | done | embeddings todo |
| Outline + section write | done | |
| Polish pass | in_progress | wire into pipelines |
| PureRAG / warm-start / simulated user | todo | |
| Wikipedia TOC mining | todo | |
| Embedding similarity | todo | Jaccard fallback today |

## Local Deep Research

| Capability | Status | Notes |
|---|---|---|
| Persistence / jobs / WS / auth API | done | sessions unused until Phase 1 |
| Strategies (5) | in_progress | 2 of 5 |
| Search engines (full set) | in_progress | 5 registered |
| LLM providers (full set) | in_progress | 3 registered |
| Document library + RAG | todo | |
| Settings persistence | todo | |
| Export md/html | done API | web UI todo |
| PDF export | todo | |
| Delete / clear history | todo | |
| MCP server (outbound) | todo | |
| News / subscriptions | todo | |
| Metrics / benchmarks | todo | |
| Chat / follow-up research | todo | |

## Explicit non-goals (product, not capability gaps)

- Vendoring upstream source trees as disconnected apps
- Per-user SQLCipher (shared Postgres + optional auth by design)
- STORM Streamlit demo (React UI is the product surface)
- Paper eval datasets on backup branches (FreshWiki construction scripts)
