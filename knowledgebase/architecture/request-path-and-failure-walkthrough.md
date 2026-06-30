# Request Path and Failure Walkthrough

This page exists because the user's real question is not:

> what components are in the ingress stack?

It is:

> when one hostname is requested, what literally happens next, which files own
> each decision, what does the receiving node know on its own, and where does
> the story stop being system-owned and turn back into private operator memory?

The neighboring smaller question this page must not collapse into is:

> what is our reverse-proxy chain?

That smaller question is too easy.
Plenty of stacks can name Cloudflare, Traefik, auth, and observability without
answering the user's actual complaint about wrong-node humiliation.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what the target request contract is
- which live runtime components already participate in that contract
- where the request path is already concrete in the root Compose surface
- where current proof stops for wrong-node and backend-loss behavior

This page is not allowed to prove:

- generic wrong-node success today
- generic peer-forward fallback today
- that HTTP and TCP share the same maturity level
- that a plausible walkthrough is the same thing as a verified drill

## The target contract

The repo's clearest target request contract is still the one named in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That target contract matters because it is stricter than:

- "DNS can hit more than one box"
- "Traefik is running on multiple nodes"
- "the request can probably be proxied somewhere"

The contract is not about reachability alone.
It is about preserved request meaning after locality fails.

## The live request actors already present in the priority runtime

The current root runtime already gives us a real request stack to talk about.

