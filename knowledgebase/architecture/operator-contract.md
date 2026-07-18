# Operator Contract and Success Criteria

This page defines the acceptance bar for `bolabaden-infra`.

It exists to stop the repo from quietly answering a smaller, calmer question
than the one the user is actually trying to solve.

The user is not really asking:

- how do we make Docker feel a little more clustered?
- how do we add nicer failover helpers?
- how do we keep Compose while scaling modestly?
- how do we write prettier architecture docs?

Those are neighboring questions.

They are not the contract.

## The contract in one sentence

Several ordinary Docker nodes should start behaving like one
request-preserving personal cloud while keeping Compose as the readable
authoring contract and refusing sacred-node folklore, fake HA language, and
heavy controller tax unless that extra complexity clearly removes a real hidden
operator burden.

That sentence is what every architecture page in this repo is supposed to
answer to.

## The dream underneath the contract

The user's deeper dream is not merely uptime.

It is dignity under wrong-node entry.

That means:

- the first hop should stop feeling lucky
- the wrong node should stop feeling humiliating
- route rescue should stop depending on private memory
- auth and middleware should stop turning to mush during peer-forward handoff
- stateful services should stop being narrated more optimistically than their
  real authority model deserves

If those things do not improve, the repo may still become more elaborate
without getting meaningfully closer to the real target.

## What this page is allowed to do

This page is allowed to:

- define success in the user's actual terms
- preserve the user's impatience with fake relief
- separate HTTP, protected HTTP, TCP, and stateful lanes
- state what a good path and a bad path must feel like
- forbid stronger language until the runtime actually earns it

This page is not allowed to:

- pretend the current runtime already satisfies the contract
- treat clearer prose as proof of operator relief
- flatten all service classes into one maturity bar
- let plural DNS impersonate end-to-end failover

## Strongest honest current answer

The docs now understand the contract more clearly than the runtime currently
proves it.

That is progress.
It is not completion.

The current runtime still leaves the hardest clauses unmet:

- generic wrong-node success is not yet proven
- backend-loss fallback is not yet proven to survive the triggering failure
- protected-route parity after handoff is not yet proven
- stateful authority is still far harsher than edge reachability
- current placement truth is still more architectural pressure than proven
  runtime authority

So the contract remains intentionally uncomfortable.

## The operator relief this contract is trying to buy

This repo is supposed to remove very specific private burdens from the
operator's head.

The contract is moving in the right direction only when fewer moments still
sound like this:

- "I know which node really hosts that service"
- "I know which peer is actually safe right now"
- "I know that helper route will disappear if the backend dies"
- "I know that the forwarded route is not semantically the same one"
- "I know this stateful hostname only looks portable from the outside"

Those private completions are the real human SPOF the repo is trying to kill.

## The anti-benchmark

The platform has failed the contract in any lane where one of these remains
true:

- the safest answer to "what runs where right now?" is still private operator
  memory
- a healthy receiving node cannot determine local versus remote service
  ownership from shared truth
- peer choice is still mostly reachability plus hope
- the rescue route disappears during the failure it exists to cover
- a forwarded protected route does not remain the same protected service
- a stateful service is described as resilient because it is reachable, not
  because write authority and promotion are defined

If one of those survives, the lane may be better understood.
It is not contract-satisfying.

## The proof-bundle standard

Before the docs claim the contract is being met in any lane, they should be
able to point to a proof bundle containing:

- the exact route or service class exercised
- the node that received the first hop
- the truth source the receiving node used
- evidence that the next decision was system-owned rather than privately
  reconstructed
- evidence that the user-visible contract stayed the same after handoff
- the explicit sentence naming which stronger clauses remain unmet

Without that bundle, the docs may be giving orientation rather than relief.

The bundle should also name the old human job it replaced.

That matters because a route can work while still leaving the operator as the
person who knows why it worked.
The contract is not satisfied until the system can expose enough of the reason
for the operator to inspect it without privately reconstructing the whole
scene.

## The v1 completion boundary

The first honest completion boundary is intentionally narrow.

For this repo, a believable v1 would mean:

1. the root Compose-first runtime remains the visible authoring surface
2. one explicit current-placement or peer-eligibility truth source exists
3. one active request path consumes that truth
4. one public node that does not host the service can still preserve a named
   stateless HTTP request
5. the same route has a backend-loss proof packet
6. one protected route has a policy-parity proof packet
7. stateful services are still documented under a separate harsher lane rather
   than borrowed into the HTTP win

That v1 would not mean the whole platform is solved.
It would mean the repo finally proved one complete transfer of a private
operator burden into system-owned truth.

