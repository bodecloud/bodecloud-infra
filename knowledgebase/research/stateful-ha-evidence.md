# Stateful HA Evidence

This page is the evidence boundary for the repo's stateful story.

It exists because this is the easiest layer for a multi-node stack to look
grown up while still secretly depending on one sacred writer, one sacred disk,
or one sacred operator memory.

The user is not asking whether a database process can come back.
The user is asking whether the system can stop privately depending on:

- one node that still owns the only real write authority
- one local path that still contains the only truthful durable copy
- one recovery ritual that still only works because the operator remembers the
  topology better than the system does

If those sentences are still true, then the stack may be reachable, restartable,
or routable, but it is not yet what the user is trying to build.

## What this page is and is not allowed to prove

This page is allowed to prove:

1. the priority runtime already contains several real stateful dependencies
2. those dependencies are still mostly expressed as node-local storage plus
   node-local authority
3. the repo already knows ingress continuity and stateful correctness are
   different problems
4. some future hardening ideas exist for several stateful classes
5. the current worktree still leaves most state authority concentrated in
   singular places

This page is not allowed to prove:

- that stateful anti-SPOF is solved
- that a TCP endpoint surviving equals preserved authority
- that a route through Traefik or HAProxy implies correct writer failover
- that a service being important and persistent makes it resilient
- that a planning document can upgrade present tense runtime truth

## The dream this page has to protect

The dream is not "stateful services exist on more than one box somewhere."

The dream is:

- a node can disappear
- a truthful writer still exists
- clients can find that writer without human folklore filling the gap
- the old writer cannot keep acting like it still owns truth
- storage survival is not secretly just one lucky disk path

Anything weaker can still feel impressive from the outside.
That is why stateful evidence has to be harsher than ingress evidence.

At ingress, fake maturity often sounds like:

- the hostname still resolved
- the reverse proxy still answered
- the route graph still looked serious

At the stateful layer, fake maturity sounds like:

- the database port still answered
- the application reconnected once
- another node could theoretically run the same service
- the docs know the right cluster nouns

None of those are enough.

## Strongest honest current answer

The strongest honest current answer from the worktree is:

- the repo already depends on real stateful systems in the live priority stack
- those systems are still mostly anchored to node-local persistence and
  node-local authority
- the planning layer understands many of the missing truths more clearly than
  the runtime currently proves
- the current stack is better at proving stateful exposure and stateful risk
  awareness than replicated authority, safe promotion, or client rediscovery

That is meaningful progress in honesty.
It is not yet stateful dignity.

## Evidence hierarchy for stateful claims

| Claim type | Highest authority | Why it outranks others | It still does not prove |
| --- | --- | --- | --- |
| What stateful services are live in the priority runtime | [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml) and active `compose/docker-compose.*.yml` fragments | these files define the tracked runtime surface | resilient authority behavior |
| What the repo already knows stateful HA would require | [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md), [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md) | these are the clearest planning and honesty surfaces | that the topology is live now |
| Why fake HA language is unacceptable here | the knowledgebase architecture and research pages | they preserve the burden standard the user is actually asking for | low-level runtime mechanics by themselves |

If a paragraph sounds stronger than this table permits, it is probably
borrowing confidence from a different layer.

## What the live runtime concretely proves today

The priority implementation still centers on:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
- [`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- [`compose/docker-compose.llm.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.llm.yml)

### 1. MongoDB is live and still a node-local authority story

Verified runtime evidence:

- root runtime declares `mongodb`
- persistent storage is bind-mounted at
  `${CONFIG_PATH:-./volumes}/mongodb/data:/data/db`
- it is exposed through Traefik TCP labels

This proves:

- MongoDB is a real dependency of the priority runtime
- durability is still tied to a local filesystem path
- transport exposure already exists

This does not prove:

- Replica Set membership
- election behavior
- safe primary promotion
- client rediscovery of the new primary
- fencing or split-brain prevention

The hidden sentence still surviving is:

> I still privately know which node and disk actually matter for MongoDB truth.

### 2. Redis is live and still mostly a single-writer durability story

Verified runtime evidence:

- root runtime declares `redis`
- persistent storage is bind-mounted at `${CONFIG_PATH:-./volumes}/redis:/data`
- it is exposed through Traefik TCP labels

This proves:

- Redis is a live stateful dependency
- the stack is not pretending Redis is merely optional cache fluff
- the runtime still reads like one durable Redis instance, not an already
  verified failover topology

This does not prove:

- Sentinel
- replica promotion
- client master rediscovery
- reconnect correctness after failover
- any shared notion of write authority transfer

