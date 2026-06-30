# Orchestration Research 2026

> **Reading boundary**: this file is a legacy planning and research artifact.
> It is useful because it shows what kinds of extra coordination layers the repo
> keeps exploring. It is **not** proof that the live root runtime has already
> adopted one of these layers as its authoritative control plane.
>
> For the current evidence-first reading, start with:
>
> - [`../knowledgebase/architecture/orchestration-options.md`](../knowledgebase/architecture/orchestration-options.md)
> - [`../knowledgebase/research/orchestration-research-2026.md`](../knowledgebase/research/orchestration-research-2026.md)
> - [`../knowledgebase/research/orchestrator-tradeoffs-evidence.md`](../knowledgebase/research/orchestrator-tradeoffs-evidence.md)

# What This Page Is Actually About

This page only helps if it stays tied to the real repo question.

The question is **not**:

> which orchestrator is most fashionable or most feature-rich?

The question is:

> what is the smallest honest coordination layer that can make multiple
> ordinary Docker nodes stop behaving like loosely related boxes, without
> forcing the operator into heavyweight platform capture before it has actually
> earned its keep?

That difference matters.

`bolabaden-infra` is not choosing from a blank sheet. The repo already has:

- a real root [`docker-compose.yml`](../docker-compose.yml)
- active includes under [`../compose/`](../compose/)
- a public ingress surface that wants any-node entry
- strong pressure toward local-first serving and peer-forward fallback
- strong refusal to call something "HA" when it only preserves reachability

So this page is not a product ranking.
It is a pressure map.

It should also be read as a record of unresolved insistence.
The repo keeps returning here because the normal market categories keep failing
to describe the exact thing the user is still asking for.
If this file is reduced to a calmer survey of options, it starts preserving the
same nouns while losing the actual question.

## The Pressure That Keeps Producing Orchestration Research

The repo keeps returning to orchestration questions because plain Compose,
while still valuable, does not fully answer the problems the user actually
cares about.

Those problems are harsher than generic "scaling":

- traffic can land on the wrong node
- the receiving node may not know enough to recover intelligently
- route-generation logic can disappear at the moment fallback is needed
- env and secret truth can drift across nodes
- the operator still has to privately remember too much placement truth
- stateful services can look reachable without becoming resilient

That is why the repo keeps exploring extra layers instead of just adding more
Compose fragments forever.

This is also why the research has to stay hostile to respectable-sounding
closure.
A platform family does not become a real answer here because:

- it sounds mature
- it is commonly adopted
- it makes the architecture diagram look more adult

Those are all ways to improve narration faster than burden ownership.

## What the User Is Actually Rejecting

The orchestration search only makes sense if the refusals stay explicit.

The user is repeatedly rejecting:

### 1. Static glue narrated as dynamic infrastructure

Examples:

- manually syncing files between nodes
- predeclaring backends and calling it failover
- templating static data and calling it service discovery

That does not solve the underlying wrong-node and drift problem.

This page should keep one harsher version of that sentence visible:
static glue becomes especially dangerous when the docs start narrating it as if
the system had already crossed into shared-truth territory.

### 2. Fake HA that stops at DNS or a proxy

Examples:

- multiple A records alone
- a still-routed hostname alone
- a reverse proxy in front of one real failure domain

Those may improve reachability.
They do not automatically improve correctness.

They also do not automatically improve dignity.
That matters in this repo because "it still answered something" is often much
less than the user is actually asking the system to preserve.

### 3. Heavyweight platform sermons that ignore the actual pain

"Just use Kubernetes" is only an answer if Kubernetes solves the right
pressure more honestly than the current stack can.

If the actual pain is still:

- missing placement truth
- brittle route persistence
- secret drift
- unclear node ownership

then jumping straight to a giant scheduler can still be the wrong move.

This page should therefore resist two equal and opposite simplifications:

- "eventually everyone serious lands on Kubernetes anyway"
- "anything bigger than Compose is philosophical betrayal"

Both moves shrink the user's real search into a more ordinary culture war than
the repo is actually conducting.

### 4. A false choice between raw Compose pain and total platform capture

This is probably the deepest pattern in the repo.

The whole research surface exists because the ecosystem too often offers only
two unsatisfying extremes:

- stay in hand-managed Compose sprawl
- accept a full new worldview all at once

The repo is trying to find the missing middle.

But even "missing middle" is not enough if it becomes just another category
label.
The real test is whether a candidate layer changes who owns the bad-day truth,
not merely where the layer sits on a complexity spectrum.

## The Real Things the Repo Is Trying to Buy

Different tools keep appearing, but the desired purchases are fairly stable.

### Placement truth

The system needs a trustworthy answer to:

- what runs where right now
- what may move
- who should answer the request if the local node is wrong

Without that, the operator remains the hidden scheduler.

That phrase should be taken literally.
One of the deepest anti-SPOF themes in this repo is that a human head can be
the most dangerous unreplicated control plane in the whole system.

### Convergence truth

The system needs a trustworthy answer to:

- what config and secret material exists on each node
- who detects drift
- who regenerates route state
- who repairs obvious mismatch

Without that, "distributed" often just means "the same mess copied around."

This repo is full of pressure against exactly that outcome:
multi-node duplication without multi-node truth still feels distributed in the
most misleading way possible.

### Ingress truth

