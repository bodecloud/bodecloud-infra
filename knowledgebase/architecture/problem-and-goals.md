# Problem, Pressure, and Goals

This page defines the actual architecture problem in `bolabaden-infra`.

The problem is not:

- "more clustering"
- "better deployment hygiene"
- "more modern infrastructure"
- "better docs about self-hosting"

The problem is:

> how do several ordinary Docker nodes become one request-preserving personal
> cloud, while `docker-compose.yml` stays readable and the system stops leaning
> on one operator to remember where the real topology truth lives?

That question is stricter than generic HA language, and smaller than a
"migrate to Kubernetes" story.

It is also harsher than many self-hosting discussions are willing to be.
This repo is not just looking for resilience.
It is looking for a system that stops humiliating the operator by revealing,
too late, that the operator is still the real keeper of placement and recovery
truth.

That means a better problem statement is not problem reduction.
A sharper benchmark is not one more serious option.
The page can become more exact and still leave the same unsolved burden
perfectly intact.

That sentence is not rhetorical decoration.
It is the practical acceptance test behind almost every future architecture
choice in this repo.

## What this page is and is not allowed to prove

This page is authoritative about:

- the real problem the repo is trying to solve
- the concrete requirement stack implied by that problem
- which adjacent answers are still too small

This page is not authoritative about:

- whether the current runtime already satisfies those requirements
- whether one future helper layer has already won
- whether naming the problem cleanly means the remaining gap is small

This page is the benchmark, not the completion report.

It has to stay that way because the repo is already coherent enough to tempt
readers into thinking the remaining gap is mostly execution.
That temptation is part of the problem.

It also cannot become a page where diagnosis quality starts impersonating
solution quality.

That warning needs to stay explicit because this repo is now good enough at
self-diagnosis to create counterfeit comfort.
If the benchmark becomes precise enough that the remaining gap starts sounding
small, then the page has started paying out reassurance it has not earned.

This repo already has enough serious language that the problem statement itself
can start behaving like one more fake comfort surface:

- the benchmark is sharp
- the gaps are named
- the goals are legible
- therefore the missing part must be mostly implementation effort

That is exactly the downgrade this page has to block.

It is not just a reading mistake.
It is one of the main ways the repo's own sophistication can start betraying
the user:

- the diagnosis is sharper
- the components are real
- the planning is serious
- the surviving humiliation starts sounding like implementation residue instead
  of the central unsolved fact

The more serious the stack looks, the easier it becomes to understate the
remaining burden-transfer problem.

## Strongest honest current answer

The repo already proves a serious Compose-first platform with a strong edge,
many services, and real planning work around failover and anti-SPOF behavior.

What it does not yet prove is the harder thing:

- that wrong-node entry is survivable generically
- that service placement truth is shared explicitly
- that fallback routes survive backend loss
- that peer forwarding preserves the same auth, middleware, and routing meaning
- that stateful services have honest failover semantics

That gap is the actual problem. Everything else is supporting detail.

More bluntly:

- the repo already has enough machinery to look serious
- it still does not have enough system-owned truth to make that seriousness
  fully trustworthy

That is why adjacent summaries keep drifting toward:

- better HA
- better orchestration
- better service discovery
- better failover language

while still missing the repo's harder sentence:

> when the request lands on the wrong healthy node, who still has to finish
> the story?

That distinction is one of the main reasons ordinary summaries keep failing the
user.

They keep summarizing the repo as if its main problem were insufficient
machinery.
The harder truth is that the repo already has enough machinery to make false
closure sound plausible.

The repo is therefore past the stage where "needs more components" is an
honest diagnosis by itself.
The danger now is that real sophistication can start borrowing credibility for
a burden transfer that never actually happened.

That sentence should remain brutal.

The repo is not starving for seriousness theater.
It is already serious enough that the docs can start accidentally helping the
same lie the user keeps rejecting:

- the stack feels mature
- therefore the remaining problem must be mostly polish

This page has to keep saying no to that move.

One more boundary belongs here plainly:

- better problem understanding is not architecture progress

The first is necessary.
The second only happens when one more decisive truth stops living in private
operator folklore.

It also has to keep refusing a subtler move:

- the dream is now reconstructed faithfully
- the benchmark is now severe and concrete
- the missing burdens are now named well
- therefore the rest starts sounding like implementation labor rather than a
  still-unmet platform contract

That conclusion is still too soothing.

## The shortest exact problem statement

The repo is trying to solve this:

> keep Compose as the main authoring and operator surface, but add just enough
> shared truth that a request landing on the wrong healthy node does not turn
> into guesswork, folklore, redirects, or fake failover.

The repo is therefore not only a hosting system.
It is a search for a smaller honest control surface.

The key word there is "honest."
The user is not just searching for a smaller layer.
They are searching for one that really carries burden instead of merely
renaming it.

That last clause has to stay literal.

The repo is not searching for:

- a cleaner abstraction
- a calmer story
- a more modern control plane

unless one of those also changes who owns the truth on the bad day.

That is why this page should stay impatient with adjacent improvements.

Better routing expression, better healthchecks, better sync, or better helper
generation are all real.
They are still not the answer if the wrong-node scene still cashes out into:

- remembered placement
- remembered safe peer choice
- remembered route meaning
- remembered stateful caveats

If a candidate layer still leaves behind one more sentence that starts with:

- "well, privately we know..."
- "in practice the operator knows..."
- "normally that hostname really belongs..."

then the repo-level problem is still alive.

The whole point is to make fewer important sentences start with:

- "well, privately we know..."
- "in practice we remember..."
- "normally that node actually..."
- "as long as the operator knows..."

## The wound behind the problem statement

The problem is not only architectural.
It is experiential.

The user keeps running into the same humiliating reveal:

- the stack sounds flexible
- the stack sounds distributed
- the stack sounds full of options
- then a real request or failure exposes that the decisive truth still lives
  in one human head

That is the wound this page has to preserve.

If that experiential layer disappears, the problem can start sounding like an
ordinary engineering optimization:

- improve routing
- improve clustering
- improve topology awareness
- improve failover

Those are all real subproblems.
They are still too polite if they stop naming the lived reveal:

> when reality gets sharp, the operator is still the hidden control plane.

That line is the real center of gravity.

If the page ever gets better organized by weakening that sentence into
"operational complexity," "topology awareness," or "coordination overhead,"
then the docs have become
cleaner and less faithful at the same time.

The same thing happens if the page becomes more comprehensive while quietly
making the wound feel more procedural than humiliating.
The user is not merely reporting friction.
The user is reporting that many grown-up-looking answers still collapse back
into private operator completion the moment the request path becomes real.

That is why this page cannot be allowed to become a polished executive summary
of the repo's pain.
If it starts sounding like the main hard work is now communication quality,
then it is already betraying the thing it is supposed to guard.

## What the user is actually trying to make impossible

The dream is easier to understand if stated as anti-goals.

The user is trying to make these moments impossible:

- "the wrong node got the request, so now I have to remember the real answer"
- "the fallback exists only as architecture narration"
- "the route still technically works, but no longer means the same thing"
- "the platform sounds distributed, but its decisive truths are still private"
- "the only serious alternative is to surrender to a giant control plane
  worldview"

That is why the problem statement must keep both halves visible:

- hidden operator burden is too high
- the ecosystem's offered escapes still often feel fake

If a solution narrative removes only the first half, it becomes fantasy.
If it removes only the second half, it becomes complaint without architecture.
The real documentation has to preserve both.

That is one reason this repo is hard to document honestly.

Most infra writing wants one of two safe modes:

- emotional distance
- or technical distance

This project cannot afford either one by itself.
The whole point is that the technical missing truths and the emotional insult
of those missing truths are the same wound seen from two angles.

## Three bars every future option has to clear

Any future helper, registry, agent, scheduler, or orchestrator candidate has
to clear all three of these bars:

### 1. Honesty bar

It must make the system easier to describe without inflating what is actually
proven.

### 2. Dignity bar

It must remove a real bad-day private reconstruction step instead of merely
making the surrounding machinery look more respectable.

### 3. Legibility bar

It must not charge so much control-plane worldview tax that the repo loses the
directness that makes the root Compose surface valuable in the first place.

This is one of the places where the knowledgebase most needs to "actually RAG"
instead of merely summarize.

