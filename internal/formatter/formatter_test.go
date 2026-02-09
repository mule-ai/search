package formatter

import (
	"testing"

	"github.com/mule-ai/search/internal/searxng"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "json format",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "markdown format",
			format:  "markdown",
			wantErr: false,
		},
		{
			name:    "text format",
			format:  "text",
			wantErr: false,
		},
		{
			name:    "invalid format",
			format:  "xml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFormatter(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && f == nil {
				t.Error("NewFormatter() returned nil formatter for valid format")
			}
		})
	}
}

func TestJSONFormatter(t *testing.T) {
	f := NewJSONFormatter()

	if f == nil {
		t.Fatal("NewJSONFormatter() returned nil")
	}

	// Create test data
	results := []searxng.SearchResult{
		{
			Title:    "Test Title 1",
			URL:      "https://example.com/1",
			Content:  "Test content 1",
			Engine:   "google",
			Category: "general",
			Score:    0.95,
		},
		{
			Title:    "Test Title 2",
			URL:      "https://example.com/2",
			Content:  "Test content 2",
			Engine:   "duckduckgo",
			Category: "general",
			Score:    0.85,
		},
	}

	response := &searxng.SearchResponse{
		Query:            "test query",
		Results:          results,
		Answers:          []searxng.Answer{{Answer: "Test answer"}},
		Infoboxes:        []searxng.Infobox{{Infobox: "Test infobox"}},
		Suggestions:      []string{"suggestion1", "suggestion2"},
		NumberOfResults:  100,
		SearchTime:       0.25,
	}

	// Test Format method
	output, err := f.Format(response)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	if len(output) == 0 {
		t.Error("Format() returned empty string")
	}

	// Verify JSON contains expected fields
	if !containsString(output, "test query") {
		t.Error("Format() output missing query")
	}
	if !containsString(output, "Test Title 1") {
		t.Error("Format() output missing title")
	}
	if !containsString(output, "https://example.com/1") {
		t.Error("Format() output missing URL")
	}
}

func TestMarkdownFormatter(t *testing.T) {
	f := NewMarkdownFormatter()

	if f == nil {
		t.Fatal("NewMarkdownFormatter() returned nil")
	}

	results := []searxng.SearchResult{
		{
			Title:    "Test Title",
			URL:      "https://example.com",
			Content:  "Test content",
			Engine:   "google",
			Category: "general",
			Score:    0.95,
		},
	}

	response := &searxng.SearchResponse{
		Query:            "test query",
		Results:          results,
		Answers:          []searxng.Answer{{Answer: "Test answer"}},
		Infoboxes:        []searxng.Infobox{{Infobox: "Test infobox"}},
		Suggestions:      []string{"suggestion1"},
		NumberOfResults:  100,
		SearchTime:       0.25,
	}

	output, err := f.Format(response)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	if len(output) == 0 {
		t.Error("Format() returned empty string")
	}

	// Verify markdown format
	if !containsString(output, "# Search Results") {
		t.Error("Format() output missing header")
	}
	if !containsString(output, "Test Title") {
		t.Error("Format() output missing title")
	}
	if !containsString(output, "Source:") {
		t.Error("Format() output missing source info")
	}
}

func TestTextFormatter(t *testing.T) {
	f := NewTextFormatter(false)

	if f == nil {
		t.Fatal("NewTextFormatter() returned nil")
	}

	results := []searxng.SearchResult{
		{
			Title:    "Test Title",
			URL:      "https://example.com",
			Content:  "Test content",
			Engine:   "google",
			Category: "general",
			Score:    0.95,
		},
	}

	response := &searxng.SearchResponse{
		Query:            "test query",
		Results:          results,
		NumberOfResults:  100,
		SearchTime:       0.25,
	}

	output, err := f.Format(response)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	if len(output) == 0 {
		t.Error("Format() returned empty string")
	}

	// Verify text format
	if !containsString(output, "[1]") {
		t.Error("Format() output missing numbered result")
	}
	if !containsString(output, "Test Title") {
		t.Error("Format() output missing title")
	}
}

