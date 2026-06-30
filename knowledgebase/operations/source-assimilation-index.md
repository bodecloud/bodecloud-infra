# Source Assimilation Index

This page is the retrieval contract behind the knowledgebase.

Its job is not merely to say "cite sources."
Its job is to stop the docs from sounding comprehensive while still quietly
answering a smaller, safer, easier question than the user actually asked.

In this repo, bad retrieval usually looks polished.

It sounds like:

- many files were consulted
- many tools were mentioned
- many options were listed
- many architecture words were used

and the user's real problem still got downgraded.

This page exists to stop that downgrade.

## The dream this page has to protect

The user's frustration is not mainly about missing documentation.

It is about a broader pattern:

- the ecosystem offers many partial tools
- each tool solves only a slice
- the operator is left privately reconstructing the real system
- docs then pretend the existence of options is the same thing as a real path

So "actually RAG this time" in `bolabaden-infra` does **not** mean:

- gather more files
- blend them into a calm voice
- smooth over contradictions
- produce a plausible average summary

It means:

1. reconstruct the real dream
2. identify what class of claim is being made
3. route that claim to the strongest authority for that class
4. say exactly what that source proves
5. say exactly what it does not prove
6. preserve contradiction instead of ironing it out
7. keep current worktree truth above elegant prose for runtime claims
8. use archive pressure to recover the user's actual complaint without letting
   the archive impersonate runtime evidence

If the docs skip that sequence, they may become longer but not more useful.

## Strongest honest current answer

The repo has four major evidence classes, and they do not carry equal weight:

1. architecture-intent and honesty surfaces
2. live root runtime surfaces
3. planning and promotion surfaces
4. archive-pressure surfaces

The biggest recurring retrieval mistake is flattening those four classes into
one blended narrative.

Once that happens, at least one of two things becomes false:

- the worktree starts sounding more complete than it is
- the user's actual complaint gets replaced by a neighboring, weaker question

This page exists to keep both failures visible.

## Source classes and what they are for

### 1. Architecture-intent and honesty surfaces

These are the strongest sources for:

- what the repo wants the runtime to do
- what the dream actually is
- where the project explicitly refuses fake HA language

