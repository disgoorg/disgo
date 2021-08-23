package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/disgo/rest/route"
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
		config.RateLimiterConfig = &rate.DefaultConfig
	}
	if config.RateLimiter == nil {
		config.RateLimiter = rate.NewLimiter(config.Logger, config.RateLimiterConfig)
	}
	return &ClientImpl{config: *config}
}

// Client allows doing requests to different endpoints
type Client interface {
	Close()
	Logger() log.Logger
	HTTPClient() *http.Client
	RateLimiter() rate.Limiter
	Config() Config
	Do(route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, opts ...RequestOpt) Error
}

type ClientImpl struct {
	config Config
}

func (c *ClientImpl) Close() {
	c.config.HTTPClient.CloseIdleConnections()
}

func (c *ClientImpl) Logger() log.Logger {
	return c.config.Logger
}

func (c *ClientImpl) HTTPClient() *http.Client {
	return c.config.HTTPClient
}

func (c *ClientImpl) RateLimiter() rate.Limiter {
	return c.config.RateLimiter
}

func (c *ClientImpl) Config() Config {
	return c.config
}

func (c *ClientImpl) retry(cRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, tries int, opts []RequestOpt) Error {
	rqURL := cRoute.URL()
	rqBuffer := &bytes.Buffer{}
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
				return NewError(nil, err)
			}
		}
		body, _ := ioutil.ReadAll(io.TeeReader(buffer, rqBuffer))
		c.Logger().Debugf("request to %s, body: %s", rqURL, string(body))
	}

	rq, err := http.NewRequest(cRoute.Method().String(), rqURL, rqBuffer)
	if err != nil {
		return NewError(nil, err)
	}

	if headers := c.Config().Headers; headers != nil {
		rq.Header = headers
	}
	rq.Header.Set("User-Agent", c.Config().UserAgent)
	if contentType != "" {
		rq.Header.Set("Content-Type", contentType)
	}

	config := &RequestConfig{Request: rq}
	config.Apply(opts)

	if config.Ctx != nil {
		// wait for rate limits
		err = c.RateLimiter().WaitBucket(config.Ctx, cRoute)
		if err != nil {
			return NewError(nil, err)
		}
		rq = rq.WithContext(config.Ctx)
	}

	rs, err := c.HTTPClient().Do(config.Request)
	if err != nil {
		return NewError(rs, err)
	}

	if err = c.RateLimiter().UnlockBucket(cRoute, rs.Header); err != nil {
		// TODO: should we maybe retry here?
		return NewError(rs, err)
	}

	if rs.Body != nil {
		buffer := &bytes.Buffer{}
		body, _ := ioutil.ReadAll(io.TeeReader(rs.Body, buffer))
		rs.Body = ioutil.NopCloser(buffer)
		c.Logger().Debugf("response from %s, code %d, body: %s", rqURL, rs.StatusCode, string(body))
	}

	switch rs.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		if rsBody != nil && rs.Body != nil {
			if err = json.NewDecoder(rs.Body).Decode(rsBody); err != nil {
				c.Logger().Errorf("error unmarshalling response. error: %s", err)
				return NewError(rs, err)
			}
		}
		return nil

	case http.StatusBadGateway:
		c.Logger().Error(discord.ErrBadGateway)
		return NewError(rs, discord.ErrBadGateway)

	case http.StatusBadRequest:
		c.Logger().Error(discord.ErrBadRequest)
		return NewError(rs, discord.ErrBadRequest)

	case http.StatusTooManyRequests:
		c.Logger().Error(discord.ErrRatelimited)
		if tries >= c.RateLimiter().Config().MaxRetries {
			return NewError(rs, discord.ErrRatelimited)
		}
		return c.retry(cRoute, rqBody, rsBody, tries+1, opts)

	case http.StatusUnauthorized:
		c.Logger().Error(discord.ErrUnauthorized)
		return NewError(rs, discord.ErrUnauthorized)

	default:
		body, _ := ioutil.ReadAll(rq.Body)
		return NewError(rs, fmt.Errorf("request to %s failed. statuscode: %d, body: %s", rq.URL, rs.StatusCode, body))
	}
}

func (c *ClientImpl) Do(cRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, opts ...RequestOpt) Error {
	return c.retry(cRoute, rqBody, rsBody, 1, opts)
}
