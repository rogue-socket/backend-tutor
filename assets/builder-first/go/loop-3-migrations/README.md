# Loop 3 — Migrations + schema evolution

**Tier mapping:** T2.5 (online migrations)
**Time:** 90–120 minutes
**Theme:** *the schema you wrote in Loop 2 was wrong; fix it without downtime.*
**Prereqs:** Loop 2 done.

## What you're building

Add a `visibility` column to the populated `links` table — values `'public'` or `'private'`, NOT NULL, default `'public'`. Do it in a way that lets you keep deploying the running service through the migration.

You'll do this **twice**:
1. **The wrong way** (`migrations/9001_visibility_wrong.sql` — name it 9001 so it sorts late and isn't accidentally applied). One statement, takes ACCESS EXCLUSIVE, blocks all reads and writes for the duration.
2. **The right way** (expand/contract): four small steps, each commits independently, no statement holds a long lock.

Both produce the same end-state schema. Only one is shippable.

## The expand/contract sequence

```
Step 1 (migration 0002): ALTER TABLE links ADD COLUMN visibility TEXT;
                         (NULL allowed; new code can write 'public' or 'private')

Step 2 (script):         UPDATE links SET visibility = 'public' WHERE visibility IS NULL
                         (in batches of 1000, sleep 100ms between batches)

Step 3 (migration 0003): ALTER TABLE links ALTER COLUMN visibility SET NOT NULL;
                         ALTER TABLE links ALTER COLUMN visibility SET DEFAULT 'public';
                         (the NOT NULL is the only locking part — fast, since column has no NULLs)

Step 4 (migration 0004): CREATE INDEX CONCURRENTLY idx_links_visibility ON links(visibility);
                         (CONCURRENTLY = no exclusive lock; takes longer wall-clock but doesn't block writes)
```

The application code rolls forward in parallel:
- Before step 1: code doesn't know about `visibility`. Inserts succeed (no NOT NULL yet).
- After step 1, before step 2: new code paths can write `visibility`; old paths still work (NULL is allowed).
- After step 2: all rows have a value.
- After step 3: NOT NULL is enforced; old code paths that don't write `visibility` will *now* fail unless step 3's `SET DEFAULT` is in place. **That's why step 3 sets default and not-null in the same migration.**
- After step 4: the index exists; queries filtering by visibility get the speedup.

## Tasks

1. **Read all four files in `migrations/` plus `9001_visibility_wrong.sql`.** Understand each.
2. **Apply step 1.** `go run .` (your migration runner from Loop 2 picks up the new file).
3. **Update your application code** to write `visibility` on Create / Update. Default to `'public'` if the request doesn't specify. Old data still has NULL — handle that on read by treating NULL as `'public'`.
4. **Run the backfill script.** `go run ./scripts/backfill`. It updates rows in batches of 1000 with 100ms between batches. While it runs, hammer the API with inserts — none should fail.
5. **Apply step 3** (NOT NULL + DEFAULT). Should complete in milliseconds since no NULLs exist.
6. **Apply step 4** (index CONCURRENTLY). Will take longer (10s–minutes for large tables) but won't block writes. While it runs, keep hammering the API with inserts/updates.
7. **Try the WRONG migration on a copy of the DB**: `psql -c "$(cat migrations/9001_visibility_wrong.sql)"` while a parallel insert script runs. Watch the inserts block.

## Stretch

- Read [Strong Migrations](https://github.com/ankane/strong_migrations) — Rails-flavored but the rules apply universally. Map each rule to a Postgres lock level.
- Use `pg_locks` to actually observe what's locked during the wrong migration: `SELECT * FROM pg_locks WHERE relation = 'links'::regclass;`
- Try `lock_timeout` — set it to 1s before the wrong migration. The migration aborts instead of holding the queue.
