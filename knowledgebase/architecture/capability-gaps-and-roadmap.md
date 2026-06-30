# Capability Gaps and Roadmap

For the evidence shaping this roadmap, start with:

- [Evidence Ledger](../research/evidence-ledger.md)
- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
- [Stateful HA Evidence](../research/stateful-ha-evidence.md)
- [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This page is not a feature wishlist.

It is the map of which missing truths still force the rest of the platform to
overclaim relative to the user's real benchmark.

The user is not asking for a tidy next-steps board.
They are asking why so many infrastructures keep pretending there are rich
options while quietly forcing one of two humiliations:

- stay in ordinary Docker forever and accept wrong-node fragility
- jump into a heavyweight control plane before smaller honest layers were even
  exhausted

So the only roadmap worth keeping here is the roadmap that answers:

> which missing truth or recovery layer is still forcing the platform to lean
> on hope, memory, stale assumptions, or theatrical confidence?

That is what "roadmap" means in this repo.

## What this page is and is not allowed to prove

This page is authoritative about:

- prioritizing missing truth layers in the order that protects honesty
- explaining which unresolved gaps still force overclaim if left untouched
- connecting current runtime evidence to the next proof-bearing promotion step
- distinguishing valuable next work from already-proven capability
- keeping stateless HTTP, protected HTTP, TCP, and stateful promotion on
  different tracks

This page is not authoritative about:

- acting like a completion report
- implying that a well-ordered roadmap means the runtime is already coherent
- flattening all recovery classes into one generic queue
- upgrading planning clarity into present-tense runtime maturity

This page is a sequencing contract.
It is not an implementation victory lap.

## Strongest honest current answer

The next work is not "whatever seems generally useful for self-hosting."

The next work is whichever missing truth layer still forces the docs and the
runtime story to rely on:

- hope
- operator memory
- stale topology assumptions
- rhetorical glue
- or helper-language that sounds stronger than the proof it owns

In the priority implementation today, that means:

- placement truth
- convergence truth
- route durability under real failure
- peer eligibility
- stateless wrong-node proof
- protected-route continuity
- keeping TCP and stateful promotion under separate honesty gates

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
9. only then decide what has actually earned promotion into a stronger control
   layer

This is not motivational sequencing.
It is the dependency chain between the dream and reality.

## What still does not count as progress

This repo has a specific problem, so it also needs a specific false-progress
filter.

The following may all be useful work and still not count as roadmap closure by
themselves:

- adding more public nodes without proving wrong-node meaning survives
- adding richer proxy logic without proving route persistence under failure
- adding more healthchecks without proving peer eligibility semantics
- adding secret-sync or file-sync helpers without showing how they gate
  substitution trust
- standing up Nomad, OpenSVC, k3s, or other controller experiments without
  proving which hidden burden they actually removed
- making docs calmer, cleaner, or more enterprise-sounding while the same weak
  runtime truths remain

That list exists because the user is explicitly frustrated by ecosystems that
keep narrating partial machinery as if the burden has already moved.

## Why the order has to stay this strict

The project wants all of these at once:

- multiple public nodes
- local-first service handling
- peer fallback when the request lands on the wrong node
- no fake HA language
- no forced Swarm or Kubernetes jump before necessary

That dream depends on a stack of truths.
If the lower truths are weak, the upper ones become performance.

The real dependency chain is:

1. know where services actually live
2. know whether peers are semantically aligned enough to substitute
3. keep the recovery route alive when the local backend dies
4. know which peers are eligible to receive traffic now
5. prove that policy, auth, and middleware survive the handoff
6. keep stateful classes under harsher rules than stateless HTTP

If that chain breaks at step 1 or 2, later success claims are already suspect.

## What this roadmap is really protecting

This roadmap is not protecting a tool preference.
It is protecting the user's real demand:

> stop pretending there are plenty of options when most supposed options either
> collapse under wrong-node pressure or force a giant control plane before the
> smaller honest answers have even been exhausted

That is why the roadmap stays narrow and harsh.

The next thing to build is not the next thing that sounds advanced.
It is the next thing whose absence still forces the stack to lie.

This also means the roadmap must actively reject seductive but misordered
moves:

- promoting to a bigger platform before placement and convergence truth are
  explicit
- narrating ingress sophistication as if it settles stateful correctness
- treating helper growth as "still just Compose" after helpers start owning
  scheduler-like truth silently
- proving happy-path forwarding before proving that the route survives backend
  death

## Read this page as integrity gates, not aspirations

Each priority below is written as:

- what hidden burden it addresses
- what the current worktree proves
- what failure signature still survives
- what proof threshold would actually close that gap
- what that closure would unlock next

That format matters because the repo does not need more dreams.
It needs clearer rules for when a bigger claim becomes legal.

## Priority 0: documentation honesty

Class:

- integrity gate

Why it is first:

If the docs overclaim, every later decision gets poisoned.
This repo cannot afford another smooth narrative that upgrades intent into
"basically working."

Current truth:

- the knowledgebase is far stricter than the older docs
- it can still regress whenever a page starts sounding like a generic HA guide
  instead of a pressure-tested reading of the worktree and archive

Failure signature:

- DNS plurality described as end-to-end failover
- helper layers described as live truth before runtime proof exists
- stateful systems described as HA because they are reachable
- route generation described as recovery without route-persistence proof

Proof threshold:

- every operator-critical page distinguishes live runtime, architecture
  intent, planning pressure, and archive pressure
- every major resilience claim makes its proof class obvious
- every weaker state stays named as weak instead of being widened into
  "solved"

Unlocks:

- trustworthy sequencing
- less self-deception
- cleaner promotion decisions later

This priority sounds editorial.
It is actually architectural, because dishonest docs make every next decision
worse.

## Priority 1: placement truth

Class:

- foundational

What this really means:

The receiving node needs a current answer to:

> is the requested service local, and if not, where does it actually live right
> now?

Without that answer the whole dream collapses into guesswork.

Current truth:

- the repo repeatedly converges on `services.yaml` or an equivalent
  current-state surface as the intended lightweight registry
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

## Why the roadmap keeps stateless HTTP as the first real proof lane

The repo absolutely cares about TCP and stateful workloads.
It still has to start with stateless HTTP because that is the narrowest class
where the core wound can be tested without borrowing credibility from a larger
claim than the system has earned.

Stateless HTTP is the first place where the repo can prove all of these at
once:

- wrong-node entry happened
- placement truth existed
- a peer was chosen for a reason
- the route survived local backend loss
- the request still meant the same thing to the user

That is not the final dream.
It is the first point where the dream stops being theoretical.

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

- the first real claim that "fallback" means more than "there was another idea
  on paper"

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
- the master plan contains sync-agent, peer broadcast, and failover-agent
  directions
- the current runtime still does not visibly prove a shared eligibility
  decision surface

Failure signature:

- the node knows another peer exists
- the node can technically reach that peer
- the node still cannot answer whether that peer is the right recovery target
  now

Proof threshold:

- peer selection is based on shared current truth, not remembered placement,
  guessed health, or stale topology lore

What counts as evidence:

- an explicit eligibility rule set
- visible health and alignment inputs
- proof that the selected peer was chosen from that rule set

Unlocks:

- stronger wrong-node drills
- less fragile peer-forward language
- a real decision surface for whether a bigger controller has earned itself

## Priority 5: one real stateless HTTP wrong-node path

Class:

- first end-to-end proof lane

What this really means:

The repo needs at least one named route where the wrong-node story stops being
speculative and becomes something the system can actually do.

Current truth:

- the dream is explicit
- the edge stack is serious
- the docs are already harsh about what does not count
- no generic wrong-node success is being claimed yet

Failure signature:

- every important sentence about wrong-node behavior remains architecture
  intent rather than runtime proof

Proof threshold:

- intentionally land on a healthy node that does not host the route locally
- show that the request still completes correctly
- show why it succeeded from inspectable system truth rather than folklore

What counts as evidence:

- controlled entry target
- visible local-versus-remote decision
- peer-choice explanation
- user-visible route success

Unlocks:

- the first legitimate "this platform feels different now" claim
- a harder baseline for later helper or orchestrator promotion

This is the first place where "multi-node" can stop sounding emotional and
start sounding operational.

## Priority 6: protected-route continuity on the same path

Class:

- semantic continuity gate

Why it follows stateless HTTP:

Protected routes are stricter.
The repo should not claim protected wrong-node success before proving plain
HTTP wrong-node success first.

Current truth:

- TinyAuth, nginx auth extensions, CrowdSec, and Traefik middleware are all
  materially present
- the worktree still does not prove local and forwarded policy meaning are the
  same

Failure signature:

- the forwarded route returns a response
- but auth, middleware, headers, or rewrites no longer behave like the same
  protected service

Proof threshold:

- compare local and peer-forwarded behavior for one protected route
- show that policy meaning remains intact rather than merely "reachable"

What counts as evidence:

- auth behavior before and after handoff
- middleware continuity
- user-visible parity of the protected surface

Unlocks:

- stronger confidence in L7 peer-forward claims
- a real reason to say the edge stack is preserving service meaning, not just
  transport

## Priority 7: TCP and stateful honesty gates

Class:

- separate harsh lanes

What this really means:

The repo must not let HTTP maturity contaminate TCP and stateful claims.

Current truth:

- TCP exposure exists
- stateful services exist
- the docs already reject calling them HA by adjacency

Failure signature:

- "we solved wrong-node HTTP" mutates into "the platform is highly available"
- a routed TCP endpoint gets overread as resilient ownership
- a stateful service gets promoted mainly because it can be reached

Proof threshold:

- per service class, define transport, ownership, replication, promotion, and
  client-behavior truth separately

What counts as evidence:

- explicit write-authority model
- explicit promotion and recovery model
- route drills that do not overclaim beyond their class

Unlocks:

- honest stateful planning
- cleaner future comparisons between helper layers and stronger orchestrators

## Priority 8: only then ask whether a stronger control layer has earned itself

Class:

- promotion gate

Why it is last:

The repo should not ask "which orchestrator wins?" before it can answer the
smaller question:

> which hidden burden is still unsolved after the narrow honest layers were
> tried?

Current truth:

- the repo is still Compose-first by default
- it is exploring OpenSVC, k3s, Nomad, helper agents, and registry ideas
- no broader control family has yet earned default promotion for the whole
  stack

Failure signature:

- a bigger tool gets chosen because it sounds adult
- the same hidden truths survive under a cleaner abstraction boundary

Proof threshold:

- the remaining unsolved burden is specific and named
- the promoted layer owns that burden better than the current stack
- the repo can explain what worldview tax it is paying and why it is worth it

Unlocks:

- justified promotion instead of prestige capture
- clearer limits on what should remain Compose-owned

## What not to build out of order

The roadmap should explicitly reject these out-of-sequence moves:

- proving bigger-cluster demos before exposing placement truth
- refining failover rhetoric before route-persistence proof
- celebrating mesh reachability before peer-eligibility truth
- broad orchestrator comparison before the missing burden is named precisely
- softening stateful language because the ingress layer got stronger

Those moves are not always useless.
They are misordered relative to the user's real problem.

## What a passing proof packet looks like for this roadmap

Every later promotion argument should be able to point to a proof packet rather
than only to a narrative.

For this repo, a strong proof packet usually means:

- the exact route or service class being discussed is named
- the node that received the request is named
- the local-versus-remote decision is visible
- the placement or eligibility truth surface that informed that decision is
  visible
- the failure event being tested is explicit
- the route persistence result is visible after the failure
- the user-visible result is captured
- the page stays explicit about what broader class was **not** proven

Without that structure, "progress" in this repo too easily collapses back into
coherent description of a future someone hopes will exist.

## The practical interpretation

If the repo has limited attention, the next highest-value work is not the most
technically impressive change.

It is the change that most reduces one of these humiliations:

- the wrong node still needs remembered topology
- the candidate peer still needs guessed semantic alignment
- the fallback route still disappears when the backend dies
- the protected route still changes meaning after handoff
- the stateful story still depends on optimism instead of authority

That is how the roadmap should be read.

## The brutal final roadmap question

Before any roadmap item is promoted, force it through this exact question:

> if this lands, which specific hidden burden becomes more system-owned and
> less operator-reconstructed?

If the answer is vague, the work may still be interesting.
It is not yet roadmap-critical for this repo.

## Bottom line

This roadmap is not trying to make the project look well managed.

It is trying to stop three things:

- helper sprawl being mistaken for solved truth
- bigger platforms being chosen before the smaller honest layers were exhausted
- polished sequencing being mistaken for runtime maturity

The roadmap only has value if it keeps the burden-transfer question brutal.
