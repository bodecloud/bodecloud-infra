# Compose-First Architecture

This page answers the question that keeps returning in different disguises:

> if the dream is already bigger than ordinary single-node Docker Compose, why
> is the repo still centered on
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> instead of already surrendering to Nomad, k3s, Kubernetes, or some other
> full control plane?

The answer is not:

- Compose is enough
- Compose solves distributed systems
- Compose is the forever platform

The answer is:

- Compose is still the clearest live authoring truth in the repo
- the user wants distributed behavior without prematurely surrendering
  readability
- the repo is trying to discover the smallest honest extra control surface
  before paying the tax of a heavyweight orchestrator

That is what "Compose-first" means here.

Not:

- "Compose solves clustering"
- "the adult answer is obviously Kubernetes"
- "distributed truth can stay implicit forever"

It means the repo is refusing to replace the mental model before isolating the
exact missing truths.

That refusal is not caution for its own sake.
It is a reaction to a repeated bait-and-switch:

- local Docker feels legible
- distribution turns that legibility into hidden operator glue
- the ecosystem responds by offering a larger worldview
- but the larger worldview often arrives before anyone has named which truth
  actually went missing

## The real argument for Compose-first

The real argument is not emotional attachment to YAML.

It is also not nostalgia for simple tooling.
The user is not trying to preserve Compose because stronger systems are scary in
the abstract.
The user is trying to preserve the last live authoring surface that still feels
inspectable before complexity disappears behind a farther-away control plane.

It is that the user is trying to preserve a very specific kind of operator
control:

- see what exists by reading the repo directly
- see how ingress, auth, healthchecks, configs, networks, and secrets are
  authored
- know which capability is missing before replacing the whole operating model
- avoid being trapped in a larger system whose abstractions answer different
  problems than the one this repo is actually facing

That is why "just use k3s" keeps reading as incomplete here.
It may eventually be correct.
It is not self-justifying.

That distinction has to stay visible because this page is one of the easiest
places for the site to become accidentally reassuring.

If "Compose-first" gets read as:

- the current stack is basically enough
- helpers already close the important gaps
- promotion is mostly preference

then the page has failed.

That is one of the most important distinctions in the repo:

- a stronger platform may eventually be right
- but "eventually right" is not the same as "already justified"

The repo keeps demanding a harder question first:

> which exact truth layer is missing, and can that layer be added without
> immediately replacing Compose as the main surface the operator reads?

That question is the entire point of Compose-first.

## What Compose-first does and does not mean

In `bolabaden-infra`, Compose-first means all of these are still true:

- the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  is the priority implementation surface
- the files under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
  are the active modular decomposition of that surface
- the live runtime should still be explainable from those files first
- helper layers should wrap, validate, enrich, generate from, or synchronize
  around that surface before replacing it

That last point is the real experiment.
The repo is not merely "stuck on Compose."
It is testing whether the missing truths can be added in a way that still keeps
Compose as the main human contract instead of downgrading it into a vestigial
artifact.

It does **not** mean:

- root Compose already expresses live cross-node placement truth
- root Compose already proves wrong-node success
- root Compose already proves failover, fallback, or stateful resilience
- the repo has already solved distributed behavior and merely prefers friendlier
  syntax
- intended `services.yaml` behavior should be narrated as if it were already a
  shipped tracked root contract

That distinction is one of the main reasons
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
matters more architecturally than `AGENTS.md`.

- `copilot-instructions.md` says what the runtime should feel like
- `AGENTS.md` says how to inspect and validate the repo

Both matter.
They do not carry the same architectural weight.

That distinction matters because earlier docs often blurred together:

- the desired runtime feeling
- the actual implementation surface
- the proof boundary

Once those blur, the docs immediately start sounding more complete than the
system actually is.

This page therefore has to keep one harsh question visible:

what exact missing truth is Compose still refusing to own, and what is the
smallest extra layer that could own it without making Compose a decorative lie?

## Why Compose is still the least dishonest live surface

Compose is still useful here for practical reasons, not nostalgic ones.

## 1. It keeps authoring truth visible

Many important truths are still directly readable:

- service definitions
- network attachments
- secrets references
- inline configs
- labels
- restart policies
- healthchecks
- mounts

That matters because the user keeps rejecting systems that hide too much
meaning behind a distant reconciliation layer before that hiding has earned its
keep.

Compose does not make distribution easy.
It does keep the current authoring truth close to the human eye.

That is not trivial.
It is one of the last sane surfaces the repo still has.

That sentence is blunt on purpose.
The user is reacting against a tool ecosystem where many "better" answers buy
coordination by making the real state harder to see until a new worldview has
already been accepted.

