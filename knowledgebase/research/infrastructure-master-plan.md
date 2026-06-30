# Infrastructure Master Plan: What It Actually Proves

This page is not a replacement for
[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md).
It exists because the master plan is one of the easiest files in the repo to
overread.

The master plan is detailed enough to sound half-implemented.
That is dangerous.

It is dangerous in a very specific way:

- a module is named
- a flow is described
- the description sounds coherent
- coherence gets mistaken for system-owned truth

That is exactly the kind of documentation drift the rest of this knowledgebase
is trying to stop.

It is also one of the easiest places to accidentally downgrade the user's real
frustration into ordinary roadmap language.
If this page turns into:

- "here is the future platform"
- "here are the modules"
- "here is the migration path"

then it has already lost the actual wound.

The wound is that local Docker starts legible, multi-node life destroys that
legibility, and most ecosystems answer by demanding a bigger worldview before
the repo has named which missing truth actually justified the surrender.

The actual question readers need answered first is:

> what does the master plan prove about the user's dream and the repo's future
> direction, and what does it still not prove in the live tracked runtime?

That distinction matters because this project is already surrounded by
infrastructure language that habitually upgrades good plans into implied
capabilities.
The docs here have to do the opposite.

They also have to do something slightly harsher than ordinary roadmap hygiene:

they have to keep saying whether the plan merely describes missing truth well,
or whether the runtime already owns any of that truth directly.

## The shortest honest reading

The master plan proves that the repo has already stopped treating its problem
as:

- "we need more services"
- "we need nicer proxy config"
- "we need a cleaner deployment story"

Instead, it proves that the repo understands the real missing layer is about
distributed truth:

- placement truth
- convergence truth
- peer-aware routing truth
- failure-preserving route truth
- restart and repair truth
- bootstrap portability

It does **not** prove those layers are live, battle-tested, or final.

It also does **not** prove that the hidden operator reconstruction burden has
already materially dropped just because the future layer is now legible on
paper.

That sentence is the main rule for reading the whole master plan.

## Why this page needs to exist at all

The user's actual complaint is not just that docs have been too simple.
It is that too many systems and explanations feel sophisticated while avoiding
the real question:

> what actually removes the hidden operator burden and wrong-node fragility,
> instead of just renaming them?

The master plan is one of the clearest places where that question gets asked
seriously.
It is also one of the easiest places for a reader to quietly hear:

- "a named agent must already exist"
- "a detailed failover flow must already be validated"
- "a coherent architecture must already be operational"

This page exists to block those mistakes.
It also exists because the repo now contains enough serious planning language
to perform coherence.
That means the docs have to keep asking a harsher question:

> what exact hidden human reconstruction burden is this planned module trying
> to remove, and what burden is it still quietly leaving behind?

That second half matters just as much as the first.

A plan can be excellent and still leave behind the one burden the user actually
cares about most:

the system still only works honestly because a human privately knows which node,
which route, or which authority is real.

## What the master plan clearly proves

## 1. The repo is deliberately searching for a thinner control layer before promoting a heavyweight orchestrator

The master plan repeatedly returns to the same structure:

- keep Docker and Compose as the human authoring layer
- add targeted sync, convergence, and failover logic where Compose alone is too
  static
- pay only for the coordination semantics that remove real operator pain

That is not procrastination.
It is architecture intent under pressure.
It is the repo trying not to lose the last readable authoring surface before it
has correctly identified which missing truths actually justify replacing it.

This matters because the user is not just skeptical of Kubernetes or Nomad as
brands.
They are skeptical of ecosystems that demand a giant worldview shift before
earning the right to solve the actual local problem.

### What this proves

- the repo is explicitly searching for a missing middle layer
- it does not accept "just use Kubernetes" as a sufficient answer
- orchestration promotion is supposed to be earned by demonstrated need, not by
  fashion

### What this does not prove

- that the thinner layer will be sufficient forever
- that the thinner layer is already implemented
- that the thinner layer will remain thinner once it owns enough truth

## 2. The repo sees cross-node truth, not extra containers, as the dominant missing capability

This is one of the strongest signals in the entire master plan.

