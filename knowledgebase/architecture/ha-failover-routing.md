# HA, Failover, and Routing

Read this page as the routing burden ledger for the priority implementation
rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

If a routing claim cannot survive contact with the current Compose runtime,
this page is supposed to reject it.

Read these alongside it:

- [Current Compose Runtime](current-compose-runtime.md)
- [Request Path and Failure Walkthrough](request-path-and-failure-walkthrough.md)
- [Failure Model and Maturity Matrix](failure-model-and-maturity.md)
- [Operator Questions and Honest Answers](../operations/operator-questions-and-honest-answers.md)
- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)

## Why this page exists

This repo is not trying to answer the easy question:

> can several machines exist, all run Docker, all publish services, and all
> have Traefik somewhere in the picture?

That question was answered the moment the repo grew past a single host.

The real question is harsher:

> if traffic lands on the wrong healthy node, can that node still preserve the
> meaning of the request without the operator privately remembering the rest of
> the story?

That is the routing benchmark.

It is why the repo keeps circling around:

- `services.yaml`
- peer-aware forwarding
- route durability
- middleware continuity
- the difference between "reachable" and "actually preserved"

## The strongest routing dream

The strongest intent surface is still
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).

It says the desired operating model is:

- no heavyweight orchestrator by default
- manual service placement remains acceptable
- any healthy public node can receive the first request
- local-first serving should happen when the target is already on that node
- peer-forward fallback should happen when the target is remote
- the edge should preserve auth, middleware, and policy meaning
- anti-SPOF pressure should not be flattened into fake HA language

That is the dream.

The file even gives the target routing contract in plain text:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That directness matters because this page is not really about "routing" in the
generic reverse-proxy sense.
It is about whether that contract survives reality.

The rest of this page is about how much of that dream the tracked runtime
actually earns.

## What "multiple Docker nodes" has to mean to be worth anything

For this repo, "multiple nodes" is worthless if it only means:

- more than one box exists
- more than one box can be reached
- more than one box can terminate TLS
- more than one box can be named in DNS

The phrase only becomes meaningful when the wrong first hop can still produce
the right service outcome without private operator completion.

That is why the routing question is socially harsher than a normal proxy
question.
The issue is not only whether packets can move.
It is whether humiliation disappears when the first receiving node is not the
owner node.

## The fake wins this page is supposed to reject

This page exists partly to reject routing stories that sound close enough but
still leave the wound alive.

Examples of fake wins:

- more than one public IP exists, therefore HA is basically present
- Traefik has the route, therefore the wrong node can preserve it
- `docker-gen-failover` wrote fallback-shaped config, therefore fallback is
  solved
- the protected page still answers, therefore the policy meaning survived
- the TCP port is reachable, therefore the stateful service became resilient

Those are not tiny wording mistakes.
They are exactly how the repo would drift back into the same ambiguity that
made the user angry in the first place.

## What the current runtime materially contains

The routing story is not hypothetical.
The priority Compose stack already ships real ingress, policy, and transport
surfaces:

- the root stack includes directly routed HTTP services such as
  `chat-analytics`, `searxng`, `code-server`, `dozzle`, `homepage`,
  `portainer`, `wishlist`, and other label-driven routes
- the root stack includes TCP-exposed surfaces such as `mongodb`, `redis`, and
  `biodecompwarehouse*` through Traefik TCP routers
- the `compose/docker-compose.coolify-proxy.yml` fragment includes the central
  edge surfaces: `traefik`, `tinyauth`, `nginx-traefik-extensions`,
  `crowdsec`, `cloudflare-ddns`, `docker-gen-failover`, `whoami`,
  `logrotate-traefik`, and `autokuma`
- protected routes already exist in the live authoring surface through labels
  like `traefik.http.routers.code-server.middlewares: nginx-auth@file`,
  `traefik.http.routers.dozzle.middlewares: nginx-auth@file`, and
  `traefik.http.routers.portainer.middlewares: nginx-auth@file`
- TinyAuth is not merely an idea; the proxy fragment defines
  `traefik.http.middlewares.tinyauth.forwardAuth.address:
  http://auth:3000/api/auth/traefik`
