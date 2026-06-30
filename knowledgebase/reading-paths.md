# Reading Paths and Retrieval Routes

This page is the practical entry map for the knowledgebase.

The site is now large enough that browsing by folder can still waste time.
The useful question is not:

> which page sounds nearby?

It is:

> which truth register answers the claim I actually need to make without
> flattening the repo into a calmer, smaller problem than the user is really
> forcing?

That difference matters.
This repo is especially vulnerable to documentation that feels organized while
quietly downgrading the actual question.

## What this page is and is not allowed to prove

This page is authoritative about:

- where a reader should start for a given question
- which pages answer runtime, dream, planning, archive-pressure, or proof
  questions best
- how to avoid answering a weaker neighboring question by accident
- how to keep retrieval disciplined instead of merely broad

This page is not authoritative about:

- proving failover behavior itself
- deciding the winning future control layer
- implying that strong navigation equals strong implementation maturity
- turning a reading route into a completion claim

## Strongest honest current answer

The knowledgebase is now broad enough that the main risk is no longer missing
pages.
The main risk is reading the right pages in the wrong order and accidentally
answering:

- "better HA"
- "better orchestration"
- "better docs"
- "what are the available options?"

instead of:

> how do several ordinary Docker nodes stop behaving like separate islands,
> while `docker-compose.yml` stays legible and wrong-node entry stops
> depending on operator folklore?

## The reading mistake this page is trying to stop

Most bad retrieval in this repo follows the same pattern:

1. start with the folder name
2. gather a reasonable-looking set of nearby pages
3. blend dream, runtime, planning, and archive pressure into one clean voice
4. produce a conclusion that sounds mature
5. quietly lose the user's actual complaint

This page exists to stop that sequence.

It also exists to stop a more subtle error:

- the stronger the site gets, the easier it is to believe the repo itself is
  already closer to closure than the evidence supports

That is false.
The site is getting better at reconstructing the pressure.
That does not mean the runtime has already solved the pressure.

## Read by claim type first

Before following a reading path, identify the class of sentence you are trying
to support.

### If the sentence is really claiming...

| Claim class | Start with | Why | It still must not imply... |
| --- | --- | --- | --- |
| what the user actually wants | [User Intent and Dream](research/user-intent-and-dream.md) | strongest reconstruction of the dream and anti-goals | that the runtime already earned that dream |
| what the current implementation actually contains | [Current Compose Runtime](architecture/current-compose-runtime.md) | strongest live-runtime inventory | that presence equals resilience under failure |
| what the stack is still missing | [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md) | strongest sequencing of unresolved burdens | that the gaps are already partially closed just because they are clearly named |
| why normal answers still feel fake | [Operator Questions and Honest Answers](operations/operator-questions-and-honest-answers.md) | strongest user-facing explanation of the hidden burden | that the critique itself proves the replacement |
| what proof is still required | [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md) | strongest proof boundary page | that a missing drill can be replaced by elegant theory |
| what future layer is being tested | [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md) | strongest candidate-layer pressure page | that a candidate therefore earned default promotion |

If you cannot route a claim this way, the claim is probably too vague to be
useful.

## Question-first reading paths

### 1. "What is the user actually trying to build?"

Read in this order:

1. [User Intent and Dream](research/user-intent-and-dream.md)
2. [Problem and Goals](architecture/problem-and-goals.md)
3. [Operator Contract and Success Criteria](architecture/operator-contract.md)
4. [Archive Pressure Patterns](research/archive-pressure-patterns.md)

Use this path when you need:

- the dream
- the anti-goals
- the negative benchmark
- the real acceptance bar
- the reason ordinary option lists keep feeling insulting

Do **not** use it as runtime proof.

### 2. "What does the current root runtime actually contain?"

Read in this order:

1. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
2. [Current Compose Runtime](architecture/current-compose-runtime.md)
3. [Compose Fragment Map](architecture/compose-fragment-map.md)
4. [Failure Model and Maturity](architecture/failure-model-and-maturity.md)
5. [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)
6. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)

Use this path when you need:

- the priority implementation
- the real root include graph
- the difference between "present in YAML" and "system-owned under failure"

This path is for live truth, not wishful extrapolation.

### 3. "Why is wrong-node traffic still the real humiliation threshold?"

Read in this order:

