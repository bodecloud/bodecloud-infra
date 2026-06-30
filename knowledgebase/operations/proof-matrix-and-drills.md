# Proof Matrix and Drill Catalog

This page is the claim firewall for `bolabaden-infra`.

It exists for the question the repo is most likely to lie about by accident:

> what exact thing are we claiming to have solved, what exact drill would prove
> it, what exact proof class have we actually reached, and what important part
> of the dream would still remain unproven even after that drill passes?

That is the real proof problem here.

The user is not asking for a generic checklist.
They are asking whether several ordinary Docker nodes can stop behaving like a
set of separate local stacks with nice branding and start behaving like one
request-preserving, operator-legible personal cloud without immediately forcing
surrender to Swarm, Kubernetes, or some other worldview-heavy control plane.

This page exists so green-looking evidence does not silently upgrade itself into
broader closure than it actually earned.

## What this page is and is not allowed to prove

This page is authoritative about:

- what proof classes exist in this repo
- what each class can honestly support
- which stronger ceilings must remain closed after a drill passes
- which exact drills would materially move the dream forward

This page is not authoritative about:

- whether a specific service has already passed a given drill
- whether one passing drill upgrades the whole stack
- whether the overall dream is satisfied just because a proof class exists

## Strongest honest current answer

The repo needs proof discipline more than it needs more confident language. A
passed drill only matters if it makes one specific fragment of the dream more
true and keeps the next unclosed ceiling visible. The danger is not missing
tests. The danger is serious-looking evidence being quietly promoted into
broader architecture closure than it actually earned.

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
- no relevant stress or failure has been exercised

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

- a known-good route was exercised
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

- ownership, replication, promotion, reconnect behavior, and storage truth were
  defined for one stateful class
- a real failure drill exercised that exact topology

Allowed claim:

- this exact stateful topology has passed this exact class of drill

Forbidden upgrade:

- the platform is now generically HA

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

Each row below should be read against that standard.

## Current matrix for the priority implementation

| Dream fragment | Current proof class | Strongest current anchors | Exact drill needed next | What still remains unproven even after that drill passes |
| --- | --- | --- | --- | --- |
| Any-node public entry is a real target | `Config present` | `.github/copilot-instructions.md`, `README.md`, `cloudflare-ddns`, master plan | Prove more than one public node can receive the first hop | That the right service is preserved after wrong-node entry |
| Local-first service is real | `Intent only` | Intent surfaces plus current Traefik-centered local routing model | Pick one stateless HTTP route and prove local serve from the honest hosting node with logs or backend identity | That wrong-node forwarding works when locality is absent |
| Wrong-node stateless HTTP requests succeed | `Intent only` | Routing philosophy plus serious edge runtime | Intentionally land a request for one stateless HTTP route on the wrong node and prove receiving-node identity, backend-node identity, and success | That the route also survives when fallback is required because the preferred local backend disappeared |
| Fallback route survives backend loss | `Intent only` | Master plan explicitly records route-persistence risk around `docker-gen-failover` | Start from a known-good route, remove the preferred local backend, and observe whether the fallback route remains present | That auth and middleware semantics stayed identical after handoff |
| Middleware and auth continuity survive peer handoff | `Config present` | TinyAuth, Nginx auth extensions, CrowdSec, Traefik middleware surfaces | Compare one protected route locally versus through intentional peer handoff and verify parity | That all protected routes now share the same continuity guarantees |
| Placement truth is live and shared | `Intent only` | Repeated `services.yaml` intent in architecture and operations pages | Introduce or expose one real placement-truth surface and prove routing or eligibility logic consumes it | That convergence, drift detection, and restart semantics are also trustworthy |
| Peer eligibility truth is live | `Intent only` | Headscale runtime plus peer-broadcast and sync direction in the master plan | Prove a receiving node chooses a peer from current tracked truth rather than folklore | That the same truth remains valid under revision drift and secret drift |
| TCP forwarding works for one named service | `Config present` | Current TCP routing such as MongoDB in root Compose | Exercise a specific TCP path end to end and prove transport plus backend identity | That the corresponding stateful semantics are safe |
| Headscale control-plane failover is real | `Intent only` | Live Headscale runtime plus leader-election and replication plan | Define exact topology, then run an owner-loss drill with continuity checks | That broader mesh or application failover is solved |
| One stateful service is honestly resilient | `Intent only` | Service-specific plan or runtime plus explicit topology | Define ownership, replicas, promotion, reconnect, and storage truth; then exercise failure | That other stateful classes are similarly mature |

