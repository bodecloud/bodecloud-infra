# DevOps Runbook

This is not a generic "how to run Docker" page.

It exists for one harder operational question:

> what exact evidence do we need before we can say the stack preserved a
> request, survived wrong-node entry, or reduced hidden operator burden,
> instead of merely looking healthy?

The neighboring smaller question it must not collapse into is:

> what commands should I run to feel informed?

That smaller question is too easy.
This repo is full of commands that can succeed one layer before the system
itself stops depending on private operator reconstruction.

## What this page is and is not allowed to prove

This page is authoritative about:

- which evidence class a claim belongs to
- which checks are weak, medium, or strong for this repo
- what order to inspect the stack in
- what exact stronger sentence is still forbidden after a successful check

This page is not authoritative about:

- whether any specific route is already resilient
- whether any specific failover path already works
- whether any specific stateful service is already safe

This is a method page.
It is not a certificate page.

## The real operator problem

The user is not mainly asking for more operational fluency.
The user is asking for less hidden operator burden.

The specific burden to keep exposing is:

- hidden topology memory
- hidden placement truth
- hidden convergence truth
- hidden proof translation after the command succeeds

If the operator still has to privately remember:

- which node really hosts the service
- whether the current node is serving locally or forwarding
- whether the peer is merely reachable or actually eligible
- whether auth, middleware, and environment still mean the same thing on the
  peer

then the runbook has not yet reached the user's real pain.

## Start every investigation with a claim sentence

Before running anything, write the claim in this form:

> I am trying to prove `<specific claim>` and I need `<proof class>` evidence.

Good examples:

- "I am trying to prove the merged root Compose graph still resolves, and I
  need authored-shape evidence."
- "I am trying to prove `wishlist.$DOMAIN` answers on a normal day, and I need
  route-behavior evidence."
- "I am trying to prove a wrong-node HTTP request preserved route meaning for
  one service, and I need wrong-node drill evidence."
- "I am trying to prove Redis remained authoritative after node loss, and I
  need stateful-correctness evidence."

Bad examples:

- "I want to make sure things are healthy."
- "I want to check HA."
- "I want to see whether failover works in general."

Those are not claims.
They are invitations to narrate comfort as completion.

## Evidence ladder

The repo needs a stricter evidence ladder than most self-hosting stacks.

| Evidence class | Typical tools | What it can honestly prove | What it still cannot prove |
|---|---|---|---|
| Authored shape | `docker-compose.yml`, `compose/`, `docker compose config --quiet` | The declared graph resolves and the priority implementation surface is inspectable | That requests survive failure |
| Local runtime health | `docker compose ps`, container healthchecks, logs | Containers are up on this node and local processes may be healthy | Wrong-node recovery, backend-loss survival, stateful correctness |
| Edge-route behavior | `curl`, headers, backend identity markers, Traefik logs | One ingress path answered and can sometimes be tied to one backend identity | Peer-forward correctness under loss |
| Wrong-node behavior | controlled drill plus route identity evidence | One named route preserved meaning when the request landed on the wrong node | Generic failover across the whole stack |
| Backend-loss behavior | controlled failure plus before/after evidence | One named route survived a backend-loss condition with known semantics | Stateful authority transfer for unrelated services |
| Stateful correctness | service-specific replication, election, reconnect, and write-path evidence | One stateful service preserved authority and client truth under the tested failure | That every other stateful surface is now HA |

The operator should name the evidence class before running the first command.

## What the priority runtime already tells you to inspect first

The root runtime is still the priority implementation surface.
Start with:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active fragments under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md) for validation constraints

Concrete examples from the current root graph:

- active include set contains:
  - `docker-compose.coolify-proxy.yml`
  - `docker-compose.docs.yml`
  - `docker-compose.firecrawl.yml`
  - `docker-compose.headscale.yml`
  - `docker-compose.llm.yml`
  - `docker-compose.metrics.yml`
  - `docker-compose.stremio-group.yml`
  - `docker-compose.warp-nat-routing.yml`
  - `docker-compose.wishlist.yml`
- root-owned shared networks include:
  - `publicnet`
  - `backend`
  - `warp-nat-net`
- root-owned live stateful services include:
  - `mongodb`
  - `redis`
- root-owned operator-significant services include:
  - `watchtower`
  - `homepage`
  - `code-server`
  - `searxng`

Those facts matter because they tell you the authored surface is real and large.
They do not by themselves tell you which burden has moved out of the
operator's head.

## The operational sequence

Use this order whenever the claim matters.

### 1. Inspect authored shape first

Goal:

> confirm what the priority implementation claims to be.

Typical checks:

```bash
docker compose config --quiet
docker compose config --services
docker compose config | rg "traefik|tinyauth|crowdsec|headscale|mongodb|redis"
```

Use this stage to answer questions like:

- is the graph valid?
- which fragments are active?
- which service names, networks, configs, and secrets are present?
- is the route even declared?

