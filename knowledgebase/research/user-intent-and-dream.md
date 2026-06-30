# User Intent, Dream State, and Control-Plane Refusal

This page exists because most infrastructure docs avoid the only question that
actually decides whether they are useful for this repo:

> what does the operator want to feel true at request time, and which kinds of
> "solutions" already fail that standard even when they sound advanced?

For `bolabaden-infra`, that is not a side question.
It is the center of gravity of the entire documentation set.

If this page gets too polite, the rest of the site drifts back toward generic
DevOps narration.
That would miss the user's real demand entirely.

## The shortest possible honest reading

The user is not mainly trying to:

- host some services
- spread containers across a few nodes
- add a reverse proxy
- collect orchestration brands
- or make the stack sound more "serious"

The real goal is much sharper:

> build a personal cloud that stays manually understandable and
> Compose-readable, but behaves like one resilient distributed platform at
> request time, without paying heavyweight orchestrator tax before that tax has
> clearly justified itself

That is the dream the repo keeps circling.

A more exact reading is:

the user wants the system to become distributed without becoming emotionally
unowned.

That is why Compose readability matters.
That is why sacred-node memory is so offensive here.
That is why giant orchestrator sermons read as non-answers.

The user is not defending convenience for its own sake.
They are defending the last surfaces where the system still feels directly
legible and directly theirs.

That reading still needs one more layer of precision.

The user is not merely asking for "lighter orchestration."
That phrase is too weak and too easily co-opted by the same fake-option market
the repo is reacting against.

The sharper demand is:

> find or build an option that keeps causal legibility, reduces sacred-node
> dependence, reduces wrong-node humiliation, and exposes current truth without
> importing a whole scheduler worldview before that worldview has clearly
> earned trust

That is why the dream cannot be summarized as "something between Compose and
Kubernetes."
Plenty of things can sit between those nouns while still preserving the same
hidden operator burden.

## What this page is and is not allowed to prove

This page is authoritative about:

- the dream the repo is trying to preserve
- the negative benchmark the user keeps applying
- the kinds of answers that still fail even when they sound sophisticated

This page is not authoritative about:

- what the current root runtime already does successfully
- whether any specific future control layer has already earned promotion
- whether the tracked implementation already satisfies the dream

That boundary matters because this page is meant to recover the real ask, not
quietly bless the current state of the implementation.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "this is the dream the repo is protecting" | repo-native intent + archive pressure | `.github/copilot-instructions.md`, `README.md`, archive-derived research pages | that the dream is already live |
| "these are the user's anti-goals" | archive pressure + repo-native intent | `archive-pressure-patterns.md`, `README.md`, `.github/copilot-instructions.md` | that rejecting bad answers automatically identifies the winning good one |
| "Compose readability matters for a deeper reason" | repo-native intent | `.github/copilot-instructions.md`, `AGENTS.md`, root runtime surface | that Compose alone already solves cross-node truth |
| "this is what the platform should feel like on the bad day" | dream reconstruction | this page plus evidence-boundary pages | that the repo has already earned that feeling |

If a sentence about the user's dream starts sounding like runtime proof, it has
crossed the wrong boundary.

The user's real test is harsher:

> after the diagram is over and the bad day starts, does this still feel like
> a real option, or did it just rename the same dependence on private operator
> interpretation?

For recurring archive patterns that reinforce this reading, also see
[`archive-pressure-patterns.md`](archive-pressure-patterns.md).
For the repo's instruction-surface authority order, also see
[`../architecture/instruction-surfaces-and-authority.md`](../architecture/instruction-surfaces-and-authority.md).

## The documentation standard this dream requires

This repo does not merely need "better docs."
It needs a stronger reconstruction standard.

The closest useful analogy is:

- ordinary infra docs act like a shallow summary over a messy archive
- this repo needs something closer to authoritative retrieval over a messy
  archive

That means the docs have to preserve things ordinary architecture writing tries
to clean up:

- contradictions between dream and proof
- shifts in confidence across different files
- partial mechanisms that matter even though they do not yet compose into one
  calm story
- repeated archive pressure that reveals what the user is *actually* rejecting

