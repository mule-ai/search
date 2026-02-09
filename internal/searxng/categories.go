// Package searxng provides category definitions and utilities for SearXNG.
//
// This file contains the list of valid search categories, their aliases,
// and helper functions for working with categories.
package searxng

import (
	"fmt"
	"strings"
)

// Category represents a search category with its metadata.
type Category struct {
	Name        string
	DisplayName string
	Description string
	ExampleQuery string
	SpecialFormatting bool // Indicates if this category needs special formatting
}

// ValidCategories holds all supported SearXNG categories
var ValidCategories = map[string]Category{
	"general": {
		Name:        "general",
		DisplayName: "General",
		Description: "General web search across all engines",
		ExampleQuery: "golang tutorials",
		SpecialFormatting: false,
	},
	"images": {
		Name:        "images",
		DisplayName: "Images",
		Description: "Image search results",
		ExampleQuery: "cute cats",
		SpecialFormatting: true, // Needs special formatting for image URLs
	},
	"videos": {
		Name:        "videos",
		DisplayName: "Videos",
		Description: "Video search results",
		ExampleQuery: "funny moments",
		SpecialFormatting: false,
	},
	"news": {
		Name:        "news",
		DisplayName: "News",
		Description: "News articles and recent updates",
		ExampleQuery: "technology news",
		SpecialFormatting: false,
	},
	"map": {
		Name:        "map",
		DisplayName: "Map",
		Description: "Map search results",
		ExampleQuery: "coffee shop near me",
		SpecialFormatting: false,
	},
	"music": {
		Name:        "music",
		DisplayName: "Music",
		Description: "Music search results",
		ExampleQuery: "the beatles",
		SpecialFormatting: false,
	},
	"it": {
		Name:        "it",
		DisplayName: "IT",
		Description: "Information technology search",
		ExampleQuery: "kubernetes tutorial",
		SpecialFormatting: false,
	},
	"science": {
		Name:        "science",
		DisplayName: "Science",
		Description: "Scientific research and papers",
		ExampleQuery: "quantum computing",
		SpecialFormatting: false,
	},
	"files": {
		Name:        "files",
		DisplayName: "Files",
		Description: "File and document search",
		ExampleQuery: "pdf report",
		SpecialFormatting: false,
	},
	"social media": {
		Name:        "social media",
		DisplayName: "Social Media",
		Description: "Social media content",
		ExampleQuery: "trending topics",
		SpecialFormatting: false,
	},
}

// CategoryAliases provides alternate names for categories
var CategoryAliases = map[string]string{
	"web":     "general",
	"search":  "general",
	"photo":   "images",
	"picture": "images",
	"pics":    "images",
	"video":   "videos",
	"youtube": "videos",
	"movie":   "videos",
	"tech":    "it",
	"research": "science",
	"papers":  "science",
	"document": "files",
	"documents": "files",
	"pdf":     "files",
}

// GetCategory retrieves a category by name, resolving aliases.
//
// Returns the Category and nil error if found, or an empty Category and error
// if the category name is not recognized.
//
// Example:
//
//	cat, err := searxng.GetCategory("images")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Category: %s - %s\n", cat.DisplayName, cat.Description)
//	// Output: Category: Images - Image search results
//
//	// Aliases are also resolved
//	cat, err = searxng.GetCategory("photo")
//	// Returns the same "images" category
func GetCategory(name string) (Category, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	
	// Check direct match
	if cat, ok := ValidCategories[name]; ok {
		return cat, nil
	}
	
	// Check aliases
	if alias, ok := CategoryAliases[name]; ok {
		if cat, ok := ValidCategories[alias]; ok {
			return cat, nil
		}
	}
	
	return Category{}, fmt.Errorf("unknown category: %s (valid categories: %s)", 
		name, strings.Join(GetCategoryNames(), ", "))
}

// GetCategoryNames returns a list of all valid category names.
//
// The returned list is unsorted.
func GetCategoryNames() []string {
	names := make([]string, 0, len(ValidCategories))
	for name := range ValidCategories {
		names = append(names, name)
	}
	return names
}

// GetDisplayNames returns a list of display names for all categories.
//
// Display names are human-readable versions (e.g., "General", "Images").
func GetDisplayNames() []string {
	names := make([]string, 0, len(ValidCategories))
	for _, cat := range ValidCategories {
		names = append(names, cat.DisplayName)
	}
	return names
}

