# Instruction Surfaces and Authority

This page exists because the repo now has enough aligned language to create a
new kind of failure:

- the dream is clear
- the warnings are clear
- the authoring rules are clear
- the roadmap is clear
- therefore the runtime must be close to solved

That conclusion is false.

This repo has become good at naming the wound.
It is still uneven at proving which parts of the wound have actually moved out
of the operator's head.

So this page is not a taxonomy exercise.
It is an authority map for deciding which files are allowed to say what.

## The exact question this page answers

This page answers:

> which files actually explain the multi-node ordinary-Docker, no-Swarm by
> default, local-first then peer-forward dream, and which files only constrain
> how we talk about, author, or plan around that dream?

The weaker neighboring question is:

> which files mention HA, failover, clusters, or Compose the most?

That weaker question is one reason earlier docs became useless.

This repo has many files that mention adjacent concepts.
Far fewer files are allowed to own the same class of truth.

## Why the authority map matters here

The user is frustrated partly because the ecosystem keeps collapsing different
truth classes together:

- architecture desire becomes runtime implication
- current runtime presence becomes distributed-behavior confidence
- planning coherence becomes "we basically know what to do"
- authoring rigor becomes proof of resilience

This repo can accidentally repeat the same move internally if it does not keep
its own authority boundaries harsh.

That is what this page is for.

## What this page is and is not allowed to prove

This page is allowed to prove:

- which files own which truth class
- which surfaces are primary for the user's real dream
- which surfaces are primary for the current runtime
- which files are only honesty or authoring constraints
- which kinds of sentences are illegal from each surface

This page is not allowed to prove:

- that repeated language across files counts as implementation corroboration
- that a planning artifact is halfway to being runtime proof
- that stricter docs imply more system-owned truth
- that the priority implementation already satisfies the dream

## The authority map

The most important repo surfaces are not equal.

| Surface | Truth class | What it is authoritative about | What it must not be upgraded into |
| --- | --- | --- | --- |
| [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md) | Architecture-dream truth | The clearest direct statement of the desired operating contract | Runtime proof that wrong-node, fallback, or stateful dignity already work |
| [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md) | Repo-level honesty wall | The repo's blunt benchmark, anti-fake-HA posture, and explanation of the user's real complaint | A present-tense implementation report |
| [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md) | Runtime-anchor and operator-surface truth | What the priority implementation surface is, how it is validated, and what constraints exist while working in the repo | The main architecture manifesto |
| [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules) | Authoring-discipline truth | Commit rules, inline config preference, and healthcheck discipline | Distributed-systems correctness, route durability, or failover proof |
| [`knowledgebase/AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/AGENTS.md) | Documentation-honesty truth | How the docs must separate runtime truth, planned truth, and archive pressure | Proof that the runtime now owns one more missing burden |
| [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml) | Live root runtime truth | What the priority implementation actually declares at the root right now | Generic multi-node correctness |
| Active files under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/) | Live fragment truth | What each live subsystem currently contributes to ingress, auth, observability, state, routing, and sidecars | That all fragments are equally live, equally trusted, or equally mature |
| [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md) | Planning and promotion truth | Which missing truth-owning layers the repo already knows it needs to consider | That the current runtime already owns those layers |
| [`knowledgebase/research/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/) and the source archive | Research-pressure truth | Why the repo keeps returning to `services.yaml`, peer-aware routing, stateful honesty, OpenSVC, Nomad, k3s, Swarm, and other options | Shipped behavior |

That table should be read as a conflict resolver.

If two files seem to disagree, the first thing to ask is not "which one sounds
more complete?"

It is:

> are they even trying to answer the same truth class?

Very often they are not.

## The strongest honest reading orders

There are two serious reading orders in this repo, depending on the question.

### If the question is: what is the user actually trying to build?

Read in this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
4. the research and archive pages
5. only then the live Compose surfaces as the reality check

That order answers:

> what dream is the repo trying to earn, and what accusations is it making
> against the usual ecosystem answers?

### If the question is: what has the priority implementation actually earned?

Read in this order:

1. [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
2. active files under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)
3. `docker compose config`
4. route-specific drills, if any exist
5. only then planning and research as contrast

That order answers:

> what is actually present, inspectable, and currently claimable?

Those two orders are both necessary because the repo's dream is much broader
than its currently proven runtime behavior.

## The clearest dream surface

The clearest architecture-dream file remains
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).

It explicitly says the repo wants:

- no central orchestrator by default
- services manually assigned to nodes
- current-state truth preferred over desired-state theater
- local-first serving
- peer-forward fallback when the target is remote
- explicit separation between L7 HTTP and L4/raw TCP
- anti-SPOF pressure without fake HA language

It also states the target request contract directly:

```text
User -> Cloudflare DNS -> any surviving node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that currently hosts it
```

That makes it the strongest answer to:

> what does the user actually want the system to feel like?

It is not the strongest answer to:

> what already works today?

That boundary has to stay explicit or the repo will start using clarity of
desire as a substitute for proof of delivery.

## The repo-level honesty wall

[`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md) is the strongest
repo-level honesty wall around the dream.

Its real job is not onboarding.
Its real job is to keep the repo from calming the problem down into easier
questions like:

- which orchestrator is best?
- how do we improve HA?
- how do we modernize the stack?

The README preserves the harsher benchmark:

> which options are still real once traffic lands on the wrong node and the
> operator is no longer allowed to privately complete the topology story?

That is why README authority is higher than AGENTS authority for the user
complaint, even though AGENTS is closer to the implementation.

## The runtime-anchor surfaces

[`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md),
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml),
and the active fragment files under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)
own the runtime-anchor layer.

They answer questions like:

- what is the priority implementation surface?
- what services, configs, networks, and secrets actually exist now?
- what commands validate the authored graph?
- what operator constraints exist while working here?

Examples that belong to this layer now:

- the root include graph
- root networks like `publicnet`, `backend`, and `warp-nat-net`
- root services like `mongodb`, `redis`, `code-server`, `searxng`, `homepage`,
  `dozzle`, and others
- edge services in `compose/docker-compose.coolify-proxy.yml` like `traefik`,
  `tinyauth`, `nginx-traefik-extensions`, `crowdsec`, `cloudflare-ddns`,
  `docker-gen-failover`, and `whoami`
- Headscale surfaces in `compose/docker-compose.headscale.yml`
- monitoring surfaces in `compose/docker-compose.metrics.yml`

Those files prove the stack is real.

They do **not** prove:

- that wrong-node routing works generically
- that backend-loss fallback survives
- that protected routes preserve policy meaning across peer-forward handoff
- that stateful services have earned HA language

That is why runtime presence must stay separate from distributed capability.

## The authoring-discipline surfaces

[`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules) and parts of
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
govern how services should be authored.

They are authoritative about things like:

- commit every change
- prefer inline Compose configs
- require actual healthchecks
- do not weaken healthchecks just to make the stack look healthy

Those rules matter.
Bad authoring discipline makes the anti-SPOF question harder to answer
honestly.

But these surfaces are not allowed to cash themselves out into:

- peer-forward correctness
- middleware continuity after wrong-node handoff
- route durability under backend loss
- stateful failover semantics

This is one of the most important boundaries in the whole repo:

good Compose hygiene is necessary, but it is not distributed-systems truth.

## The planning surfaces

[`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
and related planning artifacts are authoritative about pressure and candidate
promotion paths, not completion.

They are allowed to say:

- which missing truths the repo already knows it lacks
- which future layers might earn promotion
- which hidden operator burdens are still alive
- where the repo suspects the missing middle layer might live

They are not allowed to say:

- that those surfaces already exist in the current runtime
- that a coherent module list is halfway to shipping
- that a live root `services.yaml` or equivalent is already runtime authority
- that peer-forward or failover is already broadly proven

This matters because this repo's planning layer is unusually good.
That makes it especially dangerous.

Detailed planning can sound like partial implementation when the reader is
already hungry for relief.

## The research-pressure surfaces

The research pages and source archive are not decorative.
They preserve why the repo keeps investigating many adjacent systems without
wanting to surrender immediately to a heavyweight scheduler.

They are authoritative about:

- recurring user frustration
- recurring ecosystem dead ends
- why helper piles often still feel fake
- why wrong-node dignity keeps outranking generic "more clustered" language
- why stateful honesty remains harsher than ingress optimism

They are not authoritative about current shipped behavior.

That sounds obvious.
It is still one of the easiest mistakes to make in this repo, because the
archive often describes the real problem more vividly than the runtime proves
its current answer.

## The dangerous blended sentence

This repo has one especially dangerous sentence shape:

1. the archive names the wound correctly
2. the instruction files describe the dream clearly
3. the live Compose files show serious machinery
4. the planning docs name plausible repairs
5. the docs quietly blend all of that into a present-tense feeling of maturity

That fifth step is the failure this page exists to block.

If a paragraph depends on more than one authority class, it should say which
class is carrying which part of the claim.

Otherwise the docs start sounding better than the system actually is.

## The current bottom line

The strongest current authority reading is:

- the dream lives most clearly in `copilot-instructions.md`
- the repo-level honesty wall lives most clearly in `README.md`
- current runtime truth lives in root Compose plus active fragments
- planning authority lives in `INFRASTRUCTURE_MASTER_PLAN.md` and related docs
- archive pressure lives in research and source-archive surfaces

The repo becomes more useful, not less, when those layers stay separate.

The reader should leave this page with one strict reflex:

> before trusting a strong sentence, ask which surface was actually allowed to
> say that part
