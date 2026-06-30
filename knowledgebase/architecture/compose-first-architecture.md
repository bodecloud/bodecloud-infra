# Compose-First Architecture

This page is not here to defend Docker Compose emotionally.
It is here to answer a sharper question:

> if the repo clearly wants something bigger than one ordinary Docker host,
> why is the priority implementation still rooted in
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> instead of already giving up and calling the answer Swarm, Nomad, k3s,
> Kubernetes, or some other stronger control plane?

The answer is not:

- Compose is enough
- Compose solves distributed systems
- Compose is modern enough to not need anything else
- bigger orchestrators are bad on principle

The current honest answer is narrower:

- Compose is still the clearest live authoring surface in the repo
- the user wants multi-node behavior without immediately surrendering legibility
- the repo is still searching for the smallest truthful middle layer between
  single-node Compose and a heavyweight scheduler
- the missing truths are still more honestly visible as missing in Compose than
  they would be inside a larger, more flattering control surface

That is what "Compose-first" is supposed to mean here.

## The real accusation behind Compose-first

The user is not simply saying:

> I like Compose better than Kubernetes.

The stronger accusation is closer to:

> why do so many upgrades in infrastructure complexity still leave me holding
> the same bad-day explanatory burden?

That burden is the heart of the repo.
The user is looking for a system where several ordinary Docker nodes stop
behaving stupidly under:

- wrong-node entry
- backend loss
- peer selection
- middleware continuity
- stateful authority pressure

without immediately accepting a larger opaque system that merely hides the same
human SPOF behind more machinery.

So Compose-first is not nostalgia.
It is suspicion.

## What this page is and is not allowed to prove

This page is allowed to prove:

- why Compose remains the primary authoring contract today
- why the repo has not yet promoted a heavyweight orchestrator by default
- which truths the current Compose surface clearly does not own
- what would have to change before Compose-first stops being honest

This page is not allowed to prove:

- that the current helper stack already closes the distributed gap
- that a future middle layer has already earned promotion
- that the live runtime already behaves like the dream
- that "Compose-first" is a satisfying answer by itself

This page is about custody of truth, not completion.

## The strongest current reason Compose still survives

Compose survives here for one hard reason:

> it is still the least dishonest primary surface for the runtime the repo can
> currently show.

That sentence is intentionally conditional.

The worktree proves the live runtime is still authored through:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- Docker-native labels, configs, secrets, networks, and service definitions

That matters because it means the operator can still inspect the main runtime
surface by reading authored files instead of reverse-engineering a remote
controller's worldview.

But that alone would not justify Compose-first.
The stronger reason is that the missing distributed truths are still plainly
missing there.

The current worktree still does not clearly show a system-owned answer for:

- current placement truth
- peer-eligibility truth
- backend-loss route durability
- protected-route semantic continuity under handoff
- stateful authority and promotion truth

Compose is still honest partly because those absences are not hidden.

## Compose-first is not a compliment

This page should never make Compose-first sound like a polished architecture
destination.

It is closer to a temporary restraint:

- do not promote a stronger control plane until it really owns a missing truth
- do not hide distributed uncertainty behind richer nouns
- do not let helper sprawl quietly become the real control plane while humans
  keep pretending the truth still lives in Compose

That last point matters.
The biggest threat is not only "eventually Kubernetes."
It is also:

> a shadow control plane grows around Compose, helpers become the real decision
> surface, and Compose remains only the decorative language humans still feel
> attached to.

If that has happened, Compose-first is already lying.

## What does not justify staying Compose-first

These are weak excuses and should stay illegal:

- Compose is familiar
- Compose is readable on the happy path
- the repo already invested heavily in Compose
- larger orchestrators feel annoying
- the helper stack is getting sophisticated, so the gap must be closing

Some of those may be operationally true.
They still do not answer the only question that matters:

> is Compose still the least dishonest primary contract for the truths the
> runtime really owns today?

If the answer becomes no, emotional attachment is not a defense.

## What the live worktree proves in Compose's favor

The current runtime gives Compose a few real advantages.

### 1. The live implementation is still genuinely Compose-authored

This is not a repo that merely talks Compose while secretly living somewhere
else.

The authoring priority still centers on:

- root `docker-compose.yml`
- include fragments under `compose/`
- labels and attached runtime metadata for routing, auth, health, and service
  exposure

That means the visible authored surface is still the place where a human can
inspect what the repo thinks exists.

### 2. The runtime is broad enough that hiding it would be tempting

The repo is not a toy stack.
The active runtime already includes:

- real edge behavior in `compose/docker-compose.coolify-proxy.yml`
- real private-network and reachability pressure in
  `compose/docker-compose.headscale.yml`
- real metrics and alerting in `compose/docker-compose.metrics.yml`
- stateful application infrastructure in the core and LLM fragments
- helper logic around routing, auth, DDNS, and fallback generation

That matters because the repo is already complex enough that a fake adulthood
story would be easy to tell.

Compose-first is only useful if it keeps that temptation visible instead of
giving it cover.

### 3. The missing truths are still easier to see than to hand-wave

The strongest Compose-first argument is still this:

> the missing distributed truths are easier to point at directly in the current
> Compose surface than they would be inside a larger system that sounded more
> finished than it really was.

In other words:

- Compose does not solve the problem
- Compose still exposes the wound more honestly than a premature replacement

That is not glamorous, but it is real.

## The hidden test: where does the bad-day intelligence live?

The real Compose-first test is not:

> can the stack be described from Compose files?

It is:

> when the bad day begins, does the important decision still live in one
> operator's head, in helper folklore, or in a shared inspectable artifact?

The user is explicitly trying to escape the first two.

If a stronger layer cannot move one decisive bad-day sentence into a shared
system-owned surface, then it has not earned the surrender it asks for.

## The specific middle layer the repo is still hunting for

The dream here is not "stay on Compose forever."
The dream is:

> discover whether there is a narrower middle layer that can carry current
> placement, peer choice, and failover truth honestly before paying the full
> worldview tax of a heavyweight orchestrator.

That is why ideas like `services.yaml`, route generation, dynamic failover
helpers, and peer-aware edge logic keep returning.

The repo is not wandering randomly.
It keeps circling the same missing sentence:

> if the request lands on the wrong healthy node, how does that node know what
> to do next from shared truth rather than private operator memory?

Compose-first remains alive because that middle answer is still unearned, not
because the question went away.

## The retention packet this page should demand

Any future defense of Compose-first should be able to state a small packet:

1. what exact truth Compose still exposes better than the challenger
2. what exact burden Compose still fails to own
3. what smaller helper or artifact carries that missing burden today, if any
4. why a heavier control plane still has not earned its worldview tax
5. what exact evidence would end Compose-first

Without that packet, "Compose-first" becomes style instead of analysis.

## What would actually end Compose-first

Compose-first should end the moment another surface can prove, with current
worktree evidence plus drills, that it owns at least one decisive missing truth
more honestly than Compose does.

In this repo, decisive truths include:

- current placement truth
- peer-eligibility truth
- backend-loss route durability
- protected-route semantic continuity under handoff
- stateful authority and promotion truth

A challenger only earns promotion if it can show all of the following for at
least one of those truths:

- the new surface is not merely richer or more famous
- operators can inspect what it believes without more folklore
- one hidden human SPOF really moved into a shared artifact or controller-owned
  truth
- the new surface stays honest on the bad day, not only on the diagram

That is the standard.

## What does not end Compose-first

These do not end it by themselves:

- running Nomad, k3s, Kubernetes, OpenSVC, or another controller in a lab
- generating richer route files
- adding more helper daemons around the edge
- talking as if `services.yaml` is central without proving the live runtime
  consumes it
- making the architecture diagram look more cluster-like

Those may all be useful experiments.
They still do not answer the only serious replacement question:

> what exact burden moved out of the operator's head, and where is the proof?

## The honest bottom line

Compose-first is still the least dishonest visible authoring surface in
`bolabaden-infra`, but that is not praise.

It means:

- the dream is larger than Compose
- the heavier alternatives have not yet earned themselves
- the helper pile has not yet proven it became a truthful shared control layer
- the user is still right to feel that the ecosystem keeps forcing an ugly
  binary between legible inadequacy and opaque adulthood

So the honest sentence today is:

> the repo still leads with Compose not because Compose solved the distributed
> problem, but because the repo still has not proven a more truthful primary
> contract for the missing middle layer the user actually wants.
