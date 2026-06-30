# Request Path and Failure Walkthrough

This page exists because “HA routing” becomes vague almost immediately unless
someone forces the question back down to a literal request.

The user is not asking for a category label. The user is asking:

> when one hostname is requested, what literally happens next, who decides each
> step, what hidden assumptions are still doing work, and where exactly does
> the current proof stop?

That needs to be answered like an operator tracing a real request through a
real stack, not like a README polishing an architecture idea until the pain
disappears.

There is a second implied question under that one:

> if the request succeeds, did the system preserve the request itself, or did
> a human silently preserve the architecture around it?

That question is what separates a satisfying walkthrough from a decorative
one.

It also separates a real answer from the kind of answer the user keeps
rejecting.

The user is not asking for a prettier routing explanation.
They are asking whether the system can stop relying on hidden operator
reconstruction at the exact moment the wrong-node event becomes real.

## What this page is and is not allowed to prove

This page is allowed to:

- reconstruct the intended end-to-end request contract from repo evidence
- show where that contract currently has believable support
- expose exactly where the request story starts depending on missing machinery
- prevent architecture language from silently pretending a request path is
  already durable

This page is not allowed to:

- claim stack-wide wrong-node success
- treat a plausible walkthrough as runtime proof
- collapse HTTP route plausibility into TCP or stateful failover claims
- imply that middleware, auth, locality, peer eligibility, and application
  correctness have all already survived real failure conditions

## Quick claim router

If the question is:

- "What is the user actually trying to make a request do?" this page is one of
  the best reconstruction surfaces.
- "Where does the request story become aspirational?" this page is supposed to
  point at that seam explicitly.
- "Does the current stack already prove generic wrong-node routing?" no.
- "Can I use this page alone to claim end-to-end failover?" no. Use it with
  the proof matrix and failure-model pages.

## What a solved request path would actually feel like

The user is not trying to win an architecture argument.

They are trying to get to a point where the request path stops creating doubt.

In the solved version of this repo, an operator should be able to assume:

- the request can land on any healthy public node
- the receiving node can tell whether the service is local without relying on
  private memory
- if the service is remote, the receiving node can identify the right healthy
  peer from current truth rather than stale folklore
- the fallback path still exists when the local backend is exactly the thing
  that failed
- auth, middleware, and visible policy still mean the same thing after handoff
- the operator can explain why the request succeeded by reading tracked truth,
  not by reconstructing the architecture socially

If a path is reachable but one of those assumptions is still false, the repo is
closer to a clever routing demo than to a trustworthy request-preserving
platform.

## The sentence this page is protecting

The dream is not:

> several nodes are online and a proxy exists.

The dream is:

> a request for one service can land on any healthy node, stay local when that
> is honest, move to a healthy peer when locality is absent or broken, and
> still preserve the same service contract without pretending HTTP, TCP, and
> stateful systems are one problem.

Everything below is just a stress test of that sentence.

## Strongest honest current answer

The repo can already describe a serious request contract rather than a vague
"maybe we can proxy to peers someday" hope, but that is still weaker than real
route-specific proof. The honest current answer is that the architecture dream
is sharp, the likely direction is visible, and some proxy ingredients exist,
but the decisive test remains the uncomfortable one: a request lands on the
wrong node, the system chooses an eligible peer, policy survives the hop, and
the application still behaves correctly under that exact path.

That means this page is not merely documenting one request path.
It is helping decide whether the repo is actually building a new class of
operator experience, or only a more sophisticated version of the same old
"hope the right node gets hit" bargain.

This page should therefore be read less like a tutorial and more like a
cross-examination.

Every step should answer:

- who knows this?
- how do they know it?
- what disappears if that knowledge disappears?
- what in the current worktree actually proves the answer?

That repetition is intentional.

The user is not missing generic high-availability vocabulary.
The user is exhausted by systems that can explain themselves fluently while
still forcing the operator to privately reconstruct where the service really
lives, which node is still sacred, and which fallback story only works while
nothing important is actually broken.

That means this page should be read as scenario-based reconstruction, not as a
generic walkthrough.

The standard is closer to:

- take a concrete request
- take a concrete failure
- ask what must be true at each step
- separate what the worktree proves from what the plans merely want

If that sounds harsher than normal docs, good.
Normal docs are exactly what kept turning this repo's central question into a
cleaner but smaller neighboring one.

So this page keeps returning to the same pressure on purpose:

- wrong-node behavior
- backend-loss behavior
- policy-preserving behavior
- truth that survives operator absence

That last item matters more than it sounds.
One of the repo's hidden enemies is not just a dead node.
It is a living architecture that still needs a specific operator to remember
where the truth really lives.

