# Backlog

Durable "someday/maybe" items — distinct from session-level Unresolved (which is "next session"). Each entry: one-line item specific enough to act on cold, optional priority, **Why:** if non-obvious.

## Tier 1 — small fixes from persona-test round (high-value, fast)

*All four landed 2026-05-08 — see SKILL.md Step 2 (Q4 + routing table) and Step 2.5 (language gloss, switch-later line, optional goals probe), Step 3a adjacent-domain variant.*

- ~~F1. Gloss "scaffold" inline in SKILL.md Step 2.5 language pitch.~~ Done.
- ~~F2. Add "any of these are fine; you can switch later via /config" line after the language pitch.~~ Done.
- ~~F3. Add optional stated-goals/timeline probe to SKILL.md Step 2.5.~~ Done. Saved to `learner.stated_goals` in `progress.json`.
- ~~F4. Broaden adjacent-domain examples in SKILL.md Step 2.~~ Done — broadened examples *and* added "production systems where things have to keep running" generalization in Q4 + Step 3a variant.

## Tier 2 — design changes worth discussing first

- **F5. Add "what bit you in production?" probe to adjacent-domain mode only.** *Why:* When asked of Anika, surfaced a stale-cache story that was the perfect adjacent-domain anchor for Loop 6. Wei's session would have benefited from the same probe. Should be a standard step under `working_mode = adjacent_domain`, not asked of standard-mode learners.
- **F6. Senior-lane Step 3c should explicitly request a real shipped system.** *Why:* Marcus volunteered "fintech" and "infra"; the diagnostic anchored to those. A reserved senior wouldn't volunteer. Make it explicit: *"Pick a system you've shipped. Describe one failure mode you'd most want to revisit."* Then anchor diagnostic Q's to *that* system rather than abstract scenarios.
- **F7. Warn at language-pick time when spec-only path will be rough.** *Why:* Tyler (Node spec-only, Working lane, confident) is fine. Joseph (Java spec-only, Foundations lane, ESL hedger) is going to need much more coaching support. Trigger: if `language` is spec-only AND (`level == foundations` OR diagnostic showed heavy hedge-register), add a one-line at language-pick time: "spec-only means more tutor coaching per loop; switch to a prefilled language any time."

## Tier 3 — polish / edge cases

- **Tone-softener for confident-shallow learners (Devansh archetype).** *Why:* "Honest critic, not cheerleader" is the right rule, but "you have shipping reflexes but no instrumentation muscle" walks close to the cliff for cocky learners. Substance unchanged; consider phrasing-softener for delivery only.
- **Surface-stated-context rule should also say *use it***, not just *note it*. *Why:* Joseph's "most of my classmates are 21" was noted to session-state but not reflected in lesson-pacing decisions. The rule's current text says to play context back at assessment; should also influence pacing / emphasis throughout.
- **Add fluency-vs-knowledge axis to `progress.json` topic schema.** *Why:* Joseph's T7 gap was speed-under-pressure, not absence-of-knowledge. The current topic status (`unknown / weak / in_progress / solid / mastered`) doesn't capture this. Optional `fluency` field per-topic would let the tutor schedule drilling differently from teaching.

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

- **G-C1. Build `tests/` infrastructure.** *Why:* §21 entire section missing — single largest gap. Must include: activation test cases (does the description trigger on the right phrases?), mode-routing tests (does "design X" go to mock interview, "review my Y" to design review?), `progress.json` schema validator (Python, CI-runnable), reference fixture (`progress_valid.json`), practical-coverage tests, `run_all.py`. CI on these would have caught yesterday's "Python first-class scaffolding shipped" lie.
- **G-C2. Build the workspace viewer.** *Why:* §20 DIFFERENTIATOR — `index.html` + `manifest.json` + `python -m http.server` instruction in workspace-README so learners can browse `notes/` / `cheatsheets/` / `flashcards/` as a styled site.
- **G-C3. Add pinned-deps + staleness banner.** *Why:* §20 DIFFERENTIATOR. Equivalent of `uv.lock` + `LOOP_VERSIONS.md`. Banner in workspace-README if manifest is >6 months old. Also subsumes the existing pgx-version-bump backlog item.
- **G-C4. Multi-branch distribution.** *Why:* §20 DIFFERENTIATOR — `main`, `cc-windows`, `codex-macos`, `codex-windows`. Shared curriculum + per-host shell conventions. Heavy lift; defer until Tier A/B are done.

### Misc spec drifts (small, may not be worth fixing)

- **§3 workspace subdir naming**: skill uses `projects/`; checklist uses `exercises/`. Functionally equivalent. *Why drifted:* `projects/` reads as "the fun-stuff folder" for builder-first; `exercises/` is more neutral for foundations-first. Worth deciding rather than just renaming silently.
- **§3 goal capture optional vs MUST**: F3 (just landed) makes the goals/timeline probe optional. Checklist says MUST. The optionality is intentional ("don't gate casual learners"); flag the divergence rather than fix.
- ~~**§5 path-suggestion-by-goal table**~~ Done 2026-05-08. `references/curriculum.md` now has a "Path suggestions by stated goal" section just above the real-time path, covering interview prep (4-8 wk), first backend role, payments/billing service, real-time (cross-link), and SRE / on-call readiness. Each entry: tier walk + why-this-shape + pacing note.
- **§8 exercise tuning telemetry**: `progress-template.json` has `exercises.entries[]` but no per-entry difficulty schema (`planned_difficulty`, `observed_difficulty`, `hints_used_max_level`). Adding requires a schema bump.