The plan is not mainly:

- add more components
- add more sidecars
- add more proxy tricks

It is mainly:

- stop leaving cluster truth in operator memory

The recurring missing duties are:

- env and secret convergence
- compose or file convergence
- service placement awareness
- restart or redeploy semantics when truth changes
- failover logic that follows current reality instead of stale assumptions

That is the real center of gravity.
The repo is not starving for "more components."
It is starving for a system that stops forcing the operator to privately stitch
together:

- who is authoritative
- what the wrong node should do
- which peer is current
- whether the route that exists now will still exist after failure

### What this proves

- the user's dominant pain is hidden coordination truth
- the project already understands that human memory is currently one of the
  biggest SPOFs

### What this does not prove

- that a sync-agent style layer is already trustworthy enough to own those
  responsibilities

## 3. The plan treats request preservation as central, not decorative

The master plan spends real energy on:

- node-entry behavior
- global versus node-specific hostnames
- DDNS and route generation
- internal peer connectivity
- recovery and repair assumptions

That matters because the user is not mainly asking:

> how do I deploy services?

They are asking:

> how do I keep the request alive, honest, and semantically coherent when it
> lands on the wrong node or when the expected backend disappears?

That is a much sharper requirement than "high availability."
It is closer to:

> preserve the meaning of the request without requiring the operator to become
> the hidden routing oracle

### What this proves

- any future control layer that ignores wrong-node request preservation is
  misaligned
- routing truth is part of the core problem, not a polish layer

### What this does not prove

- that current failover generation already behaves correctly under failure

## 4. The repo wants a portable personal-cloud system, not a private ritual stack

The master plan keeps pushing toward:

- cleaner bootstrap
- lower env and secret drift
- domain abstraction
- more reusable infra shape

That reveals another important part of the dream:

the system should not only work for the original operator who already remembers
every invisible fact.

This matters because one of the user's deepest frustrations is hidden knowledge.
A system that "works" only because one person privately remembers everything is
already centralized in the worst possible place: a human head.
That is not just a portability complaint.
It is one of the repo's deepest anti-SPOF arguments.

It is also why this page cannot treat "portable architecture" and "portable
truth" as synonyms.

The repo can describe a portable architecture long before it proves that the
portable truth layer is strong enough to stop social reconstruction from being
the real control plane.

### What this proves

- documentation and control surfaces must survive portability pressure
- "works for the author because the author remembers the ghosts" is not an
  acceptable end state

### What this does not prove

- that portability or reproducibility are already solved

## What the master plan does not prove

The master plan is strong planning evidence.
It does **not** prove that any of the following are already true in the tracked
runtime:

- multi-node env and secret sync works live
- compose and file convergence is already trustworthy
- a live tracked root `services.yaml` exists and is consumed
- peer-aware forwarding is broadly validated
- failover routes survive the exact failure that made them necessary
- sync-agent or failover-agent loops are implemented and battle-tested
- watchtower-style repair semantics are dependable
- stateful services already have anti-SPOF correctness

This list should stay painful.
If it becomes easy to forget, the plan will start impersonating runtime truth.

It will also start impersonating reduced operator burden.

Those are different failures, but they usually travel together in this repo.

## What the master plan quietly assumes

The plan is coherent, but it depends on several very large assumptions.
Those assumptions have to stay visible rather than being smuggled into
"future architecture."
That is the discipline this page should model:

- do not erase contradiction to make the plan cleaner
- do not erase uncertainty to make the plan sound stronger
- do not erase partial truth to make the repo sound more finished

## Assumption 1: a lightweight agent layer can deliver a meaningful share of orchestrator value

For that to be true, the future layer must reliably answer:

- what service belongs where
- what changed
- what must restart
- what must remain untouched
- which peer is healthy and eligible enough to receive forwarded traffic

That is already dangerously close to building a control plane around Compose.

This is not criticism.
It is the actual cost model.
The project is not avoiding orchestration logic entirely.
It is trying to buy only the fragments that correspond to the exact wound.
If those fragments eventually add up to "you built a control plane anyway,"
the docs need to say that plainly.

