# Stateful HA Plan: Which Truths Get Earned First

This page is about sequencing.
It is not a product catalogue and not a permission slip to call the stack
statefully resilient early.

The question here is:

> which parts of stateful resilience can `bolabaden-infra` harden honestly
> first, and which parts must remain explicitly deferred until the surrounding
> routing, placement, and authority surfaces stop depending on private human
> truth?

That wording matters because stateful HA is where infrastructure writing most
easily becomes fraudulent without intending to.

The docs can get sharper.
The diagrams can get cleaner.
The repo can still be one bad failure away from the same old answer:

- one box still mattered
- one writer still mattered
- one recovery ritual still mattered
- one human still knew the topology better than the system

## What this page is and is not allowed to prove

This page is allowed to:

- define the order in which stateful hardening becomes honest
- separate stateful subproblems that the ecosystem keeps collapsing together
- explain which service classes deserve earlier hardening and which do not
- preserve the repo's refusal to counterfeit stateful dignity from ingress
  progress

This page is not allowed to:

- imply the current root runtime already has broad stateful HA
- treat an L4 path as proof of preserved writer authority
- suggest every stateful workload should graduate to the same cluster pattern
- let planning coherence masquerade as live capability

## Strongest honest current answer

The strongest honest current answer is:

- the repo already understands the right stateful lesson
- it has not yet earned broad stateful resilience
- the right near-term move is not "cluster everything"
- the right near-term move is to harden a small number of truths in the right
  order while refusing to lie about everything else

That is not hesitation.
That is sequencing discipline.

## The stateful trap this plan is trying to avoid

Stateful fake adulthood usually looks like:

- more replicas
- more cluster nouns
- more dashboards
- more failover diagrams

while the private truth stays basically unchanged:

> I still know which node really matters, which copy is authoritative, and how
> recovery is supposed to work.

This plan exists to stop the repo from becoming more elaborate while preserving
the same private operator burden.

## What still does not count as stateful progress here

The following still do not count as honest stateful resilience:

- a port answering after restart
- a container moving to another node
- a proxy route continuing to exist
- storage living on a path that merely sounds durable
- a replica process existing without proven authority semantics
- a cluster story that never forces client rediscovery and fencing questions

If a page becomes emotionally relieving by leaning on those signals, it is
probably softening the real problem instead of solving it.

## The repo's actual near-term stateful priorities

The repo does not read like it wants to solve every distributed storage problem
immediately.
It reads like it wants to do something much more disciplined:

1. make wrong-node and ingress behavior less humiliating
2. make placement truth less dependent on human memory
3. choose a very small number of stateful systems to harden properly
4. refuse to pretend every bind-mounted service is now anti-SPOF

That sequence matters because stateful correctness built on top of routing
ambiguity and placement folklore is just a fancier lie.

## Stage 0: Keep the accusation alive

Before any technical sequence, the repo has to preserve the accusation that
created it:

- the ecosystem keeps pretending the only serious options are manual folklore
  or a heavyweight orchestrator worldview
- "multi-node" often means little more than multiple places a request can die
- many self-hosting stacks sound flexible until the request hits the wrong node
  or the authority disk disappears

If the docs lose that accusation, the plan becomes calmer and less useful.

## Stage 1: Fix request-path truth before claiming state-path truth

The first hard prerequisite to honest stateful hardening is that the repo get
better at request-path truth:

- any healthy public node can receive the request
- that node can distinguish local service versus remote service
- remote service forwarding preserves policy, auth, and operator visibility
- wrong-node success is not just lucky local placement

Why this stage comes first:

- if requests already fail when they land on the wrong healthy node, then
  replicated state underneath does not restore dignity
- otherwise the repo risks building better writer topology under a request
  plane that still humiliates the user first

This stage still does not prove stateful HA.
It only removes one adjacent source of confusion.

## Stage 2: Make placement truth explicit instead of social

Before the repo can claim serious stateful promotion paths, it needs better
current-state truth about where things actually live now.

That means some equivalent of:

- explicit placement truth
- explicit backend truth
- explicit peer eligibility truth
- fewer cases where the operator knows more than the system can say

This is where lightweight current-state registry ideas matter.

Without this stage, stateful promotion still depends too much on private human
memory:

- where the workload was running
- which peer should receive the handoff
- which clients need to be pointed elsewhere

This stage still does not prove stateful HA either.
It only removes a large amount of invisible operator glue.

## Stage 3: Separate workload classes instead of solving "stateful" as one thing

Once request truth and placement truth improve, the next honest move is not
"turn on stateful HA for the repo."

The next move is to split workloads into service classes:

- intentionally singular for now
- candidate replicated single-writer
- quorum and election worthy
- not yet worth hardening because surrounding layers are still too fragile

That classification itself is progress because it prevents fake equivalence.

