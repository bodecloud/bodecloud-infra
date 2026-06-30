# Stateful HA Evidence

This page is the proof boundary for the repo’s stateful story.

This is the part of the architecture where "multi-node" becomes least
meaningful on its own.
You can have several nodes, a working proxy, and even plausible failover demos
while still depending on:

- one write authority
- one node-local volume
- one unreplicated datastore
- one reconnect path that breaks badly after promotion

This page exists to stop the docs from calling that situation "HA."

It also exists because state is where the fake-option problem becomes most
expensive.

At the ingress layer, a weak answer may still produce a misleadingly intact
demo.
At the stateful layer, weak answers can preserve the interface while losing the
 place where authority, durability, or correctness actually lived.

## What this page is trying to prove

This page is not trying to prove that stateful HA is solved.

It is trying to prove seven narrower things:

1. the live stack already depends on real stateful services
2. those services still lean heavily on node-local persistence
3. the planning layer clearly separates ingress continuity from data
   correctness
4. the intended stateful topologies are sharper than the live runtime currently
   proves
5. the hardest blockers are replication, authority, client discovery, and
   storage portability together
6. stateful risk is already present in the priority runtime, not just in future
   architecture discussions
7. the repo is unusually honest about stateful HA, but honesty is not the same
   thing as proof

## What the user is actually asking this page to police

The user is not asking whether a service can be restarted.
The user is asking whether the system can lose a node without losing the place
where truth lives.

That distinction is why this page exists.
Without it, almost any routed or restarted stateful service can be narrated as
"good enough HA" long before the hard parts are addressed.

That is exactly the kind of downgrade this repo is trying to stop.

The user is not asking whether a node can come back.
The user is asking whether truth itself stopped living in one brittle place.

That is why stateful documentation has to stay somewhat repetitive.
This is the place where the stack can look most alive while the real authority
still lives on one disk path, one elected node that is not actually
replaceable yet, or one reconnect story that only works in the flattering
sequence.

## Evidence classes this page relies on

### Class 1: live implementation evidence

Used for:

- which stateful services are already part of the root runtime
- where node-local storage assumptions are visible
- which services expose TCP or depend on durable backing systems

### Class 2: repo-native intent evidence

Used for:

- the repo-wide honesty boundary that ingress cannot fake data correctness

### Class 3: planned architecture evidence

Used for:

- Redis Sentinel direction
- MongoDB Replica Set direction
- cautious posture toward Postgres HA
- separation between ingress failover and stateful correctness

### Class 4: archive-pressure evidence

Used for:

- why fake HA language is unacceptable
- why node survival is not enough if state authority is still singular

## Strongest live stateful anchors

