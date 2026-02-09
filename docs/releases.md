# Release Process

This document describes the release process for the Search CLI project.

## Prerequisites

1. **GoReleaser** - Install GoReleaser for automated releases:
   ```bash
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
   ```

2. **GPG Key** - For signing releases (optional but recommended):
   ```bash
   gpg --full-generate-key
   ```

3. **GitHub Token** - Set `GITHUB_TOKEN` environment variable with repo permissions

## Making a Release

### 1. Update Version

Update version in `Makefile`:
```bash
# Update VERSION variable
VERSION=1.0.0
```

### 2. Update CHANGELOG

Create/update a `CHANGELOG.md` with the new version changes.

### 3. Commit Changes

```bash
git add .
git commit -m "Release v1.0.0"
```

### 4. Tag the Release

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin main --tags
```

### 5. GoReleaser (Manual)

For testing first:
```bash
make release-dry
```

For actual release:
```bash
goreleaser release --clean
```

### 6. Automated Release

When you push a tag, the GitHub Actions workflow (`.github/workflows/release.yml`) will:
1. Run tests
2. Build binaries for all platforms
3. Create a GitHub release
4. Generate checksums
5. Sign artifacts (if GPG key is configured)

## Homebrew Tap

After release, update the Homebrew tap:
1. Fork/clone https://github.com/mule-ai/homebrew-tap
2. Update the formula with new version and checksums
3. Submit PR

Or use GoReleaser's automatic Homebrew publishing (configured in `.goreleaser.yml`).

## Scoop Bucket

After release, update the Scoop bucket:
1. Fork/clone https://github.com/mule-ai/scoop-bucket
2. Update the manifest with new version
3. Submit PR

Or use GoReleaser's automatic Scoop publishing.

## Arch Linux AUR

For AUR packages:
1. Clone the AUR package: `git clone ssh://aur@aur.archlinux.org/search-bin.git`
2. Update `PKGBUILD` with new version and checksums
3. Update `.SRCINFO`
4. Commit and push: `git commit -am "Update to v1.0.0" && git push`

## Verification

After release, verify:

1. Binaries download and work correctly
2. Checksums match
3. Installation via package managers work
4. Version is correct

## Snapshot Builds

For testing release builds without creating a tag:

```bash
make snapshot
```

This creates binaries in `dist/` directory with version info from current git state.

## Rollback

If something goes wrong:

1. Delete the release from GitHub
2. Delete the tag:
   ```bash
   git tag -d v1.0.0
   git push --delete origin v1.0.0
   ```

3. Create a new release with the correct version

## Release Notes Template

Use `.github/RELEASE_TEMPLATE.md` as a starting point for release notes. Key sections:
- What's Changed
- Features
- Bug Fixes
- Installation instructions
- Upgrade instructions