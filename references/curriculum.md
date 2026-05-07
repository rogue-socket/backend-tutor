# Curriculum

Topic tree for backend engineering, mapped to fresh primary sources (engineering blogs, OWASP API Top 10, the 12-Factor App, Software Engineering at Google, Database Internals, Building Microservices 2e, API Design Patterns, the SRE books, and HTTP RFCs). Use this when planning a learning path, checking prerequisites, or deciding what to teach next.

## How to use this file

- Topics are listed in **rough learning order within each tier**. Tiers themselves are *roughly* ordered (T0 → T11), but plenty of cross-tier dependencies exist; respect the per-topic `prereqs`.
- Each topic has: prereqs, source pointers, key concepts, and "**understands when they can**" criteria. The "understands when" line is what you check before marking a topic done in `progress.json`.
- When a learner asks "what should I learn next?", check `progress.json` for completed topics, then suggest the earliest unmet topic with all prereqs satisfied.
- When a learner asks for a topic out of order ("teach me distributed locks" before they've done concurrency), warn them — *"this depends on T3 mutex semantics, want to cover that first or push through?"* — and honor their answer.
- **Cross-skill boundaries**: where a topic is shared with `system-design-tutor` or `ai-systems-tutor`, this file teaches the *implementation / operational* side and points the *architecture / model* side out. Don't redo what those skills own.

---

## Tier 0 — Networking & HTTP

The wire. If a learner can't reason about what actually goes over the network, every higher tier turns into cargo-culting framework code.

### 0.1 — TCP / sockets / TLS basics
- **Prereqs**: none
- **Sources**: RFC 9110 (HTTP semantics), RFC 8446 (TLS 1.3), Cloudflare blog on TLS internals, Stevens *UNIX Network Programming* (canonical but old — skim only)
- **Concepts**: 3-way handshake, TIME_WAIT, Nagle's algorithm, keep-alive, TLS handshake (client hello → cert → key exchange → finished), SNI, ALPN, mutual TLS (mTLS), cert chains, OCSP/CRL
- **Understands when they can**: explain why TLS handshake adds latency on a fresh connection and how session resumption helps; describe the difference between an idle TCP connection and a closed one and why connection pools exist.

### 0.2 — DNS
- **Prereqs**: 0.1
- **Sources**: Cloudflare *Learning* DNS series, RFC 1035 (skim)
- **Concepts**: recursive vs authoritative resolvers, A / AAAA / CNAME / MX / TXT records, TTLs, DNS caching layers (OS, resolver, browser), propagation delay, anycast DNS, split-horizon DNS, DNS-based service discovery
- **Understands when they can**: explain why a TTL change on a critical record needs to start hours before a planned cutover; describe how a misconfigured CNAME can cause outages.

### 0.3 — HTTP/1.1 semantics
- **Prereqs**: 0.1
- **Sources**: **RFC 9110** (HTTP semantics — read), **API Design Patterns** ch. 1-4 (Geewax)
- **Concepts**: methods (GET / POST / PUT / PATCH / DELETE / HEAD / OPTIONS), status code families, idempotent vs safe methods, headers (request / response / entity), content negotiation (Accept, Content-Type), conditional requests (If-Match, If-None-Match), range requests, chunked transfer encoding, Connection: keep-alive, HTTP pipelining (and why it failed)
- **Understands when they can**: pick the right method for a given operation by semantics not convention; explain why `PUT` is idempotent and `POST` is not; reach for `409 Conflict` vs `422 Unprocessable Entity` vs `400 Bad Request` correctly.

### 0.4 — HTTP/2
- **Prereqs**: 0.3
- **Sources**: **Cloudflare blog — HTTP/2 explainer**, RFC 9113
- **Concepts**: binary framing, multiplexing (no head-of-line at HTTP layer, but TCP HoL still bites), HPACK header compression, server push (deprecated in practice), stream priorities, flow control
- **Understands when they can**: explain what HTTP/2 fixes vs HTTP/1.1 and what it *doesn't* fix; describe head-of-line blocking at the TCP layer.

### 0.5 — HTTP/3 / QUIC
- **Prereqs**: 0.4
- **Sources**: RFC 9000 (QUIC), RFC 9114 (HTTP/3), Cloudflare blog on HTTP/3 deployment
- **Concepts**: QUIC over UDP, 0-RTT and 1-RTT handshakes, connection migration, no head-of-line at transport layer, TLS 1.3 baked in
- **Understands when they can**: explain why QUIC fixes TCP head-of-line blocking; describe one production gotcha (UDP middlebox interference, kernel CPU cost).

### 0.6 — WebSockets
- **Prereqs**: 0.3
- **Sources**: RFC 6455, **Discord eng blog** (scaling presence to millions), **Figma multiplayer post**
- **Concepts**: HTTP upgrade handshake, frames (text / binary / control), ping/pong, close codes, message vs frame, backpressure (when the receiver can't keep up), origin checking, subprotocols
- **Understands when they can**: walk through the upgrade handshake; explain why you need to handle backpressure explicitly (no built-in flow control beyond TCP); name two reasons a long-lived WebSocket connection breaks (idle proxy timeout, NAT rebinding).

### 0.7 — Server-Sent Events (SSE)
- **Prereqs**: 0.3
- **Sources**: WHATWG HTML Living Standard §EventSource, MDN
- **Concepts**: text/event-stream content type, single-direction (server → client), auto-reconnect with `Last-Event-ID`, line-based framing, vs WebSockets vs long-polling
- **Understands when they can**: pick SSE over WebSockets when bidirectional isn't needed and load-balancer compatibility matters; explain why SSE works through more middleboxes than WebSockets.

### 0.8 — gRPC transport
- **Prereqs**: 0.4
- **Sources**: grpc.io docs, **Building Microservices 2e** ch. 5
- **Concepts**: Protobuf wire format, HTTP/2 underneath, four streaming kinds (unary, server-stream, client-stream, bidi), deadlines (vs timeouts), interceptors / middleware, gRPC-Web for browsers
- **Understands when they can**: choose gRPC vs REST by use case (internal service-to-service, polyglot codebase, high RPC volume); explain why deadlines must propagate across services.

### 0.9 — CDNs and edge
- **Prereqs**: 0.2, 0.3
- **Sources**: Cloudflare and Fastly engineering blogs
- **Concepts**: anycast routing, edge POPs, cache hierarchy (edge / shield / origin), Vary headers, Cache-Control directives (s-maxage, stale-while-revalidate, stale-if-error), purging, CDN as DDoS shield
- **Understands when they can**: pick the right Cache-Control directives for a public vs private response; explain when `Vary: Cookie` accidentally kills your hit rate.

---

## Tier 1 — APIs

How services talk. Every wrong decision here ossifies into a public contract you can't break.

### 1.1 — REST design
- **Prereqs**: 0.3
- **Sources**: **API Design Patterns** (Geewax) ch. 5-10, **Stripe API docs** as canonical example
- **Concepts**: resources vs RPC, hierarchical URIs, plural nouns, sub-resources, batch endpoints, error envelopes (RFC 7807 problem+json), HATEOAS (and why almost no one does it), filtering / sorting / field selection
- **Understands when they can**: design a CRUD resource with sub-resources without inventing a new style every time; explain why "REST vs RPC" is a spectrum, not a binary.

### 1.2 — Idempotency
- **Prereqs**: 0.3, 1.1
- **Sources**: **Stripe blog — "Designing robust and predictable APIs with idempotency"**, **API Design Patterns** ch. 26
- **Concepts**: idempotent HTTP methods (built-in for GET/PUT/DELETE), idempotency keys for POST, server-side key storage / TTL, the "in-flight" state problem, key collision under concurrency, what to return on a duplicate
- **Understands when they can**: implement an idempotency-key system server-side that handles concurrent duplicate requests correctly (not just sequential ones); explain why caching the *response body* is part of the contract.

### 1.3 — API versioning
- **Prereqs**: 1.1
- **Sources**: **API Design Patterns** ch. 24, Stripe API versioning docs
- **Concepts**: URI versioning (`/v1/...`), header versioning, content-type versioning, date-based versioning (Stripe), breaking vs non-breaking changes, additive evolution, deprecation timelines, version pinning
- **Understands when they can**: classify a change as breaking or non-breaking by the contract it changes (request shape, response shape, error codes, semantics); design a 12-month deprecation flow.

### 1.4 — Pagination
- **Prereqs**: 1.1
- **Sources**: **API Design Patterns** ch. 21, Slack / GitHub / Stripe API pagination docs
- **Concepts**: offset/limit (and why it breaks at scale), cursor-based, keyset (seek-method), opaque cursors, total counts (and why they're often a lie)
- **Understands when they can**: explain why offset pagination becomes O(n) at deep pages; implement cursor pagination on top of an indexed sort key.

### 1.5 — Errors and error envelopes
- **Prereqs**: 0.3, 1.1
- **Sources**: **RFC 7807** (Problem Details), **API Design Patterns** ch. 25
- **Concepts**: status code semantics, error object shape (code, message, fields, request id), retryable vs non-retryable, validation errors vs business errors, security-sensitive error leakage
- **Understands when they can**: design an error envelope that distinguishes "client should retry," "client should fix their request," and "server has a bug" — without leaking internals.

### 1.6 — Rate limiting
- **Prereqs**: 1.1
- **Sources**: **Stripe blog — "Scaling your API with rate limiters"**, **Cloudflare blog**
- **Concepts**: token bucket, leaky bucket, fixed window, sliding window, sliding log; 429 + Retry-After + RateLimit-* headers (RFC 9331), per-user vs per-IP vs per-key, distributed rate limiters (Redis-based)
- **Understands when they can**: pick the right algorithm for the traffic shape; explain why a naive Redis INCR + EXPIRE has a race; reason about cost-of-attack vs cost-of-defense.

### 1.7 — Auth: sessions
- **Prereqs**: 1.1
- **Sources**: OWASP Session Management Cheat Sheet, **OWASP API Top 10 — API2:2023 (Broken Authentication)**
- **Concepts**: session ID generation (CSPRNG), server-side session store (Redis / DB), cookie attributes (Secure, HttpOnly, SameSite, Path, Domain), CSRF protection, session fixation, idle vs absolute timeouts, logout
- **Understands when they can**: name three concrete cookie-flag misconfigurations and what each opens up; explain why session rotation on privilege change matters.

### 1.8 — Auth: JWT (and when not to)
- **Prereqs**: 1.7
- **Sources**: RFC 7519, **Auth0 blog**, *"Stop using JWT for sessions"* by Sven Slootweg (read both sides)
- **Concepts**: header / payload / signature, signing algorithms (HS256 vs RS256/EdDSA), key rotation (kid), claims (iss, aud, exp, nbf, iat, sub, jti), the `alg: none` attack, token revocation problem, refresh tokens, short-vs-long expiry trade-off
- **Understands when they can**: name two reasons stateful sessions beat JWT for first-party web apps; explain when JWT genuinely earns its keep (cross-domain SSO, internal service-to-service).

### 1.9 — OAuth 2.0 and OIDC
- **Prereqs**: 1.7, 1.8
- **Sources**: RFC 6749 (OAuth 2.0), OIDC core spec, **Okta dev blog**, **OAuth 2.0 in Action** (Manning)
- **Concepts**: roles (resource owner, client, authorization server, resource server), grant types (authorization code + PKCE, client credentials, device code; *not* implicit, *not* password), scopes, ID token vs access token, OIDC discovery, JWKS, refresh token rotation, redirect URI validation
- **Understands when they can**: pick the right grant type for a given app shape (SPA, mobile, server-side web, machine-to-machine); explain why PKCE is non-negotiable for public clients.

### 1.10 — gRPC service design
- **Prereqs**: 0.8
- **Sources**: grpc.io docs, **Building Microservices 2e** ch. 5
- **Concepts**: Proto3 conventions, message evolution (field numbers as the contract, reserved fields), error model (status codes + details), deadlines, retries, load balancing (server-side vs client-side), reflection
- **Understands when they can**: evolve a `.proto` file without breaking existing clients; design a service interface that handles partial failure across a streaming RPC.

### 1.11 — GraphQL basics
- **Prereqs**: 1.1
- **Sources**: graphql.org docs, **Apollo blog**, **GitHub GraphQL API as case study**
- **Concepts**: schema (types, queries, mutations, subscriptions), resolvers, the N+1 problem, DataLoader pattern, query complexity limits, persisted queries, schema federation
- **Understands when they can**: explain why naive resolvers cause N+1 and how DataLoader fixes it; reason about when GraphQL earns its complexity (BFFs, frontend-driven data needs) and when it doesn't.

### 1.12 — Real-time API patterns
- **Prereqs**: 0.6, 0.7, 1.1
- **Sources**: **Discord engineering**, **Figma multiplayer post**, **Slack RTM/Events API design**
- **Concepts**: long-polling, SSE, WebSocket message protocols (custom vs Phoenix Channels / Socket.IO / GraphQL subscriptions), presence, broadcast, message ordering, reconnection logic with replay (Last-Event-ID, sequence numbers)
- **Understands when they can**: pick between SSE / WebSocket / long-polling by the actual requirements; design a reconnection protocol that doesn't lose or double-deliver messages.

### 1.13 — Webhooks
- **Prereqs**: 0.3, 1.2
- **Sources**: **Stripe webhook docs**, **GitHub webhook docs**
- **Concepts**: at-least-once delivery, signature verification (HMAC), replay protection (timestamp + tolerance window), retry semantics, ordering guarantees (none by default), idempotent receivers
- **Understands when they can**: implement a webhook receiver that handles duplicates and out-of-order delivery; verify signatures correctly without timing-leak comparisons.

### 1.14 — Contract testing
- **Prereqs**: 1.1
- **Sources**: **Pact docs**, **Software Engineering at Google** ch. 14, OpenAPI/Swagger specs
- **Concepts**: consumer-driven contracts, OpenAPI as a contract source, contract test vs integration test vs end-to-end, schema validation in CI
- **Understands when they can**: explain why end-to-end tests rot and contract tests don't; set up a basic Pact or schema-validation flow between two services.

---

## Tier 2 — Databases

Where the truth lives. Every backend engineer eventually owns one. Most don't until something breaks.

### 2.0 — DB types tour and pick-when matrix
- **Prereqs**: 1.1
- **Sources**: **Database Internals** (Petrov) ch. 1-2 (overview), **DDIA** ch. 2 (data models — still gold), engineering blogs from Discord, Notion, Figma on choosing their DB
- **Concepts**: relational (Postgres, MySQL), document (Mongo, DynamoDB doc-style), key-value (Redis, DynamoDB KV-style, etcd), wide-column (Cassandra, ScyllaDB), columnar (ClickHouse, BigQuery, Snowflake), time-series (TimescaleDB, InfluxDB, Prometheus), graph (Neo4j, dgraph), search (Elasticsearch, OpenSearch, Meilisearch), vector (pgvector, Pinecone, Qdrant)
- **Pick-when matrix** (covered in this lesson): "what's your workload?" → relationships + transactions = relational; key-by-key high-throughput = KV; full-text relevance = search; analytical aggregates over wide tables = columnar; semantic similarity = vector; metric ingest at high cardinality = time-series; "I don't know yet" = relational, period.
- **Understands when they can**: defend a choice of database for a given workload using QPS, data shape, consistency need, query patterns, ops cost — not just "I know Postgres."

### 2.1 — SQL fluency
- **Prereqs**: 2.0
- **Sources**: **PostgreSQL docs** (read SELECT, WITH, window functions), **Use the Index, Luke!** by Markus Winand (free online — *the* practical SQL indexing book), Modern SQL by the same author
- **Concepts**: joins (INNER / LEFT / RIGHT / FULL, semi-joins, anti-joins), GROUP BY + HAVING, window functions (ROW_NUMBER, LAG/LEAD, PARTITION BY), CTEs (recursive too), subqueries vs joins, NULL semantics
- **Understands when they can**: write a query with a window function to compute a per-user running total; reason about NULL behavior in WHERE and JOIN.

### 2.2 — Indexes and the query planner
- **Prereqs**: 2.1
- **Sources**: **Use the Index, Luke!**, **Database Internals** ch. 4, Postgres `EXPLAIN ANALYZE` docs
- **Concepts**: B-tree (default, range queries), hash (equality only, in-memory), GIN (full-text, JSONB, arrays), GiST (geo, range types), BRIN (huge sequential tables), partial indexes, covering indexes (INCLUDE), composite indexes and column order, index-only scans, the query planner (cost-based, statistics, ANALYZE), why ORM-generated queries break the planner
- **Understands when they can**: read `EXPLAIN ANALYZE` output and find the bottleneck; pick the right index for a given query and explain why; recognize when an index is *hurting* (write amplification).

### 2.3 — Transactions and isolation
- **Prereqs**: 2.1
- **Sources**: **Database Internals** ch. 5, **DDIA** ch. 7 (still gold), **Jepsen** posts (rigor on isolation levels in real DBs)
- **Concepts**: ACID, isolation levels (Read Uncommitted, Read Committed, Repeatable Read, Serializable, Snapshot Isolation), anomalies (dirty read, non-repeatable read, phantom, write skew, lost update), MVCC vs locking, SELECT FOR UPDATE / FOR SHARE / SKIP LOCKED, deadlock detection
- **Understands when they can**: explain why Postgres' "Repeatable Read" is actually Snapshot Isolation and what that misses (write skew); pick the right isolation level for a given workload without overshooting to Serializable.

### 2.4 — Schema design
- **Prereqs**: 2.1
- **Sources**: **DDIA** ch. 2-3, **Database Design for Mere Mortals** (Hernandez), real Postgres schemas from open-source projects
- **Concepts**: normalization (1NF/2NF/3NF), denormalization for read paths, surrogate vs natural keys, soft deletes (and their downsides), audit columns, ENUMs vs lookup tables, JSONB columns (when yes, when no), tenant-isolation strategies (column / schema / database)
- **Understands when they can**: design a schema for a given domain; defend choices on normalization vs denormalization with read/write ratios.

### 2.5 — Migrations
- **Prereqs**: 2.4
- **Sources**: **GitHub blog** on online migrations, **Stripe blog** on online schema changes, gh-ost / pt-online-schema-change docs, **Strong Migrations** (Ankane) — Rails-flavored but the rules apply everywhere
- **Concepts**: online vs offline migrations, expand/contract pattern, backfilling without locking, dual-write, dual-read, NOT NULL on existing columns (locking trap), index creation (CONCURRENTLY in Postgres), foreign keys, migration rollback (often impossible — use forward-fix)
- **Understands when they can**: walk through a 5-step expand/contract migration for adding a NOT NULL column to a 50M-row table without downtime; identify which Postgres DDL statements take ACCESS EXCLUSIVE locks.

### 2.6 — ORMs vs raw SQL, and N+1
- **Prereqs**: 2.1, 2.2
- **Sources**: **Use the Index, Luke!** ch. on ORMs, sqlc / Prisma / SQLAlchemy / GORM docs (compare patterns)
- **Concepts**: active record vs data mapper vs query builder vs raw, lazy loading vs eager loading, the N+1 query problem, batch loading / DataLoader, ORM-generated query pathologies, when to drop to raw SQL
- **Understands when they can**: spot N+1 in a code review; rewrite an N+1 ORM query using `JOIN` or a batch fetch.

### 2.7 — Connection pooling
- **Prereqs**: 2.1
- **Sources**: **PgBouncer docs**, **HikariCP docs** (the canonical sizing rationale), AWS RDS Proxy
- **Concepts**: server vs client pools, transaction pooling vs session pooling vs statement pooling (PgBouncer modes), pool sizing (it's *less* than you think — see HikariCP docs), prepared statements + transaction pooling incompatibility, connection limits (Postgres default 100, why that matters)
- **Understands when they can**: size a pool given concurrency * (connection-time / total-time); explain why doubling a pool size often *hurts* throughput.

### 2.8 — Postgres deep-dive
- **Prereqs**: 2.2, 2.3
- **Sources**: **Postgres docs**, **Postgres at Scale** talks from PGCon, **Notion's Postgres sharding post**
- **Concepts**: MVCC (heap + visibility map), VACUUM and autovacuum, transaction ID wraparound, WAL, replication slots, tablespaces, extensions (pgvector, pg_trgm, postgis, pg_stat_statements), JSONB indexing, partitioning (declarative)
- **Understands when they can**: diagnose bloat from long-running transactions; configure autovacuum sensibly for a high-write table; pick partitioning column for a time-series workload.

### 2.9 — Redis
- **Prereqs**: 2.0
- **Sources**: **Redis docs**, **Redis Best Practices** (mostly via official blog + talks), **Redis Streams**, the Redlock paper *and* Martin Kleppmann's critique
- **Concepts**: data types (string, hash, list, set, sorted set, stream, bitmap, hyperloglog), persistence (RDB vs AOF), eviction policies (allkeys-lru, volatile-ttl, etc.), pipelining, Lua scripting (atomic), pub/sub vs streams, Redis Cluster (slot-based sharding), Sentinel (HA)
- **Understands when they can**: pick the right Redis data type for a given access pattern; explain why Redis pub/sub is fire-and-forget and Streams isn't; reason about Redlock's safety claims.

### 2.10 — Search engines
- **Prereqs**: 2.0
- **Sources**: **Elasticsearch: The Definitive Guide** (older but free, still valid concepts), Elastic blog, OpenSearch docs
- **Concepts**: inverted indexes, analyzers (tokenizer + filters), term vs phrase queries, BM25 scoring, refresh interval, near-real-time, denormalization for search (vs joining), shard sizing (small enough to recover, large enough to pack), index lifecycle management
- **Understands when they can**: explain why Elastic / OpenSearch need denormalized documents; tune relevance with analyzer changes vs scoring boosts.

### 2.11 — Vector databases
- **Prereqs**: 2.0
- **Sources**: **pgvector docs**, Pinecone / Qdrant blogs (be skeptical of vendor framing), Lucene 9.x HNSW
- **Concepts**: dense embeddings, similarity metrics (cosine, L2, dot), exact vs approximate nearest neighbor (ANN), HNSW graph, IVFFlat, recall vs latency vs build-time trade-off, hybrid search (BM25 + vector)
- **Understands when they can**: defend "use Postgres + pgvector" vs "use a dedicated vector DB" with traffic and recall numbers; explain when ANN's recall trade-off matters and when it doesn't.

### 2.12 — Backups and disaster recovery
- **Prereqs**: 2.8
- **Sources**: **GitLab 2017 incident** (the one where backups didn't work), Postgres docs on PITR, AWS RDS / Aurora backup model
- **Concepts**: logical vs physical backups, point-in-time recovery (PITR), WAL archiving, RTO vs RPO, restore drills (untested backups don't exist), cross-region copies
- **Understands when they can**: define RTO and RPO for a hypothetical service; design a backup schedule that meets a 1-hour RPO; argue for why a quarterly restore drill is non-negotiable.

---

## Tier 3 — Concurrency & async

The land of "it works on my machine in single-threaded mode." Backend interviews lean on this disproportionately and rightly so.

### 3.1 — Concurrency models compared
- **Prereqs**: none (but: language familiarity in the learner's chosen language)
- **Sources**: **The Go Memory Model**, **Python asyncio docs**, **Node.js event loop docs**, *Concurrency in Go* (Cox-Buday)
- **Concepts**: OS threads vs lightweight threads (goroutines, virtual threads in Java 21+), event loop (single-threaded async, Node, Python asyncio), the GIL (Python), shared-memory vs message-passing, when each model wins
- **Understands when they can**: explain why CPU-bound work in Python is best multi-process and IO-bound is fine async; describe what a goroutine actually is at the OS level.

### 3.2 — Race conditions and atomicity
- **Prereqs**: 3.1
- **Sources**: **The Go Memory Model**, **Java Concurrency in Practice** (still gold), **Database Internals** ch. 5
- **Concepts**: data races, atomicity, compare-and-swap (CAS), memory ordering / happens-before, the read-modify-write trap, why `counter++` isn't atomic
- **Understands when they can**: identify a race in a code review; pick between mutex and atomic for a given operation.

### 3.3 — Mutexes, RWMutexes, and the cost of locks
- **Prereqs**: 3.2
- **Sources**: language stdlib docs, **Java Concurrency in Practice**
- **Concepts**: mutex acquisition cost, contention, lock granularity (coarse vs fine), reader-writer locks (and when they're worse than a regular mutex), reentrancy, the lock ordering rule, condition variables
- **Understands when they can**: explain why an RWMutex with mostly-readers can still bottleneck on the writer's contention with readers; reason about lock granularity trade-offs.

### 3.4 — Deadlocks, livelocks, starvation
- **Prereqs**: 3.3
- **Sources**: **Java Concurrency in Practice**, **Operating Systems** (Tanenbaum) ch. 6
- **Concepts**: the four conditions for deadlock, lock ordering as cure, deadlock detection, livelock, starvation, fairness
- **Understands when they can**: walk through a deadlock scenario in their language and fix it via lock ordering; explain why retrying naively turns deadlock into livelock.

### 3.5 — Idempotency in handlers
- **Prereqs**: 1.2, 3.2
- **Sources**: **Stripe idempotency post** (revisit), **Building Microservices 2e** ch. 12
- **Concepts**: implementing idempotency keys with a dedicated table, the upsert race, idempotency under crash (key recorded but response not stored), exactly-once delusion, the difference between idempotent operations and idempotent endpoints
- **Understands when they can**: implement an idempotent POST endpoint that survives a process crash mid-request; explain why "exactly-once" doesn't exist and what people mean when they claim it.

### 3.6 — Job queues and work distribution
- **Prereqs**: 3.1
- **Sources**: **AWS SQS docs**, **RabbitMQ docs**, **Apache Kafka** (Kreps essays), **Sidekiq / Celery / River** docs depending on language
- **Concepts**: producer / consumer / broker, topics vs queues, partitions, consumer groups, ordering (per-partition only), visibility timeouts, ack / nack, FIFO vs standard, message TTL
- **Understands when they can**: explain why Kafka's per-partition ordering doesn't give you global ordering and why that's a feature; pick SQS vs Kafka vs RabbitMQ for a given workload.

### 3.7 — Delivery semantics
- **Prereqs**: 3.5, 3.6
- **Sources**: **Kafka docs on EOS**, **Confluent blog on exactly-once**, **The Log** (Kreps)
- **Concepts**: at-most-once, at-least-once, "effectively-once" / "exactly-once" via idempotent consumers + transactional offsets, the two-generals problem, why retries are mandatory and dedup is your problem
- **Understands when they can**: defend "at-least-once + idempotent consumer" as the default for almost all workloads; explain Kafka's exactly-once semantics and what they do *not* cover.

### 3.8 — Retries with backoff and jitter
- **Prereqs**: 3.6
- **Sources**: **AWS Architecture Blog — "Exponential backoff and jitter"** (Marc Brooker — the canonical post), **Google SRE book** ch. 22
- **Concepts**: exponential backoff, full jitter vs equal jitter vs decorrelated jitter, retry budgets, retry storm prevention, idempotency as a retry prerequisite, *when not to retry* (4xx vs 5xx, deadline exceeded propagation)
- **Understands when they can**: implement decorrelated jitter and explain why it beats fixed-interval retries under load; identify cases where retrying is harmful (cascading failures).

### 3.9 — Dead-letter queues and poison messages
- **Prereqs**: 3.6, 3.8
- **Sources**: AWS SQS docs on DLQ, **Building Microservices 2e** ch. 13
- **Concepts**: poison messages (always fail), max-receive count → DLQ move, replay tooling, DLQ alarms, observability for DLQ depth
- **Understands when they can**: design a DLQ flow with replay; explain why DLQ depth is a tier-1 alert.

### 3.10 — Saga and distributed transactions
- **Prereqs**: 3.5, 3.7, 2.3
- **Sources**: **Building Microservices 2e** ch. 6, **Microservices Patterns** (Richardson) ch. 4-5, **the Outbox pattern**
- **Concepts**: 2PC (and why it's avoided), Saga pattern (orchestration vs choreography), compensating actions, the outbox pattern (transactional event publishing), inbox pattern, idempotent compensations
- **Understands when they can**: design a saga for a multi-service workflow with compensations; defend the outbox pattern over dual-write.

---

## Tier 4 — Caching

Speed and cost. Cache invalidation is one of the two hard problems for a reason.

### 4.1 — Cache patterns
- **Prereqs**: 2.0
- **Sources**: **AWS caching whitepaper**, **DDIA** ch. 11 (skim), real engineering blogs (Discord, Shopify) on cache architecture
- **Concepts**: cache-aside (lazy load), read-through, write-through, write-behind / write-back, refresh-ahead
- **Understands when they can**: pick the right pattern for a given consistency / latency / cost trade-off.

### 4.2 — Invalidation
- **Prereqs**: 4.1
- **Sources**: engineering blogs from Etsy, Pinterest, Shopify on cache invalidation; **Phil Karlton's quote** (the cliché has substance)
- **Concepts**: TTL-based, key-versioning (`user:42:v3`), pub/sub invalidation, write-through invalidation, the staleness window, soft TTL vs hard TTL
- **Understands when they can**: design an invalidation strategy that doesn't depend on perfect ordering between cache writes and DB writes.

### 4.3 — Cache stampedes
- **Prereqs**: 4.1
- **Sources**: **Discord's stampede incident postmortem**, single-flight (Go), **probabilistic early expiration paper** (Vattani et al.)
- **Concepts**: thundering herd, single-flight / request coalescing, probabilistic early refresh (XFetch), background refresh, locking strategies and their downsides
- **Understands when they can**: implement single-flight in their chosen language; explain why a 10s TTL cache on a 1000 QPS hot key without single-flight is an outage waiting to happen.

### 4.4 — Cache hierarchies
- **Prereqs**: 0.9, 4.1
- **Sources**: Cloudflare engineering blog
- **Concepts**: in-process (LRU) → distributed (Redis) → CDN edge — and why each layer's cost / staleness / hit rate differs; cache key design (don't cache the URL; cache the canonicalized request); negative caching
- **Understands when they can**: walk through the layers a cache miss traverses and the latency at each.

---

## Tier 5 — Reliability

Surviving the day everything goes wrong. Mostly disciplines, not features.

### 5.1 — Timeouts and deadlines
- **Prereqs**: 0.3
- **Sources**: **Marc Brooker's blog** (timeouts and retries), **gRPC docs on deadlines**
- **Concepts**: timeouts as a budget that propagates, deadline propagation across services, default-no-timeout libraries (the bug), client timeout < server timeout < downstream timeout (the cascade rule)
- **Understands when they can**: identify a missing timeout in a code review; design a timeout budget across a 3-hop call chain.

### 5.2 — Retries (revisit, with reliability framing)
- **Prereqs**: 3.8
- **Sources**: **SRE book** ch. 22, Marc Brooker's blog
- **Concepts**: retry budgets (token bucket on retries), client-side vs gateway-side retries (don't double-retry), retry-after as a signal
- **Understands when they can**: explain why naive retries amplify load during a partial outage; describe a retry budget.

### 5.3 — Circuit breakers and bulkheads
- **Prereqs**: 5.1, 5.2
- **Sources**: **Release It!** (Nygard) ch. 5-6, Hystrix docs (historical), Resilience4j / gobreaker
- **Concepts**: closed / open / half-open states, threshold tuning, failure isolation (bulkhead), thread pool isolation vs semaphore isolation
- **Understands when they can**: tune a circuit breaker's thresholds for a downstream that's degraded vs hard-down; explain why bulkheads protect *you* from a single misbehaving dependency.

### 5.4 — Hedged requests
- **Prereqs**: 5.1
- **Sources**: **"The Tail at Scale"** (Dean & Barroso, 2013), Google SRE book
- **Concepts**: send a second request after a delay shorter than p99 to a different replica, take the first response — trade extra cost for tail-latency cut
- **Understands when they can**: defend hedging with a cost-vs-tail-latency calculation; identify when hedging makes things worse (writes, expensive operations).

### 5.5 — Health checks
- **Prereqs**: 0.3
- **Sources**: **Kubernetes docs** on probes, **Microservices Patterns** ch. 11
- **Concepts**: liveness (am I broken — restart me) vs readiness (can I serve — pull me from LB) vs startup (give me time to warm), the dependency-cascade problem (don't fail readiness on a remote dependency), health endpoint as a load-bearing API
- **Understands when they can**: explain why a readiness probe that checks DB connectivity can take down a cluster during a DB hiccup; design probes that fail safely.

### 5.6 — Feature flags and kill switches
- **Prereqs**: none
- **Sources**: **LaunchDarkly blog**, GitHub's Scientist library, **Feature Flags book** (free)
- **Concepts**: rollout flags (gradual %), kill switches (instant disable), per-user vs per-cohort, flag debt, flag-as-config vs experimentation
- **Understands when they can**: gate a risky deploy behind a kill switch; explain why flag debt is a real cost.

### 5.7 — Deploy strategies
- **Prereqs**: 5.5, 5.6, 9.x
- **Sources**: **GitHub blog** on deployments, **Spinnaker docs**
- **Concepts**: blue / green, canary (% rollout), shadow traffic (mirror), rolling, recreate, roll-back vs roll-forward, deploy under load
- **Understands when they can**: pick the strategy by risk profile and rollback cost; design canary metrics (error rate, latency, business signal).

### 5.8 — Graceful shutdown
- **Prereqs**: 0.1, 3.1
- **Sources**: language stdlib docs (e.g., Go `os/signal`), Kubernetes docs on terminationGracePeriodSeconds
- **Concepts**: SIGTERM handling, draining inflight requests, closing connections, `preStop` hooks in K8s, leader-election relinquish, the 30-second default
- **Understands when they can**: implement graceful shutdown that drains a request queue and closes a DB pool; explain why `kill -9` is a bug, not a deploy strategy.

---

## Tier 6 — Observability & on-call

What you can see when you're paged at 3am and have 15 minutes.

### 6.1 — Structured logs
- **Prereqs**: none
- **Sources**: **12-Factor App XI (Logs)**, Datadog / Honeycomb / Splunk docs
- **Concepts**: JSON logs, log levels (and why TRACE/DEBUG should be off in prod by default), correlation IDs / request IDs, context propagation, PII scrubbing, log volume cost
- **Understands when they can**: design a log line with the fields needed to debug a request without leaking PII; trace a request across two services using correlation IDs.

### 6.2 — Metrics: counters, gauges, histograms
- **Prereqs**: none
- **Sources**: **Prometheus docs**, OpenTelemetry metrics spec
- **Concepts**: counter (monotonic), gauge (point-in-time), histogram (bucketed distribution), summary (quantiles, but client-side — usually wrong), labels and cardinality (the explosion problem), rate() and irate()
- **Understands when they can**: pick the right metric type; explain why high-cardinality labels (user_id, request_id) blow up Prometheus.

### 6.3 — RED and USE methods
- **Prereqs**: 6.2
- **Sources**: **Brendan Gregg's USE method**, **Tom Wilkie's RED method** (Weave / Grafana Labs)
- **Concepts**: RED (Rate, Errors, Duration) for request-driven services; USE (Utilization, Saturation, Errors) for resources; Four Golden Signals (SRE book — latency, traffic, errors, saturation)
- **Understands when they can**: dashboard a service using RED; reason about saturation on a thread pool or DB pool.

### 6.4 — Distributed tracing
- **Prereqs**: 6.1
- **Sources**: **OpenTelemetry docs**, **Honeycomb's blog** on observability, **Distributed Tracing in Practice** (Parker et al.)
- **Concepts**: spans, traces, parent-child relationships, W3C Trace Context (`traceparent` header), context propagation across HTTP / gRPC / queues, span attributes vs events vs links, sampling (head sampling, tail sampling), baggage
- **Understands when they can**: instrument a service to emit a span with attributes; trace a request across two services; explain why head sampling drops the bad ones unless you're careful.

### 6.5 — SLOs, SLIs, error budgets
- **Prereqs**: 6.2, 6.3
- **Sources**: **The SRE Book** ch. 4, **The SRE Workbook** ch. 1-3
- **Concepts**: SLI (a measurable proxy for "is the service good"), SLO (a target on the SLI), SLA (the contract), error budget = 1 - SLO, error budget burn rate alerts
- **Understands when they can**: write an SLI for a hypothetical API; design a 99.9% SLO and the alerts that fire on burn-rate breach.

### 6.6 — Alerting and runbooks
- **Prereqs**: 6.5
- **Sources**: **Rob Ewaschuk's "My Philosophy on Alerting"** (canonical), **SRE book** ch. 6
- **Concepts**: page-worthy = symptom + actionable + urgent, multi-window multi-burn-rate alerts (Google's pattern), runbooks (what's broken / what to check / what to do), alert fatigue, deduplication
- **Understands when they can**: classify an alert as page-worthy or ticket-worthy; write a runbook for a hypothetical "p99 latency breached SLO" alert.

### 6.7 — Postmortems and incident response
- **Prereqs**: 6.6
- **Sources**: **SRE book** ch. 15, real public postmortems (GitHub, Cloudflare, AWS), **Etsy debriefing facilitator's guide**
- **Concepts**: blameless postmortem, timeline, contributing factors (not "root cause"), action items with owners, retro vs postmortem, IC / commander roles in incident response
- **Understands when they can**: facilitate a postmortem that doesn't blame an individual; write contributing factors that point at systems, not people.

---

## Tier 7 — Performance & scale

Where capacity meets reality.

### 7.1 — Latency vs throughput, tail latency
- **Prereqs**: 6.2
- **Sources**: **"The Tail at Scale"** (Dean & Barroso), **Gil Tene's "How NOT to Measure Latency"** talk
- **Concepts**: p50 / p95 / p99 / p99.9, why averages lie, coordinated omission, the "fan-out amplifies tail" effect, latency budgets
- **Understands when they can**: explain why a 100-request fan-out service has p99 ≈ p99.99 of the backend; identify coordinated omission in a load test.

### 7.2 — Profiling
- **Prereqs**: 7.1
- **Sources**: **Brendan Gregg's flame graph posts**, language profilers (Go pprof, Python py-spy, Node clinic, Java async-profiler)
- **Concepts**: CPU profile (sampling), allocation profile, blocking profile, flame graphs, off-CPU analysis, microbenchmarks (and how they lie)
- **Understands when they can**: read a flame graph; pick CPU vs allocation vs blocking profile for a given symptom.

### 7.3 — Load testing
- **Prereqs**: 7.1
- **Sources**: **k6 docs**, **vegeta**, **Locust**, Gil Tene on load gen
- **Concepts**: open-loop (constant arrival rate) vs closed-loop (constant concurrency), realistic distributions, ramp-up, soak tests, the load-tester-is-the-bottleneck trap, production traffic shadowing
- **Understands when they can**: design a load test that mirrors production traffic patterns; explain why closed-loop tools (typical) under-report tail latency.

### 7.4 — Capacity estimation
- **Prereqs**: 7.1
- **Sources**: **DDIA** ch. 1 (still gold), Jeff Dean's "Numbers Every Programmer Should Know"
- **Concepts**: Little's Law (L = λW), back-of-envelope storage / bandwidth / QPS calculations, peak vs average (4-10x), headroom, cost as a constraint
- **Understands when they can**: estimate the disk for a write-heavy service ("10M users, 200 events/day, 1KB each, 1-year retention"); apply Little's Law to size a thread pool.

### 7.5 — Query tuning (revisit)
- **Prereqs**: 2.2
- **Sources**: **Use the Index, Luke!** (revisit), Postgres `EXPLAIN ANALYZE` deep-dive
- **Concepts**: rewriting a slow query, partial / covering / composite indexes, query planner statistics drift, parameter sniffing, workload-driven schema redesign
- **Understands when they can**: take a 5-second query and reduce it to <50ms via index + plan analysis.

### 7.6 — Hot paths and optimization discipline
- **Prereqs**: 7.2
- **Sources**: **Don Knuth's "premature optimization is the root of all evil" — in full context**, Brendan Gregg
- **Concepts**: profile first, fix the biggest, measure again; the 80/20 of perf work; benchmarking discipline; representativeness of synthetic workloads
- **Understands when they can**: defend or reject a perf change with before / after numbers from a representative workload.

---

## Tier 8 — Security

OWASP API Top 10 (2023) is the spine. Everything else hangs off it.

### 8.1 — OWASP API Security Top 10 (2023)
- **Prereqs**: 1.7
- **Sources**: **OWASP API Security Top 10 (2023)** — the full list, read it cover to cover
- **Concepts**: API1 BOLA (object-level authz), API2 Broken Authentication, API3 Broken Object Property Level Authorization, API4 Unrestricted Resource Consumption, API5 Broken Function Level Authorization, API6 Unrestricted Access to Sensitive Business Flows, API7 SSRF, API8 Security Misconfiguration, API9 Improper Inventory Management, API10 Unsafe Consumption of APIs
- **Understands when they can**: name the difference between API1 (BOLA) and API5 (BFLA) with a concrete example each; identify SSRF (API7) in a code review of an outbound HTTP call.

### 8.2 — Injection
- **Prereqs**: 2.1
- **Sources**: **OWASP Top 10 (web)** A03:2021 Injection, **Bobby Tables**
- **Concepts**: SQLi (parameterized queries, never string-concatenate), NoSQL injection (Mongo `$where`, etc.), command injection (no `exec(user_input)`), template injection
- **Understands when they can**: explain why parameterized queries beat escaping; spot a vulnerable line in a code review.

### 8.3 — XSS and CSRF (less backend, but you'll touch them)
- **Prereqs**: 1.7
- **Sources**: OWASP Cheat Sheets, **Content Security Policy** docs
- **Concepts**: stored / reflected / DOM XSS, output encoding, CSP, CSRF tokens, SameSite cookies as CSRF defense
- **Understands when they can**: pick the right defense per attack class; defend SameSite=Lax as default.

### 8.4 — Secrets management
- **Prereqs**: 9.x (12-Factor)
- **Sources**: **12-Factor App III (Config)**, AWS Secrets Manager / HashiCorp Vault docs
- **Concepts**: secrets in env vars (and the shell-history risk), secret managers, dynamic secrets, rotation, never in code / git history, scanning (gitleaks, trufflehog), encryption at rest with KMS
- **Understands when they can**: design a secret rotation flow that doesn't restart the service unnecessarily; argue against committing secrets to .env files.

### 8.5 — Encryption at rest and in transit
- **Prereqs**: 0.1
- **Sources**: AWS KMS docs, Postgres TDE docs, **Cryptography Engineering** (Ferguson/Schneier/Kohno)
- **Concepts**: TLS for in-transit (already covered in T0), at-rest encryption (full-disk vs column-level), envelope encryption, key rotation, KMS / HSM
- **Understands when they can**: explain envelope encryption (data key + master key); pick column-level encryption for PII fields.

### 8.6 — Authorization models
- **Prereqs**: 1.7
- **Sources**: **Google Zanzibar paper**, **AuthZed / SpiceDB docs**, OWASP authz cheatsheet
- **Concepts**: RBAC (roles), ABAC (attributes), ReBAC (relationships, à la Zanzibar), policy enforcement points, the "check at every layer" principle, fail-closed vs fail-open
- **Understands when they can**: pick RBAC vs ABAC by the access pattern; explain the BOLA defense ("authorize on the object, not just the route").

### 8.7 — DDoS and abuse defense
- **Prereqs**: 1.6
- **Sources**: Cloudflare's DDoS reports, Stripe's bot-defense posts, **OWASP API4 (Unrestricted Resource Consumption)**
- **Concepts**: rate limiting (revisit), CAPTCHAs (and their cost), proof-of-work, IP reputation, anomaly detection, the cost-of-attack vs cost-of-defense framing
- **Understands when they can**: design a defense in depth for a public signup endpoint; calculate cost-of-attack vs cost-of-defense.

### 8.8 — Dependency security
- **Prereqs**: 9.x (build)
- **Sources**: **GitHub's Dependabot docs**, **OWASP Dependency-Check**, Sigstore / SLSA framework, **the SolarWinds and Log4Shell incidents**
- **Concepts**: lockfiles, SCA scanning, vulnerability databases (NVD, GHSA, OSV), reproducible builds, signed artifacts, supply-chain attacks (typosquatting, dependency confusion), pinning vs floating versions
- **Understands when they can**: defend a dependency-pinning policy; identify a typosquat in a package list.

---

## Tier 9 — DevOps adjacency

Backend engineers don't have to be SREs, but they have to *speak* deploy. Treated as prerequisite, not deep ops.

### 9.1 — The 12-Factor App
- **Prereqs**: none
- **Sources**: **12factor.net** — read all 12 factors
- **Concepts**: codebase, dependencies, config, backing services, build/release/run, processes (stateless), port binding, concurrency, disposability, dev/prod parity, logs, admin processes
- **Understands when they can**: review a Dockerfile + deploy config and list 12-Factor violations.

### 9.2 — Containers
- **Prereqs**: 9.1
- **Sources**: **Docker docs**, **"Best practices for writing Dockerfiles"** (official), Google's distroless images
- **Concepts**: image layers, Dockerfile semantics (caching, ordering of COPY/RUN), multi-stage builds, base image choice (alpine vs distroless vs slim vs full), running as non-root, build arg vs env vs runtime, image scanning
- **Understands when they can**: write a multi-stage Dockerfile that produces a minimal image; explain why running as root inside a container is still bad.

### 9.3 — CI/CD
- **Prereqs**: 9.2
- **Sources**: **Software Engineering at Google** ch. 23-24, GitHub Actions / CircleCI / GitLab CI docs
- **Concepts**: build / test / release / deploy stages, artifact promotion vs rebuild, pipeline-as-code, secrets in CI, branch protection, deploy keys vs OIDC, trunk-based vs gitflow
- **Understands when they can**: design a CI pipeline that promotes a single artifact through environments; defend trunk-based development for backend services.

### 9.4 — Environment configuration
- **Prereqs**: 9.1
- **Sources**: **12-Factor III (Config)**, **Confluent's "Configuration as Code"** posts
- **Concepts**: env vars vs config files vs config service, validation at boot (fail fast), per-env overrides, feature flags vs config, the "commit your prod config to git" risk
- **Understands when they can**: design a config-loading layer that validates at startup and fails fast on invalid combinations.

### 9.5 — Infrastructure as Code (literacy)
- **Prereqs**: 9.2
- **Sources**: **Terraform docs**, **Pulumi docs**
- **Concepts**: declarative vs imperative, state file (and why it's load-bearing), plan vs apply, modules, drift, the "click in console once, capture later" anti-pattern
- **Understands when they can**: read a Terraform file and explain what it does; identify drift and how to recover.

### 9.6 — Kubernetes literacy
- **Prereqs**: 9.2
- **Sources**: **Kubernetes Up & Running**, official K8s docs
- **Concepts**: Pod / Deployment / Service / Ingress, ConfigMap / Secret, Namespace, resource requests vs limits, HPA, liveness/readiness/startup probes (revisit T5), why the control plane has to be Highly Available
- **Understands when they can**: deploy a simple service to K8s and explain each manifest; pick resource requests sensibly. **Note**: deep K8s ops belongs to a dedicated SRE skill — backend engineers stop at "I can read and modify our manifests."

---

## Tier 10 — Cloud literacy

What you reach for, why, and what it costs. Cloud-agnostic concepts; AWS as the default reference because of market coverage.

### 10.1 — Compute models
- **Prereqs**: 9.2
- **Sources**: AWS / GCP / Azure docs (compare), **"Serverless in the Wild"** paper (Microsoft)
- **Concepts**: VMs (EC2 / Compute Engine) vs containers-on-VMs (ECS/EKS, GKE) vs FaaS (Lambda / Cloud Functions / Cloud Run), cold start, scaling shape, billing granularity (per-second vs per-100ms), when serverless wins (bursty, low-traffic, event-driven) and when it loses (sustained, high-throughput, latency-sensitive)
- **Understands when they can**: pick compute model by workload shape; estimate cost for a hypothetical bursty workload across the three.

### 10.2 — Storage primitives
- **Prereqs**: 2.0
- **Sources**: AWS S3 / EBS / EFS docs, **"Building S3" essays** by Andy Warfield
- **Concepts**: object (S3, blob storage), block (EBS, persistent disks), file (EFS, FSx), storage classes (Standard / IA / Glacier), eventual consistency in S3 (now strongly consistent — but understand the history), egress cost
- **Understands when they can**: pick object vs block vs file; predict egress cost for a hypothetical multi-region copy.

### 10.3 — Managed services tour
- **Prereqs**: 2.0, 3.6
- **Sources**: AWS / GCP / Azure marketing pages, but read past them; **AWS re:Invent talks** as the real spec
- **Concepts**: managed Postgres / MySQL (RDS, Aurora, Cloud SQL), managed Mongo (DocumentDB — *not Mongo-compatible enough*; MongoDB Atlas — actual Mongo), DynamoDB / Cosmos / Spanner, queues (SQS, Pub/Sub, Service Bus), streams (Kinesis, Pub/Sub Lite, Event Hubs)
- **Understands when they can**: pick managed Postgres vs DynamoDB by access pattern, not religion.

### 10.4 — IAM
- **Prereqs**: 8.6
- **Sources**: AWS IAM docs, **"AWS IAM Policy Evaluation Logic"** doc, Google IAM docs
- **Concepts**: principals, resources, actions, policies (identity-based vs resource-based), roles vs users, AssumeRole / OIDC federation, least privilege, the "wildcard everything" trap, instance / pod identity (IRSA, Workload Identity)
- **Understands when they can**: write a policy granting least-privilege access for a specific use case; explain why instance-role chaining beats hardcoded keys.

### 10.5 — Cost shape and lock-in
- **Prereqs**: 10.1, 10.2, 10.3
- **Sources**: **"The Cost of Cloud, a Trillion Dollar Paradox"** (a16z), **Basecamp / 37signals exit-from-cloud posts** (the other side)
- **Concepts**: on-demand vs reserved vs spot, egress as the silent killer, idle waste, S3 storage class lifecycle, multi-cloud as marketing, lock-in risk per service
- **Understands when they can**: estimate monthly cost for a hypothetical service across compute / storage / egress; defend a managed-service choice with a cost number.

---

## Tier 11 — Distributed systems (implementation-level)

Cross-link with `system-design-tutor` — that skill owns *architecture-level reasoning at scale* (designing globally distributed systems from scratch). This tier owns the *code-it-and-break-it* version: set up replication, kill nodes, observe behavior.

### 11.1 — Replication: single-leader
- **Prereqs**: 2.3, 2.8
- **Sources**: **DDIA** ch. 5 (still the canonical text), Postgres docs on streaming replication
- **Concepts**: sync vs async replication, replication lag, read-your-writes, monotonic reads, follower failover, the split-brain scenario, manual vs automatic promotion
- **Understands when they can**: walk through what async replication does to read-your-writes and how to mitigate (sticky reads, version tokens); set up a Postgres primary + replica locally and induce failover.

### 11.2 — Replication: multi-leader and leaderless
- **Prereqs**: 11.1
- **Sources**: **DDIA** ch. 5, **Cassandra docs**, Riak docs
- **Concepts**: multi-leader (and conflict resolution: LWW, CRDTs, application-level merge), leaderless (Dynamo-style), quorum reads/writes (W + R > N), read repair, anti-entropy / Merkle trees
- **Understands when they can**: explain why Cassandra's W=R=ALL gives you stronger reads but kills availability; reason about CRDTs vs LWW for a counter use case.

### 11.3 — Sharding / partitioning
- **Prereqs**: 11.1
- **Sources**: **DDIA** ch. 6, **Discord's sharding post**, **Notion's sharding post**, Vitess docs
- **Concepts**: hash partitioning, range partitioning, consistent hashing, virtual nodes, hot keys / hot shards, cross-shard queries (and why they kill you), resharding
- **Understands when they can**: pick partitioning strategy by workload; explain why consistent hashing matters for elasticity; defend "design your sharding key like it's permanent — because it is."

### 11.4 — Consensus (concept-level only)
- **Prereqs**: 11.1
- **Sources**: **The Raft paper** (Ongaro & Ousterhout — read it), **Raft visualization at thesecretlivesofdata.com**, etcd docs
- **Concepts**: leader election, log replication, safety vs liveness, FLP impossibility, Raft state machine, why Paxos exists but is rarely implemented directly, why you almost never write your own consensus
- **Understands when they can**: explain Raft's leader election in 5 sentences; reason about what etcd / Zookeeper actually do for you.

### 11.5 — Distributed locking and coordination
- **Prereqs**: 3.3, 11.4
- **Sources**: **Martin Kleppmann's Redlock critique**, etcd / Zookeeper docs, **Chubby paper** (Google)
- **Concepts**: distributed lock semantics (and why they're harder than they look), fencing tokens, leader leases, the network-partition trap, Redlock's safety claims and their critics
- **Understands when they can**: explain why a TTL-based distributed lock without fencing tokens is unsafe; pick etcd over Redis for "absolutely must not double-process this."

### 11.6 — Time and clocks
- **Prereqs**: 3.2
- **Sources**: **DDIA** ch. 8, **TrueTime paper** (Spanner), Lamport's "Time, Clocks, and the Ordering of Events"
- **Concepts**: wall clocks (NTP, drift, leap seconds), monotonic clocks (always use them for durations), logical clocks (Lamport, vector), hybrid logical clocks (HLC), TrueTime
- **Understands when they can**: explain why `time.Now()` for a duration measurement is a bug; identify NTP-step incidents in a postmortem.

### 11.7 — Service discovery and load balancing
- **Prereqs**: 0.2, 11.4
- **Sources**: **Building Microservices 2e** ch. 5, Envoy / Istio docs, Consul docs
- **Concepts**: client-side vs server-side LB, DNS-based discovery, registry-based (Consul, etcd, Eureka), service mesh (Envoy as data plane, Istio/Linkerd as control plane), gRPC's built-in client LB
- **Understands when they can**: pick client-side vs server-side LB by use case; explain what a service mesh adds and what it costs.

---

## Path: Real-time systems (cross-cutting)

Not a tier — a *path* assembled from T0/T1/T3 plus a dedicated builder loop. Use when the learner has an explicit real-time goal (live chat, collaborative editing, presence, notifications).

**Sequence:**
1. T0.6 (WebSockets) + T0.7 (SSE) + T1.12 (real-time API patterns) — protocols
2. T3.1 (concurrency models) + T3.6 (queues) — concurrency for fan-out
3. **Real-time-specific topics** (covered as a sub-path):
   - **Presence systems** — heartbeat / ping, last-seen, scaling presence
   - **Broadcast and fan-out** — pub/sub layer between server instances (Redis pub/sub, NATS, Kafka), why a single-instance WS server doesn't survive a second instance
   - **Sticky sessions** — when to use them, when load balancers fight you
   - **Backpressure** — what happens when the slow consumer can't keep up; explicit window vs implicit drop
   - **Message ordering and replay** — sequence numbers, Last-Event-ID, the resync protocol
4. **Builder-first dedicated loop** (in `references/builder-first.md`) — Loop 8: real-time chat / live cursor.
- **Sources**: **Discord engineering** (presence at scale), **Figma multiplayer**, **Phoenix Channels** docs, **Slack RTM**, **Trevor Sullivan's "Building Real-Time Apps with Elixir"** (Pragmatic Programmers).

---

## Source-to-tier index

| Source | Tiers it anchors |
|---|---|
| RFC 9110 (HTTP semantics) | T0, T1 |
| RFC 8446 (TLS 1.3), RFC 9000 (QUIC) | T0 |
| **API Design Patterns** (Geewax, 2021) | T1 |
| **Stripe engineering blog** | T1 (idempotency, rate limiting, API design), T8 |
| **Cloudflare engineering blog** | T0, T4, T8 (DDoS) |
| **GitHub engineering blog** | T2 (online migrations), T6 (incidents), T8 |
| **Discord / Figma / Notion / Shopify engineering** | T0 (real-time), T2 (sharding), T4 (caching), T11 |
| **Database Internals** (Petrov, 2019) | T2 |
| **Use the Index, Luke!** (Winand) | T2, T7 |
| **DDIA** (Kleppmann, 2017) — "fresh" caveat: still the canonical text where no fresher equivalent exists | T2, T11 |
| **Building Microservices 2e** (Newman, 2021) | T1 (gRPC), T3 (saga), T11 (service discovery) |
| **Microservices Patterns** (Richardson, 2018) | T3 (saga, outbox) |
| **Software Engineering at Google** (2020) | T1 (contract testing), T9 (CI/CD) |
| **The 12-Factor App** | T9 (entire), T6 (logs), T8 (config/secrets) |
| **OWASP API Security Top 10 (2023)** | T8 (entire spine) |
| **The SRE Book** + **The SRE Workbook** (Google, free) | T5, T6 (entire) |
| **"The Tail at Scale"** (Dean & Barroso, 2013) | T5 (hedged), T7 |
| **Marc Brooker's blog** | T3 (backoff), T5 (timeouts/retries) |
| **Release It!** (Nygard, 2nd ed) | T5 (circuit breakers, bulkheads) |
| **OpenTelemetry docs** | T6 (tracing, metrics) |
| **The Raft paper** | T11 |
| **The Outbox pattern** (microservices.io) | T3 |
| **Phil Eaton / Hillel Wayne / Aphyr (Jepsen)** for distributed-systems rigor | T2 (isolation), T11 |
