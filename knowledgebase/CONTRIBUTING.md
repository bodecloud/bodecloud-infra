# Contributing to the Knowledgebase

This knowledgebase is not a generic docs folder. It is part of the control surface for `bolabaden-infra`.

That means documentation changes here are expected to do real work:

- clarify what the system actually is
- clarify what the system is trying to become
- prevent architecture claims from outrunning the evidence

If a page only sounds polished, it is not useful enough.

If a page only becomes easier to read by making the system sound more settled
than it is, it is actively harmful.

This repo needs a stricter contribution standard than normal docs because the
most dangerous regressions here are often:

- smoother summaries
- tidier explanations
- cleaner separation of topics

when those improvements quietly delete the hidden operator burden or the
remaining ownership ambiguity.

That is why contributors should not think of this directory as a normal docs
polish surface.
It is closer to a truth-discipline surface.

Many infra repos regress when contributors make them sound more professional.
This repo regresses when contributors make it sound more settled than the
runtime, proof chain, and hidden operator burden actually support.

## The documentation contract

Every meaningful infrastructure statement belongs to one of three layers:

### 1. Live implementation truth

Use this layer when the statement is proven directly by the tracked repo and merged root Compose implementation.

Examples:

- services present in `docker compose config --services`
- Traefik labels present in the root stack
- included fragments under `compose/`
- secrets, configs, ports, networks, and healthchecks defined in tracked files

### 2. Planned architecture truth

Use this layer when the repo clearly intends something, but the live root implementation does not yet fully prove it.

Examples:

- current-state registry concepts such as `services.yaml`
- OpenSVC promotion ideas
- future stateful failover patterns
- proposed convergence or control-plane mechanisms

### 3. Research-pressure truth

Use this layer when the value of a page is preserving the recurring questions, comparisons, and frustrations that drove the repo here.

Examples:

- source-archive synthesis
- orchestrator comparisons
- anti-SPOF design exploration
- reverse-proxy and failover tradeoff research

Every contribution should preserve one more distinction inside those layers:

- what the runtime itself owns
- what a human operator still has to reconstruct around the runtime

And there is a second distinction contributors should keep asking about:

- which options in the repo are genuinely different because they move truth
  ownership
- which options only look different because they rename the same hidden burden
  with a new tool or cleaner architecture story

## What not to do

Do not:

- flatten the three truth layers into one smooth story
- call a route resilient just because the proxy is up
- call a stateful service HA just because it is reachable
- call a planning artifact a live control surface
- smooth away contradictions to make the docs feel cleaner
- replace a repo-specific problem with generic "best practice" language
- narrate a path as if it is self-explaining when the explanation still depends
  on private operator memory
- let summary pages sound calmer than the proof ceiling

This repo is trying to answer a very specific question. Generic cloud-platform prose usually makes the docs worse.

It is also trying to answer a much stranger question than most infrastructure
docs ever need to hold:

can several ordinary Docker nodes behave like one request-preserving, readable,
anti-SPOF personal cloud without the intelligence of the system still living in
one person's head?

If a contribution quietly swaps that question for a more normal one, the
contribution is off target even if its individual sentences are accurate.

That includes swaps like:

- "which orchestrator is best?"
- "how should we cluster Docker?"
- "which proxy should front multiple nodes?"

Those can all be useful subquestions.
They become harmful when they displace the harder repo-shaped question they are
supposed to serve.

## How to decide where a change belongs

If the change answers "what is actually running now?":

