# Design Tensions and Contradictions

This page exists because the clean version of this project is a lie.

The clean version says:

- Compose is simple
- DNS can point at multiple nodes
- Traefik can route
- some helper glue can fill in the gaps
- therefore the platform can become multi-node and anti-SPOF without much pain

That is exactly the kind of story the user is tired of.

It is not just technically incomplete.
It is emotionally false to the actual experience that produced this repo.

The real experience is:

- the ecosystem keeps presenting "options"
- those options sound empowering at the diagram level
- the hard truths keep getting pushed back into operator memory
- then the same systems act surprised when the user says the options still do
  not feel real

The real pressure behind `bolabaden-infra` is not "please describe a nice
distributed architecture."

It is:

> why does everything feel empowering and legible while the stack is local, but
> the second requests can land on the wrong node, services can live on
> different boxes, and failure has to be absorbed instead of merely noticed,
> the option space collapses into fake HA language, hidden operator burden, and
> heavyweight orchestrators that demand surrender before they prove they are
> worth it?

This page is here to keep that pressure alive.

It names the contradictions the repo is actually wrestling with so the rest of
the docs cannot quietly "resolve" them by smoothing them over.

If this page ever becomes comfortable, it has probably failed.

That discomfort is not a tone choice.
It is part of the retrieval contract of the whole knowledgebase.

The user is not asking for architecture that sounds reasonable.
They are asking for architecture writing that keeps naming where the option
space still feels fraudulent.

## What this page is and is not allowed to prove

This page is authoritative about:

- which contradictions the repo is consciously carrying
- why smoothing those contradictions would misdescribe the project
- where dream, runtime, and control-plane growth are still in conflict

This page is not authoritative about:

- whether one side of the contradiction has already won
- whether the runtime has already resolved the tensions
- whether a future layer is already justified just because the tension is clear

This page names the pressure honestly.
It does not certify that the pressure has been relieved.

## Quick claim router

| If the sentence is really claiming... | Primary class | Strongest anchors | It still must not imply... |
| --- | --- | --- | --- |
| "this contradiction is central to the repo" | synthesis over intent, runtime, and archive pressure | `.github/copilot-instructions.md`, current runtime pages, archive-pressure pages | that the contradiction is already resolved |
| "Compose readability and distributed truth are in tension" | current runtime + intent | Compose-first pages, current runtime pages, intent surfaces | that Compose-first has already passed the bad-day test |
| "helper growth may already be control-plane growth" | contradiction analysis + planning pressure | master-plan, orchestration research, evidence pages | that helper growth is therefore already wrong or already sufficient |
| "local-first and resilient fallback want different things" | routing and failure analysis | ingress evidence pages, planning docs, this page | that either half has already been fully earned |

If a contradiction is described in a way that sounds settled, the page has
lost the reason it exists.

## The contradiction at the center of the whole repo

The dream sounds small when compressed:

> any surviving public node should be able to receive the request, decide
> whether the target is local, preserve the request if it is not, and do all of
> that without forcing the operator to accept a giant scheduler as the default
> answer

That sentence sounds modest.
It is not modest.

It sounds modest mostly because most documentation cultures are trained to talk
about this kind of demand as if it were just one more routing problem.

It is not just a routing problem.
It is a demand that the system itself stop behaving stupidly once topology,
failure, and service placement become real.

It requires at least five different kinds of truth to coexist:

- ingress truth: which nodes can honestly receive the first hop
- placement truth: where the service actually lives right now
- routing truth: what path should be taken when the service is not local
- eligibility truth: which remote target is safe, not merely reachable
- semantic truth: whether the remote answer still means the same thing

Single-node Docker largely lets the operator ignore those distinctions.
Multi-node request preservation does not.

That is why the user's frustration feels bigger than "I need better routing."

The wrong-node request is the moment the system has to reveal whether it
actually contains those truths or whether the operator was silently supplying
them all along.

That is why so many "lightweight" answers stop feeling lightweight the second
the user asks the system, rather than the operator, to answer the wrong-node
question.

That is the contradiction the repo keeps circling:

> the project wants the legibility of local Docker and the survivability of a
> distributed platform, but it does not want to accept a heavyweight control
> plane until that control plane has clearly earned the right to exist

Everything else on this page is just a sharper version of that one conflict.

That is why this page should feel less like a balanced tradeoff memo and more
like a refusal to let the docs buy calmness by deleting the wound.

