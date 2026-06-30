# DevOps Runbook

This is not a generic Docker runbook.

It is the operator discipline for one unusually specific and frustrating
problem:

> can this repo stop depending on hidden operator memory, lucky node locality,
> and vague HA language, or does it still only *look* more resilient than it
> actually is?

That is the real operational question behind `bolabaden-infra`.

The user is not mainly asking for startup commands.
They are asking for an operating method that can tell the difference between:

- a Compose-first multi-node stack that genuinely preserves requests under
  pressure
- and a pile of partly-true components that feels impressive right up until
  the request lands on the wrong box or the preferred backend disappears

This runbook exists to make that difference visible.

## What this page is and is not allowed to prove

This page is authoritative about:

- how operators should route questions to the right proof class
- what kinds of command output are too weak for stronger claims
- how to stop green output from inflating into resilience theater

This page is not authoritative about:

- whether a given route already survives the wrong node
- whether a given stateful class is already resilient
- whether the overall dream is already satisfied

This page governs operational discipline, not completion status.

## Quick claim router

| If the operator is really trying to prove... | Start with... | It still must not imply... |
| --- | --- | --- |
| authored runtime coherence | `docker compose config --quiet` and related authoring checks | that services start, routes work, or failover is real |
| local service presence and health | `docker compose config --services`, `docker compose ps`, `docker inspect` health | that local health implies cross-node truth |
| route correctness on a normal day | path request plus edge and backend identity evidence | that wrong-node or backend-loss behavior is solved |
| wrong-node or backend-loss behavior | intentional route drills and failure drills | that one passing drill proves whole-stack resilience |
| stateful correctness | topology-specific ownership, election, reconnect, and storage evidence | that generic cluster or ingress success means state is safe |

If a command class is weaker than the claim being narrated, the operator should
downgrade the claim, not upgrade the command.

## What an honest runbook means in this repo

An honest runbook here is not just a list of commands.

It is a way of refusing to let healthy-looking output outrank the real
question.

That means the runbook has to keep forcing the operator to say:

- what exact claim is under test
- which proof class the evidence belongs to
- which stronger story would still be a lie even if the command passed
- whether the explanation for success now lives more in inspected system truth
  or still in private operator recollection

## What this runbook is trying to kill

The hardest operational burden in this repo is not container startup.
It is hidden reconstruction burden.

That burden appears whenever the operator has to privately reconstruct truth
the platform should already be exposing.

In practice that means four recurring hidden burdens.

### 1. Hidden topology burden

The operator should not have to privately remember:

- which node currently hosts the service
- whether the route is local, remote, or stale
- whether the candidate peer is real fallback or merely reachable

If the system still depends on that memory, the operator's head is still the
control plane.

### 2. Hidden convergence burden

The operator should not have to guess:

- whether nodes are on the same revision
- whether env and secret surfaces are aligned
- whether auth and middleware assumptions are aligned
- whether a forwarded request would land in semantically equivalent runtime

Transport success without convergence truth is not a real recovery story.

### 3. Hidden claim burden

The operator should not have to keep translating:

- "multiple DNS records" into "but not necessarily preserved requests"
- "healthy proxy" into "but not necessarily valid failover"
- "container restarted" into "but not necessarily state-safe recovery"

If the operator is still doing that translation manually, the docs and tooling
are not honest enough yet.

### 4. Hidden proof burden

The operator should not have to reconstruct:

- which exact evidence supports the claim
- which exact failure class was exercised
- what remained unproven after the test passed

This runbook exists to stop those hidden burdens from being normalized.

## Start by naming the real question

Before running commands, state what you are actually trying to prove.

In this repo, the serious questions are usually one of these:

- does the priority merged root Compose graph still resolve?
- what services are actually in the current root runtime?
- can this node receive and interpret traffic at all?
- does a named ingress path still behave as documented?
- does one wrong-node request still succeed?
- does fallback survive local backend loss?
- does the relevant state topology remain correct after failure?

If you do not name the question first, green output will try to trick you into
claiming more than it supports.

That is one of the main failure modes the user is explicitly trying to remove.

This step is more than ceremony.

If the question is not named first, the repo tends to fall into the same trap
as the wider ecosystem:

- run enough commands to feel comforted
- then narrate that comfort as resilience

That is exactly the operational version of the documentation failure the user
keeps rejecting.
The point of this runbook is to make "what did this actually prove?" a normal
question again.

## Route the question to the right proof class first

This runbook gets much easier once the operator stops asking one giant
"is it working?" question.

Use this routing table before running anything.

