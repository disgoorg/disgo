package rate

import (
	"context"
	"errors"
	"net/http"

	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

var ErrCtxTimeout = errors.New("rate limit exceeds context deadline")

var DefaultConfig = Config{
	MaxRetries: 10,
}

type Config struct {
	Logger     log.Logger
	MaxRetries int
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithMaxRetries(maxRetries int) ConfigOpt {
	return func(config *Config) {
		config.MaxRetries = maxRetries
	}
}

type Limiter interface {
	Logger() log.Logger
	Close(ctx context.Context)
	Config() Config
	WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error
	UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error
}
