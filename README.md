# bolabaden-infra

`bolabaden-infra` is a Compose-first infrastructure repo for `bolabaden.org`.

The core ambition is not "run some containers." It is:

> make multiple ordinary Docker nodes behave like a resilient, operator-readable system without immediately paying the full tax of Swarm, Kubernetes, or another heavyweight orchestrator.

There is a second sentence that matters just as much:

> stop forcing the operator's private memory to act as the missing control
> plane when traffic lands on the wrong node or when the expected backend
> disappears.

That ambition is why this repo keeps returning to the same pressure points:

- no single node should be a mandatory public entrypoint
- requests should prefer local service instances when possible
- a node that receives traffic should be able to forward intelligently when the target service is elsewhere
- health, middleware, auth, and observability should survive failover instead of silently changing behavior
- stateful services should be treated honestly rather than declared "HA" by marketing vocabulary alone
- the control plane should only become more complicated when the extra machinery solves a real pain that Compose alone cannot

This README also has to preserve a blunt fact that the wider ecosystem keeps
trying to blur:

the repo is not short on product names, proxy names, or orchestration names.
It is short on **real options** that still feel real after the first request
lands on the wrong machine.

That is why the docs here keep sounding harsher than ordinary self-hosting
documentation.
The user frustration is not "there are no tools."
It is "too many supposed solutions stop one layer before truth actually moves
out of my head."

## Read this repo with the right negative benchmark

The user is not mainly asking:

- how do I make the docs cleaner?
- which orchestrator is best?
- what is the modern HA stack?

The user is asking a harsher question:

> which options are real once traffic lands on the wrong node, and which
> options still secretly depend on me remembering the real topology?

That negative benchmark should shape every architectural reading of this repo.

If a solution still requires the operator to privately remember:

- which node really hosts the service
- which peer is safe now
- which route survives backend loss
- which database path still owns truth

then for this repo it is still much closer to a fake option than a solved
platform capability.

## The shortest honest description

Today this repo is a serious, modular, multi-service Docker Compose stack with:

- a real root [`docker-compose.yml`](docker-compose.yml)
- a large set of included fragments under [`compose/`](compose)
- a real Traefik-centered ingress surface
- real auth, middleware, observability, and maintenance components
- clear evidence that the repo wants multi-node current-state routing and failover

It is **not yet** a finished, proven multi-node control plane.

The gap matters. The repo clearly wants node-aware routing, fallback, and anti-SPOF behavior, but current evidence still shows missing or incomplete pieces around:

- live placement truth
- trustworthy failover generation
- cross-node convergence of secrets and environment
- proven peer-aware fallback semantics
- stateful failover that preserves correctness, not just liveness

That list is the difference between:

- a stack that can be described as distributed
- and a stack that stops depending on sacred-node memory

Those are not the same achievement.

If you read nothing else, read that distinction correctly:

- the dream is clear
- the direction is serious
- the proof is still partial

## The fastest way to misread this repo

The fastest bad reading is:

1. see Cloudflare multi-node entry
2. see Traefik, auth, middleware, and observability
3. see orchestration side paths like OpenSVC, Nomad, or k3s
4. assume the remaining problem is mainly polish, automation, or platform
   choice

That reading is wrong.

The remaining problem is still missing truth ownership:

- current placement truth
- peer eligibility truth
- convergence truth
- route persistence under the relevant failure
- stateful correctness truth

That is why the repo can look sophisticated and still feel unsolved in exactly
the way the user keeps complaining about.

## Which repo files actually explain that dream

This repo has several instruction surfaces, but they are not equal.

If you are trying to understand the multi-node Docker, no-Swarm,
local-first-then-peer-forward idea, the order is:

1. [`.github/copilot-instructions.md`](.github/copilot-instructions.md)
2. [`README.md`](README.md)
3. [`AGENTS.md`](AGENTS.md)
4. [`.cursorrules`](.cursorrules)

The blunt summary is:

- `copilot-instructions.md` names the architecture dream directly
- `README.md` keeps the repo-level honesty wall around that dream
- `AGENTS.md` tells you where current implementation truth must be checked
- `.cursorrules` mostly governs service-authoring discipline, not distributed
  semantics

That distinction matters because one of the easiest ways to misunderstand this
repo is to treat repeated wording across several files as if it were proof that
the runtime already behaves the way the dream is described.

## What the repo is really trying to build

The strongest repo-native statement of intent is [`.github/copilot-instructions.md`](.github/copilot-instructions.md).

That file explicitly describes:

- multi-node Docker without Kubernetes or Swarm by default
- no central orchestrator
- distributed failover
- a lightweight `services.yaml` current-state registry
- L7 routing through Traefik
- separate L4 handling for raw TCP workloads
- Cloudflare multi-A DNS for node-level failover

Its request model is simple and important:

```text
User -> Cloudflare DNS -> any node
  local service exists -> serve locally
  service is remote    -> forward to peer that has it
```

That is the architectural dream.

The repo is not only trying to host services. It is trying to answer a more specific frustration:

> why does getting redundancy usually force operators to choose between raw Compose sprawl and a giant orchestration platform they do not actually want?

That question should be read a little more aggressively than it first sounds.

The complaint is not just that the market offers bad defaults.
It is that many market answers improve one narrow slice while still leaving the
operator to reconstruct:

- where the service really lives
- whether the wrong node knows that
- whether the fallback path survives the relevant failure
- whether auth and middleware continuity remain intact after handoff
- whether a stateful backend is honestly resilient or merely still reachable

That is the real option drought this repo is reacting to.

## The repo's real question is smaller than Kubernetes and bigger than Compose

This is one of the most important framing rules in the tree.

The repo is not asking:

- should we stay on raw Compose forever?
- should we just bite the bullet and move to Kubernetes?

It is asking:

> what is the smallest added truth-owning layer that makes wrong-node requests,
> backend-loss recovery, and hidden human topology memory stop being the
> dominant failure mode?

That is why this repo keeps circling helper agents, `services.yaml`,
OpenSVC-shaped ingress work, schema-first ideas, Nomad comparisons, and k3s
exploration at the same time.

Those are not random side quests.
They are different attempts to locate the missing middle layer without paying
 full scheduler tax before it is clearly earned.

## What is live vs what is planned

This repo now treats documentation in three layers because flattening them together is what made earlier docs ambiguous and misleading.

### 1. Live implementation truth

This is what is actually present in the worktree and the merged root Compose surface.

Examples:

- root [`docker-compose.yml`](docker-compose.yml)
- included fragments under [`compose/`](compose)
- services visible via `docker compose config`
- live Traefik labels, networks, secrets, and service definitions

### 2. Planned architecture truth

This is where the repo is clearly trying to go, based on planning and design docs.

Examples:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`docs/stateful_ha_plan.md`](docs/stateful_ha_plan.md)
- [`docs/osvc_ingress_ha.md`](docs/osvc_ingress_ha.md)

### 3. Research-pressure truth

This is the archive of repeated questions, experiments, comparisons, and frustrations that explain why the repo keeps exploring Compose, OpenSVC, Nomad, k3s, Cloudflare, Traefik, L4/L7 failover, and anti-SPOF patterns at the same time.

Examples:

- [`knowledgebase/source-archive/`](knowledgebase/source-archive)
- the synthesized research pages under [`knowledgebase/research/`](knowledgebase/research)

These layers are related, but they are not interchangeable. A planned control surface is not live proof. A research thread is not a shipped implementation. A parsed Compose graph is not resilience.

That separation is not just documentation hygiene.
It is the main defense against polished ambiguity.

This repo has enough material to sound deeply understood while still silently
shrinking the user's dream into something more normal, such as:

- better HA
- better routing
- better service discovery
- better orchestrator comparison

Those are all subproblems.
They are not the whole ask.

## The current architecture shape

At a high level, the root stack is built around these surfaces:

- public ingress and edge controls
  - Traefik
  - TinyAuth
  - CrowdSec
  - nginx-based request extensions
  - Cloudflare/DDNS integration
- observability
  - Prometheus
  - VictoriaMetrics
  - Grafana
  - Loki
  - Alertmanager
- app and media services
  - site, docs, dashboards, APIs, media and support services
- state-bearing services
  - Redis
  - MongoDB
  - RabbitMQ
  - Postgres variants
- network and egress experiments
  - WARP-related routing components

The key point is that the root stack already has a broad and real runtime surface.

The harder point is that broad runtime surface does **not** automatically imply:

- trustworthy node-aware forwarding
- trustworthy failover preservation
- trustworthy route persistence under backend disappearance
- trustworthy stateful continuity

## Known architectural pressure points

The repo's current decision surface is dominated by a handful of real issues.

### Missing tracked placement truth

The repo repeatedly describes a lightweight `services.yaml` registry, but the tracked root implementation does not currently ship a live root `services.yaml`.

That means the idea is central, but the live source of truth is still incomplete.

### Failover generation is present but not proven trustworthy

`docker-gen-failover` exists in the proxy layer, but repo planning explicitly records that it can remove routes when containers stop. That is the opposite of what a failover mechanism must do.

### DNS is only the first layer

Cloudflare multi-A DNS can help clients hit a surviving node, but it does not prove that the surviving node:

- knows where the target service lives
- forwards correctly
- preserves middleware/auth semantics
- preserves stateful correctness

### HTTP and TCP are different problems

HTTP routing through Traefik is far easier to reason about than TCP failover for Redis, MongoDB, and other state-sensitive services. The repo has to keep those classes separate or the docs become dishonest.

### Stateful HA is the real honesty wall

