// Package searxng provides a client for interacting with SearXNG search instances.
//
// It includes functionality for executing searches, parsing responses, and handling
// various SearXNG features like categories, time range filtering, and pagination.
package searxng

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mule-ai/search/internal/config"
	"github.com/mule-ai/search/internal/errors"
)

const defaultUserAgent = "search-cli/0.1.0"

// Client represents a SearXNG API client.
//
// It encapsulates the HTTP client, instance URL, and authentication credentials
// needed to communicate with a SearXNG instance.
type Client struct {
	instanceURL string
	client      *http.Client
	userAgent   string
	apiKey      string
}

// NewClient creates a new SearXNG client with the given configuration.
//
// The client is configured with the instance URL, timeout, and optional API key
// from the provided Config. Returns a ready-to-use Client instance.
//
// Example:
//
//	cfg := config.DefaultConfig()
//	cfg.Instance = "https://search.butler.ooo"
//	cfg.Timeout = 30
//	client := searxng.NewClient(cfg)
//	resp, err := client.Search(searxng.NewSearchRequest("golang"))
func NewClient(cfg *config.Config) *Client {
	return &Client{
		instanceURL: cfg.Instance,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		userAgent: defaultUserAgent,
		apiKey:    cfg.APIKey,
	}
}

// NewClientWithTimeout creates a new SearXNG client with custom timeout.
//
// This is a convenience function for creating a client with specific timeout settings
// without a full Config object.
func NewClientWithTimeout(instanceURL string, timeout time.Duration) *Client {
	return &Client{
		instanceURL: instanceURL,
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: defaultUserAgent,
	}
}

// Search executes a search query against the SearXNG API.
//
// It builds the appropriate URL with query parameters, executes the HTTP request,
// and parses the JSON response. The SearchResponse includes results, answers,
// infoboxes, and suggestions.
//
// Returns an error if:
//   - The instance URL is invalid
//   - The HTTP request fails
//   - The API returns a non-200 status code
//   - The response JSON cannot be parsed
func (c *Client) Search(req *SearchRequest) (*SearchResponse, error) {
	// Build the URL
	u, err := url.Parse(c.instanceURL)
	if err != nil {
		return nil, errors.InvalidURL(c.instanceURL).WithErr(err)
	}

	// Add the search path if not present
	if u.Path == "" || !strings.HasSuffix(u.Path, "/search") {
		u.Path = strings.TrimSuffix(u.Path, "/") + "/search"
	}

	// Build query parameters
	query := u.Query()
	query.Set("q", req.Query)
	query.Set("format", req.Format)

	// Set page number (1-indexed for SearXNG)
	query.Set("pageno", strconv.Itoa(req.Page))

	// Set language
	if len(req.Languages) > 0 {
		query.Set("language", req.Languages[0])
	}

	// Set safe search
	query.Set("safesearch", strconv.Itoa(req.SafeSearch))

	// Set categories (comma-separated)
	if len(req.Categories) > 0 {
		query.Set("categories", strings.Join(req.Categories, ","))
	}

	// Set engines (comma-separated)
	if len(req.Engines) > 0 {
		query.Set("engines", strings.Join(req.Engines, ","))
	}

	// Set time range
	if req.TimeRange != "" {
		query.Set("time_range", req.TimeRange)
	}

	u.RawQuery = query.Encode()

	// Create HTTP request
	httpReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeAPIError, "failed to create search request", err)
	}

	httpReq.Header.Set("User-Agent", c.userAgent)
	httpReq.Header.Set("Accept", "application/json")

	// Add API key if present
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Execute request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, errors.NetworkError(err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.HTTPStatusError(resp.StatusCode, resp.Status).WithVerbose(fmt.Sprintf("Response body: %s", string(body)))
	}

	// Parse response using optimized decoder
	decoder := NewOptimizedDecoder(resp.Body)
	defer decoder.Close()

	var searchResp SearchResponse
	if err := decoder.Decode(&searchResp); err != nil {
		return nil, errors.InvalidResponse(err)
	}

	// Set search time from headers if available
	if searchTime := resp.Header.Get("X-Response-Time"); searchTime != "" {
		if t, err := strconv.ParseFloat(searchTime, 64); err == nil {
			searchResp.SearchTime = t
		}
	}

	// Set pagination info
	searchResp.Page = req.Page
	searchResp.Instance = c.instanceURL

	return &searchResp, nil
}

// SearchWithConfig executes a search using individual request parameters.
//
// This is a convenience method that creates a SearchRequest from the provided
// parameters and executes it. The format parameter is used for the SearXNG API
// response format (should always be "json" for this client).
func (c *Client) SearchWithConfig(query string, results int, format string, category string, timeout int, language string, safeSearch int, page int, timeRange string) (*SearchResponse, error) {
	req := NewSearchRequest(query)

	// Always use JSON for SearXNG API response (format parameter is for our output formatter)
	req.Format = "json"

	// Set page number (default to 1)
	if page > 0 {
		req.Page = page
	}

	// Set category
	if category != "" {
		req.Categories = []string{category}
	}

	// Set language
	if language != "" {
		req.Languages = []string{language}
	}

	// Set safe search
	req.SafeSearch = safeSearch

	// Set time range
	if timeRange != "" {
		req.TimeRange = timeRange
	}

	return c.Search(req)
}

