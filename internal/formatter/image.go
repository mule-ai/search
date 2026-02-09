package formatter

import (
	"fmt"
	"strings"

	"github.com/mule-ai/search/internal/searxng"
)

// ImageFormatter handles formatting for image search results.
//
// It formats image results with the image URL prominently displayed,
// making it easy to view and access image search results.
type ImageFormatter struct {
	plainFormatter *TextFormatter
}

// NewImageFormatter creates a new image-specific formatter.
//
// The returned formatter is optimized for displaying image search results
// with the image URL, source page, and metadata clearly presented.
//
// Example:
//
//	// Search for images
//	client := searxng.NewClient(cfg)
//	results, _ := client.SearchWithConfig("cats", 10, "json", "images", 30, "en", 1, 1, "")
//
//	// Format with image formatter
//	imgFormatter := formatter.NewImageFormatter()
//	output, _ := imgFormatter.Format(results)
//	fmt.Println(output)
func NewImageFormatter() *ImageFormatter {
	return &ImageFormatter{
		plainFormatter: NewTextFormatter(false),
	}
}

// Format formats search results with special handling for images.
//
// For image results, the image URL (ImgSrc) is displayed prominently.
func (f *ImageFormatter) Format(response *searxng.SearchResponse) (string, error) {
	if len(response.Results) == 0 {
		return fmt.Sprintf("%s\n\nNo results found.", f.formatHeader(response)), nil
	}

	var sb strings.Builder
	
	// Header
	sb.WriteString(f.formatHeader(response))
	sb.WriteString("\n\n")
	
	// Results
	for i, result := range response.Results {
		sb.WriteString(f.formatImageResult(i+1, result))
		if i < len(response.Results)-1 {
			sb.WriteString("\n")
		}
	}
	
	// Footer
	sb.WriteString("\n")
	sb.WriteString(f.formatFooter(response))
	
	return sb.String(), nil
}

func (f *ImageFormatter) formatHeader(response *searxng.SearchResponse) string {
	var sb strings.Builder
	
	query := response.Query
	if query == "" {
		query = "Search"
	}
	
	sb.WriteString(query)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", len(query)))
	
	if response.NumberOfResults > 0 {
		sb.WriteString(fmt.Sprintf("\nFound %d results", response.NumberOfResults))
	}
	
	return sb.String()
}

func (f *ImageFormatter) formatImageResult(index int, result searxng.SearchResult) string {
	var sb strings.Builder
	
	// Number and title
	imageURL := searxng.GetImageURL(result)
	if imageURL != "" {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", index, result.Title))
		sb.WriteString(fmt.Sprintf("    Image: %s\n", imageURL))
	} else {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", index, result.Title))
		sb.WriteString(fmt.Sprintf("    URL: %s\n", result.URL))
	}
	
	// Source page (if different from image URL)
	if imageURL != "" && result.URL != "" && result.URL != imageURL {
		sb.WriteString(fmt.Sprintf("    Source: %s\n", result.URL))
	}
	
	// Engine and metadata
	if result.Engine != "" {
		metadata := []string{result.Engine}
		if result.Score > 0 {
			metadata = append(metadata, fmt.Sprintf("Score: %.2f", result.Score))
		}
		sb.WriteString(fmt.Sprintf("    Source: %s\n", strings.Join(metadata, " | ")))
	}
	
	// Content (truncated description)
	if result.Content != "" {
		content := strings.TrimSpace(result.Content)
		if len(content) > 200 {
			content = content[:197] + "..."
		}
		sb.WriteString(fmt.Sprintf("\n    %s\n", content))
	}
	
	return sb.String()
}

func (f *ImageFormatter) formatFooter(response *searxng.SearchResponse) string {
	var sb strings.Builder
	
	if response.SearchTime > 0 {
		sb.WriteString(fmt.Sprintf("\nSearch time: %.2fs", response.SearchTime))
	}
	
	if response.Query != "" {
		sb.WriteString(fmt.Sprintf(" | Instance: (configured instance)"))
	}
	
	return sb.String()
}