## The exact drill levels

These levels are intentionally strict. They are meant to stop partial evidence
from sounding final.

### Drill level 1: authoring coherence

Questions answered:

- does the merged graph resolve?
- does the authored stack still render coherently?

Typical evidence:

```bash
docker compose config --quiet
python3 -m mkdocs build -f mkdocs.yml --strict
```

This is necessary.
It is almost never sufficient for the user's real benchmark.

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

### Drill level 3: wrong-node proof

Questions answered:

- can a request intentionally land on the wrong healthy node and still reach
  the correct peer?

Minimum evidence:

- receiving-node identity
- proof that the receiving node did not host the target locally
- backend-node identity
- user-visible success

Allowed claim after passing:

- wrong-node success exists for this exact route and service class

Still forbidden:

- the whole stack now preserves wrong-node requests

### Drill level 4: backend-loss fallback proof

Questions answered:

- does the route required for recovery remain present once the preferred local
  backend actually disappears?

Minimum evidence:

- known-good pre-failure route
- intentional backend stop or failure
- observation that the route remained present
- user-visible success through fallback

Allowed claim after passing:

- this exact route survived backend loss

Still forbidden:

- fallback is broadly trustworthy everywhere

### Drill level 5: semantic continuity proof

Questions answered:

- did the service still behave like the same protected route after handoff?

Minimum evidence:

- local and peer-forwarded comparison
- auth parity
- middleware parity
- visible policy parity

Allowed claim after passing:

- this route kept the same user-visible contract under handoff

Still forbidden:

- all protected routes now share that guarantee

### Drill level 6: stateful topology proof

Questions answered:

- does one named stateful topology keep its promised ownership and recovery
  behavior under failure?

Minimum evidence:

- explicit write owner
- explicit replica or standby behavior
- explicit promotion or recovery sequence
- explicit storage truth
- client reconnect or failover behavior under drill

Allowed claim after passing:

- this exact topology is proven for this exact failure

Still forbidden:

- the platform is now generically HA

## Candidate drills that would actually move the repo forward

The following drills matter more than generic test count because they attack the
user's real frustration directly.

### Candidate drill A: one stateless HTTP local-first proof

Good candidate classes:

- docs
- a simple dashboard
- a utility frontend
- a route whose backend identity is easy to observe

Why it matters:

- establishes the honest local baseline before wrong-node claims begin

### Candidate drill B: one stateless HTTP wrong-node proof

Goal:

- intentionally send the route through a healthy node that does not host the
  target locally

Why it matters:

- this is the first place where "multi-node" stops being emotional theater and
  starts becoming a real property

### Candidate drill C: backend-loss proof for the same route

Goal:

- stop the preferred local backend and see whether the route still survives

Why it matters:

- directly tests whether the route-generation and fallback story collapses at
  the exact moment it is needed

### Candidate drill D: protected-route parity proof

Goal:

- compare local versus peer-forwarded behavior for a protected route

Why it matters:

- proves whether request meaning survives, not just transport

### Candidate drill E: one stateful topology honesty proof

Good candidates:

- Headscale
- MongoDB
- Redis

Why it matters:

- prevents the repo from accidentally proving only HTTP optimism and narrating
  that as general resilience

## What not to do with passing results

Do not let a passing drill become:

- stack-wide maturity
- generic HA
- proof that an orchestrator decision has been settled
- proof that stateful services inherited HTTP success

The discipline rule is simple:

> always state the next unclosed ceiling

That is how the repo avoids becoming another system that sounds comprehensive a
week before it is trustworthy.
