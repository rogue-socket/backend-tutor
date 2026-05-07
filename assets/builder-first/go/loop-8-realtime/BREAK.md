# BREAK — Loop 8

This loop has one main BREAK plus two stretch breaks.

---

## BREAK 1 — The second instance silos

**Foundations filled here:** T3.6 (pub/sub), Path: Real-time.

### The setup

You finished Phase A. The single-instance demo works: tab 1 creates a link, tab 2 sees it.

Now scale to two instances:

```bash
docker compose up --build --scale server=2
```

(Or run two `go run ./cmd/server` processes on different ports.)

Open four browser tabs:
- Tabs 1, 2 → instance A
- Tabs 3, 4 → instance B

Tab 1 creates a link. **Tab 2 sees it. Tabs 3 and 4 don't.**

### Why

The Hub is in-memory, per-process. Instance A's Hub knows nothing about instance B's connected clients. The broadcast goes to instance A's clients only.

### What you'll observe in the demo UI

- Tab 1 (instance A): submits a link → instance A's `POST /links` handler creates the link in DB, broadcasts to instance A's local Hub.
- Tab 2 (instance A): WebSocket receives the event.
- Tab 3 (instance B): `GET /links` works (DB is shared), but the WebSocket never sees the event.
- Tab 4 (instance B): same as tab 3.

This is the *single most common* gotcha in real-time backends. Every developer assumes broadcast works "everywhere"; without a bridge, it works only within one process.

### The fix

Redis pub/sub. Each instance:
1. **Subscribes** to a shared channel (`links.events`)
2. **Publishes** every local broadcast to that channel
3. On receiving from the channel, fans out to its *local* clients (without republishing — that would loop forever)

After wiring this up, all four tabs see the event.

### Why pub/sub specifically

For broadcast / fan-out, pub/sub is the right tool:
- Lightweight (no persistence; messages are fire-and-forget)
- Low-latency (sub-ms typical)
- No ordering guarantee, but for "live activity feed" use cases that's fine

For features needing persistence/replay (e.g., chat history), use Redis Streams (Loop 5's tool) plus pub/sub. For very high fan-out (millions of subscribers), use a dedicated tool (NATS, Kafka) — but Redis pub/sub holds up well to the tens of thousands of subscribers most apps actually need.

### The dedup problem

Your bridge publishes events on local broadcast and consumes from the channel on receive. Without dedup, your *own* events come back to you and get re-broadcast — to the same local clients — duplicating. And re-published, looping forever.

Two patterns:
1. **Instance ID + skip-self.** Tag each event with the publishing instance's ID; subscribers skip events where `instance_id == my_id`.
2. **Recent-IDs set.** Tag each event with a unique ID; subscribers maintain a small recent-IDs LRU; skip events whose IDs they've seen recently.

Pattern 2 is more general (handles cross-region replay, retries) but slightly heavier. Pattern 1 is simpler if you trust your bridge to never resend.

---

## BREAK 2 — Sticky sessions and the reconnect

**Foundations filled here:** Path: Real-time, T0.6.

### The setup

WebSocket connections are long-lived. With round-robin load balancing, a *reconnect* can land on a different instance — fine for stateless requests, but your real-time state (e.g., "user X is online", "user X is in chat room Y") may be local to the original instance.

To reproduce: implement a simple presence feature ("show online users"). Restart instance A. The clients that were connected to A reconnect; round-robin LB sends some to B. B doesn't know they were ever online; the presence state is wrong.

### The fix

**Sticky sessions** — load balancer routes a client consistently to the same backend, by cookie hash or IP hash.

Or — better, harder — **make presence stateless across instances.** Store presence in Redis with a TTL (heartbeat extends TTL); any instance can read and update.

The sticky-session approach is simpler but limits horizontal scale (a hot user lands on one instance forever). The Redis-presence approach scales but adds complexity.

For Loop 8's demo, sticky sessions are fine. Real production WebSocket systems often use both: sticky sessions for the easy 80% + Redis-backed shared state for the rest.

---

## BREAK 3 — Backpressure: the slow client

**Foundations filled here:** T0.6 (WebSockets backpressure).

### The setup

A client connects, then doesn't read fast enough (mobile on bad network, browser tab in background, malicious slow read). The server's send queue for that client fills.

If the server's broadcast loop *blocks* on a slow client's queue, every other client suffers. One slow client = global outage.

### The fix

Per-client bounded send buffer with non-blocking enqueue:

```go
select {
case c.send <- payload:
    // sent; carry on
default:
    // queue full; drop or close
    log.Printf("client %d slow — dropping event", c.UserID)
}
```

Two reasonable policies:
- **Drop the event.** Client misses an event but stays connected. Right for ephemeral feeds (live activity).
- **Close the connection.** Forces a reconnect, often with a replay protocol. Right for stateful feeds where missing an event matters.

Loop 5 (queue) and Loop 9 (observability) interact here: a closed-connection metric should alarm if it spikes.

---

## The takeaway

> Single-instance WebSocket servers are toys; the moment you horizontal-scale, you need a pub/sub layer. Sticky sessions and backpressure are the supporting cast.

Three rules:
1. **Every broadcast crosses instances via a pub/sub layer.**
2. **Sticky sessions for stateful real-time, or shared state via Redis.**
3. **Slow clients drop or disconnect — never block the broadcast loop.**
