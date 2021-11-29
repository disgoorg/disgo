package rest

import (
	"fmt"
	"net/http"

	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/rest/rrate"
	"github.com/DisgoOrg/log"
)

// DefaultConfig is the configuration which is used by default
var DefaultConfig = Config{
	HTTPClient:        http.DefaultClient,
	RateLimiterConfig: &rrate.DefaultConfig,
	Headers:           http.Header{},
	UserAgent:         fmt.Sprintf("DiscordBot (%s, %s)", info.GitHub, info.Version),
}

// Config is the configuration for the rest client
type Config struct {
	Logger            log.Logger
	HTTPClient        *http.Client
	RateLimiter       rrate.Limiter
	RateLimiterConfig *rrate.Config
	Headers           http.Header
	UserAgent         string
}

// ConfigOpt can be used to supply optional parameters to NewClient
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger applies a custom logger to the rest rate limiter
//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithHTTPClient applies a custom http.Client to the rest rate limiter
//goland:noinspection GoUnusedExportedFunction
func WithHTTPClient(httpClient *http.Client) ConfigOpt {
	return func(config *Config) {
		config.HTTPClient = httpClient
	}
}

// WithRateLimiter applies a custom rrate.Limiter to the rest client
//goland:noinspection GoUnusedExportedFunction
func WithRateLimiter(rateLimiter rrate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

// WithRateLimiterConfig applies a custom logger to the rest rate limiter
//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfig(rateLimiterConfig rrate.Config) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfig = &rateLimiterConfig
	}
}

// WithRateLimiterConfigOpts applies rrate.ConfigOpt for the rrate.Limiter to the rest rate limiter
//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfigOpts(opts ...rrate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &rrate.DefaultConfig
		}
		config.RateLimiterConfig.Apply(opts)
	}
}

// WithHeaders adds a custom header to all requests
//goland:noinspection GoUnusedExportedFunction
func WithHeaders(headers http.Header) ConfigOpt {
	return func(config *Config) {
		config.Headers = headers
	}
}

// WithUserAgent sets the user agent for all requests
//goland:noinspection GoUnusedExportedFunction
func WithUserAgent(userAgent string) ConfigOpt {
	return func(config *Config) {
		config.UserAgent = userAgent
	}
}
