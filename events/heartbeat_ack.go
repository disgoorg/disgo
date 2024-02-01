package events

import "github.com/snekROmonoro/disgo/gateway"

type HeartbeatAck struct {
	*GenericEvent
	gateway.EventHeartbeatAck
}
