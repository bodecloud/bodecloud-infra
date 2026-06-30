# Reading Paths and Retrieval Routes

This page exists because the knowledgebase is no longer mainly failing at
lacking content.

Its remaining failure mode is structural:

- the right pages exist
- the evidence boundaries are sharper
- the navigation can still tempt readers into browsing by topic instead of by
  the real question

That is dangerous in this repo because topic-first reading often recreates the
same old downgrade:

- "better HA"
- "better clustering"
- "better orchestrator choices"

instead of the harder repo-shaped question:

> how do several ordinary Docker nodes stop behaving like separate islands
> whose real correctness still depends on private operator reconstruction when
> the request lands on the wrong machine?

This page is the answer to that structural problem.

## What this page is and is not allowed to prove

This page is allowed to:

- route readers to the right evidence class for the question they actually have
- stop the site from being read like a generic topic taxonomy
- reduce the chance that the user’s real dream gets replaced by a smaller
  neighboring question

This page is not allowed to:

- substitute for the proof pages
- imply that navigation clarity means implementation maturity
- make a route look sufficient if it only answers the easier subproblem

## Read by the question you are really asking

### 1. "What is the user actually trying to make true?"

Read in this order:

1. [User Intent and Dream](research/user-intent-and-dream.md)
2. [Operator Contract and Success Criteria](architecture/operator-contract.md)
3. [Operator Questions and Honest Answers](operations/operator-questions-and-honest-answers.md)

Use this path when you need:

- the real dream
- the negative benchmark
- the anti-fake-options framing

Do not use this path as runtime proof.

### 2. "What does the priority runtime actually prove today?"

Read in this order:

1. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
2. [Current Compose Runtime](architecture/current-compose-runtime.md)
3. [Compose Fragment Map](architecture/compose-fragment-map.md)
4. [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

Use this path when you need:

- live runtime truth
- authoritative source ordering
- the line between config shape and stronger proof

### 3. "Why is wrong-node traffic still the humiliating threshold?"

Read in this order:

1. [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
2. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
3. [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)

Use this path when you need:

- request-path realism
- wrong-node meaning preservation
- where first-hop plurality stops being enough

### 4. "What still remains hidden in operator memory?"

Read in this order:

1. [The Missing Middle Layer](architecture/missing-middle-layer.md)
2. [Failure Model and Maturity](architecture/failure-model-and-maturity.md)
3. [Decision Paths and Promotion Rules](operations/decision-paths-and-promotion-rules.md)
4. [Source Assimilation Index](operations/source-assimilation-index.md)

Use this path when you need:

- hidden SPOFs
- burden ownership analysis
- evidence discipline

### 5. "Why are stateful services a separate honesty problem?"

Read in this order:

1. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
2. [Stateful HA Evidence](research/stateful-ha-evidence.md)
3. [Stateful HA Plan](research/stateful-ha-plan.md)

Use this path when you need:

- the difference between reachability and authority
- the live stateful risk surface
- planned stateful futures without overclaiming them

### 6. "Which future paths are real options versus renamed burden?"

Read in this order:

1. [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)
2. [Orchestration Research 2026](research/orchestration-research-2026.md)
3. [Infrastructure Master Plan](research/infrastructure-master-plan.md)
4. [OpenSVC Cluster Bootstrap](research/opensvc-cluster-bootstrap.md)
5. [Garden and k3s Exploration](research/garden-k3s-exploration-evidence.md)
6. [Nomad Exploration](research/nomad-exploration-evidence.md)

Use this path when you need:

- candidate-future comparison by burden ownership
- not just product taxonomy
- side-path evidence without premature promotion

## Read by truth register when the claim type is already known

### Live runtime truth

Start with:

- [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
- [Current Compose Runtime](architecture/current-compose-runtime.md)
- [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

### Dream and intent truth

Start with:

- [User Intent and Dream](research/user-intent-and-dream.md)
- [Operator Contract and Success Criteria](architecture/operator-contract.md)
- [Problem and Goals](architecture/problem-and-goals.md)

### Planning and promotion truth

Start with:

- [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md)
- [Infrastructure Master Plan](research/infrastructure-master-plan.md)
- [Decision Paths and Promotion Rules](operations/decision-paths-and-promotion-rules.md)

### Archive pressure and reconstruction truth

Start with:

- [Archive Pressure Patterns](research/archive-pressure-patterns.md)
- [Source Assimilation Index](operations/source-assimilation-index.md)
- [Evidence Ledger](research/evidence-ledger.md)

## Fastest route for an impatient but serious reader

If someone only has time for one high-fidelity route, use:

1. [User Intent and Dream](research/user-intent-and-dream.md)
2. [Operator Contract and Success Criteria](architecture/operator-contract.md)
3. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
4. [Current Compose Runtime](architecture/current-compose-runtime.md)
5. [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
6. [The Missing Middle Layer](architecture/missing-middle-layer.md)
7. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
8. [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

That route does the minimum necessary to preserve:

- the dream
- the acceptance bar
- the live proof boundary
- the wrong-node humiliation threshold
- the hidden burden analysis
- the stateful honesty split

## Strongest honest current answer

The strongest honest current answer is that this site should not be navigated
like a normal infrastructure docs tree.

It should be navigated like a retrieval system that is trying to stop the same
old substitution error:

- answering the smaller neighboring question
- with the wrong evidence class
- and then mistaking cleaner navigation for stronger truth

If this page helps a reader reach the right pages but still answer the calmer
question, it has still failed.
