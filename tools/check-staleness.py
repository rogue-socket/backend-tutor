#!/usr/bin/env python3
"""Print whether LOOP_VERSIONS.md's pins are stale.

Reads `last_verified` and `staleness_threshold_days` from the YAML frontmatter
of `LOOP_VERSIONS.md` (sibling to this script's parent dir). Exits:
  0 — fresh
  1 — stale (older than threshold)
  2 — file missing or malformed

The tutor's SKILL.md flow may shell out to this at workspace setup or warm
resume. The learner can also run it directly.

Stdlib only.
"""
from __future__ import annotations

import datetime as _dt
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
VERSIONS_FILE = REPO_ROOT / "LOOP_VERSIONS.md"


def parse_frontmatter(text: str) -> dict[str, str]:
    """Lenient YAML-like frontmatter parser. Only handles the keys we use."""
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n", text, flags=re.DOTALL)
    if not m:
        return {}
    out: dict[str, str] = {}
    for line in m.group(1).splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        if ":" not in line:
            continue
        k, v = line.split(":", 1)
        out[k.strip()] = v.strip()
    return out


def main(argv: list[str]) -> int:
    if not VERSIONS_FILE.is_file():
        print(f"check-staleness: {VERSIONS_FILE} not found", file=sys.stderr)
        return 2

    fm = parse_frontmatter(VERSIONS_FILE.read_text(encoding="utf-8"))
    raw_date = fm.get("last_verified")
    if not raw_date:
        print("check-staleness: LOOP_VERSIONS.md frontmatter missing last_verified", file=sys.stderr)
        return 2

    try:
        last = _dt.date.fromisoformat(raw_date)
    except ValueError:
        print(f"check-staleness: last_verified is not ISO date: {raw_date!r}", file=sys.stderr)
        return 2

    try:
        threshold = int(fm.get("staleness_threshold_days", "180"))
    except ValueError:
        threshold = 180

    today = _dt.date.today()
    age = (today - last).days

    if age <= threshold:
        print(f"OK: LOOP_VERSIONS.md last verified {raw_date} ({age} days ago, threshold {threshold})")
        return 0

    print(
        f"STALE: LOOP_VERSIONS.md last verified {raw_date} ({age} days ago, "
        f"threshold {threshold}). Pins may have moved — see the 'How to refresh' "
        f"section of LOOP_VERSIONS.md."
    )
    return 1


if __name__ == "__main__":
    sys.exit(main(sys.argv))
