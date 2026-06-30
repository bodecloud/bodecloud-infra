# Current Compose Runtime

This page is the strongest live-runtime inventory in the repository.

If a sentence says something is part of the priority implementation now, it
needs to survive contact with:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- the fragments it actively includes from
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

Everything else is weaker evidence for current runtime claims.

That strength is also why this page is dangerous.
Runtime inventories are where complicated infra stacks most easily start
sounding more relieving than they really are.
Once a repo has enough services, networks, helpers, and middlewares, simple
presence starts getting mistaken for shared truth.

This page has to keep interrupting that move.

## What this page is and is not allowed to prove

This page is authoritative about:

- what the root runtime actually includes today
- which networks, configs, secrets, and service families are materially present
- which files still define the priority implementation
- which visible runtime facts shape the real burden-transfer problem

This page is not authoritative about:

- generic wrong-node success
- durable peer-forward behavior under backend loss
- shared placement truth
- auth or middleware continuity after peer handoff
- honest stateful failover semantics

This page is a live inventory.
It is not a distributed-systems completion certificate.

## Strongest honest current answer

The root runtime already proves a serious Compose-first platform.

It clearly proves:

- the root `docker-compose.yml` is still a major implementation surface, not a
  symbolic wrapper
- the repo is actively split across included fragments rather than pretending to
  be a single-file toy stack
- the stack already spans public ingress, auth, security, observability,
  private mesh, media, search, AI tooling, TCP services, and alternate egress
  experiments
- the runtime is complex enough that words like "failover," "anti-SPOF," and
  "peer-aware" need strict custody

It still does **not** prove the thing the user actually wants:

- shared current placement truth
- trustworthy peer eligibility truth
- durable wrong-node rescue paths
- preserved auth and middleware semantics after cross-node forwarding
- stateful failover semantics that are honest rather than merely reachable

That difference is the whole reason this page exists.

The stack is broad.
The stack is real.
The stack is still not the same thing as a truth-owning multi-node platform.

## The priority runtime really does center on the root file

The root file defines the visible control surface the user is trying not to
lose too early.

The root `docker-compose.yml` currently owns or directly declares:

- shared networks
- shared configs
- shared secrets
- direct services
- the include list for the active fragment stack

That matters because the repo's dream is explicitly not:

- demote the last readable layer immediately
- then ask the user to trust a more hidden control plane on faith

The root file remains a major truth surface because it still lets the operator
inspect a large part of reality directly.

## Active include set

The root file actively includes these fragments today:

- `compose/docker-compose.coolify-proxy.yml`
- `compose/docker-compose.docs.yml`
- `compose/docker-compose.firecrawl.yml`
- `compose/docker-compose.headscale.yml`
- `compose/docker-compose.llm.yml`
- `compose/docker-compose.metrics.yml`
- `compose/docker-compose.stremio-group.yml`
- `compose/docker-compose.warp-nat-routing.yml`
- `compose/docker-compose.wishlist.yml`

That include list already says something important about the runtime:

- the edge stack is not imaginary
- the docs surface is live
- the private mesh assumptions are live
- observability is live
- LLM and research tooling are live
- alternate egress experiments are live
- media and application workloads are live

What it does **not** say is where cross-node placement truth lives.
The include graph proves breadth.
It does not by itself prove distributed dignity.

## Root-owned networks

The root runtime defines three central shared networks:

| Network | Immediate visible role | What it honestly proves | What it still does not prove |
| --- | --- | --- | --- |
| `publicnet` | ingress-adjacent and public-facing traffic | public exposure is explicitly modeled in the runtime | that all public nodes share request-time service truth |
| `backend` | internal app and support traffic | the stack distinguishes internal traffic from public entry | that cross-node service ownership is system-owned |
| `warp-nat-net` | alternate egress and NAT-routed traffic | the repo is experimenting with non-default egress and routing behavior | that cross-node failover decisions are cluster-owned |

Those networks matter because they show the repo is already modeling different
traffic realities.
They do not prove the runtime can explain wrong-node decisions for itself.

## Shared runtime surfaces declared at the root

The root file also defines shared configs and secrets that matter operationally.

Examples visible directly in the root:

- shared secret `signing_secret`
- `watchtower-config.json` mounted from `~/.docker/config.json`
- session-manager assets
- Homepage inline config content

This proves the repo is already using Compose as a real control surface for:

- shared configuration
- inline-config preference
- secret-backed service behavior
- operator-facing tooling

It still does not prove that config centralization equals shared topology truth.
The repo can be highly declarative and still leave the operator privately
finishing cross-node decisions.

## Direct root services that matter to the architecture story

The root file itself visibly carries several important workloads, including:

- `mongodb`
- `dcef`
- `chat-analytics`
- `searxng`
- `code-server`
- additional direct services further down the root file

The architecture importance is not just that these exist.
It is that the root file still hosts:

