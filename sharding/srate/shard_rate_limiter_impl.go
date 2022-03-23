package srate

import (
	"context"
	"sync"
	"time"

	"github.com/DisgoOrg/log"
	"github.com/sasha-s/go-csync"
)

var _ Limiter = (*limiterImpl)(nil)

func NewLimiter(opts ...ConfigOpt) Limiter {
	config := DefaultConfig()
	config.Apply(opts)

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
	var wg sync.WaitGroup
	r.Lock()

	for key := range r.buckets {
		wg.Add(1)
		b := r.buckets[key]
		go func() {
			defer wg.Done()
			if err := b.CLock(ctx); err != nil {
				r.Logger().Error("failed to close bucket: ", err)
			}
			b.Unlock()
		}()
	}
}

func (r *limiterImpl) Config() Config {
	return r.config
}

func (r *limiterImpl) getBucket(shardID int, create bool) *bucket {
	r.Logger().Debug("locking shard srate limiter")
	r.Lock()
	defer func() {
		r.Logger().Debug("unlocking shard srate limiter")
		r.Unlock()
	}()
	key := ShardMaxConcurrencyKey(shardID, r.config.MaxConcurrency)
	b, ok := r.buckets[key]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{
			Key: key,
		}
		r.buckets[key] = b
	}
	return b
}

func (r *limiterImpl) WaitBucket(ctx context.Context, shardID int) error {
	b := r.getBucket(shardID, true)
	r.Logger().Debugf("locking shard bucket: %+v", b)
	if err := b.CLock(ctx); err != nil {
		return err
	}

	var until time.Time
	now := time.Now()

	if b.Reset.After(now) {
		until = b.Reset
	}

	if until.After(now) {
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return context.DeadlineExceeded
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
		r.Logger().Debugf("unlocking shard bucket: %+v", b)
		b.Unlock()
	}()

	b.Reset = time.Now().Add(time.Duration(r.config.StartupDelay) * time.Second)
}

type bucket struct {
	csync.Mutex
	Key   int
	Reset time.Time
}
