package searxng

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mule-ai/search/internal/config"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		instance string
		timeout  time.Duration
		wantErr  bool
	}{
		{
			name:     "valid client",
			instance: "https://search.butler.ooo",
			timeout:  30 * time.Second,
			wantErr:  false,
		},
		{
			name:     "valid client with custom timeout",
			instance: "https://search.example.com",
			timeout:  60 * time.Second,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Instance: tt.instance,
				Timeout:  int(tt.timeout.Seconds()),
			}
			client := NewClient(cfg)
			if client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	tests := []struct {
		name     string
		instance string
		timeout  time.Duration
		wantErr  bool
	}{
		{
			name:     "valid client",
			instance: "https://search.butler.ooo",
			timeout:  30 * time.Second,
			wantErr:  false,
		},
		{
			name:     "valid client with custom timeout",
			instance: "https://search.example.com",
			timeout:  60 * time.Second,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClientWithTimeout(tt.instance, tt.timeout)
			if client == nil {
				t.Error("NewClientWithTimeout() returned nil client")
			}
		})
	}
}

func TestClientSearch(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify query parameter
		query := r.URL.Query().Get("q")
		if query != "test query" {
			t.Errorf("Expected query 'test query', got '%s'", query)
		}

		// Verify format parameter
		format := r.URL.Query().Get("format")
		if format != "json" {
			t.Errorf("Expected format 'json', got '%s'", format)
		}

		// Return mock response
		response := SearchResponse{
			Query:           "test query",
			Results: []SearchResult{
				{
					Title:    "Test Result",
					URL:      "https://example.com",
					Content:  "Test content",
					Engine:   "google",
					Category: "general",
					Score:    0.95,
				},
			},
			Answers:          []Answer{{Answer: "Test answer"}},
			Infoboxes:        []Infobox{{Infobox: "Test infobox"}},
			Suggestions:      []string{"suggestion1"},
			NumberOfResults:  100,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Create client with test server URL
	cfg := &config.Config{
		Instance: ts.URL,
		Timeout:  30,
	}
	client := NewClient(cfg)

	// Perform search
	request := NewSearchRequest("test query")
	request.Format = "json"
	request.Page = 1

	response, err := client.Search(request)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	// Verify response
	if response.Query != "test query" {
		t.Errorf("Expected query 'test query', got '%s'", response.Query)
	}

	if len(response.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(response.Results))
	}

	if response.Results[0].Title != "Test Result" {
		t.Errorf("Expected title 'Test Result', got '%s'", response.Results[0].Title)
	}

	if response.NumberOfResults != 100 {
		t.Errorf("Expected 100 results, got %d", response.NumberOfResults)
	}
}

