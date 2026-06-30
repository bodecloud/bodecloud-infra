# HA, Failover, and Routing

Read this page as the routing burden ledger for the priority implementation
rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

If a route claim cannot survive contact with the current Compose runtime, this
page is supposed to reject it.

Read these alongside it:

- [Current Compose Runtime](current-compose-runtime.md)
- [Request Path and Failure Walkthrough](request-path-and-failure-walkthrough.md)
- [Failure Model and Maturity Matrix](failure-model-and-maturity.md)
- [Operator Contract and Success Criteria](operator-contract.md)
- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)

## Why this page exists

This repo is not trying to answer the easy question:

> can several machines exist, all run Docker, all publish services, and all
> have Traefik somewhere in the picture?

That question was already answered the moment the repo grew past a single host.

The real question is harsher:

> if traffic lands on the wrong healthy node, can that node still preserve the
> meaning of the request without the operator privately remembering the rest of
> the story?

That is the routing benchmark.

It is the reason the repo keeps circling around `services.yaml`, peer-aware
fallback, route durability, and the difference between "reachable" and
"actually preserved."

## Strongest architecture intent

The strongest intent surface is still
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It says the desired operating model is:

- no heavyweight orchestrator by default
- manual service placement is acceptable
- any healthy public node can receive the first request
- local-first serving should happen when the target is already on that node
- peer-forward fallback should happen when the target is remote
- the edge should preserve auth, middleware, and policy meaning
- anti-SPOF pressure should not be flattened into fake HA language

That is the dream.

The rest of this page is about how much of that dream the tracked runtime
actually earns.

## What the current runtime materially contains

The routing story is not hypothetical. The priority Compose stack already ships
real ingress, policy, and transport surfaces:

- the root stack includes directly routed HTTP services such as
  `chat-analytics`, `searxng`, `code-server`, `dozzle`, `homepage`,
  `portainer`, and others through Traefik labels
- the root stack includes TCP-exposed stateful surfaces such as `mongodb` and
  `redis` through Traefik TCP routers
- the `compose/docker-compose.coolify-proxy.yml` fragment includes the central
  edge surfaces: `traefik`, `tinyauth`, `nginx-traefik-extensions`,
  `crowdsec`, `cloudflare-ddns`, `docker-gen-failover`, `whoami`, and
  `autokuma`
- protected routes already exist in the live authoring surface through labels
  like `traefik.http.routers.code-server.middlewares: nginx-auth@file`,
  `traefik.http.routers.dozzle.middlewares: nginx-auth@file`, and
  `traefik.http.routers.portainer.middlewares: nginx-auth@file`
- TinyAuth is not merely an idea; the proxy fragment defines
  `traefik.http.middlewares.tinyauth.forwardAuth.address:
  http://auth:3000/api/auth/traefik`
- Headscale is real and internet-facing in
  `compose/docker-compose.headscale.yml`, including both `headscale-server`
  and the `headscale` UI with Traefik routers and redirects
- `docker-gen-failover` is real and explicitly writes generated route material
  to `${CONFIG_PATH:-./volumes}/traefik/dynamic/failover-fallbacks.yaml`

So the problem is not "there is no routing stack."

The problem is that the stack still does not clearly own all the truth needed
to make wrong-node requests dignified.

## What the routing stack still does not prove

The current worktree does **not** yet prove:

- that the root runtime ships and consumes a live tracked root
  `services.yaml` or equivalent placement-truth source
- that any named HTTP route already succeeds generically after wrong-node
  entry
- that a preferred-backend loss leaves behind a durable rescue route
- that auth and middleware continuity are preserved after peer-forward handoff
- that TCP routing implies safe stateful substitution
- that Headscale, MongoDB, Redis, PostgreSQL, or RabbitMQ have inherited real
  HA dignity from the existence of routing components

Those absences are what keep "HA" from being an honest one-word summary.

## The routing classes that must not be flattened

This repo becomes misleading the moment these route classes are narrated as one
 success story.

### 1. Stateless HTTP

This is the lane most likely to get real relief first.

Candidate surfaces already exist:

- `whoami`
- `wishlist`
- `chat-analytics`
- `searxng`
- `homepage`
- parts of the site shell itself

