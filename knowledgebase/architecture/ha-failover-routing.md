# HA, Failover, and Routing

For the strongest evidence behind this page, read
[`../research/ingress-and-failover-evidence.md`](../research/ingress-and-failover-evidence.md)
first. Then read
[`request-path-and-failure-walkthrough.md`](request-path-and-failure-walkthrough.md)
for the literal request trace.

This page is not here to celebrate that the repo has Traefik, Cloudflare, more
than one node, and several services. That would be the easy version. The user
is already tired of the easy version.

The hard version is this:

> if a request lands on a healthy node that does not host the target service,
> can that node still preserve the meaning of the request instead of quietly
> becoming a dead end, a redirect machine, or a new place where sacred-node
> knowledge hides?

That is the routing question in `bolabaden-infra`.
That question is stronger than "is there failover?"
It asks whether the wrong node stops being semantically wrong.

## What this page is and is not allowed to prove

This page is authoritative about:

- which routing truth layers have to coexist before "failover" means anything
- why first-hop plurality is weaker than request-preserving recovery
- where the repo's routing story still depends on missing placement,
  eligibility, persistence, or semantic truth

This page is not authoritative about:

- whether a specific wrong-node route has already been proven end to end
- whether backend-loss recovery already survives for a given service
- whether stateful traffic has inherited HTTP honesty for free

This page explains the routing problem shape and its honesty boundaries.
It should not be used as a substitute for route-specific drills.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "DNS or any-node entry helps only with the first hop" | routing-layer analysis | this page, ingress evidence, request-path pages | that first-hop plurality is unimportant |
| "routing still lacks shared truth" | live-gap synthesis | this page, current runtime, placement discussions, evidence pages | that no routing progress exists at all |
| "one helper could still fail the real requirement" | failure analysis | this page plus ingress and failover evidence | that the whole direction is worthless |
| "request-preserving failover is stricter than route reachability" | proof-discipline synthesis | this page, proof matrix, runbook | that route reachability has no value whatsoever |

If a routing sentence starts sounding like a passed failover drill, it has
left this page's authority.

## What the user is actually rebelling against

The repo archive keeps circling the same frustration from different angles:

- multiple A records help, but they do not answer where the service actually is
- a reverse proxy helps, but it does not invent cross-node truth
- matching labels look clean, but local Docker labels are not global discovery
- a failover helper sounds promising, but if the route disappears when the
  backend disappears, the helper is part of the problem
- “just use Kubernetes” often arrives before anyone has precisely named the
  missing capability

The user is not merely asking for prettier HA language.

The user wants ordinary Docker hosts to stop behaving like isolated islands
that only appear unified as long as the operator remembers which machine is the
real one.

That last clause is the part many surrounding ecosystems keep failing to hear.
The user is not merely annoyed by node-locality.
The user is annoyed by architectures that still outsource truth to remembered
locality while offering polished new words like mesh, failover, cluster, or
service discovery.

That emotional center matters because it explains why this repo keeps rejecting
solutions that are technically adjacent but psychologically wrong. “More load
balancing” is not enough. “More boxes” is not enough. “A backup proxy in front
of the proxy” is not enough. The repo is trying to get to a place where
traffic can land anywhere and the system itself, not the operator’s private
memory, knows what to do next.

## The routing claim that would actually matter

The dream is not:

- DNS resolves to more than one machine
- Traefik can route local containers
- a node can reach another node over the network
- a helper can generate some dynamic config

The dream is:

> a request can hit any surviving public node, stay local when that is honest,
> and move to a healthy peer when locality is absent or broken, while
> preserving auth, middleware, routing policy, and operator readability

Anything weaker than that may still be useful engineering. It is just not the
thing the user keeps asking for.

That distinction has to stay explicit because this repo already contains enough
real edge machinery to make an adjacent answer sound like the real answer.
The docs therefore need to keep saying, in different ways, that "can route"
and "preserves the request contract under the bad failure" are not neighboring
confidence levels.

## The hidden difference between “multi-node” and “request-preserving”

Many infrastructures stop too early and still sound impressive.

They prove:

1. more than one machine exists
2. more than one machine can receive traffic
3. the proxy layer can see local containers

Then they start narrating resilience.

But the user’s standard is stricter:

1. the request lands on the wrong healthy machine
2. that machine has to know the service is not local
3. that machine has to know which peer currently owns the service
4. that machine has to know the peer is eligible now, not merely named in some
   stale config
5. the route needed for fallback has to survive the failure that made fallback
   necessary
6. the request has to keep meaning the same thing after handoff

That is why this repo keeps splitting the story into:

- node-entry survival
- locality truth
- peer-selection truth
- route persistence
- semantic continuity
- stateful honesty

If those collapse into one sentence called “HA,” the docs become decorative.

## Routing here has to stay layered or it becomes dishonest

This stack already contains enough moving parts to fool casual readers.

The live request path can involve:

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

That is a real edge surface. It is also exactly the kind of surface that can
sound “HA-like” while still missing the one truth layer the user actually
needs.
It is also the kind of surface where a convincing demo can make the docs
weaker:

- the first hop works
- a route answers
- the logs look serious
- and everyone forgets to ask whether the same route survives the wrong failure
  with the same meaning

The layers have to be kept separate.

### Layer 1: public node-entry reachability

Question:

- can the client reach some healthy public node at all?

This is where Cloudflare and multi-record DNS help.

What this layer can honestly buy:

- more than one node can receive the first hop
- one public node is less likely to become a sacred ingress machine
- some node loss can be survived at the DNS or edge-entry level

What it does not buy:

- service-locality truth
- peer eligibility truth
- fallback-route persistence
- stateful correctness

This matters because the archive repeatedly drifts toward the same false jump:

1. all nodes are now public
2. therefore failover is mostly solved

That jump is wrong.

### Layer 2: local edge-stack health

Question:

- once the request reaches a node, is the local edge stack healthy enough to
  make the next decision coherently?

This includes more than “Traefik is running.” The local stack may need to
preserve:

- forward-auth behavior
- middleware ordering
- security filtering
- rewrite and redirect logic
- Docker-provider visibility
- dynamic file-provider state

If this layer is broken, distributed entry above it is meaningless.

### Layer 3: locality truth

Question:

- does the receiving node actually know whether the requested service is local?

This is where the archive keeps converging on `services.yaml` or an equivalent
placement-truth surface.

The reason is simple: local labels and human memory are not enough. The repo is
not searching for `services.yaml` because files are fun. It is searching for a
surface that stops “the operator remembers where it lives” from being the real
control plane.

That is why the placement-truth conversation cannot be reduced to
"service discovery" as a fashionable category.
The repo is not looking for discovery in the abstract.
It is looking for the death of sacred remembered placement.

### Layer 4: peer-selection truth

Question:

- if the service is not local, does the receiving node know which peer
  currently hosts it and whether that peer should be trusted right now?

Not:

- the peer the docs mention
- the peer that hosted it last week
- the peer a stale template still points at

It needs current truth.

This is the actual “service discovery” problem that the archive keeps naming
after DNS and placement have already been stripped away.

### Layer 5: route persistence

Question:

- when the preferred backend disappears, does the route needed for recovery
  remain available?

This is one of the hardest honesty walls in the repo.

The stack already contains `docker-gen-failover`, and the direction is clearly
toward dynamic failover behavior. But planning evidence also records a critical
failure mode: routes can disappear when containers stop. If that remains true,
then the mechanism does not merely need polish. It fails the user’s actual
requirement at the exact moment the system should prove itself.

This is one of the sharpest places where the docs can accidentally become
harmful.
If a page describes dynamic failover as an achieved capability while this
failure mode is still live, the page is teaching the operator the wrong
survival model.

### Layer 6: semantic continuity

Question:

- after peer handoff, does the request still behave like the same service?

That means preserving:

- auth expectations
- middleware ordering
- headers and forwarding assumptions
- externally visible hostname meaning
- redirect and rewrite behavior
- operator understanding of what happened

In this repo, semantic continuity is part of routing correctness. It is not
optional polish.

It is also where "better than nothing" answers become dangerous.
If the peer path is only "close enough" while auth, headers, redirects, or
trust boundaries drift, then the stack did not preserve the request.
It preserved the outer shell of a request while changing what the service
actually meant.

## Why `docker-gen-failover` is still an honesty wall

The repo is right to keep this component under suspicion instead of treating it
as a solved story.

Its existence proves something useful:

- the repo is not naively content with local-only routing
- dynamic failover is part of the active design pressure
- the edge stack is trying to evolve beyond static local labels

