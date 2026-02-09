# Search CLI Examples

This directory contains example commands and use cases for the Search CLI tool.

## Basic Examples

### Simple Search

```bash
search "golang tutorials"
```

### Search with Multiple Words

```bash
search "machine learning with python"
```

### Search for Specific Terms

```bash
search "site:github.com golang cli"
```

## Result Count Examples

### Get More Results

```bash
search -n 20 "docker best practices"
```

### Get Fewer Results (Quick Scan)

```bash
search -n 5 "latest tech news"
```

### Maximum Results

```bash
search -n 100 "linux kernel"
```

## Output Format Examples

### JSON Output (for scripting)

```bash
search -f json "rust programming"
```

### JSON with jq (get just titles)

```bash
search -f json "golang" | jq -r '.results[].title'
```

### JSON with jq (get just URLs)

```bash
search -f json "kubernetes" | jq -r '.results[].url'
```

### JSON with jq (count by engine)

```bash
search -f json "docker" | jq '.results[] | .engine' | sort | uniq -c
```

### JSON with jq (filter by score)

```bash
search -f json "golang" | jq '.results[] | select(.score > 0.8)'
```

### Markdown Output (for documentation)

```bash
search -f markdown "microservices architecture"
```

### Save markdown to file

```bash
search -f markdown "api design" > results.md
```

### Plaintext Output (default)

```bash
search -f text "linux commands"
```

## Category Examples

### Image Search

```bash
search -c images "cute cats"
```

### Video Search

```bash
search -c videos "python tutorial"
```

### News Search

```bash
search -c news "artificial intelligence"
```

### IT/Computing Search

```bash
search -c it "kubernetes deployment"
```

### Science Search

```bash
search -c science "quantum computing"
```

### File Search

```bash
search -c files "docker compose"
```

### Music Search

```bash
search -c music "lofi beats"
```

## Time Filter Examples

### Results from Today

```bash
search --time day "breaking news"
```

### Results from This Week

```bash
search --time week "tech startup"
```

### Results from This Month

```bash
search --time month "javascript frameworks"
```

### Results from This Year

```bash
search --time year "machine learning trends"
```

## Instance Examples

### Use Default Instance

```bash
search "golang"
```

### Use Custom Instance

```bash
search -i https://searx.me "golang"
```

### Use Instance with Longer Timeout

```bash
search -i https://searx.example.com -t 60 "complex query"
```

## Language Examples

### Search in English

```bash
search -l en "python programming"
```

### Search in German

```bash
search -l de "python programmierung"
```

### Search in French

```bash
search -l fr "apprentissage automatique"
```

### Search in Spanish

```bash
search -l es "programación en python"
```

## Safe Search Examples

### Safe Search Off

```bash
search -s 0 "technical documentation"
```

### Safe Search Moderate (default)

```bash
search -s 1 "web search"
```

### Safe Search Strict

```bash
search -s 2 "family friendly content"
```

## Verbose Mode Examples

### Debug Query Issues

```bash
search -v "search query"
```

### See Which Config is Used

```bash
search -v "test"
```

### Debug Network Issues

```bash
search -v -t 60 "complex query"
```

## Browser Integration Examples

### Open First Result

```bash
search --open "github"
```

### Open Top 5 Results

```bash
search -n 5 --open-all "rust programming"
```

### Open Specific Category Results

```bash
search -c images --open-all "cute puppies"
```

## Config File Examples

### Use Custom Config

```bash
search --config ~/my-search-config.yaml "query"
```

### Test New Config

```bash
search --config /tmp/test-config.yaml "test"
```

## Combined Flag Examples

### JSON Output with Many Results

```bash
search -n 30 -f json "docker orchestration"
```

### Image Search with Time Filter

```bash
search -c images --time week "sunset photography"
```

### News Search with Verbose Output

```bash
search -c news -v "tech industry"
```

### Custom Instance with Custom Language

```bash
search -i https://searx.me -l de "linux"
```

### Strict Safe Search with Markdown Output

```bash
search -s 2 -f markdown "educational content"
```

