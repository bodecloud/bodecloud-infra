# Ingress and Failover Evidence

This page is the hard boundary around the repo's ingress story.

It exists because ingress is where infrastructure documentation most easily
turns disrespectful without sounding careless.

The ecosystem is very good at producing a stack that looks reassuring:

- Cloudflare DNS
- several public nodes
- Traefik
- healthchecks
- wildcard routes
- helper containers with words like `failover`

and still leaves the user with the same humiliating private question:

> if traffic lands on the wrong healthy node, can that node still preserve the
> meaning of the request without me privately completing the topology?

That is the ingress question in `bolabaden-infra`.

This page is not allowed to downgrade that question into softer adjacent ones:

- do multiple nodes exist?
- can traffic hit more than one IP?
- is Traefik running?
- is there some helper that sounds dynamic?

Those are ingredients.
The user is not starved for ingress ingredients.
The user is starved for a system that stops cashing those ingredients back out
as human memory the moment locality disappears.

That is why this page has to stay rude.
If it becomes smoother than the failure it is describing, it has already
started lying.

## The strongest honest current answer

The repo already has a serious ingress surface.

It does not yet prove a general, trustworthy wrong-node
request-preservation surface.

What is true right now:

- the architecture dream is explicit
- the priority runtime really does ship an edge stack
- protected HTTP routes are already a real concern
- raw TCP routes are already a real concern
- first-hop plurality is already part of the design
- the repo already knows service discovery and current placement truth are the
  real missing middle
- the current fallback helper is directionally relevant but not trustworthy

What is still not honestly proved:

- generic peer-forward success for arbitrary services
- durable route persistence after local backend loss
- middleware and auth continuity during peer fallback
- a live tracked root placement-truth surface such as actual
  `services.yaml` consumption
- safe equivalence between `reachable peer` and `eligible peer`
- raw TCP failover correctness just because L7 ideas exist

That gap is not a footnote.
That gap is the ingress story.

The user's complaint lives exactly there:
everything before it can look adult;
everything after it still too often sounds like:

> trust me, I know where it really runs.

## What this page is and is not allowed to prove

This page is allowed to prove:

1. the repo's any-node-entry and peer-forward dream is explicit
2. the current root runtime already carries a serious ingress surface
3. HTTP and TCP are both real present-tense edge concerns, but they are
   different proof classes
4. service discovery and current placement truth are still the real sticking
   points
5. the repo already knows first-hop plurality is weaker than preserved request
   meaning
6. the current fallback path has named trust gaps
7. ingress claims here must include locality, backend identity, auth, and
   middleware, not merely path routing

This page is not allowed to prove:

- that multi-node ingress is done
- that DNS plurality equals end-to-end failover
- that a healthy proxy equals wrong-node success
- that a fallback helper equals trustworthy route persistence
- that an HTTP narrative automatically applies to raw TCP services

This is an evidence boundary page, not a proxy feature catalog.

## Retrieval contract for this page

This page has to keep four source classes separate on purpose.

### Class 1: the dream

Primary anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

This class is allowed to prove:

- the repo wants any healthy node to receive the first request
- the repo wants local-first then peer-forward behavior
- the repo wants Compose-first multi-node anti-SPOF pressure without immediate
  Swarm/Kubernetes capture

This class is not allowed to prove:

- that the root runtime already does that

### Class 2: the live root runtime

Primary anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active files under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)

This class is allowed to prove:

- what the priority implementation currently ships
- which ingress components are real
- which networks, routes, middleware surfaces, and TCP/HTTP exposures exist

This class is not allowed to prove:

- distributed correctness under failure just because the edge stack is rich

### Class 3: repo-native gap naming

Primary anchors:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- architecture evidence pages in this knowledgebase

This class is allowed to prove:

- what the repo itself already admits is still missing
- where current helper layers are not trustworthy enough yet

This class is not allowed to prove:

- that the repair path is already live

### Class 4: archive pressure

Primary anchors:

- [`docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md)
- [`load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md)
- [`distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md)
- [`docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/chatgpt-exports/conversations/docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md)

This class is allowed to prove:

- what the user is actually rebelling against
- why ordinary ingress language keeps failing the user
- which private sentence the docs must keep alive until the runtime kills it

This class is not allowed to prove:

- current runtime behavior

If a paragraph blends these classes, it should say which class is doing which
work.
Otherwise the page will slowly start using archive force or dream force to fake
runtime force.

## The dream this page has to protect

The strongest intent wording lives in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md),
which states the target operating contract directly:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That contract is the standard this page uses.

So this page is not a reverse-proxy checklist.
It is a proof filter for one very specific dream:

> ordinary Docker nodes should stop behaving insultingly dumb when requests
> land on the wrong box, the preferred backend dies, or anti-SPOF pressure is
> finally supposed to mean something.

The user is not chasing prestige architecture.
The user is chasing a stack that stops making multiple nodes feel decorative on
the bad day.

## What the archive proves about the real missing layer

The archive is unusually explicit here.

### 1. The user is not asking for more Docker

`docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`
opens with the exact shape of the wound:

- can several Docker hosts be unified without immediately introducing a cluster
  manager?
- if placement is manual and Cloudflare already handles entry, then the real
  challenge is service discovery and routing

So the bad summary is:

> the user wants multi-node Docker.

No.
The user already has multi-node Docker.
The missing thing is shared current-state truth that lets the receiving node
behave correctly when locality is absent.

### 2. The user states the desired wrong-node behavior directly

The same archive thread contains one of the most revealing sentences in the
project:

- if a service is requested on node1 and it does not exist on node1, the
  system should hand it off to another owned node
- therefore the real missing layer is service discovery that feels unified

That sentence is more useful than a dozen generic HA labels.
It says exactly what the receiving node has to know.

So every ingress explanation here should keep asking:

- what runs where right now?
- how does the receiving node know that?
- what proves the chosen target is healthy and semantically valid?

### 2.5. The same source also shows the exact overclaim to avoid

The `docker-multi-node-without-swarm` thread is valuable because it does not
only contain the desired flow.
It also contains a generated draft that turns the target architecture into
premature certainty.

That draft treats the combination of:

- multiple Cloudflare A records
- per-node L4/L7 forwarding
- redundant proxy processes
- presumed stateful replication

as if it were enough to say the system has no SPOFs.

That is exactly the documentation failure this page exists to prevent.

The honest extraction is narrower:

- DNS and DDNS can reduce sacred public entrypoint pressure
- per-node proxying can become the mechanism for wrong-node request
  preservation
- a registry or discovery surface can become the missing request-time truth
  source

The dishonest upgrade is:

- therefore service-level failover is solved
- therefore TCP forwarding proves stateful continuity
- therefore "zero SPOF" is already earned

This source has to be read as both architecture pressure and cautionary
evidence.
It proves why the user wants the flow.
It also proves how easy it is for documentation to flatter the flow before the
runtime owns it.

That means this page should never say "Cloudflare plus L4/L7 forwarding
ensures no SPOFs."
It can only say:

> Cloudflare plus per-node forwarding describes the intended shape of the
> missing middle. The proof still has to show placement truth, peer
> eligibility, route durability, policy continuity, and stateful authority
> separately.

### 3. Familiar reverse-proxy names are not the missing option

`load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
shows the user explicitly asking for existing projects that really approximate
the behavior they want, and rejecting answers that stop at familiar reverse
proxy categories.

That means this repo should never sound impressed merely because Traefik,
NGINX, or HAProxy appear.

The presence of serious tools is not the same thing as the removal of the
hidden burden.

### 4. The user resists heavyweight capture for a reason

`distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
captures the other side of the pressure:

- there is no clean drop-in peer-equal Compose scaler
- glue is likely unavoidable
- the user still does not want to build a whole orchestration framework

That does not prove the current stack works.
It proves why the repo keeps circling a missing middle instead of simply
accepting Swarm or Kubernetes as the emotional end of the conversation.

### 5. Emotional pressure is part of the evidence surface

`docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`
contains the line:

- Docker feels gaslighting

That is not decorative color.
It explains why the docs must not reward the same pattern of:

- many partial tools exist
- therefore the user already has options

The emotional pressure is evidence about the acceptance bar.
It is not evidence that the runtime has crossed it.

## What the live root runtime concretely proves

The priority implementation is still the merged root graph centered on:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- [`compose/docker-compose.docs.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.docs.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- [`compose/docker-compose.metrics.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.metrics.yml)

### 1. The ingress surface is real

The root runtime defines and uses:

- `publicnet`
- `backend`
- `warp-nat-net`

That proves the stack already distinguishes public entry, internal service
traffic, and specialized routing domains.

It does not prove that the receiving node knows which peer currently hosts
which service.

That sentence matters more than most healthy Compose output.
It is the line between:

- the edge exists
- the edge can preserve meaning under wrong-node entry

### 2. Traefik is already one of the stack's real control planes

The edge fragment contains live components such as:

- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `cloudflare-ddns`
- `docker-gen-failover`

That proves:

- request correctness already depends on policy and middleware surfaces
- ingress is not just port exposure
- the edge already has multiple cooperating control layers

It does not prove:

- those layers preserve behavior together under peer-forward fallback
- route state survives local backend disappearance
- auth continuity remains intact when locality changes

### 3. Protected HTTP routes are already a real stress class

The root Compose surface includes protected routes such as Dozzle, code-server,
and Portainer, where correctness already depends on policy continuity rather
than mere reachability.

That matters because the user is not asking whether an unprotected dummy
upstream can be forwarded.
They are asking whether the same protected service still exists after handoff.

### 4. Raw TCP routes are already part of the same runtime

The root runtime also contains TCP exposure for services such as:

- MongoDB
- Redis

That proves the ingress story cannot stay purely HTTP-centric.

It also forces a stricter language wall:

- TCP route existence is real
- TCP failover truth is still unproven
- stateful correctness is stricter than TCP reachability

### 5. Headscale makes the mesh assumption real

The current runtime includes:

- `headscale-server`
- `headscale`

