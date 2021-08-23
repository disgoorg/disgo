package rest

import (
	"context"
	"net/http"
)

type RequestConfig struct {
	Request *http.Request
	Ctx     context.Context
}

type RequestOpt func(config RequestConfig) RequestConfig

func WithCtx(ctx context.Context) RequestOpt {
	return func(config RequestConfig) RequestConfig {
		config.Ctx = ctx
		return config
	}
}

func WithReason(reason string) RequestOpt {
	return func(config RequestConfig) RequestConfig {
		config.Request.Header.Set("X-Audit-Log-Reason", reason)
		return config
	}
}

func WithHeader(key string, value string) RequestOpt {
	return func(config RequestConfig) RequestConfig {
		config.Request.Header.Set(key, value)
		return config
	}
}

func applyRequestOpts(config RequestConfig, opts []RequestOpt) RequestConfig {
	for _, opt := range opts {
		config = opt(config)
	}
	return config
}
