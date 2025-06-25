package gateway

import (
	"context"
	"time"

	"github.com/sasha-s/go-csync"
)

// CommandsPerMinute is the default number of commands per minute that the Gateway will allow.
const CommandsPerMinute = 120

// RateLimiter provides handles the rate limiting logic for connecting to Discord's Gateway.
type RateLimiter interface {
	// Close gracefully closes the RateLimiter.
	// If the context deadline is exceeded, the RateLimiter will be closed immediately.
	Close(ctx context.Context)

	// Reset resets the RateLimiter to its initial state.
	Reset()

	// Wait waits for the RateLimiter to be ready to send a new message.
	// If the context deadline is exceeded, Wait will return immediately and no message will be sent.
	Wait(ctx context.Context) error

	// Unlock unlocks the RateLimiter and allows the next message to be sent.
	Unlock()
}

var _ RateLimiter = (*rateLimiterImpl)(nil)

// NewRateLimiter creates a new default RateLimiter with the given RateLimiterConfigOpt(s).
func NewRateLimiter(opts ...RateLimiterConfigOpt) RateLimiter {
	cfg := defaultRateLimiterConfig()
	cfg.apply(opts)

	return &rateLimiterImpl{
		config: cfg,
	}
}

type rateLimiterImpl struct {
	mu csync.Mutex

	reset     time.Time
	remaining int

	config rateLimiterConfig
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
