# Feature parity matrix

Synthora is a clean-room synthesis of three projects (all MIT-licensed);
capabilities were reimplemented from their published architectures, not
vendored. Sources:
[langchain-ai/open_deep_research](https://github.com/langchain-ai/open_deep_research),
[stanford-oval/storm](https://github.com/stanford-oval/storm),
[LearningCircuit/local-deep-research](https://github.com/LearningCircuit/local-deep-research).

Status: ✅ implemented · 🔶 partial · ⬜ deferred (explicit non-goal for v1)

## Open Deep Research (orchestration)

| Capability | Status | Synthora module |
|---|---|---|
| Nested graphs (top / supervisor / researcher) | ✅ | `orchestration/graphs.py` |
| Clarify-with-user interrupt → research brief | ✅ | `orchestration/nodes.py` (`clarify_with_user`, LangGraph `interrupt`) |
| Supervisor tools: ConductResearch / think / ResearchComplete | ✅ | `supervisor` + `supervisor_route` (JSON decisions, provider-agnostic) |
| Parallel researchers + concurrency cap + overflow errors | ✅ | `supervisor_tools` (`asyncio.gather`, capacity notes) |
| Isolated researcher ReAct loop | ✅ | `researcher_step` / `researcher_should_continue` |
| Compress-before-return context isolation | ✅ | `compress_research` |
| One-shot final report from brief + notes | ✅ | `final_report_generation` |
| Role-split models (5 roles) | ✅ | `RunConfig` + `ResearchContext` |
| Runtime config: env / configurable / API payload | ✅ | `studio.py`, `RunConfig`, API `config` |
| `langgraph.json` + Studio surface | ✅ | `langgraph.json`, `orchestration/studio.py` |
| MCP tool loading into researchers | ⬜ | planned: adapter in `packages/adapters` |
| Anthropic/OpenAI native web search tools | ⬜ | search engines cover retrieval |

## STORM / Co-STORM (intelligence)

| Capability | Status | Synthora module |
|---|---|---|
| Multi-perspective persona discovery | ✅ | `intelligence/perspectives.py` |
| Perspective-guided question asking | ✅ | `PerspectiveEngine.generate_questions` |
| Iterative grounded QA (decompose → retrieve → cite) | ✅ | researcher loop + strategies + discourse expert turns |
| Collaborative discourse: experts + turn policy | ✅ | `DiscourseManager` (L expert turns then moderator) |
| Moderator unknown-unknowns from unused evidence | ✅ | `rank_unused_evidence` (`sim^α · (1−sim)^(1−α)`) |
| Dynamic hierarchical mind map: insert + reorganize | ✅ | `KnowledgeMap` (similarity insert, LLM clustering at capacity K) |
| Outline-first then section-wise cited writing | ✅ | `OutlineBuilder` / `SectionWriter` |
| Polish pass (dedup + lead summary) | ✅ | `SectionWriter.polish` |
| Human steering mid-run | ✅ | steer API → `ctx.steering` → discourse user turns |
| Embedding-based similarity | 🔶 | lexical Jaccard default; `SimilarityFn` pluggable |
| Wikipedia-TOC perspective mining | 🔶 | personas mined from brief + retrieved context instead |

## Local Deep Research (platform)

| Capability | Status | Synthora module |
|---|---|---|
| Persistence: runs, sessions, artifacts, citations, knowledge maps | ✅ | `packages/persistence` (Postgres/SQLite, Alembic) |
| Background jobs + concurrency limits + lifecycle | ✅ | Redis queue + `apps/worker` (queued→running→completed/failed/cancelled) |
| REST API + status/cancel/report | ✅ | `apps/api/main.py` |
| Real-time progress (WebSocket) | ✅ | event replay + Redis pub/sub streaming |
| User management + optional auth | ✅ | `AUTH_MODE=none|session`, register/login, JWT, PBKDF2 |
| Search strategy abstraction | ✅ | `strategy_registry` (2 strategies; registry open) |
| Search engine abstraction (meta + academic) | ✅ | `search_engine_registry` (5 engines; registry open) |
| LLM provider abstraction + think-tag handling | ✅ | `llm_registry` + `strip_think_tags` |
| Research history + export | ✅ | history API/UI; export markdown + printable HTML |
| Docker Compose self-host (app + private search + local LLM) | ✅ | `docker-compose.yml` (api, worker, web, postgres, redis, searxng, ollama profile) |
| Python SDK | ✅ | `packages/sdk` (`SynthoraClient`) |
| PDF export binary | 🔶 | printable HTML export; browser print-to-PDF |
| Per-user SQLCipher encrypted DBs | ⬜ | non-goal v1 (shared Postgres chosen) |
| Local document library + RAG engine | ⬜ | `documents` table reserved |
| MCP server exposing Synthora as tools | ⬜ | planned |

## Multi-pipeline requirement

| Pipeline | Status |
|---|---|
| `fast_research` | ✅ |
| `deep_research` (ODR research + STORM synthesis + criticism) | ✅ |
| `academic_research` (lit search, citation verify, peer review, bibliography) | ✅ |
| `autonomous_research` (bounded hypothesize/investigate/gap loop) | ✅ |
