package core

import (
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/log"
)

// NewBuilder returns a new api.DisgoBuilder instance
func NewBuilder(token string) DisgoBuilder {
	return &DisgoBuilderImpl{
		DisgoConfig: DisgoConfig{Token: token},
	}
}

// DisgoBuilderImpl implementation of the api.DisgoBuilder interface
type DisgoBuilderImpl struct {
	DisgoConfig
}

// SetLogger sets logger implementation disgo should use as an _examples logrus
func (b *DisgoBuilderImpl) SetLogger(logger log.Logger) DisgoBuilder {
	b.Logger = logger
	return b
}

// SetHTTPClient sets the http.Client rest.Client uses
func (b *DisgoBuilderImpl) SetHTTPClient(httpClient *http.Client) DisgoBuilder {
	b.HTTPClient = httpClient
	return b
}

// SetRestClient sets the rest.Client rest.Service uses
func (b *DisgoBuilderImpl) SetRestClient(restClient rest.Client) DisgoBuilder {
	b.RestClient = restClient
	return b
}
func (b *DisgoBuilderImpl) SetRestClientConfig(config rest.Config) DisgoBuilder {
	b.RestClientConfig = &config
	return b
}

func (b *DisgoBuilderImpl) SetRestClientConfigOpts(opts ...rest.ConfigOpt) DisgoBuilder {
	b.RestClientConfig.Apply(opts)
	return b
}

// SetRateLimiter sets the rate.Limiter the rest.Client uses
func (b *DisgoBuilderImpl) SetRateLimiter(rateLimiter rate.Limiter) DisgoBuilder {
	b.RateLimiter = rateLimiter
	return b
}
func (b *DisgoBuilderImpl) SetRateLimiterConfig(config rate.Config) DisgoBuilder {
	b.RateLimiterConfig = &config
	return b
}

func (b *DisgoBuilderImpl) SetRateLimiterConfigOpts(opts ...rate.ConfigOpt) DisgoBuilder {
	b.RateLimiterConfig.Apply(opts)
	return b
}

// SetRestServices lets you inject your own api.Services
func (b *DisgoBuilderImpl) SetRestServices(restServices rest.Services) DisgoBuilder {
	b.RestServices = restServices
	return b
}

// SetEventManager lets you inject your own api.EventManager
func (b *DisgoBuilderImpl) SetEventManager(eventManager EventManager) DisgoBuilder {
	b.EventManager = eventManager
	return b
}

// AddEventListeners lets you add an api.EventListener to your api.EventManager
func (b *DisgoBuilderImpl) AddEventListeners(eventListeners ...EventListener) DisgoBuilder {
	for _, eventListener := range eventListeners {
		b.EventListeners = append(b.EventListeners, eventListener)
	}
	return b
}

// SetRawEventsEnabled enables/disables the events.RawGatewayEvent
func (b *DisgoBuilderImpl) SetRawEventsEnabled(enabled bool) DisgoBuilder {
	b.RawEventsEnabled = enabled
	return b
}

// SetVoiceDispatchInterceptor sets the api.VoiceDispatchInterceptor
func (b *DisgoBuilderImpl) SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) DisgoBuilder {
	b.VoiceDispatchInterceptor = voiceDispatchInterceptor
	return b
}

// SetGateway lets you inject your own api.Gateway
func (b *DisgoBuilderImpl) SetGateway(gateway gateway.Gateway) DisgoBuilder {
	b.Gateway = gateway
	return b
}

// SetGatewayConfig sets the gateway.Config the gateway.Gateway uses
func (b *DisgoBuilderImpl) SetGatewayConfig(gatewayConfig gateway.Config) DisgoBuilder {
	b.GatewayConfig = &gatewayConfig
	return b
}

func (b *DisgoBuilderImpl) SetGatewayConfigOpts(opts ...gateway.ConfigOpt) DisgoBuilder {
	b.GatewayConfig.Apply(opts)
	return b
}

// SetHTTPServer lets you inject your own api.EventManager
func (b *DisgoBuilderImpl) SetHTTPServer(httpServer httpserver.Server) DisgoBuilder {
	b.HTTPServer = httpServer
	return b
}

// SetHTTPServerConfig sets the default api.Server properties
func (b *DisgoBuilderImpl) SetHTTPServerConfig(config httpserver.Config) DisgoBuilder {
	b.HTTPServerConfig = &config
	return b
}

func (b *DisgoBuilderImpl) SetHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) DisgoBuilder {
	b.HTTPServerConfig.Apply(opts)
	return b
}

// SetCache lets you inject your own api.Caches
func (b *DisgoBuilderImpl) SetCache(cache Caches) DisgoBuilder {
	b.Caches = cache
	return b
}

// SetCacheConfig lets you inject your own CacheConfig
func (b *DisgoBuilderImpl) SetCacheConfig(config CacheConfig) DisgoBuilder {
	b.CacheConfig = &config
	return b
}

func (b *DisgoBuilderImpl) SetCacheConfigOpts(opts ...CacheConfigOpt) DisgoBuilder {
	b.CacheConfig.Apply(opts)
	return b
}

// SetAudioController lets you inject your own api.AudioController
func (b *DisgoBuilderImpl) SetAudioController(audioController AudioController) DisgoBuilder {
	b.AudioController = audioController
	return b
}

// SetEntityBuilder lets you inject your own api.EntityBuilder
func (b *DisgoBuilderImpl) SetEntityBuilder(entityBuilder EntityBuilder) DisgoBuilder {
	b.EntityBuilder = entityBuilder
	return b
}

// Build builds your api.Disgo instance
func (b *DisgoBuilderImpl) Build() (Disgo, error) {
	return buildDisgoImpl(b.DisgoConfig)
}
