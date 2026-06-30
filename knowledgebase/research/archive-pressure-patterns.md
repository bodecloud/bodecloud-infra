# Archive Pressure Patterns: What the User Keeps Forcing the Repo to Confront

This page exists because the imported plaintext archive is not just extra
background.

It is the clearest record of what the user keeps asking for after ordinary
answers fail.

That matters because this repo can only be documented honestly if it preserves
the pressure that created it:

- the user asks for multi-node behavior that feels like it should already exist
- the ecosystem keeps offering either static glue or heavyweight platform
  migration
- the user rejects those answers because the actual problem survives
- the same missing capabilities keep reappearing under different tool names

Without this page, the knowledgebase drifts back toward generic infrastructure
writing and starts forgetting why this repo is so dissatisfied in the first
place.

That forgetting is one of the main failure modes of the whole site.

The runtime, plans, and instruction surfaces can all be documented accurately
and still miss the point if the archive pressure disappears.
The archive is where the repo most clearly remembers that the real complaint is
not "I need more infrastructure nouns."
It is:

> I keep getting offered answers that sound like options while leaving the same
> hidden sacred-node burden intact

This page exists to keep that complaint from being domesticated into something
more ordinary.

## Evidence boundary for this page

This is an archive-pressure map.

It is authoritative about:

- recurring demand
- recurring refusal
- recurring disappointment
- recurring capability thresholds

It is not authoritative about:

- current live root-runtime behavior
- proof that the priority implementation already satisfies those demands
- which future option has already won

The archive is authoritative about problem shape, not about shipped runtime
truth.

That distinction matters because the archive is unusually rich.
It is tempting to let repeated pressure masquerade as proof of direction or
near-implementation.

That would be the wrong lesson.

The archive tells us what the user keeps refusing, what kinds of answers keep
failing, and what kind of system they are really trying to call into
existence.
It does not tell us that the current runtime has already become that system.

## The archive's core complaint

Across the imported material, the user is not mainly asking for more services.

They are repeatedly asking some version of:

> why is there no straightforward middle layer that makes multiple ordinary
> Docker nodes behave dynamically and resiliently without forcing me into static
> glue or a heavyweight cluster platform?

That question is more specific than "I want HA."
It is really a stacked demand:

- any node can take traffic
- locality should be preserved when possible
- wrong-node requests should still succeed
- manual upstream rewriting should not be the whole strategy
- node survival is not enough if the service path breaks
- middleware and auth should not silently disappear during fallback
- stateful systems should not be lied about
- orchestration should have to earn its complexity

The archive keeps restating the same demand with different vocabulary because
the ecosystem keeps answering adjacent questions instead of this one.

That repetition is not noise.
It is one of the strongest signals in the whole repo.

When a user keeps restating the same dissatisfaction through different tool
names, the right reading is not "they are indecisive."
The right reading is:

the ecosystem keeps forcing them back into the same false choice space.

## What this pressure is doing to the repo

The archive is not just research.
It is acting like a constraint generator.

It keeps forcing the docs and the architecture to answer:

- which truths are still implicit
- which behaviors still depend on human memory
- which resilience claims are overstated
- which platform taxes are still unjustified

That is why the archive matters so much.

The archive is effectively acting like a negative benchmark engine.
It keeps recovering the exact kinds of coherence the user no longer trusts:

- a route that exists but dies when it matters
- a platform that sounds distributed but still depends on one remembered host
- an HA story that is really just ingress plurality
- a control plane that is really just private human reconstruction with nicer
  nouns

## Pattern 1: "There should already be a standard dynamic model, but there never is"

This is one of the loudest repeating frustrations in the archive.

The user repeatedly collides with an ecosystem where:

- Compose handles single-node service definition reasonably well
- multi-node behavior becomes hand-built glue
- "dynamic" often still means "you predeclared a lot of static truth"
- every answer seems to require one more pile of integration ritual

What the user is rejecting here is not YAML itself.

They are rejecting the experience of having to reconstruct runtime truth
through:

- scattered labels
- partial templates
- node-specific edits
- undocumented conventions
- manual peer knowledge

That pressure is what keeps the repo searching for a middle layer between:

- raw Compose sprawl
- and full orchestrator worldview capture

This is where the docs must resist a common flattening move.

It is too easy to rewrite this pattern as:

- wants cleaner automation
- wants better service discovery
- wants less manual config

Those are all true and all too small.

The deeper pattern is:

the user wants the system to stop requiring private reconstruction of topology
and recovery truth just to behave like one platform under stress.

### What this pattern forces the docs to say

Do not describe the repo as if it is merely looking for better Compose hygiene.

It is looking for a more standardized dynamic behavior model than raw
Compose alone gives.

## Pattern 2: "Wrong-node requests should still work"

This is probably the single most important recurring demand in the archive.

The user does not mainly want:

- replicas behind a load balancer
- or DNS that eventually lands on a healthy node

They want something stricter:

- a request can land on a node that does not host the target service
- that node still preserves the request path
- the node forwards or falls back correctly
- the operator does not manually wire every service or node permutation

That is much stronger than ordinary load balancing.
It quietly requires all of these things:

- entry-node independence
- current placement truth
- route persistence
- coherent peer forwarding
- middleware and auth continuity
- convergence of secrets and deployment state

