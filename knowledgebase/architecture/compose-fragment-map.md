# Compose Fragment Map

This page reconstructs how the priority implementation is actually assembled
from the root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and the active files under
[`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/).

It is not here to praise modularity.
It is here to answer the harder question:

> when someone says "the stack," which files are carrying present-tense runtime
> burden, which files are carrying intent or experiments nearby, and where
> would a reader accidentally over-credit the tree as already solving the
> multi-node no-Swarm failover dream?

That matters because this repo contains several classes of Compose truth at the
same time:

- the live Compose-first runtime
- side-path Compose futures
- generated or aggregate Compose views
- optional and parked stacks
- planning and research pressure aimed beyond the live runtime

If those classes are blended together, the tree starts sounding like it offers
more real choices than it honestly does.
That is one of the exact wounds this knowledgebase is supposed to stop hiding.

## What this page is and is not allowed to prove

This page is allowed to prove:

- where the priority runtime is authored right now
- which fragments are part of the active root include graph
- which burdens the root file still owns directly
- which nearby Compose files are real but not part of the live assembly path
- where breadth of files most easily gets mistaken for breadth of honest
  options

This page is not allowed to prove:

- that wrong-node routing already works generically
- that file partitioning equals shared placement truth
- that serious edge language equals preserved request semantics
- that helper-generated fallbacks equal peer-forward dignity
- that nearby Compose options are already earned runtime choices

## The strongest honest current answer

The live runtime is still decisively Compose-first.
The root file is still the canonical human assembly surface.
The fragment set is large enough that the repo already carries distributed
systems burdens, but not yet large enough to claim distributed truth ownership.

That is the right reading.

Not:

- `this is just a simple Compose repo`

and not:

- `this is already a distributed control plane with only a few gaps`

The root runtime is substantial.
The missing middle truth layer is substantial too.

## The authority ladder for reading Compose files honestly

If the goal is to understand the live runtime, read in this order:

1. root
   [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. the active include targets named by that root file
3. runtime-boundary pages such as
   [current-compose-runtime.md](current-compose-runtime.md)
   and
   [failure-model-and-maturity.md](failure-model-and-maturity.md)
4. only then side-path Compose files, aggregate views, and planning material

If someone starts with `compose/docker-compose.everything.yml`,
`compose/docker-compose.parsed.yml`, or a future-facing fragment before reading
the root file, they are already stepping away from the repo's strongest
present-tense authority surface.

## What the root file still owns directly

The root file is not just a loader.
It still owns real platform language and real services directly.

### Shared network language

The root file defines shared networks including:

- `publicnet`
- `backend`
- `warp-nat-net`

That proves the runtime is already modeling different traffic classes
explicitly.
It does not prove that any healthy receiving node already owns the full meaning
of the request when locality is absent.

### Shared config and secret language

The root file also carries shared config and secret surfaces such as:

- secret `signing_secret`
- config `watchtower-config.json`
- session-manager assets
- inline Homepage configuration blobs

That proves Compose is being used as a real authoring and policy surface.
It does not prove shared placement truth, peer eligibility truth, or end-to-end
fallback correctness.

### Root-owned services

The root file still directly defines a serious workload set, including:

- `mongodb`
- `redis`
- `dcef`
- `chat-analytics`
- `searxng`
- `code-server`
- `homepage`
- `watchtower`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `dozzle`
- `portainer`
- `dns-server`
- `telemetry-auth`
- `bolabaden-nextjs`
- `biodecompwarehouse`
- `biodecompwarehouse-mcp`
- `biodecompwarehouse-bsim-server`
- `biodecompwarehouse-aio`
- profile-gated `session-manager`

This matters because any future promotion into a stronger middle layer has to
account for the fact that the root file is not leftovers.
It is still one of the main places where actual infrastructure truth is being
written by hand.

## The active include path

The root file currently includes these fragments:

1. `compose/docker-compose.coolify-proxy.yml`
2. `compose/docker-compose.docs.yml`
3. `compose/docker-compose.firecrawl.yml`
4. `compose/docker-compose.headscale.yml`
5. `compose/docker-compose.llm.yml`
6. `compose/docker-compose.metrics.yml`
7. `compose/docker-compose.stremio-group.yml`
8. `compose/docker-compose.warp-nat-routing.yml`
9. `compose/docker-compose.wishlist.yml`

That list is not directory trivia.
It is the current answer to:

> which fragments are actually part of the priority implementation that still
> centers on `docker-compose.yml`?

Everything else in `compose/` has to be read with lower authority unless it is
explicitly promoted.

## Why the fragment map matters to the real architecture problem

The user's frustration is not just "there are too many files."

The real frustration is:

- there are many visible layers
- many of them look like they could solve the problem
- many of them really do solve some narrow local problem
- but the overall wrong-node, fallback, and truth-ownership question is still
  easy to leave half-private

So the fragment map has to answer not only:

- what files are active?

but also:

- what kind of burden does each active fragment visibly own?
- what nearby overclaim does that burden tempt people into making?
- what private operator sentence still survives after the fragment is read?

## Fragment burden ledger

Read the fragment set as a burden ledger, not as a directory tour.

| Active fragment | It visibly owns | What that means in plain language | What the docs must still refuse |
| --- | --- | --- | --- |
| `compose/docker-compose.coolify-proxy.yml` | public ingress, auth, middleware, DDNS, edge security, helper-driven failover generation | this is where the repo most strongly tries to turn many public nodes and many routed services into one coherent front door | that sophisticated ingress already equals generic wrong-node dignity, peer correctness, or auth continuity under failure |
| `compose/docker-compose.docs.yml` | docs serving inside the same routed platform | the documentation surface is itself subject to the same ingress truth limits it describes | that clean docs hosting means the harder routes are system-owned |
| `compose/docker-compose.firecrawl.yml` | browser automation, worker coordination, queue-backed support state | this fragment proves the runtime is carrying real worker/state complexity, not just static apps | that queue and worker seriousness automatically upgrades stateful failover truth |
| `compose/docker-compose.headscale.yml` | private mesh identity and cross-node reachability assumptions | this is part of the repo's attempt to make nodes mutually reachable without pretending that reachability is enough | that node identity or mesh reachability equals valid forwarding truth |
| `compose/docker-compose.llm.yml` | AI gateway, tool exposure, caches, search/vector/state auxiliaries | this fragment widens the runtime into app-gateway and memory-bearing workloads that make distributed truth even more important | that an app gateway or proxy layer solves node-level request preservation |
| `compose/docker-compose.metrics.yml` | metrics, probes, logs, dashboards, alerts, blackbox checking | this fragment proves the repo takes observability seriously and wants visible evidence instead of pure folklore | that observability breadth equals automated recovery, shared placement truth, or verified failover |
| `compose/docker-compose.stremio-group.yml` | media and helper workload cluster with varied external dependencies | this fragment shows how quickly workload sprawl makes one-human truth ownership brittle | that large workload breadth creates coherent placement semantics by itself |
| `compose/docker-compose.warp-nat-routing.yml` | alternate egress shaping and packet-path control | this fragment proves the repo is exploring path-control and egress ownership, not just basic app hosting | that packet-path cleverness equals preserved application meaning |
| `compose/docker-compose.wishlist.yml` | a narrow lightweight app surface | this fragment provides a simpler active include among harder fragments | that an easy stateless app path proves the hard multi-node burdens are solved |

## What the active fragment set says about the repo

The active fragment set already says something more specific than "large Compose
repo."

It says the priority runtime is simultaneously carrying:

- edge and identity pressure
- private mesh pressure
- observability pressure
- state and worker pressure
- media/app sprawl pressure
- egress and network-path pressure

That is why the repo no longer feels like simple container hosting.
The operator is already being asked to reconcile the same classes of burden that
usually force a system to either:

- expose shared truth clearly

or:

- quietly make one person keep stitching the truth together by hand

The current tree proves the burden set.
It does not prove that burden transfer has actually happened.

## Root versus fragment: where the pressure really lives

One easy mistake is to imagine the root file as mere bootstrap and the fragments
as the "real stack."
That is not the honest reading.

The root file still owns:

- shared networks
- shared config and secret naming
- direct stateful services such as `mongodb` and `redis`
- operator-facing utilities such as `code-server`, `dozzle`, and `portainer`
- docs and site-adjacent services such as `bolabaden-nextjs`, `homepage`, and
  `telemetry-auth`

The fragments then widen the domain footprint from there.

So the better reading is:

- the root file is still a major workload and policy surface
- the fragments widen the platform into more burden classes
- the missing truth-owning middle layer would have to reconcile both

It does not get to ignore the root just because the fragments look more modern
or more distributed.

## Important side-path Compose files that are not in the active include path

The repo also contains meaningful Compose files that are not currently part of
the live root include graph.
They matter, but they do not carry the same authority.

| Side-path file | Why it matters | Why it becomes dangerous if overread |
| --- | --- | --- |
| `compose/docker-compose.authentik.yml` | preserves a real alternate identity direction | it can be mistaken for present auth authority when TinyAuth is the visible active edge path |
| `compose/docker-compose.l4-ingress.yml` | preserves explicit thought about L4 and raw TCP ingress | it can be mistaken for active runtime proof of TCP failover or stateful authority |
| `compose/docker-compose.nomad.yml` | proves scheduler exploration is real, not imaginary | it can be mistaken for an already-earned control-plane decision |
| `compose/docker-compose.core.yml` | shows alternate grouping instincts | it can be mistaken for the present canonical assembly path |
| `compose/docker-compose.warp.yml` and `compose/docker-compose.vpn-docker.yml` | preserve alternate egress and network approaches | they can be mistaken for the active WARP routing path |
| `compose/docker-compose.parsed.yml`, `compose/docker-compose.semiparsed.yml`, `compose/docker-compose.semifullparsed.yml`, `compose/docker-compose.everything.yml` | expose aggregate or transformed views of the stack | they can look more complete than the canonical runtime and tempt people into treating generated breadth as authority |
| `compose/docker-compose.unused.yml`, `compose/docker-compose.unsend.yml`, `compose/docker-compose.wordpress.yml`, `compose/docker-compose.plex.yml` | show parked domains, templates, and appetite for more surfaces | they inflate the perceived option space without earning present runtime authority |

## The three most dangerous reading mistakes

### Mistake 1: file count as option-space proof

`There are many Compose files` does not mean `the user now has many mature
answers`.

It may only mean:

- the repo has many experiments
- the repo has many burden classes
- the repo has many future directions
- the repo has many places where private operator truth can hide

### Mistake 2: generated view as canonical truth

Aggregate files can look more complete than the root because they are flatter
or more exhaustive-looking.

That does not make them the real human control surface.

The user explicitly cares that the priority implementation still centers on the
real `docker-compose.yml` path, not on some prettier synthetic summary.

### Mistake 3: edge sophistication as multi-node dignity

The most seductive overread in this repo is:

- Traefik is serious
- DDNS is serious
- helper-generated failover config is serious
- mesh and observability are serious
- therefore the platform must be close to honest wrong-node behavior

That does not follow.

A serious edge still may depend on one person privately knowing:

- which backend is really current
- which peer is really eligible
- which fallback path really survives failure
- which protected route really keeps its meaning during handoff

## What this page should leave behind

After reading this page, the reader should be able to say:

- which files are carrying the live priority runtime
- which files are nearby but lower-authority
- which burden classes are already present
- which overclaims each fragment most tempts people into making
- why the repo can be both substantial and still incomplete in truth ownership

The reader should not leave saying:

- `there are lots of files, so there must be lots of real answers`
- `the fragments look modular, so the architecture must be nearly solved`
- `the generated outputs look complete, so the root file is mostly historical`

## The shortest honest summary

The fragment tree proves that `bolabaden-infra` is no longer struggling with
toy problems.
It is already carrying the kinds of edge, mesh, state, observability, and
egress burdens that make a missing control plane impossible to hide forever.

What it still does not prove is that those burdens have been transferred out of
one operator's private head and into system-owned truth.
