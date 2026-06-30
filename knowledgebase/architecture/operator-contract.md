# Operator Contract and Success Criteria

This page defines the acceptance bar for `bolabaden-infra`.

It exists because the repo now has enough architecture language, enough
research, and enough partially-real machinery that it would be very easy to
sound thoughtful while still answering a smaller question than the user is
actually forcing on the project.

That smaller question is usually one of these:

- how do we make Docker feel more clustered?
- how do we get better failover?
- how do we keep Compose while scaling a little?
- how do we document the moving parts more clearly?

Those are all adjacent.
They are still smaller than the real benchmark.

The real benchmark is:

> can several ordinary Docker nodes start behaving like one request-preserving
> personal cloud, without correctness depending on sacred-node memory, lucky
> first-hop placement, or fake closure dressed up as HA?

If this page gets soft, the whole knowledgebase gets soft with it.

## What this page is and is not allowed to prove

This page is allowed to:

- define the real acceptance bar the rest of the docs answer to
- preserve the user's actual complaint instead of a calmer neighboring one
- state what success has to feel like at request time and failure time
- distinguish burden removal from better-narrated burden
- keep stateless HTTP, protected HTTP, TCP, and stateful lanes under separate
  honesty rules

This page is not allowed to:

- claim that the current runtime already satisfies the contract
- use sharper wording as substitute proof
- flatten all service classes into one maturity bar
- smuggle "mostly solved" language in through architecture tone
- let first-hop plurality impersonate request preservation

This page is the benchmark.
It is not a completion certificate.

## Strongest honest current answer

The repo now states the operator contract much more clearly than before, but
the current runtime still meets it only unevenly. The acceptance bar remains
deliberately uncomfortable: wrong-node requests must stop feeling like a
gamble, fallback must survive the failure that made fallback necessary,
protected routes must preserve the same policy meaning after handoff, and
stateful systems must stop being described with availability language they
have not earned.

That discomfort is not a tone choice.
It is a defensive measure against a specific documentation failure:

- the stack becomes easier to explain
- the architecture becomes easier to admire
- the operator is still the one privately completing the request story

If the docs stop sounding harsh before that burden moves, then the docs are
helping the lie, not correcting it.

## The dream in one exact sentence

The operator wants several ordinary Docker nodes to behave like one
request-preserving personal cloud while keeping Compose as the readable
authoring contract and refusing fake HA, sacred remembered entrypoints, and
heavyweight orchestrator tax unless that tax clearly removes a real burden.

That sentence is the center of gravity for the whole repo.

## What the contract is emotionally trying to buy

This contract is not only describing technical success.
It is describing relief.

The operator wants to stop feeling that the system only looks coherent because
they are silently carrying the missing truth inside their own head.

That means the contract is trying to buy a very specific kind of operational
relief:

- the first hop no longer feels lucky
- the wrong node no longer feels humiliating
- the rescue path no longer feels like private folklore
- the stateful lane no longer gets described with optimism it did not earn

That is why this page has to stay harsher than an ordinary success-criteria
document.

Ordinary success criteria often ask:

- does it work?
- is it redundant?
- is it automated?
- is it observable?

This repo also has to ask:

- who still had to secretly know the answer first?
- what still only looks solved because the operator can mentally repair the
  topology story on demand?

If the answer is still "the operator," then some part of the contract is still
unmet no matter how respectable the stack sounds.

## What the user is actually tired of

The user is not fundamentally asking for:

- more YAML
- more healthchecks
- more reverse proxies
- more dashboards
- more "cloud-native" product names

The recurring frustration is deeper:

- Compose feels straightforward until wrong-node behavior matters
- many HA answers stop at the first hop and call it closure
- many helper tools solve configuration churn without solving request
  preservation
- many orchestrator answers demand worldview surrender before they prove they
  are paying down the right pain
- many docs present several respectable options even when the real burden
  stayed in the operator's head

That is why this contract is about lived operational relief, not component
inventory.

## What success has to feel like

The project is not solved because:

- a proxy is running
- multiple public nodes exist
- Traefik labels are present
- a mesh is reachable
- routes can be generated on the happy path

It is solved only when the platform feels materially different during the two
moments that expose the truth:

1. when the request lands on the wrong healthy node
2. when the preferred backend disappears

"Materially different" means the operator should stop feeling like the system
quietly needs them to remember the topology in order to be coherent.

That last clause matters.
The user is not asking for a cleaner story about the same hidden burden.
They are asking for the burden itself to move.

This is the page where "move" must stay literal.

