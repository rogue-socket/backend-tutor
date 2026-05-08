"""LLM-as-judge harness for backend-tutor.

Two-model architecture:
  - Responder (Sonnet 4.6): given SKILL.md + skill list + prompt, picks a skill / mode.
  - Judge (Opus 4.7): given responder output + rubric + expected, returns
    {pass, confidence, reason} via structured output.

Skips cleanly with exit 0 if either ANTHROPIC_API_KEY or the `anthropic` SDK is
unavailable — keeps this opt-in so contributors without keys aren't blocked.

Usage:
    python3 tests/run_llm.py --smoke              # ~5 fixtures, ~$0.10
    python3 tests/run_llm.py --full               # all fixtures
    python3 tests/run_llm.py --category activation
    python3 tests/run_llm.py --strict             # exit 1 on any failure (default: informational)

Design ref: prds/2026-05-09_llm-as-judge-harness.md
"""
from __future__ import annotations

import argparse
import json
import os
import sys
from dataclasses import dataclass
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
SKILL_MD = REPO_ROOT / "SKILL.md"
FIXTURES_DIR = REPO_ROOT / "tests" / "fixtures" / "llm"
RUBRICS_DIR = REPO_ROOT / "tests" / "rubrics"

RESPONDER_MODEL = "claude-sonnet-4-6"
JUDGE_MODEL = "claude-opus-4-7"

SIBLING_SKILLS = """\
Available skills (pick the one that should handle the prompt):

1. **backend-tutor** — agent-driven course on backend engineering implementation. Onboards a learner, runs lessons, ships hands-on exercises in a chosen language (Go, Python, Node, Java, Kotlin, Rust). Triggers on requests to learn/practice/review backend code at the implementation level.
2. **system-design-tutor** — architecture-at-scale design ("design X for N million users", capacity estimation, sharding strategy at the diagram level).
3. **ai-systems-tutor** — LLM-specific infrastructure (agent loops, RAG, prompt caching, multi-agent orchestration).
4. **no skill matches** — pick this if none of the above is a good fit (e.g., debugging a specific Python syntax error, frontend question, unrelated topic).
"""

JUDGE_SCHEMA = {
    "type": "object",
    "properties": {
        "pass": {"type": "boolean", "description": "True iff the responder picked correctly per the rubric."},
        "confidence": {"type": "integer", "minimum": 1, "maximum": 5, "description": "1=ambiguous/malformed, 5=clear-cut."},
        "reason": {"type": "string", "description": "One short sentence: what the responder picked and whether it matches expected."},
    },
    "required": ["pass", "confidence", "reason"],
    "additionalProperties": False,
}


@dataclass
class Fixture:
    id: str
    category: str
    expected: str
    prompt: str
    rubric_id: str


@dataclass
class Result:
    fixture_id: str
    expected: str
    responder_output: str
    judge_pass: bool
    judge_confidence: int
    judge_reason: str
    cache_read_tokens: int


