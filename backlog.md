# Backlog

Durable "someday/maybe" items — distinct from session-level Unresolved (which is "next session"). Each entry: one-line item specific enough to act on cold, optional priority, **Why:** if non-obvious.

## What's next (as of 2026-05-09)

State: persona round 1 + Tier A/B/C audit work + spec drifts all closed. 4 branches live (`main` + 3 platform), 6/6 tests pass, decisions documented. The open items below are ranked by leverage — pick from the top.

### High leverage

1. **Persona round 2 — untested branches.** Same methodology as the §7 forced-load round (8 personas, paired pre/post analysis, write up findings under `test_findings/`). Targets: Case B cold resume (workspace exists, no `session-state.md` — 2-3 personas with varying gap), mid-lesson lane-recovery circuit breaker ("this is too basic / I'm drowning / you routed me wrong" within 1-2 lesson messages), multi-harness handoff (start in Claude Code, resume in Codex / Cursor / Copilot CLI via `~/backend-dev/session-state.md`), non-coder graceful exit. *Why high leverage:* the §7 round found 4 real content gaps in `incidents.md` plus validated the triple-belt fix. These four branches are completely untested in practice. Multi-hour effort; gets its own session.

2. **Symlink install activation test.** Open a fresh Claude Code session in any working directory and try a trigger phrase like *"start the backend course"*. Verify backend-tutor's onboarding actually fires (Step 1 → Q1-Q4 → routing). The symlink at `~/.claude/skills/backend-tutor/` was created 2026-05-08 but couldn't be self-validated mid-session because the skill list was fixed at startup. *Why high leverage:* if the description doesn't trigger, the whole skill is dead-on-arrival regardless of the rest. ~10 min if it works; could surface description-tuning bugs that block real adoption.

### Medium leverage

3. ~~**Design + implement the LLM-as-judge harness.**~~ Design landed 2026-05-09 (`prds/2026-05-09_llm-as-judge-harness.md`); implementation landed same day. `tests/run_llm.py` ships responder (Sonnet 4.6, low effort, no thinking, SKILL.md as cached system context) + judge (Opus 4.7, low effort, structured JSON via `output_config.format`). 5 bootstrap activation fixtures (`tests/fixtures/llm/activation_smoke.jsonl`) + 2 rubrics (`activation.md`, `mode-routing.md`). Skip paths verified: missing `ANTHROPIC_API_KEY` → SKIP exit 0; missing `anthropic` pkg → SKIP exit 0. Per-fixture errors don't crash the runner. Not wired into `run_all.py`. **What's still untested:** the actual responder/judge round-trip (needs a real key). What's not yet built: more bootstrap fixtures (mode-routing, more activation cases), CI wire-up. Per the PRD, mode-routing fixtures should come from persona round 2 transcripts.

4. **Loop 4 restructure check.** The `loop-4-auth/` "merge into main package" UX is the simplest pattern but means the file can't be type-checked in isolation — `go build ./...` fails standalone with a confusing error. Header comment warns; restructuring to `internal/auth/` (package `auth`) is the principled fix *if* the UX keeps biting in practice. Conditional: only do this if a learner reports the confusion in real use.

### Low leverage / deprioritized

5. **Authoring backlog** (blocks broader rollout, not personal use): Python (FastAPI) Loop 1 scaffold; Loops 2-10 spec-only mirrors; `assets/exercise-templates/` directory population. Skip until a non-Go learner actually shows up.

### Stale state to refresh periodically

- **Pin manifest staleness.** `LOOP_VERSIONS.md` `last_verified` is 2026-05-08; `tools/check-staleness.py` will start warning at the 180-day threshold (≈ 2026-11-04). Refresh: bump pgx if newer, run `go mod tidy` in each loop, update the manifest, re-date.

---

## Tier 1 — small fixes from persona-test round (high-value, fast)

*All four landed 2026-05-08 — see SKILL.md Step 2 (Q4 + routing table) and Step 2.5 (language gloss, switch-later line, optional goals probe), Step 3a adjacent-domain variant.*

