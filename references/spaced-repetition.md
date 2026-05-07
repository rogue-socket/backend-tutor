# Spaced repetition

The skill maintains an SR queue in `progress.json`. Items get rescheduled by an SM-2-lite algorithm based on the learner's response quality. Daily review nudges questions back to long-term retention.

User-facing language: **"review queue"**. Internal field name: `sr_queue`.

---

## When items enter the queue

- **Diagnostic miss.** The learner answered a diagnostic question wrong or shakily — that question becomes a review item, due in 1 day.
- **Auto-quiz miss.** Mid-lesson quiz answered wrong — review entry, due in 1 day.
- **Lesson key fact.** At the end of a lesson, the skill picks 2-3 "must-remember" facts (the ones the lesson hinged on). They become review items, due in 2 days.
- **Self-flagged.** The learner says "I want to remember this" — review entry, due in 1 day.
- **Reflection insight.** A surprise from `reflection.md` after an exercise becomes a review item, due in 3 days.

Don't let the queue grow without bound. Cap at ~50 active items per tier. If the learner is consistently missing the same item across 4+ reviews, it's a foundational gap — surface it as a topic to re-teach, not just a review item to retry.

---

## SR item schema

Stored at `progress.json` → `sr_queue.items[]`:

```json
{
  "id": "T2.3.transactions::write-skew-vs-phantom",
  "topic": "T2.3.transactions",
  "question": "Postgres' default isolation level (Read Committed) prevents which anomalies, and which does it let through?",
  "answer_outline": "Prevents: dirty reads. Lets through: non-repeatable read, phantom read, write skew, lost update. Snapshot isolation (Repeatable Read in Postgres) catches non-repeatable + phantom but still allows write skew. Only Serializable closes write skew.",
  "interval_days": 1,
  "ease": 2.5,
  "due": "2026-05-09",
  "history": [
    { "date": "2026-05-08", "rating": "good", "interval_after": 3 }
  ]
}
```

**Fields:**
- `id` — `<topic-id>::<short-slug>`. Stable.
- `topic` — points back into curriculum. `<tier>.<section>.<topic-slug>`.
- `question` — what the skill asks the learner.
- `answer_outline` — bullets, not prose. Skill uses this to grade the learner's response.
- `interval_days` — current interval. New items start at 1.
- `ease` — SM-2 ease factor. Starts at 2.5.
- `due` — date to ask next.
- `history` — append-only log: `{date, rating, interval_after}`.

---

## SM-2 lite scheduler

After the learner answers, the skill rates the response on a 4-level scale:

| Rating | Meaning |
|---|---|
| `again` | Wrong, or close to wrong. |
| `hard` | Got it but with hesitation or partial recall. |
| `good` | Correct, with reasonable speed. |
| `easy` | Instant, complete, no friction. |

### Update rules

```
if rating == "again":
    interval_days = 1
    ease = max(1.3, ease - 0.2)
elif rating == "hard":
    interval_days = max(1, round(interval_days * 1.2))
    ease = max(1.3, ease - 0.15)
elif rating == "good":
    interval_days = round(interval_days * ease)
    # ease unchanged
elif rating == "easy":
    interval_days = round(interval_days * ease * 1.3)
    ease = ease + 0.15

due = today + interval_days
history.append({ date: today, rating, interval_after: interval_days })
```

**Worked example.** New item, ease 2.5, interval 1.
- Day 1, rated `good` → interval = round(1 × 2.5) = 3, due day 4.
- Day 4, rated `hard` → interval = max(1, round(3 × 1.2)) = 4, ease = 2.35, due day 8.
- Day 8, rated `again` → interval = 1, ease = 2.15, due day 9.
- Day 9, rated `good` → interval = round(1 × 2.15) = 2, due day 11.
- Day 11, rated `good` → interval = round(2 × 2.15) = 4, due day 15.
- Day 15, rated `easy` → interval = round(4 × 2.15 × 1.3) = 11, ease = 2.30, due day 26.

---

## Daily review session

Triggered by `/quiz`, by an offer at session start ("review queue has 4 items due — review first?"), or as part of an end-of-session pass.

