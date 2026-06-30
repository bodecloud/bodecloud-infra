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
