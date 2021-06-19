package api

import (
	"net/http"

	"github.com/DisgoOrg/log"
)

// DisgoBuilder allows you to create a Disgo client through a series of methods
type DisgoBuilder interface {
	SetLogger(level log.Logger) DisgoBuilder
	SetToken(token string) DisgoBuilder
	SetHTTPClient(httpClient *http.Client) DisgoBuilder
	SetGatewayIntents(GatewayIntents ...GatewayIntents) DisgoBuilder
	SetRawGatewayEventsEnabled(enabled bool) DisgoBuilder
	SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) DisgoBuilder
	SetEntityBuilder(entityBuilder EntityBuilder) DisgoBuilder
	SetEventManager(eventManager EventManager) DisgoBuilder
	AddEventListeners(eventsListeners ...EventListener) DisgoBuilder
	SetWebhookServer(webhookServer WebhookServer) DisgoBuilder
	SetWebhookServerProperties(listenURL string, listenPort int, publicKey string) DisgoBuilder
	SetRestClient(restClient RestClient) DisgoBuilder
	SetCache(cache Cache) DisgoBuilder
	SetMemberCachePolicy(memberCachePolicy MemberCachePolicy) DisgoBuilder
	SetThreadMemberCachePolicy(threadMemberCachePolicy ThreadMemberCachePolicy) DisgoBuilder
	SetMessageCachePolicy(messageCachePolicy MessageCachePolicy) DisgoBuilder
	SetCacheFlags(cacheFlags ...CacheFlags) DisgoBuilder
	EnableCacheFlags(cacheFlags ...CacheFlags) DisgoBuilder
	DisableCacheFlags(cacheFlags ...CacheFlags) DisgoBuilder
	SetGateway(gateway Gateway) DisgoBuilder
	Build() (Disgo, error)
}
