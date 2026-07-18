# Kubernetes Deployment Status

## Reality check

This file is a historical branch-status note, not authoritative current repo
truth.

The priority implementation for `bolabaden-infra` is still the root
Compose-first runtime, not a completed Kubernetes migration.

This page is worth keeping because it records one of the branch's strongest
"deployment succeeded" narratives.
It is dangerous when read lazily because the tone can outrun what the broader
repo evidence actually proves.

For the assimilated reading, see:

- [`../knowledgebase/research/garden-k3s-exploration-evidence.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/garden-k3s-exploration-evidence.md)

## What this file can and cannot establish

This file can establish:

- what the branch considered a successful local Kubernetes deployment pass
- which tools and commands were reportedly used in that pass
- which parity claims the branch was trying to make relative to Compose

This file cannot establish:

- that the branch achieved stable repo-wide Kubernetes parity
- that the root runtime was replaced
- that HA or anti-SPOF goals were honestly solved
- that the cluster path survived adversarial wrong-node or backend-loss tests

That last distinction matters because a deployment-status file can easily make
progress sound more complete than the repo's real burden shift.

## Historical branch summary

**Date placeholder preserved from original file:** `$(date)`

### Status: historical branch claim, not authoritative current repo status

The original version of this file narrated a successful Kubernetes deployment
with close Compose parity.
The current worktree supports preserving that claim as branch history, not as
the final repo verdict.

### Branch-claimed prerequisites

The branch claimed these prerequisites were complete:

- Docker Compose deployment health on the source side
- Garden.io configuration parity work
- local Kubernetes cluster setup

Those claims are relevant because they show what the branch considered "ready
enough to attempt promotion."
They are not enough to prove that the promotion actually earned itself.

### Branch-claimed deployment process

The branch reported:

1. installed `kubectl`
2. installed `kind`
3. created a local Kubernetes cluster
4. installed Garden CLI
5. validated Garden.io configurations
6. deployed services to Kubernetes

This is useful as a process snapshot.
It should not be mistaken for current verified reproducibility.

### Branch-claimed service deployment

The branch described these service groups as deployed:

- core infrastructure
- reverse proxy services
- application services
- LLM services
- Stremio services
- metrics services
- WARP services

Again, that tells the reader what the branch was trying to move, not what the
repo has now authoritatively proven under the active Compose-first priority.

## Verification language boundary

The original page used direct success language such as:

- all configurations validated
- services deployed in dependency order
- Kubernetes resources created
- pods running and healthy

Those statements should now be read as branch-local success assertions.

They do **not** by themselves prove:

- the cluster stayed healthy after deployment
- ingress and service identity matched the user's real anti-SPOF benchmark
- wrong-node request preservation was solved
- stateful systems gained honest HA semantics

## Historical commands preserved from the branch

### Deploy to Kubernetes

```bash
cd garden.io
export KUBECONFIG=/tmp/kubeconfig
export PATH="/tmp/garden-install:$PATH"
garden deploy --env k8s
```

### Check status

```bash
export KUBECONFIG=/tmp/kubeconfig
kubectl get pods --all-namespaces
kubectl get deployments --all-namespaces
kubectl get services --all-namespaces
```

### Historical next steps

The branch recommended:

1. monitor pod health
2. check service status
3. verify ingress
4. review logs

That is still useful as a record of what the branch operators thought remained
important after the deployment claim.

## What the parity language should now mean

The original page implied:

- configurations were identical between Compose and Kubernetes
- health checks matched
- secrets mapped correctly
- volumes mapped to persistent volumes

The safer modern reading is:

- the branch *aimed* for close parity
- the branch believed it had reached a strong enough configuration state to
  claim deployment progress
- the repo as a whole still contains contradictory evidence about how far this
  path actually went

That is a much more truthful use of this file.

## Why this page still matters

Because it captures one side of a real contradiction in the repo.

Some branch pages say:

- Kubernetes direction incomplete
- bootstrap still in progress
- HA still not fully demonstrated

This page says, in effect:

- deployment worked well enough locally that the branch treated it as a
  success

The knowledgebase needs both sides of that contradiction.
What it must not do is let the strongest success tone quietly overwrite the
rest of the evidence.
