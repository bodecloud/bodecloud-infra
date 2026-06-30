# Design Tensions and Contradictions

This page exists because the clean version of this project is a lie.

The clean version sounds comforting:

- Docker Compose is readable
- Cloudflare can point at more than one node
- Traefik can route
- a few helpers can patch over the rest
- therefore the platform can become multi-node and anti-SPOF without much real
  contradiction

That is exactly the kind of story the user is tired of.

It is not only technically incomplete.
It is emotionally false to the experience that produced `bolabaden-infra`.

The real experience is harsher:

- the ecosystem keeps presenting "options"
- those options often sound empowering at the diagram level
- the hard truths keep getting pushed back into operator memory
- then the same systems act surprised when the user says the options still do
  not feel real

This page exists to keep those contradictions visible so the rest of the docs
cannot quietly buy coherence by deleting them.

## What this page is and is not allowed to prove

This page is authoritative about:

- which contradictions the repo is consciously carrying
- why smoothing those contradictions would misdescribe the project
- where dream, runtime, helper growth, and proof boundaries are still in
  conflict
- why the repo still feels like it lacks real options despite having many
  technologies in play

This page is not authoritative about:

- whether one side of a contradiction has already won
- whether the runtime has already resolved those tensions
- whether a future layer is justified just because the tension is clear
- whether better articulation means the repo is close to closure

This is a tension-preservation page.
It is not a victory page.

## Strongest honest current answer

The repo is in motion between:

- strong dream clarity
- stronger honesty boundaries than before
- a serious live Compose and ingress surface
- several candidate middle layers
- uneven proof across the actual failure lanes

That means the project is not mainly suffering from lack of cleverness.
It is suffering from real contradictions between what the operator wants to
keep legible and what the system itself would need to own before wrong-node,
backend-loss, protected-route, and stateful claims stop collapsing back into
operator interpretation.

## The contradiction at the center of the whole repo

The dream sounds small when compressed:

> any surviving public node should be able to receive the request, determine
> whether the target is local, preserve the request if it is not, and do all
> of that without forcing the operator to adopt a giant scheduler as the
> default answer

That sentence is not small.

It requires several different truths to coexist:

- public-entry truth
- locality truth
- placement truth
- peer-eligibility truth
- route-persistence truth
- policy-continuity truth
- service-class truth
- stateful-authority truth

Single-node Docker lets the operator blur many of those together.
Multi-node request preservation does not.

That is why the user’s frustration is bigger than "I need better routing."
The wrong-node request is the moment where the system has to reveal whether it
contains those truths itself or whether the operator was silently supplying
them all along.

Everything else on this page is just a sharper version of that one conflict.

## What this page is trying to prevent

This page is not balancing tradeoffs for sport.
It is blocking three kinds of drift:

- dream drift: the docs collapse the user's actual wound into generic
  architecture language
- proof drift: intent gets narrated as if it were live behavior
- control-plane drift: the repo quietly rebuilds orchestration in fragments
  while pretending nothing substantial was added

If those drifts are not resisted, the repo becomes easier to read and harder
to trust.

## Tension 1: Compose readability vs scheduler-grade expectations

### What the repo wants to keep

- the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  plus `compose/` fragments remain the primary authoring surface
- service definitions stay readable
- an operator can still inspect the platform without decoding a full cluster
  API
- manual service placement remains allowed where it still pays for itself

### What the repo also wants to become true

- any healthy public node can receive the first request
- the receiving node can tell whether the target is local
- if not local, the request can be preserved through peer-forward logic or
  equivalent recovery
- the route needed for rescue survives the failure that made rescue necessary

### Why the tension is real

Compose is good at expressing:

- local container intent
- configs
- secrets
- labels
- networks
- healthchecks

Compose does not naturally own:

- distributed placement truth
- convergence truth
- peer eligibility
- cross-node route durability
- wrong-node request preservation

So the moment the repo asks Compose to participate in those higher-order
truths, one of two things happens:

1. helper layers start accumulating around Compose
2. a stronger truth-owning substrate gets promoted

That does not make Compose-first wrong.
It makes Compose-first honest only if the repo admits the missing work instead
of describing Compose as if it natively solved distribution.