A shallow summary can preserve the nouns while dropping the pressure.
A faithful reconstruction has to preserve:

- the technical missing truths
- the emotional reason those missing truths still feel insulting
- the fact that many nearby answers fail by sounding complete one layer too soon

It also has to preserve the asymmetry between understanding and escape.
This page can understand the problem very well and the repo can still be far
from one option that has earned the right to feel believable.

That balance is one of the central duties of the whole knowledgebase.

## What still does not count as understanding the problem

This page needs its own anti-flattening filter.

The following still do not count as having reconstructed the user's actual
problem:

- rewriting "anti-SPOF" in more polished language
- treating the problem as general self-hosting complexity
- treating the problem as mostly about choosing between Compose and a scheduler
- describing the topology pressure without naming private operator
  reconstruction as the wound
- assuming that more flexibility or more options are automatically relief

Those readings all stay too generic.
The user's complaint is narrower and harsher:

> too many supposedly serious answers still leave the operator as the real
> keeper of placement, fallback, and semantic continuity truth.

That sentence is the real anti-flattening rule for the page.
If a future summary cannot survive it, the future summary is too small.

It is also the sentence that should veto most "good enough" interpretations.

The user is not asking for an architecture that is broadly respectable.
The user is asking for one that stops cashing out into private truth custody at
the exact moments that currently expose the fake option problem.

## What does not count as solving this problem

The problem is specific enough that it needs an explicit false-solution filter.

The repo has **not** solved the problem merely because:

- more public nodes exist
- DNS can hit several addresses
- a reverse proxy can see more containers
- a mesh network connects the nodes
- a helper can generate routes on the happy path
- a larger orchestrator can be demoed in isolation

The repo also has not solved the problem merely because the docs can now state
the requirement stack cleanly.
Better diagnosis is still not the same thing as shared runtime ownership.

Those may all be helpful ingredients.
None of them, by themselves, prove that request meaning survives wrong-node
entry without private operator reconstruction.

That is why the repo keeps sounding unimpressed by things that would normally
count as obvious progress.
Progress is not denied.
It is just being held against a much sharper held-out scene.

And that held-out scene is the part most ecosystems keep trying to negotiate
away.

They keep offering:

- healthier first hops
- more respectable helpers
- larger control surfaces

while hoping the user will stop asking the uglier question:

> yes, but who still privately knows what the wrong node should have done?

That sharper held-out scene is what keeps this page from collapsing into a
generic roadmap preface.

Without it, the page would quietly become:

- an inventory of pain points
- an inventory of candidate layers
- an inventory of next steps

All three would be technically useful.
None would yet be faithful enough to the user's actual benchmark.

## The hidden enemy

The hidden enemy is not lack of products.

The hidden enemy is private operator reconstruction.

That reconstruction currently shows up in questions like:

- what runs where right now?
- if this hostname hits node B, does node B know the service is actually on
  node A?
- if node A is unhealthy, what route survives and who generated it?
- if a helper proxy forwards the request, do auth and middleware still mean the
  same thing?
- if Redis, MongoDB, Headscale, or another stateful service moves, what still
  owns write authority?

As long as those questions are answered mainly from memory, the platform is
still only partially system-owned.

This is the most important compression of the whole page:

- the repo does not mainly lack components
- it mainly lacks system-owned answers to the questions that become decisive on
  the bad day

That is also why the user's frustration can coexist with a stack that already
looks technically substantial.

## What the current stack already gives

The repo already has meaningful assets:

- a real root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active include fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- a public edge layer centered on Traefik, CrowdSec, TinyAuth, and
  `nginx-traefik-extensions`
- Cloudflare participation for multi-node public entry pressure
- Headscale for private-mesh coordination assumptions
- observability through Grafana, Prometheus, VictoriaMetrics, Loki, Promtail,
  cAdvisor, Blackbox, and Alertmanager
- explicit planning around service failover, secret sync, compose sync, and
  missing middle-layer helpers

Those are real strengths.

They are still weaker than the final requirement stack because most of them are
ingredients, not shared truth.

That sentence is the correct posture for almost the whole current stack.
Useful, real, often impressive, but still not the same thing as burden
ownership.

That sentence should be treated as the posture for most of the current stack.

