# Contributing to `bolabaden-infra`

This repo is not a normal app repo and it is not a generic "homelab docs"
repo.

It is a Compose-first infrastructure repo that is trying to answer a much
sharper question:

> how do multiple ordinary Docker nodes stop behaving like loosely related
> boxes, without immediately defaulting to Swarm, Kubernetes, or another
> heavyweight control plane that has not yet earned its keep?

That means contribution quality is not measured only by whether a change looks
tidier, more automated, or more "serious."

It is measured by whether the change makes the actual anti-fake-option dream
more true, or at least makes the proof boundary more honest.

## Read the repo in the right order

Before changing architecture-facing docs, read these first:

1. [`README.md`](README.md)
2. [`.github/copilot-instructions.md`](.github/copilot-instructions.md)
3. [`AGENTS.md`](AGENTS.md)
4. [`knowledgebase/index.md`](knowledgebase/index.md)
5. [`knowledgebase/architecture/problem-and-goals.md`](knowledgebase/architecture/problem-and-goals.md)
6. [`knowledgebase/research/user-intent-and-dream.md`](knowledgebase/research/user-intent-and-dream.md)
7. [`knowledgebase/research/evidence-ledger.md`](knowledgebase/research/evidence-ledger.md)

The reason for that order is simple:

- the root repo files explain the live project boundary
- the knowledgebase explains the real dream, the real proof boundary, and the
  real user frustration

If you skip that reading order, it becomes very easy to write docs that sound
smart while shrinking the actual ask.

## The main documentation rule

Do not flatten these evidence layers together:

- live root runtime truth
- repo-native intent
- planned architecture
- research pressure from the archive

If a page turns those into one smooth present-tense story, it is making the
site worse even if the prose gets cleaner.

This repo needs retrieval-like reconstruction, not cheerful summarization.

That means a good contribution preserves:

- contradictions
- uneven proof
- repeated archive pressure
- where a mechanism is real but still does not cross the user's real benchmark

## What the user is actually asking for

The user is not mainly asking for:

- more Compose fragments
- more orchestration brands
- more node count
- more polished HA vocabulary

The user is asking for a personal cloud that stays operator-readable while:

- any surviving node can become the first hop
- wrong-node requests stop turning into private operator reconstruction
- local-first serving remains possible
- peer-forward fallback stays honest when locality is absent
- stateful systems are described with much harsher truth than stateless HTTP

That is the benchmark your changes should be judged against.

## Fake options vs real options

This repo is unusually sensitive to fake options.

A change is not a real new option just because it introduces:

- a new orchestrator
- a new proxy
- a new helper
- a new sync path
- a new health layer

It becomes a real option only if it changes where hidden truth actually lives.

Ask these questions:

- does this reduce sacred-node dependence?
- does this reduce wrong-node humiliation?
- does this externalize current truth instead of assuming private memory?
- does this preserve request meaning under failure, not just route existence?
- does this solve a concrete pain, or mostly improve the story around the same
  pain?

If the answer is mostly "it sounds more mature now," the contribution is still
too shallow.

## What to edit where

### Edit the knowledgebase when

- you are clarifying architecture meaning
- you are separating live truth from plan truth
- you are reconstructing source archive pressure
- you are explaining failover, routing, stateful HA, orchestration tradeoffs,
  or proof limits

The knowledgebase is the primary explanation surface.

### Edit legacy `docs/` pages when

- the file is still useful as a planning artifact
- it needs a stronger reading boundary
- it needs clearer links back to the knowledgebase
- it still contains phrasing that overclaims what the runtime proves

Do not treat old planning docs as the main truth surface.

### Edit root repo files when

- you are clarifying onboarding order
- you are tightening repo-wide honesty boundaries
- you are improving how contributors interpret the project

## Compose authoring rules

When editing Compose files:

- prefer inline `configs:` content over external mounted config files when
  practical
- escape literal dollar signs as `$$` inside inline Compose config content
- do not "fix" services by weakening or removing healthchecks
- keep healthchecks meaningful enough to test actual service behavior, not just
  port reachability

Those rules are not cosmetic.
They are part of the repo's refusal to paper over brittle behavior with
lighter-looking YAML.

## Validation baseline

Minimum checks for documentation-facing work:

```bash
python3 -m mkdocs build -f mkdocs.yml --strict
```

Minimum checks for Compose-facing work, when env and secrets are prepared:

```bash
docker compose config --quiet
docker compose config --services
```

Remember:

- passing MkDocs proves the site renders, not that the claims are honest
- passing Compose config proves authored shape, not distributed correctness

This repo cares about both rendering discipline and claim discipline.
The second one is rarer and more important.

## Questions every serious contribution should survive

Before considering a change "done," ask:

1. Does this make the real dream more true, or only easier to narrate?
2. Did I preserve where the evidence stops?
3. Did I accidentally upgrade first-hop survivability into full failover?
4. Did I accidentally upgrade route existence into request-preserving behavior?
5. Did I accidentally upgrade reachability into stateful correctness?
6. Did I make it harder, not easier, to mistake a fake option for a real one?

If the page or change does not survive those questions, keep going.
