# Current Compose Runtime

This page describes the strongest live implementation surface in the repo:

the root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
plus the fragments it includes directly.

That matters because this is still the priority implementation.
If a claim cannot survive contact with the root Compose entrypoint, it is not a
live-runtime claim yet.

This page is intentionally narrower than the repo's bigger dream.
It is not here to retell the aspiration.
It is here to answer a harsher question:

> what does the tracked root runtime actually give us today, and where does it
> still stop far short of the wrong-node, peer-aware, anti-SPOF behavior the
> user actually wants?

## What this page is and is not allowed to prove

This page is authoritative about:

- what the priority root Compose surface currently includes
- what the active include graph visibly expresses
- which truths are still missing from the tracked root runtime

This page is not authoritative about:

- future control-plane direction as a winner
- user dream reconstruction by itself
- wrong-node or failure-path success unless explicitly exercised elsewhere

This page is the strongest inventory of live authored runtime shape, not a
global resilience verdict.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "this is in the current root runtime" | live-root evidence | root `docker-compose.yml`, active include fragments, `docker compose config` when available | that the behavior works under failure |
| "the runtime shows serious edge and workload breadth" | live-root evidence | root graph, included fragments, network and service declarations | that breadth equals shared truth |
| "this truth is still missing from the root runtime" | live-root evidence plus negative facts | absence of root `services.yaml`, lack of exercised failure proof, page-local negative sections | that nothing informal exists anywhere |
| "this missing truth matters to the dream" | synthesis linked to stronger intent pages | this page plus architecture and research pages | that dream fulfillment is already close just because the runtime looks serious |

If a sentence starts celebrating modularity as if it were shared truth, this
page is being overread.

That distinction matters because large Compose stacks are extremely easy to
overread.

A lot of self-hosting documentation sees:

- many services
- many labels
- many networks
- many helper containers

and starts narrating "platform."

This page exists to stop that move.

It also exists to stop a softer version of the same mistake:

because the root runtime is already large, modular, and full of serious
supporting infrastructure, it becomes easy to narrate "half a distributed
platform" as if it were "most of a distributed platform."

That temptation is one of the main sources of ambiguity in this repo.

It is also one of the most dangerous interpretive shortcuts here.
The stack already contains enough real machinery that a reader can feel
reassured before the hard missing truth has even been named.

This page is also doing a reconstruction job, not just an inventory job.

The root Compose implementation is one of the richest surviving artifacts of
what the system currently believes about itself.
But rich artifacts are not the same thing as self-explaining runtime truth.

That distinction matters because the user is effectively asking the docs to do
for infrastructure what a high-fidelity retrieval system would do for a messy
personal archive:

- preserve the strongest signals
- keep contradictions visible
- refuse comforting simplifications
- and reconstruct the actual hidden pressure instead of answering a smaller
  easier question

## What this page is for

If an operator asks:

> what is actually present in the priority Compose implementation right now?

this page should answer in concrete terms:

- what the root file owns directly
- which fragments are active in the live include path
- which networks and service classes dominate the stack
- which exposure models already exist
- which parts of the larger multi-node story are still architectural pressure
  rather than proven runtime behavior

It should also answer a sharper question:

> which parts of the current runtime still encode node-local truth even while
> the surrounding docs talk about any-node entry and anti-SPOF pressure?

It should not be used to smuggle in future control-plane behavior or to widen a
local Compose truth into a whole-cluster truth.

The main smaller easier question this page refuses to answer is:

> is the root runtime big, modular, and serious?

The answer to that is obviously yes.
The harder and more useful question is:

> does the root runtime already carry enough explicit shared truth that the
> multi-node dream no longer depends on operator recollection?

That answer is much more conditional, and the entire point of this page is to
keep that condition visible.

Another way to phrase the condition is:

> does the root runtime merely look like something that ought to know the
> cluster truth, or does it already expose enough inspectable shared truth to
> stop leaning on remembered placement?

That sharper question is closer to the user's real complaint.

## The root implementation boundary

The priority implementation is the root Compose file itself.

That means the strongest live truths today are:

- shared networks are declared centrally in root
- several important services are still authored directly in root
- multiple domain fragments are part of the active include path
- the repo is Compose-first in real deployment shape, not just in rhetorical
  aspiration

