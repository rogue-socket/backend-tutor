# BREAK — Loop 5

Three breaks. Run them in order.

---

## BREAK 1 — Ack before work

**Foundations filled here:** T3.7 (delivery semantics).

### The setup

You wrote the worker like this:
```go
for _, m := range msgs {
    q.Ack(ctx, m.ID)        // ← ack first
    process(m)              // ← then do the work
}
```

Or perhaps more subtly:
```go
go process(m)               // ← spawn
q.Ack(ctx, m.ID)            // ← ack immediately
```

### What happens

Kill the worker mid-batch (`kill -9 $(pgrep worker)`). The messages were acked, so Redis considers them done. They're never redelivered. The work didn't happen.

### Reproduce

Add a 1-second sleep in `process`. Send 10 messages. Kill the worker after 2 seconds. Restart the worker. Check: how many notifications were actually sent?

Should be ~2; will be 0 if you ack-before-work.

### The fix

Ack **after** successful processing. Always. Make it impossible to forget — wrap the call in a helper that takes a closure, runs it, and acks on success.

### Why this is subtle

The "ack first" pattern is tempting because it makes the next message available faster. The latency improvement is microseconds; the correctness cost is total. **Latency-via-ack-first is never worth it.**

For libraries that auto-ack-on-receive (some MQTT libs, some Kafka client misconfigurations): turn that off. You're choosing the wrong default for any non-trivial workload.

---

## BREAK 2 — Non-idempotent consumer

**Foundations filled here:** T3.5 (idempotency in handlers), T3.7.

### The setup

You skipped the dedup check. Consumer code:
```go
process(m)                  // sends an email
q.Ack(ctx, m.ID)
```

Looks correct. It's not.

### What happens

The message is processed once. Then the worker hiccups *between* `process` and `Ack` — network blip, GC pause, Redis briefly unreachable. The ack doesn't land. After `min-idle-time`, Redis re-delivers via XPENDING / XAUTOCLAIM. The new worker instance processes the same message again — sends a second email.

User got two emails. They open a support ticket.

### Reproduce

Add a chaos hook to your worker: 10% probability of `os.Exit(1)` between `process` and `Ack`. Restart on crash. Run with 100 messages. Count emails sent vs messages enqueued. Difference = your duplicate count.

### The fix

Idempotent dedup. Two patterns, both fine:
- **Redis set** — `SADD notifications:seen <key>` with TTL; `SISMEMBER` before processing. Lightweight, eventually consistent (the SADD itself can fail mid-flight, but you get at most one duplicate per failure window).
- **Postgres unique constraint** — insert into a `processed_messages` table with `(message_key) UNIQUE`; `ON CONFLICT DO NOTHING`. Stronger consistency, transactional with the work itself if you do both in the same transaction.

For `link.created` notifications, the dedup key is the `link_id` — that's unique and stable.

### Why this is the lesson

"Exactly-once" delivery is a marketing term. The engineering reality is at-least-once + idempotent consumer. Every queue, every webhook, every retry — the receiving side is responsible for dedup. The broker can't do it for you, no matter what their docs claim.

---

## BREAK 3 — Poison message blocks the queue

**Foundations filled here:** T3.9 (DLQ).

### The setup

You enqueue a message that will *always* fail to process — a malformed payload, a reference to a deleted link, a divide-by-zero on a specific user. The worker tries it, fails, doesn't ack. Redelivery fires. Worker tries again. Fails. Forever.

Without max-deliveries + DLQ, this either:
- Blocks the queue (if you have a single worker and stop processing other messages while the bad one churns)
- Wastes resources (parallel workers happily process other messages but burn CPU on the bad one indefinitely)
- Fills logs with the same error

### Reproduce

Enqueue a "send notification for link 999999999" — a non-existent link. Worker fails on the lookup. Watch redelivery happen via `XPENDING notifications workers`. Watch deliverycount climb.

### The fix

In the worker loop, before calling `process`:
```go
if queue.MaxDeliveriesReached(m.DeliveryCount) {
    log.Printf("DLQ %s: %d deliveries", m.ID, m.DeliveryCount)
    _ = q.SendToDLQ(ctx, m, "max deliveries exceeded")
    continue
}
```

After the threshold (5 by convention; 3–10 is reasonable), the message goes to the DLQ. The DLQ is a parking lot for human investigation. **Critical:** alert on DLQ depth. A growing DLQ = something's broken; a silent DLQ = an outage you don't know about yet.

### Why this is the lesson

> A DLQ without an alarm is a graveyard.

The DLQ is not a fix; it's a quarantine. A growing DLQ is a system telling you "something I can't handle on my own." Whoever's on-call for the queue is on-call for the DLQ.

Bonus: write a small "replay" tool that reads from the DLQ, lets a human inspect each message, and either re-enqueues or discards. You'll need it.

---

## The takeaway

Three rules:
1. **Ack after work, never before.**
2. **Idempotent consumers are mandatory.** At-least-once + idempotent = "effectively once."
3. **Every queue gets a DLQ. Every DLQ gets an alarm.**
