// Package version provides version information for the search CLI.
//
// Version variables can be overridden at build time using ldflags:
//
//	go build -ldflags "-X github.com/mule-ai/search/pkg/version.Version=1.0.0 \
//	                   -X github.com/mule-ai/search/pkg/version.GitCommit=abc123 \
//	                   -X github.com/mule-ai/search/pkg/version.BuildDate=2024-01-01"
package version

// Version is the current version of search CLI.
//
// This can be overridden at build time using ldflags.
// Default: "1.0.0"
var Version = "1.0.0"

// GitCommit stores the git commit hash.
//
// Set at build time using ldflags.
// Default: "unknown"
var GitCommit = "unknown"

// BuildDate stores the build timestamp.
//
// Set at build time using ldflags.
// Default: "unknown"
var BuildDate = "unknown"

// GoVersion stores the Go version used to build.
//
// Set at build time using ldflags.
// Default: "unknown"
var GoVersion = "unknown"
