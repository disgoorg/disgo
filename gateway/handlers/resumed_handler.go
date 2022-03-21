package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerReady handles discord.GatewayEventTypeResumed
type gatewayHandlerResumed struct{}

// EventType returns the gateway.EventType
func (h *gatewayHandlerResumed) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeResumed
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerResumed) New() any {
	return nil
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerResumed) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, _ any) {
	client.EventManager().Dispatch(&events.ResumedEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
	})

}
