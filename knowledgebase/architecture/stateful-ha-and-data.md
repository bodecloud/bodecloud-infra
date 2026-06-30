# Stateful HA and Data Ownership

For the evidence boundary behind this page, start with
[Stateful HA Evidence](../research/stateful-ha-evidence.md).

This page exists because the user's real stateful question is harsher than:

> does the service still answer?

It is:

> if the current node dies, where does authoritative truth live, who is still
> allowed to write, how do clients rediscover that truth, and which parts of
> the answer are still secretly being carried by one operator?

The neighboring smaller question this page must not collapse into is:

> which databases do we run?

That smaller question is not enough.
The user is not just inventorying stateful services.
They are checking whether "multi-node" is still mostly edge theater once data
authority matters.

## What this page is and is not allowed to prove

This page is allowed to prove:

- the live priority runtime already depends on real stateful services
- those services still lean heavily on node-local authority and storage
- the repo already distinguishes ingress progress from stateful correctness
- the current worktree still falls far short of generic stateful anti-SPOF

This page is not allowed to prove:

- that TCP exposure equals HA
- that restartability equals authority preservation
- that one more node or one more bind mount equals failover dignity
- that future plans make the present runtime safe

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
- demo smoothness
- recovery options

They do not, by themselves, prove preserved authority.

## The strongest honest current answer

The strongest honest answer from the current worktree is:

- stateful risk is already live and significant
- authority is still mostly singular per service
- durability is still heavily node-local
- planning pressure is much sharper than live proof

That means the repo already knows many of the right future questions.
It does not yet prove that the main stateful classes have crossed into
trustworthy multi-node behavior.

## What the live runtime already proves

The current priority runtime already depends on all of the following stateful
classes:

- MongoDB in the root runtime
- Redis in the root runtime
- Headscale with SQLite in the active Headscale fragment
- Postgres and RabbitMQ in the active Firecrawl fragment
- Postgres in the active LLM fragment
- multiple bind-mounted volume paths under `${CONFIG_PATH:-./volumes}/...`

That matters because stateful pain is not speculative future architecture pain.
It is already present in the live stack.

## Stateful inventory by class

### 1. MongoDB in the root runtime

Live evidence:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  declares `mongodb`
- persistent data is bind-mounted at
  `${CONFIG_PATH:-./volumes}/mongodb/data:/data/db`
- the service is exposed through Traefik TCP labels

What this proves:

- MongoDB is a real dependency of the priority implementation
- durability is still tied to a node-local path
- exposure and healthchecks already exist

What it does not prove:

- replica-set semantics
- election behavior
- write authority promotion
- client topology rediscovery after failure

### 2. Redis in the root runtime

Live evidence:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  declares `redis`
- persistent data is bind-mounted at
  `${CONFIG_PATH:-./volumes}/redis:/data`
- the service is exposed through Traefik TCP labels

What this proves:

- Redis is a live stateful dependency, not just a helper idea
- the root runtime still looks like a single durable Redis instance story

What it does not prove:

- Sentinel
- replica promotion
- correct client reconnect behavior
- resilient master discovery

### 3. Headscale in the active fragment

Live evidence:

- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
  sets `database.type: sqlite`
- the database path is `/var/lib/headscale/db.sqlite`
- WAL mode is enabled
- the config itself notes that Postgres is discouraged upstream for current
  Headscale development

What this proves:

- the private mesh control plane is currently anchored to a local SQLite file
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

- the stateful surface is wider than just the obvious root databases
- application continuity already depends on a small graph of durable backends

What it does not prove:

- HA Postgres
- mirrored durable messaging semantics
- promotion-safe topology across the subgraph
- resilient multi-node truth for the whole application cluster

### 5. LLM Postgres surface

Live evidence from
[`compose/docker-compose.llm.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.llm.yml):

- `litellm-postgres` is live
- its data persists under
  `${CONFIG_PATH:-./volumes}/litellm/pgdata:/var/lib/postgresql/data`

What this proves:

- more than one Postgres-backed application surface already exists
- node-local persistence assumptions are repeating, not isolated

What it does not prove:

- HA Postgres for LLM-facing services
- promotion semantics
- reconnect correctness after node loss

## What the live runtime still shows clearly

The live worktree still shows:

- persistence is mostly node-local
- writer authority is mostly singular
- replicated authority is not the default reality
- client rediscovery behavior is not yet a generally proven surface

That is the shortest honest reason stateful claims must stay much harsher than
ingress claims.

The edge can become more survivable earlier than the data can become more
truth-owning.

## Why ingress progress and stateful progress diverge

This repo can improve user-visible continuity before it earns honest stateful
HA.

That asymmetry is useful because:

- wrong-node HTTP recovery can remove real operator pain early
- some services are mostly request-routing problems
- the edge stack already carries real policy, auth, and observability meaning

That asymmetry is dangerous because:

- the hostname can still work while the only real copy of data lived on the
  dead node
- a proxy path can preserve appearance while authority remains singular
- a convincing demo can hide the fact that one disk and one writer still
  mattered more than the rest of the topology

That is why "the hostname still works" is almost meaningless by itself for this
layer.

## The three questions every stateful claim must answer

Every serious stateful claim in this repo should answer all three:

1. Where does authoritative truth live right now?
2. What happens to that authority if the current node disappears?
3. How do clients discover the new truth without human folklore filling the
   gap?

If any of those questions is still mainly answered by operator memory, then the
system is still socially carrying part of the control plane.

## What still does not count as stateful progress

These are still weak or fake progress signals here:

- the same hostname still resolves
- the TCP port still answers
- the container restarts on the same node
- the service could be redeployed elsewhere later
- a second machine could theoretically mount similar storage
- a liveness check still passes

Those may all improve confidence.
They do not yet prove preserved authority.

## Service-class reality check

### Redis-class systems

What matters:

- primary truth
- promotion rules
- client reconnect behavior
- storage durability

What is not enough:

- one more routed endpoint
- one more container on another node

### Mongo-class systems

What matters:

- replica-set semantics
- election behavior
- consistency model
- client topology awareness

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
- promotion/fencing semantics
- storage portability
- reconnect behavior

What is not enough:

- healthchecks
- restart policy
- future portability language

## What the repo is allowed to say honestly today

The repo is allowed to say:

- stateful dependencies are real and already central
- stateful SPOF pressure is not theoretical
- the docs now separate ingress continuity from authority continuity
- some services may still deserve explicit single-writer honesty rather than
  fake resilience language

The repo is not yet allowed to say:

- stateful zero-SPOF is solved
- routed TCP equals HA
- the current stack generally survives authority-node loss without private
  operator completion

## Bottom line

The user does not want stateful services to sound more distributed.
The user wants a stricter answer to one brutal question:

> if this node dies, does the system still know who owns truth, and can it act
> on that truth without me being the missing algorithm?

The current worktree answers that question honestly in only one broad way:

not yet.

What it *does* prove is where the pain already lives:

- MongoDB
- Redis
- Headscale SQLite
- Firecrawl Postgres and RabbitMQ
- LLM Postgres
- node-local bind-mounted durability assumptions

That is enough to make stateful honesty mandatory.
It is not yet enough to make stateful HA claims adult.