- HTTP-facing services
- protected tooling surfaces
- TCP-exposed stateful services
- developer and operator surfaces

This means the priority runtime is not merely "the fragments."
The root file remains part of the real platform contract.

## Service families the runtime visibly contains today

The current runtime can be read more honestly by grouping it into service
families instead of a raw service dump.

### 1. Edge, ingress, and auth family

Most visibly anchored by `compose/docker-compose.coolify-proxy.yml`.

This family includes live pressure from tools such as:

- Traefik
- TinyAuth
- `nginx-traefik-extensions`
- CrowdSec
- Docker socket proxies
- `cloudflare-ddns`
- `docker-gen-failover`

What this family honestly proves:

- the repo has a serious HTTP edge stack
- policy, auth, and middleware are not afterthoughts
- DNS plurality and proxy experimentation are live concerns

What it still does not prove:

- generic wrong-node success
- trustworthy peer selection
- middleware continuity after cross-node forwarding
- route durability under the exact failure that makes fallback matter

### 2. Stateful core services family

Visible in the root file and fragments such as `compose/docker-compose.core.yml`
and the active root graph.

Examples include:

- MongoDB
- Redis
- Headscale-related state
- Firecrawl-supporting stateful surfaces

What this family honestly proves:

- the repo is not only about static frontends or disposable apps
- state-bearing services are already part of the priority runtime
- TCP routing pressure is already real

What it still does not prove:

- that these services are honestly HA
- that TCP reachability means safe substitution
- that state authority is no longer concentrated on a sacred node

### 3. Observability and operator family

Most visibly anchored by `compose/docker-compose.metrics.yml`, Homepage config,
and operator-facing dashboards.

What this family honestly proves:

- the stack is already instrumented and watched
- the operator has substantial runtime visibility surfaces
- the repo is serious about knowing what is alive

What it still does not prove:

- that observability equals convergence truth
- that alerts can tell the wrong node what to do
- that being able to see the failure is the same thing as the system owning the
  recovery truth

### 4. Mesh and routing-experiment family

Most visibly anchored by `compose/docker-compose.headscale.yml` and
`compose/docker-compose.warp-nat-routing.yml`.

What this family honestly proves:

- the repo is already experimenting with private-mesh communication
- non-default routing and egress are not theoretical
- the runtime is trying to shape traffic, not merely host containers

What it still does not prove:

- that mesh reachability equals eligible forwarding
- that alternate routing equals shared placement truth
- that any healthy node can preserve any route meaningfully

### 5. Workload and application family

Examples from active fragments include docs, wishlist, media, search, AI, and
developer tooling.

What this family honestly proves:

- the platform is already broad enough to create real routing and recovery
  pressure
- the bad-day problem is not hypothetical because many kinds of workloads are
  already in play

What it still does not prove:

- that breadth creates believable options
- that one route class's optimism transfers to another

## What the runtime still forces the operator to know privately

The current runtime may still force a human to know things like:

- which node actually hosts a named service right now when the request lands
  somewhere else
- whether a helper-generated rescue path would still exist after backend loss
- whether a reachable peer is merely alive or actually safe for that route's
  auth, middleware, secrets, and revision expectations
- whether a TCP-exposed stateful service is only transport-reachable or safe to
  describe with failover language
- whether a state-bearing workload still hides one sacred authority node behind
  a more distributed-looking ingress story

That list is a more faithful runtime summary than "the stack is still maturing."

The user is not saying the runtime lacks machinery.
They are saying too many decisive truths still have to be remembered,
inferred, or narrated by a human at exactly the moment a believable platform
should already be exposing them.

## What still does not count as runtime evidence

The following are all real facts and still not enough for the user's benchmark:

- many active services
- a sophisticated-looking edge stack
- multiple active fragments
- TCP routers on stateful services
- broad observability
- helpers with failover-shaped names
- more than one public-entry-looking component
- a route that can be rendered now but has not been stressed under the failure
  that would make it matter

Those facts prove the repo is operationally dense.
They do not prove the system, rather than the operator, owns the truth needed
on the bad day.

## What a runtime-backed progress packet would have to contain

Before this page supports stronger "the runtime has moved closer to the dream"
language, it should point to a packet containing:

- the exact live component or route being credited
- the narrower hidden burden it really reduces
- the runtime artifact proving that component is materially present now
- the adjacent burden it still does not reduce
- the next drill required before stronger language becomes legal

Without that packet, inventory tends to impersonate relief.

## Bottom line

The current Compose runtime is already serious enough to deserve careful,
evidence-first reading.

It proves:

- the root Compose surface matters
- the edge stack is real
- stateful and TCP pressure are real
- the repo is already broad enough that wrong-node truth is a live problem

It still does not prove:

- shared placement truth
- trustworthy peer eligibility
- durable route persistence under failure
- preserved request meaning after peer handoff
- honest stateful failover semantics

That is the honest runtime boundary today.
The stack is substantial.
The missing middle layer is still missing.
