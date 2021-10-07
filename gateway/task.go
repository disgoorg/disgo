package gateway

import (
	"context"
	"time"
)

type TaskConfig struct {
	Ctx    context.Context
	Checks []Check
	Delay  time.Duration
}

type Check func() bool

type TaskOpt func(config *TaskConfig)

func (c *TaskConfig) Apply(opts []TaskOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.Ctx == nil {
		c.Ctx = context.TODO()
	}
}

func WithCtx(ctx context.Context) TaskOpt {
	return func(config *TaskConfig) {
		config.Ctx = ctx
	}
}

func WithCheck(check Check) TaskOpt {
	return func(config *TaskConfig) {
		config.Checks = append(config.Checks, check)
	}
}

func WithDelay(delay time.Duration) TaskOpt {
	return func(config *TaskConfig) {
		config.Delay = delay
	}
}