| If the real question is... | Start with... | Next required proof if the first step passes | Do not upgrade into... |
| --- | --- | --- | --- |
| Does the priority merged root runtime still resolve? | `docker compose config --quiet` | `docker compose config --services` to inspect what actually merged | any claim about runtime behavior, routing, failover, or resilience |
| What is actually in the current root runtime? | `docker compose config --services` | `docker compose ps` and local inspection for the specific service | any claim that the service works cross-node |
| Did a local service start and stay healthy? | `docker compose ps` plus `docker inspect <service> --format='{{json .State.Health}}'` | route-targeted path proof if the service is publicly consumed | any claim that local health implies peer safety or distributed readiness |
| Does one ingress path answer normally? | route-targeted request plus receiving-node logs | backend identity proof and, when relevant, auth or middleware proof | wrong-node success, backend-loss survival, or stateful correctness |
| Does one wrong-node HTTP request still preserve the service? | intentional wrong-node request plus receiving-node and backend-node logs | backend-loss drill for that same route | generic multi-node success for all routes |
| Does fallback survive when the preferred backend disappears? | known-good route plus backend-stop or backend-loss drill | semantic continuity comparison between local and peer-forwarded behavior | broad failover claims for the stack |
| Is one stateful service actually resilient? | explicit topology inspection: owner, replica or peer, client reconnect story | real failure drill for that stateful class | "the platform is HA now" |

The point of this table is not convenience.
It is to stop the operator from burning time on the wrong command class and
then trying to inflate its meaning afterward.

## Strongest honest current answer

If an operator asks, "What is the real job of this runbook?" the shortest
defensible answer is:

> Its job is to stop operational proof from collapsing into comfort theater by
> forcing each claim to declare its proof class, its ceiling, and what still
> remains dependent on private operator reconstruction even after the command
> succeeds.

Anything softer than that will let the same fake-HA story creep back in through
green output instead of prose.

## The proof classes this runbook assumes

This runbook is downstream of
[`proof-matrix-and-drills.md`](proof-matrix-and-drills.md).

Every useful operational statement here belongs to one proof class.
Do not merge them just because the screen output looked healthy.

### 1. Authoring proof

Questions answered:

- does the tracked graph interpolate?
- do the docs render?
- is the authored configuration structurally coherent?

Typical commands:

```bash
docker compose config --quiet
python3 -m mkdocs build -f mkdocs.yml --strict
```

What this proves:

- the repo's authored surfaces still resolve

What this does not prove:

- services start
- routes work
- peers agree
- failover works
- state remains correct

This proof class is necessary.
It is almost never sufficient for what the user actually cares about.

### 2. Local runtime proof

Questions answered:

- what is locally present?
- what started?
- what healthchecks are passing?

Typical commands:

```bash
docker compose config --services
docker compose ps
docker inspect <service> --format='{{.State.Health.Status}}'
```

What this proves:

- the node has a certain local running shape

What this does not prove:

- distributed correctness
- wrong-node behavior
- peer eligibility
- fallback persistence after loss
- state safety

The user is specifically frustrated with stacks that stop at this proof class
and then start using HA vocabulary anyway.

That warning should stay front and center.

Local runtime proof is often the exact point where a system becomes seductive
enough to over-describe and not yet real enough to deserve it.

### 3. Path proof

Questions answered:

- does ingress reach the intended route?
- do auth and middleware still behave correctly?
- what backend actually answered?

This proof class requires more than a success code.
At minimum it should include:

- a route-targeted request
- the node that received the first hop
- relevant edge logs
- backend confirmation when backend identity matters

If the path's meaning matters, a `200` alone is weak evidence.

If the explanation still requires "I know that node B is the real one," the
proof is also weaker than it looks.
The runbook is not only checking whether traffic moved.
It is checking whether the system itself can increasingly explain why the move
was honest.

### 4. Failure proof

Questions answered:

- what happens when the request lands on the wrong node?
- what happens when the preferred local backend disappears?
- what happens when a candidate peer becomes unavailable or ineligible?
- does the route survive the failure long enough to matter?

This is the proof class where most infrastructure theater dies.

If a system has never been exercised here, its resilience story is still
mostly aspirational.

For this repo, that should be read literally.

A multi-node stack that has not been exercised under wrong-node or backend-loss
conditions is still mostly describing a desired mood, not a trustworthy
recovery surface.

It is also the proof class where human folklore often sneaks back in.
A failure drill that only works once the operator remembers the unofficial
topology is not useless, but it has not yet removed the hidden control plane.

### 5. State proof

Questions answered:

- who owns the state now?
- did election, promotion, or ownership transfer happen correctly?
- did clients rediscover the correct survivor?
- was durable correctness preserved rather than mere liveness?

