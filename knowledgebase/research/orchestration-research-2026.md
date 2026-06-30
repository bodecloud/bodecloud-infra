# Orchestration Research 2026: What the Repo Is Actually Trying to Buy

This page is not a generic orchestrator comparison.

If it becomes "Compose vs Nomad vs Kubernetes," it stops describing the repo.

The actual research question inside `bolabaden-infra` is much narrower and much
more hostile than that:

> what is the smallest honest layer that can make multiple ordinary Docker
> nodes behave less stupidly under drift, misplacement, wrong-node entry, and
> failure, without forcing the operator to pay a giant worldview tax before
> that tax has clearly earned itself?

That is the research question.

There is also a deeper version that this page has to protect:

> what kind of system is the user actually trying to bring into existence,
> once all the easier substitute questions are stripped away?

That matters because the orchestration search in this repo is not shopping.
It is an attempt to escape fake options.

The fake options are familiar:

- keep piling static glue onto Compose and pretend generation made it dynamic
- accept first-hop plurality as if it solved request preservation
- import a full scheduler worldview before proving the smaller missing truth
  cannot be carried another way

This page only helps if it keeps that trap visible.

## Strongest honest current answer

If a reader asks what the orchestration research is actually concluding right
now, the shortest defensible answer is:

> the repo is not yet choosing a winner so much as ranking which missing truth
> layers still force the operator to carry the platform privately. Compose
> remains the live human contract, the archive keeps rejecting both static-glue
> theater and premature worldview capture, and any stronger candidate still
> has to prove that it transfers the right burden at an acceptable readability
> and control-plane cost.

That is a research conclusion, not a product endorsement.

## What this page is and is not allowed to prove

This page is allowed to prove:

1. the orchestration search here is really a search for the smallest honest
   owner of missing distributed truth
2. the user is rejecting both static-glue theater and premature full-cluster
   worldview capture
3. Compose still matters because it remains the live human contract
4. candidate families should be ranked by burden transfer, not platform
   prestige
5. the research already knows what kinds of humiliations a promoted layer must
   kill

This page is not allowed to prove:

- that one candidate family has already won
- that a stronger control plane is automatically justified
- that a smaller layer will be sufficient forever
- that research coherence equals implementation proof

## Retrieval contract for this page

### Class 1: live-root evidence

Primary anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active Compose fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)

This class is allowed to prove:

- what the priority implementation still is
- why Compose is the live baseline being protected

It is not allowed to prove:

- that the current runtime already satisfies the desired distributed behavior

### Class 2: repo-direction evidence

