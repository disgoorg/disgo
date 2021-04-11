package api

import (
	"github.com/DisgoOrg/log"

	"github.com/DisgoOrg/disgo/api/endpoints"
)

// DisgoBuilder allows you to create a Disgo client through a series of methods
type DisgoBuilder interface {
	SetLogger(level log.Logger) DisgoBuilder
	SetToken(token endpoints.Token) DisgoBuilder
	SetIntents(intents Intents) DisgoBuilder
	SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) DisgoBuilder
	SetEntityBuilder(entityBuilder EntityBuilder) DisgoBuilder
	SetEventManager(eventManager EventManager) DisgoBuilder
	AddEventListeners(eventsListeners ...EventListener) DisgoBuilder
	SetWebhookServer(webhookServer WebhookServer) DisgoBuilder
	SetWebhookServerProperties(listenURL string, listenPort int, publicKey string) DisgoBuilder
	SetRestClient(restClient RestClient) DisgoBuilder
	SetCache(cache Cache) DisgoBuilder
	SetMemberCachePolicy(memberCachePolicy MemberCachePolicy) DisgoBuilder
	SetMessageCachePolicy(messageCachePolicy MessageCachePolicy) DisgoBuilder
	SetCacheFlags(cacheFlags CacheFlags) DisgoBuilder
	SetGateway(gateway Gateway) DisgoBuilder
	Build() (Disgo, error)
}
