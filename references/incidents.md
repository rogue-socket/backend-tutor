# Incidents

Real, public, well-documented backend incidents. Use them as **opening hooks, mid-lesson anchors, or reproduce-an-incident exercises** (see `practical-mode.md` type D). A topic without a war story is forgettable.

**Discipline:** never fabricate specifics. If you can't remember the details cleanly, point at the public postmortem URL and let the learner read it.

Format per entry:
- **Name + year + public RCA URL fragment**
- **What happened** (2-3 sentences)
- **What to teach with it** (which tier / topic it anchors)

---

## T0 — Networking & HTTP

### Cloudflare WAF outage — 2019-07-02
**RCA:** cloudflare.com/blog (search "Details of the Cloudflare outage on July 2, 2019")
**What happened:** A regex deployed to the global WAF rule set contained catastrophic backtracking (`(?:(?:\"|'|\]|\}|\\|\d|(?:nan|infinity|true|false|null|undefined|symbol|math)|\`|\-|\+)+[)]*;?((?:\s|-|~|!|{}|\|\||\+)*.*(?:.*=.*)))`). CPUs across the global edge spiked to 100%; ~27 minutes of HTTP 502s for a large fraction of the internet.
**Teach with it:** input validation, why regex on hot paths needs explicit complexity bounds, the value of canary deploys for global config changes (T0 + T5.7).

### Cloudflare BGP / dashboard outage — 2020-07-17
**RCA:** cloudflare.com/blog (search "July 17 Cloudflare outage")
**What happened:** A bad backbone configuration push misrouted traffic, blackholing ~50% of Cloudflare's global network for ~27 minutes.
**Teach with it:** config-push as a top failure mode; why dashboards depend on the same network they're supposed to manage (the "can't deploy a fix because the deploy system is also down" trap). T0 + T9.4.

### Slack outage — 2021-01-04
**RCA:** slack.engineering "Slack's Outage on January 4th 2021"
**What happened:** Cascading failure during return-from-holiday traffic surge. AWS Transit Gateway scaling lag combined with internal service discovery (Consul) issues.
**Teach with it:** capacity planning for known traffic events; what "warm-up" means for cloud infra. T7.4 + T10.

---

## T1 — APIs

### GitHub 24-hour incident — 2018-10-21
**RCA:** github.blog "October 21 post-incident analysis"
**What happened:** A 43-second cross-coast network partition caused MySQL Orchestrator to fail over a primary. When the network healed, two clusters thought they were primary; reconciling them required manually replaying ~40 minutes of writes. The user-visible incident lasted 24 hours.
**Teach with it:** split-brain in failover, why "automatic failover" is harder than it looks, the cost of write divergence on a multi-region cluster. T11.1, T1 (impact of API unavailability), T6.7 (postmortem quality).

### Stripe API design (positive case study)
**Source:** stripe.com/blog (idempotency, API versioning, error envelope posts)
**What happened:** Not an incident — a sustained design discipline. Idempotency keys, date-based versioning per account, structured error envelope, careful rate limiting.
**Teach with it:** what *good* looks like in T1.1 (REST design), T1.2 (idempotency), T1.3 (versioning), T1.5 (errors).

### Idempotency-failure pattern (across payments / webhooks)
**Sources:** stripe.com/blog "Designing robust and predictable APIs with idempotency"; AWS SQS dev guide (visibility timeout / at-least-once); GitHub webhook delivery docs; Stripe webhook signature & retry docs.
**What happened:** Few public RCAs name idempotency-failure as the root cause — most teams ship this bug, catch it in staging or absorb it as user-visible weirdness, and don't write a public postmortem. The pattern, documented across canonical guidance: client retries a non-idempotent POST (network blip, lambda timeout, default retry middleware), server treats it as a new request, duplicate side-effect lands — double-charge, duplicate notification, repeated inventory decrement. Variant: server keys the dedup table on the idempotency key alone, ignoring a body hash; a buggy or malicious client reuses the key with a *different* payload, server returns the *cached* prior response while silently accepting the new write — worst-of-both-worlds. Stripe's design treats the idempotency-key row as a state machine (`started → succeeded/failed`) so a crash mid-flight can be resolved on retry without producing a second side-effect.
**Teach with it:** the Stripe blog is the canonical *positive* anchor for T1.2 (idempotency keys server-side), T3.5 (idempotency in handlers under crash), T1.13 (webhooks — receiver-side dedup, signature verification with timing-safe comparisons). Pair with the Knight Capital entry below for the broader "duplicate execution destroys you" frame.

