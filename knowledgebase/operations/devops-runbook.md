# DevOps Runbook

This is not a generic Docker runbook.

It is the operator method for one unusually specific problem:

> how do we tell the difference between a Compose-first multi-node stack that
> actually preserves requests under pressure and one that only *looks* more
> resilient than it is?

That is the real operational question behind `bolabaden-infra`.

The user is not mainly asking for startup commands.
They are asking for a way to stop hidden operator memory from acting like the
real control plane.

This runbook exists to make that burden visible and to stop green output from
inflating into resilience theater.

It exists because operator calm is cheap in a repo like this.
A few healthy commands can make the stack feel governed.
The runbook's job is to keep asking what the operator still had to know anyway.

## What this page is and is not allowed to prove

This page is authoritative about:

- how to route questions to the right proof class
- what command classes are too weak for stronger claims
- what order to inspect the stack in
- what counts as weak, medium, and strong operational evidence in this repo
- how to keep operators from narrating comfort as closure

This page is not authoritative about:

- whether a specific ingress path is already resilient
- whether a specific stateful service is already safe
- whether the architecture dream has already been achieved

This is an operator method page, not a completion certificate.

It is also not a comfort script.
If following it mainly leaves a human feeling reassured, rather than leaving
behind inspectable evidence plus a named remaining ceiling, then it is already
drifting back toward ritual instead of operational truth.

It is also not allowed to become a page that professionalizes the hidden
operator role.

That is a subtle failure mode in a repo like this:

- the runbook gets sharper
- the questions get better
- the sequence gets more disciplined
- the operator becomes more fluent at explaining what is and is not proven

and yet the same person is still privately acting as the final bridge between
what the system exposes and what the stack actually means under stress.

That is better operational literacy.
It is not yet less hidden operator burden.

## Strongest honest current answer

The real job of this runbook is:

1. force every claim to name its proof class first
2. force every successful command to state what stronger story would still be
   a lie
3. force the operator to distinguish authored shape, local runtime health,
   route behavior, wrong-node behavior, backend-loss behavior, and stateful
   correctness
4. expose where private operator reconstruction is still doing work the system
   should eventually own itself

If the runbook does not do those things, it becomes a comfort ritual instead
of an operational tool.

That distinction matters more here than in an ordinary service repo.
The root wound is not lack of commands.
It is that commands keep succeeding one layer before the hidden operator role
actually shrinks.

So this page has to keep separating two things that are easy to confuse:

- a more competent operator
- a system that requires less operator heroism

The first can improve dramatically while the second barely moves.
This runbook has to keep making that uncomfortable on purpose.

## What still does not count as a serious runbook in this repo

This page should also say more bluntly what fake operational maturity still
looks like.

These still do not count:

- a long sequence of healthy commands with no named proof class
- a green validation pass followed by a broader resilience sentence
- a route test with no backend identity evidence
- a failover claim with no before/after semantics comparison
- a stateful reassurance story that never names write authority
- a calm operator summary that hides how much private topology memory was still
  required

The user is not mainly asking for more commands.
The user is asking for less hidden reconstruction burden.

If the runbook leaves the operator feeling informed but still privately
carrying the same sacred-node truth, then it stayed too shallow.

If it leaves the operator sounding more authoritative while still carrying the
same sacred-node truth, then it may actually be making the wound easier to hide.

That is the human test for every runbook section:

> after this step, what exact private sentence did the operator no longer need
> to finish alone?

## The dream this runbook has to protect

The repo is trying to keep the directness of Compose while pushing toward:

- any-node public entry
- local-first serving
- peer-forward fallback
- anti-SPOF pressure
- honest boundaries where the system still depends on shared operator memory

Operationally, that means the runbook cannot be satisfied by:

- "`docker compose config` passed"
- "Traefik is healthy"
- "the service answered on one node"
- "the container restarted"

Those may be necessary signals.
They are not the same as the dream the user is asking for.

Operationally, the dream is not "a nicer incident notebook."

It is:

- wrong-node entry that does not collapse into private recollection
- fallback that preserves service meaning instead of just transport
- evidence that stays inspectable after the stressful moment passes
- a platform where the operator increasingly verifies shared truth instead of
  impersonating it

