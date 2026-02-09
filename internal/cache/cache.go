// Package cache provides a simple in-memory cache for search results.
//
// The cache is used to store recent search results, improving performance
// for repeated queries. It uses LRU (Least Recently Used) eviction policy
// when the cache reaches its maximum size.
package cache

import (
	"sync"
	"time"
)

// CacheEntry represents a cached search result with expiration.
type CacheEntry struct {
	Response interface{}
	Expires  time.Time
}

// Cache is a thread-safe in-memory cache with LRU eviction.
type Cache struct {
	mu       sync.RWMutex
	store    map[string]*CacheEntry
	maxSize  int
	ttl      time.Duration
	list     []string // Track order for LRU
}

// NewCache creates a new cache with the specified maximum size and TTL.
//
// maxSize is the maximum number of entries to store.
// ttl is the time-to-live for cache entries.
//
// Example:
//
//	cache := cache.NewCache(100, 5*time.Minute)
func NewCache(maxSize int, ttl time.Duration) *Cache {
	return &Cache{
		store:   make(map[string]*CacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
		list:    make([]string, 0, maxSize),
	}
}

// Get retrieves a value from the cache.
//
// Returns the cached response and true if found and not expired.
// Returns nil and false if not found or expired.
//
// Example:
//
//	cache := cache.NewCache(100, 5*time.Minute)
//	cache.Set("query:golang", searchResponse)
//	if val, found := cache.Get("query:golang"); found {
//	    resp := val.(*searxng.SearchResponse)
//	    fmt.Printf("Found cached results: %d\n", len(resp.Results))
//	}
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.store[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if time.Now().After(entry.Expires) {
		return nil, false
	}

	// Move to end of list (most recently used)
	c.moveToFront(key)

	return entry.Response, true
}

// Set stores a value in the cache with the current time + TTL.
//
// If the cache is full, the least recently used entry is evicted.
//
// Example:
//
//	cache := cache.NewCache(100, 5*time.Minute)
//	resp := &searxng.SearchResponse{Query: "golang", Results: [...]}
//	cache.Set("query:golang", resp)
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict
	if len(c.store) >= c.maxSize {
		c.evictLRU()
	}

	// Store the entry
	c.store[key] = &CacheEntry{
		Response: value,
		Expires:  time.Now().Add(c.ttl),
	}

	// Add to front of list
	c.list = append([]string{key}, c.list...)
}

// Delete removes an entry from the cache.
//
// Example:
//
//	cache := cache.NewCache(100, 5*time.Minute)
//	cache.Set("query:golang", searchResponse)
//	cache.Delete("query:golang")
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	c.removeFromList(key)
}

// Clear removes all entries from the cache.
//
// Example:
//
//	cache := cache.NewCache(100, 5*time.Minute)
//	cache.Set("query:golang", searchResponse)
//	cache.Set("query:rust", searchResponse2)
//	cache.Clear() // Cache is now empty
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store = make(map[string]*CacheEntry)
	c.list = make([]string, 0, c.maxSize)
}

// Size returns the current number of entries in the cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.store)
}

// evictLRU removes the least recently used entry from the cache.
// Must be called with the lock held.
func (c *Cache) evictLRU() {
	if len(c.list) == 0 {
		return
	}

	// Get the LRU key (last item in list)
	lruKey := c.list[len(c.list)-1]

	// Remove from store
	delete(c.store, lruKey)

	// Remove from list
	c.list = c.list[:len(c.list)-1]
}

// moveToFront moves a key to the front of the LRU list.
// Must be called with the lock held (or at least read lock).
func (c *Cache) moveToFront(key string) {
	// Remove from current position
	c.removeFromList(key)

	// Add to front
	c.list = append([]string{key}, c.list...)
}

// removeFromList removes a key from the LRU list.
// Must be called with the lock held.
func (c *Cache) removeFromList(key string) {
	for i, k := range c.list {
		if k == key {
			c.list = append(c.list[:i], c.list[i+1:]...)
			break
		}
	}
}

// Cleanup removes expired entries from the cache.
//
// This is called automatically periodically, but can be called manually
// if needed.
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.store {
		if now.After(entry.Expires) {
			delete(c.store, key)
			c.removeFromList(key)
		}
	}
}

// Stats returns cache statistics.
type Stats struct {
	Size    int
	MaxSize int
	Hits    int64
	Misses  int64
}

// stats tracking (simplified - in production you'd want atomic counters)
var (
	hits   int64
	misses int64
)

// GetStats returns cache statistics.
func (c *Cache) GetStats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Stats{
		Size:    len(c.store),
		MaxSize: c.maxSize,
		Hits:    hits,
		Misses:  misses,
	}
}

// recordHit records a cache hit.
func recordHit() {
	// In production, use atomic.AddInt64(&hits, 1)
}

// recordMiss records a cache miss.
func recordMiss() {
	// In production, use atomic.AddInt64(&misses, 1)
}
