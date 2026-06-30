# OpenSVC Cluster Bootstrap

> **Reading boundary**: this file describes why OpenSVC is being explored and
> what a clean bootstrap path would need to guarantee. It is **not** proof that
> the priority root runtime is already OpenSVC-governed, nor proof that
> bootstrapping a few nodes into OpenSVC automatically solves wrong-node
> routing, peer fallback, or stateful HA.
>
> For the evidence-first version of this topic, read:
>
> - [`../knowledgebase/research/opensvc-cluster-bootstrap.md`](../knowledgebase/research/opensvc-cluster-bootstrap.md)
> - [`../knowledgebase/research/osvc-ingress-ha.md`](../knowledgebase/research/osvc-ingress-ha.md)
> - [`INFRASTRUCTURE_MASTER_PLAN.md`](INFRASTRUCTURE_MASTER_PLAN.md)

# Why This Exists

OpenSVC keeps reappearing in this repo because it may be one of the few
serious candidates for strengthening cluster membership and narrow failover
logic **without** forcing the whole stack to stop being Compose-first.

That is the attractive part.

The dangerous part is obvious too:

it would be easy to overread "an OpenSVC cluster exists" into claims the repo
does not yet prove, such as:

- full any-node request correctness
- shared placement truth for the whole stack
- trustworthy peer-forward routing under failure
- stateful zero-SPOF behavior

This page exists to keep that boundary hard.

## What OpenSVC Is Being Asked to Solve Here

The repo is not exploring OpenSVC because "clustering is cool."

It is exploring OpenSVC because the current stack still lacks a sufficiently
trustworthy answer to several infra-grade questions:

- how do nodes become recognized members of one coordination surface
- how can infrastructure-level failover stay more explicit than ad hoc scripts
- what stronger truth can route generation consume than "whatever containers
  were running on this one box a second ago"
- how do we remove more sacred-node behavior without immediately importing a
  full Kubernetes worldview

That is the real ask.

## What OpenSVC Would Need to Make True to Matter

OpenSVC only earns its keep here if it makes the system more trustworthy in
ways the current root runtime does not already prove.

At minimum, a useful OpenSVC-backed path would need to improve confidence in:

### 1. Membership truth

The system needs a better answer to:

- which nodes are currently real participants
- how nodes discover one another
- what cluster identity is authoritative enough to drive routing or failover

### 2. Infra failover truth

The system needs a better answer to:

- what happens when one narrow critical service disappears
- who chooses the replacement or alternate path
- what prevents the fallback route from disappearing with the dead backend

### 3. Bootstrap trust

The system needs a better answer to:

- when a new node is truly in the system
- when it is safe as a peer
- when it is safe as an ingress participant

### 4. Reduced operator memory burden

The system needs a better answer to:

- which facts can stop living in the operator's head
- which transitions can become explicit and inspectable

If OpenSVC does not improve those, it is just more machinery.

## What This Does Not Mean

Exploring or even partially standing up an OpenSVC cluster does **not** by
itself mean:

- the root [`../docker-compose.yml`](../docker-compose.yml) has stopped being
  the live authoring contract
- the repo has settled on OpenSVC as the universal future control plane
- OpenSVC now owns app placement truth for the entire stack
- the HTTP peer-forward model is fully proven
- the stateful HA story is complete

Those broader claims require stronger evidence than node join commands or
generated files.

## The Minimum Honest Bootstrap Story

If OpenSVC is used here, the bootstrap story only stays useful if it answers
the following concretely.

### Host readiness

Questions:

- is the node prepared for the chosen runtime assumptions
- are required packages, networking, storage, and container dependencies ready
- does the host match the repo's actual operating model

### Identity and membership

Questions:

- how does the node become a recognized cluster member
- what stable name does it use
- how do other nodes reach it on the private path

### Coordination handoff

Questions:

- what data from OpenSVC will downstream route-generation or failover logic
  consume
- how fresh is that data
- what happens if the data source itself is stale or partitioned

### Safety as an ingress participant

Questions:

- when may this node safely receive global traffic
- when may it safely forward to peers
- what guarantees are required before calling that safe

If those answers remain vague, bootstrap is still mostly ritual.

## A Minimal, Repo-Grounded Experimental Path

This is the most conservative way to interpret OpenSVC work in this repo.

### Step 1: treat OpenSVC as an infra-grade membership and supervision layer

At this stage, it is helping answer:

- who is in the cluster
- who is healthy enough to matter
- what infra-level services should stay alive

It is **not** yet being granted authority over the entire Compose runtime.

### Step 2: let route-generation experiments consume OpenSVC membership

At this stage, OpenSVC can become one input to:

- Traefik failover file generation
- L4 config generation
- node eligibility for fallback

This is meaningful, but still weaker than saying OpenSVC is the global runtime
truth for everything.

### Step 3: prove whether this actually lowers wrong-node stupidity

This is the key test.

If OpenSVC-backed membership cannot materially improve:

- peer selection
- route persistence
- fallback correctness
- reduced operator reconstruction tax

then it has not earned expansion.

## What Would Count as Stronger Proof Later

This repo should only narrate stronger OpenSVC claims when evidence exists for
them.

Examples of stronger evidence would include:

- a proven, repo-grounded cluster bootstrap flow using current paths and
  current scripts
- explicit proof that route generation now survives local backend death more
  reliably because it is consuming stronger membership truth
- explicit proof that node participation decisions are less manual and less
  private-memory driven
- clear boundaries on which services or domains are actually OpenSVC-governed

Without that, "OpenSVC cluster bootstrap" should remain a serious direction,
not a completed architecture claim.

## Important Repository Boundary

Older bootstrap language often implied a different repo shape or older paths.
That is not good enough now.

Read this topic through the current repository boundary:

- the priority implementation still centers on [`../docker-compose.yml`](../docker-compose.yml)
- the current docs authority lives in [`../knowledgebase/`](../knowledgebase/)
- OpenSVC remains one explored strengthening layer, not the settled whole-stack
  runtime owner

## Bottom Line

OpenSVC is interesting here because it may help solve a very specific missing
middle problem:

- stronger cluster membership
- stronger infra failover primitives
- stronger bootstrap trust

without demanding immediate full scheduler conversion.

That is why the repo keeps revisiting it.

But until the live runtime proves a stronger dependency on OpenSVC, the honest
reading stays:

- serious candidate
- meaningful experiment
- not yet the whole answer
