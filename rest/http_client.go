package rest

import (
	"context"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
	"net/http"
)

// HTTPClient allows doing requests to different endpoints
type HTTPClient interface {
	Close()
	Logger() log.Logger
	HTTPClient() *http.Client
	UserAgent() string
	Do(ctx context.Context, route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) Error
	DoWithHeaders(ctx context.Context, route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, customHeader http.Header) Error
}

// NewHTTPClient constructs a new HTTPClient with the given http.Client, log.Logger & useragent
//goland:noinspection GoUnusedExportedFunction
func NewHTTPClient(logger log.Logger, httpClient *http.Client, userAgent string) HTTPClient {
	return &HTTPClientImpl{userAgent: userAgent, httpClient: httpClient, logger: logger}
}

type HTTPClientImpl struct {
	httpClient  *http.Client
	rateLimiter rate.RateLimiter
	logger      log.Logger
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

func (c *HTTPClientImpl) UserAgent() string {
	return c.userAgent
}

func (c *HTTPClientImpl) Do(request Request) Error {
	c.rateLimiter.
}
