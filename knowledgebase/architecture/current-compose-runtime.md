# Current Compose Runtime

This page is the closest thing the repo has to a present-tense runtime map.

That does **not** mean it is allowed to sound triumphant.
In this repository, runtime inventory is one of the easiest places to start
lying by accident. Once enough services, helper containers, middlewares,
healthchecks, and dashboards exist, the stack can start sounding unified long
before it actually owns the truths that matter on the bad day.

So this page has one job:

> say exactly what the priority Compose runtime really contains now, and refuse
> to upgrade that presence into cross-node dignity it has not yet earned.

## What this page is and is not allowed to prove

This page is authoritative about:

- what the root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  actively includes today
- which networks, configs, secrets, and service families materially exist in
  the tracked runtime
- which service categories are already exerting real routing, auth,
  observability, and state pressure
- which parts of the repo are clearly live enough that words like "failover,"
  "fallback," and "anti-SPOF" need disciplined handling

This page is not authoritative about:

- generic wrong-node success
- generic peer-forward success
- backend-loss survival for arbitrary routes
- shared placement truth
- shared peer-eligibility truth
- honest stateful failover semantics

This is a runtime page.
It is not a permission slip to say the platform already behaves like one
coherent personal cloud.

## The strongest honest current answer

The current runtime already proves a large, serious, Compose-first platform.

It proves all of the following:

- the root Compose file is still a real control surface, not just a bootstrap
  wrapper
- the priority implementation is distributed across active fragments, not
  hidden behind a separate orchestrator
- the edge stack is real
- auth and middleware are real
- observability is real
- stateful services are real
- alternative routing and egress experiments are real
- the repo has enough breadth that wrong-node truth is now a live pressure, not
  a theoretical one

It still does **not** prove the deeper thing the user wants:

- that any healthy public node already knows what the request should mean
- that the receiving node already knows whether to serve locally or forward
- that the receiving node already knows which peer is actually eligible
- that a protected route keeps the same auth and middleware meaning after
  handoff
- that the stateful layer is anything other than harshly mixed in maturity

That gap is not a rhetorical nuance.
That gap is the reason the repo still feels short on real options even while it
already contains a lot of machinery.

## Why the root file still matters so much

The dream in this repo is not merely "more automation."
It is "stop making the operator be the hidden registry and hidden fallback
brain."

That is why the root file still matters.
The root file is the last large human-readable place where the operator can
still inspect a significant portion of the system without first trusting a more
abstract controller.

The root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
still directly defines or anchors:

- shared networks
- shared secrets
- shared configs
- direct services
- the active include graph for the fragment stack

That means the root file is not just "legacy Compose."
It is still one of the repo's strongest live truth surfaces.

## The active include set is already a meaningful runtime fact

The root file actively includes all of these fragments today:

- `compose/docker-compose.coolify-proxy.yml`
- `compose/docker-compose.docs.yml`
- `compose/docker-compose.firecrawl.yml`
- `compose/docker-compose.headscale.yml`
- `compose/docker-compose.llm.yml`
- `compose/docker-compose.metrics.yml`
- `compose/docker-compose.stremio-group.yml`
- `compose/docker-compose.warp-nat-routing.yml`
- `compose/docker-compose.wishlist.yml`

That include set matters because it shows the runtime is not centered on a
single service class.
It is already carrying pressure from all of these at once:

- ingress and auth
- docs publishing
- search and developer surfaces
- AI gateways and tool routing
- stateful app support
- private mesh assumptions
- media workloads
- monitoring
- alternative egress shaping

What the include set does **not** tell you is where the decisive cross-node
truth lives when a request lands on the wrong machine.

The include graph proves breadth.
It does not prove convergence.

## Root-owned networks

The root runtime defines three especially important shared networks:

| Network | Immediate visible role | What it honestly proves | What it still does not prove |
| --- | --- | --- | --- |
| `publicnet` | public-facing and ingress-adjacent traffic | public exposure is explicitly modeled in the runtime | that all public nodes already share request-time service truth |
| `backend` | app-to-app and support traffic | the stack already distinguishes internal traffic from edge traffic | that cross-node ownership and forwarding decisions are system-owned |
| `warp-nat-net` | alternate egress and NAT-routed traffic | non-default traffic shaping is materially part of the platform | that egress shaping equals shared placement or fallback truth |

This matters because the repo is already trying to shape multiple traffic
realities at once.
It is not "just a bunch of containers."
But even a careful network model can still leave the operator holding the
actual answer to "what should happen now?"

## Root-owned secrets and configs

The root file also proves the stack is using Compose as a real configuration
surface rather than only as a service launcher.

Visible examples include:

- secret `signing_secret`
- config `watchtower-config.json`
- session-manager assets
- multiple inline Homepage configuration blobs

That matters because it confirms several repo values are already live in the
runtime:

- shared configuration is real
- inline config preference is real
- secret-backed behavior is real
- operator-facing surfaces are real

It still does not prove that configuration centrality equals topology truth.
The system can be very declarative and still make the operator privately finish
the cross-node sentence.

## Direct services in the root file

The root file still directly carries important workload surfaces instead of
hiding everything inside fragments.

The visible root services include at least:

- `mongodb`
- `dcef`
- `chat-analytics`
- `searxng`
- `code-server`
- `session-manager` behind a profile gate

Those matter because they show the root file is not just holding bootstrap
plumbing.
It still hosts real application and operator-facing surfaces:

- TCP-exposed stateful service: `mongodb`
- search surface: `searxng`
- developer/operator surface: `code-server`
- analysis workload: `chat-analytics`

That is why the repo cannot honestly be read as "everything important is in the
future middle layer."
The current Compose runtime is already the platform.

## Fragment-by-fragment runtime map

The cleanest way to understand the live runtime is not to dump every service
name.
It is to group the stack by the burden each fragment is trying to absorb.

### `compose/docker-compose.coolify-proxy.yml`

This fragment is the clearest live proof that the repo is already serious about
edge behavior.

It visibly carries:

- `cloudflare-ddns`
- `nginx-traefik-extensions`
- `tinyauth`
- `crowdsec`
- `crowdsec-init`
- `traefik`
- `whoami`
- `docker-gen-failover`
- `logrotate-traefik`
- `autokuma`

This fragment honestly proves:

- the repo is already modeling any-node public entry pressure
- auth continuity is already a real concern
- middleware continuity is already a real concern
- event-reactive routing helpers are already part of the runtime
- observability of the edge layer is not optional

This fragment still does not prove:

- generic wrong-node HTTP success
- generic peer eligibility truth
- preserved protected-route meaning after peer forwarding
- that a helper named `docker-gen-failover` actually solves the exact bad day
  that made the user ask for relief

### `compose/docker-compose.headscale.yml`

This fragment proves the private mesh assumptions are live, not imaginary.

It visibly carries:

- `headscale-server`
- `headscale`
- inline Headscale config

The active config also matters because it still points to SQLite under
`/var/lib/headscale/db.sqlite`.

This fragment honestly proves:

- private-node identity and connectivity are part of the real runtime story
- the repo is already treating private mesh assumptions as infrastructure, not
  just aspiration
- "reachable peer" is no longer a theoretical category

This fragment still does not prove:

- that reachability equals valid forwarding target
- that node identity equals service ownership truth
- that the control surface itself is no longer socially singular

### `compose/docker-compose.metrics.yml`

This fragment proves the stack is heavily instrumented and that the operator is
not flying blind.

It visibly carries:

- `victoriametrics`
- `prometheus`
- `grafana`
- `node_exporter`
- `cadvisor`
- `loki`
- `promtail`
- `blackbox-exporter`
- `mongodb-exporter`
- `redis-exporter`
- `process-exporter`
- `alertmanager`