Primary files:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
- [`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)
- [`compose/docker-compose.llm.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.llm.yml)
- [`compose/docker-compose.stremio-group.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.stremio-group.yml)

## What the live files concretely prove

### 1. MongoDB is live and still anchored to node-local storage

[`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
shows `mongodb` with:

- bind mount to `${CONFIG_PATH:-./volumes}/mongodb/data`
- Traefik TCP exposure
- local healthcheck

What that proves:

- MongoDB is a real dependency of the priority runtime
- persistence is still tied to one node-local storage path

What it does not prove:

- Replica Set semantics
- promotion safety
- client-aware failover correctness

### 2. Redis is live and still aligned to single-instance local durability

The same file shows `redis` with:

- bind mount to `${CONFIG_PATH:-./volumes}/redis:/data`
- password-backed single-instance assumptions
- Traefik TCP exposure
- local healthcheck

What that proves:

- Redis durability and TCP exposure are already part of the real runtime

What it does not prove:

- Sentinel
- replica promotion
- master-aware routing
- reconnect correctness after failover

### 3. Firecrawl expands the stateful surface beyond Redis and MongoDB

[`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)
shows:

- `nuq-postgres` with bind mount to
  `${CONFIG_PATH:-./volumes}/nuq-postgres/data`
- `rabbitmq`
- `firecrawl` depending on Redis, Postgres, and RabbitMQ

What that proves:

- the stateful surface is already broader than the obvious core services
- application continuity in this repo already depends on multiple durable
  backends

What it does not prove:

- that those backends already have HA-capable authority and storage semantics

### 4. Stateful dependency is distributed across domains, not isolated to one corner

Even without claiming every service in `llm` or `stremio-group` has equal HA
priority, the live runtime already shows that durable data and node-local
storage assumptions are not confined to one tiny subsystem.

What that proves:

- stateful correctness is already a present architecture problem

What it does not prove:

- which domains should be promoted first
- which services are intentionally allowed to remain single-node

## What the live files do not prove

The live runtime does not, by itself, prove:

- Redis Sentinel or equivalent leader election
- MongoDB Replica Set health and client topology correctness
- Postgres promotion, fencing, or reconnect behavior
- storage portability beyond node-local bind mounts
- service-class-specific failover maturity

That missing proof is the main reason stateful HA has to stay separate from
ingress rhetoric.

Another way to say it:

- the current worktree proves that the repo already depends on state
- it does not prove that state has stopped depending on sacred infrastructure

That second line is what prevents a routed hostname from being mistaken for
stateful resilience.
A stable endpoint can hide a fragile authority model extremely well.
The docs therefore need to keep asking not just "can something answer?" but
"did the place where correctness lives actually stop being singular?"

That second line is the heart of the page.

Several nodes, several routes, and several helper layers still do not add up
to stateful anti-SPOF unless authority, replication, discovery, and storage all
stop collapsing into one remembered or one local place.

## Strongest intent and planning anchors

Primary files:

- [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

## What these files explicitly say

### 1. Stateful HA is not container relocation

[`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
explicitly says:

- zero-SPOF stateful services cannot be achieved by merely moving containers
- replication plus quorum are required
- raw TCP should not be conflated with hostname-routed HTTP

What that proves:

- the repo has a strong conceptual boundary around stateful truth

What it does not prove:

- that the runtime has already crossed it

This distinction matters because a repo can sound unusually wise about stateful
HA while still remaining operationally dependent on the exact SPOFs it is able
to describe so accurately.

### 2. Redis and MongoDB have named target topologies

The same planning layer explicitly points toward:

- Redis Sentinel plus master-aware routing
- MongoDB Replica Set semantics

What that proves:

- the repo has concrete opinions about what a less-fake stateful future looks
 like

What it does not prove:

- that those topologies are active now

That means the planning layer is doing something valuable but incomplete:
it is reconstructing the correct shape of the problem, not yet proving the live
runtime has inherited that shape.

### 3. Postgres is treated more cautiously than the rest

The planning layer resists pretending Postgres belongs in the same confidence
bucket as simpler stateless or lightly stateful failover work.

What that proves:

- the repo is not shouting "HA everything" without class-specific honesty

What it does not prove:

- which Postgres path will ultimately be adopted

### 4. README keeps the stateful honesty wall visible

[`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
explicitly warns that ingress cleverness cannot fake:

- replication
- election
- quorum
- reconnect behavior
- durable data safety

What that proves:

- the repo already understands the most common stateful documentation lie

What it does not prove:

- that those concerns are solved service by service

## Critical negative evidence

### Negative fact 1: live stateful dependency exists before universal HA proof

The worktree already proves that durable services are live now.

It does not prove:

- universal replication
- universal promotion logic
- universal client-topology correctness
- universal storage portability

This mismatch is one of the most important honesty boundaries in the repo.

It is also one of the easiest places for a broader platform story to cheat.

A scheduler, a cluster, or a nicer ingress path can make the stack feel more
distributed while leaving state authority just as singular as before.

### Negative fact 2: bind mounts still tie important services to one failure domain

Current live examples include:

- MongoDB data under `${CONFIG_PATH}/mongodb/data`
- Redis data under `${CONFIG_PATH}/redis`
- Postgres data under `${CONFIG_PATH}/nuq-postgres/data`

What that proves:

- node loss still has deeper consequences than ingress docs alone can express

What it does not prove:

- that the repo has no recovery path at all
- only that shared or replicated storage truth is not currently proved by the
  root runtime

### Negative fact 3: recommended topology is clearer than deployed topology

The planning layer already has strong opinions:

- Redis should grow into Sentinel-style failover
- MongoDB should grow into Replica Set semantics
- Postgres should be treated carefully or managed externally

What that proves:

- the repo’s idea of "better" is not vague

What it does not prove:

- which of those paths are already live
- which stateful services are intentionally single-node
- which migrations should happen first

### Negative fact 4: current root evidence is stronger at exposing the problem than solving it

The live runtime plus planning layer give a very sharp diagnosis:

- state is real
- node-local storage is real
- ingress cannot fake data correctness
- Redis and MongoDB need topology-aware futures

What they do not yet give is a matching end-to-end proof chain for those
futures in the root runtime.

That is not failure of documentation.
That is the current state of the implementation.

The useful role of documentation here is therefore not to make the repo sound
more mature.
It is to keep the costliest unfinished truths visible.

## Claim ledger

| Claim | Evidence class | Confidence | What it actually proves | What it still does not prove |
|---|---|---|---|---|
| The live stack already has real stateful dependencies. | Live implementation | High | Redis, MongoDB, Postgres, RabbitMQ, and node-local persistence assumptions are already in the priority runtime. | Not that they already run in resilient topologies. |
| The repo explicitly treats stateful HA as a different problem than HTTP failover. | Planned architecture + repo-native intent | High | The planning layer cleanly separates ingress continuity from datastore correctness. | Not that the runtime enforces that distinction perfectly everywhere. |
| Redis needs native failover semantics, not just exposure behind ingress. | Planned architecture | High | Sentinel-style thinking is explicit in repo planning. | Sentinel or equivalent is not proven live in the root runtime. |
| MongoDB needs Replica Set semantics and client-aware discovery. | Planned architecture | High | The intended topology is explicit. | The current runtime does not prove a working Replica Set. |
| Postgres is intentionally treated more cautiously than simpler stateful classes. | Planned architecture + README | High | The repo is refusing shallow "just HA the database too" narration. | Which long-term Postgres HA path will actually win. |
| Node-local storage remains one of the biggest hidden SPOF truths. | Live implementation + planned architecture | High | Bind-mounted data paths are visible throughout the runtime and the planning layer acknowledges the problem. | Which exact replicated or shared-storage path will be adopted later. |

## The stateful subproblems the docs must keep separate

The evidence supports at least four separate stateful subproblems:

1. replication
2. authority and election
3. client discovery and reconnect behavior
4. storage portability

If any one of those is missing, the phrase "HA datastore" becomes much weaker
than it sounds.

This is the main downgrade filter the rest of the knowledgebase should inherit
whenever it is tempted to say that a stateful class is "basically covered."

## What evidence would be required before the claims can be upgraded

Before the docs could honestly narrate a service class as statefully HA, the
repo would need evidence that demonstrates all of the following for that class:

1. replicated or otherwise resilient data topology is active
2. authority or leadership transition is defined and testable
3. clients discover or reconnect correctly after promotion
4. storage semantics survive the intended failure mode
5. the resulting behavior is specific to that service class rather than implied
   from generic clustering language

Without that chain, the only honest label remains partial or planned.

And that is not a disappointing label.
It is the label that prevents the repo from copying the same false closure it
rejects elsewhere.

## Operational reading

### What is already real

- the stack already depends on multiple durable services
- the repo knows stateful HA is different from ingress failover
- the planning layer already has concrete topology opinions

### What is still incomplete

- service-specific replication proof
- service-specific election and promotion proof
- client behavior proof after failover
- storage portability beyond node-local bind mounts

## What this page should stop future docs from claiming

Future docs should not claim:

- a routed TCP endpoint equals datastore HA
- multiple nodes automatically protect data
- a scheduler alone solves authority or quorum
- the repo has no stateful SPOFs unless service-specific topology is proven

Those are some of the most expensive false claims the repo could make, because
they would hide the exact place where the user's dream still most needs sharper
truth rather than broader language.

## Reading discipline this page imposes on the rest of the docs

If another page wants to use words like:

- resilient
- HA
- failover-safe
- zero-SPOF

for a state-bearing service, it needs to answer all of these:

1. where is the replicated data proof?
2. where is the authority or election proof?
3. where is the client discovery or reconnect proof?
4. where is the storage portability or node-loss proof?

If it cannot answer those questions, the wording must be downgraded.

## Bottom line

The evidence says the repo already carries real stateful risk, and the repo
also understands that risk unusually well.

What the evidence does not say is that the hard part is finished.
The worktree proves stateful dependency.
The planning layer proves stateful awareness.
Neither one proves replicated, election-safe, client-correct,
storage-portable stateful HA today.

That is the real outcome of this page:

the repo has already learned to name the hard stateful wound correctly, but it
is still in the process of moving that understanding out of analysis and into
live topology truth.
