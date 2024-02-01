package events

import "github.com/snekROmonoro/disgo/gateway"

type Raw struct {
	*GenericEvent
	gateway.EventRaw
}