It also provisions dashboards for many runtime surfaces, including:

- Traefik
- CrowdSec
- TinyAuth
- `cloudflare-ddns`
- `nginx-traefik-extensions`
- `docker-gen-failover`
- `headscale`
- `headscale-server`
- `watchtower`
- `homepage`
- `code-server`
- `searxng`
- `firecrawl`
- `mongodb`
- `redis`
- WARP-related services

This fragment honestly proves:

- the operator can already observe a lot of the stack
- the repo cares about route, container, log, and exporter visibility
- many named services are mature enough to be first-class observability
  subjects

This fragment still does not prove:

- that the system owns recovery truth instead of only exposing symptoms
- that alerts can tell the wrong node what the request should mean
- that visibility equals burden transfer

### `compose/docker-compose.warp-nat-routing.yml`

This fragment proves the repo is not only hosting services.
It is also deliberately shaping egress and routing behavior.

It visibly carries:

- `warp-net-init`
- `warp-nat-gateway`
- `warp_router`
- `ip-checker-warp`
- inline configs `warp-nat-setup.sh` and `warp-monitor.sh`

This fragment honestly proves:

- non-default egress routing is already part of the real platform
- helper logic is already manipulating network behavior materially
- WARP availability is treated as a runtime condition worth monitoring and
  reacting to

This fragment still does not prove:

- that alternate egress truth equals shared service placement truth
- that network repair logic is the same thing as cross-node request rescue
- that a route can preserve application meaning just because packets can still
  leave somewhere

### `compose/docker-compose.firecrawl.yml`

This fragment carries a stateful application family rather than just edge
plumbing.

It visibly carries:

- `playwright-service`
- `firecrawl`
- `nuq-postgres`
- `rabbitmq`

This fragment honestly proves:

- the runtime already includes application stacks with mixed HTTP, worker, and
  stateful support surfaces
- queueing and persistence pressure are already real

This fragment still does not prove:

- that those stateful support components are HA
- that service presence implies safe failover semantics

### `compose/docker-compose.llm.yml`

This fragment proves the repo is already trying to operate a tool-routing and
gateway-heavy AI surface, not just simple web apps.

It visibly carries:

- `open-webui`
- `mcpo`
- `model-updater`
- `litellm`
- `litellm-postgres`
- `gptr`
- `qdrant`
- `mcp-proxy`

This fragment honestly proves:

- the runtime already contains systems that depend on routing, tool exposure,
  caching, persistence, and external API coordination
- provider and gateway failover are already part of the app-level vocabulary

This fragment still does not prove:

- that node-level request preservation is solved just because app-level gateway
  failover exists
- that the AI/tooling surface has escaped hidden placement truth any more than
  the rest of the stack has

### `compose/docker-compose.stremio-group.yml`

This fragment proves the repo is already hosting a large media/service bundle
with many secrets and helper surfaces.

It visibly carries:

- `stremio`
- `flaresolverr`
- `jackett`
- `prowlarr`
- `aiostreams`
- `comet`
- `mediafusion`
- `mediaflow-proxy`
- `stremthru`
- `rclone-init`
- `rclone`

This fragment honestly proves:

- the stack already spans complicated application families rather than a small
  lab demo
- routing, auth, media proxying, and external-provider coupling are already
  real platform concerns

This fragment still does not prove:

- that breadth of workloads somehow creates unified placement truth
- that resilience conclusions from one workload class transfer automatically to
  another

### `compose/docker-compose.docs.yml` and `compose/docker-compose.wishlist.yml`

These smaller fragments still matter because they provide narrow, concrete HTTP
surfaces that can be used as proof candidates.

They visibly carry:

- `mkdocs`
- `wishlist`

These fragments honestly prove:

- the repo has simple stateless HTTP surfaces suitable for early wrong-node
  drills
- not every proof has to start with the hardest stateful or auth-heavy system

