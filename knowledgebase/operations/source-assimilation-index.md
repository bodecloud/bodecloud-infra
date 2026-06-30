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

## Source-to-packet assimilation

Archive material should not enter the knowledgebase as atmospheric context.
It should enter as pressure on a specific claim.

The useful question is not:

> what did this source talk about?

The useful question is:

> which packet field, forbidden upgrade, or surviving private sentence did this
> source make sharper?

Use this shape when assimilating any archive conversation, exported chat,
research note, or plan fragment:

```yaml
source_assimilation_packet:
  source: "knowledgebase/source-archive/<provider>/conversations/<file>.md"
  source_family: "wrong-node-routing | fallback | protected-http | raw-tcp | stateful-authority | orchestrator-promotion | operator-frustration"
  claim_class: "intent | runtime | plan | archive-pressure | proof"
  claim_pressure: "<the narrow claim this source pressures>"
  packet_field_pressured:
    - "entry_node | locality_result | placement_source | selected_peer | peer_eligibility | policy_chain | backend_loss | stateful_authority | forbidden_claim"
  strongest_sentence_supported: "<one sentence this source can support>"
  still_forbidden: "<the stronger sentence this source cannot support>"
  surviving_private_sentence: "<what the operator still has to know alone>"
  page_destination: "<page that absorbed this pressure>"
```

This packet does not have to be stored literally for every source.
But the page that uses the source should leave those answers visible enough
that a later contributor can reconstruct the custody chain.

If a source cannot be mapped to a packet field, it may still be useful
background.
It is not yet a burden-transfer source.

### Examples of correct assimilation

Correct:

- a Redis load-balancing archive pressures `raw-tcp` and
  `stateful-authority` fields
- it supports "Redis is TCP and can be routed with TCP routers"
- it forbids "Traefik TCP labels across non-Swarm Docker hosts make Redis HA"
- it leaves "I still personally know who the writer is" as the surviving
  private sentence

Correct:

- a multi-node Docker without Swarm archive pressures `placement_source`,
  `selected_peer`, and `peer_eligibility`
- it supports "the repo wants ordinary Docker nodes to share enough current
  truth for wrong-node routing"
- it forbids "manual multi-host Compose has already become a distributed
  control plane"
- it leaves "I still personally know what runs where" as the surviving private
  sentence

Incorrect:

- citing several failover conversations and then writing "the repo has a broad
  HA strategy"

That is not assimilation.
That is source laundering.

## Assimilation ledger rule

When a page uses archive pressure to justify a stronger reading, it should also
name the ledger transition:

| Source pressure | Packet field moved | Stronger sentence still illegal |
| --- | --- | --- |
| Archive says non-Swarm hosts need explicit shared truth | `placement_source` | wrong-node routing is solved |
| Archive says a fallback helper can disappear under backend loss | `backend_loss` | generated config proves failover |
| Archive says auth/middleware behavior can change across proxy handoff | `policy_chain` | protected route answered, so policy survived |
| Archive says Redis or MongoDB are TCP/stateful | `stateful_authority` | TCP reachability proves stateful HA |

This is the bridge between "RAG" and useful infrastructure documentation.
The archive is allowed to recover the wound.
It is not allowed to impersonate the runtime.

## Priority archive source map

Use this table when pulling imported conversations into a page.
It is deliberately source-level rather than topic-level.
The point is to keep each file attached to the exact burden it sharpens before
any prose tries to synthesize it.

