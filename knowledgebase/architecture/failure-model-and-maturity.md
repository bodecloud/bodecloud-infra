# Failure Model and Maturity Matrix

This page is the anti-theater page for the priority implementation rooted at
[`docker-compose.yml`](/run/media/brunner56/MyBook/Workspaces/bolabaden-infra/docker-compose.yml).

It exists to answer the harder question behind almost every architecture
discussion in this repo:

> what actually breaks today, what truth is materially live, what is still
> carried by the operator, and what exact proof would be required before
> stronger maturity language stops being a lie?

That is the real maturity problem in `bolabaden-infra`.

It is not mainly a problem of "how far along is the stack?"
It is a problem of how often the ecosystem still offers something that sounds
like maturity while leaving the most important explanation trapped in private
operator memory.

That is why ordinary maturity language is not enough here.
Ordinary maturity pages assume the reader mainly wants orientation.
This reader wants to know which parts of the system have actually stopped
making them be the private adult in the room.

That is the hidden insult this page has to keep visible:

- the user is not short on technologies
- the user is not short on respectable diagrams
- the user is not short on architecture language
- the user is short on options that stop collapsing into "you still have to
  know what is really true when the bad day arrives"

The repo does not mature as one calm platform.
It matures as several uneven truth lanes:

- public entry
- local edge execution
- locality and placement truth
- peer eligibility
- route persistence
- protected-route continuity
- service-class-specific recovery
- stateful ownership and promotion

If those are compressed into one comfort level, the docs become flatter than
the worktree.
They also become kinder than the evidence.

## What this page is and is not allowed to prove

This page is authoritative about:

- which failure domains define the real maturity problem
- how unevenly those domains are currently maturing
- which hidden operator burdens are still part of the live failure model
- what the next honest maturity step is for each domain

This page is not authoritative about:

- whether one named route has already passed its harder drills
- whether a future control-plane direction has already earned promotion
- whether one partial gain upgrades the whole stack

## Strongest honest current answer

The strongest honest current answer is that the repo has outgrown "toy stack"
status but has not yet crossed into "request-preserving platform" status.
The current implementation proves a serious Compose-first infrastructure
surface, real ingress machinery, real auth and observability layers, and clear
architectural intent.
It does not yet prove that wrong-node traffic, backend loss,
protected-route parity, or stateful ownership have become system-owned truths
rather than better-documented operator burdens.

That distinction matters because many ecosystems would already start speaking
as if this were the final stretch:

- ingress exists
- Cloudflare is involved
- more than one node matters
- helpers are appearing
- therefore the platform is "basically" maturing into HA

This page exists to say no to that move.
The operator's real complaint is that "basically" keeps cashing out into a
moment where the system still needs private human explanation to stay honest.

That is why maturity language has to stay lane-specific here.
The repo is coherent enough that one strong-looking lane can start donating
borrowed confidence to its weaker neighbors if the docs get lazy.

That risk grows as the docs improve.
The site can now explain the uneven lanes well enough that a reader may start
feeling like the hard part is mostly classification and sequencing, when the
runtime may still be leaving the same decisive completion step in the
operator's head.

This page has to keep interrupting that feeling.

That is the part most infrastructure writing keeps sanitizing away.
The user is not merely annoyed by incompleteness.
They are angry at how often incompleteness arrives wearing the costume of a
nearly finished option.

## How to read maturity in this repo

A maturity label here does not mean "how many components exist."
It means:

- what the system itself currently owns
- what the operator still has to reconstruct
- what the docs must not quietly credit to the system yet

That makes the hidden metric very simple:

> how much human reconstruction is still required before the feature behaves
> honestly on the bad day?

If the answer is "a lot," the feature is still immature here even if the stack
already looks sophisticated.

Another way to say the same thing:
the lane matures only when one fewer private translation step is required from
the operator at the exact moment the nice diagram stops being enough.

