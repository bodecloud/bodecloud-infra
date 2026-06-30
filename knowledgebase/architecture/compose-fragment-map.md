# Compose Fragment Map

This page explains how the priority implementation is partitioned **right now**.

It is not a directory tour.
It is a control-surface map for the real runtime the user keeps pointing back
to:

- the root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- the included `compose/` fragments
- the surrounding side paths that are visible in the repo but not equal in
  authority

This matters because the tree contains several different kinds of truth at the
same time:

- the live Compose-first stack
- alternate or future directions
- generated snapshot files
- parked or optional service surfaces

If the docs flatten those into one cheerful blob called "the stack," the reader
cannot tell what actually governs the runtime today.
This page is therefore doing more than mapping files.
It is reconstructing a tree that already contains several competing futures at
once, without pretending those futures have already reconciled into one runtime
truth.

That matters because one of the user's deepest frustrations is not just
complexity.
It is the feeling that the ecosystem keeps pretending the option space is
clearer than it is.

This tree is one of the places where that false clarity can easily be
manufactured:

- live runtime files
- alternate futures
- generated aggregates
- parked domains

can all be narrated as if they are one coherent stack unless the docs keep the
authority boundaries explicit.

## The live entrypoint

The live entrypoint is the root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

That file is not just a stub that points elsewhere. It:

- defines the base networks
- declares shared configs and secrets references
- owns several core services directly
- includes the active fragments that extend the runtime by domain

That means any architecture reading that starts from a derived mega-file,
semi-parsed file, or side-path fragment before starting from root Compose is
already reading the repo in the wrong order.

## What the root file tells us structurally

The root file already encodes one of the repo's strongest truths:

- this is still a Compose-first system
- but it is no longer a small, one-surface Compose system

The root owns shared network language such as:

- `publicnet`
- `backend`
- `warp-nat-net`

and then pulls in the active fragment set that creates most of the repo's real
architecture pressure.

That is the important nuance.
The root file is still the contract, but the contract is now broad enough that
it behaves like several domains glued together under one authoring surface.

That is also why the repo feels like it is already post-Compose in its
problems even while Compose remains central in authorship.
The user is not confused about where the YAML lives.
They are frustrated that the surrounding truth still does not compose into one
reliable request-preserving platform.

## Active fragments included by the root stack

These fragments are part of the current live merged stack:

| Fragment | Primary responsibility | Why it exists | Main risk surface |
|---|---|---|---|
| `compose/docker-compose.coolify-proxy.yml` | Public ingress, auth, middleware, edge security | Concentrates external exposure, identity edges, and route-generation pressure. | Routing bugs, auth ordering, failover generation, public exposure, false confidence about dynamic HA |
| `compose/docker-compose.docs.yml` | Documentation serving | Keeps docs inside the same routed infrastructure they are describing. | The docs inherit the same ingress truth gaps they are documenting |
| `compose/docker-compose.firecrawl.yml` | Crawling, browser automation, and queue-backed work | Concentrates worker coordination plus state-bearing support services. | Worker orchestration, Postgres/Rabbit state, restart semantics |
| `compose/docker-compose.headscale.yml` | Private mesh and identity backbone | Provides private node-to-node path assumptions the wider multi-node dream depends on. | Hidden SPOF risk if treated as coordination truth without redundancy |
| `compose/docker-compose.llm.yml` | LLM gateway, model surfaces, MCP-style tooling | Concentrates model-facing tooling, auxiliary datastores, and API aggregation. | Secret sprawl, drift, auxiliary state, API dependency layering |
| `compose/docker-compose.metrics.yml` | Metrics, logs, alerting, dashboards | Gives the repo visibility into itself. | Operational complexity and the risk of using observability breadth as a substitute for recovery truth |
| `compose/docker-compose.stremio-group.yml` | Media, debrid, search, and related apps | Represents one of the largest heterogeneous workload clusters in the repo. | External dependency churn, storage sprawl, network quirks, app drift |
| `compose/docker-compose.warp-nat-routing.yml` | Alternate routing and controlled egress experiments | Encodes network-level ambitions that exceed normal app hosting. | Route correctness, privilege sensitivity, hidden network complexity |
| `compose/docker-compose.wishlist.yml` | Small standalone app include | Shows the root can still absorb simple app surfaces without forcing a new domain model. | Low direct risk, but still part of naming and routing complexity |

## What this fragment set means

The root include set is already telling the story of the repo's real operating
domains:

- edge control
- private coordination
- observability
- app/workload clusters
- network experiments

That is not mere organization.
It is where the user’s real frustrations are already materialized:

- too many important things happening at the edge
- too many services depending on one Compose worldview
- too many futures visible at once
- too much pressure to invent a control plane without admitting it yet

That last line should stay sharp.
One of the most important things this page can help a contributor see is that
the repo is not simply growing "more Compose."
It is accumulating exactly the kinds of burdens that make a missing middle
truth layer unavoidable to even talk about honestly.

## Services defined directly in the root file

The root file also owns a substantial set of services directly rather than pushing everything into fragments.

Representative root-owned services include:

- `mongodb`
- `redis`
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
- several `biodecompwarehouse*` services

## Why this matters

The root file is not just glue.
It is still a major authoring surface for:

- important datastores
- dashboards
- socket mediation
- update automation
- custom project services
- operator-facing utilities