Its existence does not prove:

- that the fallback path survives backend loss
- that peer forwarding is trustworthy under real failure
- that auth and middleware stay coherent after takeover

The proof bar here is strict:

1. a request succeeds while the service is local
2. the local backend becomes unavailable
3. the route needed for recovery remains present
4. a healthy peer is selected
5. the same policy path is preserved
6. the client still experiences the same service contract

Without that chain, the honest label is still partial.

Another way to say it:

- a surviving path is not yet a preserved contract
- a returned payload is not yet the same service
- a reachable peer is not yet an eligible peer
- a named failover helper is not yet failover

## Global names and node-scoped names mean different things

This repo already encodes two different routing expectations, and the docs
should preserve that instead of flattening them.

### Global names

Examples:

- `docs.$DOMAIN`
- `grafana.$DOMAIN`
- `service.$DOMAIN`

Meaning:

- any honest healthy path that can fulfill the request should be acceptable

These names create interchangeability pressure.

### Node-scoped names

Examples:

- `docs.$TS_HOSTNAME.$DOMAIN`
- `service.$TS_HOSTNAME.$DOMAIN`
- `mongodb.$TS_HOSTNAME.$DOMAIN`

Meaning:

- locality remains part of the explicit contract

These names are not cosmetic aliases. They are one of the repo’s clearest ways
to keep “global failover identity” separate from “I want this specific node.”

If the docs blur those together, they erase one of the most practical mental
tools in the stack.

## HTTP and TCP do not deserve the same confidence

The repo already routes both classes, but they do not inherit the same proof.

### HTTP is where real early proof can be earned

HTTP is the most promising class for honest wrong-node proof because the repo
already has:

- mature hostname routing
- explicit middleware composition
- forward-auth surfaces
- health-aware load-balancing concepts
- many workloads that are comparatively easier to relocate

If the repo first proves real request-preserving multi-node behavior anywhere,
it will probably be here.

### TCP is a sharper boundary

The root runtime already exposes TCP-oriented services through Traefik TCP
routers and passthrough. That proves seriousness. It does not prove safety.

For TCP services, especially state-bearing ones:

- reachability is not resilience
- passthrough is not topology truth
- a socket accepting traffic is not stateful continuity

This is why the docs keep separate phrases for:

- TCP exposed
- TCP routable
- TCP failover-aware
- statefully resilient

Those are not the same maturity step.

## The maturity ladder that actually matches the user’s dream

The routing story makes more sense when treated as a ladder instead of a label.

### Level 1: a surviving public node is reachable

Proves:

- the client can hit some node

Does not prove:

- service success

### Level 2: a local service can be served locally

Proves:

- local-first routing works in a happy path

Does not prove:

- wrong-node success

### Level 3: the wrong node knows where the service should be

Proves:

- placement-truth logic exists at least partially

Does not prove:

- the route survives failure

### Level 4: the wrong node can forward mechanically

Proves:

- some peer handoff works

Does not prove:

- semantic continuity
- backend-loss resilience
- stateful safety

This level is where a lot of "it basically works" language comes from.
The repo needs to keep treating that language as suspect.
Mechanical forwarding is one of the most seductive half-solutions in the whole
problem space, because it makes the architecture look unified while the real
truth model is still fragmented.

### Level 5: the wrong node preserves the request semantically

Proves:

- auth, middleware, and request meaning survive handoff

This is the first level that starts to resemble the actual user dream.

### Level 6: local backend loss still preserves the request

Proves:

- the system survives the precise failure that made fallback necessary

This is the first level that deserves the word failover without apology.

## What the operator should conclude today

The honest current reading is not defeatist. It is precise:

- the ingress surface is real
- the edge stack is complex and serious
- the any-node-entry philosophy is explicit
- local-first plus peer-forward is explicit
- the stack already spans HTTP, TCP, auth, and middleware
- the pressure toward a placement-truth layer is explicit
- route persistence under failure is still a critical unresolved boundary
- stateful services still require a stricter proof model than the HTTP story
- hidden operator reconstruction is still part of the risk surface whenever the
  truth layer is implied rather than inspectable

So the right current answer is:

- this is not “just some Compose files”
- this is not yet a fully proved request-preserving multi-node platform either

That distinction is the whole point of this page.
If that sentence ever gets rounded into "mostly there," this page has failed.
