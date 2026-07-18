# Copilot Instructions for `bolabaden-infra`

This is the docs-side mirror of `.github/copilot-instructions.md`.
If this file conflicts with [`../INFRASTRUCTURE_MASTER_PLAN.md`](../INFRASTRUCTURE_MASTER_PLAN.md)
or [`../knowledgebase/architecture/instruction-surfaces-and-authority.md`](../knowledgebase/architecture/instruction-surfaces-and-authority.md),
the master plan and authority map win.

## Read This Repo Correctly

This codebase powers `bolabaden.org` as a **Compose-first, multi-node Docker
infrastructure repo**.

The important phrase is not merely "multi-node."
It is:

> multi-node Docker infrastructure that is trying to become anti-SPOF and
> peer-aware without immediately collapsing into Kubernetes, Docker Swarm, or
> another heavyweight orchestrator.

If you read the repo as "just a bunch of Compose files," you will miss the
actual goal.
If you read it as "already a finished distributed control plane," you will lie
about what the worktree currently proves.

## The Architecture Dream

The intended direction is:

- no central orchestrator by default
- services manually assigned to nodes
- current-state truth preferred over scheduler-declared desired state
- local-first serving when the requested service is already on the receiving
  node
- peer-forward fallback when the request lands on a healthy node that does not
  host the target service locally
- explicit separation between L7 HTTP behavior and L4 or raw TCP behavior
- anti-SPOF pressure without fake HA language

The core desired request model is:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That is the **target operating contract**.
Do not silently upgrade that target into "already proven runtime behavior."

## What Is Intent vs What Is Live

This file is one of the strongest **intent** surfaces in the repo.
It can support claims like:

- the repo wants multi-node Docker without Kubernetes or Docker Swarm by
  default
- the repo wants a lightweight current-state registry concept such as
  `services.yaml`
- the repo wants local-first then peer-forward request behavior
- the repo wants Cloudflare to support any-node entry rather than one sacred
  public box
- the repo wants Traefik and related edge tooling to preserve policy and
  middleware across routed services

This file does **not** prove:

- that the tracked root runtime currently ships a live root `services.yaml`
- that wrong-node requests already succeed generically
- that peer-forward fallback survives backend-loss conditions
- that middleware or auth continuity under peer fallback is fully proven
- that TCP failover is solved
- that stateful services are honestly HA just because they are reachable

If you need live proof, start with:

- `docker-compose.yml`
- `compose/docker-compose.*.yml`
- `docker compose config`
- the knowledgebase pages under `knowledgebase/architecture/`

## Current-State Registry Philosophy

The repo repeatedly converges on a lightweight current-state registry such as
`services.yaml`.

The intended meaning is:

- the system records where services actually live
- routing can consume current placement truth
- operators keep a readable source of placement knowledge
- the repo avoids a heavyweight scheduler unless it truly earns its keep

Example target shape:

```yaml
http:
  dozzle.bolabaden.org:
    backends:
      - host: node1.bolabaden.org
        port: 8080
      - host: node3.bolabaden.org
        port: 8080

tcp:
  redis-main:
    port: 6379
    backends:
      - host: node1.bolabaden.org
        port: 6379
```

Important boundary:

- treat `services.yaml` as architecture intent unless the tracked runtime
  actually ships and consumes it

## Routing Philosophy

### L7 / HTTP(S)

The intended HTTP stack centers on:

- Traefik v3
- health-aware routing
- auth continuity
- middleware continuity
- primary versus fallback behavior that remains visible to operators

The proxy layer is not supposed to merely expose containers.
It is supposed to preserve the meaning of a request path even when locality is
not available.

### L4 / TCP

Raw TCP services such as Redis or MongoDB are a different class.

Do not assume:

- HTTP failover logic automatically solves TCP forwarding
- node reachability equals stateful correctness
- a proxy path equals trustworthy failover

L4 and stateful systems require stricter language and stronger proof.

## Public Entry Philosophy

Cloudflare is part of the anti-SPOF story, but only as the first hop.

The intended node-entry model is:

- multiple public A or AAAA records
- any healthy public node can receive the first request
- no single reverse-proxy machine should quietly become the sacred public
  entrypoint

Do not confuse:

- DNS can hit more than one box

with:

- the request will be preserved correctly all the way to the right service

Those are different claims.

## Compose Authoring Expectations

### Root implementation priority

The priority implementation still centers on:

- `docker-compose.yml`
- `compose/docker-compose.*.yml`

This is a Compose-first repo in actual authoring shape, not just in rhetoric.

### Inline configs preferred

Prefer inline Compose `configs:` content over external mounted config-file
sprawl whenever practical.

Preferred:

```yaml
configs:
  example-config:
    content: |
      multiline
      config
      content
```

Avoid when a simpler inline config is possible:

```yaml
configs:
  example-config:
    file: ./path/to/file.conf
```

and especially:

```yaml
services:
  example:
    volumes:
      - ./path/to/file.conf:/etc/example/file.conf:ro
```

### Dollar signs in inline configs

In `configs:` content:

- `${VAR}` means Compose interpolation
- `$$` means a literal `$` should survive into the rendered config

Example:

```yaml
configs:
  grafana.ini:
    content: |
      instance_name = $$HOSTNAME
      port = ${GRAFANA_PORT:-3000}
```

### Healthchecks are mandatory

Never "fix" a service by removing or weakening its healthcheck.

Bad patterns:

- disabled healthchecks
- commented-out healthchecks
- TCP-only checks like `nc -z`

Preferred pattern:

```yaml
healthcheck:
  test: ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/health || exit 1"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 30s
labels:
  deunhealth.restart.on.unhealthy: "true"
```

Healthchecks are part of the no-fake-HA discipline.
They still do not prove cross-node resilience by themselves.

## Working Rules for Contributors and Agents

When making changes, preserve these reading rules:

- do not let target architecture language masquerade as current proof
- do not describe DNS redundancy as end-to-end failover
- do not describe a proxy surface as stateful HA
- do not assume local health implies peer-forward eligibility
- do not claim wrong-node success unless the route, peer choice, middleware,
  auth, and application behavior have actually been shown to survive it

When in doubt, prefer these stronger questions:

- what runs where right now?
- how would the receiving node know that?
- what survives if the local backend dies?
- what stays true if the request lands on the wrong node?
- what class of service am I talking about: stateless HTTP, TCP, or
  state-bearing?

## Practical Validation Baseline

Use these as the minimum validation surface:

```bash
docker compose config --quiet
docker compose config --services
python3 -m mkdocs build -f mkdocs.yml --strict
```

Remember:

- Compose validation requires a prepared env and secret surface
- `~/.docker/config.json` may need to exist even if it contains only `{}`
- passing validation proves authored shape, not distributed correctness

## Bottom Line

The repo's central question is not:

> how do we host more services?

It is:

> how do we keep the directness of Compose while making multiple ordinary
> Docker nodes behave less stupidly under wrong-node entry, backend loss, and
> anti-SPOF pressure?

Every contribution should either:

- make that dream more true
- make the proof boundary more honest
- or make the operator surface more readable

Anything else is infrastructure theater.
