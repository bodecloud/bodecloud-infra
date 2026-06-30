# Current-State Registry and Peer Eligibility

This page defines the `services.yaml` idea without pretending the tracked
runtime already ships it.

It should be read alongside:

- [HA, Failover, and Routing](ha-failover-routing.md)
- [Request Path and Failure Walkthrough](request-path-and-failure-walkthrough.md)
- [The Missing Middle Layer](missing-middle-layer.md)
- [Capability Gaps and Roadmap](capability-gaps-and-roadmap.md)
- [Proof Matrix and Drill Catalog](../operations/proof-matrix-and-drills.md)

## The shortest honest answer

The repo keeps circling a lightweight current-state registry because the
wrong-node problem cannot be solved from proxy syntax alone.

The receiving node needs to know:

- which service was requested
- whether that service is local right now
- if it is remote, which peer currently owns it
- whether that peer is eligible for this route class
- whether the fact is fresh enough to trust
- which artifact proves the decision after the fact

That is what `services.yaml` keeps trying to become.

It is not currently a proven runtime feature.
A repository search does not show a tracked root `services.yaml` consumed by
the priority Compose runtime.

## What this page is and is not allowed to prove

This page is allowed to define:

- what a current-state registry would need to mean in this repo
- why desired-state Compose files are not enough by themselves
- why peer eligibility is stricter than node reachability
- what packet fields would prove one wrong-node decision used shared truth
- how to keep a registry from becoming prettier operator folklore

This page is not allowed to prove:

- that root `docker-compose.yml` already consumes this registry
- that wrong-node routing already works generically
- that a file named `services.yaml` would automatically solve fallback
- that all service classes can share the same eligibility rules
- that a registry replaces stateful authority, fencing, or client rediscovery

## Why desired state is not the truth the user is asking for

Compose is still the priority authoring surface.
That does not mean Compose owns current distributed truth.

Compose can say:

- this service is defined here
- this container should exist on this host
- this label should produce this local route
- this healthcheck should evaluate this local process

The wrong-node event asks a different question:

> this request landed here right now; where is the correct service authority
> now, and why is that peer safe for this request?

That is current-state truth.
It can be informed by Compose.
It cannot be replaced by Compose.

## The three truths a registry must not blur

### 1. Placement truth

Placement truth answers:

> where does this service actually live right now?

Minimum useful fields:

```yaml
placement:
  service: "whoami"
  route_class: "stateless-http"
  owner_node: "node1"
  backend_url: "http://node1.tailnet:8080"
  source: "manual-services-yaml | docker-runtime | opensvc | nomad | generated"
  observed_at: "2026-06-30T00:00:00Z"
  freshness_seconds: 30
  convergence_state: "matched | stale | disputed | unknown"
```

Placement truth does not prove the peer is safe.
It only says where the service appears to live.

### 2. Peer eligibility truth

Peer eligibility answers:

> may this receiving node hand this route to that peer without changing the
> meaning of the request?

Minimum useful fields:

```yaml
peer_eligibility:
  peer: "node1"
  route_class: "stateless-http"
  transport: "headscale | lan | public | cloudflare-tunnel | unknown"
  health_source: "container-health | blackbox | proxy-health | unknown"
  policy_profile: "public | protected | tcp | stateful"
  policy_converged: true
  version_match: "matched | drifted | unknown"
  eligible: true
  reason: "healthy public stateless HTTP backend with matching policy profile"
```

Reachable is not eligible.
Reachable only means packets can get there.

Eligibility means the peer is appropriate for this exact route class.

### 3. Decision truth

Decision truth answers:

> what did the receiving node decide, and can an operator inspect why?

Minimum useful fields:

```yaml
routing_decision:
  hostname: "whoami.bolabaden.org"
  entry_node: "node2"
  locality_result: "remote"
  placement_source: "services.yaml"
  selected_peer: "node1"
  peer_eligibility_reason: "healthy stateless HTTP owner"
  fallback_order: ["node1", "node3"]
  policy_chain: "public"
  result: "forwarded"
  explanation_artifact: "/var/log/bolabaden-router/decision.jsonl"
```

Without decision truth, a successful request can still be a lucky request.

## The route classes need different eligibility rules

One registry format can hold multiple route classes.
One eligibility rule cannot safely govern all of them.

