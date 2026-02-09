# Manual Testing Plan for Search CLI v1.0.0

This document outlines the manual testing steps to verify all features work correctly before release.

## Prerequisites

1. Built binary available: `make build`
2. Test SearXNG instance accessible: https://search.butler.ooo
3. Clean config state: `rm -f ~/.search/config.yaml`

## Test Categories

### 1. Basic Functionality Tests

#### 1.1 No Config - Should use defaults
```bash
# Remove config to test defaults
rm -f ~/.search/config.yaml
./bin/search "golang tutorials"
# Expected: Should work with default instance and settings
```

#### 1.2 Basic Search
```bash
./bin/search "golang"
# Expected: Return results in plaintext format
```

#### 1.3 Query with Special Characters
```bash
./bin/search "golang OR rust"
./bin/search "site:github.com golang"
# Expected: Handle special query syntax correctly
```

### 2. Output Format Tests

#### 2.1 Plaintext Output (default)
```bash
./bin/search -f text "kubernetes"
# Expected: Numbered list format with URLs
```

#### 2.2 JSON Output
```bash
./bin/search -f json "docker"
# Expected: Valid JSON with results array
./bin/search -f json "docker" | jq '.results | length'
# Expected: Count of results
```

#### 2.3 Markdown Output
```bash
./bin/search -f markdown "terraform"
# Expected: Markdown formatted with # headings
```

### 3. Result Count Tests

#### 3.1 Default Results (10)
```bash
./bin/search "rust" | grep -c "^\[" || true
# Expected: 10 results
```

#### 3.2 Custom Result Count
```bash
./bin/search -n 5 "python" | grep -c "^\[" || true
# Expected: 5 results
```

#### 3.3 Maximum Results (100)
```bash
./bin/search -n 100 "java"
# Expected: Up to 100 results
```

#### 3.4 Invalid Result Count
```bash
./bin/search -n 0 "test"
# Expected: Error about invalid range
./bin/search -n 101 "test"
# Expected: Error about invalid range
```

### 4. Instance Configuration Tests

#### 4.1 Default Instance
```bash
./bin/search "search engine"
# Expected: Use https://search.butler.ooo
```

#### 4.2 Custom Instance
```bash
./bin/search -i https://searx.work "test"
# Expected: Use specified instance
```

#### 4.3 Invalid Instance URL
```bash
./bin/search -i "not-a-url" "test"
# Expected: Error about invalid URL
```

### 5. Category Tests

#### 5.1 General Category (default)
```bash
./bin/search -c general "news"
# Expected: General web results
```

#### 5.2 Images Category
```bash
./bin/search -c images "cats"
# Expected: Image results with img_src fields
```

#### 5.3 Videos Category
```bash
./bin/search -c videos "music"
# Expected: Video results
```

#### 5.4 News Category
```bash
./bin/search -c news "technology"
# Expected: News results
```

#### 5.5 Invalid Category
```bash
./bin/search -c invalid "test"
# Expected: Error about invalid category
```

### 6. Language Tests

#### 6.1 Default Language (en)
```bash
./bin/search "test"
# Expected: English results
```

#### 6.2 Custom Language
```bash
./bin/search -l de "test"
# Expected: German-focused results
./bin/search -l fr "test"
# Expected: French-focused results
```

### 7. Safe Search Tests

#### 7.1 Moderate (default)
```bash
./bin/search "test"
# Expected: Moderate filtering
```

#### 7.2 Safe Search Off
```bash
./bin/search -s 0 "test"
# Expected: No filtering
```

#### 7.3 Strict Safe Search
```bash
./bin/search -s 2 "test"
# Expected: Strict filtering
```

#### 7.4 Invalid Safe Search Level
```bash
./bin/search -s 3 "test"
# Expected: Error about invalid level
```

### 8. Timeout Tests

#### 8.1 Default Timeout (30s)
```bash
./bin/search "test"
# Expected: Completes within 30s
```

#### 8.2 Custom Timeout
```bash
./bin/search -t 10 "test"
# Expected: Use 10s timeout
```

#### 8.3 Timeout Too Short
```bash
./bin/search -t 1 "test"
# Expected: May timeout with error message
```

### 9. Pagination Tests

#### 9.1 First Page (default)
```bash
./bin/search "golang"
# Expected: Page 1 results
```

#### 9.2 Second Page
```bash
./bin/search --page 2 "golang"
# Expected: Different results from page 1
```

#### 9.3 Page Number in JSON
```bash
./bin/search -f json --page 2 "test" | jq '.page'
# Expected: 2
```

### 10. Time Range Tests

#### 10.1 No Time Filter (default)
```bash
./bin/search "ai"
# Expected: All-time results
```

#### 10.2 Past Day
```bash
./bin/search --time day "ai"
# Expected: Recent results
```

#### 10.3 Past Week
```bash
./bin/search --time week "ai"
# Expected: Results from past week
```

#### 10.4 Past Month
```bash
./bin/search --time month "ai"
# Expected: Results from past month
```

