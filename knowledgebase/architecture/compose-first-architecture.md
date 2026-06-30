# Compose-First Architecture

This page answers the question that keeps returning in different disguises:

> if the dream is already bigger than ordinary single-node Docker Compose, why
> is the repo still centered on
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> instead of already surrendering to Nomad, k3s, Kubernetes, or some other
> full control plane?

The answer is not:

- Compose is enough
- Compose solves distributed systems
- Compose is the forever platform

The answer is:

- Compose is still the clearest live authoring truth in the repo
- the user wants distributed behavior without prematurely surrendering
  readability
- the repo is trying to discover the smallest honest extra control surface
  before paying the tax of a heavyweight orchestrator

That is what "Compose-first" means here.

## What this page is and is not allowed to prove

This page is authoritative about:

- why Compose remains the primary human contract
- why the repo has not yet promoted a heavyweight orchestrator by default
- which important truths the current Compose surface still does not own

This page is not authoritative about:

- whether the current helper stack already closes those truths
- whether a future control layer has already earned promotion
- whether the current runtime already behaves like the dream

This page explains the authoring posture, not the final runtime verdict.

It also has to explain why Compose-first is not nostalgia.
The repo is not clinging to Compose because the user is scared of more capable
systems.
It is staying Compose-first only as long as Compose remains the least
dishonest place where cause, effect, and configuration meaning are still
visible.

## Strongest honest current answer

Compose-first is still the least dishonest live authoring surface in the repo,
but only because the missing distributed truths are still more honestly visible
as missing than they would be if the repo pretended a bigger control plane had
already earned itself. The important boundary is that Compose-first explains why
the repo still leads with `docker-compose.yml`; it does not prove the current
runtime already owns wrong-node truth, peer eligibility, backend-loss
durability, or stateful correctness.

That distinction matters because the ecosystem keeps trying to force a fake
binary:

- either Compose is enough
- or the user should stop asking whether anything narrower than a full control
  plane could ever have been honest

This repo exists because that binary itself feels like evidence of missing
options.

## What still does not count as justifying Compose-first

This page needs to keep one bad defense illegal.

The following still do not count as a serious justification for remaining
Compose-first:

- Compose is familiar
- Compose is readable on the happy path
- the repo already has a lot of Compose invested
- larger orchestrators feel annoying or ideological
- helpers exist, so the gap must already be closing

Some of those may be emotionally or operationally true.
They still do not answer the only defensible question:

> is Compose still the least dishonest primary contract for the truths the
> runtime actually owns today?

If the answer becomes no, sentimentality is not a valid defense.

That sentence protects against a subtle failure mode:
Compose-first can sound pragmatic and grounded while still becoming just
another ideology if it stops asking whether the bad-day intelligence has
already migrated somewhere murkier.

## What a real Compose-first retention packet would have to contain

Before the docs use stronger "Compose-first is still right here" language, they
should be able to point to a concrete packet.

That packet should contain:

- the exact truth Compose still exposes better than the stronger alternative
- the exact hidden burden Compose still does not own
- the narrower helper or artifact carrying the missing burden today, if any
- the reason a heavier control plane still has not earned the worldview tax
- the explicit sentence describing what evidence would end Compose-first as the
  default

Without that packet, Compose-first can quietly turn from honest restraint into
unexamined attachment.

This page has to be strict about that because "Compose-first" is exactly the
kind of phrase that can start sounding adult, disciplined, and reasonable even
after it has stopped answering the user's real question.

## What Compose-first is being forced to justify

Compose-first is not only being defended against heavyweight orchestrators.
It is also being interrogated against a different threat:

> the possibility that the helper pile becomes a shadow control plane while
> Compose remains only the decorative language humans still feel attached to

That threat has to stay visible.

Otherwise this page becomes a polite defense of the status quo instead of a
continuing test:

- is Compose still the least dishonest human contract?
- or is it being preserved after the actual bad-day intelligence moved
  somewhere murkier?

That second question is one of the most important anti-romantic checks in the
whole knowledgebase.

It is also an empathy check.
The user is not simply saying "I dislike complexity."
They are saying "why do so many complexity increases still leave me holding
the same explanatory burden, and why do so many simpler answers stop being
honest as soon as distribution matters?"

