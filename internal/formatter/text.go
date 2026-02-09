package formatter

import (
	"fmt"
	"strings"

	"github.com/mule-ai/search/internal/searxng"
)

// TextFormatter formats search results as plain text.
type TextFormatter struct {
	BaseFormatter
	NoColor bool // Disable colored output
}

// NewTextFormatter creates a new text formatter.
//
// The noColor parameter controls whether colored output is disabled.
//
// Example:
//
//	tf := formatter.NewTextFormatter(false) // enable colors
//	output, err := tf.Format(response)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func NewTextFormatter(noColor bool) *TextFormatter {
	return &TextFormatter{
		BaseFormatter: *NewBaseFormatter(),
		NoColor:       noColor,
	}
}

// Format formats the search results as plain text.
//
// Returns a formatted string with numbered results, URLs, and metadata.
// Returns an error if the response is nil.
//
// Example:
//
//	response := &searxng.SearchResponse{
//	    Query: "golang",
//	    Results: []searxng.SearchResult{...},
//	    NumberOfResults: 100,
//	}
//	tf := formatter.NewTextFormatter(false)
//	output, err := tf.Format(response)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func (f *TextFormatter) Format(result *searxng.SearchResponse) (string, error) {
	if result == nil {
		return "", fmt.Errorf("nil response provided")
	}

	var buf strings.Builder

	// Header
	buf.WriteString(fmt.Sprintf("%s\n", result.Query))
	buf.WriteString(strings.Repeat("=", len(result.Query)) + "\n\n")

	// Display results count - use NumberOfResults if available, otherwise use count of returned results
	totalResults := result.NumberOfResults
	if totalResults == 0 {
		totalResults = len(result.Results)
	}

	if len(result.Results) == 0 {
		buf.WriteString("No results found.\n\n")
	} else {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs", totalResults, result.SearchTime))

		// Add page info if not on first page
		if result.Page > 1 {
			buf.WriteString(fmt.Sprintf(" (Page %d)", result.Page))
		}

		buf.WriteString("\n\n")
	}

	// Results
	for i, res := range result.Results {
		f.formatResult(&buf, res, i)
		if i < len(result.Results)-1 {
			buf.WriteString("\n")
		}
	}

	// Answers
	if len(result.Answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range result.Answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	// Infoboxes
	if len(result.Infoboxes) > 0 {
		for _, infobox := range result.Infoboxes {
			// Infobox header
			buf.WriteString(fmt.Sprintf("\n## %s\n", infobox.Infobox))

			// Content if available
			if len(infobox.Content) > 0 {
				buf.WriteString(fmt.Sprintf("\n%s\n", infobox.Content))
			}

			// Attributes (key-value pairs)
			if len(infobox.Attributes) > 0 {
				buf.WriteString("\n")
				for _, attr := range infobox.Attributes {
					if len(attr.Value) > 0 {
						buf.WriteString(fmt.Sprintf("  • %s: %s\n", attr.Label, attr.Value))
					}
				}
			}

			// URLs (official website, wikipedia, etc)
			if len(infobox.URLs) > 0 {
				buf.WriteString("\n  Links:\n")
				for _, urlInfo := range infobox.URLs {
					prefix := "    "
					if urlInfo.Official {
						prefix = "  ★ "
					}
					buf.WriteString(fmt.Sprintf("%s%s\n", prefix, urlInfo.Title))
					if len(urlInfo.URL) > 0 {
						buf.WriteString(fmt.Sprintf("      %s\n", urlInfo.URL))
					}
				}
			}

			// Engine source
			if len(infobox.Engine) > 0 {
				buf.WriteString(fmt.Sprintf("\n  Source: %s\n", infobox.Engine))
			}
		}
	}

	// Suggestions
	if len(result.Suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range result.Suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String(), nil
}

func (f *TextFormatter) formatResult(buf *strings.Builder, result searxng.SearchResult, index int) {
	// Numbered title
	title := f.colorize(fmt.Sprintf("[%d] %s", index+1, result.Title), "bold")
	buf.WriteString(title + "\n")

	// URL
	buf.WriteString(fmt.Sprintf("    %s\n", result.URL))

	// Source and score
	var sourceInfo strings.Builder
	sourceInfo.WriteString("    Source: " + result.Engine)
	if result.Score > 0 {
		sourceInfo.WriteString(fmt.Sprintf(" | Score: %.2f", result.Score))
	}
	buf.WriteString(sourceInfo.String() + "\n")

	// Category if available
	if len(result.Category) > 0 {
		buf.WriteString(fmt.Sprintf("    Category: %s\n", result.Category))
	}

	// Content
	if len(result.Content) > 0 {
		buf.WriteString("\n")

		// Truncate content if too long
		content := result.Content
		if len(content) > f.Width-8 {
			content = content[:f.Width-11] + "..."
		}

		lines := strings.Split(content, "\n")
		for _, line := range lines {
			wrapped := f.MaxWidth(line, f.Width-8)
			for _, wline := range wrapped {
				buf.WriteString("    " + wline + "\n")
			}
		}
	}
}

// colorize adds ANSI color codes to text
func (f *TextFormatter) colorize(text string, style string) string {
	if f.NoColor {
		return text
	}

	codes := map[string]string{
		"bold":     "\033[1m",
		"red":      "\033[31m",
		"green":    "\033[32m",
		"yellow":   "\033[33m",
		"blue":     "\033[34m",
		"magenta":  "\033[35m",
		"cyan":     "\033[36m",
		"reset":    "\033[0m",
	}

	if code, ok := codes[style]; ok {
		return code + text + codes["reset"]
	}
	return text
}

// FormatResult formats a single result as plain text
func (f *TextFormatter) FormatResult(result searxng.SearchResult, index int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("[%d] %s\n", index+1, result.Title))
	buf.WriteString(fmt.Sprintf("    %s\n", result.URL))

	var sourceInfo strings.Builder
	sourceInfo.WriteString("    Source: " + result.Engine)
	if result.Score > 0 {
		sourceInfo.WriteString(fmt.Sprintf(" | Score: %.2f", result.Score))
	}
	buf.WriteString(sourceInfo.String() + "\n")

	if len(result.Content) > 0 {
		content := result.Content
		if len(content) > f.Width-8 {
			content = content[:f.Width-11] + "..."
		}
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			wrapped := f.MaxWidth(line, f.Width-8)
			for _, wline := range wrapped {
				buf.WriteString("    " + wline + "\n")
			}
		}
	}

	return buf.String()
}

