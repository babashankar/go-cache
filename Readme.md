# Go-Cache

A high-performance, thread-safe in-memory caching solution for Go applications.

## Features

- 🔄 Thread-safe for concurrent access
- 🕒 Optional expiration times for cached items
- 🧹 Automatic cleanup of expired items
- 💾 Efficient storage of any data type
- 🔍 Simple and intuitive API
- 🚀 High performance with minimal overhead

## Installation

```bash
go get github.com/yourusername/go-cache
```

## Basic Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/yourusername/go-cache"
)

type User struct {
	Name  string
	Email string
}

func main() {
	// Create a new cache with cleanup every 5 minutes
	c := gocache.New(5 * time.Minute)
	defer c.StopJanitor() // Important: stop the cleanup goroutine when done

	// Store a string (no expiration)
	c.Set("key", "value")

	// Store a value with expiration
	c.SetWithExpiration("tempKey", "I'll expire", 30*time.Second)

	// Retrieve a string value
	value, found := c.GetString("key")
	if found {
		fmt.Println(value) // Outputs: value
	}

	// Store a struct
	user := User{Name: "John", Email: "john@example.com"}
	c.Set("user:123", user)

	// Retrieve and unmarshal into a struct
	var retrievedUser User
	found, _ := c.Get("user:123", &retrievedUser)
	if found {
		fmt.Printf("%+v\n", retrievedUser) // Outputs: {Name:John Email:john@example.com}
	}

	// Check time-to-live (TTL)
	ttl, _ := c.TTL("tempKey")
	fmt.Printf("Expires in: %v\n", ttl) // Something less than 30 seconds

	// Check if key exists
	if c.Exists("key") {
		fmt.Println("Key exists!")
	}

	// Delete a key
	c.Delete("key")

	// Flush all keys
	c.Flush()
}
```

## API Reference

### Creating a New Cache

```go
// Create a new cache with cleanup every 5 minutes
cache := gocache.New(5 * time.Minute)

// Create a cache with no automatic cleanup
cache := gocache.New(0)
```

### Setting Values

```go
// Set with no expiration
cache.Set("key", value)

// Set with expiration
cache.SetWithExpiration("key", value, 30 * time.Second)
```

### Getting Values

```go
// Get bytes
data, found := cache.GetBytes("key")

// Get string
str, found := cache.GetString("key")

// Get struct or any other type
var user User
found, err := cache.Get("user:123", &user)
```

### Other Operations

```go
// Delete a key
cache.Delete("key")

// Check if a key exists
exists := cache.Exists("key")

// Get time-to-live for a key
ttl, err := cache.TTL("key")

// Count items in cache
count := cache.Count()

// Remove all expired items manually
cache.DeleteExpired()

// Remove all items
cache.Flush()

// Stop the cleanup goroutine (important!)
cache.StopJanitor()
```