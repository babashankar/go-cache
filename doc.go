// Package gocache provides a versatile, thread-safe in-memory caching mechanism for Go applications.
//
// The package is designed to be simple to use while providing features like:
// - Storage of any data type (serialized as bytes)
// - Optional expiration times for cached items
// - Automatic cleanup of expired items
// - Thread-safe operations for concurrent access
//
// Basic usage:
//
//	// Create a new cache with cleanup every 5 minutes
//	c := gocache.New(5 * time.Minute)
//
//	// Store a value indefinitely
//	c.Set("key", "value")
//
//	// Store a value with expiration
//	c.SetWithExpiration("tempKey", "tempValue", 30 * time.Second)
//
//	// Retrieve a string value
//	value, found := c.GetString("key")
//
//	// Retrieve and unmarshal into a struct
//	var user User
//	found, err := c.Get("user:123", &user)
//
// When done with the cache, you should call StopJanitor() to stop the cleanup goroutine.
package gocache