For stateless HTTP, wrong-node success would mean:

1. a healthy public node receives the request
2. that node knows whether the service is local or remote
3. if remote, it knows which peer is eligible now
4. it forwards without changing the user-visible meaning of the route
5. if the preferred backend disappears, the rescue path still exists

What does **not** prove this lane:

- Cloudflare can point more than one A record at the domain
- Traefik is healthy
- a local request returns `200`
- `docker-gen-failover` generated something that looks fallback-shaped

### 2. Protected HTTP

This lane is stricter because the route has policy meaning, not just transport
meaning.

The current runtime already shows protected-route intent in real authoring
surfaces:

- `code-server` uses `nginx-auth@file`
- `dozzle` uses `nginx-auth@file`
- `portainer` uses `nginx-auth@file`
- many metrics/admin surfaces also use `nginx-auth@file`
- TinyAuth exists as a live forward-auth building block

Protected wrong-node success would mean the forwarded request still preserves:

- auth challenge behavior
- forward-auth behavior
- middleware ordering
- header and trust-boundary assumptions
- the same visible route semantics from the user perspective

This is why "the page still loads" is too weak.

A forwarded protected route that answers but no longer behaves like the same
protected service is not a successful handoff.

### 3. Raw TCP

This lane already exists in the live runtime through Traefik TCP routers for
`mongodb` and `redis`, plus other TCP-oriented surfaces.

That proves real transport exposure.

It does **not** prove:

- safe peer-aware failover
- state authority transfer
- correct client semantics after rerouting
- operational dignity under node loss

TCP routing is an execution tool.
It is not a substitute for stateful ownership semantics.

### 4. Stateful routes

This is the harshest lane and the one most likely to be lied about if the docs
get lazy.

The current stack includes real stateful dependencies:

- root `mongodb`
- root `redis`
- `nuq-postgres`
- `litellm-postgres`
- `rabbitmq`
- `headscale-server`

Headscale is the clearest example of why the lane must stay strict.
Its current config still uses SQLite at `/var/lib/headscale/db.sqlite`.

That means a clean public route to Headscale is not the same thing as
multi-node state authority.

The route can be real while the authority is still singular.

## The actual routing burden

The repo keeps rediscovering the same missing middle layer:

- the node that receives a request needs current placement truth
- it also needs peer eligibility truth
- it also needs a durable route artifact that survives the failure that
  triggered fallback
- it also needs to preserve route semantics after handoff

Without those truths, the node can be reachable and still "behave stupidly"
under wrong-node entry.

This is why `services.yaml` keeps reappearing.
The repo is trying to externalize the sentence:

> where does this service actually live right now, and who says that is true?

Until the platform owns that answer, the operator is still the safest control
plane.

## The wrong-node event, decomposed honestly

When the user says several Docker nodes should behave like one cloud, the
critical event is not "two nodes exist."

It is this event:

1. Cloudflare resolves a hostname to one of several healthy public nodes.
2. The request lands on a node that does **not** host the target service.
3. That node must decide whether to serve locally or forward remotely.
4. If remote, it must choose a peer that is eligible **now**, not one that was
   probably right yesterday.
5. The forwarded request must still be the same route, not a looser
   approximation.
6. If the preferred backend is gone, the rescue route must still exist.
7. The operator should later be able to inspect the decision from artifacts,
   not folklore.

Every architecture choice in this repo is downstream of that event.

## Where the current runtime is strongest

The live stack is strongest in these areas:

- it already has serious edge machinery rather than a toy reverse proxy
- it already separates HTTP and TCP in actual runtime labels
- it already has protected-route surfaces, not merely public happy-path demos
- it already has observability around edge components through the metrics stack
- it already exposes helper/control components such as `docker-gen-failover`,
  `cloudflare-ddns`, and Headscale rather than hiding them

That means the docs do not need to pretend the repo is earlier than it is.

## Where the current runtime is still unfinished

The live stack is still unfinished in exactly the places that matter most to
the user:

- placement truth still appears architecturally central without being clearly
  rooted in a live tracked runtime source
- peer selection still looks more like a missing join than a proven platform
  behavior
- `docker-gen-failover` is directionally relevant but currently configured as a
  generated dynamic-route helper, not a proven backend-loss survivor
