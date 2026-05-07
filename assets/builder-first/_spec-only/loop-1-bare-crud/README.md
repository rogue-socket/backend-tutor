# Loop 1 — Bare CRUD (spec-only)

**You're here because** the language you picked (`learner.language` in `progress.json`) doesn't have prefilled scaffolding shipped yet. That's fine. This file is the spec; you implement against it; the tutor reviews.

If you'd prefer prefilled scaffolding, switch language to Go or Python (`/config language go`) and try again.

---

**Tier mapping:** T0 (HTTP), T1 (REST design), T3 (concurrency, peripherally)
**Time:** 90–120 minutes
**Theme:** *feel HTTP without framework magic.*

## What you're building

A single-resource HTTP API for "links" (think: a personal bookmarking service). All five verbs over `/links` and `/links/{id}`. **No routing framework** — use your language's standard library or its closest equivalent (Node: `http`, not Express; Java: `HttpServer`, not Spring; Rust: `hyper` is fine, `axum` is borderline; Kotlin: ktor is fine if you use it minimally). **No database** — an in-memory data structure. **No auth, no logs, no metrics.**

The whole point is to avoid framework magic. If you find yourself reaching for an annotation, decorator, or DI container, you've gone too far. Loop 2 introduces frameworks; Loop 1 is the wire.

## API contract

```
GET    /links            → 200, JSON array of all links
POST   /links            → 201, JSON of the new link with assigned ID
                           (415 if Content-Type isn't application/json,
                            400 if body is malformed,
                            405 on other verbs)

GET    /links/{id}       → 200, JSON of the link, or 404
PATCH  /links/{id}       → 200, JSON of the updated link, or 404
DELETE /links/{id}       → 204 no content, or 404
                           (405 on other verbs)
```

Resource shape:

```json
{
  "id": 42,
  "url": "https://example.com",
  "title": "example"
}
```

ID is server-assigned on POST, immutable thereafter.

## Required tests (≥3)

1. **Happy path:** POST creates a link; GET by the new ID retrieves it; payload matches.
2. **404 on missing:** GET `/links/9999` (or any never-created ID) returns 404.
3. **400 on malformed body:** POST with `{not json` returns 400.

Plus the **concurrency test** (the BREAK — see `BREAK.md`): 100 concurrent POSTs must result in 100 distinct links, no lost updates, no panics.

## Tasks

1. **Set up the project** in your language. `cd ~/backend-dev/projects/loop-1-bare-crud`. Initialize the project (`npm init` / `cargo init` / `mvn archetype:generate` / `gradle init` / equivalent).
2. **Build the in-memory store** with the five operations (List, Get, Create, Update, Delete). **Do not add synchronisation yet.** Loop 1's BREAK depends on you running unsynchronised first.
3. **Build the HTTP layer** — two routes (`/links` and `/links/{id}`), method-dispatch by hand, status codes per the contract.
4. **Write the three required tests.** They should pass.
5. **Write the concurrency test.** It will probably fail — see `BREAK.md`.
6. **Add synchronisation** (mutex / lock / equivalent). Re-run the concurrency test; should pass clean.
7. **Verify against `WIN.md`.**
8. **Fill in `NOTES.md`** — what surprised you, what you'd want to come back to.
9. **Quickpass** — answer `quickpass.json` to consolidate.

## Stretch (if you finish in <60 minutes)

- `GET /healthz` returning `{"status":"ok"}` — you'll need it every loop after.
- Request-log middleware (method, path, status, duration_ms) using the language's stdlib.
- `Location: /links/{id}` header on POST responses (REST convention).

## When you're stuck

- Re-read the contract.
- Look at `assets/builder-first/go/loop-1-bare-crud/main.go` as a reference implementation — it's not your language, but the *shape* (Store struct, two HTTP handlers, helper functions) translates.
- Ask the tutor: *"in [language], how would I [thing]?"* — small specific questions, not "write this for me."
