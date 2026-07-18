# Constellation Agent Docs: Read This Subtree Carefully

This subtree documents the Go-based `infra/` experiment usually referred to as
"Constellation Agent."

Historically, many files here were written as if Constellation were already a
fully proven, production-ready, zero-SPOF orchestration layer for the whole
repo.

That is too strong.

## What this subtree is useful for

Use `infra/docs/` to understand:

- what the Go control-plane effort is trying to become
- which architectural jobs it is meant to own
- what components have been coded
- how the experiment differs from the root Compose-first runtime

This subtree is valuable as:

- implementation direction
- codebase orientation
- architectural pressure evidence

## What this subtree does not prove by itself

Reading these docs does **not** prove that:

- Constellation is the adopted runtime for `bolabaden-infra`
- the repo has already achieved zero-SPOF behavior
- cross-node failover is fully validated in real deployment conditions
- stateful backends are operationally proven under this control plane
- the root `docker-compose.yml` implementation has been replaced

Those claims require live runtime evidence, not just internal design docs.

## Current relationship to the rest of the repo

The priority implementation of `bolabaden-infra` is still rooted in the top
level [`docker-compose.yml`](../docker-compose.yml) and its include set.

Constellation should currently be read as:

- a serious attempt to build a stronger coordination layer around Docker
- a candidate answer to placement truth, routing truth, and failover truth
- an experiment that still needs explicit proof boundaries

## Recommended reading order

1. [`../../README.md`](../../README.md)
2. [`../../STRATEGY.md`](../../STRATEGY.md)
3. [`../../knowledgebase/architecture/compose-first-architecture.md`](../../knowledgebase/architecture/compose-first-architecture.md)
4. [`PROJECT_SUMMARY.md`](PROJECT_SUMMARY.md)
5. [`ARCHITECTURE.md`](ARCHITECTURE.md)
6. [`ROADMAP.md`](ROADMAP.md)

## Reading rule for this subtree

Unless a document explicitly cites verified runtime evidence, treat it as one
of:

- intended design
- coded capability
- planned capability
- experimental claim

Do not automatically upgrade it to:

- production proof
- global repo truth
- complete HA verification
