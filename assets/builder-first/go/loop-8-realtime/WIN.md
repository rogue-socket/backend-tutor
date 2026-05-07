# WIN — Loop 8 done

## Code (Phase A)

- [ ] WebSocket upgrade handler with auth (session cookie validated on upgrade)
- [ ] In-memory Hub: register, unregister, broadcast
- [ ] `POST /links` triggers a broadcast
- [ ] Demo page (`static/index.html`) shows two tabs seeing each other's events

## Code (Phase B)

- [ ] Redis pub/sub bridge: publish on local broadcast, subscribe and forward on receive
- [ ] Self-events deduped (instance ID or recent-IDs set)
- [ ] Two instances locally; four browser tabs (2 per instance) all see each other's events
- [ ] Sticky sessions configured at the LB

## Verification

- [ ] BREAK 1 reproduced: two instances without bridge silo events; with bridge, all clients see all events
- [ ] BREAK 2 demonstrated: sticky sessions enabled; reconnects land on the same instance
- [ ] BREAK 3 demonstrated: slow client drops events instead of blocking the broadcast loop

## Understanding

1. **Pub/sub vs streams (Redis): when to pick which for cross-instance broadcast?**
   *Outline: pub/sub for fire-and-forget broadcast where missing a message is fine (live activity, presence, real-time stats). Streams for replay-capable, persistent fan-out (chat history, audit events). Pub/sub is lighter (no storage), lower latency. Streams give you an offset to replay from on reconnect — necessary if a momentary disconnect must not lose messages. Many real systems use both: pub/sub for hot-path low-latency, streams for backfill on reconnect.*

2. **Why does the bridge need dedup?**
   *Outline: the publishing instance is also subscribed to the channel. Without dedup, every local broadcast → publish → comes back as a remote event → broadcast locally again → publish again → infinite loop. Dedup options: tag each event with publishing instance ID and skip-self, or maintain a recent-IDs LRU and skip seen events. Both work; instance-ID is simpler when bridge topology is known.*

3. **Sticky sessions vs shared state: which makes more sense for a presence feature on WebSocket?**
   *Outline: depends on scale and feature shape. Sticky sessions (LB routes a client consistently to one backend) make presence trivially correct on one instance — but a hot user always lands there forever. Shared state (presence in Redis with heartbeat-driven TTL) lets any instance read/write — scales horizontally but adds latency and a dependency. For < 100K concurrent users, sticky is fine. For larger scales, shared state is required. Most production WebSocket systems use both: sticky for easy cases, Redis-backed state for cross-instance coordination.*

## Reflection

What surprised you? Common ones:
- The single-instance → multi-instance break is *visible* in the UI; impossible to ignore
- Dedup logic in the bridge took longer than the bridge itself
- Backpressure handling is the slow-burn correctness hit; everything looks fine until it isn't

## What's next

Loop 9 — Observability + planted bug. Wire up OpenTelemetry, debug a planted issue from telemetry alone.
