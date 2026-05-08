# AGENTS.md

This repository is a **portable tutor skill for backend engineering**. It runs in any tool-using agent — Claude Code, OpenAI Codex CLI, GitHub Copilot CLI, Cursor, Aider, etc.

**To run the tutor:** read `SKILL.md` in this directory and follow the protocol. The protocol is tool-agnostic prose ("read the file at X", "write to Y", "run the command Z"); translate to your harness's primitives. State lives entirely as files in the workspace — no MCP server, no database.

**Workspace location:** `~/backend-dev/`. Course state (progress, session-state, notes, exercises, projects, flashcards, reviews) lives there. This repo is the source-of-truth for the skill itself.

**Three lanes**, picked from a vibe check at first invocation (Step 2 of `SKILL.md`):

- **Foundations** — true entry. Open with a win, then light vocab probes.
- **Working** — the standard middle. 12-question diagnostic across the tier tree. Adjacent-domain variant for learners coming from frontend / mobile / SRE / data eng / ML / hardware.
- **Senior** — gap-fill against their actual project. 6-question diagnostic, harder, with a re-diagnostic affordance.

**Two orientations**, picked at Step 2.5 of `SKILL.md`:

- **`foundations_first`** — walk T0 → T11 in order. Reference files: `references/curriculum.md`, `references/exercise-bank.md`, `references/theory-modes.md`.
- **`builder_first`** — 10 hands-on coding loops. Reference: `references/builder-first.md`. When this orientation is picked, copy `assets/builder-first/<language>/` (or `assets/builder-first/_spec-only/` if no scaffolding exists for the chosen language) into the workspace's `projects/`.

**Reference files** in `references/` and assets in `assets/` are loaded **on demand** by the protocol — don't preload them.

**Slash commands** the user may type (parse and dispatch per `SKILL.md`):

| Command | When |
|---|---|
| `/plan`, `/start`, `/quiz`, `/continue`, `/notes`, `/config` | Anytime |
| `/loop list`, `/loop [n]`, `/loop skip`, `/loop quickpass` | Builder-first only |

Plain English also works: *"teach me X"*, *"build a Y"*, *"make this easier"*, *"pause"*. Full reference at `assets/COMMANDS.md` (copied to the workspace as `~/backend-dev/COMMANDS.md` at first-time onboarding).

**Sibling skills** that may be installed alongside this one:

- **system-design-tutor** — owns architecture-at-scale design ("design Twitter for 100M users"). Hand off when the request is at that scope.
- **ai-systems-tutor** — owns LLM-specific infrastructure (agent loops, RAG, prompt caching). Hand off when the request is about LLM internals. *Calling* an LLM from a backend service is in scope here.

**Default code language:** `go` (full prefilled scaffolding). `python` (FastAPI) is the planned secondary. Node/TS, Java/Kotlin, Rust, and others are spec-only — the learner implements against the spec; the tutor reviews. Set on `learner.language` in `progress.json`.

For installation across harnesses, see `CLAUDE.md` → *How the skill is meant to be installed*.
