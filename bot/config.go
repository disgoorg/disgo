package bot

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/internal/tokenhelper"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/disgo/voice"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(gatewayHandlers map[gateway.EventType]GatewayEventHandler, httpHandler HTTPServerEventHandler) *Config {
	return &Config{
		Logger:                 slog.Default(),
		EventManagerConfigOpts: []EventManagerConfigOpt{WithGatewayHandlers(gatewayHandlers), WithHTTPServerHandler(httpHandler)},
		MemberChunkingFilter:   MemberChunkingFilterNone,
	}
}

// Config lets you configure your Client instance.
type Config struct {
	Logger *slog.Logger

	RestClient           rest.Client
	RestClientConfigOpts []rest.ConfigOpt
	Rest                 rest.Rest

	EventManager           EventManager
	EventManagerConfigOpts []EventManagerConfigOpt

	VoiceManager           voice.Manager
	VoiceManagerConfigOpts []voice.ManagerConfigOpt

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

// WithLogger lets you inject your own logger implementing *slog.Logger.
func WithLogger(logger *slog.Logger) ConfigOpt {
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

// WithEventListenerFunc adds the given func(e E) to the default EventManager.
func WithEventListenerFunc[E Event](f func(e E)) ConfigOpt {
	return WithEventListeners(NewListenerFunc(f))
}

// WithEventListenerChan adds the given chan<- E to the default EventManager.
func WithEventListenerChan[E Event](c chan<- E) ConfigOpt {
	return WithEventListeners(NewListenerChan(c))
}

// WithGateway lets you inject your own gateway.Gateway.
func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

// WithDefaultGateway creates a gateway.Gateway with sensible defaults.
func WithDefaultGateway() ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, func(_ *gateway.Config) {})
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

// WithDefaultShardManager creates a sharding.ShardManager with sensible defaults.
func WithDefaultShardManager() ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, func(_ *sharding.Config) {})
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
func BuildClient(token string, cfg *Config, gatewayEventHandlerFunc func(client Client) gateway.EventHandlerFunc, httpServerEventHandlerFunc func(client Client) httpserver.EventHandlerFunc, os string, name string, github string, version string) (Client, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := tokenhelper.IDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("error while getting application id from token: %w", err)
	}
	client := &clientImpl{
		token:  token,
		logger: cfg.Logger,
	}

	client.applicationID = *id

	if cfg.RestClient == nil {
		// prepend standard user-agent. this can be overridden as it's appended to the front of the slice
		cfg.RestClientConfigOpts = append([]rest.ConfigOpt{
			rest.WithUserAgent(fmt.Sprintf("DiscordBot (%s, %s)", github, version)),
			rest.WithLogger(client.logger),
			func(config *rest.Config) {
				config.RateLimiterConfigOpts = append([]rest.RateLimiterConfigOpt{rest.WithRateLimiterLogger(cfg.Logger)}, config.RateLimiterConfigOpts...)
			},
		}, cfg.RestClientConfigOpts...)

		cfg.RestClient = rest.NewClient(client.token, cfg.RestClientConfigOpts...)
	}

	if cfg.Rest == nil {
		cfg.Rest = rest.New(cfg.RestClient)
	}
	client.restServices = cfg.Rest

	if cfg.VoiceManager == nil {
		cfg.VoiceManager = voice.NewManager(client.UpdateVoiceState, *id, append([]voice.ManagerConfigOpt{voice.WithLogger(cfg.Logger)}, cfg.VoiceManagerConfigOpts...)...)
	}
	client.voiceManager = cfg.VoiceManager

	if cfg.EventManager == nil {
		cfg.EventManager = NewEventManager(client, append([]EventManagerConfigOpt{WithEventManagerLogger(cfg.Logger)}, cfg.EventManagerConfigOpts...)...)
	}
	client.eventManager = cfg.EventManager

	if cfg.Gateway == nil && len(cfg.GatewayConfigOpts) > 0 {
		var gatewayRs *discord.Gateway
		gatewayRs, err = client.restServices.GetGateway()
		if err != nil {
			return nil, err
		}

		cfg.GatewayConfigOpts = append([]gateway.ConfigOpt{
			gateway.WithURL(gatewayRs.URL),
			gateway.WithLogger(cfg.Logger),
			gateway.WithOS(os),
			gateway.WithBrowser(name),
			gateway.WithDevice(name),
			func(config *gateway.Config) {
				config.RateLimiterConfigOpts = append([]gateway.RateLimiterConfigOpt{gateway.WithRateLimiterLogger(cfg.Logger)}, config.RateLimiterConfigOpts...)
			},
		}, cfg.GatewayConfigOpts...)

		cfg.Gateway = gateway.New(token, gatewayEventHandlerFunc(client), nil, cfg.GatewayConfigOpts...)
	}
	client.gateway = cfg.Gateway

	if cfg.ShardManager == nil && len(cfg.ShardManagerConfigOpts) > 0 {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = client.restServices.GetGatewayBot()
		if err != nil {
			return nil, err
		}

		shardIDs := make([]int, gatewayBotRs.Shards)
		for i := 0; i < gatewayBotRs.Shards-1; i++ {
			shardIDs[i] = i
		}

		cfg.ShardManagerConfigOpts = append([]sharding.ConfigOpt{
			sharding.WithShardCount(gatewayBotRs.Shards),
			sharding.WithShardIDs(shardIDs...),
			sharding.WithGatewayConfigOpts(
				gateway.WithURL(gatewayBotRs.URL),
				gateway.WithLogger(cfg.Logger),
				gateway.WithOS(os),
				gateway.WithBrowser(name),
				gateway.WithDevice(name),
				func(config *gateway.Config) {
					config.RateLimiterConfigOpts = append([]gateway.RateLimiterConfigOpt{gateway.WithRateLimiterLogger(cfg.Logger)}, config.RateLimiterConfigOpts...)
				},
			),
			sharding.WithLogger(cfg.Logger),
			func(config *sharding.Config) {
				config.RateLimiterConfigOpts = append([]sharding.RateLimiterConfigOpt{sharding.WithRateLimiterLogger(cfg.Logger), sharding.WithMaxConcurrency(gatewayBotRs.SessionStartLimit.MaxConcurrency)}, config.RateLimiterConfigOpts...)
			},
		}, cfg.ShardManagerConfigOpts...)

		cfg.ShardManager = sharding.New(token, gatewayEventHandlerFunc(client), cfg.ShardManagerConfigOpts...)
	}
	client.shardManager = cfg.ShardManager

	if cfg.HTTPServer == nil && cfg.PublicKey != "" {
		cfg.HTTPServerConfigOpts = append([]httpserver.ConfigOpt{
			httpserver.WithLogger(cfg.Logger),
		}, cfg.HTTPServerConfigOpts...)

		cfg.HTTPServer = httpserver.New(cfg.PublicKey, httpServerEventHandlerFunc(client), cfg.HTTPServerConfigOpts...)
	}
	client.httpServer = cfg.HTTPServer

	if cfg.MemberChunkingManager == nil {
		cfg.MemberChunkingManager = NewMemberChunkingManager(client, cfg.Logger, cfg.MemberChunkingFilter)
	}
	client.memberChunkingManager = cfg.MemberChunkingManager

	if cfg.Caches == nil {
		cfg.Caches = cache.New(cfg.CacheConfigOpts...)
	}
	client.caches = cfg.Caches

	return client, nil
}