Without this proof class, stateful HA language is still illegal.

## Validation prerequisites

The priority root Compose surface is not free to inspect.
Validation depends on a prepared environment.

At minimum expect:

- a populated env surface based on
  [`.env.example`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.env.example)
- placeholder or real secret files under `${SECRETS_PATH}`
- `~/.docker/config.json` present, even if it only contains `{}`
- meaningful `DOMAIN` and `TS_HOSTNAME` values if hostname interpolation is
  expected to mean anything
- required secrets such as `SEARXNG_SECRET` before the merged graph can fully
  resolve

This matters because one of the easiest ways to lie is to narrate the priority
runtime without actually being able to inspect it.

It also matters because partial inspection encourages over-joining.
When merged runtime truth is gated by env and secrets, operators are more
tempted to fill gaps with memory.
This runbook should keep treating that temptation as a risk surface, not as
harmless background noise.

## The first honesty wall

If the root Compose graph does not resolve, stop narrating runtime certainty.

That is the first wall.

The root implementation is still the priority runtime surface.
If it does not resolve, the next claim should be about missing env or secret
surfaces, not about resilient behavior you did not actually inspect.

Honest statements:

- "The root Compose graph does not currently resolve because `DOMAIN` is unset."
- "The merged stack cannot be inspected yet because `SEARXNG_SECRET` is
  missing."
- "We do not currently have a real service inventory from the priority merged
  stack."

Dishonest upgrades:

- "The stack probably still contains X."
- "The route should still be fine."
- "Failover is probably unaffected."

This is a hard line because the user is asking for less narrative cheating, not
more.

The runbook should keep normalizing hard stops like this.

Stopping and saying “we cannot truthfully inspect the priority runtime yet” is
more aligned with the repo’s goal than filling the gap with plausible
architecture prose.

## Minimum baseline commands

These commands keep the runbook anchored to the repo's actual current surface.

### Validate the merged root Compose graph

```bash
docker compose config --quiet
```

Use this when the question is:

- does the root implementation still parse?
- did includes, interpolation, configs, and secret references survive the
  change?

What this proves:

- the authored root graph currently resolves

What this does not prove:

- services start
- paths work
- peer forwarding works
- failover works
- nodes agree semantically

### Inspect the priority merged service inventory

```bash
docker compose config --services
```

Use this when the question is:

- what is actually in the priority merged runtime?
- did the real runtime surface change?

Strict rule:

if this command fails, do not claim certainty about the merged runtime you did
not actually inspect.

### Build the documentation strictly

```bash
python3 -m mkdocs build -f mkdocs.yml --strict
```

Fallback if local tooling is unavailable:

```bash
docker run --rm -v "$PWD:/docs" -w /docs squidfunk/mkdocs-material:latest build -f mkdocs.yml --strict
```

What this proves:

- the docs still build
- navigation and markdown structure remain coherent

What this does not prove:

- the docs are complete
- the docs are honest
- the runtime supports the claims

Strict rendering is structural proof, not truth proof.

### Read local healthchecks, but do not let them inflate

Typical commands:

```bash
docker compose ps
docker inspect <service> --format='{{json .State.Health}}'
```

Useful questions:

- did the service start?
- is the local healthcheck passing?
- what check is actually being run?

What this proves:

- local process health under local conditions

What this does not prove:

- remote peer eligibility
- route persistence
- request preservation under wrong-node conditions
- state safety

This repo is right to insist on real healthchecks.
It is also right to refuse to treat them as distributed truth by themselves.

## Baseline route inspection questions

Before saying a route "works," answer all of these:

- what hostname or entrypoint was hit?
- which node received the first hop?
- was the target local or remote from that node's perspective?
- which middleware chain applied?
- which backend actually answered?
- what logs prove that?

If you cannot answer those questions, you may have path proof for reachability,
but not for request preservation.

## Evidence bundle required before writing a stronger claim

For any non-trivial runtime statement, capture an evidence bundle instead of
one green command.

Minimum bundle for a serious route claim:

- the exact claim sentence you are trying to justify
- the proof class you think the evidence belongs to
- the hostname, route, or service identity under test
- the first-hop node identity
- the answering backend identity when backend identity matters
- the command or request used
- the logs or inspect output that prove what happened
- one sentence naming what still remains unproven even if this bundle is valid

This repo needs that last bullet because the main documentation failure is not
usually "no evidence at all."
It is "some evidence exists, so a stronger claim quietly got upgraded."

If you cannot produce the last bullet, the claim is probably already trying to
say too much.

## Wrong-node HTTP drill

This is the most important drill shape in the whole repo.

