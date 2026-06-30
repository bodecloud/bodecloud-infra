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
- `compose/docker-compose.core.yml`
- `compose/docker-compose.docs.yml`
- `compose/docker-compose.firecrawl.yml`
- `compose/docker-compose.headscale.yml`
- `compose/docker-compose.metrics.yml`
- `compose/docker-compose.llm.yml`
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
4. relevant research pages under
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
- live runtime
- missing gap
- archive pressure

Only then choose sources.

### Rule 3: let the worktree outrank elegant synthesis for runtime claims

If a planning page or archive conversation sounds better than the live compose
surface, the live compose surface still wins for "what exists now."

### Rule 4: never let planning docs impersonate implementation

Planning material is often clearer than the worktree.
That does not make it runtime proof.

### Rule 5: never let archive pressure impersonate architecture proof

The archive helps recover the real complaint.
It does not tell you which implementation path is complete.

### Rule 6: preserve contradiction

If the repo simultaneously says:

- "this is the dream"
- and "this is still missing"

do not smooth that into one moderate sentence.
Expose both and state the boundary.

### Rule 7: always ask where operator memory is still doing system work

This is one of the highest-value retrieval questions in the whole repo.

Whenever a page describes a topology, also ask:

- what truth still lives in private operator reconstruction?
- what should eventually be externalized into shared current-state knowledge?

If the docs cannot answer that, they have not really assimilated the system.

## What the source archive is best used for

The source archive is most useful when a page needs to recover one of these:

- the user's exact complaint
- the direct statement of the desired wrong-node behavior
- the reason familiar orchestration answers were not accepted
- the emotional pressure behind "fake options" and "fake HA" language

Some of the most important archive threads for this repo's identity are:

- `docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`
- `load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
- `distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
- `traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md`
- `docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md`

These files help explain why the repo exists in its current shape.
They do not remove the duty to inspect the current worktree.

## Failure modes this page is explicitly trying to prevent

### Failure mode 1: flattening everything into one "balanced" voice

This makes the docs sound mature while erasing source authority.

### Failure mode 2: answering the easier neighboring question

This is the most common way to be technically adjacent and still useless.

### Failure mode 3: letting availability-sounding words do proof work

Examples:

- "multi-node"
- "distributed"
- "fallback"
- "load balanced"
- "HA"

These words mean almost nothing unless routed to the right evidence class.

### Failure mode 4: confusing "many options exist" with "the user has a real path"

This is the emotional center of the repo.
The docs should be actively hostile to this confusion.

## Bottom line

This knowledgebase only becomes useful if it behaves like a retrieval system
with source discipline, not like a summarizer with a calm tone.

The correct assimilation standard is:

> reconstruct the user's real dream, route each claim to the strongest source
> class for that claim, preserve contradiction, and never let a nicer source
> silently overrule current worktree truth.

Anything weaker will keep producing the same kind of documentation the user is
already frustrated with:

- tidy
- plausible
- sourced
- and still not answering the real question