If the walkthrough feels unforgiving, it is because the repo is trying to stop
accepting adjacent answers as if they were the real one.

## The emotional failure this walkthrough is trying to prevent

The user is trying to get out of a very specific experience:

1. the request lands on the wrong node
2. the system looks distributed, but the real answer still lives in the
   operator’s head
3. now the operator has to remember where the service actually lives
4. now the operator has to trust that proxy glue, helper templates, and
   failover assumptions still line up
5. now “multi-node” mostly means there are more places for ambiguity to hide

That is why this page focuses on request preservation, not just reachability.

Reachability can be faked by many things:

- a still-live public node
- a static redirect
- a stale peer route
- a TCP forward that ignores policy meaning
- a standby process that answers differently from the expected service

Request preservation is stricter.
It asks whether the system preserved the thing the user thought they were
talking to.

## The request we are tracing

Use a stateless HTTP example first, because that is where the repo is most
likely to earn real proof before everything else:

```bash
curl -H "Host: dozzle.bolabaden.org" https://bolabaden.org
```

Or conceptually:

```text
client -> https://dozzle.bolabaden.org
```

The hostname is not the point. The class is:

- public HTTP
- edge-routed
- potentially local on one node and remote on another

That is the clearest place to explain the multi-node dream without cheating.

## The components that can already participate in this path

The current root runtime already shows that a real request can involve:

- Cloudflare-backed public entry assumptions
- `cloudflare-ddns`
- `traefik`
- `crowdsec`
- `crowdsec-init`
- `tinyauth`
- `nginx-traefik-extensions`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `docker-gen-failover`
- the service container itself

That matters because the stack is already complicated enough that “the proxy is
up” is almost meaningless.

The useful question is:

which of these components merely exist, and which of them actually prove the
request survives the failure under discussion?

That distinction is also why this repo cannot be judged by component count or
surface sophistication.
One of the user's deepest complaints is that the ecosystem keeps offering more
components much faster than it offers more truthful options.

That is the whole difference between:

- a stack that looks comprehensive in screenshots
- and a system that has actually started owning the missing middle layer

This repo already has enough serious moving parts that "look how much infra is
here" stopped being a meaningful answer a long time ago.

That is why the walkthrough has to behave more like a set of held-out test
cases than like a polished architecture tour.

If a route only sounds convincing in the easiest local-success story, it has
not yet earned the right to be described as part of the resilient platform
dream.

## Separate the runtime classes now or the walkthrough starts lying

### Class 1: stateless or near-stateless HTTP

Examples in the runtime include:

- `dozzle`
- `docs`
- `grafana`
- `prometheus`
- `wishlist`
- `homepage`

These are the best places to prove wrong-node recovery first.

### Class 2: raw TCP but not automatically harmless

Representative examples include:

- `redis`
- `mongodb`
- some `biodecompwarehouse*` surfaces

Here, transport and correctness already start to diverge.

### Class 3: state-bearing systems with real ownership semantics

Examples include:

- Redis when used as real state, not throwaway cache glue
- MongoDB with primary/election expectations
- RabbitMQ with queue durability and topology
- Postgres and similar critical stores

For this class, the request path stops being only a routing story.

If this split disappears, the rest of the page becomes dishonest.

## The scenario family that actually decides whether the docs are honest

This page should be read against four concrete scenario classes:

### Scenario A: local success

The request lands on a healthy node and the service is local.

This is the easiest case.
It is useful, but it is also the least revealing case.

### Scenario B: wrong node, healthy peer

The request lands on a healthy public node that does not host the service
locally, but some peer does.

This is the center of gravity of the dream.
If this case still depends on folklore, the platform is still weaker than the
docs should ever imply.

### Scenario C: wrong node, local backend disappearance

The request lands on a node whose local assumption is broken or whose local
backend just disappeared, and the fallback path is needed precisely because
the local truth failed.

This is one of the most important scenarios in the entire repo.
It is the place where many elegant stories collapse, because the same
mechanism that made the local route look real often turns out not to preserve
the recovery route once the local backend is gone.

This is the case that exposes whether the route generator and failover logic
actually survive the failure they are named after.

### Scenario D: stateful false confidence

The request can still be forwarded somewhere, but forwarding alone no longer
proves semantic correctness because the service class has ownership,
replication, election, or durability meaning.

This is the case that prevents the docs from silently upgrading liveness into
correctness.

There is also one hidden scenario running underneath all four:

### Scenario E: operator reconstruction disguised as runtime behavior

The request appears explainable after the fact, but only because the operator
already knew:

