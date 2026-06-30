# Knowledgebase Agent Notes

This file is for agents and contributors touching `knowledgebase/`.

It is intentionally excluded from the rendered MkDocs site. The public entrypoint is [`index.md`](index.md), not this file.

## Primary rule

Do not let documentation become more certain than the implementation.

That is still not strong enough on its own.

Also do not let documentation sound more system-owned than the system really
is.

This repo is unusually vulnerable to a specific failure mode:

- a config parses
- a route exists
- a service starts
- and the docs quietly upgrade that into "the multi-node architecture works"

That upgrade is exactly what this knowledgebase is meant to prevent.

There is a second failure mode sitting right beside it:

- a path looks understandable after the fact
- the page explains it smoothly
- the smooth explanation quietly depends on private operator reconstruction
- the docs start sounding like the runtime itself owns that truth

That is how useful-seeming documentation turns back into architecture theater.

## The repo's real problem

The user is not merely asking for nicer docs.

They are trying to discover whether multiple ordinary Docker nodes can be turned into something genuinely resilient and flexible without immediately defaulting to:

- Docker Swarm
- Kubernetes or k3s
- a heavyweight centralized control plane
- opaque platform abstractions that hide the true behavior

When editing docs, stay grounded in that problem.

Do not flatten it into:

- better HA in general
- better clustering in general
- better orchestrator selection
- better platform maturity language

All of those are smaller and calmer than the actual question.

## Truth layers

All substantive documentation should stay legible across these three layers:

### Live implementation truth

Grounded directly in:

- [`../docker-compose.yml`](../docker-compose.yml)
- included files under [`../compose/`](../compose)
- current root labels, configs, networks, secrets, and services
- commands such as `docker compose config`

### Planned architecture truth

Grounded in:

- [`../docs/INFRASTRUCTURE_MASTER_PLAN.md`](../docs/INFRASTRUCTURE_MASTER_PLAN.md)
- other explicit plan/design docs under `../docs/`

### Research-pressure truth

Grounded in:

- [`source-archive/`](source-archive)
- synthesized research pages under [`research/`](research)

Do not collapse these into one timeline.

Also do not let them collapse into one voice.

In this repo, a calm unified voice is often a warning sign that one of the
truth layers has been allowed to impersonate the others.

## Pages with the highest leverage

When in doubt, prefer improving these pages:

- [`index.md`](index.md)
- [`architecture/problem-and-goals.md`](architecture/problem-and-goals.md)
- [`architecture/current-compose-runtime.md`](architecture/current-compose-runtime.md)
- [`architecture/ha-failover-routing.md`](architecture/ha-failover-routing.md)
- [`architecture/stateful-ha-and-data.md`](architecture/stateful-ha-and-data.md)
- [`architecture/capability-gaps-and-roadmap.md`](architecture/capability-gaps-and-roadmap.md)
- [`operations/devops-runbook.md`](operations/devops-runbook.md)

Those pages define how the whole repo is read.

## Important recurring honesty checks

Keep these boundaries explicit unless new evidence truly changes them:

- a root `services.yaml` is conceptually central, but not currently a live tracked root source of truth
- `docker-gen-failover` exists, but has been documented as a weak or defective mechanism rather than solved HA
- `watchtower` exists, but presence is not proof of safe update automation
- DNS failover is not the same as service-level failover
- HTTP failover is not the same as stateful HA
- Traefik labels and hostnames do not by themselves prove peer-aware routing correctness

Keep one more boundary explicit everywhere:

- an apparently successful request is not self-explaining until the page can say
  what truth the runtime owned directly versus what truth a human still supplied
  by reconstruction

## Validation habit

At minimum, after knowledgebase edits run:

```bash
python3 -m mkdocs build -f mkdocs.yml --strict
```

When claims touch the live Compose surface, also verify with Compose commands.

When claims touch request preservation, fallback, or ownership of cluster truth,
do not stop at shape-validation evidence alone.

Those claims need pages to state what is directly proven, what is only planned,
and what still depends on hidden operator burden.

## Bottom line

The right standard for these docs is not "does this read smoothly?"

It is:

> would an operator make a better next decision after reading this, and would that decision still be honest if the implementation were inspected directly?

If the answer only stays persuasive while the operator privately remembers which
node is really special, the page is still not done.
