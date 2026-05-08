# Decisions

Durable architectural and methodological decisions for this project, with rationale. Each entry: dated heading, the decision in one line, then **Why:** and **How to apply:** lines. Append-only — supersede old entries with new ones rather than rewriting.

## 2026-05-07 — Three lanes (Foundations / Working / Senior), not five

Diagnostic routing puts learners into one of three lanes; adjacent-domain entrants (frontend/SRE/data-eng/ML/EE) are a `working_mode = adjacent_domain` variant of Working, not a separate lane. Non-coder is dropped entirely.

**Why:** ai-systems-tutor uses 5 lanes (beginner / bridge / middle / expert / non_coder); the persona round during this skill's design surfaced that 5 was over-engineered for the audience that comes to a *backend* tutor. Backend has fewer entry archetypes than AI systems; collapsing to 3 cuts router complexity and authoring overhead. Non-coder is dropped because backend is by definition a coding role; a non-coder asking for backend literacy is best routed elsewhere.

**How to apply:** When extending or modifying the lane router (SKILL.md Step 2), preserve the 3-lane shape. New archetypes (e.g., a different adjacent domain like ML/research) get handled by extending the adjacent-domain examples, not by adding a fourth lane. If a fourth lane ever feels necessary, that's a signal the diagnostic is doing too much work — split the diagnostic instead.

---

## 2026-05-07 — Language as a top-level learner config

`progress.json` carries `learner.language` as a first-class field. Go primary (prefilled scaffolding), Python/FastAPI secondary (prefilled), Node/TS/Java/Kotlin/Rust spec-only.

**Why:** Both sibling skills default to Python; for a backend tutor that's wrong. Senior backend job listings (Stripe, Cloudflare, Shopify samples) lean Go-heavy; backend Python (FastAPI, Django) and Node are also major; Java/Kotlin still ship enterprise; Rust is rising. Defaulting Python would mis-serve the majority of the audience and make the skill feel out-of-date for interview prep.

**How to apply:** When authoring new builder-first loops or exercises, check `learner.language` and either ship language-specific scaffolding (Go and Python) or hand a spec for the others. Don't hardcode language-specific code in reference files; keep concept-level explanations language-agnostic. When a new "first-class" language is added (e.g., authoring Python Loop 2-10 scaffolds), the assets/builder-first/<lang>/ tree mirrors assets/builder-first/go/ structurally.

---

## 2026-05-07 — Both orientations shipped (foundations-first + builder-first)

Learners pick at onboarding (Step 2.5). Foundations-first walks T0–T11 in tier order; builder-first runs the 10-loop spec growing one "links" service through deliberate breaks.

**Why:** Both audiences exist and need different starts. A foundations-first learner wants mental models before code; a builder-first learner wants momentum and learns through failure. Forcing one path loses the other. ai-systems-tutor reached the same conclusion (foundations-first vs builder-first); the design transfers cleanly.

**How to apply:** When authoring curriculum content, write tier topics that work for both orderings. The builder-first path is *not* a license to skip foundations — it's a different *order*. Foundations get filled in mid-loop as the service breaks; the tutor must enforce this with calibration probes during loops, not just at lane onboarding.

---

## 2026-05-07 — Source anchor strategy: fresh engineering blogs over DDIA-as-default

Primary sources for curriculum.md: Stripe / Cloudflare / GitHub / Discord / Figma / Shopify engineering blogs, OWASP API Top 10 (2023), 12-Factor App, *Software Engineering at Google* (2020), *Database Internals* (Petrov 2019), *Building Microservices 2e* (Newman 2021), *API Design Patterns* (Geewax 2021), the SRE books, HTTP RFCs (9110, 8446, 9000). DDIA cited only where no fresher canonical equivalent exists (T2 isolation levels, T11 replication concepts).

**Why:** User explicitly asked for "fresher and new" sources targeting "what real backend job listings ask for." DDIA (2017) is still the bible for a few topics, but using it as default anchor signals "old curriculum" to a learner who can find it everywhere; engineering blogs from companies actually shipping at scale signal current practice and pin claims to real systems.

**How to apply:** When adding a new topic to curriculum.md, search for an engineering-blog or recent-book primary source first; reach for DDIA only if nothing fresher covers the topic at sufficient depth. Always cite the source by URL or chapter pointer. Real-incident anchors (incidents.md) take priority over textbook examples — the war story is the load-bearing thing.

---

## 2026-05-07 — Cross-skill hand-off rules

When a learner asks for architecture-at-scale design ("design Twitter for 100M users"), backend-tutor hands off to system-design-tutor. When a learner asks about LLM-specific infra (agent loops, RAG, prompt caching), backend-tutor hands off to ai-systems-tutor. Backend-tutor stays on the implementation/operational side.

