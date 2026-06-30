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
It is emotionally false to the experience that produced this repo.

The real experience is:

- the ecosystem keeps presenting "options"
- those options sound empowering at the diagram level
- the hard truths keep getting pushed back into operator memory
- then the same systems act surprised when the user says the options still do
  not feel real

This page exists to keep those contradictions alive so the rest of the docs
cannot quietly buy coherence by deleting them.

## What this page is and is not allowed to prove

This page is authoritative about:

- which contradictions the repo is consciously carrying
- why smoothing those contradictions would misdescribe the project
- where dream, runtime, and control-plane growth are still in conflict

This page is not authoritative about:

- whether one side of a contradiction has already won
- whether the runtime has already resolved the tensions
- whether a future layer is already justified just because the tension is clear

## Strongest honest current answer

The project is in motion between strong authoring truth, strong dream clarity,
uneven runtime proof, and several candidate middle layers. That means the repo
is not mainly suffering from lack of cleverness. It is suffering from real
contradictions between what the operator wants to keep legible and what the
runtime would need to own before wrong-node, backend-loss, and stateful claims
stop collapsing back into operator interpretation.

## The contradiction at the center of the whole repo

The dream sounds small when compressed:

> any surviving public node should be able to receive the request, decide
> whether the target is local, preserve the request if it is not, and do all of
> that without forcing the operator to accept a giant scheduler as the default
> answer

That sentence is not modest.

It requires several different truths to coexist:

- ingress truth
- placement truth
- routing truth
- peer eligibility truth
- semantic continuity truth
- stateful authority truth

Single-node Docker lets the operator blur those together.
Multi-node request preservation does not.

That is why the user's frustration feels bigger than "I need better routing."
The wrong-node request is the moment the system has to reveal whether it
actually contains those truths or whether the operator was silently supplying
them all along.

Everything else on this page is just a sharper version of that one conflict.

## What this page is trying to prevent

This page is not balancing tradeoffs for sport.
It is blocking three kinds of drift:

- dream drift: the docs collapse the user's real frustration into generic
  architecture language
- proof drift: intent gets narrated as live behavior
- control-plane drift: the repo quietly rebuilds orchestration in fragments
  while pretending nothing substantial was added

If those drifts are not resisted, the repo becomes impossible to read
honestly.

## Tension 1: Compose readability vs scheduler-grade expectations

### What the repo wants

- the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
  and `compose/` fragments remain the primary authoring surface
- service definitions stay readable
- an operator can still inspect the stack without decoding a full cluster API
- manual service placement remains allowed where it still pays for itself

### What the repo also wants

- any public node can receive the first request
- the node can tell whether the target is local
- if not local, the request can be preserved through peer forwarding or
  equivalent recovery logic
- routes can survive the local failure that made recovery necessary

### Why the tension is real

Compose is good at expressing:

- local container intent
- networks
- configs
- secrets
- labels
- healthchecks

Compose does not naturally own:

- distributed placement truth
- convergence truth
- peer eligibility
- cross-node route persistence
- wrong-node request preservation

So the second the repo asks Compose to participate in those higher-order
truths, one of two things happens:

1. helper layers accumulate around Compose
2. a stronger control substrate gets promoted

That does not make Compose-first wrong.
It makes Compose-first honest only if the repo admits the missing work instead
of describing Compose as if it natively solved distribution.

The user's frustration is aimed directly at this gap.
They are not asking for prettier YAML.
They are asking why "simple Docker" stops offering real options once
distribution matters.

## Tension 2: "No orchestrator by default" vs "someone still has to know what runs where"

The clearest intent surface in the repo is still
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It says:

- no central orchestrator by default
- manual service placement is acceptable
- current-state truth is preferred over scheduler-declared desired state
- local-first serving is preferred
- peer-forward fallback is required when the target is not local

That is clear.
It also implies a question that never goes away:

> what tells a node, right now, where the service really lives?

The repo keeps reaching for smaller answers:

- `services.yaml`
- sync-agent
- failover-agent
- generated Traefik file-provider config
- OpenSVC-assisted inventory
- convergence helpers

Those are not random side ideas.
They are repeated attempts to avoid paying the full cost of a heavyweight
scheduler while still extracting one scheduler-grade truth:

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

That is survivable only if the docs admit it.

## Tension 3: local-first serving vs route persistence under failure

