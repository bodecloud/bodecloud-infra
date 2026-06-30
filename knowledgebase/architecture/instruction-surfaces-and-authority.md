# Instruction Surfaces and Authority

This page exists because one of the easiest ways to make this repo lie is to
let several truthful files blur into one false voice.

That blur usually sounds disciplined and reasonable:

- the repo has a clear direction
- the main instruction files agree
- the architecture is well understood
- therefore the remaining work is mostly implementation

That sequence is exactly the kind of polished overpayment the user is tired of.

The repo is not harmed mainly by lack of language.
It is harmed when language starts borrowing authority from neighboring files
that are talking about different truth classes:

- dream truth
- honesty-boundary truth
- runtime-anchor truth
- authoring-discipline truth

If those become interchangeable, the docs become smoother and less honest at
the same time.

That is not a minor editorial risk in this repo.
It is one of the main ways the user's dream gets quietly normalized into
something weaker and more familiar.

## The exact question this page answers

This page answers the recurring question:

> which repo files actually explain the multiple ordinary Docker nodes,
> no-Swarm-by-default, local-first-then-peer-forward, wrong-node-survival
> dream, and which files only support that dream indirectly?

That question matters because a contributor can read the right files and still
misread the repo if they do not know what each file is allowed to prove.

The danger is not only factual error.
The danger is a quieter drift:

- workflow guidance gets narrated as architecture law
- architecture intent gets narrated as runtime proof
- service-authoring discipline gets narrated as distributed-systems maturity
- repeated repo language gets narrated as if the missing middle layer were
  already settled

This page exists to stop that drift.

It also exists because the user keeps asking variants of the same question as
a stress test:

> does the repo itself still know where its truth lives, or has it started
> borrowing confidence from adjacent files again?

## What this page is and is not allowed to prove

This page is authoritative about:

- which instruction surfaces carry which class of truth
- how to route architectural questions to the right source
- how to keep repo agreement from being overread as runtime closure
- which instruction surfaces are central and which are merely supportive

This page is not authoritative about:

- whether the runtime already satisfies the dream it describes
- whether the future truth-owning middle layer is already implemented
- whether repeated language across files should be treated as operational proof
- whether the current stack has already earned stronger failover language

This is an authority map.
It is not a completion report.

## Strongest honest current answer

The shortest honest answer is:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  is the clearest architecture-dream file in the repo
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
  is the strongest repo-level honesty wall around that dream
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
  is the runtime/operator anchor that keeps claims tied to the priority root
  implementation
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)
  is mostly Compose-authoring discipline and should not be overpromoted into
  distributed-systems authority

Those files are aligned.
They are not equal.

If a contributor comes away thinking "they all basically say the same thing,"
then this page has failed.

Because they do not.
They are aligned around the same wound, but each one carries a different kind
of custody:

- one names the dream
- one keeps the honesty wall around the dream
- one drags claims back onto the real runtime surface
- one prevents bad Compose habits from making the truth even harder to assess

## What still does not count as understanding authority here

This page also needs to block a softer misread.

The following still do not count as having understood the repo's instruction
surfaces correctly:

- noticing that the main files use similar language
- concluding they all carry the same truth class
- using repeated wording as if it were runtime corroboration
- treating stricter style rules as distributed-systems maturity evidence
- treating implementation-anchor files as if they were the architecture thesis

Those are exactly the moves that let the docs become smoother and less honest
at the same time.

They are also the moves that make the repo sound more mature than the current
runtime has actually earned.

## What a real authority-mapping packet would have to contain

Before this page supports a stronger claim like "the repo's guidance surfaces
are now genuinely assimilated," it should be because the documentation leaves
behind a concrete packet.

That packet should contain:

- the exact file being cited
- the truth class that file carries
- the stronger claim that file is allowed to support
- the adjacent stronger claim it is not allowed to support
- the file that must be consulted next if the question shifts truth class

Without that packet, the reader can still mistake aligned files for one merged
voice.

## The ranking that should govern the whole knowledgebase

For the specific question:

> where does the repo most clearly explain the Compose-first, multi-node,
> anti-SPOF, wrong-node-preserving idea without immediately defaulting to Swarm
> or Kubernetes?

the files should be read in this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
4. [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

That order is not stylistic.
It is a truth-custody order:

- dream first
- repo-level honesty wall second
- runtime anchor third
- service-authoring discipline fourth

If the order flips, the repo usually starts smuggling runtime maturity into
files that were only supposed to protect framing or authoring discipline.

## The simplest default answer

If the user asks some version of:

> do `AGENTS.md`, `.github/copilot-instructions.md`, and `.cursorrules`
> explain the multiple-node Docker failover/fallback thing?

the correct short answer is:

- `copilot-instructions.md`: yes, directly
- `AGENTS.md`: only partially, mostly by anchoring the real runtime surface
- `.cursorrules`: mostly no, except where stricter Compose hygiene prevents
  fake resilience stories from getting even easier to write

That answer should be treated as the default compression of the authority map.

It should also be treated as a warning against overpromotion.
The repo now has enough repeated language that a casual reader can mistake
coherence for completion.

## File-by-file authority

## `.github/copilot-instructions.md`

This is the repo's clearest architecture manifesto.

It is the strongest source for claims like:

- this is a Compose-first multi-node Docker repo
- the project explicitly resists immediate capture by Swarm or Kubernetes
- there is no central orchestrator by default
- local-first serving is the desired normal path
- peer-forward fallback is part of the target operating contract
- L7 HTTP and L4/TCP must not be flattened together
- Cloudflare plurality is only the first hop, not the whole failover story
- a current-state registry concept such as `services.yaml` is central to the
  missing layer

It also contains the clearest behavior sketch in the whole repo:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That matters because it proves the dream is not vague.
It is request-contract-shaped.

### What it can legitimately prove

- the dream is explicit, not inferred
- wrong-node survival is part of the target
- the repo wants current-state truth, not only desired-state declarations
- the project is intentionally trying to avoid heavyweight control-plane tax by
  default

### What it cannot legitimately prove

- that the tracked root runtime already consumes a live root `services.yaml`
- that generic wrong-node success already exists
- that fallback survives backend loss in the current stack
- that middleware or auth continuity are already proven
- that TCP or stateful failover are solved

This file is high-authority intent.
It is not live runtime proof.

It is the strongest place to learn what the repo is trying to become.
It is also the easiest place to accidentally overread if the reader is hungry
for closure.

## `README.md`

The root README is the strongest repo-facing synthesis and honesty wall.

Its authority is different from `copilot-instructions.md`.
It is not the sharpest architecture manifesto, but it is the strongest
repo-level compression of:

- what the dream is
- what the root runtime really contains
- why the current option space still feels thinner than it looks
- what the repo must not pretend has already been solved

It is the best file for claims like:

- the repo already has a serious Compose-first platform surface
- the dream is real and central
- the hidden operator burden still matters
- the missing truth-owning middle layer is still incomplete

### What it can legitimately prove

- the repo values honesty boundaries enough to put them at root level
- the dream and the wound are both explicit
- the current stack is substantial without being overclosed

### What it cannot legitimately prove

- runtime correctness by itself
- that every caveat it names has already been repaired
- that the future control surface is settled

The README should be treated as the strongest public "do not lie to yourself"
surface in the repo.

That makes it unusually important.
Most READMEs try to sell competence.
This one also has to prevent competence theater.

## `AGENTS.md`

`AGENTS.md` matters because it drags the conversation back to the real
worktree.

It is the strongest source for claims like:

- the priority implementation is the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  plus included `compose/` fragments
- this is an infrastructure repo, not a normal application repo
- validation is constrained by real env and secret burdens
- implementation claims should tie back to actual commands and actual runtime
  surfaces

It is not the clearest explanation of the dream.
Its job is different.
It prevents abstraction drift by forcing contributors back onto the real root
runtime.

### What it can legitimately prove

- where implementation truth should be checked first
- which files are the priority authoring/runtime surfaces
- which operational constraints shape validation claims

### What it cannot legitimately prove

- the full wrong-node, peer-forward, anti-SPOF dream by itself
- that any HA path already works end to end
- that the runtime has already earned the stronger language used in the dream
  files

`AGENTS.md` is the implementation anchor.
It is not the architecture thesis.

That distinction matters because implementation-anchor files often acquire
accidental charisma.
They mention real commands, real paths, and real runtime surfaces, which can
make them sound like they implicitly prove more architecture than they do.

## `.cursorrules`

`.cursorrules` is narrower still.

This is the file most likely to be accidentally overpromoted by someone who
sees strong language and assumes they are reading distributed-architecture
law.
They usually are not.

Its main authority is around service-authoring discipline:

- prefer inline configs over scattered external config files
- require real healthchecks
- do not "fix" instability by removing healthchecks
- keep the repo from getting even easier to fake-resilience into

That matters.
In a repo trying to avoid fake HA, poor healthcheck discipline and config
sprawl absolutely make the story worse.

But `.cursorrules` still does not define the routing and failover model.
It mainly protects the repo from authoring patterns that would make honest
runtime assessment harder.

### What it can legitimately prove

- the repo takes Compose hygiene seriously
- the repo resists anti-debugging habits like disabling healthchecks
- portable, reviewable service definitions are part of its discipline

### What it cannot legitimately prove

- cross-node request preservation
- peer-selection logic
- wrong-node success
- stateful failover correctness
- the repo's full architectural dream

This file is supportive discipline, not the main truth surface.

It is valuable precisely because it is narrower.
When it is overpromoted, the repo starts mistaking strict Compose hygiene for a
distributed answer.

## Question-to-source map

When the question is:

- what is the platform trying to become at request time?
- is wrong-node preservation really part of the design target?
- is the repo intentionally resisting immediate Swarm or Kubernetes capture?

start with
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).

