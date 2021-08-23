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

		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Logger = logger

		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &rate.DefaultConfig
		}
		config.RateLimiterConfig.Logger = logger

		if config.GatewayConfig == nil {
			config.GatewayConfig = &gateway.DefaultConfig
		}
		config.GatewayConfig.Logger = logger

		if config.HTTPServerConfig == nil {
			config.HTTPServerConfig = &httpserver.DefaultConfig
		}
		config.HTTPServerConfig.Logger = logger
	}
}

func WithHTTPClient(httpClient *http.Client) DisgoOpt {
	return func(config *DisgoConfig) {
		config.HTTPClient = httpClient
	}
}

func WithRestClient(restClient rest.Client) DisgoOpt {
	return func(config *DisgoConfig) {
		config.RestClient = restClient
	}
}

func WithRestClientConfig(restClientConfig rest.Config) DisgoOpt {
	return func(config *DisgoConfig) {
		config.RestClientConfig = &restClientConfig
	}
}

func WithRestClientConfigOpts(opts ...rest.ConfigOpt) DisgoOpt {
	return func(config *DisgoConfig) {
		config.RestClientConfig.Apply(opts)
	}
}

func WithRateLimiter(rateLimiter rate.Limiter) DisgoOpt {
	return func(config *DisgoConfig) {
		config.RateLimiter = rateLimiter
	}
}

func WithRateLimiterConfig(rateLimiterConfig rate.Config) DisgoOpt {
	return func(config *DisgoConfig) {
		config.RateLimiterConfig = &rateLimiterConfig
	}
}

func WithRateLimiterConfigOpts(opts ...rate.ConfigOpt) DisgoOpt {
	return func(config *DisgoConfig) {
		config.RateLimiterConfig.Apply(opts)
	}
}

func WithRestServices(restServices rest.Services) DisgoOpt {
	return func(config *DisgoConfig) {
		config.RestServices = restServices
	}
}

func WithEventManager(eventManager EventManager) DisgoOpt {
	return func(config *DisgoConfig) {
		config.EventManager = eventManager
	}
}

func WithEventListeners(listeners ...EventListener) DisgoOpt {
	return func(config *DisgoConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

func WithRawEventsEnabled() DisgoOpt {
	return func(config *DisgoConfig) {
		config.RawEventsEnabled = true
	}
}

func WithVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) DisgoOpt {
	return func(config *DisgoConfig) {
		config.VoiceDispatchInterceptor = voiceDispatchInterceptor
	}
}

func WithGateway(gateway gateway.Gateway) DisgoOpt {
	return func(config *DisgoConfig) {
		config.Gateway = gateway
	}
}

func WithGatewayConfig(gatewayConfig gateway.Config) DisgoOpt {
	return func(config *DisgoConfig) {
		config.GatewayConfig = &gatewayConfig
	}
}

func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) DisgoOpt {
	return func(config *DisgoConfig) {
		config.GatewayConfig.Apply(opts)
	}
}

func WithHTTPServer(httpServer httpserver.Server) DisgoOpt {
	return func(config *DisgoConfig) {
		config.HTTPServer = httpServer
	}
}

func WithHTTPServerConfig(httpServerConfig httpserver.Config) DisgoOpt {
	return func(config *DisgoConfig) {
		config.HTTPServerConfig = &httpServerConfig
	}
}

func WithHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) DisgoOpt {
	return func(config *DisgoConfig) {
		config.HTTPServerConfig.Apply(opts)
	}
}

func WithCache(cache Cache) DisgoOpt {
	return func(config *DisgoConfig) {
		config.Cache = cache
	}
}

func WithCacheConfig(cacheConfig CacheConfig) DisgoOpt {
	return func(config *DisgoConfig) {
		config.CacheConfig = &cacheConfig
	}
}

func WithCacheConfigOpts(opts ...CacheConfigOpt) DisgoOpt {
	return func(config *DisgoConfig) {
		config.CacheConfig.Apply(opts)
	}
}

func WithAudioController(audioController AudioController) DisgoOpt {
	return func(config *DisgoConfig) {
		config.AudioController = audioController
	}
}

func WithEntityBuilder(entityBuilder EntityBuilder) DisgoOpt {
	return func(config *DisgoConfig) {
		config.EntityBuilder = entityBuilder
	}
}
