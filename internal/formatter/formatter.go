// Package formatter provides output formatting for search results.
//
// It supports multiple output formats including JSON, Markdown, and plain text.
// The Formatter interface allows for easy extension with new formats.
package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/mule-ai/search/internal/searxng"
)

// Formatter is the interface for all output formatters.
//
// Any formatter must implement the Format method, which takes a SearchResponse
// and returns a formatted string representation.
type Formatter interface {
	Format(result *searxng.SearchResponse) (string, error)
}

// BaseFormatter contains common formatting functionality.
//
// It provides text wrapping, truncation, and utility methods used by
// specific formatter implementations.
type BaseFormatter struct {
	Width int
}

// NewBaseFormatter creates a new base formatter with default width.
//
// The default width is 80 characters, suitable for most terminal displays.
func NewBaseFormatter() *BaseFormatter {
	return &BaseFormatter{
		Width: 80,
	}
}

// Truncate truncates a string to the specified length.
//
// If the string is shorter than length, it is returned unchanged.
// If length is 0 or negative, the string is returned unchanged.
func (f *BaseFormatter) Truncate(s string, length int) string {
	if len(s) <= length || length <= 0 {
		return s
	}
	truncated := make([]byte, length)
	copy(truncated, s)
	return string(truncated)
}

// TruncateWithEllipsis truncates a string and adds an ellipsis ("...").
//
// The total length will be at most the specified length (including the ellipsis).
// If the string fits, it is returned unchanged.
func (f *BaseFormatter) TruncateWithEllipsis(s string, length int) string {
	if len(s) <= length || length <= 0 {
		return s
	}
	if length < 4 {
		return s[:length]
	}
	return s[:length-3] + "..."
}

// MaxWidth wraps text to the specified width.
//
// It splits text into multiple lines, ensuring no line exceeds the width.
// Existing newlines are preserved. Words are not broken unless a single
// word exceeds the width.
func (f *BaseFormatter) MaxWidth(s string, width int) []string {
	if f.Width <= 0 {
		f.Width = width
	}
	lines := strings.Split(s, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) <= width {
			result = append(result, line)
			continue
		}
		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}
		currentLine := words[0]
		for _, word := range words[1:] {
			if len(currentLine)+1+len(word) <= width {
				currentLine += " " + word
			} else {
				result = append(result, currentLine)
				currentLine = word
			}
		}
		result = append(result, currentLine)
	}
	return result
}

// NewFormatter creates a formatter based on the format string.
//
// Supported formats: "json", "markdown" (or "md"), "text" (or "plaintext").
// Returns an error if the format is not recognized.
//
// Example:
//
//	f, err := formatter.NewFormatter("json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	output, err := f.Format(response)
func NewFormatter(format string) (Formatter, error) {
	return NewFormatterForCategory(format, "", false)
}

// NewFormatterForCategory creates a formatter based on format and category.
//
// This allows category-specific formatting. For example, image results get
// special treatment in markdown and text formats.
//
// The noColor flag disables colored output for text formatters.
//
// Example:
//
//	f, err := formatter.NewFormatterForCategory("markdown", "images", false)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	output, err := f.Format(response)
func NewFormatterForCategory(format string, category string, noColor bool) (Formatter, error) {
	// Check if category needs special formatting
	if searxng.NeedsSpecialFormatting(category) {
		// For image category with markdown, use special image markdown formatter
		if format == "markdown" || format == "md" {
			return &imageMarkdownFormatter{}, nil
		}
		// For image category with text, use image formatter
		if format == "text" || format == "plaintext" {
			return NewImageFormatter(), nil
		}
	}

	// Default formatting
	switch strings.ToLower(format) {
	case "json":
		return NewJSONFormatter(), nil
	case "markdown", "md":
		return NewMarkdownFormatter(), nil
	case "text", "plaintext":
		return NewTextFormatter(noColor), nil
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}

// imageMarkdownFormatter handles markdown formatting for image results.
type imageMarkdownFormatter struct{}

func (f *imageMarkdownFormatter) Format(response *searxng.SearchResponse) (string, error) {
	return FormatAsMarkdown(response)
}

// FormatResults formats the search results as a numbered list.
//
// Each result includes the title, URL, and content snippet (truncated and wrapped).
// Results are separated by blank lines.
func (f *BaseFormatter) FormatResults(results []searxng.SearchResult, width int) []string {
	var lines []string
	for i, result := range results {
		lines = append(lines, fmt.Sprintf("[%d] %s", i+1, result.Title))
		if len(result.URL) > 0 {
			lines = append(lines, fmt.Sprintf("    %s", result.URL))
		}
		if len(result.Content) > 0 {
			wrapped := f.MaxWidth(f.TruncateWithEllipsis(result.Content, width-8), width-8)
			for _, line := range wrapped {
				lines = append(lines, fmt.Sprintf("    %s", line))
			}
		}
		if i < len(results)-1 {
			lines = append(lines, "")
		}
	}
	return lines
}

// FormatResultWithSource formats a single result with source information.
//
// Includes the title, URL, source engine, score (if available), and content snippet.
// Content is wrapped to the specified width.
func (f *BaseFormatter) FormatResultWithSource(result searxng.SearchResult, index int, width int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("[%d] %s\n", index+1, result.Title))
	if len(result.URL) > 0 {
		buf.WriteString(fmt.Sprintf("    %s\n", result.URL))
	}
	if len(result.Engine) > 0 {
		buf.WriteString(fmt.Sprintf("    Source: %s", result.Engine))
		if result.Score > 0 {
			buf.WriteString(fmt.Sprintf(" | Score: %.2f", result.Score))
		}
		buf.WriteString("\n")
	}
	if len(result.Content) > 0 {
		wrapped := f.MaxWidth(f.TruncateWithEllipsis(result.Content, width-8), width-8)
		for _, line := range wrapped {
			buf.WriteString(fmt.Sprintf("    %s\n", line))
		}
	}
	return buf.String()
}

// WriteOutput writes formatted output to the specified writer.
//
// This is a convenience method for writing formatted results directly
// to an io.Writer (e.g., os.Stdout, a file, etc.).
func (f *BaseFormatter) WriteOutput(w io.Writer, format string, result *searxng.SearchResponse) error {
	formatted, err := f.FormatWithFormatter(result, format)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(formatted))
	return err
}

// FormatWithFormatter formats output using the specified format string.
//
// Supported formats: "json", "markdown", "md", "text", "plaintext".
// Returns the formatted string or an error if the format is invalid.
func (f *BaseFormatter) FormatWithFormatter(result *searxng.SearchResponse, format string) (string, error) {
	switch strings.ToLower(format) {
	case "json":
		return f.formatJSON(result)
	case "markdown", "md":
		return f.formatMarkdown(result)
	default:
		return f.formatText(result)
	}
}

// formatJSON formats results as JSON
func (f *BaseFormatter) formatJSON(result *searxng.SearchResponse) (string, error) {
	// Use the JSON formatter
	jf := NewJSONFormatter()
	return jf.Format(result)
}

// formatMarkdown formats results as Markdown
func (f *BaseFormatter) formatMarkdown(result *searxng.SearchResponse) (string, error) {
	mf := NewMarkdownFormatter()
	return mf.Format(result)
}

// formatText formats results as plain text
func (f *BaseFormatter) formatText(result *searxng.SearchResponse) (string, error) {
	tf := NewTextFormatter(false)
	return tf.Format(result)
}
