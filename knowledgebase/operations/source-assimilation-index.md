# Source Assimilation Index

This page is the retrieval contract for the knowledgebase.

Its job is not just to say "consult sources."
Its job is to stop the docs from sounding broad, careful, and well-organized
while still quietly answering a smaller question than the user actually asked.

In this repository, bad retrieval often looks polished.

It usually sounds like:

- a lot of files were read
- a lot of technologies were named
- a lot of options were summarized
- the user received a calm architecture story
- and the real burden was still edited down into something easier

This page exists to stop that downgrade.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what "actual assimilation" means in this repository
- how source classes should be separated before stronger claims are written
- why retrieval discipline is part of the repo's honesty contract

This page is not allowed to prove:

- that a page has already assimilated the repo just because it cites many files
- that source abundance itself produces better answers
- that retrieval discipline can substitute for runtime proof

This is a retrieval contract page, not a claim that the work has already been
done correctly.

## The dream retrieval has to protect

The user is not asking for:

- a larger option catalog
- a more generic architecture guide
- a smoother self-hosting narrative
- a blended summary of "things related to HA"

The user is asking for a knowledgebase that actually reconstructs:

- what they are trying to make true
- why normal answers keep failing
- which truths are already live in the worktree
- which truths are only intent or planning pressure
- which future layers are genuinely trying to remove hidden burden
- where the docs must stay harsh because state, routing, and failover are still
  incomplete

So "actually RAG this time" in `bolabaden-infra` means:

1. recover the real dream before writing
2. identify what class of claim is being made
3. route that claim to the strongest source class
4. say exactly what that source proves
5. say exactly what it does not prove
6. preserve contradiction instead of ironing it out
7. keep worktree truth above elegant prose for runtime claims
8. use archive pressure to recover the user's complaint without letting the
   archive impersonate runtime evidence

If the docs skip that sequence, they may become longer without becoming more
truthful.

## Strongest honest current answer

The repository has four major evidence classes:

1. architecture-intent and honesty surfaces
2. live root runtime surfaces
3. planning and promotion surfaces
4. archive-pressure surfaces

The biggest recurring retrieval mistake is flattening those four classes into
one blended narrative.

Once that happens, at least one of these becomes false:

- the worktree starts sounding more complete than it is
- the user's actual complaint gets replaced by a neighboring, weaker question

This page exists to keep both failure modes visible.

## What still does not count as real retrieval here

The following still do not count as "actually RAG this time" in this repo:

- reading many files without distinguishing source class
- blending archive pressure into runtime proof
- treating planning language as if it outranks the worktree
- producing a calm answer that edits the wound down into a smaller question
- summarizing related technologies without reconstructing the user's benchmark
- removing contradiction because it makes the docs cleaner

That kind of retrieval can sound thorough while still answering the wrong
question.

## Source class 1: architecture-intent and honesty surfaces

These are the strongest sources for:

- what the repo wants the runtime to do
- what the dream actually is
- what kinds of fake-HA language the repo already refuses

