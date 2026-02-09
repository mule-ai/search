#!/bin/bash
# Compare Search Results Across Instances
# Shows how different SearXNG instances return different results

if [ -z "$1" ]; then
  echo "Usage: $0 \"search query\""
  exit 1
fi

QUERY="$1"
OUTPUT_DIR="comparison-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$OUTPUT_DIR"

# List of SearXNG instances to compare
INSTANCES=(
  "https://search.butler.ooo"
  "https://searx.be"
  "https://searx.work"
)

echo "Comparing results across instances for: $QUERY"
echo ""

for instance in "${INSTANCES[@]}"; do
  echo "Testing: $instance"

  # Extract instance name for filename
  name=$(echo "$instance" | sed 's|https://||' | sed 's|http://||' | tr '/' '_')

  # Get results
  search -i "$instance" -f json -n 10 "$QUERY" > "$OUTPUT_DIR/${name}.json"

  # Count results and display summary
  count=$(jq '.results | length' "$OUTPUT_DIR/${name}.json")
  engines=$(jq -r '.results[] | .engine' "$OUTPUT_DIR/${name}.json" | sort -u | wc -l)

  echo "  Results: $count"
  echo "  Engines: $engines"
  echo ""

  sleep 1
done

echo "Comparison complete!"
echo "Results saved to: $OUTPUT_DIR"

# Generate comparison summary
echo "# Search Comparison Report" > "$OUTPUT_DIR/summary.md"
echo "**Query:** $QUERY" >> "$OUTPUT_DIR/summary.md"
echo "**Date:** $(date)" >> "$OUTPUT_DIR/summary.md"
echo "" >> "$OUTPUT_DIR/summary.md"

for instance in "${INSTANCES[@]}"; do
  name=$(echo "$instance" | sed 's|https://||' | sed 's|http://||' | tr '/' '_')
  file="$OUTPUT_DIR/${name}.json"

  if [ -f "$file" ]; then
    echo "## $instance" >> "$OUTPUT_DIR/summary.md"
    count=$(jq '.results | length' "$file")
    echo "**Results:** $count" >> "$OUTPUT_DIR/summary.md"
    echo "" >> "$OUTPUT_DIR/summary.md"
    jq -r '.results[0:3] | .[] | "- [\(.title)](\(.url))"' "$file" >> "$OUTPUT_DIR/summary.md"
    echo "" >> "$OUTPUT_DIR/summary.md"
  fi
done

echo "Summary: $OUTPUT_DIR/summary.md"