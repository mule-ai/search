// Package profiling provides benchmarking tests for performance optimization.
package profiling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mule-ai/search/internal/searxng"
)

// mockSearchResponse creates a mock response for benchmarking.
func mockSearchResponse(numResults int) *searxng.SearchResponse {
	results := make([]searxng.SearchResult, numResults)
	
	for i := 0; i < numResults; i++ {
		results[i] = searxng.SearchResult{
			Title:   fmt.Sprintf("Result Title %d", i),
			URL:     fmt.Sprintf("https://example.com/result/%d", i),
			Content: fmt.Sprintf("This is the content for result %d. It contains some text that would typically appear in a search result snippet.", i),
			Engine:  "google",
			Score:   float64(1.0 - float64(i)*0.01),
			Category: "general",
		}
	}
	
	return &searxng.SearchResponse{
		Query:            "test query",
		Results:          results,
		NumberOfResults:  numResults * 1000,
		Answers:          []searxng.Answer{},
		Infoboxes:        []searxng.Infobox{},
		Suggestions:      []string{},
	}
}

// BenchmarkJSONParsingSmall benchmarks JSON parsing for small result sets (10 results).
func BenchmarkJSONParsingSmall(b *testing.B) {
	resp := mockSearchResponse(10)
	data, _ := json.Marshal(resp)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result searxng.SearchResponse
		if err := json.Unmarshal(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJSONParsingMedium benchmarks JSON parsing for medium result sets (50 results).
func BenchmarkJSONParsingMedium(b *testing.B) {
	resp := mockSearchResponse(50)
	data, _ := json.Marshal(resp)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result searxng.SearchResponse
		if err := json.Unmarshal(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJSONParsingLarge benchmarks JSON parsing for large result sets (100 results).
func BenchmarkJSONParsingLarge(b *testing.B) {
	resp := mockSearchResponse(100)
	data, _ := json.Marshal(resp)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result searxng.SearchResponse
		if err := json.Unmarshal(data, &result); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJSONMarshal benchmarks JSON marshaling.
func BenchmarkJSONMarshal(b *testing.B) {
	resp := mockSearchResponse(50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(resp); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBufferPool benchmarks buffer pool performance.
func BenchmarkBufferPool(b *testing.B) {
	resp := mockSearchResponse(50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(resp); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkWrite benchmarks writing large output.
func BenchmarkWrite(b *testing.B) {
	resp := mockSearchResponse(100)
	data, _ := json.Marshal(resp)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := os.Stdout.Write(data); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStringFormatting benchmarks string formatting operations.
func BenchmarkStringFormatting(b *testing.B) {
	resp := mockSearchResponse(50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range resp.Results {
			_ = fmt.Sprintf("%s\n%s\n\n%s\n---", r.Title, r.URL, r.Content)
		}
	}
}

// BenchmarkStringsBuilder benchmarks using strings.Builder for output.
func BenchmarkStringsBuilder(b *testing.B) {
	resp := mockSearchResponse(50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var builder bytes.Buffer
		for _, r := range resp.Results {
			builder.WriteString(r.Title)
			builder.WriteByte('\n')
			builder.WriteString(r.URL)
			builder.WriteByte('\n')
			builder.WriteByte('\n')
			builder.WriteString(r.Content)
			builder.WriteString("\n---\n")
		}
		_ = builder.String()
	}
}

// BenchmarkTimeOperations benchmarks time operations for search tracking.
func BenchmarkTimeOperations(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = time.Now()
		_ = time.Since(time.Now().Add(-time.Second))
	}
}