package core

import (
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/sharding"
	"github.com/DisgoOrg/log"
)

// NewBotBuilder returns a new BotBuilder instance
func NewBotBuilder(token string) *BotBuilder {
	return &BotBuilder{
		Token: token,
	}
}

// BotBuilder implementation of the BotBuilder interface
type BotBuilder struct {
	Token string
	BotConfig
}

// SetLogger sets the logger implementation disgo should use as a logger
func (b *BotBuilder) SetLogger(logger log.Logger) *BotBuilder {
	b.Logger = logger
	return b
}

// SetHTTPClient sets the http.Client the rest.Client should use
func (b *BotBuilder) SetHTTPClient(httpClient *http.Client) *BotBuilder {
	b.HTTPClient = httpClient
	return b
}

// SetRestClient sets the rest.Client the rest.Service should use
func (b *BotBuilder) SetRestClient(restClient rest.Client) *BotBuilder {
	b.RestClient = restClient
	return b
}

// SetRestClientConfig sets the rest.Config the rest.Client should use
func (b *BotBuilder) SetRestClientConfig(config rest.Config) *BotBuilder {
	b.RestClientConfig = &config
	return b
}

func (b *BotBuilder) SetRestClientConfigOpts(opts ...rest.ConfigOpt) *BotBuilder {
	if b.RestClientConfig == nil {
		b.RestClientConfig = &rest.DefaultConfig
	}
	b.RestClientConfig.Apply(opts)
	return b
}

// SetRestServices lets you inject your own Services
func (b *BotBuilder) SetRestServices(restServices rest.Services) *BotBuilder {
	b.RestServices = restServices
	return b
}

// SetEventManager lets you inject your own EventManager
func (b *BotBuilder) SetEventManager(eventManager EventManager) *BotBuilder {
	b.EventManager = eventManager
	return b
}

// AddEventListeners lets you add an EventListener to your EventManager
func (b *BotBuilder) AddEventListeners(eventListeners ...EventListener) *BotBuilder {
	if b.EventManagerConfig == nil {
		b.EventManagerConfig = &DefaultEventManagerConfig
	}
	for _, eventListener := range eventListeners {
		b.EventManagerConfig.EventListeners = append(b.EventManagerConfig.EventListeners, eventListener)
	}
	return b
}

// SetRawEventsEnabled enables/disables the RawEvent
func (b *BotBuilder) SetRawEventsEnabled(enabled bool) *BotBuilder {
	if b.EventManagerConfig == nil {
		b.EventManagerConfig = &DefaultEventManagerConfig
	}
	b.EventManagerConfig.RawEventsEnabled = enabled
	return b
}

// SetVoiceDispatchInterceptor sets the VoiceDispatchInterceptor
func (b *BotBuilder) SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) *BotBuilder {
	if b.EventManagerConfig == nil {
		b.EventManagerConfig = &DefaultEventManagerConfig
	}
	b.EventManagerConfig.VoiceDispatchInterceptor = voiceDispatchInterceptor
	return b
}

// SetGateway lets you inject your own gateway.Gateway
func (b *BotBuilder) SetGateway(gateway gateway.Gateway) *BotBuilder {
	b.Gateway = gateway
	return b
}

// SetGatewayConfig sets the gateway.Config the gateway.Gateway should use
func (b *BotBuilder) SetGatewayConfig(gatewayConfig gateway.Config) *BotBuilder {
	b.GatewayConfig = &gatewayConfig
	return b
}

func (b *BotBuilder) SetGatewayConfigOpts(opts ...gateway.ConfigOpt) *BotBuilder {
	if b.GatewayConfig == nil {
		b.GatewayConfig = &gateway.DefaultConfig
	}
	b.GatewayConfig.Apply(opts)
	return b
}

// SetShardManager lets you inject your own sharding.ShardManager
func (b *BotBuilder) SetShardManager(shardManager sharding.ShardManager) *BotBuilder {
	b.ShardManager = shardManager
	return b
}

// SetShardManagerConfig sets the sharding.Config the sharding.ShardManager should use
func (b *BotBuilder) SetShardManagerConfig(shardManagerConfig sharding.Config) *BotBuilder {
	b.ShardManagerConfig = &shardManagerConfig
	return b
}

func (b *BotBuilder) SetShardMangerConfigOpts(opts ...sharding.ConfigOpt) *BotBuilder {
	if b.ShardManagerConfig == nil {
		b.ShardManagerConfig = &sharding.DefaultConfig
	}
	b.ShardManagerConfig.Apply(opts)
	return b
}

// SetHTTPServer lets you inject your own EventManager
func (b *BotBuilder) SetHTTPServer(httpServer httpserver.Server) *BotBuilder {
	b.HTTPServer = httpServer
	return b
}

// SetHTTPServerConfig sets the default httpserver.Server properties
func (b *BotBuilder) SetHTTPServerConfig(config httpserver.Config) *BotBuilder {
	b.HTTPServerConfig = &config
	return b
}

func (b *BotBuilder) SetHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) *BotBuilder {
	if b.HTTPServerConfig == nil {
		b.HTTPServerConfig = &httpserver.DefaultConfig
	}
	b.HTTPServerConfig.Apply(opts)
	return b
}

// SetCache lets you inject your own Caches
func (b *BotBuilder) SetCache(cache Caches) *BotBuilder {
	b.Caches = cache
	return b
}

// SetCacheConfig lets you inject your own CacheConfig
func (b *BotBuilder) SetCacheConfig(config CacheConfig) *BotBuilder {
	b.CacheConfig = &config
	return b
}

func (b *BotBuilder) SetCacheConfigOpts(opts ...CacheConfigOpt) *BotBuilder {
	if b.CacheConfig == nil {
		b.CacheConfig = &DefaultCacheConfig
	}
	b.CacheConfig.Apply(opts)
	return b
}

// SetAudioController lets you inject your own AudioController
func (b *BotBuilder) SetAudioController(audioController AudioController) *BotBuilder {
	b.AudioController = audioController
	return b
}

// SetMembersChunkingManager lets you inject your own MembersChunkingManager
func (b *BotBuilder) SetMembersChunkingManager(membersChunkingManager MembersChunkingManager) *BotBuilder {
	b.MembersChunkingManager = membersChunkingManager
	return b
}

// SetEntityBuilder lets you inject your own EntityBuilder
func (b *BotBuilder) SetEntityBuilder(entityBuilder EntityBuilder) *BotBuilder {
	b.EntityBuilder = entityBuilder
	return b
}

// Build builds your Bot instance
func (b *BotBuilder) Build() (*Bot, error) {
	return buildBot(b.Token, b.BotConfig)
}
