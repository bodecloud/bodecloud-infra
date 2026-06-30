# Evidence Ledger

This page is the claim governor for the whole knowledgebase.

It exists because the main documentation failure in this repo is no longer
missing information.
It is information being narrated at the wrong confidence level.

The user's complaint is not merely that previous docs were too short.
It is that they kept sounding more complete, more settled, and more proved
than the repo had actually earned.

This page exists to stop that from happening again.

## The problem this ledger is trying to prevent

`bolabaden-infra` has enough real pressure, enough real proxy machinery, enough
real research, and enough real planning artifacts to sound finished long
before it actually is.

That creates one especially dangerous failure mode:

1. archive pressure names the wound clearly
2. repo-native intent names the dream clearly
3. live runtime shows serious components
4. planning text proposes sharp repairs
5. the docs quietly fuse those into one present-tense feeling of progress

That fusion is what this page is here to block.

The repo can survive imperfect implementation.
It should not survive narrative overpayment.

## What this page is and is not allowed to prove

This page is allowed to:

- define evidence classes for the rest of the knowledgebase
- route claim types to the right authority surfaces
- explain where confidence ceilings stop
- keep contradiction visible instead of smoothing it away

This page is not allowed to:

- serve as direct runtime proof for a specific drill or route
- imply that well-governed claims are already fulfilled claims
- merge archive pressure, architecture intent, runtime artifacts, and plans
  into one present-tense capability story
- substitute confidence discipline for implementation maturity

This page is about how to speak honestly, not about pretending honesty itself
is completion.

## The four evidence classes

Every serious claim in this site should be routed through one of four evidence
classes.

## Class 1: live runtime evidence

Use for claims about what the current tracked implementation actually contains
or does.

Examples:

- which networks are defined in `docker-compose.yml`
- whether Traefik, CrowdSec, TinyAuth, or MongoDB are part of the active root
  runtime
- whether Headscale currently uses SQLite in the active fragment
- whether a service definition, healthcheck, label, or config surface is
  present now

Strong anchors:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active files under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)
- validated rendered output such as `docker compose config`
- route-specific drill results where they exist

What Class 1 cannot do by itself:

- prove the user's full dream
- prove distributed behavior the current runtime only gestures toward
- prove future failover semantics from present components alone

## Class 2: repo-native intent evidence

Use for claims about what the repo explicitly wants to become.

Examples:

- local-first then peer-forward request philosophy
- Compose-first without heavy orchestration by default
- any-node public entry desire
- refusal to merge HTTP, TCP, and stateful maturity into one fake story

Strong anchors:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

What Class 2 cannot do by itself:

- prove the root runtime already behaves that way
- justify upgrading present-tense capability language

Intent is real.
Intent is not runtime proof.

## Class 3: planning and promotion evidence

Use for claims about known gaps, proposed repairs, and future promotion logic.

Examples:

- `docker-gen-failover` being treated as an unreliable current answer
- secret sync and compose sync still being unresolved
- service failover being a named capability gap
- OpenSVC, Nomad, k3s, or similar paths being explored rather than adopted

Strong anchors:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
- research pages under
  [`knowledgebase/research/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/)
- architecture pages that explicitly frame promotion or maturity boundaries

What Class 3 cannot do by itself:

- prove a repair is live now
- upgrade roadmap detail into current capability

Detailed plans are still plans.

## Class 4: archive-pressure evidence

Use for claims about what the user keeps rejecting, what pain keeps recurring,
and which answer families still feel fake.

Examples:

- dissatisfaction with ordinary load-balancer advice
- repeated refusal of both brittle static glue and premature orchestrator
  capture
- recurring demand for wrong-node dignity and shared current truth

Strong anchors:

- [`archive-pressure-patterns.md`](archive-pressure-patterns.md)
- source archive files under
  [`knowledgebase/source-archive/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/)

What Class 4 cannot do by itself:

- prove the local repo has already solved the pain
- certify one future technical choice as locally proven

Archive pressure tells us what counts.
It does not tell us that the repo has already crossed that line.

## The anti-merger rule

This is the most important rule in the whole page:

do not let several weaker evidence classes emotionally combine into a stronger
claim than any one of them actually proves.

That means:

- Class 1 plus Class 2 does not automatically equal shipped capability
- Class 2 plus Class 3 does not automatically equal near-complete behavior
- Class 1 plus Class 3 does not automatically equal working failover
- Class 1 plus Class 4 does not automatically equal the user's dream being met

The repo is especially vulnerable to this because it has:

- strong architecture intent
- a serious live edge stack
- many explored control-plane directions
- a very clear archive of dissatisfaction

Put together, those can sound like completion.
This ledger exists to stop "sounds like completion" from replacing proof.

## The five fields every serious page should make recoverable

If a page is actually doing retrieval-like reconstruction rather than polished
summarization, a reader should be able to recover all five of these quickly:

1. the claim
2. the evidence class
3. the confidence ceiling
4. what the evidence genuinely proves
5. what it still does not prove

The fifth field is the most important.

This repo does not mainly need more confident prose.
It needs prose that stops exactly where the evidence stops, and explains why
that stop point is still painful.

## Common documentation inflation patterns to block

The ledger should stop all of these:

- live-runtime claims made from planning text
- failover claims made from proxy presence alone
- distributed-truth claims made from modular Compose authoring
- stateful HA claims made from ingress cleverness
- architecture-closure claims made from exploration artifacts
- smooth wording upgrades where "dynamic," "multi-node," or "resilient" start
  meaning more than the evidence supports

Those are not minor phrasing bugs.
They are the exact path by which the docs become satisfying for the wrong
reason.

## Practical examples of how to speak correctly

Correct:

- "the repo wants any-node entry plus peer forwarding"
- "the root runtime already has a substantial Traefik/CrowdSec/auth edge"
- "the plan clearly identifies service failover and secret sync as missing"
- "the archive shows the user rejects answers that preserve sacred-node memory"

Incorrect:

- "the stack already has resilient multi-node service failover"
- "the system already behaves like one distributed platform"
- "stateful HA is mostly in place because the services are routed"
- "the missing middle layer is effectively solved"

The difference is not tone.
It is evidence class discipline.

## Strongest honest current answer

The knowledgebase now has enough material to sound convincing nearly
everywhere, which makes this ledger more important, not less.

The main threat is no longer missing nouns.
It is adjacency masquerading as proof:

- runtime plus intent plus archive plus planning getting emotionally read as
  present capability

This page exists to keep those joins explicit and to stop the docs from
becoming more coherent by becoming less faithful.
