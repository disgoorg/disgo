package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func defaultClientConfig() clientConfig {
	return clientConfig{
		Logger:     slog.Default(),
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
		URL:        fmt.Sprintf("%sv%d", API, Version),
	}
}

type clientConfig struct {
	Logger                *slog.Logger
	HTTPClient            *http.Client
	RateLimiter           RateLimiter
	RateLimiterConfigOpts []RateLimiterConfigOpt
	URL                   string
	UserAgent             string
}

// ClientConfigOpt can be used to supply optional parameters to NewClient
type ClientConfigOpt func(config *clientConfig)

func (c *clientConfig) apply(opts []ClientConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "rest_client"))
	if c.RateLimiter == nil {
		c.RateLimiter = NewRateLimiter(c.RateLimiterConfigOpts...)
	}
}

// WithLogger applies a custom logger to the rest rate limiter
func WithLogger(logger *slog.Logger) ClientConfigOpt {
	return func(config *clientConfig) {
		config.Logger = logger
	}
}

// WithHTTPClient applies a custom http.Client to the rest rate limiter
func WithHTTPClient(httpClient *http.Client) ClientConfigOpt {
	return func(config *clientConfig) {
		config.HTTPClient = httpClient
	}
}

// WithRateLimiter applies a custom RateLimiter to the rest client
func WithRateLimiter(rateLimiter RateLimiter) ClientConfigOpt {
	return func(config *clientConfig) {
		config.RateLimiter = rateLimiter
	}
}

// WithRateLimiterConfigOpts applies RateLimiterConfigOpt to the RateLimiter
func WithRateLimiterConfigOpts(opts ...RateLimiterConfigOpt) ClientConfigOpt {
	return func(config *clientConfig) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}

// WithDefaultRateLimiterConfigOpts applies [RateLimiterConfigOpt] to the RateLimiter and prepend the options to the existing ones.
func WithDefaultRateLimiterConfigOpts(opts ...RateLimiterConfigOpt) ClientConfigOpt {
	return func(config *clientConfig) {
		config.RateLimiterConfigOpts = append(opts, config.RateLimiterConfigOpts...)
	}
}

// WithURL sets the api url for all requests
func WithURL(url string) ClientConfigOpt {
	return func(config *clientConfig) {
		config.URL = url
	}
}

// WithUserAgent sets the user agent for all requests
func WithUserAgent(userAgent string) ClientConfigOpt {
	return func(config *clientConfig) {
		config.UserAgent = userAgent
	}
}
