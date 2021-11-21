package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerVoiceServerUpdate handles core.GatewayEventVoiceServerUpdate
type gatewayHandlerVoiceServerUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerVoiceServerUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceServerUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerVoiceServerUpdate) New() interface{} {
	return &discord.VoiceServerUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceServerUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.VoiceServerUpdate)

	bot.EventManager.Dispatch(&events.VoiceServerUpdateEvent{
		GenericEvent:      events.NewGenericEvent(bot, sequenceNumber),
		VoiceServerUpdate: payload,
	})
}
