# Loop 6 — Caching + stampede

**Tier mapping:** T4 (entire)
**Time:** 90 minutes
**Theme:** *cache invalidation is hard for a reason; cache stampedes are why.*
**Prereqs:** Loops 2 + 5 done (you have Postgres + Redis).

## What you're building

Add a cache to a hot read endpoint. The "popular links" endpoint is the obvious candidate — `GET /links/popular` returns the top 50 most-viewed links across the system. The query is moderately expensive (a few joins, an aggregation); the result is the same for every user; cacheable.

Cache-aside pattern with a 60s TTL. Then **deliberately reproduce a stampede**, then fix it with single-flight.

## Architecture

```
GET /links/popular
    ↓
cache.Get("popular:v1")
    ├── hit → return cached
    └── miss → singleflight.Do("popular:v1", func() {
                    rows := db.Query(...)
                    cache.Set("popular:v1", rows, 60s)
                    return rows
               })
```

`singleflight.Do` ensures that for any given key, only one goroutine actually runs the loader function; concurrent callers wait for that result. This is the load-bearing primitive.

## New code

- `internal/cache/cache.go` — Redis Get/Set wrapper with JSON marshaling
- `internal/cache/popular.go` — the cached "popular" query, using `golang.org/x/sync/singleflight`
- (Add to existing handlers) — `GET /links/popular`

## Tasks

1. **Implement `cache.Get` / `cache.Set`** with JSON values + a `cache_version` prefix so you can flush all keys with one Redis FLUSHDB-equivalent (don't actually FLUSHDB — bump the version).
2. **Implement `popular.GetOrCompute`** — cache-aside with single-flight wrapped around the DB query.
3. **Wire `GET /links/popular`** to call it.
4. **Run BREAK 1** — the stampede. Disable single-flight; observe.
5. **Run BREAK 2** — the staleness. Update a link directly in the DB; observe stale cache.
6. **Implement invalidation on write** — when a user PATCHes or DELETEs a link, invalidate the popular cache (or bump the version).
7. **Verify against `WIN.md`.**

## The single-flight library

`golang.org/x/sync/singleflight` is the canonical Go implementation. It's tiny (~100 lines). Read the source — understanding it deeply is worth more than memorising any specific cache library.

```go
import "golang.org/x/sync/singleflight"

var sfg singleflight.Group

func GetOrCompute(ctx context.Context, ...) ([]Link, error) {
    if cached, ok := cache.Get(key); ok {
        return cached, nil
    }
    val, err, _ := sfg.Do(key, func() (any, error) {
        rows := db.Query(...)
        cache.Set(key, rows, 60*time.Second)
        return rows, nil
    })
    if err != nil { return nil, err }
    return val.([]Link), nil
}
```

The `_` in `Do`'s return is the "shared" bool — true if more than one caller waited for this result. Wire it into a metric; you'll see cache-stampede prevention in action.

## Stretch — probabilistic early refresh

After single-flight, the next failure mode is "every request between TTL expiry and the first refresh sees a miss." Probabilistic Early Refresh (XFetch — Vattani et al.) refreshes proactively before TTL expiry, smoothing the spike.

Implementation: in `GetOrCompute`, when a request reads cached data with remaining TTL near zero, *some fraction* of requests trigger a background refresh. The fraction increases as TTL approaches zero.

```go
shouldRefresh := rand.Float64() < (1 - timeRemaining / cacheDuration) ^ 2
```

Squaring the ratio makes early refreshes rare when fresh, common when stale. Implement, measure stampede behaviour with and without; this is Loop 6's stretch goal.