---

## T2 — Databases

### GitLab database deletion — 2017-01-31
**RCA:** about.gitlab.com/blog "Postmortem of database outage of January 31"
**What happened:** A tired engineer ran `rm -rf` on a primary replica thinking it was a secondary. Then discovered all five backup methods were broken: pg_dump didn't run (version mismatch), Azure disk snapshots disabled, S3 backups not configured, LVM snapshots only daily, replication lag rendered the replica useless. They restored from a 6-hour-old staging copy.
**Teach with it:** untested backups don't exist; production access controls; the "tired ops engineer at midnight" failure mode. T2.12 (backups & DR).

### Discord — Trillions of messages on ScyllaDB
**Source:** discord.com/blog "How Discord stores trillions of messages"
**What happened:** Discord migrated from Cassandra to ScyllaDB after hitting hot-partition pain (single channels with millions of messages causing GC stalls). The ScyllaDB rewrite eliminated tail-latency spikes.
**Teach with it:** picking a database for a specific workload; what "hot partition" actually feels like in prod. T2.0, T11.3.

### AWS DynamoDB — 2015-09-20 outage
**RCA:** aws.amazon.com/message/5467D2/
**What happened:** Metadata service overload caused a 5-hour DynamoDB outage in us-east-1, which cascaded to ~20+ AWS services that depend on it.
**Teach with it:** metadata services as the silent dependency; why "managed = no ops" is a lie. T10.3.

### Bad-index-in-production pattern (DDL locks / write amplification)
**Sources:** strongmigrations.com (Andrew Kane / Ankane — Rails-flavored but rules apply everywhere); github.com/github/gh-ost (MySQL online schema-change tool, written *because* of this failure mode); Postgres docs on `CREATE INDEX CONCURRENTLY`; Use the Index, Luke! (Markus Winand).
**What happened:** Two recurring failure shapes, both extensively documented in public engineering practice, neither typically getting a single named-incident RCA: **(1) Non-CONCURRENTLY index creation**: `CREATE INDEX` on Postgres without `CONCURRENTLY` (or any DDL on MySQL without an online-DDL tool) takes an `ACCESS EXCLUSIVE` lock on the table for the duration of the build — minutes-to-hours on large tables. Engineers have stalled checkout, login, and feed-write paths this way. The Strong Migrations gem refuses to run a non-concurrent index migration in CI for this reason; gh-ost exists because the same problem is industry-wide on MySQL. **(2) Write amplification from too many indexes**: every additional index on a write-heavy table turns each insert/update into N+1 disk writes plus N index-page splits. Tables with 10+ indexes routinely run inserts an order of magnitude slower than expected; batch jobs that used to take 10 minutes start taking hours, and the team blames the application before realizing the index list grew silently over years.
**Teach with it:** T2.2 (indexes and the planner — every index is a trade); T2.5 (online migrations — `CONCURRENTLY` is non-negotiable, expand/contract for `NOT NULL` adds, gh-ost-style for MySQL); T2.6 (ORM-generated index pathologies). The existence of gh-ost / pg_repack / Strong Migrations is the *evidence* the failure mode is widespread enough to justify dedicated tooling.

---

## T3 — Concurrency & async

### AWS Kinesis outage — 2020-11-25
**RCA:** aws.amazon.com/message/11201/
**What happened:** Adding capacity to the Kinesis front-end fleet caused servers to exceed the OS thread limit. Each new server tried to open a thread per peer for cluster membership, total threads scaled O(n²), hit the limit, fleet couldn't form a cluster.
**Teach with it:** why OS / runtime limits surface unexpectedly when you scale up; the difference between thread-per-connection and event-loop architectures. T3.1, T7.4.

### Knight Capital — 2012-08-01
**Sources:** SEC filing, public postmortems
**What happened:** Deployed new code to 7 of 8 servers; the 8th still had old code that reused a previously deprecated feature flag. The two code paths fought; Knight executed $7B in errant trades in 45 minutes, lost $440M, the company never recovered.
**Teach with it:** deploy hygiene, feature flag debt, the "one out of N" tail bug. T5.6, T5.7, T9.3.

