# Infrastructure Master Plan: What It Actually Proves

This page is not a replacement for
[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md).
It exists because the master plan is one of the easiest files in the repo to
overread.

The master plan is detailed enough to sound half-implemented.
That is dangerous.

The danger is very specific:

- a module is named
- a flow is described
- the description sounds coherent
- coherence gets mistaken for system-owned truth

That is exactly the drift this knowledgebase is trying to stop.

This page therefore should not be read as a friendly roadmap summary.
It is a corrective retrieval layer.

Its job is not to make the master plan feel shorter.
Its job is to stop the reader from remembering the wrong thing about the master
plan after skimming it.

## Strongest honest current answer

If a reader asks what the master plan actually buys today, the shortest
defensible answer is:

> the master plan proves the repo understands that the real missing layer is
> distributed truth ownership rather than lack of components, and it sketches a
> thinner control layer intended to externalize placement, convergence,
> failover, and portability burdens without immediately surrendering to a full
> heavyweight orchestrator. It does not prove that those truth surfaces are
> already live, trusted, or ownership-complete in the tracked runtime.

Anything calmer or stronger than that is probably overreading the plan.

## What this page is and is not allowed to prove

This page is allowed to prove:

1. the master plan is one of the strongest directional artifacts in the repo
2. the plan names the missing truth layers with unusual clarity
3. the repo is deliberately searching for a thinner control layer before
   heavyweight promotion
4. the plan understands the dominant wound as hidden distributed truth, not
   lack of infra nouns
5. the plan treats request preservation, portability, and anti-private-ritual
   operation as central

This page is not allowed to prove:

- that the planned layers are already runtime-owned
- that named modules are already active
- that a coherent sequence equals removed operator burden
- that the thinner layer will certainly be sufficient forever

It is also not allowed to let coherence impersonate ownership.

## Retrieval contract for this page

### Class 1: master-plan directional evidence

Primary anchor:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This class is allowed to prove:

- what the repo now understands clearly
- what control surfaces it wants to add
- what kinds of burden it is trying to move

It is not allowed to prove:

- that the live runtime already owns those truths

### Class 2: repo-direction context

Primary anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

This class is allowed to prove:

- the repo wants a thinner middle layer before heavyweight capture

It is not allowed to prove:

- that heavyweight promotion is permanently ruled out

### Class 3: archive-pressure context

Primary anchors:

- archive-pressure synthesis pages
- archive conversations about wrong-node pain and fake options

This class is allowed to prove:

- why the master plan had to become more explicit
- what accusation the plan is trying to answer

It is not allowed to prove:

- that naming the accusation already removes it

### Class 4: live runtime reality

Primary anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active Compose fragments

This class is allowed to prove:

- what the tracked runtime actually owns today

It is not allowed to be silently replaced by Class 1 coherence.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "the repo understands the wound clearly" | master-plan directional evidence | `docs/INFRASTRUCTURE_MASTER_PLAN.md` | that the live runtime already owns the planned truth surfaces |
| "the repo wants a thinner middle layer before heavyweight promotion" | repo-direction plus master-plan evidence | `.github/copilot-instructions.md`, `README.md`, master plan | that heavyweight promotion is permanently ruled out |
| "the hidden burden is distributed truth, not lack of tools" | master-plan plus archive pressure | master plan plus archive-pressure pages | that naming the burden already removes it |
| "specific modules are promising" | master-plan directional evidence | `docs/INFRASTRUCTURE_MASTER_PLAN.md` | that those modules are active, validated, or ownership-complete |

If a paragraph cannot be mapped cleanly here, it is probably flattening the
plan into a stronger present tense than the repo has earned.

## The fastest correct reading of the master plan

This page should let the reader recover the useful summary immediately:

- the plan proves the repo understands the wound clearly
- the wound is hidden distributed truth, not lack of infra nouns
- the plan is strong evidence for direction
- the plan is not runtime proof
- the plan only matters if it eventually removes hidden operator
  reconstruction burden rather than merely diagramming it better

If that summary is not obvious, the plan is still too easy to overread.

## What still does not count as master-plan progress

The following still do not count as meaningful completion signals:

- a planned module being named convincingly
- a sequence sounding operationally plausible
- a sync or failover flow being easy to picture
- a future control surface being cleaner than the current one
- a roadmap making the architecture feel emotionally calmer

Those are signs that the repo is thinking harder.
They are not signs that the hidden operator burden has actually moved.

This page should keep that asymmetry visible.

Better planning is real progress.
It is just not the same category of progress as one more runtime truth leaving
the operator's memory.

## What the master plan clearly proves

### 1. The repo is deliberately searching for a thinner control layer before heavyweight promotion

The master plan repeatedly returns to the same structure:

- keep Docker and Compose as the human authoring layer
- add targeted sync, convergence, and failover logic where Compose alone is too
  static
- pay only for the coordination semantics that remove real operator pain

That is not procrastination.
It is architecture intent under pressure.

The repo is not trying to preserve Compose out of nostalgia.
It is trying to preserve the last directly inspectable layer until something
else proves it can own the painful truths more honestly.

What this proves:

