# Audit — Global Tutoring-Skill Checklist vs `ending_back`

**Date:** 2026-05-08 12:10
**Source:** "Tutoring Skill: Global All-Inclusive Checklist" provided by user, consolidated from `system-design-tutor` + `ai-system-tutor` + this skill (`ending_back`).
**Method:** Section-by-section walk against `SKILL.md`, `references/*`, `assets/*`, repo root. ✅ pass / ⚠️ partial / ❌ missing.

---

## 1. Skill identity & metadata — ⚠️

| Item | Status | Note |
|---|---|---|
| YAML frontmatter (name, description, license, compatibility, metadata) | ⚠️ | Has only `name` + `description` (`SKILL.md:2-3`). Missing `license`, `compatibility`, `metadata` (author, category, domain, version). |
| Description starts with what + trigger phrases | ✅ | `SKILL.md:3` includes "Trigger phrases: ..." |
| Negative scope + sibling hand-offs | ✅ | "Do NOT use for unrelated coding tasks. For pure architecture-at-scale design ... hand off to system-design-tutor. For LLM-specific infra ... ai-systems-tutor." |
| Teaching-philosophy one-liner | ✅ | `SKILL.md:8` "you drive ... user steers when they want a detour" |

## 2. Session controller — ✅

All 5 items pass.

- Workspace location convention `~/backend-dev/`, CWD-then-home (`SKILL.md:30`).
- Three-case branch (`SKILL.md:34-38`).
- Override map table (`SKILL.md:42-61`).
- "Proposal as default, not demand" stated (`SKILL.md:42`).
- Slash commands `/plan /start /quiz /continue /notes /config /loop ...` listed.

## 3. First-time onboarding — ⚠️

| Item | Status | Note |
|---|---|---|
| Tutor drives the flow | ✅ | "**You drive the entire flow.** Don't ask the user what they want." (`SKILL.md:67`) |
| Workspace setup with subdirs | ⚠️ | Spec drift: uses `projects/` instead of checklist's `exercises/`. Functionally equivalent (builder-first hands-on code goes here) but a strict reader of the cross-skill checklist would flag. |
| README.md, progress.json, session-state.md, COMMANDS.md copied | ✅ | `SKILL.md:75-80`. |
| **Goal capture before diagnostic** | ⚠️ | F3 just landed: optional probe (`SKILL.md:154-158`). Checklist says MUST. The probe is intentionally optional to not gate casual learners — a strict reader of the global checklist would flag this as soft. |
| Track / orientation capture | ✅ | Step 2.5. |
| Lane routing via vibe check (5 questions, patterns not counts) | ✅ | Step 2. |
| Adjacent-domain mode | ✅ | `working_mode = adjacent_domain`. |
| Language preference | ✅ | `learner.language` in `progress.json`. |
| Diagnostic shape adapts to lane | ✅ | 3a foundations / 3b working 12-Q / 3c senior 6-Q. |

## 4. Diagnostic design — ✅

All 12 items pass.

- Run without permission, one question at a time, no answer reveals.
- "Strengths and gaps must be equally specific — both must cite the actual answer the learner gave." (`SKILL.md:231`)
- Inline glosses on likely-unfamiliar terms (e.g., `*Idempotency* — running an operation twice has the same effect as once`).
- Adaptive depth (strong+specific → deeper; hand-wavy → base + note mechanism gap; total miss → move on, down-weight adjacent).
- Gap classification: vocabulary / mechanism / engineering-rationale (`SKILL.md:243-247`).
- "Lead with strengths every time, even when the learner missed most of the diagnostic." (`SKILL.md:248`)
- "Particular gap: ... pick one, not three." (`SKILL.md:241`)
- "**Avoid the word 'intermediate' and any level-comparison framing.**" (`SKILL.md:231`)
- Save findings to `notes/diagnostic-YYYY-MM-DD.md` (`SKILL.md:289`).

## 5. Curriculum & path planning — ⚠️

