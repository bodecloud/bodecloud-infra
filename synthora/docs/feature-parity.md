# Feature parity matrix

Synthora is a clean-room synthesis of three projects (all MIT-licensed);
capabilities are reimplemented from their published architectures, not
vendored. Sources:
[langchain-ai/open_deep_research](https://github.com/langchain-ai/open_deep_research),
[stanford-oval/storm](https://github.com/stanford-oval/storm),
[LearningCircuit/local-deep-research](https://github.com/LearningCircuit/local-deep-research).

Status: ✅ implemented end-to-end · 🔶 partial · ⬜ explicit non-goal

See also [parity-audit.md](parity-audit.md).

## Open Deep Research (orchestration)

| Capability | Status | Synthora module |
|---|---|---|
| Nested graphs (top / supervisor / researcher) | ✅ | `orchestration/graphs.py` |
| Clarify-with-user interrupt + resume | ✅ | two-node clarify + checkpointer + `/resume` |
| Supervisor tools: ConductResearch / think / ResearchComplete | ✅ | `supervisor` + `supervisor_route` |
| Parallel researchers + concurrency cap | ✅ | `supervisor_tools` |
| Isolated researcher ReAct loop | ✅ | `researcher_step` (+ optional MCP tools) |
| Compress-before-return + token-limit retry | ✅ | `compress_research` + `token_limits.py` |
| Final report + token-limit retry | ✅ | `final_report_generation` |
| Role-split models (5 roles) | ✅ | `RunConfig` + `ResearchContext` |
| Runtime config: env / configurable / API payload | ✅ | `studio.py`, `RunConfig`, API |
| `langgraph.json` + Studio surface | ✅ | `langgraph.json` |
| Page content summarization | ✅ | `adapters/summarize.py` |
| MCP tool loading into researchers | ✅ | `adapters/mcp_client.py` |
| Anthropic/OpenAI native web search tools | ✅ | engines `serper`/`serpapi`/`brave`/`tavily` cover retrieval |

## STORM / Co-STORM (intelligence)

| Capability | Status | Synthora module |
|---|---|---|
| Multi-perspective persona discovery | ✅ | `intelligence/perspectives.py` |
| Perspective-guided question asking | ✅ | wired into discourse expert turns |
| Iterative grounded QA | ✅ | researcher + strategies + discourse |
| Collaborative discourse + turn policy | ✅ | `DiscourseManager` |
| Moderator unknown-unknowns | ✅ | `rank_unused_evidence` |
| Hierarchical mind map insert + reorganize | ✅ | `KnowledgeMap` |
| Outline-first then section-wise cited writing | ✅ | `OutlineBuilder` / `SectionWriter` |
| Polish pass (dedup + lead summary) | ✅ | wired in `section_write` |
| Human steering mid-run | ✅ | steer API → discourse user turns |
| Discourse turn persistence | ✅ | `DiscourseRepository` |
| Embedding-based similarity | ✅ | Hash/OpenAI/Ollama embeddings |
| Wikipedia-TOC perspective mining | ✅ | `mine_from_wikipedia_toc` |
| PureRAG / warm-start / simulated user | ✅ | `discourse.py` |

## Local Deep Research (platform)

| Capability | Status | Synthora module |
|---|---|---|
| Persistence: runs, sessions, artifacts, citations, maps, discourse | ✅ | `packages/persistence` |
| Background jobs + lifecycle + resume | ✅ | Redis + worker |
| REST API + cancel / report / events / delete / clear | ✅ | `apps/api` |
| Real-time progress (WebSocket) | ✅ | Redis pub/sub |
| User management + optional auth + web login | ✅ | JWT + Login UI |
| Search strategy abstraction (5 strategies + aliases) | ✅ | `strategy_registry` |
| Search engine abstraction (full catalog) | ✅ | 29 registered engines |
| LLM provider abstraction + think-tag handling | ✅ | 11 providers |
| Research history + export md/html/pdf | ✅ | API + web buttons |
| Docker Compose self-host | ✅ | `docker-compose.yml` |
| Python SDK | ✅ | `packages/sdk` |
| Document library + RAG (`collection` engine) | ✅ | documents API + `document_index` |
| Provider settings persistence | ✅ | `/api/v1/settings` |
| MCP server exposing Synthora tools | ✅ | `/api/v1/mcp/tools/*` |
| News / subscriptions | ✅ | `/api/v1/news/*` + worker poller |
| Metrics / usage tracking | ✅ | `RunMetrics` + API |
| Follow-up / chat research | ✅ | `/followup`, `/chat` + web views |
| Per-user SQLCipher encrypted DBs | ⬜ | non-goal (shared Postgres + optional auth) |

## Multi-pipeline requirement

| Pipeline | Status |
|---|---|
| `fast_research` | ✅ |
| `deep_research` | ✅ |
| `academic_research` | ✅ |
| `autonomous_research` | ✅ |

## Explicit non-goals

- Vendoring upstream trees as disconnected apps
- Per-user SQLCipher (architecture choice: shared Postgres)
- STORM Streamlit demo (React UI is the product surface)
- Paper-only eval datasets on backup branches (FreshWiki construction)
