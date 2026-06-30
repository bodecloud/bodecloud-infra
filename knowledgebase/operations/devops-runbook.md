# DevOps Runbook

This page is not here to make the operator sound competent.
It is here to stop the repo from confusing successful commands with transferred
burden.

The real operational question in `bolabaden-infra` is not:

> what commands can I run to get reassuring output?

It is:

> what exact evidence do I need before I can honestly say the system, rather
> than my own memory, preserved a request, survived a wrong-node landing, or
> reduced a real hidden SPOF?

That is a much harsher standard.
It is also the only standard that matches what the user is actually trying to
make true.

## What this page is and is not allowed to prove

This page is authoritative about:

- how to classify claims before testing them
- which evidence classes are weak, medium, and strong in this repository
- what order to inspect the runtime in
- what stronger sentence is still forbidden after each successful check

This page is not authoritative about:

- whether a specific route is already resilient
- whether a specific failover path already works
- whether a specific stateful service is already safe
- whether a future middle layer already earned promotion

This is a method page.
It is not a certificate page.

## The real output of a serious runbook pass

The real output of a pass in this repo is not:

- a green terminal
- a pile of commands
- a more confident paragraph
- a screenshot that one route answered

The real output is:

- one narrower honest sentence that became true
- one stronger sentence that is still forbidden
- one hidden operator burden that was either killed or exposed more clearly

If the pass does not identify those three things, the operator probably learned
something useful but did not yet reduce the social control plane.

## The real operator problem

The user is not mainly short on commands.
The user is short on options that remove private reconstruction burden.

The recurring hidden burden looks like this:

- hidden topology memory
- hidden placement truth
- hidden peer-eligibility truth
- hidden policy-preservation truth
- hidden state-authority truth

If a command succeeds but the operator still has to privately know:

- which node really hosts the service
- whether the current node is serving locally or forwarding
- whether the candidate peer is merely alive or actually valid
- whether the forwarded request still means the same thing
- whether the stateful writer is still singular

then the runbook has not yet reached the user's actual pain.

## The private sentence test

Every pass in this repo should be treated like a hunt for one surviving private
sentence.

Examples:

- `I still personally know which node really owns this service.`
- `I still personally know whether this node is serving locally or forwarding.`
- `I still personally know which peer is merely reachable versus eligible.`
- `I still personally know whether the forwarded request still means the same thing.`
- `I still personally know which writer is authoritative.`

The pass got stronger only if one of those sentences either:

- died completely
- shrank into a smaller surviving sentence
- or was exposed more honestly than before

If the sentence merely became harder to notice, that is infrastructure theater.

## Start every pass with one claim sentence

Before touching the runtime, write the claim in this format:

> I am trying to prove `<specific claim>` and I need `<proof class>` evidence.

Good examples:

- `I am trying to prove the merged root Compose graph still resolves, and I need authored-shape evidence.`
- `I am trying to prove wishlist.$DOMAIN answers through the current edge stack, and I need route-behavior evidence.`
- `I am trying to prove one request can land on the wrong healthy node and still preserve route meaning, and I need wrong-node drill evidence.`
- `I am trying to prove mongodb remained authoritative after a failure, and I need stateful-correctness evidence.`

Bad examples:

- `I want to make sure everything looks healthy.`
- `I want to test HA.`
- `I want to check failover in general.`

Those are not claims.
They are invitations to narrate comfort as progress.

## The evidence ladder

This repo needs a stricter evidence ladder than most homelab writeups.

| Evidence class | Typical tools | What it can honestly prove | What it still cannot prove |
| --- | --- | --- | --- |
| Authored shape | `docker-compose.yml`, `compose/`, `docker compose config` | the declared graph resolves and the priority implementation surface is inspectable | that requests survive pressure |
| Local runtime health | `docker compose ps`, healthchecks, container logs | a container is up on this node and may be healthy locally | wrong-node success, backend-loss survival, stateful correctness |
| Route behavior | `curl`, headers, Traefik logs, backend identity markers | one route answered and can sometimes be tied to one backend identity | that the same route survives the failure that makes fallback matter |
| Wrong-node drill | controlled node targeting plus route identity evidence | one named route preserved meaning after landing on a non-owner node | that the whole platform is now unified |
| Backend-loss drill | controlled failure plus before-and-after route evidence | one named route survived one named failure mode with known limits | that unrelated routes or stateful services inherited that property |
| Stateful correctness | leader, write-path, election, rediscovery, reconnect evidence | one named stateful surface preserved authority honestly | that the stateful layer as a whole is now anti-SPOF |