| Item | Status | Note |
|---|---|---|
| Topic tree with prereqs in `references/curriculum.md` | ✅ | File present. |
| Mapping topic → anchor sources | ✅ | Per CLAUDE.md "T0–T11 topic tree + sources". |
| **Path-suggestion-by-goal table** | ⚠️ | Not seen in `SKILL.md`; may be in `curriculum.md` (not opened during audit). Verify. |
| Two valid traversals (foundations-first vs builder-first) | ✅ | `decisions.md` 2026-05-07 entry. |
| Builder-first WIN/BREAK criteria | ✅ | `WIN.md` + `BREAK.md` per loop. |
| `/loop skip` + `/loop quickpass` | ✅ | `SKILL.md:60-61`. |
| "Builder-first not a license to skip foundations" | ✅ | `SKILL.md:159`. |

## 6. Mode dispatch — ✅

- Mode dispatch table (`SKILL.md:342-353`).
- "Load files only when the relevant mode is active. Never preload everything." (`SKILL.md:355`)
- Theory / practical / SR / mock interview / design review / notes / session control all present.
- Hand-off rules to sibling skills (`SKILL.md:357-360`).

## 7. Teaching philosophy — ⚠️ (1 execution gap)

All 13 items written. **Execution gap:**

- "Ground every lesson in real incidents" — rule is in `SKILL.md:384`. Yesterday's persona round (2026-05-07/08): `references/incidents.md` was opened by **0 of 8 tutor agents** across paired-context tests. The rule exists; the file isn't reliably consulted. Either inline canonical incidents into `SKILL.md` or add an explicit Step 3 hook that forces the load.

All other items pass: source anchoring, cycle-don't-camp, calibration before teaching, calibration suspension for foundations lesson 1, comprehension checks every 2-3 moves, honest critic, push for numbers, senior numbers invert, read register not just words, surface stated context, honor the explicit ask.

## 8. State persistence — ⚠️

| Item | Status | Note |
|---|---|---|
| `progress.json` schema (profile, topics, sr_queue, exercises, loops, sessions) | ✅ | `assets/progress-template.json`. |
| `session-state.md` | ✅ | Schema in `references/session-control.md`. |
| Update on lesson finish, pause, 30+min, before /compact | ✅ | `SKILL.md:432-436`. |
| Update progress.json after every meaningful interaction | ✅ | `SKILL.md:438-442`. |
| **Append-only event log** | ⚠️ | Schema uses `.entries[]` arrays implying append, but "append-only, never delete" is **not** explicitly stated. ~1 line fix in `progress-template.json` `_notes` or `references/spaced-repetition.md`. |
| **Exercise tuning telemetry** (`planned_difficulty`, `observed_difficulty`, `hints_used_max_level`) | ⚠️ | `progress-template.json` has `exercises.entries[]` but no per-entry difficulty schema. Not specified anywhere visible. |
| User-facing language differs from internal ("review queue" vs `sr_queue`) | ✅ | `SKILL.md:293`. |
| Track editable mid-course | ⚠️ | `/config` command exists; "switch tracks at any session start" not explicitly promised. |

## 9. Spaced repetition — ✅

- SR scheduler in `references/spaced-repetition.md` ("SR queue + progress.json schema" per CLAUDE.md).
- Items added on miss (`SKILL.md:291`).
- `/quiz` command (`SKILL.md:54`).
- Warm resume does overdue SR first (`SKILL.md:325`).

## 10. Resume protocol — ⚠️

| Item | Status | Note |
|---|---|---|
| Propose, don't ask | ✅ | `SKILL.md:323`. |
| Priority order (mid-lesson <14d → SR overdue → next step) | ✅ | `SKILL.md:323-326`. |
| Format string (≤4 lines) | ✅ | `SKILL.md:328`. |
| ≥14-day reminder + slash-command nudge | ✅ | `SKILL.md:334`. |
| **Mid-lesson + practical → "make easier / harder" reminder** | ❌ | Not present. Also tied to §11 below — the difficulty knob doesn't exist anywhere yet. |

## 11. Difficulty adjustment — ❌

**Whole section missing.**