#### 10.5 Past Year
```bash
./bin/search --time year "ai"
# Expected: Results from past year
```

### 11. Verbose Mode Tests

#### 11.1 Verbose Output
```bash
./bin/search -v "test"
# Expected: Additional debug information
```

#### 11.2 Verbose with JSON
```bash
./bin/search -v -f json "test"
# Expected: Verbose logs before JSON
```

### 12. Color Output Tests

#### 12.1 Auto Color (TTY detection)
```bash
./bin/search "test" | cat
# Expected: No color codes (piped)
./bin/search "test"
# Expected: Color codes if TTY
```

#### 12.2 Force No Color
```bash
./bin/search --no-color "test"
# Expected: No color codes
```

### 13. Browser Integration Tests

#### 13.1 Open First Result
```bash
# Note: May fail if no browser configured
./bin/search --open "golang"
# Expected: Opens browser (if available)
```

#### 13.2 Open All Results
```bash
./bin/search -n 3 --open-all "rust"
# Expected: Opens all results (browser dependent)
```

### 14. Config File Tests

#### 14.1 Create Default Config
```bash
rm -f ~/.search/config.yaml
./bin/search "test"
# Expected: Creates ~/.search/config.yaml with defaults
cat ~/.search/config.yaml
# Expected: YAML with instance, results, format, etc.
```

#### 14.2 Custom Config File
```bash
./bin/search --config /tmp/test-config.yaml "test"
# Expected: Use specified config file
```

#### 14.3 Config Priority
```bash
# Set default in config
echo 'format: "json"' > ~/.search/config.yaml
./bin/search "test"
# Expected: JSON output from config

# Override with flag
./bin/search -f text "test"
# Expected: Text output (flag overrides config)
```

### 15. Environment Variable Tests

#### 15.1 Instance from Env
```bash
export SEARCH_INSTANCE="https://search.butler.ooo"
./bin/search "test"
# Expected: Use env instance
```

#### 15.2 Results from Env
```bash
export SEARCH_RESULTS="5"
./bin/search "test"
# Expected: 5 results
```

#### 15.3 Format from Env
```bash
export SEARCH_FORMAT="json"
./bin/search "test"
# Expected: JSON output
```

#### 15.4 Priority: Env vs Flag
```bash
export SEARCH_FORMAT="json"
./bin/search -f text "test"
# Expected: Text output (flag overrides env)
```

### 16. Error Handling Tests

#### 16.1 Empty Query
```bash
./bin/search ""
# Expected: Error about empty query
```

#### 16.2 No Query Argument
```bash
./bin/search
# Expected: Error about missing query
```

#### 16.3 Network Error (Invalid Instance)
```bash
./bin/search -i https://this-domain-does-not-exist-12345.com "test"
# Expected: Network error with helpful message
```

#### 16.4 Invalid Config File
```bash
echo "invalid: yaml: content: [" > ~/.search/config.yaml
./bin/search "test"
# Expected: Error about invalid config YAML
```

### 17. Help and Version Tests

#### 17.1 Help Flag
```bash
./bin/search --help
# Expected: Full help output with all flags
./bin/search -h
# Expected: Same as --help
```

#### 17.2 Version Flag
```bash
./bin/search --version
# Expected: Version number
./bin/search -V
# Expected: Same as --version
```

### 18. Shell Completion Tests

#### 18.1 Bash Completion
```bash
./bin/search completion bash
# Expected: Bash completion script
```

#### 18.2 Zsh Completion
```bash
./bin/search completion zsh
# Expected: Zsh completion script
```

#### 18.3 Fish Completion
```bash
./bin/search completion fish
# Expected: Fish completion script
```

### 19. Edge Cases

#### 19.1 Very Long Query
```bash
./bin/search "$(python3 -c 'print("test " * 100)')
# Expected: Handle long query
```

#### 19.2 Unicode Characters
```bash
./bin/search "日本語"
./bin/search "🔍 emoji"
# Expected: Handle unicode correctly
```

#### 19.3 URL in Query
```bash
./bin/search "https://github.com/torvalds/linux"
# Expected: Search for the URL
```

#### 19.4 Query with Quotes
```bash
./bin.search '"exact phrase"'
# Expected: Search for exact phrase
```

### 20. Performance Tests

#### 20.1 Large Result Set
```bash
time ./bin/search -n 100 "programming"
# Expected: Completes in reasonable time
```

#### 20.2 Multiple Searches (Sequential)
```bash
time for i in {1..5}; do ./bin/search "test $i" > /dev/null; done
# Expected: Each search completes quickly
```

## Test Results Template

For each test, record:
- [ ] PASS - Test passed
- [ ] FAIL - Test failed (note issue)
- [ ] SKIP - Test skipped (reason)

## Issues Found

Document any issues discovered during testing:

1. **Issue**: Description
   - Test: Test name
   - Expected: What should happen
   - Actual: What actually happened
   - Severity: Critical/High/Medium/Low

## Sign-off

- Tester: _______________
- Date: _______________
- Version: _______________
- Ready for Release: [ ] YES [ ] NO
