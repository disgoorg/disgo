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

func defaultConfig(gatewayHandlers map[gateway.EventType]GatewayEventHandler, httpHandler HTTPServerEventHandler) config {
	return config{
		Logger:                 slog.Default(),
		EventManagerConfigOpts: []EventManagerConfigOpt{WithGatewayHandlers(gatewayHandlers), WithHTTPServerHandler(httpHandler)},
		MemberChunkingFilter:   MemberChunkingFilterNone,
	}
}

type config struct {
	Logger *slog.Logger

	RestClient           rest.Client
	RestClientConfigOpts []rest.ClientConfigOpt
	Rest                 rest.Rest
	RestConfigOpts       []rest.ConfigOpt

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

// ConfigOpt is a type alias for a function that takes a config and is used to configure your Client.
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "bot"))
}

// WithLogger lets you inject your own Logger implementing *slog.Logger.
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *config) {
		config.Logger = logger
	}
}

// WithRestClient lets you inject your own rest.Client.
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfigOpts let's you configure the default rest.Client.
func WithRestClientConfigOpts(opts ...rest.ClientConfigOpt) ConfigOpt {
	return func(config *config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

// WithRest lets you inject your own rest.Rest.
func WithRest(rest rest.Rest) ConfigOpt {
	return func(config *config) {
		config.Rest = rest
	}
}

// WithRestConfigOpts lets you configure the default rest.Rest.
func WithRestConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *config) {
		config.RestConfigOpts = append(config.RestConfigOpts, opts...)
	}
}

// WithEventManager lets you inject your own EventManager.
func WithEventManager(eventManager EventManager) ConfigOpt {
	return func(config *config) {
		config.EventManager = eventManager
	}
}

// WithEventManagerConfigOpts lets you configure the default EventManager.
func WithEventManagerConfigOpts(opts ...EventManagerConfigOpt) ConfigOpt {
	return func(config *config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, opts...)
	}
}

// WithEventListeners adds the given EventListener(s) to the default EventManager.
func WithEventListeners(eventListeners ...EventListener) ConfigOpt {
	return func(config *config) {
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
	return func(config *config) {
		config.Gateway = gateway
	}
}

// WithDefaultGateway creates a gateway.Gateway with sensible defaults.
func WithDefaultGateway() ConfigOpt {
	return func(config *config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, gateway.WithDefault())
	}
}

// WithGatewayConfigOpts lets you configure the default gateway.Gateway.
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

// WithShardManager lets you inject your own sharding.ShardManager.
func WithShardManager(shardManager sharding.ShardManager) ConfigOpt {
	return func(config *config) {
		config.ShardManager = shardManager
	}
}

// WithDefaultShardManager creates a sharding.ShardManager with sensible defaults.
func WithDefaultShardManager() ConfigOpt {
	return func(config *config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, sharding.WithDefault())
	}
}

// WithShardManagerConfigOpts lets you configure the default sharding.ShardManager.
func WithShardManagerConfigOpts(opts ...sharding.ConfigOpt) ConfigOpt {
	return func(config *config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, opts...)
	}
}

// WithHTTPServer lets you inject your own httpserver.Server.
func WithHTTPServer(httpServer httpserver.Server) ConfigOpt {
	return func(config *config) {
		config.HTTPServer = httpServer
	}
}

// WithHTTPServerConfigOpts lets you configure the default httpserver.Server.
func WithHTTPServerConfigOpts(publicKey string, opts ...httpserver.ConfigOpt) ConfigOpt {
	return func(config *config) {
		config.PublicKey = publicKey
		config.HTTPServerConfigOpts = append(config.HTTPServerConfigOpts, opts...)
	}
}

// WithCaches lets you inject your own cache.Caches.
func WithCaches(caches cache.Caches) ConfigOpt {
	return func(config *config) {
		config.Caches = caches
	}
}

// WithCacheConfigOpts lets you configure the default cache.Caches.
func WithCacheConfigOpts(opts ...cache.ConfigOpt) ConfigOpt {
	return func(config *config) {
		config.CacheConfigOpts = append(config.CacheConfigOpts, opts...)
	}
}

// WithMemberChunkingManager lets you inject your own MemberChunkingManager.
func WithMemberChunkingManager(memberChunkingManager MemberChunkingManager) ConfigOpt {
	return func(config *config) {
		config.MemberChunkingManager = memberChunkingManager
	}
}

// WithMemberChunkingFilter lets you configure the default MemberChunkingFilter.
func WithMemberChunkingFilter(memberChunkingFilter MemberChunkingFilter) ConfigOpt {
	return func(config *config) {
		config.MemberChunkingFilter = memberChunkingFilter
	}
}

func defaultHTTPServerEventHandlerFunc(client *Client) httpserver.EventHandlerFunc {
	return client.EventManager.HandleHTTPEvent
}

func defaultGatewayEventHandlerFunc(client *Client) gateway.EventHandlerFunc {
	return client.EventManager.HandleGatewayEvent
}

