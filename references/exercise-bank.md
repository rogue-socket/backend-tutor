# Exercise bank

Catalog of exercises by tier. Use this as a menu — pick the exercise that matches the topic, the learner's language, and the time budget. Each entry: name, type, time, language fit, goal, success criterion.

**Types** (from `practical-mode.md`):
- **A. Build-from-scratch** — write a tiny version of a real system
- **B. Break-and-fix** — diagnose a planted bug
- **C. Compare-two-approaches** — measure a trade-off
- **D. Reproduce-an-incident** — recreate a documented failure

**Language fit:** "All" if the exercise is portable; specific list otherwise.

**Time:** target wall-clock for the build itself, not including theory or writeup.

---

## T0 — Networking & HTTP

### T0-A1: HTTP-by-hand
- **Type:** A
- **Time:** 30 min
- **Language fit:** All
- **Goal:** open a raw TCP connection to a real HTTP server, send a `GET /` request, parse the response without any HTTP library.
- **Success criterion:** prints status line + headers + first 200 bytes of body. Bonus: handle chunked transfer encoding.

### T0-A2: TLS handshake observer
- **Type:** A
- **Time:** 45 min
- **Language fit:** All (with `openssl s_client` as the easy mode)
- **Goal:** capture a TLS 1.3 handshake with `tcpdump` or `wireshark`, walk through the messages, identify SNI in the ClientHello.
- **Success criterion:** annotated capture file or screenshot showing each handshake message with a 1-line description.

### T0-A3: WebSocket echo + backpressure
- **Type:** A
- **Time:** 60 min
- **Language fit:** Go, Python (FastAPI + websockets), Node
- **Goal:** WebSocket echo server. Then add a slow consumer scenario — server sends 10K messages, client reads slowly, observe what happens to memory.
- **Success criterion:** demonstrate the unbounded-buffer growth, then fix it with explicit backpressure (drop, throttle, or close).

### T0-D1: Reproduce the Cloudflare 2019 regex outage
- **Type:** D
- **Time:** 30 min
- **Language fit:** All
- **Goal:** write a small server that runs a vulnerable regex on user input. Hit it with a crafted payload and watch CPU pin to 100%.
- **Success criterion:** profiler shows time spent in regex backtracking; mitigation (regex timeout, RE2-style non-backtracking engine) reduces it.

---

## T1 — APIs

### T1-A1: Idempotent payments endpoint
- **Type:** A
- **Time:** 60 min
- **Language fit:** Go (scaffold), Python (scaffold), others spec-only
- **Goal:** `POST /payments` with idempotency keys backed by Postgres. Must handle: first request, duplicate request, in-flight duplicate (key recorded but response not yet stored), key collision under concurrency.
- **Success criterion:** integration test with 100 concurrent duplicate requests produces exactly one charge and the same response body to all callers.

### T1-B1: Find the broken pagination
- **Type:** B
- **Time:** 30 min
- **Language fit:** Go, Python
- **Goal:** given a server using offset/limit pagination over a 1M-row table, identify why deep pages take 5+ seconds. Convert to cursor pagination.
- **Success criterion:** p99 latency under 50ms at any depth; explain why offset is O(n) at deep pages.

### T1-C1: REST vs gRPC same service
- **Type:** C
- **Time:** 90 min
- **Language fit:** Go, Python
- **Goal:** implement the same simple service (e.g., `getUser`) twice — REST + JSON, gRPC + protobuf. Run a load test against both.
- **Success criterion:** measured p50/p99 + bytes-on-wire numbers, plus a written 3-sentence trade-off.

### T1-A2: OAuth 2.0 authorization code + PKCE flow
- **Type:** A
- **Time:** 90 min
- **Language fit:** All
- **Goal:** implement a public-client (SPA-style) auth flow against a real provider (Auth0, Okta dev, or self-hosted Keycloak). Walk through redirect, code exchange, PKCE verifier.
- **Success criterion:** working flow that retrieves an access token; explain why PKCE is mandatory for public clients.

### T1-A3: Webhook receiver with signature verification + idempotent handling
- **Type:** A
- **Time:** 60 min
- **Language fit:** All
- **Goal:** receiver for Stripe-shaped webhooks. Verify HMAC signature with constant-time compare. Handle duplicates (Stripe retries up to 3 days). Reject replays >5 min old.
- **Success criterion:** unit test for signature mismatch, duplicate event, replay attack, valid event.

