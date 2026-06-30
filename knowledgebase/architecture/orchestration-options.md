# Orchestration Options

For the evidence underneath this page, start with:

- [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
- [Orchestration Research 2026](../research/orchestration-research-2026.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This page exists to stop one lazy question from taking over the repo:

> which orchestrator is best?

That is not the real decision surface here.

The real question is:

> which extra control layer has actually earned the right to exist because it
> removes a real hidden human SPOF, a real wrong-node failure, a real stale
> topology lie, or a real state-authority ambiguity that the current
> Compose-first stack cannot solve honestly?

That is the orchestration question that matters in `bolabaden-infra`.

## The strongest honest current answer

The current default stance in this repo is:

- keep Compose as the human control surface for as long as that remains honest
- refuse heavyweight promotion by prestige alone
- add the smallest extra truth-owning layer that actually removes a hidden
  burden
- promote to a stronger controller only when the smaller layer cannot carry the
  real failure classes anymore

That stance is not anti-orchestrator.
It is anti-unearned-orchestrator.

## What this page is and is not allowed to prove

This page is allowed to prove:

- how orchestration options should be judged in this repo
- which burden a promoted layer would need to own to justify itself
- why product prestige is weaker than repo-specific burden transfer
- what candidate families are relevant to the real problem
- what proof packets are required before stronger platform language becomes
  legal

This page is not allowed to prove:

- which orchestrator has won globally
- whether the live runtime already demonstrates the promoted behavior
- whether ecosystem fashion should override repo-specific truth
- whether a more powerful controller is automatically a better answer

This is a promotion filter, not a winner declaration.

## The user's real decision is not logo selection

The user is not choosing between software brands.

They are trying to decide whether the dream can be preserved:

- ordinary Docker nodes
- Compose-first authorship
- visible inspectable behavior
- fewer secret sacred nodes
- fewer secret sacred humans
- real fallback and failover behavior
- no premature surrender to a heavyweight control plane

The orchestration decision only becomes interesting when one of those dreams
starts colliding with truths the current stack cannot own by itself.

## The accusation every option has to survive

Every candidate here should be read under the same accusation:

> are you actually removing the hidden human SPOF, or are you just replacing
> it with a more respectable story?

If the answer is the second one, the option may still be powerful.
It is not yet the right answer for this repo.

## The failure mode this page is trying to stop

The easiest way for a page like this to become dishonest is:

1. the option list gets more sophisticated
2. the controllers sound more grown up
3. the operator still privately carries the same decisive truth
4. the docs quietly call that progress

This page exists to stop prestige from masquerading as burden transfer.

## What "worldview tax" means here

`Worldview tax` is not a decorative phrase.
It means the option demands some combination of:

- more abstraction distance
- more invisible cluster state
- more controller trust
- more operational doctrine
- less direct reasoning from `docker-compose.yml` plus current runtime

That tax may be worth paying.
It is only worth paying when a real hidden burden moved in return.

## The non-negotiable questions every option must answer

Before any option gets described as the right next layer, it has to answer:

1. which exact hidden operator burden does it remove?
2. which truth does it own that the current runtime does not own cleanly?
3. how does the wrong node determine the right backend now?
4. how does the system distinguish a merely reachable peer from an eligible
   peer?
5. what evidence proves the route still exists after backend loss?
6. what evidence proves the route still means the same thing after handoff?
7. how does the option state its limits for TCP and stateful classes?
8. what worldview tax does the option impose in return?

If a candidate cannot answer those concretely, it may still be interesting.
It has not earned default gravity here.

## The private completion test

Every orchestration option should be judged by one direct question:

> after adopting this layer, what exact sentence should the operator no longer
> need to finish privately?

Examples of sentences this repo wants to kill:

- `That route only works because I know node3 is the real backend today.`
- `The proxy can reach node2, but I privately know node2 is the only safe one.`
- `The fallback config exists, but I know it disappears when the preferred
  backend dies.`
- `Redis is reachable through the edge, but I know that does not mean there is
  safe failover.`
- `The protected route still returns, but I know the auth semantics changed.`

If the layer cannot kill at least one such sentence cleanly, it has not yet
justified its existence.

## Main candidate families

### 1. Stay Compose-first and harden a shared-truth layer

Examples:

- tracked `services.yaml`
- node-local proxy plus shared registry
- peer-sync helpers
- mesh-assisted state broadcast

What this family is trying to preserve:

- `docker-compose.yml` stays central
- service intent remains human-legible
- nodes still feel like ordinary Docker machines
- the added layer only grows enough cluster truth to solve the wrong-node
  wound

What it can plausibly solve:

- current placement lookup
- inspectable local-vs-remote routing decisions
- a narrower burden transfer than a full scheduler

What it still has to prove:

- registry freshness
- disagreement handling
- peer-eligibility truth
- fallback durability under actual backend loss
- protected-route continuity
- strict limits for TCP and stateful claims

Private sentence still likely to survive unless proved otherwise:

> yes, but I still personally know whether the shared registry is current
> enough to trust.

### 2. Compose plus gossip or event-driven coordination

Examples from the archive:

- Serf-style membership and failure events
- peer-equal node agents reacting to cluster events
- distributed signals without immediately adopting a manager hierarchy

Why it stays attractive:

- it respects the user's instinct against early hierarchy
- it preserves the emotional dream of equal nodes
- it can reduce private failure detection
- it can distribute liveness signals faster than manual memory can

Why it stays risky:

- membership is not placement truth
- event distribution is not authority
- liveness is not eligibility
- gossip is much better at `who seems alive` than at
  `who is definitely the correct backend for this protected or stateful route`

Private sentence still likely to survive unless proved otherwise:

> yes, but I still personally know hearing the event is not the same thing as
> knowing the right backend.

### 3. Compose plus registry and dynamic proxy/control-plane helpers

Examples:

- Consul-like registries
- HAProxy or Envoy driven from discovery state
- helper APIs publishing backend identity and policy inputs

What this family can improve:

- service-name to backend resolution
- inspectable inputs for wrong-node forwarding
- better explanation surfaces for current choices

What it may still leave behind:

- stronger central trust concentration
- policy continuity gaps
- registry sacredness
- ambiguity about stateful authority

This family only becomes real if the receiving node can explain its choice from
shared truth rather than cultural memory.

Private sentence still likely to survive unless proved otherwise:

> yes, but I still personally know whether the registry or helper became the
> new sacred component.

### 4. Stronger orchestrators or cluster managers

Examples:

- Nomad
- OpenSVC
- k3s
- Kubernetes

What they can genuinely bring:

- scheduling and rescheduling
- cluster state
- stronger service discovery
- health-aware relocation
- more mature failover machinery

What they charge:

- a larger worldview
- more abstraction between operator and runtime
- more hidden machinery
- more distance from Compose-first legibility

These systems are not too big in the abstract.
They are only too expensive if they do not remove a repo-specific hidden burden
that smaller layers failed to remove honestly.

Private sentence still likely to survive unless proved otherwise:

> yes, but I still personally do not know whether the larger controller solved
> my wound or just made it harder to inspect.

## Promotion criteria by family

| Candidate family | It starts earning promotion when... | It has not earned promotion if... |
| --- | --- | --- |
| Compose plus shared registry | the wrong node can consult shared current placement truth and explain its choice | the registry is stale, operator-fed folklore, or not truly runtime-consumed |
| Gossip or event-driven coordination | events lead to inspectable and correct route decisions | the system knows who is alive but not who is truly eligible |
| Registry plus dynamic proxy | the proxy preserves route meaning and fallback survives backend loss | updates exist, but protected semantics or stateful caveats still rely on human memory |
| Stronger orchestrator | it removes a named hidden burden smaller layers failed to remove honestly | it mainly sounds more adult while leaving the same decisive truth privately held |

## Automatic disqualifiers for a not-yet-earned option

A candidate has not yet earned default promotion here if it mainly does one or
more of these:

- changes the control surface without naming which hidden burden moved
- improves deployment prestige while leaving wrong-node meaning weak
- improves controller power while leaving stateful authority socially held
- reduces local toil while preserving private topology reconstruction on the
  bad day
- sounds more grown up mainly because the worldview got larger
- offers recovery vocabulary without post-failure proof
- makes the diagram cleaner while keeping route semantics ambiguous

That does not make the candidate worthless.
It means the candidate has not yet answered the repo's real benchmark strongly
enough to deserve default status.

## What would actually count as orchestration progress

This repo should only narrate orchestration progress when at least one of these
becomes true in a provable way:

- the wrong node no longer needs private placement folklore
- peer choice no longer depends on remembered safety
- a fallback survives the failure that used to erase it
- a protected route preserves meaning after handoff
- a stateful surface now has explicit authority ownership and rediscovery
  semantics the runtime can explain

## What a real winner would look like

A winner in this repo would not merely be:

- more standard
- more popular
- more automatic
- more scalable on paper
- easier to compare in a blog post

A winner would make at least one previously true statement false:

- the operator is no longer the hidden placement ledger
- the wrong node no longer needs folklore to find the right backend
- peer choice no longer depends on remembered safety
- a fallback path no longer disappears during the exact failure it claims to
  solve
- protected-route meaning no longer becomes ambiguous after handoff

That is what winning has to mean here.

## The honest bottom line

The user is not starving for orchestrator names.
The user is starving for one option that can survive this accusation:

> if I stop privately finishing the topology sentence, does the platform still
> know what to do?

That is why this page has to stay harsher than a normal comparison page.
The right orchestrator, if there is one, will not merely be more capable.
It will be the first one that removes a specific humiliation instead of merely
describing it better.
