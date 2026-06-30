# Instruction Surfaces and Authority

This page exists because the repo now has enough aligned language to create a
new kind of failure:

- the dream is clear
- the warnings are clear
- the authoring rules are clear
- the roadmap is clear
- therefore the runtime must be close to solved

That conclusion is false.

The repo has become good at naming the wound.
It is still uneven at proving which parts of that wound have actually moved out
of the operator's head.

So this page is not a taxonomy exercise.
It is an authority map for deciding which files are allowed to say what without
illegally borrowing confidence from their neighbors.

## The actual question this page answers

This page answers:

> which files are really allowed to explain the multi-node ordinary-Docker,
> no-Swarm-by-default, local-first then peer-forward dream, and which files
> only constrain how we talk about, author, or plan around that dream?

It also answers the more practical question that tends to come up during real
work:

> when someone says "AGENTS, Copilot instructions, and `.cursorrules` are
> basically telling us what we want to do with multiple Docker nodes, no Swarm,
> failover, fallback, and wrong-node routing, right?", what is the exact honest
> answer?

The exact honest answer is:

> yes, but each file proves a different part of that sentence, and none of them
> proves the runtime already owns the whole behavior.

That distinction is not pedantry.
It is the difference between assimilating the repo and merely recognizing its
keywords.

The weaker neighboring question is:

> which files mention HA, failover, clusters, or Compose the most?

That weaker question is one reason earlier documentation became useless.
This repo has many files that speak near the problem.
Far fewer files are allowed to own the same class of truth.

## Why the authority map matters here

The user is frustrated partly because the broader ecosystem keeps collapsing
different truth classes together:

- architecture desire becomes runtime implication
- current runtime presence becomes distributed-behavior confidence
- planning coherence becomes "we basically know what to do"
- authoring rigor becomes proof of resilience

This repo can repeat exactly the same mistake internally if it does not keep
its authority boundaries harsh.

That is what this page is for.

## What this page is and is not allowed to prove

This page is allowed to prove:

- which files own which truth class
- which surfaces are primary for the user's actual dream
- which surfaces are primary for current runtime truth
- which files are only honesty, authoring, or planning constraints
- which kinds of sentences become illegal when they are spoken from the wrong
  surface

This page is not allowed to prove:

- that repeated language across files counts as implementation corroboration
- that a planning artifact is halfway to runtime proof
- that stricter docs imply more system-owned truth
- that the priority implementation already satisfies the dream

This is an authority wall, not a maturity badge.

## The dangerous sentence shape this page has to block

This repo has one especially dangerous sentence shape:

1. the archive names the wound correctly
2. the instruction files describe the dream clearly
3. the live Compose files show serious machinery
4. the planning docs name plausible repairs
5. the docs quietly blend all of that into a present-tense feeling of maturity

That fifth step is the failure.

If a paragraph depends on more than one authority class, it should say which
class is carrying which part of the sentence.
Otherwise the docs start sounding wiser than the runtime actually is.

## The authority map

The most important repo surfaces are not equal.

| Surface | Truth class | What it is authoritative about | What it must not be upgraded into |
| --- | --- | --- | --- |
| [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md) | Architecture-dream truth | The clearest direct statement of the desired operating contract | Runtime proof that wrong-node routing, failover, or stateful dignity already work |
| [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md) | Repo-level honesty wall | The repo's blunt benchmark, anti-fake-HA posture, and the user's actual complaint | A present-tense implementation report |
| [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md) | Runtime-anchor and operator-surface truth | What the priority implementation surface is, how it is validated, and what constraints exist while working in the repo | The main architecture manifesto |
| [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules) | Authoring-discipline truth | Commit rules, inline config preference, and healthcheck discipline | Distributed-systems correctness, route durability, or failover proof |
| [`knowledgebase/AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/AGENTS.md) | Documentation-honesty truth | How docs must separate runtime truth, planned truth, and archive pressure | Proof that the runtime now owns one more missing burden |
| [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml) | Live root runtime truth | What the priority implementation actually declares at the root right now | Generic multi-node correctness |
| Active files under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/) | Live fragment truth | What each live subsystem contributes to ingress, auth, observability, state, routing, and helpers | That all fragments are equally live, equally trusted, or equally mature |
| [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md) | Planning and promotion truth | Which missing truth-owning layers the repo already knows it needs to consider | That the current runtime already owns those layers |
| [`knowledgebase/research/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/knowledgebase/research/) and the source archive | Research-pressure truth | Why the repo keeps returning to `services.yaml`, peer-aware routing, stateful honesty, OpenSVC, Nomad, k3s, Swarm, and other options | Shipped behavior |

