// Performance profiling tool for the search CLI.
// This script runs CPU and memory profiling on key operations.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/mule-ai/search/internal/config"
	"github.com/mule-ai/search/internal/formatter"
	"github.com/mule-ai/search/internal/searxng"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run profile.go [cpu|mem|all]")
		os.Exit(1)
	}

	profileType := os.Args[1]

	switch profileType {
	case "cpu":
		profileCPU()
	case "mem":
		profileMemory()
	case "all":
		profileCPU()
		profileMemory()
	default:
		fmt.Printf("Unknown profile type: %s\n", profileType)
		os.Exit(1)
	}
}

// profileCPU runs CPU profiling on search operations
func profileCPU() {
	f, err := os.Create("/tmp/search_cpu.prof")
	if err != nil {
		fmt.Printf("Could not create CPU profile: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("Could not start CPU profile: %v\n", err)
		os.Exit(1)
	}
	defer pprof.StopCPUProfile()

	fmt.Println("Running CPU profile...")

	// Profile JSON parsing
	fmt.Println("  Profiling JSON parsing...")
	profileJSONParsing()

	// Profile formatters
	fmt.Println("  Profiling formatters...")
	profileFormatters()

	// Profile client operations
	fmt.Println("  Profiling client operations...")
	profileClientOps()

	fmt.Println("CPU profile saved to /tmp/search_cpu.prof")
	fmt.Println("Analyze with: go tool pprof /tmp/search_cpu.prof")
}

// profileMemory runs memory profiling
func profileMemory() {
	fmt.Println("Running memory profile...")

	// Force GC before profiling
	runtime.GC()

	f, err := os.Create("/tmp/search_mem.prof")
	if err != nil {
		fmt.Printf("Could not create memory profile: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Profile JSON parsing
	fmt.Println("  Profiling JSON parsing memory...")
	profileJSONParsing()

	// Profile formatters
	fmt.Println("  Profiling formatters memory...")
	profileFormatters()

	// Profile client operations
	fmt.Println("  Profiling client operations memory...")
	profileClientOps()

	// Force GC and write heap profile
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Printf("Could not write memory profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Memory profile saved to /tmp/search_mem.prof")
	fmt.Println("Analyze with: go tool pprof /tmp/search_mem.prof")
}

// profileJSONParsing profiles JSON decoding performance
func profileJSONParsing() {
	// Create a large mock response
	results := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		results[i] = map[string]interface{}{
			"title":       fmt.Sprintf("Result %d", i),
			"url":         fmt.Sprintf("https://example.com/%d", i),
			"content":     "Example content for testing performance",
			"engine":      "google",
			"category":    "general",
			"score":       0.9,
			"parsed_url":  []string{fmt.Sprintf("https://example.com/%d", i)},
			"template":    "default.html",
			"engines":     []string{"google"},
		}
	}

	mockData := map[string]interface{}{
		"query":            "test query",
		"results":          results,
		"answers":          []interface{}{},
		"infoboxes":        []interface{}{},
		"suggestions":      []string{},
		"number_of_results": 10000,
	}

	mockJSON, _ := json.Marshal(mockData)

	// Run multiple iterations for better profiling
	for i := 0; i < 100; i++ {
		var response searxng.SearchResponse
		if err := json.Unmarshal(mockJSON, &response); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// profileFormatters profiles formatter performance
func profileFormatters() {
	resp := createMockResponse(50)

	// Profile JSON formatter
	for i := 0; i < 100; i++ {
		jf := formatter.NewJSONFormatter()
		_, _ = jf.Format(resp)
	}

	// Profile Markdown formatter
	for i := 0; i < 100; i++ {
		mf := formatter.NewMarkdownFormatter()
		_, _ = mf.Format(resp)
	}

	// Profile Text formatter
	for i := 0; i < 100; i++ {
		tf := formatter.NewTextFormatter(false)
		_, _ = tf.Format(resp)
	}
}

// profileClientOps profiles client operations
func profileClientOps() {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}

	// Profile client creation
	for i := 0; i < 100; i++ {
		_ = searxng.NewClient(cfg)
	}

	// Profile request creation
	for i := 0; i < 100; i++ {
		_ = searxng.NewSearchRequest("test query")
	}

	// Profile request with options
	for i := 0; i < 100; i++ {
		req := searxng.NewSearchFromQuery("test query",
			searxng.WithPage(2),
			searxng.WithSafeSearch(1),
			searxng.WithTimeRange("week"),
		)
		_ = req
	}
}

// createMockResponse creates a mock search response
func createMockResponse(numResults int) *searxng.SearchResponse {
	results := make([]searxng.SearchResult, numResults)
	for i := 0; i < numResults; i++ {
		results[i] = searxng.SearchResult{
			Title:    fmt.Sprintf("Test Result %d", i),
			URL:      fmt.Sprintf("https://example.com/result/%d", i),
			Content:  "This is test content for the search result",
			Engine:   "google",
			Category: "general",
			Score:    0.9 - float64(i)*0.01,
		}
	}

	return &searxng.SearchResponse{
		Query:           "test query",
		Results:         results,
		NumberOfResults: numResults * 100,
		SearchTime:      0.25,
		Page:            1,
		Instance:        "https://search.butler.ooo",
	}
}

// init adds environment info to profile
func init() {
	fmt.Printf("Profiling tool for search CLI\n")
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Started at: %s\n\n", time.Now().Format(time.RFC3339))
}