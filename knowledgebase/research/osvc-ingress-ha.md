# OpenSVC Ingress HA: The Exact Bet This Generator Is Making

This page is about one specific bet in the repo:

> if node membership and live container labels can be discovered at runtime,
> Traefik routes can be generated from reality instead of being hand-written
> for every service and every node.

That bet matters because the user is not mainly asking for "dynamic config."
The user is asking for a stack where several ordinary Docker nodes stop acting
like isolated islands that happen to share a Git repo.

More specifically, the user is asking for a category of option the ecosystem
keeps failing to present cleanly:

- not static per-node folklore
- not "just pick a giant orchestrator"
- but a thinner truth-owning layer that can still answer the wrong-node event
  honestly

This generator matters because it is one of the clearest live attempts to make
that category real.

This page is still a research-and-interpretation page.
It is **not** proof that the current tracked root runtime already provides
universal wrong-node success or that OpenSVC already governs the live stack.

That boundary matters because generated ingress is one of the easiest places
for a repo like this to sound nearly solved while still leaving the deepest
truths unowned.

## What this page is and is not allowed to prove

This page is allowed to:

- explain exactly what the current OpenSVC ingress generator is trying to own
- show where generated routes really come from and what they still assume
- clarify the difference between generated config and owned distributed truth
- recover the precise bet this branch is making about local-first and peer
  fallback ingress

This page is not allowed to:

- claim universal wrong-node success
- imply OpenSVC already governs the live root runtime
- treat runtime generation as proof that DNS, peer eligibility, or semantic
  continuity are solved
- use script detail as a substitute for passed failure drills

## Quick claim router

If the question is:

- "What exactly is this generator doing?" this page is a primary answer.
- "Does generated ingress already equal proven failover?" no.
- "Why is this branch strategically interesting?" this page is meant to answer
  that precisely.
- "Can I cite this as route-level proof?" not by itself.

## What the current script literally does

Primary artifact:

- [`scripts/osvc_ingress_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_ingress_sync.py)

The current generator:

- requires `DOMAIN`
- derives the local node from `TS_HOSTNAME`, or falls back to `hostname -s`
- reads cluster node names from `om node ls --format json`
- reads running containers from `docker ps -q` plus `docker inspect`
- keeps only containers where:
  - `traefik.enable=true`
  - at least one `traefik.http.*` label exists
- chooses the backend port by:
  - preferring `traefik.http.services.*.loadbalancer.server.port`
  - otherwise falling back to the first exposed container port
- writes
  [`${CONFIG_PATH}/traefik/dynamic/failover-fallbacks.yaml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/volumes/traefik/dynamic/failover-fallbacks.yaml)
- creates two file-provider services per discovered container:
  - `<name>-direct`
  - `<name>-with-failover`
- creates two routers per discovered container:
  - `<service>.<domain>`
  - `<service>.<localnode>.<domain>`
- uses optional health hints from `kuma.healthcheck.path`,
  `kuma.healthcheck.interval`, and `kuma.healthcheck.timeout`

That is already much more specific than "there is some OpenSVC ingress idea."

It also reveals the first hard boundary:

- the generated local backend is `http://<container-name>:<port>`
- remote fallbacks are `https://<service>.<peer-node>.<domain>`

So this model is not doing deep service discovery across peer nodes directly.
It is doing:

1. local-container routing for the fast path
2. peer-hostname routing for the remote path

That distinction matters because it tells you exactly where the design is
strong and exactly where it can still lie.

It is one of the most useful exactness tests in the whole knowledgebase:
what does the generator really know, and what is it still assuming indirectly
through naming, DNS, reachability, or shared policy?

That is also why generated config is dangerous in this repo.
It can look like inspectable runtime intelligence while still smuggling
critical truths in through adjacent systems the generator does not itself own.

## Strongest honest current answer

The strongest honest current answer is that this generator is one of the
clearest live attempts in the repo to move fallback routing out of handwritten
folklore and into runtime-derived config. That is real progress. It is still
not the same thing as system-owned wrong-node truth, because the generated
result still depends on adjacent realities like DNS shape, east-west
reachability, and preserved middleware meaning that the generator does not
fully own.

## The dream this generator is trying to encode

The repo's desired request model is not ordinary load balancing.
It is:

1. a request can land on any surviving public node
2. if that node already hosts the service, it should serve locally
3. if not, the node should still have a meaningful next step
4. the operator should not have to hand-author every fallback permutation

The generator is one of the clearest attempts in the repo to encode that dream
without immediately replacing Compose with a heavyweight scheduler.

That is why this page matters beyond OpenSVC itself.
It is a live artifact of the repo trying to turn "wrong node" from a private
operator panic into a runtime-readable next step.

## The two hostname classes it creates

### Node-scoped hostnames

Pattern:

- `<service>.<node>.<domain>`

Generated router rule:

- `Host(\`<service>.<localnode>.<domain>\`)`

Meaning:

- this node is supposed to be the first authority for that hostname
- if the service is local, the request should complete locally
- if the service is not local, the failover service still has to preserve the
  request meaning through a peer hostname

This is the operator-facing precision path.
It matters for debugging and for proving whether locality-aware routing is real
or only described.

### Global hostnames

Pattern:

- `<service>.<domain>`

Generated router rule:

- `Host(\`<service>.<domain>\`)`

Meaning:

- any ingress-capable node may receive the request
- that node tries local service first
- then falls through to peer-specific hostnames for other nodes

