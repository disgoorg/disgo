package rate

import (
	"context"
	"time"

	"github.com/DisgoOrg/log"
	"github.com/sasha-s/go-csync"
)

//goland:noinspection GoUnusedExportedFunction
func NewLimiter(config *Config) Limiter {
	if config == nil {
		config = &DefaultConfig
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	return &limiterImpl{
		config: *config,
	}
}

//goland:noinspection GoNameStartsWithPackageName
type limiterImpl struct {
	csync.Mutex

	reset     time.Time
	remaining int

	config Config
}

func (r *limiterImpl) Logger() log.Logger {
	return r.config.Logger
}

func (r *limiterImpl) Close(ctx context.Context) {
	// TODO: wait for all buckets to unlock
}

func (r *limiterImpl) Config() Config {
	return r.config
}

func (r *limiterImpl) Wait(ctx context.Context) error {
	r.Logger().Debug("locking gateway rate limiter")
	if err := r.CLock(ctx); err != nil {
		return err
	}

	now := time.Now()

	if r.reset.IsZero() {
		r.reset = now.Add(time.Minute)
	}

	var until time.Time

	if r.remaining == 0 && r.reset.After(now) {
		until = r.reset
	}

	if until.After(now) {
		// TODO: do we want to return early when we know rate limit bigger than ctx deadline?
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return context.DeadlineExceeded
		}

		select {
		case <-ctx.Done():
			r.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (r *limiterImpl) Unlock() {
	r.Logger().Debug("unlocking gateway rate limiter")
	now := time.Now()
	if r.reset.Before(now) {
		r.reset = now.Add(time.Minute)
		r.remaining = r.Config().CommandsPerMinute
	}
	r.Mutex.Unlock()
}
