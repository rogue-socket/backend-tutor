# backend-tutor commands

Slash commands and natural-language overrides. Both forms work — pick whichever you remember.

---

## Slash commands

| Command | What it does |
|---|---|
| `/plan` | Show the full curriculum + your current position. Use when you want a map. |
| `/start [topic]` | Start a lesson. With no argument, starts the next planned topic. With an argument, jumps to that topic (warns about prereqs). |
| `/quiz` | Run a spaced-repetition review of items due today. Capped at 15 items per session. |
| `/continue` | Resume from `session-state.md`. The default "I'm back" command. |
| `/notes [topic]` | Generate or update reference notes for a topic. With no argument, uses the current topic. |
| `/config` | Show or edit your learner profile in `progress.json` (level, language, orientation, stated goals). |

### Builder-first only

| Command | What it does |
|---|---|
| `/loop list` | Print all 10 loops with current status. |
| `/loop [N]` | Jump to loop N. Warns if prereqs aren't met; honors override. |
| `/loop skip` | Skip the current loop after a 30-second summary. Marks `skipped` in progress.json. |
| `/loop quickpass` | 3 quiz questions on the loop's WIN criteria. Pass = mark `done` without running it. |
| `/loop reset N` | Wipe loop N's project dir and start fresh (asks for confirmation). |

---

## Natural language overrides

These work anywhere; they outrank the tutor's proposal:

| You say | Tutor does |
|---|---|
| "yes" / "ok" / "let's go" / "continue" | Execute the proposal |
| "teach me X" / "build Y" / "design Z" | Honor the detour; queue current proposal for next time |
| "quiz me" / "review first" | Run review session |
| "pause" / "stop for today" / "I have to go" | End-of-session protocol; clean checkpoint |
| "give me notes" / "write this up" / "summarize this topic" | Generate notes for current topic |
| "what's the plan?" / "where are we?" | Show curriculum position |
| "this is too basic" | Re-route up a lane; depth-check first |
| "I'm drowning" | Re-route down a lane; concrete picture |
| "you routed me to the wrong place" | Acknowledge; offer re-diagnostic; don't argue |

---

## When stuck

- "I'm stuck on this exercise" → tutor walks the hint ladder (smallest hint that unblocks)
- "this is taking too long" → tutor pair-writes the next few lines, then hands control back
- "I want to skip this exercise" → fine; flag what concept you're skipping so it can come back via SR

---

## When you want to bend the rules

- *"Skip the diagnostic"* — works at first invocation, but the tutor will ask one or two probes before each lesson instead. Calibration has to happen somewhere.
- *"Just give me the answer"* — works for review queue items if you genuinely already know them; otherwise the tutor will push back. The wrestle is the lesson.
- *"Speed up"* — the tutor will compress explanations and cut visuals. Say it again if it's still too slow.
- *"Slow down"* — tutor will pause more, ask more comprehension checks, lean on visuals.

---

## What the tutor won't do

- Write the code for you in a builder-first loop
- Skip the BREAKs in builder-first loops (the breaks *are* the lessons)
- Generate notes silently — you ask or are offered
- Suggest a context reset before saving state
- Pretend you got something right when you didn't

If any of these happen anyway, that's a bug — call it out and it'll get fixed.
