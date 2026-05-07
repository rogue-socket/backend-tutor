# BREAK — the concurrency race (language-agnostic)

**Foundations filled here:** T3.1 (concurrency models), T3.2 (data races), T3.3 (mutexes).

## The setup

You implemented the in-memory store *without* any synchronisation. The first three tests pass. Then you ran the concurrency test — 100 concurrent inserts. Something went wrong:
- A panic / runtime error / exception about concurrent access
- A wrong final count (you got 87 links instead of 100)
- A test that "passed" once, then failed the next run

This non-determinism is the lesson.

## What's happening

Two (or more) threads / tasks / goroutines / coroutines are reading and writing the same data structure at the same time without coordination. Different runtimes manifest this differently:

| Runtime | Manifestation |
|---|---|
| Go | `fatal error: concurrent map writes` panic; `go test -race` shows the data race |
| Java / Kotlin | `ConcurrentModificationException` on collections; subtler corruption on plain primitives |
| Node.js | Single-threaded JS event loop is *mostly* safe, but worker_threads + shared array buffers replicate the problem; `await` between read and write opens a window |
| Python | GIL prevents most low-level corruption, but read-modify-write (`counter += 1`) is *not* atomic — there's a bytecode boundary between read and write where another thread can run |
| Rust | The compiler refuses to compile shared mutable access without a `Mutex`, `RwLock`, or atomic; if you got here, you used `unsafe` or `Mutex::lock().unwrap()` without thinking about it |

The detail differs; the underlying cause is identical: **shared mutable state without synchronisation = undefined behaviour.**

## What to observe

Run the test multiple times. Different runs may produce different failures. **Concurrency bugs are not bugs you can rely on reproducing on demand.** This is exactly what makes them dangerous in production: by the time you notice, the timing has already passed.

## The fix

Add a synchronisation primitive: a mutex / lock / equivalent. Lock around every method that touches the shared state — both reads and writes. (Read-only methods need locks too: many runtimes don't promise visibility of one thread's writes to another thread without a synchronisation event.)

In your language:
- **Java/Kotlin:** `synchronized` block, or `ReentrantLock`, or use `ConcurrentHashMap` directly
- **Node:** if using worker_threads, use `Atomics` and `SharedArrayBuffer`; otherwise the event loop already serialises, but be careful with `await` mid-mutation
- **Rust:** wrap the data in `Mutex<T>` or `RwLock<T>`; the compiler will then guide you
- **Python:** `threading.Lock` around mutations; consider whether you actually need threads at all (asyncio is usually the better choice)
- **Other:** check the language's stdlib for `Mutex`, `Lock`, `Synchronize`, or `RWLock` patterns

Re-run the concurrency test. It should pass.

## Why this is Loop 1's break, not Loop 2's

You'll move to a real database in Loop 2, and the database will give you ACID transactions. So why care about in-memory races now?

Because **the in-memory map is not the only shared state in your service.** A connection pool is shared state. A cache is shared state. A counter for metrics is shared state. A "current leader" pointer in a leader-elected service is shared state. Every backend service has dozens of shared-state structures that are *not* the database.

Loop 1's break teaches you to *see* shared state and ask "what synchronises this?" — a habit you'll need every loop after this one.

## The takeaway

> Concurrent access to shared mutable state without synchronisation is **always** a bug. Either the runtime catches it, the type system catches it, or it ships and you find out the hard way.
