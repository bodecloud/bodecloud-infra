# Problem, Pressure, and Goals

This page exists to stop the repo from being described with weak verbs, lazy
questions, and infrastructure language that sounds smarter than the actual
understanding behind it.

The problem here is not:

- "better deployment"
- "more clustering"
- "more HA features"
- "a more modern platform"

The problem is much harsher:

the user wants several ordinary Docker machines to behave more like one
resilient, peer-aware personal cloud at request time, while refusing to accept
a heavyweight control plane as the automatic price of adulthood.

That means the real architecture question is not:

> which orchestrator should we use?

It is:

> what exact truths and recovery behaviors are required for multiple
> Compose-managed nodes to preserve requests, preserve operator clarity, and
> reduce SPOFs without importing more platform tax than the pain actually
> justifies?

That question also implies a documentation requirement:

the page must not make the problem sound cleaner than it is by pretending the
relevant truths all mature together.

They do not.
That unevenness is part of the architecture, not an editorial annoyance.

If a page does not answer that version, it is still too shallow.

The page also has to do something stronger than sounding intense.
It has to reconstruct the problem in a way that lets later pages be judged.

That means this page should answer four separate questions at once:

1. what exact pain is being treated as central
2. which truths the current stack still fails to own cleanly
3. what a real improvement would need to relocate out of private memory
4. which adjacent infrastructure answers are still too small even if they are
   technically respectable

The docs can also fail by answering a more normal adjacent question extremely
well.

They can become impressive by explaining ingress, clustering, service
discovery, or orchestrator options while still failing to reconstruct the thing
the user is actually pushing toward:

> a distributed personal cloud whose intelligence is no longer privately held
> together by one operator's memory

For the deeper intent behind this page, start with
[`../research/user-intent-and-dream.md`](../research/user-intent-and-dream.md)
and [`../research/archive-pressure-patterns.md`](../research/archive-pressure-patterns.md).
For the repo's instruction-surface hierarchy, read
[`instruction-surfaces-and-authority.md`](instruction-surfaces-and-authority.md).

## What this page is and is not allowed to prove

This page is authoritative about:

- what problem the repo is actually trying to solve
- which hidden burdens define success more accurately than generic HA language
- what the dream requires before it counts as real progress

This page is not authoritative about:

- whether the current runtime already satisfies those requirements
- whether one future control-plane path has already won
- whether naming the problem cleanly means the implementation gap is small

This page is the benchmark and requirement frame.
It is not a claim that the benchmark has been met.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "this is the actual problem, not an adjacent one" | problem reconstruction | this page, dream pages, archive pressure | that the runtime already answers it |
| "the dream implies these requirements" | benchmark synthesis | this page plus dream and evidence pages | that all requirements mature together |
| "generic options lists are still too small" | negative-benchmark synthesis | this page, archive-pressure pages, orchestrator evidence | that no partial option can still help |
| "operator memory is still the hidden control plane" | requirement-pressure judgment | this page, runbook, proof matrix, runtime pages | that nothing meaningful has been externalized yet |

If a sentence starts using a clear requirement frame as if it were proof of
delivery, it has outrun this page's authority.

## The actual problem statement

The user is not primarily trying to solve:

- how to run Docker on more than one VPS
- how to expose services through Traefik
- how to put Cloudflare in front
- how to choose an orchestrator brand

Those are subordinate questions.

The real problem is:

> how do I make multiple Docker nodes behave like one request-preserving
> system, where traffic can land anywhere and still find the right service,
> without immediately surrendering the whole problem to Kubernetes, Swarm, or
> some other desired-state empire?

That is not just a hosting question.
It is a truth-ownership question.

It includes:

- placement truth
- convergence truth
- routing truth
- failure truth
- policy continuity truth
- stateful honesty

Those are not seven equal bullets.
They are seven different truth burdens that tend to break at different times.

That asymmetry is the whole reason the repo is so easy to misunderstand.

Some layers can look strong while the more humiliating ones remain weak.
That is why this page has to behave more like a requirement breakdown than a
generic architecture intro.

What makes this problem hard is that each of those truths becomes visible at a
different moment of pain.

That is why the repo can look more mature than it feels.

Some truths are already strong:

- authoring truth
- local service truth
- edge presence

Some truths are still weak:

- wrong-node recovery truth
- cross-node semantic parity
- fallback-route survival under the relevant failure

If those are narrated as if they mature together, the docs lie.

## Requirement stack implied by the dream

Once the dream is stated honestly, the required stack becomes more concrete.

The repo is implicitly asking for a platform that can satisfy all of these:

1. any-node first hop:
   any surviving public node can accept the request
2. local-first preservation:
   if the requested service is local, the node should serve locally without
   pretending that locality does not matter
3. wrong-node preservation:
   if the service is remote, the request still has a legitimate path to the
   correct healthy peer
4. policy preservation:
   the forwarded request keeps the same auth, middleware, and routing meaning
   instead of becoming a degraded bypass path
5. survivable fallback:
   the rescue path remains available under the failure it was meant to absorb
6. truthful state language:
   stateful systems are described according to their real write, storage, and
   failover semantics rather than by the mere existence of a network route
7. inspectable ownership:
   an operator can explain why the system behaved as it did from tracked shared
   truth instead of remembered topology folklore

This requirement stack is the page's most important output.

If a future mechanism satisfies:

- `1` and `2`, but not `3`
- or `3` and `4`, but not `5`
- or `1` through `5`, but not `6`
- or all of those while still failing `7`

then it may represent progress, but it has not yet reached the repo's actual
goal state.

That is why "does failover work?" is still too weak a question.
The better question is:

which parts of this stack are genuinely system-owned now, and which parts are
still being socially reconstructed by the operator?

That is why so many options feel promising for ten minutes and then collapse:

- DNS plurality answers first-hop reachability
- reverse proxies answer some local routing questions
- healthchecks answer some local liveness questions
- schedulers answer some larger placement questions

But the user is not looking for individually impressive answers.
The user is looking for the smallest stack shape that stops producing the same
psychological failure:

> I added more infrastructure, but I still have to remember the real topology
> privately for the system to make sense.

That is one of the most important sentences in the whole architecture set.

It explains why "more infrastructure" can feel like negative progress here.
If more components arrive without relocating where the actual truth lives, the
system becomes more decorated but not more honest.

## The hidden enemy is not lack of tools

The archive sometimes sounds like the complaint is "there are not enough
options."
That is not the deeper issue.

The deeper issue is that too many available options solve one narrow slice
while leaving the operator trapped in reconstruction work.

The user keeps ending up in the same cursed pattern:

- one layer knows something
- another layer assumes something
- the operator privately remembers the thing that joins them

That is the real tax.

The dream is not "more capabilities."
The dream is "stop making the operator's private memory the missing control
plane."

That is why goal language here has to stay harsher than generic self-hosting
or DevOps writing.

If a goal sounds like:

- better HA
- better routing
- better orchestration

then it is probably already a smaller neighboring question.

That is why this project feels so much bigger than "Docker failover."

Another way to say it:

- manual service placement is acceptable
- explicit node affinity is acceptable
- class-specific limits are acceptable
- uneven maturity across services is acceptable

What is not acceptable is a stack that only appears coherent because the human
operator knows which parts are theater.

That word "theater" is exact.

The user is reacting against systems where redundancy exists in conversation,
routes exist in static text, peers exist on a diagram, and HA exists in
branding, while live behavioral truth still depends on who remembers what.

## The user’s hidden negative benchmark

The docs become much more honest once the hidden negative benchmark is stated
directly.

The user wants this to stop being true:

> the stack only behaves coherently because I privately remember which machine
> is special, which route is real, which peer currently owns the service, and
> which parts of the supposed HA story are actually just ceremony

That is the real enemy.

The dream is not merely distributed.
The dream is:

distributed without hidden sacred-node memory.

That sentence explains why so many superficially relevant answers still fail.

It also explains why the user is not satisfied by generic “options” lists.

## Strongest honest current answer

If someone asks, "What is the shortest honest way to state the repo's goal?"
the answer is:

> Make several ordinary Docker nodes behave like one request-preserving,
> operator-readable personal cloud by moving more of the wrong-node, fallback,
> policy, and state truth out of private operator memory and into explicit,
> inspectable system-owned behavior, without paying more control-plane tax than
> the real pain has earned.
Most option lists quietly assume one of the following is fine:

- one machine is still the real ingress node
- one machine is still the real data node
- one machine is still where the helper logic actually matters
- one human still privately remembers how to repair the illusion

Those are not neutral compromises in this repo.
They are the exact pattern the repo is trying to escape.

## The goal is not just a mechanism, but a transfer of burden

Most infrastructure goal statements name a feature set.
This one needs a different formulation.

The real goal is a transfer:

- from private memory to inspectable shared truth
- from wrong-node surprise to wrong-node preservation
- from diagram-level HA language to behavior-level honesty
- from "I know which machine is special" to "the system tells me which machine
  matters and why"

That transfer is why the user keeps rejecting otherwise respectable answers.

