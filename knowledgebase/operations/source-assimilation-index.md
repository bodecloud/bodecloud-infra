# Source Assimilation Index

This page is the retrieval contract for the knowledgebase.

Its job is not merely to say "use sources."
Its job is to stop the docs from sounding broad, careful, and cross-linked
while still quietly answering a smaller question than the user actually asked.

In this repo, bad retrieval often looks polished.

## The real problem this page is trying to stop

The common failure here is not lack of reading.
It is reading many true things and still answering the wrong question.

That usually happens like this:

1. the archive reconstructs the wound
2. the instruction files reconstruct the dream
3. the runtime proves serious machinery exists
4. the plans describe coherent next steps
5. the final page becomes broader while the hidden burden stays vague

This page exists to stop step 5.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what "actual assimilation" means in `bolabaden-infra`
- how source classes must stay separate before stronger claims are written
- why retrieval discipline is part of the honesty contract
- which source families matter most for the no-Swarm, wrong-node,
  burden-transfer problem

This page is not allowed to prove:

- that a page has already assimilated the repo just because it cites many
  files
- that source abundance by itself produces better answers
- that retrieval discipline can substitute for runtime proof

This is a retrieval contract page, not a completion badge.

## What "actually RAG this time" means here

In this repository, "actually RAG this time" does not mean:

- read more files
- mention more tools
- summarize more adjacent technologies
- produce a calmer architecture story
- flatten dream, runtime, plan, and archive into one clean narrative

It means:

1. recover the real dream before drafting
2. identify what class of claim the page is making
3. route that claim to the strongest source class
4. say exactly what that source class proves
5. say exactly what it does not prove
6. preserve contradiction instead of ironing it out
7. keep worktree truth above elegant prose for runtime claims
8. use archive pressure to recover the user's accusation without letting the
   archive impersonate implementation proof

If the docs skip that sequence, they may become longer without becoming more
truthful.

## The accusation retrieval must keep alive

Retrieval is only aligned here if it preserves the user's accusation in a form
the next page can still feel:

> there seem to be endless options for multi-node Docker, failover, clustering,
> proxies, overlays, discovery, and orchestration, but too many of them solve
> one visible layer and then quietly leave the operator as the hidden control
> plane when reality gets sharp.

If retrieval produces a calmer answer that no longer feels accused by that
sentence, it probably answered a neighboring question instead.

## The most useful retrieval question

The most useful assimilation question in this repo is often:

> after reading these sources, what exact sentence is still privately finished
> by the operator?

Examples:

- `I still personally know what runs where.`
- `I still personally know which peer is truly eligible.`
- `I still personally know whether the fallback survives backend loss.`
- `I still personally know whether the protected route still means the same
  thing after handoff.`
- `I still personally know who the real writer is.`

If a retrieval pass cannot name that surviving sentence, it usually has not
actually reconstructed the repo's real problem.

## The four evidence classes that matter most

This repo has four major evidence classes:

1. architecture-intent and honesty surfaces
2. live root runtime surfaces
3. planning and promotion surfaces
4. archive-pressure surfaces

The biggest recurring mistake is flattening those four classes into one blended
voice.

Once that happens, at least one of these becomes false:

- the worktree starts sounding more complete than it is
- the user's complaint gets replaced by a neighboring weaker question
- plans begin lending runtime confidence they did not earn
- archive synthesis begins sounding like implementation proof

## Why source abundance is not the same thing as assimilation

This repo is already large enough that an answer can cite:

- instruction files
- runtime files
- plan files
- archive syntheses
- research pages

and still be weak.

Assimilation only starts when the answer also states:

- which source class is carrying which part of the claim
- what the strongest consulted artifact still does not prove
- what hidden burden remained after the reading pass

Without that, the answer may be better sourced and still badly assimilated.

## Priority source map

Use this map when reconstructing the repo's actual architecture problem.

### 1. Dream and honesty surfaces

Read these first when the page is trying to recover what the repo wants to make
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
- authoring priorities
- why Compose remains central
- why heavier control layers are not allowed to win by tone alone

Do not use these alone for:

- present-tense runtime capability claims
- route-specific failover claims
- proof that a missing middle layer already exists

Likely surviving private sentence after dream-only reading:

> yes, but I still personally do not know what the runtime truly owns today.

### 2. Live runtime surfaces

Read these first when the page is claiming what the priority implementation
actually ships today.

Primary anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- merged output from `docker compose config`

Use these for:

- service presence
- network presence
- fragment inclusion
- label, config, and secret surfaces
- proof that a component is in the runtime at all

Do not use these alone for:

- generic wrong-node success
- backend-loss fallback durability
- shared placement truth
- stateful correctness

Likely surviving private sentence after runtime-only reading:

> yes, but I still personally do not know whether these components cash out
> into the burden transfer the user actually wants.

### 3. Planning and promotion surfaces

Read these first when the page is asking what the repo has already named as
missing or what a stronger layer would need to earn.

Primary anchors:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- related plan docs under `/docs`
- roadmap, proof, and architecture pages under `knowledgebase/`

