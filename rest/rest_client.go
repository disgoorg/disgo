package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/disgoorg/json/v2"

	"github.com/disgoorg/disgo/discord"
)

// NewClient constructs a new Client with the given config struct
func NewClient(botToken string, opts ...ConfigOpt) Client {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &clientImpl{
		botToken: botToken,
		config:   cfg,
	}
}

// Client allows doing requests to different endpoints
type Client interface {
	// HTTPClient returns the http.Client the rest client uses
	HTTPClient() *http.Client

	// RateLimiter returns the RateLimiter the rest client uses
	RateLimiter() RateLimiter

	// Close closes the rest client and awaits all pending requests to finish. You can use a cancelling context to abort the waiting
	Close(ctx context.Context)

	// Do makes a request to the given CompiledAPIRoute and marshals the given any as json and unmarshalls the response into the given interface
	Do(endpoint *CompiledEndpoint, rqBody any, rsBody any, opts ...RequestOpt) error
}

type clientImpl struct {
	botToken string
	config   config
}

func (c *clientImpl) Close(ctx context.Context) {
	c.config.RateLimiter.Close(ctx)
	c.config.HTTPClient.CloseIdleConnections()
}

func (c *clientImpl) HTTPClient() *http.Client {
	return c.config.HTTPClient
}

func (c *clientImpl) RateLimiter() RateLimiter {
	return c.config.RateLimiter
}

func (c *clientImpl) retry(endpoint *CompiledEndpoint, rqBody any, rsBody any, tries int, opts []RequestOpt) error {
	var (
		rawRqBody   []byte
		err         error
		contentType string
	)

	if rqBody != nil {
		switch v := rqBody.(type) {
		case *discord.MultipartBuffer:
			contentType = v.ContentType
			rawRqBody = v.Buffer.Bytes()

		case url.Values:
			contentType = "application/x-www-form-urlencoded"
			rawRqBody = []byte(v.Encode())

		default:
			contentType = "application/json"
			if rawRqBody, err = json.Marshal(rqBody); err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
		}
		c.config.Logger.Debug("new request", slog.String("endpoint", endpoint.URL), slog.String("body", string(rawRqBody)))
	}

	rq, err := http.NewRequest(endpoint.Endpoint.Method, c.config.URL+endpoint.URL, bytes.NewReader(rawRqBody))
	if err != nil {
		return err
	}

	rq.Header.Set("User-Agent", c.config.UserAgent)
	if contentType != "" {
		rq.Header.Set("Content-Type", contentType)
	}

	if endpoint.Endpoint.BotAuth {
		// add token opt to the start, so you can override it
		opts = append([]RequestOpt{WithToken(discord.TokenTypeBot, c.botToken)}, opts...)
	}

	cfg := defaultRequestConfig(rq)
	cfg.apply(opts)

	if cfg.Delay > 0 {
		timer := time.NewTimer(cfg.Delay)
		defer timer.Stop()
		select {
		case <-cfg.Ctx.Done():
			return cfg.Ctx.Err()
		case <-timer.C:
		}
	}

	// wait for rate limits
	err = c.RateLimiter().WaitBucket(cfg.Ctx, endpoint)
	if err != nil {
		return fmt.Errorf("error locking bucket in rest client: %w", err)
	}
	rq = cfg.Request.WithContext(cfg.Ctx)

	for _, check := range cfg.Checks {
		if !check() {
			_ = c.RateLimiter().UnlockBucket(endpoint, nil)
			return discord.ErrCheckFailed
		}
	}

	rs, err := c.HTTPClient().Do(rq)
	if err != nil {
		_ = c.RateLimiter().UnlockBucket(endpoint, nil)
		return fmt.Errorf("error doing request in rest client: %w", err)
	}

	if err = c.RateLimiter().UnlockBucket(endpoint, rs); err != nil {
		return fmt.Errorf("error unlocking bucket in rest client: %w", err)
	}

	var rawRsBody []byte
	if rs.Body != nil {
		if rawRsBody, err = io.ReadAll(rs.Body); err != nil {
			return fmt.Errorf("error reading response body in rest client: %w", err)
		}
		c.config.Logger.Debug("new response", slog.String("endpoint", endpoint.URL), slog.String("code", rs.Status), slog.String("body", string(rawRsBody)))
	}

	switch rs.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if rsBody != nil && rs.Body != nil {
			if err = json.Unmarshal(rawRsBody, rsBody); err != nil {
				c.config.Logger.Error("error unmarshalling response body", slog.Any("err", err), slog.String("endpoint", endpoint.URL), slog.String("code", rs.Status), slog.String("body", string(rawRsBody)))
				return fmt.Errorf("error unmarshalling response body: %w", err)
			}
		}
		return nil

	case http.StatusTooManyRequests:
		if tries >= c.RateLimiter().MaxRetries() {
			return NewError(rq, rawRqBody, rs, rawRsBody)
		}
		return c.retry(endpoint, rqBody, rsBody, tries+1, opts)

	default:
		return NewError(rq, rawRqBody, rs, rawRsBody)
	}
}

func (c *clientImpl) Do(endpoint *CompiledEndpoint, rqBody any, rsBody any, opts ...RequestOpt) error {
	return c.retry(endpoint, rqBody, rsBody, 1, opts)
}
