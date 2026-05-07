# Loop 5 — Async work via a queue

**Tier mapping:** T3.6 (queues), T3.7 (delivery semantics), T3.8 (retries), T3.9 (DLQs), T1.13 (webhooks shape)
**Time:** 120–180 minutes
**Theme:** *some things shouldn't happen on the request path.*
**Prereqs:** Loops 2 + 4 done.

## What you're building

When a user creates a link, send them a "link created" notification (email-shaped — but we'll just log it; SMTP is a distraction). Move that work off the request path onto a **background worker** consuming from **Redis Streams**.

Three properties matter:
1. **The HTTP request returns fast** — the worker does the work asynchronously
2. **At-least-once delivery** — messages survive worker crashes
3. **Idempotent consumer** — duplicate deliveries don't cause duplicate side effects

Plus a **dead-letter queue** for poison messages: if the worker fails on the same message 5 times, route it to the DLQ instead of looping forever.

## Architecture

```
HTTP POST /links
    ↓
Handler: insert link → enqueue notification → return 201
    ↓
Redis Stream "notifications"
    ↓
Worker process: XREADGROUP → process → XACK
                           ↘ on repeated failure → DLQ stream
```

## New files

- `cmd/server/main.go` — the HTTP service (Loop 4 + the enqueue call)
- `cmd/worker/main.go` — the background worker, runs as a separate process
- `internal/queue/queue.go` — enqueue / consume / ack / DLQ helpers
- `docker-compose.yml` — adds Redis to the existing Postgres setup

The split into `cmd/server` and `cmd/worker` is intentional — they're different processes with different lifecycles. The worker will scale independently; the HTTP service must not block on the worker.

## Tasks

1. **Add Redis to docker-compose.** `redis:7-alpine` on `:6379`.
2. **Implement `Enqueue`** — `XADD notifications * type "link.created" link_id $1 user_id $2`.
3. **Implement the worker loop** — `XREADGROUP GROUP workers $consumer COUNT 10 BLOCK 5000 STREAMS notifications >`.
4. **Process each message:**
   - Look up `link_id`, format the "email," log it.
   - On success → `XACK notifications workers <id>`.
   - On failure → don't ack; let it redeliver via XPENDING.
5. **Implement idempotent dedup.** Track processed IDs in a Redis set with a 24h TTL. Before processing, `SISMEMBER`; after success, `SADD`. (Alternative: store in Postgres — slightly stronger, slightly slower.)
6. **Implement DLQ routing.** Use `XPENDING` to find messages with delivery count ≥ 5; `XADD notifications.dlq` then `XACK` the original. Add a metric / log on every DLQ event.
7. **Run BREAKs 1, 2, 3.** See `BREAK.md`.

## Stretch

- Implement consumer-group rebalancing: kill one worker, watch another pick up its in-flight messages via XAUTOCLAIM.
- Add a "delayed delivery" mechanism — message with a `not_before` timestamp, worker checks before processing, re-queues with delay if not ready. (You'll roll your own; Redis Streams doesn't have built-in delays.)
- Write a tiny dashboard: `GET /admin/queue` returns stream length, consumer pending count, DLQ count.

## Why Redis Streams over RabbitMQ / Kafka / SQS

Loop 5 prioritises *learning* over *production resilience*. Streams give you the right concepts (consumer groups, ack/nack, redelivery, dead-letter pattern) with one binary, no operator, free in dev. The semantic surface is similar enough to Kafka that what you learn here transfers.

For real production work at scale, you'd likely use Kafka (high throughput, strong durability) or SQS (managed, simple). But that's Loop 10 territory, not Loop 5.
