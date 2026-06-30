# Design Tensions and Contradictions

This page is the contradiction ledger for the priority implementation rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

The user's dream is not vague:

> several ordinary Docker nodes should start behaving like one
> request-preserving personal cloud, without immediately surrendering to
> Kubernetes, Docker Swarm, or another heavyweight orchestrator.

The worktree already contains enough real machinery that the dangerous mistake
is no longer:

> this is too small to matter.

The dangerous mistake is:

> this already sounds like a clustered platform, so the remaining gaps must be
> minor.

They are not minor.
They are the exact truths that decide whether the platform owns the bad day or
the operator still does.

## What this page is and is not allowed to prove

This page is allowed to prove:

- the main contradictions the repo is consciously carrying
- which contradictions are rooted in real current files and services
- where intent is stronger than live proof
- what specific artifact would reduce each contradiction instead of merely
  narrating it better

This page is not allowed to prove:

- which option has already won
- whether a future controller is justified
- whether the current runtime already solves the contradiction
- whether sharper prose counts as burden transfer

This is a contradiction ledger, not a closure report.

## The strongest honest current answer

The repo already escaped the toy-stack phase.
It has not escaped the private-operator-completion phase.

That means the real contradictions are no longer generic tradeoffs like
`simplicity versus power`.
They are more exact:

- Compose is still the readable truth surface, but distributed truth is not yet
  system-owned.
- Cloudflare can give plural first-hop entry, but plural DNS is not preserved
  service meaning.
- Traefik can execute real edge policy, but edge execution is not shared
  placement knowledge.
- Headscale can make peers reachable, but reachability is not peer validity.
- TCP exposure can look clean, but clean exposure is not stateful authority.

If the docs flatten those contradictions into generic architecture language,
they stop describing the user's actual complaint.

## All contradictions are downstream of one event

The repo's real request contract is still:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

Every contradiction on this page is a different way that contract can collapse
back into private operator narration.

That matters because the contradictions are not abstract tensions.
They are different versions of the same accusation:

> the platform still looks coherent only because a human is privately joining
> the missing truths.

## The contradiction packet standard

Before the docs claim a contradiction is improving, they should be able to
point to a packet containing:

- the exact contradiction being reduced
- the hidden burden that previously remained private
- the new artifact, drill, or runtime surface that narrowed it
- the stronger sentence now allowed
- the sentence that is still forbidden

Without that packet, improvement is still mostly narrative.

## The central contradiction

The whole repo compresses to one difficult sentence:

> any healthy public node should be able to receive the request, determine
> whether the service is local, preserve the request if it is not, and do so
> without a giant scheduler quietly becoming the only truthful adult in the
> room.

That sentence requires the system to own several distinct truths:

- public-entry truth
- current-placement truth
- locality truth
- peer-eligibility truth
- route-persistence truth
- policy-continuity truth
- service-class truth
- stateful-authority truth

Single-node Docker lets one human blur those together.
Wrong-node traffic forces them apart.

That is why the user keeps sounding harsher than a normal infra planner.
The wrong-node event is what makes the hidden human control plane visible.

## Contradiction matrix

