# Search CLI Makefile

.PHONY: all build clean install test coverage lint vet help build-all dev release-dry snapshot profile optimize optimize-upx benchmark benchcmp

# Version information
VERSION ?= 1.0.0
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
GO_VERSION ?= $(shell go version | cut -d' ' -f3)

# Build flags for ldflags to inject version information
LDFLAGS := -ldflags=" \
	-X github.com/mule-ai/search/pkg/version.Version=$(VERSION) \
	-X github.com/mule-ai/search/pkg/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/mule-ai/search/pkg/version.BuildDate=$(BUILD_DATE) \
	-X github.com/mule-ai/search/pkg/version.GoVersion=$(GO_VERSION) \
	"

# Optimized build flags for smaller binary size
OPT_LDFLAGS := -ldflags=" \
	-X github.com/mule-ai/search/pkg/version.Version=$(VERSION) \
	-X github.com/mule-ai/search/pkg/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/mule-ai/search/pkg/version.BuildDate=$(BUILD_DATE) \
	-X github.com/mule-ai/search/pkg/version.GoVersion=$(GO_VERSION) \
	-s -w -X main.disableMetrics=true \
	"

all: build

# Cross-compilation targets
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/search-linux-amd64 ./cmd/search
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/search-linux-arm64 ./cmd/search
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/search-darwin-amd64 ./cmd/search
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/search-darwin-arm64 ./cmd/search
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/search-windows-amd64.exe ./cmd/search
	@echo "Built binaries:"
	@ls -lh bin/

# Development build with race detection
dev:
	@echo "Building development binary with race detection..."
	go build -race -o bin/search-dev ./cmd/search

# Run GoReleaser in dry-run mode
release-dry:
	@echo "Running GoReleaser in dry-run mode..."
	goreleaser release --skip-publish --skip-sign --clean

# Create a release snapshot
snapshot:
	@echo "Creating snapshot build..."
	goreleaser release --snapshot --clean

help:
	@echo "Search CLI - Available targets:"
	@echo "  build       - Build the search binary"
	@echo "  build-all   - Build for all platforms (Linux, macOS, Windows)"
	@echo "  dev         - Build development binary with race detection"
	@echo "  optimize    - Build optimized smaller binary"
	@echo "  optimize-upx - Build and compress with UPX"
	@echo "  profile     - Generate CPU and memory profiles"
	@echo "  benchmark   - Run performance benchmarks"
	@echo "  benchcmp    - Compare two benchmark results"
	@echo "  test        - Run all tests"
	@echo "  coverage    - Generate coverage report"
	@echo "  lint        - Run go linters (including golangci-lint)"
	@echo "  vet         - Run go vet"
	@echo "  install     - Install to GOBIN"
	@echo "  clean       - Clean build artifacts"
	@echo "  release-dry - Run GoReleaser in dry-run mode"
	@echo "  snapshot    - Create a snapshot build"
	@echo ""
	@echo "Build variables:"
	@echo "  VERSION=$(VERSION)"
	@echo "  GIT_COMMIT=$(GIT_COMMIT)"
	@echo "  BUILD_DATE=$(BUILD_DATE)"
	@echo "  GO_VERSION=$(GO_VERSION)"

build:
	@echo "Building search CLI..."
	@echo "Version: $(VERSION)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"
	@echo "Go version: $(GO_VERSION)"
	go build $(LDFLAGS) -o bin/search ./cmd/search

clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/search
	rm -f coverage.out coverage.html
	go clean ./...

install: build
	@echo "Installing search to /bin..."
	@sudo mkdir -p /bin
	@sudo cp bin/search /bin/search
	@sudo chmod +x /bin/search
	@echo "Installed to /bin/search"

test:
	@echo "Running tests..."
	go test -v ./...

coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep -v "total:"
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not found. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin latest" && exit 1)
	golangci-lint run --timeout=5m ./...
	go vet ./...
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

# Build optimized binary with reduced size
optimize:
	@echo "Building optimized binary..."
	@echo "Version: $(VERSION)"
	go build -trimpath $(OPT_LDFLAGS) -o bin/search ./cmd/search
	@echo "Optimized binary size:"
	@ls -lh bin/search

# Build and compress with UPX (if available)
optimize-upx: optimize
	@echo "Compressing with UPX..."
	@which upx > /dev/null 2>&1 || (echo "UPX not found. Install with: apt-get install upx (Debian/Ubuntu) or brew install upx (macOS)" && exit 1)
	@cp bin/search bin/search-backup
	upx --best --lzma --force bin/search || (mv bin/search-backup bin/search && exit 1)
	@rm bin/search-backup
	@echo "Compressed binary size:"
	@ls -lh bin/search
	@echo ""
	@echo "Testing compressed binary..."
	@bin/search --version

# Comprehensive optimization with size comparison
optimize-full:
	@echo "Running comprehensive optimization..."
	@chmod +x scripts/optimize.sh
	@VERSION=$(VERSION) scripts/optimize.sh

# Generate CPU and memory profiles
profile:
	@echo "Generating CPU profile..."
	@echo "Run: go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./..."
	@echo "Then analyze with: go tool pprof cpu.prof"
	@mkdir -p profiles
	go test -cpuprofile=profiles/cpu.prof -memprofile=profiles/mem.prof -bench=. -benchmem ./... > /dev/null 2>&1
	@echo "Profiles saved to profiles/"
	@echo "View CPU profile: go tool pprof profiles/cpu.prof"
	@echo "View memory profile: go tool pprof profiles/mem.prof"
	@echo "Compare with: go tool pprof -base profiles/old.prof profiles/new.prof"

# Run performance benchmarks
benchmark:
	@echo "Running performance benchmarks..."
	@mkdir -p profiles
	@echo "=== SearXNG Client Benchmarks ===" | tee profiles/benchmark.txt
	go test -tags=benchmark -bench=. -benchmem ./internal/searxng/... | tee -a profiles/benchmark.txt
	@echo "" | tee -a profiles/benchmark.txt
	@echo "=== Formatter Benchmarks ===" | tee -a profiles/benchmark.txt
	go test -tags=benchmark -bench=. -benchmem ./internal/formatter/... | tee -a profiles/benchmark.txt
	@echo "" | tee -a profiles/benchmark.txt
	@echo "Benchmarks complete. Results saved to profiles/benchmark.txt"

# Compare benchmark results (requires benchstat: go install golang.org/x/perf/cmd/benchstat@latest)
benchcmp:
	@echo "Comparing benchmarks..."
	@which benchstat > /dev/null 2>&1 || (echo "benchstat not found. Install with: go install golang.org/x/perf/cmd/benchstat@latest" && exit 1)
	@if [ ! -f "profiles/old.txt" ]; then \
		echo "Old benchmark file not found (profiles/old.txt)"; \
		echo "Save current benchmarks with: make benchmark > profiles/old.txt"; \
		exit 1; \
	fi
	@if [ ! -f "profiles/new.txt" ]; then \
		echo "New benchmark file not found (profiles/new.txt)"; \
		echo "Save current benchmarks with: make benchmark > profiles/new.txt"; \
		exit 1; \
	fi
	benchstat profiles/old.txt profiles/new.txt
