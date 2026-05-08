# Forced-load verification: `references/incidents.md`

**Date:** 2026-05-08
**Trigger:** Backlog item — verify §7 forced-load triple-belt (added in 2026-05-08 Tier A audit response, G-A6) actually moves the needle on incidents.md load behavior.
**Baseline:** 2026-05-07/08 paired-persona round, **0 of 8** tutor agents opened `references/incidents.md` despite the rule being present at SKILL.md:401-403.
**Post-fix triple-belt under test:**
1. Rule (SKILL.md:401-403) — *"Every lesson references at least one real-world incident from `references/incidents.md`"*
2. Forced-load hook (SKILL.md:405) — bold *"Before citing any specific incident in a lesson, read `references/incidents.md`"* + don't-fabricate clause
3. Anti-pattern (SKILL.md:607) — *"❌ Citing an incident from memory without loading `references/incidents.md` first"*
4. Plus override-map rows at SKILL.md:363 (theory lesson → incidents.md) and SKILL.md:370 (user asks for incident → incidents.md)

## Methodology

8 fresh `general-purpose` subagents spawned in parallel. Each was told it was being invoked as the backend-tutor skill, given the path to SKILL.md, given a persona context and a turn-N user message naturally requesting a real incident, and asked to act per the skill's instructions with no priming about what was being measured. Reporting was structured: every Read tool call logged in order, then the user-facing response, then a META section answering Q1 (loaded incidents.md before citing? Y/N), Q2 (cited any specific incident at all? Y/N), Q3 (what triggered the load).

Personas mirrored the original round where named in backlog: Anika, Wei, Marcus, Tyler, Joseph, Devansh. Two backfilled (Priya, Robert) to reach 8.

## Result table

| # | Persona | Topic | Loaded incidents.md FIRST? | Cited specific incident? | Source of citation |
|---|---|---|---|---|---|
| 1 | Anika | cache stampede (T4) | ✅ YES | ✅ Discord stampede | from file |
| 2 | Wei | queue redelivery / at-least-once (T3) | ✅ YES | ✅ Knight Capital 2012 + Stripe blog | from file (gap-flagged honestly) |
| 3 | Marcus | saga / outbox (T3/T11) | ✅ YES | ❌ declined | gap; offered public-source walkthrough |
| 4 | Tyler | idempotency keys (T1) | ✅ YES | ❌ declined | gap; pointed to Stripe blog |
| 5 | Joseph | indexes (T2) | ✅ YES | ❌ declined | gap; pattern + public sources |
| 6 | Devansh | DB outage (T2) | ✅ YES | ✅ GitLab 2017-01-31 | from file (full specifics) |
| 7 | Priya | named outages (T0/T1) | ✅ YES | ✅ Cloudflare 2019-07-02, GitHub 2018-10-21, FB 2021-10-04 | from file |
| 8 | Robert | circuit breakers (T5) | ✅ YES | ✅ Roblox 2021-10-28, AWS DynamoDB 2015-09-20 | from file |

**Headline:** **8/8 loaded the file before responding** (vs. 0/8 pre-fix). 5/8 cited from the file with full specifics. 3/8 *correctly identified gaps in incidents.md and refused to fabricate* — invoking the don't-fabricate clause of the forced-load rule by name in their META Q3 line.

## What worked

- The triple-belt does its job. Multiple agents cited specific lines (SKILL.md:405, SKILL.md:607, the override-map row) in their META Q3 explanations — the rule, hook, and anti-pattern reinforce each other rather than being redundant.
- The don't-fabricate clause inside the forced-load hook ("If the relevant tier section is missing or thin, say so honestly") is load-bearing. Without it, the agents who hit gaps would have likely confabulated. With it, they explicitly flagged "incidents.md doesn't have this; here's a public source instead." That's *better* behavior than the original goal — fail-honest beats fail-silent.
- The override-map row "User asks for incident / case study → references/incidents.md" (SKILL.md:370) is doing real work as a redundancy check; agents whose prompts mentioned it cited it independently of the §7 rule.

## Content gaps in `incidents.md` surfaced by the round

Three of the 8 agents found tier sections that lack a clean named incident for the topic the persona asked about. These are now actionable backlog items:

| Tier | Gap | Persona that hit it |
|---|---|---|
| T1 | No named idempotency-failure incident (e.g., a payments double-charge with public RCA) | Tyler |
| T2 | No clean "bad index in production" incident (write amplification, lock from non-CONCURRENTLY index, planner picking wrong index at scale) | Joseph |
| T3/T11 | No saga / outbox case study (Uber Cadence origin posts and eBay outbox-pattern writeups exist publicly) | Marcus |
| T3 | No queue-redelivery / non-idempotent-consumer postmortem (Wei found Knight Capital adjacent but not queue-shaped) | Wei |

Adding these would lift the "from file" citation rate from 5/8 to plausibly 8/8. Worth a separate authoring pass.

## Recommendation

Triple-belt is **validated**. No need to escalate to inlining canonical incidents in SKILL.md. Backlog item §7 verification can be marked done. Open a follow-up to fill the four content gaps in `incidents.md` so future agents on those topics cite from the file rather than declining + redirecting to public sources.

## Caveats

- Subagent self-reporting: agents told us they loaded the file via the READ_LOG section. This matches the skill's actual behavior because the prompt didn't ask them to *predict* whether they'd load it — they were asked to *act* and then list what they read. Self-report risk remains non-zero but low.
- Persona set differs slightly from original 8 (2 backfilled). Direct head-to-head comparison would need the original persona names; backfilled personas (Priya, Robert) hit topics the file covers well, which biases their cite-from-file rate up. Even discounting them, 6/6 of the original-set personas loaded the file.
- Single round. A second independent round would harden the result. Low priority given the 8/8 → 0/8 swing is already overwhelming.