Ingress cleverness can hide a lot. It cannot fake correct replication, election, quorum, reconnect behavior, or durable failover for stateful systems.

This is also why the repo keeps feeling bigger than "Docker failover."

The deeper problem is not lack of components.
It is lack of a shared truth surface that can make the first healthy receiving
node stop being a semantic gamble.

## Real options vs fake options

This repo needs a blunter filter than normal infrastructure READMEs use.

| Option class | Real only if... | Still fake if... |
| --- | --- | --- |
| Multi-node DNS entry | the receiving node can still preserve the request meaningfully | the request only works when it lands on the right node by luck |
| Proxy-centered failover | fallback survives the failure that made fallback necessary | the route vanishes when the preferred backend disappears |
| Placement registry | routing or eligibility logic actually consumes live placement truth | the registry is only architecture intent or planning rhetoric |
| Heavier orchestrator promotion | it removes a named truth gap the current layer cannot own honestly | it mostly renames the same burden while adding worldview tax |
| Stateful HA story | authority, replication, promotion, and client behavior are defined and proven | a service is merely reachable through a stable hostname or TCP path |

The user is frustrated precisely because the ecosystem offers many answers that
look like the left column while behaving like the right column.

## Recommended reading

If you want the real operator-grade explanation instead of the old flattened story, start in the knowledgebase:

- [`knowledgebase/index.md`](knowledgebase/index.md)
- [`knowledgebase/architecture/problem-and-goals.md`](knowledgebase/architecture/problem-and-goals.md)
- [`knowledgebase/architecture/current-compose-runtime.md`](knowledgebase/architecture/current-compose-runtime.md)
- [`knowledgebase/architecture/compose-first-architecture.md`](knowledgebase/architecture/compose-first-architecture.md)
- [`knowledgebase/architecture/ha-failover-routing.md`](knowledgebase/architecture/ha-failover-routing.md)
- [`knowledgebase/architecture/stateful-ha-and-data.md`](knowledgebase/architecture/stateful-ha-and-data.md)
- [`knowledgebase/architecture/capability-gaps-and-roadmap.md`](knowledgebase/architecture/capability-gaps-and-roadmap.md)
- [`knowledgebase/operations/devops-runbook.md`](knowledgebase/operations/devops-runbook.md)

If you want the shortest route through the real problem, use this order:

1. [`knowledgebase/research/user-intent-and-dream.md`](knowledgebase/research/user-intent-and-dream.md)
2. [`knowledgebase/operations/operator-questions-and-honest-answers.md`](knowledgebase/operations/operator-questions-and-honest-answers.md)
3. [`knowledgebase/architecture/request-path-and-failure-walkthrough.md`](knowledgebase/architecture/request-path-and-failure-walkthrough.md)
4. [`knowledgebase/architecture/current-compose-runtime.md`](knowledgebase/architecture/current-compose-runtime.md)
5. [`knowledgebase/architecture/missing-middle-layer.md`](knowledgebase/architecture/missing-middle-layer.md)
6. [`knowledgebase/operations/decision-paths-and-promotion-rules.md`](knowledgebase/operations/decision-paths-and-promotion-rules.md)
7. [`knowledgebase/research/evidence-ledger.md`](knowledgebase/research/evidence-ledger.md)

That order forces the dream, the wound, the live runtime, the missing layer,
and the proof ceiling to stay in view at the same time.

Those pages are the current authoritative explanation because they explicitly separate:

- what exists now
- what the repo is aiming for
- what remains unproven

## Validation commands

This repo is not a single-app project, so "it builds" is not enough.

Use these commands as the minimum baseline:

```bash
docker compose config --quiet
docker compose config --services
python3 -m mkdocs build -f mkdocs.yml --strict
```

For infra tooling under `infra/`:

```bash
cd infra
make build
make test
go vet ./...
```

Important environment assumptions:

- many env vars are required for Compose interpolation
- placeholder secret files under `${SECRETS_PATH}` may be required
- `~/.docker/config.json` must exist, even if it is just `{}`
- Go 1.24+ is required for the `infra/` module

## What this README is intentionally not doing

It is not trying to present the repo as already finished.

It is also not trying to bury the dream under generic best-practice language.

The right way to read `bolabaden-infra` is:

- the user wants genuine multi-node flexibility and resilience
- they do not want to be bullied into heavyweight orchestration before it earns its keep
- they are trying to discover the smallest control surface that removes fake HA and obvious SPOFs without destroying the directness of Compose

That is the problem this repo exists to solve.

## Bottom line

If you are contributing here, the most important discipline is not "use the right YAML syntax."

It is this:

> do not widen the architecture claim beyond the evidence.

And do not narrow the user's dream into a cleaner neighboring question just
because that cleaner question is easier to summarize.

Keep the dream explicit, keep the proof boundaries honest, and make each change answer a real operator problem instead of adding abstract infrastructure theater.