func TestClientSearchWithParameters(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query().Get("q")
		if query != "advanced query" {
			t.Errorf("Expected query 'advanced query', got '%s'", query)
		}

		page := r.URL.Query().Get("pageno")
		if page != "2" {
			t.Errorf("Expected page '2', got '%s'", page)
		}

		language := r.URL.Query().Get("language")
		if language != "de" {
			t.Errorf("Expected language 'de', got '%s'", language)
		}

		safeSearch := r.URL.Query().Get("safesearch")
		if safeSearch != "0" {
			t.Errorf("Expected safesearch '0', got '%s'", safeSearch)
		}

		categories := r.URL.Query().Get("categories")
		if categories != "images" {
			t.Errorf("Expected categories 'images', got '%s'", categories)
		}

		// Return mock response
		response := SearchResponse{
			Query:           "advanced query",
			Results:         []SearchResult{},
			NumberOfResults: 50,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Create client
	cfg := &config.Config{
		Instance: ts.URL,
		Timeout:  30,
	}
	client := NewClient(cfg)

	// Perform search with parameters
	request := NewSearchRequest("advanced query")
	request.Format = "json"
	request.Page = 2
	request.Languages = []string{"de"}
	request.SafeSearch = 0
	request.Categories = []string{"images"}

	response, err := client.Search(request)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if response.Query != "advanced query" {
		t.Errorf("Expected query 'advanced query', got '%s'", response.Query)
	}
}

func TestClientSearchError(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse int
		serverBody     string
		wantErr        bool
	}{
		{
			name:           "server error 500",
			serverResponse: 500,
			serverBody:     "Internal Server Error",
			wantErr:        true,
		},
		{
			name:           "server error 404",
			serverResponse: 404,
			serverBody:     "Not Found",
			wantErr:        true,
		},
		{
			name:           "invalid JSON",
			serverResponse: 200,
			serverBody:     "{ invalid json }",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverResponse)
				w.Write([]byte(tt.serverBody))
			}))
			defer ts.Close()

			// Create client
			cfg := &config.Config{
				Instance: ts.URL,
				Timeout:  30,
			}
			client := NewClient(cfg)

			// Perform search
			request := NewSearchRequest("test query")
			_, err := client.Search(request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientTimeout(t *testing.T) {
	// Create a test server that delays response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	// Create client with short timeout (50ms should trigger timeout)
	cfg := &config.Config{
		Instance: ts.URL,
		Timeout:  1, // 1 second timeout for a 2-second server response
	}
	client := NewClient(cfg)

	// Perform search
	request := NewSearchRequest("test query")
	_, err := client.Search(request)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestNewSearchRequest(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  *SearchRequest
	}{
		{
			name:  "basic request",
			query: "test query",
			want: &SearchRequest{
				Query:      "test query",
				Format:     "json",
				Page:       1,
				Languages:  []string{"en"},
				SafeSearch: 1,
				Categories: []string{"general"},
				Timeout:    30 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSearchRequest(tt.query)

			if got.Query != tt.want.Query {
				t.Errorf("Query = %v, want %v", got.Query, tt.want.Query)
			}

			if got.Format != tt.want.Format {
				t.Errorf("Format = %v, want %v", got.Format, tt.want.Format)
			}

			if got.Page != tt.want.Page {
				t.Errorf("Page = %v, want %v", got.Page, tt.want.Page)
			}

			if len(got.Languages) != len(tt.want.Languages) {
				t.Errorf("Languages length = %v, want %v", len(got.Languages), len(tt.want.Languages))
			}

			if got.SafeSearch != tt.want.SafeSearch {
				t.Errorf("SafeSearch = %v, want %v", got.SafeSearch, tt.want.SafeSearch)
			}
		})
	}
}

func TestClientWithCustomHeaders(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for User-Agent header
		userAgent := r.Header.Get("User-Agent")
		if !strings.Contains(userAgent, "search-cli") {
			t.Errorf("Expected User-Agent to contain 'search-cli', got '%s'", userAgent)
		}

		// Return mock response
		response := SearchResponse{
			Query:   "test",
			Results: []SearchResult{},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Create client
	cfg := &config.Config{
		Instance: ts.URL,
		Timeout:  30,
	}
	client := NewClient(cfg)

	// Perform search
	request := NewSearchRequest("test")
	_, err := client.Search(request)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
}

func TestClientRequestURL(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path
		if r.URL.Path != "/search" {
			t.Errorf("Expected path '/search', got '%s'", r.URL.Path)
		}

		// Return mock response
		response := SearchResponse{
			Query:   "test",
			Results: []SearchResult{},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Create client with URL without /search
	cfg := &config.Config{
		Instance: ts.URL,
		Timeout:  30,
	}
	client := NewClient(cfg)

	// Perform search
	request := NewSearchRequest("test")
	_, err := client.Search(request)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("empty query", func(t *testing.T) {
		request := NewSearchRequest("")
		if request.Query != "" {
			t.Errorf("Expected empty query, got '%s'", request.Query)
		}
	})

	t.Run("special characters in query", func(t *testing.T) {
		specialQuery := "test & query with \"quotes\" and <tags>"
		request := NewSearchRequest(specialQuery)

		if request.Query != specialQuery {
			t.Errorf("Expected query '%s', got '%s'", specialQuery, request.Query)
		}
	})

	t.Run("unicode in query", func(t *testing.T) {
		unicodeQuery := "test 查询 検索"
		request := NewSearchRequest(unicodeQuery)

		if request.Query != unicodeQuery {
			t.Errorf("Expected query '%s', got '%s'", unicodeQuery, request.Query)
		}
	})
}

// Benchmark tests
func BenchmarkClientSearch(b *testing.B) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := SearchResponse{
			Query: "test",
			Results: []SearchResult{
				{
					Title:   "Test Result",
					URL:     "https://example.com",
					Content: "Test content",
					Engine:  "google",
					Score:   0.95,
				},
			},
			NumberOfResults: 100,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Instance: ts.URL,
		Timeout:  30,
	}
	client := NewClient(cfg)
	request := NewSearchRequest("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.Search(request)
	}
}

// Helper function to read response body
func readBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(body), nil
}

// Tests for builder methods
func TestSearchRequestBuilder(t *testing.T) {
	t.Run("WithPage", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithPage(3))
		if request.Page != 3 {
			t.Errorf("WithPage() = %v, want 3", request.Page)
		}
	})

	t.Run("WithFormat", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithFormat("json"))
		if request.Format != "json" {
			t.Errorf("WithFormat() = %v, want 'json'", request.Format)
		}
	})

	t.Run("WithCategories", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithCategories("images", "videos"))
		if len(request.Categories) != 2 {
			t.Errorf("WithCategories() length = %v, want 2", len(request.Categories))
		}
		if request.Categories[0] != "images" {
			t.Errorf("WithCategories()[0] = %v, want 'images'", request.Categories[0])
		}
	})

	t.Run("WithEngines", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithEngines("google", "bing"))
		if len(request.Engines) != 2 {
			t.Errorf("WithEngines() length = %v, want 2", len(request.Engines))
		}
		if request.Engines[0] != "google" {
			t.Errorf("WithEngines()[0] = %v, want 'google'", request.Engines[0])
		}
	})

	t.Run("WithTimeRange", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithTimeRange("week"))
		if request.TimeRange != "week" {
			t.Errorf("WithTimeRange() = %v, want 'week'", request.TimeRange)
		}
	})

	t.Run("WithSafeSearch", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithSafeSearch(2))
		if request.SafeSearch != 2 {
			t.Errorf("WithSafeSearch() = %v, want 2", request.SafeSearch)
		}
	})

	t.Run("WithLanguage", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithLanguage("de"))
		if len(request.Languages) != 1 {
			t.Errorf("WithLanguage() length = %v, want 1", len(request.Languages))
		}
		if request.Languages[0] != "de" {
			t.Errorf("WithLanguage()[0] = %v, want 'de'", request.Languages[0])
		}
	})

	t.Run("WithTimeout", func(t *testing.T) {
		request := NewSearchFromQuery("test query", WithTimeout(60*time.Second))
		if request.Timeout != 60*time.Second {
			t.Errorf("WithTimeout() = %v, want 60s", request.Timeout)
		}
	})
}

