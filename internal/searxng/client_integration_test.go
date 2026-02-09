//go:build integration
// +build integration

package searxng

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mule-ai/search/internal/config"
)

// TestIntegrationClient tests the SearXNG client against a real instance.
// To run these tests, use: go test -tags=integration -v ./internal/searxng/...
func TestIntegrationClient(t *testing.T) {
	// Get test instance from environment or use default
	instanceURL := os.Getenv("SEARXNG_TEST_INSTANCE")
	if instanceURL == "" {
		instanceURL = "https://search.butler.ooo"
	}

	// Create config for the client
	cfg := &config.Config{
		Instance: instanceURL,
		Timeout:  30,
	}

	t.Run("BasicSearch", func(t *testing.T) {
		client := NewClient(cfg)

		req := NewSearchRequest("golang")
		req.Format = "json"
		req.Page = 1
		req.Languages = []string{"en"}

		resp, err := client.Search(req)
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}

		if resp.Query != "golang" {
			t.Errorf("Expected query 'golang', got '%s'", resp.Query)
		}

		if len(resp.Results) == 0 {
			t.Error("Expected at least one result")
		}

		// Verify first result has required fields
		if len(resp.Results) > 0 {
			result := resp.Results[0]
			if result.Title == "" {
				t.Error("Expected result to have a title")
			}
			if result.URL == "" {
				t.Error("Expected result to have a URL")
			}
		}
	})

	t.Run("SearchWithCategory", func(t *testing.T) {
		client := NewClient(cfg)

		req := NewSearchRequest("nature")
		req.Format = "json"
		req.Page = 1
		req.Categories = []string{"images"}
		req.Languages = []string{"en"}

		resp, err := client.Search(req)
		if err != nil {
			t.Fatalf("Search with category failed: %v", err)
		}

		if len(resp.Results) == 0 {
			t.Log("Warning: No image results returned (may be normal for some queries)")
		}

		// Image results should have img_src
		for _, result := range resp.Results {
			if result.ImgSrc != "" {
				return // Success - found image result
			}
		}
	})

	t.Run("SearchWithTimeRange", func(t *testing.T) {
		client := NewClient(cfg)

		req := NewSearchRequest("news")
		req.Format = "json"
		req.Page = 1
		req.Languages = []string{"en"}
		req.TimeRange = "day"

		resp, err := client.Search(req)
		if err != nil {
			t.Fatalf("Search with time range failed: %v", err)
		}

		// Just verify it doesn't error - time range filtering is instance-dependent
		_ = resp
	})

	t.Run("SearchWithPagination", func(t *testing.T) {
		client := NewClient(cfg)

		// Search page 1
		req1 := NewSearchRequest("programming")
		req1.Format = "json"
		req1.Page = 1
		req1.Languages = []string{"en"}

		resp1, err := client.Search(req1)
		if err != nil {
			t.Fatalf("Search page 1 failed: %v", err)
		}

		// Search page 2
		req2 := NewSearchRequest("programming")
		req2.Format = "json"
		req2.Page = 2
		req2.Languages = []string{"en"}

		resp2, err := client.Search(req2)
		if err != nil {
			t.Fatalf("Search page 2 failed: %v", err)
		}

		// Results may differ, but both should succeed
		t.Logf("Page 1: %d results, Page 2: %d results", len(resp1.Results), len(resp2.Results))
	})

	t.Run("SearchWithSafeSearch", func(t *testing.T) {
		client := NewClient(cfg)

		// Test safe search level 2 (strict)
		req := NewSearchRequest("test")
		req.Format = "json"
		req.Page = 1
		req.SafeSearch = 2

		resp, err := client.Search(req)
		if err != nil {
			t.Fatalf("Search with safe search failed: %v", err)
		}

		// Just verify it doesn't error
		_ = resp
	})

	t.Run("SearchWithTimeout", func(t *testing.T) {
		// Create client with short timeout
		timeoutCfg := &config.Config{
			Instance: instanceURL,
			Timeout:  1,
		}
		client := NewClient(timeoutCfg)

		req := NewSearchRequest("test")
		req.Format = "json"
		req.Page = 1

		_, err := client.Search(req)
		if err != nil {
			// Timeout is acceptable if the instance is slow
			t.Logf("Search with timeout returned error (may be expected): %v", err)
		}
	})
}

