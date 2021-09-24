package core

import (
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/sharding"
	"github.com/DisgoOrg/log"
)

type BotConfig struct {
	Logger log.Logger

	HTTPClient       *http.Client
	RestClient       rest.Client
	RestClientConfig *rest.Config
	RestServices     rest.Services

	EventManager             EventManager
	EventListeners           []EventListener
	RawEventsEnabled         bool
	VoiceDispatchInterceptor VoiceDispatchInterceptor

	Gateway       gateway.Gateway
	GatewayConfig *gateway.Config

	ShardManager       sharding.ShardManager
	ShardManagerConfig *sharding.Config

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

func WithShardManager(shardManager sharding.ShardManager) BotConfigOpt {
	return func(config *BotConfig) {
		config.ShardManager = shardManager
	}
}

func WithShardManagerConfig(shardManagerConfig sharding.Config) BotConfigOpt {
	return func(config *BotConfig) {
		config.ShardManagerConfig = &shardManagerConfig
	}
}

func WithShardManagerConfigOpts(opts ...sharding.ConfigOpt) BotConfigOpt {
	return func(config *BotConfig) {
		if config.ShardManagerConfig == nil {
			config.ShardManagerConfig = &sharding.DefaultConfig
		}
		config.ShardManagerConfig.Apply(opts)
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
			config.CacheConfig = &DefaultCacheConfig
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