## What this page is for

This page is not trying to "balance tradeoffs."
It is here to stop three kinds of drift:

- dream drift: the docs collapse the user's real frustration into generic
  architecture language
- proof drift: intent gets narrated as live behavior
- control-plane drift: the repo quietly rebuilds orchestration in fragments
  while pretending nothing substantial was added

If those drifts are not resisted, the repo becomes impossible to read honestly.

It also becomes indistinguishable from the rest of the infrastructure discourse
the user is already frustrated with: lots of components, lots of buzzwords, and
no honest accounting of where the real burden still lives.

This is why contradiction is not a documentation flaw here.
It is evidence.

The project is in motion between:

- strong authoring truth
- strong dream clarity
- uneven runtime proof
- multiple candidate middle layers

If the docs smooth that into one comfortable architecture posture, they erase
the very thing the user most needs help thinking through.

## Tension 1: Compose readability vs scheduler-grade expectations

## What the repo wants

- the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  and `compose/` fragments remain the primary authoring surface
- service definitions stay human-readable
- an operator can still inspect the stack without decoding a full cluster API
- manual service placement remains allowed where it still pays for itself

## What the repo also wants

- any public node can receive the first request
- the node can tell whether the target is local
- if not local, the request can be preserved through peer forwarding or
  equivalent recovery logic
- routes can survive the local failure that made recovery necessary

## Why the tension is real

Compose is good at expressing:

- local container intent
- networks
- configs
- secrets
- healthchecks
- dependency hints

Compose does not naturally own:

- distributed placement truth
- cluster convergence truth
- peer eligibility
- cross-node route persistence
- wrong-node request preservation

So the second the repo asks Compose to participate in those higher-order truths,
one of two things happens:

1. helper layers accumulate around Compose
2. a stronger control substrate gets promoted

That does not mean Compose-first is wrong.
It means Compose-first is only honest if the repo admits the missing work rather
than describing Compose as if it natively solved distribution.

This is one of the most important tensions to keep alive.

Compose is not being defended as a religion.
It is being defended as the last surface where the system still feels causally
legible.

The contradiction appears because the user wants to keep that legibility while
demanding behaviors that Compose alone does not close.

The user's frustration is aimed directly at this gap.
They are not asking for prettier YAML.
They are asking why "simple Docker" stops offering real options once
distribution matters.

That wording matters.
The complaint is not "Compose is imperfect."
The complaint is "the supposed alternatives keep removing legibility without
actually removing the hidden burden."

## Tension 2: "No orchestrator by default" vs "someone still has to know what runs where"

The most important intent surface in the repo is still
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It states the dream cleanly:

- no central orchestrator by default
- manual service placement is acceptable
- current-state truth is preferred over scheduler-declared desired state
- local-first serving is preferred
- peer-forward fallback is required when the target is not local

That is not vague.
It is specific.

But it quietly implies a question that never goes away:

> what tells a node, right now, where the service really lives?

The repo keeps reaching for smaller answers:

- `services.yaml`
- sync-agent
- failover-agent
- generated Traefik file-provider config
- OpenSVC-assisted inventory
- various control or convergence helpers

Those are not random ideas.
They are repeated attempts to avoid paying the full cost of a heavyweight
scheduler while still extracting one scheduler-grade truth:

> current placement

The present-tense contradiction is even sharper:

- the knowledgebase and instruction surfaces repeatedly depend on something
  like `services.yaml`
- the root worktree does not currently prove that a live tracked root
  `services.yaml` exists and is consumed by the priority runtime

So this is not only a philosophical tension.
It is an evidence tension:

the docs can currently describe the missing middle layer more clearly than the
runtime proves it exists.

That is exactly why documentation honesty has to stay strict.

There is a harsher way to say the same thing:

the repo already knows the shape of the missing answer more clearly than it has
implemented the answer itself.

That is survivable if the docs admit it.
It becomes dangerous the second the docs start narrating that missing layer as
if it were already settled runtime truth.

This is also the point where many projects start lying to themselves.
They accumulate enough architecture language to describe the missing middle
layer fluently, then mistake fluent description for implemented truth.

That is exactly the documentation trap this repo has to keep resisting.

The more clearly the missing layer can be described, the easier it becomes to
mistake clarity of desire for closeness of implementation.

This repo has to resist that move on purpose.

## Tension 3: local-first serving vs route persistence under failure

