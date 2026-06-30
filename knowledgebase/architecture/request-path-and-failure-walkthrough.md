# Request Path and Failure Walkthrough

This page exists because "HA routing" becomes fake almost immediately unless it
is forced back down to one literal request.

The user is not asking for a better category label.
They are asking:

> when one hostname is requested, what literally happens next, who knows each
> fact, where does that knowledge come from, and where exactly does current
> proof stop?

That has to be answered like an operator tracing a real path through a real
stack.
If it turns into generic architecture prose, it stops being useful.

## Strongest honest current answer

The repo already has enough edge machinery to describe a serious request path.

It does **not** yet prove the hardest part of the request path:

- that a healthy receiving node which lacks the service locally can preserve
  the request from shared current truth instead of forcing the operator to stay
  the hidden control plane

That is the seam this page exists to expose.

## What this page is and is not allowed to prove

This page is allowed to:

- reconstruct the intended request contract from repo-native sources
- show which components already participate in that contract
- identify the exact seams where the request story becomes aspirational
- explain why wrong-node success is stricter than "a proxy answered"

This page is not allowed to:

- claim generic wrong-node success today
- treat a plausible walkthrough as runtime proof
- merge HTTP, TCP, and stateful failover into one maturity story
- imply that middleware, auth, locality, peer eligibility, and application
  behavior have already survived a real failure drill

## Read this page with the right standard

The standard here is not:

> can traffic reach something?

The standard is:

> did the system preserve the request itself without the operator privately
> reconstructing the answer?

That is the user's real benchmark.

Many systems can preserve reachability.
Far fewer can preserve request meaning once the request hits the wrong node or
the local backend disappears.

## Primary evidence for this page

Use these together:

1. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. [`compose/docker-compose.coolify-proxy.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.coolify-proxy.yml)
3. [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md)
4. [`ha-failover-routing.md`](ha-failover-routing.md)
5. [`failure-model-and-maturity.md`](failure-model-and-maturity.md)
6. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
7. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

Those together let this page answer:

- what the target contract is
- what the live edge surface really is
- where the live edge story still stops

## The target request contract

The clearest repo-native contract is already stated in
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md):

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That is not yet a proof statement.
It is the target contract this walkthrough is cross-examining.

## The request we are tracing

Use a stateless HTTP service first, because that is where the repo is most
likely to earn honest progress before anything stateful does.

Conceptually:

```text
client -> https://dozzle.bolabaden.org
```

Or in operator terms:

```bash
curl -H "Host: dozzle.bolabaden.org" https://bolabaden.org
```

The concrete hostname is not the point.
The class is:

- public HTTP
- auth-bearing or middleware-bearing
- possibly local on one node and remote on another

That is the most revealing request class for the user's actual dream.

## What the live stack already proves is present

From the root runtime and active proxy fragment, a real HTTP request can
already involve all of these components:

- Cloudflare DNS assumptions at the edge
- `publicnet`, `backend`, and `warp-nat-net`
- `traefik`
- `crowdsec`
- `crowdsec-init`
- `tinyauth`
- `nginx-traefik-extensions`
- `dockerproxy-ro`
- `dockerproxy-rw`
- `docker-gen-failover`
- the target container itself

That matters because the system is already complex enough that "the proxy is
up" is almost meaningless as a status sentence.

The real question is not whether the parts exist.
It is which parts already carry shared truth and which parts still depend on
assumptions that disappear on the bad day.

## Walkthrough 1: the flattering path the repo is most likely to support

Scenario:

1. Cloudflare resolves the public hostname to a healthy public node.
2. The receiving node is reachable.
3. Traefik is healthy on that node.
4. The requested service is actually local to that same node.
5. The relevant auth and middleware chain also exists and behaves correctly.

What this path likely means in today's repo:

- DNS can plausibly land on a valid public entrypoint
- Traefik can plausibly route the request
- middleware can plausibly execute on the local path
- the target container can plausibly answer locally

What this path still does **not** prove:

- that another healthy node could have answered the same request with the same
  behavior
- that peer selection is based on live placement truth rather than convention
- that the same route survives if the local backend is the thing that fails

This is an important distinction.
Happy-path locality is not the user's real pain.
Wrong-node dignity is.

## Walkthrough 2: the real scene the user cares about

Scenario:

1. Cloudflare sends the request to node A.
2. The target service actually lives on node B.
3. Node A is healthy enough to accept traffic.
4. Node A is not healthy enough to satisfy the request locally because the
   service is absent there.
5. The system must now preserve the request without private operator rescue.

For this scene to be honestly solved, node A must know all of the following:

1. that the target service is not local
2. which peer currently owns a valid instance
3. whether that peer is healthy enough for this route class
4. whether auth and middleware meaning survive the handoff
5. whether the fallback path still exists after the failure that made fallback
   necessary

That is the real contract.

## What the current worktree gives this scene

The current worktree gives meaningful pieces, but not the full contract.

### What is strong already

- the root runtime is modular and serious, not a toy stack
- Traefik is a real edge control surface
- auth and middleware are already part of request correctness
- there is explicit repo desire for local-first then peer-forward behavior
- there is explicit desire for any-node public entry rather than one sacred box

### What is still too weak

- no tracked root `services.yaml` or equivalent live placement registry is
  proven active
- `docker-gen-failover` is specifically called out as buggy because it deletes
  routes on container stop
- peer eligibility rules are not yet proved as tracked shared truth
- wrong-node success is not yet proved route-by-route
- middleware and auth continuity across peer fallback are not yet proved

That means the system can be described as strongly oriented toward the desired
request contract.
It cannot yet be honestly described as having earned that contract.

## Where the request story currently breaks

The cleanest way to say it is:

the request story is strongest up to "a serious edge stack exists" and weakest
at "the receiving node owns enough current truth to preserve the request when
locality fails."

That seam matters more than any individual product choice.

The request story currently weakens at these exact points:

1. placement truth
   The repo wants current-state routing, but the root runtime does not yet
   prove a durable shared placement surface.
2. peer eligibility truth
   The repo does not yet prove a shared answer to "which peer is valid now for
   this route?"
3. route persistence under backend loss
   The helper most associated with failover generation is already documented as
   unreliable in the failure that matters.
4. policy continuity
   The repo clearly cares about auth and middleware, but has not yet proved
   they survive peer forwarding generically.

Those four things are what turn "interesting ingress stack" into "actually
preserved request."

## Why this page refuses to merge HTTP, TCP, and state

The user is not merely asking for packets to move.
They are asking for service meaning to survive.

That is why this page has to keep classes separate:

- stateless HTTP can plausibly be the first real wrong-node proof target
- raw TCP is harder because request semantics are thinner and client behavior
  is harsher
- stateful systems are harder still because route continuity does not equal
  authority continuity

If those classes are merged, the docs become smoother and less honest.

## What would count as real proof

This page is not the proof.
It should still define what proof would count.

For one stateless HTTP route, real proof would mean:

1. the service is intentionally absent from the receiving node
2. the request lands on that node anyway
3. the node selects the correct healthy peer from tracked truth
4. auth and middleware still behave as expected
5. the application response remains semantically the same service
6. the path still survives local-backend-loss conditions relevant to the route

Anything weaker than that may still be useful progress.
It is not yet the user's real benchmark.

## Bottom line

The repo already has a serious request model and a serious edge surface.
It also has enough explicit architecture intent to make the dream concrete:

- any-node entry
- local-first service when honest
- peer-forward when necessary
- no fake merging of HTTP, TCP, and stateful claims

What the repo still does not prove is the part the user actually keeps asking
about:

that the wrong node can preserve the request from shared current truth instead
of forcing the operator to remain the hidden control plane.
