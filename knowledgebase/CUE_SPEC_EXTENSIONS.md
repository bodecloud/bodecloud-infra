# CUE Specification Extensions

This page explains what the `x-cue` idea is really doing inside this repo.

It is not here to pretend that `x-cue` already drives the live root runtime.

It is here because `x-cue` exposes a very specific pressure:

> plain Compose is still the favored human authoring surface, but the repo keeps needing to describe coordination semantics that plain Compose does not express cleanly.

That makes `x-cue` important as architectural evidence even if the literal key name never becomes production reality.
This page is strongest when read as a pressure signal, not a schema sales
pitch.
The point is not "here is a neat extension surface."
The point is "the repo keeps needing semantics that plain Compose cannot
honestly carry by itself."

There is a second warning this page has to keep visible:

- once richer semantics get named
- naming starts sounding like shared understanding
- shared understanding starts sounding like runtime ownership
- the docs forget that a named semantic can still be socially reconstructed
  rather than system-owned

That drift is one of the easiest ways for a schema page to overstate maturity.

## What `x-cue` is actually trying to rescue

The repo's problem is not that YAML is too short.

The problem is that once the stack becomes genuinely multi-node, the operator needs to express things like:

- whether a service should prefer local serving but allow peer fallback
- whether a service is active-active, active-passive, singleton, or "do not casually move this"
- whether a service should be publicly routable, cluster-visible, or only locally reachable
- whether a workload has stronger stateful caution than a restart policy can capture
- whether a service should be promoted to a stronger control layer once risk crosses a threshold

Plain Compose can define containers, labels, networks, ports, secrets, and volumes very well.

It is much worse at expressing those higher-order semantics in a way that remains explicit and durable.

`x-cue` is one answer to that gap.

But it is only an honest answer if the page keeps distinguishing:

- semantics the repo wants to express
- semantics a future layer might consume
- semantics the current worktree actually proves are interpreted today

## Why this topic matters more than it looks

If `x-cue` is read casually, it sounds like yet another extension idea.

But in this repo it is doing something more revealing:

- it shows the repo does not actually want to abandon Compose first
- it shows the repo also knows Compose alone is semantically too weak for the dream
- it shows the repo is searching for a middle layer that adds meaning without forcing total platform migration

That is why `x-cue` matters even before implementation.

It tells you what the repo feels is missing.
It also tells you that the missing layer is not merely runtime code.
Part of the wound is descriptive: the operator lacks a durable place to say
what kind of coordination truth a service actually needs.

That matters because descriptive clarity can easily be mistaken for operational
closure.

The user is not asking for prettier ways to describe pain.
They are asking for the system to stop depending on private reconstruction of
that pain.

## What this means in the actual stack

This page gets more useful once the semantic gap is tied back to the live
workload classes in the repo instead of left at the level of abstract schema
design.

The current root runtime already contains service classes with very different
operational meaning:

- public L7 surfaces such as `homepage`, `dozzle`, `portainer`, docs, and other
  routed dashboards
- state-bearing core services such as `mongodb`, `redis`, `rabbitmq`, and
  Postgres-shaped surfaces
- mesh and identity infrastructure such as `headscale`
- worker and queue-backed surfaces such as the Firecrawl-related services
- locality- or hardware-sensitive media and LLM surfaces

Plain Compose can define all of them.
It cannot describe all of them well enough for the dream this repo is chasing.

That is the practical reason `x-cue` matters here.

## What `x-cue` actually proves

### 1. The repo wants to raise Compose's semantic ceiling before replacing Compose

This is the clearest signal.

The repo's dissatisfaction is not mainly with Compose syntax.

It is with the fact that plain Compose has no first-class vocabulary for:

- placement truth
- failover mode
- visibility intent
- stateful caution
- hardware capability hints
- control-plane promotion thresholds

`x-cue` proves that the repo wants to enrich Compose upward instead of discarding it immediately.

It does **not** prove that enrichment alone would make the stack meaningfully
less dependent on human interpretation.

### 2. The repo prefers optional semantics over mandatory platform migration

This is tightly aligned with the user's underlying frustration.

The user does not want more arbitrary abstraction just for its own sake.

The user wants more expressive power **only where it removes real pain**.

That is why `x-cue` is framed as:

- keep Compose valid
- keep authoring recognizable
- add richer semantics only where the stack has outgrown plain Docker-era assumptions

### 3. The repo is already thinking in coordination semantics, not just service definitions

The rest of the knowledgebase already separates:

- convergence truth
- placement truth
- ingress truth
- failover truth
- stateful truth

`x-cue` is the schema-level version of that same mental model.

It is an attempt to say:

- these meanings exist
- the stack keeps depending on them
- they should not remain hidden in operator folklore

### 4. The repo has not resolved whether those semantics belong in metadata, code, or both

This is one of the most important truths the older docs softened too much.

`x-cue` suggests a metadata-first future.