When the question is:

- how boldly can that dream be described without cheating?
- what honesty wall is the repo already trying to keep?
- which future-facing claims are still only intent rather than proof?

start with
[`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md).

When the question is:

- where should live implementation truth be checked first?
- which files and commands define the priority runtime?
- what real env and secret burdens constrain validation claims?

start with
[`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md).

When the question is:

- what authoring behaviors are mandatory for Compose services?
- what counts as sabotage in healthcheck handling?
- how does the repo prevent fake resilience from becoming even easier to write?

start with
[`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules).

That map matters because bad documentation rarely lies in one sentence.
It usually lies by answering a good question from the wrong authority surface.

That is the exact failure mode this page is trying to make harder:

- a reader asks an architecture question
- the repo answers from an implementation or discipline surface
- the answer sounds concrete enough to feel satisfying
- the missing truth layer stays missing anyway

## Illegal promotions

These are the promotions this page should make harder:

- "these files all agree, so the runtime gap is probably small now"
- "the repo has a clear dream, so the missing middle layer is mostly settled"
- "strict authoring rules imply distributed correctness"
- "runtime validation guidance implies failover maturity"
- "a cleaner authority hierarchy means the future control plane is nearly
  chosen"

Those are all versions of the same flattening error:

- dream truth
- honesty truth
- runtime truth
- authoring-discipline truth

get collapsed into one calm story.

The user keeps pressing this question because they are not merely ranking docs.
They are checking whether the repo itself knows where its truth lives.

That is why authority has to be handled as custody, not just convenience.
Each file is protecting the repo from a different kind of self-deception.

## Bottom line

The most important thing this page proves is not that the repo now has a tidy
documentation hierarchy.

It proves something harsher:

> the repo names the dream, the honesty wall, the runtime anchor, and the
> authoring discipline more clearly than before, but those are still different
> truth classes, and treating their agreement as runtime proof would recreate
> exactly the same polished ambiguity the user keeps rejecting.

That is why the order has to stay explicit:

- dream first
- honesty wall second
- runtime anchor third
- authoring discipline fourth

If that order fades, the repo starts remembering itself wrong.

And "remembering itself wrong" is serious here.
This project is explicitly trying to build systems that no longer require the
operator's private memory to finish the platform's story.
If the docs themselves start leaning on blurred memory and borrowed authority,
they reenact the same failure at the documentation layer.
