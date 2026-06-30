# Operator Contract and Success Criteria

This page exists because the rest of the knowledgebase can be technically
accurate and still miss the real standard the user is forcing on the repo.

That standard is not:

> is there an interesting multi-node Compose architecture here?

It is:

> does this stack start feeling like one request-preserving personal cloud
> instead of several ordinary Docker nodes whose correctness still depends on
> operator memory, lucky request placement, and architecture theater?

If this page becomes soft, the whole knowledgebase becomes soft again.

For the literal request trace, read
[`request-path-and-failure-walkthrough.md`](request-path-and-failure-walkthrough.md)
next.
This page is the higher-level acceptance contract that says what has to become
true before the deeper frustration is honestly being solved.

## What this page is and is not allowed to prove

This page is allowed to:

- define the acceptance standard the rest of the knowledgebase must answer to
- keep the user's real benchmark sharper than ordinary HA language
- distinguish burden removal from architecture theater
- state what "solved" has to feel like from an operator standpoint

This page is not allowed to:

- claim that the current runtime already satisfies the contract
- use good acceptance language as substitute proof for good implementation
- collapse stateless and stateful success criteria into one maturity claim
- make the dream look cleaner by shrinking it into an easier neighboring goal

## Quick claim router

If the question is:

- "What is the real success contract for this repo?" this page is a primary
  answer.
- "Does the current stack already meet that contract?" no. This page defines
  the bar more than it proves the bar has been met.
- "Why do the docs keep rejecting prettier but smaller interpretations?" this
  page is one of the best answers.
- "Can I cite this as runtime proof?" only indirectly. Use the request,
  failure, and evidence pages for proof.

## What this page should make impossible to miss

This page exists to stop a specific failure mode in the documentation itself.

That failure mode looks like this:

- the repo accumulates better names for the same pain
- the diagrams become cleaner
- the options become more legible
- the docs sound more decisive
- the operator burden stays fundamentally where it was

If the documentation lets that happen, then the docs are participating in the
same problem the user is rebelling against.

The user is not asking for a more articulate explanation of why multi-node
Compose is hard.
The user is asking for the architecture to stop behaving like "several
ordinary Docker hosts plus operator folklore" at the exact moment where
request preservation matters.

So this page has one job:

> keep the difference between architectural fluency and actual burden removal
> brutally visible

## The dream in one sentence

The operator wants several ordinary Docker nodes to behave like one
request-preserving personal cloud, while keeping Compose as the readable
authoring surface and refusing fake HA, hidden sacred nodes, and heavyweight
orchestrator tax unless that tax clearly removes real pain.

That sentence is the center of gravity for the whole repo.

## What the user is actually tired of

The user is not fundamentally asking for:

- more modular Compose
- more healthchecks
- more DNS records
- more containers
- prettier diagrams

The recurring frustration is deeper:

- raw Compose is readable, but too static once wrong-node entry matters
- static upstream glue keeps pushing the burden back onto operator memory
- many "HA" answers stop at the first hop
- many tools solve deployment without solving request preservation
- many heavier orchestrators demand worldview surrender before proving they are
  paying down the right pain

This is why the repo keeps circling around:

- current-state placement truth
- local-first then peer-forward routing
- Cloudflare any-node entry
- Traefik plus stronger routing truth
- failover and sync agents
- OpenSVC, Nomad, k3s, and Kubernetes as options rather than conclusions

The real demand is not:

> make it clustered

It is:

> remove the hidden human SPOFs and request-path guesswork without replacing
> the whole operator surface before that replacement has clearly earned its
> keep

## What "solved" has to feel like

The dream is not satisfied because:

- a proxy runs
- multiple A records exist
- several nodes are online
- a route template can be generated

It is satisfied only if the platform starts to feel materially different at
request time and failure time.
The key phrase there is "feel materially different."
The user is not asking for a nicer explanation of the same hidden burden.
They are asking for the burden itself to stop living in operator memory.

## Strongest honest current answer

The repo already has a much sharper operator contract than most self-hosting
stacks ever write down, but that is still a contract, not a proof of
fulfillment. The strongest honest answer is that the user's real benchmark is
now explicit enough to police the rest of the docs, while the implementation is
still only partially able to satisfy it without hidden topology memory,
wrong-node ambiguity, or uneven service-class maturity.

The failure mode to watch for here is subtle:

- docs improve
- helper components grow
- routing stories become more sophisticated
- the operator still has to silently know which node is "really" the one that
  matters

That is not progress toward the dream.
That is a better narrated version of the same dependence on private
reconstruction.

### Feeling 1: any-node entry is normal, not a gamble

The operator should be able to think:

> it is okay if the request lands on node A, B, or C first

without secretly meaning:

> I just hope the right node happened to get chosen

Any-node entry only becomes real when the first hop stops feeling like a coin
flip.

### Feeling 2: local service stays genuinely local

If the requested service already runs on the receiving node, that node should
serve it there.

That matters because:

- the fast path stays fast
- locality stays legible
- debugging stays grounded
- the stack does not invent fake cluster ritual just to look distributed