- ~~F1. Gloss "scaffold" inline in SKILL.md Step 2.5 language pitch.~~ Done.
- ~~F2. Add "any of these are fine; you can switch later via /config" line after the language pitch.~~ Done.
- ~~F3. Add optional stated-goals/timeline probe to SKILL.md Step 2.5.~~ Done. Saved to `learner.stated_goals` in `progress.json`.
- ~~F4. Broaden adjacent-domain examples in SKILL.md Step 2.~~ Done — broadened examples *and* added "production systems where things have to keep running" generalization in Q4 + Step 3a variant.

## Tier 2 — design changes worth discussing first

- ~~**F5. Add "what bit you in production?" probe to adjacent-domain mode only.**~~ Done 2026-05-08. Step 3a's adjacent-domain variant ends with an optional soft prompt — *"What's a backend-adjacent thing that bit you in production at work?"* — that, if answered, is saved to `learner.stated_context` and surfaces back in the matching tier lesson (stale-cache → T4, queue duplicate → T3.7, etc.). Soft, not required.
- ~~**F6. Senior-lane Step 3c should explicitly request a real shipped system.**~~ Done 2026-05-08. Step 3c now opens with a strong-invitation anchor probe before the six questions — *"Pick a system you've actually shipped. Name it and one failure mode you'd want to dig into."* If named, the system is saved to `learner.stated_context` and referenced in the questions where it fits. Graceful fallback if the senior declines: questions run as written.
- ~~**F7. Warn at language-pick time when spec-only path will be rough.**~~ Done 2026-05-08. Step 2.5 now surfaces a one-time heads-up when the learner picked a spec-only language AND `level == foundations`: spec-only means more from-scratch writing per loop; can switch via `/config` any time. Working/senior + spec-only is left alone (Tyler archetype handles it fine). Hedge-register-only-but-Working learners aren't covered by this gate; if F7v2 is needed for that case, surface as a separate item.

## Tier 3 — polish / edge cases

- ~~**Tone-softener for confident-shallow learners (Devansh archetype).**~~ Done 2026-05-08. Added a phrasing-softener clause under SKILL.md's "Honest critic, not cheerleader" rule: pivot "you have X but not Y" framings to "X is solid; next: Y." Substance unchanged, delivery only — same rule applies to any analogous identity-adjacent framing.
- ~~**Surface-stated-context rule should also say *use it***, not just *note it*.~~ Done 2026-05-08. The "Surface stated context" rule in SKILL.md is now titled "Surface stated context — and use it, not just note it" and explicitly distinguishes the playback half (cheap trust-buying) from the use-it half (re-ordering the curriculum walk, grounding examples in the named project, pacing for the stated demographic). Playback signals you heard them; use-it proves you're acting on the signal.
- ~~**Add fluency-vs-knowledge axis to `progress.json` topic schema.**~~ Done 2026-05-08. Optional `fluency` field per-topic with values `unknown | slow | building | fluent`. Documented in `references/spaced-repetition.md` (with worked example: T7.4 status=solid + fluency=slow), enforced in `tests/schemas/progress.schema.json` (enum), explained in `progress-template.json` `_notes`. Tutor schedules drilling rather than re-teaching when fluency=slow.

## Untested branches (run another persona round)

- **Case B cold resume.** Workspace exists, no `session-state.md`. None of 8 personas hit this. Should test: 2-3 personas resuming after a gap, varying days-since-last-session.
- **Mid-lesson lane-recovery circuit breaker.** ("This is too basic" / "I'm drowning" / "you routed me wrong" within the first 1-2 lesson messages.) Devansh's circuit-break was pre-diagnostic; mid-lesson untested.
- **Multi-harness handoff.** Started in Claude Code, resumed in Codex / Cursor / Copilot CLI. The bridge is `~/backend-dev/session-state.md`; untested in practice.
- **Non-coder branch.** Skill exits gracefully per spec. Untested experimentally; probably fine.

## Global-checklist gaps (from 2026-05-08 audit — `audits/2026-05-08_12-10-22_global-checklist.md`)

