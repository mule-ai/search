# Development Workflow

This document describes the development workflow and processes for the Search CLI project.

## Architecture Overview

### Project Structure

```
search/
├── cmd/
│   └── search/
│       └── main.go              # CLI entry point
├── internal/
│   ├── config/
│   │   ├── config.go            # Config loading/parsing
│   │   └── config_test.go
│   ├── searxng/
│   │   ├── client.go            # SearXNG API client
│   │   ├── types.go             # Data structures
│   │   └── client_test.go
│   ├── formatter/
│   │   ├── formatter.go         # Output formatter interface
│   │   ├── json.go
│   │   ├── markdown.go
│   │   ├── text.go
│   │   └── formatter_test.go
│   └── cli/
│       ├── root.go              # Root command
│       ├── flags.go             # Flag definitions
│       └── completion.go        # Shell completion
├── pkg/
│   └── version/
│       └── version.go           # Version info
├── docs/                        # Documentation
├── examples/                    # Example scripts
├── tests/                       # Integration tests
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── SPEC.md
```

### Package Responsibilities

#### `internal/config`
- **Responsibility**: Configuration file loading, parsing, and validation
- **Key Functions**: `Load()`, `Validate()`, `DefaultConfig()`
- **Dependencies**: viper, yaml

#### `internal/searxng`
- **Responsibility**: HTTP client for SearXNG API interaction
- **Key Functions**: `NewClient()`, `Search()`, `buildURL()`
- **Dependencies**: net/http, encoding/json

#### `internal/formatter`
- **Responsibility**: Output formatting in multiple formats
- **Key Functions**: `Format()`, `NewFormatter()`
- **Dependencies**: encoding/json, text/template

#### `internal/cli`
- **Responsibility**: Cobra command definitions and execution
- **Key Functions**: `NewRootCommand()`, `PreRun()`, `Run()`
- **Dependencies**: cobra, viper

#### `pkg/version`
- **Responsibility**: Version information
- **Key Functions**: `Version()`, `Commit()`, `BuildDate()`

## Data Flow

### 1. CLI Invocation Flow

```
User Input
    ↓
Cobra CLI (root.go)
    ↓
PreRun: Load Config (flags → env → file → defaults)
    ↓
Validate Input
    ↓
Run: Create SearXNG Client
    ↓
Execute Search Request
    ↓
Parse Response
    ↓
Format Output
    ↓
Display Results
```

### 2. Configuration Priority Chain

```
1. CLI Flags (highest priority)
   ↓
2. Environment Variables
   ↓
3. Config File (~/.search/config.yaml)
   ↓
4. Default Values (lowest priority)
```

### 3. Search Request Flow

```
Query + Options
    ↓
SearXNG Client
    ↓
Build HTTP Request
    ↓
Send to SearXNG Instance
    ↓
Parse JSON Response
    ↓
Return Results
```

## Development Tasks

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/mule-ai/search.git
cd search

# Install dependencies
go mod download

# Create config directory and default config
mkdir -p ~/.search
cp examples/config.yaml ~/.search/config.yaml

# Build the binary
make build

# Run tests
make test
```

### Running the CLI

```bash
# From local build
./bin/search "golang tutorials"

# Or install to PATH
make install
search "golang tutorials"
```

### Testing Different Components

```bash
# Test config loading
go test -v ./internal/config/...

# Test SearXNG client with mocks
go test -v ./internal/searxng/...

# Test formatters
go test -v ./internal/formatter/...

# Test CLI commands
go test -v ./internal/cli/...

# Run integration tests (requires live SearXNG instance)
go test -tags=integration -v ./tests/...

# Run all tests with coverage
make coverage
```

### Debugging

Enable verbose output:

```bash
search -v "your query"
```

Add debug prints in code:

```go
if config.Verbose {
    fmt.Fprintf(os.Stderr, "Debug: %v\n", value)
}
```

### Testing Against Local SearXNG Instance

```bash
# Using Docker
docker run -d -p 8888:8080 searxng/searxng

# Test against local instance
search -i http://localhost:8888 "test query"
```

## Common Development Tasks

### Adding a New CLI Flag

1. Add flag definition in `internal/cli/flags.go`:

```go
var NewFlag string