| Tension | What the repo wants to keep | What the live worktree already has | What is still missing | Why the nearby answer still feels fake | What would materially reduce the contradiction |
| --- | --- | --- | --- | --- | --- |
| Compose readability vs distributed truth ownership | Root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml) plus `compose/` stays the primary authoring contract | Active multi-fragment Compose runtime, shared networks, inline configs, real ingress, real stateful services | A live shared source of current placement and eligibility truth | `Use more Compose` stops one layer early; `abandon Compose` often demands worldview surrender before proving it moves the right burden | Introduce one consumed placement-truth surface and show a receiving node using it for a real decision |
| No orchestrator by default vs someone still knowing what runs where | Manual placement remains acceptable where it still pays for itself | Repeated intent toward `services.yaml`, sync-agent ideas, helper-generated routing, OpenSVC pressure in research docs | Present-tense proof that routing or eligibility consumes a tracked shared placement surface | The docs can describe the missing middle more clearly than the runtime proves it exists | One root-level current-state registry, or equivalent, visibly consumed by routing or peer-selection logic |
| Any-node public entry vs preserved request meaning | No sacred public box; any healthy node should accept the first hop | `cloudflare-ddns`, Cloudflare-first design, Traefik edge stack | Proof that first-hop plurality survives into correct service meaning on the wrong node | Multi-record DNS feels like failover even when the receiving node still lacks service truth | Wrong-node drill where the landing node is intentionally not the service owner and still routes correctly |
| Local-first serving vs backend-loss survival | If the service is local, keep it local and legible | Traefik local routing, real routers, healthchecks, app surfaces like `whoami`, `wishlist`, and `code-server` | Proof that the rescue path survives the failure that made rescue necessary | A local `200` is comforting but says nothing about the anti-SPOF claim | Demonstrate local success, then stop the local backend and show preserved route continuity |
| Serious ingress machinery vs missing cross-node truth | Rich edge behavior without a giant cluster API | `traefik`, `tinyauth`, `nginx-traefik-extensions`, `crowdsec`, `docker-gen-failover`, file-provider surfaces | Shared placement, shared peer eligibility, protected-route parity after handoff | A sophisticated proxy stack can flatter the reader into thinking the hard part is mostly done | Compare local and peer-forwarded behavior for one protected HTTP route |
| Mesh reachability vs peer validity | Private mesh should make peer-forward feasible | Live Headscale fragment with reachable private-node assumptions | Proof that reachability plus identity becomes `this peer is valid for this request right now` | Connectivity products often become fake service discovery by social implication alone | Peer-selection artifact that names allowed backends and the freshness rules behind them |
| HTTP routing logic vs TCP and stateful truth | Stateless HTTP should mature first without lying about harsher lanes | Traefik HTTP and TCP routers already exist for services such as `mongodb` and `redis` | Separate proof rules for transport, authority, promotion, and reconnect semantics | Clean TCP ingress is frequently overread as real HA | Keep HTTP, TCP, and stateful lanes under separate acceptance bars and drills |
| Reachability vs stateful authority | Stateful services should be spoken about much more harshly | Root `mongodb`, root `redis`, Headscale SQLite, `nuq-postgres`, `rabbitmq`, `litellm-postgres`, `qdrant` | Replication, promotion, rediscovery, and writer-authority truth | `It still answers` is not `the write owner stopped being singular` | Per service class, define write authority, failover sequence, promotion rules, and rediscovery behavior |
| Helper growth vs hidden control-plane drift | Avoid rebuilding a giant scheduler accidentally | `docker-gen-failover`, route generators, auth helpers, NAT routing, metrics, file-provider logic | A clear line between `helper` and `controller that now owns the truth` | The repo can quietly grow orchestration in fragments while pretending it is still only Compose | For each helper, state what truth it owns, what it does not own, and what would justify promoting it |

## The contradictions most likely to be overclaimed

### 1. Public entry

This is the easiest place to lie because:

> traffic can hit more than one node

feels like anti-SPOF progress.
It is only first-hop progress until the wrong-node path is proven.

### 2. Traefik-centered ingress

Traefik is one of the strongest real parts of the runtime.
That makes it easy to overcredit.

Traefik executes decisions well.
It does not by itself invent shared cross-node truth.

### 3. Stateful services

This is where almost every ecosystem cheats.
Hostname stability, TCP reachability, and restartability are not ownership,
promotion, or client rediscovery semantics.

## What these contradictions are really measuring

Each contradiction is really asking:

- where does this truth live now?
- who has to complete it when reality gets sharp?
- what artifact would move it out of folklore?

If the answer is still `the operator`, then the contradiction is still active
even if the surrounding stack looks more serious than before.

## What the docs must keep refusing

Until the corresponding proof packet exists, the docs should keep refusing
sentences like:

- `wrong-node behavior is basically solved`
- `multi-node ingress is already HA-shaped`
- `services.yaml is effectively present`
- `Headscale makes peer-forward discovery mostly handled`
- `MongoDB and Redis are resilient because they are reachable through Traefik`

Those are forbidden because they misplace custody of truth.

## What would actually make the repo feel like it has a real option

The smallest meaningful contradiction reduction is not:

> add a famous tool.

It is:

1. expose current placement truth outside private memory
2. use that truth on a receiving node to decide local versus peer route
3. prove one wrong-node stateless HTTP path
4. prove whether that same route survives backend loss
5. compare protected-route continuity across the handoff
6. keep TCP and stateful surfaces on their own harsher timelines

That sequence is the difference between:

- a better explanation of the user's wound

and:

- one less reason for the user to believe the ecosystem is still mostly fake.

Another way to say the same thing:

- a real option is not just another tool family
- a real option is one that kills one socially humiliating sentence

Examples:

- `the wrong node still needs me to explain locality`
- `the reachable peer still needs me to explain safety`
- `the fallback still needs me to explain whether it survives the failure`

If those sentences survive, the option may still be interesting, but it has not
yet become real in the way the user is demanding.

## The honest bottom line

The repo is already beyond toy infrastructure.
It is not yet beyond hidden-human completion.

That means the contradictions are not mostly about prettier diagrams, more
famous tools, or cleaner prose.
They are about whether one operator is still privately supplying the truth the
platform pretends to own.

If that answer remains yes, the contradiction is still live no matter how
serious the surrounding stack looks.
