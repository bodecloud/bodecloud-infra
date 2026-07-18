# Knowledgebase README

This file is the repo-facing doorway into the MkDocs site rooted at
[`index.md`](index.md).

Contributors usually read a repo README before they read the rendered site.
That makes this file important, because `bolabaden-infra` is easy to misread as:

- a Compose homelab
- a reverse-proxy stack
- an orchestrator comparison project
- a "future HA ideas" repo

Those are all adjacent.
They are not the real problem.

## The shortest correct reading

The repo is trying to answer this question:

> how do you keep `docker-compose.yml` as the human control surface, spread
> services across several ordinary Docker nodes, and still make wrong-node
> traffic, fallback, and anti-SPOF behavior feel like one coherent platform
> instead of one operator privately remembering the real answer?

That question should remain active while reading or editing anything in this
directory.

## What this directory is for

This directory is not here to provide generic infrastructure prose.
It is here to keep the documentation from becoming more certain, more mature,
or more platform-like than the current evidence allows.

The practical failure this directory is trying to stop is:

1. the stack looks serious
2. the docs sound coherent
3. the reader starts assuming the hidden burden must be lower
4. the runtime still needs private operator completion on the bad day

If a page helps that sequence happen, it failed even if the prose sounds smart.

## The hidden burden the docs must keep naming

The main unresolved burden is not generic complexity.
It is private custody of truths such as:

- what runs where right now
- which peer is actually valid now
- whether fallback survives backend loss
- whether a forwarded protected route still means the same thing
- whether a stateful service only looks movable from the outside

Most nearby tooling improves naming, automation, or visibility without clearly
moving those truths into the system.
That is why the site has to stay harsher than a normal architecture knowledge
base.

## How to read the site correctly

Keep these authority layers separate:

- live runtime truth:
  [`../docker-compose.yml`](../docker-compose.yml),
  [`../compose/`](../compose), and real validation commands
- repo-native intent:
  especially
  [`../.github/copilot-instructions.md`](../.github/copilot-instructions.md)
- planning truth:
  especially
  [`../docs/INFRASTRUCTURE_MASTER_PLAN.md`](../docs/INFRASTRUCTURE_MASTER_PLAN.md)
- archive and research pressure:
  [`research/`](research) and [`operations/source-assimilation-index.md`](operations/source-assimilation-index.md)

If those collapse into one calm voice, the docs start laundering intent and
pressure into present-tense capability.

## What a useful page in this directory should leave behind

After reading a page, a contributor should be able to say:

- what question the page was really answering
- which truth layer carried the answer
- which burden still remains operator-owned
- which stronger sentence is still forbidden
- what artifact or drill would make that stronger sentence legal

If the page leaves behind only:

- "the architecture is clearer"
- "the options are broader"
- "the docs feel mature"

then it is still too weak for this repo.

## Where to start

If you are entering from the repo, start with:

- [`index.md`](index.md)

Then follow this order:

1. [`architecture/problem-and-goals.md`](architecture/problem-and-goals.md)
2. [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
3. [`architecture/ha-failover-routing.md`](architecture/ha-failover-routing.md)
4. [`architecture/stateful-ha-and-data.md`](architecture/stateful-ha-and-data.md)
5. [`architecture/capability-gaps-and-roadmap.md`](architecture/capability-gaps-and-roadmap.md)
6. [`operations/devops-runbook.md`](operations/devops-runbook.md)

If you need the pressure and evidence behind those pages, go next to:

- [`research/evidence-ledger.md`](research/evidence-ledger.md)
- [`research/ingress-and-failover-evidence.md`](research/ingress-and-failover-evidence.md)
- [`research/stateful-ha-evidence.md`](research/stateful-ha-evidence.md)
- [`research/orchestrator-tradeoffs-evidence.md`](research/orchestrator-tradeoffs-evidence.md)

## What this README is and is not allowed to prove

This file is allowed to prove:

- how the knowledgebase should be entered
- what question the site is really trying to preserve
- which reading routes are safest

This file is not allowed to prove:

- that the runtime already crossed the thresholds the route describes
- that better site organization means lower operator burden
- that the repo is already close to its final architecture

This is a doorway contract, not runtime proof.
