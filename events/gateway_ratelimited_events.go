package events

import "github.com/disgoorg/disgo/gateway"

// RateLimited indicates we received a RateLimited from the gateway.Gateway
type RateLimited struct {
	*GenericEvent
	gateway.EventRateLimited
}