The system needs a trustworthy answer to:

- whether any node can safely receive traffic
- whether the receiving node can distinguish local serve from peer forward
- whether the route survives the local failure that triggered recovery
- whether middleware, auth, and policy stay coherent across the handoff

### Failover truth

The system needs a trustworthy answer to:

- who decides the next backend
- how that decision survives backend disappearance
- whether the system converges after handoff rather than silently degrading

### Lower operator reconstruction tax

This is one of the least glamorous and most important goals.

The user does not want to keep reconstructing:

- current node placement
- hostname scope
- service affinity
- stateful write ownership
- which layer is responsible when something fails

If a new layer does not reduce that burden, it has not earned itself.

That is the hidden metric behind almost every option discussion in this repo.
The question is not only "what can this platform do?"
It is also "what private explanatory labor does this platform stop asking the
operator to perform?"

## Option Families, in Repo Terms

This repo should evaluate orchestration families by what missing truth they
would own, not by how impressive their marketing page sounds.

## 1. Compose plus coordination glue

Representative repo-shaped ideas:

- `services.yaml`
- sync-agent concepts
- failover-agent concepts
- generated Traefik file-provider config
- DDNS and node-aware forwarding helpers
- Constellation-style state and convergence logic

### Why this family keeps returning

Because it preserves what the repo still values:

- Docker-native authoring
- readable service definitions
- incremental evolution
- low migration shock

### What it is good at

- authoring readability
- local iteration
- explicit operator control
- gradual extension of a Compose-first stack

### What it risks

- becoming a shadow control plane without admitting it
- scattering truth across helpers, templates, agents, and conventions
- preserving file readability while losing runtime clarity

This page should not underrate that risk.
Shadow control planes are one of the easiest ways to keep the same pain while
telling a much nicer story about it.

### When this family is still the right answer

Stay here when the dominant pain is still:

- missing truth layers
- brittle edge behavior
- convergence gaps

rather than full scheduler pressure.

## 2. OpenSVC or similar infra-grade HA tooling

This family matters because some of the repo's problems are not primarily app
orchestrator problems. They are boring infra coordination problems:

- node health
- narrow service relocation
- explicit failover ownership
- cluster membership
- stable infrastructure-grade handoff primitives

### Why it fits the repo better than generic "just use Kubernetes"

OpenSVC is interesting here because it can strengthen a specific layer:

- node membership
- service supervision
- ingress or L4 generation inputs
- infra-grade relocation primitives

without forcing the whole repo to stop being Compose-first all at once.

### What it does not solve by itself

It does **not** automatically prove:

- that the root runtime is now OpenSVC-governed end to end
- that stateful HA became solved
- that all app placement truth is now unified
- that wrong-node HTTP success is globally complete

This matters because the repo must not promote infra-grade primitives into
fake universal proof.

That warning matters beyond OpenSVC specifically.
One of the repo's recurring failure modes is taking something that genuinely
helps one narrow truth and then narrating it as if it had healed the whole
stack's honesty problem.

## 3. Nomad and similar scheduler families

This family matters when the main pressure becomes:

- workload placement
- rescheduling
- service registration
- cleaner runtime reconciliation

Nomad is attractive here because it often feels closer to the missing middle
than Kubernetes does.

### Why it is interesting

- lighter operator worldview than Kubernetes
- stronger scheduling semantics than raw Compose
- good fit when placement and runtime reconciliation become dominant pain

### Why it is not yet the default repo answer

The current repo still shows multiple plausible futures:

- stronger Compose helpers
- OpenSVC-backed infra coordination
- schema-first metadata plus generation
- agent-managed convergence
- scheduler promotion

Nomad may become right later, but the current live repo does not yet prove
that whole-stack scheduler promotion is the cleanest next move.

## 4. k3s or Kubernetes families

These become reasonable when the dominant unsolved pain is undeniably
cluster-shaped:

- service discovery
- reconciliation
- standardized scheduling
- ecosystem integration
- uniform operational primitives

### Why this family keeps appearing anyway

Because the ecosystem is strong here, and because many of the missing features
really do exist in mature form in Kubernetes land.

### Why the repo is still resisting default promotion

Because the repo is not trying to win a popularity contest.
It is trying to avoid paying:

- a large control-plane tax
- a storage and networking worldview shift
- a heavier operator model

before the stack has proved that this is actually the right pain to solve next.

## What This Page Should and Should Not Make You Conclude

### It should make you conclude:

- the repo is not indecisive for no reason
- the repo is trying to solve a real missing-middle problem
- platform choice here has to be earned by a named failure class
- the current root runtime is still Compose-first

### It should **not** make you conclude:

- that the repo has already converged on one orchestrator future
- that OpenSVC, Nomad, or Kubernetes has already been promoted into the live
  root control plane
- that stateful HA becomes solved once a scheduler is chosen
- that research and planning documents are runtime proof

## Bottom Line

The orchestration research in this repo is best read as a refusal to accept
either of these lies:

- "just keep adding Compose and scripts forever"
- "just surrender to a giant platform now and call the problem solved"

The real search is for the smallest layer that can honestly absorb:

- placement truth
- convergence truth
- ingress truth
- failover truth

without erasing the parts of Docker and Compose that still make the stack
readable.

That search is still active.
The live runtime has **not** yet proven that one family has decisively won.
