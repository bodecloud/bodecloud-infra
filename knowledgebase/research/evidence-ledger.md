# Evidence Ledger

This page is the claim governor for the whole knowledgebase.

It exists because the main documentation failure in this repo is no longer
missing information.
It is information being narrated at the wrong confidence level.

The repo has enough:

- real ingress machinery
- real multi-node pressure
- real planning depth
- real archive frustration
- real implementation complexity

to sound mature long before it is trustworthy.

This ledger exists to stop that overpayment.

## What this page is and is not allowed to prove

This page is allowed to:

- define the evidence classes the rest of the knowledgebase should use
- route claim types to the right authority surfaces
- explain confidence ceilings
- keep contradiction visible instead of smoothing it away
- force pages to say what remains unproven even when something important is
  true

This page is not allowed to:

- serve as runtime proof for a specific route or drill
- imply that disciplined narration equals fulfilled implementation
- merge archive pressure, architecture intent, runtime artifacts, and plans
  into one present-tense capability story
- substitute honesty discipline for actual maturity

This is a page about how to speak truthfully, not a page that proves the repo
is done.

## Strongest honest current answer

The biggest remaining documentation risk is not false facts.
It is adjacency masquerading as proof.

That usually happens in this order:

1. archive pressure names the wound clearly
2. instruction files name the dream clearly
3. live Compose files show serious components
4. planning docs propose sharp repairs
5. the docs quietly fuse those into one present-tense feeling of maturity

That fusion is exactly what this ledger exists to block.

## The hidden mistake this ledger is trying to stop

Most weak summaries in this repo are not wrong because they invented facts.
They are wrong because they let several weaker truths emotionally combine into
a stronger claim than any one source actually proves.

That is how documentation starts sounding like this:

- "the stack is basically anti-SPOF now"
- "the system already has a distributed control story"
- "failover is mostly there, it just needs polish"
- "the remaining work is mostly automation"

Those are the kinds of sentences this ledger is supposed to make illegal
unless there is route-specific proof strong enough to survive them.

## The four evidence classes

Every serious claim in this knowledgebase should route through one of four
classes.

### Class 1: live runtime evidence

Use for claims about what the current tracked implementation actually contains
or currently exposes.

Typical questions:

- what networks are defined in the merged root runtime?
- is Traefik part of the active edge?
- does Headscale currently use SQLite?
- is a given service in the priority implementation?

Strong anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active files under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)
- `docker compose config`
- route-specific drill output, if a drill really exists

What Class 1 can honestly prove:

- authored or inspected runtime shape
- presence of a component
- local implementation facts
- the starting condition for a stronger drill

What Class 1 cannot prove by itself:

- the user's full dream
- distributed behavior the runtime only gestures toward
- future failover semantics from present components alone
- backend-loss recovery
- stateful authority correctness

Presence is real.
Presence is not closure.

### Class 2: architecture-intent and honesty evidence

Use for claims about what the repo explicitly wants to become and what honesty
boundaries it already enforces.

Typical questions:

- what request contract is the repo trying to earn?
- why is the repo resisting immediate heavyweight orchestration?
- where does the repo explicitly reject fake HA language?

Strong anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

What Class 2 can honestly prove:

- the dream
- the anti-goals
- the repo's intended operating contract
- the repo's own honesty wall

What Class 2 cannot prove by itself:

- that the root runtime already behaves that way
- that the docs are allowed to use present-tense capability language
- that the helper layer already exists just because the dream names it clearly

Intent is real.
Intent is not runtime proof.

### Class 3: planning and promotion evidence

Use for claims about known gaps, proposed repairs, explicit capability holes,
and possible future promotion logic.

Typical questions:

- what does the repo already know is missing?
- what repair paths are being seriously explored?
- which helper layers are already called out as untrustworthy?

