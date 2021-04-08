package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericGatewayStatusEvent struct {
	GenericEvent
	Status api.GatewayStatus
}

type ConnectedEvent struct {
	GenericGatewayStatusEvent
}

type ReconnectedEvent struct {
	GenericGatewayStatusEvent
}

type ResumedEvent struct {
	GenericGatewayStatusEvent
}

type DisconnectedEvent struct {
	GenericGatewayStatusEvent
}

type ShutdownEvent struct {
	GenericGatewayStatusEvent
}