### Feeling 3: wrong-node entry is survivable

If the request lands on the wrong node, the user does not want to hear:

- the DNS worked
- Traefik is healthy
- there is a peer on the mesh
- a route template exists

They want the request to still complete correctly because the receiving node
can determine:

- the service is not local
- which peer currently hosts it
- whether that peer is eligible now
- whether the route needed for fallback survives the relevant failure
- whether auth, middleware, headers, and policy remain coherent on that path

That is the real standard.

### Feeling 4: stateful systems stop being lied about

For Redis, MongoDB, Postgres, RabbitMQ, and other state-bearing systems, the
user wants the docs to stop cheating.

The system should never imply that a database is "HA enough" merely because:

- it is reachable through a global hostname
- a TCP router exists
- another node could theoretically be pointed at it
- the container could be recreated elsewhere

For stateful systems, the operator wants blunt answers:

- who owns writes
- who replicates from whom
- how promotion works
- what breaks on node loss
- what reconnect behavior clients must expect
- whether the real failure domain is still one local disk path

### Feeling 5: the control surface stays legible

The user is not refusing orchestration because they enjoy manual pain.

They are refusing unreadable control planes that hide the system's real
behavior behind controller magic, stale metadata, or abstraction tax.

The desired system should still let an operator answer:

- what runs where right now?
- who says that is true?
- how does a receiving node choose local serve versus peer handoff?
- what changes when a backend disappears?
- which added component is paying for itself, and how?

If the answer becomes "the controller knows," the repo has moved away from its
purpose unless that complexity was explicitly justified.
If the operator still has to mentally complete the architecture at the moment
of failure, the contract has not been met no matter how many layers exist.

## The operator question this repo keeps trying not to lose

The deeper demand underneath almost every routing, failover, or orchestrator
conversation is:

> if I disappear for a week, does the system still know enough about itself to
> preserve a request, or did I merely hide my own memory behind cleaner YAML?

That is harsher than normal infra acceptance language, but it is closer to the
real frustration.

The user is not merely asking whether the platform is automatable.
They are asking whether the truth needed on the bad day has moved out of one
person's reconstruction loop and into inspected, current, durable system
knowledge.

If that does not happen, then "options" are mostly theater:

- more nodes means more places to be wrong
- more proxies means more places to hide missing truth
- more orchestration means more places to rename the same SPOF

That is why this page keeps treating operator absence as part of the failure
model.
The hidden benchmark is not "could a strong operator recover this eventually?"
It is "does the system already know enough to stop needing hidden operator
completion as part of normal request handling?"

This is the anti-benchmark the repo has to keep in view.
If the architecture still degrades into:

- "the right person knows which box to hit"
- "the operator can explain the intended fallback"
- "the proxy could be rewired quickly if needed"
- "the labels more or less imply the answer"

then the central problem is still present.
Those are all variants of hidden human completion.

## What does not count as success

The repo has to keep refusing several fake versions of success.

### Fake success 1: multiple A records means failover

What it proves:

- more than one public IP may receive traffic

What it does not prove:

- the receiving node knows where to send the request next
- the right service is reachable from the wrong node
- policy continuity survives the handoff
- stateful semantics survive the handoff

### Fake success 2: the proxy is healthy

What it proves:

- the local edge process is up

What it does not prove:

- peer selection is correct
- fallback routes persist under backend loss
- auth or middleware remain coherent after peer handoff
- the route that exists in the happy path still exists in the failure path

### Fake success 3: the Compose stack is large and modular

What it proves:

- the repo is operationally serious
- a real root implementation exists

What it does not prove:

- anti-SPOF behavior
- wrong-node success
- stateful failover correctness

### Fake success 4: there is a service-registry idea

What it proves:

- the architecture recognizes placement truth as necessary

What it does not prove:

- the tracked runtime currently ships and consumes that truth

This is why the `services.yaml` gap matters so much.

### Fake success 5: it sounds like Kubernetes, but smaller

What it proves:

- a proposed system may be controller-shaped

What it does not prove:

- the added complexity is proportionate
- the user actually wants that worldview
- the missing middle layer has been solved rather than renamed

### Fake success 6: the service stayed reachable

What it proves:

- some path still answered

What it does not prove:

- service semantics were preserved
- the same auth, middleware, and headers applied
- the correct data authority answered
- the surviving path is trustworthy under stress

This is one of the easiest lies to tell accidentally.

### Fake success 7: the docs can fluently explain the architecture

What it proves:

- the repo has accumulated real ideas
- the system has named components
- someone can narrate a plausible flow

What it does not prove:

- the runtime has one inspectable truth source for placement
- the bad-day request path is actually owned by the system
- the explanation survives turnover, drift, or operator absence

This fake success matters here because documentation quality is part of the
problem statement.
The user is explicitly rebelling against smooth summaries that reduce the
feeling of ambiguity without reducing the ambiguity itself.

So the contract is not:

> explain the architecture well enough that the pain feels smaller