def load_fixtures(category_filter: str | None, smoke: bool) -> list[Fixture]:
    if not FIXTURES_DIR.is_dir():
        return []
    pattern = "*_smoke.jsonl" if smoke else "*.jsonl"
    fixtures: list[Fixture] = []
    for path in sorted(FIXTURES_DIR.glob(pattern)):
        with path.open("r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                d = json.loads(line)
                if category_filter and d.get("category") != category_filter:
                    continue
                fixtures.append(Fixture(
                    id=d["id"],
                    category=d["category"],
                    expected=d["expected"],
                    prompt=d["prompt"],
                    rubric_id=d["rubric_id"],
                ))
    return fixtures


def load_rubric(rubric_id: str) -> str:
    path = RUBRICS_DIR / f"{rubric_id}.md"
    return path.read_text(encoding="utf-8")


def call_responder(client, skill_md: str, prompt: str) -> tuple[str, int]:
    """Returns (responder_output_text, cache_read_input_tokens)."""
    response = client.messages.create(
        model=RESPONDER_MODEL,
        max_tokens=400,
        thinking={"type": "disabled"},
        output_config={"effort": "low"},
        system=[
            {
                "type": "text",
                "text": skill_md,
                "cache_control": {"type": "ephemeral"},
            },
        ],
        messages=[
            {
                "role": "user",
                "content": (
                    f"{SIBLING_SKILLS}\n\n"
                    f"User prompt: {prompt!r}\n\n"
                    "Pick exactly one of the four options above and explain in one sentence why. "
                    "Format your response as:\n"
                    "PICK: <skill name or 'no skill matches'>\n"
                    "REASON: <one sentence>"
                ),
            },
        ],
    )
    text = next((b.text for b in response.content if b.type == "text"), "")
    cache_read = getattr(response.usage, "cache_read_input_tokens", 0) or 0
    return text, cache_read


def call_judge(client, rubric: str, fixture: Fixture, responder_output: str) -> dict:
    """Returns {pass: bool, confidence: int, reason: str}."""
    response = client.messages.create(
        model=JUDGE_MODEL,
        max_tokens=400,
        thinking={"type": "disabled"},
        output_config={
            "effort": "low",
            "format": {"type": "json_schema", "schema": JUDGE_SCHEMA},
        },
        system="You are a strict, rubric-driven grader. Apply the rubric exactly. Do not add criteria the rubric does not state.",
        messages=[
            {
                "role": "user",
                "content": (
                    f"# Rubric\n\n{rubric}\n\n"
                    f"# Fixture\n\n"
                    f"- id: {fixture.id}\n"
                    f"- expected: {fixture.expected}\n"
                    f"- prompt: {fixture.prompt!r}\n\n"
                    f"# Responder output\n\n{responder_output}\n\n"
                    "Apply the rubric and return your verdict."
                ),
            },
        ],
    )
    text = next((b.text for b in response.content if b.type == "text"), "")
    return json.loads(text)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--smoke", action="store_true", help="Run *_smoke.jsonl fixtures only (default).")
    parser.add_argument("--full", action="store_true", help="Run all fixtures.")
    parser.add_argument("--category", help="Filter fixtures by category (activation | mode-routing).")
    parser.add_argument("--strict", action="store_true", help="Exit 1 if any fixture fails (default: informational).")
    args = parser.parse_args()

    if not os.environ.get("ANTHROPIC_API_KEY"):
        print("SKIP: ANTHROPIC_API_KEY not set. LLM tests are opt-in.")
        return 0

    try:
        import anthropic
    except ImportError:
        print("SKIP: anthropic package not installed. Run `pip install anthropic` in a conda env.")
        return 0

    smoke = not args.full  # default to smoke
    fixtures = load_fixtures(args.category, smoke=smoke)
    if not fixtures:
        print(f"FAIL: no fixtures matched (category={args.category}, smoke={smoke}).")
        return 1

    if not SKILL_MD.is_file():
        print(f"FAIL: SKILL.md not found at {SKILL_MD}")
        return 1
    skill_md = SKILL_MD.read_text(encoding="utf-8")

    rubric_cache: dict[str, str] = {}
    client = anthropic.Anthropic()

    print(f"Running {len(fixtures)} fixture(s) — responder={RESPONDER_MODEL}, judge={JUDGE_MODEL}\n")
    results: list[Result] = []
    for fx in fixtures:
        if fx.rubric_id not in rubric_cache:
            rubric_cache[fx.rubric_id] = load_rubric(fx.rubric_id)
        rubric = rubric_cache[fx.rubric_id]

        try:
            responder_text, cache_read = call_responder(client, skill_md, fx.prompt)
        except Exception as e:
            print(f"  [{fx.id}] responder error: {type(e).__name__}: {e}")
            results.append(Result(fx.id, fx.expected, f"<error: {e}>", False, 1, str(e), 0))
            continue

        try:
            verdict = call_judge(client, rubric, fx, responder_text)
        except Exception as e:
            print(f"  [{fx.id}] judge error: {type(e).__name__}: {e}")
            results.append(Result(fx.id, fx.expected, responder_text, False, 1, f"judge error: {e}", cache_read))
            continue

        r = Result(
            fixture_id=fx.id,
            expected=fx.expected,
            responder_output=responder_text,
            judge_pass=bool(verdict["pass"]),
            judge_confidence=int(verdict["confidence"]),
            judge_reason=str(verdict["reason"]),
            cache_read_tokens=cache_read,
        )
        results.append(r)
        flag = "PASS" if r.judge_pass else "FAIL"
        cache_note = f" cache_read={r.cache_read_tokens}" if r.cache_read_tokens else ""
        print(f"  [{r.fixture_id}] {flag} (conf={r.judge_confidence}){cache_note} — {r.judge_reason}")

    passed = sum(1 for r in results if r.judge_pass)
    failed = len(results) - passed
    total_cache_read = sum(r.cache_read_tokens for r in results)
    print(f"\n{'=' * 60}")
    print(f"Summary: {passed} passed, {failed} failed (of {len(results)})")
    if total_cache_read:
        print(f"Cache reads: {total_cache_read} tokens (responder system prompt)")
    if failed and args.strict:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
