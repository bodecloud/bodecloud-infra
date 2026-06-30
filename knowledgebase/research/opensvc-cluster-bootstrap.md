# OpenSVC Cluster Bootstrap: Serious Side-Path, Not Settled Control Plane

This page exists to stop a very specific documentation failure:

> seeing OpenSVC scripts in the repo and quietly narrating them as if the main
> stack had already crossed into a trusted clustered runtime.

That has not been proved.

What *has* been proved is narrower and still important:

- the repo has a real OpenSVC exploration path
- that path is tied to ingress and failover experiments
- the path is trying to solve a real pain class in the repo
- the path is still only a candidate, not the live sovereign control plane

That distinction matters because OpenSVC is one of the easiest candidates to
romanticize here.

It sits exactly in the psychological space the user wants to believe exists:

- stronger than raw Compose folklore
- narrower than whole-cluster worldview capture

So the docs have to be especially careful not to upgrade "serious candidate"
into "earned answer" too early.

## What this page is and is not allowed to prove

This page is allowed to prove that OpenSVC is a serious side-path in the repo
and that it is materially tied to experiments around bootstrap, ingress, and a
stronger middle coordination layer.

It is not allowed to prove that OpenSVC already governs the main root runtime
or that it has already become the repo's sovereign placement and failover
authority.

Because OpenSVC sits so close to the user's hoped-for middle layer, this page
must be more suspicious than usual about upgrading serious experimentation into
earned runtime trust.

## Quick claim router

Use this page for claims like:

- the repo seriously explored OpenSVC as a narrower alternative to jumping
  straight to Kubernetes
- OpenSVC membership is already tied to the HTTP ingress generation experiment
- the side-path is trying to reduce memory-driven coordination burden, not just
  add one more helper script family

Do not use this page for claims like:

- OpenSVC already owns current placement truth for the root runtime
- all routing generation is already driven uniformly by OpenSVC
- the repo's anti-SPOF story is already resting on a proven OpenSVC cluster

## Strongest honest current answer

The strongest honest current answer is that OpenSVC is one of the clearest
serious candidates for the missing middle layer.

It has real install, bootstrap, deploy, and sync artifacts, and the HTTP path
already expresses local-first then peer-fallback thinking through OpenSVC node
membership.

But it still remains a side-path whose scripts reveal ambition and useful
implementation pressure more strongly than they prove sovereign runtime
ownership.

## What problem OpenSVC is being asked to solve

The repo is not asking OpenSVC to replace Docker.
It is asking whether OpenSVC could help with the missing middle layer between:

- raw Compose on several nodes
- and a full Kubernetes-like worldview

The actual pressure looks like this:

- node membership should become more explicit
- failover should stop depending so heavily on operator memory
- ingress generation should consume something stronger than folklore
- the stack should gain more infra-grade continuity without abandoning its
  Compose-first identity immediately

That is why OpenSVC keeps resurfacing.

It is not resurfacing because the repo wants more names.
It is resurfacing because the repo keeps rediscovering the same wound:
requests, node continuity, and generated routing need a stronger owner than
private operator memory, but maybe not yet a universal scheduler.

## What exists in the worktree today

Concrete artifacts:

- [`scripts/install_opensvc.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/install_opensvc.sh)
- [`scripts/bootstrap_opensvc.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/bootstrap_opensvc.sh)
- [`scripts/deploy_opensvc.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/deploy_opensvc.sh)
- [`scripts/osvc_ingress_sync.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_ingress_sync.sh)
- [`scripts/osvc_ingress_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_ingress_sync.py)
- [`scripts/osvc_l4_sync.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_l4_sync.sh)
- [`scripts/osvc_l4_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_l4_sync.py)
- [`compose/docker-compose.l4-ingress.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.l4-ingress.yml)

That is enough to say OpenSVC is a serious side-path.
It is not enough to say:

- OpenSVC already owns current placement truth
- OpenSVC already governs the main root runtime
- the repo's anti-SPOF story is already resting on a proven OpenSVC cluster

## What the scripts actually show

## 1. The bootstrap path is half cluster experiment, half Docker environment prep

The current
[`scripts/bootstrap_opensvc.sh`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/bootstrap_opensvc.sh)
does not read like a minimal "join the cluster" script.
It does all of this:

- loads `.env`
- sets repo-local defaults like `DOMAIN`, `STACK_NAME`, `CONFIG_PATH`, and
  `OPENSVC_CONFIG_DIR`
- checks Docker availability
- creates expected Docker networks such as:
  - `warp-nat-net`
  - `${STACK_NAME}_default`
  - `${STACK_NAME}_publicnet`
  - `${STACK_NAME}_backend`
  - `${STACK_NAME}_nginx_net`
- verifies or installs OpenSVC
- warns if the node is not initialized for `om node ls`

That tells us something important:

the OpenSVC path here is not just "cluster membership."
It is trying to become a broader environment bootstrap and coordination layer
around the existing Docker world.