- No "make this easier" / "make this harder" knob in `SKILL.md` or `references/practical-mode.md`.
- No semantics defined (easier = downshift scope/constraints, same topic; harder = add one realistic failure or scale constraint, same topic).
- No first-practical promise ("If this feels off-level, say 'make this easier' or 'make this harder'").
- Confirmed via grep: "easier", "harder", "difficulty" appear only in non-§11 contexts in `SKILL.md`; not at all in `practical-mode.md`.

One of the biggest single gaps. ~30 min copy work to add to `SKILL.md` + `practical-mode.md` + override map.

## 12. Mock interview mode — ✅

All 7 items present in `SKILL.md` Mock Interview Mode section: don't drive, requirements first, BoE numbers, trade-off probe, failures injected, scoring buckets, write to `reviews/YYYY-MM-DD-<system>.md`.

## 13. Design review mode — ✅

All 5 items, including the full 9-item stress-test list verbatim (10x scale, dependent service outage, hot key, thundering herd, slow downstream, primary DB failure, cache cluster loss, deploy-mid-incident, secrets compromise).

## 14. Notes generation mode — ✅

All 8 items: on-demand + end-of-topic offer, one file per topic with update-not-overwrite, strict structure (one-liner / core / trade-offs table / numbers / anchors / mistakes / artifacts), 2-min skimmable, self-contained, no transcript dump, honest about gaps, don't break flow, show + save + tell where.

## 15. Context-window management — ✅

- 60+ msg trigger, end-of-debugging, mode-switch triggers (`SKILL.md:447-453`).
- "Always write state to disk first, then suggest the command." (`SKILL.md:454`)
- Tool-agnostic phrasing: "Claude Code: `/compact`; Codex: new task; Copilot CLI: new session; Claude.ai: summary-then-new-chat."

## 16. Circuit breakers — ✅

- Tutor-side circuit breaker on senior misroute (`SKILL.md:258`): "That's on the router — let me drop to the standard diagnostic ..."
- Lane-recovery in first 1-2 lesson messages (`SKILL.md:303-308`).
- Re-diagnostic affordance in senior-lane closing (`SKILL.md:277-279`).

## 17. Anti-patterns — ⚠️

| Item | Status | Note |
|---|---|---|
| All 12 listed MUST anti-patterns | ✅ | `SKILL.md:558-571`. |
| **Answer-all-N-questions anti-pattern** | ❌ | Yesterday's persona round (Wei turn 5) flagged it. Not in `SKILL.md`. |
| `references/anti-patterns-with-examples.md` (DIFFERENTIATOR) | ❌ | File not present. Sibling `system-design-tutor` is said to have one. |

## 18. Format & tone — ✅

All 5 items: ~250-word ceiling, no-emoji-unless-user-first, real-systems anchors, diagram tier (Mermaid in chat / interactive HTML in workspace / ASCII fallback), code in `learner.language`.

## 19. Repo structure & infrastructure — ❌

| Item | Status | Note |
|---|---|---|
| `SKILL.md` | ✅ | |
| `references/` (curriculum, theory-modes, practical-mode, exercise-bank, incidents, spaced-repetition, session-control, builder-first) | ✅ | All 8 expected files present. |
| `assets/workspace-README.md`, `assets/progress-template.json`, `assets/COMMANDS.md` | ✅ | |
| `assets/exercise-templates/` | ⚠️ | Directory exists, **empty**. Already in backlog. |
| **`tests/`** | ❌ | Directory not present. Whole §21 follows from this. |
| **`LICENSE`** | ❌ | Not present. Checklist says MIT or equivalent. |
| `references/builder-first.md` (DIFFERENTIATOR) | ✅ | |
| **`references/anti-patterns-with-examples.md`** (DIFFERENTIATOR) | ❌ | |

## 20. Portability & distribution — ❌

