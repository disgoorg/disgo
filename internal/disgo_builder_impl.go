package internal

import (
	"errors"
	"net/http"

	"github.com/DisgoOrg/log"

	"github.com/DisgoOrg/disgo/api"
)

// NewBuilder returns a new api.DisgoBuilder instance
func NewBuilder(token string) api.DisgoBuilder {
	return &DisgoBuilderImpl{
		token:      token,
		cacheFlags: api.CacheFlagsDefault,
	}
}

// DisgoBuilderImpl implementation of the api.DisgoBuilder interface
type DisgoBuilderImpl struct {
	logger                   log.Logger
	token                    string
	httpClient               *http.Client
	gateway                  api.Gateway
	restClient               api.RestClient
	audioController          api.AudioController
	cache                    api.Cache
	memberCachePolicy        api.MemberCachePolicy
	messageCachePolicy       api.MessageCachePolicy
	cacheFlags               api.CacheFlags
	gatewayIntents           api.GatewayIntents
	rawGatewayEventsEnabled  bool
	entityBuilder            api.EntityBuilder
	eventManager             api.EventManager
	voiceDispatchInterceptor api.VoiceDispatchInterceptor
	webhookServer            api.WebhookServer
	listenURL                *string
	listenPort               *int
	publicKey                *string
	eventListeners           []api.EventListener
}

// SetLogger sets logger implementation disgo should use as an example logrus
func (b *DisgoBuilderImpl) SetLogger(logger log.Logger) api.DisgoBuilder {
	b.logger = logger
	return b
}

// SetToken sets the BotToken to connect to discord
func (b *DisgoBuilderImpl) SetToken(token string) api.DisgoBuilder {
	b.token = token
	return b
}

// SetHTTPClient sets the http.Client the api.RestClient uses
func (b *DisgoBuilderImpl) SetHTTPClient(httpClient *http.Client) api.DisgoBuilder {
	b.httpClient = httpClient
	return b
}

// SetGatewayIntents sets the api.GatewayIntents to connect to discord
func (b *DisgoBuilderImpl) SetGatewayIntents(gatewayIntents ...api.GatewayIntents) api.DisgoBuilder {
	b.gatewayIntents = api.GatewayIntentsNone.Add(gatewayIntents...)
	return b
}

// SetRawGatewayEventsEnabled enables/disables the events.RawGatewayEvent
func (b *DisgoBuilderImpl) SetRawGatewayEventsEnabled(enabled bool) api.DisgoBuilder {
	b.rawGatewayEventsEnabled = enabled
	return b
}

// SetEntityBuilder lets you inject your own api.EntityBuilder
func (b *DisgoBuilderImpl) SetEntityBuilder(entityBuilder api.EntityBuilder) api.DisgoBuilder {
	b.entityBuilder = entityBuilder
	return b
}

// SetEventManager lets you inject your own api.EventManager
func (b *DisgoBuilderImpl) SetEventManager(eventManager api.EventManager) api.DisgoBuilder {
	b.eventManager = eventManager
	return b
}

// AddEventListeners lets you add an api.EventListener to your api.EventManager
func (b *DisgoBuilderImpl) AddEventListeners(eventListeners ...api.EventListener) api.DisgoBuilder {
	for _, eventListener := range eventListeners {
		b.eventListeners = append(b.eventListeners, eventListener)
	}
	return b
}

// SetVoiceDispatchInterceptor sets the api.VoiceDispatchInterceptor
func (b *DisgoBuilderImpl) SetVoiceDispatchInterceptor(voiceDispatchInterceptor api.VoiceDispatchInterceptor) api.DisgoBuilder {
	b.voiceDispatchInterceptor = voiceDispatchInterceptor
	return b
}

// SetWebhookServer lets you inject your own api.EventManager
func (b *DisgoBuilderImpl) SetWebhookServer(webhookServer api.WebhookServer) api.DisgoBuilder {
	b.webhookServer = webhookServer
	return b
}

// SetWebhookServerProperties sets the default api.WebhookServer properties
func (b *DisgoBuilderImpl) SetWebhookServerProperties(listenURL string, listenPort int, publicKey string) api.DisgoBuilder {
	b.listenURL = &listenURL
	b.listenPort = &listenPort
	b.publicKey = &publicKey
	return b
}

// SetRestClient lets you inject your own api.RestClient
func (b *DisgoBuilderImpl) SetRestClient(restClient api.RestClient) api.DisgoBuilder {
	b.restClient = restClient
	return b
}

