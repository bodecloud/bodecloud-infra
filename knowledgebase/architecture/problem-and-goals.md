# Problem, Pressure, and Goals

This page defines the actual architecture problem in `bolabaden-infra`.

It is not:

- better deployment hygiene
- more modern infrastructure
- more cluster-shaped diagrams
- picking the "right" orchestrator
- writing nicer docs about self-hosting

It is:

> how do several ordinary Docker nodes become one believable request-preserving
> personal cloud while
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> remains readable and one operator stops being the hidden registry, hidden
> peer selector, hidden route explainer, and hidden memory of where the real
> topology truth lives?

That sentence is harsher than ordinary HA prose on purpose.
The user is not merely frustrated with outages.
The user is frustrated with fake options that still collapse into private human
glue once reality gets sharp.

## What this page is and is not allowed to prove

This page is allowed to prove:

- the real problem the repo is trying to solve
- the hidden burdens implied by that problem
- why many neighboring answers are too small
- why "no Swarm by default" is really a burden-accounting rule

This page is not allowed to prove:

- that the current runtime already satisfies the problem
- that one future helper layer has already won
- that a sharper benchmark makes the remaining gap small

This page is the benchmark.
It is not the completion report.

## The user's real demand

The strongest intent surfaces in the repo make the pressure explicit:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [knowledgebase/AGENTS.md](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/AGENTS.md)

Together they imply:

- do not default immediately to Swarm, Kubernetes, k3s, or another heavyweight
  scheduler worldview
- keep Compose as the real authoring and operator surface
- let any healthy public node be a plausible first hop
- serve locally when locality exists
- forward to a healthy peer when locality does not exist
- do not counterfeit HA by confusing reachability, DNS plurality, or proxy
  presence with transferred burden

That is why the real question is not:

> what orchestrator should I use?

It is closer to:

> can multiple ordinary Docker nodes become one believable platform without the
> final truth still living in one human head?

That is the dream these docs have to keep visible.

## The exact operating contract the dream points at

The most useful sentence in the repo is not "multi-node."
It is not even "anti-SPOF."

It is the contract preserved in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That is the thing the repo keeps trying to achieve without immediately
surrendering to:

- Docker Swarm
- Kubernetes or k3s
- a giant scheduler that owns more truth than it can explain

Many otherwise intelligent summaries quietly replace that contract with smaller
substitutes such as:

- multiple public nodes exist
- Traefik is present on the edge
- the route can probably be recreated elsewhere
- some kind of failover helper exists

Those are environment facts.
They are not the contract.

## What "no Swarm by default" really means

The no-Swarm or no-Kubernetes pressure is not ideology.
It is burden accounting.

The repo is really asking:

> how much shared truth can we add before we have to promote ourselves into a
> heavier control plane?

That is a different question from:

> which orchestrator has the most features?

The user is not angry because orchestrators exist.
The user is angry because too many proposed answers only feel impressive while
one operator still privately knows:

- what runs where
- which peer is valid
- which route still means the same thing after handoff
- which fallback disappears under the exact failure that made it matter

So "stay Compose-first" is not nostalgia.
It is a demand to justify every promoted control layer by the hidden sentence
it actually kills.

## The shortest exact problem statement

The repo is trying to do this:

> keep Compose as the main authoring and operator surface, but add just enough
> shared truth that a request landing on the wrong healthy node does not turn
> into guesswork, folklore, redirects, or fake failover.

The key phrase is `shared truth`.

The repo is not short on components.
It is short on smaller honest control surfaces that move decisive bad-day truth
out of private operator memory and into the system.

## What the system must eventually know for itself

If the dream becomes real, the system has to own more than "containers are up."
It has to own bad-day truths such as:

- where the requested service actually lives now
- whether the receiving node should serve locally or forward
- which peer is healthy and eligible for that specific route
- whether auth, middleware, and headers survive handoff unchanged
- whether fallback still exists under the failure that made it necessary
- whether the route class is stateless HTTP, protected HTTP, raw TCP, or a
  stateful surface requiring stricter semantics

