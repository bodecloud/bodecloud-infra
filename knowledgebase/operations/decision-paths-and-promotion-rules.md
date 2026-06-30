# Decision Paths and Promotion Rules

This page exists to answer the question that still gets left too loose even
after the architecture is explained honestly:

> given the user’s actual dream, what are the real paths forward, what does
> each path genuinely buy, what does each path still refuse to solve, and when
> has a stronger control plane actually earned promotion?

Without this page, the repo can describe the dream and the gaps accurately and
still leave the operator with the same old frustration:

- lots of interesting pieces
- no blunt mapping between pain and next layer
- no clean rule for when Compose should stay central
- no clean rule for when promotion is no longer optional

This page is supposed to fix that.

It is also supposed to stop the repo from responding to the user's frustration
the way the broader ecosystem keeps responding to it:

- by offering more categories
- by offering more products
- by offering more maturity theater
- while leaving the same hidden burden intact

That last point is the whole reason this page must stay stricter than an
ordinary architecture comparison.

The user is not mainly frustrated by lack of names.
The user is frustrated by being offered many names for paths that still do not
answer the real wound.

This page therefore has to do more than rank next steps.
It has to stop the repo from repeating the broader ecosystem pattern:

- rename the pain
- package the pain differently
- call the new package a serious option
- leave the operator privately completing the architecture on the bad day

If a path still ends there, then for this repo it is not a real option yet.

## What a useful decision page has to prevent

The user’s real frustration is not lack of information.
It is too many options that are not actually options.

This page has to keep the repo from making that worse by offering paths that:

- sound distinct
- use different tooling names
- but still leave the same hidden sacred-node memory in place

If two paths leave the same truth gap unresolved, then for this repo they are
much closer to the same path than to different ones.

That means this page has to behave like a filter for fake options.

Different packaging does not create a meaningfully different path if all of
the following stay true:

- placement truth still lives in human memory
- wrong-node entry still depends on luck or ad hoc glue
- peer eligibility is still socially reconstructed instead of runtime-legible
- fallback still means "some response happened" rather than "the same service
  contract survived"

That rule is stricter than normal architecture comparison, and it needs to be.

The user is not lacking lists of tooling options.
The user is trying to escape a world where many differently branded answers
still collapse back into:

- one remembered placement truth
- one sacred ingress truth
- one operator-only failure narrative

That collapse is the main anti-benchmark for this whole page.
If a recommendation changes tooling, controller shape, or platform vocabulary
without changing the hidden truth the operator still has to carry, then the
recommendation is mostly theater from the user's point of view.

## Read this page correctly

This is not a product ranking.

It is a routing table for decisions.

Each path below is organized around a pain class:

- what hurts now
- what layer should own that pain
- what that layer really solves
- what that layer still does not solve
- what evidence would justify promoting to the next layer

That keeps the repo from doing the usual thing:

calling in a bigger platform before the exact missing truth has been isolated.

This page therefore treats decision-making as an evidence-routing problem.

The right question is not:

> what is the most impressive next layer?

It is:

> what exact hidden burden is still not owned honestly enough by the current
> layer, and what is the smallest promotion that would really own it?

The key word is "really."
This repo does not need more semi-answers that sound structurally different
while preserving the same semantic gamble around wrong-node traffic, backend
loss, or stateful authority.

That question should be read together with the repo's strongest architecture
intent surface:
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).

It is the clearest place where the dream is stated as:

- Compose-first
- multi-node
- local-first then peer-forward
- anti-SPOF pressure
- without immediately collapsing into heavyweight scheduler capture

