# Ingress and Failover Evidence

This page is the proof boundary for the repo's ingress story.

Ingress is the easiest place for infrastructure docs to start lying by
accident.

A live reverse proxy, wildcard DNS, healthchecks, and a fallback-looking helper
can make a stack feel "HA enough" long before it can actually preserve a real
request after it lands on the wrong node or loses the local backend.

This page exists to stop that inflation.

It also exists because ingress is where the ecosystem most often gives the
user something that feels like an option while still leaving the real burden
untouched.

The repeated pattern is:

- more public entry targets
- more proxy vocabulary
- more dynamic routing language
- still no runtime truth that makes the receiving node stop being
  semantically wrong on the bad day

This page therefore has to do more than list ingress ingredients.
It has to filter against false emotional closure.
The wrong reader reaction would be:

- there is Cloudflare
- there is Traefik
- there are healthchecks
- there is a fallback-looking helper
- therefore the ingress problem is mostly handled

That reaction is exactly what this page has to prevent.

It also exists because the archive shows how quickly answers regress into edge
theater:

- multiple A records get narrated as if request preservation is done
- "the proxy is healthy" gets narrated as if wrong-node handling is solved
- tool names get recommended as if naming them closes the failover gap
- one half-working route path gets treated like a distributed service model

The user keeps pushing back on exactly that drift.

That pushback is not nitpicking.
It is the center of the repo's whole existence.

## What this page is trying to prove

This page is not trying to prove that ingress HA is complete.

It is trying to prove seven much narrower things:

1. the root runtime already has a substantial real ingress surface
2. the any-node-entry and peer-forward dream is explicit in repo-native intent
3. HTTP and raw TCP are both already part of the live edge, but they are not
   the same problem
4. tracked shared placement truth is still absent from the root runtime
5. route persistence under failure is not yet trustworthy enough to narrate as
   solved
6. middleware and auth continuity are part of routing correctness here, not a
   side concern
7. current evidence supports a serious ingress architecture, but not yet a
   fully proved distributed request-preservation surface

That last point is the entire reason this page exists.

## What the user is actually asking this page to protect

The user is not asking whether Cloudflare can send traffic to multiple nodes.

The user is asking whether any node can:

- receive a request for the wrong service
- discover the right destination
- preserve the route after local failure
- and still behave like the same service instead of a degraded workaround

That is a much harder claim.
This page exists to stop the docs from silently swapping the hard claim for the
easy one.

It also needs to stop a subtler downgrade:

swapping the user's actual question for "which ingress ingredients do we
already have?"

That downgrade sounds harmless, but it is one of the most common ways to make
the docs feel informative while still failing the request.

The user is not starved for ingredient lists.
The user is starved for a real answer about whether those ingredients have
crossed the line into preserved request meaning on the bad day.

That is the underlying anti-benchmark:

if the docs become satisfied by "traffic can arrive at more than one node,"
then they have already accepted exactly the downgrade the user is trying to
escape.

The archive evidence for that pressure is unusually direct:

- [`../source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`](../source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md)
  shows the user explicitly saying the nodes are already meant to L7 or L4 to
  one another and that what is still missing is unified service discovery.
- [`../source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`](../source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md)
  shows direct dissatisfaction with the usual load-balancer names because they
  do not, by themselves, deliver the Cloudflare-style failover behavior the
  user is really looking for.
- [`../source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md`](../source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md)
  is a good reminder that even the syntax surface of failover thinking can
  trick people into overestimating what Traefik's Docker labels alone actually
  support.

## Evidence classes this page relies on

This page uses all four evidence classes, but not equally.

## Quick claim router

Use this map before letting any ingress sentence become confident:

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "the edge stack is real and serious" | Class 1 | `docker-compose.yml`, `compose/docker-compose.coolify-proxy.yml`, `compose/docker-compose.docs.yml`, `compose/docker-compose.core.yml` | that wrong-node survival is end-to-end solved |
| "the repo wants any-node entry plus peer forwarding" | Class 2 | `.github/copilot-instructions.md`, `README.md` | that the tracked runtime already performs that contract |
| "the repo already knows where current ingress breaks" | Class 3 | `docs/osvc_ingress_ha.md`, `docs/INFRASTRUCTURE_MASTER_PLAN.md` | that the repair path is live |
| "the user rejects ordinary load-balancer answers" | Class 4 | archive conversations linked below | that ecosystem dissatisfaction itself proves a local technical result |

