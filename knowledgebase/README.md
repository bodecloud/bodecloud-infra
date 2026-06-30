# Knowledgebase README

This file is the repo-facing doorway into the MkDocs site rooted in
[`index.md`](index.md).

It exists because many contributors open `README.md` first, and this repo is
too easy to misread if they do not get the framing immediately.

It is also here because the repo is surrounded by exactly the kind of
infrastructure language that can sound smart while refusing to say the one
thing the user is actually mad about:

the stack keeps looking flexible until the real topology still has to live in
one human head.

That line is not a flourish.
It is the practical doorway test for this whole directory.
If a reader can walk through the knowledgebase and still leave mainly thinking
"this is a nice Compose homelab with some future HA ideas," then the doorway
already failed before the deeper pages had a chance to help.

## Read this directory correctly

This directory is not here to provide generic infrastructure prose.
It is here to keep one specific documentation failure from happening again:

the repo sounding more solved, more uniform, and more platform-like than the
current evidence allows.

That failure is not just cosmetic.
It is the documentation form of the same fake-option problem the repo is
trying to escape:

- one answer sounds simple but hides static glue
- another answer sounds serious but imports a giant worldview
- the real missing middle remains unnamed

`bolabaden-infra` is not merely:

- some Docker Compose files
- some Traefik labels
- some future thoughts about Swarm, Nomad, or Kubernetes

It is a focused attempt to answer a much sharper problem:

> how do you make multiple ordinary Docker nodes behave like one resilient,
> readable, anti-SPOF personal cloud without paying the full tax of a
> heavyweight orchestrator before it has truly earned its keep?

If a page loses that question, it becomes noise.
If a page answers only a tidier neighboring question, it also becomes noise.

If a page answers that question too politely, it also becomes noise.
The user is not asking for a mellow orientation guide.
They are asking for docs that can withstand frustration, contradiction, and
the feeling that every apparent option secretly preserves the same hidden tax.

That hidden tax is not generic complexity.
It is the moment where the request lands on the wrong node, a backend
disappears, or a stateful claim gets serious, and the system quietly needs the
operator to finish the sentence anyway.

## The shortest honest reconstruction of the wound

The wound is not merely "multi-node Docker is hard."

The wound is that the stack keeps sounding flexible until a real request lands
on the wrong node and the truth suddenly collapses back into:

- private placement memory
- sacred-node folklore
- unclear peer eligibility
- a fallback story that may exist only as architecture narration

The archive is useful here because it restores how plainly the user stated the
dream.
The user already accepted manual placement and already had Cloudflare plurality.
The thing they kept asking for was the missing service-discovery and
wrong-node-preservation layer.
That is why this site cannot keep drifting back to gentler questions like
"which orchestrator feels right?" without first naming what specific burden is
supposed to move.

That is why these docs cannot behave like ordinary repo docs.

Ordinary repo docs can get away with:

- broad overviews
- tool catalogs
- architecture summaries
- future-roadmap confidence

This repo cannot.

Here, a page fails if it makes the system sound more coherent than the burden
transfer really is.

## The real standard a reader needs to carry

The practical test for almost every page in this site is:

> did this page explain which truth became system-owned, or did it mostly
> explain the available tools and plans more elegantly?

If the second answer is stronger, the page is still too weak.

That is the standard contributors should carry while reading or editing this
knowledgebase.

It is also the easiest way to detect fake progress in docs form.

A documentation pass is not strong here because it:

- sounds more mature
- covers more concepts
- names more candidate platforms
- has more complete reading routes

It is strong only if it makes the repo harder to overread and easier to map to
the exact burden the user is trying to eliminate.

That means the doorway should leave the reader with one uncomfortable sentence
still active in their head:

> if the system still needs a human to privately know what node, peer, route,
> or owner is really authoritative, then the main wound is still here.

## What this page is and is not allowed to prove

This page is allowed to prove:

- how the knowledgebase should be entered without immediately downgrading the
  question
- which reading routes are safest for reconstructing the user's real demand
- that the site exists to resist fake convergence and fake option language

This page is not allowed to prove:

- that the runtime already crossed the thresholds the route describes
- that reading order itself is evidence of implementation maturity
- that the repo is already settled just because the site has become broader

This is a doorway contract, not a runtime proof page.

## What still does not count as a useful repo-facing doorway

This README should also say more explicitly what fake helpfulness still looks
like.

These still do not count:

- a contributor-facing intro that mainly catalogs folders
- a calm summary that never names the sacred-node or wrong-node pressure
- a reading order that sounds complete but does not say what stronger claim is
  still forbidden
- a doorway that makes the docs feel mature without warning that maturity is
  not the same as runtime ownership
- a repo intro that makes the problem sound like "choosing an orchestrator"

The user is not asking for a nicer README tone.
They are asking for a doorway that does not betray the core question before the
reader even reaches the deeper pages.

If this file helps someone browse but still leaves them with a smaller, calmer,
or more product-shaped question, then it is still too weak.

That means the doorway pages have to do more than orient.
They have to keep a reader from leaving with the wrong calm conclusion.

They also have to keep one sharper distinction visible than most READMEs ever
try to hold:

- a request looked successful because the system itself owned the truth needed
  to preserve it
- a request looked successful because a human operator still knew which hidden
  node, route, storage path, or fallback story to mentally reconstruct

If that distinction disappears, the docs can sound comprehensive while still
teaching the same lie the user is tired of hearing:

the stack is flexible now, so the operator burden must already be lower.

That lie is exactly what this doorway is trying to break before the reader even
clicks into architecture or operations pages.
The user is not short on flexible-looking stacks.
They are short on options that stay honest when the topology question stops
being theoretical.

In this repo, a page can be technically informative and still fail if it lets
the reader walk away thinking:

- the options are clearer than they really are
- the authority surfaces all say the same kind of thing
- the runtime is closer to the dream than current evidence supports

## Start here

If you are reading in the repo, start with:

- [`index.md`](index.md)

Then follow this order:

- [`reading-paths.md`](reading-paths.md)
- [`research/user-intent-and-dream.md`](research/user-intent-and-dream.md)
- [`architecture/problem-and-goals.md`](architecture/problem-and-goals.md)
- [`architecture/operator-contract.md`](architecture/operator-contract.md)
- [`architecture/instruction-surfaces-and-authority.md`](architecture/instruction-surfaces-and-authority.md)
- [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
- [`architecture/failure-model-and-maturity.md`](architecture/failure-model-and-maturity.md)
- [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md)
- [`architecture/ha-failover-routing.md`](architecture/ha-failover-routing.md)
- [`operations/operator-questions-and-honest-answers.md`](operations/operator-questions-and-honest-answers.md)
- [`architecture/compose-first-architecture.md`](architecture/compose-first-architecture.md)
- [`architecture/missing-middle-layer.md`](architecture/missing-middle-layer.md)
- [`architecture/capability-gaps-and-roadmap.md`](architecture/capability-gaps-and-roadmap.md)
- [`architecture/stateful-ha-and-data.md`](architecture/stateful-ha-and-data.md)
- [`operations/decision-paths-and-promotion-rules.md`](operations/decision-paths-and-promotion-rules.md)
- [`research/evidence-ledger.md`](research/evidence-ledger.md)
- [`operations/proof-matrix-and-drills.md`](operations/proof-matrix-and-drills.md)
- [`operations/devops-runbook.md`](operations/devops-runbook.md)

That order is intentional.
It starts with the question-first reading map, then the dream, then the actual
problem statement, then the success contract, then clarifies which repo files
own the dream versus the runtime anchor, then forces the live runtime and its
uneven maturity back into view, then rebuilds the wrong-node and failover
pressure, then names the missing middle layer, then closes on evidence and
proof boundaries.

That sequence matters because this repo is easiest to misunderstand when read
in the opposite order.
If someone starts from:

- orchestration options
- compose fragment lists
- implementation details

without first understanding the user's hidden negative benchmark, they will
almost always flatten the problem into "wants better HA" or "needs a cluster."
Those readings are too weak.

There is also a second-order reading failure this sequence is trying to stop:

- the docs are now broad
- broad docs feel like convergence
- convergence feels like the repo must already be mostly settled

That is also too weak.
The site is getting better at reconstructing the pressure.
That does not mean the runtime has already paid down the pressure.

## Fast question-first routes

If you are not trying to read the whole site, do not browse by folder.
Browse by the actual question you need answered.

### "What is the shortest serious way to read this repo without being fooled?"

Read:

- [`reading-paths.md`](reading-paths.md)
- [`research/user-intent-and-dream.md`](research/user-intent-and-dream.md)
- [`architecture/instruction-surfaces-and-authority.md`](architecture/instruction-surfaces-and-authority.md)

### "What is the dream, really?"

Read:

- [`research/user-intent-and-dream.md`](research/user-intent-and-dream.md)
- [`architecture/problem-and-goals.md`](architecture/problem-and-goals.md)
- [`architecture/operator-contract.md`](architecture/operator-contract.md)

### "What does the root runtime actually prove?"

Read:

- [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
- [`architecture/compose-fragment-map.md`](architecture/compose-fragment-map.md)
- [`architecture/failure-model-and-maturity.md`](architecture/failure-model-and-maturity.md)

### "Why is wrong-node entry still the humiliating threshold?"

Read:

- [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md)
- [`research/ingress-and-failover-evidence.md`](research/ingress-and-failover-evidence.md)

### "Why do most apparent options still feel fake?"

Read:

- [`architecture/missing-middle-layer.md`](architecture/missing-middle-layer.md)
- [`architecture/capability-gaps-and-roadmap.md`](architecture/capability-gaps-and-roadmap.md)
- [`research/orchestrator-tradeoffs-evidence.md`](research/orchestrator-tradeoffs-evidence.md)
- [`research/orchestration-research-2026.md`](research/orchestration-research-2026.md)

### "What is still missing even if the proxy story gets better?"

Read:

- [`architecture/stateful-ha-and-data.md`](architecture/stateful-ha-and-data.md)
- [`research/stateful-ha-evidence.md`](research/stateful-ha-evidence.md)

## Anti-cheat reading rules

Do not treat these as equivalent:

- route reachability and request preservation
- multiple public nodes and no sacred node
- a rendered stack and a distributed truth layer
- a clean future plan and a live ownership surface
- a stronger-sounding platform and a genuinely different burden transfer

This README should make one thing obvious before any contributor goes deeper:

if a page answers a calmer neighboring question than the one the user is
actually asking, it may still be useful, but it is not sufficient.

The practical test is simple:

- does the page help explain what burden became system-owned?
- or does it mostly explain the available tools more elegantly?

If the second answer is stronger, the page is still too weak for this repo.

## What a serious repo-entry packet should leave behind

This README should not merely say where to click next.
It should help a reader preserve the minimum packet needed to not misread the
repo.

That packet should include:

- the main burden the user is trying to escape
- the first truth register the reader needs
- the first runtime-facing page they should inspect
- the first dream- or archive-facing page they should inspect
- the stronger sentence that remains forbidden after the first read

If a reader can finish this file and still say only "this is a multi-node
Docker project with some HA ambitions," then the doorway failed.

If they can instead say:

- "this repo is trying to stop wrong-node recovery from depending on one human
  head"
- "it is resisting both static-glue denial and premature orchestrator
  surrender"
- "the docs now tell me exactly which stronger claims are still illegal"

then the doorway is finally doing its job.

## What this rewrite is trying to fix

Earlier docs were not mainly too short.
They were too flattening.

They mixed together:

- live Compose truth
- architecture intent
- planning documents
- exploratory control-plane ideas
- archive-derived frustration
- partial edge evidence

until everything started sounding equally real.

That flattening made the docs feel polished while removing the only distinction
that matters in this repo:

> this is the dream  
> this is the next-layer plan  
> this is what the tracked worktree actually proves

The rewritten knowledgebase is intentionally stricter.
It treats overclaiming as a documentation defect.

It also treats vagueness as a defect.
In this repo, vague prose is not neutral.
It actively pushes the reader back toward the same architecture theater the
user was already frustrated by.

There is another failure mode the rewrite is trying to stop:

- the page sounds mature
- maturity sounds like convergence
- convergence sounds like ownership
- ownership gets attributed to the runtime even when the runtime still depends
  on private operator reconstruction

That sequence is how "good docs" become another form of fake HA language.

That is why so many pages now sound harsher than a normal knowledgebase.
They are trying to preserve the actual psychological pressure of the repo:

- too many tools sound like options without changing the operator burden
- too many failure stories stop at the first hop
- too many docs narrate a distributed system before the truth layer exists

The site also needs to preserve something normal docs usually try to clean up:

- competing futures visible in the same tree
- partial evidence that does not yet compose into one calm story
- real contradictions between intent, side paths, and live proof

Those are not editing failures here.
They are part of the evidence.

That is also why this README has to be harsher than a normal contributor
doorway.
It is not enough to help readers "find the right page."
It has to stop them from importing the wrong genre assumptions in the first
place.

If a contributor arrives expecting:

- a stable platform taxonomy
- one calm architecture story
- one obvious mature promotion path

this directory should interrupt that expectation immediately.

It should also interrupt a subtler expectation:

- if the docs are now broad, serious, and highly cross-linked, the repo must
  already be closer to settled than before

That expectation is dangerous here.

This rewrite is not trying to create the feeling of settlement.
It is trying to create a better reconstruction of:

- what the user is actually demanding
- where the live worktree still falls short
- which futures are real, fake, or still too weakly proved

If a contributor leaves this directory with a calmer feeling but a blurrier
understanding of those distinctions, the docs have regressed.

## The fastest useful reading route

This knowledgebase is not best read as a taxonomy.
It is best read as a pressure chain.

If someone only has limited time, this is the shortest route that still
preserves the real dream:

1. [`research/user-intent-and-dream.md`](research/user-intent-and-dream.md)
2. [`operations/operator-questions-and-honest-answers.md`](operations/operator-questions-and-honest-answers.md)
3. [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md)
4. [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
5. [`architecture/missing-middle-layer.md`](architecture/missing-middle-layer.md)
6. [`operations/decision-paths-and-promotion-rules.md`](operations/decision-paths-and-promotion-rules.md)
7. [`research/evidence-ledger.md`](research/evidence-ledger.md)

That route exists because many other reading orders make the repo feel calmer
than it really is.
This one keeps the wound, the runtime, the gap, and the proof ceiling visible
at the same time.

The site is also trying to preserve a sharper exit condition for every major
page:

after reading it, a contributor should be able to say:

- what the user is actually trying to make true
- what the current worktree actually proves
- what is still intent, promotion work, or exploration
- which "options" are fake because they leave the same hidden burden intact

If a page cannot support that exit condition, it is still too soft.

There is one more exit condition the doorway pages should preserve:

- can the reader name the exact hidden burden each proposed layer is supposed
  to remove, rather than merely naming the tool or platform family?

Without that, the docs drift back into respectable platform browsing instead of
answering the repo's real question.

## The three truth layers you must preserve

### 1. Live implementation truth

This is what the current repo and merged Compose surface actually show.

Examples:

- root [`docker-compose.yml`](../docker-compose.yml)
- active fragments under [`../compose/`](../compose)
- inventories and validation derived from `docker compose config`

This layer can prove what is authored and present now.
It cannot prove distributed resilience just because the graph parses or the
edge is complex.

### 2. Planned architecture truth

This is what the repo is seriously trying to build next.

Examples:

- [`../docs/INFRASTRUCTURE_MASTER_PLAN.md`](../docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`../docs/stateful_ha_plan.md`](../docs/stateful_ha_plan.md)
- [`../docs/osvc_ingress_ha.md`](../docs/osvc_ingress_ha.md)

This layer can prove direction and recognized gaps.
It cannot prove the promoted layer is already live.

### 3. Research-pressure truth

This is the archive and the synthesis pages derived from it.

Examples:

- `source-archive/`
- pages under `research/`

This layer can prove the user’s repeated demands, refusals, and standards.
It cannot prove runtime behavior.

Those three layers are not editorial preferences.
They are the main mechanism preventing this knowledgebase from collapsing back
into polished ambiguity.

## The current honest summary

The current docs are now built around a harder but more useful summary:

- the dream is explicit
- the root Compose implementation is real
- the edge stack is already serious
- the repo genuinely wants anti-SPOF, any-node entry, local-first serve, and
  peer-forward fallback
- the tracked worktree still does not generically prove wrong-node request
  preservation
- the tracked worktree still does not generically prove stateful HA correctness

## The fake clarity traps this directory has to keep preventing

The knowledgebase regresses whenever one of these starts sounding like enough:

- "there are multiple public nodes now"
- "the proxy stack is sophisticated"
- "there is a registry-shaped idea"
- "there is a failover generator"
- "there are several orchestration futures in the tree"

Those can all be true and still leave the same hidden burden intact.

The real question remains:

- did the system itself take ownership of one more missing truth layer
- or did the docs only get better at narrating the absence of that layer?

That distinction is the whole maintenance game for this directory.

It also needs one extra sentence that stays active in the reader's head:

the site is not trying to sound more complete than the ecosystem's answers.
It is trying to sound more faithful to why those answers kept failing the user.

That last pair has to stay uncomfortable.
The moment it starts sounding like "basically there," the docs will be lying in
the same style the user is already exhausted by.

The same warning applies to any apparently successful path.

A success story is still incomplete here unless the page can say:

- what truth the runtime owned directly
- what truth was still socially reconstructed by the operator
- what would fail if the same request landed on a different healthy node on a
  worse day

That means the remaining problem is not “make the docs prettier.”
It is “keep naming the missing truth layers precisely enough that the next
implementation decision becomes obvious instead of theatrical.”

That sentence is probably the best maintenance heuristic for the whole site.
When editing, ask:

- did this page make the missing truth layer more visible?
- or did it just make the system sound more coherent?

Only the first one is aligned here.
That question is worth restating even more harshly:

> did this edit help reconstruct the actual hidden burden, or did it only make
> the repo sound more organized than it really is?

The second outcome is a regression in this directory, even if the prose looks
better.

## What the high-level pages must do now

The top-level pages in this directory should do more than summarize.
They need to keep four things visible at once:

- the dream
- the hidden operator tax
- the missing truth layers
- the proof ceiling

And they need to keep one more distinction visible than ordinary infra docs
normally tolerate:

- which options are genuinely different because they relocate truth ownership
- which options only feel different because they rename the same operator
  burden under a cleaner platform story

If any one of those disappears, the docs start sounding nicer and getting less
useful again.

They also need to keep a fifth thing visible:

- whether the apparent intelligence of the system is actually present in the
  runtime, or merely present in the operator reading the page

## The minimum exit condition for any rewritten summary page

After reading a top-level or summary page, a contributor should be able to say:

- what the user is actually trying to make true
- what the current worktree really proves
- which source class the page leaned on for each major claim
- which options are still fake because they relocate almost no hidden burden
- what the next real proof would need to show

If a page cannot support those five answers, it is still too soft for this
repo even if the prose is cleaner.

That is the wrong trade for this repo.
Niceness is cheap.
A site that preserves the actual failure pressure is harder to write and much
more useful.

## The pages that now carry the most weight

If you only have time for the highest-signal pages, use these:

- [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md)
  This is the literal request-path page. It forces the docs to stay concrete.
- [`operations/operator-questions-and-honest-answers.md`](operations/operator-questions-and-honest-answers.md)
  This answers the recurring archive questions directly.
- [`architecture/missing-middle-layer.md`](architecture/missing-middle-layer.md)
  This names the actual gap between raw Compose and premature orchestration.
- [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
  This explains what the root runtime really contains and what it still does
  not prove.
- [`research/evidence-ledger.md`](research/evidence-ledger.md)
  This page keeps the rest of the docs from becoming architecture theater.

## What not to do when editing these docs

Do not:

- narrate DNS reachability as if it were preserved service success
- narrate proxy presence as if it were proven failover
- narrate described `services.yaml` intent as if it were live tracked root truth
- narrate TCP exposure as if it were stateful correctness
- narrate orchestration exploration as if the repo has already chosen a final
  control plane
- narrate “multi-node” as if that alone solves the sacred-entrypoint problem

In this repo, each of those mistakes creates false clarity.

Also do not:

- narrate a named future helper or agent as if naming it proves it exists
- narrate a control-plane direction as if that direction has already earned
  promotion
- narrate "route still answered" as if that proves request meaning survived
- narrate archive frustration as if it were just mood instead of architectural
  evidence

## Maintenance rule

When infrastructure behavior changes, update the matching documentation in the
same turn.

The most likely matching pages are:

- runtime surface:
  [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
- request-path or fallback behavior:
  [`architecture/request-path-and-failure-walkthrough.md`](architecture/request-path-and-failure-walkthrough.md)
  and
  [`architecture/ha-failover-routing.md`](architecture/ha-failover-routing.md)
- stateful topology:
  [`architecture/stateful-ha-and-data.md`](architecture/stateful-ha-and-data.md)
- proof boundaries:
  [`research/evidence-ledger.md`](research/evidence-ledger.md)
- operator procedures and verification:
  [`operations/devops-runbook.md`](operations/devops-runbook.md)
  and
  [`operations/proof-matrix-and-drills.md`](operations/proof-matrix-and-drills.md)

## Rendered site behavior

MkDocs renders this knowledgebase through [`../mkdocs.yml`](../mkdocs.yml).

This `README.md` is for repo readers, not the main rendered site navigation.
The real site entrypoint remains:

- [`index.md`](index.md)

## Bottom line

This directory is now supposed to do one thing well:

help an operator understand the real dream, the real live runtime, the real
missing layer, and the real proof ceiling without having to reverse-engineer
all four from vague prose.

If it ever starts feeling like a normal infrastructure documentation set again,
that is probably a sign it has drifted away from what the user actually asked
for.