But [`infra/docs/ARCHITECTURE.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/infra/docs/ARCHITECTURE.md) pushes a materially different instinct:

- imperative over declarative
- infrastructure defined in Go code
- gossip and consensus as active coordination

That means `x-cue` is not just a schema idea.

It is also evidence of an unresolved repo-level question:

> should the missing semantics live primarily in declarative extensions, in active agent code, or in a hybrid contract where metadata describes intent and code enforces it?

That question is still open.
And it needs to stay open in the docs until the worktree genuinely settles it.
If `x-cue` gets narrated like a mostly-decided future, this page starts lying
in exactly the polished way the user is already tired of.

The page also has to keep a sharper version of the same warning:

metadata can make intention easier to inspect while still leaving ownership
unresolved.

If a future contributor can read the schema and still has to privately answer
"yes, but what really happens on the wrong node right now?", then the missing
truth layer is still missing.

## What `x-cue` does not prove

This page needs to keep hard edges.

It does **not** prove:

- that the current root runtime reads `x-cue`
- that any tracked binary in the root implementation turns `x-cue` into live behavior
- that there is a stable enforced schema with operator-safe guarantees
- that HA, failover, mesh visibility, or stateful policies can be relied on today because a field family was imagined for them

Those would be operational claims. The worktree does not support them yet.

Add one more non-proof explicitly:

- that naming richer semantics has already reduced the amount of private
  operator interpretation required during failure

## The real problem `x-cue` is trying to solve

Plain Compose is very good at saying:

- run this container
- attach these networks
- expose this port
- mount this data
- set these labels

Plain Compose is much weaker at saying:

- keep this route alive even if the local backend disappears
- prefer this node, but forward safely to peers
- do not let this singleton split brain
- this workload is state-sensitive and should be promoted or handled differently
- this service has cluster visibility but should not be public
- this workload needs GPU or node capability affinity

The repo keeps running into exactly those semantic ceilings.

`x-cue` is a way of saying:

- keep Compose for container definition
- put the missing semantics beside it
- let a stronger layer decide how to enforce them later
That sentence is powerful and dangerous at the same time.
It is powerful because it preserves Compose readability.
It is dangerous because "enforce them later" can hide a huge amount of
unsettled control-plane ownership if the docs stop being explicit.

That phrase, "enforce them later," should stay uncomfortable throughout this
page.

It is exactly where elegant schema design can start laundering unresolved
ownership into what sounds like inevitable future capability.

## The most useful `x-cue` intent families

The exact example syntax is less important than the operational meaning each family tries to capture.

### HA mode intent

This family exists because restart policies and naive replica counts do not tell the whole truth.

The repo needs language for things like:

- active-active versus active-passive
- singleton with lease or election semantics
- acceptable failover behavior
- survivability expectations

This is a real semantic need because the repo is trying to stop saying "HA" when it only means "container restarts."

### Placement intent

This family exists because the repo needs a vocabulary for:

- where a service should run
- where it may run if needed
- where it must not run
- whether peer pickup is allowed

This is where the missing live tracked `services.yaml` concept and `x-cue` concept overlap strongly.

Both are trying to make placement less implicit.

### Visibility and mesh intent

This family exists because exposure in this stack is not binary.

The repo keeps caring about distinctions like:

- local only
- cluster reachable
- public through Traefik
- private but peer forwardable
- L4 only versus L7 routable

Plain network membership alone does not fully communicate those meanings.

### Stateful caution intent

This family exists because the repo is trying to stay honest about data-bearing workloads.

The operator needs a way to mark:

- this service should not be moved casually
- this service requires stronger continuity semantics
- this service cannot be treated like stateless web traffic

This is one of the most important semantic gaps in the entire repo.

### Capability and hardware intent

This family exists because some workloads depend on:

- GPU
- transcoding devices
- storage locality
- node-specific acceleration or capacity

Those concerns already exist in the stack whether or not they are encoded formally.

### Promotion-threshold intent

This is one of the most important hidden families.

The repo keeps multiple futures alive:

- Compose plus helpers
- OpenSVC
- Nomad
- k3s
- Kubernetes
- Constellation-style cluster logic

`x-cue` can be read as a way of marking which workloads or domains have outgrown plain Compose and deserve promotion into stronger control semantics.

That is more useful than treating promotion as folklore.

## The semantic gaps by workload class

This is the concrete version of the whole page.

| Workload class | Current repo-shaped examples | Meaning the operator actually needs | Why plain Compose is too weak on its own |
|---|---|---|---|
| Public stateless HTTP surfaces | `homepage`, docs, dashboards, some app frontends | Safe local-first serving with peer fallback only if route continuity and middleware continuity stay intact. | Compose can define labels and ports, but not "prefer local, preserve route, and keep policy identical on peer handoff." |
| Stateful single-writer services | `mongodb`, writable `redis`, some Postgres-shaped services | Must not casually move, duplicate, or split-brain; write authority and promotion semantics matter more than reachability. | Compose can define a container, not "singleton with stateful caution and explicit promotion requirements." |
| Cluster-visible but not public services | internal queues, helper backends, mesh-adjacent services | Must be reachable by peers or workloads without being treated as public entrypoints. | Network attachment alone does not fully capture intended visibility or trust boundaries. |
| Worker and queue consumers | Firecrawl-adjacent workers, async processing surfaces | May be restartable or relocatable differently from the frontend or queue they depend on. | Compose can define restart policies, but not richer "pickup allowed", "queue-safe to relocate", or "locality preferred" semantics. |
| Locality- and hardware-sensitive services | media stack, transcoding surfaces, model workloads | Need explicit hints about storage affinity, hardware capability, and whether cross-node relocation is realistic or fantasy. | Compose can mount devices and volumes, but not encode the operator-level meaning of those dependencies in a durable way. |
| Infra promotion candidates | ingress helpers, failover agents, stronger coordination surfaces | Need a way to say "this domain has outgrown plain Compose and deserves stronger control semantics." | Compose has no native vocabulary for maturity, promotion threshold, or stronger coordination ownership. |

This is the actual value of the `x-cue` idea:

- not prettier config
- not extra abstraction for its own sake
- but a way to stop very different workload classes from being narrated as if
  they all live at the same coordination maturity level

That is only half the job, though.

The other half is preventing these classes from being narrated as if naming the
difference already means the runtime owns the difference.

## The design trade inside `x-cue`

`x-cue` is only attractive if it preserves the correct trade.

### What it is trying to preserve

- one primary authoring surface
- Docker and Compose continuity
- operator readability
- lower friction than a full platform migration

### What it is trying to add

- richer coordination semantics
- more explicit HA and placement intent
- future backend portability
- less hidden meaning buried in labels, scripts, and memory

### What it risks adding instead

- shadow metadata nobody actually enforces
- a second config universe that feels cleaner on paper than in operations
- semantic claims that overpromise what the live runtime can really do

That risk is exactly why the documentation must stay explicit about evidence classes.

It is also why this page should keep sounding slightly dissatisfied.

If it starts reading like "the semantic model is basically handled," then it is
quietly repeating the same mistake the repo is trying to escape everywhere
else.

## Why metadata alone would still not solve the whole problem

This page also has to keep another hard boundary:

even very good metadata would not solve the full repo problem by itself.

Why:

- metadata can describe placement intent, but not make peers converge on it by
  itself
- metadata can describe failover mode, but not keep a route alive when a
  generator deletes the route on container death
- metadata can describe stateful caution, but not create replication,
  promotion, quorum, or reconnect behavior
- metadata can describe visibility, but not by itself prove policy continuity
  on wrong-node handoff

So the best reading is hybrid:

- metadata makes missing meanings explicit
- active coordination code consumes or enforces those meanings
- the docs stay explicit about which meanings are only named, which are
  partially enforced, and which are actually live

One more distinction belongs there:

- which meanings are still being privately reconstructed by operators even after
  they have been named cleanly in metadata

Without that distinction, a schema layer can easily become one more elegant way
to overstate maturity.

## Why `x-cue` still matters if the final implementation becomes agent-driven

This is the part worth making explicit.

Even if the repo eventually decides that active coordination belongs mostly in:

- Constellation
- a failover-agent
- generated Traefik and L4 config
- or another Go-based convergence engine

the missing semantics still need a home.

The literal `x-cue` key could disappear and the problem would remain.

The problem is:

- where are these meanings written down
- how are they made explicit
- how are they kept from dissolving into tribal memory

So `x-cue` matters even if it never survives as a concrete syntax.

It is a catalog of the meanings the repo cannot avoid.

## How the repo should use this idea today

### 1. Keep the authoring-contract boundary explicit

Compose remains the live root authoring contract until a stronger layer is actually proven in the priority implementation.

### 2. Use the idea to name missing semantics, not to pretend they are solved

If the repo keeps tripping over an intent that Compose cannot describe well, the docs should name that intent directly instead of hiding it in labels, custom scripts, or assumptions.

### 3. Distinguish metadata aspirations from enforcement reality

Every semantic family should be documented in one of these classes:

- live
- partially implemented
- planned
- conceptual

That is the only honest way to keep the docs useful.

### 4. Use `x-cue` to force promotion criteria into the open

If a workload class needs more than Compose can honestly provide, the repo should say why and what stronger layer is justified.

That is far better than pretending every service belongs at the same maturity level forever.

## Current evidence verdict

The most honest reading is:

- `x-cue` is strong evidence that the repo wants Compose-plus-semantics rather than immediate Compose abandonment
- it is also strong evidence that the repo knows plain Compose lacks important coordination vocabulary
- it is most useful when read against real workload classes in the current
  stack rather than as generic schema futurism
- it does not currently prove live runtime support in the tracked root implementation
- it also does not settle whether those semantics will ultimately live in metadata, active agent logic, or a hybrid of both

So `x-cue` should be treated as a schema sketch of the repo's missing meanings, not as proof that those meanings are already enforced by the live system.