The operator should name the evidence class before running the first command.

## The burden map behind the evidence ladder

This is the more useful translation of the ladder:

| Evidence class | Hidden sentence it can sometimes kill | Hidden sentence it usually leaves alive |
| --- | --- | --- |
| Authored shape | `I do not even know whether this service is declared.` | `I still do not know whether the runtime preserves requests under pressure.` |
| Local runtime health | `I do not know whether this node is running the container at all.` | `I still do not know whether the node can rescue a wrong-node request.` |
| Route behavior | `I do not know whether this hostname currently answers.` | `I still do not know whether the answer survives wrong-node entry or backend loss.` |
| Wrong-node drill | `I do not know whether this route can preserve meaning from the wrong node.` | `I still do not know whether other routes, protected routes, or stateful routes inherited that property.` |
| Backend-loss drill | `I do not know whether this named failure still preserves this named route.` | `I still do not know whether authority, state, or unrelated paths inherited that property.` |
| Stateful correctness | `I do not know whether this stateful surface still knows who owns truth.` | `I still do not know whether the platform as a whole stopped being socially singular.` |

This is why a green result at the wrong rung often feels insulting rather than
helpful.

## The current runtime tells you what to inspect first

The strongest current runtime anchors are still:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active fragments under [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/)
- [`AGENTS.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/AGENTS.md)
- [`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md)

Concrete runtime facts worth remembering before any drill:

- active edge fragment: `docker-compose.coolify-proxy.yml`
- active mesh fragment: `docker-compose.headscale.yml`
- active metrics fragment: `docker-compose.metrics.yml`
- active egress fragment: `docker-compose.warp-nat-routing.yml`
- root-owned networks: `publicnet`, `backend`, `warp-nat-net`
- root-owned directly declared services include `mongodb`, `searxng`,
  `code-server`, `chat-analytics`, and protected admin surfaces
- tractable stateless proof candidates already exist: `whoami`, `wishlist`,
  `mkdocs`

Those facts matter because the runtime is broad enough that sloppy testing will
over-upgrade very easily.

## The route classes must stay separate during testing

Do not run one pass and narrate it as if all route classes benefited equally.

There are at least four materially different classes here:

1. stateless public HTTP routes
2. protected HTTP routes carrying auth and middleware meaning
3. raw TCP reachability surfaces
4. state-bearing services whose truth is not exhausted by successful transport

The same output can mean very different things across those classes.

For example:

- a `200 OK` on `whoami` is an excellent early stateless signal
- a `200 OK` on an authenticated admin surface is not enough unless auth and
  middleware continuity are preserved too
- a successful TCP connect is almost useless as an HA sentence by itself
- a stateful writer answering one more time can still be the wrong kind of
  success

If the class is not named, the result is already too blurry.

## Operational sequence

Use this order unless you have a very specific reason not to.

### 1. Inspect authored shape first

Goal:

> confirm what the priority implementation claims to be.

Typical checks:

```bash
docker compose config --quiet
docker compose config --services
docker compose config | rg "traefik|tinyauth|crowdsec|docker-gen-failover|headscale|mongodb|redis|warp"
```

Questions this stage can answer honestly:

- does the graph resolve?
- which fragments are active?
- is the named service or route even declared?
- which networks, configs, and secrets are present?

Forbidden upgrade after success:

> therefore failover works

Shape validation proves authored reality, not behavior under pressure.

Private sentence that usually survives:

> yes, but I still personally know whether the declared graph is more coherent
> than the runtime under failure

### 2. Inspect local runtime health second

