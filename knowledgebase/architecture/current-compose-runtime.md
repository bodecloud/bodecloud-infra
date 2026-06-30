# Current Compose Runtime

This page is the strongest live implementation surface in the repo:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- the fragments it includes directly under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

If a claim cannot survive contact with those files, it is not a live-runtime
claim yet.

## What this page is and is not allowed to prove

This page is authoritative about:

- what the priority runtime includes today
- which major service domains are really active in the root stack
- which helper components and failure signatures are already visible

This page is not authoritative about:

- generic wrong-node success
- end-to-end failover proof
- stateful correctness under peer substitution
- promoting future helper layers into current truth

This page is the strongest current inventory, not the final resilience verdict.

## Strongest honest current answer

The root runtime is already broad, serious, and operationally rich.

It clearly proves:

- Compose is still the real implementation center
- the stack already has a substantial public edge layer
- observability, auth, mesh, worker, and app domains are real
- the repo is far beyond "just some containers"

It still does **not** prove the hardest missing truths:

- shared placement truth
- durable wrong-node rescue routes
- peer eligibility truth
- policy-preserving peer handoff
- honest stateful failover

That is why the stack can look mature while the user's real complaint remains
valid.

## Root-owned networks

The root file defines three central networks:

| Network | Visible role | What it suggests | What it does not prove |
| --- | --- | --- | --- |
| `publicnet` | ingress-adjacent and externally exposed traffic | the repo is already partitioning public-facing traffic explicitly | that the same topology truth is shared across nodes |
| `backend` | east-west application and support traffic | the stack already distinguishes internal service traffic from public edge traffic | that cross-node service placement is coordinated |
| `warp-nat-net` | specialized egress / routing surface | the repo has real network-policy and alternate-routing pressure, not just app hosting | that route correctness survives node loss or helper drift |

These networks prove segmentation pressure.
They do not prove shared cluster truth.

## Active include graph

The live root stack includes these fragments directly:

| Fragment | Primary domain | Main components or concerns |
| --- | --- | --- |
| `compose/docker-compose.coolify-proxy.yml` | public edge and ingress | Traefik, CrowdSec, TinyAuth, `nginx-traefik-extensions`, `cloudflare-ddns`, `docker-gen-failover` |
| `compose/docker-compose.docs.yml` | docs serving | keeps the documentation inside the same routed surface it is describing |
| `compose/docker-compose.firecrawl.yml` | crawl / browser / queue workloads | worker-style services with their own support-state concerns |
| `compose/docker-compose.headscale.yml` | private mesh and coordination assumptions | `headscale-server`, `headscale`, mesh identity, ACL, internal routing assumptions |
| `compose/docker-compose.llm.yml` | model gateway and AI tooling | LLM surfaces, helper databases, MCP-style tooling pressure |
| `compose/docker-compose.metrics.yml` | observability | Grafana, Prometheus, VictoriaMetrics, Loki, Promtail, cAdvisor, Blackbox, Alertmanager |
| `compose/docker-compose.stremio-group.yml` | media and adjacent workloads | large heterogeneous app cluster with external dependency churn |
| `compose/docker-compose.warp-nat-routing.yml` | route experiments and controlled egress | alternate network behavior and privilege-sensitive routing |
| `compose/docker-compose.wishlist.yml` | smaller standalone services | simple app absorption without inventing a new domain model |

This include set is important because it shows the repo's real operating
domains:

- public ingress
- private mesh
- observability
- workload clusters
- network experimentation

That is already a platform-shaped stack, even though the missing shared-truth
layer is still incomplete.

## Root-owned service classes

The root file is not just glue. It still owns many direct services.

Representative root-owned services include:

- data and state: `mongodb`, `redis`
- operator tools: `dozzle`, `portainer`, `homepage`, `watchtower`
- app and utility surfaces: `searxng`, `code-server`, `dns-server`,
  `telemetry-auth`, `chat-analytics`
- site-facing workloads: `bolabaden-nextjs`
- project-specific TCP and mixed-protocol services:
  `biodecompwarehouse*`

This matters because any future promotion cannot honestly treat the root file as
legacy residue. It is still a major infrastructure surface.

## The live edge stack already present

The active proxy fragment makes the current edge story concrete.

Visible live edge components include:

- `traefik`
- `crowdsec`
- `crowdsec-init`
- `tinyauth`
- `nginx-traefik-extensions`
- `cloudflare-ddns`
- `docker-gen-failover`

That proves the repo is already serious about:

- ingress
- auth continuity
- middleware and security policy
- public DNS participation
- dynamic route generation attempts

It does **not** prove that those components already compose into durable
wrong-node recovery.

## The live mesh and observability surfaces

The runtime also clearly includes:

- `headscale-server` and `headscale` in the Headscale fragment
- `grafana`, `prometheus`, `victoriametrics`, `loki`, `promtail`,
  `cadvisor`, `blackbox-exporter`, and `alertmanager` in the metrics fragment

This matters for two reasons.

First, the repo is not missing seriousness or instrumentation.

Second, even with those surfaces present, the repo still does not automatically
own:

- current placement truth
- semantically valid peer substitution
- request-preserving failover

Observability helps measure the platform. It does not replace the missing
middle layer.

## Concrete gaps already named by repo planning

The runtime should be read alongside the master-plan gaps the repo has already
named explicitly.

Current planning material calls out:

- `docker-gen-failover` can delete routes when containers stop
- `watchtower` is configured but not functioning correctly
- Cloudflare DDNS multi-record behavior is still problematic
- secret sync is still manual
- compose sync is still manual
- service failover is not automated

Those named gaps are a useful reminder:

the runtime already has many ingredients, but the repo itself still records the
core control-plane truths as incomplete.

## What the current runtime really proves

The current root runtime proves all of these:

1. the stack is Compose-first in reality, not just in rhetoric
2. the stack already spans edge, mesh, observability, worker, and app domains
3. the platform already has enough machinery that loose language becomes
   dangerous
4. the repo is already trying to solve the multi-node problem from inside a
   serious implementation, not from a toy baseline

## What it still does not prove

The current root runtime does **not** yet prove:

1. that any healthy public node can preserve a request for a remote service
2. that a live root placement registry such as `services.yaml` exists and is
   consumed
3. that the wrong-node rescue route survives the exact backend failure that
   made rescue necessary
4. that auth and middleware remain semantically identical after peer handoff
5. that stateful services have cluster-grade authority and failover semantics

That difference is the key reading rule for this whole site:

the runtime is already real, but the shared-truth layer is still incomplete.
