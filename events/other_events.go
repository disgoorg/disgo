package events

import (
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/webhookevent"
)

type HeartbeatAck struct {
	*Event
	*GatewayEvent
	gateway.EventHeartbeatAck
}

type WebhookPing struct {
	*Event
	*WebhookEvent
}

type GatewayRaw struct {
	*Event
	*GatewayEvent
	gateway.EventRaw
}

type WebhookRaw struct {
	*Event
	*WebhookEvent
	webhookevent.EventDataRaw
}
