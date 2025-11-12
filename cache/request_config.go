package cache

import "context"

type requestConfig struct {
	Ctx context.Context
}

// RequestOpt can be used to supply optional parameters to various cache operations
type RequestOpt func(config *requestConfig)

// WithCtx applies a custom context to the cache operation
func WithCtx(ctx context.Context) RequestOpt {
	return func(config *requestConfig) {
		config.Ctx = ctx
	}
}
