// Package formatter provides optimized output formatting for search results.
//
// The optimized formatters use buffer pooling and other performance optimizations
// to reduce memory allocations and improve formatting speed for large result sets.
package formatter

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/mule-ai/search/internal/searxng"
)

// OptimizedFormatter provides high-performance formatting using buffer pools.
//
// This formatter is optimized for speed and memory efficiency when handling
// large result sets by reusing buffers and minimizing allocations.
type OptimizedFormatter struct {
	format string
}

// NewOptimizedFormatter creates an optimized formatter for the specified format.
//
// Supported formats are "json", "markdown", "md", and "text". The optimized
// formatter uses buffer pooling to reduce memory allocations.
//
// Example:
//
//	optFormatter := formatter.NewOptimizedFormatter("json")
//	output, err := optFormatter.Format(searchResponse)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func NewOptimizedFormatter(format string) *OptimizedFormatter {
	return &OptimizedFormatter{format: format}
}

// Format formats the search response using buffer pooling for efficiency.
//
// This method is optimized for performance by reusing buffers from a pool.
// It automatically selects the appropriate format method based on the
// formatter's format setting.
//
// Example:
//
//	optFormatter := formatter.NewOptimizedFormatter("text")
//	response := &searxng.SearchResponse{
//	    Query: "golang",
//	    Results: []searxng.SearchResult{...},
//	}
//	output, err := optFormatter.Format(response)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func (f *OptimizedFormatter) Format(result *searxng.SearchResponse) (string, error) {
	buf := searxng.GetBuffer()
	defer searxng.PutBuffer(buf)

	switch strings.ToLower(f.format) {
	case "json":
		return f.formatJSON(result, buf)
	case "markdown", "md":
		return f.formatMarkdown(result, buf)
	default:
		return f.formatText(result, buf)
	}
}

// formatJSON formats results as JSON using the buffer pool.
func (f *OptimizedFormatter) formatJSON(result *searxng.SearchResponse, buf *bytes.Buffer) (string, error) {
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(result); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// formatMarkdown formats results as Markdown using the buffer pool.
func (f *OptimizedFormatter) formatMarkdown(result *searxng.SearchResponse, buf *bytes.Buffer) (string, error) {
	// Write header
	buf.WriteString("# Search Results: ")
	buf.WriteString(result.Query)
	buf.WriteString("\n\n")

	// Write result count
	buf.WriteString("Found **")
	buf.WriteString(formatNumber(result.NumberOfResults))
	buf.WriteString("** results")

	if result.SearchTime > 0 {
		buf.WriteString(" in ")
		buf.WriteString(formatDuration(result.SearchTime))
	}

	buf.WriteString("\n\n")

	// Write results
	for i, r := range result.Results {
		// Title as link
		buf.WriteString("## [")
		buf.WriteString(escapeMarkdown(r.Title))
		buf.WriteString("](")
		buf.WriteString(r.URL)
		buf.WriteString(")\n")

		// Source and score
		if r.Engine != "" || r.Score > 0 {
			buf.WriteString("**Source:** ")
			if r.Engine != "" {
				buf.WriteString(capitalize(r.Engine))
			}
			if r.Score > 0 {
				buf.WriteString(" | **Score:** ")
				buf.WriteString(formatFloat(r.Score))
			}
			buf.WriteString("\n\n")
		}

		// Content
		if r.Content != "" {
			buf.WriteString(escapeMarkdown(truncateString(r.Content, 500)))
			buf.WriteString("\n\n")
		}

		if i < len(result.Results)-1 {
			buf.WriteString("---\n\n")
		}
	}

	return buf.String(), nil
}

// formatText formats results as plain text using the buffer pool.
func (f *OptimizedFormatter) formatText(result *searxng.SearchResponse, buf *bytes.Buffer) (string, error) {
	// Write header
	buf.WriteString(result.Query)
	buf.WriteString("\n")
	buf.WriteString(strings.Repeat("=", len(result.Query)))
	buf.WriteString("\n\n")

	// Write results
	for i, r := range result.Results {
		buf.WriteString("[")
		buf.WriteString(formatInt(i + 1))
		buf.WriteString("] ")
		buf.WriteString(r.Title)
		buf.WriteString("\n")

		if r.URL != "" {
			buf.WriteString("    ")
			buf.WriteString(r.URL)
			buf.WriteString("\n")
		}

		if r.Engine != "" || r.Score > 0 {
			buf.WriteString("    Source: ")
			buf.WriteString(capitalize(r.Engine))
			if r.Score > 0 {
				buf.WriteString(" | Score: ")
				buf.WriteString(formatFloat(r.Score))
			}
			buf.WriteString("\n")
		}

		if r.Content != "" {
			lines := wrapText(truncateString(r.Content, 300), 76)
			for _, line := range lines {
				buf.WriteString("    ")
				buf.WriteString(line)
				buf.WriteString("\n")
			}
		}

		if i < len(result.Results)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String(), nil
}

// Helper functions for optimized formatting.

func formatNumber(n int) string {
	if n < 1000 {
		return formatInt(n)
	}
	if n < 1000000 {
		return formatInt(n/1000) + "," + formatInt(n%1000)
	}
	s := formatInt(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	var sb strings.Builder
	sb.Grow(12)
	sb.WriteString(strconv.Itoa(n))
	return sb.String()
}

func formatFloat(f float64) string {
	var sb strings.Builder
	sb.Grow(16)
	sb.WriteString(strconv.FormatFloat(f, 'f', 2, 64))
	return sb.String()
}

func formatDuration(d float64) string {
	// Simplified - format as seconds
	return "0.00s" // Placeholder
}

func escapeMarkdown(s string) string {
	// Basic markdown escaping
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	return s
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	// Optimize: avoid strings.ToUpper allocation for single char
	if s[0] >= 'a' && s[0] <= 'z' {
		// Convert to uppercase by subtracting 32
		return string(s[0]-32) + s[1:]
	}
	return s
}

func wrapText(s string, width int) []string {
	if len(s) <= width {
		return []string{s}
	}
	// Optimized word wrapping with pre-allocation
	words := strings.Fields(s)
	// Pre-allocate capacity to avoid re-allocations (estimate: 1 line per 4 words)
	estimatedLines := (len(words) + 3) / 4
	lines := make([]string, 0, estimatedLines)

	var sb strings.Builder
	sb.Grow(width) // Pre-grow for first line

	for i, word := range words {
		wordLen := len(word)
		currentLen := sb.Len()

		// Check if we can add this word to the current line
		if currentLen == 0 {
			// First word in line
			sb.WriteString(word)
		} else if currentLen+1+wordLen <= width {
			// Add to existing line
			sb.WriteString(" ")
			sb.WriteString(word)
		} else {
			// Start new line
			lines = append(lines, sb.String())
			sb.Reset()
			sb.Grow(width)
			sb.WriteString(word)
		}

		// If this is the last word and we have content, add it
		if i == len(words)-1 && sb.Len() > 0 {
			lines = append(lines, sb.String())
		}
	}

	return lines
}