// FormatResults formats multiple results as plain text
func (f *TextFormatter) FormatResults(results []searxng.SearchResult, query string, total int, searchTime float64) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs\n\n", total, searchTime))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// FormatWithQuery formats results with a custom query string
func (f *TextFormatter) FormatWithQuery(query string, results []searxng.SearchResult, total int, searchTime float64) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs\n\n", total, searchTime))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// FormatWithAnswers formats results with answers
func (f *TextFormatter) FormatWithAnswers(query string, results []searxng.SearchResult, answers []string) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	return buf.String()
}

// FormatWithInfoboxes formats results with infoboxes
func (f *TextFormatter) FormatWithInfoboxes(query string, results []searxng.SearchResult, infoboxes []searxng.Infobox) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
		}
	}

	return buf.String()
}

// FormatWithSuggestions formats results with suggestions
func (f *TextFormatter) FormatWithSuggestions(query string, results []searxng.SearchResult, suggestions []string) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String()
}

// FormatFull formats the complete search response
func (f *TextFormatter) FormatFull(query string, results []searxng.SearchResult, answers []string, infoboxes []searxng.Infobox, suggestions []string) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	if len(infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
		}
	}

	if len(suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String()
}

// FormatHeader formats the search header
func (f *TextFormatter) FormatHeader(query string, total int, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs\n\n", total, searchTime))
	}
	return buf.String()
}

// FormatFooter formats the search footer
func (f *TextFormatter) FormatFooter(query string, total int, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("\n*Search completed in %.2fs*\n", searchTime))
	return buf.String()
}

// FormatResultsTable formats results as a table
func (f *TextFormatter) FormatResultsTable(results []searxng.SearchResult, query string, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", total))
	}

	if len(results) == 0 {
		return buf.String()
	}

	buf.WriteString(fmt.Sprintf("%-3s %-50s %-20s %8s\n", "#", "Title", "Source", "Score"))
	buf.WriteString(strings.Repeat("-", 80) + "\n")

	for i, res := range results {
		title := f.TruncateWithEllipsis(res.Title, 50)
		buf.WriteString(fmt.Sprintf("%-3d %-50s %-20s %8.2f\n", i+1, title, res.Engine, res.Score))
	}

	return buf.String()
}

// FormatSimple formats results with just title and URL
func (f *TextFormatter) FormatSimple(results []searxng.SearchResult, query string, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", total))
	}

	for i, res := range results {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, res.Title))
		buf.WriteString(fmt.Sprintf("    %s\n\n", res.URL))
	}

	return buf.String()
}

// FormatWithSource formats a result with source info
func (f *TextFormatter) FormatWithSource(result searxng.SearchResult, index int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("[%d] %s\n", index+1, result.Title))
	buf.WriteString(fmt.Sprintf("    %s\n", result.URL))
	buf.WriteString(fmt.Sprintf("    Source: %s | Score: %.2f\n", result.Engine, result.Score))
	return buf.String()
}

// FormatWithCategory formats a result with category
func (f *TextFormatter) FormatWithCategory(result searxng.SearchResult, index int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("[%d] %s\n", index+1, result.Title))
	buf.WriteString(fmt.Sprintf("    %s\n", result.URL))
	if len(result.Category) > 0 {
		buf.WriteString(fmt.Sprintf("    Category: %s\n", result.Category))
	}
	return buf.String()
}

