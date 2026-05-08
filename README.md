# backend-tutor

A portable, agent-driven course on **backend engineering** — runs in Claude Code, Codex CLI, Copilot CLI, Cursor, Aider, anywhere with file-read + file-write + shell.

![Language](https://img.shields.io/badge/language-Go%20%7C%20Python%20%7C%20Node%20%7C%20Java%20%7C%20Rust-blue) ![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey) ![License](https://img.shields.io/badge/license-MIT-green) ![Loops](https://img.shields.io/badge/builder--first-10%20loops-orange) ![Status](https://img.shields.io/badge/status-v0.1.0-yellow)

The skill OWNS the curriculum: it onboards you, drives lessons, schedules spaced-repetition reviews, runs hands-on exercises, and checkpoints state across sessions. You steer when you want a detour; the default is forward motion.

Sibling of [system-design-tutor](https://github.com/anthropics/skills) (architecture-at-scale design) and [ai-systems-tutor](https://github.com/rogue-socket/ai-system-tutor) (LLM-specific infrastructure). Backend-tutor stays on the implementation/operational side: how to wire it up, how to read the query plan, how to debug it at 3am.

## Contents

- [What it covers](#what-it-covers)
- [Three lanes, two orientations](#three-lanes-two-orientations)
- [Builder-first loops](#builder-first-loops)
- [Requirements](#requirements)
- [Install](#install)
- [Usage](#usage)
- [Slash commands](#slash-commands)
- [Repo layout](#repo-layout)
- [Workspace layout (after first run)](#workspace-layout-after-first-run)
- [Status](#status)
- [License](#license)

## What it covers

Twelve tiers (T0–T11), end-to-end:

| | Tier | What |
|---|---|---|
| T0 | Networking & HTTP | TCP, TLS, HTTP/1.1 → HTTP/3, DNS, request anatomy |
| T1 | APIs | REST, gRPC, GraphQL, WebSockets, SSE, idempotency, versioning |
| T2 | Databases | Relational, document, KV, columnar, time-series, graph, search, vector — at the implementation level |
| T3 | Concurrency & async | Threads, goroutines, async/await, channels, locks, race conditions |
| T4 | Caching | Read-through, write-through, stampede mitigation, invalidation |
| T5 | Reliability | Retries, circuit breakers, timeouts, error budgets |
| T6 | Observability & on-call | Logs, metrics, traces, alerts, postmortems |
| T7 | Performance & scale | Query plans, capacity estimation, load testing |
| T8 | Security | OWASP API Top 10, authn/authz, secrets, input validation |
| T9 | DevOps adjacency | 12-Factor App, containers, deploys, config |
| T10 | Cloud literacy | Compute / queue / storage primitives across AWS / GCP / Azure |
| T11 | Distributed systems (impl-level) | Replication, sharding, consensus — code-it-and-break-it scope |

Anchored to engineering blogs (Stripe, Cloudflare, GitHub, Discord, Figma, Shopify, Netflix, Slack, Uber), OWASP API Security Top 10 (2023), the 12-Factor App, *Software Engineering at Google*, *Database Internals*, *Building Microservices 2e*, *API Design Patterns*, the SRE books, and the relevant HTTP RFCs.

## Three lanes, two orientations

After a short pattern-based vibe check, you're routed into one of three lanes:

- **Foundations** — true entry. Opens with a win, then light vocab probes, then a concrete first lesson.
- **Working** — the standard middle. 12-question diagnostic across the tiers; adaptive depth on each answer. Adjacent-domain variant for learners coming from frontend / mobile / SRE / data eng / ML / hardware.
- **Senior** — gap-fill against your actual project. 6-question diagnostic, harder, with a re-diagnostic affordance if the router got it wrong.

Then you pick how to walk the course:

- **Foundations-first** — walk T0 → T11 in order. Mental models first, then build on top. Slower start, sturdier base.
- **Builder-first** — Loop 1 ships a working CRUD service in your chosen language. Loop 2 breaks it on purpose (concurrency). Each loop adds a layer and breaks it. Foundations get filled in *as the service breaks*.

Both end up covering similar ground; the order differs.

## Builder-first loops

| # | What you build / break |
|---|---|
| 1 | Bare CRUD — net/http, no framework, single binary |
| 2 | Persistence — pgx + Postgres, observe the concurrency break |
| 3 | Online schema migrations — expand/contract, no downtime |
| 4 | Auth — sessions over JWT, revocation, middleware |
| 5 | Async work via queue — at-least-once delivery, idempotent consumers |
| 6 | Caching + stampede — single-flight, probabilistic early refresh |
| 7 | Containerize + deploy — 12-Factor config, multi-stage Dockerfile |
| 8 | Real-time + scale out — WebSockets, presence, fan-out |
| 9 | Observability + planted bug — traces beat logs, find the bug |
| 10 | Load test + capacity + postmortem — k6, p99 budget, write the writeup |

Default scaffolding language is **Go** (full prefilled starter code with TODOs). **Python (FastAPI)** is the planned secondary. **Node/TS, Java, Kotlin, Rust** are spec-only — you implement against the spec, the tutor reviews.

## Requirements

- A tool-using agent (Claude Code, Codex CLI, Copilot CLI, Cursor, Aider, ...) — anywhere with file-read + file-write + shell.
- Whatever language toolchain you pick at onboarding (Go default — install [from go.dev](https://go.dev/dl/)).
- For builder-first Loops 2+: Docker / Docker Compose for the Postgres / Redis / queue dependencies.

No API keys. No Python deps. Backend-tutor is content-only; everything you run is your own toolchain on your own machine.

## Install

See [`INSTALL.md`](INSTALL.md) for the full matrix (macOS / Linux / Windows × Claude Code / Codex CLI / Copilot CLI / Cursor / Aider).

Quick path — macOS / Linux + Claude Code:

```bash
git clone https://github.com/rogue-socket/backend-tutor ~/Documents/backend-tutor
ln -s ~/Documents/backend-tutor ~/.claude/skills/backend-tutor
```

Then `> start the backend tutor` in any Claude Code session. For Windows, Codex CLI, or other harnesses: see `INSTALL.md`.

### Platform branches

If you want a single-path README without conditionals, switch to the branch that matches your setup:

| Branch | OS | Harness |
|---|---|---|
| [`cc-windows`](https://github.com/rogue-socket/backend-tutor/tree/cc-windows) | Windows | Claude Code |
| [`codex-macos`](https://github.com/rogue-socket/backend-tutor/tree/codex-macos) | macOS / Linux | Codex CLI |
| [`codex-windows`](https://github.com/rogue-socket/backend-tutor/tree/codex-windows) | Windows | Codex CLI |

`main` covers every combo; the platform branches strip out the conditionals.

## Usage

```
> start the backend tutor
```

The skill auto-routes from the trigger phrase: vibe check → lane → orientation + language → `~/backend-dev/` setup → first lesson. Resume any time with `> /continue` — works across harnesses (the workspace at `~/backend-dev/` is the bridge).

See `INSTALL.md` for harness-specific invocations (Codex CLI, Copilot CLI, Cursor, Aider).

## Slash commands

| Command | When |
|---|---|
| `/plan` | Show the curriculum + your current position |
| `/start [topic]` | Start a lesson on a topic (or the next planned one) |
| `/quiz` | Run a spaced-repetition review of due items |
| `/notes [topic]` | Generate or update notes for a topic |
| `/continue` | Resume from `session-state.md` |
| `/config` | Show or edit your learner profile (level, language, goals) — switchable mid-course |
| `/loop list` | Builder-first only — show all 10 loops + status |
| `/loop [n]` | Builder-first only — jump to loop N |
| `/loop skip` | Builder-first only — skip the current loop with a 30-second summary |
| `/loop quickpass` | Builder-first only — pass with 3 quick questions |

Plain English also works: *"teach me caching"*, *"build a payments idempotency handler"*, *"make this easier"*, *"pause"*. Full reference at `assets/COMMANDS.md` (copied to `~/backend-dev/COMMANDS.md` on first-time onboarding).

## Repo layout

```
SKILL.md                               router + onboarding (3 lanes, 2 orientations)
AGENTS.md                              entry point for non-Claude-Code harnesses
LICENSE                                MIT
references/
  curriculum.md                        T0–T11 topic tree + sources
  builder-first.md                     10-loop spec
  exercise-bank.md                     catalog by tier
  incidents.md                         real backend postmortems
  practical-mode.md                    exercise playbook (multi-language)
  theory-modes.md                      5 teaching modes
  session-control.md                   pause/resume/checkpoint
  spaced-repetition.md                 SR queue + progress.json schema
  anti-patterns-with-examples.md       paired bad/good for each anti-pattern
assets/
  workspace-README.md                  copied to learner's ~/backend-dev/README.md
  COMMANDS.md                          slash commands reference
  progress-template.json               initial progress.json
  builder-first/
    setup/README.md                    toolchain setup per language
    go/loop-1-bare-crud/ … loop-10-loadtest/   runnable Go scaffolds
    _spec-only/loop-1-bare-crud/       language-agnostic mirror
```

## Workspace layout (after first run)

```
~/backend-dev/
├── README.md             ← guide for the learner
├── COMMANDS.md           ← slash commands reference
├── progress.json         ← long-term progress
├── session-state.md      ← resume pointer (overwritten each session)
├── notes/                ← topic notes the tutor generates
│   └── diagrams/
├── cheatsheets/          ← per-topic quick reference cards
├── projects/             ← builder-first loops + standalone exercises
├── reviews/              ← mock interview / design review writeups
├── flashcards/           ← optional — exported decks per topic
└── meta/                 ← anything else
```

Everything in `~/backend-dev/` is yours. Nothing leaves your machine.

## Status

**v0.1.0** — usable but pre-1.0.

Shipping:
- All three lanes + both orientations + 12-tier curriculum
- 10 builder-first loops in Go (runnable, with WIN/BREAK criteria per loop)
- Loop 1 spec-only mirror for non-Go learners
- Spaced repetition (SM-2 lite), session state, multi-harness portability

Pending (see `backlog.md` for the full list):
- Python (FastAPI) builder-first scaffolding
- Loops 2–10 spec-only mirrors
- Workspace viewer, pinned dependencies, multi-branch distribution
- `tests/` infrastructure (skill activation, schema validation, mode routing)

## License

MIT. See [LICENSE](LICENSE).
