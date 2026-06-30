# Unified Orchestration Blueprint

This page exists to answer the orchestration question the repo actually keeps
asking, not the softer question that infrastructure docs usually answer.

The soft question is:

> which orchestrator should we use?

The real question is:

> what is the smallest additional control layer that makes wrong-node requests,
> backend loss, and hidden human memory stop being the dominant source of
> failure, while preserving the Compose-first surface the operator still
> trusts?

That is the real blueprint.

If this page loses that question, it turns back into generic platform
comparison prose and stops being useful.

It also becomes dangerous in a more specific way:

- the control-layer story starts sounding cleaner
- cleanliness starts sounding like shared truth
- shared truth gets inferred even when the runtime still depends on remembered
  placement and remembered fallback interpretation

That is exactly how a blueprint can sound advanced while still preserving the
same human SPOF under a more respectable vocabulary.

There is an even harsher standard underneath it:

> can the documentation reconstruct the actual system the user is trying to
> build, including the parts they have not named cleanly yet, without flattening
> that system into a safer adjacent question?

That matters because the orchestration problem in this repo is not just a tool
problem.
It is also a reconstruction problem.

The user is surrounded by answers that keep taking a real demand:

- any-node entry that is not fake
- request preservation that survives the wrong receiving node
- a multi-node Docker world that does not immediately become a giant scheduler
- honest anti-SPOF language instead of theatrical redundancy

and shrinking it into easier stories:

- "pick an orchestrator"
- "add better load balancing"
- "improve service discovery"
- "document the existing stack more clearly"

Those are all smaller than the thing being asked for.
This page is only useful if it keeps refusing that shrinkage.
It also has to refuse a subtler failure mode:

- the blueprint starts sounding coherent
- coherence starts sounding like settlement
- settlement starts sounding like a live operating model
- the unresolved ownership question quietly disappears

## What a real blueprint has to do in this repo

In a normal infrastructure repo, a blueprint can get away with listing
components, environments, and orchestration candidates.

That is not enough here.

In this repo, a real blueprint has to answer a harsher question:

> what exact missing truth is making the current system feel fake, and what is
> the smallest added layer that would make that specific truth stop depending
> on operator memory?

If a proposed layer cannot answer that question, it is not yet part of the
solution. It is just more architecture vocabulary.

It also has to answer a second question:

> what part of the user's real dream would still be missing even if this layer
> were added successfully?

That second question is what stops the blueprint from overclaiming.
Many candidate layers can improve one pain while still leaving the deeper wound
untouched.

There is a third question this page has to keep asking:

> if this layer appears to work, how much of that success is system-owned truth
> and how much is still operator reconstruction that merely became easier to
> narrate after the fact?

That question is what stops "better architecture" from being confused with
"the system truly owns the missing truth now."

The failure mode this repo keeps encountering is not "no one offered enough
technology."
It is "many technologies were offered as if they closed more of the problem
than they actually do."
That sentence is one of the most important filters in the whole knowledgebase.
It means this page cannot merely compare power, elegance, or ecosystem depth.
It has to keep reconstructing the exact wound each candidate does and does not
close.

## This page is not a chooser, it is a filter

This repo does not need another winner bracket between:

- Docker Compose
- helper agents
- CUE
- OpenSVC
- Nomad
- k3s
- Kubernetes

It needs a much harsher filter:

> which layer removes a real hidden-human SPOF, a real wrong-node request
> failure, or a real convergence failure, and which layer merely renames the
> problem while charging more worldview tax?

That is the only orchestration question worth keeping.

Another way to say the same thing is:

the winning layer is not the most powerful one.
It is the one that takes ownership of the exact truth the operator is currently
carrying privately, while damaging the operator-readable surface as little as
possible.

## The blueprint should be read as a burden-transfer map

This page becomes much easier to use when each candidate layer is judged by
one blunt question:

> what hidden burden moves out of the operator's head if this layer becomes
> real?

If the honest answer is "not much," the layer is still mostly theater for this
repo.

The main burdens to track are:

- remembered placement truth
- remembered peer safety
- remembered fallback-route behavior
- remembered convergence assumptions
- remembered stateful authority

That is the real comparison surface.
Not elegance.
Not ecosystem size.
Not how mature the diagrams look afterward.

The hidden corollary is important:

if the truth still has to be privately carried, then the layer did not win just
because it reduced the number of boxes or made the diagram look more modern.

## The dream this blueprint is trying to protect

`bolabaden-infra` does not actually want one of these simplistic stories:

