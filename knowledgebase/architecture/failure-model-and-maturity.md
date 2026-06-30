# Failure Model and Maturity Matrix

This page is the lane-by-lane bad-day maturity ledger for the priority
implementation rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

It is not here to produce one flattering score for the whole stack.

It exists to answer the question that actually matters:

> when the nice path breaks, which truths already belong to the system, and
> which ones still quietly fall back into private operator custody?

That is what maturity means in this repo.

## What this page is and is not allowed to prove

This page is allowed to prove:

- which lanes already have real runtime shape
- which lanes are still mostly intent or partial proof
- which hidden private sentences still survive in each lane
- what sort of proof packet would justify upgrading a lane

This page is not allowed to prove:

- that more visible machinery equals maturity
- that the stack has one meaningful global maturity score
- that a stronger lane in one area upgrades adjacent lanes automatically
- that helper presence, dashboards, or prose coherence are themselves proof of
  transferred truth

## The fake maturity move this page is supposed to prevent

The bad version of maturity language in this repo sounds like:

- there are more services
- there are more helpers
- there are more dashboards
- there are more edge components
- therefore the platform is becoming mature

That is exactly the move this page should interrupt.

This repo does not care about maturity as visual seriousness.
It cares about maturity as transferred bad-day truth.

## Why maturity has to stay lane-specific

The stack is already too real to be called a toy.

It has:

- multiple Compose fragments
- a substantial Traefik edge
- DNS plurality assumptions
- private-network identity surfaces
- protected admin routes
- TCP routers
- monitoring, alerting, and helper components

But none of that justifies one platform-wide maturity sentence.

The user's real contract is still:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

Every maturity label here is downstream of that event model.
If a label cannot say how much of that contract is actually system-owned for
one lane, it is too flattering.

## The hidden maturity test

A lane is only more mature if one fewer private sentence has to be completed by
the operator when something goes wrong.

That is the test.

Not:

- more labels
- more services
- more dashboards
- more helper containers
- more sophisticated prose

The test is:

> after calling this lane more mature, what exact sentence does the operator no
> longer need to know off-book?

If the answer is `none`, then the lane may be better instrumented or better
documented, but it has not matured in the way this repo actually cares about.

Another way to say the same thing:

> a lane only matures when the wrong healthy node becomes less socially
> embarrassing.

## Maturity levels

### `Intent-shaped`

The repo clearly wants the lane, but the tracked runtime does not yet prove the
behavior.

### `Runtime-shaped`

Real relevant components are present in the tracked runtime, but the lane still
depends on hidden human joins or unproven transitions.

### `Partial proof`

Some real burden has moved into system-owned behavior or the docs have captured
an important failure truth, but the decisive handoff is still incomplete.

### `Trustworthy for this lane`

The exact lane has enough narrow evidence that the docs can use stronger
language without borrowing confidence from adjacent lanes.

Most important lanes are still below that last line.

## The proof-packet rule

No lane should be upgraded unless the docs can point to a narrow proof packet
that includes:

- the exact lane being upgraded
- the hidden operator burden that used to exist
- the system-owned truth or artifact that replaced it
- the exercised failure condition, comparison, or drill
- the stronger sentence now allowed
- the sentence still forbidden

Without that packet, maturity is still mostly atmosphere.

## The matrix

