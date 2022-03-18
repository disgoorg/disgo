package bot

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

// Config lets you configure your Bot instance
// Config is the core.Bot config used to configure everything
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
	MemberChunkingManager core.MemberChunkingManager
	MemberChunkingFilter  *core.MemberChunkingFilter
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger lets you inject your own logger implementing log.Logger
//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfig(restClientConfig rest.Config) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfig = &restClientConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClientOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestServices(restServices rest.Services) ConfigOpt {
	return func(config *Config) {
		config.RestServices = restServices
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEventManager(eventManager core.EventManager) ConfigOpt {
	return func(config *Config) {
		config.EventManager = eventManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEventListeners(listeners ...core.EventListener) ConfigOpt {
	return func(config *Config) {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &core.DefaultEventManagerConfig
		}
		config.EventManagerConfig.EventListeners = append(config.EventManagerConfig.EventListeners, listeners...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRawEventsEnabled() ConfigOpt {
	return func(config *Config) {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &core.DefaultEventManagerConfig
		}
		config.EventManagerConfig.RawEventsEnabled = true
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithAsyncEventsEnabled() ConfigOpt {
	return func(config *Config) {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &core.DefaultEventManagerConfig
		}
		config.EventManagerConfig.AsyncEventsEnabled = true
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayConfig(gatewayConfig gateway.Config) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfig = &gatewayConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.GatewayConfig == nil {
			config.GatewayConfig = &gateway.DefaultConfig
		}
		config.GatewayConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardManager(shardManager sharding.ShardManager) ConfigOpt {
	return func(config *Config) {
		config.ShardManager = shardManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardManagerConfig(shardManagerConfig sharding.Config) ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfig = &shardManagerConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardManagerOpts(opts ...sharding.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.ShardManagerConfig == nil {
			config.ShardManagerConfig = &sharding.DefaultConfig
		}
		config.ShardManagerConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServer(httpServer httpserver.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServerConfig(httpServerConfig httpserver.Config) ConfigOpt {
	return func(config *Config) {
		config.HTTPServerConfig = &httpServerConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServerOpts(opts ...httpserver.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.HTTPServerConfig == nil {
			config.HTTPServerConfig = &httpserver.DefaultConfig
		}
		config.HTTPServerConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCaches(caches core.Caches) ConfigOpt {
	return func(config *Config) {
		config.Caches = caches
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCacheConfig(cacheConfig core.CacheConfig) ConfigOpt {
	return func(config *Config) {
		config.CacheConfig = &cacheConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCacheOpts(opts ...core.CacheConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.CacheConfig == nil {
			config.CacheConfig = &core.DefaultCacheConfig
		}
		config.CacheConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithAudioController(audioController core.AudioController) ConfigOpt {
	return func(config *Config) {
		config.AudioController = audioController
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMemberChunkingManager(memberChunkingManager core.MemberChunkingManager) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingManager = memberChunkingManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMemberChunkingFilter(memberChunkingFilter core.MemberChunkingFilter) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingFilter = &memberChunkingFilter
	}
}
