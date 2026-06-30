# Compose Fragment Map

This page reconstructs how the priority implementation is actually assembled
from the root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and the active files under
[`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/).

It is not here to admire modularity.
It is here to answer the harder operator question:

> when someone says "the stack," which files are really carrying the present
> burden, and which files are merely nearby futures, parked experiments, or
> misleading aggregates?

That distinction matters because this repo contains several classes of truth at
once:

- the live Compose-first runtime
- side-path Compose futures
- generated aggregate views
- optional or parked stacks
- planning and research pressure that points beyond the live runtime

If those classes are blended together, the tree starts sounding like it offers
more genuine options than it really does.
That is one of the exact wounds this knowledgebase is supposed to stop hiding.

## What this page is and is not allowed to prove

This page is allowed to prove:

- where the priority runtime is authored today
- which Compose fragments are part of the live root include path
- which major burdens still live directly in the root file
- which nearby Compose files are real but not part of the active runtime path
- where readers are most likely to confuse breadth with actual anti-SPOF
  maturity

This page is not allowed to prove:

- that wrong-node routing already works generically
- that fragment partitioning equals shared placement truth
- that edge sophistication equals peer-forward correctness
- that "many Compose files" means the repo already has a healthy option space
- that a side-path fragment has earned promotion just because it is thoughtful

## Strongest honest current answer

The live runtime is still decisively Compose-first.
The root file is still the canonical assembly surface.
The fragment set is broad enough that the repo already feels post-Compose in
its burdens, but not post-Compose in its truth ownership.

That is the real reading.

Not:

- "this is just a simple Compose repo"

and not:

- "this is already a distributed control plane with a few rough edges"

The root runtime is substantial.
The missing middle truth layer is still substantial too.

## Authority ladder

If a contributor is trying to understand how the live runtime is put together,
read in this order:

1. root
   [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. the active include targets named by the root file
3. runtime-boundary pages such as
   [current-compose-runtime.md](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/architecture/current-compose-runtime.md)
   and
   [failure-model-and-maturity.md](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/architecture/failure-model-and-maturity.md)
4. only then side-path Compose files and planning material

If someone starts with `compose/docker-compose.everything.yml`,
`compose/docker-compose.parsed.yml`, or a future-facing fragment before reading
the root file, they are already skipping the repo's main truth surface.

## What the root file still owns directly

The root file is not just a loader.
It still directly owns important runtime language and several important
services.

### Shared network language

The root file defines these shared networks:

- `publicnet`
- `backend`
- `warp-nat-net`

That proves the runtime is already modeling different traffic classes
explicitly.
It does not prove that any receiving node already knows the right service
meaning of an incoming request.

### Shared config and secret language

The root file also carries shared config and secret surfaces such as:

- secret `signing_secret`
- config `watchtower-config.json`
- session-manager assets
- inline Homepage configuration blobs

That proves Compose is being used as a real configuration surface.
It does not prove topology truth, current placement truth, or peer eligibility
truth.

### Root-owned services

The root file still directly defines a serious workload set, including:

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

This matters because the root file is still one of the main places where real
infrastructure decisions are made.
Any future promotion to a narrower controller, helper layer, or orchestrator
must account for the fact that the root file is not residue.
It is still the live platform surface.

## The active include path

The root file currently includes these fragments:

1. `compose/docker-compose.coolify-proxy.yml`
2. `compose/docker-compose.docs.yml`
3. `compose/docker-compose.firecrawl.yml`
4. `compose/docker-compose.headscale.yml`
5. `compose/docker-compose.llm.yml`
6. `compose/docker-compose.metrics.yml`
7. `compose/docker-compose.stremio-group.yml`
8. `compose/docker-compose.warp-nat-routing.yml`
9. `compose/docker-compose.wishlist.yml`

That list is not trivia.
It is the present-tense answer to "which fragments are actually part of the
priority implementation?"

Everything else in `compose/` must be read with more caution.

## Fragment burden ledger

The cleanest way to read the fragment set is not as a directory tour.
It is as a burden ledger: what does each fragment visibly own, and what nearby
claim must still be refused?

| Active fragment | It visibly owns | Evidence examples in the file tree | What the docs must still refuse |
| --- | --- | --- | --- |
| `compose/docker-compose.coolify-proxy.yml` | public ingress, auth, middleware, DDNS, edge security, helper-driven failover generation | `traefik`, `tinyauth`, `crowdsec`, `cloudflare-ddns`, `docker-gen-failover`, `autokuma`, `whoami`, file-provider output to `traefik/dynamic/failover-fallbacks.yaml` | that sophisticated ingress already equals generic wrong-node dignity or protected-route continuity |
| `compose/docker-compose.docs.yml` | docs serving inside the same routed platform | `mkdocs` is part of the active include path | that the documentation surface escapes the same ingress truth limits it is describing |
| `compose/docker-compose.firecrawl.yml` | browser automation, queue-backed workers, support state | `playwright-service`, `firecrawl`, `nuq-postgres`, `rabbitmq` | that worker and queue presence somehow upgrades stateful failover truth |
| `compose/docker-compose.headscale.yml` | private mesh identity and connectivity assumptions | `headscale-server`, `headscale`, active SQLite config at `/var/lib/headscale/db.sqlite` | that reachability or node identity equals valid forwarding truth |
| `compose/docker-compose.llm.yml` | AI gateway, tool exposure, caches, auxiliary state | `open-webui`, `litellm`, `litellm-postgres`, `qdrant`, `mcp-proxy`, `mcpo`, `gptr` | that app-level gateway failover solves node-level request preservation |
| `compose/docker-compose.metrics.yml` | metrics, logs, dashboards, probes, alerts | `prometheus`, `victoriametrics`, `grafana`, `loki`, `alertmanager`, exporters, blackbox probes, dashboards for `docker-gen-failover`, `tinyauth`, `whoami`, `headscale`, `watchtower`, `mongodb`, `redis` | that observability breadth equals automated recovery or shared runtime truth |
| `compose/docker-compose.stremio-group.yml` | media/search/debrid/helper workload cluster | heterogeneous app bundle with proxying, storage, and external-provider dependence | that workload breadth creates unified placement truth or resilience semantics |
| `compose/docker-compose.warp-nat-routing.yml` | alternate egress shaping and network behavior control | `warp-net-init`, `warp-nat-gateway`, `warp_router`, `ip-checker-warp`, inline scripts | that packet-path cleverness equals preserved application meaning |
| `compose/docker-compose.wishlist.yml` | a small narrow app surface | active lightweight app include | that a simple stateless app surface proves the hard routes are solved |

## What the active fragment set actually says about the repo

The fragment set already tells a more specific story than "large Compose
repo."

It says the priority runtime is simultaneously carrying:

- edge and identity pressure
- private mesh pressure
- observability pressure
- state and worker pressure
- media/app sprawl pressure
- network behavior pressure

That is why the repo no longer feels like simple container hosting.
The operator is already being asked to reconcile the kinds of burdens that
normally force a system to either:

- expose shared truth clearly

or:

- quietly make one human keep stitching the truth together by hand

The current tree proves the burden set.
It does not yet prove that the burden transfer has happened.

## Root versus fragment: where the pressure really lives

One easy mistake is to imagine the root file as a bootstrap layer and the
fragments as the "real stack."
That is not how this repo reads honestly.

The root file still owns:

- shared networks
- shared config and secret naming
- direct TCP-exposed stateful services such as `mongodb` and `redis`
- operator-facing utilities such as `code-server`, `dozzle`, and `portainer`
- the docs and site-adjacent world around `bolabaden-nextjs`, `homepage`, and
  `telemetry-auth`

The fragments then extend that root with additional domains.

So the better reading is:

- the root file is still a major workload and policy surface
- the fragments widen the domain footprint
- the repo's missing truth layer has to reconcile both, not replace one neat
  bucket

## Important side-path Compose files that are not in the active include path

The repo also contains meaningful Compose files that are not currently part of
the live root include graph.
They matter, but they do not carry the same authority.

| Side-path file | Why it matters | Why it is dangerous if overread |
| --- | --- | --- |
| `compose/docker-compose.authentik.yml` | preserves a real alternate identity direction | it can be mistaken for present auth authority when TinyAuth is the visible live edge path |
| `compose/docker-compose.l4-ingress.yml` | preserves separate thinking for L4 and raw TCP ingress | it can be confused with active runtime proof for TCP failover |
| `compose/docker-compose.nomad.yml` | proves scheduler exploration is real, not imaginary | it can be mistaken for an already-earned control-plane decision |
| `compose/docker-compose.core.yml` | captures alternate service grouping instincts | it can be mistaken for the present canonical assembly path |
| `compose/docker-compose.warp.yml` and `compose/docker-compose.vpn-docker.yml` | preserve alternate network and egress approaches | they can be confused with the active WARP routing path |
| `compose/docker-compose.parsed.yml`, `compose/docker-compose.semiparsed.yml`, `compose/docker-compose.semifullparsed.yml`, `compose/docker-compose.everything.yml` | show aggregate or transformed views of the stack | they can look more "complete" than the root file and tempt readers into treating generated output as canonical truth |
| `compose/docker-compose.unused.yml`, `compose/docker-compose.unsend.yml`, `compose/docker-compose.wordpress.yml`, `compose/docker-compose.plex.yml` | show template appetite, parked domains, and alternate futures | they can inflate the perceived option space without earning present authority |

## The most dangerous mistakes this page is trying to stop

### Mistake 1: file count as option-space proof

"There are many Compose files" does not mean "the user now has many mature
choices."

Some files are:

- live runtime surfaces
- exploratory futures
- generated views
- historical parking lots

Those are not interchangeable forms of maturity.

### Mistake 2: modularity as resilience

A runtime partitioned into neat fragments can still leave the operator owning
the most important answer:

> what should happen to this request on this node right now?

Fragment cleanliness helps orientation.
It does not by itself transfer that burden into the system.

### Mistake 3: aggregate output as authority

Derived files can feel more total and therefore more truthful.
In this repo, that is a trap.
The root file plus its active include path is the authority surface.
Generated aggregates are downstream views, not upstream truth.

### Mistake 4: side-path exploration as current behavior

The repo is full of real thought.
Real thought is not the same thing as promoted runtime behavior.

That distinction matters especially in a repo whose whole theme is:

- no fake closure
- no fake HA
- no fake answer just because it sounds more adult

## What this page can honestly support elsewhere

This page can support sentences like:

- "the active runtime still centers on root Compose"
- "the root file directly owns both shared infrastructure language and major
  services"
- "the current runtime is widened by nine active fragments"
- "the tree contains active, side-path, generated, and parked Compose files
  that must not be flattened into one authority class"

This page cannot support sentences like:

- "the repo already has a coherent distributed control plane"
- "the fragment graph already proves anti-SPOF behavior"
- "the generated aggregate files are a better truth source than the root file"
- "the side-path fragments mean the platform already has many mature fallback
  architectures"

## Bottom line

The fragment map proves that the priority implementation is still one
Compose-first runtime assembled from a root file plus a live include set.
It also proves that the tree contains several neighboring futures, helper
directions, and misleading aggregates that can easily be mistaken for present
authority.

That is the real documentary burden here.
The problem is not merely "too many files."
The problem is that the tree can be narrated as richer, more converged, and
more option-complete than it really is unless the docs keep saying:

- this is live
- this is nearby
- this is generated
- this is parked
- and none of those categories should be promoted by mood alone
