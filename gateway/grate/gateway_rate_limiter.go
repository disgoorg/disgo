package grate

import (
	"context"

	"github.com/disgoorg/log"
)

// Limiter provides handles the rate limiting logic for connecting to Discord's Gateway.
type Limiter interface {
	// Logger returns the logger used by the Limiter.
	Logger() log.Logger

	// Close gracefully closes the Limiter.
	// If the context deadline is exceeded, the Limiter will be closed immediately.
	Close(ctx context.Context)

	// Reset resets the Limiter to its initial state.
	Reset()

	// Wait waits for the Limiter to be ready to send a new message.
	// If the context deadline is exceeded, Wait will return immediately and no message will be sent.
	Wait(ctx context.Context) error

	// Unlock unlocks the Limiter and allows the next message to be sent.
	Unlock()
}