Use these for:

- named missing truths
- candidate repair families
- promotion thresholds
- why some stronger layers are still unearned

Do not use these alone for:

- claiming the repair is live
- claiming a candidate already won
- implying the runtime already crossed the threshold

Likely surviving private sentence after plan-only reading:

> yes, but I still personally do not know whether the repair exists outside
> the plan.

### 4. Archive-pressure surfaces

Read these first when the page is trying to reconstruct what the user is
actually rebelling against and why normal answers keep failing.

Primary anchors:

- [`knowledgebase/source-archive/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/)
- synthesis pages under `knowledgebase/research/`

High-value recurring thread families include:

- multi-node Docker without Swarm
- distributed HA orchestration
- load-balancer and failover alternatives
- Compose fork, fallback-list, and runtime watcher pressure
- shared-IP, VPN, Tailscale, WireGuard, and Cloudflare public-entry pressure
- manual-placement acceptance versus current placement truth
- Nomad, k3s, Kubernetes, and control-plane comparisons
- reverse-proxy and middleware continuity discussions
- helper-layer frustration where the helper disappears under the failure it was
  meant to absorb

Use these for:

- recurring burden-transfer complaints
- why certain candidate families keep reappearing
- why generic option lists feel smaller than they sound

Do not use these alone for:

- present-tense implementation proof
- declaring one explored option live
- route-specific success claims

Likely surviving private sentence after archive-only reading:

> yes, but I still personally do not know what the current worktree truly
> proves.

## Retrieval packets by page type

Different page types need different minimum source packets.

### Doorway or overview pages

Minimum packet:

- one dream source
- one live runtime source
- one planning source
- one archive-pressure source

Required outcome:

- the reader can tell what the repo wants
- what it ships
- what it still lacks
- why the user is still dissatisfied

Required surviving warning:

- `do not let the doorway page sound more mature than the runtime`

### Runtime pages

Minimum packet:

- root Compose file
- active fragments
- if needed, merged config output

Required outcome:

- the page says what is present now
- it explicitly refuses to upgrade component presence into distributed-capability
  proof

Required surviving warning:

- `do not confuse present components with present burden transfer`

### Architecture or roadmap pages

Minimum packet:

- one dream source
- one live runtime source
- one planning source

Required outcome:

- the page can state the gap between current runtime and target behavior
  without narrating the gap as if it were mostly closed

Required surviving warning:

- `do not let sequencing sound like implementation`

### Research or evidence pages

Minimum packet:

- one live runtime source
- one archive-pressure source
- one plan or intent source

Required outcome:

- the page preserves the accusation
- keeps source classes separate
- defines what is still illegal to claim

Required surviving warning:

- `do not let archive frustration impersonate runtime proof`

### Operator pages

Minimum packet:

- one live runtime source
- one dream source
- one archive-pressure source

Required outcome:

- the page must say what truth is still operator-owned today
- and what exact artifact would externalize it

Required surviving warning:

- `do not let operational clarity impersonate reduced burden`

Operator pages should also name the packet field they are trying to externalize.

Examples:

| If the operator page asks... | It should point at... | Stronger claim still forbidden until... |
| --- | --- | --- |
| "Can the wrong node handle this?" | `entry_node`, `locality_result`, `placement_source`, `placement_decision_packet` | a drill shows the wrong node used shared truth. |
| "Can it forward to the right peer?" | `selected_peer`, `peer_eligibility_reason`, `placement_decision_packet.peer_eligibility` | eligibility is demonstrated, not assumed from reachability. |
| "Did auth survive?" | `policy_chain`, `handoff.preserves_auth` | local and wrong-node behavior are compared. |
| "Did failover work?" | `backend_condition`, `backend_loss` | the preferred backend is actually removed or broken during the test. |
| "Is this stateful service HA?" | authority, writer, recovery, fencing fields in a stateful packet | authority transfer or honest singularity is proven. |

This keeps operator pages from turning into better explanations of the same
private burden.

## The retrieval sequence that should happen before writing

Before drafting, do this in order:

1. name the user-facing wound being reconstructed
2. name the smaller neighboring question the page might accidentally answer
3. name the claim class the page wants to make
4. pull the strongest dream surface
5. pull the strongest runtime surface
6. pull the strongest plan or archive surface needed to explain the gap
7. write down what truth is still privately owned after reading all of them
8. only then draft the page

If the writing begins before step 7, the page usually becomes too smooth.

## Archive retrieval drill

For pages that use the imported archive, "read the archive" is too vague.
Do this instead:

1. Search the source archive by the failure class, not only by the tool name.
   Use terms like `wrong-node`, `fallback`, `service discovery`, `Cloudflare`,
   `shared IP`, `VPN`, `Compose fork`, `Nomad`, `OpenSVC`, `k3s`, `stateful`,
   `middleware`, and `auth`.
2. Pick at least one thread that supports the user's accusation and one thread
   that complicates the easy answer.
3. Classify each selected thread by stack layer:
   - public entry
   - service placement
   - peer eligibility
   - route durability
   - policy continuity
   - stateful authority
4. Extract the exact user question or correction that made the ordinary answer
   inadequate.
5. Write the forbidden upgrade next to it.

The drill should produce a small table like this before prose starts:

| Thread family | Stack layer | User pressure | Ordinary dodge | Forbidden upgrade |
| --- | --- | --- | --- | --- |
| Docker multi-node without Swarm | service placement | manual placement is fine; request-time discovery is not | sell a scheduler because placement is manual | scheduling maturity proves wrong-node routing |
| Dynamic HA proxy/shared IP | public entry | first-hop SPOF and service-level SPOF are different | claim single-IP or DNS plurality solves failover | IP-level survival proves service-level correctness |
| Compose fork/fallback lists | route durability | fallback must happen at runtime, not only in YAML | propose preprocessing or prettier schema | generated config proves backend-loss behavior |

If a page cannot fill this table, it may still cite the archive, but it has not
assimilated it.

This matters because the archive is full of near-misses.
Near-misses are useful only when they are preserved as near-misses.
If they are blended into one generic HA complaint, the docs lose the user's
real standard again.

## Why "too smooth" is a real warning sign here

Smoothness is suspicious in this repo because the underlying materials are not
smooth:

- the dream is sharper than the runtime
- the runtime is richer than the proof
- the options are broader than the honest choices
- the archive is angrier than normal infra prose expects

If the draft becomes tidy by flattening those tensions, retrieval probably
failed even if it was thorough.

## What still does not count as real retrieval here

These still do not count as "actually RAG this time":

- reading many files without distinguishing source class
- blending archive pressure into runtime proof
- treating planning language as if it outranks the worktree
- producing a calm answer that edits the wound down into a smaller question
- summarizing related technologies without reconstructing the benchmark
- removing contradiction because it makes the docs cleaner
- recovering the ecosystem around the wound while leaving the wound itself
  under-described

That last failure mode is common here.
The docs can become more exhaustive and still answer the wrong question.

## The small auditable packet every retrieval pass should leave behind

Actual assimilation here should leave a small packet, not just a bigger pile
of citations.

At minimum the packet should preserve:

- the exact accusation being reconstructed
- the strongest runtime artifact consulted
- the strongest dream, plan, or archive artifact consulted
- the private burden still left over after both were read
- the stronger sentence that still stayed illegal

If a retrieval pass cannot produce that packet, then "we really read the repo
this time" is still too congratulatory for this project.

The packet should also include an archive-family field when archive pressure is
used:

```yaml
archive_family:
  thread: "docker-multi-node-without-swarm"
  stack_layer: "service placement"
  user_pressure: "manual placement is acceptable; request-time placement truth is not"
  ordinary_dodge: "recommend a scheduler because placement is manual"
  forbidden_upgrade: "scheduler vocabulary proves wrong-node service discovery"
```

This field keeps the archive from becoming decorative citation.
It forces the writer to preserve why that thread matters and which stronger
sentence it still cannot support.

For pages about ingress, routing, or failover, the packet should be concrete
enough to downgrade into the route-level and placement-decision schemas:

```yaml
assimilation_packet:
  accusation: "the operator is still the hidden control plane"
  runtime_anchor: "docker-compose.yml + active compose fragments"
  intent_anchor: ".github/copilot-instructions.md"
  archive_anchor: "source-archive threads that preserve the frustration"
  archive_family:
    thread: "docker-multi-node-without-swarm"
    stack_layer: "service placement"
    user_pressure: "manual placement is accepted; current placement truth is missing"
    ordinary_dodge: "sell orchestration because placement is manual"
    forbidden_upgrade: "orchestrator comparison proves runtime placement truth"
  legal_sentence: "the stack has serious ingredients and a precise dream"
  illegal_sentence: "the stack already proves generic wrong-node failover"
  surviving_private_sentence: "I still personally know where this route lives"
  next_runtime_packet:
    route_packet_field: placement_source
    placement_decision_packet_field: placement_source
    required_upgrade: "receiving node consults shared current placement truth"
```

The exact values will change by page.
The point is that retrieval should end with a next proof field, not a vague
recommendation to "test more."

If the next proof field is `placement_source`, the writer should usually link
to [Current-State Registry and Peer Eligibility](../architecture/current-state-registry-and-peer-eligibility.md)
and name which part of `placement_decision_packet` is still missing.
That prevents "we need service discovery" from staying too abstract.

## The honest bottom line

This repo is not mainly asking for broader summaries.
It is asking for evidence custody.

The real retrieval standard is not:

> did we read enough?

It is:

> after reading, what accusation stayed alive, what source class carried the
> answer, and what burden still remained privately held?

Actual assimilation here means:

- the accusation stayed alive
- runtime, intent, plan, and archive did not trade confidence illegally
- the docs kept naming what still remained operator-owned

If a page becomes more cross-linked, more ecosystem-rich, and more
source-aware while becoming less able to name the surviving operator-owned
truth, the page got worse.
