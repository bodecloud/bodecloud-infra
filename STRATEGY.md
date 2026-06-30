***

name: bolabaden.org Infrastructure
last_updated: 2026-06-30

***

# bolabaden.org Infrastructure Strategy

## What problem this strategy is actually solving

The problem is not "how do we run containers."

The real problem is:

> how do we make multiple ordinary Docker nodes behave like a resilient, operator-readable system without immediately paying the full tax of Swarm, Kubernetes, or another heavyweight orchestrator?

That problem breaks down into several linked pressures:

- requests may land on any node
- the relevant service may not be on the node that received the request
- routing, auth, middleware, and health semantics should survive failover
- secrets, env, and placement truth should stop depending on operator memory
- stateful services should be treated honestly instead of being lumped into the same story as stateless web routing

The manual synchronization bottleneck is the heart of the issue. The repo keeps paying tax in:

- routing drift
- service placement ambiguity
- env and secret sprawl
- partial failover behavior
- too much hidden operator knowledge

This page also has to defend against a common failure mode in infrastructure
strategy writing:

the strategy sounds stronger simply because it names the right nouns.

That is not enough for this repo.

The user is not short on nouns.
They are short on options that still feel real after the first request lands on
the wrong node, the local backend disappears, or the operator has to answer
"what actually runs where right now?" without resorting to private memory.

So this strategy is not trying to produce a respectable posture.
It is trying to keep the repo aligned with a harsher benchmark:

> does the next coordination layer actually remove hidden burden, or does it
> mostly improve the story around the same hidden burden?

There is a second benchmark implied by that one:

> when the ecosystem appears to offer many products, does this repo actually
> gain another real option, or just another better-dressed version of the same
> wound?

That question matters because the user frustration is not just "too few
features."
It is that too many supposed options stop one layer before truth actually moves
out of the operator's head.

## Guiding policy

The repo's strategy is to stay Compose-first for as long as that remains the least painful honest choice, while deliberately adding only the extra coordination layers that solve a real failure mode.

That means:

- keep Compose readable and operator-facing
- reject fake HA language
- improve ingress, placement, and convergence before pretending the system has solved state
- keep multiple stronger-control-plane options open until one clearly earns its tax

This is not anti-automation.

It is anti-premature-platform-ideology.

The goal is not to avoid a control plane forever. The goal is to avoid adopting a heavier one before the repo can justify exactly what pain it removes.

That last clause matters more than it usually does in strategy documents.

For this repo, "heavier" is not just about CPU, RAM, or setup complexity.
It is about worldview cost:

- what becomes harder to inspect directly
- what becomes harder to reason about causally
- what gets hidden behind a more authoritative control surface
- and whether that hiddenness is actually paying for itself by removing a real
  sacred-node, wrong-node, or convergence pain

If a future control plane mostly changes legitimacy, fashion, or naming, but
not the real burden location, it has still not earned itself here.

That is why this strategy should also be read as an anti-fake-option document.

It is not enough for a candidate layer to be:

- popular
- cluster-shaped
- highly automated
- widely described as HA

It has to change who owns the truth on the bad day.

## What this strategy believes

### 1. Compose is still the baseline human contract

The repo is built around the root [`docker-compose.yml`](docker-compose.yml) and included fragments under [`compose/`](compose). That remains the clearest live surface for the operator.

### 2. Compose alone is not enough for honest multi-node behavior

The stack already shows the need for stronger layers around:

- placement truth
- route generation
- health-aware failover
- DNS and ingress convergence
- bootstrap and node join consistency

### 3. The next useful layer is not automatically Kubernetes

The strategy keeps several futures open because the repo is still working out where extra complexity should live:

- Compose plus helpers and registry-driven routing
- failover-agent or Constellation-style convergence logic
- OpenSVC or similar infra-grade HA
- Nomad or a comparable lighter scheduler
- k3s or Kubernetes only when the capability threshold is clearly crossed

This also means the repo should resist category-based comfort.

It is not enough to say:

- "this is the lightweight option"
- "this is the serious option"
- "this is the scalable option"

Those labels are too easy to map onto the same fake-option market the user is
already reacting against.

The only labels that matter are closer to:

- which hidden truth would this own directly
- which hidden truth would still be socially reconstructed
- which humiliation on the bad day would still survive even after promotion

That third bullet is a better option filter than most ordinary strategy docs
ever use.

The repo should keep asking:

- after adopting this layer, what still breaks unless one human knows the
  missing context privately?

If the answer is "most of the important parts," the layer is still mostly
cosmetic no matter how mature it sounds.

### 4. State is the honesty boundary

The strategy treats ingress failover and stateful resilience as different categories.

It is acceptable to improve:

- HTTP request continuity
- local-first then peer-fallback routing
- DDNS and ingress survivability

without pretending that Redis, MongoDB, Postgres, RabbitMQ, or other state-bearing systems are thereby solved.

## Coherent action

The coherent action for this repo is to reduce operator burden in the order that most directly supports the actual dream.

### Track 1: Placement and convergence truth

Build and verify the smallest honest current-state layer that answers:

- what runs where
- what can fail over where
- what should restart locally
- what should be picked up elsewhere

This includes work around:

- a real placement-truth surface such as the repeatedly referenced `services.yaml`
- safer route generation
- sync and convergence logic
- explicit separation between live placement truth and planned placement models

Why this matters:

Without placement truth, "multi-node" remains partly folklore.

There is a stronger way to say the same thing:

without placement truth, the operator is still partially functioning as the
control plane even if the repo already contains a large amount of real
automation and edge machinery.

