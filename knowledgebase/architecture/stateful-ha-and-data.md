# Stateful HA and Data Ownership

For the evidence boundary behind this page, start with
[`../research/stateful-ha-evidence.md`](../research/stateful-ha-evidence.md).

This page exists because the user is not merely asking for URLs to keep
landing somewhere.

The harder question is:

> if one node dies, does the system still know where truth lives, who may
> write, how clients rediscover the valid topology, and whether the surviving
> path is trustworthy?

That is where most infrastructure writing starts cheating. It celebrates:

- a still-routed hostname
- a live proxy
- a restarted container
- a reachable standby node

and quietly skips the only question that actually matters:

> did the system preserve the authoritative data path, or did it only preserve
> the appearance of continuity?

That distinction is the whole reason this page exists.

It also exists because this is the place where otherwise honest docs most
easily start flattering themselves.

Ingress progress is visible.
Stateful truth is harder, slower, and uglier.

That asymmetry is exactly why this page has to stay stricter than the rest of
the site.
If it ever starts sounding like a graceful addendum to routing work, it is
probably already hiding the real wound.

The repo can absolutely improve ingress before it earns honest stateful HA.
The danger is letting early routing progress create emotional overconfidence in
the data plane.

## What this page is and is not allowed to prove

This page is allowed to:

- define the boundary between ingress progress and real stateful correctness
- explain why stateful HA is a stricter claim than "reachable from more places"
- identify the missing data-plane capabilities that still block honest
  promotion
- keep the repo from overstating resilience just because edge behavior improves

This page is not allowed to:

- imply that proxying to a database equals high availability
- treat replicated access paths as the same thing as replicated authority
- upgrade route durability into write-path correctness
- pretend the repo has already solved storage portability, leader truth, or
  failover reconciliation unless the implementation proves it

## Priority evidence stack for stateful claims

Stateful claims on this site should be routed through this order:

1. [`../research/stateful-ha-evidence.md`](../research/stateful-ha-evidence.md)
2. root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
   and relevant active fragments
3. [`../operations/proof-matrix-and-drills.md`](../operations/proof-matrix-and-drills.md)
4. [`/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)

That order matters because stateful pages are easy to falsify in two opposite
ways:

- over-credit routing and proxy progress as data-plane truth
- over-credit future plans as if they already relocate authority

This page stays honest only if current evidence outranks planned elegance.

## Quick claim router

If the question is:

- "Why is stateful HA still the hardest part?" this page is a primary answer.
- "Can ingress improvements be real even if the data plane is immature?" yes,
  and this page explains that boundary.
- "Does improved routing mean Redis, Postgres, Mongo, or object storage are now
  honestly HA?" no.
- "What has to become system-owned before stateful HA claims are real?" this
  page exists to answer that precisely.

## The frustration beneath the stateful question

The user is already exhausted with “just put it behind a proxy,” “just add DNS
records,” and “just move the container” answers.

Those answers are attractive because they produce visible progress quickly.

They are also exactly how a stack starts lying to itself:

- the port is still open
- the dashboard still looks green
- another node can technically answer
- therefore the SPOF must be gone

For stateful services, that logic is false often enough to be dangerous.

This repo is unusual because the planning layer already resists that shortcut.
The architecture docs have to resist it too.

The user is not asking for caution as a stylistic preference.
They are asking the docs not to join the fraud.

That means this page should be read as a claim boundary, not a future-features
wish list.

Its main job is to stop the rest of the knowledgebase from silently upgrading:

- reachability into authority
- liveness into correctness
- proxy success into data-plane continuity

## The hard rule this repo is trying to defend

For any state-bearing service, you do not honestly remove the SPOF by merely:

- exposing it through Traefik
- giving it a stable hostname
- copying the container definition to another node
- teaching one node to TCP-proxy to another node
- proving that another box can answer on the same port

Those things may improve:

- reachability
- recoverability
- operator ergonomics
- demo quality

They do **not** by themselves produce stateful HA.

Real stateful HA requires some combination of:

- replicated data across failure domains
- explicit authority or election semantics
- client discovery or reconnect behavior that survives promotion
- storage semantics that do not strand the workload on a dead node
- failover behavior that does not require improvised operator rescue

In this repo, that last item matters a lot.

The user is exhausted with systems that only sound adult until the important
node disappears, at which point the real recovery plan is:

- remember the topology
- remember which survivor should win
- remember which clients need to move
- remember which hostname is now lying

That is not HA.
That is human failover wearing system language.

That phrase should color every stateful claim in the repo.

If recovery still depends on a human privately remembering topology truth, then
the control plane is still partly social, even if the hostname stayed pretty.

That is the real argument behind
[`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md).

