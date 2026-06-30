# Reading Paths and Retrieval Routes

This page is the practical map for reading the knowledgebase without
accidentally answering a smaller question than the user is actually asking.

The core retrieval mistake in this repo is simple:

- find relevant pages
- blend runtime, intent, plans, and archive pressure into one voice
- produce a coherent answer
- silently lose the exact private burden still left to the operator

This page exists to stop that.

## What this page is and is not allowed to prove

This page is authoritative about:

- where to start for a given question
- which truth register each route should lean on
- how to avoid weaker neighboring questions

This page is not authoritative about:

- whether the runtime already behaves correctly
- whether one architecture path already won
- whether better retrieval means the implementation is more mature

## Start from the wound, not the folder

Do not start from folder names.
Start from the actual question and the actual wound.

For each reading pass, identify:

1. the real question you are trying to answer
2. the smaller neighboring question you must avoid
3. which truth layer is allowed to carry the answer
4. what stronger sentence must remain forbidden at the end

If you cannot name those four things, the route is too loose.

## The main truth registers

Use these deliberately:

- Runtime truth:
  root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml),
  active `compose/` fragments, and validation commands
- Intent truth:
  [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- Plan truth:
  [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- Research pressure:
  [research/evidence-ledger.md](research/evidence-ledger.md) and related evidence pages
- Archive pressure:
  [operations/source-assimilation-index.md](operations/source-assimilation-index.md)

## Route 1: "What is the user actually trying to build?"

Use when you need the dream, not present-tense proof.

Start with:

1. [architecture/problem-and-goals.md](architecture/problem-and-goals.md)
2. [research/user-intent-and-dream.md](research/user-intent-and-dream.md)
3. [architecture/design-tensions-and-contradictions.md](architecture/design-tensions-and-contradictions.md)

Primary truth register:

- intent truth

Do not upgrade into:

- `the runtime already does this`

## Route 2: "What does the current Compose runtime materially contain?"

Use when you need current implementation shape.

Start with:

1. [architecture/current-compose-runtime.md](architecture/current-compose-runtime.md)
2. [architecture/compose-fragment-map.md](architecture/compose-fragment-map.md)
3. [architecture/request-path-and-failure-walkthrough.md](architecture/request-path-and-failure-walkthrough.md)
4. [architecture/stateful-ha-and-data.md](architecture/stateful-ha-and-data.md)

Primary truth register:

- runtime truth

Do not upgrade into:

- `because the runtime is broad, the missing truth must already be system-owned`

## Route 3: "Why is wrong-node behavior the real benchmark?"

Use when a route or failover claim sounds too broad.

Start with:

1. [architecture/request-path-and-failure-walkthrough.md](architecture/request-path-and-failure-walkthrough.md)
2. [architecture/ha-failover-routing.md](architecture/ha-failover-routing.md)
3. [architecture/failure-model-and-maturity.md](architecture/failure-model-and-maturity.md)

Primary truth registers:

- runtime truth
- intent truth

Do not upgrade into:

- `one local happy-path success means the platform is basically anti-SPOF`

## Route 4: "What still has to be proven before stronger claims become honest?"

Use when you need claim boundaries and drill classes.

Start with:

1. [operations/devops-runbook.md](operations/devops-runbook.md)
2. [operations/operator-questions-and-honest-answers.md](operations/operator-questions-and-honest-answers.md)
3. [research/evidence-ledger.md](research/evidence-ledger.md)
4. [research/ingress-and-failover-evidence.md](research/ingress-and-failover-evidence.md)
5. [research/stateful-ha-evidence.md](research/stateful-ha-evidence.md)

Primary truth registers:

- proof discipline
- runtime evidence

Do not upgrade into:

- `the existence of a disciplined proof language means the implementation must be close`

## Route 5: "Why do nearby options keep feeling fake?"

Use when the repo sounds like a product comparison problem.

Start with:

1. [architecture/missing-middle-layer.md](architecture/missing-middle-layer.md)
2. [architecture/orchestration-options.md](architecture/orchestration-options.md)
3. [architecture/capability-gaps-and-roadmap.md](architecture/capability-gaps-and-roadmap.md)
4. [research/orchestrator-tradeoffs-evidence.md](research/orchestrator-tradeoffs-evidence.md)

Primary truth registers:

- intent truth
- research pressure

Do not upgrade into:

- `a more famous controller is automatically a more honest answer`

## Route 6: "What is the repo's harsh answer on stateful services?"

Use when availability language starts sounding too generous.

Start with:

1. [architecture/stateful-ha-and-data.md](architecture/stateful-ha-and-data.md)
2. [architecture/failure-model-and-maturity.md](architecture/failure-model-and-maturity.md)
3. [research/stateful-ha-evidence.md](research/stateful-ha-evidence.md)
4. [research/stateful-ha-plan.md](research/stateful-ha-plan.md)

Primary truth registers:

- runtime truth
- evidence pages

Do not upgrade into:

- `reachable through Traefik` or `restartable elsewhere` equals HA

## Route 7: "What should I edit first if a page still feels vague?"

Start with:

1. [index.md](index.md)
2. [architecture/problem-and-goals.md](architecture/problem-and-goals.md)
3. [architecture/current-compose-runtime.md](architecture/current-compose-runtime.md)
4. [research/evidence-ledger.md](research/evidence-ledger.md)
5. [operations/source-assimilation-index.md](operations/source-assimilation-index.md)

Edit pages that still fail to leave behind:

- the hidden burden still private
- the truth register carrying the answer
- the stronger sentence still forbidden
- the next artifact or drill needed for a stronger claim

## What a good reading pass should leave behind

A good reading route should leave you with a small packet:

- the request or failure class you are reasoning about
- the strongest runtime artifact you relied on
- the strongest intent, plan, or research artifact you relied on
- the private burden still surviving afterward
- the sentence still forbidden
- the next proof packet required

If you finish with only:

- `the site is organized`
- `the architecture is clearer`
- `there are several plausible directions`

then the route was too weak for this repo.
