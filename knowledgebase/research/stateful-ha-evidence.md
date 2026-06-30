# Stateful HA Evidence

This page is the proof boundary for the repo's stateful story.

It exists because "multi-node" becomes least trustworthy precisely where state
starts to matter most.

Several nodes, a public proxy, and even a plausible fallback path can still
hide a very old problem:

- one write authority
- one node-local disk path
- one unreplicated SQLite file
- one reconnect story that only works in the flattering order

That is not the kind of ambiguity this repo can afford.

At the ingress layer, weak answers can still look impressive.
At the stateful layer, weak answers can preserve the interface while losing the
place where correctness actually lived.

This page therefore exists to stop the docs from calling reachability "HA."

## The dream this page has to protect

The user is not asking whether stateful services can be restarted.

The user is asking whether the system can stop depending on sacred authority
points:

- one host that owns the writable truth
- one disk path that cannot disappear safely
- one topology that only behaves correctly if the operator already knows where
  truth lives

That is why stateful HA in this repo has to stay separate from:

- public entry
- reverse-proxy success
- basic healthchecks
- "container came back up"

The dream is not distributed-looking state.
The dream is state that has stopped secretly living in one brittle place.

## Strongest honest current answer

The priority runtime already depends on real stateful services.

Those services are still largely anchored to node-local authority and node-local
storage.

The planning layer is sharper about the required stateful futures than the live
runtime currently proves.

The repo already knows that ingress continuity and stateful correctness are not
the same thing.

It does **not** yet prove that the main stateful classes have crossed into
trustworthy anti-SPOF topologies.

That means the current stack is stronger at proving:

- stateful exposure
- stateful importance
- stateful risk awareness

than it is at proving:

- replicated authority
- promotion safety
- reconnect correctness
- storage portability
- fencing or quorum discipline

## What this page is and is not allowed to prove

This page is allowed to prove:

1. the live root runtime already depends on real stateful systems
2. several of those systems still reveal node-local durability assumptions
3. Headscale is currently a singleton control-plane concern
4. the planning layer already separates ingress failover from stateful
   correctness
5. the repo has concrete future ideas for some stateful classes
6. the current worktree still leaves state authority concentrated in specific
   places
7. fake HA language is especially dangerous here

This page is not allowed to prove:

- that stateful HA is solved
- that routed TCP endpoints equal resilient authority
- that a service being reachable through Traefik means the backing state is
  resilient
- that live node-local persistence is acceptable just because future HA plans
  exist

## What still does not count as stateful evidence here

The following still do not count as trustworthy stateful proof:

- an endpoint answering over TCP
- a durable-looking bind mount path existing
- an application reconnecting once after a restart
- a future HA note sitting beside a singleton live deployment
- a proxy or ingress path staying up while authority stays singular
- a service looking distributed from the outside while write truth stays local

This page exists precisely because those weaker signals are easy to narrate as
progress when the hard authority question is still unresolved.

## Evidence hierarchy for stateful claims

| Claim type | Highest authority | Why it outranks others | It still does not prove |
| --- | --- | --- | --- |
| What stateful systems are really in the priority runtime | `docker-compose.yml`, active `compose/docker-compose.*.yml` | this is the live implementation surface | resilient topology behavior |
| What the repo already knows stateful HA requires | `docs/stateful_ha_plan.md`, `docs/INFRASTRUCTURE_MASTER_PLAN.md`, `README.md` | these are the clearest honesty and promotion surfaces | that the future topology is live |
| Why the user rejects weak HA language | archive-pressure and source archive pages | they restore the true standard | technical completion |

If a stateful paragraph sounds stronger than this table allows, it is probably
borrowing confidence from ingress.

## What the live runtime concretely proves

