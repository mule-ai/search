package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mule-ai/search/internal/searxng"
)

// CachedClient wraps a SearXNG client with caching functionality.
type CachedClient struct {
	client *searxng.Client
	cache  *Cache
}

// NewCachedClient creates a new cached SearXNG client.
//
// The cache will store up to maxCache entries for the specified TTL.
// Recommended values are maxCache=100 and ttl=5*time.Minute for most use cases.
//
// Example:
//
//	client := searxng.NewClient(cfg)
//	cached := cache.NewCachedClient(client, 100, 5*time.Minute)
//	resp, err := cached.Search(req)
func NewCachedClient(client *searxng.Client, maxCache int, ttl time.Duration) *CachedClient {
	return &CachedClient{
		client: client,
		cache:  NewCache(maxCache, ttl),
	}
}

// Search executes a search query, using the cache if available.
//
// The cache key is generated from the search request parameters.
// Cached results are returned immediately without an API call.
func (cc *CachedClient) Search(req *searxng.SearchRequest) (*searxng.SearchResponse, error) {
	// Generate cache key
	key := cacheKey(req)

	// Try to get from cache
	if cached, found := cc.cache.Get(key); found {
		recordHit()
		if resp, ok := cached.(*searxng.SearchResponse); ok {
			return resp, nil
		}
		// If type assertion fails, treat as cache miss
		recordMiss()
	} else {
		recordMiss()
	}

	// Execute search - bypass cache and call client directly
	resp, err := cc.client.Search(req)
	if err != nil {
		return nil, err
	}

	// Store in cache
	cc.cache.Set(key, resp)

	return resp, nil
}

// cacheKey generates a unique cache key from a search request.
func cacheKey(req *searxng.SearchRequest) string {
	// Create a hash of the request parameters
	h := sha256.New()
	h.Write([]byte(req.Query))
	h.Write([]byte(fmt.Sprintf("%d", req.Page)))
	h.Write([]byte(req.Format))

	for _, cat := range req.Categories {
		h.Write([]byte(cat))
	}

	for _, lang := range req.Languages {
		h.Write([]byte(lang))
	}

	h.Write([]byte(fmt.Sprintf("%d", req.SafeSearch)))
	h.Write([]byte(req.TimeRange))

	// Return hex string (first 16 chars is enough for uniqueness)
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// GetCache returns the underlying cache for direct access.
//
// This allows you to clear the cache, get stats, or perform other operations.
//
// Example:
//
//	cached.GetCache().Clear()
//	stats := cached.GetCache().GetStats()
func (cc *CachedClient) GetCache() *Cache {
	return cc.cache
}

// GetClient returns the underlying SearXNG client.
func (cc *CachedClient) GetClient() *searxng.Client {
	return cc.client
}

// ClearCache clears all cached entries.
func (cc *CachedClient) ClearCache() {
	cc.cache.Clear()
}

// GetStats returns cache statistics.
func (cc *CachedClient) GetStats() Stats {
	return cc.cache.GetStats()
}
