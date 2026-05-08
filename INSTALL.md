# Install — Windows + Claude Code

This is the **`cc-windows`** branch. It ships only the Windows + Claude Code install path. For the full matrix (macOS / Linux / other harnesses), switch to [`main`](https://github.com/rogue-socket/backend-tutor/tree/main).

---

## PowerShell (no admin / Developer Mode required)

Use a directory junction. Junctions work for any user on standard Windows without elevated permissions, unlike symlinks.

```powershell
git clone https://github.com/rogue-socket/backend-tutor "$env:USERPROFILE\Documents\backend-tutor"

$src = "$env:USERPROFILE\Documents\backend-tutor"
$dst = "$env:USERPROFILE\.claude\skills\backend-tutor"
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.claude\skills" | Out-Null
cmd /c mklink /J "$dst" "$src"
```

Verify:

```powershell
Get-Item $dst | Select-Object Name, Target
# Target should print the source path
```

## Use it

```
> start the backend tutor
```

The skill auto-routes from the trigger phrase: vibe check → lane → orientation + language → workspace setup → first lesson.

The workspace lives at `%USERPROFILE%\backend-dev\`. Resume any time with `> /continue`.

## Update

When upstream `main` changes, update by pulling the source repo (the junction stays pointed at the live source):

```powershell
cd "$env:USERPROFILE\Documents\backend-tutor"
git pull origin cc-windows
```