---

## T2 — Databases

### T2-A1: Schema design + first migration
- **Type:** A
- **Time:** 60 min
- **Language fit:** Any with a Postgres client
- **Goal:** design schema for a multi-user todo app (users, lists, items, sharing). Write the migration in Postgres. Add appropriate indexes.
- **Success criterion:** schema review against a checklist (PKs, FKs, NOT NULL where appropriate, indexes for likely queries, no JSONB blobs of stuff that should be columns).

### T2-A2: EXPLAIN ANALYZE deep dive
- **Type:** A
- **Time:** 45 min
- **Language fit:** All (Postgres CLI)
- **Goal:** seed a Postgres DB with 1M rows. Run a slow query, read `EXPLAIN ANALYZE`, identify the bottleneck (seq scan, missing index, bad join order), fix.
- **Success criterion:** query down from >1s to <50ms with the right index; explain the plan in plain English.

### T2-B1: Find the N+1
- **Type:** B
- **Time:** 30 min
- **Language fit:** Go (sqlc / pgx), Python (SQLAlchemy)
- **Goal:** given a list endpoint that does N+1 queries, find it via query log, fix with a single JOIN or `IN (...)` batch.
- **Success criterion:** query count goes from N+1 to 1 (or 2); response time falls correspondingly.

### T2-A3: Online migration — add NOT NULL column
- **Type:** A
- **Time:** 90 min
- **Language fit:** All (Postgres)
- **Goal:** add a `NOT NULL` column to a 1M-row table without locking. Use the expand/contract pattern.
- **Success criterion:** zero downtime measured by parallel-running load tester (no failed inserts during migration).

### T2-D1: Cache stampede reproduction
- **Type:** D
- **Time:** 60 min
- **Language fit:** Go, Python (with Redis)
- **Goal:** hot key with 60s TTL, 1000 concurrent requests at expiry. Watch the DB get hammered. Add single-flight; watch it recover.
- **Success criterion:** with naive cache-aside, DB QPS spikes to 1000+ at expiry; with single-flight, DB QPS stays at 1.

### T2-A4: Postgres replication failover drill
- **Type:** A
- **Time:** 90 min
- **Language fit:** All (Postgres + docker-compose)
- **Goal:** set up primary + async replica with `pg_basebackup`. Run continuous writes. Kill the primary. Promote the replica. Observe what gets lost (replication lag).
- **Success criterion:** numerical answer to "how many writes did we lose?" — and explain why the answer depends on `synchronous_commit`.

---

## T3 — Concurrency & async

### T3-B1: The counter race
- **Type:** B
- **Time:** 20 min
- **Language fit:** All
- **Goal:** given a `counter++` handler with 100 concurrent decrements, observe missing decrements. Fix with mutex, then with atomic, then with a DB row-level lock. Compare throughput.
- **Success criterion:** all three approaches produce correct totals; learner can describe when each is appropriate.

### T3-A1: Idempotent Redis Streams consumer
- **Type:** A
- **Time:** 90 min
- **Language fit:** Go, Python
- **Goal:** consumer group reading from a Redis stream. Process each message at-least-once with idempotent dedup (e.g., processed-ID set with TTL). Survive consumer crash mid-processing.
- **Success criterion:** kill -9 the consumer mid-batch; on restart, no duplicate side effects, no lost messages.

### T3-D1: Retry storm
- **Type:** D
- **Time:** 45 min
- **Language fit:** All
- **Goal:** client retries every 5xx with no backoff. Server is degraded (50% errors, 1s latency). Observe load amplification (10x). Add exponential backoff + jitter. Add a retry budget.
- **Success criterion:** with naive retries, server load is 10x baseline; with budget, load is bounded; explain why jitter matters.

### T3-A2: Saga with compensation
- **Type:** A
- **Time:** 120 min
- **Language fit:** Go, Python
- **Goal:** orchestrate a 3-step saga (charge → reserve → ship). When step 2 fails, compensate step 1. Idempotent compensations.
- **Success criterion:** every step is idempotent; injected failure at any point leaves the system in a consistent state (no orphaned charges or reservations).

