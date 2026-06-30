# Proof Matrix and Drill Catalog

This page exists for the question the repo is most likely to lie about by
accident:

> what exact thing are we claiming to have solved, what exact drill would prove
> it, what exact proof class have we actually reached, and what important part
> of the dream would still remain unproven even after that drill passes?

That is the real center of `bolabaden-infra`.

The user is not asking for a generic infrastructure verification checklist.
They are asking whether several ordinary Docker nodes can stop behaving like a
collection of separate local stacks with nice branding and start behaving like
one request-preserving, operator-legible personal cloud without immediately
forcing surrender to Swarm, Kubernetes, or some other worldview-heavy control
plane.

This page exists because that kind of project can become dishonest very
easily.

It becomes dishonest when:

- config presence gets narrated as live behavior
- happy-path success gets narrated as failure absorption
- route reachability gets narrated as request preservation
- peer reachability gets narrated as peer eligibility
- TCP transport gets narrated as stateful correctness
- a surviving `200` gets narrated as preserved semantics

This page is the repo's claim firewall against those upgrades.

This page therefore has to behave like a proof router, not just a proof list.

For each claim it should force four questions:

1. what exact dream fragment is being claimed
2. what exact proof class is required before that fragment may be narrated more
   strongly
3. what weaker evidence classes must be explicitly refused
4. what part of the dream would still remain unproven even after the drill
   passes

It is also where the repo tries to stop doing the documentation equivalent of
answering the wrong question well.

The user is not merely asking for more tests.
The user is asking for proof that survives contact with the real hidden dream:

- fewer fake options
- less sacred-node folklore
- less operator-only topology memory
- more runtime truth that can explain itself under stress

It also has to stay aligned with the repo's real authority order.

For this exact question:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  is the clearest intent surface for the multi-node Docker, anti-SPOF,
  local-first then peer-forward dream
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
  is mainly a repo-operability and validation surface
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)
  is mainly an authoring-discipline surface

That distinction matters because good proof language here cannot pretend every
instruction file carries the same kind of truth.

If a page starts treating authoring preference as architecture proof, the docs
are already drifting again.

## What this page is and is not allowed to prove

This page is authoritative about:

- what proof classes exist in this repo
- what each class can honestly support
- what stronger claim ceilings must remain closed after a drill passes

This page is not authoritative about:

- whether a specific service or route has actually passed a given drill
- whether one passed drill upgrades the whole stack
- whether the broader dream is satisfied just because a proof class exists

This page is the proof-language contract.
It is not the proof result itself.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "this behavior is only intent so far" | proof classification | intent surfaces plus this matrix | that the intent is weak or optional |
| "this route is wrong-node proven" | drill-result classification | route-specific evidence plus this matrix | that the platform now generically preserves wrong-node traffic |
| "this evidence class is too weak for that claim" | proof-discipline judgment | this matrix, runbook, supporting drills | that the weaker evidence is useless |
| "the next ceiling is still unclosed" | proof sequencing | this matrix plus the concrete route or topology evidence | that earlier proof was fake simply because it was partial |

If a sentence starts using a proof class as if it were a stack-wide verdict, it
has already exceeded this page's authority.

## What a good proof page has to do here

A normal checklist asks whether a feature works.

This repo needs a harsher standard:

> what exact part of the user’s dream did this evidence make more true, and
> what exact part is still unproven even if the test passed?

That is why the proof classes below are intentionally unfriendly to casual
overclaiming.

They are also intentionally repetitive.

The repetition is a guard against the easiest failure mode in this project:
green-looking evidence quietly getting promoted into a stronger architecture
claim because everyone is tired and the stack already sounds sophisticated.

Another way to say it:

the matrix is not here to reward the stack for being elaborate.
It is here to keep pressure on whether the elaborate parts actually removed a
hidden burden the user hates.

## What this page is trying to stop

The biggest structural risk in this repo is not that nothing works.
It is that many parts work *just enough* to tempt the docs into saying more
than the system has earned.

