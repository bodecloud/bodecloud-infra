# Operator Contract and Success Criteria

This page defines the acceptance bar for `bolabaden-infra`.

It is not here to celebrate sophistication.
It is here to stop the repo from quietly answering a smaller question than the
user actually cares about.

The smaller neighboring questions are:

- how do we make Docker feel more clustered?
- how do we add better failover helpers?
- how do we keep Compose while scaling a little?
- how do we document the moving parts more clearly?

Those are all adjacent.
They are not the contract.

## The contract in one sentence

Several ordinary Docker nodes should start behaving like one
request-preserving personal cloud while keeping Compose as the readable
authoring contract and refusing sacred-node folklore, fake HA language, and
heavy controller tax unless that tax clearly removes a real burden.

That is the benchmark every architecture page answers to.

## What this page is and is not allowed to prove

This page is allowed to:

- define what success actually means
- preserve the user's real complaint instead of a calmer neighboring one
- state what success has to feel like during wrong-node entry and backend loss
- separate stateless HTTP, protected HTTP, TCP, and stateful lanes
- forbid stronger sentences until the runtime earns them

This page is not allowed to:

- claim the current runtime already satisfies the contract
- treat better prose as proof of relief
- flatten all service classes into one maturity bar
- let first-hop plurality impersonate request preservation

## Strongest honest current answer

The docs are now much clearer about the contract than the runtime is.
That is useful, but it is not relief.

The runtime still meets the contract only unevenly.
The main unmet clauses are still the decisive ones:

- wrong-node requests are not yet generically proven
- backend-loss fallback is not yet proven to survive the failure that triggers
  it
- protected-route parity after handoff is not yet proven
- stateful authority is still much harsher than edge reachability

So the contract remains intentionally uncomfortable.

## What the contract is trying to buy

The contract is not only technical.
It is trying to buy a specific kind of operator relief:

- the first hop no longer feels lucky
- the wrong node no longer feels humiliating
- the rescue path no longer feels like private folklore
- a protected route remains the same protected route after handoff
- stateful services stop being described more optimistically than their writer
  authority deserves

If those do not move, the contract is not met no matter how adult the stack
sounds.

## The private completion test

The contract only becomes real when fewer important moments still sound like
this:

- "I know which node really has that service"
- "I know which peer is actually safe for that handoff"
- "I know that helper loses the route if the wrong thing dies"
- "I know the forwarded route is not semantically the same one"
- "I know that stateful hostname only looks portable from the outside"

Those are the exact private completions the system is supposed to eliminate.

## The anti-benchmark

The platform has failed the contract in any lane where one of these remains
true:

- the safest answer to "what runs where right now?" is still private operator
  memory
- a healthy wrong node cannot determine local versus remote service ownership
- peer choice is mostly reachability plus hope
- the rescue route disappears during the failure it is meant to absorb
- a protected forwarded route is not the same protected route from the user's
  perspective
- a stateful service is described as resilient because it is reachable, not
  because authority and promotion are defined

If one of those survives, the lane may be better understood.
It is not contract-satisfying.

## Proof bundle standard

Before the docs claim the contract is being met in any lane, they should point
to a bundle containing:

- the exact route or service class exercised
- the node that received the first hop
- the truth source the receiving node used
- evidence that the next decision was system-owned rather than privately
  reconstructed
- evidence that the user-visible contract stayed the same after handoff
- the explicit sentence naming what stronger clauses are still unmet

Without that bundle, the docs may still be describing orientation rather than
relief.

## Acceptance criteria by lane

### 1. Any-node entry stops feeling like a coin flip

Acceptance means:

- more than one public node can receive the first hop
- the docs do not confuse that with end-to-end failover
- the next decision after the first hop is moving from private memory toward
  shared system-owned truth

What does not count:

- multiple DNS records by themselves
- one successful route from the preferred node
- "ingress is redundant now"

### 2. Honest locality remains possible

If the target service is local to the receiving node, the platform should keep
that path local.

Acceptance means:

- the node can determine locality honestly
- local-first behavior is visible, not accidental
- the docs do not narrate peer-forwarding as if locality stopped mattering

### 3. Wrong-node entry becomes survivable

This is the center of gravity for the whole repo.

If the request lands on a healthy node that does not host the service locally,
the receiving node should be able to determine:

- that the service is not local
- where it actually lives now
- which peer is eligible now
- whether the route survives the relevant failure
- whether the user-visible meaning survives the handoff

If those still require private completion by the operator, the central
contract is still unmet.

### 4. Protected routes preserve policy meaning after handoff

For protected HTTP surfaces, success is stricter than "something answered."

The following should survive peer handoff:

- auth challenge behavior
- forward-auth behavior
- middleware ordering
- security filtering
- path and header behavior that define what the route actually is

If a forwarded route returns content but no longer behaves like the same
protected service, that is failure.

### 5. Fallback survives the failure that made fallback necessary

This clause exists because the repo already knows the `docker-gen-failover`
trap:

- it is directionally relevant
- it can still delete routes when the backend disappears

Fallback that evaporates on the bad day is theater.

### 6. The control surface stays legible

The user is not refusing heavier orchestration out of nostalgia.
They are refusing to prepay abstraction tax for a control plane that still
cannot prove it removed the right burden.

Acceptance means an operator can still answer:

- what runs where right now?
- who says that is true?
- how did this node choose local serve versus peer handoff?
- what changed when the preferred backend disappeared?
- which extra component is actually paying for itself?

If the answer becomes "the controller knows" without inspectable proof, the
repo has drifted away from its purpose unless the extra complexity clearly
removed a real burden.

### 7. Stateful services stop being lied about

For `redis`, `mongodb`, `nuq-postgres`, `litellm-postgres`, `rabbitmq`,
`headscale`, and similar surfaces, the questions are harsher:

- who owns writes?
- who replicates from whom?
- how does promotion work?
- what breaks on node loss?
- how do clients rediscover the real authority?
- is the real failure domain still one local disk path?

Stable hostnames, TCP forwarding, and "it still answers" do not satisfy this
lane.

## What stronger sentences are still forbidden

Until matching proof bundles exist, the docs should keep forbidding sentences
such as:

- "wrong-node behavior is basically handled now"
- "fallback is mostly solved for HTTP"
- "the platform now behaves like one cloud"
- "the only remaining issue is broader automation"
- "the stateful lane is resilient enough"

Those sentences are exactly how a repo starts flattering itself before relief
actually arrives.

## What would count as real progress

The smallest contract-faithful progress sequence is:

1. externalize placement truth from private memory
2. prove one wrong-node stateless HTTP route
3. prove whether the same route survives backend loss
4. prove protected-route continuity for one named service
5. keep TCP and stateful lanes on separate harsher tracks

That is the difference between:

- a better architecture story

and:

- one more burden actually leaving the operator's head
