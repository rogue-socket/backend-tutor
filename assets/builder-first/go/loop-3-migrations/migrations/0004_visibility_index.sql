-- Loop 3 — Step 4: add an index, the safe way.
--
-- CREATE INDEX CONCURRENTLY:
--   - No ACCESS EXCLUSIVE lock — concurrent reads and writes proceed.
--   - Takes 2–3x longer wall-clock than a regular CREATE INDEX.
--   - Cannot run inside a transaction. Migration runners must apply this
--     outside a tx (your runner needs to detect this and skip the BEGIN/COMMIT).
--   - If it fails partway, leaves an INVALID index — drop and retry.
--
-- Without CONCURRENTLY, this would block all writes on `links` for the duration
-- of the build (potentially minutes for a large table).

CREATE INDEX CONCURRENTLY idx_links_visibility ON links(visibility);
