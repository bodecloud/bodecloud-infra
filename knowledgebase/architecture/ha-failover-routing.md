# HA, Failover, and Routing

Read this page as the routing truth map for the priority implementation rooted
at [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

If you want the deeper evidence stack first, read:

- [`../research/ingress-and-failover-evidence.md`](../research/ingress-and-failover-evidence.md)
- [`request-path-and-failure-walkthrough.md`](request-path-and-failure-walkthrough.md)
- [`current-compose-runtime.md`](current-compose-runtime.md)

This page exists because "HA" becomes meaningless almost immediately in this
repo unless routing is split into the layers the user is actually angry about.

The user is not asking:

- can more than one node receive traffic?
- is Traefik present?
- are there healthchecks?
- is there a helper that sounds failover-shaped?

The user is asking:

> if traffic lands on the wrong healthy machine, can that machine still
> preserve the service contract without requiring private operator memory?

That is the routing question here.

## What this page is and is not allowed to prove

This page is authoritative about:

- how routing has to be decomposed before "failover" means anything useful
- which routing surfaces are materially live in the current Compose runtime
- which routing truths are still missing, social, or only planned
- why first-hop plurality is much weaker than request-preserving recovery

This page is not authoritative about:

- whether a specific hostname has already passed a real wrong-node drill
- whether backend-loss fallback is already durable for a named route
- whether stateful traffic has inherited HTTP optimism
- whether one working helper path upgrades the whole platform

## Strongest honest current answer

The root runtime already has a serious edge stack and enough moving parts to
justify real routing work: Cloudflare-oriented public entry assumptions,
Traefik, TinyAuth, CrowdSec, Nginx auth extensions, Docker-socket proxies, and
`docker-gen-failover` are all materially present. What is still missing is the
thing the user actually wants that machinery to buy: shared placement truth,
trustworthy peer eligibility, durable fallback-route persistence, and proof
that a request keeps meaning the same thing after wrong-node handoff.

That is the whole difference between:

- a stack that can route many things locally
- and a stack that stops gambling on node locality

## The current routing surface in the priority implementation

The priority root runtime is not abstract. It is the merged surface created by
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
plus included fragments such as:

- [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- [`compose/docker-compose.docs.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.docs.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- [`compose/docker-compose.metrics.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.metrics.yml)

At the root level, the merged graph already defines three important network
surfaces:

- `publicnet`
- `backend`
- `warp-nat-net`

That alone does not prove resilience, but it does prove the routing problem is
already materially encoded in the runtime rather than only imagined in plans.

Within the edge fragment, the current routing surface includes:

- `cloudflare-ddns`
- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `crowdsec-init`
- `docker-gen-failover`
- `dockerproxy-ro`
- `dockerproxy-rw`

Within the priority runtime, there are already concrete service classes that
stress routing differently:

- stateless or near-stateless HTTP surfaces such as docs, dashboards, and
  utility frontends
- protected HTTP surfaces that depend on TinyAuth, Nginx auth middleware, and
  CrowdSec behavior staying coherent
- raw TCP surfaces such as MongoDB, which is already routed by Traefik TCP
  labels in the root graph
- stateful control-plane surfaces such as Headscale, whose single-node reality
  is already openly acknowledged in the master plan

This matters because "routing" here is not one problem.

## The routing layers that must stay separate

If these layers get collapsed into one sentence called "HA," the docs become
decorative again.

### 1. Public node-entry reachability

Question:

- can a client reach some healthy public node at all?

What the repo already has:

- Cloudflare-oriented public entry assumptions
- `cloudflare-ddns`
- a long-standing goal of multi-node public entry rather than one sacred box
- repo-native intent in [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  that explicitly says "User -> Cloudflare DNS -> any surviving node"

What this layer can honestly buy:

- more than one node can plausibly receive the first hop
- ingress does not have to be psychologically concentrated on one sacred host
- some node loss at the public edge may be survivable

What this layer does not buy:

- service locality truth
- peer eligibility truth
- fallback-route persistence
- policy continuity
- stateful correctness

This is why "multiple A records" is not even close to the final answer.

### 2. Local edge-stack health

Question:

- once traffic reaches a node, is the local edge stack coherent enough to make
  the next decision honestly?

What the repo already has:

- Traefik as the real L7 execution surface
- TinyAuth for forward-auth behavior
- `nginx-traefik-extensions` as additional auth and middleware glue
- CrowdSec as the active security and filtering surface
- Docker provider plus file provider in Traefik
- healthchecks across key edge components

What this layer can honestly buy:

- local route execution
- local auth and middleware handling
- local TLS and certificate handling
- serious edge behavior rather than one symbolic reverse proxy mention

What it does not buy:

- cross-node knowledge
- trustworthy peer choice
- proof that local and forwarded behavior remain semantically identical

If this layer is broken, distributed entry above it is useless.
If this layer is healthy, the distributed problem is still not solved.

### 3. Locality truth

Question:

- does the receiving node actually know whether the requested service is local?

What the repo already has:

- local Docker labels
- local Docker provider visibility
- a clear architectural desire for a tracked current-state registry such as
  `services.yaml`

What the repo does not yet prove:

- a live tracked root [`services.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/services.yaml)
  consumed by routing decisions
- a shared placement-truth surface that outranks operator memory

Why this matters:

The recurring `services.yaml` pressure is not about loving files. It is about
stopping "the operator remembers where the service really lives" from being the
real control plane.

Without locality truth, a wrong-node request cannot become an honest decision.

### 4. Peer-selection truth

Question:

- if the service is not local, does the receiving node know which peer is
  actually valid right now?

This is stricter than:

- can nodes talk over Headscale or Tailscale?
- does a peer hostname exist?
- did the same service run on that peer last week?

What the repo already has:

- Headscale as a materially live mesh/control-plane component in
  [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- explicit planning around peer broadcast, leader election, and node-aware
  coordination in [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

What the repo does not yet prove:

- a live peer-eligibility decision surface consumed by the edge layer
- a trustworthy answer to "which peer should I use now?" that is stronger than
  static config or social knowledge

Mesh reachability is helpful.
Mesh reachability is not current truth ownership.

### 5. Fallback-route persistence

Question:

- when the preferred backend disappears, does the route needed for recovery
  still exist?

This is one of the hardest and most important seams in the current runtime.

The repo already knows about a specific live weakness:

- `docker-gen-failover` is present in the edge stack
- the master plan explicitly records that the current approach deletes routes
  when a container stops, which defeats the point of fallback

This is not a minor caveat.
It is exactly the kind of failure the user is trying to stop:

> the system looked dynamic until the bad event arrived, then the dynamic route
> vanished with the thing it was supposed to route around

That is why the docs keep treating `docker-gen-failover` as both meaningful and
dangerous:

- meaningful because it shows the repo is trying to generate fallback-aware
  Traefik config
- dangerous because the current generation model is recorded as losing routes
  at the wrong moment

### 6. Policy continuity

Question:

- if the request is handed to a peer, does it still behave like the same
  service from the user's perspective?

In this repo, routing correctness includes:

- auth continuity
- middleware continuity
- visible policy continuity
- header and rewrite behavior continuity

That is why the edge stack is not "just Traefik."
It is Traefik plus TinyAuth plus Nginx auth extensions plus CrowdSec plus file
provider plus helper-generated dynamic config.

The user is not asking for "some response happened."
The user is asking whether the same protected route stays the same protected
route after handoff.

### 7. Service-class honesty

Question:

- are we telling the truth about what kind of traffic we are talking about?

The root runtime already proves that service classes differ sharply:

- HTTP services can at least plausibly move toward local-first then peer-forward
- TCP services such as MongoDB already exist in the root graph, but TCP
  forwarding does not equal stateful safety
- Headscale is a real example of a control-plane singleton that the docs
  already refuse to narrate as magically HA

This is why the routing story has to fork by service class:

- stateless HTTP can plausibly earn the first serious wrong-node proof
- protected HTTP must additionally prove auth and middleware continuity
- TCP needs separate transport and client-behavior proof
- stateful systems need ownership, replication, promotion, and reconnect truth

## What the current stack can honestly claim today

The current stack can honestly claim:

- the root Compose graph already contains a real multi-surface ingress system
- Traefik is an actual execution surface, not a decorative proxy mention
- Cloudflare-style any-node entry is explicit architecture intent
- Headscale, auth, middleware, and observability are part of the real runtime
- the repo has already identified concrete failover blockers rather than merely
  gesturing at HA

The current stack cannot honestly claim:

- generic wrong-node request preservation
- durable backend-loss fallback
- live shared placement truth
- trustworthy peer eligibility truth
- middleware and auth continuity under peer fallback
- stateful HA just because traffic can still reach a port

## The routing claim that would actually matter

The claim that would materially satisfy the user is not:

- "Traefik is configured"
- "Cloudflare points at multiple places"
- "Headscale exists"
- "a helper can generate dynamic config"

It is:

> any surviving public node can receive a request, determine whether the target
> is local, preserve the request if it is not, survive the failure that made
> fallback necessary, and keep auth, middleware, and operator readability
> intact throughout that handoff

Anything smaller may still be good engineering.
It is simply not the full routing problem this repo is trying to solve.

## What should happen next

The next routing work that actually pays down the user's pain is narrower than
"pick an orchestrator" and stronger than "add more labels."

The priority sequence is:

1. expose or introduce a live placement-truth surface strong enough to answer
   "what runs where right now?"
2. prove one stateless HTTP route locally on the honest hosting node
3. prove that same route through intentional wrong-node entry
4. prove whether the route survives when the preferred local backend actually
   disappears
5. compare local versus peer-forwarded behavior for auth and middleware parity
6. only then narrate stronger routing maturity

That is the smallest honest path out of fake HA language for the priority
implementation.