The 22-section "Tutoring Skill: Global All-Inclusive Checklist" run against `ending_back`. Items are tiered by effort. Plan: walk Tier A first, then B, then C across separate sessions.

### Tier A — copy edits (~2 hours total, can land in one session)

*All six landed 2026-05-08. SKILL.md got the difficulty knob in the override map + a new "Difficulty adjustment" subsection under Step 3 philosophy; the answer-all-N and incidents-from-memory anti-patterns; the easier/harder warm-resume reminder; expanded YAML frontmatter (license, compatibility, metadata); a switchable-mid-course note on `/config`; and a forced-load hook on the "Ground every lesson in real incidents" subsection. `references/practical-mode.md` got a full "Difficulty knob" section with first-practical promise and telemetry schema. `progress-template.json` `_notes` covers append-only invariant, mid-course profile switching, and the exercise telemetry fields.*

- ~~G-A1. Add §11 difficulty adjustment to SKILL.md + practical-mode.md.~~ Done.
- ~~G-A2. Add answer-all-N-questions anti-pattern to SKILL.md.~~ Done.
- ~~G-A3. Add "make easier / harder" reminder to warm-resume when resuming a practical exercise.~~ Done.
- ~~G-A4. Expand SKILL.md YAML frontmatter.~~ Done.
- ~~G-A5. Add "append-only event log" and "track switchable mid-course" lines.~~ Done.
- ~~G-A6. Add explicit "load `references/incidents.md` before citing an incident" hook in SKILL.md Step 3.~~ Done.

### Tier B — files to create (~1 day total)

*All four landed 2026-05-08. `LICENSE` (MIT, 2026 rogue-socket — matches the sibling), root-level `AGENTS.md` for non-Claude-Code harnesses (workspace location, three lanes, two orientations, slash commands, sibling-skill hand-offs, default language), `references/anti-patterns-with-examples.md` with paired bad/good for all 14 anti-patterns in `SKILL.md` (cross-linked from the SKILL.md anti-patterns block + reference-files list), and Windows directory-junction install instructions added to `CLAUDE.md` (PowerShell `mklink /J`, no admin / Developer Mode required).*

- ~~G-B1. Create `LICENSE` (MIT or equivalent).~~ Done.
- ~~G-B2. Create `AGENTS.md`.~~ Done.
- ~~G-B3. Create `references/anti-patterns-with-examples.md`.~~ Done.
- ~~G-B4. Add Windows install instructions to CLAUDE.md or workspace-README.md.~~ Done — went into `CLAUDE.md` since that's where the install one-liner already lived.

### Tier C — real work (days each, separate sessions)

- **G-C1. Build `tests/` infrastructure.** MVP shipped 2026-05-08: `tests/` now has 4 deterministic tests + `run_all.py` + JSON Schema for `progress.json` + valid fixture + README. Coverage: practical-coverage (catches claim/disk drift in builder-first scaffolds — the motivating bug), reference-presence (every `references/*.md` mentioned in SKILL.md exists), progress-schema (template + fixture validate against draft-7 schema), frontmatter (required fields, semver, non-empty harness/platform lists, lenient YAML to mirror harness parsers). **4/4 pass clean** after downgrading the SKILL.md:151 Python claim from "prefilled scaffolding shipped" to "planned (...; spec-only fallback applies until then)" — the test suite caught the drift on its first run, exactly as §21 motivated.
  - **Deferred** to a follow-up session: activation tests (does description trigger on "start the backend course" but not "help me debug Python"), mode-routing tests (does "review my schema" route to design review, "design X at 100M users" hand off to system-design-tutor), CI integration. These need an LLM-as-judge harness; separate design pass.
  - **Run:** `python3 tests/run_all.py`. Dep: `pip install jsonschema` (stdlib for everything else).
