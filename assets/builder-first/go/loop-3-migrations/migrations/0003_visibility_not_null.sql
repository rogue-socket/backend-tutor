-- Loop 3 — Step 3: make the column NOT NULL with a default.
--
-- Run this AFTER the backfill script has set every row's visibility.
--
-- Why these two statements together:
--   - SET NOT NULL: brief ACCESS EXCLUSIVE; cheap because the constraint validates
--     against current data (which has no NULLs after step 2).
--   - SET DEFAULT: also brief; ensures inserts from old code paths that don't
--     specify visibility still work.
--
-- DO NOT run this before the backfill — SET NOT NULL on a column with any NULLs fails.

ALTER TABLE links ALTER COLUMN visibility SET NOT NULL;
ALTER TABLE links ALTER COLUMN visibility SET DEFAULT 'public';
