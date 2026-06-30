# Operator Questions and Honest Answers

This page exists because `bolabaden-infra` is not mainly short on components.
It is short on answers that still feel true after a request lands on the wrong
node, a preferred backend disappears, or a stateful service answers one more
time than it deserves to.

That is the standard here.

The answer must still feel honest after the bad day begins.

## What this page is and is not allowed to prove

This page is allowed to:

- restate the operator's real questions in the sharper form the repo actually
  needs
- answer those questions from the strongest evidence the repo currently has
- explain why the common neighboring answer still stops one layer too early
- name the next missing artifact, drill, or proof packet required for a
  stronger sentence

This page is not allowed to:

- narrate target architecture as if it were already runtime proof
- confuse first-hop plurality with request preservation
- confuse route execution with route truth
- use bigger orchestration nouns as a substitute for burden transfer
- make the repo sound more solved than the worktree can currently support

This page is a burden ledger.
It is not a confidence theater page.

## The real question behind all the smaller questions

The operator is not fundamentally asking:

- which orchestrator exists
- which proxy is most mature
- which HA product sounds most enterprise
- which cluster stack is most fashionable this year

The operator is asking:

> how do several ordinary Docker nodes become one believable
> request-preserving platform while
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> remains the real authoring surface and one human stops being the hidden
> registry, hidden peer selector, hidden fallback explainer, and hidden memory
> of what lives where?

That is why ordinary option lists keep feeling smaller than they sound.

## The answer discipline used on this page

Every serious answer here should leave behind five things:

1. the hidden burden the operator is actually naming
2. the strongest current evidence class behind the answer
3. why the nearby common answer is still too small
4. the next artifact or drill that would allow a stronger sentence
5. the private sentence the operator still has to finish alone today

If an answer does not leave those five things behind, it is still too close to
shopping advice.

## Strongest honest current summary

The repo is not mainly suffering from lack of tooling.
It is suffering from lack of smaller honest control surfaces that move
decisive bad-day truth out of private operator memory without immediately
forcing the whole system into Swarm, Nomad, k3s, Kubernetes, or another
larger worldview that has not yet proved it deserves to hide that much truth.

That is why the ecosystem can feel full while still failing to feel useful.

## Question 1: What is the user actually trying to make true?

### Hidden burden

The operator is still privately carrying all of these:

- placement memory
- peer eligibility judgment
- fallback interpretation
- route meaning continuity
- stateful caveat memory

### Strongest current evidence class

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  states the target operating contract directly
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
  confirms the priority implementation is still Compose-first
- the current runtime already proves a real edge stack, real stateful
  services, a private mesh layer, and multi-fragment Compose reality
- the knowledgebase now keeps first hop, wrong-node, backend-loss, and
  stateful truth separated instead of blending them

### Why the nearby common answer is still too small

The user is not merely asking for:

- better load balancing
- anti-SPOF in the vague sense
- a more mature cluster
- fewer Compose files

The real target is harsher:

> multiple ordinary Docker nodes should start behaving like one
> request-preserving personal cloud without the operator having to privately
> finish the sentence about what the request was supposed to mean.

### What would allow a stronger answer

- one shared placement-truth surface actually consumed by routing or forwarding
- one narrow wrong-node HTTP proof
- one backend-loss proof that shows the rescue path still means the same thing

### Private sentence still surviving today

> yes, but I still personally know what the request should have meant better
> than the system does

## Question 2: Why do ordinary HA answers keep feeling fake here?

### Hidden burden

The operator still has to privately know one or more of:

- where the service really lives
- which peer is actually safe
- whether the rescue route survives backend loss
- whether a protected forwarded request still means the same thing
- whether a reachable stateful service is merely answerable rather than
  honestly movable

### Strongest current evidence class

- the live runtime already contains Traefik, TinyAuth, CrowdSec,
  `docker-gen-failover`, Headscale, WARP routing, metrics, and stateful
  services
- the archive-pressure pages repeatedly show the complaint is not lack of
  products, but lack of burden-transferring options

### Why the nearby common answer is still too small

Common answers such as:

- `point Cloudflare at more nodes`
- `just use Traefik`
- `add healthchecks`
- `use service discovery`
- `move to Kubernetes`

often improve one visible layer while leaving the decisive topology sentence
privately owned.

### What would allow a stronger answer

- one route-specific proof showing exactly what truth moved from the operator
  into the system

### Private sentence still surviving today

> yes, but I still personally know the real answer when the platform is
> supposed to act coherent

## Question 3: Is Traefik the answer to the multi-node problem here?

### Hidden burden

Routing execution is being confused with routing truth.

### Strongest current evidence class

Traefik is materially live in the priority runtime and fronts real surfaces.
It already sits beside:

- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- file-provider config
- TCP routers for services such as `mongodb`, `redis`, and
  `biodecompwarehouse*`

### Why the nearby common answer is still too small

Traefik absolutely buys real things:

- TLS termination
- HTTP routing
- middleware execution
- auth integration
- TCP exposure

Traefik does **not** by itself buy:

- shared current placement truth
- peer-eligibility truth
- backend-loss route persistence
- honest stateful authority transfer

The repo already knows how to expose a route.
The repo is trying to stop privately narrating what that route means.

### What would allow a stronger answer

- a receiving node using shared placement truth to choose local versus remote
- one protected route comparison between local execution and peer-forwarded
  execution

### Private sentence still surviving today

> yes, but I still personally know whether Traefik's next hop is actually the
> right peer

## Question 4: Why is Cloudflare not the answer by itself?

### Hidden burden

Plural first-hop reachability is being confused with preserved request meaning.

### Strongest current evidence class

- Cloudflare is part of the explicit public-entry philosophy
- `cloudflare-ddns` is live in the current edge stack
- multi-record public entry is clearly part of the repo's anti-SPOF pressure

### Why the nearby common answer is still too small

Cloudflare can help with:

- more than one public record
- first-hop resilience
- exposure management

Cloudflare cannot, by itself, tell the receiving node:

- whether the service is local
- where the real backend lives now
- which peer is eligible for this route
- whether the rescue path still survives the failure that made it necessary

### What would allow a stronger answer

- one wrong-node route proven after intentionally landing on a non-owner node

### Private sentence still surviving today

> yes, but I still personally know that DNS plurality did not solve the real
> request-preservation problem

## Question 5: Why does Headscale not solve service discovery by itself?

### Hidden burden

Reachability and identity are being mistaken for placement truth and peer
validity.

### Strongest current evidence class

