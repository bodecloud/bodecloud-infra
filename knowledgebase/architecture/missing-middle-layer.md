# The Missing Middle Layer

This page exists because the user's real complaint is not Docker Compose is
limited.

The complaint is much sharper:

> why does everything feel flexible and empowering while the stack is local,
> but the second requests can land on the wrong node, services are distributed
> on purpose, DNS stops being singular, and failure has to be absorbed instead
> of merely noticed, the market suddenly offers either fragile glue or giant
> orchestrators that want total compliance?

That is the reason this repo keeps searching for a missing middle.

If this page is weak, the whole knowledgebase collapses into a false binary:

- keep hacking raw Compose forever
- or admit that real infrastructure means handing the problem to Nomad, k3s,
  Kubernetes, or something similarly heavyweight

That binary is exactly what the project is trying to challenge.

It is not just technically unsatisfying.
It is psychologically insulting to the user.

The binary says:

- either keep all the hidden glue yourself
- or accept that real adulthood means losing the system to a larger ideology

This page matters because the user is refusing that framing itself, not merely
choosing one side of it.

It also matters because this is one of the easiest pages to accidentally ruin
with smart-sounding synthesis.

If this page becomes:

- a market comparison
- a lightweight-versus-heavyweight feature table
- or a vague essay about orchestration maturity

then it has already answered a smaller neighboring question than the one this
repo is actually asking.

This page has to behave more like reconstruction from pressure than summary
from categories.

It also has to do something stricter than merely defend "lighter than
Kubernetes."
It has to identify the smallest additional layer that would genuinely reduce
the user's reconstruction burden instead of just relocating it behind nicer
language.

## The real problem statement

The project is not primarily asking:

- how do we deploy more containers?
- how do we make YAML cleaner?
- how do we imitate cloud-native best practices?

It is asking:

> what is the smallest honest control surface that lets several ordinary Docker
> nodes behave like one request-preserving platform, without forcing the
> operator to accept a heavyweight scheduler before it has proven that smaller
> answers are inadequate?

That is why the repo keeps converging on the same pressure points:

- manual service placement is acceptable
- hidden operator memory is not acceptable
- multiple public nodes are desirable
- DNS plurality alone is not enough
- wrong-node request loss is unacceptable
- stateful honesty matters more than flattering HA language

Those pressures appear across:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- the archive material around multi-node Docker, dynamic routing, failover,
  Compose alternatives, and orchestrator tradeoffs

The recurring frustration is consistent:

single-node Docker feels like ownership;
multi-node Docker too often feels like punishment.

The missing middle layer is the repo's attempt to find a better answer than
just accept that punishment or just surrender to a giant scheduler.

That is why this page should not read like a market comparison.
It should read like a reconstruction of the specific thing the user feels is
missing from today's option space:

not "another orchestration product,"
but a legible truth-owning layer that shows its work.

That phrase matters.

The wanted layer does not merely have to coordinate.
It has to be explainable.

That explainability requirement is what keeps this page from collapsing into a
tool-family comparison.

The user is not only asking:

- which system can route
- which system can fail over
- which system can schedule

They are also asking:

- where does the new truth live
- who can inspect it
- what still depends on remembered operator lore
- what disappears when the local backend disappears

If a middle layer does not answer those questions cleanly, it may still be an
interesting product, but it is not the category this repo is trying to find.

That means the page should be read with a harsher standard than ordinary
architecture prose:

- preserve what the current runtime does not own
- preserve what the planning layer wants to promote
- preserve why the user keeps rejecting answers that still sound “serious”
  while leaving sacred-node memory intact

If those are smoothed into one clean middle-ground story, the page becomes
less faithful exactly where it needs to be most careful.

If it solves problems by disappearing truth into a black box, the user may gain
functionality while losing the thing they are actually defending:
operator-readable trust.

## The user's hidden complaint about the market

The user is not only criticizing the repo.
They are criticizing the shape of the option space itself.

