# Evidence Ledger

This page is the claim governor for the whole knowledgebase.

It exists because the main documentation failure in this repo is not missing
information.
It is information being narrated at the wrong confidence level.

The user's complaint is not mainly that earlier docs were too short.
It is that they kept sounding more complete, more settled, and more proved than
the repo had actually earned.

There is a more specific version of that complaint underneath the wording:

the docs kept reconstructing a smaller dream than the user actually has.

They would often preserve pieces of the stack, pieces of the pressure, and
pieces of the future direction, but they would fail to preserve the full
negative benchmark:

- any healthy node should be a legitimate first hop
- wrong-node requests should stop collapsing back into sacred host knowledge
- anti-SPOF language should correspond to where truth and authority really live
- the system should not need a giant scheduler worldview unless that worldview
  has clearly earned itself

That is why this ledger exists.
It is not just preventing factual error.
It is preventing a subtler kind of dishonesty:

> the documentation becomes more coherent by becoming less faithful to the real
> ask

So this page answers one strict question:

> for each major architecture claim, what evidence supports it, what class of
> evidence is it, and exactly where does that evidence stop?

If a page cannot answer that cleanly, it is still too loose for this project.

This ledger is therefore not an appendix.
It is one of the main devices stopping the site from drifting back into
"interesting architecture summary" mode.

It is also one of the pages that most directly decides whether this docs set is
actually doing retrieval-like reconstruction or merely polished summarization.

If the rest of the site is trying to preserve the real dream, this page is the
surface that says what preservation means operationally:

- preserve the evidence class
- preserve the confidence ceiling
- preserve the contradiction
- preserve the exact place where proof stops

It also has to preserve something easier to lose:

the difference between a claim that sounds adjacent to the user's dream and a
claim that actually crosses one of the user's real negative benchmarks.

That matters because this repo is no longer endangered mainly by missing facts.
It is endangered by flattering adjacency.

Examples:

- "multi-node entry exists" can sound close to "wrong-node requests stop being
  humiliating"
- "dynamic config exists" can sound close to "the recovery route survives the
  failure that made recovery necessary"
- "stateful service is reachable elsewhere" can sound close to "authority and
  data ownership survived coherently"

Those are not tiny wording differences.
They are exactly the gaps the user keeps reacting to.

This ledger exists so those gaps stay named instead of being smoothed into
confidence.

## What this ledger is trying to prevent

This ledger exists to block recurring documentation failures:

- live-runtime claims made from planning text
- failover claims made from proxy presence alone
- distributed-truth claims made from modular Compose authoring
- stateful HA claims made from ingress cleverness
- architecture-closure claims made from exploratory artifacts
- smooth vocabulary upgrades where "dynamic," "multi-node," or "resilient"
  suddenly start meaning more than the evidence supports

Another way to say the same thing is:

`bolabaden-infra` has enough real pressure, enough real edge machinery, and
enough real future-path exploration to sound finished long before it actually
is.

The repo can survive imperfect implementation.
It cannot survive docs that erase which parts are still imperfect.

The deeper failure being guarded against is:

> the docs become coherent by becoming less faithful

That is the wrong kind of coherence for this repo.

This is one of the most important rules in the whole site:

coherence is not automatically a virtue here.

If coherence was purchased by shrinking the dream, smoothing the contradiction,
or silently upgrading one evidence class into another, the page has improved
its tone while damaging its truthfulness.

There is also a temporal aspect the docs must preserve.

The archive, current runtime, and planning layers are not interchangeable
snapshots of one stable system.
They are fragments from a moving argument about what this infrastructure should
become.

If a page fuses those fragments into a single smooth present-tense story, it
does not merely blur evidence classes.
It erases the actual development pressure that makes the repo intelligible.

## How to use this ledger

Every serious page in the knowledgebase should be readable through five fields:

- the claim
- the evidence class
- the confidence
- what the evidence genuinely proves
- what it still does not prove

The last field is the most important one.

