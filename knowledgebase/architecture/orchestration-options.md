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

That is the orchestration question that matters in this repo.

## The user's real decision, stated plainly

The user is not choosing between software logos.

They are trying to decide whether the dream can be preserved:

- ordinary Docker nodes
- Compose-first authorship
- visible, inspectable behavior
- fewer secret sacred nodes
- fewer secret sacred humans
- real fallback and failover behavior
- no premature surrender to a heavyweight control plane

The orchestration decision only becomes interesting when one of those dreams
starts colliding with truths the current stack cannot own by itself.

## The user's accusation against every option

Every candidate here should be read under the same accusation:

> are you actually removing the hidden human SPOF, or are you just replacing
> it with a more respectable story?

If the answer is the second one, the option may still be powerful.
It is not yet the right answer for this repo.

## What this page is and is not allowed to prove

This page is authoritative about:

- how orchestration options should be judged in this repo
- which burden a promoted layer would need to own to justify itself
- why product prestige is weaker than repo-specific burden transfer
- what kinds of candidate families are relevant to the actual problem
- what proof packets are required before stronger platform language becomes
  legal

This page is not authoritative about:

- which orchestrator has won globally
- whether the live runtime already demonstrates the promoted behavior
- whether ecosystem fashion should override repo-specific truth
- whether a more powerful controller is automatically a better answer

This is a promotion filter, not a winner declaration.

## The failure mode this page is trying to stop

The easiest way for a page like this to become dishonest is:

1. the option list gets more sophisticated
2. the controllers sound more grown up
3. the operator still privately carries the same decisive truth
4. the docs quietly call that progress

This page exists to stop prestige from masquerading as burden transfer.

## The current default stance

The current default stance in this repo is:

- keep Compose as the human control surface for as long as that remains honest
- refuse heavyweight promotion by prestige alone
- add the smallest extra truth-owning layer that actually removes a hidden
  burden
- promote to a stronger controller only when the smaller layer cannot carry the
  real failure classes anymore

That stance is not anti-orchestrator.

It is anti-unearned-orchestrator:

- anti "looks more adult"
- anti "more clustered must be better"
- anti "industry standard" as a substitute for proof
- anti replacing one hidden human dependency with one hidden controller myth

## What "worldview tax" really means here

`Worldview tax` is not just a dramatic phrase.
It means the option demands some combination of:

- more abstraction distance
- more invisible cluster state
- more controller trust
- more operational doctrine
- less direct reasoning from `docker-compose.yml` plus the runtime

That tax may be worth paying.
It is only worth paying when a real hidden burden moved in return.

## The non-negotiable questions every option must answer

Before any option gets described as the right next layer, it has to answer:

1. Which exact hidden operator burden does it remove?
2. Which truth does it own that the current runtime does not own cleanly?
3. How does the wrong node determine the right backend now?
4. How does the system distinguish a merely reachable peer from an eligible
   peer?
5. What evidence proves the route still exists after backend loss?
6. What evidence proves the route still means the same thing after handoff?
7. How does the option state its limits for TCP and stateful classes?
8. What worldview tax does the option impose in return?

If a candidate cannot answer those concretely, it may still be interesting, but
it has not earned default gravity here.

## The private completion test

Every orchestration option should be judged by a direct question:

> after adopting this layer, what exact sentence should the operator no longer
> need to finish privately?

Examples of the sentences this repo wants to kill:

- "That route only works because I know node3 is the real backend today."
- "The proxy can reach node2, but I privately know node2 is the only safe one."
- "The fallback config exists, but I know it disappears when the preferred
  backend dies."
- "Redis is reachable through the edge, but I know that does not mean there is
  safe failover."
- "The protected route still returns, but I know the auth semantics changed."

If the layer cannot kill at least one such sentence cleanly, it has not yet
justified its existence.

## The private-sentence benchmark

The useful benchmark is not:

- how many features does this system have?
- how standard is it?
- how complete is its cluster story?
- how easy is it to compare in generic infra discourse?

The useful benchmark is:

> what exact sentence should stop being true after we adopt it?

Examples:

- `I still personally know which node really hosts this service.`
- `I still personally know which peer is valid rather than merely alive.`
- `I still personally know the fallback disappears when the preferred backend dies.`
- `I still personally know the protected route changed meaning after handoff.`
- `I still personally know the stateful writer is singular.`

## The main candidate families

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
- the added layer only grows enough cluster truth to solve the wrong-node wound

What it can plausibly solve:

- current placement lookup
- inspectable local-vs-remote routing decisions
- a narrower, more honest burden transfer than a full scheduler

What it still must prove:

- registry freshness
- explicit disagreement handling
- peer eligibility truth
- fallback durability under actual backend loss
- protected-route continuity
- strict boundaries for TCP and stateful claims

This is the family most emotionally aligned with the repo.
It is also the family most vulnerable to fake closure through clever glue.

Private sentence still likely to survive unless proved otherwise:

