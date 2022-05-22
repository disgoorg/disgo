package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
func (h *gatewayHandlerResumed) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, _ any) {
	client.EventManager().DispatchEvent(&events.Resumed{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
	})

}
