# Operator Contract and Success Criteria

This page defines the success contract for `bolabaden-infra`.

It exists because the rest of the knowledgebase can become smarter, more
concrete, and more honest while still missing the one standard the user is
actually forcing on the repo.

That standard is not:

> is there an interesting Compose-first multi-node architecture here?

It is:

> does this stack start behaving like one request-preserving personal cloud
> instead of several ordinary Docker nodes whose correctness still depends on
> operator memory, lucky request placement, and fake closure?

If this page gets soft, the whole site gets soft again.

## What this page is and is not allowed to prove

This page is allowed to:

- define the acceptance bar the rest of the docs answer to
- keep the user's real benchmark sharper than generic HA language
- distinguish burden removal from architecture theater
- state what "solved enough to feel different" has to mean

This page is not allowed to:

- claim that the current runtime already satisfies the contract
- use good acceptance language as substitute proof
- flatten stateless, protected, TCP, and stateful classes into one bar
- quietly shrink the dream into an easier neighboring goal

## Strongest honest current answer

The repo now has a much sharper operator contract than ordinary self-hosting
stacks usually write down, but the implementation still only partially meets
it. The most important success criteria remain intentionally uncomfortable:
wrong-node requests must stop feeling like a gamble, fallback must survive the
failure that made fallback necessary, protected routes must preserve policy
meaning, and stateful services must stop being described by marketing-level
availability vocabulary.

## The dream in one sentence

The operator wants several ordinary Docker nodes to behave like one
request-preserving personal cloud while keeping Compose as the readable
authoring surface and refusing fake HA, hidden sacred nodes, and heavyweight
orchestrator tax unless that tax clearly removes real pain.

That sentence is the center of gravity for the whole repo.

## What the user is actually tired of

The user is not fundamentally asking for:

- more YAML
- more healthchecks
- more load balancers
- more dashboards
- more platform brands

The recurring frustration is deeper:

- Compose is readable until wrong-node routing matters
- static glue keeps pushing burden back into operator memory
- many "HA" answers stop at the first hop
- many tools solve deployment without solving request preservation
- many heavier orchestrators demand worldview surrender before proving they are
  paying down the right pain

That is why the contract here is about lived operational relief, not just
component inventory.

## What "solved" has to feel like

This repo is not solved because:

- a proxy runs
- multiple A records exist
- several nodes are online
- a route template can be generated

It is solved only when the platform feels materially different at request time
and failure time.

The key phrase is "materially different."
The user is not asking for a better narrated version of the same hidden burden.
They are asking for the burden itself to move.

## The actual acceptance criteria

Each criterion below is a real bar. The runtime can meet some and miss others.
Nothing about one class of success upgrades the rest.

### 1. Any-node entry stops feeling like a coin flip

The operator should be able to think:

> it is okay if traffic lands on node A, B, or C first

without secretly meaning:

> I hope the right host happened to be chosen.

Acceptance means:

- more than one public node can receive the first hop
- the docs do not confuse first-hop plurality with full request preservation
- the next decision after first hop is system-owned more than memory-owned

### 2. Local service remains genuinely local

If the requested service already runs on the receiving node, the platform
should serve it there.

Why this matters:

- locality stays fast
- locality stays inspectable
- the stack does not invent fake cluster ritual just to look distributed

Acceptance means:

- the receiving node can determine locality honestly
- local serving is visible rather than accidental
- the docs do not narrate global forwarding as if local-first no longer matters

### 3. Wrong-node entry becomes survivable

If the request lands on a node that does not host the target service locally,
the user does not want to hear:

- the proxy is healthy
- the mesh exists
- a peer hostname is reachable
- the DNS was plural

They want the request to still complete correctly because the receiving node
can determine:

- the service is not local
- which peer currently hosts it
- whether that peer is eligible now
- whether the route needed for fallback survives the relevant failure
- whether the user-visible service contract still means the same thing after
  handoff

This is one of the central acceptance criteria of the whole repo.

### 4. Protected routes keep the same policy meaning after handoff

In this repo, routing correctness is not only transport correctness.

For protected surfaces, the following must survive peer handoff:

- auth challenge behavior
- forward-auth behavior
- middleware ordering
- security filtering
- headers and rewrites that matter to the route's meaning

If a forwarded route still returns "something" but no longer behaves like the
same protected route, that is not success.

### 5. Fallback survives the failure that made fallback necessary

This criterion exists because the repo already knows about a specific trap:

- `docker-gen-failover` is meaningful
- the master plan also records that it can currently delete routes when a
  container stops

That means the system is not allowed to claim success unless the route required
for recovery remains present during actual backend loss.

Fallback that evaporates under real failure is theater.

### 6. The control surface stays legible

The operator is not refusing orchestration because they enjoy manual pain.
They are refusing unreadable control planes that hide causal truth behind
controller magic or stale metadata.

Acceptance means an operator can still answer:

- what runs where right now?
- who says that is true?
- how did this node decide local serve versus peer handoff?
- what changed when the backend disappeared?
- which extra component is actually paying for itself?

If the answer becomes "the controller knows" without inspectable proof, the
repo has moved away from its purpose unless the extra complexity was clearly
earned.

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

A stateful service is not "HA enough" merely because:

- a global hostname exists
- a TCP router exists
- another node could theoretically be pointed at it
- the container could be restarted elsewhere

### 8. The operator can disappear briefly without the platform losing its mind

One of the repo's deepest questions is:

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
- remembered restart caveats
- remembered hidden route exceptions
- remembered stateful folklore

## Service-class-specific success criteria

The acceptance contract is not one size fits all.

### Stateless HTTP

Success means:

- any-node entry is normal
- local-first serving is preserved when honest
- wrong-node requests can still complete through the right peer
- backend-loss fallback can be exercised and remains present

This is the lane most likely to earn the first serious proof.

### Protected HTTP

Success means everything above plus:

- auth parity
- middleware parity
- visible policy parity

Protected-route success is stricter than plain reachability.

### Raw TCP

Success means:

- transport reaches the right backend
- the receiving node does not lie about what class of recovery is happening
- the docs keep transport proof separate from stateful correctness proof

### Stateful systems

Success means:

- ownership is explicit
- replication is explicit
- promotion rules are explicit
- storage truth is explicit
- client reconnect expectations are explicit
- failure drills are specific to that topology

This is the harshest lane and should stay the slowest to overclaim.

## What definitely does not count as success

The following are useful but insufficient:

- `docker compose config --quiet`
- clean merged Compose output
- healthy local containers
- happy-path public routing
- peer reachability
- "a route answered"
- "the dashboard looked good"

Those may support narrower claims.
They do not satisfy the operator contract by themselves.

## The smallest honest path toward meeting the contract

The repo does not need a giant worldview import before it learns anything.
The smallest honest path is narrower:

1. expose live placement truth strong enough to answer "what runs where now?"
2. prove one local stateless HTTP route on the correct node
3. prove the same route through intentional wrong-node entry
4. prove whether fallback survives actual backend loss
5. compare protected-route policy continuity across that same path
6. only then narrate stronger maturity for the HTTP lane
7. keep TCP and stateful classes on their own harsher tracks

That sequence is not the whole dream.
It is the smallest sequence that would make the dream more true rather than
just more articulate.