That is more serious than a toy experiment.
It is also more dangerous to overstate, because it starts to look like a shadow
control plane before it has earned that label.

That is one of the most important observations in this page.

The repo is not only evaluating whether OpenSVC can help.
It is also exposing the threshold where a "narrow helper" starts becoming a
real control surface that ought to be named honestly.

## 2. HTTP route generation currently does use OpenSVC membership

The HTTP side, in
[`scripts/osvc_ingress_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_ingress_sync.py),
currently:

- reads node names from `om node ls --format json`
- reads live Docker containers and Traefik HTTP labels
- writes `${CONFIG_PATH}/traefik/dynamic/failover-fallbacks.yaml`
- prefers local backend first
- then emits peer-node hostnames as fallbacks

This is one of the clearest artifacts in the repo expressing:

- any-node entry
- local-first service behavior
- peer-aware fallback

So yes, OpenSVC is materially tied to the HTTP ingress experiment.

That is real progress, not just branding.
But it is still progress inside an experiment, not closure over the main
runtime.

## 3. The L4 path is related, but not identical

One of the easiest ways for the docs to get sloppy is to blur the TCP side into
the same truth story as the HTTP side.

The current
[`scripts/osvc_l4_sync.py`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/scripts/osvc_l4_sync.py)
does **not** currently get node truth from OpenSVC membership.
It gets node IPs from:

- `tailscale status --json`

Then it:

- discovers opt-in services by Docker labels like `osvc.l4.enable=true`
- generates per-port HAProxy frontends and backends
- writes `${CONFIG_PATH}/haproxy/haproxy.cfg`

That distinction matters.
The L4 path is still aligned with the broader OpenSVC/failover direction, but
the current script is actually anchored to Tailscale mesh truth, not directly
to `om node ls`.

So future docs need to stop saying "OpenSVC drives all routing generation"
unless the scripts truly say that.

That kind of precision is exactly what the older docs kept losing.
This rewrite is trying to preserve it everywhere.

## 4. Docker remains central

Nothing in this side-path suggests a clean break from Docker.
The actual pattern is additive:

- Docker still runs workloads
- Compose still defines the main workload surface
- OpenSVC is being evaluated as a stronger coordination layer around those
  facts

That alignment is why this path matters emotionally as well as technically.
It is one of the few candidate layers that tries to respect the user's desire
for stronger behavior without immediately demanding a total platform religion
change.

That is why this side-path should not be dismissed just because it is not the
main runtime.
It is one of the clearest artifacts of the repo trying to find a genuinely
smaller but still truth-owning layer.

## What this side-path really proves

### It proves the repo is serious about a stronger middle layer

This is not just product comparison prose.
The worktree contains:

- install path
- bootstrap path
- deploy path
- HTTP generation path
- TCP generation path

That is enough to prove a serious attempt to make node-aware failover and
generated ingress more infra-like.

### It proves the repo is searching for cluster truth without jumping straight to Kubernetes

That is one of the most central tensions in the whole project:

- Compose alone is not enough
- but "just use Kubernetes" is not accepted as closure

OpenSVC is one of the clearest artifacts of that middle-layer search.

### It proves the control-layer problem is broader than simple routing

Because `bootstrap_opensvc.sh` also prepares Docker networks and wider runtime
assumptions, this path is already drifting toward:

- environment convergence
- bootstrap consistency
- network shape normalization

That means OpenSVC is not only about failover.
It is being explored as part of a broader "stop making the operator reconstruct
everything by hand" effort.

That sentence is one of the most useful compressions of the branch.

## What it still does not prove

It does **not** prove that:

- a production OpenSVC cluster is already deployed and trusted for the main
  stack
- OpenSVC is already the final placement authority
- the HTTP generator is battle-tested under repeated real node loss
- the TCP path is already statefully correct just because HAProxy can present a
  stable port
- the root Compose-first runtime now depends on OpenSVC strongly enough to call
  it the live control plane

Those are the exact overclaims this page exists to prevent.

## The real dependency chain behind this candidate

For this path to become truly valuable, all of these would have to hold:

1. node membership truth is accurate
2. Docker runtime discovery is accurate
3. generated routing files are correct
4. peer reachability is reliable
5. forwarded requests preserve semantics
6. stateful services behind L4 frontends have real replicated truth beneath the
   frontend

If any one of those breaks, the repo may still have a useful helper layer, but
it does not yet have the fully trusted middle layer the user is actually asking
for.

That is why this page has to sound a little uncomfortable.
The user is not looking for another path that feels promising while the real
truth chain remains implicit.

## Bottom line

The OpenSVC bootstrap path proves that `bolabaden-infra` is seriously probing a
stronger coordination layer that sits between raw multi-node Compose and
heavier orchestration. That is real progress. It also proves that the repo is
willing to let node truth and generated routing become more automatic without
abandoning Docker as the main workload substrate. What it does not prove is
that OpenSVC has already earned the right to be narrated as the present control
plane. Right now it is a serious candidate, not closure.