- ~~**G-C2. Build the workspace viewer.**~~ Done 2026-05-08. `assets/workspace-viewer/` now has `index.html` (sidebar grouped by category, marked.js from CDN, system-aware light/dark), `manifest.template.json` (schema example), and `regenerate-manifest.py` (stdlib-only, walks `notes/cheatsheets/flashcards/sessions/reviews/meta`, derives title from H1 with filename fallback). SKILL.md Step 1 grew a step 6 to copy the viewer into `~/backend-dev/viewer/` and run regenerate at workspace setup. workspace-README.md got a "Browse your workspace" section. New `tests/test_viewer_assets.py` (HTML parses, template valid, regenerate runs empty + populated) — 5/5 pass.
- ~~**G-C3. Add pinned-deps + staleness banner.**~~ Done 2026-05-08. New `LOOP_VERSIONS.md` (root) with YAML frontmatter (`last_verified`, `staleness_threshold_days: 180`) and per-loop dep tables. New `tools/check-staleness.py` (stdlib): exit 0/1/2 for fresh/stale/error. SKILL.md Step 1.7 runs the check at workspace setup; Warm/Cold Resume runs it too and surfaces a one-line banner only when stale. `assets/workspace-README.md` got a "Toolchain pin freshness" pointer. New `tests/test_loop_versions.py` (6/6 in suite). Subsumed pgx bump: `loop-2-persist/go.mod` now `pgx v5.7.5` with `go 1.24.0` directive (forced by transitive deps); `go vet` + `go build` clean. Setup README floor bumped from "1.22 or newer" → "1.24 or newer."
- ~~**G-C4. Multi-branch distribution.**~~ Done 2026-05-08. Four branches live on origin: `main` (universal), `cc-windows` (Windows + Claude Code), `codex-macos` (macOS/Linux + Codex CLI), `codex-windows` (Windows + Codex CLI). Each platform branch carries one commit on top of main: a custom `INSTALL.md` (single-path) plus a small README banner naming the branch. New `INSTALL.md` on main owns the full multi-platform/multi-harness matrix; main's README slimmed to a quick path + a Branches table linking the platform branches. New `tools/sync-platform-branches.sh` rebases the platform branches onto main and pushes with `--force-with-lease`; dry-run validated all three rebase clean. Conflicts only happen if main itself touches `INSTALL.md` or the README banner region.

### Misc spec drifts (small, may not be worth fixing)

- ~~**§3 workspace subdir naming**~~ Decided 2026-05-09. Kept `projects/`. See `decisions.md` 2026-05-09 *Workspace subdir naming*.
- ~~**§3 goal capture optional vs MUST**~~ Decided 2026-05-09. Kept optional (intentional divergence from audit). See `decisions.md` 2026-05-09 *Goal/timeline capture is optional*.
- ~~**§5 path-suggestion-by-goal table**~~ Done 2026-05-08. `references/curriculum.md` now has a "Path suggestions by stated goal" section just above the real-time path, covering interview prep (4-8 wk), first backend role, payments/billing service, real-time (cross-link), and SRE / on-call readiness. Each entry: tier walk + why-this-shape + pacing note.
- ~~**§8 exercise tuning telemetry**~~ Done (already shipped). Schema bump landed in the Tier A round: `tests/schemas/progress.schema.json` now has `planned_difficulty` (enum), `observed_difficulty` (enum), `hints_used_max_level` (1-6), and `adjustments[]` on each `exercises.entries[]` item. `assets/progress-template.json` `_notes` documents the four fields. The backlog item was stale-as-described.

## Authoring backlog (deprioritized for this user but blocks broader rollout)

- **Python (FastAPI) Loop 1 scaffold.** SKILL.md says Python is first-class but only Go ships. *Why deferred:* user is Go-only personally; Python is a "skill for other learners" task.
- **Loops 2-10 spec-only mirrors.** Currently only Loop 1 has `assets/builder-first/_spec-only/`. Non-Go/non-Python learners get one loop's worth of spec, then nothing. *Why deferred:* same as above — affects future learners on Node/Java/Rust/etc.
- **`assets/exercise-templates/` directory population.** SKILL.md mentions it for non-builder-first learners doing standalone tier exercises. Currently empty. Could fill in per-tier as foundations-first learners reach each tier.

## Verification / quality