Forbidden upgrade after success:

> therefore failover works

That sentence is still a lie after shape validation.

### 2. Inspect local runtime health second

Goal:

> confirm what this node is actually running and how healthy it looks locally.

Typical checks:

```bash
docker compose ps
docker inspect --format '{{.State.Health.Status}}' traefik
docker logs --tail=200 traefik
docker logs --tail=200 tinyauth
docker logs --tail=200 crowdsec
```

Use this stage to answer:

- is the local container running?
- is the local healthcheck green?
- are local edge services erroring?
- is the expected local backend even present?

Forbidden upgrade after success:

> therefore the route is resilient

A healthy container on one node is still only local-runtime evidence.

### 3. Inspect edge-route behavior third

Goal:

> prove that one route answered and tie that answer to a backend identity if
> possible.

Typical checks:

```bash
curl -I https://service.$DOMAIN
curl -sv https://service.$DOMAIN
docker logs --tail=200 traefik
```

Prefer evidence that answers:

- which hostname was used?
- which router handled it?
- which service/backend did Traefik think it sent to?
- can the response be tied to node identity or app identity?

Forbidden upgrade after success:

> therefore wrong-node routing is solved

A working happy-path request is not a wrong-node proof.

### 4. Only then run a wrong-node drill

Goal:

> prove one specific service can preserve request meaning when traffic lands on
> a healthy node that does not host the service locally.

Minimum packet for this drill:

- exact hostname tested
- receiving node identity
- actual backend node identity
- evidence of how the receiving node knew where to send the request
- evidence that auth and middleware meaning stayed the same
- explicit sentence naming what was still not proven

If the drill still depends on private operator recollection of placement, say
that explicitly.
That means the hidden control plane is still present.

Forbidden upgrade after success:

> therefore any-node entry now works generically

One service-class proof is still one service-class proof.

### 5. Run backend-loss drills separately

Goal:

> prove the route still behaved honestly after the expected local backend went
> away.

Minimum packet:

- before/after route behavior
- exact failure introduced
- whether fallback was local restart, peer-forward, or operator intervention
- whether the route kept the same policy and auth meaning

Forbidden upgrade after success:

> therefore the platform has real failover

The question is always:

failover of what, under which exact failure, with which remaining hidden human
burden?

### 6. Treat stateful drills as their own category

Goal:

> prove one stateful service preserved authority, not just reachability.

Minimum packet:

- authoritative write location before failure
- what happened to authority after failure
- how clients rediscovered the correct topology
- what storage or election mechanism made the claim honest

For current repo examples, keep this boundary explicit:

- `mongodb` being exposed through Traefik TCP is not MongoDB HA
- `redis` answering on the expected port is not Redis Sentinel or promotion
  semantics
- `headscale` using SQLite in the active fragment is a serious persistence fact,
  not a cosmetic implementation detail

Forbidden upgrade after success:

> therefore stateful SPOF is solved

Stateful claims stay the harshest claims in this repo.

## Weak, medium, and strong evidence in practice

### Weak evidence

Examples:

- `docker compose config --quiet` passes
- the service is in `docker compose ps`
- the container healthcheck is green
- the hostname answers once

Weak evidence is still useful.
It proves the stack is not imaginary.
It does not prove the stack preserved meaning under pressure.

### Medium evidence

Examples:

- route behavior can be tied to a known backend identity
- Traefik logs confirm which router and service handled the request
- a backend restart or local outage was exercised and the route still answered
  with named remaining limits

Medium evidence is where the repo becomes more interesting.
It is still not enough for broad anti-SPOF claims.

### Strong evidence

Examples:

- a named wrong-node drill with backend identity and policy continuity evidence
- a named backend-loss drill with explicit before/after semantics
- a stateful drill that proves authority transfer or honest singularity

Strong evidence is not "green output."
It is an inspectable proof packet with a named ceiling.

## What still does not count as a serious runbook result

These are still invalid outcomes:

- a green command with no named proof class
- a route test with no backend identity
- a failover claim with no explicit failure introduced
- a stateful reassurance story that never names authority
- an operator summary that hides how much topology truth was still private

Those are exactly the outcomes that make the docs and tooling sound more mature
without making the platform less socially manual.

## Required close-out after every operational pass

After each investigation, write these four sentences explicitly:

1. `Claim tested:` what exact claim did we test?
2. `Evidence class:` what class of evidence did we gather?
3. `What this proves:` what narrower sentence is now honest?
4. `What is still forbidden:` what stronger sentence would still be a lie?

If those four sentences are missing, the operator will almost always drift into
story inflation.

## Bottom line

The right test for this runbook is not:

> did it make the operator sound competent?

The right test is:

> after following it, what exact private topology sentence did the operator no
> longer have to finish alone?

If the answer is "none yet, but the operator understands the stack better,"
that is still honest progress.

It is not yet the user's dream.
