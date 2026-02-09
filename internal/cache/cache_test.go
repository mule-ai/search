// Package cache provides tests for the caching functionality.
package cache

import (
	"testing"
	"time"

	"github.com/mule-ai/search/internal/config"
	"github.com/mule-ai/search/internal/searxng"
)

// TestNewCache tests creating a new cache instance.
func TestNewCache(t *testing.T) {
	cache := NewCache(100, 5*time.Minute)
	if cache == nil {
		t.Fatal("NewCache returned nil")
	}
	if cache.maxSize != 100 {
		t.Errorf("expected maxSize 100, got %d", cache.maxSize)
	}
	if cache.ttl != 5*time.Minute {
		t.Errorf("expected ttl 5m, got %v", cache.ttl)
	}
}

// TestCacheSetGet tests basic set and get operations.
func TestCacheSetGet(t *testing.T) {
	cache := NewCache(10, 5*time.Minute)

	// Test set and get
	testData := &searxng.SearchResponse{
		Query:   "test query",
		Results: []searxng.SearchResult{},
	}
	cache.Set("test-key", testData)

	// Get should return the data
	result, found := cache.Get("test-key")
	if !found {
		t.Fatal("expected to find cached data")
	}
	if result != testData {
		t.Error("returned data doesn't match cached data")
	}
}

// TestCacheExpiration tests that entries expire after TTL.
func TestCacheExpiration(t *testing.T) {
	cache := NewCache(10, 10*time.Millisecond)

	testData := &searxng.SearchResponse{
		Query: "test query",
	}
	cache.Set("test-key", testData)

	// Should be found immediately
	_, found := cache.Get("test-key")
	if !found {
		t.Fatal("expected to find cached data immediately")
	}

	// Wait for expiration
	time.Sleep(15 * time.Millisecond)

	// Should not be found after expiration
	_, found = cache.Get("test-key")
	if found {
		t.Error("expected cache entry to be expired")
	}
}

// TestCacheLRU tests LRU eviction when cache is full.
func TestCacheLRU(t *testing.T) {
	cache := NewCache(3, 5*time.Minute)

	// Fill cache
	for i := 0; i < 3; i++ {
		cache.Set("key"+string(rune('0'+i)), &searxng.SearchResponse{
			Query: "query" + string(rune('0'+i)),
		})
	}

	// Access key0 to make it recently used
	cache.Get("key0")

	// Add one more entry, should evict key1 (least recently used)
	cache.Set("key3", &searxng.SearchResponse{
		Query: "query3",
	})

	// key0 should still exist
	_, found := cache.Get("key0")
	if !found {
		t.Error("expected key0 to still exist (was recently used)")
	}

	// key1 should be evicted
	_, found = cache.Get("key1")
	if found {
		t.Error("expected key1 to be evicted (least recently used)")
	}

	// key3 should exist
	_, found = cache.Get("key3")
	if !found {
		t.Error("expected key3 to exist (just added)")
	}
}

// TestCacheDelete tests deleting entries.
func TestCacheDelete(t *testing.T) {
	cache := NewCache(10, 5*time.Minute)

	testData := &searxng.SearchResponse{
		Query: "test query",
	}
	cache.Set("test-key", testData)

	// Verify it exists
	_, found := cache.Get("test-key")
	if !found {
		t.Fatal("expected to find cached data")
	}

	// Delete it
	cache.Delete("test-key")

	// Should not exist now
	_, found = cache.Get("test-key")
	if found {
		t.Error("expected cache entry to be deleted")
	}
}

// TestCacheClear tests clearing all entries.
func TestCacheClear(t *testing.T) {
	cache := NewCache(10, 5*time.Minute)

	// Add some entries
	for i := 0; i < 5; i++ {
		cache.Set("key"+string(rune('0'+i)), &searxng.SearchResponse{
			Query: "query" + string(rune('0'+i)),
		})
	}

	if cache.Size() != 5 {
		t.Errorf("expected cache size 5, got %d", cache.Size())
	}

	// Clear cache
	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", cache.Size())
	}

	// Verify no entries exist
	for i := 0; i < 5; i++ {
		_, found := cache.Get("key" + string(rune('0'+i)))
		if found {
			t.Errorf("expected key%d to be cleared", i)
		}
	}
}

// TestNewCachedClient tests creating a cached client.
func TestNewCachedClient(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 100, 5*time.Minute)

	if cached == nil {
		t.Fatal("NewCachedClient returned nil")
	}
	if cached.client != client {
		t.Error("cached client doesn't wrap the provided client")
	}
	if cached.cache == nil {
		t.Error("cached client doesn't have a cache")
	}
}

