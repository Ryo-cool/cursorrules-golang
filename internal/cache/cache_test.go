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
			cache := NewCache()

			// Test Set
			err := cache.Set(tt.key, tt.value, tt.ttl)
			if err != nil {
				t.Errorf("Set() error = %v", err)
				return
			}

			// For expired key test, wait for expiration
			if tt.wantErr {
				time.Sleep(tt.ttl * 2)
			}

			// Test Get
			got, err := cache.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.value {
				t.Errorf("Get() got = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestCacheEviction(t *testing.T) {
	cache := NewCache()
	maxItems := 1000

	// Test cache eviction policy
	for i := 0; i < maxItems+100; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := cache.Set(key, i, time.Minute)
		if err != nil {
			t.Errorf("Set() error = %v", err)
		}
	}

	// Verify cache size doesn't exceed max capacity
	if cache.Size() > maxItems {
		t.Errorf("Cache size %d exceeds max capacity %d", cache.Size(), maxItems)
	}
}

// ... additional test cases for cache optimization strategies ...
