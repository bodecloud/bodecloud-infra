# Orchestration Research 2026: What the Repo Is Actually Trying to Buy

This page is not a generic orchestrator comparison.

If it becomes "Compose vs Nomad vs Kubernetes," it stops describing the repo.

The actual research question inside `bolabaden-infra` is much more specific:

> what is the smallest honest layer that can make multiple ordinary Docker
> nodes behave less stupidly under drift, misplacement, wrong-node entry, and
> failure, without forcing the operator into a giant worldview shift before
> that tax has clearly earned itself?

That is the question.

There is also a deeper version of the same question that the research has to
respect:

> what kind of system is the user actually trying to bring into existence, once
> all the easier substitute questions are stripped away?

That matters because the orchestration search in this repo is not a shopping
exercise.
It is an attempt to stop being trapped between fake options.

The fake options are familiar:

- keep piling static glue onto Compose and pretend generation made it dynamic
- accept first-hop plurality as if it solved request preservation
- import a full scheduler worldview before proving the smaller missing truth
  cannot be carried another way

This page only helps if it keeps that trap visible.

Not:

- what platform is most popular
- what platform is most powerful
- what platform wins benchmark arguments
- what platform makes the stack sound more "serious"

The user is already tired of those substitutions.

The user is not short on nouns.
They are short on answers that remain honest after the request lands on the
wrong machine and the hidden sacred-node knowledge becomes relevant.

Another way to say it:

the user is not really asking for an orchestrator recommendation.
They are asking why the modern infra ecosystem keeps acting like there are only
two serious adult lives available:

- eternal static Docker plus human glue
- or total adoption of a scheduler worldview

This page should keep treating that absence of believable middle options as one
of the primary wounds, not as a side remark.

## This page is not choosing products, it is ranking missing truths

The repo does not really need "the best orchestrator."
It needs a better accounting of what is missing.

The right research question is:

> which missing truth is still forcing the operator to carry the system in
> their head, and which candidate family would own that truth without charging
> more abstraction tax than the pain is worth?

That question is much harder than a feature checklist.
It is also much closer to what the user is actually asking.

The point of the ranking is not to crown a winner.
It is to stop the repo from buying a giant worldview to compensate for a
smaller truth that was never named precisely enough.

## Evidence classes used here

This page mixes four evidence classes.
They should stay separate.

### 1. Live-root evidence

What the priority implementation already proves from the root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and its active includes.

### 2. Repo-direction evidence

What repo-native surfaces clearly say the project wants, even when the root
runtime does not yet fully implement it.

