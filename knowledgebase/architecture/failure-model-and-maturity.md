# Failure Model and Maturity Matrix

This page is the hard reality check for the architecture.

It exists to answer the pressure behind most of the user’s questions:

> what actually breaks today, what real protection is present, what is only
> intended, and what exact evidence would be required before stronger language
> stops being a lie?

That pressure matters because this repo is trying to escape several bad
infrastructure habits at the same time:

- saying HA when the proof is only “a proxy exists”
- saying multi-node when the proof is only “the stack was copied to more than
  one machine”
- saying dynamic when the proof is only “config is generated from more inputs”
- saying resilient when ingress, placement truth, convergence truth, and
  stateful correctness are all maturing at different speeds

This page exists to stop those upgrades from happening for free.

It also has to do something more procedural than older maturity pages:

- define which domain is being judged
- say which truth layer is actually live there
- name the hidden operator tax still present there
- name the next honest maturity step without narrating two steps at once

It also exists to stop a subtler failure:

the stack can become more mature in vocabulary faster than it becomes mature in
truth-handling.

That is one of the user's deepest frustrations.
The surrounding ecosystem is full of systems that can explain themselves with
the right words long before they can preserve a request honestly under wrong-
node entry, backend loss, or stateful failure.

This matrix is here so the repo does not make that same move against itself.

It also exists because maturity language is one of the easiest ways for docs
to become flatter than the truth.

If this page becomes:

- a feature inventory
- a roadmap in table form
- or a smoother version of the architecture pages

then it has stopped doing its actual job.

Its job is to preserve the unevenness between domains strongly enough that the
reader cannot accidentally believe the stack matures as one thing.

It is also here because the user is not mainly asking whether features exist.
They are asking whether those features stop the system from collapsing back
into hidden human interpretation the moment stress becomes real.

That is why this page has to stay harsher than a normal maturity matrix.

## What this page is and is not allowed to prove

This page is authoritative about:

- which failure domains define the repo's real maturity problem
- how unevenly those domains can mature without justifying broader comfort
- where hidden operator reconstruction is still part of the live failure model

This page is not authoritative about:

- whether a specific route or workload has already passed its harder drills
- whether a future control-plane direction has already earned promotion
- whether maturity language should override narrower runtime or proof pages

This page is a domain-by-domain honesty matrix, not a substitute for route-
level or topology-level proof.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "this domain is still immature in a specific way" | synthesis over runtime, evidence, and operator burden | current runtime pages, proof matrix, research evidence, this page's domain framing | that nothing useful exists in the domain |
| "one domain matured faster than another" | maturity synthesis | this matrix plus supporting evidence pages | that the stack matures as one clean platform |
| "hidden operator reconstruction is still part of the failure model" | contradiction plus proof analysis | this matrix, runbook, proof pages, intent surfaces | that the operator burden has no partial reductions anywhere |
| "the next honest maturity step is X" | roadmap plus proof discipline | domain row, supporting evidence, promotion rules | that the next step is already complete or universally sufficient |

If a maturity sentence starts sounding like a global resilience verdict, this
page is being overread.

## Strongest honest current answer

The strongest honest current answer is that the repo does not mature as one
platform. It matures as several uneven truth lanes that can easily be confused
for one another. HTTP ingress, wrong-node recovery, peer-selection truth,
operator reconstruction burden, and stateful correctness are progressing at
different rates. Any maturity story that compresses them into one comfort level
is already flatter than the evidence.

## What this page is really measuring

This is not a generic infrastructure maturity matrix.

It is a maturity matrix for one very specific dream:

> several ordinary Docker nodes behave like one request-preserving,
> operator-readable personal cloud without importing a heavyweight orchestrator
> before that tax is clearly justified

Every row should be read through that lens.

If a capability sounds impressive but does not move that dream forward, it is
secondary here.
If a capability sounds mature but still leaves the operator carrying hidden
truth in their head, it is still immature here.

If a capability makes the documentation sound more finished while the runtime
still depends on remembered exceptions, it is still immature here.

That rule matters more than the label names themselves.

The labels are just compression.
The real content is whether the human operator is still silently performing the
missing control-plane work.

That means each row should be interpreted as:

- what the system itself owns
- what the operator still owns
- what the docs must not quietly credit to the system yet

That last rule matters because this repo is especially vulnerable to maturity
theater.

It already has:

- real edge components
- real Compose structure
- real plans
- real option exploration

All of that is enough to make a weaker doc set sound mature long before the
user's real benchmark has been crossed.