## 2. It handles heterogeneous workloads without inventing one giant worldview

The root runtime already spans several workload classes:

- public ingress and auth
- observability
- operator dashboards
- browser automation and crawling
- AI tooling
- state-bearing services
- mesh and identity-adjacent surfaces
- app and media workloads

Compose is genuinely good at expressing this kind of heterogeneity.

Its weakness is not service diversity.
Its weakness is distributed truth.

That is one of the most useful corrections this page can make.
Compose is not failing because the stack got broad or interesting.
Compose is failing because broad and interesting stacks eventually demand shared
truth that Compose does not naturally own.

That distinction matters because it explains why the repo is not running away
from Compose for authoring reasons.
It is searching beyond Compose for truth-ownership reasons.

## 3. It preserves Docker-native operator ergonomics

The repo still gets real value from:

- `docker compose config`
- direct service inspection
- label-driven Traefik wiring
- local iteration without cluster bootstrap or reconciliation delay

That is not a side benefit.
It is part of the design goal.

The user wants the system to become more resilient without immediately becoming
less legible.

This is the real meaning of Compose-first.

## The four truths Compose does not collapse for you

One reason earlier docs felt useless is that they kept calling several
different problems "the platform."

That flattening hid the actual missing layer.

These truths need to stay separate.

## 1. Authoring truth

Questions answered:

- what services exist?
- what labels, networks, secrets, configs, and healthchecks are authored?
- what does the root include graph activate?

This is where Compose is already strong.

## 2. Placement truth

Questions answered:

- what runs where right now?
- which node is local for a given service?
- what source of truth should a wrong-node receiver consult before forwarding?

This is where the repo-native concept of `services.yaml` matters.
The idea is clear across the README, instruction files, and planning docs.