func AddNewFlag(cmd *cobra.Command) {
    cmd.Flags().StringVarP(&NewFlag, "new-flag", "n", "", "Description")
    viper.BindPFlag("new_flag", cmd.Flags().Lookup("new-flag"))
}
```

2. Add to config struct in `internal/config/config.go`:

```go
type Config struct {
    NewFlag string `yaml:"new_flag" mapstructure:"new_flag"`
}
```

3. Update `SPEC.md` with documentation

### Adding a New Output Format

1. Create new formatter in `internal/formatter/`:

```go
// xml.go
package formatter

type XMLFormatter struct{}

func (f *XMLFormatter) Format(resp *SearchResponse) (string, error) {
    // implementation
}

func init() {
    formatters["xml"] = &XMLFormatter{}
}
```

2. Add "xml" to valid formats list

3. Add tests in `internal/formatter/formatter_test.go`

### Modifying the SearXNG Client

1. Update request types in `internal/searxng/types.go`
2. Modify client methods in `internal/searxng/client.go`
3. Update tests to handle new parameters
4. Test against live instance: `go test -tags=integration ./internal/searxng/...`

### Adding a New Search Category

1. Add category constant in `internal/cli/flags.go`:

```go
const (
    CategoryGeneral = "general"
    CategoryImages  = "images"
    CategoryNewType = "newtype" // Add here
)
```

2. Add validation

3. Update formatter to handle new result type if needed

## Build and Release Process

### Local Build

```bash
# Standard build
make build

# Build with version info
make build VERSION=1.0.0 COMMIT=abc123 DATE=$(date -u +%Y-%m-%dT%H:%M:%S)

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o bin/search-linux-amd64 ./cmd/search
GOOS=darwin GOARCH=arm64 go build -o bin/search-darwin-arm64 ./cmd/search
GOOS=windows GOARCH=amd64 go build -o bin/search-windows-amd64.exe ./cmd/search
```

### Testing Before Release

```bash
# Run full test suite
make test

# Check test coverage
make coverage

# Run linters
make lint

# Manual smoke tests
./bin/search -v "test"
./bin/search -f json "test"
./bin/search -f markdown "test"
./bin/search -n 5 "test"
```

### Creating a Release

1. Update version in `pkg/version/version.go`
2. Update CHANGELOG.md
3. Create git tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions will build and publish release

## Performance Considerations

### Bottlenecks

1. **Network latency**: SearXNG API calls are the primary bottleneck
2. **JSON parsing**: For large result sets (>100 results)
3. **Terminal output**: Markdown rendering can be slow for TTY

### Optimization Strategies

1. Use streaming for large JSON responses
2. Implement result caching for repeated queries
3. Lazy load formatter results
4. Use buffer pools for string concatenation

### Benchmarking

```bash
# Run benchmarks
go test -bench=. ./...

# Benchmark specific package
go test -bench=BenchmarkFormat ./internal/formatter/...

# Profile memory
go test -memprofile=mem.prof ./internal/searxng/...
go tool pprof mem.prof
```

## Code Review Guidelines

### What to Look For

1. **Correctness**: Does it work as intended?
2. **Error Handling**: Are errors handled gracefully?
3. **Testing**: Is there adequate test coverage?
4. **Documentation**: Are functions documented?
5. **Style**: Does it follow Go conventions?
6. **Security**: Are inputs validated and sanitized?
7. **Performance**: Are there obvious performance issues?

### Review Checklist

- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] New code has tests
- [ ] Documentation updated
- [ ] No breaking changes (or documented if intentional)
- [ ] Error messages are user-friendly
- [ ] No hardcoded values (use config)
- [ ] Follows semantic versioning

## Troubleshooting

### Common Issues

**Build fails:**
```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

**Tests fail:**
```bash
# Update test dependencies
go mod download

# Run with verbose output
go test -v ./...
```

**Config not loading:**
```bash
# Check config path
ls -la ~/.search/config.yaml

# Validate YAML syntax
cat ~/.search/config.yaml

# Enable verbose output
search -v "test"
```

**SearXNG API errors:**
```bash
# Test instance directly
curl "https://search.butler.ooo/search?q=test&format=json"

# Check timeout
search -t 60 "test"
```

## Resources

### Internal Documentation
- [SPEC.md](../SPEC.md) - Full specification
- [README.md](../README.md) - User documentation
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Contribution guidelines
- [configuration.md](configuration.md) - Configuration details

### External Resources
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)
- [SearXNG API Documentation](https://searxng.github.io/searxng/dev/search_api.html)
- [Effective Go](https://golang.org/doc/effective_go.html)
