# Loop 9 — Observability + planted bug

**Tier mapping:** T6 (entire), T7 peripheral
**Time:** 120–150 minutes
**Theme:** *find the bug from telemetry alone.*
**Prereqs:** Loops 7 + 8 done.

## What you're building

Three layers of telemetry, each with a specific job:

1. **Structured logs** with correlation IDs propagated across services
2. **Metrics** (Prometheus model — counter, gauge, histogram) with the RED method dashboard
3. **Distributed traces** (OpenTelemetry) showing a request's path through HTTP handler → DB query → cache lookup → queue enqueue

Then: the tutor injects a **planted bug** — something subtle (e.g., a 1% slow path, or wrong-tenant data leakage under concurrency). Your job is to find it from telemetry alone, no `println`, no debugger.

## Stack

- **Logs:** stdlib `log/slog` — JSON output, structured fields, correlation ID middleware
- **Metrics:** `github.com/prometheus/client_golang` — counter, histogram, gauge
- **Traces:** `go.opentelemetry.io/otel` — auto-instrumented HTTP handler, DB query spans, manual spans where needed
- **Local stack:** Prometheus + Grafana + Tempo (or Jaeger) via docker-compose

## New files

- `internal/telemetry/log.go` — slog setup, correlation ID middleware
- `internal/telemetry/metrics.go` — Prometheus registry, helpers
- `internal/telemetry/trace.go` — OTel SDK setup, propagation
- `dashboards/red.json` — Grafana dashboard (importable)
- `docker-compose.yml` — adds Prometheus, Grafana, Tempo

## Tasks

### 1. Logs — correlation IDs

```go
// Middleware: pull X-Request-ID header (or generate one), inject into ctx
// and into every log line via slog.With.

func WithRequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rid := r.Header.Get("X-Request-ID")
        if rid == "" {
            rid = newID() // 16-char random
        }
        ctx := context.WithValue(r.Context(), ctxRequestID, rid)
        w.Header().Set("X-Request-ID", rid)
        logger := slog.With("request_id", rid)
        ctx = context.WithValue(ctx, ctxLogger, logger)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

Every log line includes the request_id. Cross-service: propagate the header on outbound calls.

### 2. Metrics — RED

```go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "http_requests_total"},
        []string{"method", "route", "status"},
    )
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
        },
        []string{"method", "route", "status"},
    )
)
```

Note `route` not `path` — `route = "/links/{id}"` not `/links/42`. **High-cardinality labels (user_id, IP, raw path) are the silent killer of Prometheus.** Stretch goal: deliberately mis-add a `user_id` label, watch your scrape time balloon, fix.

### 3. Traces — OpenTelemetry

```go
// Auto-instrument HTTP handlers with otelhttp.NewHandler.
// Manual span around DB queries: trace.SpanFromContext(ctx) ; tracer.Start(ctx, "db.query").
// Add attributes: db.statement, db.system, http.status_code, etc.
// Propagate the trace context on outbound calls (otelhttp.NewTransport).
```

Run Tempo locally; query traces in Grafana. A request's full path becomes a flame chart.

### 4. SLO + alert

Define a single SLO: 99.5% availability over 28 days for `GET /links`.

Two alerts:
- **Fast-burn**: burn rate > 14.4× over 1 hour → page (would burn the whole error budget in 2 days).
- **Slow-burn**: burn rate > 6× over 6 hours → ticket (will burn budget in a week if sustained).

Standard Google SRE multi-window multi-burn-rate alert recipe; copy from the SRE Workbook chapter.

### 5. The planted bug

Tell the tutor "I'm ready for the planted bug." The tutor picks one of:

- **Phantom 1% slow:** randomly, one in 100 requests adds a 2-second delay. p50 looks fine, p99 is bad. Find it from histograms.
- **Wrong-tenant data leakage:** under concurrent load, ~0.1% of `GET /links` returns the wrong user's links. Find it via traces showing the wrong owner_id; correlate with the bug (a missing `WHERE owner_id = $1` in one code path).
- **Slow query that only fires on Sundays:** a query branch hits an unindexed column when a feature flag is on. Find it via DB query histograms tagged by query name.
- **Memory leak in the WS hub:** disconnected clients aren't unregistered cleanly; gauge shows connected-client count grows over time. Find it via a gauge and a runtime metric.

You debug from telemetry alone. The tutor coaches via Socratic questions ("what does the trace say about the slow span?"), not the answer.

## Stretch — cardinality explosion

Add a label `user_id` to `httpRequestsTotal`. Run for 5 minutes with 1000 distinct test users. Look at:

```bash
curl -s http://localhost:9090/api/v1/status/tsdb | jq '.data.headStats.numSeries'
```

Watch the series count balloon. Fix: remove `user_id`. Discuss why low-cardinality labels (status, route) are fine and high-cardinality (user_id, request_id) are not.

Rule of thumb: any label whose cardinality could exceed ~1000 is wrong for Prometheus. (Logs and traces are where high-cardinality data goes; *not* metrics.)
