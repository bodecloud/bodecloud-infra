# The Missing Middle Layer

This page answers the most important architecture question in the repo:

> if plain Compose becomes too dependent on private operator memory, and
> immediate promotion to Kubernetes, Swarm, Nomad, or another larger control
> plane feels too expensive or too premature, what is the smaller truth-owning
> layer the repo is actually searching for?

That is the missing middle.

It is not a product name.
It is not a prestige category.
It is not just "something lighter than Kubernetes."

It is a burden-transfer shape.

## What this page is and is not allowed to prove

This page is authoritative about:

- the capability shape the repo is repeatedly converging on
- which hidden burdens the middle layer must actually remove
- which candidate families are already visible in the repo
- why many respectable-looking answers still fail this repo's benchmark
- how to tell whether a helper, registry, or control plane is genuinely paying
  down the wound rather than renaming it

This page is not authoritative about:

- claiming the winner has already been found
- claiming the live runtime already owns this layer
- equating "smaller than Kubernetes" with "good enough"
- treating one candidate family as promoted just because its logic sounds
  coherent

This page describes the wanted layer.
It does not certify a final chosen implementation.

## The shortest honest definition

The missing middle is:

the smallest added layer that makes decisive bad-day truths system-owned
instead of privately remembered.

That means the middle is **not** defined by:

- size
- popularity
- how friendly the UI feels
- whether it is marketed as lightweight
- whether it sits "between Compose and Kubernetes" on a diagram

It is defined by whether one more humiliating private completion step actually
leaves the operator.

## The hidden duties that layer must remove

The user is trying to stop being all of the following at once:

- the hidden service registry
- the hidden failover brain
- the hidden drift detector
- the hidden routing explainer
- the hidden memory of what lives where

If a candidate layer does not actually move those duties into the system, then
it is not the middle this repo is searching for, no matter how elegant or
fashionable it sounds.

## Why "missing middle" is not just a complexity bracket

One of the easiest mistakes in this repo is to hear "missing middle" and think
it means:

- medium complexity
- medium size
- medium opinionation
- halfway between Compose and Kubernetes

That is not the real meaning here.

The real meaning is harsher:

- a layer small enough not to demand worldview surrender too early
- but strong enough to own truths that currently live only in human memory

That is why some apparently small helpers still fail the benchmark:

- they reduce repetition
- they improve generation
- they improve expression
- but they do not remove the decisive private explanation step

And it is why some larger systems may eventually pass:

- they centralize more
- they impose more worldview
- but they may truly own the right truths if they can prove the burden moved

## The truths that currently have nowhere honest to live

The missing middle exists because the repo keeps needing truths like these:

| Truth | Why the system needs it | Why plain Compose does not own it cleanly |
| --- | --- | --- |
| Current placement | the wrong node must know where the service actually lives now | Compose defines desired containers per host, not a shared live placement map |
| Current peer eligibility | not every reachable peer is safe for every route | local health checks do not define cross-node suitability |
| Current rescue-route validity | fallback only matters if it survives failure | rendered config is weaker than post-failure route truth |
| Current route-class meaning | HTTP, protected HTTP, TCP, and stateful paths need different semantics | Compose can expose all of them, but not classify their substitution rules |
| Current explanation for local vs remote choice | operators need inspectable reasoning after the fact | folklore often substitutes for operator-readable routing truth |

The middle layer is missing because those truths still have nowhere honest to
live.

## What the repo already has without that layer

The repo already has:

- a large real Compose runtime
- a serious Traefik, CrowdSec, TinyAuth, and nginx-auth-bearing edge layer
- Headscale as a real private-mesh assumption
- planning pressure toward `services.yaml`, peer sync, and failover helpers
- research into OpenSVC, Nomad, k3s, and related paths

Those are meaningful ingredients.
They are not yet the same thing as one promoted layer that cleanly owns the
decisive truths across the priority implementation.

## Candidate families the repo keeps circling