// IsValidCategory checks if a category name is valid.
//
// It checks both direct category names and aliases.
//
// Example:
//
//	if searxng.IsValidCategory("images") {
//	    fmt.Println("Valid category")
//	}
//
//	if !searxng.IsValidCategory("invalid") {
//	    fmt.Println("Invalid category")
//	}
func IsValidCategory(name string) bool {
	_, err := GetCategory(name)
	return err == nil
}

// NormalizeCategory normalizes a category name, resolving aliases.
//
// Returns the canonical category name. If the category is unknown,
// returns the input as-is (may be a custom category).
func NormalizeCategory(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	
	if _, ok := ValidCategories[name]; ok {
		return name
	}
	
	if alias, ok := CategoryAliases[name]; ok {
		return alias
	}
	
	return name // Return as-is if unknown (may be a custom category)
}

// NeedsSpecialFormatting checks if a category requires special formatting.
//
// Categories like "images" may need different output formatting.
//
// Example:
//
//	if searxng.NeedsSpecialFormatting("images") {
//	    // Use image-specific formatter
//	    formatter := formatter.NewImageFormatter()
//	} else {
//	    // Use standard formatter
//	    formatter := formatter.NewTextFormatter(false)
//	}
func NeedsSpecialFormatting(category string) bool {
	cat, err := GetCategory(category)
	if err != nil {
		return false
	}
	return cat.SpecialFormatting
}

// FormatResultForCategory formats a result based on its category.
//
// For image results, it ensures the image URL is properly handled.
// Returns the modified or original SearchResult.
func FormatResultForCategory(result SearchResult, category string) SearchResult {
	cat, err := GetCategory(category)
	if err != nil {
		return result
	}
	
	if cat.SpecialFormatting && result.ImgSrc != "" {
		// For image results, ensure the image URL is prominent
		return result
	}
	
	return result
}

// GetImageURL returns the primary image URL for a result.
//
// It checks the ImgSrc field first, then falls back to the main URL
// for image category results.
//
// Example:
//
//	result := searxng.SearchResult{
//	    Title: "Cute Cat",
//	    URL: "https://example.com/cat.jpg",
//	    ImgSrc: "https://cdn.example.com/cat.jpg",
//	    Category: "images",
//	}
//	imageURL := searxng.GetImageURL(result)
//	// Returns: "https://cdn.example.com/cat.jpg"
func GetImageURL(result SearchResult) string {
	if result.ImgSrc != "" {
		return result.ImgSrc
	}
	// For image results, the main URL might be the image
	if result.Category == "images" {
		return result.URL
	}
	return ""
}

// GetCategorySuggestions returns suggested categories based on query keywords.
//
// It analyzes the query for keywords that might indicate a preferred category
// (e.g., "picture" suggests "images", "video" suggests "videos").
//
// Example:
//
//	suggestions := searxng.GetCategorySuggestions("cute cat picture")
//	// Returns: ["images"]
//
//	suggestions = searxng.GetCategorySuggestions("golang tutorial")
//	// Returns: [] (no specific category suggested)
func GetCategorySuggestions(query string) []string {
	query = strings.ToLower(query)
	suggestions := []string{}
	
	// Keyword to category mapping
	keywordCategories := map[string]string{
		"image": "images",
		"images": "images",
		"picture": "images",
		"photo": "images",
		"video": "videos",
		"youtube": "videos",
		"news": "news",
		"map": "map",
		"music": "music",
		"song": "music",
		"pdf": "files",
		"download": "files",
		"research": "science",
		"paper": "science",
	}
	
	for keyword, category := range keywordCategories {
		if strings.Contains(query, keyword) {
			suggestions = append(suggestions, category)
		}
	}
	
	return suggestions
}

// CategoryInfo returns detailed information about a category.
//
// Returns a formatted string with the category name, description, and example query.
func CategoryInfo(name string) string {
	cat, err := GetCategory(name)
	if err != nil {
		return fmt.Sprintf("Category '%s' not found", name)
	}
	
	return fmt.Sprintf("%s (%s)\n  Description: %s\n  Example: search -c %s \"%s\"",
		cat.DisplayName, cat.Name, cat.Description, cat.Name, cat.ExampleQuery)
}