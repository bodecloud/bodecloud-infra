# The Missing Middle Layer

This page answers the most important architecture question in the repo:

> if plain Compose becomes too dependent on private operator memory, and
> immediate promotion to Kubernetes, Swarm, Nomad, or another larger control
> plane feels too expensive or too premature, what is the smaller truth-owning
> layer the repo is actually searching for?

That is the missing middle.

It is not a product name.
It is not a prestige category.
It is not "lighter than Kubernetes" by itself.

It is a burden-transfer shape.

This page only makes sense if the burden is stated brutally:

the user is trying to stop being the hidden registry, hidden failover brain,
hidden drift detector, hidden routing explainer, and hidden memory of what
lives where.

If a candidate layer does not actually move those duties into the system, then
it is not the middle this repo is searching for, no matter how elegant or
fashionable it sounds.

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

## Strongest honest current answer

The repo is looking for the smallest added layer that makes these truths
system-owned instead of remembered:

- where a service lives now
- whether a peer is eligible now
- whether the rescue route still exists now
- whether candidate peers are converged enough to substitute safely
- why the system made the routing or failover decision it made

The repo already has:

- a large real Compose runtime
- a serious Traefik, CrowdSec, and auth-bearing edge surface
- Headscale as a real private-mesh assumption
- planning pressure toward `services.yaml`, peer sync, and failover helpers
- research into OpenSVC, Nomad, k3s, and related paths

What it still does **not** have is one promoted layer that cleanly owns those
truths across the priority implementation.

That is why so many surrounding ideas still feel partial:

- they improve first-hop reachability
- or they improve routing expression
- or they improve service ownership
- or they improve scheduling semantics

but they do not necessarily remove the exact hidden burden the user is
complaining about.

## What still does not count as finding the middle

This page needs to reject a very specific kind of fake progress.

The following still do not count as having found the missing middle for this
repo:

- identifying several promising candidate families
- finding something smaller than Kubernetes
- finding something more dynamic than static Compose
- giving helpers nicer names for placement, sync, or failover
- building enough glue that the repo feels "cluster-ish"

Those things may all be relevant.
They still do not answer the only question that matters:

> which hidden burden stops living in the operator after this layer exists?

If that answer is still vague, the repo has found more machinery, not the
middle.

## What a real missing-middle proof packet would have to contain

Before this page supports a stronger claim like "the middle is becoming real,"
it should point to a concrete packet.

That packet should contain:

- the exact hidden burden being transferred
- the artifact or runtime surface that now owns that burden
- the route or service class where the transfer becomes visible
- the failure condition or decision point that proves the burden moved
- the explicit sentence naming which adjacent burdens still remain private

Without that packet, the docs are still mostly describing an attractive shape
instead of a transferred responsibility.

## Why this layer is needed at all

The repo does not need a missing middle because clustering is fashionable.
It needs one because several critical facts are still too easy to answer only
from private human reconstruction.

Those facts include:

- what runs where right now
- which peer is the right candidate for this service right now
- whether that peer is merely reachable or actually safe to receive traffic
- whether the route needed for rescue survives backend disappearance
- whether auth and middleware keep the same meaning during handoff
- whether a stateful service is genuinely movable or merely reachable

As long as those answers remain mostly social, the operator is still the
effective control plane.

That is the thing this page is trying to kill.

## Why DNS and proxy sophistication are still not enough

The repo already has Cloudflare participation and a strong live edge stack.
That still does not satisfy the missing middle because:

- Cloudflare can choose a healthy first-hop node, but not the correct backend
  placement inside the user's topology
- Traefik can express rich routing and middleware locally, but not magically
  infer trustworthy cross-node placement truth
- healthchecks can detect some failures, but not automatically answer whether a
  remote peer is semantically equivalent for this traffic
- richer edge behavior can make the platform sound more complete than its
  distributed truth actually is

This is why the repo keeps separating:

- first-hop plurality
- request preservation
- policy continuity
- stateful authority

Those are different truths.
The middle layer is needed because the current stack can already express some
of them without fully owning the rest.

## The repo has already named the gap directly