The request model the repo wants is simple to say:

```text
request hits any node
  service is local  -> serve locally
  service is remote -> preserve path and forward to healthy peer
```

The first half is intuitive.
The second half is where most systems start bluffing.

The bluff usually takes the form of:

- the route exists in config
- the proxy can technically target a peer
- the platform therefore has failover

The repo keeps insisting those are different claims.

This is exactly where the priority runtime shows a live contradiction:

- `docker-gen-failover` materially exists
- the master plan explicitly records that the current model can delete routes
  when a container stops

So the repo already knows a route can look dynamic while still failing at the
exact moment the bad day arrives.

That contradiction matters because "local-first" and "survives local backend
loss" are not the same maturity tier.

## Tension 4: serious ingress machinery vs missing distributed truth

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
- the missing cross-node truths are still substantial enough that confidence
  can outrun proof

This is one of the deepest documentation traps in the repo.
The user is not frustrated because nothing exists.
They are frustrated because a lot exists while the hidden burden remains too
intact.

## Tension 5: mesh reachability vs peer eligibility truth

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

## Tension 6: anti-SPOF instinct vs stateful reality

The repo is strongly anti-SPOF by instinct.
That instinct is healthy for ingress and stateless HTTP.

It becomes misleading the second stateful systems are described with the same
tone.

For TCP or stateful surfaces, much harsher truths appear:

- who owns writes?
- what storage path is still sacred?
- how does promotion work?
- what reconnect behavior do clients need?
- what does "survived" even mean for this topology?

The contradiction is not that stateful HA is impossible.
It is that ingress progress and stateful correctness diverge sharply enough that
one can be improving while the other remains largely theoretical.

That is why the docs have to keep saying:

- stateless HTTP may earn the first serious wrong-node proof
- TCP is a different lane
- stateful systems must be proved topology by topology

## Tension 7: helper growth vs "still just Compose"

The repo wants the smallest extra truth-owning layer possible.
That is reasonable.

But helper growth has its own danger:

- `docker-gen-failover`
- sync-agent ideas
- failover-agent ideas
- state registries
- peer broadcast
- drift checks

At some point this can stop being "just Compose plus some glue" and start being
"an unacknowledged control plane that is less inspectable than either raw
Compose or a real orchestrator."

That is a live contradiction, not a theoretical one.

The docs therefore need to keep asking:

- is this helper actually removing hidden-human SPOFs?
- or is it mainly relocating complexity into a narrower, less visible layer?

## Tension 8: the repo can name the wound more clearly than it can yet heal it

This is one of the harshest contradictions, but it is also one of the most
important.

The current knowledgebase can now name quite precisely:

- why wrong-node behavior is the real threshold
- why `services.yaml` keeps reappearing
- why route persistence is different from route existence
- why protected-route continuity matters
- why stateful honesty must stay harsher

That does not mean the runtime has already crossed those thresholds.

It means the repo already understands the shape of the missing answer more
clearly than it has implemented that answer.

That can be healthy if admitted.
It becomes dangerous only when clear articulation gets mistaken for nearness of
completion.

## The contradiction the user is actually reacting to

The user is not merely stuck between "simple" and "advanced."
They are stuck between:

- a readable local-first world that falls apart under wrong-node pressure
- and larger worlds that often demand surrender before proving they remove the
  right pain

That is why so many surrounding options still feel fake.
They often solve one layer while leaving the same hidden burden where it was.

The docs fail if they smooth this into:

- normal platform tradeoffs
- ordinary maturity language
- generic service-discovery problems

This repo is reacting to a narrower and more hostile reality:

too many answers sound like options while still depending on operator folklore.

## What this page should make impossible to say lazily

After reading this page, it should be harder to say:

- "the remaining work is mostly polish"
- "Compose is basically enough now"
- "the mesh means discovery is mostly handled"
- "the helper stack probably already closes the important gaps"
- "if one HTTP route works, the rest is mostly a matter of scale"

Those are exactly the kinds of sentences that re-domesticate the problem into a
cleaner but smaller one.

## The blunt reading

The repo wants the legibility of Compose, the survivability of a distributed
platform, and the honesty of saying "not yet" until the system itself owns more
of the explanation for why a wrong-node request still works.

That is the contradiction.
It is not resolved.
And if the docs start sounding calm about it too early, they have stopped being
useful for the actual dream.
