# Capability Gaps and Roadmap

For the evidence shaping this roadmap, start with:

- [Evidence Ledger](../research/evidence-ledger.md)
- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
- [Stateful HA Evidence](../research/stateful-ha-evidence.md)
- [Orchestrator and Control-Plane Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
- [Infrastructure Master Plan: What It Actually Proves](../research/infrastructure-master-plan.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This page is not a morale board.
It is not a "future enhancements" list.
It is not here to make the repo sound more adult.

It exists to answer one specific question:

> which missing truth layer still forces the operator to privately finish the
> story when a request lands on the wrong node, the preferred backend
> disappears, a protected route crosses nodes, or a stateful service loses its
> current authority?

That is what "roadmap" means in this repo.

## What this page is and is not allowed to prove

This page is allowed to prove:

- which hidden burdens still remain system-external
- which missing truth layers are upstream of everything else
- why some future-looking work is actually adjacent rather than next
- why HTTP, protected HTTP, TCP, and stateful paths cannot share one fake
  maturity queue

This page is not allowed to prove:

- that a good sequence means the platform is already coherent
- that a helper, agent, or orchestrator is automatically the next answer
- that all anti-SPOF work is one linear ladder
- that later steps deserve promotion before earlier truths exist

## The user's real roadmap question

The user is not asking:

- what is fashionable in homelabs
- what sounds like clustering
- what gets the repo closer to Kubernetes vocabulary
- what adds the most machinery per diagram

The user is asking:

> what exact missing truth still forces me to personally know the answer better
> than the system does?

That question has to stay brutal.
If it gets softened into "what should we build next?" then the roadmap becomes
generic immediately.

## The roadmap is downstream of one exact contract

The roadmap only stays honest if it stays tied to the repo's clearest explicit
contract in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

Every gap below is really a sentence about why that contract is not yet fully
system-owned.

Most roadmap confusion begins when the contract gets replaced by smaller,
blurrier phrases like:

- improve failover
- add resilience
- remove SPOFs
- make it more clustered

Those phrases are downstream.
They are too weak to prioritize honestly on their own.

## Strongest honest current answer

The strongest honest current answer is that the next work is whichever truth
layer most directly removes one of these private operator burdens:

- remembered service location
- remembered peer safety
- remembered route persistence under backend loss
- remembered policy continuity after handoff
- remembered writer authority for stateful surfaces

That means the next serious gaps are still mostly about truth ownership before
they are about product choice:

- documentation honesty and proof custody
- placement truth
- convergence truth
- route durability under real failure
- peer-eligibility truth
- one narrow wrong-node stateless HTTP proof
- protected-route continuity on that same path
- strict separation between HTTP, TCP, and stateful promotion rules
- stateful promotion only after authority semantics are explicit

That list is not just technical.
It is a list of places where the user's frustration is still objectively right.

## The hidden-sentence test

Every roadmap item should be read as a sentence-killer.

If a candidate step lands, what exact private sentence becomes less true?

Examples of the sentences that still survive today:

- `I still personally know what runs where.`
- `I still personally know which healthy peer is actually safe.`
- `I still personally know whether the fallback route survives the failure that made it matter.`
- `I still personally know whether this protected route still means the same thing after handoff.`
- `I still personally know who owns truth for this stateful workload.`

If a step does not kill or materially shrink one of those, it may still be
useful, but it is not yet the next burden-removal step.

## The filter every candidate step must survive

Every helper, registry, sync loop, middle layer, or orchestrator candidate
should have to answer this sentence:

> after this lands, what exact humiliating thing will the operator no longer
> have to know privately?

Answers like these are not good enough:

- more flexibility
- more scalability
- more dynamism
- closer to industry standard
- more cluster-like behavior

Those are aesthetic or ecosystem answers.
This repo needs burden answers.

## The shortest honest sequence

If the repo wants the shortest burden-faithful route, it still looks like
this:

1. keep the docs hostile to fake certainty
2. establish placement truth
3. establish convergence truth
4. make route continuity survive real backend loss
5. define peer eligibility as something stricter than reachability
6. prove one full wrong-node stateless HTTP path
7. prove protected-route continuity on that same path
8. keep HTTP, TCP, and stateful tracks under separate proof rules
9. only then decide what has truly earned promotion into a stronger control
   layer

This is not motivational sequencing.
It is dependency ordering between the dream and reality.

## Why the order starts by being mean to the docs

Starting with documentation honesty can sound like procrastination unless the
real failure mode is understood.

It comes first because:

- the repo is already complex enough to flatter itself accidentally
- broader platform nouns are already easy to borrow
- one false stronger sentence can contaminate every later decision

If the map lies first, good implementation work becomes harder to judge
honestly.

## Why placement truth comes before almost everything else

Without placement truth:

- wrong-node success is mostly folklore
- peer-forward logic can still choose from remembered assumptions
- failure analysis still depends on the operator reconstructing what ran where

The placement layer matters because the user is sick of stacks where a tool
claims to be distributed while the human still privately tracks real service
location.

This is one of the deepest hidden burdens in the repo.

## Why convergence truth is its own gap

Placement truth by itself can still rot into a polite lie.

If the system says a workload lives somewhere, but node drift, stale Compose
state, or unsynchronized fragments have already moved reality elsewhere, then:

- the registry becomes theater
- route generation inherits bad facts
- wrong-node recovery still collapses into human verification

That is why convergence truth is not a duplicate of placement truth.
It is the question of whether the system's claim about current reality can be
trusted at all.

## Why route durability under real failure is a separate gate

Even if placement and convergence truth improve, the repo still has to answer:

> when the preferred backend actually disappears, does the route meaning
> survive, or does the fallback vanish precisely when it becomes relevant?

This matters because the stack already knows about the helper trap where
generated failover config sounds clever until a stop event deletes the route
that was supposed to preserve dignity.

Until a named route survives an intentional backend-loss event, failover prose
should stay suspicious.

## Why peer eligibility is stricter than reachability

A live peer is not automatically a safe peer.

The missing question is:

> what shared truth makes this peer eligible to receive this request right now?

That truth might depend on:

- current service placement
- policy continuity
- middleware capability
- version or drift state
- workload class
- state authority constraints

Without explicit peer eligibility truth, forwarding remains dangerously close
to "pick a node that seems alive."

## Why the first proof should be narrow

The repo does not need ten proofs first.
It needs one narrow, uncompromising one.

That first proof should be:

- stateless
- HTTP
- wrong-node on purpose
- evidence-rich

Because that is the smallest surface where the system can stop sounding
theoretical and start proving one burden actually moved.

The first proof should not quietly expand into TCP or stateful claims.
That is how fake maturity starts.

## Why protected-route continuity comes immediately after

A plain wrong-node success packet is not enough if the protected routes still
cheat.

The user is not only asking:

- did the request arrive?

They are also asking:

- did it remain the same protected route with the same auth and middleware
  meaning after crossing nodes?

That is why protected-route continuity is not optional polish.
It is one of the most meaningful tests of whether peer forwarding preserved the
semantics of the request instead of just producing a response.

## Why TCP and stateful tracks must stay demoted longer

TCP and stateful paths are easier to overclaim than stateless HTTP.

TCP is easier to overclaim because:

- reachability sounds more impressive than it is
- the port answering is easy to confuse with meaningful continuity

Stateful workloads are easier to overclaim because:

- persistence sounds close to resilience
- a replicated-looking topology sounds close to authority transfer
- people start borrowing cluster language before the writer question is owned

So these tracks must stay under harsher proof rules.

## The roadmap lanes

| Gap | What it really asks | What false completion looks like | What would honestly move it |
| --- | --- | --- | --- |
| Documentation honesty and proof custody | do the docs stop letting softer evidence impersonate stronger truth? | the docs sound mature | lane-specific claim routing and burden-faithful rewrites |
| Placement truth | can the system say what runs where right now? | the operator can usually answer | one live tracked placement authority consumed by routing or eligibility logic |
| Convergence truth | can the system trust its own placement claims? | a registry exists on paper | inspectable sync or drift evidence trusted by the forwarding layer |
| Route durability | does fallback survive the actual failure event? | helper config exists | a named route survives preferred-backend loss with preserved semantics |
| Peer eligibility truth | which peers are safe for which requests right now? | the peer responds to health checks | forwarding decisions made from shared truth stricter than reachability |
| Stateless wrong-node HTTP proof | can one deliberately misplaced HTTP request still complete correctly? | the architecture sounds right | force one request onto the wrong node and preserve the evidence |
| Protected-route continuity | does auth and middleware meaning survive handoff? | the route returns a response | compare local and wrong-node protected-route behavior |
| TCP separation | can transport continuity be described without pretending it solved state? | the port answered | explicit transport-only proof and explicit forbidden stronger claims |
| Stateful authority proof | who writes, who promotes, who fences, who is rediscovered? | the datastore is persistent and reachable | workload-specific proof packets with authority, failure, promotion, fencing, and client behavior |

## What is adjacent but not next

Many things may still be useful later, but they are not yet next if they land
before the truths above:

- a richer custom middle layer
- a new scheduler
- Nomad promotion
- OpenSVC promotion
- k3s or Kubernetes promotion
- more dashboards
- more automation for its own sake

These may eventually earn their place.
They do not earn it just by sounding like more options.

The user is frustrated precisely because the ecosystem keeps offering
option-shaped things that still collapse into the same human SPOF.

## Bottom line

This roadmap is not "what features should we add next?"

It is:

> what exact truths still live outside the system, and in what order do they
> have to move inward before stronger anti-SPOF language becomes honest?

Until that order is respected, the repo may gain more machinery without getting
closer to the user's real dream.
