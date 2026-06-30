# Legacy Docs Index

This `docs/` directory is no longer the primary documentation surface for understanding `bolabaden-infra`.

The current authoritative documentation lives under:

- [`knowledgebase/`](../knowledgebase/)
- especially [`knowledgebase/index.md`](../knowledgebase/index.md)

This matters because earlier docs in this repo flattened too many things together:

- live Compose behavior
- planned failover architecture
- research pressure from the source archive
- optimistic platform language that sounded more complete than the implementation really was

The newer knowledgebase separates those layers explicitly.

It also does something harsher than a normal docs index is expected to do:

- it keeps the user's real dream large
- it keeps the available options under suspicion
- it refuses to let planning coherence impersonate runtime proof

That is necessary because this repo has enough material to sound mature while
still leaving the same hidden operator reconstruction burden in place.

## Start in the knowledgebase

If you need to understand what the repo is actually trying to build, read these first:

1. [`../knowledgebase/index.md`](../knowledgebase/index.md)
2. [`../knowledgebase/architecture/problem-and-goals.md`](../knowledgebase/architecture/problem-and-goals.md)
3. [`../knowledgebase/architecture/current-compose-runtime.md`](../knowledgebase/architecture/current-compose-runtime.md)
4. [`../knowledgebase/architecture/compose-first-architecture.md`](../knowledgebase/architecture/compose-first-architecture.md)
5. [`../knowledgebase/architecture/ha-failover-routing.md`](../knowledgebase/architecture/ha-failover-routing.md)
6. [`../knowledgebase/architecture/stateful-ha-and-data.md`](../knowledgebase/architecture/stateful-ha-and-data.md)
7. [`../knowledgebase/architecture/capability-gaps-and-roadmap.md`](../knowledgebase/architecture/capability-gaps-and-roadmap.md)
8. [`../knowledgebase/operations/devops-runbook.md`](../knowledgebase/operations/devops-runbook.md)

Those pages are now the best explanation of the actual repo question:

> how do you make multiple ordinary Docker nodes feel resilient, peer-aware, and low-bullshit without immediately falling into Swarm, Kubernetes, or another heavyweight control plane?

There is a second question hiding under that one:

> when the ecosystem seems to offer lots of products and patterns, which of
> them are real options here and which are just better-dressed versions of the
> same wound?

If you skip that second question, the repo becomes much easier to summarize and
much less useful to understand.

## What this `docs/` directory is still good for

This directory still contains important repo artifacts. They just should not be mistaken for the whole story.

That warning matters more than it usually would.

Many files in this directory are detailed enough to sound like platform
closure:

- a master plan can sound like a control plane already exists
- a failover plan can sound like failover already works
- an orchestration comparison can sound like the future path is basically
  chosen

The knowledgebase exists precisely to stop those upgrades.

### Planning anchors

- [Infrastructure Master Plan](INFRASTRUCTURE_MASTER_PLAN.md)
- [Stateful HA Plan](stateful_ha_plan.md)
- [OpenSVC Ingress HA](osvc_ingress_ha.md)
- [Orchestration Research 2026](orchestration_research_2026.md)

These are valuable because they reveal where the repo wants to go. They are not the same thing as proof of current live behavior.

They are also not automatically "real options" just because they are
well-described.

A real option in this repo is one that actually moves truth out of private
operator memory.
If a planning artifact cannot plausibly do that, it is still mostly a pressure
map rather than a route to closure.

### Product- or subsystem-specific docs

- [Maintenance Guide](MAINTENANCE.md)
- [OTLP Quickstart](OTLP_QUICKSTART.md)
- [KotorModSync Telemetry Setup](KOTORMODSYNC_TELEMETRY_SETUP.md)
- [KotorModSync Client Integration](KOTORMODSYNC_CLIENT_INTEGRATION.md)
- [KotorModSync Security Summary](KOTORMODSYNC_SECURITY_SUMMARY.md)

These are still useful, but they sit beside the larger infrastructure question rather than replacing it.

### Planning history

- [`plans/`](plans/)
- [`brainstorms/`](brainstorms/)
- [`residual-review-findings/`](residual-review-findings/)

These files help explain decisions, but they should not be treated as current runtime truth unless a live implementation page or verification artifact proves the same claim.

## Reading guide by question

If your question is "what is actually running now?":

- [`../knowledgebase/architecture/current-compose-runtime.md`](../knowledgebase/architecture/current-compose-runtime.md)
- [`../knowledgebase/architecture/compose-fragment-map.md`](../knowledgebase/architecture/compose-fragment-map.md)

If your question is "what does failover really mean here?":

- [`../knowledgebase/architecture/ha-failover-routing.md`](../knowledgebase/architecture/ha-failover-routing.md)
- [`../knowledgebase/research/ingress-and-failover-evidence.md`](../knowledgebase/research/ingress-and-failover-evidence.md)

If your question is "how honest is the stateful HA story?":

- [`../knowledgebase/architecture/stateful-ha-and-data.md`](../knowledgebase/architecture/stateful-ha-and-data.md)
- [`../knowledgebase/research/stateful-ha-evidence.md`](../knowledgebase/research/stateful-ha-evidence.md)
- [Stateful HA Plan](stateful_ha_plan.md)

If your question is "why not just pick Kubernetes, Swarm, Nomad, or OpenSVC?":

- [`../knowledgebase/architecture/orchestration-options.md`](../knowledgebase/architecture/orchestration-options.md)
- [`../knowledgebase/research/orchestrator-tradeoffs-evidence.md`](../knowledgebase/research/orchestrator-tradeoffs-evidence.md)
- [Orchestration Research 2026](orchestration_research_2026.md)

If your question is "why does this still feel like there are no real options
even though there are plenty of nouns?":

- [`../knowledgebase/operations/operator-questions-and-honest-answers.md`](../knowledgebase/operations/operator-questions-and-honest-answers.md)
- [`../knowledgebase/architecture/problem-and-goals.md`](../knowledgebase/architecture/problem-and-goals.md)
- [`../knowledgebase/architecture/missing-middle-layer.md`](../knowledgebase/architecture/missing-middle-layer.md)

## Important boundary

Do not use this directory as evidence that the repo already has:

- fully proven multi-node failover
- a live tracked root `services.yaml` current-state registry
- universal peer-aware request success
- zero-SPOF stateful behavior
- one settled future control plane

Those are exactly the kinds of overclaims the knowledgebase was rewritten to avoid.

The docs fail again if they let a reader walk away thinking:

- "there are plenty of options, we just need to pick one"
- "the right proxy plus DNS probably already solves most of it"
- "a more serious control plane will automatically remove the real burden"

Those are the same adjacent answers the user is already reacting against.

## Bottom line

Use this folder as:

- a planning archive
- a set of subsystem docs
- a source layer

Use the knowledgebase as:

- the current authoritative explanation
- the honesty boundary
- the place where live truth, planned truth, and research pressure are kept separate on purpose

And use this `docs/` tree as a source archive of serious planning pressure, not
as a calm summary that has already turned the pressure into closure.