Purpose:

- prove that a request can land on a node that does not host the target
  service locally
- still preserve the intended route contract

Minimum evidence:

1. identify the receiving node
2. identify the node that actually hosts the target backend
3. send the request through the wrong-node path intentionally
4. capture edge logs from the receiving node
5. capture backend logs or equivalent identity proof from the answering node
6. confirm auth and middleware continuity if that route depends on them

What this drill proves when it passes:

- one specific wrong-node HTTP path is real

What it does not automatically prove:

- all routes behave the same way
- backend-loss recovery works
- TCP classes are solved
- stateful services are covered

Additional honesty checks:

- if the receiving node only succeeded because the operator privately chose the
  hidden real host, the drill is weaker than it looks
- if backend identity cannot be shown, the drill is closer to route
  reachability than to request preservation
- if auth or middleware semantics matter and were not checked, the drill is
  not yet semantic continuity proof

That distinction has to remain explicit because this repo is trying to escape
exactly the habit of widening claims too cheaply.

## Backend-loss failover drill

Purpose:

- determine whether the route required for recovery survives the failure that
  made recovery necessary

Minimum shape:

1. start from a known-good route
2. verify whether the route is local-first under healthy conditions
3. stop, kill, or otherwise remove the preferred local backend
4. inspect whether the route remains present
5. attempt the request again
6. prove whether traffic reached a healthy peer or simply lost the route

This is the drill that exposes the difference between:

- dynamic-looking routing
- and actually resilient routing

If the route disappears with the local backend, the failover claim stays weak
no matter how good the happy path looked.

## Convergence drift check

Purpose:

- answer whether forwarded traffic would still land in semantically equivalent
  runtime

Questions to inspect:

- are node revisions aligned?
- are required env surfaces aligned?
- are relevant secrets aligned?
- are auth and middleware dependencies aligned?
- is there visible evidence of drift?

If these questions are unanswered, cross-node forwarding should still be spoken
of cautiously.

That caution is not overkill.
It is the difference between transport-level optimism and operator-grade truth.

## Stop conditions that should interrupt narration

The runbook should normalize stopping early when the proof ceiling is reached.

Stop and downgrade the claim when:

- the merged root Compose graph does not resolve
- the service inventory cannot be rendered from the priority runtime
- the receiving node cannot be identified
- the answering backend cannot be identified for a backend-sensitive claim
- the route answered, but the proof still depends on remembered unofficial
  topology
- the failure drill changed the request path, but no one verified whether the
  same policy stack still applied
- a stateful service is reachable, but ownership or client rediscovery is
  still unclear

These are not annoyances.
They are the exact places where the repo drifts back toward polished ambiguity.

## Default wording rules after a drill

After any drill, write the result in this order:

1. what exact route or service class was exercised
2. what exact failure or non-failure condition was exercised
3. what proof class was actually reached
4. what stronger sentence is still forbidden

Example shape:

> One stateless HTTP route was intentionally sent to a wrong receiving node and
> still completed through a healthy peer. That is wrong-node proof for that
> exact route. It is not yet backend-loss proof, semantic continuity proof for
> all protected routes, or any kind of stateful HA proof.

That style may sound repetitive.
It is supposed to.
Repetition is one of the few reliable defenses against the exact inflation
pattern the user keeps pushing back on.

## Stateful service drill boundary

For stateful systems, stop reusing the HTTP mental model.

Before calling anything resilient, identify:

- who owns writes
- what replication state exists
- what promotion means
- how clients rediscover the correct owner
- what evidence proves durable correctness after failure

If the answers are weak, the honest output is not:

- "stateful HA is close enough"

It is:

- "the service may be reachable, but stateful correctness is not yet proven"

That sentence is much more useful to this project than false reassurance.

## What to write down after every serious drill

Every real drill should end with a short record containing:

- what question was being tested
- what proof class was obtained
- what exact failure class was exercised
- what succeeded
- what remained unknown
- what larger claim is still not legal

This matters because the repo is trying to escape memory-based infrastructure.
If drill conclusions live only in the operator's head, the system has not
actually matured.

## The operator stance this repo needs

The right operator stance for `bolabaden-infra` is not cynical and not
credulous.

It is:

- respect authoring proof
- demand runtime proof
- separate local health from distributed correctness
- separate transport success from semantic continuity
- separate HTTP success from stateful safety
- treat unclear evidence as incomplete, not as success

That stance matches what the user is actually trying to build:

not just a working stack,
but a stack where the available options are real, the proof boundaries are
visible, and wrong-node or backend-loss pressure stops being the moment when
all the supposed flexibility suddenly evaporates.
