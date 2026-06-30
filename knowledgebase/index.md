# bolabaden Infrastructure Knowledgebase

This site exists for one question only:

> how do you keep
> [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
> as the real human control surface, spread services across several ordinary
> Docker nodes, and make any healthy public node a believable first hop without
> secretly leaving one operator as the missing control plane?

That is the real subject of `bolabaden-infra`.

Everything else in this site is only useful if it helps answer that question
more honestly.

If a page answers a calmer neighboring question, the page is still failing even
when the page is technically correct.

## The thing earlier docs kept doing wrong

The earlier failure was not mainly:

- not enough words
- not enough pages
- not enough diagrams
- not enough technology comparisons

The real failure was:

1. the docs sounded organized
2. the docs sounded balanced
3. the docs quietly rewrote the user's demand into a smaller problem
4. the hidden operator burden became harder to see

That is how infrastructure documentation becomes polished and useless at the
same time.

The user is not mainly asking for:

- generic self-hosting guidance
- generic reverse-proxy patterns
- generic high availability vocabulary
- generic orchestrator comparison charts
- generic "mature homelab" structure

The user is asking for the exact place where multi-node Docker usually stops
being honest.

## The user's dream, stated without calming it down

The dream is not just:

- run Docker on more than one machine
- have more than one public IP or DNS record
- place services on different nodes
- avoid Kubernetes because it is annoying
- bolt on enough helpers that the stack sounds clustered

The dream is:

- keep Compose near the center of authorship
- keep the nodes ordinary enough that the stack is still understandable
- let Cloudflare point at more than one surviving public node
- preserve the meaning of the request when it lands on the wrong node
- preserve auth, middleware, and request semantics during peer forwarding
- stop treating placement truth as something one human just remembers
- stop replacing one sacred node with one different sacred node
- stop replacing one sacred node with one sacred human
- only pay heavyweight orchestration cost if that cost actually kills a real
  hidden burden

That is not a generic HA dream.
It is a very specific anti-fake-HA dream.

## The core accusation the site has to preserve

The accusation is:

> there are endless tools for multi-node Docker, overlays, service discovery,
> reverse proxies, failover, ingress, mesh, middlewares, and orchestration, but
> too many of them only make the stack look distributed while quietly leaving
> one person as the private source of truth when the bad day starts.

If a page forgets that accusation, it gets less useful even when it gets more
complete.

That accusation should stay active while reading every page:

> does this explanation still depend on a human privately knowing which node is
> special, which peer is safe, which backend is current, which fallback is
> real, or which answer is still authoritative under failure?

If the answer is yes, then the page is still describing an operator-owned truth
instead of a system-owned truth.

## The private-sentence benchmark

Use one brutal benchmark for every page:

> after reading this page, what exact private sentence is still being finished
> by the operator?

Examples:

- `I still personally know which node is the real one.`
- `I still personally know which peer can safely answer this request.`
- `I still personally know whether fallback still exists under actual failure.`
- `I still personally know whether auth and middleware survive peer forwarding.`
- `I still personally know whether the stateful answer is trustworthy or merely reachable.`

If a page cannot leave that sentence behind clearly, the page may still be
interesting while still being too soft for the real problem.

## The shortest honest architecture statement

The strongest intent surface in the repo is
[`.github/copilot-instructions.md`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/.github/copilot-instructions.md).
It says the intended direction is:

- no central orchestrator by default
- services manually assigned to nodes
- current-state truth preferred over scheduler-declared desired state
- local-first serving when the service already lives on the receiving node
- peer-forward fallback when the receiving node is healthy but the service is
  remote
- explicit separation between HTTP routing and raw TCP or stateful behavior
- anti-SPOF pressure without fake HA language

Its core request contract is:

```text
User -> Cloudflare DNS -> any surviving public node
  service is local  -> serve locally
  service is remote -> forward to healthy peer that really hosts it
```

That is the dream.
It is not the same thing as present-tense proof.

The gap between those two things is not documentation noise.
It is the whole architecture wound.

## What this repo already proves

The priority runtime is still rooted in:

- [`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml)
- active fragments under
  [`compose/`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose)

That runtime already proves a serious Compose-first platform with:

- a real Traefik-bearing edge and auth surface
- public and private network segmentation
- stateful services such as MongoDB, Redis, PostgreSQL, RabbitMQ, and Qdrant
- observability, probes, dashboards, and operator surfaces
- Headscale and mesh pressure
- alternate routing and egress experiments
- multiple domain clusters living under one Compose-first authorship surface

That runtime is not fake.

What it still does not prove is the exact thing the user is most angry about:

- that any healthy public node can accept the request and preserve its meaning
  when the service is remote
- that shared placement truth is explicit instead of privately remembered
- that peer eligibility is system-owned instead of guessed
- that fallback survives the failure that made fallback necessary
- that auth, middleware, and application semantics survive peer handoff
- that L4 or stateful surfaces escaped singular authority instead of merely
  becoming reachable through a serious-looking stack

The runtime is serious.
The missing burden transfer is serious too.

## The three truth layers that must stay visible

Every useful page in this site has to preserve all three of these at once:

1. the dream is specific
2. the runtime is real
3. the truth-owning middle layer is still incomplete

Lose the first one and the page becomes a calmer neighboring question.

Lose the second one and the page becomes architecture theater.

Lose the third one and the page starts sounding like the stack is mostly solved
because it looks sophisticated.

That third failure is the most common one in multi-node infrastructure docs.

## The actual wound under all the subtopics

Most of the subtopics in this repo are not independent subjects.
They are different places where the same wound shows up.

That wound is:

the operator is still acting like the missing control plane.

That hidden operator job currently includes things like:

- remembering what runs where right now
- remembering which public node is truly safe
- remembering which peer is truly safe
- remembering whether a generated fallback path still exists under failure
- remembering whether a route is merely syntactically rendered or actually
  semantically preserved
- remembering whether a stateful answer is authoritative or just reachable

Many pages sound mature while still leaving that human job intact.
This site is supposed to stop that from being easy to forget.

## What a fake good summary sounds like

These sentences sound good while still being too small:

- `the repo is exploring several HA strategies`
- `the repo is evolving toward a middle layer`
- `the repo has a strong multi-node foundation`
- `the repo is comparing lightweight versus heavyweight orchestration`
- `the repo has several ingress and fallback options`

Those statements are not always false.
They are just weak unless they also preserve:

> the operator is still the missing control plane on the bad day

If they do not preserve that sentence, they are the exact kind of summary the
user already hates.

## How this site should be used

This knowledgebase should be used less like a normal doc portal and more like a
strict retrieval system.

For any question, keep asking:

1. what exact question am I answering?
2. what smaller neighboring question would be easier to answer?
3. which truth layer is carrying this answer?
4. what stronger sentence is still forbidden?
5. what exact private sentence is still being carried by the operator?

If those answers are not visible, the reading pass is too vague.

## Start here

Use these pages first depending on what you need:

- [reading-paths.md](reading-paths.md)
  when you need the shortest honest retrieval route for a specific question
- [architecture/problem-and-goals.md](architecture/problem-and-goals.md)
  when you need the real target without pretending it is already proven
- [architecture/current-compose-runtime.md](architecture/current-compose-runtime.md)
  when you need the present implementation shape
- [architecture/request-path-and-failure-walkthrough.md](architecture/request-path-and-failure-walkthrough.md)
  when you need to understand why wrong-node behavior is the real benchmark
- [architecture/missing-middle-layer.md](architecture/missing-middle-layer.md)
  when you need the actual burden that helpers and side systems keep failing to
  kill
- [operations/devops-runbook.md](operations/devops-runbook.md)
  when you need proof classes instead of architecture poetry
- [operations/source-assimilation-index.md](operations/source-assimilation-index.md)
  when you need to know which source surfaces are shaping the narrative and how
  they should be separated

## The shortest honest summary of the whole site

`bolabaden-infra` is trying to answer a very specific question:

how far can a Compose-first, multi-node, anti-SPOF Docker platform go before it
either transfers truth ownership into a real middle layer or admits it still
depends on one human privately finishing the platform?

That is the real question.
This site only stays useful if it keeps answering that question instead of a
nicer one.
