package core

import (
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/log"
)

type DisgoConfig struct {
	Token  string
	Logger log.Logger

	HTTPClient *http.Client

	RestClient       rest.Client
	RestClientConfig *rest.Config

	RateLimiter       rate.Limiter
	RateLimiterConfig *rate.Config

	RestServices rest.Services

	EventManager             EventManager
	EventListeners           []EventListener
	RawEventsEnabled         bool
	VoiceDispatchInterceptor VoiceDispatchInterceptor

	Gateway       gateway.Gateway
	GatewayConfig *gateway.Config

	HTTPServer       httpserver.Server
	HTTPServerConfig *httpserver.Config

	Cache       Cache
	CacheConfig *CacheConfig

	AudioController AudioController
	EntityBuilder   EntityBuilder
}

type DisgoOpt func(config *DisgoConfig)

func (c *DisgoConfig) Apply(opts []DisgoOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) DisgoOpt {
	return func(config *DisgoConfig) {
		config.Logger = logger
	}
}

func WithHTTPClient(httpClient *http.Client) DisgoOpt {
	return func(config *DisgoConfig) {
		config.HTTPClient = httpClient
	}
}