- the repo is explicitly searching for a missing middle layer
- it does not accept "just use Kubernetes" as a sufficient answer
- orchestration promotion is supposed to be earned by demonstrated need, not by
  fashion

What this does not prove:

- that the thinner layer will be sufficient forever
- that the thinner layer is already implemented
- that the thinner layer will remain thin once it owns enough truth

### 2. The repo sees cross-node truth, not extra containers, as the dominant missing capability

This is one of the strongest signals in the master plan.

The plan is not mainly:

- add more components
- add more sidecars
- add more proxy tricks

It is mainly:

- stop leaving cluster truth in operator memory

The recurring missing duties are:

- env and secret convergence
- compose or file convergence
- service placement awareness
- restart or redeploy semantics when truth changes
- failover logic that follows current reality instead of stale assumptions

What this proves:

- the user's dominant pain is hidden coordination truth
- the project already understands that human memory is currently one of the
  biggest SPOFs

What this does not prove:

- that a sync-agent style layer is already trustworthy enough to own those
  responsibilities

### 3. The plan treats request preservation as central, not decorative

The master plan spends real energy on:

- node-entry behavior
- global versus node-specific hostnames
- DDNS and route generation
- internal peer connectivity
- recovery and repair assumptions

That matters because the user is not mainly asking:

> how do I deploy services?

They are asking:

> how do I keep the request alive, honest, and semantically coherent when it
> lands on the wrong node or when the expected backend disappears?

What this proves:

- any future control layer that ignores wrong-node request preservation is
  misaligned
- routing truth is part of the core problem, not a polish layer

What this does not prove:

- that current failover generation already behaves correctly under failure

### 4. The repo wants a portable personal-cloud system, not a private ritual stack

The master plan keeps pushing toward:

- cleaner bootstrap
- lower env and secret drift
- domain abstraction
- more reusable infra shape

That reveals another important part of the dream:

the system should not only work for the original operator who already remembers
every invisible fact.

What this proves:

- documentation and control surfaces must survive portability pressure
- "works for the author because the author remembers the ghosts" is not an
  acceptable end state

What this does not prove:

- that portability or reproducibility are already solved

## The burden-transfer scorecard for the planned modules

The master plan is easiest to overread when modules are treated as inherently
good because they are specific and coherent.

They should instead be judged by what hidden burden they are trying to remove.

| Planned module or theme | Hidden burden it is trying to remove | What still remains unproven even if the plan is coherent |
| --- | --- | --- |
| env and secret sync | private operator memory about cross-node parity | whether forwarded requests actually land in semantically equivalent runtime |
| compose and file sync | silent drift between nodes | whether the resulting nodes still expose the same useful recovery behavior under failure |
| placement truth such as `services.yaml` | remembered service location | whether routing or eligibility logic actually consumes live placement truth |
| failover-agent or routing repair logic | routes that only look dynamic while the primary is healthy | whether recovery paths survive the exact backend-loss event that made them necessary |
| bootstrap portability | one operator remaining the only person who knows the ritual | whether portability also preserves readability rather than hiding truth behind automation |
| repair and restart semantics | private judgment about what should restart or stay stable | whether the automation can make those decisions without quietly becoming a larger controller |

That table matters because the user is not mainly asking for more architecture.
They are asking for less private burden at request time and failure time.

## What the master plan does not prove

The master plan is strong planning evidence.
It does **not** prove that any of the following are already true in the tracked
runtime:

- multi-node env and secret sync works live
- compose and file convergence is already trustworthy
- a live tracked root `services.yaml` exists and is consumed
- peer-aware forwarding is broadly validated
- failover routes survive the exact failure that made them necessary
- sync-agent or failover-agent loops are implemented and battle-tested
- watchtower-style repair semantics are dependable
- stateful services already have anti-SPOF correctness

This list should stay painful.
If it becomes easy to forget, the plan will start impersonating runtime truth.

## The main overreading risk this page is trying to stop

The master plan is detailed enough to simulate settlement.

That means it can create a false feeling that:

- the missing layer is already mostly designed
- therefore the missing layer is already mostly known
- therefore the remaining gap is mostly implementation effort

That sequence is dangerous here.

The user's real frustration is not that no one ever drew the missing layer.
It is that many drawn layers still do not prove they would move enough truth
out of human memory to become a real option.

So the hard question after every planned module remains:

- if this module existed tomorrow, what exact hidden burden would still remain
  socially reconstructed?

If that question is not asked, the plan becomes too flattering.

## What a real promotion packet for the master plan would need

Before the master plan can be read as more than strong directional evidence, a
promotion packet would need to show:

- which planned truth surface is now live
- which exact operator reconstruction burden it removed
- which current root runtime artifact consumes that truth
- which failure case is now handled without private human memory
- what new burden the added layer introduced in exchange

If that packet cannot be assembled, the correct reading stays:

- the plan names the wound well
- the plan may point toward the repair
- the runtime still has to earn ownership

## Bottom line

The master plan matters because it names the wound precisely enough that the
repo can stop pretending it merely needs more tools.

But naming the wound is not the same thing as removing it.
Until the tracked runtime can show that one more important truth has actually
left private operator memory and become inspectable system-owned behavior, the
master plan remains strong directional evidence rather than runtime proof.
