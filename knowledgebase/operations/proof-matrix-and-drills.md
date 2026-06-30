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

## The highest-leverage next drills

If the goal is to move the dream rather than produce broad-looking evidence,
the next best drills are:

1. one intentional wrong-node stateless HTTP route
2. one backend-loss drill for that same route
3. one protected-route parity comparison
4. one service-class-specific stateful authority packet

That order matters.
It keeps the repo from borrowing confidence from narrower wins.

## What a passed drill still does not mean

This page should keep a few false promotions illegal:

- a local `200` does not mean wrong-node dignity
- one wrong-node success does not mean shared truth is solved
- one fallback success does not mean the platform behaves like one cloud
- TCP reachability does not mean stateful resilience
- a well-designed proof matrix does not mean the runtime has already matured

The matrix is only useful if it keeps those ceilings visible in advance.
