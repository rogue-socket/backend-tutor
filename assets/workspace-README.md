# Backend Dev — your workspace

This directory holds your personal progress for the **backend-tutor** course. Everything here is yours — notes, exercises, projects, flashcards. The tutor reads and writes these files between sessions.

## Layout

```
~/backend-dev/
├── README.md                ← you are here
├── COMMANDS.md              ← slash commands + override reference
├── progress.json            ← long-term progress (topics done, weak points, review queue)
├── session-state.md         ← where you left off (overwritten each session)
├── notes/                   ← topic notes the tutor generates ("/notes <topic>")
│   └── diagrams/            ← interactive HTML diagrams
├── cheatsheets/             ← per-topic quick reference cards
├── projects/                ← builder-first loops + standalone exercises live here
│   └── loop-N-<slug>/       ← one folder per builder-first loop
├── reviews/                 ← mock interview / design review writeups
├── flashcards/              ← optional — exported flashcard decks per topic
├── viewer/                  ← static workspace viewer (see "Browse your workspace")
└── meta/                    ← anything else the tutor wants to keep around
```

## Two state files that matter

**`progress.json`** — the durable record. What topics you've covered, your confidence level, your review queue, your loop progress, your stated goals. The tutor reads this every session. Don't hand-edit unless you know what you're doing; use `/config` to update learner profile fields.

**`session-state.md`** — the resume pointer. Short markdown file the tutor overwrites at every checkpoint. Tells the next session exactly where you left off. Human-readable; you can read it cold to remember where you were.

## Common commands

- `/plan` — show the full curriculum + your current position
- `/start [topic]` — start a lesson on a topic (or the next planned one)
- `/quiz` — run a spaced-repetition review of due items
- `/notes [topic]` — generate or update notes for a topic
- `/continue` — resume from `session-state.md`
- `/config` — show or edit your learner profile (level, language, goals)

**Builder-first only:**
- `/loop list` — show all 10 loops + status
- `/loop N` — jump to loop N
- `/loop skip` — skip the current loop (with a 30-second summary)
- `/loop quickpass` — pass the loop with 3 quick questions

Plain English works too: *"teach me caching"*, *"review my schema"*, *"pause for today"*. Full reference at `COMMANDS.md`.

## Pausing and resuming

Just say "pause" or "I have to go." The tutor checkpoints state and stops cleanly. Next session, run `/continue` (or just say hi) and the tutor picks up exactly where you left off — even if you're using a different harness (Claude Code, Cursor, Codex CLI, etc.).

## Multi-harness portability

Everything in `~/backend-dev/` is portable across tool-using agents. If you started in Claude Code and want to continue in Codex CLI, the workspace is the bridge. The tutor reads `session-state.md` regardless of which harness it's running in.

## When something goes sideways

- **`session-state.md` got deleted** — tutor will read `progress.json` and propose a recalibration
- **`progress.json` got corrupted** — keep a backup before editing manually; failing that, the tutor can rebuild from your notes + brief recalibration
- **You skipped 3 weeks** — run `/quiz` first; the SR queue will be hungry
- **You want to switch language or orientation mid-course** — say so; the tutor will handle the migration cleanly

## Browse your workspace

Your notes, cheatsheets, and flashcards are plain markdown — readable anywhere. The bundled `viewer/` is a tiny static site that renders them with a sidebar so you can skim everything at once.

**One-time setup:** `viewer/` is already in your workspace. Whenever you've added or edited files, regenerate the manifest:

```bash
cd ~/backend-dev
python3 viewer/regenerate-manifest.py
```

**Serve it:**

```bash
cd ~/backend-dev
python3 -m http.server 8000
# open http://localhost:8000/viewer/
```

The viewer reads `viewer/manifest.json` and fetches markdown files relative to the workspace root. Markdown rendering uses `marked` from a CDN, so the first load needs internet; after that it's cacheable. If `manifest.json` is missing, the viewer tells you to run the regenerate script.

## Privacy

This is a local directory. Nothing here is sent anywhere unless you choose to share it. The tutor reads what's in here as part of running sessions; that's it.
