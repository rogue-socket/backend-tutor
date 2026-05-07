# BREAK — Loop 6

Two breaks. The first is the canonical cache-failure mode; the second is the canonical cache-correctness mode.

---

## BREAK 1 — The cache stampede

**Foundations filled here:** T4.3 (cache stampedes).

### The setup

You implemented cache-aside *without* single-flight. The cache has a 10-second TTL (shortened for the demo from the production 60s). Run a load test:

```bash
hey -n 1000 -c 100 -z 60s 'http://localhost:8080/links/popular'
# or use any concurrent client
```

### What you'll see

For the first ~10 seconds, almost all requests hit the cache. Latency: sub-ms.

At T=10s (TTL expiry), the cache returns "miss" — every concurrent request runs the DB query. The DB sees a stampede: ~100 simultaneous identical queries. DB latency spikes (or worse, DB connection pool exhausts). Some requests fail with timeouts. Others succeed slowly.

After all the simultaneous queries finish, the cache is repopulated. Sub-ms latency resumes — until the next TTL expiry, which restarts the stampede.

### Reproduce in metrics

Wire up a query counter on the `popular` query. Without single-flight: count goes from 0 → 100+ at every TTL boundary. With single-flight: count goes from 0 → 1 per TTL boundary.

### The fix

`golang.org/x/sync/singleflight`:

```go
val, err, shared := p.sfg.Do(popularKey, func() (any, error) {
    rows, err := p.loader(ctx)
    if err == nil {
        b, _ := json.Marshal(rows)
        p.rdb.Set(ctx, popularKey, b, popularTTL)
    }
    return rows, err
})
```

After this fix, rerun the load test. DB query count at TTL boundary: 1.

### Why this is the lesson

Stampedes are the single most common cache-related outage. Discord, Reddit, Wikipedia, and basically every site with a hot key has had one in production. The pattern is universal:
- Hot key
- TTL expires
- Concurrent miss
- Backend gets hammered

The cure is universal too: only one loader runs at a time per key. Single-flight is the simplest implementation; locks (Redis SET NX EX) work too at the cost of extra round-trips.

---

## BREAK 2 — The staleness window

**Foundations filled here:** T4.2 (cache invalidation).

### The setup

You PATCH a link's title via `PATCH /links/{id}` — your handler updates the DB but doesn't invalidate the popular-links cache. The next request to `GET /links/popular` returns the *old* title, because the cache hasn't expired yet.

### What you'll observe

```
PATCH /links/42 {"title":"new title"}    → 200, DB has 'new title'
GET /links/popular                        → returns the old title for link 42
... wait 60s ...
GET /links/popular                        → returns 'new title'
```

That gap (here, up to 60s) is the **staleness window** — the period during which the cache returns data the DB has already changed.

### Why this matters

For some workloads, 60s of staleness is fine (popular-links rankings, leaderboards, dashboards). For others, it's catastrophic (a user's own link they just edited; a privacy setting toggled to "private" but still visible publicly).

The cure is one of:
1. **Invalidate on write.** PATCH/DELETE handler calls `cache.Del(popularKey)`. Simple; correct; introduces tight coupling between the write path and the cache key list.
2. **Version bump.** Maintain a "version" key (e.g., `cache:popular:version`). Reads include the version in the cache key (`popular:v3`). On write, bump the version. Old keys naturally expire from cache via TTL.
3. **Pub/Sub invalidation.** Writer publishes to a channel; cache subscribers invalidate. Useful at multi-instance scale (Loop 8 territory).
4. **Just-in-time recomputation with short TTL.** Accept the staleness window; tune TTL to the business tolerance.

### The right choice

For a "popular links" cache: **invalidate on write** is overkill (lots of writes; the rankings are already approximate; users don't notice 60s lag in popularity). **Short TTL** is the right call.

For a "user's own links" cache (which you don't have but might add): **invalidate on write** is mandatory.

The pattern: **match the invalidation strategy to the consistency requirement of the data, not to a generic best practice.**

### Implementation note: the timing trap

Even with "invalidate on write," there's a race:
1. T0: request A reads from cache (miss), starts the DB query
2. T1: request B writes new data to DB, then invalidates the cache (no-op, cache is empty)
3. T2: request A's query returns *old* data (stale read from before T1's write)
4. T3: request A writes the (stale) result to the cache
5. The cache now contains stale data, with no signal to invalidate again

This is the *cache-aside race* — and it's load-bearing for staleness. The fix is either:
- Read-your-writes invalidation: re-invalidate after the cache write happens, OR
- Use `set if not exists` semantics so the slow loader doesn't overwrite a fresh value, OR
- Accept it and shorten the TTL

This is one of those things every senior backend engineer has been bitten by.

---

## The takeaway

> A naive cache is an outage waiting for the right traffic spike.

Two non-negotiable habits:
1. **Single-flight on miss** — for any cache, any key. Free correctness, almost-free performance.
2. **Invalidation strategy chosen deliberately** — TTL alone is fine for some workloads, wrong for others. Don't default; decide.
