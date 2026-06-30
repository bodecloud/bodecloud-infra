# Proof Matrix and Drill Catalog

This page is the claim firewall for `bolabaden-infra`.

It exists for the question the repo is most likely to lie about by accident:

> what exact thing are we claiming to have solved, what exact drill would prove
> it, what exact proof class have we actually reached, and what important part
> of the dream would still remain unproven even after that drill passes?

That is the real proof problem here.

The user is not asking for a generic checklist.
They are asking whether several ordinary Docker nodes can stop behaving like a
set of separate local stacks with nicer branding and start behaving like one
request-preserving, operator-legible personal cloud without immediately
forcing surrender to Swarm, Kubernetes, or some other worldview-heavy control
plane.

This page exists so green-looking evidence does not quietly upgrade itself
into broader closure than it actually earned.

That danger has to stay emotionally explicit.
The repo is no longer at risk of only looking too primitive.
It is now equally at risk of looking mature enough that narrow evidence starts
getting socially promoted into a larger feeling of completion.

## What this page is and is not allowed to prove

This page is authoritative about:

- what proof classes exist in this repo
- what each class can honestly support
- which stronger ceilings must remain closed after a drill passes
- which exact drills would materially move the dream forward
- which kinds of passing results are still too narrow to widen into stack-wide
  confidence

This page is not authoritative about:

- whether a specific service has already passed a given drill
- whether one passing drill upgrades the whole stack
- whether the overall dream is satisfied just because a proof class exists

This is a proof-boundary page, not a success report.

## Strongest honest current answer

The repo needs proof discipline more than it needs more confident language.
A passed drill only matters if it makes one specific fragment of the dream more
true and keeps the next unclosed ceiling visible.

The danger is not missing tests.
The danger is serious-looking evidence being quietly promoted into broader
architecture closure than it actually earned.

That is why this page cannot read like a standard QA matrix.
The user is not just asking whether routes can be exercised.
The user is asking whether the proof discipline is strong enough to resist the
ecosystem's favorite lie: one convincing success signal equals one meaningful
option.

## The thing this page is trying to stop

Most fake closure in this repo follows a repeating pattern:

1. config exists
2. the service starts
3. a happy-path request returns `200`
4. the surrounding prose starts sounding distributed
5. the operator's hidden burden remains almost unchanged

This page exists to break that sequence.

It also exists to make the break a little uncomfortable.
If a drill result feels satisfying before it has named the next still-closed
ceiling, then the page is already helping fake closure along.

It forces every claimed success to answer:

- what exact path or topology was exercised?
- what exact failure class was exercised?
- what stronger sentence is still forbidden?
- which hidden operator burden still survives even after the drill passes?

If those questions disappear, the site becomes another fake HA story.

## The dream fragments being measured

This matrix only makes sense if the dream stays visible.

The dream is not:

- more nodes
- more YAML
- more healthchecks
- more cluster vocabulary

The dream is:

> any surviving public node can receive the request, determine whether the
> target is local, preserve the request if it is not, survive the failure that
> made fallback necessary, keep policy and auth coherent on the fallback path,
> and do all of that without the operator's head secretly remaining the real
> control plane

Each proof class below should be read against that standard.

That means the matrix is not primarily measuring technical cleverness.
It is measuring how much of the operator's private explanatory burden the
system has really taken over, one narrow lane at a time.

## The proof classes

Use these classes consistently and conservatively.

### `Intent only`

Meaning:

- the behavior is clearly wanted
- repo-native intent surfaces describe it directly
- there is no meaningful live proof yet

Allowed claim:

- this is a real target

Forbidden upgrade:

- the runtime already behaves this way

### `Config present`

Meaning:

- the tracked config contains ingredients for the behavior
- the authored system is leaning in that direction

Allowed claim:

- the implementation has been authored toward this outcome

Forbidden upgrade:

- the behavior now exists under live conditions

### `Happy-path runtime`

Meaning:

- the path or service works under nominal conditions
- no relevant wrong-node, backend-loss, or state-topology stress has been
  exercised

Allowed claim:

- this path works in normal conditions

Forbidden upgrade:

- the platform now preserves this path under wrong-node or backend-loss stress

### `Wrong-node proven`

Meaning:

- a specific request was intentionally sent through a node that did not host
  the target locally
- receiving-node identity was observed
- backend-node identity was observed
- the request still completed through the correct peer

Allowed claim:

- this exact wrong-node path is real

Forbidden upgrade:

- the platform now generically preserves wrong-node traffic