### T3-D2: Deadlock drill
- **Type:** D
- **Time:** 30 min
- **Language fit:** All (DB)
- **Goal:** two transactions update rows A and B in opposite orders. Watch one get killed by the deadlock detector.
- **Success criterion:** reproduce reliably; fix by lock-ordering; explain why retry alone isn't sufficient (livelock).

---

## T4 — Caching

### T4-A1: Cache-aside vs write-through
- **Type:** C
- **Time:** 60 min
- **Language fit:** Go, Python
- **Goal:** implement both for a `getUser` endpoint backed by Redis + Postgres. Measure read latency, write latency, staleness window after a write.
- **Success criterion:** measured numbers + a written 3-sentence trade-off.

### T4-D1: The expiring hot key
- **Type:** D
- **Time:** 45 min
- **Language fit:** Go, Python
- **Goal:** see T2-D1 (same exercise; appears in both tier banks because it teaches both DB and cache topics).

### T4-A2: Probabilistic early refresh
- **Type:** A
- **Time:** 60 min
- **Language fit:** Go, Python
- **Goal:** implement XFetch (Vattani et al.) for a hot cache key. Compare hit-rate / staleness vs naive TTL.
- **Success criterion:** demonstrate that early-refresh smooths the spike at expiry; learner can explain the math.

---

## T5 — Reliability

### T5-A1: Circuit breaker from scratch
- **Type:** A
- **Time:** 60 min
- **Language fit:** Go (no library), Python
- **Goal:** implement a circuit breaker (closed / open / half-open) without using a library. Wrap a flaky downstream. Tune thresholds.
- **Success criterion:** under sustained downstream failure, breaker opens within N requests; in half-open, allows exactly one probe; explain why "fail fast" is a feature.

### T5-A2: Timeout budget across 3 hops
- **Type:** A
- **Time:** 45 min
- **Language fit:** All
- **Goal:** caller → service A → service B → service C. Each hop has its own timeout. Design budgets so the caller's deadline propagates correctly (gRPC deadlines or X-Request-Deadline header).
- **Success criterion:** when service C is slow, the caller times out at the right time; service A doesn't keep working on a request the caller has abandoned.

### T5-D1: Health-check cascade
- **Type:** D
- **Time:** 30 min
- **Language fit:** All
- **Goal:** readiness probe checks DB connectivity. DB has a 5-second hiccup. All replicas fail readiness simultaneously, the entire service is pulled from LB, hiccup turns into outage.
- **Success criterion:** reproduce; fix (probe a local cache of the last-good DB ping; or don't check DB in readiness at all — that's what monitors are for).

### T5-A3: Graceful shutdown with in-flight requests
- **Type:** A
- **Time:** 45 min
- **Language fit:** All
- **Goal:** SIGTERM handler that drains in-flight requests, closes DB pool, exits within 30s.
- **Success criterion:** load tester running while you SIGTERM the server; zero failed requests; clean exit.

---

## T6 — Observability & on-call

### T6-A1: Structured logs + correlation IDs
- **Type:** A
- **Time:** 45 min
- **Language fit:** All
- **Goal:** middleware that injects a request ID, propagates it to downstream calls, includes it in every log line. Trace a request across two services.
- **Success criterion:** given a request ID, can grep both services' logs and see the full path.

### T6-A2: OpenTelemetry tracing across two services
- **Type:** A
- **Time:** 90 min
- **Language fit:** Go, Python
- **Goal:** instrument two services with OTel. Run them locally with Jaeger or Tempo. Trace a request that hops both. Add span attributes that would help debug a slowdown.
- **Success criterion:** Jaeger UI shows a multi-service trace with attributes; learner can identify the slowest span.

### T6-A3: SLO + multi-window burn-rate alerts
- **Type:** A
- **Time:** 60 min
- **Language fit:** All (Prometheus + Grafana)
- **Goal:** define a 99.9% availability SLO for a service. Configure Prometheus alerts for fast-burn (1-hour window) and slow-burn (6-hour window). Test by injecting failures.
- **Success criterion:** alerts fire at the right thresholds; explain why one-window alerts page too late or too noisily.

