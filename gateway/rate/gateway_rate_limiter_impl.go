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

func (l *limiterImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *limiterImpl) Close(ctx context.Context) {
	_ = l.CLock(ctx)
	l.Unlock()
}

func (l *limiterImpl) Config() Config {
	return l.config
}

func (l *limiterImpl) Wait(ctx context.Context) error {
	l.Logger().Debug("locking gateway rate limiter")
	if err := l.CLock(ctx); err != nil {
		return err
	}

	now := time.Now()

	var until time.Time

	if l.remaining == 0 && l.reset.After(now) {
		until = l.reset
	}

	if until.After(now) {
		// TODO: do we want to return early when we know rate limit bigger than ctx deadline?
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return context.DeadlineExceeded
		}

		select {
		case <-ctx.Done():
			l.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (l *limiterImpl) Unlock() {
	l.Logger().Debug("unlocking gateway rate limiter")
	now := time.Now()
	if l.reset.Before(now) {
		l.reset = now.Add(time.Minute)
		l.remaining = l.Config().CommandsPerMinute
	}
	l.Mutex.Unlock()
}
