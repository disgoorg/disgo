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
	"github.com/DisgoOrg/log"
)

// HTTPClient allows doing requests to different endpoints
type HTTPClient interface {
	Close()
	Logger() log.Logger
	HTTPClient() *http.Client
	RateLimiter() rate.RateLimiter
	UserAgent() string
	Do(ctx context.Context, route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) Error
}

// NewHTTPClient constructs a new HTTPClient with the given http.Client, log.Logger & useragent
//goland:noinspection GoUnusedExportedFunction
func NewHTTPClient(logger log.Logger, httpClient *http.Client, rateLimiter rate.RateLimiter, headers http.Header, userAgent string) HTTPClient {
	if logger == nil {
		logger = log.Default()
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if rateLimiter == nil {
		rateLimiter = rate.NewRateLimiter(logger)
	}
	return &HTTPClientImpl{
		logger:      logger,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		headers:     headers,
		userAgent:   userAgent,
	}
}

type HTTPClientImpl struct {
	logger      log.Logger
	httpClient  *http.Client
	rateLimiter rate.RateLimiter
	headers     http.Header
	userAgent   string
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

func (c *HTTPClientImpl) UserAgent() string {
	return c.userAgent
}

func (c *HTTPClientImpl) Do(ctx context.Context, cRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) Error {
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

	if c.headers != nil {
		rq.Header = c.headers
	}
	rq.Header.Set("User-Agent", c.UserAgent())
	if contentType != "" {
		rq.Header.Set("Content-Type", contentType)
	}

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return NewError(rs, err)
	}

	// catch rate limit errors here
	// TODO: we should retry on hitting 429
	if err = c.rateLimiter.UnlockBucket(cRoute, rs.Header); err != nil {
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

	case http.StatusUnauthorized:
		c.Logger().Error(discord.ErrUnauthorized)
		return NewError(rs, discord.ErrUnauthorized)

	default:
		body, _ := ioutil.ReadAll(rq.Body)
		return NewError(rs, fmt.Errorf("request to %s failed. statuscode: %d, body: %s", rq.URL, rs.StatusCode, body))
	}
}
