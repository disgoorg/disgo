package rest

import (
	"context"
	"net/http"
)

// NewNoopRateLimiter return a new noop RateLimiter.
func NewNoopRateLimiter() RateLimiter {
	return &noopRateLimiter{}
}

type noopRateLimiter struct{}

func (l *noopRateLimiter) MaxRetries() int { return 0 }

func (l *noopRateLimiter) Close(_ context.Context) {}

func (l *noopRateLimiter) Reset() {}

func (l *noopRateLimiter) WaitBucket(_ context.Context, _ *CompiledEndpoint) error { return nil }

func (l *noopRateLimiter) UnlockBucket(_ *CompiledEndpoint, _ *http.Response) error { return nil }