- ~~`go vet` + `go build` on all 10 Go loop scaffolds.~~ Done 2026-05-08 with Go 1.26.3. Loops 1, 2 build standalone (Loop 2 tests skip cleanly without Postgres). Loops 3, 5q, 5w, 6, 8, 9 build standalone in synthesized modules. Loop 4 vet-clean but fails standalone build (no `main()` — expected for a merge-into-main delta; clarifying header comment added). `loop-2-persist/go.sum` now generated and `loop-2-persist/go.mod` carries indirect deps from `tidy`. `.gitignore` added in `assets/builder-first/go/` for the loop1/loop2 binaries.
- ~~**Bump pgx version pin.**~~ Done 2026-05-08 as part of G-C3. Bumped `v5.5.5` → `v5.7.5`; held back from latest (`v5.9.2` at time of bump) to keep a wider Go-version compat window.
- **Restructure Loop 4 into `internal/auth/` (package auth)** *if* the merge-into-main UX proves confusing in practice. *Why:* the current "drop into your main package" pattern is the simplest UX but means the file can't be type-checked in isolation; a learner running `go build ./...` in `loop-4-auth/` hits a confusing error. Header comment now warns; restructuring is the principled fix if it keeps biting.
- **Test the symlink install end-to-end.** Symlink created 2026-05-08 (`~/.claude/skills/backend-tutor` → `~/Documents/ending_back`). Activation test still pending — needs a fresh Claude Code session to verify the description triggers, onboarding flow runs, and references lazy-load. The current session's skill list was fixed at startup, so this can't be self-validated mid-conversation.
- ~~**Verify §7 forced-load hook actually triggers in practice.**~~ Done 2026-05-08. Post-fix paired-persona round: **8/8 loaded `references/incidents.md` before citing specifics** (vs. 0/8 pre-fix). 5/8 cited from the file with full specifics; 3/8 correctly identified file gaps and invoked the don't-fabricate clause rather than confabulate. Triple-belt validated; no need to escalate to inlining. Full writeup: `test_findings/2026-05-08_17-11-00_incidents-forced-load-verify.md`. Surfaced 4 content gaps in `incidents.md` itself — see new authoring item below.

## `incidents.md` content gaps (surfaced by 2026-05-08 §7 verification round)

*All four landed 2026-05-08. Each entry is framed as a pattern / positive-case-study (matching the existing Stripe T1 precedent) rather than a fabricated single-incident RCA, since public postmortems naming these specific failure modes are thin — the lessons live in canonical guidance docs and engineering-blog positive case studies.*

- ~~**T1 — idempotency-failure incident.**~~ Done. Pattern entry citing Stripe blog (positive anchor) + AWS SQS / GitHub / Stripe webhook docs; covers retry-without-dedup and key-without-body-hash variants.
- ~~**T2 — bad-index-in-production incident.**~~ Done. Pattern entry citing Strong Migrations + gh-ost + Postgres `CONCURRENTLY` docs + Use the Index Luke; covers DDL-lock and write-amplification variants. Tooling existence is the evidence the failure mode is widespread.
- ~~**T3 — queue-redelivery / non-idempotent-consumer postmortem.**~~ Done. Pattern entry citing AWS SQS / Lambda / Kafka / webhook docs; covers SQS visibility-timeout, Kafka rebalance, and webhook 5xx-retry variants. Universal mitigation: dedup table keyed on stable application-domain message ID.
- ~~**T3/T11 — saga / outbox case study.**~~ Done. Positive-case-study entry citing Uber Cadence/Temporal origin + microservices.io outbox reference + Debezium + Richardson's Microservices Patterns + Newman's Building Microservices 2e. Senior-grade build-vs-adopt frame.

## Remote / distribution state

- ~~**Local `main` ahead of `origin/main`.**~~ Pushed 2026-05-08. `origin` is `https://github.com/rogue-socket/backend-tutor.git`; remote repo confirmed live; main fast-forwarded clean from `7ed630c` → `80c1769`.
- **Seed Tier C as GitHub issues.** Tier C items (G-C2 workspace viewer, G-C3 pinned-deps + staleness banner, G-C4 multi-branch distribution) are each "days, separate sessions" — good fit for `/to-issues`. Remote now exists, so this is unblocked.
