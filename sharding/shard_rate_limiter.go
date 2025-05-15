package sharding

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/sasha-s/go-csync"
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
	// If WaitBucket fails, UnlockBucket should not be called.
	UnlockBucket(shardID int)
}

// ShardMaxConcurrencyKey returns the bucket the given shardID with maxConcurrency belongs to.
func ShardMaxConcurrencyKey(shardID int, maxConcurrency int) int {
	return shardID % maxConcurrency
}

var _ RateLimiter = (*rateLimiterImpl)(nil)

// NewRateLimiter creates a new default RateLimiter with the given RateLimiterConfigOpt(s).
func NewRateLimiter(opts ...RateLimiterConfigOpt) RateLimiter {
	cfg := defaultRateLimiterConfig()
	cfg.apply(opts)

	return &rateLimiterImpl{
		buckets: make(map[int]*bucket),
		config:  cfg,
	}
}

type rateLimiterImpl struct {
	mu sync.Mutex

	buckets map[int]*bucket
	config  rateLimiterConfig
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
