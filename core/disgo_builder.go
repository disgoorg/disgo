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

	SetRestClient(restHTTPClient rest.Client) DisgoBuilder
	SetRestClientConfig(config rest.Config) DisgoBuilder
	SetRestClientConfigOpts(opts ...rest.ConfigOpt) DisgoBuilder

	SetRateLimiter(rateLimiter rate.Limiter) DisgoBuilder
	SetRateLimiterConfig(config rate.Config) DisgoBuilder
	SetRateLimiterConfigOpts(opts ...rate.ConfigOpt) DisgoBuilder

	SetRestServices(restServices rest.Services) DisgoBuilder

	SetEventManager(eventManager EventManager) DisgoBuilder
	AddEventListeners(eventsListeners ...EventListener) DisgoBuilder
	SetRawEventsEnabled(enabled bool) DisgoBuilder
	SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) DisgoBuilder

	SetGateway(gateway gateway.Gateway) DisgoBuilder
	SetGatewayConfig(config gateway.Config) DisgoBuilder
	SetGatewayConfigOpts(opts ...gateway.ConfigOpt) DisgoBuilder

	SetHTTPServer(httpServer httpserver.Server) DisgoBuilder
	SetHTTPServerConfig(config httpserver.Config) DisgoBuilder
	SetHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) DisgoBuilder

	SetCache(cache Caches) DisgoBuilder
	SetCacheConfig(config CacheConfig) DisgoBuilder
	SetCacheConfigOpts(opts ...CacheConfigOpt) DisgoBuilder

	SetAudioController(audioController AudioController) DisgoBuilder
	SetEntityBuilder(entityBuilder EntityBuilder) DisgoBuilder

	Build() (Disgo, error)
}