// Test client getter and setter methods
func TestClientGettersSetters(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := NewClient(cfg)

	t.Run("GetInstance", func(t *testing.T) {
		instance := client.GetInstance()
		if instance != "https://search.butler.ooo" {
			t.Errorf("GetInstance() = %v, want 'https://search.butler.ooo'", instance)
		}
	})

	t.Run("SetUserAgent", func(t *testing.T) {
		customUA := "MyCustomAgent/1.0"
		client.SetUserAgent(customUA)
		ua := client.GetUserAgent()
		if !strings.Contains(ua, customUA) {
			t.Errorf("GetUserAgent() = %v, want to contain %v", ua, customUA)
		}
	})

	t.Run("SetAPIKey", func(t *testing.T) {
		apiKey := "test-api-key-123"
		client.SetAPIKey(apiKey)
		key := client.GetAPIKey()
		if key != apiKey {
			t.Errorf("GetAPIKey() = %v, want %v", key, apiKey)
		}
	})
}

// Test ValidateInstance
func TestValidateInstance(t *testing.T) {
	t.Run("valid HTTPS instance", func(t *testing.T) {
		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		cfg := &config.Config{
			Instance: ts.URL,
			Timeout:  30,
		}
		client := NewClient(cfg)
		err := client.ValidateInstance()
		if err != nil {
			t.Errorf("ValidateInstance() error = %v, wantErr false", err)
		}
	})

	t.Run("invalid URL - no scheme", func(t *testing.T) {
		cfg := &config.Config{
			Instance: "search.example.com",
			Timeout:  30,
		}
		client := NewClient(cfg)
		err := client.ValidateInstance()
		if err == nil {
			t.Error("ValidateInstance() expected error for URL without scheme")
		}
	})

	t.Run("invalid URL - spaces", func(t *testing.T) {
		cfg := &config.Config{
			Instance: "https://search example.com",
			Timeout:  30,
		}
		client := NewClient(cfg)
		err := client.ValidateInstance()
		if err == nil {
			t.Error("ValidateInstance() expected error for URL with spaces")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		cfg := &config.Config{
			Instance: "",
			Timeout:  30,
		}
		client := NewClient(cfg)
		err := client.ValidateInstance()
		if err == nil {
			t.Error("ValidateInstance() expected error for empty URL")
		}
	})
}

// Test IsReachable
func TestIsReachable(t *testing.T) {
	t.Run("reachable instance", func(t *testing.T) {
		// Create a test server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		cfg := &config.Config{
			Instance: ts.URL,
			Timeout:  5,
		}
		client := NewClient(cfg)

		reachable := client.IsReachable()
		if !reachable {
			t.Error("IsReachable() = false, want true")
		}
	})

	t.Run("unreachable instance", func(t *testing.T) {
		// Use a non-routable IP address
		cfg := &config.Config{
			Instance: "http://192.0.2.1:9999", // TEST-NET-1, should be unreachable
			Timeout:  1,
		}
		client := NewClient(cfg)

		reachable := client.IsReachable()
		if reachable {
			t.Error("IsReachable() = true, expected false for unreachable instance")
		}
	})
}

