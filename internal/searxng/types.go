// Package searxng provides data structures for SearXNG API requests and responses.
//
// These types map to the JSON format used by SearXNG instances.
package searxng

import (
	"encoding/json"
	"time"
)

// SearchResult represents a single search result from SearXNG.
//
// It includes the result title, URL, content snippet, source engine,
// category, relevance score, and optional image source.
//
// Example:
//
//	result := searxng.SearchResult{
//	    Title:   "A Tour of Go",
//	    URL:     "https://go.dev/tour/",
//	    Content: "Welcome to a tour of the Go programming language...",
//	    Engine:  "google",
//	    Category: "general",
//	    Score:   0.95,
//	}
type SearchResult struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Content     string   `json:"content"`
	Engine      string   `json:"engine"`
	Category    string   `json:"category"`
	Score       float64  `json:"score"`
	ImgSrc      string   `json:"img_src,omitempty"`
	ParsedURL   []string `json:"parsed_url,omitempty"`
	Template    string   `json:"template,omitempty"`
	PublishedDate *time.Time `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling for SearchResult.
//
// This handles edge cases like missing fields and provides default values.
func (sr *SearchResult) UnmarshalJSON(data []byte) error {
	// Use type alias to avoid recursion
	type Alias SearchResult
	aux := &struct {
		Score interface{} `json:"score"`
		*Alias
	}{
		Alias: (*Alias)(sr),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle score which can be number or missing
	switch v := aux.Score.(type) {
	case float64:
		sr.Score = v
	case nil:
		sr.Score = 0.0
	}

	// Ensure slices are initialized
	if sr.ParsedURL == nil {
		sr.ParsedURL = []string{}
	}

	return nil
}

// Answer represents a direct answer from SearXNG.
//
// SearXNG sometimes returns instant answers (e.g., calculations, definitions)
// separate from the search results.
//
// Example:
//
//	answer := searxng.Answer{
//	    Answer: "42",
//	    Engine: "answer",
//	    URL:    "https://en.wikipedia.org/wiki/42_(number)",
//	}
type Answer struct {
	Answer      string   `json:"answer"`
	URL         string   `json:"url,omitempty"`
	Engine      string   `json:"engine,omitempty"`
	ParsedURL   []string `json:"parsed_url,omitempty"`
	Template    string   `json:"template,omitempty"`
}

// Infobox represents an infobox from SearXNG.
//
// Infoboxes are structured information panels (like Wikipedia sidebars)
// that provide detailed information about a topic.
//
// Example:
//
//	infobox := searxng.Infobox{
//	    Infobox: "Golang",
//	    Content: "Go is a statically typed, compiled programming language...",
//	    Attributes: []searxng.Attribute{
//	        {Label: "Designed by", Value: "Robert Griesemer, Rob Pike, Ken Thompson"},
//	        {Label: "First appeared", Value: "2009"},
//	    },
//	}
type Infobox struct {
	Infobox     string   `json:"infobox"`
	ID          string   `json:"id,omitempty"`
	Content     string   `json:"content,omitempty"`
	ImgSrc      string   `json:"img_src,omitempty"`
	URLs        []URLInfo `json:"urls,omitempty"`
	Engine      string   `json:"engine,omitempty"`
	URL         string   `json:"url,omitempty"`
	Template    string   `json:"template,omitempty"`
	ParsedURL   []string `json:"parsed_url,omitempty"`
	Title       string   `json:"title,omitempty"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Engines     []string `json:"engines,omitempty"`
	Positions   string   `json:"positions,omitempty"`
	Score       float64  `json:"score,omitempty"`
	Category    string   `json:"category,omitempty"`
	PublishedDate *time.Time `json:"publishedDate,omitempty"`
	Attributes  []Attribute `json:"attributes,omitempty"`
}

// URLInfo represents a URL within an infobox.
//
// It includes the URL, display title, and an optional flag indicating
// if this is the official/primary source.
//
// Example:
//
//	urlInfo := searxng.URLInfo{
//	    Title:   "Official Website",
//	    URL:     "https://go.dev",
//	    Official: true,
//	}
type URLInfo struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Official bool  `json:"official,omitempty"`
}

