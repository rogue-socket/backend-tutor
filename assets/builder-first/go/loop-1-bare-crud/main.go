// Loop 1 — bare CRUD for a "links" service. In-memory only. No framework.
//
// This file is a starter scaffold. Search for "TODO" — that's where you write code.
// Do not import a router library; the whole point is to feel net/http directly.
//
// Run:
//   go run .
// Test:
//   go test ./...
// Race detector (you'll need this for Loop 1's BREAK):
//   go test -race ./...

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Link is the resource we're CRUDing.
type Link struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// Store holds links in memory.
//
// TODO 1: pick a data structure. A `map[int]Link` is the obvious choice; a
// `[]Link` slice works too but makes Get/Update/Delete O(n). Either is fine for
// Loop 1 — pick one and move on.
//
// TODO 2: add a sync.Mutex (or sync.RWMutex) field. Loop 1's BREAK is what
// happens when you forget this. We'll deliberately leave it out first, run
// `go test -race`, see the race, then add the mutex. Don't add it yet — see
// BREAK.md.
//
// TODO 3: an int counter for generating new IDs.
type Store struct {
	// TODO: fields go here.
}

// List returns all links. Order is undefined for a map; sort by ID if you want
// stable output for tests.
func (s *Store) List() []Link {
	// TODO: lock, copy, return.
	panic("TODO: implement List")
}

// Get returns the link with the given ID. The bool is false if the ID isn't found.
func (s *Store) Get(id int) (Link, bool) {
	// TODO
	panic("TODO: implement Get")
}

// Create inserts a new link. Returns the link with the freshly-assigned ID.
func (s *Store) Create(url, title string) Link {
	// TODO
	panic("TODO: implement Create")
}

// Update patches an existing link. The bool is false if the ID isn't found.
// For Loop 1, treat this as a full replace (not a real PATCH semantics — that
// can wait for Loop 2 or later).
func (s *Store) Update(id int, url, title string) (Link, bool) {
	// TODO
	panic("TODO: implement Update")
}

// Delete removes a link. The bool is false if the ID wasn't there.
func (s *Store) Delete(id int) bool {
	// TODO
	panic("TODO: implement Delete")
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP handlers
//
// We have two routes:
//   /links       — collection: GET (list), POST (create)
//   /links/{id}  — item:       GET, PATCH, DELETE
//
// Everything else gets 404 (unknown route) or 405 (wrong verb on a known route).
// ─────────────────────────────────────────────────────────────────────────────

func handleLinks(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// TODO: write s.List() as JSON, status 200.
		case http.MethodPost:
			// TODO:
			//   1. require Content-Type: application/json (else 415)
			//   2. decode the body into a struct {URL, Title} (else 400)
			//   3. validate URL is non-empty (else 400)
			//   4. s.Create(...)
			//   5. write the new link as JSON, status 201
			_ = strconv.Atoi // remove this line once you remove the panic
			panic("TODO: POST /links")
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleLinkByID(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path is "/links/{id}". Strip the prefix and parse the ID.
		// TODO: extract the ID. If the trailing segment is empty or non-integer,
		// return 404 (per "REST design says unknown URL = 404").
		idStr := strings.TrimPrefix(r.URL.Path, "/links/")
		_ = idStr // delete once you parse it

		switch r.Method {
		case http.MethodGet:
			// TODO: s.Get(id), 200 or 404
			panic("TODO: GET /links/{id}")
		case http.MethodPatch:
			// TODO: decode body, s.Update, 200 or 404
			panic("TODO: PATCH /links/{id}")
		case http.MethodDelete:
			// TODO: s.Delete, 204 or 404
			panic("TODO: DELETE /links/{id}")
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers — write these once, use everywhere.
// ─────────────────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// At this point we've already written the status; logging is the most
		// we can do.
		log.Printf("encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	// Simple error envelope for Loop 1. We'll get fancier in Loop 2+.
	writeJSON(w, status, map[string]string{"error": msg})
}

// ─────────────────────────────────────────────────────────────────────────────
// main
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	store := newStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/links", handleLinks(store))
	mux.HandleFunc("/links/", handleLinkByID(store))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// newStore returns an initialised Store. Move the construction logic here so
// main and tests share it.
func newStore() *Store {
	// TODO: initialise your fields.
	return &Store{}
}
