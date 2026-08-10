package cache

import "context"

// DefaultRequestConfig creates a new RequestConfig with a default context
func DefaultRequestConfig() RequestConfig {
	return RequestConfig{
		Ctx: context.TODO(),
	}
}

// RequestConfig is the configuration for cache requests
type RequestConfig struct {
	// Ctx applies a custom context to the cache operation
	// The default cache implementation ignores the context
	Ctx context.Context
}

// RequestOpt can be used to supply optional parameters to various cache operations
type RequestOpt func(config *RequestConfig)

// WithCtx applies a custom context to the cache operation
func WithCtx(ctx context.Context) RequestOpt {
	return func(config *RequestConfig) {
		config.Ctx = ctx
	}
}
