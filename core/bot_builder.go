package core

import (
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/log"
)

// NewBotBuilder returns a new api.BotBuilder instance
func NewBotBuilder(token string) *BotBuilder {
	return &BotBuilder{
		Token: token,
	}
}

// BotBuilder implementation of the api.BotBuilder interface
type BotBuilder struct {
	Token string
	BotConfig
}

// SetLogger sets logger implementation disgo should use as an _examples logrus
func (b *BotBuilder) SetLogger(logger log.Logger) *BotBuilder {
	b.Logger = logger
	return b
}

// SetHTTPClient sets the http.Client rest.Client uses
func (b *BotBuilder) SetHTTPClient(httpClient *http.Client) *BotBuilder {
	b.HTTPClient = httpClient
	return b
}

// SetRestClient sets the rest.Client rest.Service uses
func (b *BotBuilder) SetRestClient(restClient rest.Client) *BotBuilder {
	b.RestClient = restClient
	return b
}
func (b *BotBuilder) SetRestClientConfig(config rest.Config) *BotBuilder {
	b.RestClientConfig = &config
	return b
}

func (b *BotBuilder) SetRestClientConfigOpts(opts ...rest.ConfigOpt) *BotBuilder {
	b.RestClientConfig.Apply(opts)
	return b
}

// SetRateLimiter sets the rate.Limiter the rest.Client uses
func (b *BotBuilder) SetRateLimiter(rateLimiter rate.Limiter) *BotBuilder {
	b.RateLimiter = rateLimiter
	return b
}
func (b *BotBuilder) SetRateLimiterConfig(config rate.Config) *BotBuilder {
	b.RateLimiterConfig = &config
	return b
}

func (b *BotBuilder) SetRateLimiterConfigOpts(opts ...rate.ConfigOpt) *BotBuilder {
	b.RateLimiterConfig.Apply(opts)
	return b
}

// SetRestServices lets you inject your own api.Services
func (b *BotBuilder) SetRestServices(restServices rest.Services) *BotBuilder {
	b.RestServices = restServices
	return b
}

// SetEventManager lets you inject your own api.EventManager
func (b *BotBuilder) SetEventManager(eventManager EventManager) *BotBuilder {
	b.EventManager = eventManager
	return b
}

// AddEventListeners lets you add an api.EventListener to your api.EventManager
func (b *BotBuilder) AddEventListeners(eventListeners ...EventListener) *BotBuilder {
	for _, eventListener := range eventListeners {
		b.EventListeners = append(b.EventListeners, eventListener)
	}
	return b
}

// SetRawEventsEnabled enables/disables the events.RawGatewayEvent
func (b *BotBuilder) SetRawEventsEnabled(enabled bool) *BotBuilder {
	b.RawEventsEnabled = enabled
	return b
}

// SetVoiceDispatchInterceptor sets the api.VoiceDispatchInterceptor
func (b *BotBuilder) SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) *BotBuilder {
	b.VoiceDispatchInterceptor = voiceDispatchInterceptor
	return b
}

// SetGateway lets you inject your own api.Gateway
func (b *BotBuilder) SetGateway(gateway gateway.Gateway) *BotBuilder {
	b.Gateway = gateway
	return b
}

// SetGatewayConfig sets the gateway.Config the gateway.Gateway uses
func (b *BotBuilder) SetGatewayConfig(gatewayConfig gateway.Config) *BotBuilder {
	b.GatewayConfig = &gatewayConfig
	return b
}

func (b *BotBuilder) SetGatewayConfigOpts(opts ...gateway.ConfigOpt) *BotBuilder {
	b.GatewayConfig.Apply(opts)
	return b
}

// SetHTTPServer lets you inject your own api.EventManager
func (b *BotBuilder) SetHTTPServer(httpServer httpserver.Server) *BotBuilder {
	b.HTTPServer = httpServer
	return b
}

// SetHTTPServerConfig sets the default api.Server properties
func (b *BotBuilder) SetHTTPServerConfig(config httpserver.Config) *BotBuilder {
	b.HTTPServerConfig = &config
	return b
}

func (b *BotBuilder) SetHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) *BotBuilder {
	b.HTTPServerConfig.Apply(opts)
	return b
}

// SetCache lets you inject your own api.Caches
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
	b.CacheConfig.Apply(opts)
	return b
}

// SetAudioController lets you inject your own api.AudioController
func (b *BotBuilder) SetAudioController(audioController AudioController) *BotBuilder {
	b.AudioController = audioController
	return b
}

// SetEntityBuilder lets you inject your own api.EntityBuilder
func (b *BotBuilder) SetEntityBuilder(entityBuilder EntityBuilder) *BotBuilder {
	b.EntityBuilder = entityBuilder
	return b
}

// Build builds your api.Bot instance
func (b *BotBuilder) Build() (*Bot, error) {
	return buildBot(b.Token, b.BotConfig)
}
