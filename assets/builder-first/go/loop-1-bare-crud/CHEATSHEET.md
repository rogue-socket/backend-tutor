# Loop 1 cheatsheet

Quick reference for the patterns you'll need. Not a tutorial — a reminder.

## net/http skeleton

```go
mux := http.NewServeMux()
mux.HandleFunc("/links", handleLinks(store))
mux.HandleFunc("/links/", handleLinkByID(store))   // trailing slash = prefix match
http.ListenAndServe(":8080", mux)
```

`HandleFunc("/foo", h)` matches *exactly* `/foo`. `HandleFunc("/foo/", h)` matches `/foo/` and any subpath (`/foo/123`).

## Read JSON body

```go
var req struct {
    URL   string `json:"url"`
    Title string `json:"title"`
}
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    writeError(w, 400, "invalid json")
    return
}
```

`json.NewDecoder` reads streaming; for small bodies it's basically the same as `json.Unmarshal(io.ReadAll(...))`. Both fine here.

## Write JSON response

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(link)
```

**Order matters.** `WriteHeader` flushes the status line; you can't set headers or change the status after that. Always: headers → WriteHeader → body.

## Status codes you'll need

| Code | When |
|---|---|
| 200 OK | Successful GET / PATCH (with body) |
| 201 Created | Successful POST creating a new resource |
| 204 No Content | Successful DELETE (no body) |
| 400 Bad Request | Body is malformed (e.g., bad JSON) |
| 404 Not Found | Unknown URL or unknown resource ID |
| 405 Method Not Allowed | Known URL, wrong method |
| 415 Unsupported Media Type | Wrong / missing Content-Type for a body-bearing request |

`415` is the most-forgotten of the bunch. If you accept JSON, reject anything that isn't `application/json` (or starts with `application/json;` for `application/json; charset=utf-8`).

## Mutex pattern

```go
type Store struct {
    mu     sync.Mutex
    links  map[int]Link
    nextID int
}

func (s *Store) Create(url, title string) Link {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.nextID++
    l := Link{ID: s.nextID, URL: url, Title: title}
    s.links[s.nextID] = l
    return l
}
```

`defer s.mu.Unlock()` is the standard idiom — it releases on every return path including panics. **Always.**

## Race detector

```bash
go test -race ./...
```

Use it in CI. It catches races by instrumenting every memory access — slower at runtime (~5–10x), but priceless. If you only ever remember one Go test flag, make it `-race`.

## Common gotchas

- **Trailing slash on routes.** `mux.HandleFunc("/links", h)` doesn't match `/links/`. The stdlib mux is intentionally strict.
- **Forgetting to set Content-Type before WriteHeader.** Browsers and proxies will guess; don't let them.
- **Reading the request body twice.** `r.Body` is a single-use stream. Read it, parse it, you're done. If you need it twice, `io.ReadAll` it into a buffer first.
- **Unhandled errors from `json.Encoder.Encode`.** Encoding can fail (network closed, etc.). Log it; you can't change the status anymore.
- **Calling `panic` in production handlers.** `net/http` recovers from a handler panic and serves 500, but the request is lost and the panic is logged. Better: write a 500 explicitly with structured information.

## When the tutor asks "what surprised you?"

There's almost always a surprise in Loop 1. Common ones:
- net/http is *much* simpler than people remember after years of frameworks
- Routing in the stdlib is intentionally minimal (this changes in Go 1.22+ with the new `http.ServeMux` patterns)
- `sync.Mutex` defaults to unlocked; you don't need to construct it
- The race detector is built in; no third-party tool needed
- Tests run fast (sub-second for this scope) without any test framework

If your surprise is *"none of this surprised me"*, you may be in the wrong lane — say so, the tutor will re-route up.