| Lane | What the current worktree materially proves | Current maturity | Private sentence still surviving | What would honestly move the lane |
| --- | --- | --- | --- | --- |
| Public first-hop plurality | The architecture dream clearly targets multiple public nodes; `cloudflare-ddns` is live; `.github/copilot-instructions.md` explicitly treats any-node entry as a design goal | Runtime-shaped | `Plural DNS still does not tell me what the healthy node should do next.` | Show traffic arriving at more than one public node while keeping the docs explicit that first-hop plurality is still weaker than wrong-node success |
| Local edge execution | The stack already has live edge services such as `traefik`, `tinyauth`, `nginx-traefik-extensions`, `crowdsec`, `homepage`, `dozzle`, `code-server`, `whoami`, and protected metrics surfaces | Runtime-shaped | `A local success still does not tell me whether distributed truth moved.` | Prove one named local route with its real policy stack visible from ingress to backend |
| Placement truth | Intent surfaces keep converging on a `services.yaml`-like current-state registry | Intent-shaped | `I still personally know what runs where.` | Introduce one live tracked placement authority consumed by routing or eligibility logic |
| Convergence truth | Plans and research preserve sync pressure for secrets, Compose state, and node agreement | Intent-shaped | `Even if the system says where something runs, I still need to know whether that claim is stale.` | Expose inspectable drift or sync truth that the forwarding layer can actually trust |
| Peer-eligibility truth | Headscale is live; private peer connectivity matters; peer-aware routing is a first-class pressure | Intent-shaped | `Reachable still is not the same as safe.` | Show one peer choice made from shared eligibility truth stricter than reachability |
| Stateless wrong-node HTTP | The desired contract is explicit and the runtime contains plausible stateless HTTP candidates | Intent-shaped | `The wrong-node success story still mostly lives in architecture prose.` | Force one request onto the wrong healthy node and show correct completion |
| Backend-loss HTTP survival | The repo already knows helper presence is weaker than route durability and already tracks the helper trap where failover can vanish under stop events | Partial proof | `I still privately do not know whether the fallback survives the failure that made it matter.` | Re-run a named route while stopping the preferred backend and preserve the evidence |
| Protected-route continuity | The runtime ships `nginx-auth@file` and TinyAuth on real routes | Runtime-shaped | `A peer-forwarded route may still not be the same protected route.` | Compare one protected route locally and after wrong-node handoff with auth and middleware meaning preserved |
| TCP forwarding | The root runtime already exposes `mongodb` and `redis` through Traefik TCP routers | Runtime-shaped | `A TCP port answering is still not authority continuity.` | Keep transport proof separate from any stateful resilience claims |
| Headscale control-plane resilience | Headscale is live, externally routed, and useful | Runtime-shaped | `Headscale truth still lives on one SQLite authority path.` | Define and prove authority transition before speaking of Headscale HA |
| Stateful databases and queues | `mongodb`, `redis`, `nuq-postgres`, `litellm-postgres`, `rabbitmq`, and `qdrant` are real runtime dependencies | Intent-shaped | `I still personally know the writer is singular or undefined.` | Treat each stateful class separately and prove authority, promotion, fencing, and client rediscovery |
| Operator inspectability | Docs and dashboards are better than before; service visibility exists across several surfaces | Partial proof | `I can still be forced to privately explain why the distributed-looking decision was valid.` | Surface inspectable evidence for locality, peer choice, fallback origin, and authority state |

## What the matrix is really scoring

This matrix is not scoring:

- elegance
- ecosystem status
- quantity of machinery
- similarity to a cluster product

It is scoring how much these sentences still survive:

- `I still personally know where this service really lives.`
- `I still personally know which peer is actually safe right now.`
- `I still personally know whether the fallback is real or theatrical.`
- `I still personally know whether this protected route is still the same route after handoff.`
- `I still personally know whether the writer is still singular.`

If those sentences do not shrink, the maturity label should stay harsh.

## The current maturity story in plain English

The repo is already beyond "one reverse proxy and some containers."

But it is still below the level where the platform itself can be trusted to
carry the most important distributed questions without private human memory.

The plain-English state is:

- the edge is real
- the anti-SPOF dream is coherent
- the lane decomposition is increasingly honest
- the burden transfer is still incomplete in the most important failure modes

That is why this page refuses one global badge.

The problem is not that nothing exists.
The problem is that the decisive distributed truth is still too often implied,
remembered, or reconstructed rather than owned.

## The lanes easiest to flatter by accident

The following lanes are especially vulnerable to inflated language:

- public entry, because plural DNS sounds more finished than it is
- Traefik-centered ingress, because route execution sounds like route truth
- protected routes, because successful response codes hide policy discontinuity
- helper-driven fallback, because generated config sounds closer to survival
- stateful services, because reachability sounds much closer to resilience than
  it really is

## What a lane upgrade should actually feel like

A real lane upgrade should feel less like:

- the docs got sharper
- the diagrams got better
- the platform sounds more serious

and more like:

- one bad-day question now has a system-owned answer
- one private operator sentence became weaker
- one category of fake confidence became illegal

If the change does not feel like that, the lane probably improved cosmetically
more than it matured structurally.

## Bottom line

This matrix is not here to reassure the reader that the stack is maturing
nicely.

It is here to keep forcing the harder question:

> when the nice path breaks, what truth is actually owned by the system now,
> and what truth still falls back into one human's head?

That is the only maturity scale in this repo that matters enough to trust.
