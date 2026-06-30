# Request Path and Failure Walkthrough

This page reconstructs the request path the user is actually asking about.

Not:

> what proxy do we use?

and not:

> what components exist at the edge?

The real question is:

> when a hostname is requested, what literally happens next, which runtime
> surfaces own each step, what truth does the receiving node have on its own,
> and where does the story stop being system-owned and turn back into private
> operator memory?

That is the only honest way to talk about "multi-node without Swarm" in this
repo.

## What this page is and is not allowed to prove

This page is allowed to prove:

- the target request contract the repo keeps trying to earn
- which live runtime components already participate in that contract
- which parts of the request path are materially present in the current Compose
  runtime
- where proof stops today for wrong-node, peer-forward, and backend-loss
  behavior
- why route class matters: plain HTTP, protected HTTP, TCP, and stateful lanes
  are not one maturity surface

This page is not allowed to prove:

- generic wrong-node success today
- generic peer-forward correctness today
- generic backend-loss survival today
- that a plausible walkthrough is the same thing as a validated drill
- that HTTP optimism transfers into TCP or stateful authority

This is a decomposition page.
It is not a green-check page.

## The fake walkthrough this page is supposed to stop

The bad version of this page would do something like this:

1. DNS resolves
2. Traefik matches
3. auth middleware exists
4. helper containers exist
5. therefore the distributed request story is mostly in place

That is exactly the walkthrough shape this page is meant to reject.

A detailed walkthrough can still be useless if it quietly skips the moment
where the receiving node stops knowing and the operator starts privately
completing the story.

## The target contract

The clearest contract still lives in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That contract is much stricter than:

- more than one public node exists
- Traefik is running
- some helper can generate fallback-looking config
- a route can probably be proxied somewhere

The contract is about preserved request meaning after locality fails.

That word "meaning" matters.
The user is not merely asking for packets to arrive.
They are asking for the wrong healthy node to stop acting stupid when it is not
the owner node.

## The one event this page is really about

This entire repo can be reduced to one embarrassing event:

1. a public request lands on a healthy node
2. the requested service is not actually local there
3. the node must decide what to do
4. the decision must be right for this route class
5. the explanation must exist outside one operator's head

If a page, helper, or architecture proposal does not improve that event, it is
probably orbiting the problem instead of solving it.

## The live request actors already present now

From the root runtime and active fragments, the present request path already
contains these real actors:

### Public-entry and naming actors

- `cloudflare-ddns`
- Cloudflare-first public entry assumptions

### Main L7 execution actors

- `traefik`
- Traefik router/service labels across root and fragments

### Policy and auth actors

- `nginx-traefik-extensions`
- `tinyauth`
- protected routes using `nginx-auth@file`
- `crowdsec`

### Fallback-intent and helper actors

- `docker-gen-failover`
- generated dynamic output to
  `${CONFIG_PATH:-./volumes}/traefik/dynamic/failover-fallbacks.yaml`

### Private reachability actors

- `headscale-server`
- `headscale`

### Proof-candidate actors

- `whoami`
- `wishlist`
- `mkdocs`

The runtime also already includes real routed service classes:

- simple stateless HTTP:
  - `whoami`
  - `wishlist`
  - `mkdocs`
  - `searxng`
  - `homepage`
- protected HTTP:
  - `code-server`
  - `dozzle`
  - `portainer`
  - metrics/admin routes such as `prometheus`, `cadvisor`, `alertmanager`
  - `mcpo`
- TCP/stateful transport:
  - `mongodb`
  - `redis`
  - multiple `biodecompwarehouse*` routes
- state-bearing app subgraphs:
  - `nuq-postgres`
  - `rabbitmq`
  - `litellm-postgres`
  - `qdrant`
  - Headscale SQLite authority

So the request-path problem is not hypothetical.
The missing layer matters because the runtime is already broad enough to make
the missing truth socially expensive.

## The private sentence test for this page

At the end of any walkthrough section, the reader should be able to answer:

> what exact sentence would still need private operator completion if this
> request landed on the wrong healthy node right now?

If the section leaves behind only components and flow language, but not that
surviving private sentence, the walkthrough is still too soft.

## The one question to keep asking at every hop

At every stage of the path, keep asking:

> what does the receiving node know from shared, inspectable, system-owned
> truth, and what would still have to be privately supplied by the operator if
> the local backend were absent?

If that question disappears, the walkthrough can stay technical while still
quietly turning back into theater.

## Walkthrough 1: normal-day stateless HTTP

Use `whoami.$DOMAIN`, `wishlist.$DOMAIN`, or `mkdocs.$DOMAIN` as the cleanest
mental model.

These are the most useful first-pass routes because they are:

- HTTP
- narrow
- low-state
- good proof candidates
- less likely to hide authority ambiguity

They are also the fairest early test of the repo's promise.
If the stack cannot make one low-state HTTP route survive wrong-node entry with
an inspectable explanation, then louder HA language elsewhere is almost
certainly premature.

### Hop 1: DNS and public entry

What is live now:

- Cloudflare DDNS is part of the active edge fragment
- the repo explicitly rejects one sacred public box as the only front door

What this honestly proves:

- first-hop plurality is part of the live design
- more than one node is intended to be a plausible public entry candidate

What it still does not prove:

- the receiving node knows where the service actually lives now
- the receiving node knows whether the target is local or remote
- the receiving node knows which peer is valid if the service is remote

This is the first place many ecosystems overclaim.
Plural DNS is necessary.
Plural DNS is not preserved request meaning.

Private sentence still surviving after Hop 1:

> traffic can hit more than one node, but I still personally know that this
> says nothing about what the chosen node should do next

### Hop 2: Traefik accepts the hostname

What is live now:

- Traefik is materially present
- root and fragment services carry real routers and service labels

What this honestly proves:

- the receiving node can parse the hostname
- the receiving node can match a local router definition
- the route exists in authored runtime truth

What it still does not prove:

- that the router owns cross-node placement truth
- that the receiving node can choose a remote backend from shared current
  truth

This hop is route execution.
It is not yet route truth.

Private sentence still surviving after Hop 2:

> the route matched, but I still personally know that route execution is
> weaker than distributed backend truth

### Hop 2.5: the missing decision point

This is the hop many infra writeups skip because it is socially awkward.

The missing question is:

> once the route is known, what tells this node whether it is the owner node or
> merely the first-hop node?

That is where placement truth and peer-eligibility truth would have to appear
if the contract were becoming real.

Without that layer, the walkthrough silently mutates into:

1. hostname matched
2. a human knows what should happen
3. therefore the route story feels understandable

That is the exact move this repo is trying to escape.

### Hop 3: middleware and auth policy are applied

What is live now:

- `nginx-traefik-extensions` participates in forward-auth behavior
- `tinyauth` provides real auth participation
- `crowdsec` is in the edge stack
- protected routes already use `nginx-auth@file`

What this honestly proves:

- the repo already treats request meaning as more than host-to-container
  forwarding
- auth, trust boundaries, and middleware continuity are live concerns

What it still does not prove:

- that the same semantics survive peer-forward handoff
- that the remote peer has converged enough config, secrets, and policy
  expectations to preserve the same route meaning

This is why protected HTTP has to stay stricter than a plain stateless route.

Private sentence still surviving after Hop 3:

> the same middleware names exist, but I still personally know that this does
> not prove the remote handoff preserves the same route meaning

### Hop 4: local service execution

If the target service is local to the receiving node, the story gets much
stronger.

What this honestly proves:

- the local happy path works for that route class
- the edge and app are at least locally connected
- the service is not merely theoretical

What it still does not prove:

- wrong-node dignity
- peer-eligibility truth
- backend-loss rescue

The local happy path is real progress.
It is not the user's actual benchmark.

Private sentence still surviving after Hop 4:

> the owner node can serve the route, but I still personally know that this
> says very little about wrong-node dignity

