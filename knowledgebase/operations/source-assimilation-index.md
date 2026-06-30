# Source Assimilation Index

This page is the retrieval contract behind the knowledgebase.

Its job is simple:

> stop the docs from sounding deeply informed while still quietly answering a
> smaller, safer version of the user's actual problem.

`bolabaden-infra` contains several different evidence classes at once:

- live root Compose runtime
- repo-native dream and honesty surfaces
- planning docs with stronger language than current proof supports
- orchestration side paths that may matter later
- archive pressure that explains why the ordinary answers keep failing

If those are flattened into one voice, the site becomes polished and useless.

## What this page is and is not allowed to prove

This page is authoritative about:

- which evidence classes outrank others for different claims
- what "actually RAG this time" means in this repo
- how to reconstruct the architecture dream without smuggling in false proof

This page is not authoritative about:

- whether a specific runtime claim is true
- whether one summary page already solved the reconstruction problem
- whether archive pressure can stand in for live implementation evidence

This page defines the retrieval rules. It is not itself runtime proof.

## Strongest honest current answer

The right retrieval strategy for this repo is not "cite many files." It is:

1. identify what class of claim is being made
2. route that claim to the highest-authority source class first
3. say what that source proves
4. say what it does not prove
5. preserve contradiction instead of ironing it out
6. keep the current worktree above elegant prose for runtime claims
7. use archive pressure to recover the user's real complaint without letting
   the archive impersonate runtime truth

If a page skips that sequence, it can become longer, calmer, and more wrong.

## Primary source classes

### 1. Dream and authority sources

These define what the repo is trying to become.

Start with:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
4. [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)

Use this class when the claim is:

- what the platform wants to do at request time
- whether the repo resists immediate Swarm or Kubernetes promotion
- which instruction files actually carry architecture authority

### 2. Live runtime sources

These define the strongest current implementation truth.

Start with:

1. root [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. active fragments under
   [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
3. validation guidance and constraints in
   [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)

Important current live surfaces include:

- root networks: `publicnet`, `backend`, `warp-nat-net`
- root-owned services such as `mongodb`, `redis`, `searxng`, `code-server`,
  `homepage`, `watchtower`, `dozzle`, `portainer`, `dns-server`,
  `telemetry-auth`, `bolabaden-nextjs`, and `biodecompwarehouse*`
- active include fragments:
  - `compose/docker-compose.coolify-proxy.yml`
  - `compose/docker-compose.docs.yml`
  - `compose/docker-compose.firecrawl.yml`
  - `compose/docker-compose.headscale.yml`
  - `compose/docker-compose.llm.yml`
  - `compose/docker-compose.metrics.yml`
  - `compose/docker-compose.stremio-group.yml`
  - `compose/docker-compose.warp-nat-routing.yml`
  - `compose/docker-compose.wishlist.yml`

Use this class when the claim is:

- what exists today
- what the root stack actually owns
- what helper surfaces are present now

### 3. Planning and promotion sources

These define what the repo is seriously considering or trying to earn next.

Start with:

1. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
2. [`docs/stateful_ha_plan.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/stateful_ha_plan.md)
3. [`docs/osvc_ingress_ha.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/osvc_ingress_ha.md)
4. research pages such as
   [Ingress and Failover Evidence](../research/ingress-and-failover-evidence.md),
   [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md),
   and [Stateful HA Evidence](../research/stateful-ha-evidence.md)

Use this class when the claim is:

- what missing layer might be added
- what bugs or gaps the repo has already named explicitly
- what future promotion candidates exist

Examples already named in planning material:

- `docker-gen-failover` route deletion on container stop
- non-functional `watchtower`
- Cloudflare DDNS multi-record problems
- manual secret sync
- manual compose sync
- unproven service failover
- Headscale HA / leader-election ideas

### 4. Archive and pressure sources

These explain why the ordinary answers still feel fake to the user.

Start with:

- [Archive Pressure Patterns](../research/archive-pressure-patterns.md)
- [User Intent and Dream](../research/user-intent-and-dream.md)
- [Evidence Ledger](../research/evidence-ledger.md)

Use this class when the claim is:

- what the user is rebelling against
- why "just use Kubernetes" is not an adequate answer by itself
- why the repo keeps circling a missing middle layer

## Claim routing rules

| If the sentence is really claiming... | Start with | It still must not imply... |
| --- | --- | --- |
| "this exists in the live stack" | root Compose runtime | that it survives the relevant failure |
| "this is the repo's target behavior" | `copilot-instructions.md` and `README.md` | that the runtime already earned it |
| "this future layer is being seriously explored" | planning docs and research pages | that it has been promoted or proven |
| "the user is rejecting this whole category of answer" | archive-pressure pages | that pressure alone identifies the winning implementation |

If a page sounds balanced because it mixes these classes without naming them,
the retrieval contract has failed.

## What "actually RAG this time" means here

In this repo, real retrieval means all of the following:

1. begin with the question, not the folder name
2. identify the claim type before looking for evidence
3. route live claims to the worktree first
4. route dream claims to the instruction surfaces first
5. route future claims to planning material first
6. bring in archive pressure only to restore the real complaint, not to fake
   proof
7. say where the system still depends on human reconstruction

That last point matters a lot.

The user is not only asking whether the facts can be found.
They are asking whether the system itself owns enough shared truth that the
operator no longer has to privately join everything together.

If the docs cannot say where private operator memory is still acting like the
real control plane, they have not assimilated the repo deeply enough.
