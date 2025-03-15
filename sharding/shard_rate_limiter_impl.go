package sharding

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/sasha-s/go-csync"
)

var _ RateLimiter = (*rateLimiterImpl)(nil)

// NewRateLimiter creates a new default RateLimiter with the given RateLimiterConfigOpt(s).
func NewRateLimiter(opts ...RateLimiterConfigOpt) RateLimiter {
	config := DefaultRateLimiterConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "sharding_rate_limiter"))

	return &rateLimiterImpl{
		buckets: make(map[int]*bucket),
		config:  *config,
	}
}

type rateLimiterImpl struct {
	mu sync.Mutex

	buckets map[int]*bucket
	config  RateLimiterConfig
}

func (r *rateLimiterImpl) Close(ctx context.Context) {
	r.config.Logger.Debug("closing shard rate limiter")
	var wg sync.WaitGroup
	r.mu.Lock()

	for key := range r.buckets {
		wg.Add(1)
		b := r.buckets[key]
		go func() {
			defer wg.Done()
			if err := b.mu.CLock(ctx); err != nil {
				r.config.Logger.Error("failed to close bucket", slog.Any("err", err))
			}
			b.mu.Unlock()
		}()
	}
	wg.Wait()
}

func (r *rateLimiterImpl) getBucket(shardID int, create bool) *bucket {
	r.config.Logger.Debug("locking shard rate limiter")
	r.mu.Lock()
	defer func() {
		r.config.Logger.Debug("unlocking shard rate limiter")
		r.mu.Unlock()
	}()

	key := ShardMaxConcurrencyKey(shardID, r.config.MaxConcurrency)
	b, ok := r.buckets[key]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{key: key}
		r.buckets[key] = b
	}
	return b
}

func (r *rateLimiterImpl) WaitBucket(ctx context.Context, shardID int) error {
	b := r.getBucket(shardID, true)
	r.config.Logger.Debug("locking shard bucket", slog.Int("key", b.key))
	if err := b.mu.CLock(ctx); err != nil {
		return err
	}

	now := time.Now()
	if b.reset.Before(now) {
		return nil
	}

	if deadline, ok := ctx.Deadline(); ok && b.reset.After(deadline) {
		b.mu.Unlock()
		return context.DeadlineExceeded
	}

	select {
	case <-ctx.Done():
		b.mu.Unlock()
		return ctx.Err()
	case <-time.After(b.reset.Sub(now)):
		return nil
	}
}

func (r *rateLimiterImpl) UnlockBucket(shardID int) {
	b := r.getBucket(shardID, false)
	if b == nil {
		return
	}

	defer func() {
		r.config.Logger.Debug("unlocking shard bucket", slog.Int("key", b.key), slog.Time("reset", b.reset))
		b.mu.Unlock()
	}()
	b.reset = time.Now().Add(r.config.IdentifyWait)
}

// bucket represents a rate-limiting bucket for a shard group.
type bucket struct {
	mu    csync.Mutex
	key   int
	reset time.Time
}
