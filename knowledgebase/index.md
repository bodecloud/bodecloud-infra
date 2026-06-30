# bolabaden Infrastructure Knowledgebase

This site exists to answer one unusually stubborn infrastructure question:

> how do you keep
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> as the real human control surface, spread services across several ordinary
> Docker nodes, and still make wrong-node traffic, fallback, and anti-SPOF
> behavior feel like one coherent platform instead of one operator privately
> remembering the real answer?

That is the real subject of `bolabaden-infra`.

This site is not primarily:

- a generic self-hosting handbook
- a reverse-proxy cookbook
- an orchestrator comparison catalog
- a generic `modern homelab` notebook

Those topics only matter when they help answer the real question above.

## The accusation this site has to preserve

The site only stays useful if it keeps the user's real accusation visible:

> there seem to be endless options for multi-node Docker, failover, overlays,
> service discovery, ingress, and orchestration, but too many of them solve
> one visible layer and then quietly leave the operator as the hidden control
> plane when reality gets sharp.

If a page becomes more polished while losing that accusation, it got worse.

## The shortest correct reading

The strongest intent surface in the repo is
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It says the intended direction is:

- no central orchestrator by default
- services manually assigned to nodes
- current-state truth preferred over scheduler-declared desired state
- local-first serving when the requested service already lives on the receiving
  node
- peer-forward fallback when the receiving node is healthy but the service is
  remote
- explicit separation between HTTP routing and raw TCP or stateful behavior
- anti-SPOF pressure without fake HA language

That target request contract is:

```text
User -> Cloudflare DNS -> any surviving public node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that really hosts it
```

That is the dream.
It is not the same thing as current proof.

## The three-part checksum for the whole site

Everything in this knowledgebase should keep this sentence intact:

1. the dream is specific
2. the runtime is real
3. the truth-owning middle layer is still incomplete

If a page loses one of those three, it starts lying in one direction or
another.

## What the repo already proves

The priority implementation is still rooted in:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

That runtime already proves a serious Compose-first platform with:

- a real Traefik and auth-bearing edge layer
- public and private network segmentation
- stateful services such as MongoDB, Redis, PostgreSQL, RabbitMQ, and Qdrant
- observability and operator surfaces
- Headscale and mesh pressure
- alternate routing and egress experiments

What it still does **not** prove is the thing the user actually cares about:

- that any healthy public node can accept a request and preserve it correctly
  when the service is remote
- that placement truth is shared explicitly instead of remembered
- that peer eligibility is system-owned rather than guessed
- that fallback paths survive the failure that makes fallback matter
- that auth, middleware, and request semantics survive peer handoff
- that stateful services are genuinely resilient rather than merely reachable

The runtime is serious.
The missing burden transfer is still brutally concrete.

## The hidden operator job the site has to keep naming

The shortest honest summary of the current wound is:

the operator is still acting like the missing control plane.

That hidden job currently includes things like:

- remembering what runs on which node right now
- remembering which peer is actually valid right now
- remembering whether a generated fallback route would still exist under real
  failure
- remembering whether a reachable answer is only transport-reachable or
  semantically valid
- remembering which stateful surfaces still hide a sacred authority node

If a page forgets that hidden job, it will almost always overstate the system.

## The four truth layers

This site only stays honest if it keeps these truth layers separate:

| Truth layer | Main examples | What it is allowed to answer | What it is not allowed to answer |
| --- | --- | --- | --- |
| Live runtime truth | `docker-compose.yml`, active `compose/` fragments, `docker compose config` | what the priority implementation actually ships now | what the system would do if missing layers already existed |
| Repo-native intent | `.github/copilot-instructions.md`, `knowledgebase/AGENTS.md` | what the repo is clearly trying to become | whether the runtime already behaves that way |
| Planning truth | `docs/INFRASTRUCTURE_MASTER_PLAN.md`, stateful and ingress plan docs | what gaps are already named and what promotion paths exist | that a planned repair is already active |
| Archive and research pressure | `knowledgebase/research/`, `knowledgebase/source-archive/` | why the user keeps rejecting partial answers and which patterns recur | present-tense runtime proof |

The site gets misleading the moment those four layers blend into one calm
voice.

## The core architecture problem

The repo is not just asking for `better HA`.
It is asking for a smaller honest middle layer between:

- static multi-node Docker glue that still depends on private folklore
- and a heavyweight orchestrator worldview that has not yet earned the right to
  hide that much truth

That middle layer would need to make at least some of these truths
system-owned:

- current placement
- current peer eligibility
- current rescue-route validity
- current route-class meaning
- current explanation for why the system chose local versus remote

If a candidate does not move those truths into the system, it has not solved
the user's real complaint, no matter how mature it sounds.

## The pressure already visible in the archive

The archive is consistent about where the real pain lives:

- multi-node Docker without Swarm keeps converging on service discovery as the
  hard missing piece
- distributed HA orchestration keeps showing that fully peer-equal behavior is
  rare without custom glue
- load-balancer and failover alternatives keep showing that many products
  solve only one slice of the wound
- Nomad, k3s, Kubernetes, and related explorations show that stronger control
  planes may help, but they bring a larger worldview that still has to prove
  it kills the right private burden

That is why the user does not feel like they lack product names.
They feel like they lack honest options.

## How to use this site

If you want the shortest route through the knowledgebase:

1. [Problem, Pressure, and Goals](architecture/problem-and-goals.md)
2. [Current Compose Runtime](architecture/current-compose-runtime.md)
3. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
4. [The Missing Middle Layer](architecture/missing-middle-layer.md)
5. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
6. [Failure Model and Maturity Matrix](architecture/failure-model-and-maturity.md)
7. [Orchestration Options](architecture/orchestration-options.md)
8. [DevOps Runbook](operations/devops-runbook.md)

If you want the proof pressure underneath those pages:

- [Evidence Ledger](research/evidence-ledger.md)
- [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)
- [Stateful HA Evidence](research/stateful-ha-evidence.md)
- [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)
- [Source Assimilation Index](operations/source-assimilation-index.md)

If you want the navigation logic spelled out first:

- [Reading Paths](reading-paths.md)

## What a good page in this site should leave behind

After reading a useful page here, a contributor should be able to answer:

- what concrete question the page was solving
- which truth layer carried the answer
- which burden still remains operator-owned
- which stronger sentence is still forbidden
- what artifact, drill, or runtime change would make that stronger sentence
  legal

If a page leaves behind only:

- `the architecture is clearer`
- `the options are broader`
- `the platform sounds more mature`

then it is still too weak for this repo.

## What this site must keep refusing

This knowledgebase should keep several false conclusions illegal:

- `there are multiple public nodes, so failover is basically solved`
- `Traefik exists, so wrong-node forwarding must be close`
- `a helper has failover in the name, so the platform owns fallback now`
- `a route can be rendered, so it must survive the failure that made it
  matter`
- `the docs are clearer now, so the platform must be closer to adulthood`

Those are exactly the overreads this site exists to interrupt.

The useful conclusion is not that the repo is simple.
It is that the repo is serious, the dream is specific, and the missing burden
transfer is still brutally concrete.
