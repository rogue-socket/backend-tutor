"""Verify practical-coverage claims in SKILL.md match on-disk reality.

The bug this catches: SKILL.md claiming a language ships "prefilled scaffolding"
when the corresponding `assets/builder-first/<lang>/loop-N-*` directories don't
exist (or are empty). The original §21 audit motivation was "would have caught
yesterday's Python first-class scaffolding shipped lie."

Approach:
1. Parse the language-claim block in SKILL.md (the `> - **<lang>** — <claim>` lines).
2. For each language claimed "prefilled scaffolding shipped", assert all 10
   loop directories exist on disk and contain at least one source file.
3. For spec-only languages, check the shared `_spec-only/` mirror has Loop 1
   at minimum (the floor; full mirror is a known authoring backlog item, so
   missing 2-10 is a warning, not a fail).
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
SKILL_MD = REPO_ROOT / "SKILL.md"
BUILDER_FIRST = REPO_ROOT / "assets" / "builder-first"
EXPECTED_LOOPS = list(range(1, 11))


def parse_language_claims(text: str) -> dict[str, str]:
    """Return {language_key: status} where status is 'prefilled' or 'spec-only'.

    Looks for blockquote lines like:
        > - **Go** — prefilled scaffolding shipped ...
        > - **Python (FastAPI)** — prefilled scaffolding shipped
        > - **Node / TypeScript** — supported, but you'll implement against a spec ...
        > - **Java / Kotlin** — supported, spec-only
    """
    claims: dict[str, str] = {}
    pattern = re.compile(r"^>\s*-\s*\*\*(.+?)\*\*\s*[—-]\s*(.+)$", re.M)
    for m in pattern.finditer(text):
        langs_raw = m.group(1).strip()
        body = m.group(2).strip().lower()
        if "prefilled scaffolding shipped" in body:
            status = "prefilled"
        elif "spec-only" in body or "implement against a spec" in body:
            status = "spec-only"
        else:
            continue
        # Split on '/' (handles "Node / TypeScript", "Java / Kotlin"); drop
        # parenthetical qualifiers like "(FastAPI)".
        for part in langs_raw.split("/"):
            name = part.split("(")[0].strip().lower()
            if not name or name == "other":
                continue
            # Pick a stable filesystem-friendly key.
            key = {"typescript": "node"}.get(name, name)
            # First claim wins; prefilled trumps spec-only if both seen.
            existing = claims.get(key)
            if existing is None or (existing == "spec-only" and status == "prefilled"):
                claims[key] = status
    return claims


def loop_dir_for(lang: str, loop_n: int) -> Path | None:
    """Find `loop-<N>-*` directory under `assets/builder-first/<lang>/`."""
    parent = BUILDER_FIRST / lang
    if not parent.is_dir():
        return None
    matches = sorted(parent.glob(f"loop-{loop_n}-*"))
    return matches[0] if matches else None


def dir_has_source_file(d: Path) -> bool:
    """True if `d` contains at least one regular file (any extension)."""
    if not d.is_dir():
        return False
    for child in d.rglob("*"):
        if child.is_file():
            return True
    return False


def main() -> int:
    if not SKILL_MD.is_file():
        print(f"FAIL: SKILL.md not found at {SKILL_MD}")
        return 1

    text = SKILL_MD.read_text(encoding="utf-8")
    claims = parse_language_claims(text)
    if not claims:
        print("FAIL: parsed zero language claims from SKILL.md — parser regression?")
        return 1

    print(f"Parsed {len(claims)} language claims from SKILL.md: {claims}")

    failures: list[str] = []
    warnings: list[str] = []

    # Hard checks: prefilled-language coverage.
    for lang, status in sorted(claims.items()):
        if status != "prefilled":
            continue
        lang_dir = BUILDER_FIRST / lang
        if not lang_dir.is_dir():
            failures.append(
                f"SKILL.md claims '{lang}' ships prefilled scaffolding, but "
                f"{lang_dir.relative_to(REPO_ROOT)} does not exist."
            )
            continue
        for n in EXPECTED_LOOPS:
            d = loop_dir_for(lang, n)
            if d is None:
                failures.append(
                    f"'{lang}' is claimed prefilled but loop-{n}-* directory is missing under "
                    f"{lang_dir.relative_to(REPO_ROOT)}/"
                )
            elif not dir_has_source_file(d):
                failures.append(
                    f"'{lang}' loop-{n} directory exists at {d.relative_to(REPO_ROOT)} but is empty."
                )

    # Soft checks: spec-only mirror has at least Loop 1.
    spec_only_dir = BUILDER_FIRST / "_spec-only"
    if any(s == "spec-only" for s in claims.values()):
        if not spec_only_dir.is_dir():
            failures.append(
                f"SKILL.md claims spec-only languages are supported, but "
                f"{spec_only_dir.relative_to(REPO_ROOT)} does not exist."
            )
        else:
            if loop_dir_for("_spec-only", 1) is None:
                failures.append(
                    f"spec-only mirror exists but Loop 1 is missing — that's the floor."
                )
            missing_high = [n for n in EXPECTED_LOOPS[1:] if loop_dir_for("_spec-only", n) is None]
            if missing_high:
                warnings.append(
                    f"spec-only mirror missing loops {missing_high} — known authoring backlog item, not a hard fail."
                )

    for w in warnings:
        print(f"WARN: {w}")
    if failures:
        for f in failures:
            print(f"FAIL: {f}")
        print(f"\ntest_practical_coverage: {len(failures)} failure(s), {len(warnings)} warning(s)")
        return 1
    print(f"PASS: test_practical_coverage ({len(warnings)} warning(s))")
    return 0


if __name__ == "__main__":
    sys.exit(main())
