// Package validation provides input validation for the search CLI.
//
// It validates queries, configuration values, and search parameters,
// returning detailed errors with suggestions when validation fails.
package validation

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mule-ai/search/internal/errors"
	"github.com/mule-ai/search/internal/searxng"
)

// ValidFormats is the list of supported output formats.
var ValidFormats = []string{"text", "json", "markdown"}

// ValidSafeSearchLevels is the list of valid safe search levels.
var ValidSafeSearchLevels = []int{0, 1, 2}

// ValidationError represents a validation error with context.
//
// It includes the field name, the invalid value, a descriptive message,
// and an optional suggestion for fixing the error.
type ValidationError struct {
	Field      string
	Value      interface{}
	Message    string
	Suggestion string
}

func (e ValidationError) Error() string {
	msg := fmt.Sprintf("validation error for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
	if e.Suggestion != "" {
		msg += fmt.Sprintf("\nSuggestion: %s", e.Suggestion)
	}
	return msg
}

// ValidateQuery checks if the search query is valid.
//
// It ensures the query is not empty and does not exceed 1000 characters.
// Returns an error with a suggestion if validation fails.
//
// Example:
//
//	err := validation.ValidateQuery("golang tutorials")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateQuery(query string) error {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return errors.EmptyQuery()
	}
	if len(trimmed) > 1000 {
		return errors.New(errors.ErrCodeInvalidRange, "query is too long (max 1000 characters)").
			WithSuggestion("Shorten your search query to 1000 characters or less").
			WithVerbose(fmt.Sprintf("Current query length: %d characters", len(query)))
	}
	return nil
}

// ValidateResultCount checks if the result count is within valid range.
//
// The valid range is 1-100 results per page.
//
// Example:
//
//	err := validation.ValidateResultCount(20)
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateResultCount(count int) error {
	if count < 1 {
		return errors.InvalidRange("results", 1, 100, count)
	}
	if count > 100 {
		return errors.InvalidRange("results", 1, 100, count)
	}
	return nil
}

// ValidateTimeout checks if the timeout is within valid range.
//
// The valid range is 1-300 seconds (5 minutes).
//
// Example:
//
//	err := validation.ValidateTimeout(30)
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateTimeout(timeout int) error {
	if timeout < 1 {
		return ValidationError{
			Field:   "timeout",
			Value:   timeout,
			Message: "timeout must be at least 1 second",
		}
	}
	if timeout > 300 {
		return ValidationError{
			Field:   "timeout",
			Value:   timeout,
			Message: "timeout cannot exceed 300 seconds (5 minutes)",
		}
	}
	return nil
}

// ValidateInstanceURL checks if the instance URL is valid.
//
// It ensures the URL is properly formatted, uses http or https scheme,
// and includes a valid host.
//
// Example:
//
//	err := validation.ValidateInstanceURL("https://search.butler.ooo")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateInstanceURL(instanceURL string) error {
	if instanceURL == "" {
		return ValidationError{
			Field:   "instance",
			Value:   instanceURL,
			Message: "instance URL cannot be empty",
		}
	}

	// Parse URL to validate format
	parsedURL, err := url.Parse(instanceURL)
	if err != nil {
		return ValidationError{
			Field:   "instance",
			Value:   instanceURL,
			Message: fmt.Sprintf("invalid URL format: %v", err),
		}
	}

	// Check scheme - only https and http are allowed
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ValidationError{
			Field:   "instance",
			Value:   instanceURL,
			Message: "instance URL must use http or https scheme",
		}
	}

	// Enforce HTTPS for non-localhost instances (security requirement)
	if parsedURL.Scheme == "http" {
		host := parsedURL.Hostname()
		if host != "localhost" && host != "127.0.0.1" && host != "::1" {
			return ValidationError{
				Field:      "instance",
				Value:      instanceURL,
				Message:    "HTTPS is required for all non-localhost instances",
				Suggestion: "Use https:// instead of http://",
			}
		}
	}

	// Check host
	if parsedURL.Host == "" {
		return ValidationError{
			Field:   "instance",
			Value:   instanceURL,
			Message: "instance URL must include a host",
		}
	}

	return nil
}

