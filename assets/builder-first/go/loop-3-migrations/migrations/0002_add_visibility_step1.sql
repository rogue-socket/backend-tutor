-- Loop 3 — Step 1: add the column nullable.
--
-- Lock taken: ACCESS EXCLUSIVE on `links`, but only briefly (metadata-only change).
-- On Postgres 11+ this is fast even on huge tables because there's no rewrite.
-- Concurrent reads and writes are blocked for milliseconds; survivable for any service
-- that handles transient DB hiccups.
--
-- Why not include NOT NULL DEFAULT in the same statement?
--   - On Postgres 11+, ADD COLUMN with a CONSTANT default is also metadata-only.
--   - But: NOT NULL on a column with no default = full table rewrite OR pre-existing rows must already conform.
--   - We'll add NOT NULL in step 3, AFTER backfill ensures there are no NULLs.

ALTER TABLE links ADD COLUMN visibility TEXT;
