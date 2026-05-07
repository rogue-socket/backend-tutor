# BREAK — Loop 10

The whole loop *is* the break. Whatever falls over first under load is your subject.

## Common patterns to expect

The first thing that breaks is rarely the thing you predicted. Some common ones:

### Connection pool exhaustion (likely first)

Your DB pool was sized for dev. Under sustained read load, you hit it. p99 spikes; DB is calm; pool stats show saturation. Loop 2's BREAK 1 at scale.

**Telemetry signature:** `pgxpool.Stat().AcquireWaitDuration()` rising; DB CPU low; application latency rising.

**Fix:** size the pool sensibly. HikariCP rule of thumb: pool size = N * (cores * 2 + spindles), where N is the number of app instances. For Postgres on modern SSD: pool size = 2-3× CPU cores per instance. *Less* than you think.

### Cache stampede (very likely)

Your cache TTL is 60s. Under load, every 60s the popular-links cache expires; concurrent requests hammer the DB. p99 spikes periodically. Loop 6's BREAK 1 at scale.

**Telemetry signature:** sawtooth pattern in DB query rate. Latency spikes precisely synced to TTL boundaries.

**Fix:** single-flight on cache miss. You probably already implemented this in Loop 6, but it might not be wired everywhere it should be.

### N+1 you didn't catch

Loop 2's N+1 was on tags. There's probably another one — auth check on every link, owner lookup, something. Under load, the per-request query count multiplies traffic on the DB.

**Telemetry signature:** DB QPS scales super-linearly with HTTP RPS. Trace shows N child spans per request.

**Fix:** batch the lookup. If it's auth, cache the user in the request context.

### GC pause spikes

Go's GC is good but allocates show up. Heavy JSON marshaling on hot paths can produce GC pauses that show up as p99 latency.

**Telemetry signature:** `go_gc_duration_seconds` correlated with p99 spikes. Allocation rate (`go_memstats_alloc_bytes_per_sec`) high.

**Fix:** sync.Pool for hot allocations, encoding/json → easyjson or sonic for high-throughput, set GOGC if memory headroom allows.

### Single-threaded bottleneck

Some piece of your code is gated by a global lock or a single goroutine. CPU is fine, throughput plateaus.

**Telemetry signature:** RPS plateaus before CPU saturates. Latency rises while RPS stays flat (queueing).

**Fix:** find the lock or the goroutine; partition the state, sharded locks, or work-stealing pool.

### Worker queue backup

Background work falls behind. Queue depth grows. SLO on "notification sent within 1 minute" breaches.

**Telemetry signature:** `queue_depth` gauge climbing. Worker CPU high or worker DB queries slow.

**Fix:** scale the worker pool, batch the work, or increase the worker's per-message efficiency.

## How to find which one

Order of operations:

1. **Look at the dashboards FIRST.** RED on the HTTP service. USE on the resources (pool, CPU, memory, GC).
2. **Trace a slow request.** Which span is the offender?
3. **Read the code only after** telemetry has narrowed the search.

This is the same discipline as Loop 9. Loop 10 is where it pays off.

## What if nothing breaks?

You're under-loading. Push harder. Real services break somewhere — if yours doesn't at 1000 RPS, push to 5000, then 10000. Or your test is closed-loop and hiding tail latency; switch to open-loop arrival rate.

If after honest open-loop testing nothing breaks: **interesting**. You're either over-provisioned or you've already done good design work in Loops 1-9. Capacity-plan the next bottleneck — at what RPS would you expect each layer to fail next? Validate one of those predictions experimentally.

## Postmortem discipline

Whatever you find, write it up. The postmortem is the artifact; the writing is where the learning lives.

The format from `postmortem-template.md` is a starting point. Modify for your situation. The non-negotiables:
- **Timeline with quotes from telemetry.** Memory + screenshot.
- **Multiple contributing factors.** Almost never one cause.
- **Concrete action items.** "Improve observability" is fiction.
- **Honest "what didn't go well."** This is where the learning is.
