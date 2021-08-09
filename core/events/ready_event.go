package events

import (
	"github.com/DisgoOrg/disgo/gateway"
)

// ReadyEvent indicates we received the ReadyEvent from the api.Gateway
type ReadyEvent struct {
	*GenericEvent
	*gateway.ReadyGatewayEvent
}
