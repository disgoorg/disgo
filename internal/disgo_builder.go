package internal

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

func NewBuilder(token string) api.DisgoBuilder {
	return DisgoBuilderImpl{
		logLevel: log.InfoLevel,
		token: &token,
	}
}

type DisgoBuilderImpl struct {
	logLevel       log.Level
	token          *string
	gateway        api.Gateway
	restClient     api.RestClient
	intents        api.Intent
	eventManager   api.EventManager
	eventListeners []api.EventListener
	cache          api.Cache
}

func (b DisgoBuilderImpl) SetLogLevel(logLevel log.Level) api.DisgoBuilder {
	b.logLevel = logLevel
	return b
}

func (b DisgoBuilderImpl) SetToken(token string) api.DisgoBuilder {
	b.token = &token
	return b
}

func (b DisgoBuilderImpl) SetIntents(intents api.Intent) api.DisgoBuilder {
	b.intents = intents
	return b
}

func (b DisgoBuilderImpl) SetEventManager(eventManager api.EventManager) api.DisgoBuilder {
	b.eventManager = eventManager
	return b
}

func (b DisgoBuilderImpl) AddEventListeners(eventListeners ...api.EventListener) api.DisgoBuilder {
	for _, eventListener := range eventListeners {
		b.eventListeners = append(b.eventListeners, eventListener)
	}
	return b
}

func (b DisgoBuilderImpl) SetRestClient(restClient api.RestClient) api.DisgoBuilder {
	b.restClient = restClient
	return b
}

func (b DisgoBuilderImpl) SetCache(cache api.Cache) api.DisgoBuilder {
	b.cache = cache
	return b
}

func (b DisgoBuilderImpl) SetGateway(gateway api.Gateway) api.DisgoBuilder {
	b.gateway = gateway
	return b
}

func (b DisgoBuilderImpl) Build() (api.Disgo, error) {
	log.SetLevel(b.logLevel)

	disgo := &DisgoImpl{}
	if b.token == nil {
		return nil, errors.New("please specify the token")
	}
	disgo.token = *b.token

	if b.gateway == nil {
		disgo.gateway = newGatewayImpl(disgo)
	} else {
		disgo.gateway = b.gateway
	}

	if b.restClient == nil {
		disgo.restClient = newRestClientImpl(*b.token)
	} else {
		disgo.restClient = b.restClient
	}

	disgo.intents = b.intents

	if b.eventManager == nil {
		disgo.eventManager = newEventManagerImpl(disgo, b.eventListeners)

	} else {
		disgo.eventManager = b.eventManager
	}

	if b.cache == nil {
		disgo.cache = newCacheImpl()
	} else {
		disgo.cache = b.cache
	}

	return disgo, nil
}