- which node really hosted the service
- which helper state was current
- which peer was actually eligible
- which route was only safe in the easiest version of failure

If that is what happened, then the walkthrough is not documenting system
behavior.
It is documenting a successful human reconstruction of missing system truth.

Those four scenarios are not a formal test suite.
They are the minimum mental harness the docs need if they want to stop
pretending broad coherence where only partial proof exists.

## Walkthrough 1: the easiest honest case

This is the case the repo most wants to make boring:

1. the client reaches a healthy public node
2. the requested service is already local
3. the local edge stack is healthy enough to route it
4. the service responds

This case matters, but it is not the benchmark that decides whether the dream
has been achieved.

It mainly proves the system is not broken in the most flattering direction of
travel.

### Step 1: the client reaches some public node

At the very front, the request reaches a public node the architecture wants to
be eligible for entry.

The repo strongly supports this as intent:

- more than one public node should be able to receive the first hop
- Cloudflare-backed entry should reduce sacred-ingress-machine pressure
- no single reverse-proxy box should silently become the real system

What is strongly supported:

- first-hop plurality is explicit
- any-node entry is explicit
- Cloudflare is part of the story

What is not proved by this step:

- that the service will actually succeed

This is the first place where weaker infra writing usually cheats.

### Step 2: the request enters the local edge stack

Once the request lands, the local edge stack has to process it coherently.

For HTTP paths, that can include:

- Traefik router selection
- auth decisions through `tinyauth`
- CrowdSec-related behavior
- nginx extension handling
- local or dynamic backend selection

What the worktree proves strongly:

- the edge stack is real
- auth and middleware are part of the live path
- ingress is not a toy one-container proxy story

What it does not yet prove:

- that all of the same semantics survive peer fallback

### Step 3: the node concludes the service is local

This is the clean fast path.

The receiving node effectively decides:

> I host the requested service here, so I should serve it here.

That is the desired first choice because it preserves:

- locality
- fewer hops
- lower latency
- clearer debugging
- a more readable operator model

What the repo clearly wants:

- local-first whenever possible

What the current docs still cannot upgrade to settled runtime truth:

- one shared, auditable placement-truth mechanism explaining exactly how that
  locality decision is known today

That missing explanation is one reason `services.yaml` keeps coming back.

### Step 4: the local backend responds

If the request completes here, the system feels calm.

That is useful proof.

But it is still only:

- happy-path ingress proof
- not wrong-node proof
- not backend-loss proof

If the docs blur those together, they overstate maturity.

That blur is one of the main reasons infrastructure writing becomes
emotionally useless in a repo like this.

If a local happy path is described with the same confidence as a wrong-node or
backend-loss path, the easiest evidence trains the operator to believe the
hardest claim.

That training effect is one of the main documentation failures this
knowledgebase is trying to undo.
In a repo with this much real machinery, the local success story is dangerous
precisely because it is true.
It creates just enough confidence to make the unresolved cases harder to see.

## Walkthrough 2: the real dream case

This is the case the user actually cares about most:

1. the client reaches a healthy public node
2. the target service is not local
3. the receiving node still succeeds by handing the request to the correct
   healthy peer
4. the request still behaves like the same service

This is where the repo either becomes a real request-preserving platform or
proves it still is not one yet.

### Step 1: the request lands on the wrong node

The receiving node is healthy enough to accept traffic, but it does not host
the target service locally.

In this architecture, that is not an edge case. It is one of the main desired
operating behaviors.

If a doc page only explains the local happy path, it has already missed the
actual question.

### Step 2: the node must know the service is not local

This is the first hard control question:

> how does the node know the request cannot be satisfied locally?

Weak answers include:

- the operator remembers
- the hostname merely hints at it
- local labels imply something
- the docs assume everyone knows where the service lives

Those are not good enough, because they keep sacred-node knowledge alive under
new branding.

If the receiving node only behaves correctly because humans around it already
know the placement truth, then the repo has not removed the SPOF. It has moved
the SPOF into shared operator memory.

This is one of the most important sentences in the whole knowledgebase.

The user is not only trying to avoid one dead server.
They are trying to avoid one remembered server.

That is why the wrong-node request is such a useful benchmark.

It forces the system to answer whether the architecture truth is:

- in tracked runtime knowledge
- in generated helper state
- in a durable registry
- or still mostly in the operator's private reconstruction loop

This is why `services.yaml` matters so much in the archive. The file itself is
not magic. The pressure behind it is the real thing: the repo keeps
rediscovering that shared placement truth has to exist somewhere explicit.

The current worktree strongly documents the need for that layer. It does not
yet strongly prove one settled live implementation of it at the root runtime.

### Step 3: the node must know which peer currently hosts the service

