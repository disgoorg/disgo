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
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// HTTPClient allows doing requests to different endpoints
type HTTPClient interface {
	Close()
	Logger() log.Logger
	HTTPClient() *http.Client
	UserAgent() string
	Do(route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) Error
	DoWithHeaders(route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, customHeader http.Header) Error
}

// NewHTTPClient constructs a new HTTPClient with the given http.Client, log.Logger & useragent
//goland:noinspection GoUnusedExportedFunction
func NewHTTPClient(logger log.Logger, httpClient *http.Client, userAgent string) HTTPClient {
	return &HTTPClientImpl{userAgent: userAgent, httpClient: httpClient, logger: logger}
}

type HTTPClientImpl struct {
	httpClient *http.Client
	logger     log.Logger
	userAgent  string
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

func (c *HTTPClientImpl) UserAgent() string {
	return c.userAgent
}

func (c *HTTPClientImpl) Do(route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) Error {
	return c.DoWithHeaders(route, rqBody, rsBody, nil)
}

func (c *HTTPClientImpl) DoWithHeaders(compiledRoute *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, customHeader http.Header) Error {
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
		c.Logger().Debugf("request to %s, body: %s", compiledRoute.URL(), string(body))
	}

	rq, err := http.NewRequest(compiledRoute.Method().String(), compiledRoute.URL(), rqBuffer)
	if err != nil {
		return NewError(nil, err)
	}

	if customHeader != nil {
		rq.Header = customHeader
	}
	rq.Header.Set("User-Agent", c.UserAgent())
	if contentType != "" {
		rq.Header.Set("Content-Type", contentType)
	}

	rs, err := c.httpClient.Do(rq)
	if err != nil {
		return NewError(rs, err)
	}

	if rs.Body != nil {
		buffer := &bytes.Buffer{}
		body, _ := ioutil.ReadAll(io.TeeReader(rs.Body, buffer))
		rs.Body = ioutil.NopCloser(buffer)
		c.Logger().Debugf("response from %s, code %d, body: %s", compiledRoute.URL(), rs.StatusCode, string(body))
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

	case http.StatusTooManyRequests:
		c.Logger().Error(discord.ErrRatelimited)
		return NewError(rs, discord.ErrRatelimited)

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
