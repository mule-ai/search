# Build & Release Documentation

This document describes the build system and release infrastructure for the Search CLI project.

## Build System

### Makefile Targets

The project uses a comprehensive `Makefile` with the following targets:

#### Core Targets
- `make build` - Build the binary for the current platform
- `make test` - Run all tests
- `make install` - Install to `$GOBIN` or `$HOME/go/bin`
- `make clean` - Clean build artifacts

#### Quality Targets
- `make lint` - Run golangci-lint and go vet
- `make vet` - Run go vet only
- `make coverage` - Generate code coverage report

#### Release Targets
- `make build-all` - Build for all platforms (Linux/macOS/Windows)
- `make dev` - Build development binary with race detection
- `make release-dry` - Run GoReleaser in dry-run mode
- `make snapshot` - Create snapshot builds

### Cross-Compilation

The `make build-all` target builds binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Binaries are placed in the `bin/` directory with platform-specific names:
- `search-linux-amd64`
- `search-linux-arm64`
- `search-darwin-amd64`
- `search-darwin-arm64`
- `search-windows-amd64.exe`

### Version Information

Version information is injected at build time using ldflags:
- `Version` - From `VERSION` variable (default: 0.1.0)
- `GitCommit` - Git commit hash (short)
- `BuildDate` - ISO 8601 timestamp
- `GoVersion` - Go version used for build

Example:
```bash
make build VERSION=1.0.0
```

## Continuous Integration

### GitHub Actions

The project uses GitHub Actions for CI/CD:

#### Test Workflow (`.github/workflows/test.yml`)
Runs on every push and pull request to `main` and `develop` branches:
1. Set up Go 1.21
2. Download and verify dependencies
3. Run `go vet`
4. Check formatting
5. Run tests with race detection
6. Generate coverage report
7. Check coverage threshold (80%)
8. Upload coverage to Codecov
9. Build binary

#### Release Workflow (`.github/workflows/release.yml`)
Runs on tag pushes (e.g., `v1.0.0`):
1. Set up Go
2. Import GPG key for signing
3. Run GoReleaser
4. Build cross-platform binaries
5. Create GitHub release
6. Upload release artifacts

## Release Automation

### GoReleaser Configuration

GoReleaser is configured in `.goreleaser.yml` with:

**Build Settings:**
- Supports Linux, macOS, Windows
- Architectures: amd64, arm64, 386
- Stripped binaries with optimized ldflags
- Trimpath for cleaner builds

**Archive Settings:**
- Platform-specific naming
- Includes LICENSE, README.md, docs/
- ZIP format for Windows

**Checksums:**
- SHA-256 checksums for all artifacts
- Stored in `*_checksums.txt`

**Signing:**
- Optional GPG signing of checksums
- Requires `GPG_PRIVATE_KEY` and `PASSPHRASE` secrets

**Package Managers:**
- Homebrew (tap: mule-ai/homebrew-tap)
- Scoop (bucket: mule-ai/scoop-bucket)
- Nix (flakes: mule-ai/nix-flakes)

**Changelog:**
- Auto-generated from commits
- Excludes docs, test, ci, chore commits

### Package Manager Files

#### Homebrew Formula
Located at `homebrew/search.rb`:
- Defines installation method
- Specifies dependencies
- Includes verification test

#### Scoop Manifest
Located at `scoop/search.json`:
- Windows-specific configuration
- Auto-updates from GitHub releases

#### Arch Linux AUR
Located at `arch/PKGBUILD`:
- PKGBUILD template for creating AUR package
- Supports x86_64 and aarch64

## Installation Methods

### Users can install via:

1. **Homebrew** (macOS/Linux)
   ```bash
   brew install mule-ai/tap/search
   ```

2. **Scoop** (Windows)
   ```powershell
   scoop install mule-ai/scoop-bucket/search
   ```

3. **Arch Linux AUR**
   ```bash
   paru -S search-bin
   ```

4. **Go**
   ```bash
   go install github.com/mule-ai/search/cmd/search@latest
   ```

5. **From Source**
   ```bash
   git clone https://github.com/mule-ai/search
   cd search
   make install
   ```

## Release Process Summary

1. Update version in `Makefile`
2. Update `CHANGELOG.md`
3. Commit changes
4. Tag release: `git tag -a v1.0.0`
5. Push tag: `git push origin main --tags`
6. GitHub Actions automatically:
   - Builds binaries
   - Creates release
   - Publishes to package managers

See `docs/releases.md` for detailed release process.

## Development Build

For development, use:

```bash
# Build with race detection
make dev

# Quick build
make build

# Test everything
make test

# Check code quality
make lint

# Full coverage
make coverage
```

## Troubleshooting

### Build Issues

**Missing golangci-lint:**
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
```

**Version not updating:**
Make sure to pass VERSION variable:
```bash
make build VERSION=1.0.0
```

### CI Issues

**Tests failing locally:**
```bash
make test
make coverage
```

**Coverage below 80%:**
Check `coverage.html` for uncovered areas.

### Release Issues

**GoReleaser dry-run:**
```bash
make release-dry
```

**Tag already exists:**
Delete and recreate:
```bash
git tag -d v1.0.0
git push --delete origin v1.0.0
```

## Security

- GPG signing for releases (optional)
- SHA-256 checksums for all artifacts
- Dependency scanning via Dependabot (if configured)
- Code coverage threshold enforcement

## Performance

- Stripped binaries reduce size ~30%
- Cross-compilation for multiple platforms
- Parallel builds in CI
- Artifact caching in GitHub Actions