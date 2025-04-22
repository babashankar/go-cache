package gocache

import (
	"testing"
	"time"
)

type testStruct struct {
	Name string
	Age  int
}

func TestCacheBasic(t *testing.T) {
	c := New(time.Minute)
	defer c.StopJanitor()

	// Test string
	err := c.Set("test", "value")
	if err != nil {
		t.Fatalf("Error setting value: %v", err)
	}

	val, found := c.GetString("test")
	if !found {
		t.Fatal("Expected to find key 'test'")
	}
	if val != "value" {
		t.Fatalf("Expected 'value', got '%s'", val)
	}

	// Test struct
	item := testStruct{Name: "John", Age: 30}
	err = c.Set("item", item)
	if err != nil {
		t.Fatalf("Error setting struct: %v", err)
	}

	var retrievedItem testStruct
	found, err = c.Get("item", &retrievedItem)
	if !found {
		t.Fatal("Expected to find key 'item'")
	}
	if err != nil {
		t.Fatalf("Error getting struct: %v", err)
	}
	if retrievedItem.Name != "John" || retrievedItem.Age != 30 {
		t.Fatalf("Retrieved item doesn't match: %+v", retrievedItem)
	}
}

func TestCacheExpiration(t *testing.T) {
	c := New(100 * time.Millisecond) // Fast cleanup for tests
	defer c.StopJanitor()

	// Set item with 300ms expiration
	c.SetWithExpiration("expire", "soon", 300*time.Millisecond)

	// Should be available immediately
	val, found := c.GetString("expire")
	if !found {
		t.Fatal("Item should be available")
	}
	if val != "soon" {
		t.Fatalf("Expected 'soon', got '%s'", val)
	}

	// Wait for expiration
	time.Sleep(400 * time.Millisecond)

	// Should be gone now
	_, found = c.GetString("expire")
	if found {
		t.Fatal("Item should have expired")
	}

	// Test TTL
	c.SetWithExpiration("ttl-test", "value", 5*time.Second)
	ttl, err := c.TTL("ttl-test")
	if err != nil {
		t.Fatalf("Error getting TTL: %v", err)
	}
	if ttl <= 0 || ttl > 5*time.Second {
		t.Fatalf("TTL should be between 0 and 5 seconds, got %v", ttl)
	}

	// Test no expiration
	c.Set("forever", "value")
	ttl, err = c.TTL("forever")
	if err != nil {
		t.Fatalf("Error getting TTL: %v", err)
	}
	if ttl != -1 {
		t.Fatalf("TTL should be -1 for non-expiring items, got %v", ttl)
	}
}

func TestCacheDelete(t *testing.T) {
	c := New(time.Minute)
	defer c.StopJanitor()

	c.Set("delete-me", "value")
	c.Delete("delete-me")

	_, found := c.GetString("delete-me")
	if found {
		t.Fatal("Item should be deleted")
	}

	// Test flush
	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Flush()

	if c.Count() != 0 {
		t.Fatalf("Cache should be empty after flush, has %d items", c.Count())
	}
}

func TestRawBytes(t *testing.T) {
	c := New(time.Minute)
	defer c.StopJanitor()

	// Test storing and retrieving raw bytes
	rawData := []byte{0x01, 0x02, 0x03, 0x04}
	c.Set("binary", rawData)

	retrieved, found := c.GetBytes("binary")
	if !found {
		t.Fatal("Expected to find binary data")
	}

	if len(retrieved) != len(rawData) {
		t.Fatalf("Binary data length mismatch: expected %d, got %d", len(rawData), len(retrieved))
	}

	for i, b := range rawData {
		if retrieved[i] != b {
			t.Fatalf("Binary data mismatch at position %d: expected %d, got %d", i, b, retrieved[i])
		}
	}
}
