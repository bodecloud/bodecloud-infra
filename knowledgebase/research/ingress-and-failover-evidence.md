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

## Strongest honest current answer

The repo already has a serious ingress surface.

It does **not** yet prove a general, trustworthy wrong-node request-preservation
surface.

What is true today:

- the root runtime has real public and private network segmentation
- Traefik is a live control surface
- auth and middleware are already part of edge correctness, not optional
  decoration
- Cloudflare plurality is part of the intended first-hop anti-SPOF story
- the docs and archive are extremely explicit that service discovery is still
  the missing middle layer
- planned failover glue has known trust problems

What is still not honestly proved:

- generic peer-forward success for arbitrary services
- durable route persistence after local backend loss
- middleware and auth continuity during peer-forward fallback
- a tracked root runtime source of current placement truth like live
  `services.yaml` consumption
- safe equivalence between "reachable peer" and "eligible fallback target"
- raw TCP failover correctness just because L7 ideas exist

That gap is the center of the ingress story, not an implementation detail.

## What this page is and is not allowed to prove

This page is allowed to prove:

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

This page is not allowed to prove:

- that multi-node ingress is "done"
- that DNS plurality equals end-to-end failover
- that a healthy proxy equals wrong-node success
- that a fallback helper equals trustworthy route persistence
- that an HTTP narrative automatically applies to raw TCP services

## Evidence hierarchy for ingress claims

Use this order every time.

| Claim type | Highest authority | Why it outranks others | It still does not prove |
| --- | --- | --- | --- |
| What the repo wants ingress to become | `.github/copilot-instructions.md`, `README.md` | these define the architecture dream and honesty wall | that the root runtime already does it |
| What the root stack actually ships today | `docker-compose.yml`, `compose/docker-compose.*.yml` | these are the priority implementation surface | distributed correctness under failure |
| What the repo already knows is missing | `docs/osvc_ingress_ha.md`, `docs/INFRASTRUCTURE_MASTER_PLAN.md` | these name the actual gaps | that the repair path is already live |
| Why the ordinary answers keep failing the user | source archive conversations | these restore the real complaint | runtime proof |

If a paragraph blends more than one row, say which row is doing which work.

## What the archive proves about the real missing layer

The archive is unusually explicit about the user's actual question.

### 1. The repo is not looking for "more Docker"

[`docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`](../source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md)
opens with the exact shape of the problem:

- "Can I just unify a bunch of Docker hosts with some clever networking +
  load balancers, without introducing a cluster manager?"
- requests need "some dynamic mechanism to update routing when containers
  move/come/go"
- if placement is manual and Cloudflare already handles entry, then "your real
  challenge is service discovery/routing"

That matters because it kills the most common bad summary:

> the user wants multi-node Docker

No. The user already has multi-node Docker.
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

### 3. "Existing load balancers" are not the user's real missing option

[`load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`](../source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md)
shows the user asking for "the closest existing projects" that replace
Cloudflare's paid failover behavior.

The same conversation immediately concedes that familiar names like Traefik,
NGINX, and HAProxy may not fully satisfy the need.

That is critical context.
It means docs in this repo should never sound impressed merely because those
names appear in the stack.

### 4. Traefik label failover was explicitly tried and explicitly failed

[`traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md`](../source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md)
contains a very concrete dead end:

- attempted label:
  `traefik.http.routers.whoami.failover.service: whoami-servers@file`
- resulting runtime error:
  `field not found, node: failover`

That matters because it stops a common lie:

> Traefik already has failover, we just need to wire it up better

The archive shows that the naive "just use Traefik labels" path already hit a
real wall.

### 5. The user is resisting heavyweight orchestrators for a reason

[`distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`](../source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md)
captures the other side of the pressure:

- K3s still uses leader election
- Docker Swarm still has manager nodes
- "you will probably need to build glue, but not a full orchestration
  framework"

That does not prove the current stack works.
It proves why the repo keeps circling a missing middle layer instead of just
accepting Swarm or Kubernetes as the answer.

### 6. The emotional pressure is part of the requirements

[`docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`](../source-archive/chatgpt-exports/conversations/docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md)
includes the line:

