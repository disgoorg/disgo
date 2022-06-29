package bot

import (
	"fmt"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/internal/tokenhelper"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/log"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(gatewayHandlers map[discord.GatewayEventType]GatewayEventHandler, httpHandler HTTPServerEventHandler) *Config {
	return &Config{
		Logger:                 log.Default(),
		EventManagerConfigOpts: []EventManagerConfigOpt{WithGatewayHandlers(gatewayHandlers), WithHTTPServerHandler(httpHandler)},
		MemberChunkingFilter:   MemberChunkingFilterNone,
	}
}

// Config lets you configure your Client instance.
type Config struct {
	Logger log.Logger

	RestClient           rest.Client
	RestClientConfigOpts []rest.ConfigOpt
	Rest                 rest.Rest

	EventManager           EventManager
	EventManagerConfigOpts []EventManagerConfigOpt

	Gateway           gateway.Gateway
	GatewayConfigOpts []gateway.ConfigOpt

	ShardManager           sharding.ShardManager
	ShardManagerConfigOpts []sharding.ConfigOpt

	HTTPServer           httpserver.Server
	PublicKey            string
	HTTPServerConfigOpts []httpserver.ConfigOpt

	Caches          cache.Caches
	CacheConfigOpts []cache.ConfigOpt

	MemberChunkingManager MemberChunkingManager
	MemberChunkingFilter  MemberChunkingFilter
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Client.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger lets you inject your own logger implementing log.Logger.
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithRestClient lets you inject your own rest.Client.
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfigOpts let's you configure the default rest.Client.
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

// WithRest lets you inject your own rest.Rest.
func WithRest(rest rest.Rest) ConfigOpt {
	return func(config *Config) {
		config.Rest = rest
	}
}

// WithEventManager lets you inject your own EventManager.
func WithEventManager(eventManager EventManager) ConfigOpt {
	return func(config *Config) {
		config.EventManager = eventManager
	}
}

// WithEventManagerConfigOpts lets you configure the default EventManager.
func WithEventManagerConfigOpts(opts ...EventManagerConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, opts...)
	}
}

// WithEventListeners adds the given EventListener(s) to the default EventManager.
func WithEventListeners(eventListeners ...EventListener) ConfigOpt {
	return func(config *Config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, WithListeners(eventListeners...))
	}
}

// WithEventListenerFunc adds the given ListenerFunc(s) to the default EventManager.
func WithEventListenerFunc[E Event](listenerFunc func(e E)) ConfigOpt {
	return WithEventListeners(NewListenerFunc(listenerFunc))
}

// WithGateway lets you inject your own gateway.Gateway.
func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

// WithGatewayConfigOpts lets you configure the default gateway.Gateway.
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

// WithShardManager lets you inject your own sharding.ShardManager.
func WithShardManager(shardManager sharding.ShardManager) ConfigOpt {
	return func(config *Config) {
		config.ShardManager = shardManager
	}
}

// WithShardManagerConfigOpts lets you configure the default sharding.ShardManager.
func WithShardManagerConfigOpts(opts ...sharding.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, opts...)
	}
}

// WithHTTPServer lets you inject your own httpserver.Server.
func WithHTTPServer(httpServer httpserver.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

// WithHTTPServerConfigOpts lets you configure the default httpserver.Server.
func WithHTTPServerConfigOpts(publicKey string, opts ...httpserver.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.PublicKey = publicKey
		config.HTTPServerConfigOpts = append(config.HTTPServerConfigOpts, opts...)
	}
}

// WithCaches lets you inject your own cache.Caches.
func WithCaches(caches cache.Caches) ConfigOpt {
	return func(config *Config) {
		config.Caches = caches
	}
}

// WithCacheConfigOpts lets you configure the default cache.Caches.
func WithCacheConfigOpts(opts ...cache.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.CacheConfigOpts = append(config.CacheConfigOpts, opts...)
	}
}

// WithMemberChunkingManager lets you inject your own MemberChunkingManager.
func WithMemberChunkingManager(memberChunkingManager MemberChunkingManager) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingManager = memberChunkingManager
	}
}

// WithMemberChunkingFilter lets you configure the default MemberChunkingFilter.
func WithMemberChunkingFilter(memberChunkingFilter MemberChunkingFilter) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingFilter = memberChunkingFilter
	}
}