The clearest planning pressure is in
[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md),
which explicitly records these as still missing:

- universal wrong-node success
- a live tracked root `services.yaml` current-state registry
- trustworthy route persistence under local backend failure
- automated service failover between nodes

Those are not side quests.
They are the exact truths the missing middle would need to own.

The master plan also records:

- `docker-gen-failover` is broken because it can delete routes when containers
  stop
- Cloudflare DDNS presence is not the same thing as full failover
- Headscale is single-node today
- secret sync and Compose sync are still manual

That list is basically a proof that the middle layer has not been earned yet.

## The burden-transfer test

This is the most important test on the page:

before calling anything the missing middle, ask:

> which specific hidden explanation will the system own after this exists that
> the operator currently has to supply from memory?

If the answer is weak, then the thing may still be useful.
It is not the missing middle for this repo.

That means candidate layers should be judged less by:

- popularity
- modernity
- how small or large they sound
- how impressive the diagrams are

and more by:

- which hidden burden they concretely remove
- which truths become inspectable and shared
- which truths remain private, stale, or guess-driven
- whether operator legibility improves or decays

## The easiest way this page can still fail

The most common failure mode here is subtle:

the page starts sounding like it has already narrowed the solution to a small
respectable cluster of tool families, and that feeling of narrowed possibility
gets mistaken for relief.

That is still too soft for this repo.

The user is not mainly asking for a curated shortlist.
The user is asking for a system that stops quietly requiring remembered rescue
knowledge on the bad day.

If the shortlist still leaves that burden mostly intact, this page has
organized the wound without treating it.

## What a real middle answer would have to leave behind in the repo

The user is not only asking for a better idea.
They are asking for a system that leaves behind clearer runtime evidence.

So a real middle answer should eventually leave behind artifacts like:

- a tracked placement-truth surface such as `services.yaml` or an equally
  inspectable generated state surface
- an explicit peer-eligibility surface showing why a node is eligible, not just
  that it answered a ping
- visible convergence markers that show whether peers are close enough in
  config, secrets, and revision to substitute honestly
- durable routing artifacts or generated config that survive the local failure
  they are meant to absorb
- drill output, status surfaces, or logs that explain why a request was served
  locally, forwarded, or denied

The filename choices are not sacred.
The evidence classes are.

If a proposal cannot name what new inspectable artifacts it would cause the
repo to have, then it is still too abstract to count as a serious middle-layer
answer here.

## What still does not count as inspectable artifacts

Even this page's artifact language needs a harsher filter.

The following still do not count as the kind of inspectable evidence this repo
needs:

- a generated file with no documented consumer
- a status page that repeats what the operator already had to infer
- a peer list that shows reachability but not semantic eligibility
- dynamic config that disappears under the failure it is meant to absorb
- logs that explain only that forwarding happened, not why it was trustworthy

The artifact has to make the system more authoritative, not merely more
verbose.

## The wanted layer, stated as responsibilities instead of brands

The middle layer has to own responsibility, not only configuration syntax.

At minimum, it has to own six families of truth.

### 1. Placement truth

The repo repeatedly converges on `services.yaml` as the simplest mental model:

- which service exists on which node
- which protocol class it belongs to
- which ports or hostnames matter
- whether it has multiple backends
- whether the route is local-first only, peer-forwardable, or harsher

The filename is not sacred.
The property set is.

The truth has to be:

- explicit
- current
- shared
- inspectable
- consumed by routing logic instead of merely documented for humans

If the registry exists but routing still depends on remembered host placement,
the middle layer has failed.

### 2. Eligibility truth

A peer being reachable is not the same thing as a peer being safe.

The middle layer has to distinguish:

- node exists
- node is reachable
- service is present
- service is healthy
- service is converged enough
- service is semantically eligible for this traffic right now

That last line matters because the repo keeps refusing fake equivalence
between:

- a port that answers
- and a service that can honestly preserve the same behavior

Without eligibility truth, wrong-node forwarding is still decorated guessing.

### 3. Route durability truth

The repo's own planning material records that `docker-gen-failover` can remove
Traefik routes when containers stop.

