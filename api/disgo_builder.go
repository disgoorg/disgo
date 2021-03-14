package api

import log "github.com/sirupsen/logrus"

// DisgoBuilder allows you to create a Disgo client through a series of methods
type DisgoBuilder interface {
	SetLogLevel(level log.Level) DisgoBuilder
	SetToken(string) DisgoBuilder
	SetIntents(Intent) DisgoBuilder
	SetEventManager(EventManager) DisgoBuilder
	AddEventListeners(...EventListener) DisgoBuilder
	SetRestClient(RestClient) DisgoBuilder
	SetCache(Cache) DisgoBuilder
	SetMemberCachePolicy(MemberCachePolicy) DisgoBuilder
	SetGateway(Gateway) DisgoBuilder
	Build() (Disgo, error)
}
