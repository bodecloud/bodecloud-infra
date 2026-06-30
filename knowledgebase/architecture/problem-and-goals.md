# Problem, Pressure, and Goals

This page defines the actual architecture problem in `bolabaden-infra`.

The problem is not:

- "more clustering"
- "better deployment hygiene"
- "more mature infrastructure"
- "better docs about self-hosting"
- "choosing the right orchestrator"

The problem is:

> how do several ordinary Docker nodes become one request-preserving personal
> cloud while
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> stays readable and the system stops leaning on one operator to remember where
> the real topology truth lives?

That question is stricter than generic HA language and smaller than "migrate to
Kubernetes."
It is also much closer to the user's actual frustration than most infrastructure
summaries are willing to be.

The repo is not just looking for more resilience.
It is looking for a system that stops humiliating the operator by revealing,
too late, that the operator was still the real keeper of placement, failover,
and route-meaning truth.

## What this page is and is not allowed to prove

This page is authoritative about:

- the real problem the repo is trying to solve
- the concrete requirement stack implied by that problem
- which neighboring answers are still too small
- why the no-Swarm dream is not reducible to generic HA language

This page is not authoritative about:

- whether the current runtime already satisfies those requirements
- whether one proposed helper layer has already won
- whether a sharper problem statement makes the remaining gap small

This page is the benchmark, not the completion report.

## The user's real demand

The repo's strongest intent surfaces,
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md),
[README.md](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md),
and
[knowledgebase/AGENTS.md](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/AGENTS.md),
make the pressure explicit:

- do not default immediately to Swarm, Kubernetes, k3s, or another heavyweight
  control plane
- keep Compose as the real authoring and operator surface
- make any healthy public node a plausible first hop
- serve locally when locality exists
- forward to a healthy peer when locality does not exist
- do not fake HA by confusing reachability, DNS plurality, or proxy presence
  with genuine burden transfer

That is why the real question is not "what orchestrator should I use?"

It is closer to:

> can multiple ordinary Docker nodes become one believable platform without the
> final truth still living in one human head?

That is the dream this knowledgebase has to keep visible.

## The exact operating contract the dream is pointing at

The most useful sentence in the repo is not "multi-node" and it is not
"anti-SPOF."

It is the request contract preserved in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That is what the repo keeps trying to achieve without immediately surrendering
to:

- Docker Swarm
- Kubernetes or k3s
- a giant scheduler that owns more truth than it can explain

This matters because a lot of otherwise intelligent infrastructure writing
quietly swaps that contract for smaller substitutes such as:

- "multiple public nodes exist"
- "Traefik is present on the edge"
- "the route can probably be recreated elsewhere"
- "there is some kind of failover helper"

Those are pieces of the environment.
They are not the contract.

## What "no Swarm by default" really means here

The no-Swarm or no-Kubernetes pressure should not be read as ideology.
It is a burden-accounting rule.

The repo is effectively asking:

> how much shared truth can we add before we have to promote ourselves into a
> heavier control plane?

That is a different question from:

> which orchestrator has the most features?

The user is not angry because orchestrators exist.
The user is angry because so many proposed answers only feel impressive while
one operator still privately knows:

- what runs where
- which peer is valid
- which route still means the same thing after handoff
- which fallback disappears under the exact failure that made it matter

That is why "stay Compose-first" is not nostalgia.
It is a demand to justify every promoted control layer by the hidden burden it
actually removes.

## The shortest exact problem statement

The repo is trying to do this:

> keep Compose as the main authoring and operator surface, but add just enough
> shared truth that a request landing on the wrong healthy node does not turn
> into guesswork, folklore, redirects, or fake failover.

The key phrase is "shared truth."
The repo is not starved for components.
It is starved for smaller honest control surfaces that actually move truth out
of private operator memory and into the system.

## What the system must eventually know for itself

If the dream becomes real, the system has to own more than "containers are up."
It has to own decisive bad-day truths such as:

- where the requested service actually lives now
- whether the receiving node should serve locally or forward
- which peer is healthy and eligible for that specific route
- whether auth, middleware, and headers survive the handoff unchanged
- whether fallback still exists under the failure that made it necessary
- whether the route class is stateless HTTP, protected HTTP, raw TCP, or a
  stateful surface that needs stricter semantics

The user is not merely asking for these truths to be writable somewhere.
They are asking for them to become inspectable enough that the wrong receiving
node does not need a human narrator to finish the story.

If those truths still cash out into:

- "well, privately we know where it runs"
- "in practice the operator knows which peer is safe"
- "normally that hostname really belongs to node X"

then the central problem is still alive, even if the surrounding stack looks
much more serious.

## Why common answers still feel too small

Many nearby answers improve one layer while quietly leaving the decisive burden
intact.

Examples:

- DNS plurality helps more than one public node receive traffic, but it does
  not prove the wrong node can preserve the request meaningfully
- Traefik helps with HTTP routing, but its presence alone does not prove
  peer-forward continuity or stateful correctness
- healthchecks improve local truth, but they do not by themselves define
  peer eligibility
- file sync, secret sync, and helper generation can reduce drift, but they do
  not automatically create trustworthy current placement truth
- Swarm, Nomad, OpenSVC, k3s, or Kubernetes may eventually earn a place, but
  only if they remove a concrete hidden burden rather than merely replacing one
  kind of opacity with another

That is why the user's frustration is not just "there are too many options."
It is:

> too many options solve one visible layer and then quietly leave the operator
> as the hidden control plane when reality gets sharp.

## The wound behind the problem

This page should preserve the lived failure scene, not smooth it away.

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

## What this repo is actually trying to make impossible

The repo is trying to make this scene stop being normal:

1. a request hits a healthy public node
2. the requested service is not local to that node
3. the operator still has to privately know where the real backend lives
4. the operator still has to privately know whether forwarding is safe
5. the operator still has to privately know whether auth, middleware, or data
   semantics survive the handoff

That is the hidden work the platform is supposed to absorb.

If one more proposed layer does not make that scene less dependent on private
memory, then it is not yet attacking the central pain.

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

The acceptance tests for future architecture work therefore sound like:

- can the receiving node determine locality or remote ownership from shared
  truth rather than memory?
- can it choose a peer from eligibility truth rather than reachability alone?
- can it preserve the same protected route semantics after forwarding?
- can the fallback still work after the preferred backend disappears?
- can stronger claims stay honest under direct inspection of the runtime?

Another way to say the same thing:

- can the node explain its distributed decision from shared artifacts?
- can the operator inspect the explanation after the fact?
- can the explanation survive without socially reconstructing the topology?

If the answer to those is still "not yet," then the remaining gap is not
cosmetic.
It is still the central unsolved fact.

## Bottom line

The repo is not just building infrastructure.
It is searching for the smallest honest middle layer between:

- static multi-node Docker glue that still depends on private folklore
- and a heavyweight orchestrator worldview that has not yet earned the right to
  hide that much truth

That is the benchmark every later page should inherit.

Better understanding of the problem is necessary.
It is still not architecture progress by itself.
Architecture progress only happens when one more decisive bad-day sentence stops
living in private operator memory.