The priority implementation still centers on:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
- [`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)

### 1. MongoDB is live and still node-local

[`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
shows `mongodb` with:

- a bind mount under `${CONFIG_PATH:-./volumes}/mongodb/data`
- TCP exposure
- local healthchecks

This proves:

- MongoDB is not a hypothetical future consideration
- durable truth currently still ties back to one local storage path per node

It does **not** prove:

- Replica Set semantics
- promotion safety
- write authority failover
- client topology correctness after failure

### 2. Redis is live and still mostly a single-instance durability story

The same core compose surface shows `redis` with:

- local `/data` persistence under `${CONFIG_PATH:-./volumes}/redis`
- password-backed single-instance assumptions
- TCP exposure
- local healthchecks

This proves:

- Redis is a real stateful dependency of the root stack
- it is not currently documented as if it already has Sentinel-backed
  authority failover

It does **not** prove:

- Sentinel
- replica promotion
- application reconnect correctness
- correct master discovery after a node event

### 3. Firecrawl extends the stateful surface beyond the obvious core

[`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)
adds:

- `nuq-postgres`
- `rabbitmq`
- application dependency on `redis`, `nuq-postgres`, and `rabbitmq`

`nuq-postgres` itself uses a bind mount under
`${CONFIG_PATH:-./volumes}/nuq-postgres/data`.

This proves:

- the stateful surface is wider than just MongoDB and Redis
- application continuity already depends on several durable backends

It does **not** prove:

- HA Postgres
- mirrored durable storage
- failover-safe messaging semantics
- resilient promotion across this subgraph

### 4. Headscale is currently a singleton control-plane problem

[`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
is explicit:

- database type is `sqlite`
- path is `/var/lib/headscale/db.sqlite`
- WAL is enabled

[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
is equally explicit:

- Headscale is single-node today
- if the node running it goes down, the mesh loses its control plane
- a future HA direction would require leader-election discipline plus
  Litestream-backed replication for the SQLite database

This is one of the clearest examples in the repo of the difference between:

- a live service
- a useful service
- a resilient authority surface

Headscale clearly satisfies the first two.
The docs themselves say the third is still missing.

## What the live runtime does not yet prove

The current worktree does not, by itself, prove:

- Redis Sentinel or equivalent master election
- MongoDB Replica Set health and client failover behavior
- failover-safe Postgres promotion
- RabbitMQ cluster semantics
- fencing against split-brain or double-writer conditions
- state portability independent of one node-local bind mount
- operator-independent discovery of who currently owns write truth

That is why stateful HA needs its own proof standard here.

Without that separation, the stack could look better than it actually is
because:

- a route still resolves
- the container restarts
- or the endpoint answers

None of those things tells you whether authority stopped being singular.

## What a real stateful evidence packet would need

For this repo, a serious stateful packet would need to show:

- who the writer or authoritative leader was before failure
- what failure was introduced intentionally
- how promotion, election, or recovery occurred
- what storage substrate preserved the authoritative dataset
- what clients observed during and after the event
- what evidence proves continuity of authority rather than mere reachability

Without that packet, the docs are only allowed to say:

- "stateful dependency exists"
- "stateful risk is acknowledged"
- or "stateful hardening is planned"

not "stateful HA is demonstrated."

## What the planning layer already knows

The planning docs are much sharper than generic "someday we should cluster
things" language.

### 1. `stateful_ha_plan.md` already rejects container relocation as a fake answer

[`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
states that zero-SPOF stateful services do not come from simply moving
containers around.

It explicitly separates:

- HTTP hostname routing
- plain TCP forwarding
- replication and quorum requirements

That is important because it keeps the repo from pretending that the ingress
story automatically solves state.

### 2. `INFRASTRUCTURE_MASTER_PLAN.md` already names the concrete stateful gaps

The master plan calls out several real pressure points:

- Headscale singleton control plane
- SQLite at `/var/lib/headscale/db.sqlite`
- Litestream as a possible replication layer for Headscale's SQLite file
- missing automated service failover between nodes
- route failover and state failover being different modules

The plan does not just say "HA is missing."
It says where authority still lives and what classes of repair might earn trust
later.

### 3. The repo is already stricter than most ecosystems about state language

This is one of the strongest positive signals in the project.

The docs and planning surfaces repeatedly refuse to equate:

- node survival
- endpoint reachability
- container restart
- proxy continuity

with:

- durable truth survival

That honesty is valuable.
It still is not the same thing as completion.

## The real blocker set for stateful anti-SPOF

The hard problem is not "make state available on another node."
The hard problem is making all of these true together:

1. replicated authority exists
2. clients can discover the current authority safely
3. writes do not fork during partial failure
4. storage is portable or replicated enough to survive node loss
5. the promotion path is more trustworthy than operator memory
6. application reconnect behavior matches the promoted topology

If any one of those remains private operator knowledge, the operator is still
acting as part of the control plane.

That is exactly what this repo is trying to reduce.

## Service classes the docs should keep separate

Stateful services in this repo are not one undifferentiated category.

At minimum, keep these classes distinct:

- control-plane singleton state
  - current example: Headscale
- cache or coordination stores that may need promotion-aware clients
  - current example: Redis
- document or application databases with stronger durability expectations
  - current example: MongoDB
- application-specific relational stores
  - current example: `nuq-postgres`
- messaging systems whose correctness depends on topology semantics, not just
  open ports
  - current example: RabbitMQ

If the docs flatten these into "stateful services," they stop being useful.

## What this page should force every reader to ask

Before accepting any stateful HA sentence, ask:

1. where does write truth live right now?
2. what storage path holds it?
3. what exact mechanism replicates or protects it?
4. who decides promotion?
5. how do clients discover the new authority?
6. what prevents double-writer behavior?
7. what evidence proves this beyond route reachability?

If the page cannot help answer those questions, it has not really described
stateful HA.

It has only described stateful exposure.

## Bottom line

The current worktree proves that `bolabaden-infra` already depends on
stateful systems across several domains.

It also proves that several of those systems still depend on node-local storage
or singleton authority surfaces.

The planning layer is already honest about what real stateful anti-SPOF would
require.

The runtime does **not** yet prove that those requirements are satisfied.

That means the right summary is:

- stateful risk is real now
- stateful anti-SPOF is not yet honestly earned

This page should be read as a guardrail against one of the easiest lies in all
of infrastructure writing:

> the service stayed reachable, therefore the truth survived

In this repo, that sentence remains unproven unless authority itself stopped
being singular.
