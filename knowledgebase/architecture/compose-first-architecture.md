# Compose-First Architecture

This page answers the question that keeps returning in different disguises:

> if the dream is already bigger than ordinary single-node Docker Compose, why
> is the repo still centered on
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> instead of already surrendering to Nomad, k3s, Kubernetes, OpenSVC, or some
> other stronger control plane?

The answer is not:

- Compose is enough
- Compose solves distributed systems
- Compose is the forever platform

The answer is:

- Compose is still the clearest live human authoring contract in the repo
- the user wants distributed behavior without immediately surrendering
  inspectability
- the repo is trying to discover the smallest honest extra truth-owning layer
  before paying the worldview tax of a heavyweight orchestrator

That is what "Compose-first" means here.

## Why this page has to stay conditional

This page should not sound loyal to Compose.
It should sound conditional and slightly suspicious of it.

Otherwise "Compose-first" turns into an identity surface instead of what it is
supposed to be:

- a temporary truth-custody discipline
- a demand that the next larger surface actually earn itself
- a refusal to hide bad-day logic behind bigger nouns before the burden moves

If Compose stops being the least dishonest visible surface, this page should be
one of the first places willing to say so.

That exit-mindedness is more important than the slogan itself.

## What this page is and is not allowed to prove

This page is authoritative about:

- why Compose remains the primary human contract today
- why the repo has not yet promoted a heavyweight orchestrator by default
- which important truths the current Compose surface still does not own
- what would actually end Compose-first honestly

This page is not authoritative about:

- whether the current helper pile already closes those truths
- whether a future control layer has already earned promotion
- whether the current runtime already behaves like the dream

This page explains the authoring posture, not the final runtime verdict.

## Strongest honest current answer

Compose-first is still the least dishonest live authoring surface in the repo,
but only because the missing distributed truths are still more honestly visible
as missing there than they would be if the repo pretended a larger control
plane had already earned itself.

The important boundary is:

- Compose-first explains why the repo still leads with `docker-compose.yml`
- it does **not** prove that the runtime already owns wrong-node truth,
  peer eligibility, backend-loss durability, or stateful correctness

That distinction matters because the ecosystem keeps forcing a fake binary:

- either Compose is enough
- or the user should stop asking whether anything narrower than a full control
  plane could ever have been honest

This repo exists because that binary itself feels like evidence of missing
options.

The user is not merely comparing tools.
The user is reacting to an ecosystem that keeps shrinking the menu down to:

- keep the legible but inadequate thing
- or accept the larger hidden thing and stop asking whether there should have
  been a more truthful middle step

## What does not justify remaining Compose-first

This page needs to keep one bad defense illegal.

These still do **not** count as a serious justification for remaining
Compose-first:

- Compose is familiar
- Compose is readable on the happy path
- the repo already has a lot of Compose invested
- larger orchestrators feel annoying or ideological
- helper layers exist, so the gap must already be closing

Some of those may be emotionally or operationally true.
They still do not answer the only defensible question:

> is Compose still the least dishonest primary contract for the truths the
> runtime actually owns today?

If the answer becomes no, sentimentality is not a valid defense.

That warning matters because Compose can become its own fake adulthood too.
It can sound grounded, pragmatic, and readable while still leaving the same
bad-day burden untouched.

## What the live worktree actually proves

The current worktree proves several important things in Compose's favor.

### 1. The real runtime is still authored as Compose

The priority implementation still runs through:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- an include graph of active fragments
- Docker-native labels, networks, configs, secrets, and service definitions

That proves:

- the live operator surface is still file-first and Docker-native
- the runtime is still legible from the authored source of truth

### 2. The runtime is already broad enough that hiding it would be tempting

The root runtime already includes:

- real edge behavior in `compose/docker-compose.coolify-proxy.yml`
- real control-plane and mesh pressure in
  `compose/docker-compose.headscale.yml`
- real monitoring and alerting in `compose/docker-compose.metrics.yml`
- alternate orchestration exploration in old or side-path compose files
- helper layers around routing, auth, and maintenance

That proves the repo is not naively "just a few Compose files."

It also proves the opposite danger:

the repo is already broad enough that a shadow control plane could quietly grow
around Compose while humans continue pretending Compose is still where the
truth lives.

### 3. The missing truths are still honestly visible as missing

