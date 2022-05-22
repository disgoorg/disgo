package srate

import (
	"context"

	"github.com/disgoorg/log"
)

// Limiter limits how many shards can log in to Discord at the same time.
type Limiter interface {
	// Logger returns the logger the Limiter uses
	Logger() log.Logger

	// Close gracefully closes the Limiter.
	// If the context deadline is exceeded, the Limiter will be closed immediately.
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
