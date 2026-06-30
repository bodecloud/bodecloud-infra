# bolabaden Infrastructure Knowledgebase

This site exists for one hard question:

> how do you keep
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> as the real human control surface, spread services across multiple ordinary
> Docker nodes, and still make wrong-node traffic, fallback, and anti-SPOF
> behavior feel like one coherent platform instead of one operator secretly
> remembering the real topology?

That is the real question.

This repo is not mainly about:

- generic self-hosting
- generic "high availability"
- tidy orchestrator comparison
- prettier Compose patterns
- collecting more modern infrastructure nouns

Those are all adjacent.
They are not the main wound this repo is trying to close.

## What this site is and is not allowed to prove

This site is authoritative about:

- the repo's actual dream
- the current root Compose implementation surface
- the difference between live runtime truth, repo-native intent, planning
  pressure, and archive pressure
- the concrete gaps between today's stack and genuine wrong-node recovery
- which claims still require proof before stronger language becomes legal

This site is not authoritative about:

- claiming that the current runtime already behaves like the dream
- turning a clear architecture explanation into failover proof
- promoting research or plans into shipped behavior
- narrating a larger control plane as already justified just because the
  current stack is uncomfortable

The site should help a reader leave with the right map of reality, not the
most optimistic story.

That means this front page is not allowed to become a high-gloss compression of
the whole repo.
It has to stay a routing surface that keeps the dream, the burden, and the
still-missing proof all visible at once.

If this page ever becomes the place where the repo starts sounding more solved
than it really is, then the nicest-looking page in the site will also be the
most misleading one.

## Strongest honest current answer

`bolabaden-infra` already contains a serious Compose-first platform:

- a real root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active includes under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- a substantial Traefik, CrowdSec, TinyAuth, and nginx-auth edge stack
- observability, maintenance, and operator surfaces
- private-mesh pressure through Headscale
- repeated architecture pressure toward any-node entry, peer-aware routing,
  and anti-SPOF behavior

What it still does **not** prove is the part the user actually cares about
most:

- that any healthy node can accept a request and preserve it correctly when
  the service is remote
- that placement truth is shared explicitly instead of remembered
- that peer eligibility is system-owned rather than guessed from reachability
- that fallback routes survive the failure they are meant to absorb
- that middleware, auth, and request semantics survive peer handoff
- that stateful services are truly resilient rather than merely reachable

The dream is clear.
The stack is real.
The missing truth-owning middle layer is still incomplete.

That three-part sentence has to stay intact.

Most weak summaries in this repo fail by dropping one clause:

- they keep the dream and the stack, then quietly forget the missing layer
- they keep the stack and the gaps, then quietly forget why the dream hurts
- they keep the dream and the gaps, then quietly flatten the current runtime
  into something thinner and less real than it is

This front page should keep all three clauses visible at once.

## What still does not count as a good front door for this repo

This page should also say more bluntly what fake orientation still looks like.

These still do not count:

- a broad overview that never names the hidden burden explicitly
- a technically correct summary that makes the problem sound like generic HA
- a polished entry page that leaves "wrong-node" as a detail instead of the
  humiliation threshold
- an options-rich tone that quietly implies the operator burden must already be
  falling
- a reading route that never says which stronger sentence is still forbidden

The user is not mainly asking for a better welcome page.
The user is asking for a front door that does not immediately downgrade the
real wound into a calmer neighboring question.

If this page sounds helpful while still making the repo feel more solved than
the evidence allows, then it is still part of the problem.

That danger is not hypothetical.
Polished ambiguity is one of the easiest ways infrastructure docs become
misleading without ever containing a single obviously false sentence.

A front page can stay technically accurate while still producing the wrong
emotional conclusion:

- there is a lot here
- the architecture is serious
- the gaps are clearly named
- therefore the hard part must already be mostly under control

This page is not allowed to produce that conclusion.

## What kind of site this is

This is not a normal architecture site.

Normal architecture sites try to:

- smooth contradictions
- merge intent and implementation into one calm voice
- present options as if the option space itself were proof of progress
- reward coherence even when the worktree still depends on hidden operator
  memory