Primary anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`infra/docs/ARCHITECTURE.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/infra/docs/ARCHITECTURE.md)

This class is allowed to prove:

- what kind of system the repo is trying to grow into
- that the repo is searching for thinner truth-owning layers before
  heavyweight promotion

It is not allowed to prove:

- that the desired layer already exists in the tracked runtime

### Class 3: archive-pressure evidence

Primary anchors:

- archive-pressure synthesis pages
- imported archive conversations about wrong-node routing, orchestration, and
  frustration

This class is allowed to prove:

- what the user keeps rejecting
- why ordinary comparison logic is too weak for this repo

It is not allowed to prove:

- that current runtime behavior has crossed the archive's threshold

### Class 4: non-proof planning and exploration

Primary anchors:

- option pages
- exploration pages
- design and blueprint docs

This class is allowed to clarify direction.
It is not allowed to impersonate runtime ownership.

If a paragraph starts using Class 2 or Class 3 force as if it were Class 1
force, the page has already become too flattering.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "the research is trying to buy missing truth, not prestige" | repo-direction plus archive pressure | `.github/copilot-instructions.md`, master plan, archive-pressure pages | that the winning purchase has already been made |
| "the live root runtime constrains the option space" | live-root evidence | root `docker-compose.yml`, active Compose fragments | that the current runtime already satisfies the desired behavior |
| "the user rejects static glue and premature worldview capture" | archive pressure | archive-pressure pages and archive conversations | that no stronger control plane will ever be justified |
| "a candidate family matters" | planning and research framing | this page, option pages, exploratory artifacts | that the candidate is already promoted or production-trusted |

If the page starts sounding like a product verdict without naming the burden
being transferred, it has become too shallow for this repo.

## What the user is actually rejecting

The orchestration research only makes sense if the refusals stay explicit.

The user is not merely saying:

- I want HA
- I want clustering
- I want orchestration

The user is rejecting a repeating pattern of fake choice.

### Refusal 1: static glue disguised as dynamic infrastructure

The repo repeatedly runs into answers that still cash out as:

- predeclare everything
- manually sync the files
- keep a human in the loop for topology truth
- call it dynamic because a template or generator was involved

That is not the dream here.

The user is explicitly trying to escape systems that only look smarter than
static config while still depending on hidden operator knowledge.

This is why templating, generation, or GitOps language by itself cannot close
the wound.
Those may all help.
They are still near-misses if they do not create a new shared runtime truth
that the wrong node can actually rely on under stress.

### Refusal 2: wrong-node traffic that still depends on sacred human memory

The desired system is not just "load balanced."
It is supposed to behave like this:

- any healthy node can receive the request
- the receiving node can tell whether the target is local
- if not local, it can preserve the request through a legitimate remote path
- the operator does not privately remember all placement truth to make that
  work

That is a stricter demand than ordinary reverse-proxy redundancy.

If the receiving node still depends on:

- implied topology
- stale registry truth
- operator-maintained rescue conventions
- semantically uncertain middleware inheritance

then the system is still asking the user to do the same hidden labor beneath a
larger stack.

The bar is not "route traffic somewhere healthy."
The bar is "stop the wrong receiving node from becoming a semantic error."

### Refusal 3: "just use Kubernetes" as sermon instead of answer

The repo is not anti-Kubernetes because Kubernetes is fake.

It is anti-Kubernetes-as-default-sermon because the user does not want to pay:

- control-plane tax
- storage and state tax
- workflow tax
- worldview tax

before there is solid evidence that smaller, more legible answers cannot carry
the required truths.

If a big orchestrator is the right answer, this repo wants the justification to
sound like:

- here is the exact hidden truth it owns better
- here is the exact humiliation that stops being true

Until then, "just use Kubernetes" is not rigor.
It is evasion.

### Refusal 4: fake HA that never reaches placement, policy, or state truth

The repo repeatedly rejects resilience stories that stop at:

- multiple public nodes
- DNS failover
- green health indicators
- route existence

and never answer:

- what runs where right now
- what happens on the wrong node
- whether middleware and auth survive the handoff
- whether state continuity is real

Without this refusal, the whole project degrades into infrastructure theater.

### Refusal 5: the forced choice between raw Compose pain and total platform capture

The ecosystem keeps offering two unsatisfying extremes:

- stay in hand-managed Compose sprawl
- surrender into a full desired-state cluster worldview

The orchestration research exists because the user does not accept that as the
only serious choice set.

The missing middle is the real object of research.
Not a brand.
Not a benchmark leader.
Not a currently fashionable controller.

The object is:

the smallest additional truth-owning layer that stops the system from
functioning as a distributed illusion.

## The four plus one questions every orchestration answer must survive

Every candidate family must survive these questions:

1. what hidden burden moves out of the operator's head if this layer becomes
   real?
2. what part of the wrong-node problem does it actually solve?
3. what new worldview or control-plane tax does it import?
4. what part of the user's real dream still remains unowned even if it works
   exactly as advertised?
5. is the repo about to buy a control plane larger than the exact wound it is
   healing?

If a candidate cannot survive those questions, it is still closer to platform
theater than to a real answer for this repo.

## The orchestration anti-benchmark

The repo is not evaluating candidates against a fantasy perfect platform.
It is evaluating them against a harsher anti-benchmark:

- the request still only works if it lands on the unofficial right node
- the operator still privately remembers where the service really lives
- fallback still looks real only while the preferred backend is healthy
- stateful systems still gain a nicer edge story without gaining trustworthy
  authority or replication truth

If a candidate mostly preserves those conditions, then even a respectable
platform is still a weak answer in this repo's terms.

## What the current worktree already proves before any promotion

### Strong proof: Compose is still the human contract

The root runtime is still the root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml),
not Nomad jobs, not Kubernetes manifests, and not an OpenSVC-owned runtime.

That matters because every future move is being judged against a Compose-first
baseline, not a clean-room abstraction contest.

The thing being protected is not Docker branding.
It is the readability and directness Compose still gives before distribution
becomes real.

### Strong proof: the stack has already outgrown naive single-node mental models

The active include set and root-owned services already span:

- public ingress
- observability
- private coordination
- heterogeneous app clusters
- stateful services
- network and egress experiments

That means the orchestration problem is not a thought exercise.
The repo has already entered the zone where:

- operator memory
- static assumptions
- node-local hacks

become architecture tax.

### Strong proof: the repo wants any-node entry semantics

[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
states the intended flow directly:

- Cloudflare DNS points to multiple nodes
- any node may receive the request
- if the service is local, serve locally
- otherwise forward to the peer that hosts it

That desired behavior is much more specific than generic "high availability."
It is about request-time correctness, not merely component plurality.

## What a real orchestration promotion packet would need

Before this research could justify promoting a whole candidate family, the
packet would need to show:

- the specific missing truth that the current runtime still cannot carry
- the exact candidate layer that now owns that truth better
- what wrong-node, drift, or failover behavior improved measurably
- what operator reconstruction burden still remains afterward
- what new control-plane, storage, or workflow burden was accepted

Without that packet, this page should stay a decision filter, not a disguised
verdict.

## Bottom line

The repo does not need more orchestration names.
It needs one more option that remains believable after wrong-node entry,
backend loss, and hidden operator glue become real.

If an orchestration layer wins, it will not win by prestige.
It will win because one more previously private bad-day sentence stopped being
true.
