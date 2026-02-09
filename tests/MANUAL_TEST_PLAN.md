# Manual Test Plan for Search CLI

This document outlines the manual testing procedures for the search CLI tool.

## Test Environment

- **Binary**: `/data/jbutler/git/mule-ai/search/search`
- **SearXNG Instance**: https://search.butler.ooo
- **Test Date**: 2026-02-09
- **Config Location**: ~/.search/config.yaml

## Test Categories

### 1. Basic Functionality Tests

#### Test 1.1: Basic Search
```bash
./search "golang tutorials"
```
**Expected**: Returns results in plaintext format with numbered list

#### Test 1.2: Search with Result Limit
```bash
./search -n 5 "rust programming"
```
**Expected**: Returns exactly 5 results

#### Test 1.3: Help Command
```bash
./search --help
```
**Expected**: Displays comprehensive help with all flags

#### Test 1.4: Version Command
```bash
./search --version
```
**Expected**: Displays version information

### 2. Output Format Tests

#### Test 2.1: JSON Output
```bash
./search -f json "docker containers"
```
**Expected**: Valid JSON with query, results array, and metadata

#### Test 2.2: Markdown Output
```bash
./search -f markdown "kubernetes"
```
**Expected**: Markdown formatted output with headers and links

#### Test 2.3: Plaintext Output (default)
```bash
./search -f text "linux commands"
```
**Expected**: Plaintext numbered format

### 3. Configuration Tests

#### Test 3.1: Custom Instance
```bash
./search -i https://search.butler.ooo "test query"
```
**Expected**: Queries specified instance

#### Test 3.2: Custom Config File
```bash
./search --config /tmp/test-config.yaml "test"
```
**Expected**: Uses specified config file

#### Test 3.3: Environment Variables
```bash
export_SEARCH_INSTANCE=https://search.butler.ooo
export_SEARCH_RESULTS=15
./search "env test"
```
**Expected**: Uses environment variable values

### 4. Search Options Tests

#### Test 4.1: Category Search - Images
```bash
./search -c images "mountains"
```
**Expected**: Returns image results with img_src fields

#### Test 4.2: Category Search - Videos
```bash
./search -c videos "tutorial"
```
**Expected**: Returns video results

#### Test 4.3: Category Search - News
```bash
./search -c news "technology"
```
**Expected**: Returns news results

#### Test 4.4: Time Range Filter
```bash
./search --time day "breaking news"
```
**Expected**: Filters results to last day

#### Test 4.5: Language Filter
```bash
./search -l de "deutsch"
```
**Expected**: Searches with German language preference

#### Test 4.6: Safe Search Levels
```bash
./search -s 0 "test query"    # Off
./search -s 1 "test query"    # Moderate (default)
./search -s 2 "test query"    # Strict
```
**Expected**: Respects safe search setting

### 5. Pagination Tests

#### Test 5.1: Page 2 Results
```bash
./search --page 2 "golang"
```
**Expected**: Returns second page of results

#### Test 5.2: Combined Pagination and Limit
```bash
./search --page 3 -n 5 "python"
```
**Expected**: Returns 5 results from page 3

### 6. Browser Integration Tests

#### Test 6.1: Open First Result
```bash
./search --open "golang"  # Note: This will try to open a browser
```
**Expected**: Opens first result in default browser

### 7. Error Handling Tests

#### Test 7.1: Empty Query
```bash
./search ""
```
**Expected**: Error message about empty query

#### Test 7.2: Invalid Instance URL
```bash
./search -i https://invalid.example.com "test"
```
**Expected**: Network error with helpful message

#### Test 7.3: Invalid Result Count
```bash
./search -n 500 "test"
```
**Expected**: Validation error about max result count

#### Test 7.4: Invalid Format
```bash
./search -f invalid "test"
```
**Expected**: Error listing valid formats

#### Test 7.5: Invalid Timeout
```bash
./search -t 500 "test"
```
**Expected**: Validation error about max timeout

### 8. Verbose Mode Tests

#### Test 8.1: Verbose Output
```bash
./search -v "test query"
```
**Expected**: Shows detailed information about request, timing, etc.

### 9. Shell Completion Tests

#### Test 9.1: Bash Completion
```bash
./search completion bash > /tmp/search_completion.bash
source /tmp/search_completion.bash
./search <TAB><TAB>
```
**Expected**: Shows subcommands and flags

#### Test 9.2: Zsh Completion
```bash
./search completion zsh > /tmp/_search
```
**Expected**: Generates valid zsh completion script

#### Test 9.3: Fish Completion
```bash
./search completion fish > /tmp/search.fish
```
**Expected**: Generates valid fish completion script

### 10. Colored Output Tests

#### Test 10.1: Auto Color Detection
```bash
./search "color test"
```
**Expected**: Uses colors when output is to TTY

#### Test 10.2: Force No Color
```bash
./search --no-color "color test"
```
**Expected**: No color codes in output

### 11. API Key Tests

#### Test 11.1: API Key via CLI Flag
```bash
./search --api-key "test-key" "query"
```
**Expected**: Sends Authorization header with Bearer token

#### Test 11.2: API Key via Config
```bash
# Add api_key to config.yaml
./search "query"
```
**Expected**: Uses API key from config

### 12. Edge Case Tests

#### Test 12.1: Special Characters in Query
```bash
./search "c++ templates"
./search "node.js async"
./search "golang & rust"
```
**Expected**: Properly handles special characters

#### Test 12.2: Unicode Query
```bash
./search "日本語"
./search "emoji 😊"
```
**Expected**: Properly handles Unicode characters

#### Test 12.3: Long Query
```bash
./search "this is a very long search query with many words to test if the system handles it correctly"
```
**Expected**: Processes long query without issues

#### Test 12.4: No Results Scenario
```bash
./search "asdfghjklzxcvbnm1234567890unlikelytohavematches"
```
**Expected**: Returns empty results with appropriate message

### 13. Integration Tests

#### Test 13.1: Pipe to jq
```bash
./search -f json "test" | jq '.results | length'
```
**Expected**: Valid JSON that can be parsed by jq

#### Test 13.2: Pipe to grep
```bash
./search "test" | grep -i "test"
```
**Expected**: Plaintext output is grep-able

#### Test 13.3: Output Redirection
```bash
./search "test" > /tmp/search_output.txt
cat /tmp/search_output.txt
```
**Expected**: Output is correctly written to file

### 14. Performance Tests

#### Test 14.1: Large Result Set
```bash
time ./search -n 100 "popular topic"
```
**Expected**: Completes in reasonable time (<5 seconds)

#### Test 14.2: Timeout Test
```bash
time ./search -t 5 "test"
```
**Expected**: Respects timeout setting

## Test Results Template

For each test, record:
- ✅ PASS - Test completed successfully
- ❌ FAIL - Test failed with description
- ⚠️  PARTIAL - Test partially passed with notes
- ⏭️  SKIP - Test was skipped

---

## Summary

Total Tests: XX
Passed: XX
Failed: XX
Skipped: XX

### Issues Found
1. [List any issues discovered during testing]

### Recommendations
1. [Any recommendations for improvements]
