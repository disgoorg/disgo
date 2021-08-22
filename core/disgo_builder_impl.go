package core

import (
	"fmt"
	"io"
	"net/http"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/disgo/util"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

// NewBuilder returns a new api.DisgoBuilder instance
func NewBuilder(token string) DisgoBuilder {
	return &DisgoBuilderImpl{
		token: token,
	}
}

// DisgoBuilderImpl implementation of the api.DisgoBuilder interface
type DisgoBuilderImpl struct {
	logger log.Logger

	httpClient        *http.Client
	restClient        rest.Client
	restClientConfig  *rest.Config
	rateLimiter       rate.RateLimiter
	rateLimiterConfig *rate.Config
	restServices      rest.Services

	token string

	eventManager             EventManager
	eventListeners           []EventListener
	rawEventsEnabled         bool
	voiceDispatchInterceptor VoiceDispatchInterceptor

	gateway       gateway.Gateway
	gatewayConfig *gateway.Config

	httpServer       httpserver.Server
	httpServerConfig *httpserver.Config

	cache       Cache
	cacheConfig *CacheConfig

	audioController AudioController
	entityBuilder   EntityBuilder
}

// SetLogger sets logger implementation disgo should use as an _examples logrus
func (b *DisgoBuilderImpl) SetLogger(logger log.Logger) DisgoBuilder {
	b.logger = logger
	return b
}

// SetHTTPClient sets the http.Client rest.Client uses
func (b *DisgoBuilderImpl) SetHTTPClient(httpClient *http.Client) DisgoBuilder {
	b.httpClient = httpClient
	return b
}

// SetRestClient sets the rest.Client rest.Service uses
func (b *DisgoBuilderImpl) SetRestClient(restClient rest.Client) DisgoBuilder {
	b.restClient = restClient
	return b
}
func (b *DisgoBuilderImpl) SetRestClientConfig(config rest.Config) DisgoBuilder {
	b.restClientConfig = &config
	return b
}

// SetRateLimiter sets the rate.RateLimiter the rest.Client uses
func (b *DisgoBuilderImpl) SetRateLimiter(rateLimiter rate.RateLimiter) DisgoBuilder {
	b.rateLimiter = rateLimiter
	return b
}
func (b *DisgoBuilderImpl) SetRateLimiterConfig(config rate.Config) DisgoBuilder {
	b.rateLimiterConfig = &config
	return b
}

// SetRestServices lets you inject your own api.Services
func (b *DisgoBuilderImpl) SetRestServices(restServices rest.Services) DisgoBuilder {
	b.restServices = restServices
	return b
}

// SetEventManager lets you inject your own api.EventManager
func (b *DisgoBuilderImpl) SetEventManager(eventManager EventManager) DisgoBuilder {
	b.eventManager = eventManager
	return b
}

// AddEventListeners lets you add an api.EventListener to your api.EventManager
func (b *DisgoBuilderImpl) AddEventListeners(eventListeners ...EventListener) DisgoBuilder {
	for _, eventListener := range eventListeners {
		b.eventListeners = append(b.eventListeners, eventListener)
	}
	return b
}

// SetRawEventsEnabled enables/disables the events.RawGatewayEvent
func (b *DisgoBuilderImpl) SetRawEventsEnabled(enabled bool) DisgoBuilder {
	b.rawEventsEnabled = enabled
	return b
}

// SetVoiceDispatchInterceptor sets the api.VoiceDispatchInterceptor
func (b *DisgoBuilderImpl) SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) DisgoBuilder {
	b.voiceDispatchInterceptor = voiceDispatchInterceptor
	return b
}

// SetGateway lets you inject your own api.Gateway
func (b *DisgoBuilderImpl) SetGateway(gateway gateway.Gateway) DisgoBuilder {
	b.gateway = gateway
	return b
}

// SetGatewayConfig sets the gateway.Config the gateway.Gateway uses
func (b *DisgoBuilderImpl) SetGatewayConfig(gatewayConfig gateway.Config) DisgoBuilder {
	b.gatewayConfig = &gatewayConfig
	return b
}

// SetHTTPServer lets you inject your own api.EventManager
func (b *DisgoBuilderImpl) SetHTTPServer(httpServer httpserver.Server) DisgoBuilder {
	b.httpServer = httpServer
	return b
}

