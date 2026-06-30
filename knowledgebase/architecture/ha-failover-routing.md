# HA, Failover, and Routing

Read this page as the routing truth map for the priority implementation rooted
at [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

If you want the deeper evidence stack first, read:

- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
- [Request Path and Failure Walkthrough](request-path-and-failure-walkthrough.md)
- [Current Compose Runtime](current-compose-runtime.md)
- [Operator Contract and Success Criteria](operator-contract.md)

This page exists because "HA" becomes almost useless in this repo unless
routing is split into the separate truths the user is actually asking about.

The user is not mainly asking:

- can more than one node receive traffic?
- is Traefik present?
- do healthchecks exist?
- is there a helper that sounds failover-shaped?

The harder question is:

> if traffic lands on the wrong healthy machine, can that machine still
> preserve the service contract without requiring private operator memory?

That is the routing problem here.

## What this page is and is not allowed to prove

This page is authoritative about:

- the routing classes that matter in this repo
- the target any-node-entry contract the repo keeps converging on
- the current runtime surfaces that already participate in routing
- the truths that still remain operator-owned
- the difference between first-hop plurality and request-preserving recovery

This page is not authoritative about:

- whether a specific hostname has already passed a full wrong-node drill
- whether backend-loss fallback is already durable for a named route
- whether stateful traffic inherited optimism from the HTTP story
- whether one working helper path upgrades the whole platform

This is a routing decomposition page, not a victory lap.

## The routing dream in one sentence

The strongest intent surface is
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It describes the desired operating contract as:

- any healthy public node can receive the first request
- if the service is local, that node serves it locally
- if the service is remote, that node forwards to a healthy peer that
  currently hosts it
- the edge should preserve auth, middleware, and policy meaning rather than
  merely expose containers

That target is anti-SPOF and peer-aware without immediately collapsing into
Swarm, Kubernetes, k3s, or another heavyweight orchestrator.

It is also only intent unless the runtime proves it.

## The routing classes that must stay separate

Routing in this repo splits into at least four classes.
The documentation gets misleading the moment these collapse into one HA story.

### 1. Stateless HTTP

This is the easiest class to describe and the easiest class to overclaim from.

Wrong-node success for stateless HTTP would mean:

- a healthy public node receives the request
- that node determines whether the service is local or remote
- if remote, it chooses an eligible peer
- forwarding preserves the route semantics
- if the preferred backend disappears, a still-valid fallback can take over

Things that do **not** prove stateless HTTP failover here:

- Cloudflare can resolve to multiple public IPs
- Traefik is healthy
- a local route returns `200`
- a helper renders fallback-shaped config

### 2. Protected HTTP

This class is stricter than ordinary stateless HTTP because the system has to
preserve:

- auth behavior
- middleware order
- headers and trust boundaries
- the exact meaning of the protected route

Wrong-node success here means more than "the response still comes back."
It means the forwarded route is still the same protected service in
user-visible, auth-visible, and policy-visible terms.

### 3. Raw TCP

This class is different enough that HTTP optimism cannot spill into it.

A TCP router for Redis or MongoDB does **not** prove:

- peer-aware service failover
- semantic continuity
- safe substitution between backends
- state correctness

It proves transport exposure, not full replacement safety.

### 4. Stateful surfaces

This is the strictest class.
State-bearing services do not become honestly HA because:

- they are reachable through Traefik
- more than one node exists
- a proxy can point at more than one backend

Stateful routing only becomes meaningful when authority, durability,
substitution safety, and recovery semantics are explicit.

## What the live runtime already contains

The root runtime already contains enough routing and ingress machinery to make
the problem real rather than hypothetical.

Visible anchors include:

- Cloudflare-oriented public-entry assumptions
- Traefik
- TinyAuth
- `nginx-traefik-extensions`
- CrowdSec
- Docker socket proxies
- `docker-gen-failover`
- Headscale mesh assumptions
- TCP routers for services such as MongoDB and Redis

That is a real routing stack.

What it still does **not** clearly own is the thing the user actually wants
that machinery to buy:

- shared placement truth
- trustworthy peer eligibility truth
- durable fallback-route persistence under failure
- proof that a request keeps the same meaning after wrong-node handoff
- proof that stateful authority survives any routing story being told

That missing truth-owning middle layer is why the route story still feels
unfinished.

## The exact truths routing needs

For the routing dream to become believable, the platform has to own more than
"containers are up."
It needs at least these truths:

| Truth | Why it matters | What currently fakes it | What real ownership would look like |
| --- | --- | --- | --- |
| Placement truth | the receiving node must know where the service actually lives now | hostnames, folklore, static assumptions | a tracked current-state source the receiving node can consult |
| Peer eligibility truth | not every reachable peer is safe for every route | "that node seems alive" | health, revision, policy, and route-class checks |
| Route durability truth | fallback only matters if it survives the failure that triggered it | rendered config while everything is healthy | proof that the route still exists after backend loss |
| Semantic continuity truth | forwarded requests must still mean the same thing | transport success or `200` alone | preserved auth, middleware, headers, and user-visible route behavior |
| Stateful authority truth | TCP reachability is weaker than HA | proxy presence or multiple nodes | explicit ownership, durability, and recovery semantics |

If the system still cashes those truths out into private human memory, the
central routing problem is still alive.

## The wrong-node path, step by step

The path the repo wants is:

1. Cloudflare sends the request to any healthy public node.
2. The receiving node decides whether the target service is local or remote.
3. If local, it serves locally.
4. If remote, it chooses an eligible peer from shared truth.
5. The forwarded route preserves the same auth, middleware, and header meaning.
6. If the preferred backend disappears, the system still has a valid fallback
   path.
7. The operator can later inspect why the route was valid without reconstructing
   the story from memory.

The path the repo is still trying to kill is:

1. Cloudflare sends the request to a healthy node.
2. The target service is not local.
3. The operator privately knows where it probably runs.
4. The operator privately knows whether forwarding is safe.
5. The route works only because the operator completed the story off-book.

That second path is the real SPOF the user keeps complaining about.

## What still does not count as HA or failover here

The following may be real progress.
They still do not count as meaningful HA or failover on their own in this repo:

- more than one public node can receive the first hop
- Traefik is present and healthy
- a helper generates fallback-shaped route material
- a TCP router exists for a stateful service
- a local protected route returns `200`
- a mesh exists between nodes
- the docs can explain the route classes more clearly

All of those can improve posture.
None of them satisfy the user's benchmark unless they also reduce the need for
private placement memory, preserve request meaning on the wrong node, and
survive the failure that made fallback necessary.

The checksum question remains:

> yes, but who still had to know the real answer first?

If the answer is still "the operator," then the route may be real software and
still fake relief.

## Why Cloudflare and Traefik are necessary but insufficient

Cloudflare is part of the anti-SPOF story as the first hop.
Traefik is central to the HTTP routing story.
Both matter.
Neither solves the whole problem alone.

Cloudflare can help more than one node become a first hop.
That is weaker than proving the wrong node can preserve the request correctly.

Traefik can express routers, middlewares, services, and TCP paths.
That is weaker than proving:

- the wrong node knows the right backend
- the chosen peer is actually eligible
- the fallback survives backend loss
- the route preserves meaning after handoff

The problem is not lack of ingress components.
It is lack of a shared truth layer that lets the receiving node make the
correct decision without folklore.

## What the archive keeps confirming

The archive is consistent about where the problem really lives:

- `docker-multi-node-without-swarm__...` keeps converging on service discovery
  as the hard missing piece once placement and DNS plurality are accepted
- that same thread frames the real problem as mapping "service name" to
  "where is it running right now"
- `load-balancer-failover-alternatives__...` shows that even feature-rich
  failover products often solve only one edge slice and still leave broader
  route meaning unresolved

That is why the routing burden in this repo is not "find one better proxy."
It is "find or build the smallest layer that lets the wrong node know the
right answer."

## What a real proof packet would need

For this repo, a serious routing proof packet would need to show all of the
following together:

- the exact entry node that first received the request
- whether the requested service was local or remote at that moment
- the current placement truth source the receiving node consulted
- the peer-selection logic that chose the fallback target
- preserved middleware and auth behavior across the handoff
- backend health or loss conditions during the test
- operator-readable artifacts explaining why the route was valid

If any of those are missing, the packet may still be useful debugging
evidence, but it does not close the user's real routing question.

## Bottom line

The routing dream is simple to say:

- any healthy node may receive the request
- local services serve locally
- remote services forward safely
- fallback survives failure
- the operator no longer has to narrate the topology from memory

The current runtime already has enough ingress machinery to make that dream
plausible.
It still lacks the shared truth layer that would make the wrong node know why
it is safe to act.
