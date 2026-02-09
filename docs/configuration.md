# Configuration Guide

This guide covers all aspects of configuring the Search CLI tool.

## Config File Location

The global configuration file is located at:

```
~/.search/config.yaml
```

## Config File Structure

```yaml
# SearXNG instance URL (required)
instance: "https://search.butler.ooo"

# Default number of results to return (default: 10)
results: 10

# Default output format: json, markdown, or text (default: text)
format: "text"

# Optional: API key if instance requires authentication
api_key: ""

# Optional: Request timeout in seconds (default: 30)
timeout: 30

# Optional: Specific categories to search (default: all)
# Options: general, images, videos, news, map, music, it, science, files, etc.
categories:
  - "general"

# Optional: Language preference (default: en)
language: "en"

# Optional: Safe search setting (default: moderate)
# Options: 0 (off), 1 (moderate), 2 (strict)
safe_search: 1
```

## Configuration Precedence

Settings are applied in the following order (highest to lowest priority):

1. **CLI Flags** - Command-line arguments override everything
2. **Environment Variables** - Environment variables override config file
3. **Config File** - `~/.search/config.yaml` or custom path
4. **Defaults** - Built-in default values

### Example Precedence

```bash
# Default from config: results=10
search "query"                    # Returns 10 results

# Environment variable overrides config
export SEARCH_RESULTS=20
search "query"                    # Returns 20 results

# CLI flag overrides everything
search -n 50 "query"              # Returns 50 results
```

## Environment Variables

All configuration options can be set via environment variables:

| Environment Variable | Config Key | Description | Default |
|---------------------|------------|-------------|---------|
| `SEARCH_INSTANCE` | `instance` | SearXNG instance URL | https://search.butler.ooo |
| `SEARCH_RESULTS` | `results` | Number of results | 10 |
| `SEARCH_FORMAT` | `format` | Output format | text |
| `SEARCH_API_KEY` | `api_key` | API key for auth | empty |
| `SEARCH_TIMEOUT` | `timeout` | Request timeout (seconds) | 30 |
| `SEARCH_LANGUAGE` | `language` | Language code | en |
| `SEARCH_SAFE` | `safe_search` | Safe search level | 1 |

### Using Environment Variables

```bash
# Set instance URL
export SEARCH_INSTANCE="https://searx.me"

# Set number of results
export SEARCH_RESULTS="25"

# Set output format
export SEARCH_FORMAT="json"

# Now use the CLI
search "golang tutorials"
```

### Environment Variables in Files

For persistent settings, add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
# ~/.bashrc or ~/.zshrc
export SEARCH_INSTANCE="https://search.butler.ooo"
export SEARCH_RESULTS="15"
export SEARCH_FORMAT="text"
export SEARCH_LANGUAGE="en"
```

## Creating a Custom Config File

You can specify a custom config file location using the `--config` flag:

```bash
search --config /path/to/custom-config.yaml "query"
```

This is useful for:

- Using different instances for different projects
- Testing configurations without modifying the global config
- Team-specific configurations

## Default Config Generation

If the config file doesn't exist, Search CLI will create one with default values:

```yaml
instance: https://search.butler.ooo
results: 10
format: text
api_key: ""
timeout: 30
language: en
safe_search: 1
```

## Config Validation

Search CLI validates configuration on startup:

- **instance**: Must be a valid URL
- **results**: Must be between 1 and 100
- **format**: Must be one of: `json`, `markdown`, `text`
- **timeout**: Must be between 1 and 300 seconds
- **safe_search**: Must be 0, 1, or 2

### Invalid Config Handling

If your config is invalid, Search CLI will:

1. Display an error message describing the issue
2. Show the invalid value and location
3. Suggest how to fix it
4. Fall back to defaults for that setting

Example error:

```
Error: Invalid config file ~/.search/config.yaml:
  Line 5: results must be between 1 and 100, got 150
