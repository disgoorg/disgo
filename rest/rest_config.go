package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// DefaultConfig is the configuration which is used by default
func DefaultConfig() *Config {
	return &Config{
		Logger:     slog.Default(),
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
		URL:        fmt.Sprintf("%sv%d", API, Version),
	}
}

// Config is the configuration for the rest client
type Config struct {
	Logger                *slog.Logger
	HTTPClient            *http.Client
	RateLimiter           RateLimiter
	RateLimiterConfigOpts []RateLimiterConfigOpt
	URL                   string
	UserAgent             string
}

// ConfigOpt can be used to supply optional parameters to NewClient
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.RateLimiter == nil {
		c.RateLimiter = NewRateLimiter(c.RateLimiterConfigOpts...)
	}
}

// WithLogger applies a custom logger to the rest rate limiter
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithHTTPClient applies a custom http.Client to the rest rate limiter
func WithHTTPClient(httpClient *http.Client) ConfigOpt {
	return func(config *Config) {
		config.HTTPClient = httpClient
	}
}

// WithRateLimiter applies a custom RateLimiter to the rest client
func WithRateLimiter(rateLimiter RateLimiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

// WithRateLimiterConfigOpts applies RateLimiterConfigOpt to the RateLimiter
func WithRateLimiterConfigOpts(opts ...RateLimiterConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}

// WithURL sets the api url for all requests
func WithURL(url string) ConfigOpt {
	return func(config *Config) {
		config.URL = url
	}
}

// WithUserAgent sets the user agent for all requests
func WithUserAgent(userAgent string) ConfigOpt {
	return func(config *Config) {
		config.UserAgent = userAgent
	}
}