This site has to do the opposite.

It has to preserve:

- the user's frustration with fake options
- the difference between a first hop and preserved request meaning
- the difference between a serious edge stack and a solved distributed system
- the difference between planning a registry and the runtime actually owning
  current placement truth
- the difference between reachable TCP and honest stateful authority

If those distinctions disappear, the site can sound broad and useful while
still teaching the same lie the user is tired of hearing:

> there are lots of options now, so the hidden burden must already be lower

## The dream in one brutally concrete scene

The dream is not abstract HA.
It is not "better clustering."
It is not "more infra options."

It is one concrete scene becoming normal:

1. a request lands on a healthy public node
2. that node does **not** host the target service locally
3. the node still knows what the request means
4. the node still knows which peer is eligible
5. the request still preserves auth, middleware, and service meaning
6. the operator does **not** have to privately remember the topology first

That is the scene.

The repo only gets to call itself closer to the dream when that scene becomes
more system-owned and less operator-owned.

If a page does not help a reader reason about that scene, it may still be
technically useful, but it is not yet close enough to the actual wound.

This scene is also the empathy test for the site.

If a page remains correct while making that scene feel secondary, already
manageable, or mostly handled by implication, then the page is still answering
the wrong human problem.

## What the user is actually rebelling against

The user is not mainly rebelling against Docker Compose itself.
They are rebelling against a repeating pattern in the self-hosting and infra
ecosystem:

- one option is "stay simple" but secretly means static glue plus private
  human reconstruction
- one option is "get serious" but secretly means accept a full orchestrator
  worldview before the missing truth was even named clearly
- everything in between sounds half-real until a bad day reveals the same
  sacred-node folklore underneath

That is why this site keeps returning to the phrase "missing middle."

The missing middle does **not** mean:

- a medium-complexity product
- a softer Kubernetes
- a nicer sync script
- a route generator with better branding

It means:

- a thinner truth-owning layer
- that removes the specific hidden burden the operator currently carries
- without charging more worldview tax than the burden is worth

That distinction is the center of the whole site.

## What a satisfying option would actually feel like

A satisfying option in this repo would not merely add features.
It would change where the truth lives.

It would make at least some of these statements become true:

- "what runs where right now?" can be answered by the system directly
- "can this node accept the request safely?" is not guessed from surface
  reachability alone
- "who should receive the request if local service is absent?" is no longer
  private operator folklore
- "does this fallback still preserve the request meaning?" becomes drillable
  instead of assumed
- "is this service actually resilient or only reachable?" becomes harder to
  lie about

That is why the repo keeps rejecting respectable but unsatisfying closure.

The docs should therefore help a reader feel the difference between:

- more moving parts
- more present-tense confidence
- and a genuine transfer of burden out of one human head

Only the third one counts.

That is also the rule for site polish and retrieval quality.

The knowledgebase should not be optimized toward:

- sounding broader
- sounding calmer
- sounding more authoritative in the abstract

unless those improvements also preserve the still-unmoved burden with equal
clarity.

## The four truth registers you have to keep separate

Most documentation drift in this repo happens when these get flattened into
one blended narrative.

### 1. Live runtime truth

Use this when the claim is:

> what actually exists in the priority implementation today?

Primary anchors:

- root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active include fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)

### 2. Repo-native dream and honesty truth

Use this when the claim is:

> what is the platform trying to become, and what honesty wall is the repo
> already insisting on?

Primary anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

### 3. Planning and promotion truth

Use this when the claim is:

> which missing layer has already been named, explored, or proposed?

Primary anchors:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)
- [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)
- [Stateful HA Evidence](research/stateful-ha-evidence.md)

### 4. Archive-pressure truth

Use this when the claim is:

> what is the user actually rebelling against, and why do normal answers keep
> feeling fake?

Primary anchors:

- [User Intent and Dream](research/user-intent-and-dream.md)
- [Archive Pressure Patterns](research/archive-pressure-patterns.md)
- [Evidence Ledger](research/evidence-ledger.md)
- [Source Assimilation Index](operations/source-assimilation-index.md)

These registers overlap.
They are not interchangeable.