The hidden sentence still surviving is:

> I still privately know Redis truth is basically where the current writer is.

Archive-derived boundary:

- `redis-url-and-load-balancing__68a914f8-d47c-8324-8734-bc1f17507bac.md`
  separates three things that infrastructure docs often blur:
  - `redis://` is raw Redis protocol over TCP, not HTTP
  - Traefik can route Redis only through TCP routers, not HTTP routers
  - Traefik's Docker provider on ordinary non-Swarm Docker sees only the local
    Docker daemon

That source matters because it is not merely saying "Redis is stateful."
It is identifying a practical false shortcut:

> if every node has the same Docker labels, Traefik will discover and balance
> Redis globally across non-Swarm hosts.

That shortcut is false.
Identical labels across independent Docker daemons do not create global backend
truth.
They create repeated local declarations unless some other registry, file
provider, Consul/etcd surface, generated config, or orchestrator-grade provider
supplies cross-node knowledge.

Even if cross-node TCP forwarding is later implemented, the stateful claim
still stays narrow.
TCP forwarding can prove transport reachability.
It cannot prove:

- which Redis instance is the current writer
- whether a replica was promoted
- whether clients rediscovered the promoted master
- whether independent Redis servers accidentally diverged
- whether stale writers were fenced

The legal sentence is:

> Redis can be exposed as a TCP service, and a future registry-backed forwarding
> layer could route Redis connections to a selected peer.

The illegal sentence is:

> Traefik TCP labels across ordinary non-Swarm Docker hosts make Redis
> multi-node HA.

### 3. Qdrant is live and persistent, but persistence is not authority transfer

Verified runtime evidence:

- `compose/docker-compose.core.yml` declares `qdrant`
- storage is bind-mounted at
  `${CONFIG_PATH:-./volumes}/qdrant/storage:/qdrant/storage`
- the service is exposed over HTTP through Traefik

This proves:

- vector state is already part of the live stack
- durable storage is still expressed as a local path
- reachability is already better than authority semantics

This does not prove:

- cluster membership
- quorum or replica behavior
- write-order correctness under node failure
- recovery semantics that would let a newcomer identify the authoritative copy

Qdrant matters here because vector stores are especially easy to narrate as
"already distributed enough" while still behaving like local durable state.

### 4. Firecrawl expands the stateful surface into a small stateful graph

Verified runtime evidence:

- `compose/docker-compose.firecrawl.yml` depends on `redis`, `nuq-postgres`,
  and `rabbitmq`
- `nuq-postgres` stores data under
  `${CONFIG_PATH:-./volumes}/nuq-postgres/data:/var/lib/postgresql/data`

This proves:

- the stateful surface is wider than MongoDB and Redis
- continuity for this subsystem already depends on more than one durable
  backend
- messaging durability and relational durability are already part of the live
  problem

This does not prove:

- HA Postgres
- safe queue continuity under failover
- correct recovery ordering across Redis, Postgres, and RabbitMQ
- promotion discipline across that subgraph

This is important because "the app came back" is especially weak evidence when
the app's truth depends on several stateful subsystems whose authority contracts
have not been proven together.

### 5. Headscale is currently a singleton control-plane authority

Verified runtime evidence:

- `compose/docker-compose.headscale.yml` configures `database.type: sqlite`
- the database path is `/var/lib/headscale/db.sqlite`
- WAL mode is enabled

Planning evidence:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
  explicitly treats Headscale as single-node today
- the same planning surface names leader-election plus replication as future
  work, not present proof

This proves:

- the mesh control plane is currently anchored to one local SQLite authority
- the repo is at least honest enough to admit that this remains singular

This does not prove:

- automatic control-plane replacement
- shared authority across nodes
- cluster-safe promotion
- mesh continuity after losing the Headscale authority node

Headscale is one of the cleanest examples in the repo of the difference between:

- a live service
- a valuable service
- a resilient authority surface

The first two are already true.
The third is still aspirational.

### 6. The LLM stack adds another relational authority plus vector state

Verified runtime evidence:

- `compose/docker-compose.llm.yml` declares `litellm-postgres`
- the same fragment declares `qdrant`
- application components depend on `litellm-postgres` and `redis`

This proves:

- the LLM side of the stack is also already state-bearing
- more than one subsystem already depends on Postgres-class truth
- vector and relational state are both live concerns, not future ambitions

This does not prove:

- a failover-ready Postgres topology
- client failover behavior for those apps
- authority transfer for vector state
- that the LLM layer has stopped inheriting node-local storage assumptions

## The cross-cutting burdens the runtime still leaves unresolved