If a paragraph crosses rows, it should say which row is doing which work.

## Class 1: live implementation evidence

Used for:

- active edge components
- routed HTTP services
- routed TCP services
- live network and label surfaces

## Class 2: repo-native intent evidence

Used for:

- the intended request model
- the no-heavy-orchestrator routing philosophy
- the role of `services.yaml`

## Class 3: planned architecture evidence

Used for:

- known failover weaknesses
- route-persistence concerns
- future sync and failover direction
- the ingress versus stateful split

## Class 4: archive-pressure evidence

Used for:

- why "node reachable" is not enough
- why wrong-node success is the real threshold
- why fake HA language is specifically unacceptable in this repo

## Strongest live ingress anchors

Primary files:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- [`compose/docker-compose.docs.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.docs.yml)
- [`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
- [`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)

## What the live files concretely prove

## 1. The root runtime has a real ingress network shape

[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
defines and uses:

- `publicnet`
- `backend`
- `warp-nat-net`

What that proves:

- the ingress story is not abstract prose layered on a flat single-network
  stack
- the root runtime already partitions public entry, internal traffic, and
  specialized routing concerns

What it does **not** prove:

- that node-to-node ingress truth is coordinated correctly
- that the receiving node knows what to do with a wrong-node request

## 2. Traefik is a real runtime control surface, not a symbolic proxy mention

[`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
contains:

- live `traefik`
- Docker provider wiring
- file-provider usage
- TinyAuth integration
- CrowdSec-related config
- forward-auth middleware labels
- `docker-gen-failover`

What that proves:

- the edge layer is already one of the stack's real control planes
- ingress correctness already depends on auth, middleware, and helper
  components, not just one router

What it does **not** prove:

- that the helper path survives failure correctly
- that auth and middleware remain semantically stable during peer fallback

That second point matters a lot.
The user is not merely asking for route continuity.
They are asking for request-meaning continuity.

That last distinction is one of the easiest places for infrastructure writing
to become fake.
A request that still gets *some* answer is not automatically the same request
contract surviving.

This is one of the strongest places where the docs need to stay psychologically
accurate, not just technically accurate.

Many stacks can make the operator feel relief by producing a response.
Far fewer can prove the same request stayed the same request after locality,
backend availability, and policy path stopped lining up cleanly.

That is why ingress evidence here should be read as a held-out test against the
promise encoded in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md),
not as a generic reverse-proxy feature review.

The question is not whether Traefik, Cloudflare, or the Compose edge looks
serious.
The question is whether they currently make the Compose-first any-node-entry
dream more true without silently reinstalling a sacred public node in practice.

That means this page should always keep four subquestions separate:

1. Can traffic hit more than one node?
2. Can the receiving node know whether it is the right node?
3. Can the receiving node find a healthy eligible peer if it is the wrong one?
4. Can the request preserve auth, middleware, and service meaning while doing
   that?

Most ecosystem answers stop after question 1.
Many stronger-looking answers stop after question 2.
This repo is still trying to earn honest yeses to questions 3 and 4.

## 3. The docs site is routed through the same ingress story

[`compose/docker-compose.docs.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.docs.yml)
shows `mkdocs` with:

- `publicnet`
- Traefik HTTP router labels
- Traefik service port labels
- healthcheck path and interval

What that proves:

- the documentation surface is not outside the architecture being described

What it does **not** prove:

- that the docs route itself would survive wrong-node handling or backend-loss
  conditions

## 4. The live edge already mixes HTTP and raw TCP concerns

[`compose/docker-compose.core.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.core.yml)
shows:

- `mongodb` with Traefik TCP router labels
- `redis` with Traefik TCP router labels

[`compose/docker-compose.firecrawl.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.firecrawl.yml)
shows:

- `firecrawl` on `backend` and `publicnet`
- Traefik HTTP router labels
- healthcheck-backed dependencies on Redis, Postgres, RabbitMQ, and Playwright

What that proves:

- the live ingress surface already spans both L7 and L4/TCP exposure classes
- app-level complexity is already attached to the edge surface

What it does **not** prove:

- that HTTP and TCP should be spoken about with the same failover confidence
- that TCP-exposed state-bearing services have meaningful continuity semantics

This distinction matters because one of the easiest ways to lie is to let a
stronger HTTP story leak into stateful TCP confidence.

## What the live files do not prove

The current root runtime does not, by itself, prove:

- that any node can always identify the correct healthy peer
- that the route required for fallback survives local backend loss
- that peer eligibility is based on trusted current-state convergence
- that auth and middleware remain semantically identical during fallback
- that the ingress story is equally mature for HTTP and raw TCP

That missing proof should be read as a routing-fidelity gap, not a cosmetic
gap.
The danger is not merely that the docs would be slightly optimistic.
The danger is that they would declare relief before the system has actually
stopped cheating with hidden human interpretation.

This is also where the docs have to resist turning "the ingress stack is
substantial" into "the ingress problem is mostly solved."

Those statements are not close.
The first is real current-state evidence.
The second is exactly the kind of smoother, more flattering summary this
knowledgebase is trying to stop.

That missing proof is not a detail.
It is the main honesty boundary for this part of the repo.

This is exactly where the archive keeps forcing the same correction:

- first-hop reachability is not the same as service continuity
- a route existing in one provider is not the same as recovery surviving the
  relevant failure
- "Traefik can route it" is not the same as "the wrong node can preserve the
  request truthfully"

Another way to phrase it:

- the worktree proves a real ingress stack exists
- it does not yet prove that wrong-node handling is trustworthy end to end

That is the real compression of the whole page.

The ingress layer is already too serious to dismiss.
It is still not allowed to narrate itself as if seriousness had already become
preserved meaning.

## Strongest intent and planning anchors

Primary files:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

## What these files explicitly say

## 1. The desired request model is deliberate

[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
explicitly describes:

- no central orchestrator
- `services.yaml` as service-registry concept
- L7 Traefik handling
- separate L4 handling
- Cloudflare node-level failover
- any-node entry with local serve or peer forward

What that proves:

- wrong-node request survival is part of the architecture dream, not a later
  reinterpretation

What it does **not** prove:

- that the root runtime already behaves that way under failure

## 2. README preserves the honesty wall

[`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
explicitly keeps open:

- placement truth gaps
- failover generation gaps
- convergence gaps
- stateful correctness gaps

What that proves:

- the repo already knows ingress prose can overclaim

What it does **not** prove:

- that the current runtime has crossed the proof threshold the README is
  warning about

## Why tool-name answers keep failing this repo

One of the most repeated archive frustrations is that the ecosystem keeps
answering this problem with product categories:

- reverse proxy
- load balancer
- service discovery
- ingress HA

Those categories are not useless.
They just keep arriving one level too early.

The user is not asking:

> what should I put in front of my nodes?

They are asking:

> what live truth lets the receiving node stop being semantically wrong when
> the request lands there first?

That is why naming Traefik, NGINX, HAProxy, or Cloudflare features never closes
the matter by itself.
Those tools can participate in the answer.
They are not the answer until they are tied to current placement truth, route
survival, and policy continuity.

That sentence is one of the deepest translation keys in the docs.

The user is not asking for a thing called "ingress HA."
The user is asking for the smallest truthful machinery that stops the first
healthy wrong node from becoming an interpretive error.

That should also control future wording.

Good wording:

- "multi-node first-hop entry is supported as intent"
- "the live edge is substantial"
- "shared placement truth is still absent from the tracked root runtime"
- "route-preserving peer fallback is still below proof threshold"

Bad wording:

- "ingress HA is basically there"
- "the stack already supports distributed failover"
- "Traefik plus Cloudflare solves the problem"
- "dynamic generation closes wrong-node routing"

## 3. The planning layer sharpens the subproblem split

[`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
separates:

- node-scoped hostnames
- global hostnames
- Cloudflare LB / VIP / round-robin entry options
- raw TCP as a separate problem class needing L4 handling and stateful honesty

[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
records:

- `docker-gen-failover` deletes routes on container stop
- `services.yaml` should exist as distributed registry
- sync-agent and failover-agent are the planned repair direction

What that proves:

- the repo already knows exactly where the current ingress story breaks

What it does **not** prove:

- that the repair mechanisms are live

This is another place where the docs need to resist a familiar temptation:

- precise diagnosis
- plus plausible repair direction
- becomes narrated as near-complete solution

The user has already seen too many systems make that move.

## Critical negative evidence

## Negative fact 1: tracked placement truth is still absent

Current root worktree check:

- `services.yaml`: absent

What that proves:

- the placement-registry concept is central
- the tracked root runtime still lacks the shared truth layer the routing dream
  expects

What it does **not** prove:

- that no informal placement knowledge exists anywhere

The point is not that operators know nothing.
The point is that the shared tracked truth surface is still missing from the
priority runtime.

That missing surface is not paperwork.
It is one of the main differences between:

- "I can still probably make this work"
- and "the runtime now carries enough of its own truth to stop depending on
  private recollection"

## Negative fact 2: route generation exists without route-persistence proof

Current evidence split:

- live runtime contains `docker-gen-failover`
- planning docs say its behavior defeats failover when containers stop

What that proves:

- dynamic generation exists as a live tactic
- dynamic generation is not yet trustworthy failover evidence

What it does **not** prove:

- that the replacement path is active

## Negative fact 3: node entry is easier than request correctness

Intent and planning sources strongly support:

- Cloudflare-backed node entry

But they do not prove that the receiving node:

- knows current placement truth
- preserves middleware semantics
- forwards correctly after local backend failure

This is exactly why ingress claims must stay layered rather than celebratory.

## Negative fact 4: auth and middleware are part of the failover problem, not decorations

Current live edge structure already involves:

- TinyAuth
- forward-auth middleware
- CrowdSec-related surfaces
- Traefik file and Docker provider combinations

What that proves:

- route correctness here is not just transport

What it does **not** prove:

- that those semantics are already preserved across fallback paths

This is one of the most important differences between the user's standard and a
generic proxy-failover standard.

The generic standard asks whether packets and responses survived.
The user's standard asks whether the architecture stopped cheating.

## Strongest honest current answer

If a reader asks, "What do we actually know right now?" the shortest
defensible answer is:

> The tracked Compose-first runtime already has a serious ingress stack and the
> repo explicitly wants any-node entry with local-first or peer-forward
> handling, but the priority runtime still lacks shared tracked placement truth
> and still does not prove end-to-end wrong-node request preservation with
> stable auth, middleware, and service semantics.

Anything stronger than that needs stronger evidence than this page currently
has.

## The strongest honest current claim

The strongest current ingress claim the repo can support is:

> the priority runtime already contains a substantial, serious ingress and
> edge-control stack, and the desired any-node-entry plus peer-forward behavior
> is explicit in repo-native intent, but the tracked runtime still does not
> prove shared placement truth, durable fallback-route persistence, or end-to-end
> wrong-node request preservation

That is a large claim.
It is also much smaller than "ingress HA is solved."

That difference is the whole reason this page exists.

## Bottom line

The repo's ingress story is real enough to deserve serious attention and strict
enough documentation.

But the worktree still stops short of proving the exact behavior the user
cares about most:

- a request lands on the wrong node
- the local backend may be gone
- the receiving node still knows the right next step
- the route survives
- the policy survives
- the request still means the same thing

Until that chain is proven, ingress should be described as:

- structurally serious
- architecturally aligned
- still proof-limited

That is the honest ceiling.

It is also the page's main reconstruction result:

the ingress problem here is not exposure, not first hop, and not proxy
selection in the abstract.
It is whether the first healthy wrong node has enough durable truth to stop
turning the whole platform back into a private operator memory game.