Not:

- which peer hosted it yesterday
- which peer the docs hope hosts it
- which peer a stale template still names

It needs a current answer.

This is the actual “service discovery” problem that remains after the user has
already dismissed manual placement and DNS plurality as the main issue.

The archive states this plainly: if scheduling is manual and Cloudflare already
handles first-hop plurality, the hard remaining problem is routing truth.

### Step 4: the node must know the peer is eligible now

Even if the peer is the intended host, that is still weaker than:

- healthy
- semantically compatible
- carrying coherent config and secrets
- safe to receive the request right now
- preserving the same edge assumptions

Without convergence truth, wrong-node forwarding becomes:

- mechanically plausible
- operationally seductive
- but not yet trustworthy

A node being reachable is cheap.
A node being eligible for this exact request, with the right revision, secrets,
middleware assumptions, and backend state, is the real platform question.

That sentence needs to stay strict because this is where "more options" often
collapse into disguised ambiguity.
If there are many possible peers but no strong shared truth about which peer is
honest now, then the system has more branches without more certainty.

This is where many stacks quietly replace truth with hope:

- the named peer is probably up
- the route is probably still there
- auth is probably equivalent there
- the middleware path is probably close enough

The user is specifically asking this repo to stop building on "probably."

That request is deeper than "please be more robust."

It is closer to:

> stop giving me options that only sound like system behavior while still
> depending on tacit operator memory and lucky convergence underneath

### Step 5: the handoff route must still exist

Once the node decides to hand the request off, the route needed for that
handoff must still exist under the exact failure that made fallback necessary.

This is why `docker-gen-failover` remains a central honesty boundary.

Planning evidence already records the dangerous failure mode: route generation
can remove the route when containers stop. If that is still true, the docs
cannot honestly say:

> fallback routing is trustworthy under backend loss

until route persistence is fixed and shown.

If the route disappears when the backend disappears, the system did not fail
over.
It revealed that the failover story was attached to the thing that died.

### Step 6: the request must still preserve policy

A peer-forwarded request is not an honest success if:

- local path requires auth but peer path bypasses it
- middleware ordering changes
- headers or forwarding assumptions drift
- the same URL now means a weaker service contract

In this repo, the question is not “did the peer return `200`?”

The real question is:

> did the same service contract survive the handoff, or did the system merely
> return something useful-looking?

A superficially successful peer path can still be wrong if:

- the forward-auth policy drifted
- the middleware chain changed meaning
- the backend answered with a different trust boundary
- the visible hostname stayed stable while the semantic contract weakened

This is why "HTTP 200" is too weak as evidence.
It proves response.
It does not yet prove preserved service identity.

### Step 7: the operator must be able to explain the success without folklore

Even after the request succeeds, one more question remains:

> can a tired operator explain why this worked by reading tracked truth, or do
> they still have to remember who the real host was and which helpers were
> lying less than usual?

This step often gets omitted because it sounds post hoc.
It is not post hoc here.
It is part of the architecture contract.

If successful routing still depends on unwritten memory to be intelligible,
then the human SPOF has not been removed.

> did the peer return the same effective service contract?

That includes the meaning of auth, middleware, headers, and hostname identity.

This deserves to stay blunt:

- a response code can survive while the real policy contract dies
- a proxy handoff can preserve transport while weakening meaning
- a user-visible success can still be an architectural failure if the wrong-node
  path behaves like a different product

That is why the user keeps sounding dissatisfied with “it still returned 200”
style evidence.

They are not chasing green lights.
They are chasing preserved meaning.

That is also why ordinary synthetic demos feel insulting here.

A demo that proves a peer can answer something is cheap.
A demo that proves the same request still means the same thing after locality
breaks is much closer to the actual dream.

### Current maturity of the wrong-node dream

Today this whole walkthrough is:

- explicit as architecture
- structurally plausible from the edge shape
- not yet proved end to end

That is the honest state.

Another way to say it:

- the repo has already escaped the "toy reverse proxy" phase
- it has not yet escaped the "operator still secretly stitches the truth
  together" phase

## Walkthrough 3: local backend loss after the contract already exists

This is the walkthrough class that most punishes vague wording.

The repo can sound very advanced while still failing here, because the stack
already contains helpers, generators, multiple public nodes, and dynamic edge
components.

That only makes precision more necessary.
If the docs say "failover" when the route vanishes with the primary backend,
the docs are actively making the platform harder to reason about.

This is the path that separates clever routing from actual failover:

1. the request would normally prefer the local backend
2. the local backend disappears or becomes unhealthy
3. the route needed for recovery must remain
4. the peer path must take over
5. policy must remain coherent

