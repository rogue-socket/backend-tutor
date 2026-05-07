# BREAK — Loop 2

Two breaks here. Run them in order; both are foundational.

---

## BREAK 1 — Connection pool exhaustion

**Foundations filled here:** T2.7 (connection pooling), T7.4 (capacity).

### The setup

You've configured `pgxpool` with `MaxConns = 2`. (In real life you'd start higher; we're shrinking it to make the failure visible.) Now run a load script that fires 50 concurrent slow queries:

```bash
hey -n 50 -c 50 'http://localhost:8080/links'
# or use any concurrent client; even a quick bash:
for i in $(seq 1 50); do curl -s 'http://localhost:8080/links' > /dev/null & done
```

To make queries "slow", temporarily add `pg_sleep(2)` to the List query, OR add an artificial 2s delay in the handler before the DB call.

### What you'll see

The first 2 requests get connections immediately. The next 48 wait. Some succeed once the first 2 release; others time out (default `pgxpool` connection acquisition timeout is 30s, but request timeout will fire first).

The application sees this as **request timeouts**, not as DB errors. From the application logs you'd think the DB was slow. From the DB's logs you'd see only 2 connections, looking lazy.

### Why this is the lesson

In production, pool exhaustion is a top-5 outage cause. Symptoms:
- Application latency spikes
- DB looks underutilized (low CPU, few connections)
- Restarts "fix" it briefly, then it returns

Mitigations (in order):
1. **Right-size the pool.** Per HikariCP's docs: `pool_size = (cores * 2) + effective_spindle_count`. Postgres's *server-side* default `max_connections` is 100 — your *client-side* pool typically wants to be a small fraction of that times your replica count.
2. **Acquisition timeout.** Set it shorter than your request timeout so the failure is "couldn't get a conn" not "request hung."
3. **Query timeouts.** A long-running query that holds a connection is what causes pool starvation in the first place.

### Restore

Once you've felt it, set `MaxConns = 10` and remove the `pg_sleep`.

---

## BREAK 2 — N+1 on tags

**Foundations filled here:** T2.6 (ORMs and N+1).

### The setup

You implemented `ListWithTagsNaive` per TODO 10 — for each link, a separate query for its tags. Run the integration test:

```bash
go test -run TestListWithTagsIsNotNPlusOne -v ./...
```

It seeds 10 links with 3 tags each, then calls the list endpoint while counting queries. The naive impl emits **11 queries** (1 for the list + 1 per link).

### What you'll see in logs

If you've enabled query logging (`pgxpool` `Tracer` hook):

```
SELECT id, url, title FROM links ORDER BY id LIMIT 50
SELECT t.name FROM tags t JOIN link_tags lt ON lt.tag_id = t.id WHERE lt.link_id = $1 -- $1=1
SELECT t.name FROM tags t JOIN link_tags lt ON lt.tag_id = t.id WHERE lt.link_id = $1 -- $1=2
SELECT t.name FROM tags t JOIN link_tags lt ON lt.tag_id = t.id WHERE lt.link_id = $1 -- $1=3
... (×N)
```

That's the N+1 query pattern. Each round-trip to the DB has fixed cost (~0.5ms locally, 5–50ms across a network); 100 round-trips = 50ms–5s, before any actual work.

### The fix

Two equally good approaches:

**Option A — single query with JOIN + aggregation:**
```sql
SELECT l.id, l.url, l.title, COALESCE(array_agg(t.name) FILTER (WHERE t.name IS NOT NULL), '{}')
FROM links l
LEFT JOIN link_tags lt ON lt.link_id = l.id
LEFT JOIN tags t       ON t.id = lt.tag_id
GROUP BY l.id
ORDER BY l.id
LIMIT 50
```

**Option B — two queries, batch the second:**
```sql
-- 1: get the link IDs
SELECT id, url, title FROM links ORDER BY id LIMIT 50;
-- 2: get all tags for those IDs in one shot
SELECT lt.link_id, t.name
FROM link_tags lt JOIN tags t ON t.id = lt.tag_id
WHERE lt.link_id = ANY($1)
```
Then assemble in Go. Sometimes simpler than the JOIN+GROUP, especially when the per-link list is large.

Both produce 1–2 queries instead of N+1.

### Why this is the lesson

ORMs make N+1 *trivial* to write (`for link in links: link.tags.all()`). Many production codebases have N+1 lurking under code that looks innocuous. The cure is:

1. **Query logging in dev** — see every query, count surprises
2. **Eager-load by default** for known relations
3. **Tests that assert query counts** for hot paths (the test above)

This loop's takeaway: **once a list endpoint exists, every relation is a potential N+1.** Plan the load pattern before you write the code.
