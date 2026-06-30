# HA, Failover, and Routing

Read this page as the routing truth map for the priority implementation rooted
at [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

If you want the deeper evidence stack first, read:

- [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
- [Request Path and Failure Walkthrough](request-path-and-failure-walkthrough.md)
- [Current Compose Runtime](current-compose-runtime.md)
- [Operator Contract and Success Criteria](operator-contract.md)

This page exists because "HA" becomes meaningless almost immediately in this
repo unless routing is decomposed into the separate truths the user is actually
angry about.

That anger should stay specific.
The user is not mainly furious that routing is hard.
They are furious that routing stories keep sounding complete one layer before
the receiving node can actually explain, in shared system terms, why the
request is still safe on the wrong machine.

The grounding for that is concrete, not hypothetical.

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
  states the target contract directly:
  any surviving public node should be able to receive the request, serve it
  locally if the service is local, or forward it to a healthy peer if it is
  not.
- [`docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/chatgpt-exports/conversations/docker-multi-node-without-swarm__68a916ef-b554-832a-aa13-dee8b95de50f.md)
  captures the same pressure in plainer words:
  manual placement is acceptable, DNS plurality already exists, and the real
  problem is service discovery and request preservation when traffic lands on
  the wrong node.
- [`load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/source-archive/chatgpt-exports/conversations/load-balancer-failover-alternatives__68252e5b-7218-8006-8857-2e46d731e299.md)
  shows the user explicitly rejecting answers that stop at familiar proxy
  nouns.

The user is not asking:

- can more than one node receive traffic?
- is Traefik present?
- are there healthchecks?
- is there a helper that sounds failover-shaped?

The user is asking:

> if traffic lands on the wrong healthy machine, can that machine still
> preserve the service contract without requiring private operator memory?

That is the routing question here.

The reason the wording has to stay that severe is that the user is not merely
frustrated with downtime.
They are frustrated with how many respectable infrastructure options solve one
routing layer loudly while quietly leaving the operator as the hidden
cross-node interpreter.

That last phrase is the real routing wound in this repo:
the operator remains the final translator between node entry, placement truth,
policy continuity, and fallback legitimacy.

## What this page is and is not allowed to prove

This page is authoritative about:

- how routing has to be decomposed before "failover" means anything useful
- which routing layers are materially live in the current Compose runtime
- which routing truths are still missing, social, or only planned
- why first-hop plurality is much weaker than request-preserving recovery
- why HTTP, TCP, and stateful failover must not be flattened together

This page is not authoritative about:

- whether a specific hostname has already passed a real wrong-node drill
- whether backend-loss fallback is already durable for a named route
- whether stateful traffic inherited optimism from the HTTP story
- whether one working helper path upgrades the whole platform

This is a routing decomposition page, not a route-success report.
It also cannot become a reassurance page.

The user is not mainly asking for cleaner routing terminology.
They are asking why so many HA and failover stories keep sounding plausible
until the request lands on the wrong healthy node and the operator is still
the only thing that knows what should happen next.

That means this page has to keep one humiliating question alive all the way
through:

> after the first hop lands on the wrong node, what exact truth does that node
> itself own, and what exact truth is still being borrowed from the operator?

## Strongest honest current answer

The root runtime already has a serious edge stack and enough moving parts to
make the routing problem real rather than hypothetical:

- Cloudflare-oriented public-entry assumptions
- Traefik
- TinyAuth
- `nginx-traefik-extensions`
- CrowdSec
- Docker-socket proxies
- `docker-gen-failover`
- Headscale mesh assumptions
- TCP routers for services such as MongoDB and Redis

What is still missing is the thing the user actually wants that machinery to
buy:

- shared placement truth
- trustworthy peer eligibility truth
- durable fallback-route persistence
- proof that the request keeps meaning the same thing after wrong-node handoff
- proof that stateful authority survives any routing story being told

That is the whole difference between:

- a stack that can route many things locally
- and a stack that stops gambling on node locality

That second line is closer to the user's dream than ordinary HA vocabulary.
The dream is not just more reachability.
The dream is that wrong-node entry stops being the moment where the system
reveals it was still counting on human folklore.

If a routing summary sounds sophisticated while still leaving wrong-node entry
as the moment where the operator privately completes the story, then the page is
still rewarding the wrong thing.

That is the anti-theater rule for this page:
if the receiving node still depends on human folklore to know who should serve,
who is eligible, and whether the handoff preserves the same protected service,
the routing layer is still being socially completed.

## What still does not count as HA or failover here

This page should make the common overreads illegal.

The following still do not count as meaningful HA or failover in this repo:

- more than one public node can receive the first hop
- Traefik is present and healthy
- a helper generates fallback-shaped route material
- a TCP router exists for a stateful service
- a local protected route returns `200`
- a mesh exists between nodes

All of those may be real progress.
None of them are yet the user's actual benchmark unless they also reduce the
need for private placement memory, preserve request meaning on the wrong node,
and survive the failure that made fallback necessary.

This is where normal HA vocabulary becomes actively misleading.

Many ecosystems will happily call the above:

- resilience
- failover
- high availability
- distributed ingress

The repo has to keep asking the harsher follow-up:

- yes, but who still had to know the real answer first?

That is not rhetorical flourish.
It is the shortest honest checksum for whether "HA routing" has become one more
label for partial machinery plus private operator settlement.

They also do not satisfy the deeper complaint:

> I am tired of options that solve one routing layer and leave the next layer
> as my private job.

That is why this page has to keep decomposing routing into narrower truths
instead of letting "HA routing" sound like one solved category.

The ecosystem habit this page is resisting is simple:

- solve one visible routing layer
- let the label "HA" drift upward
- leave the next layer as the user's private job

## What a real routing proof packet would have to contain

If this page ever supports stronger routing claims, it should be because an
actual route-level proof packet exists.

That packet should include:

- the exact route class, such as stateless HTTP, protected HTTP, raw TCP, or a
  named stateful surface
- the entry node and backend node identities
- the source of placement and peer truth used for the handoff
- the failure condition introduced, if fallback is being claimed
- the policy or middleware comparison, if semantic continuity is being claimed
- the explicit statement of which stronger routing class is still unproven

Examples:

- a protected HTTP route can prove wrong-node dignity for that route class
  without proving TCP or stateful correctness
- a TCP route can prove transport continuity without proving write authority
- a stateful drill can prove one exact topology without upgrading the whole
  platform into generic HA

Without packets like those, routing prose is still too easy to overread as
stack-wide resilience.

The packet requirement is stricter here because the user is already surrounded
by too many plausible routing stories that fail exactly when they finally
matter.

What the user is starved for is not another routing recipe.
It is one believable packet that proves the request remained the same service
after locality failed, and that the explanation for that survival lived in the
system rather than in one operator's head.

This is also why route-class separation matters so much.

The repo cannot afford summaries that do this:

- one HTTP success becomes "failover works"
- one TCP router becomes "stateful is covered too"
- one helper artifact becomes "the missing middle is basically here"

Those are exactly the overreads the user is already exhausted by.

## The current routing surface in the priority implementation

The priority root runtime is the merged surface created by:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
- [`compose/docker-compose.docs.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.docs.yml)
- [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- [`compose/docker-compose.metrics.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.metrics.yml)

At the root level, the merged graph already defines three important network
surfaces:

- `publicnet`
- `backend`
- `warp-nat-net`

That alone does not prove resilience.
It does prove the routing problem is already materially encoded in the runtime,
not just imagined in planning docs.

The repo is not theorizing from an empty lab.
It is already carrying enough live routing complexity that the missing burden
feels offensive rather than hypothetical.

That distinction matters because the frustration is not beginner confusion.
The user is not staring at an empty repo asking abstract clustering questions.
They are staring at a stack that is already serious enough to make the
remaining hidden burden feel insulting instead of understandable.

Within the edge fragment, the current routing surface includes:

- `cloudflare-ddns`
- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `crowdsec-init`
- `docker-gen-failover`
- `dockerproxy-ro`
- `dockerproxy-rw`

Within the root runtime there are already concrete route classes that stress
routing differently:

- stateless or near-stateless HTTP surfaces such as docs and utility frontends
- protected HTTP surfaces such as Dozzle that depend on auth and middleware
  continuity
- raw TCP surfaces such as MongoDB and Redis already exposed through Traefik
  TCP routers in the root graph
- stateful control-plane surfaces such as Headscale, whose single-node reality
  is already openly acknowledged in the planning layer

That matters because "routing" here is not one problem.

It is also why ordinary "best HA proxy" or "best failover tool" advice is too
small for this repo.
The user is not short on proxies.
The user is short on options that preserve meaning instead of merely
redirecting packets.

That sentence is the whole routing benchmark:
preserve meaning, not just transport.
Everything else on this page is here to stop the docs from quietly forgetting
that.

## The routing layers that must stay separate

If these layers get collapsed into one sentence called "HA," the docs become
decorative again.

They also become emotionally false.
They start describing reassurance instead of burden movement.

### 1. Public node-entry reachability

Question:

- can a client reach some healthy public node at all?

What the repo already has:

- Cloudflare-oriented public-entry assumptions
- `cloudflare-ddns`
- a long-standing goal of multi-node public entry instead of one sacred public
  box
- repo-native intent that explicitly says "User -> Cloudflare DNS -> any
  surviving node"

What this layer can honestly buy:

- more than one node can plausibly receive the first hop
- ingress does not have to be psychologically concentrated on one sacred host
- some node loss at the public edge may be survivable

What this layer does not buy:

- service locality truth
- peer eligibility truth
- fallback-route persistence
- policy continuity
- stateful correctness

This is why "multiple A records" is not even close to the final answer.

It is often the first place people feel tempted to stop thinking because it
looks distributed enough to market.

That temptation is one of the exact pressures this knowledgebase is fighting.
The user has already seen enough marketed distribution.
What is being demanded here is distribution that still tells the truth after
the landing node is wrong.

That last clause is the anti-theater core of the whole routing story.

The page should keep readers from mistaking:

- more nodes
- richer ingress
- more failover nouns

for the much narrower achievement the user is actually looking for:

- a wrong-node machine that can still tell the truth without human rescue

### 2. Local edge-stack health

Question:

- once traffic reaches a node, is the local edge stack coherent enough to make
  the next decision honestly?

What the repo already has:

- Traefik as the real L7 execution surface
- TinyAuth for forward-auth behavior
- `nginx-traefik-extensions` as auth and middleware glue
- CrowdSec as active security and filtering logic
- Docker and file provider surfaces
- health-bearing edge components

What this layer can honestly buy:

- local route execution
- local auth and middleware handling
- local TLS and certificate handling
- serious edge behavior instead of one symbolic proxy

What it does not buy:

- cross-node knowledge
- trustworthy peer choice
- proof that local and forwarded behavior remain semantically identical

If this layer is broken, distributed entry above it is useless.
If this layer is healthy, the distributed problem is still not solved.

This is a recurring pattern across the repo:

- a healthy local edge is necessary
- a healthy local edge is impressive
- a healthy local edge is still not the same thing as wrong-node rescue

This distinction is one of the repo's most important anti-fake-option rules.
Many systems become locally impressive enough that the user gets told the
remaining gap is small.
The user's whole point is that the remaining gap is the actual problem.

This page should not let local sophistication become emotional closure.

That is one of the easiest ways routing documentation starts lying while
remaining technically accurate.

### 3. Locality truth

Question:

- does the receiving node actually know whether the requested service is local?

What the repo already has:

- local Docker labels
- local Docker provider visibility
- route material in the root and edge runtime
- repeated architectural pressure toward a tracked current-state registry such
  as `services.yaml`

What the repo does not yet prove:

- a live tracked root `services.yaml` consumed by routing decisions
- a shared placement-truth surface that outranks operator memory

Why this matters:

The recurring `services.yaml` pressure is not about loving files.
It is about stopping "the operator remembers where the service really lives"
from being the real control plane.

Without locality truth, a wrong-node request cannot become an honest decision.

That is why `services.yaml` keeps haunting the repo.
Not because YAML is fashionable, but because some externalized present truth
has to exist somewhere before the system can stop pretending the operator's
memory is not the real control plane.

### 4. Peer-selection truth

Question:

- if the service is not local, does the receiving node know which peer is
  actually valid right now?

This is stricter than:

- can nodes talk over Headscale?
- does a peer hostname exist?
- did the same service run on that peer recently?

What the repo already has:

- Headscale as a materially live mesh component
- explicit planning around peer broadcast, leader election, and node-aware
  coordination in the master plan

What the repo does not yet prove:

- a live peer-eligibility decision surface consumed by the edge layer
- a trustworthy answer to "which peer should I use now?" that is stronger than
  static config or social memory

Mesh reachability is helpful.
Mesh reachability is not current truth ownership.

That difference is where many "peer-aware" stacks keep cheating.
Peer contact is not peer judgment.
A reachable peer is not yet a peer the receiving node can honestly trust for
this route right now.

That last sentence is one of the main reasons so many "peer-aware" options
still feel fake to the user.
Contact keeps getting sold as if contact had already become judgment.

That is one of the cleanest expressions of the user's complaint in the whole
routing layer.

The user is tired of surfaces that can contact each other, ping each other,
discover each other, or list each other while still not being able to answer:

- should this exact request go there now?
- and if it does, will it still mean the same thing?

### 5. Fallback-route persistence

Question:

- when the preferred backend disappears, does the route needed for recovery
  still exist?

This is one of the hardest and most important seams in the current runtime.

The repo already knows about a specific live weakness:

- `docker-gen-failover` is present in the edge stack
- the master plan explicitly records that the current approach deletes routes
  when containers stop

This is not a minor caveat.
It is exactly the kind of failure the user is trying to stop:

> the system looked dynamic until the bad event arrived, then the dynamic route
> vanished with the thing it was supposed to route around

That is why the docs keep treating `docker-gen-failover` as both meaningful and
dangerous:

- meaningful because it shows the repo is trying to generate fallback-aware
  Traefik config
- dangerous because the current generation model is recorded as losing routes
  at the wrong moment

That "wrong moment" language should be read literally.
The repo is criticizing solutions that advertise adaptability until the bad
event arrives and then discard the very route that was supposed to preserve
dignity.

This is not just a technical annoyance.
It is one of the betrayal patterns the user keeps reacting to:

- the option sounds dynamic
- the option sounds modern
- the option sounds like failover
- the exact failure reveals the option never fully owned the promise

### 6. Policy continuity

Question:

- if traffic is forwarded, is it still the same protected service?

This layer is easy to under-document and central to the repo's actual demand.

For protected HTTP surfaces, continuity has to include:

- TinyAuth behavior
- `nginx-traefik-extensions` forward-auth behavior
- middleware ordering
- security filtering implications
- headers, rewrites, redirects, and path behavior that define route identity

That means a forwarded route is not proven merely because:

- the upstream answered
- a login page still appears
- some `200` exists

The route has to stay meaningfully the same route.

Otherwise the stack has not preserved the service contract.
It has only preserved enough surface behavior to sound successful in a shallow
status update.

That is why protected-route routing has to stay this severe.
The user is not asking for a success report that can survive on status codes.
They are asking whether the system still deserves to say "this request worked"
without smuggling in a human who already knows what silently changed.

### 7. Stateful truth

Question:

- even if routing works, does the route still point to a service with honest
  authority semantics?

This is where HTTP optimism often becomes a lie.

The root runtime already includes state-bearing services such as:

- MongoDB
- Redis
- Headscale

Those services can already be discussed in routing terms because they have live
route material.
They cannot be broadly discussed in HA terms unless the repo also answers:

- who owns writes?
- what replicates from what?
- how does promotion work?
- what reconnect behavior do clients need?
- what remains tied to one local disk path?

That is why this page refuses to let TCP and stateful surfaces inherit success
from one HTTP drill.

This repo is not anti-HTTP proof.
It just refuses to let one HTTP proof become a laundering mechanism for much
harsher service classes.

## What the repo can already say without lying

The repo can already say all of these honestly:

- the priority implementation has a real and non-trivial routing surface
- the stack already distinguishes public entry, internal traffic, and special
  routing domains
- the edge is already policy-bearing rather than a dumb proxy
- the repo explicitly wants local-first then peer-forward behavior
- the repo explicitly knows that Cloudflare plurality is weaker than preserved
  request meaning
- the repo explicitly knows its current failover helper is not yet trustworthy

That is already meaningful.
It is still far from generic closure.

That gap between meaningful progress and believable option is one of the
hardest things the docs need to preserve.
Too many infrastructure discussions erase that gap the moment enough machinery
appears.

## What the repo still cannot say honestly

The repo still cannot say all of these without stronger route-specific proof:

- "wrong-node routing is basically solved"
- "the stack is anti-SPOF now"
- "multi-record DNS plus Traefik gives us failover"
- "protected routes keep their meaning after peer handoff"
- "TCP and stateful failover are covered by the same routing story"
- "`docker-gen-failover` already gives durable rescue paths"

Those sentences are exactly the kind of overpayment this knowledgebase is
trying to stop.

The user has already paid that overpayment too many times in the surrounding
space:

- a component exists
- a helper exists
- a route exists
- therefore the burden must have moved

This page is here to keep saying that last step is exactly what remains
unproven.

## The real pressure this page must preserve

The user is not frustrated because there are no tools.
They are frustrated because too many respectable tools answer only one routing
layer at a time and then let documentation speak as if the whole burden moved.

That last clause matters.
The real enemy here is not only missing functionality.
It is partial functionality plus overconfident narration.

The real routing pressure in this repo is:

- first hop should not be sacred
- locality should stay real
- wrong-node entry should not collapse into folklore
- rescue routes should not evaporate under the bad day
- protected routes should stay semantically coherent
- stateful services should stay under harsher truth rules

That list is not only a technical decomposition.
It is a record of the ways the available options keep becoming smaller than
they first sound.

If any of those get blurred together under "HA," the docs regress.

## Bottom line

This repo already has a serious routing stack.
It does not yet have a proven request-preserving routing truth stack.

The difference is exactly this:

- reachability is partially real now
- request preservation is still only partially system-owned

Everything in `bolabaden-infra` keeps circling that gap because it is the real
one.

Another way to say the same thing:

- the repo is not starved for components
- it is starved for options that remain honest after wrong-node entry, backend
  loss, and stateful consequences are all allowed into the same sentence

That sentence is close to the center of the whole documentation effort.
If the site fully internalizes it, the rest of the pages can stay honest.
If the site forgets it, the docs keep regressing into polished ambiguity.
