# Feature parity matrix

Synthora is a clean-room synthesis of three projects (all MIT-licensed);
capabilities are reimplemented from their published architectures, not
vendored. Sources:
[langchain-ai/open_deep_research](https://github.com/langchain-ai/open_deep_research),
[stanford-oval/storm](https://github.com/stanford-oval/storm),
[LearningCircuit/local-deep-research](https://github.com/LearningCircuit/local-deep-research).

Status: ✅ implemented end-to-end · 🔶 partial / in progress · ⬜ not yet shipped

See also [parity-audit.md](parity-audit.md) for the living checklist.

## Open Deep Research (orchestration)

| Capability | Status | Synthora module |
|---|---|---|
| Nested graphs (top / supervisor / researcher) | ✅ | `orchestration/graphs.py` |
| Clarify-with-user interrupt | 🔶 | interrupt present; checkpointer + resume shipping |
| Supervisor tools: ConductResearch / think / ResearchComplete | ✅ | `supervisor` + `supervisor_route` |
| Parallel researchers + concurrency cap | ✅ | `supervisor_tools` |
| Isolated researcher ReAct loop | ✅ | `researcher_step` |
| Compress-before-return context isolation | ✅ | `compress_research` |
| Final report from brief + notes | ✅ | `final_report_generation` |
| Role-split models (5 roles) | ✅ | `RunConfig` + `ResearchContext` |
| Runtime config: env / configurable / API payload | ✅ | `studio.py`, `RunConfig`, API |
| `langgraph.json` + Studio surface | ✅ | `langgraph.json` |
| Token-limit retry / truncation | 🔶 | shipping |
| Page content summarization | 🔶 | shipping |
| MCP tool loading into researchers | 🔶 | shipping |
| Anthropic/OpenAI native web search tools | 🔶 | shipping |

## STORM / Co-STORM (intelligence)

| Capability | Status | Synthora module |
|---|---|---|
| Multi-perspective persona discovery | ✅ | `intelligence/perspectives.py` |
| Perspective-guided question asking | 🔶 | method exists; wiring into discourse shipping |
| Iterative grounded QA | ✅ | researcher + strategies + discourse |
| Collaborative discourse + turn policy | ✅ | `DiscourseManager` |
| Moderator unknown-unknowns | ✅ | `rank_unused_evidence` |
| Hierarchical mind map insert + reorganize | ✅ | `KnowledgeMap` |
| Outline-first then section-wise cited writing | ✅ | `OutlineBuilder` / `SectionWriter` |
| Polish pass (dedup + lead summary) | 🔶 | method exists; pipeline wiring shipping |
| Human steering mid-run | ✅ | steer API → discourse user turns |
| Discourse turn persistence | 🔶 | shipping |
| Embedding-based similarity | 🔶 | Jaccard default; embeddings shipping |
| Wikipedia-TOC perspective mining | 🔶 | shipping |
| PureRAG / warm-start / simulated user | 🔶 | shipping |

## Local Deep Research (platform)

| Capability | Status | Synthora module |
|---|---|---|
| Persistence: runs, artifacts, citations, knowledge maps | ✅ | `packages/persistence` |
| Sessions linked to runs | 🔶 | repo exists; API/UI shipping |
| Background jobs + lifecycle | ✅ | Redis + worker |
| REST API + cancel / report / events | ✅ | `apps/api` |
| Resume after clarify interrupt | 🔶 | shipping |
| Real-time progress (WebSocket) | ✅ | Redis pub/sub |
| User management + optional auth | ✅ | JWT + PBKDF2 |
| Auth / export / resume web UI | 🔶 | shipping |
| Search strategy abstraction (5 strategies) | 🔶 | 2 registered; rest shipping |
| Search engine abstraction (full catalog) | 🔶 | 5 registered; catalog shipping |
| LLM provider abstraction + think-tag handling | 🔶 | 3 registered; catalog shipping |
| Research history + export md/html | ✅ | API; web export button shipping |
| PDF export binary | 🔶 | shipping |
| Docker Compose self-host | ✅ | `docker-compose.yml` |
| Python SDK | ✅ | `packages/sdk` |
| Document library + RAG | 🔶 | table reserved; shipping |
| Provider settings persistence | 🔶 | table reserved; shipping |
| Delete run / clear history | 🔶 | shipping |
| MCP server exposing Synthora | 🔶 | shipping |
| News / subscriptions | 🔶 | shipping |
| Metrics / benchmarks | 🔶 | shipping |
| Follow-up / chat research | 🔶 | shipping |
| Per-user SQLCipher encrypted DBs | ⬜ | non-goal (shared Postgres) |

## Multi-pipeline requirement

| Pipeline | Status |
|---|---|
| `fast_research` | ✅ |
| `deep_research` | ✅ |
| `academic_research` | ✅ |
| `autonomous_research` | ✅ |
