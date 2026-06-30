# Source Assimilation Index

This page is the retrieval contract for the knowledgebase.

Its job is not just to say "consult sources."
Its job is to stop the docs from sounding broad, careful, and well-organized
while still quietly answering a smaller question than the user actually asked.

In this repository, bad retrieval often looks polished.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what "actual assimilation" means in this repository
- how source classes should be separated before stronger claims are written
- why retrieval discipline is part of the repo's honesty contract
- which source families matter most for the no-Swarm, wrong-node, burden
  transfer problem

This page is not allowed to prove:

- that a page has already assimilated the repo just because it cites many files
- that source abundance itself produces better answers
- that retrieval discipline can substitute for runtime proof

This is a retrieval contract page, not a claim that the work has already been
done correctly.

## What "actually RAG this time" means here

In `bolabaden-infra`, "actually RAG this time" does **not** mean:

- read more files
- mention more tools
- summarize more adjacent technologies
- produce a calmer architecture story

It means:

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

## The source families that matter most

This repo has four major evidence classes:

1. architecture-intent and honesty surfaces
2. live root runtime surfaces
3. planning and promotion surfaces
4. archive-pressure surfaces

The biggest recurring retrieval mistake is flattening those four classes into
one blended narrative.

Once that happens, at least one of these becomes false:

- the worktree starts sounding more complete than it is
- the user's actual complaint gets replaced by a neighboring, weaker question

## Priority source map

Use this map when reconstructing the repo's actual architecture problem.

### 1. Dream and honesty surfaces

Read first when the page is trying to recover what the repo wants to make
true.

Primary anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
- [`knowledgebase/AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/AGENTS.md)
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

Use these for:

- the target operating contract
- honesty boundaries
- repo-authoring priorities
- how the repo wants Compose to remain central

Do **not** use these alone for:

- present-tense runtime capability claims
- route-specific failover claims
- proof that a missing middle layer already exists

### 2. Live runtime surfaces

Read first when the page is claiming what the priority implementation actually
ships today.

Primary anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- merged output from `docker compose config`

Use these for:

- service presence
- network presence
- config and secret surfaces
- active fragment membership
- proof that a component is in the runtime at all

Do **not** use these alone for:

- generic wrong-node success
- backend-loss fallback durability
- shared placement truth
- stateful correctness

### 3. Planning and promotion surfaces

Read first when the page is asking what the repo has already named as missing
or what a stronger layer would need to earn.

Primary anchors:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- related plan docs under `/docs`
- roadmap, proof, and promotion pages in the knowledgebase

Use these for:

- named missing truths
- candidate repair families
- promotion thresholds
- why some stronger layers are still unearned

Do **not** use these alone for:

- claiming the repair is live
- claiming a candidate already won
- implying the runtime has already crossed the threshold

### 4. Archive-pressure surfaces

Read first when the page is trying to reconstruct what the user is really
rebelling against and why ordinary answers keep failing.

Primary anchors:

- [`knowledgebase/source-archive/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/)
- synthesis pages under [`research/`](../research/evidence-ledger.md)

High-value recurring threads include:

- `docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`
- `distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md`
- `load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`
- `nomad-multi-node-failover__68765e45-1ec4-8006-9179-5ef176d7a90f.md`

Use these for:

- recurring burden-transfer complaints
- why certain candidate families keep recurring
- why generic option lists feel smaller than they sound

Do **not** use these alone for:

- present-tense implementation proof
- declaring one explored option live
- route-specific success claims

## Retrieval packets by page type

Different page types need different minimum source packets.

### Doorway or overview pages

Minimum packet:

- one dream source
- one live runtime source
- one planning source
- one archive-pressure source

Required outcome:

- the reader can tell what the repo wants, what it ships, what it still lacks,
  and why the user is still dissatisfied

### Runtime pages

Minimum packet:

- root Compose file
- active fragments
- if needed, merged config output

Required outcome:

- the page says what is present now and explicitly refuses to over-upgrade that
  presence into distributed capability proof

### Architecture or roadmap pages

Minimum packet:

- one dream source
- one live runtime source
- one planning source

Required outcome:

- the page can state the gap between current runtime and target behavior
  without narrating the gap as if it were mostly closed

### Research or evidence pages

Minimum packet:

- one live runtime source
- one archive-pressure source
- one plan or intent source

Required outcome:

- the page preserves the accusation, keeps source classes separate, and defines
  what is still illegal to claim

## What still does not count as real retrieval here

The following still do not count as "actually RAG this time" in this repo:

- reading many files without distinguishing source class
- blending archive pressure into runtime proof
- treating planning language as if it outranks the worktree
- producing a calm answer that edits the wound down into a smaller question
- summarizing related technologies without reconstructing the user's benchmark
- removing contradiction because it makes the docs cleaner
- recovering the ecosystem around the wound while leaving the wound itself
  under-described

That last failure mode is common here.
The docs can become more exhaustive and still answer the wrong question.

## The small auditable packet every retrieval pass should leave behind

Actual assimilation here should leave behind a small packet, not just a bigger
stack of citations.

At minimum the packet should preserve:

- the exact user-facing accusation being reconstructed
- the strongest runtime artifact consulted
- the strongest dream or archive artifact consulted
- the private burden still left over after both were read
- the stronger sentence that still stayed illegal

If the retrieval pass cannot produce that packet, then "we really read the repo
this time" is still too self-congratulatory for this project.

## Bottom line

This repo is not mainly asking for broader summaries.
It is asking for evidence custody.

Actual assimilation here means:

- we preserved the accusation
- we did not let runtime, intent, plan, and archive trade confidence illegally
- we kept naming what still stayed operator-owned

If a page becomes more cross-linked, more source-aware, and more ecosystem-rich
while becoming less able to name the surviving operator-owned truth, the page
got worse.