Goal:

> confirm what this node is actually running and how healthy it looks locally.

Typical checks:

```bash
docker compose ps
docker inspect --format '{{.State.Health.Status}}' traefik
docker inspect --format '{{.State.Health.Status}}' searxng
docker logs --tail=200 traefik
docker logs --tail=200 tinyauth
docker logs --tail=200 crowdsec
docker logs --tail=200 docker-gen-failover
```

Questions this stage can answer honestly:

- is the local container actually running?
- is the local healthcheck green?
- is the local edge layer erroring?
- is the backend even present on this node?

Forbidden upgrade after success:

> therefore the route is resilient

Local health remains local evidence.

Private sentence that usually survives:

> yes, but I still personally know that a healthy local container says almost
> nothing about wrong-node rescue

### 3. Inspect route behavior third

Goal:

> prove one route answered and tie it to backend identity if possible.

Typical checks:

```bash
curl -I https://wishlist.$DOMAIN
curl -sv https://wishlist.$DOMAIN
curl -sv https://whoami.$DOMAIN
docker logs --tail=200 traefik
```

Prefer evidence that lets you answer:

- which hostname was hit?
- which router handled it?
- which backend or service did Traefik think it used?
- can the response be tied to node identity or app identity?

Forbidden upgrade after success:

> therefore wrong-node routing is solved

Happy-path route success is not a wrong-node proof.

Private sentence that usually survives:

> yes, but I still personally know this may only work because the request hit
> the lucky node on a good day

### 4. Only then run a wrong-node drill

Goal:

> prove one specific route preserves meaning when traffic lands on a healthy
> node that does not host the service locally.

Start with stateless HTTP candidates before protected or stateful routes.
The best current early candidates are usually:

- `whoami`
- `wishlist`
- `mkdocs`

Minimum proof packet:

- exact hostname tested
- receiving node identity
- actual backend node identity
- evidence of how the receiving node decided local versus remote
- evidence that the response still meant the same thing
- explicit sentence naming what was still not proven

If the drill still depends on private operator recollection of placement, say
that explicitly.
That means the hidden control plane is still human.

Forbidden upgrade after success:

> therefore any-node entry now works generically

One successful route is one successful route.

Private sentence that usually survives:

> yes, but I still personally know whether the rest of the platform would act
> this coherent without me

### 5. Run backend-loss drills separately

Goal:

> prove one named route behaves honestly after the expected backend goes away.

Minimum proof packet:

- before-and-after route behavior
- exact failure introduced
- whether recovery was local restart, peer forwarding, or operator
  intervention
- whether auth and middleware meaning stayed the same
- what still remained human knowledge

Forbidden upgrade after success:

> therefore the platform has failover

The real sentence always has to be narrower:

failover of what, under which exact failure, with which remaining human
burden?

Private sentence that usually survives:

> yes, but I still personally know how much of the rescue path was still me

### 6. Treat stateful drills as a separate discipline

Goal:

> prove one stateful service preserved authority, not just reachability.

Current repo examples that require especially harsh honesty:

- `mongodb`
- `redis`
- Headscale with SQLite
- `nuq-postgres`
- `litellm-postgres`
- `rabbitmq`

Minimum packet:

- authoritative writer or leader before failure
- what happened to authority after failure
- how clients rediscovered the correct topology
- what storage, replication, or election mechanism makes the claim honest
- whether the outcome was true continuity, manual intervention, or singular
  restart

Forbidden upgrades after success:

- `mongodb is HA because Traefik TCP exposed it`
- `redis is safe because the port answered`
- `headscale is redundant because more than one node can reach it`
- `stateful SPOF is solved`

Stateful claims stay the harshest claims in this repo.

Private sentence that usually survives:

> yes, but I still personally know who owns truth and how clients should
> recover

## Weak, medium, and strong evidence in practice

### Weak evidence

Examples:

- `docker compose config --quiet` passes
- the service appears in `docker compose ps`
- the healthcheck is green
- the hostname answers once

