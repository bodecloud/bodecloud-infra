# Stateful HA and Data Ownership

For the deeper evidence boundary behind this page, start with
[Stateful HA Evidence](../research/stateful-ha-evidence.md).

This page exists because the user's real stateful question is harsher than:

> does the service still answer?

It is:

> if the current node dies, where does authoritative truth live, who is still
> allowed to write, how do clients rediscover that truth, and which parts of
> the answer are still secretly being carried by one operator?

The nearby smaller question this page must not collapse into is:

> which databases do we run?

That smaller question is not useless.
It is just nowhere near enough.
The user is not inventorying durable services for fun.
They are checking whether "multi-node" is still mostly ingress theater once
data authority matters.

## What this page is and is not allowed to prove

This page is allowed to prove:

- the live priority runtime already depends on real stateful services
- those services still reveal strong node-local authority and storage
  assumptions
- the repo already distinguishes ingress progress from stateful correctness
- the current worktree still falls far short of generic stateful anti-SPOF
- some service classes deserve explicit single-writer honesty rather than
  softened resilience language

This page is not allowed to prove:

- that TCP exposure equals HA
- that restartability equals authority preservation
- that one more node or one more bind mount equals failover dignity
- that future planning language makes the present runtime safe
- that a route staying up means the truth stayed up

## The failure mode this page is trying to stop

The most dangerous flattering sentence in a repo like this is:

> the service still worked

That sentence can hide at least four very different realities:

- the same node never actually failed
- the same local writer restarted
- a stale or secondary surface answered
- the only truthful copy of data was still singular and lucky

This page exists to stop stateful continuity from being narrated in transport
language.

## The hard rule this page is defending

For any state-bearing service, you do not honestly remove the SPOF merely by:

- routing it through Traefik
- giving it a stable hostname
- copying its service definition to another node
- restarting it elsewhere later
- proving that the port still answers

Those may improve:

- reachability
- operator ergonomics
- recovery options
- demo smoothness

They do not, by themselves, prove preserved authority.

## The private sentence stateful pages must kill

For state-bearing systems, the hidden sentence is almost always some variant
of:

> I still personally know who the real writer is, which copy is authoritative,
> and how the clients are supposed to recover.

If that sentence survives untouched, the SPOF has not been removed in the way
the user actually cares about.

At best, the repo may have improved:

- exposure
- recovery speed
- restart ergonomics
- operator visibility

Those are real gains.
They are not yet authority transfer.

## Strongest honest current answer

The strongest honest answer from the current worktree is:

- stateful risk is already live and significant
- authority is still mostly singular per service
- durability is still heavily node-local
- planning pressure is sharper than live proof

That means the repo already knows many of the right future questions.
It does not yet prove that the main stateful classes have crossed into
trustworthy multi-node behavior.

## The smallest honest stateful sentence versus the fake bigger sentence

This page should force the distinction between these two:

- smaller honest sentence:
  `this service is durable, important, and currently singular in authority`
- fake bigger sentence:
  `this service is effectively HA because more than one node can expose or reach it`

The first sentence may feel disappointing.
It is still much more useful than the second if the second becomes false on the
first real failure.

## What the live runtime already proves

The current priority runtime already depends on all of the following stateful
classes:

- MongoDB in the root runtime
- Redis in the root runtime
- Headscale with SQLite in the active Headscale fragment
- Postgres and RabbitMQ in the active Firecrawl fragment
- Postgres in the active LLM fragment
- Qdrant in the active LLM fragment
- multiple bind-mounted data paths under `${CONFIG_PATH:-./volumes}/...`

That matters because stateful pain is not hypothetical future-architecture
pain.
It is already present in the live stack.

## The four questions every stateful claim should answer

Every serious stateful claim in this repo should answer all four:

1. Where does authoritative truth live right now?
2. What happens to that authority if the current node disappears?
3. How do clients discover the new truth without human folklore filling the
   gap?
4. What prevents two competing truth-owners from both acting like the writer?

If those questions are still mainly answered by operator memory, the system is
still socially carrying part of the control plane.

## The stateful burden map

These are the main hidden burdens the user is trying to remove at this layer:

| Burden | What it really asks | What fake progress looks like |
| --- | --- | --- |
| Writer identity | who may safely accept writes right now? | a port answered |
| Authority continuity | who owns truth after failure? | the service restarted |
| Client rediscovery | how do dependents find the new truth? | DNS or TCP still points somewhere |
| Fencing | what stops the old writer from still acting alive? | the old node is probably down |
| Storage truth | did durable state move, replicate, or remain singular? | the compose file exists on more than one node |

This map is the reason stateful docs must stay much meaner than ingress docs.

## Stateful inventory by service class

### 1. MongoDB in the root runtime

Live evidence:

- root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  declares `mongodb`
- persistent data is bind-mounted at
  `${CONFIG_PATH:-./volumes}/mongodb/data:/data/db`
- the service is exposed through Traefik TCP labels

What this proves:

- MongoDB is a real dependency of the priority runtime
- durability is still tied to a node-local path
- transport exposure already exists

What it does not prove:

- Replica Set semantics
- election behavior
- write-authority promotion
- client topology rediscovery after failure
- fencing or split-brain prevention

Private sentence still surviving:

> yes, but I still personally know MongoDB authority is basically where that
> node and disk were

### 2. Redis in the root runtime

Live evidence:

- root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  declares `redis`
- persistent data is bind-mounted at
  `${CONFIG_PATH:-./volumes}/redis:/data`
- the service is exposed through Traefik TCP labels

What this proves:

- Redis is a live stateful dependency, not a conceptual helper
- the runtime still looks like a single durable Redis instance story

What it does not prove:

- Sentinel
- replica promotion
- correct client reconnect behavior
- resilient master discovery
- state authority transfer under node loss

Private sentence still surviving:

> yes, but I still personally know Redis truth is still mostly one writer

### 3. Headscale in the active fragment

Live evidence:

- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
  sets `database.type: sqlite`
- the database path is `/var/lib/headscale/db.sqlite`
- WAL mode is enabled
- the service is internet-facing through the edge stack

What this proves:

- the private-mesh control plane is currently anchored to a local SQLite file
- Headscale is a singleton authority concern today, even if the service is
  highly useful

What it does not prove:

- shared control-plane authority
- automatic leader replacement
- cluster-safe promotion
- mesh continuity after authority-node loss

This is one of the sharpest examples in the repo of the difference between:

- a live service
- a valuable service
- a resilient authority surface

Private sentence still surviving:

> yes, but I still personally know the mesh control plane still bottoms out in
> one SQLite-backed authority node

### 4. Firecrawl subgraph

Live evidence from
[`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml):

- application dependencies include:
  - `redis`
  - `nuq-postgres`
  - `rabbitmq`
- `nuq-postgres` persists under
  `${CONFIG_PATH:-./volumes}/nuq-postgres/data:/var/lib/postgresql/data`
- `rabbitmq` is a real live runtime service, not just a planned future

What this proves:

- the stateful surface is wider than the obvious root databases
- application continuity already depends on a graph of durable backends

What it does not prove:

- HA Postgres
- mirrored durable messaging semantics
- promotion-safe topology across the subgraph
- resilient multi-node truth for the whole application family

Private sentence still surviving:

> yes, but I still personally know this subgraph still depends on singular
> durable authorities

### 5. LLM Postgres and vector-state surface

Live evidence from
[`compose/docker-compose.llm.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.llm.yml):

- `litellm-postgres` is live
- its data persists under
  `${CONFIG_PATH:-./volumes}/litellm/pgdata:/var/lib/postgresql/data`
- `qdrant` is live
- its storage persists under
  `${CONFIG_PATH:-./volumes}/qdrant/storage:/qdrant/storage`

What this proves:

- more than one Postgres-backed or durable AI-adjacent surface already exists
- node-local persistence assumptions are repeating, not isolated
- statefulness in this repo is no longer only about classic infra databases

What it does not prove:

- HA Postgres for LLM-facing services
- promotion semantics
- reconnect correctness after node loss
- safe vector-store authority failover

Private sentence still surviving:

> yes, but I still personally know these AI-adjacent durable systems are still
> node-local truth holders

## What the live runtime still shows clearly

The live worktree still shows:

- persistence is mostly node-local
- writer authority is mostly singular
- replicated authority is not the default runtime truth
- client rediscovery behavior is not yet a generally proven surface

That is the shortest honest reason stateful claims must stay much harsher than
ingress claims.

The edge can become more survivable earlier than the data can become more
truth-owning.

## Why stateful progress often feels slower than everything else

The user is frustrated because many other layers offer partial relief quickly:

- more public entry nodes
- better proxy behavior
- cleaner auth and middleware
- nicer observability

Stateful truth is slower because it asks for something uglier and more
constraining:

- explicit writer rules
- explicit failure ownership
- explicit replacement authority
- explicit client reconvergence

That slower pace is not bureaucracy.
It is the cost of no longer lying about where truth lives.

## Why ingress progress and stateful progress diverge

This repo can improve user-visible continuity before it earns honest stateful
HA.

That asymmetry is useful because:

- wrong-node HTTP recovery can remove real operator pain early
- some services are mostly request-routing problems
- the edge stack already carries real policy, auth, and observability meaning

That asymmetry is dangerous because:

- the hostname can still work while the only authoritative copy of data lived
  on the dead node
- a proxy path can preserve appearance while authority remains singular
- a convincing demo can hide the fact that one disk and one writer still
  mattered more than the rest of the topology

