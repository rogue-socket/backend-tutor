# Loop 8 — Real-time + scale out

**Tier mapping:** T0.6 (WebSockets), T1.12 (real-time API patterns), T3.6 (pub/sub), Path: Real-time
**Time:** 180 minutes (split across 2 sittings if needed)
**Theme:** *the live feature is the easy part; making it survive a second instance is the hard part.*
**Prereqs:** Loops 4 + 5 + 7 done.

## What you're building

A **live "who's adding links right now"** feature: when any user creates a link, every other connected user sees a toast / notification in real time. Implemented as WebSockets fanning out from the server.

Two phases:

**Phase A (single instance):** WebSocket server that broadcasts to all connected clients. Two browser tabs see each other.

**Phase B (multi-instance):** Run two app instances behind a load balancer. Tabs connected to instance A don't see events from tabs on instance B. **This is the BREAK.** Fix with Redis pub/sub bridging the instances.

This loop is *the* most visually satisfying break — open four browser tabs, watch one tab's action propagate to the others (or fail to), see the architecture matter immediately.

## Architecture

```
Phase A:
    Client ←→ WebSocket ←→ Server (in-memory hub)
    
Phase B:
    Client A ←→ WS ←→ Server-1 ─┐
    Client B ←→ WS ←→ Server-1 ─┤
                                ├──→ Redis Pub/Sub channel "links.events"
    Client C ←→ WS ←→ Server-2 ─┤
    Client D ←→ WS ←→ Server-2 ─┘
```

Each server subscribes to the channel; on receive, it fans out to its own connected WebSocket clients. Each server publishes to the channel when one of *its* clients triggers an event (e.g., creates a link).

## New files

- `internal/ws/hub.go` — the in-memory Hub: connected clients, broadcast, register, unregister
- `internal/ws/handler.go` — the HTTP upgrade handler
- `internal/ws/bridge.go` — the Redis pub/sub bridge (Phase B)
- `static/index.html` — a tiny client page so you can demo it in a browser

## Tasks

### Phase A — single instance

1. **Implement the Hub.** Goroutine-safe set of `*Client` connections. Register, Unregister, Broadcast methods. Use `chan` for the broadcast queue, not direct method calls — avoids head-of-line blocking on a slow client.
2. **Implement the upgrade handler.** Use `gorilla/websocket` or `nhooyr.io/websocket`. Auth on the upgrade — read the session cookie, validate, attach `userID` to the Client.
3. **Wire `POST /links` to publish to the Hub.** When a link is created, `hub.Broadcast(LinkCreatedEvent{...})`.
4. **Add a static demo page.** Two tabs open at `http://localhost:8080/demo` show each other's link creates.

### Phase B — multi-instance

5. **Run two instances locally.**
   ```bash
   docker compose up --scale server=2
   ```
   Behind nginx or a simple load balancer (or just point one tab at `:8080` and another at `:8081` if you bind two ports).

6. **Verify the BREAK.** Tabs on different instances don't see each other's events.

7. **Implement the Redis bridge.** When a Hub broadcasts an event, also publish to Redis channel `links.events`. Each Hub subscribes to the channel and broadcasts incoming events to its local clients.

8. **Avoid the broadcast loop.** When the bridge receives an event from Redis, broadcast it locally — but don't republish it to Redis. Use a per-event ID + a recent-set to detect "I sent this myself."

9. **Sticky sessions.** Configure the load balancer to route by client (cookie-based or IP-hash). Without sticky sessions, reconnects can land on a different instance, breaking presence state. Demonstrate this once, then enable stickiness.

## Stretch

- **Replay on reconnect.** Each event has a sequence number. Client tracks the last-seen seq. On reconnect, client sends `?since=N`; server sends events from a Redis Stream backlog with id > N. Now reconnects don't lose events.
- **Backpressure.** A slow client can't keep up. Implement: per-client send buffer; if it fills, drop the oldest event (or close the connection). Without this, one slow client can OOM the server.
- **Presence.** Track who's online; broadcast "user X joined" / "left." Heartbeat ping every 30s; unregister on missed pong.

## The fun part

Open four browser tabs side by side. Tabs 1 and 2 connected to instance A; tabs 3 and 4 to instance B. Without the bridge: 1 and 2 see each other's events; 3 and 4 see each other's; nothing crosses. With the bridge: every tab sees every event.

This is a teaching moment that *no amount of explanation can substitute for*. The architecture limitation is visible in the UI.