This is the public any-node dream in its clearest current script form.

## What this design gets right

The strongest thing here is not OpenSVC branding.
It is the refusal to hardcode a giant per-service, per-node routing matrix.

That refusal matters because static route files become one of the easiest ways
for a multi-node Docker stack to lie:

- a route still exists in config
- the backend reality moved
- the operator forgets the mismatch
- the docs keep narrating continuity that no longer exists

This generator is trying to make routing follow live evidence instead of stale
operator memory.

That is the real anti-SPOF move here, more than any brand name.
It is trying to move request truth out of recollection and into something the
runtime can derive.

But it should only be praised to the exact degree that the derivation is real.
If DNS, peer eligibility, or middleware continuity are still socially
reconstructed outside the generator, then the docs should say that plainly
instead of letting "generated" impersonate "owned."

## What this design quietly depends on

The generator is only as honest as its weakest truth surface.

### Truth surface 1: node membership

`om node ls --format json` has to be accurate enough to describe the nodes that
should be considered eligible peer targets.

### Truth surface 2: local runtime discovery

`docker inspect` plus Traefik labels have to accurately reveal:

- which services are running
- which services are HTTP-routed
- which backend port Traefik should actually use

### Truth surface 3: DNS shape

The peer fallback URLs only mean anything if the naming layer exists:

- `<node>.<domain>`
- `*.<node>.<domain>`
- some viable global-entry strategy for `<service>.<domain>`

Without that, the generated fallback URLs are elegant fiction.

### Truth surface 4: east-west reachability

The receiving node has to be able to reach peer nodes reliably enough that
`https://<service>.<peer>.<domain>` is a real recovery path instead of a second
guess.

### Truth surface 5: semantic continuity

The request has to keep meaning the same thing after the peer hop:

- auth
- middleware ordering
- headers
- challenge flow
- app-visible request assumptions

This is one of the main places where a system can "work" and still betray the
user's real requirement.

That is why the user keeps rejecting explanations that stop at successful
transport.
The hard part is not whether the request can move.
It is whether it stays the same request in a way the user would recognize as
honest.

That phrase should stay central.
The user's frustration is not solved by transport alone.
It is solved only when the system stops quietly changing the meaning or trust
conditions of the request just because it landed on the wrong node first.

## What this script does not currently prove

It does **not** prove that:

- the local backend disappears and the global route still survives correctly
- the peer-specific hostname definitely lands on the right node during failure
- auth and middleware continuity remain stable under the peer hop
- peer hostnames are consuming a trustworthy current-state placement registry
- OpenSVC has become the final authoritative placement truth for the stack

That is why this page has to stay harsher than the original `docs/osvc_ingress_ha.md`.

The repo has too many surfaces that can now produce a convincing fallback demo.
This page exists to keep convincing demo and preserved meaning from being
treated like synonyms.

## Where the route-generation model is strongest

It is strongest for:

- stateless or mostly stateless HTTP services
- services already described cleanly by Traefik labels
- environments where per-node wildcard DNS exists
- situations where removing manual route sprawl is already a big win even
  before universal failover proof exists

In those cases, the generator can materially reduce operator reconstruction
tax.

That reduction is important.
Even partial truth automation matters in this repo, as long as the docs do not
upgrade partial reduction into universal closure.

## Where it is still weak

It is weak when:

- a service's local and remote semantics are not interchangeable
- middleware behavior is order-sensitive
- peer hostnames can exist before placement truth is actually trustworthy
- the local router disappears at exactly the moment the fallback should matter
- app behavior changes subtly when the request is served remotely

This is why generated config is not the same thing as preserved request
meaning.

It is one of the strongest recurring rules in the whole rewrite.

## The DNS question this page must keep honest

The generator only solves the "what should Traefik try next?" part.
It does not solve "which public node did the client reach first?"

The repo still has several possible first-hop models:

- Cloudflare Load Balancer
- self-hosted VIP / keepalived-style entry
- round-robin DNS

Those are not equivalent.
They change:

- health awareness
- cache behavior
- how fast dead nodes stop receiving traffic
- how much fake HA language the docs are allowed to use

That last point is exactly why first-hop strategy can never be treated as
decorative in this repo.

So no document should say "OpenSVC ingress HA is solved" without also naming
the actual first-hop strategy.

## Why the TCP boundary matters so much

The current generator is explicitly HTTP-only.
That is a healthy constraint.

The repo also has:

- [`scripts/osvc_l4_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_l4_sync.py)
- [`compose/docker-compose.l4-ingress.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.l4-ingress.yml)

But that L4 path is a different class of claim entirely.
For TCP services, reachable frontend ports still say nothing by themselves
about:

- primary election
- replica correctness
- write continuity
- state survival

The repo is right to keep those categories separate.

That separation is not pedantry.
It is one of the main things preventing the knowledgebase from sounding more
complete than the runtime has earned.

It is also one of the main things preserving the user's actual dream instead of
flattening it into "the cluster has a proxy now, therefore failover is mostly
handled."

## Bottom line

The OpenSVC ingress path matters because it is one of the clearest attempts in
the repo to turn multi-node ingress from static folklore into generated,
runtime-shaped behavior. Its real value is not "dynamic YAML." Its value is
that it tries to preserve the user's actual dream: any node can receive the
request, the local node should win when it can, and peer fallback should not
require hand-maintained per-service sprawl. What remains unproven is whether
that generated model stays semantically trustworthy when the local path breaks
and the peer path becomes the only path that matters.