### `Fallback-route proven`

Meaning:

- a known-good path was exercised
- the preferred local backend actually disappeared or was stopped
- the route required for recovery remained present and usable long enough to
  preserve the request

Allowed claim:

- this exact route survived backend loss in this exact drill

Forbidden upgrade:

- fallback is broadly solved everywhere

### `Semantic continuity proven`

Meaning:

- local and fallback behavior were compared where users would actually notice
- auth, middleware, headers, and visible policy remained meaningfully the same

Allowed claim:

- this path stayed the same protected service under handoff

Forbidden upgrade:

- all protected routes now have parity guarantees

### `Stateful topology proven`

Meaning:

- ownership, replication, promotion, reconnect behavior, and storage truth
  were defined for one stateful class
- a real failure drill exercised that exact topology

Allowed claim:

- this exact stateful topology has passed this exact class of drill

Forbidden upgrade:

- the platform is now generically HA

## The anti-widening rule

This is the most important rule on the page:

every passed drill must still name the stronger sentence that remains illegal.

Examples:

- a wrong-node drill for one protected HTTP route does not prove all HTTP routes
- a fallback-route drill does not prove policy continuity
- a policy-parity drill does not prove route durability
- a TCP drill does not prove stateful authority
- a Headscale owner-loss drill does not prove general mesh or app failover

If a page stops naming the next still-closed ceiling, it is already starting to
overpay the evidence.

That overpayment is exactly what most surrounding infrastructure discourse
rewards.
This repo has to keep being ruder than that on purpose.

## What still does not count as proof in this repo

This page should also make the false-proof patterns explicit.

The following still do not count as meaningful proof of the user's dream:

- a service starting successfully
- one healthy local `200` response
- a green healthcheck plus calm logs
- multiple public nodes existing at the same time
- a fallback config file existing on disk
- a route generator that looks dynamic until the preferred backend actually dies
- a TCP connection succeeding without any service-authority or topology truth

Those may all be useful ingredients.
They are still weaker than proof that the system owns more of the important
truth on the bad day.

That last phrase is the real retrieval target for the page.
Not "did the route answer?"
Not "did the topology look plausible?"
But "did the system itself own more of the decisive truth after this drill than
it owned before?"

This matters because the whole repo is vulnerable to one specific lie:

- a path worked once
- therefore the platform feels distributed

That is exactly the lie this page is meant to prevent.

And it is a tempting lie precisely because the repo is now rich enough to make
partial success feel emotionally expensive to keep qualifying.
This page exists so the repo keeps qualifying it anyway.

## What a passing drill should leave behind

Every serious drill in this repo should produce a proof packet, not just a
result.

That packet should include:

- the exact route or service class exercised
- the exact topology before the drill
- the exact failure or stress introduced
- receiving-node and backend-node identity where relevant
- the visible user-facing result
- the stronger sentence that still remains illegal afterward
- the surviving hidden operator burden, if any

If a drill cannot leave that packet behind, it is too easy for later prose to
inflate it.

This repo does not just need more tests.
It needs drills whose outputs are narrow enough to stay honest and rich enough
to resist flattering reinterpretation.

## The current matrix for the priority implementation

| Dream fragment | Current proof class | Strongest current anchors | Exact drill needed next | What still remains unproven even after that drill passes |
| --- | --- | --- | --- | --- |
| Any-node public entry is a real target | `Config present` | `.github/copilot-instructions.md`, `README.md`, `cloudflare-ddns`, master plan | Prove more than one public node can receive the first hop | That the right service is preserved after wrong-node entry |
| Local-first service is real | `Intent only` | intent surfaces plus current Traefik-centered local routing model | Pick one stateless protected HTTP route and prove local serve from the honest hosting node with backend identity evidence | That wrong-node forwarding works when locality is absent |
| Wrong-node stateless protected HTTP requests succeed | `Intent only` | routing philosophy plus serious edge runtime | Intentionally land one protected HTTP route on the wrong node and prove receiving-node identity, backend-node identity, route success, and auth continuity | That the route also survives when the preferred backend disappears |
| Fallback route survives backend loss | `Intent only` | master plan explicitly records route-persistence risk around `docker-gen-failover` | Start from a known-good route, remove the preferred backend, and observe whether the fallback route remains present and usable | That auth and middleware semantics stayed identical after handoff |
| Middleware and auth continuity survive peer handoff | `Config present` | TinyAuth, Nginx auth extensions, CrowdSec, Traefik middleware surfaces | Compare one protected route locally versus through intentional peer handoff and verify parity | That all protected routes now share the same continuity guarantees |
| Placement truth is live and shared | `Intent only` | repeated `services.yaml` pressure in README, master plan, and architecture pages | Introduce or expose one real placement-truth surface and prove routing or eligibility logic consumes it | That convergence, restart semantics, and drift handling are trustworthy |
| Peer eligibility truth is live | `Intent only` | Headscale runtime plus peer-broadcast and sync direction in the master plan | Prove a receiving node chooses a peer from current tracked truth rather than folklore | That the same truth remains valid under revision drift and secret drift |
| TCP forwarding works for one named service | `Config present` | current TCP routing such as MongoDB and Redis in root Compose | Exercise a specific TCP path end to end and prove transport plus backend identity | That the corresponding stateful semantics are safe |
| Headscale control-plane failover is real | `Intent only` | live Headscale runtime plus leader-election and replication plan | Define exact topology, then run an owner-loss drill with continuity checks | That broader mesh or application failover is solved |
| One stateful service is honestly resilient | `Intent only` | service-specific plan or runtime plus explicit topology | Define ownership, replicas, promotion, reconnect, and storage truth; then exercise failure | That other stateful classes are similarly mature |