### Track 2: Ingress continuity without fake guarantees

Harden the any-node request model so the stack gets better at:

- local-first serving
- peer forwarding
- middleware continuity
- auth continuity
- health-aware route survival

This includes work around:

- Traefik routing semantics
- replacement of known-brittle failover generation paths like `docker-gen-failover`
- Cloudflare node-entry behavior
- private peer reachability

Why this matters:

This is the nearest honest resilience layer for the repo.

It is also the layer most likely to produce false confidence if it improves
first.

That is why this track has to stay tied to:

- wrong-node behavior
- route persistence under the exact failure that triggered fallback
- semantic continuity after peer handoff

rather than settling for "more nodes answer now" language.

It is also the place where the surrounding ecosystem most often tries to sell
the user an illusion of choice.

There are endless ways to make more nodes accept traffic.
There are far fewer ways to make those nodes preserve request meaning once they
accept it.

This track should therefore be judged by meaning preservation, not traffic
plurality.

### Track 3: Bootstrap and environment convergence

Reduce the amount of hidden memory required to bring up a node or keep one aligned.

This includes work around:

- bootstrap consistency
- secret and env hydration
- node identity and mesh join
- reducing tacit operator steps

Why this matters:

If joining or recovering a node still depends on private memory, the system keeps carrying human SPOFs even when containers look distributed.

This is one of the least glamorous parts of the strategy and one of the most
important.

A system can look advanced at the edge while still be spiritually single-node
because only one operator memory can really reconstruct it after a messy
recovery.

### Track 4: Honest stateful resilience

Treat state-bearing systems as first-class architecture problems rather than as extensions of the HTTP routing story.

This includes work around:

- replica and election semantics
- storage topology
- failover correctness
- client reconnect behavior
- L4 exposure boundaries

Why this matters:

This is where the repo either becomes genuinely trustworthy or quietly starts lying.

It is also where fake-option language becomes most dangerous.

For stateful systems, a new proxy, cluster, or replicated storage layer may
create new machinery without yet creating a trustworthy answer.

### Track 5: Control-plane justification and escalation

Keep stronger orchestration options alive, but judge them by explicit pain removed rather than industry fashion.

This includes work around:

- documenting promotion thresholds
- clarifying whether missing semantics belong in metadata, active agent logic, or a stronger scheduler
- preserving optionality without hiding indecision

Why this matters:

The repo is not choosing products. It is choosing where complexity should live.

It is also choosing which kinds of hiddenness are acceptable.

Some hiddenness is worth buying if it removes a concrete burden and leaves a
clear inspectable contract behind.
Some hiddenness is merely repackaged surrender.

The strategy needs to keep those two things separate or the whole repo drifts
back into ordinary infrastructure consumerism.

That drift would look like progress from far away.
It would also recreate the same user frustration under more respectable names.

## Success metrics

These metrics are only useful if they map to the actual anti-SPOF goal.

- **Request continuity under node loss**
  Percent of representative HTTP requests that still succeed when an entry node or serving node fails.

- **Failover correctness, not just detection**
  Time from service or node failure to restored healthy routing, with middleware/auth behavior still preserved.

- **Placement truth confidence**
  How often the operator can answer "what runs where right now?" from authoritative artifacts instead of inference or memory.

- **Configuration and secret convergence time**
  Time required for updated configuration or secret material to become consistent across the nodes that need it.

- **Bootstrap memory removed**
  Number of hidden manual decisions still required to turn a fresh node into a trustworthy participant.

- **Stateful recovery realism**
  Verified service-specific recovery behavior for Redis, MongoDB, Postgres, RabbitMQ, and similar systems, measured separately from ingress success.

## Milestones

### Near-term

- make the failover and route-generation story less brittle than the current `docker-gen-failover` path
- tighten placement-truth and convergence documentation until the next implementation surface is unambiguous
- keep proving the difference between DNS-level survivability, HTTP-level failover, and stateful correctness

### Medium-term

- establish a genuinely trustworthy placement and convergence layer
- reduce node bootstrap and recovery memory burden materially
- prove representative failover behavior on real service classes rather than only on architecture diagrams

### Longer-term

- decide which workloads remain Compose-first
- decide which workloads require stronger scheduler or HA semantics
- justify any move toward OpenSVC, Nomad, k3s, Kubernetes, or deeper Constellation ownership with concrete pain removed

## What this strategy is explicitly not doing

- claiming the repo already has universal multi-node failover
- treating DNS failover as equivalent to service resilience
- treating stateless routing improvements as proof of stateful HA
- committing the whole repo to Kubernetes or Swarm by ideology
- pretending the future control surface is already fully settled

## One-line strategy

Keep the stack Compose-first and operator-readable while deliberately adding the smallest honest coordination layers needed to remove SPOFs, survive request-path failures, and reduce hidden operator burden.

## Short message

`bolabaden-infra` is trying to build a personal-cloud middle ground: stronger than raw multi-node Compose sprawl, lighter than full orchestrator religion, and honest about where routing success ends and real HA begins.

The phrase "middle ground" is still slightly too gentle on its own.

This repo is really searching for an honesty frontier:

- more explicit shared truth than raw Compose sprawl
- less worldview capture than heavyweight orchestrator doctrine
- less fake comfort than the surrounding ecosystem keeps trying to sell

If the repo keeps that sentence intact, it stays aligned with the real dream.
If it loses that sentence, it will likely regress into one of two failure
modes:

- documentation that sounds mature while leaving the same burden in place
- promotion pressure toward a larger platform before the real thresholds have
  actually been crossed
