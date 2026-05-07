# Loop 10 — Load test + capacity + postmortem

**Tier mapping:** T7 (entire), T6.7 (postmortems)
**Time:** 180 minutes
**Theme:** *break it on purpose under load, then write the RCA.*
**Prereqs:** Loops 7 + 9 done.

## What you're building

A k6 load test that ramps your service from 0 to 1000 RPS over 5 minutes, soaks at peak for 5, and identifies whatever breaks first. Then a postmortem on whatever you found.

This is the capstone. By Loop 10, you have a service with persistence, auth, async work, caching, deploy, real-time, and observability. Now you'll discover that *something* in there has a capacity ceiling lower than you thought.

## Stack

- **k6** (`brew install k6`) — the load tester. Scripts in JS.
- **Service running** in your `docker compose up` from Loop 7 (with telemetry from Loop 9).
- **Optional:** deploy to Fly.io / Render / Railway and run k6 from your laptop against the public URL — closer to production, slower iterate cycle.

## Tasks

### 1. Define the load profile

```js
// load.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    ramp_then_soak: {
      executor: 'ramping-arrival-rate',
      startRate: 10,
      timeUnit: '1s',
      preAllocatedVUs: 100,
      maxVUs: 500,
      stages: [
        { duration: '5m', target: 1000 },   // ramp 0→1000 RPS
        { duration: '5m', target: 1000 },   // soak at 1000 RPS
        { duration: '1m', target: 0 },      // ramp down
      ],
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],         // <1% errors
    http_req_duration: ['p(99)<500'],       // p99 < 500ms
  },
};

const BASE = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Mix of read-heavy traffic (90% reads, 10% writes) — adjust to your service's
  // traffic shape.
  if (Math.random() < 0.9) {
    const r = http.get(`${BASE}/links/popular`);
    check(r, { 'popular 200': (r) => r.status === 200 });
  } else {
    const r = http.post(`${BASE}/links`, JSON.stringify({
      url: 'https://x', title: 'x',
    }), { headers: { 'Content-Type': 'application/json' } });
    check(r, { 'create 201': (r) => r.status === 201 });
  }
}
```

**Open-loop, not closed-loop.** `ramping-arrival-rate` produces requests at a target *rate*, not a target concurrency. Closed-loop (constant VUs) hides tail latency under back-pressure; open-loop is closer to real production traffic.

### 2. Run baseline

```bash
k6 run load.js
```

Expect (probably):
- p50 latency: rises modestly
- p99 latency: spikes once you hit a bottleneck
- Error rate: spikes once something breaks

Capture the metrics in your Grafana dashboard from Loop 9. Screenshot the spike.

### 3. Identify the bottleneck

The first thing to break is rarely the thing you expected. Common culprits:
- **Connection pool exhaustion** (DB or Redis). `pgxpool.Stat()` shows pool saturation.
- **Single-threaded handler.** A handler that holds a global lock or a single goroutine; CPU is fine but throughput plateaus.
- **DB query slowness on a missing index** that you didn't notice in dev because dev had 100 rows and prod has 1M.
- **GC pauses** under allocation pressure.
- **WebSocket hub backpressure** (Loop 8 territory) if you mixed WS traffic in.
- **Cache stampede** at TTL boundaries (Loop 6 territory) if you didn't single-flight.

The exercise is to *find* what broke, fix it (or cap it), and re-run.

### 4. Write the postmortem

In `~/backend-dev/projects/loop-10-loadtest/postmortem.md`:

```markdown
# links service — load test postmortem — YYYY-MM-DD

## Summary
[1-2 sentences. What broke, what was the user-visible behaviour, how long until recovery.]

## Timeline
- [HH:MM] load test started, ramping to 1000 RPS
- [HH:MM] first error rate above 1%
- [HH:MM] p99 above 1s
- [HH:MM] root contributing factor identified
- [HH:MM] fix applied; recovery confirmed

## Contributing factors
[Plural, system-level. Not "root cause" — the term implies a single owner.]
1. [E.g., "DB connection pool sized at 10, well below the saturation point at 750 RPS read traffic"]
2. [E.g., "popular-links cache TTL of 60s caused a stampede every minute"]
3. [E.g., "no rate limiting on POST endpoint"]

## What went well
- [E.g., "the SLO alert (Loop 9) fired correctly within 90 seconds of the breach"]
- [E.g., "graceful shutdown (Loop 7) meant no requests dropped during the redeploy"]

## What didn't go well
- [E.g., "the runbook didn't cover this specific failure mode"]
- [E.g., "trace context wasn't propagated to the worker, making cross-service debugging harder than it should have been"]

## Action items
- [ ] [Concrete change with owner — even if 'me, by Friday']
- [ ] [Test added to prevent regression — load test in CI?]
- [ ] [Documentation / runbook update]
- [ ] [Observability gap closed (e.g., "add a metric for pool acquisition wait time")]
```

The form matters. Postmortems are an artifact, not a meeting. Future-you (or your team) reads this when something *similar* happens.

### 5. Verify the fix

Re-run the load test with the fix applied. Quantify the improvement:
- Old: p99 = 2s at 800 RPS, errors at 950 RPS
- New: p99 = 200ms at 1000 RPS, no errors

That delta is the value you produced this loop.

## Stretch

- **Reproduce a specific incident from `incidents.md`.** Cache stampede (Loop 6's BREAK at scale). Retry storm (caller has aggressive retries; downstream slows). Connection pool exhaustion (Loop 2's BREAK at scale).
- **Add load testing to CI.** A 30-second smoke test at moderate RPS; fail builds on regression. (Slow; usually nightly.)
- **Capacity plan.** Given the bottleneck you found, what's the max sustainable RPS? At what point does the next bottleneck kick in? Document.
