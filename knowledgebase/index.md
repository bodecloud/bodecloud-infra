# bolabaden Infrastructure Knowledgebase

This site exists to answer one unusually stubborn infrastructure question:

> how do you keep
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> as the real human control surface, spread services across several ordinary
> Docker nodes, and still make wrong-node traffic, fallback, and anti-SPOF
> behavior feel like one coherent platform instead of one operator privately
> remembering the real answer?

That is the real subject of `bolabaden-infra`.

It is also the question every page in this site should be forced to serve.

If a page is informative but answers a smaller question than that one, the page
is still failing the user.

## The failure mode this site has to prevent

The failure mode is not:

- the docs are too short
- the docs are too scattered
- the docs are too technical

The failure mode is:

1. the site sounds organized
2. the site sounds mature
3. the user accusation gets softened into a calmer neighboring question
4. the hidden operator burden becomes easier to forget

That is how a documentation set becomes polished and useless at the same time.

This site is not primarily:

- a generic self-hosting handbook
- a reverse-proxy cookbook
- an orchestrator comparison catalog
- a generic `modern homelab` notebook

Those topics only matter when they help answer the real question above.

## The user's dream, stated as directly as possible

The dream is not just:

- "run Docker on more than one machine"
- "have multiple ingress points"
- "add enough glue that the stack feels clustered"
- "avoid Kubernetes because it is annoying"

The dream is:

- keep `docker-compose.yml` close to the center of human authorship
- keep the nodes ordinary enough that the system is still understandable
- avoid inventing one new sacred node to replace the old one
- avoid inventing one new sacred human to replace the old one
- make any healthy public node a believable first hop
- preserve the meaning of the request even when it lands on the wrong node
- keep failover and fallback from degenerating into private folklore
- only pay heavier orchestration cost if it actually kills one of those hidden
  burdens

This matters because many infrastructure summaries quietly rewrite the dream
into something easier:

- generic high availability
- generic clustering
- generic orchestration selection
- generic self-hosting maturity

Those are not the same thing.

## The private-sentence benchmark for the whole site

The whole site should be read under a simple benchmark:

> after reading this page, what exact private sentence is still being finished
> by the operator?

Examples:

- `I still personally know which node is the real one.`
- `I still personally know which peer is safe.`
- `I still personally know whether the fallback is real.`
- `I still personally know whether the route that answered still means the same thing.`
- `I still personally know whether the stateful answer is authoritative.`

If a page cannot make that clearer, it may still be informative while still
being weaker than the user needs.

## The accusation this site has to preserve

The site only stays useful if it keeps the user's real accusation visible:

> there seem to be endless options for multi-node Docker, failover, overlays,
> service discovery, ingress, and orchestration, but too many of them solve
> one visible layer and then quietly leave the operator as the hidden control
> plane when reality gets sharp.

If a page becomes more polished while losing that accusation, it got worse.

The accusation should stay active while reading every page:

> does this explanation still depend on a human privately knowing which node is
> special, which peer is safe, which backend is current, or which fallback is
> fake?

If the answer is yes, then the docs are not yet describing a system-owned
truth.

## Why the accusation matters more than tone

Many infra docs get rewarded for sounding:

- even-handed
- balanced
- practical
- calm

Those traits are fine until they erase the accusation.

In this repo, the accusation is the integrity check.
If the page sounds nicer by making the wound easier to forget, the page got
worse.

## The shortest correct reading

The strongest intent surface in the repo is
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It says the intended direction is:

- no central orchestrator by default
- services manually assigned to nodes
- current-state truth preferred over scheduler-declared desired state
- local-first serving when the requested service already lives on the receiving
  node
- peer-forward fallback when the receiving node is healthy but the service is
  remote
- explicit separation between HTTP routing and raw TCP or stateful behavior
- anti-SPOF pressure without fake HA language

That target request contract is:

```text
User -> Cloudflare DNS -> any surviving public node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that really hosts it
```

That is the dream.
It is not the same thing as current proof.

That gap between dream and proof is not a documentation nuisance.
It is the whole architecture problem.