## Walkthrough 2: protected HTTP

Now use a route such as:

- `code-server.$DOMAIN`
- `dozzle.$DOMAIN`
- `portainer.$DOMAIN`
- `prometheus.$DOMAIN`
- `mcpo.$DOMAIN`

These are the revealing routes because the system has to preserve not just
transport, but meaning:

- auth behavior
- middleware order
- trust boundaries
- forward-auth behavior
- header assumptions
- visible route semantics

### What the live runtime already gives this class

The current runtime already gives protected routes:

- real Traefik routers
- `nginx-auth@file`
- TinyAuth participation
- middleware-bearing edge surfaces
- observability around edge and helper components

### What still remains unproven

The current worktree still does not prove:

- that a protected route forwarded to another node still means the same thing
- that the peer has converged enough secrets and config to preserve that
  meaning
- that the same protected route survives backend loss with the same visible
  contract

That is why "the page still loads" is too weak.
A route that answers but no longer behaves like the same protected service is
not a successful handoff.

## Walkthrough 2.5: the false comfort transition

This repo repeatedly risks one exact transition:

- a route exists locally
- the route is protected
- the edge stack is serious
- therefore a peer-forward version of the route must be close enough

That transition is illegitimate until the docs can show:

- the source of placement truth
- the source of peer-eligibility truth
- the continuity of middleware and auth semantics
- the continuity of the route under backend loss if failover is being claimed

Without those, the walkthrough is still leaning on emotional plausibility.

## Walkthrough 3: wrong-node entry

This is the scene the user actually cares about.

Assume all of the following:

1. Cloudflare sends traffic to a healthy public node.
2. Traefik, CrowdSec, and auth all look healthy.
3. The requested service is not actually local to that node.

At that moment, the repo becomes very good at describing the problem and still
not strong enough to prove the general solution.

### What the receiving node can likely know already

From the live edge stack, the receiving node can likely know:

- the hostname
- the route class
- the middleware chain it would normally apply
- whether the edge stack itself looks healthy

### What the current worktree does not yet prove the node knows on its own

The current worktree does not yet prove the receiving node has shared truth
for:

- exact live placement of the target service
- whether the relevant peer is healthy enough to receive the request
- whether that peer has converged secrets, environment, and auth semantics
- whether the forwarded request would still mean the same thing after handoff

That is the missing middle layer in one scene.

Private sentence still surviving in the wrong-node scene:

> the wrong node can see the request shape, but I still personally know facts
> the runtime has not yet proven it owns

## Walkthrough 4: backend-loss scene

There is a second trap that must stay separate from wrong-node entry:

> even if a wrong-node handoff looks plausible while everything is healthy, the
> rescue path may still disappear once the preferred backend actually fails

This matters because a helper can look excellent right up until the exact event
it is meant to absorb.

The repo already records this pressure around `docker-gen-failover`.

So backend-loss has to stay distinct from wrong-node success.

Another way to say the same thing:

- wrong-node success asks whether the first-hop mistake can be corrected
- backend-loss success asks whether the correction path survives the damage

The user is frustrated because many options solve one of those socially and not
the other operationally.

Wrong-node success means:

- the request landed on a non-owner node
- the system still chose a valid backend

Backend-loss success means:

- the preferred backend disappeared
- the rescue path still existed
- the route still meant the same thing
- the operator did not have to privately complete the story

Those proofs overlap.
They are not the same proof.

Private sentence still surviving in the backend-loss scene:

> the route looked recoverable while healthy, but I still personally know the
> rescue path may evaporate during the exact failure event

## Why helper presence still is not enough

The repo already contains several things that make the story sound closer to
solved:

- `docker-gen-failover`
- substantial Traefik dynamic and middleware logic
- Headscale-backed private networking assumptions
- route-rich Compose fragments

But helper presence alone still does not count as wrong-node proof.

## The actual deliverable this page is pushing toward

The real deliverable is not a smoother explanation of a hypothetical path.
It is one preserved route packet that can answer all of these:

- which node received the request first?
- why was it the wrong node?
- what shared truth named the real owner?
- why was the chosen peer valid now?
- what visible route semantics stayed the same?
- what artifact lets a second operator inspect that answer later?

Until one packet like that exists, this page is still primarily a warning page.

The user is specifically tired of:

- there is a failover helper
- there is a private mesh
- there is a route generator
- there is a dynamic proxy

being narrated as if the hard part had therefore moved out of the operator's
head.

It has not, unless the receiving node can show where it got the truth for the
handoff decision.

That is the burden-transfer test in its most practical form:

- not `is there a helper?`
- not `is there a generated file?`
- not `is there private connectivity?`

But:

> can the receiving node show the truth source behind the decision without a
> human privately narrating the rest?

## Why route class must stay separate

This repo becomes misleading the moment these route classes are narrated as one
shared maturity story.

### 1. Stateless HTTP

This is the lane most likely to earn real relief first.

Why:

- route meaning is visible at L7
- middleware and auth can be inspected
- one good drill can prove a narrow but real win

Best early proof candidates:

- `whoami`
- `wishlist`
- `mkdocs`

### 2. Protected HTTP

This lane is stricter because the route has policy meaning, not just transport
meaning.

Better future candidates:

- `code-server`
- `dozzle`
- `portainer`
- `prometheus`
- `mcpo`

### 3. Raw TCP

This lane already exists through live TCP routers for:

- `mongodb`
- `redis`
- `biodecompwarehouse*`

It proves real transport exposure.
It does not prove:

- peer-aware failover
- semantic continuity
- safe substitution
- stateful correctness

### 4. Stateful surfaces

This is the strictest lane of all.

A stateful service does not become honestly HA because:

- it is reachable through Traefik
- more than one node exists
- a proxy can point to more than one backend

Stateful routing becomes meaningful only when authority, durability,
promotion, reconnect, and rediscovery are explicit.

## The exact point where the story turns back into operator memory

The request story turns back into private operator memory the moment the
runtime cannot answer all of these on its own:

- where does this service really live right now?
- is the candidate peer merely alive, or actually valid for this route?
- if forwarded, will the route still mean the same thing?
- if the preferred backend disappears, what path still exists?
- if the route is stateful, who is the actual authority?

That is the real SPOF reading in this repo.

The user is not only worried about a dead server.
The user is worried about the operator still being the hidden explainer when
the system was supposed to behave like one cloud.

## What a real walkthrough-backed proof packet would need

Before this page can support stronger language, it needs a proof packet with:

- the exact hostname
- the receiving node identity
- the actual backend identity
- the source of placement truth used for the decision
- the source of peer-eligibility truth used for the decision
- evidence that policy and auth meaning stayed the same
- whether the route also survived backend loss
- the sentence that is still forbidden even after the drill

Without that packet, a walkthrough is still disciplined imagination rather than
validated relief.

## The exact stronger sentences still forbidden

Even after this walkthrough, the following are still forbidden today:

- `wrong-node HTTP is basically solved`
- `fallback is mostly handled now`
- `protected route handoff is mostly a middleware problem`
- `TCP routing is close enough to stateful failover`
- `the request path already behaves like one cloud`

## Bottom line

The current runtime already contains a serious request stack:

- Cloudflare-oriented public entry
- Traefik routing
- TinyAuth and forward-auth participation
- CrowdSec enforcement
- Headscale mesh assumptions
- helper-driven fallback intent
- a mix of simple HTTP, protected HTTP, TCP, and stateful services

That is enough to make request preservation a real platform problem.

What the current worktree still does not prove is the part the user actually
cares about:

- that the wrong healthy node already knows what the request should mean
- that it already knows which peer is valid
- that the route still means the same thing after handoff
- that the rescue path survives backend loss
- that TCP and stateful routes have escaped operator folklore

That is the honest boundary today.
The runtime can already execute a lot.
The missing win is still the point where the receiving node can explain itself
without a human finishing the sentence first.