// TestIntegrationMockServer tests the client against a local mock server.
// This ensures the client correctly parses responses.
func TestIntegrationMockServer(t *testing.T) {
	// Create a mock SearXNG server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		if r.URL.Path != "/search" {
			t.Errorf("Expected path /search, got %s", r.URL.Path)
		}

		query := r.URL.Query().Get("q")
		_ = r.URL.Query().Get("format") // format is read but not used in verification

		// Return mock response
		mockResp := map[string]interface{}{
			"query": query,
			"results": []map[string]interface{}{
				{
					"title":       "Test Result",
					"url":         "https://example.com/test",
					"content":     "This is a test result snippet",
					"engine":      "test",
					"category":    "general",
					"score":       0.95,
					"parsed_url":  []string{"https://example.com", "test"},
					"template":    "default.html",
					"engines":     []string{"test"},
				},
			},
			"answers":           []interface{}{},
			"infoboxes":         []interface{}{},
			"suggestions":       []string{"test suggestion"},
			"number_of_results": 1000,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	t.Run("MockServerSearch", func(t *testing.T) {
		cfg := &config.Config{
			Instance: server.URL,
			Timeout:  30,
		}
		client := NewClient(cfg)

		req := NewSearchRequest("test query")
		req.Format = "json"
		req.Page = 1

		resp, err := client.Search(req)
		if err != nil {
			t.Fatalf("Mock server search failed: %v", err)
		}

		if resp.Query != "test query" {
			t.Errorf("Expected query 'test query', got '%s'", resp.Query)
		}

		if len(resp.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(resp.Results))
		}

		result := resp.Results[0]
		if result.Title != "Test Result" {
			t.Errorf("Expected title 'Test Result', got '%s'", result.Title)
		}
		if result.URL != "https://example.com/test" {
			t.Errorf("Expected URL 'https://example.com/test', got '%s'", result.URL)
		}
		if result.Engine != "test" {
			t.Errorf("Expected engine 'test', got '%s'", result.Engine)
		}
		if result.Score != 0.95 {
			t.Errorf("Expected score 0.95, got %f", result.Score)
		}
	})

	t.Run("MockServerEmptyResults", func(t *testing.T) {
		// Create a server that returns empty results
		emptyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mockResp := map[string]interface{}{
				"query":              "empty",
				"results":            []interface{}{},
				"answers":            []interface{}{},
				"infoboxes":          []interface{}{},
				"suggestions":        []string{},
				"number_of_results":  0,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockResp)
		}))
		defer emptyServer.Close()

		cfg := &config.Config{
			Instance: emptyServer.URL,
			Timeout:  30,
		}
		client := NewClient(cfg)

		req := NewSearchRequest("empty")
		req.Format = "json"
		req.Page = 1

		resp, err := client.Search(req)
		if err != nil {
			t.Fatalf("Empty results search failed: %v", err)
		}

		if len(resp.Results) != 0 {
			t.Errorf("Expected 0 results, got %d", len(resp.Results))
		}
	})

	t.Run("MockServerError", func(t *testing.T) {
		// Create a server that returns an error
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Internal server error",
			})
		}))
		defer errorServer.Close()

		cfg := &config.Config{
			Instance: errorServer.URL,
			Timeout:  30,
		}
		client := NewClient(cfg)

		req := NewSearchRequest("test")
		req.Format = "json"
		req.Page = 1

		_, err := client.Search(req)
		if err == nil {
			t.Error("Expected error from error server")
		}
	})

	t.Run("MockServerInvalidJSON", func(t *testing.T) {
		// Create a server that returns invalid JSON
		badServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer badServer.Close()

		cfg := &config.Config{
			Instance: badServer.URL,
			Timeout:  30,
		}
		client := NewClient(cfg)

		req := NewSearchRequest("test")
		req.Format = "json"
		req.Page = 1

		_, err := client.Search(req)
		if err == nil {
			t.Error("Expected error from invalid JSON response")
		}
	})
}
