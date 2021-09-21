package sharding

import (
	"context"

	"github.com/DisgoOrg/log"
	"github.com/pkg/errors"
)

var ErrCtxTimeout = errors.New("rate limit exceeds context deadline")

type ShardRateLimiter interface {
	Logger() log.Logger
	Close(ctx context.Context)
	Config() Config
	WaitBucket(ctx context.Context, shardID int) error
	UnlockBucket(shardID int) error
}

func ShardRateLimitKey(shardID int, maxConcurrency int) int {
	return shardID % maxConcurrency
}
