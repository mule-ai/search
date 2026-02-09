# Release Notes for v{{VERSION}}

## What's Changed

### Features
- Feature 1 description
- Feature 2 description

### Bug Fixes
- Bug fix 1 description
- Bug fix 2 description

### Performance
- Performance improvement 1

### Documentation
- Documentation update 1

## Installation

### From Binary
Download the appropriate binary for your platform from the [Assets](https://github.com/mule-ai/search/releases/tag/v{{VERSION}}) section below.

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

### From Source
```bash
git clone https://github.com/mule-ai/search.git
cd search
make install
```

## Verification

Verify the installation:
```bash
search --version
```

## Upgrade

### Homebrew
```bash
brew upgrade search
```

### Scoop
```powershell
scoop update search
```

### From Source
```bash
cd search
git fetch
git checkout v{{VERSION}}
make install
```

## Full Changelog
https://github.com/mule-ai/search/compare/v{{PREVIOUS_VERSION}}...v{{VERSION}}