// ValidateInstance checks if the instance URL is valid and reachable.
//
// It validates the URL format (scheme and host) and attempts a simple
// HTTP GET request to verify connectivity.
//
// Returns an error if the URL is malformed or the instance is unreachable.
func (c *Client) ValidateInstance() error {
	u, err := url.Parse(c.instanceURL)
	if err != nil {
		return errors.InvalidURL(c.instanceURL).WithErr(err)
	}

	// Check scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.InvalidURL(c.instanceURL).WithVerbose("URL scheme must be http or https")
	}

	// Check host
	if u.Host == "" {
		return errors.InvalidURL(c.instanceURL).WithVerbose("URL must include a host")
	}

	// Try to fetch the root endpoint
	httpReq, err := http.NewRequest(http.MethodGet, c.instanceURL, nil)
	if err != nil {
		return errors.Wrap(errors.ErrCodeAPIError, "failed to create validation request", err)
	}

	httpReq.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return errors.NetworkError(err)
	}
	defer resp.Body.Close()

	return nil
}

// ParseResults parses raw JSON response bytes into a SearchResult slice.
//
// This is useful for processing responses that have already been fetched.
// Returns an error if the JSON is invalid or doesn't match the expected format.
func ParseResults(data []byte) ([]SearchResult, error) {
	var resp SearchResponse
	if err := UnmarshalResponse(data, &resp); err != nil {
		return nil, errors.InvalidResponse(err)
	}
	return resp.Results, nil
}

// ParseResponse parses raw JSON response bytes into a SearchResponse.
//
// This is useful for processing responses that have already been fetched.
// Returns an error if the JSON is invalid or doesn't match the expected format.
func ParseResponse(data []byte) (*SearchResponse, error) {
	var resp SearchResponse
	if err := UnmarshalResponse(data, &resp); err != nil {
		return nil, errors.InvalidResponse(err)
	}
	return &resp, nil
}

// NewSearchFromQuery creates a search request from a query string with optional modifiers.
//
// This function uses the functional options pattern for flexible request configuration.
// Example:
//
//	req := NewSearchFromQuery("golang",
//	    WithPage(2),
//	    WithSafeSearch(2),
//	    WithTimeRange("week"))
func NewSearchFromQuery(query string, options ...func(*SearchRequest)) *SearchRequest {
	req := NewSearchRequest(query)
	for _, opt := range options {
		opt(req)
	}
	return req
}

// WithPage sets the page number for the search request.
//
// Pages are 1-indexed in the SearXNG API.
func WithPage(page int) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.Page = page
	}
}

// WithFormat sets the response format for the search request.
//
// Valid values are "json", "rss", or "csv". For this client, "json" is recommended.
func WithFormat(format string) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.Format = format
	}
}

// WithCategories sets the search categories for the request.
//
// Categories include: general, images, videos, news, map, music, it, science, files, etc.
// Multiple categories can be specified.
func WithCategories(categories ...string) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.Categories = categories
	}
}

// WithEngines sets the specific search engines to use.
//
// This limits the search to only the specified engines (e.g., "google", "duckduckgo").
// Multiple engines can be specified.
func WithEngines(engines ...string) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.Engines = engines
	}
}

// WithTimeRange sets the time range filter for the search.
//
// Valid values are: "day", "week", "month", "year".
func WithTimeRange(timeRange string) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.TimeRange = timeRange
	}
}

// WithSafeSearch sets the safe search level for the request.
//
// Valid values are: 0 (off), 1 (moderate), 2 (strict).
func WithSafeSearch(level int) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.SafeSearch = level
	}
}

// WithLanguage sets the language for the search request.
//
// Uses standard language codes (e.g., "en", "de", "fr").
func WithLanguage(language string) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.Languages = []string{language}
	}
}

// WithTimeout sets the request timeout duration.
//
// If not set, the default timeout from the client's HTTP client is used.
func WithTimeout(timeout time.Duration) func(*SearchRequest) {
	return func(req *SearchRequest) {
		req.Timeout = timeout
	}
}

// Do executes a search request and returns the response.
//
// This is an alias for Search() provided for convenience.
func (c *Client) Do(req *SearchRequest) (*SearchResponse, error) {
	return c.Search(req)
}

// GetInstance returns the configured SearXNG instance URL.
func (c *Client) GetInstance() string {
	return c.instanceURL
}

// SetUserAgent sets the user agent string for HTTP requests.
func (c *Client) SetUserAgent(ua string) {
	c.userAgent = ua
}

// GetUserAgent returns the current user agent string.
func (c *Client) GetUserAgent() string {
	return c.userAgent
}

// SetAPIKey sets the API key for authentication.
//
// If set, the key will be sent as a Bearer token in the Authorization header.
func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}

// GetAPIKey returns the current API key.
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// IsReachable checks if the SearXNG instance is reachable.
//
// Returns true if ValidateInstance() succeeds, false otherwise.
func (c *Client) IsReachable() bool {
	err := c.ValidateInstance()
	return err == nil
}
