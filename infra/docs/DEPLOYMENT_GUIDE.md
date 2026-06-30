# Constellation Agent Deployment Guide

## Read this before treating it like production guidance

This file explains how Constellation appears intended to be deployed.

It should **not** be read as proof that:

- Constellation is the adopted runtime for the repo
- the deployment path is fully validated end to end
- the resulting cluster is automatically zero-SPOF

Use it as an implementation-oriented guide for experimentation, auditing, and
future verification work.

## What Constellation is trying to provide

Constellation is intended to give Docker nodes stronger shared behavior around:

- node and service discovery
- routing generation
- selected leadership or lease decisions
- DNS updates
- health-aware coordination

That makes it a candidate control layer above raw Compose.

## Preconditions for using this guide honestly

Before following this guide in anger, the operator should be clear about three
things:

1. The root repo is still Compose-first.
2. Constellation is a parallel control-plane effort, not automatically the
   authoritative runtime.
3. The deployment outcome must be verified by runtime evidence, not inferred
   from the existence of this guide.

## Environment assumptions

This subtree expects roughly:

- Linux nodes with `systemd`
- Docker available
- Go available for building `infra/`
- Tailscale or equivalent peer reachability assumptions
- Cloudflare credentials if DNS behavior is being exercised

Those are environment assumptions, not proof that the full cluster behavior is
sound.

## Intended deployment flow

At a high level, the intended flow appears to be:

1. Build the Constellation binaries from `infra/`.
2. Install or copy the binaries to candidate nodes.
3. Configure secrets and environment material needed for cluster coordination.
4. Start the agent on one or more nodes.
5. Verify cluster formation, routing behavior, and any leadership-style actions.
6. Only then test whether the control layer usefully improves the actual
   multi-node Docker problem.

## What must be verified after deployment

A successful install is not enough.

The following questions matter more than "did the service start":

- did nodes actually discover each other
- is there a stable cluster view
- is routing generated from shared state rather than local guesswork
- does failure produce sane fallback behavior
- what services, if any, are safe to entrust to this control layer

## Relationship to the root stack

If Constellation is deployed, the operator should keep explicit notes on:

- which traffic still flows through root Compose-managed infrastructure
- which services are only being observed by Constellation
- which services are actually being acted on by Constellation
- whether authority is shared, overlapping, or transferred

Without that, deployment creates another hidden-truth problem instead of fixing
one.

## Better framing for this guide

Treat this document as:

- a deployment hypothesis
- a staging and verification aid
- a control-plane experiment manual

Do not treat it as:

- blanket production approval
- proof of complete HA
- proof that stateful correctness is solved

## Recommended companion docs

- [`README.md`](README.md)
- [`PROJECT_SUMMARY.md`](PROJECT_SUMMARY.md)
- [`SYSTEM_STATUS.md`](SYSTEM_STATUS.md)
- [`ROADMAP.md`](ROADMAP.md)
- [`../../STRATEGY.md`](../../STRATEGY.md)