// TestCachedClientGetCache tests accessing the underlying cache.
func TestCachedClientGetCache(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 100, 5*time.Minute)

	cache := cached.GetCache()
	if cache == nil {
		t.Fatal("GetCache returned nil")
	}

	// Test that we can use the cache
	testData := &searxng.SearchResponse{Query: "test"}
	cache.Set("test", testData)

	result, found := cache.Get("test")
	if !found {
		t.Fatal("expected to find cached data")
	}
	if result != testData {
		t.Error("returned data doesn't match")
	}
}

// TestCachedClientClearCache tests clearing the cache.
func TestCachedClientClearCache(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 100, 5*time.Minute)

	// Add some data
	cached.GetCache().Set("test", &searxng.SearchResponse{Query: "test"})
	if cached.GetCache().Size() != 1 {
		t.Errorf("expected cache size 1, got %d", cached.GetCache().Size())
	}

	// Clear cache
	cached.ClearCache()
	if cached.GetCache().Size() != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", cached.GetCache().Size())
	}
}

// TestCachedClientGetStats tests getting cache statistics.
func TestCachedClientGetStats(t *testing.T) {
	cfg := &config.Config{
		Instance: "https://search.butler.ooo",
		Timeout:  30,
	}
	client := searxng.NewClient(cfg)
	cached := NewCachedClient(client, 100, 5*time.Minute)

	stats := cached.GetStats()
	if stats.MaxSize != 100 {
		t.Errorf("expected MaxSize 100, got %d", stats.MaxSize)
	}
	if stats.Size != 0 {
		t.Errorf("expected Size 0, got %d", stats.Size)
	}
}

// TestCacheKey tests cache key generation.
func TestCacheKey(t *testing.T) {
	req := &searxng.SearchRequest{
		Query:      "test query",
		Page:       2,
		Format:     "json",
		Categories: []string{"images"},
		Languages:  []string{"en"},
		SafeSearch: 1,
		TimeRange:  "week",
	}

	key1 := cacheKey(req)
	key2 := cacheKey(req)

	// Same request should produce same key
	if key1 != key2 {
		t.Errorf("same request produced different keys: %s != %s", key1, key2)
	}

	// Different request should produce different key
	req2 := &searxng.SearchRequest{
		Query:  "different query",
		Page:   1,
		Format: "json",
	}
	key3 := cacheKey(req2)
	if key1 == key3 {
		t.Error("different request produced same key")
	}
}

// TestCacheMoveToFront tests LRU list management.
func TestCacheMoveToFront(t *testing.T) {
	cache := NewCache(5, 5*time.Minute)

	// Add some entries
	cache.Set("key1", &searxng.SearchResponse{Query: "query1"})
	cache.Set("key2", &searxng.SearchResponse{Query: "query2"})
	cache.Set("key3", &searxng.SearchResponse{Query: "query3"})

	// Access key1 to move it to front
	cache.Get("key1")

	// Add more entries to fill cache and trigger eviction
	// Order before: [key1, key3, key2] (after accessing key1)
	// Adding key4, key5, key6 will evict key2 (least recently used)
	cache.Set("key4", &searxng.SearchResponse{Query: "query4"})
	cache.Set("key5", &searxng.SearchResponse{Query: "query5"})
	cache.Set("key6", &searxng.SearchResponse{Query: "query6"})

	// key1 should still exist because it was accessed recently
	_, found := cache.Get("key1")
	if !found {
		t.Error("expected key1 to still exist")
	}

	// key2 should be evicted (least recently used)
	_, found = cache.Get("key2")
	if found {
		t.Error("expected key2 to be evicted")
	}
}

// TestCacheCleanup tests expired entry cleanup.
func TestCacheCleanup(t *testing.T) {
	cache := NewCache(10, 10*time.Millisecond)

	// Add entries
	cache.Set("key1", &searxng.SearchResponse{Query: "query1"})
	cache.Set("key2", &searxng.SearchResponse{Query: "query2"})

	if cache.Size() != 2 {
		t.Errorf("expected cache size 2, got %d", cache.Size())
	}

	// Wait for expiration
	time.Sleep(15 * time.Millisecond)

	// Cleanup should remove expired entries
	cache.Cleanup()

	if cache.Size() != 0 {
		t.Errorf("expected cache size 0 after cleanup, got %d", cache.Size())
	}
}