**Why:** The three skills could overlap heavily; without explicit boundaries each skill would either duplicate content or pull learners into its own framing of every question. Clear boundaries mean each skill stays focused and learners know which tool to reach for. Concretely: "build a sharded service" is backend-tutor; "design a globally distributed sharded datastore for 100M users" is system-design-tutor.

**How to apply:** SKILL.md's mode-dispatch table includes a "hand-off rules" section. When the tutor detects a request outside its scope, it should *suggest* the sibling skill explicitly rather than try to answer. Cross-link in references (e.g., curriculum.md T11 cross-links to system-design-tutor for architecture-level reasoning).

---

## 2026-05-07 — T11 Distributed Systems is implementation-level only

T11 in the curriculum covers replication, sharding, consensus, distributed locking, time/clocks, service discovery — at the *code-it-and-break-it* level (set up Postgres replication, simulate failover, observe split-brain). Architecture-level reasoning ("design a multi-region datastore for X scale") is system-design-tutor's territory.

**Why:** Without this boundary, T11 would balloon into "system design lite" and would directly compete with system-design-tutor. Backend engineers need both perspectives, but they need to know which is which. The implementation level is where backend-tutor adds value distinct from system-design-tutor; the architecture level is where it would duplicate badly.

**How to apply:** When adding to T11, ask "is this code-it / operate-it / observe-it?" If yes, it belongs here. If it's "draw the architecture", route to system-design-tutor instead. The curriculum.md T11 section explicitly cross-links the boundary.

---

## 2026-05-07 — Real-time as a cross-cutting path, not a tier

WebSockets, SSE, long-polling, real-time message protocols, presence, fan-out — assembled as a *path* through T0/T1/T3 plus a dedicated builder-first Loop 8, rather than its own tier.

**Why:** Real-time concerns cross multiple tiers (transport in T0, API patterns in T1, concurrency/pub-sub in T3). Carving out a separate tier would either force duplication of T0/T1/T3 content or fragment the relevant content unnaturally. A path that walks the relevant pieces in order, plus a hands-on loop, gives the same coverage with cleaner curriculum structure.

**How to apply:** When a learner has explicit real-time goals (live chat, collaborative editing, presence systems), route them through the path defined at the bottom of curriculum.md, with Loop 8 as the practical capstone. Don't create T12; instead extend the existing tiers as needed.

---

## 2026-05-07 — Skill source repo doubles as the project; install via symlink

`~/Documents/ending_back/` is both the git-versioned skill source and the directory that gets symlinked into `~/.claude/skills/backend-tutor/`. Mirrors how ai-systems-tutor is installed.

**Why:** This pattern keeps the skill's authored content under version control independently of any single user's `~/.claude/` installation. Updates flow through git (`git pull` updates the source; the symlink reflects changes immediately). Distribution to other users is `git clone + ln -s`. The alternative — authoring directly inside `~/.claude/skills/` — couples skill authoring to one user's machine state.

**How to apply:** Don't author skill content directly in `~/.claude/skills/backend-tutor/`. Always edit in `~/Documents/ending_back/` and rely on the symlink. When packaging for distribution to other users, ship the repo URL + a one-line `ln -s` command (mirroring ai-systems-tutor's install instructions).

---

## 2026-05-08 — Difficulty knob: one constraint at a time, never combined

When the learner says "make this easier" or "make this harder", apply exactly one constraint shift — never combine multiple downshifts or multiple upshifts in the same re-pitch. Easier downshifts pick from a fixed list (mock the dependency / shrink dataset / drop failure injection / narrow success criterion / pre-write boilerplate). Harder upshifts pick from a fixed list (one realistic failure / one scale constraint / remove a piece of scaffolding / add an SLO or budget).

**Why:** Calibration drifts unrecoverably when multiple constraints move at once. A learner who says "easier" wants the *concept* to land, not the entire exercise reshaped — mocking the dependency *and* shrinking the dataset *and* dropping the failure injection together changes what's being taught, not just how much. Keeping the knob to one shift per pull also makes the telemetry interpretable: `observed_difficulty: easier_than_planned` with `adjustments: ["mocked Postgres"]` is actionable signal for the next exercise; the same with three changes is noise. The first-practical promise is stated *once per session* (not per-exercise) for the same reason — repeated reminders teach the learner that "make this easier" is a normal mode rather than a calibration tool.

**How to apply:** When the knob is pulled, pick the single change that's most likely to land the concept (easier) or surface the next layer (harder), apply it, log it as one entry to `adjustments[]`. If the learner is *still* off-level after one shift, that's a calibration miss worth logging as a separate `adjustments[]` entry with the second shift — not a license to combine them in a single re-pitch. Refer to `references/practical-mode.md` *Difficulty knob* for the canonical constraint lists; don't invent new constraint categories on the fly without updating that file.

