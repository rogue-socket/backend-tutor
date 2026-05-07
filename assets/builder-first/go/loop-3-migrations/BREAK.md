# BREAK — Loop 3

**Foundations filled here:** T2.5 (online migrations), T2.3 (transactions and isolation, peripherally).

## The setup

You have a populated `links` table from Loop 2. You're going to apply two migrations:
1. `9001_visibility_wrong.sql` — the obvious one-statement migration
2. `0002` → `0004` — the safe expand/contract sequence

To feel the difference, you'll run a parallel insert script in another terminal.

## The parallel inserter

In a second terminal, before running either migration:

```bash
while true; do
  curl -s -X POST -H 'Content-Type: application/json' \
    -d '{"url":"https://x","title":"x"}' \
    http://localhost:8080/links > /dev/null
  echo -n .
  sleep 0.05
done
```

This fires ~20 inserts/second. Watch the dots stream.

## BREAK 1 — apply the wrong migration

Apply `9001_visibility_wrong.sql` to a *copy* of the DB (not your dev DB — you don't want to recover from this):

```bash
docker compose exec db createdb -U app links_scratch
docker compose exec db pg_dump -U app links | docker compose exec -T db psql -U app links_scratch
```

Then in your inserter, point it at `links_scratch` (or just hammer SELECT against it). And in a third terminal:

```bash
docker compose exec db psql -U app links_scratch -c '\timing on'
docker compose exec db psql -U app links_scratch -c 'ALTER TABLE links ADD COLUMN visibility TEXT NOT NULL DEFAULT '"'"'public'"'"';'
```

### What you'll see

- On a small table (< 100K rows on a fast disk): the migration completes in tens of milliseconds; you barely notice.
- On a larger table (1M+ rows): you'll see the inserts pause. Run the migration timing — it could be hundreds of ms to multiple seconds. During that window, every query against `links` is queued behind the ACCESS EXCLUSIVE lock.

To make the lock *visible* even on a small table, hold a transaction open in another psql session:

```sql
BEGIN;
SELECT * FROM links LIMIT 1;
-- (don't COMMIT; this holds an ACCESS SHARE lock)
```

Then run the wrong migration. **It will block, indefinitely, waiting for the SELECT's lock to release.** Worse, every subsequent query (including the inserts) queues behind the migration. One sleeping psql + one in-flight migration = total outage.

### Why this is the lesson

Lock queueing is the silent killer. The migration itself isn't slow; the *queue forming behind it* is what causes the outage. A 50ms migration with a 5-minute idle transaction in front of it = 5 minutes of blocked queries.

Standard mitigation: `SET lock_timeout = '1s'; BEGIN; ALTER TABLE ...; COMMIT;` — if the lock can't be acquired in 1 second, the migration aborts cleanly instead of queueing.

## BREAK 2 — apply the safe sequence

Restart your dev DB. Apply `0002`, then run the backfill script, then `0003`, then `0004`. Run the parallel inserter throughout.

Expected:
- **Step 0002 (ALTER TABLE ADD COLUMN):** invisible; inserts don't pause.
- **Backfill script:** runs for as long as your dataset takes; inserts continue uninterrupted.
- **Step 0003 (SET NOT NULL + DEFAULT):** brief lock; might cause a few queued inserts but completes in milliseconds.
- **Step 0004 (CREATE INDEX CONCURRENTLY):** takes longer (seconds to minutes for large tables) but inserts continue. Without CONCURRENTLY, this would block writes for the entire build.

You should see **zero failed inserts** during the safe sequence.

## The takeaway

> The hardest part of a schema migration is "the old code is still running." Expand/contract is the discipline that lets you ship the schema and the code rollouts independently.

Standard rules (memorise these):

1. **Don't combine column-add with NOT NULL + DEFAULT** unless the table is small *and* you've staged the deploy carefully.
2. **NOT NULL** requires a backfill or a default — if backfill, do it in batches.
3. **Indexes on populated tables**: always `CREATE INDEX CONCURRENTLY`. Always.
4. **Foreign keys on populated tables**: `ALTER TABLE ... ADD CONSTRAINT ... NOT VALID;` then `ALTER TABLE ... VALIDATE CONSTRAINT ...;`. Two steps.
5. **DROP COLUMN**: cheap, but commits the column shape forever (Postgres keeps the dropped column metadata; pg_dump/restore is the only way to actually reclaim space).
6. **Always set `lock_timeout`** before any DDL on a populated table.
