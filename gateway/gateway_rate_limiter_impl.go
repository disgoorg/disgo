package gateway

import (
	"context"
	"log/slog"
	"time"

	"github.com/sasha-s/go-csync"
)

// NewRateLimiter creates a new default RateLimiter with the given RateLimiterConfigOpt(s).
func NewRateLimiter(opts ...RateLimiterConfigOpt) RateLimiter {
	config := DefaultRateLimiterConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "gateway_rate_limiter"))

	return &rateLimiterImpl{
		config: *config,
	}
}

type rateLimiterImpl struct {
	mu csync.Mutex

	reset     time.Time
	remaining int

	config RateLimiterConfig
}

func (l *rateLimiterImpl) Close(ctx context.Context) {
	_ = l.mu.CLock(ctx)
}

func (l *rateLimiterImpl) Reset() {
	l.reset = time.Time{}
	l.remaining = 0
	l.mu = csync.Mutex{}
}

func (l *rateLimiterImpl) Wait(ctx context.Context) error {
	l.config.Logger.Debug("locking gateway rate limiter")
	if err := l.mu.CLock(ctx); err != nil {
		return err
	}

	now := time.Now()

	var until time.Time

	if l.remaining == 0 && l.reset.After(now) {
		until = l.reset
	}

	if until.After(now) {
		select {
		case <-ctx.Done():
			l.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (l *rateLimiterImpl) Unlock() {
	l.config.Logger.Debug("unlocking gateway rate limiter")
	now := time.Now()
	if l.reset.Before(now) {
		l.reset = now.Add(time.Minute)
		l.remaining = l.config.CommandsPerMinute
	}
	l.mu.Unlock()
}