The repo does not need to apologize for having real assets.
It needs to stop mistaking real assets for proof that the decisive truth has
already left the operator's head.

## What a serious problem-definition proof packet would have to contain

If this page ever supports the claim that the repo has correctly identified the
real problem rather than a neighboring one, it should be because the
documentation leaves behind a proof packet.

That packet should include:

- the exact hidden burden being named
- the route or service class where that burden becomes visible
- the system-owned truth that is missing today
- the human reconstruction step that still fills the gap today
- the narrower artifact or drill that would prove the burden has moved
- the explicit statement of which stronger burdens still remain

Without that packet, even a sharp problem statement can still drift back into
"interesting infra essay" territory.

## Requirement stack implied by the dream

If the dream is taken seriously, the system needs all of these:

1. **Any-node first hop**
   A request can land on any surviving public node.

2. **Local-first handling**
   If the requested service is local, the node serves it locally instead of
   pretending locality does not matter.

3. **Explicit placement truth**
   If the service is not local, the receiving node has a trustworthy current
   source for where it really lives.

4. **Peer eligibility truth**
   The receiving node can tell which peer is healthy and semantically eligible
   now, not just historically configured.

5. **Fallback-route durability**
   The route needed for rescue remains alive under the failure that made rescue
   necessary.

6. **Policy preservation**
   Auth, middleware, routing policy, and visible service meaning survive the
   handoff.

7. **Stateful honesty**
   Redis, MongoDB, Headscale, databases, and other state-bearing systems are
   described by their real authority and failover semantics, not by mere
   reachability.

8. **Inspectable ownership**
   An operator can explain why the request succeeded by reading tracked shared
   truth, not by privately reconstructing the topology.

If a design satisfies only `1` and `2`, it is still not enough.
If it satisfies `1` through `6` but cheats on `7`, it is still not enough.
If it satisfies all of those but still fails `8`, the system is still leaning
on folklore.

That is why the problem statement must stay stricter than generic anti-SPOF
language.
Without the full stack above, "resilience" becomes too easy to say and too hard
to trust.

## What progress would actually look like in this repo

A real step forward would not just add components.
It would make at least one previously private answer become system-readable.

Examples of genuine progress would look like:

- a request landing on the wrong node can now be explained from tracked
  placement truth instead of memory
- peer eligibility is now derived from a shared surface rather than inferred
  socially
- the rescue path survives the exact failure that used to delete it
- forwarded requests can be shown to preserve auth and middleware meaning
- stateful service claims get narrower, stricter, and more provable instead of
  broader and more flattering

## Why generic options lists still fail

Many neighboring answers are technically respectable and still too small here.

### "Use more DNS"

Cloudflare multi-record DNS helps with first-hop plurality.
It does not tell the wrong node where the service is, whether the peer is
healthy, or whether policy survives the handoff.

### "Use a reverse proxy"

Traefik is a real edge asset and clearly central to this repo.
Local container discovery is still not the same thing as cross-node placement
truth.

### "Use a helper generator"

The repo's own planning docs call out `docker-gen-failover` as broken because
it can delete routes when containers stop. A helper that removes the rescue
path under failure is part of the problem, not proof of the solution.

### "Use Kubernetes / k3s / Nomad / OpenSVC"

These may become valid promotions later. They are not automatically the answer
just because they are larger or more respectable.

The repo's standard is harsher:

- what truth do they own?
- what burden do they really remove?
- what new worldview tax do they impose?
- do they solve the exact wrong-node and anti-folklore problem, or mostly
  relocate it?

That harsher standard is what keeps this repo from collapsing back into the
same two humiliations the user keeps rejecting:

- static glue that still needs private human topology memory
- heavyweight control planes adopted before smaller honest answers were
  exhausted

## The practical goal

The practical goal is not "be more like cloud-native infrastructure."

It is:

> externalize enough placement, health, routing, and policy truth that the bad
> day stops depending on sacred remembered nodes and private rescue knowledge.

That is the real acceptance bar for every future control-plane, agent, helper,
registry, or orchestrator decision in this repo.

For the exact acceptance bar that follows from this page, continue to
[Operator Contract and Success Criteria](./operator-contract.md).