## What this site must keep forcing the reader to remember

Three things must stay simultaneously visible:

- the dream is specific
- the runtime is real
- the hidden truth-owning layer is still incomplete

If one of those disappears, the reader usually falls into one of three wrong
stories:

- only aspiration, no implementation
- only implementation, no accusation
- only seriousness-of-stack, no missing burden transfer

## The three-part checksum for the whole site

Everything in this knowledgebase should keep this sentence intact:

1. the dream is specific
2. the runtime is real
3. the truth-owning middle layer is still incomplete

If a page loses one of those three, it starts lying in one direction or
another.

The most common failure modes are:

- the dream disappears and the page answers only a calmer neighboring question
- the runtime disappears and the page becomes architecture theater
- the missing middle disappears and the page sounds like the problem is mostly
  solved because the stack looks serious

Those failures are why this site has to be read almost like a RAG retrieval
system with strict source boundaries, not like an ordinary documentation site.

## What the repo already proves

The priority implementation is still rooted in:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

That runtime already proves a serious Compose-first platform with:

- a real Traefik and auth-bearing edge layer
- public and private network segmentation
- stateful services such as MongoDB, Redis, PostgreSQL, RabbitMQ, and Qdrant
- observability and operator surfaces
- Headscale and mesh pressure
- alternate routing and egress experiments

What it still does **not** prove is the thing the user actually cares about:

- that any healthy public node can accept a request and preserve it correctly
  when the service is remote
- that placement truth is shared explicitly instead of remembered
- that peer eligibility is system-owned rather than guessed
- that fallback paths survive the failure that makes fallback matter
- that auth, middleware, and request semantics survive peer handoff
- that stateful services are genuinely resilient rather than merely reachable

The runtime is serious.
The missing burden transfer is still brutally concrete.

## What a fake good summary sounds like

These still sound good while being too small:

- `the repo is exploring several HA strategies`
- `the repo is evolving toward a middle layer`
- `the repo has a solid multi-node foundation`
- `the repo is deciding between lightweight and heavyweight orchestration`

Each one becomes more truthful only if it also preserves:

> the operator is still the missing control plane on the bad day

## How to read this site without making it useless

Before trusting any section, ask four questions:

1. what exact question is this page answering?
2. what smaller neighboring question would be easier to answer but less useful?
3. which truth layer is actually carrying the answer?
4. what stronger sentence is still forbidden when the page ends?

If the page does not leave those answers behind, it is still too soft.

## The hidden operator job the site has to keep naming

The shortest honest summary of the current wound is:

the operator is still acting like the missing control plane.

That hidden job currently includes things like:

- remembering what runs on which node right now
- remembering which peer is actually valid right now
- remembering whether a generated fallback route would still exist under real
  failure
- remembering whether a reachable answer is only transport-reachable or
  semantically valid
- remembering which stateful surfaces still hide a sacred authority node

If a page forgets that hidden job, it will almost always overstate the system.

That is why the knowledgebase has to keep circling back to the same harsh
point:

> the current wound is not lack of components, but lack of system-owned truth
> during the bad day

## The four truth layers

This site only stays honest if it keeps these truth layers separate:

| Truth layer | Main examples | What it is allowed to answer | What it is not allowed to answer |
| --- | --- | --- | --- |
| Live runtime truth | `docker-compose.yml`, active `compose/` fragments, `docker compose config` | what the priority implementation actually ships now | what the system would do if missing layers already existed |
| Repo-native intent | `.github/copilot-instructions.md`, `knowledgebase/AGENTS.md` | what the repo is clearly trying to become | whether the runtime already behaves that way |
| Planning truth | `docs/INFRASTRUCTURE_MASTER_PLAN.md`, stateful and ingress plan docs | what gaps are already named and what promotion paths exist | that a planned repair is already active |
| Archive and research pressure | `knowledgebase/research/`, `knowledgebase/source-archive/` | why the user keeps rejecting partial answers and which patterns recur | present-tense runtime proof |

The site gets misleading the moment those four layers blend into one calm
voice.

## Why the site has to act like evidence custody, not normal docs

