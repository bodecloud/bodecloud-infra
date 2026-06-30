# The Missing Middle Layer

This page answers the hardest architecture question in the repo:

> if plain Docker Compose leaves too much decisive truth inside the operator's
> head, but immediate surrender to Swarm, Nomad, OpenSVC, k3s, or Kubernetes
> still feels premature, what is the smaller added layer that could honestly
> remove the wound?

That smaller added layer is what this repo calls the missing middle.

It is not a brand.
It is not a product slot.
It is not shorthand for "lightweight orchestrator."

It is the first extra layer that can make one more bad-day truth genuinely
system-owned instead of privately remembered.

## The user's real dream

The user is not merely asking for:

- cleaner Compose files
- prettier dashboards
- one more reverse proxy trick
- "some HA"
- a lighter Kubernetes

The real dream is harsher and more specific:

- keep the human control surface close to ordinary Docker and Compose
- keep nodes feeling like ordinary machines rather than sacred cluster pets
- avoid premature lock-in to a heavyweight control-plane worldview
- stop one node from being secretly special
- stop one human from being secretly special
- let requests arrive on the wrong node without instantly turning correctness
  into folklore
- let failover mean more than "I can imagine how it might work"
- preserve the ability to reason about the system directly from the repo and
  runtime instead of trusting an opaque platform story

The middle layer exists only because plain Compose does not satisfy that dream
once the platform wants real multi-node dignity.

## What this page is and is not allowed to prove

This page is authoritative about:

- what the repo is actually searching for when it says "middle layer"
- which hidden burdens the added layer must remove to count
- what kinds of candidate families keep reappearing in the archive
- why many respectable-looking answers still fail this repo's benchmark
- how to tell the difference between helpful glue and true burden transfer

This page is not authoritative about:

- claiming the winner has already been found
- claiming the live runtime already owns these truths
- treating "smaller than Kubernetes" as success
- mistaking clearer explanation for completed architecture

This is a doctrine page for evaluation.
It is not a completion certificate.

## The shortest honest definition

The missing middle is:

> the smallest added layer that causes one more decisive multi-node truth to
> stop living only in private operator memory

That means the middle is **not** defined by:

- smaller binary count
- simpler marketing
- friendlier UI
- "lighter" branding
- being halfway between Compose and Kubernetes on a comparison chart

It is defined by whether the system can now carry one more sentence that the
operator previously had to finish by hand.

## The exact wound the middle layer must remove

Without a real middle layer, the user is still forced to be all of the
following at once:

- the hidden service registry
- the hidden placement ledger
- the hidden failover brain
- the hidden peer eligibility judge
- the hidden routing explainer
- the hidden drift reconciler
- the hidden memory of which nodes are equal in theory but not equal in
  practice

That is the wound.

The wound is not "manual work" in the abstract.
It is not "insufficient automation."
It is not "lack of enterprise features."

It is this:

> the platform still needs a human to privately remember what truth really
> governs request handling after the healthy happy-path story stops being
> enough

## The accusation the platform must survive

Every candidate middle layer must survive the user's actual accusation:

> you still only know this works because you privately know which node is
> special, which peer is safe, which backend is current, and which fallback is
> fake

If the accusation remains true after adopting a candidate layer, then the
middle is still missing.

## Why "middle" does not mean medium complexity

The easiest way to misunderstand this page is to think "middle" means:

- medium complexity
- medium scope
- medium opinionation
- medium amount of clustering

That is not what it means here.

The repo's meaning is sharper:

- small enough that Compose can still remain a meaningful human control surface
- strong enough that the operator is no longer the final keeper of crucial
  topology truth

That is why many small helpers still fail:

- they reduce repetition
- they improve config generation
- they improve route expression
- they reduce local toil

But they still leave the operator carrying the decisive sentence privately.

That is also why larger systems may eventually pass:

- they cost more
- they impose more worldview
- they hide more machinery

But if they truly move the decisive truth into the runtime, then they may have
earned that cost.

## The truths that currently have nowhere honest to live

The middle layer is missing because the platform keeps needing truths like
these and does not yet have a clean, current, inspectable place to hold them:

| Truth the runtime needs | Why it matters on the bad day | Why plain Compose does not own it cleanly |
| --- | --- | --- |
| Current placement truth | the receiving node must know where the service really lives now | Compose defines desired containers per host, not a shared live placement map |
| Current peer eligibility truth | reachable is weaker than safe, current, and semantically valid | local health checks do not define cross-node route suitability |
| Current fallback validity truth | rescue routes only matter if they survive the failure event | rendered config is weaker than post-failure route reality |
| Current route-class truth | stateless HTTP, protected HTTP, TCP, and stateful paths need different semantics | Compose can expose them all, but not classify substitution legitimacy |
| Current policy continuity truth | auth, middleware, and trust boundaries must survive handoff | local route definitions do not prove remote semantic equivalence |
| Current explanation truth | operators need inspectable reasoning after the fact | folklore often stands in for a runtime-owned decision trail |
| Current authority truth for stateful services | state needs more than ingress exposure and reachability | Compose does not supply distributed authority or recovery meaning |

These are not theoretical niceties.
They are exactly the truths that turn a platform from "looks clustered" into
"can survive embarrassment honestly."

## What the repo already has without the middle

The repo already has real ingredients:

- a substantial Compose-first runtime
- a non-trivial Traefik edge surface
- CrowdSec, TinyAuth, nginx-auth, and protected route semantics
- Headscale as a real private mesh assumption
- planning pressure toward tracked placement surfaces such as `services.yaml`
- research pressure from Compose-only, Nomad, OpenSVC, k3s, and other paths

Those are not fake.

But they are still ingredients, not proof that the runtime owns the right
multi-node truths today.

