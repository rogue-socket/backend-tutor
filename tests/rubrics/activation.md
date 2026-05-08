# Activation rubric

Used by the LLM judge to grade activation-test responses.

## What's being judged

The responder was given:
1. SKILL.md (backend-tutor's full content) as system context
2. A list of available skills (backend-tutor, system-design-tutor, ai-systems-tutor, "no skill matches")
3. A user prompt

It was asked to pick the skill (or "no skill matches") that should handle the prompt, and explain why in one sentence.

## Pass criteria

A response **passes** if it picks the skill matching the fixture's `expected` field:

- `expected: "trigger"` → response must select `backend-tutor`
- `expected: "no-trigger"` → response must select something other than `backend-tutor` (a sibling skill, or "no skill matches")

The reasoning quality is *secondary* — we're testing whether the description triggers correctly, not whether the model writes good prose.

## Fail criteria

- Picks `backend-tutor` when `expected: "no-trigger"` (false positive — backend-tutor steals queries from siblings or unrelated)
- Picks something else when `expected: "trigger"` (false negative — backend-tutor misses real learners)
- Refuses or hedges without committing to a pick (treat as fail; ambiguity is a description-tuning bug)

## Confidence

- 5 — clear-cut match or non-match per the criteria above
- 3 — borderline: response picks correctly but the reasoning suggests it could have gone either way
- 1 — response is malformed, hedges heavily, or grading required heavy interpretation

Low-confidence passes are still passes; flag them in `reason` so a human can re-check.

## Output schema

`{pass: bool, confidence: 1-5, reason: str}`. `reason` should be one short sentence — what the responder picked and whether that matches `expected`.
