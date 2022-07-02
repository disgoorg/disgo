package events

import "github.com/disgoorg/disgo/gateway"

type Raw struct {
	*GenericEvent
	gateway.EventRaw
}
