# WIN — Loop 5 done

## Code

- [ ] Server enqueues notifications via `XADD`; HTTP returns < 50ms even when the worker is offline
- [ ] Worker is a separate `cmd/worker` process; reads via `XREADGROUP`; acks after successful processing
- [ ] Idempotent dedup (Redis set or Postgres `processed_messages` table)
- [ ] Max-deliveries (5) enforced; messages exceeding it routed to `notifications.dlq`
- [ ] Worker handles SIGTERM gracefully — drains in-flight, closes Redis client, exits clean
- [ ] Metrics / logs for: enqueue count, process success, process failure, DLQ count, queue depth

## Verification

- [ ] BREAK 1 reproduced: ack-before-work loses messages on crash; ack-after-work doesn't
- [ ] BREAK 2 reproduced: chaos-hooked worker produces duplicate side effects without dedup; zero with dedup
- [ ] BREAK 3 reproduced: a poison message lands in the DLQ after 5 deliveries; an alert (log line OK for now) fires

## Understanding

1. **At-least-once vs at-most-once vs "exactly-once" — name a real-world workload for each.**
   *Outline: at-least-once = the right default for almost everything. Notifications, webhook delivery, billing events. Idempotent consumer required. At-most-once = telemetry pings, game state updates where stale > duplicate. Lossy is acceptable. "Exactly-once" = doesn't exist over a network; what people mean is at-least-once + transactional consumer (e.g., Kafka EOS within a single Kafka topic-partition + idempotent producer). Outside that narrow case, claims of exactly-once are either lies or marketing reframings of "at-least-once + dedup."*

2. **You ack after work. The worker crashes mid-XACK. What happens?**
   *Outline: the message stays in XPENDING with the original consumer's name. After `min-idle-time` (configured at consumer-group level), XAUTOCLAIM (or another worker calling XREADGROUP with `>`) picks it up. The new worker re-runs `process`. Idempotent dedup catches it. Net: at-least-once delivery preserved; correctness preserved if dedup is correct.*

3. **The DLQ has 47 messages and growing. What's your investigation order?**
   *Outline: (1) Inspect a sample — what's common? Same user? Same link_id? Same error message? (2) If a code bug — fix the bug, replay the DLQ. (3) If bad data — purge the offenders, add a producer-side validation. (4) If a downstream outage — pause replay until the downstream is healthy. (5) Always: post-incident, audit whether the alarm fired before the user noticed.*

## Reflection

What surprised you? Common ones:
- The dedup set is the load-bearing thing; the queue is just the transport.
- Crash-during-ack is a real and frequent scenario, not a theoretical one.
- DLQs are easy to build and easy to forget about.

## What's next

Loop 6 — Caching + stampede. Redis cache-aside on a hot read endpoint, the cache stampede in action, and single-flight as the canonical fix.
