#!/bin/bash
# Batch Search Script
# Search multiple queries and save results

# Check if query file is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <queries.txt>"
  echo "Format: One query per line in queries.txt"
  exit 1
fi

QUERY_FILE="$1"
OUTPUT_DIR="search-results-$(date +%Y%m%d-%H%M%S)"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "Running batch search from: $QUERY_FILE"
echo "Results will be saved to: $OUTPUT_DIR"
echo ""

# Counter
count=0

# Read each line and search
while IFS= read -r query; do
  # Skip empty lines and comments
  [[ -z "$query" || "$query" =~ ^#.*$ ]] && continue

  count=$((count + 1))
  echo "[$count] Searching: $query"

  # Sanitize filename
  filename=$(echo "$query" | tr ' ' '_' | tr -d '[:punct:]' | cut -c1-50)

  # Run search and save results
  search -f json "$query" > "$OUTPUT_DIR/${filename}.json"

  # Also save markdown for easy reading
  search -f markdown "$query" > "$OUTPUT_DIR/${filename}.md"

  # Small delay between requests
  sleep 1
done < "$QUERY_FILE"

echo ""
echo "Batch search complete!"
echo "Results saved to: $OUTPUT_DIR"
echo "Total searches: $count"