## The hidden metric is operator reconstruction tax

Most maturity pages only ask:

> does the feature exist?

This repo has to ask a harder question:

> how much human reconstruction is still required before the feature becomes
> usable honestly?

That is the real anti-SPOF question here.

The operator is still part of the failure model if the system only behaves
coherently because the operator privately remembers:

- what runs where
- which peer is currently valid
- which route is real
- which “HA” claim is only social knowledge
- which fallback path is still theoretical

So this page is measuring both:

- runtime maturity
- how much hidden human control-plane labor is still being performed by hand

That second metric is not optional. It is the entire reason the user remains
frustrated even when the stack already contains real software, real routes, and
real documentation.

The user is not frustrated because nothing exists.
The user is frustrated because too many things exist in a way that still
requires a private human interpreter before the system becomes coherent.

That sentence should color every row in the matrix.

The opposite of maturity in this repo is not emptiness.
It is sophistication that still depends on folklore.

This page should therefore be treated as a matrix for folklore-removal as much
as for feature presence.

That is the central reading rule for this matrix.

It should also be used against the rest of the docs:

- does this page show real maturity?
- or did it just describe more sophisticated-looking components while folklore
  remains the real join layer?

## The fake-option trap this matrix is trying to prevent

The archive keeps circling one recurring disappointment:

- small-node Docker gives direct control but keeps wrong-node fragility
- “just use DNS” helps first hop but not preserved service meaning
- “just use a reverse proxy” helps routing but not cross-node truth
- “just use Kubernetes” or “just use Nomad” often arrives before the missing
  capability has been isolated

That means the user keeps getting offered options that are not really options.
They are either:

- weak partial answers that still require sacred-node memory
- or giant worldview imports that demand trust before they have earned it

The point of the matrix is not merely to rank components.
It is to identify exactly where the repo is still trapped between those two bad
choices.

Some futures really do offer new truth ownership.
Others mostly offer a better story about the same hidden burden.
This matrix exists to keep those two futures from looking equally mature.

That means the matrix has to stay tied to bad-day behavior.

If a capability looks mature only while:

- traffic lands locally
- the preferred backend still exists
- the state owner has not disappeared
- or the operator is still present to explain the weird exception

then the maturity claim is still too flattering.

That means the matrix is not neutral about "options."
It is explicitly trying to separate:

- real additional truth ownership
- from answers that merely sound like optional futures while leaving the same
  hidden burden intact

This matrix is here to keep the repo from recreating that same fake choice in
its own docs.

## The failure domains that define the whole repo

Before the matrix, name the dominant failure domains directly.

### Failure domain 1: edge-path failure

Questions:

- can a request reach a healthy public node?
- can that node parse, challenge, and route the request coherently?
- do auth and middleware still mean the same thing under failure?

### Failure domain 2: placement-truth failure

Questions:

- does the system know where services actually live right now?
- is routing consuming auditable truth or social knowledge?
- is the operator still the hidden scheduler?

### Failure domain 3: convergence-truth failure

Questions:

- do nodes agree on env, secrets, revision, and helper-generated state?
- is a reachable peer actually eligible to receive forwarded traffic?

### Failure domain 4: path-persistence failure

Questions:

- if the preferred local backend disappears, does the recovery path remain?
- or does the mechanism that promised failover delete the route needed to use
  it?

### Failure domain 5: semantic-continuity failure

Questions:

- if a request lands on the wrong node and gets handed to a peer, does it still
  behave like the same service path?
- or does auth, middleware, or visible policy silently change?

### Failure domain 6: stateful-correctness failure

Questions:

- who owns writes?
- what promotion or election semantics exist?
- what reconnect behavior survives node loss?
- what data becomes ambiguous or stranded?

The user’s dream only becomes real when these domains stop drifting apart.

If ingress looks advanced but placement truth is still social knowledge, the
dream is not real yet.
If peer forwarding exists but eligibility truth is still guessed, the dream is
not real yet.
If TCP routes exist but write ownership remains ambiguous, the dream is not
real yet.

## How to read the matrix

### Truth layers

Each row below distinguishes among:

1. live root Compose truth
2. planned multi-node architecture truth
3. research-pressure truth from the repo’s accumulated frustration and archive

If those layers collapse together, the matrix becomes smooth and useless.

Smoothness is the danger here.
This repo does not need prettier maturity language.
It needs language that preserves where the system is still split across live
runtime truth, planned truth, and operator-held truth.

There is also a practical reading rule:

- `Present` means materially live in a recognizable form, not emotionally
  solved
- `Partial` means some burden relocated, important burden still socially held
- `Planned` means the repo knows the missing truth, but the runtime does not
  own it yet
- `Absent` means important to the dream, not materially system-owned today

This is another place where the docs should act more like held-out test cases
than like architecture marketing.

The operator-held truth part is the most important one.

Most infrastructure matrices ignore it.
This one cannot, because operator reconstruction tax is one of the main things
the repo is trying to eliminate.

## What should move a row upward

A row should move upward only when:

- a stronger proof class was earned for the relevant domain
- the system now owns more of the explanation for success under stress
- the operator has to privately reconstruct less of the governing truth

A row should not move upward merely because:

- another helper exists
- another future was documented
- another path sounds more serious
- another layer can now be named

### Confidence levels

- `High`
  - The current worktree directly supports the statement.
- `Medium`
  - The direction is well evidenced, but implementation or proof is incomplete.
- `Low`
  - Mostly exploratory, aspirational, or weakly verified.

### Maturity labels

- `Present`
  - Materially live in the tracked implementation in a recognizable form.
- `Partial`
  - Real pieces exist, but failure behavior is not trustworthy enough yet.
- `Planned`
  - The behavior is described clearly, but not materially live as dependable
    runtime truth.
- `Absent`
  - Important to the dream, but not materially present in tracked current
    state.

## Maturity matrix