## What Compose-first means in the priority implementation

In `bolabaden-infra`, Compose-first means all of these are still true:

- the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  is the priority implementation surface
- active behavior is decomposed through included fragments in
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- the merged runtime should still be explainable from those files first
- helper layers should wrap, validate, generate from, or synchronize around
  that surface before replacing it

The root graph already proves this is not decorative rhetoric:

- it defines the shared networks `publicnet`, `backend`, and `warp-nat-net`
- it includes the major runtime fragments rather than hiding them behind a
  remote controller API
- it already expresses heterogeneous workload classes directly

That matters because the user keeps rejecting systems that hide too much
meaning behind a farther-away reconciliation layer before that hiding has
earned its keep.

That rejection is not only taste.
It is a substantive architecture position:

> a farther-away control surface should not become the default answer until it
> owns a named truth better than the readable surface it is replacing

## What Compose-first definitely does not mean

It does not mean:

- root Compose already expresses live cross-node placement truth
- root Compose already proves wrong-node success
- root Compose already proves fallback durability under backend loss
- root Compose already proves protected-route semantic continuity
- root Compose already proves stateful resilience
- intended `services.yaml` behavior should be narrated as if it were already a
  shipped tracked root contract

That distinction matters because earlier docs often blurred together:

- desired runtime feeling
- actual implementation surface
- proof boundary

Once those blur, the docs immediately start sounding more complete than the
system actually is.

That is exactly the documentation habit the user is pushing back on across the
whole repo.

## Why Compose is still the least dishonest live surface

Compose is still useful here for practical reasons, not nostalgic ones.

### 1. It keeps authoring truth visible

Many important truths are still directly readable in the root graph:

- service definitions
- network attachments
- labels
- secrets references
- inline configs
- healthchecks
- restart policies
- mounts

That is not trivial.
It is one of the last sane surfaces the repo still has.

The user is reacting against an ecosystem where many "better" answers buy
coordination by making the real state harder to see until a new worldview has
already been accepted.

Compose does not make distribution easy.
It does keep current authoring truth close to the human eye.

That is why Compose-first still carries moral weight here.
Not because YAML is sacred, but because hidden state and hidden decisions have
to earn themselves instead of arriving pre-forgiven.

### 2. It handles heterogeneous workloads without inventing one giant worldview

The current root runtime already spans sharply different workload classes:

- public ingress and auth
- observability
- operator dashboards
- browser automation and crawling
- AI tooling
- mesh and identity surfaces
- raw TCP services
- state-bearing services

Compose is genuinely good at expressing this heterogeneity.

Its weakness is not workload diversity.
Its weakness is distributed truth ownership.

That distinction matters because it explains why the repo is not fleeing
Compose for authoring reasons.
It is searching beyond Compose for truth-ownership reasons.

### 3. It preserves Docker-native operator ergonomics

The repo still gets real value from being able to ask questions in a very
direct way:

- what is in the merged graph?
- which fragment introduced this service?
- what networks is it on?
- what auth or middleware labels exist?
- what port or backend is being routed locally?

That style of inspection still matters because the user does not want to lose
the causal chain between authored config and live behavior unless the extra
abstraction clearly buys something real.

### 4. It keeps the control-plane tax explicit

As long as Compose remains primary, any extra helper layer has to justify
itself against a readable baseline.

That is healthy.
It forces harder questions:

- what truth does this helper own?
- is that truth already implicit in the operator's head?
- does this helper remove a real hidden-human SPOF?
- or does it mostly create a more respectable-looking story?

This is one reason Compose-first is still useful even when it is not sufficient.
It keeps the cost of promotion visible.

Keeping the cost visible is part of how the repo refuses fake options.
If every larger layer arrived already coated in inevitability, the user would
lose the last meaningful place from which to ask whether the new tax actually
buys the right truth.

## The exact truths Compose does not naturally own

This page is pointless if it does not name the missing truths plainly.

Compose is good at authored local intent.
Compose does not naturally own:

- distributed placement truth
- convergence truth across nodes
- peer eligibility truth
- cross-node route persistence under backend loss
- wrong-node request preservation
- service-class-specific stateful authority

Those are not tiny caveats.
They are the heart of the user's complaint.