The repo is not short on candidate families.
It is short on one that has clearly earned itself.

### 1. Static glue plus better local proxies

Examples:

- richer Traefik or nginx expression
- helper-generated fallback config
- more labels, more includes, more templating

What this family can help with:

- cleaner route expression
- better local proxy behavior
- improved operator legibility

What it usually fails to solve:

- shared live placement truth
- cross-node peer eligibility truth
- inspectable wrong-node decision logic

### 2. Lightweight registry plus local proxying

Examples:

- a tracked `services.yaml`
- file-backed service maps
- small sync agents
- peer-broadcast updates over mesh links

What this family is trying to solve:

- "service name -> where is it running right now"
- giving every node a shared placement view
- keeping Compose as the main authoring layer

Why it is attractive here:

- it attacks the exact wound visible in the archive
- it may stay narrower than a scheduler
- it keeps the control surface close to Docker and Compose

What it still must prove:

- registry freshness and correctness under failure
- peer eligibility beyond simple reachability
- route durability after backend loss
- stateful truth beyond HTTP routing

### 3. Gossip and event-driven Compose glue

Examples from the archive:

- Serf-like membership and failure events
- node agents reacting to gossip events
- peer-equal failure detection without immediately adopting a full scheduler

What this family can help with:

- membership awareness
- node-failure signals
- peer-to-peer event distribution

What it still does not give for free:

- strong shared truth
- conflict-free decisions
- state authority
- safe failover semantics for complex protected or stateful routes

This family is attractive when the user wants "all nodes equal" and does not
want Swarm.
It is dangerous if eventual gossip gets mistaken for authoritative decision
truth.

### 4. Stronger orchestrator or cluster-control families

Examples:

- Nomad
- OpenSVC
- k3s or Kubernetes-derived paths

What they can potentially own better:

- scheduling and rescheduling
- service discovery
- stronger cluster state
- automatic failover mechanics

What they cost:

- larger control-plane worldview
- more hidden machinery
- more abstraction distance from the readable Compose surface

That cost is acceptable only if they remove a real hidden burden that the
smaller families cannot remove honestly.

## What the archive keeps confirming

The archive keeps converging on the same points:

- `docker-multi-node-without-swarm__...` frames service discovery as the hard
  missing piece once manual placement and DNS plurality are accepted
- `distributed-ha-orchestration__...` makes clear that fully peer-equal,
  leaderless orchestration usually requires custom glue rather than an
  off-the-shelf miracle
- `nomad-multi-node-failover__...` shows that stronger orchestrators can help,
  but they bring a larger worldview and still need proof against this repo's
  exact burden

That means the missing middle is not a naming failure.
It is a truth-ownership failure.

## Tests a candidate must pass to count as the middle

A candidate layer starts earning the name only if it can answer:

1. Can the receiving node determine locality or remote ownership from shared
   truth rather than memory?
2. Can it choose a peer from eligibility truth rather than reachability alone?
3. Can it explain why the route remained valid after handoff?
4. Can it preserve protected-route semantics rather than just transport?
5. Can it keep a fallback alive after the preferred backend disappears?
6. Can it do all of that without forcing premature surrender to a larger
   control-plane worldview?

If the answer is still "not yet," then the layer may be useful but it has not
yet become the missing middle this repo needs.

## What still does not count as finding the middle

The following still do not count as having found the missing middle:

- identifying several promising candidate families
- finding something smaller than Kubernetes
- finding something more dynamic than static Compose
- giving helpers nicer names for placement, sync, or failover
- building enough glue that the repo feels cluster-ish
- becoming much better at describing what the middle should do

All of that may be real progress.
None of it proves the operator stopped being the final keeper of truth.

## Bottom line

The middle layer is not "the nicer orchestrator."
It is the first added layer that can survive this accusation:

> did one more private topology truth actually stop living only in the
> operator?

Until the answer becomes yes, the missing middle is still missing, no matter
how complete the surrounding explanation sounds.