Anything weaker can still be useful progress, but it should not be called the
missing middle layer.

## V1 completion audit table

This table is the contract-level audit.
It exists so the repo cannot quietly call the knowledgebase, a plan, or a
plausible helper the first real completion boundary.

| V1 clause | Evidence that would prove it | Current honest status | Why completion is still illegal |
| --- | --- | --- | --- |
| Compose remains the visible authoring surface | Root `docker-compose.yml`, active `compose/` fragments, and validation output still define the priority runtime. | Materially true as the current architecture baseline. | This proves authorship shape, not multi-node burden transfer. |
| One explicit placement or peer-eligibility truth source exists | A tracked or generated source such as `services.yaml`, runtime registry output, OpenSVC/Nomad inventory, or equivalent is consumed by routing or eligibility logic. | Defined as a requirement, not proven as a consumed root-runtime authority. | A design for placement truth is weaker than a receiving node using it. |
| One active request path consumes that truth | A request packet shows the receiving node read the placement source and used it to choose local serve or peer handoff. | Not yet proven. | Without consumption evidence, the source can still be decorative topology memory. |
| One non-owner public node preserves one stateless HTTP request | A `route_packet` and `placement_decision_packet` show first-hop node, owner node, selected peer, peer eligibility, service identity, and preserved meaning. | Not yet proven. | A local `200`, DNS plurality, or route existence does not prove wrong-node dignity. |
| The same route has backend-loss proof | A backend-loss packet shows the preferred backend was removed or broken and the rescue path remained usable. | Not yet proven. | A fallback file, helper container, or healthcheck does not prove survival under the failure that matters. |
| One protected route has policy-parity proof | Local and wrong-node/peer-forwarded protected-route behavior are compared for auth challenge, middleware, headers, and trust boundary. | Not yet proven. | "The page loaded" is not proof that the protected route stayed the same protected route. |
| Stateful services remain under a harsher lane | Stateful docs and proof packets keep reachability separate from write authority, promotion, fencing, storage, and rediscovery. | Partly true as documentation discipline. | Documentation discipline does not prove any stateful service became HA. |

The table should be read harshly.
Only the first and last clauses have meaningful documentation support today,
and even those are support for boundaries, not proof of the v1 dream.

The first time this table should become more optimistic is after a narrow route
drill leaves behind packets that a second operator can inspect without
reconstructing the topology from memory.

## The v1 non-goals

These are not required for the first honest completion boundary:

- full Kubernetes-grade scheduling
- generic service mobility for every workload
- stateful HA for every database-shaped service
- automatic migration of all Compose services
- one dashboard that explains every failure mode
- a single universal controller for every traffic class

Those may become future work.
They are not the first honest proof target.

The first target is smaller and stricter:

> prove that one healthy wrong node can stop needing the operator for one real
> route, then keep the claim narrow.

## Acceptance criteria by lane

### 1. Any-node entry stops being a coin flip

Acceptance means:

- more than one public node can receive the first hop
- the docs do not confuse that with full failover
- the next decision after the first hop is increasingly based on shared
  current truth rather than private recollection

This is why Cloudflare DDNS and plural A records matter, but only as the first
step.

They help with sacred-ingress reduction.
They do not finish the contract.

### 2. Locality remains first-class

If the requested service is already local to the receiving node, that path
should stay local.

Acceptance means:

- the node can determine locality honestly
- local-first behavior is deliberate, not accidental
- peer-forwarding does not get narrated as though locality has become
  irrelevant

The user is not trying to erase locality.
The user is trying to make locality non-fragile.

### 3. Wrong-node entry becomes survivable

This is the center of gravity for the whole repo.

If the request lands on a healthy node that does not host the service locally,
the receiving node should be able to determine:

- that the service is not local
- where it actually lives now
- which peer is eligible now
- whether the route survives the relevant failure
- whether the user-visible meaning survives the handoff

If any of that still depends on private completion by the operator, the core
contract is still unmet.

### 4. Protected routes remain the same protected routes

This lane is stricter than plain stateless HTTP.

For protected routes, acceptance means peer-forward handoff preserves:

- auth challenge behavior
- forward-auth behavior
- middleware order
- filtering and trust-boundary assumptions
- path/header semantics that define what the route actually is

The live runtime already gives this lane teeth because real surfaces already
use `nginx-auth@file`, and TinyAuth already exists as a live forward-auth
component.

That means the repo has no excuse to flatten protected-route proof into "it
returned content."

### 5. Fallback survives the failure that made fallback necessary

This clause exists because the repo already knows one of its real traps:

- `docker-gen-failover` is directionally relevant
- it is a real part of the runtime
- helper existence is still weaker than route durability

