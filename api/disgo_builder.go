package api

import log "github.com/sirupsen/logrus"

type DisgoBuilder interface {
	SetLogLevel(level log.Level) DisgoBuilder
	SetToken(string) DisgoBuilder
	SetIntents(Intent) DisgoBuilder
	SetEventManager(EventManager) DisgoBuilder
	AddEventListeners(...EventListener) DisgoBuilder
	SetRestClient(RestClient) DisgoBuilder
	SetCache(Cache) DisgoBuilder
	SetGateway(Gateway) DisgoBuilder
	Build() (Disgo, error)
}
