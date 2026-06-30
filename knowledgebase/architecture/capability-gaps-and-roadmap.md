# Capability Gaps and Roadmap

For the evidence shaping this roadmap, start with:

- [Evidence Ledger](../research/evidence-ledger.md)
- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
- [Stateful HA Evidence](../research/stateful-ha-evidence.md)
- [Orchestrator and Control-Plane Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This page is not a feature list.
It is not a self-hosting "next steps" board.
It is not a comfort document.

It is the sequence of missing truth layers that still force the platform to
lean on:

- private memory
- stale assumptions
- helper language that sounds stronger than it is
- visible machinery that still cannot cash out into bad-day dignity

That is what "roadmap" means in this repo.

## The exact roadmap question

The only roadmap question that matters here is:

> which missing truth layer still forces the operator to personally remember,
> infer, narrate, or reconstruct the answer during wrong-node entry,
> backend loss, or stateful recovery?

If a candidate next step cannot answer that question concretely, it is not a
real next step yet.

It may still be a good idea.
It is just not the next burden-removal step.

## What this page is and is not allowed to prove

This page is authoritative about:

- prioritizing missing truth layers in the order that protects honesty
- identifying which unresolved gaps still make stronger sentences illegal
- connecting current runtime evidence to the next proof-bearing promotion step
- keeping stateless HTTP, protected HTTP, TCP, and stateful tracks separate
- naming what sort of control surface might actually earn promotion next

This page is not authoritative about:

- acting like a completion report
- implying that a good sequence means the runtime is already coherent
- flattening all recovery classes into one generic queue
- turning roadmap quality into present-tense platform maturity

This is a sequencing contract, not a morale artifact.

## Strongest honest current answer

The next work is not:

- whatever is fashionable in self-hosting
- whatever sounds clustered
- whatever makes the diagram more adult
- whatever adds one more helper daemon

The next work is whichever missing truth layer still forces the repo to rely on:

- remembered service location
- remembered peer safety
- remembered route persistence
- remembered state authority
- vague helper confidence instead of inspectable proof

In the priority implementation today, that means the next serious gaps are:

- placement truth
- convergence truth
- route durability under real failure
- peer eligibility truth
- one narrow wrong-node stateless HTTP proof
- protected-route continuity on that same route
- keeping TCP and stateful promotion under harsher separate gates

That list is not merely technical.
It is a list of places where the user's frustration is still objectively
correct.

## The filter every candidate layer must survive

Every helper, registry, orchestrator, sync loop, mesh trick, or middle layer
should have to answer this sentence:

> after this lands, what exact humiliating thing will the operator no longer
> have to know privately?

If the answer is vague, the option is still vague.

If the answer is:

- "it gives us more flexibility"
- "it makes things more dynamic"
- "it is more scalable"
- "it is closer to industry standard"

then the hidden burden probably did not move.

That filter is necessary because the ecosystem is full of answers that improve
terminology, automation, or elegance while still leaving the same human SPOF in
place.

## The shortest honest sequence

If the repo wants the shortest burden-faithful order, it still looks like
this:

1. keep the docs hostile to fake certainty
2. establish placement truth
3. establish convergence truth
4. make routes survive the actual backend-loss event
5. define peer eligibility as something stricter than reachability
6. prove one full wrong-node stateless HTTP path
7. prove protected-route continuity on that same path
8. keep HTTP, TCP, and stateful tracks under separate proof rules
9. only then decide what has truly earned promotion into a stronger control
   layer

That order is not motivational sequencing.

It is the dependency chain between the dream and reality.

## Gap 1: documentation honesty and proof custody

This sounds softer than the others, but it is first because the repo is now
serious enough that docs can overpay confidence long before the runtime earns
it.

This gap exists to stop pages from doing any of these:

- treating DNS plurality as end-to-end failover
- treating route presence as route durability
- treating a plan as current runtime truth
- treating explanation quality as smaller implementation burden

What closure would actually mean:

- the major pages keep dream, runtime, planning, and archive pressure separate
- stronger claims stay illegal until backed by route-level or class-level proof

What closure still does **not** mean:

- documentation quality equals runtime progress

## Gap 2: placement truth

This is still the first major runtime gap.

The receiving node needs a shared answer to:

- what service lives where right now?
- which node is authoritative for this route right now?

The repo keeps converging on a lightweight current-state registry such as
`services.yaml`.
The important part is not the filename.
The important part is that placement truth becomes:

- shared
- inspectable
- current enough to matter
- consumable by the receiving node

What closure would mean:

- a live runtime-consumed placement truth surface exists
- wrong-node routing no longer depends on remembered service location

What still remains illegal after partial progress:

- generic wrong-node success
- peer-forward continuity
- stateful HA language

## Gap 3: convergence truth

Placement truth alone is not enough.

The platform also needs a believable story for how that truth stays current
enough during:

- deployment drift
- secret drift
- node churn
- partial failure
- split visibility

This gap exists because stale truth can be as misleading as missing truth.

The questions behind it are:

- who updates the shared truth?
- how is the update propagated?
- how stale can it get before forwarding becomes dishonest?
- how does the runtime expose disagreement?

What closure would mean:

- the operator is no longer privately carrying the freshest topology answer
- the platform has a believable story for currentness rather than static hope

What still remains illegal after closure:

- route durability claims without failure drills
- protected-route continuity claims

## Gap 4: route durability under failure

This is where many respectable-looking stacks collapse.

The route cannot merely exist while everything is healthy.
It has to survive the failure that made fallback necessary.

This gap is especially alive here because `docker-gen-failover` is already part
of the runtime and already part of the repo's skepticism.

The current repo already preserves the harder warning:

- helper existence is weaker than route persistence
- a fallback surface can still disappear during the exact failure it is meant
  to rescue

What closure would mean:

- one named route survives a real fallback-triggering failure without private
  human reconstruction

What still remains illegal after partial closure:

- class-wide HTTP failover claims
- TCP equivalence claims
- stateful promotion

## Gap 5: peer eligibility truth

Peer reachability is too weak.

The platform needs a stricter answer to:

- which peer is healthy?
- which peer actually hosts the right service revision?
- which peer is acceptable for that route's auth, middleware, and policy
  expectations?

This matters because a reachable peer can still be the wrong peer.

What closure would mean:

- the receiving node chooses from eligibility truth rather than reachability
  plus hope

What still remains illegal after closure:

- protected-route continuity claims without semantic comparison
- stateful substitution claims

## Gap 6: one wrong-node stateless HTTP proof

The roadmap needs one brutally narrow success packet before broader language
becomes believable.

That packet should prove one full stateless HTTP wrong-node path end to end:

- the request lands on a healthy node that does not host the service locally
- the receiving node discovers the real backend from shared truth
- forwarding succeeds
- the route survives the relevant failure condition if fallback is being
  claimed
- the operator does not have to privately complete the answer

This is intentionally narrow.

One honest narrow proof is more valuable here than ten architecture-shaped
options.

What still remains illegal after closure:

- protected-route continuity claims
- TCP equivalence claims
- stateful promotion

## Gap 7: protected-route continuity

After one stateless HTTP route is proven, the next pressure is not "declare the
HTTP story solved."

It is:

> does the same thing remain true when the route has policy meaning?

The runtime already has real protected-route surfaces:

- `code-server` with `nginx-auth@file`
- `dozzle` with `nginx-auth@file`
- `portainer` with `nginx-auth@file`
- various metrics/admin routes with auth middleware
- TinyAuth as a live auth surface in the edge fragment

This means the protected-route gap is not hypothetical.

What closure would mean:

- one named protected route behaves like the same protected route after
  wrong-node handoff
- auth, middleware ordering, and trust assumptions survive the handoff

What still remains illegal after closure:

- TCP parity claims
- stateful equivalence claims
- whole-platform "one cloud" language

## Gap 8: keeping TCP and stateful services on harsher tracks

This repo becomes dishonest very quickly if HTTP success bleeds into TCP or
stateful dignity.

The current runtime already exposes TCP surfaces like:

- `mongodb`
- `redis`

It also already depends on stateful components like:

- `headscale-server`
- `nuq-postgres`
- `litellm-postgres`
- `rabbitmq`

Headscale remains the clearest current warning because the active config still
uses SQLite at `/var/lib/headscale/db.sqlite`.

That means:

- an externally reachable control-plane route is real
- distributed state authority is still a separate problem

What closure would mean:

- per-stateful-class authority, promotion, persistence, and rediscovery
  semantics are explicit

What still remains illegal before that:

- "HA" language for stateful surfaces based only on ingress or TCP exposure

## Gap 9: promotion only after burden transfer

Only after the earlier gaps are narrower should the repo decide whether a
stronger control surface has truly earned promotion.

Candidate layers might include:

- a narrow runtime-consumed registry layer
- stronger route-generation or sync logic
- OpenSVC-shaped ingress control
- Nomad-like scheduling and service discovery
- k3s or Kubernetes for broader desired-state control

But no candidate earns promotion because it sounds serious.

It only earns promotion if it can prove:

- one named hidden operator burden actually moved
- the new layer stays inspectable
- the bad-day story is more honest, not just more automatic-looking

## The current bottom line

The strongest current roadmap sentence is:

> the next work is not broader orchestration in the abstract; it is whichever
> missing truth layer most directly stops the operator from being the last
> untracked dependency during wrong-node entry, backend loss, or stateful
> recovery

That is the roadmap this repo is allowed to keep.