The user is already exhausted by that exact pattern in the wider tooling
ecosystem:

- DNS redundancy gets called HA
- a healthy reverse proxy gets called failover
- a route template gets called dynamic orchestration
- a reachable database gets called resilient
- a bigger platform gets called the answer before the actual missing truth is
  even named

This page exists so `bolabaden-infra` does not recreate that same pattern in a
better accent.

That phrase matters.

If the docs here merely become a more articulate version of the same
ecosystem-wide bluff, then the repo has failed at the exact thing it is trying
to repair.

## The user-visible dream this matrix is actually measuring

The matrix only makes sense if the repo's real dream stays visible.

The dream is not:

- "more nodes"
- "more YAML"
- "more reverse proxies"
- "more healthchecks"
- "more clustering vocabulary"

The dream is:

> any surviving public node can receive the request, determine whether the
> target is local, preserve the request if it is not, survive the failure that
> made fallback necessary, keep policy and auth coherent on the fallback path,
> and do all of that without the operator's head secretly remaining the real
> control plane

That is the standard every claim below is being judged against.

The hidden anti-benchmark is also important:

the docs are failing if an operator can read a passed drill and still not know
whether the proof covered:

- wrong-node entry
- backend-loss fallback
- semantic continuity
- or only happy-path reachability

The page is also failing if a passed drill still lets a reader feel reassured
while remaining unable to answer the harsher question:

> did this prove a real option, or only prove a nicer-looking version of the
> same old operator-dependent workaround?

That phrasing matters because the repo is not only trying to prove that a path
can work once.
It is trying to prove that the system itself now owns more of the explanation
for why the path worked.

This is the repo's version of asking whether a reconstruction actually captured
the important inner structure, or only learned how to sound convincing from
the outside.

## The proof classes

The repo needs harsh proof language because soft proof language is how fake HA
 gets normalized.

Use these proof classes consistently.

There is also an ordering rule:

- never narrate a claim using a stronger proof class than the strongest drill
  actually earned
- always state the next unclosed ceiling

### `Intent only`

Meaning:

- the behavior is clearly wanted
- the docs or architecture surfaces describe it directly
- there is no meaningful live implementation proof yet

Allowed claim:

- this is a real target

Forbidden upgrade:

- the runtime already behaves this way

### `Config present`

Meaning:

- tracked config contains ingredients for the behavior
- the authored system is leaning in this direction

Allowed claim:

- the implementation has been authored toward this outcome

Forbidden upgrade:

- the behavior now exists under live conditions

### `Happy-path runtime`

Meaning:

- the path or service works under nominal conditions
- no meaningful stress or failure condition has been exercised

Allowed claim:

- the tested path works in normal conditions

Forbidden upgrade:

- the path is resilient

### `Wrong-node proven`

Meaning:

- a request intentionally landed on a node that did not host the target
  service locally
- the request still completed correctly for that specific route

Allowed claim:

- this exact wrong-node path is real

Forbidden upgrade:

- the platform now generically preserves wrong-node requests

Additional warning:

- if the route only succeeded because the operator privately knew which node
  was the real host, the drill proved recovery by human reconstruction, not by
  system-owned truth

This distinction matters because the user is specifically trying to escape
“one demo equals one platform” thinking.

It also matters because this repo is already rich enough in edge machinery to
generate many flattering demos that still leave the governing truth gap
untouched.

### `Backend-loss proven`

Meaning:

- fallback was required because the preferred local backend actually
  disappeared or became unavailable
- the route needed for recovery survived long enough to matter

Allowed claim:

- this exact route survives this exact backend-loss class

Forbidden upgrade:

- peer-fallback is broadly solved everywhere
- the failover mechanism is trustworthy in general

### `Semantic continuity proven`

Meaning:

- the peer-handoff path preserved auth, middleware, headers, and externally
  visible policy rather than merely returning a response

Allowed claim:

- this path stays meaningfully the same service under fallback

