# BREAK — Loop 7

Three breaks. All three are common production failures.

---

## BREAK 1 — Readiness probe takes the cluster down

**Foundations filled here:** T5.5 (health checks).

### The setup

You set `/readyz` to ping Postgres on every check. Kubernetes (or your LB) probes every 5 seconds. Postgres has a 5-second hiccup (transient blip — it happens).

Every replica fails its readiness probe within a single check interval. Every replica gets pulled from the load balancer simultaneously. The service is fully unavailable for ~10 seconds — much longer than the underlying DB hiccup.

### Why this happens

Coupling readiness to a *shared dependency* turns one dependency's hiccup into a service outage. The DB hiccup affects every replica equally; correlated failure → correlated unhealthy → no replicas left.

### The fix

Several non-exclusive options:

1. **Don't fail readiness on remote dependencies.** Cache the last-good ping result; readiness reflects the cached state with a stale-OK window.
2. **Slack the threshold.** Require N consecutive failures before pulling, with N tuned to be longer than typical hiccups.
3. **Distinguish "I can serve some traffic" from "I can serve all traffic."** Some endpoints don't need the DB; readiness shouldn't fail for them.

The conceptual point: **readiness should reflect the replica's local capability**, not the global health of every dependency. Monitors and alerts handle dependency health; readiness handles "should this replica receive new traffic."

### Test

Inject a 5-second DB hiccup (`docker compose pause db; sleep 5; docker compose unpause db`). Watch readiness behaviour. With naive impl: all replicas pulled, total outage. With cached/lagged impl: traffic continues, possibly slower.

---

## BREAK 2 — `kill -9` mid-request

**Foundations filled here:** T5.8 (graceful shutdown).

### The setup

You ship a service that handles SIGTERM correctly... or you think you do. The deploy system (or `docker stop`) sends SIGTERM, waits 30 seconds, then SIGKILL. If you ignore SIGTERM (or don't drain), every in-flight request at the SIGKILL moment is lost.

To reproduce: start the service. Hammer it with `hey -z 60s -c 50 http://localhost:8080/links`. In another terminal: `docker kill -s SIGKILL <container>` (note: `kill` defaults to SIGTERM; we're forcing SIGKILL to simulate a deploy bug).

Without graceful shutdown: every in-flight request at the kill moment dropped → ~50 errors.

### What "graceful" actually means

```go
srv.Shutdown(ctx)     // 1. stop accepting new connections
                      // 2. let in-flight requests complete
                      // 3. close the listener
                      // 4. return when all connections idle, or ctx times out
```

This is built into `net/http`. The bug is usually one of:
- Not wiring it up to a signal handler
- Not giving it a long-enough context (5 seconds is often too short for slow downstream calls)
- Closing the DB pool *before* shutdown returns (in-flight requests can't query)

### Test

After the fix, `docker stop` (sends SIGTERM) and watch the load tester. Should see zero errors during a clean shutdown.

---

## BREAK 3 — Secret in image / hardcoded config

**Foundations filled here:** T8.4 (secrets), T9.1 (12-Factor III).

### The setup

You wrote:
```dockerfile
ENV DATABASE_URL=postgres://app:hunter2@db:5432/links
```

Or worse, baked the credential into a `config.json` that's `COPY`'d into the image.

### What happens

`docker history --no-trunc <image>` shows every line of the Dockerfile, including the ENV line. Anyone with read access to the image (which often means: anyone, eventually) has the credential.

Or: scan the image with Trivy / Grype / dive. Find the credential.

### The fix

12-Factor III: config in the environment, *injected at runtime*, not at build time. Per-environment configs as runtime flags or env vars; secrets via a secret manager (AWS Secrets Manager, GCP Secret Manager, HashiCorp Vault, or just `fly secrets set` / `kubectl create secret`).

```bash
# Right:
docker run -e DATABASE_URL=... links:v1

# Or:
fly secrets set DATABASE_URL=...
fly deploy
```

### Detection in CI

```bash
gitleaks detect --source . --report-path /tmp/gitleaks.json
trufflehog filesystem .
```

Add to CI; fail builds on findings.

---

## The takeaway

> The 12-Factor App is operational law, not philosophy.

Three rules:
1. **Readiness is local capability, not global health.**
2. **Graceful shutdown is the difference between a clean deploy and a 0.5% error spike every time.**
3. **Secrets at runtime, never at build time.**
