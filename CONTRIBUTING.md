# Contributing to Search CLI

Thank you for your interest in contributing to the Search CLI project! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow:

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Go 1.21 or later
- Make (optional, for build automation)
- Git

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/search.git
   cd search
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/mule-ai/search.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Build the project:
   ```bash
   make build
   # or
   go build -o bin/search ./cmd/search
   ```

## Development Workflow

### 1. Create a Branch

Create a new branch for your contribution from the `main` branch:

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Your Changes

- Write code following the [Coding Standards](#coding-standards)
- Add tests for new functionality
- Update documentation as needed
- Ensure all tests pass: `make test`

### 3. Test Your Changes

Run the test suite:

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific package tests
go test ./internal/config/...
go test ./internal/searxng/...
go test ./internal/formatter/...
go test ./internal/cli/...
```

### 4. Commit Your Changes

Write clear, descriptive commit messages following the guidelines in [Commit Messages](#commit-messages).

### 5. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 6. Create a Pull Request

- Go to the original repository on GitHub
- Click "New Pull Request"
- Provide a clear description of your changes
- Link any related issues

## Coding Standards

### Go Code Style

Follow standard Go conventions as described in [Effective Go](https://golang.org/doc/effective_go.html):

- Use `gofmt` for formatting: `gofmt -w .`
- Use `golint` for linting: `golint ./...`
- Use `go vet` for static analysis: `go vet ./...`

### Naming Conventions

- **Package names**: Short, lowercase, single words when possible
- **Constants**: `PascalCase` or `UPPER_SNAKE_CASE` for exported constants
- **Variables**: `camelCase` for local variables, `PascalCase` for exported
- **Functions**: `PascalCase` for exported, `camelCase` for private

### File Organization

- One package per directory
- Keep files focused on a single responsibility
- Use subdirectories for related functionality
- Test files should be named `*_test.go` and co-located with source

### Documentation

- All exported functions must have godoc comments
- Include usage examples in godoc comments
- Add package-level documentation at the top of each package
- Update README.md for user-facing changes

Example godoc comment:

```go
// Search performs a search query against the configured SearXNG instance.
//
// The query parameter is the search string. Options can be provided to
// customize the search behavior.
//
// Example:
//
//	results, err := client.Search(ctx, "golang tutorials", searxng.Options{
//	    Language: "en",
//	    Results:  10,
//	})
//
// Returns a SearchResponse with results or an error if the request fails.
func (c *Client) Search(ctx context.Context, query string, opts Options) (*SearchResponse, error) {
	// implementation
}
```

### Error Handling

- Always handle errors, never ignore them
- Wrap errors with context using `fmt.Errorf` or `errors.Wrap`
- Return errors from functions when they occur
- Create custom error types for package-specific errors

Example:

```go
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

## Testing Guidelines

### Test Coverage

We aim for >80% code coverage. All new code should include appropriate tests.

### Writing Tests

- Write unit tests for all exported functions
- Use table-driven tests for multiple test cases
- Mock external dependencies (HTTP clients, file system)
- Test both success and failure paths

Example table-driven test:

```go
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Instance: "https://search.example.com",
				Results:  10,
			},
			wantErr: false,
		},
		{
			name: "missing instance",
			config: Config{
				Results: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

### Integration Tests

Integration tests should run against a test SearXNG instance. These are marked with a build tag:

```go
//go:build integration
// +build integration

func TestIntegrationSearch(t *testing.T) {
	// integration test code
}
```

Run integration tests:
```bash
go test -tags=integration ./...
```

### Benchmark Tests

For performance-critical code, write benchmark tests:

```go
func BenchmarkFormatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Format(results)
	}
}
```

## Commit Messages

Follow conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or tooling changes

### Examples

```
feat(cli): add --open flag to open results in browser

Implement browser integration to open the first search result
directly in the system's default web browser.

Closes #42
```

```
fix(formatter): handle nil results in markdown formatter

Previously, the markdown formatter would panic when receiving
nil results. This commit adds nil checks and returns an error
instead.

Fixes #38
```

```
docs: update installation instructions

Clarify the installation from source instructions and add
troubleshooting section for common build issues.
```

## Pull Request Process

### Before Submitting

1. Ensure your code passes all tests: `make test`
2. Check code coverage: `make coverage`
3. Run linters: `make lint`
4. Update documentation if needed
5. Rebase your branch on latest main: `git rebase upstream/main`

### Pull Request Description

Include:

- **Summary**: Brief description of changes
- **Motivation**: Why this change is needed
- **Changes**: List of files/areas modified
- **Testing**: How you tested your changes
- **Related Issues**: Link to related issues

Example:

```
## Summary
Adds support for search categories (images, videos, news, etc.)

## Motivation
Users need to search specific content types without visiting the web UI.

## Changes
- Added --category flag to CLI
- Added category validation
- Updated formatter to handle image results
- Added tests for category selection

## Testing
- Added unit tests for category validation
- Manually tested with all category types
- Verified error handling for invalid categories

## Related Issues
Closes #15
```

### Review Process

1. Automated checks (CI) will run
2. Maintainers will review your code
3. Address review comments
4. Once approved, a maintainer will merge

### Feedback

Be open to feedback on your pull request. Maintainers may request:

- Code changes for consistency
- Additional test coverage
- Documentation updates
- Performance improvements

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## Questions?

If you have questions about contributing:

- Open an issue with your question
- Start a discussion in the repository's Discussions tab
- Check existing issues and pull requests for similar topics

Thank you for your contributions!