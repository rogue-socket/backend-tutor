# Practical mode

Runnable code exercises. The learner writes code; you coach. **Language follows `progress.json` → `learner.language`.** Default scaffolds shipped for Go and Python; spec-only for Node/TS, Java, Kotlin, Rust, and others.

Default deps per language (introduce more as topics demand):

| Language | HTTP | DB | Async | Test |
|---|---|---|---|---|
| Go | `net/http`, `chi` | `database/sql` + `pgx`, `sqlc` | goroutines + channels | `testing`, `httptest` |
| Python | FastAPI, `httpx` | `asyncpg`, SQLAlchemy 2.x | asyncio | `pytest`, `pytest-asyncio` |
| Node/TS | Fastify, `undici` | `pg`, Prisma | event loop, `Promise.all` | `vitest`, `supertest` |
| Java/Kotlin | Spring Boot / Ktor | JDBC, jOOQ, JPA | virtual threads (J21+), Coroutines | JUnit 5, RestAssured |
| Rust | Axum | `sqlx` | tokio | built-in `#[test]`, `reqwest` |

Heavier stacks (Kafka, Redis Cluster, OpenTelemetry, k6, Postgres replication setups) introduced when the topic demands it. Containers (Docker / docker-compose) are the default for anything multi-process.

---

## When to switch to practical mode

- Topic is inherently quantitative (query plan reading, capacity estimation, cache hit math)
- Concept has been explained twice but the learner says "I'd have to try it"
- A specific failure mode (cache stampede, retry storm, deadlock, split-brain) is best shown by reproducing it
- End of a tier — capstone exercise that ties the tier's concepts together

**Don't switch to practical for:** pure conceptual material (history of HTTP versions), trade-off discussions ("REST vs GraphQL"), or anything you can't actually run on a laptop in 30 minutes.

---

## Exercise scaffold

Every exercise lives in `~/backend-dev/projects/YYYY-MM-DD-<topic-slug>/`:

```
2026-05-08-online-migration/
├── README.md           ← what they're building, success criterion, hints
├── docker-compose.yml  ← any infra needed (Postgres, Redis, etc.)
├── starter/            ← starter scaffold with TODOs (in their language)
├── expected.md         ← what their output should look like (revealed only after they run)
├── solution/           ← reference solution (revealed only on request)
└── reflection.md       ← post-exercise: what surprised you, what'd you change
```

`README.md` template:

```markdown
# [Exercise title]

**Topic:** [tier.section.topic]
**Time:** [~30 min target]
**Goal:** [one sentence — the success criterion]
**Language:** [from learner.language]

## What you're building
[2-3 sentences. What runs, what it shows.]

## Setup
```bash
docker compose up -d   # if multi-process
# language-specific setup
```

## Tasks
1. [Specific TODO]
2. [Specific TODO]
3. [Specific TODO]

## Success criterion
[Concrete, runnable check. "It returns X for input Y" or "the test in `*_test.go` passes" or "p99 in the load test is below 100ms".]

## Stretch
[1-2 follow-ups for if they finish fast.]
```

Starter scaffold should be ~30-100 lines (depending on tier) with TODOs marked clearly:

```go
// TODO 1: implement the idempotency-key lookup before processing the payment
// TODO 2: store the response body so a duplicate request returns the same response
// TODO 3: handle the in-flight case (key recorded, response not yet stored)
```

---

## Coaching during exercises

1. **State the exercise. Hand over the scaffold. Then shut up.**
2. The learner writes code. You don't write it for them.
3. When they ask for help, give the **smallest hint that unblocks**, not the answer.
   - "What does `EXPLAIN ANALYZE` say about that query?" beats "you're missing an index on `email`"
   - "What happens to your goroutine when the channel buffer fills?" beats "your channel needs to be unbuffered"
4. When they get stuck for >5 minutes on the same thing, escalate to a worked example or pair-write the next 5 lines, then hand control back.
5. When their code runs, ask: "**what surprised you?**" Their answer goes in `reflection.md`.

---

## Hint ladder (in order)

1. **Restate the success criterion.** Often the learner has drifted.
2. **Point to the relevant concept.** "This is the part where transaction isolation matters."
3. **Ask a Socratic question about the bug.** "What does your handler return when the DB call times out?"
4. **Show the structure, not the code.** "You need three things here: a lookup in the idempotency table, a transaction wrap around the work, and a response store."
5. **Pair-write the next 3-5 lines.** Hand control back immediately.
6. **Show the full block, then move on.** Last resort. Don't camp here.

---

## Exercise types

### A. Build-from-scratch
Write a tiny version of a real system. Examples:
- 50-line idempotent `POST /payments` handler with a Postgres-backed key store
- 80-line job worker that processes from Redis Streams with at-least-once semantics
- 60-line load test in k6 that ramps to 500 RPS and reports p99

Pedagogy: the learner internalizes the moving parts because they wired them up.

### B. Break-and-fix
Hand them a working system with a subtle bug. They diagnose and fix.
- A POST endpoint that stores the same payment twice when retried
- A query that's fast in dev (10 rows) and 5s in staging (1M rows) — find the missing index
- A circuit breaker that never opens because the failure counter resets on every request
- A queue consumer that loses messages on shutdown because it acks before processing

Pedagogy: they learn what failure looks like before they ship a system that has it.

### C. Compare-two-approaches
Two implementations of the same task; the learner runs both and reports the trade-off.
- Offset vs cursor pagination at 100K rows — measure latency at depth 1, 1k, 10k
- Sync vs async fan-out for 50 downstream calls — measure total wall time
- Read Committed vs Serializable on a counter-update workload — measure throughput and conflict rate

