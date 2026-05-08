"""SKILL.md YAML frontmatter sanity checks.

Catches the bug where required frontmatter fields go missing during edits, or
where the version field stops being a valid semver, or where harness/platform
lists collapse to empty.

Note: uses a line-based parser rather than `yaml.safe_load`. Strict YAML rejects
unquoted descriptions that contain colon-space sequences (the SKILL.md
description has `Trigger phrases: "start the backend course"` mid-string),
even though every agent-harness frontmatter parser in the wild accepts them.
We mirror the harness's leniency rather than imposing strict YAML on the
description; structured sub-blocks (compatibility, metadata) get strict
parsing because they're well-behaved.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
SKILL_MD = REPO_ROOT / "SKILL.md"

REQUIRED_TOP_LEVEL = ["name", "description", "license", "compatibility", "metadata"]
REQUIRED_COMPAT = ["harnesses", "platforms"]
REQUIRED_METADATA = ["author", "category", "domain", "version"]
SEMVER_RE = re.compile(r"^\d+\.\d+\.\d+(?:[-+].+)?$")


def extract_frontmatter(text: str) -> str | None:
    if not text.startswith("---\n"):
        return None
    end = text.find("\n---", 4)
    if end == -1:
        return None
    return text[4:end]


def parse_frontmatter(fm_text: str) -> dict:
    """Parse top-level keys with a tolerant line-based approach.

    A top-level key is a line matching `^(\\w+):\\s*(.*)$` with no indentation.
    Indented lines (`^  ...`) belong to the most recent top-level key as
    nested fields. Inline list values (`[a, b, c]`) are split on commas.
    """
    out: dict = {}
    current_key: str | None = None
    nested_block: list[str] = []

    def flush_nested(key: str, lines: list[str]) -> None:
        if not lines:
            return
        nested: dict = {}
        for line in lines:
            m = re.match(r"^\s+(\w+):\s*(.*)$", line)
            if not m:
                continue
            k, v = m.group(1), m.group(2).strip()
            if v.startswith("[") and v.endswith("]"):
                items = [i.strip() for i in v[1:-1].split(",") if i.strip()]
                nested[k] = items
            else:
                nested[k] = v
        out[key] = nested

    for line in fm_text.splitlines():
        if not line.strip():
            continue
        if line.startswith(" ") or line.startswith("\t"):
            if current_key is not None:
                nested_block.append(line)
            continue
        if current_key is not None and nested_block:
            flush_nested(current_key, nested_block)
            nested_block = []
        m = re.match(r"^(\w+):\s*(.*)$", line)
        if not m:
            continue
        key, value = m.group(1), m.group(2)
        if value == "":
            current_key = key
            out[key] = None
        else:
            out[key] = value
            current_key = key
            nested_block = []
    if current_key is not None and nested_block:
        flush_nested(current_key, nested_block)
    return out


def main() -> int:
    if not SKILL_MD.is_file():
        print(f"FAIL: SKILL.md not found at {SKILL_MD}")
        return 1

    text = SKILL_MD.read_text(encoding="utf-8")
    fm_text = extract_frontmatter(text)
    if fm_text is None:
        print("FAIL: SKILL.md has no YAML frontmatter delimited by `---` markers.")
        return 1

    fm = parse_frontmatter(fm_text)
    failures: list[str] = []

    for key in REQUIRED_TOP_LEVEL:
        if key not in fm:
            failures.append(f"missing top-level frontmatter field: {key}")

    if isinstance(fm.get("name"), str):
        if fm["name"] != "backend-tutor":
            failures.append(f"frontmatter name is {fm['name']!r}, expected 'backend-tutor'")
    elif "name" in fm:
        failures.append("frontmatter name is missing or not a string")

    desc = fm.get("description")
    if not isinstance(desc, str) or len(desc) < 100:
        failures.append("frontmatter description is missing or shorter than 100 chars")

    compat = fm.get("compatibility")
    if isinstance(compat, dict):
        for key in REQUIRED_COMPAT:
            if key not in compat:
                failures.append(f"missing compatibility.{key}")
            elif not isinstance(compat[key], list) or not compat[key]:
                failures.append(f"compatibility.{key} must be a non-empty list")
    elif compat is not None and "compatibility" in fm:
        failures.append("compatibility must be a mapping (got scalar)")

    md = fm.get("metadata")
    if isinstance(md, dict):
        for key in REQUIRED_METADATA:
            if key not in md:
                failures.append(f"missing metadata.{key}")
        version = md.get("version")
        if isinstance(version, str):
            if not SEMVER_RE.match(version):
                failures.append(f"metadata.version {version!r} is not semver (X.Y.Z)")
        elif version is not None:
            failures.append("metadata.version must be a string")
    elif md is not None and "metadata" in fm:
        failures.append("metadata must be a mapping (got scalar)")

    if failures:
        for f in failures:
            print(f"FAIL: {f}")
        print(f"\ntest_frontmatter: {len(failures)} failure(s)")
        return 1
    print("PASS: test_frontmatter")
    return 0


if __name__ == "__main__":
    sys.exit(main())
