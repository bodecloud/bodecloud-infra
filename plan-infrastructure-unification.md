# Infrastructure Unification Plan

This file used to be a long essay about the pain of scaling multi-node
infrastructure.
Some of that diagnosis was real.
As a working document, it was too indirect and too easy to overread.

This rewrite keeps the core point and removes the theatrical sprawl.

## The actual problem

The repo does not mainly suffer from a lack of containers.

It suffers from too many partial truths spread across different surfaces:

- Compose definitions
- include selection
- env and secret material
- DNS and ingress behavior
- route-generation helpers
- node placement assumptions
- stateful-service behavior
- operator memory

That fragmentation is what "infrastructure unification" is supposed to fix.

This page also has to guard against a specific misunderstanding:

unification can sound like a cleanliness project when read too quickly.

For this repo it is not mainly about cleanliness.
It is about reducing the amount of private operator interpretation required
before the stack becomes coherent under stress.

That means the word "unification" should be heard as:

- less folklore
- less remembered exception handling
- less sacred-node dependence
- more inspectable shared truth at the moment a request, failure, or failover
  decision becomes real

## What unification means here

In `bolabaden-infra`, unification does **not** mean:

- invent one magical control surface
- flatten everything into a giant YAML or giant Go program
- assume Kubernetes, Nomad, OpenSVC, or Constellation automatically solves the
  whole stack

It means:

> reduce the amount of hidden operator reconstruction required to understand
> what runs where, how requests are routed, what fails over, and what only
> appears healthy.

That is the real job.

And that job is stricter than "make the platform feel more integrated."

This repo is full of things that can feel more integrated while still leaving
the same hidden burden in place.

So a unification move only counts if it changes where the burden actually
lives, not just how pleasantly the stack can be described afterward.

That also explains why the repo can feel starved for options while being
surrounded by tools.

The missing thing is not product availability.
It is availability of options that actually relocate truth out of private
operator memory.

Without that relocation, "unification" can become another polished non-answer:

- cleaner diagrams
- more respectable product nouns
- more platform mood
- the same underlying reconstruction tax

## The dream behind the word

The user is clearly asking for a middle layer that feels more coherent than raw
multi-node Compose, but less doctrinaire than immediately adopting a heavyweight
orchestrator.

The dream is roughly:

- any node can receive public traffic
- if the service exists locally, serve it locally
- if the service lives elsewhere, forward to the peer that hosts it
- keep routing, middleware, and auth semantics coherent across that handoff
- stop depending on private operator memory for placement truth
- stop pretending stateful HA is solved just because a web request still loads

That is the unification target.

There is a sharper reading of the same target:

the user is not asking for a single platform identity.
They are asking for fewer moments where the platform stops making sense unless
one human silently supplies the missing context.

That is the deeper demand behind the frustration with the lack of options.

The repo is not only asking "what can we adopt?"
It is asking:

- what here would actually count as a real option instead of fake-option
  theater?

That is why this file should keep focusing on truth classes instead of tool
categories.

## The five truths that actually need unifying

### 1. Convergence truth

Do the right nodes have the right:

- compose fragments
- secret material
- env values
- generated route data
- support files

If not, failover remains a performance.

### 2. Placement truth

Can the system answer:

- what runs where right now
- what host should answer this request
- what is eligible to take over

If not, the operator is still the hidden scheduler.

### 3. Ingress truth

Can any node receive a request and decide honestly whether to:

- serve locally
- forward to a peer
- fail because the service is actually unavailable

If not, the "any-node" story is still partly imaginary.

### 4. Failover truth

When a backend or node disappears:

- does the route survive
- is there a real next candidate
- does the system converge after handoff

If not, the repo still has decorative HA.

### 5. State truth

When writes and durable topology matter:

- who is primary
- who can safely take over
- what replication exists
- what reconnect behavior is expected

If this is fuzzy, the system may look resilient at L7 while still being
architecturally fragile.

## Why the repo keeps exploring multiple futures

Different coordination layers help with different missing truths.

### Compose plus helpers

Good for:

- readable authoring
- low migration shock
- small targeted improvements

Weak when:

- placement truth remains implicit
- route survival depends on brittle generators
- stateful behavior needs stronger coordination

It is also weak when the helper layer becomes fluent enough to sound like a
control plane before it has actually accepted control-plane accountability.

That is one of the biggest danger patterns in this repo:

- the helper layer sounds coherent
- the docs start narrating coherence as capability
- the operator still privately supplies the missing truth anyway

### OpenSVC-style HA

Good for:

- infra-grade failover
- service takeover semantics
- stronger continuity at the node/service boundary

Still not a free answer for:

- developer workflow clarity
- full placement ergonomics
- stateful correctness across every backend

### Nomad

Good for:

- scheduler-backed placement truth
- relocation semantics
- clearer runtime state

Cost:

- stronger control-plane worldview
- another platform to operate

### k3s / Kubernetes

Good for:

- rich control-plane semantics
- workload scheduling
- broad ecosystem support

Cost:

- a much heavier operational and storage model
- easy over-adoption when only part of the stack truly needs it

### Constellation / custom agent layers

Good for:

- preserving repo-specific behavior
- experimenting with a Compose-adjacent control plane
- owning exactly the semantics the repo cares about

Risk:

- reinventing a control plane without enough proof
- drifting into "claimed complete" territory before runtime validation exists

There is an even more repo-specific version of that risk:

- the custom layer feels psychologically perfect because it preserves local
  ownership language
- but the actual runtime truth is still too weak, too scattered, or too
  unverified to justify the confidence it encourages

That is why custom unification work must be judged even more harshly than
off-the-shelf platform adoption.
It is much easier for homegrown layers to inherit the user's preferred
language while still failing the user's real benchmark.

That is also why this file should not flatter custom work merely for sounding
closer to the repo's emotional vocabulary.

Local language fit is not proof.
Burden movement is proof.

## Current honest reading

The current repo already proves:

- Compose is still the primary human contract
- the user wants any-node entry with local-first then peer-forward behavior
- the repo no longer trusts raw Compose alone to express all runtime truth

The current repo does **not** yet prove:

- universally trusted placement truth
- broadly validated peer-aware failover
- stateful anti-SPOF behavior across all critical backends
- one settled control-plane choice

That means this file should never be used as evidence that the repo has already
become "one system" in the stronger sense the user actually wants.

It should be used as evidence that the repo understands the shape of the
remaining fragmentation more clearly than many surrounding tools do.

That difference matters.

Understanding the fragmentation clearly is a prerequisite for finding a real
option.
It is not the same thing as already having one.

## Practical plan

The most useful unification work remains:

1. Make placement truth explicit.
2. Make route generation depend on that truth instead of folklore.
3. Reduce bootstrap and secret-convergence drift.
4. Separate ingress resilience from stateful resilience in both code and docs.
5. Promote only the semantics that have earned a stronger platform.

## Where the fuller repo-grounded version lives

The rendered knowledgebase version of this topic is here:

- [`knowledgebase/research/plan-infrastructure-unification.md`](knowledgebase/research/plan-infrastructure-unification.md)

That page is the better maintained explanatory surface.
This root file now exists to keep a concise, honest, operator-usable statement
near the repo root.
