# User Intent, Dream State, and Control-Plane Refusal

This page is the recovery surface for the question the rest of the site keeps
trying to answer too politely:

> what is the user actually trying to make feel normal, and what kinds of
> "solutions" still fail that standard even when they sound professional?

For `bolabaden-infra`, this is not flavor text.
It is the center of the project.

If this page gets watered down, the rest of the docs become technically
organized while still missing the real point:

- the user is not lacking tools
- the user is not lacking product categories
- the user is not lacking diagrams
- the user is lacking an option that still feels honest after the request lands
  on the wrong node

That last line is the whole reason this repository exists.

## The shortest honest reading

The user is not mainly trying to:

- spread some containers across multiple machines
- make the stack sound more distributed
- collect one more reverse proxy, scheduler, or cluster product
- upgrade from "small homelab" to "serious homelab" language

The real ask is much sharper:

> build a personal cloud that stays Compose-readable and operator-owned, but
> behaves like one resilient system at request time, without paying a giant
> scheduler tax before that scheduler has clearly earned itself.

That sentence still needs one more cut.

The user is not only asking for "something between Compose and Kubernetes."
That summary is too soft.
There are many things "between" those nouns that still preserve the same
private burden:

- the operator still remembers where the service really lives
- the operator still remembers which node is sacred
- the operator still remembers which fallback is currently safe
- the operator still remembers which stateful system is only pretending to be
  resilient

The dream is not "middleweight orchestration."
The dream is:

> remove sacred-node memory and wrong-node humiliation while preserving causal
> legibility.

That is the real benchmark.

## What this page is and is not allowed to prove

This page is authoritative about:

- the dream the repository is trying to preserve
- the anti-goals the user keeps rejecting
- the emotional and operational benchmark future work has to cross

This page is not authoritative about:

- what the current root runtime already proves
- whether any proposed control layer has already earned promotion
- whether the tracked implementation already satisfies the dream

This page is a demand reconstruction surface.
It is not a runtime proof page.

## Primary anchors for this page

Use these files in this order when reconstructing the dream:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`archive-pressure-patterns.md`](archive-pressure-patterns.md)
4. [`../architecture/instruction-surfaces-and-authority.md`](../architecture/instruction-surfaces-and-authority.md)

Why this order matters:

- `copilot-instructions.md` is the clearest statement of the dream
- `README.md` keeps the repo-level honesty boundary around that dream
- `archive-pressure-patterns.md` preserves what kinds of "answers" the user
  already considers fake
- the authority page explains why not every instruction surface should be read
  as architecture proof

## The dream in one concrete scene

The best way to reconstruct the user's intent is not with tooling names.
It is with a literal bad-day request:

1. Cloudflare sends a request to a healthy public node.
2. That node does not host the target service locally.
3. The receiving node still knows what service the request means.
4. The receiving node still knows which healthy peer is eligible now.
5. The request still preserves auth, middleware, and visible service meaning.
6. The operator does not have to reconstruct the answer privately first.

That is the dream.

If a proposed architecture still fails that scene, the user will correctly
experience it as the same old system wearing better clothes.

## What the user wants the platform to feel like

Most infrastructure writing stops at topology.
This repo only makes sense if it also preserves desired runtime feeling.

The user wants the platform to feel like this:

- traffic can land on any surviving node without that being a gamble
- a local service stays local when that is honest
- a remote service still works when the request lands on the wrong node
- fallback still feels like the same service, not a semantically degraded
  workaround
- middleware and auth do not silently disappear during fallback
- operator surfaces stay inspectable instead of becoming invisible cluster
  folklore
- stateful systems are described much more harshly than stateless ones

That feeling is not sentimental.
It is the held-out evaluation surface for the whole repo.

## What the user keeps rejecting

The user keeps rejecting two answer families.

### Rejected family 1: static glue that still leaves truth in the operator's head

Examples:

- hand-maintained per-node upstream tables
- hardcoded peer maps for every service
- "dynamic" patterns that still require manual edits for each placement change
- failover stories that only work while a specific person still remembers the
  topology

These answers fail because they preserve the hidden burden while improving the
story around it.

### Rejected family 2: heavyweight worldview capture before it has earned trust

Examples:

- "just use Kubernetes"
- "just use Swarm"
- "just use Nomad/Consul" without proving it preserves the desired operator
  surface
- platforms that solve one real pain while importing a much larger opaque
  control plane

These answers fail because they often relocate the burden rather than making it
more legible and more honestly owned.

The dream refuses both extremes.

## The real complaint under the architecture language

The user's recurring complaint is not just "self-hosting is hard."
It is more specific:

1. Docker and Compose feel empowering while the system is small.
2. The moment multi-node routing and failover matter, truth starts leaking into
   remembered host placement and ad hoc peer glue.
3. The surrounding ecosystem offers either brittle hand wiring or a heavy
   orchestrator worldview.
4. Neither option feels like a real answer to the actual wound.

That wound is:

> the platform stops being directly legible at exactly the moment resilience is
> supposed to become more real.

This is why the repo cares so much about Compose readability.
Not because YAML is sacred.
Because causal legibility is sacred.

## The hidden dependency the user is trying to kill

Underneath most of the repo's frustration is one specific pain:

self-hosting tools feel empowering right up until the system starts depending
on things the operator privately remembers.

That hidden dependency can look like:

- one remembered public entrypoint
- one remembered placement truth
- one remembered "safe peer"
- one remembered rescue ritual
- one remembered warning about a stateful service that only looks portable

The user wants less of the system's truth to live there.

That is why many ordinary answers still feel insulting here.
They improve:

- naming
- automation
- presentation
- cluster vocabulary

while leaving the same private truth burden intact.

## What counts as a real option in this repo

For this project, an option is only real if it makes at least one of these
things less true:

- wrong-node entry still collapses back into private operator knowledge
- fallback still depends on remembered placement
- auth and middleware still become uncertain during handoff
- stateful resilience is still mostly branding
- the operator still cannot answer "what runs where right now?" from shared
  inspectable truth

If a proposed path does not materially reduce one of those burdens, then from
the user's point of view it is mostly theater even if it is technically
respectable.

## What this page should force the rest of the site to do

Every serious page in the knowledgebase should preserve all of the following:

- the real request-time benchmark
- the distinction between live proof and dream reconstruction
- the fact that wrong-node dignity is a first-class goal
- the fact that stateful honesty is stricter than stateless continuity
- the fact that larger control planes must earn their opacity

If a page becomes smoother by shrinking one of those things, it got worse.

## Strongest honest current answer

The dream is not ambiguous.
The user wants a multi-node Docker system that behaves like one resilient
platform at request time without defaulting to heavy orchestration, and they
specifically want to stop acting as the hidden control plane when requests land
on the wrong node or a backend disappears.

The current repo absolutely preserves that dream.
It does not yet prove the dream is live.

That distinction has to stay visible everywhere.
