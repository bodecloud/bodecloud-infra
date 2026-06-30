# Current Compose Runtime

This page is the strongest live implementation inventory in the repository.

If a sentence claims something is part of the priority implementation now, it
has to survive contact with:

- root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- the fragments directly included from
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

Everything else is weaker evidence for live-runtime claims.

That strength is exactly why this page is dangerous.
Runtime inventories are where infrastructure documentation most often starts
sounding useful before it has actually answered the operator's real question.
Once a stack is broad enough, mere presence starts getting misheard as
relief.
This page has to keep interrupting that reflex.

The user already has the broad stack.
What they still do not have is confidence that the stack itself, rather than
the operator, owns the decisive truth when a request lands on the wrong node
or a preferred backend disappears.

That distinction needs to stay almost repetitive on this page because runtime
inventories are one of the easiest places for "present" to get confused with
"owned."
The runtime can visibly contain more and more machinery while the most
important sentence still begins with:

- "in practice the operator knows..."
- "privately we know that node X is the real one..."
- "normally that fallback only works because we remember..."

If this page stops making those unfinished sentences visible, then the
inventory is laundering presence into adulthood.

## What this page is and is not allowed to prove

This page is authoritative about:

- what the root stack includes today
- which networks and service domains are actually active
- which helper layers are already present in the priority implementation
- which visible implementation facts shape the real documentation burden

This page is not authoritative about:

- generic wrong-node success
- durable peer-forward routing during backend loss
- cross-node policy preservation
- stateful correctness under failover
- promoting planning material into live truth

This is the strongest current inventory.
It is not a distributed-systems completion certificate.

It also cannot be allowed to become a comfort certificate.
The user's complaint is not that too little exists.
The complaint is that a lot already exists and too much of the important truth
can still collapse back into private operator explanation once the bad day
arrives.

That sentence should govern the whole page.
If the inventory makes the reader feel calmer without making the hidden
operator burden materially narrower, the inventory has started performing
reassurance instead of preserving truth.

## Strongest honest current answer

The current runtime is already serious enough that shallow language becomes
dangerous.

It clearly proves:

- Compose is still the real implementation center
- the root file is not symbolic or legacy residue
- the stack already spans public ingress, private mesh, observability,
  operator tools, app hosting, TCP services, documentation, AI tooling, and
  route experimentation
- the repo has enough machinery that words like "failover" or "anti-SPOF"
  must be used carefully

It still does **not** prove the hardest missing truths:

- shared current placement truth
- durable wrong-node rescue routes
- semantically valid peer selection
- preserved auth and middleware after peer handoff
- honest stateful failover semantics

That difference is the whole reason this page exists.

The runtime is already rich enough to tempt the reader into the wrong
conclusion:

- the stack is broad
- the edge is serious
- the helpers are real
- therefore the options must be becoming real too

This page exists to stop that jump.

It exists because this repo is already past the stage where absence is the
main problem.
The danger now is partial presence being socially overread as one more burden
having moved when it really has not.

More bluntly:
"look how much is already here" is one of the main ways the ecosystem keeps
dressing up partial answers as if they had finally become believable options.

## What still does not count as runtime evidence

This page is where documentation inflation becomes especially tempting.

The following still do not count as strong runtime evidence for the user's real
goal:

- a large number of active services
- a sophisticated-looking edge stack
- multiple fragments and networks
- TCP labels on stateful services
- observability breadth
- helper presence with failover-shaped names
- an any-node-looking ingress story that still depends on remembered backend
  ownership
- a route that can be rendered today but has not been stressed under the
  failure that would make the route matter

Those facts prove the repo is serious and operationally dense.
They still do not prove that the system, rather than the operator, owns the
truth needed on the bad day.

That difference is not pedantry.
It is the whole emotional content of the project.
The user is not asking for more evidence that the stack is substantial.
The stack is already substantial enough to feel insulting when the final truth
still lives in private operator memory.

That is why a runtime inventory in this repo has to do more than enumerate.
It has to keep saying which part of the runtime is real, which part is merely
adjacent to the dream, and which part still leaves the operator acting as the
missing control plane.

That is why partial presence is more dangerous here than simple absence.
Absence tells the truth faster.
Partial presence is how a stack starts sounding complete before it has actually
removed the humiliation threshold the user keeps pointing at.

## What the operator still has to know privately today

This page also has to say, without softening it, what the current runtime
still forces a human to know personally.

Today the operator may still need to know things like:

- which node currently hosts the real copy of one named HTTP service when the
  request lands somewhere else
