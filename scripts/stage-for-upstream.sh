#!/usr/bin/env bash
# Stage files intended for the upstream PR. Run from repo root after all improvements are done.
set -euo pipefail
cd "$(dirname "$0")/.."

git add \
  .gitignore \
  ftdc_exporter/ftdc/normalize.go \
  ftdc_exporter/ftdc/iterator.go \
  ftdc_exporter/ftdc/normalize_test.go \
  ftdc_exporter/ftdc/testdata/metrics_to_get.txt \
  metrics_to_get.txt \
  grafana/dashboards/dashboard.json \
  run.sh \
  scripts/

echo "Staged for upstream PR:"
git diff --cached --stat
echo
echo "Review:  git diff --cached"
echo "Commit:  git commit -m \"\$(cat <<'EOF'"
echo "Fix MongoDB 8.x FTDC metrics and WiredTiger Tickets dashboard."
echo
echo "Add queues.execution ticket fields, flatten nested FTDC metric documents,"
echo "and limit the WiredTiger Tickets panel to read/write available tickets."
echo "EOF"
echo ")\""
echo "Publish: ./scripts/publish-to-upstream.sh"