// Test ParseResults
func TestParseResults(t *testing.T) {
	t.Run("parse valid JSON results", func(t *testing.T) {
		jsonData := `{
			"query": "test",
			"results": [
				{
					"title": "Test Result",
					"url": "https://example.com",
					"content": "Test content",
					"engine": "google",
					"category": "general",
					"score": 0.95
				}
			],
			"number_of_results": 100
		}`

		results, err := ParseResults([]byte(jsonData))
		if err != nil {
			t.Errorf("ParseResults() error = %v", err)
		}
		if len(results) != 1 {
			t.Errorf("ParseResults() length = %v, want 1", len(results))
		}
		if results[0].Title != "Test Result" {
			t.Errorf("ParseResults() title = %v, want 'Test Result'", results[0].Title)
		}
	})

	t.Run("parse invalid JSON", func(t *testing.T) {
		invalidJSON := `{ invalid json }`

		_, err := ParseResults([]byte(invalidJSON))
		if err == nil {
			t.Error("ParseResults() expected error for invalid JSON, got nil")
		}
	})
}

// Test ParseResponse
func TestParseResponse(t *testing.T) {
	t.Run("parse HTTP response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := SearchResponse{
				Query: "test",
				Results: []SearchResult{
					{
						Title:    "Test Result",
						URL:      "https://example.com",
						Content:  "Test content",
						Engine:   "google",
						Category: "general",
						Score:    0.95,
					},
				},
				NumberOfResults: 100,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer ts.Close()

		// Make a direct HTTP request
		resp, err := http.Get(ts.URL + "/search?q=test&format=json")
		if err != nil {
			t.Fatalf("Failed to make test request: %v", err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		parsed, err := ParseResponse(body)
		if err != nil {
			t.Errorf("ParseResponse() error = %v", err)
		}
		if parsed.Query != "test" {
			t.Errorf("ParseResponse() query = %v, want 'test'", parsed.Query)
		}
		if len(parsed.Results) != 1 {
			t.Errorf("ParseResponse() length = %v, want 1", len(parsed.Results))
		}
	})
}

// Test NewSearchFromQuery
func TestNewSearchFromQuery(t *testing.T) {
	query := "test query"
	request := NewSearchFromQuery(query)

	if request.Query != query {
		t.Errorf("NewSearchFromQuery() query = %v, want %v", request.Query, query)
	}
	if request.Format != "json" {
		t.Errorf("NewSearchFromQuery() format = %v, want 'json'", request.Format)
	}
}

// Test SearchWithConfig
func TestSearchWithConfig(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify parameters from config
		query := r.URL.Query().Get("q")
		if query != "config test" {
			t.Errorf("Expected query 'config test', got '%s'", query)
		}

		language := r.URL.Query().Get("language")
		if language != "de" {
			t.Errorf("Expected language 'de', got '%s'", language)
		}

		safeSearch := r.URL.Query().Get("safesearch")
		if safeSearch != "2" {
			t.Errorf("Expected safesearch '2', got '%s'", safeSearch)
		}

		categories := r.URL.Query().Get("categories")
		if categories != "images" {
			t.Errorf("Expected categories 'images', got '%s'", categories)
		}

		response := SearchResponse{
			Query:           "config test",
			Results:         []SearchResult{},
			NumberOfResults: 50,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Instance:   ts.URL,
		Timeout:    30,
		Language:   "de",
		SafeSearch: 2,
		Categories: []string{"images"},
	}

	client := NewClient(cfg)

	response, err := client.SearchWithConfig("config test", 10, "json", "images", 30, "de", 2, 1, "")
	if err != nil {
		t.Fatalf("SearchWithConfig() error = %v", err)
	}

	if response.Query != "config test" {
		t.Errorf("SearchWithConfig() query = %v, want 'config test'", response.Query)
	}
}