If a page feels cleaner because it silently dropped those things, the page got
worse.

This matters because the user is not complaining about tone.
They are complaining about a deeper failure:

systems and docs keep sounding flexible until the real topology still has to
be reconstructed privately.

So this page should be read as a rule for the rest of the site:

> preserve the hidden negative benchmark, preserve the uneven evidence, and do
> not answer a smaller neighboring problem just because it produces smoother
> prose

It also has to preserve the difference between:

- a real option
- a fake option
- a future option that may earn reality later

That distinction is one of the main ways the site avoids sounding more hopeful
than the evidence supports.

## What the user wants the system to feel like

Most infrastructure docs stop at nouns.
This project only makes sense if it also states the desired runtime feeling.

The user wants the system to feel like this:

- traffic can land on any surviving node without that being a gamble
- a local service stays local instead of being swallowed by fake cluster ritual
- a remote service still works when the request lands on the wrong node
- the route used during failure still feels like the same service, not a
  semantically degraded fallback path
- operator surfaces remain readable instead of disappearing behind invisible
  control logic
- stateful systems are described with much harsher honesty than stateless ones

That feeling should be treated as a held-out evaluation surface.
The docs should keep asking:

- does a proposed answer reduce humiliation on the wrong node?
- does it externalize truth, or just automate around hidden truth?
- does it preserve causal legibility, or merely relocate it?
- does it sound like relief only because it renamed the same burden?

That runtime feeling is the real benchmark.
If a design sounds clever but does not produce that feeling, it is still
missing the point.

This matters because many answers can be technically competent and still feel
wrong in the exact way the user is rebelling against:

- first hop is still a gamble
- fallback still feels ceremonial
- stateful HA still feels like branding
- the operator still has to reconstruct the truth privately

The benchmark here is not only topological.
It is whether the platform stops behaving like a distributed illusion.

That means the dream is not satisfied by broader infrastructure language,
cleaner diagrams, or more named alternatives.
It is only satisfied when the wrong-node path, the fallback path, and the
stateful honesty path stop depending on private repair knowledge.

## The actual complaint

The user's frustration follows a repeated pattern:

1. Docker and Compose feel readable and empowering at small scale
2. the moment multi-node routing, placement, and failover matter, that
   readability mutates into manual glue and hidden truth
3. the conventional answer becomes:
   - just use Kubernetes
   - just use Swarm
   - just use a platform
4. the operator is pushed into a rotten trade:
   - keep brittle hand-maintained coordination
   - or accept a heavyweight desired-state worldview that removes some pain
     while importing a different class of complexity

This repo exists because the user does not accept that trade as the only adult
answer.

That is why the project should be read as a search for a missing middle:

- more dynamic than raw Compose
- more readable than a controller empire
- more honest than fake HA language

That "missing middle" phrase matters, but it is still not enough on its own.

The user is not mainly searching for:

- a market category
- a respectable new tool family
- a balanced-looking architecture diagram

The user is searching for a different honesty frontier.

They want more of the real topology burden to be carried by inspectable shared
truth and less of it to be carried by remembered placement, remembered rescue
rituals, and remembered sacred hosts.

That is why so many conventional answers still feel insulting.
They often do one of two things:

- preserve the hidden burden while improving the story
- or replace the story with a stronger worldview before proving the hidden
  burden was actually removed

The dream refuses both moves.

## The deepest emotional center of the project

Underneath the architecture vocabulary, this project is reacting to one very
specific pain:

self-hosting tools keep feeling empowering right up until the moment the system
depends on things the operator privately remembers.

That hidden dependency can look like:

- one remembered machine
- one remembered route
- one remembered placement fact
- one remembered secret-sync caveat
- one remembered "do not restart that here" warning
- one remembered topology truth the platform never exposed

That is why the repo keeps pushing so hard on:

- anti-sacred-node design
- no-fake-failover language
- lighter-than-Kubernetes coordination
- operator-readable truth

This is not just technical preference.
It is a refusal to accept systems that only feel flexible until distribution or
failure becomes real.