The market keeps offering:

- local tools that feel empowering until distribution matters
- proxy recipes that stop at first-hop plurality
- cluster platforms that demand worldview conversion before proving they solve
  the exact pain
- "HA" answers that quietly preserve sacred-node memory underneath

The missing middle layer is the name for the category the user wishes existed
more cleanly:

- enough truth to survive wrong-node traffic
- enough inspectability to stay readable
- enough explicitness to stop private topology memory from being the real
  control plane
- but not so much machinery that the operator loses the system to abstractions

This is the important nuance many tool comparisons miss.

The user is not merely asking for something lighter than Kubernetes.
They are asking for something that still behaves honestly on a bad day:

- when a request lands on the wrong node
- when the preferred backend is gone
- when a helper route is stale
- when state did not actually move with the workload
- when the operator needs to understand what happened quickly

That is why this page matters.

Bad-day behavior is the real benchmark.

Many platforms look competent on a green dashboard.
The user's standard is closer to:

> when reality goes sideways, can the system still explain itself without
> forcing me to reconstruct hidden topology lore from memory?

This page should therefore be read with one hard rule:

bad-day inspectability outranks respectable architecture vocabulary.

If a candidate solution sounds mature but becomes harder to explain under:

- wrong-node entry
- peer loss
- stale placement truth
- secret drift
- route disappearance
- stateful ownership conflict

then it is not a real middle answer for this repo.

## What middle layer means here

This phrase becomes useless if it just means some automation.

Here it has a hard meaning:

it is the smallest extra truth-bearing layer that closes the gap between:

- human-readable Compose authoring
- manually chosen service placement
- multiple public entry nodes
- and real wrong-node request preservation

without claiming more than it proves.

The middle layer is not supposed to make the system sound modern.
It is supposed to make the system stop lying.

That sentence is close to the real one-line definition of the whole repo.

That is the cleanest one-line definition on this page.

If a future control layer routes better, schedules better, or self-heals better
but still leaves the same sacred facts socially held, then it may be useful,
but it is not yet the missing middle this repo is looking for.

## Minimum properties of a real middle layer

To stop this phrase from drifting, the page needs a stricter checklist.

A real middle layer for this repo would need to provide:

1. explicit placement truth:
   where the service currently lives is recorded in a shared inspectable form
2. explicit peer eligibility truth:
   the system can distinguish between a reachable peer and a semantically valid
   peer for the class of service being forwarded
3. fallback-path durability:
   the route needed for rescue does not vanish with the local failure it is
   supposed to compensate for
4. request-meaning preservation:
   forwarded traffic keeps the same auth, middleware, and service identity
   semantics
5. class-sensitive honesty:
   stateless and stateful services are not given the same emotional grade just
   because both can be proxied
6. operator readability:
   when something fails, the explanation can be recovered from tracked system
   surfaces rather than ritual memory

That is the minimum bar.

If a proposed layer only gives:

- service discovery without fallback durability
- forwarding without policy continuity
- stateful reachability without stateful honesty
- automation without readability

then it is still a partial answer, not the missing middle.

This is what keeps the category from dissolving into "some automation plus a
registry."

## What would still be fake middle-layer theater

The repo also needs a negative checklist, because fake middle layers are often
exactly what look attractive at first glance.

A candidate is still theater if:

- it introduces shared config but not shared live truth
- it forwards requests but loses middleware, auth, or request identity under
  that handoff
- it centralizes the real answer into one helper node or one helper process
  while preserving multi-node rhetoric
- it can recover only while the operator already knows which peer is special
- it narrates stateful backends as "covered" because a proxy target exists
- it reduces YAML pain while leaving bad-day causal explanation harder than
  before

That last failure matters a lot.

This repo is not trying to trade one form of toil for another prettier one.
It is trying to reduce the number of truths that only exist socially.

That means a good middle layer must be understandable in failure language, not
just architecture language.

