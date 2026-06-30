# Ingress and Failover Evidence

This page is the hard boundary around the repo's ingress story.

It exists because ingress is where infrastructure documentation most easily
becomes emotionally misleading.

The stack can have:

- Cloudflare DNS
- several public nodes
- Traefik
- healthchecks
- wildcard routes
- fallback-sounding helpers

and still fail the real question the user keeps asking:

> if traffic lands on the wrong surviving node, can that node still preserve
> the meaning of the request instead of merely answering somehow?

That is the actual ingress question in `bolabaden-infra`.

This page is not allowed to downgrade that question into:

- do multiple nodes exist?
- can traffic arrive at more than one IP?
- is Traefik running?
- is there some kind of fallback helper?

Those are ingredients.
The user is asking for preserved service meaning under wrong-node entry and
backend loss.

## What this page is and is not allowed to prove

This page is authoritative about:

1. the repo's any-node-entry and peer-forward dream is explicit
2. the current root runtime already carries a serious ingress surface
3. HTTP and TCP both already exist at the edge, but they are different proof
   classes
4. service discovery and current placement truth are still the real sticking
   points
5. the repo already knows that first-hop plurality is not the same thing as
   preserved requests
6. the current fallback path has named trust gaps
7. ingress documentation here must include auth, middleware, locality, and
   backend identity, not just path routing

This page is not authoritative about:

- that multi-node ingress is done
- that DNS plurality equals end-to-end failover
- that a healthy proxy equals wrong-node success
- that a fallback helper equals trustworthy route persistence
- that an HTTP narrative automatically applies to raw TCP services

This is an evidence boundary page, not a proxy feature catalog.

## Strongest honest current answer

The repo already has a serious ingress surface.

It does **not** yet prove a general, trustworthy wrong-node
request-preservation surface.

What is true today:

- the root runtime has real public and private network segmentation
- Traefik is a live control surface
- auth and middleware are already part of edge correctness, not optional
  decoration
- Cloudflare plurality is part of the intended first-hop anti-SPOF story
- the docs and archive are explicit that service discovery is still the
  missing middle layer
- planned failover glue has known trust problems

What is still not honestly proved:

- generic peer-forward success for arbitrary services
- durable route persistence after local backend loss
- middleware and auth continuity during peer-forward fallback
- a tracked root runtime source of current placement truth such as live
  `services.yaml` consumption
- safe equivalence between "reachable peer" and "eligible fallback target"
- raw TCP failover correctness just because L7 ideas exist

That gap is the center of the ingress story, not a minor implementation note.

## The dream this page has to protect

The architecture dream is explicit in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

- Compose-first
- multi-node
- anti-SPOF pressure
- no immediate collapse into Swarm, Kubernetes, or another heavyweight control
  plane
- local-first serving
- peer-forward fallback when the receiving node does not host the target
  service locally

The target request contract is stated there directly:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That contract is the standard this page uses.

This page is therefore not a reverse-proxy feature checklist.
It is a proof filter for one very specific dream:

> ordinary Docker nodes should behave less stupidly when requests land on the
> wrong box, the preferred backend dies, or the stack is under anti-SPOF
> pressure.

## Evidence hierarchy for ingress claims

Use this order every time.

| Claim type | Highest authority | Why it outranks others | It still does not prove |
| --- | --- | --- | --- |
| What the repo wants ingress to become | `.github/copilot-instructions.md`, `README.md` | these define the architecture dream and honesty wall | that the root runtime already does it |
| What the root stack actually ships today | `docker-compose.yml`, `compose/docker-compose.*.yml` | these are the priority implementation surface | distributed correctness under failure |
| What the repo already knows is missing | `docs/INFRASTRUCTURE_MASTER_PLAN.md`, ingress research pages | these name the actual gaps | that the repair path is already live |
| Why ordinary answers keep failing the user | source archive conversations | these restore the real complaint | runtime proof |

If a paragraph blends more than one row, the docs should say which row is doing
which work.

## What the archive proves about the actual missing layer

The archive is unusually explicit about the user's real question.

### 1. The repo is not looking for "more Docker"

`docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`
opens with the exact shape of the problem:

- "Can I just unify a bunch of Docker hosts with some clever networking +
  load balancers, without introducing a cluster manager?"
- requests need "some dynamic mechanism to update routing when containers
  move/come/go"
- if placement is manual and Cloudflare already handles entry, then "your real
  challenge is service discovery/routing"

That matters because it kills the most common bad summary:

> the user wants multi-node Docker

No.
The user already has multi-node Docker.
The missing thing is shared current-state truth that lets the receiving node
behave correctly when locality is absent.

### 2. The user states the desired wrong-node behavior directly

The same archive file includes the clearest single sentence in the whole
project:

- "the nodes are setup in a way where if any service was requested on node1
  but it doesn't exist on node1 ... it'll l7/l4 ... to the other nodes i own.
  So all I really need is service discovery to be loadbalanced/unified."

That sentence is more informative than a generic "high availability" label.
It tells you exactly what the operator wants the receiving node to know.

This is why ingress writing in this repo has to keep asking:

- what runs where right now?
- how does the receiving node know that?
- what proves the receiving node is choosing a healthy, semantically valid
  target?

### 3. Familiar load balancer names are not the missing option

`load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
shows the user asking for existing projects that really approximate the kind of
behavior they want, and explicitly rejecting answers that stop at familiar
reverse-proxy categories.

That is critical context.
It means docs in this repo should never sound impressed merely because Traefik,
NGINX, or HAProxy appear in the stack.

The presence of serious tools is not the same thing as the removal of the
hidden burden.

### 4. The user is resisting heavyweight orchestrators for a reason

`distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
captures the other side of the pressure:

- there is no clean drop-in peer-equal Compose scaler
- glue is likely unavoidable
- the user still does not want to build a whole orchestration framework

That does not prove the current stack works.
It proves why the repo keeps circling a missing middle layer instead of just
accepting Swarm or Kubernetes as the answer.

### 5. Emotional pressure is part of the retrieval surface

`docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`
contains the line:

- "Docker feels gaslighting"

That is not background color.
It explains why the docs must not reward the same "many partial tools exist,
therefore you have options" pattern the user is explicitly rejecting.

The emotional pressure is evidence about the acceptance bar, not evidence that
the runtime already crossed it.

## What the live root runtime concretely proves

The priority implementation is still the merged root graph centered on:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- [`compose/docker-compose.docs.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.docs.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- [`compose/docker-compose.metrics.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.metrics.yml)

### 1. The ingress surface is real, not rhetorical

The root runtime defines and uses:

- `publicnet`
- `backend`
- `warp-nat-net`

This proves the root stack already distinguishes public entry, internal service
traffic, and specialized routing domains.

It does **not** prove that the receiving node knows which peer currently hosts
which service.

### 2. Traefik is already one of the stack's real control planes

The edge fragment contains live components such as:

- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `cloudflare-ddns`
- `docker-gen-failover`

This proves:

- request correctness already depends on policy and middleware surfaces
- ingress is not just port exposure
- the edge already has multiple cooperating control layers

It does **not** prove:

- that those layers preserve behavior under peer-forward fallback
- that route state survives local backend disappearance
- that auth continuity remains the same when locality changes

### 3. Protected HTTP routes are already a real stress class

The root runtime clearly includes routes such as Dozzle that already depend on
policy continuity, not just reachability.

For example, the root Compose surface shows:

- a `dozzle` service
- Traefik route material for that service
- middleware attachment through `nginx-auth@file`

That matters because it makes the wrong-node problem harder and more concrete.
The user is not asking whether an unprotected dummy upstream can be forwarded.
They are asking whether the same protected service still exists after handoff.

### 4. Raw TCP routes are already part of the same runtime

The root runtime also clearly contains TCP exposure for services such as:

- MongoDB
- Redis

That proves the ingress story cannot stay purely HTTP-centric.
It also proves the docs must stay strict:

- TCP route existence is real
- TCP failover truth is still unproven
- stateful correctness is stricter than TCP reachability

### 5. Headscale makes the mesh assumption real, not hypothetical

The current runtime already includes:

- `headscale-server`
- `headscale`

That means the repo already depends on a real private-mesh worldview.
But the planning layer also records that Headscale is single-node today.

So the correct summary is:

- mesh assumptions are real
- mesh control-plane HA is not yet proved

### 6. Cloudflare plurality is present, but the repo itself warns that it is
not enough

The stack includes `favonia/cloudflare-ddns`, but
[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
explicitly says current Cloudflare DDNS presence is not the same thing as full
multi-node request failover.

Without the planning layer, someone could say:

- Cloudflare DDNS exists
- therefore any-node ingress is basically handled

The plan explicitly rejects that inflation.

### 7. The current fallback helper is not trustworthy enough to narrate as solved

The same plan calls out a very specific gap:

- `docker-gen-failover` removes Traefik routes when containers stop
- route persistence under local backend failure is still missing
- automated service failover between nodes is still missing

That is one of the most important live-reading constraints in the repo.

The stack contains a fallback-looking subsystem.
The strongest repo-native planning doc says that subsystem currently defeats
the very property it was meant to provide.

So the honest summary is:

- fallback glue exists
- fallback trust does not

### 8. The root runtime still lacks tracked shared placement truth

The intent surfaces repeatedly converge on lightweight current-state truth like
`services.yaml` or equivalent shared discovery state.

The README is explicit:

- the repo keeps circling a lightweight `services.yaml` current-state registry
- the tracked root implementation does not currently ship a live root
  `services.yaml`

That matters because the missing thing is not a better load balancer.
It is shared current-state knowledge.

## The exact claim boundaries this evidence supports

This evidence stack supports all of these sentences:

- the repo's any-node-entry and peer-forward dream is explicit
- the current runtime already has a serious ingress surface
- protected HTTP routes are a real present-tense concern
- TCP routes are real present-tense concerns
- service discovery and placement truth are still the real sticking points
- the repo already knows first-hop plurality is much weaker than preserved
  request meaning
- the current fallback helper is directionally relevant but not yet trustworthy

This evidence stack does **not** support these sentences:

- wrong-node requests already succeed generically
- durable route persistence is solved
- protected-route policy continuity is proven under peer handoff
- Cloudflare plus Traefik already gives end-to-end failover
- a reachable peer is automatically an eligible peer
- TCP or stateful failover is solved because HTTP ideas exist

Those stronger claims need route- or topology-specific drills.

## Why this page stays harder on language than normal proxy docs

Ordinary proxy documentation is often happy to stop at:

- route configured
- backend healthy
- wildcard DNS working
- helper generated config

That language is too weak for this repo because the user's actual complaint is
that respectable layers keep getting stacked without the hidden operator burden
really moving.

So this page has to stay harder:

- first hop is not enough
- local proxy health is not enough
- fallback-looking config is not enough
- service discovery talked about in prose is not enough
- a live mesh is not enough

Only shared truth plus route-specific drills can close the gap.

## Bottom line

The ingress story in `bolabaden-infra` is already serious enough to deserve
careful documentation.
It is still too incomplete to deserve casual failover language.

The live runtime proves a strong edge surface.
The archive proves exactly why that is still not the same thing as a solved
wrong-node request-preservation surface.

That distinction is the most important thing this page is protecting.
