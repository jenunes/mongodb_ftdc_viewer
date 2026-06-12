#!/usr/bin/env bash
# Push the current branch to your fork and open a PR against devops-land/mongodb_ftdc_viewer.
# Prerequisites: gh auth login -h github.com, and at least one commit on the branch.
set -euo pipefail
cd "$(dirname "$0")/.."

UPSTREAM="devops-land/mongodb_ftdc_viewer"
BRANCH="$(git branch --show-current)"

if [[ "$BRANCH" == "main" ]]; then
  echo "Create a feature branch first, e.g.: git checkout -b contrib/improvements"
  exit 1
fi

if ! gh auth status -h github.com &>/dev/null; then
  echo "GitHub CLI is not authenticated. Run: gh auth login -h github.com"
  exit 1
fi

if git diff --quiet && git diff --cached --quiet; then
  :
else
  echo "You have unstaged or uncommitted changes. Commit first, or run ./scripts/stage-for-upstream.sh"
  exit 1
fi

if ! git rev-list --count "upstream/main..HEAD" | grep -qv '^0$'; then
  echo "No commits ahead of upstream/main. Stage and commit your changes first."
  exit 1
fi

if ! git remote get-url origin &>/dev/null; then
  GH_USER="$(gh api user --jq .login)"
  echo "Creating fork and adding remote: origin -> https://github.com/${GH_USER}/mongodb_ftdc_viewer.git"
  gh repo fork "$UPSTREAM" --remote=false --clone=false
  git remote add origin "https://github.com/${GH_USER}/mongodb_ftdc_viewer.git"
fi

GH_USER="$(gh api user --jq .login)"

git push -u origin "$BRANCH"
gh pr create \
  --repo "$UPSTREAM" \
  --base main \
  --head "${GH_USER}:${BRANCH}" \
  --title "Fix MongoDB 8.x FTDC metrics and dashboard improvements" \
  --body "$(cat <<'EOF'
## Summary
- Support MongoDB 8.x ticket metrics under `serverStatus.queues.execution.*`
- Flatten nested FTDC metric documents in the exporter
- Simplify WiredTiger Tickets panel to read/write available tickets only

## Test plan
- [ ] Re-ingest FTDC data from MongoDB 8.0.17 and confirm WiredTiger Tickets panel shows data
- [ ] `go test ./ftdc/...` in `ftdc_exporter`
- [ ] Verify legacy MongoDB versions still show ticket metrics where applicable
EOF
)"
