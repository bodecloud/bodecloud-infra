# The Missing Middle Layer

This page answers a very specific question:

> if raw Compose is too dependent on private operator memory, and immediate
> Kubernetes or Swarm promotion is too expensive or too ideological, what is
> the smaller truth-owning layer the repo is actually searching for?

That is the "missing middle" in `bolabaden-infra`.

It is not a brand name. It is a capability shape.

## What this page is and is not allowed to prove

This page is authoritative about:

- what category of helper layer the repo actually wants
- which burdens that layer must remove to count as real progress
- why many adjacent solutions still feel fake to this user

This page is not authoritative about:

- claiming that the repo has already found the winner
- promoting one research path into the live runtime
- pretending that "lighter than Kubernetes" is a sufficient definition

This page defines the wanted shape, not the final chosen implementation.

## Strongest honest current answer

The repo is searching for a layer that owns enough shared truth to make
wrong-node recovery, peer selection, and bad-day explanation explicit, while
still keeping `docker-compose.yml` legible and manually assignable.

That layer is still incomplete.

The repo already has:

- real Compose runtime breadth
- a strong edge layer
- planning around failover helpers, secret sync, compose sync, and Headscale
- research into OpenSVC, Nomad, k3s, and related alternatives

What it does not yet have is a promoted shared-truth layer that cleanly
answers:

- where the service lives now
- whether the peer is eligible now
- whether the route survives now
- whether the handoff preserves meaning now

## The wanted layer, in plain terms

The missing middle is the smallest added surface that would make these things
system-owned instead of remembered:

- service placement truth
- peer eligibility truth
- fallback-route truth
- convergence truth for secrets, env, and deployment shape
- operator-readable explanation of what happened

That layer does **not** need to be a full scheduler by default.
It **does** need to stop private topology memory from being the real control
plane.

## What disqualifies a fake middle answer

A candidate is not the missing middle for this repo if it mainly does one or
more of these:

- improves ingress reachability without solving wrong-node semantics
- hides placement truth behind magic without making it inspectable
- uses bigger-cluster vocabulary while the operator still has to privately join
  the facts together
- forwards traffic but drops auth, middleware, or visible service meaning
- treats stateful services as HA because some route still answers
- sounds smaller than Kubernetes while still demanding worldview-scale trust

The repo can still learn from such candidates.
They just have not solved the wound this page is tracking.

## What the current repo evidence says the missing layer probably includes

Based on:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- current root Compose and include surfaces

the missing layer likely needs these concrete capabilities.

### 1. Placement registry or equivalent current-state truth

The repo repeatedly converges on `services.yaml` as the simplest mental model:

- what service lives on which node
- which ports or classes it exposes
- whether there are multiple backends

The important part is not the filename.
The important part is that the truth be:

- explicit
- current
- readable
- shared

### 2. Peer health and eligibility model

Reachability is not enough.
The layer has to distinguish:

- peer exists
- peer is reachable
- peer is healthy for this service
- peer is semantically safe to receive this traffic

Without that distinction, wrong-node forwarding becomes a more decorated guess.

### 3. Route materialization that survives failure

Planning material explicitly calls out `docker-gen-failover` as broken because
it can delete routes when containers stop.

That means the missing layer cannot merely generate dynamic config.
It has to generate or preserve rescue routes in a way that survives the failure
that made rescue necessary.

### 4. Convergence for secrets, env, and deployment shape

The master plan also names manual secret sync and manual compose sync as real
gaps.

That matters because peer substitution is dishonest if:

- the peer lacks the same secrets
- the peer lacks the same env assumptions
- the peer runs a materially different service definition

So the missing middle needs some convergence story, even if it is lighter than
a full scheduler.

### 5. Operator-readable explanation

This repo does not only want automation.
It wants a system that can still explain itself.

That means the added layer should make it easier to answer:

- why this node forwarded
- why this peer was selected
- why the route still existed
- why the request kept the same meaning

If the helper gets smarter while the explanation gets darker, the repo has
gained machinery and lost trust.

## Candidate families already visible in the repo

These are the main candidate families already present in planning or research.

| Candidate family | What it might solve | Why it is still incomplete |
| --- | --- | --- |
| lightweight registry plus helper agents | placement truth, sync, redeploy, failover logic | still must prove durability, eligibility, and semantic preservation |
| OpenSVC-shaped ingress or service supervision | stronger service ownership and failover semantics | still needs proof that it removes burden instead of renaming it |
| Nomad-style promotion | scheduling and health-aware placement | must earn its tax and not just widen the worldview |
| k3s / Kubernetes promotion | broad cluster truth and scheduling machinery | may solve many layers, but the repo treats the cost and loss of legibility as real |
| improved proxy automation alone | route generation and local edge behavior | not enough unless it owns placement, durability, and explanation too |

## The real success test

The missing middle has only been found when this scene stops feeling fragile:

1. request lands on a healthy public node
2. that node does not host the service locally
3. the node knows where the service lives now
4. the node knows which peer is eligible now
5. the route survives the local failure
6. auth and middleware still mean the same thing
7. the operator can explain all of that from shared tracked truth

Anything weaker may still be useful engineering.
It is not yet the missing middle this repo is looking for.
