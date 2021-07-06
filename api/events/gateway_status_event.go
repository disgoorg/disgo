package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericGatewayStatusEvent is called upon receiving ConnectedEvent, ReconnectedEvent, ResumedEvent, DisconnectedEvent or ShutdownEvent
type GenericGatewayStatusEvent struct {
	*GenericEvent
	Status api.GatewayStatus
}

// ConnectedEvent indicates disgo connected to the api.Gateway
type ConnectedEvent struct {
	*GenericGatewayStatusEvent
}

// ReconnectedEvent indicates disgo reconnected to the api.Gateway
type ReconnectedEvent struct {
	*GenericGatewayStatusEvent
}

// ResumedEvent indicates disgo resumed to the api.Gateway
type ResumedEvent struct {
	*GenericGatewayStatusEvent
}

// DisconnectedEvent indicates disgo disconnected to the api.Gateway
type DisconnectedEvent struct {
	*GenericGatewayStatusEvent
}