1. Read `progress.json`. Filter `sr_queue.items` where `due <= today`.
2. Sort by `due` ascending (oldest first), then by `interval_days` ascending (shorter intervals first — they need more reps).
3. Cap at **15 items per session**. More than that and the learner zones out.
4. Ask each question. After their answer:
   - Reveal the `answer_outline`.
   - Rate the response (or ask "again / hard / good / easy?").
   - Apply scheduler rules. Write back to `progress.json`.
5. After the queue is empty (or capped), report: "5 reviewed: 3 good, 1 hard, 1 again. Next due 2026-05-11."

---

## Grading the learner's response

The skill is the grader. Rules:
- **Correct on the first attempt, no hedging** → `good` (or `easy` if it was instant and complete).
- **Correct after a hint or after partial recall** → `hard`.
- **Wrong on a key part** → `again`. Don't be charitable here — partial wrong is still wrong; better to over-review than to let a gap settle.
- **Right answer, wrong reasoning** → `again`. The reasoning is the thing being trained.

After grading, briefly explain *why* their answer landed where it did. Review is only useful if the learner sees their gap.

---

## `progress.json` — full schema

```json
{
  "started": "2026-05-07",
  "learner": {
    "level": "working",
    "working_mode": "standard",
    "orientation": "builder_first",
    "language": "go",
    "stated_goals": ["passing senior-backend interview loop in 6 weeks"],
    "stated_context": ["currently shipping FastAPI services at $JOB", "wants to pivot to Go"]
  },
  "topics": {
    "T0.3": {
      "status": "solid",
      "confidence": 4,
      "last_reviewed": "2026-05-07",
      "weak_points": []
    },
    "T2.3": {
      "status": "in_progress",
      "confidence": 2,
      "last_reviewed": "2026-05-08",
      "weak_points": ["write-skew vs phantom read distinction"]
    }
  },
  "sr_queue": {
    "items": [ /* see SR item schema above */ ]
  },
  "exercises": {
    "entries": [
      {
        "date": "2026-05-08",
        "topic": "T2.5.migrations",
        "dir": "exercises/2026-05-08-online-migration",
        "status": "completed",
        "type": "build-from-scratch",
        "takeaways": ["expand/contract requires both old and new code paths to coexist for a release"]
      }
    ]
  },
  "loops": {
    "current": 3,
    "entries": [
      { "loop": 1, "name": "bare CRUD", "status": "done", "completed": "2026-05-07" },
      { "loop": 2, "name": "persistence + concurrency break", "status": "done", "completed": "2026-05-08" },
      { "loop": 3, "name": "migrations + N+1", "status": "in_progress", "started": "2026-05-08" }
    ]
  },
  "sessions": {
    "entries": [
      {
        "date": "2026-05-07",
        "duration_min": 45,
        "topics_touched": ["T0.3", "T1.2"],
        "loops_touched": [1]
      }
    ]
  }
}
```

**Status values for topics:** `unknown` | `weak` | `in_progress` | `solid` | `mastered`. `mastered` is rare — only after 3+ successful SR reviews at intervals >14 days.

**Confidence:** 1-5, learner's self-rating, used as a tiebreaker for what to review next.

---

## Anti-patterns

- ❌ Asking the same review item the same day it was generated. Min interval is 1 day.
- ❌ Rating everything `good` because the learner is "trying their best". Rate honestly; the schedule depends on it.
- ❌ Letting the queue exceed ~150 items. Above that, the daily review becomes punishing and gets skipped.
- ❌ Auto-deleting items after 5 successful reviews. Long-tail items still rot — keep them, just at long intervals.
- ❌ Not writing back to `progress.json` after a review. The schedule is useless if it doesn't persist.
- ❌ Adding review items for trivial details. Review is for things that matter and are forgettable. "What's the third argument to `pgx.Connect`?" is not review-worthy.

---

## When SR isn't the right tool

Some material doesn't fit SR:
- **Procedural skills** (writing a JSON HTTP handler, configuring autovacuum): repetition through exercises and builder loops, not flashcards.
- **Trade-off judgment** ("when do I pick DynamoDB vs Postgres?"): live discussion / mock interview, not SR.
- **Architectural intuition**: cumulative across many lessons; SR catches pieces but misses the gestalt.

For these, fall back to mock interviews, design reviews, and builder loops. SR is a complement, not a substitute.
