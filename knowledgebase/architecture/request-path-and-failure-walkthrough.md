# Request Path and Failure Walkthrough

This page reconstructs the request path the user is actually asking about.

Not:

> what proxy do we use?

and not:

> what containers exist at the edge?

The real question is:

> when a hostname is requested, what literally happens next, what truth does
> the receiving node own at each hop, and at what point does the explanation
> turn back into private operator memory?

That is the only honest way to document "multi-node without Swarm" in this
repo.

## What this page is and is not allowed to prove

This page is allowed to prove:

- the target request contract the repo keeps trying to earn
- which live runtime components already participate in that contract
- which route classes are materially present in the current Compose runtime
- where proof stops today for wrong-node, peer-forward, and backend-loss
  behavior
- which private sentence still survives after each hop

This page is not allowed to prove:

- generic wrong-node success today
- generic peer-forward correctness today
- generic backend-loss survival today
- that a plausible walkthrough equals a validated drill
- that HTTP optimism transfers into TCP or stateful authority

This is a decomposition page.
It is not a green-check page.

## The fake walkthrough this page has to stop

The bad version of this page would say:

1. DNS resolves
2. Traefik matches
3. auth middleware exists
4. helper containers exist
5. therefore the distributed request path is mostly in place

That shape is exactly what this page exists to prevent.

A walkthrough can be detailed and still be useless if it quietly skips the
moment where the receiving node stops knowing and the operator starts finishing
the sentence privately.

## The target contract

The clearest contract still lives in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That is stricter than:

- more than one public node exists
- the hostname matches
- a helper can generate fallback-looking config
- the route can probably be proxied somewhere

The contract is about preserved request meaning after locality fails.

## The one event the whole repo keeps circling

The entire architecture pressure can be reduced to one embarrassing event:

1. a public request lands on a healthy node
2. the requested service is not local there
3. that node must decide what to do
4. the decision must be correct for the route class
5. the explanation must exist outside one operator's head

If a page, helper, or architecture proposal does not improve that event, it is
probably orbiting the problem instead of solving it.

## The live request actors already present now

From the root runtime and active fragments, the request path already contains
real actors.

### Public-entry and naming actors

- `cloudflare-ddns`
- Cloudflare-first public entry assumptions

### Main L7 execution actors

- `traefik`
- Traefik router and service labels across root and fragment files

### Policy and auth actors

- `nginx-traefik-extensions`
- `tinyauth`
- `crowdsec`
- protected routes using `nginx-auth@file`

### Fallback and helper actors

- `docker-gen-failover`
- generated output to
  `${CONFIG_PATH:-./volumes}/traefik/dynamic/failover-fallbacks.yaml`

### Private reachability actors

- `headscale-server`
- `headscale`

### Proof-candidate services

- `whoami`
- `wishlist`
- `mkdocs`

The runtime also already contains real route classes:

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
  - `prometheus`
  - `alertmanager`
  - `mcpo`
- TCP transport:
  - `mongodb`
  - `redis`
  - `biodecompwarehouse*`
- state-bearing app subgraphs:
  - `nuq-postgres`
  - `rabbitmq`
  - `litellm-postgres`
  - `qdrant`
  - Headscale SQLite authority

So the request-path problem is not theoretical.
The stack is already large enough that missing truth has become socially
expensive.

## The private-sentence test for this page

At the end of each walkthrough section, the reader should be able to answer:

> what exact sentence would still need private operator completion if this
> request landed on the wrong healthy node right now?

If a section leaves behind only components and flow language, but not that
surviving sentence, the walkthrough is still too soft.

## The one question to keep asking at every hop

At every stage, keep asking:

> what does the receiving node know from shared inspectable system-owned truth,
> and what would still have to be privately supplied by the operator if the
> local backend were absent?

If that question disappears, the walkthrough can stay technical while still
quietly collapsing back into theater.

## Walkthrough 1: normal-day stateless HTTP

