package rate

import (
	"context"

	"github.com/DisgoOrg/log"
)

type Limiter interface {
	Logger() log.Logger
	Close(ctx context.Context)
	Config() Config
	WaitBucket(ctx context.Context, shardID int) error
	UnlockBucket(shardID int)
}

func ShardMaxConcurrencyKey(shardID int, maxConcurrency int) int {
	return shardID % maxConcurrency
}
