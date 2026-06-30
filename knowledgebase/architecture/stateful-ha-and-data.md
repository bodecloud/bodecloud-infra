# Stateful HA and Data Ownership

For the evidence boundary behind this page, start with
[`../research/stateful-ha-evidence.md`](../research/stateful-ha-evidence.md).

This page exists because the user is not merely asking whether hostnames still
answer.
They are asking something harsher:

> if one node dies, does the system still know where authoritative state lives,
> who may write, how clients rediscover the valid topology, and whether the
> surviving path is actually trustworthy?

That is where most infrastructure writing starts cheating.
It celebrates:

- a still-routed hostname
- a reachable standby
- a restarted container
- a proxy that still accepts connections

while skipping the only question that matters:

> did the system preserve authority, or only the appearance of continuity?

That distinction is why this page has to stay blunt.

## What this page is and is not allowed to prove

This page is allowed to:

- define the boundary between ingress progress and real stateful correctness
- explain why stateful HA is much stricter than "reachable from more places"
- identify which capabilities still block honest stateful claims
- keep the rest of the docs from over-crediting routing progress

This page is not allowed to:

- imply that TCP proxying equals stateful HA
- treat replicated access paths as replicated authority
- upgrade a stable hostname into write-path correctness
- pretend the repo has already solved leader truth, storage portability, or
  failover reconciliation unless the current worktree proves it

## The hard rule this repo is defending

For any state-bearing service, you do not honestly remove the SPOF merely by:

- exposing it through Traefik
- giving it a stable hostname
- copying the service definition to another node
- teaching one node to TCP-forward to another node
- showing that another box can answer on the same port

Those things may improve:

- reachability
- operator ergonomics
- restartability
- demo quality

They do not, by themselves, produce stateful HA.

Real stateful HA requires some combination of:

- replicated data across failure domains
- explicit authority or election semantics
- client rediscovery or reconnect behavior that survives promotion
- storage semantics that do not strand truth on one dead node
- recovery behavior that does not collapse back into improvised operator rescue

That last line matters especially here.
The user is exhausted with systems whose real failover procedure is:

1. remember which node used to be primary
2. decide which survivor should win now
3. remember which clients must move
4. remember which hostname is suddenly lying

That is not HA.
That is human failover wearing system language.

## What the live runtime already proves

The current priority runtime already proves real stateful risk exists.

From the root runtime and active fragments, the repo already depends on:

- MongoDB
- Redis
- Headscale with SQLite in the active headscale fragment
- Postgres in additional active subsystems such as Firecrawl and other service
  fragments
- multiple node-local bind-mount durability assumptions

That means stateful pain is not hypothetical future architecture pain.
It is already part of the live stack.

## What the live runtime still shows clearly

The strongest honest reading of the current runtime is:

- persistence is still heavily node-local
- authority is still mostly singular per service
- replication is not yet the default reality
- client topology correctness is not yet a generally proved surface

That is visible in concrete ways:

- MongoDB still binds data locally
- Redis still behaves like a single durable instance, not a proved Sentinel
  or promoted topology
- Headscale's active config still uses SQLite and explicitly notes Postgres is
  discouraged upstream for current Headscale development
- other services depend on Postgres or RabbitMQ but are not yet documented as
  replicated authority models

That is why this page has to refuse flattering summaries.

## Why ingress progress and stateful progress diverge so sharply

This repo can honestly improve ingress before it earns honest stateful HA.

That asymmetry is useful and dangerous.

Useful because:

- better ingress can improve user-visible survivability early
- wrong-node HTTP recovery can remove real operator pain before stateful HA is
  solved
- some services really are closer to stateless routing problems than to
  authority problems

Dangerous because:

- a service can still look "up" while its only real copy of data lives on the
  dead node
- peer forwarding can preserve the appearance of continuity while truth remains
  singular
- a convincing demo can hide the most important unresolved failure

That is why "the hostname still works" can be an actively misleading sentence
for this layer.

For stateless HTTP it may be real progress.
For stateful systems it can mean the edge stayed pretty while correctness did
not survive.

## The three questions every stateful claim must answer

Every serious stateful claim in this repo should answer all three:

1. Where does authoritative truth live right now?
2. What happens to that authority if the current node disappears?
3. How do clients discover the new truth without human folklore filling the
   gap?

If any of those questions is still mainly answered by operator memory, then the
system is still socially carrying part of the control plane.

That is the phrase to keep in mind:

social control plane.

It explains why many "HA-looking" answers still fail the user's standard.

## Service-class reality check

Stateful services should be read by class, not with one generic resilience
label.

### Redis-class systems

What matters:

- master or primary truth
- promotion semantics
- client reconnect behavior
- storage durability

What is not enough:

- one more routed endpoint
- one more container on another node

### Mongo-class systems

What matters:

- replica-set semantics
- election behavior
- consistency and client topology awareness
- data locality across failure domains

What is not enough:

- TCP exposure through Traefik
- a local healthcheck

### SQLite-backed control-plane services such as current Headscale

What matters:

- whether the database is still fundamentally local to one node
- how leadership or singleton truth is represented
- whether a replacement node can become authoritative without manual rescue

What is not enough:

- the service being reachable through the edge
- the admin surface responding

### Postgres-backed application subsystems

What matters:

- fencing and promotion semantics
- storage portability
- reconnect behavior
- which subsystem actually needs HA versus explicit single-writer honesty

What is not enough:

- healthchecks
- a restart policy
- "could be moved elsewhere later"

## What the user actually wants from the stateful story

The user is not asking for stateful services to sound more distributed.
They are asking for a stricter answer:

> if this node dies, does the system still know who owns truth, and can it act
> on that truth without me being the missing algorithm?

That is the stateful benchmark.

It is why the docs should sometimes prefer explicit single-writer honesty over
fake resilience language.

A clearly described single-node authority model is more honest than a routed
multi-node illusion that still depends on remembered rescue steps.

## What would count as real progress here

Real progress in this layer would look like:

1. service-class-specific authority models, not generic "stateful HA" slogans
2. explicit replication or single-writer honesty per critical service
3. storage and promotion semantics that survive node loss
4. client rediscovery behavior that does not depend on operator folklore
5. drills that prove correctness, not just liveness

Until then, the repo should keep saying:

- ingress may improve earlier
- stateful truth remains harsher

That unevenness is not a docs flaw here.
Pretending it has already been smoothed away would be the flaw.

## Strongest honest current answer

The live stack already carries meaningful stateful risk across multiple service
domains, and the repo is unusually honest about the difference between routed
reachability and replicated authority.

What it does not yet prove is the thing that matters most:

that truth itself has stopped living in one brittle place, with the operator
acting as the missing failover logic when that place disappears.