### Queue-redelivery / non-idempotent consumer pattern
**Sources:** AWS SQS dev guide (visibility timeout, at-least-once, FIFO vs standard); AWS Lambda + SQS + DLQ docs; Stripe / GitHub / Shopify webhook delivery docs; Kafka docs on consumer rebalance + offset commit semantics; Confluent's "Exactly-Once Semantics" blog (read alongside Tyler Treat's "You Cannot Have Exactly-Once Delivery" critique).
**What happened:** Public RCAs naming "queue redelivery to a non-idempotent consumer" as the root cause are rare — most teams that ship this bug either catch it in staging or absorb it as user-visible weirdness without writing a postmortem. The pattern, documented across every queue system's docs: **(SQS variant)** visibility timeout expires before the consumer ack's (worker GC pause, slow downstream call, oversize message rerouted to a slow path), broker re-delivers, a second consumer processes the same message → duplicate charge / duplicate email / duplicate inventory decrement. **(Kafka variant)** consumer rebalance during a deploy or a failed liveness probe → uncommitted offsets re-delivered to the new consumer in the group; same duplicate side-effect. **(Webhook variant)** receiver returns 5xx or times out at the load balancer; sender retries per its retry policy (Stripe, GitHub, Shopify all retry); receiver processes the same event twice. The mitigation is universal: idempotent consumers via a dedup table keyed on a *stable* message ID (the application-domain ID, not the broker-assigned receipt handle, which changes across redeliveries).
**Teach with it:** T3.6 (queues — partition / consumer-group semantics), T3.7 (delivery semantics — at-least-once is the default; "exactly-once" is a marketing term that costs you throughput and still requires idempotent consumers downstream), T3.5 (idempotent handlers as the universal queue safety net). Pair with the AWS Kinesis 2020 entry to emphasize that the queue cluster itself can also fail in cluster-formation ways — idempotent consumers are necessary but not sufficient.

---

## T4 — Caching

### Discord — cache stampede on a hot key
**Source:** various Discord engineering talks / blog posts
**What happened:** A popular channel's cache key expired during a high-traffic event; thousands of concurrent requests slammed the DB to recompute, the DB couldn't keep up, the cache backfill stalled, the stampede self-perpetuated.
**Teach with it:** the cache-stampede pattern; single-flight / probabilistic early refresh as the cure. T4.3.

### Facebook 2021-10-04 (BGP / DNS cascade — partly cache-relevant)
**RCA:** engineering.fb.com "More details about the October 4 outage"
**What happened:** A BGP withdrawal from Facebook's backbone took the authoritative DNS servers off the internet. DNS resolvers globally started thrashing, causing a thundering-herd retry surge that hit recovery hard. Internal services that used DNS for service discovery couldn't talk to each other; engineers couldn't badge into the buildings.
**Teach with it:** DNS cache TTL behavior under origin failure; "negative caching" of DNS failures; physical-world dependencies (badge readers). T0.2, T4.2, T5.5.

---

## T5 — Reliability

### Robinhood — 2020-03-02 / 03 GME-adjacent capacity
**RCA:** robinhood.engineering blog
**What happened:** Robinhood hit a thundering-herd surge from a market event and ran out of compute / DB capacity, taking the trading platform offline during market hours.
**Teach with it:** capacity planning for known-bursty events; what "scale" means for a regulated workload. T7.4, T5.

### Facebook — 2021-10-04 (revisit)
**Teach with it:** the "deploy a fix when the deploy system is also down" trap. T5.7.

### Roblox — 2021-10-28 to 31 (73-hour outage)
**RCA:** blog.roblox.com "Roblox Return to Service 10/28-10/31 2021"
**What happened:** Consul cluster overwhelmed by streaming workload from a new feature. Recovery from Consul state was slow (hours per attempt) because of contention on Consul's KV store. Compounded by a separate bug in Consul itself. 73 hours.
**Teach with it:** dependency criticality; recovery time as a design constraint; the limits of off-the-shelf coordination services at scale. T11.4, T11.5, T6.7.

---

## T6 — Observability & on-call

### Honeycomb's "production excellence" stories
**Source:** honeycomb.io blog
**What happened:** Various — Honeycomb has consistently good public writing on tracing-driven debugging.
**Teach with it:** what tracing gives you that logs and metrics don't (T6.4); span attribute design.

