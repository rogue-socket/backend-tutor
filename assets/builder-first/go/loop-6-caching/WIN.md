# WIN — Loop 6 done

## Code

- [ ] `GET /links/popular` cached in Redis with cache-aside
- [ ] Single-flight wraps the loader; concurrent misses produce 1 DB query, not N
- [ ] Invalidation on write (PATCH/DELETE on links) — choice documented in NOTES.md
- [ ] Hit / miss counts visible in logs or metrics

## Verification

- [ ] BREAK 1 reproduced: stampede before single-flight; ≤1 query per TTL after
- [ ] BREAK 2 reproduced: staleness window without invalidation; bounded after fix

## Understanding

1. **Single-flight prevents N concurrent loads on cache miss. What's the cost?**
   *Outline: latency for N-1 of those concurrent callers — they wait for the one in-flight loader. Trade-off: marginal latency cost for the waiters vs catastrophic backend load if you don't. Worth it almost universally. The exception is when the loader can fail or hang — then a stuck loader keeps everyone waiting; mitigation is per-call timeouts inside the loader.*

2. **Cache-aside has a known race: a slow read in flight can write stale data over fresh. Walk through it.**
   *Outline: T0 reader misses, queries DB. T1 writer updates DB, invalidates cache (no-op, empty). T2 reader's query returns pre-T1 data. T3 reader writes it to cache. Cache is now stale with no invalidation pending. Mitigations: keep cache TTL short relative to write rate; use SET NX so slow writes don't overwrite; or read the data version and only write if version matches.*

3. **Why is "version bump" sometimes preferred over `cache.Del`?**
   *Outline: in a multi-instance / multi-region cache, `Del` may not reach all replicas instantly. Bumping the version key is a single atomic write; old keys remain in cache but are unreferenced and naturally expire. Trade-off: old keys consume memory until TTL. Worth it when invalidation must be effectively instant across many readers; not worth it when memory matters.*

## Reflection

What surprised you? Common ones:
- Stampedes happen at *every* TTL boundary, not just the first
- Single-flight is so simple it's almost embarrassing
- "Invalidate on write" sounds simple but the race is real

## What's next

Loop 7 — Containerize + deploy. Multi-stage Dockerfile, full-stack docker-compose, graceful shutdown, real cloud deploy.