That is why this page should not over-weight [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
or [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)
when deciding architectural promotion.

Those files matter, but they are not the primary expression of the
distributed-systems dream.

## The shortest possible answer

Right now the most evidence-aligned default still is:

1. stay Compose-first
2. close the middle-layer truth gaps
3. prove one real wrong-node HTTP path
4. prove backend-loss route persistence
5. only then promote specific domains that still hurt in ways
   Compose-plus-helpers cannot solve honestly

That default is not caution for its own sake.
It is the only path that matches the repo’s current state and the user’s real
refusal pattern.

The refusal pattern matters.

The user is not refusing orchestration because orchestration is hard.
The user is refusing premature worldview capture that renames the same old gap
instead of closing it.

That is why this page should feel biased against promotion theater.
The user is not asking for the next step that sounds mature.
The user is asking for the next step that makes the wrong-node and hidden-truth
problem less true.

## The deepest decision rule

The repo should not promote because a bigger tool exists.
It should promote because the current layer has become the hidden tax.

That means:

- do not promote because cluster sounds more mature
- do not promote because the edge stack is already sophisticated
- do not promote because stateful systems feel scary in the abstract
- do not promote because helper logic looks inelegant from a distance

Promote only when the current layer is now forcing the operator to carry too
much hidden truth in their head, or is itself becoming a disguised controller
without earning the honesty of naming that fact.

That is the harshest and most useful decision rule in the repo:

promote when the current layer is lying by omission, not when a shinier layer
exists.

There is a companion rule:

do not preserve the current layer out of sentimentality either.

If Compose-plus-helpers stops being the least dishonest surface and becomes a
socially maintained pseudo-control-plane that only still sounds simpler, then
"staying Compose-first" can become just as fake as premature platform capture.

This can be phrased even more bluntly:

promote when a path becomes fake by leaving the same burden unnamed.
Do not promote merely because a new platform lets the docs sound calmer.

That is the decision rule most likely to preserve the real dream.

Without it, "more options" quickly becomes another way of saying
"more substitutes that do not actually answer the user's question."

## Path 1: Compose-first plus explicit truth helpers

Use this path when:

- manual service placement is still acceptable
- Compose readability is still buying real value
- the dominant pain is wrong-node routing and drift, not scheduler-scale churn
- the operator wants to preserve direct file-level authorship

This path means:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  and `compose/docker-compose.*.yml` remain the priority implementation
- a narrow control surface is added behind them
- that surface owns truth, routing generation, and convergence checks rather
  than whole-platform scheduling

This is the path most aligned with the user's deepest request:

not "never build a control plane,"
but "build the smallest truth-owning layer that actually stops the fake-option
cycle."

That last clause has to stay sharp.
If this path produces a placement file, a sync helper, a route generator, and a
convergence checker, but the operator still has to mentally decide which one is
the real truth on the bad day, then the path has not escaped the fake-option
cycle.
It has only organized it better.

### What this path must contain if it is going to stop being hand-wavy

#### A. One explicit placement-truth surface

The repo keeps pointing at `services.yaml` for a reason.

That surface needs to answer:

- which service identity exists
- which node currently hosts it
- whether it is node-local, replicated, or globally available
- what ports and protocols matter

#### B. One convergence-truth surface

The repo needs auditable node state for at least:

- Compose or repo revision
- secret and env freshness
- auth and middleware dependencies
- deployment freshness

#### C. One route-generation surface that survives failure

The first-hop proxy cannot depend on route generation that disappears when the
local backend disappears.

That is why replacing or hardening `docker-gen-failover` matters so much.

#### D. One explicit peer-eligibility model

It must be possible to explain:

- why a peer was eligible
- why a peer was excluded
- what health definition was actually applied

#### E. One proof-first stateless HTTP drill

This path is not real until one representative service proves:

- any-node first hop
- wrong-node forwarding
- backend-loss route survival
- middleware and auth continuity

This requirement is deliberately narrow because the user is tired of vague
“supports failover” claims.

One route proven honestly is worth more in this repo than twenty abstract
features that still leave the real path unexercised.

That asymmetry matters here.
A stack with more features is not ahead of the user's dream if the bad-day
request path is still narratively stronger than operationally owned.

### What this path solves well

- keeps Compose central
- preserves readability
- solves the current truth gaps first
- matches the repo’s strongest instruction surfaces

### What this path still does not solve by itself

- stateful election and write ownership semantics
- generic auto-placement
- large-scale reconciliation logic
- universal rescheduling after node loss

Those are not failures of the path.
They are its honest boundaries.

Those boundaries are valuable, because an honest partial answer is exactly what
the repo is trying to preserve against ecosystems that oversell broad ones.

### When this path stops being enough

Promote beyond this path when one of these becomes the dominant pain:

- manual placement itself is now a bigger tax than the helper layer
- service churn is too high for explicit placement plus helper truth
- the helper layer is silently becoming an orchestrator anyway
- stateful failover needs deeper topology ownership than helper agents can
  honestly provide

If those are not true yet, stay here.

This page should be read as slightly biased toward staying here until the pain
becomes concrete.

That bias is not conservatism for its own sake.
It exists because the repo still has not fully earned a broader promotion, and
because broad promotion too early would erase the exact missing layer the user
is trying to isolate.

That is why this path is not merely "the conservative option."
It is the path most likely to preserve explanatory fidelity while the repo
still learns what the actual missing truth layer has to do.

But staying here is only justified while the path is still getting more
truthful.
If the helper layer keeps growing while the same core questions still resolve
to operator reconstruction, then this path stops being disciplined and starts
being evasive.

## Path 2: infra-grade HA promotion for ingress or identity-critical pieces

Use this path when:

- the most painful problem is not app scheduling
- the most painful problem is keeping a small number of critical edge or
  identity surfaces alive
- the operator wants boring HA primitives before broad platform capture

Typical candidates:

- ingress identity continuity
- VIP-like behavior
- critical reverse-proxy presence
- tightly-scoped infra roles where failover logic is narrower than full app
  scheduling

This is where OpenSVC-, keepalived-, or HAProxy-style promotion becomes much
more reasonable than full scheduler promotion.

This path is important because it is one of the few honest “more tooling”
answers that does not pretend every problem is an app-scheduling problem.

It also matters psychologically.

One source of the user's frustration is being forced into giant solution
domains when the real pain may still be much narrower:

- keep first-hop truth alive
- keep identity surfaces boring
- stop one small critical layer from staying sacred

The limit is important:
removing a sacred first hop is not the same as removing sacred application
truth.
This path earns itself when the narrow pain is actually narrow.
It becomes fake when its first-hop improvements get narrated as if the whole
service surface became request-preserving.

### What this path solves well

- strong first-hop or critical infra failover semantics
- boring HA for narrow domains
- explicit ownership of infra-grade liveness

### What this path still does not solve by itself

- application placement truth for the whole stack
- generic wrong-node request preservation for every service
- stateful application correctness
- broad developer workflow unification

This path is excellent when the repo needs infra-grade resilience for a narrow
slice, not when it needs a universal application control plane.

## Path 3: lighter scheduler promotion

Use this path when:

- the dominant pain has truly become placement, rescheduling, and deployment
  automation
- the operator still wants something lighter and more legible than Kubernetes
- the repo is ready to pay some scheduler tax in exchange for real lifecycle
  ownership

This path says:

- the middle layer taught us where the pain really is
- now a real workload scheduler is justified for some or all stateless
  services
- the repo still wants a simpler operator contract than a full Kubernetes
  ecosystem

The key phrase is “for some or all stateless services.”

If this path starts being narrated as a universal answer to stateful truth,
shared storage, and protocol-specific failover, the docs are drifting again.

This path also fails if it gets promoted mainly because the helper mesh became
embarrassing to explain.
Embarrassment is not evidence.
The scheduler should be promoted because it owns a named pain more honestly,
not because the existing story became hard to narrate cleanly.

### What this path solves well

- workload placement automation
- service lifecycle ownership
- rescheduling after failure
- a clearer answer to where should this run

### What this path still does not solve automatically

- stateful correctness
- storage ownership
- ingress semantics unless integrated carefully
- honest multi-protocol failover policy

The repo should only promote here after it can say:

- the helper-layer truth gaps are understood
- scheduler promotion is solving a named pain, not vaguely modernizing the
  stack

That last sentence is one of the main anti-theater rules in the repo.

## Path 4: full desired-state platform promotion

Use this path only when:

- the repo truly wants desired-state reconciliation as a first-class contract
- operator tax is justified by the ecosystem gained
- service churn, controllers, storage operators, and platform-wide scheduling
  genuinely outweigh the Compose-first legibility loss

This path is not forbidden.
It is just the highest-tax option in the repo’s universe.

That tax is not only operational.
It is also paid in:

- lost directness
- farther-away truth surfaces
- more controller-shaped explanations
- higher risk that the docs trust the platform story more than the bad-day
  evidence

### What this path solves well

- desired-state scheduling and reconciliation
- broad operator ecosystem
- mature controller patterns
- strong platform conventions

### What this path still does not magically solve

- the need for honest stateful topology design
- cross-domain anti-SPOF storytelling
- the difference between ingress continuity and application correctness
- the need to define what the service still means under failure

This repo should never promote here just because multi-node Docker feels
annoying.
That is not a sufficient reason.

The earned reason would have to sound much more specific:

- explicit truth layers are understood
- helper composition has already grown control-plane-sized
- desired-state reconciliation is now solving a named dominant pain
- the operator is actually getting clarity back, not just ecosystem gravity

This page should keep treating "clarity back" as the key requirement.

If promotion makes the stack more fashionable but not more self-explaining,
then the user has not actually gained an option.

That sentence should be treated as a hard gate, not a nice side effect.
The user's complaint is fundamentally about hidden truth and fake options, not
about lack of ecosystem prestige.

## Path 5: service-class-specific stateful promotion

Use this path when the real pain is not generic routing at all.

Examples:

- Redis needs Sentinel, Cluster, or protocol-aware proxying
- MongoDB needs replica-set truth and client rediscovery
- RabbitMQ needs explicit clustering or an honest single-node boundary
- shared storage needs a real answer about locking, ownership, and durability

This path matters because stateful systems are exactly where infrastructure
language becomes most misleading.

This is also where the repo most needs permission to be uneven.

It is acceptable for:

- stateless HTTP to mature first
- Redis to get a sharper answer than RabbitMQ
- MongoDB to be promoted differently from Postgres

Uniformity is not the goal.
Honest service-class-specific correctness is.

This page should stay willing to look uneven if unevenness is the more honest
state.
The user wants a truthful map of what is real, weak, or still fictionalized,
not a symmetrical platform story.

This is one of the biggest places where generic platform language harms the
repo.

The user does not need every subsystem to mature at the same rate.
The user needs each subsystem to stop pretending it is more solved than it is.

The repo should prefer:

- narrower honest promotion per stateful class

over:

- one universal claim that the platform is now HA

## The actual decision map

Use this blunt mapping.

It is intentionally blunt because the user has already seen what happens when
tooling guidance is too willing to blur adjacent pains together.

### If the pain is:

> requests land on the wrong node and the node cannot figure out where to send
> them

Use:

- Path 1 first

Do not jump directly to:

- Path 4

unless the missing truth layer has already proved too complex to keep narrow.

### If the pain is:

> one or two critical ingress or identity surfaces still feel sacred even
> before app scheduling becomes the main problem

Use:

- Path 2

This is exactly where infra-grade HA primitives can earn their keep without
capturing the whole stack.

### If the pain is:

> stateless service placement and rescheduling are now a bigger tax than manual
> placement plus helper truth

Use:

- Path 3

This is the point where a lighter scheduler starts to solve a real problem
instead of merely sounding modern.

### If the pain is:

> the repo genuinely wants a full desired-state platform with controller
> semantics and ecosystem depth

Use:

- Path 4

But only after saying explicitly which pain remained unsolved by Paths 1, 2,
or 3.

### If the pain is:

> a specific stateful service still has fake HA language around it

Use:

- Path 5

This repo should promote stateful truth per service class long before it claims
stateful anti-SPOF globally.

## The promotion rules

The repo should only promote upward when the lower layer has become the hidden
tax itself.

That means:

- do not promote because the bigger platform exists
- do not promote because cluster language sounds cleaner than helper logic
- do not promote because the edge stack is already sophisticated
- do not promote because stateful systems are intimidating in the abstract

Promote only when the current layer is now hiding too much truth or demanding
too much manual compensation.

If the helper layer becomes opaque, coordination-heavy, and difficult to audit,
then the honest move may be to promote it into a named scheduler instead of
pretending it is still somehow simpler because it has fewer logos.

The inverse is also true.
If a promoted platform still leaves the operator carrying the same private
wrong-node, route-survival, or stateful-authority knowledge, then the platform
has not actually earned its promotion in this repo even if it is powerful in
general.

That is another way of saying:

- small is not automatically honest
- big is not automatically dishonest
- the real question is which layer now best exposes and owns the truth the
  operator would otherwise have to carry privately

## The honest current recommendation

The current best recommendation is still:

- keep Compose central
- build the missing middle layer honestly
- prove one real stateless HTTP wrong-node path
- prove backend-loss route survival
- keep stateful promotion narrower and later
- let stronger control planes earn themselves by named pain class

That is the route most aligned with the user’s actual dream.

It is also the route least likely to destroy the evidence the repo still needs
in order to understand itself honestly.

## Bottom line

The repo does not need a bigger answer first.
It needs a more precise one first.

The right promotion rule is:

> strengthen only the layer that is currently forcing the operator to carry too
> much hidden truth in their head

That is how the repo avoids replacing one kind of architecture theater with a
larger one.
