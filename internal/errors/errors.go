// Package errors provides structured error handling for the search CLI.
//
// It defines error codes for scripting, user-friendly error messages,
// and suggestions for common error scenarios.
package errors

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ErrorCode represents a unique error code for scripting.
//
// Error codes are used to identify specific error types in scripts and automation.
type ErrorCode string

const (
	// Config errors
	ErrCodeConfigNotFound    ErrorCode = "CONFIG_NOT_FOUND"
	ErrCodeConfigInvalid     ErrorCode = "CONFIG_INVALID"
	ErrCodeConfigParseError  ErrorCode = "CONFIG_PARSE_ERROR"

	// Network errors
	ErrCodeNetworkTimeout    ErrorCode = "NETWORK_TIMEOUT"
	ErrCodeNetworkUnreachable ErrorCode = "NETWORK_UNREACHABLE"
	ErrCodeConnectionRefused ErrorCode = "CONNECTION_REFUSED"
	ErrCodeDNSFailed         ErrorCode = "DNS_FAILED"

	// API errors
	ErrCodeAPIError          ErrorCode = "API_ERROR"
	ErrCodeAPIUnavailable    ErrorCode = "API_UNAVAILABLE"
	ErrCodeInvalidResponse   ErrorCode = "INVALID_RESPONSE"

	// Input errors
	ErrCodeEmptyQuery        ErrorCode = "EMPTY_QUERY"
	ErrCodeInvalidFormat     ErrorCode = "INVALID_FORMAT"
	ErrCodeInvalidURL        ErrorCode = "INVALID_URL"
	ErrCodeInvalidRange      ErrorCode = "INVALID_RANGE"

	// Result errors
	ErrCodeNoResults         ErrorCode = "NO_RESULTS"
	ErrCodeEmptyResults      ErrorCode = "EMPTY_RESULTS"
)

// SearchError is a structured error with user-friendly messages.
//
// It includes an error code, message, suggestion, underlying error, and
// verbose details for debugging.
type SearchError struct {
	Code       ErrorCode
	Message    string
	Suggestion string
	Err        error
	Verbose    string
}

func (e *SearchError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", e.Code, e.Message))
	if e.Suggestion != "" {
		sb.WriteString(fmt.Sprintf("\nSuggestion: %s", e.Suggestion))
	}
	if e.Err != nil {
		sb.WriteString(fmt.Sprintf("\nDetails: %v", e.Err))
	}
	return sb.String()
}

func (e *SearchError) Unwrap() error {
	return e.Err
}

// New creates a new SearchError with the given code and message.
func New(code ErrorCode, message string) *SearchError {
	return &SearchError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context.
func Wrap(code ErrorCode, message string, err error) *SearchError {
	return &SearchError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithSuggestion adds a suggestion to the error
func (e *SearchError) WithSuggestion(suggestion string) *SearchError {
	e.Suggestion = suggestion
	return e
}

// WithVerbose adds verbose information
func (e *SearchError) WithVerbose(verbose string) *SearchError {
	e.Verbose = verbose
	return e
}

// WithErr adds an underlying error
func (e *SearchError) WithErr(err error) *SearchError {
	e.Err = err
	return e
}

// Error constructors for common scenarios

// Config errors
// ConfigNotFound creates an error for missing config files.
func ConfigNotFound(path string) *SearchError {
	return &SearchError{
		Code:       ErrCodeConfigNotFound,
		Message:    fmt.Sprintf("Configuration file not found: %s", path),
		Suggestion: "A default configuration will be used. Run 'search --help' for configuration options.",
	}
}

func ConfigInvalid(err error) *SearchError {
	return &SearchError{
		Code:       ErrCodeConfigInvalid,
		Message:    "Configuration file is invalid",
		Suggestion: "Check your config file syntax at ~/.search/config.yaml",
		Err:        err,
	}
}

func ConfigParseError(line int, err error) *SearchError {
	return &SearchError{
		Code:       ErrCodeConfigParseError,
		Message:    fmt.Sprintf("Failed to parse configuration at line %d", line),
		Suggestion: "Verify YAML syntax in your config file",
		Err:        err,
	}
}

// Network errors
func NetworkError(err error) *SearchError {
	code := ErrCodeNetworkUnreachable
	message := "Network error occurred"
	suggestion := "Check your internet connection and try again"

	// Provide specific error handling based on error type
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		code = ErrCodeNetworkTimeout
		message = "Request timed out"
		suggestion = "The search request took too long. Try increasing --timeout or check your network"
	} else if strings.Contains(err.Error(), "connection refused") {
		code = ErrCodeConnectionRefused
		message = "Connection refused by the server"
		suggestion = "The SearXNG instance may be down. Try a different instance with --instance"
	} else if strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "dns") {
		code = ErrCodeDNSFailed
		message = "DNS lookup failed"
		suggestion = "Check the instance URL and your DNS configuration"
	}

	return &SearchError{
		Code:       code,
		Message:    message,
		Suggestion: suggestion,
		Err:        err,
	}
}

