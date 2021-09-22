package rate

import (
	"context"
	"sync"
	"time"

	"github.com/DisgoOrg/log"
)

var _ Limiter = (*limiterImpl)(nil)

func NewLimiter(config *Config) Limiter {
	if config == nil {
		config = &DefaultConfig
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	return &limiterImpl{
		buckets: map[int]*bucket{},
		config:  *config,
	}
}

type limiterImpl struct {
	sync.Mutex

	buckets map[int]*bucket
	config  Config
}

func (r *limiterImpl) Logger() log.Logger {
	return r.config.Logger
}

func (r *limiterImpl) Close(ctx context.Context) {

}

func (r *limiterImpl) Config() Config {
	return r.config
}

func (r *limiterImpl) getBucket(shardID int, create bool) *bucket {
	r.Logger().Debug("locking shard rate limiter")
	r.Lock()
	defer func() {
		r.Logger().Debug("unlocking shard rate limiter")
		r.Unlock()
	}()
	key := ShardMaxConcurrencyKey(shardID, r.config.MaxConcurrency)
	b, ok := r.buckets[key]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{}
		r.buckets[key] = b
	}
	return b
}

func (r *limiterImpl) WaitBucket(ctx context.Context, shardID int) error {
	b := r.getBucket(shardID, true)
	r.Logger().Debugf("locking bucket: %+v", b)
	b.Lock()

	var until time.Time
	now := time.Now()

	if b.Reset.After(now) {
		until = b.Reset
	}

	if until.After(now) {
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return ErrCtxTimeout
		}

		select {
		case <-ctx.Done():
			b.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (r *limiterImpl) UnlockBucket(shardID int) {
	b := r.getBucket(shardID, false)
	if b == nil {
		return
	}
	defer func() {
		r.Logger().Debugf("unlocking bucket: %+v", b)
		b.Unlock()
	}()

	b.Reset = time.Now().Add(5 * time.Second)
}

type bucket struct {
	sync.Mutex
	Reset time.Time
}