That last line is the whole operational horizon.
The runbook is not here to make the operator more heroic.
It is here to make heroism less necessary and more obviously measurable when it
is still required.

That means "good runbook quality" is not enough.
The page cannot merely teach someone to narrate the system honestly.
It also has to keep pointing at the exact place where the system still refuses
to own what the operator just learned how to explain more clearly.

That is the bar this page has to protect.

## What the runbook is trying to kill

The real recurring failure mode in this repo is hidden reconstruction burden.

That burden appears whenever the operator has to privately remember or infer
what the system should be exposing explicitly.

There are four recurring versions of that burden.

### 1. Hidden topology burden

The operator should not have to privately remember:

- which node currently hosts the service
- whether the request path is local or remote
- whether the peer being targeted is merely reachable or actually the right
  backend

If the answer still lives mostly in someone's head, the operator's head is
still part of the control plane.

### 2. Hidden convergence burden

The operator should not have to guess:

- whether nodes are on the same revision
- whether secrets and env surfaces still match
- whether auth and middleware assumptions are equivalent across nodes
- whether a peer-forwarded request lands in semantically comparable runtime

Transport success without convergence truth is not a real recovery story.

### 3. Hidden claim burden

The operator should not have to keep mentally translating:

- "multiple DNS records" into "not yet preserved requests"
- "healthy proxy" into "not yet wrong-node proof"
- "container restarted" into "not yet state-safe recovery"

If the operator must keep doing that privately, the docs and tooling are still
too flattering.

### 4. Hidden proof burden

The operator should not have to reconstruct after the fact:

- which evidence class a claim belonged to
- which failure mode was actually exercised
- what remained unproven even after the command passed

That is one of the central reasons this runbook exists.

## Start by naming the real question

Before running anything, say what you are actually trying to prove.

In this repo, the serious questions are usually one of these:

- does the merged root Compose graph still resolve?
- what services are actually present in the priority runtime?
- is a named service locally healthy on this node?
- does one documented ingress path answer on a normal day?
- what happens when the request lands on the wrong node?
- what survives when the preferred local backend disappears?
- did a stateful service keep correctness, not just reachability?

If the question is not named first, the operator will almost always drift
into:

- run enough commands to feel reassured
- then narrate that reassurance as resilience

That drift is not just a habit problem.
It is one of the main ways hidden operator knowledge keeps surviving while the
repo sounds increasingly sophisticated.

That is exactly the failure pattern the user is trying to get away from.

This sequence is therefore not a performance rubric for impressive operators.
It is a measuring stick for how much private closure work still survives after
the operator has done everything correctly.

## The operational sequence

When investigating or validating anything serious, use this order:

1. name the claim
2. name the proof class required
3. inspect authored shape first
4. inspect local runtime second
5. inspect route behavior third
6. only then move into wrong-node, backend-loss, or stateful drills
7. after every success, state the next stronger sentence that is still
   forbidden

That last step is where most runbooks fail.

It is also the step most likely to keep this repo honest.
The first successful command is rarely the lie.
The lie usually begins with the next paragraph.

## Proof classes

This repo becomes much easier to operate once "is it working?" is split into
smaller claim types.

### 1. Authoring proof

Questions answered:

- does the tracked graph interpolate?
- do the docs still render?
- is the authored configuration structurally coherent?

Commands:

```bash
docker compose config --quiet
docker compose config --services
python3 -m mkdocs build -f mkdocs.yml --strict
```

What this proves:

- the authored surfaces still resolve
- the merged root graph still has a coherent shape

What this does **not** prove:

- containers started
- routes work
- peers agree
- failover exists
- state is safe

This is the minimum proof class, not the final one.

That reminder should be emotionally prominent, not procedural boilerplate.
In this repo, authoring proof is especially seductive because the authored
graph is already rich enough to look like half the platform has arrived.

### 2. Local runtime proof

Questions answered:

- what is present on this node?
- what started?
- what healthchecks are passing?

Commands:

```bash
docker compose ps
docker inspect <container> --format='{{json .State.Health}}'
docker logs --tail=100 <container>
```

What this proves:

- this node currently has a particular running shape
- the local healthcheck surface is reporting something specific

What this does **not** prove:

- wrong-node behavior
- peer-forward eligibility
- fallback persistence
- stateful correctness

