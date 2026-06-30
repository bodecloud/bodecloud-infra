# Current Compose Runtime

This page is the present-tense runtime inventory for the priority
implementation rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

It exists because this repo now has enough real machinery to produce a very
specific kind of lie:

> the stack is already broad, instrumented, routed, authenticated, and
> helper-driven, so it must already be close to acting like one coherent
> multi-node personal cloud

That does not follow.

This page is therefore not trying to sound impressive.
It is trying to say exactly what the live runtime contains, where the evidence
is strongest, and what the reader is still forbidden to upgrade from mere
presence into cross-node dignity.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what the root Compose runtime actively includes now
- which networks, configs, secrets, and service families are definitely in the
  tracked runtime
- which named surfaces are materially live enough that words like "failover,"
  "fallback," and "anti-SPOF" must now be handled with discipline
- which runtime burdens are already real rather than hypothetical

This page is not allowed to prove:

- generic wrong-node success
- generic peer-forward success
- backend-loss survival for arbitrary routes
- shared placement truth
- shared peer-eligibility truth
- honest stateful failover semantics

This page inventories the runtime.
It does not certify the dream.

## Strongest honest current answer

The current runtime already proves a serious Compose-first platform.

It proves all of these are materially present:

- a broad root Compose control surface
- an active fragment graph
- real edge routing
- real auth and middleware
- real observability
- real private-mesh assumptions
- real stateful services
- real helper containers with failover-shaped jobs
- real alternate routing and egress logic

It still does not prove the decisive thing the user is actually asking for:

- that any healthy public node already knows what a request should mean
- that the receiving node already knows whether to serve locally or forward
- that the receiving node already knows which peer is currently eligible
- that a protected route keeps the same auth and middleware meaning after a
  handoff
- that stateful services have escaped singular write-authority truth

That gap is the whole reason the repo can feel both serious and still
insufficient at the same time.

## Root runtime facts that are definitely live

### Active include set

The root runtime actively includes:

- `compose/docker-compose.coolify-proxy.yml`
- `compose/docker-compose.docs.yml`
- `compose/docker-compose.firecrawl.yml`
- `compose/docker-compose.headscale.yml`
- `compose/docker-compose.llm.yml`
- `compose/docker-compose.metrics.yml`
- `compose/docker-compose.stremio-group.yml`
- `compose/docker-compose.warp-nat-routing.yml`
- `compose/docker-compose.wishlist.yml`

That proves the runtime is already spanning several domains at once.
It does not prove those domains have converged on shared request-time truth.

### Root networks

The root runtime defines:

- `publicnet`
- `backend`
- `warp-nat-net`

These prove the runtime already distinguishes between public exposure, internal
traffic, and alternate routed/egress behavior.
They do not prove the system has a shared answer to "which node should serve
this request right now?"

### Root configs and secrets

Visible root-level config and secret surfaces include:

- secret `signing_secret`
- config `watchtower-config.json`
- session-manager assets
- multiple inline Homepage configuration blobs

That proves Compose is acting as a real control surface rather than just a
launcher.
It does not prove that the control surface has escaped private operator memory
for placement and fallback truth.

### Root-owned services

Visible root services include at least:

- `mongodb`
- `redis`
- `dcef`
- `chat-analytics`
- `searxng`
- `code-server`
- `homepage`
- `watchtower`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `dozzle`
- `portainer`
- `dns-server`
- `telemetry-auth`
- `bolabaden-nextjs`
- `biodecompwarehouse`
- `biodecompwarehouse-mcp`
- `biodecompwarehouse-bsim-server`
- `biodecompwarehouse-aio`
- profile-gated `session-manager`

That proves the root file is still one of the main live truth surfaces in the
repo.
It is not a bootstrap wrapper around the "real" system.

## Runtime evidence that already matters for the anti-SPOF conversation

The repo is already dense enough that certain runtime details are no longer
small implementation trivia.
They are part of the live proof boundary.

### Protected HTTP routes are already real

The root runtime already shows protected operator surfaces using
`nginx-auth@file`, including:

- `code-server`
- `dozzle`
- `portainer`
- metrics-side surfaces such as `prometheus`, `cadvisor`, and `alertmanager`

This matters because the user is not only asking for packet delivery.
They are asking for preserved request meaning.

So the real question is not only:

- can another node answer somehow?

It is also:

- does the route preserve the same auth and middleware contract when it stops
  being local?

The live runtime proves those protected routes exist.
It does not yet prove that peer rescue preserves their meaning.

### Raw TCP exposure is already real

The root runtime already exposes TCP services through Traefik, including:

- `mongodb`
- `redis`
- multiple `biodecompwarehouse*` surfaces

This matters because the repo is already carrying L4 and stateful pressure,
not just friendly stateless HTTP pressure.

It also forces a stricter documentation boundary:

- TCP reachability is not stateful authority
- passthrough routing is not promotion logic
- clean TLS/TCP exposure is not failover truth

### Edge helper growth is already real

The active edge stack already includes helper surfaces such as:

- `docker-gen-failover`
- `cloudflare-ddns`
- `autokuma`
- `logrotate-traefik`
- `nginx-traefik-extensions`

The `docker-gen-failover` container is especially important because it writes a
dynamic file to:

- `${CONFIG_PATH:-./volumes}/traefik/dynamic/failover-fallbacks.yaml`

using Docker events and label filtering.

That proves the repo is already growing helper behavior around Compose.
It does not prove that the helper owns trustworthy route survival under the
exact backend-loss event that makes the user angry.

### Headscale-backed private reachability is already real

The active Headscale fragment still points to SQLite:

- `database.type: sqlite`
- `sqlite.path: /var/lib/headscale/db.sqlite`

and exposes both:

- `headscale-server`
- `headscale`

This proves private mesh identity and reachability are live assumptions in the
runtime.
It does not prove that reachability, identity, and service validity have
become one shared truth.

### Observability of edge and helper surfaces is already real

The metrics stack already provisions dashboards and probes for many surfaces,
including:

- `traefik`
- `crowdsec`
- `tinyauth`
- `whoami`
- `docker-gen-failover`
- `nginx-traefik-extensions`
- `headscale`
- `headscale-server`
- `watchtower`
- `homepage`
- `code-server`
- `searxng`
- `mongodb`
- `redis`
- WARP-related services

That proves the operator can already observe a lot of the runtime.
It does not prove the runtime can already make the right recovery decision
without the operator privately finishing the story.

## Fragment-by-fragment runtime ledger

The best way to read the current runtime is not as one flat service dump.
It is as a ledger of which burden each fragment is already trying to absorb.

### `compose/docker-compose.coolify-proxy.yml`

This fragment visibly carries:

- `traefik`
- `tinyauth`
- `crowdsec`
- `crowdsec-init`
- `cloudflare-ddns`
- `nginx-traefik-extensions`
- `whoami`
- `docker-gen-failover`
- `logrotate-traefik`
- `autokuma`

It also shows real forward-auth and middleware pressure.

This fragment honestly proves:

- the edge stack is real
- auth continuity matters in the live runtime
- middleware continuity matters in the live runtime
- helper-driven route generation is already part of the implementation
- "wrong node" is no longer a theoretical concern

It still does not prove:

- generic wrong-node HTTP success
- generic protected-route continuity after peer forwarding
- shared placement truth
- shared peer-eligibility truth

### `compose/docker-compose.headscale.yml`

This fragment visibly carries:

- `headscale-server`
- `headscale`
- active inline Headscale configuration

It honestly proves:

- private mesh assumptions are live
- reachable private peers are part of the present platform story
- the multi-node dream is already relying on more than public ingress

It still does not prove:

- that reachability equals valid forwarding target
- that node identity equals service ownership truth
- that the control surface itself is no longer socially singular

### `compose/docker-compose.metrics.yml`

This fragment visibly carries:

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

It also provisions dashboards and probe/alert coverage for many core runtime
surfaces.

It honestly proves:

- the operator is not flying blind
- edge, helper, datastore, and app surfaces are already observable enough to
  be named first-class operational subjects

It still does not prove:

- automated recovery truth
- route-preserving rescue behavior
- that visibility equals burden transfer

### `compose/docker-compose.warp-nat-routing.yml`

This fragment visibly carries:

- `warp-net-init`
- `warp-nat-gateway`
- `warp_router`
- `ip-checker-warp`
- inline setup and monitoring scripts

It honestly proves:

- the repo is already trying to own network behavior, not just app hosting
- non-default egress shaping is part of the live platform
- helper logic is already materially affecting packet paths

It still does not prove:

- shared service placement truth
- request-preserving peer rescue
- that route repair at the network layer equals preserved application meaning

### `compose/docker-compose.firecrawl.yml`

This fragment visibly carries:

- `playwright-service`
- `firecrawl`
- `nuq-postgres`
- `rabbitmq`

It honestly proves:

- queue-backed and state-bearing application families are already part of the
  runtime
- worker and persistence pressure are already real

It still does not prove:

- HA truth for those support services
- safe stateful failover semantics

### `compose/docker-compose.llm.yml`

This fragment visibly carries:

- `open-webui`
- `mcpo`
- `model-updater`
- `litellm`
- `litellm-postgres`
- `gptr`
- `qdrant`
- `mcp-proxy`

It honestly proves:

- gateway-heavy and tool-routing-heavy workloads are already live
- auxiliary state, provider coordination, and tool exposure are already part
  of the platform

It still does not prove:

- that app-level gateway failover solves node-level wrong-node behavior
- that the AI/tool surface has escaped hidden placement truth

### `compose/docker-compose.stremio-group.yml`

This fragment visibly carries a heterogeneous workload cluster including:

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

It honestly proves:

- the repo is already operating large application bundles with external
  dependency churn, secrets, storage, and routing pressure
- this is not a toy lab stack

It still does not prove:

- unified placement truth
- transferable resilience semantics across workload families

### `compose/docker-compose.docs.yml` and `compose/docker-compose.wishlist.yml`

These smaller fragments visibly carry:

- `mkdocs`
- `wishlist`

They honestly prove:

- the repo has narrow stateless HTTP surfaces that can serve as early proof
  candidates for a real wrong-node drill

They still do not prove:

- that the protected routes are solved
- that a narrow HTTP win upgrades the harder lanes

## Service families that are definitely live

If the runtime is grouped by burden rather than file location, the live stack
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

- ingress, policy, protection, and helper-driven route shaping are already
  first-class runtime concerns

This family does not prove:

- that a wrong-node request preserves the same route semantics

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

- replication
- promotion
- rediscovery
- honest failover semantics

### 4. Routing and egress experiment family

Examples:

- `warp-nat-gateway`
- `warp_router`
- WARP monitor/setup logic

This family proves:

- the platform is already trying to own network behavior, not merely start
  containers

This family does not prove:

- that network cleverness has become shared request meaning

### 5. Stateless HTTP proof-candidate family

Examples:

- `whoami`
- `wishlist`
- `mkdocs`

This family proves:

- the repo has tractable early candidates for narrow wrong-node drills

This family does not prove:

- that protected or stateful routes are therefore solved

## What the runtime still forces a human to know privately

This is still the most important runtime summary.

The live Compose platform may still force the operator to know things such as:

- which node actually hosts a named service right now
- whether the current request is being served locally or via a helper-driven
  rescue path
- whether the reachable peer is merely alive or actually valid for that route
- whether the forwarded route preserves auth, middleware, and app-expectation
  meaning
- whether a stateful service is merely transport-reachable while still being
  singular in the ways that matter

That list explains the missing layer better than any vague sentence about
"maturity."

## What still does not count as runtime proof

These are all real runtime facts and still not enough for the user's actual
benchmark:

- many active services
- multiple public-looking surfaces
- healthchecks
- helper containers with failover-shaped names
- a sophisticated edge stack
- dashboards for nearly everything
- TCP routers for stateful services
- multiple public DNS targets
- private mesh identity
- route success on a normal day

Those facts prove density.
They do not prove that the system has stopped humiliating the operator on the
bad day.

## What a stronger runtime progress packet would need

Before this page can support a stronger sentence like "the runtime has moved
closer to the dream," it needs a packet containing:

- the exact route or service being credited
- the exact hidden burden that route or service now removes
- the runtime artifact proving the component is materially present
- the drill proving behavior under pressure
- the explicit ceiling on what still remains unproven

Without that packet, runtime inventory becomes comfort prose.

## Bottom line

The current Compose runtime is substantial.
It is real.
It is already carrying serious edge, routing, state, policy, observability,
and network-behavior pressure.

It honestly proves:

- the root Compose surface still matters
- the fragment graph is live
- the edge stack is live
- the mesh assumptions are live
- the monitoring stack is live
- the stateful layer is materially present
- the wrong-node problem is now a real platform problem rather than a
  hypothetical one

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
