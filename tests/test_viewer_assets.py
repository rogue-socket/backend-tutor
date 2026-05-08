"""Workspace viewer assets sanity check.

- `assets/workspace-viewer/index.html` parses as HTML.
- `assets/workspace-viewer/manifest.template.json` is valid JSON with the
  required top-level keys (`generated_at`, `entries`).
- `assets/workspace-viewer/regenerate-manifest.py` runs against an empty
  fake workspace and produces a valid manifest with `entries == []`.
"""

from __future__ import annotations

import json
import shutil
import subprocess
import sys
import tempfile
from html.parser import HTMLParser
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
VIEWER_DIR = REPO_ROOT / "assets" / "workspace-viewer"


class _StrictParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self.error_msg: str | None = None

    def error(self, message: str) -> None:  # py<3.10 compat
        self.error_msg = message


def check_html() -> list[str]:
    failures: list[str] = []
    path = VIEWER_DIR / "index.html"
    if not path.is_file():
        return [f"missing: {path.relative_to(REPO_ROOT)}"]
    text = path.read_text(encoding="utf-8")
    parser = _StrictParser()
    try:
        parser.feed(text)
        parser.close()
    except Exception as e:
        failures.append(f"index.html failed to parse: {e}")
    if "manifest.json" not in text:
        failures.append("index.html doesn't reference manifest.json")
    if "marked" not in text:
        failures.append("index.html doesn't import a markdown renderer")
    return failures


def check_template() -> list[str]:
    failures: list[str] = []
    path = VIEWER_DIR / "manifest.template.json"
    if not path.is_file():
        return [f"missing: {path.relative_to(REPO_ROOT)}"]
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as e:
        return [f"manifest.template.json invalid JSON: {e}"]
    for key in ("generated_at", "entries"):
        if key not in data:
            failures.append(f"manifest.template.json missing key: {key}")
    if not isinstance(data.get("entries"), list):
        failures.append("manifest.template.json: entries must be a list")
    return failures


def check_regenerate_script() -> list[str]:
    failures: list[str] = []
    script = VIEWER_DIR / "regenerate-manifest.py"
    if not script.is_file():
        return [f"missing: {script.relative_to(REPO_ROOT)}"]

    with tempfile.TemporaryDirectory() as tmp:
        workspace = Path(tmp)
        viewer = workspace / "viewer"
        viewer.mkdir()
        shutil.copy(script, viewer / "regenerate-manifest.py")

        proc = subprocess.run(
            [sys.executable, str(viewer / "regenerate-manifest.py")],
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            failures.append(
                f"regenerate-manifest.py exited {proc.returncode} on empty workspace: "
                f"{proc.stderr.strip() or proc.stdout.strip()}"
            )
            return failures

        manifest_path = viewer / "manifest.json"
        if not manifest_path.is_file():
            failures.append("regenerate-manifest.py didn't produce manifest.json")
            return failures
        try:
            data = json.loads(manifest_path.read_text(encoding="utf-8"))
        except json.JSONDecodeError as e:
            failures.append(f"generated manifest.json invalid: {e}")
            return failures
        if data.get("entries") != []:
            failures.append(f"empty workspace should yield empty entries, got: {data.get('entries')}")

        # Now drop a markdown file in and re-run; expect 1 entry.
        (workspace / "notes").mkdir()
        (workspace / "notes" / "indexes.md").write_text("# Indexes\n\nbody\n", encoding="utf-8")
        proc = subprocess.run(
            [sys.executable, str(viewer / "regenerate-manifest.py")],
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            failures.append(f"regenerate-manifest.py failed on populated workspace: {proc.stderr}")
            return failures
        data = json.loads(manifest_path.read_text(encoding="utf-8"))
        entries = data.get("entries", [])
        if len(entries) != 1:
            failures.append(f"expected 1 entry after adding one .md file, got {len(entries)}")
        elif entries[0].get("title") != "Indexes" or entries[0].get("path") != "notes/indexes.md":
            failures.append(f"unexpected entry shape: {entries[0]}")

    return failures


def main() -> int:
    if not VIEWER_DIR.is_dir():
        print(f"FAIL: {VIEWER_DIR.relative_to(REPO_ROOT)} not found")
        return 1

    failures: list[str] = []
    failures += check_html()
    failures += check_template()
    failures += check_regenerate_script()

    if failures:
        for f in failures:
            print(f"FAIL: {f}")
        print(f"\ntest_viewer_assets: {len(failures)} failure(s)")
        return 1
    print("PASS: test_viewer_assets (html parses, template valid, regenerate script works empty + populated)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
