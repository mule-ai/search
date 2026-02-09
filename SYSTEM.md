# Search CLI - Agent Instructions

## Basic Usage

```bash
search "<query>"
```

Example:
```bash
search "golang tutorials"
```

## Output Formats

The tool supports three output formats:

### Text (Default)
Human-readable plain text with colored output:

```bash
search "machine learning"
```

### JSON
Machine-readable JSON format for parsing:

```bash
search -f json "rust programming"
```

JSON output structure:
```json
{
  "query": "rust programming",
  "results": [...],
  "infoboxes": [...],
  "answers": [],
  "suggestions": [],
  "total_results": 0,
  "metadata": {
    "instance": "https://search.butler.ooo",
    "search_time": "0.41s"
  }
}
```

### Markdown
Formatted markdown output:

```bash
search -f markdown "python async"
```

## Key Flags

### Results Control
- `-n, --results <int>` - Number of results to return (default: 10)
  ```bash
  search -n 20 "docker containers"
  ```

- `--page <int>` - Page number for pagination (default: 1)
  ```bash
  search --page 2 "kubernetes"
  ```

### Category Selection
- `-c, --category <string>` - Search category (default: "general")

Available categories:
- `general` - Default web search
- `images` - Image search
- `videos` - Video search
- `news` - News articles
- `map` - Map results
- `music` - Music search
- `it` - IT/tech results
- `science` - Scientific results
- `files` - File search

```bash
search -c images "cute cats"
search -c news "technology trends"
```

### Language
- `-l, --language <code>` - Language code (default: "en")
  ```bash
  search -l de "machine learning"
  ```

### Time Range
- `--time <range>` - Filter by time period: `day`, `week`, `month`, `year`
  ```bash
  search --time week "ai breakthrough"
  ```

### Instance Selection
- `-i, --instance <url>` - Use a specific SearXNG instance
  ```bash
  search -i https://searx.work "privacy tools"
  ```

Find more instances: https://searx.space/

### Safe Search
- `-s, --safe <level>` - Safe search level (0=off, 1=moderate, 2=strict)
  ```bash
  search -s 0 "medical research"
  ```

### Caching
- `--cache` - Enable caching (default: true)
- `--no-cache` - Disable caching for this request
- `--clear-cache` - Clear cache before searching
- `--cache-stats` - Show cache statistics

```bash
search --no-cache "latest news"
```

### Output Options
- `--no-color` - Disable colored output
- `--open` - Open first result in browser
- `--open-all` - Open all results in browser

### Performance & Debugging
- `-v, --verbose` - Enable verbose output
- `-t, --timeout <seconds>` - Request timeout (default: 30)
  ```bash
  search -v -t 60 "complex query"
  ```

## Understanding Results

### When No Results Are Found

The default SearXNG instance may have engines blocked by CAPTCHA. However, the tool will still display:

1. **Infoboxes** - Structured information from sources like Wikidata, Wikipedia
   - Includes descriptions, attributes, official links
   - Often provides comprehensive topic overviews

2. **Answers** - Direct answers from computational engines

3. **Suggestions** - Alternative search queries

Example output when engines are blocked:
```
go programming
==============

No results found.

## Go

programming language developed by Google and the open-source community

  • Inception: Tuesday, November 10, 2009
  • Developer: The Go Authors, Robert Griesemer, Rob Pike, Google

  Links:
  ★ Official website
      https://go.dev
  ...
```

### Infobox Components

Infoboxes contain:
- **Title/Name** - The subject name
- **Content** - Description or summary
- **Attributes** - Key-value pairs (founded date, CEO, version, etc.)
- **Links** - Official websites, Wikipedia, repositories, social media
  - Links marked with ★ are official/primary sources

## Best Practices for Agents

### 1. Use JSON for Programmatic Access
When you need to parse results:
```bash
search -f json "query" | jq '.results[] | .title, .url'
```

### 2. Adjust Result Count Based on Task
- For quick lookups: `search -n 5 "query"`
- For comprehensive research: `search -n 30 "query"`

### 3. Use Time Filters for Current Events
```bash
search --time day "breaking news"
search --time week "tech updates"
```

### 4. Leverage Categories
```bash
search -c news "election"
search -c images "sunset"
search -c videos "tutorial"
```

### 5. Handle Rate Limiting
If the default instance is blocked:
- Use `--no-cache` to bypass cached empty results
- Try a different instance: `search -i https://searx.work "query"`
- Check https://searx.space/ for available instances

### 6. Combine with Pipes
```bash
# Get URLs only
search -f json "golang" | jq -r '.results[].url'

# Count results
search -f json "topic" | jq '.total_results'

# Extract infobox data
search -f json "company" | jq '.infoboxes[].attributes'
```

## Common Workflows

### Research a Topic
```bash
# Get comprehensive results with infoboxes
search -n 20 "artificial intelligence"

# Get recent news
search -c news --time week "ai regulation"

# Find official documentation
search "golang official documentation"
```

### Find Resources
```bash
# Tutorials
search "golang tutorial beginner"

# Libraries/Packages
search "python http client library"

# Examples
search "react hooks examples"
```

### Quick Facts
```bash
# The infobox often provides structured facts
search "python programming language"
search "linux kernel"
```

### Troubleshooting
```bash
# Error messages
search "docker error container already exists"

# Version info
search "go 1.21 release notes"

# Documentation
search "kubernetes pod networking"
```

## Performance Tips

1. **Use caching** - Repeated queries are faster with cache enabled (default)
2. **Adjust timeout** - Increase timeout for complex queries: `-t 60`
3. **Limit results** - Use `-n` to reduce bandwidth for quick lookups
4. **Choose specific categories** - Faster than general search when you know the type

## Troubleshooting

### "No results found" but infoboxes appear
This is normal when search engines are rate-limited. The infoboxes still provide valuable structured information.

### Slow responses
- Try a different instance: `-i https://searx.work`
- Reduce timeout: `-t 15`
- Reduce result count: `-n 5`

### Instance unavailable
Find alternative instances at https://searx.space/ and use with `-i`

## Advanced Usage

### Shell Integration
```bash
# Fuzzy search with fzf
search -f json "topic" | jq -r '.results[].title' | fzf

# Open results in browser
search --open "search query"
```

### Batch Queries
```bash
for term in "golang" "rust" "python"; do
    echo "### $term"
    search -n 5 "$term"
done
```

### Logging Searches
```bash
search -v "query" 2>&1 | tee search_log.txt
```

## Summary

- **Binary:** `/bin/search`
- **Default format:** Text (use `-f json` for parsing)
- **Default instance:** `https://search.butler.ooo`
- **Infoboxes display** even when main results are empty
- **Categories** available for targeted searches
- **JSON output** recommended for agent integration
- **Multiple instances** available if default is rate-limited

For more help:
```bash
search --help
search help        # Show available commands
search version     # Show version info
```