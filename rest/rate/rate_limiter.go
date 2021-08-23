package rate

import (
	"context"
	"errors"
	"net/http"

	"github.com/DisgoOrg/disgo/rest/route"
)

var ErrCtxTimeout = errors.New("rate limit exceeds context deadline")

var DefaultConfig = Config{
	MaxRetries: 10,
}

type Config struct {
	MaxRetries int
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithMaxRetries(maxRetries int) ConfigOpt {
	return func(config *Config) {
		config.MaxRetries = maxRetries
	}
}

type Limiter interface {
	Close(ctx context.Context)
	Config() Config
	WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error
	UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error
}