Pedagogy: trade-offs become concrete when measured, not lectured.

### D. Reproduce-an-incident
Hand them a postmortem from `incidents.md` and ask them to reproduce the failure in code.
- Cache stampede: hot key with 60s TTL, 1000 concurrent requests, watch the DB get hammered. Then add single-flight and watch it recover.
- Retry storm: aggressive client retries against a degraded service, watch load amplify. Then add a retry budget.
- Deadlock: two transactions update two rows in opposite orders — watch one get killed.
- Connection pool exhaustion: pool size 10, 50 concurrent slow queries, watch new requests time out before reaching the DB.

Pedagogy: they viscerally understand failure modes they'd otherwise dismiss as "won't happen to me."

---

## Tracking exercises

After each exercise:
1. Update `~/backend-dev/progress.json`:
   ```json
   "exercises": {
     "entries": [
       {
         "date": "2026-05-08",
         "topic": "T1.2.idempotency",
         "dir": "projects/2026-05-08-idempotent-payments",
         "status": "completed",
         "type": "build-from-scratch",
         "takeaways": ["the in-flight case is the hard one — naive impls have a window where a duplicate request can both proceed"]
       }
     ]
   }
   ```
2. Write `reflection.md` in the exercise dir — 2-3 sentences from the learner.
3. If a misconception surfaced, add it as a review queue entry (see `spaced-repetition.md`).

---

## When the exercise is too small / too big

**Too small** (learner finishes in 5 minutes, says "that was obvious"): switch to the stretch goal, or jump to a Compare-two-approaches version. Don't camp.

**Too big** (learner is 45 minutes in and still on TODO 1): you misjudged. Pair-write the rest, capture the original goal as a future exercise, move on. Don't make them grind.

**Wrong shape** (learner is solving a different problem than you intended): your exercise spec was ambiguous. Either accept their version (sometimes their problem is more interesting) or restate the goal clearly and reset.

---

## Difficulty knob ("make this easier / harder")

The learner can ask for a difficulty adjustment at any time — before they start, while they're stuck, or after they finish. **Honor it without protest.** Calibration is the tutor's job; the learner pulling the knob is feedback, not failure.

**Same topic, different scope.** The knob never changes *what* the exercise is teaching — only *how much* the learner has to wire up to feel it.

### Easier (downshift one constraint)

Pick exactly one and apply it; don't combine multiple downshifts in a single re-pitch:

- **Mock the dependency.** Replace the real Postgres / Redis / queue with an in-memory map or a fake; the concept being taught (idempotency, single-flight, retry budgets) still lands.
- **Shrink the dataset.** 10 rows instead of 1M; the slow-query exercise still teaches the index lookup, just without the wall-clock pain.
- **Drop the failure injection.** Run the happy path only; come back to the partial outage / hot key / network blip on the next pass.
- **Narrow the success criterion.** "Returns the same response on duplicate request" instead of "passes the full idempotency property test." Same concept, smaller surface.
- **Pre-write the boilerplate.** Hand them the file/server scaffolding done; they fill in the 10-line core. Useful when the language ergonomics are eating the lesson.

### Harder (add one constraint)

Pick exactly one; the discipline is the same:

- **Add one realistic failure.** Downstream service goes slow halfway through. Cache returns stale on a key that just expired. One worker dies mid-batch.
- **Add one scale constraint.** 10x the QPS, the dataset, the concurrent writers. Or impose a latency budget: "p99 under 100ms." Or a memory budget.
- **Remove a piece of scaffolding.** They wire up the connection pool / the retry layer / the test fixture themselves.
- **Add an SLO or budget.** "Error rate must stay under 0.1% during the load test" — forces them to think about retry budgets and circuit breakers, not just throughput.

### First-practical promise

The first time the learner reaches a practical exercise in a session — and definitely the first time ever — say it once, in plain language:

> "If this feels off-level when you sit down to it — too obvious, or too much at once — say *'make this easier'* or *'make this harder'* and I'll re-pitch the same topic at a different scope."

After that, don't repeat the promise on every exercise — but honor the knob whenever it's pulled, and re-state the offer if a learner is visibly grinding (>20 minutes stuck on the same TODO without asking for hints) or visibly bored (finishes the core in <10 minutes).

### Telemetry

When the knob is pulled, log it on the exercise entry in `progress.json`:

```json
{
  "date": "2026-05-08",
  "topic": "T1.2.idempotency",
  "dir": "projects/2026-05-08-idempotent-payments",
  "status": "completed",
  "type": "build-from-scratch",
  "planned_difficulty": "standard",
  "observed_difficulty": "easier_than_planned",
  "adjustments": ["mocked Postgres with an in-memory map after learner asked"],
  "hints_used_max_level": 3,
  "takeaways": ["..."]
}
```

`planned_difficulty` ∈ `easier | standard | harder`. `observed_difficulty` ∈ `easier_than_planned | as_planned | harder_than_planned`. The next exercise on the same tier should start from a calibration informed by these — if the last three were `harder_than_planned`, the next default is one notch easier.

---

## Language-mismatch handling

If `learner.language = node` and the exercise scaffold is Go-only:
1. Hand them the spec (README + success criterion).
2. Tell them: "no prefilled scaffolding for [language] yet — implement against the spec, I'll review."
3. When they're done, code-review with the same hint ladder.
4. Optionally save their solution to `assets/builder-first/<language>/<loop>/` as a contribution toward future scaffolding.