### T6-D1: Page-worthy vs ticket-worthy
- **Type:** B-style (paper exercise)
- **Time:** 30 min
- **Language fit:** N/A
- **Goal:** given 10 example alerts, classify each as page / ticket / delete. Justify each.
- **Success criterion:** matches the rubric (symptom + actionable + urgent = page; everything else lower); explain the costs of each misclassification.

---

## T7 — Performance & scale

### T7-A1: Load test with k6 + identify the bottleneck
- **Type:** A
- **Time:** 90 min
- **Language fit:** Service in Go/Python; load tester in k6 (JS)
- **Goal:** ramp from 0 to 1000 RPS against a service with a planted bottleneck (small connection pool, missing index, sync I/O on hot path). Find it via metrics + profile.
- **Success criterion:** identify the bottleneck before fixing it; measure the improvement.

### T7-A2: Capacity estimation paper exercise
- **Type:** A
- **Time:** 30 min
- **Language fit:** N/A
- **Goal:** "10M users, 200 events/day, 1KB each, 1-year retention." Calculate disk, peak QPS (4x avg), and a sensible Postgres + S3 split. Show your math.
- **Success criterion:** numerical answers + a sentence on the assumption set.

### T7-B1: The slow query
- **Type:** B
- **Time:** 60 min
- **Language fit:** All (Postgres)
- **Goal:** given a 5-second query in production-shaped data, get it to <50ms. Use `EXPLAIN ANALYZE`, indexes, query rewrites, and (if needed) schema changes.
- **Success criterion:** meets the latency target; explain each change.

### T7-A3: pprof / py-spy profile + flame graph
- **Type:** A
- **Time:** 45 min
- **Language fit:** Go (pprof), Python (py-spy)
- **Goal:** profile a CPU-bound handler, render a flame graph, identify the hot path.
- **Success criterion:** flame graph attached; sentence-level explanation of the top 3 frames.

---

## T8 — Security

### T8-D1: Capital One SSRF reproduction (in a sandbox!)
- **Type:** D
- **Time:** 60 min
- **Language fit:** All
- **Goal:** server that fetches user-supplied URLs (an "image proxy"). Reproduce SSRF to a fake metadata service running on `169.254.169.254`. Then mitigate.
- **Success criterion:** initial code leaks "credentials"; mitigation (allowlist hosts, block link-local ranges, IMDSv2 token requirement) closes the hole.

### T8-A1: Parameterized queries vs string concatenation
- **Type:** A / B hybrid
- **Time:** 30 min
- **Language fit:** All
- **Goal:** show a query vulnerable to SQLi (`f"WHERE name = '{name}'"`). Inject `'; DROP TABLE users; --`. Convert to parameterized; observe the attack fails.
- **Success criterion:** before/after demo; explain why parameterization is structural, not just escaping.

### T8-A2: JWT verification with key rotation
- **Type:** A
- **Time:** 90 min
- **Language fit:** Go, Python
- **Goal:** middleware that verifies JWTs against a JWKS endpoint. Handles key rotation (kid). Rejects `alg: none` and algorithm confusion attacks.
- **Success criterion:** integration test covers: valid token, expired token, wrong issuer, `alg: none`, RS256 verified with HS256 key (algorithm confusion).

### T8-A3: Secret scanning + rotation drill
- **Type:** A
- **Time:** 45 min
- **Language fit:** All
- **Goal:** wire up a secret manager (AWS Secrets Manager, HashiCorp Vault dev mode) for a DB password. Rotate the secret without restarting the service.
- **Success criterion:** service picks up the new secret within N seconds of rotation, no failed connections during rotation.

---

## T9 — DevOps adjacency

### T9-A1: Multi-stage Dockerfile + minimal image
- **Type:** A
- **Time:** 45 min
- **Language fit:** Go, Python (others adapt)
- **Goal:** Dockerfile that produces a <50MB image (Go) or a clean slim Python image. Multi-stage build, non-root user.
- **Success criterion:** image size + Trivy scan output; non-root verified.

### T9-A2: 12-Factor audit
- **Type:** B-style
- **Time:** 30 min
- **Language fit:** N/A
- **Goal:** given a sample app + Dockerfile + deploy config (planted with violations), list which factors are violated and how to fix each.
- **Success criterion:** matches the rubric (catches at least 6 of the 8 planted violations).