That table is a conflict resolver.

If two files seem to agree, do not ask first:

> doesn't that mean the claim is probably true?

Ask instead:

> are these files even allowed to speak the same class of truth?

Very often they are not.

## Answering the instruction-file question directly

The instruction surfaces should be read as a joined contract with separated
burdens.

If the question is:

> do AGENTS, Copilot instructions, and `.cursorrules` explain the multi-node
> Docker, no-Swarm, failover/fallback direction?

The answer is:

> yes, but `copilot-instructions.md` explains the dream, `AGENTS.md` explains
> where current implementation truth has to be checked, and `.cursorrules`
> explains how changes must be authored without weakening the stack. They are
> not three equal witnesses for runtime capability.

Read them this way:

| File | What it contributes to the answer | What it cannot contribute |
| --- | --- | --- |
| `.github/copilot-instructions.md` | The target behavior: ordinary Docker nodes, no central orchestrator by default, local-first serving, peer-forward fallback, L7/L4 separation, anti-SPOF pressure without fake HA language. | Proof that any generic wrong-node request currently succeeds. |
| `AGENTS.md` | The implementation anchor: root `docker-compose.yml`, `compose/`, Go infra tooling, telemetry auth, validation commands, and environment gotchas. | The full architecture manifesto or evidence that the target behavior is live. |
| `.cursorrules` | The authoring discipline: commit changes, prefer inline Compose configs, preserve meaningful healthchecks, and avoid papering over service failure. | Distributed correctness, route durability, middleware continuity, or stateful failover. |

So the instruction files do explain the desired direction.
They do not, by themselves, answer the harder operational question:

> when traffic lands on a healthy node that does not host the service, what
> exact shared truth does that node consult, which peer does it choose, and what
> proof shows the route still means the same thing after handoff?

That second question must come from runtime evidence, proof packets, and
route-specific drills.
If a page answers it only by pointing back at instructions, it is borrowing
confidence from intent.

## The instruction surfaces as a burden split

The three recurring instruction files divide the user's frustration into
different burdens:

- `copilot-instructions.md` names the missing behavior.
- `AGENTS.md` names the surfaces that must be inspected before claiming that
  behavior is live.
- `.cursorrules` names the minimum authoring hygiene required not to make the
  behavior less trustworthy while changing Compose.

That burden split matters because the user is not asking for a slogan like
"multi-node Docker failover."
The user is asking for a system that stops forcing one human to be:

- the placement registry
- the peer selector
- the fallback explainer
- the middleware-continuity auditor
- the stateful-authority caveat keeper

The instruction surfaces can help preserve that demand.
They cannot retire those human jobs until live mechanisms and proof packets
show that the work moved into the system.

## The strongest honest reading orders

There are two serious reading orders in this repo.

### If the question is: what is the user actually trying to build?

Read in this order:

1. [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
2. [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)
3. [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
4. the research and archive pages
5. only then the live Compose surfaces as the reality check

That order answers:

> what dream is the repo trying to earn, and what accusation is it making
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

The repo needs both orders because the dream is much broader than the currently
proven runtime behavior.

## The clearest dream surfaces

The clearest architecture-dream surfaces are still:

- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- [`README.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/README.md)

`copilot-instructions.md` is the sharpest direct statement of the operating
contract.
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

That boundary must stay explicit or the repo will start using clarity of
desire as a substitute for proof of delivery.

## The repo-level honesty wall

The README matters because it keeps the repo from calming the problem down into
easier questions like:

- which orchestrator is best?
- how do we improve HA?
- how do we modernize the stack?

Its real benchmark is harsher:

> which options are still real once traffic lands on the wrong node and the
> operator is no longer allowed to privately complete the topology story?

That is why README authority is higher than AGENTS authority for the user's
complaint even though AGENTS is closer to implementation work.

## The runtime-anchor surfaces

The runtime-anchor layer is owned by:

- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active files under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)

These surfaces answer questions like:

- what is the priority implementation surface?
- what services, configs, networks, and secrets actually exist now?
- what commands validate the authored graph?
- what operator constraints exist while working here?

Examples that belong here now:

- root services such as `mongodb`, `redis`, `code-server`, `searxng`,
  `homepage`, and `dozzle`
- edge services such as `traefik`, `tinyauth`,
  `nginx-traefik-extensions`, `crowdsec`, `cloudflare-ddns`,
  `docker-gen-failover`, and `whoami`
- Headscale surfaces in `compose/docker-compose.headscale.yml`
- monitoring surfaces in `compose/docker-compose.metrics.yml`

These files prove the stack is real.
They do not prove:

- generic wrong-node routing correctness
- backend-loss persistence
- protected-route parity after peer-forward handoff
- stateful HA dignity

That is why runtime presence must stay separate from distributed-capability
claims.

## The authoring-discipline surfaces

The authoring-discipline layer is owned mainly by:

- [`.cursorrules`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.cursorrules)
- the authoring sections inside
  [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)

These surfaces are authoritative about:

- commit every change
- prefer inline Compose configs
- require actual healthchecks
- do not weaken healthchecks just to make the stack look healthy

Those rules matter because bad authoring discipline makes the anti-SPOF problem
harder to answer honestly.

But these surfaces are not allowed to cash themselves out into:

- peer-forward correctness
- middleware continuity after wrong-node handoff
- route durability under backend loss
- stateful failover semantics

This is one of the most important boundaries in the repo:

> good Compose hygiene is necessary, but it is not distributed-systems truth.

## The planning surfaces

The planning layer is owned by:

- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- related planning artifacts under `docs/`
- roadmap and gap pages in the knowledgebase

These surfaces are allowed to say:

- which missing truths the repo already knows it lacks
- which future layers might earn promotion
- which hidden operator burdens are still alive
- where the repo suspects the missing middle may live

They are not allowed to say:

- that those surfaces already exist in the current runtime
- that a coherent module list is halfway to shipping
- that a live root `services.yaml` or equivalent is already runtime authority
- that peer-forward or failover is already broadly proven

This matters because this repo's planning layer is unusually good.
That makes it unusually dangerous.

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
- why stateful honesty stays harsher than ingress optimism

They are not authoritative about current shipped behavior.

That sounds obvious.
It is still one of the easiest mistakes to make here, because the archive often
describes the real problem more vividly than the runtime proves its current
answer.

## The private sentence test for authority

Whenever a strong sentence appears in the docs, ask:

> which surface was actually allowed to say this part, and what private
> sentence would still survive if I forced the claim back through the runtime?

Examples:

- if the sentence came from dream surfaces, the surviving private sentence is
  often:
  `yes, but I still do not know what the runtime owns today`
- if the sentence came from runtime surfaces, the survivor is often:
  `yes, but I still do not know whether the system can explain itself on the
  bad day`
- if the sentence came from plans:
  `yes, but I still do not know whether the repair exists outside the plan`
- if the sentence came from archive pressure:
  `yes, but I still do not know what the current worktree truly proves`

If that reflex disappears, the authority map has failed.

## The honest bottom line

The strongest current authority reading is:

- the dream lives most clearly in `copilot-instructions.md`
- the repo-level honesty wall lives most clearly in `README.md`
- current runtime truth lives in root Compose plus active fragments
- planning authority lives in `INFRASTRUCTURE_MASTER_PLAN.md` and related docs
- archive pressure lives in research and source-archive surfaces

The docs get better, not worse, when those layers stay separate.

The reflex this page should leave behind is:

> before trusting a strong sentence, ask which surface was actually allowed to
> say that part.