This is the strongest current argument for Compose-first.

The current runtime still visibly lacks several decisive truths:

- live current placement truth
- peer eligibility truth
- backend-loss route durability
- protected-route semantic continuity after peer-forward handoff
- stateful authority and promotion truth

Compose-first is still honest partly because those absences are still visible.

A larger control plane that only sounded more adult while leaving those same
burdens socially reconstructed would be less honest, not more.

## The shadow-control-plane threat

Compose-first is not only being tested against heavyweight orchestrators.

It is also being tested against a different threat:

> the possibility that helper layers become the real control plane while
> Compose remains only the decorative language humans still feel attached to

That threat has to stay visible.

Otherwise this page turns into a polite defense of the status quo instead of a
continuing test:

- is Compose still the least dishonest human contract?
- or has the actual bad-day intelligence already migrated somewhere murkier?

This is one of the repo's most important anti-romantic checks.

The user is not simply saying "I dislike complexity."
The user is saying:

> why do so many complexity increases still leave me holding the same
> explanatory burden, and why do so many simpler answers stop being honest as
> soon as distribution matters?

That is why Compose-first cannot be framed as a tasteful preference.

## What a real Compose-first retention packet would require

Before the docs use stronger "Compose-first is still right here" language, they
should be able to point to a concrete retention packet containing:

- the exact truth Compose still exposes better than the stronger alternative
- the exact hidden burden Compose still does not own
- the narrower helper or artifact carrying the missing burden today, if any
- the reason a heavier control plane still has not earned the worldview tax
- the explicit sentence describing what evidence would end Compose-first

Without that packet, Compose-first can quietly turn from honest restraint into
unexamined attachment.

The whole point is not to defend Compose culturally.
The point is to defend the last surface that is still more honest than the
larger answers surrounding it, and to stop defending it the moment that is no
longer true.

## What would actually end Compose-first

This page needs a real exit condition.

Compose-first should stop being the default the moment another surface can
prove, with current-worktree evidence plus drills, that it owns at least one
decisive missing truth more honestly than root Compose does.

In this repo, the decisive truths are things like:

- current placement truth
- peer eligibility truth
- backend-loss route durability
- protected-route semantic continuity under handoff
- stateful authority and promotion truth

That means Compose-first should end only after a challenger can show:

- the new control surface is not merely richer, farther away, or more famous
- operators can still inspect what it believes with less social translation
- one hidden human SPOF actually moved into a shared system-owned artifact
- the new surface stays honest on the bad day rather than only on the happy
  path

The user is not asking to keep Compose forever.
The user is asking for the next layer to earn the surrender it requests.

## What does not end Compose-first

These do **not** end Compose-first by themselves:

- standing up Nomad, k3s, Kubernetes, OpenSVC, or another controller in a lab
- generating richer route files
- adding more helper daemons around the edge
- making `services.yaml` sound central without proving the runtime consumes it
- making the architecture diagram look more cluster-like

Those may all be valid experiments.

They still do not answer the only serious replacement question:

> what exact burden does the new surface own better than readable Compose, and
> where is the evidence that this remains true under wrong-node or backend-loss
> stress?

If this page forgets that question, Compose-first collapses into branding and
its replacements collapse into prestige.

## What Compose-first means in this repo specifically

In `bolabaden-infra`, Compose-first means all of these are still true:

- the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  is the priority implementation surface
- active behavior is decomposed through included fragments in
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)
- the merged runtime should still be explainable from those files first
- helper layers should wrap, validate, generate from, or synchronize around
  that surface before replacing it

The root graph already proves this is not decorative rhetoric:

- it defines shared networks like `publicnet`, `backend`, and `warp-nat-net`
- it includes the major runtime fragments directly
- it already expresses heterogeneous workload classes directly
- it already exposes ingress, auth, observability, and stateful surfaces in a
  way humans can still inspect without controller-specific translation

That is a real advantage.

## The current bottom line

The strongest current sentence this page allows is:

> Compose-first is still the default because root Compose remains the clearest
> inspectable human contract for the current runtime, while the decisive missing
> distributed truths are still more honestly visible as missing there than they
> would be in a larger control plane that has not yet earned promotion

That sentence is intentionally conditional.

It defends Compose only as long as Compose remains the least dishonest place to
see what the system really knows and what it still does not know.
