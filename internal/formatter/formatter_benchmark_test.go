package formatter

import (
	"testing"

	"github.com/mule-ai/search/internal/searxng"
)

// BenchmarkJSONFormatter benchmarks the JSON formatter.
func BenchmarkJSONFormatter(b *testing.B) {
	resp := createMockResponse(50)
	f := NewJSONFormatter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Format(resp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarkdownFormatter benchmarks the Markdown formatter.
func BenchmarkMarkdownFormatter(b *testing.B) {
	resp := createMockResponse(50)
	f := NewMarkdownFormatter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Format(resp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTextFormatter benchmarks the text formatter.
func BenchmarkTextFormatter(b *testing.B) {
	resp := createMockResponse(50)
	f := NewTextFormatter(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Format(resp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOptimizedJSON benchmarks the optimized JSON formatter.
func BenchmarkOptimizedJSON(b *testing.B) {
	resp := createMockResponse(50)
	f := NewOptimizedFormatter("json")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Format(resp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOptimizedMarkdown benchmarks the optimized Markdown formatter.
func BenchmarkOptimizedMarkdown(b *testing.B) {
	resp := createMockResponse(50)
	f := NewOptimizedFormatter("markdown")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Format(resp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOptimizedText benchmarks the optimized text formatter.
func BenchmarkOptimizedText(b *testing.B) {
	resp := createMockResponse(50)
	f := NewOptimizedFormatter("text")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := f.Format(resp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFormatterSmall benchmarks all formatters with small result sets.
func BenchmarkFormatterSmall(b *testing.B) {
	resp := createMockResponse(10)

	b.Run("JSON", func(b *testing.B) {
		f := NewJSONFormatter()
		for i := 0; i < b.N; i++ {
			_, _ = f.Format(resp)
		}
	})

	b.Run("Markdown", func(b *testing.B) {
		f := NewMarkdownFormatter()
		for i := 0; i < b.N; i++ {
			_, _ = f.Format(resp)
		}
	})

	b.Run("Text", func(b *testing.B) {
		f := NewTextFormatter(false)
		for i := 0; i < b.N; i++ {
			_, _ = f.Format(resp)
		}
	})
}

// BenchmarkFormatterLarge benchmarks all formatters with large result sets.
func BenchmarkFormatterLarge(b *testing.B) {
	resp := createMockResponse(100)

	b.Run("JSON", func(b *testing.B) {
		f := NewJSONFormatter()
		for i := 0; i < b.N; i++ {
			_, _ = f.Format(resp)
		}
	})

	b.Run("Markdown", func(b *testing.B) {
		f := NewMarkdownFormatter()
		for i := 0; i < b.N; i++ {
			_, _ = f.Format(resp)
		}
	})

	b.Run("Text", func(b *testing.B) {
		f := NewTextFormatter(false)
		for i := 0; i < b.N; i++ {
			_, _ = f.Format(resp)
		}
	})
}

// createMockResponse creates a mock search response for benchmarking.
func createMockResponse(numResults int) *searxng.SearchResponse {
	results := make([]searxng.SearchResult, numResults)

	for i := 0; i < numResults; i++ {
		results[i] = searxng.SearchResult{
			Title:   "Test Result Title for Benchmarking Performance",
			URL:     "https://example.com/test-result-url",
			Content: "This is a test content snippet that represents what would be returned from a SearXNG search instance. It contains enough text to test formatting performance.",
			Engine:  "google",
			Score:   0.95 - float64(i)*0.01,
			Category: "general",
		}
	}

	return &searxng.SearchResponse{
		Query:           "benchmark test query",
		Results:         results,
		NumberOfResults: numResults * 1000,
		SearchTime:      0.25,
		Page:            1,
		Instance:        "https://search.butler.ooo",
	}
}