### T9-A3: CI pipeline with promoted artifact
- **Type:** A
- **Time:** 90 min
- **Language fit:** All
- **Goal:** GitHub Actions / GitLab CI pipeline: build once → test → push to registry → deploy to staging → manual approval → deploy to prod. Same artifact through all stages.
- **Success criterion:** working pipeline; explain why "build once, promote" beats "rebuild per env".

---

## T10 — Cloud literacy

### T10-A1: Lambda vs Fargate vs EC2 cost+latency
- **Type:** C
- **Time:** 90 min
- **Language fit:** All
- **Goal:** deploy the same simple service to all three. Measure cold-start, p99, and monthly cost at 100 / 10K / 1M requests/day.
- **Success criterion:** numbers + a recommendation tied to traffic shape.

### T10-A2: Least-privilege IAM policy
- **Type:** A
- **Time:** 45 min
- **Language fit:** N/A (cloud)
- **Goal:** write an IAM policy for a service that reads from one S3 bucket, writes to one DynamoDB table, publishes to one SQS queue. No wildcards.
- **Success criterion:** policy passes `aws iam simulate-principal-policy`; explain why each statement.

---

## T11 — Distributed systems

### T11-A1: Postgres async replication + induced split-brain
- **Type:** A / D hybrid
- **Time:** 90 min
- **Language fit:** All (Postgres)
- **Goal:** primary + async replica. Network-partition the replica. Promote it (now you have two primaries). Heal the partition. Observe write divergence.
- **Success criterion:** identify the diverged writes; describe what manual reconciliation looks like; explain why naive automatic failover causes this.

### T11-A2: Consistent hash ring with virtual nodes
- **Type:** A
- **Time:** 60 min
- **Language fit:** All
- **Goal:** implement consistent hashing with virtual nodes. Show what fraction of keys move when a node is added (should be ~1/N).
- **Success criterion:** measured reshuffle ratio close to theoretical; explain the role of virtual nodes vs physical nodes.

### T11-D1: Distributed lock without fencing tokens
- **Type:** D
- **Time:** 60 min
- **Language fit:** Go, Python
- **Goal:** Redis-based distributed lock with TTL. Process A acquires the lock, GC-pauses for longer than the TTL, lock expires, process B acquires it, process A wakes up and writes. Observe the corruption.
- **Success criterion:** reproduce the corruption; fix with fencing tokens (monotonic counter checked on every write).

### T11-A3: Raft visualization walkthrough
- **Type:** A (paper-and-screen)
- **Time:** 45 min
- **Language fit:** N/A
- **Goal:** use thesecretlivesofdata.com Raft visualization. Drive it through: leader election, log replication, partition + recovery.
- **Success criterion:** learner can explain each phase in five sentences; identify what `currentTerm` does and why it matters.

---

## Path: Real-time

### RT-A1: Live cursor / chat
- **Type:** A
- **Time:** 120 min (the dedicated builder loop)
- **Language fit:** Go, Python (FastAPI), Node
- **Goal:** WebSocket server that broadcasts cursor positions / chat messages to all connected clients. Single instance.
- **Success criterion:** two browser tabs see each other's cursors in real time.

### RT-A2: Scale RT-A1 to two server instances
- **Type:** A (continuation of RT-A1)
- **Time:** 90 min
- **Goal:** add a second server instance. Notice that clients on different instances don't see each other. Add a Redis pub/sub layer to broadcast across instances.
- **Success criterion:** four browser tabs split across two server instances all see each other.

### RT-D1: Sticky-session lock-in
- **Type:** D
- **Time:** 45 min
- **Goal:** show that a load balancer without sticky sessions causes WebSocket reconnects to land on the wrong instance, breaking presence state. Add sticky sessions; show the fix.
- **Success criterion:** demonstrate both the broken and fixed states.

---

## How to pick an exercise

- **Mid-lesson, concept just landed:** type B (break-and-fix) — fastest signal that they understood.
- **End of a tier:** type A (build-from-scratch) — capstone.
- **Trade-off discussion:** type C (compare-two) — measurement settles arguments.
- **Engineering-rationale gap (the "why does this exist" gap):** type D (reproduce-an-incident) — viscerally answers "why."

If the learner's language doesn't have prefilled scaffolding, hand them the spec from this file + the README template in `practical-mode.md`. Same hint ladder applies.
