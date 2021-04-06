package internal

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo/api"
)

// NewBuilder returns a new api.DisgoBuilder instance
func NewBuilder(token string) api.DisgoBuilder {
	return DisgoBuilderImpl{
		logLevel:   log.InfoLevel,
		token:      &token,
		cacheFlags: api.CacheFlagsDefault,
	}
}

// DisgoBuilderImpl implementation of the api.DisgoBuilder interface
type DisgoBuilderImpl struct {
	logLevel           log.Level
	token              *string
	gateway            api.Gateway
	restClient         api.RestClient
	cache              api.Cache
	memberCachePolicy  api.MemberCachePolicy
	messageCachePolicy api.MessageCachePolicy
	cacheFlags         api.CacheFlags
	intents            api.Intents
	eventManager       api.EventManager
	webhookServer      api.WebhookServer
	listenURL          *string
	listenPort         *int
	publicKey          *string
	eventListeners     []api.EventListener
}

// SetLogLevel sets logrus.Level of logrus
func (b DisgoBuilderImpl) SetLogLevel(logLevel log.Level) api.DisgoBuilder {
	b.logLevel = logLevel
	return b
}

// SetToken sets the token to connect to discord
func (b DisgoBuilderImpl) SetToken(token string) api.DisgoBuilder {
	b.token = &token
	return b
}

// SetIntents sets the api.Intents to connect to discord
func (b DisgoBuilderImpl) SetIntents(intents api.Intents) api.DisgoBuilder {
	b.intents = intents
	return b
}

// SetEventManager lets you inject your own api.EventManager
func (b DisgoBuilderImpl) SetEventManager(eventManager api.EventManager) api.DisgoBuilder {
	b.eventManager = eventManager
	return b
}

// AddEventListeners lets you add an api.EventListener to your api.EventManager
func (b DisgoBuilderImpl) AddEventListeners(eventListeners ...api.EventListener) api.DisgoBuilder {
	for _, eventListener := range eventListeners {
		b.eventListeners = append(b.eventListeners, eventListener)
	}
	return b
}

// SetWebhookServer lets you inject your own api.EventManager
func (b DisgoBuilderImpl) SetWebhookServer(webhookServer api.WebhookServer) api.DisgoBuilder {
	b.webhookServer = webhookServer
	return b
}

// SetWebhookServerProperties sets the default api.WebhookServer properties
func (b DisgoBuilderImpl) SetWebhookServerProperties(listenURL string, listenPort int, publicKey string) api.DisgoBuilder {
	b.listenURL = &listenURL
	b.listenPort = &listenPort
	b.publicKey = &publicKey
	return b
}

// SetRestClient lets you inject your own api.RestClient
func (b DisgoBuilderImpl) SetRestClient(restClient api.RestClient) api.DisgoBuilder {
	b.restClient = restClient
	return b
}

// SetCache lets you inject your own api.Cache
func (b DisgoBuilderImpl) SetCache(cache api.Cache) api.DisgoBuilder {
	b.cache = cache
	return b
}

// SetMemberCachePolicy lets you set your own api.MemberCachePolicy
func (b DisgoBuilderImpl) SetMemberCachePolicy(memberCachePolicy api.MemberCachePolicy) api.DisgoBuilder {
	b.memberCachePolicy = memberCachePolicy
	return b
}

// SetMessageCachePolicy lets you set your own api.MessageCachePolicy
func (b DisgoBuilderImpl) SetMessageCachePolicy(messageCachePolicy api.MessageCachePolicy) api.DisgoBuilder {
	b.messageCachePolicy = messageCachePolicy
	return b
}

// SetCacheFlags lets you set the api.CacheFlags
func (b DisgoBuilderImpl) SetCacheFlags(cacheFlags api.CacheFlags) api.DisgoBuilder {
	b.cacheFlags = cacheFlags
	return b
}

// EnableCacheFlags lets you enable certain api.CacheFlags
func (b DisgoBuilderImpl) EnableCacheFlags(cacheFlags api.CacheFlags) api.DisgoBuilder {
	b.cacheFlags.Add(cacheFlags)
	return b
}

// DisableCacheFlags lets you disable certain api.CacheFlags
func (b DisgoBuilderImpl) DisableCacheFlags(cacheFlags api.CacheFlags) api.DisgoBuilder {
	b.cacheFlags.Remove(cacheFlags)
	return b
}

// SetGateway lets you inject your own api.Gateway
func (b DisgoBuilderImpl) SetGateway(gateway api.Gateway) api.DisgoBuilder {
	b.gateway = gateway
	return b
}

// Build builds your api.Disgo instance
func (b DisgoBuilderImpl) Build() (api.Disgo, error) {
	log.SetLevel(b.logLevel)

	disgo := &DisgoImpl{}
	if b.token == nil {
		return nil, errors.New("please specify the token")
	}
	disgo.token = *b.token

	id, err := IDFromToken(disgo.token)
	if err != nil {
		log.Errorf("error while getting application id from token: %s", err)
		return nil, err
	}

	disgo.applicationID = *id

	if b.gateway == nil {
		b.gateway = newGatewayImpl(disgo)
	}
	disgo.gateway = b.gateway

	if b.restClient == nil {
		b.restClient = newRestClientImpl(disgo, *b.token)
	}
	disgo.restClient = b.restClient

	disgo.intents = b.intents

	if b.eventManager == nil {
		b.eventManager = newEventManagerImpl(disgo, b.eventListeners)
	}
	disgo.eventManager = b.eventManager

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
