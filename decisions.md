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