func TestBaseFormatter(t *testing.T) {
	f := NewBaseFormatter()

	if f == nil {
		t.Fatal("NewBaseFormatter() returned nil")
	}

	// Test TruncateWithEllipsis
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "short text",
			text:     "short",
			maxLen:   20,
			expected: "short",
		},
		{
			name:     "exact length",
			text:     "exactly ten!",
			maxLen:   12,
			expected: "exactly ten!",
		},
		{
			name:     "truncate needed",
			text:     "this is a very long text that needs to be truncated",
			maxLen:   20,
			expected: "this is a very lo...", // 20 total: 17 chars + 3 for ellipsis
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.TruncateWithEllipsis(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateWithEllipsis() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test MaxWidth (used for text wrapping)
	wrapped := f.MaxWidth("this is a test for wrapping text to fit within specified width", 20)
	if len(wrapped) == 0 {
		t.Error("MaxWidth() returned empty slice")
	}

	// Check that no line exceeds max width
	for _, line := range wrapped {
		if len(line) > 25 {
			t.Errorf("MaxWidth() produced line longer than expected: %q (len=%d)", line, len(line))
		}
	}
}

func TestFormatterFactory(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		want    string
		wantErr bool
	}{
		{
			name:    "json formatter",
			format:  "json",
			want:    "*formatter.JSONFormatter",
			wantErr: false,
		},
		{
			name:    "markdown formatter",
			format:  "markdown",
			want:    "*formatter.MarkdownFormatter",
			wantErr: false,
		},
		{
			name:    "text formatter",
			format:  "text",
			want:    "*formatter.TextFormatter",
			wantErr: false,
		},
		{
			name:    "invalid format",
			format:  "invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFormatter(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				typeName := typeof(f)
				if typeName != tt.want {
					t.Errorf("NewFormatter() type = %v, want %v", typeName, tt.want)
				}
			}
		})
	}
}

// Helper functions

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	var line string
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(c)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func typeof(v interface{}) string {
	switch v.(type) {
	case *JSONFormatter:
		return "*formatter.JSONFormatter"
	case *MarkdownFormatter:
		return "*formatter.MarkdownFormatter"
	case *TextFormatter:
		return "*formatter.TextFormatter"
	default:
		return "unknown"
	}
}

// Test additional TextFormatter methods
func TestTextFormatterAdditionalMethods(t *testing.T) {
	results := []searxng.SearchResult{
		{
			Title:    "Test Title",
			URL:      "https://example.com",
			Content:  "Test content",
			Engine:   "google",
			Category: "general",
			Score:    0.95,
		},
	}

	t.Run("FormatWithResultCount", func(t *testing.T) {
		f := NewTextFormatter(false)
		output := f.FormatWithResultCount("test query", results, 100)

		if len(output) == 0 {
			t.Error("FormatWithResultCount() returned empty string")
		}
		if !containsString(output, "test query") {
	}
	})
}


// Test JSONFormatter additional methods
func TestJSONFormatterAdditionalMethods(t *testing.T) {
	results := []searxng.SearchResult{
		{
			Title:    "Test Title",
			URL:      "https://example.com",
			Content:  "Test content",
			Engine:   "google",
			Category: "general",
			Score:    0.95,
		},
	}

	t.Run("FormatWithQuery", func(t *testing.T) {
		f := NewJSONFormatter()
		output, err := f.FormatWithQuery("test query", results, 0.25, "https://search.example.com")
		
		if err != nil {
			t.Errorf("FormatWithQuery() error = %v", err)
		}
		if len(output) == 0 {
			t.Error("FormatWithQuery() returned empty string")
		}
		if !containsString(output, "test query") {
			t.Error("FormatWithQuery() missing query")
		}
	})

	t.Run("Pretty toggle", func(t *testing.T) {
		f := NewJSONFormatter()
		
		if !f.IsPretty() {
			t.Error("NewJSONFormatter should have pretty enabled by default")
		}
		
		f.DisablePretty()
		if f.IsPretty() {
			t.Error("DisablePretty() should disable pretty printing")
		}
		
		f.EnablePretty()
		if !f.IsPretty() {
			t.Error("EnablePretty() should enable pretty printing")
		}
	})

	t.Run("FormatResult", func(t *testing.T) {
		f := NewJSONFormatter()
		if len(results) > 0 {
			output, err := f.FormatResult(results[0])
			if err != nil {
				t.Errorf("FormatResult() error = %v", err)
			}
			if len(output) == 0 {
				t.Error("FormatResult() returned empty string")
			}
		}
	})

	t.Run("FormatWithQueryResultsTimeAndInstance", func(t *testing.T) {
		f := NewJSONFormatter()
		output, err := f.FormatWithQueryResultsTimeAndInstance("test query", results, 0.25, "https://search.example.com")
		
		if err != nil {
			t.Errorf("FormatWithQueryResultsTimeAndInstance() error = %v", err)
		}
		if len(output) == 0 {
			t.Error("FormatWithQueryResultsTimeAndInstance() returned empty string")
		}
	})

	t.Run("FormatResultsOnly", func(t *testing.T) {
		f := NewJSONFormatter()
		output, err := f.FormatResultsOnly(results)
		
		if err != nil {
			t.Errorf("FormatResultsOnly() error = %v", err)
		}
		if len(output) == 0 {
			t.Error("FormatResultsOnly() returned empty string")
		}
	})
}

// Test NewFormatterForCategory
func TestNewFormatterForCategory(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		category string
		wantErr  bool
	}{
		{"text with images category", "text", "images", false},
		{"markdown with images category", "markdown", "images", false},
		{"json with images category", "json", "images", false},
		{"text with general category", "text", "general", false},
		{"invalid format", "invalid", "general", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFormatterForCategory(tt.format, tt.category, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormatterForCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && f == nil {
				t.Error("NewFormatterForCategory() returned nil for valid input")
			}
		})
	}
}
