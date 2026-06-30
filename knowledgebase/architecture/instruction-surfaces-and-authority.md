# Instruction Surfaces and Authority

This page exists because one of the easiest ways to ruin this repo is to let
different instruction files borrow authority they do not actually have.

That mistake shows up fast:

- workflow guidance gets narrated as architecture law
- architecture intent gets narrated as runtime proof
- authoring discipline gets narrated as distributed correctness

The result is always the same:

the repo starts sounding more solved than it is.

That is exactly what the user does not want.
It is also one of the easiest ways for the docs to recreate the same
option-space fraud the user keeps reacting against:

- several files sound aligned
- alignment gets mistaken for one coherent implemented truth
- one coherent truth gets mistaken for runtime behavior
- the actual authority boundary disappears

They are already frustrated by infrastructure ecosystems that pretend the
option space is richer than it really is.
This docs set cannot repeat that behavior by flattening its own instruction
surfaces into one fake voice.

## The question this page answers

The concrete recurring question is:

> which repo instruction files actually define the multi-node Docker, no-Swarm,
> wrong-node-survival, failover/fallback idea, and what does each file have the
> authority to claim?

That sounds narrow.
It is not narrow.

If the authority map is wrong, every other page drifts.

## Priority source stack for the exact question the user keeps asking

For the narrow but crucial question:

> where does the repo most clearly explain the multiple ordinary Docker nodes,
> no-Swarm-by-default, failover/fallback, wrong-node-survival direction?

the source stack should be read in this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
4. [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

That order is not cosmetic.

It means:

- the dream is named first
- the repo-level honesty wall comes next
- the runtime anchor comes after that
- the authoring discipline comes last

If a contributor flips that order, they will usually start smuggling runtime
authority into files that were only meant to protect writing discipline or repo
framing.

## What this page is and is not allowed to prove

This page is authoritative about:

- which instruction surfaces carry which class of truth
- how to route architecture, runtime, and authoring questions to the right file
- why aligned repo language must not be mistaken for implemented distributed truth

This page is not authoritative about:

- whether the runtime already satisfies the dream those files describe
- whether the missing middle layer is already implemented or settled
- whether repeated language across files should be narrated as proof

This page is an authority map.
It is not a runtime verdict.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "the dream is clearly named in-repo" | repo-native intent | `.github/copilot-instructions.md`, `README.md`, this page | that the runtime already behaves that way |
| "the current runtime anchor lives here" | live-repo validation surface | `AGENTS.md`, root `docker-compose.yml`, compose fragments | that validation guidance equals architecture closure |
| "this authoring rule supports honesty" | authoring-discipline reasoning | `.cursorrules`, Compose files, this page | that authoring discipline proves cross-node behavior |
| "these files are aligned but unequal" | authority synthesis | this page plus the named files | that aligned language is therefore implementation proof |

If an authority sentence starts sounding like a system-completion sentence, it
has crossed out of this page's jurisdiction.

## The shortest honest answer

If someone asks:

> do `AGENTS.md`, `.github/copilot-instructions.md`, and `.cursorrules`
> explain the multiple-Docker-node, no-Swarm, failover/fallback model?

The correct answer is:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):
  yes, directly, and more clearly than any other repo file
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md):
  only partially, mostly by anchoring attention to the actual root runtime
  surface
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules):
  mostly no, except where Compose hygiene helps prevent fake resilience claims

The blunt version is:

- if you want the dream, read `copilot-instructions.md`
- if you want the runtime anchor, read `AGENTS.md`
- if you want authoring discipline, read `.cursorrules`

If a contributor leaves this page still thinking all three files "basically say
the same thing," then this page has failed.

That is the shortest honest answer to the exact question the user keeps
pressing.
It is also a useful test for the rest of the knowledgebase:
if another page cannot preserve that unevenness, it is probably smoothing the
repo into a more coherent story than the evidence supports.

That answer also needs an operational form.
Readers should not only remember which file is strongest.
They should know which file to reach for when a specific class of question
appears.

## Question-to-source map

When the question is:

- what is the platform trying to become at request time?
- is wrong-node preservation part of the design target?
- is the repo intentionally resisting immediate Swarm or Kubernetes capture?

start with
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).

When the question is:

- how boldly can the dream be described without cheating?
- what honesty wall does the repo itself keep trying to preserve?
- which future-facing claims are acknowledged as intent rather than proof?

start with
[`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md).

When the question is:

- where should live implementation truth be checked first?
- what files and commands define the runtime anchor?
- what real env and secret burden constrains validation claims?

start with
[`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md).

When the question is:

- what authoring patterns are mandatory for Compose services?
- what counts as sabotage in healthcheck or config handling?
- how does the repo try to avoid making fake resilience even easier to write?

start with
[`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules).

This matters because a lot of bad documentation does not lie in one sentence.
It lies by answering a good question from the wrong authority surface.

## The false-convergence problem this page has to stop

This repo now has several instruction surfaces that are much better aligned
than they used to be.

That improvement is real.
It is also dangerous.

Alignment across instruction files can too easily be overread as:

- the repo knows what it wants
- therefore the missing layer is basically designed
- therefore the runtime gap is mostly implementation work

That sequence is exactly the kind of false convergence this page exists to
interrupt.

The right conclusion is narrower:

- the repo's dream is easier to reconstruct
- the authority boundaries are clearer
- the runtime still has not earned the stronger behavior claims those files are
  aiming at

If that last line disappears, this page starts helping the docs become calm,
ordered, and wrong again.

That is the exact failure mode to watch for:

- intent gets used to close proof gaps
- workflow guidance gets used to imply architecture maturity
- authoring discipline gets used to imply distributed correctness
- aligned language gets used to imply that the missing middle is already
  mostly chosen

The repo is more coherent than before.
That does not mean it is more finished than before.

It is important because people naturally want to say:

- "all three files point in the same direction"

That statement is true only in a weak sense.
It is false if used to imply that all three files *explain* the no-Swarm,
wrong-node, peer-forward architecture equally well.

They do not matter equally.
Treating them as equal is one of the easiest ways to misunderstand the repo.

## What is really at stake

The user's real complaint is not:

- "please rank repo docs"

It is:

> please stop giving me systems and explanations that sound rich and flexible
> until I try to pin down where the actual truth lives

That complaint applies inside this repo too.

Every important page in the knowledgebase eventually depends on answering:

- which file names the dream most clearly?
- which file keeps us tied to the real runtime?
- which file enforces discipline without overclaiming architecture?
- which file mostly governs Compose hygiene rather than distributed truth?

That is the real purpose of this page.
The point is not to be tidy about documentation hierarchy.
The point is to stop the repo from remembering itself wrong.

## The authority ranking the docs should use

For the multi-node Docker without immediate orchestrator surrender story, the
highest-signal files are:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

They should be read in that order unless the question is specifically about
runtime validation or authoring rules.

The simplest summary is:

> `copilot-instructions.md` says the dream most clearly; `README.md` explains
> the dream and its honesty walls at repo level; `AGENTS.md` ties claims back to
> the actual runtime surface; `.cursorrules` is mostly service-authoring
> discipline, not architecture explanation

That summary should be treated as the default answer when the user asks some
version of:

- "which file actually explains what we are trying to do?"
- "which file is just supporting discipline?"
- "which file proves less than it sounds like it does?"

That ranking should drive every future documentation judgment.
It should also be read as a defense against a specific temptation:

- treat repeated repo language as proof
- treat repo-wide agreement as implementation
- treat service-authoring rules as distributed-systems semantics

It should also be read as a defense against one subtler temptation:

- treat a clean authority hierarchy as proof that the future control surface is
  mostly settled

The user keeps asking this because they are not really ranking docs.
They are checking whether the repo itself knows where its truth lives.

That is why the order has to stay explicit:

- dream first
- honesty wall second
- runtime anchor third
- authoring discipline fourth

## Trust matrix

| File | Best used for | Strongest claims it can support | Claims it cannot support alone |
| --- | --- | --- | --- |
| `.github/copilot-instructions.md` | Architecture intent | No central orchestrator by default, multi-node Docker worldview, local-first then peer-forward, shared placement-truth concept, separate L7/L4 thinking | Live end-to-end failover proof |
| `README.md` | Repo-facing synthesis and honesty boundary | The dream is real, the root runtime is substantial, and the proof gaps still matter | Runtime correctness by itself |
| `AGENTS.md` | Runtime anchor | Root `docker-compose.yml` plus `compose/` are the priority implementation surface and validation must stay tied to real env/secret requirements | The full HA model by itself |
| `.cursorrules` | Compose hygiene and authoring discipline | Real healthchecks, inline configs, portable reviewable service definitions, self-healing posture | Cross-node request preservation, peer-forwarding logic, or failover correctness |

This is not just a neat table.
It is the boundary between honest docs and docs that start sounding like sales
material.

## Strongest honest current answer

If someone asks, "What is the most important thing this page proves?" the
shortest defensible answer is:

> The repo now names the dream, the honesty walls, the runtime anchor, and the
> authoring discipline much more clearly than before, but those are still
> different truth classes, and treating their agreement as runtime proof would
> recreate exactly the same polished ambiguity the user keeps rejecting.

It is also a defense against a specific flattening error:

- architecture intent
- runtime anchor
- authoring discipline

are not interchangeable categories.

## File-by-file authority

## `.github/copilot-instructions.md`

This is the repo's clearest architecture manifesto.

If someone asks the user's question in the most direct form:

> which file actually explains the multiple Docker nodes, no-Swarm,
> local-first-then-peer-forward idea?

this is the answer.

It says, directly:

- this is a multi-node Docker infrastructure repo
- the project wants to become anti-SPOF without defaulting to Swarm or
  Kubernetes
- there is no central orchestrator by default
- services are manually assigned to nodes
- current-state truth matters more than scheduler-declared desire
- local-first service is preferred
- peer-forward fallback is part of the desired baseline
- L7 HTTP and L4/TCP are different problem classes
- Cloudflare multi-record DNS is part of node-entry resilience, not the whole
  story

It also contains the cleanest request-contract sketch in the repo:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That sketch matters because it proves the dream is not vague.
It is behavior-shaped.

This file does not merely say:

- there are several nodes
- Traefik exists
- Cloudflare exists

It says the repo wants a very particular kind of system:

- first hop can be any healthy node
- locality should stay visible instead of being hidden behind cluster ritual
- wrong-node entry should not be fatal
- peer-forwarding is normal desired behavior, not an exotic future luxury
This matters because it proves the repo's dream is behavior-shaped enough to be
judged, not just broad enough to be admired.

### What it can legitimately prove

- the dream is explicit, not inferred
- the project is intentionally resisting premature heavy-orchestrator capture
- wrong-node survival is part of the design target
- a shared placement-truth surface such as `services.yaml` is central to the
  concept

### What it cannot legitimately prove

- that the root runtime currently consumes a live tracked `services.yaml`
- that peer-forwarding works end to end today
- that fallback survives backend-loss events
- that middleware and auth continuity are already proven
- that L4/TCP failover is solved
- that stateful systems are honestly HA

This file is high-authority intent.
It is not runtime proof.

## `README.md`

The root README is the strongest public synthesis surface.

Its authority is different from `copilot-instructions.md`.
It is less raw as architecture law, but stronger as the repo's main honesty
boundary for:

- what the project is trying to make true
- what the root runtime already contains
- what is still planning, aspiration, or partial proof
- why the documentation must keep those states separate

### What it can legitimately prove

- the dream matters enough to sit at repo root
- the repo knows it is navigating a missing-middle-layer problem
- overclaiming is already recognized as dangerous

### What it cannot legitimately prove

- runtime behavior by itself
- that every caveat it mentions has been operationally closed
- that the future control surface is settled

The README is the strongest repo-facing "do not lie to yourself" file.
The repo already has enough parts to sound ambitious.
The harder job is preserving a truthful compression of what those parts do not
yet buy.

It does **not** say that the repo has already resolved which truth-owning
layer will ultimately earn promotion.
A README can be extremely clear about the wound while the cure is still
genuinely unsettled.

## `AGENTS.md`

`AGENTS.md` matters because it drags the conversation back to the real worktree.

It is also important to say what it does not do:

it does not explain the dream as clearly as `copilot-instructions.md`.

Its job is not to give the whole philosophy of the project.
Its job is to stop abstraction drift.

It makes clear:

- this is a Compose infrastructure repo, not a normal app repo
- the priority implementation is the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  plus the fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- validation is constrained by real env and secret requirements
- claims about the system should tie back to real commands and real runtime
  surfaces

What it does **not** do is explain the multi-node failover/fallback dream with
the same specificity as `copilot-instructions.md`.

It says enough to confirm that this is a multi-node Compose repo and that the
root runtime matters.
It does not carry the full philosophical burden of the no-Swarm, any-node,
wrong-node-survival model.

That matters because one of the user's biggest complaints is fake option-space.
`AGENTS.md` is one of the files that helps resist that by keeping the repo
grounded in actual inspectable surfaces.

### What it can legitimately prove

- where implementation truth should be checked first
- that the repo is still Compose-first in operator reality
- that environment and secret requirements are part of the real runtime burden

### What it cannot legitimately prove

- the whole no-cluster peer-forward model by itself
- that the dream is already implemented
- that any HA path works end to end

`AGENTS.md` is an implementation anchor.
It is not the architecture thesis.
That distinction has to stay sharp because implementation anchors are
structurally attractive to overread.

## `.cursorrules`

`.cursorrules` is narrower still.

This is the file most likely to be accidentally overpromoted by people who see
strong wording and assume they are looking at architecture.
They are mostly not.

Its job is to enforce service-authoring behavior that prevents the
infrastructure story from becoming even more brittle than it already is.

It mostly says:

- prefer inline configs over external file sprawl
- require real healthchecks
- do not "fix" services by deleting healthchecks
- keep a self-healing posture where appropriate

That is useful.
It is also much narrower than people sometimes assume when they lump it together
with the architecture instruction files.

In fact, one of the strongest reasons not to overread `.cursorrules` is that it
also contains generic material that is not really about this repo's distributed
architecture question at all.
Its presence is a reminder that not every instruction surface is pure signal.

That is not trivial.
In a repo trying to escape fake HA, healthcheck discipline and portable config
surfaces matter a lot.

But this file still does not define the overall routing and failover model.

### What it can legitimately prove

- the repo wants health-aware, reviewable, portable service definitions
- self-healing posture is intentional
- resilient authoring discipline matters to the project

### What it cannot legitimately prove

- cross-node failover correctness
- placement truth
- peer-forward decision logic
- distributed request preservation

Treat it as authoring discipline, not control-plane proof.
That sentence should remain blunt because `.cursorrules` is exactly the kind of
file people naturally overpromote once they see strong wording about health and
Compose rigor.

## How these files combine into one coherent reading

If you read them in the right order, the repo becomes much less ambiguous.

That sentence needs an explicit warning attached to it:

- less ambiguous is not the same thing as more implemented
- coherent reading order is not the same thing as settled control-plane choice
- conceptual agreement is not the same thing as runtime-owned distributed truth

This page should help the reader recover a cleaner map of the repo's internal
argument.
It should not make that argument sound more finished than it is.

## 1. `copilot-instructions.md` names the dream

It defines the desired operating contract:

- no central orchestrator by default
- multi-node Docker
- current-state truth
- local-first service
- peer-forward fallback
- separate L7 and L4 handling

This file says what the project wants the platform to feel like at request
time.

## 2. `README.md` says the dream is real but not fully proven

It lifts the architecture intent into repo-level language while preserving the
honesty wall between:

- what is live
- what is planned
- what is still missing

This file says how boldly the dream may be described without cheating.

## 3. `AGENTS.md` says where to look before claiming anything

It forces the reader back to:

- root Compose
- included fragments
- real validation constraints

This file says where current truth must be checked.

## 4. `.cursorrules` says how not to sabotage the story

It keeps service definitions from degenerating into:

- fake health
- config sprawl
- review-hostile setups
- brittle local workarounds

This file says how service authoring should behave while the project is still
Compose-first.

## The single most important reading rule for these files together

When these instruction surfaces agree, the right conclusion is usually:

- the repo has a sharper and more internally consistent dream than before

The wrong conclusion is:

- the missing truth layer is therefore mostly known, mostly chosen, or mostly
  operational

That wrong conclusion is what this page should keep making harder.

## What this means for the rest of the knowledgebase

Any page that cites these instruction surfaces should follow these rules:

- use `copilot-instructions.md` to support architecture intent
- use `README.md` to support repo-level framing and honesty boundaries
- use `AGENTS.md` to support runtime-surface prioritization
- use `.cursorrules` to support authoring discipline
- do not let any one of them silently claim runtime proof it does not own

It should then apply one more pass:

- ask whether the chosen file is being used for the claim class it actually
  owns
- ask whether a stronger source exists for that sentence
- ask whether the paragraph quietly upgrades intent into behavior
- ask whether the result makes the repo sound more converged than the current
  worktree proves

That last rule is the whole point.

The user is trying to build a system where the truth is readable.
The docs have to model that same discipline.

## The real takeaway

The instruction surfaces are not just repo trivia.
They are the repo's internal map of where truth lives.
If the docs flatten that map, they stop being a guide and start becoming
another partial-truth surface the operator has to mentally correct.

The decisive takeaway is:

> if you want to know what this project is trying to become, read
> `copilot-instructions.md`; if you want to know how honestly that dream may be
> narrated at repo level, read `README.md`; if you want to know where current
> implementation truth must be checked, read `AGENTS.md`; if you want to know
> how service authoring is supposed to avoid making the whole thing more brittle,
> read `.cursorrules`

That is the authority model the knowledgebase should keep enforcing.

It is also the shortest honest answer to the recurring complaint underneath the
question:

the repo is allowed to have several aligned instruction files, but it is not
allowed to pretend that alignment erases the difference between the file that
names the dream, the file that forces proof back to the worktree, and the file
that mostly keeps YAML hygiene from making the situation worse.
