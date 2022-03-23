package rest

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/disgoorg/disgo/json"
	"github.com/pkg/errors"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/disgo/rest/rrate"
	"github.com/disgoorg/log"
)

// NewClient constructs a new Client with the given Config struct
//goland:noinspection GoUnusedExportedFunction
func NewClient(botToken string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)

	return &clientImpl{botToken: botToken, config: *config}
}

// Client allows doing requests to different endpoints
type Client interface {
	// Logger returns the logger the rest client uses
	Logger() log.Logger

	// HTTPClient returns the http.Client the rest client uses
	HTTPClient() *http.Client

	// RateLimiter returns the rrate.Limiter the rest client uses
	RateLimiter() rrate.Limiter
	// Config returns the Config the rest client uses
	Config() Config

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

func (c *clientImpl) RateLimiter() rrate.Limiter {
	return c.config.RateLimiter
}

func (c *clientImpl) Config() Config {
	return c.config
}

func (c *clientImpl) retry(cRoute *route.CompiledAPIRoute, rqBody any, rsBody any, tries int, opts []RequestOpt) error {
	rqURL := cRoute.URL()
	var rawRqBody []byte
	var err error
	var contentType string

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
				return errors.Wrap(err, "failed to marshal request body")
			}
		}
		c.Logger().Tracef("request to %s, body: %s", rqURL, string(rawRqBody))
	}

	rq, err := http.NewRequest(cRoute.APIRoute.Method().String(), rqURL, bytes.NewReader(rawRqBody))
	if err != nil {
		return err
	}

	rq.Header.Set("User-Agent", c.Config().UserAgent)
	if contentType != "" {
		rq.Header.Set("Content-Type", contentType)
	}

	var (
		tokenType discord.TokenType
		token     string
	)

	if cRoute.APIRoute.NeedsBotAuth() {
		tokenType = discord.TokenTypeBot
		token = c.botToken
	}

	config := DefaultRequestConfig(rq, tokenType, token)
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
		return errors.Wrap(err, "error locking bucket in rest client")
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
		return errors.Wrap(err, "error doing request in rest client")
	}

	if err = c.RateLimiter().UnlockBucket(cRoute, rs.Header); err != nil {
		// TODO: should we maybe retry here?
		return errors.Wrap(err, "error unlocking bucket in rest client")
	}

	var rawRsBody []byte
	if rs.Body != nil {
		if rawRsBody, err = ioutil.ReadAll(rs.Body); err != nil {
			return errors.Wrap(err, "error reading response body in rest client")
		}
		c.Logger().Tracef("response from %s, code %d, body: %s", rqURL, rs.StatusCode, string(rawRsBody))
	}

	switch rs.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if rsBody != nil && rs.Body != nil {
			if err = json.Unmarshal(rawRsBody, rsBody); err != nil {
				wErr := errors.Wrap(err, "error unmarshalling response body")
				c.Logger().Error(wErr)
				return NewErrorErr(rq, rawRqBody, rs, rawRsBody, wErr)
			}
		}
		return nil

	case http.StatusBadGateway, http.StatusUnauthorized:
		return NewError(rq, rawRqBody, rs, rawRsBody)

	case http.StatusTooManyRequests:
		if tries >= c.RateLimiter().Config().MaxRetries {
			return NewError(rq, rawRqBody, rs, rawRsBody)
		}
		return c.retry(cRoute, rqBody, rsBody, tries+1, opts)

	default:
		var v discord.APIError
		if err = json.Unmarshal(rawRsBody, &v); err != nil {
			return errors.Wrap(err, "error unmarshalling error response body")
		}
		return NewErrorAPIErr(rq, rawRqBody, rs, rawRsBody, v)
	}
}

func (c *clientImpl) Do(cRoute *route.CompiledAPIRoute, rqBody any, rsBody any, opts ...RequestOpt) error {
	return c.retry(cRoute, rqBody, rsBody, 1, opts)
}

func applyAuthHeader(tokenType discord.TokenType, token string, opts []RequestOpt) []RequestOpt {
	return append(opts, WithHeader("authorization", tokenType.Apply(token)))
}
