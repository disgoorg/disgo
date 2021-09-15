package core

import (
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/disgo/sharding"
	"github.com/DisgoOrg/log"
)

type BotConfig struct {
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

	ShardManager       sharding.ShardManager
	ShardManagerConfig *sharding.Config
	GatewayConfig      *gateway.Config
	GatewayFunc        func(config gateway.Config) gateway.Gateway

	HTTPServer       httpserver.Server
	HTTPServerConfig *httpserver.Config

	Caches      Caches
	CacheConfig *CacheConfig

	AudioController AudioController
	EntityBuilder   EntityBuilder
}

type BotConfigOpt func(config *BotConfig)

func (c *BotConfig) Apply(opts []BotConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) BotConfigOpt {
	return func(config *BotConfig) {
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

func WithHTTPClient(httpClient *http.Client) BotConfigOpt {
	return func(config *BotConfig) {
		config.HTTPClient = httpClient
	}
}

func WithRestClient(restClient rest.Client) BotConfigOpt {
	return func(config *BotConfig) {
		config.RestClient = restClient
	}
}

func WithRestClientConfig(restClientConfig rest.Config) BotConfigOpt {
	return func(config *BotConfig) {
		config.RestClientConfig = &restClientConfig
	}
}

func WithRestClientConfigOpts(opts ...rest.ConfigOpt) BotConfigOpt {
	return func(config *BotConfig) {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Apply(opts)
	}
}

func WithRateLimiter(rateLimiter rate.Limiter) BotConfigOpt {
	return func(config *BotConfig) {
		config.RateLimiter = rateLimiter
	}
}

func WithRateLimiterConfig(rateLimiterConfig rate.Config) BotConfigOpt {
	return func(config *BotConfig) {
		config.RateLimiterConfig = &rateLimiterConfig
	}
}

func WithRateLimiterConfigOpts(opts ...rate.ConfigOpt) BotConfigOpt {
	return func(config *BotConfig) {
		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &rate.DefaultConfig
		}
		config.RateLimiterConfig.Apply(opts)
	}
}

func WithRestServices(restServices rest.Services) BotConfigOpt {
	return func(config *BotConfig) {
		config.RestServices = restServices
	}
}

func WithEventManager(eventManager EventManager) BotConfigOpt {
	return func(config *BotConfig) {
		config.EventManager = eventManager
	}
}

func WithEventListeners(listeners ...EventListener) BotConfigOpt {
	return func(config *BotConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

func WithRawEventsEnabled() BotConfigOpt {
	return func(config *BotConfig) {
		config.RawEventsEnabled = true
	}
}

func WithVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) BotConfigOpt {
	return func(config *BotConfig) {
		config.VoiceDispatchInterceptor = voiceDispatchInterceptor
	}
}

func WithGateway(gateway gateway.Gateway) BotConfigOpt {
	return func(config *BotConfig) {
		config.Gateway = gateway
	}
}

func WithGatewayConfig(gatewayConfig gateway.Config) BotConfigOpt {
	return func(config *BotConfig) {
		config.GatewayConfig = &gatewayConfig
	}
}

func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) BotConfigOpt {
	return func(config *BotConfig) {
		if config.GatewayConfig == nil {
			config.GatewayConfig = &gateway.DefaultConfig
		}
		config.GatewayConfig.Apply(opts)
	}
}

func WithHTTPServer(httpServer httpserver.Server) BotConfigOpt {
	return func(config *BotConfig) {
		config.HTTPServer = httpServer
	}
}

func WithHTTPServerConfig(httpServerConfig httpserver.Config) BotConfigOpt {
	return func(config *BotConfig) {
		config.HTTPServerConfig = &httpServerConfig
	}
}

func WithHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) BotConfigOpt {
	return func(config *BotConfig) {
		if config.HTTPServerConfig == nil {
			config.HTTPServerConfig = &httpserver.DefaultConfig
		}
		config.HTTPServerConfig.Apply(opts)
	}
}

func WithCaches(caches Caches) BotConfigOpt {
	return func(config *BotConfig) {
		config.Caches = caches
	}
}

func WithCacheConfig(cacheConfig CacheConfig) BotConfigOpt {
	return func(config *BotConfig) {
		config.CacheConfig = &cacheConfig
	}
}

func WithCacheConfigOpts(opts ...CacheConfigOpt) BotConfigOpt {
	return func(config *BotConfig) {
		if config.CacheConfig == nil {
			config.CacheConfig = &DefaultConfig
		}
		config.CacheConfig.Apply(opts)
	}
}

func WithAudioController(audioController AudioController) BotConfigOpt {
	return func(config *BotConfig) {
		config.AudioController = audioController
	}
}

func WithEntityBuilder(entityBuilder EntityBuilder) BotConfigOpt {
	return func(config *BotConfig) {
		config.EntityBuilder = entityBuilder
	}
}