The docs should only become stronger when some truth actually moves from:

- private operator reconstruction

to:

- system-owned, inspectable, shared decision surfaces

Anything weaker than that may still be progress.
It is not yet relief.

## What still does not count as operator relief

This page should say the quiet failure mode out loud.

The following still do not count as the operator contract moving in the user's
favor:

- the docs sound more rigorous
- the stack is easier to describe at a high level
- several candidate solutions are compared respectfully
- a helper reduces repetition but not bad-day ambiguity
- a stronger controller is proposed without proving which burden leaves the
  operator

Those things may improve orientation.
They do not yet change the lived moment where the operator still has to answer
the real question from private memory.

This distinction has to stay severe.

The repo is not starved for sophisticated-looking machinery.
It is starved for proof that the sophisticated-looking machinery actually moved
the humiliating part of the job out of one human head.

## What a real operator-relief proof bundle would have to contain

Before the docs claim the contract is being met in any lane, they should point
to a bundle that contains:

- the exact route or service class exercised
- the node where the request landed first
- the truth source the receiving node used
- the evidence that the next decision was system-owned rather than
  operator-reconstructed
- the evidence that the user-visible contract stayed the same after handoff
- the explicit sentence naming which stronger contract clauses remain unmet

Without that bundle, the docs may still be describing aspiration or
orientation rather than relief.

## The negative benchmark

Before listing acceptance criteria, keep the anti-benchmark visible.

The platform has failed this contract if any of these remain true in the lane
being discussed:

- the safest answer to "what runs where right now?" is still private operator
  memory
- a healthy wrong node cannot tell whether the requested service is local
- peer selection is mostly reachability plus hope
- the fallback route disappears during the failure it was meant to absorb
- a protected forwarded route no longer behaves like the same protected route
- a stateful service is being called resilient because it is reachable, not
  because authority and promotion are defined

If one of those survives, the problem may be narrowed.
It is not closed.

## The actual acceptance criteria

Each criterion below is real.
Success in one lane does not silently upgrade the rest.

### 1. Any-node entry stops feeling like a coin flip

The operator should be able to think:

> it is okay if traffic lands on node A, B, or C first

without secretly meaning:

> I hope the right machine happened to be chosen.

Acceptance means:

- more than one public node can receive the first hop
- the docs do not confuse first-hop plurality with full request preservation
- the next decision after first hop is increasingly system-owned rather than
  memory-owned

What does **not** count:

- multiple A records by themselves
- a single successful route from a preferred node
- vague language like "ingress is redundant now"

### 2. Local service remains honestly local

If the target service already runs on the receiving node, the platform should
serve it there.

Why this matters:

- locality stays fast
- locality stays inspectable
- the platform does not invent unnecessary global indirection for aesthetics

Acceptance means:

- the receiving node can determine locality honestly
- local-first behavior is visible rather than accidental
- the docs do not narrate peer forwarding as if locality stopped mattering

### 3. Wrong-node entry becomes survivable

If the request lands on a healthy node that does not host the service locally,
the user does not want to hear:

- the proxy is healthy
- the mesh exists
- a peer hostname is reachable
- the DNS is plural

They want the request to still complete correctly because the receiving node
can determine:

- the service is not local
- where it actually lives now
- which peer is eligible now
- whether the recovery path survives the relevant failure
- whether the user-visible service meaning survives the handoff

This is the central acceptance criterion of the whole repo.

It is central because it is where most fake closure finally breaks.

Many infra stories can sound nearly complete right up until this moment.
Then the real test appears:

- can the receiving node decide correctly without a human secretly supplying
  the topology?

If not, the platform may be rich, but it is still not the thing the user
actually came here trying to build.

### 4. Protected routes keep the same policy meaning after handoff

In this repo, routing correctness is not only transport correctness.

For protected surfaces, the following must survive peer handoff:

- auth challenge behavior
- forward-auth behavior
- middleware ordering
- security filtering
- headers, rewrites, and path behavior that define what the route actually is

If a forwarded route still returns "something" but no longer behaves like the
same protected service, that is not success.

### 5. Fallback survives the failure that made fallback necessary

This criterion exists because the repo already knows about a specific trap:

- `docker-gen-failover` is meaningful and directionally relevant
- the master plan also records that it can currently delete routes when a
  container stops

That means the system is not allowed to claim success unless the route needed
for recovery remains present during actual backend loss.

Fallback that evaporates under the bad day is theater.

### 6. The control surface stays legible