This boundary matters because the user is not asking for a decorative
multi-node story.

They are asking:

> how much of the anti-SPOF, any-node-entry, peer-forward dream is already
> materially visible in the tracked root runtime?

The honest answer is:

- a lot of edge scaffolding is already visible
- a lot of naming and routing pressure is already visible
- the live shared placement truth is still missing
- the live wrong-node truth is still incomplete
- the strongest distributed claims still live more in intent and planning than
  in a proven shared runtime substrate

This page has to preserve those proportions.

Those proportions are emotionally important, not just technically important.

If the docs flatten them, the user gets exactly the same failure pattern that
the repo is rebelling against elsewhere:

- lots of serious-looking machinery
- lots of named options
- still no clean answer to what the receiving node knows on the bad day

The root runtime matters so much because this repo is no longer allowed the
luxury of vague aspiration.
It already has enough real ingress policy, real state, real helper surfaces,
and real operator tooling to punish sloppy language immediately.

That means the docs cannot grade on a curve just because the stack is clearly
serious.
Seriousness is already proven.
The unresolved question is whether seriousness has become inspectable
multi-node truth or is still largely sophisticated node-local truth.

There is also a subtler danger this page has to prevent:

the more serious the runtime looks, the easier it becomes for readers to
accidentally invert the burden of proof.

Instead of asking:

> what exact shared truth is now inspectable?

they start asking:

> with this much machinery present, surely the remaining gaps cannot be that
> important anymore?

That inversion is poison for this repo.

The user's actual complaint is precisely that the remaining gaps are the whole
problem, even after a great deal of impressive machinery is already present.

That is why this page has to keep saying the same thing from different angles:

- seriousness is already proven
- breadth is already proven
- modularity is already proven
- shared placement truth is not
- wrong-node truth is not
- request-preserving bad-day truth is not

That means this page has to keep teaching a disciplined reading habit:

- visible machinery is real
- modularity is real
- supporting edge infrastructure is real
- those things still do not become shared placement truth, wrong-node truth,
  or request-preserving failover just because they accumulated to a serious
  size

If the page stops saying that in multiple ways, it becomes easier to read the
runtime as "almost a platform" instead of "a serious stack still missing the
explicit shared truth that would make the bad day coherent."

## Root-owned networks

The root file defines three core networks:

- `publicnet`
- `backend`
- `warp-nat-net`

Those names already reveal three different runtime concerns:

- `publicnet`: externally reachable or ingress-adjacent traffic
- `backend`: east-west internal service traffic
- `warp-nat-net`: specialized egress or routing handling

That proves the stack is already partitioned by control surface, not just by
application category.

What it does **not** prove:

- that cross-node network truth is coordinated correctly
- that wrong-node forwarding logic is already trustworthy
- that these network names by themselves mean the same topology story on every
  participating host

The networks prove segmentation pressure.
They do not prove shared routing truth.

## Root include path

The root file includes these fragments directly:

- `compose/docker-compose.coolify-proxy.yml`
- `compose/docker-compose.docs.yml`
- `compose/docker-compose.firecrawl.yml`
- `compose/docker-compose.headscale.yml`
- `compose/docker-compose.llm.yml`
- `compose/docker-compose.metrics.yml`
- `compose/docker-compose.stremio-group.yml`
- `compose/docker-compose.warp-nat-routing.yml`
- `compose/docker-compose.wishlist.yml`

This proves something important about the current runtime:

- the stack is already modular
- the modularity is still expressed through Compose includes
- there is still no clearly promoted shared cluster substrate above that include
  graph

That is one reason the repo still feels like it lives in the "missing middle"
zone.

The authoring surface is serious.
The distributed truth surface is still weaker than the authoring surface.

That sentence captures a lot of the user's actual frustration.

The repo already has enough Compose structure to look like a platform from far
away.
The pain begins when a request lands on the wrong node and that structural
richness still does not answer where the service actually lives now.

The stack already looks like something that should know a great deal.
What it still lacks is not categories.
It lacks a clear, shared answer to what a receiving node knows when the
requested service is elsewhere.

## What can be proven directly from file inspection

Even before interpolation succeeds, the root file already proves quite a lot.
That proof is incomplete, but it is not vague.

