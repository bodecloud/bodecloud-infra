# Garden and k3s Exploration: What This Branch Actually Taught

This page exists because `garden.io/` is one of the biggest remaining places
where the repo can still sound far more settled than it really is.

If those files are read carelessly, they suggest a clean story:

- Compose was migrated to Garden
- Kubernetes or k3s HA is effectively solved
- zero-SPOF architecture is already configured and nearly done

That is not an honest reading of the current evidence.

The honest reading is harsher:

- `garden.io/` preserves a serious alternate direction
- it contains useful implementation pressure and scripts
- it also contains many self-contradictory status documents
- it does **not** prove that the Compose-first root runtime has been replaced
- it does **not** prove that Kubernetes or k3s has already closed the repo's
  real failover and statefulness questions

This page matters because this branch is one of the easiest places to confuse
strong aspiration with strong closure.

The branch contains enough serious work that it can emotionally feel like a
solved migration even when the evidence still reads more like repeated
bootstrap struggle plus partial insight.

That feeling matters because it is one of the main ways infra documentation
lies without technically inventing artifacts:

- there really are many scripts
- there really are ambitious status files
- there really are cluster-shaped assets

But the emotional effect of all that serious-looking work can exceed what the
current runtime truth actually supports.

This page should keep resisting that drift.

## What this page is and is not allowed to prove

This page is allowed to prove that the repo seriously tested Garden and
k3s-shaped answers and that those efforts produced real implementation pressure,
scripts, and lessons about scheduler-backed truth.

It is not allowed to prove that Kubernetes-shaped tooling already replaced the
Compose-first runtime or that this branch already solved the repo's real
failover, ingress, and statefulness questions.

This page also has to resist a very common overread:

- many serious branch artifacts exist
- several docs use strong HA language
- therefore the branch must have been close enough that the rest is mostly
  cleanup

The current evidence does not support that jump.

## What still does not count as Garden or k3s closure

The following still do not count as serious completion in this branch:

- strong HA wording inside branch-local status files
- migration assets existing under `garden.io/`
- anti-affinity or Longhorn templates existing
- a cluster bootstrap working once in a partial environment
- zero-SPOF claims that still sit beside repair or join notes
- the branch feeling more modern than Compose

Those all show real investment.
They do not yet prove that this branch earned replacement status over the
Compose-first runtime.

## Quick claim router

Use this page for claims like:

- the repo explored scheduler-backed and k3s-shaped answers seriously enough to
  learn from them
- the branch shows how cluster tooling can reduce some missing truths while
  introducing new bootstrap, identity, and repair burdens
- the repo did not reject Kubernetes-shaped options out of ignorance

Do not use this page for claims like:

- the Compose-first root runtime was already superseded by Garden or k3s
- zero-SPOF architecture was already operationally achieved in this branch
- cluster membership and stateful HA were already sufficiently settled that the
  rest was just cleanup

## Strongest honest current answer

The strongest honest current answer is that this branch is valuable evidence,
not closure.

It proves a serious attempt to buy scheduler-owned truth with a lighter
Kubernetes-shaped stack, and it also proves that cluster bootstrap, identity,
repair, and storage truth become a major burden surface of their own.

So the branch matters because it teaches what a stronger layer costs here, not
because it has already delivered the user's full dream.

## What a real Garden or k3s promotion packet would need

Before this branch could be described as an earned future rather than high-value
evidence, a promotion packet would need to show:

- stable multi-node control-plane reality, not bootstrap-in-progress language
- which hidden placement or failover burden the branch removed in practice
- whether wrong-node request meaning improved, not just deployment structure
- what storage and state authority story actually survived node loss
- what new cluster tax the repo accepted in exchange

Without that packet, the correct reading stays:

- serious exploration
- serious implementation pressure
- still not sovereign runtime closure

## Why this branch matters at all

This exploration is not random noise.

It records a real attempt to buy several things the Compose-first runtime still
does not prove:

- stronger placement and deployment machinery
- scheduler-owned service distribution
- ingress and replica semantics that are less hand-wired
- a path away from manually reconstructed cluster truth
- a possible answer to "wrong-node entry" without inventing an entirely custom
  Compose control plane

That makes `garden.io/` important research-pressure and implementation-pressure
evidence even when it is not the active priority runtime.

It also makes the branch unusually valuable as negative evidence.

It shows what happened when the repo tried to buy a stronger truth layer by
moving upward into scheduler-backed machinery:
some real pain moved, some real capability improved, and an entirely new class
of bootstrap and control-plane truth immediately appeared.

## What the strongest Garden and k3s files claim

Representative files:

- [`../../garden.io/README-Garden-Migration.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/garden.io/README-Garden-Migration.md)
- [`../../garden.io/k8s-ha-config/COMPLETE-HA-SUMMARY.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/garden.io/k8s-ha-config/COMPLETE-HA-SUMMARY.md)
- [`../../garden.io/k8s-ha-config/ZERO-SPOF-IMPLEMENTATION.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/garden.io/k8s-ha-config/ZERO-SPOF-IMPLEMENTATION.md)

Those files contain very strong language such as:

- "migrated from Docker Compose to Garden"
- "all HA infrastructure configured and ready"
- "zero SPOF"
- "5-node Kubernetes cluster setup"
- "3-node control plane with embedded etcd"

If taken literally, that would mean the repo already crossed out one of its
biggest open questions.

The branch's own internal contradictions show that this conclusion is too
strong.

Those contradictions should not be treated like embarrassment or noise.
They are part of the evidence.

They show the exact boundary between:

- "this path may eventually earn itself"
- and "this path already solved the wound"

## What the same branch also admits

The strongest contradictory evidence comes from the same folder:

- additional nodes "joining in progress"
- single-node or in-progress control plane notes
- etcd IP mismatch repair steps
- node registration failures
- Tailscale dependency and bootstrap friction
- multiple documents whose "ready" or "complete" status still ends with
  remaining join, install, or validation steps

That means the branch does not read like shipped HA.
It reads like:

- partial cluster bootstrap
- serious script accumulation
- repeated attempts to normalize k3s and HA setup
- changing understanding of what was broken versus what was actually solved

## What this exploration really proves

### 1. The repo seriously tested Kubernetes-shaped answers

This is not hypothetical comparison prose.

The worktree contains:

- Garden project structure
- service-level `.garden.yml` files
- k3s and HA helper scripts
- ingress-related documents
- cluster bootstrap and verification scripts
- Longhorn and anti-affinity templates

That proves a real attempt existed to move toward a scheduler-owned model.

That matters because it keeps the rest of the docs honest.

The repo is not rejecting Kubernetes-shaped answers out of ignorance.
It pushed on them hard enough to discover that they import their own truth
surface, repair burden, and operational taxes.

### 2. k3s was attractive because it looked like the smallest Kubernetes tax

The branch pressure is not "Kubernetes because enterprise."

It is much more repo-shaped than that:

- k3s looked lighter than full Kubernetes
- Garden looked like a way to improve dependency and deployment structure
- the branch tried to preserve local development ergonomics while buying
  stronger cluster semantics

That is important because it matches the user's real frustration with the
ecosystem:

> if a stronger layer is required, can it be the smallest one that earns its
> keep?

That question is effectively the same one running through the rest of the
knowledgebase.

Garden and k3s matter here not because they won, but because they are one of
the strongest concrete attempts to answer that question without immediately
jumping to the heaviest possible worldview.

### 3. The branch discovered that cluster bootstrap truth is its own problem

A lot of the Garden and k3s material quietly proves something the main
knowledgebase must not forget:

cluster tooling does not eliminate bootstrap pain.
It relocates it.

The branch shows recurring pain around:

- node registration
- IP identity changes
- Tailscale integration
- etcd state repair
- multi-node control plane join correctness
- storage and replication setup

That matters because "just use k3s" is only honest if those taxes are worth the
new capabilities.

This is one of the clearest counterexamples to shallow advice in the whole
repo.

Scheduler-backed tooling can absolutely reduce some missing truths.
It can also create new ones around bootstrap, identity, repair, and cluster
health that must be carried just as honestly.

That is especially important for this repo because the user is not only trying
to buy more capable software.
They are trying to stop the feeling that there are no genuine options between
"static Docker forever" and "accept a giant new ritual stack."

This branch does not prove that k3s solved that complaint.
It proves the repo tested one candidate middle path and discovered that even
the lighter Kubernetes-shaped answer still arrives carrying a substantial new
truth surface of its own.

### 4. This branch still does not prove the user's full dream

Even if a five-node k3s cluster had been fully healthy, that would still not
finish the repo's actual question automatically.

The hard questions would still include:

- what is the trusted request path under wrong-node entry
- how middleware and auth semantics survive node-to-node handoff
- whether ingress survives local backend disappearance the way the user wants
- how stateful services replicate and fail over honestly
- whether the operator surface is still readable enough to be worth it

That means Kubernetes-shaped tooling can help, but it does not by itself
eliminate fake-HA risk.

The user is not asking for a branch that sounds more complete.
The user is asking for a path where the runtime really becomes less dependent on
private operator stitching.

A cluster only helps if it pays that tax down instead of relocating it.

## Why these files felt dangerous in the old docs

The old doc style was especially vulnerable here because it tended to flatten
these three statements into one:

- "a Kubernetes branch exists"
- "many HA configs were created"
- "the repo solved HA by moving to Garden or k3s"

Only the first two are supported.
The third is not.

That flattening is exactly the kind of document-level lie this rewrite is
trying to remove everywhere.

That distinction is exactly why this page exists.

Another dangerous flattening to keep rejecting here is:

- "the branch contains strong HA language"
- therefore "the branch must have been close enough that the rest is just
  cleanup"

No.

In this repo, cleanup language is one of the easiest ways to smuggle enormous
remaining burden out of view.
If node identity, control-plane membership, storage truth, or ingress behavior
are still unresolved, then the branch is not "basically done."
It is still inside the costly part of the problem.

## The honest status of the Garden and k3s branch

As of this documentation pass, the best evidence-aligned summary is:

- this branch is **important exploration**
- it demonstrates real appetite for scheduler-backed deployment
- it contains useful assets, scripts, and implementation pressure
- it does **not** supersede the root Compose-first runtime
- it does **not** prove shipped zero-SPOF behavior
- it does **not** close the stateful correctness problem
- it should be read as an alternate direction under investigation, not as the
  current victorious architecture

## What operators should actually take from it

Take this branch seriously for:

- what Compose pain classes pushed the repo toward Kubernetes-shaped tooling
- what bootstrap friction appeared immediately once that move was attempted
- what HA language needs to be demoted from "done" to "attempted"
- what partial assets might still be reusable later

Also take it seriously as a warning about documentation tone.

If another page cites this branch mainly to make the repo sound richer,
broader, or more inevitable, that page is using the branch dishonestly.
The most faithful use of this branch is narrower:

- it proves real experimentation happened
- it preserves reusable cluster assets
- it records the taxes that appeared immediately
- it refuses to let "cluster exists on disk" masquerade as "the dream is now
  operational"

Take it seriously, in other words, as evidence of what this promotion path
really costs and what it really promises, not as proof that the promotion has
already paid off.

Do **not** take it as proof that:

- the main runtime has been migrated
- the cluster path already works end to end
- zero-SPOF has been demonstrated
- Kubernetes or k3s has already earned whole-stack promotion

## Bottom line

`garden.io/` is not irrelevant.
It is also not closure.

It is evidence that the repo pushed hard on the Kubernetes-shaped middle-ground
question and learned something useful:

> a scheduler may solve some real problems here, but it also imports a new
> truth surface that must be bootstrapped, repaired, and operated honestly
> before it can claim to have simplified anything.

That is the main reconstruction result from this branch.
