// Package main is the entry point for the search CLI application.
//
// It initializes the CLI command and executes it, handling any errors
// by printing them to stderr and exiting with a non-zero status code.
package main

import (
	"fmt"
	"os"

	"github.com/mule-ai/search/cmd/search/cli"
)

// main is the application entry point.
//
// It executes the CLI and exits with status 1 if an error occurs.
func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