### Step 1: local backend is preferred

The planning layer already points toward the intended preference model:

- local backend is the fast path
- peer backends are fallback or lower weight

That matches the repo’s philosophy:

- preserve locality when possible
- leave the node only when necessary

### Step 2: local backend becomes unhealthy or disappears

The live stack already contains real healthchecks and health-aware components.

That proves seriousness.

What it does not prove yet is the most important thing:

> when the local backend becomes unhealthy, does the recovery route survive?

### Step 3: route persistence becomes the whole problem

This is where the docs have to stop smoothing the story.

If the failover generator removes the route when the local container stops, the
system does not merely need refinement. It fails the central promise.

That is why `docker-gen-failover` is still treated as warning evidence, not as
finished proof.

This is also where the user’s frustration with “options” becomes easiest to
understand.

Lots of approaches can generate a fallback route.
Far fewer can guarantee that the route survives the exact local backend loss
that made it necessary.

That is the real test of whether an option was ever an option.

If the route only exists while the preferred local backend is healthy, then the
fallback story was not a second choice.
It was just a flattering description of the first one.

### Step 4: the peer path would need to take over

The planning layer clearly wants:

- a service registry that does not delete truth when containers stop
- health-aware peer selection
- local preference with healthy-peer fallback

That is strong planning evidence.

It is not yet live failure proof.

### Step 5: the request still has to mean the same thing

If a backend-loss path returns a different auth experience, a weaker
middleware chain, or a semantically degraded service, then the failover story
is not complete.

This repo is correct to treat that as architectural failure, not minor drift.

## Walkthrough 4: why TCP and stateful systems do not inherit the HTTP story

This is where many infrastructure docs become actively misleading.

### Raw TCP reachability is weaker than service correctness

A TCP router proving that bytes can move between client and backend does not
prove:

- state ownership
- election semantics
- replication correctness
- client rediscovery
- safe write continuity

### Redis example

If a client connects to `redis.$DOMAIN`, a successful TCP path still does not
answer:

- who is primary
- whether replicas exist
- how failover is triggered
- how clients rediscover authority
- whether split datasets are avoided

So when the repo talks about Redis, it has to stop pretending HTTP ingress
logic is the whole story.

### MongoDB example

A routed MongoDB connection is weaker than:

- a working Replica Set
- election behavior
- correct client discovery semantics

The endpoint being reachable is not the same thing as replica-set truth being
preserved.

### Operator rule

Read TCP and stateful claims with a stricter filter:

- HTTP can often be proved path by path
- stateful systems must be proved topology by topology

That is not pessimism. It is the only honest standard.

## What the current worktree proves strongly in this walkthrough

The current worktree strongly proves:

- the repo has a substantial real edge stack
- the any-node-entry dream is explicit
- local-first plus peer-forward is explicit
- HTTP and TCP are already first-class in the ingress story
- auth and middleware are genuinely in the live request path
- route persistence under failure is recognized as a real problem
- `services.yaml` is recognized as a necessary missing truth surface

It also strongly proves that the repo already understands the question better
than many generic HA guides do.

What it does not yet prove is that understanding has fully crossed over into a
runtime that can answer the question for itself.

The problem is not absence of architectural pressure.
The problem is that several key truths are still pressure rather than settled
runtime proof.

## What the current worktree still does not prove strongly enough

The current worktree does not yet strongly prove:

- a generic wrong-node HTTP success path
- a backend-loss path where route persistence clearly survives
- a live shared placement-truth surface consumed by routing
- a fully proved semantic-continuity path across peer handoff
- a stateful zero-SPOF story for Redis, MongoDB, RabbitMQ, or similar data
  systems

It also does not yet strongly prove that an operator can explain the whole
wrong-node path from one tracked truth source instead of synthesizing it from
edge components, Compose fragments, helper intent, and memory.

That list is not negativity. It is what keeps the docs from flattering the
stack into lying.

## The shortest operator takeaway

Today the repo should be read like this:

- first-hop plurality is part of the architecture
- local-first service is part of the architecture
- wrong-node recovery is part of the architecture
- the edge machinery is real
- the truth and failure-proof layers behind that dream are still maturing

The practical reading is:

the repo already has enough real edge infrastructure that the dream is no
longer theoretical, but it still lacks enough shared truth and failure proof to
pretend the wrong-node path is generically solved.

That unresolved wrong-node path is not one bug among many.
It is the clearest place where the repo can still accidentally collapse back
into the same fake-option landscape that made the project necessary.

So the system is not “just a bunch of Compose files.”

But it is also not yet entitled to narrate itself as a fully proved
request-preserving multi-node platform.
