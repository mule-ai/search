#!/bin/bash
# Daily News Report Generator
# Fetches latest news on specified topics and creates a report

REPORT_FILE="daily-report-$(date +%Y%m%d).md"
TEMP_DIR=$(mktemp -d)

# Topics to track
TOPICS=(
  "artificial intelligence"
  "machine learning"
  "programming languages"
  "kubernetes"
  "cybersecurity"
)

echo "# Daily Tech News Report" > "$REPORT_FILE"
echo "**Date:** $(date '+%Y-%m-%d')" >> "$REPORT_FILE"
echo "**Generated:** $(date '+%H:%M:%S')" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

for topic in "${TOPICS[@]}"; do
  echo "Fetching news for: $topic"

  echo "## $topic" >> "$REPORT_FILE"
  echo "" >> "$REPORT_FILE"

  # Get news from today, markdown format
  search -c news --time day -f json "$topic" | \
    jq -r '.results[0:5] | .[] | "- [\(.title)](\(.url))"' >> "$REPORT_FILE"

  echo "" >> "$REPORT_FILE"
  echo "---" >> "$REPORT_FILE"
  echo "" >> "$REPORT_FILE"

  sleep 2
done

echo "Report saved to: $REPORT_FILE"

# Optional: Display report
if command -v less &> /dev/null; then
  less "$REPORT_FILE"
fi

rm -rf "$TEMP_DIR"