From [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and
[`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml),
the live actors already include:

- Cloudflare DDNS:
  - `cloudflare-ddns`
- L7 ingress:
  - `traefik`
- policy-preserving forward-auth layer:
  - `nginx-traefik-extensions`
- identity/auth surface:
  - `tinyauth`
- edge filtering and enforcement:
  - `crowdsec`
- private mesh assumption:
  - `headscale-server` and related Headscale surfaces
- route test helper:
  - `whoami`
- fallback experiment surface:
  - `docker-gen-failover`

The root runtime also proves that actual routed services already exist across
both HTTP and TCP classes:

- HTTP examples:
  - `code-server`
  - `searxng`
  - `wishlist`
  - `whoami`
- TCP examples:
  - `mongodb`
  - `redis`

That means the request-path question is not hypothetical.
The stack already has enough live edge machinery for the missing middle layer to
matter.

## The exact question to keep asking at each hop

At every stage, the real question is:

> what does the receiving node know by shared system-owned truth, and what
> would still have to be privately supplied by the operator if the local
> backend were absent?

If that question is skipped, a walkthrough can still sound technical while
quietly collapsing back into theater.

## Literal walkthrough: normal-day stateless HTTP request

Use a service like `wishlist.$DOMAIN` or `whoami.$DOMAIN` as the mental model.

### Hop 1: DNS and public entry

Live evidence:

- the repo uses Cloudflare DDNS in
  [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- the architecture dream explicitly rejects one sacred public reverse-proxy box

What this hop proves:

- the repo wants more than one public entry candidate
- DNS plurality is part of the anti-SPOF story

What it does not prove:

- the receiving node knows where the target service really lives
- the request can preserve meaning if the local backend is absent

### Hop 2: Traefik accepts the hostname

Live evidence:

- Traefik is the real L7 execution surface
- services such as `wishlist` and `whoami` carry real Traefik router/service
  labels

What this hop proves:

- the node can parse a hostname and select a local router definition
- the runtime already has real route declarations, not just design notes

What it does not prove:

- the router has cross-node placement truth
- the node can pick a remote peer correctly from shared current-state knowledge

### Hop 3: middleware and auth are applied

Live evidence:

- `nginx-traefik-extensions` exposes `nginx-auth` forward-auth behavior
- `tinyauth` provides an auth service at `/api/auth/traefik`
- `crowdsec` is wired into Traefik as a real edge filtering layer

What this hop proves:

- the repo already treats request meaning as more than just host-to-container
  forwarding
- auth and middleware continuity are first-class concerns in the live runtime

What it does not prove:

- those same semantics survive a cross-node handoff
- the peer node has the same secrets, auth expectations, and middleware
  assumptions at the decisive moment

### Hop 4: local service execution

If the requested service is actually local to the receiving node, the story is
much stronger.

What this hop can honestly prove:

- the local happy path works for that service class
- edge and app are at least locally connected

What it still does not prove:

- wrong-node dignity
- peer eligibility truth
- backend-loss persistence

The local happy path is real progress.
It is not the user's whole question.

## Literal walkthrough: wrong-node scene

This is the scene the user actually cares about.

Assume:

1. DNS sends traffic to a healthy public node
2. Traefik, CrowdSec, and auth all look respectable
3. the requested service is not actually local to that node

At that exact moment, the live repo is very good at describing the problem and
not yet strong enough to prove the general solution.

### What the node can likely know already

From the live edge stack, the node can likely know:

- the hostname
- the route class
- the middleware chain it would normally apply
- whether the request arrived at a healthy edge surface

### What the repo still does not prove the node knows on its own

The current worktree does not yet prove the node has shared current truth for:

- exact live placement of the target service
- whether the relevant peer is healthy enough to receive the request
- whether that peer has converged secrets, environment, and auth semantics
- whether the peer-forwarded request would still mean the same thing after
  handoff

That is the missing middle layer in one scene.

## Why helper presence still is not enough

The repo already contains things that make the wrong-node story sound closer:

- `docker-gen-failover`
- substantial Traefik dynamic and middleware logic
- Headscale-backed private networking assumptions
- route-rich Compose fragments

But helper presence alone still does not count as wrong-node proof.

The user is specifically tired of:

- "there is a failover helper"
- "there is a private mesh"
- "there is a proxy plugin"
- "there is a route generator"

being narrated as if the hard part had therefore moved out of the operator's
head.

It has not, unless the receiving node can show where it got the truth for the
handoff decision.

## HTTP versus TCP request meaning

This page also needs to keep one hard separation visible.

### Stateless or mostly stateless HTTP routes

The dream here is plausible earlier.

Why:

- route meaning is more visible
- auth and middleware can be observed at L7
- one good drill can prove a narrower but real win

### TCP and stateful routes

The dream is much harsher here.

Examples already present in the root runtime:

- `mongodb`
- `redis`

Why the bar is higher:

- HostSNI/TCP passthrough does not solve authority
- peer-forwarding a TCP stream does not prove stateful correctness
- client rediscovery and write-path truth matter far more

That is why a request walkthrough for `whoami.$DOMAIN` cannot be overpromoted
into a maturity claim for `redis.$DOMAIN` or `mongodb.$DOMAIN`.

## What a real request-path proof packet would need

For one protected stateless HTTP route, a serious proof packet would need:

- exact hostname tested
- receiving node identity
- backend node identity
- source of placement truth used by the receiving node
- evidence that the middleware chain stayed the same
- evidence that auth meaning stayed the same
- evidence that the route survived the exact failure introduced
- a sentence naming what was still not proven outside that route class

Without that packet, a clean walkthrough is still analysis, not runtime-owned
request preservation.

## What still does not count as a real request-path answer

These still fail the user's standard:

- naming Cloudflare, Traefik, TinyAuth, and CrowdSec in order
- showing a service answers locally
- proving the route works on the same node that hosts the backend
- showing more than one DNS record exists
- describing a likely peer-forward design without proving shared placement
  truth

Those are all useful.
None of them answers the user's actual accusation:

> does the healthy wrong node know enough on its own, or am I still the hidden
> registry and route explainer?

## Bottom line

The live request path is already serious enough to be worth tracing:

- Cloudflare DDNS is real
- Traefik is real
- TinyAuth and nginx auth are real
- CrowdSec is real
- Headscale is real
- HTTP and TCP routes are real

What is still not generally proven is the handoff scene that matters most:

the healthy wrong node preserving the request from shared truth instead of
forcing the operator to privately finish the logic.

That is why this page should be read less like a diagram and more like a
custody audit:

at which hop does the request stop being system-owned and start depending on
one human remembering what the stack really means?