- "there is no control plane, only Compose"
- "the repo has already chosen Kubernetes, Nomad, OpenSVC, or an agent mesh"
- "the user just wants better load balancing"

What it wants is harsher and more specific:

- keep Compose as the readable authoring contract
- keep manual placement acceptable unless stronger automation truly earns its
  keep
- allow any healthy public node to be the first hop
- preserve local-first service when locality exists
- preserve the request when traffic lands on the wrong healthy node
- stop pretending DNS plurality or proxy presence means failover is solved
- stop describing stateful systems as anti-SPOF when only the ingress path is
  redundant

That is the problem this blueprint is trying to solve.

The dream is not merely "build a homelab that is more serious."
It is closer to:

> make several ordinary Docker machines behave like one personal cloud that does
> not lie about where its intelligence lives

That is why this repo is so hostile to sacred nodes, hidden placement memory,
and fake HA wording.
The user is not asking for prettier control surfaces.
They are asking for the intelligence of the system to stop living in one
person's head while still remaining legible.
That last clause matters.
The repo is not just anti-hidden-truth.
It is anti-hidden-truth that reappears inside a more sophisticated control
plane and calls itself success.

## What the blueprint must preserve

Any future control layer is only aligned if it preserves the things the user is
actually defending.

### 1. Compose remains the operator-readable authoring surface

The root implementation still centers on:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- included fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

That is not temporary trivia.
It is the current operator contract.

Any stronger layer that requires the operator to stop understanding the stack
through those files has to prove it is paying down a real pain that Compose
cannot solve honestly.

### 2. First-hop plurality is necessary but not sufficient

The user does not want one sacred public box.

Cloudflare and multi-record public entry matter because they allow:

- more than one node to receive the first request
- node loss without immediate public dead-end

But the blueprint must not stop there.
First-hop plurality is only the start of the story.

### 3. Local-first service is part of the dream, not an optimization detail

If the requested service is already on the receiving node, that node should
just serve it.

That matters because the user is trying to preserve:

- directness
- locality
- debuggability
- operator trust in the fast path

The repo is not trying to hide everything behind a cluster fiction for the
sake of looking modern.

### 4. Wrong-node requests are the real orchestration test

This is the deepest coordination requirement in the whole repo.

The blueprint only earns its name if it explains how a receiving node can
determine:

- the service is not local
- which peer currently hosts it
- whether that peer is eligible now
- whether the route still exists under the relevant failure
- whether auth, middleware, headers, and visible request semantics stay intact

This is the real anti-stupidity layer the user keeps asking for.

It is also where the docs have to resist a very common flattening move:

- "the node can probably forward it"
- "the registry could point at the right peer"
- "Traefik or the helper stack should be able to preserve the route"

Those are all weaker than what the user is actually asking.

The real standard is:

- can the node know enough current truth to hand the request off correctly
- can it know that without private operator memory closing the remaining gaps
- can it preserve visible semantics while doing so
- can the docs say exactly which of those statements are live, planned, or still
  only inferred

This is also the point where many otherwise respectable infrastructure options
stop being real options for this repo.

An option that only improves:

- first-hop plurality
- local happy-path routing
- cleaner static config

while still leaving the wrong-node path socially reconstructed is not solving
the thing the user is actually mad about.

### 5. Stateful systems stay in a separate honesty bucket

The blueprint must never let HTTP request preservation quietly blur into
stateful correctness.

Any control layer still has to answer, separately:

- who owns writes
- how replication or election works
- how clients rediscover the real authority
- what breaks on node loss
- whether the real failure domain is still one disk or one host

If the blueprint loses that separation, it becomes fake HA again.

It also becomes much easier to tell the comforting but wrong story that "the
hard part is now mostly ingress."

For this repo that story is a regression.
The whole point is to stop letting a stronger edge posture impersonate
ownership of the deeper system truth.

### 6. The docs must preserve the user's hidden negative benchmark

The negative benchmark in this repo is not "worse than Kubernetes."
It is:

- the route only works if requests land on the special node
- the operator still remembers where the real service lives
- the fallback exists only on paper or only at the first hop
- middleware, auth, or request semantics break when the handoff becomes real
- the system sounds distributed while its real coordination still depends on
  social memory

If a future control layer leaves those conditions mostly intact, then however
modern it sounds, it has not crossed the user's real bar.

## The candidate-layer scorecard

This repo needs a sharper filter than "lighter versus heavier orchestrator."