This repo does not mainly need more confident prose.
It needs prose that stops exactly where the evidence stops.

And it needs prose that says why stopping there matters.

This repo is full of components, plans, and experiments that can make a page
sound smarter simply by being mentioned together.
The ledger is what prevents "mentioned together" from turning into "proved
together."

It also needs prose that preserves what kind of thing is being evidenced:

- live runtime truth
- explicit repo desire
- promotion direction
- recurring user refusal
- outside-looking possibility that still has not crossed into local proof

If those are not kept apart, "evidence" becomes a flattering mood instead of a
discipline.

Another rule belongs here explicitly:

the ledger should not only stop false claims.
It should stop wrong *classes* of satisfaction.

In this repo, many failures begin when a page quietly treats one of the
following as emotionally equivalent to the real ask:

- first-hop survivability
- rich local proxy policy
- generated config
- cross-node reachability
- a plausible future control layer

All of those may matter.
None of them automatically satisfy the user's actual benchmark.

The benchmark remains harsher:

- can the wrong node remain honest?
- can the system externalize current truth?
- can recovery survive without private operator recollection becoming the real
  control plane?

## The evidence classes

The knowledgebase uses four evidence classes.
They are ordered by what they are allowed to prove.

## Class 1: live implementation evidence

Primary sources:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active fragments under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- direct render and inventory reads from `docker compose config`, when env and
  secrets are sufficiently prepared

This class can prove:

- what the tracked root runtime currently includes
- which services, labels, networks, configs, volumes, and secrets are authored
- which runtime surfaces are visibly present in the worktree

This class cannot prove by itself:

- that peer forwarding works under failure
- that routes persist when local backends disappear
- that cross-node placement truth converges
- that a TCP-exposed datastore is meaningfully HA
- that stateful failover is coordinated correctly

This is the class most likely to be overused.
Live components are real.
In this repo, "real component exists" is usually only the beginning of the hard
question, not the end of it.

That is because the user's complaint is rarely "the component is missing."
It is usually one of these harsher complaints:

- the component still depends on a sacred node
- the component still assumes the operator remembers placement truth
- the component still fails once the wrong machine receives the request
- the component exists, but only one human can explain how it remains correct
  after drift or failure

That is especially true for:

- `docker-gen-failover`
- multi-node DNS entry
- Traefik TCP exposure
- rich observability surfaces

Each of those can make the stack look more mature than the proof actually is.

That is why live implementation evidence has to be read with an extra
question:

> does this artifact prove preserved behavior under the relevant failure, or
> does it merely prove that another real component exists?

That question should often be followed immediately by another:

> if this component disappeared from the docs, would the actual runtime truth
> get weaker, or would only the story get less impressive?

That second question is crucial because this repo contains many real components
whose presence increases sophistication faster than it increases honest
cross-node truth.

## Class 2: repo-native intent evidence

Primary sources:

- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

This class can prove:

- what operating model the repo most clearly intends
- that the implementation is still Compose-first
- that resilience-oriented authoring discipline is intentional
- that the dream is explicit rather than inferred

This class cannot prove:

- that the intended operating model is fully shipped
- that a described control surface is actually live
- that the runtime behaves as intended under failure

This class matters especially here because the dream is unusually clear.
But clarity of dream is still not runtime proof.

It is still crucial evidence, though, because the user's actual standards are
not recoverable from the Compose graph alone.

It is also crucial because this repo would become easier to misread in the
opposite direction without it.

If Class 1 alone governed the site, readers could conclude:

- the current runtime surface is the whole problem
- the user's harsher dream is mostly inferred later
- the repo is mainly one more large Compose stack

That would be wrong.
The repo-native intent surfaces are what keep the docs tied to the real
benchmark instead of shrinking down to whatever the current runtime already
makes easy to say.

If this class were ignored, the docs would regress into a different lie:

> the current runtime surface is the whole problem statement

That would be just as misleading as letting planning text impersonate shipped
behavior.

## Class 3: planned architecture evidence