That means "better understanding the lane" and "the lane matured" have to stay
separate achievements.
If the explanation got clearer while the private translation step stayed
intact, then the docs improved and the lane did not.

That is why this page should feel harsher than normal maturity language.
In many projects, maturity means polish, breadth, or adoption.
Here, maturity means the system has stopped humiliating the operator by
revealing too late that the operator was still the real keeper of placement,
eligibility, fallback, or stateful truth.

## Maturity levels used on this page

Use these labels narrowly.

### `Intent-shaped`

The repo has named the target clearly, but current runtime proof is weak or
absent.

### `Runtime-shaped`

Real components are live and the current stack visibly leans in that direction,
but the system still depends on hidden interpretation or unproved joins.

### `Partial proof`

Some narrower evidence exists or the failure shape is clearly isolated, but the
important next ceiling remains open.

### `Trustworthy for this lane`

Not a global victory. It means this exact lane has enough specific proof that
the docs can speak more strongly about it.

At the moment, most important lanes are still below that last level.

That is not failure of documentation.
It is the documentation refusing to pretend that naming the missing truths is
the same thing as system-owning them.

That refusal is part of the architecture discipline here.
If the docs become too eager to celebrate "directionally close" surfaces, they
start reproducing the exact thing the repo was created to push back on:
respectable prose that still leaves the operator holding the last real answer.

This page should therefore feel less like a platform scorecard and more like a
map of where the hidden operator control plane is still alive.

## What still does not count as maturity

This page needs a direct filter against maturity theater.

The following still do not count as a lane becoming materially more mature:

- more labels, helper names, or route objects exist
- more nodes can be named in the docs
- a route works on the preferred node
- a failure story sounds more plausible
- one lane got sharper while another lane was silently upgraded with it
- the reader now feels more oriented about why the lane is hard
- the docs can now narrate the failure boundary more elegantly

In this repo, maturity only moves when a hidden reconstruction burden actually
shrinks.

That phrase should be treated as literal gating logic, not just good rhetoric.
If the burden did not shrink, the lane did not mature.
It may have become more instrumented, more legible, or more platform-like.
Those are different claims and should stay different.

That difference is one of the hardest honesty rules in the whole site:

- more platform-like is not the same as more system-owned
- more system-owned is not the same as more trustworthy
- more trustworthy in one lane is not the same as global maturity

The repo has to keep repeating that because the surrounding ecosystem keeps
selling a softer story:

- a bigger stack feels more mature
- more nodes feel more mature
- a cleverer edge layer feels more mature
- a more famous platform feels more mature

All of that may be directionally useful.
None of it answers the user's actual benchmark unless one more explanation
stops needing to live in the operator's head.

## The private completion test for each maturity lane

Every lane on this page should also survive one direct test:

after calling the lane more mature, what exact sentence does the operator no
longer need to finish privately?

Examples:

- "yes, but I still know which node really has this route"
- "yes, but I still know which peer is actually safe for this handoff"
- "yes, but I still know the fallback route disappears if that backend dies"
- "yes, but I still know this stateful surface still has one sacred owner"
- "yes, but I still know which policy-bearing route is only cosmetically
  transferable and which one is actually semantically safe"

If the lane still leaves one of those sentences mostly intact, then the lane
may be better instrumented, better documented, or more impressive-looking, but
it is not yet mature in the sense this repo actually cares about.

That matters because ordinary maturity talk keeps rewarding exactly the wrong
signals here:

- more machinery
- more polish
- more vocabulary
- less honest pressure on the final private completion step

This page exists to keep that substitution illegal.

## What a lane-specific maturity proof packet would have to contain

Before any lane graduates to stronger wording, the docs should be able to
point to a packet that contains:

- the exact lane being upgraded
- the previous hidden operator burden
- the new system-owned truth or artifact
- the failure condition or comparison that was exercised
- the boundary sentence naming what the lane still does not prove

Without that packet, "more mature now" is still mostly atmosphere.

