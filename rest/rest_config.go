package rest

import (
	"fmt"
	"net/http"

	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/rest/rrate"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	HTTPClient:        http.DefaultClient,
	RateLimiterConfig: &rrate.DefaultConfig,
	Headers:           http.Header{},
	UserAgent:         fmt.Sprintf("DiscordBot (%s, %s)", info.GitHub, info.Version),
}

type Config struct {
	Logger            log.Logger
	HTTPClient        *http.Client
	RateLimiter       rrate.Limiter
	RateLimiterConfig *rrate.Config
	Headers           http.Header
	UserAgent         string
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

func WithHTTPClient(httpClient *http.Client) ConfigOpt {
	return func(config *Config) {
		config.HTTPClient = httpClient
	}
}

func WithRateLimiter(rateLimiter rrate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

func WithRateLimiterConfig(rateLimiterConfig rrate.Config) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfig = &rateLimiterConfig
	}
}

func WithRateLimiterConfigOpts(opts ...rrate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &rrate.DefaultConfig
		}
		config.RateLimiterConfig.Apply(opts)
	}
}

func WithHeaders(headers http.Header) ConfigOpt {
	return func(config *Config) {
		config.Headers = headers
	}
}

func WithUserAgent(userAgent string) ConfigOpt {
	return func(config *Config) {
		config.UserAgent = userAgent
	}
}
