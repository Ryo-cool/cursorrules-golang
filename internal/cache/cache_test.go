package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheOperations(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		ttl     time.Duration
		wantErr bool
	}{
		{
			name:    "basic set and get",
			key:     "test-key",
			value:   "test-value",
			ttl:     time.Minute,
			wantErr: false,
		},
		{
			name:    "expired key",
			key:     "expired-key",
			value:   "expired-value",
			ttl:     time.Millisecond,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(Config{
				MaxItems:        100,
				CleanupInterval: time.Minute,
			})

			// Test Set
			cache.Set(tt.key, tt.value, tt.ttl)

			// For expired key test, wait for expiration
			if tt.wantErr {
				time.Sleep(tt.ttl * 2)
			}

			// Test Get
			got, exists := cache.Get(tt.key)
			if exists == tt.wantErr {
				t.Errorf("Get() exists = %v, wantErr %v", exists, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.value {
				t.Errorf("Get() got = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestCacheEviction(t *testing.T) {
	maxItems := 5
	cache := New(Config{
		MaxItems:        maxItems,
		CleanupInterval: time.Minute,
	})

	// Test cache eviction policy
	for i := 0; i < maxItems+3; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, i, time.Minute)
	}

	// Verify cache size doesn't exceed max capacity
	stats := cache.GetStats()
	if stats["size"].(int) > maxItems {
		t.Errorf("Cache size %d exceeds max capacity %d", stats["size"], maxItems)
	}
}

func TestCacheStats(t *testing.T) {
	cache := New(Config{
		MaxItems:        100,
		CleanupInterval: time.Minute,
	})

	// Add some items and perform operations
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)

	// Get existing and non-existing items
	cache.Get("key1")
	cache.Get("key2")
	cache.Get("non-existent")

	// Check stats
	stats := cache.GetStats()
	if stats["hit_count"].(uint64) != 2 {
		t.Errorf("Expected hit count of 2, got %v", stats["hit_count"])
	}
	if stats["miss_count"].(uint64) != 1 {
		t.Errorf("Expected miss count of 1, got %v", stats["miss_count"])
	}
}

// ... additional test cases for cache optimization strategies ...
