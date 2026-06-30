# CUE Bootstrap Protocol

This page explains what the CUE bootstrap concept is really trying to fix.

It is not a description of a finished bootstrap engine in the current root implementation.

That distinction matters because the repo already has real bootstrap artifacts:

- [`cloud-init-bootstrap.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/cloud-init-bootstrap.sh)
- [`scripts/bootstrap_opensvc.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/bootstrap_opensvc.sh)
- environment and secret assumptions spread across the Compose surface
- planned node registration and failover-agent steps in [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

The point of this page is not to deny that bootstrap work exists.

The point is to say that the existing bootstrap reality and the desired bootstrap experience are still far apart.
That gap has to stay emotionally visible.
If bootstrap gets narrated as "mostly a scripting problem," this page has
already shrunk the real issue into something too small.

There is a second shrinkage this page has to keep resisting:

- bootstrap steps get listed clearly
- the list starts sounding like a converged join contract
- the contract starts sounding like the runtime itself owns node truth
- the hidden operator reconstruction burden disappears from view

That is how a bootstrap page can become more reassuring while becoming less
honest.

## The user's real bootstrap frustration

The user is not complaining because a few shell commands are annoying.

The deeper frustration is that "bring up a node" still often means:

- remember which base host assumptions matter
- remember which secrets must already exist
- remember which env files and inline config dependencies must align
- remember which services should run locally versus just be reachable remotely
- remember how ingress, DDNS, mesh identity, and peer forwarding are supposed to converge
- remember which parts of the docs are live truth and which are future plans

That is not just setup friction.

It is evidence that bootstrap, convergence, and failover are not yet one coherent experience.
It is also evidence that node join is still too dependent on human
reconstruction.
That dependence is one of the main hidden SPOFs the repo is trying to remove.

## What the bootstrap concept is actually trying to make true

The bootstrap idea only becomes useful when described in operator terms instead of abstract control-plane terms.

The desired future node join looks something like this:

1. a new node is prepared consistently
2. it becomes a known peer with stable identity
3. it receives the right secret and env material
4. it learns what it is allowed or expected to run
5. it can participate in any-node ingress safely
6. it does not require a pile of tacit operator memory to become trustworthy

That is a much higher bar than "Docker installed successfully."
It is really closer to:

> a new node becomes a believable participant in the same truth model the rest
> of the system claims to use

That sentence has to be read adversarially.

If the node only becomes believable because the operator still privately knows:

- which secrets really matter here
- which services are only supposed to be reachable remotely
- which ingress assumptions are still hand-maintained
- which peers are trusted despite the docs not proving why

then bootstrap is still narratively stronger than operationally owned.

## What the CUE bootstrap concept actually proves

### 1. The repo wants bootstrap to become state-driven instead of memory-driven

This is the clearest signal.

The concept assumes that a newly joined node should derive enough truth from the system rather than from the operator's memory.

That means bootstrap is implicitly trying to unify:

- host prep
- identity
- secret hydration
- service and placement truth
- ingress convergence
- post-join trustworthiness

This aligns strongly with the master plan's emphasis on:

- sync agents
- failover-agent ideas
- lightweight registry concepts such as `services.yaml`
- reduced operator burden

### 2. The repo does not want bootstrap to be a one-off ritual disconnected from steady-state behavior

This is a real architectural preference.

The desired model is not:

- one set of scripts for day zero
- a totally different system for day two

The desired model is:

- the same truth model used for runtime coordination should help bring nodes into that runtime cleanly
- joining a node should not feel like an exception path that depends on private tribal knowledge

That matters because the user's frustration is fundamentally about having too few trustworthy layers, not too few scripts.

### 3. The repo sees environment assembly as one of the real coordination problems

This is bigger than convenience.

The current stack already requires:

- many env vars for Compose interpolation
- placeholder secret files
- config continuity across fragments
- Docker and proxy assumptions that have to line up correctly

That means bootstrap is not merely about machine provisioning.

It is about preventing environment drift from becoming another SPOF.

## What the CUE bootstrap concept does not prove

This page has to stay hard-edged.

It does **not** prove:

- that `cue bootstrap` exists as a real root command
- that the root runtime currently has a live node-join engine hydrating state automatically
- that a tracked root `services.yaml` is already present and consumed by bootstrap logic
- that secrets, routing, and placement are already unified under one authoritative join process
- that Tailscale or Headscale identity, Traefik routing, DNS updates, and service placement already converge through one proven engine

Those would all require stronger live evidence than the current worktree provides.

## The current live bootstrap story

The current story is more fragmented and more honest:

- there are real bootstrap scripts
- there are real Compose definitions
- there are real plans for sync, failover, service registry, and node registration
- there are exploratory stronger-cluster paths such as OpenSVC and Constellation

What the current story does **not** yet provide is a single trustworthy answer to:

- when is a node truly "in the system"
- when is it safe as an ingress participant
- when is its placement truth known
- when are its secrets and runtime assumptions definitely converged

That gap is exactly why bootstrap remains such an important design topic.

It is also why this page cannot let "we have bootstrap assets" drift into "node
join is basically solved."

The repo does have assets.
It still does not have one plainly evidenced answer to whether a newly joined
node knows enough truth on its own to stop being socially babysat.

## The bootstrap layers the repo is trying to unify

The concept becomes much clearer when separated into concrete layers.

### Host readiness

Questions this layer must answer:

- is the machine prepared for the chosen runtime assumptions
- are storage, kernel, networking, and container prerequisites satisfied
- can the host safely run the intended workload classes

The repo has real assets here, but not a single universally trusted convergence contract.

### Node identity and mesh membership

Questions this layer must answer:

- how does the node become a recognized peer
- how do peers address it privately
- what stable name does it have in the cluster
- what trust assumptions depend on that identity

This layer matters because any-node ingress is worthless if nodes cannot reliably reach one another on the private path.

### Secret and env hydration

Questions this layer must answer:

- what material belongs on this node
- how is that material sourced
- how is drift corrected
- how can a new node avoid copy-paste divergence

This is one of the sharpest unresolved pains in the repo.

It is also one of the places where false maturity is easiest to narrate.

Env files, secret placeholders, and helper scripts can make a node look
provisioned while still leaving the operator to reconstruct:

- what was essential
- what was only copied for convenience
- what is actually authoritative when drift appears

### Placement truth hydration

Questions this layer must answer:

- what services should run here
- what services should never run here
- what services are only peer-reachable from here
- where is the authoritative current-state record stored

This is where the missing live tracked `services.yaml` becomes crucial.

Without placement truth, bootstrap can only ever be partially trustworthy.

### Ingress and routing hydration

Questions this layer must answer:

- what traffic can this node receive publicly
- what can it serve locally
- what can it forward to peers
- how do middleware, auth, and health semantics stay consistent during that process

This is not "post-bootstrap polish." It is part of whether node join produced a useful participant at all.

## Why the bootstrap concept matters strategically

Even if literal CUE never becomes the mechanism, the concept still reveals the repo's real standard for success.

That standard is not:

- the host is reachable
- Docker starts
- Compose runs once

The real bar is closer to:

- a node joins with much less operator memory burden
- the node becomes trustworthy as a peer and ingress participant
- bootstrap establishes the same truth the steady-state system will continue to use
- the system stops depending on private operator reconstruction for obvious cluster facts

That is a serious and worthwhile bar.

## The hidden risk inside the bootstrap dream

The attractive part of the concept is also the dangerous part.

A system that owns:

- host prep
- identity
- secret hydration
- placement truth
- ingress generation
- failover readiness
- runtime convergence

is already most of a control plane.
That means this page should be read as a warning as much as an aspiration.
Bootstrap unification can genuinely reduce pain.
It can also quietly become another place where the repo builds an implicit
controller while still pretending it has not chosen one.

And there is an even more specific danger:

bootstrap can become the place where the repo starts acting as though truth is
owned because truth is generated.

Those are not the same thing.

A generator can still leave the operator privately responsible for knowing:

- when the generated state is stale
- what to trust when generated outputs disagree
- whether a joined node is truly eligible for ingress or only configured enough
  to look plausible

That means the repo cannot evaluate bootstrap unification by aesthetics alone.

It has to ask:

- does this really remove operator pain
- does it reduce hidden knowledge
- does it give safer convergence
- or does it just relocate complexity into custom code and more implied machinery

That is the hard standard the repo should keep applying.

## The relation between bootstrap and the unresolved future-control-plane split

This is the part generic docs usually skip.

The repo is not only deciding how nodes should join.

It is also deciding **who owns the truth after they join**:

- richer declarative metadata
- a registry plus generated config
- a failover or convergence agent
- Constellation-style cluster logic
- or a stronger external platform

Bootstrap cannot be fully settled until that ownership question is clearer.

That does not make the concept useless.

It means bootstrap is one of the places where the repo's unresolved orchestration split becomes impossible to ignore.

## What should be carried forward from this document

### 1. Judge bootstrap by operator memory removed

A future bootstrap path is only better if it removes hidden steps the operator currently has to remember.

It is not better merely because the remembered steps are now hidden behind a
cleaner join command.

### 2. Judge bootstrap by runtime trust established

A node is not done because Docker is installed.

It is only meaningfully bootstrapped when:

- identity is known
- secret/env material is coherent
- placement truth is clear enough
- ingress behavior is predictable

And even those need harsher wording in practice:

- "clear enough" means not privately reconstructed from scattered docs
- "predictable" means not dependent on remembering which node is really special

### 3. Keep bootstrap claims separated by evidence class

Docs should continue separating:

- existing scripts and current host-prep truth
- planned bootstrap improvements
- speculative unified bootstrap futures

Without that separation, the docs become exactly the kind of false reassurance the user is already tired of.

## Current evidence verdict

The most honest reading is:

- the CUE bootstrap concept is strong evidence that the repo wants less manual environment assembly and less hidden operator memory
- it is also strong evidence that bootstrap is considered part of the larger placement, failover, and control-plane problem
- it does not currently prove a finished root node-join engine or a fully unified bootstrap/runtime truth model

So this document should be treated as a demanding design target for trustworthy node join, not as proof that the target has already been met.
