package rest

import (
	"context"
	"net/http"
	"time"
)

type RequestConfig struct {
	Request *http.Request
	Ctx     context.Context
	Checks  []Check
	Delay   time.Duration
}

type Check func() bool

type RequestOpt func(config *RequestConfig)

func (c *RequestConfig) Apply(opts []RequestOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.Ctx == nil {
		c.Ctx = context.TODO()
	}
}

func WithCtx(ctx context.Context) RequestOpt {
	return func(config *RequestConfig) {
		config.Ctx = ctx
	}
}

func WithCheck(check Check) RequestOpt {
	return func(config *RequestConfig) {
		config.Checks = append(config.Checks, check)
	}
}

func WithDelay(delay time.Duration) RequestOpt {
	return func(config *RequestConfig) {
		config.Delay = delay
	}
}

func WithReason(reason string) RequestOpt {
	return func(config *RequestConfig) {
		config.Request.Header.Set("X-Audit-Log-Reason", reason)
	}
}

func WithHeader(key string, value string) RequestOpt {
	return func(config *RequestConfig) {
		config.Request.Header.Set(key, value)
	}
}
