package formatter

import (
	"fmt"
	"strings"

	"github.com/mule-ai/search/internal/searxng"
)

// MarkdownFormatter formats search results as Markdown.
type MarkdownFormatter struct {
	BaseFormatter
}

// NewMarkdownFormatter creates a new Markdown formatter.
//
// Example:
//
//	mf := formatter.NewMarkdownFormatter()
//	output, err := mf.Format(response)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{
		BaseFormatter: *NewBaseFormatter(),
	}
}

// Format formats the search results as Markdown.
//
// Returns a formatted string with proper Markdown headings, links, and formatting.
// Returns an error if the response is nil.
//
// Example:
//
//	response := &searxng.SearchResponse{
//	    Query: "golang",
//	    Results: []searxng.SearchResult{...},
//	    NumberOfResults: 100,
//	}
//	mf := formatter.NewMarkdownFormatter()
//	output, err := mf.Format(response)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
func (f *MarkdownFormatter) Format(result *searxng.SearchResponse) (string, error) {
	if result == nil {
		return "", fmt.Errorf("nil response provided")
	}

	var buf strings.Builder

	// Header
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", result.Query))
	
	// Display results count - use NumberOfResults if available, otherwise use count of returned results
	totalResults := result.NumberOfResults
	if totalResults == 0 {
		totalResults = len(result.Results)
	}
	
	if totalResults > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs", totalResults, result.SearchTime))
		
		// Add page info if not on first page
		if result.Page > 1 {
			buf.WriteString(fmt.Sprintf(" (Page %d)", result.Page))
		}
		
		buf.WriteString("\n\n")
	} else {
		buf.WriteString("No results found\n\n")
	}

	// Results
	for i, res := range result.Results {
		f.formatResult(&buf, res, i)
		if i < len(result.Results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	// Answers
	if len(result.Answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range result.Answers {
			buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
		}
	}

	// Infoboxes
	if len(result.Infoboxes) > 0 {
		buf.WriteString("\n## Infoboxes\n\n")
		for _, infobox := range result.Infoboxes {
			buf.WriteString(fmt.Sprintf("- %s\n", infobox.Infobox))
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

func (f *MarkdownFormatter) formatResult(buf *strings.Builder, result searxng.SearchResult, index int) {
	// Title as heading
	buf.WriteString(fmt.Sprintf("## [%s](%s)\n", f.escapeMarkdown(result.Title), result.URL))

	// Source and score
	var sourceInfo strings.Builder
	sourceInfo.WriteString(fmt.Sprintf("**Source:** %s", result.Engine))
	if result.Score > 0 {
		sourceInfo.WriteString(fmt.Sprintf(" | **Score:** %.2f", result.Score))
	}
	buf.WriteString(sourceInfo.String() + "\n\n")

	// Content
	if len(result.Content) > 0 {
		buf.WriteString(f.TruncateWithEllipsis(result.Content, f.Width-4) + "\n\n")
	}

	// Category if available
	if len(result.Category) > 0 {
		buf.WriteString(fmt.Sprintf("*Category: %s*\n", result.Category))
	}
}

// escapeMarkdown escapes markdown special characters
func (f *MarkdownFormatter) escapeMarkdown(s string) string {
	chars := []string{"\\", "`", "*", "_", "{", "}", "[", "]", "(", ")", "#", "+", "-", ".", "!", "|"}
	for _, c := range chars {
		s = strings.ReplaceAll(s, c, "\\"+c)
	}
	return s
}

// FormatResult formats a single result as Markdown
func (f *MarkdownFormatter) FormatResult(result searxng.SearchResult, index int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("## [%s](%s)\n", f.escapeMarkdown(result.Title), result.URL))

	var sourceInfo strings.Builder
	sourceInfo.WriteString(fmt.Sprintf("**Source:** %s", result.Engine))
	if result.Score > 0 {
		sourceInfo.WriteString(fmt.Sprintf(" | **Score:** %.2f", result.Score))
	}
	buf.WriteString(sourceInfo.String() + "\n\n")

	if len(result.Content) > 0 {
		buf.WriteString(f.TruncateWithEllipsis(result.Content, f.Width-4) + "\n\n")
	}

	if len(result.Category) > 0 {
		buf.WriteString(fmt.Sprintf("*Category: %s*\n", result.Category))
	}

	return buf.String()
}

// FormatResults formats multiple results as Markdown
func (f *MarkdownFormatter) FormatResults(results []searxng.SearchResult, query string, total int, searchTime float64) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs\n\n", total, searchTime))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	return buf.String()
}

// FormatWithQuery formats results with a custom query string
func (f *MarkdownFormatter) FormatWithQuery(query string, results []searxng.SearchResult, total int, searchTime float64) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs\n\n", total, searchTime))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	return buf.String()
}

// FormatWithAnswers formats results with answers
func (f *MarkdownFormatter) FormatWithAnswers(query string, results []searxng.SearchResult, answers []searxng.Answer) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
		}
	}

	return buf.String()
}

