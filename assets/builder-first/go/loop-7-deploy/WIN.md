# WIN — Loop 7 done

## Code

- [ ] Multi-stage Dockerfile; final image < 20 MB; non-root user
- [ ] `docker compose up --build` runs the full stack with one command
- [ ] `/healthz` and `/readyz` endpoints implemented; readyz cached with sensible staleness
- [ ] SIGTERM handled; in-flight requests drain; clean exit within 30s
- [ ] All config via env vars; no secrets in code, image, or git history
- [ ] Trivy scan: zero critical CVEs

## Verification

- [ ] BREAK 1 reproduced: naive readyz pulls cluster on DB hiccup; cached readyz doesn't
- [ ] BREAK 2 reproduced: `docker stop` (SIGTERM) drains cleanly; `docker kill -SIGKILL` doesn't
- [ ] BREAK 3 reproduced: `docker history` reveals secrets in old impl; nothing in fixed impl

## Deploy

- [ ] Service is reachable on the public internet at a stable URL
- [ ] DB and Redis are managed services (not your laptop)
- [ ] Secrets configured via the platform's secret store
- [ ] Health endpoints visible to the platform's LB

## Understanding

1. **You see a 1% error rate spike during every deploy. The deploy is rolling, 5 minutes total. Where do you look first?**
   *Outline: graceful shutdown. The 1% is the in-flight requests at the moment each replica is terminated. If `srv.Shutdown` isn't wired, or the timeout is too short, or the LB removes the replica from rotation slower than the replica stops accepting requests — every replica replacement loses the in-flight tail. Confirm: log SIGTERM receipt; measure shutdown duration; check LB drain timing relative to terminationGracePeriodSeconds.*

2. **Distroless images have no shell. How do you write a Docker healthcheck?**
   *Outline: either (a) include a small healthcheck binary (`/healthcheck` that does an HTTP GET to `/healthz`), (b) implement a `-healthcheck` flag on your main binary that does the check internally and exits 0/1, or (c) use a non-distroless base image that has wget/curl. Most production teams use (a) or (b); (c) defeats the security benefit.*

3. **Why is `terminationGracePeriodSeconds` (Kubernetes) typically 30s, and what should your shutdown timeout be?**
   *Outline: 30s is a balance between "long enough for typical request lifetimes" and "short enough that a stuck pod doesn't block deploys." Your in-app shutdown context should be slightly less than this (e.g., 25s) so you finish before SIGKILL. If your typical request is longer than 30s (long polls, large file uploads), bump both — or split the long-running paths to a separate service with its own SLO.*

## Reflection

What surprised you? Common ones:
- Distroless's lack of shell breaks half your habits but eliminates whole categories of attack
- Graceful shutdown is one of those things that's correct or wrong, never partially right
- "Just use 12-Factor" is only obvious in retrospect

## What's next

Loop 8 — Real-time + scale out. WebSockets, the second-instance break, Redis pub/sub.
