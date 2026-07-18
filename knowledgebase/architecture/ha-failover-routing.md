# HA, Failover, and Routing

Read this page as the routing burden ledger for the priority implementation
rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

If a routing claim cannot survive contact with the current Compose runtime,
this page is supposed to reject it.

Read these alongside it:

- [Current Compose Runtime](current-compose-runtime.md)
- [Request Path and Failure Walkthrough](request-path-and-failure-walkthrough.md)
- [Current-State Registry and Peer Eligibility](current-state-registry-and-peer-eligibility.md)
- [Failure Model and Maturity Matrix](failure-model-and-maturity.md)
- [Operator Questions and Honest Answers](../operations/operator-questions-and-honest-answers.md)
- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)

## The actual question this page answers

This repo is not trying to answer the easy question:

> can multiple machines run Docker, publish services, and use Traefik?

That question stopped being interesting the moment the stack became multi-node.

The real question is harsher:

> if traffic lands on the wrong healthy node, can that node still preserve the
> meaning of the request without the operator privately completing the story?

That is the benchmark.

Everything else is secondary:

- extra IPs
- DNS plurality
- Traefik presence
- generated fallback config
- container health

All of those can exist while the human remains the real control plane.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what the routing dream actually is
- which ingress, auth, helper, and transport surfaces are live now
- which route classes exist in the current runtime
- which distributed routing truths remain unowned
- why "HA" is still too generous as a one-word summary

This page is not allowed to prove:

- generic wrong-node success today
- generic peer-forward correctness today
- generic backend-loss survival today
- auth or middleware continuity after cross-node handoff
- TCP or stateful resilience just because a port is routed

This page should stay narrower than ambition and harsher than marketing.

## The routing dream, stated without smoothing

The strongest intent surface remains
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).

Its target operating contract is:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That contract matters because it is stronger than:

- more than one node exists
- more than one node is reachable
- Traefik can parse the host
- a dynamic helper can generate fallback-looking config

It demands preserved request meaning after locality fails.

That means all of the following have to become system-owned rather than
socially reconstructed:

- where the service actually lives now
- which peer is eligible now
- whether the forwarded route still means the same thing
- whether the rescue path survives preferred-backend loss

## "Multiple Docker nodes" means nothing unless the first-hop mistake is survivable

For this repo, the phrase "multiple nodes" is worthless if it only means:

- more than one box exists
- more than one box can terminate TLS
- more than one node can appear in DNS
- more than one node runs the proxy stack

The phrase only becomes meaningful when the wrong first hop can still lead to
the right service outcome without one operator silently carrying the answer.

That is the entire wound.

## The fake wins this page has to keep illegal

This page exists partly to reject flattering half-claims such as:

- more than one public IP exists, therefore HA is basically present
- Traefik has the route, therefore the wrong node can preserve it
- `docker-gen-failover` wrote generated fallback config, therefore failover is
  solved
- the protected page still loads, therefore policy meaning survived handoff
- the TCP port is reachable, therefore the stateful service became resilient

Those are not harmless wording mistakes.
They are the exact path back to the ambiguity that made the user furious.

## What the current runtime materially contains

The routing stack is real.
The current worktree already ships serious ingress machinery.

### Main L7 edge and policy surfaces

