# Instruction Surfaces and Authority

This page exists because the repo now has enough aligned language to create a
new kind of confusion:

- the dream is clear
- the warnings are clear
- the authoring rules are clear
- therefore the stack must be close to solved

That conclusion is false.

The user is not checking whether the repo sounds thoughtful.
The user is checking whether the repo still knows which files are allowed to
state:

- the dream
- the honesty boundaries
- the current runtime
- the authoring rules
- the pressure from research and prior frustration

If those truth classes blur together, the docs become smoother and less useful.

## The exact question this page answers

This page answers:

> which files actually explain the multiple ordinary Docker nodes,
> no-Swarm-by-default, local-first-then-peer-forward dream, and which files
> only constrain how we talk about or author that dream?

The neighboring smaller question it must not collapse into is:

> which files mention HA, failover, or Compose the most?

That smaller question is too weak.
This repo is full of files that mention related concepts without owning the
same truth.

## What this page is and is not allowed to prove

This page is allowed to prove:

- which files own which truth class
- which files are primary versus supporting for the multi-node Compose dream
- which file should win when two truthful files are talking about different
  layers
- which kinds of claims are forbidden from each surface

This page is not allowed to prove:

- that the runtime already satisfies the dream
- that repeated language across files counts as implementation corroboration
- that a planning artifact is live proof
- that a stricter writing standard means the system itself owns more truth

## Authority map

The repo's most important authority surfaces are not equal.

| Surface | Truth class | What it is authoritative about | What it must not be upgraded into |
|---|---|---|---|
| [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md) | Architecture-dream truth | The clearest statement of the actual desired operating contract | Runtime proof that wrong-node recovery already works |
| [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md) | Repo-level honesty wall | The repo-wide framing for anti-SPOF, anti-fake-HA, and the user's actual complaint | A live implementation report |
| [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md) | Runtime-anchor and working-surface truth | The priority implementation surface, validation commands, and operator/build constraints | The main architecture manifesto |
| [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules) | Authoring-discipline truth | Commit discipline, inline config preferences, healthcheck requirements | Distributed-systems correctness or failover semantics |
| [`knowledgebase/AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/AGENTS.md) | Documentation-honesty truth | How docs must separate live truth, planned truth, and research pressure | Proof that the implementation has moved one more burden out of the operator's head |
| [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml) | Live root runtime truth | What the priority implementation actually declares at the root | Generic multi-node correctness |
| [`compose/docker-compose.*.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/) | Live fragment truth | What each active or inactive subsystem currently contributes | That all fragments are equally live or equally trusted |
| [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md) | Planned architecture truth | Which missing truth-owning layers the repo has been explicitly considering | That the root runtime already owns those layers |
| [`knowledgebase/research/`](../research/) and source archive | Research-pressure truth | Why the repo keeps returning to service discovery, failover, routing, anti-SPOF, Nomad, OpenSVC, k3s, and other options | Shipped behavior |

## The strongest honest reading order

If the reader wants to understand the actual dream first, the reading order is:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml) and the active fragments
4. [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
5. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
6. the research and archive pages

If the reader wants to know what is safe to claim about the current runtime,
the order changes:

1. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. active `compose/` fragments
3. `docker compose config`
4. runtime command output and drills
5. only then the planning and research surfaces as contrast

That split matters.

One order answers:

> what is the repo trying to become?

The other answers:

> what has the priority implementation actually earned?

Those are not the same question.

## The clearest dream surface

The clearest architecture-dream file is
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).

It explicitly states:

- no central orchestrator by default
- services manually assigned to nodes
- current-state truth preferred over scheduler-declared desired state
- local-first serving
- peer-forward fallback when the request lands on a healthy node without the
  local service
- explicit separation between L7 HTTP behavior and L4/raw TCP behavior
- anti-SPOF pressure without fake HA language

It also states the target request contract directly:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That file is therefore the strongest answer to:

> what is the user actually trying to build?

It is not the strongest answer to:

> what already works today?

## The repo-level honesty wall

[`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md) is the strongest repo-level honesty wall around
the dream.

Its main job is not to introduce every subsystem.
Its main job is to keep the dream from being normalized into a calmer generic
question like:

- "which orchestrator is best?"
- "how do we improve HA?"
- "how do we modernize the stack?"

The README keeps the harsher framing visible:

> which options are still real once traffic lands on the wrong node and the
> operator is no longer allowed to privately complete the topology story?

That is why the README has more architectural authority than `AGENTS.md`, even
though `AGENTS.md` is closer to the current implementation.

## The runtime-anchor surfaces

[`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md), [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml),
and the active fragments under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/) own the runtime
anchor layer.

Their job is to answer questions like:

- what is the priority implementation surface?
- what services, configs, secrets, and networks actually exist now?
- what commands are required to validate the authored graph?
- what toolchain or environment assumptions exist for working in this repo?

Examples from the current root runtime that belong to this layer:

- the root includes active fragments such as:
  - `docker-compose.coolify-proxy.yml`
  - `docker-compose.headscale.yml`
  - `docker-compose.metrics.yml`
  - `docker-compose.warp-nat-routing.yml`
  - `docker-compose.wishlist.yml`
- the root directly declares live shared networks like:
  - `publicnet`
  - `backend`
  - `warp-nat-net`
- the root directly exposes stateful services like:
  - `mongodb`
  - `redis`
- the root directly wires observable or policy-relevant services like:
  - `watchtower`
  - `homepage`
  - `code-server`
  - `searxng`

Those files tell us the stack is real, broad, and operationally serious.
They do not tell us the wrong-node problem is solved.

## The authoring-discipline surfaces

[`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules) and parts of
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
govern how services should be authored.

They are authoritative about things like:

- commit every change
- prefer inline Compose configs
- require real healthchecks
- do not weaken healthchecks to make the graph look healthy

Those rules matter because bad authoring habits make the anti-SPOF question
harder to answer honestly.

They do **not** prove:

- that peer-forward fallback works
- that middleware survives wrong-node handoff
- that a healthcheck implies cross-node resilience
- that a service with Traefik labels has earned HA language

This is one of the repo's most important authority boundaries.
Good Compose hygiene is necessary.
It is not the same thing as distributed-systems truth.

## The planning surfaces

[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
and related `docs/` planning artifacts are authoritative about pressure and
direction, not completion.

They are useful because they keep naming the same missing burdens:

- secret and env convergence across nodes
- service placement truth
- failover and redeploy strategy
- DNS plurality that does not collapse into one sacred public node
- auth and middleware continuity
- stateful authority honesty

Those docs are allowed to say:

- which missing middle layers the repo has already recognized
- which candidate control surfaces might earn promotion later
- which burdens still live in operator memory

They are not allowed to say:

- that those surfaces already exist in the root implementation
- that a coherent plan is one step away from working runtime custody

## The research-pressure surfaces

The research pages and source archive are not decorative.
They explain why the repo keeps investigating so many adjacent systems without
wanting to surrender immediately to a heavyweight scheduler.

They are authoritative about:

- repeated patterns of user frustration
- ecosystem dead ends
- tradeoff pressure between Compose, Nomad, OpenSVC, k3s, Swarm, Cloudflare,
  and helper layers
- the real accusation behind "there are no real options"

They should not be overread as implementation evidence.

Their job is to preserve pressure and comparative reasoning.
Their job is not to certify the current stack.

## Conflict-resolution rules

When two files appear to support different conclusions, use these rules.

### Rule 1: Live runtime beats plan

If the plan says a service registry should exist but the root runtime does not
ship or consume it, the runtime wins for current-state claims.

### Rule 2: Dream surfaces beat calmer paraphrases

If a summary page sounds more generic than
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md) or
[`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md), trust the sharper dream statement.

### Rule 3: Authoring rules do not become failover proof

If [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules) or a service fragment enforces healthchecks or inline configs,
that still does not upgrade into route-survival or peer-forward correctness.

### Rule 4: Research pressure explains intent, not state

If the archive shows the repo strongly wants a feature, that still does not
make the feature live.

### Rule 5: Knowledgebase honesty rules constrain every other page

If a page sounds smoother than the current implementation deserves,
[`knowledgebase/AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/AGENTS.md) should be treated as the corrective.

## The biggest ways contributors still misread authority here

These still count as misunderstandings:

- "Several files say similar things, so the runtime is probably close."
- "The repo has a clear architecture story, so the missing part is mostly
  implementation detail."
- "The plan, README, and runtime all point in the same direction, so that
  direction must already be semi-owned by the stack."
- "Healthchecks, labels, and Traefik routers imply the system is already
  failover-capable in spirit."
- "A planning file that feels emotionally relieving is evidence of technical
  closure."

Those are exactly the moves that convert alignment into counterfeit maturity.

## Bottom line

The authority split in this repo is brutally simple:

- `copilot-instructions.md` names the dream
- `README.md` protects the dream from being softened
- `docker-compose.yml` and active fragments describe the live runtime surface
- `AGENTS.md` anchors work back to that runtime
- `.cursorrules` enforces Compose discipline
- the plan and research files explain pressure and candidate futures

Those surfaces support each other.
They do not add up to one stronger proof class.

That sentence needs to stay explicit because the repo is now coherent enough
to fool readers in a more sophisticated way:

it can sound fully assimilated before it is actually less dependent on hidden
operator knowledge.