## Scripting Examples

### Extract and Open URLs

```bash
search -f json "golang tutorial" | jq -r '.results[0].url' | xargs open
```

### Save Results to File

```bash
search -f json -n 50 "machine learning" > results.json
```

### Count Results

```bash
search -f json "docker" | jq '.results | length'
```

### Get Highest Scored Results

```bash
search -f json "kubernetes" | jq '.results | sort_by(.score) | reverse | .[0:5]'
```

### Search Multiple Queries

```bash
for query in "golang" "rust" "python"; do
  echo "=== $query ==="
  search -n 5 "$query"
done
```

### Parallel Searches

```bash
search "golang" & search "rust" & search "python" & wait
```

## Integration with Other Tools

### Pipe to fzf for Interactive Selection

```bash
search -f json "docker" | jq -r '.results[] | "\(.title) \t\(.url)"' | fzf
```

### Send Results to Less

```bash
search -n 50 "linux commands" | less
```

### Grep Through Results

```bash
search -f json "programming" | jq -r '.results[].content' | grep -i "tutorial"
```

### Create HTML from Markdown

```bash
search -f markdown "api design" | pandoc -f markdown -o results.html
```

## Troubleshooting Examples

### Test Connection

```bash
search -v "test query"
```

### Test Different Instance

```bash
search -i https://searx.me "test"
```

### Increase Timeout for Slow Connections

```bash
search -t 60 "complex query"
```

### Debug Config Issues

```bash
search -v --config ~/.search/config.yaml "test"
```

## Real-World Use Cases

### Research Paper Search

```bash
search -c science -n 20 "quantum entanglement recent papers"
```

### Programming Documentation

```bash
search -c it -f json "python async await" | jq -r '.results[0].url' | xargs open
```

### News Monitoring

```bash
search -c news --time day "artificial intelligence regulation"
```

### Image Assets for Design

```bash
search -c images -n 30 "minimalist logo design"
```

### Video Tutorial Discovery

```bash
search -c videos "react hooks tutorial"
```

### Technical Documentation Lookup

```bash
search -c it -s 0 "kubernetes ingress nginx"
```

### Competitor Analysis

```bash
search -n 30 "competitor company name"
```

### Market Research

```bash
search -c news --time month "industry trends 2024"
```

### Learning Path Planning

```bash
search "golang learning path" "rust learning path" "python learning path"
```

## Advanced Examples

### Batch Search from File

```bash
cat queries.txt | xargs -I {} search -n 5 "{}"
```

### Search and Email Results

```bash
search -f json "topic" | mail -s "Search Results" user@example.com
```

### Create Daily Search Report

```bash
#!/bin/bash
date > report.txt
echo "=== Tech News ===" >> report.txt
search -c news --time day "technology" >> report.txt
echo "=== AI Updates ===" >> report.txt
search -c news --time day "artificial intelligence" >> report.txt
cat report.txt | mail -s "Daily Search Report" user@example.com
```

### Search History Pattern

```bash
HISTTIMEFORMAT="%s " history | grep "search " | awk '{print $2}' | sort | uniq -c | sort -rn | head
```

## Tips and Tricks

### Use Quotes for Multi-word Queries

```bash
# Good
search "machine learning"

# Bad (searches for "machine" only)
search machine learning
```

### Combine with Shell Aliases

```bash
# Add to ~/.bashrc or ~/.zshrc
alias s='search'
alias si='search -i https://searx.me'
alias sj='search -f json'
alias sm='search -f markdown'
```

### Use Tab Completion

```bash
search --format <TAB>  # Shows: json, markdown, text
search --category <TAB>  # Shows all categories
```

### Quick Results Preview

```bash
search -n 3 "topic"  # Just get top 3 results
```

### Find Specific Domains

```bash
search "site:reddit.com golang tips"
search "site:stackoverflow.com kubernetes error"
```

## See Also

- [Configuration Guide](../docs/configuration.md)
- [Main README](../README.md)
- [SPEC.md](../SPEC.md) for detailed technical information