## Authoring backlog (deprioritized for this user but blocks broader rollout)

- **Python (FastAPI) Loop 1 scaffold.** SKILL.md says Python is first-class but only Go ships. *Why deferred:* user is Go-only personally; Python is a "skill for other learners" task.
- **Loops 2-10 spec-only mirrors.** Currently only Loop 1 has `assets/builder-first/_spec-only/`. Non-Go/non-Python learners get one loop's worth of spec, then nothing. *Why deferred:* same as above — affects future learners on Node/Java/Rust/etc.
- **`assets/exercise-templates/` directory population.** SKILL.md mentions it for non-builder-first learners doing standalone tier exercises. Currently empty. Could fill in per-tier as foundations-first learners reach each tier.

## Verification / quality

- ~~`go vet` + `go build` on all 10 Go loop scaffolds.~~ Done 2026-05-08 with Go 1.26.3. Loops 1, 2 build standalone (Loop 2 tests skip cleanly without Postgres). Loops 3, 5q, 5w, 6, 8, 9 build standalone in synthesized modules. Loop 4 vet-clean but fails standalone build (no `main()` — expected for a merge-into-main delta; clarifying header comment added). `loop-2-persist/go.sum` now generated and `loop-2-persist/go.mod` carries indirect deps from `tidy`. `.gitignore` added in `assets/builder-first/go/` for the loop1/loop2 binaries.
- **Bump pgx version pin.** `loop-2-persist/go.mod` pins `v5.5.5` (mid-2024); current is likely v5.7+. Not blocking but stale.
- **Restructure Loop 4 into `internal/auth/` (package auth)** *if* the merge-into-main UX proves confusing in practice. *Why:* the current "drop into your main package" pattern is the simplest UX but means the file can't be type-checked in isolation; a learner running `go build ./...` in `loop-4-auth/` hits a confusing error. Header comment now warns; restructuring is the principled fix if it keeps biting.
- **Test the symlink install end-to-end.** `ln -s ~/Documents/ending_back ~/.claude/skills/backend-tutor` → invoke skill → verify it actually onboards a learner. Untested.
- ~~**Verify §7 forced-load hook actually triggers in practice.**~~ Done 2026-05-08. Post-fix paired-persona round: **8/8 loaded `references/incidents.md` before citing specifics** (vs. 0/8 pre-fix). 5/8 cited from the file with full specifics; 3/8 correctly identified file gaps and invoked the don't-fabricate clause rather than confabulate. Triple-belt validated; no need to escalate to inlining. Full writeup: `test_findings/2026-05-08_17-11-00_incidents-forced-load-verify.md`. Surfaced 4 content gaps in `incidents.md` itself — see new authoring item below.

## `incidents.md` content gaps (surfaced by 2026-05-08 §7 verification round)

Three of 8 personas in the §7 verify round hit topic-tier intersections where `references/incidents.md` had no named incident, forcing the tutor to honest-decline + redirect to public sources. Adding these would lift the from-file citation rate from 5/8 toward 8/8 on a re-run.

- **T1 — idempotency-failure incident.** A payments double-charge or webhook redelivery RCA with date + amount + named system. Stripe's blog is the positive case study; need a *negative* one.
- **T2 — bad-index-in-production incident.** Write amplification under a too-many-indexes table, or a non-CONCURRENTLY `CREATE INDEX` that took an `ACCESS EXCLUSIVE` lock and stalled checkout. GitHub or Shopify blogs likely have one.
- **T3 — queue-redelivery / non-idempotent-consumer postmortem.** SQS visibility-timeout expiry → duplicate processing, or webhook retry → duplicate side-effect. Public RCAs are thin; conference talks exist.
- **T3/T11 — saga / outbox case study.** Uber Cadence/Temporal origin posts and eBay's outbox-pattern writeups are the canonical sources; pick one, summarize the failure mode it solved.

## Remote / distribution state

- **Local `main` is 10 commits ahead of `origin/main`** (as of 2026-05-08 13:03). Not pushed. *Why holding:* user hasn't said push.
- **Remote repo may not exist yet.** README and AGENTS.md reference `https://github.com/rogue-socket/backend-tutor`; this URL is aspirational. `gh repo create rogue-socket/backend-tutor --public --source=. --remote=origin` is the next move if so. Confirm before pushing.
- **Seed Tier C as GitHub issues.** `audits/2026-05-08_*-global-checklist.md` Tier C items (G-C1 tests/, G-C2 workspace viewer, G-C3 pinned-deps + staleness banner, G-C4 multi-branch distribution) are each "days, separate sessions" — good fit for `/to-issues`. Do this after the remote is created.