Read in this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
4. [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

These files do **not** have equal weight.

Their real hierarchy is:

- `copilot-instructions.md`
  - strongest statement of the no-Swarm/no-Kubernetes, local-first,
    peer-forward, anti-SPOF dream
- `README.md`
  - strongest repo-level honesty wall about what the stack still does not
    prove
- `AGENTS.md`
  - operator/setup/runtime constraints, not the deepest architecture dream
- `.cursorrules`
  - authoring discipline and Compose hygiene, not core architecture proof

If a page treats those four files as if they all define the same thing, the
page has already started flattening the repo.

### 2. Live root runtime surfaces

These are the strongest sources for:

- what the priority implementation actually is
- what the root stack really owns today
- what helper layers are live now

Start with:

1. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. active fragments under
   [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

Important current fragments include:

- `compose/docker-compose.coolify-proxy.yml`
- `compose/docker-compose.docs.yml`
- `compose/docker-compose.firecrawl.yml`
- `compose/docker-compose.headscale.yml`
- `compose/docker-compose.llm.yml`
- `compose/docker-compose.metrics.yml`
- `compose/docker-compose.stremio-group.yml`
- `compose/docker-compose.warp-nat-routing.yml`
- `compose/docker-compose.wishlist.yml`

These files outrank planning docs whenever the claim is:

- this exists now
- this is active in the root stack
- this network, service, or helper is part of the priority runtime

They do **not** outrank planning docs for:

- what the repo wishes those services did under failure
- whether a helper layer is already trustworthy

### 3. Planning and promotion surfaces

These are the strongest sources for:

- what gaps the repo has already named explicitly
- what future promotion candidates exist
- what bug or capability boundary the project already understands

Start with:

1. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
2. [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
3. [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
4. research pages under
   [`knowledgebase/research/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research)

These files are where the repo already says things like:

- route persistence under backend loss is still missing
- automated service failover between nodes is still missing
- `docker-gen-failover` currently deletes routes when containers stop
- Cloudflare DDNS presence is not the same thing as full multi-node failover
- Headscale is single-node today
- stateful HA needs replication and topology-specific correctness, not just
  container movement

These planning sources are allowed to explain what is missing.
They are not allowed to silently promote themselves into live proof.

### 4. Archive-pressure surfaces

These are the strongest sources for:

- what the user is actually rebelling against
- why ordinary "just use X" answers keep failing
- where the missing middle layer keeps reappearing

Start with:

- [User Intent and Dream](../research/user-intent-and-dream.md)
- [Archive Pressure Patterns](../research/archive-pressure-patterns.md)
- [Evidence Ledger](../research/evidence-ledger.md)
- the source archive conversations under
  [`knowledgebase/source-archive/chatgpt-exports/conversations/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/chatgpt-exports/conversations)

This class is essential because the repo is not only technical.
It is also a response to a recurring ecosystem failure:

> the tooling landscape keeps offering fragments and naming them as if the
> operator's real burden has been removed.

Archive pressure restores that complaint.
It does **not** prove runtime correctness.

## Claim router

Use this table before writing any summary page.

| If the sentence is really claiming... | Start with | Why | It still must not imply... |
| --- | --- | --- | --- |
| "this is the dream of the repo" | `copilot-instructions.md`, `README.md` | strongest architecture-intent and honesty sources | that the runtime already earned it |
| "this exists in the live stack today" | root Compose runtime | strongest implementation truth | that it survives the relevant failure |
| "the repo already knows this is missing or broken" | planning docs | strongest named-gap sources | that the repair is already live |
| "the user is rejecting this whole category of answer" | archive-pressure sources | strongest statement of the actual complaint | that complaint itself proves the implementation |

If a sentence cannot be routed through this table, it is probably too vague to
be useful.

## Retrieval packs for the most common real questions

These are the highest-signal local bundles for the questions this repo keeps
returning to. They exist because broad archive scans produce volume faster than
understanding.

### Pack A: "Why does Compose stop feeling honest the moment multi-node matters?"

Use when the real question is:

- why is the user so hostile to static glue?
- why is `docker-compose.yml` still sacred even though it is now painful?
- why does the repo keep insisting the missing problem is not "better YAML"?

Start with:

1. [User Intent and Dream](../research/user-intent-and-dream.md)
2. [Problem and Goals](../architecture/problem-and-goals.md)
3. `source-archive/chatgpt-exports/conversations/docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`
4. `source-archive/chatgpt-exports/conversations/docker-compose-multi-server-setup__67f73c50-150c-8006-8408-c03db2d8d287.md`

What this pack proves well:

- the user does not merely dislike Compose syntax
- the real break happens when runtime truth escapes local node boundaries
- advice that says "just run Traefik and DDNS" keeps stopping one layer too
  early

What it must not be widened into:

- Compose itself is useless
- remote Docker sockets or overlay mesh automatically solve the hidden burden

### Pack B: "What counts as a fake failover answer in this repo?"

Use when the real question is:

- why are Cloudflare, Traefik, and more healthchecks still not enough?
- why does the repo keep saying wrong-node dignity is the real threshold?
- why is a helper like `docker-gen-failover` both important and distrusted?

Start with:

1. [HA, Failover, and Routing](../architecture/ha-failover-routing.md)
2. [Request Path and Failure Walkthrough](../architecture/request-path-and-failure-walkthrough.md)
3. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
4. `source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
5. `source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md`

What this pack proves well:

- the user is not looking for a generic reverse proxy
- Traefik configuration alone does not grant failover semantics
- the repo already knows route persistence under backend loss is one of the
  hardest seams

What it must not be widened into:

- the current runtime has no meaningful ingress machinery
- all proxy-based fallback ideas are permanently invalid

### Pack C: "Does some orchestrator already earn default promotion?"

Use when the real question is:

- should the repo just move to Nomad, k3s, OpenSVC, or something else?
- is the user only resisting mature tools out of taste?
- what exact burden would a bigger control plane have to own better?

Start with:

1. [Orchestration Options](../architecture/orchestration-options.md)
2. [The Missing Middle Layer](../architecture/missing-middle-layer.md)
3. [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
4. `source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
5. `source-archive/chatgpt-exports/conversations/nomad-multi-node-failover__68765e45-1ec4-8006-9179-5ef176d7a90f.md`
6. [Garden and k3s Exploration Evidence](../research/garden-k3s-exploration-evidence.md)
7. [Nomad Exploration Evidence](../research/nomad-exploration-evidence.md)

What this pack proves well:

- the repo is not anti-orchestrator in principle
- "respectable platform family" is weaker than "owns a named hidden burden"
- the user is explicitly searching for a middle layer, not merely for a brand

What it must not be widened into:

- every stronger control plane is equally bad
- peer-equal or leaderless rhetoric is automatically the right answer

### Pack D: "What does the current worktree actually own today?"

Use when the real question is:

- what is active in the priority implementation?
- which pieces are live versus planned versus archived?
- where should route or topology claims be checked first?

Start with:

1. [Instruction Surfaces and Authority](../architecture/instruction-surfaces-and-authority.md)
2. [Current Compose Runtime](../architecture/current-compose-runtime.md)
3. [Compose Fragment Map](../architecture/compose-fragment-map.md)
4. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
5. active `compose/docker-compose.*.yml` fragments

What this pack proves well:

- the root runtime is serious and broad
- specific services, networks, and helper components are really present
- the current stack already encodes real routing, mesh, and observability
  pressure

What it must not be widened into:

- broad wrong-node resilience is already live
- stateful HA can be inferred from component presence

## What "full assimilation" means in this repo

Full assimilation is not achieved when the docs can repeat the repository.

It is achieved when the docs can explain, with source discipline:

- what the user actually wants
- what the current worktree really does
- where the worktree still falls short
- which missing layers the repo keeps circling
- which nearby explanations would sound reasonable but still miss the point

That last item is critical.

Many weak summaries in this repo have been wrong not because they invented
facts, but because they answered a neighboring question instead of the real
one.

Examples of neighboring but weaker questions:

- "How do I host many services?"
- "What reverse proxies exist?"
- "What orchestrators are popular?"
- "Can Docker Compose run on more than one host?"

Those are all easier than the real question:

> how do we keep the directness of Compose while making multiple ordinary
> Docker nodes behave less stupidly under wrong-node entry, backend loss, and
> anti-SPOF pressure?

If the docs drift back to the easier questions, assimilation has failed.

## Retrieval rules for documentation work

When rewriting or summarizing, follow this sequence.

### Rule 1: start with the user question, not the folder name

Do not open files merely because they are nearby.
Open files because they can answer a specific claim the page needs to make.

### Rule 2: identify claim type before gathering evidence

Decide whether the page is trying to say:

- dream
- runtime
- planning gap
- archive complaint
- proof boundary

If you do not identify the claim type first, you will almost certainly start
letting unlike evidence classes borrow authority from each other.

### Rule 3: prefer the smallest strong bundle over the broadest sweep

In this repo, six sharply chosen files usually beat sixty vaguely adjacent
ones.

Breadth matters only after the strongest sources for the exact claim have
already been read.

### Rule 4: preserve contradiction on purpose

If runtime, planning, and archive pressure disagree in tone, keep that visible.
The disagreement is often more informative than a blended summary.

### Rule 5: always name the ceiling

Every serious page should make it obvious:

- what the sources let it say confidently
- what stronger sentence would still be a lie

If that ceiling is missing, the page is probably overclaiming.

## The easiest ways to fail at RAG in this repo

These failure modes are now common enough to name explicitly.

### Failure mode 1: option laundering

The docs list enough products and frameworks that the reader feels the option
space has become healthy, even though no option has yet removed the real hidden
burden.

### Failure mode 2: intent inflation

The docs read a strong dream file and start narrating current behavior as if
the runtime had already earned it.

### Failure mode 3: archive inflation

The docs read a large number of pain threads and start narrating the repo as if
the repeated pain itself were evidence of implementation progress.

### Failure mode 4: runtime flattery

The docs see a serious stack and start narrating its presence as if presence
alone made wrong-node, fallback, or stateful semantics trustworthy.

### Failure mode 5: plan sedation

The docs find detailed plans and mistake the clarity of the repair path for
partial completion of the repair itself.

## Bottom line

The retrieval standard in `bolabaden-infra` is not "did we read a lot?"

It is:

> did we reconstruct the user's real demand, route each claim to the right
> authority class, preserve the repo's own contradictions, and refuse to let
> serious-looking adjacent evidence merge into broader confidence than the
> current worktree has actually earned?

If the answer is no, then the docs may be longer than before and still not
actually more honest.
