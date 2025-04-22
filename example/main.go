// File: example/main.go
package main

import (
	"fmt"
	"time"

	gocache "git.source.akamai.com/~bsn/mock-bootstrapper.git/github/go-cache"
	// Replace with your actual import path
)

// Example user struct
type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	// Create a new cache with cleanup every minute
	c := gocache.New(1 * time.Minute)
	defer c.StopJanitor() // Important: stop the cleanup goroutine when done

	// Example 1: Store a simple string
	c.Set("greeting", "Hello, world!")
	greeting, found := c.GetString("greeting")
	if found {
		fmt.Println("Greeting:", greeting)
	}

	// Example 2: Store with expiration
	c.SetWithExpiration("temporary", "I'll be gone soon", 5*time.Second)

	// Check TTL
	ttl, _ := c.TTL("temporary")
	fmt.Printf("'temporary' expires in %v\n", ttl)

	// Wait a bit and check again
	time.Sleep(2 * time.Second)
	ttl, _ = c.TTL("temporary")
	fmt.Printf("'temporary' now expires in %v\n", ttl)

	// Example 3: Store and retrieve a struct
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	c.Set("user:1", user)

	var retrievedUser User
	found, err := c.Get("user:1", &retrievedUser)
	if found && err == nil {
		fmt.Printf("Retrieved user: %+v\n", retrievedUser)
	}

	// Example 4: Store raw bytes
	rawData := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" in ASCII
	c.Set("raw", rawData)

	data, found := c.GetBytes("raw")
	if found {
		fmt.Printf("Raw data: %v\n", data)
		fmt.Printf("As string: %s\n", string(data))
	}

	// Example 5: Check if an item exists
	if c.Exists("greeting") {
		fmt.Println("'greeting' exists in cache")
	}

	// Example 6: Delete an item
	c.Delete("greeting")
	if !c.Exists("greeting") {
		fmt.Println("'greeting' was successfully deleted")
	}

	// Example 7: Count items in cache
	fmt.Printf("Cache contains %d items\n", c.Count())

	// Example 8: Flush the cache
	c.Flush()
	fmt.Printf("After flush, cache contains %d items\n", c.Count())
}
