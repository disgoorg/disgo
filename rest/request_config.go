package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
)

// DefaultRequestConfig creates a new RequestConfig with the given http.Request and a default context
func DefaultRequestConfig(rq *http.Request) RequestConfig {
	return RequestConfig{
		Request: rq,
		Ctx:     context.TODO(),
	}
}

// RequestConfig is the configuration for rest requests
type RequestConfig struct {
	// Request is the http.Request to be used for the request
	Request *http.Request
	// Ctx applies a custom context to the request
	Ctx context.Context
	// Checks is a list of checks to be executed right before a request is made
	Checks []Check
	// Delay applies a delay to the request
	Delay time.Duration
}

// Check is a function which gets executed right before a request is made
type Check func() bool

// RequestOpt can be used to supply optional parameters to Client.Do
type RequestOpt func(config *RequestConfig)

func (c *RequestConfig) apply(opts []RequestOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.Ctx == nil {
		c.Ctx = context.TODO()
	}
}

// WithCtx applies a custom context to the request
func WithCtx(ctx context.Context) RequestOpt {
	return func(config *RequestConfig) {
		config.Ctx = ctx
	}
}

// WithCheck adds a new check to the request
func WithCheck(check Check) RequestOpt {
	return func(config *RequestConfig) {
		config.Checks = append(config.Checks, check)
	}
}

// WithDelay applies a delay to the request
func WithDelay(delay time.Duration) RequestOpt {
	return func(config *RequestConfig) {
		config.Delay = delay
	}
}

// WithHeader adds a custom header to the request
func WithHeader(key string, value string) RequestOpt {
	return func(config *RequestConfig) {
		config.Request.Header.Set(key, value)
	}
}

// WithReason adds a reason header to the request. Not all discord endpoints support this
func WithReason(reason string) RequestOpt {
	reason = strings.ReplaceAll(url.QueryEscape(reason), "+", " ")
	return WithHeader("X-Audit-Log-Reason", reason)
}

// WithDiscordLocale adds the X-Discord-Locale header with the passed locale to the request
func WithDiscordLocale(locale discord.Locale) RequestOpt {
	return WithHeader("X-Discord-Locale", locale.Code())
}

// WithToken adds the Authorization header with the passed token to the request
func WithToken(tokenType discord.TokenType, token string) RequestOpt {
	return WithHeader("Authorization", tokenType.Apply(token))
}

// WithQueryParam applies a custom query parameter to the request
func WithQueryParam(param string, value any) RequestOpt {
	return func(config *RequestConfig) {
		values := config.Request.URL.Query()
		values.Add(param, fmt.Sprint(value))
		config.Request.URL.RawQuery = values.Encode()
	}
}