Fallback that evaporates on the bad day is theater.

The user does not want theatrical HA.
The user wants the burden moved.

### 6. The control surface stays inspectable

The repo is Compose-first for a reason.

The user is not refusing Swarm, k3s, Nomad, or Kubernetes out of stubborn
purism.
The user is refusing to prepay controller tax for a control plane that still
cannot prove it removed the right burden.

Acceptance means an operator can still answer:

- what runs where right now?
- who says that is true?
- how did this node choose local serve versus peer handoff?
- what changed when the preferred backend disappeared?
- which extra component is actually paying for itself?

If the answer becomes "the controller knows" without inspectable proof, the
repo has drifted away from its purpose unless that controller clearly bought
back a real hidden-human SPOF.

### 7. Stateful services stop being lied about

This is the harshest lane.

For `mongodb`, `redis`, `nuq-postgres`, `litellm-postgres`, `rabbitmq`,
`headscale-server`, and similar surfaces, acceptance requires answers to
questions like:

- who owns writes?
- who replicates from whom?
- how does promotion work?
- what actually breaks on node loss?
- how do clients rediscover the real authority?
- is the real failure domain still one disk path?

Headscale is the clearest current warning because its config still points to
SQLite at `/var/lib/headscale/db.sqlite`.

That means a healthy HTTP route to Headscale is not a proof of Headscale HA.

Stable hostnames, TCP routing, and "it still answers" do not satisfy this
lane.

## What stronger sentences remain forbidden

Until matching proof bundles exist, the docs should keep forbidding sentences
such as:

- "wrong-node behavior is basically handled now"
- "fallback is mostly solved for HTTP"
- "the platform now behaves like one cloud"
- "the remaining work is mostly automation"
- "the stateful lane is resilient enough"

Those are exactly the sentences that turn progress into self-flattery before
relief actually arrives.

## What real progress looks like

The smallest contract-faithful progress sequence is:

1. externalize placement truth from private memory
2. prove one wrong-node stateless HTTP route
3. prove whether the same route survives backend loss
4. prove protected-route continuity for one named service
5. keep TCP and stateful lanes on separate harsher tracks

## Contract ledger

Use this ledger to keep future edits honest.

| Lane | Minimum proof | Operator sentence that should die | Stronger sentence still forbidden |
| --- | --- | --- | --- |
| Any-node entry | More than one public node receives traffic and first-hop identity is observed | `I know which public node is the real one.` | `Any-node entry preserves the request end to end.` |
| Locality | A receiving node can explain local versus remote service ownership from an explicit truth source | `I know whether this service is local here.` | `All placement truth is system-owned.` |
| Wrong-node HTTP | One named stateless route succeeds through a non-owner first hop with peer identity observed | `I know which peer should answer this wrong-node request.` | `Wrong-node routing is generically solved.` |
| Backend loss | The same route survives loss of the preferred backend through the intended rescue path | `I know whether fallback survives the failure.` | `HTTP failover is broadly solved.` |
| Protected HTTP | One protected route preserves auth, middleware, headers, and visible behavior across handoff | `I know whether this forwarded route is still the same protected service.` | `Protected-route continuity is solved across the stack.` |
| TCP | A raw TCP route is reachable and explicitly labeled as transport-only proof | `I know whether this port is reachable.` | `Stateful authority is resilient.` |
| Stateful authority | One service defines and exercises writer, replica, promotion, and rediscovery behavior | `I know who owns writes after failure.` | `The general stateful lane is anti-SPOF.` |

This ledger is deliberately unforgiving.
It prevents one lane's success from borrowing dignity for another lane.

That sequence matters because it moves the repo from:

- a stronger explanation of the dream

to:

- one more hidden human burden actually leaving the operator's head

## What this contract says about future orchestration choices

This contract does not ban Nomad, k3s, Kubernetes, Swarm-like behavior, or any
other control layer in principle.

It does ban adopting them for the wrong reason.

A future controller only earns its keep here if it can prove it materially
improves one or more of these:

- current placement truth
- peer eligibility truth
- wrong-node request preservation
- backend-loss route durability
- protected-route semantic continuity
- inspectable stateful authority

If it only makes the stack look more adult while leaving those burdens in the
operator's head, then it violated the contract even if it introduced more
industry-standard nouns.

## The strongest honest contract sentence today

Today the strongest honest sentence is:

> the repo now describes the user's real benchmark more faithfully, but the
> runtime still meets that benchmark unevenly because the most important
> wrong-node, backend-loss, protected-route, and stateful authority clauses are
> not yet broadly proven
