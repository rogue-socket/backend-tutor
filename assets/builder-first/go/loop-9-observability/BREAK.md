# BREAK — Loop 9

## The main BREAK — the planted bug

The tutor injects a real bug. You debug it from telemetry alone — no `println`, no debugger, no stepping through code. You may read code, but only after the telemetry has narrowed the search space.

The four canonical planted bugs (your tutor will pick one):

### Bug A: phantom 1% slow

A handler has a `if rand.Intn(100) == 0 { time.Sleep(2*time.Second) }`. Hidden in a code path the obvious tests don't cover.

**Telemetry signature:**
- p50 latency: normal
- p99 latency: 2s+
- The histogram buckets show a clear bimodal distribution: most requests under 50ms, a tail at 2s
- Trace search: filter for slow spans → see the exact code path

**The lesson:** averages lie. p50 = healthy, p99 = on fire. RED dashboards must show histograms, not just rates.

### Bug B: wrong-tenant data leakage

Under concurrent load, ~0.1% of `GET /links` returns *another user's* links. Cause: a connection-pooled prepared statement reused across tenants without `WHERE owner_id = $1` properly bound.

**Telemetry signature:**
- Logs: a small fraction of requests where `request.user_id != response.first_link.owner_id`
- Traces: the DB query span shows the bind parameters; spot the wrong tenant
- Metric: a counter you instrument to count "owner_id mismatch" events when you suspect it

**The lesson:** correctness bugs hide in low percentiles too. Telemetry that surfaces "this happened to user X at time T" is worth more than averages.

### Bug C: Sunday slow query

A feature flag flips on for some users → triggers an unindexed query branch. Sunday is when the flag rollout cohort is highest. Slow queries spike then.

**Telemetry signature:**
- `db_query_duration_seconds{query="search_links"}` p99 jumps Sunday morning
- Cross-reference: the query name is logged with feature flags evaluated; correlation appears
- pg_stat_statements shows the query running with seq scans

**The lesson:** observability must be wide enough to catch correlations you didn't predict. Don't tag queries only by name; include the relevant feature-flag context as a span attribute (low-cardinality flag = OK as label, but ideally on traces).

### Bug D: WS hub memory leak

Clients disconnect; the hub doesn't clean up the entry; goroutines accumulate.

**Telemetry signature:**
- `ws_clients_connected` gauge climbs monotonically
- `go_goroutines` (from runtime metrics) climbs in lock-step
- Memory metric trends up

**The lesson:** gauges are how you spot accumulation. Counters tell you rate; gauges tell you state. Both matter.

## Smaller stretch BREAK — cardinality explosion

You add `user_id` as a Prometheus label, thinking "I'd love to see per-user request rates."

```go
httpRequestsTotal.WithLabelValues(method, route, status, userID).Inc()
```

Run a load test with 10K distinct test users. Watch your Prometheus:

```bash
curl -s http://localhost:9090/api/v1/status/tsdb | jq '.data.headStats.numSeries'
```

Series count balloons. Scrape time grows. Memory usage grows. Eventually Prometheus OOMs.

**Why:** every unique combination of label values is a separate time series. `user_id` with 10K values × 5 routes × 3 status codes × 3 methods = 450K series. From this one well-intentioned line.

**Fix:** remove the label. Per-user observability goes in logs and traces (where high-cardinality is the design), not metrics.

**Rule:** label cardinality ≤ ~1000 per label. Total series ≤ low millions per Prometheus instance. Plan accordingly.

## How to debug from telemetry

When the tutor says "the bug is in," resist the urge to read code first. The order:

1. **Look at the dashboards.** RED. Where is the anomaly? Latency? Errors? Saturation?
2. **Drill into the histogram.** Is it bimodal? Heavy tail? Specific buckets growing?
3. **Filter traces.** Find a slow / errored example. Walk the spans.
4. **Read the trace's spans.** Which span is slow? Which is missing?
5. **Now read the code** for that specific span.

The discipline matters. Backend engineers without it default to "I'll just add some `println`" — which produces no signal at scale and trains bad habits.

## The takeaway

> Observability isn't telemetry. Telemetry is a precondition; *being able to answer questions you didn't pre-plan* is observability.

Three rules:
1. **Histograms over averages** for latency. p99 hides under p50.
2. **High-cardinality data in logs and traces, low-cardinality in metrics.** Cardinality is the silent cost.
3. **Correlate across signals.** Logs say "what happened"; metrics say "is the system healthy"; traces say "what's *this specific request* doing." Use all three.
