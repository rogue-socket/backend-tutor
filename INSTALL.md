# Install — macOS / Linux + Codex CLI

This is the **`codex-macos`** branch. It ships only the macOS / Linux + OpenAI Codex CLI install path. For the full matrix (Windows / Claude Code / other harnesses), switch to [`main`](https://github.com/rogue-socket/backend-tutor/tree/main).

---

## Clone

```bash
git clone https://github.com/rogue-socket/backend-tutor ~/Documents/backend-tutor
```

Codex reads `AGENTS.md` from the project's working directory — there's no symlink step. The repo *is* the install.

## Use it

```bash
cd ~/Documents/backend-tutor
codex "start the backend tutor"
```

The agent picks up `AGENTS.md`, which routes it to `SKILL.md`. From there it runs the vibe check → lane → orientation + language → workspace setup → first lesson.

The workspace lives at `~/backend-dev/`. Resume any time with `/continue` (works across harnesses — the workspace is the bridge).

## Update

```bash
cd ~/Documents/backend-tutor
git pull origin codex-macos
```

## Note for projects

If you'd rather invoke the tutor from inside another project's directory (so it has cwd context for codebase-specific questions), copy `AGENTS.md` into that project:

```bash
cp ~/Documents/backend-tutor/AGENTS.md /path/to/your/project/
```

`AGENTS.md` is a thin pointer at `SKILL.md`; copying it into a project lets Codex find the tutor while sitting in your own codebase.
