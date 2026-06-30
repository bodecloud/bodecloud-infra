# Current Compose Runtime

This page is the present-tense runtime inventory for the priority
implementation rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

It exists because the repo now has enough real machinery to produce a very
specific lie:

> the stack is already broad, instrumented, routed, authenticated, and
> helper-driven, so it must already be close to acting like one coherent
> multi-node personal cloud.

That does not follow.

This page is therefore not trying to sound impressive.
It is trying to say exactly what the live runtime contains, where the evidence
is strong, and what the reader is still forbidden to upgrade from mere presence
into cross-node dignity.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what the root Compose runtime actively includes now
- which service families, networks, configs, and secrets are definitely live in
  the tracked runtime
- which runtime burdens are already real rather than hypothetical
- which surfaces are serious enough that failover language must now become more
  disciplined

This page is not allowed to prove:

- generic wrong-node success
- generic peer-forward success
- backend-loss survival for arbitrary routes
- shared placement truth
- shared peer-eligibility truth
- honest stateful authority transfer

This page inventories the runtime.
It does not certify the dream.

## Strongest honest current answer

The current runtime already proves a serious Compose-first platform.

It materially proves the presence of:

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
- that a protected route keeps the same auth and middleware meaning after
  handoff
- that stateful services have escaped singular authority truth

That is why the runtime can feel both serious and still insufficient at the
same time.

## What the runtime already proves about the user's dream

The runtime is strong enough to prove that the user is not imagining a platform
from nothing.

It proves the repo already has:

- enough edge machinery that wrong-node behavior is a real concern
- enough protected surfaces that middleware continuity is no longer optional
- enough TCP and stateful surfaces that sloppy HA language becomes dangerous
- enough helpers and observability that "we have no options" is no longer the
  real problem

The frustration is not lack of components.
It is lack of smaller honest control surfaces that turn those components into
inspectable distributed decisions.

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

This proves the runtime already spans several domains at once.
It does not prove those domains have converged on shared request-time truth.

### Root networks

The root runtime defines:

- `publicnet`
- `backend`
- `warp-nat-net`

This proves the runtime already distinguishes public exposure, internal traffic,
and alternate routing or egress behavior.
It does not prove the system has one shared answer to:

> which node should serve this request right now?

### Root configs and secrets

Visible root-level config and secret surfaces include:

- secret `signing_secret`
- config `watchtower-config.json`
- session-manager assets
- multiple inline Homepage configuration blobs

This proves Compose is already acting as a real control surface rather than
just a launcher.
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

This proves the root file is still one of the main live truth surfaces in the
repo.
It is not a mere bootstrap wrapper around some other hidden platform.

## Runtime evidence that already matters for the anti-SPOF conversation

The runtime is already dense enough that some details are no longer small
implementation trivia.
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

So the live question is not only:

- can another node answer somehow?

It is also:

- does the route preserve the same auth and middleware contract when it stops
  being local?

The runtime proves those protected routes exist.
It does not yet prove that peer rescue preserves their meaning.

### Raw TCP exposure is already real

The root runtime already exposes TCP services through Traefik, including:

- `mongodb`
- `redis`
- multiple `biodecompwarehouse*` surfaces

This matters because the repo is already carrying L4 and stateful pressure, not
just friendly stateless HTTP pressure.

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

`docker-gen-failover` matters especially because it writes a dynamic file under:

- `${CONFIG_PATH:-./volumes}/traefik/dynamic/failover-fallbacks.yaml`

using Docker events and label filtering.

This proves the repo is already growing helper behavior around Compose.
It does not prove the helper owns trustworthy route survival under the exact
backend-loss event that makes the user angry.

This is one of the clearest examples of why the docs must stay strict.
The runtime already contains enough moving parts to impersonate a solution.

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

### Observability is already real, but observability is not ownership

The metrics stack already provisions dashboards and probes for many runtime
surfaces, including edge, auth, helper, and stateful systems.

That proves the repo is serious enough to inspect itself.
It does not prove the repo is yet serious enough for the system to answer its
own hardest distributed questions.

This distinction matters because dashboards are one of the easiest ways to
mistake visible machinery for transferred burden.

## The runtime burden map already visible from inventory alone

Even before the repo proves wrong-node or failover behavior, the inventory
already forces these burden classes into the open:

| Burden class | Why runtime presence already makes it real | Why presence still is not enough |
| --- | --- | --- |
| Wrong-node request meaning | multiple fragments and edge surfaces already imply nontrivial route interpretation | the node still may not know what the request should mean |
| Protected-route continuity | auth and middleware are already attached to real routes | handoff may still change route meaning |
| TCP pressure | L4 surfaces already exist | transport still is not authority |
| Stateful pressure | databases, queues, and vector stores are already live | persistence still is not failover dignity |
| Helper-truth pressure | helper daemons are already shaping behavior | helper presence still is not proof of correctness |

## What the runtime is still forbidden from implying

Because the runtime is already broad, a reader may be tempted to upgrade
presence into stronger claims.

The following upgrades are still forbidden:

- "Traefik is present, therefore wrong-node routing must be close"
- "Cloudflare is plural, therefore first-hop resilience must be close"
- "`docker-gen-failover` exists, therefore backend-loss rescue must be close"
- "Headscale is live, therefore peer truth must be close"
- "Redis and MongoDB are exposed, therefore stateful anti-SPOF must be close"

Those are exactly the kinds of upgrades this page is supposed to stop.

## Bottom line

The current runtime proves that the repo is serious enough for the anti-SPOF
dream to matter now.

It does not prove that the dream is already partially owned in the most
important ways.

The correct present-tense reading is:

- broad runtime
- real pressure
- real helper growth
- real protected and stateful surfaces
- still missing the shared truth that would let the healthy wrong node stop
  needing a human narrator
