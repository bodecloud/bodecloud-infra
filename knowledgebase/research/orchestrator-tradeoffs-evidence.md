# Orchestrator and Control-Plane Tradeoffs Evidence

This page is the proof boundary for the repo's platform-choice story.

It also exists to stop a very specific doc failure:

the point where "there are several possible futures here" gets rewritten into
"there is one obvious future and the docs are just slowly getting comfortable
enough to admit it."

That smoother story is exactly what this repo has to resist unless the evidence
really earns it.

The real question is not:

> Which orchestrator is best in the abstract?

It is:

> Which extra layer of machinery is actually justified by the specific hidden
> human SPOFs, wrong-node failures, convergence failures, and state-truth gaps
> this stack is already carrying?

That distinction matters because shallow infrastructure advice collapses the
question too early:

- use Kubernetes because it is powerful
- use Nomad because it is lighter
- use Swarm because it is Docker-native
- use OpenSVC because it solves HA

Those answers are not automatically wrong.
They are just not specific enough for this repo.

This page also exists because the archive keeps showing something more specific
than ordinary platform indecision.

The user repeatedly collides with the same sequence:

- Compose stays readable until distribution becomes real
- wrong-node entry exposes missing placement, routing, and convergence truth
- proposed fixes either hardcode more fragile glue or import a new worldview
- the answer starts sounding decisive before it has explained what exact wound
  it is healing

If that sequence is forgotten, the repo gets narrated like a normal
"which orchestrator should I pick?" exercise.
That would be false.

It would also flatten the user's actual frustration into consumer choice,
which is almost the opposite of what the repo is trying to recover.

The user is not trying to shop more elegantly.
The user is trying to stop being trapped between:

- fragile glue that keeps sacred-node truth private
- and giant worldview packages that demand surrender before they have proved
  they heal the right wound

That means this page cannot read like a buyer's guide with better prose.
It has to read like an argument about burden ownership.
If a platform path is described without naming which hidden burden it removes,
the description is too vague for this repo.

That is why this page has to stay hostile to fake differentiation.

Many future paths only become real options if they reduce a different hidden
burden than the current path.
If they mostly rename the same burden, then for this repo they are much closer
to the same option than to distinct ones.

## What this page is and is not allowed to prove

This page is allowed to prove that the repo is comparing future control layers
by burden ownership, not by branding, ecosystem prestige, or generic feature
lists.

It is not allowed to crown a winner or to imply that one platform path has
already earned whole-stack promotion in the current evidence.

It is also not allowed to confuse "there are several futures" with "all of
those futures are equally justified."

More concretely, it is allowed to prove eight narrower things:

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

That last point is one of the most important reconstruction results in the
whole knowledgebase.

The repo is not merely comparing products.
It is watching for the moment when its own "smaller" solution stops being small
in any honest sense.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "Compose is still the live baseline" | Class 1 | `docker-compose.yml`, active Compose fragments | that Compose is sufficient for every failure class |
| "the repo wants multi-node Docker without immediate heavyweight capture" | Class 2 | `.github/copilot-instructions.md`, `README.md` | that scheduler promotion is off the table forever |
| "the repo already knows which missing capabilities hurt" | Class 3 | `docs/INFRASTRUCTURE_MASTER_PLAN.md`, `docs/osvc_ingress_ha.md`, `docs/stateful_ha_plan.md` | that one future has already earned promotion |
| "the user keeps rejecting fake closure" | Class 4 | linked archive threads | that hesitation is random or that all options deserve equal weight forever |

This page should never allow "there are several futures" to become "all
futures are equally justified."

## What this page should let a reader answer immediately

After reading this page, a reader should not still be wondering:

- are these paths genuinely different?
- what exact hidden burden does each one remove?
- which ones are only different in branding or worldview?

Those answers need to be explicit.

The quickest useful summary is:

- Compose-first is still the live baseline
- helper growth is real and may already be approaching shadow-control-plane
  territory
- infra-grade HA promotion and scheduler promotion solve different wounds
- stateful truth remains separate from both
- no current evidence proves that one whole-stack future has already earned
  promotion over the others

If that summary cannot be recovered quickly, the page is still too soft.

It is also watching for a second threshold:

the point where "keeping options open" stops being evidence discipline and
starts becoming a way to avoid naming which specific burden should be promoted
next.

This page should keep both thresholds visible.

## What this page is protecting against

The user is tired of being forced into one of two fake choices:

