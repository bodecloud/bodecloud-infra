# Capability Gaps and Roadmap

For the evidence shaping this roadmap, start with:

- [`../research/evidence-ledger.md`](../research/evidence-ledger.md)
- [`../research/ingress-and-failover-evidence.md`](../research/ingress-and-failover-evidence.md)
- [`../research/stateful-ha-evidence.md`](../research/stateful-ha-evidence.md)
- [`../research/orchestrator-tradeoffs-evidence.md`](../research/orchestrator-tradeoffs-evidence.md)

This page is not a feature roadmap.

It is the map of what still makes the current architecture partly fake.

That wording is intentional.
The user is not asking for a tidy backlog.
They are asking why so many infrastructure stacks offer the feeling of options
while quietly forcing a terrible choice:

- stay in small-node Docker forever and accept wrong-node fragility
- or jump straight into a large control plane that demands allegiance before it
  has earned trust

The only roadmap worth keeping here is the one that answers:

> which missing truth or recovery layer is still forcing the rest of the system
> to overclaim?

That means this page should be read less like "what do we build next?" and
more like "what do we have to make impossible to keep lying about next?"

If that sounds harsher than a normal roadmap, that is because the user is not
asking for a normal roadmap.
The user is asking for an exit from fake options, fake HA language, and hidden
topology memory.

If a roadmap item does not reduce overclaim, hidden operator burden, or request
loss under real failure, it is secondary.

That standard is intentionally harsher than an ordinary platform roadmap.

The user is not asking for "the next sensible improvements."
The user is asking for the exact sequence that makes the system stop collapsing
back into fake options, remembered placement, and wrong-node ambiguity.

## What this page is and is not allowed to prove

This page is allowed to:

- prioritize missing capability layers in the order that protects honesty
- explain which gaps still force overclaim if left unresolved
- connect design intent to promotion work that would make the repo more true
- distinguish "valuable next work" from "already proven runtime behavior"

This page is not allowed to:

- serve as a completion report
- imply that a neat roadmap means the architecture is already coherent in
  runtime
- flatten distinct gaps across HTTP, TCP, and stateful systems into one generic
  backlog
- turn documentation clarity into false maturity claims

## Quick claim router

If the question is:

- "What should happen next if we care about honest anti-SPOF progress?" this
  page is a primary answer.
- "Which missing layer still makes the docs overstate reality?" this page
  should call that out directly.
- "Does having this roadmap mean the repo is already operationally converged?"
  no.
- "Can one backlog item solve the whole dream?" no. The roadmap exists because
  several independent truth layers still need promotion.

## The shortest honest roadmap

The current best order remains:

1. keep the docs hostile to fake certainty
2. establish placement truth
3. establish convergence truth
4. make routes survive local backend loss
5. define peer eligibility as something stricter than reachability
6. prove one full wrong-node HTTP path end to end
7. keep HTTP, TCP, and stateful classes separate
8. only then decide what has actually earned promotion into a stronger control
   layer

This is not a motivational sequence.
It is the dependency chain between the dream and reality.

That is why the order may look narrower than a normal engineering backlog.
It is optimized around recovering the missing truth layer, not around producing
the most impressive-looking milestone first.

## Strongest honest current answer

The next work is not "whatever seems generally useful for self-hosting." The
next work is whichever missing truth layer still forces the docs to rely on
hope, operator memory, or rhetorical glue. Today that usually means placement
truth, peer eligibility and convergence, route durability under wrong-node
conditions, and keeping HTTP, TCP, and stateful promotion in separate honesty
lanes so one partial win does not counterfeit the rest.

## Why the order is this strict

The project wants all of these things at once:

- multiple public nodes
- local-first service handling
- peer fallback when the request lands on the wrong node
- no fake HA language
- no forced Swarm or Kubernetes jump before necessary

That dream depends on a stack of truths.
If the lower truths are weak, the upper truths become performance.

This is one of the most important constraints in the whole knowledgebase.

The repo already has enough moving parts that it can perform coherence well.
The roadmap exists to stop performed coherence from being mistaken for earned
behavior.

