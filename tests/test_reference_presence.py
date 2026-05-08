"""Every `references/*.md` file mentioned in SKILL.md must exist on disk.

Catches the bug where SKILL.md references a file that hasn't been authored, or
where a reference file is renamed without updating the routing/forced-load hooks.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
SKILL_MD = REPO_ROOT / "SKILL.md"
REFS_DIR = REPO_ROOT / "references"


def main() -> int:
    if not SKILL_MD.is_file():
        print(f"FAIL: SKILL.md not found at {SKILL_MD}")
        return 1
    if not REFS_DIR.is_dir():
        print(f"FAIL: references/ not found at {REFS_DIR}")
        return 1

    text = SKILL_MD.read_text(encoding="utf-8")
    mentioned = sorted(set(re.findall(r"references/([a-zA-Z0-9_\-]+\.md)", text)))
    if not mentioned:
        print("FAIL: SKILL.md mentions zero reference files — parser regression?")
        return 1

    failures: list[str] = []
    on_disk = {p.name for p in REFS_DIR.glob("*.md")}

    for name in mentioned:
        if name not in on_disk:
            failures.append(f"SKILL.md mentions references/{name} but file is missing.")

    # Also flag orphan files in references/ that SKILL.md never points to —
    # those are either dead, or someone added a file without wiring it in.
    orphans = sorted(on_disk - set(mentioned))
    if orphans:
        print(f"WARN: orphan reference files (present on disk, not mentioned in SKILL.md): {orphans}")

    if failures:
        for f in failures:
            print(f"FAIL: {f}")
        print(f"\ntest_reference_presence: {len(failures)} failure(s)")
        return 1
    print(f"PASS: test_reference_presence ({len(mentioned)} reference(s) mentioned, all present)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