### Cloudflare 2019 WAF (revisit)
**Teach with it:** why CPU saturation alarms must be page-worthy; what a *good* postmortem reads like (Cloudflare's RCA is a model). T6.6, T6.7.

---

## T7 — Performance & scale

### "The Tail at Scale" — Dean & Barroso, 2013
**Source:** Communications of the ACM
**What happened:** Not an incident — the foundational paper on why p99 latency matters at scale and how fan-out amplifies it.
**Teach with it:** tail latency, hedged requests, the math behind "100 sequential calls means p99 ≈ p99.99 of one call." T7.1, T5.4.

### Discord — Cassandra → ScyllaDB (revisit)
**Teach with it:** GC pause as a tail-latency source; why language-runtime choice is a perf decision. T7.1, T7.2.

---

## T8 — Security

### Capital One — 2019-07-19 (Paige Thompson)
**Source:** US DOJ filings; AWS public statements; Krebs on Security analysis
**What happened:** Misconfigured WAF on Capital One's AWS-hosted application allowed an attacker to perform an SSRF attack, retrieving IAM credentials from the EC2 instance metadata service (IMDSv1). Credentials had `s3:ListBucket` and `s3:GetObject` over the bank's data buckets. ~100M customer records exfiltrated.
**Teach with it:** **OWASP API7 (SSRF) — the canonical example.** Also T8.6 (least-privilege IAM); T10.4 (instance metadata). The mitigation: IMDSv2 (session-token-required), egress filtering, scoped IAM roles.

### Equifax — 2017 (Apache Struts CVE-2017-5638)
**Source:** GAO report; public postmortems
**What happened:** Unpatched Apache Struts vulnerability exploited via crafted Content-Type header in a file-upload form. ~147M people's PII exfiltrated over ~10 weeks before detection.
**Teach with it:** dependency security (T8.8); detection latency; why patch-management is engineering work, not IT work.

### Log4Shell — 2021-12-09 (CVE-2021-44228)
**Source:** apache.org / NVD / countless writeups
**What happened:** Log4j 2.x evaluated `${jndi:ldap://...}` strings inside log messages. Attackers controlling any logged input — User-Agent, header values, form fields — could trigger remote code execution.
**Teach with it:** input handling at every layer (logs are user input!); transitive dependencies; the "zero-day → patch in 24h" exercise. T8.2, T8.8.

### SolarWinds Orion — 2020
**Source:** CISA report; FireEye disclosure
**What happened:** Attacker compromised SolarWinds' build system and inserted a backdoor into the Orion update artifact, which was code-signed and shipped to ~18,000 customers including US gov agencies.
**Teach with it:** supply-chain security (T8.8); reproducible builds; SLSA / Sigstore. Why build-system security is on the critical path.

---

## T9 — DevOps adjacency

### Knight Capital (revisit)
**Teach with it:** deploy automation that doesn't actually verify all hosts converged. T9.3.

### GitLab 2017 (revisit)
**Teach with it:** environment parity (12-Factor X); the "I thought this was staging" failure mode. T9.1.

---

## T11 — Distributed systems

### GitHub 2018 (revisit)
**Teach with it:** automatic failover gone wrong; manual replay required. T11.1, T11.5 (consensus / leader leases).

### Roblox 2021 (revisit)
**Teach with it:** Consul / etcd / Zookeeper as load-bearing infrastructure; recovery from coordination-service failure is itself a coordination problem. T11.4, T11.5.

### AWS S3 — 2017-02-28
**RCA:** aws.amazon.com/message/41926/
**What happened:** Engineer running a debug command typo'd a parameter and accidentally took down a larger set of S3 servers than intended in us-east-1. Restart took 4 hours because the S3 metadata subsystem hadn't been restarted in years and was slow to come up cold.
**Teach with it:** the "we haven't tested cold starts" failure mode; why production access tools must be safer than `--all`. T2.12, T9.

---

## How to use incidents in lessons

- **As an opening hook:** "In 2019 Cloudflare took half the internet down with a regex. Today's lesson is about why." (T0)
- **As a mid-lesson anchor:** explain the concept, then "this is what happened to GitHub in 2018." (T11)
- **As a reproduce-an-incident exercise:** see `practical-mode.md` type D.
- **As a postmortem-writing exercise:** hand the learner the timeline, ask them to write the "contributing factors" section before reading the published one.

**Don't:** invent specifics. If you'd be guessing about a number, name, or sequence, say so or skip the detail.
