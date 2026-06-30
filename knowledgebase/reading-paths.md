# Reading Paths and Retrieval Routes

This page is not navigation chrome.
It is a retrieval discipline page for a repo whose problem is now large enough
that a reader can gather many true fragments and still produce the wrong answer.

The main retrieval failure looks like this:

1. gather some runtime facts
2. gather some intent and roadmap language
3. gather some research pressure
4. blend them into one coherent voice
5. accidentally answer a smaller, calmer question than the user is really asking

That is exactly how documentation starts sounding comprehensive while still
failing the operator.

This page exists to stop that.

## The reading discipline in one sentence

Do not read from folder names.
Read from the wound.

That means every pass should start by naming:

1. the actual question
2. the smaller neighboring question you must refuse
3. the truth register allowed to carry the answer
4. the stronger sentence that must still remain forbidden
5. the private sentence that might still survive afterward

If you cannot name those five things, the retrieval path is already too loose.

## The main truth registers

Use these deliberately and do not let one impersonate another.

### Runtime truth

Primary surfaces:

- root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active `compose/` include fragments
- `docker compose config` class validation

This register can prove:

- what is materially authored in the priority implementation
- what services, networks, labels, configs, and fragments exist
- what the Compose-first runtime shape actually contains

This register cannot prove by itself:

- generic wrong-node dignity
- preserved semantics under peer forwarding
- trustworthy stateful authority

### Intent truth

Primary surfaces:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [architecture/instruction-surfaces-and-authority.md](architecture/instruction-surfaces-and-authority.md)

This register can prove:

- what the repo is explicitly trying to become
- which kinds of fake-HA language the repo itself is trying to avoid
- why no-Swarm and no-heavyweight-orchestrator-by-default is part of the goal

This register cannot prove:

- that the runtime already behaves that way
- that repeated intention language equals convergence

### Plan truth

Primary surfaces:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- roadmap-oriented architecture and research pages

This register can prove:

- what work the repo thinks still needs to happen
- what staged promotion logic exists

This register cannot prove:

- that planned work is already earning live claims

### Research pressure

Primary surfaces:

- [research/evidence-ledger.md](research/evidence-ledger.md)
- related evidence pages

This register can prove:

- what outside patterns or tradeoffs are pressing on the architecture
- what stronger claim would still need better evidence

This register cannot prove:

- that the repo has already internalized the best external answer

### Archive pressure

Primary surfaces:

- [operations/source-assimilation-index.md](operations/source-assimilation-index.md)

This register can prove:

- which sources keep recurring in the repo's documentation pressure
- how to separate source classes without flattening them into one voice

This register cannot prove:

- that volume of source material equals runtime maturity

## Route 1: "What is the user actually trying to build?"

Use this when the real need is the dream, not present-tense proof.

Start with:

1. [architecture/problem-and-goals.md](architecture/problem-and-goals.md)
2. [research/user-intent-and-dream.md](research/user-intent-and-dream.md)
3. [architecture/design-tensions-and-contradictions.md](architecture/design-tensions-and-contradictions.md)
4. [architecture/instruction-surfaces-and-authority.md](architecture/instruction-surfaces-and-authority.md)

Primary truth register:

- intent truth

Question this route is allowed to answer:

- what exact platform dream is being protected here?

Smaller neighboring question it must refuse:

- what generic HA or clustering story would sound similar?

Stronger sentence still forbidden afterward:

- `the runtime already behaves this way end to end`

Private sentence likely still surviving:

- `I still personally know which parts of the dream are still only intent.`

## Route 2: "What does the current Compose runtime materially contain?"

Use this when the need is current implementation shape and burden distribution.

Start with:

1. [architecture/current-compose-runtime.md](architecture/current-compose-runtime.md)
2. [architecture/compose-fragment-map.md](architecture/compose-fragment-map.md)
3. [architecture/request-path-and-failure-walkthrough.md](architecture/request-path-and-failure-walkthrough.md)
4. [architecture/stateful-ha-and-data.md](architecture/stateful-ha-and-data.md)

Primary truth register:

- runtime truth

Question this route is allowed to answer:

- what serious platform surface already exists in the priority Compose runtime?

Smaller neighboring question it must refuse:

- does serious-looking runtime breadth already imply shared truth ownership?

Stronger sentence still forbidden afterward:

- `because the stack is broad and modular, wrong-node success must be close`

Private sentence likely still surviving:

- `I still personally know which sophisticated parts are only sophisticated locally.`

## Route 3: "Why is wrong-node behavior the real benchmark?"

Use this when a route, failover, or anti-SPOF claim starts sounding too broad.

Start with:

1. [architecture/request-path-and-failure-walkthrough.md](architecture/request-path-and-failure-walkthrough.md)
2. [architecture/ha-failover-routing.md](architecture/ha-failover-routing.md)
3. [architecture/current-state-registry-and-peer-eligibility.md](architecture/current-state-registry-and-peer-eligibility.md)
4. [architecture/failure-model-and-maturity.md](architecture/failure-model-and-maturity.md)
5. [research/ingress-and-failover-evidence.md](research/ingress-and-failover-evidence.md)

Primary truth registers:

- runtime truth
- intent truth

Question this route is allowed to answer:

- why is wrong-node dignity stricter than "reachable ingress"?

Smaller neighboring question it must refuse:

- is there some fallback-looking mechanism somewhere in the tree?

Stronger sentence still forbidden afterward:

- `a rendered route is close enough to a preserved request`

Private sentence likely still surviving:

- `I still personally know which peer forward really preserves meaning.`

## Route 4: "What still has to be proven before stronger claims become honest?"

Use this when the issue is proof class, drills, or claim boundaries.

Start with:

