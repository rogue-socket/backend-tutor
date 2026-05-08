# Mode-routing rubric

Used by the LLM judge to grade mode-routing-test responses.

## What's being judged

The responder was given:
1. SKILL.md (backend-tutor's full content) as system context
2. A user prompt assumed to be *inside* an active backend-tutor session
3. The mode-dispatch table from SKILL.md as a list of modes (theory, design-review, builder-first-loop, handoff-system-design, handoff-ai-systems, ...)

It was asked to pick the mode that should handle the prompt and explain why in one sentence.

## Pass criteria

Response passes if it picks the mode matching the fixture's `expected` field. Hand-off fixtures (`expected: "handoff-system-design"` or `"handoff-ai-systems"`) require the responder to recognize the request is out of scope and route to the sibling skill.

## Fail criteria

- Picks the wrong mode (e.g., theory mode for a "design X for 100M users" prompt that should hand off to system-design-tutor)
- Tries to handle a request in-skill that should hand off
- Refuses or hedges without committing to a pick

## Confidence + output schema

Same as `activation.md`: `{pass: bool, confidence: 1-5, reason: str}`.