If a sentence cannot say which register it belongs to, that sentence is still
too weak to trust.

## The shortest serious reading route

If you only have time for one pass and do not want to be fooled by calm
architecture language, use this order:

1. [User Intent and Dream](research/user-intent-and-dream.md)
2. [Problem and Goals](architecture/problem-and-goals.md)
3. [Operator Contract and Success Criteria](architecture/operator-contract.md)
4. [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
5. [Current Compose Runtime](architecture/current-compose-runtime.md)
6. [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
7. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
8. [The Missing Middle Layer](architecture/missing-middle-layer.md)
9. [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md)
10. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
11. [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

That route keeps all of these visible at once:

- what the user actually wants
- what the repo really contains
- where the hidden burden still lives
- what the missing layer actually is
- what still lacks proof
- why stateful surfaces stay under harsher rules

This route still must not be overread into:

- evidence that the route already survives the relevant failures
- proof that the missing middle is already selected
- permission to narrate the runtime as mostly converged

This is the shortest serious route, not the shortest flattering route.

If a reader follows it and still leaves thinking:

- the repo mainly needs polish
- the repo mainly needs a tool choice
- the repo mainly needs a few more integrations

then the route was either followed too casually or the front door still failed
to hold the real wound in place.

## Fast routes by real question

Use these when you already know what you need and want to avoid folder
browsing.

| If you need to know... | Start here |
| --- | --- |
| what the user is actually trying to make true | [User Intent and Dream](research/user-intent-and-dream.md) |
| what the root runtime really contains today | [Current Compose Runtime](architecture/current-compose-runtime.md) |
| which file most clearly states the multiple-node, no-Swarm-by-default dream | [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md) |
| why wrong-node entry is still the humiliating threshold | [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md) |
| why Cloudflare plus Traefik is still weaker than real failover | [HA, Failover, and Routing](architecture/ha-failover-routing.md) |
| what helper or control surface the repo is actually searching for | [The Missing Middle Layer](architecture/missing-middle-layer.md) |
| what success would actually have to mean | [Operator Contract and Success Criteria](architecture/operator-contract.md) |
| what still has to become system-owned before stronger claims are legal | [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md) |
| why Redis, MongoDB, Headscale, and similar services need harsher language | [Stateful HA and Data](architecture/stateful-ha-and-data.md) |
| what proof is still missing before the docs can speak more strongly | [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md) |

## What a real first-pass packet should leave behind

This front page should not merely orient.
It should help a reader leave behind an auditable first-pass packet.

A serious first pass should preserve:

- the real question the reader came with
- the truth register or registers they actually needed
- the strongest runtime-facing artifact they consulted
- the strongest dream- or archive-facing artifact they consulted
- the contradiction that still remained open after the first pass
- the next proof packet required before a stronger claim becomes honest

If the first pass only leaves behind "I understand the architecture better
now," it is still too weak for this repo.

It should leave behind something harsher and more useful:

- the exact humiliation threshold the system still fails at
- the exact hidden burden still living in operator memory
- the exact stronger sentence the docs are still not allowed to say
- the exact next proof packet that would be required before saying it

## The easiest bad read

The easiest way to misunderstand this repo is:

1. see Cloudflare and more than one public node
2. see Traefik, CrowdSec, TinyAuth, dashboards, helpers, and many services
3. see OpenSVC, Nomad, and k3s exploration
4. conclude that the remaining problem is mostly automation or polish

That reading is wrong.

The remaining problem is still the hard one:

- current placement truth
- convergence truth
- peer eligibility truth
- fallback-route durability
- cross-node semantic parity
- stateful honesty
- moving those truths out of one operator's head

This site is only useful if it keeps those gaps visible instead of smoothing
them away.

## The anti-cheat rule for the whole site

Before trusting any confident sentence in this knowledgebase, ask:

> is this sentence describing a real truth the current worktree owns, or is it
> describing a future layer, a research pressure, or an archive-derived
> complaint as if they were already one thing?

If the answer is unclear, the sentence is still too weak.

That weakness is not harmless here.
It is how polished documentation turns back into fake closure.
