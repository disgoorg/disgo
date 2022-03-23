package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/disgoorg/disgo/discord"
)

func DefaultRequestConfig(rq *http.Request) *RequestConfig {
	return &RequestConfig{
		Request: rq,
		Ctx:     context.TODO(),
	}
}

// RequestConfig are additional options for the request
type RequestConfig struct {
	Request   *http.Request
	Ctx       context.Context
	Checks    []Check
	Delay     time.Duration
	TokenType discord.TokenType
	Token     string
}

// Check is a function which gets executed right before a request is made
type Check func() bool

// RequestOpt can be used to supply optional parameters to Client.Do
type RequestOpt func(config *RequestConfig)

// Apply applies the given RequestOpt(s) to the RequestConfig & sets the context if none is set
func (c *RequestConfig) Apply(opts []RequestOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.Ctx == nil {
		c.Ctx = context.TODO()
	}
}

// WithCtx applies a custom context to the request
//goland:noinspection GoUnusedExportedFunction
func WithCtx(ctx context.Context) RequestOpt {
	return func(config *RequestConfig) {
		config.Ctx = ctx
	}
}

// WithCheck adds a new check to the request
//goland:noinspection GoUnusedExportedFunction
func WithCheck(check Check) RequestOpt {
	return func(config *RequestConfig) {
		config.Checks = append(config.Checks, check)
	}
}

// WithDelay applies a delay to the request
//goland:noinspection GoUnusedExportedFunction
func WithDelay(delay time.Duration) RequestOpt {
	return func(config *RequestConfig) {
		config.Delay = delay
	}
}

// WithReason adds a reason header to the request. Not all discord endpoints support this
//goland:noinspection GoUnusedExportedFunction
func WithReason(reason string) RequestOpt {
	return func(config *RequestConfig) {
		config.Request.Header.Set("X-Audit-Log-Reason", reason)
	}
}

// WithHeader adds a custom header to the request
//goland:noinspection GoUnusedExportedFunction
func WithHeader(key string, value string) RequestOpt {
	return func(config *RequestConfig) {
		config.Request.Header.Set(key, value)
	}
}

// WithQueryParam applies a custom query parameter to the request
//goland:noinspection GoUnusedExportedFunction
func WithQueryParam(param string, value any) RequestOpt {
	return func(config *RequestConfig) {
		values := config.Request.URL.Query()
		values.Add(param, fmt.Sprint(value))
		config.Request.URL.RawQuery = values.Encode()
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithToken(tokenType discord.TokenType, token string) RequestOpt {
	return func(config *RequestConfig) {
		config.TokenType = tokenType
		config.Token = token
	}
}
