# Kubernetes Deployment - Configuration Complete

## Reality check

This document records branch-level configuration validation work.
It is useful, but it is not proof that the repo's Kubernetes path is fully
deployed, healthy, equal to the root Compose-first runtime, or proven under
failure.

Read it as a configuration-progress artifact, not as closure.

For the assimilated branch reading, see:

- [`../knowledgebase/research/garden-k3s-exploration-evidence.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/garden-k3s-exploration-evidence.md)

## What this file can and cannot establish

This file can establish:

- that the branch spent real effort fixing Garden YAML and API-version issues
- that the branch reached a configuration state it considered ready for a
  deployment attempt
- which known blockers still remained visible at that moment

This file cannot establish:

- that a full cluster deployment succeeded
- that ingress worked cleanly after deployment
- that branch validation translated into runtime parity
- that the broader anti-SPOF problem was solved by this milestone

That boundary is important because "configuration complete" can easily be read
as "architecture problem solved," which would be false here.

## Status: historical branch claim

The original file framed the result as:

> all Garden.io configurations validated and ready for Kubernetes deployment

That is still useful as a branch statement.
It should now be read more carefully:

- the branch believed the configuration layer was mostly ready
- the branch still documented at least one ingress-related blocker
- readiness to deploy is not the same thing as proof after deployment

## Historical completed tasks

### YAML syntax fixes

The branch reported:

- command-array syntax fixes around `||`
- conversion of inline arrays to multi-line format
- nested template-string fixes

### API version corrections

The branch reported:

- `apiVersion` corrections from `garden.io/v2` to `garden.io/v0`
- fixes across 59 configuration files

### Configuration validation

The branch reported:

- successful validation of all Garden.io configurations
- only deprecation warnings remaining

### Local cluster setup

The branch reported:

- local `kind` cluster creation
- `kubectl` installation
- Garden CLI installation
- configured cluster access

These details still matter because they show the branch was not imaginary.
They do not prove the repo should narrate the Kubernetes path as promoted.

## Historical known issue

The original page called out an ingress-controller admission-webhook issue.

That point matters more than it may seem.
It is one of the branch clues that "configuration complete" was still not the
same thing as "runtime done."

The branch described:

- admission webhook secret required
- kind-cluster limitation or manual setup burden
- possible need for manual secret creation or webhook disablement

This is exactly the kind of detail that prevents a clean triumphant reading.

## Historical next steps

The branch proposed:

1. resolve the ingress-controller issue
2. deploy services
3. monitor the cluster

That should now be read as evidence that the branch still had a meaningful
distance between validated config and trustworthy runtime.

## Historical command reference

### Validate configuration

```bash
garden validate --env k8s
```

### Deploy to Kubernetes

```bash
garden deploy --env k8s
```

### Check cluster status

```bash
kubectl cluster-info
kubectl get nodes
```

### View pods

```bash
kubectl get pods --all-namespaces
```

## What the parity notes should now mean

The original page closed with claims like:

- health checks are comprehensive and matching
- secrets are properly configured
- volumes are mapped correctly
- dependencies are defined

The safer reading is:

- the branch believed it had achieved strong configuration alignment
- that alignment is meaningful as preparation evidence
- it is still weaker than runtime proof, failure proof, or repo-wide
  architecture closure

## Why preserve this file

Because the repo's real story includes partial readiness claims that were
stronger than "nothing worked" but weaker than "the path is now settled."

This file is one of those middle artifacts.

The docs get worse if they delete it and pretend the branch never approached
readiness.
They also get worse if they preserve it and quietly let "configuration
complete" impersonate "problem complete."
