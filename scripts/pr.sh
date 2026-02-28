#!/usr/bin/env bash
set -euo pipefail

TITLE="${1:-}"
if [ -z "$TITLE" ]; then
  echo "Usage: ./scripts/pr.sh \"<type>: ...\""
  echo "Types: feat / fix / docs / refactor / test / chore / ci / build / perf / hotfix"
  exit 1
fi

# Require clean working tree (avoid accidental context switches)
if [ -n "$(git status --porcelain)" ]; then
  echo "Working tree is not clean. Commit/stash your changes first."
  git status --porcelain
  exit 1
fi

# Determine default branch via GitHub (fallback: main)
default_branch="$(gh repo view --json defaultBranchRef -q .defaultBranchRef.name 2>/dev/null || true)"
if [ -z "${default_branch}" ]; then
  default_branch="main"
fi

# Update default branch safely
git checkout "$default_branch"
git pull --rebase

# Detect branch prefix from title (fallback: feat)
prefix="feat"
case "$TITLE" in
  feat:*|feat\(*)           prefix="feat" ;;
  fix:*|fix\(*)             prefix="fix" ;;
  docs:*|docs\(*)           prefix="docs" ;;
  refactor:*|refactor\(*)   prefix="refactor" ;;
  test:*|test\(*)           prefix="test" ;;
  chore:*|chore\(*)         prefix="chore" ;;
  ci:*|ci\(*)               prefix="ci" ;;
  build:*|build\(*)         prefix="build" ;;
  perf:*|perf\(*)           prefix="perf" ;;
  hotfix:*|hotfix\(*)       prefix="hotfix" ;;
esac

# Create branch name
slug=$(echo "$TITLE" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g' | sed -E 's/^-|-$//g')
branch="${prefix}/${slug}-$(date +%Y%m%d%H%M%S)"
git checkout -b "$branch"

echo "Now implement your changes, then run:"
echo "  ./scripts/pr-finish.sh \"${TITLE}\""