> yes, but I still personally know whether the shared registry is current
> enough to trust

### 2. Compose plus gossip or event-driven coordination

Examples from the archive:

- Serf-style membership and failure events
- peer-equal node agents reacting to cluster events
- distributed signals without immediately adopting a manager hierarchy

Why it remains attractive:

- it respects the user's instinct against early hierarchy
- it preserves the emotional dream of equal nodes
- it can reduce private failure detection
- it can distribute signals faster than manual memory can

Why it remains risky:

- membership is not placement truth
- event distribution is not authority
- liveness is not eligibility
- gossip is much better at "who seems alive" than at "who is definitely the
  correct backend for this protected or stateful route"

This family is often the closest answer to:

> I do not want Swarm, but I also do not want to write an entire orchestrator

It still usually needs another truth-owning layer on top of the event flow.

Private sentence still likely to survive unless proved otherwise:

> yes, but I still personally know hearing the event is not the same thing as
> knowing the right backend

### 3. Compose plus registry and dynamic proxy/control-plane helpers

Examples:

- Consul-like registries
- HAProxy or Envoy driven from service-discovery state
- helper APIs publishing backend identity and policy inputs

What this family can improve:

- service-name to backend resolution
- explicit runtime-fed routing
- inspectable inputs for wrong-node forwarding
- better explanation surfaces for current choices

What it may still leave behind:

- stronger central trust concentration
- policy continuity gaps
- registry sacredness
- ambiguity about stateful authority

This family becomes a real candidate only if the receiving node can explain its
choice from shared truth rather than cultural memory.

Private sentence still likely to survive unless proved otherwise:

> yes, but I still personally know whether the registry or helper became the
> new sacred component

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
> my wound or just made it harder to inspect

## What the archive pressure means for these choices

The archive does not show random indecision.
It shows repeated pressure around the same fault line.

### `docker-multi-node-without-swarm__...`

This line of thinking accepts early that:

- manual placement may be acceptable
- Cloudflare multi-A DNS can help first-hop plurality
- plurality does not solve service ownership

The remaining hard problem becomes:

- mapping "service name" to "where is it really running right now?"

That means the orchestration question here is already narrower than "full
scheduler or not."
It is specifically about current placement truth and wrong-node rescue.

### `distributed-ha-orchestration__...`

This line keeps the peer-equal dream visible:

- all nodes equal
- any node can notice failure
- action should not depend on one sacred manager

The uncomfortable answer it keeps surfacing is:

- fully leaderless, drop-in orchestration for Compose is rare
- existing tools can get close
- glue is still usually required

That matters because it stops the docs from promising an off-the-shelf miracle
where the archive never actually found one.

## Why "just use the stronger orchestrator" is still too small

That answer is still too small because it skips the repo's actual decision:

- which truth is missing now?
- could a smaller layer own it honestly?
- if not, which larger layer owns it in a visibly better way?
- what exact sentence dies after the migration?

Without those answers, the recommendation is still mostly aesthetic.

### `nomad-multi-node-failover__...`

This pressure shows what stronger orchestrators can bring:

- cluster scheduling
- replication
- rescheduling
- service discovery

It also shows why "Nomad exists" is not automatically an answer:

- it changes the control-plane worldview
- it moves further from Compose-first authorship
- it still has to prove repo-specific burden transfer

### `load-balancer-failover-alternatives__...`

This pressure exposes another trap:

- a rich failover product can solve one edge slice well
- but still fail the broader service-ownership and route-meaning problem

That is why "better failover features" and "better orchestration answer" must
remain separate concepts in this repo.

## Promotion criteria by family

Use this table as the harsh filter.

| Candidate family | It starts earning promotion when... | It has not earned promotion if... |
| --- | --- | --- |
| Compose plus shared registry | the wrong node can consult shared current placement truth and explain its choice | the registry is stale, operator-fed folklore, or not truly runtime-consumed |
| Gossip or event-driven coordination | events lead to inspectable and correct route decisions | the system knows who is alive but not who is truly eligible |
| Registry plus dynamic proxy | the proxy preserves route meaning and fallback survives backend loss | updates exist, but protected semantics or stateful caveats still rely on human memory |
| Stronger orchestrator | it removes a named hidden burden smaller layers failed to remove honestly | it mostly sounds more adult while leaving the same decisive truth privately held |

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

That is what "winning" has to mean here.

## Bottom line

The repo does not need more orchestration names.
It needs one more option that remains believable after wrong-node entry,
backend loss, and hidden operator glue become real.

If an orchestration layer wins, it will not win by prestige.

It will win because one more previously private bad-day sentence stopped being
true.

## Bottom line

The user is not starving for orchestrator names.
The user is starving for one option that can survive the accusation:

> if I stop privately finishing the topology sentence, does the platform still
> know what to do?

That is why this page has to stay harsher than a normal comparison page.
The right orchestrator, if there is one, will not merely be more capable.
It will be the first one that removes a specific humiliation instead of merely
describing it better.