Most ordinary answers in the archive solve one slice of that and then talk as
if they solved the whole thing.

### What this pattern forces the docs to say

Never collapse:

- node-entry survival
- peer forwarding
- route persistence
- request continuity
- full request correctness

into one vague phrase like "HA routing."

That warning should govern more than this page.

The entire site fails if it starts describing:

- entry-node plurality
- peer reachability
- route generation
- request preservation
- semantic continuity

as if they were one maturity tier.

## Pattern 3: "Do not hardcode everything, but also do not immediately jump to Consul, Swarm, or Kubernetes"

This is the central refusal pattern in the archive.

The user repeatedly rejects both extremes.

### Rejected extreme A: static glue everywhere

- hardcoded per-service upstream tables
- hand-maintained peer lists
- solutions that only stay dynamic if the operator edits every change

### Rejected extreme B: heavyweight control-plane surrender

- "just use Kubernetes"
- "just use Swarm"
- "just use Consul"
- "just use Nomad" when the answer never explains how it preserves the user's
  desired operator surface

The archive keeps forcing a narrower target:

- dynamic behavior
- lighter coordination
- more direct ownership of what the system thinks is true

This is one of the clearest places where the user's dream stops sounding like a
normal product comparison.

They are not only asking for lighter machinery.
They are asking for machinery whose truth can still be inspected and whose
benefit can still be named without hand-waving.

That is why repo-native ideas keep reappearing:

- `services.yaml`
- sync-agent
- failover-agent
- file-generated Traefik config
- Compose-first plus stronger truth layers

### What this pattern forces the docs to say

Optionality is not indecision here.

It is the result of the user refusing both:

- static pain
- and ideological platform capture

## Pattern 4: "DNS failover is not the whole story"

The archive repeatedly rejects the shallow version of resilience:

- multiple A records
- some client-side redistribution
- and then calling the problem solved

The user keeps pushing past that because DNS only helps with:

- landing on a surviving node

It does not prove:

- the surviving node knows where the service is
- the request path survives
- auth and middleware survive
- stateful correctness survives

This distinction keeps resurfacing whenever an answer stops too early and
declares victory at the ingress edge.

### What this pattern forces the docs to say

Cloudflare multi-A DNS, Cloudflare LB, keepalived VIPs, and similar entry-layer
mechanisms should be documented as:

- ingress-layer tools
- not proof of end-to-end service continuity

This is exactly the kind of distinction generic self-hosting discourse keeps
blurring.
The archive matters because it keeps restoring the user's impatience with that
blur.

## Pattern 5: "I want equality between nodes if possible"

The archive shows repeated resistance to sacred infrastructure roles:

- sacred ingress box
- sacred manager
- sacred leader that quietly becomes the real control point

This does not mean every service must be symmetric.
It means the repo keeps trying to avoid a system where anti-SPOF language is
built on top of one quietly indispensable machine.

That is why "equality between nodes" is not a cosmetic preference.
It is a resistance to the pattern where the architecture diagram says
distributed while one host or one remembered role still carries the real
intelligence of the platform.

### What this pattern forces the docs to say

Whenever a role becomes operationally sacred, the docs should name it plainly
instead of smuggling it through "cluster" vocabulary.

## Pattern 6: "Stateful honesty matters more than architectural aesthetics"

The archive repeatedly resists the move where stateless routing progress gets
widened into a general resilience story.

The user keeps pushing on questions like:

- who owns writes?
- what replica or quorum model exists?
- what survives node loss?
- how do clients rediscover the surviving topology?

That is why the repo has to keep separating:

- ingress continuity
- request continuity
- stateful correctness

### What this pattern forces the docs to say

Never use successful HTTP routing language to imply Redis, MongoDB, RabbitMQ,
or Postgres now have honest failover semantics.

## Pattern 7: "The docs themselves are part of the problem if they get too congratulatory"

The archive pressure is not only about the runtime.
It is also about narration.

The user keeps rejecting answers that:

- smooth over contradictions
- upgrade intent into proof
- treat modularity as if it were already orchestration
- act as though the lack of options has already been solved

That means the docs are part of the control surface.
If they speak too broadly, they recreate the same problem at the documentation
layer.

### What this pattern forces the docs to say

The knowledgebase should read like a map of exact pain, exact evidence, and
exact missing layers, not a congratulatory essay about modern infrastructure.

## What the archive keeps forcing the architecture to admit

Across these patterns, the archive keeps forcing the repo to admit three things
at once:

1. raw Compose is too static once wrong-node routing and failover become real
2. the ecosystem's standard answer often jumps too quickly to heavyweight
   control planes
3. the missing middle layer is not imaginary; it is the real subject of the
   project

That is why the repo keeps circling:

- current-state registry ideas
- sync-agent and failover-agent ideas
- OpenSVC exploration
- Nomad exploration
- Garden, k3s, and Kubernetes exploration

This is not indecision for its own sake.
It is the shape of a repo trying to find the smallest honest layer that removes
the right burdens.

## Bottom line

The archive keeps saying the same thing in different ways:

> the user does not want prettier YAML, and they do not want premature platform
> surrender. They want a real middle layer that makes multiple Docker nodes
> behave like one request-preserving, operator-readable system without lying
> about state or failover.

If the rest of the knowledgebase stops reflecting that pressure, it stops
reflecting `bolabaden-infra`.