## Why routing progress and stateful progress diverge so sharply

This repo can meaningfully improve ingress and still remain brutally exposed in
stateful domains.

That asymmetry is both useful and dangerous.

Useful because:

- better ingress can improve user-visible survivability early
- wrong-node HTTP recovery can solve real pain before the entire platform is
  “done”
- operator-facing services can become much less brittle even before deep data
  HA exists

Dangerous because:

- a service can still look “up” while its only real copy of data lives on one
  dead node
- peer-forwarding can preserve the appearance of continuity while truth remains
  singular
- a good demo can create more false comfort than a plainly broken system

This is why "the hostname still works" can be actively misleading.

For stateless HTTP, that may be real progress.
For stateful systems, it can mean the most important failure has merely been
hidden behind a still-responding edge.

That is why this page keeps hitting the brakes.

## Automatic disqualifiers for fake stateful closure

For this repo, a stateful answer has not earned honest HA language if it
mainly depends on:

- a stable hostname without replicated authority
- TCP reachability without explicit write ownership
- standby presence without client rediscovery semantics
- container mobility without durable storage portability
- operator memory to decide which survivor is actually valid

Any one of those can still be useful operationally.
None of them are enough to justify flattering stateful resilience language by
themselves.

It also explains why this page has to preserve uneven maturity instead of
tidying it away.

Some ingress paths may legitimately improve earlier.
Some stateful surfaces may remain intentionally crude or explicitly
single-writer for a while.

That unevenness is not a documentation defect here.
Pretending it has already resolved would be the defect.

## Strongest honest current answer

Ingress and peer-forward behavior can improve materially before the stateful
story becomes honest, but that does not make the stateful problem smaller. The
real barrier is still replicated authority: durable placement truth, write-path
rules, leader or primary semantics where needed, storage portability, and
recovery behavior that survives actual node loss instead of merely surviving a
proxy hop.

## What the user actually wants from the stateful story

The user is not asking for stateful services to sound more distributed.

They are asking for a much stricter answer:

> if the machine I unconsciously rely on disappears, does the data plane still
> have an explicit answer for where truth lives, who may write, how recovery
> works, and what is still unsafe?

That is the standard that keeps this page useful.

The most frustrating fake option in this space is the one that says:

- the hostname still resolves
- the TCP port is still open
- another node is online
- therefore the SPOF must be gone

That fake option is exactly what this page is trying to outlaw in the
knowledgebase.

## Stateful HA is four different problems, not one

The docs become mushy when they talk about “the database” or “stateful
services” as if that were one concern.

It is at least four.

### 1. Replication

Question:

- does the data exist in more than one meaningful failure domain?

If the answer is no, the stack still depends on one copy of reality. That is
not resilient state. It is a nicer route to one failure domain.

### 2. Authority and election

Question:

- who is currently allowed to accept authoritative writes, and how does that
  authority move?

Without an explicit answer, failover becomes:

- split-brain
- stale reads
- writes to the wrong place
- confused operator rescue

### 3. Client discovery

Question:

- how do applications discover the valid topology after failure or promotion?

This is why a pretty hostname is not enough for systems like MongoDB,
primary/replica Redis, or clustered queues.

### 4. Storage portability

Question:

- if the workload moves, does the durable state move safely, remain available,
  or stay replicated at all?

This is the quiet place where node-local bind mounts can shrink the whole dream
back down to one machine with better marketing.

That sentence should be read directly against the priority runtime.

The root Compose file already contains real local bind-mounted state for
services such as MongoDB and Redis.
So this page is not warning about an abstract future risk.
It is naming a current boundary in the priority implementation.

## The operator reading model this page wants to enforce

Every time someone says “this service is HA now,” this page wants four blunt
questions forced in front of them:

1. Where is the data replicated?
2. Who currently owns writes?
3. How do clients find the valid surviving topology after failure?
4. What still depends on one host path, one node-local volume, or one manual
   rescue action?

If those answers are missing, the honest language is weaker than “HA,” even if
the service remains reachable.

## Quick maturity matrix for the current repo

This table is deliberately blunt.

