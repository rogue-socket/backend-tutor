// Integration tests for Loop 2.
//
// These tests assume Postgres is running locally on :5432 (via `docker compose up -d`).
// They use a separate test database to avoid polluting dev data.
//
// Run:
//   docker compose up -d
//   go test ./...
//
// First three tests cover the basic CRUD flow over Postgres. The fourth is
// Loop 2's BREAK 2: it asserts query count to catch N+1.

package main

import (
	"context"
	"os"
	"testing"
)

// testStore returns a fresh Store backed by the test database. The schema is
// migrated; tables are truncated between tests via the t.Cleanup hook.
func testStore(t *testing.T) (*Store, context.Context) {
	t.Helper()
	if os.Getenv("DATABASE_URL_TEST") == "" {
		os.Setenv("DATABASE_URL_TEST", "postgres://app:app@localhost:5432/links_test?sslmode=disable")
	}
	// TODO: build a pool against DATABASE_URL_TEST, run migrations, return.
	// Add a t.Cleanup that TRUNCATEs links, tags, link_tags between tests.
	t.Skip("TODO: implement testStore — see TODO 1 in main.go for the pool wiring")
	return nil, context.Background()
}

func TestCreateThenGet(t *testing.T) {
	s, ctx := testStore(t)
	got, err := s.Create(ctx, "https://example.com", "example")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if got.ID == 0 {
		t.Fatalf("create returned ID=0")
	}
	got2, found, err := s.Get(ctx, got.ID)
	if err != nil || !found {
		t.Fatalf("get: %v, found=%v", err, found)
	}
	if got2.URL != "https://example.com" {
		t.Fatalf("get returned URL=%q", got2.URL)
	}
}

func TestUpdateMissingIsNotFound(t *testing.T) {
	s, ctx := testStore(t)
	_, found, err := s.Update(ctx, 9999, "x", "y")
	if err != nil {
		t.Fatalf("update missing: %v", err)
	}
	if found {
		t.Fatalf("update missing: want found=false")
	}
}

func TestDeleteIdempotent(t *testing.T) {
	s, ctx := testStore(t)
	created, _ := s.Create(ctx, "x", "y")
	ok, _ := s.Delete(ctx, created.ID)
	if !ok {
		t.Fatalf("first delete should report ok=true")
	}
	ok2, _ := s.Delete(ctx, created.ID)
	if ok2 {
		t.Fatalf("second delete should report ok=false (already gone)")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// BREAK 2: the N+1 detector.
//
// This test seeds 10 links each with 3 tags, then calls ListWithTagsNaive
// while counting queries. With the naive impl, it'll see 1 + 10 = 11 queries.
// After fixing with ListWithTagsBatched, it should see 1–2.
//
// Implementation note: pgx supports a Tracer hook on the pool config. Wire one
// up in testStore() that increments a counter on each query, expose the
// counter to this test.
// ─────────────────────────────────────────────────────────────────────────────

func TestListWithTagsIsNotNPlusOne(t *testing.T) {
	s, ctx := testStore(t)

	// Seed: 10 links, 3 tags each.
	for i := 0; i < 10; i++ {
		l, err := s.Create(ctx, "https://x", "x")
		if err != nil {
			t.Fatal(err)
		}
		for _, tag := range []string{"a", "b", "c"} {
			if err := s.AddTag(ctx, l.ID, tag); err != nil {
				t.Fatal(err)
			}
		}
	}

	// TODO: reset the query counter, call s.ListWithTagsBatched, assert count <= 2.
	// While the BREAK is intentional, this test is what tells you the FIX worked.
	t.Skip("TODO: wire up the query counter via pgx Tracer, then assert <=2 queries")
}
