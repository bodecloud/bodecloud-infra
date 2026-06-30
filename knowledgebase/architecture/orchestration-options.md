# Orchestration Options

For the evidence underneath this page, start with:

- [Orchestrator Tradeoffs Evidence](../research/orchestrator-tradeoffs-evidence.md)
- [Orchestration Research 2026](../research/orchestration-research-2026.md)
- [`docs/INFRASTRUCTURE_MASTER_PLAN.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docs/INFRASTRUCTURE_MASTER_PLAN.md)

This page exists to stop one lazy question from taking over the repo:

> which orchestrator is best?

That is not the real decision surface here.

The real question is:

> which extra layer of control has actually earned the right to exist because
> it removes a real hidden human SPOF, a real wrong-node failure, or a real
> convergence failure that the current Compose-first stack cannot solve
> honestly?

That is the decision surface that matters.

## What this page is and is not allowed to prove

This page is authoritative about:

- how orchestration options should be judged in this repo
- which burden a promoted layer would need to own to justify its existence
- why respectable platform names are weaker than domain-specific truth
  ownership
- which candidate families are relevant to the repo's actual pressure

This page is not authoritative about:

- whether one orchestration path has already won globally
- whether the current runtime already demonstrates the promoted behavior
- whether broader ecosystem prestige should override repo-specific evidence

This page is a promotion filter, not a final platform verdict.

## The current default stance

The default stance in this repo is:

- keep Compose as the human control surface for as long as it remains honest
- refuse heavyweight promotion by prestige alone
- add the smallest extra truth-owning layer that actually removes a hidden
  burden
- promote to a stronger controller only when the smaller layer cannot carry the
  real failure classes anymore

That stance is not anti-orchestrator.
It is anti-unearned-orchestrator.

## The questions every option must answer

Before an option gets described as the right next layer, it has to answer:

1. Which exact hidden operator burden does it remove?
2. Which truth does it own that the current runtime does not own cleanly?
3. How does the wrong node know the right backend now?
4. How does the system distinguish a merely reachable peer from an eligible
   peer?
5. What post-failure evidence would prove the route still exists and still
   means the same thing?
6. What worldview tax does the option impose in return?

If a candidate cannot answer those concretely, it may still be interesting, but
it has not earned default gravity here.

## The main candidate families

### 1. Stay Compose-first and harden the shared-truth layer

Examples:

- tracked `services.yaml`
- node-local proxy plus shared registry
- peer-sync helpers
- mesh-assisted state broadcast

What this family is trying to preserve:

- `docker-compose.yml` remains the real control surface
- service placement stays human-legible
- the system grows just enough cluster truth to solve wrong-node routing

What it can plausibly solve:

- current placement lookup
- some inspectable local-vs-remote routing decisions
- narrower burden transfer without a full scheduler

What it still must prove:

- registry freshness
- peer eligibility truth
- durable fallback after backend loss
- stateful caveats

This is the family most aligned with the repo's emotional and operational
default.
It is also the family most vulnerable to fake closure through clever glue.

### 2. Compose plus gossip or event-driven coordination

Examples from the archive:

- Serf-style membership and failure events
- peer-equal node agents reacting to cluster events
- distributed signals without immediately adopting a manager hierarchy

Why it stays attractive:

- it respects the user's instinct against early hierarchy
- it can reduce private failure detection
- it keeps node equality more visible than server-client schedulers do

Why it stays risky:

- membership is not the same thing as authoritative placement truth
- event distribution is not the same thing as safe action ownership
- gossip is good at "who seems alive" and weaker at "who is definitely the
  right backend for this protected or stateful route"

This family is often the closest answer to "I do not want Swarm, but I also do
not want to build an entire orchestrator."
It still usually needs glue and proof discipline.

### 3. Compose plus service-discovery and control-plane helpers

Examples:

- Consul-like registries
- HAProxy or Envoy driven from service-discovery state
- helper APIs that publish backend identity to each node

What this family can improve:

- service-name to backend resolution
- dynamic routing input
- central or replicated registry semantics

What it changes emotionally:

- the user stops asking the wrong node to rely on folklore
- the system begins to own more current-state explanation

What it can still leave behind:

- central trust concentration
- policy continuity gaps
- stateful ambiguity

This family can be a real middle layer if the registry is trustworthy and the
receiving node can explain its choice without private memory.

### 4. Stronger orchestrators or cluster managers

Examples:

- Nomad
- OpenSVC
- k3s
- Kubernetes

What they can bring:

- scheduling and rescheduling
- native cluster state
- stronger service discovery
- health-aware relocation
- richer failover mechanics

What they charge:

- a larger worldview
- more abstraction between operator and runtime
- more hidden machinery behind the control surface

This cost is acceptable only if they remove a concrete burden the smaller
families cannot remove honestly.

## The archive pressure on these choices

The archive does not show random indecision.
It shows repeated pressure around the same exact fault line.

### `docker-multi-node-without-swarm__...`

This thread accepts two things early:

- manual placement is fine
- Cloudflare multi-A records solve first-hop plurality better than they solve
  service ownership

The remaining hard problem becomes:

- mapping "service name" to "where is it running right now"

That means the repo's orchestration question is already narrower than
"full scheduler or not."
It is specifically about current placement truth and wrong-node rescue.

### `distributed-ha-orchestration__...`

This thread keeps the peer-equal dream visible:

- all nodes equal
- any node can detect failure and take action
- decisions emerge from distributed agreement rather than one sacred leader

The answer it keeps surfacing is uncomfortable but useful:

- fully leaderless, drop-in orchestration for Compose is rare
- existing tools can get close
- glue is still usually required

That matters because it stops the docs from promising a neat off-the-shelf
miracle.

### `nomad-multi-node-failover__...`

This thread shows what stronger orchestrators can offer:

- cluster scheduling
- replication
- rescheduling
- service discovery

It also shows why "Nomad exists" is not by itself a repo answer:

- it changes the control-plane worldview
- it pushes the repo further from the readable Compose-first authoring model
- it still needs proof against the repo's exact burden-transfer benchmark

### `load-balancer-failover-alternatives__...`

This thread is useful because it exposes another trap:

- a rich failover product may solve one edge slice well
- but still not solve the whole route-meaning and service-ownership problem

That is why "better failover features" and "better orchestration answer" must
stay separate here.

## Promotion criteria by candidate family

Use this table as the harsh filter.

| Candidate family | It starts earning default gravity when... | It has not earned default gravity if... |
| --- | --- | --- |
| Compose plus registry | the wrong node can consult shared current placement truth and explain its choice | the registry is still folklore, stale, or only operator-maintained |
| Compose plus gossip/events | event signals lead to inspectable, correct route decisions without private memory | the system knows who is alive but not who is truly eligible |
| Registry plus dynamic proxy | the proxy preserves route meaning and fallback survives backend loss | routing updates exist but protected semantics or stateful caveats still rely on operator memory |
| Stronger orchestrator | it removes a named hidden burden the smaller layers failed to remove honestly | it mainly sounds more adult while leaving the same decisive truth privately held |

## Automatic disqualifiers for a not-yet-earned option

A candidate has not yet earned default promotion here if it mainly does one or
more of these:

- changes the control surface without naming which hidden burden moved
- improves deployment prestige while leaving wrong-node meaning weak
- improves controller power while leaving stateful authority socially held
- reduces local toil while preserving private topology reconstruction on the
  bad day
- sounds more adult mainly because the worldview got larger

That does not make the candidate worthless.
It means the candidate has not yet answered the repo's real benchmark strongly
enough to deserve default status.

## The private completion test

Every candidate should be judged by a direct question:

> after adopting this layer, what exact sentence should the operator no longer
> need to finish privately?

Examples of the sentence this repo wants to kill:

- "that route only works because I know node3 is the real backend today"
- "the proxy can reach node2, but I privately know node2 is the only safe one"
- "the fallback file exists, but I know it disappears when the preferred
  backend stops"
- "Redis is reachable through the edge, but I know that does not mean it is
  safe to call it failover"

If the layer cannot kill at least one such sentence cleanly, it has not
justified itself yet.

## Bottom line

The repo does not need more platform names.
It needs one more option that remains believable after wrong-node entry,
backend loss, and hidden operator glue become real.

That means the winning orchestration layer, if there is one, will not win by
prestige.
It will win by making one more previously private bad-day sentence stop being
true.
