package gateway

import (
	"context"
)

// CommandsPerMinute is the default number of commands per minute that the Gateway will allow.
const CommandsPerMinute = 120

// RateLimiter provides handles the rate limiting logic for connecting to Discord's Gateway.
type RateLimiter interface {
	// Close gracefully closes the RateLimiter.
	// If the context deadline is exceeded, the RateLimiter will be closed immediately.
	Close(ctx context.Context)

	// Reset resets the RateLimiter to its initial state.
	Reset()

	// Wait waits for the RateLimiter to be ready to send a new message.
	// If the context deadline is exceeded, Wait will return immediately and no message will be sent.
	Wait(ctx context.Context) error

	// Unlock unlocks the RateLimiter and allows the next message to be sent.
	Unlock()
}