// FormatWithInfoboxes formats results with infoboxes
func (f *MarkdownFormatter) FormatWithInfoboxes(query string, results []searxng.SearchResult, infoboxes []searxng.Infobox) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
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
func (f *MarkdownFormatter) FormatWithSuggestions(query string, results []searxng.SearchResult, suggestions []string) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
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
func (f *MarkdownFormatter) FormatFull(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
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
func (f *MarkdownFormatter) FormatHeader(query string, total int, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs\n\n", total, searchTime))
	} else {
		buf.WriteString("No results found\n\n")
	}
	return buf.String()
}

// FormatFooter formats the search footer
func (f *MarkdownFormatter) FormatFooter(query string, total int, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("\n*Search completed in %.2fs*\n", searchTime))
	return buf.String()
}

// FormatResultsTable formats results as a markdown table
func (f *MarkdownFormatter) FormatResultsTable(results []searxng.SearchResult, query string, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", total))
	} else {
		buf.WriteString("No results found\n\n")
	}

	if len(results) == 0 {
		return buf.String()
	}

	buf.WriteString("| # | Title | Source | Score |\n")
	buf.WriteString("|---|-------|--------|-------|\n")

	for i, res := range results {
		title := f.TruncateWithEllipsis(res.Title, 50)
		buf.WriteString(fmt.Sprintf("| %d | [%s](%s) | %s | %.2f |\n", i+1, f.escapeMarkdown(title), res.URL, res.Engine, res.Score))
	}

	return buf.String()
}

// FormatSimple formats results with just title and URL
func (f *MarkdownFormatter) FormatSimple(results []searxng.SearchResult, query string, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", total))
	}

	for i, res := range results {
		buf.WriteString(fmt.Sprintf("%d. [%s](%s)\n", i+1, f.escapeMarkdown(res.Title), res.URL))
	}

	return buf.String()
}

// FormatWithSource formats a result with source info
func (f *MarkdownFormatter) FormatWithSource(result searxng.SearchResult, index int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("## [%s](%s)\n", f.escapeMarkdown(result.Title), result.URL))
	buf.WriteString(fmt.Sprintf("**Source:** %s | **Score:** %.2f\n\n", result.Engine, result.Score))
	if len(result.Content) > 0 {
		buf.WriteString(f.TruncateWithEllipsis(result.Content, f.Width-4) + "\n")
	}
	return buf.String()
}

// FormatWithCategory formats a result with category
func (f *MarkdownFormatter) FormatWithCategory(result searxng.SearchResult, index int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("## [%s](%s)\n", f.escapeMarkdown(result.Title), result.URL))
	if len(result.Category) > 0 {
		buf.WriteString(fmt.Sprintf("*Category: %s*\n\n", result.Category))
	}
	if len(result.Content) > 0 {
		buf.WriteString(f.TruncateWithEllipsis(result.Content, f.Width-4) + "\n")
	}
	return buf.String()
}