That is devastating for any serious failover story.

It means the middle layer cannot merely generate dynamic config.
It has to generate or preserve rescue routes in a way that survives the exact
failure that makes rescue necessary.

If the fallback path disappears at the moment of local failure, the helper is
not a failover layer.
It is part of the outage.

### 4. Convergence truth

The repo also records manual secret sync and manual Compose sync as still-open
gaps.

That matters because peer substitution is dishonest if:

- the peer lacks the same secrets
- the peer lacks the same env assumptions
- the peer runs a materially different service definition
- the peer has drifted in image version, middleware, or route labels
- the peer looks healthy but is semantically on another world

So even a "small" middle layer needs some convergence story.
It does not need to be a full scheduler by definition.
It does need to stop cross-node substitution from being a silent drift gamble.

### 5. Policy-continuity truth

The user does not only want packets to arrive somewhere.
They want the request to keep meaning the same thing.

For protected HTTP, that means the middle layer needs enough truth to preserve:

- auth expectations
- middleware behavior
- headers and rewrites that define route identity
- security posture relevant to the route

If the forwarding story survives only at the transport layer, it still has not
solved the user’s real complaint.

### 6. Explanation truth

The user does not only want automation.
They want the system to stay inspectable.

That means the middle layer has to make it easier to answer:

- why did this node serve locally?
- why did it forward?
- why was that peer chosen?
- why did the rescue route still exist?
- why did the request keep the same meaning?

If the helper gets smarter while explanations get darker, the repo has gained
machinery and lost trust.

## What disqualifies a fake middle answer

A candidate is **not** the missing middle for this repo if it mainly does one
or more of these:

- improves ingress reachability without solving wrong-node meaning
- hides placement truth behind magic without making it inspectable
- uses larger cluster vocabulary while the operator still has to privately join
  the important facts together
- forwards traffic but leaves middleware, auth, or service meaning ambiguous
- treats stateful services as resilient because a route still answers
- sounds smaller than Kubernetes while quietly importing a scheduler's
  worldview tax anyway
- gives the repo more moving parts while leaving the same hidden burden intact

The repo can still learn from such candidates.
They just have not solved the wound this page is tracking.

## What the user is actually asking these candidates to stop doing

The repo keeps revisiting candidate families because ordinary self-hosting
answers keep doing one insulting thing:

- they offer more nouns
- they offer more diagrams
- they offer more cluster vocabulary
- they offer more traffic paths

while still leaving the operator with the same emergency burden:

- remember which node is real for this service
- remember whether the peer is actually aligned
- remember whether the rescue route will still exist after failure
- remember whether the protected route still means the same thing after handoff

That is why this page keeps sounding harsher than a normal comparison page.
The user is not shopping for technologies.
They are trying to stop being the unpaid control plane.

## Candidate families already visible in the repo

The repo is not searching blindly.
Several candidate families are already visible in planning and research.

## Candidate family 1: lightweight registry plus helper agents

This is the most obvious middle-shaped path already visible in the repo.

It usually implies:

- a shared current-state file or registry such as `services.yaml`
- peer broadcast or sync
- helper logic for failure detection and redeploy
- route generation or route updates from tracked truth

Why it is attractive:

- it preserves Compose as the primary authoring surface
- it externalizes placement truth without immediately promoting a full
  scheduler
- it stays closer to the repo's demand for inspectable ownership

Why it is still incomplete:

- it still has to prove route durability
- it still has to prove safe peer eligibility
- it still has to prove policy preservation
- it still has to prove convergence semantics that are stronger than “git pull
  happened”
- it still risks quietly becoming a scheduler in disguise

The risk here is not only failure.
It is invisible growth into an unacknowledged control plane.

## Candidate family 2: OpenSVC-shaped service ownership

The OpenSVC research pages are interesting because they feel closer to a
truth-owning supervision layer than a pure proxy trick.

Why this family is attractive:

- stronger service ownership and supervision semantics
- better fit for explicit resource and service responsibility than pure
  sidecar-style glue
- possible way to answer wrong-node rescue without inventing every primitive
  from scratch

