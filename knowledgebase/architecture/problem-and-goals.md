# Problem, Pressure, and Goals

This page defines the actual architecture problem in `bolabaden-infra`.

The problem is not:

- "more clustering"
- "better deployment hygiene"
- "more modern infrastructure"
- "better docs about self-hosting"

The problem is:

> how do several ordinary Docker nodes become one request-preserving personal
> cloud, while `docker-compose.yml` stays readable and the system stops leaning
> on one operator to remember where the real topology truth lives?

That question is stricter than generic HA language, and smaller than a
"migrate to Kubernetes" story.

## What this page is and is not allowed to prove

This page is authoritative about:

- the real problem the repo is trying to solve
- the concrete requirement stack implied by that problem
- which adjacent answers are still too small

This page is not authoritative about:

- whether the current runtime already satisfies those requirements
- whether one future helper layer has already won
- whether naming the problem cleanly means the remaining gap is small

This page is the benchmark, not the completion report.

## Strongest honest current answer

The repo already proves a serious Compose-first platform with a strong edge,
many services, and real planning work around failover and anti-SPOF behavior.

What it does not yet prove is the harder thing:

- that wrong-node entry is survivable generically
- that service placement truth is shared explicitly
- that fallback routes survive backend loss
- that peer forwarding preserves the same auth, middleware, and routing meaning
- that stateful services have honest failover semantics

That gap is the actual problem. Everything else is supporting detail.

## The shortest exact problem statement

The repo is trying to solve this:

> keep Compose as the main authoring and operator surface, but add just enough
> shared truth that a request landing on the wrong healthy node does not turn
> into guesswork, folklore, redirects, or fake failover.

The repo is therefore not only a hosting system.
It is a search for a smaller honest control surface.

## What still does not count as understanding the problem

This page needs its own anti-flattening filter.

The following still do not count as having reconstructed the user's actual
problem:

- rewriting "anti-SPOF" in more polished language
- treating the problem as general self-hosting complexity
- treating the problem as mostly about choosing between Compose and a scheduler
- describing the topology pressure without naming private operator
  reconstruction as the wound
- assuming that more flexibility or more options are automatically relief

Those readings all stay too generic.
The user's complaint is narrower and harsher:

> too many supposedly serious answers still leave the operator as the real
> keeper of placement, fallback, and semantic continuity truth.

## What does not count as solving this problem

The problem is specific enough that it needs an explicit false-solution filter.

The repo has **not** solved the problem merely because:

- more public nodes exist
- DNS can hit several addresses
- a reverse proxy can see more containers
- a mesh network connects the nodes
- a helper can generate routes on the happy path
- a larger orchestrator can be demoed in isolation

Those may all be helpful ingredients.
None of them, by themselves, prove that request meaning survives wrong-node
entry without private operator reconstruction.

## The hidden enemy

The hidden enemy is not lack of products.

The hidden enemy is private operator reconstruction.

That reconstruction currently shows up in questions like:

- what runs where right now?
- if this hostname hits node B, does node B know the service is actually on
  node A?
- if node A is unhealthy, what route survives and who generated it?
- if a helper proxy forwards the request, do auth and middleware still mean the
  same thing?
- if Redis, MongoDB, Headscale, or another stateful service moves, what still
  owns write authority?

As long as those questions are answered mainly from memory, the platform is
still only partially system-owned.

## What the current stack already gives

The repo already has meaningful assets:

- a real root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active include fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- a public edge layer centered on Traefik, CrowdSec, TinyAuth, and
  `nginx-traefik-extensions`
- Cloudflare participation for multi-node public entry pressure
- Headscale for private-mesh coordination assumptions
- observability through Grafana, Prometheus, VictoriaMetrics, Loki, Promtail,
  cAdvisor, Blackbox, and Alertmanager
- explicit planning around service failover, secret sync, compose sync, and
  missing middle-layer helpers

Those are real strengths.

They are still weaker than the final requirement stack because most of them are
ingredients, not shared truth.

## What a serious problem-definition proof packet would have to contain

If this page ever supports the claim that the repo has correctly identified the
real problem rather than a neighboring one, it should be because the
documentation leaves behind a proof packet.

That packet should include:

- the exact hidden burden being named
- the route or service class where that burden becomes visible
- the system-owned truth that is missing today
- the human reconstruction step that still fills the gap today
- the narrower artifact or drill that would prove the burden has moved
- the explicit statement of which stronger burdens still remain

Without that packet, even a sharp problem statement can still drift back into
"interesting infra essay" territory.

## Requirement stack implied by the dream

If the dream is taken seriously, the system needs all of these:

1. **Any-node first hop**
   A request can land on any surviving public node.

2. **Local-first handling**
   If the requested service is local, the node serves it locally instead of
   pretending locality does not matter.

3. **Explicit placement truth**
   If the service is not local, the receiving node has a trustworthy current
   source for where it really lives.

4. **Peer eligibility truth**
   The receiving node can tell which peer is healthy and semantically eligible
   now, not just historically configured.

5. **Fallback-route durability**
   The route needed for rescue remains alive under the failure that made rescue
   necessary.

6. **Policy preservation**
   Auth, middleware, routing policy, and visible service meaning survive the
   handoff.

7. **Stateful honesty**
   Redis, MongoDB, Headscale, databases, and other state-bearing systems are
   described by their real authority and failover semantics, not by mere
   reachability.

8. **Inspectable ownership**
   An operator can explain why the request succeeded by reading tracked shared
   truth, not by privately reconstructing the topology.

If a design satisfies only `1` and `2`, it is still not enough.
If it satisfies `1` through `6` but cheats on `7`, it is still not enough.
If it satisfies all of those but still fails `8`, the system is still leaning
on folklore.

That is why the problem statement must stay stricter than generic anti-SPOF
language.
Without the full stack above, "resilience" becomes too easy to say and too hard
to trust.

## Why generic options lists still fail

Many neighboring answers are technically respectable and still too small here.

### "Use more DNS"

Cloudflare multi-record DNS helps with first-hop plurality.
It does not tell the wrong node where the service is, whether the peer is
healthy, or whether policy survives the handoff.

### "Use a reverse proxy"

Traefik is a real edge asset and clearly central to this repo.
Local container discovery is still not the same thing as cross-node placement
truth.

### "Use a helper generator"

The repo's own planning docs call out `docker-gen-failover` as broken because
it can delete routes when containers stop. A helper that removes the rescue
path under failure is part of the problem, not proof of the solution.

### "Use Kubernetes / k3s / Nomad / OpenSVC"

These may become valid promotions later. They are not automatically the answer
just because they are larger or more respectable.

The repo's standard is harsher:

- what truth do they own?
- what burden do they really remove?
- what new worldview tax do they impose?
- do they solve the exact wrong-node and anti-folklore problem, or mostly
  relocate it?

That harsher standard is what keeps this repo from collapsing back into the
same two humiliations the user keeps rejecting:

- static glue that still needs private human topology memory
- heavyweight control planes adopted before smaller honest answers were
  exhausted

## The practical goal

The practical goal is not "be more like cloud-native infrastructure."

It is:

> externalize enough placement, health, routing, and policy truth that the bad
> day stops depending on sacred remembered nodes and private rescue knowledge.

That is the real acceptance bar for every future control-plane, agent, helper,
registry, or orchestrator decision in this repo.

For the exact acceptance bar that follows from this page, continue to
[Operator Contract and Success Criteria](./operator-contract.md).
