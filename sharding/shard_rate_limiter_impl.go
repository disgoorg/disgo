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
		buckets: map[int]*bucket{},
		config:  *config,
	}
}

type rateLimiterImpl struct {
	mu sync.Mutex

	buckets map[int]*bucket
	config  RateLimiterConfig
}

func (r *rateLimiterImpl) Close(ctx context.Context) {
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

		b = &bucket{
			Key: key,
		}
		r.buckets[key] = b
	}
	return b
}

func (r *rateLimiterImpl) WaitBucket(ctx context.Context, shardID int) error {
	b := r.getBucket(shardID, true)
	r.config.Logger.Debug("locking shard bucket", slog.Int("key", b.Key), slog.Time("reset", b.Reset))
	if err := b.mu.CLock(ctx); err != nil {
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
			b.mu.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (r *rateLimiterImpl) UnlockBucket(shardID int) {
	b := r.getBucket(shardID, false)
	if b == nil {
		return
	}
	defer func() {
		r.config.Logger.Debug("unlocking shard bucket", slog.Int("key", b.Key), slog.Time("reset", b.Reset))
		b.mu.Unlock()
	}()

	b.Reset = time.Now().Add(5 * time.Second)
}

type bucket struct {
	mu    csync.Mutex
	Key   int
	Reset time.Time
}