Use `whoami.$DOMAIN`, `wishlist.$DOMAIN`, or `mkdocs.$DOMAIN` as the cleanest
mental model.

These are useful first-pass routes because they are:

- HTTP
- low-state
- narrow
- good proof candidates
- less likely to hide state-authority ambiguity

They are also the fairest early test of the repo's promise.
If the stack cannot make one low-state HTTP route survive wrong-node entry with
an inspectable explanation, louder HA language elsewhere is premature.

### Hop 1: DNS and public entry

What is live now:

- Cloudflare DDNS is in the active edge fragment
- the repo explicitly rejects one sacred public box as the only front door

What this honestly proves:

- first-hop plurality is part of the live design
- more than one node is intended to be a valid public entry candidate

What it still does not prove:

- the receiving node knows where the service lives now
- the receiving node knows whether the target is local or remote
- the receiving node knows which peer is valid if the target is remote

Private sentence still surviving after Hop 1:

> traffic can hit more than one node, but I still personally know that this
> says nothing about what the chosen node should do next.

### Hop 2: Traefik accepts the hostname

What is live now:

- Traefik is materially present
- root and fragment services carry real routers and services

What this honestly proves:

- the receiving node can parse the hostname
- the receiving node can match a route definition
- the route exists in authored runtime truth

What it still does not prove:

- that the route owns cross-node placement truth
- that the receiving node can choose a remote backend from shared current truth

This hop is route execution.
It is not route truth.

Private sentence still surviving after Hop 2:

> the route matched, but I still personally know that route execution is
> weaker than distributed backend truth.

### Hop 2.5: the missing decision point

This is the hop infra prose loves to skip because it is socially awkward.

The missing question is:

> once the route is known, what tells this node whether it is the owner node
> or merely the first-hop node?

That is where placement truth and peer-eligibility truth would have to appear
if the contract were becoming real.

Without that layer, the walkthrough quietly mutates into:

1. the hostname matched
2. a human knows what should happen
3. therefore the route story feels understandable

That mutation is exactly what this repo is trying to escape.

### Hop 3: middleware and auth policy are applied

What is live now:

- `nginx-traefik-extensions` participates in forward-auth behavior
- `tinyauth` participates in real auth handling
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

Private sentence still surviving after Hop 3:

> the same middleware names exist, but I still personally know that this does
> not prove the remote handoff preserves the same route meaning.

### Hop 4: local service execution

If the target service is local to the receiving node, the story gets stronger.

What this honestly proves:

- the local happy path works for that route class
- the edge and the app are at least locally connected
- the service is not merely theoretical

What it still does not prove:

- wrong-node dignity
- peer-eligibility truth
- backend-loss rescue

The local happy path is real progress.
It is still smaller than the user's benchmark.

Private sentence still surviving after Hop 4:

> the owner node can serve the route, but I still personally know that this
> says very little about wrong-node dignity.

## Walkthrough 2: protected HTTP

Now use a route such as:

- `code-server.$DOMAIN`
- `dozzle.$DOMAIN`
- `portainer.$DOMAIN`
- `prometheus.$DOMAIN`
- `mcpo.$DOMAIN`

These are more revealing because the system must preserve not just transport,
but meaning:

- auth behavior
- middleware order
- trust boundaries
- forward-auth behavior
- header assumptions
- visible route semantics

### What the runtime already gives this class

The current runtime already gives protected routes:

- real Traefik routers
- `nginx-auth@file`
- TinyAuth participation
- middleware-bearing edge surfaces
- observability around edge and helper components

### What remains unproven

The current worktree still does not prove:

- that a protected route forwarded to another node still means the same thing
- that the peer has converged enough secrets and config to preserve that
  meaning
- that the same protected route survives backend loss with the same visible
  contract

That is why "the page still loads" is too weak.
A route that answers but no longer behaves like the same protected service is
not a successful handoff.

## Walkthrough 2.5: the false-comfort transition

This repo repeatedly risks one exact illegitimate transition:

- a route exists locally
- the route is protected
- the edge stack is serious
- therefore a peer-forward version of the route must be close enough

