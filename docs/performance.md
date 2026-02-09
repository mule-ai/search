# Performance Optimization Guide

This document describes the performance optimizations implemented in the search CLI and provides guidance for further optimization work.

## Implemented Optimizations

### 1. Buffer Pooling

**Location:** `internal/searxng/pool.go`

We use `sync.Pool` to reuse buffers and byte slices, reducing memory allocations and GC pressure.

```go
// GetBuffer retrieves a buffer from the pool
buf := searxng.GetBuffer()
defer searxng.PutBuffer(buf)
```

**Benefits:**
- Reduces heap allocations
- Lowers garbage collection overhead
- Improves performance for repeated operations

### 2. Optimized JSON Decoding

**Location:** `internal/searxng/decoder.go`

Streaming JSON decoder with buffer pooling for efficient parsing of SearXNG responses.

```go
// Use the optimized decoder
resp, err := searxng.DecodeResponse(responseBody)
```

**Benefits:**
- Lower memory usage for large responses
- Faster JSON parsing with streaming decoder
- Automatic buffer reuse

### 3. Optimized Formatters

**Location:** `internal/formatter/optimized.go`

Formatters that use buffer pooling instead of string concatenation.

```go
formatter := NewOptimizedFormatter("json")
output, err := formatter.Format(response)
```

**Benefits:**
- Reduced string allocations
- Faster formatting for large result sets
- Lower memory footprint

### 4. Optimized Builds

**Location:** `Makefile`

Build targets that produce smaller, faster binaries.

```bash
# Optimized build (smaller binary)
make optimize

# Optimized + UPX compressed
make optimize-upx
```

**Optimizations applied:**
- `-trimpath`: Remove file system paths from binary
- `-s -w`: Strip debug information and reduce DWARF table size
- Build-time optimizations

**Binary sizes:**
- Standard build: ~15-20 MB
- Optimized build: ~10-12 MB
- UPX compressed: ~3-5 MB

## Profiling

### CPU Profiling

Generate a CPU profile to identify bottlenecks:

```bash
# Generate profile
make profile

# View profile
go tool pprof profiles/cpu.prof

# Interactive analysis
go tool pprof -http=:8080 profiles/cpu.prof
```

### Memory Profiling

Generate a memory profile:

```bash
# Generate profile
go test -memprofile=mem.prof -bench=. ./...

# View profile
go tool pprof profiles/mem.prof

# View top allocations
go tool pprof -top profiles/mem.prof
```

### Benchmarking

Run benchmarks to track performance:

```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkJSONParsing -benchmem ./...

# Compare before/after
go test -bench=. -benchmem ./... > old.txt
# Make changes
go test -bench=. -benchmem ./... > new.txt
benchstat old.txt new.txt
```

## Performance Targets

### Goals

- **Startup time:** < 50ms
- **Search latency:** < 500ms (excluding network)
- **Memory usage:** < 50MB for typical queries
- **Binary size:** < 10MB (uncompressed)

### Current Benchmarks

Run `go test -bench=. -benchmem ./internal/...` to see current performance.

#### JSON Parsing (10 results)
- Target: < 100 µs per operation
- Target: < 5 allocations per operation

#### JSON Parsing (100 results)
- Target: < 500 µs per operation
- Target: < 10 allocations per operation

#### Formatting (50 results)
- Target: < 1 ms per operation
- Target: < 20 allocations per operation

## Optimization Strategies

### Completed

1. ✅ Buffer pooling for JSON encoding/decoding
2. ✅ Streaming JSON decoder
3. ✅ Optimized formatters using bytes.Buffer
4. ✅ Build flags for smaller binary size

### Future Opportunities

1. **Connection pooling** - Reuse HTTP connections for multiple searches
2. **Response caching** - Cache common queries
3. **Parallel processing** - Process multiple searches concurrently
4. **Lazy evaluation** - Defer formatting until output needed
5. **String interning** - Reduce duplicate string allocations
6. **Custom JSON marshaling** - Implement MarshalJSON for complex types

## Monitoring Performance

### In Production

```bash
# Enable profiling flags
SEARCH_PROFILE_CPU=cpu.prof SEARCH_PROFILE_MEM=mem.prof search "query"

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Continuous Benchmarking

Add to CI/CD:

```yaml
- name: Run benchmarks
  run: |
    go test -bench=. -benchmem ./... > bench.txt
    # Upload to benchmark tracking service
```

## Tips for Further Optimization

1. **Profile before optimizing** - Always use pprof to identify real bottlenecks
2. **Measure impact** - Use benchmarks to verify improvements
3. **Consider trade-offs** - Optimization vs. code readability
4. **Test with real data** - Benchmarks should match production workloads
5. **Monitor over time** - Performance can regress with new features

## Tools

- **pprof**: CPU and memory profiling
- **benchstat**: Benchmark comparison tool (`go install golang.org/x/perf/cmd/benchstat@latest`)
- **go test -bench**: Run benchmarks
- **upx**: Binary compression tool

## Resources

- [Go Profiling](https://go.dev/doc/diagnostics#profiling)
- [Optimizing Go code](https://go.dev/doc/diagnostics#optimization)
- [sync.Pool documentation](https://pkg.go.dev/sync#Pool)
- [GoReleaser](https://goreleaser.com/) for release optimization
