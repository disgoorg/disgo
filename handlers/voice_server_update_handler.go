package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerVoiceServerUpdate handles discord.GatewayEventTypeVoiceServerUpdate
type gatewayHandlerVoiceServerUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerVoiceServerUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceServerUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerVoiceServerUpdate) New() any {
	return &discord.VoiceServerUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceServerUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.VoiceServerUpdate)

	client.EventManager().DispatchEvent(&events.VoiceServerUpdate{
		GenericEvent:      events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceServerUpdate: payload,
	})
}