- stay in fragile Compose forever
- or surrender the whole problem to Kubernetes, Swarm, or another grand
  platform before the repo has earned that move

This page exists to preserve the middle truth:

the repo has real pain that plain Compose does not solve honestly, but that
does not mean every extra layer of machinery is justified equally.

Another way to say it:

the repo is not trying to avoid control planes forever.
It is trying to avoid paying scheduler-scale cost before the smaller missing
truths have been isolated well enough to justify that bill.

That is why the page has to stay argumentative instead of tidy.

If the docs make this decision surface look too clean, they risk hiding the
main insight:
many proposed upgrades differ in branding much more than they differ in whether
they actually reduce the operator's private truth burden.

Another way to say it:

- the user is not asking for permission to stay naive forever
- the user is also not willing to pay for a giant scheduler just to avoid
  naming smaller missing truths precisely

That refusal pattern is one of the strongest findings in the whole archive.
If a future page smooths it into ordinary indecision, the reconstruction has
already failed.

That is not indecision.
It is discipline.

## Fake differentiation versus real differentiation

This repo needs a stricter filter than normal platform comparisons use.

Two options are not meaningfully different here unless they reduce
meaningfully different hidden burdens.

That means:

- different deployment syntax is not enough
- different ecosystem size is not enough
- different controller branding is not enough
- different HA vocabulary is not enough

The differentiator that matters is:

- who now owns placement truth?
- who now owns convergence truth?
- who now owns peer eligibility?
- who now owns route persistence?
- who now owns stateful authority?

One more differentiator belongs here:

- which new worldview costs are now being imposed in exchange?

That matters because the user is not refusing platform tax in the abstract.
They are refusing to pay worldview-scale cost before the specific missing truth
is named tightly enough to justify the bill.

If two futures answer those with roughly the same "still mostly the operator,"
then they are much closer to one fake choice than to two real ones.

That is why this page should read less like a buyer's guide and more like an
evidence ledger for when control-plane growth has actually earned itself.

## Evidence classes this page relies on

### Class 1: live implementation evidence

Used for:

- confirming the runtime is still Compose-first
- showing optional futures exist directly in the worktree
- proving the stack is already broader than trivial one-node Compose

### Class 2: repo-native intent evidence

Used for:

- confirming the repo's no-heavy-orchestrator default
- confirming the any-node, local-first, peer-forward dream

### Class 3: planned architecture evidence

Used for:

- identifying actual missing capability classes
- showing which stronger layers are being considered seriously
- separating ingress, placement, convergence, and stateful pain

### Class 4: archive-pressure evidence

Used for:

- explaining why the repo keeps multiple futures alive
- proving the user rejects both static glue theater and premature worldview
  capture

## Strongest live anchors

Primary files:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- [`compose/docker-compose.metrics.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.metrics.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- [`compose/docker-compose.warp-nat-routing.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.warp-nat-routing.yml)
- [`compose/docker-compose.nomad.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.nomad.yml)

## What the live worktree concretely proves

### 1. The runtime is still authored as Compose

The current operator surface is still driven from:

- root `docker-compose.yml`
- the root include graph
- Docker-native labels, networks, configs, and service definitions

What that proves:

- the active control surface is still file-first and Docker-native

What it does not prove:

- that Compose remains sufficient for every future failure class

### 2. The repo is already broader than simple local Compose

The worktree already contains:

- alternate Compose fragments
- non-primary orchestration artifacts such as the Nomad fragment
- multi-domain workload groupings
- generated-config expectations in the proxy layer

What that proves:

- the repo is still Compose-first, but it is no longer honestly describable as
  "just a compose file"
- the stack is already carrying enough moving truth that pretending raw
  Compose is the whole control surface would be dishonest

What it does not prove:

- that the right next layer has already been chosen

This is where the repo needs to be careful not to confuse preserved futures
with earned futures.

### 3. The worktree is preserving optional futures directly

Artifacts like
[`compose/docker-compose.nomad.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.nomad.yml)
matter even if they are not the priority live path.

What that proves:

- control-plane optionality is not just talk

What it does not prove:

- that those futures are equally mature or equally aligned

## What the live worktree does not prove

The worktree does not, by itself, prove:

- which future control plane should win
- that helper growth can continue forever without becoming its own hidden
  platform
- that Nomad, OpenSVC, k3s, or any other path has already earned promotion

That absence of proof is central here.

The current worktree is best read as:

- Compose is still the live baseline
- helper pressure is real
- control-plane pressure is real
- platform closure is still unresolved

