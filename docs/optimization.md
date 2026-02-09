# Binary Optimization Guide

This guide explains how to build and optimize the `search` CLI binary for smaller size and better distribution.

## Overview

The search CLI binary can be significantly reduced in size through optimization techniques:

- **Standard build**: ~15-20 MB (with debug info)
- **Optimized build**: ~10-12 MB (stripped debug info)
- **UPX compressed**: ~4-6 MB (compressed executable)

## Optimization Methods

### 1. Standard Build

The default build includes debug information and symbols:

```bash
go build -o bin/search ./cmd/search
```

### 2. Optimized Build

Strip debug symbols and reduce binary size:

```bash
go build -trimpath -ldflags="-s -w" -o bin/search ./cmd/search
```

**Flags explained:**
- `-trimpath`: Removes file system paths from the binary
- `-ldflags="-s -w"`:
  - `-s`: Strip symbol table
  - `-w`: Strip DWARF debug info

### 3. UPX Compression

UPX (Ultimate Packer for eXecutables) compresses the executable while maintaining the ability to run it directly.

#### Install UPX

**Linux (Debian/Ubuntu):**
```bash
sudo apt-get install upx
```

**macOS:**
```bash
brew install upx
```

**Fedora/RHEL:**
```bash
sudo dnf install upx
```

**Arch Linux:**
```bash
sudo pacman -S upx
```

#### Compress the Binary

```bash
# Basic compression
upx bin/search

# Best compression (slower)
upx --best bin/search

# Best compression with LZMA (slowest, smallest)
upx --best --lzma bin/search
```

#### Decompress if Needed

```bash
upx -d bin/search
```

## Using the Makefile

The project provides convenient Makefile targets:

### `make optimize`

Build an optimized binary without debug info:

```bash
make optimize
```

### `make optimize-upx`

Build and compress with UPX:

```bash
make optimize-upx
```

This will:
1. Build an optimized binary
2. Compress it with UPX using best LZMA compression
3. Verify the compressed binary still works

### `make optimize-full`

Run comprehensive optimization with size comparison:

```bash
make optimize-full
```

This will:
1. Build standard, optimized, and compressed versions
2. Show size comparison table
3. Test the compressed binary
4. Display recommendations

## Using the Optimization Script

For detailed optimization analysis, use the provided script:

```bash
./scripts/optimize.sh
```

This script creates three binaries:
- `bin/search-standard` - Standard build
- `bin/search-optimized` - Optimized build (stripped)
- `bin/search` - UPX compressed

It also provides a size comparison table and tests the final binary.

## GoReleaser Integration

The `.goreleaser.yml` configuration includes automatic UPX compression for releases:

```yaml
upx:
  - enabled: true
    goos:
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
    compress:
      - --best
      - --lzma
```

When you create a release, GoReleaser automatically:
1. Builds optimized binaries for all platforms
2. Compresses Linux and macOS builds with UPX
3. Creates checksums and signatures

## Size Comparison

Typical binary sizes (Linux amd64):

| Build Type | Size | Reduction |
|------------|------|-----------|
| Standard | ~18 MB | - |
| Optimized | ~11 MB | 39% |
| UPX Compressed | ~4.5 MB | 75% |

## Performance Impact

### Startup Time

- Standard: ~5ms
- Optimized: ~4ms (slightly faster)
- UPX Compressed: ~15-20ms (decompression overhead)

The decompression overhead is negligible for CLI usage (milliseconds).

### Runtime Performance

No impact on runtime performance. The binary is decompressed into memory on first launch.

## Recommendations

### For Development

Use standard builds for faster compilation:

```bash
go build -o bin/search ./cmd/search
```

### For Distribution

Use optimized builds:

```bash
make optimize
```

### For Releases

Let GoReleaser handle optimization automatically:

```bash
make snapshot
```

### For Minimal Size

Use UPX compression:

```bash
make optimize-upx
```

## Troubleshooting

### UPX Fails to Compress

Some Go binaries may have sections that UPX cannot compress. If UPX fails:

1. Try a different compression method:
   ```bash
   upx --best bin/search  # Without LZMA
   ```

2. Use the optimized build without UPX:
   ```bash
   make optimize
   ```

### Compressed Binary Doesn't Work

If the compressed binary fails to run:

1. Decompress and test:
   ```bash
   upx -d bin/search
   bin/search --version
   ```

2. Rebuild with less aggressive compression:
   ```bash
   upx --2 bin/search  # Fast compression
   ```

3. Verify the original binary works:
   ```bash
   make optimize
   bin/search --version
   ```

### Antivirus False Positives

UPX-compressed executables may trigger some antivirus software. If this is a concern:

1. Use the optimized build without UPX
2. Sign the binary (if distributing)
3. Submit to antivirus vendors as a false positive

## Cross-Platform Notes

### Windows

UPX compression works on Windows but may trigger antivirus software more frequently. Consider using optimized builds without UPX for Windows distributions.

### macOS

UPX compression works well on macOS. However, when notarizing for distribution, you may need to use the optimized build without UPX.

### Linux

UPX compression works excellently on Linux and is recommended for all distributions.

## CI/CD Integration

### GitHub Actions

The release workflow automatically uses GoReleaser with UPX compression:

```yaml
- name: Run GoReleaser
  run: make release-dry
```

### Manual Release

```bash
# Create a snapshot build
make snapshot

# Full release (requires GitHub token)
goreleaser release
```

## Further Optimization

If you need even smaller binaries, consider:

1. **Reduce dependencies**: Audit and remove unused dependencies
2. **Go 1.20+**: Use the new toolchain for better optimization
3. **Static analysis**: Use `go tool nm` to find large symbols
4. **Profile-guided optimization**: (Experimental) Use PGO in Go 1.20+

## Resources

- [UPX Documentation](https://upx.github.io/)
- [Go Binary Size](https://go.dev/doc/diagnostics#tool_asm)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Go Build Constraints](https://go.dev/ref/mod#build-constraints)
