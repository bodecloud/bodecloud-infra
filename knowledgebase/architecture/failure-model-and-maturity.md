# Failure Model and Maturity Matrix

This page is the lane-by-lane maturity map for the priority implementation
rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

It is not asking "how impressive is the stack?"
It is asking:

> when the good path breaks, which truths are already system-owned, which still
> fall back into private operator custody, and what exact proof would move one
> more lane out of that private space?

That is the maturity question that matters in this repo.

## What this page is and is not allowed to prove

This page is authoritative about:

- the real failure lanes the repo is trying to mature
- what the current worktree materially proves in each lane
- which hidden operator burden still survives in each lane
- what next proof would justify stronger language

This page is not authoritative about:

- a global maturity score
- future controller choices being already justified
- one lane silently upgrading another
- "directionally close" being good enough

## Strongest honest current answer

The repo has a serious Compose-first runtime and a serious anti-SPOF dream.
It does not yet have one platform-wide proof that the system owns wrong-node
truth, backend-loss truth, protected-route continuity, and stateful authority.

So the only honest maturity model is lane-specific.

The current worktree is:

- beyond toy-stack status
- beyond single-proxy simplification
- still below "request-preserving personal cloud" status

That is why this page separates lanes instead of flattening them.

## How maturity is measured here

A lane matures only when one fewer private reconstruction step is required from
the operator on the bad day.

That means:

- more services does not equal more maturity
- more labels does not equal more maturity
- a better explanation does not equal more maturity
- a stronger local route does not equal more maturity elsewhere

The hidden test is simple:

> after calling this lane more mature, what exact sentence does the operator no
> longer need to finish privately?

If the answer is "none," then the lane may be better documented or more
instrumented, but it is not more mature in the sense this repo cares about.

## Maturity labels

### `Intent-shaped`

The lane is clearly named in the docs and architecture dream, but the current
runtime does not yet prove the behavior.

### `Runtime-shaped`

Real components are live and relevant, but the lane still depends on hidden
operator interpretation or missing joins.

### `Partial proof`

Some meaningful evidence or known failure isolation exists, but the decisive
burden transfer is still incomplete.

### `Trustworthy for this lane`

This exact lane has enough narrow proof that the docs can speak more strongly
without borrowing confidence from neighboring lanes.

Most important lanes are still below that last level.

## Lane proof packet standard

Before any lane gets upgraded, the docs should be able to point to a packet
containing:

- the lane being upgraded
- the previous hidden operator burden
- the new system-owned truth or artifact
- the exercised failure condition or comparison
- the stronger sentence now allowed
- the sentence that remains forbidden

Without that packet, maturity is still mostly atmosphere.

## The matrix

