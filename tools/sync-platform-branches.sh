#!/usr/bin/env bash
# Rebase each platform branch onto main, push if clean.
#
# The platform branches (cc-windows, codex-macos, codex-windows) carry one
# commit each on top of main: a custom INSTALL.md plus a small README banner.
# When main moves, those branches need to be rebased forward.
#
# Conflicts only happen if main itself touches INSTALL.md or the README banner
# region — fix those by hand once and re-run.
#
# Usage:
#   tools/sync-platform-branches.sh            # rebase + push all
#   tools/sync-platform-branches.sh --dry-run  # rebase locally, don't push
#   tools/sync-platform-branches.sh cc-windows # one branch only
set -euo pipefail

PLATFORM_BRANCHES=(cc-windows codex-macos codex-windows)
DRY_RUN=0
TARGETS=()

for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1 ;;
    -h|--help)
      sed -n '2,/^set -/p' "$0" | sed 's/^# \?//'
      exit 0
      ;;
    *) TARGETS+=("$arg") ;;
  esac
done

if [ ${#TARGETS[@]} -eq 0 ]; then
  TARGETS=("${PLATFORM_BRANCHES[@]}")
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "error: working tree has uncommitted changes; commit or stash first." >&2
  exit 1
fi

START_BRANCH=$(git rev-parse --abbrev-ref HEAD)
trap 'git checkout "$START_BRANCH" >/dev/null 2>&1 || true' EXIT

echo "fetching origin…"
git fetch origin --quiet

echo "rebasing main onto origin/main…"
git checkout main --quiet
git rebase origin/main

for branch in "${TARGETS[@]}"; do
  echo
  echo "=== $branch ==="
  if ! git show-ref --verify --quiet "refs/heads/$branch"; then
    echo "  skip: local branch $branch doesn't exist (create it first)"
    continue
  fi
  git checkout "$branch" --quiet
  if git rebase main; then
    echo "  rebased clean"
    if [ "$DRY_RUN" -eq 0 ]; then
      git push --force-with-lease origin "$branch"
      echo "  pushed"
    else
      echo "  (dry run — not pushing)"
    fi
  else
    echo "  CONFLICT on $branch — fix by hand and run 'git rebase --continue', then re-run this script"
    git rebase --abort
    exit 1
  fi
done

echo
echo "all branches synced."
