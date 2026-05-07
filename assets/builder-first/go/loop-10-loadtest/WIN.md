# WIN — Loop 10 done

## Code

- [ ] `load.js` parameterised; runs against local or deployed instance
- [ ] Open-loop arrival rate (not closed-loop VU count)
- [ ] Realistic mix of read/write/auth-required endpoints
- [ ] k6 thresholds defined; failures fail the run

## Verification

- [ ] You ran the load test; *something* broke (or you pushed until it did)
- [ ] You identified the bottleneck from telemetry, not code reading
- [ ] You applied a fix (or a documented cap)
- [ ] You re-ran; quantified the improvement (before/after numbers)

## Postmortem

- [ ] Written up using `postmortem-template.md` (or your modification)
- [ ] Multiple contributing factors named, none of them "human error"
- [ ] Action items concrete and assigned (even if to yourself)
- [ ] "What didn't go well" honestly populated

## Understanding

1. **Open-loop vs closed-loop load testing — when does each lie?**
   *Outline: closed-loop (constant VU count) hides tail latency under back-pressure: when the service slows, VUs queue, fewer requests are issued, total RPS drops, and the dashboard shows lower load than reality. Open-loop (constant arrival rate) issues at the target rate regardless of latency, surfacing tail latency as it would be in production. Closed-loop is fine for "how much can a single user push," open-loop for "how does the service behave at this RPS." Most production load tests should be open-loop. Gil Tene's "How NOT to Measure Latency" talk is the canonical takedown of closed-loop measurement.*

2. **Little's Law: L = λW. What does it tell you for sizing a service?**
   *Outline: L = average concurrent requests in flight, λ = arrival rate (RPS), W = average request duration. So L = λ × W. If your service handles 500 RPS with 100ms average latency, L = 50 — you have 50 requests in flight on average. Bound L by your concurrency capacity (pool size, threads, goroutines). If your DB pool is 10 and your queries average 50ms, max RPS = 10 / 0.05 = 200, before queueing kicks in. Useful for back-of-envelope capacity questions.*

3. **A postmortem says "the root cause was that the engineer mistyped a parameter." What's wrong with that, and what's a better framing?**
   *Outline: 'human error' as root cause is almost always factually true and almost always useless — humans will keep making errors. The interesting question is what about the *system* let that error matter: was there no confirmation prompt? No automated check? No type system / linting / linter / code review catching it? A better framing: 'a typo in a CLI parameter went unverified because the production tool accepted any input without confirmation. The contributing factor is the tool's permissive interface, not the typo.' The action item then targets the system, not the person — which is a system you can improve.*

## Reflection

What surprised you? Common ones:
- The first bottleneck wasn't the one you'd have bet on
- Postmortem writing is where 80% of the learning is
- Telemetry from Loop 9 was the difference between "I see what happened" and "I'm guessing"

## What's next

You've completed all 10 loops. The links service has been through CRUD, persistence, migrations, auth, async work, caching, deploy, real-time, observability, and load testing. You've broken it intentionally in every loop and seen each failure recover.

What's next is your call:
- **Repeat with a different domain.** Pick a different small product idea, run loops 1–10 again. Most things you learned generalise; many don't until you re-derive them.
- **Specialise.** Take one tier (T2 databases? T11 distributed systems?) and go deep — `/start T2.8` and follow the curriculum.
- **Read incidents** in `references/incidents.md` and reproduce them. The reproduce-an-incident exercise type is among the most educational.
- **Ship something real.** The links service is small enough; deploy it, give it a domain name, post it somewhere, let real traffic hit it. Production teaches what staging can't.
- **Mock interview.** `/start mock-interview` — the tutor will run a senior-backend interview against your stated direction.