| Lane | What the current worktree materially proves | Current maturity | Hidden operator burden still present | Next honest maturity step |
| --- | --- | --- | --- | --- |
| Public node-entry reachability | The architecture dream explicitly targets any-node first hop; `cloudflare-ddns` is live; first-hop plurality is a first-class design pressure in `.github/copilot-instructions.md` | Runtime-shaped | The operator still cannot treat first-hop plurality as preserved service meaning | Show traffic arriving at more than one public node without upgrading that to wrong-node success |
| Local edge execution | The edge stack is real: `traefik`, `tinyauth`, `nginx-traefik-extensions`, `crowdsec`, routers, middlewares, healthchecks, local app routes | Runtime-shaped | Local success can still be mistaken for cross-node truth ownership | Prove one named local route with its real policy stack visible end to end |
| Placement and locality truth | The docs repeatedly converge on a `services.yaml`-like registry, but the priority runtime does not yet prove a live consumed root registry | Intent-shaped | "What runs where right now?" is still most safely answered from private memory | Introduce one tracked placement-truth surface visibly consumed by routing or eligibility logic |
| Peer eligibility truth | Headscale is live, private reachability is not hypothetical, and research surfaces name sync and peer-broadcast pressure | Intent-shaped | Reachable is not the same as semantically safe; the receiving node still lacks a proven shared answer to "which peer is valid now?" | Demonstrate one peer-selection decision from shared current truth rather than folklore |
| Fallback-route persistence | `docker-gen-failover` exists and the repo already records that it can delete routes when a backend dies | Partial proof | The rescue path can still disappear during the exact failure it is meant to absorb | Harden or replace route generation, then re-run a backend-loss drill |
| Protected-route continuity | TinyAuth, Nginx forward-auth extensions, CrowdSec, and Traefik middleware surfaces are all live | Runtime-shaped | A forwarded route may still not mean the same protected service after handoff | Compare one protected route locally and through peer-forward path, preserving auth and middleware behavior |
| Stateless HTTP wrong-node success | The target contract is explicit and there are suitable live HTTP surfaces such as `whoami`, `wishlist`, and `code-server` | Intent-shaped | Wrong-node success is still a hope, not a proven property | Force one request onto the wrong healthy node and show correct completion |
| Backend-loss HTTP survival | The repo already distinguishes this from wrong-node rescue and already documents the route-deletion trap | Intent-shaped | Even a good wrong-node story may still collapse when the preferred backend disappears | Stop the preferred backend and prove whether the rescue path survives with the same visible contract |
| TCP forwarding | The root runtime already exposes TCP routes such as `mongodb` and `redis` through Traefik TCP | Runtime-shaped | Transport reachability is easy to overread as stateful resilience | Split proof into transport success, client behavior, and authority semantics |
| Headscale control-plane resilience | Headscale is live and important, but active config still uses SQLite at `/var/lib/headscale/db.sqlite` | Runtime-shaped | The mesh still depends on a singular state authority | Only speak of HA after a real authority transition path exists |
| Stateful databases and queues | Root `mongodb`, root `redis`, Firecrawl `nuq-postgres` and `rabbitmq`, and `litellm-postgres` are all real current dependencies | Intent-shaped | Write authority, replication, promotion, reconnect, and disk truth remain singular or undefined | Per service class, define authority, promotion, failover, and rediscovery semantics |
| Convergence and drift control | Research and plans clearly name secret sync, compose sync, and node-alignment pressure | Intent-shaped | A forwarded request may still land on a semantically different node revision or secret surface | Expose drift detection and prove nodes agree on the truth peer-forward decisions depend on |

## The lanes most likely to be overclaimed

### Public entry

Plural DNS and multiple public nodes are real gains.
They are still only first-hop gains until the receiving node can preserve
service meaning.

### Traefik-centered ingress

Traefik is a strong execution surface.
That makes it easy to overcredit.
It executes routing well; it does not by itself create distributed truth.

### Stateful services

This is where most infra stacks cheat hardest.
Reachability, clean TCP exposure, and container restart behavior do not answer
ownership or promotion.

## What "more mature than before" is allowed to mean

The docs may use stronger maturity language only when one of these becomes
true:

- current placement truth is externalized from private memory
- the receiving node can justify peer choice from shared current truth
- the fallback route survives the exact failure that used to delete it
- a protected forwarded route demonstrably preserves the same policy meaning
- a stateful surface gains explicit authority and promotion semantics

The docs must not upgrade maturity language because:

- the stack looks larger
- the prose sounds clearer
- a happy-path `200` happened
- a controller idea sounds more serious

## The current failure model in plain English

Today the platform can still fail in ways that directly offend the user's
benchmark:

- traffic can land on a healthy node that still lacks trustworthy shared
  placement truth
- a rescue mechanism can be present and still disappear during the relevant
  failure
- a protected route can still answer differently after handoff even when the
  response code stays green
- TCP reachability can still be mistaken for stateful dignity
- the operator can still remain the safest keeper of topology truth

That last item is the real anti-SPOF reading.

## What would materially change the maturity story

The smallest meaningful maturity sequence is:

1. move present-tense placement truth out of private memory
2. prove one stateless wrong-node HTTP route
3. prove whether the same route survives backend loss
4. compare protected-route continuity across that handoff
5. only then promote stronger HTTP-lane language
6. keep TCP and stateful lanes on their own harsher proof tracks

Until then, the honest sentence is:

> the repo has real runtime machinery and sharper truth custody than before,
> but its most important failure lanes are still maturing separately rather
> than composing into one platform-wide victory
