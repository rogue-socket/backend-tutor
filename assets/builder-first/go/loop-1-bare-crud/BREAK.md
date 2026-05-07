# BREAK — the concurrency race

**Foundations filled here:** T3.1 (concurrency models), T3.2 (data races), T3.3 (mutexes).

## The setup

You implemented the `Store` methods without any locking. The first three tests pass. Then you ran:

```bash
go test -race ./...
```

…and `TestStoreConcurrency` panics, hangs, or reports `WARNING: DATA RACE`.

## What's happening

Two (or more) goroutines are reading and writing the same map at the same time. Go's runtime detects this and either:

1. **Panics with `fatal error: concurrent map writes`** — the runtime's safety net. Visible without `-race`.
2. **Reports `WARNING: DATA RACE`** under `-race` — the race detector saw two unsynchronised accesses to the same memory.
3. **Silently corrupts** if you got lucky with timing — bytes written by one goroutine get half-overwritten by another. This is the worst case; it's why you run `-race` in CI.

This isn't a Go quirk. It's the same in every shared-memory concurrent runtime — Java, C++, Rust, Python (yes, even with the GIL — read-modify-write is still split across multiple bytecodes). The map is shared mutable state; concurrent access without synchronisation is undefined behaviour.

## What to observe

Run the test a few times. Different runs may produce different errors:
- A panic on map writes
- A race detector report pointing at the exact two lines (one read, one write)
- A test that "passes" (got lucky) but still fails the next run

This non-determinism is the lesson. **Concurrency bugs are not bugs you can rely on reproducing on demand.**

## The fix

Add a `sync.Mutex` (or `sync.RWMutex`) to the `Store` struct. Lock in every method that touches the map. Be paranoid — read paths need locks too, because Go's memory model doesn't promise visibility of writes across goroutines without synchronisation.

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

After the fix, `go test -race ./...` should pass clean.

## Why this is Loop 1's break, not Loop 2's

You'll move to a real database in Loop 2, and the database will give you ACID transactions for free. So why care about in-memory races now?

Because **the in-memory map is not the only shared state in your service.** A connection pool is shared state. A cache is shared state. A counter for metrics is shared state. A "current leader" pointer in a leader-elected service is shared state. Every backend service has dozens of shared-state structures that are not the database.

Loop 1's break teaches you to *see* shared state and ask "what synchronises this?" — a habit you'll need every loop after this one.

## The takeaway

> Concurrent access to shared mutable state without synchronisation is **always** a bug. Either the runtime catches it, the race detector catches it, or it ships and you find out the hard way.

Add this to your `NOTES.md` if it landed.

## Stretch — if you want to go deeper now

- Run with `GOMAXPROCS=1`. Does the race still happen? (It shouldn't — explain why.)
- Replace `sync.Mutex` with `sync.RWMutex`. Measure throughput on a read-heavy workload (lots of `List`, occasional `Create`). Is it faster?
- Read [The Go Memory Model](https://go.dev/ref/mem). Don't try to memorize it — read it once, know it's there.
