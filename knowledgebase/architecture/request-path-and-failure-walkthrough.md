# Request Path and Failure Walkthrough

This page exists because the user's real question is not:

> what components are in the ingress stack?

It is:

> when one hostname is requested, what literally happens next, which files own
> each decision, what does the receiving node know on its own, and where does
> the story stop being system-owned and turn back into private operator memory?

That is a much harder question than "what reverse proxy do we use?"
It is also the question the repo keeps circling every time ordinary HA answers
start feeling fake.

## What this page is and is not allowed to prove

This page is allowed to prove:

- what the target request contract is
- which live runtime components already participate in that contract
- where the request path is already concrete in the active Compose runtime
- where present proof stops for wrong-node and backend-loss behavior

This page is not allowed to prove:

- generic wrong-node success today
- generic peer-forward fallback today
- that HTTP and TCP share the same maturity level
- that a plausible walkthrough is the same thing as a verified drill

This is a request-path decomposition page.
It is not a drill result.

## The target contract

The clearest target request contract is still the one named in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That contract matters because it is stricter than:

- `more than one public node exists`
- `Traefik is running`
- `the route can probably be proxied somewhere`

The contract is not about reachability alone.
It is about preserved request meaning after locality fails.

## The live request actors already present now

The current root runtime already gives us a real chain to talk about.

