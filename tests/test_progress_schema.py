"""Validate `assets/progress-template.json` and `tests/fixtures/progress_valid.json`
against `tests/schemas/progress.schema.json`.

Catches the bug where the template drifts from the documented schema in
`references/spaced-repetition.md`, or where a session-state-writing change
omits a required field.
"""

from __future__ import annotations

import json
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
SCHEMA_PATH = REPO_ROOT / "tests" / "schemas" / "progress.schema.json"
TEMPLATE_PATH = REPO_ROOT / "assets" / "progress-template.json"
FIXTURE_PATH = REPO_ROOT / "tests" / "fixtures" / "progress_valid.json"


def main() -> int:
    try:
        from jsonschema import Draft7Validator
    except ImportError:
        print("FAIL: jsonschema not installed. Run: pip install jsonschema")
        return 1

    failures: list[str] = []

    if not SCHEMA_PATH.is_file():
        print(f"FAIL: schema missing at {SCHEMA_PATH}")
        return 1
    schema = json.loads(SCHEMA_PATH.read_text())
    validator = Draft7Validator(schema)

    targets = [
        ("template", TEMPLATE_PATH),
        ("fixture", FIXTURE_PATH),
    ]

    for label, path in targets:
        if not path.is_file():
            failures.append(f"{label} missing at {path}")
            continue
        try:
            data = json.loads(path.read_text())
        except json.JSONDecodeError as e:
            failures.append(f"{label} at {path} is not valid JSON: {e}")
            continue
        errors = sorted(validator.iter_errors(data), key=lambda e: list(e.absolute_path))
        if errors:
            for err in errors:
                loc = "/".join(str(p) for p in err.absolute_path) or "<root>"
                failures.append(f"{label} schema violation at {loc}: {err.message}")

    if failures:
        for f in failures:
            print(f"FAIL: {f}")
        print(f"\ntest_progress_schema: {len(failures)} failure(s)")
        return 1
    print(f"PASS: test_progress_schema (template + fixture validate clean)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