The dependency chain is:

1. know where services actually live
2. know whether nodes are semantically aligned enough to substitute
3. keep the recovery route alive when the local backend dies
4. know which peers are eligible to receive traffic
5. prove that policy, auth, and middleware survive that handoff
6. keep stateful classes under stricter rules than stateless HTTP

If that chain breaks at step 1 or 2, later success claims are already suspect.

## What this roadmap is actually trying to protect

This roadmap is not protecting a tool preference.
It is protecting the user’s real demand:

> stop pretending there are plenty of options when most of the supposed options
> either collapse under wrong-node pressure or force a giant control plane
> before the smaller honest answers have even been exhausted

That sentence is the real meaning of "roadmap" in this repo.

Not:

- which thing should be built next because it is common
- which thing should be built next because it sounds mature

But:

- which missing truth is still forcing the rest of the stack to lie

That is why the roadmap stays narrow and harsh.

The next thing to build is not the next thing that sounds advanced.
It is the next thing whose absence still forces the stack to overclaim.

This also means the roadmap must actively reject certain seductive but
misordered moves:

- promoting to a bigger platform before placement and convergence truth are
  explicit
- narrating ingress sophistication as if it closes stateful correctness
- treating helper growth as "still just Compose" after it starts owning
  scheduler-like truth silently

## Read this as integrity gates, not aspirations

Each priority below is written in the format:

- current truth
- failure signature
- proof threshold
- what it unlocks

That format matters because the repo does not need more architecture dreams.
It needs clearer gates for when a bigger claim becomes legal.

The page should therefore be read as a falsification queue.

Each priority is trying to remove one remaining excuse the stack currently has
for sounding more complete than it is.

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
  and archive pressure
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

- the repo repeatedly converges on `services.yaml` or an equivalent state
  surface as the intended lightweight registry
- the tracked root runtime still does not prove a live root `services.yaml`
  exists and is consumed by the priority implementation

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

This is also why placement truth comes before almost every glamorous proposal.
Without it, many later options differ mostly in presentation while still
depending on the same private operator reconstruction.

## Priority 2: convergence truth

Class:

- foundational

Why it matters:

A request surviving transport does not mean the platform preserved meaning.

If forwarded traffic reaches a peer with the wrong secrets, wrong revision,
wrong auth surface, or wrong middleware assumptions, the user still experiences
failure even if the request technically "worked."

Current truth:

- the archive and plans repeatedly acknowledge env, secret, and revision drift
  as first-class problems
- the current multi-node story still leaves too much semantic alignment inside
  the operator’s head

Failure signature:

- the peer does host the service
- the peer is reachable
- the remote answer still diverges because runtime meaning is not aligned

Proof threshold:

- auditable distribution of env and secret state
- visible node-by-node revision or deployment state
- documented restart versus no-restart semantics for changes
- inspectable drift detection

What counts as evidence:

- node-state reporting
- an explicit source of truth for convergence-critical surfaces
- documented behavior when nodes diverge

Unlocks:

- semantically trustworthy forwarding
- safer multi-node rollouts
- fewer "same service name, different behavior" traps

This is one of the repo's hidden hardest layers because it answers a question
many other stacks skip:

> even if the peer is reachable, is it actually the same service in a way the
> user would trust?

## Priority 3: route persistence under backend loss

Class:

- integrity gate

Why it matters:

If the recovery route dies at the same moment the local backend dies, there is
no failover story.
There is only an attractive diagram.

Current truth:

- `docker-gen-failover` and related routing ideas exist in the repo's
  architecture story
- planning material already records the deletion-risk problem when local
  backends disappear

Failure signature:

- the local backend stops
- route regeneration runs
- the recovery path vanishes with the local target

Proof threshold:

- fallback routes survive independently of local backend liveness
- generator or replacement logic behaves correctly across stop/die events
- at least one real backend-loss drill proves route continuity

What counts as evidence:

- generated config that remains valid after backend disappearance
- drill logs
- a concrete route that still succeeds while the original backend is down

Unlocks:

- honest use of the word failover
- meaningful backend-loss testing
- a non-decorative peer-forward story