An operator should eventually be able to ask:

- why did this node serve locally?
- why did it forward?
- why did it refuse to forward?
- why did the fallback route survive?
- why is this still unsafe for state?

and get answers from tracked system truth rather than social memory.

## Why raw Compose alone stops being enough

Compose remains the correct live authoring surface for this repo today.
That is an evidence-based statement, not nostalgia.

Compose is genuinely good at:

- readable service definitions
- explicit networks
- configs and secrets
- healthchecks
- modular composition
- per-node deployment clarity

But the key question changes in a distributed world.

Single-node Docker asks:

> what should run on this machine?

The user's real distributed problem asks:

> what should happen when traffic for a non-local service lands on this machine
> right now, under partial failure, without inventing fake certainty?

Compose does not naturally own the answers to:

- where the service actually lives right now
- whether peer nodes are converged enough to substitute
- whether the recovery route survives local failure
- whether the remote peer is eligible for forwarded traffic
- whether auth and middleware meaning survive the hop

That is not a flaw in Compose.
It is the point where a local runtime description stops being a sufficient
distributed truth source.

That distinction is important because this page is not anti-Compose.
It is anti-fantasy.

Compose remains the right authoring anchor until a stronger layer proves that
it relocates enough of the missing truth burden to justify itself.

The real challenge is not to "graduate from Compose."
It is to stop asking Compose to pretend it already owns distributed truths that
it does not naturally own.

This is also why the user keeps sounding more frustrated than a normal
"should I use Nomad or k3s?" conversation would imply.

The pain is not only that Compose stops being enough.
The pain is that the gap after Compose is full of answers that either:

- hide too much
- assume too much
- centralize too much
- or solve one truth while forcing the operator to manually reconstruct three
  others

That last failure mode is the one the docs must keep visible.

The user is not seeking a magical tool that solves everything.
They are seeking a control layer that does not cheat by paying down one burden
while re-hiding the rest under smoother language.

## How future candidates should be judged

This page should make later comparison pages harsher, not softer.

Before any proposed control layer earns emotional legitimacy here, it should be
interrogated with questions like:

- what new truth surface does it create that the operator can actually inspect?
- which current hidden burden does it relocate out of memory?
- what happens when traffic lands on the wrong healthy node?
- what happens when the preferred local backend is gone?
- what remains true for stateless HTTP services that is still false for TCP or
  stateful systems?
- which part of the claimed resilience is behavior and which part is still
  narration?

If the answer set is still fuzzy, then the solution may still be worth
researching, but it has not yet earned the status of "the missing middle."

This is another place where the docs need to act more like a held-out test
than like a summary:

- does the proposed layer actually move truth out of private memory?
- does it survive a wrong-node request without folklore?
- does it explain itself on a bad day?

If the answer to those is still unclear, the layer is still descriptive hope,
not yet the missing middle the user is searching for.

## The minimum job description of the missing middle

If the repo wants to remain Compose-first while growing into a real multi-node
request-preserving platform, the middle layer has to own a small but serious
set of truths.

Not all possible truths.
Only the ones the rest of the dream collapses without.

### 1. Placement truth

It must answer:

> what runs where right now?

Not:

- what should run where according to a plan
- what used to run there yesterday
- what the operator vaguely remembers

This is why `services.yaml` or an equivalent registry keeps resurfacing.

Without placement truth:

- a node cannot tell whether the target is local
- a node cannot know which peer should receive forwarded traffic
- route generation remains theater

### 2. Convergence truth

It must answer:

> if I forward this request to another node, am I forwarding into the same
> semantic service or just a similarly named container?

That means visibility into enough of the following to make cross-node fallback
honest:

- revision state
- env drift
- secret drift
- auth and middleware dependencies
- deployment freshness

Without convergence truth, distributed fallback becomes a confidence trick:

transport succeeds
while meaning quietly diverges.

That warning should echo across the whole site.

