package sharding

import (
	"context"
)

// MaxConcurrency is the default number of shards that can log in at the same time.
const MaxConcurrency = 1

// RateLimiter limits how many shards can log in to Discord at the same time.
type RateLimiter interface {
	// Close gracefully closes the RateLimiter.
	// If the context deadline is exceeded, the RateLimiter will be closed immediately.
	Close(ctx context.Context)

	// WaitBucket waits for the given shardID bucket to be available for new logins.
	// If the context deadline is exceeded, WaitBucket will return immediately and no login will be attempted.
	WaitBucket(ctx context.Context, shardID int) error

	// UnlockBucket unlocks the given shardID bucket.
	UnlockBucket(shardID int)
}

// ShardMaxConcurrencyKey returns the bucket the given shardID with maxConcurrency belongs to.
func ShardMaxConcurrencyKey(shardID int, maxConcurrency int) int {
	return shardID % maxConcurrency
}
