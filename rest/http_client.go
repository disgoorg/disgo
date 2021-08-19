package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/disgo/util"
	"github.com/DisgoOrg/log"
)

// Client allows doing requests to different endpoints
type Client interface {
	Close()
	Logger() log.Logger
	HTTPClient() *http.Client
	RateLimiter() rate.RateLimiter
	Config() Config
	Do(ctx context.Context, route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) Error
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
	return &HTTPClientImpl{
		logger:      logger,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		config:      *config,
	}
}

var DefaultConfig = Config{
	Headers:   nil,
	UserAgent: fmt.Sprintf("DiscordBot (%s, %s)", util.GitHub, util.Version),
}

type Config struct {
	Headers   http.Header
	UserAgent string
}

type HTTPClientImpl struct {
	logger      log.Logger
	config      Config
	httpClient  *http.Client
	rateLimiter rate.RateLimiter
}

func (c *HTTPClientImpl) Close() {
	c.httpClient.CloseIdleConnections()
}

func (c *HTTPClientImpl) Logger() log.Logger {
	return c.logger
}

func (c *HTTPClientImpl) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *HTTPClientImpl) RateLimiter() rate.RateLimiter {
	return c.rateLimiter
}

func (c *HTTPClientImpl) Config() Config {
	return c.config
}

func (c *HTTPClientImpl) retry(ctx context.Context, cRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, tries int) Error {
	// wait for rate limits
	err := c.rateLimiter.WaitBucket(ctx, cRoute)
	if err != nil {
		return NewError(nil, err)
	}

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

	rs, err := c.httpClient.Do(rq)
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
		return c.retry(ctx, cRoute, rqBody, rsBody, tries+1)

	case http.StatusUnauthorized:
		c.Logger().Error(discord.ErrUnauthorized)
		return NewError(rs, discord.ErrUnauthorized)

	default:
		body, _ := ioutil.ReadAll(rq.Body)
		return NewError(rs, fmt.Errorf("request to %s failed. statuscode: %d, body: %s", rq.URL, rs.StatusCode, body))
	}
}

func (c *HTTPClientImpl) Do(ctx context.Context, cRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) Error {
	return c.retry(ctx, cRoute, rqBody, rsBody, 1)
}
