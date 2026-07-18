# Decision Paths and Promotion Rules

This page is the promotion gate for architecture choices in `bolabaden-infra`.

It does not ask "which tool looks more serious?"
It asks:

> which next move actually removes one private bad-day sentence the operator
> currently has to finish for the platform?

If a proposed promotion cannot answer that, it is still prestige-shaped rather
than burden-shaped.

## What this page is and is not allowed to prove

This page is authoritative about:

- how to evaluate the next layer choice
- what counts as too weak a reason to promote
- what a stronger layer would have to genuinely own

This page is not authoritative about:

- the repo's final architecture already being chosen
- a specific promotion already being complete
- roadmap clarity equaling shipped capability

## Strongest honest current answer

The current evidence still supports this default path:

1. stay Compose-first for authoring
2. stop lying about what Compose alone does not own
3. add the smallest shared-truth layer that makes placement and peer choice
   less social
4. prove one wrong-node stateless HTTP path
5. prove backend-loss survival for that same route
6. only then promote specific pain domains that still remain fake or too
   operator-dependent

That is not timidity.
It is the smallest path that still respects the user's refusal pattern.

## Promotion rule

Do not promote because a bigger tool exists.
Promote because the current layer has become the hidden tax.

The harsh version of the rule is:

> after this promotion, what exact sentence does the operator no longer need to
> finish privately?

If that sentence cannot be written concretely, promotion is still premature.

## Weak reasons to promote

These are not enough:

- "the ecosystem standard is Kubernetes now"
- "the helper layer is awkward to summarize"
- "there are enough moving parts that a scheduler feels inevitable"
- "we can draw a cleaner diagram with a stronger control plane"
- "the bigger system looks more dynamic"

All of those can be true while the operator still privately owns:

- placement truth
- peer validity
- fallback durability
- protected-route meaning
- stateful authority

## Promotion packet standard

A real promotion packet should contain:

- the named hidden burden being reduced
- the truth source being externalized or promoted
- the artifact that carries that truth now
- the drill that proves behavior changed under stress
- the operator-visible inspection path for that decision
- the sentence that still remains outside the promoted layer

Without that packet, "promotion" is still mostly rhetoric.

The packet also has to identify the authority it is replacing.

For this repo, that usually means naming one of these current private
authorities:

- operator memory of current placement
- operator judgment of peer safety
- operator interpretation of generated fallback files
- operator knowledge of middleware and auth continuity
- operator caveats around stateful authority

If a candidate promotion cannot name the private authority it is replacing, it
is not yet a promotion.
It is only a new component.

Use a structured packet when a tool, helper, registry, generated file, or
orchestrator starts being described as an answer rather than research:

```yaml
promotion_packet:
  candidate_layer: "<tool, helper, registry, generated artifact, or controller>"
  promoted_from: "research | helper | docs-only | local-runtime | manual-ops"
  hidden_burden_removed: "<one private sentence the operator no longer owns>"
  replaced_private_authority: "<operator memory, judgment, caveat, or manual reconstruction>"
  new_truth_authority:
    artifact: "<file, API, runtime state, controller, log, or generated config>"
    freshness_or_convergence: "<how stale or drifted truth is detected>"
  runtime_consumer: "<route path, proxy, sync loop, agent, controller, or CLI>"
  drill_required: "<wrong-node, backend-loss, policy-parity, stateful-authority, or other>"
  drill_packet: "<route_packet, placement_decision_packet, backend_loss_packet, stateful_authority_packet>"
  operator_inspection_path: "<command, log, page, file, or dashboard proving why>"
  legal_new_claim: "<one narrow claim now allowed>"
  still_forbidden:
    - "<stronger claim still illegal>"
```

If the packet cannot name `runtime_consumer`, it is still only documentation or
inventory.
If it cannot name `drill_packet`, it is still not stress-proven.
If it cannot name `hidden_burden_removed`, it is probably solving a neighboring
problem instead of the user's actual wound.

## Minimum v1 promotion boundary

The first promotion does not need to solve the whole repo.
It does need to make one complete burden transfer visible.

A contract-faithful v1 promotion should include:

1. one explicit placement or peer-eligibility truth source
2. one consumer of that truth in the request path or route-generation path
3. one stateless HTTP route that exercises local versus remote behavior
4. one backend-loss drill for that same route
5. one operator inspection path that explains why the receiving node made the
   decision it made

That is the smallest useful unit because it proves a real transfer of
responsibility.

A smaller packet can still be good implementation progress, but it should not
be allowed to upgrade the docs into saying the middle layer has arrived.

The minimum v1 packet is intentionally stateless and HTTP-shaped.
That is not because TCP or stateful workloads are unimportant.
It is because the first promotion has to prove the system can absorb one
operator burden without immediately borrowing confidence from harder lanes.

## Rejection rule for fake options

An option should be rejected or kept as research-only when it mostly answers one
of these easier questions:

- can we make the diagram cleaner?
- can we make deployment feel more modern?
- can we reduce the number of manual commands?
- can we make the proxy config more dynamic?
- can we get a dashboard that looks more cluster-like?

Those are useful only if they serve the harder question:

> which private bad-day sentence stops being true?

If the answer is still "none yet," the option may remain useful research, but
it has not earned promotion.

## The main candidate paths

### Path 1: Stay Compose-first and add a lightweight shared-truth layer

Use when:

- the main missing burden is current placement and peer-eligibility truth
- the edge stack is already capable enough locally
- the repo still wants readable authoring and direct inspection

This path must eventually answer:

- where does current placement truth live?
- what consumes it?
- how does a receiving node choose local versus peer route?

Promote this path if it can delete the sentence:

> I still have to remember where it really lives

Do not promote this path merely because a `services.yaml`-shaped artifact
exists.
Promote it when some runtime path consumes that truth and changes behavior
because of it.

Minimum promotion packet for this path:

- `hidden_burden_removed`: `I still have to remember where it really lives`
- `new_truth_authority`: current-state registry or generated runtime state
- `runtime_consumer`: route-generation path, edge proxy, or routing decision
  layer
- `drill_packet`: `placement_decision_packet` plus a route packet for one
  stateless HTTP route

### Path 2: Harden helpers without promoting a full controller

Use when:

- the missing burden is route persistence or convergence drift
- the repo already has useful helpers but they are not yet trustworthy

This path must eventually answer:

- which helper owns which truth?
- what failure scene does it survive?
- how is its decision visible to the operator?

Promote this path if it can delete the sentence:

> I still have to remember whether fallback survives the real failure

Do not promote this path merely because a helper renders a file.
Promote it when the generated decision remains correct during the failure scene
it exists to handle.

Minimum promotion packet for this path:

- `hidden_burden_removed`: `I still have to remember whether fallback survives
  the real failure`
- `new_truth_authority`: generated config plus an inspectable survival record
- `runtime_consumer`: the proxy or route path that keeps using the fallback
  after preferred-backend loss
- `drill_packet`: backend-loss packet for one named route

### Path 3: Promote a stronger ingress or service-coordination layer

Use when:

- lightweight truth surfaces exist
- wrong-node and backend-loss drills show the next ceiling is no longer
  placement truth but execution durability

This path must eventually answer:

- what exact burden the stronger ingress layer now owns
- why a smaller helper stack could not carry it honestly first

Promote this path if it can delete the sentence:

> I still have to remember whether the same protected-route meaning survives
> handoff

Do not promote this path merely because a stronger ingress layer can route to
more places.
Promote it when it preserves service meaning, policy, and eligibility under a
wrong-node handoff.

Minimum promotion packet for this path:

- `hidden_burden_removed`: `I still have to remember whether the same protected
  route meaning survives handoff`
- `new_truth_authority`: ingress/service-coordination state visible outside
  one operator's memory
- `runtime_consumer`: the actual protected route path
- `drill_packet`: route packet plus policy-parity evidence

### Path 4: Promote a heavier orchestrator or cluster worldview

Use only when:

- the smaller shared-truth and helper layers were actually tried or honestly
  ruled out
- a heavier system can be shown to own a burden the repo still cannot carry
  without private folklore

This path must answer:

- which exact burden did the smaller layer fail to externalize?
- which exact system-owned truth does the heavier platform now provide?
- how can the operator inspect that truth rather than just trust the
  controller?

Promote this path if it can delete a sentence like:

> I still have to remember which peer is actually safe and why the platform
> chose it

Do not promote this path merely because Compose has become uncomfortable.
Promote it only when the discomfort is proven to come from a burden a heavier
system would actually own better.

Minimum promotion packet for this path:

- `hidden_burden_removed`: the exact private sentence the smaller layer failed
  to kill
- `new_truth_authority`: scheduler, service registry, controller state, or
  cluster API that owns the missing truth
- `runtime_consumer`: request path, scheduling path, state authority path, or
  operator workflow that actually consumes that truth
- `drill_packet`: the narrow packet matching the burden, not a general
  "cluster works" result

## Demotion and quarantine rules

Not every serious experiment deserves to stay in the active story.

Demote or quarantine a path when:

- it adds a new dashboard but leaves truth ownership private
- it adds a new controller but still requires the operator to explain the bad
  day afterward
- it makes one easy route cleaner while making protected routes or stateful
  lanes easier to overclaim
- it creates a second source of truth without proving which one wins under
  drift
- it requires more worldview tax than the burden it removes

Demotion does not mean the experiment was worthless.
It means the docs should stop letting it sound like an earned option.

## Promotion ledger format

Use this short ledger for any future promotion decision:

| Field | Required answer |
| --- | --- |
| Candidate layer | The concrete tool, helper, file, service, or controller being promoted |
| Promoted from | Research, helper, docs-only, local runtime, or manual ops |
| Private sentence removed | The exact operator-owned sentence that should die |
| Replaced private authority | The human memory, judgment, caveat, or reconstruction being replaced |
| New truth authority | The artifact or component that now owns that truth |
| Runtime consumer | The code, service, proxy, or route path that consumes it |
| Drill passed | The stress scene that proved behavior changed |
| Drill packet | The proof packet that records the stress result and ceiling |
| Inspection path | How an operator can see why the decision happened |
| Legal new claim | The one narrow stronger sentence now allowed |
| Still forbidden | The stronger claim that remains illegal afterward |

If any row is empty, the promotion is not yet real.

## What still does not count as a real option

An option is still too weak if it mainly changes:

- vocabulary
- controller prestige
- dashboard confidence
- deployment ergonomics

while leaving unchanged:

- where placement truth lives
- who determines peer eligibility
- whether wrong-node requests stay meaningful
- whether fallback survives backend loss
- whether stateful authority is still singular

That is the repo's core anti-fake-option rule.

## The next honest question for any proposal

Before accepting any proposal, ask:

1. what burden is still private today?
2. what artifact would externalize it?
3. what drill proves that artifact changes behavior on the bad day?
4. what stronger sentence becomes legal after that drill?
5. what sentence is still forbidden?

If the proposal cannot survive those five questions, it has not earned
promotion yet.