Weak evidence is still useful.
It proves the stack is not imaginary.
It does not prove the stack preserved meaning under pressure.

### Medium evidence

Examples:

- route behavior can be tied to a known backend identity
- Traefik logs confirm which router and service handled the request
- a backend restart or local outage was exercised and the route still answered
  with named limits

Medium evidence is where the repo starts becoming genuinely interesting.
It is still not broad anti-SPOF proof.

### Strong evidence

Examples:

- a named wrong-node drill with backend identity and policy continuity
  evidence
- a named backend-loss drill with explicit before-and-after semantics
- a stateful drill that proves authority transfer honestly, or proves honest
  singularity instead of pretending otherwise

Strong evidence is not green output.
It is a proof packet with a named ceiling.

## The fake-adult failure mode this page is trying to prevent

The most common documentation failure in infra repos like this is:

1. a route answers
2. the prose becomes more mature
3. the hidden human burden remains almost unchanged
4. the repo starts sounding solved

This page exists to break that pattern.

The repo should never upgrade because:

- the command list looked serious
- the stack name sounded enterprise
- the proxy logs were verbose
- the operator now feels better oriented

Orientation matters.
It is not the same thing as burden transfer.

## Example claim packets

These are the kinds of close-outs this repo actually needs.

### Example: authored shape

- `Claim tested:` the current root graph still includes the edge, mesh,
  metrics, docs, and WARP fragments.
- `Evidence class:` authored shape.
- `What this proves:` the priority implementation surface still materially
  contains those layers.
- `What is still forbidden:` saying those layers already cooperate into
  generic wrong-node success.

### Example: stateless route

- `Claim tested:` `wishlist.$DOMAIN` answers through the current Traefik
  stack.
- `Evidence class:` route behavior.
- `What this proves:` one public HTTP route answered and can be inspected at
  the edge.
- `What is still forbidden:` saying the same route would survive wrong-node
  entry or backend loss.

### Example: wrong-node drill

- `Claim tested:` one request for `whoami.$DOMAIN` landed on a non-owner node
  and still reached the correct backend.
- `Evidence class:` wrong-node drill.
- `What this proves:` one stateless route preserved meaning under one named
  topology condition.
- `What is still forbidden:` saying protected routes, TCP routes, or stateful
  routes inherited that property.

### Example: stateful honesty

- `Claim tested:` MongoDB still answered after one failure.
- `Evidence class:` local runtime health or route behavior, unless authority
  transfer was explicitly shown.
- `What this proves:` one node-local MongoDB surface remained reachable or
  restarted.
- `What is still forbidden:` saying MongoDB authority, election, or rediscovery
  became multi-node safe.

## What still does not count as a serious runbook result

These are still invalid outcomes:

- a green command with no named claim
- a route test with no backend identity
- a failover claim with no explicit failure introduced
- a stateful reassurance story that never names authority
- an operator summary that hides how much topology truth was still private

Those are exactly the outcomes that make the docs sound more adult while the
platform remains socially manual.

## What the user is actually asking the runbook to do

The user is not begging for more checklists.
The user is asking for a runbook that makes several ordinary Docker nodes feel
less humiliatingly dependent on private operator memory.

That means the runbook is only good if, after using it, the operator can say:

- `the system now owns this smaller piece of truth`
- `this stronger sentence is still forbidden`
- `this is the next exact burden to externalize`

Anything softer than that may still be useful operations writing.
It is not yet aligned with the user's actual dream.

## Required close-out after every operational pass

Every pass should end with these four sentences written explicitly:

1. `Claim tested:` what exact claim was tested?
2. `Evidence class:` what class of evidence was gathered?
3. `What this proves:` what narrower sentence is now honest?
4. `What is still forbidden:` what stronger sentence would still be a lie?

If those four lines are missing, story inflation will usually begin
immediately.

## Bottom line

The test for this runbook is not:

> did it make the operator feel informed?

The test is:

> after following it, what exact private topology sentence did the operator no
> longer have to finish alone?

If the answer is:

> none yet, but the runtime is better understood

that is still honest progress.

It is not yet the user's dream.