It should also be read as a warning that helper growth is already carrying a
governance burden.
If helper growth starts owning placement truth, convergence truth, peer
eligibility, and repair semantics, then the repo is no longer merely keeping
Compose simple.
It is building a control plane and must be judged like one.

## What each broad family is actually buying in burden terms

This page needs this map because otherwise tradeoff language drifts back into
feature language.

| Family | Hidden burden it is mainly trying to remove | Burden it still leaves behind if oversold |
| --- | --- | --- |
| Compose-first plus helpers | remembered placement, remembered failover glue, some convergence work | helper sprawl can still leave private peer judgment and route semantics in the operator's head |
| Infra-grade HA substrate | sacred ingress or identity surfaces, first-hop fragility | ordinary service wrong-node preservation and stateful truth may still remain separate problems |
| Scheduler promotion | remembered placement and rescheduling burden for selected workloads | ingress meaning, stateful authority, and service-class-specific truth can still remain unresolved |
| Full desired-state platform | broad lifecycle, reconciliation, and ecosystem burden | the user can still be left with a harder-to-read system whose real stateful and semantic truths are not magically solved |

That table matters because it makes one recurring doc mistake harder:
pretending all "stronger" futures are buying the same kind of relief.

That unresolved closure is not an embarrassment.
It is one of the most faithful facts the docs can preserve.

## Promotion questions this page should force

Before any future gets promoted in prose, the page should be able to answer:

1. Which hidden human reconstruction burden does this path reduce first?
2. Which burdens remain operator-owned afterward?
3. Which service classes actually benefit from this promotion?
4. Does this path reduce wrong-node ambiguity, or just relocate it?
5. Does this path reduce stateful authority ambiguity, or only supervision
   ambiguity?
6. Is the repo paying for a worldview larger than the wound it is healing?

If those questions cannot be answered, the future is not mature enough to be
spoken about as if it had earned itself.