They still do not prove:

- that the protected edge behaves correctly
- that passing a small stateless route test upgrades the whole platform

## Service families that are definitely live now

If the runtime is grouped by burden rather than by file, the live platform
already contains all of these families:

### 1. Edge and policy family

Examples:

- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `cloudflare-ddns`
- `docker-gen-failover`

This family proves:

- the repo already treats HTTP entry, policy, and protection as first-class
  runtime concerns

This family does not prove:

- that the wrong node can preserve a protected route correctly

### 2. Operator and observability family

Examples:

- `homepage`
- `grafana`
- `prometheus`
- `victoriametrics`
- `alertmanager`
- `autokuma`
- exporters and dashboards

This family proves:

- the operator can already see a lot of runtime reality

This family does not prove:

- that the system can already act coherently without private operator
  completion

### 3. Stateful core family

Examples:

- `mongodb`
- `redis`
- Headscale SQLite
- `nuq-postgres`
- `rabbitmq`
- `litellm-postgres`
- `qdrant`

This family proves:

- the repo is already carrying real write-authority and persistence pressure

This family does not prove:

- that these services are honestly multi-writer, replicated, or failover-safe

### 4. Routing and egress experiment family

Examples:

- `warp-nat-gateway`
- `warp_router`
- WARP monitor/setup logic

This family proves:

- the platform is already trying to own network behavior, not merely container
  hosting

This family does not prove:

- that network cleverness has already become shared request meaning

### 5. Stateless HTTP proof-candidate family

Examples:

- `whoami`
- `wishlist`
- `mkdocs`

This family proves:

- the repo has tractable early candidates for narrow wrong-node drills

This family does not prove:

- that the harder protected or stateful paths are therefore solved

## What the runtime still forces a human to know privately

This is the part that matters most.

The current Compose runtime may still force the operator to know things such as:

- which node actually hosts a named service right now
- whether the request is being served locally or via a helper-driven rescue path
- whether the reachable peer is merely alive or actually valid for that route
- whether the forwarded route still preserves auth, middleware, secrets, and
  app revision expectations
- whether a stateful service is transport-reachable but still singular in all
  the ways that matter

That list is a better summary of the missing platform layer than any sentence
about "maturity."

The repo is not short on runtime components.
It is still short on system-owned shared truth.

## What still does not count as runtime evidence

The following are all real runtime facts and still not enough for the user's
benchmark:

- many active services
- multiple public-looking surfaces
- healthchecks
- helper containers with failover-shaped names
- a sophisticated edge stack
- dashboards for nearly everything
- TCP routers on stateful services
- multiple A records at DNS
- private mesh identity
- route success on a normal day

Those facts prove the repo is operationally dense.
They do not prove the platform knows enough to stop humiliating the operator on
the bad day.

## What a stronger runtime progress packet would need

Before this page can legally support stronger language like "the runtime has
moved closer to the dream," it needs a narrow packet containing:

- the exact route or service being credited
- the exact hidden burden that route or service now removes
- the runtime artifact proving the component is materially present
- the drill proving behavior under pressure
- the explicit ceiling on what still remains unproven

Without that packet, runtime inventory turns into comfort prose.

## Bottom line

The current Compose runtime is substantial.
It is real.
It is already carrying serious edge, routing, state, policy, and observability
pressure.

It honestly proves:

- the root Compose surface still matters
- the fragment graph is live
- the edge stack is live
- the mesh assumptions are live
- the monitoring stack is live
- the stateful layer is materially present
- the wrong-node problem is now a real platform problem, not a hypothetical one

It still does not honestly prove:

- shared placement truth
- shared peer-eligibility truth
- generic wrong-node success
- generic backend-loss route survival
- preserved protected-route meaning after peer handoff
- honest stateful failover semantics

That is the live runtime boundary today.
The stack is broad enough to deserve respect.
The missing layer is still the layer that would let the system answer the hard
question for itself.
