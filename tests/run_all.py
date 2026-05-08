"""Run every test_*.py in this directory; exit non-zero on any failure.

Usage:
    python3 tests/run_all.py

Each test is a standalone script that prints PASS/FAIL/WARN lines to stdout
and exits 0 (pass) or 1 (fail). This runner shells out so any single test
crashing won't take the others down.
"""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path

TESTS_DIR = Path(__file__).resolve().parent


def main() -> int:
    test_files = sorted(p for p in TESTS_DIR.glob("test_*.py") if p.is_file())
    if not test_files:
        print("FAIL: no test files matching test_*.py found in tests/")
        return 1

    results: list[tuple[str, int]] = []
    print(f"Running {len(test_files)} test(s)\n")
    for tf in test_files:
        print(f"--- {tf.name} ---")
        proc = subprocess.run(
            [sys.executable, str(tf)],
            cwd=TESTS_DIR.parent,
            capture_output=True,
            text=True,
        )
        sys.stdout.write(proc.stdout)
        if proc.stderr:
            sys.stderr.write(proc.stderr)
        results.append((tf.name, proc.returncode))
        print()

    passed = [name for name, rc in results if rc == 0]
    failed = [name for name, rc in results if rc != 0]
    print("=" * 60)
    print(f"Summary: {len(passed)} passed, {len(failed)} failed")
    for name in failed:
        print(f"  FAILED: {name}")
    return 0 if not failed else 1


if __name__ == "__main__":
    sys.exit(main())
