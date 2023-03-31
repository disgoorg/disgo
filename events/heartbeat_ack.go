package events

import "github.com/disgoorg/disgo/gateway"

type HeartbeatAck struct {
	*GenericEvent
	gateway.EventHeartbeatAck
}