// BuildClient creates a new Client instance with the given Token, config, Gateway handlers, http handlers os, name, github & version.
func BuildClient(
	token string,
	otps []ConfigOpt,
	gatewayHandlers map[gateway.EventType]GatewayEventHandler,
	httpHandler HTTPServerEventHandler,
	os string,
	name string,
	github string,
	version string,
) (*Client, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := tokenhelper.IDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("error while getting application id from Token: %w", err)
	}

	cfg := defaultConfig(gatewayHandlers, httpHandler)
	cfg.apply(otps)

	client := &Client{
		Token:         token,
		Logger:        cfg.Logger,
		ApplicationID: *id,
	}

	if cfg.RestClient == nil {
		// prepend standard user-agent. this can be overridden as it's appended to the front of the slice
		cfg.RestClientConfigOpts = append([]rest.ClientConfigOpt{
			rest.WithUserAgent(fmt.Sprintf("DiscordBot (%s, %s)", github, version)),
			rest.WithLogger(client.Logger),
			rest.WithDefaultRateLimiterConfigOpts(
				rest.WithRateLimiterLogger(cfg.Logger),
			),
		}, cfg.RestClientConfigOpts...)

		cfg.RestClient = rest.NewClient(client.Token, cfg.RestClientConfigOpts...)
	}

	if cfg.Rest == nil {
		cfg.Rest = rest.New(cfg.RestClient, cfg.RestConfigOpts...)
	}
	client.Rest = cfg.Rest

	if cfg.VoiceManager == nil {
		cfg.VoiceManager = voice.NewManager(client.UpdateVoiceState, *id, append([]voice.ManagerConfigOpt{voice.WithLogger(cfg.Logger)}, cfg.VoiceManagerConfigOpts...)...)
	}
	client.VoiceManager = cfg.VoiceManager

	if cfg.EventManager == nil {
		cfg.EventManager = NewEventManager(client, append([]EventManagerConfigOpt{WithEventManagerLogger(cfg.Logger)}, cfg.EventManagerConfigOpts...)...)
	}
	client.EventManager = cfg.EventManager

	if cfg.Gateway == nil && len(cfg.GatewayConfigOpts) > 0 {
		var gatewayRs *discord.Gateway
		gatewayRs, err = client.Rest.GetGateway()
		if err != nil {
			return nil, err
		}

		cfg.GatewayConfigOpts = append([]gateway.ConfigOpt{
			gateway.WithURL(gatewayRs.URL),
			gateway.WithLogger(cfg.Logger),
			gateway.WithOS(os),
			gateway.WithBrowser(name),
			gateway.WithDevice(name),
			gateway.WithDefaultRateLimiterConfigOpts(
				gateway.WithRateLimiterLogger(cfg.Logger),
			),
		}, cfg.GatewayConfigOpts...)

		cfg.Gateway = gateway.New(token, defaultGatewayEventHandlerFunc(client), cfg.GatewayConfigOpts...)
	}
	client.Gateway = cfg.Gateway

	if cfg.ShardManager == nil && len(cfg.ShardManagerConfigOpts) > 0 {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = client.Rest.GetGatewayBot()
		if err != nil {
			return nil, err
		}

		shardIDs := make([]int, gatewayBotRs.Shards)
		for i := range gatewayBotRs.Shards {
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
			),
			sharding.WithLogger(cfg.Logger),
			sharding.WithDefaultIdentifyRateLimiterConfigOpt(
				gateway.WithIdentifyMaxConcurrency(gatewayBotRs.SessionStartLimit.MaxConcurrency),
				gateway.WithIdentifyRateLimiterLogger(cfg.Logger),
			),
		}, cfg.ShardManagerConfigOpts...)

		cfg.ShardManager = sharding.New(token, defaultGatewayEventHandlerFunc(client), cfg.ShardManagerConfigOpts...)
	}
	client.ShardManager = cfg.ShardManager

	if cfg.HTTPServer == nil && cfg.PublicKey != "" {
		cfg.HTTPServerConfigOpts = append([]httpserver.ConfigOpt{
			httpserver.WithLogger(cfg.Logger),
		}, cfg.HTTPServerConfigOpts...)

		cfg.HTTPServer, err = httpserver.New(cfg.PublicKey, defaultHTTPServerEventHandlerFunc(client), cfg.HTTPServerConfigOpts...)
		if err != nil {
			return nil, fmt.Errorf("error while creating http server: %w", err)
		}
	}
	client.HTTPServer = cfg.HTTPServer

	if cfg.MemberChunkingManager == nil {
		cfg.MemberChunkingManager = NewMemberChunkingManager(client, cfg.Logger, cfg.MemberChunkingFilter)
	}
	client.MemberChunkingManager = cfg.MemberChunkingManager

	if cfg.Caches == nil {
		cfg.Caches = cache.New(cfg.CacheConfigOpts...)
	}
	client.Caches = cfg.Caches

	return client, nil
}