The request model the repo wants is simple to say:

```text
request hits any node
  service is local  -> serve locally
  service is remote -> preserve path and forward to healthy peer
```

The first half is intuitive.
The second half is where most systems start bluffing.

That bluff is one of the archive's main recurring enemies.

It often takes the form of:

- the route exists in config
- the proxy can technically target a peer
- the platform therefore has failover

The repo keeps insisting that those are different claims.

The archive and planning surfaces already record the key failure:

- the local backend dies
- route generation runs again
- the route disappears along with the local container

That is the exact opposite of what failover is supposed to mean.

So two separate truths have to stay visible:

- local-first is the preferred fast path
- fallback-path persistence is the required correctness path

If only the first is solved, the platform remains a fair-weather system.
It works while locality holds and collapses the moment locality breaks.

That is why "local-first" is not the prize.
It is only the fast path.

The real prize is preserving request meaning when locality breaks without
forcing the operator to privately restitch the platform in their head.

That is not a minor defect.
It is the main complaint.

The contradiction here is severe:

- local-first wants runtime truth to stay close to the receiving node
- resilient fallback wants routing truth to outlive the local backend

Any layer that cannot satisfy both is still a partial answer.

## Strongest honest current answer

If a reader asks, "What is the deepest contradiction this repo is living
inside?" the shortest defensible answer is:

> The repo wants the legibility of Compose and the survivability of a
> distributed platform, but it does not want to accept heavyweight worldview
> tax until that tax has proved it owns the exact missing truths better than a
> thinner layer could. That contradiction is still active; it has not been
> resolved by current runtime seriousness or by planning clarity alone.

Anything tidier than that is probably deleting the real pressure.

A design that only works while locality holds is not anti-SPOF.
It is locality with nicer branding.

That line is harsh on purpose because the entire archive keeps converging on the
same disappointment:

- impressive normal-case routing
- followed by disappearance of the exact path that failure made necessary

That is not graceful degradation.
It is recovery theater.

## Tension 4: ingress equality vs real node asymmetry

The anti-SPOF instinct in the repo is correct:

- no sacred public box
- multiple public records
- any surviving node should be able to receive the first request

But ingress equality is not the same thing as full node equality.

In the actual runtime world:

- not every node hosts the same services
- not every service class is equally movable
- not every node has the same write authority or state ownership
- secret drift and revision drift still matter
- policy continuity can diverge even when transport works

So the docs have to keep the layer boundary hard:

- ingress plurality is an entry property
- service correctness is a routing, policy, and state property

If those get blurred, the repo starts saying things like:

- "multiple DNS records means HA"
- "any node can receive traffic, therefore the platform is equalized"
- "a healthy reverse proxy means failover works"

Those are exactly the kinds of false options the user is reacting against.

The blunt reading rule here is:

entry equality is cheap compared to semantic equality.

Getting multiple nodes to receive the first packet is not the hard part.
Getting those nodes to preserve the same service meaning after that packet
lands is the hard part.

This distinction is one of the repo's most important anti-bullshit filters.
If a proposal mainly solves first-hop plurality but leaves service meaning to
be reconstructed by hand, it is still an incomplete answer no matter how
impressive the edge diagram looks.

## Tension 5: anti-SPOF instinct vs stateful reality

The project is right to hate single points of failure.
But anti-SPOF pressure becomes misleading the second stateful systems are
treated like stateless HTTP.

For stateless HTTP, a lighter recovery path is plausible:

- multiple ingress nodes
- placement truth
- peer-eligible forwarding
- middleware continuity

For TCP or stateful systems, much harder truths appear:

- who owns writes right now
- whether replicas are actually caught up
- what promotion means
- whether clients rediscover the right survivor
- whether "reachable" still means "correct"

The contradiction is not that stateful HA is impossible.
It is that the same language keeps getting borrowed for drastically different
problem classes.

That is why this repo keeps insisting on separate handling for:

- L7 HTTP behavior
- L4 or raw TCP behavior
- stateful correctness

This is one of the healthiest instincts in the entire documentation set.
If that separation disappears, the repo will start describing reachability as
durability.

And once that happens, anti-SPOF language becomes actively misleading instead
of merely premature.

## Tension 6: operator control vs hidden operator tax

One of the user's deepest complaints is not "I want more automation."
It is:

> why do the supposed options keep turning into systems where the operator
> still has to privately remember topology, failure rules, and semantic drift,
> except now the stack is also more opaque?

That is the most important human tension in the repo.

Manual placement is acceptable.
Hidden dependence on operator memory is not.

Those are different things.

An honest Compose-first system can still allow:

- manual node assignment
- deliberate placement
- explicit service-class exceptions

But it cannot keep requiring the operator to mentally reconstruct:

- which peer currently hosts what
- whether remote state is converged
- which routes vanished after failure
- whether a fallback path preserved auth and middleware
- whether a stateful service is merely up or actually safe

Once the system depends on those hidden reconstructions, the platform is still
centralized.
The operator's head has become the control plane.

That is precisely the condition the user wants this project to escape.

This is probably the most important contradiction in the whole repo because it
is the easiest one to miss.

Manual placement is not the enemy.
Unreadable hidden truth is the enemy.

That is the key distinction the repo must keep defending.
The user's dream is not "remove humans from everything."
It is "stop requiring a private inner map before the platform becomes honest."

## Tension 7: smaller helpers vs accidental control-plane reassembly

The repo does not want to import Swarm, Nomad, or Kubernetes by reflex.
That reluctance is reasonable.

But there is an unavoidable counter-question:

> at what point do enough small helpers become a control plane in every way
> that matters except name and self-awareness?

Every proposed helper carries part of that burden:

- membership
- placement
- route generation
- convergence
- eligibility
- failure reaction

There is nothing wrong with decomposing those responsibilities.
In fact, decomposition may be exactly the right path.

The danger is narrative dishonesty:

- calling the result "just Compose"
- pretending the helper mesh is simpler than a real orchestrator because it was
  assembled incrementally
- refusing to compare its tax honestly against Nomad, OpenSVC, or k3s once it
  reaches control-plane size

The repo should absolutely try the smaller path first.
It should not pretend the smaller path stays small forever by definition.

This is where fake optionality becomes especially dangerous.
A helper mesh can feel like freedom because it avoided adopting a named
orchestrator, while still spending the same categories of complexity:

- membership truth
- placement truth
- convergence truth
- eligibility truth
- failure-reaction truth

At that point the question is no longer "did we install Kubernetes?"
It is "did we rebuild control-plane responsibility without admitting it?"

The mature reading is not:

- helpers good
- orchestrators bad

It is:

- every helper that owns truth, eligibility, convergence, or failure reaction
  is spending control-plane budget
- once enough of that budget has been spent, the repo should compare the
  resulting helper mesh honestly against Nomad, OpenSVC, k3s, or other stronger
  answers

That honest comparison is one of the main things the user has been denied by
typical advice.
Most recommendations either:

- understate the helper tax
- or overstate the inevitability of a heavyweight cluster

This repo is trying to earn a more precise answer than either reflex.

## The real question behind all these tensions

All of these contradictions collapse into one operator-grade question:

> what is the smallest extra truth-bearing layer that actually removes the
> wrong-node and backend-loss lies without replacing Docker readability with a
> larger, equally frustrating fiction?

That is the user's real question.

Not:

- "what is the coolest orchestrator?"
- "what is the most cloud-native answer?"
- "how do we make the docs sound strategic?"

But:

> where is the boundary between useful added structure and one more system that
> steals legibility while still failing to preserve the request under real
> pressure?

That boundary is what makes this repo interesting.

The project is not merely looking for an implementation.
It is trying to discover the smallest amount of added structure that still
feels like a real option instead of another bait-and-switch.

That is why the documentation has to sound more like a reconstruction of a
problem the user has been forced to think about for too long, and less like a
generic architecture memo.

If the knowledgebase forgets that question, it stops reflecting the project.

## What this means for every other page

Every serious page in the docs should now be read through this lens:

- if a page proposes a layer, what missing truth is that layer claiming to own?
- if a page claims resilience, what exact failure class has been survived?
- if a page promises simplicity, whose hidden burden is actually paying for it?
- if a page resists a heavyweight orchestrator, what substitute tax is it
  accepting instead?

Those are the only questions that keep the repo aligned with the dream behind
it.

The dream is not "avoid Kubernetes because Kubernetes is ugly."
The dream is:

> retain real options, real legibility, and real operator control while still
> making a wrong-node or backend-loss event survivable in a way that is honest,
> testable, and not secretly held together by wishful thinking

That is the standard these docs should keep enforcing.