- edit [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
- possibly edit [`architecture/compose-fragment-map.md`](architecture/compose-fragment-map.md)

If the change answers "what does failover really mean here?":

- edit [`architecture/ha-failover-routing.md`](architecture/ha-failover-routing.md)
- possibly edit [`research/ingress-and-failover-evidence.md`](research/ingress-and-failover-evidence.md)

If the change answers "what is still missing or unproven?":

- edit [`architecture/failure-model-and-maturity.md`](architecture/failure-model-and-maturity.md)
- edit [`architecture/capability-gaps-and-roadmap.md`](architecture/capability-gaps-and-roadmap.md)

If the change answers "how should contributors verify claims?":

- edit [`operations/devops-runbook.md`](operations/devops-runbook.md)
- possibly edit [`research/evidence-ledger.md`](research/evidence-ledger.md)

If the change affects the repo's overall framing:

- edit [`index.md`](index.md)
- possibly edit the repo-root [`../README.md`](../README.md)

If the change affects how contributors should judge truth claims or narration
discipline:

- edit [`AGENTS.md`](AGENTS.md)
- edit [`operations/source-assimilation-index.md`](operations/source-assimilation-index.md)

## Retrieval workflow for contributors

Before writing or revising a page, stop and fill in this mental checklist:

1. What is the real question this page is supposed to answer?
2. What smaller neighboring question would be easier to answer but would miss
   the user's real frustration?
3. Which truth layer actually supports the answer?
4. Which evidence only supports pressure, comparison, or future direction?
5. What would count as overclaiming on this page?

Examples of the wrong substitution:

- "What is the best orchestrator?" instead of "Which layer would own the
  specific missing truth without charging too much worldview tax?"
- "Can traffic hit more than one node?" instead of "Has the wrong healthy node
  stopped being semantically wrong?"
- "Can the service restart?" instead of "Did truth stop living in one brittle
  place?"

If you cannot answer the checklist cleanly, do not start drafting yet.

## Evidence expectations

Use the smallest honest claim your evidence supports.

Examples:

- `docker compose config --quiet` proves config validity, not resilience
- `docker compose config --services` proves inclusion in the merged stack, not healthy runtime
- `docker compose ps` proves runtime presence, not successful routing
- successful HTTP requests prove route reachability, not peer failover
- successful peer failover for stateless HTTP does not prove stateful HA

Also apply this test:

- a successful-looking scenario does not prove the system owned the truth needed
  to make it succeed; it may only prove that an operator could still reconstruct
  the right answer during or after the fact

When in doubt, narrow the claim rather than widening the evidence.

Also use this test:

- if the page sounds more mature after your rewrite, did the evidence get
  stronger, or did the narration merely get smoother?

In this repo, smoother narration without stronger evidence is usually a
regression.

## Writing rules

- Prefer direct, plain language over soft architecture jargon.
- Preserve important contradictions instead of editing them away.
- Explain why a component matters, not just that it exists.
- Name the failure mode a mechanism is supposed to solve.
- If a page contains examples that are illustrative rather than live, label them clearly.
- If a repo-native artifact is missing, say it is missing.
- If the runtime still depends on socially reconstructed truth, say that
  bluntly.
- If a page gets clearer by asking a smaller question than the user is asking,
  undo that move.

## Validation

Before finishing documentation edits, run:

```bash
python3 -m mkdocs build -f mkdocs.yml --strict
```

If your changes talk about the live Compose surface, also run:

```bash
docker compose config --quiet
docker compose config --services
```

Use stronger checks when the claim requires them.

If the claim is about wrong-node handling, fallback ownership, policy
continuity, or the disappearance of sacred-node knowledge, a strict MkDocs pass
and a Compose render are necessary but not sufficient.

Those claims need stronger evidence and narrower language discipline.

## Final standard

A good contribution to this knowledgebase should make at least one of these more explicit:

- the user's real infrastructure dream
- the repo's actual current implementation boundary
- the difference between live proof and planned ambition
- the next honest move the repo should make
- whether the apparent intelligence of the system is actually present in the
  runtime or still living in the operator's head

An excellent contribution also makes at least one fake option less seductive.

That means it becomes harder for a reader to walk away thinking:

- "multi-node" itself solved something
- a named future helper already exists because it is described cleanly
- a bigger platform is probably the answer simply because the current state
  feels awkward
- a route surviving the easy case means the request path is broadly preserved

If it does not improve one of those, it is probably not done yet.