A platform can add:

- more routing logic
- more healthchecks
- more registries
- more automation
- or more controller machinery

and still fail the goal if the important explanations remain tacit.

This is also why "more options" keeps feeling like a non-answer.
The missing thing is not choice volume.
It is burden relocation.

If the operator remains the secret place where the system becomes coherent,
the goal has not been met.

## Why this feels bigger than just make failover work

The user is not only frustrated that failover is hard.
They are frustrated that every apparent fix reveals another hidden missing
layer.

The pressure chain is the real architecture:

1. DNS can reach more than one node
2. but the receiving node may not host the target service
3. so the node needs placement truth
4. but placement truth is not enough if peers are semantically out of sync
5. so the system also needs convergence truth
6. but convergence is not enough if the route needed for fallback disappears
   with the local backend
7. and none of that solves stateful correctness by itself

That is why the repo keeps rediscovering the same missing middle layer.

The project does not feel hard because it lacks components.
It feels hard because too many crucial truths are still implicit.

This is the place where a lot of lighter-weight documentation quietly fails.
It explains the chain, but it does not force the reader to ask which links are
currently system-owned and which are still folklore.

That distinction has to remain active on every later page.

That is why the benchmark here is harsher than “can we fail over?”

The real benchmark is:

> can we describe, inspect, and trust the path from any-node entry to correct
> service behavior without socially reconstructing the architecture every time
> a request lands on the wrong box or a backend disappears?

If the answer is still no, the docs should not pretend that another tool,
proxy, or helper has “mostly solved” the problem.

That rule should harden the entire architecture section.

The real question is not whether the stack has moved in the right direction.
It is whether the named mechanism actually reduces the need for private
operator reconstruction at the exact moment the architecture is under stress.

## The repo's desired operating model

The repo-native instruction surfaces, README, and planning layer converge on a
real operating model that should be treated as first-class.

For this architecture question:

- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

are stronger intent surfaces than
[`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md),
which is mostly about repo execution, validation, and environment handling.

That distinction matters here because one of the easiest ways to write useless
docs is to let broad agreement across repo files sound like runtime maturity.

This page should therefore be read alongside:

- [`instruction-surfaces-and-authority.md`](instruction-surfaces-and-authority.md)
- [`../operations/source-assimilation-index.md`](../operations/source-assimilation-index.md)

The first stops authority flattening.
The second stops retrieval flattening.
This page depends on both.

That distinction matters because weaker docs keep blurring:

- how to work in the repo
- what the system is trying to become

Those are not the same question.

## Compose remains the main authoring surface

The root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and the included files under
[`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
remain the real implementation anchor.

That matters because the repo is explicitly resisting a full worldview shift
before it has isolated the exact missing truths.

The desire is:

- keep service definitions readable
- keep operator intent inspectable
- keep local iteration direct
- add abstraction only when it pays down a real pain class

That last point is one of the deepest architectural rules in the repo.

The user is not trying to preserve Compose because YAML is emotionally
comforting.
The user is trying to preserve the last surface where:

- placement intent is still inspectable
- service identity is still obvious
- mistakes are still close to their cause
- the system has not yet disappeared into platform folklore

This is not stubborn attachment to YAML.
It is an attempt to preserve the last surface that still feels directly legible
before distribution starts lying.

If the repo is described as simply "preferring Compose," the reader misses the
real point.

Compose is standing in for inspectability, ownership, causal closeness, and a
refusal of abstraction that cannot name its benefit honestly.

## No central orchestrator by default

This is not supposed to quietly mutate into:

- Kubernetes-lite
- Swarm with different branding
- an unnamed scheduler hiding behind helper scripts

The strongest direct statement still lives in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md),
which frames the repo as multi-node Docker infrastructure aiming at anti-SPOF
and peer-aware behavior without immediately collapsing into Kubernetes or
Swarm.

That phrasing matters because the user is not anti-orchestrator in some
ideological sense.
They are anti-premature-surrender.

That distinction is crucial.

The repo is not claiming:

- no orchestrator should ever exist
- all control planes are bad
- all schedulers are overkill forever

The repo is claiming something narrower and much more defensible:

> if a new layer is introduced, it should earn its existence by owning a
> missing truth more honestly than the current stack does, not merely by
> sounding more mature.

## Current-state truth, not scheduler-declared truth

The repo keeps leaning toward a lightweight current-state registry such as
`services.yaml`.

That is not a cosmetic preference.
It reflects a different philosophy:

- services are manually assigned
- the system records where they actually live
- routing consumes that current truth
- operator control stays primary
- a heavyweight reconciler does not get to overwrite intent automatically

