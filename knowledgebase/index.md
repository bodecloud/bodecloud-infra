# bolabaden Infrastructure Knowledgebase

This knowledgebase exists because `bolabaden-infra` is exactly the kind of
repo that can be summarized correctly and still understood completely wrong.

It also exists because the user is not merely under-informed.
They are frustrated by an ecosystem that keeps sounding like it offers choices
while repeatedly collapsing back into the same two insults:

- keep hand-maintaining private topology truth
- or accept a larger orchestrator worldview before it has proven which smaller
  missing layer actually mattered

If that emotional center disappears, the docs start sounding professional and
becoming useless again.

If a page says only:

- multi-node Docker
- failover
- anti-SPOF
- maybe Nomad
- maybe k3s
- maybe OpenSVC
- maybe helper agents

then that page has already failed.

That wording is too generic.
It turns the repo into an ordinary platform-comparison problem.
It erases the real pressure that made the repo exist in the first place.

The real question here is much sharper:

> how do you make several ordinary Docker nodes behave like one resilient
> personal cloud at request time, without lying about failover, without
> sacrificing readability too early, and without being forced into a whole
> orchestrator ideology just to stop one node from becoming sacred?

Everything in this site is supposed to orbit that question.

The site should also keep one harsher follow-up visible:

> when a request succeeds, did the system preserve the request, or did a human
> quietly preserve the architecture around it?

That is the question that stops a polished knowledgebase from becoming one more
beautifully organized version of the same old ambiguity.

## What this site is and is not allowed to prove

This site is allowed to:

- reconstruct the user's real anti-fake-options dream from the repo and archive
- separate runtime truth, intent truth, planning truth, and archive pressure
- show where the current implementation is strong, partial, contradictory, or
  still aspirational
- help a reader stop answering the smaller neighboring question by accident

This site is not allowed to:

- turn documentation clarity into an implied completion claim
- flatten competing futures into one falsely coherent platform story
- use architecture intent as proof that the runtime already behaves that way
- become easier to skim by quietly discarding contradiction, uneven maturity,
  or missing proof

## Quick site router

Start here if you want:

- the real benchmark and acceptance bar:
  [Operator Contract and Success Criteria](architecture/operator-contract.md)
- the literal wrong-node and degraded-backend request story:
  [Request Path and Failure Walkthrough](architecture/request-path-and-failure-walkthrough.md)
- the real runtime authority surfaces:
  [Instruction Surfaces and Authority](architecture/instruction-surfaces-and-authority.md)
  and [Current Compose Runtime](architecture/current-compose-runtime.md)
- the archive-grounded reconstruction of the dream:
  [User Intent and Dream](research/user-intent-and-dream.md)
- the hard proof boundaries:
  [Proof Matrix and Drill Catalog](operations/proof-matrix-and-drills.md)

## How to read this site without being fooled

This site is only useful if the reader keeps several truth registers separate
instead of letting them blend into one impressive-sounding story.

The main registers are:

- live runtime truth from the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  and the included fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- repo-native architecture intent from files such as
  [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  and [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- planning and promotion work from docs that describe what the repo is trying
  to earn next
- archive pressure from conversations and notes that reveal which kinds of
  answers the user keeps rejecting

Those registers are related, but they are not interchangeable.

If a page sounds persuasive because it quietly lets:

- intent impersonate runtime
- promotion work impersonate present capability
- archive frustration impersonate implementation proof
- runtime fragments impersonate solved architecture

then that page has failed, even if every sentence sounds reasonable.

This is the closest thing the knowledgebase has to a site-wide reading
contract:

- use the runtime to prove runtime claims
- use intent files to reconstruct the dream
- use planning files to explain what is being promoted
- use archive pressure to recover what simpler summaries would otherwise erase

If a page cannot say which register a claim belongs to, it is still too vague
for this repo.

There is a stricter rule underneath that sentence:

the site is not allowed to become clearer by making the user's dream smaller.
It is also not allowed to become cleaner by discarding the contradictions,
partial evidence, and competing futures that make the dream difficult.
In this repo, those messy edges are not editorial debris.
They are part of the truth surface.

That means this site has to be read as reconstruction, not as summary.

If the site becomes easier to skim by quietly merging together:

- live runtime truth
- architecture intent
- promotion work
- exploratory futures
- archive pressure

then it may become smoother while becoming less useful.
That is exactly the failure this rewrite is trying to prevent.

That is what the earlier docs kept doing.
They would preserve parts of the topology, parts of the plans, and parts of the
language, but they would quietly downgrade the real ask into something more
normal:

- better HA
- better routing
- better orchestration options
- better service discovery

Those are all subproblems.
They are not the full dream.

## Strongest honest current answer

This site is now much closer to a real reconstruction surface than to the
earlier polished-but-useless summary layer, but it only stays truthful when the
reader keeps the claim classes separate. The strongest honest answer for the
whole knowledgebase is that the dream is now being named far more faithfully
than before, while the runtime still only proves pieces of that dream unevenly
across HTTP, TCP, and stateful classes.

The full dream is closer to:

> build a distributed personal cloud whose intelligence no longer secretly
> lives in one operator's head, without surrendering the last readable
> operator-facing surfaces too early

If that reconstruction disappears, the docs may still be polished while
remaining useless.
That is the most important reading rule for the whole site:

- do not answer a smaller neighboring question just because it is easier
- do not swap the real wound for a tidy summary of related subproblems
- do not let rhetorical coherence outrank faithful reconstruction

## What a real answer would actually have to own

The site should not merely describe the dream.
It should keep the implied requirement stack visible.

A real answer would have to own at least these truths:

1. first-hop truth:
   multiple healthy public nodes can receive the first request without one
   silently remaining the sacred entrypoint
2. placement truth:
   the receiving node can determine where the requested service actually lives
   now, not where someone hoped it would live
3. eligibility truth:
   the system can distinguish a merely reachable peer from a peer that is
   semantically safe to forward to
4. route-survival truth:
   the fallback path needed during wrong-node or degraded-node behavior does
   not disappear with the local backend it was supposed to rescue
5. policy continuity truth:
   auth, middleware, headers, and request meaning survive the handoff instead
   of degrading into a looser emergency path
6. stateful honesty:
   stateful services are not narrated as solved just because a proxy can still
   reach something over the network
7. operator-legibility truth:
   the explanation for success or failure lives in inspectable shared surfaces
   rather than private memory and remembered repair rituals

If a proposed future only improves one or two of those truths, it may still be
useful, but it is not yet the full answer this repo keeps trying to recover.

This is why so much infrastructure writing feels irrelevant here.
It often stops after:

- ingress plurality
- service exposure
- controller presence
- health status
- or a respectable-looking HA diagram

Those can all exist while the operator still privately carries the most
important truths.

That is the difference between "more infrastructure" and "less hidden burden."

## What still does not count as solving the problem

The site also needs an explicit anti-cheat section, because this repo is full
of answers that can look mature while quietly preserving the same wound.

The dream is not satisfied by:

- Cloudflare being able to hit more than one node while the next hop is still
  a semantic gamble
- Traefik being present while wrong-node request preservation remains partial
  or unproven
- a placement registry existing on paper while the live runtime does not
  actually own and consume it
- helper automation existing while the real topology still has to be mentally
  reconstructed by the operator
- a stateful backend being reachable from more than one node while write
  ownership, storage truth, and failure semantics remain singular
- a respectable orchestrator being introduced before it has proven which
  smaller missing truths it actually relocates out of private memory

This is the negative benchmark the whole site should preserve:

if a reader leaves with more named options but still cannot tell whether the
wrong healthy node has stopped being architecturally humiliating, the docs are
still too smooth.

The archive pressure behind that question is not hypothetical.
It shows up directly in conversations where the user says, in effect:

- manual placement is fine
- Cloudflare multi-A first hop is already fine
- what is missing is unified service discovery and wrong-node request
  preservation
- the usual reverse-proxy and load-balancer names do not solve that by
  themselves

That is why this site should not read like a normal self-hosting README.
The user is not short on nouns.
They are short on answers that remain honest after the request hits the wrong
machine.

Another way to say it:

the repo is not merely searching for a better implementation.
It is searching for a more truthful category of answer than the ecosystem
normally offers.

That category would let the operator say:

- I still understand the system directly
- the first healthy node is not a semantic gamble
- I am not secretly the only place where the real topology lives

If the site stops preserving that deeper search, it can still look extremely
good while once again feeling useless to the user.

The priority implementation still starts from the real root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and the included fragments under
[`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose).
Nothing in this knowledgebase is allowed to outrank that worktree for runtime
claims.

## The dream is not "more infrastructure"

The user is not mainly saying:

- Docker Compose feels limited
- I need more cluster features
- I should probably upgrade to something bigger

The complaint is harsher than that.
It is closer to this:

- Cloudflare can send traffic to more than one node, but that still does not
  tell the receiving node what to do next
- Traefik can route and enforce policy, but that still does not prove the
  route survives the wrong failure
- a failover generator can exist, but that still does not mean the recovery
  path survives backend disappearance
- a service can be exposed through a stable hostname while still hiding one
  real writer, one real disk, or one real sacred backend
- "just use Kubernetes" sounds like paying a giant worldview tax before
  proving which smaller missing layer would actually solve the real pain

So the real deficit is not "features."
It is missing truth:

- placement truth
- peer eligibility truth
- convergence truth
- route-persistence truth
- stateful correctness truth

The system already has many components.
It still lacks enough *shared* truth to make those components feel like one
coherent service surface under the wrong request and the wrong failure.

That last phrase matters.
Many stacks feel coherent only while:

- requests happen to land locally
- the expected backend is still present
- the operator still remembers which host is special

This repo exists because that bargain feels fake.

That fakery is the emotional center of the whole project.

The user is not just saying "this is complicated."
They are saying:

> the system keeps feeling empowering right up until the moment wrong-node
> entry, fallback, and stateful reality matter, and then the only visible
> answers suddenly become either private glue or giant worldview capture

That is why the site should be judged by a harsher held-out test:

after reading it, can someone state which current options are real, which are
fake, and which still only exist as intent or promotion work?

If the answer is still no, then the site is still too polite for this repo.

That is why the docs have to preserve the pressure, not just the nouns.

The user is not short on component names, platform names, or general
availability vocabulary.
They are short on answers that remain honest once the request lands on the
wrong node and the explanation for success can no longer hide inside an
operator's memory.

## What the user is actually trying to make true

The instruction files, root README, planning docs, and archive pressure all
converge on the same desired runtime feeling:

1. any healthy public node can be the first hop
2. if the requested service is already on that node, it serves locally
3. if the request lands on the wrong healthy node, that node can still
   preserve the request by forwarding it to the correct healthy peer
4. auth, middleware, headers, and routing meaning survive that handoff
5. stateful systems are described with the same honesty as stateless ones

The dream is not satisfied because several boxes are online.
It is not satisfied because the proxy is green.
It is not satisfied because a route exists in config.

It is only satisfied when those five sentences stop being hopes and start
becoming trusted behavior.

And even then, one more condition remains:
the trusted behavior has to be explainable from inspected truth surfaces rather
than reconstructed folklore.

That is the standard the rest of this site is now trying to preserve.
Not:

- "looks clustered"
- "sounds orchestrated"
- "contains enough components to feel serious"

But:

- can the wrong healthy node stop being semantically wrong?

That is also why the site keeps returning to the repo's authority order.
For architecture intent, the clearest surface is still
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
For live implementation truth, the root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
still outrank calmer prose.
For repo-operability constraints, [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
matters.
For authoring discipline, [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)
matters.

If those roles get flattened into one voice, the whole site starts sounding
more implemented than the repo proves.

That sentence is one of the most important tests in the entire site.

It is not asking whether another node can answer something.
It is asking whether the platform has enough shared truth that the first
receiving node stops being an architectural gamble.

## The hidden enemy: operator reconstruction tax

This repo keeps circling one recurring human problem:

the system still depends too much on things the operator privately knows.

Examples:

- which node actually hosts a service right now
- whether a hostname is global or node-scoped
- whether a fallback route is live, generated, stale, or imaginary
- whether peers share the same env, secret, and revision state
- whether a reachable TCP frontend is hiding a fake HA story

That invisible reconstruction burden is one of the repo's main enemies.

So the deeper goal is not merely "availability."
It is:

> reduce the amount of hidden human truth required for several Docker nodes to
> behave like one request-preserving platform.

## If you only remember four things, remember these

This page should not require a second pass before the reader can answer the
main operator questions.

The shortest high-fidelity reading is:

1. the dream is not "better HA"; it is "stop making the operator's private
   topology memory the real control plane"
2. the live root Compose runtime is real, large, and serious, but it still
   does not broadly prove wrong-node request preservation or stateful
   correctness
3. the real missing layer is not "more tooling" in the abstract; it is shared
   placement, peer-eligibility, convergence, and route-persistence truth
4. many apparent options are still fake options because they only rename the
   same hidden burden under a different platform story

If a future rewrite makes those four points harder to recover, it is drifting
again no matter how polished it sounds.

## The fake-option filter

This repo has too many adjacent answers that sound different while preserving
the same wound.

Use this filter before treating anything as a real path forward.

| Apparent option | Why it sounds promising | Why it is still fake for this repo unless more is true |
| --- | --- | --- |
| More public DNS records | first-hop plurality sounds like resilience | the receiving node may still not know where the service lives or whether fallback preserves meaning |
| A stronger reverse proxy setup | proxy language sounds like routing truth | the proxy may still only understand local backends, stale peers, or happy-path routing |
| A service-registry idea on paper | it names the missing placement layer | the tracked root runtime may still not ship or consume it |
| A failover generator | it implies recovery behavior exists | the recovery route may still disappear when the backend disappears |
| A bigger orchestrator | it sounds like the grown-up answer | it may charge worldview tax before proving which exact truth layer it is actually removing |
| A reachable TCP frontend | it makes stateful services look globally present | it may still hide one real writer, one real disk path, or one real authority with weak failure semantics |

That table matters because the user is not just short on working code.
They are short on options that stay real after the first wrong-node request or
the first serious backend-loss event.

That is also why smoother docs can still be a regression here.
If a page lowers the felt ambiguity by hiding the remaining reconstruction tax,
it is making the site more pleasant while making the repo harder to reason
about honestly.

That is why the repo keeps pushing on the missing middle layer.

## What earlier documentation got wrong

The earlier documentation was not mainly too short.
It was too flattening.

It kept merging together:

- live root Compose reality
- architecture intent
- helper-agent and `services.yaml` ideas
- orchestration explorations
- archive frustration
- partial ingress evidence

until everything started sounding equally real.

That is exactly how docs become calm, polished, and useless.

That flattening was not just inaccurate.
It quietly taught the reader that the stack was closer to settlement than the
current worktree and proof surfaces actually support.

The operator does not need a smoother story.
The operator needs the docs to keep answering one hard question:

> if traffic lands on the wrong node during the wrong failure, what exact live
> truth keeps that request from dying, mutating, or quietly changing meaning?

If a page cannot answer that, it is still too shallow for this repo.

It is also too shallow if it answers a smaller neighboring question perfectly.

The site can still fail by being technically accurate about ingress, Compose
fragments, or orchestration candidates while still refusing to reconstruct the
architecture pressure that made those pages necessary in the first place.

That is also the most useful way to read all future architecture choices here.
The deciding question is not:

- is this tool lighter?
- is this tool more modern?
- is this tool more popular?

It is:

- which missing truth does this tool take ownership of, and does that
  ownership actually reduce hidden operator burden?

## The shortest honest reading of the repo

As of June 30, 2026:

- the dream is explicit
- the anti-SPOF instinct is explicit
- the live root Compose stack is real and substantial
- the edge stack is already serious
- the Compose-first preference is real
- the repo still does not broadly prove wrong-node request preservation
- the repo still does not broadly prove stateful zero-SPOF correctness

So the correct current reading is:

- this is a serious Compose-first infrastructure repo
- carrying real multi-node and failover pressure
- with a partially real and partially planned peer-aware routing story
- whose strongest proofs are still at the edge and authoring layers
- while the missing placement, convergence, and failure-truth layers remain the
  real obstacle

That should feel a little uncomfortable.
If it does not, the page is probably drifting back into architecture theater.
That discomfort is not a presentation flaw.
It is one of the few things keeping the docs aligned with the actual state of
the worktree and the actual severity of the user's complaint.

That discomfort is intentional.
The user is not asking for reassuring infrastructure prose.
They are asking for docs that can keep telling the truth even when the truth is
that the dream is sharper than the runtime, the edge is stronger than the
middle, and the missing truth layer is still the main obstacle.

## The authority order that should shape every high-level claim

The site becomes useless again the moment these source classes start getting
flattened.

For high-level pages, keep this order explicit:

1. current worktree and merged root Compose evidence for runtime claims
2. `.github/copilot-instructions.md` for the sharpest statement of the dream
3. planning docs for what the repo is seriously trying next
4. research and archive synthesis for why the usual ecosystem answers kept
   failing the user

That means:

- do not let a plan answer whether something is live
- do not let a runtime inventory answer what the user is really frustrated by
- do not let archive pressure answer what currently works
- do not let one smooth synthesis page outrank the stubborn negative artifact
  that still says the truth layer is missing

This entry page should keep that ranking visible because summary pages are
exactly where evidence classes quietly get merged into one reassuring voice.

## What victory would actually feel like

This site should keep repeating the end-state in operator language, because
that is the easiest thing for generic infra prose to erase.

Real progress would feel like this:

- first-hop node choice stops feeling like a gamble
- the operator can say where a service lives without consulting tribal memory
- a wrong-node request still feels like "the same service" after handoff
- backend loss does not destroy the route needed for recovery
- stateful services stop being described by edge symptoms and start being
  described by topology truth
- extra control-plane machinery, if adopted, feels earned instead of imposed

That is a much stricter target than "the cluster works."

## What the next real proof would have to look like

The next meaningful leap is not another philosophy page.
It is one route proving more of the dream than the current stack proves today.

The most valuable next proof would look like this:

- one stateless HTTP route
- request intentionally lands on a healthy node that does not host the service
  locally
- that node proves it knows the service is remote
- that node proves why the chosen peer is eligible now
- the request succeeds through the peer path
- auth, middleware, and visible request meaning are compared, not assumed
- the same route is later re-tested under backend-loss conditions so the route
  needed for recovery is not only a happy-path artifact

Anything weaker may still be useful engineering.
It is just not yet the proof that most directly answers the user's real demand.

## The live obstacle is still a missing truth layer

The useful way to read the repo is not "what features exist?"
It is "which truth layers are already strong enough, and which ones still
collapse back into operator memory or hopeful prose?"

### Strong enough to speak plainly about

- root Compose is the real authoring surface
- included fragments under `compose/` are the live decomposition of that
  surface
- Traefik-centered ingress is real
- auth and middleware surfaces are real
- observability surfaces are real
- HTTP and TCP exposure are both real
- the repo clearly wants local-first plus peer-forward behavior
- the repo clearly wants a lightweight shared placement truth instead of
  jumping straight to a heavyweight scheduler

### Still missing or still too weak to narrate as solved

- tracked live root placement truth such as `services.yaml`
- trustworthy wrong-node routing based on current placement truth
- peer eligibility and convergence truth
- route persistence when the local backend disappears
- cross-node certainty around env, secret, and revision coherence
- stateful failover semantics that deserve the phrase "high availability"

That list is the real backlog regardless of whether the eventual answer uses:

- better Compose-side helper logic
- generated configs
- OpenSVC
- Nomad
- k3s
- Kubernetes
- or some hybrid control layer

The user is not mainly demanding more options.
The user is demanding documentation that stops obscuring which exact option
closes which exact truth gap.

That is why the entry page cannot behave like a normal docs homepage.
Normal homepages compress.
This one has to discriminate:

- what is already real
- what is direction
- what is exploratory pressure
- what only sounds like an option because it still preserves the same hidden
  burden under cleaner wording

That is the difference between this site and the kind of README the user is
already tired of.
A normal README lists tools and components.
This one is supposed to preserve the shape of the wound.

## At-a-glance operator table

| Operator question | Best current answer | Confidence | Read next |
| --- | --- | --- | --- |
| What is the dream in one sentence? | Any healthy node should be able to receive a request, serve locally when possible, and preserve the request through a healthy peer when the service is remote, without fake HA language. | High | [`architecture/operator-contract.md`](architecture/operator-contract.md), [`research/user-intent-and-dream.md`](research/user-intent-and-dream.md) |
| What exact middle layer is missing? | A narrow placement, peer-eligibility, and route-persistence layer between raw Compose authoring and whole-platform promotion. | High | [`architecture/missing-middle-layer.md`](architecture/missing-middle-layer.md), [`operations/decision-paths-and-promotion-rules.md`](operations/decision-paths-and-promotion-rules.md) |
| What does the real root runtime currently prove? | A large live Compose stack with serious ingress, auth, observability, TCP/L7 exposure, and modular decomposition. | High | [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md), [`architecture/compose-fragment-map.md`](architecture/compose-fragment-map.md) |
| Is generic wrong-node HTTP success broadly proved? | No. The target contract is clear, but the proof chain is still incomplete. | High | [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md), [`architecture/ha-failover-routing.md`](architecture/ha-failover-routing.md) |
| Is tracked live placement truth already present at repo root? | No. The concept is central, but the tracked root worktree still does not expose it as live truth. | High | [`architecture/problem-and-goals.md`](architecture/problem-and-goals.md), [`research/evidence-ledger.md`](research/evidence-ledger.md) |
| Is stateful HA solved? | No. The docs now treat ingress continuity and data correctness as different claims on purpose. | High | [`architecture/stateful-ha-and-data.md`](architecture/stateful-ha-and-data.md), [`research/stateful-ha-evidence.md`](research/stateful-ha-evidence.md) |
| Do the docs now answer the recurring frustrated operator questions directly? | Yes. The Q&A and walkthrough pages now answer them directly instead of making the reader infer them. | High | [`operations/operator-questions-and-honest-answers.md`](operations/operator-questions-and-honest-answers.md), [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md) |
| Is the remaining problem “polish”? | No. The remaining problem is missing truth layers, not missing prose. | High | [`architecture/capability-gaps-and-roadmap.md`](architecture/capability-gaps-and-roadmap.md), [`research/evidence-ledger.md`](research/evidence-ledger.md) |

## Read the site as a pressure chain, not a taxonomy

Do not read this documentation like a feature encyclopedia.
Read it like a failure-analysis path.

That instruction is not a stylistic preference.
It is one of the main defenses against the exact failure the user is angry
about:

- tidy taxonomies make the space look richer than it is
- category pages make adjacent answers sound interchangeable
- the real missing layer vanishes into a map of related tools

This site should keep refusing that collapse.

### Step 1: understand the user's actual demand

- [`research/user-intent-and-dream.md`](research/user-intent-and-dream.md)
- [`architecture/operator-contract.md`](architecture/operator-contract.md)

These pages answer what the system is supposed to feel like from the operator
and request-path perspective.

### Step 2: force the request path to stay concrete

- [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md)
- [`operations/operator-questions-and-honest-answers.md`](operations/operator-questions-and-honest-answers.md)

These pages answer what must happen when traffic lands on the wrong node and
why first-hop reachability is not the same thing as preserved service success.

### Step 3: understand why Compose is still the live center

- [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
- [`architecture/compose-first-architecture.md`](architecture/compose-first-architecture.md)

These pages explain why Compose is still the least dishonest live authoring
surface even though the dream already exceeds plain local Compose.

### Step 4: understand the missing middle

- [`architecture/missing-middle-layer.md`](architecture/missing-middle-layer.md)
- [`operations/decision-paths-and-promotion-rules.md`](operations/decision-paths-and-promotion-rules.md)
- [`architecture/capability-gaps-and-roadmap.md`](architecture/capability-gaps-and-roadmap.md)

These pages answer which layer is missing between raw Compose and whole-platform
promotion, and when heavier orchestration actually earns its keep.

### Step 5: keep the honesty walls intact

- [`architecture/ha-failover-routing.md`](architecture/ha-failover-routing.md)
- [`architecture/stateful-ha-and-data.md`](architecture/stateful-ha-and-data.md)
- [`research/evidence-ledger.md`](research/evidence-ledger.md)

These pages stop the docs from quietly upgrading:

- DNS redundancy into end-to-end failover
- proxy presence into request preservation
- routability into stateful correctness

They also stop the site from pretending the user is basically asking for
"better HA."
They are asking for something stricter:

- a personal cloud that stops depending on private sacred-node memory
- without being bullied into platform ideology before the ideology has earned
  its keep

## The central distinction this site will keep repeating

In this repo, the most important distinction is still:

> this is the dream  
> this is the plan  
> this is what the tracked worktree proves today

Every major page in this knowledgebase exists to stop those three sentences
from collapsing into one smooth architecture story.

That is what "actually RAG this time" means here:

- recover the real demand surface from the instruction files and archive
- rank sources by authority
- preserve negative facts
- preserve unresolved contradictions
- and stop claims exactly where proof stops

That is what "actually RAG this time" means for this repo in practice:

- read the archive as pressure, not decoration
- read the instruction surfaces as intent hierarchy, not generic contributor
  boilerplate
- read `docker-compose.yml` as the current implementation anchor, not as a
  relic waiting to be superseded
- read every orchestration option against the same wrong-node and hidden-truth
  benchmark

## Bottom line

If an operator asks:

> can I already trust `bolabaden-infra` to behave like one resilient personal
> cloud across multiple ordinary Docker nodes?

the current best answer is:

- not yet generically
- but the repo is now much clearer about what has to become true before that
  claim stops being theater

That is the point of this documentation set.