The user is not asking for prettier YAML.
They are asking why "simple Docker" stops offering real options once
distribution matters.

## Tension 2: "No orchestrator by default" vs "someone still has to know what runs where"

The clearest intent surface in the repo is still
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It explicitly says:

- no central orchestrator by default
- manual service placement is acceptable
- current-state truth is preferred over scheduler-declared desired state
- local-first serving is preferred
- peer-forward fallback is part of the target contract

That is clear.
It also implies the question that never goes away:

> what tells a node, right now, where the service really lives?

The repo keeps reaching for smaller answers:

- `services.yaml`
- sync-agent ideas
- failover-agent ideas
- generated Traefik file-provider config
- OpenSVC-assisted inventory
- drift and convergence helpers

Those are not random side ideas.
They are repeated attempts to avoid paying full heavyweight-orchestrator tax
while still extracting one scheduler-grade truth:

> current placement

The present-tense contradiction is sharp:

- the docs repeatedly depend on something like `services.yaml`
- the priority worktree does not currently prove a live tracked root
  [`services.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/services.yaml)
  exists and is consumed by routing

So this is not only a philosophical tension.
It is an evidence tension:

the repo can describe the missing middle layer more clearly than the runtime
proves it exists.

That is survivable only if the docs admit it directly.

## Tension 3: first hop plurality vs preserved request meaning

Cloudflare multi-record entry is part of the dream.
It is also one of the easiest places for the docs to become fraudulent.

Plural public entry can make the stack feel distributed because:

- more than one node is reachable
- DNS does not point at one sacred box
- a failure at the public edge feels less final

But the user is not asking only for more than one reachable first hop.
They are asking whether the request still keeps its meaning after landing on a
healthy node that does not host the target locally.

That means first-hop plurality is in tension with the deeper requirement:

- the receiving node must know what the request means
- it must know where the service lives now
- it must know which peer is acceptable now
- it must know whether the rescue route survives the failure

Without that, DNS plurality remains emotionally satisfying and operationally
thinner than the language around it suggests.

## Tension 4: local-first serving vs route persistence under failure

The request model the repo wants is simple to say:

```text
request hits any node
  service is local  -> serve locally
  service is remote -> preserve path and forward to healthy peer
```

The first half is intuitive.
The second half is where systems start bluffing.

The bluff usually takes one of these forms:

- the route exists in config
- the proxy can technically target a peer
- the mesh is reachable
- therefore the platform has failover

The repo keeps insisting those are different claims.

This is where the priority runtime shows a live contradiction:

- `docker-gen-failover` materially exists
- the master plan explicitly records that the current approach can delete
  routes when containers stop

So the repo already knows a route can look dynamic while still failing at the
exact moment the bad day arrives.

That contradiction matters because:

- local-first behavior
- wrong-node survival
- backend-loss survival

are three different maturity tiers, not one.

## Tension 5: serious ingress machinery vs missing distributed truth

The priority runtime already has a real edge stack:

- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `cloudflare-ddns`
- `docker-gen-failover`

This is enough machinery to sound mature.
That is precisely why it is dangerous.

The contradiction is:

- the edge stack is real enough to inspire confidence
- the missing cross-node truths are still serious enough that confidence can
  outrun proof

That is one of the deepest traps in the whole knowledgebase.
The user is not frustrated because nothing exists.
They are frustrated because a lot exists while the hidden burden remains too
intact.

## Tension 6: mesh reachability vs peer eligibility truth

Headscale is materially live through
[`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml).

That proves the repo is not pretending node-to-node connectivity exists.
It does not settle:

- which peer currently owns the requested service
- whether that peer is on an acceptable revision and secret surface
- whether it should be trusted right now for this route

So the contradiction is:

- transport and identity are becoming real
- routing truth is still weaker than transport truth

The ecosystem often stops here and calls that "service discovery."
The user keeps rejecting exactly that stop point.

## Tension 7: protected-route semantics vs "some response happened"

This repo is not only about transport.
It is also about whether the request still counts as the same protected route
after handoff.

That means protected-route success is stricter than:

- a backend answered
- the HTTP status is green
- a login page still appeared
- a forwarded response looked plausible