// ValidateFormat checks if the format is supported.
//
// Valid formats are: text, json, markdown.
//
// Example:
//
//	err := validation.ValidateFormat("json")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateFormat(format string) error {
	format = strings.ToLower(strings.TrimSpace(format))
	for _, validFormat := range ValidFormats {
		if format == validFormat {
			return nil
		}
	}
	return ValidationError{
		Field:   "format",
		Value:   format,
		Message: fmt.Sprintf("format must be one of: %s", strings.Join(ValidFormats, ", ")),
	}
}

// ValidateSafeSearch checks if the safe search level is valid.
//
// Valid levels are: 0 (off), 1 (moderate), 2 (strict).
//
// Example:
//
//	err := validation.ValidateSafeSearch(1)
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateSafeSearch(safeSearch int) error {
	for _, validLevel := range ValidSafeSearchLevels {
		if safeSearch == validLevel {
			return nil
		}
	}
	return ValidationError{
		Field:   "safeSearch",
		Value:   safeSearch,
		Message: "safe search level must be 0 (off), 1 (moderate), or 2 (strict)",
	}
}

// ValidateLanguage checks if the language code is valid.
//
// It performs a basic ISO 639-1 check (2-letter codes) or ISO 639-1 with
// region code (5-letter codes like en_US).
//
// Example:
//
//	err := validation.ValidateLanguage("en")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateLanguage(language string) error {
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		return ValidationError{
			Field:   "language",
			Value:   language,
			Message: "language code cannot be empty",
		}
	}
	if len(language) != 2 && len(language) != 5 {
		return ValidationError{
			Field:   "language",
			Value:   language,
			Message: "language code must be 2-letter (ISO 639-1) or 5-letter (with region, e.g., en_US)",
		}
	}
	return nil
}

// ValidatePageNumber checks if the page number is valid.
//
// Valid page numbers are 1-50.
//
// Example:
//
//	err := validation.ValidatePageNumber(2)
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidatePageNumber(page int) error {
	if page < 1 {
		return ValidationError{
			Field:   "page",
			Value:   page,
			Message: "page number must be at least 1",
		}
	}
	if page > 50 {
		return ValidationError{
			Field:   "page",
			Value:   page,
			Message: "page number cannot exceed 50",
		}
	}
	return nil
}

// ValidateTimeRange checks if the time range is valid.
//
// Valid values are: day, week, month, year.
// Empty string is allowed (no time filter).
//
// Example:
//
//	err := validation.ValidateTimeRange("week")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateTimeRange(timeRange string) error {
	if timeRange == "" {
		return nil // Optional field
	}
	
	timeRange = strings.ToLower(strings.TrimSpace(timeRange))
	validRanges := []string{"day", "week", "month", "year"}
	for _, validRange := range validRanges {
		if timeRange == validRange {
			return nil
		}
	}
	return ValidationError{
		Field:   "timeRange",
		Value:   timeRange,
		Message: fmt.Sprintf("time range must be one of: %s", strings.Join(validRanges, ", ")),
	}
}

// ValidateCategory checks if the category is valid.
//
// It normalizes category aliases and checks against known SearXNG categories.
// Unknown categories are allowed (not an error) as SearXNG may support custom categories.
func ValidateCategory(category string) error {
	if category == "" {
		return ValidationError{
			Field:   "category",
			Value:   category,
			Message: "category cannot be empty",
		}
	}
	
	// Normalize the category (resolve aliases)
	normalized := searxng.NormalizeCategory(category)
	
	// Check if it's a known category
	if !searxng.IsValidCategory(normalized) {
		// Don't error for unknown categories as SearXNG may support more
		// Just warn the user via verbose output if they wanted to know
		return nil
	}
	
	return nil
}