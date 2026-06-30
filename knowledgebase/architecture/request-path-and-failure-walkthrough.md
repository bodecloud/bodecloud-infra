# Request Path and Failure Walkthrough

This page exists because "HA routing" becomes fake almost immediately unless it
is forced back down to one literal request through one literal sequence of
decisions.

The user is not asking for a more elegant category label.
They are asking:

> when one hostname is requested, what literally happens next, who knows each
> fact, where does that knowledge come from, which part is already real in the
> priority runtime, and where exactly does current proof stop?

That has to be answered like an operator tracing a real path through a real
stack.
If it turns into generic architecture prose, it stops being useful.

## What this page is and is not allowed to prove

This page is authoritative about:

- the target request contract the repo is trying to earn
- which live runtime components already participate in that contract
- the exact seams where the request story becomes aspirational
- why wrong-node dignity is stricter than "the proxy answered"
- why request preservation has to be decomposed into distinct truth checks

This page is not authoritative about:

- generic wrong-node success today
- route-by-route proof for the whole stack
- automatic parity between HTTP, TCP, and stateful surfaces
- pretending a plausible walkthrough is the same thing as a verified drill

This is a request-trace and boundary page.
It is not a success report.

## Strongest honest current answer

The repo already has enough live machinery to describe a serious request path:

- Cloudflare-oriented first-hop plurality
- Traefik as the real L7 execution surface
- TinyAuth and `nginx-traefik-extensions` as policy-bearing edge logic
- CrowdSec as active edge filtering
- Headscale as a real private-mesh assumption
- root-declared TCP routes for services such as MongoDB and Redis

What the repo still does **not** prove is the hardest part of the path:

- that a healthy receiving node which lacks the service locally can preserve
  the request from shared current truth instead of forcing the operator to stay
  the hidden control plane

That seam is the entire point of this page.

## What still does not count as a request-path answer

This repo needs a harsher standard than "the path sounds traceable."

The following still do not count as a real answer to the user's request-path
question:

- naming the components in order
- showing that the first hop can hit more than one node
- proving the local happy path only
- describing a plausible peer-forward story without a shared placement authority
- assuming fallback semantics from helper presence alone
- treating one forwarded request as enough without naming what still remained
  operator-supplied

Those things can make the walkthrough clearer.
They do not yet prove the system has stopped depending on sacred-node or
private-topology memory.

## What a real request-path proof packet would have to contain

If this page ever supports a stronger claim than "the repo understands the
seam clearly," it should be because a real proof packet exists.

For a stateless protected HTTP route, that packet would need artifacts like:

- the exact hostname or route class exercised
- the receiving-node identity
- the backend-node identity
- the source of locality or placement truth used for the decision
- the evidence that policy and auth remained the same after handoff
- the evidence that the route survived the failure condition being claimed
- the explicit sentence about what broader route classes were still not proven

Without a packet like that, a clean walkthrough is still analysis, not route
ownership proof.

## Read this page with the correct standard

The standard here is not:

> can traffic reach something?

The standard is:

> did the system preserve the request itself without the operator privately
> reconstructing the answer?

That means preserving:

- intended destination
- route class
- locality truth
- peer eligibility truth
- middleware and auth meaning
- backend identity
- recovery behavior under local failure

Many systems preserve reachability.
Far fewer preserve request meaning once the request lands on the wrong node or
the local backend disappears.

## Primary evidence for this page

Use these together:

1. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
3. [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
4. [Current Compose Runtime](current-compose-runtime.md)
5. [HA, Failover, and Routing](ha-failover-routing.md)
6. [Operator Contract and Success Criteria](operator-contract.md)
7. [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
8. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
9. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

Those sources together answer:

- what the target contract is
- what the live edge surface really is
- where the live request story currently stops
- why fallback and state stay under stricter language

## The target request contract

The clearest repo-native contract is already stated in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That is not a proof statement.
It is the target contract this walkthrough is cross-examining.

The key thing to notice is that the contract has three distinct decision
moments:

1. first-hop node selection
2. local-versus-remote service determination
3. semantically safe rescue when locality is absent

Most weak docs only narrate step `1`.
The user is angry about steps `2` and `3`.

## The request class we should trace first

Use a stateless or near-stateless protected HTTP route first, because that is
the place where the repo can plausibly earn honest progress before stateful
surfaces do.

A concrete class example is:

- `https://dozzle.$DOMAIN`

That route is useful because it is:

- public HTTP
- protected by edge policy
- already present in the root runtime
- easy to describe in terms of locality, middleware, and wrong-node behavior

The goal is not "prove Dozzle specifically right now."
The goal is to pick a route class where the whole user complaint becomes
visible.

## Step 1: first-hop plurality

What happens first in the dream:

1. Cloudflare resolves the requested public name
2. the request lands on one healthy public node

What the live stack already proves:

- `cloudflare-ddns` is part of the active edge fragment
- the repo explicitly wants more than one public node to be a valid first hop
- the planning layer explicitly warns that Cloudflare DDNS presence is not the
  same thing as full multi-node request failover

What this step honestly buys:

- public-node plurality as a first-hop pressure
- a serious anti-sacred-entrypoint intention

What it does **not** buy:

- service locality truth
- peer eligibility truth
- preserved request semantics after wrong-node entry

This distinction matters because "multiple A records" is not one step away
from the dream.
It is one early ingredient of the dream.

## Step 2: local edge-stack execution

Once the request reaches a node, the local edge stack becomes responsible for
the next decision.

In the live runtime, that edge stack already includes:

- `traefik`
- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- `dockerproxy-ro`
- `dockerproxy-rw`
- file-provider and Docker-provider style route material in the proxy fragment

That proves the repo already treats request correctness as more than just:

- port reachability
- TLS termination
- "some reverse proxy exists"

The request already depends on policy-bearing surfaces.

That is why wrong-node success has to preserve more than transport.
If the fallback path changes the meaning of auth, middleware, or headers, the
user will experience it as a different service even if it still returns `200`.

## Step 3: local-versus-remote truth

This is the first truly difficult step.

The receiving node has to determine whether the service is actually local.

The live runtime already provides:

- local Docker labels
- local Docker provider visibility
- rich route material for HTTP and TCP services

The repo also repeatedly expresses a desire for:

- a shared current-state registry such as `services.yaml`

What the worktree does **not** yet prove:

- a live root `services.yaml` is present and consumed
- another explicit shared placement-truth surface outranks operator memory

That means the current best honest sentence is:

> the repo already knows that locality truth must become shared, but the
> priority runtime does not yet prove that the receiving node owns that truth
> from an active shared registry.

This is why the user keeps refusing generic "more load balancing" answers.
Without locality truth, the wrong-node machine is still epistemically weak.

## Step 4: local happy path

The flattering path through the stack is:

1. Cloudflare lands the request on node A
2. node A is healthy
3. Traefik and related edge logic are healthy on node A
4. the target service is local to node A
5. auth and middleware execute locally
6. the service answers locally

This path is important because it proves something real:

- the stack already has serious local route execution
- policy-bearing local routes can be real
- the platform is not a toy

But the flattering path is not the wound.

A locally served route does **not** prove:

- another node could preserve the same request
- wrong-node service discovery works
- peer fallback preserves the same semantics
- the rescue path survives local backend disappearance

This repo only starts answering the user's real complaint when locality fails.

## Step 5: wrong-node entry

Now trace the real scene:

1. Cloudflare lands the request on node A
2. node A is healthy enough to receive traffic
3. the service actually lives on node B
4. node A must preserve the request without private operator rescue

For that to be honestly solved, node A must know all of the following:

1. the service is not local
2. which peer currently owns a valid backend
3. whether that peer is healthy enough for this route class
4. whether auth and middleware semantics survive handoff
5. whether the route needed for rescue still exists under failure

That is already much stricter than:

- "there is another node"
- "Headscale can reach another node"
- "Traefik is present"

The repo currently has meaningful pieces of this scene, but not closure.

## What the worktree already gives the wrong-node scene

### Strong pieces already present

- the dream itself is explicit and repeated in repo-native intent surfaces
- the runtime already includes a serious public edge surface
- the runtime already includes private-mesh assumptions via Headscale
- the runtime already includes protected HTTP surfaces and TCP surfaces
- the planning layer already names the missing registry and failover gaps

### Weak or still-unproven pieces

- no live root `services.yaml` or equivalent shared placement surface is proven
- peer eligibility rules are not yet proven as tracked current truth
- `docker-gen-failover` is explicitly recorded as deleting routes when
  containers stop
- no generic route-by-route wrong-node proof is present
- middleware and auth continuity across peer forwarding are not yet broadly
  proven

This means the repo can honestly claim strong directional convergence.
It cannot honestly claim generic wrong-node dignity yet.

## Step 6: backend disappearance after wrong-node pressure

This is where many "failover" stories die.

Take the already difficult wrong-node scene and add one more condition:

- the preferred backend disappears or becomes unavailable

Now the receiving node needs not only placement truth but route durability.

The repo's own planning material records the most important negative evidence
here:

- `docker-gen-failover` is present
- `docker-gen-failover` can delete routes when containers stop
- automated service failover between nodes is still missing

That means the current route story weakens exactly where the user cares most:

> the recovery path may not survive the event that made recovery necessary.

This is why fallback-route durability has to be treated as its own separate
truth, not as part of "proxy health."

## Step 7: policy continuity after handoff

Even if a request can be forwarded, a protected route still has another hard
requirement:

- it must keep behaving like the same protected route

In this stack, that means handoff has to preserve at least:

- TinyAuth behavior
- `nginx-traefik-extensions` forward-auth behavior
- middleware ordering
- CrowdSec filtering implications
- visible route semantics such as redirects, headers, and path handling

This matters because the user is not asking for:

- "can a byte stream still move?"

They are asking for:

- "can the same service still exist meaningfully after wrong-node entry?"

If forwarding silently strips or changes the route's protection model, that is
not real success.

## Step 8: why this page refuses to merge HTTP, TCP, and state

The repo has to keep three separate request classes visible.

### Stateless or near-stateless HTTP

This is the first lane that can plausibly earn honest wrong-node drills.

Why:

- route identity is visible
- middleware and auth can be compared
- backend identity can be observed

### Raw TCP

The root runtime already includes TCP exposure for services such as MongoDB and
Redis via Traefik TCP routers.

That proves TCP exists at the edge.
It does **not** prove:

- TCP wrong-node preservation is solved
- TCP fallback semantics are trustworthy
- state-bearing consequences are acceptable

### Stateful services

For stateful services, request continuity is not enough.
The operator still needs answers for:

- write authority
- replication direction
- promotion behavior
- reconnect semantics
- storage truth

That is why stateful routes stay under much harsher language everywhere else in
the knowledgebase.

## What a real wrong-node proof would actually need

This page is not performing the drill, but it should make the proof
requirements explicit.

For one stateless protected HTTP route, a real wrong-node proof would need:

1. evidence that the request was intentionally sent to a node that does not
   host the route locally
2. evidence of the receiving node's identity
3. evidence of the answering backend's identity
4. evidence that the route still completed successfully
5. evidence that auth and middleware behavior remained consistent enough to
   still count as the same route

For backend-loss fallback, add:

6. evidence that the preferred backend actually disappeared
7. evidence that the route required for rescue remained present long enough to
   preserve the request

Until that exists, this page remains a truthful trace plus an honesty wall.

## Bottom line

The current request story is strongest up to:

- first-hop plurality exists as a real pressure
- a serious edge stack exists as a real runtime surface
- locality matters and is already reflected in the live stack

The request story is weakest at the exact point the user cares about most:

- the receiving node still does not clearly prove that it owns enough current
  truth to preserve the request when locality fails and the backend situation
  gets worse

That is why this repo keeps circling the same problem.
It is not missing proxies.
It is missing system-owned request truth.
