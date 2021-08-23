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
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

var DefaultConfig = Config{
	Headers:   nil,
	UserAgent: fmt.Sprintf("DiscordBot (%s, %s)", info.GitHub, info.Version),
}

type Config struct {
	Headers   http.Header
	UserAgent string
}

// NewClient constructs a new Client with the given http.Client, log.Logger & useragent
//goland:noinspection GoUnusedExportedFunction
func NewClient(logger log.Logger, httpClient *http.Client, rateLimiter rate.RateLimiter, config *Config) Client {
	if logger == nil {
		logger = log.Default()
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if rateLimiter == nil {
		rateLimiter = rate.NewRateLimiter(logger, nil)
	}
	if config == nil {
		config = &DefaultConfig
	}
	return &ClientImpl{
		logger:      logger,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		config:      *config,
	}
}

// Client allows doing requests to different endpoints
type Client interface {
	Close()
	Logger() log.Logger
	HTTPClient() *http.Client
	RateLimiter() rate.RateLimiter
	Config() Config
	Do(route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, opts ...RequestOpt) Error
}

type ClientImpl struct {
	logger      log.Logger
	config      Config
	httpClient  *http.Client
	rateLimiter rate.RateLimiter
}

func (c *ClientImpl) Close() {
	c.httpClient.CloseIdleConnections()
}

func (c *ClientImpl) Logger() log.Logger {
	return c.logger
}

func (c *ClientImpl) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *ClientImpl) RateLimiter() rate.RateLimiter {
	return c.rateLimiter
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

	config := applyRequestOpts(RequestConfig{Request: rq}, opts)

	if config.Ctx != nil {
		// wait for rate limits
		err = c.rateLimiter.WaitBucket(config.Ctx, cRoute)
		if err != nil {
			return NewError(nil, err)
		}
	}

	rs, err := c.httpClient.Do(config.Request)
	if err != nil {
		return NewError(rs, err)
	}

	if err = c.rateLimiter.UnlockBucket(cRoute, rs.Header); err != nil {
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