- Headscale is live and internet-facing in
  [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml),
  including both `headscale-server` and the `headscale` UI
- `docker-gen-failover` is real and explicitly writes generated route material
  to `${CONFIG_PATH:-./volumes}/traefik/dynamic/failover-fallbacks.yaml`

So the problem is not "there is no routing stack."

The problem is that the stack still does not clearly own all the truth needed
to make wrong-node requests dignified.

## Why `docker-gen-failover` does not get to claim more than it earned

`docker-gen-failover` is exactly the kind of component that can flatter a repo
into sounding closer to solved than it is.

Its presence proves:

- the repo is actively trying to generate fallback-aware routing material
- failover is a tracked concern in the live authoring surface
- the edge is not being treated as purely static

Its presence does not prove:

- that the generated rescue route survives the preferred-backend failure
- that the wrong node knows when to trust the generated route
- that a protected route preserves the same meaning after the handoff
- that the helper has escaped the hidden-control-plane trap

That distinction matters because the user is specifically tired of options that
look dynamic while still requiring a human to know when they are lying.

## What the routing stack still does not prove

The current worktree does **not** yet prove:

- that the root runtime ships and consumes a live tracked root
  `services.yaml` or equivalent placement-truth source
- that any named HTTP route already succeeds generically after wrong-node
  entry
- that preferred-backend loss leaves behind a durable rescue route
- that auth and middleware continuity are preserved after peer-forward handoff
- that TCP routing implies safe stateful substitution
- that Headscale, MongoDB, Redis, PostgreSQL, RabbitMQ, or Qdrant inherited
  real HA dignity from the existence of routing components

Those absences are what keep "HA" from being an honest one-word summary.

## The routing burden in one table

| Truth the platform needs | Why it matters | Strongest current evidence | Current gap |
| --- | --- | --- | --- |
| Placement truth | the receiving node must know where the service lives now | recurring `services.yaml` pressure in repo intent surfaces | no clear live tracked root placement authority consumed by routing |
| Peer-eligibility truth | not every reachable peer is safe for every route | Headscale and private-mesh assumptions are real | reachability is not yet proven equivalent to safe route eligibility |
| Route-durability truth | fallback only matters if it survives the failure | `docker-gen-failover` is live and generates dynamic Traefik material | helper presence is weaker than backend-loss persistence proof |
| Semantic-continuity truth | wrong-node handoff must preserve route meaning | protected-route middleware surfaces are real | continuity across peer-forward handoff remains unproven |
| Stateful-authority truth | TCP exposure is not the same as HA | TCP routers for Redis and MongoDB exist | authority, promotion, and write dignity remain separate unresolved problems |

Until those truths become more system-owned, the operator remains the safest
distributed control plane in the stack.

## The route classes that must never be flattened

This repo becomes misleading the moment these lanes are narrated as one shared
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

- Cloudflare can point more than one record at the domain
- Traefik is healthy
- a local request returns `200`
- `docker-gen-failover` generated something that looks fallback-shaped

This is the lane where the repo has the best chance of earning a real early
victory.
But that is only true if the victory is narrow and embarrassing enough to be
believable:

- one named route
- one wrong first hop
- one inspectable peer decision
- one preserved outcome

Anything broader than that is probably bluff again.

### 2. Protected HTTP

This lane is stricter because the route has policy meaning, not just transport
meaning.

The current runtime already shows protected-route intent in real authoring
surfaces:

- `code-server` uses `nginx-auth@file`
- `dozzle` uses `nginx-auth@file`
- `portainer` uses `nginx-auth@file`
- multiple metrics and admin surfaces also use `nginx-auth@file`
- TinyAuth exists as a live forward-auth building block

Protected wrong-node success would mean the forwarded request still preserves:

- auth challenge behavior
- forward-auth behavior
- middleware ordering
- header and trust-boundary assumptions
- the same visible route semantics from the user's point of view

This is why `the page still loads` is too weak.

A forwarded protected route that answers but no longer behaves like the same
protected service is not a successful handoff.

### 3. Raw TCP