This is one of the biggest temptation zones in the whole repo.
Local runtime proof is where a stack can begin to feel complete while still
leaving the real distributed question unanswered.

### 3. Route correctness proof

Questions answered:

- does a known ingress path answer on a normal day?
- which backend actually answered?
- were the expected auth and middleware surfaces involved?

Typical evidence:

- request output
- receiving-node logs
- backend identity evidence
- auth or middleware logs when relevant

What this proves:

- one documented path still works in a specific scenario

What this does **not** prove:

- wrong-node success
- backend-loss survival
- full service-class resilience

### 4. Failure-path proof

Questions answered:

- does the request survive when it lands on the wrong node?
- does fallback survive when the preferred backend disappears?
- does the post-fallback behavior remain semantically equivalent?

Typical evidence:

- intentional wrong-node request
- receiving-node logs
- peer backend logs
- before/after route identity comparison
- backend stop or simulated backend-loss drill

What this proves:

- one specific failure path was exercised and behaved in a particular way

What this does **not** prove:

- generic stack-wide resilience
- all routes share the same recovery properties

### 5. Stateful correctness proof

Questions answered:

- who owns write truth now?
- what replicates it?
- how do clients discover current authority?
- what happens after promotion or node loss?

Typical evidence:

- service-specific topology inspection
- replication status
- promotion evidence
- reconnect behavior
- storage or authority ownership evidence

What this proves:

- correctness of one stateful class in one tested topology

What this does **not** prove:

- generic "the platform is HA now"

## Claim router

Use this before reaching for commands.

| Real question | Start with | Next proof if that passes | Must not be upgraded into |
| --- | --- | --- | --- |
| Does the priority merged root runtime still resolve? | `docker compose config --quiet` | `docker compose config --services` | claims about runtime behavior or resilience |
| What is actually part of the current root runtime? | `docker compose config --services` | `docker compose ps` | claims about cross-node correctness |
| Did a local service start and stay healthy? | `docker compose ps` and `docker inspect` | route-targeted proof for public services | claims about peer-forward or failover readiness |
| Does one ingress path answer normally? | targeted request plus logs | backend identity and policy-path proof | claims about wrong-node or backend-loss survival |
| Does one wrong-node request still preserve the service? | intentional wrong-node request plus logs on both nodes | backend-loss drill on the same service | generic multi-node success |
| Does fallback survive local backend loss? | known-good path plus backend-loss drill | semantic comparison before vs after | stack-wide failover completion |
| Is one stateful class actually resilient? | topology-specific ownership and replication evidence | real failure drill for that service class | broad HA claims for all state |

## Practical validation baseline

These are the minimum repo-native checks that should stay routine:

```bash
docker compose config --quiet
docker compose config --services
python3 -m mkdocs build -f mkdocs.yml --strict
```

Interpret them narrowly:

- they prove authored shape
- they do not prove runtime behavior
- they definitely do not prove distributed correctness

## The minimum serious operator notes after each validation

After any meaningful command or drill, record four things:

1. what exact claim you were testing
2. what exact proof class the evidence reached
3. what stronger sentence is still forbidden
4. where hidden operator reconstruction still survived

If those notes are missing, the command output will almost always get
overread later.

## What a real runbook proof packet would have to contain

The repo should not leave behind isolated command output.
It should leave behind proof packets that another operator can inspect without
reconstructing the whole event socially.

A serious packet should contain:

- the exact question being tested
- the exact proof class reached
- the exact node or nodes involved
- the route or service identity that was supposed to survive
- the strongest sentence that became more honest
- the stronger sentence that is still forbidden
- the place where hidden operator memory still had to fill a gap

For stronger drills, the packet should also contain:

- before/after backend identity evidence
- logs from both the receiving node and the selected peer when peer handoff is
  involved
- the specific failure introduced on purpose
- the semantic difference, if any, before and after the failure

If the packet cannot be inspected later without a narrator, the narrator is
still part of the control plane.

## What not to do

Do not let the runbook become:

- a list of comfort commands
- a generic Docker health guide
- a proof substitute for wrong-node and backend-loss drills
- a place where stateful services inherit HTTP optimism

The whole point is to keep the operator from feeling calmer than the evidence
allows.

The runbook succeeds only when it makes overstatement harder than honesty.