The missing truths are not peripheral gaps around an otherwise complete
platform.
They are the reason the platform still does not feel like it offers the kind
of believable option the user keeps searching for.

This is why Compose-first cannot be read as a quiet endorsement of the current
helper strategy.
It is better read as:

- Compose is still the least dishonest readable surface
- the helper and promotion story is still under cross-examination
- the wrong-node and backend-loss contracts are still the real judge

That is why Compose-first remains honest only while it remains conditional.
The moment it becomes identity rather than interrogation, it becomes another
smaller and safer story than the user asked for.

## Where the priority runtime already shows this tension

The current root runtime gives concrete examples of why Compose-first is both
valuable and insufficient.

### Example 1: serious edge execution, missing distributed truth

The edge stack is already real:

- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `cloudflare-ddns`
- `docker-gen-failover`

That proves the repo is not doing toy ingress work.
It also proves the danger:

a serious local routing and policy surface can make distributed truth feel
closer than it is.

### Example 2: real mesh and control-plane surfaces, still not shared placement truth

Headscale is materially live through
[`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml).

That means peer communication and identity are not hypothetical.
But it still does not prove the receiving node has a trustworthy shared answer
to:

- where does the target service live right now?
- which peer is eligible now?

So even a live mesh does not settle the missing middle layer by itself.

The ecosystem is full of steps that feel like progress while still leaving the
final judgment trapped in a human brain.
Compose-first is only worth preserving if it helps prevent those half-solutions
from being over-credited.

### Example 3: helper growth already pressures the boundary

The repo already contains or plans around helpers that smell like a partial
control plane:

- `docker-gen-failover`
- secret and compose sync concepts
- peer broadcast concepts
- `services.yaml` pressure
- failover-agent and sync-agent ideas in the master plan

This is exactly why Compose-first must stay interrogative rather than
celebratory.

If the real bad-day intelligence migrates into these helpers without becoming
more inspectable than a heavyweight controller, then Compose-first becomes a
decorative comfort blanket instead of a real operator contract.

That sentence should stay harsh.
The user is not asking to be comforted by readable syntax.
They are asking for a system where readable syntax still points at the place
where the truth actually lives.

## The real argument against premature promotion

The repo is not saying:

- Kubernetes is bad
- Nomad is wrong
- OpenSVC is cheating

The repo is saying:

> do not import a larger worldview before you can name the exact missing truth
> it would own and prove that smaller layers could not own it honestly

That is a much narrower and more defensible standard.

It also explains why "just use k3s" keeps reading as incomplete here.
It may eventually be right.
It is not self-justifying.

The user keeps demanding a harder question first:

> which exact truth layer is missing, and can that layer be added without
> immediately replacing Compose as the main surface the operator reads?

That question is the entire point of Compose-first.

If the repo ever loses sight of that, "Compose-first" will become a polished
but much less honest phrase.

## The anti-benchmark this page must keep visible

If the answer becomes:

> Compose is still the authoring surface, but the actual bad-day intelligence
> now lives in several helpers, some generated state, and operator intuition
> about which outputs matter

then Compose-first has failed on its own terms.

The user is not trying to preserve Compose as a cultural artifact.
They are trying to preserve a surface where the system still feels causally
owned.

If that ownership has already migrated somewhere murkier, the docs need to say
so instead of romanticizing Docker-native syntax.

## What would actually justify moving beyond Compose-first

The repo should promote a stronger control layer only when it can answer:

- which missing truth does this layer own directly?
- what hidden-human SPOF does it remove?
- what narrower alternatives were exhausted first?
- how does it stay more inspectable or more honest than the burden it replaces?

That means promotion is not blocked forever.
It is gated by honesty.

## The blunt reading

Compose-first is not the solution.
It is the current least-dishonest authoring stance while the repo is still
trying to identify the smallest extra truth-owning layer that would make
wrong-node requests, backend-loss recovery, and hidden topology memory stop
being the dominant failure mode.

That is why the repo is still centered on `docker-compose.yml`.
Not because Compose secretly solved the problem, but because the problem has
not yet earned a bigger answer honestly enough.

Another way to say the same thing:

- Compose-first is not the dream
- Compose-first is the refusal to pretend the dream has already justified a
  larger concealment layer