Atmosphere is not neutral here.
Atmosphere is how people get sold on options that still betray them under the
first serious routing or failover question.

That packet requirement is deliberately stricter than ordinary docs practice.
The user is already surrounded by explanations that sound plausible before the
failure drill and feel evasive after it.
This repo cannot afford to become another source of the same tone.

## The matrix

| Domain | What the current worktree materially proves | Current maturity | Hidden operator burden still present | Next honest maturity step |
| --- | --- | --- | --- | --- |
| Public node-entry reachability | The repo explicitly targets any-node entry; `cloudflare-ddns` is live in the edge stack; multi-node entry is a first-class idea in `.github/copilot-instructions.md` and `README.md` | Runtime-shaped | Public plurality is not the same thing as preserved request meaning; current DDNS behavior is still recorded as a real limitation in the master plan | Prove multiple public nodes can receive traffic without narrating that as end-to-end failover |
| Local edge execution | Traefik, TinyAuth, Nginx auth extensions, CrowdSec, file-provider config, and healthchecks are materially live | Runtime-shaped | Local edge health can still be mistaken for cross-node truth; auth and middleware continuity across peer handoff remain unproved | Prove one protected route locally with the exact policy surface visible |
| Placement and locality truth | The docs repeatedly converge on `services.yaml` as the needed current-state registry, but the priority runtime does not yet prove a live consumed root registry | Intent-shaped | The operator still remains the safest source of "what runs where right now?" | Introduce or expose one live placement-truth surface consumed by routing or eligibility logic |
| Peer eligibility truth | Headscale is materially live; peer communication is not purely hypothetical; the master plan names sync-agent and peer broadcast directions | Intent-shaped | Reachability can still be confused with correctness; the system does not yet visibly own "which peer is valid now?" | Prove one peer-selection decision is based on current tracked truth rather than folklore |
| Fallback-route persistence | `docker-gen-failover` exists, but the master plan explicitly records that the current model can delete routes when a container stops | Partial proof | The fallback surface is still vulnerable to disappearing at the exact moment it is needed | Replace or harden route generation, then exercise backend-loss drills |
| Protected-route semantic continuity | The edge stack already includes TinyAuth, Nginx auth middleware, CrowdSec, and Traefik middleware surfaces | Runtime-shaped | A forwarded request may not yet behave like the same service from the user's perspective | Compare local versus peer-forwarded behavior for one protected route and prove auth plus middleware parity |
| Stateless HTTP wrong-node success | The dream is explicit and the runtime is mature enough to make this a real near-term target, but no generic wrong-node proof is claimed | Intent-shaped | Wrong-node success still feels like a hope rather than a property | Prove one named stateless HTTP route through intentional wrong-node entry |
| Backend-loss fallback for HTTP | The docs already isolate this as distinct from wrong-node reachability and record the known route-persistence bug | Intent-shaped | A route that answers happily today may still disappear or semantically degrade during real loss | Stop the local backend and prove whether the peer-forward path survives with the same user-visible contract |
| TCP forwarding | The root graph already routes at least some TCP traffic such as MongoDB through Traefik TCP labels | Runtime-shaped | Transport reachability may be overread as service resilience | Define proof separately for transport success, client behavior, and ownership semantics |
| Headscale control plane | Headscale is materially live, public, and exposed through Traefik; the master plan is blunt that it remains effectively singleton today | Runtime-shaped | The mesh still depends on one active control-plane owner and SQLite-backed state | Prove leader election and data continuity only after a real HA path exists |
| Stateful databases and queues | Stateful surfaces exist in the root graph or surrounding plans, but the docs already reject branding them as HA by adjacency | Intent-shaped | Storage ownership, replication, promotion, reconnect semantics, and disk truth are still the real failure domain | Per service class, define write authority, replica semantics, failover sequence, and recovery proof |
| Convergence and drift control | The master plan clearly names secret sync, compose sync, and node-alignment pressure | Intent-shaped | A forwarded request may still land in semantically different runtime because node revisions or secrets differ | Expose a drift-detection surface and prove nodes agree on the inputs that matter for peer fallback |

