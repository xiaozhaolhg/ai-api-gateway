package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestRedisCache tests the Redis cache implementation
// NOTE: These tests require a running Redis instance
// In a full test setup, you would use testcontainers or a mock
func TestRedisCache(t *testing.T) {
	// Skip if Redis is not available
	t.Skip("Redis cache tests require running Redis instance - skip for MVP")

	ctx := context.Background()

	// Create a mock Redis client for testing
	// In production, use testcontainers or a real Redis instance
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	cache := &RedisCache{client: client}
	defer cache.Close()

	t.Run("SetAndGet", func(t *testing.T) {
		key := "test:key"
		value := "test-value"

		err := cache.Set(ctx, key, value, 60)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		retrieved, err := cache.Get(ctx, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved != value {
			t.Errorf("Expected %s, got %s", value, retrieved)
		}
	})

	t.Run("GetNonExistent", func(t *testing.T) {
		_, err := cache.Get(ctx, "nonexistent:key")
		if err != nil {
			t.Errorf("Get should return empty string for non-existent key, got error: %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:delete"
		value := "value-to-delete"

		cache.Set(ctx, key, value, 60)
		cache.Delete(ctx, key)

		retrieved, _ := cache.Get(ctx, key)
		if retrieved != "" {
			t.Errorf("Expected empty string after delete, got %s", retrieved)
		}
	})

	t.Run("ClearPrefix", func(t *testing.T) {
		prefix := "test:clear:"
		keys := []string{"test:clear:1", "test:clear:2", "test:clear:3"}

		for _, key := range keys {
			cache.Set(ctx, key, "value", 60)
		}

		err := cache.ClearPrefix(ctx, prefix)
		if err != nil {
			t.Fatalf("ClearPrefix failed: %v", err)
		}

		for _, key := range keys {
			retrieved, _ := cache.Get(ctx, key)
			if retrieved != "" {
				t.Errorf("Expected empty string after ClearPrefix for key %s", key)
			}
		}
	})

	t.Run("TTLExpiration", func(t *testing.T) {
		key := "test:ttl"
		value := "value-with-ttl"

		// Set with 1 second TTL
		cache.Set(ctx, key, value, 1)

		// Should be available immediately
		retrieved, _ := cache.Get(ctx, key)
		if retrieved != value {
			t.Errorf("Expected value before TTL expiration")
		}

		// Wait for expiration
		time.Sleep(2 * time.Second)

		// Should be expired
		retrieved, _ = cache.Get(ctx, key)
		if retrieved != "" {
			t.Errorf("Expected empty string after TTL expiration")
		}
	})
}
