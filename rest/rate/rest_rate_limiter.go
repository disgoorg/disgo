package rate

import (
	"context"
	"net/http"

	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

type Limiter interface {
	Logger() log.Logger
	Close(ctx context.Context)
	Config() Config
	WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error
	UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error
}
