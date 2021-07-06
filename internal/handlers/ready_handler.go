package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// ReadyHandler handles api.ReadyGatewayEvent
type ReadyHandler struct{}

// Event returns the raw gateway event Event
func (h *ReadyHandler) Event() api.GatewayEventType {
	return api.GatewayEventReady
}

// New constructs a new payload receiver for the raw gateway event
func (h *ReadyHandler) New() interface{} {
	return &api.ReadyGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ReadyHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	readyEvent, ok := i.(*api.ReadyGatewayEvent)
	if !ok {
		return
	}

	disgo.Cache().CacheUser(&readyEvent.SelfUser)

	for i := range readyEvent.Guilds {
		disgo.Cache().CacheGuild(readyEvent.Guilds[i])
	}

	disgo.EventManager().Dispatch(&events.ReadyEvent{
		GenericEvent:      events.NewGenericEvent(disgo, sequenceNumber),
		ReadyGatewayEvent: readyEvent,
	})

}