## Priority 4: peer eligibility truth

Class:

- foundational

Why it matters:

The repo cannot keep treating "reachable" as a synonym for "safe to receive
traffic."

Current truth:

- the docs already distinguish local, remote, and stateful concerns more
  sharply than before
- runtime proof is still too weak to claim that peers are chosen from a
  trustworthy eligibility model

Failure signature:

- the system can name a remote peer
- but not prove that peer is semantically current, policy-compatible, and safe
  for the class of request being forwarded

Proof threshold:

- eligibility criteria are explicit
- runtime or generation logic consumes them
- at least one representative service path proves the chosen peer was not just
  alive, but acceptable

What counts as evidence:

- a peer-selection model
- logs or generated state showing why a peer was chosen
- failure drills showing bad peers are excluded

Unlocks:

- more trustworthy wrong-node routing
- less hidden operator judgment during failure

## Priority 5: one end-to-end wrong-node HTTP proof

Class:

- pivotal proof gate

Why it matters:

This is the first place the repo can stop sounding like an argument and start
becoming a demonstrated system.

Current truth:

- the dream is explicit
- the edge stack is already rich and serious
- the architecture is most likely to earn real proof first on stateless HTTP

Failure signature:

- all the right words exist
- but there is still no concrete drill where one healthy wrong node receives
  the request and the same service contract survives through a healthy peer

Proof threshold:

- one representative stateless HTTP hostname
- request begins on a node that does not host the service locally
- the node identifies the service as remote
- a healthy peer is selected
- auth, middleware, and externally visible behavior remain coherent
- the request completes successfully

What counts as evidence:

- a documented drill procedure
- logs from the receiving node and selected peer
- proof that the path remained semantically correct, not merely reachable

Unlocks:

- one genuinely defensible "wrong-node HTTP works" claim
- a real baseline for future backend-loss drills
- less temptation to overstate the whole stack prematurely

## Priority 6: keep HTTP, TCP, and stateful classes separate

Class:

- anti-fake-HA guardrail

Why it matters:

The easiest way to accidentally lie is to prove one HTTP path and then narrate
the whole platform as if TCP and stateful services are morally adjacent.

Current truth:

- the docs now separate these classes much better than before
- the runtime still exposes all three classes through the same broad
  infrastructure story

Failure signature:

- a good HTTP failover result gets widened into "the stack is resilient"
- TCP exposure gets mistaken for TCP continuity
- stateful routing gets mistaken for stateful correctness

Proof threshold:

- every critical page keeps the classes separate
- every roadmap decision names which class it affects
- stateful claims remain gated by stateful proof, not ingress proof

Unlocks:

- honest communications
- better promotion rules
- less false comfort

## Priority 7: explicit promotion rules for stronger control layers

Class:

- strategic

Why it matters:

The repo does not need one more broad orchestrator debate.
It needs a disciplined rule for when a stronger layer has finally earned the
right to exist.

Current truth:

- Compose-first remains the live baseline
- OpenSVC, Nomad, k3s, Kubernetes, and helper-surface futures all remain open
- the docs now describe these options more honestly, but promotion is still not
  settled

Failure signature:

- the helper layer quietly becomes a shadow control plane
- or a heavyweight platform gets adopted before the exact earned pain is named

Proof threshold:

- each stronger layer is mapped to named failure classes it would remove
- the operator cost of that layer is explicit
- one domain at a time can be promoted instead of the whole repo jumping

What counts as evidence:

- domain-specific promotion criteria
- control-plane tax stated plainly
- proof that a promoted layer removes something the current path cannot remove
  honestly

Unlocks:

- less ideology
- better sequencing
- a cleaner boundary between "explored" and "earned"

## Bottom line

This roadmap is not here to make the repo sound ambitious.

It is here to identify the exact missing truths that still force the current
stack to overclaim.

The ordering remains harsh because the user’s real benchmark is harsh:

> traffic can land anywhere, still find the right service without sacred-node
> memory, and never pretend reachability is the same thing as preserved
> meaning

That is what the roadmap should keep serving.
