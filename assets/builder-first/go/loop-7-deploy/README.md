# Loop 7 — Containerize + deploy

**Tier mapping:** T9 (entire — 12-Factor, containers, CI/CD), T5.5 (health checks), T5.8 (graceful shutdown)
**Time:** 120–180 minutes
**Theme:** *runs on my laptop is not a deploy strategy.*
**Prereqs:** Loops 2 + 4 + 5 + 6 done.

## What you're building

A real Docker image of your service, plus a complete `docker-compose.yml` running the full stack (server + worker + Postgres + Redis), plus a deploy to a real cloud — easy mode: Fly.io / Render / Railway. Hard mode: AWS ECS or EC2 with Terraform.

Plus three operational essentials added to your service: health endpoints, graceful shutdown, and 12-Factor-clean configuration.

## New files

- `Dockerfile` — multi-stage build, distroless or scratch base, non-root user
- `docker-compose.yml` — full stack (already touched in Loops 2, 5, 6 — finalize here)
- `.dockerignore` — keep credentials, build artifacts, and bloat out of the build context
- (Add to existing) `GET /healthz` and `GET /readyz` handlers
- (Add to existing) graceful shutdown in `main.go`

## Tasks

### 1. Dockerfile (multi-stage)

```dockerfile
# Stage 1: build
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/worker ./cmd/worker

# Stage 2: runtime — distroless static, no shell, no apt, ~2 MB
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
COPY --from=build /out/worker /worker
COPY migrations /migrations
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
```

Build, scan, run:

```bash
docker build -t links:dev .
trivy image links:dev          # any critical CVEs? fail if so
docker run --rm -p 8080:8080 -e DATABASE_URL=... links:dev
```

Image should be ~10–20 MB.

### 2. Health endpoints

```go
// /healthz — am I alive (the process is running)? Cheap, no external deps.
func handleHealthz(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, 200, map[string]string{"status": "ok"})
}

// /readyz — am I ready to serve traffic? Checks critical dependencies.
// IMPORTANT: don't fail readyz on every dependency hiccup. See BREAK 1.
func handleReadyz(db *pgxpool.Pool, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
        defer cancel()

        if err := db.Ping(ctx); err != nil {
            writeJSON(w, 503, map[string]string{"status": "db_unavailable"})
            return
        }
        if err := rdb.Ping(ctx).Err(); err != nil {
            // Note: cache failure shouldn't necessarily fail readyz —
            // the service can degrade gracefully. Decide per dependency.
            writeJSON(w, 503, map[string]string{"status": "redis_unavailable"})
            return
        }
        writeJSON(w, 200, map[string]string{"status": "ready"})
    }
}
```

Difference matters:
- **liveness** failure → orchestrator restarts the container
- **readiness** failure → orchestrator pulls the container from the load balancer (no restart)

Confusing them is BREAK 1.

### 3. Graceful shutdown

```go
srv := &http.Server{Addr: ":8080", Handler: mux}

go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("listen: %v", err)
    }
}()

sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
<-sigCh
log.Println("shutdown signal received")

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
    log.Printf("shutdown error: %v", err)
}
db.Close()
rdb.Close()
log.Println("clean exit")
```

Key behaviours:
- `srv.Shutdown(ctx)` stops accepting new connections, lets in-flight requests complete
- 30-second budget aligns with Kubernetes default `terminationGracePeriodSeconds`
- Close DB / Redis pools *after* HTTP shutdown — they may be in use by in-flight requests

### 4. 12-Factor cleanup

Audit your service against the 12 factors. Most likely violations:
- **III. Config** — anything hardcoded (DB URL, Redis addr, log level)?
- **VI. Processes** — any local-disk state? (Sessions in DB ✓; uploads, if any, must go to object storage)
- **IX. Disposability** — graceful shutdown ✓; fast startup ≤ a few seconds?
- **XI. Logs** — writing to stdout, not log files? (You should be.)

### 5. Deploy

Easy mode (10 minutes): Fly.io.

```bash
brew install flyctl
fly launch                     # detects Dockerfile, asks a few questions
fly postgres create            # provisions a managed DB
fly redis create               # provisions managed Redis
fly secrets set DATABASE_URL=... REDIS_URL=...
fly deploy
```

Hard mode (90+ minutes): AWS ECS Fargate with Terraform. Worth doing once for the experience; not necessary for Loop 7's WIN.

## Stretch

- **Image scan in CI.** GitHub Actions: `trivy image --severity CRITICAL --exit-code 1 links:${{ github.sha }}`. Fail builds on critical CVEs.
- **Reproducible builds.** Pin Go version exactly; pin base image by digest (`distroless/static-debian12@sha256:...`).
- **Multi-arch image.** `docker buildx build --platform linux/amd64,linux/arm64 ...` — useful for ARM cloud (cheaper) and Apple Silicon dev.
