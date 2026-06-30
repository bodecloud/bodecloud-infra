# Operator Questions and Honest Answers

This page turns the recurring operator questions in the repo into burden-faithful
answers.

It does not exist to recommend products politely.
It exists to answer the sharper version of the question:

> which truth is still privately owned at the exact moment the platform is
> supposed to act coherent?

That is the standard that keeps this page useful.

## What this page is and is not allowed to prove

This page is allowed to:

- restate recurring questions in the sharper form the repo actually needs
- explain why common answers still stop one layer too early
- tie answers to current repo evidence
- say what artifact or drill would be required before a stronger answer
  becomes legal

This page is not allowed to:

- imply the runtime is nearly solved just because the questions are sharper
- replace burden transfer with better tooling advice
- blur first-hop plurality into request preservation
- use orchestration nouns as if they answered the benchmark by themselves

## Answer format used on this page

Each serious answer on this page should leave behind:

1. the hidden burden the operator is actually talking about
2. the strongest current evidence class behind the answer
3. why the nearby common answer is still too small
4. the next artifact or drill that would allow a stronger sentence
5. the private sentence the operator still has to finish today

Without those five pieces, the answer is still too close to shopping advice.

## Strongest honest current answer

The repo is not mainly suffering from lack of tools.
It is suffering from lack of options that relocate the right truths out of
operator memory without either:

- stopping one layer early
- or importing a much heavier worldview than has actually proved it pays for
  itself

That is why the ecosystem can feel full of products and still empty of real
answers.

## Question 1: What is the user actually trying to make true?

The user is trying to make several ordinary Docker nodes behave like one
request-preserving personal cloud without immediately paying full Swarm,
Kubernetes, or similar controller tax.

In repo-native terms, that means all of the following:

- any healthy public node can take the first hop
- a local service stays local when that is honest
- a wrong-node request still completes correctly
- the receiving node has current truth about placement and peer eligibility
- fallback survives real backend loss
- protected routes preserve auth and middleware meaning after handoff
- stateful services are treated far more harshly than stateless ones

Smaller goals may still be useful.
They are not the full ask.

## Question 2: Why do ordinary HA answers keep feeling fake here?

### Hidden burden

The operator still has to privately know one or more of:

- where the service really lives
- which peer is actually safe
- whether the rescue route survives backend loss
- whether a forwarded protected route still means the same thing
- whether a stateful surface only looks movable from the outside

### Strongest current evidence class

- `.github/copilot-instructions.md` states the real target contract
- root `docker-compose.yml` and active fragments prove a serious live edge and
  service surface
- knowledgebase/runtime pages now separate first hop, wrong-node, backend-loss,
  and stateful lanes explicitly

### Why the nearby common answer is still too small

Common answers such as:

- "point Cloudflare at more nodes"
- "just use Traefik"
- "add more healthchecks"
- "use service discovery"
- "use Kubernetes"

often solve one slice while leaving the decisive truth privately carried.

### What would allow a stronger answer

- one shared placement-truth surface consumed by routing or peer-selection logic
- one wrong-node stateless HTTP drill
- one backend-loss drill

### Private sentence still surviving today

> yes, but I still personally know which node really has it

## Question 3: Is Traefik the answer to the multi-node problem here?

### Hidden burden

The operator needs shared current truth, not only routing execution.

### Strongest current evidence class

Traefik is materially live through the priority runtime and already fronts real
HTTP and TCP surfaces.
It works with:

- `tinyauth`
- `nginx-traefik-extensions`
- `crowdsec`
- file-provider config
- TCP routers for services such as `mongodb` and `redis`

### Why the nearby common answer is still too small

Traefik buys:

- local routing execution
- TLS termination
- middleware execution
- auth integration
- HTTP and TCP exposure

It does not, by itself, buy:

- shared current placement truth
- shared peer eligibility truth
- route persistence under backend loss
- stateful authority semantics

### What would allow a stronger answer

- a receiving node using shared placement truth to choose local versus remote
- a protected route comparison between local and peer-forwarded execution

### Private sentence still surviving today

> yes, but I still personally know whether Traefik's next hop is actually the
> right peer

## Question 4: Why is Cloudflare not the answer by itself?

### Hidden burden

Plural first-hop reachability is not the same as preserved request meaning.

### Strongest current evidence class

- Cloudflare is part of the repo's explicit public-entry philosophy
- `cloudflare-ddns` is live in the edge stack

### Why the nearby common answer is still too small

Cloudflare can help:

- multiple public records
- first-hop resilience
- public exposure management

It cannot, by itself, tell the receiving node:

- whether the target service is local
- where it actually lives now
- which peer is eligible now
- whether the rescue path survives the relevant failure

### What would allow a stronger answer

- one wrong-node route proven after intentionally landing traffic on a non-owner
  node