It is:

> explain it in a way that keeps the unresolved burden visible until the
> burden is actually removed

### Fake success 8: the burden moved from YAML confusion to agent or helper confusion

What it proves:

- the repo has introduced a smarter middle layer
- there is now a dedicated component with a strong name
- the architecture may be closer to current-state awareness

What it does not prove:

- the operator no longer has to trust invisible internal reasoning
- the helper's truth inputs are current, inspectable, and durable
- request preservation is owned by the system rather than merely rephrased

This matters because a "thin" helper can still become a new sacred box, a new
private memory store, or a new unverifiable controller if the repo stops
insisting on inspectable truth.

## The contracts the repo actually has to satisfy

The dream becomes much easier to evaluate when decomposed into contracts that
match real operator questions.

### Contract 1: node-entry contract

Claim:

> any healthy public node can be the first hop

Minimum strong evidence:

- more than one public node is intentionally eligible for entry
- the entry story is not secretly anchored to one real proxy box
- representative request paths can start on more than one node

What still does not satisfy it:

- multiple DNS records with no verified useful behavior after arrival

### Contract 2: locality contract

Claim:

> if the target service is already local, the receiving node serves it locally

Minimum strong evidence:

- one representative service path is demonstrably local-first
- the docs can explain why the node considered it local

What still does not satisfy it:

- merely observing that a local service responded once

### Contract 3: wrong-node preservation contract

Claim:

> if the request lands on the wrong node, that node can still preserve the
> request through a healthy peer without operator rescue

Minimum strong evidence:

- one real stateless HTTP path proves wrong-node success
- the receiving node demonstrably knew the service was remote
- peer choice was based on current enough truth
- the forwarded path preserved auth, middleware, and externally visible
  request meaning

What still does not satisfy it:

- the peer was reachable
- a proxy rule existed
- a manual retry succeeded after operator intervention

### Contract 4: backend-loss contract

Claim:

> when the preferred local backend disappears, the recovery path remains
> available and the request still succeeds through the same contract

Minimum strong evidence:

- route persistence survives backend stop/die conditions
- the recovery route does not vanish with the original container
- the same request still completes correctly through the peer path

What still does not satisfy it:

- a failover generator exists in theory
- recovery works only while the primary is still healthy
- the route survived only because the operator re-rendered or reattached it

### Contract 5: stateful honesty contract

Claim:

> state-bearing systems are described only as resilient when authority,
> replication, and client behavior are actually defined and tested

Minimum strong evidence:

- clear write authority
- clear replication or election model
- client discovery or reconnect story
- storage and promotion behavior under node loss

What still does not satisfy it:

- TCP exposure
- a stable hostname
- a container that can be restarted elsewhere

### Contract 6: readability contract

Claim:

> the operator can still inspect the system and understand where the truth
> lives without decoding a hidden control plane

Minimum strong evidence:

- placement, convergence, and routing truth surfaces are inspectable
- the docs can name who owns each important truth
- helper layers are narrow enough to audit instead of mythologize

What still does not satisfy it:

- "the automation handles it"
- "the controller knows"
- "it is generated somewhere"

### Contract 7: anti-folklore contract

Claim:

> the system's critical request-path truth is not primarily carried by shared
> operator memory

Minimum strong evidence:

- the docs can name the exact tracked or inspectable surface for placement
  truth
- the docs can name the exact tracked or inspectable surface for route
  eligibility truth
- a bad-day operator does not need to privately remember the "real node" to
  explain why recovery succeeded

What still does not satisfy it:

- the same two maintainers always know what to do
- the stack is survivable only after manual recollection
- the truth exists only as a social convention plus labels

These contracts are intentionally narrower than a grand "solved distributed
systems" claim.
But they are also stricter than a normal homelab success story.
The repo does not get credit for being thoughtfully modular if the wrong-node
path is still a semantic gamble.

## The repo’s real success sentence

If this repo ever honestly reaches the user’s dream, the defensible sentence
will be something like:

> multiple ordinary Docker nodes can receive the same public traffic surface,
> keep local services local, preserve wrong-node stateless HTTP requests
> through healthy peers without sacred-node memory, and speak much more
> honestly about where stateful truth still lives

That sentence is narrower than "zero SPOF everywhere."

It is also far more real.

## The single most important success test

If the user asks the rude but correct question:

> if a healthy request lands on the wrong healthy node, is that still a normal
> request or did we just enter a recovery ritual?

the finished system should answer:

> that is still a normal request

and the docs should be able to defend that answer without appealing to:

- maintainer intuition
- hoped-for peer reachability
- stale placement assumptions
- "it usually works"
- a hidden sacred front door

Until that answer is true, the documentation should keep sounding unfinished
where the system is unfinished.

## Bottom line

The operator contract here is not "be modern."

It is:

> stop requiring the operator’s private memory to complete the architecture at
> request time

That is the bar.

If a future change, document, helper, or promoted control layer does not move
the repo toward that bar, it is not solving the user’s real problem, no matter
how sophisticated it sounds.