// FormatWithMetadata formats a result with metadata
func (f *MarkdownFormatter) FormatWithMetadata(result searxng.SearchResult, index int, metadata map[string]interface{}) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("## [%s](%s)\n", f.escapeMarkdown(result.Title), result.URL))

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
		buf.WriteString("**" + strings.Join(info, " | ") + "**\n\n")
	} else {
		buf.WriteString("\n")
	}

	if len(result.Content) > 0 {
		buf.WriteString(f.TruncateWithEllipsis(result.Content, f.Width-4) + "\n")
	}

	return buf.String()
}

// FormatWithAnswersAndInfoboxes formats with answers and infoboxes
func (f *MarkdownFormatter) FormatWithAnswersAndInfoboxes(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
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
func (f *MarkdownFormatter) FormatWithAnswersAndSuggestions(query string, results []searxng.SearchResult, answers []searxng.Answer, suggestions []string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
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
func (f *MarkdownFormatter) FormatWithInfoboxesAndSuggestions(query string, results []searxng.SearchResult, infoboxes []searxng.Infobox, suggestions []string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
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
func (f *MarkdownFormatter) FormatWithAllMetadata(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
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
func (f *MarkdownFormatter) FormatResponse(result *searxng.SearchResponse) string {
	return f.FormatWithAllMetadata(
		result.Query,
		result.Results,
		result.Answers,
		result.Infoboxes,
		result.Suggestions,
	)
}

// FormatCompact formats results in a compact style
func (f *MarkdownFormatter) FormatCompact(results []searxng.SearchResult, query string, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", total))
	}

	for i, res := range results {
		buf.WriteString(fmt.Sprintf("### %d. %s\n", i+1, f.escapeMarkdown(res.Title)))
		buf.WriteString(fmt.Sprintf("*%s*\n\n", res.URL))
	}

	return buf.String()
}

// FormatWithSearchTime formats with custom search time
func (f *MarkdownFormatter) FormatWithSearchTime(query string, results []searxng.SearchResult, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs\n\n", len(results), searchTime))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	return buf.String()
}

// FormatWithInstance formats with custom instance URL
func (f *MarkdownFormatter) FormatWithInstance(query string, results []searxng.SearchResult, instance string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", len(results)))
		buf.WriteString(fmt.Sprintf("*Instance: %s*\n\n", instance))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	return buf.String()
}

// FormatWithAllOptions formats with all available options
func (f *MarkdownFormatter) FormatWithAllOptions(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs\n\n", len(results), searchTime))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
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
func (f *MarkdownFormatter) FormatWithResultCount(query string, results []searxng.SearchResult, total int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results\n\n", total))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	return buf.String()
}

// FormatWithSearchTimeAndInstance formats with search time and instance
func (f *MarkdownFormatter) FormatWithSearchTimeAndInstance(query string, results []searxng.SearchResult, searchTime float64, instance string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if len(results) > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs\n\n", len(results), searchTime))
		buf.WriteString(fmt.Sprintf("*Instance: %s*\n\n", instance))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	return buf.String()
}

// FormatWithAll formats with all options
func (f *MarkdownFormatter) FormatWithAll(query string, results []searxng.SearchResult, answers []searxng.Answer, infoboxes []searxng.Infobox, suggestions []string, total int, searchTime float64) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Search Results: %s\n\n", query))
	if total > 0 {
		buf.WriteString(fmt.Sprintf("Found **%d** results in %.2fs\n\n", total, searchTime))
	} else {
		buf.WriteString("No results found\n\n")
	}

	for i, res := range results {
		f.formatResult(&buf, res, i)
		if i < len(results)-1 {
			buf.WriteString("\n---\n\n")
		}
	}

	if len(answers) > 0 {
		buf.WriteString("\n## Answers\n\n")
		for _, answer := range answers {
buf.WriteString(fmt.Sprintf("- %s\n", answer.Answer))
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
func (f *MarkdownFormatter) FormatResponseWithAll(result *searxng.SearchResponse, instance string) string {
	return f.FormatWithAll(
		result.Query,
		result.Results,
		result.Answers,
		result.Infoboxes,
		result.Suggestions,
		result.NumberOfResults,
		result.SearchTime,
	)
}
