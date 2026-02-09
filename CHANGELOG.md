# Changelog

All notable changes to the `search` CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-09

### Added
- Initial stable release of search CLI
- SearXNG API integration with configurable instances
- Multiple output formats: JSON, Markdown, and Plaintext
- Configuration file support (`~/.search/config.yaml`)
- Environment variable support for CI/CD integration
- Command-line flags for all configuration options
- Shell completion support (bash, zsh, fish)
- Search category filtering (general, images, videos, news, map, music, it, science, files)
- Pagination support with `--page` flag
- Time range filtering (day, week, month, year)
- Colored output for terminals with `--no-color` option
- Browser integration with `--open` and `--open-all` flags
- Verbose mode for debugging
- Comprehensive error handling with user-friendly messages
- Input validation for all parameters
- Homebrew formula for macOS installation
- Scoop manifest for Windows installation
- Arch Linux AUR package support
- Cross-platform builds (Linux, macOS, Windows) for multiple architectures
- UPX compression for smaller binary sizes
- Comprehensive test suite with >80% code coverage
- Benchmark tests for performance monitoring

### Configuration
- Default SearXNG instance: `https://search.butler.ooo`
- Config precedence: CLI flags > environment variables > config file > defaults
- Configurable timeout (default: 30 seconds)
- Configurable result count (default: 10, range: 1-100)
- Configurable language (default: en)
- Configurable safe search levels (0=off, 1=moderate, 2=strict)

### Documentation
- Complete README with installation and usage examples
- Configuration guide in `docs/configuration.md`
- Example commands in `examples/` directory
- CONTRIBUTING.md with development workflow
- Go doc comments on all exported functions
- Performance benchmarks in `BENCHMARKS.md`
- Performance optimization notes in `PERFORMANCE.md`

### Testing
- Unit tests for all packages
- Integration tests against live SearXNG instances
- Benchmark tests for performance monitoring
- Code coverage reporting
- CI/CD integration with GitHub Actions

[1.0.0]: https://github.com/mule-ai/search/releases/tag/v1.0.0