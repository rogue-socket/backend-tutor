# PRD — LLM-as-judge test harness

**Status:** initial implementation shipped 2026-05-09 (`tests/run_llm.py`, 5 activation fixtures, 2 rubrics). Skip paths verified; live responder/judge round-trip not yet exercised. Mode-routing fixtures still TODO (waiting on persona-round-2 transcripts per the original "When to build" note).
**Owner:** rogue-socket
**Date:** 2026-05-09
**Closes backlog item:** #3 from "What's next" — design the LLM-as-judge harness for activation + mode-routing tests deferred from G-C1.

---

## Problem

`tests/run_all.py` ships 6 deterministic tests that catch structural regressions (frontmatter, JSON schema, reference presence, manifest pinning, viewer assets, practical-coverage). They can't catch the failure modes that depend on *what an LLM does* given the skill content:

1. **Activation regressions.** Does the SKILL.md description trigger backend-tutor on `"start the backend course"` and *not* on `"help me debug Python"`? A poorly-tuned description either misses real learners or steals queries from sibling skills (system-design-tutor, ai-systems-tutor).
2. **Mode-routing regressions.** Inside the skill, does the dispatcher pick the right mode? `"review my schema"` should hit design review; `"design Twitter for 100M users"` should hand off to system-design-tutor; `"teach me caching"` should load `references/theory-modes.md` + `references/incidents.md`.

Both are LLM-judged because the failure mode is "the model picked wrong," not "the file is malformed." The §7 forced-load round caught one such failure manually (incidents.md opened by 0/8 personas) — at one-shot scale, manual works; at every-PR scale, it doesn't.

## Goals

- Catch activation drift when SKILL.md frontmatter / description changes.
- Catch mode-routing regressions when SKILL.md mode-dispatch table or reference set changes.
- Run cheaply enough to be useful pre-PR (seconds-to-minutes; cents-to-dollars).
- Stay opt-in — never block the deterministic suite.

## Non-goals

- Full end-to-end onboarding simulation (Q1-Q4 → routing → first lesson). Out of scope; that's manual persona-round territory.
- Regression-testing every prose change in SKILL.md. Judges grade routing decisions, not writing quality.
- Running on every push in CI. Pre-PR + nightly + on-demand only.

---

## Architecture

Two models, two roles.

```
                       ┌─────────────────┐
fixtures/*.json  ───►  │  Responder LLM  │  ───► raw response
                       │  (sees SKILL.md │
                       │   as context)   │
                       └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
rubric/*.md      ───►  │   Judge LLM     │  ───► {pass: bool, reason: str}
                       │  (different     │
                       │   model)        │
                       └─────────────────┘
                                │
                                ▼
                       results/<run-id>.jsonl
```

### Responder

- Receives the user prompt as input + SKILL.md as system context.
- For activation tests: also receives a *list* of available skills (backend-tutor + 2-3 plausible siblings + a "no skill matches" option) and is asked to pick. This isolates the routing decision from full-onboarding noise.
- For mode-routing tests: assumes backend-tutor is already active, asks which mode/reference set the dispatcher should select.
- Low temperature (0.0-0.3) for repeatability.

### Judge

- A **different** model from the responder (mitigates self-grading bias). Default pairing: responder = Sonnet, judge = Opus. Cost ratio sane; judge stronger.
- Receives the responder's output + the rubric for this fixture + the expected outcome.
- Returns structured JSON: `{pass: bool, confidence: 1-5, reason: str}`. Schema-validated on parse.
- Temperature 0.0 (rubric-graded judging should be near-deterministic).

### Fixtures

```
tests/fixtures/llm/
  activation/
    should_trigger.jsonl       # prompts backend-tutor SHOULD claim
    should_not_trigger.jsonl   # prompts a sibling or no-skill should claim
  mode-routing/
    theory.jsonl               # → references/theory-modes.md + incidents.md
    design-review.jsonl        # → design review mode
    handoff-system-design.jsonl   # → system-design-tutor
    handoff-ai-systems.jsonl   # → ai-systems-tutor
    builder-first.jsonl        # → references/builder-first.md
```

Each line: `{prompt, expected, rubric_id, notes?}`. ~10-15 fixtures per file to start; grow as regressions surface.

### Rubrics

```
tests/rubrics/
  activation.md      # what counts as "triggered backend-tutor"
  mode-routing.md    # how to judge mode picks; cross-skill handoff criteria
```

Versioned. When SKILL.md description changes meaningfully, rubric may need update — surface this with a `rubric_version` field in fixtures + a check that fixture and rubric versions match.

---

## Wire-in

**Not** in `tests/run_all.py`. Separate runner:

```
tests/run_llm.py [--smoke | --full] [--category activation | mode-routing | all]
                 [--judge-model opus] [--responder-model sonnet]
                 [--seed N] [--out results/<run-id>.jsonl]
```

- `--smoke` runs ~5 representative fixtures per category (~30 sec, ~$0.10).
- `--full` runs everything (~3 min, ~$1-3).
- Skips cleanly with exit 0 + a `SKIP` log line if `ANTHROPIC_API_KEY` is unset (so contributors without keys aren't blocked).

Suggested cadence:
- **Pre-PR (manual):** `tests/run_llm.py --smoke` when SKILL.md or references/ change.
- **Nightly (cron/Action):** `--full` against `main`. Surface failures to an issue, don't gate merges.
- **Never:** on every push.

---

## Decisions

### Responder ≠ judge
Self-grading inflates pass rates. Always pair models. If only one provider is available, use the same model with explicit "you are evaluating someone else's work" framing as a fallback — and flag in results that pair-different was unavailable.

### Judge returns structured output
`{pass, confidence, reason}` is the contract. Free-form judge responses lead to ambiguous CI signal. Use the SDK's structured-output / tool-use mode to enforce.

### Fixtures are JSONL, not Python
JSONL is editable by non-coders, diffable in PRs, easy to grow. A fixture is data, not code.

### Rubric is markdown, not code
Same reasoning — humans iterate on rubrics. The judge prompt template loads `rubrics/<id>.md` verbatim.

### Cost-aware by default
Smoke is the default mode. Full requires explicit opt-in. Per-run cost is logged with results.

### Failure handling
A fixture failure is informational by default — does *not* fail CI. Failure log goes to a results file; humans triage. Hard-fail mode (`--strict`) for the rare case (release branch, before-tagging) where regressions should block.

---

## Cost / flakiness mitigations

- Temperature 0 on both responder + judge.
- Judge sees rubric, not free-form "is this good?" — graded judging is more stable than vibe judging.
- Multi-sample for borderline cases: if judge `confidence < 4`, re-run with seed+1 and require majority pass. Adds cost; only kicks in when needed.
- Fixture set kept small (≤ ~50 total). LLM-as-judge isn't a coverage tool; it's a regression alarm.
- Results JSONL is committed to a `results/` branch (or gh-pages) so trend-over-time is visible without polluting main history.

---

## Open questions (decide when implementing)

1. **Which Anthropic SDK pattern?** Plain `messages.create` vs the agent-SDK (`claude-agent-sdk`) vs the Files API for SKILL.md context. Probably plain messages with prompt caching on the SKILL.md system block, since that block is large and rarely changes.
2. **Prompt caching.** SKILL.md is ~600 lines and gets re-sent on every fixture. With prompt caching (`cache_control: ephemeral`), 50 fixtures × 1 call/each = 1 cache write + 49 hits = ~70-90% cost reduction. Worth wiring in from day one.
3. **Where does the judge's rubric live in its prompt?** As system context or as part of the user message? System makes it cacheable; user message lets us swap rubrics per-fixture without breaking the cache. Probably system + per-fixture user-message override slot.
4. **Reproducibility seed.** SDK supports `seed` param? If not, accept that LLM-as-judge has irreducible variance; quote variance in the results summary.
5. **CI integration.** GitHub Actions workflow with `ANTHROPIC_API_KEY` as a repo secret; nightly cron + on-demand `workflow_dispatch`. Or external (a Modal/Replit cron). Defer until the harness itself works.
6. **Bootstrap fixtures.** Where do the first 30-50 prompts come from? Two sources: (a) the persona-round transcripts (real prompts learners produced), (b) hand-authored "obvious-trigger / obvious-non-trigger" pairs. Start hand-authored; grow from real transcripts as personal use accumulates.

## Out of scope (explicitly)

- LLM-judging of *content quality* (does the lesson teach indexes well?). That's a different problem; needs different rubrics and probably human review.
- Cross-skill handoff *execution* tests (does system-design-tutor actually pick up?). Tests one skill in isolation; coordination tests are a separate harness.
- Anti-pattern detection in tutor responses ("did the tutor answer all 5 questions when it should have asked one at a time?"). High-value but harder to rubric. Phase 2.

## Effort estimate

- Initial implementation: 1-2 sessions (~4-8 hours).
  - `tests/run_llm.py` runner: ~150 LOC.
  - 30-50 bootstrap fixtures: ~1 hour to author.
  - 2 rubric files: ~30 min.
  - GitHub Actions wire-up: ~30 min.
- Steady-state cost: ~$0.10 per `--smoke` run, ~$1-3 per `--full` nightly. Manageable.

## When to build

After persona round 2 lands. The persona round will surface real prompt patterns that should become fixtures; building the harness first means hand-fabricating fixtures that may not match how learners actually phrase things.