- "Docker feels gaslighting"

That is not background color.
It explains why the docs must not reward the same "there are many partial tools
therefore you have options" pattern that the user is explicitly rejecting.

## What the live root runtime concretely proves

The priority implementation is still the merged root graph centered on:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- [`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
- [`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)

### 1. The ingress surface is real, not rhetorical

[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
defines and uses:

- `publicnet`
- `backend`
- `warp-nat-net`

This proves the root stack already distinguishes public entry, internal service
traffic, and specialized routing domains.

It does **not** prove that the receiving node knows which peer currently hosts
which service.

### 2. Traefik is already one of the stack's real control planes

[`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
contains live edge components such as:

- `traefik`
- `tinyauth`
- `crowdsec`
- `cloudflare-ddns`
- `docker-gen-failover`

This proves:

- request correctness already depends on policy and middleware surfaces
- ingress is not just port exposure
- the edge stack already has multiple cooperating control layers

It does **not** prove:

- that these layers preserve behavior under peer-forward fallback
- that route state survives local backend disappearance
- that auth continuity remains the same when locality changes

### 3. Cloudflare plurality is present, but the repo itself warns that it is
not enough

The stack includes `favonia/cloudflare-ddns`, but
[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
explicitly says current Cloudflare DDNS presence is not the same thing as full
multi-node request failover.

This is one of the best examples of why the repo needs source-weighted docs.

Without the planning layer, someone could easily say:

- Cloudflare DDNS exists
- therefore any-node ingress is basically handled

The plan explicitly rejects that inflation.

### 4. The current fallback helper is not trustworthy enough to narrate as
solved

The same plan calls out a very specific gap:

- `docker-gen-failover` removes Traefik routes when containers stop
- route persistence under local backend failure is still missing
- automated service failover between nodes is still missing

That is one of the most important live-reading constraints in the whole repo.

The stack contains a fallback-looking subsystem.
The strongest repo-native planning doc says that subsystem currently defeats
the very property it was meant to provide.

So the honest summary is:

- fallback glue exists
- fallback trust does not

### 5. The root runtime still lacks tracked shared placement truth

The intent surfaces repeatedly converge on lightweight current-state truth like
`services.yaml` or equivalent shared discovery state.

The live root runtime does not currently prove that such a registry is both:

- tracked in the priority implementation
- and actually consumed as the placement truth for peer-forward routing

This is the cleanest single explanation for why the dream is still not proved.

The missing layer is not "more reverse proxy."
It is shared current-state knowledge.

## HTTP versus TCP: do not flatten them

The repo already carries both HTTP and raw TCP edge concerns.

That does **not** mean one solution class covers both.

[`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
is explicit that:

- HTTP(S) can route by hostname and already has meaningful ingress work
- plain TCP like Redis cannot be treated the same way without stricter
  constraints such as TLS/SNI or topology-specific forwarding

So any ingress sentence that sounds like "the routing layer solves failover"
must say which service class it means:

- stateless HTTP
- raw TCP
- or state-bearing traffic that happens to cross the edge

## What the current docs must force every reader to ask

Before accepting any ingress claim, ask:

1. what node received the request?
2. was the target service local or remote to that node?
3. where did the receiving node learn the current placement truth?
4. what makes that truth fresher than operator memory?
5. what survives if the preferred local backend disappears?
6. what stays true about auth and middleware after peer forwarding?
7. is this claim about HTTP path preservation or a raw TCP workaround?

If the docs do not help answer those questions, they are still operating at the
ingredient-list level the user is trying to escape.

## Bottom line

The current worktree proves that `bolabaden-infra` is already serious about
ingress architecture.

It does **not** yet prove that ingress has crossed the line into:

- generic wrong-node success
- route persistence under backend loss
- shared current-state service discovery
- semantically preserved peer-forward behavior

That missing middle layer is the real ingress problem.

The user's complaint is not "there are too few reverse proxies."
It is:

> the ecosystem keeps offering fragments while leaving service discovery,
> request preservation, and honest anti-SPOF behavior as my problem.

This page should be read as the repo's refusal to pretend that gap is already
closed.