The repo already has policy-bearing edge pieces:

- TinyAuth
- nginx forward-auth extensions
- CrowdSec
- Traefik middleware

The contradiction is:

- the edge stack is rich enough to define route identity
- that same richness makes wrong-node parity harder to prove honestly

If the fallback path subtly changes auth, headers, middleware ordering, or
security posture, then the route did not really survive as the same route.

That is why "some response happened" is one of the most dangerous fake-success
standards in the repo.

## Tension 8: anti-SPOF instinct vs stateful reality

The repo is strongly anti-SPOF by instinct.
That instinct is healthy for ingress and stateless HTTP.

It becomes misleading the moment stateful systems are narrated with the same
tone.

For TCP or stateful surfaces, much harsher truths appear:

- who owns writes?
- what storage path is still sacred?
- how does promotion work?
- what reconnect behavior do clients need?
- what does "survived" even mean for this topology?

The contradiction is not that stateful HA is impossible.
It is that ingress progress and stateful correctness diverge sharply enough
that one can be improving while the other remains largely theoretical.

That is why the docs keep having to say:

- stateless HTTP may earn the first serious wrong-node proof
- TCP is a different lane
- stateful systems must be proved topology by topology

## Tension 9: helper growth vs "still just Compose"

The repo wants the smallest extra truth-owning layer possible.
That is reasonable.

But helper growth has its own danger:

- `docker-gen-failover`
- sync-agent ideas
- failover-agent ideas
- state registries
- peer broadcast
- drift checks

At some point this can stop being:

- "just Compose plus a bit of glue"

and start being:

- "an unacknowledged control plane that is less inspectable than either raw
  Compose or a real orchestrator"

That is a live contradiction, not a theoretical one.

The repo therefore has to keep asking:

- is this helper actually removing hidden-human SPOFs?
- is it externalizing truth?
- or is it relocating complexity into a narrower, less visible layer?

If that question disappears, the project can become more complex without
becoming more honest.

## Tension 10: the repo can name the wound more clearly than it can yet heal it

This is one of the harshest contradictions and one of the most important.

The current knowledgebase can now name quite precisely:

- why wrong-node behavior is the real threshold
- why `services.yaml` keeps reappearing
- why route persistence is different from route existence
- why protected-route continuity matters
- why stateful honesty must stay harsher

That does not mean the runtime has crossed those thresholds.

It means the repo understands the shape of the missing answer more clearly than
it has implemented that answer.

That can be healthy if admitted.
It becomes dangerous only when clear articulation gets mistaken for nearness of
completion.

## The contradiction the user is actually reacting to

The user is not mainly stuck between "simple" and "advanced."
They are stuck between:

- a readable local-first world that falls apart under wrong-node pressure
- larger worlds that often demand surrender before proving they remove the
  right pain

That is why so many surrounding options still feel fake.
They often solve one layer while leaving the same hidden burden where it was.

The docs fail if they smooth this into:

- normal platform tradeoffs
- ordinary maturity language
- generic service-discovery problems

This repo is reacting to a narrower and more hostile reality:

too many answers sound like options while still depending on operator folklore.

## The burden-transfer test

Every future helper, registry, platform, or orchestration decision should be
read through one question:

> which hidden explanation will the system own after this exists that the
> operator currently has to supply from memory?

If the answer is weak, then the proposal may still be useful, but it is not
yet paying down the core contradiction.

That test is more important here than ordinary elegance, popularity, or
platform fashion.

## What this page should make impossible to say lazily

After reading this page, it should be harder to say:

- "the remaining work is mostly polish"
- "Compose is basically enough now"
- "the mesh means discovery is mostly handled"
- "the helper stack probably already closes the important gaps"
- "if one HTTP route works, the rest is mostly a matter of scale"
- "the project just needs a little more automation"

Those are exactly the sentences that domesticate the problem into a smaller,
cleaner, and less honest one.

## Bottom line

The repo wants:

- the legibility of Compose
- the survivability of a distributed platform
- and the honesty of saying "not yet" until the system itself owns more of the
  explanation for why a wrong-node request still works

That is the contradiction.
It is not resolved.

And if the docs start sounding calm about it too early, they have stopped
being useful for the actual dream.
