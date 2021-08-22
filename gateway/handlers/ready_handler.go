package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// ReadyHandler handles discord.GatewayEventTypeReady
type ReadyHandler struct{}

// EventType returns the gateway.EventType
func (h *ReadyHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeReady
}

// New constructs a new payload receiver for the raw gateway event
func (h *ReadyHandler) New() interface{} {
	return &discord.ReadyGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ReadyHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	readyEvent, ok := i.(*discord.ReadyGatewayEvent)
	if !ok {
		return
	}

	disgo.EntityBuilder().CreateSelfUser(readyEvent.SelfUser, core.CacheStrategyYes)

	for _, guild := range readyEvent.Guilds {
		disgo.EntityBuilder().CreateGuild(guild, core.CacheStrategyYes)
	}

	eventManager.Dispatch(&events.ReadyEvent{
		GenericEvent:      events.NewGenericEvent(disgo, sequenceNumber),
		ReadyGatewayEvent: readyEvent,
	})

}