---

## 2026-05-08 — Forced-load triple-belt: rule + hook + anti-pattern

When a SKILL.md rule needs a specific reference file to be loaded before it can be followed (e.g., "ground every lesson in a real incident" needs `references/incidents.md`), wire it three ways: (1) the rule itself in the relevant section, (2) an explicit "load X first" hook paragraph in the same section, (3) a paired anti-pattern in the Anti-patterns block ("❌ Citing X from memory without loading the file first"). Don't rely on the rule alone.

**Why:** The 2026-05-08 paired-persona round showed `references/incidents.md` was opened by 0 of 8 tutor agents despite the rule being present at `SKILL.md:384`. A rule that's authored but never followed is worse than no rule — it produces fabrications (wrong dates, wrong services, wrong root causes) that the learner repeats in interviews. The triple-belt approach acknowledges that LLM-driven tutors don't always follow file-load instructions when the topic is in their training data; an anti-pattern explicitly framing memory-recitation as a failure mode catches the cases where the hook gets ignored. Rule alone is a 0/8; rule + hook + anti-pattern is unverified-but-believed-better — needs another persona round to confirm.

**How to apply:** When adding a rule that depends on a reference file (e.g., a new rule "ground every security topic in OWASP API Top 10 from `references/security-anchors.md`"), add all three pieces in the same edit: the rule in its section, the hook in the same section, the anti-pattern in the Anti-patterns block. Verify in the next persona round. If a rule keeps getting violated despite the triple-belt, escalate to inlining the canonical content in SKILL.md rather than relying on file-load.

---

## 2026-05-08 — Paired bad/good examples beat lists for anti-pattern calibration

The Anti-patterns block in SKILL.md is a list of one-liners. The companion file `references/anti-patterns-with-examples.md` pairs each one-liner with a short bad example and a short good replacement. The pair is the durable artifact; the list is the index.

**Why:** A list of anti-patterns is easy to skim and easy to forget. A paired bad/good example forces both a contentful contrast and a model of what the good move looks like. The bad/good split is deliberately *contentful*, not stylistic — the bad example often has nothing wrong with the prose, only with the *move* the tutor is making (e.g., answering a multi-part question with N/2 answers; teeing up the next step before finishing the current one). This shape transfers across other skills as a standard reference pattern. The audit framed it as a DIFFERENTIATOR vs siblings; the durable reason is calibration quality, not differentiation.

**How to apply:** When adding a new anti-pattern to SKILL.md's Anti-patterns block, also append a paired entry to `references/anti-patterns-with-examples.md` with: a short concrete-context setup, a bad example with the move framed, and a good replacement that names what the move *should* have been. Keep both short — pattern-matching beats exhaustiveness. The block-list and the pair-file should always carry the same items; if they drift, the pair-file is the source of truth.

---

## 2026-05-09 — Workspace subdir naming: `projects/`, not `exercises/`

The hands-on subdirectory inside `~/backend-dev/` is named `projects/`. The audit-checklist proposed `exercises/`; we kept `projects/`.

**Why:** `projects/` reads as "the fun-stuff folder" for builder-first learners — Loop 1 ships a working CRUD service that grows across loops, which is *one project that evolves*, not a sequence of disconnected exercises. `exercises/` is more neutral and would work for foundations-first standalone exercises, but builder-first is the orientation we expect to drive the most engagement, and the framing matters: a learner sees `~/backend-dev/projects/loop-3-migrations/` and gets the right mental model that this is real, owned code. Keeping `projects/` also avoided a sweep across SKILL.md, references, assets, and workspace-README that would have churned for naming-only reasons.

**How to apply:** Don't rename `projects/` to `exercises/` (the inverse drift would be cosmetic-cost-only). When adding a foundations-first standalone exercise that doesn't fit the builder-first single-project model, drop it under `projects/exercises/<slug>/` rather than spinning up a sibling top-level dir — the `projects/` name is broad enough to cover both cases. If a learner asks why their workspace says "projects" when they're doing standalone exercises, the answer is "same dir, two flavors of work."

---

## 2026-05-09 — Goal/timeline capture is optional, not required

SKILL.md Step 2.5's stated-goals/timeline probe ("interview prep / first backend role / payments service / SRE readiness / etc.") is optional. The audit checklist said capture MUST happen. We diverge intentionally.

**Why:** Casual learners ("I just want to learn backend, no specific deadline") don't have a goal-shaped answer to give, and forcing them to produce one creates friction at the exact moment we want momentum (Step 2.5 is right before the diagnostic, before any value has been delivered). The probe is *useful* when the learner has a goal — it routes them through `references/curriculum.md` § "Path suggestions by stated goal" — but useful-when-present doesn't justify required-for-all. Persona round F3 surfaced this; the optional-with-fallback design landed there.