That means any future promotion to a stronger control layer cannot honestly
treat the root file as residue.
It contains real infrastructure decisions today.

## The architectural reading hidden in the fragment map

The fragment layout reveals at least five meaningful domains.
Those domains are the reason the repo no longer feels like "just Docker
Compose," even though Docker Compose is still the live baseline.
The interesting part is not merely that the domains exist.
It is that they create pressure for a stronger control layer while still
refusing to agree on which stronger layer has actually earned promotion.

## 1. Edge control plane

Dominant surfaces:

- Traefik
- CrowdSec
- TinyAuth
- nginx-based extensions
- DDNS
- route generation

Why this domain matters:

- it controls what the outside world actually experiences
- it is where wrong-node success either becomes real or collapses into theater
- it is also where the repo carries one of its clearest known integrity failures: route generation that is not yet trustworthy under backend stop conditions

## 2. Private coordination plane

Dominant surfaces:

- Headscale
- DNS and mesh-related services
- WARP routing helpers

Why this domain matters:

- the multi-node dream depends on trustworthy node-to-node private paths
- but these services do not automatically equal placement truth, convergence truth, or failover correctness

This is a good example of why "there is a mesh" is not enough evidence by
itself.

## 3. Observability plane

Dominant surfaces:

- Prometheus
- VictoriaMetrics
- Grafana
- Loki
- exporters
- Alertmanager

Why this domain matters:

- the repo has already invested significantly in seeing itself
- but observability breadth does not equal automated recovery truth

This is another place where the docs must refuse false inflation.
A visible failure is still a failure until recovery semantics are real.

## 4. Workload clusters

Dominant surfaces:

- media cluster
- llm cluster
- firecrawl/browser-automation cluster
- custom project apps

Why this domain matters:

- this is where heterogeneity becomes scheduling, storage, and drift pressure
- these clusters are also where control-plane promotion decisions are most likely to be earned later

## 5. Network and egress experimentation

Dominant surfaces:

- WARP routing
- alternate route-control helpers
- traffic-path experiments

Why this domain matters:

- the repo's dream is not only "host apps," but also "shape how requests and egress behave across nodes"
- that means network experimentation is not peripheral; it is part of the main anti-SPOF and anti-sacred-node pressure

## Important non-live and side-path Compose files

The tree also contains Compose files that are not part of the active root
include path, but still matter for interpretation.

| File | Why it matters |
|---|---|
| `compose/docker-compose.authentik.yml` | Signals a real alternate identity direction that should not be confused with the current primary auth path. |
| `compose/docker-compose.l4-ingress.yml` | Proves the repo treats plain TCP and L4 handling as a separate domain from HTTP ingress. |
| `compose/docker-compose.nomad.yml` | Shows scheduler/control-plane exploration is not theoretical. |
| `compose/docker-compose.core.yml` | Preserves a simpler or alternate foundational service grouping that still reflects important design thinking. |
| `compose/docker-compose.warp.yml`, `compose/docker-compose.vpn-docker.yml` | Preserve alternate network and egress strategies. |
| `compose/docker-compose.parsed.yml`, `compose/docker-compose.semiparsed.yml`, `compose/docker-compose.semifullparsed.yml`, `compose/docker-compose.everything.yml` | Derived or aggregated views that can look authoritative if read carelessly. Some even carry hardcoded-looking IP-era assumptions. They are not the canonical live entrypoint. |
| `compose/docker-compose.unused.yml`, `compose/docker-compose.unsend.yml`, `compose/docker-compose.wordpress.yml`, `compose/docker-compose.plex.yml` | Parked or optional surfaces that still reveal how broad the repo's template appetite and alternate futures have been. |

## The most dangerous reading mistake

The most dangerous reading mistake is to assume:

- if a Compose file exists, it is part of the live stack
- if a generated aggregate exists, it is the canonical truth
- if a fragment describes a control-plane future, the runtime already behaves that way

This page exists to keep those categories separate.

It also exists to stop a subtler mistake:

- assuming that because several futures are visible, the repo must already
  have a rich and healthy option space

Not necessarily.
Some of those futures are genuinely different.
Some are variations on the same hidden-burden story.
One of the documentation job here is to stop that distinction from being
blurred by file count alone.

## What this page actually proves

The live stack is not confusing merely because there are "too many files."

It is confusing because several futures are visible in the tree at once:

- Compose-first live operation
- exploratory control-plane work
- alternate auth and ingress directions
- generated aggregate snapshots
- optional exhaustive app surfaces

That is why the knowledgebase has to keep separating:

- current live truth
- side-path or research truth
- generated or archival views

Otherwise the repo becomes impossible to reason about under pressure.
Under calm conditions, a messy tree can still feel navigable.
Under failure or migration pressure, confusing current truth with side-path
truth is exactly how the operator gets forced back into hidden reconstruction.

That is especially important here because the user is not asking for generic
tidiness.
The user is asking which surfaces are actually carrying the no-Swarm,
multi-node, anti-SPOF pressure today.

## Bottom line

The fragment map shows that the repo is already partitioned into several real
domains, while still insisting that root Compose remains the canonical runtime
contract.

That is exactly the tension the rest of the documentation has to preserve:

- this is still one Compose-first system
- but it is already carrying enough domain pressure that future control layers,
  promotions, and narrower platforms are visible in the tree

If the docs flatten that into either "just Compose" or "already a distributed platform," they stop telling the truth.
