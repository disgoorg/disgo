package rrate

import (
	"context"
	"net/http"

	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// Limiter can be used to supply your own rate limit implementation
type Limiter interface {
	// Logger returns the logger the rate limiter uses
	Logger() log.Logger

	// Config returns the Config the rate limiter uses
	Config() Config

	// Close closes the rate limiter and awaits all pending requests to finish. You can use a cancelling context to abort the waiting
	Close(ctx context.Context)

	// WaitBucket waits for the given bucket to be available for new requests & locks it
	WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error

	// UnlockBucket unlocks the given bucket and calculates the rate limit for the next request
	UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error
}