- whether a route generated by helper logic is merely present right now or
  would still survive after the preferred backend disappears
- whether a reachable peer is only alive on the mesh or actually acceptable
  for the requested route's auth, middleware, secrets, and revision surface
- whether a TCP endpoint exposed through Traefik is merely transport-reachable
  or actually safe to describe with failover language
- whether a state-bearing surface such as MongoDB, Redis, or Headscale still
  hides a sacred authority node behind a more distributed-looking ingress story
- whether a request that reaches the "right" backend after forwarding is still
  the same protected service in terms of middleware, headers, policy, and
  user-visible semantics

That list is a more faithful summary of the user's complaint than "the stack is
still maturing."

It is also the shortest honest answer to what the runtime still expects one
human to quietly finish for it.

The user is not saying the runtime lacks machinery.
They are saying too many decisive truths still need to be remembered,
inferred, or reconstructed by a human at the exact moment a believable
platform should already be exposing them.

That is why runtime breadth cannot be narrated as relief on its own.
If the operator still has to privately complete the sentence
"yes, but what this node should really do right now is...",
then the runtime is still borrowing adulthood from human memory.

## What a runtime-backed progress packet would have to contain

Before this page is used to support stronger "the runtime is moving toward the
dream" language, it should point to a concrete packet.

That packet should include:

- the exact live component or path being credited
- the narrower burden it really reduces
- the artifact in the root runtime that proves the component is materially
  present now
- the adjacent burden it still does not reduce
- the next drill or comparison required before stronger language is allowed

Without that packet, inventory tends to impersonate relief.

That line should be read literally.
Inventory is one of the oldest ways infra docs accidentally overpay
confidence:

- they show how much is present
- presence gets mistaken for burden transfer
- burden transfer gets narrated before it exists

There is a harsher way to say the same thing:
inventory is where a complicated stack starts trying to cash out dignity it
has not yet earned.
This knowledgebase should stay suspicious of that move on purpose.

The anti-slop question for every subsection on this page is therefore:

> what exact truth did the runtime gain here, and what exact bad-day sentence
> would the operator still have to privately complete anyway?

## The root file is still a major infrastructure surface

The root file is not just a launcher for fragments.
It still owns:

- networks
- configs
- secrets
- direct services
- operator-facing utility surfaces
- some of the most sensitive mixed-protocol workloads in the repo

That matters because one of the easiest bad reads in this repository is:

> the real system must already live somewhere else because the stack has been
> split into fragments.

That read is false.

The root `docker-compose.yml` is still a major truth surface.
Any future middle layer or orchestrator promotion has to respect that instead
of treating the root file as a leftover artifact.

That matters because the user is explicitly resisting systems that demote the
last readable surface before proving the new hidden surface has earned the
trust being demanded.

One of the user's core complaints is that many ecosystems demand surrender of
the last legible layer before proving the new hidden layer is actually more
honest.
The root Compose surface still matters because it lets the operator inspect a
large part of reality without being forced to take a distributed fairy tale on
faith.

## Root-owned networks

The root runtime defines three central networks:

| Network | Immediate visible role | Why it matters | What it still does not prove |
| --- | --- | --- | --- |
| `warp-nat-net` | controlled alternate egress and routing behavior | the repo is already experimenting with network truth and non-default egress | that cross-node route decisions are cluster-owned |
| `publicnet` | ingress-adjacent traffic and externally exposed surfaces | public-facing behavior is already explicitly partitioned | that all public nodes share the same request-time service truth |
| `backend` | internal app and support traffic | the stack clearly distinguishes internal from public traffic | that cross-node placement and peer eligibility are system-owned |

These networks prove segmentation pressure.
They do not prove a shared multi-node control plane.

## Active include graph

The root file directly includes the following fragments:

| Fragment | Primary domain | Why it matters |
| --- | --- | --- |
| `compose/docker-compose.coolify-proxy.yml` | public ingress and edge policy | this is the strongest live edge surface: Traefik, CrowdSec, TinyAuth, `nginx-traefik-extensions`, `cloudflare-ddns`, `docker-gen-failover`, Autokuma, log rotation |
| `compose/docker-compose.docs.yml` | documentation serving | the docs themselves live inside the same routed stack they are describing |
| `compose/docker-compose.firecrawl.yml` | browser and crawl workloads | proves the repo is already carrying worker-like and queue-like complexity |
| `compose/docker-compose.headscale.yml` | private mesh and coordination assumptions | Headscale is not theoretical; the runtime depends on it as part of the private-network story |
| `compose/docker-compose.llm.yml` | AI and model-gateway workloads | the runtime includes higher-complexity app surfaces beyond ordinary self-hosting |
| `compose/docker-compose.metrics.yml` | observability | Grafana, Prometheus, VictoriaMetrics, Loki, Promtail, Blackbox, cAdvisor, exporters, dashboards, alerts |
| `compose/docker-compose.stremio-group.yml` | media and streaming-adjacent workloads | large heterogeneous service pressure remains part of the real stack |
| `compose/docker-compose.warp-nat-routing.yml` | alternate-routing and egress experiments | networking behavior is part of the repo's live design surface, not an afterthought |
| `compose/docker-compose.wishlist.yml` | smaller absorbable app surfaces | the stack keeps growing by absorbing services into the same core runtime |