// SetHTTPServerConfig sets the default api.Server properties
func (b *DisgoBuilderImpl) SetHTTPServerConfig(config httpserver.Config) DisgoBuilder {
	b.httpServerConfig = &config
	return b
}

// SetCache lets you inject your own api.Cache
func (b *DisgoBuilderImpl) SetCache(cache Cache) DisgoBuilder {
	b.cache = cache
	return b
}

// SetCacheConfig lets you inject your own CacheConfig
func (b *DisgoBuilderImpl) SetCacheConfig(config CacheConfig) DisgoBuilder {
	b.cacheConfig = &config
	return b
}

// SetAudioController lets you inject your own api.AudioController
func (b *DisgoBuilderImpl) SetAudioController(audioController AudioController) DisgoBuilder {
	b.audioController = audioController
	return b
}

// SetEntityBuilder lets you inject your own api.EntityBuilder
func (b *DisgoBuilderImpl) SetEntityBuilder(entityBuilder EntityBuilder) DisgoBuilder {
	b.entityBuilder = entityBuilder
	return b
}

// Build builds your api.Disgo instance
func (b *DisgoBuilderImpl) Build() (Disgo, error) {
	disgo := &DisgoImpl{}

	if b.token == "" {
		return nil, discord.ErrNoBotToken
	}
	disgo.token = b.token

	id, err := util.IDFromToken(disgo.token)
	if err != nil {
		disgo.Logger().Errorf("error while getting application id from BotToken: %s", err)
		return nil, err
	}
	// TODO: figure out how we handle different application & client ids
	disgo.applicationID = *id
	disgo.clientID = *id

	if b.logger == nil {
		b.logger = log.Default()
	}
	disgo.logger = b.logger

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	if b.rateLimiter == nil {
		b.rateLimiter = rate.NewRateLimiter(disgo.logger, b.rateLimiterConfig)
	}

	if b.restClientConfig == nil {
		b.restClientConfig = &rest.DefaultConfig
	}

	if b.restClientConfig.Headers == nil {
		b.restClientConfig.Headers = http.Header{}
	}

	if _, ok := b.restClientConfig.Headers["authorization"]; !ok {
		b.restClientConfig.Headers["authorization"] = []string{fmt.Sprintf("Bot %s", b.token)}
	}

	if b.restClient == nil {
		b.restClient = rest.NewClient(b.logger, b.httpClient, b.rateLimiter, b.restClientConfig)
	}

	if b.restServices == nil {
		b.restServices = rest.NewServices(disgo.logger, b.restClient)
	}
	disgo.restServices = b.restServices

	if b.eventManager == nil {
		b.eventManager = NewEventManager(disgo, b.eventListeners)
	}
	disgo.eventManager = b.eventManager

	if b.gatewayConfig != nil {
		if b.gateway == nil {
			b.gateway = gateway.New(disgo.logger, b.restServices, b.token, *b.gatewayConfig, func(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader) {
				disgo.EventManager().HandleGateway(gatewayEventType, sequenceNumber, payload)
			})
		}
		disgo.gateway = b.gateway
	}

	if b.httpServerConfig == nil {
		if b.httpServer == nil {
			b.httpServer = httpserver.New(disgo.logger, *b.httpServerConfig, func(gatewayEventType discord.GatewayEventType, responseChannel chan discord.InteractionResponse, payload io.Reader) {
				disgo.EventManager().HandleHTTP(gatewayEventType, responseChannel, payload)
			})
		}
		disgo.httpServer = b.httpServer
	}

	if b.audioController == nil {
		b.audioController = NewAudioController(disgo)
	}
	disgo.audioController = b.audioController

	if b.entityBuilder == nil {
		b.entityBuilder = NewEntityBuilder(disgo)
	}
	disgo.entityBuilder = b.entityBuilder

	disgo.voiceDispatchInterceptor = b.voiceDispatchInterceptor

	if b.cacheConfig == nil {
		b.cacheConfig = &CacheConfig{
			CacheFlags:         CacheFlagsDefault,
			MemberCachePolicy:  MemberCachePolicyDefault,
			MessageCachePolicy: MessageCachePolicyDefault,
		}
	}

	if b.cache == nil {
		b.cache = NewCache(disgo, *b.cacheConfig)
	}
	disgo.cache = b.cache

	return disgo, nil
}