This is one of the most important differences between the user's dream and a
default cluster worldview.

But the tracked root worktree still does not ship a live
[`services.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/services.yaml).

So one of the main reading rules for this knowledgebase must remain:

- treat `services.yaml` as explicit architecture intent
- do not treat it as already-proven runtime truth

This kind of distinction is exactly what the user means by actually think.

The user is asking for documentation that can survive a skeptical reread after
the excitement fades.

That means every page has to keep asking:

- what is intention?
- what is live?
- what still depends on operator recollection?
- what has merely been reworded rather than solved?

## Any-node entry with local-first serve and peer-forward fallback

The central behavioral contract is still:

```text
User -> Cloudflare DNS -> any healthy node
  if service exists locally:
    serve locally
  else:
    forward to healthy peer that hosts it now
```

That sentence looks small.
It hides most of the hard part of the entire repo.

Because to make it real, the receiving node needs to know:

- whether the target is local
- where it actually lives if not local
- whether the remote candidate is eligible now
- whether the route survives the failure that made fallback necessary
- whether auth and middleware meaning survive the hop

And that list still understates the emotional demand beneath the architecture.

The user is not just asking for technical success.
They want the request path to stop feeling suspicious.

That means the system should eventually answer wrong-node events in a way that
feels boring, legible, and unsurprising, not clever, magical, or dependent on
remembered trivia.

That is why the user keeps rejecting shallow answers.

## The pressure chain

The repo becomes much easier to reason about when the pressure is described as
a chain instead of a slogan.

### Pressure 1: any-node entry

If more than one node can receive public traffic, first-hop survivability
improves.

That immediately creates the next question:

- how does the receiving node know what to do with the request?

### Pressure 2: local-first decision

If the target service is local, the node should serve locally.

That preserves:

- the fast path
- simpler debugging
- lower latency
- a more direct operator model

But the moment the service is not local, the system needs stronger truth than a
local Compose file can provide.

### Pressure 3: wrong-node success

If the target service is remote, the receiving node must know:

- that the service is not local
- which peer currently hosts it
- whether that peer is eligible now
- whether the route survives the local failure condition
- whether the service contract remains the same after handoff

This is the actual anti-SPOF problem, not simply "more nodes."

### Pressure 4: backend-loss survival

Even if wrong-node forwarding works once, that is still weaker than:

> when the preferred local backend dies, the recovery path remains available

That is why route persistence is such a hard boundary in the repo.

### Pressure 5: stateful honesty

Even if the HTTP story improves, stateful systems are still their own class.

That means the architecture must keep a hard line between:

- ingress continuity
- stateless request continuity
- stateful correctness

This line is not optional.
It is one of the healthiest instincts in the entire repo.

## The goals, stated the hard way

The goals should not be phrased as generic aspirations.
They should be stated in the language of the pain they are trying to remove.

### Goal 1: remove sacred-node dependence without pretending first-hop plurality is enough

Success means:

- any healthy public node can receive traffic
- that does not quietly collapse into one box that everyone still treats as the
  real entrypoint

Failure mode to resist:

- calling DNS plurality anti-SPOF when request preservation is still weak

### Goal 2: make wrong-node requests survivable

Success means:

- a request landing on the wrong node is not fatal by default
- the node has current enough truth to make a correct next decision

Failure mode to resist:

- narrating routes, proxies, or DNS records as if they already proved this

### Goal 3: preserve readability while adding only the truths that have earned their keep

Success means:

- Compose remains the primary human authoring contract
- additional layers are explainable in terms of the exact burden they remove

Failure mode to resist:

- importing control-plane tax just to sound more serious

### Goal 4: stop hiding the control plane inside the operator's head

Success means:

- topology, placement, and convergence truths become inspectable
- fewer critical facts live only in memory

Failure mode to resist:

- systems that are technically multi-node but still operationally depend on one
  person privately remembering the real topology

### Goal 5: keep stateful honesty much harsher than stateless comfort

Success means:

- stateful services are not promoted by the same weak evidence that might be
  enough for a first HTTP path drill

Failure mode to resist:

- confusing reachability with resilience

## Bottom line

The problem this repo is trying to solve is not:

- how to deploy more containers

It is:

> how to preserve real options, real operator clarity, and real request
> continuity when multiple Docker nodes are deliberately exposed to the same
> traffic surface, without immediately paying for a giant orchestrator that has
> not yet proven it is the smallest honest answer

That is the pressure every other page should inherit.

If a page softens that into generic HA language, it is already drifting away
from the user's actual dream.
