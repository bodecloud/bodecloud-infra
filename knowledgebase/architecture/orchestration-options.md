# Orchestration Options: Which Extra Layer Has Actually Earned the Right to Exist?

For the evidence underneath this page, start with
[`../research/orchestrator-tradeoffs-evidence.md`](../research/orchestrator-tradeoffs-evidence.md)
and [`../research/orchestration-research-2026.md`](../research/orchestration-research-2026.md).

This page exists to stop one lazy question from taking over the repo:

> which orchestrator is best?

That is not the real decision surface here.

The real question is:

> which extra layer of control has actually earned the right to exist because
> it removes a real hidden human SPOF, a real wrong-node failure, or a real
> convergence failure that the current Compose-first stack cannot solve
> honestly?

That is the decision surface that matters.

If this page ever turns into a normal "platform comparison" page, it has
already failed.

It fails in a second, more subtle way when it starts sounding like the main
work is simply picking between respectable tool families.

That is smaller than the real problem.

The real problem is deciding which layer, if any, has earned the right to own
the truths the operator is still privately carrying:

- where the service really lives now
- which peer is actually eligible now
- whether the fallback route still exists now
- whether the request keeps its meaning after handoff

The user is not short on platform names.
The user is short on options that still feel honest after wrong-node entry,
backend loss, and hidden operator glue become real.

That difference is the whole reason this page exists.

If it becomes a neat feature table, it will have already degraded into the same
kind of answer the user keeps rejecting elsewhere.

This page therefore has to be harsher than a normal "architecture options"
document.
It has to stop the repo from confusing:

- more categories
- more tool families
- more respectable platform names

with:

- a genuinely different answer to the hidden-truth problem

If a candidate layer still leaves the operator privately reconstructing the
same bad-day answer, then from the user's point of view the option family has
mostly changed costumes, not outcomes.

## What this page is and is not allowed to prove

This page is authoritative about:

- how orchestration options should be judged in this repo
- which burden a promoted layer would need to own to justify its existence
- why respectable platform names are weaker than domain-specific truth ownership

This page is not authoritative about:

- whether one orchestration path has already won globally
- whether the current runtime already demonstrates the promoted behavior
- whether broader ecosystem prestige should override repo-specific evidence

This page is a promotion filter, not a final platform verdict.

## Priority decision stack for this page

When this page evaluates whether some extra layer has earned the right to
exist, it should route the question through this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`../architecture/current-compose-runtime.md`](current-compose-runtime.md)
3. [`../architecture/problem-and-goals.md`](problem-and-goals.md)
4. [`../research/orchestrator-tradeoffs-evidence.md`](../research/orchestrator-tradeoffs-evidence.md)
5. [`../operations/proof-matrix-and-drills.md`](../operations/proof-matrix-and-drills.md)

That order matters.

If this page starts from product families, it collapses into market prose.
If it starts from runtime alone, it under-reconstructs the dream.
If it starts from plans alone, it over-promotes futures into present gravity.

This page only stays honest if the dream, the live baseline, the hidden wound,
the candidate evidence, and the proof ceiling are all present at once.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "Compose-first is still the default" | current strategy stance | this page, `.github/copilot-instructions.md`, current runtime pages | that Compose already solves the missing truths |
| "a stronger layer has not yet earned itself" | promotion-threshold judgment | this page plus orchestrator evidence and proof pages | that stronger layers are forbidden in principle |
| "this candidate owns a specific hidden burden better" | option-family evaluation | this page plus candidate evidence | that the candidate therefore solves the whole dream |
| "prestige or maturity language is too weak" | anti-theater judgment | this page, archive pressure, proof matrix | that professional tooling has no place here |

If a sentence starts sounding like "pick the most mature orchestrator," it has
already left this page's decision surface.

## Automatic disqualifiers for a not-yet-earned option

A candidate option has not yet earned default promotion here if it mainly does
one or more of these:

- changes the control surface without naming which hidden burden moved
- improves deployment prestige while leaving wrong-node meaning weak
- improves controller power while leaving stateful authority socially held
- reduces local toil while preserving private topology reconstruction on the
  bad day
- sounds more adult mainly because the worldview got larger

That does not make the candidate worthless.
It means the candidate has not yet answered the repo's real benchmark strongly
enough to deserve default status.

## The shortest honest answer

The current default stance is still:

> stay Compose-first, close the missing truth layers as narrowly as possible,
> and delay whole-stack orchestrator promotion until a specific domain clearly
> proves it needs a stronger control plane

That is not hesitation.
It is the most evidence-aligned answer available right now.

It is also the least insulting answer available right now.

Too much infrastructure advice effectively says:

- accept brittle manual glue
- or surrender to a giant control plane
- and stop pretending there should have been anything meaningful in between

This repo exists because the user does not accept that as the only adult
decision surface.

The repo is not starting from zero.
It already has:

- a real root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- a broad set of included fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- real ingress, auth, observability, and state-bearing surfaces
- an explicit no-heavy-orchestrator default in
  [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- a real multi-node dream built around any-node entry, local-first serve, and
  peer-forward fallback

What it still does not have is the missing middle layer that makes those
sentences trustworthy under failure.

That phrase matters.
The docs are not trying to preserve optionality for aesthetic reasons.
They are trying to preserve a search space where the answer is not always
forced to collapse into:

- underpowered static glue
- or overpowered worldview replacement

That middle search space is not theoretical.
It is the main object the repo is trying to defend.

That middle search space is also easy to erase with mature-sounding docs.
If this page starts reading like the repo is simply choosing between
respectable platform families, it has already shrunk the dream into a smaller,
more ordinary decision than the one the user is actually forcing.

That is also why this page must keep separating real options from relabeled
non-options.

For this repo, a platform path is not meaningfully different just because the
tool name changes.
It becomes meaningfully different only if it owns a different hidden burden
than the current path and makes that burden less dependent on private operator
reconstruction.

That standard needs to stay severe.
Many platform comparisons fail this repo precisely because they over-credit
surface differences and under-credit continuity of burden.
If the same sacred truths survive under a new control plane, then the new
control plane has not really answered the user's complaint yet.

That last clause should stay brutal.

If a new layer mostly improves:

- naming
- legitimacy
- ecosystem reassurance
- post-hoc explanation quality

while still leaving the same sacred truths privately carried, then it has not
become a real new option for this repo.

## Strongest honest current answer

If someone asks, "What is the shortest honest answer about orchestration
options right now?" the answer is:

> No broader orchestration family has yet earned default promotion for the repo
> as a whole; the current evidence still supports staying Compose-first and
> forcing any stronger layer to justify itself by owning a named hidden burden
> better than the current stack, not by sounding more adult or more complete.

## What the user is actually trying to avoid

The orchestration story makes no sense unless the user’s refusals stay visible.

The user is not merely trying to avoid complexity.
They are trying to avoid bad complexity.

More specifically, they are trying to avoid complexity that steals legibility
without actually removing the invisible burden that made the old option fail.

That is the key anti-benchmark for every promoted layer:

- did it really remove the burden
- or did it just hide it behind a cleaner abstraction boundary

### Refusal 1: raw Compose pain disguised as freedom

Compose still feels good when:

- the operator knows what runs where
- requests land on the right node
- the stack is small enough that remembered facts are still manageable

It stops feeling good when:

- wrong-node requests become normal
- placement truth drifts
- failover paths depend on private operator memory
- service semantics start living in fragile static glue

The user does not want docs that romanticize that pain.

### Refusal 2: heavyweight orchestrator ideology before the pain is named precisely

The user is tired of being told:

- just use Kubernetes
- just use Nomad
- just use Swarm
- just use a platform

before the docs can answer a simpler question:

- which exact missing truth is the new layer buying?

If that question is not answered first, orchestrator advice becomes religion
instead of engineering.

This page should therefore be read as a filter against ideology, not as a
shopping guide.

It should also be read as a filter against theatrical differentiation.

If two proposed futures both leave:

- sacred-node ingress truth
- socially reconstructed placement truth
- fake HA closure around stateful surfaces

then they are much closer to the same answer than to different answers, no
matter how different the product pages sound.

That is one of the repo's deepest themes.
The problem is not that Kubernetes, Nomad, OpenSVC, or helpers are forbidden.
The problem is that advice too often arrives as ideology before it arrives as a
domain-specific explanation of what truth gets better.

It also often arrives as comfort before it arrives as ownership.

The page should keep asking:

- what exact truth would this layer own directly
- what truth would still remain reconstructed by humans
- what humiliation would still survive on the bad day even after promotion

### Refusal 3: fake HA closure

The user is also tired of success stories that stop at:

- multiple DNS records
- a healthy proxy
- a route that exists while the primary is healthy
- a green healthcheck

None of those prove:

- wrong-node request preservation
- peer eligibility truth
- middleware continuity
- stateful correctness

So this page has to remain hostile to recommendation styles that sound good
because they imply closure faster than the evidence earns closure.
That is another way of saying:

- calmness is not proof
- maturity vocabulary is not proof
- controller presence is not proof

That is why this page has to stay organized around pain classes, not product
features.

## The only promotion rule that matters

The repo does not need a popularity contest between platforms.
It needs a harsh promotion rule:

1. identify the dominant failure class
2. identify whether the current Compose-plus-helper path can solve it honestly
3. identify whether a stronger layer removes that failure class or merely
   renames it
4. promote only the domains whose pain clearly earns the new tax

If a page cannot say what tax is being paid and what failure it removes, it is
still too shallow.

That is the standard every later orchestration recommendation should inherit.

This standard is deliberately harsher than normal tooling guidance because the
user has already seen too many "solutions" that merely rename the burden.

It is also harsher because the docs themselves are part of the failure mode.
A smoother recommendation can be actively harmful here if it causes the reader
to feel the choice got simpler while the bad-day request path is still socially
reconstructed.

## The problem classes that actually decide the answer

The important question is not "what features does this platform have?"
The important question is "which missing truth would this platform own?"

### 1. Authoring and inspectability

Questions:

- how are services defined?
- how much hidden machinery is required to understand what will happen?
- how easy is local iteration and debugging?

This is where Compose is still strongest.

The repo still gets real value from:

- readable YAML
- domain-based fragments
- Docker-native workflows
- direct file inspection

Any promoted layer has to justify what it takes away here.

Because once readability is spent, the replacement has to do more than look
modern.
It has to give back enough truth to justify the loss.

### 2. Placement and convergence truth

Questions:

- where does a service live right now?
- how do nodes agree on that truth?
- how do env, secret, image, and routing assumptions converge?
- what happens when a node disappears or drifts?

This is where raw Compose starts to require either:

- stronger helpers
- a real registry
- or promotion into a stronger control layer

But this page should keep one warning attached to all three:
if the resulting truth is still fragmented enough that a knowledgeable operator
has to mentally fuse several partial signals into "what is actually true right
now?", then the layer has still not earned full trust.

### 3. Traffic and request continuity

Questions:

- can any healthy node receive traffic and still preserve success?
- can the receiving node decide between local serve and peer handoff?
- can the route survive the failure that made fallback necessary?
- do middleware, auth, and headers remain coherent after handoff?

This is where the repo’s real pain sits.

That point deserves to stay emotionally explicit:

the repo is not primarily suffering from a lack of orchestrator features.
It is suffering from the absence of a trustworthy answer to "what should the
wrong node do right now without me privately carrying the topology?"

That sentence needs one sharper continuation:

if a candidate layer cannot answer that in system-owned terms, then however
powerful it is, it still has not earned promotion for the main pain this repo
is organized around.

That question is stronger than "which scheduler is better."
It is the most useful compression of the whole repo.

### 4. Infra-grade failover

Questions:

- can a narrow critical service relocate or survive node loss explicitly?
- does the repo need boring infra-style HA primitives instead of a universal
  application platform?

This is where OpenSVC-style promotion becomes attractive.

### 5. Stateful correctness

Questions:

- who owns writes?
- how is promotion coordinated?
- how do clients rediscover authority?
- what storage semantics survive node loss?

This is the class the repo most strongly refuses to blur into generic
"orchestrator solved it" language.

No scheduler alone closes this.

## Baseline facts before promoting anything

Before talking about platform families, the repo already proves several things.

### Fact 1: the live runtime is still concretely Compose-first

The current operator surface is still rooted in:

- root `docker-compose.yml`
- the include graph
- Docker-native labels, networks, and service definitions

That means every promotion decision is being judged against a Compose-first
baseline, not a blank sheet.

### Fact 2: the stack is already broad enough that invisible operator truth is now an architectural problem

The active runtime spans:

- ingress and edge policy
- private coordination surfaces
- observability
- heterogeneous app groups
- state-bearing services
- network and egress experimentation

That means the missing middle layer is not theoretical.
The stack is already large enough that remembered facts and hand-managed
convergence have become part of the architecture.

That is why the option question is no longer theoretical.
The repo is already paying complexity tax.
The unresolved question is which complexity is actually buying truth and which
complexity is still just theater.

Another way to say the same thing:

which complexity would still leave a future contributor able to write a smooth
page that only works because they already know which node is the real one?

If the answer is "quite a lot of it," then the promoted layer is not yet
buying enough truth.

### Fact 3: the repo still does not broadly prove the key truths a promoted layer would need to own

The current stack still does not broadly prove:

- live tracked placement truth
- auditable convergence truth
- durable route persistence under failure
- broad wrong-node request success
- stateful anti-SPOF correctness

That is why the first answer is still "name the missing truth" rather than
"pick a winner."

## Option family 1: Compose-first plus stronger helpers

This is not the absence of a strategy.
It is the current strategy.

It is also the category the user most wants to be real.
If this family can be made honest enough, it preserves the most of what still
feels good about Compose while buying back the truths that raw Compose keeps
outsourcing to the operator.

Representative repo-shaped ideas include:

- `services.yaml`
- sync-agent and failover-agent directions
- generated Traefik file-provider config
- Constellation-style convergence helpers
- DDNS and node-aware route generation

### What this family preserves

- readable service definitions
- Docker-native workflow
- direct local iteration
- explicit operator placement control
- relatively low worldview tax

### What this family is genuinely buying

If done well, it buys:

- explicit placement truth without full scheduler capture
- explicit route generation from current state
- narrower convergence logic
- a path toward wrong-node recovery without abandoning Compose

The key phrase is "if done well."
This page should not let that phrase evaporate into optimism.
The archive and runtime both show how easy it is for helper-heavy systems to
sound dynamic while still asking the operator to privately arbitrate between
stale, partial, or socially remembered truths.

That "if done well" qualifier matters a lot.
The archive is full of examples where helper-heavy approaches sound dynamic but
still leave the user editing, remembering, or inferring the most important
facts by hand.

### What this family is weak at unless it grows real teeth

- placement truth can remain half-social
- convergence truth can remain half-social
- route persistence can fail exactly when needed
- helper logic can become a hidden control plane while the docs still call the
  system just Compose

### When this family remains the right path

Stay here when:

- the main pain is still coordination truth more than scheduler power
- direct file-level control is still buying real value
- the helper layer is still explainable and auditable
- the repo can still show which specific truth surface each helper owns

### The real warning

Compose-plus-helpers stops being the honest answer when:

- placement truth remains mostly manual
- convergence truth remains mostly social
- wrong-node success depends on tribal knowledge
- route generation remains brittle under real failure
- the helper layer effectively becomes a custom orchestrator in disguise

At that point the repo must either formalize that control plane honestly or
promote the right domain into something stronger.

That is one of the harshest but most necessary conclusions in the whole
knowledgebase:

if the helper mesh owns enough truth, then pretending the system is still "just
Compose" becomes as dishonest as pretending Kubernetes would have been free.

That is one of the most important anti-benchmarks in the whole knowledgebase.
It prevents the repo from treating helper accumulation as automatically safer
simply because it happened in smaller steps.

## Option family 2: OpenSVC or other infra-grade HA tooling

This family matters because some of the repo’s hardest pains are
infra-failover-shaped, not scheduler-shaped.

That distinction deserves more weight than generic orchestration discussions
usually give it.
The user is not simply asking for "better clustering."
They are asking for the exact layer that is still sacred to stop being sacred,
and sometimes that layer is much narrower than an application platform.

That distinction is one of the least appreciated parts of the user's demand.
Sometimes the pain is not "schedule everything better."
Sometimes the pain is "this narrow critical surface must actually survive a node
loss without turning into mythology."

### What this family preserves

- narrow domain promotion
- explicit failover semantics
- a direct mapping between node loss and recovery behavior
- less pressure to turn every service into an application-platform problem

### What this family is genuinely good at

- ingress survivability
- explicit service takeover
- targeted anti-SPOF promotion
- operator-readable HA behavior for narrow critical domains

### What it is not

- a universal application platform
- a replacement for all Compose authoring
- a shortcut around data-topology truth
- proof that stateful correctness is solved just because service takeover is
  cleaner

### When this family becomes the best fit

Promote here when:

- the dominant pain is explicit node or service failover for a narrow critical
  domain
- the affected surface is small enough to isolate
- the key question is "what happens when this service or node dies?" rather
  than "how do I schedule everything?"

### Why it fits this repo better than generic docs usually admit

The archive shows the user is more hostile to fake HA than to extra tooling.

That makes infra-grade HA promotion attractive when it solves one precise
failure class directly instead of demanding whole-stack worldview surrender.

This is exactly why OpenSVC-style thinking remains attractive here even though
generic docs often skip past it.
It offers a way to spend complexity on one painful truth instead of importing a
whole political system for the stack.

That is why this family should be judged by discipline, not breadth.
If it starts being narrated as a near-universal resilience answer, the docs are
sliding back into the same fake-closure pattern the user is tired of.

## Option family 3: Nomad

Nomad keeps returning because it is the clearest scheduler-shaped answer that
does not instantly impose full Kubernetes worldview capture.

That still does not make it self-justifying.
It makes it a serious candidate once the repo can honestly say the main pain is
now scheduler pain rather than missing-middle-layer pain.

### What this family preserves better than Kubernetes

- a smaller conceptual footprint
- clearer scheduler-centric reasoning
- less ecosystem sprawl by default

### What this family is genuinely good at

- placement
- relocation
- restart coordination
- turning current-state scheduling into something more operational than social

### What it may buy for this repo

Nomad is strongest when the real pain has become:

- service placement
- relocation after host loss
- restart orchestration
- job-oriented convergence

rather than:

- raw ingress failover semantics by themselves

### What it still does not buy automatically

- middleware continuity under peer handoff
- stateful correctness
- a free answer to all ingress questions
- freedom from storage and topology design

### When this family earns promotion

Promote here when:

- scheduler pain is now clearer than helper growth pain
- the repo is ready to let a scheduler own placement truth
- the operator accepts a job-and-cluster worldview for the promoted domain

Nomad is a meaningful candidate when the helper layer is starting to look like
a scheduler anyway.

The key phrase is "look like a scheduler anyway."
The repo should not promote Nomad because it is fashionable.
It should promote Nomad if it is already paying scheduler-like tax in a more
fragile and less self-aware form.

That sentence is the real gate.
Nomad is not here because it is an elegant middle ground in the abstract.
It is here because a repo can drift into scheduler-shaped suffering long before
it admits that scheduler ownership may now be cheaper than pretending the
helper mesh is still "small."

## Option family 4: k3s or Kubernetes

This family should only be judged honestly, not reflexively.

That means resisting both:

- "Kubernetes solves real infrastructure, so eventually just admit it"
- "Kubernetes is evil, so never let it into the conversation"

Neither reflex is serious enough for this repo.

This page should stay resistant to both comfort reflexes:

- "eventually just admit Kubernetes"
- "never let Kubernetes near anything readable"

Both are ways of skipping the harder question about which truth the layer would
own more honestly than the current stack.

### What this family is genuinely good at

- scheduler-backed placement
- broad ecosystem integration
- mature cluster primitives
- richer controllers for workloads that truly need them

### What this family costs

- a large worldview shift
- storage and networking tax
- controller complexity
- a bigger abstraction gap from current Compose authoring

### When this family actually earns promotion

Promote here only when:

- the promoted domain truly needs scheduler-backed placement and controller
  machinery
- the helper path is no longer staying thin
- the repo has crossed from "missing middle layer" pain into "we are already
  rebuilding a cluster platform badly" pain

### What this family should not be allowed to do in the docs

It should not be presented as:

- automatically solving wrong-node success
- automatically solving ingress semantics
- automatically solving stateful correctness
- automatically being the mature answer without naming which failure class it
  is removing

Kubernetes only earns itself here if the docs can answer a question more
specific than "it is the standard":

which truth did the repo fail to own any other way, and why is the controller
tax now cheaper than the helper tax?

That cost comparison is the whole point.
If the docs cannot state it concretely, then "Kubernetes is mature" is just a
prestige argument, not an engineering argument.

## Option family 5: Swarm

Swarm matters mostly because it is the obvious Docker-native temptation.

### What it superficially offers

- familiar Docker shape
- scheduler flavor without full Kubernetes tax
- easier mental migration from Compose

### Why it still does not obviously win here

The repo’s real pressure is not "I wish Docker had a cluster mode."
It is:

- wrong-node request preservation
- truthful failover
- convergence truth
- stateful honesty

That is why Swarm keeps failing to become the obvious answer.
It is near the current worldview, but nearness is not the same as removing the
specific humiliations that made the current worldview insufficient.

Swarm does not automatically answer those well enough to earn default victory,
especially when the user is already wary of fake closure.

That is why Swarm remains more temptation than conclusion.
It is close enough to Docker to sound like the obvious compromise, but the
repo's real question is not "what feels Docker-shaped?"
It is "what actually removes the fake-option burden?"

## The harsh cross-check for every option

Before any promoted layer is described as the answer, ask:

1. Which exact missing truth will this layer own?
2. How will the operator know that truth is current?
3. What tax does the layer import?
4. What failure class does it remove that the current layer cannot remove
   honestly?
5. Does it preserve the user’s actual demand surface?

That demand surface is still:

- any node can receive traffic
- local services stay local
- wrong-node requests still succeed
- middleware and auth remain coherent
- stateful systems are only described as resilient when they truly are
- the operator is not forced to carry invisible cluster truth privately

If a candidate layer does not improve those conditions, it has not earned
itself.

This is the repo's version of a proof gate.
Tooling should not be promoted because it sounds more adult.
It should be promoted because one previously human-held truth becomes
materially more visible, current, and trustworthy.

This page should keep those three words together:

- visible
- current
- trustworthy

Many options can improve one or two of them while still failing the third.
The user's complaint lives exactly in that gap.

## Bottom line

This repo should not promote an orchestrator because:

- it is popular
- it is lighter
- it is Docker-native
- it is powerful

It should promote a layer only when it can say, concretely:

> this layer owns this missing truth, removes this hidden human SPOF, and pays
> for itself more honestly than staying Compose-first would

Right now the most honest answer remains:

- keep Compose as the live baseline
- keep making the missing truth layers explicit
- promote only the domains whose failure class clearly earns a stronger layer

That is not indecision.
It is the only non-theatrical orchestration policy the repo currently supports.

And if that answer ever starts sounding timid, remember what the repo is
reacting against:

- fake closure
- sacred-node cosplay
- hidden control-plane labor
- platform adoption that cannot explain what exact wound it is healing