This repo is not satisfied with transport success alone.
It is trying to guard against systems where the request technically survives
while the semantic contract quietly changes underneath it.

This is one of the places the repo is sharper than normal infrastructure
writing.

The user is not satisfied by "the peer answered."
They are asking whether the peer answered as the same service contract or as a
different local accident with the same name.

### 3. Route persistence under failure

It must ensure the route needed for recovery survives the local failure that
made recovery necessary.

That is the point of the repeated warning around `docker-gen-failover` and
similar ideas.

If the local container dies and the recovery route disappears with it, then the
platform does not have failover.
It has a diagram of failover.

That distinction is one of the most important in the entire documentation set.

The user is reacting against systems that are architecturally flattering before
they are behaviorally honest.
This page should keep naming that difference bluntly.

This is one of the most important distinctions in the entire repo.

### 4. Peer eligibility truth

It must distinguish between:

- a peer that exists
- a peer that is locally healthy
- a peer that is safe to receive forwarded traffic

Those are not the same thing.

This is where many systems start overclaiming because reachable feels close
enough to good.
It is not close enough.

Without eligibility truth, wrong-node forwarding can succeed at transport level
while still producing user-facing breakage or policy drift.

### 5. Policy continuity

It must preserve the meaning of the request path, not merely deliver bytes to a
remote process.

For HTTP that includes:

- host semantics
- middleware continuity
- auth continuity
- TLS and header assumptions
- policy identity at the edge

If a forwarded request returns 200 but dropped the expected auth or middleware
meaning, the system did not actually preserve the request.

This repo is smarter than many generic HA guides on this point.
It repeatedly insists that reachable backend is not the same as preserved
contract.

### 6. Service-class honesty

It must keep different service classes separate.

That means staying honest about:

- stateless HTTP
- raw TCP or L4 services
- stateful systems

The middle layer may help all three.
It cannot narrate all three the same way.

This matters because the easiest way to fake progress is:

- solve one HTTP route
- talk as if TCP and stateful HA are morally adjacent

They are not.

This is the point where many "middle layer" dreams quietly become fake.

If the proposed layer only tells a convincing story for stateless HTTP, that
may still be progress.
But it should then be narrated as:

- an HTTP middle layer
- not a universal HA answer
- not a solved data-plane answer
- not a reason to stop thinking about state authority

## What the middle layer is not supposed to do

The user's frustration matters here because it tells us what failure modes to
avoid.

The missing middle is **not** supposed to:

- become a stealth orchestrator while the docs still call the stack just
  Compose
- hide topology so thoroughly that debugging requires reverse-engineering the
  helper layer itself
- flatten HTTP, TCP, and stateful systems into one badge labeled HA
- confuse Cloudflare multi-node entry with end-to-end request preservation
- replace operator control with opaque automation that still fails under real
  pressure

The project is not anti-automation.
It is anti-automation that charges cognitive rent without paying down the
actual pain.

This distinction should remain visible everywhere in the repo.

The user's real test for a helper is not:

> does it automate something?

It is:

> after adding it, do I understand the distributed behavior more clearly than
> before, or did I just move the mystery into another box?

## Why Cloudflare is necessary but insufficient

The repo is right to care about multi-node public entry.

Cloudflare and plural A or AAAA records help with:

- removing the sacred public box
- allowing any healthy public node to become the first hop
- making ingress less hostage to one machine

But Cloudflare does not solve the missing middle by itself.

It does not answer:

- where the service actually lives
- whether peers are converged
- whether the route survives local backend loss
- whether auth and middleware continuity survive the hop
- whether a stateful service remains correct under failover

Cloudflare solves plurality at the first hop.
The missing middle solves truth after the first hop.

Those are different jobs.

## The main candidate forms of the middle layer

The repo does not need one final universal answer yet.
It does need a sober reading of the main families of answer.

### Form 1: Compose-first plus explicit truth helpers

This is still the strongest pressure direction already visible in the repo.

Examples:

- `services.yaml`
- sync-agent
- failover-agent
- generated Traefik dynamic config
- CUE-shaped or generator-shaped control ideas

This family says:

- keep Compose as service intent
- keep manual placement where it still makes sense
- add narrowly scoped truth and routing helpers behind the scenes

Why this form is attractive:

- preserves authoring readability
- resists premature orchestrator capture
- lets the project solve only the missing truths it actually needs

Why this form is dangerous:

- helpers can quietly accumulate into a de facto control plane
- evidence can lag far behind architecture language
- debugging burden can move from Docker into bespoke glue

This danger matters more in this repo than in a generic homelab.

If a helper mesh eventually owns:

- placement truth
- peer eligibility
- route generation
- failure reaction
- convergence checks
- promotion logic

then the honest question stops being "can we stay Compose-first?"
The honest question becomes:

> have we already built a control plane and just chosen not to name it yet?

This is probably the best near-term direction if the repo keeps the proof
boundaries brutally honest.

### Form 2: HA substrate plus Compose-facing runtime

OpenSVC and related substrate ideas live in this space.

This family says:

- keep service authoring closer to Compose
- let a stronger cluster substrate own membership, failover, or
  routing-relevant truth

Why it is attractive:

- it can shoulder harder recovery and HA burdens without immediately importing
  full Kubernetes semantics
- it may solve some truth layers more explicitly than custom helper meshes

Why it is dangerous:

- it still introduces a real control plane
- it can become another abstraction the operator must decode
- it must still prove that it reduced net burden instead of moving it

That last point is the real filter.

The user is not looking for "less Kubernetes."
They are looking for a layer whose added complexity can be defended in the
language of specific missing truths rather than in the language of platform
prestige.

This form matters because it may be the genuine middle between bespoke glue and
full cloud-native orchestration.

### Form 3: scheduler promotion

Nomad, k3s, and Kubernetes live here.

This family says:

- stop trying to reconstruct distributed truth by hand
- promote a scheduler or cluster API to own placement, health, and service
  discovery at a larger scale

Why it is attractive:

- it offers a more complete answer to the missing truth problem
- it can make route, placement, and eligibility stories more native

Why it is resisted:

- it imports significant abstraction and tax
- it often demands surrender of the legibility that made Compose attractive in
  the first place
- it is frequently recommended too early, before smaller honest answers are
  exhausted

The repo's stance is not never.
It is:

> not before the need is proven

And the need is not proven merely because the user is frustrated.
It is proven when the current layer keeps failing the same concrete request,
placement, convergence, or policy questions and the next layer can answer them
more honestly.

That distinction matters.

## The real evaluation question

The middle layer should not be judged by style, popularity, or benchmark
prestige.

It should be judged by this:

> did it remove hidden operator burden and false resilience claims without
> replacing them with a larger, equally frustrating fiction?

That question is the filter for every candidate:

- bespoke helper mesh
- OpenSVC-like substrate
- Nomad
- k3s
- Kubernetes

The question is never merely:

> can this tool route traffic?

It is:

> which missing truth does this tool own, how inspectable is that truth, what
> tax does it import, and does that tax finally pay down the wrong-node and
> backend-loss lies that Compose alone cannot fix?

## What the user is really trying to force into existence

The user is not just asking for a better document.
They are trying to force a more honest infrastructure category into existence:

a category where:

- Docker readability survives
- multiple nodes are first-class
- wrong-node requests are not fatal by default
- Cloudflare plurality is only the first hop, not the whole story
- stateful honesty is preserved
- operator understanding is not sacrificed to buy survivability

That is the dream behind the missing middle layer.

The repo does not yet fully prove that dream in runtime form.
But it does prove the dream is not naive or vague.

It is a precise demand:

> give me the smallest truth-owning layer that makes distributed Docker stop
> behaving like a trap, while refusing to pretend that larger orchestrators are
> automatically the only adult answer

That is what this page is protecting.
