# bolabaden Infrastructure Knowledgebase

This site exists to answer one unusually stubborn infrastructure question:

> how do you keep
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> as the real human control surface, spread services across several ordinary
> Docker nodes, and still make wrong-node traffic, fallback, and anti-SPOF
> behavior feel like one coherent platform instead of one operator privately
> remembering the real answer?

That is the real subject of `bolabaden-infra`.

This is not primarily:

- a generic self-hosting site
- a broad HA survey
- an orchestrator comparison hub
- a "how to modernize your homelab" handbook

Those are neighboring topics.
They are not the main wound the repo is trying to close.

The repo's recurring complaint is sharper than that:

- ordinary Docker answers often stop at static glue, DNS plurality, or local
  proxying
- heavyweight answers often jump straight to Swarm, Kubernetes, k3s, Nomad, or
  another larger control plane
- both categories can sound respectable while still leaving the operator as the
  final keeper of placement, failover, and route-meaning truth

The user's dream is to find out whether there is an honest middle layer:

- still Compose-first
- still human-legible
- still based on ordinary Docker nodes
- but no longer dependent on one person silently remembering what should happen
  when traffic lands on the wrong healthy machine

## The architecture dream

The strongest intent surface in the repo is
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It says the intended direction is:

- no central orchestrator by default
- services manually assigned to nodes
- current-state truth preferred over scheduler-declared desired state
- local-first serving when the requested service already lives on the receiving
  node
- peer-forward fallback when the receiving node is healthy but the service is
  remote
- explicit separation between HTTP routing and raw TCP or stateful behavior
- anti-SPOF pressure without fake HA language

The target request contract is therefore:

```text
User -> Cloudflare DNS -> any surviving public node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that really hosts it
```

That contract is the dream.
It is not the same thing as current proof.

## What the site is for

This site is supposed to help a reader answer three different questions without
collapsing them together:

1. what does the root Compose runtime actually contain today?
2. what architecture is the repo clearly trying to grow into?
3. which missing truth layers still force the operator to privately complete
   the story during wrong-node entry, backend loss, or stateful recovery?

Most infrastructure docs fail here by making the stack sound progressively more
adult as the terminology gets better.
This site is trying to do the opposite:
improve clarity without laundering partial machinery into fake closure.

## What this site is and is not allowed to prove

This site is authoritative about:

- the repo's actual dream
- the current root Compose implementation surface
- the difference between live runtime truth, repo-native intent, planning
  pressure, and archive pressure
- the exact gaps between today's stack and genuine wrong-node recovery
- which stronger claims still need proof before they become legal

This site is not allowed to:

- claim the current runtime already behaves like the dream
- turn a good explanation into failover proof
- promote plans or research into shipped behavior
- imply that better organization means the platform now owns more truth

The user does not mainly need one more smooth infrastructure story.
They need the docs to stop helping the same lie that keeps happening elsewhere:

- the stack looks serious
- the options sound rich
- the route names sound mature
- therefore the system must finally be close to handling the bad day on its own

That conclusion is exactly what this site has to keep blocking when the
evidence does not support it.

## Strongest honest current answer

`bolabaden-infra` already contains a serious Compose-first platform:

- a real root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active includes under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)
- a substantial Traefik, CrowdSec, TinyAuth, and nginx-auth edge layer
- operator, observability, and maintenance surfaces
- private-mesh pressure through Headscale
- repeated repo pressure toward any-node entry, peer-aware routing, and
  anti-SPOF behavior

What it still does **not** prove is the thing the user actually cares about:

- that any healthy public node can accept a request and preserve it correctly
  when the service is remote
- that placement truth is shared explicitly instead of remembered
- that peer eligibility is system-owned rather than guessed from reachability
- that fallback routes survive the failure they are meant to absorb
- that auth, middleware, and request semantics survive peer handoff
- that stateful services are genuinely resilient rather than merely reachable

The dream is clear.
The stack is real.
The missing truth-owning middle layer is still incomplete.

That three-part sentence is the checksum for the whole site.
If a page preserves only two of the three, it usually becomes misleading in the
same way the user is already exhausted by.

## The hidden job the operator still performs

The shortest honest summary of the current wound is:

the operator is still acting like the missing control plane.

That hidden job includes things like:

- remembering what runs on which node right now
- remembering which peer is actually safe to forward to
- remembering whether the fallback path still exists under failure
- remembering whether a reachable answer is only transport-reachable or
  semantically valid
- remembering which stateful surfaces still hide a sacred authority node

Most nearby tools improve machinery, naming, or automation without clearly
shrinking that role.
That is why so many respectable answers still feel useless in this repo.

## How to read the site correctly

Keep four authority layers separate:

- live runtime truth:
  the root
  [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml),
  included files under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose),
  and checks such as `docker compose config`
- repo-native intent:
  especially
  [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)
- planning truth:
  especially
  [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)
- archive and research pressure:
  the pages under [research](research/evidence-ledger.md) and the archive map in
  [Source Assimilation Index](operations/source-assimilation-index.md)

The site gets misleading the moment those four layers blend into one calm
voice.

## Start here

If you want the shortest route through the site:

1. [Problem, Pressure, and Goals](architecture/problem-and-goals.md)
2. [Current Compose Runtime](architecture/current-compose-runtime.md)
3. [HA, Failover, and Routing](architecture/ha-failover-routing.md)
4. [Stateful HA and Data](architecture/stateful-ha-and-data.md)
5. [Capability Gaps and Roadmap](architecture/capability-gaps-and-roadmap.md)
6. [DevOps Runbook](operations/devops-runbook.md)

If you want the proof pressure behind those pages:

- [Evidence Ledger](research/evidence-ledger.md)
- [Ingress and Failover Evidence](research/ingress-and-failover-evidence.md)
- [Stateful HA Evidence](research/stateful-ha-evidence.md)
- [Orchestrator Tradeoffs Evidence](research/orchestrator-tradeoffs-evidence.md)

## What still does not count

This site should keep a few false conclusions illegal:

- "there are multiple public nodes, so failover is basically solved"
- "Traefik exists, so wrong-node forwarding must be close"
- "a helper has failover in the name, so the platform owns fallback now"
- "a route can be rendered, so it must survive the failure that made it matter"
- "the docs are clearer now, so the platform must be closer to adulthood"

Those are exactly the overreads this site is supposed to interrupt.

The useful emotional conclusion is not that the repo is simple.
It is that the repo is serious, the dream is specific, and the missing burden
transfer is still brutally concrete.
