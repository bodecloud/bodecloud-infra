# Proof Matrix and Drill Catalog

This page is the claim firewall for `bolabaden-infra`.

It exists to answer one question:

> what exact thing are we claiming, what exact drill would prove it, what proof
> class have we reached, and what stronger sentence is still forbidden even if
> that drill passes?

That is the only way to keep a serious-looking stack from overpromoting narrow
evidence into broader closure.

## What this page is and is not allowed to prove

This page is authoritative about:

- proof classes used in this repo
- what each class can honestly support
- which drills would materially move the dream forward
- which ceilings remain closed after a drill passes

This page is not authoritative about:

- whether a specific service already passed a drill
- whether one passing drill upgrades the whole stack
- whether the overall dream is satisfied

## Strongest honest current answer

The repo needs proof discipline more than it needs more confident language.
Most fake closure here follows the same pattern:

1. config exists
2. service starts
3. happy-path request returns `200`
4. surrounding prose starts sounding distributed
5. the operator still carries the decisive truth privately

This page exists to break that sequence.

The deeper reason it exists is that the repo no longer lacks moving parts.
It lacks proof packets that survive the user's actual accusation:

> after all these tools, am I still the hidden control plane?

## Proof classes

### `Intent only`

Meaning:

- the behavior is clearly wanted
- repo-native intent surfaces describe it
- there is no meaningful live proof yet

Allowed claim:

- this is a real target

Forbidden upgrade:

- the runtime already behaves this way

### `Config present`

Meaning:

- the tracked config contains ingredients for the behavior
- the authored system leans in that direction

Allowed claim:

- the implementation has been authored toward this outcome

Forbidden upgrade:

- the behavior exists under live conditions

### `Happy-path runtime`

Meaning:

- the path or service works under nominal conditions
- no relevant wrong-node, backend-loss, or authority stress has been exercised

Allowed claim:

- this path works in normal conditions

Forbidden upgrade:

- the platform preserves it under stress

### `Wrong-node proven`

Meaning:

- a request was intentionally sent through a node that did not host the target
  locally
- receiving-node identity was observed
- backend-node identity was observed
- the request still completed through the correct peer

Allowed claim:

- this exact wrong-node path is real

Forbidden upgrade:

- wrong-node traffic is generically solved

This class matters more than the others because it is the first class that
begins testing the repo's actual promise instead of merely its runtime shape.

### `Fallback-route proven`

Meaning:

- the preferred backend disappeared
- the rescue path remained present
- the request still completed via the intended fallback path

Allowed claim:

- this exact route survives this exact backend-loss scene

Forbidden upgrade:

- HTTP fallback is broadly solved

### `Policy-parity proven`

Meaning:

- local and peer-forwarded versions of one protected route were compared
- auth, middleware, and visible route behavior remained equivalent enough to
  call it the same protected service

Allowed claim:

- this exact protected route preserved policy meaning under this handoff

Forbidden upgrade:

- protected-route continuity is now solved across the stack

### `Stateful authority proven`

Meaning:

- write owner, replica model, promotion path, and client rediscovery behavior
  were explicitly defined and exercised for one service class

Allowed claim:

- this exact stateful authority model is real for this service class

Forbidden upgrade:

- the stateful lane is broadly resilient now

## Drill packet standard

Every serious drill should leave behind a packet containing:

- the route or service class exercised
- the node that received the first hop
- the truth source the receiving node used
- the failure scene exercised
- what the drill proved
- what stronger sentence remains forbidden

Without that packet, the result is still too easy to overread.

The packet should also answer one socially important question:

- which humiliating private sentence stopped being true because of this drill?

If the packet cannot answer that, it is probably still a technical anecdote
instead of real burden transfer.

## Route proof packet template

Use this template when a drill touches HTTP routing, peer forwarding, or
fallback.

| Field | Required evidence |
| --- | --- |
| Route under test | Exact hostname, path, service name, and whether it is public, protected, or internal |
| Intended owner | The node or backend expected to own the service before failure |
| First-hop node | The public node that actually received the request |
| Locality result | Whether the first-hop node served locally or chose a peer |
| Truth source used | The file, registry, generated config, label set, or runtime API consulted |
| Peer chosen | The backend peer selected and why it was eligible |
| Policy surface | Auth, middleware, headers, and trust-boundary behavior observed |
| Failure scene | Wrong-node entry, backend loss, peer loss, route drift, or policy comparison |
| User-visible result | Status, content identity, headers, and whether the request still meant the same thing |
| Private sentence removed | The exact operator sentence the system no longer needs |
| Still forbidden | The larger claim this drill still does not prove |

This template is intentionally verbose.
The repo does not need more impressive success anecdotes.
It needs evidence packets that a later contributor cannot easily overread.

## Stateful proof packet template

