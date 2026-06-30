# Capability Gaps and Roadmap

For the evidence shaping this roadmap, start with:

- [`../research/evidence-ledger.md`](../research/evidence-ledger.md)
- [`../research/ingress-and-failover-evidence.md`](../research/ingress-and-failover-evidence.md)
- [`../research/stateful-ha-evidence.md`](../research/stateful-ha-evidence.md)
- [`../research/orchestrator-tradeoffs-evidence.md`](../research/orchestrator-tradeoffs-evidence.md)
- [`/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This page is not a feature backlog.

It is the map of what still makes the current architecture partly fake relative
to the user's real benchmark.

The user is not asking for a tidy next-steps board.
They are asking why so many infrastructures offer the feeling of options while
quietly forcing one of two humiliations:

- stay in ordinary Docker forever and accept wrong-node fragility
- jump into a heavyweight control plane before smaller honest layers were even
  exhausted

So the only roadmap worth keeping here is the roadmap that answers:

> which missing truth or recovery layer is still forcing the rest of the system
> to overclaim?

That is the meaning of "roadmap" in this repo.

## What this page is and is not allowed to prove

This page is allowed to:

- prioritize missing truth layers in the order that protects honesty
- explain which gaps still force overclaim if left unresolved
- connect runtime evidence to the next proof-bearing promotion step
- distinguish valuable next work from already-proven capability

This page is not allowed to:

- act like a completion report
- imply that a well-ordered roadmap means the runtime is already coherent
- flatten HTTP, TCP, stateful, and control-plane lanes into one generic queue
- upgrade planning clarity into present-tense implementation maturity

## Strongest honest current answer

The next work is not "whatever seems generally useful for self-hosting." The
next work is whichever missing truth layer still forces the docs to rely on
hope, operator memory, stale topology assumptions, or rhetorical glue. In the
priority implementation today, that means placement truth, convergence truth,
peer eligibility, route durability under real failure, and keeping stateless,
protected, TCP, and stateful promotion under separate honesty gates.

## The shortest honest roadmap

The current best order remains:

1. keep the docs hostile to fake certainty
2. establish placement truth
3. establish convergence truth
4. make routes survive local backend loss
5. define peer eligibility as something stricter than reachability
6. prove one full wrong-node stateless HTTP path end to end
7. prove protected-route continuity on that same path
8. keep HTTP, TCP, and stateful classes under separate proof rules
9. only then decide what has earned promotion into a stronger control layer

This is not a motivational sequence.
It is the dependency chain between the dream and reality.

## Why the order is this strict

The project wants all of these things at once:

- multiple public nodes
- local-first service handling
- peer fallback when the request lands on the wrong node
- no fake HA language
- no forced Swarm or Kubernetes jump before necessary

That dream depends on a stack of truths.
If the lower truths are weak, the upper truths become performance.

The dependency chain is:

1. know where services actually live
2. know whether peers are semantically aligned enough to substitute
3. keep the recovery route alive when the local backend dies
4. know which peers are eligible to receive traffic now
5. prove that policy, auth, and middleware survive the handoff
6. keep stateful classes under harsher rules than stateless HTTP

If that chain breaks at step 1 or 2, later success claims are already suspect.

## What this roadmap is actually protecting

This roadmap is not protecting a tool preference.
It is protecting the user's real demand:

> stop pretending there are plenty of options when most of the supposed options
> either collapse under wrong-node pressure or force a giant control plane
> before the smaller honest answers have even been exhausted

That is why the roadmap stays narrow and harsh.

The next thing to build is not the next thing that sounds advanced.
It is the next thing whose absence still forces the stack to lie.

This also means the roadmap must actively reject seductive but misordered moves:

- promoting to a bigger platform before placement and convergence truth are
  explicit
- narrating ingress sophistication as if it settles stateful correctness
- treating helper growth as "still just Compose" after helpers start owning
  scheduler-like truth silently

## Read this as integrity gates, not aspirations

Each priority below is written in the format:

- current truth
- failure signature
- proof threshold
- what it unlocks

That format matters because the repo does not need more architecture dreams.
It needs clearer gates for when a bigger claim becomes legal.

## Priority 0: documentation honesty

Class:

- integrity gate

Why it is first:

If the docs overclaim, every later decision gets poisoned.
This repo cannot afford one more smooth narrative that upgrades intent into
"basically working."

Current truth:

- the knowledgebase is much stricter than the earlier docs
- it can still regress whenever a page starts sounding like a generic HA guide
  instead of a pressure-tested reading of the worktree and archive

Failure signature:

- DNS plurality described as end-to-end failover
- helper layers described as live truth before runtime proof exists
- stateful systems described as HA because they are reachable
- route generation described as recovery without route-persistence proof

Proof threshold:

- every operator-critical page distinguishes live runtime, architecture intent,
  planning pressure, and archive pressure
- every major resilience claim makes its proof class obvious
- every weaker state stays named as weak instead of being widened into
  "solved"

Unlocks:

- trustworthy sequencing
- less self-deception
- cleaner promotion decisions later

## Priority 1: placement truth

Class:

- foundational

What this really means:

The receiving node needs a current answer to:

> is the requested service local, and if not, where does it actually live right
> now?

Without that answer the whole dream collapses into guesswork.

Current truth:

- the repo repeatedly converges on `services.yaml` or an equivalent current-
  state surface as the intended lightweight registry
- the tracked priority runtime still does not prove a live root
  [`services.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/services.yaml)
  exists and is consumed by routing

