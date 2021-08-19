package core

import (
	"net/http"

	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/rate"
	"github.com/DisgoOrg/log"
)

// DisgoBuilder allows you to create a Disgo client through a series of methods
type DisgoBuilder interface {
	SetLogger(level log.Logger) DisgoBuilder

	SetHTTPClient(httpClient *http.Client) DisgoBuilder
	SetRestHTTPClient(restHTTPClient rest.HTTPClient) DisgoBuilder
	SetRateLimiter(rateLimiter rate.RateLimiter) DisgoBuilder
	SetRestServices(restServices rest.Services) DisgoBuilder

	SetEventManager(eventManager EventManager) DisgoBuilder
	AddEventListeners(eventsListeners ...EventListener) DisgoBuilder
	SetRawEventsEnabled(enabled bool) DisgoBuilder
	SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) DisgoBuilder

	SetGateway(gateway gateway.Gateway) DisgoBuilder
	SetGatewayConfig(config gateway.Config) DisgoBuilder

	SetHTTPServer(httpServer httpserver.Server) DisgoBuilder
	SetHTTPServerConfig(config httpserver.Config) DisgoBuilder

	SetCache(cache Cache) DisgoBuilder
	SetCacheConfig(config CacheConfig) DisgoBuilder

	SetAudioController(audioController AudioController) DisgoBuilder
	SetEntityBuilder(entityBuilder EntityBuilder) DisgoBuilder

	Build() (Disgo, error)
}
