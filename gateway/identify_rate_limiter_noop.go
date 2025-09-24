package gateway

import (
	"context"
)

var _ IdentifyRateLimiter = (*noopIdentifyRateLimiter)(nil)

// NewNoopIdentifyRateLimiter creates a new noop RateLimiter.
func NewNoopIdentifyRateLimiter() IdentifyRateLimiter {
	return &noopIdentifyRateLimiter{}
}

type noopIdentifyRateLimiter struct{}

func (r *noopIdentifyRateLimiter) Close(_ context.Context)             {}
func (r *noopIdentifyRateLimiter) Wait(_ context.Context, _ int) error { return nil }
func (r *noopIdentifyRateLimiter) Unlock(_ int)                        {}
