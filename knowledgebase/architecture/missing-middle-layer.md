# The Missing Middle Layer

This page answers the repository's most important architectural question:

> if plain Compose becomes too dependent on private operator memory, and
> immediate promotion to Kubernetes, Swarm, Nomad, or another larger control
> plane feels too expensive or too premature, what is the smaller truth-owning
> layer the repo is actually searching for?

That is the "missing middle."

It is not a brand name.
It is a burden-removal shape.

## What this page is and is not allowed to prove

This page is authoritative about:

- the capability shape the repo is repeatedly converging on
- which burdens the middle layer must actually remove
- which candidate families are visible in the repo
- why many plausible solutions still fail this repo's benchmark

This page is not authoritative about:

- claiming the winner has already been found
- claiming the live runtime already owns this layer
- equating "lighter than Kubernetes" with "good enough"

This page describes the wanted layer, not the final chosen implementation.

## Strongest honest current answer

The repo is looking for the smallest added layer that makes these truths
system-owned instead of remembered:

- where a service lives now
- whether a peer is eligible now
- whether a rescue route still exists now
- whether peers are converged enough to substitute safely
- why the system made the routing or failover decision it made

The repo already has:

- a large real Compose runtime
- a serious Traefik/CrowdSec/auth edge surface
- Headscale as a real private-mesh assumption
- planning pressure toward `services.yaml`, peer sync, and failover helpers
- research into OpenSVC, Nomad, k3s, and related paths

What it still does **not** have is a promoted layer that cleanly owns those
truths across the whole priority implementation.

## Why this layer is needed at all

The repo does not need a missing middle because clustering is fashionable.
It needs one because several facts are still too easy to answer only from
private human reconstruction.

Those facts include:

- what runs where right now
- which node is the right peer for this specific service
- whether that peer is merely reachable or actually safe to receive traffic
- whether the route needed for fallback survives backend disappearance
- whether auth and middleware keep the same meaning during handoff
- whether a stateful service is actually movable or only reachable

As long as those answers remain mostly social, the operator is still the hidden
control plane.

That is the real thing this page is trying to kill.

## Why DNS and proxies are not enough by themselves

The repo already has Cloudflare participation and a strong live edge stack.
That still does not satisfy the missing middle because:

- Cloudflare can choose a healthy first-hop node, but not the correct backend
  service placement inside the user's topology
- Traefik can express rich routing and middleware locally, but not magically
  infer trustworthy cross-node placement truth
- health checks can detect some failures, but not automatically answer whether
  a remote peer is semantically equivalent for this traffic

This is why the repo keeps separating:

- first-hop plurality
- request preservation
- stateful authority

The middle layer is needed because those are different truths.

## The repo has already named the gap directly

The clearest planning pressure is in
[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md),
which explicitly records these as still missing:

- universal wrong-node success
- a live tracked root `services.yaml` current-state registry
- trustworthy route persistence under local backend failure
- automated service failover between nodes

Those are not side quests.
They are exactly the truths the missing middle would need to own.

The master plan also records:

- `docker-gen-failover` is broken because it can delete routes when containers
  stop
- Cloudflare DDNS presence is not the same thing as full failover
- Headscale is single-node today
- secret sync and compose sync are still manual

That list is basically a map of why the middle layer has not been earned yet.

## The wanted layer, stated as responsibilities instead of brands

The middle layer has to own responsibility, not just naming.

At minimum, it has to own five families of truth.

### 1. Placement truth

The repo repeatedly converges on `services.yaml` as the simplest mental model:

- which service exists on which node
- which protocol class it belongs to
- which ports or hostnames matter
- whether it has multiple backends

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
- service is semantically eligible for this kind of traffic right now

That last line matters because the repo keeps refusing fake equivalence between:

- a port that answers
- and a service that can honestly preserve the same behavior

Without eligibility truth, wrong-node forwarding is only decorated guessing.

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
- the peer has drifted in image version, sidecar assumptions, or route labels

So even a "small" middle layer needs some convergence story.
It does not need to be a full scheduler by definition.
It does need to stop cross-node substitution from being a silent drift gamble.

### 5. Explanation truth

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
- uses bigger cluster vocabulary while the operator still has to privately join
  the facts together
- forwards traffic but leaves middleware, auth, or service meaning ambiguous
- treats stateful services as resilient because a route still answers
- sounds smaller than Kubernetes while quietly importing a scheduler's worldview
  tax anyway

The repo can still learn from such candidates.
They just have not solved the wound this page is tracking.

## Candidate families already visible in the repo

The repo is not searching blindly.
Several candidate families are already visible in planning and research.

### Candidate family 1: lightweight registry plus helper agents

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
- it still risks quietly becoming a scheduler in disguise

### Candidate family 2: OpenSVC-shaped service ownership

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

### Candidate family 3: Nomad-style promotion

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

### Candidate family 4: k3s / Kubernetes promotion

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

### Candidate family 5: improved proxy automation alone

This family tries to keep the solution close to the edge stack.

Why it is attractive:

- minimal platform change
- direct effect on ingress behavior
- seems intuitive when the pain is phrased as routing

Why it is still incomplete:

- proxy automation alone does not own placement truth
- it does not solve convergence truth
- it often collapses stateful correctness into reachability
- the `docker-gen-failover` bug is already evidence that proxy-only automation
  can fail precisely at the wrong time

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

## The hidden risk this page must keep visible

One of the biggest risks in this repo is that the "small missing middle"
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
the hidden registry, hidden failover brain, hidden drift detector, and hidden
explainer for wrong-node behavior.

The repo already knows many of the capabilities that layer must own.
It does not yet prove that one implementation has earned promotion as the
answer.
