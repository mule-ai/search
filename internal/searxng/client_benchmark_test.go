//go:build benchmark
// +build benchmark

package searxng

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mule-ai/search/internal/config"
)

// BenchmarkSearchRequestCreation benchmarks creating search requests
func BenchmarkSearchRequestCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewSearchRequest("golang programming")
	}
}

// BenchmarkJSONParsing benchmarks parsing SearXNG JSON responses
func BenchmarkJSONParsing(b *testing.B) {
	// Create a realistic mock response
	mockResponse := `{
		"query": "golang",
		"results": [
			{
				"title": "A Tour of Go",
				"url": "https://go.dev/tour/",
				"content": "Welcome to a tour of the Go programming language. The Tour is divided into a list of modules that you can access by clicking on A Tour of Go on the left side of the page.",
				"engine": "google",
				"parsed_url": ["https://go.dev/tour/"],
				"template": "default.html",
				"engines": ["google"]
			},
			{
				"title": "The Go Programming Language",
				"url": "https://go.dev/",
				"content": "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
				"engine": "duckduckgo",
				"parsed_url": ["https://go.dev/"],
				"template": "default.html",
				"engines": ["duckduckgo"]
			},
			{
				"title": "Go by Example",
				"url": "https://gobyexample.com/",
				"content": "Go by Example is a hands-on introduction to Go using annotated, example programs. Check out the first example or browse the full list below.",
				"engine": "bing",
				"parsed_url": ["https://gobyexample.com/"],
				"template": "default.html",
				"engines": ["bing"]
			}
		],
		"answers": [],
		"infoboxes": [],
		"suggestions": ["golang tutorial", "golang download", "golang vs python"],
		"number_of_results": 1250000
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var response SearchResponse
		if err := json.Unmarshal([]byte(mockResponse), &response); err != nil {
			b.Fatalf("Failed to parse JSON: %v", err)
		}
	}
}

// BenchmarkJSONParsingLarge benchmarks parsing large JSON responses
func BenchmarkJSONParsingLarge(b *testing.B) {
	// Create a large mock response with 50 results
	results := make([]map[string]interface{}, 50)
	for i := 0; i < 50; i++ {
		results[i] = map[string]interface{}{
			"title":       "Example Result " + string(rune('0'+i%10)),
			"url":         "https://example.com/result/" + string(rune('0'+i%10)),
			"content":     "This is example content for result number " + string(rune('0'+i%10)) + ". It contains some text to simulate a real search result.",
			"engine":      "google",
			"parsed_url":  []string{"https://example.com/result/" + string(rune('0'+i%10))},
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
		"number_of_results": 5000,
	}

	mockJSON, _ := json.Marshal(mockData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var response SearchResponse
		if err := json.Unmarshal(mockJSON, &response); err != nil {
			b.Fatalf("Failed to parse JSON: %v", err)
		}
	}
}

// BenchmarkNewClient benchmarks creating new clients
func BenchmarkNewClient(b *testing.B) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewClient(cfg)
	}
}

// BenchmarkNewClientWithTimeout benchmarks creating clients with custom timeout
func BenchmarkNewClientWithTimeout(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewClientWithTimeout("https://search.butler.ooo", 30*time.Second)
	}
}

// BenchmarkResultProcessing benchmarks processing search results
func BenchmarkResultProcessing(b *testing.B) {
	// Create mock results
	results := make([]SearchResult, 20)
	for i := 0; i < 20; i++ {
		results[i] = SearchResult{
			Title:    "Example Result",
			URL:      "https://example.com",
			Content:  "Example content",
			Engine:   "google",
			Category: "general",
			Score:    0.9,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate result processing
		for _, result := range results {
			_ = result.Title
			_ = result.URL
			_ = result.Content
			_ = result.Engine
		}
	}
}

// BenchmarkFilterResults benchmarks filtering results by criteria
func BenchmarkFilterResults(b *testing.B) {
	results := make([]SearchResult, 50)
	for i := 0; i < 50; i++ {
		engine := "google"
		if i%3 == 0 {
			engine = "duckduckgo"
		} else if i%3 == 1 {
			engine = "bing"
		}
		results[i] = SearchResult{
			Title:    "Example Result",
			URL:      "https://example.com",
			Content:  "Example content",
			Engine:   engine,
			Category: "general",
			Score:    float64(i) / 50.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Filter results by engine
		filtered := make([]SearchResult, 0)
		for _, result := range results {
			if result.Engine == "google" {
				filtered = append(filtered, result)
			}
		}
		_ = filtered
	}
}

// BenchmarkSortResults benchmarks sorting results by score
func BenchmarkSortResults(b *testing.B) {
	results := make([]SearchResult, 50)
	for i := 0; i < 50; i++ {
		results[i] = SearchResult{
			Title:    "Example Result",
			URL:      "https://example.com",
			Content:  "Example content",
			Engine:   "google",
			Category: "general",
			Score:    float64(i) / 50.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a copy for sorting
		sorted := make([]SearchResult, len(results))
		copy(sorted, results)

		// Simple bubble sort for benchmarking
		for j := 0; j < len(sorted)-1; j++ {
			for k := 0; k < len(sorted)-j-1; k++ {
				if sorted[k].Score < sorted[k+1].Score {
					sorted[k], sorted[k+1] = sorted[k+1], sorted[k]
				}
			}
		}
	}
}

// BenchmarkResponseValidation benchmarks validating search responses
func BenchmarkResponseValidation(b *testing.B) {
	response := &SearchResponse{
		Query: "golang",
		Results: []SearchResult{
			{
				Title:    "Example",
				URL:      "https://example.com",
				Content:  "Content",
				Engine:   "google",
				Category: "general",
			},
		},
		NumberOfResults: 1000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Validate response
		if response.Query == "" {
			b.Fatal("Query should not be empty")
		}
		if response.Results == nil {
			b.Fatal("Results should not be nil")
		}
	}
}

// BenchmarkCategoryMapping benchmarks category name mapping
func BenchmarkCategoryMapping(b *testing.B) {
	categories := []string{"general", "images", "videos", "news", "music", "it", "science", "files"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Map categories (simulated)
		for _, cat := range categories {
			_ = cat // Simulate processing
		}
	}
}

// BenchmarkSearchRequestOptions benchmarks creating search requests with various options
func BenchmarkSearchRequestOptions(b *testing.B) {
	options := []*SearchRequest{
		NewSearchRequest("golang"),
		{
			Query:      "rust programming",
			Format:     "json",
			Page:       2,
			Languages:  []string{"en", "de"},
			SafeSearch: 1,
			Categories: []string{"general", "it"},
		},
		{
			Query:      "kubernetes",
			Format:     "json",
			Page:       1,
			TimeRange:  "week",
			SafeSearch: 0,
			Categories: []string{"general"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, req := range options {
			_ = req.Query
			_ = req.Format
			_ = req.Categories
		}
	}
}

// BenchmarkRequestWithAllOptions benchmarks search request with all options set
func BenchmarkRequestWithAllOptions(b *testing.B) {
	req := &SearchRequest{
		Query:      "complex search query",
		Format:     "json",
		Page:       3,
		Languages:  []string{"en", "de", "fr"},
		SafeSearch: 2,
		Categories: []string{"general", "it", "science"},
		Engines:    []string{"google", "duckduckgo"},
		TimeRange:  "month",
		Timeout:    60 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Access all fields
		_ = req.Query
		_ = req.Format
		_ = req.Page
		_ = req.Languages
		_ = req.SafeSearch
		_ = req.Categories
		_ = req.Engines
		_ = req.TimeRange
		_ = req.Timeout
	}
}

// BenchmarkInfoboxParsing benchmarks parsing infobox data
func BenchmarkInfoboxParsing(b *testing.B) {
	infobox := Infobox{
		Infobox:  "Example Infobox",
		Content:  "Infobox content",
		Engine:   "wikipedia",
		URLs:     []URLInfo{{Title: "Example", URL: "https://example.com", Official: true}},
		Template: "default.html",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = infobox.Infobox
		_ = infobox.Content
		_ = len(infobox.URLs)
		_ = infobox.Engine
	}
}

// BenchmarkAnswerParsing benchmarks parsing answer data
func BenchmarkAnswerParsing(b *testing.B) {
	answer := Answer{
		Answer:    "42",
		URL:       "https://example.com/answer",
		Engine:    "wikipedia",
		ParsedURL: []string{"https://example.com/answer"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = answer.Answer
		_ = answer.URL
		_ = answer.Engine
	}
}

// BenchmarkSuggestionProcessing benchmarks processing suggestions
func BenchmarkSuggestionProcessing(b *testing.B) {
	suggestions := []string{
		"golang tutorial",
		"golang download",
		"golang vs python",
		"golang jobs",
		"golang framework",
		"golang course",
		"golang book",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Process suggestions
		for _, suggestion := range suggestions {
			_ = suggestion
		}
	}
}

// BenchmarkScoreCalculation benchmarks score-related operations
func BenchmarkScoreCalculation(b *testing.B) {
	results := make([]SearchResult, 100)
	for i := 0; i < 100; i++ {
		results[i] = SearchResult{
			Title:    "Result",
			URL:      "https://example.com",
			Score:    float64(i) / 100.0,
			Engine:   "google",
			Category: "general",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		totalScore := 0.0
		for _, result := range results {
			totalScore += result.Score
		}
		_ = totalScore
	}
}

// BenchmarkURLParsing benchmarks URL field processing
func BenchmarkURLParsing(b *testing.B) {
	result := SearchResult{
		URL:       "https://example.com/path/to/page?query=value",
		ParsedURL: []string{"https://example.com", "/path/to/page", "query=value"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.URL
		_ = len(result.ParsedURL)
	}
}

// BenchmarkEngineListProcessing benchmarks engine list operations
func BenchmarkEngineListProcessing(b *testing.B) {
	engines := []string{"google", "duckduckgo", "bing", "brave", "wikipedia"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Join engines (simulated)
		result := ""
		for j, engine := range engines {
			if j > 0 {
				result += ","
			}
			result += engine
		}
		_ = result
	}
}

// BenchmarkTimeRangeOptions benchmarks time range processing
func BenchmarkTimeRangeOptions(b *testing.B) {
	timeRanges := []string{"day", "week", "month", "year"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tr := range timeRanges {
			_ = tr // Simulate validation/processing
		}
	}
}