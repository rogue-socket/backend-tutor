// Backfill script for Loop 3, step 2.
//
// Updates `visibility` to 'public' for any links where it's NULL, in batches of
// 1000, with a 100ms pause between batches. The pause is the throttle —
// without it, this script can saturate the DB on large tables and starve the
// application of connections.
//
// Run:
//   go run ./scripts/backfill
//
// Idempotent: rerunning is safe (only NULLs get updated). Stop and resume safe:
// uses a WHERE filter, no per-row tracking needed.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	batchSize = 1000
	pause     = 100 * time.Millisecond
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app:app@localhost:5432/links?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	totalUpdated := 0
	for {
		// TODO: write the UPDATE.
		//
		// Two reasonable shapes:
		//
		// Shape A — RETURNING with a CTE for batching:
		//   WITH to_update AS (
		//     SELECT id FROM links WHERE visibility IS NULL LIMIT $1 FOR UPDATE SKIP LOCKED
		//   )
		//   UPDATE links SET visibility = 'public'
		//   WHERE id IN (SELECT id FROM to_update)
		//   RETURNING id
		//
		// Shape B — UPDATE with a subquery limit (simpler, slightly worse under
		// concurrent backfills):
		//   UPDATE links SET visibility = 'public'
		//   WHERE id IN (SELECT id FROM links WHERE visibility IS NULL LIMIT $1)
		//
		// Both work. Shape A is safer if you ever run two backfills in parallel
		// (e.g., during a botched re-run); FOR UPDATE SKIP LOCKED prevents
		// double-counting.
		//
		// Implement, then count rowsAffected.
		var rowsAffected int
		_ = pool // remove once you've written the query

		if rowsAffected == 0 {
			break
		}
		totalUpdated += rowsAffected
		fmt.Printf("updated %d (total %d)\n", rowsAffected, totalUpdated)
		time.Sleep(pause)
	}

	fmt.Printf("done — %d rows updated\n", totalUpdated)
}