Failure signature:

- traffic lands on a healthy node
- the node cannot answer service placement deterministically
- peer-forward behavior depends on stale assumptions or private operator memory

Proof threshold:

- one explicit auditable placement-truth surface exists
- route generation or runtime logic demonstrably consumes it
- node-local versus remote service identity can be derived from current truth,
  not inference

What counts as evidence:

- a tracked root artifact or equally explicit generated state surface
- documented update semantics
- proof that runtime or generators actually read it

Unlocks:

- honest local-first logic
- non-theatrical wrong-node handling
- meaningful global-hostname semantics

This is why placement truth comes before almost every glamorous proposal.
Without it, later options differ mostly in presentation while still depending
on the same private reconstruction.

## Priority 2: convergence truth

Class:

- foundational

What this really means:

Even if a node knows where the service "should" be, fallback is still fake if
the candidate peer is semantically out of alignment.

Current truth:

- the master plan explicitly names secret sync and Compose sync as missing
  layers
- current node alignment is still partly social and operational rather than a
  live shared truth surface

Failure signature:

- a peer is reachable
- the peer answers
- but the route still lands on different env, secrets, middleware, or revision
  semantics than the operator assumed

Proof threshold:

- a node can prove that the peer it is about to use is on an acceptable
  revision and secret surface for the relevant route

What counts as evidence:

- visible revision or generation markers
- documented convergence rules
- drift detection or equivalence checks that actually gate eligibility

Unlocks:

- honest peer substitution
- cleaner failover drills
- less social coordination burden during recovery

Without convergence truth, wrong-node recovery is still partly luck wearing a
better badge.

## Priority 3: route persistence under backend loss

Class:

- routing integrity gate

What this really means:

The route required for recovery has to survive the failure that made recovery
necessary.

Current truth:

- `docker-gen-failover` is materially live in the edge stack
- the master plan explicitly records that the current model can delete routes
  when a container stops

Failure signature:

- the fallback path appears to exist while the backend is healthy
- the backend disappears
- the route evaporates exactly when it is needed

Proof threshold:

- at least one real backend-loss drill proves the route remains present and
  usable

What counts as evidence:

- pre-failure route identity
- intentional backend stop or failure
- observation that the route stayed present
- user-visible success through the surviving path

Unlocks:

- the first real claim that "fallback" means more than "there was another
  idea on paper"

This priority is narrower than "fix failover."
It is stronger than "dynamic config exists."

## Priority 4: peer eligibility truth

Class:

- routing decision gate

What this really means:

The receiving node needs something stricter than reachability when choosing a
peer.

