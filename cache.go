package gocache

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Item represents a cache item with value and expiration
type Item struct {
	Value      []byte // Store all values as byte slices
	Expiration int64  // 0 means no expiration
	Created    int64
}

// Cache is a thread-safe in-memory key:value store with optional expiration
type Cache struct {
	items           map[string]Item
	mu              sync.RWMutex
	cleanupInterval time.Duration
	stopCleanup     chan bool
}

// New creates a new Cache with the provided cleanup interval
// cleanupInterval: 0 means no automatic cleanup
func New(cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		items:           make(map[string]Item),
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan bool),
	}

	// Start the janitor if cleanup interval > 0
	if cleanupInterval > 0 {
		go cache.startJanitor()
	}

	return cache
}

// Set adds an item to the cache with no expiration
func (c *Cache) Set(key string, value interface{}) error {
	return c.SetWithExpiration(key, value, 0) // 0 means no expiration
}

// SetWithExpiration adds an item to the cache with a specific expiration time
func (c *Cache) SetWithExpiration(key string, value interface{}, duration time.Duration) error {
	var bytes []byte
	var err error

	// Convert the value to []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		// Use JSON for everything else
		bytes, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	var expiration int64
	if duration <= 0 {
		// 0 or negative means no expiration
		expiration = 0
	} else {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = Item{
		Value:      bytes,
		Expiration: expiration,
		Created:    time.Now().UnixNano(),
	}
	c.mu.Unlock()

	return nil
}

// GetBytes retrieves raw byte data from the cache
func (c *Cache) GetBytes(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

// Get retrieves and unmarshals an item from the cache
func (c *Cache) Get(key string, target interface{}) (bool, error) {
	bytes, found := c.GetBytes(key)
	if !found {
		return false, nil
	}

	// If target is nil, just return found status
	if target == nil {
		return true, nil
	}

	// Handle string target specially for efficiency
	if strPtr, ok := target.(*string); ok {
		*strPtr = string(bytes)
		return true, nil
	}

	// Unmarshal for other types
	return true, json.Unmarshal(bytes, target)
}

// GetString gets a string value from the cache
func (c *Cache) GetString(key string) (string, bool) {
	bytes, found := c.GetBytes(key)
	if !found {
		return "", false
	}
	return string(bytes), true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Exists checks if a key exists in the cache and is not expired
func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return false
	}

	// Check if the item has expired
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return false
	}

	return true
}

// Flush removes all items from the cache
func (c *Cache) Flush() {
	c.mu.Lock()
	c.items = make(map[string]Item)
	c.mu.Unlock()
}

// Count returns the number of items in the cache (including expired items)
func (c *Cache) Count() int {
	c.mu.RLock()
	count := len(c.items)
	c.mu.RUnlock()
	return count
}

// DeleteExpired deletes all expired items from the cache
func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// startJanitor starts the cleanup goroutine
func (c *Cache) startJanitor() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// StopJanitor stops the cleanup goroutine
func (c *Cache) StopJanitor() {
	if c.cleanupInterval > 0 {
		c.stopCleanup <- true
	}
}

// TTL returns the time to live for a key
func (c *Cache) TTL(key string) (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return 0, errors.New("key not found")
	}

	if item.Expiration == 0 {
		return -1, nil // No expiration, return -1 to indicate infinite TTL
	}

	now := time.Now().UnixNano()
	if now > item.Expiration {
		return 0, errors.New("key expired")
	}

	return time.Duration(item.Expiration - now), nil
}