That is also why generic documentation styles fail here.
They describe components and plans but not the psychological break where
ownership turns into folklore.

That break is central to the dream.
If it is omitted, the repo gets misread.

## The exact anti-goals

To understand the dream, it helps to make the anti-goals explicit.

The user does **not** want:

- a single reverse-proxy node that becomes the real sacred box while the docs
  still say "multi-node"
- static upstream tables that still depend on private operator memory
- DNS redundancy narrated as if it were preserved service meaning
- a route that technically answers while silently dropping auth, middleware, or
  request identity
- a giant controller surface adopted mainly because the ecosystem says that is
  what "serious infrastructure" looks like
- stateful services narrated by liveness symptoms instead of topology truth
- a docs set that grants every proposed future equal emotional legitimacy just
  because each one can be described in sophisticated language

These refusals are not decorative.
They are the negative shape of the desired platform.

They also explain why the user sounds harsher than a typical self-hosting
operator.

The ecosystem keeps offering mature-sounding answers that add surface machinery
while leaving the same hidden sacred facts intact.
This repo exists as a refusal of that bait-and-switch.

## Strongest honest current answer

If a reader asks, "What is this repo actually trying to achieve?" the shortest
defensible answer is:

> It is trying to build a Compose-readable personal cloud that can behave like
> one resilient platform at request time, especially on the wrong node and on
> the bad day, without accepting either brittle static glue or heavyweight
> scheduler worldview tax before that tax has clearly earned itself.

That is a dream statement, not a runtime-completion statement.

## The clearest repo-native expression of the dream

The highest-signal files remain:

- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

They are not equal.

The blunt authority order is:

1. `copilot-instructions.md` names the dream most directly
2. `README.md` keeps the repo-level honesty wall around that dream
3. `AGENTS.md` anchors claims back to the live root runtime
4. `.cursorrules` mostly constrains authoring discipline
5. `INFRASTRUCTURE_MASTER_PLAN.md` names promoted future mechanisms

If those surfaces are flattened into one voice, the dream immediately starts
sounding more implemented than the worktree proves.

It also starts sounding more settled than the repo's real argument allows.

That matters because the docs are not only preserving "what the repo plans."
They are preserving a live dispute inside the repo:

- what the user actually refuses
- what the runtime really proves
- what the planning layer keeps trying to add
- and which futures may still turn out to be fake options once they are tested
  under wrong-node, fallback, and stateful pressure

If that live dispute disappears, the site can still look comprehensive while
quietly becoming much less faithful to the real ask.

Taken together, they define a very recognizable desired shape:

- no central orchestrator by default
- multi-node Docker as the main world
- Compose retained as the main authoring surface
- Cloudflare as first-hop node-entry infrastructure
- Traefik as the main L7 request surface
- separate treatment for L4/TCP and stateful concerns
- a lightweight current-state registry concept such as `services.yaml`
- helper agents or narrow control surfaces before wholesale platform capture

The sharpest language still lives in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md),
which explicitly describes:

- services manually assigned to nodes
- a current-state registry concept
- local serve if present
- peer proxy if remote

The master plan turns that dream into named future mechanisms:

- `bolabaden-sync-agent`
- `bolabaden-failover-agent`
- distributed placement truth
- multi-record Cloudflare DNS
- replacement of brittle failover generation

So the dream is not vague.
It is already implementation-shaped, even though the live runtime does not yet
prove the whole thing.

The docs fail whenever they blur those into one present-tense story:

- the dream is sharper than a brainstorm
- the plans are closer to the dream than the runtime is
- the runtime is still weaker than the dream

That unevenness is not a documentation defect.
It is one of the most important surviving truths in the repo.

## The behavior the user is actually asking for

Across repo-native docs, plans, and the archive, the desired operating contract
keeps converging on the same behavior.

## 1. Any healthy node can be the first hop

The system should not depend on one sacred public box.

Cloudflare multi-record DNS matters, but only as the first hop.
The deeper property is:

- any surviving public node can receive the request

That is the start of the dream, not the completion of it.

## 2. Local service should stay local

If the target service is already on the receiving node, the system should not
pretend it needs a full cluster worldview to serve it.

