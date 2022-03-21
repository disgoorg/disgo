package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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
func (h *gatewayHandlerVoiceServerUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.VoiceServerUpdate)

	bot.EventManager().Dispatch(&events.VoiceServerUpdateEvent{
		GenericEvent:      events.NewGenericEvent(bot, sequenceNumber),
		VoiceServerUpdate: payload,
	})
}