The active proxy fragment
[`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
contains live edge and helper surfaces including:

- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `cloudflare-ddns`
- `docker-gen-failover`
- `whoami`
- `autokuma`
- `logrotate-traefik`

The same fragment also defines live auth and middleware surfaces such as:

- `traefik.http.middlewares.nginx-auth.forwardAuth.address:
  http://nginx-traefik-extensions:80/auth`
- `traefik.http.middlewares.tinyauth.forwardAuth.address:
  http://auth:3000/api/auth/traefik`
- `docker-gen-failover` writing dynamic material to
  `/traefik/dynamic/failover-fallbacks.yaml`

### Protected HTTP route surfaces

The root Compose surface already attaches protected-route middleware to real
services such as:

- `code-server` via `traefik.http.routers.code-server.middlewares:
  nginx-auth@file`
- `dozzle` via `traefik.http.routers.dozzle.middlewares: nginx-auth@file`
- `portainer` via `traefik.http.routers.portainer.middlewares:
  nginx-auth@file`
- additional admin and metrics surfaces in the metrics and LLM fragments

### TCP transport surfaces

The current runtime also contains real TCP routers, including:

- `mongodb`
- `redis`
- multiple `biodecompwarehouse*` routes

The root stack and core fragment both carry `traefik.tcp.*` routing labels for
those surfaces.

### Private-network and peer-reachability surfaces

The active Headscale fragment
[`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
contains both:

- `headscale-server`
- `headscale`

So the repo is not missing ambition, edge machinery, or private-mesh pressure.

## The problem is not absence of machinery

The problem is this:

> the stack contains many things that sound like distributed capability while
> still not clearly owning the truth needed to make the distributed decision
> honestly.

That is why the docs have to keep speaking like a burden ledger instead of an
overview.

## What the routing stack still does not prove

The current worktree still does not prove:

- a live tracked root placement authority such as consumed `services.yaml`
- a generic wrong-node HTTP route that already succeeds end to end
- a durable rescue path that survives preferred-backend loss
- preserved middleware and auth meaning after peer-forward handoff
- that reachable peers are equivalent to eligible peers
- that TCP routing implies stateful dignity

Those gaps are the difference between "complex ingress" and "honest failover."

## The truths the routing layer still needs

The routing problem is not one truth.
It is several truths that are often illegally blended together.

| Truth the platform needs | Why it matters | Strongest current evidence | What is still missing |
| --- | --- | --- | --- |
| Placement truth | the receiving node must know where the service lives now | repeated pressure toward `services.yaml` and peer-aware routing in intent surfaces | no clearly live tracked root placement artifact consumed by routing |
| Peer-eligibility truth | not every reachable peer is valid for every route | Headscale and private-network assumptions are live | reachability still looks weaker than proven eligibility |
| Route-durability truth | fallback only matters if it survives the failure | `docker-gen-failover` is real and writes dynamic config | generated config is weaker than backend-loss survival proof |
| Semantic-continuity truth | forwarded protected routes must still mean the same thing | real auth and middleware surfaces exist | cross-node continuity of those semantics is unproven |
| Stateful-authority truth | TCP reachability is not HA | TCP routers for MongoDB and Redis exist | authority, promotion, and write ownership remain unresolved |

Until more of those become system-owned, the operator remains the safest
distributed control plane in the stack.

## The missing routing decision object

The routing dream needs a small decision object, not just more proxy syntax.

In implementation terms, the receiving node needs to be able to answer this
shape before a route can be called peer-aware:

```yaml
routing_decision:
  hostname: whoami.example.test
  route_class: stateless-http
  entry_node: node2
  service_locality: remote
  placement_truth:
    source: services.yaml | osvc | nomad | generated-runtime-state | unknown
    service_owner: node1
    freshness: unproven
  peer_eligibility:
    peer: node1
    transport: headscale | public | lan | unknown
    health: unproven
    policy_converged: unproven
  handoff:
    preserves_auth: unproven
    preserves_middleware: unproven
    preserves_headers: unproven
    preserves_service_identity: unproven
  backend_loss:
    tested: false
    surviving_backend: unproven
  explanation_artifact: unproven
```

That object does not have to be literal YAML in v1.
But some equivalent record has to exist if the stack is going to stop making
the operator the hidden routing algorithm.

Right now, the current runtime supplies many ingredients for that object.
It does not yet supply the object itself.

The absence is not theoretical.
A repository search currently finds no tracked root `services.yaml` placement
authority in the priority Compose runtime; the only `services`-named YAML file
found is under the Garden/k8s exploration area, which is not the root Compose
contract.

## The route classes must stay separate

The repo becomes dishonest as soon as these lanes are narrated like one shared
success story.

### 1. Stateless HTTP

This is the lane most likely to earn a real early win.

Candidate proof surfaces already exist:

- `whoami`
- `wishlist`
- `mkdocs`
- `homepage`
- `searxng`

For this lane, wrong-node success would mean:

1. a healthy public node receives the request
2. that node knows whether the service is local or remote
3. if remote, it can choose an eligible peer from shared current truth
4. the route still means the same thing after forwarding
5. if failover is claimed, the rescue path survives preferred-backend loss

What does not prove this lane:

- Cloudflare can hit more than one record
- Traefik is healthy
- a local request returns `200`
- `docker-gen-failover` generated something fallback-shaped

This lane should be won narrowly or not claimed at all.

### 2. Protected HTTP

This lane is stricter because the route has policy meaning, not just transport
meaning.

The runtime already shows real protected-route intent through:

- `nginx-auth@file`
- TinyAuth
- forward-auth wiring
- protected admin and metrics surfaces

Protected wrong-node success would have to preserve:

- auth challenge behavior
- middleware order
- trust-boundary assumptions
- header behavior
- the visible semantics of the route

This is why "the page still loads" is too weak.
A forwarded protected route that answers but no longer behaves like the same
protected service is not a successful handoff.

### 3. Raw TCP

This lane already exists through Traefik TCP routers for:

- `mongodb`
- `redis`
- `biodecompwarehouse*`

That proves transport exposure.
It does not prove:

- safe peer-aware failover
- state-authority transfer
- correct client semantics after rerouting
- dignity under node loss

TCP routing is an execution tool.
It is not a substitute for authority.

### 4. Stateful routes and state-bearing services

This is the lane most likely to be lied about if the docs relax.

The live stack includes real state-bearing services such as:

- `mongodb`
- `redis`
- `nuq-postgres`
- `litellm-postgres`
- `rabbitmq`
- `qdrant`
- `headscale-server`

Headscale is the cleanest example of why this lane must stay strict.
Its server config still uses SQLite at `/var/lib/headscale/db.sqlite`.

That means a working public route to Headscale is not the same thing as
multi-node authority.
The route can be plural while the authority is still singular.

So the honest reading stays:

> if write authority, promotion rules, recovery order, and data truth still
> need a human narrator, the service is still singular in the way that matters.

## The one event this entire repo keeps orbiting

When the user says several Docker nodes should behave less stupidly, the
critical event is not "two nodes exist."

It is this:

1. Cloudflare resolves a hostname to a healthy public node
2. the request lands on a node that does not host the target service locally
3. that node must decide what to do
4. the decision must be correct for that route class
5. the explanation must exist outside one operator's head

That is the routing problem.

If a tool, helper, or page does not improve that event, it is probably solving
something adjacent instead.

## Why Cloudflare is necessary but insufficient

Cloudflare is part of the anti-SPOF story because it helps prevent one sacred
public node from becoming the only entrypoint.

But Cloudflare only buys first-hop plurality.
It does not answer:

- what if the chosen node does not host the service?
- how does that node know where the service lives now?
- how does it know whether forwarding is safe?
- how does it keep auth and middleware meaning intact?
- what survives when the preferred backend is gone?

Plural DNS is real.
It is still smaller than preserved requests.

## Why Traefik is necessary but insufficient

Traefik is central to the live runtime shape.
The stack already uses it for:

- HTTP routers
- middleware attachment
- TCP routers
- protected admin surfaces
- redirects and service exposure

That matters.

But Traefik is an execution plane.
It routes based on what it knows.
The missing question remains:

> who supplied the distributed truth needed for the right cross-node decision?

If the answer is still one operator or a half-manual helper path, the system
still has a real hidden human SPOF.

## Why `docker-gen-failover` still cannot claim victory

`docker-gen-failover` is relevant and worth preserving in the docs.
It is one of the most concrete live attempts at dynamic fallback in the stack.

Its presence proves:

- failover is a live authoring concern
- the edge is not purely static
- the repo is trying to externalize fallback behavior

Its presence does not prove:

- that the rescue route survives the preferred backend failing
- that the receiving node knows when the generated route is trustworthy
- that protected-route semantics survive the handoff
- that the helper escaped the shadow-control-plane trap

The repo's own research pressure already records the harsher warning:
the helper can still look like failover and then evaporate during the exact
failure event that made fallback matter.

That makes it a perfect example of the user's complaint.

## The proof packet this page actually wants

This repo needs fewer broad routing claims and more narrow proof packets.

The minimal packet looks like this:

1. name one route
2. name the node that receives the request first
3. show that the node does not host the target locally
4. show what shared truth told it where to go
5. show why the chosen peer was eligible rather than merely reachable
6. show that the forwarded route still meant the same thing
7. if failover is claimed, show the route after the preferred backend died
8. show the evidence without private operator narration

That packet is intentionally small because the repo does not need another large
architecture story nearly as much as it needs one honest route-level success.

The first acceptable packet should probably use a boring stateless HTTP route,
not a heroic stateful service.

Good early candidates:

| Candidate | Why it is useful | Why it is still not enough |
| --- | --- | --- |
| `whoami` | Minimal route semantics; good wrong-node smoke target. | Does not prove auth continuity or stateful behavior. |
| `mkdocs` | Real docs service; still low-state. | Does not prove protected admin behavior. |
| `wishlist` | More app-like than `whoami`, still HTTP. | May introduce app assumptions before the routing layer is isolated. |

Bad first candidates:

| Candidate | Why not first |
| --- | --- |
| `code-server` / `portainer` / `dozzle` | Protected-route semantics make the proof stricter; use after a simple route packet exists. |
| `redis` / `mongodb` | TCP and stateful authority questions will pollute the first routing proof. |
| `headscale` | Control-plane/state authority makes it a poor first proof even though it matters strategically. |

## What the docs are allowed to say today

The docs are allowed to say:

- the repo has a serious Compose-first ingress stack
- the any-node-entry dream is explicit and coherent
- the routing problem has already been decomposed correctly
- wrong-node dignity is the real benchmark
- stateless HTTP, protected HTTP, TCP, and stateful lanes must stay separate

The docs are not allowed to say:

- wrong-node behavior is basically solved
- fallback is mostly handled now
- Cloudflare plus Traefik removed the hidden operator SPOF
- TCP routing implies stateful resilience
- the platform already behaves like one cloud in the user's intended sense

## The honest bottom line

The current stack contains real edge machinery, real anti-SPOF intent, and
real peer-aware pressure.

It does not yet prove that the most important routing truths have moved out of
private human custody.

So the honest routing sentence today is:

> `bolabaden-infra` already has serious multi-node ingress machinery, but the
> decisive truths for wrong-node entry, protected-route continuity, backend-loss
> survival, and stateful dignity are still only partially system-owned.