The operator is not refusing heavier orchestration because they enjoy manual
pain.
They are refusing unreadable control planes that hide causal truth before they
have clearly earned that hiding.

Acceptance means an operator can still answer:

- what runs where right now?
- who says that is true?
- how did this node decide local serve versus peer handoff?
- what changed when the backend disappeared?
- which extra component is actually paying for itself?

If the answer becomes "the controller knows" without inspectable proof, the
repo has drifted away from its purpose unless the added complexity clearly
removed a real burden.

This clause exists because unreadable seriousness is one of the user's main
anti-goals.

The user is not refusing heavier orchestration out of nostalgia.
They are refusing to prepay abstraction tax for a control plane that still
cannot prove it removed the specific bad-day dependence on remembered sacred
nodes, remembered route exceptions, and remembered recovery folklore.

### 7. Stateful services stop being lied about

For Redis, MongoDB, Postgres, RabbitMQ, Headscale, and similar systems, the
docs must stay much harsher.

The operator wants blunt answers:

- who owns writes?
- who replicates from whom?
- how does promotion work?
- what breaks on node loss?
- what reconnect behavior should clients expect?
- is the real failure domain still one local disk path?

A stateful surface is not "HA enough" merely because:

- a global hostname exists
- a TCP router exists
- another node could theoretically be pointed at it
- the container could be restarted elsewhere

This criterion also protects the whole site from a common collapse:

- "the stateless story improved, so we can speak more confidently about the
  whole platform"

No.

The stateful lane is where the docs have to become less forgiving, not more.
If this lane stays unresolved, the site is still obligated to keep saying so
bluntly.

### 8. The operator can disappear briefly without the platform forgetting itself

One of the deepest repo questions is:

> if I disappear for a week, does the system still know enough about itself to
> preserve a request, or did I only hide my own memory behind cleaner YAML?

That is a legitimate acceptance criterion here.

It means the important truth should increasingly live in:

- tracked current-state surfaces
- inspectable routing inputs
- visible health and eligibility logic
- explicit service-class contracts

and less in:

- remembered sacred nodes
- remembered route exceptions
- remembered restart caveats
- remembered stateful folklore

## What absolutely does not count as success

This repo needs an explicit "does not count" section because too many
neighboring infrastructures keep getting congratulated for smaller wins.

The following are useful, but they do not satisfy this contract by themselves:

- "the dashboard is green"
- "Cloudflare points at multiple boxes"
- "Traefik is healthy"
- "Headscale peers can see each other"
- "a generated route exists while the backend is healthy"
- "the service worked from the node that already hosted it"
- "the orchestrator has a clean story about desired state"

Those may be prerequisites.
They are not closure.

## Service-class-specific success criteria

The operator contract is not one size fits all.

### Stateless HTTP

Success means:

- any-node entry is normal
- local-first serving is preserved when honest
- wrong-node requests can still complete through the right peer
- backend-loss fallback can be exercised and remains present

This is the lane most likely to earn the first serious proof.

### Protected HTTP

Success means everything stateless HTTP needs, plus:

- forward-auth continuity
- middleware parity
- preserved security and rewrite semantics
- user-visible behavior that still feels like the same service

This lane is harder than plain HTTP and should be narrated more harshly.

### TCP

Success means transport claims are separated from stronger service claims.

The docs must distinguish:

- socket reachability
- router correctness
- client reconnect semantics
- ownership semantics behind the routed endpoint

TCP success is not automatically application success.

### Stateful services

Success means:

- authority is explicit
- replication semantics are explicit
- promotion behavior is explicit
- failure and recovery sequence is explicit
- data-loss and split-brain risks are named, not implied away

This is the harshest lane and should remain the slowest to earn stronger
language.

## What "good enough to feel different" really means

The operator does not need every lane solved before the platform feels more
real.
But the platform does need enough truth movement that the experience changes.

The smallest meaningful psychological shift would be:

1. current placement truth becomes inspectable
2. one stateless HTTP wrong-node path is genuinely proven
3. the same route is exercised through backend loss
4. one protected route is compared across local and forwarded behavior
5. the docs remain harsh about TCP and stateful surfaces

At that point the user could reasonably say:

> this is still incomplete, but it is no longer just a better narrated pile of
> ordinary Docker boxes

That is a real threshold.

## The brutal final question

Every future helper, registry, sync layer, agent, or orchestrator promotion
should survive this exact question:

> after this change, what hidden burden moved out of the operator's head and
> into the system itself?

If the answer is vague, the work may still be interesting.
It has not satisfied the operator contract.
