# Stateful HA Plan: Where the Repo Refuses to Lie

This page is about sequencing, not product shopping.

It asks a much stricter question:

> Which parts of stateful resilience is `bolabaden-infra` actually prepared to
> solve honestly first, and which parts must stay explicitly deferred until the
> surrounding ingress, placement, and failover layers stop being fragile?

That question matters because stateful HA is where infrastructure docs most
often become fraudulent without meaning to.
It is also where the repo would be easiest to "finish" on paper while leaving
the real pain intact.
For stateful systems especially, the docs need to preserve the stricter rule
behind this rewrite:

- do not compress the hard problem into a cleaner adjacent problem
- do not discard contradiction because it complicates the story
- do not turn partial likeness into authoritative sameness

## Why this page exists

The user is clearly frustrated with fake availability stories.

Stateful systems are where the fake version of HA is easiest to accidentally
sell:

- the container restarts
- the port answers
- the dashboard turns green
- the proxy route exists

but:

- write authority is unclear
- client topology assumptions are wrong
- storage is still node-local truth
- promotion behavior is untested
- failover correctness is wishful

This page exists to keep the repo from telling that lie.
It also exists because this is one of the few places where the user is
explicitly demanding that the docs not help the system gaslight itself.

That sentence matters because stateful infrastructure is where the whole
project could most easily retreat into fake adulthood:

- more replicas
- more ports
- more diagrams
- more cluster words

while the operator still privately knows that one machine, one disk path, or
one manual recovery ritual remains the real authority.

This page should keep refusing that retreat.

## Evidence classes used here

This page blends:

- live-root evidence from the active Compose-first runtime
- repo-direction evidence from planning and architecture docs
- stateful-side-path evidence from L4 ingress and OpenSVC-related work

It does **not** claim that the entire root stack already has broad stateful HA.

## The strongest principle the repo already understands

The single most important principle in the current docs is:

> moving a container or preserving a port is not the same as making state
> survive node loss.

That instinct is consistent across:

- [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
- [`knowledgebase/architecture/stateful-ha-and-data.md`](../architecture/stateful-ha-and-data.md)
- [`knowledgebase/research/stateful-ha-evidence.md`](stateful-ha-evidence.md)

It is one of the healthiest instincts in the entire repo, because it directly
resists shallow HA theater.
That instinct is more than a caveat.
It is one of the repo's clearest philosophical boundaries:

the system is allowed to become more reachable before it becomes more honest
about state only if the docs never confuse those two kinds of progress.

## The repo's actual near-term stateful priority

The repo does **not** read like it wants to solve every distributed-storage
problem immediately.

It reads like it wants to do something much more disciplined:

1. make ingress and wrong-node behavior less fragile
2. make placement and failover truth less dependent on human memory
3. choose a small number of critical stateful systems to harden properly
4. refuse to pretend that every bind-mounted service is now anti-SPOF

That is not hesitation. It is sequencing discipline.
It is the repo refusing the usual trap where the desire for a complete story
causes people to overstate what the first stateful hardening step will buy.

## Why sequencing matters so much here

If the repo tries to solve all of these at once:

- any-node ingress
- peer forwarding
- failover route generation
- service relocation
- replicated storage
- app-specific state semantics

then every category becomes harder to reason about, and the system starts lying
through sheer complexity overload.

The current plan is strongest when it refuses that trap.

It also lines up with the user's deeper frustration about lack of genuine
options.
Stateful infrastructure is one of the clearest places where the ecosystem
pretends there are only two choices:

- stay naive and node-local
- or immediately adopt the heaviest distributed answer in sight

This repo is trying to make a narrower, more service-specific path visible
instead.

## Current implementation evidence that shapes the plan

The live runtime already includes meaningful state-bearing services directly in
the priority implementation and active fragments, including:

- `redis`
- `mongodb`
- `rabbitmq`
- `litellm-postgres`
- `nuq-postgres`
- many bind-mounted app and config paths under `${CONFIG_PATH}`-style storage

The repo also contains a separate L4 ingress path:

- [`compose/docker-compose.l4-ingress.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.l4-ingress.yml)
- [`scripts/osvc_l4_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_l4_sync.py)

That matters because it proves the repo is already distinguishing between:

- "can a client reach a TCP endpoint?"
- and "is the data model behind that endpoint resilient?"

Those are not the same question.
That distinction should keep reappearing because the user has already seen too
many toolchains erase it.

## The subproblems this plan has to keep separate

If future docs collapse these into one checkbox labeled "HA," the repo will
regress.

## 1. Ingress reachability

Can the client reach a live endpoint at all?

This is the most superficial layer.

It matters, but it proves the least.

## 2. Role discovery

Can the client or frontend path discover the correct:

- primary
- master
- replica set seed list
- valid backend role

Without role discovery, a reachable port can still point at the wrong truth.

## 3. Data continuity

Does the authoritative dataset survive node loss and continue from the right
source?

This is the real stateful HA question.

## 4. Storage substrate truth

Is the storage implementation consistent with the failover story being told?

Examples:

- bind mounts on one node do not become resilient because proxy logic got
  smarter
- shared filesystems impose their own failure and correctness semantics
- database-native replication and quorum are different from filesystem
  replication

## 5. Client failover behavior

Even if the backend topology is correct, do the actual clients reconnect,
rediscover, and continue correctly?

This is one of the easiest layers to forget and one of the easiest ways to
declare success too early.

## Redis: why it reads like an early serious candidate

The repo points toward a Redis path involving:

- Sentinel-style master election
- HAProxy or L4 routing toward the active master

That is a pragmatic fit for the repo's actual desire:

- one logical Redis entrypoint
- automatic primary promotion
- understandable L4 routing semantics

It is attractive because it moves the repo beyond single-node Redis without
forcing immediate adoption of a more complex distributed data model.

### What that path could prove if actually implemented and tested

- a single node loss does not inherently destroy Redis availability
- write traffic can converge toward a promoted master
- clients do not require manual rewiring after every failover event

### What it would still not prove by itself

- that every Redis-using application in the repo tolerates failover correctly
- that every client reconnect strategy is safe
- that persistence and replica lag tradeoffs have been accepted consciously

Redis is therefore a good early candidate, but not a universal proof of stateful
HA maturity.

## MongoDB: why the repo treats it differently

The repo's Mongo direction implies replica-set semantics, not just frontend
port failover.

That is the correct instinct.

MongoDB resilience is not honestly solved by:

- a VIP
- a shared port
- or a simple TCP switch

MongoDB needs topology-aware thinking:

- multiple members
- client seed-list or SRV awareness
- election behavior
- replica lag and write concern decisions

This is healthier than pretending that a dumb front door makes Mongo resilient.

## Postgres: why the repo stays cautious

The repo's tone around Postgres is more restrained, and that restraint is one
of the better signs in the whole knowledgebase.

Postgres HA usually imports heavier machinery:

- Patroni-style coordination
- repmgr-style promotion logic
- consensus layers such as etcd or similar
- sharper correctness risks around promotion, lag, and split-brain handling

The repo appears to understand that Postgres is not the best place to perform a
premature victory lap.

That caution encodes a valuable priority judgment:

> do not enter the most complex stateful HA domain first while ingress,
> placement, and convergence truth are still under active repair.

## RabbitMQ and message systems: the implied caution

Even where a broker can be reached through L4 or forwarded through another node,
that still leaves major questions open:

- queue durability
- mirrored or quorum queue behavior
- reconnect semantics
- consumer topology assumptions

These systems are not solved merely because they fit behind a stable port.

## Files and bind-mounted application state: the plan's hardest honesty test

This may be the hardest truth in the repo.

Many services can look portable at the container layer while their real identity
still lives in:

- bind-mounted config
- local media or app state
- node-local databases
- mutable files under `${CONFIG_PATH}` and similar trees

That means much apparent multi-node flexibility still collapses back into local
disk truth.

The repo is right to resist pretending otherwise.
This is probably one of the deepest frustrations under the whole project:
the stack can look wide, modern, and peer-capable while its most important
state still collapses back into one host path and one remembered rescue path.

That is also why "more options" has felt fake to the user.
Many supposed options only diversify the front-door story while leaving the
real state authority exactly where it already was.

### Why the repo is right to defer universal distributed-filesystem ambition

Trying to solve all of this simultaneously would mean solving:

- ingress
- placement
- failover routing
- replicated storage
- app-specific state semantics

in one jump.

That is exactly how fake completeness happens.

## The L4 ingress path: what it changes and what it does not

The presence of:

- [`scripts/osvc_l4_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_l4_sync.py)
- [`compose/docker-compose.l4-ingress.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.l4-ingress.yml)

shows that the repo is trying to expose TCP services behind generated HAProxy
frontends and Tailscale-addressed backends.

That is useful.

But this path only changes:

- how a client reaches a backend port
- how multiple backends can be presented behind one logical TCP surface

It does **not** by itself solve:

- primary election
- replica consistency
- write continuity
- storage durability
- application-aware recovery logic

That distinction should remain brutally explicit in every future stateful page.

## The sequencing judgment the current plan encodes

The repo appears to be making a deliberate priority call.

### Solve earlier

- honest ingress behavior
- wrong-node routing that does not collapse into theater
- placement awareness
- opt-in TCP frontend generation
- targeted hardening of Redis and MongoDB first

### Treat cautiously or defer

- universal shared storage stories
- broad Postgres HA complexity
- stack-wide claims that all bind-mounted services are now resilient
- pretending every TCP-reachable service has become state-safe

This is exactly the kind of sequencing the repo should be making.
It is also a reminder that progress has to stay service-specific here.
One hardened route or one replicated datastore should not bleed into a
stack-wide stateful maturity claim.

And one promoted tool should not be allowed to claim ownership of every
stateful burden just because it participates in the path.

The repo should keep asking, per service:

- who owns authority
- who owns discovery
- who owns promotion
- who owns storage truth
- who owns client reconnection assumptions

If those answers are still fragmented or private, the docs should say so even
when the surrounding stack looks much more impressive.

## What stronger proof would be required before broad stateful HA claims

Before the repo can claim broad stateful HA, it would need stronger evidence of:

- live replicated deployments for selected datastores
- tested promotion and rejoin behavior
- client configurations that actually tolerate failover
- storage semantics that survive node loss
- operational runbooks for split-brain, stale replica, rejoin, and resync cases
- proof that ingress continuity and state continuity are aligned rather than
  accidentally decoupled

Without that, "stateful HA" should remain:

- a plan category
- a service-specific claim
- or a tightly scoped experiment

not a broad description of the stack.
The repo needs language that can tolerate uneven maturity without treating that
unevenness as failure.
Otherwise every partial success gets narratively inflated just to avoid
admitting the remaining wound.

## What this page is trying to prevent

This page exists to prevent four specific lies.

### Lie 1: reachable port equals resilient service

It does not.

### Lie 2: failover proxy equals stateful continuity

It does not.

### Lie 3: database replicas equal application-level correctness

They do not by themselves.

### Lie 4: one hardened datastore path means the repo has solved stateful HA in general

It has not.

## Bottom line

The stateful HA plan is strongest when read as a boundary-setting prioritization
document, not a shopping list. It shows that the repo already understands the
difference between reachable endpoints and durable state, and that it is trying
to harden the most important service classes first instead of pretending all
local data can become cluster-safe overnight. That is one of the most valuable
forms of honesty in `bolabaden-infra`, and the docs should preserve it rather
than smoothing it away.
