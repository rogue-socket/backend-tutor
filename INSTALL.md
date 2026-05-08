# Install — Windows + Codex CLI

This is the **`codex-windows`** branch. It ships only the Windows + OpenAI Codex CLI install path. For the full matrix (macOS / Linux / Claude Code / other harnesses), switch to [`main`](https://github.com/rogue-socket/backend-tutor/tree/main).

---

## Clone (PowerShell)

```powershell
git clone https://github.com/rogue-socket/backend-tutor "$env:USERPROFILE\Documents\backend-tutor"
```

Codex reads `AGENTS.md` from the project's working directory — there's no junction or symlink step. The repo *is* the install.

## Use it

```powershell
cd "$env:USERPROFILE\Documents\backend-tutor"
codex "start the backend tutor"
```

The agent picks up `AGENTS.md`, which routes it to `SKILL.md`. From there it runs the vibe check → lane → orientation + language → workspace setup → first lesson.

The workspace lives at `%USERPROFILE%\backend-dev\`. Resume any time with `/continue` (works across harnesses — the workspace is the bridge).

## Update

```powershell
cd "$env:USERPROFILE\Documents\backend-tutor"
git pull origin codex-windows
```

## Note for projects

If you'd rather invoke the tutor from inside another project's directory (so it has cwd context for codebase-specific questions), copy `AGENTS.md` into that project:

```powershell
Copy-Item "$env:USERPROFILE\Documents\backend-tutor\AGENTS.md" "C:\path\to\your\project\"
```

`AGENTS.md` is a thin pointer at `SKILL.md`; copying it into a project lets Codex find the tutor while sitting in your own codebase.