That is why "the hostname still works" is almost meaningless by itself for
this layer.

## What still does not count as stateful progress

These are still weak or fake progress signals here:

- the same hostname still resolves
- the TCP port still answers
- the container restarts on the same node
- the service could be redeployed elsewhere later
- a second machine could theoretically mount similar storage
- a liveness check still passes
- a proxy can still reach something at the same address

Those may improve confidence.
They do not yet prove preserved authority.

## What a fake stateful upgrade usually sounds like

These are the kinds of sentences this page should aggressively reject:

- `MongoDB is redundant because Traefik can still forward TCP`
- `Redis is resilient because another node can run the same service definition`
- `Headscale is HA enough because the mesh usually works`
- `Postgres is portable because the data directory is bind-mounted`
- `RabbitMQ survived because the process came back`
- `Qdrant is fine because the index directory persists`

Every one of those may describe a nearby operational convenience.
None of them, by themselves, answer the authority question.

## Service-class reality check

### Redis-class systems

What matters:

- primary truth
- promotion rules
- client reconnect behavior
- storage durability
- safe discovery of the current writer

What is not enough:

- one more routed endpoint
- one more container on another node

### Mongo-class systems

What matters:

- replica-set semantics
- election behavior
- consistency model
- client topology awareness
- safe write-owner transition

What is not enough:

- TCP exposure through Traefik
- a local healthcheck

### SQLite-backed control-plane services

Headscale is the active example here.

What matters:

- whether the SQLite file is still fundamentally local authority
- whether a replacement node can become authoritative without operator rescue
- whether leadership is represented by a real system-owned rule

What is not enough:

- the admin surface responding
- the mesh working on a normal day

### Postgres-backed application subsystems

What matters:

- single-writer honesty versus actual HA
- promotion and fencing semantics
- storage portability
- reconnect behavior
- durable truth for the whole subgraph, not one container in isolation

What is not enough:

- healthchecks
- restart policy
- future portability language

### Queue and vector-store subsystems

What matters:

- durable message or index truth
- recovery ordering
- writer and reader reconvergence
- client behavior after topology change

What is not enough:

- the management UI loading
- a metrics target staying up
- one process starting again

## The proof packet each stateful class would need before stronger language

The repo should expect at least this packet before graduating a stateful
sentence:

1. `Authority before failure:` who was the writer or leader?
2. `Failure introduced:` what exactly died or was removed?
3. `Authority after failure:` who became authoritative, by which mechanism?
4. `Fencing:` what prevented stale authority from continuing to act?
5. `Client rediscovery:` how did dependents find the new truth?
6. `Storage truth:` what durable state replicated, moved, or remained singular?
7. `Remaining forbidden sentence:` what stronger claim is still not earned?

Without that packet, most stateful optimism in this repo should still be
treated as architectural desire rather than proof.

## What the repo is allowed to say honestly today

The repo is allowed to say:

- stateful dependencies are real and already central
- stateful SPOF pressure is not theoretical
- the docs now separate ingress continuity from authority continuity
- some services deserve explicit single-writer honesty rather than softened
  resilience language

The repo is not yet allowed to say:

- stateful zero-SPOF is solved
- routed TCP equals HA
- the current stack generally survives authority-node loss without private
  operator completion
- the data plane has become as mature as the ingress story

## What a real stateful progress packet would need

Before this page can support stronger language, a proof packet should show:

- the exact service class
- where writer authority lived before failure
- what artifact or system decided the replacement authority
- how the old authority was fenced or invalidated
- how clients rediscovered the new truth
- what durable storage truth moved, replicated, or remained singular
- what sentence is still forbidden after the test

Without that packet, the docs may describe better future ideas while the
runtime still leaves authority socially singular.

## Bottom line

The user does not want stateful services to sound more distributed.
The user wants a stricter answer to one brutal question:

> if this node dies, does the system still know who owns truth, and can it act
> on that truth without me being the missing algorithm?

The current worktree answers that question honestly in only one broad way:

not yet.

What it does prove is where the pain already lives:

- MongoDB
- Redis
- Headscale SQLite
- Firecrawl Postgres and RabbitMQ
- LLM Postgres
- Qdrant
- node-local bind-mounted durability assumptions

That is enough to make stateful honesty mandatory.
It is not yet enough to make stateful HA claims adult.

## Bottom line

The user is not demanding that every database magically become enterprise-grade
today.
The user is demanding that the docs stop pretending stateful dignity exists
where only transport dignity exists.

The strongest current honest sentence is still:

> the repo has real stateful dependencies, real authority pressure, and real
> node-local truth, but it does not yet broadly prove that authority survives
> node loss without private operator completion.

That sentence is harsher than normal infra prose.
It is also much closer to what the user is actually trying to learn.
