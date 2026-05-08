# tests/

Deterministic structural tests for the `backend-tutor` skill source.
**Not** behavior tests — those need an LLM-as-judge harness (deferred).

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

## What's deliberately not tested here

- **Activation routing** — does the skill description trigger on "start the backend course" but not on "help me debug a Python script"? Needs an LLM-as-judge harness; deferred to a separate session.
- **Mode routing** — does "review my schema" route to design review, "design Twitter at scale" route to system-design-tutor handoff, "teach me X" go to theory mode? Same harness dependency.
- **CI integration** — defer until tests are stable locally and the layout settles.

## Layout

```
tests/
  README.md
  run_all.py
  test_practical_coverage.py
  test_reference_presence.py
  test_progress_schema.py
  test_frontmatter.py
  schemas/
    progress.schema.json
  fixtures/
    progress_valid.json
```

## Adding a test

Drop a new `test_<name>.py` script next to the others. `run_all.py` picks it up
automatically by glob. Follow the pattern: print human-readable PASS/FAIL/WARN
lines, return 0 on success, 1 on failure.