The user is not merely asking for those truths to be writable somewhere.
They are asking for them to become inspectable enough that the wrong receiving
node does not need a human narrator to finish the story.

If those truths still cash out into:

- "privately we know where it runs"
- "in practice the operator knows which peer is safe"
- "normally that hostname really belongs to node X"

then the central problem is still alive, no matter how serious the surrounding
stack looks.

## Why common answers still feel too small

Many nearby answers improve one layer while quietly leaving the decisive burden
intact.

Examples:

- DNS plurality lets more than one public node receive traffic, but it does not
  prove the wrong node can preserve the request meaningfully
- Traefik helps with HTTP routing, but its presence alone does not prove
  peer-forward continuity or stateful correctness
- healthchecks improve local truth, but they do not by themselves define
  peer eligibility
- sync loops can reduce drift, but they do not automatically create trustworthy
  current placement truth
- Swarm, Nomad, OpenSVC, k3s, or Kubernetes may eventually earn a place, but
  only if they remove a concrete hidden burden rather than replacing one kind
  of opacity with another

That is why the user's frustration is not just "there are too many options."
It is:

> too many options solve one visible layer and then quietly leave the operator
> as the hidden control plane when reality gets sharp.

## The wound behind the problem

This page should preserve the lived failure scene instead of smoothing it away.

The user keeps hitting the same reveal:

1. the stack sounds flexible
2. the stack sounds distributed
3. the stack sounds full of serious options
4. a real request or failure arrives
5. the decisive truth still lives in one human head

That is why phrases like:

- operational complexity
- topology awareness
- coordination overhead

are too polite if they replace the sharper truth:

> when reality gets sharp, the operator is still the hidden control plane.

That sentence is not melodrama.
It is the most faithful summary of the repo's actual pain.

## What the repo is trying to make impossible

The repo is trying to make this scene stop being normal:

1. a request hits a healthy public node
2. the requested service is not local to that node
3. the operator still has to privately know where the real backend lives
4. the operator still has to privately know whether forwarding is safe
5. the operator still has to privately know whether auth, middleware, or data
   semantics survive the handoff

That is the hidden work the platform is supposed to absorb.

If one more proposed layer does not make that scene less dependent on private
memory, it is not attacking the central pain yet.

## Acceptance tests implied by the problem

The benchmark implied by the user's dream is harsher than "does the stack look
more professional?"

A believable step forward should narrow at least one of these:

- remembered placement
- remembered safe peer choice
- remembered route meaning
- remembered fallback durability
- remembered stateful caveats

In other words, a candidate layer matters only if it changes who owns the truth
on the bad day.

That means future architecture work should be judged by questions like:

- can the receiving node determine locality or remote ownership from shared
  truth rather than memory?
- can it choose a peer from eligibility truth rather than reachability alone?
- can it preserve the same protected-route semantics after forwarding?
- can the fallback survive the exact failure that made it relevant?
- can the docs still say what remains forbidden after the improvement lands?

## Goals, stated brutally

The goals are not:

- sound clustered
- look modern
- get closer to industry standard
- use bigger platform nouns sooner

The goals are:

1. make the healthy wrong node less dependent on one human's head
2. move placement, eligibility, and fallback truth into inspectable system
   state
3. preserve the meaning of the request, not just packet delivery
4. keep proof rules strict enough that stateful, TCP, and protected-route
   claims do not borrow confidence from simpler lanes
5. justify any promoted control layer by the exact private sentence it removes

## Bottom line

The user's dream is not "multi-node infrastructure" in the abstract.

It is:

> a believable request-preserving personal cloud made out of ordinary Docker
> nodes, where Compose remains inspectable and one human stops being the place
> where the final answer still secretly lives.

If a page, tool, or architecture move loses that sentence, it is already
drifting away from the point of the repo.