| Candidate layer | What it can honestly buy | What it still does not buy by itself | Failure mode if oversold |
| --- | --- | --- | --- |
| Compose-first plus thin helpers | preserves direct authoring readability; can host placement, eligibility, and route-generation truth narrowly | stateful topology ownership, generic scheduling, broad reconciliation | helper mesh becomes a shadow control plane while the docs still pretend it is "just Compose" |
| Schema-first extension such as CUE-style semantics | stronger declared intent, cleaner generation surfaces, better explicit meaning than scattered labels | live convergence, live peer judgment, failure-time route persistence | the repo starts sounding more explicit while runtime truth is still not owned |
| Agent-first or active-control layer | active observation, sync, route generation, and convergence reactions become possible | honest simplicity, stateful semantics, and bounded control-plane scope | the agents quietly become an orchestrator in disguise without paying down enough pain |
| Narrow infra-grade HA promotion such as OpenSVC-shaped ingress work | can remove sacred ingress or identity surfaces without promoting the whole app layer | app-level wrong-node preservation for all services, stateful correctness | first-hop continuity gets narrated as if the whole service surface is now resilient |
| Lighter scheduler promotion such as Nomad-shaped workloads | real placement and rescheduling ownership for selected stateless classes | stateful truth, ingress semantics, and protocol-specific failover by default | the repo calls scheduling maturity "platform maturity" too early |
| Full desired-state platform such as k3s or Kubernetes | broad reconciliation, ecosystem depth, controller patterns, stronger platform conventions | automatic honesty about state, request meaning, or cross-domain SPOF language | worldview capture arrives before the exact pain it removes is named plainly |

This table is intentionally unfair to broad answers.
The repo has already seen too many broad answers that improve a category while
still leaving the same hidden burden intact at request time.

## What the user is explicitly refusing

The orchestration future is shaped as much by refusals as by aspirations.

### Refusal 1: "just use the big orchestrator"

The user is not refusing Kubernetes, k3s, or Nomad because they are
unfamiliar.
The user is refusing the lazy move where a huge control-plane worldview gets
introduced before the docs can prove which exact pain that worldview is paying
down.

The recurring archive pressure is:

- raw Compose becomes too static
- common answers leap immediately to heavy orchestration
- the real missing middle layer never gets named clearly

This repo exists because that leap feels dishonest and coercive.

The archive makes something else clear too:

the user is not refusing Kubernetes because they are attached to simplicity as
an aesthetic.
They are refusing being pushed into a giant worldview before anyone can explain
why the smaller missing truth layer cannot be solved more directly.

### Refusal 2: DNS theater

The user is openly tired of answers that stop at:

- multiple A or AAAA records
- round-robin entry
- "some healthy node will answer"

That is not the same thing as preserved service success.

A lot of otherwise competent infrastructure writing still treats a multi-node
first hop as if it already reduced the most humiliating failure:

the request worked only because everyone involved still knew which node was the
real one.

### Refusal 3: hidden human coordination as the real control plane

The user does not want the real cluster truth to keep living in:

- one remembered machine
- one remembered route
- one remembered storage path
- one remembered service location
- one remembered private fact in the operator's head

If the system still depends on those, the control plane is still human memory.

That is why this blueprint cannot be read as a product-comparison page.

The user is not simply trying to avoid Kubernetes.
They are trying to avoid one more system where the real intelligence still
lives outside the runtime, except now it is harder to inspect.

This is the deepest refusal of the whole repo.

The user can tolerate complexity.
What they do not tolerate is complexity that still leaves the same sacred human
memory in charge, only with nicer diagrams around it.

That is the quality bar this page has to keep visible.

If a future layer mostly improves:

- aesthetics
- component legitimacy
- configuration regularity
- post-hoc explainability for people already familiar with the stack

while still leaving real handoff truth privately carried, then it has not met
the user's bar even if most infrastructure readers would call it a major
upgrade.

### Refusal 4: stateful marketing language

If a database is reachable through a global name but still depends on one real
writer, one promotion ritual, one real disk path, or one fragile replication
assumption, the user wants the docs to say that bluntly.

The blueprint has to preserve that hostility to fake HA language.

### Refusal 5: documentation that becomes coherent by shrinking the question

This page also has to refuse a documentation failure mode.

The easiest way to make a repo like this sound intelligent is to replace the
real demand with a more conventional one:

- "needs service discovery"
- "needs a scheduler"
- "needs a stronger ingress layer"
- "needs standard clustering"

Those may all be adjacent truths.
None of them are automatically the real truth.