## The lanes most likely to be overclaimed

Three lanes are especially vulnerable to documentation inflation.

### 1. Public entry

This is the easiest place to overclaim because multi-A records and multiple
reachable nodes feel emotionally satisfying. They are still only the first hop.

### 2. Traefik-centered ingress

Traefik is one of the strongest real components in the whole stack. That is
exactly why it is dangerous to overread it. A strong local routing surface can
make cross-node gaps feel smaller than they are.

### 3. Stateful services

These are where the ecosystem most often cheats. Reachability, TCP forwarding,
or even container restart behavior do not settle ownership, replication,
promotion, or reconnect truth.

There is a deeper reason these three lanes are so easy to lie about:

- they each produce a visible success signal early
- that visible signal feels like relief
- readers want relief
- then the hidden burden remains untouched underneath the relief

This page has to be strong enough to interrupt that emotional shortcut.

That shortcut matters because many surfaces in this repo already look mature
enough to a casual reader.
The risk is no longer that the system looks too primitive.
The risk is that it looks adult enough to be overcredited before the
wrong-node and backend-loss drills have actually transferred ownership of the
hard truth.

## What "more mature than before" is allowed to mean

The docs may say a lane has become more mature only if the relevant hidden
operator burden has actually shrunk.

Examples of honest maturity improvements:

- the system exposes current placement truth instead of requiring remembered
  placement
- the receiving node can choose a peer from shared current truth rather than
  from static folklore
- fallback routes survive the exact failure that used to delete them
- local and peer-forwarded protected behavior have been compared and shown to
  preserve the same visible contract

Examples of dishonest maturity upgrades:

- more components exist
- more labels exist
- more plans exist
- a route works locally
- a happy-path `200` happened

## The current failure model in plain English

Today the stack can still fail in ways that feel especially offensive to the
user:

- traffic may reach a healthy node without that node having trustworthy shared
  truth about where the service should really go
- a fallback-looking helper may still lose the route during the failure it was
  meant to absorb
- a protected route may still answer differently after handoff even if the
  response code stays green
- stateful surfaces may still be described more confidently than their storage
  and promotion model deserves
- the operator may still be the safest place where the important topology truth
  lives

That last point is the anti-SPOF reading that matters most here.

## What would materially change the maturity story

The smallest meaningful maturity jump is not "pick Kubernetes" or "add another
proxy." It is:

1. move current placement truth out of private memory
2. prove one real stateless wrong-node HTTP path
3. prove whether the same route survives backend loss
4. compare protected-route continuity across that handoff
5. only then promote stronger language for the HTTP lane
6. keep TCP and stateful lanes on their own harsher timelines

Until those things happen, the mature thing to say is not "the stack is mostly
resilient now."

The mature thing to say is:

> the stack is getting better at naming the right problems and has real
> machinery around them, but the most important truths are still maturing as
> separate lanes rather than one platform-wide victory.

That sentence is intentionally unsatisfying.
It should be unsatisfying.
The whole point of the repo is that the available options keep sounding more
satisfying than they deserve.

## The most dangerous false upgrade

The easiest false upgrade in this repo is:

> one HTTP lane got clearer, therefore the platform is now broadly HA-shaped

That sentence is still wrong here.

Even if one protected or stateless HTTP path becomes trustworthy, the docs must
still keep separate pressure on:

- TCP transport versus service authority
- backend-loss fallback versus wrong-node happy-path rescue
- stateful ownership and promotion
- singleton control-plane realities

If those boundaries blur, this page stops serving its only job.

The larger dream is not "eventually say HA with better caveats."
The larger dream is that the user eventually gets at least one path through
this system that feels like a real option rather than another intelligent
story about why the operator still has to carry the truth alone.