It also proves several negative facts that earlier docs were too polite about.

This section should therefore be read with a stricter habit than normal
Compose inspection.

Do not ask:

- how rich does the stack look?
- how many concerns are clearly represented?
- how many helper layers have already been introduced?

Ask instead:

- what does a receiving node know without operator recollection?
- what does root Compose still leave socially reconstructed?
- which bad-day answers are still missing even after the supporting machinery
  is acknowledged fully?

That is the reading discipline that keeps this page aligned with the user's
actual frustration.

## Negative fact 1: root Compose is not a placement registry

The root entrypoint proves:

- service definitions
- labels
- networks
- include shape
- dependency hints
- exposure intent

It does **not** prove:

- where a service lives across multiple nodes right now
- which peer should receive a wrong-node forward
- which peer is currently eligible for fallback

This is the place where the docs have to resist the strongest form of
infrastructure theater:

the belief that enough labels, helpers, and routes eventually add up to a
placement story even when no explicit shared placement story exists.

That belief is exactly how large Compose systems become emotionally
misleading.
They accumulate so much intent that readers start crediting them with runtime
knowledge they have never actually externalized.

## Strongest honest current answer

If a reader asks, "What does the current Compose runtime actually buy us
today?" the shortest defensible answer is:

> It buys a real, serious, modular live authoring surface with substantial edge
> and workload machinery already present in the root runtime, but it still does
> not buy explicit shared placement truth, trustworthy wrong-node behavior, or
> a self-explaining bad-day recovery contract in the tracked root implementation.

Anything stronger than that needs exercised proof from outside plain file
inspection.

That missing truth is exactly why `services.yaml` keeps reappearing in the repo
story.

## Negative fact 2: root Compose is not yet a wrong-node survival contract

The root runtime contains enough ingress-adjacent machinery that it can tempt
readers into a stronger conclusion:

- a public node can receive the request
- Traefik and supporting helpers are present
- peer-aware behavior is clearly part of the repo's direction
- therefore wrong-node entry must be substantially solved already

That conclusion is still too generous.

Root Compose does not yet prove, by itself:

- that a receiving node can always distinguish local versus remote ownership
- that it can always choose the correct healthy peer at request time
- that the route needed for recovery survives the local failure that triggered
  recovery
- that middleware and auth continuity remain intact after cross-node handoff

This matters because "the edge stack looks serious" is one of the easiest
places for documentation to over-credit the runtime.

The user is not looking for evidence that the edge stack looks expensive,
modular, or familiar.
The user is looking for evidence that the wrong machine stops being the start
of a human-reconstruction exercise.

The root runtime does not yet prove that globally.

## Negative fact 3: root Compose is still stronger at authorship than at shared truth

One of the clearest repo-wide patterns is now hard to miss:

- authoring truth is rich
- domain separation is rich
- service coverage is rich
- supporting infrastructure is rich
- shared current-state truth is still weaker than all of those

This is not a casual observation.
It is one of the most important explanations for why the user remains
frustrated despite the stack already being substantial.

The repo already knows how to describe many concerns.
It is still maturing toward a world where nodes can answer those concerns with
shared inspectable truth instead of remembered placement, remembered peer
eligibility, and remembered rescue routes.

That is the practical reason this page refuses to flatter the root runtime.
The stack deserves credit for seriousness.
It does not deserve automatic credit for having crossed the missing-middle
threshold just because its Compose surface is now large and disciplined.

This is one of the sharpest current gaps between:

- "large serious stack"
- and "distributed request-preserving platform"

It is also the gap between:

- a system that looks like it should know where things live
- and a system that can prove where things live without appealing to memory or
  folklore

## Negative fact 2: no live tracked root `services.yaml` currently closes that gap

The repo repeatedly describes a lightweight `services.yaml` current-state
registry.

The tracked root worktree still does not ship one.

So the current runtime boundary is:

- strong on Compose authoring truth
- weak on shared live placement truth

Another blunt way to say it:

the root runtime is already large enough to feel distributed while still being
small enough to depend on remembered placement.

That is not a minor documentation nuance.
It is one of the main hidden failure modes the user is asking this whole
knowledgebase to preserve instead of smoothing over.

That is one of the most misleading architectural middle states possible.

