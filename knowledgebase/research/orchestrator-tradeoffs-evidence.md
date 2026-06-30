# Orchestrator and Control-Plane Tradeoffs Evidence

This page is the proof boundary for the repo's platform-choice story.

It exists to stop one recurring documentation failure:

the point where
"there are several possible futures here"
gets rewritten into
"there is one obvious future and the docs are just slowly getting honest
enough to admit it."

That smoother story is exactly what this repo has to resist unless the
evidence really earns it.

The danger is not merely picking the wrong orchestrator.
The danger is letting "more adult architecture" become another synonym for:

> the same hidden burden still lands on the operator, but now the burden has
> better branding and more prestigious nouns around it.

That sentence should govern the whole page.

The real question is not:

> Which orchestrator is best in the abstract?

It is:

> Which extra layer of machinery is actually justified by the specific hidden
> human SPOFs, wrong-node failures, convergence failures, and state-truth gaps
> this stack is already carrying?

## Strongest honest current answer

The repo is not comparing futures by branding.
It is comparing them by burden ownership.

The current worktree proves:

- Compose is still the live baseline
- plain Compose is already carrying named burdens it does not close honestly
- the repo is already experimenting with shadow-control-plane behavior around
  Compose
- several future layers are being considered for different reasons

The current worktree does **not** prove:

- that one whole-stack future has already earned promotion
- that "lighter than Kubernetes" is automatically a good answer
- that keeping multiple futures alive is just indecision

The real threshold is whether a candidate removes one more specific hidden
burden that currently keeps the operator acting as the control plane.

## What this page is and is not allowed to prove

This page is allowed to prove:

1. the current runtime is still concretely Compose-first
2. the repo is already carrying pain that plain Compose does not close
3. the preserved alternatives map to different missing truth classes
4. the user is refusing fake closure, not randomly hesitating
5. the repo is already partway toward inventing a control plane around Compose
6. no current evidence proves that one platform has already earned
   whole-stack promotion
7. the decision must stay tied to named failure classes, not platform
   prestige
8. the real threshold is whether the missing middle layer stays narrow or
   silently grows into a scheduler in disguise

This page is not allowed to:

- crown a winner
- treat exploration artifacts as implementation
- imply that all futures are equally justified

## Retrieval contract for this page

### Class 1: live implementation baseline

Primary anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active Compose fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)

This class is allowed to prove:

- what the repo currently ships
- what the actual authoring surface still is

It is not allowed to prove:

- that Compose remains sufficient for every failure class

### Class 2: intent and architecture dream

Primary anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

This class is allowed to prove:

- the repo wants multi-node Docker without immediate heavyweight capture
- scheduler promotion is not the default emotional end state

It is not allowed to prove:

- that scheduler promotion is permanently off the table

### Class 3: repo-native gap naming

Primary anchors:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
- [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)

This class is allowed to prove:

- which missing truths are already painful enough to name
- which present helper layers are still too weak

It is not allowed to prove:

- that one future has already won

### Class 4: archive pressure

Primary anchors:

- archive comparison threads
- archive frustration threads
- archive orchestration threads

This class is allowed to prove:

- the user is refusing fake closure
- platform comparison here is not an abstract product-literacy exercise

It is not allowed to prove:

- current runtime correctness

If a paragraph moves from one class to another, the docs should say so.
Otherwise dream force and archive force will quietly impersonate implementation
force.

## Claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "Compose is still the live baseline" | live implementation | `docker-compose.yml`, active Compose fragments | that Compose is sufficient for every failure class |
| "the repo wants multi-node Docker without immediate heavyweight capture" | intent | `.github/copilot-instructions.md`, `README.md` | that scheduler promotion is forbidden forever |
| "the repo already knows which missing capabilities hurt" | repo-native gap naming | master plan and HA planning docs | that one future has already earned promotion |
| "the user keeps rejecting fake closure" | archive pressure | archive threads and archive-pressure synthesis | that hesitation is random or unserious |

