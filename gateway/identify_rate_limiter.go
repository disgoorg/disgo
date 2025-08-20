package gateway

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/sasha-s/go-csync"
)

// DefaultMaxConcurrency is the default number of shards that can log in at the same time.
const DefaultMaxConcurrency = 1

// IdentifyRateLimiter limits how many shards can log in to Discord at the same time.
type IdentifyRateLimiter interface {
	// Close gracefully closes the RateLimiter.
	// If the context deadline is exceeded, the RateLimiter will be closed immediately.
	Close(ctx context.Context)

	// Wait waits for the given shardID bucket to be available for new logins.
	// If the context deadline is exceeded, Wait will return immediately and no login will be attempted.
	Wait(ctx context.Context, shardID int) error

	// Unlock unlocks the given shardID bucket.
	// If Wait fails, Unlock should not be called.
	Unlock(shardID int)
}

// MaxConcurrencyKey returns the bucket the given shardID with maxConcurrency belongs to.
func MaxConcurrencyKey(shardID int, maxConcurrency int) int {
	return shardID % maxConcurrency
}

var _ IdentifyRateLimiter = (*identifyRateLimiterImpl)(nil)

// NewIdentifyRateLimiter creates a new default RateLimiter with the given IdentifyRateLimiterConfigOpt(s).
func NewIdentifyRateLimiter(opts ...IdentifyRateLimiterConfigOpt) IdentifyRateLimiter {
	cfg := defaultIdentifyRateLimiterConfig()
	cfg.apply(opts)

	return &identifyRateLimiterImpl{
		buckets: make(map[int]*identifyBucket),
		config:  cfg,
	}
}

type identifyRateLimiterImpl struct {
	mu sync.Mutex

	buckets map[int]*identifyBucket
	config  identifyRateLimiterConfig
}

func (r *identifyRateLimiterImpl) Close(ctx context.Context) {
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

func (r *identifyRateLimiterImpl) getBucket(shardID int, create bool) *identifyBucket {
	r.config.Logger.Debug("locking shard rate limiter")
	r.mu.Lock()
	defer func() {
		r.config.Logger.Debug("unlocking shard rate limiter")
		r.mu.Unlock()
	}()

	key := MaxConcurrencyKey(shardID, r.config.MaxConcurrency)
	b, ok := r.buckets[key]
	if !ok {
		if !create {
			return nil
		}

		b = &identifyBucket{key: key}
		r.buckets[key] = b
	}
	return b
}

func (r *identifyRateLimiterImpl) Wait(ctx context.Context, shardID int) error {
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

func (r *identifyRateLimiterImpl) Unlock(shardID int) {
	b := r.getBucket(shardID, false)
	if b == nil {
		return
	}

	defer func() {
		r.config.Logger.Debug("unlocking shard bucket", slog.Int("key", b.key), slog.Time("reset", b.reset))
		b.mu.Unlock()
	}()
	b.reset = time.Now().Add(r.config.Wait)
}

// identifyBucket represents a rate-limiting bucket for a shard group.
type identifyBucket struct {
	mu    csync.Mutex
	key   int
	reset time.Time
}