- Headscale is materially live in
  [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- the active fragment still uses SQLite at `/var/lib/headscale/db.sqlite`
- the mesh layer clearly matters to the intended peer-forward model

### Why the nearby common answer is still too small

Headscale can give the repo:

- private-node connectivity
- stable identity
- real mesh plumbing

It does not by itself answer:

- which node currently owns the requested service
- which peer is valid for this exact route
- whether the control plane itself has stopped being socially singular

### What would allow a stronger answer

- a routing or placement layer that consumes peer identity plus current
  service ownership

### Private sentence still surviving today

> yes, but I still personally know that a reachable peer is not yet a proven
> valid backend

## Question 6: Why does `services.yaml` keep reappearing?

### Hidden burden

The operator keeps acting as the safest current-state registry in the system.

### Strongest current evidence class

- the docs repeatedly converge on `services.yaml`
- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  explicitly names lightweight current-state-registry philosophy
- the current runtime still does not prove a tracked root `services.yaml`
  consumed by route selection logic

### Why the nearby common answer is still too small

Without some shared placement-truth surface, every other answer remains
undersized:

- Cloudflare only gets traffic onto a node
- Traefik only executes the routing graph it already knows
- Headscale only makes peers reachable
- helpers only look convincing until the wrong-node question is asked

### What would allow a stronger answer

- one tracked placement registry or equivalent truth surface consumed by
  routing, forwarding, or selection logic in the priority runtime

### Private sentence still surviving today

> yes, but I still personally know what runs where better than the system does

## Question 7: Why are helpers like `docker-gen-failover` still potentially fake comfort?

### Hidden burden

A helper can reduce repetition without absorbing the exact failure that makes
it matter.

### Strongest current evidence class

- `docker-gen-failover` is materially present in the live edge fragment
- the current generated file path is
  `/traefik/dynamic/failover-fallbacks.yaml`
- the docs already record that helper-generated routes can disappear exactly
  when the preferred backend disappears

### Why the nearby common answer is still too small

A helper often looks like relief because it:

- generates config
- reduces manual edits
- reacts to events

But if the helper fails during the exact backend-loss condition it is meant to
absorb, then it did not move the real burden.

### What would allow a stronger answer

- a backend-loss drill showing that the rescue route still exists and still
  means the same thing after the failure

### Private sentence still surviving today

> yes, but I still personally know the helper does not yet survive the exact
> bad day I care about

## Question 8: Why are stateful services treated so much more harshly?

### Hidden burden

Write authority, replication truth, promotion flow, reconnect behavior, and
rediscovery truth are still the real failure domains.

### Strongest current evidence class

The live runtime already contains:

- root `mongodb`
- `redis`
- Headscale SQLite
- `nuq-postgres`
- `litellm-postgres`
- `rabbitmq`
- `qdrant`

### Why the nearby common answer is still too small

The following do not answer the real question:

- stable hostnames
- TCP exposure
- successful healthchecks
- restartability
- `we can move it later`

Those can improve reachability and operations without changing who owns truth.

### What would allow a stronger answer

- for each stateful class: explicit write owner, replica model, promotion
  flow, reconnect expectations, and rediscovery behavior

### Private sentence still surviving today

> yes, but I still personally know the real writer is singular

## Question 9: Why is "just use Nomad / Swarm / k3s / Kubernetes" still not a complete answer?

### Hidden burden

The repo needs burden transfer, not merely a more adult-sounding control
plane.

### Strongest current evidence class

- the archive-pressure material repeatedly explores Swarm, Nomad, k3s,
  Kubernetes, OpenSVC, and related directions
- the repo's honesty surfaces explicitly say heavier control layers must earn
  themselves

### Why the nearby common answer is still too small

A stronger orchestrator might eventually deserve adoption.
But the repo is not asking:

- which orchestrator exists
- which orchestrator is popular
- which orchestrator is technically capable in the abstract

It is asking:

> which option removes the humiliating moment where a healthy node still needs
> a human to explain what the request should have meant?

If the answer does not name that transferred burden, it is still too broad.

### What would allow a stronger answer

- a burden-by-burden promotion matrix showing what truth a candidate layer
  would actually own

### Private sentence still surviving today

> yes, but I still do not know whether this solves my problem or just replaces
> it with a bigger worldview

## Question 10: What is the most useful next proof to chase?

### Hidden burden

The repo still lacks one humiliatingly concrete proof that a receiving node
can act correctly without private operator completion.

### Strongest current evidence class

The runtime already has:

- simple stateless HTTP candidates such as `whoami`, `wishlist`, and `mkdocs`
- real edge policy surfaces
- real wrong-node architectural pressure

### Why the nearby common answer is still too small

`pick a more mature platform` is too broad until the repo first proves what
burden is actually being moved.

### What would allow a stronger answer

The best next proof packet is still:

1. expose one shared placement-truth surface
2. intentionally land traffic on the wrong healthy node
3. prove one stateless HTTP route still completes correctly
4. then kill the preferred backend and prove whether the same route survives

### Private sentence still surviving today

> yes, but I still personally know the system has not yet passed the
> humiliating wrong-node test

## What still does not count as an honest answer

These still do not count:

- naming a product without naming which burden it would actually own
- answering `why are there no real options?` with a tool list
- treating first-hop plurality as request preservation
- treating clearer prose as stronger evidence
- recommending a controller without naming which private sentence it would
  kill

That is the protection mechanism for this page.

## The packet a genuinely useful answer should leave behind

After reading any serious answer in this repo, the operator should know:

- which truth is still private
- why the obvious nearby answer stops one layer too early
- what next artifact would externalize that truth
- what sentence remains forbidden until that artifact exists

If the answer leaves only a better shopping list, it failed.