```

## Common Configurations

### Privacy-Focused Config

```yaml
instance: "https://search.butler.ooo"
results: 20
format: "text"
timeout: 30
language: "en"
safe_search: 2  # Strict
```

### Developer Config

```yaml
instance: "https://search.butler.ooo"
results: 15
format: "json"  # Easy to parse with jq
timeout: 60     # Longer timeout for complex queries
language: "en"
safe_search: 0  # Off for technical searches
```

### Quick Search Config

```yaml
instance: "https://search.butler.ooo"
results: 5      # Fewer results for faster scanning
format: "text"
timeout: 15
language: "en"
safe_search: 1
```

### Research Config

```yaml
instance: "https://search.butler.ooo"
results: 50     # More results for comprehensive research
format: "markdown"  # Easy to read and document
timeout: 60
language: "en"
safe_search: 1
```

## Using Multiple Instances

You can create multiple config files for different SearXNG instances:

```bash
# Default config (~/.search/config.yaml)
instance: "https://search.butler.ooo"

# Work config (~/.search-work.yaml)
instance: "https://searx.work.example.com"
api_key: "work-api-key"

# Personal config (~/.search-personal.yaml)
instance: "https://searxng.privacy.example.com"
```

Usage:

```bash
# Use default config
search "work stuff"

# Use work config
search --config ~/.search-work.yaml "work stuff"

# Use personal config
search --config ~/.search-personal.yaml "personal stuff"
```

## Config File Tips

### Comments in YAML

You can add comments to your config file using `#`:

```yaml
# My SearXNG instance
instance: "https://search.butler.ooo"

# Number of results (1-100)
results: 20

# Output format: json, markdown, or text
format: "text"
```

### Testing Your Config

Use verbose mode to see which config values are being used:

```bash
search -v "test query"
```

Output will show:

```
Config loaded from: /home/user/.search/config.yaml
Using instance: https://search.butler.ooo
Results: 20
Format: text
Language: en
Timeout: 30s
```

### Reloading Config

The config is loaded each time you run `search`. Simply edit your config file and the next command will use the new values. No need to restart anything.

## Troubleshooting Config Issues

### Config Not Being Used

1. Check if a custom config is being set: `search --help`
2. Look for environment variables that might override: `env | grep SEARCH`
3. Use verbose mode to see which config is loaded: `search -v "test"`

### Permission Errors

If you get permission errors:

```bash
# Fix permissions
chmod 600 ~/.search/config.yaml
```

### YAML Syntax Errors

Common YAML mistakes:

```yaml
# WRONG: No space after colon
instance:"https://search.butler.ooo"

# RIGHT: Space after colon
instance: "https://search.butler.ooo"

# WRONG: Using = for assignment
instance = "https://search.butler.ooo"

# WRONG: Quotes around keys
"instance": "https://search.butler.ooo"
```

Use a YAML validator to check your config:

```bash
# Using Python
python3 -c "import yaml; yaml.safe_load(open('~/.search/config.yaml'))"

# Using Ruby
ruby -ryaml -e "puts YAML.load_file('~/.search/config.yaml').inspect"
```

## Example: Full Config File

```yaml
# ~/.search/config.yaml
# Search CLI Configuration File

# SearXNG instance to use
instance: "https://search.butler.ooo"

# Number of search results to return (1-100)
results: 20

# Output format: json, markdown, or text
format: "text"

# API key (if instance requires authentication)
api_key: ""

# Request timeout in seconds (1-300)
timeout: 30

# Language code (en, de, fr, es, etc.)
language: "en"

# Safe search: 0=off, 1=moderate, 2=strict
safe_search: 1

# Optional: Define specific categories
# categories:
#   - general
#   - videos
#   - images
```

## Related Documentation

- [Installation Guide](../README.md#installation)
- [Usage Examples](../examples/README.md)
- [Output Formats](../README.md#output-formats)
- [Troubleshooting](../README.md#troubleshooting)