package bot

import (
	"github.com/DisgoOrg/disgo/cache"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding"
	"github.com/DisgoOrg/disgo/gateway/sharding/srate"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/internal/tokenhelper"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
	"github.com/pkg/errors"
)

// Config lets you configure your Client instance
// Config is the core.Client config used to configure everything
type Config struct {
	Logger log.Logger

	RestClient           rest.Client
	RestClientConfigOpts []rest.ConfigOpt
	RestServices         rest.Services

	EventManager           EventManager
	EventManagerConfigOpts []EventManagerConfigOpt

	Gateway           gateway.Gateway
	GatewayConfigOpts []gateway.ConfigOpt

	ShardManager           sharding.ShardManager
	ShardManagerConfigOpts []sharding.ConfigOpt

	HTTPServer           httpserver.Server
	HTTPServerConfigOpts []httpserver.ConfigOpt

	Caches          cache.Caches
	CacheConfigOpts []cache.ConfigOpt

	AudioController       AudioController
	MemberChunkingManager MemberChunkingManager
	MemberChunkingFilter  *MemberChunkingFilter
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
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestServices(restServices rest.Services) ConfigOpt {
	return func(config *Config) {
		config.RestServices = restServices
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEventManager(eventManager EventManager) ConfigOpt {
	return func(config *Config) {
		config.EventManager = eventManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEventManagerConfigOpts(opts ...EventManagerConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardManager(shardManager sharding.ShardManager) ConfigOpt {
	return func(config *Config) {
		config.ShardManager = shardManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardManagerConfigOpts(opts ...sharding.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServer(httpServer httpserver.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.HTTPServerConfigOpts = append(config.HTTPServerConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCaches(caches cache.Caches) ConfigOpt {
	return func(config *Config) {
		config.Caches = caches
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCacheConfigOpts(opts ...cache.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.CacheConfigOpts = append(config.CacheConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithAudioController(audioController AudioController) ConfigOpt {
	return func(config *Config) {
		config.AudioController = audioController
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMemberChunkingManager(memberChunkingManager MemberChunkingManager) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingManager = memberChunkingManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMemberChunkingFilter(memberChunkingFilter MemberChunkingFilter) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingFilter = &memberChunkingFilter
	}
}

func BuildClient(token string, config Config, gatewayEventHandlerFunc func(client Client) gateway.EventHandlerFunc, httpServerEventHandlerFunc func(client Client) httpserver.EventHandlerFunc) (Client, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := tokenhelper.IDFromToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting application id from BotToken")
	}
	client := &ClientImpl{
		BotToken: token,
	}

	// TODO: figure out how we handle different application & client ids
	client.BotApplicationID = *id
	client.BotClientID = *id

	if config.Logger == nil {
		config.Logger = log.Default()
	}
	client.BotLogger = config.Logger

	if config.RestClient == nil {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		if config.RestClientConfig.Logger == nil {
			config.RestClientConfig.Logger = config.Logger
		}
		config.RestClient = rest.NewClient(client.BotToken, config.RestClientConfigOpts)
	}

	if config.RestServices == nil {
		config.RestServices = rest.NewServices(config.RestClient)
	}
	client.RestServices = config.RestServices

	if config.EventManager == nil {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &DefaultEventManagerConfig
		}

		config.EventManager = NewEventManager(client, config.EventManagerConfig)
	}
	client.BotEventManager = config.EventManager

	if config.Gateway == nil && config.GatewayConfig != nil {
		var gatewayRs *discord.Gateway
		gatewayRs, err = client.RestServices.GatewayService().GetGateway()
		if err != nil {
			return nil, err
		}
		if config.GatewayConfig.Logger == nil {
			config.GatewayConfig.Logger = client.BotLogger
		}
		config.Gateway = gateway.New(token, gatewayRs.URL, 0, 0, config.GatewayConfig)
	}
	client.BotGateway = config.Gateway

	if config.ShardManager == nil && config.ShardManagerConfig != nil {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = client.RestServices.GatewayService().GetGatewayBot()
		if err != nil {
			return nil, err
		}

		if config.ShardManagerConfig.RateLimiterConfig == nil {
			config.ShardManagerConfig.RateLimiterConfig = &srate.DefaultConfig
		}
		if config.ShardManagerConfig.RateLimiterConfig.Logger == nil {
			config.ShardManagerConfig.RateLimiterConfig.Logger = config.Logger
		}
		if config.ShardManagerConfig.RateLimiterConfig.MaxConcurrency == 0 {
			config.ShardManagerConfig.RateLimiterConfig.MaxConcurrency = gatewayBotRs.SessionStartLimit.MaxConcurrency
		}

		// apply recommended shard count
		if !config.ShardManagerConfig.CustomShards {
			config.ShardManagerConfig.ShardCount = gatewayBotRs.Shards
			config.ShardManagerConfig.Shards = sharding.NewIntSet()
			for i := 0; i < gatewayBotRs.Shards; i++ {
				config.ShardManagerConfig.Shards.Add(i)
			}
		}
		if config.ShardManager == nil {
			config.ShardManager = sharding.New(token, gatewayBotRs.URL, config.ShardManagerConfig)
		}
	}
	client.BotShardManager = config.ShardManager

	if config.HTTPServer == nil && config.HTTPServerConfig != nil {
		if config.HTTPServerConfig.Logger == nil {
			config.HTTPServerConfig.Logger = config.Logger
		}
		config.HTTPServer = httpserver.New(config.HTTPServerConfig)
	}
	client.BotHTTPServer = config.HTTPServer

	if config.AudioController == nil {
		config.AudioController = NewAudioController(client)
	}
	client.BotAudioController = config.AudioController

	if config.MemberChunkingManager == nil {
		if config.MemberChunkingFilter == nil {
			config.MemberChunkingFilter = &MemberChunkingFilterNone
		}
		config.MemberChunkingManager = NewMemberChunkingManager(client, *config.MemberChunkingFilter)
	}
	client.BotMemberChunkingManager = config.MemberChunkingManager

	if config.Caches == nil {
		if config.CacheConfig == nil {
			config.CacheConfig = &cache.DefaultConfig
		}
		config.Caches = cache.NewCaches(*config.CacheConfig)
	}
	client.BotCaches = config.Caches

	return client, nil
}