## Drill levels

These levels are intentionally strict.
They are designed to stop partial evidence from sounding final.

### Drill level 1: authoring coherence

Questions answered:

- does the merged graph resolve?
- does the authored stack still render coherently?

Typical evidence:

```bash
docker compose config --quiet
python3 -m mkdocs build -f mkdocs.yml --strict
```

Necessary?
Yes.

Sufficient for the user's real benchmark?
Almost never.

### Drill level 2: happy-path runtime

Questions answered:

- does the named local route or service work under nominal conditions?
- which node or backend answered?

Typical evidence:

- route request
- local logs
- backend identity evidence
- health status

This still does not prove wrong-node behavior or backend-loss recovery.

### Drill level 3: wrong-node path proof

Questions answered:

- was the request intentionally received by a non-hosting node?
- did the receiving node preserve the request by selecting the correct peer?
- can backend identity and receiving-node identity both be shown?

Typical evidence:

- forced first-hop placement
- route success
- receiving-node evidence
- backend-node evidence

This still does not prove route durability under backend loss.

### Drill level 4: backend-loss fallback proof

Questions answered:

- did the preferred backend actually disappear?
- did the rescue route remain present and usable?
- did the request keep working through the failure window?

Typical evidence:

- before/after backend state
- route-material or control-surface evidence
- successful request during degraded topology

This still does not prove policy parity unless that comparison is also made.

### Drill level 5: semantic continuity proof

Questions answered:

- did the fallback path keep the same protected-service meaning?
- did auth, middleware, headers, and visible policy remain meaningfully the
  same?

Typical evidence:

- local versus forwarded comparison
- auth challenge comparison
- middleware/header comparison
- user-visible route parity evidence

This still does not prove stateful safety.

### Drill level 6: stateful topology proof

Questions answered:

- who owns writes?
- how does replication work?
- how does promotion work?
- what reconnect behavior do clients experience?
- did a real failure drill exercise that exact topology?

Typical evidence:

- topology definition
- replication or storage evidence
- failure drill results
- client reconnect or consistency evidence

This is the harshest proof class because it closes the fewest lies at once.

## The exact next-drill logic

When choosing the next drill, prefer the path that removes the most ambiguity
with the least fake widening.

That usually means this order:

1. one stateless protected HTTP route
2. the same route under intentional wrong-node entry
3. the same route under backend loss
4. the same route with explicit policy-parity comparison
5. one named TCP route
6. one stateful topology with explicit authority semantics

Why this order:

- it keeps route identity visible
- it keeps policy continuity visible
- it delays stateful overclaiming
- it produces evidence the rest of the docs can narrate honestly

## The strongest sentences that are still illegal today

Until stronger drills exist, these sentences should remain forbidden:

- "wrong-node traffic is solved"
- "the stack is anti-SPOF now"
- "fallback is broadly handled"
- "auth continuity under handoff is proven"
- "TCP failover is covered"
- "stateful services are resilient"

The whole point of this page is to keep those ceilings closed until route- or
topology-specific evidence earns them.

## Bottom line

The repo does not need more casual green checks.
It needs drills that make one exact fragment of the dream more true while
keeping the next unclosed ceiling visible.

That is what turns evidence into trust instead of architecture theater.