Forbidden upgrade:

- stateful service classes are now covered too
- every protected route now shares the same semantic contract

This is one of the most important proof classes in the whole repo because it
guards against the easiest lie:

- the request succeeded
- therefore the same service survived

Those are not the same claim.

### `State-proven`

Meaning:

- ownership, replication, election or promotion, reconnect behavior, and
  durable correctness have all been shown for the relevant stateful class

Allowed claim:

- this stateful class now has an honest resilience story

Forbidden upgrade:

- the entire platform is now anti-SPOF in a general sense

## The claims that actually define whether the dream is getting closer

These are not generic cloud capabilities.
They are the exact claims that answer whether the repo is becoming less
dependent on hidden human memory, lucky request placement, and architecture
theater.

Another way to read the table below is:

- which claims actually reduce the operator's reconstruction tax
- which claims only decorate the same tax with better tooling

That means the table should not be read as feature progress first.
It should be read as burden-relocation progress first.

| Claim | What the claim really means | Current proof class | Strongest current evidence | Exact next drill | What that drill still would not prove |
| --- | --- | --- | --- | --- | --- |
| Any healthy public node can be the first hop | More than one node can intentionally receive ingress without one sacred public reverse-proxy box | Config present | Cloudflare plurality and anti-sacred-node intent are explicit in [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md), [README.md](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md), and the root runtime surface | Force one representative HTTP route through at least two distinct public node-entry targets and correlate edge logs on both nodes | That the receiving node can preserve the request correctly when the service is not local |
| Local-first service is real | A receiving node can determine that a target is local and serve it there without pretending a scheduler made the decision | Intent only | The request contract is explicit in `copilot-instructions.md` and the architecture pages | Pick one stateless HTTP route and prove local service from the hosting node with logs showing the backend remained local | That wrong-node forwarding works when locality is absent |
| Wrong-node HTTP requests succeed | A node that does not host the service locally can still preserve the request by handing it to the correct peer | Intent only | The dream is explicit; the root Traefik-centered stack is mature enough to make this worth proving | Intentionally land a request for one stateless HTTP route on the wrong node and prove receiving-node identity, backend-node identity, and user-visible success | That the route also survives when fallback is required because the preferred local backend disappeared |
| Fallback route survives backend loss | The route required for recovery remains present after the preferred local backend disappears | Partial | Planning and knowledgebase pages explicitly record route-persistence risk around generated failover behavior | Start from a known-good route, remove the preferred local backend, and prove whether the fallback route remained present long enough to preserve the request | That auth and middleware semantics stayed identical after handoff |
| Middleware and auth continuity survive peer handoff | A peer-forwarded request still behaves like the same protected route rather than a semantically different shortcut | Partial | TinyAuth, CrowdSec, and related edge-policy surfaces are materially live in the root stack | Compare one protected route locally versus through intentional peer fallback and verify auth challenge, middleware, headers, and visible policy parity | That all protected routes now share the same continuity guarantees |
| Placement truth is live and shared | The system has an auditable answer to "what runs where right now?" that is stronger than operator folklore | Intent only | `services.yaml` remains central in architecture intent, but the priority runtime does not currently prove a live tracked root [`services.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/services.yaml) consumed by routing | Introduce or expose one real placement-truth surface and prove that routing or eligibility logic consumes it | That convergence, drift detection, and restart semantics are trustworthy across revisions and secret changes |
| Node convergence is strong enough to trust peer traffic | A forwarded request lands on a semantically compatible peer, not merely a reachable one | Planned | Env and secret requirements are explicit in [AGENTS.md](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md), and the runbook already names convergence as a current gap | Expose node-by-node revision, env, secret, and policy readiness for one service group and prove the chosen peer is actually eligible | That automatic relocation or broader scheduling is now safe |
| TCP routing is real | Raw TCP services are actually exposed through the live ingress or control surface | Happy-path runtime | TCP routers for Redis, MongoDB, and related services exist in the root stack | Validate one representative TCP service through its intended routed path | That the TCP path is resilient or semantically trustworthy under node failure |
| TCP failover is meaningful | A TCP client survives the intended failure class without ambiguous or unsafe topology behavior | Intent only | The repo clearly wants L4 handling and keeps TCP separate from HTTP in its intent surfaces | Define and run one L4 failure drill for a non-critical TCP workload first | That state-bearing TCP services now have honest HA |
| Redis resilience is real | Redis has explicit topology truth, failover semantics, and reconnect behavior | Planned | Redis is live and planning pages already treat Sentinel-like or equivalent topology truth as the minimum honest bar | Prove a Redis primary and replica or equivalent topology with failover and client reconnect behavior | That MongoDB, RabbitMQ, or other stateful classes are equally mature |
| MongoDB resilience is real | Replica-set election and client discovery preserve correctness after node loss | Planned | MongoDB is live and the repo already names replica-set semantics as the minimum honest threshold | Prove working replica-set behavior, election, and client discovery under failure | That the broader platform is generically anti-SPOF |
| RabbitMQ resilience is real | Queue durability and failover expectations are explicit and verified | Planned | RabbitMQ exists in the runtime, but the current docs do not show strong failure proof | Either explicitly bound RabbitMQ as single-node or prove clustered recovery behavior | That all stateful services now share one resilience story |
| Compose remains the readable human contract while the control surface grows | Added coordination actually pays down the right pain instead of becoming another opaque platform | Partial | Root Compose remains the real authoring and deploy surface, and the repo keeps refusing premature worldview capture | After adding any helper layer, prove an operator can still answer what runs where, why, and how the route was chosen | That a larger scheduler or HA layer has already earned whole-stack promotion |

## The drill ladder the repo should speak in publicly

The repo should not speak about drills as if "tested" were one thing.
There are levels, and the whole point is to stop later levels from being
implied by earlier ones.

The public language rule should be mechanical:

- if the strongest drill was level 2, do not speak with level 3 nouns
- if the strongest drill was level 3, do not imply level 4 survivability
- if the strongest drill was level 4, do not imply level 5 semantic parity
- if the strongest drill was level 5, do not imply level 6 stateful honesty

### Drill level 0: architecture paper

What happened:

- the intended behavior is described clearly

Allowed claim:

- the target operating contract is explicit

Forbidden upgrade:

- the runtime already does this

### Drill level 1: config-shape proof

What happened:

- the tracked config contains ingredients for the behavior

Allowed claim:

- the implementation is authored in this direction

Forbidden upgrade:

- the behavior survives runtime or failure

### Drill level 2: happy-path runtime proof

What happened:

- the path or service worked in nominal conditions

Allowed claim:

- the tested path works in normal conditions

Forbidden upgrade:

- the path is resilient

### Drill level 3: wrong-node proof

What happened:

- the request intentionally landed on a node that did not host the service
  locally and still completed

Allowed claim:

- wrong-node success exists for this exact route and service class

Forbidden upgrade:

- the whole stack now preserves wrong-node requests

### Drill level 4: backend-loss proof

What happened:

- the preferred local backend disappeared and the path still survived

Allowed claim:

- this route survives the specific failure it claims to absorb

Forbidden upgrade:

- all policy and middleware behavior remains identical unless checked

### Drill level 5: semantic continuity proof

What happened:

- local and peer-fallback behavior were compared where users would actually
  notice

Allowed claim:

- this request path remains semantically stable under handoff

Forbidden upgrade:

- stateful systems are now solved too

### Drill level 6: stateful correctness proof

What happened:

- the relevant state topology survived with correct ownership and client
  behavior

Allowed claim:

- this specific stateful class has a meaningful resilience story now

Forbidden upgrade:

- the whole platform is now generically anti-SPOF

The repo should keep using these levels in public language because they slow
down exactly the kind of vague summarization the user hates.

They force every “we proved X” statement to answer:

- proved under what condition
- for what route or class
- with what remaining ceiling

That is what turns proof from a badge into an explanatory surface.

The repo does not only need to know that something passed.
It needs to know what private burden the passing proof actually removed.

The winning question after any drill is not only:

- did it pass?

It is:

- which reason for operator reconstruction just became less true?

That last sentence is a stricter requirement than it sounds.
A drill is stronger here when it removes one more reason the operator has to
silently complete the architecture from memory after the command finishes.

## Recommended proof order

The repo should not try to prove everything at once.
The right order is the one that pays down the user's actual frustration chain
instead of chasing the most flattering demo first.

That is why the recommended order looks narrower than a normal platform
roadmap.

It is optimized for recovering the missing truth layer, not for maximizing how
many green checks can appear early.

That optimization should stay visible because it explains why the proof order
is intentionally narrower and less flattering than a normal demo roadmap.

### 1. Prove one representative stateless HTTP route

Pick a route that is:

- publicly reachable
- easy to correlate in logs
- not state-critical
- ideally already behind meaningful auth or middleware

Why first:

- this is the cleanest honest path to proving any-node entry, locality, and
  wrong-node preservation without cheating on state

It is also the fastest way to force the system to answer whether it has started
owning placement truth outside the operator's head.

### 2. Prove backend-loss persistence for that same route

Why second:

- otherwise the repo can still accidentally prove wrong-node success only in a
  version of the world where the preferred local backend never actually
  disappeared

### 3. Prove semantic continuity for one protected route

Why third:

- a surviving `200` from a peer is not the same thing as preserving the
  route's meaning

### 4. Expose or prove real placement truth

Why fourth:

- otherwise too much of the multi-node story still depends on human memory
  rather than auditable current state

That ordering also resists a quieter failure:
adding a registry-shaped surface before the repo has been forced to prove what
that surface must actually explain during wrong-node or backend-loss pressure.

This ordering is useful precisely because it resists a common simplification:
that a registry surface should be discussed abstractly before the repo has
proved what that surface must be able to explain under stress.

The page should keep reminding the reader:

- a registry file is not yet a trusted truth layer just because it exists
- it becomes one only when a stressed route can be explained from it

### 5. Move into one stateful class at a time

Suggested order:

1. Redis
2. MongoDB
3. RabbitMQ
4. storage and shared data surfaces

Why:

- each class has different ownership, election, reconnect, and storage
  semantics
- flattening them into one "stateful HA" checkbox would recreate exactly the
  fake closure the user is trying to escape

## Questions this page is supposed to make more uncomfortable

If this page is doing its job, it should make it uncomfortable to say any of
these sentences without much stronger evidence:

- wrong-node requests work
- failover is solved
- the service registry is live
- the nodes are converged
- the databases are resilient
- the control surface is now mature enough

It should also make it uncomfortable to mistake emotional relief for proof.

A cleaner dashboard, a more elegant generator, or a more modern helper layer
may make the stack feel better.
This page exists to keep "feels more platform-like" from being mistaken for
"has actually reduced the hidden truth burden."

If a future change really wants to say one of those things, it should first
update the relevant row above with:

- stronger evidence
- the exact drill that was run
- the exact failure class that was exercised
- the exact limit of what that drill still does not prove

## Bottom line

This page exists because the repo is not trying to become a better liar.
It is trying to become harder to fool.

The real win is not "we ran more tests."
The real win is:

> the available options stop collapsing into theater, the proof boundaries stay
> visible, and the user can finally tell which parts of the platform are real
> and which parts are still only architecture pressure

That is the standard the rest of the docs should keep inheriting.

## Related pages

Use this page together with:

- [`../architecture/operator-contract.md`](../architecture/operator-contract.md)
- [`../architecture/failure-model-and-maturity.md`](../architecture/failure-model-and-maturity.md)
- [`devops-runbook.md`](devops-runbook.md)
- [`source-assimilation-index.md`](source-assimilation-index.md)
