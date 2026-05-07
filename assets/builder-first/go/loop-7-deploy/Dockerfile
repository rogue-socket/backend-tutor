# Loop 7 — multi-stage Dockerfile.
#
# Goals:
#   - Final image < 20 MB (distroless static)
#   - Non-root user
#   - Reproducible (pin tags; pin sha256 in production)
#   - No shell, no package manager → smaller attack surface

# ─── Stage 1: build ──────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS build

WORKDIR /src

# Cache deps in a separate layer so dep changes don't invalidate the build cache.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO disabled = fully static binary, copyable into distroless/static.
# -trimpath: removes filesystem paths from the binary (reproducibility).
# -ldflags '-s -w': strip debug info, smaller binary.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/worker ./cmd/worker

# ─── Stage 2: runtime ────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/server /server
COPY --from=build /out/worker /worker
COPY migrations /migrations

USER nonroot

EXPOSE 8080

# Default: run the server. Override with `docker run ... /worker` for the worker.
ENTRYPOINT ["/server"]
