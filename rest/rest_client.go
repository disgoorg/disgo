package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/disgoorg/disgo/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/log"
)

// NewClient constructs a new Client with the given Config struct
func NewClient(botToken string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)

	config.RateLimiter.Reset()

	return &clientImpl{botToken: botToken, config: *config}
}

// Client allows doing requests to different endpoints
type Client interface {
	// Logger returns the logger the rest client uses
	Logger() log.Logger

	// HTTPClient returns the http.Client the rest client uses
	HTTPClient() *http.Client

	// RateLimiter returns the rrate.RateLimiter the rest client uses
	RateLimiter() RateLimiter

	// Close closes the rest client and awaits all pending requests to finish. You can use a cancelling context to abort the waiting
	Close(ctx context.Context)

	// Do makes a request to the given route.CompiledAPIRoute and marshals the given any as json and unmarshalls the response into the given interface
	Do(route *route.CompiledAPIRoute, rqBody any, rsBody any, opts ...RequestOpt) error
}

type clientImpl struct {
	botToken string
	config   Config
}

func (c *clientImpl) Close(ctx context.Context) {
	c.config.RateLimiter.Close(ctx)
	c.config.HTTPClient.CloseIdleConnections()
}

func (c *clientImpl) Logger() log.Logger {
	return c.config.Logger
}

func (c *clientImpl) HTTPClient() *http.Client {
	return c.config.HTTPClient
}

func (c *clientImpl) RateLimiter() RateLimiter {
	return c.config.RateLimiter
}

func (c *clientImpl) retry(cRoute *route.CompiledAPIRoute, rqBody any, rsBody any, tries int, opts []RequestOpt) error {
	var (
		rqURL       = cRoute.URL()
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
		c.Logger().Tracef("request to %s, body: %s", rqURL, string(rawRqBody))
	}

	rq, err := http.NewRequest(cRoute.APIRoute.Method().String(), rqURL, bytes.NewReader(rawRqBody))
	if err != nil {
		return err
	}

	rq.Header.Set("User-Agent", c.config.UserAgent)
	if contentType != "" {
		rq.Header.Set("Content-Type", contentType)
	}

	if cRoute.APIRoute.NeedsBotAuth() {
		// add token opt to the start, so you can override it
		opts = append([]RequestOpt{WithToken(discord.TokenTypeBot, c.botToken)}, opts...)
	}

	config := DefaultRequestConfig(rq)
	config.Apply(opts)

	if config.Delay > 0 {
		timer := time.NewTimer(config.Delay)
		defer timer.Stop()
		select {
		case <-config.Ctx.Done():
			return config.Ctx.Err()
		case <-timer.C:
		}
	}

	// wait for rate limits
	err = c.RateLimiter().WaitBucket(config.Ctx, cRoute)
	if err != nil {
		return fmt.Errorf("error locking bucket in rest client: %w", err)
	}
	rq = rq.WithContext(config.Ctx)

	for _, check := range config.Checks {
		if !check() {
			_ = c.RateLimiter().UnlockBucket(cRoute, nil)
			return discord.ErrCheckFailed
		}
	}

	rs, err := c.HTTPClient().Do(config.Request)
	if err != nil {
		_ = c.RateLimiter().UnlockBucket(cRoute, nil)
		return fmt.Errorf("error doing request in rest client: %w", err)
	}

	if err = c.RateLimiter().UnlockBucket(cRoute, rs.Header); err != nil {
		// TODO: should we maybe retry here?
		return fmt.Errorf("error unlocking bucket in rest client: %w", err)
	}

	var rawRsBody []byte
	if rs.Body != nil {
		if rawRsBody, err = io.ReadAll(rs.Body); err != nil {
			return fmt.Errorf("error reading response body in rest client: %w", err)
		}
		c.Logger().Tracef("response from %s, code %d, body: %s", rqURL, rs.StatusCode, string(rawRsBody))
	}

	switch rs.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if rsBody != nil && rs.Body != nil {
			if err = json.Unmarshal(rawRsBody, rsBody); err != nil {
				wErr := fmt.Errorf("error unmarshalling response body: %w", err)
				c.Logger().Error(wErr)
				return wErr
			}
		}
		return nil

	case http.StatusTooManyRequests:
		if tries >= c.RateLimiter().MaxRetries() {
			return NewError(rq, rawRqBody, rs, rawRsBody)
		}
		return c.retry(cRoute, rqBody, rsBody, tries+1, opts)

	default:
		return NewError(rq, rawRqBody, rs, rawRsBody)
	}
}

func (c *clientImpl) Do(cRoute *route.CompiledAPIRoute, rqBody any, rsBody any, opts ...RequestOpt) error {
	return c.retry(cRoute, rqBody, rsBody, 1, opts)
}
