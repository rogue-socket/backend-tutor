# WIN — Loop 3 done

## Code

- [ ] Migrations 0002, 0003, 0004 all applied; the `links` table has a NOT NULL `visibility` column with default `'public'` and an index on it
- [ ] `scripts/backfill/backfill.go` runs idempotently (rerunning is a no-op)
- [ ] Application code writes `visibility` on Create / Update; reads handle NULL (in case backfill is incomplete)
- [ ] Migration runner handles `CREATE INDEX CONCURRENTLY` (skips wrapping it in a transaction)

## Verification (you tested this)

- [ ] Parallel inserter ran during the safe sequence — zero failed inserts
- [ ] Confirmed (on a scratch DB) that the wrong migration *would* have blocked under load
- [ ] You set `lock_timeout` to 1s before any DDL and felt the migration abort cleanly when blocked

## Understanding

You can answer these:

1. **Why does `ALTER TABLE ADD COLUMN visibility TEXT NOT NULL DEFAULT 'public'` work in milliseconds on Postgres 11+ but historically required a full rewrite?**
   *Outline: pre-PG11, every existing row had to be physically updated to add the new column with its default. PG11+ stores the default as metadata; old rows return the default at read-time; new rows write it directly. The NOT NULL + DEFAULT combination becomes metadata-only. (NOT NULL without a DEFAULT still requires that no existing rows are NULL — but PG also tracks this as a per-column "has-no-NULLs" flag now.)*

2. **You ran `CREATE INDEX CONCURRENTLY` and it failed partway. What state is the index in, and how do you recover?**
   *Outline: an INVALID index. Postgres marks it `indisvalid = false` in `pg_index`. Queries don't use it. You can't just CREATE INDEX with the same name. Recovery: `DROP INDEX CONCURRENTLY <name>; CREATE INDEX CONCURRENTLY <name> ...;`. (Also: investigate why it failed — usually a unique-violation in the data or a deadlock with an unrelated session.)*

3. **Why does `lock_timeout = 1s` matter for production migrations even when the migration itself is fast?**
   *Outline: locks queue. A 50ms migration behind a 5-minute idle transaction blocks every subsequent query for those 5 minutes — total outage. `lock_timeout` makes the migration abort instead of queueing. You retry instead of taking down the service. Standard practice in any production migration framework.*

## Reflection

What surprised you in this loop? Common surprises:
- The wrong migration is *fast*. The outage is from queueing, not the work itself.
- Backfill scripts are unceremonious — small loops, no fancy framework needed.
- "Online" migrations are entirely about *coexistence* — the old code path has to keep working.

## What's next

Loop 4 — Auth. Sessions or JWT, password hashing, cookie flags, and three breaks: algorithm confusion, CSRF, token logging.
