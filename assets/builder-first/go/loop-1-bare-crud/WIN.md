# WIN — Loop 1 done

You've finished Loop 1 when **all** of these are true. Don't move to Loop 2 until they are.

## Code

- [ ] `Store` has a mutex; `go test -race ./...` passes clean
- [ ] All 5 verbs work: GET list, GET item, POST create, PATCH update, DELETE
- [ ] Status codes are right:
  - 200 on successful GET / PATCH
  - 201 on successful POST (with the new resource as body — bonus: `Location` header)
  - 204 on successful DELETE
  - 400 on malformed body
  - 404 on unknown ID or unknown route
  - 405 on wrong method against a known route
  - 415 on missing/wrong `Content-Type` for POST/PATCH
- [ ] Error responses are JSON, not HTML or plain text
- [ ] Whole thing is under ~200 lines of Go

## Tests

- [ ] At least 3 tests pass (the three scaffolded ones count, but adding a fourth — your choice — solidifies the muscle memory)
- [ ] `go test -race ./...` passes
- [ ] Tests are self-contained — no dependence on a running server

## Understanding (the part the tutor will probe)

You can answer these without looking anything up:

1. **Why does the race detector flag the unsynchronised map access — what specifically does it see?**
   *Outline: it instruments memory access and detects two accesses to the same address from different goroutines without an intervening synchronisation event (lock, channel send/receive, sync.WaitGroup, etc.). The Go memory model defines "happens-before" via these events; absence of a happens-before edge between two accesses to the same address = a race.*

2. **Why is `sync.Mutex` enough — why don't we need atomic operations or a channel-based design?**
   *Outline: a mutex provides the happens-before edge required for memory visibility, and is the simplest tool that works. Atomics are for single-word operations (counters); channels are for ownership transfer or coordination; mutexes are for protecting a region of code that touches multiple fields. The Store touches multiple fields per operation (the map and the counter), so a mutex is the right tool.*

3. **Why 404 on `GET /links/abc` (non-integer ID), not 400?**
   *Outline: REST design says the URL identifies a resource. An unparseable resource identifier means "no such resource" — that's 404. 400 is for malformed requests where the URL is fine but the body is bad. Different layer.*

## Reflection (1 paragraph in `NOTES.md`)

What surprised you in this loop? Pick one thing — concurrency, status code semantics, the shape of `net/http`, anything — and write 2–3 sentences. The tutor will read this and pick the next loop's emphasis based on it.

## What's next

Loop 2 — Persistence + the concurrency break. The in-memory map gets replaced by Postgres, and you'll discover that "the DB handles concurrency for me" is mostly true and partly a lie.
