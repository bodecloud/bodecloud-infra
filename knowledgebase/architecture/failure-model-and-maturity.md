# Failure Model and Maturity Matrix

This page is the failure-lane maturity ledger for the priority implementation
rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

It is not trying to produce one flattering score for the whole stack.

It is trying to answer the question that actually matters here:

> when the nice path breaks, which truths already belong to the system, and
> which ones still quietly fall back into private operator custody?

That is the maturity standard for this repo.

## Why maturity has to stay lane-specific

This stack is already too real to be called a toy.

It has:

- multiple Compose fragments
- a substantial Traefik edge
- DNS plurality assumptions
- private-network identity surfaces
- protected admin routes
- TCP routers
- monitoring, alerting, and helper components

But none of that justifies a platform-wide maturity sentence by itself.

The user is not asking whether the repo looks advanced.
The user is asking whether the bad day still requires the operator to privately
finish the story.

That answer is still different by lane.

## The hidden maturity test

A lane is only more mature if one fewer private sentence has to be completed by
the operator when something goes wrong.

That is the test.

Not:

- more services
- more labels
- more dashboards
- more helper containers
- more sophisticated prose

The test is:

> after calling this lane more mature, what exact sentence does the operator no
> longer need to know off-book?

If the answer is "none," then the lane may be better instrumented or better
documented, but it has not matured in the way this repo actually cares about.

## Maturity levels

### `Intent-shaped`

The repo clearly wants the lane, but the current runtime does not yet prove the
behavior.

### `Runtime-shaped`

Real relevant components are present in the tracked runtime, but the lane still
depends on hidden human joins or unproven transitions.

### `Partial proof`

Some meaningful burden has moved into system-owned behavior or the repo has
captured an important failure truth, but the decisive handoff is still
incomplete.

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

Without that packet, "maturity" is still mostly atmosphere.

## The matrix

| Lane | What the current worktree materially proves | Current maturity | Hidden operator burden still present | What would move the lane honestly |
| --- | --- | --- | --- | --- |
| Public first-hop plurality | The architecture dream clearly targets multiple public nodes; `cloudflare-ddns` is live; `.github/copilot-instructions.md` explicitly treats any-node entry as a real design goal | Runtime-shaped | The operator still cannot equate first-hop plurality with preserved service meaning | Show traffic arriving at more than one public node while keeping the docs honest that this is still weaker than wrong-node success |
| Local edge execution | The stack already has live edge services such as `traefik`, `tinyauth`, `nginx-traefik-extensions`, `crowdsec`, `whoami`, `dozzle`, `code-server`, `homepage`, and protected metrics/admin routes | Runtime-shaped | Local route success can still be mistaken for cross-node truth ownership | Prove one named local route with its actual policy stack visible from ingress to backend |
| Placement truth | The repo's intent surfaces keep converging on a `services.yaml`-like registry as a lightweight current-state source of truth | Intent-shaped | "What runs where right now?" is still safest when answered from private memory | Introduce one live tracked placement authority consumed by routing or eligibility logic |
| Peer eligibility truth | Headscale is live; private peer connectivity is part of the real stack; peer-aware routing is a first-class design pressure | Intent-shaped | Reachable peers are not yet the same thing as semantically safe peers | Show one peer choice made from shared current truth rather than folklore |
| Stateless wrong-node HTTP | The desired contract is explicit; the runtime already contains plausible surfaces like `whoami`, `wishlist`, `homepage`, `chat-analytics`, and `searxng` | Intent-shaped | Wrong-node success is still more architectural story than proven property | Force one request onto the wrong healthy node and show correct completion |
| Backend-loss HTTP survival | The repo already knows helper presence is weaker than route durability and already tracks the `docker-gen-failover` trap | Partial proof | The rescue route can still disappear during the failure it is supposed to cover | Re-run a named route while stopping the preferred backend and preserve the evidence |
| Protected-route continuity | The runtime already ships `nginx-auth@file` on real routes and TinyAuth as a live auth component | Runtime-shaped | A peer-forwarded route may still stop being the same protected service | Compare one protected route locally and after wrong-node handoff with auth and middleware behavior preserved |
| TCP forwarding | The root runtime already exposes `mongodb` and `redis` through Traefik TCP routers | Runtime-shaped | Transport reachability is easy to overread as stateful dignity | Separate transport success proof from authority and failover claims |
| Headscale control-plane resilience | `headscale-server` and `headscale` are live, externally routed, and monitored | Runtime-shaped | Current config still roots state in `/var/lib/headscale/db.sqlite`, so public reachability is still not authority redundancy | Define and prove authority transition before speaking of Headscale HA |
| Stateful databases and queues | `mongodb`, `redis`, `nuq-postgres`, `litellm-postgres`, and `rabbitmq` are real runtime dependencies | Intent-shaped | Write authority, promotion, persistence, and rediscovery semantics remain singular or undefined | Treat each stateful class separately and define authority, promotion, and client rediscovery |
| Drift and convergence control | Research and plans already preserve sync pressure for secrets, Compose state, and node agreement | Intent-shaped | A wrong-node request can still land on a semantically different node revision | Expose drift-detection truth that peer-forward decisions can actually trust |
| Operator inspectability | The docs and metrics are better than before; dashboards exist for key services including TinyAuth, Headscale, and `docker-gen-failover` | Partial proof | The operator can still be required to privately explain why the distributed decision was valid | Surface inspectable evidence for locality, peer choice, and fallback route origin |

