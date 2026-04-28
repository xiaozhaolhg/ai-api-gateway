package port

import "context"

// Cache defines the interface for caching routing data.
// This follows Clean Architecture principles - the domain defines the contract,
// and infrastructure provides the implementation (Redis).
type Cache interface {
	// Get retrieves a value from the cache by key.
	// Returns the value and a boolean indicating if the key was found.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value in the cache with an optional TTL (time-to-live) in seconds.
	// If ttl is 0, the key has no expiration.
	Set(ctx context.Context, key string, value string, ttl int) error

	// Delete removes a specific key from the cache.
	Delete(ctx context.Context, key string) error

	// ClearPrefix removes all keys matching the given prefix.
	// This is used for cache invalidation (e.g., clearing all router:* keys).
	ClearPrefix(ctx context.Context, prefix string) error

	// Close closes the cache connection and releases resources.
	Close() error
}