Read in this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
4. [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

These files do **not** have equal weight.

Their real hierarchy is:

- `copilot-instructions.md`
  - strongest statement of the multi-node, no-Swarm/no-Kubernetes-by-default,
    local-first, peer-forward dream
- `README.md`
  - strongest repo-level honesty wall about what the live stack still does not
    prove
- `AGENTS.md`
  - strongest operator/setup/runtime constraints surface
- `.cursorrules`
  - strongest authoring-discipline and Compose-hygiene surface

If a page treats those four as if they are all saying the same thing, that page
has already started flattening the repo.

## Source class 2: live root runtime surfaces

These are the strongest sources for:

- what the priority implementation actually is
- what is active in the root stack now
- which helpers and services are already part of the live runtime

Start with:

1. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. active fragments under
   [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

High-value current fragments include:

- `compose/docker-compose.coolify-proxy.yml`
- `compose/docker-compose.docs.yml`
- `compose/docker-compose.firecrawl.yml`
- `compose/docker-compose.headscale.yml`
- `compose/docker-compose.llm.yml`
- `compose/docker-compose.metrics.yml`
- `compose/docker-compose.stremio-group.yml`
- `compose/docker-compose.warp-nat-routing.yml`
- `compose/docker-compose.wishlist.yml`

These surfaces outrank planning docs whenever the claim is:

- this service is live now
- this network is part of the priority runtime now
- this helper already exists now
- this domain is really part of the root implementation now

They do **not** outrank planning docs for claims like:

- this helper is already trustworthy under failure
- this route survives backend disappearance
- this service is now honestly anti-SPOF

## Source class 3: planning and promotion surfaces

These are the strongest sources for:

- what capability gaps the repo has already named
- what failure signatures are already known
- what future layers are being considered for promotion

Start with:

1. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
2. [Infrastructure Master Plan Evidence](../research/infrastructure-master-plan.md)
3. [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
4. [Stateful HA Evidence](../research/stateful-ha-evidence.md)
5. [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
6. [OpenSVC Ingress HA](../research/osvc-ingress-ha.md)
7. [Nomad Exploration Evidence](../research/nomad-exploration-evidence.md)
8. [Garden k3s Exploration Evidence](../research/garden-k3s-exploration-evidence.md)

These are the surfaces where the repo already says things like:

- universal wrong-node success is still missing
- a live `services.yaml` registry is still unproven
- `docker-gen-failover` can delete routes when containers stop
- Cloudflare DDNS presence is not the same thing as full failover
- Headscale is single-node today
- automated service failover between nodes is still missing
- stateful HA requires topology-specific authority and replication truth

These sources are allowed to explain what is missing.
They are not allowed to promote themselves into live proof.

## Source class 4: archive-pressure surfaces

These are the strongest sources for:

- what the user is actually rebelling against
- why normal answer families keep sounding fake
- where the missing middle keeps reappearing

Start with:

- [User Intent and Dream](../research/user-intent-and-dream.md)
- [Archive Pressure Patterns](../research/archive-pressure-patterns.md)
- [Evidence Ledger](../research/evidence-ledger.md)
- the source archive conversations under
  [`knowledgebase/source-archive/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive)

This class is essential because the repo is not only technical.
It is also a reaction against an ecosystem pattern:

> many partial tools exist, and docs often narrate their existence as if the
> operator's real burden has already been removed.

Archive pressure restores that complaint.
It does **not** prove runtime correctness.

## Claim router

Use this table before writing any page or paragraph that sounds confident.

| If the sentence is really claiming... | Start with | Why | It still must not imply... |
| --- | --- | --- | --- |
| "this is the dream of the repo" | `copilot-instructions.md`, `README.md` | strongest architecture-intent and honesty surfaces | that the runtime already earned it |
| "this exists in the live stack today" | root Compose runtime | strongest implementation truth | that it survives the relevant failure |
| "the repo already knows this is missing or broken" | planning surfaces | strongest named-gap evidence | that the repair is already live |
| "the user rejects this whole family of answers" | archive-pressure surfaces | strongest statement of the actual complaint | that the complaint itself proves implementation |
| "this service class needs harsher language" | runtime plus stateful evidence pages | strongest combined truth about protocol and state class | that all stateful closure is already solved |

If a sentence cannot be routed through this table, it is probably still too
vague to trust.

## What a real retrieval proof packet would need

Before a new page or major rewrite can claim it truly assimilated the repo, the
packet should make recoverable:

- which source classes were used
- which files supplied dream reconstruction versus runtime proof
- which contradictions were preserved instead of flattened
- what the strongest live evidence actually was
- what the next stronger artifact would need to prove beyond the current pass

If the rewrite cannot leave that packet behind, it may still be polished, but
it has not actually assimilated the documentation honestly.

## What retrieval still does not count in this repo

Because the user explicitly asked for the docs to "actually RAG this time,"
this page needs to say what false retrieval still looks like.

These still do not count as good retrieval:

- reading many files but collapsing them into a calmer, smaller question than
  the user actually asked
- retrieving architecture intent and then talking as though it were runtime
  proof
- retrieving runtime presence and then narrating backend-loss or wrong-node
  behavior that was never actually shown
- retrieving archive frustration and then using it as a shortcut to claim the
  repo is already near closure
- retrieving many options and technologies without naming which hidden burden
  each one is supposed to remove

The user is not asking for broad thematic relevance.
The user is asking for burden-faithful reconstruction.

If retrieval does not preserve the exact wound, it is still a downgrade even if
it looks comprehensive.

## What a successful retrieval packet should leave behind

In this repo, retrieval should leave behind more than a summary.

A successful packet should let a later reader recover all of these:

- the real question being answered
- the source classes that were consulted
- the strongest artifact for the runtime-facing portion of the answer
- the strongest archive or intent source for the pressure-facing portion
- the contradiction that was preserved instead of blended away
- the next missing artifact that would be required to upgrade the claim

That final item matters because otherwise retrieval becomes archival fluency
without architectural consequence.

The docs should not merely remember what the user wants.
They should also remember what proof would have to exist before stronger claims
about that dream become honest.

## Retrieval packs for the most common real questions

Broad archive scans produce volume faster than understanding.
These packs exist to keep retrieval deliberate.

### Pack A: "Why does Compose stop feeling honest the moment multi-node matters?"

Use when the real question is:

- why is the user still clinging to Compose?
- why is the user so hostile to static glue?
- why is the problem not just "better YAML"?

Start with:

1. [User Intent and Dream](../research/user-intent-and-dream.md)
2. [Problem and Goals](../architecture/problem-and-goals.md)
3. `source-archive/chatgpt-exports/conversations/docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`
4. `source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`

What this pack proves well:

- the user does not merely dislike Compose syntax
- the real break happens when runtime truth escapes local node boundaries
- service discovery and wrong-node routing are the real pressure points

What it must not be widened into:

- Compose itself is useless
- any remote-proxy or overlay trick automatically solves the burden

### Pack B: "What counts as a fake failover answer in this repo?"

Use when the real question is:

- why are Cloudflare and Traefik still not enough?
- why does the repo keep saying wrong-node dignity is the threshold?
- why is `docker-gen-failover` both central and distrusted?

Start with:

1. [HA, Failover, and Routing](../architecture/ha-failover-routing.md)
2. [Request Path and Failure Walkthrough](../architecture/request-path-and-failure-walkthrough.md)
3. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
4. `source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
5. `source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md`

What this pack proves well:

- first-hop redundancy is not the same thing as preserved request meaning
- dynamic route generation can still fail at the critical moment
- proxy category fit does not equal solved failover semantics

What it must not be widened into:

- all proxy-driven failover is worthless
- the edge stack has no value

### Pack C: "What is the missing middle actually trying to own?"

Use when the real question is:

- why is the repo still searching for a helper layer?
- why is `services.yaml` such a recurring idea?
- what truth would a promoted middle layer need to own?

Start with:

1. [The Missing Middle Layer](../architecture/missing-middle-layer.md)
2. [Infrastructure Master Plan Evidence](../research/infrastructure-master-plan.md)
3. [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
4. `source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
5. `source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`

What this pack proves well:

- the missing middle is a capability shape, not a product label
- service discovery and eligibility truth are the real pressure points
- the repo is trying to avoid both static glue and premature platform capture

What it must not be widened into:

- the repo has already chosen the winning implementation

### Pack D: "Why do stateful services stay under harsher rules?"

Use when the real question is:

- why are Redis, MongoDB, and Headscale treated differently from Dozzle or
  static frontends?
- why does the repo keep refusing broad HA language?

Start with:

1. [Stateful HA and Data](../architecture/stateful-ha-and-data.md)
2. [Stateful HA Evidence](../research/stateful-ha-evidence.md)
3. [Stateful HA Plan](../research/stateful-ha-plan.md)
4. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

What this pack proves well:

- stateful correctness is about authority, replication, and write truth
- route reachability is not the same thing as honest failover
- Headscale is still a singleton control-plane concern today

What it must not be widened into:

- no stateful service can ever be improved before total cluster closure

### Pack E: "What does the current root runtime actually contain?"

Use when the real question is:

- what is real in the priority implementation now?
- which fragments and services are central to the live stack?
- why is the root runtime already serious enough that docs must stay precise?

Start with:

1. [Current Compose Runtime](../architecture/current-compose-runtime.md)
2. root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
3. the included `compose/docker-compose.*.yml` fragments

What this pack proves well:

- the stack is already broad and platform-shaped
- the root file is still a major authoring surface
- ingress, mesh, observability, and mixed protocol domains are all live now

What it must not be widened into:

- the distributed truth problem is therefore already solved

## Minimum retrieval workflow before rewriting a page

Before replacing or heavily revising an infra page, use this sequence:

1. identify the page's primary claim type
2. route it through the claim router above
3. read at least one source from the strongest class
4. read at least one source that states what remains unproven
5. if the page reconstructs user frustration, read an archive-pressure source
6. only then write the page

If steps `3` through `5` are skipped, the page is very likely to flatten the
repo into a smoother but weaker story.

## Anti-cheat rules for assimilation

Never let:

- an archive complaint become runtime proof
- a planning document become implementation proof
- a proxy feature become full failover proof
- a DNS plurality story become request-preservation proof
- an option catalog become evidence that the operator's burden is lower

Never summarize a source class without also naming its boundary.

Never use "the repo wants X" and "the runtime proves X" interchangeably.

If contradiction disappears entirely, you probably retrieved too shallowly.

## Bottom line

Good retrieval in `bolabaden-infra` does not mean "more files."
It means routing each claim to the right authority, preserving the repo's own
honesty boundaries, and keeping the user's real complaint alive instead of
replacing it with a calmer adjacent problem.

If the docs do that, they become useful.
If they do not, they become infrastructure theater with citations.
