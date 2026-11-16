package gateway

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sasha-s/go-csync"
)

// ErrRateLimiterClosed is returned by RateLimiter.Wait after RateLimiter.Close has been called.
var ErrRateLimiterClosed = errors.New("rate limiter closed")

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
	// Close marks the RateLimiter closed and waits for a send already admitted by Wait to finish.
	// Wait calls that have not been admitted fail with ErrRateLimiterClosed.
	// If the context deadline is exceeded, Close returns immediately; new Wait calls are still rejected.
	Close(ctx context.Context)

	// Wait waits for the RateLimiter to be ready to send a new message.
	// If the context deadline is exceeded, Wait will return immediately and no message will be sent.
	// After Close, Wait returns ErrRateLimiterClosed.
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
	mu sync.Mutex

	// rateMu is held from a successful Wait until Unlock.
	rateMu    csync.Mutex
	reset     time.Time
	remaining int

	config rateLimiterConfig

	// closed rejects Wait calls once Close has started.
	// pending counts Wait calls admitted before closed was set.
	closed  bool
	pending sync.WaitGroup
}

func (l *rateLimiterImpl) Close(ctx context.Context) {
	l.mu.Lock()
	l.closed = true
	l.mu.Unlock()

	// An admitted send still holds rateMu until Unlock, so wait without taking it.
	done := make(chan struct{})
	go func() {
		l.pending.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
	}
}

// admit reserves a slot on l.pending. The caller must Done it, either from Unlock or from a Wait path that returns before Unlock.
func (l *rateLimiterImpl) admit() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return ErrRateLimiterClosed
	}
	l.pending.Add(1)
	return nil
}

// Note: this function updates internal state and must be called from a lock state
func (l *rateLimiterImpl) isRateLimited(now time.Time, commandType RateLimiterCommandType) bool {
	if now.After(l.reset) {
		l.reset = now.Add(time.Minute)
		l.remaining = l.config.CommandsPerMinute
	}

	return l.remaining <= 0 || (l.remaining < l.config.ReservedCommandSlots && commandType != InternalCommandType)
}

func (l *rateLimiterImpl) Wait(ctx context.Context, commandType RateLimiterCommandType) (err error) {
	if err = l.admit(); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			l.pending.Done()
		}
	}()

	l.config.Logger.Debug("locking gateway rate limiter")

	for {
		if err = l.rateMu.CLock(ctx); err != nil {
			return err
		}

		now := time.Now()
		if l.isRateLimited(now, commandType) {
			l.rateMu.Unlock()
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
	defer l.pending.Done()
	l.config.Logger.Debug("unlocking gateway rate limiter")
	l.rateMu.Unlock()
}
