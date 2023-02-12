package sharding

import (
	"context"
)

var _ RateLimiter = (*noopRateLimiter)(nil)

// NewNoopRateLimiter creates a new noop RateLimiter.
func NewNoopRateLimiter() RateLimiter {
	return &noopRateLimiter{}
}

type noopRateLimiter struct{}

func (r *noopRateLimiter) Close(_ context.Context)                   {}
func (r *noopRateLimiter) WaitBucket(_ context.Context, _ int) error { return nil }
func (r *noopRateLimiter) UnlockBucket(_ int)                        {}
