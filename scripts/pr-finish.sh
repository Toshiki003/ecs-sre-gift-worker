#!/usr/bin/env bash
set -euo pipefail

TITLE="${1:-}"
if [ -z "$TITLE" ]; then
  echo "Usage: ./scripts/pr-finish.sh \"feat: ...\""
  exit 1
fi

get_default_branch() {
  local branch=""

  if command -v gh >/dev/null 2>&1; then
    branch="$(gh repo view --json defaultBranchRef -q .defaultBranchRef.name 2>/dev/null || true)"
  fi

  if [ -z "$branch" ]; then
    branch="$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed -E 's@^refs/remotes/origin/@@' || true)"
  fi

  if [ -z "$branch" ]; then
    branch="main"
  fi

  printf '%s' "$branch"
}

# Ensure there are changes
if [ -z "$(git status --porcelain)" ]; then
  echo "No changes to commit."
  exit 0
fi

current_branch="$(git rev-parse --abbrev-ref HEAD)"
default_branch="$(get_default_branch)"

if [ "$current_branch" = "$default_branch" ] || [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ]; then
  echo "Refusing to run on protected/default branch: ${current_branch}"
  echo "Create a feature branch first (example: ./scripts/pr.sh \"feat: ...\")."
  exit 1
fi

if [ "${ALLOW_ANY_BRANCH:-false}" != "true" ] && ! [[ "$current_branch" =~ ^(feat|fix|docs|refactor|test|chore|ci|build|perf|hotfix)/ ]]; then
  echo "Refusing to run on branch '${current_branch}' (unexpected naming)."
  echo "Allowed prefixes: feat/, fix/, docs/, refactor/, test/, chore/, ci/, build/, perf/, hotfix/"
  echo "If this is intentional, rerun with: ALLOW_ANY_BRANCH=true ./scripts/pr-finish.sh \"${TITLE}\""
  exit 1
fi

# --- Label auto-detection from title prefix ---
label=""
case "$TITLE" in
  feat:*|feat\(*) label="feat" ;;
  fix:*|fix\(*)   label="fix" ;;
  docs:*|docs\(*) label="docs" ;;
  refactor:*|refactor\(*) label="refactor" ;;
  test:*|test\(*) label="test" ;;
  chore:*|chore\(*) label="chore" ;;
esac

# --- Commit with Co-Authored-By ---
git add -A
git commit -m "$TITLE

Co-Authored-By: Claude <noreply@anthropic.com>"

git push -u origin HEAD

# --- Build PR body from template ---
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEMPLATE="${REPO_ROOT}/.github/PULL_REQUEST_TEMPLATE.md"

pr_body_args=()
if [ -f "$TEMPLATE" ]; then
  pr_body_args=(--body-file "$TEMPLATE")
else
  pr_body_args=(--fill)
fi

# --- Create PR ---
set +e
pr_create_output="$(gh pr create --title "$TITLE" "${pr_body_args[@]}" 2>&1)"
pr_create_status=$?
set -e

if [ $pr_create_status -ne 0 ]; then
  echo "Failed to create PR:"
  echo "$pr_create_output"
  exit $pr_create_status
fi

pr_url="$(printf '%s\n' "$pr_create_output" | tail -n1 | tr -d '\r')"
if ! printf '%s' "$pr_url" | grep -Eq '^https://github\.com/.+/pull/[0-9]+$'; then
  fallback_url="$(gh pr view --json url -q .url 2>/dev/null || true)"
  if printf '%s' "$fallback_url" | grep -Eq '^https://github\.com/.+/pull/[0-9]+$'; then
    pr_url="$fallback_url"
  else
    echo "PR URL could not be parsed safely."
    echo "$pr_create_output"
    exit 1
  fi
fi

echo "PR created: ${pr_url}"

# Apply detected label
if [ -n "$label" ]; then
  gh pr edit "$pr_url" --add-label "$label" >/dev/null 2>&1 || true
fi

# Always add ai-review label for automated review
gh pr edit "$pr_url" --add-label "ai-review" >/dev/null 2>&1 || true

echo "Done."