It is also why
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
matters more here than [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
or [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules).

Those other files help explain how to work in the repo.
The Copilot instructions are the clearer statement of what the repo is trying
to become without lying about already being there.

That hierarchy matters because a lot of infra repos have several guidance
files, but not all of them are architecture witnesses.
For this question:

- `copilot-instructions.md` is the strongest intent surface
- `AGENTS.md` is supporting operational context
- `.cursorrules` is mostly implementation discipline and hygiene

## Strongest intent and planning anchors

Primary files:

- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
- [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)

## What these files explicitly say

### 1. The repo wants resilience without immediate whole-stack orchestrator capture

The strongest intent surfaces repeatedly support:

- multi-node ordinary Docker by default
- no central orchestrator by default
- local-first and peer-aware routing ambition
- lightweight registry and helper-agent direction

What that proves:

- whole-stack orchestrator promotion is not the current default answer

What it does not prove:

- that promotion will never be justified

### 2. The repo already knows where Compose is insufficient

The planning layer is explicit about missing capability classes:

- placement truth
- failover generation and route persistence
- convergence of secrets and runtime state
- reliable any-node entry behavior
- stateful correctness

What that proves:

- the platform-choice question is being asked from real pain, not curiosity
- the repo already understands the problem is about truth ownership, not just
  deployment polish

That second line is the real translation key.

The repo is not fundamentally trying to pick a deployment framework.
It is trying to decide who or what owns enough truth to make wrong-node,
failover, and convergence behavior stop depending on operator folklore.

That is why this page has to stay strict about fake differentiation.
Some options differ mostly in syntax and ecosystem.
Others differ in where truth ownership moves.
Only the second kind of difference really matters here.

What it does not prove:

- that all those pain classes should be solved by the same tool

### 3. OpenSVC-style HA remains attractive for infra-grade reasons

[`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
shows why infra-grade HA thinking remains attractive:

- dynamic ingress generation
- infra-style relocation and failover logic
- explicit L4/L7 boundary discipline

What that proves:

- some of the repo's hardest pain is infra-failover-shaped, not purely
  scheduler-shaped

What it does not prove:

- that OpenSVC is already the chosen global control plane

### 4. Stateful truth remains separate from scheduler choice

[`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
keeps forcing the same distinction:

- placement and supervision are one problem
- authority, election, replication, and storage truth are another

What that proves:

- no scheduler choice alone can honestly close the repo's HA story

What it does not prove:

- which scheduler or infra tool should handle the non-stateful layers

## Archive-pressure anchors

Representative archive clusters include:

- Compose frustration and dynamic alternatives
- Nomad versus Kubernetes comparisons
- Swarm versus Nomad and Swarm versus Kubernetes threads
- ClusterLabs or infra-grade HA exploration
- repeated no-hardcoded-node, no-fake-failover pressure

The sharpest archive examples for this page are:

- [`../source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`](../source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md),
  where the user explicitly narrows the problem to service discovery and
  wrong-node forwarding after manual placement and Cloudflare multi-A entry are
  already accepted.
- [`../source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`](../source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md),
  where the user asks for equal nodes, distributed reaction, and no-Swarm
  scaling of `docker-compose.yml` without building a full orchestrator from
  scratch.
- [`../source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`](../source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md),
  where the user rejects the usual "Traefik, NGINX, HAProxy" answer because
  those names alone do not close the actual failover and policy-preservation
  gap.

## What the archive proves

The archive keeps returning to the same refusal pattern:

- do not stop at DNS-only answers
- do not jump to a giant platform if the gain is not specific
- do not pretend scheduler choice alone solves ingress or storage truth
- do not confuse "more components" with "less sacred-node knowledge"

What that proves:

- the repo keeps multiple futures open because the problem is genuinely
  multidimensional
- the user is not randomly platform-shopping; they keep rediscovering the same
  missing truth layer under different tool names

That repeated rediscovery is exactly why the docs have to sound more like
reconstruction than recommendation.

The archive is not just a pile of comparisons.
It is a negative benchmark engine showing which answers keep sounding fresh
while leaving the same hidden wound intact.

What it does not prove:

- that indecision is cost-free
- that every preserved option deserves equal long-term attention

## The capability classes that should not be collapsed

The platform question only makes sense when the pain classes are separated.

### Class 1: authoring and inspectability

Where Compose is strongest:

- file-first changes
- human-readable service definitions
- ordinary diffs
- lower abstraction tax

### Class 2: placement and convergence

Where Compose starts to run out of honesty:

- what runs where
- how nodes converge on that truth
- how placement changes are reflected safely

### Class 3: ingress and failover semantics

Where helper layers, OpenSVC-style approaches, or stronger platforms become
relevant:

- any-node entry
- local-first or peer-forward correctness
- route persistence while services flap

### Class 4: stateful resilience

Where no scheduler alone is enough:

- write authority
- election and promotion
- client discovery
- storage durability and portability

Different platform choices solve different subsets of those classes.

## Strongest honest current answer

The shortest defensible platform answer today is:

> The live repo is still concretely Compose-first, the helper and coordination
> burden is already real enough that plain Compose is no longer the full story,
> several future control-plane families remain genuinely relevant because they
> address different missing truth classes, and no current evidence proves that
> one whole-stack future has already earned promotion over the others.

That answer is less emotionally satisfying than a winner.
It is also much closer to the truth the repo has currently earned.
That is why the decision surface remains open.

It is also why the repo should be suspicious of any page that makes the choice
sound singular too early.

The user's real frustration is partly that the outside ecosystem keeps bundling
separate burdens together:

- routing and wrong-node correctness
- service placement authority
- convergence and repair
- stateful authority and durability

Then it markets one product choice as if all of those must naturally move
together.

This page should keep refusing that collapse unless the evidence actually shows
that one promotion step pays down multiple burdens honestly instead of merely
renaming them under one logo.

## Evidence-backed reading of each path

### Compose-first plus helper layers

What the evidence proves:

- this is still the baseline strategy
- it preserves the operator contract the repo currently trusts
- it is the least-disruptive path for explicit placement truth, generated
  routing, and narrow convergence logic

Why it remains emotionally attractive:

- it keeps the current human authoring surface alive
- it preserves locality instead of burying it
- it tries to buy only the missing truths instead of importing a whole new
  worldview all at once

What the evidence does not prove:

- that this path has already become trustworthy enough to close the
  orchestration problem

The known `docker-gen-failover` route-loss issue is a major example of why
that stronger claim would be premature.

It is also the clearest warning that helper growth can become shadow
control-plane growth.
If Compose-plus-helpers ends up owning placement truth, convergence truth,
peer eligibility, and recovery behavior, then the repo has not escaped
orchestration logic.
It has only decomposed it.

And if the docs still call that state "basically still plain Docker," the docs
are helping hide the real complexity instead of reconstructing it.

That is one of the most important anti-fake-option observations in the repo.

If the helper path keeps expanding, then "we avoided orchestration" may become
less honest than "we built our own control plane slowly and do not want to name
it yet."

That sentence should remain one of the main promotion tripwires in the repo.
It is not anti-helper rhetoric.
It is the exact point where the docs have to stop pretending the burden
transfer is still cheap.

### OpenSVC or other infra-grade HA tooling

What the evidence proves:

- infra-grade takeover logic remains attractive for narrow critical domains
- the repo's ingress and failover pain is not purely scheduler-shaped

Why that matters:

- some of the user's deepest frustration is not "I cannot schedule enough
  containers"
- it is "reachable" keeps getting confused with "actually recoverable"
- infra-style HA tooling spends complexity on that narrower wound more
  directly than Kubernetes does

What the evidence does not prove:

- that infra-grade HA should become the dominant runtime for all services

That boundary matters because infra-grade HA can be the right answer to the
actual wound without being a universal answer to everything nearby.

### Nomad

What the evidence proves:

- Nomad remains attractive when placement, relocation, and restart semantics
  start mattering more than preserving a pure file-first helper path
- it offers a scheduler-shaped answer without instantly importing full
  Kubernetes worldview capture

Why it stays alive in the repo:

- it looks like the least insulting scheduler-shaped answer
- it may admit that placement and restart authority need a real owner without
  erasing the repo's desire for a more legible operator surface

That is why Nomad remains emotionally important here.
It represents a possibility that the repo could promote into something more
explicit without being captured by the heaviest worldview first.

What the evidence does not prove:

- that Nomad already wins the repo's ingress semantics or stateful truth
  problem

The honest question for Nomad is not whether it is lighter than Kubernetes.
It is whether the repo is already paying scheduler-like tax in a more fragile,
less self-aware form.

### k3s or Kubernetes

What the evidence proves:

- stronger scheduler-backed and controller-backed futures remain plausible
- the repo has kept those futures open because some domains may eventually earn
  that tax

What the evidence does not prove:

- that this tax has already been justified for the whole stack

Kubernetes only earns itself here if the docs can explain something stricter
than:

- it is powerful
- it is common
- it eventually solves lots of things

It has to explain which hidden human SPOFs, wrong-node failures, or
convergence burdens it removes more honestly than the thinner alternatives.

### Swarm

What the evidence proves:

- Docker-native cluster temptation is obvious and recurrent

What the evidence does not prove:

- that Docker-native flavor alone answers the specific failure classes the repo
  keeps fighting

## The hidden threshold this evidence keeps pointing at

The repo is very close to a harder question:

> if the helper path ends up owning placement truth, convergence truth,
> route-generation truth, peer eligibility truth, and failover truth, has the
> repo effectively built a custom control plane around Compose already?

That question matters because it becomes the real threshold for promotion.

If the custom helper layer stays narrow and auditable, it may remain the right
answer.

If it keeps expanding into scheduler-like responsibilities, then selected
promotion into Nomad, k3s, Kubernetes, or infra-grade HA may finally earn its
keep.

That promotion still should not be narrated as "the platform decision is over"
unless the repo can also say which burdens were *not* solved by that
promotion.

Otherwise the docs just recreate the same problem at a more sophisticated
level:

- one path gets described as adult and settled
- unresolved burdens are pushed out of frame
- the operator quietly keeps carrying them anyway

## The user's real benchmark for "earned promotion"

A future platform earns promotion here only if it demonstrably reduces at least
one of the repo's real hidden taxes:

- sacred-node memory
- stale placement facts
- route loss under backend failure
- private operator reconstruction during wrong-node traffic
- uncontrolled growth of helper logic into a shadow control plane

If a platform mostly adds:

- prestige
- fashionable vocabulary
- new moving parts
- broader abstraction

without materially reducing those taxes, then it has not earned promotion in
this repo, regardless of how standard it sounds elsewhere.

This is probably the harshest useful rule in the page.

Standardization is not the same thing as alignment.
Popularity is not the same thing as healing the correct wound.

## Bottom line

This page does not prove a winner.
It proves a stricter reading rule:

- the repo has real control-plane pressure
- the repo still values the Compose-first operator contract
- the user is refusing fake HA and fake closure
- the remaining decision must stay tied to named truth gaps and named failure
  classes

That is the main outcome this page is protecting:

the repo keeps the option space open only where the evidence says the wound is
still unresolved, not because indecision is comfortable, and not because every
tooling path deserves the same respect forever.

That is the actual tradeoffs evidence.
Anything softer turns the platform story back into architecture theater.
