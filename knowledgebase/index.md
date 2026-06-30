# bolabaden Infrastructure Knowledgebase

This site is for one question:

> how do you keep `docker-compose.yml` as the real center of gravity, spread
> services across multiple ordinary Docker nodes, and still make wrong-node
> traffic, fallback, and anti-SPOF behavior feel like one coherent platform
> instead of one operator remembering the real topology?

That is the real question. The repo is not mainly about:

- general self-hosting
- generic high availability
- "which orchestrator is best?"
- "how do I make Compose prettier?"

Those are adjacent questions. They are not the core problem this repo is trying
to solve.

## What this site is and is not allowed to prove

This site is authoritative about:

- the repo's actual architecture dream
- the current root Compose implementation surface
- the difference between live runtime truth, intent truth, and planning truth
- the concrete gaps between today's stack and real wrong-node recovery

This site is not authoritative about:

- claiming that the current runtime already behaves like the dream
- turning a clean explanation into proof of failover
- promoting research or planning docs into shipped behavior

The site should help a reader leave with the right map of reality, not the most
optimistic story.

## Strongest honest current answer

`bolabaden-infra` already contains a serious Compose-first platform:

- a real root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active includes under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- a substantial Traefik, CrowdSec, TinyAuth, and nginx-auth edge layer
- observability and maintenance surfaces
- private-mesh work through Headscale
- repeated planning pressure toward any-node entry, peer-aware routing, and
  anti-SPOF behavior

What it does **not** yet prove is the part the user actually cares about most:

- that any healthy node can accept a request and preserve it correctly when the
  service is remote
- that placement truth is shared explicitly instead of remembered
- that fallback routes survive the failure they are meant to absorb
- that middleware, auth, and request semantics survive peer handoff
- that stateful services are truly resilient rather than merely reachable

The dream is clear. The stack is real. The missing middle truth layer is still
incomplete.

## Read this site in the right order

If you only have time for one serious pass, use this route:

1. [User Intent and Dream](research/user-intent-and-dream.md)
2. [Problem and Goals](architecture/problem-and-goals.md)
3. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
4. [Current Compose Runtime](architecture/current-compose-runtime.md)
5. [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
6. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
7. [The Missing Middle Layer](architecture/missing-middle-layer.md)
8. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
9. [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

That route preserves:

- what the user wants
- what the repo already is
- what still depends on private operator memory
- what still lacks proof

## Fast site router

Use these shortcuts when you already know the question.

| If you need to know... | Start here |
| --- | --- |
| what the user is actually trying to make true | [User Intent and Dream](research/user-intent-and-dream.md) |
| what the root runtime really contains today | [Current Compose Runtime](architecture/current-compose-runtime.md) |
| which file has the strongest authority for the multi-node no-Swarm dream | [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md) |
| why wrong-node entry is still the humiliating threshold | [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md) |
| why Cloudflare plus Traefik is still weaker than real failover | [HA, Failover, and Routing](architecture/ha-failover-routing.md) |
| what helper/control layer the repo is actually searching for | [The Missing Middle Layer](architecture/missing-middle-layer.md) |
| why Redis, MongoDB, Headscale, and other stateful services need harsher language | [Stateful HA and Data](architecture/stateful-ha-and-data.md) |
| what still needs proof before stronger claims are legal | [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md) |

## The three truth registers that matter most

### 1. Live runtime truth

Use this when the claim is "what is actually implemented right now?"

Primary sources:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active include fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)

### 2. Repo-native intent truth

Use this when the claim is "what is the platform trying to become?"

Primary sources:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

### 3. Planning and pressure truth

Use this when the claim is "what missing layer is being explored or promoted?"

Primary sources:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)
- [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)
- [Stateful HA Evidence](research/stateful-ha-evidence.md)
- [Source Assimilation Index](operations/source-assimilation-index.md)

These registers overlap. They are not interchangeable.

## The easiest bad read

The easiest way to misunderstand this repo is:

1. see Cloudflare and multiple public nodes
2. see Traefik, CrowdSec, TinyAuth, dashboards, helpers, and many services
3. see Nomad, k3s, and OpenSVC exploration
4. conclude that the remaining problem is mostly polish or automation

That reading is wrong.

The remaining problem is still the hard one:

- current placement truth
- peer eligibility truth
- fallback-route durability
- cross-node semantic parity
- stateful honesty
- getting those truths out of one operator's head

This site is only useful if it keeps those gaps visible instead of smoothing
them away.
