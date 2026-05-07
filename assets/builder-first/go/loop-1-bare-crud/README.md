# Loop 1 — Bare CRUD

**Tier mapping:** T0 (HTTP), T1 (REST design), T3 (concurrency, peripherally)
**Time:** 90–120 minutes
**Theme:** *feel HTTP without framework magic.*

## What you're building

A single-resource HTTP API for "links" (think: a personal bookmarking service). All five verbs over `/links` and `/links/{id}`. No framework — `net/http` only. No database — an in-memory map. No auth, no logs, no metrics. Just the request/response cycle, by hand.

## Why this is Loop 1

Every backend service is, underneath all the abstraction, a program that listens on a port, parses HTTP, does some work, and writes a response. Frameworks (Gin, Echo, FastAPI, Express, Spring) hide that loop behind decorators, middlewares, and DI containers — which is fine for shipping, but disastrous for *understanding*. Loop 1 makes you feel HTTP at the wire.

## Run it

```bash
cd loop-1-bare-crud
go run .
```

Should print `listening on :8080`. Then in another terminal:

```bash
curl -X POST http://localhost:8080/links \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","title":"example"}'
# → 201 Created with the new link

curl http://localhost:8080/links
# → 200 OK with the list
```

## Test it

```bash
go test ./...           # 4 tests, three pre-written + the concurrency test
go test -race ./...     # the race-detector run — Loop 1's BREAK relies on this
```

## Tasks (in order)

1. **Read `main.go` top to bottom.** Don't write code yet. Get the structure in your head.
2. **Implement `Store` methods** (`List`, `Get`, `Create`, `Update`, `Delete`) — without a mutex. We'll add one after the BREAK.
3. **Implement `handleLinks`** — GET list and POST create.
4. **Implement `handleLinkByID`** — GET, PATCH, DELETE on `/links/{id}`.
5. **Run `go test`.** First three tests should pass.
6. **Run `go test -race`.** The concurrency test will fail or panic. **This is the BREAK.** Read `BREAK.md` for what to observe.
7. **Add a `sync.Mutex` to `Store`.** Run `go test -race` again — it should pass clean.
8. **Verify against `WIN.md`.** All criteria checked.
9. **Fill in `NOTES.md`** with what surprised you. Keep it short.
10. **Quickpass** — answer the three questions in `quickpass.json` to consolidate. The tutor will run these for you on `/loop quickpass`.

## Stretch (if you finish in <60 minutes)

- Add a `GET /healthz` endpoint that returns `{"status":"ok"}` with status 200. (You'll need it for every later loop anyway.)
- Add a request log middleware that prints `method path status duration_ms` for each request. (No external libraries; wrap your handler.)
- Make the JSON output deterministic — sort `List()` by `ID` so tests don't flake.
- Add a `Location: /links/{id}` header on POST responses (REST convention).

## Hints (only read when stuck)

<details>
<summary>Hint: parsing the ID from /links/{id}</summary>

Use `strings.TrimPrefix(r.URL.Path, "/links/")` to get the trailing segment, then `strconv.Atoi`. If the result is empty or doesn't parse, return 404. (Why 404 and not 400? Because the URL is the resource locator; an unparseable URL = unknown resource.)
</details>

<details>
<summary>Hint: 405 vs 404</summary>

`/links` accepts GET and POST. PATCH on `/links` (not `/links/{id}`) → 405. GET on `/widgets` → 404 (route doesn't exist at all).
</details>

<details>
<summary>Hint: the mutex placement</summary>

Don't lock inside `Get`/`List` and forget `Create`/`Update`/`Delete`. Lock anywhere the map is touched, both reads and writes. `sync.RWMutex` is fine but `sync.Mutex` is simpler and probably faster at this scale. Pick simple.
</details>