| Service class | Live worktree truth | Target direction in repo plans | Honest current language |
| --- | --- | --- | --- |
| Redis | Live, durable, node-local, TCP-exposed | Sentinel-style primary/replica with promotion-aware access | Single-writer durability with future HA direction, not solved HA |
| MongoDB | Live, durable, node-local, TCP-exposed | Replica Set semantics and topology-aware clients | Real state dependency with sharper intended future than current proof |
| Postgres | Live in some app domains, still storage-sensitive | Careful, maybe managed or narrowly self-hosted HA | Important state; repo is right to stay cautious |
| RabbitMQ | Live dependency in app surfaces | Clustered or explicitly bounded message durability | Stateful message backbone, not disposable glue |
| Bind-mounted service data generally | Widespread node-local assumption | Selective replication or explicit node affinity | Hidden SPOF pressure routing cannot erase |

The reading rule for this table is:

the “honest current language” column matters more than the “target direction”
column.

That rule is broader than the table.

Across the whole site, "honest current language" has to outrank:

- good future topology
- sophisticated failover vocabulary
- and any tempting sentence that sounds calmer than the evidence really is

The repo already has many target directions. The thing that makes the
knowledgebase worth trusting is refusing to narrate future topology as present
safety.

## Service-by-service reading of the current direction

The repo does not need every stateful service to have the same future.
It does need each one to be spoken about honestly.

That unevenness is not a weakness.
It is the only adult way to describe state.

Some services may eventually justify replication and promotion logic.
Some may justify explicit node affinity instead.
Some may justify managed externalization.

Uniformity is not the goal.
Truth is.

### Redis

The planning direction points toward:

- one writable primary
- one or more replicas
- Sentinel-style election
- optional HAProxy or equivalent in front of the active writer

Why that direction fits this repo:

- many workloads want one logical Redis more than they want Redis Cluster
  complexity
- Sentinel-style failover can preserve a single writable contract while
  removing one sacred data node

What is not enough:

- one Redis behind DNS failover
- two Redis containers with no election
- TCP reachability narrated as write continuity
- a standby that still requires manual re-pointing during stress

What would count as honest proof:

- active replica topology
- live election or promotion evidence
- client reconnection or retry-path proof
- confirmed write correctness after failover

These are not paperwork requirements.
They are the concrete evidence needed to show that the data plane still knows
where authority lives after failure.

What “confirmed write correctness” means here is not vague:

- the old writer is no longer ambiguously writable
- clients discover or are routed to the new truth cleanly
- writes after failover land on the authoritative survivor
- operators do not have to improvise the meaning of the cluster under stress

### MongoDB

The planning direction is sharper than the live proof:

- Replica Set semantics
- multiple members
- proper seed-list or SRV discovery

Why this matters:

- MongoDB clients survive primary changes only when topology is expressed
  correctly
- blind TCP forwarding is not the same thing as client-aware primary movement

What is not enough:

- one MongoDB container behind a stable hostname
- one TCP router whose backend happens to change
- another node that could start MongoDB if the operator intervenes

What would count as honest proof:

- a live Replica Set
- verified election behavior
- clients using topology-aware connection strings
- application behavior confirmed after primary movement

This is one of the clearest examples of why “TCP failover” is too weak a phrase
for this repo.

MongoDB does not merely need reachability. It needs topology-aware continuity.

That phrase applies more broadly than MongoDB.

The user is not trying to keep ports open.
They are trying to keep truth rediscoverable.

### Postgres

The repo is cautious here on purpose, and that caution is mature.

The effective stance is:

- use managed Postgres where that reduces chaos honestly
- or accept single-node Postgres until the surrounding control surface is less
  fragile

Why that is reasonable:

- Postgres HA is real and possible
- but it immediately drags in replication, fencing, promotion policy,
  reconnect behavior, backup discipline, and operator burden

That burden matters because the repo is explicitly comparing every new layer
against a simpler question:

> did this layer remove more hidden operator rescue than it introduced?

In this repo, pretending Postgres HA is easy would create more false comfort
than real resilience.

That caution is not a lack of ambition. It is one of the few places where the
repo already behaves like an adult system: it refuses to market a hard state
problem as if it were just another proxy exercise.

### RabbitMQ

RabbitMQ should be treated as real stateful infra, not as disposable glue.

Questions that matter:

- what persistence mode is acceptable?
- does queue leadership need clustering?
- what message-loss behavior is acceptable during failover?
- what retry behavior do clients need when topology changes?

The point is not that every queue needs the same answer.
The point is that "it is only a queue" is exactly how state risk gets
trivialized until it becomes production truth.

What is not enough:

- “the container restarts”
- “another node can reach it”
- “it is only a queue”

### Bind-mounted and storage-path-heavy services

Any service depending heavily on `${CONFIG_PATH}` or similar node-local paths
is still constrained by the node’s storage truth unless:

- the storage is intentionally shared
- the storage is intentionally replicated
- the service is explicitly declared node-affine and tolerated as such

This is one of the largest hidden SPOF boundaries in the repo.

## Node-local bind mounts are simple and dangerous at the same time

Why they remain attractive:

- simple
- inspectable
- Docker-native
- easy to reason about locally

Why they remain a constraint:

- relocation is hard
- failover can strand state
- node death still matters more than the proxy layer makes it appear
- backup and restore semantics remain strongly node-shaped

This does not mean “node-local storage is forbidden.”

It means the docs should stop letting node-local storage hide behind
distributed-looking ingress.

Many anti-SPOF narratives fail right here:

- the network path becomes distributed
- the storage path remains spiritually single-machine
- the architecture starts looking modern while still depending on one sacred
  disk layout or one manual rebuild sequence

## Shared or replicated storage is not universal absolution

Possible future directions in the archive and plans include:

- CephFS
- GlusterFS
- NFS in narrow bounded cases
- replicated block storage with deliberate single-writer policy

What this page should stop people from doing:

- defaulting to “just use distributed storage everywhere”

Shared storage solves some problems and introduces others. It is not a universal
forgiveness layer for weak topology design.

## Practical sequencing for this repo

The planning layer implies a sensible order, and the docs should say it
plainly.

### 1. Classify state honestly

Each service should be described as one of:

- truly stateless
- stateful but intentionally node-local for now
- stateful and important enough to justify real HA work

Without this, every discussion turns into vague maximalism.

### 2. Stop inflating ingress progress into data progress

There is limited value in claiming the stateful story is almost solved while:

- peer-aware routing is still maturing
- placement truth is still incompletely materialized
- cross-node convergence is still weak

### 3. Promote the most central stores first

The clearest first-class HA candidates remain:

- Redis
- MongoDB

Why:

- they are central enough to matter broadly
- their target topologies are already described more concretely than most of
  the rest of the stack

### 4. Treat universal shared storage as a later, narrower decision

Some services may need shared or replicated storage later.

That does **not** mean the whole repo should commit to one storage doctrine
prematurely.

### 5. Keep recoverability and anti-SPOF separate

The repo can still improve meaningfully before every critical state path is
fully resilient.

Examples of valuable but weaker outcomes:

- better backups
- faster node rebuilds
- explicit node affinity
- clearer promotion runbooks
- bounded failover on selected stores

Those are worth having.

They are just not the same claim as true anti-SPOF state.

## What “zero SPOF” should mean here

For a state-bearing service, that phrase should only be used when all of the
following are true:

- data is replicated across failure domains
- authority or election semantics are explicit
- failover does not depend on improvised manual rescue
- clients can discover or reconnect to the valid topology
- the intended node-loss scenario preserves the service contract honestly

If those conditions are not met, the honest language is weaker:

- improved recoverability
- bounded node affinity
- partial failover
- better reachability

Those outcomes still matter. They just are not the same thing.

This is the language discipline the rest of the repo should inherit:

- use `zero SPOF` rarely
- use `recoverable`, `node-affine`, `partially failed over`, or
  `bounded-single-node` when those are the truer statements
- prefer a sharp smaller truth over a grander fuzzy one

## Current truth from the worktree

The current worktree supports these claims:

- multiple real state-bearing services are already part of the priority runtime
- the planning layer explicitly separates stateful HA from ingress HA
- Redis and MongoDB have clearer future directions than the live stack proves
- node-local storage assumptions still shape the real runtime strongly

The current worktree does **not** prove:

- replicated topologies for all critical stores
- universal client failover correctness
- a solved shared-storage strategy
- broad end-to-end zero-SPOF claims for stateful services

It also does not yet prove that an operator can lose a primary node and still
derive authoritative write safety entirely from tracked system truth rather
than remembered rescue procedure.

That is the honesty wall this page exists to preserve.

## Bottom line

The dream in this repo is not “make the dashboard stay green after a container
dies.”

The dream is:

> stop worshipping one machine as the place where the real system lives

For stateless services, routing can get part of the way there relatively
quickly.

For stateful services, routing only gets to the beginning of the real work.

If the docs forget that, they stop documenting `bolabaden-infra` and start
telling the same fake HA story the user is already tired of hearing.