### Private sentence still surviving today

> yes, but I still personally know that DNS redundancy did not solve the real
> request-preservation problem

## Question 5: Why does Headscale not solve service discovery by itself?

### Hidden burden

Reachability and identity are being mistaken for peer validity and current
placement truth.

### Strongest current evidence class

- Headscale is materially live through
  [`compose/docker-compose.headscale.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/compose/docker-compose.headscale.yml)
- the active config still uses SQLite at `/var/lib/headscale/db.sqlite`

### Why the nearby common answer is still too small

Headscale can help the repo with:

- private-node connectivity
- stable node identity
- mesh assumptions becoming real instead of imaginary

It does not by itself answer:

- which node currently owns the requested service
- which peer is eligible for this exact route
- whether the control plane itself has ceased to be singular

### What would allow a stronger answer

- a routing or eligibility surface that consumes peer identity plus current
  service ownership

### Private sentence still surviving today

> yes, but I still personally know that a reachable peer is not yet a proven
> valid backend

## Question 6: Why does `services.yaml` keep reappearing in the docs?

### Hidden burden

The operator keeps being the safest current-state registry.

### Strongest current evidence class

- repeated knowledgebase references to `services.yaml`
- `.github/copilot-instructions.md` explicitly names lightweight
  current-state-registry philosophy
- current runtime still does not prove a live tracked root `services.yaml`
  consumed by routing

### Why the nearby common answer is still too small

Without some shared placement-truth surface, every other answer stays weaker:

- Cloudflare only gets traffic onto a node
- Traefik only executes local routing decisions
- Headscale only makes peers reachable
- helpers only look plausible until the wrong-node question gets asked

### What would allow a stronger answer

- one tracked placement registry or equivalent truth surface, visibly consumed
  by some routing or selection logic in the priority runtime

### Private sentence still surviving today

> yes, but I still personally know what runs where right now better than the
> system does

## Question 7: Why are stateful services treated so much more harshly?

### Hidden burden

Write authority, replication, promotion, reconnect, and rediscovery truth are
still the real failure domains.

### Strongest current evidence class

The live runtime already contains:

- root `mongodb`
- root `redis`
- `headscale` SQLite
- Firecrawl `nuq-postgres`
- `rabbitmq`
- `litellm-postgres`

### Why the nearby common answer is still too small

The following do not answer the real question:

- stable hostnames
- TCP exposure
- successful healthchecks
- restartability
- "we can move it later"

Those can improve reachability and operations without changing who owns truth.

### What would allow a stronger answer

- per service class: explicit write owner, replica model, promotion flow,
  reconnect expectations, and rediscovery behavior

### Private sentence still surviving today

> yes, but I still personally know the real writer is singular

## Question 8: Why can a helper still be fake comfort?

### Hidden burden

A helper can reduce repetition without removing bad-day ambiguity.

### Strongest current evidence class

`docker-gen-failover` is materially present, and the repo already records that
it can delete routes when the backend stops.

### Why the nearby common answer is still too small

A helper often looks like progress because it:

- generates config
- reduces manual edits
- reacts to events

But if it fails during the exact failure it is meant to absorb, then it did
not move the real burden.

### What would allow a stronger answer

- backend-loss drill showing the rescue route still exists and still means the
  same thing

### Private sentence still surviving today

> yes, but I still personally know that the helper fails on the exact bad day
> I cared about

## Question 9: What is the most useful next proof to chase?

### Hidden burden

The repo still lacks one narrow, humiliatingly concrete proof that a receiving
node can act correctly without private operator completion.

### Strongest current evidence class

The runtime already has:

- good stateless HTTP candidates such as `whoami` and `wishlist`
- real edge policy surfaces
- real wrong-node architectural pressure

### Why the nearby common answer is still too small

"Pick a more mature platform" is too broad until the repo first proves what
burden is actually being moved.

### What would allow a stronger answer

The best next proof packet is:

1. expose one shared placement-truth surface
2. intentionally land a request on the wrong healthy node
3. prove one stateless HTTP route still completes correctly
4. then kill the preferred backend and prove whether the same route survives

### Private sentence still surviving today

> yes, but I still personally know the system has not yet passed the humiliating
> wrong-node test

## What still does not count as an honest answer

These still do not count:

- naming a more mature product category without naming which burden it would
  actually own
- answering "why are there no real options?" with a list of tools
- treating first-hop plurality as request preservation
- treating clearer prose as stronger evidence
- recommending a controller without also saying which private sentence it would
  kill

That is the whole protection mechanism for this page.

## What a genuinely useful answer should leave behind

After reading any serious answer in this repo, the operator should know:

- which truth is still private
- why the obvious nearby answer stops one layer too early
- what next artifact would externalize that truth
- what sentence remains forbidden until that artifact exists

If the answer leaves only a better shopping list, it failed.