Use this template when a claim touches MongoDB, Redis, Postgres-shaped services,
RabbitMQ, Headscale, Qdrant, or any other state-bearing service.

```yaml
stateful_authority_packet:
  claim_tested: "stateful authority under failure"
  service: "<one exact stateful service>"
  authority_before: "<writer/leader/source of truth before failure>"
  failure_introduced: "<exact node, process, disk, network, or backend failure>"
  authority_after: "<writer/leader/source of truth after failure>"
  client_observation: "<what dependent clients saw before/during/after>"
  rediscovery_mechanism: "<DNS, seed list, Sentinel, driver, registry, manual, none>"
  fencing_or_split_brain_guard: "<mechanism, or none>"
  storage_truth: "<replication, backup, snapshot, shared storage, singular disk>"
  operator_intervention_required: true
  result: "pass | fail | honest-singularity | inconclusive"
  what_this_proves: "<one narrow sentence>"
  what_is_still_forbidden: "<larger HA sentence still illegal>"
```

If the result is `honest-singularity`, the packet still has value.
It records that the repo removed ambiguity without pretending the workload
became resilient.

TCP reachability should never be used as a shortcut for this packet.

## Drill matrix

| Claim under test | Minimum drill | Current ceiling if it passes | Stronger sentence still forbidden afterward |
| --- | --- | --- | --- |
| More than one public node can receive traffic | Send traffic to at least two public nodes and observe first-hop identity | First-hop plurality is real | Wrong-node request preservation is solved |
| One stateless HTTP route works locally | Hit a local route such as `whoami` or `wishlist` on its owning node | Local happy path is real | Wrong-node and backend-loss survival are solved |
| One stateless HTTP route survives wrong-node entry | Force first hop onto a non-owner healthy node and observe correct peer completion | This exact wrong-node path is real | Shared placement and peer eligibility are generically solved |
| One stateless HTTP route survives backend loss | Remove the preferred backend after route establishment and observe rescue-path continuity | This exact fallback path is real | HTTP failover is broadly solved |
| One protected route keeps the same policy meaning after handoff | Compare local and peer-forwarded auth, middleware, headers, and visible behavior | This exact protected route preserved semantics | Protected-route continuity is solved across the stack |
| TCP route is reachable | Connect to a TCP-exposed service such as `mongodb` or `redis` through the intended entrypoint | Transport path is real | Stateful resilience or authority transfer is solved |
| One stateful service owns authority correctly under failure | Define and exercise writer, replica, promotion, and rediscovery behavior | This exact authority model is real | The general stateful lane is now anti-SPOF |

## The matrix is really an anti-bluff device

This matrix exists because this repo is especially easy to overread in three
ways:

1. config gets mistaken for behavior
2. behavior gets mistaken for distributed truth
3. one distributed truth gets mistaken for broad platform closure

The matrix should make those jumps feel socially illegal before they feel
technically tempting.

## The highest-leverage next drills

If the goal is to move the dream rather than produce broad-looking evidence,
the next best drills are:

1. one intentional wrong-node stateless HTTP route
2. one backend-loss drill for that same route
3. one protected-route parity comparison
4. one service-class-specific stateful authority packet

That order matters.
It keeps the repo from borrowing confidence from narrower wins.

It also tracks the user's actual pain gradient:

- first prove the wrong node stops acting stupid for one low-state route
- then prove the rescue path survives the failure that made it matter
- then prove policy meaning survives the same style of handoff
- only then spend serious honesty on stateful authority

## What a passed drill still does not mean

This page should keep a few false promotions illegal:

- a local `200` does not mean wrong-node dignity
- one wrong-node success does not mean shared truth is solved
- one fallback success does not mean the platform behaves like one cloud
- TCP reachability does not mean stateful resilience
- a well-designed proof matrix does not mean the runtime has already matured

The matrix is only useful if it keeps those ceilings visible in advance.

## The exact contract every serious drill is testing

Every real drill in this repo should be reducible back to this:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

The drill is only interesting if it tests one uncertain part of that contract:

- first-hop plurality
- locality determination
- peer selection
- fallback durability
- policy continuity
- state authority

Otherwise the drill may still be useful, but it is not yet testing the dream
the user is actually asking the repo to earn.

## Evidence naming convention

When a proof packet is recorded, name it by the behavior it actually proves,
not the ambition it gestures toward.

Good names:

- `whoami-wrong-node-http-proof`
- `wishlist-backend-loss-http-proof`
- `dozzle-protected-route-policy-parity-proof`
- `redis-tcp-reachability-only`
- `headscale-state-authority-not-proven`

Weak names:

- `ha-test`
- `failover-working`
- `cluster-proof`
- `zero-spof-validation`

The name should make the ceiling visible before anyone opens the packet.
