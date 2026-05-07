# Loop 2 — Persistence + the concurrency break

**Tier mapping:** T2 (Postgres, indexes, transactions, connection pools), T3 (concurrency revisit)
**Time:** 120–150 minutes
**Theme:** *the in-memory map was a lie.*
**Prereqs:** Loop 1 done.

## What you're building

Replace Loop 1's in-memory `Store` with Postgres-backed persistence. Same HTTP API surface (the handlers don't change). New concerns: connection pooling, schema migrations, query plans, transactions, and the related-data N+1 problem.

You'll also add a `tags` resource related to links — that's what sets up the N+1 BREAK.

## Run it

```bash
docker compose up -d         # starts Postgres on :5432
go run .                     # runs migrations on startup, then serves on :8080
```

The first time you run `go run .`, it'll execute `migrations/0001_init.sql`. Inspect with:

```bash
docker compose exec db psql -U app links -c '\d links'
```

## Tasks

1. **Read `main.go`.** New imports: `pgxpool`, `pgx/v5`. Same handler shape as Loop 1.
2. **Wire up the connection pool.** TODO 1.
3. **Implement the migration runner.** TODO 2 — read SQL files in `migrations/` in order, apply each in a transaction, track applied versions in a `schema_migrations` table.
4. **Implement the Store methods** against Postgres. TODOs 3–7. Use `pgxpool.Pool.Query` / `QueryRow` / `Exec`.
5. **Add tags.** New table, new endpoints `GET /links/{id}/tags`, `POST /links/{id}/tags`. The list endpoint `GET /links` should return links *with their tags inline*. **First: implement this naively** — for each link in the list, run a separate query for its tags. This is the N+1 setup for BREAK 2.
6. **Run the integration tests.** They assume `docker compose up -d` has been run.
7. **Run BREAK 1** (pool exhaustion). See `BREAK.md`.
8. **Run BREAK 2** (N+1). See `BREAK.md`. Fix with a single JOIN or `IN (...)` batch.
9. **Verify against `WIN.md`.**

## Stretch

- Add an index on `links.created_at`. Run `EXPLAIN ANALYZE SELECT ... ORDER BY created_at DESC LIMIT 50` before and after. Save the plans.
- Use `database/sql` instead of pgx for one method, then revert. Compare ergonomics.
- Set `pgxpool.Config.MaxConns = 2` and rerun the load test from BREAK 1. Lower the pool until you see the difference clearly.

## Hints

<details>
<summary>Hint: where the migration runner goes</summary>

Run it once at startup, before `ListenAndServe`. Read `migrations/*.sql` sorted by filename, query `schema_migrations` for already-applied versions, apply the missing ones in a transaction each. ~50 lines.
</details>

<details>
<summary>Hint: spotting N+1</summary>

Set `pgxpool.Config.ConnConfig.Tracer = ...` to log every query. Then hit `GET /links`. If you see N+1 queries (one for the list, one per link for tags), you've reproduced the bug.
</details>