// BuildClient creates a new Client instance with the given token, Config, gateway handlers, http handlers os, name, github & version.
func BuildClient(token string, config Config, gatewayEventHandlerFunc func(client Client) gateway.EventHandlerFunc, httpServerEventHandlerFunc func(client Client) httpserver.EventHandlerFunc, os string, name string, github string, version string) (Client, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := tokenhelper.IDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("error while getting application id from token: %w", err)
	}
	client := &clientImpl{
		token:  token,
		logger: config.Logger,
	}

	client.applicationID = *id

	if config.RestClient == nil {
		// prepend standard user-agent. this can be overridden as it's appended to the front of the slice
		config.RestClientConfigOpts = append([]rest.ConfigOpt{
			rest.WithUserAgent(fmt.Sprintf("DiscordBot (%s, %s)", github, version)),
			rest.WithLogger(client.logger),
			func(config *rest.Config) {
				config.RateRateLimiterConfigOpts = append([]rest.RateLimiterConfigOpt{rest.WithRateLimiterLogger(client.logger)}, config.RateRateLimiterConfigOpts...)
			},
		}, config.RestClientConfigOpts...)

		config.RestClient = rest.NewClient(client.token, config.RestClientConfigOpts...)
	}

	if config.Rest == nil {
		config.Rest = rest.New(config.RestClient)
	}
	client.restServices = config.Rest

	if config.EventManager == nil {
		config.EventManager = NewEventManager(client, config.EventManagerConfigOpts...)
	}
	client.eventManager = config.EventManager

	if config.Gateway == nil && config.GatewayConfigOpts != nil {
		var gatewayRs *discord.Gateway
		gatewayRs, err = client.restServices.GetGateway()
		if err != nil {
			return nil, err
		}

		config.GatewayConfigOpts = append([]gateway.ConfigOpt{
			gateway.WithGatewayURL(gatewayRs.URL),
			gateway.WithLogger(client.logger),
			gateway.WithOS(os),
			gateway.WithBrowser(name),
			gateway.WithDevice(name),
			func(config *gateway.Config) {
				config.RateRateLimiterConfigOpts = append([]gateway.RateLimiterConfigOpt{gateway.WithRateLimiterLogger(client.logger)}, config.RateRateLimiterConfigOpts...)
			},
		}, config.GatewayConfigOpts...)

		config.Gateway = gateway.New(token, gatewayEventHandlerFunc(client), nil, config.GatewayConfigOpts...)
	}
	client.gateway = config.Gateway

	if config.ShardManager == nil && config.ShardManagerConfigOpts != nil {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = client.restServices.GetGatewayBot()
		if err != nil {
			return nil, err
		}

		shardIDs := make([]int, gatewayBotRs.Shards)
		for i := 0; i < gatewayBotRs.Shards-1; i++ {
			shardIDs[i] = i
		}

		config.ShardManagerConfigOpts = append([]sharding.ConfigOpt{
			sharding.WithShardCount(gatewayBotRs.Shards),
			sharding.WithShardIDs(shardIDs...),
			sharding.WithGatewayConfigOpts(
				gateway.WithGatewayURL(gatewayBotRs.URL),
				gateway.WithLogger(client.logger),
				gateway.WithOS(os),
				gateway.WithBrowser(name),
				gateway.WithDevice(name),
				func(config *gateway.Config) {
					config.RateRateLimiterConfigOpts = append([]gateway.RateLimiterConfigOpt{gateway.WithRateLimiterLogger(client.logger)}, config.RateRateLimiterConfigOpts...)
				},
			),
			sharding.WithLogger(client.logger),
			func(config *sharding.Config) {
				config.RateRateLimiterConfigOpts = append([]sharding.RateLimiterConfigOpt{sharding.WithRateLimiterLogger(client.logger), sharding.WithMaxConcurrency(gatewayBotRs.SessionStartLimit.MaxConcurrency)}, config.RateRateLimiterConfigOpts...)
			},
		}, config.ShardManagerConfigOpts...)

		config.ShardManager = sharding.New(token, gatewayEventHandlerFunc(client), config.ShardManagerConfigOpts...)
	}
	client.shardManager = config.ShardManager

	if config.HTTPServer == nil && config.PublicKey != "" && config.HTTPServerConfigOpts != nil {
		config.HTTPServerConfigOpts = append([]httpserver.ConfigOpt{
			httpserver.WithLogger(client.logger),
		}, config.HTTPServerConfigOpts...)

		config.HTTPServer = httpserver.New(config.PublicKey, httpServerEventHandlerFunc(client), config.HTTPServerConfigOpts...)
	}
	client.httpServer = config.HTTPServer

	if config.MemberChunkingManager == nil {
		config.MemberChunkingManager = NewMemberChunkingManager(client, config.MemberChunkingFilter)
	}
	client.memberChunkingManager = config.MemberChunkingManager

	if config.Caches == nil {
		config.Caches = cache.New(config.CacheConfigOpts...)
	}
	client.caches = config.Caches

	return client, nil
}
