// Package cache provides integration tests for caching functionality.
package cache

import (
	"testing"
	"time"

	"github.com/mule-ai/search/internal/config"
	"github.com/mule-ai/search/internal/searxng"
)

// TestCachedClientCreation tests creating a cached client with various configs.
func TestCachedClientCreation(t *testing.T) {
	tests := []struct {
		name      string
		cacheSize int
		ttl       time.Duration
	}{
		{"small cache", 10, 5 * time.Minute},
		{"medium cache", 100, 10 * time.Minute},
		{"large cache", 1000, 30 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Instance: "https://search.butler.ooo",
				Timeout:  30,
			}
			client := searxng.NewClient(cfg)
			cached := NewCachedClient(client, tt.cacheSize, tt.ttl)

			stats := cached.GetStats()
			if stats.MaxSize != tt.cacheSize {
				t.Errorf("expected MaxSize %d, got %d", tt.cacheSize, stats.MaxSize)
			}
		})
	}
}

// TestCachedClientCacheAccess tests accessing the underlying cache.
func TestCachedClientCacheAccess(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 10, 5*time.Minute)

	// Get the underlying cache
	cache := cached.GetCache()
	if cache == nil {
		t.Fatal("GetCache returned nil")
	}

	// Verify we can use it directly
	testData := &searxng.SearchResponse{
		Query: "test",
		Results: []searxng.SearchResult{
			{Title: "Test Result", URL: "https://example.com"},
		},
	}
	cache.Set("test-key", testData)

	result, found := cache.Get("test-key")
	if !found {
		t.Fatal("expected to find cached data")
	}

	resp := result.(*searxng.SearchResponse)
	if resp.Query != "test" {
		t.Errorf("expected query 'test', got '%s'", resp.Query)
	}
	if len(resp.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.Results))
	}
}

// TestCachedClientClearCacheIntegration tests the ClearCache method in an integration context.
func TestCachedClientClearCacheIntegration(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 10, 5*time.Minute)

	// Add some test data
	cached.GetCache().Set("key1", &searxng.SearchResponse{Query: "query1"})
	cached.GetCache().Set("key2", &searxng.SearchResponse{Query: "query2"})

	stats := cached.GetStats()
	if stats.Size != 2 {
		t.Errorf("expected cache size 2, got %d", stats.Size)
	}

	// Clear cache
	cached.ClearCache()

	stats = cached.GetStats()
	if stats.Size != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", stats.Size)
	}
}

// TestCachedClientGetStatsIntegration tests the GetStats method in an integration context.
func TestCachedClientGetStatsIntegration(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 100, 5*time.Minute)

	// Initially empty
	stats := cached.GetStats()
	if stats.Size != 0 {
		t.Errorf("expected Size 0, got %d", stats.Size)
	}
	if stats.MaxSize != 100 {
		t.Errorf("expected MaxSize 100, got %d", stats.MaxSize)
	}

	// Add some data
	cached.GetCache().Set("key1", &searxng.SearchResponse{Query: "query1"})
	cached.GetCache().Set("key2", &searxng.SearchResponse{Query: "query2"})

	stats = cached.GetStats()
	if stats.Size != 2 {
		t.Errorf("expected Size 2, got %d", stats.Size)
	}
}

// TestCachedClientGetClient tests accessing the underlying SearXNG client.
func TestCachedClientGetClient(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 10, 5*time.Minute)

	underlying := cached.GetClient()
	if underlying == nil {
		t.Fatal("GetClient returned nil")
	}
	if underlying != client {
		t.Error("GetClient didn't return the original client")
	}
}

// TestCacheKeyDifferentRequests tests that different requests produce different keys.
func TestCacheKeyDifferentRequests(t *testing.T) {
	requests := []*searxng.SearchRequest{
		{Query: "golang", Page: 1, Format: "json"},
		{Query: "golang", Page: 2, Format: "json"},
		{Query: "rust", Page: 1, Format: "json"},
		{Query: "golang", Page: 1, Format: "rss"},
		{Query: "golang", Page: 1, Format: "json", Categories: []string{"images"}},
	}

	keys := make(map[string]bool)
	for _, req := range requests {
		key := cacheKey(req)
		if keys[key] {
			t.Errorf("duplicate key generated for different request: %s", key)
		}
		keys[key] = true
	}

	if len(keys) != len(requests) {
		t.Errorf("expected %d unique keys, got %d", len(requests), len(keys))
	}
}

// TestCacheWithRealConfig tests cache configuration loading.
func TestCacheWithRealConfig(t *testing.T) {
	// Create a config with cache settings
	cfg := &config.Config{
		Instance:     "https://search.butler.ooo",
		Timeout:      30,
		CacheEnabled: true,
		CacheSize:    50,
		CacheTTL:     600, // 10 minutes
	}

	// Create a client
	client := searxng.NewClient(cfg)

	// Wrap with caching using config values
	cached := NewCachedClient(client, cfg.CacheSize, time.Duration(cfg.CacheTTL)*time.Second)

	// Verify cache is configured correctly
	stats := cached.GetStats()
	if stats.MaxSize != cfg.CacheSize {
		t.Errorf("expected MaxSize %d, got %d", cfg.CacheSize, stats.MaxSize)
	}

	// Verify we can add entries
	testData := &searxng.SearchResponse{
		Query: "test",
		Results: []searxng.SearchResult{
			{Title: "Test Result", URL: "https://example.com"},
		},
	}
	cached.GetCache().Set("test", testData)

	if cached.GetCache().Size() != 1 {
		t.Errorf("expected cache size 1, got %d", cached.GetCache().Size())
	}
}

// TestCacheKeyConsistency tests that cache keys are consistent across calls.
func TestCacheKeyConsistency(t *testing.T) {
	req := &searxng.SearchRequest{
		Query:      "test query",
		Page:       1,
		Format:     "json",
		Categories: []string{"general"},
		Languages:  []string{"en"},
		SafeSearch: 1,
		TimeRange:  "",
	}

	// Generate keys multiple times
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = cacheKey(req)
	}

	// All keys should be the same
	for i := 1; i < len(keys); i++ {
		if keys[i] != keys[0] {
			t.Errorf("key[%d] != key[0]: %s != %s", i, keys[i], keys[0])
		}
	}
}

// TestCacheKeyHashLength tests that cache keys are the expected length.
func TestCacheKeyHashLength(t *testing.T) {
	req := &searxng.SearchRequest{
		Query:  "test",
		Page:   1,
		Format: "json",
	}

	key := cacheKey(req)

	// cacheKey returns first 16 chars of SHA256 hash
	if len(key) != 16 {
		t.Errorf("expected key length 16, got %d", len(key))
	}

	// Should be valid hex
	for _, c := range key {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("key contains invalid hex character: %c", c)
		}
	}
}

// BenchmarkCacheKey benchmarks cache key generation.
func BenchmarkCacheKey(b *testing.B) {
	req := &searxng.SearchRequest{
		Query:      "benchmark test query",
		Page:       1,
		Format:     "json",
		Categories: []string{"general"},
		Languages:  []string{"en"},
		SafeSearch: 1,
		TimeRange:  "week",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cacheKey(req)
	}
}