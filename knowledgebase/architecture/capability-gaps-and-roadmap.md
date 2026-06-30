# Capability Gaps and Roadmap

For the evidence shaping this roadmap, start with:

- [Evidence Ledger](../research/evidence-ledger.md)
- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
- [Stateful HA Evidence](../research/stateful-ha-evidence.md)
- [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This page is not a feature wishlist.

It is the map of which missing truth layers still force the platform to
overclaim relative to the user's actual benchmark.

The user is not asking for a tidy next-steps board.
They are asking why so many infrastructures still collapse into one of two
dissatisfying outcomes:

- stay in ordinary Docker and keep wrong-node fragility plus private topology
  memory
- jump into a heavyweight orchestrator before smaller honest layers were even
  exhausted

So the only roadmap worth keeping here is the roadmap that answers:

> which missing truth layer still forces the platform to lean on hope, memory,
> stale assumptions, or theatrical confidence?

That is what "roadmap" means in this repo.

## What this page is and is not allowed to prove

This page is authoritative about:

- prioritizing missing truth layers in the order that protects honesty
- explaining which unresolved gaps still force stronger claims to remain
  illegal
- connecting current runtime evidence to the next proof-bearing promotion step
- separating stateless HTTP, protected HTTP, TCP, and stateful promotion tracks
- explaining what kind of control surface might actually earn the right to
  exist next

This page is not authoritative about:

- acting like a completion report
- implying that a well-ordered roadmap means the runtime is already coherent
- flattening all recovery classes into one generic queue
- upgrading planning quality into present-tense platform maturity

This page is a sequencing contract, not a morale document.

## Strongest honest current answer

The next work is not "whatever seems generally useful for self-hosting."

The next work is whichever missing truth layer still forces the docs and the
runtime story to rely on:

- operator memory
- stale topology assumptions
- helper language that sounds stronger than the proof it owns
- visible machinery that cannot yet cash out into real bad-day decisions

In the priority implementation today, that means the next serious gaps are:

- placement truth
- convergence truth
- route durability under real failure
- peer eligibility truth
- wrong-node proof for one stateless HTTP path
- protected-route continuity on that same path
- keeping TCP and stateful promotion under separate honesty gates

That list is best read as a list of places where the user's frustration is
still objectively correct.

## The roadmap question every candidate layer must survive

Every helper, middle layer, or orchestrator experiment should have to answer:

> after this lands, what exact humiliating thing will the operator no longer
> have to personally remember, infer, narrate, or reconstruct during wrong-node
> entry, backend loss, or stateful recovery?

If the answer is vague, the option is still vague.
If the answer is "it gives us more flexibility" or "it makes things more
dynamic," the burden probably did not move.
If the answer only sounds good on a diagram, it is still mostly theater.

This repo needs that filter because the ecosystem is full of nearby answers
that improve elegance, terminology, or automation while leaving the same human
burden intact.

## The gaps in the order that actually protects honesty

The shortest honest order still looks like this:

1. keep the docs hostile to fake certainty
2. establish placement truth
3. establish convergence truth
4. make routes survive local backend loss
5. define peer eligibility as something stricter than reachability
6. prove one full wrong-node stateless HTTP path end to end
7. prove protected-route continuity on that same path
8. keep HTTP, TCP, and stateful classes under separate proof rules
9. only then decide what has actually earned promotion into a stronger control
   layer

That order is not motivational sequencing.
It is the dependency chain between the dream and reality.

## Gap 1: Documentation honesty and proof custody

This gap sounds softer than the others, but it is first because the repo is now
serious enough that documentation can overpay confidence long before the runtime
earns it.

This gap exists to stop pages from quietly doing any of these:

- treating DNS plurality as end-to-end failover
- treating route presence as route durability
- treating a planning artifact as current runtime truth
- treating a better explanation as a smaller gap

What closure would actually mean:

- the main pages keep dream, runtime, and missing-middle truth separate
- stronger claims stay illegal until backed by route-level or class-level proof

What it still does **not** mean:

- documentation quality equals runtime progress

## Gap 2: Placement truth

This is the first major runtime gap.

The receiving node needs a shared answer to:

- what service lives where right now?
- which node is currently authoritative for this route?

The repo's intent surfaces repeatedly converge on a lightweight current-state
registry such as `services.yaml`.
The important part is not the filename.
The important part is that current placement truth becomes shared, inspectable,
and usable by the receiving node.

What closure would mean:

- a real runtime-consumed placement truth surface exists
- wrong-node routing no longer depends on remembered service location

What still remains illegal even after partial progress:

- generic wrong-node success
- peer-forward continuity
- stateful HA claims

## Gap 3: Convergence truth

Placement truth alone is not enough.
The platform also needs confidence that the shared truth is current enough and
converges fast enough to matter during change and failure.

This gap exists because stale truth can be as misleading as missing truth.

Questions behind this gap:

- how is placement updated?
- who publishes or reconciles it?
- how stale can it get before forwarding becomes dishonest?
- what happens during node churn, partial failure, or split visibility?

What closure would mean:

- the operator is no longer privately carrying the freshest topology truth
- the runtime has a believable story for currentness rather than static hope

What still remains illegal after closure:

- route durability claims without failure drills
- auth or middleware continuity claims

## Gap 4: Route durability under failure

This is where many respectable stacks collapse.

The route cannot merely exist while everything is healthy.
It has to survive the actual failure that made fallback necessary.

This gap is especially live in this repo because helper logic such as
`docker-gen-failover` has already been documented as weak or defective rather
than solved HA.

Questions behind this gap:

- does the wrong-node rescue path still exist after the preferred backend dies?
- does route material persist or vanish under container stop or backend loss?
- does the receiving node still know what to do after the local favorite path
  disappears?

What closure would mean:

- one route class survives a real fallback-triggering failure without private
  human reconstruction

What still remains illegal after partial closure:

- class-wide failover claims across all HTTP routes
- any TCP or stateful promotion

## Gap 5: Peer eligibility truth

Peer reachability is too weak.

The platform needs a stricter answer to:

- which peer is healthy?
- which peer is actually hosting the right service revision?
- which peer is acceptable for that route's auth, middleware, and policy
  expectations?

This matters because a reachable peer can still be the wrong peer.

What closure would mean:

- the receiving node can choose from eligibility truth rather than mesh
  reachability alone

What still remains illegal after closure:

- protected-route continuity claims without semantic comparison
- stateful substitution claims

## Gap 6: Wrong-node proof for one stateless HTTP path

The roadmap needs one brutally narrow success packet before broader language
becomes believable.

That packet should prove one full stateless HTTP wrong-node path end to end:

- request lands on a healthy node that does not host the service locally
- the receiving node discovers the real backend from shared truth
- forwarding works
- the route survives the relevant failure condition if fallback is being claimed
- the operator does not have to privately complete the answer

This is intentionally narrow.
One narrow honest proof is more valuable here than ten diagram-shaped options.

What still remains illegal after closure:

- protected-route continuity claims
- TCP equivalence claims
- stateful promotion

## Gap 7: Protected-route continuity

After one stateless HTTP proof, the next harder problem is protected HTTP.

This gap exists because a route returning a response is weaker than a route
preserving the same protected service semantics.

Protected-route proof has to compare:

- auth behavior
- middleware order and effect
- headers and trust boundaries
- user-visible service meaning before and after handoff

What closure would mean:

- one protected route survives wrong-node forwarding without semantic drift

What still remains illegal after closure:

- class-wide protected-route confidence
- TCP and stateful carryover

## Gap 8: Separate honesty tracks for TCP and stateful surfaces

This gap exists to stop HTTP optimism from leaking into classes where it does
not belong.

The repo already contains TCP-exposed stateful services such as MongoDB and
Redis.
That makes the separation rule live, not theoretical.

Questions behind this gap:

- is the proof only about transport?
- is write authority preserved?
- is substitution safe?
- is the service merely reachable or honestly resilient?

What closure would mean:

- TCP proof packets and stateful proof packets have their own gates
- HTTP success no longer socially upgrades those classes

What still remains illegal after closure:

- generic "the platform is HA now" language

## What kind of stronger control layer could actually earn promotion

Only after the earlier gaps move does it become honest to ask whether a
stronger middle layer or orchestrator has earned the right to exist.

That could eventually include things like:

- a lightweight runtime-consumed placement layer
- a more explicit convergence mechanism
- stricter service discovery or peer registry logic
- or, later, experiments with Nomad, OpenSVC, k3s, Kubernetes, or another
  controller

But the promotion rule stays the same:

- the new layer must remove a concrete hidden burden
- not merely make the architecture story cleaner

The repo is not anti-controller on principle.
It is anti-unearned opacity.

## What still does not count as progress

The following may all be useful and still not count as roadmap closure by
themselves:

- adding more public nodes without proving wrong-node meaning survives
- adding richer proxy logic without proving route durability under failure
- adding more healthchecks without proving peer eligibility semantics
- adding sync helpers without showing how they reduce hidden substitution trust
- standing up Nomad, OpenSVC, k3s, or Kubernetes experiments without proving
  which hidden burden they actually removed
- making the docs calmer, cleaner, or more enterprise-sounding while the same
  runtime truth gaps remain
- making the roadmap more comprehensive while the operator still has to
  privately finish the decisive sentence on the bad day

That last point matters more than it should in a normal repo.
This project is already sophisticated enough that better planning language can
start performing as if it were one more solved capability.
This page has to keep blocking that move.

## What a real roadmap-promotion packet would have to contain

Before this page supports stronger language like "this next step has clearly
earned priority," it should point to a packet containing:

- the exact hidden burden the priority is meant to remove
- the current worktree evidence that the burden still survives
- the narrower artifact, drill, or truth surface that would close it
- the stronger class of claim that closure would unlock next
- the explicit statement of what still remains illegal after that priority lands

Without that packet, even a good sequence can still collapse back into
aspiration.

## Bottom line

This roadmap is not here to make the repo sound strategic.
It is here to keep the docs from illegally spending confidence before the
corresponding truth layer has moved.

The repo already has many sophisticated options nearby.
What it is still starved for is proof that one of those options actually moved
a truth boundary instead of merely enriching the story around it.