1. [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
2. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
3. [Operator Questions and Honest Answers](operations/operator-questions-and-honest-answers.md)
4. [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)
5. `source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
6. `source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md`

Use this path when you need:

- literal request-path reasoning
- why first-hop plurality is weaker than preserved request meaning
- where peer-forward routing still lacks proof
- why the user's frustration is sharper than "needs better load balancing"

### 4. "What helper layer is the repo actually looking for?"

Read in this order:

1. [The Missing Middle Layer](architecture/missing-middle-layer.md)
2. [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md)
3. [Decision Paths and Promotion Rules](operations/decision-paths-and-promotion-rules.md)
4. [Source Assimilation Index](operations/source-assimilation-index.md)
5. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

Use this path when you need:

- the actual missing control surface
- why the repo is not satisfied with raw Compose or immediate orchestrator
  surrender
- the burden-ownership test future candidates must pass

### 5. "Why are stateful services a much harsher problem?"

Read in this order:

1. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
2. [Stateful HA Evidence](research/stateful-ha-evidence.md)
3. [Stateful HA Plan](research/stateful-ha-plan.md)
4. [Failure Model and Maturity](architecture/failure-model-and-maturity.md)

Use this path when you need:

- the difference between liveness and authority
- why Redis, MongoDB, Headscale, and databases cannot inherit HTTP optimism
- why stateful promotion must remain much slower and harsher

### 6. "Which future paths are real options and which are still theater?"

Read in this order:

1. [Orchestration Options](architecture/orchestration-options.md)
2. [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)
3. [Infrastructure Master Plan](research/infrastructure-master-plan.md)
4. [Garden and k3s Exploration Evidence](research/garden-k3s-exploration-evidence.md)
5. [Nomad Exploration Evidence](research/nomad-exploration-evidence.md)
6. [OpenSVC Cluster Bootstrap](research/opensvc-cluster-bootstrap.md)
7. [OpenSVC Ingress HA](research/osvc-ingress-ha.md)
8. `source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
9. `source-archive/chatgpt-exports/conversations/nomad-multi-node-failover__68765e45-1ec4-8006-9179-5ef176d7a90f.md`

Use this path when you need:

- candidate middle layers
- honest tradeoffs
- future directions without mistaking them for current truth

### 7. "How do I keep my own summary from lying?"

Read in this order:

1. [Source Assimilation Index](operations/source-assimilation-index.md)
2. [Evidence Ledger](research/evidence-ledger.md)
3. [Archive Pressure Patterns](research/archive-pressure-patterns.md)
4. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)

Use this path when you need:

- retrieval discipline
- source hierarchy
- archive pressure without runtime inflation
- a check against blending unlike truth classes into one neat narrative

### 8. "Which archive threads should I actually read instead of vaguely searching everything?"

Read in this order:

1. [Archive Pressure Patterns](research/archive-pressure-patterns.md)
2. [Source Assimilation Index](operations/source-assimilation-index.md)
3. then one of these focused bundles:

- Compose pain bundle:
  - `source-archive/chatgpt-exports/conversations/docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`
  - `source-archive/chatgpt-exports/conversations/docker-compose-multi-server-setup__67f73c50-150c-8006-8408-c03db2d8d287.md`
- failover realism bundle:
  - `source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
  - `source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md`
- orchestrator pressure bundle:
  - `source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
  - `source-archive/chatgpt-exports/conversations/nomad-multi-node-failover__68765e45-1ec4-8006-9179-5ef176d7a90f.md`

Use this path when you need:

- actual archive reconstruction instead of generic RAG theater
- the user complaint in its sharper native form
- smaller high-signal bundles instead of one giant undirected archive pass

## Read by truth register

### Live runtime truth

Start with:

- [Current Compose Runtime](architecture/current-compose-runtime.md)
- [Compose Fragment Map](architecture/compose-fragment-map.md)
- [Failure Model and Maturity](architecture/failure-model-and-maturity.md)
- [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

Use this register for:

- what exists now
- which helpers are materially live
- which lanes are still only runtime-shaped rather than trustworthy

### Intent and dream truth

Start with:

- [User Intent and Dream](research/user-intent-and-dream.md)
- [Problem and Goals](architecture/problem-and-goals.md)
- [Operator Contract and Success Criteria](architecture/operator-contract.md)

Use this register for:

- what the user actually wants
- what would count as genuine relief
- which nearby goals are still too small

### Planning and promotion truth

Start with:

- [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md)
- [Infrastructure Master Plan](research/infrastructure-master-plan.md)
- [Decision Paths and Promotion Rules](operations/decision-paths-and-promotion-rules.md)

Use this register for:

- named missing layers
- sequencing pressure
- when a stronger control surface might earn itself

### Archive and reconstruction truth

Start with:

- [Archive Pressure Patterns](research/archive-pressure-patterns.md)
- [Source Assimilation Index](operations/source-assimilation-index.md)
- [Operator Questions and Honest Answers](operations/operator-questions-and-honest-answers.md)

Use this register for:

- why normal answers keep failing
- which bundles are worth reading closely
- how to stop replacing the real complaint with a tidier one

## Bottom line

The right reading path in this repo is never just "start in the most relevant
folder."

It is:

> identify the class of claim, pick the smallest strong bundle that can answer
> it honestly, and keep runtime truth, intent truth, planning truth, and
> archive-pressure truth from quietly merging into one stronger story than the
> current worktree has actually earned.

If that discipline disappears, the site can still feel organized while losing
the real question again.
