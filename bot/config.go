package bot

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/sharding"
	"github.com/DisgoOrg/log"
)

type Config struct {
	Logger log.Logger

	RestClient       rest.Client
	RestClientConfig *rest.Config
	RestServices     rest.Services

	EventManager       core.EventManager
	EventManagerConfig *core.EventManagerConfig

	Gateway       gateway.Gateway
	GatewayConfig *gateway.Config

	ShardManager       sharding.ShardManager
	ShardManagerConfig *sharding.Config

	HTTPServer       httpserver.Server
	HTTPServerConfig *httpserver.Config

	Caches      core.Caches
	CacheConfig *core.CacheConfig

	AudioController       core.AudioController
	EntityBuilder         core.EntityBuilder
	MemberChunkingManager core.MemberChunkingManager
	MemberChunkingFilter  *core.MemberChunkingFilter
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

func WithRestClientConfig(restClientConfig rest.Config) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfig = &restClientConfig
	}
}

func WithRestClientOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Apply(opts)
	}
}

func WithRestServices(restServices rest.Services) ConfigOpt {
	return func(config *Config) {
		config.RestServices = restServices
	}
}

func WithEventManager(eventManager core.EventManager) ConfigOpt {
	return func(config *Config) {
		config.EventManager = eventManager
	}
}

func WithEventListeners(listeners ...core.EventListener) ConfigOpt {
	return func(config *Config) {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &core.DefaultEventManagerConfig
		}
		config.EventManagerConfig.EventListeners = append(config.EventManagerConfig.EventListeners, listeners...)
	}
}

func WithRawEventsEnabled() ConfigOpt {
	return func(config *Config) {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &core.DefaultEventManagerConfig
		}
		config.EventManagerConfig.RawEventsEnabled = true
	}
}

func WithVoiceDispatchInterceptor(voiceDispatchInterceptor core.VoiceDispatchInterceptor) ConfigOpt {
	return func(config *Config) {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &core.DefaultEventManagerConfig
		}
		config.EventManagerConfig.VoiceDispatchInterceptor = voiceDispatchInterceptor
	}
}

func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

func WithGatewayConfig(gatewayConfig gateway.Config) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfig = &gatewayConfig
	}
}

func WithGatewayOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.GatewayConfig == nil {
			config.GatewayConfig = &gateway.DefaultConfig
		}
		config.GatewayConfig.Apply(opts)
	}
}

func WithShardManager(shardManager sharding.ShardManager) ConfigOpt {
	return func(config *Config) {
		config.ShardManager = shardManager
	}
}

func WithShardManagerConfig(shardManagerConfig sharding.Config) ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfig = &shardManagerConfig
	}
}

func WithShardManagerOpts(opts ...sharding.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.ShardManagerConfig == nil {
			config.ShardManagerConfig = &sharding.DefaultConfig
		}
		config.ShardManagerConfig.Apply(opts)
	}
}

func WithHTTPServer(httpServer httpserver.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

func WithHTTPServerConfig(httpServerConfig httpserver.Config) ConfigOpt {
	return func(config *Config) {
		config.HTTPServerConfig = &httpServerConfig
	}
}

func WithHTTPServerOpts(opts ...httpserver.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.HTTPServerConfig == nil {
			config.HTTPServerConfig = &httpserver.DefaultConfig
		}
		config.HTTPServerConfig.Apply(opts)
	}
}

func WithCaches(caches core.Caches) ConfigOpt {
	return func(config *Config) {
		config.Caches = caches
	}
}

func WithCacheConfig(cacheConfig core.CacheConfig) ConfigOpt {
	return func(config *Config) {
		config.CacheConfig = &cacheConfig
	}
}

func WithCacheOpts(opts ...core.CacheConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.CacheConfig == nil {
			config.CacheConfig = &core.DefaultCacheConfig
		}
		config.CacheConfig.Apply(opts)
	}
}

func WithEntityBuilder(entityBuilder core.EntityBuilder) ConfigOpt {
	return func(config *Config) {
		config.EntityBuilder = entityBuilder
	}
}

func WithAudioController(audioController core.AudioController) ConfigOpt {
	return func(config *Config) {
		config.AudioController = audioController
	}
}

func WithMemberChunkingManager(memberChunkingManager core.MemberChunkingManager) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingManager = memberChunkingManager
	}
}

func WithMemberChunkingFilter(memberChunkingFilter core.MemberChunkingFilter) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingFilter = &memberChunkingFilter
	}
}