// Attribute represents a key-value pair within an infobox.
//
// Infoboxes use attributes to display structured data such as
// "Founded: 2008" or "CEO: John Doe".
//
// Example:
//
//	attr := searxng.Attribute{
//	    Label: "Founded",
//	    Value: "2008",
//	}
type Attribute struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Entity string `json:"entity,omitempty"`
}

// SearchResponse represents the complete response from a SearXNG API query.
//
// It includes search results, answers, infoboxes, suggestions, and metadata
// about the search (query, result count, timing, pagination).
//
// Example:
//
//	response := &searxng.SearchResponse{
//	    Query: "golang",
//	    Results: []searxng.SearchResult{...},
//	    NumberOfResults: 1250000,
//	    SearchTime: 0.24,
//	    Answers: []searxng.Answer{...},
//	    Suggestions: []string{"golang tutorial", "golang jobs"},
//	}
type SearchResponse struct {
	Query            string         `json:"query"`
	Results          []SearchResult `json:"results"`
	Answers          []Answer       `json:"answers"`
	Infoboxes        []Infobox      `json:"infoboxes"`
	Suggestions      []string       `json:"suggestions"`
	Corrections      []string       `json:"corrections,omitempty"`
	UnresponsiveEngines [][]string   `json:"unresponsive_engines,omitempty"`
	// Number of results from the API
	NumberOfResults int    `json:"number_of_results"`
	// Search time in seconds
	SearchTime float64 `json:"-"`
	// Pagination info
	Page     int    `json:"page,omitempty"`
	Instance string `json:"-"` // Instance URL for display
}

// UnmarshalJSON implements custom JSON unmarshaling for SearchResponse.
//
// This implementation validates the response and handles edge cases like
// missing or malformed data from SearXNG instances.
func (sr *SearchResponse) UnmarshalJSON(data []byte) error {
	// Use type alias to avoid recursion
	type Alias SearchResponse
	aux := &struct {
		NumberOfResults interface{} `json:"number_of_results"`
		*Alias
	}{
		Alias: (*Alias)(sr),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle number_of_results which can be string or number
	switch v := aux.NumberOfResults.(type) {
	case float64:
		sr.NumberOfResults = int(v)
	case string:
		// Some instances return number as string
		var n int
		if err := json.Unmarshal([]byte(v), &n); err == nil {
			sr.NumberOfResults = n
		}
	case nil:
		// If not present, default to length of results
		sr.NumberOfResults = len(sr.Results)
	}

	// Ensure slices are initialized to avoid nil panics
	if sr.Results == nil {
		sr.Results = []SearchResult{}
	}
	if sr.Answers == nil {
		sr.Answers = []Answer{}
	}
	if sr.Infoboxes == nil {
		sr.Infoboxes = []Infobox{}
	}
	if sr.Suggestions == nil {
		sr.Suggestions = []string{}
	}
	if sr.Corrections == nil {
		sr.Corrections = []string{}
	}
	if sr.UnresponsiveEngines == nil {
		sr.UnresponsiveEngines = [][]string{}
	}

	return nil
}

// SearchRequest represents a search request to the SearXNG API.
//
// It contains all parameters that can be sent to the /search endpoint,
// including the query, format, pagination, filters, and engine selection.
//
// Example:
//
//	req := &searxng.SearchRequest{
//	    Query: "golang",
//	    Format: "json",
//	    Page: 1,
//	    Languages: []string{"en"},
//	    SafeSearch: 1,
//	    Categories: []string{"general"},
//	}
//	client := searxng.NewClient(cfg)
//	resp, err := client.Search(req)
type SearchRequest struct {
	Query       string
	Format      string // "json", "rss", "csv"
	Page        int
	Languages   []string
	SafeSearch  int
	Categories  []string
	Engines     []string
	TimeRange   string
	Timeout     time.Duration
}

// NewSearchRequest creates a new search request with default values.
//
// The defaults are:
//   - Format: "json"
//   - Page: 1
//   - Languages: ["en"]
//   - SafeSearch: 1 (moderate)
//   - Categories: ["general"]
//   - Timeout: 30 seconds
//
// The query parameter is required and should not be empty.
func NewSearchRequest(query string) *SearchRequest {
	return &SearchRequest{
		Query:      query,
		Format:     "json",
		Page:       1,
		Languages:  []string{"en"},
		SafeSearch: 1,
		Categories: []string{"general"},
		Timeout:    30 * time.Second,
	}
}
