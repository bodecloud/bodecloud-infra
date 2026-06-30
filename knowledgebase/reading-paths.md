# Reading Paths and Retrieval Routes

This page is the practical entry map for the knowledgebase.

The site is large enough now that topic browsing can still waste time. The
useful question is not "which folder sounds close?" It is "which truth register
answers the question I actually have?"

## What this page is and is not allowed to prove

This page is authoritative about:

- where a reader should start for a given question
- which pages answer runtime, intent, planning, or proof questions best
- how to avoid collapsing the repo into a calmer, smaller problem

This page is not authoritative about:

- proving failover behavior itself
- deciding the winning future control layer
- implying that good navigation means strong implementation maturity

## Strongest honest current answer

The knowledgebase is now broad enough that the main risk is no longer missing
pages. The main risk is reading the right pages in the wrong order and
accidentally answering:

- "better HA"
- "better orchestration"
- "better docs"

instead of:

> how do several ordinary Docker nodes stop behaving like separate islands,
> while `docker-compose.yml` stays legible and wrong-node entry stops depending
> on operator folklore?

## Question-first reading paths

### 1. "What is the user actually trying to build?"

Read in this order:

1. [User Intent and Dream](research/user-intent-and-dream.md)
2. [Problem and Goals](architecture/problem-and-goals.md)
3. [Operator Contract and Success Criteria](architecture/operator-contract.md)

Use this path when you need:

- the dream
- the anti-goals
- the success criteria

Do not use it as runtime proof.

### 2. "What does the current root runtime actually contain?"

Read in this order:

1. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
2. [Current Compose Runtime](architecture/current-compose-runtime.md)
3. [Compose Fragment Map](architecture/compose-fragment-map.md)
4. [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

Use this path when you need:

- the priority implementation
- the real root include graph
- the gap between "present in YAML" and "proven under failure"

### 3. "Why is wrong-node traffic still the real threshold?"

Read in this order:

1. [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
2. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
3. [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)

Use this path when you need:

- literal request-path reasoning
- why multi-A DNS is weaker than request preservation
- where peer-forward routing still lacks proof

### 4. "What helper layer is the repo actually looking for?"

Read in this order:

1. [The Missing Middle Layer](architecture/missing-middle-layer.md)
2. [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md)
3. [Decision Paths and Promotion Rules](operations/decision-paths-and-promotion-rules.md)
4. [Source Assimilation Index](operations/source-assimilation-index.md)

Use this path when you need:

- the actual missing control surface
- why the repo is not satisfied with raw Compose or immediate Kubernetes
- the burden-ownership test future candidates must pass

### 5. "Why are stateful services a stricter problem?"

Read in this order:

1. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
2. [Stateful HA Evidence](research/stateful-ha-evidence.md)
3. [Stateful HA Plan](research/stateful-ha-plan.md)

Use this path when you need:

- the difference between liveness and authority
- why Redis, MongoDB, Headscale, and databases cannot inherit HTTP optimism

### 6. "Which future paths are real options and which are still theater?"

Read in this order:

1. [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)
2. [Infrastructure Master Plan](research/infrastructure-master-plan.md)
3. [Garden and k3s Exploration Evidence](research/garden-k3s-exploration-evidence.md)
4. [Nomad Exploration Evidence](research/nomad-exploration-evidence.md)
5. [OpenSVC Cluster Bootstrap](research/opensvc-cluster-bootstrap.md)
6. [OpenSVC Ingress HA](research/osvc-ingress-ha.md)

Use this path when you need:

- candidate middle layers
- honest tradeoffs
- future directions without mistaking them for live truth

## Read by truth register

### Live runtime truth

Start with:

- [Current Compose Runtime](architecture/current-compose-runtime.md)
- [Compose Fragment Map](architecture/compose-fragment-map.md)
- [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

### Intent and dream truth

Start with:

- [User Intent and Dream](research/user-intent-and-dream.md)
- [Problem and Goals](architecture/problem-and-goals.md)
- [Operator Contract and Success Criteria](architecture/operator-contract.md)

### Planning and promotion truth

Start with:

- [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md)
- [Infrastructure Master Plan](research/infrastructure-master-plan.md)
- [Decision Paths and Promotion Rules](operations/decision-paths-and-promotion-rules.md)

### Archive and reconstruction truth

Start with:

- [Archive Pressure Patterns](research/archive-pressure-patterns.md)
- [Source Assimilation Index](operations/source-assimilation-index.md)
- [Evidence Ledger](research/evidence-ledger.md)

## Fastest route for an impatient serious reader

If someone only wants the shortest route to the real situation, use:

1. [User Intent and Dream](research/user-intent-and-dream.md)
2. [Problem and Goals](architecture/problem-and-goals.md)
3. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
4. [Current Compose Runtime](architecture/current-compose-runtime.md)
5. [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
6. [The Missing Middle Layer](architecture/missing-middle-layer.md)
7. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
8. [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

That is the minimum path that keeps:

- the dream
- the root implementation
- the wrong-node problem
- the missing middle
- the stateful split
- the proof boundary

visible at the same time.
