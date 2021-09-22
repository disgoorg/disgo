package rate

import (
	"context"

	"github.com/DisgoOrg/log"
	"github.com/pkg/errors"
)

var ErrCtxTimeout = errors.New("rate limit exceeds context deadline")

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
