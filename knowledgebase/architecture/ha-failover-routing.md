# HA, Failover, and Routing

Read this page as the routing truth map for the priority implementation rooted
at [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

If you want the deeper evidence stack first, read:

- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
- [Request Path and Failure Walkthrough](request-path-and-failure-walkthrough.md)
- [Current Compose Runtime](current-compose-runtime.md)
- [Operator Contract and Success Criteria](operator-contract.md)

This page exists because "HA" becomes almost useless in this repo unless
routing is broken into the separate truths the user is actually asking about.

The user is not mainly asking:

- can more than one node receive traffic?
- is Traefik present?
- do healthchecks exist?
- is there a helper that sounds failover-shaped?

The harder question is:

> if traffic lands on the wrong healthy machine, can that machine still
> preserve the service contract without requiring private operator memory?

That is the routing problem here.

## The routing dream, in the repo's own words

The strongest intent surface is
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It describes the target operating contract as:

- any healthy public node can receive the first request
- if the service is local, that node serves it locally
- if the service is remote, that node forwards to a healthy peer that currently
  hosts it
- the edge should preserve auth, middleware, and policy meaning rather than
  merely expose containers

That target is anti-SPOF and peer-aware without immediately collapsing into
Swarm, Kubernetes, k3s, or another heavyweight orchestrator.

It is also only intent unless the runtime proves it.

## What this page is and is not allowed to prove

This page is authoritative about:

- how routing has to be decomposed before "failover" means anything useful
- which routing classes exist in the current Compose runtime
- which truths are materially live, which are planned, and which are still
  operator-owned
- why first-hop plurality is weaker than request-preserving recovery
- why HTTP, TCP, and stateful failover must not be flattened together

This page is not authoritative about:

- whether a specific hostname has already passed a real wrong-node drill
- whether backend-loss fallback is already durable for a named route
- whether stateful traffic inherited optimism from the HTTP story
- whether one working helper path upgrades the whole platform

This is a routing decomposition page, not a victory lap.

## The routing classes that must stay separate

Routing in this repo splits into at least four different classes:

### 1. Stateless HTTP

This is the easiest class to talk about and the most dangerous one to
overclaim from.

What it would mean:

- wrong-node entry is acceptable
- the receiving node can discover the real backend
- forwarding keeps the route semantics intact
- backend loss can trigger a valid fallback rather than a dead route

What does **not** prove it:

- Cloudflare can reach more than one node
- Traefik is healthy
- a local route returns `200`
- a helper renders fallback-shaped config

### 2. Protected HTTP

This is stricter than ordinary stateless HTTP because the platform has to
preserve:

- auth behavior
- middleware order
- headers and trust boundaries
- the exact meaning of the protected route

Wrong-node success here means more than "the response still comes back."
It means the forwarded route is still the same protected service in user-visible
and policy-visible terms.

### 3. Raw TCP

This class is different enough that HTTP optimism cannot spill into it.

A TCP router for Redis or MongoDB does **not** prove:

- peer-aware service failover
- semantic continuity
- state correctness
- safe substitution between backends

It only proves transport is possible through a proxy surface.

### 4. Stateful surfaces

This is the strictest class.
State-bearing services do not become honestly HA because:

- they are reachable through Traefik
- multiple nodes exist
- a proxy can point at more than one host

Stateful routing only becomes meaningful when authority, data durability,
substitution safety, and recovery semantics are all explicit.

## Strongest honest current answer

The root runtime already contains enough edge and routing machinery to make the
problem real rather than hypothetical:

- Cloudflare-oriented public-entry assumptions
- Traefik
- TinyAuth
- `nginx-traefik-extensions`
- CrowdSec
- Docker socket proxies
- `docker-gen-failover`
- Headscale mesh assumptions
- TCP routers for services such as MongoDB and Redis

That is a real routing stack.

What it still does **not** clearly own is the thing the user actually wants
that machinery to buy:

- shared placement truth
- trustworthy peer eligibility truth
- durable fallback-route persistence under failure
- proof that a request keeps the same meaning after wrong-node handoff
- proof that stateful authority survives any routing story being told

That is the missing middle layer.

The repo already knows how to speak about routes.
It still does not universally know how to let the receiving node explain, in
shared system terms, why the route is safe when locality fails.

## What still does not count as HA or failover here

The following may be real progress.
They still do not count as meaningful HA or failover on their own in this repo:

- more than one public node can receive the first hop
- Traefik is present and healthy
- a helper generates fallback-shaped route material
- a TCP router exists for a stateful service
- a local protected route returns `200`
- a mesh exists between nodes
- the docs can now explain the route classes more clearly

All of those can improve posture.
None of them satisfy the user's benchmark unless they also reduce the need for
private placement memory, preserve request meaning on the wrong node, and
survive the failure that made fallback necessary.

The harsher checksum question is:

> yes, but who still had to know the real answer first?

If the answer is still "the operator," then the route may be real software and
still fake relief.

## The exact hidden routing jobs the operator may still be doing

Today the operator may still have to know things like:

- which node currently hosts the real copy of a named HTTP service
- whether a helper-generated rescue path would still exist after the preferred
  backend disappears
- whether a reachable peer is merely alive on the mesh or actually acceptable
  for the route's auth, middleware, secrets, and revision surface
- whether a TCP endpoint is merely transport-reachable or safe to describe with
  failover language
- whether a state-bearing surface still hides a sacred authority node behind a
  more distributed-looking ingress story

That list is why "there are a lot of routing options" is not the same thing as
relief in this repo.

## The missing truths that matter most

The routing story only becomes believable when the platform owns more of these:

### Placement truth

The receiving node needs shared, current knowledge of where the service
actually lives now.

The repo keeps converging on a lightweight current-state registry such as
`services.yaml`, but the docs must still treat that as intent unless the root
runtime really ships and consumes it.

### Peer eligibility truth

Reachability is not enough.
The platform needs a stricter answer to:

- which peer is healthy?
- which peer is serving the right revision?
- which peer is acceptable for this route's auth and middleware expectations?

### Route durability truth

A fallback route only matters if it survives the failure that made fallback
necessary.
A route rendered while everything is still healthy is not enough.

### Semantic continuity truth

The request has to keep meaning the same thing after forwarding.
That includes:

- auth continuity
- middleware continuity
- trust-boundary continuity
- user-visible route continuity

Without that, the route may be alive but not preserved.

## Why Cloudflare and Traefik are necessary but insufficient

Cloudflare is part of the anti-SPOF story as the first hop.
Traefik is central to the HTTP routing story.
Both are important.
Neither solves the whole problem alone.

Cloudflare can help more than one node become a first hop.
That is weaker than proving the wrong node can preserve the request correctly.

Traefik can express routers, middlewares, services, and TCP paths.
That is weaker than proving:

- the wrong node knows the right backend
- the peer is actually eligible
- the fallback survives the backend loss
- the route preserves meaning after handoff

This is exactly where normal HA language becomes misleading:
layer one gets solved loudly, and the remaining hidden operator work gets
quietly edited out.

## What a real routing proof packet would have to contain

If a future page supports stronger routing claims, it should point to a real
route-level proof packet containing:

- the exact route class:
  stateless HTTP, protected HTTP, raw TCP, or a named stateful surface
- the entry node and backend node identities
- the source of placement and peer truth used for the handoff
- the failure condition introduced, if fallback is being claimed
- the policy or middleware comparison, if semantic continuity is being claimed
- the explicit statement of which stronger routing class is still unproven

That is the minimum needed to stop "route exists" from impersonating "route is
trustworthy under pressure."

## Bottom line

The routing dream in this repo is not merely "more ingress flexibility."
It is:

- any healthy node can receive traffic
- the node can decide local versus remote honestly
- it can preserve the request meaning when forwarding
- and the operator no longer has to privately finish the story

The current stack is serious enough to make that dream plausible.
It is not yet proven enough to let the docs speak as if the dream already owns
the bad day.