Current truth:

- Headscale is materially live
- the master plan talks about sync-agent, health reports, peer broadcast, and
  service state
- the current runtime does not yet prove a narrow trusted eligibility answer
  that the receiving node can consume

Failure signature:

- the platform knows some peers exist
- it cannot prove which one is correct now for the requested service

Proof threshold:

- one auditable peer-selection decision based on current tracked truth rather
  than folklore

What counts as evidence:

- visible eligibility inputs
- visible selection rule
- route decision trace or equivalent proof that the chosen peer came from the
  shared truth surface

Unlocks:

- honest wrong-node forwarding
- better separation between "reachable" and "safe"

## Priority 5: one real wrong-node stateless HTTP proof

Class:

- proof gate

What this really means:

The repo needs at least one route that stops being a theory.

Current truth:

- the priority runtime has a serious enough edge stack to make this worth
  proving
- the docs already isolate the exact claim sharply
- no generic wrong-node proof is currently claimed

Failure signature:

- the architecture sounds convincing
- no single route has been intentionally exercised through the wrong healthy
  node and traced end to end

Proof threshold:

- intentionally land one stateless HTTP route on the wrong node
- prove receiving-node identity
- prove backend-node identity
- prove user-visible success

Unlocks:

- one genuinely defensible "wrong-node HTTP works" claim
- a real baseline for backend-loss drills

This is the first place the repo can make the user’s central pain materially
smaller rather than merely better described.

## Priority 6: protected-route continuity

Class:

- semantic integrity gate

What this really means:

For protected routes, transport is not enough.
The forwarded path must still behave like the same protected service.

Current truth:

- TinyAuth, Nginx auth extensions, CrowdSec, and Traefik middleware are all
  materially live
- the docs clearly name policy continuity as part of routing correctness
- no route-specific parity proof is currently claimed

Failure signature:

- the route still answers
- but auth, middleware, headers, or visible policy meaning diverge after
  peer handoff

Proof threshold:

- compare one protected route locally versus through intentional peer handoff
- verify auth challenge, middleware, headers, and visible policy parity

Unlocks:

- the first honest claim that a peer-forwarded request remained the same
  protected route, not just a successful shortcut

## Priority 7: keep HTTP, TCP, and stateful lanes separate

Class:

- anti-inflation gate

What this really means:

One HTTP success must not counterfeit the rest of the platform.

Current truth:

- the root runtime already contains HTTP, protected HTTP, TCP, and stateful
  surfaces
- the docs already know these are different lanes
- the biggest remaining risk is emotional overread after a narrower win

Failure signature:

- a stateless HTTP route succeeds
- the docs start sounding as if TCP and stateful systems are morally adjacent

Proof threshold:

- separate proof matrices and route language remain in force after narrower
  wins

Unlocks:

- honest maturation without platform-wide bluff

## Priority 8: only then consider stronger control-layer promotion

Class:

- promotion gate

What this really means:

Bigger platforms should be justified by a concrete burden they remove, not by
prestige or generic best practice.

Current truth:

- the repo has real exploration around OpenSVC, Nomad, k3s, and other control
  surfaces
- the current pain is still more sharply about wrong-node truth and hidden
  memory than about generic scheduler scale

Failure signature:

- a heavier platform is chosen before the narrower missing layers were made
  explicit and partially tested

Proof threshold:

- the repo can say exactly which unresolved burden the promoted layer removes
  that smaller layers could not remove honestly

Unlocks:

- a stronger control plane that has actually earned trust instead of merely
  arriving first

## The roadmap in blunt English

The repo does not yet need "the best orchestration story."
It needs the smallest sequence that makes the wrong-node path less humiliating,
the fallback path less ceremonial, and the operator's private topology memory
less central.

That is why the roadmap is not:

- add more helpers
- add more nodes
- add more dashboards
- pick the most respected cluster

It is:

- make truth explicit
- make substitution honest
- make fallback survive the real failure
- prove one route for real
- refuse to let that one proof inflate into platform-wide closure
