---
last_verified: 2026-05-08
verified_with: go1.26.3 darwin/arm64
staleness_threshold_days: 180
---

# LOOP_VERSIONS.md

What versions the shipped builder-first scaffolding was last verified against. If `last_verified` above is older than `staleness_threshold_days` (default 180), pins may have moved on — re-verify before assuming the scaffolding is current.

The tutor checks this file's date at workspace setup and warm-resume. You can also run the staleness check manually:

```bash
python3 ~/.claude/skills/backend-tutor/tools/check-staleness.py
```

(For non-Claude-Code harnesses, point that command at wherever you cloned the skill.)

---

## Toolchain floors

| Language | Floor | Reason |
|---|---|---|
| Go | 1.24 | loop-2-persist's pgx v5.7.5 transitive deps require it |
| Python | 3.12 | planned (FastAPI scaffolding not yet shipped) |
| Node | 20 LTS | spec-only |
| Java | 21 (Temurin) | spec-only — virtual threads |
| Rust | latest stable | spec-only |

Bumping a floor: update both this table *and* `assets/builder-first/setup/README.md`. The tests cross-check that loops actually build with the floor stated here.

---

## Per-loop pinned dependencies

### Go

**loop-1-bare-crud**

- stdlib only (no external deps)
- `go.mod` carries `go 1.22` directive — anything 1.22+ works

**loop-2-persist**

- `github.com/jackc/pgx/v5` — `v5.7.5` (last verified 2026-05-08; latest at verification was v5.9.2 — held back to v5.7 to keep a wider Go-version compat window)
- `go.mod` carries `go 1.24.0` directive (forced by pgx transitive deps)
- Postgres: scaffolding tested against Postgres 16 in Docker; 14+ should work

**loops 3-10**

- Build as deltas onto loop-2's module — same dep tree applies. Loops introduce additional packages (`github.com/redis/go-redis/v9`, `github.com/IBM/sarama` for Kafka, etc.) but those aren't pinned in go.mod files because the loops are merge-deltas, not standalone modules. Cross-check the README of each loop for specific suggested versions.

### Python (planned, not shipped)

When FastAPI scaffolding lands, this section will pin: `fastapi`, `uvicorn`, `sqlalchemy`, `alembic`, `pytest`, `httpx`. Lock format will be `uv.lock` (uv-managed).

### Node / Java / Kotlin / Rust (spec-only)

No prefilled scaffolding ships for these languages — the learner implements against the loop spec. No version pins to maintain here; the loop README points at canonical libraries (Express/Fastify, Spring Boot, Ktor, Axum) without pinning specific versions.

---

## How to refresh

When `last_verified` goes stale (>180 days):

1. From the repo root, run `go mod tidy` in each `assets/builder-first/go/loop-*/` directory that has a `go.mod`. Note any version bumps and whether the `go` directive moved.
2. If the `go` directive moved up, update the `Toolchain floors` table above and the setup README.
3. Run `python3 tests/run_all.py` and confirm all tests still pass — particularly `test_loop_versions.py`, which cross-checks this manifest against actual `go.mod` contents.
4. Update the `last_verified` date in this file's frontmatter.
5. Commit with a `LOOP_VERSIONS.md: refresh pins YYYY-MM-DD` message.