Why it is still incomplete:

- not yet proven in the live runtime
- still must show that it removes burden instead of renaming it
- must prove that it keeps the operator surface legible enough to count as a
  real middle rather than just a different control plane
- still owes explicit answers for policy continuity and stateful harshness

This family is promising only if it can stay narrow relative to the truths it
claims to own.

## Candidate family 3: Nomad-style promotion

Nomad appears repeatedly because it offers more scheduler truth than plain
Compose without looking as worldview-heavy as Kubernetes.

Why this family is attractive:

- clearer placement and health semantics than ad hoc helper glue
- real scheduling primitives
- lighter reputation than Kubernetes

Why it is still incomplete:

- it still imports a scheduler worldview and therefore owes a burden-removal
  justification
- it does not automatically prove preserved wrong-node semantics for this repo's
  actual stack
- it may solve more placement than the user needs while still demanding trust
  across more opaque machinery

Nomad only counts as middle if it is genuinely the smallest honest answer left,
not merely the smaller famous orchestrator.

## Candidate family 4: k3s / Kubernetes promotion

This family is not irrational.
It is just expensive in the repo's value system.

Why this family is attractive:

- broad cluster truth ownership
- mature ingress, service-discovery, and scheduling machinery
- fewer DIY coordination primitives

Why it is still incomplete:

- it carries the largest worldview cost
- it risks reducing Compose legibility too early
- it has not yet been proven to be the smallest honest answer to this repo's
  actual wound
- it can easily become the default answer before the narrower burden-transfer
  question was answered properly

This family may still win later.
This page simply refuses to let it win by prestige alone.

## Candidate family 5: improved proxy automation alone

This family tries to keep the solution close to the edge stack.

Why it is attractive:

- minimal platform change
- direct effect on ingress behavior
- intuitive when the pain is phrased narrowly as routing

Why it is still incomplete:

- proxy automation alone does not own placement truth
- it does not solve convergence truth
- it often collapses stateful correctness into reachability
- the `docker-gen-failover` bug already proves proxy-only automation can fail
  at the worst possible time

This family is useful only if it becomes part of a larger truth-owning story.
On its own, it is too close to the exact category of partial answer the user
already distrusts.

## The real success test

The missing middle has only been found when this scene stops feeling fragile:

1. a request lands on a healthy public node
2. the node does not host the service locally
3. the node knows where the service lives now
4. the node knows which peer is eligible now
5. the rescue route survives the local failure
6. auth and middleware still preserve the same meaning
7. the operator can explain all of that from shared tracked truth

Anything weaker may still be useful engineering.
It is not yet the middle layer this repo is searching for.

## What would actually count as promotion evidence

No candidate family should be called "the missing middle" until the repo can
show a proof packet that includes all of the following for at least one real
service path:

- the current placement truth the system consumed
- the reason a remote peer was or was not eligible
- evidence that the rescue route survived the local backend loss
- evidence that the forwarded route preserved auth or middleware meaning when
  relevant
- evidence that an operator can reconstruct the decision from shared surfaces
  instead of private memory

This matters because a lot of infra work looks persuasive right up until the
moment you ask:

> what exact evidence would prove the hidden burden moved out of the operator
> and into the system?

If the answer is still vague, the candidate has not earned promotion.

## The hidden risk this page must keep visible

One of the biggest risks in the repo is that the "small missing middle"
quietly grows into a scheduler in disguise.

That does not automatically make it wrong.
It does mean the docs should stay explicit about when the repo starts paying:

- worldview cost
- control-plane opacity
- operator legibility loss
- central coordination assumptions

The middle is only honestly "middle" if the added layer stays narrow relative
to the exact truths it is supposed to own.

## Bottom line

The missing middle is not a product waiting to be discovered.

It is the smallest truth-owning layer that would stop the operator from being
the:

- hidden registry
- hidden failover brain
- hidden drift detector
- hidden peer-eligibility judge
- hidden routing explainer

for wrong-node behavior.

The repo already knows many of the responsibilities that layer must own.
It does not yet prove that one implementation has earned promotion as the
answer.