So private-mesh assumptions are not hypothetical.
But the planning layer also records that Headscale is still single-node today.

The correct summary is therefore:

- mesh assumptions are real
- mesh-control-plane HA is not yet proved

### 6. Cloudflare plurality is present, but the repo says that is not enough

The stack includes Cloudflare DDNS, but the planning and instruction surfaces
explicitly reject the idea that Cloudflare presence alone equals full multi-node
request failover.

Without those honesty surfaces, someone could say:

- Cloudflare DDNS exists
- therefore any-node ingress is basically handled

The repo itself rejects that inflation.

### 7. The current fallback helper is not trustworthy enough to narrate as solved

The same evidence stack keeps calling out a specific gap:

- `docker-gen-failover` is real
- it is directionally relevant
- route persistence under local backend failure is still missing
- automated service failover between nodes is still missing

So the honest summary is:

- fallback glue exists
- fallback trust does not

### 8. The root runtime still lacks tracked shared placement truth

The intent surfaces repeatedly converge on lightweight current-state truth like
`services.yaml` or equivalent shared discovery state.

The repo is explicit that the tracked root implementation does not currently
ship a live root `services.yaml`.

That matters because the missing thing is not a better load balancer.
It is shared current-state knowledge.

## What still does not count as ingress evidence here

The following are too weak to count as meaningful ingress proof in this repo:

- several A or AAAA records existing at Cloudflare
- a healthy Traefik dashboard
- multiple reverse-proxy containers running
- a generated fallback file existing on disk
- a peer being pingable or reachable over Tailscale
- one service answering after manual operator nudging

These are ingredients or local signals.
They are not proof that the wrong healthy node can preserve request meaning
under pressure.

This page should stay hostile to ingredient inflation.
The repo already has enough ingredients to sound almost finished to the wrong
reader.
The job here is to make that reader uncomfortable again.

## What a real ingress proof packet would have to contain

For this repo, a serious ingress proof packet would need to show all of the
following together:

- the exact entry node that first received the request
- whether the requested service was local or remote at that moment
- the current placement truth source the receiving node consulted
- the peer-selection logic that chose the fallback target
- preserved middleware and auth behavior across the handoff
- backend health or loss conditions during the test
- operator-readable artifacts explaining why the route was valid

If any of those are missing, the packet may still help debug the system.
It does not close the user's real ingress question.

The operator-readable artifact matters especially.
The user is not only asking for successful routing.
They are asking for a system that can later explain why the route was valid
without requiring the operator to retell the topology from memory.

## Evidence packet applied to the current root stack

Applying the evidence-ledger source-custody format to ingress gives this
current packet:

| Packet field | Current answer |
| --- | --- |
| User wound | A healthy wrong node should not need the operator to remember where the service really lives. |
| Runtime anchor | `docker-compose.yml` and active `compose/` fragments contain Traefik, auth/middleware surfaces, Cloudflare DDNS, Headscale, TCP routers, and `docker-gen-failover`. |
| Intent anchor | `.github/copilot-instructions.md` states the any-node, local-first, peer-forward target contract. |
| Plan/gap anchor | Master-plan and architecture pages keep naming current placement truth and route durability as missing. |
| Archive anchor | Imported multi-node-without-swarm, load-balancer-failover, distributed-HA, and Compose-frustration threads explain why ordinary answers feel fake. |
| Legal sentence | The repo has serious ingress machinery and a precise peer-forward dream, but generic wrong-node request preservation remains unproved. |
| Illegal sentence | Cloudflare plus Traefik plus `docker-gen-failover` already gives trustworthy multi-node failover. |
| Next proof | A route-specific packet for one stateless HTTP route where a non-hosting entry node uses inspectable placement truth to choose an eligible peer and preserve route meaning. |

This packet is intentionally unsatisfying.
It is unsatisfying in the same place the runtime is still incomplete.

That is the correct outcome.
The documentation should not soothe the exact gap the implementation has not
closed.

## Exact claim boundaries this evidence supports

This evidence stack supports these sentences:

- the repo's any-node-entry and peer-forward dream is explicit
- the current runtime already has a serious ingress surface
- protected HTTP routes are a real present-tense concern
- TCP routes are a real present-tense concern
- service discovery and placement truth are still the real sticking points
- the repo already knows first-hop plurality is weaker than preserved request
  meaning
- the current fallback helper is directionally relevant but not yet
  trustworthy

This evidence stack does not support these sentences:

- wrong-node requests already succeed generically
- durable route persistence is solved
- protected-route policy continuity is proved under peer handoff
- Cloudflare plus Traefik already gives end-to-end failover
- a reachable peer is automatically an eligible peer
- TCP or stateful failover is solved because HTTP ideas exist

Those stronger claims need route-specific or topology-specific drills.

## The honest bottom line

This repo does not lack ingress nouns.
It lacks one more system-owned truth layer between:

- first-hop plurality
- and preserved request meaning

Until the runtime can show that the wrong healthy node knows where the service
actually lives, knows the target is semantically valid, preserves policy during
handoff, and leaves behind an inspectable explanation artifact, the ingress
story remains promising but incomplete.