Strong anchors:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
- [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
- research pages under
  [`knowledgebase/research/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/)

What Class 3 can honestly prove:

- the repo already sees the gap
- the gap is not being ignored
- a future path is being considered seriously
- a current helper is already distrusted for a named reason

What Class 3 cannot prove by itself:

- that the repair is live
- that the current runtime already inherited the planned behavior
- that a proposed helper has earned trust just because the plan is detailed

Detailed plans are still plans.

### Class 4: archive-pressure evidence

Use for claims about what the user keeps rejecting, what pain keeps recurring,
and which answer families still feel fake.

Typical questions:

- why does "just use X" keep failing here?
- what exactly is the user rebelling against?
- why does the repo keep talking about wrong-node dignity and sacred-node
  memory?

Strong anchors:

- [Archive Pressure Patterns](archive-pressure-patterns.md)
- [User Intent and Dream](user-intent-and-dream.md)
- source archive files under
  [`knowledgebase/source-archive/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/)

What Class 4 can honestly prove:

- the standard the user is holding the system to
- the kinds of answers that should not satisfy the docs
- the recurring shape of the user's disappointment

What Class 4 cannot prove by itself:

- that the local repo already solved the pain
- that one future technical choice has already been earned
- that the runtime is close just because the archive pressure is clear

Archive pressure tells us what counts.
It does not certify completion.

## Claim routing matrix

Use this before writing any substantial paragraph.

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "this exists in the current priority runtime" | Class 1 | root Compose runtime | that it survives the relevant failure |
| "this is the dream or operating contract" | Class 2 | `copilot-instructions.md`, `README.md` | that the worktree already earned it |
| "the repo already knows this is missing or broken" | Class 3 | planning docs and gap pages | that the repair is active |
| "the user rejects this whole answer family" | Class 4 | archive-pressure and source archive | that rejection alone proves the better design |
| "this lane is getting closer but still unsafe" | Class 1 + 3 | runtime plus named gap | that "closer" is the same as proven |
| "this helper matters but is not yet trustworthy" | Class 1 + 3 | live helper presence plus explicit distrust in plans | that helper presence equals route durability |

If a sentence cannot be routed cleanly, it is usually too vague or doing too
much work at once.

## The anti-merger rule

This is the most important rule on the page:

do not let several weaker evidence classes emotionally combine into a stronger
claim than any one of them actually proves.

That means:

- Class 1 plus Class 2 does not equal shipped capability
- Class 2 plus Class 3 does not equal near-complete behavior
- Class 1 plus Class 3 does not equal working failover
- Class 1 plus Class 4 does not equal the user's dream being met
- Class 2 plus Class 4 does not equal justified orchestrator promotion

This repo is especially vulnerable to this because it has:

- strong architecture intent
- a serious live edge stack
- many explored control-plane directions
- a very clear archive of dissatisfaction

Put together, those can sound like completion.
This ledger exists to stop "sounds like completion" from replacing proof.

## The five fields every serious page should make recoverable

If a page is behaving like a retrieval surface rather than a polished essay, a
reader should be able to recover all five of these quickly:

1. the claim
2. the evidence class
3. the confidence ceiling
4. what the evidence genuinely proves
5. what it still does not prove

The fifth field is the most important.

This repo does not mainly need more confident prose.
It needs prose that stops exactly where the evidence stops and explains why
that stop point is still painful.

## Common inflation patterns this ledger should block

The ledger should stop all of these:

- live-runtime claims made from planning text
- failover claims made from proxy presence alone
- distributed-truth claims made from modular Compose authoring
- stateful HA claims made from ingress cleverness
- architecture-closure claims made from exploration artifacts
- smooth wording upgrades where "dynamic," "multi-node," or "resilient" begin
  to mean more than the evidence supports
- emotional satisfaction with the dream being mistaken for progress in the
  worktree

Those are not minor wording bugs.
They are how the docs become satisfying for the wrong reason.

## Practical examples

### Correct

- "the repo wants any-node entry plus peer forwarding"
- "the root runtime already has a substantial Traefik/auth/CrowdSec edge"
- "the master plan explicitly records service failover as still missing"
- "the archive shows the user rejects answers that preserve sacred-node
  memory"
- "the current stack has fallback-shaped glue, but the plans explicitly record
  that it can lose routes under container stop"

### Incorrect

- "the stack already has resilient multi-node service failover"
- "the system already behaves like one distributed platform"
- "stateful HA is mostly in place because the services are routed"
- "the missing middle layer is effectively solved"
- "because several sources all point the same way, the runtime is close"

The difference is not tone.
It is evidence-class discipline.

## What the ledger should make the rest of the docs do

Every serious page should be forced to answer:

1. what question is this page really answering?
2. which source class is doing the heavy lifting?
3. what nearby stronger sentence would still be a lie?
4. where does the current answer still depend on operator reconstruction?
5. what would have to be observed next before the stronger sentence becomes
   legal?

That fourth question matters because the user's complaint is not only about
missing features.
It is about private operator memory still doing system work.

If the docs cannot say where that hidden burden still lives, they have not
really assimilated the repo.

## The bottom-line test

The knowledgebase now has enough material to sound convincing almost
everywhere.
That makes this ledger more important, not less.

The main threat is no longer missing nouns.
It is:

- runtime
- plus dream
- plus plan
- plus archive pressure

being emotionally read as one present capability story.

This page exists to keep those joins explicit and to stop the docs from
becoming more coherent by becoming less faithful.