From [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
and
[`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml),
the live request actors already include:

- public-entry updater:
  - `cloudflare-ddns`
- main L7 execution surface:
  - `traefik`
- policy-preserving forward-auth layer:
  - `nginx-traefik-extensions`
- auth provider:
  - `tinyauth`
- edge filtering and enforcement:
  - `crowdsec`
- helper with fallback intent:
  - `docker-gen-failover`
- private mesh assumption:
  - `headscale-server` and related Headscale surfaces
- route test helper:
  - `whoami`

The wider runtime also proves that real routed services already exist across
multiple service classes:

- simple stateless HTTP candidates:
  - `wishlist`
  - `whoami`
  - `mkdocs`
- richer HTTP surfaces:
  - `searxng`
  - `code-server`
  - `open-webui`
  - `gptr`
- TCP/stateful surfaces:
  - `mongodb`
  - `redis`
  - `rabbitmq`
  - `nuq-postgres`

So the request-path question is not hypothetical.
The runtime is already broad enough that the missing middle layer matters.

## The one question to keep asking at every hop

At every stage, keep asking:

> what does the receiving node know from shared system-owned truth, and what
> would still have to be privately supplied by the operator if the local
> backend were absent?

If that question disappears, the walkthrough can stay technical while still
quietly drifting back into theater.

## Literal walkthrough: normal-day stateless HTTP request

Use `wishlist.$DOMAIN` or `whoami.$DOMAIN` as the mental model for the cleanest
first pass.

These are useful because they are:

- HTTP
- comparatively simple
- good proof candidates
- less likely to hide stateful ambiguity

### Hop 1: DNS and public entry

Live evidence:

- Cloudflare DDNS is part of the active edge fragment
- the repo explicitly rejects one sacred public reverse-proxy box

What this hop honestly proves:

- first-hop plurality is part of the live design
- more than one node is intended to be a plausible public entry candidate

What it still does not prove:

- the receiving node knows where the target service actually lives now
- the receiving node can preserve meaning if the service is not local

This is where many explanations stop too early.
DNS plurality is necessary.
It is not the same thing as request preservation.

### Hop 2: Traefik accepts the hostname

Live evidence:

- Traefik is materially present now
- `wishlist`, `whoami`, `searxng`, `code-server`, and many other surfaces carry
  real router/service labels

What this hop honestly proves:

- the receiving node can parse the hostname
- the receiving node can match a local router definition
- the route exists in the authored runtime

What it still does not prove:

- the router owns cross-node placement truth
- the receiving node can select a remote peer from shared current truth

This hop is route execution.
It is not yet route truth.

### Hop 3: middleware and auth are applied

Live evidence:

- `nginx-traefik-extensions` exposes forward-auth behavior used by Traefik
- `tinyauth` provides real auth participation
- `crowdsec` is wired into the edge stack

What this hop honestly proves:

- the repo already treats request meaning as more than host-to-container
  forwarding
- auth, trust boundaries, and middleware continuity are live concerns

What it still does not prove:

- that the same semantics survive peer forwarding
- that the remote peer has the same secrets, middleware assumptions, and auth
  expectations at the decisive moment

This is why protected HTTP is stricter than a plain stateless route.

### Hop 4: local service execution

If the requested service is actually local to the receiving node, the story is
much stronger.

What this hop can honestly prove:

- the local happy path works for that route class
- edge and app are at least locally connected
- the service is not merely theoretical

What it still does not prove:

- wrong-node dignity
- peer-eligibility truth
- backend-loss persistence

The local happy path is real progress.
It is not the user's whole question.

## Literal walkthrough: protected HTTP request

Now use a route like `code-server.$DOMAIN`, `grafana.$DOMAIN`, or `mcpo.$DOMAIN`
as the mental model.

These are more revealing because the route has to preserve not just transport,
but meaning:

- auth behavior
- middleware order
- trust boundaries
- any forwarded headers or policy assumptions

### What the runtime already gives this class

The current stack already gives protected routes:

- Traefik router definitions
- auth participation through `nginx-auth@file` and TinyAuth
- middleware surfaces
- observability through logs and dashboards

### What still remains unproven

The current worktree still does not prove:

- that a protected route forwarded to another node still means the same thing
- that the peer has converged enough secrets and config to preserve that meaning
- that the same route survives backend loss with the same visible contract

This is why protected HTTP cannot inherit confidence from a simple `whoami`
success.

## Literal walkthrough: wrong-node scene

This is the scene the user actually cares about.

Assume all of the following:

1. DNS sends traffic to a healthy public node
2. Traefik, CrowdSec, and auth all look healthy
3. the requested service is not actually local to that node

At that exact moment, the repo is very good at describing the problem and not
yet strong enough to prove the general solution.

### What the receiving node can likely know already

From the live edge stack, the receiving node can likely know:

- the hostname
- the route class
- the middleware chain it would normally apply
- whether the edge stack itself looks healthy

### What the current worktree does not yet prove the node knows on its own

The current worktree does not yet prove the receiving node has shared truth
for:

- exact live placement of the target service
- whether the relevant peer is healthy enough to receive the request
- whether that peer has converged secrets, environment, and auth semantics
- whether the forwarded request would still mean the same thing after handoff

That is the missing middle layer in one scene.

## Why helper presence still is not enough

The repo already contains several things that make the wrong-node story sound
closer to solved:

- `docker-gen-failover`
- substantial Traefik dynamic and middleware logic
- Headscale-backed private networking assumptions
- route-rich Compose fragments

But helper presence alone still does not count as wrong-node proof.

The user is specifically tired of:

- `there is a failover helper`
- `there is a private mesh`
- `there is a route generator`
- `there is a dynamic proxy`

being narrated as if the hard part had therefore moved out of the operator's
head.

It has not, unless the receiving node can show where it got the truth for the
handoff decision.

## Backend-loss scene

There is a second trap that is easy to hide inside a clean walkthrough:

> even if a wrong-node handoff looks plausible while everything is healthy, the
> rescue path may still disappear once the preferred backend actually fails.

This matters because a route generator or helper can look excellent right up
until the exact condition it is meant to absorb.

The repo already records this pressure around `docker-gen-failover`.

So backend-loss needs to stay separate from wrong-node success.

Wrong-node success means:

- the request landed on a non-owner node
- the system still chose a valid backend

Backend-loss success means:

- the preferred backend disappeared
- the rescue path still existed
- the route still meant the same thing
- the operator did not have to privately complete the story

Those are related.
They are not the same proof.

## HTTP versus TCP request meaning

This page also has to keep one hard separation visible.

### Stateless or mostly stateless HTTP routes

The dream is most plausible here first.

Why:

- route meaning is visible at L7
- middleware and auth can be inspected
- one good drill can prove a narrow but real win

This is why early proof candidates should usually come from:

- `whoami`
- `wishlist`
- `mkdocs`

### Protected HTTP routes

This class is stricter.

A `200` is not enough.
The question is whether the same protected service meaning survives handoff.

Good future candidates include:

- `code-server`
- `grafana`
- `mcpo`

### Raw TCP routes

This class is different enough that HTTP optimism cannot spill into it.

A TCP router for `mongodb` or `redis` does **not** prove:

- peer-aware service failover
- semantic continuity
- safe substitution between backends
- stateful correctness

It proves transport exposure, not full replacement safety.

### Stateful surfaces

This is the strictest class of all.

A stateful service does not become honestly HA because:

- it is reachable through Traefik
- more than one node exists
- a proxy can point to more than one backend

Stateful routing becomes meaningful only when authority, durability,
substitution safety, promotion, reconnect, and rediscovery are explicit.

## The exact place the story turns back into operator memory

The story turns back into private operator memory the moment the runtime cannot
answer all of these on its own:

- where does this service really live right now?
- is the candidate peer just alive, or valid for this route?
- if forwarded, will the route still mean the same thing?
- if the preferred backend disappears, what path still exists?
- if the route is stateful, who is the actual authority?

That is the real SPOF reading in this repo.

The user is not only worried about a dead server.
The user is worried about the operator still being the hidden explainer when
the system was supposed to behave like one cloud.

## What a real walkthrough-backed proof packet would need

Before this page can support stronger language, it needs a proof packet with:

- the exact hostname
- the receiving node identity
- the actual backend identity
- the source of placement truth used for the decision
- the source of peer-eligibility truth used for the decision
- evidence that policy/auth meaning stayed the same
- whether the route also survived backend loss
- the sentence that is still forbidden even after the drill

Without that packet, a walkthrough is still mostly disciplined imagination.

## Bottom line

The current runtime already contains a serious request stack:

- Cloudflare-oriented public entry
- Traefik routing
- TinyAuth and forward-auth participation
- CrowdSec enforcement
- Headscale mesh assumptions
- helper-driven fallback intent
- a mix of simple HTTP, protected HTTP, TCP, and stateful services

That is enough to make the request-preservation problem real.

What the current worktree still does not prove is the part the user actually
cares about:

- that the wrong healthy node already knows what the request should mean
- that it already knows which peer is valid
- that the route still means the same thing after handoff
- that the rescue path survives backend loss
- that TCP and stateful routes have escaped operator folklore

That is the honest boundary today.
The runtime can already execute a lot.
The missing win is still the point where the receiving node can explain itself
without a human finishing the sentence first.