func HTTPStatusError(statusCode int, status string) *SearchError {
	code := ErrCodeAPIError
	message := fmt.Sprintf("HTTP error: %s", status)
	suggestion := "The SearXNG instance returned an error. Try again later or use a different instance"

	if statusCode >= 500 {
		code = ErrCodeAPIUnavailable
		message = fmt.Sprintf("Server error (%d)", statusCode)
		suggestion = "The SearXNG instance is experiencing issues. Try again later"
	} else if statusCode == 404 {
		message = "Search endpoint not found (404)"
		suggestion = "The instance may not support the search API. Check the instance URL"
	} else if statusCode == 401 || statusCode == 403 {
		message = fmt.Sprintf("Access denied (%d)", statusCode)
		suggestion = "The instance may require authentication. Check if an API key is needed"
	}

	return &SearchError{
		Code:       code,
		Message:    message,
		Suggestion: suggestion,
	}
}

// API errors
func InvalidResponse(err error) *SearchError {
	return &SearchError{
		Code:       ErrCodeInvalidResponse,
		Message:    "Invalid response from SearXNG instance",
		Suggestion: "The instance may be misconfigured or incompatible. Try a different instance",
		Err:        err,
	}
}

func APIError(message string) *SearchError {
	return &SearchError{
		Code:       ErrCodeAPIError,
		Message:    fmt.Sprintf("SearXNG API error: %s", message),
		Suggestion: "Check the SearXNG instance logs or try a different instance",
	}
}

// Input validation errors
func EmptyQuery() *SearchError {
	return &SearchError{
		Code:       ErrCodeEmptyQuery,
		Message:    "Query cannot be empty",
		Suggestion: "Provide a search query, e.g., 'search golang tutorials'",
		Verbose:    "Tip: Use quotes for phrases: 'search \"machine learning\"'",
	}
}

func InvalidFormat(format string) *SearchError {
	return &SearchError{
		Code:       ErrCodeInvalidFormat,
		Message:    fmt.Sprintf("Invalid output format: %s", format),
		Suggestion: "Valid formats are: json, markdown, text",
	}
}

func InvalidURL(url string) *SearchError {
	return &SearchError{
		Code:       ErrCodeInvalidURL,
		Message:    fmt.Sprintf("Invalid URL: %s", url),
		Suggestion: "URL must start with http:// or https://",
	}
}

func InvalidRange(param string, min, max, value int) *SearchError {
	return &SearchError{
		Code:       ErrCodeInvalidRange,
		Message:    fmt.Sprintf("Invalid %s: %d (must be between %d and %d)", param, value, min, max),
		Suggestion: fmt.Sprintf("Use a value between %d and %d", min, max),
	}
}

// Result errors
func NoResults(query string) *SearchError {
	return &SearchError{
		Code:       ErrCodeNoResults,
		Message:    fmt.Sprintf("No results found for: %s", query),
		Suggestion: "Try different search terms, check spelling, or remove filters",
	}
}

// IsSearchError checks if an error is a SearchError
func IsSearchError(err error) (*SearchError, bool) {
	if err == nil {
		return nil, false
	}
	if searchErr, ok := err.(*SearchError); ok {
		return searchErr, true
	}
	return nil, false
}

// GetErrorCode returns the error code if available
func GetErrorCode(err error) ErrorCode {
	if searchErr, ok := err.(*SearchError); ok {
		return searchErr.Code
	}
	return ""
}

// HandleError prints error to stderr in a user-friendly format
//
// Example:
//
//	err := errors.EmptyQuery()
//	errors.HandleError(err, false) // basic error
//	errors.HandleError(err, true)  // with verbose details
func HandleError(err error, verbose bool) {
	if err == nil {
		return
	}

	searchErr, ok := err.(*SearchError)
	if !ok {
		// Not a SearchError, print as-is
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Print the error message
	fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m %s\n", searchErr.Message)

	// Print suggestion if available
	if searchErr.Suggestion != "" {
		fmt.Fprintf(os.Stderr, "\033[33m💡 %s\033[0m\n", searchErr.Suggestion)
	}

	// Print verbose details if requested
	if verbose && searchErr.Verbose != "" {
		fmt.Fprintf(os.Stderr, "\n\033[36mVerbose:\033[0m %s\n", searchErr.Verbose)
	}

	// Print underlying error in verbose mode
	if verbose && searchErr.Err != nil {
		fmt.Fprintf(os.Stderr, "\n\033[36mUnderlying error:\033[0m %v\n", searchErr.Err)
	}
}

// WrapNetworkError wraps network errors with appropriate context
func WrapNetworkError(err error) error {
	if err == nil {
		return nil
	}

	// Check for specific error types
	if urlErr, ok := err.(*url.Error); ok {
		return NetworkError(urlErr.Err)
	}

	if netErr, ok := err.(net.Error); ok {
		return NetworkError(netErr)
	}

	// HTTP errors
	if httpErr, ok := err.(*HTTPResponseError); ok {
		return HTTPStatusError(httpErr.StatusCode, httpErr.Status)
	}

	// Generic network error
	return NetworkError(err)
}

// HTTPResponseError represents an HTTP error
type HTTPResponseError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HTTPResponseError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Status)
}

// CheckHTTPError checks if an HTTP response indicates an error
//
// Example:
//
//	resp, err := http.Get("https://example.com")
//	if err != nil {
//	    return err
//	}
//	defer resp.Body.Close()
//	if err := errors.CheckHTTPError(resp); err != nil {
//	    return err
//	}
func CheckHTTPError(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return &HTTPResponseError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(body),
		}
	}
	return nil
}