This include graph matters because it proves the repo is already operating as a
platform-shaped Compose system.
The missing problem is not "how do we stop being toy-scale?"
The missing problem is "how do we stop being platform-shaped without lying
about distributed truth?"

That second sentence is the more important one.
The repo is not begging to become serious.
It is already surrounded by serious-sounding options.
What it still lacks is an option that remains believable after wrong-node
entry, backend loss, and stateful consequences stop being hypothetical.

That distinction should stay central:
the repo is already sophisticated enough to deserve serious language.
It is still not allowed to spend that seriousness as if sophistication itself
were a substitute for shared current truth.

That is the real threshold.
Not "does the stack resemble something platform-like?"
Not "does it contain components respected by serious operators?"
But "when locality fails or the request lands on the wrong node, does the
system itself carry enough truth to avoid turning the operator back into the
private interpreter of reality?"

## The live edge stack already present

The edge fragment already includes:

- `traefik`
- `crowdsec`
- `crowdsec-init`
- `tinyauth`
- `nginx-traefik-extensions`
- `cloudflare-ddns`
- `docker-gen-failover`
- `whoami`
- `logrotate-traefik`
- `autokuma`

That list is important because it proves the repo already has:

- L7 ingress
- security and bouncer logic
- forward-auth pressure
- DNS automation pressure
- dynamic routing pressure
- service probing and uptime pressure

It also proves why lazy documentation is dangerous.
Once a stack contains this much edge machinery, it becomes very easy to
accidentally narrate "serious ingress" as "solved cross-node request
preservation."

This runtime does not justify that leap.

That inflation pressure is exactly why this page cannot stop at component
enumeration.
When readers see Traefik, TinyAuth, CrowdSec, DDNS, docker-gen, probes, and
middleware, they naturally start filling in missing adulthood on the stack's
behalf.
This page has to keep refusing that generosity.

That leap is one of the main things the knowledgebase is trying to outlaw:
serious ingress is not the same thing as request-preserving truth ownership.

## The live mesh surface already present

The Headscale fragment already contains:

- `headscale-server`
- `headscale`
- inline Headscale configuration
- Traefik-facing route and middleware material around the Headscale surfaces

The planning docs explicitly record that Headscale is single-node today and
still a singleton control-plane concern.

So the correct reading is:

- Headscale is real and central to the repo's private-mesh assumptions
- Headscale is not yet a proven HA control-plane surface

This is exactly the kind of distinction the docs must keep visible.

## The live observability surface already present

The metrics fragment is huge, and that matters.

It already includes:

- `victoriametrics`
- `prometheus`
- `grafana`
- `loki`
- `promtail`
- `blackbox-exporter`
- `cadvisor`
- `alertmanager`
- `process-exporter`
- `mongodb-exporter`
- `redis-exporter`
- many dashboards, alert rules, and recording rules

That proves the repo is not missing seriousness.
It is also a warning:

> strong observability does not automatically mean strong distributed truth.

The system can be richly monitored and still fail the user's real benchmark if
wrong-node forwarding, placement truth, or stateful authority remain socially
reconstructed.

Observability can explain the pain after the fact while still doing nothing to
move who had to privately know the truth before the fact.
That is exactly the kind of impressive-but-insufficient option the user keeps
rejecting.

## Root-owned direct services

The root file still declares many direct services rather than delegating
everything to fragments.

Representative examples include:

- data and state:
  - `mongodb`
  - `redis`
- app and utility surfaces:
  - `dcef`
  - `chat-analytics`
  - `searxng`
  - `code-server`
  - `session-manager`
  - `bolabaden-nextjs`
  - `telemetry-auth`
  - `dns-server`
- operator and maintenance surfaces:
  - `dozzle`
  - `homepage`
  - `portainer`
  - `watchtower`
  - `dockerproxy-ro`
  - `dockerproxy-rw`