But the tracked root implementation still does **not** ship a live root
[`services.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/services.yaml).

That means placement truth is still much stronger as architecture intent than
as runtime fact.

This is one of the central reasons the docs must stay careful.
The repo clearly knows it needs shared placement truth.
It still has not promoted that truth into the tracked root runtime.

That gap is exactly where the hidden human SPOF keeps surviving.
If the system still needs an operator to privately complete the answer to
"where does this service live right now?", then Compose-first has preserved
readability without yet preserving distributed request truth.

## 3. Traffic truth

Questions answered:

- how does a node expose services?
- how do hostnames, middleware, auth, and proxy paths behave?
- what logic is trying to preserve or regenerate routes?

This is where the repo is already operationally serious.

Live components already include:

- `traefik`
- `crowdsec`
- `tinyauth`
- `nginx-traefik-extensions`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `docker-gen-failover`
- `cloudflare-ddns`

This layer is closest to the dream of:

> any node can receive the request and still get it to the right place

But "closest" is not "proved."

It is also the easiest layer to over-credit emotionally.
Once a stack has Traefik, auth, middleware, failover helpers, and multiple
public nodes, readers start granting it missing truth by intuition.
This page has to keep refusing that gift.

## 4. Failure truth

Questions answered:

- if the local backend disappears, does the route needed for recovery stay
  alive?
- if the request lands on the wrong node, does that node know enough to forward
  correctly?
- if peer forwarding happens, do auth and middleware semantics survive?
- if a node disappears, is the system merely degraded or actually resilient?

This is the layer the user is least willing to let the docs fake.

Compose-first only stays honest if this layer remains painfully visible.

That pain is part of the value of this page.
If failure truth gets smoothed over, the repo starts sounding like all the
other infrastructure writing the user already distrusts.

## What Compose is not solving

Compose is strong at authoring truth and weak at distributed truth.

It does not natively solve:

- cross-node placement truth
- peer eligibility and health convergence
- route persistence under backend loss
- secret and env convergence across nodes
- revision awareness between nodes
- leader election or quorum
- stateful topology correctness

This weakness is not a minor edge case.
It is the central reason the repo keeps searching for more options.

And the options only count as real if they reduce this weakness without merely
moving it into a more opaque place.

Put more bluntly:

the user is not frustrated because Compose needs one more convenience wrapper.
The user is frustrated because Compose alone does not answer the sacred-machine
question:

> if traffic lands on the "wrong" box, what exact live truth lets that box stop
> being the wrong one?

That is the question Compose-first is trying to honor without cheating.

That is also why the user keeps sounding dissatisfied with ordinary advice.
Ordinary advice often solves an adjacent problem:

- better packaging
- better deployment ritual
- better platform posture

The user keeps returning to a stricter one:

- what concrete live truth lets the receiving node stop being semantically
  wrong?

That phrasing matters because it blocks a common downgrade.
The user is not simply asking for a forwarding mechanism.
They are asking for the removal of "wrongness" as a property that only a human
can interpret away after the fact.

## How the repo is stretching Compose without replacing it

The repo is not standing still.
It is already trying to build the missing middle around Compose.

## 1. Modular includes instead of one impossible file

The root file includes active fragments for:

- proxying
- docs
- metrics
- Headscale
- firecrawl
- LLM surfaces
- media surfaces
- warp routing
- wishlist surfaces

This keeps the live stack partitioned without pretending those partitions are
already independent cluster domains.

It is still one authored world.

That matters because the repo is trying to push complexity outward only where it
earns itself.
The authored world should not be replaced merely because a cluster API would
look more conventionally serious.

## 2. Rich ingress behavior around Docker-native labels

Traefik is not merely publishing ports here.
It is carrying a large part of the repo's distributed ambition:

- global names
- node-scoped names
- auth chaining
- middleware chaining
- failover generation pressure
- entrypoint differentiation

That matters because the repo is trying to make request-time behavior smarter
without immediately replacing Docker-native authoring.

## 3. Current-state-registry thinking without full scheduler promotion

The repeated `services.yaml` concept matters because it shows the repo is
trying to introduce one of the most important missing truths:

- where things actually live now

without immediately importing a full desired-state scheduler worldview.

This is one of the clearest expressions of the user's dream:

add the truth that has earned itself, not the whole empire surrounding that
truth.

It is also why Compose-first should not be misread as conservatism for its own
sake.
It is a demand that each added layer justify itself against a very specific
pain:

- hidden placement folklore
- hidden peer-eligibility folklore
- hidden bad-day reconstruction burden

That sentence may be the cleanest summary of Compose-first in the entire
knowledgebase.
It explains why the repo keeps reaching for narrow truth-bearing layers instead
of total platform replacement.

## 4. Helper and generator pressure instead of immediate platform capture

The repo has already explored ideas like:

- sync-agent
- failover-agent
- route generation
- OpenSVC-based ingress coordination
- Constellation state surfaces
- CUE-shaped control semantics

That is not random experimentation.
It is the repo repeatedly asking:

> can we buy the missing truths narrowly enough that Compose stays the main
> human contract?

That question is the heart of the Compose-first strategy.

## The real risk of Compose-first

Compose-first is not automatically virtuous.

It becomes dishonest the moment it is used as cover for unresolved burden.

Its real danger is that the repo may gradually rebuild a control plane in small
pieces while continuing to describe itself as "still just Compose."

There is a second danger too:
the repo may accumulate just enough helper intelligence that the docs start
speaking as if the missing truth has already been externalized when it has only
been spread across more components.

That would be dishonest for exactly the same reason premature Kubernetes
surrender can be dishonest:

the actual tax would no longer match the story being told.

So Compose-first only remains a credible strategy if the repo keeps asking:

- which truth does this helper own?
- how inspectable is that truth?
- how much new operator burden did this add?
- is this still a narrow helper, or have we quietly assembled a de facto
  orchestrator?

Those questions are not optional.
They are the safety rails.

Without them, Compose-first becomes one more fake option:

- still speaking like local Docker
- while secretly demanding cluster-grade reasoning from the operator

## The real reason the user keeps resisting immediate orchestrator promotion

The user is not just trying to avoid complexity.
They are trying to avoid paying for the wrong complexity too early.

What they want is not:

- "never use a stronger platform"

It is:

> do not force me to adopt a giant worldview until the missing truths are clear
> enough that I can say exactly what that worldview is buying me

That is a much more rigorous standard than ordinary self-hosting advice uses.

It is also the correct standard for this repo.

The user is not refusing growth.
They are refusing premature surrender to a worldview that cannot yet explain its
exchange rate clearly enough.

## Bottom line

Compose-first in `bolabaden-infra` does not mean Compose is sufficient.

It means Compose is still the least dishonest live authoring surface while the
project isolates the exact missing truths required for:

- wrong-node request preservation
- peer-aware routing
- convergence visibility
- failure-path persistence
- operator-readable distributed behavior

The moment another control layer can prove it owns those truths with less net
burden than the growing helper mesh around Compose, promotion becomes a serious
question.

Until then, Compose-first is not denial.
It is the repo's way of refusing to hide the problem behind a larger ideology
before the precise shape of the problem has been fully named.

That is why this page should keep feeling slightly unsatisfied.
Compose-first is only honest if it preserves the discomfort that the system
still lacks one clean shared answer on the bad day.

That is why this page should not read like a defense of Compose as a forever
answer.
It should read like a defense of intellectual honesty while the repo is still
trying to force a missing middle layer into existence.
