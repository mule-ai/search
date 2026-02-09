# Architecture Decisions

This document records significant architectural decisions made during the development of the Search CLI tool.

## Table of Contents

- [Decision 001: Go as Implementation Language](#decision-001-go-as-implementation-language)
- [Decision 002: Cobra + Viper for CLI](#decision-002-cobra--viper-for-cli)
- [Decision 003: Configuration Precedence Chain](#decision-003-configuration-precedence-chain)
- [Decision 004: JSON for SearXNG API](#decision-004-json-for-searxng-api)
- [Decision 005: Formatter Interface Pattern](#decision-005-formatter-interface-pattern)
- [Decision 006: Error Handling Strategy](#decision-006-error-handling-strategy)
- [Decision 007: Package Structure](#decision-007-package-structure)
- [Decision 008: Testing Strategy](#decision-008-testing-strategy)
- [Decision 009: Version Management](#decision-009-version-management)
- [Decision 010: Colored Output Detection](#decision-010-colored-output-detection)

---

## Decision 001: Go as Implementation Language

**Status**: Accepted

**Context**: Need to choose a programming language for implementing a CLI tool that interacts with a web API.

**Decision**: Use Go (Golang) for implementation.

**Rationale**:
- **Single binary distribution**: Go compiles to a single static binary, simplifying distribution
- **Cross-compilation**: Easy to build for multiple platforms (Linux, macOS, Windows) from a single codebase
- **Performance**: Good performance for I/O-bound operations with goroutines
- **Standard library**: Excellent HTTP client and JSON parsing built-in
- **CLI ecosystem**: Strong CLI framework support (Cobra, Viper)
- **Type safety**: Compiled language catches errors at build time
- **Maintenance**: Popular language with good tooling (gofmt, go vet, etc.)

**Consequences**:
- Positive: Easy distribution, good performance, strong tooling
- Positive: No runtime dependencies required
- Negative: Slightly longer development time compared to scripting languages
- Negative: Larger binary size compared to C-based tools

**Alternatives Considered**:
- **Python**: Faster development but requires runtime/dependencies
- **Rust**: Better performance but steeper learning curve
- **Node.js**: Good ecosystem but requires Node.js runtime

---

## Decision 002: Cobra + Viper for CLI

**Status**: Accepted

**Context**: Need a CLI framework that handles commands, flags, and configuration management.

**Decision**: Use Cobra for command framework and Viper for configuration management.

**Rationale**:
- **Industry standard**: Cobra + Viper are the de-facto standard for Go CLIs (used by Docker, Kubernetes, Hugo)
- **Feature-rich**: Built-in support for flags, subcommands, help text, shell completion
- **Integration**: Cobra and Viper are designed to work together
- **Flag binding**: Viper can bind to Cobra flags automatically
- **Config file support**: Viper supports multiple formats (JSON, YAML, TOML)
- **Environment variables**: Viper provides automatic environment variable mapping

**Consequences**:
- Positive: Reduced boilerplate code
- Positive: Consistent user experience with other Go CLIs
- Positive: Automatic shell completion generation
- Negative: Adds two external dependencies
- Negative: Some complexity in flag precedence logic

**Alternatives Considered**:
- **Standard library flag package**: Too basic, requires more code
- **urfave/cli**: Simpler but less feature-rich than Cobra
- **kingpin**: Good alternative but smaller community

---

## Decision 003: Configuration Precedence Chain

**Status**: Accepted

**Context**: Users need multiple ways to configure the CLI (flags, environment, config file). Must define precedence order.

**Decision**: Precedence order: CLI flags > Environment variables > Config file > Default values

**Rationale**:
- **Standard 12-factor app pattern**: Matches common practice for configuration
- **User expectations**: CLI flags are most direct (user typed them)
- **Environment variables**: Useful for CI/CD automation
- **Config file**: Persistent user preferences
- **Defaults**: Reasonable behavior out-of-the-box

**Implementation**:
```go
// Order of precedence (highest to lowest):
// 1. CLI flags (direct user input)
// 2. Environment variables (SEARCH_*)
// 3. Config file (~/.search/config.yaml)
// 4. Default values (hardcoded)
```

**Consequences**:
- Positive: Flexible configuration options
- Positive: Works well in automated environments
- Negative: Complexity in testing all combinations
- Negative: Users might be confused about which value is being used

**Alternatives Considered**:
- **Flags only**: Too restrictive for frequent customization
- **Config file only**: Not flexible for CI/CD
- **Environment only**: Not user-friendly for interactive use

---

## Decision 004: JSON for SearXNG API

**Status**: Accepted

**Context**: SearXNG supports multiple response formats (JSON, RSS, CSV). Need to choose one for parsing.

**Decision**: Use JSON format for all API responses.

**Rationale**:
- **Native parsing**: Go has excellent JSON support in standard library
- **Structured data**: JSON provides nested structures for complex results
- **Type safety**: Can map directly to Go structs
- **Extensibility**: Easy to add new fields without breaking changes
- **Performance**: JSON parsing is fast and efficient
- **Standard**: Most modern APIs use JSON

**Consequences**:
- Positive: Straightforward parsing to Go structs
- Positive: Good error messages for malformed JSON
- Positive: Can handle nested data structures
- Negative: Verbose compared to CSV
- Negative: Slightly larger response size

**Alternatives Considered**:
- **RSS/Atom**: Good for feed readers but less structured
- **CSV**: Simple but lacks nested structure support
- **HTML scraping**: Too fragile, not an official API format

---

## Decision 005: Formatter Interface Pattern

**Status**: Accepted

**Context**: Need to support multiple output formats (JSON, Markdown, Plaintext) with extensibility.

**Decision**: Define a `Formatter` interface and implement separate formatters for each output type.

**Rationale**:
- **Open/closed principle**: Easy to add new formats without modifying existing code
- **Interface segregation**: Clean separation between data and presentation
- **Testing**: Each formatter can be tested independently
- **Factory pattern**: Easy to instantiate formatters by name
- **Type safety**: Compile-time checking of formatter implementations

**Implementation**:
```go
type Formatter interface {
    Format(resp *SearchResponse) (string, error)
}

func NewFormatter(format string) (Formatter, error) {
    switch format {
    case "json":
        return &JSONFormatter{}, nil
    case "markdown":
        return &MarkdownFormatter{}, nil
    case "text":
        return &TextFormatter{}, nil
    default:
        return nil, fmt.Errorf("unknown format: %s", format)
    }
}
```

**Consequences**:
- Positive: Easy to add new output formats
- Positive: Consistent interface across formatters
- Positive: Isolated testing
- Negative: Slightly more code than simple switch statement
- Negative: Must maintain interface contract

**Alternatives Considered**:
- **Single format function**: Hard to extend, mixed concerns
- **Template-based**: Flexible but harder to type-check
- **Code generation**: Overkill for simple formatters

---

## Decision 006: Error Handling Strategy

**Status**: Accepted

**Context**: Need a consistent approach to error handling across the codebase.

**Decision**: Use idiomatic Go error handling with error wrapping and custom error types.

**Rationale**:
- **Idiomatic Go**: Follows standard Go conventions
- **Error wrapping**: Preserves error context through the call stack
- **User-friendly messages**: Separate internal errors from user-facing messages
- **Type assertions**: Allows specific error handling when needed

**Implementation**:
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// Custom error types for specific scenarios
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}
```

**Consequences**:
- Positive: Clear error flow through the codebase
- Positive: Users see helpful error messages
- Positive: Can check error types for handling
- Negative: More verbose than exception-based systems
- Negative: Must remember to wrap errors at each layer

**Alternatives Considered**:
- **panic/recover**: Only for truly unrecoverable errors
- **Custom error package**: Overkill for this project size
- **Error codes**: Useful but less idiomatic in Go

---

## Decision 007: Package Structure

**Status**: Accepted

**Context**: Need to organize code into logical packages with clear responsibilities.

**Decision**: Use standard Go project layout with `internal/` for private packages and `pkg/` for public packages.

**Rationale**:
- **Standard layout**: Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- **Internal packages**: `internal/` prevents external imports, enforces API boundaries
- **Clear separation**: Each package has a single responsibility
- **Testability**: Packages can be tested independently
- **Discoverability**: Standard layout makes code easy to find

**Structure**:
```
cmd/           - Application entry points
internal/      - Private application code
  config/      - Configuration management
  searxng/     - SearXNG API client
  formatter/   - Output formatters
  cli/         - CLI commands and flags
pkg/           - Public libraries
  version/     - Version information
```

**Consequences**:
- Positive: Clear organization and boundaries
- Positive: Easy to navigate
- Positive: `internal/` prevents unintended external dependencies
- Negative: More directory levels than flat structure
- Negative: Must plan package boundaries carefully

**Alternatives Considered**:
- **Flat structure**: Harder to organize as project grows
- **Domain-driven**: Overkill for this project size
- **Package-by-layer**: Mixes concerns (e.g., all models together)

---

## Decision 008: Testing Strategy

**Status**: Accepted

**Context**: Need comprehensive testing to ensure code quality and reliability.

**Decision**: Multi-layer testing approach with unit tests, integration tests, and benchmarks.

**Strategy**:

1. **Unit Tests** (Primary)
   - Test each package in isolation
   - Mock external dependencies (HTTP, filesystem)
   - Aim for >80% code coverage
   - Use table-driven tests for multiple scenarios

2. **Integration Tests** (Secondary)
   - Use build tags (`//go:build integration`)
   - Run against test SearXNG instance
   - Test full request/response cycle
   - Not run in normal test suite

3. **Benchmark Tests** (Performance)
   - Profile critical paths (HTTP, formatting)
   - Track performance over time
   - Identify regressions

**Implementation**:
```go
// Unit test with mock
func TestClient_Search(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(mockResponse)
    }))
    defer server.Close()

    client := NewClient(Config{Instance: server.URL})
    results, err := client.Search(ctx, "test")
    // assertions...
}

// Integration test
//go:build integration
func TestIntegrationRealSearch(t *testing.T) {
    client := NewClient(Config{Instance: "https://search.butler.ooo"})
    results, err := client.Search(ctx, "test")
    // assertions...
}
```

**Consequences**:
- Positive: Comprehensive test coverage
- Positive: Fast unit tests (no real HTTP)
- Positive: Integration tests verify real API
- Negative: Must maintain mock data
- Negative: Integration tests require network access

**Alternatives Considered**:
- **Only unit tests**: Would miss real-world issues
- **Only integration tests**: Too slow and unreliable
- **Property-based testing**: Overkill for this project

---

## Decision 009: Version Management

**Status**: Accepted

**Context**: Need to track and display version information for debugging and releases.

**Decision**: Use build-time ldflags to inject version information into the binary.

**Rationale**:
- **Single source of truth**: Version defined in one place
- **No import cycles**: No need to import own package for version
- **Git integration**: Can use git tags and commit hashes
- **Flexible**: Different version info per build
- **Standard practice**: Used by many Go projects

**Implementation**:
```go
// pkg/version/version.go
package version

var (
    Version   = "dev"         // Set by ldflags
    Commit    = "unknown"     // Set by ldflags
    BuildDate = "unknown"     // Set by ldflags
)

func Print() string {
    return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
}
```

**Build command**:
```bash
go build -ldflags="-X 'github.com/mule-ai/search/pkg/version.Version=1.0.0' -X 'github.com/mule-ai/search/pkg/version.Commit=abc123' -X 'github.com/mule-ai/search/pkg/version.BuildDate=$(date)'"
```

**Consequences**:
- Positive: Exact version info in binary
- Positive: Useful for debugging
- Positive: Can track which commit deployed
- Negative: Requires build script or makefile
- Negative: Complex build command

**Alternatives Considered**:
- **Hardcoded version**: Would require editing code for each release
- **Git tags only**: Can't query tags from binary easily
- **Separate version file**: Risk of version mismatch

---

## Decision 010: Colored Output Detection

**Status**: Accepted

**Context**: Want to provide colored terminal output for better readability, but only when appropriate.

**Decision**: Automatically detect TTY support and provide `--no-color` flag to override.

**Rationale**:
- **Better UX**: Colors improve readability in terminal
- **Smart defaults**: Detect if output is going to a terminal
- **User control**: Flag to disable colors for piping/redirection
- **Standards compliant**: Respect `NO_COLOR` environment variable

**Implementation**:
```go
func shouldColorize() bool {
    // Check --no-color flag
    if noColor {
        return false
    }
    // Check NO_COLOR environment variable
    if os.Getenv("NO_COLOR") != "" {
        return false
    }
    // Check if stdout is a TTY
    if !isatty.IsTerminal(os.Stdout.Fd()) {
        return false
    }
    return true
}
```

**Consequences**:
- Positive: Better user experience in terminals
- Positive: Doesn't break scripts/pipelines
- Positive: Respects user preferences
- Negative: Adds dependency (isatty or similar)
- Negative: More complex output logic

**Alternatives Considered**:
- **Always colored**: Breaks scripts and file output
- **Never colored**: Less readable in terminal
- **Config file only**: Not flexible enough

---

## Template for Future Decisions

When adding new architectural decisions, use this template:

```markdown
## Decision XXX: [Title]

**Status**: Proposed | Accepted | Deprecated | Superseded

**Context**: [What is the issue we're facing]

**Decision**: [What we decided]

**Rationale**: [Why we made this decision]

**Implementation**: [How we implemented it]

**Consequences**:
- Positive: [Benefits]
- Negative: [Drawbacks]

**Alternatives Considered**:
- [Alternative 1]: [Pros/cons]
- [Alternative 2]: [Pros/cons]
```