### Likely honest classification from the current repo

#### A. Intentionally singular for now, but said brutally

Examples:

- Headscale control-plane authority
- some app-specific Postgres or SQLite workloads whose client behavior and
  promotion semantics are not yet owned

Honest sentence:

> This workload is important, persistent, and still singular enough that the
> docs must not call it anti-SPOF.

That is not failure.
That is the repo refusing to counterfeit dignity.

#### B. Candidate replicated single-writer workloads

Examples:

- Redis
- MongoDB

Why these rise earlier:

- the repo already names plausible topologies for them
- their role in the stack is central enough to matter
- they are easier to reason about than "all volumes for all apps"

But even here, the page must stay harsh:

- Redis Sentinel or Mongo Replica Set vocabulary does not equal proof
- client rediscovery matters as much as server promotion
- fencing matters as much as election

#### C. Stateful graph workloads that need end-to-end recovery thinking

Examples:

- Firecrawl with Redis + Postgres + RabbitMQ

Why this comes later:

- the workload is not one datastore problem
- it is a dependency graph problem
- "the app came back" is too weak when the subsystems recover out of order or
  preserve contradictory truths

These workloads should not be the first place the repo tries to feel mature.

#### D. Volume and shared-storage claims that should stay deferred

Examples:

- broad bind-mount interchangeability
- "just make the storage shared" narratives

Why this comes later:

- shared storage can multiply complexity faster than it removes pain
- the repo already knows node-local volumes are a real limitation
- pretending they are easy to generalize early is one of the fastest ways to
  produce infrastructure theater

## Stage 4: Pick one or two stateful workloads and earn the proof packet

Only after the earlier stages should the repo harden a narrow stateful set in a
way that earns stronger language.

The proof packet for any chosen workload needs at least:

1. the authority model before failure
2. the failure deliberately introduced
3. the promotion or election mechanism
4. the evidence that the old authority stopped being authoritative
5. what the clients saw before, during, and after failover
6. what storage truth preserved the authoritative dataset
7. what contradiction still remains

Without that packet, the repo may say:

- candidate topology exists
- routing path exists
- hardening plan exists

It may not yet say:

- this workload is honestly resilient

For this repo, the packet should be recorded as:

```yaml
stateful_authority_packet:
  claim_tested: "stateful authority under failure"
  service: "<one exact service, not the whole stateful lane>"
  authority_before: "<writer/leader/source of truth before failure>"
  failure_introduced: "<exact failure performed deliberately>"
  authority_after: "<writer/leader/source of truth after failure>"
  client_observation: "<dependent client behavior, not just service logs>"
  rediscovery_mechanism: "<how clients found the new writer, or manual/none>"
  fencing_or_split_brain_guard: "<mechanism, or none>"
  storage_truth: "<replication, backup, shared storage, or singular disk>"
  operator_intervention_required: true
  result: "pass | fail | honest-singularity | inconclusive"
  what_this_proves: "<narrow earned sentence>"
  what_is_still_forbidden: "<larger stateful claim still illegal>"
```

The first useful stateful packet does not have to be a victory lap.
An `honest-singularity` packet for Headscale or a Postgres-backed app can be
valuable if it removes private folklore and says exactly which manual authority
ritual still exists.

## Stage 5: Only then talk about broader stateful dignity

The repo earns broader stateful confidence only after at least one narrow
workload demonstrates:

- preserved authority
- correct promotion
- correct client rediscovery
- fencing
- storage survival

That matters because the user does not need another document that sounds like
it has many options.
The user needs proof that at least one option stopped being fake.

## Why "cluster everything" is the wrong first move

If the repo tries to solve all of this at once:

- any-node ingress
- peer forwarding
- current-state registry
- Redis failover
- Mongo failover
- Postgres failover
- queue semantics
- shared storage

then every category becomes harder to reason about, and the docs become easier
to fake simply because the system is now too complicated to inspect honestly.

The correct response to that complexity is not to write calmer prose.
It is to reduce the number of truths being promoted at once.

## The smallest honest roadmap for stateful dignity

If this repo stays faithful to the user's actual dream, the smallest honest
stateful roadmap looks something like:

1. preserve the accusation and proof boundaries
2. improve wrong-node request preservation first
3. make placement truth explicit
4. classify stateful workloads by honesty level
5. choose one or two candidate replicated single-writer workloads
6. earn proof packets for those workloads
7. keep every other stateful workload brutally honest until it earns more

That roadmap is less dramatic than "build a cluster."
It is much closer to the kind of option set the user is actually demanding.

## Bottom line

The user's dream is not a repo that sounds more distributed.
It is a repo where a request, a writer, and a failure no longer force one human
to secretly know the answer first.

This plan should only promote steps that make that sentence less true.
Everything else is choreography.
