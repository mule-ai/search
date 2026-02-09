# Search CLI

A powerful command-line search tool written in Go that queries SearXNG instances and returns formatted results.

[![Go Reference](https://pkg.go.dev/badge/github.com/mule-ai/search.svg)](https://pkg.go.dev/github.com/mule-ai/search)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Stable Release v1.0.0** - Search CLI is now production-ready! [View Changelog](CHANGELOG.md)

## Features

- 🔍 **Fast & Private**: Uses SearXNG as the backend for private, metasearch
- 🎨 **Multiple Output Formats**: JSON, Markdown, and Plaintext output
- ⚙️ **Configurable**: Supports config files, environment variables, and CLI flags
- 🌐 **Multi-Instance**: Use any SearXNG instance
- 📝 **Colored Output**: Optional colored terminal output for better readability
- 🔧 **Shell Completion**: Bash, Zsh, and Fish completion support

## Installation

### Homebrew (macOS/Linux)

```bash
brew install mule-ai/tap/search
```

### Scoop (Windows)

```powershell
scoop bucket add mule-ai https://github.com/mule-ai/scoop-bucket
scoop install search
```

### Arch Linux (AUR)

```bash
paru -S search-bin
# or
yay -S search-bin
```

### From Binary

Download the latest release binary from the [releases page](https://github.com/mule-ai/search/releases).

### From Source

```bash
git clone https://github.com/mule-ai/search
cd search
make install
```

### Using Go

```bash
go install github.com/mule-ai/search/cmd/search@latest
```

## Quick Start

```bash
# Basic search
search "golang tutorials"

# Specify number of results
search -n 20 "machine learning"

# JSON output for scripting
search -f json "rust programming" | jq '.results[] | .title'

# Markdown output
search -f markdown "kubernetes best practices"

# Use specific category
search -c images "cute cats"

# Verbose mode
search -v "search query"
```

## Configuration

Search CLI uses a configuration file located at `~/.search/config.yaml`:

```yaml
# SearXNG instance URL
instance: "https://search.butler.ooo"

# Number of results to return
results: 10

# Output format: json, markdown, or text
format: "text"

# API key (if instance requires authentication)
api_key: ""

# Request timeout in seconds
timeout: 30

# Language preference
language: "en"

# Safe search: 0 (off), 1 (moderate), 2 (strict)
safe_search: 1
```

### Configuration Precedence

CLI flags > Environment variables > Config file > Defaults

### Environment Variables

```bash
export SEARCH_INSTANCE="https://search.butler.ooo"
export SEARCH_RESULTS="20"
export SEARCH_FORMAT="json"
export SEARCH_TIMEOUT="30"
export SEARCH_LANGUAGE="en"
export SEARCH_SAFE="1"
```

## Usage

### Basic Syntax

```bash
search [flags] <query>
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--instance` | `-i` | SearXNG instance URL | From config |
| `--results` | `-n` | Number of results (1-100) | 10 |
| `--format` | `-f` | Output format | text |
| `--category` | `-c` | Search category | general |
| `--timeout` | `-t` | Timeout in seconds | 30 |
| `--language` | `-l` | Language code | en |
| `--safe` | `-s` | Safe search level (0-2) | 1 |
| `--page` | | Page number | 1 |
| `--time` | | Time filter (day/week/month/year) | |
| `--config` | | Custom config file path | ~/.search/config.yaml |
| `--verbose` | `-v` | Enable verbose output | false |
| `--no-color` | | Disable colored output | false |
| `--open` | | Open first result in browser | false |
| `--open-all` | | Open all results in browser | false |
| `--help` | `-h` | Show help | |
| `--version` | `-V` | Show version | |

### Search Categories

- `general` - General web search
- `images` - Image search
- `videos` - Video search
- `news` - News search
- `map` - Map search
- `music` - Music search
- `it` - IT/Computing
- `science` - Science
- `files` - File search

### Output Formats

#### JSON Format

```json
{
  "query": "golang tutorials",
  "total_results": 1250000,
  "results": [
    {
      "title": "A Tour of Go",
      "url": "https://go.dev/tour/",
      "content": "Welcome to a tour of the Go programming language...",
      "engine": "google",
      "category": "general",
      "score": 0.95
    }
  ],
  "metadata": {
    "search_time": "0.24s",
    "instance": "https://search.butler.ooo"
  }
}
```

#### Markdown Format

```markdown
# Search Results: golang tutorials

Found **1,250,000** results in 0.24s

## [A Tour of Go](https://go.dev/tour/)
**Source:** Google | **Score:** 0.95

Welcome to a tour of the Go programming language...
```

#### Plaintext Format

```
golang tutorials
==============

[1] A Tour of Go
    https://go.dev/tour/
    Source: Google | Score: 0.95

    Welcome to a tour of the Go programming language...
```

## Examples

### Search with specific number of results

```bash
search -n 20 "docker best practices"
```

### Search in specific category

```bash
search -c images "mountain landscapes"
search -c videos "cats funny"
```

### Filter by time range

```bash
search --time day "latest tech news"
search --time week "ai developments"
```

### JSON output for scripting

```bash
# Get just URLs from results
search -f json "golang" | jq -r '.results[].url'

# Count results by engine
search -f json "golang" | jq '.results[] | .engine' | sort | uniq -c
```

### Open results in browser

```bash
# Open first result
search --open "github"

# Open all results
search -n 5 --open-all "rust programming"
```

## Shell Completion

Generate completion scripts:

```bash
# Bash
search completion bash > /etc/bash_completion.d/search

# Zsh
search completion zsh > /usr/local/share/zsh/site-functions/_search

# Fish
search completion fish > ~/.config/fish/completions/search.fish
```

## Development

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mule-ai/search
cd search

# Standard build
make build

# Optimized build (smaller binary)
make optimize

# Build and compress with UPX (smallest binary)
make optimize-upx

# Run tests
make test

# Install
make install
```

### Binary Optimization

The project supports multiple build optimization levels:

| Build Type | Size | Command |
|------------|------|---------|
| Standard | ~15 MB | `make build` |
| Optimized | ~10 MB | `make optimize` |
| UPX Compressed | ~5 MB | `make optimize-upx` |

For detailed optimization analysis, use `make optimize-full` or run `./scripts/optimize.sh` directly.

See [docs/optimization.md](docs/optimization.md) for more details on binary optimization.

### Project Structure

```
search/
├── cmd/
│   └── search/
│       └── main.go          # CLI entry point
├── internal/
│   ├── cli/                 # CLI commands and flags
│   ├── config/              # Configuration management
│   ├── formatter/           # Output formatters
│   └── searxng/             # SearXNG API client
├── pkg/
│   └── version/             # Version information
├── go.mod
├── go.sum
├── Makefile
├── SPEC.md
└── README.md
```

## Troubleshooting

### Connection Errors

If you experience connection errors:

1. Check your internet connection
2. Verify the SearXNG instance URL: `search -i https://searx.me "test"`
3. Increase timeout: `search -t 60 "query"`

### No Results

If you get no results:

1. Try a different query
2. Check if the instance is working: `curl https://search.butler.ooo/search?q=test&format=json`
3. Try different categories: `search -c videos "query"`

### Config Issues

If you have config issues:

1. Check your config at `~/.search/config.yaml`
2. Use verbose mode: `search -v "query"`
3. Try with a custom config: `search --config /path/to/config.yaml "query"`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Performance & Benchmarks

The search CLI is optimized for performance. Key optimizations include:

- **Optimized JSON Parsing**: Custom unmarshalers are ~26% faster than standard library
- **Buffer Pooling**: Reduces allocations for JSON decoding and formatting
- **Optimized Formatters**: 32-66% faster formatting with fewer allocations
- **Efficient String Building**: Uses `strings.Builder` throughout
- **Streaming Support**: Memory-efficient processing of large responses

For detailed information about JSON parsing optimizations, see [docs/json-parsing-optimizations.md](docs/json-parsing-optimizations.md).

### Running Benchmarks

```bash
# Run all performance benchmarks
make benchmark

# Run specific package benchmarks
go test -tags=benchmark -bench=. -benchmem ./internal/searxng/...
go test -tags=benchmark -bench=. -benchmem ./internal/formatter/...
```

### Performance Profiles

Generate CPU and memory profiles to analyze performance:

```bash
make profile
go tool pprof -http=:8080 profiles/cpu.prof
```

For detailed performance information, see [PERFORMANCE.md](PERFORMANCE.md).

## License

MIT License - see LICENSE file for details

## Default SearXNG Instance

This tool uses https://search.butler.ooo as the default SearXNG instance. You can change this in your config or via the `--instance` flag.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
- Powered by [SearXNG](https://searxng.org/)
