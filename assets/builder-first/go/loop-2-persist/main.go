// Loop 2 — Postgres-backed CRUD with a tags relation.
//
// What's new vs Loop 1:
//   - pgxpool for connection pooling (TODO 1)
//   - migrations/ runner (TODO 2)
//   - Postgres-backed Store methods (TODOs 3–7)
//   - tags table + endpoints (TODOs 8–10) — the N+1 setup
//
// Run:
//   docker compose up -d   # starts Postgres
//   go run .

package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Link struct {
	ID    int      `json:"id"`
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags,omitempty"`
}

type Store struct {
	pool *pgxpool.Pool
}

// TODO 1: NewStore
//
// Read DATABASE_URL from env (default to a local docker-compose value if unset
// — fine for dev). Configure the pool — start with MaxConns=10. Ping the DB
// at startup; if it can't connect, fail fast with a clear error.
func NewStore(ctx context.Context) (*Store, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app:app@localhost:5432/links?sslmode=disable"
	}
	// TODO: parse config, set MaxConns, build pool, ping.
	_ = dsn
	return nil, errors.New("TODO: implement NewStore")
}

// TODO 2: RunMigrations
//
// Read all *.sql files in migrations/, sorted by filename. Apply any that
// aren't already in the schema_migrations table. Each migration in its own
// transaction. ~50 lines.
//
// schema_migrations table:
//   CREATE TABLE IF NOT EXISTS schema_migrations (
//     version TEXT PRIMARY KEY,
//     applied_at TIMESTAMPTZ DEFAULT now()
//   );
func (s *Store) RunMigrations(ctx context.Context, dir string) error {
	return errors.New("TODO: implement RunMigrations")
}

// ─────────────────────────────────────────────────────────────────────────────
// Store methods — same signatures as Loop 1, now backed by Postgres.
// ─────────────────────────────────────────────────────────────────────────────

// TODO 3: List — pagination optional but a LIMIT 50 is non-optional.
//
// Without LIMIT, a hostile or accidental large dataset turns this into a
// denial-of-service. Add it now; you'll thank yourself in Loop 7.
func (s *Store) List(ctx context.Context) ([]Link, error) {
	return nil, errors.New("TODO: implement List")
}

// TODO 4: Get
func (s *Store) Get(ctx context.Context, id int) (Link, bool, error) {
	return Link{}, false, errors.New("TODO: implement Get")
}

// TODO 5: Create — returns the link with the freshly-assigned ID.
//
// Use INSERT ... RETURNING id. Wrap nothing in an explicit transaction unless
// you need atomicity across multiple statements — Postgres autocommits.
func (s *Store) Create(ctx context.Context, url, title string) (Link, error) {
	return Link{}, errors.New("TODO: implement Create")
}

// TODO 6: Update
func (s *Store) Update(ctx context.Context, id int, url, title string) (Link, bool, error) {
	return Link{}, false, errors.New("TODO: implement Update")
}

// TODO 7: Delete
func (s *Store) Delete(ctx context.Context, id int) (bool, error) {
	return false, errors.New("TODO: implement Delete")
}

// ─────────────────────────────────────────────────────────────────────────────
// Tags — the N+1 setup.
// ─────────────────────────────────────────────────────────────────────────────

// TODO 8: TagsForLink — single query for one link's tags.
func (s *Store) TagsForLink(ctx context.Context, linkID int) ([]string, error) {
	return nil, errors.New("TODO: implement TagsForLink")
}

// TODO 9: AddTag
func (s *Store) AddTag(ctx context.Context, linkID int, tag string) error {
	return errors.New("TODO: implement AddTag")
}

// TODO 10: ListWithTagsNaive
//
// Return all links, each populated with its tags. Implement this NAIVELY —
// one query for the link list, then one query per link for tags. This is
// the N+1 your tests should catch.
//
// After you've reproduced N+1 in tests, write ListWithTagsBatched as a second
// method using a single JOIN or an IN (...) batch on tags. The handler should
// switch to the batched version once tests confirm the count drops to 1–2.
func (s *Store) ListWithTagsNaive(ctx context.Context) ([]Link, error) {
	return nil, errors.New("TODO: implement ListWithTagsNaive")
}

// TODO 11 (after BREAK 2): ListWithTagsBatched
func (s *Store) ListWithTagsBatched(ctx context.Context) ([]Link, error) {
	return nil, errors.New("TODO: implement ListWithTagsBatched")
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP — same shape as Loop 1; handlers now take ctx from r.Context().
// ─────────────────────────────────────────────────────────────────────────────

func handleLinks(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// TODO: switch between ListWithTagsNaive (BREAK 2) and ListWithTagsBatched (after fix).
			panic("TODO: GET /links")
		case http.MethodPost:
			// TODO: same as Loop 1 but call s.Create.
			panic("TODO: POST /links")
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleLinkByID(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Match /links/{id} or /links/{id}/tags.
		path := strings.TrimPrefix(r.URL.Path, "/links/")
		parts := strings.Split(path, "/")
		// TODO: parse parts[0] as int ID; if len(parts) == 2 && parts[1] == "tags",
		// dispatch to handleLinkTags.
		_, err := strconv.Atoi(parts[0])
		if err != nil {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		panic("TODO: dispatch /links/{id} and /links/{id}/tags")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// helpers (unchanged from Loop 1)
// ─────────────────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func main() {
	ctx := context.Background()

	store, err := NewStore(ctx)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	if err := store.RunMigrations(ctx, "migrations"); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/links", handleLinks(store))
	mux.HandleFunc("/links/", handleLinkByID(store))

	log.Printf("listening on :8080 (db OK, migrations applied)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// silence unused import warning during scaffolding — remove when you implement.
var _ = pgx.ErrNoRows