If the docs become clearer by quietly asking a smaller question than the user
is asking, they have failed even if every sentence is technically accurate.

## The coordination duties the missing layer must take over

The repo does not need "orchestration" in the abstract.
It needs a control surface that takes over specific duties ordinary Compose does
not close on its own.

### Duty 1: placement truth

Some live source has to answer:

- what runs where right now?

This is why `services.yaml` keeps reappearing across the repo.
It is not a random implementation detail.
It is the clearest expression of the missing placement-truth layer.

### Duty 2: peer eligibility truth

Knowing a service nominally lives on a peer is not enough.
The receiving node needs to know whether the peer is currently safe to use for:

- health
- compatible config
- secret parity
- policy parity
- route viability

This duty is what separates “service discovery” from “service recovery.”

Lots of stacks can discover a candidate peer.
Far fewer can prove that the peer is valid for this request now.

### Duty 3: route persistence under failure

The blueprint must stop the classic failover lie:

- a fallback route exists while the primary is healthy
- then disappears exactly when the primary fails

That is why route generation and failover replacement mechanisms have to be
judged by failure-time behavior, not by happy-path presence.

### Duty 4: convergence truth

If nodes disagree on:

- env values
- secrets
- image revisions
- middleware assumptions
- service placement

then a forwarded request may still "work" mechanically while violating the
operator's actual expectation.

### Duty 5: auditability

A real middle layer must leave the operator able to explain:

- why a request stayed local
- why it forwarded
- why a peer was chosen
- why a route disappeared
- which failure class is actually solved versus merely masked

If the answer becomes "the controller decided," the layer is becoming
misaligned unless it has clearly earned that complexity.

This is one of the most important anti-theater rules in the whole repo.

The user is willing to tolerate more machinery if it actually removes hidden
burden.
They are not willing to tolerate a black box that replaces one kind of
guesswork with another.

## The three strongest orchestration instincts in the repo

The tree does not show one settled orchestration future.
It shows several serious instincts reacting to the same pressure.

### 1. Compose-first extension

Representative surfaces:

- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`architecture/compose-first-architecture.md`](architecture/compose-first-architecture.md)

This instinct says:

- keep Compose as the language the operator still sees first
- add the missing truth layers beside it
- do not import a whole scheduler unless the thinner path is exhausted

What it buys:

- readability
- locality
- lower worldview tax

What it still lacks alone:

- live placement truth
- live peer eligibility truth
- durable wrong-node recovery

### 2. Schema-first extension

Representative surfaces:

- [`CUE_SPEC_EXTENSIONS.md`](CUE_SPEC_EXTENSIONS.md)
- [`CUE_BOOTSTRAP_PROTOCOL.md`](CUE_BOOTSTRAP_PROTOCOL.md)

This instinct says:

- Compose is still the familiar shape
- but higher-order semantics should stop living in scattered labels and tribal
  memory
- describe HA mode, placement, visibility, dependencies, and recovery intent
  explicitly

What it buys:

- stronger declared semantics
- better generation opportunities
- cleaner operator-readable intent

What it still does not buy by itself:

- active runtime convergence
- active peer health judgment
- actual route persistence under failure

### 3. Agent-first or active-control extension

Representative surfaces:

- [`infra/docs/ARCHITECTURE.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/infra/docs/ARCHITECTURE.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`research/infrastructure-master-plan.md`](research/infrastructure-master-plan.md)

This instinct says:

- some parts of the problem are active coordination problems, not just config
  description problems
- software should observe and maintain cluster truth instead of leaving it in
  the operator's head
- failover and sync behaviors may need explicit agents rather than only
  declarative metadata

What it buys:

- active convergence potential
- stronger route-generation inputs
- a path toward runtime truth rather than only design intent

What it risks:

- recreating an orchestrator in disguise
- growing a private control plane that is no longer obviously simpler than
  selected scheduler promotion

## What each path would still leave unresolved

The blueprint needs this section because this repo is especially vulnerable to
"good enough in one domain" quietly becoming "settled overall."

### If the repo stays Compose-first with thin helpers

Still unresolved unless separately solved:

- stateful election, promotion, and reconnect truth
- whether helper-layer convergence logic is now controller-sized
- whether wrong-node proof survives backend loss rather than only happy-path
  remote forwarding

### If the repo promotes a narrow infra HA layer

Still unresolved unless separately solved:

- whether application services can be preserved from the wrong node
- whether edge continuity hides stateful fragility
- whether the app layer still depends on remembered placement

### If the repo promotes a scheduler

Still unresolved unless separately solved:

- whether ingress semantics stay legible
- whether stateful classes are genuinely improved or only rescheduled
- whether the user actually gained clarity rather than just a new worldview

### If the repo promotes a full desired-state platform

Still unresolved unless separately solved:

- honest cross-domain SPOF language
- service-class-specific data truth
- the risk that the runtime is now harder to read while still needing private
  expert reconstruction on a bad day

## What this blueprint already proves

### Proof 1: plain Compose is not enough on its own

The repo still trusts Compose as the best current authoring surface.
It no longer trusts plain Compose to answer the whole multi-node request-time
problem.

That is why the tree keeps generating:

- `services.yaml` ideas
- route-generation work
- CUE semantics
- sync-agent and failover-agent planning
- OpenSVC experiments
- orchestration comparisons

### Proof 2: the repo is searching for a missing middle layer, not merely stalling a Kubernetes decision

The archive pressure does not read like indecision.
It reads like repeated rejection of a false binary:

- brittle static glue
- or total platform capture

This page exists because the repo keeps trying to name the thinner layer that
should exist between those poles.

### Proof 3: any final platform choice must be justified by the wrong-node problem

The most important test for any promoted layer is not:

- can it schedule containers?

It is:

- can it make wrong-node request preservation materially real
- without lying about stateful correctness
- and without introducing more worldview tax than the pain it removes

That is the true acceptance test.

## The next honest sequence

The blueprint points toward a more disciplined implementation order than
"evaluate orchestrators forever."

### Step 1: make live placement truth explicit for one real path

Whether it is called `services.yaml` or not, the repo needs one tracked,
operator-readable, live-consumed placement answer for at least one service
class.

That is the first point where the docs stop depending on the operator's private
memory.

### Step 2: prove one stateless HTTP wrong-node path end to end

Not just:

- first-hop plurality
- peer reachability
- a route that looks plausible

But:

- wrong-node detection
- peer choice
- forwarding
- middleware and auth continuity
- logs proving why the request still succeeded

That is the first point where the docs stop depending on hope.

### Step 3: decide whether the thin layer is still meaningfully thin

If the helper layer starts absorbing:

- placement truth
- health truth
- route persistence
- convergence
- promotion logic
- topology truth

then the repo must ask whether it is still meaningfully lighter than promoting
selected domains into Nomad, k3s, Kubernetes, or another stronger platform.

This is the point where "avoid the big orchestrator" stops being a principle
and starts needing a cost accounting.

That cost accounting should be explicit and a little ruthless.

If the helper mesh now owns:

- placement truth
- eligibility truth
- route persistence
- convergence judgments
- failure reaction

then the repo should stop flattering itself with “still just Compose” language
and compare that helper mesh honestly against stronger named platforms.

### Step 4: keep stateful promotion separate

Even if stateless HTTP wrong-node behavior becomes real, stateful services still
need their own proof and promotion path.

If this step is skipped, the repo recreates fake HA at the exact layer that
hurts users most.

## The immediate architectural question after this page

After reading this blueprint, the next question should not be:

- "so which orchestrator wins?"

It should be:

- "which specific truth layer can we make live next in the smallest possible
  way, and what stronger claims would still remain forbidden afterward?"

That question keeps the repo aligned with the user's real demand.
It is also the fastest way to stop the docs from quietly turning blueprint
coherence into fake settlement.

## The promotion rule hiding underneath the whole blueprint

The blueprint can be compressed into one harsh rule:

> do not promote a new control layer because it sounds modern, cleaner, or
> more complete; promote it only when the current layer has become the hidden
> tax the operator is still paying in memory, guesswork, and failure-time
> ambiguity

This is the sentence the rest of the docs should keep using as a filter.

It is the simplest way to preserve the user’s actual dream without collapsing
into either:

- anti-orchestrator ideology
- or premature orchestrator surrender

That is the real unifying rule underneath Compose-first, schema-first,
agent-first, OpenSVC-shaped, Nomad-shaped, and Kubernetes-shaped futures.

## Bottom line

The unified orchestration blueprint is not:

- "pick the coolest orchestrator"
- "pretend Compose is enough forever"
- "stay lightweight no matter the cost"

It is:

> preserve the Compose-first surface the operator still trusts, then add only
> the smallest amount of explicit placement, convergence, and failover truth
> needed to make wrong-node requests stop being a gamble.

That is the blueprint the rest of the repo keeps converging on, even though the
tracked runtime still does not prove the whole thing live.
