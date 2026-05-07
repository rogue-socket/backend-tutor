# [Service] — [date] load test postmortem

> Copy this file to `postmortem-YYYY-MM-DD.md` and fill in. One file per significant load-test incident; treat each as practice for the real thing.

## Summary

[1-2 sentences. What broke, what was the user-visible behaviour, how long until recovery.]

## Impact

- **Duration:** [HH:MM:SS to HH:MM:SS]
- **User-visible:** [percent of requests affected; what they saw]
- **Data integrity:** [any data loss, corruption, or inconsistency]
- **Internal cost:** [DB pressure, paged engineers, etc.]

## Timeline

| Time | Event |
|------|-------|
| HH:MM | Load test started, ramping to 1000 RPS |
| HH:MM | Error rate first crossed 1% |
| HH:MM | Alert fired (or didn't — note here) |
| HH:MM | Investigation began |
| HH:MM | Hypothesis 1: ... |
| HH:MM | Hypothesis 1 ruled out via [evidence] |
| HH:MM | Hypothesis 2: ... |
| HH:MM | Root contributing factor identified |
| HH:MM | Fix applied |
| HH:MM | Recovery confirmed (error rate, p99 back to baseline) |

## Contributing factors

> Plural, system-level. Avoid the single-cause / single-owner framing; failures are usually a chain of small choices.

1. **[Factor name].** [Description.] [Evidence: dashboard/log/trace pointing at it.]
2. **[Factor name].** [Description.] [Evidence.]
3. **[Factor name].** [Description.] [Evidence.]

## What went well

- [E.g., the SLO alert from Loop 9 fired within 90 seconds]
- [E.g., graceful shutdown from Loop 7 meant zero failed requests during redeploy]
- [E.g., the runbook had the right first-step]

## What didn't go well

- [E.g., the runbook didn't cover this specific failure mode]
- [E.g., observability gap: I had logs and traces but no metric on the load-bearing thing]
- [E.g., took 30 minutes to find what would have been obvious with one missing dashboard panel]

## Action items

- [ ] **[Concrete change.]** Owner: [name/me]. Due: [date]. Why: [the thing this prevents].
- [ ] **[Test added to prevent regression.]** [E.g., "load test in CI at 100 RPS with 50% read/50% write mix"]
- [ ] **[Documentation / runbook update.]** [Specific section, specific page.]
- [ ] **[Observability gap closed.]** [What metric / log / dashboard.]

## Lessons

[2-3 sentences. The single most important thing future-you should remember about this incident. Concrete enough to be searchable; honest enough to be useful.]

---

## Notes for filling this in

- **Avoid "human error" as a contributing factor.** It's almost never wrong but always useless. The interesting question is: what about the system let the human error matter? E.g., "engineer ran wrong command" → "production access tools accept the same commands as dev tools without confirmation."
- **Quote your evidence.** "p99 spiked at 14:32" without a screenshot is a memory; with a screenshot it's a fact.
- **Action items are concrete or fictional.** "Improve observability" is fictional. "Add `db_pool_acquire_wait_seconds` histogram to the RED dashboard" is concrete.
- **Lessons should make future-you faster.** If the lesson is "be more careful," delete it. If the lesson is "always check pool stats before adding capacity," keep it.