// FormatWithMetadata formats a result with metadata
func (f *TextFormatter) FormatWithMetadata(result searxng.SearchResult, index int, metadata map[string]interface{}) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("[%d] %s\n", index+1, result.Title))
	buf.WriteString(fmt.Sprintf("    %s\n", result.URL))

	var info []string
	if len(result.Engine) > 0 {
		info = append(info, fmt.Sprintf("Source: %s", result.Engine))
	}
	if result.Score > 0 {
		info = append(info, fmt.Sprintf("Score: %.2f", result.Score))
	}
	if len(metadata) > 0 {
		for k, v := range metadata {
			info = append(info, fmt.Sprintf("%s: %v", k, v))
		}
	}

	if len(info) > 0 {
		buf.WriteString("    " + strings.Join(info, " | ") + "\n")
	}

	return buf.String()
}

// FormatWithAnswersAndInfoboxes formats with answers and infoboxes
func (f *TextFormatter) FormatWithAnswersAndInfoboxes(query string, results []searxng.SearchResult, answers []string, infoboxes []searxng.Infobox) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	if len(infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
		}
	}

	return buf.String()
}

// FormatWithAnswersAndSuggestions formats with answers and suggestions
func (f *TextFormatter) FormatWithAnswersAndSuggestions(query string, results []searxng.SearchResult, answers []string, suggestions []string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	if len(suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String()
}

// FormatWithInfoboxesAndSuggestions formats with infoboxes and suggestions
func (f *TextFormatter) FormatWithInfoboxesAndSuggestions(query string, results []searxng.SearchResult, infoboxes []searxng.Infobox, suggestions []string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
		}
	}

	if len(suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String()
}

// FormatWithAllMetadata formats with all metadata fields
func (f *TextFormatter) FormatWithAllMetadata(query string, results []searxng.SearchResult, answers []string, infoboxes []searxng.Infobox, suggestions []string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")

	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	if len(infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
		}
	}

	if len(suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String()
}

// FormatResponse formats a SearchResponse struct
func (f *TextFormatter) FormatResponse(result *searxng.SearchResponse) string {
	// Extract answer strings from Answer structs
	answerStrings := make([]string, len(result.Answers))
	for i, a := range result.Answers {
		answerStrings[i] = a.Answer
	}

	return f.FormatWithAllMetadata(
		result.Query,
		result.Results,
		answerStrings,
		result.Infoboxes,
		result.Suggestions,
	)
}

// FormatCompact formats results in a compact style
func (f *TextFormatter) FormatCompact(results []searxng.SearchResult, query string, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", total))
	}

	for i, res := range results {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, res.Title))
		buf.WriteString(fmt.Sprintf("    %s\n", res.URL))
		if len(res.Content) > 0 {
			buf.WriteString(fmt.Sprintf("    %s\n", f.TruncateWithEllipsis(res.Content, 76)))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// FormatWithSearchTime formats with custom search time
func (f *TextFormatter) FormatWithSearchTime(query string, results []searxng.SearchResult, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs\n\n", len(results), searchTime))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// FormatWithInstance formats with custom instance URL
func (f *TextFormatter) FormatWithInstance(query string, results []searxng.SearchResult, instance string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", len(results)))
		buf.WriteString(fmt.Sprintf("* Instance: %s\n", instance))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// FormatWithAllOptions formats with all available options
func (f *TextFormatter) FormatWithAllOptions(query string, results []searxng.SearchResult, answers []string, infoboxes []searxng.Infobox, suggestions []string, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs\n\n", len(results), searchTime))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	if len(infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
		}
	}

	if len(suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String()
}

// FormatWithResultCount formats with custom result count
func (f *TextFormatter) FormatWithResultCount(query string, results []searxng.SearchResult, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results\n\n", total))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// FormatWithSearchTimeAndInstance formats with search time and instance
func (f *TextFormatter) FormatWithSearchTimeAndInstance(query string, results []searxng.SearchResult, searchTime float64, instance string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs\n\n", len(results), searchTime))
		buf.WriteString(fmt.Sprintf("* Instance: %s\n", instance))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// FormatWithAll formats with all options
func (f *TextFormatter) FormatWithAll(query string, results []searxng.SearchResult, answers []string, infoboxes []searxng.Infobox, suggestions []string, total int, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s\n", query))
	buf.WriteString(strings.Repeat("=", len(query)) + "\n\n")
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found %d results in %.2fs\n\n", total, searchTime))
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer))
		}
	}

	if len(infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
		}
	}

	if len(suggestions) > 0 {
		buf.WriteString("\n## Suggestions\n\n")
		for _, suggestion := range suggestions {
			buf.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
	}

	return buf.String()
}

// FormatResponseWithAll formats a SearchResponse with all metadata
func (f *TextFormatter) FormatResponseWithAll(result *searxng.SearchResponse, instance string) string {
	// Extract answer strings from Answer structs
	answerStrings := make([]string, len(result.Answers))
	for i, a := range result.Answers {
		answerStrings[i] = a.Answer
	}

	return f.FormatWithAll(
		result.Query,
		result.Results,
		answerStrings,
		result.Infoboxes,
		result.Suggestions,
		result.NumberOfResults,
		result.SearchTime,
	)
}
