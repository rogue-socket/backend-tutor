# WIN — Loop 1 done (spec-only)

You've finished Loop 1 when **all** of these are true. Don't move to Loop 2 until they are.

## Code

- [ ] Store has synchronisation; the concurrency test passes consistently across multiple runs (`for i in {1..10}; do <run-tests>; done`)
- [ ] All 5 verbs work: GET list, GET item, POST create, PATCH update, DELETE
- [ ] Status codes match the contract in `README.md`
  - 200 on successful GET / PATCH
  - 201 on successful POST
  - 204 on successful DELETE
  - 400 on malformed body
  - 404 on unknown ID or unknown route
  - 405 on wrong method against a known route
  - 415 on missing/wrong `Content-Type` for POST/PATCH
- [ ] Error responses are JSON, not HTML or plain text
- [ ] Whole thing is small — well under your language's "reasonable single-file size"

## Tests

- [ ] At least 3 functional tests pass
- [ ] The concurrency test passes consistently
- [ ] Tests are self-contained — no dependence on a running server

## Understanding (the part the tutor will probe)

You can answer these without looking anything up:

1. **What concurrency primitive did you use, and why that one?**
   *Outline: name the primitive (mutex, lock, semaphore, channel, atomics) and the reason — usually "simplest tool that gives me mutual exclusion across reads and writes." Atomics are for single-word ops; channels are for ownership transfer; mutexes are for protecting multi-field state.*

2. **Why does the read path also need synchronisation, even though no two reads conflict?**
   *Outline: visibility, not mutual exclusion. Without a synchronisation event, one thread's writes may not be visible to another thread (memory model / cache coherence). The lock acquire / release establishes the happens-before relationship. Some runtimes (Go) document this in their memory model; others (Java with the JMM, C++ with std::memory_order) do too.*

3. **Why 404 on `GET /links/abc` (non-integer ID), not 400?**
   *Outline: REST design says the URL identifies a resource. An unparseable resource identifier means "no such resource" — that's 404. 400 is for malformed requests where the URL is fine but the body or query is bad. Different layer.*

## Reflection (1 paragraph in `NOTES.md`)

What surprised you in this loop? Pick one thing and write 2–3 sentences. The tutor reads this and picks the next loop's emphasis based on it.

## What's next

Loop 2 — Persistence + the concurrency break. The in-memory store gets replaced by Postgres, and you'll discover that "the DB handles concurrency for me" is mostly true and partly a lie.