Equal narratability is not equal justification.
Some futures are still mostly atmosphere until they can show one more truth
becoming inspectably system-owned.

## What the live worktree concretely proves

### 1. The runtime is still authored as Compose

The current operator surface is still driven from:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- the root include graph
- Docker-native labels, networks, configs, and service definitions

That proves:

- the active control surface is still file-first and Docker-native

It does **not** prove:

- that Compose remains sufficient for every future failure class

### 2. The repo is already broader than simple local Compose

The worktree already contains more than one future-facing direction:

- serious edge behavior in
  [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- mesh and control-plane pressure in
  [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- observability and metrics surfaces in
  [`compose/docker-compose.metrics.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.metrics.yml)
- alternate orchestration exploration such as
  [`compose/docker-compose.nomad.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.nomad.yml)
- additional routing/control fragments such as
  [`compose/docker-compose.warp-nat-routing.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.warp-nat-routing.yml)

That proves:

- the repo is not naively "just a few containers"
- the pressure toward a broader control plane is already visible in the
  worktree

It does **not** prove:

- that any one of those futures has already won

### 3. The repo is already inventing helper control-plane behavior around Compose

The planning layer names missing truths that are not pure
"start this container" problems:

- secret sync
- compose sync
- service registry via `services.yaml`
- automated service failover
- Headscale-assisted peer broadcast

That matters because the repo is not debating control planes from outside.
It is already partially building one around Compose.

So the real tradeoff is not simply:

- Compose
- versus orchestrator

It is:

- which control-plane growth is narrow and honest
- versus which growth silently becomes a scheduler in disguise

That distinction matters more than "lightweight versus heavyweight."
A heavyweight layer can still be the more honest answer if it truly owns the
right truth.
A lightweight layer can still be fake if it mostly preserves the operator's
private settlement role under gentler vocabulary.

## What the planning layer says the real wounds are

[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
is already unusually explicit about the real burden set.

It names all of these as still missing:

- universal wrong-node success
- a live tracked root `services.yaml` current-state registry
- trustworthy route persistence under local backend failure
- automated service failover between nodes

It also names concrete current pain:

- Headscale is single-node today
- Cloudflare DDNS is present but not the same as full multi-node failover
- `docker-gen-failover` currently deletes routes on container stop

That is the right lens for comparing future layers.

The question becomes:

- which candidate most honestly owns one or more of those missing truths?

Not:

- which candidate is most famous?

This repo should keep saying that because fame is one of the easiest ways to
turn capability language into false reassurance.

## What the archive contributes

The archive shows that the user is not randomly "not ready" for an
orchestrator.
The user is refusing fake closure.

Important archive facts:

- `docker-multi-node-without-swarm__...` keeps converging on service discovery
  as the hard missing piece
- the same thread is especially important because the user explicitly accepts
  manual placement and narrows the real unsolved problem to request-time
  service discovery
- the same thread describes Nomad as lighter than Kubernetes but still frames
  the real issue as current placement truth
- the same pressure repeatedly returns to "service name -> where is it running
  right now"
- `forking-docker-compose__...` shows that the user was not asking for prettier
  YAML; they were looking for runtime fallback behavior such as failover lists,
  startup-failure rescue, and hybrid fallback targets without maintaining a
  permanent Compose fork
- `dynamic-ha-proxy-setup__...` shows why IP-level or shared-entry arguments
  must be separated from service-level correctness; even if first-hop SPOF is
  reduced, the wrong-node service decision can still be unowned
- `distributed-ha-orchestration__...` explicitly notes that K3s still uses
  leader election and Swarm still has manager nodes
- `docker-compose-frustration__...` captures the "Docker feels gaslighting"
  pressure that explains why superficial option lists are not enough

These sources do not pick a winner.
They do explain why the repo keeps several futures alive without that being
mere indecision.

The stronger reading is:

- the repo has not yet seen enough evidence to let one future legally
  impersonate relief
- the preserved futures are different attempts to buy specific truth ownership,
  not a pile of interchangeable platform names
- the smallest faithful answer may be a current-state/fallback decision layer,
  but only if it actually consumes runtime truth and acts during failure

## What still does not count as orchestration evidence here

These are still too weak to count as promotion evidence on their own:

- a cleaner deployment syntax
- a smoother local demo
- a broader feature matrix
- a controller that can describe more nouns
- a candidate that sounds "lighter than Kubernetes"
- calmer diagrams while the operator still privately closes the same route,
  placement, or failover gaps

Those things may matter later.
They do not yet answer the user's real question.

That real question should stay abrasive here:

- which layer now owns more of the missing truth?
- which hidden burden became inspectably less human?
- which failure class became less humiliating on the bad day?

If none of those changed, the repo did not actually gain a more honest future.

## What a real promotion packet would have to contain

No future path should get elevated just because it sounds increasingly coherent
in docs.

For this repo, a serious promotion packet would need a bundle of evidence such
as:

- the exact failure class being narrowed or removed
- the truth layer being promoted from operator memory into system-owned or
  system-consumable state
- the concrete artifact that carries that truth, such as a live placement
  registry, stronger failover ownership, or inspectable peer-eligibility data
- a wrong-node or backend-loss drill showing what the system actually does now
- operator-legibility proof showing where a human can inspect the decision path
- an explicit statement of what the promoted layer still does not solve

That packet does not need whole-stack closure.
It does need to prove that at least one hidden burden became materially less
private.

That is the minimum seriousness threshold for this page:
not which future feels most plausible,
but which future can prove that one more bad-day explanation no longer depends
on operator folklore.

## Fake differentiation versus real differentiation

This repo needs a stricter filter than ordinary platform-comparison pages use.

Two options are not meaningfully different here unless they reduce
meaningfully different hidden burdens.

That means these are not enough on their own:

- different deployment syntax
- different ecosystem size
- different controller branding
- different HA vocabulary

The differentiators that actually matter are:

- who now owns placement truth?
- who now owns convergence truth?
- who now owns peer eligibility?
- who now owns route persistence?
- who now owns stateful authority?
- what worldview cost is being imposed in exchange?

That last question matters because the user is not refusing platform tax in the
abstract.
They are refusing to pay scheduler-scale cost before the specific missing truth
is named tightly enough to justify the bill.

## Candidate families already visible in the repo

| Candidate family | What it might solve | Why it is still incomplete |
| --- | --- | --- |
| lightweight registry plus helper agents | placement truth, sync, redeploy, failover glue | still must prove durability, eligibility, and semantic preservation |
| OpenSVC-shaped ingress or service supervision | stronger service ownership and failover semantics | still needs proof that it removes burden instead of renaming it |
| Nomad-style promotion | scheduling and health-aware placement | must earn its worldview tax and not just widen the control plane |
| k3s / Kubernetes promotion | broad cluster truth and scheduling machinery | may solve many layers, but the repo treats lost legibility and worldview tax as real costs |
| improved proxy automation alone | route generation and local edge behavior | not enough unless it owns placement, durability, and explanation too |

## Success test

A future path has only earned promotion when it makes at least one current
hidden burden materially less true in shared, inspectable system behavior.

That means the path has to reduce one or more of these:

- wrong-node requests collapsing back into private operator memory
- remembered placement still being the real registry
- route persistence still failing under local backend loss
- secrets and env drift still making peer substitution dishonest
- stateful authority still secretly remaining singular

If a candidate does not do that, then for this repo it is mostly a renamed
option, not a new answer.

## Bottom line

This repo is not asking:

> which orchestrator is best?

It is asking:

> what extra layer of truth ownership is actually justified by the failure
> classes we have already named?

That is why the docs must stay hostile to fake closure.

The real threshold is not popularity or even feature depth.
It is whether the next layer removes a real hidden burden instead of merely
relocating it under more prestigious vocabulary.