Primary sources:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
- [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
- [`docs/brainstorms/20260604-failover-agent-exploration.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/brainstorms/20260604-failover-agent-exploration.md)

This class can prove:

- which capability gaps the repo already recognizes
- where current mechanisms are known to be weak
- what future control surfaces are being promoted seriously
- how the repo is thinking about ingress versus stateful HA

This class cannot prove:

- that promoted mechanisms are active today
- that chosen direction is final
- that future architecture has crossed into runtime proof

This is where the repo's "there must be a better middle layer" pressure is most
visible, but it is still planning evidence, not production proof.

Planning evidence matters a lot in this repo.
It tells us what the worktree already knows it is missing.

That is one reason planning evidence is unusually valuable here:

it often exposes where the repo is already smarter about its own missing truth
than the current runtime is.

That does not upgrade the plan into proof.
But it does help the docs resist a different kind of false simplification:

> if the runtime does not prove it yet, the repo must not really understand the
> gap yet either

That would also be wrong.

But that is still different from proving the missing layer has crossed into the
runtime and stayed there under stress.

## Class 4: research-pressure evidence

Primary sources:

- archive corpus under `knowledgebase/source-archive/`
- archive-derived synthesis pages such as
  [`archive-pressure-patterns.md`](archive-pressure-patterns.md)

This class can prove:

- what problem shape keeps recurring
- what kinds of solutions the user repeatedly rejects
- what standards the docs must keep visible

This class cannot prove:

- implementation
- deployment maturity
- runtime behavior
- correctness under failure

This class matters because the archive keeps recovering what the user is really
rejecting:

- fake HA language
- sacred public boxes
- premature platform tax
- ambiguity about whether the route really survives the wrong failure

The archive also recovers a subtler demand:

the user does not merely want more options.
They want the option space to stop hiding the real missing truth behind
pleasant wording.

The archive also proves something about documentation style:

the user is not comforted by pages that sound balanced while evading the actual
wound.
If the docs politely rotate around the real pressure without naming it, they
become another version of the same ecosystem failure the repo is reacting to.

That sentence applies to this ledger too.

If the ledger itself ever starts sounding neutral about evidence inflation, it
has already stopped governing claims and started participating in the same
problem.

## Reading discipline

The practical meaning of the evidence classes is:

- positive artifacts do not automatically outrank missing artifacts
- a live component does not automatically outrank a live contradiction about
  that component
- elegant design language never outranks a weaker current implementation
- operator frustration recovered from the archive can sharpen the question, but
  still cannot impersonate shipped proof

The archive can tell us what the user is really rebelling against.
It cannot tell us that the current runtime has defeated it.

If a future page sounds cleaner by collapsing these boundaries, that
cleanliness should be treated as a warning, not a win.

That warning is especially important in the current repo state because so many
pieces are already real enough to be rhetorically dangerous:

- the stack is large
- the edge is serious
- the plans are specific
- the option space is broad

All of that makes it easy to sound comprehensive while still failing the
user's harder benchmark.

That warning matters because this project naturally invites a deceptive
writing style:

- the repo has enough moving pieces to sound sophisticated
- the plans are specific enough to sound inevitable
- the intent surfaces are strong enough to sound shipped
- the archive is rich enough to sound like proof of direction

The ledger exists to keep those strengths from becoming narrative fraud.

## The strongest current anchors

If a maintainer only checks a small set of files before editing docs, these are
the best current anchors.

## Strongest live anchors

1. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
3. [`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
4. [`compose/docker-compose.docs.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.docs.yml)
5. [`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)

## Strongest intent and gap anchors

6. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
7. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
8. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
9. [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
10. [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)

These anchors are not equal in what they are allowed to prove.

- live Compose files are strongest for runtime truth
- `copilot-instructions.md` is strongest for the intended no-heavy-orchestrator
  request contract
- planning docs are strongest for named future control surfaces and recognized
  gap structure

If those roles get merged into one generic "architecture" category, the docs
drift again.

## Current critical negative facts

Absence is evidence in this repo.
Some of the most important truths are the things the current worktree still
does not prove.

## Negative fact 1: tracked root `services.yaml` is absent

What this proves:

- the service-registry concept is architecturally central
- the tracked root runtime still does not expose that registry as live placement
  truth

What should be inferred carefully:

- many pages may speak as if such a truth surface is conceptually present
- the evidence ceiling here still says the priority runtime has not promoted it
  to a tracked live root artifact

Why this negative fact matters:

the repo's central dream keeps converging on placement truth that is readable,
current, and not secretly human-only.
If the tracked root runtime still lacks that artifact, the docs are not allowed
to narrate the core orchestration wound as if it is already structurally
closed.

What this does not prove:

- that operators have no placement knowledge at all

It only proves that the shared tracked source of truth is still missing from
the priority runtime.

## Negative fact 2: `docker-gen-failover` is present, but the repo does not trust it as finished HA

Current evidence split:

- the live proxy fragment contains `docker-gen-failover`
- planning docs explicitly describe route-loss and fallback brittleness

What this proves:

- generator-driven failover is live as an implementation attempt
- generator presence is not the same thing as trusted request preservation

What this does not prove:

- that the replacement path is already active

## Negative fact 3: stateful services are already part of the priority runtime

Current live examples include:

- `mongodb`
- `redis`
- `nuq-postgres`
- `rabbitmq`

What this proves:

- stateful correctness is not a future-only concern
- any HA language in this repo must distinguish ingress resilience from data
  correctness

What this does not prove:

- that any of those services already have meaningful multi-node failover

## Negative fact 4: the truth surfaces are still mixed even inside the "better middle layer" experiments

Current examples:

- the HTTP OpenSVC generator reads node names from `om node ls --format json`
- the current L4 sync generator reads node IPs from `tailscale status --json`
- both are meaningful
- neither fact means the repo has already converged on one final control-plane
  truth source

What this proves:

- the repo is actively experimenting with stronger routing truth
- the experiments are still mixed and should not be narrated as a unified solved
  substrate

This matters because mixed truth sources are not just an implementation detail.
They are part of the current proof ceiling.

What this does not prove:

- that the side-paths are incoherent or worthless

It only proves that they are still experiments, not closure.

## The most important current claims and their ceilings

This is the shortest useful summary of where the strongest evidence stops
today.

| Claim | Strongest evidence class currently available | What that evidence genuinely proves | What still remains unproven |
| --- | --- | --- | --- |
| The repo wants no-heavy-orchestrator multi-node Docker | Repo-native intent evidence | The dream is explicit and central, not inferred | That the full operating contract is live |
| The root stack is large, serious, and Compose-first | Live implementation evidence | The current runtime surface is substantial and real | That the stack already has distributed truth and resilient wrong-node behavior |
| Cloudflare and Traefik are central to the ingress story | Live implementation plus repo-native intent evidence | The edge surface is materially real and intentionally central | That ingress survives backend-loss and peer-handoff semantics end to end |
| A lightweight current-state registry is architecturally central | Repo-native intent plus planned architecture evidence | `services.yaml` is a major design concept | That a tracked root registry is live and consumed by the runtime |
| The repo already distrusts fake failover stories | Repo-native intent plus planned architecture plus archive pressure | The honesty boundary is intentional | That the live stack already crossed those proof boundaries |
| Stateful HA is a separate problem class | Repo-native intent, planning, and live stateful presence | The repo understands ingress success is not data correctness | That Redis, MongoDB, RabbitMQ, or Postgres have meaningful proven failover behavior |

## Bottom line

This page exists to enforce one discipline:

> let evidence establish the size of the claim, not the other way around

If a future doc change wants the architecture to sound more complete than this
ledger supports, it should either bring stronger evidence or make a smaller
claim.

That is not caution for its own sake.
It is the only way these docs stay faithful to the user's real demand: stop
pretending the option space is richer, cleaner, or more solved than the actual
worktree and proof surfaces warrant.
