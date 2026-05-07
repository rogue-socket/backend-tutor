# WIN — Loop 2 done

## Code

- [ ] All Loop 1 endpoints work, now backed by Postgres
- [ ] Data survives a `docker compose down && docker compose up` cycle (volume-persisted)
- [ ] `migrations/0001_init.sql` runs on startup; re-running is a no-op
- [ ] `tags` resource added: `GET /links/{id}/tags`, `POST /links/{id}/tags`
- [ ] `GET /links` returns links with their tags inline, in 1–2 queries (BREAK 2 fixed)
- [ ] Connection pool sized explicitly (~10 for dev), with a documented reason

## Tests

- [ ] All Loop 1 tests still pass against Postgres
- [ ] `TestListWithTagsIsNotNPlusOne` asserts query count ≤ 2

## Understanding

You can answer these:

1. **What does `pgx`'s `pgxpool.Pool` give you that a single `*pgx.Conn` doesn't?**
   *Outline: a bounded set of reusable connections, with acquire/release semantics. Solves the "every request opens a new TCP+TLS connection" problem and bounds DB load.*

2. **You set `MaxConns=2` in BREAK 1 and saw request timeouts. Why doesn't the DB show errors?**
   *Outline: the application can't get a connection from the pool, so the request never reaches the DB. The DB sees only the 2 long-running queries; from its perspective there's no problem.*

3. **Your `ListWithTagsBatched` uses `WHERE link_id = ANY($1)` with a Go slice. What happens when the slice has 10K entries?**
   *Outline: the parameter is sent as a single Postgres array; the planner generates a hash join or nested loop depending on stats. There's no parameter-count limit (unlike a literal `IN (a, b, c, ...)` list which hits Postgres's 32K argument limit at extremes). Practical limit is bandwidth / planner cost, not protocol.*

## Reflection

In `NOTES.md`, write 2–3 sentences on what surprised you. Common surprises:
- Postgres-via-Docker is *fast* (sub-ms queries from your laptop)
- N+1 is invisible without query logging
- Pool sizing is a real decision with real consequences

## What's next

Loop 3 — Migrations. You'll add a NOT NULL column to a populated table without downtime, learning expand/contract.
