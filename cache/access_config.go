package cache

import "context"

type accessConfig struct {
	Ctx context.Context
}

// AccessOpt can be used to supply optional parameters to various cache operations
type AccessOpt func(config *accessConfig)

// WithCtx applies a custom context to the cache operation
func WithCtx(ctx context.Context) AccessOpt {
	return func(config *accessConfig) {
		config.Ctx = ctx
	}
}

func resolveAccessConfig(opts []AccessOpt) *accessConfig {
	cfg := &accessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}
	return cfg
}