This lane already exists in the live runtime through Traefik TCP routers for
`mongodb`, `redis`, and other TCP-oriented surfaces.

That proves real transport exposure.

It does **not** prove:

- safe peer-aware failover
- state-authority transfer
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
- `qdrant`
- `headscale-server`

Headscale is the clearest example of why the lane must stay strict.
Its current config still uses SQLite at `/var/lib/headscale/db.sqlite`.

That means a clean public route to Headscale is not the same thing as
multi-node authority.

The route can be real while the authority is still singular.

This is where a lot of self-hosting conversations become actively misleading.
They talk as if a reachable control plane, a public hostname, or a replica-ish
story means the state escaped the node.

For this repo, the stricter reading is:

> if the write authority, promotion rule, recovery order, and data truth still
> need a human narrator, the service is still singular in the exact way that
> matters.

## The wrong-node event decomposed honestly

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

## The route-level proof packet this repo actually needs

For this repo, a routing claim only starts to become real when it can be
reduced to one named proof packet:

1. name one route
2. name the node that receives the request
3. show that the node does not host the target locally
4. show what shared truth told it where to go
5. show why the chosen peer was eligible rather than merely reachable
6. show that the forwarded route still meant the same thing
7. if failover is claimed, show the route after the preferred backend died
8. show the resulting evidence without requiring private operator narration

That packet is intentionally narrow.

This repo does not need another hundred broad routing promises nearly as much
as it needs one route-level packet that survives embarrassment honestly.

## The human question hiding inside the routing question

The routing contract is not only technical.
It is also social:

> if someone wakes up later and asks why the wrong healthy node still produced
> the right answer, does the system have the explanation or does one operator?

That is why this page keeps sounding harsher than normal proxy documentation.
The user's real complaint is not that Traefik lacks features.
It is that the surrounding stack keeps offering "options" without truly
removing the hidden human SPOF.

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

If the answer remains `the operator` or `a half-manual helper path`, then the
system still has a real human SPOF even though the proxy layer looks serious.

## Why `docker-gen-failover` does not get automatic credit

The repo is right to keep paying attention to `docker-gen-failover`.
It is one of the most concrete attempts at dynamic fallback in the current
runtime.

But the current Compose evidence shows it as:

- a generated-file helper
- fed from Docker events
- writing failover material to a dynamic Traefik file
- not continuously proven as a backend-loss survivor

Research pages in this repo already preserve the harder warning:

> `docker-gen-failover` can still delete routes on container stop and therefore
> fail the exact backend-loss scenario that made fallback matter.

That makes it an excellent example of the user's complaint.

A thing can sound like failover, generate fallback-shaped config, and still
evaporate during the exact failure that made fallback necessary.

So the docs should keep speaking of it as:

- highly relevant
- concrete
- worth measuring
- not yet victory

## Where the current routing story is strongest

The live stack is strongest in these areas:

- it already has serious edge machinery rather than a toy reverse proxy
- it already separates HTTP and TCP in actual runtime labels
- it already has protected-route surfaces, not merely public happy-path demos
- it already has observability around edge components through the metrics stack
- it already exposes helper and control components such as
  `docker-gen-failover`, `cloudflare-ddns`, and Headscale rather than hiding
  them

That means the docs do not need to pretend the repo is earlier than it is.

## Where the current routing story is still unfinished

The live stack is still unfinished in exactly the places that matter most to
the user:

- placement truth still appears architecturally central without being clearly
  rooted in a live tracked runtime source
- peer selection still looks more like a missing join than a proven platform
  behavior
- `docker-gen-failover` is directionally relevant but still a helper, not a
  proven backend-loss survivor
- protected-route continuity is inferable from labels but not yet proven
  across wrong-node handoff
- stateful services still have transport exposure long before they have
  multi-node authority dignity

That is why the repo still feels constrained by hidden human SPOFs even though
it already has many moving parts.

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
- Cloudflare plus Traefik removed the hidden operator SPOF
- TCP routing implies stateful resilience
- the platform already behaves like one cloud in the user's intended sense
- the helper stack has already become a real distributed control plane

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