This knowledgebase is closer to a retrieval-and-proof surface than to a normal
reference site because the main risk is not forgetting facts.

The main risk is illegal confidence transfer:

- from dream to runtime
- from plan to implementation
- from archive pressure to present-tense proof
- from component presence to burden transfer

## The core architecture problem

The repo is not just asking for `better HA`.
It is asking for a smaller honest middle layer between:

- static multi-node Docker glue that still depends on private folklore
- and a heavyweight orchestrator worldview that has not yet earned the right to
  hide that much truth

That middle layer would need to make at least some of these truths
system-owned:

- current placement
- current peer eligibility
- current rescue-route validity
- current route-class meaning
- current explanation for why the system chose local versus remote

If a candidate does not move those truths into the system, it has not solved
the user's real complaint, no matter how mature it sounds.

## What this site should leave a serious reader with

After a good pass through the site, the reader should not merely think:

- `there are many options`
- `the architecture is clearer`
- `the docs are more comprehensive`

They should be able to say:

- which truth is still privately carried
- which source class proved the present runtime
- which source class only proved intent or planning
- which stronger sentence is still forbidden
- what exact artifact or drill would kill the next private sentence

## The pressure already visible in the archive

The archive is consistent about where the real pain lives:

- multi-node Docker without Swarm keeps converging on service discovery as the
  hard missing piece
- distributed HA orchestration keeps showing that fully peer-equal behavior is
  rare without custom glue
- load-balancer and failover alternatives keep showing that many products
  solve only one slice of the wound
- Nomad, k3s, Kubernetes, and related explorations show that stronger control
  planes may help, but they bring a larger worldview that still has to prove
  it kills the right private burden

That is why the user does not feel like they lack product names.
They feel like they lack honest options.

## How to use this site

If you want the shortest route through the knowledgebase, use this as the main
reader path:

1. [Problem, Pressure, and Goals](architecture/problem-and-goals.md)
2. [Current Compose Runtime](architecture/current-compose-runtime.md)
3. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
4. [The Missing Middle Layer](architecture/missing-middle-layer.md)
5. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
6. [Failure Model and Maturity Matrix](architecture/failure-model-and-maturity.md)
7. [Orchestration Options](architecture/orchestration-options.md)
8. [DevOps Runbook](operations/devops-runbook.md)

Use that path when the real question is:

> what is the actual dream, what is already real, what still remains privately
> held, and what would have to change before stronger language becomes honest?

If you want the proof pressure underneath those pages, read:

- [Evidence Ledger](research/evidence-ledger.md)
- [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)
- [Stateful HA Evidence](research/stateful-ha-evidence.md)
- [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)
- [Source Assimilation Index](operations/source-assimilation-index.md)

Use those when the real question is:

> which claims are actually supported by runtime evidence and which are still
> being socially completed by the operator?

If you want the navigation logic and retrieval discipline first, read:

- [Reading Paths](reading-paths.md)
- [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)

## What a good page in this site should leave behind

After reading a useful page here, a contributor should be able to answer:

- what concrete question the page was solving
- which truth layer carried the answer
- which burden still remains operator-owned
- which stronger sentence is still forbidden
- what artifact, drill, or runtime change would make that stronger sentence
  legal

If a page leaves behind only:

- `the architecture is clearer`
- `the options are broader`
- `the platform sounds more mature`

then it is still too weak for this repo.

## What this site must keep refusing

This knowledgebase should keep several false conclusions illegal:

- `there are multiple public nodes, so failover is basically solved`
- `Traefik exists, so wrong-node forwarding must be close`
- `a helper has failover in the name, so the platform owns fallback now`
- `a route can be rendered, so it must survive the failure that made it
  matter`
- `the docs are clearer now, so the platform must be closer to adulthood`

Those are exactly the overreads this site exists to interrupt.

The useful conclusion is not that the repo is simple.
It is that the repo is serious, the dream is specific, and the missing burden
transfer is still brutally concrete.

If the site ever becomes easier to admire than to audit, it has drifted away
from what the user actually asked for.
