# Stateful HA and Data Ownership

For the harder evidence boundary behind this page, start with
[Stateful HA Evidence](../research/stateful-ha-evidence.md).

This page is not a database inventory.
It is the architectural reading of what the user is really asking once the repo
stops talking about ingress and starts talking about truth.

The harsh question is:

> if the current node disappears, where does authoritative truth live, who may
> still write, how do clients discover that writer, and what prevents the old
> writer from still acting alive?

If a page does not answer those questions, then it may still be useful, but it
is not yet answering the wound that created `bolabaden-infra`.

## What this page is and is not allowed to prove

This page is allowed to prove:

- the live priority runtime already depends on several state-bearing systems
- those systems still reveal strong node-local authority assumptions
- the repo already understands that ingress continuity is not state continuity
- different service classes deserve different honesty and different hardening
  paths

This page is not allowed to prove:

- that TCP exposure equals HA
- that restartability equals authority preservation
- that a second node plus a proxy equals stateful dignity
- that a future plan makes present tense storage safe

## The private sentence these pages must kill

For stateful systems, the hidden sentence is usually some version of:

> I still personally know who the real writer is, which copy is authoritative,
> and how everyone is supposed to recover if the current node disappears.

If that sentence survives, then the stack may have improved:

- reachability
- recovery speed
- operator ergonomics
- visibility

Those gains are real.
They are not yet authority transfer.

## The difference between public continuity and truth continuity

The repo's multi-node dream begins at public entry:

- any healthy public node should be able to receive the request
- a local service should be served locally
- a remote service should be forwarded to a healthy peer

That already solves a real pain.
It still does not solve the deeper stateful question.

For stateful services, "the request still got somewhere" is not the same as:

- the right writer still exists
- the writer still has the truthful data
- clients can find it without folklore
- the old writer cannot act like it still owns truth

That is why this layer needs meaner language than the routing layer.

## Strongest honest current answer

The strongest honest current answer is:

- stateful risk is already live in the root runtime
- most authority is still singular per workload
- most storage truth is still expressed as a node-local path
- the repo's planning pressure is ahead of the runtime's authority proof

That means the project is already better at diagnosing the stateful burden than
at claiming it has been solved.

## The four questions every stateful claim should answer

Every serious stateful paragraph in this repo should answer all four:

1. Where does authoritative truth live right now?
2. What happens to that truth if the current authority node disappears?
3. How do clients discover the new authority without human memory filling the
   gap?
4. What prevents two competing copies from both acting like the writer?

If those questions are not answered, then the control plane is still partly
social.

## The authority ladder

The repo should treat stateful language as a ladder, not a switch.

| Level | What has been shown | What is still forbidden |
| --- | --- | --- |
| 1. Reachable | a TCP or HTTP endpoint answers | authority, correctness, or failover |
| 2. Locally healthy | the process passes a real healthcheck on its current node | cross-node survival |
| 3. Recoverable by operator | a human can restore or restart the workload with known steps | automatic or system-owned authority transfer |
| 4. Replicated or backed up | another copy, replica, or backup exists | that the copy is current, writable, or safe to promote |
| 5. Controlled writer | the system can name the current writer or leader | that clients can find it after failure |
| 6. Rediscovered by clients | clients reconnect to the correct new authority | that split brain and stale writes are impossible |
| 7. Drilled and packeted | a failure drill records authority, promotion, rediscovery, fencing, and storage truth | broad stateful anti-SPOF for unrelated workloads |

Most current priority-runtime stateful services are still below level 5.
That is not an insult to the stack.
It is the point of the documentation.

## `stateful_authority_packet`

Any future document, runbook, or implementation note that wants stronger
stateful language should carry a packet shaped like this:

```yaml
stateful_authority_packet:
  claim_tested: "stateful authority under failure"
  service: "redis | mongodb | headscale | postgres | rabbitmq | qdrant"
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

The packet is deliberately more annoying than a status table.
It forces the repo to say whether it proved authority transfer or merely
documented a better recovery ritual.

## Service classes in the current stack

The live stack already includes several distinct state classes:

| Service class | Live examples | What matters most |
| --- | --- | --- |
| Document / key-value primary stores | MongoDB, Redis | writer identity and client rediscovery |
| Relational backing stores | `nuq-postgres`, `litellm-postgres` | promotion correctness, storage truth, reconnect behavior |
| Queue / broker state | RabbitMQ | ordering, durability, delivery semantics under failure |
| Local-file control-plane authority | Headscale with SQLite | singular control-plane truth and promotion discipline |
| Vector state | Qdrant | durable index truth, replica semantics, recovery ordering |

The point of this table is not "we run a lot of databases."
The point is that each of these classes fails differently, lies differently,
and deserves different wording.

## The stateful burden map

| Burden | Real question | False comfort signal |
| --- | --- | --- |
| Writer identity | who may accept writes right now? | the service answered |
| Authority continuity | who owns truth after failure? | the container restarted |
| Client rediscovery | how do dependents find the new writer? | DNS still landed somewhere |
| Fencing | what stops the old writer from acting alive? | the old node is probably down |
| Storage truth | did truth replicate or remain singular? | there is another node available |

This burden map is why stateful docs cannot use the same emotional shortcuts as
availability marketing.

## Honest readings of the current service classes

### MongoDB

Current honest reading:

- MongoDB is live, persistent, and important
- the current runtime expresses that persistence through a node-local data path
- the repo does not yet prove Replica Set behavior, primary election, or client
  rediscovery

That means the correct present tense is not:

> MongoDB is effectively HA because it is routable and durable

The correct present tense is:

> MongoDB is a live stateful authority surface whose durable truth is still
> mostly tied to one node at a time.

### Redis

Current honest reading:

- Redis is a real dependency, not decorative cache language
- the runtime still looks like one durable Redis writer story
- the repo does not yet prove Sentinel, promotion correctness, or master
  rediscovery by clients

That means the correct present tense is:

> Redis is singular enough that a human can still know where the real writer is
> better than the system itself can express it.

### Headscale

Current honest reading:

- Headscale's control plane is useful and already integrated
- it is currently backed by one SQLite authority path
- this is one of the clearest places where a valuable service is not yet a
  resilient authority surface

The important thing here is not to call SQLite itself embarrassing.
The important thing is to refuse to pretend that a singleton control plane has
become anti-SPOF just because the surrounding ingress and mesh story got
smarter.

### Firecrawl and related Postgres / RabbitMQ surfaces

Current honest reading:

- these subsystems already depend on several durable backends at once
- app continuity here is really a graph problem, not a single container
  problem
- "the app came back" is weak evidence if Postgres, Redis, and RabbitMQ fail
  differently and recover in the wrong order

This is one of the easiest places for infrastructure docs to accidentally tell
a half-truth:

> the service recovered

Maybe.
But the real question is whether the authoritative graph recovered correctly.

### Qdrant and vector state

Current honest reading:

- vector state is already part of the live runtime
- the runtime still expresses its truth through a local storage path
- the repo does not yet prove clustered vector authority semantics

Vector systems are especially prone to fake certainty because they often live
near AI tooling that already sounds modern and distributed.
That branding should not be allowed to do architectural work.

## Not every workload deserves the same stateful strategy

One of the most important corrections this repo needs to preserve is that
"stateful HA" is not one singular feature to turn on.

Different workloads may need different honest answers:

- some should remain intentionally singular for now, with brutal honesty
- some may deserve a replicated single-writer topology
- some may deserve quorum and election machinery
- some may not be worth hardening until ingress, placement truth, and client
  behavior become less fragile

The ecosystem often pressures operators into pretending every stateful service
has only two respectable end states:

- stay embarrassingly manual
- or graduate into full cluster theater

The whole point of this repo is that the user is sick of fake binaries like
that.

## What "better than today" is allowed to mean

For a stateful workload, it is honest to say "better than today" if the repo
has genuinely improved one of these:

- clearer placement truth
- faster recovery with preserved single-writer honesty
- better operator visibility into who currently owns authority
- better proof discipline around failure behavior

It is dishonest to let "better than today" quietly expand into:

- anti-SPOF
- resilient failover
- authority continuity

unless the proof packet actually earns those phrases.

## The minimum dignified sentence for each workload

Until stronger evidence exists, the minimum dignified sentence for a workload
should usually sound like one of these:

- this workload is stateful, important, and still singular in authority
- this workload has a candidate replicated topology, but it is not proven live
- this workload is reachable across nodes, but authority continuity is still
  unproven

Those sentences are less flattering.
They are also the kind of sentences that let the user reason honestly about the
next move.

## Bottom line

The user is not asking for friendlier wording about persistent services.
The user is asking whether the system can stop hiding truth in one node, one
disk, one writer, or one operator's head.

The current worktree does not yet prove that.

What it does prove is that the repo already knows the right question, already
contains the real stateful pain surfaces, and already has enough self-awareness
to refuse transport theater as a substitute for preserved authority.