**How to apply:** When the learner doesn't volunteer a goal/timeline, don't re-ask. Save `learner.stated_goals: null` and proceed. If they later mention one mid-course (mock interview, deadline drops), capture it then via `/config` and adjust pacing. Don't gate diagnostic entry on goal capture under any circumstances. If a future audit-style review re-flags this as MUST, rebut with the persona-round evidence.

---

## 2026-05-09 — Multi-branch distribution: small per-branch divergence + sync script

`main` is the universal branch. Three platform branches (`cc-windows`, `codex-macos`, `codex-windows`) each carry exactly one commit on top of `main`: a rewritten `INSTALL.md` (single-path) plus a small README banner naming the branch. `tools/sync-platform-branches.sh` rebases the platform branches onto main and pushes with `--force-with-lease`. Conflicts only arise when main itself touches `INSTALL.md` or the README banner region.

**Why:** §20 of the audit checklist called for multi-branch distribution as a DIFFERENTIATOR. The naive shape — branches that diverge across many files — would create perpetual rebase tax for one maintainer. Constraining divergence to two files (`INSTALL.md` + README's top region) keeps the maintenance cost near zero while still delivering the §20 benefit (a Windows + Codex CLI learner gets a single-path README without conditionals). The full multi-platform matrix lives in main's `INSTALL.md` so search-engine landings on main aren't broken; the platform branches are the leaner experience for users who already know what combo they're on.

**How to apply:** When updating install/onboarding content, prefer touching files *outside* `INSTALL.md` and the README banner — those changes auto-propagate to the platform branches via rebase. When you must touch `INSTALL.md`, expect to either (a) update the platform branches' `INSTALL.md` files in matching ways, or (b) accept conflicts during the next sync run and resolve by hand once. Run `tools/sync-platform-branches.sh --dry-run` after any main commit that touches install content; run without `--dry-run` to push. Don't add a fifth platform branch without re-evaluating whether the linear divergence cost still pencils out — at some count (probably ~6+) the matrix model breaks down and a generated-from-main approach becomes cheaper.

---

## 2026-05-09 — Test harnesses are split: deterministic stays stdlib-only; LLM-as-judge is a separate, opt-in runner

`tests/run_all.py` runs only the deterministic structural tests (frontmatter parse, JSON schema validation, reference-file presence, manifest pinning, viewer assets, practical-coverage). It depends on stdlib + `jsonschema` only. `tests/run_llm.py` is a separate runner for the LLM-as-judge behavior tests (activation, mode-routing). It's never wired into `run_all.py`'s default discovery, depends on `anthropic` + `ANTHROPIC_API_KEY`, and skips cleanly with exit 0 if either is missing.

**Why:** The two test categories solve different problems with incompatible operational profiles. Deterministic tests catch structural regressions (a missing reference file, a mistyped schema field) and need to be cheap, fast, and runnable on every push without API keys or quota — the kind of test a contributor without an Anthropic key can run, and CI can run on every PR for free. LLM-as-judge tests catch description/routing regressions that the deterministic suite can't see, but they take seconds-to-minutes per run, cost cents-to-dollars, and require an API key. Mixing the two would force one of three bad outcomes: (1) couple deterministic CI to an external API key, exposing a free-tier-blocking dependency for a free-tier-friendly suite; (2) require contributors without keys to skip parts of the deterministic run, making "run the tests" mean different things to different people; (3) silently make failures look the same regardless of cause, so a structural bug and a description regression read identically in CI output. Keeping them separate preserves a clean fast-path where the deterministic suite is the gate everyone shares, and the LLM suite is the deeper check run pre-PR or nightly. The naming convention reinforces this: `run_all.py` is "all the cheap tests," not literally all of them.

**How to apply:** When adding a new structural test (anything decidable from file contents alone — JSON schema, regex, file presence, frontmatter shape), drop it as `tests/test_<name>.py` and `run_all.py` will pick it up automatically. When adding an LLM-graded test (anything that requires an LLM to evaluate output quality — does the description trigger correctly, does mode routing pick the right reference set, does a tutor response avoid an anti-pattern), add fixtures to `tests/fixtures/llm/*.jsonl` and a rubric to `tests/rubrics/<id>.md`; `run_llm.py` picks them up. Don't add `import anthropic` to anything under the deterministic-test discovery path. Don't add `tests/run_llm.py` (or any LLM-judge entry point) to the `run_all.py` invocation. If a future test needs *both* a structural component and an LLM component, split it into two test files, one per harness, and cross-reference in comments. The split is also a CI-shape invariant: deterministic runs on every push; LLM runs pre-PR + nightly + on-demand only.