| Domain | Current maturity | Confidence | Live truth today | Hidden operator tax still present | What failure still looks like | What would actually prove the next honest maturity step |
| --- | --- | --- | --- | --- | --- | --- |
| Compose authoring surface | Present | High | Root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml) plus included fragments remain the main human contract | The operator still has to mentally reconstruct a distributed system from many fragments | Authoring drift and conceptual duplication accumulate faster than shared runtime truth | A higher-level authoring surface that emits the same inspectable live shape without hiding the system |
| Compose as real deployment entrypoint | Present | High | Compose is still the actual merge and deploy surface, not just rhetoric | Multi-node meaning still depends on planning, helper logic, and human reconstruction outside Compose itself | Host loss and relocation still collapse into manual choreography or brittle glue | A thin proven coordination layer that preserves Compose readability while adding trustworthy node-aware behavior |
| Compose as full distributed control plane | Partial | High | Compose coordinates local runtime shape well, but not shared placement truth, convergence truth, or distributed failover truth | The operator is still the hidden scheduler and sometimes the hidden registry | The stack can look unified while remaining node-local where it matters, which is exactly how fake options survive | Proof that distributed truth and failure handling work without a hidden heavyweight scheduler |
| Placement truth / current-state registry | Absent | High | `services.yaml` is central in architecture intent, but there is no live tracked root [`services.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/services.yaml) powering the priority runtime | The operator still privately answers “what runs where right now?” | Wrong-node entry can still arrive before any formally tracked current-state answer exists, so the human remains the final source of topology truth | A live tracked registry or equivalent truth surface that routing logic actually consumes |
| Repo, env, and secret convergence truth | Planned | High | [`.env.example`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.env.example), `${SECRETS_PATH}`, file-backed secrets, and sync-agent style planning all exist | The operator still has to infer whether peers are semantically compatible | Forwarded requests can land on a reachable peer that is stale, mismatched, or half-prepared, which turns “reachable” into another fake safety signal | Auditable sync and drift detection, visible node-by-node readiness, and explicit restart semantics |
| HTTP ingress presence | Present | High | Traefik, TinyAuth, CrowdSec, nginx extensions, docs routing, and many service hostnames are materially live | The operator still has to distinguish ingress presence from request preservation | Public ingress can exist while distributed request meaning is still fragile | Endpoint-level proof that representative HTTP paths are healthy and policy-correct under normal operation |
| Any-node first hop | Partial | Medium | Cloudflare/DDNS and anti-sacred-node intent are strong and repo-native | The operator still has to prove that the surviving node can do something useful after receiving the request | DNS may hand traffic to a surviving node that cannot preserve the route | Verified multi-node first-hop tests for representative HTTP paths |
| Local-first service semantics | Planned | Medium | The architecture contract clearly prefers local service when available | Locality still has to be inferred rather than proven through auditable truth | A node may serve locally today, but the system cannot yet explain that decision cleanly | Demonstrated local serve behavior tied to explicit current-state locality truth |
| Wrong-node HTTP preservation | Partial | Medium | The dream is explicit and the edge stack is mature enough to make it a real target | The operator still has to reconstruct whether the receiving node knew the service was remote, who the peer was, and why | A request landing on the wrong node may still depend on static assumptions or incomplete route truth, which means the user still has fewer real options than the docs want to imply | One representative stateless HTTP route proven through wrong-node entry with logs showing peer selection and successful handoff |
| Fallback route persistence after backend loss | Partial | High | `docker-gen-failover` is live and also documented as a route-deletion risk | The operator still has to fear that the failover mechanism removes the path it promised to preserve | The local backend disappears and the fallback route disappears with it | A repaired or replaced route-generation mechanism proven under backend-loss drills |
| Middleware and auth continuity under peer handoff | Partial | Medium | TinyAuth, CrowdSec, and middleware chains are materially real | The operator still cannot assume local and peer-forwarded behavior mean the same thing | A request may still return success while challenge, auth, or policy semantics drift, producing the most dangerous kind of false confidence | Failure testing showing local and peer-forwarded behavior match where users can notice |
| Observability presence | Present | High | The stack has real monitoring, metrics, logging, exporters, and dashboards | Observability can describe failure sooner than the stack can safely absorb it | The operator learns more clearly that the system is failing, but still has to perform the recovery manually | Evidence that signals drive trusted remediation or at least unambiguous operator decisions |
| Node revision visibility | Planned | Medium | The repo clearly cares about sync and convergence, but does not yet expose one clear revision-truth surface | The operator still has to infer whether nodes are equivalent enough for peer forwarding | Different nodes can behave differently while the docs still speak as if the runtime were singular | Visible node revision and deployment-state reporting tied to affected-service semantics |
| Service redeploy and relocation | Planned | Medium | The dream of re-establishing services elsewhere is clear in planning and architecture docs | The operator still carries the burden of route surgery and recovery choreography | Host loss still means manual intervention or fragile helper logic, so the platform still depends on remembered recovery rituals | Proof that a representative stateless service can be restored or moved without silent policy drift |
| TCP ingress presence | Present | High | TCP routers for Redis, MongoDB, and biodecompwarehouse-related services are live in the root stack | The operator must keep reminding everyone that TCP exposure is not TCP resilience | A socket stays reachable while application topology or data semantics remain unsafe | Protocol-aware tests that separate simple exposure from meaningful continuity |
| TCP failover semantics | Planned | Medium | The repo intentionally keeps L4 separate from HTTP and clearly cares about raw TCP routing | The operator still has to define what good failover means for each TCP class | A routed TCP service may remain reachable while still being wrong or unsafe | Service-class-specific L4 designs plus real failure drills |
| Redis resilience | Planned | High | Redis is live and the planning layer already treats Sentinel-like or equivalent truth as the honest threshold | The operator still has to privately answer who is primary and what clients should rediscover | Redis can remain a single-node trust failure dressed up as routed infrastructure | Proven primary and replica semantics, failover, and client reconnect behavior |
| MongoDB resilience | Planned | High | MongoDB is live and replica-set semantics are already named as the minimum honest story | The operator still has to infer election and discovery behavior | MongoDB can still be a SPOF or a fragile client-discovery surface | A working replica set with verified election and client discovery under node loss |
| RabbitMQ resilience | Planned | Medium | RabbitMQ is part of the live runtime but not yet documented as a proven resilient queue domain | The operator still has to privately know whether queue durability is clustered or merely hoped for | Queueing can still depend on one node or weak persistence assumptions | Either an explicit bounded single-node contract or clustered recovery behavior with proof |
| Shared and replicated storage strategy | Low | Medium | Many node-local volume and bind-mount patterns exist; some sync tooling exists, but no trustworthy shared-state mobility story yet | The operator still carries the burden of knowing which data can move, replicate, or strand safely | Node loss can strand data or restart a service on the wrong data surface | A storage taxonomy that separates disposable cache, replica-safe state, and genuinely critical shared state |
| Headscale and private mesh continuity | Partial | Medium | Headscale is live and strategically central because peer-forwarding needs private inter-node connectivity | The operator still has to ask whether the private coordination path itself has become a hidden sacred dependency | The anti-SPOF story may still contain a quiet coordination SPOF | Either explicit redundancy or an explicit bounded acceptance of the remaining dependency |
| Update automation | Partial | High | `watchtower` is live and also documented as non-functional or untrustworthy in current form | The operator still has to watch the updater instead of trusting it | Updates can fail silently or break continuity in a way worse than having no automation | Health-aware rollout behavior, visible restart semantics, and a rollback story |
| Control-plane optionality | Present | High | Compose, CUE, OpenSVC, Nomad, k3s, and Kubernetes futures all remain intentionally open | The operator still has to decide where extra complexity belongs instead of being handed one forced answer | Optionality can decay into indecision if no domain-specific promotion triggers exist | Explicit promotion rules showing which pain class earns which stronger layer |

## The deepest current truths

These are the conclusions the matrix keeps pushing toward.

### 1. The dream is more mature than the live truth surface

The repo already knows, with unusual clarity, what kind of system it wants:

- any-node entry
- local-first serving
- peer-forward fallback
- lightweight current-state truth
- no mandatory heavyweight orchestrator

Several of the mechanisms required to make that true are still:

- missing
- brittle
- only partially implemented
- or only planned

That means the repo is in an unusual but honest state:

- the dream is not vague
- the runtime is not fake
- the middle layers are still underbuilt relative to the dream

That combination is exactly why the docs must stay strict.

If they become flattering here, they will turn genuine progress into fake
closure.

That is not a failure of imagination.
It is the proof boundary the docs must keep visible.

The repo is not confused about what it wants.
It is still trying to earn the right to say that the runtime itself understands
that desire instead of only the humans around it.

That is close to the core maturity test for the whole knowledgebase:

is the runtime itself becoming a better explainer of the dream, or are the
humans merely becoming better explainers around the runtime?

### 2. Edge sophistication is ahead of shared truth

Ingress, auth, middleware, and observability are already serious.
Placement truth and convergence truth are less mature.

That means the stack can look more cluster-like from the outside than it really
is internally.

This is one of the biggest reasons the repo keeps feeling like it has lots of
pieces but not enough trustworthy options.

The user's frustration is not "there are too few tools."
It is "there are too many tools that still fail to remove the same hidden human
burden."

That is why real components alone do not normalize these rows.
Real components can still participate in fake-option maturity theater.

### 3. HTTP maturity is ahead of stateful maturity

The repo has a clearer path to proving stateless HTTP wrong-node success than
it does to proving honest Redis, MongoDB, RabbitMQ, or shared-storage
resilience.

That separation has to stay explicit because the user is specifically
frustrated by architectures that blur them together.

### 4. Presence is not the same as interchangeable continuity

Many services are live.
Many hostnames are live.
Several TCP routes are live.

What remains incomplete is the stronger promise:

> a request path remains semantically interchangeable across nodes under failure

That stronger promise is the thing the user actually cares about.

### 5. Hidden operator tax is part of the maturity score

This matrix is not only measuring whether features exist.
It is also measuring how much implicit control-plane labor the operator is
still performing by hand.

That is why `Present` never automatically means good enough.

The user’s frustration is not irrational. It is what it feels like when a
system looks distributed but still relies on private human glue at the most
important moments.

This point matters enough to repeat plainly:

if the system still needs the operator to remember the truth before the system
can act correctly, then the operator is still a SPOF even if several machines
are technically online.

That line should be treated as the matrix's harshest invariant.

## What to trust today

Trust most:

- the existence and seriousness of the live Compose topology
- the maturity of ingress and observability surfaces
- the clarity of the repo’s anti-heavy-orchestrator dream
- the fact that stateful HA is already treated as a separate harder problem

Trust less:

- any blanket zero-SPOF claim
- any implication that a live shared registry already powers placement truth
- any suggestion that DNS failover alone equals service continuity
- any suggestion that routed TCP datastores are already resilient because they
  are reachable

## The matrix’s most important correction to older docs

One concrete example of why this page exists:

- older prose could easily say “OpenSVC plus generated ingress plus HAProxy”
  as if those were one clean maturity step
- the current worktree shows those truth surfaces are still mixed together more
  loosely than that

For example:

- the HTTP generator reads node names from OpenSVC membership
- the current L4 generator reads node IPs from `tailscale status --json`
- both are meaningful
- neither fact should be inflated into “the cluster truth layer is solved”

That is the kind of nuance the maturity matrix is supposed to preserve.

## Bottom line

This page is the operator’s lie detector.

If a future doc change implies the architecture is cleaner, safer, more
automatic, or more failover-ready than this matrix supports, the burden is on
that edit to show stronger evidence.

If the evidence is not there, the maturity level should not move.

And if a future edit makes the architecture sound cleaner by deleting mixed
truth, remembered exceptions, or unresolved failure semantics, that edit has
made the page worse even if it made it easier to read.
