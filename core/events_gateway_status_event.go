package core

import (
	"github.com/DisgoOrg/disgo/gateway"
)

// GenericGatewayStatusEvent is called upon receiving ConnectedEvent, ReconnectedEvent, ResumedEvent, DisconnectedEvent or ShutdownEvent
type GenericGatewayStatusEvent struct {
	*GenericEvent
	Status gateway.Status
}

// ConnectedEvent indicates disgo connected to the core.Gateway
type ConnectedEvent struct {
	*GenericGatewayStatusEvent
}

// ReconnectedEvent indicates disgo reconnected to the core.Gateway
type ReconnectedEvent struct {
	*GenericGatewayStatusEvent
}

// ResumedEvent indicates disgo resumed to the core.Gateway
type ResumedEvent struct {
	*GenericGatewayStatusEvent
}

// DisconnectedEvent indicates disgo disconnected to the core.Gateway
type DisconnectedEvent struct {
	*GenericGatewayStatusEvent
}
