package rest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/DisgoOrg/disgo/json"
	"github.com/pkg/errors"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/disgo/rest/rrate"
	"github.com/DisgoOrg/log"
)

// NewClient constructs a new Client with the given Config struct
//goland:noinspection GoUnusedExportedFunction
func NewClient(config *Config) Client {
	if config == nil {
		config = &DefaultConfig
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}
	if config.RateLimiterConfig == nil {
		config.RateLimiterConfig = &rrate.DefaultConfig
	}
	if config.RateLimiter == nil {
		config.RateLimiter = rrate.NewLimiter(config.RateLimiterConfig)
	}
	return &clientImpl{config: *config}
}

// Client allows doing requests to different endpoints
type Client interface {
	Close()
	Logger() log.Logger
	HTTPClient() *http.Client
	RateLimiter() rrate.Limiter
	Config() Config
	Do(route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, opts ...RequestOpt) error
}

type clientImpl struct {
	config Config
}

func (c *clientImpl) Close() {
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

func (c *clientImpl) retry(cRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, tries int, opts []RequestOpt) error {
	rqURL := cRoute.URL()
	rqBuffer := &bytes.Buffer{}
	var rawRqBody []byte
	var contentType string

	if rqBody != nil {
		var buffer *bytes.Buffer
		switch v := rqBody.(type) {
		case *discord.MultipartBuffer:
			contentType = v.ContentType
			buffer = v.Buffer

		case url.Values:
			contentType = "application/x-www-form-urlencoded"
			buffer = bytes.NewBufferString(v.Encode())

		default:
			contentType = "application/json"
			buffer = &bytes.Buffer{}
			err := json.NewEncoder(buffer).Encode(rqBody)
			if err != nil {
				return err
			}
		}
		rawRqBody, _ = ioutil.ReadAll(io.TeeReader(buffer, rqBuffer))
		c.Logger().Debugf("request to %s, body: %s", rqURL, string(rawRqBody))
	}

	rq, err := http.NewRequest(cRoute.APIRoute.Method().String(), rqURL, rqBuffer)
	if err != nil {
		return err
	}

	// write all headers to the request
	if headers := c.Config().Headers; headers != nil {
		for key, values := range headers {
			for _, value := range values {
				rq.Header.Add(key, value)
			}
		}
	}
	rq.Header.Set("User-Agent", c.Config().UserAgent)
	if contentType != "" {
		rq.Header.Set("Content-Type", contentType)
	}

	config := &RequestConfig{Request: rq}
	config.Apply(opts)

	if config.Delay > 0 {
		select {
		case <-config.Ctx.Done():
			return config.Ctx.Err()
		case <-time.After(config.Delay):
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
		buffer := &bytes.Buffer{}
		rawRsBody, _ = ioutil.ReadAll(io.TeeReader(rs.Body, buffer))
		rs.Body = ioutil.NopCloser(buffer)
		c.Logger().Debugf("response from %s, code %d, body: %s", rqURL, rs.StatusCode, string(rawRsBody))
	}

	switch rs.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if rsBody != nil && rs.Body != nil {
			if err = json.NewDecoder(rs.Body).Decode(rsBody); err != nil {
				wErr := errors.Wrap(err, "error unmarshalling response body")
				c.Logger().Error(wErr)
				return NewErrorErr(rq, rawRqBody, rs, rawRqBody, wErr)
			}
		}
		return nil

	case http.StatusBadGateway, http.StatusUnauthorized:
		return NewError(rq, rawRqBody, rs, rawRqBody)

	case http.StatusTooManyRequests:
		if tries >= c.RateLimiter().Config().MaxRetries {
			return NewError(rq, rawRqBody, rs, rawRqBody)
		}
		return c.retry(cRoute, rqBody, rsBody, tries+1, opts)

	default:
		var v discord.APIError
		if err = json.NewDecoder(rs.Body).Decode(&v); err != nil {
			return errors.Wrap(err, "error unmarshalling error response body")
		}
		return NewErrorAPIErr(rq, rawRqBody, rs, rawRqBody, v)
	}
}

func (c *clientImpl) Do(cRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, opts ...RequestOpt) error {
	return c.retry(cRoute, rqBody, rsBody, 1, opts)
}
