package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-api-gateway/router-service/internal/domain/port"
	"github.com/redis/go-redis/v9"
)

// RedisCache implements the Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
}

// Config holds Redis connection configuration.
type Config struct {
	Address     string
	Password    string
	DB          int
	DefaultTTL  int // Default TTL in seconds
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(config Config) (port.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get retrieves a value from the cache by key.
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key not found, return empty string
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// Set stores a value in the cache with an optional TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl int) error {
	var err error
	if ttl > 0 {
		err = c.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
	} else {
		err = c.client.Set(ctx, key, value, 0).Err()
	}
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Delete removes a specific key from the cache.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// ClearPrefix removes all keys matching the given prefix.
func (c *RedisCache) ClearPrefix(ctx context.Context, prefix string) error {
	iter := c.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys with prefix %s: %w", prefix, err)
	}
	return nil
}

// Close closes the Redis connection.
func (c *RedisCache) Close() error {
	return c.client.Close()
}
