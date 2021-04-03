package api

import log "github.com/sirupsen/logrus"

// DisgoBuilder allows you to create a Disgo client through a series of methods
type DisgoBuilder interface {
	SetLogLevel(level log.Level) DisgoBuilder
	SetToken(token string) DisgoBuilder
	SetIntents(intents Intents) DisgoBuilder
	SetEventManager(eventManager EventManager) DisgoBuilder
	SetVoiceDispatchInterceptor(VoiceDispatchInterceptor) DisgoBuilder
	AddEventListeners(eventsListeners ...EventListener) DisgoBuilder
	SetWebhookServer(webhookServer WebhookServer) DisgoBuilder
	SetWebhookServerProperties(listenURL string, listenPort int, publicKey string) DisgoBuilder
	SetRestClient(restClient RestClient) DisgoBuilder
	SetCache(cache Cache) DisgoBuilder
	SetMessageCachePolicy(messageCachePolicy MessageCachePolicy) DisgoBuilder
	SetMemberCachePolicy(memberCachePolicy MemberCachePolicy) DisgoBuilder
	SetCacheVoiceStates(cacheVoiceStates bool) DisgoBuilder
	SetCacheRoles(cacheRoles bool) DisgoBuilder
	SetCacheChannels(cacheChannels bool) DisgoBuilder
	SetCacheEmotes(cacheEmotes bool) DisgoBuilder
	SetGateway(gateway Gateway) DisgoBuilder
	Build() (Disgo, error)
}