## The current maturity story in plain English

The repo is already beyond "one reverse proxy and some containers."

But it is still below the level where the platform itself can be trusted to
carry the most important distributed questions without human private memory.

The current state is best described like this:

- the edge is real
- the anti-SPOF dream is coherent
- the lane decomposition is increasingly honest
- the burden transfer is still incomplete in the most important failure modes

That is why this page refuses a single global badge.

## The lanes most likely to be overclaimed

### Public entry

Plural DNS and several healthy public nodes are meaningful gains.

They still do not prove:

- locality truth
- peer eligibility truth
- route continuity
- backend-loss survival

### Traefik-centered ingress

Traefik deserves credit for being a real execution surface in the live stack.

It does not deserve automatic credit for creating distributed current truth.

### Protected routes

The presence of middleware labels is better than nothing.

It is still weaker than proving that peer-forward handoff preserves the same
protected route semantics.

### Stateful services

This is where most self-hosted platforms flatter themselves.

Clean hostnames, TCP exposure, and restart policy are not substitutes for
authority, promotion, or write dignity.

Headscale is the clearest warning because its current config still resolves to
SQLite on one state path.

## What this repo has already learned correctly

The repo has already learned several important negative truths:

- service discovery, not mere placement, is the hard missing piece once manual
  placement and plural DNS are accepted
- route helpers can exist and still fail the exact scenario they are meant to
  rescue
- protected routes cannot be judged only by response code
- stateful services must stay on a harsher proof track than HTTP ingress
- any future orchestrator or controller has to justify itself against the
  hidden operator SPOF, not just against aesthetic desire for "more cluster"

Those lessons matter because they stop the docs from exaggerating the current
stack.

## What maturity is allowed to mean here

The docs can only say a lane is more mature when one of these becomes true:

- placement truth leaves private memory
- peer eligibility becomes system-owned
- a rescue route survives the exact failure that used to delete it
- a protected forwarded route demonstrably preserves the same policy meaning
- a stateful surface gains explicit authority and promotion semantics
- a drift check proves nodes are aligned on the truths routing decisions depend
  on

The docs must not claim maturity because:

- the stack got bigger
- the diagrams got cleaner
- the happy path works
- a future orchestrator sounds credible

## The smallest honest maturation sequence

If the repo wants one realistic v1 maturity arc, it looks like this:

1. externalize present-tense placement truth
2. prove one stateless wrong-node HTTP route
3. repeat that route under preferred-backend loss
4. compare one protected route before and after peer handoff
5. keep TCP and stateful services on separate harsher tracks
6. only then promote stronger lane language

That sequence matters because it removes hidden sentences from the operator's
head instead of just adding more infrastructure nouns.

## Strongest honest maturity sentence today

The strongest current sentence this page allows is:

> the repo has a serious Compose-first runtime and a serious anti-SPOF design
> pressure, but its most important failure lanes are still maturing separately
> because the platform does not yet own all the truth needed for wrong-node,
> backend-loss, protected-route, and stateful correctness claims