| Source | Source family | Packet fields pressured | Strongest sentence supported | Stronger sentence still forbidden | Surviving private sentence |
| --- | --- | --- | --- | --- | --- |
| `source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md` | wrong-node-routing, placement truth | `entry_node`, `placement_source`, `selected_peer`, `peer_eligibility` | The user accepts manual placement but needs any receiving node to discover where the requested service currently lives. | Cloudflare DNS plurality plus per-node forwarding proves no-SPOF service routing. | I still personally know which node hosts the service and why that peer is eligible. |
| `source-archive/chatgpt-exports/conversations/docker-compose-frustration__695af0ff-0f74-8326-a73f-adcb574fa3b3.md` | operator-frustration, Compose legibility | `placement_source`, `forbidden_claim`, `operator_burden` | Compose opacity becomes painful when Docker hides lifecycle and ownership state behind project names, stopped containers, and global name conflicts. | Better Compose command hygiene solves the multi-node control-plane wound. | I still personally reconstruct what Docker really owns versus what Compose thinks it owns. |
| `source-archive/chatgpt-exports/conversations/docker-compose-multi-server-setup__67f73c50-150c-8006-8408-c03db2d8d287.md` | Compose distribution, remote-host glue | `placement_source`, `selected_peer`, `forbidden_claim` | Multi-server Compose patterns usually become remote sockets, sync scripts, or proxy glue unless a separate layer owns current truth. | A multi-server Compose recipe is already a distributed platform. | I still personally know which host is authoritative for each service. |
| `source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md` | fallback, advanced origin selection | `backend_loss`, `selected_peer`, `peer_eligibility`, `forbidden_claim` | The desired load balancer behavior includes health-aware origin choice, stickiness, circuit breaking, and fallback semantics beyond ordinary proxy presence. | A proxy or load-balancer-shaped service proves preserved request meaning under backend loss. | I still personally know what failover behavior should happen and whether it happened. |
| `source-archive/chatgpt-exports/conversations/traefik-service-failover-setup__689d5598-9720-832e-a891-ff57340bcd9c.md` | Traefik fallback, route durability | `backend_loss`, `policy_chain`, `selected_peer` | Traefik can express routing and fallback-shaped configuration, but syntax and runtime semantics are easy to overread. | A Traefik failover stanza proves route meaning, middleware, auth, and backend-loss behavior survived. | I still personally know whether the fallback path preserved the same service contract. |
| `source-archive/chatgpt-exports/conversations/distributed-ha-orchestration__685f4402-f304-8006-afcc-4802fd494bcc.md` | orchestrator-promotion, peer-equal coordination | `placement_source`, `peer_eligibility`, `forbidden_claim` | The user is looking for narrow coordination or peer-equal failure action before accepting a full orchestration worldview. | Mentioning Serf, Raft, Nomad, or gossip proves the repo should promote a new control plane. | I still personally know whether the candidate kills a specific hidden burden or only renames it. |
| `source-archive/chatgpt-exports/conversations/redis-url-and-load-balancing__68a914f8-d47c-8324-8734-bc1f17507bac.md` | raw-tcp, stateful-authority | `stateful_authority`, `backend_loss`, `forbidden_claim` | Redis is raw TCP and can be proxied at L4, but stateful authority is a stricter claim than reachability. | TCP forwarding proves Redis HA, writer safety, client rediscovery, or split-brain avoidance. | I still personally know who the writer is and whether failover is safe. |
| `source-archive/identity-exports/grok/conversations/traefik-ha-failover-without-swarm__b07a47ad-7c7a-4d7b-91f4-f5603dbf093a.md` | no-Swarm Traefik HA, wrong-node routing | `entry_node`, `selected_peer`, `policy_chain`, `forbidden_claim` | The same no-Swarm Traefik question recurs across providers, which makes the problem shape stronger but does not add runtime proof. | Cross-provider agreement proves current implementation maturity. | I still personally know which repeated advice was merely plausible versus locally proven. |

Do not treat this table as an exhaustive index of the archive.
Treat it as the current priority map for the repo's central wound.
If a new source matters more, add it by naming the packet field it pressures and
the stronger sentence it still forbids.

### How to read the table without laundering the source

Each row has to be read in both directions.

Forward reading:

- source file to source family
- source family to packet field
- packet field to page destination
- page destination to next proof requirement

Reverse reading:

- proposed page claim back to packet field
- packet field back to source family
- source family back to exact source file
- exact source file back to the forbidden upgrade it still cannot support

If only the forward direction exists, the page can cite sources while still
laundering them.
The reverse direction is what lets a later contributor downgrade a sentence
from proof back to pressure when it got too strong.

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