- mixed-protocol and project-specific surfaces:
  - `biodecompwarehouse`
  - `biodecompwarehouse-mcp`
  - `biodecompwarehouse-bsim-server`
  - `biodecompwarehouse-aio`

That service list matters for three reasons.

First, the root file is still an active authoring surface.

Second, the stack already spans:

- stateless HTTP
- TCP and database endpoints
- operator dashboards
- mixed app surfaces
- documentation surfaces
- niche project services

Third, every claim about "one unified platform" has to survive all of those
service classes, not just the easiest HTTP examples.

That is why this page cannot be satisfied by sounding comprehensive.
The user has already seen too many comprehensive stories that quietly borrow
their confidence from the easiest route class and then spread that tone over
everything else.

## Why the runtime already creates documentation pressure

This runtime is complicated enough that simple summaries become misleading.

A useful docs pass now has to answer at least these questions:

- which services are local-only concerns and which participate in a broader
  any-node story?
- which services are routed through Traefik and which expose TCP directly?
- which services are state-bearing enough that failover language must stay
  harsh?
- which helpers are active but untrusted?
- which routes are generated, which are static, and which are only planned?

Without those distinctions, the docs can sound coherent while quietly
flattening away the user's whole complaint.

That complaint is not abstract.
It is the repeated experience that "we have a stack for that" keeps turning
into "you still have to be the actual cross-node interpreter for that."

## What this runtime already rules out

The current runtime is already strong enough to rule out a few lazy readings.

It rules out:

- "the repo is still basically a simple homelab stack"
- "the real problem must just be adding more services"
- "the docs can stay high-level because the implementation is still small"
- "the root Compose surface is only scaffolding for some future real platform"

It does **not** rule in:

- generic wrong-node success
- durable backend-loss recovery
- policy-preserving peer handoff
- stateful HA that survives contact with write authority

That split is worth stating plainly because it is where many infrastructure
writeups become emotionally useless.
They are willing to admit the stack is complicated, but they still phrase the
conclusion as if complication itself were already a respectable substitute for
the missing options.

This page should keep refusing that substitution.

## The most important negative evidence already named by the repo

The current runtime already contains negative evidence that the docs must keep
visible instead of smoothing over:

- `docker-gen-failover` is present but already documented elsewhere as losing
  routes under the exact failure condition where rescue should matter
- Headscale is materially live but openly treated as effectively singleton
  today
- stateful surfaces are present in the root graph, but presence does not
  settle write authority, promotion, or replica semantics
- no live shared placement-truth artifact in the root runtime has yet earned
  stronger routing language

This matters because the current stack is not only defined by what is present.
It is also defined by which present helpers and surfaces still fail the user's
benchmark under stress.

The live runtime has to be read next to the planning material, because the
planning layer records several crucial gaps that affect how the runtime should
be described.

The strongest named gaps include:

- `docker-gen-failover` can delete routes when containers stop
- `watchtower` is configured but has not been functioning correctly
- Cloudflare DDNS presence is not the same thing as full multi-node request
  failover
- multi-record DDNS behavior is still a known problem
- secret sync is still manual
- compose sync is still manual
- automated service failover between nodes is still missing
- a live root `services.yaml` registry is still unproven

Those gaps matter because they keep the runtime honest.

The stack is already rich.
The repo itself still says the middle control-plane truths are incomplete.

That incompleteness is not a minor implementation detail.
It is the reason the runtime can already look more option-rich than it really
is under stress.

## What the runtime clearly proves

The current root runtime clearly proves:

1. the repo is Compose-first in reality, not just in rhetoric
2. the root file and include graph together form a large active platform
3. the system already carries serious ingress, mesh, observability, and
   operator complexity
4. the docs have to distinguish service classes, not just list containers
5. the repo is solving a hard problem from inside a real stack, not from a toy
   lab

## What the runtime still does not prove

The current root runtime does **not** yet prove:

1. any healthy public node can generically preserve a request for a remote
   service
2. a live placement registry such as `services.yaml` exists and is actively
   consumed
3. the rescue route survives the failure that makes rescue necessary
4. auth and middleware remain semantically identical after peer handoff
5. TCP and stateful services have cluster-grade authority and failover truth
6. the operator can explain cross-node behavior from shared tracked truth
   instead of private reconstruction

That distinction is the key reading rule for the entire site:

the current runtime is already real, but the shared-truth layer is still
incomplete.

If that sentence ever gets softened, the whole site starts telling a smaller
and safer story than the one the user is actually forcing the repo to face.