Why this matters:

- the fast path stays fast
- the mental model stays direct
- debugging stays local when locality is real
- the project avoids fake abstraction just to look modern

This is one of the subtle parts of the user's demand.
They are not asking for "distributed" as a performative aesthetic.
They are asking for distribution that does not erase the clarity of locality.

## 3. Wrong-node requests should still succeed

This is the true center of gravity of the repo.

The user keeps asking for a system where a request landing on the wrong node
still works because the receiving node can:

- determine that the target service is not local
- know which peer currently hosts it
- know that the peer is eligible now, not just in a stale config
- preserve the route when local failure is the reason fallback is needed
- preserve auth, middleware, headers, and policy continuity

This is also the best compression of the user's real benchmark.

Not:

- "can the site answer from more than one box?"

But:

- "can the wrong node stop being semantically wrong?"

That is the actual no-heavy-orchestrator distributed behavior the repo is
trying to force into existence.

This is also the place where ordinary docs begin cheating.
They replace "wrong-node success" with easier claims like:

- DNS reaches more than one box
- the proxy is up
- a peer is technically reachable
- the route still exists in config

The user is explicitly not asking for those watered-down versions.

That sentence should govern the rest of the site.

The watered-down versions are exactly the kinds of claims that sound useful
while still failing the real benchmark:

- more than one box can answer
- a fallback route exists
- the proxy is dynamic
- the platform is HA-capable

## 4. The operator surface stays readable

The user is not asking for less automation because they enjoy manual chores.
They are asking for readability because hidden control logic becomes its own
kind of tax.

The desired system should still be explainable in plain terms:

- what runs where
- what source of truth says that
- how a node decides between local serve and peer forward
- what disappears when a backend or node disappears
- what layer of complexity is actually paying for itself

That last point matters a lot.

The dream is not "never introduce more structure."
It is:

> never introduce more structure without being able to explain which real pain
> it is buying down

The user can tolerate complexity that shows its work.
They are rejecting complexity that arrives as branding, social pressure, or
fashion instead of as an honest answer to a named missing truth.

## 5. Stateful truth stays brutally honest

The user is openly hostile to fake HA vocabulary.

For state-bearing systems, the docs and architecture must answer:

- who owns writes
- what topology exists
- how promotion or election works
- how clients rediscover topology
- what really survives node loss

That is why this repo has to keep separating:

- ingress continuity
- request continuity
- stateful correctness

If those are flattened, the docs start sounding like every other stack that
confuses reachability with resilience.

## What the user is explicitly refusing

The dream is only half the picture.
The refusals matter just as much.

The user is refusing:

- sacred-node architectures that still cosplay as multi-node
- static truth disguised as dynamic routing
- DNS theater marketed as failover
- opaque control planes that replace one hidden burden with another
- platform adoption as a ritual act instead of a proven necessity
- docs that smooth over the difference between stateless paths and stateful
  ownership

These refusals are not side preferences.
They are the force shaping the whole repo.

## The real category the user is trying to create

What the user actually wants is not well-served by the usual categories.

They are trying to force a more honest category into view:

a personal cloud that:

- stays Compose-readable
- supports multi-node first-hop entry
- preserves wrong-node requests without pretending that DNS solved everything
- adds only the coordination truths that have truly earned their keep
- refuses to narrate stateful systems with the same softness as stateless ones

That is why "missing middle layer" is such a recurring concept.

It is not just a tool gap.
It is a category gap.

The user is reacting to the fact that too much of the current ecosystem says:

- local Docker is easy
- real distributed behavior belongs to heavyweight orchestrators
- therefore the space in between barely matters

This repo exists precisely because that in-between space matters.

## The simplest final summary

If this page has to collapse into one sentence, it should be this:

> the user is trying to build a personal cloud where multi-node Docker stops
> turning into fake options and hidden operator burden the moment requests can
> land on the wrong box, while also refusing to hand the whole problem to a
> heavyweight orchestrator before that surrender has clearly earned itself

That sentence should haunt every other page in the knowledgebase.

If a page forgets it, the docs will start getting gentle and generic again.
