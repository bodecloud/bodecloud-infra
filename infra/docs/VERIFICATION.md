# Constellation Agent Verification Runbook

This runbook is for verifying Constellation behavior if you choose to exercise
the `infra/` control-plane path.

It is **not** evidence that the verification has already been completed.

## What this runbook is for

Use this file to answer:

- did the cluster actually form
- did nodes share the expected truth
- did routing and leadership behavior respond sanely
- what failed when real node or service loss was introduced

That is a much more useful purpose than repeating "zero-SPOF" slogans.

## What verification must keep separate

Do not collapse these into one checkbox:

### 1. Cluster formation

Can the nodes discover each other and maintain a stable shared view?

### 2. Control-plane behavior

Can leadership, lease, DNS, or routing-generation behavior avoid obvious split
brain or node-local blindness?

### 3. Stateless request-path failover

Can representative HTTP or TCP traffic still succeed when the original serving
path disappears?

### 4. Stateful correctness

Do data-bearing services remain correct under failover, not just reachable?

That last category must be verified separately.

## Minimum verification categories

If this subtree is going to earn stronger trust, verification should cover at
least:

1. peer discovery and cluster membership
2. stable leader or lease behavior where relevant
3. routing generation from shared state
4. representative wrong-node entry behavior
5. backend loss and route survival
6. separate stateful-service failover tests

## Verification posture

Expected outcomes should be recorded as:

- proven
- contradicted
- incomplete
- not yet tested

Not as:

- probably works
- implied by architecture
- assumed from implementation effort

## Why this matters in this repo

`bolabaden-infra` is specifically trying to avoid fake-HA storytelling.

A verification document that starts by assuming the answer is already "zero
SPOF" is the wrong shape for this repo.

The job here is to produce evidence strong enough to justify later claims, not
to front-load those claims into the verification instructions themselves.