| Route class | Placement answer | Eligibility answer | Still forbidden |
| --- | --- | --- | --- |
| Stateless HTTP | which node serves this app now | peer is healthy and route meaning is preserved | protected, TCP, or stateful claims |
| Protected HTTP | which node serves this protected app now | auth, middleware, headers, and trust boundary match | "page loaded" as policy proof |
| Raw TCP | which node exposes the TCP service now | transport path is valid for that protocol | stateful correctness |
| Stateful | who owns authority, not merely where a port answers | writer, promotion, fencing, storage, and rediscovery are valid | registry-driven HA without authority proof |

This is why the first useful registry proof should be a boring stateless HTTP
route.

## What `services.yaml` must not become

A registry would fail this repo if it becomes:

- an aspirational desired-state file that operators update when they remember
- a prettier version of the same private topology memory
- a static backend list with no freshness or convergence evidence
- a DNS inventory mistaken for service ownership
- a Traefik label mirror mistaken for cross-node truth
- a place where stateful services are made to look portable because their names
  fit in a table

The danger is not that `services.yaml` is too simple.
The danger is that it can look exactly like the missing middle while still
leaving the hard part in one human's head.

## Minimum useful v1

The smallest honest v1 is not a full orchestrator.
It is one consumed placement surface for one low-risk route.

Minimum v1:

1. one stateless HTTP service selected as the test subject
2. one explicit owner node
3. one non-owner entry node
4. one shared placement artifact or runtime-generated equivalent
5. one peer eligibility record
6. one route decision artifact
7. one wrong-node drill proving the receiving node used the shared truth
8. one forbidden-claims note saying what the drill still does not prove

Good first candidates:

- `whoami`
- `mkdocs`
- `wishlist`

Bad first candidates:

- `redis`
- `mongodb`
- `headscale`
- `portainer`
- `code-server`
- anything where auth, state, or admin power will distract from placement truth

## `placement_decision_packet`

Use this packet when testing whether a route used shared current-state truth:

```yaml
placement_decision_packet:
  claim_tested: "wrong-node placement decision from shared current-state truth"
  route: "whoami.bolabaden.org"
  route_class: "stateless-http"
  entry_node: "node2"
  expected_owner: "node1"
  placement_source:
    type: "services.yaml | generated-runtime-state | opensvc | nomad | other"
    path_or_endpoint: "<file path or API endpoint>"
    observed_at: "<timestamp>"
    freshness_seconds: 30
    convergence_state: "matched | stale | disputed | unknown"
  locality_result: "local | remote | unknown"
  selected_peer: "node1"
  peer_eligibility:
    health: "healthy | unhealthy | unknown"
    policy_profile: "public | protected | tcp | stateful"
    version_match: "matched | drifted | unknown"
    eligible: true
    reason: "<why this peer was safe for this route class>"
  request_result:
    status: 200
    service_identity: "<body/header/probe proving the intended service answered>"
    preserved_meaning: true
  explanation_artifact: "<log, generated config, decision file, or command output>"
  what_this_proves: "<one exact route used shared placement truth once>"
  what_is_still_forbidden: "<generic wrong-node, protected-route, TCP, or stateful claim>"
```

If this packet cannot be filled out, the repo may still have useful routing
ideas.
It does not yet have proof that placement truth left private memory.

## How this relates to `route_packet`

`route_packet` proves the request behavior.
`placement_decision_packet` proves the placement decision behind it.

They overlap on purpose.

The important distinction is:

- `route_packet`: did the request still mean the same thing?
- `placement_decision_packet`: did the receiving node use shared current truth
  to choose the peer?

A good wrong-node proof needs both.

## The honest current ceiling

The current worktree proves that the repo deeply understands the need for this
layer.

It does not yet prove that the layer is live.

The strongest current sentence is:

> `services.yaml` or an equivalent current-state registry is one of the
> clearest missing middle candidates because it attacks placement folklore
> directly, but the priority Compose runtime has not yet earned the right to
> treat it as live authority.

That sentence is less satisfying than saying "we have a registry."
It is also the sentence that keeps the dream intact.

## Bottom line

The registry is not the dream.

The dream is that the wrong healthy node can stop asking one human:

> where does this service really live, which peer is safe, and what did I just
> decide?

A current-state registry is valuable only if it makes that question
inspectable, fresh, and boring for at least one real route.
