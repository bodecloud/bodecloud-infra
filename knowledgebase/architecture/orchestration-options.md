# Orchestration Options: Which Extra Layer Has Actually Earned the Right to Exist?

For the evidence underneath this page, start with
[`../research/orchestrator-tradeoffs-evidence.md`](../research/orchestrator-tradeoffs-evidence.md),
[`../research/orchestration-research-2026.md`](../research/orchestration-research-2026.md),
and the current planning surface in
[`/docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md).

This page exists to stop one lazy question from taking over the repo:

> which orchestrator is best?

That is not the real decision surface here.

The real question is:

> which extra layer of control has actually earned the right to exist because
> it removes a real hidden human SPOF, a real wrong-node failure, or a real
> convergence failure that the current Compose-first stack cannot solve
> honestly?

That is the decision surface that matters.

If this page ever turns into a normal "platform comparison" page, it has
already failed.

It fails in a second, more subtle way when it starts sounding like the main
work is simply picking between respectable tool families.

That is smaller than the real problem.

The real problem is deciding which layer, if any, has earned the right to own
the truths the operator is still privately carrying:

- where the service really lives now
- which peer is actually eligible now
- whether the fallback route still exists now
- whether the request keeps its meaning after handoff

The user is not short on platform names.
The user is short on options that still feel honest after wrong-node entry,
backend loss, and hidden operator glue become real.

That is a stronger complaint than "there are too many choices."
The complaint is that many choices stop feeling like choices once the user
asks the humiliating question:

> after all the adult-sounding layers are in place, why am I still the thing
> that knows what is really true right now?

That difference is the whole reason this page exists.

## What this page is and is not allowed to prove

This page is authoritative about:

- how orchestration options should be judged in this repo
- which burden a promoted layer would need to own to justify its existence
- why respectable platform names are weaker than domain-specific truth
  ownership

This page is not authoritative about:

- whether one orchestration path has already won globally
- whether the current runtime already demonstrates the promoted behavior
- whether broader ecosystem prestige should override repo-specific evidence

This page is a promotion filter, not a final platform verdict.

It also needs to behave like a disappointment filter.
The user has already seen too many ecosystem answers that become smaller the
moment they are forced to speak in terms of actual burden transfer instead of
general platform capability.

## Priority decision stack for this page

When this page evaluates whether some extra layer has earned the right to
exist, it should route the question through this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`current-compose-runtime.md`](current-compose-runtime.md)
3. [`problem-and-goals.md`](problem-and-goals.md)
4. [`../research/orchestrator-tradeoffs-evidence.md`](../research/orchestrator-tradeoffs-evidence.md)
5. [`../operations/proof-matrix-and-drills.md`](../operations/proof-matrix-and-drills.md)
6. [Archive Pressure Patterns](../research/archive-pressure-patterns.md)

That order matters.

If this page starts from product families, it collapses into market prose.
If it starts from runtime alone, it under-reconstructs the dream.
If it starts from plans alone, it over-promotes futures into present gravity.

This page only stays honest if the dream, the live baseline, the hidden wound,
the candidate evidence, and the proof ceiling are all present at once.

If any one of those drops out, the page immediately gets worse in a predictable
way:

- lose the dream and the page turns into market prose
- lose the runtime and it becomes fantasy
- lose the wound and it becomes generic architecture advice
- lose the evidence and it becomes recommendation theater
- lose the proof ceiling and it starts promoting prestige as if prestige were
  closure

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "Compose-first is still the default" | current strategy stance | this page, `.github/copilot-instructions.md`, current runtime pages | that Compose already solves the missing truths |
| "a stronger layer has not yet earned itself" | promotion-threshold judgment | this page plus orchestrator evidence and proof pages | that stronger layers are forbidden in principle |
| "this candidate owns a specific hidden burden better" | option-family evaluation | this page plus candidate evidence | that the candidate therefore solves the whole dream |
| "prestige or maturity language is too weak" | anti-theater judgment | this page, archive pressure, proof matrix | that professional tooling has no place here |

If a sentence starts sounding like "pick the most mature orchestrator," it has
already left this page's decision surface.

That is because "mature" is one of the most dangerous words in this problem
space.
It too easily means:

- more ecosystem
- more automation
- more abstraction
- more seriousness of tone

when the user is actually asking for something much narrower and harsher:

- one more option that remains believable after the bad day starts

## Automatic disqualifiers for a not-yet-earned option

A candidate option has not yet earned default promotion here if it mainly does
one or more of these:

- changes the control surface without naming which hidden burden moved
- improves deployment prestige while leaving wrong-node meaning weak
- improves controller power while leaving stateful authority socially held
- reduces local toil while preserving private topology reconstruction on the
  bad day
- sounds more adult mainly because the worldview got larger

That does not make the candidate worthless.
It means the candidate has not yet answered the repo's real benchmark strongly
enough to deserve default status.

This is where many recommendation documents fail the user.
They treat "worth considering" and "has earned default gravity" as nearly the
same thing.
This repo has to keep those very far apart.

## What still does not count as an orchestration decision here

This page should also say more directly what fake decisiveness looks like.

The following still do not count as a serious architecture decision:

- concluding that a bigger platform is inevitable because the docs now describe
  many moving parts
- treating helper sprawl as proof that a scheduler has already earned itself
- preferring the option with the broadest ecosystem without naming the exact
  hidden burden it removes
- picking the most inspectable option while skipping whether it actually owns
  enough truth on the bad day
- narrating "several good candidates remain" when most candidates still leave
  the same private burden in different packaging

This matters because the repo is not suffering from lack of reputable product
families.
It is suffering from lack of options that survive the user's actual benchmark.

The benchmark is not "which one scales better?"
It is not "which one has the biggest community?"
It is not even "which one is most likely to be broadly correct in general?"

It is much more personal and much more brutal:

> which one actually stops this whole topic from repeatedly collapsing back
> into my own private burden when a request lands on the wrong node or a local
> dependency dies?

## What a default-promotion packet would have to contain

Before any extra layer earns default status here, the docs should be able to
point to a concrete promotion packet.

That packet should include:

- the named burden being transferred
- the truth surface the candidate now owns
- the artifact carrying that truth
- the drill or failure condition that shows the burden moved
- the operator-visible inspection path that keeps the new layer legible
- the explicit statement of what still remains outside the candidate's reach

Examples:

- if a shared-truth helper layer wins, show where placement and peer eligibility
  live and how the edge consumes them
- if a service-supervision layer wins, show which fallback or takeover behavior
  became less private under failure
- if a full scheduler wins, show which narrower layer failed first and why the
  worldview tax is now justified by named pain instead of prestige

Without that packet, the page is still just warming the reader toward a product
family.

This packet requirement protects against one of the oldest infrastructure
failures in the world: confusing recommendation with evidence.
The user does not need another persuasive walk through the product landscape.
The user needs one reason to believe a new layer earned the right to exist in
their life specifically.

## The shortest honest answer

The current default stance is still:

> stay Compose-first, close the missing truth layers as narrowly as possible,
> and delay whole-stack orchestrator promotion until a specific domain clearly
> proves it needs a stronger control plane

That is not hesitation.
It is the most evidence-aligned answer available right now.

It is also the least insulting answer available right now.

Too much infrastructure advice effectively says:

- accept brittle manual glue
- or surrender to a giant control plane
- and stop pretending there should have been anything meaningful in between

This repo exists because the user does not accept that as the only adult
decision surface.

That refusal is not stubbornness.
It is the central philosophical pressure of the project.
The repo is trying to defend the idea that there should be a meaningful
middle ground between:

- brittle private glue
- and total worldview surrender

If no middle ground can survive honest proof, then the repo should say that.
But it should only say it after the honest proof, not before.

## The candidate families the repo is actually circling

The repo is not comparing infinite possibilities.
It keeps returning to a small number of families that each promise to own a
different slice of the missing truth.

### Family 1: narrower Compose-adjacent helper layers

Typical shapes:

- `services.yaml` or equivalent current-state registry
- sync-agent
- failover-agent
- file-generated Traefik dynamic config
- helper daemons that watch Docker, health, and git state

Why this family is attractive:

- preserves Compose as the main human control surface
- keeps the control plane legible
- targets the actual missing truths directly
- can be introduced incrementally

Why this family still fails easily:

- route generation can die with the backend it should route around
- peer eligibility can remain folklore in code form
- convergence of secrets, env, and revision state can still be weak
- helper sprawl can quietly become an orchestrator in disguise

There is a subtler risk too:

- this family can feel emotionally correct because it preserves the user's
  desired ergonomics
- while still failing technically to move enough truth out of private memory

That is why this family remains attractive and suspicious at the same time.

Current repo evidence:

- the master plan explicitly names `services.yaml`, sync, secret convergence,
  and service failover gaps
- `docker-gen-failover` is already treated as directionally relevant but
  operationally untrustworthy
- archive pressure repeatedly searches for a middle layer that is narrower than
  Kubernetes but more truthful than static glue

Current verdict:

- still the default search space
- not yet a proven winner

That is probably the most honest present-tense answer in the whole page.
The repo keeps circling this family because it most directly respects the
user's dream.
The repo keeps refusing to crown it because respecting the dream is not yet
the same thing as proving the dream can carry the burden honestly.

### Family 2: stronger service-supervision and placement systems

Typical shapes:

- OpenSVC-style supervision
- Nomad-style placement and rescheduling
- Consul-adjacent service discovery

Why this family is attractive:

- better ownership of health and service state
- more native answers to reschedule and failover questions
- less need for bespoke helper logic

Why this family still has to earn itself:

- it can move truth out of the operator's head while also moving it into a
  darker, less inspectable control surface
- it can solve scheduling without preserving the desired operator surface
- it can add legitimate machinery before the repo has proven which domain
  actually needs that much machinery

This family is where the user starts getting offered something that might
actually move a real burden, but also risks becoming one more layer that the
operator has to trust before it has emotionally earned that trust.

Archive pressure this family must answer:

- `distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
  shows the user actively searching for peer-equal or narrow coordination,
  rather than merely for "mature enterprise tooling"
- `nomad-multi-node-failover__68765e45-1ec4-8006-9179-5ef176d7a90f.md`
  shows that a Swarm-like story only matters if it answers the real burden
  instead of just offering a familiar scheduler shape

Current verdict:

- meaningful candidates for named domains
- not yet justified as the repo-wide default answer

That is a narrower and more respectful answer than "probably Nomad/OpenSVC
eventually."
The repo is not trying to predict the winner for sport.
It is trying to avoid forcing a bigger worldview onto the operator before the
smaller world has been given an honest chance to prove where it actually
breaks.

### Family 3: full scheduler / cluster worldview

Typical shapes:

- k3s
- Kubernetes
- larger garden-style cluster promotion

Why this family is attractive:

- broad ownership of service placement and health
- powerful scheduling, rollout, and control-plane capabilities
- mature ecosystem around service classes and cluster operations

Why this family is dangerous in this repo:

- it can replace one hidden burden with another more opaque one
- it can destroy the directness the user is explicitly trying to keep
- it can tempt the docs into declaring closure because the platform name now
  sounds adult enough

That last risk matters more here than in most projects.
The user is explicitly reacting to an ecosystem where adulthood of tone too
often substitutes for adulthood of burden ownership.

Current repo evidence:

- the repo has real k3s and Garden exploration
- the archive includes repeated curiosity about k3s, k8s, and cluster tools
- the instruction surfaces still preserve a no-heavy-orchestrator-by-default
  stance

Current verdict:

- absolutely relevant as a future promotion path
- still not the default answer the repo has earned today

This family may eventually win.
If it does, it should win because narrower layers were honestly pressed until
they failed a named burden threshold, not because Kubernetes or k3s merely
made the whole story sound complete enough to stop asking harder questions.

## The actual burden checklist a candidate must beat

A candidate layer only becomes a real option for this repo if it improves one
or more of these burdens materially and inspectably:

1. placement truth
   Can a receiving node answer "what runs where right now?" from shared truth?

2. peer eligibility truth
   Can it know which peer is healthy and semantically valid now?

3. route persistence under failure
   Does the rescue path survive the exact failure that makes rescue necessary?

4. convergence truth
   Are secrets, env, and deployment shape close enough that peer substitution
   is meaningful rather than accidental?

5. policy continuity
   Do auth, middleware, headers, and visible service meaning survive handoff?

6. stateful honesty
   Does the candidate improve authority, promotion, and write-path truth for
   stateful workloads instead of only improving ingress?

7. operator legibility
   Can an operator still explain what happened from tracked shared truth rather
   than private reconstruction?

If a candidate mostly helps one of those while worsening several others, it has
not earned default promotion yet.

This checklist is intentionally closer to a betrayal test than a feature test.
The repo already has enough feature language.
What it needs is a way to detect when a supposedly better option still leaves
the user stranded in the same explanatory role.

## What promotion evidence would actually look like

For this repo, a candidate layer earns promotion only when the docs can point
to a proof packet instead of a persuasive comparison.

That proof packet should eventually include things like:

- the named burden the candidate was supposed to remove
- the exact route, service class, or coordination domain where that burden was
  tested
- the shared truth surface the candidate now owns
- the failure or wrong-node condition that was exercised
- the post-failure or post-handoff result
- the explicit limits on what broader classes were **not** proven yet

Examples:

- if the candidate claims better placement truth, show the placement surface it
  owns and how runtime decisions consume it
- if it claims better peer eligibility, show why the chosen peer was valid and
  why another was not
- if it claims stateful progress, show authority, promotion, and client
  rediscovery truth rather than only route continuity

Without that structure, this page is too easy to overread as "the repo is
getting warmer toward product X."

## Why the archive makes this page harsher

This page has to stay harsher than a normal architecture-options page because
the archive keeps recording the same disappointment:

- static glue is too brittle
- bigger platforms are too eager to demand trust
- many answers sound dynamic until the user asks where the actual truth lives

That is why "respectable ecosystem answer" is not enough here.

The user is not merely trying to avoid complexity.
They are trying to avoid bad complexity:

- complexity that steals legibility
- complexity that relocates rather than removes burden
- complexity that sounds safer than it really is

That is also why this page must keep a visible distinction between:

- a candidate that owns more truth
- and a candidate that merely owns more machinery

The first might deserve promotion.
The second only deserves caution.

That distinction is the emotional center of the page.
The user is not allergic to power.
The user is allergic to paying more worldview tax for machinery that still
fails to answer the real burden question honestly.

## Bottom line

No broader orchestration family has yet earned default promotion for the repo
as a whole.

The current evidence still supports staying Compose-first and forcing any
stronger layer to justify itself by owning a named hidden burden better than
the current stack, not by sounding more adult, more modern, or more complete.

That is the real orchestration answer in `bolabaden-infra` today.

It is also the repo's way of refusing a false consolation:

- that the problem is basically solved once a bigger platform is mentioned
- that the missing options were imaginary
- or that the user's frustration mainly came from not having heard enough
  mature product names yet

This page should keep making that consolation unavailable.