// FormatAsMarkdown formats image results as markdown with embedded images.
//
// This creates markdown output with embedded image references, suitable for
// rendering in markdown viewers or HTML. Each result includes the image
// reference and a link to the source page.
//
// Example:
//
//	// Search for images
//	client := searxng.NewClient(cfg)
//	results, _ := client.SearchWithConfig("cats", 10, "json", "images", 30, "en", 1, 1, "")
//
//	// Format as markdown
//	md, _ := formatter.FormatAsMarkdown(results)
//	fmt.Println(md)
//	// Output includes embedded images like:
//	// ## [1. Cute Cat](https://source.page)
//	//
//	// ![Cute Cat](https://image.url/cat.jpg)
func FormatAsMarkdown(response *searxng.SearchResponse) (string, error) {
	if len(response.Results) == 0 {
		return fmt.Sprintf("# Search Results: %s\n\nNo results found.", response.Query), nil
	}

	var sb strings.Builder
	
	// Header
	sb.WriteString(fmt.Sprintf("# Search Results: %s\n\n", response.Query))
	
	if response.NumberOfResults > 0 {
		sb.WriteString(fmt.Sprintf("Found **%d** results", response.NumberOfResults))
		if response.SearchTime > 0 {
			sb.WriteString(fmt.Sprintf(" in %.2fs\n\n", response.SearchTime))
		} else {
			sb.WriteString("\n\n")
		}
	}
	
	// Results
	for i, result := range response.Results {
		sb.WriteString(formatImageMarkdown(i+1, result))
		if i < len(response.Results)-1 {
			sb.WriteString("\n")
		}
	}
	
	return sb.String(), nil
}

func formatImageMarkdown(index int, result searxng.SearchResult) string {
	var sb strings.Builder
	
	imageURL := searxng.GetImageURL(result)
	sourceURL := result.URL
	if sourceURL == "" || sourceURL == imageURL {
		sourceURL = imageURL
	}
	
	// Title with link to source page
	sb.WriteString(fmt.Sprintf("## [%d. %s](%s)\n\n", index, result.Title, sourceURL))
	
	// Embedded image markdown
	if imageURL != "" {
		sb.WriteString(fmt.Sprintf(`
![%s](%s)

`, result.Title, imageURL))
	}
	
	// Metadata
	metadata := []string{}
	if result.Engine != "" {
		metadata = append(metadata, fmt.Sprintf("**Source:** %s", result.Engine))
	}
	if result.Score > 0 {
		metadata = append(metadata, fmt.Sprintf("**Score:** %.2f", result.Score))
	}
	if len(metadata) > 0 {
		sb.WriteString(strings.Join(metadata, " | "))
		sb.WriteString("\n\n")
	}
	
	// Content/description
	if result.Content != "" {
		sb.WriteString(result.Content)
		sb.WriteString("\n\n")
	}
	
	sb.WriteString("---\n\n")
	
	return sb.String()
}

// FormatAsJSONWithImages formats image results as JSON with explicit image field.
//
// This uses the standard JSON formatter, which includes the img_src field
// for image results, making it easy to extract image URLs programmatically.
//
// Example:
//
//	// Search for images
//	client := searxng.NewClient(cfg)
//	results, _ := client.SearchWithConfig("cats", 10, "json", "images", 30, "en", 1, 1, "")
//
//	// Format as JSON
//	json, _ := formatter.FormatAsJSONWithImages(results)
//	fmt.Println(json)
//	// Output includes img_src field:
//	// {"results": [{"title": "Cute Cat", "img_src": "https://...", ...}]}
func FormatAsJSONWithImages(response *searxng.SearchResponse) (string, error) {
	// Use the standard JSON formatter but ensure image URLs are prominent
	jsonFormatter := NewJSONFormatter()
	return jsonFormatter.Format(response)
}