| Item | Status | Note |
|---|---|---|
| Tool-agnostic protocol prose | ✅ | `SKILL.md:14`: "Translate to your harness's tool primitives." |
| State-as-files (no MCP, no DB) | ✅ | `SKILL.md:14`: "State lives entirely as files in the workspace." |
| **`AGENTS.md` for non-Claude-Code harnesses** | ❌ | Not present. |
| **Per-platform install guide (macOS/Linux symlink + Windows directory junction)** | ❌ | CLAUDE.md shows only macOS/Linux `ln -s`; no Windows. |
| **Multi-branch distribution** (main, cc-windows, codex-macos, codex-windows) | ❌ | Single `main` branch. |
| **Workspace viewer** (index.html + manifest.json + `python -m http.server`) | ❌ | Not in `workspace-README.md`, not in `assets/`. |
| **Pinned dependencies + staleness warning** (`uv.lock`, `LOOP_VERSIONS.md`, ">6mo old" banner) | ❌ | Go `go.mod` files exist but no manifest/staleness mechanism. `loop-2-persist/go.mod` pgx pin is stale (already in backlog). |

## 21. Testing & validation — ❌

**Whole section missing.** No `tests/` directory; nothing in CI.

| Item | Status |
|---|---|
| Skill activation test cases | ❌ |
| Mode routing test cases | ❌ |
| `progress.json` schema validator (Python script, runnable in CI) | ❌ |
| Reference fixtures (filled-in `progress_valid.json`) | ❌ |
| Practical coverage tests (every tier represented, required tags) | ❌ |
| Automated structural test suite (`run_all.py`) | ❌ |

Single biggest gap. CI on schema + activation + routing would have caught the SKILL.md "Python first-class scaffolding shipped" lie that yesterday's persona round caught manually.

## 22. Hand-off & scope discipline — ✅

- Sibling-skill awareness: explicit hand-off to `system-design-tutor` and `ai-systems-tutor`.
- Cross-link, don't duplicate: `SKILL.md:20`, CLAUDE.md "Don't duplicate content from siblings; cross-link."
- Out-of-scope honest: `SKILL.md:360` "Frontend / mobile / pure ML training → out of scope; say so honestly."

---

## Roll-up — biggest misses ranked by leverage

1. **§21 testing & validation** — entire section missing; single largest gap. CI on schema + activation + routing.
2. **§11 difficulty adjustment** — entire section missing from SKILL.md / practical-mode.md. ~30 min.
3. **§19/20 distribution infrastructure** — no LICENSE, no AGENTS.md, no Windows guide, no workspace viewer, no pinned-deps + staleness. Each is small individually; together they're real work.
4. **§7 incidents.md execution gap** — file exists, never gets opened by the tutor in practice (yesterday's persona evidence). Either inline canonical incidents into SKILL.md or add a forced-load hook.
5. **§17 anti-pattern misses** — answer-all-N-questions not in SKILL.md; `anti-patterns-with-examples.md` doesn't exist.
6. **Small spec drifts** — §1 frontmatter (license/compatibility/metadata), §3 goal capture is optional vs MUST, §8 append-only + exercise tuning telemetry not specified, §10 missing easier/harder reminder, §3 `projects/` vs `exercises/` directory naming.

## Suggested order of attack

**Tier A — copy edits, ~2 hours total:**
1. §11 difficulty adjustment (SKILL.md + practical-mode.md + override map).
2. §17 add answer-all-N anti-pattern to `SKILL.md`.
3. §10 add easier/harder reminder to warm-resume practical-exercise resume.
4. §1 expand YAML frontmatter (license, compatibility, metadata).
5. §8 add "append-only" + "track switchable mid-course" lines.
6. §7 add an explicit "When citing an incident, load `references/incidents.md` first" hook in Step 3.

**Tier B — files to create, ~1 day total:**
7. `LICENSE` — MIT or equivalent.
8. `AGENTS.md` — minimal, points non-Claude-Code harnesses at the SKILL.md flow.
9. `references/anti-patterns-with-examples.md` — paired correct/incorrect for each anti-pattern.
10. Windows install instructions in CLAUDE.md or workspace-README.

**Tier C — real work, days each:**
11. §21 entire test suite (`tests/` dir + activation tests + schema validator + run_all.py + CI hookup).
12. Workspace viewer (index.html + manifest.json + http.server instructions).
13. Pinned-deps + staleness banner + LOOP_VERSIONS.md.
14. Multi-branch distribution (main + cc-windows + codex-macos + codex-windows).

Tier A is a clean afternoon; Tier B fills another. Tier C warrants its own sessions.
