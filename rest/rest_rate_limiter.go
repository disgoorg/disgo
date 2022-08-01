package rest

import (
	"context"
	"net/http"
)

// RateLimiter can be used to supply your own rate limit implementation
type RateLimiter interface {
	// MaxRetries returns the maximum number of retries the client should do
	MaxRetries() int

	// Close gracefully closes the RateLimiter.
	// If the context deadline is exceeded, the RateLimiter will be closed immediately.
	Close(ctx context.Context)

	// Reset resets the rate limiter to its initial state
	Reset()

	// WaitBucket waits for the given bucket to be available for new requests & locks it
	WaitBucket(ctx context.Context, endpoint *CompiledEndpoint) error

	// UnlockBucket unlocks the given bucket and calculates the rate limit for the next request
	UnlockBucket(endpoint *CompiledEndpoint, rs *http.Response) error
}