It is also one of the most frustrating states for a user who is explicitly
looking for more real options, not more ways to rename the same hidden burden.

That sentence should stay painful, because it is one of the most important
proof ceilings in the whole project.

## Negative fact 3: edge machinery existing is not the same thing as failover being correct

The active stack clearly includes substantial edge machinery.

That proves seriousness.
It does not prove that the exact path needed during local backend failure stays
alive when the failure actually happens.

This is one of the main places where infrastructure docs like to cheat:

- proxy exists
- helper exists
- generator exists
- therefore failover is "basically there"

The user is explicitly asking the docs not to do that.

## Negative fact 4: merged-runtime certainty is gated by env and secrets

The natural command for merged-runtime proof is:

```bash
docker compose config --services
```

But in the current repo state, that command does not cleanly resolve unless
enough environment is present.

A local attempt fails without values such as:

- `DOMAIN`
- `TS_HOSTNAME`
- required secrets like `SEARXNG_SECRET`

So the strongest honest rule is:

- file inspection proves root structure and active include path
- full merged inventory proof still depends on an adequately populated env and
  secret surface

That matters because if the merged graph cannot be rendered from the repo alone,
the docs should not pretend the repo itself is a self-contained proof of the
whole runtime.

## The root file is not just glue

It still owns significant services directly, including representative examples
like:

- `mongodb`
- `redis`
- `searxng`
- `code-server`
- `homepage`
- `watchtower`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `dozzle`
- `portainer`
- `telemetry-auth`
- `bolabaden-nextjs`
- several `biodecompwarehouse*` services

It also owns several concrete node-affined storage assumptions directly.

Representative examples visible in the root file include:

- MongoDB bind-mounting `${CONFIG_PATH:-./volumes}/mongodb/data:/data/db`
- Redis bind-mounting `${CONFIG_PATH:-./volumes}/redis:/data`
- code-server bind-mounting `${CONFIG_PATH:-./volumes}/code-server/dev/config:/config`
- SearxNG bind-mounting local config and cache paths under `${CONFIG_PATH:-./volumes}`

That matters because the current runtime is not merely "capable of state."
It already embeds real local gravity into the priority implementation.

Any routing or failover narrative has to be measured against that gravity.

That means the root file remains a major authoring surface for:

- core datastores
- socket mediation
- operator dashboards
- update automation
- custom project services
- public-facing application entrypoints

Any future control-plane promotion has to respect that reality.

The repo is not starting from a small blank slate.
It is starting from a broad, already meaningful Compose contract.

## The edge surface is already substantial

The live request path is not just a lone Traefik container.
The active root implementation already points to an edge or control cluster
around:

- `traefik`
- `crowdsec`
- `crowdsec-init`
- `tinyauth`
- `nginx-traefik-extensions`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `docker-gen-failover`
- `cloudflare-ddns`

This is one reason generic infrastructure prose becomes dangerous here:

there are already enough moving parts for "proxy is running" to mean almost
nothing by itself.

More importantly, this edge surface proves the repo is already trying to own:

- node entry
- auth continuity
- middleware continuity
- route generation
- service exposure policy
- the first half of any future wrong-node recovery path

That does **not** make wrong-node success real yet.
It does prove the repo is structurally building toward it rather than merely
talking about it.

This is one reason the repo has to be documented so carefully.

Once a stack reaches this level of edge sophistication, people start
retroactively granting it capabilities it has not earned yet.
The docs have to resist that social upgrade just as hard as they resist
technical overclaiming.

This page exists partly to deny exactly that social upgrade.
The user is not missing admiration for the stack.
They are missing a system that stops requiring mental completion on the bad
day.

It also proves something subtler:

the repo is already trying to preserve more than reachability.
The edge surface is visibly trying to preserve:

- auth meaning
- middleware identity
- node-scoped versus global hostname semantics
- Docker visibility boundaries through socket proxies
- generated route behavior

Those are the concerns of a system that is approaching a real distributed edge
problem, not merely exposing containers.

But those concerns existing is still not the same thing as them being unified
by one shared, current, inspectable truth surface.
That gap is what keeps the root runtime from being allowed to inherit the full
anti-SPOF claim yet.

## The runtime is not one homogeneous thing

One of the easiest ways to misunderstand this repo is to read the stack as one
flat service list.

