# Session control

Multi-session protocol: how to start, pause, resume, and manage context across days/weeks. The learner's progress lives in `~/backend-dev/`. Two files matter most: `progress.json` (the durable record) and `session-state.md` (the resume pointer).

---

## `session-state.md` schema

A single markdown file at `~/backend-dev/session-state.md`. Overwritten at every checkpoint. Keep it short — under 50 lines. The skill reads this at the start of every Warm Resume session.

```markdown
# Session state

**Last updated:** 2026-05-08 14:32
**Last session:** 2026-05-08

## Where we left off

[1-2 sentences. "Mid-lesson on T2.3 transactions and isolation. Just finished walking through write-skew with a worked example; about to start SELECT FOR UPDATE / SKIP LOCKED for queue patterns." Be specific enough that another harness could resume cold.]

## What's next

1. [Next planned step]
2. [Step after that]
3. [Step after that]

## Active review queue (top 5 due)

- [topic] — due 2026-05-08 — interval 3d, ease 2.4
- [topic] — due 2026-05-09 — interval 7d, ease 2.6
- ...

(Internally this is the spaced-repetition queue — `sr_queue` in `progress.json`. Use "review queue" in user-facing language.)

## Active builder loop (if `orientation = builder_first`)

- Loop: [N] — [name]
- State: [in-progress | broken | done]
- Path: `~/backend-dev/projects/loop-N-<slug>/`
- Last action: [1 sentence]

## Open threads

- [Anything the learner asked to revisit later]
- [Any partial exercise in `projects/<slug>/` that wasn't finished]

## Diagnostic edge

[The "particular gap" line from the diagnostic, kept current as gaps get filled.]
```

---

## When to checkpoint

Update `session-state.md` AND `progress.json` whenever:
- A lesson finishes
- An exercise or builder loop step finishes
- The user signals pause ("stop for today", "I have to go", `/pause`)
- 30+ minutes pass without a checkpoint
- You're about to suggest context compaction or a new chat
- Major mode switch (theory → practical, study → mock interview)

**State writes always happen BEFORE you suggest the user run `/compact`, start a new chat, or close the session.** A checkpoint after the conversation is gone is useless.

---

## Warm Resume protocol (full version)

When `~/backend-dev/session-state.md` exists, this is the standard "user is back" flow.

### Step 1: Read state

Read `progress.json` and `session-state.md`. Note:
- Days since last session
- Mid-lesson, mid-loop, or mid-exercise state
- Top of review queue (overdue items)
- Open threads

### Step 2: Pick a proposal

Priority order:
1. **Mid-lesson / mid-loop from <14 days ago** → resume that exact thing
2. **Review queue has overdue items** → propose review session, then continue forward
3. **Clear next curriculum step** → announce it
4. **Gap is 14+ days** → propose a brief review of the last topic before moving forward

### Step 3: One-paragraph welcome-back

Format:

> "Welcome back. Last time we [where we left off]. Today: [resume X], then [next planned step]. Review queue has [N] items due — let's knock those out first. Sound good?"

Max 4 lines. **Don't ask "what would you like to work on?"** — propose. The user can override.

### Step 4: Wait for confirmation or override

| User says | Do |
|---|---|
| "yes" / "ok" / "continue" / "let's go" | Execute the proposal |
| "actually let's do Y" | Honor the override; queue your proposal for next time |
| "I have to leave in 10 min" | Run shortened version — quick review, no new lesson |
| "where am I?" | Show full curriculum position (read from `progress.json`) |

### Step 5: Execute

Once they confirm, **don't re-explain or preamble**. Begin the lesson / review / loop step immediately.

---

## Pause protocol (mid-session)

Triggered by: "pause", "stop", "I have to go", "let's pick up tomorrow", `/pause`, or implicit signals (long silence, a dropped reply).

