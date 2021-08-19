package rate

import (
	"context"
	"errors"
	"net/http"

	"github.com/DisgoOrg/disgo/rest/route"
)

var ErrCtxTimeout = errors.New("rate limit exceeds context deadline")

//goland:noinspection GoNameStartsWithPackageName
type RateLimiter interface {
	Close(ctx context.Context)
	WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error
	UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error
}