Important examples:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- [`infra/docs/ARCHITECTURE.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/infra/docs/ARCHITECTURE.md)

Important boundary inside this class:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  is the clearest statement of the desired any-node, local-first, peer-forward
  runtime dream
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
  is stronger as repo-shape and validation context than as architecture proof
- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)
  is mostly authoring discipline and Compose hygiene, not the deepest
  architecture intent surface

That distinction matters because the repo has several "instruction" files, but
they are not equal witnesses for the same claim.

### 3. Archive-pressure evidence

Imported plaintext and synthesis pages showing what the user keeps asking for
after ordinary self-hosting answers fail them.

### 4. Non-proof planning evidence

Blueprints, brainstorms, and option pages that clarify direction but do not
prove shipped behavior.

This page is mostly about repo-direction and archive-pressure evidence, with
live-root reality used as the constraint.

That mix is deliberate.

If the page used only live-root evidence, it would underspecify the user's
actual standard.
If it used only archive pressure and planning evidence, it would drift into
theory theater.

The research has to hold all three at once:

- what the user really wants
- what the repo clearly intends
- what the tracked runtime still does not prove

## What the user is actually rejecting

The orchestration research only makes sense if the refusals stay explicit.

The user is not merely saying:

- "I want HA"
- "I want clustering"
- "I want orchestration"

The user is rejecting a repeating pattern of fake choice.

The archive pressure keeps showing that the user is not simply frustrated by
complexity.
They are frustrated by answer spaces that collapse too quickly into:

- a static answer dressed up as dynamic
- a serious-sounding platform that still leaves the same human burden intact
- a documentation voice that treats partial route survival as architectural
  closure

### Refusal 1: static glue disguised as dynamic infrastructure

The repo repeatedly runs into answers that still boil down to:

- predeclare everything
- manually sync the files
- keep a human in the loop for topology truth
- call it dynamic because a template or generator was involved

That is not the dream here.

The user is explicitly trying to escape systems that only *look* smarter than
static config while still depending on hidden operator knowledge.

This is why so many apparently "reasonable" answers fail this repo.
They improve surface flexibility without relocating where the real intelligence
of the system lives.

### Refusal 2: wrong-node traffic that still depends on sacred human memory

The desired system is not just "load balanced."
It is supposed to behave like this:

- any healthy node can receive the request
- the receiving node can tell whether the target is local
- if not local, it can preserve the request through a legitimate remote path
- the operator does not privately remember all placement truth to make that
  work

That is a much stricter demand than ordinary reverse-proxy redundancy.

It is also why the user keeps sounding harsher than normal infra consumers.
The bar is not "route traffic somewhere healthy."
The bar is "stop the wrong receiving node from becoming a semantic error."

### Refusal 3: "just use Kubernetes" as a sermon instead of an answer

The repo is not anti-Kubernetes because Kubernetes is fake.

It is anti-Kubernetes-as-default-sermon because the user does not want to pay:

- control-plane tax
- storage and state tax
- workflow tax
- worldview tax

before there is solid evidence that smaller, more legible answers cannot carry
the required truths.

That is not stubbornness.
It is the repo's central evaluation principle.

The user is effectively saying:

if a big orchestrator is the right answer, prove it by naming the exact hidden
truth it takes over better than a thinner layer could.

Until then, "just use Kubernetes" is not rigor.
It is an evasion.

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

Without this refusal, the whole project degrades into pleasant infrastructure
theater.

### Refusal 5: the forced choice between raw Compose pain and total platform capture

This is probably the deepest refusal of all.

The ecosystem keeps offering two unsatisfying extremes:

- stay in hand-managed Compose sprawl
- or surrender into a full desired-state cluster worldview

The orchestration research exists because the user does not accept that as the
only serious choice set.

That is why the repo keeps searching for a missing middle rather than merely a
better product recommendation.

This missing middle is the real object of research.

Not a brand.
Not a benchmark leader.
Not a currently fashionable control plane.

The object is:

the smallest additional truth-owning layer that stops the system from
functioning as a distributed illusion.

## The four questions every orchestration answer must survive

This page should not require the reader to reconstruct the evaluation criteria
from the rest of the site.

Every candidate family must survive these four questions:

1. what hidden burden moves out of the operator's head if this layer becomes
   real?
2. what part of the wrong-node problem does it actually solve?
3. what new worldview or control-plane tax does it import?
4. what part of the user's real dream still remains unowned even if it works
   exactly as advertised?

If a candidate cannot survive those four questions, it is still closer to
platform theater than to a real answer for this repo.

## The orchestration anti-benchmark

The repo is not evaluating candidates against a fantasy perfect platform.
It is evaluating them against a harsher anti-benchmark:

- the request still only works if it lands on the unofficial right node
- the operator still privately remembers where the service really lives
- fallback still looks real only while the preferred backend is healthy
- stateful systems still gain a nicer edge story without gaining trustworthy
  authority or replication truth

If a candidate mostly preserves those conditions, then even a very respectable
platform is still a weak answer in this repo's terms.

## What the current worktree already proves before any platform promotion

Before comparing option families, the repo already proves several important
things.

### Strong proof: Compose is still the human contract

The root runtime is still the root
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml),
not Nomad jobs, not Kubernetes manifests, and not an OpenSVC-owned runtime.

That matters because every future orchestration move is being judged against a
Compose-first baseline, not a clean-room abstraction contest.

The thing being protected is not "Docker branding."
It is the readability and directness Compose still gives before distribution
becomes real.

Compose matters here because it is still one of the last places the operator
can read the system directly.
Any promotion away from that surface has to justify the epistemic loss, not
just the operational gain.

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

It is about request-time correctness, not just component plurality.

That distinction is easy to blur and critical to preserve.

Many systems can claim redundancy because more than one machine is involved.
Far fewer can honestly claim that the meaning of the request survives when the
wrong machine is the first one to receive it.

### Strong proof: the repo no longer trusts raw Compose to express the whole truth

The recurring appearances of:

- `services.yaml`
- failover-agent planning
- `docker-gen-failover`
- OpenSVC ingress experiments
- Constellation cluster-state work
- CUE and `x-cue` ideas

all point at the same conclusion:

Compose is still trusted as the authoring baseline,
but not as the complete coordination model.

That is one of the most important truths in the entire repository.

This is the sentence many pages still fail to honor strongly enough:

the repo is Compose-first in authoring, but already post-Compose in the
problems it is trying to solve.

That is why naive "stay with Compose" answers and premature "adopt the full
orchestrator" answers can both miss the point in opposite directions.

That sentence should also be read against the user's deeper frustration:

the platform market keeps offering cleaner labels much faster than it offers a
truthful middle contract.

So when this repo keeps searching for a missing layer, it is not being vague.
It is trying to recover a category the ecosystem keeps flattening.

## What the orchestration search is really trying to buy

Different pages name different candidate tools.
The underlying desired purchases are much more stable.

### Purchase 1: placement truth

The repo needs a reliable answer to:

- what runs where
- what can move
- what should answer this request
- what should happen if the receiving node is not the host

Without that, the operator remains the hidden scheduler.

### Purchase 2: convergence truth

The repo needs a reliable answer to:

- which config and secret material exists on which node
- who notices drift
- who regenerates routing
- who repairs obvious mismatch

Without convergence truth, distribution just means copying the same confusion
to more machines.

### Purchase 3: ingress truth

The repo needs a reliable answer to:

- whether any node can safely receive traffic
- whether routing survives placement changes
- whether routes survive local backend failure
- whether auth and middleware stay coherent across handoff

This is much stricter than "the proxy is up."

### Purchase 4: failover truth

The repo needs a reliable answer to:

- what the next backend is
- who chooses it
- whether the recovery path survives the local failure
- whether the platform converges after the handoff rather than silently
  degrading

### Purchase 5: lower operator reconstruction tax

The user's frustration is not only about outages.
It is about too many invisible truths still living in their head:

- current placement
- which hostnames are global versus node-scoped
- which peers are safe
- which routes are real versus aspirational
- which stateful systems are merely reachable versus actually resilient

The orchestration research is trying to buy back that mental rent.

## What the research should make impossible to miss

The most important finding is not "there are several plausible futures."
It is:

- several futures exist
- but they are not equally different
- many of them only become real options once they relocate a genuinely
  different hidden burden

That is why this page has to be stricter than a normal comparison.
The user is not short on paths that differ in packaging.
They are short on paths that differ in where the truth actually lives on a bad
day.

## Option families, but in repo terms

The right question is not:

> what does each product do?

The right question is:

> which missing truth would each family own, and what new burden would that
> import?

## Family 1: Compose plus coordination glue

Typical repo-shaped examples:

- `services.yaml`
- sync-agent concepts
- generated Traefik file-provider config
- DDNS updates
- health-aware failover helpers

### Why this family keeps returning

Because it promises the closest thing to the user's ideal:

- keep Compose readable
- keep manual placement where it still pays off
- add only the missing truths that raw Docker does not natively provide

### What this family is trying to buy

- explicit placement truth
- narrow convergence logic
- lighter routing generation
- wrong-node request preservation without total platform capture

### What this family risks

- helper sprawl
- bespoke control-plane complexity
- invisible coupling
- a system that calls itself "still just Compose" after it has quietly grown
  into a significant orchestration layer

This family is attractive because it resists premature worldview tax.
It is dangerous because it can accidentally rebuild a scheduler in fragments.

### The real test for this family

This family only stays honest while it remains narrower than the burden it is
trying to remove.

The moment it owns all of these at once:

- placement truth
- convergence truth
- peer eligibility
- route persistence
- repair decisions

the repo has to stop flattering itself with "still just Compose" language and
compare the helper surface honestly against stronger named control planes.

## Family 2: HA substrate plus Compose-facing runtime

OpenSVC and adjacent substrate ideas live here.

This family says:

- keep higher-level service authoring closer to Compose
- let a stronger HA substrate own membership, placement-adjacent truth, or
  ingress/failover coordination

### What this family is trying to buy

- stronger failure-handling semantics
- more explicit cluster truth
- a middle path between raw glue and full Kubernetes-style controller culture

### What this family risks

- another control plane to learn
- more abstraction than the repo may want
- a new dependency that still has to prove it reduced operator burden overall

This family matters because it might be the real missing middle for some
domains, especially where helper sprawl starts approaching control-plane size.

### The real test for this family

This family is only a real answer when the repo's dominant pain is genuinely
infra-grade:

- sacred ingress
- sacred identity
- sacred first-hop behavior

If the real pain remains wrong-node application request preservation for
ordinary services, then this family can still be valuable without being the
whole answer.

## Family 3: scheduler promotion

Nomad, k3s, and Kubernetes live here.

This family says:

- stop reconstructing distributed truth by hand
- let a scheduler or cluster API own more of placement, health, and service
  identity

### What this family is trying to buy

- more native distributed truth
- stronger service identity and placement semantics
- less bespoke glue

### What this family costs

- major worldview tax
- higher platform complexity
- diminished Compose-level directness
- a more explicit control plane with more ritual and abstraction

The repo's stance is not "never."
It is "prove which hidden burden is large enough to deserve this tax."

### The real test for this family

This family becomes real once manual placement plus helper truth is clearly
the bigger burden than scheduler worldview cost.

Until then, this family remains vulnerable to the same problem the user keeps
rejecting:

- a larger platform gets introduced
- the docs feel calmer
- but the exact missing truth that earned the promotion is still not named
  concretely enough

That is not maturity.
It is a cleaner version of the same ambiguity.
It is:

> not before the need is proven by the failure of smaller honest answers

That distinction has to stay explicit.

It also has to stay granular.

"Promotion" here should not be narrated as one giant irreversible moral event.
Different domains may earn stronger ownership at different times:

- ingress and wrong-node preservation
- placement and restart authority
- convergence and cluster truth
- stateful election and replication

If one page starts talking as if a single tool promotion obviously closes all
four domains at once, that page is overselling the architecture.

## The real evaluation question

Every platform family should be judged by the same hard question:

> did it remove hidden operator burden and false resilience claims without
> replacing them with a larger, equally frustrating fiction?

That is the test for:

- helper meshes
- OpenSVC-style substrates
- Nomad
- k3s
- Kubernetes

The question is never just:

> can it route traffic?

It is:

> which missing truth does it own, how inspectable is that truth, what tax
> does it import, and does that tax finally pay down the wrong-node and
> backend-loss lies that Compose alone cannot fix?

## The purchase map the repo should keep using

This is the compressed decision table underneath the whole research effort.

| Missing truth | Compose plus helpers | HA substrate | Scheduler promotion |
| --- | --- | --- | --- |
| Placement truth | Possible, but risks folklore unless the registry becomes real and consumed | Partial, depending on how much service identity the substrate owns | Strongest native answer |
| Convergence truth | Possible, but easy to turn into invisible custom glue | Better for infra-level convergence than app-level coherence by default | Strong if the scheduler truly owns the runtime contract |
| Wrong-node request preservation | Possible, but only if routing, peer selection, and policy continuity all become real | Promising for narrower infra-grade ingress/failover domains | Strongest broad answer once the platform is fully adopted |
| Route persistence under backend loss | Fragile unless the helper layer is very disciplined | More naturally aligned with explicit failover ownership | Often stronger, but at the cost of a full platform worldview |
| Stateful correctness | Weak by itself | Mixed; infra-grade HA helps, but state still needs service-specific truth | Mixed; better placement and lifecycle, but state semantics still remain service-specific |
| Operator readability | Strongest | Medium | Weakest by default |
| Worldview tax | Lowest at first, but can rise invisibly | Medium | Highest |

This table is not a winner declaration.
It is a reminder that different layers buy different truths at different
prices.

It is also not a promise that the repo will eventually converge on exactly one
family for every problem.

One of the user's strongest recurring complaints is that option spaces get
flattened too early.
This repo should preserve the possibility that:

- a thinner layer may own some HTTP placement and forwarding truths
- a stronger HA substrate may own narrower infra failover truths
- stateful services may still need service-specific authority models that do
  not collapse neatly into the same answer

That is not indecision.
It is a refusal to let cleaner architecture rhetoric erase real differences in
what each burden actually is.

## The orchestration search is really protecting something

This research is not protecting a brand preference.
It is protecting the user's actual dream:

- real options instead of fake ones
- real readability instead of hidden sacred knowledge
- real request preservation instead of DNS theater
- real honesty around stateful systems instead of flattering liveness language
- real justification before a giant control plane becomes the default answer

That is why this page exists.

The user is trying to force a more honest infrastructure category into view:

a category where distributed Docker does not immediately collapse into either:

- brittle glue
- or total orchestrator capture

This page should keep that pressure explicit.
