# CLAUDE.md — ending_back (backend-tutor skill source)

This repo *is* a Claude Code / agent-harness skill: **backend-tutor**, a sibling of `system-design-tutor` and `ai-systems-tutor`. There is no application to run from this repo. The "code" here is the skill content — `SKILL.md` (router), `references/` (lazy-loaded mode files), `assets/` (workspace seeds + builder-first scaffolding the learner copies into `~/backend-dev/`).

## Layout

```
SKILL.md                              router + onboarding (3 lanes, 2 orientations)
references/
  curriculum.md                       T0–T11 topic tree + sources
  builder-first.md                    10-loop spec
  exercise-bank.md                    catalog by tier
  incidents.md                        real backend postmortems
  practical-mode.md                   exercise playbook (multi-language)
  theory-modes.md                     5 teaching modes
  session-control.md                  pause/resume/checkpoint
  spaced-repetition.md                SR queue + progress.json schema
assets/
  workspace-README.md                 copied to learner's ~/backend-dev/README.md
  COMMANDS.md                         slash commands reference
  progress-template.json              initial progress.json
  builder-first/
    setup/README.md                   toolchain setup per language
    go/loop-1-bare-crud/ … loop-10-loadtest/   runnable Go scaffolds
    _spec-only/loop-1-bare-crud/      language-agnostic mirror (Loop 1 only currently)
```

## How the skill is meant to be installed

### macOS / Linux

```bash
ln -s ~/Documents/ending_back ~/.claude/skills/backend-tutor
```

Mirrors how `ai-systems-tutor` is installed (source at `~/Documents/ai-system-tutor/`, symlinked to `~/.claude/skills/ai-systems-tutor/`).

### Windows (PowerShell, no admin / Developer Mode required)

Use a directory junction — junctions work on standard Windows for any user, unlike symlinks which need admin or Developer Mode:

```powershell
# Adjust the source path to wherever you cloned the repo
$src = "$env:USERPROFILE\Documents\ending_back"
$dst = "$env:USERPROFILE\.claude\skills\backend-tutor"

New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.claude\skills" | Out-Null
cmd /c mklink /J "$dst" "$src"
```

Verify:

```powershell
Get-Item $dst | Select-Object Name, Target
# Target should print the source path
```

Other harnesses (Codex CLI, Copilot CLI, Cursor, Aider) don't need the junction — they read `AGENTS.md` from the project's working directory. Either `cd` into the cloned repo before invoking the agent, or copy `AGENTS.md` into the project where you're running it.

## Invariants

- `SKILL.md` is the router; it must not load reference files speculatively. References are loaded only when the relevant mode is active.
- Workspace path is `~/backend-dev/`. Hardcoded across SKILL.md, references, and assets — change in all three if it ever moves.
- Default language for builder-first scaffolding is Go; Python is the planned secondary (not yet shipped).
- Loops 1 + 2 are standalone Go modules. Loops 3+ are deltas to the learner's evolving project.

## Sibling skills

- `system-design-tutor` (~/.claude/skills/system-design-tutor/) — owns architecture-at-scale design. Backend-tutor hands off "design X for 100M users" requests to it.
- `ai-systems-tutor` (~/Documents/ai-system-tutor/) — owns LLM-specific infra. Backend-tutor hands off agent-loop / RAG questions to it.

Don't duplicate content from siblings; cross-link.

## Recent commits

```
af198ed Initial commit
```

(Build is uncommitted as of writing — see git status.)