That transition stays illegal until the docs can show:

- the source of placement truth
- the source of peer-eligibility truth
- the continuity of middleware and auth semantics
- the continuity of the route under backend loss if failover is being claimed

Without those, the walkthrough is still leaning on emotional plausibility.

## Walkthrough 3: wrong-node entry

This is the scene the user actually cares about.

Assume all of the following:

1. Cloudflare sends traffic to a healthy public node
2. Traefik, auth, and edge helpers all look healthy
3. the requested service is not local to that node

At that moment, the repo becomes very good at describing the problem and still
not strong enough to prove the general solution.

### What the receiving node can likely know already

From the live edge stack, the receiving node can likely know:

- the hostname
- the route class
- the middleware chain it would normally apply
- whether the edge stack itself looks healthy

### What the worktree does not yet prove the node knows on its own

The current worktree does not yet prove the receiving node has shared truth
for:

- exact live placement of the target service
- whether the relevant peer is healthy enough to receive the request
- whether that peer has converged secrets, environment, and policy semantics
- whether the forwarded request would still mean the same thing after handoff

That is the missing middle layer in one scene.

Private sentence still surviving in the wrong-node scene:

> the wrong node can see the request shape, but I still personally know facts
> the runtime has not yet proven it owns.

## Walkthrough 4: backend-loss scene

This is a separate trap and it must stay separate.

Even if wrong-node forwarding looks plausible while everything is healthy, the
rescue path may still disappear once the preferred backend actually fails.

That is why `docker-gen-failover` matters and why it still does not get to
claim victory.

Another way to say the distinction:

- wrong-node success asks whether the first-hop mistake can be corrected
- backend-loss success asks whether the correction path survives the damage

The user is frustrated because many options solve one of those socially and not
the other operationally.

Wrong-node success would mean:

- the request landed on a non-owner node
- the system still chose a valid backend

Backend-loss success would mean:

- the preferred backend disappeared
- the rescue path still existed
- the route still meant the same thing
- the operator did not have to privately finish the story

Those proofs overlap.
They are not the same proof.

Private sentence still surviving in the backend-loss scene:

> the route looked recoverable while healthy, but I still personally know the
> rescue path may evaporate during the exact failure event.

## Walkthrough 5: TCP and stateful lanes

This is where many infrastructure discussions become dishonest fastest.

The runtime already contains:

- Traefik TCP routers for `mongodb`, `redis`, and other raw transport surfaces
- state-bearing systems such as `nuq-postgres`, `rabbitmq`,
  `litellm-postgres`, `qdrant`, and Headscale

That proves there is real transport machinery.
It does not prove:

- safe peer-aware failover
- write-authority transfer
- correct client behavior after rerouting
- promotion rules
- operational dignity under node loss

Headscale is the cleanest example.
The current fragment still points its authority to SQLite at
`/var/lib/headscale/db.sqlite`.

So even if the route is clean and the hostname is plural, the authority can
still be singular.

Private sentence still surviving in the stateful lane:

> the service is reachable, but I still personally know who the real writer is
> and what happens if that node disappears.

## The actual deliverable this page is pushing toward

The real deliverable is not a smoother explanation of a hypothetical flow.
It is one preserved route packet that can answer all of these:

- which node received the request first?
- why was it the wrong node?
- what shared truth named the real owner?
- why was the chosen peer eligible now?
- what visible route semantics stayed the same?
- what artifact lets a second operator inspect the answer later?
- if failover is claimed, what survived after the preferred backend died?

Until one packet like that exists, this page is primarily a warning page.

## The honest bottom line

The stack already contains real public-entry pressure, real edge machinery,
real auth layers, real helper layers, and real stateful surfaces.

What it does not yet prove is that the decisive handoff truth has moved out of
private human custody.

So the honest sentence today is:

> the repo can already explain the distributed request problem in concrete
> runtime terms, but the most important decisions in the wrong-node and
> backend-loss path are still only partially system-owned.
