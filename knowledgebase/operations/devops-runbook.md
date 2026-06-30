# DevOps Runbook

This is not a generic Docker runbook.

It is the operator discipline for one unusually specific problem:

> how do we tell the difference between a Compose-first multi-node stack that
> actually preserves requests under pressure and one that only *looks* more
> resilient than it is?

That is the operational question behind `bolabaden-infra`.

The user is not mainly asking for startup commands.
The user is asking for a way to stop hidden operator memory from acting like
the real control plane.

This runbook exists to make that burden visible and to stop green output from
inflating into resilience theater.

## The dream this runbook has to protect

The repo is trying to keep the directness of Compose while pushing toward:

- any-node public entry
- local-first serving
- peer-forward fallback
- anti-SPOF pressure
- honest boundaries where the system still depends on shared operator memory

Operationally, that means the runbook cannot be satisfied by:

- "`docker compose config` passed"
- "Traefik is healthy"
- "the service answered on one node"
- "the container restarted"

Those may be necessary signals.
They are not the same as the dream the user is asking for.

## Strongest honest current answer

The real job of this runbook is:

1. force every claim to name its proof class first
2. force every successful command to state what stronger story would still be a
   lie
3. force the operator to distinguish authored shape, local runtime health,
   route behavior, wrong-node behavior, backend-loss behavior, and stateful
   correctness
4. expose where private operator reconstruction is still doing work the system
   should eventually own itself

If the runbook does not do those things, it becomes a comfort ritual instead
of an operational tool.

## What this page is and is not allowed to prove

This page is authoritative about:

- how to route questions to the right proof class
- what command classes are too weak for stronger claims
- what order to inspect the stack in
- what counts as weak, medium, and strong evidence in this repo

This page is not authoritative about:

- whether a specific ingress path is already resilient
- whether a specific stateful service is already safe
- whether the architecture dream has already been achieved

This is an operator method page, not a completion certificate.

## What the runbook is trying to kill

The real recurring failure mode in this repo is hidden reconstruction burden.

That burden appears whenever the operator has to privately remember or infer
what the system should be exposing explicitly.

There are four recurring versions of that burden.

### 1. Hidden topology burden

The operator should not have to privately remember:

- which node currently hosts the service
- whether the request path is local or remote
- whether the peer being targeted is merely reachable or actually the right
  backend

If the answer still lives mostly in someone's head, the operator's head is
still part of the control plane.

### 2. Hidden convergence burden

The operator should not have to guess:

- whether nodes are on the same revision
- whether secrets and env surfaces still match
- whether auth and middleware assumptions are equivalent across nodes
- whether a peer-forwarded request lands in semantically comparable runtime

Transport success without convergence truth is not a real recovery story.

### 3. Hidden claim burden

The operator should not have to keep mentally translating:

- "multiple DNS records" into "not yet preserved requests"
- "healthy proxy" into "not yet wrong-node proof"
- "container restarted" into "not yet state-safe recovery"

If the operator must keep doing that privately, the docs and tooling are still
too flattering.

### 4. Hidden proof burden

The operator should not have to reconstruct after the fact:

- which evidence class a claim belonged to
- which failure mode was actually exercised
- what remained unproven even after the command passed

That is one of the central reasons this runbook exists.

## Start by naming the real question

Before running anything, say what you are actually trying to prove.

In this repo, the serious questions are usually one of these:

- does the merged root Compose graph still resolve?
- what services are actually present in the priority runtime?
- is a named service locally healthy on this node?
- does one documented ingress path answer on a normal day?
- what happens when the request lands on the wrong node?
- what survives when the preferred local backend disappears?
- did a stateful service keep correctness, not just reachability?

If the question is not named first, the operator will almost always drift into:

- run enough commands to feel reassured
- then narrate that reassurance as resilience

That is exactly the failure pattern the user is trying to get away from.

## Proof classes

This repo becomes much easier to operate once "is it working?" is split into
smaller claim types.

### 1. Authoring proof

Questions answered:

- does the tracked graph interpolate?
- do the docs still render?
- is the authored configuration structurally coherent?

Commands:

```bash
docker compose config --quiet
docker compose config --services
python3 -m mkdocs build -f mkdocs.yml --strict
```

What this proves:

- the authored surfaces still resolve
- the merged root graph still has a coherent shape

What this does **not** prove:

- containers started
- routes work
- peers agree
- failover exists
- state is safe

This is the minimum proof class, not the final one.

### 2. Local runtime proof

Questions answered:

- what is present on this node?
- what started?
- what healthchecks are passing?

Commands:

```bash
docker compose ps
docker inspect <container> --format='{{json .State.Health}}'
docker logs --tail=100 <container>
```

What this proves:

- this node currently has a particular running shape
- the local healthcheck surface is reporting something specific

What this does **not** prove:

- wrong-node behavior
- peer-forward eligibility
- fallback persistence
- stateful correctness

This is one of the biggest temptation zones in the whole repo.

Local runtime proof is exactly where a stack can begin to feel complete while
still leaving the real distributed question unanswered.

### 3. Route correctness proof

Questions answered:

- does a known ingress path answer on a normal day?
- which backend actually answered?
- were the expected auth and middleware surfaces involved?

Typical evidence:

- request output
- receiving-node logs
- backend identity evidence
- auth or middleware logs when relevant

What this proves:

- one documented path still works in a specific scenario

What this does **not** prove:

- wrong-node success
- backend-loss survival
- full service-class resilience

### 4. Failure-path proof

Questions answered:

- does the request survive when it lands on the wrong node?
- does fallback survive when the preferred backend disappears?
- does the post-fallback behavior remain semantically equivalent?