1. Acknowledge in one sentence. Don't summarize the whole session.
2. Update `session-state.md`. Be specific about *exactly* where you stopped.
3. Update `progress.json`:
   - If a topic was being taught: bump `last_reviewed`, update `confidence` and `weak_points`
   - If an exercise / loop was in flight: log it as `in_progress` with the directory path
   - Add an entry to `sessions.entries` with date, duration if known, topics touched
4. Tell the user where state is saved: "Saved to `session-state.md`. Run `/continue` next time."
5. Do NOT continue the lesson. Stop cleanly.

---

## End-of-session protocol (clean finish)

Triggered when: a topic naturally completes and the user signals ready-to-stop, or when you propose stopping (e.g. "good place to pause for today").

1. **Notes pass.** For any topic covered in this session that doesn't have `notes/<topic-slug>.md`, offer to generate notes. Don't generate without offering.
2. **Quick review pass** (optional, ~3 min). Pick 2-3 review-queue items due in the next 2 days and quiz them. This bakes today's material in.
3. **Checkpoint** as above.
4. **Preview next session.** One line: "Next: [topic]." Goes into `session-state.md`.
5. Sign off.

---

## Cold Resume (Case B — workspace exists, session-state missing)

Workspace exists but `session-state.md` is missing or unreadable. Skip workspace setup, but you still need to know where to start.

1. Read `progress.json`.
2. If it has meaningful progress (any topic with status `in_progress` or `solid`+):
   > "I have your `progress.json` but no session pointer — looks like the session-state file was deleted or this is a new harness. You're at [tier/topic from progress.json]. Want me to pick up there, or do a quick recalibration first?"
3. If `progress.json` is near-empty (just an init): run **diagnostic-lite** — 3-4 questions across whichever tiers had any topic touched — then propose a starting point and write a fresh `session-state.md`.

---

## Context-window awareness

Long sessions accumulate noise. The model gets slower and dumber as context fills. Proactively offer to checkpoint and reset when:

- 60+ messages in and about to start a new sub-topic
- Long debugging session is over and you're moving to new material
- Major mode switch (theory → practical, or study → mock interview)
- The conversation has visibly degraded (model repeating itself, missing references)

### Reset commands by harness

- **Claude Code:** `/compact` (preserves a summary, drops detail)
- **Codex / Codex CLI:** start a new task; `~/backend-dev/session-state.md` is the bridge
- **Copilot CLI:** new session; same bridge
- **Cursor / Aider:** new chat; same bridge
- **Claude.ai:** "summarize this conversation, then I'll start a new chat"

### Pre-reset protocol (always the same)

1. Stop the lesson.
2. Update `session-state.md` with extra detail — include any open question, half-baked example, decisions made this session.
3. Update `progress.json`.
4. Tell the user: "I'm going to suggest a context reset. State is saved. Run `/compact` (or start a new session). Then say `/continue` and we'll pick up clean."
5. Wait for them to do it. Don't continue talking.

---

## Multi-harness handoff

The user might run a session in Claude Code today, Codex tomorrow. The bridge is the workspace.

- Anything written to `~/backend-dev/` is portable.
- Nothing harness-specific (CC's transcript, Codex's history, etc.) crosses over.
- `session-state.md` is the contract: any harness reading it can resume.

If the user says "I was working on this in [other harness]" and `session-state.md` doesn't reflect it, that means the prior harness didn't checkpoint cleanly. Read `progress.json` for the durable state, ask the user 1-2 questions to recover the in-flight context, then write a fresh `session-state.md`.

---

## Anti-patterns

- ❌ Suggesting `/compact` before writing state
- ❌ "Let me summarize what we did today" — the user can read the diff / their own scrollback
- ❌ Asking "what would you like to do next session?" at end-of-session — the skill picks
- ❌ Long welcome-back messages — 4 lines max
- ❌ Re-explaining the curriculum every session
- ❌ Skipping checkpoints because "we'll do it at the end"
- ❌ Writing to `progress.json` without reading it first (you'll clobber prior state)