Across all of those service classes, the same deeper burdens remain:

| Burden | What the user is really asking | What fake progress looks like |
| --- | --- | --- |
| Writer identity | who is allowed to accept writes right now? | a port answered |
| Authority continuity | who owns truth after failure? | the service restarted |
| Client rediscovery | how do dependents find the new truth? | DNS still landed somewhere |
| Fencing | what prevents the old writer from still acting alive? | the old node is probably gone |
| Storage truth | did authority replicate, move, or stay singular? | more than one node has a compose file |

If these questions are still mostly answered by operator memory, then part of
the control plane is still social instead of explicit.

## What still does not count as stateful evidence here

The following still do not count as trustworthy stateful proof:

- an endpoint answering over TCP or HTTP
- a bind-mounted volume existing at a durable-looking path
- an app reconnecting once after restart
- a second node being theoretically capable of running the same service
- a future HA note sitting beside a singleton live deployment
- a proxy path staying alive while writer truth stays singular

This page should be hostile to paragraphs that become calmer by using words
like:

- durable
- redundant
- resilient
- highly available

without forcing the reader to answer writer identity, promotion rules, client
behavior, and storage truth.

## What a real proof packet would need before stronger language is allowed

Before any workload in this repo can be described as honestly resilient, the
proof packet should identify at least:

1. the authoritative writer or leader model before failure
2. the exact failure introduced on purpose
3. the promotion or election mechanism that selected the new writer
4. the evidence that the old writer stopped being authoritative
5. what clients observed before, during, and after failover
6. what storage substrate preserved the authoritative dataset
7. what contradiction still remains, if any

Without that packet, the repo is only allowed to say:

- the workload is important
- the workload is persistent
- the workload has a candidate topology
- the workload has a hardening plan

It is not yet allowed to say:

- the workload is honestly HA
- the workload is anti-SPOF
- the workload survives node loss with preserved authority

### Packet schema for this repo

Use this exact shape when collecting evidence for a stateful claim:

```yaml
stateful_authority_packet:
  claim_tested: "stateful authority under failure"
  service: "redis | mongodb | headscale | postgres | rabbitmq | qdrant"
  authority_before: "<writer/leader/source of truth before failure>"
  failure_introduced: "<exact node, process, disk, network, or backend failure>"
  authority_after: "<writer/leader/source of truth after failure>"
  client_observation: "<what dependent clients saw before/during/after>"
  rediscovery_mechanism: "<DNS, seed list, Sentinel, driver, registry, manual, none>"
  fencing_or_split_brain_guard: "<mechanism, or none>"
  storage_truth: "<replication, backup, snapshot, shared storage, singular disk>"
  operator_intervention_required: true
  result: "pass | fail | honest-singularity | inconclusive"
  what_this_proves: "<one narrow sentence>"
  what_is_still_forbidden: "<larger HA sentence still illegal>"
```

The most useful result value in the current worktree may often be
`honest-singularity`.
That result is not a failure of the documentation.
It is the documentation refusing to promote a singleton writer into a fake
cluster.

### Current-runtime packet ceilings

These are not completed proof packets.
They are the highest honest packet ceilings visible from the current runtime.

| Service | Current ceiling | Why it stops there |
| --- | --- | --- |
| MongoDB | reachable, persistent, TCP-exposed singleton authority | no tracked proof of Replica Set election, client seed-list/SRV rediscovery, or fencing |
| Redis | reachable, persistent, TCP-exposed singleton writer | no tracked proof of Sentinel, promoted master discovery, or reconnect correctness |
| Headscale | useful singleton control-plane authority | SQLite-backed local authority is explicit; leader election and replication remain future work |
| Firecrawl Postgres/RabbitMQ/Redis graph | real multi-backend stateful dependency graph | no tracked proof that the graph recovers in a correct order under partial backend loss |
| LiteLLM Postgres/Redis/Qdrant surfaces | real relational, cache, and vector state pressure | no tracked proof of failover-ready Postgres, cache authority transfer, or vector cluster semantics |
| Qdrant | reachable persistent vector store | no tracked proof of clustered replica semantics or authoritative index recovery |

These ceilings align with the instruction surfaces:
Compose is the priority implementation, and the project is trying to become
peer-aware without pretending that exposed state equals stateful HA.

## Bottom line

The worktree already proves that state matters now.

It does not yet prove that the main stateful classes have escaped singular
authority, singular storage truth, or socially reconstructed failover.

That is the correct harsh answer.

It is less comforting than a distributed-looking diagram.
It is also much closer to what the user is actually trying to learn.
