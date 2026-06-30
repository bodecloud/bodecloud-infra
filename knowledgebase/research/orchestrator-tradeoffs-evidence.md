# Orchestrator and Control-Plane Tradeoffs Evidence

This page is the proof boundary for the repo's platform-choice story.

It exists to stop one very specific documentation failure:

the point where "there are several possible futures here" gets rewritten into
"there is one obvious future and the docs are just slowly getting comfortable
enough to admit it."

That smoother story is exactly what this repo has to resist unless the evidence
really earns it.

And the repo has to resist it even when the smoother story starts sounding
more adult than the rough one.
The user is explicitly pushing against an ecosystem where adulthood is too
often faked by enlarging the worldview rather than actually relocating the
hidden burden.

The real question is not:

> Which orchestrator is best in the abstract?

It is:

> Which extra layer of machinery is actually justified by the specific hidden
> human SPOFs, wrong-node failures, convergence failures, and state-truth gaps
> this stack is already carrying?

That distinction is the whole page.

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

The real threshold is not which product sounds right.
It is whether a candidate removes the specific hidden burden that currently
keeps the operator acting as the control plane.

That sentence should dominate the whole page.
If a comparison sounds sophisticated but stops naming the surviving operator
burden, it has already drifted back into product-literacy theater.

## What this page is and is not allowed to prove

This page is allowed to prove:

1. the current runtime is still concretely Compose-first
2. the repo is already carrying pain that plain Compose does not close
3. the preserved alternatives map to different missing truth classes
4. the user is refusing fake closure, not randomly hesitating
5. the repo is already partway toward inventing a control plane around Compose
6. no current evidence proves that one platform has already earned whole-stack
   promotion
7. the decision must stay tied to named failure classes, not platform prestige
8. the real threshold is whether the missing middle layer stays narrow or
   silently grows into a scheduler in disguise

This page is not allowed to:

- crown a winner
- treat exploration artifacts as implementation
- imply that all futures are equally justified

## Claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "Compose is still the live baseline" | Class 1 | `docker-compose.yml`, active Compose fragments | that Compose is sufficient for every failure class |
| "the repo wants multi-node Docker without immediate heavyweight capture" | Class 2 | `.github/copilot-instructions.md`, `README.md` | that scheduler promotion is off the table forever |
| "the repo already knows which missing capabilities hurt" | Class 3 | `docs/INFRASTRUCTURE_MASTER_PLAN.md`, `docs/osvc_ingress_ha.md`, `docs/stateful_ha_plan.md` | that one future has already earned promotion |
| "the user keeps rejecting fake closure" | Class 4 | archive threads and archive-pressure pages | that hesitation is random or that all options deserve equal weight |

This page should never allow "there are several futures" to become "all
futures are equally justified."

That is one of the softer lies advanced documentation likes to tell.
Equal narratability is not equal justification.
Some futures are still mostly atmosphere unless they can show one more truth
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

- serious edge behavior in `compose/docker-compose.coolify-proxy.yml`
- mesh and control-plane pressure in `compose/docker-compose.headscale.yml`
- observability and metrics surfaces in `compose/docker-compose.metrics.yml`
- alternate orchestration exploration such as `compose/docker-compose.nomad.yml`
- additional routing/control fragments such as
  `compose/docker-compose.warp-nat-routing.yml`

That proves:

- the repo is not naively "just a few containers"
- the pressure toward a broader control plane is already visible in the
  worktree

It does **not** prove:

- that any one of those futures has already won

### 3. The repo is already inventing helper control-plane behavior around Compose

The master plan names several missing truths that are not pure container-start
problems:

- secret sync
- compose sync
- service registry via `services.yaml`
- automated service failover
- Headscale-assisted peer broadcast

That matters because it means the repo is not debating control planes from
outside.
It is already partially building one around Compose.

The real tradeoff is therefore not:

- Compose
- versus orchestrator

It is:

- which control-plane growth is narrow and honest
- versus which growth silently becomes a scheduler in disguise

That is the real differentiator this repo cares about.
Not only "lighter versus heavier,"
but "does the new layer stay narrow enough to feel like an earned answer
instead of another worldview tax pretending to be inevitable?"

## What the planning layer says the real wounds are

The planning docs are already unusually explicit about the actual burden set.