## The private sentences this repo is trying to kill

The missing middle is easier to understand if stated as the private sentences
the repo wants to eliminate.

Examples:

- "That route only works because I privately know node3 hosts the real backend
  today."
- "The proxy can reach node2, but I privately know node2 is the only one with
  the safe service revision."
- "The fallback file exists, but I privately know it disappears during the
  exact backend-loss event it claims to solve."
- "The service is visible on multiple nodes, but I privately know only one
  node is authoritative."
- "The protected route still responds, but I privately know the auth semantics
  are no longer the same."
- "Redis is reachable remotely, but I privately know that does not mean there
  is real failover."

If an added layer cannot kill at least one sentence like that cleanly, it has
not become the middle this repo is searching for.

## Candidate families the repo keeps circling

The repo is not short on options.
It is short on one that has fully earned belief.

### 1. Better static glue plus stronger local proxy expression

Examples:

- richer Traefik labels and file fragments
- helper-generated fallback config
- more templating, includes, and render-time composition

What this family can improve:

- local readability
- local route expression
- local operator convenience
- shape validation

What it usually cannot own by itself:

- shared current placement truth
- shared peer eligibility truth
- post-failure route persistence truth
- remote semantic continuity

This family often makes the platform easier to describe.
That is not the same thing as making the platform own one more decisive truth.

### 2. Lightweight registry plus Compose-first routing

Examples:

- tracked `services.yaml`
- file-backed placement maps
- small sync agents
- mesh-distributed backend identity updates

What this family is trying to do:

- answer "where does this service actually live right now?"
- keep Compose as the main authoring surface
- give every receiving node shared placement truth

Why this family is emotionally aligned with the repo:

- it attacks the user's real wound directly
- it keeps the system close to ordinary Docker
- it avoids immediate worldview surrender

What it still must prove before belief is earned:

- currentness under drift and churn
- explicit disagreement handling
- peer eligibility beyond simple reachability
- route survival during actual backend loss
- clear limits for protected and stateful classes

This family is the most obvious middle-layer candidate.
It is also the easiest place to fake progress with convincing glue.

### 3. Gossip and event-driven coordination

Examples from the archive:

- Serf-like membership signals
- peer-equal agents reacting to node events
- mesh-wide distribution of liveness or service updates

What this family can help with:

- membership visibility
- failure detection
- peer-to-peer signal flow
- preserving the emotional appeal of node equality

What it does not provide for free:

- conflict-free authority
- strong current placement truth
- route eligibility semantics
- policy continuity guarantees
- stateful authority meaning

This family is attractive because it sounds like the peer-equal dream.
It becomes dangerous when "every node hears the event" gets mistaken for
"every node now owns the right truth."

### 4. Registry plus dynamic proxy or control-plane helpers

Examples:

- Consul-like registries
- dynamic HAProxy or Envoy inputs
- helper APIs publishing backend identity or policy state

What this family can plausibly improve:

- service-name to backend resolution
- explicit runtime-fed routing
- inspectable resolution inputs

What it may still risk:

- new hidden sacred components
- stronger central trust concentration
- policy drift between nodes
- stateful ambiguity if ingress success gets over-interpreted

This can become a real middle layer if the receiving node can explain its
choice from shared truth instead of folklore.

### 5. Stronger orchestrator or cluster-control families

Examples:

- Nomad
- OpenSVC
- k3s
- Kubernetes

What they may genuinely own better:

- scheduling and rescheduling
- cluster state
- service discovery
- health-aware relocation
- stronger failover mechanics

What they charge:

- larger worldview
- more abstraction distance from Compose
- more hidden machinery
- more chances to confuse "controller exists" with "truth moved"

These systems are not disqualified.
They are simply too expensive to promote early unless the smaller candidates
fail to remove the wound honestly.

## What the archive pressure keeps saying

The archive keeps converging on the same point:

- Compose can describe many single-node truths very well
- multi-A DNS can help first-hop plurality
- reverse proxies can make edge behavior richer
- helpers can make shapes cleaner

But the real wound reappears as soon as the platform asks:

- where does the service live now?
- which peer is actually safe now?
- what still remains valid after the preferred backend dies?
- what semantics survive the handoff?
- who owns state authority now?

That means the missing middle is not a terminology gap.
It is a truth-ownership gap.

## What a candidate must prove to count as the middle

A candidate starts earning the name only if it can answer all of these with
runtime-owned truth rather than private operator reconstruction:

1. Can the receiving node determine local versus remote ownership from shared
   current truth?
2. Can it choose a peer from eligibility truth rather than reachability alone?
3. Can it explain why the chosen route remains valid after handoff?
4. Can it preserve protected-route meaning rather than only transport?
5. Can it keep the fallback alive through the failure event that justified it?
6. Can it expose disagreement, staleness, or uncertainty honestly instead of
   bluffing certainty?
7. Can it state where its authority ends for TCP and stateful classes?
8. Can it do all of that without imposing more worldview tax than the burden
   transfer justifies?

If the answer is still "not yet," then the middle is still missing.

## What still does not count as finding the middle

The following are real progress, but they do not yet count as having found the
middle:

- identifying several promising product families
- finding something smaller than Kubernetes
- making the repo feel more cluster-like
- generating more dynamic config
- improving route legibility
- building helper glue that still depends on operator memory
- becoming better at describing the wound

All of that can be useful.
None of it proves the decisive truth left the operator.

## Bottom line

The missing middle is not "the nicer orchestrator."

It is the first added layer that can survive this accusation:

> did one more crucial multi-node truth stop living only in the operator?

Until the answer becomes yes, the middle layer is still missing no matter how
complete, modern, or coherent the surrounding explanation sounds.
