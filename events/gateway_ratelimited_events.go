package events

import "github.com/disgoorg/disgo/gateway"

// GatewayRateLimited indicates we received a RateLimited from the gateway.Gateway
type GatewayRateLimited struct {
	*GenericEvent
	gateway.EventRateLimited
}