1. [operations/devops-runbook.md](operations/devops-runbook.md)
2. [operations/operator-questions-and-honest-answers.md](operations/operator-questions-and-honest-answers.md)
3. [research/evidence-ledger.md](research/evidence-ledger.md)
4. [research/ingress-and-failover-evidence.md](research/ingress-and-failover-evidence.md)
5. [research/stateful-ha-evidence.md](research/stateful-ha-evidence.md)

Primary truth registers:

- proof discipline
- runtime evidence

Question this route is allowed to answer:

- what exact proof packet is still missing?

Smaller neighboring question it must refuse:

- do the docs at least sound careful enough?

Stronger sentence still forbidden afterward:

- `because the proof language is disciplined, the implementation must be close`

Private sentence likely still surviving:

- `I still personally know which drill results matter most because they have not become routine system truth yet.`

## Route 5: "Why do nearby options keep feeling fake?"

Use this when the repo starts sounding like a product comparison problem.

Start with:

1. [architecture/missing-middle-layer.md](architecture/missing-middle-layer.md)
2. [architecture/current-state-registry-and-peer-eligibility.md](architecture/current-state-registry-and-peer-eligibility.md)
3. [architecture/orchestration-options.md](architecture/orchestration-options.md)
4. [architecture/capability-gaps-and-roadmap.md](architecture/capability-gaps-and-roadmap.md)
5. [research/orchestrator-tradeoffs-evidence.md](research/orchestrator-tradeoffs-evidence.md)

Primary truth registers:

- intent truth
- research pressure

Question this route is allowed to answer:

- why do so many adjacent tools still leave the same hidden burden alive?

Smaller neighboring question it must refuse:

- which orchestrator is most famous, lightest, or most modern?

Stronger sentence still forbidden afterward:

- `a more dynamic controller automatically kills the real wound`

Private sentence likely still surviving:

- `I still personally know which burden a new tool actually removes instead of merely renaming it.`

## Route 6: "What is the repo's harsh answer on stateful services?"

Use this when availability language starts sounding too generous.

Start with:

1. [architecture/stateful-ha-and-data.md](architecture/stateful-ha-and-data.md)
2. [architecture/failure-model-and-maturity.md](architecture/failure-model-and-maturity.md)
3. [research/stateful-ha-evidence.md](research/stateful-ha-evidence.md)
4. [research/stateful-ha-plan.md](research/stateful-ha-plan.md)

Primary truth registers:

- runtime truth
- evidence pages

Question this route is allowed to answer:

- what is actually known versus merely reachable for stateful surfaces?

Smaller neighboring question it must refuse:

- can the service be exposed, restarted, or contacted somewhere?

Stronger sentence still forbidden afterward:

- `reachable through a serious stack equals resilient state authority`

Private sentence likely still surviving:

- `I still personally know whether the stateful answer is authoritative or just currently reachable.`

## Route 7: "What should I edit first if the docs still feel vague?"

Use this when the knowledgebase sounds polished but not reconstructive enough.

Start with:

1. [index.md](index.md)
2. [architecture/problem-and-goals.md](architecture/problem-and-goals.md)
3. [architecture/current-compose-runtime.md](architecture/current-compose-runtime.md)
4. [research/evidence-ledger.md](research/evidence-ledger.md)
5. [operations/source-assimilation-index.md](operations/source-assimilation-index.md)
6. [architecture/instruction-surfaces-and-authority.md](architecture/instruction-surfaces-and-authority.md)

Edit pages that still fail to leave behind:

- the exact hidden burden still private
- the truth register carrying the answer
- the stronger sentence still forbidden
- the next artifact or drill required for a stronger claim
- the smaller neighboring question the page accidentally answered instead

The page is also weak if it sounds intelligent while still leaving no answer to:

- what exact operator sentence survives when the page ends?

## Route 8: "How do I actually use this site like a RAG pass?"

Use this when you need a synthesis packet instead of a casual browse.

Sequence:

1. identify the exact request or failure class
2. choose one primary route above
3. choose one contrast route that could falsify an overly smooth answer
4. extract the strongest runtime artifact you can name
5. extract the strongest intent, plan, or research artifact you can name
6. write down the private sentence still surviving afterward
7. write down the stronger sentence still forbidden
8. write down the smaller neighboring question you successfully refused

If you skip the contrast route, the answer usually becomes too smooth.
That is one of the main reasons the earlier docs felt useless.

## Recommended route pairs for common questions

Use these pairings when the answer is likely to drift into overclaiming.

| Real question | Primary route | Contrast route | Why the contrast matters |
| --- | --- | --- | --- |
| `What is this repo trying to become?` | Route 1 | Route 2 | stops dream language from impersonating runtime proof |
| `What is really live today?` | Route 2 | Route 4 | stops breadth-of-runtime from impersonating passed drills |
| `Does the repo actually solve failover?` | Route 3 | Route 6 | stops HTTP-looking answers from hiding stateful and L4 gaps |
| `Do we need Nomad, k3s, Swarm, or something else?` | Route 5 | Route 2 | stops controller discussion from ignoring the current Compose burden |
| `What should be documented next?` | Route 7 | Route 4 | stops cleanup edits from losing proof discipline |

## What a good reading pass must leave behind

A good pass should leave a compact packet containing:

- the exact request or failure class
- the strongest runtime artifact used
- the strongest intent, plan, or research artifact used
- the private sentence still surviving afterward
- the stronger sentence still forbidden
- the next proof packet required
- the smaller neighboring question successfully avoided

If the pass leaves behind only:

- a coherent story
- a list of pages
- a technology comparison
- a sense that the repo is sophisticated

then the pass is still too weak for this codebase.

## The shortest honest rule

Retrieval in this repo is successful only when it makes the hidden operator job
more visible, not less.