// SetAudioController lets you inject your own api.AudioController
func (b *DisgoBuilderImpl) SetAudioController(audioController api.AudioController) api.DisgoBuilder {
	b.audioController = audioController
	return b
}

// SetCache lets you inject your own api.Cache
func (b *DisgoBuilderImpl) SetCache(cache api.Cache) api.DisgoBuilder {
	b.cache = cache
	return b
}

// SetMemberCachePolicy lets you set your own api.MemberCachePolicy
func (b *DisgoBuilderImpl) SetMemberCachePolicy(memberCachePolicy api.MemberCachePolicy) api.DisgoBuilder {
	b.memberCachePolicy = memberCachePolicy
	return b
}

// SetMessageCachePolicy lets you set your own api.MessageCachePolicy
func (b *DisgoBuilderImpl) SetMessageCachePolicy(messageCachePolicy api.MessageCachePolicy) api.DisgoBuilder {
	b.messageCachePolicy = messageCachePolicy
	return b
}

// SetCacheFlags lets you set the api.CacheFlags
func (b *DisgoBuilderImpl) SetCacheFlags(cacheFlags ...api.CacheFlags) api.DisgoBuilder {
	b.cacheFlags = api.CacheFlagsNone.Add(cacheFlags...)
	return b
}

// EnableCacheFlags lets you enable certain api.CacheFlags
func (b *DisgoBuilderImpl) EnableCacheFlags(cacheFlags ...api.CacheFlags) api.DisgoBuilder {
	b.cacheFlags.Add(cacheFlags...)
	return b
}

// DisableCacheFlags lets you disable certain api.CacheFlags
func (b *DisgoBuilderImpl) DisableCacheFlags(cacheFlags ...api.CacheFlags) api.DisgoBuilder {
	b.cacheFlags.Remove(cacheFlags...)
	return b
}

// SetGateway lets you inject your own api.Gateway
func (b *DisgoBuilderImpl) SetGateway(gateway api.Gateway) api.DisgoBuilder {
	b.gateway = gateway
	return b
}

// Build builds your api.Disgo instance
func (b *DisgoBuilderImpl) Build() (api.Disgo, error) {

	disgo := &DisgoImpl{
		logger:                  b.logger,
		rawGatewayEventsEnabled: b.rawGatewayEventsEnabled,
	}
	if b.token == "" {
		return nil, errors.New("please specify the BotToken")
	}
	disgo.botToken = b.token

	id, err := IDFromToken(disgo.botToken)
	if err != nil {
		disgo.Logger().Errorf("error while getting application id from BotToken: %s", err)
		return nil, err
	}

	disgo.applicationID = *id
	disgo.clientID = *id

	if b.gateway == nil {
		b.gateway = newGatewayImpl(disgo)
	}
	disgo.gateway = b.gateway

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	if b.restClient == nil {
		b.restClient = newRestClientImpl(disgo, b.httpClient)
	}
	disgo.restClient = b.restClient

	if b.audioController == nil {
		b.audioController = newAudioControllerImpl(disgo)
	}
	disgo.audioController = b.audioController

	disgo.gatewayIntents = b.gatewayIntents

	if b.entityBuilder == nil {
		b.entityBuilder = newEntityBuilderImpl(disgo)
	}
	disgo.entityBuilder = b.entityBuilder

	if b.eventManager == nil {
		b.eventManager = newEventManagerImpl(disgo, b.eventListeners)
	}
	disgo.eventManager = b.eventManager

	disgo.voiceDispatchInterceptor = b.voiceDispatchInterceptor

	if b.webhookServer == nil && b.listenURL != nil && b.listenPort != nil && b.publicKey != nil {
		b.webhookServer = newWebhookServerImpl(disgo, *b.listenURL, *b.listenPort, *b.publicKey)
	}
	disgo.webhookServer = b.webhookServer

	if b.cache == nil {
		if b.memberCachePolicy == nil {
			b.memberCachePolicy = api.MemberCachePolicyDefault
		}
		if b.messageCachePolicy == nil {
			b.messageCachePolicy = api.MessageCachePolicyDefault
		}
		b.cache = newCacheImpl(disgo, b.memberCachePolicy, b.messageCachePolicy, b.cacheFlags)
	}
	disgo.cache = b.cache

	return disgo, nil
}
