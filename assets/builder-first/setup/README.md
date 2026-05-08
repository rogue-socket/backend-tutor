# Setup — first time

Three setup steps the tutor *will not run for you*. They're yours; the tutor coaches when you hit a snag.

The instructions branch on the language you picked (`learner.language` in `progress.json`).

---

## 1. Install the language toolchain

### Go (recommended for builder-first)

- macOS: `brew install go` — confirm with `go version` (need 1.24 or newer; loop-2's pgx pin requires it). See `LOOP_VERSIONS.md` at the repo root for the current pin manifest.
- Linux: official tarball at https://go.dev/dl/ — extract to `/usr/local/go`, add to PATH
- Windows: official installer at https://go.dev/dl/

### Python

- Use **uv** for project + venv management: `curl -LsSf https://astral.sh/uv/install.sh | sh`
- Confirm with `uv --version`. Python itself is managed via uv per project (`uv python install 3.12`).

### Node / TypeScript

- Use a version manager: **fnm** (`brew install fnm`) or nvm. Avoid system Node.
- Install Node 20 LTS or newer: `fnm install 20 && fnm use 20`.
- For TypeScript: project-local via `npm i -D typescript tsx`.

### Java / Kotlin

- Use **SDKMAN!** (`curl -s https://get.sdkman.io | bash`).
- `sdk install java 21-tem` — Java 21 is the first LTS with virtual threads.
- For Kotlin: `sdk install kotlin`.

### Rust

- `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh`.
- Confirm with `rustc --version` and `cargo --version`.

---

## 2. Install Docker (for Loop 2 onward)

Loop 1 doesn't need Docker. Loop 2 onward (Postgres, Redis, etc.) does.

- macOS: **OrbStack** is the lighter alternative (`brew install orbstack`). Docker Desktop also works.
- Linux: install Docker Engine via your distro's package manager + the `docker compose` plugin.
- Windows: Docker Desktop with WSL2 backend.

Confirm with `docker run --rm hello-world` and `docker compose version`.

---

## 3. Verify the workspace

```bash
ls ~/backend-dev/
# should show: README.md, COMMANDS.md, progress.json, session-state.md, notes/, projects/, ...

cat ~/backend-dev/progress.json | head -20
# should show your learner profile (level, orientation, language, etc.)
```

If `~/backend-dev/projects/loop-1-bare-crud/` exists, scaffolding was copied successfully and you're ready to start Loop 1.

If the projects directory is empty and `learner.orientation = builder_first`, the scaffolding copy failed during onboarding. Tell the tutor — it'll re-copy from the skill assets.

---

## When something fails

- **`go: command not found`** after install → restart your shell, or `source ~/.zshrc` / `~/.bashrc`
- **Docker daemon not running** → start Docker Desktop / OrbStack / `sudo systemctl start docker`
- **Port 8080 already in use** → another service is bound there. Either stop it, or change the port in `main.go` (search for `:8080`)
- **`go test -race` not detecting the race** → make sure you actually have unsynchronised access; sometimes the test is too fast to trigger. Add a `time.Sleep(time.Microsecond)` mid-Create temporarily, re-run.

For anything else: just ask the tutor. *"I'm stuck on setup step 1, here's the error: [paste]"* — that's the right shape.

---

## Once setup is done

Open Loop 1's directory and start:

```bash
cd ~/backend-dev/projects/loop-1-bare-crud
cat README.md     # the loop's overview
cat BREAK.md      # what we'll deliberately break
cat WIN.md        # success criteria
```

Then start writing code. The tutor is watching; ask questions when you have them.