Typical evidence:

- intentional wrong-node request
- receiving-node logs
- peer backend logs
- before/after route identity comparison
- backend stop or simulated backend-loss drill

What this proves:

- one specific failure path was exercised and behaved in a particular way

What this does **not** prove:

- generic stack-wide resilience
- all routes share the same recovery properties

### 5. Stateful correctness proof

Questions answered:

- who owns write truth now?
- what replicates it?
- how do clients discover current authority?
- what happens after promotion or node loss?

Typical evidence:

- service-specific topology inspection
- replication status
- promotion evidence
- reconnect behavior
- storage or authority ownership evidence

What this proves:

- correctness of one stateful class in one tested topology

What this does **not** prove:

- generic "the platform is HA now"

## Claim router

Use this before reaching for commands.

| Real question | Start with | Next proof if that passes | Must not be upgraded into |
| --- | --- | --- | --- |
| Does the priority merged root runtime still resolve? | `docker compose config --quiet` | `docker compose config --services` | claims about runtime behavior or resilience |
| What is actually part of the current root runtime? | `docker compose config --services` | `docker compose ps` | claims about cross-node correctness |
| Did a local service start and stay healthy? | `docker compose ps` and `docker inspect` | route-targeted proof for public services | claims about peer-forward or failover readiness |
| Does one ingress path answer normally? | targeted request plus logs | backend identity and policy-path proof | claims about wrong-node or backend-loss survival |
| Does one wrong-node request still preserve the service? | intentional wrong-node request plus logs on both nodes | backend-loss drill on the same service | generic multi-node success |
| Does fallback survive local backend loss? | known-good path plus backend-loss drill | semantic comparison before vs after | stack-wide failover completion |
| Is one stateful class actually resilient? | topology-specific ownership and replication evidence | real failure drill for that service class | broad HA claims for all state |

## Practical validation baseline

These are the minimum repo-native checks that should stay routine:

```bash
docker compose config --quiet
docker compose config --services
python3 -m mkdocs build -f mkdocs.yml --strict
```

Interpret them narrowly:

- they prove authored shape
- they do not prove runtime behavior
- they definitely do not prove distributed correctness

Also remember the repo-specific constraints:

- Compose validation needs the expected env and secrets surface
- `~/.docker/config.json` may need to exist even if it only contains `{}`
- passing config validation still says nothing about wrong-node handling

## Operator sequence for honest inspection

When the question is not yet fully clear, use this order.

### Step 1: confirm the authored root graph

Use:

```bash
docker compose config --quiet
docker compose config --services
```

Purpose:

- confirm that the priority implementation still merges coherently
- see what the root runtime currently claims to own

Do not claim from this step:

- runtime liveness
- route success
- failover
- stateful safety

### Step 2: identify the exact service class

Before going further, decide whether the target is:

- stateless HTTP
- raw TCP
- state-bearing service
- control-plane service

This matters because the evidence thresholds differ drastically.

HTTP route success is not enough for TCP.
TCP reachability is not enough for state.
Control-plane availability is not the same thing as replicated authority.

### Step 3: inspect local runtime truth

Use:

```bash
docker compose ps
docker inspect <container> --format='{{json .State.Health}}'
docker logs --tail=100 <container>
```

Purpose:

- confirm what is actually local
- learn what the local health and startup surfaces say

Interpretation rule:

- if the service is unhealthy, stop pretending the distributed question is the
  first problem
- if the service is healthy, do not let that health get promoted into
  distributed proof

### Step 4: prove the normal-day route

For public services, gather:

- request result
- receiving-node evidence
- backend identity
- auth or middleware evidence when relevant

Purpose:

- prove what happened on the normal path

Do not claim from this step:

- wrong-node success
- backend-loss survival

### Step 5: exercise the failure path intentionally

If the real question is failover, the operator must stop using flattering
traffic.

Exercise at least one of:

- request lands on a node that lacks the local service
- preferred backend disappears after a known-good path is established

Evidence required:

- receiving-node logs
- backend-node logs
- route identity before and after
- visible continuity or discontinuity of auth and middleware behavior

If those are missing, the result is at best partial.

### Step 6: treat stateful systems separately

For stateful systems, require answers to:

- where does write authority live now?
- what replicates it?
- who decides promotion?
- how do clients discover the new authority?
- what storage survives node loss?

If the operator only has endpoint reachability, they do not yet have a
stateful answer.

## The most important anti-theater rule

Every successful command should be followed by one sentence:

> what stronger claim would still be false even though this passed?

Examples:

- `docker compose config --quiet` passed:
  the stack may still fail on wrong-node entry
- Traefik is healthy:
  peer-forward routing may still be wrong or stale
- service endpoint answered:
  backend-loss survival may still be untested
- Redis is reachable:
  write authority may still be singular

If that sentence cannot be produced, the operator is already slipping into
comfort theater.

## What "good" would eventually look like

This runbook is working when the operator no longer has to privately stitch
together:

- what is running where
- which peer is eligible
- what happens after locality breaks
- what still depends on sacred storage or sacred nodes

That does not mean the dream is finished today.
It means the operator method is aligned with the dream instead of hiding its
absence.

## Bottom line

This repo does not mainly need more commands.
It needs commands to stay honest about what class of truth they belong to.

The core operating principle is:

> do not let authored shape, local health, or one successful path impersonate
> distributed correctness.

If the runbook enforces that discipline, it becomes useful.
If it relaxes that discipline, it becomes just another way to narrate partial
infrastructure as if it were already resilient.
