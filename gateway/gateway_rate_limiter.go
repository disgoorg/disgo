package gateway

import (
	"context"
	"time"

	"github.com/sasha-s/go-csync"
)

// CommandsPerMinute is the default number of commands per minute that the Gateway will allow.
const CommandsPerMinute = 120

// ReservedCommandSlots is the default number of CommandsPerMinute that the library
// will reserve for high priority events, like heartbeats
const ReservedCommandSlots = 3

// RateLimiterCommandType represents the type of wait performed by the rate-limiter
type RateLimiterCommandType int

const (
	NormalCommandType RateLimiterCommandType = iota
	InternalCommandType
)

// RateLimiter provides handles the rate limiting logic for connecting to Discord's Gateway.
type RateLimiter interface {
	// Close gracefully closes the RateLimiter.
	// If the context deadline is exceeded, the RateLimiter will be closed immediately.
	Close(ctx context.Context)

	// Reset resets the RateLimiter to its initial state.
	Reset()

	// Wait waits for the RateLimiter to be ready to send a new message.
	// If the context deadline is exceeded, Wait will return immediately and no message will be sent.
	Wait(ctx context.Context, commandType RateLimiterCommandType) error

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

// Note: this function updates internal state and must be called from a lock state
func (l *rateLimiterImpl) isRateLimited(now time.Time, commandType RateLimiterCommandType) bool {
	if now.After(l.reset) {
		l.reset = now.Add(time.Minute)
		l.remaining = l.config.CommandsPerMinute
	}

	return l.remaining <= 0 || (l.remaining < l.config.ReservedCommandSlots && commandType != InternalCommandType)
}

func (l *rateLimiterImpl) Wait(ctx context.Context, commandType RateLimiterCommandType) error {
	l.config.Logger.Debug("locking gateway rate limiter")

	for {
		if err := l.mu.CLock(ctx); err != nil {
			return err
		}

		now := time.Now()
		if l.isRateLimited(now, commandType) {
			l.mu.Unlock()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(l.reset.Sub(now)):
			}
			continue
		}
		break
	}

	l.remaining--
	return nil
}

func (l *rateLimiterImpl) Unlock() {
	l.config.Logger.Debug("unlocking gateway rate limiter")
	l.mu.Unlock()
}
