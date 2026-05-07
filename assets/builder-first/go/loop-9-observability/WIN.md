# WIN — Loop 9 done

## Code

- [ ] Structured logs (slog or equivalent) with `request_id` correlation across handlers
- [ ] Prometheus metrics: RED for HTTP, db query histogram, ws_clients_connected gauge, queue_depth gauge
- [ ] OpenTelemetry traces: HTTP handler auto-instrumented, manual spans on DB and cache calls, propagation on outbound HTTP
- [ ] One SLO defined (99.5% over 28 days for some endpoint)
- [ ] Multi-window multi-burn-rate alerts configured
- [ ] Local Grafana + Prometheus + Tempo running via docker-compose

## Verification

- [ ] You found the planted bug from telemetry alone — `git log` to confirm you didn't read the bug code first
- [ ] Cardinality stretch reproduced: high-cardinality label balloons series count; removal restores Prometheus
- [ ] One incident log line traceable from logs → metric anomaly → trace span

## Understanding

1. **You see p50 = 50ms, p99 = 2s. Averages are 60ms. What does this tell you, and where do you investigate?**
   *Outline: bimodal distribution. Most requests fast, a tail of slow ones. Averages lie because the tail is small. Investigate: histogram buckets — is there a clear secondary peak? Filter traces for duration > 1s — what's the slow span? Look for: a slow downstream that fails closed (timeout retry), a code path under a feature flag, an N+1 that only appears on certain inputs, GC pauses (correlate with `go_gc_duration_seconds`).*

2. **Logs vs metrics vs traces — name a question best answered by each.**
   *Outline: Logs — "what happened to request X?" (high-cardinality detail; user-specific; slow but rich). Metrics — "is the system healthy *right now*?" (low-cardinality summary; cheap to query at scale; fast). Traces — "what was this request doing across services?" (per-request structure; show fan-out, ordering, timing of dependencies). Use the three together — one signal alone is insufficient.*

3. **Sampling: head sampling vs tail sampling on traces. What's the trade-off?**
   *Outline: head sampling (decide at trace start whether to keep it) is cheap and predictable but biased — drops randomly, including the bad ones you most want. Tail sampling (collect everything, decide at end based on duration / errors / etc.) is expensive (you pay for trace data even if you drop it) but keeps the interesting traces. Practical trade-off: head sampling for high-volume baseline + always-on for errored / slow traces.*

## Reflection

What surprised you? Common ones:
- The discipline of "telemetry first, code second" feels slow until it isn't
- Cardinality is *the* metric design constraint — everything else is secondary
- Traces are addictive once you have them

## What's next

Loop 10 — Load test + capacity + postmortem. Break the system on purpose under load; write the RCA.
