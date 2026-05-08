# tests/

Two harnesses:
- **Deterministic suite** (`run_all.py`) — structural tests, stdlib + `jsonschema` only.
- **LLM-as-judge suite** (`run_llm.py`) — opt-in behavior tests; needs `ANTHROPIC_API_KEY` and the `anthropic` package. Skips cleanly if either is missing.

## Run

```bash
python3 tests/run_all.py
```

Or invoke a single test:

```bash
python3 tests/test_practical_coverage.py
```

Each test is a standalone Python 3 script that prints `PASS`/`FAIL`/`WARN` lines
and exits 0 on pass, 1 on fail.

## Dependencies

- `jsonschema` — used by `test_progress_schema.py`

```bash
pip install jsonschema
```

The other tests use only the Python 3 stdlib. `test_frontmatter.py` deliberately avoids `pyyaml`: strict YAML rejects unquoted descriptions containing `: ` sequences, but every agent-harness frontmatter parser in the wild accepts them — so the test mirrors harness leniency rather than imposing strict spec.

If a dependency is missing, the relevant test fails fast with a clear install hint.

## What each test asserts

| Test | What it catches |
|---|---|
| `test_practical_coverage.py` | Languages SKILL.md claims ship "prefilled scaffolding" actually have all 10 loop directories on disk with non-empty contents. **The motivating bug:** SKILL.md claiming Python is shipped when no Python scaffolds exist. |
| `test_reference_presence.py` | Every `references/*.md` file mentioned in SKILL.md exists on disk. Warns about orphan files (on disk but not mentioned). |
| `test_progress_schema.py` | `assets/progress-template.json` and `tests/fixtures/progress_valid.json` both validate against `tests/schemas/progress.schema.json` (Draft 7). |
| `test_frontmatter.py` | SKILL.md YAML frontmatter has required fields (`name`, `description`, `license`, `compatibility`, `metadata`); description is >=100 chars; harnesses and platforms are non-empty lists; `metadata.version` is semver. |

## LLM-as-judge suite (`run_llm.py`)

Opt-in behavior tests for activation and mode-routing — the parts the deterministic suite can't reach.

```bash
# Smoke (~5 fixtures, ~$0.10):
ANTHROPIC_API_KEY=... python3 tests/run_llm.py --smoke

# Full sweep:
python3 tests/run_llm.py --full

# Filter:
python3 tests/run_llm.py --category activation

# Release-gate mode (exit 1 on any failure; default is informational):
python3 tests/run_llm.py --strict
```

Architecture (`prds/2026-05-09_llm-as-judge-harness.md`): responder = Sonnet 4.6 (with SKILL.md as cached system context), judge = Opus 4.7 (with structured JSON output via `output_config.format`). Different models on purpose — self-grading inflates pass rates.

Skips cleanly with exit 0 if `ANTHROPIC_API_KEY` is unset or the `anthropic` package isn't installed. **Not wired into `run_all.py`** — the deterministic suite stays stdlib-only and CI-cheap.

Adding fixtures: drop new lines into `tests/fixtures/llm/*.jsonl`. Smoke runs use `*_smoke.jsonl` only; full runs match all `*.jsonl`. Each line: `{id, category, expected, prompt, rubric_id}`. Rubrics live in `tests/rubrics/<rubric_id>.md`.

## CI integration

Deferred. Run locally for now.

## Layout

```
tests/
  README.md
  run_all.py                       deterministic runner
  run_llm.py                       LLM-as-judge runner (opt-in)
  test_*.py                        deterministic tests
  schemas/
    progress.schema.json
  fixtures/
    progress_valid.json
    llm/
      activation_smoke.jsonl       smoke fixtures (run by --smoke / default)
  rubrics/
    activation.md
    mode-routing.md
```

## Adding a test

Drop a new `test_<name>.py` script next to the others. `run_all.py` picks it up
automatically by glob. Follow the pattern: print human-readable PASS/FAIL/WARN
lines, return 0 on success, 1 on failure.
