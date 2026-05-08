"""LOOP_VERSIONS.md sanity check.

- Frontmatter parses; has `last_verified` (ISO date) and `staleness_threshold_days` (int).
- Every `assets/builder-first/go/loop-*/go.mod` is referenced somewhere in the body
  (catches the case where someone adds a new loop and forgets to document its pins).
- `tools/check-staleness.py` runs against the file without crashing.

This test does NOT enforce that the file is fresh — that would make CI fail
on the calendar. Use the script at runtime if you want that signal.
"""

from __future__ import annotations

import datetime as _dt
import re
import subprocess
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
VERSIONS_FILE = REPO_ROOT / "LOOP_VERSIONS.md"
GO_LOOPS_DIR = REPO_ROOT / "assets" / "builder-first" / "go"
STALENESS_SCRIPT = REPO_ROOT / "tools" / "check-staleness.py"


def parse_frontmatter(text: str) -> dict[str, str]:
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n", text, flags=re.DOTALL)
    if not m:
        return {}
    out: dict[str, str] = {}
    for line in m.group(1).splitlines():
        line = line.strip()
        if not line or line.startswith("#") or ":" not in line:
            continue
        k, v = line.split(":", 1)
        out[k.strip()] = v.strip()
    return out


def main() -> int:
    failures: list[str] = []

    if not VERSIONS_FILE.is_file():
        print(f"FAIL: {VERSIONS_FILE.relative_to(REPO_ROOT)} not found")
        return 1

    text = VERSIONS_FILE.read_text(encoding="utf-8")
    fm = parse_frontmatter(text)

    raw_date = fm.get("last_verified")
    if not raw_date:
        failures.append("frontmatter missing last_verified")
    else:
        try:
            _dt.date.fromisoformat(raw_date)
        except ValueError:
            failures.append(f"last_verified is not ISO date: {raw_date!r}")

    threshold = fm.get("staleness_threshold_days")
    if threshold is None:
        failures.append("frontmatter missing staleness_threshold_days")
    else:
        try:
            n = int(threshold)
            if n <= 0:
                failures.append(f"staleness_threshold_days must be positive, got {n}")
        except ValueError:
            failures.append(f"staleness_threshold_days is not an int: {threshold!r}")

    # Every shipped Go loop dir should be name-checked in the body.
    if GO_LOOPS_DIR.is_dir():
        for loop_dir in sorted(GO_LOOPS_DIR.iterdir()):
            if not loop_dir.is_dir():
                continue
            if not (loop_dir / "go.mod").is_file():
                continue
            slug = loop_dir.name
            if slug not in text:
                failures.append(f"LOOP_VERSIONS.md doesn't mention shipped loop: {slug}")

    # Run the staleness script — should not crash on parse, regardless of fresh/stale.
    if STALENESS_SCRIPT.is_file():
        proc = subprocess.run(
            [sys.executable, str(STALENESS_SCRIPT)],
            capture_output=True,
            text=True,
        )
        if proc.returncode == 2:
            failures.append(f"check-staleness.py errored on parse (exit 2): {proc.stderr.strip()}")
    else:
        failures.append(f"missing: {STALENESS_SCRIPT.relative_to(REPO_ROOT)}")

    if failures:
        for f in failures:
            print(f"FAIL: {f}")
        print(f"\ntest_loop_versions: {len(failures)} failure(s)")
        return 1
    print(
        f"PASS: test_loop_versions (last_verified={fm.get('last_verified')}, "
        f"threshold={fm.get('staleness_threshold_days')}d)"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