[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
names all of these as still missing:

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

not:

- which candidate is most famous?

This repo should keep saying that because fame is one of the most dangerous
forms of false reassurance in the whole control-plane conversation.
Famous platforms can still leave the user's actual accusation unanswered if
the docs stop too early at capability language.

## What the archive contributes

The archive shows that the user is not randomly "not ready" for an
orchestrator.
The user is refusing fake closure.

Important archive facts:

- `docker-multi-node-without-swarm__...` repeatedly converges on service
  discovery as the hard missing piece
- the same thread describes Nomad as lighter than Kubernetes but still frames
  the real issue as current placement truth
- that thread also mentions Consul and gossip-based approaches, but keeps the
  burden centered on "service name -> where is it running right now"
- `distributed-ha-orchestration__...` explicitly notes that K3s still uses
  leader election and Swarm still has manager nodes
- `docker-compose-frustration__...` captures the "Docker feels gaslighting"
  pressure that explains why superficial option lists are not enough

These sources do not pick the winner.
They do explain why the repo keeps several futures alive without that being
mere indecision.

That matters because indecision is the lazy reading.
The stronger reading is that the repo is refusing to award adulthood to a
candidate before it can prove one more humiliating failure class has actually
become less privately human.

## What still does not count as orchestration evidence here

This repo needs a harsher evidence filter than ordinary platform-comparison
docs.

The following are still too weak to count as serious promotion evidence on
their own:

- a cleaner deployment syntax
- a smoother local demo
- a broader feature matrix
- a controller that can describe more nouns
- a candidate that sounds "lighter than Kubernetes"
- a candidate that makes the diagrams calmer while the operator still privately
  closes the same route, placement, or failover gaps

Those things may matter later.
They do not yet answer the user's real question.

That question should keep sounding a little abrasive here.
The user is not waiting for a prettier matrix.
The user is waiting for some candidate to stop making "just accept a larger
system" sound like a social threat instead of an earned technical answer.

The real question is still:

- which layer now owns more of the missing truth
- which hidden burden became inspectably less human
- which failure class became less humiliating on the bad day

If none of those changed, then this repo did not actually gain a more honest
future.

## What a real promotion packet would have to contain

No future path should get elevated just because it sounds increasingly
coherent in docs.

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

That packet does not need to prove whole-stack closure.
It does need to prove that at least one hidden burden became materially less
private.

That is the whole bar.
Not a better roadmap.
Not a better ecosystem fit.
Not a calmer explanation.
One more important truth has to stop living as private operator folklore.

If a candidate cannot leave that kind of artifact bundle behind, then it has
not yet earned stronger language than "still being considered."

## Fake differentiation versus real differentiation

This repo needs a stricter filter than ordinary platform comparison pages use.

Two options are not meaningfully different here unless they reduce
meaningfully different hidden burdens.

That means:

- different deployment syntax is not enough
- different ecosystem size is not enough
- different controller branding is not enough
- different HA vocabulary is not enough

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

These are the main candidate families already present in planning or research.

| Candidate family | What it might solve | Why it is still incomplete |
| --- | --- | --- |
| lightweight registry plus helper agents | placement truth, sync, redeploy, failover glue | still must prove durability, eligibility, and semantic preservation |
| OpenSVC-shaped ingress or service supervision | stronger service ownership and failover semantics | still needs proof that it removes burden instead of renaming it |
| Nomad-style promotion | scheduling and health-aware placement | must earn its worldview tax and not just widen the control plane |
| k3s / Kubernetes promotion | broad cluster truth and scheduling machinery | may solve many layers, but the repo treats lost legibility and central worldview cost as real |
| improved proxy automation alone | route generation and local edge behavior | not enough unless it owns placement, durability, and explanation too |

## The success test

A future path has only earned promotion when it makes at least one of the
current hidden burdens materially less true in shared, inspectable system
behavior.

That means the path has to do more than sound plausible.
It has to reduce one or more of these:

- wrong-node requests collapsing back into private operator memory
- remembered placement still being the real registry
- route persistence still failing under local backend loss
- secrets and env drift still making peer substitution dishonest
- stateful authority still secretly remaining singular

If a candidate does not do that, then for this repo it is mostly a renamed
option, not a new answer.

## Bottom line

This repo is not asking "which orchestrator is best?"

It is asking:

> what extra layer of truth ownership is actually justified by the failure
> classes we have already named?

That is why the docs must stay hostile to fake closure.

The real threshold is not popularity or even feature depth.
It is whether the next layer removes a real hidden burden instead of merely
relocating it under more prestigious vocabulary.