And if those fragments still leave the operator privately finishing the most
important decisions, the docs need to say that plainly too.

The repo's real question is not just whether a thin layer stays thin.
It is whether the layer actually takes ownership of the truths it claims to
purchase.

### What this proves

- the repo is not trying to avoid orchestration logic entirely
- it is trying to decompose orchestration logic into narrower, more legible
  parts

### What this does not prove

- that the decomposed form will remain cheaper or clearer than promoting some
  domains into a stronger scheduler-backed platform

## Assumption 2: ingress truth can be generated from runtime truth without becoming brittle again

For that to work, generated routing needs:

- stable placement input
- stable health input
- route persistence under backend disappearance
- middleware and auth continuity across peer handoff

That is a much harder requirement than "render a config file."

### What this proves

- the repo correctly treats route generation as a core answer, not a side task

### What this does not prove

- that the current generation path is already safe

## Assumption 3: portability can coexist with operator readability

The plan wants:

- reproducible bootstrap
- forkable infra
- explicit control surfaces
- bounded worldview tax

Those goals are all good.
They also pull against each other.

The more portable and automated the system becomes, the easier it is to hide
important truth behind the portability layer itself.

### What this proves

- portability here is not just a convenience feature
- it is part of the structural pressure shaping the control layer

### What this does not prove

- that the tradeoffs are already resolved

## The highest-signal modules in the plan

The best way to read the master plan is not by module names alone.
It is by asking:

> what hidden burden is each proposed module trying to buy back from the
> operator?

## Secret and env sync

This is not polish.
It is an attempt to remove one of the quietest cross-node SPOFs:

- the forwarded request reaches a node that does not actually share the same
  env, secret, or auth assumptions

### What this proves

- the repo already understands convergence truth is central

### What this does not prove

- that the nodes are already equivalent enough to support trusted fallback

## Compose and file sync

This is the same pressure in another form:

- a distributed stack is only as trustworthy as the degree to which each node
  is running the same intended world

### What this proves

- the repo does not accept manual drift as a harmless inconvenience

### What this does not prove

- that sync semantics are already implemented safely

## `services.yaml` and placement truth

This is arguably the most important conceptual module in the entire plan.

It is the repo's repeated answer to:

> where does a receiving node learn what is local and what is remote right now?

Without an answer to that question, wrong-node request preservation remains a
wish.

### What this proves

- the repo understands that placement truth is the first real missing layer

### What this does not prove

- that the live root runtime already ships and consumes that truth surface

## Failover-agent and routing repair logic

This module exists because the user is tired of systems that look dynamic until
the local backend actually dies.

The plan is clearly trying to buy:

- failure-aware regeneration
- durable fallback paths
- recovery logic that survives the event that triggered it

### What this proves

- the repo is targeting actual failure absorption, not just pretty routing
  diagrams

### What this does not prove

- that the current fallback path survives backend-loss events in live runtime

## Bootstrap protocol

This module exists because the operator should not be the only place where the
system still makes sense.

That is not just convenience.
It is an anti-hidden-knowledge requirement.

### What this proves

- the project wants less private ritual and more explicit reproducibility

### What this does not prove

- that bootstrap portability is already mature

## What this page means for the rest of the knowledgebase

Whenever another page cites the master plan, it should do so carefully:

- cite it as strong future-direction evidence
- cite it as proof that the user is searching for a missing middle layer
- cite it as evidence that cross-node truth is the real problem
- do not cite it as proof that the mechanisms are live

That last rule matters the most.

The master plan is valuable precisely because it sees the problem clearly.
It becomes dangerous when clarity of diagnosis gets mistaken for proof of cure.

## The real takeaway

The master plan proves that the repo is not confused about what hurts.

It knows the pain is not "more YAML" or "more containers."
It knows the pain is:

- wrong-node fragility
- hidden placement truth
- hidden convergence truth
- failure paths that disappear at the moment they matter
- operator knowledge that still acts as the real cluster manager

What the master plan does not prove is that those pains are already solved.

That is exactly why it is one of the most important documents in the tree.
It is valuable because it names the missing layer honestly.
It should not be allowed to impersonate the presence of that layer before the
runtime has earned it.
