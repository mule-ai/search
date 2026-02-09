#!/bin/bash
# Extract URLs from search results
# Useful for piping into other tools

if [ -z "$1" ]; then
  echo "Usage: $0 \"search query\" [options]"
  echo ""
  echo "Options:"
  echo "  --open     Open URLs in browser"
  echo "  --save     Save URLs to file"
  echo "  --unique   Remove duplicates"
  echo ""
  echo "Examples:"
  echo "  $0 \"golang tutorial\""
  echo "  $0 \"docker\" --unique"
  echo "  $0 \"kubernetes\" --save urls.txt"
  exit 1
fi

QUERY="$1"
shift
OPTIONS="$@"

# Get results as JSON and extract URLs
URLS=$(search -f json "$QUERY" | jq -r '.results[].url')

# Handle options
if [[ "$OPTIONS" == *"--unique"* ]]; then
  URLS=$(echo "$URLS" | sort -u)
fi

if [[ "$OPTIONS" == *"--save"* ]];
  # Extract filename
  SAVE_FILE=$(echo "$OPTIONS" | grep -oP '(?<=--save\s)[^\s]+' || echo "urls.txt")
  echo "$URLS" > "$SAVE_FILE"
  echo "Saved $(echo "$URLS" | wc -l) URLs to $SAVE_FILE"
  exit 0
fi

if [[ "$OPTIONS" == *"--open"* ]]; then
  # Detect OS and open URLs
  case "$(uname -s)" in
    Linux*)
      while read -r url; do
        xdg-open "$url" 2>/dev/null &
      done <<< "$URLS"
      ;;
    Darwin*)
      while read -r url; do
        open "$url" 2>/dev/null &
      done <<< "$URLS"
      ;;
    MING*|MSYS*|CYGWIN*)
      while read -r url; do
        start "$url" 2>/dev/null
      done <<< "$URLS"
      ;;
  esac
  echo "Opening $(echo "$URLS" | wc -l) URLs..."
  exit 0
fi

# Default: print URLs
echo "$URLS"