- protected-route continuity is inferable from labels but not yet proven across
  wrong-node handoff
- stateful services still have transport exposure long before they have
  multi-node authority dignity

That is why the repo still feels constrained by hidden human SPOFs even though
it already has many moving parts.

## Why Cloudflare is necessary but insufficient

Cloudflare is a real part of the anti-SPOF story because it helps prevent one
sacred public node from becoming the only entrypoint.

But Cloudflare only buys first-hop plurality.

It does not, by itself, answer:

- what if the chosen node does not host the service?
- how does that node know where the service actually lives now?
- how does it know whether forwarding is safe?
- how does it keep auth and middleware meaning intact?
- what survives when the preferred backend is gone?

If the docs let plural DNS impersonate preserved requests, they are describing
a smaller success than the user asked for.

## Why Traefik is necessary but insufficient

Traefik is central to the repo's actual runtime shape.
The stack already uses it for:

- HTTP routers
- middleware attachment
- TCP routers
- service port exposure
- redirects
- protected admin surfaces

That matters.

But Traefik is an execution plane.
It does not magically create distributed truth.

It can route based on what it knows.
The unresolved question is still:

> who gave it the current truth needed to make the right distributed decision?

If the answer remains "the operator" or "a half-manual helper path," then the
system still has a real human SPOF even though the proxy layer looks serious.

## Why `docker-gen-failover` does not get automatic credit

The repo is right to keep paying attention to `docker-gen-failover`.
It is one of the most concrete attempts at dynamic fallback in the current
runtime.

But the current Compose evidence shows it as:

- a generated-file helper
- fed from Docker events
- writing failover material to a dynamic Traefik file
- not continuously claimed as a proven backend-loss survivor

Research pages in this repo already preserve the harder warning:

> `docker-gen-failover` currently deletes routes on container stop

That makes it an excellent example of the user's complaint.

A thing can sound like failover, generate fallback-shaped config, and still
evaporate during the exact failure that made fallback necessary.

So the docs should keep speaking of it as:

- highly relevant
- concrete
- worth measuring
- not yet victory

## The route truths the platform still needs to own

For the routing dream to become believable, the platform needs to own at least
these truths explicitly:

| Truth | Why it matters | Current strongest evidence | Current gap |
| --- | --- | --- | --- |
| Placement truth | the receiving node must know where the service lives now | repeated `services.yaml` pressure in repo intent surfaces | no clear live tracked root placement authority consumed by routing |
| Peer eligibility truth | not every reachable peer is safe for every route | Headscale and private-network assumptions are real | reachability is not yet proven equivalent to safe route eligibility |
| Route durability truth | fallback matters only if it survives the failure | `docker-gen-failover` is live | helper existence is weaker than backend-loss persistence proof |
| Semantic continuity truth | wrong-node handoff must preserve route meaning | protected-route middleware surfaces exist | continuity across peer-forward handoff remains unproven |
| Stateful authority truth | TCP exposure is not the same as HA | TCP routers for Redis and MongoDB exist | authority, promotion, and write-dignity remain separate problems |

## What this page allows the docs to say today

The docs are allowed to say:

- the repo has a serious Compose-first multi-node ingress stack
- the any-node-entry dream is explicit and coherent
- the routing problem has already been decomposed correctly
- wrong-node dignity is the real benchmark
- protected routes and stateful routes must stay on harsher tracks than happy
  stateless HTTP

The docs are **not** allowed to say:

- wrong-node behavior is basically solved
- fallback is mostly handled now
- Cloudflare plus Traefik has removed the hidden operator SPOF
- TCP routing implies stateful resilience
- the platform now behaves like one cloud in the user's intended sense

## What would materially change the routing story

The smallest sequence that would actually change this page is:

1. expose one live current-state placement truth surface
2. prove one stateless wrong-node HTTP route end to end
3. re-run that route under preferred-backend loss
4. compare one protected route before and after peer-forward handoff
5. keep TCP and stateful claims on stricter independent tracks

That sequence matters because it moves the burden out of the operator's head.

Until that happens, the honest routing sentence is:

> the repo has real ingress machinery, real anti-SPOF intent, and real
> peer-aware pressure, but the most important distributed routing truths are
> still only partially system-owned
