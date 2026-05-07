// Tests for Loop 1.
//
// Three tests are pre-written: happy-path POST, 404 on a missing GET, and
// 400 on malformed JSON. After these pass, write at least one more test
// of your choosing (the WIN.md criteria require ≥3 — these three count, but
// adding a fourth solidifies the muscle memory).
//
// Run:
//   go test ./...
// Race detector (Loop 1's BREAK relies on this):
//   go test -race ./...

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestPostThenGet(t *testing.T) {
	store := newStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/links", handleLinks(store))
	mux.HandleFunc("/links/", handleLinkByID(store))

	// POST /links
	body := bytes.NewBufferString(`{"url":"https://example.com","title":"example"}`)
	req := httptest.NewRequest(http.MethodPost, "/links", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /links: want 201, got %d (body=%s)", rec.Code, rec.Body.String())
	}
	var created Link
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("POST response not valid JSON: %v", err)
	}
	if created.ID == 0 || created.URL != "https://example.com" {
		t.Fatalf("POST response: unexpected payload %+v", created)
	}

	// GET /links/{id}
	req2 := httptest.NewRequest(http.MethodGet, "/links/"+itoa(created.ID), nil)
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("GET /links/%d: want 200, got %d", created.ID, rec2.Code)
	}
}

func TestGetMissingIs404(t *testing.T) {
	store := newStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/links/", handleLinkByID(store))

	req := httptest.NewRequest(http.MethodGet, "/links/9999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET missing: want 404, got %d", rec.Code)
	}
}

func TestPostBadJSONIs400(t *testing.T) {
	store := newStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/links", handleLinks(store))

	req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewBufferString(`{not json`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST bad JSON: want 400, got %d", rec.Code)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Loop 1 BREAK — the concurrency race.
//
// Run with `go test -race -run TestStoreConcurrency`. Without a mutex on the
// Store's fields, this test will trigger Go's race detector or panic with
// "fatal error: concurrent map writes". After you add the mutex, this test
// should pass cleanly with -race.
// ─────────────────────────────────────────────────────────────────────────────

func TestStoreConcurrency(t *testing.T) {
	store := newStore()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			store.Create("https://example.com/"+itoa(i), "title "+itoa(i))
		}(i)
	}
	wg.Wait()

	got := len(store.List())
	if got != 100 {
		t.Fatalf("after 100 concurrent Creates, want 100 links, got %d", got)
	}
}

// itoa is a tiny helper so we don't pull in strconv just for tests.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