The current runtime is easier to reason about as service classes, because the
failure modes and maturity ceilings are radically different by class.

## Class 1: ingress, auth, and edge-control infrastructure

Representative services:

- `traefik`
- `crowdsec`
- `crowdsec-init`
- `tinyauth`
- `nginx-traefik-extensions`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `docker-gen-failover`
- `cloudflare-ddns`

This class matters disproportionately because it is where:

- public entry
- auth behavior
- middleware continuity
- route generation
- node and global hostname semantics

all collide.

It is also the class closest to the user's actual dream of:

- any-node entry
- local-first service
- peer-forward fallback

Which is exactly why it must not be overclaimed.

It is also the class most likely to generate false confidence.

A stack with Traefik, CrowdSec, forward auth, dynamic routing helpers, socket
proxies, Cloudflare, and dashboards feels like it ought to already know how to
survive wrong-node entry.

"Ought to" is not evidence.
This page exists to break that emotional inference.

The user has clearly run into too many ecosystems where "ought to" gets sold as
"therefore does."
This page is supposed to make that move harder inside this repo.

## Class 2: observability and diagnostics

Representative services:

- `prometheus`
- `victoriametrics`
- `grafana`
- `alertmanager`
- `loki`
- `promtail`
- `blackbox-exporter`
- `node_exporter`
- `process-exporter`
- `cadvisor`
- `redis-exporter`
- `mongodb-exporter`
- `dozzle`
- `autokuma`
- `logrotate-traefik`

The repo is already observability-rich.

That proves the project is not naive about runtime introspection.
It does **not** prove the observed truths are yet unified enough to drive
cross-node routing safely.

## Class 3: state-bearing core services

Representative services:

- `mongodb`
- `redis`
- `nuq-postgres`
- `rabbitmq`

This class matters because it prevents the documentation from pretending the
stack is mostly stateless.

Stateful correctness is not a future appendage here.
It is already part of the priority runtime surface.

That means every broad resilience claim has to be judged against a harder
standard.

## Class 4: public-facing or operator-facing application surfaces

Representative services:

- `bolabaden-nextjs`
- `homepage`
- `code-server`
- `searxng`
- `firecrawl`
- docs site surfaces

This is the most obvious user-facing class, but it is not the only class that
matters.

One of the reasons earlier docs felt too simple is that they read the whole
system through this class alone, which makes the stack sound cleaner and less
stateful than it actually is.

That simplification also makes the repo sound closer to solved than it is,
because the public app surfaces are the easiest part of the story to narrate
cleanly.

## What the current runtime genuinely supports as claims

The strongest claims the current runtime can support are:

- the repo is Compose-first in real deployment shape
- the root stack is broad and serious, not toy-sized
- a substantial ingress/auth/middleware surface is live in authored config
- stateful services already exist in the priority runtime
- the stack is already rich enough that hidden operator truth becomes a serious
  architectural problem

The current runtime still cannot, by itself, support stronger claims like:

- wrong-node HTTP is broadly solved
- peer fallback survives local backend loss
- cross-node placement truth is live and shared
- auth and middleware continuity under fallback are proven
- raw TCP exposure equals resilient TCP failover
- stateful services have meaningful multi-node correctness

It also cannot yet support a deeper operator claim:

- that one tracked root truth surface is enough to explain wrong-node behavior
  without synthesizing Compose, fragments, helper intent, dashboards, and human
  recollection

Those proof ceilings are the entire point of this page.

That last ceiling matters more than it might initially sound.

The user is not just asking for infrastructure that works after enough careful
reading.
The user is asking for infrastructure that eventually stops requiring that kind
of private reconstruction to feel trustworthy at all.

## Bottom line

The current Compose runtime proves that the repo is already serious, broad, and
edge-aware.

It does **not** prove that the repo has crossed the line from:

- large modular Compose system

to:

- request-preserving distributed personal cloud

That gap is not a minor missing feature.
It is the central truth deficit the whole project is trying to close.

This page should keep that gap visible every time the runtime starts sounding
too complete just because it is large.

The root runtime is therefore best read as:

- a strong body of authored intent
- a serious operational surface
- an incomplete but meaningful witness to the user's dream
- not yet the final authority on distributed placement and preserved request
  meaning
